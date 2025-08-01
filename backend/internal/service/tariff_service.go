package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/fastenmind/fastener-api/internal/models"
	"github.com/fastenmind/fastener-api/internal/repository"
	"github.com/google/uuid"
)

type TariffService interface {
	// HS Codes
	SearchHSCodes(params map[string]interface{}) ([]models.HSCode, int64, error)
	GetHSCode(code string) (*models.HSCode, error)
	ValidateHSCode(code, country string) (bool, string, []models.HSCode, error)
	GetCommonHSCodes(category string) ([]models.HSCode, error)
	
	// Tariff Calculation
	CalculateTariff(req TariffCalculationRequest) (*TariffCalculationResult, error)
	BatchCalculateTariff(items []TariffCalculationRequest) ([]TariffCalculationResult, error)
	
	// Trade Agreements
	GetTradeAgreements(countries []string) ([]models.TradeAgreement, error)
	
	// History
	GetCalculationHistory(companyID uuid.UUID, limit int) ([]models.TariffCalculation, error)
}

type TariffCalculationRequest struct {
	CompanyID              uuid.UUID `json:"-"`
	UserID                 uuid.UUID `json:"-"`
	HSCode                 string    `json:"hs_code"`
	FromCountry            string    `json:"from_country"`
	ToCountry              string    `json:"to_country"`
	ProductValue           float64   `json:"product_value"`
	Quantity               float64   `json:"quantity"`
	Unit                   string    `json:"unit"`
	WeightKG               float64   `json:"weight_kg,omitempty"`
	Currency               string    `json:"currency"`
	Incoterm               string    `json:"incoterm,omitempty"`
	PreferentialTreatment  bool      `json:"preferential_treatment"`
}

type TariffCalculationResult struct {
	HSCode             string                 `json:"hs_code"`
	FromCountry        string                 `json:"from_country"`
	ToCountry          string                 `json:"to_country"`
	ProductValue       float64                `json:"product_value"`
	TariffRate         *models.TariffRate     `json:"tariff_rate"`
	CalculatedTariff   float64                `json:"calculated_tariff"`
	EffectiveRate      float64                `json:"effective_rate"`
	Currency           string                 `json:"currency"`
	CalculationDetails CalculationDetails     `json:"calculation_details"`
	Warnings           []string               `json:"warnings,omitempty"`
}

type CalculationDetails struct {
	BaseValue             float64 `json:"base_value"`
	AdValoremDuty         float64 `json:"ad_valorem_duty,omitempty"`
	SpecificDuty          float64 `json:"specific_duty,omitempty"`
	TotalDuty             float64 `json:"total_duty"`
	PreferentialApplied   bool    `json:"preferential_applied"`
	PreferentialSavings   float64 `json:"preferential_savings,omitempty"`
}

type tariffService struct {
	repo repository.TariffRepository
}

func NewTariffService(repo repository.TariffRepository) TariffService {
	return &tariffService{repo: repo}
}

func (s *tariffService) SearchHSCodes(params map[string]interface{}) ([]models.HSCode, int64, error) {
	return s.repo.FindHSCodes(params)
}

func (s *tariffService) GetHSCode(code string) (*models.HSCode, error) {
	return s.repo.GetHSCode(code)
}

func (s *tariffService) ValidateHSCode(code, country string) (bool, string, []models.HSCode, error) {
	// Check if HS code exists
	hsCode, err := s.repo.GetHSCode(code)
	if err != nil {
		// If not found, suggest similar codes
		suggestions, _, _ := s.repo.FindHSCodes(map[string]interface{}{
			"search":    code[:4], // Search by first 4 digits
			"is_active": true,
			"page_size": 5,
		})
		return false, "HS Code not found", suggestions, nil
	}
	
	if !hsCode.IsActive {
		return false, "HS Code is inactive", nil, nil
	}
	
	// Country-specific validation could be added here
	// For now, just return valid
	return true, "", nil, nil
}

func (s *tariffService) GetCommonHSCodes(category string) ([]models.HSCode, error) {
	params := map[string]interface{}{
		"is_active": true,
		"page_size": 20,
	}
	
	if category != "" {
		params["category"] = category
	}
	
	codes, _, err := s.repo.FindHSCodes(params)
	return codes, err
}

