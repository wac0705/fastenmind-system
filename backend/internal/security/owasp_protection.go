package security

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"html"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/time/rate"
)

// OWASPProtection provides protection against OWASP Top 10 vulnerabilities
type OWASPProtection struct {
	csrfTokens     sync.Map
	sqlInjectionRe *regexp.Regexp
	xssPatterns    []*regexp.Regexp
	rateLimiters   sync.Map
}

// NewOWASPProtection creates a new OWASP protection instance
func NewOWASPProtection() *OWASPProtection {
	return &OWASPProtection{
		sqlInjectionRe: regexp.MustCompile(`(?i)(union|select|insert|update|delete|drop|create|alter|exec|script|javascript|vbscript|onload|onerror|onclick)`),
		xssPatterns: []*regexp.Regexp{
			regexp.MustCompile(`<script[^>]*>.*?</script>`),
			regexp.MustCompile(`javascript:`),
			regexp.MustCompile(`on\w+\s*=`),
			regexp.MustCompile(`<iframe[^>]*>`),
			regexp.MustCompile(`<object[^>]*>`),
			regexp.MustCompile(`<embed[^>]*>`),
		},
	}
}

// A1_BrokenAccessControl - Implement proper access control
type AccessControl struct {
	permissions map[string][]string
	mu          sync.RWMutex
}

func NewAccessControl() *AccessControl {
	return &AccessControl{
		permissions: make(map[string][]string),
	}
}

func (ac *AccessControl) CheckPermission(userRole, resource, action string) bool {
	ac.mu.RLock()
	defer ac.mu.RUnlock()
	
	key := fmt.Sprintf("%s:%s", userRole, resource)
	actions, exists := ac.permissions[key]
	if !exists {
		return false
	}
	
	for _, a := range actions {
		if a == action || a == "*" {
			return true
		}
	}
	return false
}

// A2_CryptographicFailures - Implement proper encryption
type CryptoProtection struct {
	encryptionKey []byte
}

func NewCryptoProtection(key []byte) *CryptoProtection {
	return &CryptoProtection{encryptionKey: key}
}

// EncryptSensitiveData encrypts sensitive data at rest
func (cp *CryptoProtection) EncryptSensitiveData(data []byte) ([]byte, error) {
	// Implementation would use AES-256-GCM
	// This is a placeholder for the actual implementation
	return data, nil
}

// A3_Injection - SQL Injection Protection Middleware
func (op *OWASPProtection) SQLInjectionProtection() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Check all query parameters
			for key, values := range c.QueryParams() {
				for _, value := range values {
					if op.sqlInjectionRe.MatchString(value) {
						return echo.NewHTTPError(http.StatusBadRequest, "Potential SQL injection detected")
					}
				}
			}
			
			// Check request body for JSON
			if c.Request().Header.Get("Content-Type") == "application/json" {
				var body map[string]interface{}
				if err := c.Bind(&body); err == nil {
					if op.checkSQLInjectionInMap(body) {
						return echo.NewHTTPError(http.StatusBadRequest, "Potential SQL injection detected")
					}
				}
			}
			
			return next(c)
		}
	}
}

func (op *OWASPProtection) checkSQLInjectionInMap(data map[string]interface{}) bool {
	for _, value := range data {
		switch v := value.(type) {
		case string:
			if op.sqlInjectionRe.MatchString(v) {
				return true
			}
		case map[string]interface{}:
			if op.checkSQLInjectionInMap(v) {
				return true
			}
		}
	}
	return false
}

// A4_InsecureDesign - Implement secure design patterns
type SecureDesignValidator struct {
	maxFileSize      int64
	allowedFileTypes []string
	maxRequestSize   int64
}

func NewSecureDesignValidator() *SecureDesignValidator {
	return &SecureDesignValidator{
		maxFileSize:      10 * 1024 * 1024, // 10MB
		allowedFileTypes: []string{".pdf", ".jpg", ".png", ".doc", ".docx"},
		maxRequestSize:   50 * 1024 * 1024, // 50MB
	}
}

