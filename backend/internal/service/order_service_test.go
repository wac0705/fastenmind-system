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

// MockOrderRepository is a mock implementation of OrderRepository
type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Create(ctx context.Context, order *models.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderRepository) Update(ctx context.Context, order *models.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderRepository) List(ctx context.Context, params repository.ListOrdersParams) ([]*models.Order, int64, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*models.Order), args.Get(1).(int64), args.Error(2)
}

func (m *MockOrderRepository) GetByOrderNo(ctx context.Context, orderNo string) (*models.Order, error) {
	args := m.Called(ctx, orderNo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderRepository) ExistsByOrderNo(ctx context.Context, orderNo string) (bool, error) {
	args := m.Called(ctx, orderNo)
	return args.Bool(0), args.Error(1)
}

func (m *MockOrderRepository) CreateWithItems(ctx context.Context, order *models.Order, items []models.OrderItem) error {
	args := m.Called(ctx, order, items)
	return args.Error(0)
}

func (m *MockOrderRepository) UpdateStatus(ctx context.Context, orderID uuid.UUID, status models.OrderStatus, updatedBy uuid.UUID) error {
	args := m.Called(ctx, orderID, status, updatedBy)
	return args.Error(0)
}

// Test suite
type OrderServiceTestSuite struct {
	service  *service.OrderService
	mockRepo *MockOrderRepository
	ctx      context.Context
}

func setupOrderTestSuite() *OrderServiceTestSuite {
	mockRepo := new(MockOrderRepository)
	svc := service.NewOrderService(mockRepo, nil) // DB is not needed for unit tests
	
	return &OrderServiceTestSuite{
		service:  svc,
		mockRepo: mockRepo,
		ctx:      context.Background(),
	}
}

func TestOrderService_CreateFromQuote(t *testing.T) {
	suite := setupOrderTestSuite()
	
	t.Run("successful creation from quote", func(t *testing.T) {
		// Arrange
		quoteID := uuid.New()
		createdBy := uuid.New()
		
		quote := &models.Quote{
			BaseModel:      models.BaseModel{ID: quoteID},
			QuoteNo:        "Q-2024-001",
			CustomerID:     uuid.New(),
			InquiryID:      uuid.New(),
			TotalAmount:    decimal.NewFromFloat(1000),
			Currency:       "USD",
			PaymentTerms:   "Net 30",
			DeliveryTerms:  "FOB Shanghai",
			LeadTime:       "4-6 weeks",
			Status:         models.QuoteStatusApproved,
			Items: []models.QuoteItem{
				{
					ProductName:   "Hex Bolt",
					Specification: "M10x50",
					Quantity:      100,
					Unit:          "PCS",
					UnitPrice:     decimal.NewFromFloat(10),
					TotalPrice:    decimal.NewFromFloat(1000),
				},
			},
		}
		
		req := service.CreateOrderFromQuoteRequest{
			QuoteID:       quoteID,
			CustomerPO:    "PO-12345",
			DeliveryDate:  time.Now().Add(45 * 24 * time.Hour),
			ShippingAddr:  "123 Main St, City",
			BillingAddr:   "456 Finance Ave, City",
			SpecialNotes:  "Handle with care",
			CreatedBy:     createdBy,
		}
		
		suite.mockRepo.On("ExistsByOrderNo", suite.ctx, mock.AnythingOfType("string")).Return(false, nil)
		suite.mockRepo.On("CreateWithItems", suite.ctx, mock.AnythingOfType("*models.Order"), mock.AnythingOfType("[]models.OrderItem")).Return(nil)
		
		// Act
		result, err := suite.service.CreateFromQuote(suite.ctx, req, quote)
		
		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, req.CustomerPO, result.CustomerPO)
		assert.Equal(t, quote.CustomerID, result.CustomerID)
		assert.Equal(t, quote.TotalAmount, result.TotalAmount)
		assert.Equal(t, models.OrderStatusPending, result.Status)
		assert.Len(t, result.Items, 1)
		
		suite.mockRepo.AssertExpectations(t)
	})
	
	t.Run("quote not approved", func(t *testing.T) {
		// Arrange
		quote := &models.Quote{
			Status: models.QuoteStatusDraft,
		}
		
		req := service.CreateOrderFromQuoteRequest{
			QuoteID:    uuid.New(),
			CustomerPO: "PO-12345",
			CreatedBy:  uuid.New(),
		}
		
		// Act
		result, err := suite.service.CreateFromQuote(suite.ctx, req, quote)
		
		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "quote must be approved")
		
		suite.mockRepo.AssertExpectations(t)
	})
}

