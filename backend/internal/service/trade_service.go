package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/fastenmind/fastener-api/internal/models"
	"github.com/fastenmind/fastener-api/internal/repository"
)

// TradeService handles trade-related business logic
type TradeService interface {
	// Tariff Code Management
	CreateTariffCode(ctx context.Context, code *models.TariffCode) error
	GetTariffCode(ctx context.Context, id uuid.UUID) (*models.TariffCode, error)
	ListTariffCodes(ctx context.Context, params map[string]interface{}) ([]*models.TariffCode, int64, error)
	UpdateTariffCode(ctx context.Context, code *models.TariffCode) error
	DeleteTariffCode(ctx context.Context, id uuid.UUID) error

	// Tariff Rate Management
	CreateTariffRate(ctx context.Context, rate *models.TariffRate) error
	GetTariffRate(ctx context.Context, id uuid.UUID) (*models.TariffRate, error)
	ListTariffRates(ctx context.Context, params map[string]interface{}) ([]*models.TariffRate, error)
	UpdateTariffRate(ctx context.Context, rate *models.TariffRate) error
	GetEffectiveTariffRate(ctx context.Context, tariffCodeID uuid.UUID, originCountry, destCountry string, date time.Time) (*models.TariffRate, error)

	// Shipment Management
	CreateShipment(ctx context.Context, shipment *models.Shipment) error
	GetShipment(ctx context.Context, id uuid.UUID) (*models.Shipment, error)
	ListShipments(ctx context.Context, params map[string]interface{}) ([]*models.Shipment, int64, error)
	UpdateShipment(ctx context.Context, shipment *models.Shipment) error
	CreateShipmentEvent(ctx context.Context, event *models.ShipmentEvent) error
	GetShipmentEvents(ctx context.Context, shipmentID uuid.UUID) ([]*models.ShipmentEvent, error)

	// Document Management
	CreateTradeDocument(ctx context.Context, doc *models.TradeDocument) error
	GetTradeDocument(ctx context.Context, id uuid.UUID) (*models.TradeDocument, error)
	ListTradeDocuments(ctx context.Context, params map[string]interface{}) ([]*models.TradeDocument, error)
	UpdateTradeDocument(ctx context.Context, doc *models.TradeDocument) error

	// Letter of Credit Management
	CreateLetterOfCredit(ctx context.Context, lc *models.LetterOfCredit) error
	GetLetterOfCredit(ctx context.Context, id uuid.UUID) (*models.LetterOfCredit, error)
	ListLettersOfCredit(ctx context.Context, params map[string]interface{}) ([]*models.LetterOfCredit, int64, error)
	UpdateLetterOfCredit(ctx context.Context, lc *models.LetterOfCredit) error

	// Exchange Rate Management
	GetLatestExchangeRate(ctx context.Context, fromCurrency, toCurrency string) (*models.ExchangeRate, error)
	CreateExchangeRate(ctx context.Context, rate *models.ExchangeRate) error

	// Compliance Check
	CreateComplianceCheck(ctx context.Context, check *models.TradeComplianceCheck) error
	GetComplianceChecks(ctx context.Context, resourceType, resourceID string) ([]*models.TradeComplianceCheck, error)

	// Tariff Calculation
	CalculateTariff(ctx context.Context, params models.TariffCalculationParams) (*models.TariffCalculationResult, error)
}

// NewTradeService creates a new trade service
func NewTradeService(tradeRepo repository.TradeRepository) TradeService {
	return NewTradeServiceImpl(tradeRepo)
}