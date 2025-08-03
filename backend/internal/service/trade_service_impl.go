package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/fastenmind/fastener-api/internal/models"
	"github.com/fastenmind/fastener-api/internal/repository"
)

var (
	// ErrTariffCodeNotFound is returned when a tariff code is not found
	ErrTariffCodeNotFound = errors.New("tariff code not found")
	// ErrShipmentNotFound is returned when a shipment is not found
	ErrShipmentNotFound = errors.New("shipment not found")
	// ErrDuplicateHSCode is returned when a duplicate HS code is detected
	ErrDuplicateHSCode = errors.New("HS code already exists")
)

// TradeServiceImpl implements TradeService interface
type TradeServiceImpl struct {
	tradeRepo repository.TradeRepository
}

// NewTradeServiceImpl creates a new trade service implementation
func NewTradeServiceImpl(tradeRepo repository.TradeRepository) TradeService {
	return &TradeServiceImpl{
		tradeRepo: tradeRepo,
	}
}

// Tariff Code Management

func (s *TradeServiceImpl) CreateTariffCode(ctx context.Context, code *models.TariffCode) error {
	// Check if HS code already exists
	existing, _ := s.tradeRepo.GetTariffCodeByHSCode(ctx, code.CompanyID, code.HSCode)
	if existing != nil {
		return ErrDuplicateHSCode
	}

	code.ID = uuid.New()
	code.CreatedAt = time.Now()
	code.UpdatedAt = time.Now()

	return s.tradeRepo.CreateTariffCode(ctx, code)
}

func (s *TradeServiceImpl) GetTariffCode(ctx context.Context, id uuid.UUID) (*models.TariffCode, error) {
	code, err := s.tradeRepo.GetTariffCode(ctx, id)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, ErrTariffCodeNotFound
		}
		return nil, err
	}
	return code, nil
}

func (s *TradeServiceImpl) ListTariffCodes(ctx context.Context, params map[string]interface{}) ([]*models.TariffCode, int64, error) {
	return s.tradeRepo.ListTariffCodes(ctx, params)
}

func (s *TradeServiceImpl) UpdateTariffCode(ctx context.Context, code *models.TariffCode) error {
	// Check if tariff code exists
	existing, err := s.tradeRepo.GetTariffCode(ctx, code.ID)
	if err != nil {
		if err == repository.ErrNotFound {
			return ErrTariffCodeNotFound
		}
		return err
	}

	// Check if HS code changed and already exists
	if existing.HSCode != code.HSCode {
		duplicate, _ := s.tradeRepo.GetTariffCodeByHSCode(ctx, code.CompanyID, code.HSCode)
		if duplicate != nil && duplicate.ID != code.ID {
			return ErrDuplicateHSCode
		}
	}

	code.UpdatedAt = time.Now()
	return s.tradeRepo.UpdateTariffCode(ctx, code)
}

func (s *TradeServiceImpl) DeleteTariffCode(ctx context.Context, id uuid.UUID) error {
	_, err := s.tradeRepo.GetTariffCode(ctx, id)
	if err != nil {
		if err == repository.ErrNotFound {
			return ErrTariffCodeNotFound
		}
		return err
	}

	return s.tradeRepo.DeleteTariffCode(ctx, id)
}

// Tariff Rate Management

func (s *TradeServiceImpl) CreateTariffRate(ctx context.Context, rate *models.TariffRate) error {
	rate.ID = uuid.New()
	rate.CreatedAt = time.Now()
	rate.UpdatedAt = time.Now()

	return s.tradeRepo.CreateTariffRate(ctx, rate)
}

func (s *TradeServiceImpl) GetTariffRate(ctx context.Context, id uuid.UUID) (*models.TariffRate, error) {
	return s.tradeRepo.GetTariffRate(ctx, id)
}

func (s *TradeServiceImpl) ListTariffRates(ctx context.Context, params map[string]interface{}) ([]*models.TariffRate, error) {
	return s.tradeRepo.ListTariffRates(ctx, params)
}

func (s *TradeServiceImpl) UpdateTariffRate(ctx context.Context, rate *models.TariffRate) error {
	rate.UpdatedAt = time.Now()
	return s.tradeRepo.UpdateTariffRate(ctx, rate)
}

