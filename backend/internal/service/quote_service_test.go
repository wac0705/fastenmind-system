package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/fastenmind/fastener-api/internal/domain/models"
	"github.com/fastenmind/fastener-api/internal/repository"
	"github.com/fastenmind/fastener-api/internal/service"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockQuoteRepository is a mock implementation of QuoteRepository
type MockQuoteRepository struct {
	mock.Mock
}

func (m *MockQuoteRepository) Create(ctx context.Context, quote *models.Quote) error {
	args := m.Called(ctx, quote)
	return args.Error(0)
}

func (m *MockQuoteRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Quote, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Quote), args.Error(1)
}

func (m *MockQuoteRepository) Update(ctx context.Context, quote *models.Quote) error {
	args := m.Called(ctx, quote)
	return args.Error(0)
}

func (m *MockQuoteRepository) List(ctx context.Context, params repository.ListQuotesParams) ([]*models.Quote, int64, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*models.Quote), args.Get(1).(int64), args.Error(2)
}

func (m *MockQuoteRepository) GetByQuoteNo(ctx context.Context, quoteNo string) (*models.Quote, error) {
	args := m.Called(ctx, quoteNo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Quote), args.Error(1)
}

func (m *MockQuoteRepository) ExistsByQuoteNo(ctx context.Context, quoteNo string) (bool, error) {
	args := m.Called(ctx, quoteNo)
	return args.Bool(0), args.Error(1)
}

func (m *MockQuoteRepository) CreateWithItems(ctx context.Context, quote *models.Quote, items []models.QuoteItem) error {
	args := m.Called(ctx, quote, items)
	return args.Error(0)
}

// Test suite
type QuoteServiceTestSuite struct {
	service  *service.QuoteService
	mockRepo *MockQuoteRepository
	ctx      context.Context
}

func setupQuoteTestSuite() *QuoteServiceTestSuite {
	mockRepo := new(MockQuoteRepository)
	svc := service.NewQuoteService(mockRepo, nil) // DB is not needed for unit tests
	
	return &QuoteServiceTestSuite{
		service:  svc,
		mockRepo: mockRepo,
		ctx:      context.Background(),
	}
}

func TestQuoteService_Create(t *testing.T) {
	suite := setupQuoteTestSuite()
	
	t.Run("successful creation with items", func(t *testing.T) {
		// Arrange
		items := []service.QuoteItemRequest{
			{
				ProductName:    "Hex Bolt M10x50",
				Specification:  "Grade 8.8, Zinc Plated",
				Quantity:       1000,
				Unit:           "PCS",
				UnitPrice:      decimal.NewFromFloat(0.25),
				MaterialCost:   decimal.NewFromFloat(0.15),
				ProcessingCost: decimal.NewFromFloat(0.05),
				Currency:       "USD",
			},
			{
				ProductName:    "Flat Washer M10",
				Specification:  "DIN 125, Zinc Plated",
				Quantity:       1000,
				Unit:           "PCS",
				UnitPrice:      decimal.NewFromFloat(0.05),
				MaterialCost:   decimal.NewFromFloat(0.03),
				ProcessingCost: decimal.NewFromFloat(0.01),
				Currency:       "USD",
			},
		}
		
		req := service.CreateQuoteRequest{
			CompanyID:     uuid.New(),
			InquiryID:     uuid.New(),
			CustomerID:    uuid.New(),
			PreparedBy:    uuid.New(),
			ValidUntil:    time.Now().Add(30 * 24 * time.Hour),
			PaymentTerms:  "Net 30",
			DeliveryTerms: "FOB Shanghai",
			LeadTime:      "4-6 weeks",
			Items:         items,
		}
		
		suite.mockRepo.On("ExistsByQuoteNo", suite.ctx, mock.AnythingOfType("string")).Return(false, nil)
		suite.mockRepo.On("CreateWithItems", suite.ctx, mock.AnythingOfType("*models.Quote"), mock.AnythingOfType("[]models.QuoteItem")).Return(nil)
		
		// Act
		result, err := suite.service.Create(suite.ctx, req)
		
		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, req.CustomerID, result.CustomerID)
		assert.Equal(t, models.QuoteStatusDraft, result.Status)
		assert.Len(t, result.Items, 2)
		assert.Equal(t, decimal.NewFromFloat(300), result.TotalAmount) // (0.25*1000) + (0.05*1000)
		
		suite.mockRepo.AssertExpectations(t)
	})
	
	t.Run("duplicate quote number", func(t *testing.T) {
		// Arrange
		req := service.CreateQuoteRequest{
			CompanyID:  uuid.New(),
			InquiryID:  uuid.New(),
			CustomerID: uuid.New(),
			PreparedBy: uuid.New(),
			ValidUntil: time.Now().Add(30 * 24 * time.Hour),
			Items:      []service.QuoteItemRequest{},
		}
		
		suite.mockRepo.On("ExistsByQuoteNo", suite.ctx, mock.AnythingOfType("string")).Return(true, nil)
		
		// Act
		result, err := suite.service.Create(suite.ctx, req)
		
		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "quote number already exists")
		
		suite.mockRepo.AssertExpectations(t)
	})
}

