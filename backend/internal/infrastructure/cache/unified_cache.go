package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// UnifiedCacheStrategy provides a unified caching layer
type UnifiedCacheStrategy struct {
	layers       []CacheLayer
	strategies   map[string]*CacheStrategy
	invalidators map[string][]string // key pattern -> dependent patterns
	mu           sync.RWMutex
	config       *UnifiedCacheConfig
}

// UnifiedCacheConfig holds unified cache configuration
type UnifiedCacheConfig struct {
	// Layer configuration
	EnableMultiLayer bool
	WriteThrough     bool
	ReadThrough      bool
	
	// Consistency
	StrongConsistency bool
	EventualConsistency time.Duration
	
	// Performance
	AsyncWrites      bool
	BatchInvalidation bool
	PreloadEnabled   bool
	
	// Monitoring
	EnableMetrics    bool
	EnableTracing    bool
}

// CacheLayer represents a cache layer (L1, L2, L3)
type CacheLayer struct {
	Name     string
	Cache    Cache
	Priority int
	TTL      time.Duration
	MaxSize  int64
}

// CacheStrategy defines caching strategy for specific data types
type CacheStrategy struct {
	Name              string
	KeyPattern        string
	TTL               time.Duration
	RefreshAhead      bool
	RefreshWindow     time.Duration
	CompressionEnabled bool
	EncryptionEnabled bool
	Tags              []string
	InvalidationRules []InvalidationRule
}

// InvalidationRule defines when to invalidate cache
type InvalidationRule struct {
	Event       string   // "create", "update", "delete"
	Entity      string   // "inquiry", "quote", "order"
	Patterns    []string // cache key patterns to invalidate
	Cascade     bool     // cascade invalidation
}

// NewUnifiedCacheStrategy creates a new unified cache strategy
func NewUnifiedCacheStrategy(config *UnifiedCacheConfig) *UnifiedCacheStrategy {
	return &UnifiedCacheStrategy{
		strategies:   make(map[string]*CacheStrategy),
		invalidators: make(map[string][]string),
		config:       config,
	}
}

// AddLayer adds a cache layer
func (u *UnifiedCacheStrategy) AddLayer(layer CacheLayer) {
	u.mu.Lock()
	defer u.mu.Unlock()
	
	u.layers = append(u.layers, layer)
	// Sort by priority
	u.sortLayers()
}

// RegisterStrategy registers a cache strategy
func (u *UnifiedCacheStrategy) RegisterStrategy(strategy *CacheStrategy) {
	u.mu.Lock()
	defer u.mu.Unlock()
	
	u.strategies[strategy.Name] = strategy
	
	// Build invalidation map
	for _, rule := range strategy.InvalidationRules {
		for _, pattern := range rule.Patterns {
			u.invalidators[pattern] = append(u.invalidators[pattern], strategy.KeyPattern)
		}
	}
}

// Get retrieves value from cache using appropriate strategy
func (u *UnifiedCacheStrategy) Get(ctx context.Context, key string) (interface{}, error) {
	strategy := u.getStrategy(key)
	
	// Try each layer
	for _, layer := range u.layers {
		value, err := layer.Cache.Get(ctx, key)
		if err == nil && value != nil {
			// Found in this layer, populate higher layers if enabled
			if u.config.EnableMultiLayer {
				u.populateHigherLayers(ctx, key, value, layer.Priority)
			}
			
			// Check if refresh-ahead is needed
			if strategy != nil && strategy.RefreshAhead {
				u.checkRefreshAhead(ctx, key, strategy)
			}
			
			return value, nil
		}
	}
	
	// Read-through if enabled
	if u.config.ReadThrough {
		return u.readThrough(ctx, key)
	}
	
	return nil, ErrCacheMiss
}

