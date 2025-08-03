package database

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

// QueryOptimizer 查詢優化器
type QueryOptimizer struct {
	db              *gorm.DB
	enableCache     bool
	cacheExpiration time.Duration
	slowQueryTime   time.Duration
}

// NewQueryOptimizer 創建查詢優化器
func NewQueryOptimizer(db *gorm.DB) *QueryOptimizer {
	return &QueryOptimizer{
		db:              db,
		enableCache:     true,
		cacheExpiration: 5 * time.Minute,
		slowQueryTime:   200 * time.Millisecond,
	}
}

// OptimizeQuery 優化查詢
func (o *QueryOptimizer) OptimizeQuery(db *gorm.DB) *gorm.DB {
	// 啟用查詢緩存
	if o.enableCache {
		db = db.Clauses(clause.Locking{Strength: "SHARE"})
	}

	// 啟用慢查詢日誌
	db = db.Session(&gorm.Session{
		Logger: logger.Default.LogMode(logger.Info),
	})

	return db
}

// WithIndex 使用索引提示
func (o *QueryOptimizer) WithIndex(db *gorm.DB, indexName string) *gorm.DB {
	// PostgreSQL doesn't support index hints directly
	// This is more for documentation/clarity
	return db
}

// WithPreload 優化預加載
func (o *QueryOptimizer) WithPreload(db *gorm.DB, preloads ...string) *gorm.DB {
	for _, preload := range preloads {
		// 使用 Join 預加載以減少查詢次數
		if strings.Contains(preload, ".") {
			db = db.Preload(preload)
		} else {
			db = db.Joins(preload)
		}
	}
	return db
}

// WithBatchSize 設置批次大小
func (o *QueryOptimizer) WithBatchSize(db *gorm.DB, size int) *gorm.DB {
	return db.Session(&gorm.Session{
		PrepareStmt:              true,
		SkipDefaultTransaction:   true,
		DisableNestedTransaction: true,
	})
}

// WithPartition 使用分區表
func (o *QueryOptimizer) WithPartition(db *gorm.DB, partition string) *gorm.DB {
	return db.Table(fmt.Sprintf("%s PARTITION (%s)", db.Statement.Table, partition))
}

// QueryAnalyzer 查詢分析器
type QueryAnalyzer struct {
	db *gorm.DB
}

// NewQueryAnalyzer 創建查詢分析器
func NewQueryAnalyzer(db *gorm.DB) *QueryAnalyzer {
	return &QueryAnalyzer{db: db}
}

// Explain 解釋查詢計劃
func (a *QueryAnalyzer) Explain(query *gorm.DB) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	
	sql := query.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx
	})
	
	err := a.db.Raw("EXPLAIN ANALYZE " + sql).Scan(&results).Error
	return results, err
}

// GetSlowQueries 獲取慢查詢
func (a *QueryAnalyzer) GetSlowQueries(threshold time.Duration) ([]SlowQuery, error) {
	var slowQueries []SlowQuery
	
	// PostgreSQL 慢查詢
	err := a.db.Raw(`
		SELECT 
			query,
			calls,
			total_time,
			mean_time,
			max_time
		FROM pg_stat_statements
		WHERE mean_time > ?
		ORDER BY mean_time DESC
		LIMIT 100
	`, threshold.Milliseconds()).Scan(&slowQueries).Error
	
	return slowQueries, err
}

// SlowQuery 慢查詢記錄
type SlowQuery struct {
	Query     string  `json:"query"`
	Calls     int64   `json:"calls"`
	TotalTime float64 `json:"total_time"`
	MeanTime  float64 `json:"mean_time"`
	MaxTime   float64 `json:"max_time"`
}

// IndexAdvisor 索引建議器
type IndexAdvisor struct {
	db *gorm.DB
}

// NewIndexAdvisor 創建索引建議器
func NewIndexAdvisor(db *gorm.DB) *IndexAdvisor {
	return &IndexAdvisor{db: db}
}

// SuggestIndexes 建議索引
func (a *IndexAdvisor) SuggestIndexes(table string) ([]IndexSuggestion, error) {
	var suggestions []IndexSuggestion
	
	// 分析表的查詢模式
	err := a.db.Raw(`
		SELECT 
			schemaname,
			tablename,
			attname,
			n_distinct,
			correlation
		FROM pg_stats
		WHERE tablename = ?
		AND n_distinct > 100
		AND correlation < 0.1
		ORDER BY n_distinct DESC
	`, table).Scan(&suggestions).Error
	
	return suggestions, err
}

// IndexSuggestion 索引建議
type IndexSuggestion struct {
	Schema      string  `json:"schema"`
	Table       string  `json:"table"`
	Column      string  `json:"column"`
	Cardinality float64 `json:"cardinality"`
	Correlation float64 `json:"correlation"`
}

