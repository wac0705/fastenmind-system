package resources

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"
	
	"gorm.io/gorm"
)

// DBManager 數據庫連接管理器
type DBManager struct {
	connections map[string]*ManagedDB
	mu          sync.RWMutex
}

// ManagedDB 被管理的數據庫連接
type ManagedDB struct {
	db        *gorm.DB
	sqlDB     *sql.DB
	name      string
	createdAt time.Time
	lastUsed  time.Time
	mu        sync.Mutex
}

// NewDBManager 創建新的數據庫管理器
func NewDBManager() *DBManager {
	return &DBManager{
		connections: make(map[string]*ManagedDB),
	}
}

// RegisterDB 註冊數據庫連接
func (m *DBManager) RegisterDB(name string, db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	
	// 設置連接池參數
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)
	
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.connections[name] = &ManagedDB{
		db:        db,
		sqlDB:     sqlDB,
		name:      name,
		createdAt: time.Now(),
		lastUsed:  time.Now(),
	}
	
	return nil
}

// GetDB 獲取數據庫連接
func (m *DBManager) GetDB(name string) (*gorm.DB, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	conn, exists := m.connections[name]
	if !exists {
		return nil, fmt.Errorf("database connection '%s' not found", name)
	}
	
	conn.mu.Lock()
	conn.lastUsed = time.Now()
	conn.mu.Unlock()
	
	return conn.db, nil
}

// WithTransaction 在事務中執行操作
func (m *DBManager) WithTransaction(ctx context.Context, name string, fn func(*gorm.DB) error) error {
	db, err := m.GetDB(name)
	if err != nil {
		return err
	}
	
	// 開始事務
	tx := db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	
	// 確保事務最終會被處理
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()
	
	// 執行操作
	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}
	
	// 提交事務
	return tx.Commit().Error
}

// HealthCheck 健康檢查
func (m *DBManager) HealthCheck(ctx context.Context) map[string]bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	results := make(map[string]bool)
	
	for name, conn := range m.connections {
		// 使用 Ping 檢查連接
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		err := conn.sqlDB.PingContext(ctx)
		cancel()
		
		results[name] = err == nil
	}
	
	return results
}

// Stats 獲取統計信息
func (m *DBManager) Stats() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	stats := make(map[string]interface{})
	connections := make([]map[string]interface{}, 0, len(m.connections))
	
	for name, conn := range m.connections {
		conn.mu.Lock()
		dbStats := conn.sqlDB.Stats()
		connInfo := map[string]interface{}{
			"name":              name,
			"created_at":        conn.createdAt,
			"last_used":         conn.lastUsed,
			"open_connections":  dbStats.OpenConnections,
			"in_use":            dbStats.InUse,
			"idle":              dbStats.Idle,
			"wait_count":        dbStats.WaitCount,
			"wait_duration":     dbStats.WaitDuration.String(),
			"max_idle_closed":   dbStats.MaxIdleClosed,
			"max_lifetime_closed": dbStats.MaxLifetimeClosed,
		}
		conn.mu.Unlock()
		connections = append(connections, connInfo)
	}
	
	stats["connections"] = connections
	stats["total_connections"] = len(m.connections)
	
	return stats
}

// CloseDB 關閉指定數據庫連接
func (m *DBManager) CloseDB(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	conn, exists := m.connections[name]
	if !exists {
		return nil
	}
	
	err := conn.sqlDB.Close()
	delete(m.connections, name)
	
	return err
}

// CloseAll 關閉所有數據庫連接
func (m *DBManager) CloseAll() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	var firstErr error
	
	for name, conn := range m.connections {
		if err := conn.sqlDB.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
		delete(m.connections, name)
	}
	
	return firstErr
}

// QueryContext 執行查詢（帶上下文和資源管理）
func QueryContext(ctx context.Context, db *gorm.DB, dest interface{}, query string, args ...interface{}) error {
	// 創建一個用於查詢的子上下文
	queryCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	
	// 執行查詢
	result := db.WithContext(queryCtx).Raw(query, args...).Scan(dest)
	
	return result.Error
}

// ExecContext 執行語句（帶上下文和資源管理）
func ExecContext(ctx context.Context, db *gorm.DB, query string, args ...interface{}) error {
	// 創建一個用於執行的子上下文
	execCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	
	// 執行語句
	result := db.WithContext(execCtx).Exec(query, args...)
	
	return result.Error
}

// PreparedStmtManager 預編譯語句管理器
type PreparedStmtManager struct {
	stmts map[string]*sql.Stmt
	db    *sql.DB
	mu    sync.RWMutex
}

// NewPreparedStmtManager 創建預編譯語句管理器
func NewPreparedStmtManager(db *sql.DB) *PreparedStmtManager {
	return &PreparedStmtManager{
		stmts: make(map[string]*sql.Stmt),
		db:    db,
	}
}

// Prepare 準備語句
func (m *PreparedStmtManager) Prepare(name, query string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// 如果已存在，先關閉
	if stmt, exists := m.stmts[name]; exists {
		stmt.Close()
	}
	
	// 準備新語句
	stmt, err := m.db.Prepare(query)
	if err != nil {
		return err
	}
	
	m.stmts[name] = stmt
	return nil
}

// Get 獲取預編譯語句
func (m *PreparedStmtManager) Get(name string) (*sql.Stmt, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	stmt, exists := m.stmts[name]
	if !exists {
		return nil, fmt.Errorf("prepared statement '%s' not found", name)
	}
	
	return stmt, nil
}

// Close 關閉所有預編譯語句
func (m *PreparedStmtManager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	var firstErr error
	
	for name, stmt := range m.stmts {
		if err := stmt.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
		delete(m.stmts, name)
	}
	
	return firstErr
}