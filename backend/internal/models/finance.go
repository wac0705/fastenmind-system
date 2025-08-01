package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Invoice represents a sales invoice
type Invoice struct {
	ID                uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID         uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	InvoiceNo         string     `gorm:"not null;unique" json:"invoice_no"`
	Type              string     `gorm:"not null" json:"type"`               // sales, purchase, credit_note, debit_note
	Status            string     `gorm:"not null" json:"status"`             // draft, issued, sent, paid, partial_paid, overdue, cancelled
	
	// Reference
	OrderID           *uuid.UUID `gorm:"type:uuid" json:"order_id"`
	CustomerID        *uuid.UUID `gorm:"type:uuid" json:"customer_id"`
	SupplierID        *uuid.UUID `gorm:"type:uuid" json:"supplier_id"`
	
	// Dates
	IssueDate         time.Time  `json:"issue_date"`
	DueDate           time.Time  `json:"due_date"`
	PaymentDate       *time.Time `json:"payment_date"`
	
	// Amounts
	SubTotal          float64    `json:"sub_total"`
	TaxRate           float64    `json:"tax_rate"`
	TaxAmount         float64    `json:"tax_amount"`
	DiscountRate      float64    `json:"discount_rate"`
	DiscountAmount    float64    `json:"discount_amount"`
	TotalAmount       float64    `json:"total_amount"`
	PaidAmount        float64    `json:"paid_amount"`
	BalanceAmount     float64    `json:"balance_amount"`
	Currency          string     `gorm:"default:'USD'" json:"currency"`
	ExchangeRate      float64    `gorm:"default:1" json:"exchange_rate"`
	
	// Payment Info
	PaymentTerms      string     `json:"payment_terms"`
	PaymentMethod     string     `json:"payment_method"`
	BankAccount       string     `json:"bank_account"`
	
	// Additional Info
	Notes             string     `json:"notes"`
	InternalNotes     string     `json:"internal_notes"`
	
	// Timestamps
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	CreatedBy         uuid.UUID  `gorm:"type:uuid;not null" json:"created_by"`
	
	// Relations
	Company           *Company   `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Order             *Order     `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Customer          *Customer  `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Supplier          *Supplier  `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	Creator           *User      `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

func (i *Invoice) BeforeCreate(tx *gorm.DB) error {
	i.ID = uuid.New()
	return nil
}

// InvoiceItem represents items in an invoice
type InvoiceItem struct {
	ID                uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	InvoiceID         uuid.UUID  `gorm:"type:uuid;not null" json:"invoice_id"`
	Description       string     `gorm:"not null" json:"description"`
	Quantity          float64    `gorm:"not null" json:"quantity"`
	Unit              string     `json:"unit"`
	UnitPrice         float64    `gorm:"not null" json:"unit_price"`
	TotalPrice        float64    `gorm:"not null" json:"total_price"`
	TaxRate           float64    `json:"tax_rate"`
	TaxAmount         float64    `json:"tax_amount"`
	
	// Reference
	OrderItemID       *uuid.UUID `gorm:"type:uuid" json:"order_item_id"`
	InventoryID       *uuid.UUID `gorm:"type:uuid" json:"inventory_id"`
	
	// Timestamps
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	
	// Relations
	Invoice           *Invoice   `gorm:"foreignKey:InvoiceID" json:"invoice,omitempty"`
	OrderItem         *OrderItem `gorm:"foreignKey:OrderItemID" json:"order_item,omitempty"`
	Inventory         *Inventory `gorm:"foreignKey:InventoryID" json:"inventory,omitempty"`
}

func (ii *InvoiceItem) BeforeCreate(tx *gorm.DB) error {
	ii.ID = uuid.New()
	return nil
}

// Payment represents a payment transaction
type Payment struct {
	ID                uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID         uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	PaymentNo         string     `gorm:"not null;unique" json:"payment_no"`
	Type              string     `gorm:"not null" json:"type"`               // incoming, outgoing
	Status            string     `gorm:"not null" json:"status"`             // pending, completed, failed, cancelled
	
	// Reference
	InvoiceID         *uuid.UUID `gorm:"type:uuid" json:"invoice_id"`
	CustomerID        *uuid.UUID `gorm:"type:uuid" json:"customer_id"`
	SupplierID        *uuid.UUID `gorm:"type:uuid" json:"supplier_id"`
	
	// Payment Details
	PaymentDate       time.Time  `json:"payment_date"`
	Amount            float64    `gorm:"not null" json:"amount"`
	Currency          string     `gorm:"default:'USD'" json:"currency"`
	ExchangeRate      float64    `gorm:"default:1" json:"exchange_rate"`
	PaymentMethod     string     `gorm:"not null" json:"payment_method"`      // cash, check, bank_transfer, credit_card
	
	// Bank Details
	BankName          string     `json:"bank_name"`
	BankAccount       string     `json:"bank_account"`
	TransactionNo     string     `json:"transaction_no"`
	CheckNo           string     `json:"check_no"`
	
	// Additional Info
	Notes             string     `json:"notes"`
	AttachmentPath    string     `json:"attachment_path"`
	
	// Timestamps
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	CreatedBy         uuid.UUID  `gorm:"type:uuid;not null" json:"created_by"`
	
	// Relations
	Company           *Company   `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Invoice           *Invoice   `gorm:"foreignKey:InvoiceID" json:"invoice,omitempty"`
	Customer          *Customer  `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Supplier          *Supplier  `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	Creator           *User      `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	p.ID = uuid.New()
	return nil
}

// Expense represents an expense record
type Expense struct {
	ID                uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID         uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	ExpenseNo         string     `gorm:"not null;unique" json:"expense_no"`
	Category          string     `gorm:"not null" json:"category"`           // material, labor, overhead, utility, maintenance, etc.
	SubCategory       string     `json:"sub_category"`
	Status            string     `gorm:"not null" json:"status"`             // draft, submitted, approved, paid, rejected
	
	// Reference
	SupplierID        *uuid.UUID `gorm:"type:uuid" json:"supplier_id"`
	DepartmentID      *uuid.UUID `gorm:"type:uuid" json:"department_id"`
	ProjectID         *uuid.UUID `gorm:"type:uuid" json:"project_id"`
	
	// Expense Details
	ExpenseDate       time.Time  `json:"expense_date"`
	Description       string     `gorm:"not null" json:"description"`
	Amount            float64    `gorm:"not null" json:"amount"`
	TaxAmount         float64    `json:"tax_amount"`
	TotalAmount       float64    `gorm:"not null" json:"total_amount"`
	Currency          string     `gorm:"default:'USD'" json:"currency"`
	
	// Payment Info
	PaymentMethod     string     `json:"payment_method"`
	PaymentStatus     string     `json:"payment_status"`                      // unpaid, paid
	PaidDate          *time.Time `json:"paid_date"`
	PaidBy            *uuid.UUID `gorm:"type:uuid" json:"paid_by"`
	
	// Approval
	SubmittedBy       uuid.UUID  `gorm:"type:uuid;not null" json:"submitted_by"`
	SubmittedAt       time.Time  `json:"submitted_at"`
	ApprovedBy        *uuid.UUID `gorm:"type:uuid" json:"approved_by"`
	ApprovedAt        *time.Time `json:"approved_at"`
	
	// Additional Info
	ReceiptNo         string     `json:"receipt_no"`
	AttachmentPath    string     `json:"attachment_path"`
	Notes             string     `json:"notes"`
	
	// Timestamps
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	
	// Relations
	Company           *Company   `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Supplier          *Supplier  `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	Submitter         *User      `gorm:"foreignKey:SubmittedBy" json:"submitter,omitempty"`
	Approver          *User      `gorm:"foreignKey:ApprovedBy" json:"approver,omitempty"`
	Payer             *User      `gorm:"foreignKey:PaidBy" json:"payer,omitempty"`
}

func (e *Expense) BeforeCreate(tx *gorm.DB) error {
	e.ID = uuid.New()
	return nil
}

// AccountReceivable represents AR records
type AccountReceivable struct {
	ID                uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID         uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	CustomerID        uuid.UUID  `gorm:"type:uuid;not null" json:"customer_id"`
	InvoiceID         uuid.UUID  `gorm:"type:uuid;not null" json:"invoice_id"`
	
	// Amounts
	InvoiceAmount     float64    `json:"invoice_amount"`
	PaidAmount        float64    `json:"paid_amount"`
	BalanceAmount     float64    `json:"balance_amount"`
	Currency          string     `json:"currency"`
	
	// Dates
	InvoiceDate       time.Time  `json:"invoice_date"`
	DueDate           time.Time  `json:"due_date"`
	LastPaymentDate   *time.Time `json:"last_payment_date"`
	
	// Aging
	DaysOverdue       int        `json:"days_overdue"`
	AgingCategory     string     `json:"aging_category"`                      // current, 30days, 60days, 90days, over90days
	
	// Status
	Status            string     `json:"status"`                              // open, partial, paid, written_off
	CollectionStatus  string     `json:"collection_status"`                   // normal, warning, critical
	
	// Timestamps
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	
	// Relations
	Company           *Company   `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Customer          *Customer  `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Invoice           *Invoice   `gorm:"foreignKey:InvoiceID" json:"invoice,omitempty"`
}

func (ar *AccountReceivable) BeforeCreate(tx *gorm.DB) error {
	ar.ID = uuid.New()
	return nil
}

// AccountPayable represents AP records
type AccountPayable struct {
	ID                uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID         uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	SupplierID        uuid.UUID  `gorm:"type:uuid;not null" json:"supplier_id"`
	InvoiceID         uuid.UUID  `gorm:"type:uuid;not null" json:"invoice_id"`
	
	// Amounts
	InvoiceAmount     float64    `json:"invoice_amount"`
	PaidAmount        float64    `json:"paid_amount"`
	BalanceAmount     float64    `json:"balance_amount"`
	Currency          string     `json:"currency"`
	
	// Dates
	InvoiceDate       time.Time  `json:"invoice_date"`
	DueDate           time.Time  `json:"due_date"`
	LastPaymentDate   *time.Time `json:"last_payment_date"`
	
	// Status
	Status            string     `json:"status"`                              // open, partial, paid, disputed
	PaymentPriority   string     `json:"payment_priority"`                    // high, medium, low
	
	// Timestamps
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	
	// Relations
	Company           *Company   `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Supplier          *Supplier  `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	Invoice           *Invoice   `gorm:"foreignKey:InvoiceID" json:"invoice,omitempty"`
}

func (ap *AccountPayable) BeforeCreate(tx *gorm.DB) error {
	ap.ID = uuid.New()
	return nil
}

// BankAccount represents company bank accounts
type BankAccount struct {
	ID                uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID         uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	AccountName       string     `gorm:"not null" json:"account_name"`
	AccountNo         string     `gorm:"not null" json:"account_no"`
	BankName          string     `gorm:"not null" json:"bank_name"`
	BankCode          string     `json:"bank_code"`
	BranchName        string     `json:"branch_name"`
	SwiftCode         string     `json:"swift_code"`
	
	// Account Info
	AccountType       string     `json:"account_type"`                        // checking, savings, credit
	Currency          string     `gorm:"default:'USD'" json:"currency"`
	CurrentBalance    float64    `json:"current_balance"`
	AvailableBalance  float64    `json:"available_balance"`
	
	// Status
	Status            string     `gorm:"default:'active'" json:"status"`       // active, inactive, closed
	IsDefault         bool       `gorm:"default:false" json:"is_default"`
	
	// Additional Info
	Notes             string     `json:"notes"`
	
	// Timestamps
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	
	// Relations
	Company           *Company   `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
}

func (ba *BankAccount) BeforeCreate(tx *gorm.DB) error {
	ba.ID = uuid.New()
	return nil
}

// FinancialPeriod represents accounting periods
type FinancialPeriod struct {
	ID                uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID         uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	Name              string     `gorm:"not null" json:"name"`
	Type              string     `gorm:"not null" json:"type"`               // monthly, quarterly, yearly
	StartDate         time.Time  `json:"start_date"`
	EndDate           time.Time  `json:"end_date"`
	Status            string     `gorm:"default:'open'" json:"status"`        // open, closed, locked
	IsCurrent         bool       `gorm:"default:false" json:"is_current"`
	
	// Timestamps
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	ClosedAt          *time.Time `json:"closed_at"`
	ClosedBy          *uuid.UUID `gorm:"type:uuid" json:"closed_by"`
	
	// Relations
	Company           *Company   `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Closer            *User      `gorm:"foreignKey:ClosedBy" json:"closer,omitempty"`
}

func (fp *FinancialPeriod) BeforeCreate(tx *gorm.DB) error {
	fp.ID = uuid.New()
	return nil
}