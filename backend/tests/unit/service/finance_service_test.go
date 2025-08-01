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

// MockFinanceRepository is a mock implementation of FinanceRepository
type MockFinanceRepository struct {
	mock.Mock
}

func (m *MockFinanceRepository) CreateInvoice(ctx context.Context, invoice *model.Invoice) error {
	args := m.Called(ctx, invoice)
	return args.Error(0)
}

func (m *MockFinanceRepository) GetInvoice(ctx context.Context, id uuid.UUID) (*model.Invoice, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Invoice), args.Error(1)
}

func (m *MockFinanceRepository) UpdateInvoice(ctx context.Context, invoice *model.Invoice) error {
	args := m.Called(ctx, invoice)
	return args.Error(0)
}

func (m *MockFinanceRepository) ListInvoices(ctx context.Context, filter map[string]interface{}) ([]*model.Invoice, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Invoice), args.Error(1)
}

func (m *MockFinanceRepository) CreatePayment(ctx context.Context, payment *model.Payment) error {
	args := m.Called(ctx, payment)
	return args.Error(0)
}

func (m *MockFinanceRepository) GetPayment(ctx context.Context, id uuid.UUID) (*model.Payment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Payment), args.Error(1)
}

func (m *MockFinanceRepository) ListPayments(ctx context.Context, filter map[string]interface{}) ([]*model.Payment, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Payment), args.Error(1)
}

func (m *MockFinanceRepository) CreateExpense(ctx context.Context, expense *model.Expense) error {
	args := m.Called(ctx, expense)
	return args.Error(0)
}

func (m *MockFinanceRepository) GetExpense(ctx context.Context, id uuid.UUID) (*model.Expense, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Expense), args.Error(1)
}

func (m *MockFinanceRepository) UpdateExpense(ctx context.Context, expense *model.Expense) error {
	args := m.Called(ctx, expense)
	return args.Error(0)
}

func (m *MockFinanceRepository) ListExpenses(ctx context.Context, filter map[string]interface{}) ([]*model.Expense, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Expense), args.Error(1)
}

func (m *MockFinanceRepository) GetReceivables(ctx context.Context, filter map[string]interface{}) ([]*model.Receivable, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Receivable), args.Error(1)
}

func (m *MockFinanceRepository) GetFinancialSummary(ctx context.Context, companyID uuid.UUID, dateRange string) (*service.FinancialSummary, error) {
	args := m.Called(ctx, companyID, dateRange)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.FinancialSummary), args.Error(1)
}

func TestFinanceService_CreateInvoice_Success(t *testing.T) {
	// Arrange
	mockFinanceRepo := new(MockFinanceRepository)
	financeService := service.NewFinanceService(mockFinanceRepo)
	
	orderID := uuid.New()
	companyID := uuid.New()
	
	createRequest := &service.CreateInvoiceRequest{
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
	
	mockFinanceRepo.On("CreateInvoice", mock.Anything, mock.AnythingOfType("*model.Invoice")).Return(nil)
	
	// Act
	result, err := financeService.CreateInvoice(context.Background(), createRequest, companyID)
	
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, orderID, result.OrderID)
	assert.Equal(t, "INV-2024-001", result.InvoiceNo)
	assert.Equal(t, "pending", result.Status)
	assert.Equal(t, 10000.00, result.Amount)
	assert.Equal(t, 10500.00, result.TotalAmount)
	
	mockFinanceRepo.AssertExpectations(t)
}

