package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/fastenmind/fastener-api/internal/models"
	"gorm.io/gorm"
)

type tradeRepositoryImpl struct {
	db *gorm.DB
}

// NewTradeRepositoryImpl creates a new trade repository implementation
func NewTradeRepositoryImpl(db *gorm.DB) TradeRepository {
	return &tradeRepositoryImpl{db: db}
}

// Tariff Code Management

func (r *tradeRepositoryImpl) CreateTariffCode(ctx context.Context, code *models.TariffCode) error {
	return r.db.WithContext(ctx).Create(code).Error
}

func (r *tradeRepositoryImpl) GetTariffCode(ctx context.Context, id uuid.UUID) (*models.TariffCode, error) {
	var code models.TariffCode
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&code).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &code, nil
}

func (r *tradeRepositoryImpl) GetTariffCodeByHSCode(ctx context.Context, companyID uuid.UUID, hsCode string) (*models.TariffCode, error) {
	var code models.TariffCode
	err := r.db.WithContext(ctx).Where("company_id = ? AND hs_code = ?", companyID, hsCode).First(&code).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &code, nil
}

func (r *tradeRepositoryImpl) ListTariffCodes(ctx context.Context, params map[string]interface{}) ([]*models.TariffCode, int64, error) {
	var codes []*models.TariffCode
	var total int64

	query := r.db.WithContext(ctx).Model(&models.TariffCode{})

	// Apply filters
	if companyID, ok := params["company_id"].(uuid.UUID); ok {
		query = query.Where("company_id = ?", companyID)
	}
	if search, ok := params["search"].(string); ok && search != "" {
		query = query.Where("hs_code LIKE ? OR description LIKE ? OR description_en LIKE ?", 
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if category, ok := params["category"].(string); ok && category != "" {
		query = query.Where("category = ?", category)
	}
	if isActive, ok := params["is_active"].(bool); ok {
		query = query.Where("is_active = ?", isActive)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if page, ok := params["page"].(int); ok && page > 0 {
		if pageSize, ok := params["page_size"].(int); ok && pageSize > 0 {
			offset := (page - 1) * pageSize
			query = query.Offset(offset).Limit(pageSize)
		}
	}

	// Fetch data
	if err := query.Order("hs_code ASC").Find(&codes).Error; err != nil {
		return nil, 0, err
	}

	return codes, total, nil
}

func (r *tradeRepositoryImpl) UpdateTariffCode(ctx context.Context, code *models.TariffCode) error {
	return r.db.WithContext(ctx).Save(code).Error
}

func (r *tradeRepositoryImpl) DeleteTariffCode(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.TariffCode{}, "id = ?", id).Error
}

// Tariff Rate Management

func (r *tradeRepositoryImpl) CreateTariffRate(ctx context.Context, rate *models.TariffRate) error {
	return r.db.WithContext(ctx).Create(rate).Error
}

func (r *tradeRepositoryImpl) GetTariffRate(ctx context.Context, id uuid.UUID) (*models.TariffRate, error) {
	var rate models.TariffRate
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&rate).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &rate, nil
}

func (r *tradeRepositoryImpl) ListTariffRates(ctx context.Context, params map[string]interface{}) ([]*models.TariffRate, error) {
	var rates []*models.TariffRate

	query := r.db.WithContext(ctx).Model(&models.TariffRate{})

	// Apply filters
	if tariffCodeID, ok := params["tariff_code_id"].(uuid.UUID); ok {
		query = query.Where("tariff_code_id = ?", tariffCodeID)
	}
	if countryCode, ok := params["country_code"].(string); ok && countryCode != "" {
		query = query.Where("country_code = ?", countryCode)
	}
	if tradeType, ok := params["trade_type"].(string); ok && tradeType != "" {
		query = query.Where("trade_type = ?", tradeType)
	}
	if isActive, ok := params["is_active"].(bool); ok {
		query = query.Where("is_active = ?", isActive)
	}

	// Fetch data
	if err := query.Order("country_name ASC").Find(&rates).Error; err != nil {
		return nil, err
	}

	return rates, nil
}

func (r *tradeRepositoryImpl) UpdateTariffRate(ctx context.Context, rate *models.TariffRate) error {
	return r.db.WithContext(ctx).Save(rate).Error
}

func (r *tradeRepositoryImpl) GetEffectiveTariffRate(ctx context.Context, tariffCodeID uuid.UUID, originCountry, destCountry string, date time.Time) (*models.TariffRate, error) {
	var rate models.TariffRate
	err := r.db.WithContext(ctx).
		Where("tariff_code_id = ? AND country_code = ? AND trade_type = ? AND is_active = ? AND valid_from <= ? AND (valid_to IS NULL OR valid_to >= ?)", 
			tariffCodeID, destCountry, "import", true, date, date).
		Order("agreement_type DESC, valid_from DESC"). // Prefer FTA/preferential rates
		First(&rate).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &rate, nil
}

// Shipment Management

func (r *tradeRepositoryImpl) CreateShipment(ctx context.Context, shipment *models.Shipment) error {
	return r.db.WithContext(ctx).Create(shipment).Error
}

func (r *tradeRepositoryImpl) GetShipment(ctx context.Context, id uuid.UUID) (*models.Shipment, error) {
	var shipment models.Shipment
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Documents").
		Preload("Events").
		Where("id = ?", id).
		First(&shipment).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &shipment, nil
}

func (r *tradeRepositoryImpl) ListShipments(ctx context.Context, params map[string]interface{}) ([]*models.Shipment, int64, error) {
	var shipments []*models.Shipment
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Shipment{})

	// Apply filters
	if companyID, ok := params["company_id"].(uuid.UUID); ok {
		query = query.Where("company_id = ?", companyID)
	}
	if status, ok := params["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if transportMode, ok := params["transport_mode"].(string); ok && transportMode != "" {
		query = query.Where("transport_mode = ?", transportMode)
	}
	if search, ok := params["search"].(string); ok && search != "" {
		query = query.Where("shipment_no LIKE ? OR carrier_tracking_no LIKE ?", 
			"%"+search+"%", "%"+search+"%")
	}
	if fromDate, ok := params["from_date"].(time.Time); ok {
		query = query.Where("created_at >= ?", fromDate)
	}
	if toDate, ok := params["to_date"].(time.Time); ok {
		query = query.Where("created_at <= ?", toDate)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if page, ok := params["page"].(int); ok && page > 0 {
		if pageSize, ok := params["page_size"].(int); ok && pageSize > 0 {
			offset := (page - 1) * pageSize
			query = query.Offset(offset).Limit(pageSize)
		}
	}

	// Fetch data with preloads
	if err := query.
		Preload("Order").
		Order("created_at DESC").
		Find(&shipments).Error; err != nil {
		return nil, 0, err
	}

	return shipments, total, nil
}

func (r *tradeRepositoryImpl) UpdateShipment(ctx context.Context, shipment *models.Shipment) error {
	return r.db.WithContext(ctx).Save(shipment).Error
}

func (r *tradeRepositoryImpl) CreateShipmentEvent(ctx context.Context, event *models.ShipmentEvent) error {
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *tradeRepositoryImpl) GetShipmentEvents(ctx context.Context, shipmentID uuid.UUID) ([]*models.ShipmentEvent, error) {
	var events []*models.ShipmentEvent
	err := r.db.WithContext(ctx).
		Where("shipment_id = ?", shipmentID).
		Order("event_time DESC").
		Find(&events).Error
	return events, err
}

// Document Management

func (r *tradeRepositoryImpl) CreateTradeDocument(ctx context.Context, doc *models.TradeDocument) error {
	return r.db.WithContext(ctx).Create(doc).Error
}

func (r *tradeRepositoryImpl) GetTradeDocument(ctx context.Context, id uuid.UUID) (*models.TradeDocument, error) {
	var doc models.TradeDocument
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&doc).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &doc, nil
}

