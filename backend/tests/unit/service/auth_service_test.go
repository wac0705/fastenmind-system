package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"github.com/fastenmind/fastener-api/internal/model"
	"github.com/fastenmind/fastener-api/internal/service"
)

// MockAccountRepository is a mock implementation of AccountRepository
type MockAccountRepository struct {
	mock.Mock
}

func (m *MockAccountRepository) GetByUsername(ctx context.Context, username string) (*model.Account, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Account), args.Error(1)
}

func (m *MockAccountRepository) GetByEmail(ctx context.Context, email string) (*model.Account, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Account), args.Error(1)
}

func (m *MockAccountRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Account, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Account), args.Error(1)
}

func (m *MockAccountRepository) Create(ctx context.Context, account *model.Account) error {
	args := m.Called(ctx, account)
	return args.Error(0)
}

func (m *MockAccountRepository) Update(ctx context.Context, account *model.Account) error {
	args := m.Called(ctx, account)
	return args.Error(0)
}

func (m *MockAccountRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID, loginTime time.Time) error {
	args := m.Called(ctx, id, loginTime)
	return args.Error(0)
}

func TestAuthService_Login_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockAccountRepository)
	authService := service.NewAuthService(mockRepo)
	
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	
	expectedAccount := &model.Account{
		ID:           uuid.New(),
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
		FullName:     "Test User",
		Role:         "engineer",
		IsActive:     true,
		CompanyID:    uuid.New(),
	}
	
	mockRepo.On("GetByUsername", mock.Anything, "testuser").Return(expectedAccount, nil)
	mockRepo.On("UpdateLastLogin", mock.Anything, expectedAccount.ID, mock.AnythingOfType("time.Time")).Return(nil)
	
	// Act
	result, err := authService.Login(context.Background(), "testuser", "password123")
	
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedAccount.ID, result.Account.ID)
	assert.Equal(t, expectedAccount.Username, result.Account.Username)
	assert.NotEmpty(t, result.AccessToken)
	assert.NotEmpty(t, result.RefreshToken)
	
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_InvalidUsername(t *testing.T) {
	// Arrange
	mockRepo := new(MockAccountRepository)
	authService := service.NewAuthService(mockRepo)
	
	mockRepo.On("GetByUsername", mock.Anything, "nonexistent").Return(nil, assert.AnError)
	
	// Act
	result, err := authService.Login(context.Background(), "nonexistent", "password123")
	
	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid credentials")
	
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	// Arrange
	mockRepo := new(MockAccountRepository)
	authService := service.NewAuthService(mockRepo)
	
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	
	expectedAccount := &model.Account{
		ID:           uuid.New(),
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
		FullName:     "Test User",
		Role:         "engineer",
		IsActive:     true,
		CompanyID:    uuid.New(),
	}
	
	mockRepo.On("GetByUsername", mock.Anything, "testuser").Return(expectedAccount, nil)
	
	// Act
	result, err := authService.Login(context.Background(), "testuser", "wrongpassword")
	
	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid credentials")
	
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_InactiveAccount(t *testing.T) {
	// Arrange
	mockRepo := new(MockAccountRepository)
	authService := service.NewAuthService(mockRepo)
	
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	
	expectedAccount := &model.Account{
		ID:           uuid.New(),
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
		FullName:     "Test User",
		Role:         "engineer",
		IsActive:     false, // Inactive account
		CompanyID:    uuid.New(),
	}
	
	mockRepo.On("GetByUsername", mock.Anything, "testuser").Return(expectedAccount, nil)
	
	// Act
	result, err := authService.Login(context.Background(), "testuser", "password123")
	
	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "account is disabled")
	
	mockRepo.AssertExpectations(t)
}

func TestAuthService_ValidateToken_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockAccountRepository)
	authService := service.NewAuthService(mockRepo)
	
	accountID := uuid.New()
	companyID := uuid.New()
	
	expectedAccount := &model.Account{
		ID:        accountID,
		Username:  "testuser",
		Email:     "test@example.com",
		FullName:  "Test User",  
		Role:      "engineer",
		IsActive:  true,
		CompanyID: companyID,
	}
	
	// Generate a valid token
	loginResult, _ := authService.Login(context.Background(), "testuser", "password123")
	token := loginResult.AccessToken
	
	mockRepo.On("GetByID", mock.Anything, accountID).Return(expectedAccount, nil)
	
	// Act
	claims, err := authService.ValidateToken(context.Background(), token)
	
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, accountID.String(), claims.Subject)
	assert.Equal(t, companyID.String(), claims.CompanyID)
	assert.Equal(t, "engineer", claims.Role)
	
	mockRepo.AssertExpectations(t)
}

func TestAuthService_RefreshToken_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockAccountRepository)
	authService := service.NewAuthService(mockRepo)
	
	accountID := uuid.New()
	companyID := uuid.New()
	
	expectedAccount := &model.Account{
		ID:        accountID,
		Username:  "testuser",
		Email:     "test@example.com",
		FullName:  "Test User",
		Role:      "engineer",
		IsActive:  true,
		CompanyID: companyID,
	}
	
	// Generate a valid refresh token
	loginResult, _ := authService.Login(context.Background(), "testuser", "password123")
	refreshToken := loginResult.RefreshToken
	
	mockRepo.On("GetByID", mock.Anything, accountID).Return(expectedAccount, nil)
	
	// Act
	result, err := authService.RefreshToken(context.Background(), refreshToken)
	
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.AccessToken)
	assert.NotEmpty(t, result.RefreshToken)
	assert.Equal(t, expectedAccount.ID, result.Account.ID)
	
	mockRepo.AssertExpectations(t)
}