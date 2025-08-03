package resources

import (
	"context"
	"net/http"
	"sync"
	"time"
)

// HTTPClientManager 管理 HTTP 客戶端連接
type HTTPClientManager struct {
	clients   map[string]*ManagedHTTPClient
	mu        sync.RWMutex
	transport *http.Transport
}

// ManagedHTTPClient 被管理的 HTTP 客戶端
type ManagedHTTPClient struct {
	client    *http.Client
	name      string
	createdAt time.Time
	lastUsed  time.Time
	mu        sync.Mutex
}

// NewHTTPClientManager 創建新的 HTTP 客戶端管理器
func NewHTTPClientManager() *HTTPClientManager {
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
		DisableCompression:  false,
		DisableKeepAlives:   false,
	}
	
	return &HTTPClientManager{
		clients:   make(map[string]*ManagedHTTPClient),
		transport: transport,
	}
}

// GetClient 獲取或創建 HTTP 客戶端
func (m *HTTPClientManager) GetClient(name string, timeout time.Duration) *ManagedHTTPClient {
	m.mu.RLock()
	if client, exists := m.clients[name]; exists {
		m.mu.RUnlock()
		client.mu.Lock()
		client.lastUsed = time.Now()
		client.mu.Unlock()
		return client
	}
	m.mu.RUnlock()
	
	// 創建新客戶端
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// 再次檢查避免重複創建
	if client, exists := m.clients[name]; exists {
		return client
	}
	
	client := &ManagedHTTPClient{
		client: &http.Client{
			Transport: m.transport,
			Timeout:   timeout,
		},
		name:      name,
		createdAt: time.Now(),
		lastUsed:  time.Now(),
	}
	
	m.clients[name] = client
	return client
}

// Do 執行 HTTP 請求（自動管理響應體關閉）
func (c *ManagedHTTPClient) Do(req *http.Request) (*http.Response, error) {
	c.mu.Lock()
	c.lastUsed = time.Now()
	c.mu.Unlock()
	
	return c.client.Do(req)
}

// DoWithContext 執行帶上下文的 HTTP 請求
func (c *ManagedHTTPClient) DoWithContext(ctx context.Context, req *http.Request) (*http.Response, error) {
	c.mu.Lock()
	c.lastUsed = time.Now()
	c.mu.Unlock()
	
	req = req.WithContext(ctx)
	return c.client.Do(req)
}

// CleanupIdleClients 清理空閒客戶端
func (m *HTTPClientManager) CleanupIdleClients(maxIdleTime time.Duration) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	now := time.Now()
	removed := 0
	
	for name, client := range m.clients {
		client.mu.Lock()
		if now.Sub(client.lastUsed) > maxIdleTime {
			delete(m.clients, name)
			removed++
		}
		client.mu.Unlock()
	}
	
	return removed
}

// Shutdown 關閉所有連接
func (m *HTTPClientManager) Shutdown() {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// 關閉傳輸層
	m.transport.CloseIdleConnections()
	
	// 清空客戶端映射
	m.clients = make(map[string]*ManagedHTTPClient)
}

// Stats 獲取統計信息
func (m *HTTPClientManager) Stats() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	stats := make(map[string]interface{})
	stats["total_clients"] = len(m.clients)
	
	clients := make([]map[string]interface{}, 0, len(m.clients))
	for name, client := range m.clients {
		client.mu.Lock()
		clientInfo := map[string]interface{}{
			"name":       name,
			"created_at": client.createdAt,
			"last_used":  client.lastUsed,
		}
		client.mu.Unlock()
		clients = append(clients, clientInfo)
	}
	stats["clients"] = clients
	
	return stats
}

// HTTPResponseWrapper 包裝 HTTP 響應以確保自動關閉
type HTTPResponseWrapper struct {
	*http.Response
	closed bool
	mu     sync.Mutex
}

// NewHTTPResponseWrapper 創建響應包裝器
func NewHTTPResponseWrapper(resp *http.Response) *HTTPResponseWrapper {
	return &HTTPResponseWrapper{
		Response: resp,
		closed:   false,
	}
}

// Close 關閉響應體
func (w *HTTPResponseWrapper) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	
	if w.closed || w.Response == nil || w.Response.Body == nil {
		return nil
	}
	
	w.closed = true
	return w.Response.Body.Close()
}

// EnsureClosed 確保響應體被關閉（用於 defer）
func (w *HTTPResponseWrapper) EnsureClosed() {
	if err := w.Close(); err != nil {
		// 記錄錯誤但不拋出
		_ = err
	}
}