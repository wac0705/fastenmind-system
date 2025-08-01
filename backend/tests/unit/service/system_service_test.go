package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/fastenmind/fastener-api/internal/model"
	"github.com/fastenmind/fastener-api/internal/service"
)

// MockSystemRepository is a mock implementation of SystemRepository
type MockSystemRepository struct {
	mock.Mock
}

func (m *MockSystemRepository) CreateUser(ctx context.Context, user *model.Account) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockSystemRepository) GetUser(ctx context.Context, id uuid.UUID) (*model.Account, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Account), args.Error(1)
}

func (m *MockSystemRepository) UpdateUser(ctx context.Context, user *model.Account) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockSystemRepository) ListUsers(ctx context.Context, companyID uuid.UUID) ([]*model.Account, error) {
	args := m.Called(ctx, companyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Account), args.Error(1)
}

func (m *MockSystemRepository) GetSystemSettings(ctx context.Context, companyID uuid.UUID) (*model.SystemSettings, error) {
	args := m.Called(ctx, companyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.SystemSettings), args.Error(1)
}

func (m *MockSystemRepository) UpdateSystemSettings(ctx context.Context, settings *model.SystemSettings) error {
	args := m.Called(ctx, settings)
	return args.Error(0)
}

func TestSystemService_CreateUser_Success(t *testing.T) {
	// Arrange
	mockSystemRepo := new(MockSystemRepository)
	systemService := service.NewSystemService(mockSystemRepo)
	
	companyID := uuid.New()
	
	createRequest := &service.CreateUserRequest{
		Username:    "newuser",
		Email:       "newuser@test.com",
		Password:    "password123",
		FullName:    "New User",
		Role:        "engineer",
		PhoneNumber: "123-456-7890",
	}
	
	mockSystemRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*model.Account")).Return(nil)
	
	// Act
	result, err := systemService.CreateUser(context.Background(), createRequest, companyID)
	
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "newuser", result.Username)
	assert.Equal(t, "newuser@test.com", result.Email)
	assert.Equal(t, "New User", result.FullName)
	assert.Equal(t, "engineer", result.Role)
	assert.True(t, result.IsActive)
	assert.NotEmpty(t, result.PasswordHash)
	
	mockSystemRepo.AssertExpectations(t)
}

func TestSystemService_UpdateUserRole_Success(t *testing.T) {
	// Arrange
	mockSystemRepo := new(MockSystemRepository)
	systemService := service.NewSystemService(mockSystemRepo)
	
	userID := uuid.New()
	
	existingUser := &model.Account{
		ID:       userID,
		Username: "testuser",
		Email:    "test@example.com",
		FullName: "Test User",
		Role:     "engineer",
		IsActive: true,
	}
	
	mockSystemRepo.On("GetUser", mock.Anything, userID).Return(existingUser, nil)
	mockSystemRepo.On("UpdateUser", mock.Anything, mock.AnythingOfType("*model.Account")).Return(nil)
	
	// Act
	err := systemService.UpdateUserRole(context.Background(), userID, "manager")
	
	// Assert
	assert.NoError(t, err)
	
	mockSystemRepo.AssertExpectations(t)
}

func TestSystemService_ActivateUser_Success(t *testing.T) {
	// Arrange
	mockSystemRepo := new(MockSystemRepository)
	systemService := service.NewSystemService(mockSystemRepo)
	
	userID := uuid.New()
	
	existingUser := &model.Account{
		ID:       userID,
		Username: "testuser",
		IsActive: false,
	}
	
	mockSystemRepo.On("GetUser", mock.Anything, userID).Return(existingUser, nil)
	mockSystemRepo.On("UpdateUser", mock.Anything, mock.AnythingOfType("*model.Account")).Return(nil)
	
	// Act
	err := systemService.ActivateUser(context.Background(), userID)
	
	// Assert
	assert.NoError(t, err)
	
	mockSystemRepo.AssertExpectations(t)
}

