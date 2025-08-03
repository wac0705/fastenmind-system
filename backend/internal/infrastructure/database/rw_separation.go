package database

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fastenmind/fastener-api/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/plugin/dbresolver"
)

// ReadWriteDB provides read/write separation for database operations
type ReadWriteDB struct {
	writeDB    *gorm.DB
	readDBs    []*gorm.DB
	roundRobin uint64
}

// NewReadWriteDB creates a new read/write separated database connection
func NewReadWriteDB(writeConfig config.DBConnectionConfig, readConfigs []config.DBConnectionConfig) (*ReadWriteDB, error) {
	// Initialize write database
	writeDB, err := gorm.Open(postgres.Open(writeConfig.DSN()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to write database: %w", err)
	}

	// Initialize read databases
	readDBs := make([]*gorm.DB, 0, len(readConfigs))
	for i, readConfig := range readConfigs {
		readDB, err := gorm.Open(postgres.Open(readConfig.DSN()), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("failed to connect to read database %d: %w", i, err)
		}
		readDBs = append(readDBs, readDB)
	}

	// If no read replicas, use write DB for reads
	if len(readDBs) == 0 {
		readDBs = append(readDBs, writeDB)
	}

	return &ReadWriteDB{
		writeDB: writeDB,
		readDBs: readDBs,
	}, nil
}

// ConfigureDBResolver configures GORM's built-in DB resolver for read/write separation
func ConfigureDBResolver(db *gorm.DB, readConfigs []config.DBConnectionConfig) error {
	// Create read replicas
	replicas := make([]gorm.Dialector, 0, len(readConfigs))
	for _, cfg := range readConfigs {
		dialector := postgres.Open(cfg.DSN())
		replicas = append(replicas, dialector)
	}

	// Register resolver with read replicas
	err := db.Use(
		dbresolver.Register(dbresolver.Config{
			Replicas: replicas,
			Policy:   dbresolver.RandomPolicy{},
		}).
			SetMaxIdleConns(10).
			SetMaxOpenConns(100).
			SetConnMaxLifetime(time.Hour),
	)

	return err
}

// Write returns the write database connection
func (rw *ReadWriteDB) Write() *gorm.DB {
	return rw.writeDB
}

// Read returns a read database connection using round-robin
func (rw *ReadWriteDB) Read() *gorm.DB {
	if len(rw.readDBs) == 1 {
		return rw.readDBs[0]
	}

	// Round-robin selection
	index := atomic.AddUint64(&rw.roundRobin, 1) % uint64(len(rw.readDBs))
	return rw.readDBs[index]
}

// ReadPreferPrimary returns the write DB for consistent reads
func (rw *ReadWriteDB) ReadPreferPrimary() *gorm.DB {
	return rw.writeDB
}

// Transaction executes a function within a database transaction
func (rw *ReadWriteDB) Transaction(fn func(tx *gorm.DB) error) error {
	return rw.writeDB.Transaction(fn)
}

// Close closes all database connections
func (rw *ReadWriteDB) Close() error {
	// Close write DB
	if sqlDB, err := rw.writeDB.DB(); err == nil {
		if err := sqlDB.Close(); err != nil {
			return fmt.Errorf("failed to close write database: %w", err)
		}
	}

	// Close read DBs
	for i, readDB := range rw.readDBs {
		if readDB == rw.writeDB {
			continue // Skip if read DB is same as write DB
		}
		if sqlDB, err := readDB.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				return fmt.Errorf("failed to close read database %d: %w", i, err)
			}
		}
	}

	return nil
}

// ReadWriteMiddleware is a GORM plugin that automatically routes queries
type ReadWriteMiddleware struct {
	rw *ReadWriteDB
}

// NewReadWriteMiddleware creates a new read/write middleware
func NewReadWriteMiddleware(rw *ReadWriteDB) *ReadWriteMiddleware {
	return &ReadWriteMiddleware{rw: rw}
}

// Name returns the plugin name
func (m *ReadWriteMiddleware) Name() string {
	return "read_write_middleware"
}

// Initialize initializes the plugin
func (m *ReadWriteMiddleware) Initialize(db *gorm.DB) error {
	// Register callbacks for read operations
	db.Callback().Query().Before("gorm:query").Register("rw:route_read", m.routeRead)
	db.Callback().Row().Before("gorm:row").Register("rw:route_read", m.routeRead)
	
	// Register callbacks for write operations
	db.Callback().Create().Before("gorm:create").Register("rw:route_write", m.routeWrite)
	db.Callback().Update().Before("gorm:update").Register("rw:route_write", m.routeWrite)
	db.Callback().Delete().Before("gorm:delete").Register("rw:route_write", m.routeWrite)
	
	return nil
}

