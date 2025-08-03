package errors

import (
	"errors"
	"net/http"
	"testing"
)

func TestNewAppError(t *testing.T) {
	tests := []struct {
		name       string
		errType    ErrorType
		message    string
		statusCode int
	}{
		{
			name:       "Validation error",
			errType:    ErrorTypeValidation,
			message:    "Invalid input",
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Not found error",
			errType:    ErrorTypeNotFound,
			message:    "Resource not found",
			statusCode: http.StatusNotFound,
		},
		{
			name:       "Internal error",
			errType:    ErrorTypeInternal,
			message:    "Something went wrong",
			statusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := New(tt.errType, tt.message, tt.statusCode)
			
			if err.Type != tt.errType {
				t.Errorf("Expected type %s, got %s", tt.errType, err.Type)
			}
			if err.Message != tt.message {
				t.Errorf("Expected message %s, got %s", tt.message, err.Message)
			}
			if err.StatusCode != tt.statusCode {
				t.Errorf("Expected status code %d, got %d", tt.statusCode, err.StatusCode)
			}
			if len(err.Stack) == 0 {
				t.Error("Expected stack trace to be captured")
			}
		})
	}
}

func TestWrapError(t *testing.T) {
	originalErr := errors.New("database connection failed")
	wrappedErr := Wrap(originalErr, ErrorTypeDatabase, "Failed to fetch user")
	
	if wrappedErr.Type != ErrorTypeDatabase {
		t.Errorf("Expected type %s, got %s", ErrorTypeDatabase, wrappedErr.Type)
	}
	if wrappedErr.Cause != originalErr {
		t.Error("Expected original error to be preserved as cause")
	}
	if wrappedErr.StatusCode != http.StatusBadGateway {
		t.Errorf("Expected status code %d, got %d", http.StatusBadGateway, wrappedErr.StatusCode)
	}
	
	// Test unwrap
	unwrapped := wrappedErr.Unwrap()
	if unwrapped != originalErr {
		t.Error("Unwrap should return the original error")
	}
}

func TestAppErrorWithDetails(t *testing.T) {
	err := NewValidationError("Validation failed")
	details := map[string]interface{}{
		"field": "email",
		"value": "invalid-email",
	}
	
	err.WithDetails(details)
	
	if err.Details == nil {
		t.Error("Expected details to be set")
	}
	if err.Details["field"] != "email" {
		t.Errorf("Expected field to be 'email', got %v", err.Details["field"])
	}
}

func TestAppErrorWithCode(t *testing.T) {
	err := NewBadRequestError("Invalid request")
	err.WithCode("ERR_INVALID_FORMAT")
	
	if err.Code != "ERR_INVALID_FORMAT" {
		t.Errorf("Expected code 'ERR_INVALID_FORMAT', got %s", err.Code)
	}
}

func TestPredefinedErrors(t *testing.T) {
	tests := []struct {
		name       string
		errFunc    func(string) *AppError
		message    string
		errType    ErrorType
		statusCode int
	}{
		{
			name:       "NewValidationError",
			errFunc:    NewValidationError,
			message:    "Invalid email",
			errType:    ErrorTypeValidation,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "NewNotFoundError",
			errFunc:    NewNotFoundError,
			message:    "User not found",
			errType:    ErrorTypeNotFound,
			statusCode: http.StatusNotFound,
		},
		{
			name:       "NewUnauthorizedError",
			errFunc:    NewUnauthorizedError,
			message:    "Invalid token",
			errType:    ErrorTypeUnauthorized,
			statusCode: http.StatusUnauthorized,
		},
		{
			name:       "NewForbiddenError",
			errFunc:    NewForbiddenError,
			message:    "Access denied",
			errType:    ErrorTypeForbidden,
			statusCode: http.StatusForbidden,
		},
		{
			name:       "NewConflictError",
			errFunc:    NewConflictError,
			message:    "Resource already exists",
			errType:    ErrorTypeConflict,
			statusCode: http.StatusConflict,
		},
		{
			name:       "NewTimeoutError",
			errFunc:    NewTimeoutError,
			message:    "Request timeout",
			errType:    ErrorTypeTimeout,
			statusCode: http.StatusRequestTimeout,
		},
		{
			name:       "NewRateLimitError",
			errFunc:    NewRateLimitError,
			message:    "Too many requests",
			errType:    ErrorTypeRateLimit,
			statusCode: http.StatusTooManyRequests,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.errFunc(tt.message)
			
			if err.Type != tt.errType {
				t.Errorf("Expected type %s, got %s", tt.errType, err.Type)
			}
			if err.Message != tt.message {
				t.Errorf("Expected message %s, got %s", tt.message, err.Message)
			}
			if err.StatusCode != tt.statusCode {
				t.Errorf("Expected status code %d, got %d", tt.statusCode, err.StatusCode)
			}
		})
	}
}

