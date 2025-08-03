package monitoring

import (
	"context"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
)

// Dashboard 監控儀表板
type Dashboard struct {
	metrics      *Metrics
	healthChecker *HealthChecker
	statsCache   *StatsCache
	mu           sync.RWMutex
}

// NewDashboard 創建監控儀表板
func NewDashboard(metrics *Metrics, healthChecker *HealthChecker) *Dashboard {
	return &Dashboard{
		metrics:       metrics,
		healthChecker: healthChecker,
		statsCache:    NewStatsCache(30 * time.Second),
	}
}

// StatsCache 統計快取
type StatsCache struct {
	data      map[string]interface{}
	expiresAt time.Time
	ttl       time.Duration
	mu        sync.RWMutex
}

// NewStatsCache 創建統計快取
func NewStatsCache(ttl time.Duration) *StatsCache {
	return &StatsCache{
		data: make(map[string]interface{}),
		ttl:  ttl,
	}
}

// Get 獲取快取數據
func (c *StatsCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if time.Now().After(c.expiresAt) {
		return nil, false
	}

	value, exists := c.data[key]
	return value, exists
}

// Set 設置快取數據
func (c *StatsCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value
	c.expiresAt = time.Now().Add(c.ttl)
}

// Clear 清除快取
func (c *StatsCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = make(map[string]interface{})
	c.expiresAt = time.Time{}
}

// DashboardStats 儀表板統計數據
type DashboardStats struct {
	System    SystemStats    `json:"system"`
	Business  BusinessStats  `json:"business"`
	Database  DatabaseStats  `json:"database"`
	Cache     CacheStats     `json:"cache"`
	Health    HealthStats    `json:"health"`
	Timestamp time.Time      `json:"timestamp"`
}

// SystemStats 系統統計
type SystemStats struct {
	CPUUsage      float64 `json:"cpu_usage"`
	MemoryUsage   float64 `json:"memory_usage"`
	MemoryTotal   uint64  `json:"memory_total"`
	MemoryUsed    uint64  `json:"memory_used"`
	Goroutines    int     `json:"goroutines"`
	OpenFiles     int     `json:"open_files"`
	Uptime        string  `json:"uptime"`
	Version       string  `json:"version"`
}

// BusinessStats 業務統計
type BusinessStats struct {
	TotalRequests      int64   `json:"total_requests"`
	SuccessRate        float64 `json:"success_rate"`
	AverageResponseTime float64 `json:"average_response_time"`
	ActiveUsers        int     `json:"active_users"`
	TodayQuotations    int     `json:"today_quotations"`
	TodayRevenue       float64 `json:"today_revenue"`
}

// DatabaseStats 資料庫統計
type DatabaseStats struct {
	ActiveConnections int     `json:"active_connections"`
	IdleConnections   int     `json:"idle_connections"`
	SlowQueries       int     `json:"slow_queries"`
	ErrorRate         float64 `json:"error_rate"`
}

// CacheStats 快取統計
type CacheStats struct {
	HitRate    float64 `json:"hit_rate"`
	Size       int     `json:"size"`
	Evictions  int64   `json:"evictions"`
	MemoryUsed int64   `json:"memory_used"`
}

// HealthStats 健康統計
type HealthStats struct {
	Overall    string                    `json:"overall"`
	Components map[string]HealthStatus   `json:"components"`
}

// GetStats 獲取統計數據
func (d *Dashboard) GetStats() (*DashboardStats, error) {
	// 檢查快取
	if cached, found := d.statsCache.Get("dashboard_stats"); found {
		return cached.(*DashboardStats), nil
	}

	stats := &DashboardStats{
		Timestamp: time.Now(),
	}

	// 獲取系統統計
	systemStats, err := d.getSystemStats()
	if err != nil {
		return nil, err
	}
	stats.System = *systemStats

	// 獲取業務統計
	stats.Business = d.getBusinessStats()

	// 獲取資料庫統計
	stats.Database = d.getDatabaseStats()

	// 獲取快取統計
	stats.Cache = d.getCacheStats()

	// 獲取健康狀態
	stats.Health = d.getHealthStats()

	// 快取結果
	d.statsCache.Set("dashboard_stats", stats)

	return stats, nil
}

