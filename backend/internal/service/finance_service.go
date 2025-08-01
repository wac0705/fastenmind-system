package service

import (
	"fmt"
	"time"
	"github.com/fastenmind/fastener-api/internal/models"
	"github.com/fastenmind/fastener-api/internal/repository"
	"github.com/google/uuid"
)

type FinanceService interface {
	// Invoice operations
	CreateInvoice(invoice *models.Invoice) error
	UpdateInvoice(invoice *models.Invoice) error
	GetInvoice(id uuid.UUID) (*models.Invoice, error)
	GetInvoiceByNo(invoiceNo string) (*models.Invoice, error)
	ListInvoices(companyID uuid.UUID, params map[string]interface{}) ([]models.Invoice, int64, error)
	GenerateInvoiceFromOrder(orderID uuid.UUID, userID uuid.UUID) (*models.Invoice, error)
	
	// Payment operations
	ProcessPayment(payment *models.Payment) error
	GetPayment(id uuid.UUID) (*models.Payment, error)
	ListPayments(companyID uuid.UUID, params map[string]interface{}) ([]models.Payment, int64, error)
	GetPaymentsByInvoice(invoiceID uuid.UUID) ([]models.Payment, error)
	
	// Invoice items
	GetInvoiceItems(invoiceID uuid.UUID) ([]models.InvoiceItem, error)
	
	// Expense operations
	CreateExpense(expense *models.Expense) error
	UpdateExpense(expense *models.Expense) error
	ApproveExpense(id uuid.UUID, approverID uuid.UUID) error
	RejectExpense(id uuid.UUID, approverID uuid.UUID, reason string) error
	GetExpense(id uuid.UUID) (*models.Expense, error)
	ListExpenses(companyID uuid.UUID, params map[string]interface{}) ([]models.Expense, int64, error)
	
	// AR/AP operations
	GetAccountReceivables(companyID uuid.UUID, params map[string]interface{}) ([]models.AccountReceivable, error)
	GetAccountPayables(companyID uuid.UUID, params map[string]interface{}) ([]models.AccountPayable, error)
	GetARSummary(companyID uuid.UUID) (*ARAPSummary, error)
	GetAPSummary(companyID uuid.UUID) (*ARAPSummary, error)
	
	// Bank account operations
	CreateBankAccount(account *models.BankAccount) error
	UpdateBankAccount(account *models.BankAccount) error
	GetBankAccount(id uuid.UUID) (*models.BankAccount, error)
	ListBankAccounts(companyID uuid.UUID) ([]models.BankAccount, error)
	
	// Financial period operations
	CreateFinancialPeriod(period *models.FinancialPeriod) error
	CloseFinancialPeriod(id uuid.UUID, closerID uuid.UUID) error
	GetCurrentPeriod(companyID uuid.UUID) (*models.FinancialPeriod, error)
	ListFinancialPeriods(companyID uuid.UUID) ([]models.FinancialPeriod, error)
	
	// Reports
	GetFinancialDashboard(companyID uuid.UUID) (*FinancialDashboard, error)
	GetCashFlowReport(companyID uuid.UUID, startDate, endDate time.Time) (*CashFlowReport, error)
	GetAgingReport(companyID uuid.UUID, reportType string) (*AgingReport, error)
}

type financeService struct {
	financeRepo  repository.FinanceRepository
	orderRepo    repository.OrderRepository
	customerRepo repository.CustomerRepository
	supplierRepo repository.SupplierRepository
}

func NewFinanceService(
	financeRepo repository.FinanceRepository,
	orderRepo repository.OrderRepository,
	customerRepo repository.CustomerRepository,
	supplierRepo repository.SupplierRepository,
) FinanceService {
	return &financeService{
		financeRepo:  financeRepo,
		orderRepo:    orderRepo,
		customerRepo: customerRepo,
		supplierRepo: supplierRepo,
	}
}

// Invoice operations
func (s *financeService) CreateInvoice(invoice *models.Invoice) error {
	// Generate invoice number
	invoice.InvoiceNo = s.generateInvoiceNo(invoice.CompanyID, invoice.Type)
	
	// Calculate totals
	invoice.TotalAmount = invoice.SubTotal - invoice.DiscountAmount + invoice.TaxAmount
	invoice.BalanceAmount = invoice.TotalAmount - invoice.PaidAmount
	
	// Set initial status
	if invoice.Status == "" {
		invoice.Status = "draft"
	}
	
	if err := s.financeRepo.CreateInvoice(invoice); err != nil {
		return err
	}
	
	// Create AR/AP record if issued
	if invoice.Status == "issued" {
		if err := s.createARAPRecord(invoice); err != nil {
			return err
		}
	}
	
	return nil
}

