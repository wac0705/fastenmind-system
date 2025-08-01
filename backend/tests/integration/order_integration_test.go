package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
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

type OrderIntegrationTestSuite struct {
	suite.Suite
	db           *gorm.DB
	echo         *echo.Echo
	orderHandler *handler.OrderHandler
	cleanup      func()
	
	// Test data
	company  *model.Company
	customer *model.Customer
	sales    *model.Account
	engineer *model.Account
	quote    *model.Quote
}

func (suite *OrderIntegrationTestSuite) SetupSuite() {
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
		&model.Customer{},
		&model.Inquiry{},
		&model.Quote{},
		&model.Order{},
		&model.OrderItem{},
	)
	suite.Require().NoError(err)
	
	// Setup services and handlers
	orderRepo := repository.NewOrderRepository(db)
	quoteRepo := repository.NewQuoteRepository(db)
	orderService := service.NewOrderService(orderRepo, quoteRepo)
	suite.orderHandler = handler.NewOrderHandler(orderService)
	
	// Setup Echo
	suite.echo = echo.New()
	
	// Setup cleanup function
	suite.cleanup = func() {
		// Clean up test data in correct order due to foreign key constraints
		db.Exec("DELETE FROM order_items")
		db.Exec("DELETE FROM orders")
		db.Exec("DELETE FROM quotes")
		db.Exec("DELETE FROM inquiries")
		db.Exec("DELETE FROM customers")
		db.Exec("DELETE FROM accounts")
		db.Exec("DELETE FROM companies")
	}
}

func (suite *OrderIntegrationTestSuite) TearDownSuite() {
	if suite.cleanup != nil {
		suite.cleanup()
	}
}

func (suite *OrderIntegrationTestSuite) SetupTest() {
	// Clean up before each test
	suite.cleanup()
	
	// Create test data
	suite.setupTestData()
}

func (suite *OrderIntegrationTestSuite) setupTestData() {
	// Create test company
	suite.company = &model.Company{
		Code:    "TEST001",
		Name:    "Test Company",
		Country: "US",
		Type:    "headquarters",
	}
	err := suite.db.Create(suite.company).Error
	suite.Require().NoError(err)
	
	// Create test customer
	suite.customer = &model.Customer{
		CompanyID:    suite.company.ID,
		CustomerCode: "CUST001",
		Name:         "Test Customer",
		Country:      "US",
		Currency:     "USD",
		IsActive:     true,
	}
	err = suite.db.Create(suite.customer).Error
	suite.Require().NoError(err)
	
	// Create test sales account
	suite.sales = &model.Account{
		CompanyID:    suite.company.ID,
		Username:     "sales1",
		Email:        "sales@test.com",
		PasswordHash: "$2a$10$XYZ...",
		FullName:     "Sales User",
		Role:         "sales",
		IsActive:     true,
	}
	err = suite.db.Create(suite.sales).Error
	suite.Require().NoError(err)
	
	// Create test engineer account
	suite.engineer = &model.Account{
		CompanyID:    suite.company.ID,
		Username:     "engineer1",
		Email:        "engineer@test.com",
		PasswordHash: "$2a$10$XYZ...",
		FullName:     "Engineer User",
		Role:         "engineer",
		IsActive:     true,
	}
	err = suite.db.Create(suite.engineer).Error
	suite.Require().NoError(err)
	
	// Create test quote (accepted)
	suite.quote = &model.Quote{
		QuoteNo:      "QUO-2024-001",
		CompanyID:    suite.company.ID,
		CustomerID:   suite.customer.ID,
		EngineerID:   suite.engineer.ID,
		Status:       "accepted",
		TotalAmount:  10000.00,
		Currency:     "USD",
		ValidUntil:   time.Now().AddDate(0, 1, 0),
		DeliveryDays: 30,
		PaymentTerms: "T/T 30 days",
	}
	err = suite.db.Create(suite.quote).Error
	suite.Require().NoError(err)
}

