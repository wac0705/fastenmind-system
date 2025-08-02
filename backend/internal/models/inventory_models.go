package models

import (
	"time"
	"github.com/google/uuid"
)

// Material represents a material in inventory
type Material struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CompanyID       uuid.UUID  `json:"company_id"`
	MaterialCode    string     `json:"material_code"`
	MaterialName    string     `json:"material_name"`
	Category        string     `json:"category"`
	Unit            string     `json:"unit"`
	Description     string     `json:"description"`
	IsActive        bool       `json:"is_active"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// InventoryTransaction represents an inventory transaction
type InventoryTransaction struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	MaterialID      uuid.UUID  `json:"material_id"`
	TransactionType string     `json:"transaction_type"`
	Quantity        float64    `json:"quantity"`
	ReferenceType   string     `json:"reference_type"`
	ReferenceID     uuid.UUID  `json:"reference_id"`
	Notes           string     `json:"notes"`
	CreatedBy       uuid.UUID  `json:"created_by"`
	CreatedAt       time.Time  `json:"created_at"`
}

// Inventory represents inventory for a material
type Inventory struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	MaterialID      uuid.UUID  `json:"material_id"`
	Quantity        float64    `json:"quantity"`
	ReorderLevel    float64    `json:"reorder_level"`
	ReorderQuantity float64    `json:"reorder_quantity"`
	Location        string     `json:"location"`
	LastUpdated     time.Time  `json:"last_updated"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// Supplier represents a supplier
type Supplier struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CompanyID       uuid.UUID  `json:"company_id"`
	SupplierCode    string     `json:"supplier_code"`
	SupplierName    string     `json:"supplier_name"`
	ContactPerson   string     `json:"contact_person"`
	ContactEmail    string     `json:"contact_email"`
	ContactPhone    string     `json:"contact_phone"`
	Address         string     `json:"address"`
	Country         string     `json:"country"`
	IsActive        bool       `json:"is_active"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// PurchaseOrder represents a purchase order
type PurchaseOrder struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CompanyID       uuid.UUID  `json:"company_id"`
	PONumber        string     `json:"po_number"`
	SupplierID      uuid.UUID  `json:"supplier_id"`
	OrderDate       time.Time  `json:"order_date"`
	ExpectedDate    time.Time  `json:"expected_date"`
	TotalAmount     float64    `json:"total_amount"`
	Currency        string     `json:"currency"`
	Status          string     `json:"status"`
	Notes           string     `json:"notes"`
	CreatedBy       uuid.UUID  `json:"created_by"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}