func (s *TradeServiceImpl) GetEffectiveTariffRate(ctx context.Context, tariffCodeID uuid.UUID, originCountry, destCountry string, date time.Time) (*models.TariffRate, error) {
	return s.tradeRepo.GetEffectiveTariffRate(ctx, tariffCodeID, originCountry, destCountry, date)
}

// Shipment Management

func (s *TradeServiceImpl) CreateShipment(ctx context.Context, shipment *models.Shipment) error {
	shipment.ID = uuid.New()
	shipment.ShipmentNo = s.generateShipmentNumber(shipment.CompanyID)
	shipment.CreatedAt = time.Now()
	shipment.UpdatedAt = time.Now()

	return s.tradeRepo.CreateShipment(ctx, shipment)
}

func (s *TradeServiceImpl) GetShipment(ctx context.Context, id uuid.UUID) (*models.Shipment, error) {
	shipment, err := s.tradeRepo.GetShipment(ctx, id)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, ErrShipmentNotFound
		}
		return nil, err
	}
	return shipment, nil
}

func (s *TradeServiceImpl) ListShipments(ctx context.Context, params map[string]interface{}) ([]*models.Shipment, int64, error) {
	return s.tradeRepo.ListShipments(ctx, params)
}

func (s *TradeServiceImpl) UpdateShipment(ctx context.Context, shipment *models.Shipment) error {
	// Check if shipment exists
	_, err := s.tradeRepo.GetShipment(ctx, shipment.ID)
	if err != nil {
		if err == repository.ErrNotFound {
			return ErrShipmentNotFound
		}
		return err
	}

	shipment.UpdatedAt = time.Now()
	return s.tradeRepo.UpdateShipment(ctx, shipment)
}

func (s *TradeServiceImpl) CreateShipmentEvent(ctx context.Context, event *models.ShipmentEvent) error {
	event.ID = uuid.New()
	event.RecordedAt = time.Now()
	event.CreatedAt = time.Now()

	return s.tradeRepo.CreateShipmentEvent(ctx, event)
}

func (s *TradeServiceImpl) GetShipmentEvents(ctx context.Context, shipmentID uuid.UUID) ([]*models.ShipmentEvent, error) {
	return s.tradeRepo.GetShipmentEvents(ctx, shipmentID)
}

// Document Management

func (s *TradeServiceImpl) CreateTradeDocument(ctx context.Context, doc *models.TradeDocument) error {
	doc.ID = uuid.New()
	doc.CreatedAt = time.Now()
	doc.UpdatedAt = time.Now()

	return s.tradeRepo.CreateTradeDocument(ctx, doc)
}

func (s *TradeServiceImpl) GetTradeDocument(ctx context.Context, id uuid.UUID) (*models.TradeDocument, error) {
	return s.tradeRepo.GetTradeDocument(ctx, id)
}

func (s *TradeServiceImpl) ListTradeDocuments(ctx context.Context, params map[string]interface{}) ([]*models.TradeDocument, error) {
	return s.tradeRepo.ListTradeDocuments(ctx, params)
}

func (s *TradeServiceImpl) UpdateTradeDocument(ctx context.Context, doc *models.TradeDocument) error {
	doc.UpdatedAt = time.Now()
	return s.tradeRepo.UpdateTradeDocument(ctx, doc)
}

// Letter of Credit Management

func (s *TradeServiceImpl) CreateLetterOfCredit(ctx context.Context, lc *models.LetterOfCredit) error {
	lc.ID = uuid.New()
	lc.CreatedAt = time.Now()
	lc.UpdatedAt = time.Now()

	return s.tradeRepo.CreateLetterOfCredit(ctx, lc)
}

func (s *TradeServiceImpl) GetLetterOfCredit(ctx context.Context, id uuid.UUID) (*models.LetterOfCredit, error) {
	return s.tradeRepo.GetLetterOfCredit(ctx, id)
}

func (s *TradeServiceImpl) ListLettersOfCredit(ctx context.Context, params map[string]interface{}) ([]*models.LetterOfCredit, int64, error) {
	return s.tradeRepo.ListLettersOfCredit(ctx, params)
}