// QueryCache 查詢緩存
type QueryCache struct {
	cache map[string]CacheEntry
	mu    sync.RWMutex
	ttl   time.Duration
}

// CacheEntry 緩存條目
type CacheEntry struct {
	Data      interface{}
	ExpiresAt time.Time
}

// NewQueryCache 創建查詢緩存
func NewQueryCache(ttl time.Duration) *QueryCache {
	cache := &QueryCache{
		cache: make(map[string]CacheEntry),
		ttl:   ttl,
	}
	
	// 啟動清理協程
	go cache.cleanup()
	
	return cache
}

// Get 獲取緩存
func (c *QueryCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	entry, exists := c.cache[key]
	if !exists || time.Now().After(entry.ExpiresAt) {
		return nil, false
	}
	
	return entry.Data, true
}

// Set 設置緩存
func (c *QueryCache) Set(key string, data interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.cache[key] = CacheEntry{
		Data:      data,
		ExpiresAt: time.Now().Add(c.ttl),
	}
}

// Delete 刪除緩存
func (c *QueryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	delete(c.cache, key)
}

// cleanup 清理過期緩存
func (c *QueryCache) cleanup() {
	ticker := time.NewTicker(c.ttl / 2)
	defer ticker.Stop()
	
	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, entry := range c.cache {
			if now.After(entry.ExpiresAt) {
				delete(c.cache, key)
			}
		}
		c.mu.Unlock()
	}
}

// BatchProcessor 批次處理器
type BatchProcessor struct {
	db        *gorm.DB
	batchSize int
}

// NewBatchProcessor 創建批次處理器
func NewBatchProcessor(db *gorm.DB, batchSize int) *BatchProcessor {
	return &BatchProcessor{
		db:        db,
		batchSize: batchSize,
	}
}

// CreateInBatches 批次創建
func (p *BatchProcessor) CreateInBatches(records interface{}) error {
	return p.db.CreateInBatches(records, p.batchSize).Error
}

// UpdateInBatches 批次更新
func (p *BatchProcessor) UpdateInBatches(model interface{}, updates map[string]interface{}, batchSize int) error {
	var offset int
	for {
		result := p.db.Model(model).
			Limit(batchSize).
			Offset(offset).
			Updates(updates)
		
		if result.Error != nil {
			return result.Error
		}
		
		if result.RowsAffected == 0 {
			break
		}
		
		offset += batchSize
	}
	
	return nil
}

// DeleteInBatches 批次刪除
func (p *BatchProcessor) DeleteInBatches(model interface{}, condition interface{}, batchSize int) error {
	var offset int
	for {
		var ids []uint
		err := p.db.Model(model).
			Where(condition).
			Limit(batchSize).
			Offset(offset).
			Pluck("id", &ids).Error
		
		if err != nil {
			return err
		}
		
		if len(ids) == 0 {
			break
		}
		
		err = p.db.Delete(model, ids).Error
		if err != nil {
			return err
		}
		
		offset += batchSize
	}
	
	return nil
}

// OptimizeTable 優化表
func OptimizeTable(db *gorm.DB, tableName string) error {
	// PostgreSQL VACUUM ANALYZE
	return db.Exec(fmt.Sprintf("VACUUM ANALYZE %s", tableName)).Error
}

// CreateIndex 創建索引
func CreateIndex(db *gorm.DB, tableName, indexName string, columns []string) error {
	columnStr := strings.Join(columns, ", ")
	sql := fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON %s (%s)", indexName, tableName, columnStr)
	return db.Exec(sql).Error
}

// CreatePartialIndex 創建部分索引
func CreatePartialIndex(db *gorm.DB, tableName, indexName string, columns []string, condition string) error {
	columnStr := strings.Join(columns, ", ")
	sql := fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON %s (%s) WHERE %s", 
		indexName, tableName, columnStr, condition)
	return db.Exec(sql).Error
}

// CreateCompositeIndex 創建複合索引
func CreateCompositeIndex(db *gorm.DB, tableName, indexName string, columns []string) error {
	return CreateIndex(db, tableName, indexName, columns)
}

// GetTableStatistics 獲取表統計信息
func GetTableStatistics(db *gorm.DB, tableName string) (TableStats, error) {
	var stats TableStats
	
	// 獲取表大小
	err := db.Raw(`
		SELECT 
			pg_size_pretty(pg_total_relation_size(?)) as total_size,
			pg_size_pretty(pg_relation_size(?)) as table_size,
			pg_size_pretty(pg_indexes_size(?)) as indexes_size,
			(SELECT count(*) FROM ` + tableName + `) as row_count
	`, tableName, tableName, tableName).Scan(&stats).Error
	
	return stats, err
}

// TableStats 表統計信息
type TableStats struct {
	TotalSize   string `json:"total_size"`
	TableSize   string `json:"table_size"`
	IndexesSize string `json:"indexes_size"`
	RowCount    int64  `json:"row_count"`
}