func TestFinanceService_ProcessPayment_Success(t *testing.T) {
	// Arrange
	mockFinanceRepo := new(MockFinanceRepository)
	financeService := service.NewFinanceService(mockFinanceRepo)
	
	invoiceID := uuid.New()
	companyID := uuid.New()
	
	// Mock existing invoice
	invoice := &model.Invoice{
		ID:          invoiceID,
		InvoiceNo:   "INV-2024-001",
		Status:      "pending",
		TotalAmount: 10500.00,
		PaidAmount:  0.00,
		Currency:    "USD",
	}
	
	paymentRequest := &service.ProcessPaymentRequest{
		InvoiceID:     invoiceID,
		Amount:        5000.00,
		PaymentMethod: "bank_transfer",
		PaymentDate:   time.Now(),
		Reference:     "TXN-123456",
		Notes:         "Partial payment",
	}
	
	mockFinanceRepo.On("GetInvoice", mock.Anything, invoiceID).Return(invoice, nil)
	mockFinanceRepo.On("CreatePayment", mock.Anything, mock.AnythingOfType("*model.Payment")).Return(nil)
	mockFinanceRepo.On("UpdateInvoice", mock.Anything, mock.AnythingOfType("*model.Invoice")).Return(nil)
	
	// Act
	result, err := financeService.ProcessPayment(context.Background(), paymentRequest, companyID)
	
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, invoiceID, result.InvoiceID)
	assert.Equal(t, 5000.00, result.Amount)
	assert.Equal(t, "bank_transfer", result.PaymentMethod)
	assert.Equal(t, "completed", result.Status)
	
	mockFinanceRepo.AssertExpectations(t)
}

func TestFinanceService_ProcessPayment_OverPayment(t *testing.T) {
	// Arrange
	mockFinanceRepo := new(MockFinanceRepository)
	financeService := service.NewFinanceService(mockFinanceRepo)
	
	invoiceID := uuid.New()
	companyID := uuid.New()
	
	// Mock existing invoice with some paid amount
	invoice := &model.Invoice{
		ID:          invoiceID,
		InvoiceNo:   "INV-2024-001",
		Status:      "partial",
		TotalAmount: 10500.00,
		PaidAmount:  8000.00, // Already paid 8000
		Currency:    "USD",
	}
	
	// Try to pay more than remaining amount
	paymentRequest := &service.ProcessPaymentRequest{
		InvoiceID:     invoiceID,
		Amount:        5000.00, // This would exceed total amount
		PaymentMethod: "bank_transfer",
		PaymentDate:   time.Now(),
	}
	
	mockFinanceRepo.On("GetInvoice", mock.Anything, invoiceID).Return(invoice, nil)
	
	// Act
	result, err := financeService.ProcessPayment(context.Background(), paymentRequest, companyID)
	
	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "payment amount exceeds remaining balance")
	
	mockFinanceRepo.AssertExpectations(t)
}

func TestFinanceService_CreateExpense_Success(t *testing.T) {
	// Arrange
	mockFinanceRepo := new(MockFinanceRepository)
	financeService := service.NewFinanceService(mockFinanceRepo)
	
	companyID := uuid.New()
	userID := uuid.New()
	
	createRequest := &service.CreateExpenseRequest{
		Category:    "travel",
		Amount:      500.00,
		Currency:    "USD",
		Description: "Business trip to customer site",
		ExpenseDate: time.Now(),
		ReceiptURL:  "https://example.com/receipt.jpg",
	}
	
	mockFinanceRepo.On("CreateExpense", mock.Anything, mock.AnythingOfType("*model.Expense")).Return(nil)
	
	// Act
	result, err := financeService.CreateExpense(context.Background(), createRequest, companyID, userID)
	
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "travel", result.Category)
	assert.Equal(t, 500.00, result.Amount)
	assert.Equal(t, "pending", result.Status)
	assert.Equal(t, userID, result.SubmittedBy)
	
	mockFinanceRepo.AssertExpectations(t)
}

