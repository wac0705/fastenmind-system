package monitoring

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMetrics(t *testing.T) {
	metrics := NewMetrics()

	// Test HTTP request recording
	t.Run("RecordHTTPRequest", func(t *testing.T) {
		metrics.RecordHTTPRequest("GET", "/api/test", 200, 100*time.Millisecond, 1024)
		// Metrics are recorded asynchronously, so we can't directly assert values
		// In production, we would query Prometheus
	})

	// Test business operation recording
	t.Run("RecordBusinessOperation", func(t *testing.T) {
		metrics.RecordBusinessOperation("create_quote", true)
		metrics.RecordBusinessOperation("create_quote", false)
	})

	// Test business error recording
	t.Run("RecordBusinessError", func(t *testing.T) {
		metrics.RecordBusinessError("create_quote", "validation_error")
		metrics.RecordBusinessError("create_quote", "database_error")
	})

	// Test quotation recording
	t.Run("RecordQuotation", func(t *testing.T) {
		metrics.RecordQuotation("customer1", "approved", 10000.0)
		metrics.RecordQuotation("customer2", "pending", 5000.0)
	})

	// Test DB query recording
	t.Run("RecordDBQuery", func(t *testing.T) {
		metrics.RecordDBQuery("select", "customers", 50*time.Millisecond)
		metrics.RecordDBQuery("insert", "quotes", 100*time.Millisecond)
	})

	// Test cache operations
	t.Run("CacheOperations", func(t *testing.T) {
		metrics.RecordCacheHit("query_cache")
		metrics.RecordCacheMiss("query_cache")
		metrics.CacheEvictions.Inc()
	})

	// Test system metrics update
	t.Run("UpdateSystemMetrics", func(t *testing.T) {
		metrics.UpdateSystemMetrics(1024*1024*100, 25.5, 50, 100)
	})
}

func TestHealthChecker(t *testing.T) {
	checker := NewHealthChecker()

	// Register health checks
	checker.RegisterCheck("database", func(ctx context.Context) error {
		// Simulate database health check
		return nil
	})

	checker.RegisterCheck("cache", func(ctx context.Context) error {
		// Simulate cache health check
		return nil
	})

	// Run health checks
	ctx := context.Background()
	results := checker.CheckHealth(ctx)

	assert.Len(t, results, 2)
	assert.Equal(t, "healthy", results["database"].Status)
	assert.Equal(t, "healthy", results["cache"].Status)
}

func TestStatsCache(t *testing.T) {
	cache := NewStatsCache(100 * time.Millisecond)

	// Test Set and Get
	t.Run("SetAndGet", func(t *testing.T) {
		cache.Set("test_key", "test_value")
		
		value, found := cache.Get("test_key")
		assert.True(t, found)
		assert.Equal(t, "test_value", value)
	})

	// Test expiration
	t.Run("Expiration", func(t *testing.T) {
		cache.Set("expire_key", "expire_value")
		
		// Wait for expiration
		time.Sleep(150 * time.Millisecond)
		
		value, found := cache.Get("expire_key")
		assert.False(t, found)
		assert.Nil(t, value)
	})

	// Test Clear
	t.Run("Clear", func(t *testing.T) {
		cache.Set("clear_key", "clear_value")
		cache.Clear()
		
		value, found := cache.Get("clear_key")
		assert.False(t, found)
		assert.Nil(t, value)
	})
}

func TestAlertManager(t *testing.T) {
	alertManager := NewAlertManager()

	// Add alert rules
	rule := AlertRule{
		Name:       "high_error_rate",
		Expression: "error_rate > 0.05",
		Duration:   5 * time.Minute,
		Severity:   "critical",
		Labels: map[string]string{
			"team": "backend",
		},
		Annotations: map[string]string{
			"summary": "High error rate detected",
		},
	}
	
	alertManager.AddRule(rule)

	// Add notifier (mock)
	mockNotifier := &MockNotifier{}
	alertManager.AddNotifier(mockNotifier)

	// Test that rules and notifiers are added
	assert.Len(t, alertManager.rules, 1)
	assert.Len(t, alertManager.notifiers, 1)
}

// MockNotifier for testing
type MockNotifier struct {
	notifications []Alert
}

func (m *MockNotifier) Notify(alert Alert) error {
	m.notifications = append(m.notifications, alert)
	return nil
}

func BenchmarkMetrics(b *testing.B) {
	metrics := NewMetrics()

	b.Run("RecordHTTPRequest", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			metrics.RecordHTTPRequest("GET", "/api/test", 200, 100*time.Millisecond, 1024)
		}
	})

	b.Run("RecordBusinessOperation", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			metrics.RecordBusinessOperation("test_op", i%2 == 0)
		}
	})

	b.Run("UpdateSystemMetrics", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			metrics.UpdateSystemMetrics(1024*1024*100, 25.5, 50, 100)
		}
	})
}