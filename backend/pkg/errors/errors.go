package errors

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

// ErrorType 錯誤類型
type ErrorType string

const (
	// 業務錯誤類型
	ErrorTypeValidation   ErrorType = "VALIDATION_ERROR"
	ErrorTypeNotFound     ErrorType = "NOT_FOUND"
	ErrorTypeUnauthorized ErrorType = "UNAUTHORIZED"
	ErrorTypeForbidden    ErrorType = "FORBIDDEN"
	ErrorTypeConflict     ErrorType = "CONFLICT"
	ErrorTypeBadRequest   ErrorType = "BAD_REQUEST"
	
	// 系統錯誤類型
	ErrorTypeInternal     ErrorType = "INTERNAL_ERROR"
	ErrorTypeDatabase     ErrorType = "DATABASE_ERROR"
	ErrorTypeNetwork      ErrorType = "NETWORK_ERROR"
	ErrorTypeTimeout      ErrorType = "TIMEOUT_ERROR"
	ErrorTypeRateLimit    ErrorType = "RATE_LIMIT_ERROR"
	
	// 外部服務錯誤
	ErrorTypeExternal     ErrorType = "EXTERNAL_SERVICE_ERROR"
	ErrorTypeIntegration  ErrorType = "INTEGRATION_ERROR"
)

// AppError 應用程式錯誤
type AppError struct {
	Type       ErrorType              `json:"type"`
	Message    string                 `json:"message"`
	Code       string                 `json:"code,omitempty"`
	Details    map[string]interface{} `json:"details,omitempty"`
	StatusCode int                    `json:"-"`
	Cause      error                  `json:"-"`
	Stack      []string               `json:"-"`
}

// Error 實現 error 介面
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Type, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Unwrap 實現 errors.Unwrap
func (e *AppError) Unwrap() error {
	return e.Cause
}

// WithDetails 添加詳細信息
func (e *AppError) WithDetails(details map[string]interface{}) *AppError {
	e.Details = details
	return e
}

// WithCode 添加錯誤代碼
func (e *AppError) WithCode(code string) *AppError {
	e.Code = code
	return e
}

// New 創建新的應用程式錯誤
func New(errType ErrorType, message string, statusCode int) *AppError {
	return &AppError{
		Type:       errType,
		Message:    message,
		StatusCode: statusCode,
		Stack:      captureStack(),
	}
}

// Wrap 包裝錯誤
func Wrap(err error, errType ErrorType, message string) *AppError {
	if appErr, ok := err.(*AppError); ok {
		// 如果已經是 AppError，保留原始堆疊
		appErr.Message = message + ": " + appErr.Message
		return appErr
	}
	
	return &AppError{
		Type:       errType,
		Message:    message,
		Cause:      err,
		StatusCode: getStatusCode(errType),
		Stack:      captureStack(),
	}
}

// 預定義錯誤建構函數
func NewValidationError(message string) *AppError {
	return New(ErrorTypeValidation, message, http.StatusBadRequest)
}

func NewNotFoundError(message string) *AppError {
	return New(ErrorTypeNotFound, message, http.StatusNotFound)
}

func NewUnauthorizedError(message string) *AppError {
	return New(ErrorTypeUnauthorized, message, http.StatusUnauthorized)
}

func NewForbiddenError(message string) *AppError {
	return New(ErrorTypeForbidden, message, http.StatusForbidden)
}

func NewConflictError(message string) *AppError {
	return New(ErrorTypeConflict, message, http.StatusConflict)
}

func NewBadRequestError(message string) *AppError {
	return New(ErrorTypeBadRequest, message, http.StatusBadRequest)
}

func NewInternalError(message string) *AppError {
	return New(ErrorTypeInternal, message, http.StatusInternalServerError)
}

func NewDatabaseError(message string, err error) *AppError {
	return Wrap(err, ErrorTypeDatabase, message)
}

func NewNetworkError(message string, err error) *AppError {
	return Wrap(err, ErrorTypeNetwork, message)
}

func NewTimeoutError(message string) *AppError {
	return New(ErrorTypeTimeout, message, http.StatusRequestTimeout)
}

func NewRateLimitError(message string) *AppError {
	return New(ErrorTypeRateLimit, message, http.StatusTooManyRequests)
}

func NewExternalServiceError(service string, err error) *AppError {
	return Wrap(err, ErrorTypeExternal, fmt.Sprintf("External service error: %s", service))
}

// IsAppError 檢查是否為應用程式錯誤
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// GetAppError 獲取應用程式錯誤
func GetAppError(err error) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	return nil
}

// getStatusCode 根據錯誤類型獲取狀態碼
func getStatusCode(errType ErrorType) int {
	switch errType {
	case ErrorTypeValidation, ErrorTypeBadRequest:
		return http.StatusBadRequest
	case ErrorTypeNotFound:
		return http.StatusNotFound
	case ErrorTypeUnauthorized:
		return http.StatusUnauthorized
	case ErrorTypeForbidden:
		return http.StatusForbidden
	case ErrorTypeConflict:
		return http.StatusConflict
	case ErrorTypeTimeout:
		return http.StatusRequestTimeout
	case ErrorTypeRateLimit:
		return http.StatusTooManyRequests
	case ErrorTypeDatabase, ErrorTypeNetwork, ErrorTypeExternal, ErrorTypeIntegration:
		return http.StatusBadGateway
	default:
		return http.StatusInternalServerError
	}
}

// captureStack 捕獲堆疊信息
func captureStack() []string {
	var stack []string
	
	// 跳過前3個調用（captureStack、New/Wrap、調用者）
	for i := 3; i < 10; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}
		
		// 過濾掉運行時和標準庫
		if strings.Contains(file, "runtime/") {
			continue
		}
		
		stack = append(stack, fmt.Sprintf("%s:%d %s", file, line, fn.Name()))
	}
	
	return stack
}

// ErrorResponse 錯誤響應
type ErrorResponse struct {
	Success bool                   `json:"success"`
	Error   ErrorInfo              `json:"error"`
	Meta    map[string]interface{} `json:"meta,omitempty"`
}

// ErrorInfo 錯誤信息
type ErrorInfo struct {
	Type    string                 `json:"type"`
	Message string                 `json:"message"`
	Code    string                 `json:"code,omitempty"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// ToResponse 轉換為錯誤響應
func (e *AppError) ToResponse() ErrorResponse {
	return ErrorResponse{
		Success: false,
		Error: ErrorInfo{
			Type:    string(e.Type),
			Message: e.Message,
			Code:    e.Code,
			Details: e.Details,
		},
	}
}

// ValidationError 驗證錯誤詳情
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   interface{} `json:"value,omitempty"`
}

// NewValidationErrorWithDetails 創建帶詳情的驗證錯誤
func NewValidationErrorWithDetails(message string, errors []ValidationError) *AppError {
	details := make(map[string]interface{})
	details["validation_errors"] = errors
	
	return New(ErrorTypeValidation, message, http.StatusBadRequest).
		WithDetails(details)
}