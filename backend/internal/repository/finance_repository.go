package repository

import (
	"time"
	"github.com/fastenmind/fastener-api/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FinanceRepository interface {
	// Invoice operations
	CreateInvoice(invoice *models.Invoice) error
	UpdateInvoice(invoice *models.Invoice) error
	GetInvoice(id uuid.UUID) (*models.Invoice, error)
	GetInvoiceByNo(invoiceNo string) (*models.Invoice, error)
	ListInvoices(companyID uuid.UUID, params map[string]interface{}) ([]models.Invoice, int64, error)
	CreateInvoiceItem(item *models.InvoiceItem) error
	GetInvoiceItems(invoiceID uuid.UUID) ([]models.InvoiceItem, error)
	
	// Payment operations
	CreatePayment(payment *models.Payment) error
	UpdatePayment(payment *models.Payment) error
	GetPayment(id uuid.UUID) (*models.Payment, error)
	ListPayments(companyID uuid.UUID, params map[string]interface{}) ([]models.Payment, int64, error)
	GetPaymentsByInvoice(invoiceID uuid.UUID) ([]models.Payment, error)
	
	// Expense operations
	CreateExpense(expense *models.Expense) error
	UpdateExpense(expense *models.Expense) error
	GetExpense(id uuid.UUID) (*models.Expense, error)
	ListExpenses(companyID uuid.UUID, params map[string]interface{}) ([]models.Expense, int64, error)
	
	// AR/AP operations
	CreateAccountReceivable(ar *models.AccountReceivable) error
	UpdateAccountReceivable(ar *models.AccountReceivable) error
	GetAccountReceivable(id uuid.UUID) (*models.AccountReceivable, error)
	ListAccountReceivables(companyID uuid.UUID, params map[string]interface{}) ([]models.AccountReceivable, error)
	GetARByCustomer(customerID uuid.UUID) ([]models.AccountReceivable, error)
	
	CreateAccountPayable(ap *models.AccountPayable) error
	UpdateAccountPayable(ap *models.AccountPayable) error
	GetAccountPayable(id uuid.UUID) (*models.AccountPayable, error)
	ListAccountPayables(companyID uuid.UUID, params map[string]interface{}) ([]models.AccountPayable, error)
	GetAPBySupplier(supplierID uuid.UUID) ([]models.AccountPayable, error)
	
	// Bank account operations
	CreateBankAccount(account *models.BankAccount) error
	UpdateBankAccount(account *models.BankAccount) error
	GetBankAccount(id uuid.UUID) (*models.BankAccount, error)
	ListBankAccounts(companyID uuid.UUID) ([]models.BankAccount, error)
	
	// Financial period operations
	CreateFinancialPeriod(period *models.FinancialPeriod) error
	UpdateFinancialPeriod(period *models.FinancialPeriod) error
	GetFinancialPeriod(id uuid.UUID) (*models.FinancialPeriod, error)
	GetCurrentPeriod(companyID uuid.UUID) (*models.FinancialPeriod, error)
	ListFinancialPeriods(companyID uuid.UUID) ([]models.FinancialPeriod, error)
}

type financeRepository struct {
	db *gorm.DB
}

func NewFinanceRepository(db interface{}) FinanceRepository {
	gormDB, ok := db.(*gorm.DB)
	if !ok {
		panic("invalid database type, expected *gorm.DB")
	}
	return &financeRepository{db: gormDB}
}

// Invoice operations
func (r *financeRepository) CreateInvoice(invoice *models.Invoice) error {
	return r.db.Create(invoice).Error
}

func (r *financeRepository) UpdateInvoice(invoice *models.Invoice) error {
	return r.db.Save(invoice).Error
}

func (r *financeRepository) GetInvoice(id uuid.UUID) (*models.Invoice, error) {
	var invoice models.Invoice
	err := r.db.Preload("Order").
		Preload("Customer").
		Preload("Supplier").
		Preload("Creator").
		First(&invoice, id).Error
	if err != nil {
		return nil, err
	}
	return &invoice, nil
}

func (r *financeRepository) GetInvoiceByNo(invoiceNo string) (*models.Invoice, error) {
	var invoice models.Invoice
	err := r.db.Where("invoice_no = ?", invoiceNo).First(&invoice).Error
	if err != nil {
		return nil, err
	}
	return &invoice, nil
}