// Set stores value in cache using appropriate strategy
func (u *UnifiedCacheStrategy) Set(ctx context.Context, key string, value interface{}, ttl ...time.Duration) error {
	strategy := u.getStrategy(key)
	
	// Determine TTL
	var cacheTTL time.Duration
	if len(ttl) > 0 {
		cacheTTL = ttl[0]
	} else if strategy != nil {
		cacheTTL = strategy.TTL
	} else {
		cacheTTL = 1 * time.Hour // default
	}
	
	// Apply compression if enabled
	if strategy != nil && strategy.CompressionEnabled {
		value = u.compress(value)
	}
	
	// Apply encryption if enabled
	if strategy != nil && strategy.EncryptionEnabled {
		value = u.encrypt(value)
	}
	
	// Write to all layers
	if u.config.WriteThrough {
		return u.writeThrough(ctx, key, value, cacheTTL)
	}
	
	// Write to primary layer only
	if len(u.layers) > 0 {
		return u.layers[0].Cache.Set(ctx, key, value, cacheTTL)
	}
	
	return errors.New("no cache layers configured")
}

// Delete removes value from all cache layers
func (u *UnifiedCacheStrategy) Delete(ctx context.Context, key string) error {
	var lastErr error
	
	// Delete from all layers
	for _, layer := range u.layers {
		if err := layer.Cache.Delete(ctx, key); err != nil {
			lastErr = err
		}
	}
	
	// Cascade invalidation if needed
	u.cascadeInvalidation(ctx, key)
	
	return lastErr
}

// InvalidatePattern invalidates all keys matching pattern
func (u *UnifiedCacheStrategy) InvalidatePattern(ctx context.Context, pattern string) error {
	// Get dependent patterns
	u.mu.RLock()
	dependents := u.invalidators[pattern]
	u.mu.RUnlock()
	
	patterns := append([]string{pattern}, dependents...)
	
	if u.config.BatchInvalidation {
		return u.batchInvalidate(ctx, patterns)
	}
	
	// Invalidate each pattern
	for _, p := range patterns {
		for _, layer := range u.layers {
			if invalidator, ok := layer.Cache.(PatternInvalidator); ok {
				invalidator.InvalidatePattern(ctx, p)
			}
		}
	}
	
	return nil
}

// InvalidateByTags invalidates all entries with specified tags
func (u *UnifiedCacheStrategy) InvalidateByTags(ctx context.Context, tags ...string) error {
	for _, layer := range u.layers {
		if tagInvalidator, ok := layer.Cache.(TagInvalidator); ok {
			if err := tagInvalidator.InvalidateByTags(ctx, tags...); err != nil {
				return err
			}
		}
	}
	return nil
}

// Helper methods

func (u *UnifiedCacheStrategy) sortLayers() {
	// Sort layers by priority (lower number = higher priority)
	for i := 0; i < len(u.layers)-1; i++ {
		for j := 0; j < len(u.layers)-i-1; j++ {
			if u.layers[j].Priority > u.layers[j+1].Priority {
				u.layers[j], u.layers[j+1] = u.layers[j+1], u.layers[j]
			}
		}
	}
}

func (u *UnifiedCacheStrategy) getStrategy(key string) *CacheStrategy {
	u.mu.RLock()
	defer u.mu.RUnlock()
	
	// Find matching strategy by key pattern
	for _, strategy := range u.strategies {
		if matchesPattern(key, strategy.KeyPattern) {
			return strategy
		}
	}
	
	return nil
}

func (u *UnifiedCacheStrategy) populateHigherLayers(ctx context.Context, key string, value interface{}, foundPriority int) {
	// Populate cache layers with higher priority (lower number)
	for _, layer := range u.layers {
		if layer.Priority < foundPriority {
			go layer.Cache.Set(ctx, key, value, layer.TTL)
		}
	}
}

func (u *UnifiedCacheStrategy) checkRefreshAhead(ctx context.Context, key string, strategy *CacheStrategy) {
	// Check if key is within refresh window
	for _, layer := range u.layers {
		if ttlChecker, ok := layer.Cache.(TTLChecker); ok {
			ttl, err := ttlChecker.TTL(ctx, key)
			if err == nil && ttl < strategy.RefreshWindow {
				// Trigger async refresh
				go u.refreshKey(ctx, key)
			}
			break
		}
	}
}

func (u *UnifiedCacheStrategy) readThrough(ctx context.Context, key string) (interface{}, error) {
	// This would be implemented with actual data loaders
	return nil, ErrCacheMiss
}

