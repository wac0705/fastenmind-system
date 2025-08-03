package middleware

import (
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

// ResponseCleanup 確保響應體被正確關閉的中間件
func ResponseCleanup() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 處理請求
			err := next(c)
			
			// 確保請求體被關閉
			if c.Request().Body != nil {
				defer c.Request().Body.Close()
			}
			
			return err
		}
	}
}

// HTTPClientCleanup 為 HTTP 客戶端響應添加自動清理
type CleanupHTTPClient struct {
	client *http.Client
}

// NewCleanupHTTPClient 創建帶清理功能的 HTTP 客戶端
func NewCleanupHTTPClient(client *http.Client) *CleanupHTTPClient {
	if client == nil {
		client = http.DefaultClient
	}
	return &CleanupHTTPClient{client: client}
}

// Do 執行請求並返回自動清理的響應
func (c *CleanupHTTPClient) Do(req *http.Request) (*CleanupResponse, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	
	return &CleanupResponse{Response: resp}, nil
}

// CleanupResponse 包裝 HTTP 響應以確保自動清理
type CleanupResponse struct {
	*http.Response
	closed bool
}

// Body 返回響應體（確保只讀取一次）
func (r *CleanupResponse) Body() io.ReadCloser {
	return &cleanupReadCloser{
		ReadCloser: r.Response.Body,
		onClose: func() {
			r.closed = true
		},
	}
}

// EnsureClosed 確保響應體被關閉
func (r *CleanupResponse) EnsureClosed() {
	if !r.closed && r.Response != nil && r.Response.Body != nil {
		r.Response.Body.Close()
		r.closed = true
	}
}

// cleanupReadCloser 包裝 io.ReadCloser 以跟踪關閉狀態
type cleanupReadCloser struct {
	io.ReadCloser
	onClose func()
}

func (c *cleanupReadCloser) Close() error {
	err := c.ReadCloser.Close()
	if c.onClose != nil {
		c.onClose()
	}
	return err
}

// WithHTTPCleanup 為 Echo 上下文添加 HTTP 客戶端清理功能
func WithHTTPCleanup(client *http.Client) echo.MiddlewareFunc {
	cleanupClient := NewCleanupHTTPClient(client)
	
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 將清理客戶端添加到上下文
			c.Set("cleanup_http_client", cleanupClient)
			
			// 創建響應清理列表
			responses := make([]*CleanupResponse, 0)
			c.Set("cleanup_responses", &responses)
			
			// 在請求結束時清理所有響應
			defer func() {
				if resps, ok := c.Get("cleanup_responses").(*[]*CleanupResponse); ok {
					for _, resp := range *resps {
						resp.EnsureClosed()
					}
				}
			}()
			
			return next(c)
		}
	}
}

// GetCleanupHTTPClient 從上下文獲取清理 HTTP 客戶端
func GetCleanupHTTPClient(c echo.Context) *CleanupHTTPClient {
	if client, ok := c.Get("cleanup_http_client").(*CleanupHTTPClient); ok {
		return client
	}
	return NewCleanupHTTPClient(nil)
}

// RegisterCleanupResponse 註冊需要清理的響應
func RegisterCleanupResponse(c echo.Context, resp *CleanupResponse) {
	if responses, ok := c.Get("cleanup_responses").(*[]*CleanupResponse); ok {
		*responses = append(*responses, resp)
	}
}