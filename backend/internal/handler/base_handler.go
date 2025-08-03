package handler

import (
	"context"
	"database/sql"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/fastenmind/fastener-api/pkg/resources"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// BaseHandler 提供資源管理的基礎處理器
type BaseHandler struct {
	rm *resources.ResourceManager
}

// NewBaseHandler 創建基礎處理器
func NewBaseHandler() *BaseHandler {
	return &BaseHandler{
		rm: resources.GetGlobalResourceManager(),
	}
}

// WithTransaction 在事務中執行操作
func (h *BaseHandler) WithTransaction(c echo.Context, fn func(tx *gorm.DB) error) error {
	ctx := c.Request().Context()
	return h.rm.Databases().WithTransaction(ctx, "main", fn)
}

// WithHTTPClient 使用管理的 HTTP 客戶端
func (h *BaseHandler) WithHTTPClient(c echo.Context, name string, timeout time.Duration, fn func(client *resources.ManagedHTTPClient) error) error {
	client := h.rm.HTTPClients().GetClient(name, timeout)
	return fn(client)
}

// DoHTTPRequest 執行 HTTP 請求並自動清理響應
func (h *BaseHandler) DoHTTPRequest(c echo.Context, req *http.Request) (*http.Response, error) {
	client := h.rm.HTTPClients().GetClient("default", 30*time.Second)
	
	ctx := c.Request().Context()
	resp, err := client.DoWithContext(ctx, req)
	if err != nil {
		return nil, err
	}
	
	// 註冊響應清理
	c.Set("defer_cleanup_response", resp)
	
	return resp, nil
}

// ReadFile 讀取文件並自動清理
func (h *BaseHandler) ReadFile(c echo.Context, path string, handler func(io.Reader) error) error {
	ctx := c.Request().Context()
	return resources.ReadFileWithCleanup(ctx, path, handler)
}

// WriteFile 寫入文件並自動清理
func (h *BaseHandler) WriteFile(c echo.Context, path string, handler func(io.Writer) error) error {
	ctx := c.Request().Context()
	return resources.WriteFileWithCleanup(ctx, path, handler)
}

// WithTempFile 使用臨時文件並自動清理
func (h *BaseHandler) WithTempFile(c echo.Context, pattern string, fn func(file *os.File) error) error {
	mf, err := h.rm.Files().CreateTempFile("", pattern)
	if err != nil {
		return err
	}
	
	// 獲取文件路徑
	filePath := mf.Path()
	
	// 確保文件被關閉和刪除
	defer func() {
		h.rm.Files().CloseFile(filePath)
	}()
	
	// 獲取底層文件對象
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	
	return fn(file)
}

// QueryWithContext 執行數據庫查詢並管理資源
func (h *BaseHandler) QueryWithContext(c echo.Context, dest interface{}, query string, args ...interface{}) error {
	db, err := h.rm.Databases().GetDB("main")
	if err != nil {
		return err
	}
	
	ctx := c.Request().Context()
	return resources.QueryContext(ctx, db, dest, query, args...)
}

// ExecWithContext 執行數據庫語句並管理資源
func (h *BaseHandler) ExecWithContext(c echo.Context, query string, args ...interface{}) error {
	db, err := h.rm.Databases().GetDB("main")
	if err != nil {
		return err
	}
	
	ctx := c.Request().Context()
	return resources.ExecContext(ctx, db, query, args...)
}

// CleanupHandler 清理處理器中間件
func CleanupHandler() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 創建清理函數列表
			cleanups := make([]func(), 0)
			c.Set("cleanups", &cleanups)
			
			// 在處理完成後執行所有清理函數
			defer func() {
				if cls, ok := c.Get("cleanups").(*[]func()); ok {
					for _, cleanup := range *cls {
						cleanup()
					}
				}
				
				// 清理 HTTP 響應
				if resp, ok := c.Get("defer_cleanup_response").(*http.Response); ok && resp != nil && resp.Body != nil {
					resp.Body.Close()
				}
			}()
			
			return next(c)
		}
	}
}

// RegisterCleanup 註冊清理函數
func RegisterCleanup(c echo.Context, cleanup func()) {
	if cleanups, ok := c.Get("cleanups").(*[]func()); ok {
		*cleanups = append(*cleanups, cleanup)
	}
}

// WithContext 創建帶超時的上下文
func WithContext(c echo.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	
	return context.WithTimeout(c.Request().Context(), timeout)
}

// HandlePanic 處理 panic 並確保資源清理
func HandlePanic(c echo.Context) {
	if r := recover(); r != nil {
		// 記錄 panic
		c.Logger().Error("Panic recovered:", r)
		
		// 確保資源被清理
		if cleanups, ok := c.Get("cleanups").(*[]func()); ok {
			for _, cleanup := range *cleanups {
				func() {
					defer func() {
						if r := recover(); r != nil {
							c.Logger().Error("Panic in cleanup:", r)
						}
					}()
					cleanup()
				}()
			}
		}
		
		// 重新拋出 panic
		panic(r)
	}
}

// PreparedStatements 預編譯語句管理
type PreparedStatements struct {
	stmts map[string]*sql.Stmt
	db    *sql.DB
}

// NewPreparedStatements 創建預編譯語句管理器
func NewPreparedStatements(db *sql.DB) *PreparedStatements {
	return &PreparedStatements{
		stmts: make(map[string]*sql.Stmt),
		db:    db,
	}
}

// Prepare 準備語句
func (ps *PreparedStatements) Prepare(name, query string) error {
	stmt, err := ps.db.Prepare(query)
	if err != nil {
		return err
	}
	
	// 如果已存在，先關閉舊的
	if old, exists := ps.stmts[name]; exists {
		old.Close()
	}
	
	ps.stmts[name] = stmt
	return nil
}

// Get 獲取預編譯語句
func (ps *PreparedStatements) Get(name string) (*sql.Stmt, bool) {
	stmt, exists := ps.stmts[name]
	return stmt, exists
}

// Close 關閉所有語句
func (ps *PreparedStatements) Close() {
	for _, stmt := range ps.stmts {
		stmt.Close()
	}
	ps.stmts = make(map[string]*sql.Stmt)
}