func (r *tradeRepositoryImpl) ListTradeDocuments(ctx context.Context, params map[string]interface{}) ([]*models.TradeDocument, error) {
	var docs []*models.TradeDocument

	query := r.db.WithContext(ctx).Model(&models.TradeDocument{})

	// Apply filters
	if resourceType, ok := params["resource_type"].(string); ok && resourceType != "" {
		query = query.Where("resource_type = ?", resourceType)
	}
	if resourceID, ok := params["resource_id"].(string); ok && resourceID != "" {
		query = query.Where("resource_id = ?", resourceID)
	}
	if docType, ok := params["document_type"].(string); ok && docType != "" {
		query = query.Where("document_type = ?", docType)
	}
	if status, ok := params["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}

	// Fetch data
	if err := query.Order("created_at DESC").Find(&docs).Error; err != nil {
		return nil, err
	}

	return docs, nil
}

func (r *tradeRepositoryImpl) UpdateTradeDocument(ctx context.Context, doc *models.TradeDocument) error {
	return r.db.WithContext(ctx).Save(doc).Error
}

// Letter of Credit Management

func (r *tradeRepositoryImpl) CreateLetterOfCredit(ctx context.Context, lc *models.LetterOfCredit) error {
	return r.db.WithContext(ctx).Create(lc).Error
}

func (r *tradeRepositoryImpl) GetLetterOfCredit(ctx context.Context, id uuid.UUID) (*models.LetterOfCredit, error) {
	var lc models.LetterOfCredit
	err := r.db.WithContext(ctx).
		Preload("Utilizations").
		Where("id = ?", id).
		First(&lc).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &lc, nil
}

func (r *tradeRepositoryImpl) ListLettersOfCredit(ctx context.Context, params map[string]interface{}) ([]*models.LetterOfCredit, int64, error) {
	var lcs []*models.LetterOfCredit
	var total int64

	query := r.db.WithContext(ctx).Model(&models.LetterOfCredit{})

	// Apply filters
	if companyID, ok := params["company_id"].(uuid.UUID); ok {
		query = query.Where("company_id = ?", companyID)
	}
	if status, ok := params["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if lcType, ok := params["type"].(string); ok && lcType != "" {
		query = query.Where("type = ?", lcType)
	}
	if search, ok := params["search"].(string); ok && search != "" {
		query = query.Where("lc_number LIKE ? OR applicant_name LIKE ? OR beneficiary_name LIKE ?", 
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if page, ok := params["page"].(int); ok && page > 0 {
		if pageSize, ok := params["page_size"].(int); ok && pageSize > 0 {
			offset := (page - 1) * pageSize
			query = query.Offset(offset).Limit(pageSize)
		}
	}

	// Fetch data
	if err := query.Order("created_at DESC").Find(&lcs).Error; err != nil {
		return nil, 0, err
	}

	return lcs, total, nil
}

func (r *tradeRepositoryImpl) UpdateLetterOfCredit(ctx context.Context, lc *models.LetterOfCredit) error {
	return r.db.WithContext(ctx).Save(lc).Error
}

// Exchange Rate Management

func (r *tradeRepositoryImpl) GetLatestExchangeRate(ctx context.Context, fromCurrency, toCurrency string) (*models.ExchangeRate, error) {
	var rate models.ExchangeRate
	err := r.db.WithContext(ctx).
		Where("from_currency = ? AND to_currency = ?", fromCurrency, toCurrency).
		Order("effective_date DESC").
		First(&rate).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &rate, nil
}

func (r *tradeRepositoryImpl) CreateExchangeRate(ctx context.Context, rate *models.ExchangeRate) error {
	return r.db.WithContext(ctx).Create(rate).Error
}

// Compliance Check

func (r *tradeRepositoryImpl) CreateComplianceCheck(ctx context.Context, check *models.TradeComplianceCheck) error {
	return r.db.WithContext(ctx).Create(check).Error
}

func (r *tradeRepositoryImpl) GetComplianceChecks(ctx context.Context, resourceType, resourceID string) ([]*models.TradeComplianceCheck, error) {
	var checks []*models.TradeComplianceCheck
	err := r.db.WithContext(ctx).
		Where("resource_type = ? AND resource_id = ?", resourceType, resourceID).
		Order("checked_at DESC").
		Find(&checks).Error
	return checks, err
}