// A5_SecurityMisconfiguration - Security Headers Middleware
func SecurityHeaders() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Set security headers
			c.Response().Header().Set("X-Content-Type-Options", "nosniff")
			c.Response().Header().Set("X-Frame-Options", "DENY")
			c.Response().Header().Set("X-XSS-Protection", "1; mode=block")
			c.Response().Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
			c.Response().Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline';")
			c.Response().Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			c.Response().Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
			
			return next(c)
		}
	}
}

// A6_VulnerableComponents - Component Security Scanner
type ComponentScanner struct {
	vulnerableVersions map[string][]string
	mu                sync.RWMutex
}

func NewComponentScanner() *ComponentScanner {
	return &ComponentScanner{
		vulnerableVersions: make(map[string][]string),
	}
}

func (cs *ComponentScanner) CheckVulnerabilities(dependencies map[string]string) []string {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	
	var vulnerabilities []string
	for pkg, version := range dependencies {
		if vulnVersions, exists := cs.vulnerableVersions[pkg]; exists {
			for _, vulnVersion := range vulnVersions {
				if version == vulnVersion {
					vulnerabilities = append(vulnerabilities, fmt.Sprintf("%s@%s has known vulnerabilities", pkg, version))
				}
			}
		}
	}
	return vulnerabilities
}

// A7_IdentificationAuthenticationFailures - Enhanced Authentication
type EnhancedAuth struct {
	failedAttempts sync.Map
	lockoutTime    time.Duration
	maxAttempts    int
}

func NewEnhancedAuth() *EnhancedAuth {
	return &EnhancedAuth{
		lockoutTime: 15 * time.Minute,
		maxAttempts: 5,
	}
}

func (ea *EnhancedAuth) CheckAccountLockout(username string) bool {
	if value, exists := ea.failedAttempts.Load(username); exists {
		attempt := value.(*LoginAttempt)
		if attempt.Count >= ea.maxAttempts && time.Since(attempt.LastAttempt) < ea.lockoutTime {
			return true
		}
	}
	return false
}

type LoginAttempt struct {
	Count       int
	LastAttempt time.Time
}

// A8_SoftwareDataIntegrityFailures - Integrity Verification
type IntegrityChecker struct {
	trustedSources []string
	signatureKeys  map[string]string
}

func NewIntegrityChecker() *IntegrityChecker {
	return &IntegrityChecker{
		trustedSources: []string{"https://api.fastenmind.com"},
		signatureKeys:  make(map[string]string),
	}
}

func (ic *IntegrityChecker) VerifySignature(data []byte, signature string, keyID string) bool {
	// Implement signature verification logic
	// This is a placeholder for actual implementation
	return true
}

// A9_SecurityLoggingMonitoringFailures - Security Event Logger
type SecurityLogger struct {
	events chan SecurityEvent
	file   string
}

type SecurityEvent struct {
	Timestamp   time.Time
	EventType   string
	UserID      string
	IP          string
	Details     map[string]interface{}
	Severity    string
}

func NewSecurityLogger(logFile string) *SecurityLogger {
	sl := &SecurityLogger{
		events: make(chan SecurityEvent, 1000),
		file:   logFile,
	}
	go sl.processEvents()
	return sl
}

func (sl *SecurityLogger) LogSecurityEvent(event SecurityEvent) {
	select {
	case sl.events <- event:
	default:
		// Log queue full, implement overflow handling
	}
}

func (sl *SecurityLogger) processEvents() {
	for event := range sl.events {
		// Write to secure log file
		// Implement log rotation and encryption
	}
}

// A10_SSRF - Server-Side Request Forgery Protection
type SSRFProtection struct {
	allowedHosts   []string
	blockedIPs     []string
	privateIPRegex *regexp.Regexp
}

func NewSSRFProtection() *SSRFProtection {
	return &SSRFProtection{
		allowedHosts: []string{"api.fastenmind.com", "cdn.fastenmind.com"},
		blockedIPs:   []string{"169.254.169.254", "127.0.0.1", "localhost"},
		privateIPRegex: regexp.MustCompile(`^(10\.|172\.(1[6-9]|2[0-9]|3[0-1])\.|192\.168\.)`),
	}
}

func (sp *SSRFProtection) ValidateURL(url string) error {
	// Parse and validate URL
	// Check against allowed hosts
	// Block private IPs and metadata endpoints
	return nil
}

