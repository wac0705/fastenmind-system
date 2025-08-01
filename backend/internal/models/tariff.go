package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// HSCode represents harmonized system codes for products
type HSCode struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Code          string    `gorm:"not null;uniqueIndex" json:"code"`
	Description   string    `json:"description"`
	DescriptionEN string    `json:"description_en"`
	Unit          string    `json:"unit"`
	Category      string    `json:"category"`
	ParentCode    string    `json:"parent_code,omitempty"`
	IsActive      bool      `gorm:"default:true" json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (h *HSCode) BeforeCreate(tx *gorm.DB) error {
	h.ID = uuid.New()
	return nil
}

// TariffRate represents tariff rates between countries
type TariffRate struct {
	ID                     uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	HSCode                 string     `gorm:"not null;index" json:"hs_code"`
	FromCountry            string     `gorm:"not null;index" json:"from_country"`
	ToCountry              string     `gorm:"not null;index" json:"to_country"`
	RateType               string     `gorm:"not null" json:"rate_type"` // ad_valorem, specific, compound
	RateValue              float64    `json:"rate_value"`                 // percentage for ad_valorem
	SpecificRate           float64    `json:"specific_rate,omitempty"`    // amount per unit for specific
	Currency               string     `json:"currency,omitempty"`
	Unit                   string     `json:"unit,omitempty"`
	EffectiveFrom          time.Time  `json:"effective_from"`
	EffectiveTo            *time.Time `json:"effective_to,omitempty"`
	PreferentialRate       float64    `json:"preferential_rate,omitempty"`
	PreferentialConditions string     `json:"preferential_conditions,omitempty"`
	Notes                  string     `json:"notes,omitempty"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
}

func (t *TariffRate) BeforeCreate(tx *gorm.DB) error {
	t.ID = uuid.New()
	return nil
}

// TradeAgreement represents trade agreements between countries
type TradeAgreement struct {
	ID                  uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	AgreementName       string    `gorm:"not null" json:"agreement_name"`
	MemberCountries     []string  `gorm:"type:text[]" json:"member_countries"`
	EffectiveDate       time.Time `json:"effective_date"`
	PreferentialRates   bool      `json:"preferential_rates"`
	CertificateRequired bool      `json:"certificate_required"`
	Description         string    `json:"description,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

func (t *TradeAgreement) BeforeCreate(tx *gorm.DB) error {
	t.ID = uuid.New()
	return nil
}

// TariffCalculation represents a tariff calculation record
type TariffCalculation struct {
	ID                     uuid.UUID   `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID              uuid.UUID   `gorm:"type:uuid;not null" json:"company_id"`
	UserID                 uuid.UUID   `gorm:"type:uuid;not null" json:"user_id"`
	QuoteID                *uuid.UUID  `gorm:"type:uuid" json:"quote_id,omitempty"`
	HSCode                 string      `gorm:"not null" json:"hs_code"`
	FromCountry            string      `gorm:"not null" json:"from_country"`
	ToCountry              string      `gorm:"not null" json:"to_country"`
	ProductValue           float64     `json:"product_value"`
	Quantity               float64     `json:"quantity"`
	Unit                   string      `json:"unit"`
	WeightKG               float64     `json:"weight_kg,omitempty"`
	Currency               string      `json:"currency"`
	Incoterm               string      `json:"incoterm,omitempty"`
	PreferentialTreatment  bool        `json:"preferential_treatment"`
	TariffRateID           *uuid.UUID  `gorm:"type:uuid" json:"tariff_rate_id,omitempty"`
	CalculatedTariff       float64     `json:"calculated_tariff"`
	EffectiveRate          float64     `json:"effective_rate"`
	CalculationDetails     interface{} `gorm:"type:jsonb" json:"calculation_details"`
	Warnings               []string    `gorm:"type:text[]" json:"warnings,omitempty"`
	CreatedAt              time.Time   `json:"created_at"`
	
	// Relations
	Company    *Company    `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	User       *User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Quote      *Quote      `gorm:"foreignKey:QuoteID" json:"quote,omitempty"`
	TariffRate *TariffRate `gorm:"foreignKey:TariffRateID" json:"tariff_rate,omitempty"`
}

func (t *TariffCalculation) BeforeCreate(tx *gorm.DB) error {
	t.ID = uuid.New()
	return nil
}