func (s *financeService) UpdateInvoice(invoice *models.Invoice) error {
	// Recalculate totals
	invoice.TotalAmount = invoice.SubTotal - invoice.DiscountAmount + invoice.TaxAmount
	invoice.BalanceAmount = invoice.TotalAmount - invoice.PaidAmount
	
	// Update status based on payment
	if invoice.BalanceAmount <= 0 && invoice.Status != "cancelled" {
		invoice.Status = "paid"
	} else if invoice.PaidAmount > 0 && invoice.Status != "cancelled" {
		invoice.Status = "partial_paid"
	}
	
	return s.financeRepo.UpdateInvoice(invoice)
}

func (s *financeService) GetInvoice(id uuid.UUID) (*models.Invoice, error) {
	return s.financeRepo.GetInvoice(id)
}

func (s *financeService) GetInvoiceByNo(invoiceNo string) (*models.Invoice, error) {
	return s.financeRepo.GetInvoiceByNo(invoiceNo)
}

func (s *financeService) ListInvoices(companyID uuid.UUID, params map[string]interface{}) ([]models.Invoice, int64, error) {
	return s.financeRepo.ListInvoices(companyID, params)
}

func (s *financeService) GenerateInvoiceFromOrder(orderID uuid.UUID, userID uuid.UUID) (*models.Invoice, error) {
	// Get order details
	order, err := s.orderRepo.GetOrder(orderID)
	if err != nil {
		return nil, err
	}
	
	if order.Status != "ready_to_ship" && order.Status != "shipped" {
		return nil, fmt.Errorf("order must be ready to ship or shipped to generate invoice")
	}
	
	// Create invoice
	invoice := &models.Invoice{
		CompanyID:    order.CompanyID,
		Type:         "sales",
		Status:       "issued",
		OrderID:      &order.ID,
		CustomerID:   &order.CustomerID,
		IssueDate:    time.Now(),
		DueDate:      time.Now().AddDate(0, 0, 30), // Default 30 days payment terms
		SubTotal:     order.SubTotal,
		TaxRate:      order.TaxRate,
		TaxAmount:    order.TaxAmount,
		TotalAmount:  order.TotalAmount,
		PaidAmount:   0,
		BalanceAmount: order.TotalAmount,
		Currency:     order.Currency,
		ExchangeRate: order.ExchangeRate,
		PaymentTerms: order.PaymentTerms,
		CreatedBy:    userID,
	}
	
	if err := s.CreateInvoice(invoice); err != nil {
		return nil, err
	}
	
	// Create invoice items from order items
	orderItems, err := s.orderRepo.GetOrderItems(orderID)
	if err != nil {
		return nil, err
	}
	
	for _, orderItem := range orderItems {
		invoiceItem := &models.InvoiceItem{
			InvoiceID:   invoice.ID,
			Description: fmt.Sprintf("%s - %s", orderItem.ProductName, orderItem.Specification),
			Quantity:    orderItem.Quantity,
			Unit:        orderItem.Unit,
			UnitPrice:   orderItem.UnitPrice,
			TotalPrice:  orderItem.TotalPrice,
			OrderItemID: &orderItem.ID,
		}
		
		if err := s.financeRepo.CreateInvoiceItem(invoiceItem); err != nil {
			return nil, err
		}
	}
	
	// Update order status
	order.InvoiceID = &invoice.ID
	order.Status = "invoiced"
	if err := s.orderRepo.UpdateOrder(order); err != nil {
		return nil, err
	}
	
	return invoice, nil
}

// Payment operations
func (s *financeService) ProcessPayment(payment *models.Payment) error {
	// Generate payment number
	payment.PaymentNo = s.generatePaymentNo(payment.CompanyID)
	payment.Status = "pending"
	
	// Create payment and update invoice
	if err := s.financeRepo.CreatePayment(payment); err != nil {
		return err
	}
	
	// Mark payment as completed (in real system, this would be after bank confirmation)
	payment.Status = "completed"
	return s.financeRepo.UpdatePayment(payment)
}