func TestOrderService_UpdateStatus(t *testing.T) {
	suite := setupOrderTestSuite()
	
	t.Run("successful status update to confirmed", func(t *testing.T) {
		// Arrange
		orderID := uuid.New()
		updatedBy := uuid.New()
		
		order := &models.Order{
			BaseModel: models.BaseModel{ID: orderID},
			Status:    models.OrderStatusPending,
		}
		
		suite.mockRepo.On("GetByID", suite.ctx, orderID).Return(order, nil)
		suite.mockRepo.On("UpdateStatus", suite.ctx, orderID, models.OrderStatusConfirmed, updatedBy).Return(nil)
		
		// Act
		err := suite.service.UpdateStatus(suite.ctx, orderID, models.OrderStatusConfirmed, updatedBy, "Payment received")
		
		// Assert
		assert.NoError(t, err)
		suite.mockRepo.AssertExpectations(t)
	})
	
	t.Run("invalid status transition", func(t *testing.T) {
		// Arrange
		orderID := uuid.New()
		updatedBy := uuid.New()
		
		order := &models.Order{
			BaseModel: models.BaseModel{ID: orderID},
			Status:    models.OrderStatusCancelled, // Already cancelled
		}
		
		suite.mockRepo.On("GetByID", suite.ctx, orderID).Return(order, nil)
		
		// Act
		err := suite.service.UpdateStatus(suite.ctx, orderID, models.OrderStatusConfirmed, updatedBy, "")
		
		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid status transition")
		
		suite.mockRepo.AssertExpectations(t)
	})
}

func TestOrderService_List(t *testing.T) {
	suite := setupOrderTestSuite()
	
	t.Run("successful list with date range filter", func(t *testing.T) {
		// Arrange
		startDate := time.Now().Add(-30 * 24 * time.Hour)
		endDate := time.Now()
		
		params := service.ListOrdersParams{
			CompanyID: uuid.New(),
			StartDate: &startDate,
			EndDate:   &endDate,
			Page:      1,
			PageSize:  10,
		}
		
		orders := []*models.Order{
			{
				BaseModel:   models.BaseModel{ID: uuid.New()},
				OrderNo:     "ORD-2024-001",
				CustomerPO:  "PO-123",
				TotalAmount: decimal.NewFromFloat(1000),
				Status:      models.OrderStatusConfirmed,
			},
			{
				BaseModel:   models.BaseModel{ID: uuid.New()},
				OrderNo:     "ORD-2024-002",
				CustomerPO:  "PO-124",
				TotalAmount: decimal.NewFromFloat(2000),
				Status:      models.OrderStatusInProduction,
			},
		}
		
		expectedRepoParams := repository.ListOrdersParams{
			CompanyID: params.CompanyID,
			StartDate: params.StartDate,
			EndDate:   params.EndDate,
			Offset:    0,
			Limit:     10,
		}
		
		suite.mockRepo.On("List", suite.ctx, expectedRepoParams).Return(orders, int64(2), nil)
		
		// Act
		result, err := suite.service.List(suite.ctx, params)
		
		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Items, 2)
		assert.Equal(t, int64(2), result.Total)
		assert.Equal(t, 1, result.Page)
		assert.Equal(t, 10, result.PageSize)
		assert.Equal(t, 1, result.TotalPages)
		
		suite.mockRepo.AssertExpectations(t)
	})
}

