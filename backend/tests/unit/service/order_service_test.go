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

// MockOrderRepository is a mock implementation of OrderRepository
type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

func (m *MockOrderRepository) Create(ctx context.Context, order *model.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderRepository) Update(ctx context.Context, order *model.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderRepository) List(ctx context.Context, filter map[string]interface{}) ([]*model.Order, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Order), args.Error(1)
}

func (m *MockOrderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockOrderRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string, notes string) error {
	args := m.Called(ctx, id, status, notes)
	return args.Error(0)
}

func (m *MockOrderRepository) GetItems(ctx context.Context, orderID uuid.UUID) ([]*model.OrderItem, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.OrderItem), args.Error(1)
}

func TestOrderService_CreateOrder_Success(t *testing.T) {
	// Arrange
	mockOrderRepo := new(MockOrderRepository)
	mockQuoteRepo := new(MockQuoteRepository)
	orderService := service.NewOrderService(mockOrderRepo, mockQuoteRepo)
	
	quoteID := uuid.New()
	customerID := uuid.New()
	companyID := uuid.New()
	
	quote := &model.Quote{
		ID:          quoteID,
		QuoteNo:     "QUO-2024-001",
		CustomerID:  customerID,
		Status:      "accepted",
		TotalAmount: 10000.00,
		Currency:    "USD",
		ValidUntil:  time.Now().AddDate(0, 1, 0),
	}
	
	createRequest := &service.CreateOrderRequest{
		QuoteID:         quoteID,
		PONumber:        "PO-2024-001",
		DeliveryDate:    time.Now().AddDate(0, 2, 0),
		DeliveryMethod:  "海運",
		ShippingAddress: "123 Main St, City, Country",
		PaymentTerms:    "T/T 30 days",
		DownPayment:     3000.00,
		Notes:           "Test order",
	}
	
	mockQuoteRepo.On("GetByID", mock.Anything, quoteID).Return(quote, nil)
	mockOrderRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.Order")).Return(nil)
	
	// Act
	result, err := orderService.CreateOrder(context.Background(), createRequest, companyID)
	
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, quoteID, result.QuoteID)
	assert.Equal(t, customerID, result.CustomerID)
	assert.Equal(t, "pending", result.Status)
	assert.Equal(t, "pending", result.PaymentStatus)
	assert.Equal(t, createRequest.PONumber, result.PONumber)
	assert.Equal(t, quote.TotalAmount, result.TotalAmount)
	
	mockQuoteRepo.AssertExpectations(t)
	mockOrderRepo.AssertExpectations(t)
}

func TestOrderService_CreateOrder_QuoteNotAccepted(t *testing.T) {
	// Arrange
	mockOrderRepo := new(MockOrderRepository)
	mockQuoteRepo := new(MockQuoteRepository)
	orderService := service.NewOrderService(mockOrderRepo, mockQuoteRepo)
	
	quoteID := uuid.New()
	companyID := uuid.New()
	
	quote := &model.Quote{
		ID:      quoteID,
		QuoteNo: "QUO-2024-001",
		Status:  "draft", // Not accepted
	}
	
	createRequest := &service.CreateOrderRequest{
		QuoteID:  quoteID,
		PONumber: "PO-2024-001",
	}
	
	mockQuoteRepo.On("GetByID", mock.Anything, quoteID).Return(quote, nil)
	
	// Act
	result, err := orderService.CreateOrder(context.Background(), createRequest, companyID)
	
	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "quote must be accepted")
	
	mockQuoteRepo.AssertExpectations(t)
}

func TestOrderService_GetOrder_Success(t *testing.T) {
	// Arrange
	mockOrderRepo := new(MockOrderRepository)
	mockQuoteRepo := new(MockQuoteRepository)
	orderService := service.NewOrderService(mockOrderRepo, mockQuoteRepo)
	
	orderID := uuid.New()
	
	expectedOrder := &model.Order{
		ID:            orderID,
		OrderNo:       "ORD-2024-001",
		Status:        "confirmed",
		PaymentStatus: "partial",
		TotalAmount:   10000.00,
		Currency:      "USD",
		CreatedAt:     time.Now(),
	}
	
	mockOrderRepo.On("GetByID", mock.Anything, orderID).Return(expectedOrder, nil)
	
	// Act
	result, err := orderService.GetOrder(context.Background(), orderID)
	
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedOrder.ID, result.ID)
	assert.Equal(t, expectedOrder.OrderNo, result.OrderNo)
	assert.Equal(t, expectedOrder.Status, result.Status)
	
	mockOrderRepo.AssertExpectations(t)
}

func TestOrderService_UpdateOrderStatus_Success(t *testing.T) {
	// Arrange
	mockOrderRepo := new(MockOrderRepository)
	mockQuoteRepo := new(MockQuoteRepository)
	orderService := service.NewOrderService(mockOrderRepo, mockQuoteRepo)
	
	orderID := uuid.New()
	
	existingOrder := &model.Order{
		ID:      orderID,
		OrderNo: "ORD-2024-001",
		Status:  "confirmed",
	}
	
	mockOrderRepo.On("GetByID", mock.Anything, orderID).Return(existingOrder, nil)
	mockOrderRepo.On("UpdateStatus", mock.Anything, orderID, "in_production", "Production started").Return(nil)
	
	// Act
	err := orderService.UpdateOrderStatus(context.Background(), orderID, "in_production", "Production started")
	
	// Assert
	assert.NoError(t, err)
	
	mockOrderRepo.AssertExpectations(t)
}

