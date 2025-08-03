package models

import (
	"time"
)

// ExchangeRate 匯率
type ExchangeRate struct {
	ID            string    `json:"id" gorm:"primaryKey"`
	CompanyID     string    `json:"company_id" gorm:"index"`
	FromCurrency  string    `json:"from_currency"`
	ToCurrency    string    `json:"to_currency"`
	Rate          float64   `json:"rate"`
	RateType      string    `json:"rate_type"` // buy, sell, mid, official
	EffectiveDate time.Time `json:"effective_date"`
	ValidDate     time.Time `json:"valid_date"`
	Source        string    `json:"source"` // manual, api, bank
	IsActive      bool      `json:"is_active" gorm:"default:true"`
	CreatedBy     string    `json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (ExchangeRate) TableName() string {
	return "exchange_rates"
}