func (s *financeService) GetPayment(id uuid.UUID) (*models.Payment, error) {
	return s.financeRepo.GetPayment(id)
}

func (s *financeService) ListPayments(companyID uuid.UUID, params map[string]interface{}) ([]models.Payment, int64, error) {
	return s.financeRepo.ListPayments(companyID, params)
}

func (s *financeService) GetPaymentsByInvoice(invoiceID uuid.UUID) ([]models.Payment, error) {
	return s.financeRepo.GetPaymentsByInvoice(invoiceID)
}

func (s *financeService) GetInvoiceItems(invoiceID uuid.UUID) ([]models.InvoiceItem, error) {
	return s.financeRepo.GetInvoiceItems(invoiceID)
}

// Expense operations
func (s *financeService) CreateExpense(expense *models.Expense) error {
	expense.ExpenseNo = s.generateExpenseNo(expense.CompanyID)
	expense.Status = "draft"
	expense.PaymentStatus = "unpaid"
	expense.SubmittedAt = time.Now()
	
	// Calculate total
	expense.TotalAmount = expense.Amount + expense.TaxAmount
	
	return s.financeRepo.CreateExpense(expense)
}

func (s *financeService) UpdateExpense(expense *models.Expense) error {
	// Recalculate total
	expense.TotalAmount = expense.Amount + expense.TaxAmount
	return s.financeRepo.UpdateExpense(expense)
}

func (s *financeService) ApproveExpense(id uuid.UUID, approverID uuid.UUID) error {
	expense, err := s.financeRepo.GetExpense(id)
	if err != nil {
		return err
	}
	
	if expense.Status != "submitted" {
		return fmt.Errorf("expense must be in submitted status to approve")
	}
	
	now := time.Now()
	expense.Status = "approved"
	expense.ApprovedBy = &approverID
	expense.ApprovedAt = &now
	
	return s.financeRepo.UpdateExpense(expense)
}

func (s *financeService) RejectExpense(id uuid.UUID, approverID uuid.UUID, reason string) error {
	expense, err := s.financeRepo.GetExpense(id)
	if err != nil {
		return err
	}
	
	if expense.Status != "submitted" {
		return fmt.Errorf("expense must be in submitted status to reject")
	}
	
	now := time.Now()
	expense.Status = "rejected"
	expense.ApprovedBy = &approverID
	expense.ApprovedAt = &now
	expense.Notes = fmt.Sprintf("Rejected: %s", reason)
	
	return s.financeRepo.UpdateExpense(expense)
}

func (s *financeService) GetExpense(id uuid.UUID) (*models.Expense, error) {
	return s.financeRepo.GetExpense(id)
}

func (s *financeService) ListExpenses(companyID uuid.UUID, params map[string]interface{}) ([]models.Expense, int64, error) {
	return s.financeRepo.ListExpenses(companyID, params)
}

// AR/AP operations
func (s *financeService) GetAccountReceivables(companyID uuid.UUID, params map[string]interface{}) ([]models.AccountReceivable, error) {
	ars, err := s.financeRepo.ListAccountReceivables(companyID, params)
	if err != nil {
		return nil, err
	}
	
	// Update aging info
	for i := range ars {
		s.financeRepo.UpdateAccountReceivable(&ars[i])
	}
	
	return ars, nil
}

func (s *financeService) GetAccountPayables(companyID uuid.UUID, params map[string]interface{}) ([]models.AccountPayable, error) {
	return s.financeRepo.ListAccountPayables(companyID, params)
}

