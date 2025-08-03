package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Response structure
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Success sends a success response
func Success(c echo.Context, data interface{}, message ...string) error {
	resp := Response{
		Success: true,
		Data:    data,
	}
	if len(message) > 0 {
		resp.Message = message[0]
	}
	return c.JSON(http.StatusOK, resp)
}

// SuccessWithMessage sends a success response with message
func SuccessWithMessage(c echo.Context, message string, data interface{}) error {
	return c.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error sends an error response
func Error(c echo.Context, status int, message string, err ...error) error {
	resp := Response{
		Success: false,
		Error:   message,
	}
	// If an error is provided, we can optionally log it but not send it to the client
	// This maintains backward compatibility while accepting the error parameter
	return c.JSON(status, resp)
}

// Created sends a created response
func Created(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    data,
	})
}

// NoContent sends a no content response
func NoContent(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

// PaginationResponse structure for paginated responses
type PaginationResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Error      string      `json:"error,omitempty"`
	Pagination Pagination  `json:"pagination"`
}

// Pagination structure
type Pagination struct {
	Total       int `json:"total"`
	Page        int `json:"page"`
	Limit       int `json:"limit"`
	TotalPages  int `json:"total_pages"`
}

// SuccessWithPagination sends a success response with pagination
func SuccessWithPagination(c echo.Context, message string, data interface{}, total, page, limit int) error {
	totalPages := (total + limit - 1) / limit
	return c.JSON(http.StatusOK, PaginationResponse{
		Success: true,
		Message: message,
		Data:    data,
		Pagination: Pagination{
			Total:      total,
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
		},
	})
}