// routeRead routes read operations to read replicas
func (m *ReadWriteMiddleware) routeRead(db *gorm.DB) {
	// Check if there's a locking clause (FOR UPDATE/SHARE)
	if _, exists := db.Statement.Clauses["FOR"]; exists {
		// FOR UPDATE/SHARE queries should go to primary
		db.Statement.ConnPool = m.rw.ReadPreferPrimary().Statement.ConnPool
		return
	}
	
	db.Statement.ConnPool = m.rw.Read().Statement.ConnPool
}

// routeWrite routes write operations to primary
func (m *ReadWriteMiddleware) routeWrite(db *gorm.DB) {
	db.Statement.ConnPool = m.rw.Write().Statement.ConnPool
}

// ReadWriteRepository provides base repository with read/write separation
type ReadWriteRepository struct {
	rw *ReadWriteDB
}

// NewReadWriteRepository creates a new read/write repository
func NewReadWriteRepository(rw *ReadWriteDB) *ReadWriteRepository {
	return &ReadWriteRepository{rw: rw}
}

// ReadDB returns a database connection for read operations
func (r *ReadWriteRepository) ReadDB() *gorm.DB {
	return r.rw.Read()
}

// WriteDB returns a database connection for write operations
func (r *ReadWriteRepository) WriteDB() *gorm.DB {
	return r.rw.Write()
}

// ConsistentReadDB returns primary DB for consistent reads
func (r *ReadWriteRepository) ConsistentReadDB() *gorm.DB {
	return r.rw.ReadPreferPrimary()
}

// Transaction executes a function within a transaction
func (r *ReadWriteRepository) Transaction(fn func(tx *gorm.DB) error) error {
	return r.rw.Transaction(fn)
}

// QueryOptions provides options for query execution
type QueryOptions struct {
	// UseWriteDB forces the query to use the write database
	UseWriteDB bool
	
	// ConsistentRead ensures read-after-write consistency
	ConsistentRead bool
	
	// LockForUpdate adds FOR UPDATE clause
	LockForUpdate bool
	
	// LockForShare adds FOR SHARE clause
	LockForShare bool
}

// ApplyQueryOptions applies query options to a GORM query
func ApplyQueryOptions(db *gorm.DB, opts QueryOptions) *gorm.DB {
	if opts.LockForUpdate {
		db = db.Clauses(clause.Locking{Strength: "UPDATE"})
	} else if opts.LockForShare {
		db = db.Clauses(clause.Locking{Strength: "SHARE"})
	}
	
	return db
}

// LoadBalancer interface for custom load balancing strategies
type LoadBalancer interface {
	// Next returns the next database connection to use
	Next(connections []*gorm.DB) *gorm.DB
	
	// MarkHealthy marks a connection as healthy
	MarkHealthy(db *gorm.DB)
	
	// MarkUnhealthy marks a connection as unhealthy
	MarkUnhealthy(db *gorm.DB)
}

// WeightedRoundRobinBalancer implements weighted round-robin load balancing
type WeightedRoundRobinBalancer struct {
	weights    map[*gorm.DB]int
	current    map[*gorm.DB]int
	unhealthy  map[*gorm.DB]bool
	mu         sync.RWMutex
	lastIndex  int
}

// NewWeightedRoundRobinBalancer creates a new weighted round-robin balancer
func NewWeightedRoundRobinBalancer(weights map[*gorm.DB]int) *WeightedRoundRobinBalancer {
	current := make(map[*gorm.DB]int)
	for db := range weights {
		current[db] = 0
	}
	
	return &WeightedRoundRobinBalancer{
		weights:   weights,
		current:   current,
		unhealthy: make(map[*gorm.DB]bool),
	}
}

// Next returns the next database connection based on weights
func (b *WeightedRoundRobinBalancer) Next(connections []*gorm.DB) *gorm.DB {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	// Filter healthy connections
	healthy := make([]*gorm.DB, 0)
	for _, conn := range connections {
		if !b.unhealthy[conn] {
			healthy = append(healthy, conn)
		}
	}
	
	if len(healthy) == 0 {
		// All unhealthy, return first connection
		return connections[0]
	}
	
	// Simple weighted round-robin
	var selected *gorm.DB
	maxRatio := -1.0
	
	for _, conn := range healthy {
		weight := b.weights[conn]
		if weight == 0 {
			weight = 1 // Default weight
		}
		
		current := b.current[conn]
		ratio := float64(current) / float64(weight)
		
		if selected == nil || ratio < maxRatio {
			selected = conn
			maxRatio = ratio
		}
	}
	
	b.current[selected]++
	return selected
}

// MarkHealthy marks a connection as healthy
func (b *WeightedRoundRobinBalancer) MarkHealthy(db *gorm.DB) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.unhealthy, db)
}

// MarkUnhealthy marks a connection as unhealthy
func (b *WeightedRoundRobinBalancer) MarkUnhealthy(db *gorm.DB) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.unhealthy[db] = true
}