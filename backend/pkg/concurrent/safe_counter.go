package concurrent

import (
	"sync/atomic"
	"unsafe"
)

// SafeCounter 線程安全的計數器
type SafeCounter struct {
	value int64
}

// NewSafeCounter 創建新的安全計數器
func NewSafeCounter() *SafeCounter {
	return &SafeCounter{}
}

// Increment 增加
func (c *SafeCounter) Increment() int64 {
	return atomic.AddInt64(&c.value, 1)
}

// Decrement 減少
func (c *SafeCounter) Decrement() int64 {
	return atomic.AddInt64(&c.value, -1)
}

// Add 添加指定值
func (c *SafeCounter) Add(delta int64) int64 {
	return atomic.AddInt64(&c.value, delta)
}

// Get 獲取當前值
func (c *SafeCounter) Get() int64 {
	return atomic.LoadInt64(&c.value)
}

// Set 設置值
func (c *SafeCounter) Set(value int64) {
	atomic.StoreInt64(&c.value, value)
}

// CompareAndSwap 比較並交換
func (c *SafeCounter) CompareAndSwap(old, new int64) bool {
	return atomic.CompareAndSwapInt64(&c.value, old, new)
}

// Reset 重置為 0
func (c *SafeCounter) Reset() {
	atomic.StoreInt64(&c.value, 0)
}

// SafeFloat64 線程安全的浮點數
type SafeFloat64 struct {
	value uint64
}

// NewSafeFloat64 創建新的安全浮點數
func NewSafeFloat64(initial float64) *SafeFloat64 {
	sf := &SafeFloat64{}
	sf.Set(initial)
	return sf
}

// Get 獲取值
func (f *SafeFloat64) Get() float64 {
	return float64FromUint64(atomic.LoadUint64(&f.value))
}

// Set 設置值
func (f *SafeFloat64) Set(value float64) {
	atomic.StoreUint64(&f.value, uint64FromFloat64(value))
}

// Add 添加值
func (f *SafeFloat64) Add(delta float64) float64 {
	for {
		old := atomic.LoadUint64(&f.value)
		new := uint64FromFloat64(float64FromUint64(old) + delta)
		if atomic.CompareAndSwapUint64(&f.value, old, new) {
			return float64FromUint64(new)
		}
	}
}

// float64FromUint64 從 uint64 轉換為 float64
func float64FromUint64(u uint64) float64 {
	return *(*float64)(unsafe.Pointer(&u))
}

// uint64FromFloat64 從 float64 轉換為 uint64
func uint64FromFloat64(f float64) uint64 {
	return *(*uint64)(unsafe.Pointer(&f))
}