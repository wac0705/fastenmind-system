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