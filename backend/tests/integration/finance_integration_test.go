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

type FinanceIntegrationTestSuite struct {
	suite.Suite
	db             *gorm.DB
	echo           *echo.Echo
	financeHandler *handler.FinanceHandler
	cleanup        func()
	
	// Test data
	company  *model.Company
	customer *model.Customer
	order    *model.Order
	finance  *model.Account
	manager  *model.Account
}

func (suite *FinanceIntegrationTestSuite) SetupSuite() {
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
		&model.Order{},
		&model.Invoice{},
		&model.Payment{},
		&model.Expense{},
	)
	suite.Require().NoError(err)
	
	// Setup services and handlers
	financeRepo := repository.NewFinanceRepository(db)
	financeService := service.NewFinanceService(financeRepo)
	suite.financeHandler = handler.NewFinanceHandler(financeService)
	
	// Setup Echo
	suite.echo = echo.New()
	
	// Setup cleanup function
	suite.cleanup = func() {
		// Clean up test data in correct order due to foreign key constraints
		db.Exec("DELETE FROM payments")
		db.Exec("DELETE FROM expenses")
		db.Exec("DELETE FROM invoices")
		db.Exec("DELETE FROM orders")
		db.Exec("DELETE FROM customers")
		db.Exec("DELETE FROM accounts")
		db.Exec("DELETE FROM companies")
	}
}

func (suite *FinanceIntegrationTestSuite) TearDownSuite() {
	if suite.cleanup != nil {
		suite.cleanup()
	}
}

func (suite *FinanceIntegrationTestSuite) SetupTest() {
	// Clean up before each test
	suite.cleanup()
	
	// Create test data
	suite.setupTestData()
}

func (suite *FinanceIntegrationTestSuite) setupTestData() {
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
	
	// Create test finance account
	suite.finance = &model.Account{
		CompanyID:    suite.company.ID,
		Username:     "finance1",
		Email:        "finance@test.com",
		PasswordHash: "$2a$10$XYZ...",
		FullName:     "Finance User",
		Role:         "finance",
		IsActive:     true,
	}
	err = suite.db.Create(suite.finance).Error
	suite.Require().NoError(err)
	
	// Create test manager account
	suite.manager = &model.Account{
		CompanyID:    suite.company.ID,
		Username:     "manager1",
		Email:        "manager@test.com",
		PasswordHash: "$2a$10$XYZ...",
		FullName:     "Manager User",
		Role:         "manager",
		IsActive:     true,
	}
	err = suite.db.Create(suite.manager).Error
	suite.Require().NoError(err)
	
	// Create test order
	suite.order = &model.Order{
		OrderNo:       "ORD-2024-001",
		CompanyID:     suite.company.ID,
		CustomerID:    suite.customer.ID,
		Status:        "completed",
		PaymentStatus: "pending",
		TotalAmount:   10000.00,
		Currency:      "USD",
	}
	err = suite.db.Create(suite.order).Error
	suite.Require().NoError(err)
}

