package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/fastenmind/fastener-api/internal/domain/models"
	"github.com/fastenmind/fastener-api/internal/repository"
	"github.com/fastenmind/fastener-api/internal/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockInquiryRepository is a mock implementation of InquiryRepository
type MockInquiryRepository struct {
	mock.Mock
}

func (m *MockInquiryRepository) Create(ctx context.Context, inquiry *models.Inquiry) error {
	args := m.Called(ctx, inquiry)
	return args.Error(0)
}

func (m *MockInquiryRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Inquiry, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Inquiry), args.Error(1)
}

func (m *MockInquiryRepository) Update(ctx context.Context, inquiry *models.Inquiry) error {
	args := m.Called(ctx, inquiry)
	return args.Error(0)
}

func (m *MockInquiryRepository) List(ctx context.Context, params repository.ListInquiriesParams) ([]*models.Inquiry, int64, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*models.Inquiry), args.Get(1).(int64), args.Error(2)
}

func (m *MockInquiryRepository) GetByInquiryNo(ctx context.Context, inquiryNo string) (*models.Inquiry, error) {
	args := m.Called(ctx, inquiryNo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Inquiry), args.Error(1)
}

func (m *MockInquiryRepository) ExistsByInquiryNo(ctx context.Context, inquiryNo string) (bool, error) {
	args := m.Called(ctx, inquiryNo)
	return args.Bool(0), args.Error(1)
}

// Test suite
type InquiryServiceTestSuite struct {
	service  *service.InquiryService
	mockRepo *MockInquiryRepository
	ctx      context.Context
}

func setupTestSuite() *InquiryServiceTestSuite {
	mockRepo := new(MockInquiryRepository)
	svc := service.NewInquiryService(mockRepo, nil) // DB is not needed for unit tests
	
	return &InquiryServiceTestSuite{
		service:  svc,
		mockRepo: mockRepo,
		ctx:      context.Background(),
	}
}

func TestInquiryService_Create(t *testing.T) {
	suite := setupTestSuite()
	
	t.Run("successful creation", func(t *testing.T) {
		// Arrange
		req := service.CreateInquiryRequest{
			CompanyID:       uuid.New(),
			CustomerID:      uuid.New(),
			SalesID:         uuid.New(),
			ProductCategory: "Fasteners",
			ProductName:     "Hex Bolt",
			Quantity:        1000,
			Unit:            "PCS",
			RequiredDate:    time.Now().Add(30 * 24 * time.Hour),
			Incoterm:        "FOB",
			PaymentTerms:    "Net 30",
		}
		
		suite.mockRepo.On("ExistsByInquiryNo", suite.ctx, mock.AnythingOfType("string")).Return(false, nil)
		suite.mockRepo.On("Create", suite.ctx, mock.AnythingOfType("*models.Inquiry")).Return(nil)
		
		// Act
		result, err := suite.service.Create(suite.ctx, req)
		
		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, req.ProductName, result.ProductName)
		assert.Equal(t, req.Quantity, result.Quantity)
		assert.Equal(t, models.InquiryStatusPending, result.Status)
		
		suite.mockRepo.AssertExpectations(t)
	})
	
	t.Run("duplicate inquiry number", func(t *testing.T) {
		// Arrange
		req := service.CreateInquiryRequest{
			CompanyID:       uuid.New(),
			CustomerID:      uuid.New(),
			SalesID:         uuid.New(),
			ProductCategory: "Fasteners",
			ProductName:     "Hex Bolt",
			Quantity:        1000,
			Unit:            "PCS",
			RequiredDate:    time.Now().Add(30 * 24 * time.Hour),
			Incoterm:        "FOB",
		}
		
		suite.mockRepo.On("ExistsByInquiryNo", suite.ctx, mock.AnythingOfType("string")).Return(true, nil)
		
		// Act
		result, err := suite.service.Create(suite.ctx, req)
		
		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "inquiry number already exists")
		
		suite.mockRepo.AssertExpectations(t)
	})
}

func TestInquiryService_AssignEngineer(t *testing.T) {
	suite := setupTestSuite()
	
	t.Run("successful assignment", func(t *testing.T) {
		// Arrange
		inquiryID := uuid.New()
		engineerID := uuid.New()
		assignedBy := uuid.New()
		
		inquiry := &models.Inquiry{
			BaseModel: models.BaseModel{ID: inquiryID},
			Status:    models.InquiryStatusPending,
		}
		
		suite.mockRepo.On("GetByID", suite.ctx, inquiryID).Return(inquiry, nil)
		suite.mockRepo.On("Update", suite.ctx, mock.AnythingOfType("*models.Inquiry")).Return(nil)
		
		// Act
		result, err := suite.service.AssignEngineer(suite.ctx, inquiryID, engineerID, assignedBy, "High priority")
		
		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, engineerID, *result.AssignedEngineerID)
		assert.Equal(t, models.InquiryStatusAssigned, result.Status)
		assert.NotNil(t, result.AssignedAt)
		
		suite.mockRepo.AssertExpectations(t)
	})
	
	t.Run("inquiry not found", func(t *testing.T) {
		// Arrange
		inquiryID := uuid.New()
		engineerID := uuid.New()
		assignedBy := uuid.New()
		
		suite.mockRepo.On("GetByID", suite.ctx, inquiryID).Return(nil, gorm.ErrRecordNotFound)
		
		// Act
		result, err := suite.service.AssignEngineer(suite.ctx, inquiryID, engineerID, assignedBy, "")
		
		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, service.ErrNotFound, err)
		
		suite.mockRepo.AssertExpectations(t)
	})
	
	t.Run("already assigned", func(t *testing.T) {
		// Arrange
		inquiryID := uuid.New()
		engineerID := uuid.New()
		assignedBy := uuid.New()
		existingEngineerID := uuid.New()
		
		inquiry := &models.Inquiry{
			BaseModel:          models.BaseModel{ID: inquiryID},
			Status:             models.InquiryStatusAssigned,
			AssignedEngineerID: &existingEngineerID,
		}
		
		suite.mockRepo.On("GetByID", suite.ctx, inquiryID).Return(inquiry, nil)
		
		// Act
		result, err := suite.service.AssignEngineer(suite.ctx, inquiryID, engineerID, assignedBy, "")
		
		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "already assigned")
		
		suite.mockRepo.AssertExpectations(t)
	})
}

