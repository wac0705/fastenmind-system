package concurrent

import (
	"sync"
)

// SafeMap 線程安全的 map
type SafeMap[K comparable, V any] struct {
	mu    sync.RWMutex
	items map[K]V
}

// NewSafeMap 創建新的安全 map
func NewSafeMap[K comparable, V any]() *SafeMap[K, V] {
	return &SafeMap[K, V]{
		items: make(map[K]V),
	}
}

// Set 設置值
func (m *SafeMap[K, V]) Set(key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.items[key] = value
}

// Get 獲取值
func (m *SafeMap[K, V]) Get(key K) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, exists := m.items[key]
	return value, exists
}

// Delete 刪除值
func (m *SafeMap[K, V]) Delete(key K) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.items, key)
}

// Len 獲取長度
func (m *SafeMap[K, V]) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.items)
}

// Range 遍歷所有項目
func (m *SafeMap[K, V]) Range(f func(key K, value V) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	for k, v := range m.items {
		if !f(k, v) {
			break
		}
	}
}

// Keys 返回所有鍵
func (m *SafeMap[K, V]) Keys() []K {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	keys := make([]K, 0, len(m.items))
	for k := range m.items {
		keys = append(keys, k)
	}
	return keys
}

// Values 返回所有值
func (m *SafeMap[K, V]) Values() []V {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	values := make([]V, 0, len(m.items))
	for _, v := range m.items {
		values = append(values, v)
	}
	return values
}

// Clear 清空 map
func (m *SafeMap[K, V]) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.items = make(map[K]V)
}

// LoadOrStore 載入或存儲
func (m *SafeMap[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if existing, ok := m.items[key]; ok {
		return existing, true
	}
	
	m.items[key] = value
	return value, false
}

// LoadAndDelete 載入並刪除
func (m *SafeMap[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	value, loaded = m.items[key]
	if loaded {
		delete(m.items, key)
	}
	return
}

// CompareAndSwap 比較並交換
func (m *SafeMap[K, V]) CompareAndSwap(key K, old, new V, equal func(V, V) bool) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	existing, ok := m.items[key]
	if !ok {
		return false
	}
	
	if equal(existing, old) {
		m.items[key] = new
		return true
	}
	
	return false
}

// Swap 交換值
func (m *SafeMap[K, V]) Swap(key K, value V) (previous V, loaded bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	previous, loaded = m.items[key]
	m.items[key] = value
	return
}

// Copy 複製 map
func (m *SafeMap[K, V]) Copy() map[K]V {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	copy := make(map[K]V, len(m.items))
	for k, v := range m.items {
		copy[k] = v
	}
	return copy
}