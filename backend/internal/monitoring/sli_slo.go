package monitoring

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/google/uuid"
)

// SLIManager manages Service Level Indicators
type SLIManager struct {
	indicators map[string]*SLI
	collectors map[string]SLICollector
	storage    SLIStorage
	mu         sync.RWMutex
	config     *SLIConfig
}

// SLIConfig holds SLI configuration
type SLIConfig struct {
	CollectionInterval time.Duration
	RetentionPeriod    time.Duration
	AlertEnabled       bool
	AlertWebhook       string
	MetricsExport      bool
}

// SLI represents a Service Level Indicator
type SLI struct {
	ID          uuid.UUID                `json:"id"`
	Name        string                   `json:"name"`
	Type        SLIType                  `json:"type"`
	Description string                   `json:"description"`
	Query       string                   `json:"query"`
	Threshold   float64                  `json:"threshold"`
	Unit        string                   `json:"unit"`
	Tags        map[string]string        `json:"tags"`
	Config      map[string]interface{}   `json:"config"`
	CreatedAt   time.Time                `json:"created_at"`
	UpdatedAt   time.Time                `json:"updated_at"`
}

// SLIType defines the type of SLI
type SLIType string

const (
	SLITypeAvailability SLIType = "availability"
	SLITypeLatency      SLIType = "latency"
	SLITypeThroughput   SLIType = "throughput"
	SLITypeErrorRate    SLIType = "error_rate"
	SLITypeSaturation   SLIType = "saturation"
	SLITypeCustom       SLIType = "custom"
)

// SLICollector collects SLI metrics
type SLICollector interface {
	Collect(ctx context.Context, sli *SLI) (float64, error)
	Validate(sli *SLI) error
}

// SLIStorage stores SLI data
type SLIStorage interface {
	Store(ctx context.Context, data *SLIDataPoint) error
	Query(ctx context.Context, filter SLIQueryFilter) ([]*SLIDataPoint, error)
	Aggregate(ctx context.Context, filter SLIQueryFilter, aggregation AggregationType) (*SLIAggregation, error)
}

// SLIDataPoint represents a single SLI measurement
type SLIDataPoint struct {
	SLIID     uuid.UUID         `json:"sli_id"`
	Value     float64           `json:"value"`
	Timestamp time.Time         `json:"timestamp"`
	Tags      map[string]string `json:"tags"`
	Metadata  map[string]interface{} `json:"metadata"`
}

// SLIQueryFilter filters SLI data queries
type SLIQueryFilter struct {
	SLIIDs    []uuid.UUID
	StartTime time.Time
	EndTime   time.Time
	Tags      map[string]string
	Limit     int
}

// AggregationType defines aggregation types
type AggregationType string

const (
	AggregationAverage   AggregationType = "average"
	AggregationSum       AggregationType = "sum"
	AggregationMin       AggregationType = "min"
	AggregationMax       AggregationType = "max"
	AggregationP50       AggregationType = "p50"
	AggregationP90       AggregationType = "p90"
	AggregationP95       AggregationType = "p95"
	AggregationP99       AggregationType = "p99"
	AggregationCount     AggregationType = "count"
)

// SLIAggregation represents aggregated SLI data
type SLIAggregation struct {
	Value     float64                    `json:"value"`
	Count     int64                      `json:"count"`
	StartTime time.Time                  `json:"start_time"`
	EndTime   time.Time                  `json:"end_time"`
	Breakdown map[string]*SLIAggregation `json:"breakdown,omitempty"`
}

// NewSLIManager creates a new SLI manager
func NewSLIManager(config *SLIConfig, storage SLIStorage) *SLIManager {
	return &SLIManager{
		indicators: make(map[string]*SLI),
		collectors: make(map[string]SLICollector),
		storage:    storage,
		config:     config,
	}
}

// RegisterSLI registers a new SLI
func (sm *SLIManager) RegisterSLI(sli *SLI) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sli.ID == uuid.Nil {
		sli.ID = uuid.New()
	}
	
	sli.CreatedAt = time.Now()
	sli.UpdatedAt = time.Now()

	// Validate SLI with collector
	if collector, exists := sm.collectors[string(sli.Type)]; exists {
		if err := collector.Validate(sli); err != nil {
			return fmt.Errorf("SLI validation failed: %w", err)
		}
	}

	sm.indicators[sli.Name] = sli
	return nil
}

