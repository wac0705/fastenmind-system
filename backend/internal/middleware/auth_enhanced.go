package middleware

import (
	"crypto/subtle"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/fastenmind/fastener-api/internal/infrastructure/cache"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// EnhancedAuthConfig defines enhanced authentication configuration
type EnhancedAuthConfig struct {
	// Token validation
	ValidateToken func(token string) (map[string]interface{}, error)
	
	// Permission checker
	CheckPermission func(user map[string]interface{}, resource, action string) bool
	
	// Cache for permissions
	Cache *cache.RedisCache
	
	// API key validation
	ValidateAPIKey func(key string) (map[string]interface{}, error)
	
	// Multi-factor authentication
	RequireMFA bool
	ValidateMFA func(user map[string]interface{}, code string) bool
}

// EnhancedAuth returns an enhanced authentication middleware
func EnhancedAuth(config EnhancedAuthConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Check for API key first
			apiKey := c.Request().Header.Get("X-API-Key")
			if apiKey != "" {
				user, err := config.ValidateAPIKey(apiKey)
				if err != nil {
					return echo.NewHTTPError(http.StatusUnauthorized, "invalid API key")
				}
				c.Set("user", user)
				c.Set("auth_type", "api_key")
				return next(c)
			}
			
			// Check for Bearer token
			auth := c.Request().Header.Get("Authorization")
			if auth == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization header")
			}
			
			parts := strings.Split(auth, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization format")
			}
			
			token := parts[1]
			user, err := config.ValidateToken(token)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}
			
			// Check MFA if required
			if config.RequireMFA {
				mfaCode := c.Request().Header.Get("X-MFA-Code")
				if mfaCode == "" || !config.ValidateMFA(user, mfaCode) {
					return echo.NewHTTPError(http.StatusUnauthorized, "MFA validation failed")
				}
			}
			
			c.Set("user", user)
			c.Set("auth_type", "bearer")
			
			return next(c)
		}
	}
}

// RequirePermission checks if user has required permission
func RequirePermission(resource, action string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user")
			if user == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
			}
			
			userMap, ok := user.(map[string]interface{})
			if !ok {
				return echo.NewHTTPError(http.StatusInternalServerError, "invalid user data")
			}
			
			// Check permission
			permission := fmt.Sprintf("%s:%s", resource, action)
			if !hasPermission(userMap, permission) {
				return echo.NewHTTPError(http.StatusForbidden, "insufficient permissions")
			}
			
			return next(c)
		}
	}
}

// hasPermission checks if user has a specific permission
func hasPermission(user map[string]interface{}, permission string) bool {
	// Check role-based permissions
	if role, ok := user["role"].(string); ok {
		rolePermissions := getRolePermissions(role)
		for _, p := range rolePermissions {
			if p == permission || p == "*" {
				return true
			}
		}
	}
	
	// Check user-specific permissions
	if permissions, ok := user["permissions"].([]string); ok {
		for _, p := range permissions {
			if p == permission || p == "*" {
				return true
			}
		}
	}
	
	return false
}

// getRolePermissions returns permissions for a role
func getRolePermissions(role string) []string {
	rolePermissions := map[string][]string{
		"admin": {"*"}, // Admin has all permissions
		"manager": {
			"inquiry:*", "quote:*", "order:*", "customer:*",
			"report:read", "system:read",
		},
		"engineer": {
			"inquiry:read", "inquiry:update",
			"quote:create", "quote:read", "quote:update",
			"product:read", "process:read",
		},
		"sales": {
			"inquiry:create", "inquiry:read", "inquiry:update",
			"quote:read", "order:create", "order:read",
			"customer:create", "customer:read", "customer:update",
		},
		"viewer": {
			"inquiry:read", "quote:read", "order:read",
			"customer:read", "product:read", "report:read",
		},
	}
	
	return rolePermissions[role]
}

// SecureHeaders adds security headers to responses
func SecureHeaders() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Security headers
			c.Response().Header().Set("X-Content-Type-Options", "nosniff")
			c.Response().Header().Set("X-Frame-Options", "DENY")
			c.Response().Header().Set("X-XSS-Protection", "1; mode=block")
			c.Response().Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			c.Response().Header().Set("Content-Security-Policy", "default-src 'self'")
			
			// Remove server header
			c.Response().Header().Del("Server")
			
			return next(c)
		}
	}
}

// APIKeyAuth validates API key authentication
func APIKeyAuth(validateKey func(key string) (map[string]interface{}, error)) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			apiKey := c.Request().Header.Get("X-API-Key")
			if apiKey == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing API key")
			}
			
			user, err := validateKey(apiKey)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid API key")
			}
			
			c.Set("user", user)
			c.Set("auth_type", "api_key")
			
			return next(c)
		}
	}
}

