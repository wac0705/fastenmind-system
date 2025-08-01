package model

import "github.com/google/uuid"

// Customer represents a customer
type Customer struct {
	Base
	CompanyID       uuid.UUID `json:"company_id" db:"company_id"`
	CustomerCode    string    `json:"customer_code" db:"customer_code"`
	Name            string    `json:"name" db:"name"`
	NameEn          *string   `json:"name_en,omitempty" db:"name_en"`
	ShortName       *string   `json:"short_name,omitempty" db:"short_name"`
	Country         string    `json:"country" db:"country"`
	TaxID           *string   `json:"tax_id,omitempty" db:"tax_id"`
	Address         *string   `json:"address,omitempty" db:"address"`
	ShippingAddress *string   `json:"shipping_address,omitempty" db:"shipping_address"`
	ContactPerson   *string   `json:"contact_person,omitempty" db:"contact_person"`
	ContactPhone    *string   `json:"contact_phone,omitempty" db:"contact_phone"`
	ContactEmail    *string   `json:"contact_email,omitempty" db:"contact_email"`
	PaymentTerms    *string   `json:"payment_terms,omitempty" db:"payment_terms"`
	CreditLimit     *float64  `json:"credit_limit,omitempty" db:"credit_limit"`
	Currency        string    `json:"currency" db:"currency"`
	IsActive        bool      `json:"is_active" db:"is_active"`
	
	// Relations
	Company *Company `json:"company,omitempty"`
}