package autoscaling

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/google/uuid"
)

// AutoScaler manages auto-scaling operations
type AutoScaler struct {
	policies   map[string]*ScalingPolicy
	metrics    MetricsProvider
	executor   ScalingExecutor
	config     *AutoScalerConfig
	mu         sync.RWMutex
	stopCh     chan struct{}
}

// AutoScalerConfig holds auto-scaler configuration
type AutoScalerConfig struct {
	EvaluationInterval time.Duration
	CooldownPeriod     time.Duration
	MetricsWindow      time.Duration
	MaxScaleUpRate     float64 // Maximum scale up rate per evaluation
	MaxScaleDownRate   float64 // Maximum scale down rate per evaluation
	StabilizationWindow time.Duration
}

// ScalingPolicy defines scaling rules for a service
type ScalingPolicy struct {
	ID              uuid.UUID            `json:"id"`
	Name            string               `json:"name"`
	ServiceName     string               `json:"service_name"`
	MinReplicas     int                  `json:"min_replicas"`
	MaxReplicas     int                  `json:"max_replicas"`
	TargetCPU       float64              `json:"target_cpu"`        // Target CPU utilization %
	TargetMemory    float64              `json:"target_memory"`     // Target memory utilization %
	TargetRPS       float64              `json:"target_rps"`        // Target requests per second
	CustomMetrics   []CustomMetricTarget `json:"custom_metrics"`    // Custom metric targets
	Behaviors       *ScalingBehavior     `json:"behaviors"`         // Scaling behaviors
	Enabled         bool                 `json:"enabled"`
	CreatedAt       time.Time            `json:"created_at"`
	UpdatedAt       time.Time            `json:"updated_at"`
	LastScaledAt    time.Time            `json:"last_scaled_at"`
	CurrentReplicas int                  `json:"current_replicas"`
}

// CustomMetricTarget defines custom metric scaling target
type CustomMetricTarget struct {
	Name       string  `json:"name"`
	Target     float64 `json:"target"`
	Type       string  `json:"type"` // "value", "utilization", "average"
	Weight     float64 `json:"weight"` // Weight in scaling decision
}

// ScalingBehavior defines scaling behavior rules
type ScalingBehavior struct {
	ScaleUp   *ScaleUpBehavior   `json:"scale_up"`
	ScaleDown *ScaleDownBehavior `json:"scale_down"`
}

// ScaleUpBehavior defines scale-up behavior
type ScaleUpBehavior struct {
	StabilizationWindow time.Duration    `json:"stabilization_window"`
	Policies            []ScalingRulePolicy `json:"policies"`
}

// ScaleDownBehavior defines scale-down behavior
type ScaleDownBehavior struct {
	StabilizationWindow time.Duration    `json:"stabilization_window"`
	Policies            []ScalingRulePolicy `json:"policies"`
}

// ScalingRulePolicy defines a scaling rule policy
type ScalingRulePolicy struct {
	Type          string        `json:"type"`          // "Pods", "Percent"
	Value         int           `json:"value"`         // Number of pods or percentage
	PeriodSeconds time.Duration `json:"period_seconds"` // Period for this rule
}

// MetricsProvider provides metrics for scaling decisions
type MetricsProvider interface {
	GetCPUUtilization(ctx context.Context, serviceName string, window time.Duration) (float64, error)
	GetMemoryUtilization(ctx context.Context, serviceName string, window time.Duration) (float64, error)
	GetRequestsPerSecond(ctx context.Context, serviceName string, window time.Duration) (float64, error)
	GetCustomMetric(ctx context.Context, serviceName, metricName string, window time.Duration) (float64, error)
	GetCurrentReplicas(ctx context.Context, serviceName string) (int, error)
}

// ScalingExecutor executes scaling operations
type ScalingExecutor interface {
	Scale(ctx context.Context, serviceName string, targetReplicas int) error
	GetServiceInfo(ctx context.Context, serviceName string) (*ServiceInfo, error)
}

// ServiceInfo represents service information
type ServiceInfo struct {
	Name            string            `json:"name"`
	CurrentReplicas int               `json:"current_replicas"`
	ReadyReplicas   int               `json:"ready_replicas"`
	Labels          map[string]string `json:"labels"`
	Annotations     map[string]string `json:"annotations"`
}