func (r *financeRepository) ListInvoices(companyID uuid.UUID, params map[string]interface{}) ([]models.Invoice, int64, error) {
	var invoices []models.Invoice
	var total int64
	
	query := r.db.Model(&models.Invoice{}).Where("company_id = ?", companyID)
	
	// Apply filters
	if invoiceType, ok := params["type"].(string); ok && invoiceType != "" {
		query = query.Where("type = ?", invoiceType)
	}
	
	if status, ok := params["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	
	if customerID, ok := params["customer_id"].(string); ok && customerID != "" {
		query = query.Where("customer_id = ?", customerID)
	}
	
	if supplierID, ok := params["supplier_id"].(string); ok && supplierID != "" {
		query = query.Where("supplier_id = ?", supplierID)
	}
	
	// Date range filters
	if startDate, ok := params["start_date"].(string); ok && startDate != "" {
		query = query.Where("issue_date >= ?", startDate)
	}
	
	if endDate, ok := params["end_date"].(string); ok && endDate != "" {
		query = query.Where("issue_date <= ?", endDate)
	}
	
	// Overdue filter
	if overdue, ok := params["overdue"].(bool); ok && overdue {
		query = query.Where("due_date < ? AND status NOT IN ?", time.Now(), []string{"paid", "cancelled"})
	}
	
	if search, ok := params["search"].(string); ok && search != "" {
		query = query.Where("invoice_no LIKE ?", "%"+search+"%")
	}
	
	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Apply pagination
	if page, ok := params["page"].(int); ok && page > 0 {
		pageSize := 20
		if ps, ok := params["page_size"].(int); ok && ps > 0 {
			pageSize = ps
		}
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}
	
	// Apply sorting
	sortBy := "created_at"
	if sb, ok := params["sort_by"].(string); ok && sb != "" {
		sortBy = sb
	}
	sortOrder := "DESC"
	if so, ok := params["sort_order"].(string); ok && so != "" {
		sortOrder = so
	}
	query = query.Order(sortBy + " " + sortOrder)
	
	// Load with relations
	if err := query.
		Preload("Customer").
		Preload("Supplier").
		Find(&invoices).Error; err != nil {
		return nil, 0, err
	}
	
	return invoices, total, nil
}

func (r *financeRepository) CreateInvoiceItem(item *models.InvoiceItem) error {
	return r.db.Create(item).Error
}

func (r *financeRepository) GetInvoiceItems(invoiceID uuid.UUID) ([]models.InvoiceItem, error) {
	var items []models.InvoiceItem
	err := r.db.Where("invoice_id = ?", invoiceID).
		Preload("Inventory").
		Find(&items).Error
	return items, err
}

// Payment operations
func (r *financeRepository) CreatePayment(payment *models.Payment) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Create payment
		if err := tx.Create(payment).Error; err != nil {
			return err
		}
		
		// Update invoice if linked
		if payment.InvoiceID != nil {
			var invoice models.Invoice
			if err := tx.First(&invoice, payment.InvoiceID).Error; err != nil {
				return err
			}
			
			invoice.PaidAmount += payment.Amount
			invoice.BalanceAmount = invoice.TotalAmount - invoice.PaidAmount
			
			if invoice.BalanceAmount <= 0 {
				invoice.Status = "paid"
				now := time.Now()
				invoice.PaymentDate = &now
			} else if invoice.PaidAmount > 0 {
				invoice.Status = "partial_paid"
			}
			
			if err := tx.Save(&invoice).Error; err != nil {
				return err
			}
			
			// Update AR/AP
			if invoice.Type == "sales" && invoice.CustomerID != nil {
				var ar models.AccountReceivable
				if err := tx.Where("invoice_id = ?", invoice.ID).First(&ar).Error; err == nil {
					ar.PaidAmount = invoice.PaidAmount
					ar.BalanceAmount = invoice.BalanceAmount
					ar.LastPaymentDate = &payment.PaymentDate
					if ar.BalanceAmount <= 0 {
						ar.Status = "paid"
					} else {
						ar.Status = "partial"
					}
					tx.Save(&ar)
				}
			} else if invoice.Type == "purchase" && invoice.SupplierID != nil {
				var ap models.AccountPayable
				if err := tx.Where("invoice_id = ?", invoice.ID).First(&ap).Error; err == nil {
					ap.PaidAmount = invoice.PaidAmount
					ap.BalanceAmount = invoice.BalanceAmount
					ap.LastPaymentDate = &payment.PaymentDate
					if ap.BalanceAmount <= 0 {
						ap.Status = "paid"
					} else {
						ap.Status = "partial"
					}
					tx.Save(&ap)
				}
			}
		}
		
		return nil
	})
}

func (r *financeRepository) UpdatePayment(payment *models.Payment) error {
	return r.db.Save(payment).Error
}