func TestIsAppError(t *testing.T) {
	appErr := NewInternalError("Server error")
	stdErr := errors.New("standard error")
	
	if !IsAppError(appErr) {
		t.Error("Expected IsAppError to return true for AppError")
	}
	if IsAppError(stdErr) {
		t.Error("Expected IsAppError to return false for standard error")
	}
}

func TestGetAppError(t *testing.T) {
	appErr := NewInternalError("Server error")
	stdErr := errors.New("standard error")
	
	result := GetAppError(appErr)
	if result == nil {
		t.Error("Expected GetAppError to return the AppError")
	}
	if result != appErr {
		t.Error("Expected GetAppError to return the same AppError instance")
	}
	
	result = GetAppError(stdErr)
	if result != nil {
		t.Error("Expected GetAppError to return nil for standard error")
	}
}

func TestErrorResponse(t *testing.T) {
	err := NewValidationError("Invalid input").
		WithCode("VALIDATION_001").
		WithDetails(map[string]interface{}{
			"field": "age",
			"min":   18,
		})
	
	response := err.ToResponse()
	
	if response.Success != false {
		t.Error("Expected Success to be false")
	}
	if response.Error.Type != string(ErrorTypeValidation) {
		t.Errorf("Expected error type %s, got %s", ErrorTypeValidation, response.Error.Type)
	}
	if response.Error.Message != "Invalid input" {
		t.Errorf("Expected message 'Invalid input', got %s", response.Error.Message)
	}
	if response.Error.Code != "VALIDATION_001" {
		t.Errorf("Expected code 'VALIDATION_001', got %s", response.Error.Code)
	}
	if response.Error.Details == nil {
		t.Error("Expected details to be included in response")
	}
}

func TestNewValidationErrorWithDetails(t *testing.T) {
	validationErrors := []ValidationError{
		{
			Field:   "email",
			Message: "Invalid email format",
			Value:   "not-an-email",
		},
		{
			Field:   "age",
			Message: "Must be at least 18",
			Value:   15,
		},
	}
	
	err := NewValidationErrorWithDetails("Validation failed", validationErrors)
	
	if err.Type != ErrorTypeValidation {
		t.Errorf("Expected type %s, got %s", ErrorTypeValidation, err.Type)
	}
	if err.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, err.StatusCode)
	}
	
	if err.Details == nil {
		t.Fatal("Expected details to be set")
	}
	
	if errs, ok := err.Details["validation_errors"].([]ValidationError); ok {
		if len(errs) != 2 {
			t.Errorf("Expected 2 validation errors, got %d", len(errs))
		}
	} else {
		t.Error("Expected validation_errors in details")
	}
}

func TestErrorString(t *testing.T) {
	// Test error without cause
	err1 := NewInternalError("Server error")
	expected1 := "INTERNAL_ERROR: Server error"
	if err1.Error() != expected1 {
		t.Errorf("Expected error string '%s', got '%s'", expected1, err1.Error())
	}
	
	// Test error with cause
	cause := errors.New("connection refused")
	err2 := Wrap(cause, ErrorTypeDatabase, "Database error")
	expected2 := "DATABASE_ERROR: Database error (caused by: connection refused)"
	if err2.Error() != expected2 {
		t.Errorf("Expected error string '%s', got '%s'", expected2, err2.Error())
	}
}