func TestOrderService_UpdateOrderStatus_InvalidTransition(t *testing.T) {
	// Arrange
	mockOrderRepo := new(MockOrderRepository)
	mockQuoteRepo := new(MockQuoteRepository)
	orderService := service.NewOrderService(mockOrderRepo, mockQuoteRepo)
	
	orderID := uuid.New()
	
	existingOrder := &model.Order{
		ID:      orderID,
		OrderNo: "ORD-2024-001",
		Status:  "completed", // Cannot change from completed
	}
	
	mockOrderRepo.On("GetByID", mock.Anything, orderID).Return(existingOrder, nil)
	
	// Act
	err := orderService.UpdateOrderStatus(context.Background(), orderID, "pending", "Invalid transition")
	
	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid status transition")
	
	mockOrderRepo.AssertExpectations(t)
}

func TestOrderService_ListOrders_Success(t *testing.T) {
	// Arrange
	mockOrderRepo := new(MockOrderRepository)
	mockQuoteRepo := new(MockQuoteRepository)
	orderService := service.NewOrderService(mockOrderRepo, mockQuoteRepo)
	
	companyID := uuid.New()
	
	expectedOrders := []*model.Order{
		{
			ID:        uuid.New(),
			OrderNo:   "ORD-2024-001",
			Status:    "confirmed",
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			OrderNo:   "ORD-2024-002",
			Status:    "in_production",
			CreatedAt: time.Now(),
		},
	}
	
	filter := map[string]interface{}{
		"company_id": companyID,
		"status":     "confirmed",
	}
	
	mockOrderRepo.On("List", mock.Anything, filter).Return(expectedOrders, nil)
	
	// Act
	result, err := orderService.ListOrders(context.Background(), filter)
	
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, expectedOrders[0].OrderNo, result[0].OrderNo)
	assert.Equal(t, expectedOrders[1].OrderNo, result[1].OrderNo)
	
	mockOrderRepo.AssertExpectations(t)
}

func TestOrderService_GetOrderItems_Success(t *testing.T) {
	// Arrange
	mockOrderRepo := new(MockOrderRepository)
	mockQuoteRepo := new(MockQuoteRepository)
	orderService := service.NewOrderService(mockOrderRepo, mockQuoteRepo)
	
	orderID := uuid.New()
	
	expectedItems := []*model.OrderItem{
		{
			ID:          uuid.New(),
			OrderID:     orderID,
			PartNo:      "HB-M8-20",
			Description: "Hex Bolt M8x20",
			Quantity:    5000,
			UnitPrice:   0.15,
			TotalPrice:  750.00,
		},
		{
			ID:          uuid.New(),
			OrderID:     orderID,
			PartNo:      "HB-M10-25",
			Description: "Hex Bolt M10x25",
			Quantity:    3000,
			UnitPrice:   0.25,
			TotalPrice:  750.00,
		},
	}
	
	mockOrderRepo.On("GetItems", mock.Anything, orderID).Return(expectedItems, nil)
	
	// Act
	result, err := orderService.GetOrderItems(context.Background(), orderID)
	
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, expectedItems[0].PartNo, result[0].PartNo)
	assert.Equal(t, expectedItems[1].PartNo, result[1].PartNo)
	
	mockOrderRepo.AssertExpectations(t)
}

func TestOrderService_CalculateOrderPayment_Success(t *testing.T) {
	// Arrange
	mockOrderRepo := new(MockOrderRepository)
	mockQuoteRepo := new(MockQuoteRepository)
	orderService := service.NewOrderService(mockOrderRepo, mockQuoteRepo)
	
	// Test payment calculation
	totalAmount := 10000.00
	downPaymentPercent := 30.0
	
	// Act
	downPayment, remainingAmount := orderService.CalculateOrderPayment(totalAmount, downPaymentPercent)
	
	// Assert
	expectedDownPayment := totalAmount * (downPaymentPercent / 100)
	expectedRemaining := totalAmount - expectedDownPayment
	
	assert.Equal(t, expectedDownPayment, downPayment)
	assert.Equal(t, expectedRemaining, remainingAmount)
}

func TestOrderService_ValidateStatusTransition(t *testing.T) {
	// Arrange
	mockOrderRepo := new(MockOrderRepository)
	mockQuoteRepo := new(MockQuoteRepository)
	orderService := service.NewOrderService(mockOrderRepo, mockQuoteRepo)
	
	testCases := []struct {
		currentStatus string
		newStatus     string
		shouldPass    bool
	}{
		{"pending", "confirmed", true},
		{"confirmed", "in_production", true},
		{"in_production", "quality_check", true},
		{"quality_check", "ready_to_ship", true},
		{"ready_to_ship", "shipped", true},
		{"shipped", "delivered", true},
		{"delivered", "completed", true},
		{"completed", "pending", false}, // Invalid transition
		{"cancelled", "confirmed", false}, // Invalid transition
		{"pending", "completed", false}, // Skip steps invalid
	}
	
	for _, tc := range testCases {
		// Act
		isValid := orderService.ValidateStatusTransition(tc.currentStatus, tc.newStatus)
		
		// Assert
		assert.Equal(t, tc.shouldPass, isValid, 
			"Status transition from %s to %s should be %v", 
			tc.currentStatus, tc.newStatus, tc.shouldPass)
	}
}