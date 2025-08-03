package database

import (
	"context"
	"crypto/md5"
	"encoding/binary"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/fastenmind/fastener-api/internal/config"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ShardingStrategy defines how data is distributed across shards
type ShardingStrategy interface {
	// GetShardKey returns the shard key for a given record
	GetShardKey(record interface{}) (string, error)
	
	// GetShardID returns the shard ID for a given shard key
	GetShardID(shardKey string) int
	
	// GetShardForQuery returns the shard(s) to query based on conditions
	GetShardForQuery(conditions map[string]interface{}) []int
}

// ShardManager manages database shards
type ShardManager struct {
	shards    map[int]*Shard
	strategy  ShardingStrategy
	mu        sync.RWMutex
	config    *ShardingConfig
}

// Shard represents a database shard
type Shard struct {
	ID       int
	Name     string
	DB       *gorm.DB
	Config   config.DBConnectionConfig
	Weight   int
	ReadOnly bool
	Active   bool
}

// ShardingConfig holds sharding configuration
type ShardingConfig struct {
	// Number of shards
	NumShards int
	
	// Sharding strategy
	Strategy string // "hash", "range", "geo", "custom"
	
	// Shard key field
	ShardKeyField string
	
	// Replication factor
	ReplicationFactor int
	
	// Auto-rebalancing
	EnableAutoRebalance bool
	RebalanceInterval   time.Duration
	
	// Cross-shard query settings
	EnableCrossShardQuery bool
	MaxParallelQueries    int
}

// NewShardManager creates a new shard manager
func NewShardManager(config *ShardingConfig) *ShardManager {
	sm := &ShardManager{
		shards: make(map[int]*Shard),
		config: config,
	}
	
	// Set strategy based on config
	switch config.Strategy {
	case "hash":
		sm.strategy = NewHashShardingStrategy(config.NumShards, config.ShardKeyField)
	case "range":
		sm.strategy = NewRangeShardingStrategy(config.ShardKeyField)
	case "geo":
		sm.strategy = NewGeoShardingStrategy()
	default:
		sm.strategy = NewHashShardingStrategy(config.NumShards, config.ShardKeyField)
	}
	
	return sm
}

// AddShard adds a new shard
func (sm *ShardManager) AddShard(shard *Shard) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	
	if _, exists := sm.shards[shard.ID]; exists {
		return fmt.Errorf("shard %d already exists", shard.ID)
	}
	
	// Connect to shard database
	dsn := shard.Config.DSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to shard %d: %w", shard.ID, err)
	}
	
	shard.DB = db
	shard.Active = true
	sm.shards[shard.ID] = shard
	
	return nil
}

// RemoveShard removes a shard
func (sm *ShardManager) RemoveShard(shardID int) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	
	shard, exists := sm.shards[shardID]
	if !exists {
		return fmt.Errorf("shard %d not found", shardID)
	}
	
	// Close database connection
	if sqlDB, err := shard.DB.DB(); err == nil {
		sqlDB.Close()
	}
	
	delete(sm.shards, shardID)
	return nil
}

// GetShard returns a specific shard
func (sm *ShardManager) GetShard(shardID int) (*Shard, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	
	shard, exists := sm.shards[shardID]
	if !exists {
		return nil, fmt.Errorf("shard %d not found", shardID)
	}
	
	if !shard.Active {
		return nil, fmt.Errorf("shard %d is not active", shardID)
	}
	
	return shard, nil
}

// GetShardForRecord returns the shard for a specific record
func (sm *ShardManager) GetShardForRecord(record interface{}) (*Shard, error) {
	shardKey, err := sm.strategy.GetShardKey(record)
	if err != nil {
		return nil, err
	}
	
	shardID := sm.strategy.GetShardID(shardKey)
	return sm.GetShard(shardID)
}

// ExecuteOnShard executes a function on a specific shard
func (sm *ShardManager) ExecuteOnShard(ctx context.Context, shardID int, fn func(*gorm.DB) error) error {
	shard, err := sm.GetShard(shardID)
	if err != nil {
		return err
	}
	
	return fn(shard.DB.WithContext(ctx))
}

// ExecuteOnAllShards executes a function on all shards
func (sm *ShardManager) ExecuteOnAllShards(ctx context.Context, fn func(*gorm.DB) error, parallel bool) error {
	sm.mu.RLock()
	shards := make([]*Shard, 0, len(sm.shards))
	for _, shard := range sm.shards {
		if shard.Active {
			shards = append(shards, shard)
		}
	}
	sm.mu.RUnlock()
	
	if parallel {
		return sm.executeParallel(ctx, shards, fn)
	}
	
	return sm.executeSequential(ctx, shards, fn)
}

