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

// MockQuoteRepository is a mock implementation of QuoteRepository
type MockQuoteRepository struct {
	mock.Mock
}

func (m *MockQuoteRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Quote, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Quote), args.Error(1)
}

func (m *MockQuoteRepository) Create(ctx context.Context, quote *model.Quote) error {
	args := m.Called(ctx, quote)
	return args.Error(0)
}

func (m *MockQuoteRepository) Update(ctx context.Context, quote *model.Quote) error {
	args := m.Called(ctx, quote)
	return args.Error(0)
}

func (m *MockQuoteRepository) List(ctx context.Context, filter map[string]interface{}) ([]*model.Quote, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Quote), args.Error(1)
}

func (m *MockQuoteRepository) GetByInquiryID(ctx context.Context, inquiryID uuid.UUID) (*model.Quote, error) {
	args := m.Called(ctx, inquiryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Quote), args.Error(1)
}

func (m *MockQuoteRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockInquiryRepository is a mock implementation for inquiry operations
type MockInquiryRepository struct {
	mock.Mock
}

func (m *MockInquiryRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Inquiry, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Inquiry), args.Error(1)
}

func TestQuoteService_CreateQuote_Success(t *testing.T) {
	// Arrange
	mockQuoteRepo := new(MockQuoteRepository)
	mockInquiryRepo := new(MockInquiryRepository)
	quoteService := service.NewQuoteService(mockQuoteRepo, mockInquiryRepo)
	
	inquiryID := uuid.New()
	customerID := uuid.New()
	engineerID := uuid.New()
	companyID := uuid.New()
	
	inquiry := &model.Inquiry{
		ID:            inquiryID,
		InquiryNo:     "INQ-2024-001",
		CustomerID:    customerID,
		ProductName:   "Hex Bolt M8x20",
		ProductCategory: "Bolts",
		Quantity:      10000,
		Unit:          "pcs",
		RequiredDate:  time.Now().AddDate(0, 1, 0),
		Status:        "assigned",
	}
	
	createRequest := &service.CreateQuoteRequest{
		InquiryID:     inquiryID,
		EngineerID:    engineerID,
		MaterialCost:  1000.00,
		ProcessCost:   500.00,
		SurfaceCost:   200.00,
		PackagingCost: 100.00,
		ShippingCost:  300.00,
		OverheadRate:  0.15,
		ProfitRate:    0.20,
		Currency:      "USD",
		ValidUntil:    time.Now().AddDate(0, 1, 0),
		DeliveryDays:  30,
		PaymentTerms:  "T/T 30 days",
		Notes:         "Test quote",
	}
	
	mockInquiryRepo.On("GetByID", mock.Anything, inquiryID).Return(inquiry, nil)
	mockQuoteRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.Quote")).Return(nil)
	
	// Act
	result, err := quoteService.CreateQuote(context.Background(), createRequest, companyID, engineerID)
	
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, inquiryID, result.InquiryID)
	assert.Equal(t, engineerID, result.EngineerID)
	assert.Equal(t, "draft", result.Status)
	assert.Equal(t, createRequest.Currency, result.Currency)
	assert.True(t, result.TotalCost > 0)
	assert.True(t, result.UnitPrice > 0)
	
	mockInquiryRepo.AssertExpectations(t)
	mockQuoteRepo.AssertExpectations(t)
}

func TestQuoteService_CreateQuote_InquiryNotFound(t *testing.T) {
	// Arrange
	mockQuoteRepo := new(MockQuoteRepository)
	mockInquiryRepo := new(MockInquiryRepository)
	quoteService := service.NewQuoteService(mockQuoteRepo, mockInquiryRepo)
	
	inquiryID := uuid.New()
	engineerID := uuid.New()
	companyID := uuid.New()
	
	createRequest := &service.CreateQuoteRequest{
		InquiryID:     inquiryID,
		EngineerID:    engineerID,
		MaterialCost:  1000.00,
		ProcessCost:   500.00,
	}
	
	mockInquiryRepo.On("GetByID", mock.Anything, inquiryID).Return(nil, assert.AnError)
	
	// Act
	result, err := quoteService.CreateQuote(context.Background(), createRequest, companyID, engineerID)
	
	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "inquiry not found")
	
	mockInquiryRepo.AssertExpectations(t)
}