func (s *financeService) GetARSummary(companyID uuid.UUID) (*ARAPSummary, error) {
	ars, err := s.financeRepo.ListAccountReceivables(companyID, map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	
	summary := &ARAPSummary{
		Currency: "USD",
	}
	
	for _, ar := range ars {
		summary.TotalAmount += ar.InvoiceAmount
		summary.PaidAmount += ar.PaidAmount
		summary.BalanceAmount += ar.BalanceAmount
		
		switch ar.AgingCategory {
		case "current":
			summary.Current += ar.BalanceAmount
		case "30days":
			summary.Days30 += ar.BalanceAmount
		case "60days":
			summary.Days60 += ar.BalanceAmount
		case "90days":
			summary.Days90 += ar.BalanceAmount
		case "over90days":
			summary.Over90 += ar.BalanceAmount
		}
		
		if ar.BalanceAmount > 0 {
			summary.OpenItems++
		}
	}
	
	return summary, nil
}

func (s *financeService) GetAPSummary(companyID uuid.UUID) (*ARAPSummary, error) {
	aps, err := s.financeRepo.ListAccountPayables(companyID, map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	
	summary := &ARAPSummary{
		Currency: "USD",
	}
	
	for _, ap := range aps {
		summary.TotalAmount += ap.InvoiceAmount
		summary.PaidAmount += ap.PaidAmount
		summary.BalanceAmount += ap.BalanceAmount
		
		if ap.BalanceAmount > 0 {
			summary.OpenItems++
		}
	}
	
	return summary, nil
}

// Bank account operations
func (s *financeService) CreateBankAccount(account *models.BankAccount) error {
	account.Status = "active"
	return s.financeRepo.CreateBankAccount(account)
}

func (s *financeService) UpdateBankAccount(account *models.BankAccount) error {
	return s.financeRepo.UpdateBankAccount(account)
}

func (s *financeService) GetBankAccount(id uuid.UUID) (*models.BankAccount, error) {
	return s.financeRepo.GetBankAccount(id)
}

func (s *financeService) ListBankAccounts(companyID uuid.UUID) ([]models.BankAccount, error) {
	return s.financeRepo.ListBankAccounts(companyID)
}

// Financial period operations
func (s *financeService) CreateFinancialPeriod(period *models.FinancialPeriod) error {
	period.Status = "open"
	
	// Set as current if no other current period exists
	current, err := s.financeRepo.GetCurrentPeriod(period.CompanyID)
	if err != nil || current == nil {
		period.IsCurrent = true
	}
	
	return s.financeRepo.CreateFinancialPeriod(period)
}

func (s *financeService) CloseFinancialPeriod(id uuid.UUID, closerID uuid.UUID) error {
	period, err := s.financeRepo.GetFinancialPeriod(id)
	if err != nil {
		return err
	}
	
	if period.Status == "closed" {
		return fmt.Errorf("period is already closed")
	}
	
	now := time.Now()
	period.Status = "closed"
	period.ClosedAt = &now
	period.ClosedBy = &closerID
	
	return s.financeRepo.UpdateFinancialPeriod(period)
}

func (s *financeService) GetCurrentPeriod(companyID uuid.UUID) (*models.FinancialPeriod, error) {
	return s.financeRepo.GetCurrentPeriod(companyID)
}

func (s *financeService) ListFinancialPeriods(companyID uuid.UUID) ([]models.FinancialPeriod, error) {
	return s.financeRepo.ListFinancialPeriods(companyID)
}

// Reports
func (s *financeService) GetFinancialDashboard(companyID uuid.UUID) (*FinancialDashboard, error) {
	dashboard := &FinancialDashboard{
		Currency: "USD",
		Date:     time.Now(),
	}
	
	// Get current month date range
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, -1)
	
	// Get revenue (sales invoices)
	revenueParams := map[string]interface{}{
		"type":       "sales",
		"start_date": startOfMonth.Format("2006-01-02"),
		"end_date":   endOfMonth.Format("2006-01-02"),
	}
	invoices, _, _ := s.financeRepo.ListInvoices(companyID, revenueParams)
	
	for _, inv := range invoices {
		if inv.Status != "cancelled" {
			dashboard.Revenue += inv.TotalAmount
		}
	}
	
	// Get expenses
	expenseParams := map[string]interface{}{
		"start_date": startOfMonth.Format("2006-01-02"),
		"end_date":   endOfMonth.Format("2006-01-02"),
		"status":     "approved",
	}
	expenses, _, _ := s.financeRepo.ListExpenses(companyID, expenseParams)
	
	for _, exp := range expenses {
		dashboard.Expenses += exp.TotalAmount
	}
	
	// Calculate profit
	dashboard.Profit = dashboard.Revenue - dashboard.Expenses
	
	// Get AR/AP summaries
	dashboard.ARSummary, _ = s.GetARSummary(companyID)
	dashboard.APSummary, _ = s.GetAPSummary(companyID)
	
	// Get cash balance from bank accounts
	accounts, _ := s.financeRepo.ListBankAccounts(companyID)
	for _, acc := range accounts {
		dashboard.CashBalance += acc.CurrentBalance
	}
	
	return dashboard, nil
}

func (s *financeService) GetCashFlowReport(companyID uuid.UUID, startDate, endDate time.Time) (*CashFlowReport, error) {
	report := &CashFlowReport{
		StartDate: startDate,
		EndDate:   endDate,
		Currency:  "USD",
	}
	
	// Get all payments in date range
	paymentParams := map[string]interface{}{
		"start_date": startDate.Format("2006-01-02"),
		"end_date":   endDate.Format("2006-01-02"),
		"status":     "completed",
	}
	payments, _, _ := s.financeRepo.ListPayments(companyID, paymentParams)
	
	for _, payment := range payments {
		if payment.Type == "incoming" {
			report.CashInflows += payment.Amount
			
			if payment.InvoiceID != nil {
				report.OperatingInflows += payment.Amount
			} else {
				report.OtherInflows += payment.Amount
			}
		} else {
			report.CashOutflows += payment.Amount
			
			if payment.InvoiceID != nil {
				report.OperatingOutflows += payment.Amount
			} else {
				report.OtherOutflows += payment.Amount
			}
		}
	}
	
	report.NetCashFlow = report.CashInflows - report.CashOutflows
	
	return report, nil
}

func (s *financeService) GetAgingReport(companyID uuid.UUID, reportType string) (*AgingReport, error) {
	report := &AgingReport{
		ReportType: reportType,
		Currency:   "USD",
		Date:       time.Now(),
	}
	
	if reportType == "receivable" {
		ars, _ := s.financeRepo.ListAccountReceivables(companyID, map[string]interface{}{})
		for _, ar := range ars {
			if ar.BalanceAmount > 0 {
				item := AgingItem{
					CustomerID:   &ar.CustomerID,
					InvoiceNo:    ar.Invoice.InvoiceNo,
					InvoiceDate:  ar.InvoiceDate,
					DueDate:      ar.DueDate,
					Amount:       ar.InvoiceAmount,
					PaidAmount:   ar.PaidAmount,
					Balance:      ar.BalanceAmount,
					DaysOverdue:  ar.DaysOverdue,
					AgingBucket:  ar.AgingCategory,
				}
				
				if ar.Customer != nil {
					item.CustomerName = ar.Customer.Name
				}
				
				report.Items = append(report.Items, item)
				
				// Update totals
				switch ar.AgingCategory {
				case "current":
					report.Current += ar.BalanceAmount
				case "30days":
					report.Days30 += ar.BalanceAmount
				case "60days":
					report.Days60 += ar.BalanceAmount
				case "90days":
					report.Days90 += ar.BalanceAmount
				case "over90days":
					report.Over90 += ar.BalanceAmount
				}
			}
		}
	} else {
		aps, _ := s.financeRepo.ListAccountPayables(companyID, map[string]interface{}{})
		for _, ap := range aps {
			if ap.BalanceAmount > 0 {
				item := AgingItem{
					SupplierID:   &ap.SupplierID,
					InvoiceNo:    ap.Invoice.InvoiceNo,
					InvoiceDate:  ap.InvoiceDate,
					DueDate:      ap.DueDate,
					Amount:       ap.InvoiceAmount,
					PaidAmount:   ap.PaidAmount,
					Balance:      ap.BalanceAmount,
				}
				
				if ap.Supplier != nil {
					item.SupplierName = ap.Supplier.Name
				}
				
				report.Items = append(report.Items, item)
			}
		}
	}
	
	report.Total = report.Current + report.Days30 + report.Days60 + report.Days90 + report.Over90
	
	return report, nil
}

// Helper methods
func (s *financeService) generateInvoiceNo(companyID uuid.UUID, invoiceType string) string {
	prefix := "INV"
	if invoiceType == "purchase" {
		prefix = "PINV"
	} else if invoiceType == "credit_note" {
		prefix = "CN"
	} else if invoiceType == "debit_note" {
		prefix = "DN"
	}
	
	return fmt.Sprintf("%s-%s-%06d", prefix, time.Now().Format("200601"), time.Now().Unix()%1000000)
}

func (s *financeService) generatePaymentNo(companyID uuid.UUID) string {
	return fmt.Sprintf("PAY-%s-%06d", time.Now().Format("200601"), time.Now().Unix()%1000000)
}

func (s *financeService) generateExpenseNo(companyID uuid.UUID) string {
	return fmt.Sprintf("EXP-%s-%06d", time.Now().Format("200601"), time.Now().Unix()%1000000)
}

func (s *financeService) createARAPRecord(invoice *models.Invoice) error {
	if invoice.Type == "sales" && invoice.CustomerID != nil {
		ar := &models.AccountReceivable{
			CompanyID:     invoice.CompanyID,
			CustomerID:    *invoice.CustomerID,
			InvoiceID:     invoice.ID,
			InvoiceAmount: invoice.TotalAmount,
			PaidAmount:    invoice.PaidAmount,
			BalanceAmount: invoice.BalanceAmount,
			Currency:      invoice.Currency,
			InvoiceDate:   invoice.IssueDate,
			DueDate:       invoice.DueDate,
			Status:        "open",
		}
		return s.financeRepo.CreateAccountReceivable(ar)
	} else if invoice.Type == "purchase" && invoice.SupplierID != nil {
		ap := &models.AccountPayable{
			CompanyID:      invoice.CompanyID,
			SupplierID:     *invoice.SupplierID,
			InvoiceID:      invoice.ID,
			InvoiceAmount:  invoice.TotalAmount,
			PaidAmount:     invoice.PaidAmount,
			BalanceAmount:  invoice.BalanceAmount,
			Currency:       invoice.Currency,
			InvoiceDate:    invoice.IssueDate,
			DueDate:        invoice.DueDate,
			Status:         "open",
			PaymentPriority: "medium",
		}
		return s.financeRepo.CreateAccountPayable(ap)
	}
	
	return nil
}

// Report structs
type ARAPSummary struct {
	TotalAmount   float64 `json:"total_amount"`
	PaidAmount    float64 `json:"paid_amount"`
	BalanceAmount float64 `json:"balance_amount"`
	Current       float64 `json:"current"`
	Days30        float64 `json:"days_30"`
	Days60        float64 `json:"days_60"`
	Days90        float64 `json:"days_90"`
	Over90        float64 `json:"over_90"`
	OpenItems     int     `json:"open_items"`
	Currency      string  `json:"currency"`
}

type FinancialDashboard struct {
	Revenue     float64      `json:"revenue"`
	Expenses    float64      `json:"expenses"`
	Profit      float64      `json:"profit"`
	CashBalance float64      `json:"cash_balance"`
	ARSummary   *ARAPSummary `json:"ar_summary"`
	APSummary   *ARAPSummary `json:"ap_summary"`
	Currency    string       `json:"currency"`
	Date        time.Time    `json:"date"`
}

type CashFlowReport struct {
	StartDate        time.Time `json:"start_date"`
	EndDate          time.Time `json:"end_date"`
	CashInflows      float64   `json:"cash_inflows"`
	CashOutflows     float64   `json:"cash_outflows"`
	NetCashFlow      float64   `json:"net_cash_flow"`
	OperatingInflows float64   `json:"operating_inflows"`
	OperatingOutflows float64  `json:"operating_outflows"`
	OtherInflows     float64   `json:"other_inflows"`
	OtherOutflows    float64   `json:"other_outflows"`
	Currency         string    `json:"currency"`
}

type AgingReport struct {
	ReportType string      `json:"report_type"`
	Items      []AgingItem `json:"items"`
	Current    float64     `json:"current"`
	Days30     float64     `json:"days_30"`
	Days60     float64     `json:"days_60"`
	Days90     float64     `json:"days_90"`
	Over90     float64     `json:"over_90"`
	Total      float64     `json:"total"`
	Currency   string      `json:"currency"`
	Date       time.Time   `json:"date"`
}

type AgingItem struct {
	CustomerID   *uuid.UUID `json:"customer_id,omitempty"`
	CustomerName string     `json:"customer_name,omitempty"`
	SupplierID   *uuid.UUID `json:"supplier_id,omitempty"`
	SupplierName string     `json:"supplier_name,omitempty"`
	InvoiceNo    string     `json:"invoice_no"`
	InvoiceDate  time.Time  `json:"invoice_date"`
	DueDate      time.Time  `json:"due_date"`
	Amount       float64    `json:"amount"`
	PaidAmount   float64    `json:"paid_amount"`
	Balance      float64    `json:"balance"`
	DaysOverdue  int        `json:"days_overdue"`
	AgingBucket  string     `json:"aging_bucket"`
}