func (s *TradeServiceImpl) UpdateLetterOfCredit(ctx context.Context, lc *models.LetterOfCredit) error {
	lc.UpdatedAt = time.Now()
	return s.tradeRepo.UpdateLetterOfCredit(ctx, lc)
}

// Exchange Rate Management

func (s *TradeServiceImpl) GetLatestExchangeRate(ctx context.Context, fromCurrency, toCurrency string) (*models.ExchangeRate, error) {
	return s.tradeRepo.GetLatestExchangeRate(ctx, fromCurrency, toCurrency)
}

func (s *TradeServiceImpl) CreateExchangeRate(ctx context.Context, rate *models.ExchangeRate) error {
	rate.ID = uuid.New().String()
	rate.CreatedAt = time.Now()
	rate.UpdatedAt = time.Now()

	return s.tradeRepo.CreateExchangeRate(ctx, rate)
}

// Compliance Check

func (s *TradeServiceImpl) CreateComplianceCheck(ctx context.Context, check *models.TradeComplianceCheck) error {
	check.ID = uuid.New()
	check.CheckedAt = time.Now()
	check.CreatedAt = time.Now()

	return s.tradeRepo.CreateComplianceCheck(ctx, check)
}

func (s *TradeServiceImpl) GetComplianceChecks(ctx context.Context, resourceType, resourceID string) ([]*models.TradeComplianceCheck, error) {
	return s.tradeRepo.GetComplianceChecks(ctx, resourceType, resourceID)
}

// Helper methods

func (s *TradeServiceImpl) generateShipmentNumber(companyID uuid.UUID) string {
	// Generate shipment number based on date and sequence
	return fmt.Sprintf("SHP-%s-%06d", time.Now().Format("20060102"), time.Now().Unix()%1000000)
}

// Tariff Calculation

func (s *TradeServiceImpl) CalculateTariff(ctx context.Context, params models.TariffCalculationParams) (*models.TariffCalculationResult, error) {
	// Get tariff code
	tariffCode, err := s.tradeRepo.GetTariffCodeByHSCode(ctx, params.CompanyID, params.HSCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get tariff code: %w", err)
	}

	// Get effective tariff rate
	rate, err := s.tradeRepo.GetEffectiveTariffRate(ctx, tariffCode.ID, params.OriginCountry, params.DestCountry, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to get tariff rate: %w", err)
	}

	// Calculate duties and taxes
	result := &models.TariffCalculationResult{
		HSCode:        params.HSCode,
		Description:   tariffCode.Description,
		OriginCountry: params.OriginCountry,
		DestCountry:   params.DestCountry,
		Quantity:      params.Quantity,
		Unit:          tariffCode.Unit,
		UnitValue:     params.UnitValue,
		TotalValue:    params.Quantity * params.UnitValue,
		Currency:      params.Currency,
		TariffRate:    rate.Rate,
		RateType:      rate.RateType,
	}

	// Calculate base duty
	if rate.RateType == "ad_valorem" {
		result.BaseDuty = result.TotalValue * (rate.Rate / 100)
	} else if rate.RateType == "specific" {
		result.BaseDuty = params.Quantity * rate.Rate
	}

	// Apply minimum/maximum duty
	if rate.MinimumDuty > 0 && result.BaseDuty < rate.MinimumDuty {
		result.BaseDuty = rate.MinimumDuty
	}
	if rate.MaximumDuty > 0 && result.BaseDuty > rate.MaximumDuty {
		result.BaseDuty = rate.MaximumDuty
	}

	// Calculate VAT
	result.VATRate = tariffCode.VAT
	result.VAT = (result.TotalValue + result.BaseDuty) * (tariffCode.VAT / 100)

	// Calculate excise tax
	result.ExciseTaxRate = tariffCode.ExciseTax
	result.ExciseTax = result.TotalValue * (tariffCode.ExciseTax / 100)

	// Total taxes
	result.TotalTax = result.BaseDuty + result.VAT + result.ExciseTax
	result.TotalAmount = result.TotalValue + result.TotalTax

	// Trade agreement
	result.TradeAgreement = rate.AgreementType

	return result, nil
}