func (u *UnifiedCacheStrategy) writeThrough(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if u.config.AsyncWrites {
		// Async write to all layers
		var wg sync.WaitGroup
		errCh := make(chan error, len(u.layers))
		
		for _, layer := range u.layers {
			wg.Add(1)
			go func(l CacheLayer) {
				defer wg.Done()
				if err := l.Cache.Set(ctx, key, value, ttl); err != nil {
					errCh <- err
				}
			}(layer)
		}
		
		wg.Wait()
		close(errCh)
		
		// Return first error if any
		for err := range errCh {
			return err
		}
		
		return nil
	}
	
	// Sync write to all layers
	for _, layer := range u.layers {
		if err := layer.Cache.Set(ctx, key, value, ttl); err != nil {
			return err
		}
	}
	
	return nil
}

func (u *UnifiedCacheStrategy) cascadeInvalidation(ctx context.Context, key string) {
	strategy := u.getStrategy(key)
	if strategy == nil {
		return
	}
	
	for _, rule := range strategy.InvalidationRules {
		if rule.Cascade {
			for _, pattern := range rule.Patterns {
				u.InvalidatePattern(ctx, pattern)
			}
		}
	}
}

func (u *UnifiedCacheStrategy) batchInvalidate(ctx context.Context, patterns []string) error {
	// Group patterns by layer
	for _, layer := range u.layers {
		if batchInvalidator, ok := layer.Cache.(BatchInvalidator); ok {
			if err := batchInvalidator.BatchInvalidate(ctx, patterns); err != nil {
				return err
			}
		}
	}
	return nil
}

func (u *UnifiedCacheStrategy) refreshKey(ctx context.Context, key string) {
	// Implement actual refresh logic
	// This would fetch fresh data and update cache
}

func (u *UnifiedCacheStrategy) compress(value interface{}) interface{} {
	// Implement compression
	return value
}

func (u *UnifiedCacheStrategy) encrypt(value interface{}) interface{} {
	// Implement encryption
	return value
}

// Interfaces for advanced cache features

type PatternInvalidator interface {
	InvalidatePattern(ctx context.Context, pattern string) error
}

type TagInvalidator interface {
	InvalidateByTags(ctx context.Context, tags ...string) error
}

type TTLChecker interface {
	TTL(ctx context.Context, key string) (time.Duration, error)
}

type BatchInvalidator interface {
	BatchInvalidate(ctx context.Context, patterns []string) error
}

// PredefinedStrategies provides common cache strategies

func InquiryCacheStrategy() *CacheStrategy {
	return &CacheStrategy{
		Name:               "inquiry",
		KeyPattern:         "inquiry:*",
		TTL:                15 * time.Minute,
		RefreshAhead:       true,
		RefreshWindow:      2 * time.Minute,
		CompressionEnabled: false,
		Tags:               []string{"inquiry", "business"},
		InvalidationRules: []InvalidationRule{
			{
				Event:    "update",
				Entity:   "inquiry",
				Patterns: []string{"inquiry:*", "inquiry:list:*"},
				Cascade:  true,
			},
			{
				Event:    "create",
				Entity:   "quote",
				Patterns: []string{"inquiry:%s:quotes"},
			},
		},
	}
}

func QuoteCacheStrategy() *CacheStrategy {
	return &CacheStrategy{
		Name:               "quote",
		KeyPattern:         "quote:*",
		TTL:                30 * time.Minute,
		RefreshAhead:       true,
		RefreshWindow:      5 * time.Minute,
		CompressionEnabled: true,
		Tags:               []string{"quote", "business"},
		InvalidationRules: []InvalidationRule{
			{
				Event:    "update",
				Entity:   "quote",
				Patterns: []string{"quote:*", "quote:list:*"},
				Cascade:  true,
			},
			{
				Event:    "approve",
				Entity:   "quote",
				Patterns: []string{"quote:%s", "order:pending:*"},
			},
		},
	}
}

func UserCacheStrategy() *CacheStrategy {
	return &CacheStrategy{
		Name:               "user",
		KeyPattern:         "user:*",
		TTL:                1 * time.Hour,
		RefreshAhead:       false,
		CompressionEnabled: false,
		EncryptionEnabled:  true,
		Tags:               []string{"user", "auth"},
		InvalidationRules: []InvalidationRule{
			{
				Event:    "update",
				Entity:   "user",
				Patterns: []string{"user:*", "user:permissions:*"},
				Cascade:  true,
			},
			{
				Event:    "logout",
				Entity:   "session",
				Patterns: []string{"user:%s:session:*"},
			},
		},
	}
}