func TestQuoteService_SubmitForApproval(t *testing.T) {
	suite := setupQuoteTestSuite()
	
	t.Run("successful submission", func(t *testing.T) {
		// Arrange
		quoteID := uuid.New()
		submittedBy := uuid.New()
		
		quote := &models.Quote{
			BaseModel:   models.BaseModel{ID: quoteID},
			Status:      models.QuoteStatusDraft,
			TotalAmount: decimal.NewFromFloat(1000),
		}
		
		suite.mockRepo.On("GetByID", suite.ctx, quoteID).Return(quote, nil)
		suite.mockRepo.On("Update", suite.ctx, mock.AnythingOfType("*models.Quote")).Return(nil)
		
		// Act
		result, err := suite.service.SubmitForApproval(suite.ctx, quoteID, submittedBy)
		
		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, models.QuoteStatusPendingApproval, result.Status)
		assert.NotNil(t, result.SubmittedAt)
		assert.Equal(t, submittedBy, *result.SubmittedBy)
		
		suite.mockRepo.AssertExpectations(t)
	})
	
	t.Run("quote not found", func(t *testing.T) {
		// Arrange
		quoteID := uuid.New()
		submittedBy := uuid.New()
		
		suite.mockRepo.On("GetByID", suite.ctx, quoteID).Return(nil, gorm.ErrRecordNotFound)
		
		// Act
		result, err := suite.service.SubmitForApproval(suite.ctx, quoteID, submittedBy)
		
		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, service.ErrNotFound, err)
		
		suite.mockRepo.AssertExpectations(t)
	})
	
	t.Run("invalid status transition", func(t *testing.T) {
		// Arrange
		quoteID := uuid.New()
		submittedBy := uuid.New()
		
		quote := &models.Quote{
			BaseModel: models.BaseModel{ID: quoteID},
			Status:    models.QuoteStatusApproved, // Already approved
		}
		
		suite.mockRepo.On("GetByID", suite.ctx, quoteID).Return(quote, nil)
		
		// Act
		result, err := suite.service.SubmitForApproval(suite.ctx, quoteID, submittedBy)
		
		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "can only submit draft quotes")
		
		suite.mockRepo.AssertExpectations(t)
	})
}

