package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// QuoteVersion represents a version history of a quote
type QuoteVersion struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	QuoteID       uuid.UUID `gorm:"type:uuid;not null" json:"quote_id"`
	VersionNumber int       `gorm:"not null" json:"version_number"`
	
	// Cost snapshot
	MaterialCost    float64 `json:"material_cost"`
	ProcessCost     float64 `json:"process_cost"`
	SurfaceCost     float64 `json:"surface_cost"`
	HeatTreatCost   float64 `json:"heat_treat_cost"`
	PackagingCost   float64 `json:"packaging_cost"`
	ShippingCost    float64 `json:"shipping_cost"`
	TariffCost      float64 `json:"tariff_cost"`
	OverheadRate    float64 `json:"overhead_rate"`
	ProfitRate      float64 `json:"profit_rate"`
	TotalCost       float64 `json:"total_cost"`
	UnitPrice       float64 `json:"unit_price"`
	
	// Version info
	ChangeSummary   string    `json:"change_summary,omitempty"`
	CreatedBy       uuid.UUID `gorm:"type:uuid;not null" json:"created_by_id"`
	CreatedAt       time.Time `json:"created_at"`
	
	// Relations
	Quote           *Quote    `gorm:"foreignKey:QuoteID" json:"quote,omitempty"`
	Creator         *User     `gorm:"foreignKey:CreatedBy" json:"created_by,omitempty"`
}

func (qv *QuoteVersion) BeforeCreate(tx *gorm.DB) error {
	qv.ID = uuid.New()
	return nil
}

// QuoteActivity represents activity log for a quote
type QuoteActivity struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	QuoteID     uuid.UUID `gorm:"type:uuid;not null" json:"quote_id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Action      string    `gorm:"not null" json:"action"` // created, updated, submitted, reviewed, approved, rejected, sent
	Description string    `json:"description,omitempty"`
	Metadata    JSONB     `gorm:"type:jsonb" json:"metadata,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	
	// Relations
	Quote       *Quote    `gorm:"foreignKey:QuoteID" json:"quote,omitempty"`
	User        *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (qa *QuoteActivity) BeforeCreate(tx *gorm.DB) error {
	qa.ID = uuid.New()
	return nil
}