func ConfigCacheStrategy() *CacheStrategy {
	return &CacheStrategy{
		Name:               "config",
		KeyPattern:         "config:*",
		TTL:                24 * time.Hour,
		RefreshAhead:       true,
		RefreshWindow:      1 * time.Hour,
		CompressionEnabled: false,
		Tags:               []string{"config", "system"},
		InvalidationRules: []InvalidationRule{
			{
				Event:    "update",
				Entity:   "config",
				Patterns: []string{"config:*"},
				Cascade:  true,
			},
		},
	}
}

// CacheWarmer provides cache warming functionality
type CacheWarmer struct {
	cache      *UnifiedCacheStrategy
	loaders    map[string]DataLoader
	schedule   map[string]time.Duration
	mu         sync.RWMutex
	stopCh     chan struct{}
}

type DataLoader func(ctx context.Context) (map[string]interface{}, error)

// NewCacheWarmer creates a new cache warmer
func NewCacheWarmer(cache *UnifiedCacheStrategy) *CacheWarmer {
	return &CacheWarmer{
		cache:    cache,
		loaders:  make(map[string]DataLoader),
		schedule: make(map[string]time.Duration),
		stopCh:   make(chan struct{}),
	}
}

// RegisterLoader registers a data loader
func (cw *CacheWarmer) RegisterLoader(name string, loader DataLoader, interval time.Duration) {
	cw.mu.Lock()
	defer cw.mu.Unlock()
	
	cw.loaders[name] = loader
	cw.schedule[name] = interval
}

// Start starts the cache warmer
func (cw *CacheWarmer) Start(ctx context.Context) {
	cw.mu.RLock()
	defer cw.mu.RUnlock()
	
	for name, interval := range cw.schedule {
		go cw.runLoader(ctx, name, interval)
	}
}

// Stop stops the cache warmer
func (cw *CacheWarmer) Stop() {
	close(cw.stopCh)
}

func (cw *CacheWarmer) runLoader(ctx context.Context, name string, interval time.Duration) {
	// Initial load
	cw.loadData(ctx, name)
	
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			cw.loadData(ctx, name)
		case <-cw.stopCh:
			return
		case <-ctx.Done():
			return
		}
	}
}

func (cw *CacheWarmer) loadData(ctx context.Context, name string) {
	cw.mu.RLock()
	loader, exists := cw.loaders[name]
	cw.mu.RUnlock()
	
	if !exists {
		return
	}
	
	data, err := loader(ctx)
	if err != nil {
		// Log error
		return
	}
	
	// Warm cache with loaded data
	for key, value := range data {
		cw.cache.Set(ctx, key, value)
	}
}

// Helper functions

func matchesPattern(key, pattern string) bool {
	// Simple pattern matching
	// Implement more sophisticated matching if needed
	if pattern == "*" {
		return true
	}
	
	if strings.HasSuffix(pattern, "*") {
		prefix := strings.TrimSuffix(pattern, "*")
		return strings.HasPrefix(key, prefix)
	}
	
	return key == pattern
}

// CacheKeyBuilder provides consistent cache key generation
type CacheKeyBuilder struct {
	prefix string
}

func NewCacheKeyBuilder(prefix string) *CacheKeyBuilder {
	return &CacheKeyBuilder{prefix: prefix}
}

func (ckb *CacheKeyBuilder) Build(parts ...interface{}) string {
	keyParts := []string{ckb.prefix}
	
	for _, part := range parts {
		switch v := part.(type) {
		case string:
			keyParts = append(keyParts, v)
		case uuid.UUID:
			keyParts = append(keyParts, v.String())
		case int, int32, int64:
			keyParts = append(keyParts, fmt.Sprintf("%d", v))
		default:
			keyParts = append(keyParts, fmt.Sprintf("%v", v))
		}
	}
	
	return strings.Join(keyParts, ":")
}

// Common cache key builders
var (
	InquiryKeyBuilder  = NewCacheKeyBuilder("inquiry")
	QuoteKeyBuilder    = NewCacheKeyBuilder("quote")
	OrderKeyBuilder    = NewCacheKeyBuilder("order")
	UserKeyBuilder     = NewCacheKeyBuilder("user")
	ConfigKeyBuilder   = NewCacheKeyBuilder("config")
)