func TestQuoteService_Approve(t *testing.T) {
	suite := setupQuoteTestSuite()
	
	t.Run("successful approval", func(t *testing.T) {
		// Arrange
		quoteID := uuid.New()
		approvedBy := uuid.New()
		
		quote := &models.Quote{
			BaseModel: models.BaseModel{ID: quoteID},
			Status:    models.QuoteStatusPendingApproval,
		}
		
		suite.mockRepo.On("GetByID", suite.ctx, quoteID).Return(quote, nil)
		suite.mockRepo.On("Update", suite.ctx, mock.AnythingOfType("*models.Quote")).Return(nil)
		
		// Act
		result, err := suite.service.Approve(suite.ctx, quoteID, approvedBy, "Approved for customer")
		
		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, models.QuoteStatusApproved, result.Status)
		assert.NotNil(t, result.ApprovedAt)
		assert.Equal(t, approvedBy, *result.ApprovedBy)
		
		suite.mockRepo.AssertExpectations(t)
	})
}

func TestQuoteService_List(t *testing.T) {
	suite := setupQuoteTestSuite()
	
	t.Run("successful list with filters", func(t *testing.T) {
		// Arrange
		params := service.ListQuotesParams{
			CompanyID:  uuid.New(),
			CustomerID: uuid.Nil,
			Status:     "draft",
			Page:       1,
			PageSize:   20,
		}
		
		quotes := []*models.Quote{
			{
				BaseModel:   models.BaseModel{ID: uuid.New()},
				QuoteNo:     "Q-2024-001",
				TotalAmount: decimal.NewFromFloat(1000),
				Status:      models.QuoteStatusDraft,
			},
			{
				BaseModel:   models.BaseModel{ID: uuid.New()},
				QuoteNo:     "Q-2024-002",
				TotalAmount: decimal.NewFromFloat(2000),
				Status:      models.QuoteStatusDraft,
			},
		}
		
		expectedRepoParams := repository.ListQuotesParams{
			CompanyID: params.CompanyID,
			Status:    &params.Status,
			Offset:    0,
			Limit:     20,
		}
		
		suite.mockRepo.On("List", suite.ctx, expectedRepoParams).Return(quotes, int64(2), nil)
		
		// Act
		result, err := suite.service.List(suite.ctx, params)
		
		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Items, 2)
		assert.Equal(t, int64(2), result.Total)
		assert.Equal(t, 1, result.Page)
		assert.Equal(t, 20, result.PageSize)
		assert.Equal(t, 1, result.TotalPages)
		
		suite.mockRepo.AssertExpectations(t)
	})
}

// Benchmark tests
func BenchmarkQuoteService_Create(b *testing.B) {
	suite := setupQuoteTestSuite()
	
	req := service.CreateQuoteRequest{
		CompanyID:  uuid.New(),
		InquiryID:  uuid.New(),
		CustomerID: uuid.New(),
		PreparedBy: uuid.New(),
		ValidUntil: time.Now().Add(30 * 24 * time.Hour),
		Items: []service.QuoteItemRequest{
			{
				ProductName: "Test Product",
				Quantity:    100,
				UnitPrice:   decimal.NewFromFloat(10),
				Currency:    "USD",
			},
		},
	}
	
	suite.mockRepo.On("ExistsByQuoteNo", mock.Anything, mock.AnythingOfType("string")).Return(false, nil)
	suite.mockRepo.On("CreateWithItems", mock.Anything, mock.AnythingOfType("*models.Quote"), mock.AnythingOfType("[]models.QuoteItem")).Return(nil)
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_, _ = suite.service.Create(context.Background(), req)
	}
}

