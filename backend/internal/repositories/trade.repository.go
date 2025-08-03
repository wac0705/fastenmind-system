package repositories

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/fastenmind/fastener-api/internal/models"
)

type TradeRepository struct {
	db *gorm.DB
}

func NewTradeRepository(db *gorm.DB) *TradeRepository {
	return &TradeRepository{db: db}
}

// TariffCode methods
func (r *TradeRepository) CreateTariffCode(tariffCode *models.TariffCode) error {
	return r.db.Create(tariffCode).Error
}

func (r *TradeRepository) GetTariffCode(id uuid.UUID) (*models.TariffCode, error) {
	var tariffCode models.TariffCode
	err := r.db.Preload("Company").Preload("Creator").
		First(&tariffCode, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &tariffCode, nil
}

func (r *TradeRepository) GetTariffCodesByCompany(companyID uuid.UUID, hsCode, category string, isActive *bool) ([]models.TariffCode, error) {
	var tariffCodes []models.TariffCode
	query := r.db.Where("company_id = ?", companyID)

	if hsCode != "" {
		// 使用安全的 LIKE 查詢，轉義特殊字符
		escapedHSCode := strings.ReplaceAll(hsCode, "%", "\\%")
		escapedHSCode = strings.ReplaceAll(escapedHSCode, "_", "\\_")
		query = query.Where("hs_code ILIKE ?", "%"+escapedHSCode+"%")
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	err := query.Preload("Company").Preload("Creator").
		Order("hs_code ASC").Find(&tariffCodes).Error
	return tariffCodes, err
}

func (r *TradeRepository) GetTariffCodeByHSCode(companyID uuid.UUID, hsCode string) (*models.TariffCode, error) {
	var tariffCode models.TariffCode
	err := r.db.Where("company_id = ? AND hs_code = ? AND is_active = ?", companyID, hsCode, true).
		First(&tariffCode).Error
	if err != nil {
		return nil, err
	}
	return &tariffCode, nil
}

func (r *TradeRepository) UpdateTariffCode(tariffCode *models.TariffCode) error {
	return r.db.Save(tariffCode).Error
}

func (r *TradeRepository) DeleteTariffCode(id uuid.UUID) error {
	return r.db.Delete(&models.TariffCode{}, id).Error
}

// TariffRate methods
func (r *TradeRepository) CreateTariffRate(tariffRate *models.TariffRate) error {
	return r.db.Create(tariffRate).Error
}

func (r *TradeRepository) GetTariffRate(id uuid.UUID) (*models.TariffRate, error) {
	var tariffRate models.TariffRate
	err := r.db.Preload("Company").Preload("TariffCode").Preload("Creator").
		First(&tariffRate, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &tariffRate, nil
}

func (r *TradeRepository) GetTariffRatesByTariffCode(tariffCodeID uuid.UUID, countryCode, tradeType, agreementType string) ([]models.TariffRate, error) {
	var tariffRates []models.TariffRate
	query := r.db.Where("tariff_code_id = ? AND is_active = ?", tariffCodeID, true)

	if countryCode != "" {
		query = query.Where("country_code = ?", countryCode)
	}
	if tradeType != "" {
		query = query.Where("trade_type = ?", tradeType)
	}
	if agreementType != "" {
		query = query.Where("agreement_type = ?", agreementType)
	}

	// Filter by validity period
	now := time.Now()
	query = query.Where("valid_from <= ? AND (valid_to IS NULL OR valid_to >= ?)", now, now)

	err := query.Preload("Company").Preload("TariffCode").Preload("Creator").
		Order("rate ASC").Find(&tariffRates).Error
	return tariffRates, err
}

func (r *TradeRepository) GetTariffRatesByCompany(companyID uuid.UUID, countryCode, tradeType string) ([]models.TariffRate, error) {
	var tariffRates []models.TariffRate
	query := r.db.Where("company_id = ? AND is_active = ?", companyID, true)

	if countryCode != "" {
		query = query.Where("country_code = ?", countryCode)
	}
	if tradeType != "" {
		query = query.Where("trade_type = ?", tradeType)
	}

	err := query.Preload("Company").Preload("TariffCode").Preload("Creator").
		Order("country_code ASC, rate ASC").Find(&tariffRates).Error
	return tariffRates, err
}

func (r *TradeRepository) UpdateTariffRate(tariffRate *models.TariffRate) error {
	return r.db.Save(tariffRate).Error
}

func (r *TradeRepository) DeleteTariffRate(id uuid.UUID) error {
	return r.db.Delete(&models.TariffRate{}, id).Error
}

// TradeDocument methods
func (r *TradeRepository) CreateTradeDocument(document *models.TradeDocument) error {
	return r.db.Create(document).Error
}

func (r *TradeRepository) GetTradeDocument(id uuid.UUID) (*models.TradeDocument, error) {
	var document models.TradeDocument
	err := r.db.Preload("Company").Preload("Creator").Preload("Approver").
		First(&document, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &document, nil
}

func (r *TradeRepository) GetTradeDocumentsByCompany(companyID uuid.UUID, documentType, status string) ([]models.TradeDocument, error) {
	var documents []models.TradeDocument
	query := r.db.Where("company_id = ?", companyID)

	if documentType != "" {
		query = query.Where("document_type = ?", documentType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Preload("Company").Preload("Creator").Preload("Approver").
		Order("created_at DESC").Find(&documents).Error
	return documents, err
}

func (r *TradeRepository) GetTradeDocumentsByShipment(shipmentID uuid.UUID) ([]models.TradeDocument, error) {
	var documents []models.TradeDocument
	err := r.db.Table("trade_documents").
		Joins("JOIN shipment_documents ON trade_documents.id = shipment_documents.trade_document_id").
		Where("shipment_documents.shipment_id = ?", shipmentID).
		Preload("Company").Preload("Creator").Preload("Approver").
		Find(&documents).Error
	return documents, err
}

func (r *TradeRepository) UpdateTradeDocument(document *models.TradeDocument) error {
	return r.db.Save(document).Error
}

func (r *TradeRepository) DeleteTradeDocument(id uuid.UUID) error {
	return r.db.Delete(&models.TradeDocument{}, id).Error
}

// Shipment methods
func (r *TradeRepository) CreateShipment(shipment *models.Shipment) error {
	return r.db.Create(shipment).Error
}

func (r *TradeRepository) GetShipment(id uuid.UUID) (*models.Shipment, error) {
	var shipment models.Shipment
	err := r.db.Preload("Company").Preload("Order").Preload("Creator").
		Preload("Documents").Preload("Items").Preload("Events").
		First(&shipment, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &shipment, nil
}

func (r *TradeRepository) GetShipmentByNumber(companyID uuid.UUID, shipmentNo string) (*models.Shipment, error) {
	var shipment models.Shipment
	err := r.db.Where("company_id = ? AND shipment_no = ?", companyID, shipmentNo).
		Preload("Company").Preload("Order").Preload("Creator").
		Preload("Documents").Preload("Items").Preload("Events").
		First(&shipment).Error
	if err != nil {
		return nil, err
	}
	return &shipment, nil
}

func (r *TradeRepository) GetShipmentsByCompany(companyID uuid.UUID, shipmentType, status, method string) ([]models.Shipment, error) {
	var shipments []models.Shipment
	query := r.db.Where("company_id = ?", companyID)

	if shipmentType != "" {
		query = query.Where("type = ?", shipmentType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if method != "" {
		query = query.Where("method = ?", method)
	}

	err := query.Preload("Company").Preload("Order").Preload("Creator").
		Order("created_at DESC").Find(&shipments).Error
	return shipments, err
}

func (r *TradeRepository) GetShipmentsByOrder(orderID uuid.UUID) ([]models.Shipment, error) {
	var shipments []models.Shipment
	err := r.db.Where("order_id = ?", orderID).
		Preload("Company").Preload("Order").Preload("Creator").
		Preload("Items").Order("created_at DESC").Find(&shipments).Error
	return shipments, err
}

func (r *TradeRepository) UpdateShipment(shipment *models.Shipment) error {
	return r.db.Save(shipment).Error
}

func (r *TradeRepository) DeleteShipment(id uuid.UUID) error {
	return r.db.Select("Items", "Events").Delete(&models.Shipment{ID: id}).Error
}

// ShipmentItem methods
func (r *TradeRepository) CreateShipmentItem(item *models.ShipmentItem) error {
	item.TotalWeight = item.Quantity * item.UnitWeight
	item.TotalValue = item.Quantity * item.UnitValue
	return r.db.Create(item).Error
}

func (r *TradeRepository) GetShipmentItem(id uuid.UUID) (*models.ShipmentItem, error) {
	var item models.ShipmentItem
	err := r.db.Preload("Company").Preload("Shipment").Preload("Product").
		First(&item, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *TradeRepository) GetShipmentItemsByShipment(shipmentID uuid.UUID) ([]models.ShipmentItem, error) {
	var items []models.ShipmentItem
	err := r.db.Where("shipment_id = ?", shipmentID).
		Preload("Company").Preload("Shipment").Preload("Product").
		Order("created_at ASC").Find(&items).Error
	return items, err
}

func (r *TradeRepository) UpdateShipmentItem(item *models.ShipmentItem) error {
	item.TotalWeight = item.Quantity * item.UnitWeight
	item.TotalValue = item.Quantity * item.UnitValue
	return r.db.Save(item).Error
}

func (r *TradeRepository) DeleteShipmentItem(id uuid.UUID) error {
	return r.db.Delete(&models.ShipmentItem{}, id).Error
}

// ShipmentEvent methods
func (r *TradeRepository) CreateShipmentEvent(event *models.ShipmentEvent) error {
	return r.db.Create(event).Error
}

func (r *TradeRepository) GetShipmentEvent(id uuid.UUID) (*models.ShipmentEvent, error) {
	var event models.ShipmentEvent
	err := r.db.Preload("Company").Preload("Shipment").Preload("Creator").
		First(&event, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *TradeRepository) GetShipmentEventsByShipment(shipmentID uuid.UUID) ([]models.ShipmentEvent, error) {
	var events []models.ShipmentEvent
	err := r.db.Where("shipment_id = ?", shipmentID).
		Preload("Company").Preload("Shipment").Preload("Creator").
		Order("event_time DESC").Find(&events).Error
	return events, err
}

func (r *TradeRepository) GetShipmentEventsByLocation(companyID uuid.UUID, location string, eventType string) ([]models.ShipmentEvent, error) {
	var events []models.ShipmentEvent
	query := r.db.Where("company_id = ?", companyID)

	if location != "" {
		// 使用安全的 LIKE 查詢，轉義特殊字符
		escapedLocation := strings.ReplaceAll(location, "%", "\\%")
		escapedLocation = strings.ReplaceAll(escapedLocation, "_", "\\_")
		query = query.Where("location ILIKE ?", "%"+escapedLocation+"%")
	}
	if eventType != "" {
		query = query.Where("event_type = ?", eventType)
	}

	err := query.Preload("Company").Preload("Shipment").Preload("Creator").
		Order("event_time DESC").Find(&events).Error
	return events, err
}

func (r *TradeRepository) UpdateShipmentEvent(event *models.ShipmentEvent) error {
	return r.db.Save(event).Error
}

func (r *TradeRepository) DeleteShipmentEvent(id uuid.UUID) error {
	return r.db.Delete(&models.ShipmentEvent{}, id).Error
}

// LetterOfCredit methods
func (r *TradeRepository) CreateLetterOfCredit(lc *models.LetterOfCredit) error {
	lc.AvailableAmount = lc.Amount - lc.UtilizedAmount
	return r.db.Create(lc).Error
}

func (r *TradeRepository) GetLetterOfCredit(id uuid.UUID) (*models.LetterOfCredit, error) {
	var lc models.LetterOfCredit
	err := r.db.Preload("Company").Preload("Creator").
		Preload("Shipments").Preload("Utilizations").
		First(&lc, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &lc, nil
}

func (r *TradeRepository) GetLetterOfCreditByNumber(companyID uuid.UUID, lcNumber string) (*models.LetterOfCredit, error) {
	var lc models.LetterOfCredit
	err := r.db.Where("company_id = ? AND lc_number = ?", companyID, lcNumber).
		Preload("Company").Preload("Creator").
		Preload("Utilizations").First(&lc).Error
	if err != nil {
		return nil, err
	}
	return &lc, nil
}

func (r *TradeRepository) GetLetterOfCreditsByCompany(companyID uuid.UUID, lcType, status string) ([]models.LetterOfCredit, error) {
	var lcs []models.LetterOfCredit
	query := r.db.Where("company_id = ?", companyID)

	if lcType != "" {
		query = query.Where("type = ?", lcType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Preload("Company").Preload("Creator").
		Order("issue_date DESC").Find(&lcs).Error
	return lcs, err
}

func (r *TradeRepository) GetExpiringLetterOfCredits(companyID uuid.UUID, days int) ([]models.LetterOfCredit, error) {
	var lcs []models.LetterOfCredit
	expiryDate := time.Now().AddDate(0, 0, days)
	
	err := r.db.Where("company_id = ? AND status IN (?) AND expiry_date <= ?", 
		companyID, []string{"issued", "advised", "confirmed"}, expiryDate).
		Preload("Company").Preload("Creator").
		Order("expiry_date ASC").Find(&lcs).Error
	return lcs, err
}

func (r *TradeRepository) UpdateLetterOfCredit(lc *models.LetterOfCredit) error {
	lc.AvailableAmount = lc.Amount - lc.UtilizedAmount
	return r.db.Save(lc).Error
}

func (r *TradeRepository) DeleteLetterOfCredit(id uuid.UUID) error {
	return r.db.Select("Utilizations").Delete(&models.LetterOfCredit{ID: id}).Error
}

// LCUtilization methods
func (r *TradeRepository) CreateLCUtilization(utilization *models.LCUtilization) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Create utilization record
		if err := tx.Create(utilization).Error; err != nil {
			return err
		}

		// Update LC utilized amount
		var lc models.LetterOfCredit
		if err := tx.First(&lc, utilization.LCID).Error; err != nil {
			return err
		}

		lc.UtilizedAmount += utilization.Amount
		lc.AvailableAmount = lc.Amount - lc.UtilizedAmount

		// Check if LC is fully utilized
		if lc.AvailableAmount <= 0 {
			lc.Status = "utilized"
		}

		return tx.Save(&lc).Error
	})
}

func (r *TradeRepository) GetLCUtilization(id uuid.UUID) (*models.LCUtilization, error) {
	var utilization models.LCUtilization
	err := r.db.Preload("Company").Preload("LC").Preload("Shipment").Preload("Creator").
		First(&utilization, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &utilization, nil
}

func (r *TradeRepository) GetLCUtilizationsByLC(lcID uuid.UUID, status string) ([]models.LCUtilization, error) {
	var utilizations []models.LCUtilization
	query := r.db.Where("lc_id = ?", lcID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Preload("Company").Preload("LC").Preload("Shipment").Preload("Creator").
		Order("utilized_at DESC").Find(&utilizations).Error
	return utilizations, err
}

func (r *TradeRepository) UpdateLCUtilization(utilization *models.LCUtilization) error {
	return r.db.Save(utilization).Error
}

// TradeCompliance methods
func (r *TradeRepository) CreateTradeCompliance(compliance *models.TradeCompliance) error {
	return r.db.Create(compliance).Error
}

func (r *TradeRepository) GetTradeCompliance(id uuid.UUID) (*models.TradeCompliance, error) {
	var compliance models.TradeCompliance
	err := r.db.Preload("Company").Preload("Creator").Preload("Updater").
		Preload("Checks").First(&compliance, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &compliance, nil
}

func (r *TradeRepository) GetTradeCompliancesByCompany(companyID uuid.UUID, complianceType, countryCode, status string) ([]models.TradeCompliance, error) {
	var compliances []models.TradeCompliance
	query := r.db.Where("company_id = ?", companyID)

	if complianceType != "" {
		query = query.Where("compliance_type = ?", complianceType)
	}
	if countryCode != "" {
		query = query.Where("country_code = ?", countryCode)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Filter by validity period
	now := time.Now()
	query = query.Where("valid_from <= ? AND (valid_to IS NULL OR valid_to >= ?)", now, now)

	err := query.Preload("Company").Preload("Creator").Preload("Updater").
		Order("severity DESC, created_at DESC").Find(&compliances).Error
	return compliances, err
}

func (r *TradeRepository) GetActiveTradeCompliances(companyID uuid.UUID) ([]models.TradeCompliance, error) {
	var compliances []models.TradeCompliance
	now := time.Now()
	
	err := r.db.Where("company_id = ? AND status = ? AND valid_from <= ? AND (valid_to IS NULL OR valid_to >= ?)", 
		companyID, "active", now, now).
		Preload("Company").Preload("Creator").
		Order("severity DESC").Find(&compliances).Error
	return compliances, err
}

func (r *TradeRepository) UpdateTradeCompliance(compliance *models.TradeCompliance) error {
	return r.db.Save(compliance).Error
}

func (r *TradeRepository) DeleteTradeCompliance(id uuid.UUID) error {
	return r.db.Select("Checks").Delete(&models.TradeCompliance{ID: id}).Error
}

// ComplianceCheck methods
func (r *TradeRepository) CreateComplianceCheck(check *models.ComplianceCheck) error {
	return r.db.Create(check).Error
}

func (r *TradeRepository) GetComplianceCheck(id uuid.UUID) (*models.ComplianceCheck, error) {
	var check models.ComplianceCheck
	err := r.db.Preload("Company").Preload("Compliance").
		Preload("Checker").Preload("Resolver").
		First(&check, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &check, nil
}

func (r *TradeRepository) GetComplianceChecksByCompliance(complianceID uuid.UUID, result string) ([]models.ComplianceCheck, error) {
	var checks []models.ComplianceCheck
	query := r.db.Where("compliance_id = ?", complianceID)

	if result != "" {
		query = query.Where("result = ?", result)
	}

	err := query.Preload("Company").Preload("Compliance").
		Preload("Checker").Preload("Resolver").
		Order("checked_at DESC").Find(&checks).Error
	return checks, err
}

func (r *TradeRepository) GetComplianceChecksByResource(companyID uuid.UUID, resourceType string, resourceID uuid.UUID) ([]models.ComplianceCheck, error) {
	var checks []models.ComplianceCheck
	err := r.db.Where("company_id = ? AND resource_type = ? AND resource_id = ?", 
		companyID, resourceType, resourceID).
		Preload("Company").Preload("Compliance").
		Preload("Checker").Preload("Resolver").
		Order("checked_at DESC").Find(&checks).Error
	return checks, err
}

func (r *TradeRepository) GetFailedComplianceChecks(companyID uuid.UUID) ([]models.ComplianceCheck, error) {
	var checks []models.ComplianceCheck
	err := r.db.Where("company_id = ? AND result IN (?) AND resolved_at IS NULL", 
		companyID, []string{"failed", "warning"}).
		Preload("Company").Preload("Compliance").
		Preload("Checker").Order("checked_at DESC").Find(&checks).Error
	return checks, err
}

func (r *TradeRepository) UpdateComplianceCheck(check *models.ComplianceCheck) error {
	return r.db.Save(check).Error
}

// ExchangeRate methods
func (r *TradeRepository) CreateExchangeRate(rate *models.ExchangeRate) error {
	return r.db.Create(rate).Error
}

func (r *TradeRepository) GetExchangeRate(id uuid.UUID) (*models.ExchangeRate, error) {
	var rate models.ExchangeRate
	err := r.db.Preload("Company").Preload("Creator").
		First(&rate, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &rate, nil
}

func (r *TradeRepository) GetExchangeRatesByCompany(companyID uuid.UUID, fromCurrency, toCurrency, rateType string) ([]models.ExchangeRate, error) {
	var rates []models.ExchangeRate
	query := r.db.Where("company_id = ? AND is_active = ?", companyID, true)

	if fromCurrency != "" {
		query = query.Where("from_currency = ?", fromCurrency)
	}
	if toCurrency != "" {
		query = query.Where("to_currency = ?", toCurrency)
	}
	if rateType != "" {
		query = query.Where("rate_type = ?", rateType)
	}

	err := query.Preload("Company").Preload("Creator").
		Order("valid_date DESC").Find(&rates).Error
	return rates, err
}

func (r *TradeRepository) GetLatestExchangeRate(companyID uuid.UUID, fromCurrency, toCurrency, rateType string) (*models.ExchangeRate, error) {
	var rate models.ExchangeRate
	err := r.db.Where("company_id = ? AND from_currency = ? AND to_currency = ? AND rate_type = ? AND is_active = ?", 
		companyID, fromCurrency, toCurrency, rateType, true).
		Order("valid_date DESC").First(&rate).Error
	if err != nil {
		return nil, err
	}
	return &rate, nil
}

func (r *TradeRepository) UpdateExchangeRate(rate *models.ExchangeRate) error {
	return r.db.Save(rate).Error
}

func (r *TradeRepository) DeleteExchangeRate(id uuid.UUID) error {
	return r.db.Delete(&models.ExchangeRate{}, id).Error
}

// TradeRegulation methods
func (r *TradeRepository) CreateTradeRegulation(regulation *models.TradeRegulation) error {
	return r.db.Create(regulation).Error
}

func (r *TradeRepository) GetTradeRegulation(id uuid.UUID) (*models.TradeRegulation, error) {
	var regulation models.TradeRegulation
	err := r.db.Preload("Company").Preload("Creator").Preload("Reviewer").
		First(&regulation, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &regulation, nil
}

func (r *TradeRepository) GetTradeRegulationsByCountry(companyID *uuid.UUID, countryCode, regulationType, status string) ([]models.TradeRegulation, error) {
	var regulations []models.TradeRegulation
	query := r.db.Where("country_code = ?", countryCode)

	if companyID != nil {
		query = query.Where("company_id = ? OR company_id IS NULL", *companyID)
	} else {
		query = query.Where("company_id IS NULL")
	}

	if regulationType != "" {
		query = query.Where("regulation_type = ?", regulationType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Filter by effective period
	now := time.Now()
	query = query.Where("effective_date <= ? AND (expiry_date IS NULL OR expiry_date >= ?)", now, now)

	err := query.Preload("Company").Preload("Creator").Preload("Reviewer").
		Order("effective_date DESC").Find(&regulations).Error
	return regulations, err
}

func (r *TradeRepository) GetActiveTradeRegulations(companyID *uuid.UUID) ([]models.TradeRegulation, error) {
	var regulations []models.TradeRegulation
	now := time.Now()
	query := r.db.Where("status = ? AND effective_date <= ? AND (expiry_date IS NULL OR expiry_date >= ?)", 
		"active", now, now)

	if companyID != nil {
		query = query.Where("company_id = ? OR company_id IS NULL", *companyID)
	} else {
		query = query.Where("company_id IS NULL")
	}

	err := query.Preload("Company").Preload("Creator").
		Order("country_code ASC, regulation_type ASC").Find(&regulations).Error
	return regulations, err
}

func (r *TradeRepository) UpdateTradeRegulation(regulation *models.TradeRegulation) error {
	return r.db.Save(regulation).Error
}

func (r *TradeRepository) DeleteTradeRegulation(id uuid.UUID) error {
	return r.db.Delete(&models.TradeRegulation{}, id).Error
}

// TradeAgreement methods
func (r *TradeRepository) CreateTradeAgreement(agreement *models.TradeAgreement) error {
	return r.db.Create(agreement).Error
}

func (r *TradeRepository) GetTradeAgreement(id uuid.UUID) (*models.TradeAgreement, error) {
	var agreement models.TradeAgreement
	err := r.db.Preload("Company").Preload("Creator").
		First(&agreement, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &agreement, nil
}

func (r *TradeRepository) GetTradeAgreementByCode(companyID *uuid.UUID, agreementCode string) (*models.TradeAgreement, error) {
	var agreement models.TradeAgreement
	query := r.db.Where("agreement_code = ?", agreementCode)

	if companyID != nil {
		query = query.Where("company_id = ? OR company_id IS NULL", *companyID)
	} else {
		query = query.Where("company_id IS NULL")
	}

	err := query.Preload("Company").Preload("Creator").First(&agreement).Error
	if err != nil {
		return nil, err
	}
	return &agreement, nil
}

func (r *TradeRepository) GetTradeAgreementsByType(companyID *uuid.UUID, agreementType, status string) ([]models.TradeAgreement, error) {
	var agreements []models.TradeAgreement
	query := r.db.Where("1=1")

	if companyID != nil {
		query = query.Where("company_id = ? OR company_id IS NULL", *companyID)
	} else {
		query = query.Where("company_id IS NULL")
	}

	if agreementType != "" {
		query = query.Where("agreement_type = ?", agreementType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Filter by effective period
	now := time.Now()
	query = query.Where("effective_date <= ? AND (expiry_date IS NULL OR expiry_date >= ?)", now, now)

	err := query.Preload("Company").Preload("Creator").
		Order("effective_date DESC").Find(&agreements).Error
	return agreements, err
}

func (r *TradeRepository) GetActiveTradeAgreements(companyID *uuid.UUID) ([]models.TradeAgreement, error) {
	var agreements []models.TradeAgreement
	now := time.Now()
	query := r.db.Where("status = ? AND effective_date <= ? AND (expiry_date IS NULL OR expiry_date >= ?)", 
		"active", now, now)

	if companyID != nil {
		query = query.Where("company_id = ? OR company_id IS NULL", *companyID)
	} else {
		query = query.Where("company_id IS NULL")
	}

	err := query.Preload("Company").Preload("Creator").
		Order("agreement_type ASC, name ASC").Find(&agreements).Error
	return agreements, err
}

func (r *TradeRepository) UpdateTradeAgreement(agreement *models.TradeAgreement) error {
	return r.db.Save(agreement).Error
}

func (r *TradeRepository) DeleteTradeAgreement(id uuid.UUID) error {
	return r.db.Delete(&models.TradeAgreement{}, id).Error
}

// Analytics and Statistics methods
func (r *TradeRepository) GetTradeStatistics(companyID uuid.UUID, startDate, endDate time.Time) (map[string]interface{}, error) {
	var stats map[string]interface{} = make(map[string]interface{})

	// Shipment statistics
	var shipmentStats struct {
		TotalShipments    int64   `json:"total_shipments"`
		ImportShipments   int64   `json:"import_shipments"`
		ExportShipments   int64   `json:"export_shipments"`
		InTransitShipments int64  `json:"in_transit_shipments"`
		DeliveredShipments int64  `json:"delivered_shipments"`
		TotalValue        float64 `json:"total_value"`
		TotalWeight       float64 `json:"total_weight"`
	}

	r.db.Model(&models.Shipment{}).
		Where("company_id = ? AND created_at BETWEEN ? AND ?", companyID, startDate, endDate).
		Select("COUNT(*) as total_shipments").Scan(&shipmentStats)

	r.db.Model(&models.Shipment{}).
		Where("company_id = ? AND type = ? AND created_at BETWEEN ? AND ?", companyID, "import", startDate, endDate).
		Select("COUNT(*) as import_shipments").Scan(&shipmentStats.ImportShipments)

	r.db.Model(&models.Shipment{}).
		Where("company_id = ? AND type = ? AND created_at BETWEEN ? AND ?", companyID, "export", startDate, endDate).
		Select("COUNT(*) as export_shipments").Scan(&shipmentStats.ExportShipments)

	r.db.Model(&models.Shipment{}).
		Where("company_id = ? AND status = ? AND created_at BETWEEN ? AND ?", companyID, "in_transit", startDate, endDate).
		Select("COUNT(*) as in_transit_shipments").Scan(&shipmentStats.InTransitShipments)

	r.db.Model(&models.Shipment{}).
		Where("company_id = ? AND status = ? AND created_at BETWEEN ? AND ?", companyID, "delivered", startDate, endDate).
		Select("COUNT(*) as delivered_shipments").Scan(&shipmentStats.DeliveredShipments)

	r.db.Model(&models.Shipment{}).
		Where("company_id = ? AND created_at BETWEEN ? AND ?", companyID, startDate, endDate).
		Select("COALESCE(SUM(customs_value), 0) as total_value, COALESCE(SUM(gross_weight), 0) as total_weight").
		Scan(&shipmentStats)

	stats["shipments"] = shipmentStats

	// LC statistics
	var lcStats struct {
		TotalLCs       int64   `json:"total_lcs"`
		ActiveLCs      int64   `json:"active_lcs"`
		UtilizedLCs    int64   `json:"utilized_lcs"`
		ExpiringLCs    int64   `json:"expiring_lcs"`
		TotalAmount    float64 `json:"total_amount"`
		UtilizedAmount float64 `json:"utilized_amount"`
	}

	r.db.Model(&models.LetterOfCredit{}).
		Where("company_id = ? AND created_at BETWEEN ? AND ?", companyID, startDate, endDate).
		Select("COUNT(*) as total_lcs, COALESCE(SUM(amount), 0) as total_amount, COALESCE(SUM(utilized_amount), 0) as utilized_amount").
		Scan(&lcStats)

	r.db.Model(&models.LetterOfCredit{}).
		Where("company_id = ? AND status IN (?)", companyID, []string{"issued", "advised", "confirmed"}).
		Select("COUNT(*) as active_lcs").Scan(&lcStats.ActiveLCs)

	r.db.Model(&models.LetterOfCredit{}).
		Where("company_id = ? AND status = ?", companyID, "utilized").
		Select("COUNT(*) as utilized_lcs").Scan(&lcStats.UtilizedLCs)

	expiryDate := time.Now().AddDate(0, 0, 30) // 30 days from now
	r.db.Model(&models.LetterOfCredit{}).
		Where("company_id = ? AND status IN (?) AND expiry_date <= ?", companyID, []string{"issued", "advised", "confirmed"}, expiryDate).
		Select("COUNT(*) as expiring_lcs").Scan(&lcStats.ExpiringLCs)

	stats["letter_of_credits"] = lcStats

	// Compliance statistics
	var complianceStats struct {
		TotalCompliances  int64 `json:"total_compliances"`
		ActiveCompliances int64 `json:"active_compliances"`
		FailedChecks      int64 `json:"failed_checks"`
		WarningChecks     int64 `json:"warning_checks"`
		PassedChecks      int64 `json:"passed_checks"`
	}

	r.db.Model(&models.TradeCompliance{}).
		Where("company_id = ?", companyID).
		Select("COUNT(*) as total_compliances").Scan(&complianceStats.TotalCompliances)

	r.db.Model(&models.TradeCompliance{}).
		Where("company_id = ? AND status = ?", companyID, "active").
		Select("COUNT(*) as active_compliances").Scan(&complianceStats.ActiveCompliances)

	r.db.Model(&models.ComplianceCheck{}).
		Where("company_id = ? AND result = ? AND created_at BETWEEN ? AND ?", companyID, "failed", startDate, endDate).
		Select("COUNT(*) as failed_checks").Scan(&complianceStats.FailedChecks)

	r.db.Model(&models.ComplianceCheck{}).
		Where("company_id = ? AND result = ? AND created_at BETWEEN ? AND ?", companyID, "warning", startDate, endDate).
		Select("COUNT(*) as warning_checks").Scan(&complianceStats.WarningChecks)

	r.db.Model(&models.ComplianceCheck{}).
		Where("company_id = ? AND result = ? AND created_at BETWEEN ? AND ?", companyID, "passed", startDate, endDate).
		Select("COUNT(*) as passed_checks").Scan(&complianceStats.PassedChecks)

	stats["compliance"] = complianceStats

	return stats, nil
}

func (r *TradeRepository) GetShipmentsByCountry(companyID uuid.UUID, startDate, endDate time.Time) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	rows, err := r.db.Model(&models.Shipment{}).
		Where("company_id = ? AND created_at BETWEEN ? AND ?", companyID, startDate, endDate).
		Select("origin_country, dest_country, COUNT(*) as shipment_count, SUM(customs_value) as total_value").
		Group("origin_country, dest_country").
		Having("COUNT(*) > 0").
		Order("shipment_count DESC").Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var originCountry, destCountry string
		var shipmentCount int
		var totalValue float64

		if err := rows.Scan(&originCountry, &destCountry, &shipmentCount, &totalValue); err != nil {
			continue
		}

		results = append(results, map[string]interface{}{
			"origin_country":  originCountry,
			"dest_country":    destCountry,
			"shipment_count":  shipmentCount,
			"total_value":     totalValue,
		})
	}

	return results, nil
}

func (r *TradeRepository) GetTopTradingPartners(companyID uuid.UUID, startDate, endDate time.Time, limit int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	// Import partners
	importRows, err := r.db.Model(&models.Shipment{}).
		Where("company_id = ? AND type = ? AND created_at BETWEEN ? AND ?", companyID, "import", startDate, endDate).
		Select("origin_country as country, COUNT(*) as shipment_count, SUM(customs_value) as total_value").
		Group("origin_country").
		Order("total_value DESC").
		Limit(limit).Rows()

	if err == nil {
		defer importRows.Close()
		for importRows.Next() {
			var country string
			var shipmentCount int
			var totalValue float64

			if err := importRows.Scan(&country, &shipmentCount, &totalValue); err != nil {
				continue
			}

			results = append(results, map[string]interface{}{
				"country":        country,
				"type":          "import",
				"shipment_count": shipmentCount,
				"total_value":    totalValue,
			})
		}
	}

	// Export partners
	exportRows, err := r.db.Model(&models.Shipment{}).
		Where("company_id = ? AND type = ? AND created_at BETWEEN ? AND ?", companyID, "export", startDate, endDate).
		Select("dest_country as country, COUNT(*) as shipment_count, SUM(customs_value) as total_value").
		Group("dest_country").
		Order("total_value DESC").
		Limit(limit).Rows()

	if err == nil {
		defer exportRows.Close()
		for exportRows.Next() {
			var country string
			var shipmentCount int
			var totalValue float64

			if err := exportRows.Scan(&country, &shipmentCount, &totalValue); err != nil {
				continue
			}

			results = append(results, map[string]interface{}{
				"country":        country,
				"type":          "export",
				"shipment_count": shipmentCount,
				"total_value":    totalValue,
			})
		}
	}

	return results, nil
}