package model

import (
	"time"

	"github.com/google/uuid"
)

// Product represents a product in the system
type Product struct {
	Base
	CompanyID    uuid.UUID `json:"company_id" db:"company_id"`
	ProductCode  string    `json:"product_code" db:"product_code"`
	ProductName  string    `json:"product_name" db:"product_name"`
	Description  string    `json:"description" db:"description"`
	Category     string    `json:"category" db:"category"`
	Specification string   `json:"specification" db:"specification"`
	Material     string    `json:"material" db:"material"`
	Weight       float64   `json:"weight" db:"weight"`
	Unit         string    `json:"unit" db:"unit"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	
	// Relations
	Company *Company `json:"company,omitempty"`
}