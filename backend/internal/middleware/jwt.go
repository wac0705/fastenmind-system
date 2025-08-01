package middleware

import (
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/google/uuid"
)

// JWTClaims represents the JWT claims
type JWTClaims struct {
	UserID    string `json:"user_id"`
	CompanyID string `json:"company_id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	jwt.RegisteredClaims
}

// JWT returns a JWT authentication middleware
func JWT(secretKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get token from header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(401, "missing authorization header")
			}

			// Check Bearer prefix
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				return echo.NewHTTPError(401, "invalid authorization header format")
			}

			tokenString := tokenParts[1]

			// Parse and validate token
			token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
				// Validate signing method
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(secretKey), nil
			})

			if err != nil {
				return echo.NewHTTPError(401, "invalid token: "+err.Error())
			}

			// Get claims
			claims, ok := token.Claims.(*JWTClaims)
			if !ok || !token.Valid {
				return echo.NewHTTPError(401, "invalid token claims")
			}

			// Set user context
			// Convert string IDs to UUID for consistency with handlers
			if userID, err := uuid.Parse(claims.UserID); err == nil {
				c.Set("user_id", userID)
			} else {
				c.Set("user_id", claims.UserID)
			}
			
			if companyID, err := uuid.Parse(claims.CompanyID); err == nil {
				c.Set("company_id", companyID)
			} else {
				c.Set("company_id", claims.CompanyID)
			}
			
			c.Set("email", claims.Email)
			c.Set("role", claims.Role)

			return next(c)
		}
	}
}

// GetUserID extracts user ID from context
func GetUserID(c echo.Context) string {
	if userID, ok := c.Get("user_id").(string); ok {
		return userID
	}
	return ""
}

// GetCompanyID extracts company ID from context
func GetCompanyID(c echo.Context) string {
	if companyID, ok := c.Get("company_id").(string); ok {
		return companyID
	}
	return ""
}

// GetRole extracts user role from context
func GetRole(c echo.Context) string {
	if role, ok := c.Get("role").(string); ok {
		return role
	}
	return ""
}

// Auth is an alias for JWT middleware with a default secret key
// In production, this should use a proper secret from configuration
func Auth() echo.MiddlewareFunc {
	// TODO: Get secret from configuration
	return JWT("your-secret-key")
}