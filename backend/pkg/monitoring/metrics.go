package monitoring

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/labstack/echo/v4"
)

// Metrics 監控指標
type Metrics struct {
	// HTTP 請求指標
	HTTPRequestsTotal   *prometheus.CounterVec
	HTTPRequestDuration *prometheus.HistogramVec
	HTTPResponseSize    *prometheus.HistogramVec

	// 業務指標
	BusinessOperations  *prometheus.CounterVec
	BusinessErrors      *prometheus.CounterVec
	ActiveUsers         prometheus.Gauge
	QuotationValue      *prometheus.HistogramVec

	// 資料庫指標
	DBConnectionsActive prometheus.Gauge
	DBConnectionsIdle   prometheus.Gauge
	DBQueryDuration     *prometheus.HistogramVec
	DBErrors            *prometheus.CounterVec

	// 快取指標
	CacheHits         *prometheus.CounterVec
	CacheMisses       *prometheus.CounterVec
	CacheSize         prometheus.Gauge
	CacheEvictions    prometheus.Counter

	// 系統指標
	SystemMemoryUsage  prometheus.Gauge
	SystemCPUUsage     prometheus.Gauge
	SystemGoroutines   prometheus.Gauge
	SystemOpenFiles    prometheus.Gauge

	// 自定義指標
	CustomMetrics map[string]prometheus.Collector
}

// NewMetrics 創建新的監控指標
func NewMetrics() *Metrics {
	return &Metrics{
		// HTTP 指標
		HTTPRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "path", "status"},
		),
		HTTPRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path"},
		),
		HTTPResponseSize: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_response_size_bytes",
				Help:    "HTTP response size in bytes",
				Buckets: prometheus.ExponentialBuckets(100, 10, 8),
			},
			[]string{"method", "path"},
		),

		// 業務指標
		BusinessOperations: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "business_operations_total",
				Help: "Total number of business operations",
			},
			[]string{"operation", "status"},
		),
		BusinessErrors: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "business_errors_total",
				Help: "Total number of business errors",
			},
			[]string{"operation", "error_type"},
		),
		ActiveUsers: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "active_users",
				Help: "Number of active users",
			},
		),
		QuotationValue: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "quotation_value_usd",
				Help:    "Quotation value in USD",
				Buckets: prometheus.ExponentialBuckets(100, 10, 10),
			},
			[]string{"customer", "status"},
		),

		// 資料庫指標
		DBConnectionsActive: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "db_connections_active",
				Help: "Number of active database connections",
			},
		),
		DBConnectionsIdle: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "db_connections_idle",
				Help: "Number of idle database connections",
			},
		),
		DBQueryDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "db_query_duration_seconds",
				Help:    "Database query duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"query_type", "table"},
		),
		DBErrors: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "db_errors_total",
				Help: "Total number of database errors",
			},
			[]string{"operation", "error_type"},
		),

		// 快取指標
		CacheHits: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "cache_hits_total",
				Help: "Total number of cache hits",
			},
			[]string{"cache_name"},
		),
		CacheMisses: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "cache_misses_total",
				Help: "Total number of cache misses",
			},
			[]string{"cache_name"},
		),
		CacheSize: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "cache_size_items",
				Help: "Number of items in cache",
			},
		),
		CacheEvictions: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "cache_evictions_total",
				Help: "Total number of cache evictions",
			},
		),

		// 系統指標
		SystemMemoryUsage: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "system_memory_usage_bytes",
				Help: "System memory usage in bytes",
			},
		),
		SystemCPUUsage: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "system_cpu_usage_percent",
				Help: "System CPU usage percentage",
			},
		),
		SystemGoroutines: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "system_goroutines",
				Help: "Number of goroutines",
			},
		),
		SystemOpenFiles: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "system_open_files",
				Help: "Number of open files",
			},
		),

		CustomMetrics: make(map[string]prometheus.Collector),
	}
}

// RecordHTTPRequest 記錄 HTTP 請求
func (m *Metrics) RecordHTTPRequest(method, path string, status int, duration time.Duration, size int) {
	statusStr := fmt.Sprintf("%d", status)
	m.HTTPRequestsTotal.WithLabelValues(method, path, statusStr).Inc()
	m.HTTPRequestDuration.WithLabelValues(method, path).Observe(duration.Seconds())
	m.HTTPResponseSize.WithLabelValues(method, path).Observe(float64(size))
}

// RecordBusinessOperation 記錄業務操作
func (m *Metrics) RecordBusinessOperation(operation string, success bool) {
	status := "success"
	if !success {
		status = "failure"
	}
	m.BusinessOperations.WithLabelValues(operation, status).Inc()
}