func TestInquiryService_List(t *testing.T) {
	suite := setupTestSuite()
	
	t.Run("successful list with pagination", func(t *testing.T) {
		// Arrange
		params := service.ListInquiriesParams{
			CompanyID: uuid.New(),
			Page:      2,
			PageSize:  10,
			Status:    "pending",
		}
		
		inquiries := []*models.Inquiry{
			{
				BaseModel:   models.BaseModel{ID: uuid.New()},
				InquiryNo:   "INQ-001",
				ProductName: "Hex Bolt",
				Status:      models.InquiryStatusPending,
			},
			{
				BaseModel:   models.BaseModel{ID: uuid.New()},
				InquiryNo:   "INQ-002",
				ProductName: "Flat Washer",
				Status:      models.InquiryStatusPending,
			},
		}
		
		expectedRepoParams := repository.ListInquiriesParams{
			CompanyID: params.CompanyID,
			Status:    &params.Status,
			Offset:    10, // (page-1) * pageSize
			Limit:     10,
		}
		
		suite.mockRepo.On("List", suite.ctx, expectedRepoParams).Return(inquiries, int64(25), nil)
		
		// Act
		result, err := suite.service.List(suite.ctx, params)
		
		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Items, 2)
		assert.Equal(t, int64(25), result.Total)
		assert.Equal(t, 2, result.Page)
		assert.Equal(t, 10, result.PageSize)
		assert.Equal(t, 3, result.TotalPages)
		
		suite.mockRepo.AssertExpectations(t)
	})
}

// Benchmark tests
func BenchmarkInquiryService_Create(b *testing.B) {
	suite := setupTestSuite()
	
	req := service.CreateInquiryRequest{
		CompanyID:       uuid.New(),
		CustomerID:      uuid.New(),
		SalesID:         uuid.New(),
		ProductCategory: "Fasteners",
		ProductName:     "Hex Bolt",
		Quantity:        1000,
		Unit:            "PCS",
		RequiredDate:    time.Now().Add(30 * 24 * time.Hour),
		Incoterm:        "FOB",
	}
	
	suite.mockRepo.On("ExistsByInquiryNo", mock.Anything, mock.AnythingOfType("string")).Return(false, nil)
	suite.mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Inquiry")).Return(nil)
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_, _ = suite.service.Create(context.Background(), req)
	}
}

// Table-driven tests
func TestInquiryService_ValidateInquiryData(t *testing.T) {
	tests := []struct {
		name    string
		req     service.CreateInquiryRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid inquiry",
			req: service.CreateInquiryRequest{
				CompanyID:       uuid.New(),
				CustomerID:      uuid.New(),
				SalesID:         uuid.New(),
				ProductCategory: "Fasteners",
				ProductName:     "Hex Bolt",
				Quantity:        1000,
				Unit:            "PCS",
				RequiredDate:    time.Now().Add(30 * 24 * time.Hour),
				Incoterm:        "FOB",
			},
			wantErr: false,
		},
		{
			name: "missing product name",
			req: service.CreateInquiryRequest{
				CompanyID:       uuid.New(),
				CustomerID:      uuid.New(),
				SalesID:         uuid.New(),
				ProductCategory: "Fasteners",
				ProductName:     "",
				Quantity:        1000,
				Unit:            "PCS",
				RequiredDate:    time.Now().Add(30 * 24 * time.Hour),
				Incoterm:        "FOB",
			},
			wantErr: true,
			errMsg:  "product name is required",
		},
		{
			name: "invalid quantity",
			req: service.CreateInquiryRequest{
				CompanyID:       uuid.New(),
				CustomerID:      uuid.New(),
				SalesID:         uuid.New(),
				ProductCategory: "Fasteners",
				ProductName:     "Hex Bolt",
				Quantity:        0,
				Unit:            "PCS",
				RequiredDate:    time.Now().Add(30 * 24 * time.Hour),
				Incoterm:        "FOB",
			},
			wantErr: true,
			errMsg:  "quantity must be greater than zero",
		},
		{
			name: "past required date",
			req: service.CreateInquiryRequest{
				CompanyID:       uuid.New(),
				CustomerID:      uuid.New(),
				SalesID:         uuid.New(),
				ProductCategory: "Fasteners",
				ProductName:     "Hex Bolt",
				Quantity:        1000,
				Unit:            "PCS",
				RequiredDate:    time.Now().Add(-24 * time.Hour),
				Incoterm:        "FOB",
			},
			wantErr: true,
			errMsg:  "required date must be in the future",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ValidateInquiryData(tt.req)
			
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}