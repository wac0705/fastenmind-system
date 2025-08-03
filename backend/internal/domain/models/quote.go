package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// QuoteStatus represents the status of a quote
type QuoteStatus string

const (
	QuoteStatusDraft     QuoteStatus = "draft"
	QuoteStatusSubmitted QuoteStatus = "submitted"
	QuoteStatusApproved  QuoteStatus = "approved"
	QuoteStatusRejected  QuoteStatus = "rejected"
	QuoteStatusExpired   QuoteStatus = "expired"
	QuoteStatusOrdered   QuoteStatus = "ordered"
)

// Quote represents a price quotation
type Quote struct {
	BaseModel
	QuoteNo            string          `gorm:"type:varchar(50);uniqueIndex;not null" json:"quote_no"`
	CompanyID          uuid.UUID       `gorm:"type:uuid;not null" json:"company_id"`
	Company            Company         `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	CustomerID         uuid.UUID       `gorm:"type:uuid;not null" json:"customer_id"`
	Customer           Customer        `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	InquiryID          *uuid.UUID      `gorm:"type:uuid" json:"inquiry_id"`
	Inquiry            *Inquiry        `gorm:"foreignKey:InquiryID" json:"inquiry,omitempty"`
	SalesID            uuid.UUID       `gorm:"type:uuid;not null" json:"sales_id"`
	Sales              Account         `gorm:"foreignKey:SalesID" json:"sales,omitempty"`
	EngineerID         uuid.UUID       `gorm:"type:uuid;not null" json:"engineer_id"`
	Engineer           Account         `gorm:"foreignKey:EngineerID" json:"engineer,omitempty"`
	Status             QuoteStatus     `gorm:"type:varchar(20);not null;default:'draft'" json:"status"`
	ValidityDays       int             `gorm:"default:30" json:"validity_days"`
	ExpiryDate         time.Time       `json:"expiry_date"`
	Currency           string          `gorm:"type:varchar(3);not null" json:"currency"`
	ExchangeRate       decimal.Decimal `gorm:"type:decimal(10,6)" json:"exchange_rate"`
	Incoterm           string          `gorm:"type:varchar(10);not null" json:"incoterm"`
	PaymentTerms       string          `gorm:"type:varchar(100)" json:"payment_terms"`
	DeliveryTerms      string          `gorm:"type:text" json:"delivery_terms"`
	SubTotal           decimal.Decimal `gorm:"type:decimal(15,2);not null" json:"sub_total"`
	DiscountPercent    decimal.Decimal `gorm:"type:decimal(5,2)" json:"discount_percent"`
	DiscountAmount     decimal.Decimal `gorm:"type:decimal(15,2)" json:"discount_amount"`
	TaxPercent         decimal.Decimal `gorm:"type:decimal(5,2)" json:"tax_percent"`
	TaxAmount          decimal.Decimal `gorm:"type:decimal(15,2)" json:"tax_amount"`
	ShippingCost       decimal.Decimal `gorm:"type:decimal(15,2)" json:"shipping_cost"`
	TotalAmount        decimal.Decimal `gorm:"type:decimal(15,2);not null" json:"total_amount"`
	ProfitMargin       decimal.Decimal `gorm:"type:decimal(5,2)" json:"profit_margin"`
	Notes              string          `gorm:"type:text" json:"notes"`
	InternalNotes      string          `gorm:"type:text" json:"internal_notes"`
	ApprovedBy         *uuid.UUID      `gorm:"type:uuid" json:"approved_by"`
	ApprovedAt         *time.Time      `json:"approved_at"`
	RejectionReason    string          `gorm:"type:text" json:"rejection_reason"`
	
	// Relationships
	Items      []QuoteItem      `gorm:"foreignKey:QuoteID" json:"items,omitempty"`
	Revisions  []QuoteRevision  `gorm:"foreignKey:QuoteID" json:"revisions,omitempty"`
}

func (Quote) TableName() string {
	return "quotes"
}

// QuoteItem represents an item in a quote
type QuoteItem struct {
	BaseModel
	QuoteID         uuid.UUID       `gorm:"type:uuid;not null" json:"quote_id"`
	Quote           Quote           `gorm:"foreignKey:QuoteID" json:"quote,omitempty"`
	ProductID       *uuid.UUID      `gorm:"type:uuid" json:"product_id"`
	Product         *Product        `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	ItemNo          int             `gorm:"not null" json:"item_no"`
	Description     string          `gorm:"type:text;not null" json:"description"`
	Specifications  string          `gorm:"type:jsonb" json:"specifications"`
	Quantity        decimal.Decimal `gorm:"type:decimal(15,3);not null" json:"quantity"`
	Unit            string          `gorm:"type:varchar(20);not null" json:"unit"`
	UnitPrice       decimal.Decimal `gorm:"type:decimal(15,4);not null" json:"unit_price"`
	DiscountPercent decimal.Decimal `gorm:"type:decimal(5,2)" json:"discount_percent"`
	Amount          decimal.Decimal `gorm:"type:decimal(15,2);not null" json:"amount"`
	CostPrice       decimal.Decimal `gorm:"type:decimal(15,4)" json:"cost_price"`
	ProfitMargin    decimal.Decimal `gorm:"type:decimal(5,2)" json:"profit_margin"`
	LeadTimeDays    int             `json:"lead_time_days"`
	Notes           string          `gorm:"type:text" json:"notes"`
}

func (QuoteItem) TableName() string {
	return "quote_items"
}

// QuoteRevision tracks quote version history
type QuoteRevision struct {
	BaseModel
	QuoteID      uuid.UUID `gorm:"type:uuid;not null" json:"quote_id"`
	Quote        Quote     `gorm:"foreignKey:QuoteID" json:"quote,omitempty"`
	RevisionNo   int       `gorm:"not null" json:"revision_no"`
	RevisionData string    `gorm:"type:jsonb;not null" json:"revision_data"`
	ChangedBy    uuid.UUID `gorm:"type:uuid;not null" json:"changed_by"`
	ChangedUser  Account   `gorm:"foreignKey:ChangedBy" json:"changed_user,omitempty"`
	ChangeReason string    `gorm:"type:text" json:"change_reason"`
}

func (QuoteRevision) TableName() string {
	return "quote_revisions"
}