// RegisterCollector registers an SLI collector
func (sm *SLIManager) RegisterCollector(sliType SLIType, collector SLICollector) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.collectors[string(sliType)] = collector
}

// StartCollection starts SLI data collection
func (sm *SLIManager) StartCollection(ctx context.Context) {
	ticker := time.NewTicker(sm.config.CollectionInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			sm.collectAll(ctx)
		case <-ctx.Done():
			return
		}
	}
}

// collectAll collects all registered SLIs
func (sm *SLIManager) collectAll(ctx context.Context) {
	sm.mu.RLock()
	indicators := make([]*SLI, 0, len(sm.indicators))
	for _, sli := range sm.indicators {
		indicators = append(indicators, sli)
	}
	sm.mu.RUnlock()

	var wg sync.WaitGroup
	for _, sli := range indicators {
		wg.Add(1)
		go func(s *SLI) {
			defer wg.Done()
			sm.collectSLI(ctx, s)
		}(sli)
	}
	wg.Wait()
}

// collectSLI collects a single SLI
func (sm *SLIManager) collectSLI(ctx context.Context, sli *SLI) {
	sm.mu.RLock()
	collector, exists := sm.collectors[string(sli.Type)]
	sm.mu.RUnlock()

	if !exists {
		// Log error: no collector for SLI type
		return
	}

	value, err := collector.Collect(ctx, sli)
	if err != nil {
		// Log error
		return
	}

	dataPoint := &SLIDataPoint{
		SLIID:     sli.ID,
		Value:     value,
		Timestamp: time.Now(),
		Tags:      sli.Tags,
		Metadata:  map[string]interface{}{
			"sli_name": sli.Name,
			"sli_type": sli.Type,
		},
	}

	if err := sm.storage.Store(ctx, dataPoint); err != nil {
		// Log error
		return
	}
}

// GetSLIData gets SLI data for a time range
func (sm *SLIManager) GetSLIData(ctx context.Context, sliName string, startTime, endTime time.Time) ([]*SLIDataPoint, error) {
	sm.mu.RLock()
	sli, exists := sm.indicators[sliName]
	sm.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("SLI %s not found", sliName)
	}

	filter := SLIQueryFilter{
		SLIIDs:    []uuid.UUID{sli.ID},
		StartTime: startTime,
		EndTime:   endTime,
	}

	return sm.storage.Query(ctx, filter)
}

// GetSLIAggregation gets aggregated SLI data
func (sm *SLIManager) GetSLIAggregation(ctx context.Context, sliName string, startTime, endTime time.Time, aggregation AggregationType) (*SLIAggregation, error) {
	sm.mu.RLock()
	sli, exists := sm.indicators[sliName]
	sm.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("SLI %s not found", sliName)
	}

	filter := SLIQueryFilter{
		SLIIDs:    []uuid.UUID{sli.ID},
		StartTime: startTime,
		EndTime:   endTime,
	}

	return sm.storage.Aggregate(ctx, filter, aggregation)
}

// SLOManager manages Service Level Objectives
type SLOManager struct {
	objectives map[string]*SLO
	sliManager *SLIManager
	alerter    SLOAlerter
	mu         sync.RWMutex
	config     *SLOConfig
}

// SLOConfig holds SLO configuration
type SLOConfig struct {
	EvaluationInterval time.Duration
	AlertCooldown      time.Duration
	BurnRateThresholds map[string]float64 // window -> threshold
}

