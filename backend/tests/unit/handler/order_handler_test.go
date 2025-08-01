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

// MockOrderService is a mock implementation of OrderService
type MockOrderService struct {
	mock.Mock
}

func (m *MockOrderService) CreateOrder(ctx context.Context, request *service.CreateOrderRequest, companyID uuid.UUID) (*model.Order, error) {
	args := m.Called(ctx, request, companyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

func (m *MockOrderService) GetOrder(ctx context.Context, id uuid.UUID) (*model.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

func (m *MockOrderService) UpdateOrder(ctx context.Context, id uuid.UUID, request *service.UpdateOrderRequest) (*model.Order, error) {
	args := m.Called(ctx, id, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

func (m *MockOrderService) UpdateOrderStatus(ctx context.Context, id uuid.UUID, status, notes string) error {
	args := m.Called(ctx, id, status, notes)
	return args.Error(0)
}

func (m *MockOrderService) ListOrders(ctx context.Context, filter map[string]interface{}) ([]*model.Order, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Order), args.Error(1)
}

func (m *MockOrderService) GetOrderItems(ctx context.Context, orderID uuid.UUID) ([]*model.OrderItem, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.OrderItem), args.Error(1)
}

func (m *MockOrderService) DeleteOrder(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestOrderHandler_CreateOrder_Success(t *testing.T) {
	// Arrange
	e := echo.New()
	mockOrderService := new(MockOrderService)
	orderHandler := handler.NewOrderHandler(mockOrderService)
	
	quoteID := uuid.New()
	companyID := uuid.New()
	
	createRequest := handler.CreateOrderRequest{
		QuoteID:         quoteID,
		PONumber:        "PO-2024-001",
		DeliveryDate:    time.Now().AddDate(0, 2, 0),
		DeliveryMethod:  "海運",
		ShippingAddress: "123 Main St",
		PaymentTerms:    "T/T 30 days",
		DownPayment:     3000.00,
		Notes:           "Test order",
	}
	
	expectedOrder := &model.Order{
		ID:              uuid.New(),
		OrderNo:         "ORD-2024-001",
		QuoteID:         quoteID,
		Status:          "pending",
		PaymentStatus:   "pending",
		PONumber:        "PO-2024-001",
		TotalAmount:     10000.00,
		Currency:        "USD",
		DeliveryDate:    createRequest.DeliveryDate,
		DeliveryMethod:  "海運",
		ShippingAddress: "123 Main St",
		PaymentTerms:    "T/T 30 days",
		DownPayment:     3000.00,
		Notes:           "Test order",
	}
	
	mockOrderService.On("CreateOrder", mock.Anything, mock.AnythingOfType("*service.CreateOrderRequest"), companyID).Return(expectedOrder, nil)
	
	requestBody, _ := json.Marshal(createRequest)
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	
	// Mock JWT claims
	c.Set("user_id", uuid.New().String())
	c.Set("company_id", companyID.String())
	c.Set("role", "sales")
	
	// Act
	err := orderHandler.CreateOrder(c)
	
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	
	var response model.Order
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedOrder.OrderNo, response.OrderNo)
	assert.Equal(t, expectedOrder.Status, response.Status)
	
	mockOrderService.AssertExpectations(t)
}

func TestOrderHandler_CreateOrder_ValidationError(t *testing.T) {
	// Arrange
	e := echo.New()
	mockOrderService := new(MockOrderService)
	orderHandler := handler.NewOrderHandler(mockOrderService)
	
	// Invalid request - missing required fields
	createRequest := handler.CreateOrderRequest{
		// Missing QuoteID and PONumber
		DeliveryDate: time.Now().AddDate(0, 2, 0),
	}
	
	requestBody, _ := json.Marshal(createRequest)
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	
	// Mock JWT claims
	c.Set("user_id", uuid.New().String())
	c.Set("company_id", uuid.New().String())
	c.Set("role", "sales")
	
	// Act
	err := orderHandler.CreateOrder(c)
	
	// Assert
	assert.Error(t, err)
	httpError, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, httpError.Code)
}

func TestOrderHandler_GetOrder_Success(t *testing.T) {
	// Arrange
	e := echo.New()
	mockOrderService := new(MockOrderService)
	orderHandler := handler.NewOrderHandler(mockOrderService)
	
	orderID := uuid.New()
	
	expectedOrder := &model.Order{
		ID:            orderID,
		OrderNo:       "ORD-2024-001",
		Status:        "confirmed",
		PaymentStatus: "partial",
		TotalAmount:   10000.00,
		Currency:      "USD",
	}
	
	mockOrderService.On("GetOrder", mock.Anything, orderID).Return(expectedOrder, nil)
	
	req := httptest.NewRequest(http.MethodGet, "/orders/"+orderID.String(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(orderID.String())
	
	// Mock JWT claims
	c.Set("user_id", uuid.New().String())
	c.Set("company_id", uuid.New().String())
	c.Set("role", "engineer")
	
	// Act
	err := orderHandler.GetOrder(c)
	
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	
	var response model.Order
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedOrder.ID, response.ID)
	assert.Equal(t, expectedOrder.OrderNo, response.OrderNo)
	
	mockOrderService.AssertExpectations(t)
}

func TestOrderHandler_UpdateOrderStatus_Success(t *testing.T) {
	// Arrange
	e := echo.New()
	mockOrderService := new(MockOrderService)
	orderHandler := handler.NewOrderHandler(mockOrderService)
	
	orderID := uuid.New()
	
	updateRequest := handler.UpdateOrderStatusRequest{
		Status: "in_production",
		Notes:  "Production started",
	}
	
	mockOrderService.On("UpdateOrderStatus", mock.Anything, orderID, "in_production", "Production started").Return(nil)
	
	requestBody, _ := json.Marshal(updateRequest)
	req := httptest.NewRequest(http.MethodPut, "/orders/"+orderID.String()+"/status", bytes.NewBuffer(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(orderID.String())
	
	// Mock JWT claims
	c.Set("user_id", uuid.New().String())
	c.Set("company_id", uuid.New().String())
	c.Set("role", "manager")
	
	// Act
	err := orderHandler.UpdateOrderStatus(c)
	
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	
	mockOrderService.AssertExpectations(t)
}

func TestOrderHandler_ListOrders_Success(t *testing.T) {
	// Arrange
	e := echo.New()
	mockOrderService := new(MockOrderService)
	orderHandler := handler.NewOrderHandler(mockOrderService)
	
	companyID := uuid.New()
	
	expectedOrders := []*model.Order{
		{
			ID:      uuid.New(),
			OrderNo: "ORD-2024-001",
			Status:  "confirmed",
		},
		{
			ID:      uuid.New(),
			OrderNo: "ORD-2024-002",
			Status:  "in_production",
		},
	}
	
	filter := map[string]interface{}{
		"company_id": companyID.String(),
	}
	
	mockOrderService.On("ListOrders", mock.Anything, filter).Return(expectedOrders, nil)
	
	req := httptest.NewRequest(http.MethodGet, "/orders", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	
	// Mock JWT claims
	c.Set("user_id", uuid.New().String())
	c.Set("company_id", companyID.String())
	c.Set("role", "sales")
	
	// Act
	err := orderHandler.ListOrders(c)
	
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	
	var response []model.Order
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	
	mockOrderService.AssertExpectations(t)
}

func TestOrderHandler_GetOrderItems_Success(t *testing.T) {
	// Arrange
	e := echo.New()
	mockOrderService := new(MockOrderService)
	orderHandler := handler.NewOrderHandler(mockOrderService)
	
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
	}
	
	mockOrderService.On("GetOrderItems", mock.Anything, orderID).Return(expectedItems, nil)
	
	req := httptest.NewRequest(http.MethodGet, "/orders/"+orderID.String()+"/items", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(orderID.String())
	
	// Mock JWT claims
	c.Set("user_id", uuid.New().String())
	c.Set("company_id", uuid.New().String())
	c.Set("role", "engineer")
	
	// Act
	err := orderHandler.GetOrderItems(c)
	
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	
	var response []model.OrderItem
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, expectedItems[0].PartNo, response[0].PartNo)
	
	mockOrderService.AssertExpectations(t)
}

func TestOrderHandler_DeleteOrder_Success(t *testing.T) {
	// Arrange
	e := echo.New()
	mockOrderService := new(MockOrderService)
	orderHandler := handler.NewOrderHandler(mockOrderService)
	
	orderID := uuid.New()
	
	mockOrderService.On("DeleteOrder", mock.Anything, orderID).Return(nil)
	
	req := httptest.NewRequest(http.MethodDelete, "/orders/"+orderID.String(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(orderID.String())
	
	// Mock JWT claims - only admin can delete
	c.Set("user_id", uuid.New().String())
	c.Set("company_id", uuid.New().String())
	c.Set("role", "admin")
	
	// Act
	err := orderHandler.DeleteOrder(c)
	
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)
	
	mockOrderService.AssertExpectations(t)
}

func TestOrderHandler_DeleteOrder_Forbidden(t *testing.T) {
	// Arrange
	e := echo.New()
	mockOrderService := new(MockOrderService)
	orderHandler := handler.NewOrderHandler(mockOrderService)
	
	orderID := uuid.New()
	
	req := httptest.NewRequest(http.MethodDelete, "/orders/"+orderID.String(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(orderID.String())
	
	// Mock JWT claims - sales role cannot delete
	c.Set("user_id", uuid.New().String())
	c.Set("company_id", uuid.New().String())
	c.Set("role", "sales")
	
	// Act
	err := orderHandler.DeleteOrder(c)
	
	// Assert
	assert.Error(t, err)
	httpError, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusForbidden, httpError.Code)
}