package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/fastenmind/fastener-api/internal/handler"
	"github.com/fastenmind/fastener-api/internal/model"
	"github.com/fastenmind/fastener-api/internal/service"
)

// MockAuthService is a mock implementation of AuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Login(ctx context.Context, username, password string) (*service.LoginResult, error) {
	args := m.Called(ctx, username, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.LoginResult), args.Error(1)
}

func (m *MockAuthService) ValidateToken(ctx context.Context, token string) (*service.JWTClaims, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.JWTClaims), args.Error(1)
}

func (m *MockAuthService) RefreshToken(ctx context.Context, refreshToken string) (*service.LoginResult, error) {
	args := m.Called(ctx, refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.LoginResult), args.Error(1)
}

func (m *MockAuthService) Logout(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func TestAuthHandler_Login_Success(t *testing.T) {
	// Arrange
	e := echo.New()
	mockAuthService := new(MockAuthService)
	authHandler := handler.NewAuthHandler(mockAuthService)
	
	loginRequest := handler.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}
	
	loginResult := &service.LoginResult{
		Account: &model.Account{
			ID:        uuid.New(),
			Username:  "testuser",
			Email:     "test@example.com",
			FullName:  "Test User",
			Role:      "engineer",
			CompanyID: uuid.New(),
		},
		AccessToken:  "access_token_123",
		RefreshToken: "refresh_token_123",
		ExpiresIn:    3600,
	}
	
	mockAuthService.On("Login", mock.Anything, "testuser", "password123").Return(loginResult, nil)
	
	requestBody, _ := json.Marshal(loginRequest)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	
	// Act
	err := authHandler.Login(c)
	
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	
	var response handler.LoginResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, loginResult.AccessToken, response.AccessToken)
	assert.Equal(t, loginResult.RefreshToken, response.RefreshToken)
	assert.Equal(t, loginResult.Account.Username, response.User.Username)
	
	mockAuthService.AssertExpectations(t)
}

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	// Arrange
	e := echo.New()
	mockAuthService := new(MockAuthService)
	authHandler := handler.NewAuthHandler(mockAuthService)
	
	loginRequest := handler.LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}
	
	mockAuthService.On("Login", mock.Anything, "testuser", "wrongpassword").Return(nil, assert.AnError)
	
	requestBody, _ := json.Marshal(loginRequest)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	
	// Act
	err := authHandler.Login(c)
	
	// Assert
	assert.Error(t, err)
	
	mockAuthService.AssertExpectations(t)
}

func TestAuthHandler_Login_ValidationError(t *testing.T) {
	// Arrange
	e := echo.New()
	mockAuthService := new(MockAuthService)
	authHandler := handler.NewAuthHandler(mockAuthService)
	
	// Invalid request with missing username
	loginRequest := handler.LoginRequest{
		Password: "password123",
	}
	
	requestBody, _ := json.Marshal(loginRequest)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	
	// Act
	err := authHandler.Login(c)
	
	// Assert
	assert.Error(t, err)
	httpError, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, httpError.Code)
}

func TestAuthHandler_RefreshToken_Success(t *testing.T) {
	// Arrange
	e := echo.New()
	mockAuthService := new(MockAuthService)
	authHandler := handler.NewAuthHandler(mockAuthService)
	
	refreshRequest := handler.RefreshTokenRequest{
		RefreshToken: "refresh_token_123",
	}
	
	loginResult := &service.LoginResult{
		Account: &model.Account{
			ID:        uuid.New(),
			Username:  "testuser",
			Email:     "test@example.com",
			FullName:  "Test User",
			Role:      "engineer",
			CompanyID: uuid.New(),
		},
		AccessToken:  "new_access_token_123",
		RefreshToken: "new_refresh_token_123",
		ExpiresIn:    3600,
	}
	
	mockAuthService.On("RefreshToken", mock.Anything, "refresh_token_123").Return(loginResult, nil)
	
	requestBody, _ := json.Marshal(refreshRequest)
	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewBuffer(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	
	// Act
	err := authHandler.RefreshToken(c)
	
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	
	var response handler.LoginResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, loginResult.AccessToken, response.AccessToken)
	assert.Equal(t, loginResult.RefreshToken, response.RefreshToken)
	
	mockAuthService.AssertExpectations(t)
}

func TestAuthHandler_RefreshToken_InvalidToken(t *testing.T) {
	// Arrange
	e := echo.New()
	mockAuthService := new(MockAuthService)
	authHandler := handler.NewAuthHandler(mockAuthService)
	
	refreshRequest := handler.RefreshTokenRequest{
		RefreshToken: "invalid_refresh_token",
	}
	
	mockAuthService.On("RefreshToken", mock.Anything, "invalid_refresh_token").Return(nil, assert.AnError)
	
	requestBody, _ := json.Marshal(refreshRequest)
	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewBuffer(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	
	// Act
	err := authHandler.RefreshToken(c)
	
	// Assert
	assert.Error(t, err)
	
	mockAuthService.AssertExpectations(t)
}

func TestAuthHandler_Logout_Success(t *testing.T) {
	// Arrange
	e := echo.New()
	mockAuthService := new(MockAuthService)
	authHandler := handler.NewAuthHandler(mockAuthService)
	
	mockAuthService.On("Logout", mock.Anything, "access_token_123").Return(nil)
	
	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	req.Header.Set("Authorization", "Bearer access_token_123")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	
	// Act
	err := authHandler.Logout(c)
	
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	
	var response map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Successfully logged out", response["message"])
	
	mockAuthService.AssertExpectations(t)
}

func TestAuthHandler_GetProfile_Success(t *testing.T) {
	// Arrange
	e := echo.New()
	mockAuthService := new(MockAuthService)
	authHandler := handler.NewAuthHandler(mockAuthService)
	
	userID := uuid.New()
	companyID := uuid.New()
	
	claims := &service.JWTClaims{
		Subject:   userID.String(),
		CompanyID: companyID.String(),
		Role:      "engineer",
		Username:  "testuser",
		ExpiresAt: time.Now().Add(time.Hour),
	}
	
	mockAuthService.On("ValidateToken", mock.Anything, "access_token_123").Return(claims, nil)
	
	req := httptest.NewRequest(http.MethodGet, "/auth/profile", nil)
	req.Header.Set("Authorization", "Bearer access_token_123")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	
	// Act
	err := authHandler.GetProfile(c)
	
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	
	var response handler.ProfileResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, claims.Subject, response.ID)
	assert.Equal(t, claims.Username, response.Username)
	assert.Equal(t, claims.Role, response.Role)
	
	mockAuthService.AssertExpectations(t)
}