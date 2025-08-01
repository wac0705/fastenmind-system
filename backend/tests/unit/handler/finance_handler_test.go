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

// MockFinanceService is a mock implementation of FinanceService
type MockFinanceService struct {
	mock.Mock
}

func (m *MockFinanceService) CreateInvoice(ctx context.Context, request *service.CreateInvoiceRequest, companyID uuid.UUID) (*model.Invoice, error) {
	args := m.Called(ctx, request, companyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Invoice), args.Error(1)
}

func (m *MockFinanceService) GetInvoice(ctx context.Context, id uuid.UUID) (*model.Invoice, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Invoice), args.Error(1)
}

func (m *MockFinanceService) ListInvoices(ctx context.Context, filter map[string]interface{}) ([]*model.Invoice, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Invoice), args.Error(1)
}

func (m *MockFinanceService) ProcessPayment(ctx context.Context, request *service.ProcessPaymentRequest, companyID uuid.UUID) (*model.Payment, error) {
	args := m.Called(ctx, request, companyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Payment), args.Error(1)
}

func (m *MockFinanceService) ListPayments(ctx context.Context, filter map[string]interface{}) ([]*model.Payment, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Payment), args.Error(1)
}

func (m *MockFinanceService) CreateExpense(ctx context.Context, request *service.CreateExpenseRequest, companyID uuid.UUID, userID uuid.UUID) (*model.Expense, error) {
	args := m.Called(ctx, request, companyID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Expense), args.Error(1)
}

func (m *MockFinanceService) ApproveExpense(ctx context.Context, expenseID uuid.UUID, approverID uuid.UUID, notes string) error {
	args := m.Called(ctx, expenseID, approverID, notes)
	return args.Error(0)
}

func (m *MockFinanceService) RejectExpense(ctx context.Context, expenseID uuid.UUID, approverID uuid.UUID, notes string) error {
	args := m.Called(ctx, expenseID, approverID, notes)
	return args.Error(0)
}

func (m *MockFinanceService) ListExpenses(ctx context.Context, filter map[string]interface{}) ([]*model.Expense, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Expense), args.Error(1)
}

func (m *MockFinanceService) GetReceivables(ctx context.Context, filter map[string]interface{}) ([]*model.Receivable, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Receivable), args.Error(1)
}

func (m *MockFinanceService) GetFinancialSummary(ctx context.Context, companyID uuid.UUID, dateRange string) (*service.FinancialSummary, error) {
	args := m.Called(ctx, companyID, dateRange)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.FinancialSummary), args.Error(1)
}

