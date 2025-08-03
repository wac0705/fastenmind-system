package models

import "github.com/google/uuid"

// TariffCalculationParams represents parameters for tariff calculation
type TariffCalculationParams struct {
	CompanyID     uuid.UUID `json:"company_id"`
	HSCode        string    `json:"hs_code"`
	OriginCountry string    `json:"origin_country"`
	DestCountry   string    `json:"dest_country"`
	Quantity      float64   `json:"quantity"`
	UnitValue     float64   `json:"unit_value"`
	Currency      string    `json:"currency"`
}

// TariffCalculationResult represents the result of tariff calculation
type TariffCalculationResult struct {
	HSCode         string  `json:"hs_code"`
	Description    string  `json:"description"`
	OriginCountry  string  `json:"origin_country"`
	DestCountry    string  `json:"dest_country"`
	Quantity       float64 `json:"quantity"`
	Unit           string  `json:"unit"`
	UnitValue      float64 `json:"unit_value"`
	TotalValue     float64 `json:"total_value"`
	Currency       string  `json:"currency"`
	TariffRate     float64 `json:"tariff_rate"`
	RateType       string  `json:"rate_type"`
	BaseDuty       float64 `json:"base_duty"`
	VATRate        float64 `json:"vat_rate"`
	VAT            float64 `json:"vat"`
	ExciseTaxRate  float64 `json:"excise_tax_rate"`
	ExciseTax      float64 `json:"excise_tax"`
	TotalTax       float64 `json:"total_tax"`
	TotalAmount    float64 `json:"total_amount"`
	TradeAgreement string  `json:"trade_agreement,omitempty"`
}