func TestOrderService_GetByID(t *testing.T) {
	suite := setupOrderTestSuite()
	
	t.Run("successful retrieval", func(t *testing.T) {
		// Arrange
		orderID := uuid.New()
		order := &models.Order{
			BaseModel:  models.BaseModel{ID: orderID},
			OrderNo:    "ORD-2024-001",
			CustomerPO: "PO-123",
			Status:     models.OrderStatusConfirmed,
		}
		
		suite.mockRepo.On("GetByID", suite.ctx, orderID).Return(order, nil)
		
		// Act
		result, err := suite.service.GetByID(suite.ctx, orderID)
		
		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, order.OrderNo, result.OrderNo)
		assert.Equal(t, order.CustomerPO, result.CustomerPO)
		
		suite.mockRepo.AssertExpectations(t)
	})
	
	t.Run("order not found", func(t *testing.T) {
		// Arrange
		orderID := uuid.New()
		suite.mockRepo.On("GetByID", suite.ctx, orderID).Return(nil, gorm.ErrRecordNotFound)
		
		// Act
		result, err := suite.service.GetByID(suite.ctx, orderID)
		
		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, service.ErrNotFound, err)
		
		suite.mockRepo.AssertExpectations(t)
	})
}

// Benchmark tests
func BenchmarkOrderService_CreateFromQuote(b *testing.B) {
	suite := setupOrderTestSuite()
	
	quote := &models.Quote{
		BaseModel:     models.BaseModel{ID: uuid.New()},
		QuoteNo:       "Q-2024-001",
		CustomerID:    uuid.New(),
		TotalAmount:   decimal.NewFromFloat(1000),
		Status:        models.QuoteStatusApproved,
		PaymentTerms:  "Net 30",
		DeliveryTerms: "FOB",
		Items: []models.QuoteItem{
			{
				ProductName: "Test Product",
				Quantity:    100,
				UnitPrice:   decimal.NewFromFloat(10),
				TotalPrice:  decimal.NewFromFloat(1000),
			},
		},
	}
	
	req := service.CreateOrderFromQuoteRequest{
		QuoteID:      quote.ID,
		CustomerPO:   "PO-12345",
		DeliveryDate: time.Now().Add(30 * 24 * time.Hour),
		CreatedBy:    uuid.New(),
	}
	
	suite.mockRepo.On("ExistsByOrderNo", mock.Anything, mock.AnythingOfType("string")).Return(false, nil)
	suite.mockRepo.On("CreateWithItems", mock.Anything, mock.AnythingOfType("*models.Order"), mock.AnythingOfType("[]models.OrderItem")).Return(nil)
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_, _ = suite.service.CreateFromQuote(context.Background(), req, quote)
	}
}

