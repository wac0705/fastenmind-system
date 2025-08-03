package middleware

import (
	"github.com/fastenmind/fastener-api/pkg/errors"
	"github.com/labstack/echo/v4"
)

// ErrorMiddleware 錯誤處理中間件
func ErrorMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 執行下一個處理器
			err := next(c)
			
			// 如果有錯誤，使用統一的錯誤處理器
			if err != nil {
				errors.ErrorHandler(err, c)
			}
			
			return nil
		}
	}
}

// SetupErrorHandling 設置錯誤處理
func SetupErrorHandling(e *echo.Echo) {
	// 設置自定義錯誤處理器
	e.HTTPErrorHandler = errors.ErrorHandler
	
	// 添加恢復中間件
	e.Use(errors.RecoverMiddleware())
	
	// 添加錯誤處理中間件
	e.Use(ErrorMiddleware())
}