func TestSystemService_DeactivateUser_Success(t *testing.T) {
	// Arrange
	mockSystemRepo := new(MockSystemRepository)
	systemService := service.NewSystemService(mockSystemRepo)
	
	userID := uuid.New()
	
	existingUser := &model.Account{
		ID:       userID,
		Username: "testuser",
		IsActive: true,
	}
	
	mockSystemRepo.On("GetUser", mock.Anything, userID).Return(existingUser, nil)
	mockSystemRepo.On("UpdateUser", mock.Anything, mock.AnythingOfType("*model.Account")).Return(nil)
	
	// Act
	err := systemService.DeactivateUser(context.Background(), userID)
	
	// Assert
	assert.NoError(t, err)
	
	mockSystemRepo.AssertExpectations(t)
}

func TestSystemService_ListUsers_Success(t *testing.T) {
	// Arrange
	mockSystemRepo := new(MockSystemRepository)
	systemService := service.NewSystemService(mockSystemRepo)
	
	companyID := uuid.New()
	
	expectedUsers := []*model.Account{
		{
			ID:       uuid.New(),
			Username: "user1",
			Email:    "user1@test.com",
			Role:     "engineer",
			IsActive: true,
		},
		{
			ID:       uuid.New(),
			Username: "user2",
			Email:    "user2@test.com",
			Role:     "manager",
			IsActive: true,
		},
	}
	
	mockSystemRepo.On("ListUsers", mock.Anything, companyID).Return(expectedUsers, nil)
	
	// Act
	result, err := systemService.ListUsers(context.Background(), companyID)
	
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, "user1", result[0].Username)
	assert.Equal(t, "user2", result[1].Username)
	
	mockSystemRepo.AssertExpectations(t)
}

func TestSystemService_UpdateSystemSettings_Success(t *testing.T) {
	// Arrange
	mockSystemRepo := new(MockSystemRepository)
	systemService := service.NewSystemService(mockSystemRepo)
	
	companyID := uuid.New()
	
	updateRequest := &service.UpdateSystemSettingsRequest{
		CompanyName:        "Updated Company Name",
		Currency:           "EUR",
		Language:           "zh-TW",
		Timezone:           "Asia/Taipei",
		EmailNotifications: true,
		SMSNotifications:   false,
		TaxRate:            0.08,
	}
	
	existingSettings := &model.SystemSettings{
		CompanyID: companyID,
		Currency:  "USD",
		Language:  "en",
		Timezone:  "UTC",
		TaxRate:   0.05,
	}
	
	mockSystemRepo.On("GetSystemSettings", mock.Anything, companyID).Return(existingSettings, nil)
	mockSystemRepo.On("UpdateSystemSettings", mock.Anything, mock.AnythingOfType("*model.SystemSettings")).Return(nil)
	
	// Act
	result, err := systemService.UpdateSystemSettings(context.Background(), companyID, updateRequest)
	
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Updated Company Name", result.CompanyName)
	assert.Equal(t, "EUR", result.Currency)
	assert.Equal(t, "zh-TW", result.Language)
	assert.Equal(t, 0.08, result.TaxRate)
	
	mockSystemRepo.AssertExpectations(t)
}

func TestSystemService_ValidateUserRole_Success(t *testing.T) {
	// Arrange
	mockSystemRepo := new(MockSystemRepository)
	systemService := service.NewSystemService(mockSystemRepo)
	
	validRoles := []string{"admin", "manager", "engineer", "sales", "finance", "viewer"}
	invalidRoles := []string{"invalid", "test", "unknown"}
	
	// Act & Assert - Valid roles
	for _, role := range validRoles {
		isValid := systemService.ValidateUserRole(role)
		assert.True(t, isValid, "Role '%s' should be valid", role)
	}
	
	// Act & Assert - Invalid roles
	for _, role := range invalidRoles {
		isValid := systemService.ValidateUserRole(role)
		assert.False(t, isValid, "Role '%s' should be invalid", role)
	}
}

