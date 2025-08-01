package model

import (
	"time"

	"github.com/google/uuid"
)

// Quote represents a quotation
type Quote struct {
	Base
	QuoteNo          string     `json:"quote_no" db:"quote_no"`
	InquiryID        uuid.UUID  `json:"inquiry_id" db:"inquiry_id"`
	CompanyID        uuid.UUID  `json:"company_id" db:"company_id"`
	CustomerID       uuid.UUID  `json:"customer_id" db:"customer_id"`
	EngineerID       uuid.UUID  `json:"engineer_id" db:"engineer_id"`
	Status           string     `json:"status" db:"status"`
	MaterialCost     float64    `json:"material_cost" db:"material_cost"`
	ProcessCost      float64    `json:"process_cost" db:"process_cost"`
	SurfaceCost      float64    `json:"surface_cost" db:"surface_cost"`
	HeatTreatCost    float64    `json:"heat_treat_cost" db:"heat_treat_cost"`
	PackagingCost    float64    `json:"packaging_cost" db:"packaging_cost"`
	ShippingCost     float64    `json:"shipping_cost" db:"shipping_cost"`
	TariffCost       float64    `json:"tariff_cost" db:"tariff_cost"`
	OverheadRate     float64    `json:"overhead_rate" db:"overhead_rate"`
	ProfitRate       float64    `json:"profit_rate" db:"profit_rate"`
	TotalCost        float64    `json:"total_cost" db:"total_cost"`
	UnitPrice        float64    `json:"unit_price" db:"unit_price"`
	Currency         string     `json:"currency" db:"currency"`
	ValidUntil       time.Time  `json:"valid_until" db:"valid_until"`
	DeliveryDays     int        `json:"delivery_days" db:"delivery_days"`
	PaymentTerms     string     `json:"payment_terms" db:"payment_terms"`
	Notes            *string    `json:"notes,omitempty" db:"notes"`
	ReviewedBy       *uuid.UUID `json:"reviewed_by,omitempty" db:"reviewed_by"`
	ReviewedAt       *time.Time `json:"reviewed_at,omitempty" db:"reviewed_at"`
	ApprovedBy       *uuid.UUID `json:"approved_by,omitempty" db:"approved_by"`
	ApprovedAt       *time.Time `json:"approved_at,omitempty" db:"approved_at"`
}