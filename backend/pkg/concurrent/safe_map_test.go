package concurrent

import (
	"sync"
	"testing"
)

func TestSafeMap_BasicOperations(t *testing.T) {
	m := NewSafeMap[string, int]()

	// Test Set and Get
	m.Set("key1", 100)
	value, exists := m.Get("key1")
	if !exists {
		t.Error("Expected key1 to exist")
	}
	if value != 100 {
		t.Errorf("Expected value 100, got %d", value)
	}

	// Test Get non-existent key
	_, exists = m.Get("key2")
	if exists {
		t.Error("Expected key2 to not exist")
	}

	// Test Delete
	m.Delete("key1")
	_, exists = m.Get("key1")
	if exists {
		t.Error("Expected key1 to be deleted")
	}

	// Test Len
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)
	if m.Len() != 3 {
		t.Errorf("Expected length 3, got %d", m.Len())
	}
}

func TestSafeMap_ConcurrentAccess(t *testing.T) {
	m := NewSafeMap[int, string]()
	var wg sync.WaitGroup
	
	// Number of goroutines
	numGoroutines := 100
	numOperations := 1000

	// Concurrent writes
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := id*numOperations + j
				m.Set(key, "value")
			}
		}(i)
	}
	wg.Wait()

	// Check all values were set
	expectedLen := numGoroutines * numOperations
	if m.Len() != expectedLen {
		t.Errorf("Expected length %d, got %d", expectedLen, m.Len())
	}

	// Concurrent reads
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := id*numOperations + j
				_, exists := m.Get(key)
				if !exists {
					t.Errorf("Key %d should exist", key)
				}
			}
		}(i)
	}
	wg.Wait()

	// Concurrent deletes
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := id*numOperations + j
				m.Delete(key)
			}
		}(i)
	}
	wg.Wait()

	// Check all values were deleted
	if m.Len() != 0 {
		t.Errorf("Expected length 0, got %d", m.Len())
	}
}

func TestSafeMap_Range(t *testing.T) {
	m := NewSafeMap[string, int]()
	expected := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
		"e": 5,
	}

	// Set values
	for k, v := range expected {
		m.Set(k, v)
	}

	// Test Range
	count := 0
	m.Range(func(key string, value int) bool {
		count++
		if expected[key] != value {
			t.Errorf("Expected %s=%d, got %s=%d", key, expected[key], key, value)
		}
		return true
	})

	if count != len(expected) {
		t.Errorf("Expected to iterate over %d items, got %d", len(expected), count)
	}

	// Test Range with early termination
	count = 0
	m.Range(func(key string, value int) bool {
		count++
		return count < 3 // Stop after 3 iterations
	})

	if count != 3 {
		t.Errorf("Expected to iterate over 3 items, got %d", count)
	}
}

func TestSafeMap_LoadOrStore(t *testing.T) {
	m := NewSafeMap[string, int]()

	// First call should store
	actual, loaded := m.LoadOrStore("key1", 100)
	if loaded {
		t.Error("Expected loaded to be false on first call")
	}
	if actual != 100 {
		t.Errorf("Expected actual value 100, got %d", actual)
	}

	// Second call should load
	actual, loaded = m.LoadOrStore("key1", 200)
	if !loaded {
		t.Error("Expected loaded to be true on second call")
	}
	if actual != 100 {
		t.Errorf("Expected actual value 100, got %d", actual)
	}
}

func TestSafeMap_CompareAndSwap(t *testing.T) {
	m := NewSafeMap[string, int]()
	m.Set("key1", 100)

	equal := func(a, b int) bool { return a == b }

	// Successful swap
	swapped := m.CompareAndSwap("key1", 100, 200, equal)
	if !swapped {
		t.Error("Expected successful swap")
	}
	value, _ := m.Get("key1")
	if value != 200 {
		t.Errorf("Expected value 200 after swap, got %d", value)
	}

	// Failed swap (old value doesn't match)
	swapped = m.CompareAndSwap("key1", 100, 300, equal)
	if swapped {
		t.Error("Expected failed swap")
	}
	value, _ = m.Get("key1")
	if value != 200 {
		t.Errorf("Expected value 200 (unchanged), got %d", value)
	}

	// Failed swap (key doesn't exist)
	swapped = m.CompareAndSwap("key2", 0, 100, equal)
	if swapped {
		t.Error("Expected failed swap for non-existent key")
	}
}

func TestSafeMap_Keys(t *testing.T) {
	m := NewSafeMap[string, int]()
	expected := []string{"a", "b", "c"}

	for _, k := range expected {
		m.Set(k, 1)
	}

	keys := m.Keys()
	if len(keys) != len(expected) {
		t.Errorf("Expected %d keys, got %d", len(expected), len(keys))
	}

	// Check all expected keys are present
	keyMap := make(map[string]bool)
	for _, k := range keys {
		keyMap[k] = true
	}
	for _, k := range expected {
		if !keyMap[k] {
			t.Errorf("Expected key %s not found", k)
		}
	}
}

func TestSafeMap_Values(t *testing.T) {
	m := NewSafeMap[string, int]()
	expected := []int{1, 2, 3}

	for i, v := range expected {
		m.Set(string(rune('a'+i)), v)
	}

	values := m.Values()
	if len(values) != len(expected) {
		t.Errorf("Expected %d values, got %d", len(expected), len(values))
	}

	// Check all expected values are present
	valueMap := make(map[int]bool)
	for _, v := range values {
		valueMap[v] = true
	}
	for _, v := range expected {
		if !valueMap[v] {
			t.Errorf("Expected value %d not found", v)
		}
	}
}

func TestSafeMap_Clear(t *testing.T) {
	m := NewSafeMap[string, int]()

	// Add some values
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)

	// Clear the map
	m.Clear()

	if m.Len() != 0 {
		t.Errorf("Expected length 0 after clear, got %d", m.Len())
	}

	// Check values are gone
	_, exists := m.Get("a")
	if exists {
		t.Error("Expected key 'a' to not exist after clear")
	}
}

func BenchmarkSafeMap_Set(b *testing.B) {
	m := NewSafeMap[int, int]()
	b.ResetTimer()
	
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			m.Set(i, i)
			i++
		}
	})
}

func BenchmarkSafeMap_Get(b *testing.B) {
	m := NewSafeMap[int, int]()
	
	// Pre-populate the map
	for i := 0; i < 10000; i++ {
		m.Set(i, i)
	}
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			m.Get(i % 10000)
			i++
		}
	})
}