func TestFinanceService_ApproveExpense_Success(t *testing.T) {
	// Arrange
	mockFinanceRepo := new(MockFinanceRepository)
	financeService := service.NewFinanceService(mockFinanceRepo)
	
	expenseID := uuid.New()
	approverID := uuid.New()
	
	expense := &model.Expense{
		ID:          expenseID,
		Status:      "pending",
		Amount:      500.00,
		Category:    "travel",
		Description: "Business trip",
	}
	
	mockFinanceRepo.On("GetExpense", mock.Anything, expenseID).Return(expense, nil)
	mockFinanceRepo.On("UpdateExpense", mock.Anything, mock.AnythingOfType("*model.Expense")).Return(nil)
	
	// Act
	err := financeService.ApproveExpense(context.Background(), expenseID, approverID, "Approved for reimbursement")
	
	// Assert
	assert.NoError(t, err)
	
	mockFinanceRepo.AssertExpectations(t)
}

func TestFinanceService_RejectExpense_Success(t *testing.T) {
	// Arrange
	mockFinanceRepo := new(MockFinanceRepository)
	financeService := service.NewFinanceService(mockFinanceRepo)
	
	expenseID := uuid.New()
	approverID := uuid.New()
	
	expense := &model.Expense{
		ID:          expenseID,
		Status:      "pending",
		Amount:      500.00,
		Category:    "travel",
		Description: "Business trip",
	}
	
	mockFinanceRepo.On("GetExpense", mock.Anything, expenseID).Return(expense, nil)
	mockFinanceRepo.On("UpdateExpense", mock.Anything, mock.AnythingOfType("*model.Expense")).Return(nil)
	
	// Act
	err := financeService.RejectExpense(context.Background(), expenseID, approverID, "Insufficient documentation")
	
	// Assert
	assert.NoError(t, err)
	
	mockFinanceRepo.AssertExpectations(t)
}

func TestFinanceService_GetReceivables_Success(t *testing.T) {
	// Arrange
	mockFinanceRepo := new(MockFinanceRepository)
	financeService := service.NewFinanceService(mockFinanceRepo)
	
	companyID := uuid.New()
	
	expectedReceivables := []*model.Receivable{
		{
			ID:            uuid.New(),
			CustomerID:    uuid.New(),
			InvoiceID:     uuid.New(),
			Amount:        5000.00,
			Currency:      "USD",
			DueDate:       time.Now().AddDate(0, 1, 0),
			DaysOverdue:   0,
			Status:        "outstanding",
			CollectionPriority: "medium",
		},
		{
			ID:            uuid.New(),
			CustomerID:    uuid.New(),
			InvoiceID:     uuid.New(),
			Amount:        2000.00,
			Currency:      "USD",
			DueDate:       time.Now().AddDate(0, -1, 0), // Overdue
			DaysOverdue:   30,
			Status:        "overdue",
			CollectionPriority: "high",
		},
	}
	
	filter := map[string]interface{}{
		"company_id": companyID,
		"status":     "outstanding",
	}
	
	mockFinanceRepo.On("GetReceivables", mock.Anything, filter).Return(expectedReceivables, nil)
	
	// Act
	result, err := financeService.GetReceivables(context.Background(), filter)
	
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, "outstanding", result[0].Status)
	assert.Equal(t, "overdue", result[1].Status)
	assert.Equal(t, "high", result[1].CollectionPriority)
	
	mockFinanceRepo.AssertExpectations(t)
}

func TestFinanceService_GetFinancialSummary_Success(t *testing.T) {
	// Arrange
	mockFinanceRepo := new(MockFinanceRepository)
	financeService := service.NewFinanceService(mockFinanceRepo)
	
	companyID := uuid.New()
	dateRange := "current_month"
	
	expectedSummary := &service.FinancialSummary{
		Revenue: service.RevenueSummary{
			TotalRevenue:    100000.00,
			InvoicedAmount:  90000.00,
			PaidAmount:      80000.00,
			OutstandingAmount: 10000.00,
		},
		Expenses: service.ExpenseSummary{
			TotalExpenses:   25000.00,
			PendingExpenses: 5000.00,
			ApprovedExpenses: 20000.00,
		},
		Receivables: service.ReceivableSummary{
			TotalReceivables: 15000.00,
			OverdueAmount:    3000.00,
			CurrentAmount:    12000.00,
		},
		CashFlow: service.CashFlowSummary{
			NetCashFlow:    55000.00,
			CashInflow:     80000.00,
			CashOutflow:    25000.00,
		},
	}
	
	mockFinanceRepo.On("GetFinancialSummary", mock.Anything, companyID, dateRange).Return(expectedSummary, nil)
	
	// Act
	result, err := financeService.GetFinancialSummary(context.Background(), companyID, dateRange)
	
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 100000.00, result.Revenue.TotalRevenue)
	assert.Equal(t, 80000.00, result.Revenue.PaidAmount)
	assert.Equal(t, 25000.00, result.Expenses.TotalExpenses)
	assert.Equal(t, 15000.00, result.Receivables.TotalReceivables)
	assert.Equal(t, 55000.00, result.CashFlow.NetCashFlow)
	
	mockFinanceRepo.AssertExpectations(t)
}

