package resources

import (
	"context"
	"sync"
	"time"
)

// ResourceManager 統一資源管理器
type ResourceManager struct {
	httpManager *HTTPClientManager
	fileManager *FileManager
	dbManager   *DBManager
	
	cleanupInterval time.Duration
	cleanupStop     chan struct{}
	cleanupWg       sync.WaitGroup
	
	mu sync.RWMutex
}

// NewResourceManager 創建新的資源管理器
func NewResourceManager() *ResourceManager {
	return &ResourceManager{
		httpManager:     NewHTTPClientManager(),
		fileManager:     NewFileManager(),
		dbManager:       NewDBManager(),
		cleanupInterval: 5 * time.Minute,
		cleanupStop:     make(chan struct{}),
	}
}

// HTTPClients 獲取 HTTP 客戶端管理器
func (m *ResourceManager) HTTPClients() *HTTPClientManager {
	return m.httpManager
}

// Files 獲取文件管理器
func (m *ResourceManager) Files() *FileManager {
	return m.fileManager
}

// Databases 獲取數據庫管理器
func (m *ResourceManager) Databases() *DBManager {
	return m.dbManager
}

// StartCleanup 啟動定期清理
func (m *ResourceManager) StartCleanup(ctx context.Context) {
	m.cleanupWg.Add(1)
	go m.cleanupRoutine(ctx)
}

// cleanupRoutine 清理例程
func (m *ResourceManager) cleanupRoutine(ctx context.Context) {
	defer m.cleanupWg.Done()
	
	ticker := time.NewTicker(m.cleanupInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-m.cleanupStop:
			return
		case <-ticker.C:
			m.performCleanup()
		}
	}
}

// performCleanup 執行清理
func (m *ResourceManager) performCleanup() {
	// 清理空閒的 HTTP 客戶端
	m.httpManager.CleanupIdleClients(10 * time.Minute)
	
	// 清理空閒的文件句柄
	m.fileManager.CleanupIdleFiles(5 * time.Minute)
}

// Shutdown 關閉資源管理器
func (m *ResourceManager) Shutdown(ctx context.Context) error {
	// 停止清理例程
	close(m.cleanupStop)
	
	// 等待清理例程結束
	done := make(chan struct{})
	go func() {
		m.cleanupWg.Wait()
		close(done)
	}()
	
	select {
	case <-done:
		// 正常結束
	case <-ctx.Done():
		// 超時
	}
	
	// 關閉所有資源
	var errors []error
	
	// 關閉 HTTP 客戶端
	m.httpManager.Shutdown()
	
	// 關閉所有文件
	if err := m.fileManager.CloseAllFiles(); err != nil {
		errors = append(errors, err)
	}
	
	// 關閉所有數據庫連接
	if err := m.dbManager.CloseAll(); err != nil {
		errors = append(errors, err)
	}
	
	if len(errors) > 0 {
		return errors[0]
	}
	
	return nil
}

// Stats 獲取資源統計信息
func (m *ResourceManager) Stats() map[string]interface{} {
	stats := make(map[string]interface{})
	
	stats["http_clients"] = m.httpManager.Stats()
	stats["files"] = m.fileManager.Stats()
	stats["databases"] = m.dbManager.Stats()
	
	return stats
}

// WithResources 在資源上下文中執行函數
func WithResources(ctx context.Context, fn func(context.Context, *ResourceManager) error) error {
	rm := NewResourceManager()
	
	// 啟動清理
	cleanupCtx, cancelCleanup := context.WithCancel(context.Background())
	rm.StartCleanup(cleanupCtx)
	
	// 確保資源被清理
	defer func() {
		cancelCleanup()
		
		// 給予 5 秒時間進行優雅關閉
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		rm.Shutdown(shutdownCtx)
	}()
	
	// 執行函數
	return fn(ctx, rm)
}

// GlobalResourceManager 全局資源管理器實例
var (
	globalRM   *ResourceManager
	globalOnce sync.Once
	globalMu   sync.RWMutex
)

// GetGlobalResourceManager 獲取全局資源管理器
func GetGlobalResourceManager() *ResourceManager {
	globalOnce.Do(func() {
		globalRM = NewResourceManager()
		
		// 啟動清理
		ctx := context.Background()
		globalRM.StartCleanup(ctx)
	})
	
	return globalRM
}

// ShutdownGlobalResourceManager 關閉全局資源管理器
func ShutdownGlobalResourceManager(ctx context.Context) error {
	globalMu.Lock()
	defer globalMu.Unlock()
	
	if globalRM != nil {
		return globalRM.Shutdown(ctx)
	}
	
	return nil
}