package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Customer represents a customer entity
type Customer struct {
	BaseModel
	CompanyID        uuid.UUID       `gorm:"type:uuid;not null" json:"company_id"`
	Company          Company         `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	CustomerCode     string          `gorm:"type:varchar(50);not null" json:"customer_code"`
	Name             string          `gorm:"type:varchar(100);not null" json:"name"`
	NameEn           string          `gorm:"type:varchar(100)" json:"name_en"`
	ShortName        string          `gorm:"type:varchar(50)" json:"short_name"`
	Country          string          `gorm:"type:varchar(2);not null" json:"country"`
	TaxID            string          `gorm:"type:varchar(50)" json:"tax_id"`
	Address          string          `gorm:"type:text" json:"address"`
	ShippingAddress  string          `gorm:"type:text" json:"shipping_address"`
	ContactPerson    string          `gorm:"type:varchar(100)" json:"contact_person"`
	ContactPhone     string          `gorm:"type:varchar(50)" json:"contact_phone"`
	ContactEmail     string          `gorm:"type:varchar(100)" json:"contact_email"`
	PaymentTerms     string          `gorm:"type:varchar(50)" json:"payment_terms"`
	CreditLimit      decimal.Decimal `gorm:"type:decimal(15,2)" json:"credit_limit"`
	Currency         string          `gorm:"type:varchar(3);not null;default:'USD'" json:"currency"`
	IsActive         bool            `gorm:"default:true" json:"is_active"`
	
	// Relationships
	Inquiries []Inquiry `gorm:"foreignKey:CustomerID" json:"inquiries,omitempty"`
	Orders    []Order   `gorm:"foreignKey:CustomerID" json:"orders,omitempty"`
	Quotes    []Quote   `gorm:"foreignKey:CustomerID" json:"quotes,omitempty"`
}

func (Customer) TableName() string {
	return "customers"
}

// CustomerTransactionTerms represents customer-specific transaction terms per company
type CustomerTransactionTerms struct {
	BaseModel
	CustomerID      uuid.UUID `gorm:"type:uuid;not null" json:"customer_id"`
	Customer        Customer  `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	CompanyID       uuid.UUID `gorm:"type:uuid;not null" json:"company_id"`
	Company         Company   `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Incoterm        string    `gorm:"type:varchar(10)" json:"incoterm"`
	Currency        string    `gorm:"type:varchar(3)" json:"currency"`
	Port            string    `gorm:"type:varchar(100)" json:"port"`
	Country         string    `gorm:"type:varchar(2)" json:"country"`
	PaymentTermDays int       `json:"payment_term_days"`
	Notes           string    `gorm:"type:text" json:"notes"`
}

func (CustomerTransactionTerms) TableName() string {
	return "customer_transaction_terms"
}