func (suite *OrderIntegrationTestSuite) TestCreateOrder_Success() {
	createRequest := handler.CreateOrderRequest{
		QuoteID:         suite.quote.ID,
		PONumber:        "PO-2024-001",
		DeliveryDate:    time.Now().AddDate(0, 2, 0),
		DeliveryMethod:  "海運",
		ShippingAddress: "123 Main St, City, Country",
		PaymentTerms:    "T/T 30 days",
		DownPayment:     3000.00,
		Notes:           "Integration test order",
	}
	
	requestBody, _ := json.Marshal(createRequest)
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)
	
	// Mock JWT claims
	c.Set("user_id", suite.sales.ID.String())
	c.Set("company_id", suite.company.ID.String())
	c.Set("role", "sales")
	
	err := suite.orderHandler.CreateOrder(c)
	suite.NoError(err)
	suite.Equal(http.StatusCreated, rec.Code)
	
	var response model.Order
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal(suite.quote.ID, response.QuoteID)
	suite.Equal(suite.customer.ID, response.CustomerID)
	suite.Equal("pending", response.Status)
	suite.Equal("pending", response.PaymentStatus)
	suite.Equal(createRequest.PONumber, response.PONumber)
	
	// Verify order was saved to database
	var savedOrder model.Order
	err = suite.db.First(&savedOrder, response.ID).Error
	suite.NoError(err)
	suite.Equal(response.ID, savedOrder.ID)
}

func (suite *OrderIntegrationTestSuite) TestCreateOrder_QuoteNotAccepted() {
	// Update quote status to draft
	suite.quote.Status = "draft"
	err := suite.db.Save(suite.quote).Error
	suite.Require().NoError(err)
	
	createRequest := handler.CreateOrderRequest{
		QuoteID:  suite.quote.ID,
		PONumber: "PO-2024-001",
	}
	
	requestBody, _ := json.Marshal(createRequest)
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)
	
	// Mock JWT claims
	c.Set("user_id", suite.sales.ID.String())
	c.Set("company_id", suite.company.ID.String())
	c.Set("role", "sales")
	
	err = suite.orderHandler.CreateOrder(c)
	suite.Error(err)
}