func (s *tariffService) CalculateTariff(req TariffCalculationRequest) (*TariffCalculationResult, error) {
	// Validate input
	if req.HSCode == "" || req.FromCountry == "" || req.ToCountry == "" {
		return nil, errors.New("HS code, from country, and to country are required")
	}
	
	if req.ProductValue <= 0 {
		return nil, errors.New("product value must be greater than 0")
	}
	
	// Get effective tariff rate
	now := time.Now()
	rate, err := s.repo.GetEffectiveTariffRate(req.HSCode, req.FromCountry, req.ToCountry, now)
	if err != nil {
		// If no specific rate found, try to find general rate or return zero tariff
		return s.calculateWithZeroTariff(req, "No tariff rate found for this route")
	}
	
	// Calculate tariff based on rate type
	result := &TariffCalculationResult{
		HSCode:       req.HSCode,
		FromCountry:  req.FromCountry,
		ToCountry:    req.ToCountry,
		ProductValue: req.ProductValue,
		TariffRate:   rate,
		Currency:     req.Currency,
		Warnings:     []string{},
	}
	
	details := CalculationDetails{
		BaseValue: req.ProductValue,
	}
	
	// Check for preferential treatment
	effectiveRate := rate.RateValue
	if req.PreferentialTreatment && rate.PreferentialRate > 0 {
		originalRate := effectiveRate
		effectiveRate = rate.PreferentialRate
		details.PreferentialApplied = true
		
		// Add conditions warning if any
		if rate.PreferentialConditions != "" {
			result.Warnings = append(result.Warnings, 
				fmt.Sprintf("Preferential rate conditions: %s", rate.PreferentialConditions))
		}
	}
	
	// Calculate based on rate type
	switch rate.RateType {
	case "ad_valorem":
		// Percentage of value
		details.AdValoremDuty = req.ProductValue * (effectiveRate / 100)
		details.TotalDuty = details.AdValoremDuty
		
	case "specific":
		// Fixed amount per unit
		if req.Quantity <= 0 {
			return nil, errors.New("quantity required for specific duty calculation")
		}
		details.SpecificDuty = rate.SpecificRate * req.Quantity
		details.TotalDuty = details.SpecificDuty
		
	case "compound":
		// Both ad valorem and specific
		details.AdValoremDuty = req.ProductValue * (effectiveRate / 100)
		if req.Quantity > 0 && rate.SpecificRate > 0 {
			details.SpecificDuty = rate.SpecificRate * req.Quantity
		}
		details.TotalDuty = details.AdValoremDuty + details.SpecificDuty
		
	default:
		return nil, fmt.Errorf("unknown rate type: %s", rate.RateType)
	}
	
	// Calculate preferential savings if applied
	if details.PreferentialApplied && rate.RateValue != rate.PreferentialRate {
		normalDuty := req.ProductValue * (rate.RateValue / 100)
		details.PreferentialSavings = normalDuty - details.TotalDuty
		if details.PreferentialSavings < 0 {
			details.PreferentialSavings = 0
		}
	}
	
	result.CalculatedTariff = details.TotalDuty
	result.EffectiveRate = details.TotalDuty / req.ProductValue
	result.CalculationDetails = details
	
	// Add warnings based on incoterm
	if req.Incoterm != "" {
		switch req.Incoterm {
		case "DDP":
			result.Warnings = append(result.Warnings, "DDP: Seller bears all tariff costs")
		case "DDU", "DAP":
			result.Warnings = append(result.Warnings, "Buyer responsible for import duties")
		}
	}
	
	// Save calculation record
	calcRecord := &models.TariffCalculation{
		CompanyID:             req.CompanyID,
		UserID:                req.UserID,
		HSCode:                req.HSCode,
		FromCountry:           req.FromCountry,
		ToCountry:             req.ToCountry,
		ProductValue:          req.ProductValue,
		Quantity:              req.Quantity,
		Unit:                  req.Unit,
		WeightKG:              req.WeightKG,
		Currency:              req.Currency,
		Incoterm:              req.Incoterm,
		PreferentialTreatment: req.PreferentialTreatment,
		TariffRateID:          &rate.ID,
		CalculatedTariff:      result.CalculatedTariff,
		EffectiveRate:         result.EffectiveRate,
		CalculationDetails:    details,
		Warnings:              result.Warnings,
	}
	
	if err := s.repo.CreateCalculation(calcRecord); err != nil {
		// Log error but don't fail the calculation
		result.Warnings = append(result.Warnings, "Failed to save calculation history")
	}
	
	return result, nil
}

func (s *tariffService) calculateWithZeroTariff(req TariffCalculationRequest, warning string) (*TariffCalculationResult, error) {
	result := &TariffCalculationResult{
		HSCode:           req.HSCode,
		FromCountry:      req.FromCountry,
		ToCountry:        req.ToCountry,
		ProductValue:     req.ProductValue,
		CalculatedTariff: 0,
		EffectiveRate:    0,
		Currency:         req.Currency,
		CalculationDetails: CalculationDetails{
			BaseValue: req.ProductValue,
			TotalDuty: 0,
		},
		Warnings: []string{warning},
	}
	
	return result, nil
}

func (s *tariffService) BatchCalculateTariff(items []TariffCalculationRequest) ([]TariffCalculationResult, error) {
	results := make([]TariffCalculationResult, 0, len(items))
	
	for _, item := range items {
		result, err := s.CalculateTariff(item)
		if err != nil {
			// Add error as warning and continue
			result = &TariffCalculationResult{
				HSCode:       item.HSCode,
				FromCountry:  item.FromCountry,
				ToCountry:    item.ToCountry,
				ProductValue: item.ProductValue,
				Currency:     item.Currency,
				Warnings:     []string{fmt.Sprintf("Calculation error: %v", err)},
			}
		}
		results = append(results, *result)
	}
	
	return results, nil
}

func (s *tariffService) GetTradeAgreements(countries []string) ([]models.TradeAgreement, error) {
	return s.repo.FindTradeAgreements(countries)
}

func (s *tariffService) GetCalculationHistory(companyID uuid.UUID, limit int) ([]models.TariffCalculation, error) {
	return s.repo.GetCalculationHistory(companyID, limit)
}