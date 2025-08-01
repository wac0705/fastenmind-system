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
func Success(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
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
func Error(c echo.Context, status int, message string) error {
	return c.JSON(status, Response{
		Success: false,
		Error:   message,
	})
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