// RecordBusinessError 記錄業務錯誤
func (m *Metrics) RecordBusinessError(operation, errorType string) {
	m.BusinessErrors.WithLabelValues(operation, errorType).Inc()
}

// RecordQuotation 記錄報價
func (m *Metrics) RecordQuotation(customer string, status string, value float64) {
	m.QuotationValue.WithLabelValues(customer, status).Observe(value)
}

// RecordDBQuery 記錄資料庫查詢
func (m *Metrics) RecordDBQuery(queryType, table string, duration time.Duration) {
	m.DBQueryDuration.WithLabelValues(queryType, table).Observe(duration.Seconds())
}

// RecordDBError 記錄資料庫錯誤
func (m *Metrics) RecordDBError(operation, errorType string) {
	m.DBErrors.WithLabelValues(operation, errorType).Inc()
}

// RecordCacheHit 記錄快取命中
func (m *Metrics) RecordCacheHit(cacheName string) {
	m.CacheHits.WithLabelValues(cacheName).Inc()
}

// RecordCacheMiss 記錄快取未命中
func (m *Metrics) RecordCacheMiss(cacheName string) {
	m.CacheMisses.WithLabelValues(cacheName).Inc()
}

// UpdateSystemMetrics 更新系統指標
func (m *Metrics) UpdateSystemMetrics(memoryUsage, cpuUsage float64, goroutines, openFiles int) {
	m.SystemMemoryUsage.Set(memoryUsage)
	m.SystemCPUUsage.Set(cpuUsage)
	m.SystemGoroutines.Set(float64(goroutines))
	m.SystemOpenFiles.Set(float64(openFiles))
}

// RegisterCustomMetric 註冊自定義指標
func (m *Metrics) RegisterCustomMetric(name string, collector prometheus.Collector) error {
	if err := prometheus.Register(collector); err != nil {
		return err
	}
	m.CustomMetrics[name] = collector
	return nil
}

// PrometheusHandler 返回 Prometheus HTTP 處理器
func PrometheusHandler() echo.HandlerFunc {
	return echo.WrapHandler(promhttp.Handler())
}

// MetricsMiddleware Echo 中間件
func MetricsMiddleware(metrics *Metrics) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			
			// 處理請求
			err := next(c)
			
			// 記錄指標
			duration := time.Since(start)
			status := c.Response().Status
			size := c.Response().Size
			
			// 清理路徑參數
			path := c.Path()
			if path == "" {
				path = c.Request().URL.Path
			}
			
			metrics.RecordHTTPRequest(
				c.Request().Method,
				path,
				status,
				duration,
				int(size),
			)
			
			return err
		}
	}
}

// HealthChecker 健康檢查器
type HealthChecker struct {
	checks map[string]HealthCheck
	mu     sync.RWMutex
}

// HealthCheck 健康檢查函數
type HealthCheck func(ctx context.Context) error

// NewHealthChecker 創建健康檢查器
func NewHealthChecker() *HealthChecker {
	return &HealthChecker{
		checks: make(map[string]HealthCheck),
	}
}

// RegisterCheck 註冊健康檢查
func (h *HealthChecker) RegisterCheck(name string, check HealthCheck) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.checks[name] = check
}

// CheckHealth 執行健康檢查
func (h *HealthChecker) CheckHealth(ctx context.Context) map[string]HealthStatus {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	results := make(map[string]HealthStatus)
	
	for name, check := range h.checks {
		status := HealthStatus{
			Status: "healthy",
		}
		
		checkCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		err := check(checkCtx)
		cancel()
		
		if err != nil {
			status.Status = "unhealthy"
			status.Error = err.Error()
		}
		
		results[name] = status
	}
	
	return results
}

// HealthStatus 健康狀態
type HealthStatus struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

// AlertManager 告警管理器
type AlertManager struct {
	rules     []AlertRule
	notifiers []Notifier
	mu        sync.RWMutex
}

// AlertRule 告警規則
type AlertRule struct {
	Name        string
	Expression  string
	Duration    time.Duration
	Severity    string
	Labels      map[string]string
	Annotations map[string]string
}

// Notifier 通知器接口
type Notifier interface {
	Notify(alert Alert) error
}

// Alert 告警
type Alert struct {
	Rule      AlertRule
	Value     float64
	Timestamp time.Time
}

// NewAlertManager 創建告警管理器
func NewAlertManager() *AlertManager {
	return &AlertManager{
		rules:     make([]AlertRule, 0),
		notifiers: make([]Notifier, 0),
	}
}

// AddRule 添加告警規則
func (a *AlertManager) AddRule(rule AlertRule) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.rules = append(a.rules, rule)
}

// AddNotifier 添加通知器
func (a *AlertManager) AddNotifier(notifier Notifier) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.notifiers = append(a.notifiers, notifier)
}