func TestQuoteService_GetQuote_Success(t *testing.T) {
	// Arrange
	mockQuoteRepo := new(MockQuoteRepository)
	mockInquiryRepo := new(MockInquiryRepository)
	quoteService := service.NewQuoteService(mockQuoteRepo, mockInquiryRepo)
	
	quoteID := uuid.New()
	inquiryID := uuid.New()
	
	expectedQuote := &model.Quote{
		ID:            quoteID,
		QuoteNo:       "QUO-2024-001",
		InquiryID:     inquiryID,
		Status:        "draft",
		MaterialCost:  1000.00,
		ProcessCost:   500.00,
		TotalCost:     2000.00,
		UnitPrice:     0.20,
		Currency:      "USD",
		CreatedAt:     time.Now(),
	}
	
	mockQuoteRepo.On("GetByID", mock.Anything, quoteID).Return(expectedQuote, nil)
	
	// Act
	result, err := quoteService.GetQuote(context.Background(), quoteID)
	
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedQuote.ID, result.ID)
	assert.Equal(t, expectedQuote.QuoteNo, result.QuoteNo)
	assert.Equal(t, expectedQuote.Status, result.Status)
	
	mockQuoteRepo.AssertExpectations(t)
}

func TestQuoteService_UpdateQuoteStatus_Success(t *testing.T) {
	// Arrange
	mockQuoteRepo := new(MockQuoteRepository)
	mockInquiryRepo := new(MockInquiryRepository)
	quoteService := service.NewQuoteService(mockQuoteRepo, mockInquiryRepo)
	
	quoteID := uuid.New()
	
	existingQuote := &model.Quote{
		ID:        quoteID,
		QuoteNo:   "QUO-2024-001",
		Status:    "draft",
		CreatedAt: time.Now(),
	}
	
	mockQuoteRepo.On("GetByID", mock.Anything, quoteID).Return(existingQuote, nil)
	mockQuoteRepo.On("Update", mock.Anything, mock.AnythingOfType("*model.Quote")).Return(nil)
	
	// Act
	err := quoteService.UpdateQuoteStatus(context.Background(), quoteID, "pending_approval", "Ready for review")
	
	// Assert
	assert.NoError(t, err)
	
	mockQuoteRepo.AssertExpectations(t)
}

func TestQuoteService_CalculateTotalCost(t *testing.T) {
	// Arrange
	mockQuoteRepo := new(MockQuoteRepository)
	mockInquiryRepo := new(MockInquiryRepository)
	quoteService := service.NewQuoteService(mockQuoteRepo, mockInquiryRepo)
	
	costs := service.QuoteCosts{
		MaterialCost:  1000.00,
		ProcessCost:   500.00,
		SurfaceCost:   200.00,
		HeatTreatCost: 100.00,
		PackagingCost: 50.00,
		ShippingCost:  150.00,
		TariffCost:    100.00,
	}
	
	overheadRate := 0.15
	profitRate := 0.20
	
	// Act
	totalCost, unitPrice := quoteService.CalculateTotalCost(costs, overheadRate, profitRate, 10000)
	
	// Assert
	directCost := 1000.00 + 500.00 + 200.00 + 100.00 + 50.00 + 150.00 + 100.00 // 2100
	expectedOverhead := directCost * overheadRate                                  // 315
	subtotal := directCost + expectedOverhead                                      // 2415
	expectedTotalCost := subtotal * (1 + profitRate)                              // 2898
	expectedUnitPrice := expectedTotalCost / 10000                                // 0.2898
	
	assert.Equal(t, expectedTotalCost, totalCost)
	assert.Equal(t, expectedUnitPrice, unitPrice)
}

func TestQuoteService_ListQuotes_Success(t *testing.T) {
	// Arrange
	mockQuoteRepo := new(MockQuoteRepository)
	mockInquiryRepo := new(MockInquiryRepository)
	quoteService := service.NewQuoteService(mockQuoteRepo, mockInquiryRepo)
	
	companyID := uuid.New()
	
	expectedQuotes := []*model.Quote{
		{
			ID:        uuid.New(),
			QuoteNo:   "QUO-2024-001",
			Status:    "draft",
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			QuoteNo:   "QUO-2024-002", 
			Status:    "approved",
			CreatedAt: time.Now(),
		},
	}
	
	filter := map[string]interface{}{
		"company_id": companyID,
		"status":     "draft",
	}
	
	mockQuoteRepo.On("List", mock.Anything, filter).Return(expectedQuotes, nil)
	
	// Act
	result, err := quoteService.ListQuotes(context.Background(), filter)
	
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, expectedQuotes[0].QuoteNo, result[0].QuoteNo)
	assert.Equal(t, expectedQuotes[1].QuoteNo, result[1].QuoteNo)
	
	mockQuoteRepo.AssertExpectations(t)
}