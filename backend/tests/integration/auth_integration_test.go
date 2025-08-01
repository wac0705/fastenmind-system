package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"github.com/fastenmind/fastener-api/internal/handler"
	"github.com/fastenmind/fastener-api/internal/model"
	"github.com/fastenmind/fastener-api/internal/repository"
	"github.com/fastenmind/fastener-api/internal/service"
	"github.com/fastenmind/fastener-api/pkg/database"
)

type AuthIntegrationTestSuite struct {
	suite.Suite
	db          *gorm.DB
	echo        *echo.Echo
	authHandler *handler.AuthHandler
	cleanup     func()
}

func (suite *AuthIntegrationTestSuite) SetupSuite() {
	// Setup test database connection
	config := &database.Config{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "password",
		DBName:   "fastenmind_test",
		SSLMode:  "disable",
	}
	
	db, err := database.Connect(config)
	suite.Require().NoError(err)
	suite.db = db
	
	// Auto-migrate tables
	err = db.AutoMigrate(
		&model.Company{},
		&model.Account{},
	)
	suite.Require().NoError(err)
	
	// Setup services and handlers
	accountRepo := repository.NewAccountRepository(db)
	authService := service.NewAuthService(accountRepo)
	suite.authHandler = handler.NewAuthHandler(authService)
	
	// Setup Echo
	suite.echo = echo.New()
	
	// Setup cleanup function
	suite.cleanup = func() {
		// Clean up test data
		db.Exec("DELETE FROM accounts")
		db.Exec("DELETE FROM companies")
	}
}

func (suite *AuthIntegrationTestSuite) TearDownSuite() {
	if suite.cleanup != nil {
		suite.cleanup()
	}
}

func (suite *AuthIntegrationTestSuite) SetupTest() {
	// Clean up before each test
	suite.cleanup()
}

func (suite *AuthIntegrationTestSuite) TestAuthFlow_Complete() {
	// Create test company
	company := &model.Company{
		Code:    "TEST001",
		Name:    "Test Company",
		Country: "US",
		Type:    "headquarters",
	}
	err := suite.db.Create(company).Error
	suite.Require().NoError(err)
	
	// Create test account
	account := &model.Account{
		CompanyID:    company.ID,
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "$2a$10$XYZ...", // Pre-hashed password for "password123"
		FullName:     "Test User",
		Role:         "engineer",
		IsActive:     true,
	}
	err = suite.db.Create(account).Error
	suite.Require().NoError(err)
	
	// Test 1: Login with valid credentials
	loginRequest := handler.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}
	
	requestBody, _ := json.Marshal(loginRequest)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)
	
	err = suite.authHandler.Login(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
	
	var loginResponse handler.LoginResponse
	err = json.Unmarshal(rec.Body.Bytes(), &loginResponse)
	suite.NoError(err)
	suite.NotEmpty(loginResponse.AccessToken)
	suite.NotEmpty(loginResponse.RefreshToken)
	suite.Equal("testuser", loginResponse.User.Username)
	
	// Test 2: Get profile with access token
	req = httptest.NewRequest(http.MethodGet, "/auth/profile", nil)
	req.Header.Set("Authorization", "Bearer "+loginResponse.AccessToken)
	rec = httptest.NewRecorder()
	c = suite.echo.NewContext(req, rec)
	
	err = suite.authHandler.GetProfile(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
	
	var profileResponse handler.ProfileResponse
	err = json.Unmarshal(rec.Body.Bytes(), &profileResponse)
	suite.NoError(err)
	suite.Equal("testuser", profileResponse.Username)
	suite.Equal("engineer", profileResponse.Role)
	
	// Test 3: Refresh token
	refreshRequest := handler.RefreshTokenRequest{
		RefreshToken: loginResponse.RefreshToken,
	}
	
	requestBody, _ = json.Marshal(refreshRequest)
	req = httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewBuffer(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = suite.echo.NewContext(req, rec)
	
	err = suite.authHandler.RefreshToken(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
	
	var refreshResponse handler.LoginResponse
	err = json.Unmarshal(rec.Body.Bytes(), &refreshResponse)
	suite.NoError(err)
	suite.NotEmpty(refreshResponse.AccessToken)
	suite.NotEmpty(refreshResponse.RefreshToken)
	
	// Test 4: Logout
	req = httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	req.Header.Set("Authorization", "Bearer "+refreshResponse.AccessToken)
	rec = httptest.NewRecorder()
	c = suite.echo.NewContext(req, rec)
	
	err = suite.authHandler.Logout(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *AuthIntegrationTestSuite) TestLogin_InvalidCredentials() {
	// Create test company
	company := &model.Company{
		Code:    "TEST001",
		Name:    "Test Company",
		Country: "US",
		Type:    "headquarters",
	}
	err := suite.db.Create(company).Error
	suite.Require().NoError(err)
	
	// Create test account
	account := &model.Account{
		CompanyID:    company.ID,
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "$2a$10$XYZ...", // Pre-hashed password for "password123"
		FullName:     "Test User",
		Role:         "engineer",
		IsActive:     true,
	}
	err = suite.db.Create(account).Error
	suite.Require().NoError(err)
	
	// Test with wrong password
	loginRequest := handler.LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}
	
	requestBody, _ := json.Marshal(loginRequest)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)
	
	err = suite.authHandler.Login(c)
	suite.Error(err)
	
	// Test with non-existent user
	loginRequest = handler.LoginRequest{
		Username: "nonexistent",
		Password: "password123",
	}
	
	requestBody, _ = json.Marshal(loginRequest)
	req = httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = suite.echo.NewContext(req, rec)
	
	err = suite.authHandler.Login(c)
	suite.Error(err)
}

func (suite *AuthIntegrationTestSuite) TestLogin_InactiveAccount() {
	// Create test company
	company := &model.Company{
		Code:    "TEST001",
		Name:    "Test Company",
		Country: "US",
		Type:    "headquarters",
	}
	err := suite.db.Create(company).Error
	suite.Require().NoError(err)
	
	// Create inactive test account
	account := &model.Account{
		CompanyID:    company.ID,
		Username:     "inactiveuser",
		Email:        "inactive@example.com",
		PasswordHash: "$2a$10$XYZ...", // Pre-hashed password for "password123"
		FullName:     "Inactive User",
		Role:         "engineer",
		IsActive:     false, // Inactive account
	}
	err = suite.db.Create(account).Error
	suite.Require().NoError(err)
	
	loginRequest := handler.LoginRequest{
		Username: "inactiveuser",
		Password: "password123",
	}
	
	requestBody, _ := json.Marshal(loginRequest)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)
	
	err = suite.authHandler.Login(c)
	suite.Error(err)
}

func TestAuthIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(AuthIntegrationTestSuite))
}