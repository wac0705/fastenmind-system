package model

import (
	"time"
	"github.com/google/uuid"
)

// CustomerFilter represents filter options for listing customers
type CustomerFilter struct {
	CompanyID uuid.UUID
	Search    string
	Country   string
	IsActive  bool
	Page      int
	PageSize  int
}

// CreateCustomerRequest represents the request to create a customer
type CreateCustomerRequest struct {
	CompanyID       uuid.UUID `json:"-"`
	CustomerCode    string    `json:"customer_code" validate:"required,min=1,max=50"`
	Name            string    `json:"name" validate:"required,min=1,max=200"`
	NameEn          *string   `json:"name_en,omitempty"`
	ShortName       *string   `json:"short_name,omitempty"`
	Country         string    `json:"country" validate:"required,len=2"`
	TaxID           *string   `json:"tax_id,omitempty"`
	Address         *string   `json:"address,omitempty"`
	ShippingAddress *string   `json:"shipping_address,omitempty"`
	ContactPerson   *string   `json:"contact_person,omitempty"`
	ContactPhone    *string   `json:"contact_phone,omitempty"`
	ContactEmail    *string   `json:"contact_email,omitempty" validate:"omitempty,email"`
	PaymentTerms    *string   `json:"payment_terms,omitempty"`
	CreditLimit     *float64  `json:"credit_limit,omitempty"`
	Currency        string    `json:"currency" validate:"required,len=3"`
	IsActive        bool      `json:"is_active"`
	CreatedBy       uuid.UUID `json:"-"`
}

// UpdateCustomerRequest represents the request to update a customer
type UpdateCustomerRequest struct {
	ID              uuid.UUID `json:"-"`
	CompanyID       uuid.UUID `json:"-"`
	CustomerCode    string    `json:"customer_code" validate:"required,min=1,max=50"`
	Name            string    `json:"name" validate:"required,min=1,max=200"`
	NameEn          *string   `json:"name_en,omitempty"`
	ShortName       *string   `json:"short_name,omitempty"`
	Country         string    `json:"country" validate:"required,len=2"`
	TaxID           *string   `json:"tax_id,omitempty"`
	Address         *string   `json:"address,omitempty"`
	ShippingAddress *string   `json:"shipping_address,omitempty"`
	ContactPerson   *string   `json:"contact_person,omitempty"`
	ContactPhone    *string   `json:"contact_phone,omitempty"`
	ContactEmail    *string   `json:"contact_email,omitempty" validate:"omitempty,email"`
	PaymentTerms    *string   `json:"payment_terms,omitempty"`
	CreditLimit     *float64  `json:"credit_limit,omitempty"`
	Currency        string    `json:"currency" validate:"required,len=3"`
	IsActive        bool      `json:"is_active"`
	UpdatedBy       uuid.UUID `json:"-"`
}

// CustomerStatistics represents customer statistics
type CustomerStatistics struct {
	TotalInquiries   int     `json:"total_inquiries"`
	TotalQuotes      int     `json:"total_quotes"`
	TotalOrders      int     `json:"total_orders"`
	TotalRevenue     float64 `json:"total_revenue"`
	CreditLimit      float64 `json:"credit_limit"`
	CreditUsed       float64 `json:"credit_used"`
	CreditAvailable  float64 `json:"credit_available"`
	LastOrderDate    *time.Time `json:"last_order_date,omitempty"`
	AverageOrderValue float64 `json:"average_order_value"`
}

// CreditHistory represents a customer's credit history entry
type CreditHistory struct {
	ID              uuid.UUID `json:"id"`
	CustomerID      uuid.UUID `json:"customer_id"`
	TransactionType string    `json:"transaction_type"` // order, payment, adjustment
	TransactionID   uuid.UUID `json:"transaction_id"`
	Amount          float64   `json:"amount"`
	Balance         float64   `json:"balance"`
	Description     string    `json:"description"`
	CreatedAt       time.Time `json:"created_at"`
}

// CustomerResponse represents the response with customer details
type CustomerResponse struct {
	Customer
	Statistics *CustomerStatistics `json:"statistics,omitempty"`
}