func (sm *ShardManager) executeParallel(ctx context.Context, shards []*Shard, fn func(*gorm.DB) error) error {
	errCh := make(chan error, len(shards))
	var wg sync.WaitGroup
	
	// Limit parallelism
	semaphore := make(chan struct{}, sm.config.MaxParallelQueries)
	
	for _, shard := range shards {
		wg.Add(1)
		go func(s *Shard) {
			defer wg.Done()
			
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			
			if err := fn(s.DB.WithContext(ctx)); err != nil {
				errCh <- fmt.Errorf("shard %d: %w", s.ID, err)
			}
		}(shard)
	}
	
	wg.Wait()
	close(errCh)
	
	// Collect errors
	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}
	
	if len(errs) > 0 {
		return fmt.Errorf("errors on %d shards: %v", len(errs), errs)
	}
	
	return nil
}

func (sm *ShardManager) executeSequential(ctx context.Context, shards []*Shard, fn func(*gorm.DB) error) error {
	for _, shard := range shards {
		if err := fn(shard.DB.WithContext(ctx)); err != nil {
			return fmt.Errorf("shard %d: %w", shard.ID, err)
		}
	}
	return nil
}

// HashShardingStrategy implements hash-based sharding
type HashShardingStrategy struct {
	numShards     int
	shardKeyField string
}

func NewHashShardingStrategy(numShards int, shardKeyField string) *HashShardingStrategy {
	return &HashShardingStrategy{
		numShards:     numShards,
		shardKeyField: shardKeyField,
	}
}

func (h *HashShardingStrategy) GetShardKey(record interface{}) (string, error) {
	// Use reflection to get shard key field value
	// This is a simplified implementation
	switch r := record.(type) {
	case map[string]interface{}:
		if val, ok := r[h.shardKeyField]; ok {
			return fmt.Sprintf("%v", val), nil
		}
	}
	
	return "", errors.New("shard key not found")
}

func (h *HashShardingStrategy) GetShardID(shardKey string) int {
	// Use consistent hashing
	hash := md5.Sum([]byte(shardKey))
	hashInt := binary.BigEndian.Uint32(hash[:4])
	return int(hashInt % uint32(h.numShards))
}

func (h *HashShardingStrategy) GetShardForQuery(conditions map[string]interface{}) []int {
	// Check if shard key is in conditions
	if shardKeyValue, ok := conditions[h.shardKeyField]; ok {
		shardKey := fmt.Sprintf("%v", shardKeyValue)
		return []int{h.GetShardID(shardKey)}
	}
	
	// Return all shards if shard key not specified
	shards := make([]int, h.numShards)
	for i := range shards {
		shards[i] = i
	}
	return shards
}

// RangeShardingStrategy implements range-based sharding
type RangeShardingStrategy struct {
	shardKeyField string
	ranges        []ShardRange
}

type ShardRange struct {
	ShardID int
	MinKey  interface{}
	MaxKey  interface{}
}

func NewRangeShardingStrategy(shardKeyField string) *RangeShardingStrategy {
	return &RangeShardingStrategy{
		shardKeyField: shardKeyField,
		ranges:        []ShardRange{}, // Initialize with ranges
	}
}

func (r *RangeShardingStrategy) GetShardKey(record interface{}) (string, error) {
	// Implementation similar to hash strategy
	return "", nil
}

func (r *RangeShardingStrategy) GetShardID(shardKey string) int {
	// Find appropriate range
	for _, rng := range r.ranges {
		// Compare shardKey with range bounds
		// Simplified - actual implementation would depend on key type
		return rng.ShardID
	}
	return 0
}

func (r *RangeShardingStrategy) GetShardForQuery(conditions map[string]interface{}) []int {
	// Analyze conditions to determine which ranges to query
	return []int{0}
}

// GeoShardingStrategy implements geography-based sharding
type GeoShardingStrategy struct {
	regions map[string]int // region -> shardID
}

func NewGeoShardingStrategy() *GeoShardingStrategy {
	return &GeoShardingStrategy{
		regions: map[string]int{
			"us-east":  0,
			"us-west":  1,
			"eu-west":  2,
			"asia-pac": 3,
		},
	}
}

func (g *GeoShardingStrategy) GetShardKey(record interface{}) (string, error) {
	// Extract region from record
	return "", nil
}

func (g *GeoShardingStrategy) GetShardID(shardKey string) int {
	if shardID, ok := g.regions[shardKey]; ok {
		return shardID
	}
	return 0 // Default shard
}

func (g *GeoShardingStrategy) GetShardForQuery(conditions map[string]interface{}) []int {
	// Check for region in conditions
	return []int{0}
}

// ShardedRepository provides sharded repository operations
type ShardedRepository struct {
	shardManager *ShardManager
	tableName    string
}

func NewShardedRepository(shardManager *ShardManager, tableName string) *ShardedRepository {
	return &ShardedRepository{
		shardManager: shardManager,
		tableName:    tableName,
	}
}

// Create creates a record in the appropriate shard
func (sr *ShardedRepository) Create(ctx context.Context, record interface{}) error {
	shard, err := sr.shardManager.GetShardForRecord(record)
	if err != nil {
		return err
	}
	
	return shard.DB.WithContext(ctx).Table(sr.tableName).Create(record).Error
}

