package swagger

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// Config Swagger配置
type Config struct {
	Title       string
	Description string
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
}

// Handler 創建Swagger文檔處理器
func Handler(cfg Config) echo.HandlerFunc {
	return echoSwagger.WrapHandler
}

// SecurityDefinitions 定義安全方案
type SecurityDefinitions struct {
	Bearer BearerAuth `json:"bearer,omitempty"`
}

// BearerAuth JWT Bearer認證
type BearerAuth struct {
	Type             string `json:"type"`
	Scheme           string `json:"scheme"`
	BearerFormat     string `json:"bearerFormat,omitempty"`
	Description      string `json:"description,omitempty"`
	Name             string `json:"name,omitempty"`
	In               string `json:"in,omitempty"`
}

// NewBearerAuth 創建Bearer認證配置
func NewBearerAuth() BearerAuth {
	return BearerAuth{
		Type:         "http",
		Scheme:       "bearer",
		BearerFormat: "JWT",
		Description:  "Enter the token with the `Bearer ` prefix, e.g. \"Bearer abcde12345\"",
		Name:         "Authorization",
		In:           "header",
	}
}

// Response 通用響應結構
type Response struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Operation successful"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse 錯誤響應結構
type ErrorResponse struct {
	Success bool                   `json:"success" example:"false"`
	Message string                 `json:"message" example:"Invalid request"`
	Error   string                 `json:"error,omitempty" example:"validation_error"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// PaginationResponse 分頁響應結構
type PaginationResponse struct {
	Success    bool        `json:"success" example:"true"`
	Message    string      `json:"message" example:"Data retrieved successfully"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

// Pagination 分頁信息
type Pagination struct {
	Page       int   `json:"page" example:"1"`
	PageSize   int   `json:"page_size" example:"20"`
	TotalPages int   `json:"total_pages" example:"5"`
	Total      int64 `json:"total" example:"100"`
}

// RegisterRoutes 註冊Swagger路由
func RegisterRoutes(e *echo.Echo, cfg Config) {
	// Swagger文檔路由
	e.GET("/swagger/*", Handler(cfg))
	
	// 重定向根路徑到Swagger UI
	e.GET("/swagger", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
}

// GenerateSpec 生成OpenAPI規範
func GenerateSpec(cfg Config) string {
	spec := fmt.Sprintf(`{
    "swagger": "2.0",
    "info": {
        "title": "%s",
        "description": "%s",
        "version": "%s",
        "contact": {
            "name": "API Support",
            "email": "support@fastenmind.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        }
    },
    "host": "%s",
    "basePath": "%s",
    "schemes": %v,
    "consumes": ["application/json"],
    "produces": ["application/json"],
    "securityDefinitions": {
        "Bearer": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header",
            "description": "Enter the token with the 'Bearer ' prefix, e.g. 'Bearer abcde12345'"
        }
    },
    "security": [
        {
            "Bearer": []
        }
    ]
}`, cfg.Title, cfg.Description, cfg.Version, cfg.Host, cfg.BasePath, cfg.Schemes)
	
	return spec
}

// Tag 定義API標籤
type Tag struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Tags 預定義的API標籤
var Tags = []Tag{
	{Name: "Auth", Description: "Authentication endpoints"},
	{Name: "Customers", Description: "Customer management"},
	{Name: "Suppliers", Description: "Supplier management"},
	{Name: "Products", Description: "Product management"},
	{Name: "Materials", Description: "Material management"},
	{Name: "Processes", Description: "Process management"},
	{Name: "Quotations", Description: "Quotation management"},
	{Name: "ExchangeRates", Description: "Exchange rate management"},
	{Name: "Reports", Description: "Report generation"},
	{Name: "Admin", Description: "Administrative functions"},
}

// Common response examples
const (
	SuccessExample = `{
    "success": true,
    "message": "Operation completed successfully",
    "data": {}
}`

	ErrorExample = `{
    "success": false,
    "message": "An error occurred",
    "error": "error_code",
    "details": {}
}`

	UnauthorizedExample = `{
    "success": false,
    "message": "Unauthorized",
    "error": "unauthorized"
}`

	ValidationErrorExample = `{
    "success": false,
    "message": "Validation failed",
    "error": "validation_error",
    "details": {
        "field_name": "error message"
    }
}`

	NotFoundExample = `{
    "success": false,
    "message": "Resource not found",
    "error": "not_found"
}`
)