// IPWhitelist restricts access to whitelisted IPs
func IPWhitelist(whitelist []string) echo.MiddlewareFunc {
	// Convert to map for O(1) lookup
	allowed := make(map[string]bool)
	for _, ip := range whitelist {
		allowed[ip] = true
	}
	
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			clientIP := c.RealIP()
			
			if !allowed[clientIP] {
				return echo.NewHTTPError(http.StatusForbidden, "access denied")
			}
			
			return next(c)
		}
	}
}

// RequestSigning validates request signatures
type RequestSigningConfig struct {
	// Secret key for signing
	Secret string
	
	// Header containing signature
	SignatureHeader string
	
	// Headers to include in signature
	SignedHeaders []string
	
	// Timestamp tolerance
	TimestampTolerance time.Duration
}

// RequestSigning validates request signatures
func RequestSigning(config RequestSigningConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			signature := c.Request().Header.Get(config.SignatureHeader)
			if signature == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing signature")
			}
			
			// Validate timestamp
			timestamp := c.Request().Header.Get("X-Timestamp")
			if timestamp == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing timestamp")
			}
			
			reqTime, err := time.Parse(time.RFC3339, timestamp)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "invalid timestamp format")
			}
			
			if time.Since(reqTime) > config.TimestampTolerance {
				return echo.NewHTTPError(http.StatusUnauthorized, "request expired")
			}
			
			// Build string to sign
			var parts []string
			parts = append(parts, c.Request().Method)
			parts = append(parts, c.Request().URL.Path)
			
			for _, header := range config.SignedHeaders {
				value := c.Request().Header.Get(header)
				parts = append(parts, fmt.Sprintf("%s:%s", header, value))
			}
			
			stringToSign := strings.Join(parts, "\n")
			
			// Calculate expected signature
			expectedSig := calculateHMAC(stringToSign, config.Secret)
			
			// Constant time comparison
			if subtle.ConstantTimeCompare([]byte(signature), []byte(expectedSig)) != 1 {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid signature")
			}
			
			return next(c)
		}
	}
}

// calculateHMAC calculates HMAC-SHA256 signature
func calculateHMAC(data, secret string) string {
	// Implementation depends on crypto package
	// This is a placeholder
	return "calculated-hmac"
}

// SessionAuth provides session-based authentication
type SessionAuth struct {
	Store *cache.RedisCache
	TTL   time.Duration
}

// Middleware returns session authentication middleware
func (s *SessionAuth) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get session ID from cookie
			cookie, err := c.Cookie("session_id")
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "no session")
			}
			
			// Validate session ID format
			sessionID, err := uuid.Parse(cookie.Value)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid session")
			}
			
			// Get session from cache
			ctx := c.Request().Context()
			sessionKey := fmt.Sprintf("session:%s", sessionID.String())
			
			var session map[string]interface{}
			err = s.Store.GetJSON(ctx, sessionKey, &session)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "session expired")
			}
			
			// Extend session TTL
			s.Store.Expire(ctx, sessionKey, s.TTL)
			
			// Set user in context
			c.Set("user", session["user"])
			c.Set("session", session)
			
			return next(c)
		}
	}
}

// CreateSession creates a new session
func (s *SessionAuth) CreateSession(c echo.Context, user map[string]interface{}) error {
	sessionID := uuid.New()
	session := map[string]interface{}{
		"id":         sessionID.String(),
		"user":       user,
		"created_at": time.Now(),
		"ip":         c.RealIP(),
		"user_agent": c.Request().UserAgent(),
	}
	
	// Store session
	ctx := c.Request().Context()
	sessionKey := fmt.Sprintf("session:%s", sessionID.String())
	err := s.Store.SetJSON(ctx, sessionKey, session, s.TTL)
	if err != nil {
		return err
	}
	
	// Set cookie
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID.String(),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(s.TTL),
	}
	c.SetCookie(cookie)
	
	return nil
}

// DestroySession destroys a session
func (s *SessionAuth) DestroySession(c echo.Context) error {
	cookie, err := c.Cookie("session_id")
	if err != nil {
		return nil // No session to destroy
	}
	
	// Delete from cache
	ctx := c.Request().Context()
	sessionKey := fmt.Sprintf("session:%s", cookie.Value)
	s.Store.Delete(ctx, sessionKey)
	
	// Clear cookie
	cookie.MaxAge = -1
	c.SetCookie(cookie)
	
	return nil
}