func (suite *FinanceIntegrationTestSuite) TestCreateInvoice_Success() {
	createRequest := handler.CreateInvoiceRequest{
		OrderID:     suite.order.ID,
		InvoiceNo:   "INV-2024-001",
		Amount:      10000.00,
		Currency:    "USD",
		TaxRate:     0.05,
		TaxAmount:   500.00,
		TotalAmount: 10500.00,
		DueDate:     time.Now().AddDate(0, 1, 0),
		Notes:       "Integration test invoice",
	}
	
	requestBody, _ := json.Marshal(createRequest)
	req := httptest.NewRequest(http.MethodPost, "/finance/invoices", bytes.NewBuffer(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)
	
	// Mock JWT claims
	c.Set("user_id", suite.finance.ID.String())
	c.Set("company_id", suite.company.ID.String())
	c.Set("role", "finance")
	
	err := suite.financeHandler.CreateInvoice(c)
	suite.NoError(err)
	suite.Equal(http.StatusCreated, rec.Code)
	
	var response model.Invoice
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal(suite.order.ID, response.OrderID)
	suite.Equal("INV-2024-001", response.InvoiceNo)
	suite.Equal("pending", response.Status)
	suite.Equal(10000.00, response.Amount)
	suite.Equal(10500.00, response.TotalAmount)
	
	// Verify invoice was saved to database
	var savedInvoice model.Invoice
	err = suite.db.First(&savedInvoice, response.ID).Error
	suite.NoError(err)
	suite.Equal(response.ID, savedInvoice.ID)
}

func (suite *FinanceIntegrationTestSuite) TestPaymentWorkflow_Complete() {
	// First create an invoice
	invoice := &model.Invoice{
		OrderID:     suite.order.ID,
		CompanyID:   suite.company.ID,
		InvoiceNo:   "INV-2024-001",
		Status:      "pending",
		Amount:      10000.00,
		TaxAmount:   500.00,
		TotalAmount: 10500.00,
		PaidAmount:  0.00,
		Currency:    "USD",
		DueDate:     time.Now().AddDate(0, 1, 0),
	}
	err := suite.db.Create(invoice).Error
	suite.Require().NoError(err)
	
	// Process partial payment
	paymentRequest := handler.ProcessPaymentRequest{
		InvoiceID:     invoice.ID,
		Amount:        5000.00,
		PaymentMethod: "bank_transfer",
		PaymentDate:   time.Now(),
		Reference:     "TXN-123456",
		Notes:         "Partial payment",
	}
	
	requestBody, _ := json.Marshal(paymentRequest)
	req := httptest.NewRequest(http.MethodPost, "/finance/payments", bytes.NewBuffer(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)
	
	// Mock JWT claims
	c.Set("user_id", suite.finance.ID.String())
	c.Set("company_id", suite.company.ID.String())
	c.Set("role", "finance")
	
	err = suite.financeHandler.ProcessPayment(c)
	suite.NoError(err)
	suite.Equal(http.StatusCreated, rec.Code)
	
	var paymentResponse model.Payment
	err = json.Unmarshal(rec.Body.Bytes(), &paymentResponse)
	suite.NoError(err)
	suite.Equal(invoice.ID, paymentResponse.InvoiceID)
	suite.Equal(5000.00, paymentResponse.Amount)
	suite.Equal("completed", paymentResponse.Status)
	
	// Verify invoice status was updated
	var updatedInvoice model.Invoice
	err = suite.db.First(&updatedInvoice, invoice.ID).Error
	suite.NoError(err)
	suite.Equal("partial", updatedInvoice.Status)
	suite.Equal(5000.00, updatedInvoice.PaidAmount)
	
	// Process remaining payment
	finalPaymentRequest := handler.ProcessPaymentRequest{
		InvoiceID:     invoice.ID,
		Amount:        5500.00, // Remaining amount including tax
		PaymentMethod: "bank_transfer",
		PaymentDate:   time.Now(),
		Reference:     "TXN-789012",
		Notes:         "Final payment",
	}
	
	requestBody, _ = json.Marshal(finalPaymentRequest)
	req = httptest.NewRequest(http.MethodPost, "/finance/payments", bytes.NewBuffer(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	
	rec = httptest.NewRecorder()
	c = suite.echo.NewContext(req, rec)
	
	// Mock JWT claims
	c.Set("user_id", suite.finance.ID.String())
	c.Set("company_id", suite.company.ID.String())
	c.Set("role", "finance")
	
	err = suite.financeHandler.ProcessPayment(c)
	suite.NoError(err)
	suite.Equal(http.StatusCreated, rec.Code)
	
	// Verify invoice is now fully paid
	err = suite.db.First(&updatedInvoice, invoice.ID).Error
	suite.NoError(err)
	suite.Equal("paid", updatedInvoice.Status)
	suite.Equal(10500.00, updatedInvoice.PaidAmount)
}

func (suite *FinanceIntegrationTestSuite) TestExpenseApprovalWorkflow() {
	// Create expense
	expenseRequest := handler.CreateExpenseRequest{
		Category:    "travel",
		Amount:      500.00,
		Currency:    "USD",
		Description: "Business trip to customer site",
		ExpenseDate: time.Now(),
		ReceiptURL:  "https://example.com/receipt.jpg",
	}
	
	requestBody, _ := json.Marshal(expenseRequest)
	req := httptest.NewRequest(http.MethodPost, "/finance/expenses", bytes.NewBuffer(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)
	
	// Mock JWT claims - employee submitting expense
	c.Set("user_id", suite.finance.ID.String())
	c.Set("company_id", suite.company.ID.String())
	c.Set("role", "employee")
	
	err := suite.financeHandler.CreateExpense(c)
	suite.NoError(err)
	suite.Equal(http.StatusCreated, rec.Code)
	
	var expenseResponse model.Expense
	err = json.Unmarshal(rec.Body.Bytes(), &expenseResponse)
	suite.NoError(err)
	suite.Equal("travel", expenseResponse.Category)
	suite.Equal(500.00, expenseResponse.Amount)
	suite.Equal("pending", expenseResponse.Status)
	
	// Manager approves expense
	approvalRequest := handler.ExpenseApprovalRequest{
		Notes: "Approved for reimbursement",
	}
	
	requestBody, _ = json.Marshal(approvalRequest)
	req = httptest.NewRequest(http.MethodPost, "/finance/expenses/"+expenseResponse.ID.String()+"/approve", bytes.NewBuffer(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	
	rec = httptest.NewRecorder()
	c = suite.echo.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(expenseResponse.ID.String())
	
	// Mock JWT claims - manager approving
	c.Set("user_id", suite.manager.ID.String())
	c.Set("company_id", suite.company.ID.String())
	c.Set("role", "manager")
	
	err = suite.financeHandler.ApproveExpense(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
	
	// Verify expense was approved in database
	var updatedExpense model.Expense
	err = suite.db.First(&updatedExpense, expenseResponse.ID).Error
	suite.NoError(err)
	suite.Equal("approved", updatedExpense.Status)
	suite.Equal(&suite.manager.ID, updatedExpense.ApprovedBy)
}

func (suite *FinanceIntegrationTestSuite) TestListInvoicesWithFilters() {
	// Create multiple invoices
	invoices := []*model.Invoice{
		{
			OrderID:     suite.order.ID,
			CompanyID:   suite.company.ID,
			InvoiceNo:   "INV-2024-001",
			Status:      "pending",
			Amount:      10000.00,
			Currency:    "USD",
		},
		{
			OrderID:     suite.order.ID,
			CompanyID:   suite.company.ID,
			InvoiceNo:   "INV-2024-002",
			Status:      "paid",
			Amount:      5000.00,
			Currency:    "USD",
		},
		{
			OrderID:     suite.order.ID,
			CompanyID:   suite.company.ID,
			InvoiceNo:   "INV-2024-003",
			Status:      "pending",
			Amount:      8000.00,
			Currency:    "USD",
		},
	}
	
	for _, invoice := range invoices {
		err := suite.db.Create(invoice).Error
		suite.Require().NoError(err)
	}
	
	// Test filtering by status
	req := httptest.NewRequest(http.MethodGet, "/finance/invoices?status=pending", nil)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)
	
	// Mock JWT claims
	c.Set("user_id", suite.finance.ID.String())
	c.Set("company_id", suite.company.ID.String())
	c.Set("role", "finance")
	
	err := suite.financeHandler.ListInvoices(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
	
	var response []model.Invoice
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Len(response, 2) // Should return 2 pending invoices
	
	for _, invoice := range response {
		suite.Equal("pending", invoice.Status)
	}
}

func (suite *FinanceIntegrationTestSuite) TestGetReceivables_AgingAnalysis() {
	// Create invoices with different due dates
	now := time.Now()
	invoices := []*model.Invoice{
		{
			OrderID:     suite.order.ID,
			CompanyID:   suite.company.ID,
			InvoiceNo:   "INV-2024-001",
			Status:      "pending",
			Amount:      10000.00,
			TotalAmount: 10000.00,
			PaidAmount:  0.00,
			Currency:    "USD",
			DueDate:     now.AddDate(0, 0, 10), // Due in 10 days - current
		},
		{
			OrderID:     suite.order.ID,
			CompanyID:   suite.company.ID,
			InvoiceNo:   "INV-2024-002",
			Status:      "overdue",
			Amount:      5000.00,
			TotalAmount: 5000.00,
			PaidAmount:  0.00,
			Currency:    "USD",
			DueDate:     now.AddDate(0, 0, -15), // 15 days overdue
		},
		{
			OrderID:     suite.order.ID,
			CompanyID:   suite.company.ID,
			InvoiceNo:   "INV-2024-003",
			Status:      "overdue",
			Amount:      3000.00,
			TotalAmount: 3000.00,
			PaidAmount:  0.00,
			Currency:    "USD",
			DueDate:     now.AddDate(0, 0, -45), // 45 days overdue
		},
	}
	
	for _, invoice := range invoices {
		err := suite.db.Create(invoice).Error
		suite.Require().NoError(err)
	}
	
	// Get receivables
	req := httptest.NewRequest(http.MethodGet, "/finance/receivables", nil)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)
	
	// Mock JWT claims
	c.Set("user_id", suite.finance.ID.String())
	c.Set("company_id", suite.company.ID.String())
	c.Set("role", "finance")
	
	err := suite.financeHandler.GetReceivables(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
	
	var response []model.Receivable
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Len(response, 3)
	
	// Verify aging analysis
	currentCount := 0
	overdueCount := 0
	
	for _, receivable := range response {
		if receivable.Status == "outstanding" {
			currentCount++
		} else if receivable.Status == "overdue" {
			overdueCount++
			suite.True(receivable.DaysOverdue > 0)
		}
	}
	
	suite.Equal(1, currentCount)
	suite.Equal(2, overdueCount)
}

func (suite *FinanceIntegrationTestSuite) TestFinancialSummary_Integration() {
	// Create test data for financial summary
	
	// Create invoices
	invoices := []*model.Invoice{
		{
			OrderID:     suite.order.ID,
			CompanyID:   suite.company.ID,
			InvoiceNo:   "INV-2024-001",
			Status:      "paid",
			Amount:      10000.00,
			TotalAmount: 10500.00,
			PaidAmount:  10500.00,
			Currency:    "USD",
		},
		{
			OrderID:     suite.order.ID,
			CompanyID:   suite.company.ID,
			InvoiceNo:   "INV-2024-002",
			Status:      "pending",
			Amount:      5000.00,
			TotalAmount: 5250.00,
			PaidAmount:  0.00,
			Currency:    "USD",
		},
	}
	
	for _, invoice := range invoices {
		err := suite.db.Create(invoice).Error
		suite.Require().NoError(err)
	}
	
	// Create expenses
	expenses := []*model.Expense{
		{
			CompanyID:   suite.company.ID,
			Category:    "travel",
			Amount:      500.00,
			Currency:    "USD",
			Status:      "approved",
			SubmittedBy: suite.finance.ID,
		},
		{
			CompanyID:   suite.company.ID,
			Category:    "office",
			Amount:      200.00,
			Currency:    "USD",
			Status:      "pending",
			SubmittedBy: suite.finance.ID,
		},
	}
	
	for _, expense := range expenses {
		err := suite.db.Create(expense).Error
		suite.Require().NoError(err)
	}
	
	// Get financial summary
	req := httptest.NewRequest(http.MethodGet, "/finance/summary?date_range=current_month", nil)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)
	
	// Mock JWT claims
	c.Set("user_id", suite.finance.ID.String())
	c.Set("company_id", suite.company.ID.String())
	c.Set("role", "finance")
	
	err := suite.financeHandler.GetFinancialSummary(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
	
	var response service.FinancialSummary
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	suite.NoError(err)
	
	// Verify summary calculations
	suite.Equal(15750.00, response.Revenue.InvoicedAmount) // 10500 + 5250
	suite.Equal(10500.00, response.Revenue.PaidAmount)     // Only paid invoice
	suite.Equal(5250.00, response.Revenue.OutstandingAmount) // Pending invoice
	
	suite.Equal(700.00, response.Expenses.TotalExpenses)     // 500 + 200
	suite.Equal(500.00, response.Expenses.ApprovedExpenses) // Only approved
	suite.Equal(200.00, response.Expenses.PendingExpenses)  // Only pending
	
	suite.Equal(9800.00, response.CashFlow.NetCashFlow) // 10500 - 700
}

func (suite *FinanceIntegrationTestSuite) TestFinancePermissions() {
	// Create test invoice
	invoice := &model.Invoice{
		OrderID:   suite.order.ID,
		CompanyID: suite.company.ID,
		InvoiceNo: "INV-2024-001",
		Status:    "pending",
		Amount:    10000.00,
		Currency:  "USD",
	}
	err := suite.db.Create(invoice).Error
	suite.Require().NoError(err)
	
	// Test permissions for different roles
	testCases := []struct {
		role           string
		operation      string
		shouldAccess   bool
		expectedStatus int
	}{
		{"admin", "create_invoice", true, http.StatusCreated},
		{"finance", "create_invoice", true, http.StatusCreated},
		{"manager", "create_invoice", true, http.StatusCreated},
		{"sales", "create_invoice", false, http.StatusForbidden},
		{"engineer", "create_invoice", false, http.StatusForbidden},
		{"admin", "view_summary", true, http.StatusOK},
		{"finance", "view_summary", true, http.StatusOK},
		{"manager", "view_summary", true, http.StatusOK},
		{"sales", "view_summary", false, http.StatusForbidden},
		{"engineer", "view_summary", false, http.StatusForbidden},
	}
	
	for _, tc := range testCases {
		var req *http.Request
		var expectedMethod string
		
		if tc.operation == "create_invoice" {
			createRequest := handler.CreateInvoiceRequest{
				OrderID:   suite.order.ID,
				InvoiceNo: "INV-TEST-001",
				Amount:    1000.00,
				Currency:  "USD",
			}
			requestBody, _ := json.Marshal(createRequest)
			req = httptest.NewRequest(http.MethodPost, "/finance/invoices", bytes.NewBuffer(requestBody))
			expectedMethod = "CreateInvoice"
		} else if tc.operation == "view_summary" {
			req = httptest.NewRequest(http.MethodGet, "/finance/summary", nil)
			expectedMethod = "GetFinancialSummary"
		}
		
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := suite.echo.NewContext(req, rec)
		
		// Mock JWT claims with different roles
		c.Set("user_id", suite.finance.ID.String())
		c.Set("company_id", suite.company.ID.String())
		c.Set("role", tc.role)
		
		var err error
		if expectedMethod == "CreateInvoice" {
			err = suite.financeHandler.CreateInvoice(c)
		} else if expectedMethod == "GetFinancialSummary" {
			err = suite.financeHandler.GetFinancialSummary(c)
		}
		
		if tc.shouldAccess {
			suite.NoError(err, "Role %s should be able to %s", tc.role, tc.operation)
			suite.Equal(tc.expectedStatus, rec.Code)
		} else {
			suite.Error(err, "Role %s should not be able to %s", tc.role, tc.operation)
		}
	}
}

func TestFinanceIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(FinanceIntegrationTestSuite))
}