// getSystemStats 獲取系統統計
func (d *Dashboard) getSystemStats() (*SystemStats, error) {
	stats := &SystemStats{
		Goroutines: runtime.NumGoroutine(),
		Version:    runtime.Version(),
	}

	// CPU 使用率
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err == nil && len(cpuPercent) > 0 {
		stats.CPUUsage = cpuPercent[0]
	}

	// 內存使用
	memInfo, err := mem.VirtualMemory()
	if err == nil {
		stats.MemoryTotal = memInfo.Total
		stats.MemoryUsed = memInfo.Used
		stats.MemoryUsage = memInfo.UsedPercent
	}

	// 打開文件數
	proc, err := process.NewProcess(int32(os.Getpid()))
	if err == nil {
		files, _ := proc.OpenFiles()
		stats.OpenFiles = len(files)
	}

	// 運行時間
	stats.Uptime = time.Since(startTime).String()

	// 更新指標
	d.metrics.UpdateSystemMetrics(
		float64(stats.MemoryUsed),
		stats.CPUUsage,
		stats.Goroutines,
		stats.OpenFiles,
	)

	return stats, nil
}

// getBusinessStats 獲取業務統計
func (d *Dashboard) getBusinessStats() BusinessStats {
	// 這裡應該從實際的業務數據源獲取
	// 這是示例實現
	return BusinessStats{
		TotalRequests:       1000,
		SuccessRate:         0.95,
		AverageResponseTime: 150.5, // ms
		ActiveUsers:         50,
		TodayQuotations:     25,
		TodayRevenue:        150000.00,
	}
}

// getDatabaseStats 獲取資料庫統計
func (d *Dashboard) getDatabaseStats() DatabaseStats {
	// 這裡應該從實際的資料庫連接池獲取
	// 這是示例實現
	return DatabaseStats{
		ActiveConnections: 10,
		IdleConnections:   40,
		SlowQueries:       2,
		ErrorRate:         0.01,
	}
}

// getCacheStats 獲取快取統計
func (d *Dashboard) getCacheStats() CacheStats {
	// 這裡應該從實際的快取系統獲取
	// 這是示例實現
	return CacheStats{
		HitRate:    0.85,
		Size:       1000,
		Evictions:  100,
		MemoryUsed: 1024 * 1024 * 50, // 50MB
	}
}

// getHealthStats 獲取健康狀態
func (d *Dashboard) getHealthStats() HealthStats {
	ctx := context.Background()
	components := d.healthChecker.CheckHealth(ctx)
	
	overall := "healthy"
	for _, status := range components {
		if status.Status != "healthy" {
			overall = "unhealthy"
			break
		}
	}

	return HealthStats{
		Overall:    overall,
		Components: components,
	}
}

// RegisterRoutes 註冊儀表板路由
func (d *Dashboard) RegisterRoutes(e *echo.Echo) {
	// 主儀表板頁面
	e.GET("/dashboard", d.handleDashboard)
	
	// API 端點
	api := e.Group("/api/monitoring")
	api.GET("/stats", d.handleGetStats)
	api.GET("/health", d.handleHealthCheck)
	api.GET("/metrics", PrometheusHandler())
}

// handleDashboard 處理儀表板頁面請求
func (d *Dashboard) handleDashboard(c echo.Context) error {
	// 返回儀表板 HTML 頁面
	return c.HTML(http.StatusOK, dashboardHTML)
}

// handleGetStats 處理獲取統計數據請求
func (d *Dashboard) handleGetStats(c echo.Context) error {
	stats, err := d.GetStats()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, stats)
}

// handleHealthCheck 處理健康檢查請求
func (d *Dashboard) handleHealthCheck(c echo.Context) error {
	ctx := c.Request().Context()
	results := d.healthChecker.CheckHealth(ctx)
	
	status := http.StatusOK
	for _, result := range results {
		if result.Status != "healthy" {
			status = http.StatusServiceUnavailable
			break
		}
	}

	return c.JSON(status, results)
}

// WebSocketHub WebSocket 連接管理器
type WebSocketHub struct {
	clients    map[*WebSocketClient]bool
	broadcast  chan []byte
	register   chan *WebSocketClient
	unregister chan *WebSocketClient
	mu         sync.RWMutex
}

// WebSocketClient WebSocket 客戶端
type WebSocketClient struct {
	hub  *WebSocketHub
	conn *websocket.Conn
	send chan []byte
}

// NewWebSocketHub 創建 WebSocket Hub
func NewWebSocketHub() *WebSocketHub {
	return &WebSocketHub{
		broadcast:  make(chan []byte),
		register:   make(chan *WebSocketClient),
		unregister: make(chan *WebSocketClient),
		clients:    make(map[*WebSocketClient]bool),
	}
}