func (suite *OrderIntegrationTestSuite) TestGetOrder_Success() {
	// Create test order first
	order := &model.Order{
		OrderNo:         "ORD-2024-001",
		QuoteID:         suite.quote.ID,
		CompanyID:       suite.company.ID,
		CustomerID:      suite.customer.ID,
		Status:          "confirmed",
		PaymentStatus:   "partial",
		PONumber:        "PO-2024-001",
		TotalAmount:     10000.00,
		Currency:        "USD",
		DeliveryDate:    time.Now().AddDate(0, 2, 0),
		DeliveryMethod:  "海運",
		ShippingAddress: "123 Main St",
		PaymentTerms:    "T/T 30 days",
		DownPayment:     3000.00,
		PaidAmount:      3000.00,
	}
	err := suite.db.Create(order).Error
	suite.Require().NoError(err)
	
	req := httptest.NewRequest(http.MethodGet, "/orders/"+order.ID.String(), nil)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(order.ID.String())
	
	// Mock JWT claims
	c.Set("user_id", suite.engineer.ID.String())
	c.Set("company_id", suite.company.ID.String())
	c.Set("role", "engineer")
	
	err = suite.orderHandler.GetOrder(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
	
	var response model.Order
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal(order.ID, response.ID)
	suite.Equal(order.OrderNo, response.OrderNo)
	suite.Equal(order.Status, response.Status)
}

func (suite *OrderIntegrationTestSuite) TestOrderStatusWorkflow() {
	// Create test order
	order := &model.Order{
		OrderNo:       "ORD-2024-001",
		QuoteID:       suite.quote.ID,
		CompanyID:     suite.company.ID,
		CustomerID:    suite.customer.ID,
		Status:        "pending",
		PaymentStatus: "pending",
		TotalAmount:   10000.00,
		Currency:      "USD",
	}
	err := suite.db.Create(order).Error
	suite.Require().NoError(err)
	
	// Test status progression: pending -> confirmed -> in_production -> shipped -> delivered -> completed
	statusTransitions := []struct {
		newStatus string
		notes     string
	}{
		{"confirmed", "Order confirmed by customer"},
		{"in_production", "Production started"},
		{"quality_check", "Quality inspection"},
		{"ready_to_ship", "Ready for shipping"},
		{"shipped", "Shipped to customer"},
		{"delivered", "Delivered to customer"},
		{"completed", "Order completed"},
	}
	
	for _, transition := range statusTransitions {
		updateRequest := handler.UpdateOrderStatusRequest{
			Status: transition.newStatus,
			Notes:  transition.notes,
		}
		
		requestBody, _ := json.Marshal(updateRequest)
		req := httptest.NewRequest(http.MethodPut, "/orders/"+order.ID.String()+"/status", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := suite.echo.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(order.ID.String())
		
		// Mock JWT claims
		c.Set("user_id", suite.engineer.ID.String())
		c.Set("company_id", suite.company.ID.String())
		c.Set("role", "manager")
		
		err = suite.orderHandler.UpdateOrderStatus(c)
		suite.NoError(err, "Failed to transition to status: %s", transition.newStatus)
		suite.Equal(http.StatusOK, rec.Code)
		
		// Verify status was updated in database
		var updatedOrder model.Order
		err = suite.db.First(&updatedOrder, order.ID).Error
		suite.NoError(err)
		suite.Equal(transition.newStatus, updatedOrder.Status)
	}
}

func (suite *OrderIntegrationTestSuite) TestListOrders_WithFilters() {
	// Create multiple test orders
	orders := []*model.Order{
		{
			OrderNo:    "ORD-2024-001",
			QuoteID:    suite.quote.ID,
			CompanyID:  suite.company.ID,
			CustomerID: suite.customer.ID,
			Status:     "confirmed",
			Currency:   "USD",
		},
		{
			OrderNo:    "ORD-2024-002",
			QuoteID:    suite.quote.ID,
			CompanyID:  suite.company.ID,
			CustomerID: suite.customer.ID,
			Status:     "in_production",
			Currency:   "USD",
		},
		{
			OrderNo:    "ORD-2024-003",
			QuoteID:    suite.quote.ID,
			CompanyID:  suite.company.ID,
			CustomerID: suite.customer.ID,
			Status:     "completed",
			Currency:   "USD",
		},
	}
	
	for _, order := range orders {
		err := suite.db.Create(order).Error
		suite.Require().NoError(err)
	}
	
	// Test listing all orders
	req := httptest.NewRequest(http.MethodGet, "/orders", nil)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)
	
	// Mock JWT claims
	c.Set("user_id", suite.engineer.ID.String())
	c.Set("company_id", suite.company.ID.String())
	c.Set("role", "engineer")
	
	err := suite.orderHandler.ListOrders(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
	
	var response []model.Order
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Len(response, 3)
	
	// Test filtering by status
	req = httptest.NewRequest(http.MethodGet, "/orders?status=confirmed", nil)
	rec = httptest.NewRecorder()
	c = suite.echo.NewContext(req, rec)
	
	// Mock JWT claims
	c.Set("user_id", suite.engineer.ID.String())
	c.Set("company_id", suite.company.ID.String())
	c.Set("role", "engineer")
	
	err = suite.orderHandler.ListOrders(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
	
	var filteredResponse []model.Order
	err = json.Unmarshal(rec.Body.Bytes(), &filteredResponse)
	suite.NoError(err)
	suite.Len(filteredResponse, 1)
	suite.Equal("confirmed", filteredResponse[0].Status)
}

func (suite *OrderIntegrationTestSuite) TestOrderItemsManagement() {
	// Create test order
	order := &model.Order{
		OrderNo:    "ORD-2024-001",
		QuoteID:    suite.quote.ID,
		CompanyID:  suite.company.ID,
		CustomerID: suite.customer.ID,
		Status:     "confirmed",
		Currency:   "USD",
	}
	err := suite.db.Create(order).Error
	suite.Require().NoError(err)
	
	// Create order items
	items := []*model.OrderItem{
		{
			OrderID:     order.ID,
			PartNo:      "HB-M8-20",
			Description: "Hex Bolt M8x20",
			Material:    "Steel",
			Quantity:    5000,
			UnitPrice:   0.15,
			TotalPrice:  750.00,
		},
		{
			OrderID:     order.ID,
			PartNo:      "HB-M10-25",
			Description: "Hex Bolt M10x25",
			Material:    "Stainless Steel",
			Quantity:    3000,
			UnitPrice:   0.25,
			TotalPrice:  750.00,
		},
	}
	
	for _, item := range items {
		err := suite.db.Create(item).Error
		suite.Require().NoError(err)
	}
	
	// Test getting order items
	req := httptest.NewRequest(http.MethodGet, "/orders/"+order.ID.String()+"/items", nil)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(order.ID.String())
	
	// Mock JWT claims
	c.Set("user_id", suite.engineer.ID.String())
	c.Set("company_id", suite.company.ID.String())
	c.Set("role", "engineer")
	
	err = suite.orderHandler.GetOrderItems(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
	
	var response []model.OrderItem
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Len(response, 2)
	suite.Equal("HB-M8-20", response[0].PartNo)
	suite.Equal("HB-M10-25", response[1].PartNo)
}

func (suite *OrderIntegrationTestSuite) TestOrderPermissions() {
	// Create test order
	order := &model.Order{
		OrderNo:    "ORD-2024-001",
		QuoteID:    suite.quote.ID,
		CompanyID:  suite.company.ID,
		CustomerID: suite.customer.ID,
		Status:     "confirmed",
		Currency:   "USD",
	}
	err := suite.db.Create(order).Error
	suite.Require().NoError(err)
	
	// Test permissions for different roles
	testCases := []struct {
		role           string
		operation      string
		shouldAccess   bool
		expectedStatus int
	}{
		{"admin", "view", true, http.StatusOK},
		{"manager", "view", true, http.StatusOK},
		{"engineer", "view", true, http.StatusOK},
		{"sales", "view", true, http.StatusOK},
		{"viewer", "view", true, http.StatusOK},
		{"admin", "delete", true, http.StatusNoContent},
		{"manager", "delete", true, http.StatusNoContent},
		{"engineer", "delete", false, http.StatusForbidden},
		{"sales", "delete", false, http.StatusForbidden},
		{"viewer", "delete", false, http.StatusForbidden},
	}
	
	for _, tc := range testCases {
		var req *http.Request
		var expectedMethod string
		
		if tc.operation == "view" {
			req = httptest.NewRequest(http.MethodGet, "/orders/"+order.ID.String(), nil)
			expectedMethod = "GetOrder"
		} else if tc.operation == "delete" {
			req = httptest.NewRequest(http.MethodDelete, "/orders/"+order.ID.String(), nil)
			expectedMethod = "DeleteOrder"
		}
		
		rec := httptest.NewRecorder()
		c := suite.echo.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(order.ID.String())
		
		// Mock JWT claims with different roles
		c.Set("user_id", suite.engineer.ID.String())
		c.Set("company_id", suite.company.ID.String())
		c.Set("role", tc.role)
		
		var err error
		if expectedMethod == "GetOrder" {
			err = suite.orderHandler.GetOrder(c)
		} else if expectedMethod == "DeleteOrder" {
			err = suite.orderHandler.DeleteOrder(c)
		}
		
		if tc.shouldAccess {
			suite.NoError(err, "Role %s should be able to %s", tc.role, tc.operation)
			suite.Equal(tc.expectedStatus, rec.Code)
		} else {
			suite.Error(err, "Role %s should not be able to %s", tc.role, tc.operation)
		}
	}
}

func TestOrderIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(OrderIntegrationTestSuite))
}