// FindByID finds a record by ID across all shards
func (sr *ShardedRepository) FindByID(ctx context.Context, id uuid.UUID, dest interface{}) error {
	errCh := make(chan error, 1)
	found := make(chan bool, 1)
	
	// Search in parallel across shards
	err := sr.shardManager.ExecuteOnAllShards(ctx, func(db *gorm.DB) error {
		err := db.Table(sr.tableName).Where("id = ?", id).First(dest).Error
		if err == nil {
			found <- true
			return nil
		}
		if err != gorm.ErrRecordNotFound {
			errCh <- err
		}
		return nil
	}, true)
	
	select {
	case <-found:
		return nil
	case err := <-errCh:
		return err
	default:
		if err != nil {
			return err
		}
		return gorm.ErrRecordNotFound
	}
}

// Query executes a query across relevant shards
func (sr *ShardedRepository) Query(ctx context.Context, conditions map[string]interface{}, dest interface{}) error {
	shardIDs := sr.shardManager.strategy.GetShardForQuery(conditions)
	
	var allResults []interface{}
	var mu sync.Mutex
	
	// Query relevant shards
	var wg sync.WaitGroup
	errCh := make(chan error, len(shardIDs))
	
	for _, shardID := range shardIDs {
		wg.Add(1)
		go func(sid int) {
			defer wg.Done()
			
			shard, err := sr.shardManager.GetShard(sid)
			if err != nil {
				errCh <- err
				return
			}
			
			var results []interface{}
			query := shard.DB.WithContext(ctx).Table(sr.tableName)
			
			// Apply conditions
			for key, value := range conditions {
				query = query.Where(key+" = ?", value)
			}
			
			if err := query.Find(&results).Error; err != nil {
				errCh <- err
				return
			}
			
			mu.Lock()
			allResults = append(allResults, results...)
			mu.Unlock()
		}(shardID)
	}
	
	wg.Wait()
	close(errCh)
	
	// Check for errors
	for err := range errCh {
		return err
	}
	
	// Merge results
	// This is simplified - actual implementation would handle type conversion
	dest = allResults
	
	return nil
}

// ShardMigration handles shard migration operations
type ShardMigration struct {
	shardManager *ShardManager
}

func NewShardMigration(shardManager *ShardManager) *ShardMigration {
	return &ShardMigration{shardManager: shardManager}
}

// MigrateData migrates data between shards
func (sm *ShardMigration) MigrateData(ctx context.Context, fromShardID, toShardID int, filter func(interface{}) bool) error {
	fromShard, err := sm.shardManager.GetShard(fromShardID)
	if err != nil {
		return err
	}
	
	toShard, err := sm.shardManager.GetShard(toShardID)
	if err != nil {
		return err
	}
	
	// This is a simplified implementation
	// Actual implementation would:
	// 1. Read data in batches
	// 2. Apply filter
	// 3. Write to new shard
	// 4. Verify data integrity
	// 5. Delete from old shard
	
	return nil
}

// RebalanceShards rebalances data across shards
func (sm *ShardMigration) RebalanceShards(ctx context.Context) error {
	// Analyze data distribution
	// Identify hot shards
	// Plan migrations
	// Execute migrations
	return nil
}

// ShardMonitor monitors shard health and performance
type ShardMonitor struct {
	shardManager *ShardManager
	metrics      sync.Map
}

type ShardMetrics struct {
	ShardID         int
	ConnectionCount int
	QueryCount      int64
	ErrorCount      int64
	AvgLatency      time.Duration
	DataSize        int64
	LastHealthCheck time.Time
	IsHealthy       bool
}

func NewShardMonitor(shardManager *ShardManager) *ShardMonitor {
	return &ShardMonitor{
		shardManager: shardManager,
	}
}

// HealthCheck performs health check on all shards
func (sm *ShardMonitor) HealthCheck(ctx context.Context) map[int]bool {
	health := make(map[int]bool)
	
	sm.shardManager.mu.RLock()
	defer sm.shardManager.mu.RUnlock()
	
	for shardID, shard := range sm.shardManager.shards {
		// Ping database
		sqlDB, err := shard.DB.DB()
		if err != nil {
			health[shardID] = false
			continue
		}
		
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		err = sqlDB.PingContext(ctx)
		cancel()
		
		health[shardID] = err == nil
		
		// Update metrics
		metrics := &ShardMetrics{
			ShardID:         shardID,
			LastHealthCheck: time.Now(),
			IsHealthy:       err == nil,
		}
		
		sm.metrics.Store(shardID, metrics)
	}
	
	return health
}

// GetMetrics returns metrics for all shards
func (sm *ShardMonitor) GetMetrics() map[int]*ShardMetrics {
	metrics := make(map[int]*ShardMetrics)
	
	sm.metrics.Range(func(key, value interface{}) bool {
		shardID := key.(int)
		metrics[shardID] = value.(*ShardMetrics)
		return true
	})
	
	return metrics
}