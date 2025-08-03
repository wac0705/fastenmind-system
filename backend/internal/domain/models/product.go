package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ProductStatus represents the status of a product
type ProductStatus string

const (
	ProductStatusActive       ProductStatus = "active"
	ProductStatusDiscontinued ProductStatus = "discontinued"
	ProductStatusDraft        ProductStatus = "draft"
)

// Product represents a product entity
type Product struct {
	BaseModel
	CompanyID       uuid.UUID       `gorm:"type:uuid;not null" json:"company_id"`
	Company         Company         `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	ProductCode     string          `gorm:"type:varchar(50);not null" json:"product_code"`
	Name            string          `gorm:"type:varchar(200);not null" json:"name"`
	NameEn          string          `gorm:"type:varchar(200)" json:"name_en"`
	Category        string          `gorm:"type:varchar(50);not null" json:"category"`
	SubCategory     string          `gorm:"type:varchar(50)" json:"sub_category"`
	Description     string          `gorm:"type:text" json:"description"`
	Specifications  string          `gorm:"type:jsonb" json:"specifications"`
	Unit            string          `gorm:"type:varchar(20);not null" json:"unit"`
	Weight          decimal.Decimal `gorm:"type:decimal(10,3)" json:"weight"`
	WeightUnit      string          `gorm:"type:varchar(10)" json:"weight_unit"`
	HSCode          string          `gorm:"type:varchar(20)" json:"hs_code"`
	StandardCost    decimal.Decimal `gorm:"type:decimal(15,4)" json:"standard_cost"`
	Currency        string          `gorm:"type:varchar(3);default:'USD'" json:"currency"`
	LeadTimeDays    int             `json:"lead_time_days"`
	MOQ             int             `json:"moq"` // Minimum Order Quantity
	Status          ProductStatus   `gorm:"type:varchar(20);default:'active'" json:"status"`
	DrawingFiles    []string        `gorm:"type:text[]" json:"drawing_files"`
	
	// Relationships
	InquiryItems []InquiryItem `gorm:"foreignKey:ProductID" json:"inquiry_items,omitempty"`
	QuoteItems   []QuoteItem   `gorm:"foreignKey:ProductID" json:"quote_items,omitempty"`
	OrderItems   []OrderItem   `gorm:"foreignKey:ProductID" json:"order_items,omitempty"`
}

func (Product) TableName() string {
	return "products"
}

// ProductProcess represents manufacturing process for a product
type ProductProcess struct {
	BaseModel
	ProductID      uuid.UUID       `gorm:"type:uuid;not null" json:"product_id"`
	Product        Product         `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	ProcessID      uuid.UUID       `gorm:"type:uuid;not null" json:"process_id"`
	Process        Process         `gorm:"foreignKey:ProcessID" json:"process,omitempty"`
	Sequence       int             `gorm:"not null" json:"sequence"`
	CycleTime      decimal.Decimal `gorm:"type:decimal(10,2)" json:"cycle_time"`
	SetupTime      decimal.Decimal `gorm:"type:decimal(10,2)" json:"setup_time"`
	CostPerUnit    decimal.Decimal `gorm:"type:decimal(15,4)" json:"cost_per_unit"`
	Notes          string          `gorm:"type:text" json:"notes"`
}

func (ProductProcess) TableName() string {
	return "product_processes"
}