func TestFinanceService_CalculateInvoiceTax_Success(t *testing.T) {
	// Arrange
	mockFinanceRepo := new(MockFinanceRepository)
	financeService := service.NewFinanceService(mockFinanceRepo)
	
	// Test tax calculation
	amount := 10000.00
	taxRate := 0.05 // 5%
	
	// Act
	taxAmount, totalAmount := financeService.CalculateInvoiceTax(amount, taxRate)
	
	// Assert
	expectedTaxAmount := amount * taxRate
	expectedTotalAmount := amount + expectedTaxAmount
	
	assert.Equal(t, expectedTaxAmount, taxAmount)
	assert.Equal(t, expectedTotalAmount, totalAmount)
}

func TestFinanceService_CalculateAgingAnalysis_Success(t *testing.T) {
	// Arrange
	mockFinanceRepo := new(MockFinanceRepository)
	financeService := service.NewFinanceService(mockFinanceRepo)
	
	receivables := []*model.Receivable{
		{Amount: 1000.00, DaysOverdue: 0},    // Current
		{Amount: 2000.00, DaysOverdue: 15},   // 1-30 days
		{Amount: 1500.00, DaysOverdue: 45},   // 31-60 days
		{Amount: 800.00, DaysOverdue: 75},    // 61-90 days
		{Amount: 500.00, DaysOverdue: 120},   // Over 90 days
	}
	
	// Act
	analysis := financeService.CalculateAgingAnalysis(receivables)
	
	// Assert
	assert.Equal(t, 1000.00, analysis.Current)
	assert.Equal(t, 2000.00, analysis.Days1to30)
	assert.Equal(t, 1500.00, analysis.Days31to60)
	assert.Equal(t, 800.00, analysis.Days61to90)
	assert.Equal(t, 500.00, analysis.Over90Days)
	assert.Equal(t, 5800.00, analysis.Total)
}

func TestFinanceService_ValidatePaymentAmount_Success(t *testing.T) {
	// Arrange
	mockFinanceRepo := new(MockFinanceRepository)
	financeService := service.NewFinanceService(mockFinanceRepo)
	
	testCases := []struct {
		totalAmount    float64
		paidAmount     float64
		paymentAmount  float64
		shouldBeValid  bool
	}{
		{10000.00, 0.00, 5000.00, true},     // Valid partial payment
		{10000.00, 0.00, 10000.00, true},    // Valid full payment
		{10000.00, 5000.00, 5000.00, true},  // Valid remaining payment
		{10000.00, 8000.00, 3000.00, false}, // Overpayment
		{10000.00, 0.00, -1000.00, false},   // Negative amount
		{10000.00, 0.00, 0.00, false},       // Zero amount
	}
	
	for _, tc := range testCases {
		// Act
		isValid := financeService.ValidatePaymentAmount(tc.totalAmount, tc.paidAmount, tc.paymentAmount)
		
		// Assert
		assert.Equal(t, tc.shouldBeValid, isValid,
			"Payment validation failed for total: %.2f, paid: %.2f, payment: %.2f",
			tc.totalAmount, tc.paidAmount, tc.paymentAmount)
	}
}