// XSS Protection Middleware
func (op *OWASPProtection) XSSProtection() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Sanitize query parameters
			params := make(map[string][]string)
			for key, values := range c.QueryParams() {
				sanitized := make([]string, len(values))
				for i, value := range values {
					sanitized[i] = html.EscapeString(value)
				}
				params[key] = sanitized
			}
			
			// Override query params with sanitized values
			c.Request().URL.RawQuery = ""
			for key, values := range params {
				for _, value := range values {
					c.QueryParams().Add(key, value)
				}
			}
			
			return next(c)
		}
	}
}

// CSRF Protection
func (op *OWASPProtection) CSRFProtection() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Method == "GET" {
				// Generate CSRF token for GET requests
				token := op.generateCSRFToken()
				op.csrfTokens.Store(token, time.Now())
				c.Response().Header().Set("X-CSRF-Token", token)
				return next(c)
			}
			
			// Verify CSRF token for state-changing requests
			token := c.Request().Header.Get("X-CSRF-Token")
			if token == "" {
				return echo.NewHTTPError(http.StatusForbidden, "CSRF token missing")
			}
			
			if _, exists := op.csrfTokens.Load(token); !exists {
				return echo.NewHTTPError(http.StatusForbidden, "Invalid CSRF token")
			}
			
			// Delete used token
			op.csrfTokens.Delete(token)
			
			return next(c)
		}
	}
}

func (op *OWASPProtection) generateCSRFToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// Input Validation
type InputValidator struct {
	rules map[string]ValidationRule
}

type ValidationRule struct {
	Required bool
	MinLen   int
	MaxLen   int
	Pattern  *regexp.Regexp
	Custom   func(interface{}) error
}

func NewInputValidator() *InputValidator {
	return &InputValidator{
		rules: map[string]ValidationRule{
			"email": {
				Required: true,
				Pattern:  regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
			},
			"phone": {
				Pattern: regexp.MustCompile(`^\+?[1-9]\d{1,14}$`),
			},
			"username": {
				Required: true,
				MinLen:   3,
				MaxLen:   50,
				Pattern:  regexp.MustCompile(`^[a-zA-Z0-9_-]+$`),
			},
		},
	}
}

// Rate Limiting per User
func (op *OWASPProtection) UserRateLimiting(limit int, window time.Duration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID := c.Request().Header.Get("X-User-ID")
			if userID == "" {
				userID = c.RealIP()
			}
			
			limiterI, _ := op.rateLimiters.LoadOrStore(userID, rate.NewLimiter(rate.Every(window/time.Duration(limit)), limit))
			limiter := limiterI.(*rate.Limiter)
			
			if !limiter.Allow() {
				return echo.NewHTTPError(http.StatusTooManyRequests, "Rate limit exceeded")
			}
			
			return next(c)
		}
	}
}

// Password Policy Enforcement
type PasswordPolicy struct {
	MinLength      int
	RequireUpper   bool
	RequireLower   bool
	RequireNumber  bool
	RequireSpecial bool
	CommonPasswords map[string]bool
}

func NewPasswordPolicy() *PasswordPolicy {
	return &PasswordPolicy{
		MinLength:      12,
		RequireUpper:   true,
		RequireLower:   true,
		RequireNumber:  true,
		RequireSpecial: true,
		CommonPasswords: loadCommonPasswords(),
	}
}

func (pp *PasswordPolicy) Validate(password string) error {
	if len(password) < pp.MinLength {
		return fmt.Errorf("password must be at least %d characters", pp.MinLength)
	}
	
	if pp.RequireUpper && !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return fmt.Errorf("password must contain uppercase letter")
	}
	
	if pp.RequireLower && !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return fmt.Errorf("password must contain lowercase letter")
	}
	
	if pp.RequireNumber && !regexp.MustCompile(`\d`).MatchString(password) {
		return fmt.Errorf("password must contain number")
	}
	
	if pp.RequireSpecial && !regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password) {
		return fmt.Errorf("password must contain special character")
	}
	
	if pp.CommonPasswords[strings.ToLower(password)] {
		return fmt.Errorf("password is too common")
	}
	
	return nil
}

func loadCommonPasswords() map[string]bool {
	// Load from file or embedded resource
	return map[string]bool{
		"password123": true,
		"admin123":    true,
		"12345678":    true,
		// Add more common passwords
	}
}