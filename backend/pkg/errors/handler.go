package errors

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// ErrorHandler 全局錯誤處理器
func ErrorHandler(err error, c echo.Context) {
	// 如果響應已經發送，不再處理
	if c.Response().Committed {
		return
	}

	// 處理 Echo 錯誤
	if he, ok := err.(*echo.HTTPError); ok {
		err = handleEchoError(he)
	}

	// 處理應用程式錯誤
	var appErr *AppError
	if IsAppError(err) {
		appErr = GetAppError(err)
	} else {
		// 轉換其他錯誤類型
		appErr = convertError(err)
	}

	// 記錄錯誤
	logError(c, appErr)

	// 發送錯誤響應
	c.JSON(appErr.StatusCode, appErr.ToResponse())
}

// handleEchoError 處理 Echo 框架錯誤
func handleEchoError(he *echo.HTTPError) *AppError {
	message := "Internal server error"
	if msg, ok := he.Message.(string); ok {
		message = msg
	}

	switch he.Code {
	case http.StatusBadRequest:
		return NewBadRequestError(message)
	case http.StatusUnauthorized:
		return NewUnauthorizedError(message)
	case http.StatusForbidden:
		return NewForbiddenError(message)
	case http.StatusNotFound:
		return NewNotFoundError(message)
	case http.StatusConflict:
		return NewConflictError(message)
	case http.StatusTooManyRequests:
		return NewRateLimitError(message)
	default:
		return NewInternalError(message)
	}
}

// convertError 轉換其他錯誤類型為應用程式錯誤
func convertError(err error) *AppError {
	// 處理 GORM 錯誤
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return NewNotFoundError("Record not found")
	}
	
	if errors.Is(err, gorm.ErrInvalidTransaction) {
		return NewDatabaseError("Invalid transaction", err)
	}
	
	if errors.Is(err, gorm.ErrNotImplemented) {
		return NewInternalError("Feature not implemented")
	}
	
	if errors.Is(err, gorm.ErrMissingWhereClause) {
		return NewBadRequestError("Missing where clause")
	}
	
	if errors.Is(err, gorm.ErrUnsupportedRelation) {
		return NewInternalError("Unsupported relation")
	}
	
	if errors.Is(err, gorm.ErrPrimaryKeyRequired) {
		return NewBadRequestError("Primary key required")
	}
	
	if errors.Is(err, gorm.ErrModelValueRequired) {
		return NewBadRequestError("Model value required")
	}
	
	if errors.Is(err, gorm.ErrInvalidData) {
		return NewBadRequestError("Invalid data")
	}
	
	if errors.Is(err, gorm.ErrUnsupportedDriver) {
		return NewInternalError("Unsupported database driver")
	}
	
	if errors.Is(err, gorm.ErrRegistered) {
		return NewInternalError("Model already registered")
	}
	
	if errors.Is(err, gorm.ErrInvalidField) {
		return NewBadRequestError("Invalid field")
	}
	
	if errors.Is(err, gorm.ErrEmptySlice) {
		return NewBadRequestError("Empty slice not allowed")
	}
	
	if errors.Is(err, gorm.ErrDryRunModeUnsupported) {
		return NewInternalError("Dry run mode unsupported")
	}

	// 處理標準 SQL 錯誤
	if errors.Is(err, sql.ErrNoRows) {
		return NewNotFoundError("No data found")
	}
	
	if errors.Is(err, sql.ErrConnDone) {
		return NewDatabaseError("Database connection closed", err)
	}
	
	if errors.Is(err, sql.ErrTxDone) {
		return NewDatabaseError("Transaction already completed", err)
	}

	// 處理字符串錯誤消息
	errMsg := err.Error()
	lowerMsg := strings.ToLower(errMsg)
	
	// 資料庫相關錯誤
	if strings.Contains(lowerMsg, "duplicate") || strings.Contains(lowerMsg, "unique") {
		return NewConflictError("Duplicate entry")
	}
	
	if strings.Contains(lowerMsg, "foreign key") {
		return NewBadRequestError("Foreign key constraint violation")
	}
	
	if strings.Contains(lowerMsg, "connection") || strings.Contains(lowerMsg, "dial") {
		return NewDatabaseError("Database connection error", err)
	}
	
	if strings.Contains(lowerMsg, "timeout") {
		return NewTimeoutError("Operation timed out")
	}
	
	// 驗證相關錯誤
	if strings.Contains(lowerMsg, "validation") || strings.Contains(lowerMsg, "invalid") {
		return NewValidationError(errMsg)
	}
	
	// 權限相關錯誤
	if strings.Contains(lowerMsg, "unauthorized") || strings.Contains(lowerMsg, "authentication") {
		return NewUnauthorizedError(errMsg)
	}
	
	if strings.Contains(lowerMsg, "forbidden") || strings.Contains(lowerMsg, "permission") {
		return NewForbiddenError(errMsg)
	}
	
	// 默認為內部錯誤
	return NewInternalError("An unexpected error occurred")
}