// Table-driven tests
func TestQuoteService_ValidateQuoteData(t *testing.T) {
	tests := []struct {
		name    string
		req     service.CreateQuoteRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid quote",
			req: service.CreateQuoteRequest{
				CompanyID:     uuid.New(),
				InquiryID:     uuid.New(),
				CustomerID:    uuid.New(),
				PreparedBy:    uuid.New(),
				ValidUntil:    time.Now().Add(30 * 24 * time.Hour),
				PaymentTerms:  "Net 30",
				DeliveryTerms: "FOB",
				Items: []service.QuoteItemRequest{
					{
						ProductName: "Product 1",
						Quantity:    100,
						UnitPrice:   decimal.NewFromFloat(10),
						Currency:    "USD",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "no items",
			req: service.CreateQuoteRequest{
				CompanyID:  uuid.New(),
				InquiryID:  uuid.New(),
				CustomerID: uuid.New(),
				PreparedBy: uuid.New(),
				ValidUntil: time.Now().Add(30 * 24 * time.Hour),
				Items:      []service.QuoteItemRequest{},
			},
			wantErr: true,
			errMsg:  "at least one item is required",
		},
		{
			name: "past valid until date",
			req: service.CreateQuoteRequest{
				CompanyID:  uuid.New(),
				InquiryID:  uuid.New(),
				CustomerID: uuid.New(),
				PreparedBy: uuid.New(),
				ValidUntil: time.Now().Add(-24 * time.Hour),
				Items: []service.QuoteItemRequest{
					{
						ProductName: "Product 1",
						Quantity:    100,
						UnitPrice:   decimal.NewFromFloat(10),
						Currency:    "USD",
					},
				},
			},
			wantErr: true,
			errMsg:  "valid until date must be in the future",
		},
		{
			name: "invalid item quantity",
			req: service.CreateQuoteRequest{
				CompanyID:  uuid.New(),
				InquiryID:  uuid.New(),
				CustomerID: uuid.New(),
				PreparedBy: uuid.New(),
				ValidUntil: time.Now().Add(30 * 24 * time.Hour),
				Items: []service.QuoteItemRequest{
					{
						ProductName: "Product 1",
						Quantity:    0,
						UnitPrice:   decimal.NewFromFloat(10),
						Currency:    "USD",
					},
				},
			},
			wantErr: true,
			errMsg:  "item quantity must be greater than zero",
		},
		{
			name: "negative unit price",
			req: service.CreateQuoteRequest{
				CompanyID:  uuid.New(),
				InquiryID:  uuid.New(),
				CustomerID: uuid.New(),
				PreparedBy: uuid.New(),
				ValidUntil: time.Now().Add(30 * 24 * time.Hour),
				Items: []service.QuoteItemRequest{
					{
						ProductName: "Product 1",
						Quantity:    100,
						UnitPrice:   decimal.NewFromFloat(-10),
						Currency:    "USD",
					},
				},
			},
			wantErr: true,
			errMsg:  "unit price must be positive",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ValidateQuoteData(tt.req)
			
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

// Test quote revision functionality
func TestQuoteService_CreateRevision(t *testing.T) {
	suite := setupQuoteTestSuite()
	
	t.Run("successful revision creation", func(t *testing.T) {
		// Arrange
		originalID := uuid.New()
		createdBy := uuid.New()
		
		original := &models.Quote{
			BaseModel:     models.BaseModel{ID: originalID},
			QuoteNo:       "Q-2024-001",
			Status:        models.QuoteStatusApproved,
			RevisionNo:    1,
			CustomerID:    uuid.New(),
			InquiryID:     uuid.New(),
			PaymentTerms:  "Net 30",
			DeliveryTerms: "FOB",
		}
		
		suite.mockRepo.On("GetByID", suite.ctx, originalID).Return(original, nil)
		suite.mockRepo.On("ExistsByQuoteNo", suite.ctx, "Q-2024-001-R2").Return(false, nil)
		suite.mockRepo.On("CreateWithItems", suite.ctx, mock.AnythingOfType("*models.Quote"), mock.AnythingOfType("[]models.QuoteItem")).Return(nil)
		
		// Act
		result, err := suite.service.CreateRevision(suite.ctx, originalID, createdBy, "Customer requested changes")
		
		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Q-2024-001-R2", result.QuoteNo)
		assert.Equal(t, 2, result.RevisionNo)
		assert.Equal(t, models.QuoteStatusDraft, result.Status)
		assert.Equal(t, originalID, *result.OriginalQuoteID)
		
		suite.mockRepo.AssertExpectations(t)
	})
}