package handler

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// getUserIDFromContext gets the user ID from the context
func getUserIDFromContext(c echo.Context) uuid.UUID {
	userID, _ := c.Get("user_id").(uuid.UUID)
	return userID
}

// getCompanyIDFromContext gets the company ID from the context
func getCompanyIDFromContext(c echo.Context) uuid.UUID {
	companyID, _ := c.Get("company_id").(uuid.UUID)
	return companyID
}

// getUserIDFromContextWithError gets the user ID from the context with error handling
func getUserIDFromContextWithError(c echo.Context) (uuid.UUID, error) {
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return uuid.Nil, echo.NewHTTPError(400, "Invalid user ID")
	}
	return userID, nil
}

// getCompanyIDFromContextWithError gets the company ID from the context with error handling
func getCompanyIDFromContextWithError(c echo.Context) (uuid.UUID, error) {
	companyID, ok := c.Get("company_id").(uuid.UUID)
	if !ok {
		return uuid.Nil, echo.NewHTTPError(400, "Invalid company ID")
	}
	return companyID, nil
}