// logError 記錄錯誤
func logError(c echo.Context, err *AppError) {
	logger := c.Logger()
	
	// 構建錯誤上下文
	errorContext := map[string]interface{}{
		"type":       err.Type,
		"message":    err.Message,
		"code":       err.Code,
		"path":       c.Request().URL.Path,
		"method":     c.Request().Method,
		"request_id": c.Response().Header().Get(echo.HeaderXRequestID),
	}
	
	// 添加用戶信息（如果有）
	if userID := c.Get("user_id"); userID != nil {
		errorContext["user_id"] = userID
	}
	
	if companyID := c.Get("company_id"); companyID != nil {
		errorContext["company_id"] = companyID
	}
	
	// 根據錯誤類型決定日誌級別
	switch err.Type {
	case ErrorTypeInternal, ErrorTypeDatabase, ErrorTypeNetwork:
		// 記錄堆疊信息
		if len(err.Stack) > 0 {
			errorContext["stack"] = err.Stack
		}
		if err.Cause != nil {
			errorContext["cause"] = err.Cause.Error()
		}
		logger.Error(errorContext)
	case ErrorTypeValidation, ErrorTypeBadRequest, ErrorTypeNotFound:
		logger.Debug(errorContext)
	default:
		logger.Warn(errorContext)
	}
}

// RecoverMiddleware 恢復中間件
func RecoverMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					var panicErr *AppError
					if err, ok := r.(error); ok {
						panicErr = Wrap(err, ErrorTypeInternal, "A system error occurred")
					} else {
						panicErr = NewInternalError("A system error occurred")
					}
					
					// 添加 panic 詳情
					panicErr = panicErr.
						WithCode("PANIC").
						WithDetails(map[string]interface{}{
							"panic": r,
						})
					
					// 使用錯誤處理器
					ErrorHandler(panicErr, c)
				}
			}()
			
			return next(c)
		}
	}
}

// ValidationErrorHandler 處理驗證錯誤
func ValidationErrorHandler(c echo.Context, err error) error {
	// 解析驗證錯誤並創建詳細的錯誤響應
	validationErrors := parseValidationErrors(err)
	
	if len(validationErrors) > 0 {
		return NewValidationErrorWithDetails("Validation failed", validationErrors)
	}
	
	return NewValidationError(err.Error())
}

// parseValidationErrors 解析驗證錯誤
func parseValidationErrors(err error) []ValidationError {
	var validationErrors []ValidationError
	
	// 這裡可以根據使用的驗證庫來解析錯誤
	// 例如，如果使用 go-playground/validator
	
	errStr := err.Error()
	if strings.Contains(errStr, "required") {
		field := extractFieldName(errStr)
		validationErrors = append(validationErrors, ValidationError{
			Field:   field,
			Message: "This field is required",
		})
	}
	
	return validationErrors
}

// extractFieldName 從錯誤消息中提取字段名
func extractFieldName(errMsg string) string {
	// 簡單實現，實際應根據驗證庫的錯誤格式來解析
	parts := strings.Split(errMsg, " ")
	if len(parts) > 0 {
		return parts[0]
	}
	return "unknown"
}