func TestFinanceHandler_CreateInvoice_Success(t *testing.T) {
	// Arrange
	e := echo.New()
	mockFinanceService := new(MockFinanceService)
	financeHandler := handler.NewFinanceHandler(mockFinanceService)
	
	orderID := uuid.New()
	companyID := uuid.New()
	
	createRequest := handler.CreateInvoiceRequest{
		OrderID:     orderID,
		InvoiceNo:   "INV-2024-001",
		Amount:      10000.00,
		Currency:    "USD",
		TaxRate:     0.05,
		TaxAmount:   500.00,
		TotalAmount: 10500.00,
		DueDate:     time.Now().AddDate(0, 1, 0),
		Notes:       "Test invoice",
	}
	
	expectedInvoice := &model.Invoice{
		ID:          uuid.New(),
		OrderID:     orderID,
		InvoiceNo:   "INV-2024-001",
		Status:      "pending",
		Amount:      10000.00,
		TaxAmount:   500.00,
		TotalAmount: 10500.00,
		Currency:    "USD",
		DueDate:     createRequest.DueDate,
		Notes:       "Test invoice",
	}
	
	mockFinanceService.On("CreateInvoice", mock.Anything, mock.AnythingOfType("*service.CreateInvoiceRequest"), companyID).Return(expectedInvoice, nil)
	
	requestBody, _ := json.Marshal(createRequest)
	req := httptest.NewRequest(http.MethodPost, "/finance/invoices", bytes.NewBuffer(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	
	// Mock JWT claims
	c.Set("user_id", uuid.New().String())
	c.Set("company_id", companyID.String())
	c.Set("role", "finance")
	
	// Act
	err := financeHandler.CreateInvoice(c)
	
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	
	var response model.Invoice
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedInvoice.InvoiceNo, response.InvoiceNo)
	assert.Equal(t, expectedInvoice.Amount, response.Amount)
	
	mockFinanceService.AssertExpectations(t)
}

func TestFinanceHandler_ProcessPayment_Success(t *testing.T) {
	// Arrange
	e := echo.New()
	mockFinanceService := new(MockFinanceService)
	financeHandler := handler.NewFinanceHandler(mockFinanceService)
	
	invoiceID := uuid.New()
	companyID := uuid.New()
	
	paymentRequest := handler.ProcessPaymentRequest{
		InvoiceID:     invoiceID,
		Amount:        5000.00,
		PaymentMethod: "bank_transfer",
		PaymentDate:   time.Now(),
		Reference:     "TXN-123456",
		Notes:         "Partial payment",
	}
	
	expectedPayment := &model.Payment{
		ID:            uuid.New(),
		InvoiceID:     invoiceID,
		Amount:        5000.00,
		Currency:      "USD",
		PaymentMethod: "bank_transfer",
		Status:        "completed",
		Reference:     "TXN-123456",
	}
	
	mockFinanceService.On("ProcessPayment", mock.Anything, mock.AnythingOfType("*service.ProcessPaymentRequest"), companyID).Return(expectedPayment, nil)
	
	requestBody, _ := json.Marshal(paymentRequest)
	req := httptest.NewRequest(http.MethodPost, "/finance/payments", bytes.NewBuffer(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	
	// Mock JWT claims
	c.Set("user_id", uuid.New().String())
	c.Set("company_id", companyID.String())
	c.Set("role", "finance")
	
	// Act
	err := financeHandler.ProcessPayment(c)
	
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	
	var response model.Payment
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedPayment.Amount, response.Amount)
	assert.Equal(t, expectedPayment.PaymentMethod, response.PaymentMethod)
	
	mockFinanceService.AssertExpectations(t)
}

func TestFinanceHandler_CreateExpense_Success(t *testing.T) {
	// Arrange
	e := echo.New()
	mockFinanceService := new(MockFinanceService)
	financeHandler := handler.NewFinanceHandler(mockFinanceService)
	
	companyID := uuid.New()
	userID := uuid.New()
	
	expenseRequest := handler.CreateExpenseRequest{
		Category:    "travel",
		Amount:      500.00,
		Currency:    "USD",
		Description: "Business trip to customer site",
		ExpenseDate: time.Now(),
		ReceiptURL:  "https://example.com/receipt.jpg",
	}
	
	expectedExpense := &model.Expense{
		ID:          uuid.New(),
		Category:    "travel",
		Amount:      500.00,
		Currency:    "USD",
		Description: "Business trip to customer site",
		Status:      "pending",
		SubmittedBy: userID,
	}
	
	mockFinanceService.On("CreateExpense", mock.Anything, mock.AnythingOfType("*service.CreateExpenseRequest"), companyID, userID).Return(expectedExpense, nil)
	
	requestBody, _ := json.Marshal(expenseRequest)
	req := httptest.NewRequest(http.MethodPost, "/finance/expenses", bytes.NewBuffer(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	
	// Mock JWT claims
	c.Set("user_id", userID.String())
	c.Set("company_id", companyID.String())
	c.Set("role", "employee")
	
	// Act
	err := financeHandler.CreateExpense(c)
	
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	
	var response model.Expense
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedExpense.Category, response.Category)
	assert.Equal(t, expectedExpense.Amount, response.Amount)
	
	mockFinanceService.AssertExpectations(t)
}

func TestFinanceHandler_ApproveExpense_Success(t *testing.T) {
	// Arrange
	e := echo.New()
	mockFinanceService := new(MockFinanceService)
	financeHandler := handler.NewFinanceHandler(mockFinanceService)
	
	expenseID := uuid.New()
	approverID := uuid.New()
	
	approvalRequest := handler.ExpenseApprovalRequest{
		Notes: "Approved for reimbursement",
	}
	
	mockFinanceService.On("ApproveExpense", mock.Anything, expenseID, approverID, "Approved for reimbursement").Return(nil)
	
	requestBody, _ := json.Marshal(approvalRequest)
	req := httptest.NewRequest(http.MethodPost, "/finance/expenses/"+expenseID.String()+"/approve", bytes.NewBuffer(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(expenseID.String())
	
	// Mock JWT claims - manager can approve
	c.Set("user_id", approverID.String())
	c.Set("company_id", uuid.New().String())
	c.Set("role", "manager")
	
	// Act
	err := financeHandler.ApproveExpense(c)
	
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	
	mockFinanceService.AssertExpectations(t)
}

func TestFinanceHandler_RejectExpense_Success(t *testing.T) {
	// Arrange
	e := echo.New()
	mockFinanceService := new(MockFinanceService)
	financeHandler := handler.NewFinanceHandler(mockFinanceService)
	
	expenseID := uuid.New()
	approverID := uuid.New()
	
	rejectionRequest := handler.ExpenseApprovalRequest{
		Notes: "Insufficient documentation",
	}
	
	mockFinanceService.On("RejectExpense", mock.Anything, expenseID, approverID, "Insufficient documentation").Return(nil)
	
	requestBody, _ := json.Marshal(rejectionRequest)
	req := httptest.NewRequest(http.MethodPost, "/finance/expenses/"+expenseID.String()+"/reject", bytes.NewBuffer(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(expenseID.String())
	
	// Mock JWT claims - manager can reject
	c.Set("user_id", approverID.String())
	c.Set("company_id", uuid.New().String())
	c.Set("role", "manager")
	
	// Act
	err := financeHandler.RejectExpense(c)
	
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	
	mockFinanceService.AssertExpectations(t)
}

func TestFinanceHandler_GetReceivables_Success(t *testing.T) {
	// Arrange
	e := echo.New()
	mockFinanceService := new(MockFinanceService)
	financeHandler := handler.NewFinanceHandler(mockFinanceService)
	
	companyID := uuid.New()
	
	expectedReceivables := []*model.Receivable{
		{
			ID:         uuid.New(),
			CustomerID: uuid.New(),
			Amount:     5000.00,
			Currency:   "USD",
			Status:     "outstanding",
		},
	}
	
	filter := map[string]interface{}{
		"company_id": companyID.String(),
	}
	
	mockFinanceService.On("GetReceivables", mock.Anything, filter).Return(expectedReceivables, nil)
	
	req := httptest.NewRequest(http.MethodGet, "/finance/receivables", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	
	// Mock JWT claims
	c.Set("user_id", uuid.New().String())
	c.Set("company_id", companyID.String())
	c.Set("role", "finance")
	
	// Act
	err := financeHandler.GetReceivables(c)
	
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	
	var response []model.Receivable
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, expectedReceivables[0].Amount, response[0].Amount)
	
	mockFinanceService.AssertExpectations(t)
}

func TestFinanceHandler_GetFinancialSummary_Success(t *testing.T) {
	// Arrange
	e := echo.New()
	mockFinanceService := new(MockFinanceService)
	financeHandler := handler.NewFinanceHandler(mockFinanceService)
	
	companyID := uuid.New()
	dateRange := "current_month"
	
	expectedSummary := &service.FinancialSummary{
		Revenue: service.RevenueSummary{
			TotalRevenue:      100000.00,
			InvoicedAmount:    90000.00,
			PaidAmount:        80000.00,
			OutstandingAmount: 10000.00,
		},
		Expenses: service.ExpenseSummary{
			TotalExpenses:    25000.00,
			PendingExpenses:  5000.00,
			ApprovedExpenses: 20000.00,
		},
		CashFlow: service.CashFlowSummary{
			NetCashFlow: 55000.00,
			CashInflow:  80000.00,
			CashOutflow: 25000.00,
		},
	}
	
	mockFinanceService.On("GetFinancialSummary", mock.Anything, companyID, dateRange).Return(expectedSummary, nil)
	
	req := httptest.NewRequest(http.MethodGet, "/finance/summary?date_range=current_month", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	
	// Mock JWT claims
	c.Set("user_id", uuid.New().String())
	c.Set("company_id", companyID.String())
	c.Set("role", "finance")
	
	// Act
	err := financeHandler.GetFinancialSummary(c)
	
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	
	var response service.FinancialSummary
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedSummary.Revenue.TotalRevenue, response.Revenue.TotalRevenue)
	assert.Equal(t, expectedSummary.CashFlow.NetCashFlow, response.CashFlow.NetCashFlow)
	
	mockFinanceService.AssertExpectations(t)
}

func TestFinanceHandler_ListInvoices_WithFilters(t *testing.T) {
	// Arrange
	e := echo.New()
	mockFinanceService := new(MockFinanceService)
	financeHandler := handler.NewFinanceHandler(mockFinanceService)
	
	companyID := uuid.New()
	
	expectedInvoices := []*model.Invoice{
		{
			ID:        uuid.New(),
			InvoiceNo: "INV-2024-001",
			Status:    "pending",
			Amount:    10000.00,
		},
	}
	
	filter := map[string]interface{}{
		"company_id": companyID.String(),
		"status":     "pending",
	}
	
	mockFinanceService.On("ListInvoices", mock.Anything, filter).Return(expectedInvoices, nil)
	
	req := httptest.NewRequest(http.MethodGet, "/finance/invoices?status=pending", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	
	// Mock JWT claims
	c.Set("user_id", uuid.New().String())
	c.Set("company_id", companyID.String())
	c.Set("role", "finance")
	
	// Act
	err := financeHandler.ListInvoices(c)
	
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	
	var response []model.Invoice
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, "pending", response[0].Status)
	
	mockFinanceService.AssertExpectations(t)
}

func TestFinanceHandler_PermissionCheck(t *testing.T) {
	// Arrange
	e := echo.New()
	mockFinanceService := new(MockFinanceService)
	financeHandler := handler.NewFinanceHandler(mockFinanceService)
	
	// Test cases for different roles and operations
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
		{"admin", "approve_expense", true, http.StatusOK},
		{"manager", "approve_expense", true, http.StatusOK},
		{"finance", "approve_expense", true, http.StatusOK},
		{"employee", "approve_expense", false, http.StatusForbidden},
		{"sales", "approve_expense", false, http.StatusForbidden},
	}
	
	for _, tc := range testCases {
		// Create appropriate request based on operation
		var req *http.Request
		var expectedMethod string
		
		if tc.operation == "create_invoice" {
			createRequest := handler.CreateInvoiceRequest{
				OrderID:   uuid.New(),
				InvoiceNo: "INV-TEST-001",
				Amount:    1000.00,
				Currency:  "USD",
			}
			requestBody, _ := json.Marshal(createRequest)
			req = httptest.NewRequest(http.MethodPost, "/finance/invoices", bytes.NewBuffer(requestBody))
			expectedMethod = "CreateInvoice"
		} else if tc.operation == "approve_expense" {
			approvalRequest := handler.ExpenseApprovalRequest{
				Notes: "Test approval",
			}
			requestBody, _ := json.Marshal(approvalRequest)
			req = httptest.NewRequest(http.MethodPost, "/finance/expenses/"+uuid.New().String()+"/approve", bytes.NewBuffer(requestBody))
			expectedMethod = "ApproveExpense"
		}
		
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		
		if tc.operation == "approve_expense" {
			c.SetParamNames("id")
			c.SetParamValues(uuid.New().String())
		}
		
		// Mock JWT claims with different roles
		c.Set("user_id", uuid.New().String())
		c.Set("company_id", uuid.New().String())
		c.Set("role", tc.role)
		
		// Mock service calls if access should be granted
		if tc.shouldAccess {
			if expectedMethod == "CreateInvoice" {
				expectedInvoice := &model.Invoice{
					ID:        uuid.New(),
					InvoiceNo: "INV-TEST-001",
					Amount:    1000.00,
					Status:    "pending",
				}
				mockFinanceService.On("CreateInvoice", mock.Anything, mock.AnythingOfType("*service.CreateInvoiceRequest"), mock.AnythingOfType("uuid.UUID")).Return(expectedInvoice, nil).Once()
			} else if expectedMethod == "ApproveExpense" {
				mockFinanceService.On("ApproveExpense", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("string")).Return(nil).Once()
			}
		}
		
		var err error
		if expectedMethod == "CreateInvoice" {
			err = financeHandler.CreateInvoice(c)
		} else if expectedMethod == "ApproveExpense" {
			err = financeHandler.ApproveExpense(c)
		}
		
		if tc.shouldAccess {
			assert.NoError(t, err, "Role %s should be able to %s", tc.role, tc.operation)
			assert.Equal(t, tc.expectedStatus, rec.Code)
		} else {
			assert.Error(t, err, "Role %s should not be able to %s", tc.role, tc.operation)
		}
	}
}