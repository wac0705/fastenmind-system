package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// InvoiceStatus represents the status of an invoice
type InvoiceStatus string

const (
	InvoiceStatusDraft     InvoiceStatus = "draft"
	InvoiceStatusIssued    InvoiceStatus = "issued"
	InvoiceStatusPaid      InvoiceStatus = "paid"
	InvoiceStatusOverdue   InvoiceStatus = "overdue"
	InvoiceStatusCancelled InvoiceStatus = "cancelled"
)

// Invoice represents a sales invoice
type Invoice struct {
	BaseModel
	InvoiceNo       string          `gorm:"type:varchar(50);uniqueIndex;not null" json:"invoice_no"`
	OrderID         uuid.UUID       `gorm:"type:uuid;not null" json:"order_id"`
	Order           Order           `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	CompanyID       uuid.UUID       `gorm:"type:uuid;not null" json:"company_id"`
	Company         Company         `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	CustomerID      uuid.UUID       `gorm:"type:uuid;not null" json:"customer_id"`
	Customer        Customer        `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Status          InvoiceStatus   `gorm:"type:varchar(20);not null;default:'draft'" json:"status"`
	InvoiceDate     time.Time       `json:"invoice_date"`
	DueDate         time.Time       `json:"due_date"`
	Currency        string          `gorm:"type:varchar(3);not null" json:"currency"`
	ExchangeRate    decimal.Decimal `gorm:"type:decimal(10,6)" json:"exchange_rate"`
	SubTotal        decimal.Decimal `gorm:"type:decimal(15,2);not null" json:"sub_total"`
	DiscountPercent decimal.Decimal `gorm:"type:decimal(5,2)" json:"discount_percent"`
	DiscountAmount  decimal.Decimal `gorm:"type:decimal(15,2)" json:"discount_amount"`
	TaxPercent      decimal.Decimal `gorm:"type:decimal(5,2)" json:"tax_percent"`
	TaxAmount       decimal.Decimal `gorm:"type:decimal(15,2)" json:"tax_amount"`
	ShippingCost    decimal.Decimal `gorm:"type:decimal(15,2)" json:"shipping_cost"`
	TotalAmount     decimal.Decimal `gorm:"type:decimal(15,2);not null" json:"total_amount"`
	PaidAmount      decimal.Decimal `gorm:"type:decimal(15,2)" json:"paid_amount"`
	BalanceAmount   decimal.Decimal `gorm:"type:decimal(15,2)" json:"balance_amount"`
	PaymentTerms    string          `gorm:"type:varchar(100)" json:"payment_terms"`
	Notes           string          `gorm:"type:text" json:"notes"`
	
	// Relationships
	Items    []InvoiceItem    `gorm:"foreignKey:InvoiceID" json:"items,omitempty"`
	Payments []Payment        `gorm:"foreignKey:InvoiceID" json:"payments,omitempty"`
}

func (Invoice) TableName() string {
	return "invoices"
}

// InvoiceItem represents an item in an invoice
type InvoiceItem struct {
	BaseModel
	InvoiceID       uuid.UUID       `gorm:"type:uuid;not null" json:"invoice_id"`
	Invoice         Invoice         `gorm:"foreignKey:InvoiceID" json:"invoice,omitempty"`
	OrderItemID     *uuid.UUID      `gorm:"type:uuid" json:"order_item_id"`
	OrderItem       *OrderItem      `gorm:"foreignKey:OrderItemID" json:"order_item,omitempty"`
	ItemNo          int             `gorm:"not null" json:"item_no"`
	Description     string          `gorm:"type:text;not null" json:"description"`
	Quantity        decimal.Decimal `gorm:"type:decimal(15,3);not null" json:"quantity"`
	Unit            string          `gorm:"type:varchar(20);not null" json:"unit"`
	UnitPrice       decimal.Decimal `gorm:"type:decimal(15,4);not null" json:"unit_price"`
	DiscountPercent decimal.Decimal `gorm:"type:decimal(5,2)" json:"discount_percent"`
	Amount          decimal.Decimal `gorm:"type:decimal(15,2);not null" json:"amount"`
	TaxPercent      decimal.Decimal `gorm:"type:decimal(5,2)" json:"tax_percent"`
	TaxAmount       decimal.Decimal `gorm:"type:decimal(15,2)" json:"tax_amount"`
	Notes           string          `gorm:"type:text" json:"notes"`
}

func (InvoiceItem) TableName() string {
	return "invoice_items"
}

// Payment represents a payment record
type Payment struct {
	BaseModel
	PaymentNo       string          `gorm:"type:varchar(50);uniqueIndex;not null" json:"payment_no"`
	InvoiceID       uuid.UUID       `gorm:"type:uuid;not null" json:"invoice_id"`
	Invoice         Invoice         `gorm:"foreignKey:InvoiceID" json:"invoice,omitempty"`
	CompanyID       uuid.UUID       `gorm:"type:uuid;not null" json:"company_id"`
	Company         Company         `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	PaymentDate     time.Time       `json:"payment_date"`
	Amount          decimal.Decimal `gorm:"type:decimal(15,2);not null" json:"amount"`
	Currency        string          `gorm:"type:varchar(3);not null" json:"currency"`
	ExchangeRate    decimal.Decimal `gorm:"type:decimal(10,6)" json:"exchange_rate"`
	PaymentMethod   string          `gorm:"type:varchar(50)" json:"payment_method"`
	ReferenceNo     string          `gorm:"type:varchar(100)" json:"reference_no"`
	BankName        string          `gorm:"type:varchar(100)" json:"bank_name"`
	Notes           string          `gorm:"type:text" json:"notes"`
}

func (Payment) TableName() string {
	return "payments"
}