// ScalingDecision represents a scaling decision
type ScalingDecision struct {
	ServiceName     string                 `json:"service_name"`
	CurrentReplicas int                    `json:"current_replicas"`
	TargetReplicas  int                    `json:"target_replicas"`
	Reason          string                 `json:"reason"`
	Metrics         map[string]float64     `json:"metrics"`
	Timestamp       time.Time              `json:"timestamp"`
	Action          ScalingAction          `json:"action"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// ScalingAction defines the type of scaling action
type ScalingAction string

const (
	ScaleUp   ScalingAction = "scale_up"
	ScaleDown ScalingAction = "scale_down"
	NoAction  ScalingAction = "no_action"
)

// NewAutoScaler creates a new auto-scaler
func NewAutoScaler(config *AutoScalerConfig, metrics MetricsProvider, executor ScalingExecutor) *AutoScaler {
	if config == nil {
		config = &AutoScalerConfig{
			EvaluationInterval:  30 * time.Second,
			CooldownPeriod:      5 * time.Minute,
			MetricsWindow:       5 * time.Minute,
			MaxScaleUpRate:      1.0,  // 100% increase max
			MaxScaleDownRate:    0.5,  // 50% decrease max
			StabilizationWindow: 2 * time.Minute,
		}
	}

	return &AutoScaler{
		policies: make(map[string]*ScalingPolicy),
		metrics:  metrics,
		executor: executor,
		config:   config,
		stopCh:   make(chan struct{}),
	}
}

// RegisterPolicy registers a scaling policy
func (as *AutoScaler) RegisterPolicy(policy *ScalingPolicy) error {
	as.mu.Lock()
	defer as.mu.Unlock()

	if policy.ID == uuid.Nil {
		policy.ID = uuid.New()
	}

	policy.CreatedAt = time.Now()
	policy.UpdatedAt = time.Now()

	// Validate policy
	if err := as.validatePolicy(policy); err != nil {
		return err
	}

	as.policies[policy.ServiceName] = policy
	return nil
}

// Start starts the auto-scaler
func (as *AutoScaler) Start(ctx context.Context) {
	ticker := time.NewTicker(as.config.EvaluationInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			as.evaluateAll(ctx)
		case <-as.stopCh:
			return
		case <-ctx.Done():
			return
		}
	}
}

// Stop stops the auto-scaler
func (as *AutoScaler) Stop() {
	close(as.stopCh)
}

// evaluateAll evaluates all registered policies
func (as *AutoScaler) evaluateAll(ctx context.Context) {
	as.mu.RLock()
	policies := make([]*ScalingPolicy, 0, len(as.policies))
	for _, policy := range as.policies {
		if policy.Enabled {
			policies = append(policies, policy)
		}
	}
	as.mu.RUnlock()

	var wg sync.WaitGroup
	for _, policy := range policies {
		wg.Add(1)
		go func(p *ScalingPolicy) {
			defer wg.Done()
			as.evaluatePolicy(ctx, p)
		}(policy)
	}
	wg.Wait()
}

// evaluatePolicy evaluates a single scaling policy
func (as *AutoScaler) evaluatePolicy(ctx context.Context, policy *ScalingPolicy) {
	// Check cooldown period
	if time.Since(policy.LastScaledAt) < as.config.CooldownPeriod {
		return
	}

	// Get current metrics
	metrics, err := as.collectMetrics(ctx, policy)
	if err != nil {
		// Log error
		return
	}

	// Calculate desired replicas
	decision := as.calculateScalingDecision(policy, metrics)

	// Execute scaling if needed
	if decision.Action != NoAction {
		err := as.executeScaling(ctx, decision)
		if err != nil {
			// Log error
			return
		}

		// Update policy state
		as.mu.Lock()
		policy.LastScaledAt = time.Now()
		policy.CurrentReplicas = decision.TargetReplicas
		policy.UpdatedAt = time.Now()
		as.mu.Unlock()
	}
}

// collectMetrics collects metrics for a service
func (as *AutoScaler) collectMetrics(ctx context.Context, policy *ScalingPolicy) (map[string]float64, error) {
	metrics := make(map[string]float64)

	// Get current replicas
	currentReplicas, err := as.metrics.GetCurrentReplicas(ctx, policy.ServiceName)
	if err != nil {
		return nil, err
	}
	metrics["current_replicas"] = float64(currentReplicas)

	// Get CPU utilization
	if policy.TargetCPU > 0 {
		cpu, err := as.metrics.GetCPUUtilization(ctx, policy.ServiceName, as.config.MetricsWindow)
		if err != nil {
			return nil, err
		}
		metrics["cpu_utilization"] = cpu
	}

	// Get memory utilization
	if policy.TargetMemory > 0 {
		memory, err := as.metrics.GetMemoryUtilization(ctx, policy.ServiceName, as.config.MetricsWindow)
		if err != nil {
			return nil, err
		}
		metrics["memory_utilization"] = memory
	}

	// Get requests per second
	if policy.TargetRPS > 0 {
		rps, err := as.metrics.GetRequestsPerSecond(ctx, policy.ServiceName, as.config.MetricsWindow)
		if err != nil {
			return nil, err
		}
		metrics["requests_per_second"] = rps
	}

	// Get custom metrics
	for _, customMetric := range policy.CustomMetrics {
		value, err := as.metrics.GetCustomMetric(ctx, policy.ServiceName, customMetric.Name, as.config.MetricsWindow)
		if err != nil {
			continue // Skip failed custom metrics
		}
		metrics[customMetric.Name] = value
	}

	return metrics, nil
}

// calculateScalingDecision calculates the scaling decision based on metrics
func (as *AutoScaler) calculateScalingDecision(policy *ScalingPolicy, metrics map[string]float64) *ScalingDecision {
	currentReplicas := int(metrics["current_replicas"])
	targetReplicas := currentReplicas

	var reasons []string
	var action ScalingAction = NoAction

	// CPU-based scaling
	if policy.TargetCPU > 0 && metrics["cpu_utilization"] > 0 {
		cpuRatio := metrics["cpu_utilization"] / policy.TargetCPU
		cpuTargetReplicas := int(math.Ceil(float64(currentReplicas) * cpuRatio))
		
		if cpuTargetReplicas > targetReplicas {
			targetReplicas = cpuTargetReplicas
			reasons = append(reasons, fmt.Sprintf("CPU utilization %.1f%% > target %.1f%%", 
				metrics["cpu_utilization"], policy.TargetCPU))
		}
	}

	// Memory-based scaling
	if policy.TargetMemory > 0 && metrics["memory_utilization"] > 0 {
		memoryRatio := metrics["memory_utilization"] / policy.TargetMemory
		memoryTargetReplicas := int(math.Ceil(float64(currentReplicas) * memoryRatio))
		
		if memoryTargetReplicas > targetReplicas {
			targetReplicas = memoryTargetReplicas
			reasons = append(reasons, fmt.Sprintf("Memory utilization %.1f%% > target %.1f%%", 
				metrics["memory_utilization"], policy.TargetMemory))
		}
	}

	// RPS-based scaling
	if policy.TargetRPS > 0 && metrics["requests_per_second"] > 0 {
		rpsRatio := metrics["requests_per_second"] / policy.TargetRPS
		rpsTargetReplicas := int(math.Ceil(float64(currentReplicas) * rpsRatio))
		
		if rpsTargetReplicas > targetReplicas {
			targetReplicas = rpsTargetReplicas
			reasons = append(reasons, fmt.Sprintf("RPS %.1f > target %.1f", 
				metrics["requests_per_second"], policy.TargetRPS))
		}
	}

	// Custom metrics scaling
	for _, customMetric := range policy.CustomMetrics {
		if value, exists := metrics[customMetric.Name]; exists && value > 0 {
			var ratio float64
			switch customMetric.Type {
			case "utilization":
				ratio = value / customMetric.Target
			case "average":
				ratio = value / customMetric.Target
			default: // "value"
				ratio = value / customMetric.Target
			}

			customTargetReplicas := int(math.Ceil(float64(currentReplicas) * ratio * customMetric.Weight))
			if customTargetReplicas > targetReplicas {
				targetReplicas = customTargetReplicas
				reasons = append(reasons, fmt.Sprintf("%s %.1f > target %.1f", 
					customMetric.Name, value, customMetric.Target))
			}
		}
	}

	// Apply min/max constraints
	if targetReplicas < policy.MinReplicas {
		targetReplicas = policy.MinReplicas
		if targetReplicas > currentReplicas {
			reasons = append(reasons, fmt.Sprintf("Enforcing minimum replicas %d", policy.MinReplicas))
		}
	}

	if targetReplicas > policy.MaxReplicas {
		targetReplicas = policy.MaxReplicas
		reasons = append(reasons, fmt.Sprintf("Enforcing maximum replicas %d", policy.MaxReplicas))
	}

	// Apply rate limiting
	targetReplicas = as.applyRateLimits(currentReplicas, targetReplicas)

	// Determine action
	if targetReplicas > currentReplicas {
		action = ScaleUp
	} else if targetReplicas < currentReplicas {
		action = ScaleDown
	}

	reasonText := "No scaling needed"
	if len(reasons) > 0 {
		reasonText = fmt.Sprintf("Scaling due to: %v", reasons)
	}

	return &ScalingDecision{
		ServiceName:     policy.ServiceName,
		CurrentReplicas: currentReplicas,
		TargetReplicas:  targetReplicas,
		Reason:          reasonText,
		Metrics:         metrics,
		Timestamp:       time.Now(),
		Action:          action,
		Metadata: map[string]interface{}{
			"policy_id": policy.ID.String(),
		},
	}
}

// applyRateLimits applies rate limiting to scaling decisions
func (as *AutoScaler) applyRateLimits(currentReplicas, targetReplicas int) int {
	if targetReplicas > currentReplicas {
		// Scale up rate limiting
		maxIncrease := int(math.Ceil(float64(currentReplicas) * as.config.MaxScaleUpRate))
		if targetReplicas-currentReplicas > maxIncrease {
			return currentReplicas + maxIncrease
		}
	} else if targetReplicas < currentReplicas {
		// Scale down rate limiting
		maxDecrease := int(math.Ceil(float64(currentReplicas) * as.config.MaxScaleDownRate))
		if currentReplicas-targetReplicas > maxDecrease {
			return currentReplicas - maxDecrease
		}
	}

	return targetReplicas
}

// executeScaling executes the scaling decision
func (as *AutoScaler) executeScaling(ctx context.Context, decision *ScalingDecision) error {
	if decision.Action == NoAction {
		return nil
	}

	return as.executor.Scale(ctx, decision.ServiceName, decision.TargetReplicas)
}

// validatePolicy validates a scaling policy
func (as *AutoScaler) validatePolicy(policy *ScalingPolicy) error {
	if policy.ServiceName == "" {
		return fmt.Errorf("service name is required")
	}

	if policy.MinReplicas < 1 {
		return fmt.Errorf("minimum replicas must be at least 1")
	}

	if policy.MaxReplicas < policy.MinReplicas {
		return fmt.Errorf("maximum replicas must be greater than or equal to minimum replicas")
	}

	if policy.TargetCPU < 0 || policy.TargetCPU > 100 {
		return fmt.Errorf("target CPU must be between 0 and 100")
	}

	if policy.TargetMemory < 0 || policy.TargetMemory > 100 {
		return fmt.Errorf("target memory must be between 0 and 100")
	}

	if policy.TargetRPS < 0 {
		return fmt.Errorf("target RPS must be non-negative")
	}

	return nil
}

// GetPolicyStatus gets the current status of a scaling policy
func (as *AutoScaler) GetPolicyStatus(serviceName string) (*PolicyStatus, error) {
	as.mu.RLock()
	policy, exists := as.policies[serviceName]
	as.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("policy for service %s not found", serviceName)
	}

	return &PolicyStatus{
		PolicyID:        policy.ID,
		ServiceName:     policy.ServiceName,
		Enabled:         policy.Enabled,
		CurrentReplicas: policy.CurrentReplicas,
		MinReplicas:     policy.MinReplicas,
		MaxReplicas:     policy.MaxReplicas,
		LastScaledAt:    policy.LastScaledAt,
		UpdatedAt:       policy.UpdatedAt,
	}, nil
}

// PolicyStatus represents the current status of a scaling policy
type PolicyStatus struct {
	PolicyID        uuid.UUID `json:"policy_id"`
	ServiceName     string    `json:"service_name"`
	Enabled         bool      `json:"enabled"`
	CurrentReplicas int       `json:"current_replicas"`
	MinReplicas     int       `json:"min_replicas"`
	MaxReplicas     int       `json:"max_replicas"`
	LastScaledAt    time.Time `json:"last_scaled_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Kubernetes Auto-Scaler Implementation

// KubernetesMetricsProvider implements MetricsProvider for Kubernetes
type KubernetesMetricsProvider struct {
	config *K8sConfig
}

// K8sConfig holds Kubernetes configuration
type K8sConfig struct {
	Namespace     string
	MetricsServer string
	PrometheusURL string
}

// NewKubernetesMetricsProvider creates a new Kubernetes metrics provider
func NewKubernetesMetricsProvider(config *K8sConfig) *KubernetesMetricsProvider {
	return &KubernetesMetricsProvider{config: config}
}

// GetCPUUtilization gets CPU utilization from Kubernetes metrics
func (k *KubernetesMetricsProvider) GetCPUUtilization(ctx context.Context, serviceName string, window time.Duration) (float64, error) {
	// Implementation would query Kubernetes metrics API or Prometheus
	// This is a placeholder implementation
	return 50.0, nil
}

// GetMemoryUtilization gets memory utilization from Kubernetes metrics
func (k *KubernetesMetricsProvider) GetMemoryUtilization(ctx context.Context, serviceName string, window time.Duration) (float64, error) {
	// Implementation would query Kubernetes metrics API or Prometheus
	return 60.0, nil
}

// GetRequestsPerSecond gets RPS from monitoring system
func (k *KubernetesMetricsProvider) GetRequestsPerSecond(ctx context.Context, serviceName string, window time.Duration) (float64, error) {
	// Implementation would query application metrics
	return 100.0, nil
}

// GetCustomMetric gets custom metric value
func (k *KubernetesMetricsProvider) GetCustomMetric(ctx context.Context, serviceName, metricName string, window time.Duration) (float64, error) {
	// Implementation would query custom metrics from Prometheus or other sources
	return 75.0, nil
}

// GetCurrentReplicas gets current replica count from Kubernetes
func (k *KubernetesMetricsProvider) GetCurrentReplicas(ctx context.Context, serviceName string) (int, error) {
	// Implementation would query Kubernetes API
	return 3, nil
}

// KubernetesScalingExecutor implements ScalingExecutor for Kubernetes
type KubernetesScalingExecutor struct {
	config *K8sConfig
}

// NewKubernetesScalingExecutor creates a new Kubernetes scaling executor
func NewKubernetesScalingExecutor(config *K8sConfig) *KubernetesScalingExecutor {
	return &KubernetesScalingExecutor{config: config}
}

// Scale scales a Kubernetes deployment
func (k *KubernetesScalingExecutor) Scale(ctx context.Context, serviceName string, targetReplicas int) error {
	// Implementation would use Kubernetes API to scale deployment
	// kubectl scale deployment/serviceName --replicas=targetReplicas -n namespace
	return nil
}

// GetServiceInfo gets service information from Kubernetes
func (k *KubernetesScalingExecutor) GetServiceInfo(ctx context.Context, serviceName string) (*ServiceInfo, error) {
	// Implementation would query Kubernetes API for deployment info
	return &ServiceInfo{
		Name:            serviceName,
		CurrentReplicas: 3,
		ReadyReplicas:   3,
		Labels:          map[string]string{"app": serviceName},
		Annotations:     map[string]string{},
	}, nil
}

// Predefined Scaling Policies

// FastenMindAPIScalingPolicy returns a scaling policy for FastenMind API
func FastenMindAPIScalingPolicy() *ScalingPolicy {
	return &ScalingPolicy{
		Name:         "fastenmind_api_scaling",
		ServiceName:  "fastenmind-api",
		MinReplicas:  2,
		MaxReplicas:  20,
		TargetCPU:    70.0,
		TargetMemory: 80.0,
		TargetRPS:    100.0,
		CustomMetrics: []CustomMetricTarget{
			{
				Name:   "inquiry_queue_length",
				Target: 50.0,
				Type:   "value",
				Weight: 1.0,
			},
		},
		Behaviors: &ScalingBehavior{
			ScaleUp: &ScaleUpBehavior{
				StabilizationWindow: 60 * time.Second,
				Policies: []ScalingRulePolicy{
					{
						Type:          "Percent",
						Value:         100,
						PeriodSeconds: 60 * time.Second,
					},
				},
			},
			ScaleDown: &ScaleDownBehavior{
				StabilizationWindow: 300 * time.Second,
				Policies: []ScalingRulePolicy{
					{
						Type:          "Percent",
						Value:         50,
						PeriodSeconds: 60 * time.Second,
					},
				},
			},
		},
		Enabled: true,
	}
}

// FastenMindWorkerScalingPolicy returns a scaling policy for background workers
func FastenMindWorkerScalingPolicy() *ScalingPolicy {
	return &ScalingPolicy{
		Name:        "fastenmind_worker_scaling",
		ServiceName: "fastenmind-worker",
		MinReplicas: 1,
		MaxReplicas: 10,
		TargetCPU:   60.0,
		CustomMetrics: []CustomMetricTarget{
			{
				Name:   "job_queue_length",
				Target: 20.0,
				Type:   "value",
				Weight: 1.5,
			},
			{
				Name:   "processing_time_p95",
				Target: 30000.0, // 30 seconds in milliseconds
				Type:   "average",
				Weight: 1.0,
			},
		},
		Behaviors: &ScalingBehavior{
			ScaleUp: &ScaleUpBehavior{
				StabilizationWindow: 30 * time.Second,
				Policies: []ScalingRulePolicy{
					{
						Type:          "Pods",
						Value:         2,
						PeriodSeconds: 30 * time.Second,
					},
				},
			},
			ScaleDown: &ScaleDownBehavior{
				StabilizationWindow: 600 * time.Second, // 10 minutes
				Policies: []ScalingRulePolicy{
					{
						Type:          "Pods",
						Value:         1,
						PeriodSeconds: 120 * time.Second,
					},
				},
			},
		},
		Enabled: true,
	}
}