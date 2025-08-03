package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/fastenmind/fastener-api/pkg/validator"
	"github.com/labstack/echo/v4"
)

// ValidationMiddleware 輸入驗證中間件
func ValidationMiddleware() echo.MiddlewareFunc {
	v := validator.New()
	
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 只對有請求體的方法進行驗證
			if c.Request().Method == http.MethodPost || 
			   c.Request().Method == http.MethodPut || 
			   c.Request().Method == http.MethodPatch {
				
				// 檢查 Content-Type
				contentType := c.Request().Header.Get("Content-Type")
				if !strings.Contains(contentType, "application/json") &&
				   !strings.Contains(contentType, "multipart/form-data") {
					return echo.NewHTTPError(http.StatusBadRequest, "Invalid content type")
				}
			}
			
			// 設置驗證器到上下文
			c.Set("validator", v)
			
			return next(c)
		}
	}
}

// ValidateRequest 驗證請求結構體
func ValidateRequest(c echo.Context, req interface{}) error {
	// 綁定請求數據
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}
	
	// 獲取驗證器
	v, ok := c.Get("validator").(*validator.Validator)
	if !ok {
		v = validator.New()
	}
	
	// 執行驗證
	if err := v.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	return nil
}

// ValidateQueryParams 驗證查詢參數
func ValidateQueryParams(c echo.Context, params interface{}) error {
	// 綁定查詢參數
	if err := c.Bind(params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid query parameters")
	}
	
	// 獲取驗證器
	v, ok := c.Get("validator").(*validator.Validator)
	if !ok {
		v = validator.New()
	}
	
	// 執行驗證
	if err := v.Validate(params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	return nil
}

// SanitizeString 清理字符串輸入
func SanitizeString(input string) string {
	// 移除前後空白
	input = strings.TrimSpace(input)
	
	// 移除控制字符
	var builder strings.Builder
	for _, r := range input {
		if r >= 32 && r != 127 { // 可打印字符
			builder.WriteRune(r)
		}
	}
	
	return builder.String()
}

// ValidateID 驗證 ID 參數
func ValidateID(c echo.Context, paramName string) (string, error) {
	id := c.Param(paramName)
	if id == "" {
		return "", echo.NewHTTPError(http.StatusBadRequest, "Missing "+paramName)
	}
	
	// 驗證 UUID 格式
	v := validator.New()
	if err := v.ValidateVar(id, "uuid"); err != nil {
		return "", echo.NewHTTPError(http.StatusBadRequest, "Invalid "+paramName+" format")
	}
	
	return id, nil
}

// ValidatePagination 驗證分頁參數
func ValidatePagination(c echo.Context) (page, pageSize int, err error) {
	// 預設值
	page = 1
	pageSize = 20
	
	// 解析頁碼
	if p := c.QueryParam("page"); p != "" {
		var pageInt int
		if _, err := fmt.Sscanf(p, "%d", &pageInt); err == nil && pageInt > 0 {
			page = pageInt
		}
	}
	
	// 解析頁面大小
	if ps := c.QueryParam("page_size"); ps != "" {
		var sizeInt int
		if _, err := fmt.Sscanf(ps, "%d", &sizeInt); err == nil && sizeInt > 0 {
			pageSize = sizeInt
			// 限制最大頁面大小
			if pageSize > 100 {
				pageSize = 100
			}
		}
	}
	
	return page, pageSize, nil
}

// ValidateSortParams 驗證排序參數
func ValidateSortParams(c echo.Context, allowedFields []string) (sortBy string, sortOrder string, err error) {
	sortBy = c.QueryParam("sort_by")
	sortOrder = strings.ToUpper(c.QueryParam("sort_order"))
	
	// 預設值
	if sortBy == "" {
		sortBy = "created_at"
	}
	if sortOrder == "" {
		sortOrder = "DESC"
	}
	
	// 驗證排序欄位
	isValidField := false
	for _, field := range allowedFields {
		if sortBy == field {
			isValidField = true
			break
		}
	}
	
	if !isValidField {
		return "", "", echo.NewHTTPError(http.StatusBadRequest, "Invalid sort field")
	}
	
	// 驗證排序方向
	if sortOrder != "ASC" && sortOrder != "DESC" {
		return "", "", echo.NewHTTPError(http.StatusBadRequest, "Invalid sort order")
	}
	
	return sortBy, sortOrder, nil
}

// ValidateDateRange 驗證日期範圍
func ValidateDateRange(c echo.Context) (startDate, endDate *time.Time, err error) {
	startStr := c.QueryParam("start_date")
	endStr := c.QueryParam("end_date")
	
	if startStr != "" {
		start, err := time.Parse("2006-01-02", startStr)
		if err != nil {
			return nil, nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid start date format")
		}
		startDate = &start
	}
	
	if endStr != "" {
		end, err := time.Parse("2006-01-02", endStr)
		if err != nil {
			return nil, nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid end date format")
		}
		endDate = &end
	}
	
	// 驗證日期範圍邏輯
	if startDate != nil && endDate != nil && startDate.After(*endDate) {
		return nil, nil, echo.NewHTTPError(http.StatusBadRequest, "Start date must be before end date")
	}
	
	return startDate, endDate, nil
}