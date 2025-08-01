package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Quote represents a quotation
type Quote struct {
	ID               uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	QuoteNo          string     `gorm:"not null;unique" json:"quote_no"`
	InquiryID        uuid.UUID  `gorm:"type:uuid;not null" json:"inquiry_id"`
	CompanyID        uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	CustomerID       uuid.UUID  `gorm:"type:uuid;not null" json:"customer_id"`
	SalesID          uuid.UUID  `gorm:"type:uuid;not null" json:"sales_id"`
	EngineerID       uuid.UUID  `gorm:"type:uuid;not null" json:"engineer_id"`
	Status           string     `gorm:"not null" json:"status"` // draft, pending_review, under_review, approved, sent, accepted, rejected, expired, cancelled
	
	// Costs
	MaterialCost     float64    `json:"material_cost"`
	ProcessCost      float64    `json:"process_cost"`
	SurfaceCost      float64    `json:"surface_cost"`
	HeatTreatCost    float64    `json:"heat_treat_cost"`
	PackagingCost    float64    `json:"packaging_cost"`
	ShippingCost     float64    `json:"shipping_cost"`
	TariffCost       float64    `json:"tariff_cost"`
	
	// Pricing
	OverheadRate     float64    `json:"overhead_rate"`
	ProfitRate       float64    `json:"profit_rate"`
	TotalCost        float64    `json:"total_cost"`
	UnitPrice        float64    `json:"unit_price"`
	Currency         string     `gorm:"default:'USD'" json:"currency"`
	
	// Terms
	ValidUntil       time.Time  `json:"valid_until"`
	ValidityDays     int        `json:"validity_days"`
	DeliveryDays     int        `json:"delivery_days"`
	DeliveryTerms    string     `json:"delivery_terms"`
	PaymentTerms     string     `json:"payment_terms"`
	Remarks          *string    `json:"remarks,omitempty"`
	Notes            string     `json:"notes,omitempty"`
	
	// Additional fields
	CreatedBy        uuid.UUID  `gorm:"type:uuid" json:"created_by"`
	TemplateID       *uuid.UUID `gorm:"type:uuid" json:"template_id,omitempty"`
	CurrentVersionID *uuid.UUID `gorm:"type:uuid" json:"current_version_id,omitempty"`
	TotalAmount      float64    `json:"total_amount"`
	
	// Workflow fields
	SubmittedAt      *time.Time `json:"submitted_at,omitempty"`
	ReviewedAt       *time.Time `json:"reviewed_at,omitempty"`
	ReviewerID       *uuid.UUID `gorm:"type:uuid" json:"reviewer_id,omitempty"`
	ReviewComments   string     `json:"review_comments,omitempty"`
	SentAt           *time.Time `json:"sent_at,omitempty"`
	SentByID         *uuid.UUID `gorm:"type:uuid" json:"sent_by_id,omitempty"`
	
	// Timestamps
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	
	// Relations
	Inquiry          *Inquiry   `gorm:"foreignKey:InquiryID" json:"inquiry,omitempty"`
	Company          *Company   `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Customer         *Customer  `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Sales            *User      `gorm:"foreignKey:SalesID" json:"sales,omitempty"`
	Engineer         *User      `gorm:"foreignKey:EngineerID" json:"engineer,omitempty"`
	Reviewer         *User      `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`
	SentBy           *User      `gorm:"foreignKey:SentByID" json:"sent_by,omitempty"`
}

func (q *Quote) BeforeCreate(tx *gorm.DB) error {
	q.ID = uuid.New()
	return nil
}