func (r *financeRepository) GetPayment(id uuid.UUID) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Preload("Invoice").
		Preload("Customer").
		Preload("Supplier").
		Preload("Creator").
		First(&payment, id).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *financeRepository) ListPayments(companyID uuid.UUID, params map[string]interface{}) ([]models.Payment, int64, error) {
	var payments []models.Payment
	var total int64
	
	query := r.db.Model(&models.Payment{}).Where("company_id = ?", companyID)
	
	// Apply filters
	if paymentType, ok := params["type"].(string); ok && paymentType != "" {
		query = query.Where("type = ?", paymentType)
	}
	
	if status, ok := params["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	
	if paymentMethod, ok := params["payment_method"].(string); ok && paymentMethod != "" {
		query = query.Where("payment_method = ?", paymentMethod)
	}
	
	// Date range filters
	if startDate, ok := params["start_date"].(string); ok && startDate != "" {
		query = query.Where("payment_date >= ?", startDate)
	}
	
	if endDate, ok := params["end_date"].(string); ok && endDate != "" {
		query = query.Where("payment_date <= ?", endDate)
	}
	
	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Apply pagination
	if page, ok := params["page"].(int); ok && page > 0 {
		pageSize := 20
		if ps, ok := params["page_size"].(int); ok && ps > 0 {
			pageSize = ps
		}
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}
	
	// Apply sorting
	query = query.Order("payment_date DESC")
	
	// Load with relations
	if err := query.
		Preload("Invoice").
		Preload("Customer").
		Preload("Supplier").
		Find(&payments).Error; err != nil {
		return nil, 0, err
	}
	
	return payments, total, nil
}

func (r *financeRepository) GetPaymentsByInvoice(invoiceID uuid.UUID) ([]models.Payment, error) {
	var payments []models.Payment
	err := r.db.Where("invoice_id = ?", invoiceID).
		Order("payment_date DESC").
		Find(&payments).Error
	return payments, err
}

// Expense operations
func (r *financeRepository) CreateExpense(expense *models.Expense) error {
	return r.db.Create(expense).Error
}

func (r *financeRepository) UpdateExpense(expense *models.Expense) error {
	return r.db.Save(expense).Error
}

func (r *financeRepository) GetExpense(id uuid.UUID) (*models.Expense, error) {
	var expense models.Expense
	err := r.db.Preload("Supplier").
		Preload("Submitter").
		Preload("Approver").
		Preload("Payer").
		First(&expense, id).Error
	if err != nil {
		return nil, err
	}
	return &expense, nil
}

func (r *financeRepository) ListExpenses(companyID uuid.UUID, params map[string]interface{}) ([]models.Expense, int64, error) {
	var expenses []models.Expense
	var total int64
	
	query := r.db.Model(&models.Expense{}).Where("company_id = ?", companyID)
	
	// Apply filters
	if category, ok := params["category"].(string); ok && category != "" {
		query = query.Where("category = ?", category)
	}
	
	if status, ok := params["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	
	if paymentStatus, ok := params["payment_status"].(string); ok && paymentStatus != "" {
		query = query.Where("payment_status = ?", paymentStatus)
	}
	
	// Date range filters
	if startDate, ok := params["start_date"].(string); ok && startDate != "" {
		query = query.Where("expense_date >= ?", startDate)
	}
	
	if endDate, ok := params["end_date"].(string); ok && endDate != "" {
		query = query.Where("expense_date <= ?", endDate)
	}
	
	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Apply pagination
	if page, ok := params["page"].(int); ok && page > 0 {
		pageSize := 20
		if ps, ok := params["page_size"].(int); ok && ps > 0 {
			pageSize = ps
		}
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}
	
	// Apply sorting
	query = query.Order("expense_date DESC")
	
	// Load with relations
	if err := query.
		Preload("Supplier").
		Preload("Submitter").
		Find(&expenses).Error; err != nil {
		return nil, 0, err
	}
	
	return expenses, total, nil
}

// AR operations
func (r *financeRepository) CreateAccountReceivable(ar *models.AccountReceivable) error {
	return r.db.Create(ar).Error
}

func (r *financeRepository) UpdateAccountReceivable(ar *models.AccountReceivable) error {
	// Update aging info
	if ar.Status != "paid" && ar.Status != "written_off" {
		daysOverdue := int(time.Since(ar.DueDate).Hours() / 24)
		if daysOverdue > 0 {
			ar.DaysOverdue = daysOverdue
			if daysOverdue <= 30 {
				ar.AgingCategory = "30days"
			} else if daysOverdue <= 60 {
				ar.AgingCategory = "60days"
			} else if daysOverdue <= 90 {
				ar.AgingCategory = "90days"
			} else {
				ar.AgingCategory = "over90days"
			}
		} else {
			ar.DaysOverdue = 0
			ar.AgingCategory = "current"
		}
	}
	
	return r.db.Save(ar).Error
}

func (r *financeRepository) GetAccountReceivable(id uuid.UUID) (*models.AccountReceivable, error) {
	var ar models.AccountReceivable
	err := r.db.Preload("Customer").
		Preload("Invoice").
		First(&ar, id).Error
	if err != nil {
		return nil, err
	}
	return &ar, nil
}

func (r *financeRepository) ListAccountReceivables(companyID uuid.UUID, params map[string]interface{}) ([]models.AccountReceivable, error) {
	var ars []models.AccountReceivable
	query := r.db.Where("company_id = ?", companyID)
	
	// Apply filters
	if status, ok := params["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	
	if agingCategory, ok := params["aging_category"].(string); ok && agingCategory != "" {
		query = query.Where("aging_category = ?", agingCategory)
	}
	
	if customerID, ok := params["customer_id"].(string); ok && customerID != "" {
		query = query.Where("customer_id = ?", customerID)
	}
	
	err := query.
		Preload("Customer").
		Preload("Invoice").
		Order("due_date ASC").
		Find(&ars).Error
		
	return ars, err
}

func (r *financeRepository) GetARByCustomer(customerID uuid.UUID) ([]models.AccountReceivable, error) {
	var ars []models.AccountReceivable
	err := r.db.Where("customer_id = ? AND status NOT IN ?", customerID, []string{"paid", "written_off"}).
		Preload("Invoice").
		Order("due_date ASC").
		Find(&ars).Error
	return ars, err
}

// AP operations
func (r *financeRepository) CreateAccountPayable(ap *models.AccountPayable) error {
	return r.db.Create(ap).Error
}

func (r *financeRepository) UpdateAccountPayable(ap *models.AccountPayable) error {
	return r.db.Save(ap).Error
}

func (r *financeRepository) GetAccountPayable(id uuid.UUID) (*models.AccountPayable, error) {
	var ap models.AccountPayable
	err := r.db.Preload("Supplier").
		Preload("Invoice").
		First(&ap, id).Error
	if err != nil {
		return nil, err
	}
	return &ap, nil
}

func (r *financeRepository) ListAccountPayables(companyID uuid.UUID, params map[string]interface{}) ([]models.AccountPayable, error) {
	var aps []models.AccountPayable
	query := r.db.Where("company_id = ?", companyID)
	
	// Apply filters
	if status, ok := params["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	
	if priority, ok := params["payment_priority"].(string); ok && priority != "" {
		query = query.Where("payment_priority = ?", priority)
	}
	
	if supplierID, ok := params["supplier_id"].(string); ok && supplierID != "" {
		query = query.Where("supplier_id = ?", supplierID)
	}
	
	err := query.
		Preload("Supplier").
		Preload("Invoice").
		Order("due_date ASC").
		Find(&aps).Error
		
	return aps, err
}

func (r *financeRepository) GetAPBySupplier(supplierID uuid.UUID) ([]models.AccountPayable, error) {
	var aps []models.AccountPayable
	err := r.db.Where("supplier_id = ? AND status NOT IN ?", supplierID, []string{"paid"}).
		Preload("Invoice").
		Order("due_date ASC").
		Find(&aps).Error
	return aps, err
}

// Bank account operations
func (r *financeRepository) CreateBankAccount(account *models.BankAccount) error {
	return r.db.Create(account).Error
}

func (r *financeRepository) UpdateBankAccount(account *models.BankAccount) error {
	return r.db.Save(account).Error
}

func (r *financeRepository) GetBankAccount(id uuid.UUID) (*models.BankAccount, error) {
	var account models.BankAccount
	err := r.db.First(&account, id).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *financeRepository) ListBankAccounts(companyID uuid.UUID) ([]models.BankAccount, error) {
	var accounts []models.BankAccount
	err := r.db.Where("company_id = ? AND status = ?", companyID, "active").
		Order("is_default DESC, account_name ASC").
		Find(&accounts).Error
	return accounts, err
}

// Financial period operations
func (r *financeRepository) CreateFinancialPeriod(period *models.FinancialPeriod) error {
	return r.db.Create(period).Error
}

func (r *financeRepository) UpdateFinancialPeriod(period *models.FinancialPeriod) error {
	return r.db.Save(period).Error
}

func (r *financeRepository) GetFinancialPeriod(id uuid.UUID) (*models.FinancialPeriod, error) {
	var period models.FinancialPeriod
	err := r.db.First(&period, id).Error
	if err != nil {
		return nil, err
	}
	return &period, nil
}

func (r *financeRepository) GetCurrentPeriod(companyID uuid.UUID) (*models.FinancialPeriod, error) {
	var period models.FinancialPeriod
	err := r.db.Where("company_id = ? AND is_current = ?", companyID, true).
		First(&period).Error
	if err != nil {
		return nil, err
	}
	return &period, nil
}

func (r *financeRepository) ListFinancialPeriods(companyID uuid.UUID) ([]models.FinancialPeriod, error) {
	var periods []models.FinancialPeriod
	err := r.db.Where("company_id = ?", companyID).
		Order("start_date DESC").
		Find(&periods).Error
	return periods, err
}