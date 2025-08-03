package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// InquiryStatus represents the status of an inquiry
type InquiryStatus string

const (
	InquiryStatusPending   InquiryStatus = "pending"
	InquiryStatusAssigned  InquiryStatus = "assigned"
	InquiryStatusQuoted    InquiryStatus = "quoted"
	InquiryStatusRejected  InquiryStatus = "rejected"
	InquiryStatusCancelled InquiryStatus = "cancelled"
)

// Inquiry represents a customer inquiry
type Inquiry struct {
	BaseModel
	InquiryNo           string        `gorm:"type:varchar(50);uniqueIndex;not null" json:"inquiry_no"`
	CompanyID           uuid.UUID     `gorm:"type:uuid;not null" json:"company_id"`
	Company             Company       `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	CustomerID          uuid.UUID     `gorm:"type:uuid;not null" json:"customer_id"`
	Customer            Customer      `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	SalesID             uuid.UUID     `gorm:"type:uuid;not null" json:"sales_id"`
	Sales               Account       `gorm:"foreignKey:SalesID" json:"sales,omitempty"`
	Status              InquiryStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	ProductCategory     string        `gorm:"type:varchar(50);not null" json:"product_category"`
	ProductName         string        `gorm:"type:varchar(200);not null" json:"product_name"`
	DrawingFiles        []string      `gorm:"type:text[]" json:"drawing_files"`
	Quantity            int           `gorm:"not null" json:"quantity"`
	Unit                string        `gorm:"type:varchar(20);not null" json:"unit"`
	RequiredDate        time.Time     `json:"required_date"`
	Incoterm            string        `gorm:"type:varchar(10);not null" json:"incoterm"`
	DestinationPort     string        `gorm:"type:varchar(100)" json:"destination_port"`
	DestinationAddress  string        `gorm:"type:text" json:"destination_address"`
	PaymentTerms        string        `gorm:"type:varchar(100)" json:"payment_terms"`
	SpecialRequirements string        `gorm:"type:text" json:"special_requirements"`
	AssignedEngineerID  *uuid.UUID    `gorm:"type:uuid" json:"assigned_engineer_id"`
	AssignedEngineer    *Account      `gorm:"foreignKey:AssignedEngineerID" json:"assigned_engineer,omitempty"`
	AssignedAt          *time.Time    `json:"assigned_at"`
	QuoteID             *uuid.UUID    `gorm:"type:uuid" json:"quote_id"`
	Quote               *Quote        `gorm:"foreignKey:QuoteID" json:"quote,omitempty"`
	QuotedAt            *time.Time    `json:"quoted_at"`
	
	// Relationships
	Items []InquiryItem `gorm:"foreignKey:InquiryID" json:"items,omitempty"`
}

func (Inquiry) TableName() string {
	return "inquiries"
}

// InquiryItem represents an item in an inquiry
type InquiryItem struct {
	BaseModel
	InquiryID      uuid.UUID       `gorm:"type:uuid;not null" json:"inquiry_id"`
	Inquiry        Inquiry         `gorm:"foreignKey:InquiryID" json:"inquiry,omitempty"`
	ProductID      *uuid.UUID      `gorm:"type:uuid" json:"product_id"`
	Product        *Product        `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	ItemNo         int             `gorm:"not null" json:"item_no"`
	Description    string          `gorm:"type:text;not null" json:"description"`
	Specifications string          `gorm:"type:jsonb" json:"specifications"`
	Quantity       decimal.Decimal `gorm:"type:decimal(15,3);not null" json:"quantity"`
	Unit           string          `gorm:"type:varchar(20);not null" json:"unit"`
	TargetPrice    decimal.Decimal `gorm:"type:decimal(15,4)" json:"target_price"`
	Currency       string          `gorm:"type:varchar(3)" json:"currency"`
	DrawingFiles   []string        `gorm:"type:text[]" json:"drawing_files"`
	Notes          string          `gorm:"type:text" json:"notes"`
}

func (InquiryItem) TableName() string {
	return "inquiry_items"
}