// SLO represents a Service Level Objective
type SLO struct {
	ID            uuid.UUID       `json:"id"`
	Name          string          `json:"name"`
	Description   string          `json:"description"`
	SLIName       string          `json:"sli_name"`
	Target        float64         `json:"target"`        // e.g., 99.9 for 99.9%
	TimeWindow    time.Duration   `json:"time_window"`   // e.g., 30 days
	AlertRules    []SLOAlertRule  `json:"alert_rules"`
	Tags          map[string]string `json:"tags"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

// SLOAlertRule defines when to alert for SLO violations
type SLOAlertRule struct {
	Name        string        `json:"name"`
	Severity    AlertSeverity `json:"severity"`
	BurnRate    float64       `json:"burn_rate"`     // Multiple of error budget burn rate
	TimeWindow  time.Duration `json:"time_window"`   // Time window to evaluate
	Threshold   float64       `json:"threshold"`     // Threshold for alerting
}

// AlertSeverity defines alert severity levels
type AlertSeverity string

const (
	AlertSeverityInfo     AlertSeverity = "info"
	AlertSeverityWarning  AlertSeverity = "warning"
	AlertSeverityCritical AlertSeverity = "critical"
)

// SLOAlerter sends SLO alerts
type SLOAlerter interface {
	SendAlert(ctx context.Context, alert *SLOAlert) error
}

// SLOAlert represents an SLO alert
type SLOAlert struct {
	ID          uuid.UUID                `json:"id"`
	SLOName     string                   `json:"slo_name"`
	RuleName    string                   `json:"rule_name"`
	Severity    AlertSeverity            `json:"severity"`
	Message     string                   `json:"message"`
	Value       float64                  `json:"value"`
	Threshold   float64                  `json:"threshold"`
	Timestamp   time.Time                `json:"timestamp"`
	Tags        map[string]string        `json:"tags"`
	Metadata    map[string]interface{}   `json:"metadata"`
}

// NewSLOManager creates a new SLO manager
func NewSLOManager(config *SLOConfig, sliManager *SLIManager, alerter SLOAlerter) *SLOManager {
	return &SLOManager{
		objectives: make(map[string]*SLO),
		sliManager: sliManager,
		alerter:    alerter,
		config:     config,
	}
}

// RegisterSLO registers a new SLO
func (som *SLOManager) RegisterSLO(slo *SLO) error {
	som.mu.Lock()
	defer som.mu.Unlock()

	if slo.ID == uuid.Nil {
		slo.ID = uuid.New()
	}

	slo.CreatedAt = time.Now()
	slo.UpdatedAt = time.Now()

	// Validate that SLI exists
	som.sliManager.mu.RLock()
	_, exists := som.sliManager.indicators[slo.SLIName]
	som.sliManager.mu.RUnlock()

	if !exists {
		return fmt.Errorf("SLI %s not found", slo.SLIName)
	}

	som.objectives[slo.Name] = slo
	return nil
}

// StartEvaluation starts SLO evaluation
func (som *SLOManager) StartEvaluation(ctx context.Context) {
	ticker := time.NewTicker(som.config.EvaluationInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			som.evaluateAll(ctx)
		case <-ctx.Done():
			return
		}
	}
}

// evaluateAll evaluates all SLOs
func (som *SLOManager) evaluateAll(ctx context.Context) {
	som.mu.RLock()
	objectives := make([]*SLO, 0, len(som.objectives))
	for _, slo := range som.objectives {
		objectives = append(objectives, slo)
	}
	som.mu.RUnlock()

	for _, slo := range objectives {
		som.evaluateSLO(ctx, slo)
	}
}

// evaluateSLO evaluates a single SLO
func (som *SLOManager) evaluateSLO(ctx context.Context, slo *SLO) {
	now := time.Now()
	startTime := now.Add(-slo.TimeWindow)

	// Get current SLI value
	aggregation, err := som.sliManager.GetSLIAggregation(ctx, slo.SLIName, startTime, now, AggregationAverage)
	if err != nil {
		// Log error
		return
	}

	currentValue := aggregation.Value
	errorBudget := 100.0 - slo.Target
	burnRate := (100.0 - currentValue) / errorBudget

	// Evaluate alert rules
	for _, rule := range slo.AlertRules {
		ruleStartTime := now.Add(-rule.TimeWindow)
		
		ruleAggregation, err := som.sliManager.GetSLIAggregation(ctx, slo.SLIName, ruleStartTime, now, AggregationAverage)
		if err != nil {
			continue
		}

		ruleValue := ruleAggregation.Value
		ruleBurnRate := (100.0 - ruleValue) / errorBudget

		if ruleBurnRate >= rule.BurnRate {
			alert := &SLOAlert{
				ID:        uuid.New(),
				SLOName:   slo.Name,
				RuleName:  rule.Name,
				Severity:  rule.Severity,
				Message:   fmt.Sprintf("SLO %s burn rate %.2f exceeds threshold %.2f", slo.Name, ruleBurnRate, rule.BurnRate),
				Value:     ruleBurnRate,
				Threshold: rule.BurnRate,
				Timestamp: now,
				Tags:      slo.Tags,
				Metadata: map[string]interface{}{
					"sli_value":     ruleValue,
					"error_budget":  errorBudget,
					"time_window":   rule.TimeWindow.String(),
				},
			}

			if err := som.alerter.SendAlert(ctx, alert); err != nil {
				// Log error
			}
		}
	}
}

// GetSLOStatus gets current SLO status
func (som *SLOManager) GetSLOStatus(ctx context.Context, sloName string) (*SLOStatus, error) {
	som.mu.RLock()
	slo, exists := som.objectives[sloName]
	som.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("SLO %s not found", sloName)
	}

	now := time.Now()
	startTime := now.Add(-slo.TimeWindow)

	aggregation, err := som.sliManager.GetSLIAggregation(ctx, slo.SLIName, startTime, now, AggregationAverage)
	if err != nil {
		return nil, err
	}

	currentValue := aggregation.Value
	errorBudget := 100.0 - slo.Target
	usedErrorBudget := 100.0 - currentValue
	remainingErrorBudget := errorBudget - usedErrorBudget
	burnRate := usedErrorBudget / errorBudget

	status := &SLOStatus{
		SLOName:               slo.Name,
		Target:                slo.Target,
		CurrentValue:          currentValue,
		ErrorBudget:           errorBudget,
		UsedErrorBudget:       usedErrorBudget,
		RemainingErrorBudget:  remainingErrorBudget,
		BurnRate:              burnRate,
		TimeToExhaustion:      som.calculateTimeToExhaustion(burnRate, remainingErrorBudget, slo.TimeWindow),
		IsHealthy:             currentValue >= slo.Target,
		LastEvaluated:         now,
	}

	return status, nil
}

// SLOStatus represents current SLO status
type SLOStatus struct {
	SLOName               string        `json:"slo_name"`
	Target                float64       `json:"target"`
	CurrentValue          float64       `json:"current_value"`
	ErrorBudget           float64       `json:"error_budget"`
	UsedErrorBudget       float64       `json:"used_error_budget"`
	RemainingErrorBudget  float64       `json:"remaining_error_budget"`
	BurnRate              float64       `json:"burn_rate"`
	TimeToExhaustion      time.Duration `json:"time_to_exhaustion"`
	IsHealthy             bool          `json:"is_healthy"`
	LastEvaluated         time.Time     `json:"last_evaluated"`
}

// calculateTimeToExhaustion calculates when error budget will be exhausted
func (som *SLOManager) calculateTimeToExhaustion(burnRate, remainingBudget float64, timeWindow time.Duration) time.Duration {
	if burnRate <= 0 || remainingBudget <= 0 {
		return time.Duration(math.MaxInt64) // Never
	}

	hoursToExhaustion := remainingBudget / (burnRate * 100.0) * timeWindow.Hours()
	return time.Duration(hoursToExhaustion * float64(time.Hour))
}

// Common SLI Collectors

// AvailabilitySLICollector collects availability metrics
type AvailabilitySLICollector struct {
	httpClient HTTPClient
}

type HTTPClient interface {
	Get(ctx context.Context, url string) (*HTTPResponse, error)
}

type HTTPResponse struct {
	StatusCode int
	Duration   time.Duration
	Error      error
}

func NewAvailabilitySLICollector(client HTTPClient) *AvailabilitySLICollector {
	return &AvailabilitySLICollector{httpClient: client}
}

func (c *AvailabilitySLICollector) Collect(ctx context.Context, sli *SLI) (float64, error) {
	endpoint, ok := sli.Config["endpoint"].(string)
	if !ok {
		return 0, fmt.Errorf("endpoint not configured for availability SLI")
	}

	resp, err := c.httpClient.Get(ctx, endpoint)
	if err != nil {
		return 0, err // 0% availability
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return 100, nil // 100% availability
	}

	return 0, nil // 0% availability
}

func (c *AvailabilitySLICollector) Validate(sli *SLI) error {
	if _, ok := sli.Config["endpoint"]; !ok {
		return fmt.Errorf("endpoint configuration required for availability SLI")
	}
	return nil
}

// LatencySLICollector collects latency metrics
type LatencySLICollector struct {
	httpClient HTTPClient
}

func NewLatencySLICollector(client HTTPClient) *LatencySLICollector {
	return &LatencySLICollector{httpClient: client}
}

func (c *LatencySLICollector) Collect(ctx context.Context, sli *SLI) (float64, error) {
	endpoint, ok := sli.Config["endpoint"].(string)
	if !ok {
		return 0, fmt.Errorf("endpoint not configured for latency SLI")
	}

	resp, err := c.httpClient.Get(ctx, endpoint)
	if err != nil {
		return 0, err
	}

	return resp.Duration.Seconds() * 1000, nil // Return latency in milliseconds
}

func (c *LatencySLICollector) Validate(sli *SLI) error {
	if _, ok := sli.Config["endpoint"]; !ok {
		return fmt.Errorf("endpoint configuration required for latency SLI")
	}
	return nil
}

// Predefined SLIs and SLOs

func FastenMindAvailabilitySLI() *SLI {
	return &SLI{
		Name:        "fastenmind_availability",
		Type:        SLITypeAvailability,
		Description: "FastenMind system availability",
		Threshold:   99.9,
		Unit:        "percent",
		Tags: map[string]string{
			"service":     "fastenmind",
			"environment": "production",
		},
		Config: map[string]interface{}{
			"endpoint": "https://api.fastenmind.com/health",
		},
	}
}

func FastenMindLatencySLI() *SLI {
	return &SLI{
		Name:        "fastenmind_api_latency",
		Type:        SLITypeLatency,
		Description: "FastenMind API response latency",
		Threshold:   200, // 200ms
		Unit:        "milliseconds",
		Tags: map[string]string{
			"service":     "fastenmind",
			"environment": "production",
		},
		Config: map[string]interface{}{
			"endpoint": "https://api.fastenmind.com/api/v1/inquiries",
		},
	}
}

func FastenMindAvailabilitySLO() *SLO {
	return &SLO{
		Name:        "fastenmind_availability_slo",
		Description: "FastenMind system should be available 99.9% of the time",
		SLIName:     "fastenmind_availability",
		Target:      99.9,
		TimeWindow:  30 * 24 * time.Hour, // 30 days
		AlertRules: []SLOAlertRule{
			{
				Name:       "fast_burn",
				Severity:   AlertSeverityCritical,
				BurnRate:   14.4, // 5 minutes of 100% failure in 1 hour
				TimeWindow: 1 * time.Hour,
				Threshold:  0.02,
			},
			{
				Name:       "slow_burn",
				Severity:   AlertSeverityWarning,
				BurnRate:   1.0, // Normal burn rate
				TimeWindow: 24 * time.Hour,
				Threshold:  0.1,
			},
		},
		Tags: map[string]string{
			"service":     "fastenmind",
			"environment": "production",
		},
	}
}

func FastenMindLatencySLO() *SLO {
	return &SLO{
		Name:        "fastenmind_latency_slo",
		Description: "FastenMind API should respond within 200ms for 95% of requests",
		SLIName:     "fastenmind_api_latency",
		Target:      95.0, // 95% of requests under 200ms
		TimeWindow:  7 * 24 * time.Hour, // 7 days
		AlertRules: []SLOAlertRule{
			{
				Name:       "latency_degradation",
				Severity:   AlertSeverityWarning,
				BurnRate:   2.0,
				TimeWindow: 2 * time.Hour,
				Threshold:  0.05,
			},
		},
		Tags: map[string]string{
			"service":     "fastenmind",
			"environment": "production",
		},
	}
}