func TestSystemService_GenerateUserPassword_Success(t *testing.T) {
	// Arrange
	mockSystemRepo := new(MockSystemRepository)
	systemService := service.NewSystemService(mockSystemRepo)
	
	// Act
	password := systemService.GenerateUserPassword()
	
	// Assert
	assert.NotEmpty(t, password)
	assert.GreaterOrEqual(t, len(password), 8) // Minimum password length
	assert.LessOrEqual(t, len(password), 16)   // Maximum password length
}

func TestSystemService_HashPassword_Success(t *testing.T) {
	// Arrange
	mockSystemRepo := new(MockSystemRepository)
	systemService := service.NewSystemService(mockSystemRepo)
	
	plainPassword := "password123"
	
	// Act
	hashedPassword, err := systemService.HashPassword(plainPassword)
	
	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)
	assert.NotEqual(t, plainPassword, hashedPassword)
	assert.Contains(t, hashedPassword, "$2a$") // bcrypt prefix
}

func TestSystemService_ValidatePassword_Success(t *testing.T) {
	// Arrange
	mockSystemRepo := new(MockSystemRepository)
	systemService := service.NewSystemService(mockSystemRepo)
	
	plainPassword := "password123"
	hashedPassword, _ := systemService.HashPassword(plainPassword)
	
	// Act
	isValid := systemService.ValidatePassword(plainPassword, hashedPassword)
	
	// Assert
	assert.True(t, isValid)
}

func TestSystemService_ValidatePassword_Invalid(t *testing.T) {
	// Arrange
	mockSystemRepo := new(MockSystemRepository)
	systemService := service.NewSystemService(mockSystemRepo)
	
	plainPassword := "password123"
	wrongPassword := "wrongpassword"
	hashedPassword, _ := systemService.HashPassword(plainPassword)
	
	// Act
	isValid := systemService.ValidatePassword(wrongPassword, hashedPassword)
	
	// Assert
	assert.False(t, isValid)
}

func TestSystemService_GetSystemHealth_Success(t *testing.T) {
	// Arrange
	mockSystemRepo := new(MockSystemRepository)
	systemService := service.NewSystemService(mockSystemRepo)
	
	// Act
	health := systemService.GetSystemHealth(context.Background())
	
	// Assert
	assert.NotNil(t, health)
	assert.NotEmpty(t, health.Status)
	assert.NotZero(t, health.Timestamp)
	assert.NotNil(t, health.Database)
	assert.NotNil(t, health.Memory)
	assert.NotNil(t, health.Disk)
}

func TestSystemService_GetAuditLogs_Success(t *testing.T) {
	// Arrange
	mockSystemRepo := new(MockSystemRepository)
	systemService := service.NewSystemService(mockSystemRepo)
	
	companyID := uuid.New()
	filter := service.AuditLogFilter{
		StartDate:  time.Now().AddDate(0, -1, 0),
		EndDate:    time.Now(),
		UserID:     nil,
		Action:     "",
		Page:       1,
		Limit:      50,
	}
	
	expectedLogs := []*service.AuditLog{
		{
			ID:          uuid.New(),
			UserID:      uuid.New(),
			Action:      "user.login",
			Resource:    "auth",
			IPAddress:   "192.168.1.1",
			UserAgent:   "Test Browser",
			Timestamp:   time.Now(),
			Details:     map[string]interface{}{"success": true},
		},
	}
	
	mockSystemRepo.On("GetAuditLogs", mock.Anything, companyID, filter).Return(expectedLogs, 1, nil)
	
	// Act
	result, total, err := systemService.GetAuditLogs(context.Background(), companyID, filter)
	
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	assert.Equal(t, 1, total)
	assert.Equal(t, "user.login", result[0].Action)
	
	mockSystemRepo.AssertExpectations(t)
}