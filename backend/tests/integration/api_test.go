package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fastenmind/fastener-api/internal/config"
	"github.com/fastenmind/fastener-api/internal/handler"
	"github.com/fastenmind/fastener-api/internal/middleware"
	"github.com/fastenmind/fastener-api/internal/repository"
	"github.com/fastenmind/fastener-api/internal/service"
	"github.com/fastenmind/fastener-api/pkg/database"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type APITestSuite struct {
	suite.Suite
	echo      *echo.Echo
	db        *database.DBWrapper
	authToken string
	companyID uuid.UUID
	userID    uuid.UUID
}

func (suite *APITestSuite) SetupSuite() {
	// Initialize test configuration
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "test",
			Password: "test",
			DBName:   "fastenmind_test",
			SSLMode:  "disable",
		},
		JWT: config.JWTConfig{
			SecretKey:              "test-secret-key",
			AccessTokenExpiration:  "15m",
			RefreshTokenExpiration: "7d",
		},
	}

	// Initialize database
	db, err := database.NewWrapper(cfg.Database)
	suite.Require().NoError(err)
	suite.db = db

	// Run migrations
	suite.Require().NoError(db.AutoMigrate())

	// Initialize Echo
	e := echo.New()
	middleware.SetupErrorHandling(e)
	e.Use(middleware.Logger())
	e.Use(middleware.RequestID())
	e.Use(middleware.ValidationMiddleware())

	// Initialize repositories, services, and handlers
	repos := repository.NewRepositories(db.GormDB)
	services := service.NewServices(repos, cfg, db.GormDB)
	handlers := handler.NewHandlers(services)

	// Setup routes
	api := e.Group("/api/v1")
	
	// Auth routes
	api.POST("/auth/login", handlers.Auth.Login)
	api.POST("/auth/register", handlers.Auth.Register)
	
	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.JWT(cfg.JWT.SecretKey))
	
	// Customer routes
	protected.GET("/customers", handlers.Customer.List)
	protected.POST("/customers", handlers.Customer.Create)
	protected.GET("/customers/:id", handlers.Customer.Get)
	protected.PUT("/customers/:id", handlers.Customer.Update)
	protected.DELETE("/customers/:id", handlers.Customer.Delete)

	suite.echo = e

	// Create test user and login
	suite.setupTestUser()
}

func (suite *APITestSuite) TearDownSuite() {
	// Clean up database
	suite.db.Close()
}

func (suite *APITestSuite) setupTestUser() {
	// Create test company
	suite.companyID = uuid.New()
	
	// Create test user
	suite.userID = uuid.New()
	
	// Generate auth token
	// In a real test, you would create the user in the database and get a real token
	suite.authToken = "test-token"
}

func (suite *APITestSuite) TestHealthCheck() {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	rec := httptest.NewRecorder()
	
	suite.echo.ServeHTTP(rec, req)
	
	assert.Equal(suite.T(), http.StatusOK, rec.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "ok", response["status"])
}

func (suite *APITestSuite) TestCreateCustomer() {
	customer := map[string]interface{}{
		"name":         "Test Customer",
		"email":        "test@example.com",
		"phone":        "+1-555-123-4567",
		"country":      "US",
		"currency":     "USD",
		"payment_terms": 30,
		"credit_limit": 10000,
		"address": map[string]interface{}{
			"street":     "123 Test Street",
			"city":       "Test City",
			"state":      "TS",
			"country":    "US",
			"postal_code": "12345",
		},
	}
	
	body, _ := json.Marshal(customer)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/customers", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+suite.authToken)
	rec := httptest.NewRecorder()
	
	suite.echo.ServeHTTP(rec, req)
	
	assert.Equal(suite.T(), http.StatusCreated, rec.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), response["id"])
	assert.Equal(suite.T(), "Test Customer", response["name"])
}

func (suite *APITestSuite) TestCreateCustomerValidation() {
	testCases := []struct {
		name       string
		customer   map[string]interface{}
		wantStatus int
	}{
		{
			name: "Missing required field",
			customer: map[string]interface{}{
				"email": "test@example.com",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid email",
			customer: map[string]interface{}{
				"name":  "Test Customer",
				"email": "invalid-email",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid phone",
			customer: map[string]interface{}{
				"name":  "Test Customer",
				"email": "test@example.com",
				"phone": "123",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid country code",
			customer: map[string]interface{}{
				"name":    "Test Customer",
				"email":   "test@example.com",
				"country": "USA", // Should be 2-letter code
			},
			wantStatus: http.StatusBadRequest,
		},
	}
	
	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.customer)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/customers", bytes.NewReader(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set(echo.HeaderAuthorization, "Bearer "+suite.authToken)
			rec := httptest.NewRecorder()
			
			suite.echo.ServeHTTP(rec, req)
			
			assert.Equal(t, tc.wantStatus, rec.Code)
		})
	}
}

func (suite *APITestSuite) TestListCustomers() {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/customers?page=1&page_size=10", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+suite.authToken)
	rec := httptest.NewRecorder()
	
	suite.echo.ServeHTTP(rec, req)
	
	assert.Equal(suite.T(), http.StatusOK, rec.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Contains(suite.T(), response, "data")
	assert.Contains(suite.T(), response, "pagination")
}

func (suite *APITestSuite) TestRateLimiting() {
	// Test that rate limiting works
	customer := map[string]interface{}{
		"name": "Test Customer",
	}
	body, _ := json.Marshal(customer)
	
	// Make many requests quickly
	for i := 0; i < 100; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/customers", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+suite.authToken)
		rec := httptest.NewRecorder()
		
		suite.echo.ServeHTTP(rec, req)
		
		// At some point, we should get rate limited
		if rec.Code == http.StatusTooManyRequests {
			return
		}
	}
	
	// If we didn't get rate limited, the test should fail
	suite.T().Error("Expected rate limiting to kick in")
}

func (suite *APITestSuite) TestConcurrentRequests() {
	// Test concurrent access to ensure thread safety
	done := make(chan bool, 10)
	
	for i := 0; i < 10; i++ {
		go func(id int) {
			defer func() { done <- true }()
			
			req := httptest.NewRequest(http.MethodGet, "/api/v1/customers", nil)
			req.Header.Set(echo.HeaderAuthorization, "Bearer "+suite.authToken)
			rec := httptest.NewRecorder()
			
			suite.echo.ServeHTTP(rec, req)
			
			assert.Equal(suite.T(), http.StatusOK, rec.Code)
		}(i)
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			suite.T().Error("Timeout waiting for concurrent requests")
		}
	}
}

func TestAPITestSuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}