// Run 運行 WebSocket Hub
func (h *WebSocketHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// 程序啟動時間
var startTime = time.Now()

// 儀表板 HTML 模板
const dashboardHTML = `
<!DOCTYPE html>
<html>
<head>
    <title>FastenMind Monitoring Dashboard</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 20px; background: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; }
        .header { background: #333; color: white; padding: 20px; margin-bottom: 20px; }
        .grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 20px; }
        .card { background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .card h3 { margin-top: 0; color: #333; }
        .metric { display: flex; justify-content: space-between; margin: 10px 0; }
        .metric-value { font-weight: bold; color: #2196F3; }
        .status-healthy { color: #4CAF50; }
        .status-unhealthy { color: #F44336; }
        .chart { height: 200px; margin-top: 20px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>FastenMind Monitoring Dashboard</h1>
            <p>Real-time system monitoring and metrics</p>
        </div>
        
        <div class="grid">
            <div class="card">
                <h3>System Health</h3>
                <div id="health-status"></div>
            </div>
            
            <div class="card">
                <h3>System Metrics</h3>
                <div id="system-metrics"></div>
            </div>
            
            <div class="card">
                <h3>Business Metrics</h3>
                <div id="business-metrics"></div>
            </div>
            
            <div class="card">
                <h3>Database Status</h3>
                <div id="database-status"></div>
            </div>
            
            <div class="card">
                <h3>Cache Performance</h3>
                <div id="cache-performance"></div>
            </div>
            
            <div class="card">
                <h3>Response Time</h3>
                <canvas id="response-time-chart" class="chart"></canvas>
            </div>
        </div>
    </div>
    
    <script>
        // Fetch and update dashboard data
        async function updateDashboard() {
            try {
                const response = await fetch('/api/monitoring/stats');
                const data = await response.json();
                
                updateHealthStatus(data.health);
                updateSystemMetrics(data.system);
                updateBusinessMetrics(data.business);
                updateDatabaseStatus(data.database);
                updateCachePerformance(data.cache);
            } catch (error) {
                console.error('Error updating dashboard:', error);
            }
        }
        
        function updateHealthStatus(health) {
            const container = document.getElementById('health-status');
            container.innerHTML = '<div class="metric"><span>Overall Status</span><span class="metric-value status-' + 
                health.overall + '">' + health.overall.toUpperCase() + '</span></div>';
        }
        
        function updateSystemMetrics(system) {
            const container = document.getElementById('system-metrics');
            container.innerHTML = 
                '<div class="metric"><span>CPU Usage</span><span class="metric-value">' + system.cpu_usage.toFixed(1) + '%</span></div>' +
                '<div class="metric"><span>Memory Usage</span><span class="metric-value">' + system.memory_usage.toFixed(1) + '%</span></div>' +
                '<div class="metric"><span>Goroutines</span><span class="metric-value">' + system.goroutines + '</span></div>' +
                '<div class="metric"><span>Uptime</span><span class="metric-value">' + system.uptime + '</span></div>';
        }
        
        function updateBusinessMetrics(business) {
            const container = document.getElementById('business-metrics');
            container.innerHTML = 
                '<div class="metric"><span>Active Users</span><span class="metric-value">' + business.active_users + '</span></div>' +
                '<div class="metric"><span>Success Rate</span><span class="metric-value">' + (business.success_rate * 100).toFixed(1) + '%</span></div>' +
                '<div class="metric"><span>Today Quotations</span><span class="metric-value">' + business.today_quotations + '</span></div>' +
                '<div class="metric"><span>Today Revenue</span><span class="metric-value">$' + business.today_revenue.toLocaleString() + '</span></div>';
        }
        
        function updateDatabaseStatus(database) {
            const container = document.getElementById('database-status');
            container.innerHTML = 
                '<div class="metric"><span>Active Connections</span><span class="metric-value">' + database.active_connections + '</span></div>' +
                '<div class="metric"><span>Idle Connections</span><span class="metric-value">' + database.idle_connections + '</span></div>' +
                '<div class="metric"><span>Slow Queries</span><span class="metric-value">' + database.slow_queries + '</span></div>' +
                '<div class="metric"><span>Error Rate</span><span class="metric-value">' + (database.error_rate * 100).toFixed(2) + '%</span></div>';
        }
        
        function updateCachePerformance(cache) {
            const container = document.getElementById('cache-performance');
            container.innerHTML = 
                '<div class="metric"><span>Hit Rate</span><span class="metric-value">' + (cache.hit_rate * 100).toFixed(1) + '%</span></div>' +
                '<div class="metric"><span>Cache Size</span><span class="metric-value">' + cache.size + '</span></div>' +
                '<div class="metric"><span>Evictions</span><span class="metric-value">' + cache.evictions + '</span></div>' +
                '<div class="metric"><span>Memory Used</span><span class="metric-value">' + (cache.memory_used / 1024 / 1024).toFixed(1) + ' MB</span></div>';
        }
        
        // Update dashboard every 5 seconds
        updateDashboard();
        setInterval(updateDashboard, 5000);
    </script>
</body>
</html>
`

