package query

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"go.uber.org/zap"
)

// InMemoryQueryBus 內存查詢總線實現
type InMemoryQueryBus struct {
	handlers    map[string]interface{}
	middlewares []Middleware
	logger      *zap.Logger
	mu          sync.RWMutex
}

// NewInMemoryQueryBus 創建內存查詢總線
func NewInMemoryQueryBus(logger *zap.Logger) *InMemoryQueryBus {
	return &InMemoryQueryBus{
		handlers:    make(map[string]interface{}),
		middlewares: make([]Middleware, 0),
		logger:      logger,
	}
}

// Register 註冊查詢處理器
func (b *InMemoryQueryBus) Register(queryType string, handler interface{}) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	if _, exists := b.handlers[queryType]; exists {
		return fmt.Errorf("handler for query type %s already registered", queryType)
	}
	
	// 驗證處理器類型
	handlerType := reflect.TypeOf(handler)
	if handlerType.Kind() != reflect.Func {
		return errors.New("handler must be a function")
	}
	
	// 驗證函數簽名
	if handlerType.NumIn() != 2 {
		return errors.New("handler must have exactly 2 parameters: context and query")
	}
	
	if handlerType.In(0) != reflect.TypeOf((*context.Context)(nil)).Elem() {
		return errors.New("handler's first parameter must be context.Context")
	}
	
	// 驗證返回值
	if handlerType.NumOut() != 2 {
		return errors.New("handler must return exactly 2 values: result and error")
	}
	
	if !handlerType.Out(1).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		return errors.New("handler's second return value must be error")
	}
	
	b.handlers[queryType] = handler
	b.logger.Info("Query handler registered", zap.String("query_type", queryType))
	
	return nil
}

// Send 發送查詢
func (b *InMemoryQueryBus) Send(ctx context.Context, q Query) (interface{}, error) {
	b.mu.RLock()
	handler, exists := b.handlers[q.GetQueryType()]
	b.mu.RUnlock()
	
	if !exists {
		return nil, fmt.Errorf("no handler registered for query type: %s", q.GetQueryType())
	}
	
	// 記錄查詢執行
	b.logger.Info("Executing query",
		zap.String("query_id", q.GetQueryID().String()),
		zap.String("query_type", q.GetQueryType()))
	
	// 執行處理器
	handlerValue := reflect.ValueOf(handler)
	queryValue := reflect.ValueOf(q)
	ctxValue := reflect.ValueOf(ctx)
	
	// 調用處理器
	results := handlerValue.Call([]reflect.Value{ctxValue, queryValue})
	
	// 處理返回值
	if len(results) != 2 {
		return nil, errors.New("handler returned unexpected number of values")
	}
	
	// 檢查錯誤
	if !results[1].IsNil() {
		err := results[1].Interface().(error)
		b.logger.Error("Query execution failed",
			zap.String("query_id", q.GetQueryID().String()),
			zap.String("query_type", q.GetQueryType()),
			zap.Error(err))
		return nil, err
	}
	
	result := results[0].Interface()
	
	b.logger.Info("Query executed successfully",
		zap.String("query_id", q.GetQueryID().String()),
		zap.String("query_type", q.GetQueryType()))
	
	return result, nil
}

// Use 添加中間件
func (b *InMemoryQueryBus) Use(middleware Middleware) {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	b.middlewares = append(b.middlewares, middleware)
}

// CachingMiddleware 緩存中間件
func CachingMiddleware(cache Cache, ttl time.Duration) Middleware {
	return func(next interface{}) interface{} {
		return func(ctx context.Context, q Query) (interface{}, error) {
			// 生成緩存鍵
			cacheKey := fmt.Sprintf("query:%s:%v", q.GetQueryType(), q)
			
			// 嘗試從緩存獲取
			if cached, found := cache.Get(cacheKey); found {
				return cached, nil
			}
			
			// 執行查詢
			handlerFunc := next.(func(context.Context, Query) (interface{}, error))
			result, err := handlerFunc(ctx, q)
			if err != nil {
				return nil, err
			}
			
			// 存入緩存
			cache.Set(cacheKey, result, ttl)
			
			return result, nil
		}
	}
}

// LoggingMiddleware 日誌中間件
func LoggingMiddleware(logger *zap.Logger) Middleware {
	return func(next interface{}) interface{} {
		return func(ctx context.Context, q Query) (interface{}, error) {
			logger.Info("Query received",
				zap.String("query_id", q.GetQueryID().String()),
				zap.String("query_type", q.GetQueryType()),
				zap.Time("timestamp", q.GetTimestamp()))
			
			start := time.Now()
			
			handlerFunc := next.(func(context.Context, Query) (interface{}, error))
			result, err := handlerFunc(ctx, q)
			
			duration := time.Since(start)
			
			if err != nil {
				logger.Error("Query failed",
					zap.String("query_id", q.GetQueryID().String()),
					zap.String("query_type", q.GetQueryType()),
					zap.Duration("duration", duration),
					zap.Error(err))
			} else {
				logger.Info("Query completed",
					zap.String("query_id", q.GetQueryID().String()),
					zap.String("query_type", q.GetQueryType()),
					zap.Duration("duration", duration))
			}
			
			return result, err
		}
	}
}

// MetricsMiddleware 指標中間件
func MetricsMiddleware(recordMetric func(queryType string, success bool, duration float64)) Middleware {
	return func(next interface{}) interface{} {
		return func(ctx context.Context, q Query) (interface{}, error) {
			start := time.Now()
			
			handlerFunc := next.(func(context.Context, Query) (interface{}, error))
			result, err := handlerFunc(ctx, q)
			
			duration := time.Since(start).Seconds()
			success := err == nil
			
			recordMetric(q.GetQueryType(), success, duration)
			
			return result, err
		}
	}
}

// Cache 緩存接口
type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, ttl time.Duration)
	Delete(key string)
	Clear()
}

// InMemoryCache 內存緩存實現
type InMemoryCache struct {
	data map[string]cacheEntry
	mu   sync.RWMutex
}

type cacheEntry struct {
	value     interface{}
	expiresAt time.Time
}

// NewInMemoryCache 創建內存緩存
func NewInMemoryCache() Cache {
	cache := &InMemoryCache{
		data: make(map[string]cacheEntry),
	}
	
	// 啟動清理協程
	go cache.cleanup()
	
	return cache
}

// Get 獲取緩存
func (c *InMemoryCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	entry, exists := c.data[key]
	if !exists || time.Now().After(entry.expiresAt) {
		return nil, false
	}
	
	return entry.value, true
}

// Set 設置緩存
func (c *InMemoryCache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.data[key] = cacheEntry{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
}

// Delete 刪除緩存
func (c *InMemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	delete(c.data, key)
}

// Clear 清空緩存
func (c *InMemoryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.data = make(map[string]cacheEntry)
}

// cleanup 清理過期緩存
func (c *InMemoryCache) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, entry := range c.data {
			if now.After(entry.expiresAt) {
				delete(c.data, key)
			}
		}
		c.mu.Unlock()
	}
}