// Table-driven tests for order validation
func TestOrderService_ValidateOrderData(t *testing.T) {
	tests := []struct {
		name    string
		req     service.CreateOrderFromQuoteRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid order",
			req: service.CreateOrderFromQuoteRequest{
				QuoteID:      uuid.New(),
				CustomerPO:   "PO-12345",
				DeliveryDate: time.Now().Add(30 * 24 * time.Hour),
				ShippingAddr: "123 Main St",
				BillingAddr:  "456 Finance Ave",
				CreatedBy:    uuid.New(),
			},
			wantErr: false,
		},
		{
			name: "missing customer PO",
			req: service.CreateOrderFromQuoteRequest{
				QuoteID:      uuid.New(),
				CustomerPO:   "",
				DeliveryDate: time.Now().Add(30 * 24 * time.Hour),
				CreatedBy:    uuid.New(),
			},
			wantErr: true,
			errMsg:  "customer PO is required",
		},
		{
			name: "past delivery date",
			req: service.CreateOrderFromQuoteRequest{
				QuoteID:      uuid.New(),
				CustomerPO:   "PO-12345",
				DeliveryDate: time.Now().Add(-24 * time.Hour),
				CreatedBy:    uuid.New(),
			},
			wantErr: true,
			errMsg:  "delivery date must be in the future",
		},
		{
			name: "missing shipping address",
			req: service.CreateOrderFromQuoteRequest{
				QuoteID:      uuid.New(),
				CustomerPO:   "PO-12345",
				DeliveryDate: time.Now().Add(30 * 24 * time.Hour),
				ShippingAddr: "",
				CreatedBy:    uuid.New(),
			},
			wantErr: true,
			errMsg:  "shipping address is required",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ValidateOrderData(tt.req)
			
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

// Test order status transitions
func TestOrderService_ValidateStatusTransition(t *testing.T) {
	tests := []struct {
		name        string
		fromStatus  models.OrderStatus
		toStatus    models.OrderStatus
		valid       bool
		description string
	}{
		{
			name:        "pending to confirmed",
			fromStatus:  models.OrderStatusPending,
			toStatus:    models.OrderStatusConfirmed,
			valid:       true,
			description: "Valid transition from pending to confirmed",
		},
		{
			name:        "confirmed to in_production",
			fromStatus:  models.OrderStatusConfirmed,
			toStatus:    models.OrderStatusInProduction,
			valid:       true,
			description: "Valid transition from confirmed to in production",
		},
		{
			name:        "in_production to ready_to_ship",
			fromStatus:  models.OrderStatusInProduction,
			toStatus:    models.OrderStatusReadyToShip,
			valid:       true,
			description: "Valid transition from in production to ready to ship",
		},
		{
			name:        "ready_to_ship to shipped",
			fromStatus:  models.OrderStatusReadyToShip,
			toStatus:    models.OrderStatusShipped,
			valid:       true,
			description: "Valid transition from ready to ship to shipped",
		},
		{
			name:        "shipped to delivered",
			fromStatus:  models.OrderStatusShipped,
			toStatus:    models.OrderStatusDelivered,
			valid:       true,
			description: "Valid transition from shipped to delivered",
		},
		{
			name:        "pending to cancelled",
			fromStatus:  models.OrderStatusPending,
			toStatus:    models.OrderStatusCancelled,
			valid:       true,
			description: "Orders can be cancelled from pending",
		},
		{
			name:        "delivered to confirmed",
			fromStatus:  models.OrderStatusDelivered,
			toStatus:    models.OrderStatusConfirmed,
			valid:       false,
			description: "Cannot go back from delivered to confirmed",
		},
		{
			name:        "cancelled to confirmed",
			fromStatus:  models.OrderStatusCancelled,
			toStatus:    models.OrderStatusConfirmed,
			valid:       false,
			description: "Cannot reactivate cancelled orders",
		},
		{
			name:        "shipped to pending",
			fromStatus:  models.OrderStatusShipped,
			toStatus:    models.OrderStatusPending,
			valid:       false,
			description: "Cannot go back to pending from shipped",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ValidateOrderStatusTransition(tt.fromStatus, tt.toStatus)
			
			if tt.valid {
				assert.NoError(t, err, tt.description)
			} else {
				assert.Error(t, err, tt.description)
			}
		})
	}
}

// Test shipping tracking
func TestOrderService_UpdateShippingInfo(t *testing.T) {
	suite := setupOrderTestSuite()
	
	t.Run("successful shipping info update", func(t *testing.T) {
		// Arrange
		orderID := uuid.New()
		updatedBy := uuid.New()
		
		order := &models.Order{
			BaseModel: models.BaseModel{ID: orderID},
			Status:    models.OrderStatusReadyToShip,
		}
		
		shippingInfo := service.ShippingInfo{
			Carrier:        "DHL",
			TrackingNumber: "1234567890",
			EstimatedDate:  time.Now().Add(5 * 24 * time.Hour),
			ActualDate:     nil,
		}
		
		suite.mockRepo.On("GetByID", suite.ctx, orderID).Return(order, nil)
		suite.mockRepo.On("Update", suite.ctx, mock.AnythingOfType("*models.Order")).Return(nil)
		
		// Act
		result, err := suite.service.UpdateShippingInfo(suite.ctx, orderID, shippingInfo, updatedBy)
		
		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, shippingInfo.Carrier, result.ShippingCarrier)
		assert.Equal(t, shippingInfo.TrackingNumber, result.TrackingNumber)
		assert.NotNil(t, result.EstimatedDelivery)
		
		suite.mockRepo.AssertExpectations(t)
	})
}