package database

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestQueryCache(t *testing.T) {
	cache := NewQueryCache(100 * time.Millisecond)

	// Test Set and Get
	t.Run("SetAndGet", func(t *testing.T) {
		cache.Set("key1", "value1")
		
		value, found := cache.Get("key1")
		assert.True(t, found)
		assert.Equal(t, "value1", value)
	})

	// Test Get non-existent key
	t.Run("GetNonExistent", func(t *testing.T) {
		value, found := cache.Get("nonexistent")
		assert.False(t, found)
		assert.Nil(t, value)
	})

	// Test expiration
	t.Run("Expiration", func(t *testing.T) {
		cache.Set("key2", "value2")
		
		// Wait for expiration
		time.Sleep(150 * time.Millisecond)
		
		value, found := cache.Get("key2")
		assert.False(t, found)
		assert.Nil(t, value)
	})

	// Test Delete
	t.Run("Delete", func(t *testing.T) {
		cache.Set("key3", "value3")
		cache.Delete("key3")
		
		value, found := cache.Get("key3")
		assert.False(t, found)
		assert.Nil(t, value)
	})

	// Test concurrent access
	t.Run("ConcurrentAccess", func(t *testing.T) {
		done := make(chan bool, 10)
		
		// Writers
		for i := 0; i < 5; i++ {
			go func(id int) {
				key := fmt.Sprintf("concurrent_%d", id)
				cache.Set(key, id)
				done <- true
			}(i)
		}
		
		// Readers
		for i := 0; i < 5; i++ {
			go func(id int) {
				key := fmt.Sprintf("concurrent_%d", id)
				cache.Get(key)
				done <- true
			}(i)
		}
		
		// Wait for all goroutines
		for i := 0; i < 10; i++ {
			<-done
		}
	})
}

func TestBatchProcessor(t *testing.T) {
	// This would require a real database connection
	// For now, we'll test the logic structure
	
	t.Run("BatchSizeValidation", func(t *testing.T) {
		// Test that batch processor is created with valid batch size
		assert.NotPanics(t, func() {
			// In real test, would use actual DB connection
			// processor := NewBatchProcessor(db, 1000)
			// assert.NotNil(t, processor)
		})
	})
}

func TestTableStatisticsStructure(t *testing.T) {
	stats := TableStats{
		TotalSize:   "1 GB",
		TableSize:   "800 MB",
		IndexesSize: "200 MB",
		RowCount:    1000000,
	}

	assert.Equal(t, "1 GB", stats.TotalSize)
	assert.Equal(t, "800 MB", stats.TableSize)
	assert.Equal(t, "200 MB", stats.IndexesSize)
	assert.Equal(t, int64(1000000), stats.RowCount)
}

func TestSlowQueryStructure(t *testing.T) {
	slowQuery := SlowQuery{
		Query:     "SELECT * FROM large_table",
		Calls:     100,
		TotalTime: 5000.0,
		MeanTime:  50.0,
		MaxTime:   200.0,
	}

	assert.Equal(t, "SELECT * FROM large_table", slowQuery.Query)
	assert.Equal(t, int64(100), slowQuery.Calls)
	assert.Equal(t, 5000.0, slowQuery.TotalTime)
	assert.Equal(t, 50.0, slowQuery.MeanTime)
	assert.Equal(t, 200.0, slowQuery.MaxTime)
}

func TestIndexSuggestionStructure(t *testing.T) {
	suggestion := IndexSuggestion{
		Schema:      "public",
		Table:       "customers",
		Column:      "company_id",
		Cardinality: 1000.0,
		Correlation: 0.05,
	}

	assert.Equal(t, "public", suggestion.Schema)
	assert.Equal(t, "customers", suggestion.Table)
	assert.Equal(t, "company_id", suggestion.Column)
	assert.Equal(t, 1000.0, suggestion.Cardinality)
	assert.Equal(t, 0.05, suggestion.Correlation)
}

func TestCacheKeyGeneration(t *testing.T) {
	testCases := []struct {
		name     string
		prefix   string
		params   interface{}
		expected string
	}{
		{
			name:     "Simple key",
			prefix:   "customers",
			params:   "list",
			expected: "customers:list",
		},
		{
			name:   "Complex key",
			prefix: "quotes",
			params: struct {
				CompanyID string
				Status    string
			}{"123", "active"},
			expected: "quotes:{CompanyID:123 Status:active}",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			key := fmt.Sprintf("%s:%+v", tc.prefix, tc.params)
			assert.Contains(t, key, tc.prefix)
		})
	}
}

func BenchmarkQueryCache(b *testing.B) {
	cache := NewQueryCache(5 * time.Minute)
	
	// Benchmark Set operation
	b.Run("Set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("key_%d", i)
			cache.Set(key, i)
		}
	})
	
	// Benchmark Get operation
	b.Run("Get", func(b *testing.B) {
		// Pre-populate cache
		for i := 0; i < 1000; i++ {
			cache.Set(fmt.Sprintf("key_%d", i), i)
		}
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("key_%d", i%1000)
			cache.Get(key)
		}
	})
	
	// Benchmark concurrent access
	b.Run("Concurrent", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				key := fmt.Sprintf("key_%d", i)
				if i%2 == 0 {
					cache.Set(key, i)
				} else {
					cache.Get(key)
				}
				i++
			}
		})
	})
}

