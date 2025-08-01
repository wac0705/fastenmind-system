package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/fastenmind/fastener-api/internal/models"
	"github.com/fastenmind/fastener-api/internal/repositories"
)

type TradeService struct {
	tradeRepo *repositories.TradeRepository
	userRepo  *repositories.UserRepository
}

func NewTradeService(tradeRepo *repositories.TradeRepository, userRepo *repositories.UserRepository) *TradeService {
	return &TradeService{
		tradeRepo: tradeRepo,
		userRepo:  userRepo,
	}
}

// TariffCode Service Methods
type CreateTariffCodeRequest struct {
	HSCode             string  `json:"hs_code" validate:"required"`
	Description        string  `json:"description" validate:"required"`
	DescriptionEN      string  `json:"description_en"`
	Category           string  `json:"category"`
	Unit               string  `json:"unit"`
	BaseRate           float64 `json:"base_rate"`
	PreferentialRate   float64 `json:"preferential_rate"`
	VAT                float64 `json:"vat"`
	ExciseTax          float64 `json:"excise_tax"`
	ImportRestriction  map[string]interface{} `json:"import_restriction"`
	ExportRestriction  map[string]interface{} `json:"export_restriction"`
	RequiredCerts      []string `json:"required_certs"`
}

type UpdateTariffCodeRequest struct {
	Description        *string  `json:"description"`
	DescriptionEN      *string  `json:"description_en"`
	Category           *string  `json:"category"`
	Unit               *string  `json:"unit"`
	BaseRate           *float64 `json:"base_rate"`
	PreferentialRate   *float64 `json:"preferential_rate"`
	VAT                *float64 `json:"vat"`
	ExciseTax          *float64 `json:"excise_tax"`
	ImportRestriction  map[string]interface{} `json:"import_restriction"`
	ExportRestriction  map[string]interface{} `json:"export_restriction"`
	RequiredCerts      []string `json:"required_certs"`
	IsActive           *bool    `json:"is_active"`
}

func (s *TradeService) CreateTariffCode(userID uuid.UUID, req CreateTariffCodeRequest) (*models.TariffCode, error) {
	user, err := s.userRepo.GetUser(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Convert maps and arrays to JSON strings
	importRestrictionJSON, _ := json.Marshal(req.ImportRestriction)
	exportRestrictionJSON, _ := json.Marshal(req.ExportRestriction)
	requiredCertsJSON, _ := json.Marshal(req.RequiredCerts)

	tariffCode := &models.TariffCode{
		CompanyID:         user.CompanyID,
		HSCode:           req.HSCode,
		Description:      req.Description,
		DescriptionEN:    req.DescriptionEN,
		Category:         req.Category,
		Unit:             req.Unit,
		BaseRate:         req.BaseRate,
		PreferentialRate: req.PreferentialRate,
		VAT:              req.VAT,
		ExciseTax:        req.ExciseTax,
		ImportRestriction: string(importRestrictionJSON),
		ExportRestriction: string(exportRestrictionJSON),
		RequiredCerts:    string(requiredCertsJSON),
		IsActive:         true,
		CreatedBy:        userID,
	}

	if err := s.tradeRepo.CreateTariffCode(tariffCode); err != nil {
		return nil, fmt.Errorf("failed to create tariff code: %w", err)
	}

	return s.tradeRepo.GetTariffCode(tariffCode.ID)
}

func (s *TradeService) GetTariffCode(id uuid.UUID) (*models.TariffCode, error) {
	return s.tradeRepo.GetTariffCode(id)
}

func (s *TradeService) GetTariffCodesByCompany(companyID uuid.UUID, hsCode, category string, isActive *bool) ([]models.TariffCode, error) {
	return s.tradeRepo.GetTariffCodesByCompany(companyID, hsCode, category, isActive)
}

func (s *TradeService) UpdateTariffCode(id, userID uuid.UUID, req UpdateTariffCodeRequest) (*models.TariffCode, error) {
	tariffCode, err := s.tradeRepo.GetTariffCode(id)
	if err != nil {
		return nil, fmt.Errorf("tariff code not found: %w", err)
	}

	// Update fields if provided
	if req.Description != nil {
		tariffCode.Description = *req.Description
	}
	if req.DescriptionEN != nil {
		tariffCode.DescriptionEN = *req.DescriptionEN
	}
	if req.Category != nil {
		tariffCode.Category = *req.Category
	}
	if req.Unit != nil {
		tariffCode.Unit = *req.Unit
	}
	if req.BaseRate != nil {
		tariffCode.BaseRate = *req.BaseRate
	}
	if req.PreferentialRate != nil {
		tariffCode.PreferentialRate = *req.PreferentialRate
	}
	if req.VAT != nil {
		tariffCode.VAT = *req.VAT
	}
	if req.ExciseTax != nil {
		tariffCode.ExciseTax = *req.ExciseTax
	}
	if req.ImportRestriction != nil {
		importRestrictionJSON, _ := json.Marshal(req.ImportRestriction)
		tariffCode.ImportRestriction = string(importRestrictionJSON)
	}
	if req.ExportRestriction != nil {
		exportRestrictionJSON, _ := json.Marshal(req.ExportRestriction)
		tariffCode.ExportRestriction = string(exportRestrictionJSON)
	}
	if req.RequiredCerts != nil {
		requiredCertsJSON, _ := json.Marshal(req.RequiredCerts)
		tariffCode.RequiredCerts = string(requiredCertsJSON)
	}
	if req.IsActive != nil {
		tariffCode.IsActive = *req.IsActive
	}

	if err := s.tradeRepo.UpdateTariffCode(tariffCode); err != nil {
		return nil, fmt.Errorf("failed to update tariff code: %w", err)
	}

	return s.tradeRepo.GetTariffCode(id)
}

func (s *TradeService) DeleteTariffCode(id uuid.UUID) error {
	return s.tradeRepo.DeleteTariffCode(id)
}

// TariffRate Service Methods
type CreateTariffRateRequest struct {
	TariffCodeID  string    `json:"tariff_code_id" validate:"required"`
	CountryCode   string    `json:"country_code" validate:"required"`
	CountryName   string    `json:"country_name" validate:"required"`
	Rate          float64   `json:"rate" validate:"required"`
	RateType      string    `json:"rate_type" validate:"required"`
	MinimumDuty   float64   `json:"minimum_duty"`
	MaximumDuty   float64   `json:"maximum_duty"`
	Currency      string    `json:"currency"`
	TradeType     string    `json:"trade_type" validate:"required"`
	AgreementType string    `json:"agreement_type"`
	ValidFrom     time.Time `json:"valid_from" validate:"required"`
	ValidTo       *time.Time `json:"valid_to"`
}

func (s *TradeService) CreateTariffRate(userID uuid.UUID, req CreateTariffRateRequest) (*models.TariffRate, error) {
	user, err := s.userRepo.GetUser(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	tariffCodeID, err := uuid.Parse(req.TariffCodeID)
	if err != nil {
		return nil, fmt.Errorf("invalid tariff code ID: %w", err)
	}

	// Verify tariff code exists
	_, err = s.tradeRepo.GetTariffCode(tariffCodeID)
	if err != nil {
		return nil, fmt.Errorf("tariff code not found: %w", err)
	}

	if req.Currency == "" {
		req.Currency = "USD"
	}

	tariffRate := &models.TariffRate{
		CompanyID:     user.CompanyID,
		TariffCodeID:  tariffCodeID,
		CountryCode:   req.CountryCode,
		CountryName:   req.CountryName,
		Rate:          req.Rate,
		RateType:      req.RateType,
		MinimumDuty:   req.MinimumDuty,
		MaximumDuty:   req.MaximumDuty,
		Currency:      req.Currency,
		TradeType:     req.TradeType,
		AgreementType: req.AgreementType,
		ValidFrom:     req.ValidFrom,
		ValidTo:       req.ValidTo,
		IsActive:      true,
		CreatedBy:     userID,
	}

	if err := s.tradeRepo.CreateTariffRate(tariffRate); err != nil {
		return nil, fmt.Errorf("failed to create tariff rate: %w", err)
	}

	return s.tradeRepo.GetTariffRate(tariffRate.ID)
}

func (s *TradeService) GetTariffRatesByTariffCode(tariffCodeID uuid.UUID, countryCode, tradeType, agreementType string) ([]models.TariffRate, error) {
	return s.tradeRepo.GetTariffRatesByTariffCode(tariffCodeID, countryCode, tradeType, agreementType)
}

func (s *TradeService) GetTariffRatesByCompany(companyID uuid.UUID, countryCode, tradeType string) ([]models.TariffRate, error) {
	return s.tradeRepo.GetTariffRatesByCompany(companyID, countryCode, tradeType)
}

// Shipment Service Methods
type CreateShipmentRequest struct {
	OrderID            *string   `json:"order_id"`
	ShipmentNo         string    `json:"shipment_no" validate:"required"`
	Type               string    `json:"type" validate:"required"`
	Method             string    `json:"method" validate:"required"`
	CarrierName        string    `json:"carrier_name"`
	TrackingNo         string    `json:"tracking_no"`
	ContainerNo        string    `json:"container_no"`
	ContainerType      string    `json:"container_type"`
	OriginCountry      string    `json:"origin_country" validate:"required"`
	OriginPort         string    `json:"origin_port"`
	OriginAddress      string    `json:"origin_address"`
	DestCountry        string    `json:"dest_country" validate:"required"`
	DestPort           string    `json:"dest_port"`
	DestAddress        string    `json:"dest_address"`
	EstimatedDeparture *time.Time `json:"estimated_departure"`
	EstimatedArrival   *time.Time `json:"estimated_arrival"`
	GrossWeight        float64   `json:"gross_weight"`
	NetWeight          float64   `json:"net_weight"`
	Volume             float64   `json:"volume"`
	PackageCount       int       `json:"package_count"`
	PackageType        string    `json:"package_type"`
	InsuranceValue     float64   `json:"insurance_value"`
	InsuranceCurrency  string    `json:"insurance_currency"`
	FreightCost        float64   `json:"freight_cost"`
	FreightCurrency    string    `json:"freight_currency"`
	CustomsValue       float64   `json:"customs_value"`
	CustomsCurrency    string    `json:"customs_currency"`
	SpecialInstructions string   `json:"special_instructions"`
	InternalNotes      string    `json:"internal_notes"`
	Items              []CreateShipmentItemRequest `json:"items"`
}

type CreateShipmentItemRequest struct {
	ProductID     *string `json:"product_id"`
	HSCode        string  `json:"hs_code"`
	ProductName   string  `json:"product_name" validate:"required"`
	Description   string  `json:"description"`
	Quantity      float64 `json:"quantity" validate:"required"`
	Unit          string  `json:"unit" validate:"required"`
	UnitWeight    float64 `json:"unit_weight"`
	UnitValue     float64 `json:"unit_value"`
	Currency      string  `json:"currency"`
	CountryOrigin string  `json:"country_origin"`
	Manufacturer  string  `json:"manufacturer"`
}

type UpdateShipmentRequest struct {
	Status              *string    `json:"status"`
	TrackingNo          *string    `json:"tracking_no"`
	ActualDeparture     *time.Time `json:"actual_departure"`
	ActualArrival       *time.Time `json:"actual_arrival"`
	TotalDuty           *float64   `json:"total_duty"`
	TotalTax            *float64   `json:"total_tax"`
	SpecialInstructions *string    `json:"special_instructions"`
	InternalNotes       *string    `json:"internal_notes"`
}

func (s *TradeService) CreateShipment(userID uuid.UUID, req CreateShipmentRequest) (*models.Shipment, error) {
	user, err := s.userRepo.GetUser(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	var orderID *uuid.UUID
	if req.OrderID != nil && *req.OrderID != "" {
		id, err := uuid.Parse(*req.OrderID)
		if err != nil {
			return nil, fmt.Errorf("invalid order ID: %w", err)
		}
		orderID = &id
	}

	shipment := &models.Shipment{
		CompanyID:           user.CompanyID,
		ShipmentNo:          req.ShipmentNo,
		OrderID:             orderID,
		Type:                req.Type,
		Status:              "pending",
		Method:              req.Method,
		CarrierName:         req.CarrierName,
		TrackingNo:          req.TrackingNo,
		ContainerNo:         req.ContainerNo,
		ContainerType:       req.ContainerType,
		OriginCountry:       req.OriginCountry,
		OriginPort:          req.OriginPort,
		OriginAddress:       req.OriginAddress,
		DestCountry:         req.DestCountry,
		DestPort:            req.DestPort,
		DestAddress:         req.DestAddress,
		EstimatedDeparture:  req.EstimatedDeparture,
		EstimatedArrival:    req.EstimatedArrival,
		GrossWeight:         req.GrossWeight,
		NetWeight:           req.NetWeight,
		Volume:              req.Volume,
		PackageCount:        req.PackageCount,
		PackageType:         req.PackageType,
		InsuranceValue:      req.InsuranceValue,
		InsuranceCurrency:   req.InsuranceCurrency,
		FreightCost:         req.FreightCost,
		FreightCurrency:     req.FreightCurrency,
		CustomsValue:        req.CustomsValue,
		CustomsCurrency:     req.CustomsCurrency,
		SpecialInstructions: req.SpecialInstructions,
		InternalNotes:       req.InternalNotes,
		CreatedBy:           userID,
	}

	if err := s.tradeRepo.CreateShipment(shipment); err != nil {
		return nil, fmt.Errorf("failed to create shipment: %w", err)
	}

	// Create shipment items
	for _, itemReq := range req.Items {
		var productID *uuid.UUID
		if itemReq.ProductID != nil && *itemReq.ProductID != "" {
			id, err := uuid.Parse(*itemReq.ProductID)
			if err == nil {
				productID = &id
			}
		}

		if itemReq.Currency == "" {
			itemReq.Currency = "USD"
		}

		item := &models.ShipmentItem{
			CompanyID:     user.CompanyID,
			ShipmentID:    shipment.ID,
			ProductID:     productID,
			HSCode:        itemReq.HSCode,
			ProductName:   itemReq.ProductName,
			Description:   itemReq.Description,
			Quantity:      itemReq.Quantity,
			Unit:          itemReq.Unit,
			UnitWeight:    itemReq.UnitWeight,
			UnitValue:     itemReq.UnitValue,
			Currency:      itemReq.Currency,
			CountryOrigin: itemReq.CountryOrigin,
			Manufacturer:  itemReq.Manufacturer,
		}

		if err := s.tradeRepo.CreateShipmentItem(item); err != nil {
			return nil, fmt.Errorf("failed to create shipment item: %w", err)
		}
	}

	// Create initial shipment event
	event := &models.ShipmentEvent{
		CompanyID:   user.CompanyID,
		ShipmentID:  shipment.ID,
		EventType:   "created",
		Status:      "completed",
		Description: "Shipment created",
		EventTime:   time.Now(),
		RecordedAt:  time.Now(),
		Source:      "manual",
		CreatedBy:   userID,
	}

	if err := s.tradeRepo.CreateShipmentEvent(event); err != nil {
		// Log error but don't fail the shipment creation
		fmt.Printf("Warning: failed to create shipment event: %v\n", err)
	}

	return s.tradeRepo.GetShipment(shipment.ID)
}

func (s *TradeService) GetShipment(id uuid.UUID) (*models.Shipment, error) {
	return s.tradeRepo.GetShipment(id)
}

func (s *TradeService) GetShipmentsByCompany(companyID uuid.UUID, shipmentType, status, method string) ([]models.Shipment, error) {
	return s.tradeRepo.GetShipmentsByCompany(companyID, shipmentType, status, method)
}

func (s *TradeService) UpdateShipment(id, userID uuid.UUID, req UpdateShipmentRequest) (*models.Shipment, error) {
	shipment, err := s.tradeRepo.GetShipment(id)
	if err != nil {
		return nil, fmt.Errorf("shipment not found: %w", err)
	}

	// Track status changes for events
	var statusChanged bool
	var oldStatus string
	if req.Status != nil && *req.Status != shipment.Status {
		statusChanged = true
		oldStatus = shipment.Status
		shipment.Status = *req.Status
	}

	// Update other fields
	if req.TrackingNo != nil {
		shipment.TrackingNo = *req.TrackingNo
	}
	if req.ActualDeparture != nil {
		shipment.ActualDeparture = req.ActualDeparture
	}
	if req.ActualArrival != nil {
		shipment.ActualArrival = req.ActualArrival
	}
	if req.TotalDuty != nil {
		shipment.TotalDuty = *req.TotalDuty
	}
	if req.TotalTax != nil {
		shipment.TotalTax = *req.TotalTax
	}
	if req.SpecialInstructions != nil {
		shipment.SpecialInstructions = *req.SpecialInstructions
	}
	if req.InternalNotes != nil {
		shipment.InternalNotes = *req.InternalNotes
	}

	if err := s.tradeRepo.UpdateShipment(shipment); err != nil {
		return nil, fmt.Errorf("failed to update shipment: %w", err)
	}

	// Create status change event if status changed
	if statusChanged {
		event := &models.ShipmentEvent{
			CompanyID:   shipment.CompanyID,
			ShipmentID:  shipment.ID,
			EventType:   "status_change",
			Status:      "completed",
			Description: fmt.Sprintf("Status changed from %s to %s", oldStatus, *req.Status),
			EventTime:   time.Now(),
			RecordedAt:  time.Now(),
			Source:      "manual",
			CreatedBy:   userID,
		}

		if err := s.tradeRepo.CreateShipmentEvent(event); err != nil {
			fmt.Printf("Warning: failed to create status change event: %v\n", err)
		}
	}

	return s.tradeRepo.GetShipment(id)
}

// ShipmentEvent Service Methods
type CreateShipmentEventRequest struct {
	EventType   string    `json:"event_type" validate:"required"`
	Status      string    `json:"status" validate:"required"`
	Location    string    `json:"location"`
	Description string    `json:"description"`
	Longitude   float64   `json:"longitude"`
	Latitude    float64   `json:"latitude"`
	EventTime   time.Time `json:"event_time"`
	Source      string    `json:"source"`
}

func (s *TradeService) CreateShipmentEvent(shipmentID, userID uuid.UUID, req CreateShipmentEventRequest) (*models.ShipmentEvent, error) {
	shipment, err := s.tradeRepo.GetShipment(shipmentID)
	if err != nil {
		return nil, fmt.Errorf("shipment not found: %w", err)
	}

	if req.Source == "" {
		req.Source = "manual"
	}

	event := &models.ShipmentEvent{
		CompanyID:   shipment.CompanyID,
		ShipmentID:  shipmentID,
		EventType:   req.EventType,
		Status:      req.Status,
		Location:    req.Location,
		Description: req.Description,
		Longitude:   req.Longitude,
		Latitude:    req.Latitude,
		EventTime:   req.EventTime,
		RecordedAt:  time.Now(),
		Source:      req.Source,
		CreatedBy:   userID,
	}

	if err := s.tradeRepo.CreateShipmentEvent(event); err != nil {
		return nil, fmt.Errorf("failed to create shipment event: %w", err)
	}

	return s.tradeRepo.GetShipmentEvent(event.ID)
}

func (s *TradeService) GetShipmentEventsByShipment(shipmentID uuid.UUID) ([]models.ShipmentEvent, error) {
	return s.tradeRepo.GetShipmentEventsByShipment(shipmentID)
}

// LetterOfCredit Service Methods
type CreateLetterOfCreditRequest struct {
	LCNumber           string     `json:"lc_number" validate:"required"`
	Type               string     `json:"type" validate:"required"`
	Amount             float64    `json:"amount" validate:"required"`
	Currency           string     `json:"currency" validate:"required"`
	ApplicantName      string     `json:"applicant_name" validate:"required"`
	ApplicantAddress   string     `json:"applicant_address"`
	BeneficiaryName    string     `json:"beneficiary_name" validate:"required"`
	BeneficiaryAddress string     `json:"beneficiary_address"`
	IssuingBank        string     `json:"issuing_bank" validate:"required"`
	AdvisingBank       string     `json:"advising_bank"`
	ConfirmingBank     string     `json:"confirming_bank"`
	IssueDate          time.Time  `json:"issue_date" validate:"required"`
	ExpiryDate         time.Time  `json:"expiry_date" validate:"required"`
	LastShipmentDate   *time.Time `json:"last_shipment_date"`
	PartialShipment    bool       `json:"partial_shipment"`
	Transhipment       bool       `json:"transhipment"`
	PortOfLoading      string     `json:"port_of_loading"`
	PortOfDischarge    string     `json:"port_of_discharge"`
	Description        string     `json:"description"`
	Documents          []string   `json:"documents"`
	Terms              map[string]interface{} `json:"terms"`
}

func (s *TradeService) CreateLetterOfCredit(userID uuid.UUID, req CreateLetterOfCreditRequest) (*models.LetterOfCredit, error) {
	user, err := s.userRepo.GetUser(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Convert slices and maps to JSON strings
	documentsJSON, _ := json.Marshal(req.Documents)
	termsJSON, _ := json.Marshal(req.Terms)

	lc := &models.LetterOfCredit{
		CompanyID:          user.CompanyID,
		LCNumber:           req.LCNumber,
		Type:               req.Type,
		Status:             "draft",
		Amount:             req.Amount,
		Currency:           req.Currency,
		ApplicantName:      req.ApplicantName,
		ApplicantAddress:   req.ApplicantAddress,
		BeneficiaryName:    req.BeneficiaryName,
		BeneficiaryAddress: req.BeneficiaryAddress,
		IssuingBank:        req.IssuingBank,
		AdvisingBank:       req.AdvisingBank,
		ConfirmingBank:     req.ConfirmingBank,
		IssueDate:          req.IssueDate,
		ExpiryDate:         req.ExpiryDate,
		LastShipmentDate:   req.LastShipmentDate,
		PartialShipment:    req.PartialShipment,
		Transhipment:       req.Transhipment,
		PortOfLoading:      req.PortOfLoading,
		PortOfDischarge:    req.PortOfDischarge,
		Description:        req.Description,
		Documents:          string(documentsJSON),
		Terms:              string(termsJSON),
		UtilizedAmount:     0,
		AvailableAmount:    req.Amount,
		CreatedBy:          userID,
	}

	if err := s.tradeRepo.CreateLetterOfCredit(lc); err != nil {
		return nil, fmt.Errorf("failed to create letter of credit: %w", err)
	}

	return s.tradeRepo.GetLetterOfCredit(lc.ID)
}

func (s *TradeService) GetLetterOfCredit(id uuid.UUID) (*models.LetterOfCredit, error) {
	return s.tradeRepo.GetLetterOfCredit(id)
}

func (s *TradeService) GetLetterOfCreditsByCompany(companyID uuid.UUID, lcType, status string) ([]models.LetterOfCredit, error) {
	return s.tradeRepo.GetLetterOfCreditsByCompany(companyID, lcType, status)
}

func (s *TradeService) GetExpiringLetterOfCredits(companyID uuid.UUID, days int) ([]models.LetterOfCredit, error) {
	return s.tradeRepo.GetExpiringLetterOfCredits(companyID, days)
}

// LCUtilization Service Methods
type CreateLCUtilizationRequest struct {
	LCID         string                 `json:"lc_id" validate:"required"`
	ShipmentID   *string                `json:"shipment_id"`
	Amount       float64                `json:"amount" validate:"required"`
	Currency     string                 `json:"currency" validate:"required"`
	Description  string                 `json:"description"`
	DocumentsRef map[string]interface{} `json:"documents_ref"`
	UtilizedAt   time.Time              `json:"utilized_at"`
}

func (s *TradeService) CreateLCUtilization(userID uuid.UUID, req CreateLCUtilizationRequest) (*models.LCUtilization, error) {
	user, err := s.userRepo.GetUser(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	lcID, err := uuid.Parse(req.LCID)
	if err != nil {
		return nil, fmt.Errorf("invalid LC ID: %w", err)
	}

	// Verify LC exists and has sufficient available amount
	lc, err := s.tradeRepo.GetLetterOfCredit(lcID)
	if err != nil {
		return nil, fmt.Errorf("letter of credit not found: %w", err)
	}

	if lc.AvailableAmount < req.Amount {
		return nil, fmt.Errorf("insufficient available amount in letter of credit")
	}

	var shipmentID *uuid.UUID
	if req.ShipmentID != nil && *req.ShipmentID != "" {
		id, err := uuid.Parse(*req.ShipmentID)
		if err != nil {
			return nil, fmt.Errorf("invalid shipment ID: %w", err)
		}
		shipmentID = &id
	}

	documentsRefJSON, _ := json.Marshal(req.DocumentsRef)

	utilization := &models.LCUtilization{
		CompanyID:    user.CompanyID,
		LCID:         lcID,
		ShipmentID:   shipmentID,
		Amount:       req.Amount,
		Currency:     req.Currency,
		Description:  req.Description,
		DocumentsRef: string(documentsRefJSON),
		Status:       "pending",
		UtilizedAt:   req.UtilizedAt,
		CreatedBy:    userID,
	}

	if err := s.tradeRepo.CreateLCUtilization(utilization); err != nil {
		return nil, fmt.Errorf("failed to create LC utilization: %w", err)
	}

	return s.tradeRepo.GetLCUtilization(utilization.ID)
}

func (s *TradeService) GetLCUtilizationsByLC(lcID uuid.UUID, status string) ([]models.LCUtilization, error) {
	return s.tradeRepo.GetLCUtilizationsByLC(lcID, status)
}

// Compliance Service Methods
type CreateComplianceCheckRequest struct {
	ComplianceID uuid.UUID `json:"compliance_id" validate:"required"`
	ResourceType string    `json:"resource_type" validate:"required"`
	ResourceID   uuid.UUID `json:"resource_id" validate:"required"`
	CheckType    string    `json:"check_type" validate:"required"`
}

func (s *TradeService) RunComplianceCheck(userID uuid.UUID, req CreateComplianceCheckRequest) (*models.ComplianceCheck, error) {
	user, err := s.userRepo.GetUser(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Get compliance rule
	compliance, err := s.tradeRepo.GetTradeCompliance(req.ComplianceID)
	if err != nil {
		return nil, fmt.Errorf("compliance rule not found: %w", err)
	}

	// Run compliance check logic (simplified)
	result := "passed"
	score := 100.0
	var issues []string
	var recommendations []string

	// Simulate compliance checking based on type and severity
	switch compliance.Severity {
	case "critical":
		if req.CheckType == "automatic" {
			// Simulate some failures for critical compliance
			score = 85.0
			if score < 90 {
				result = "warning"
				issues = append(issues, "Some critical compliance requirements may not be met")
				recommendations = append(recommendations, "Review all critical compliance documentation")
			}
		}
	case "high":
		score = 95.0
	}

	issuesJSON, _ := json.Marshal(issues)
	recommendationsJSON, _ := json.Marshal(recommendations)

	check := &models.ComplianceCheck{
		CompanyID:       user.CompanyID,
		ComplianceID:    req.ComplianceID,
		ResourceType:    req.ResourceType,
		ResourceID:      req.ResourceID,
		CheckType:       req.CheckType,
		Result:          result,
		Score:           score,
		Issues:          string(issuesJSON),
		Recommendations: string(recommendationsJSON),
		CheckedAt:       time.Now(),
		CheckedBy:       &userID,
	}

	if err := s.tradeRepo.CreateComplianceCheck(check); err != nil {
		return nil, fmt.Errorf("failed to create compliance check: %w", err)
	}

	return s.tradeRepo.GetComplianceCheck(check.ID)
}

func (s *TradeService) GetComplianceChecksByResource(companyID uuid.UUID, resourceType string, resourceID uuid.UUID) ([]models.ComplianceCheck, error) {
	return s.tradeRepo.GetComplianceChecksByResource(companyID, resourceType, resourceID)
}

func (s *TradeService) GetFailedComplianceChecks(companyID uuid.UUID) ([]models.ComplianceCheck, error) {
	return s.tradeRepo.GetFailedComplianceChecks(companyID)
}

// ExchangeRate Service Methods
type CreateExchangeRateRequest struct {
	FromCurrency string    `json:"from_currency" validate:"required"`
	ToCurrency   string    `json:"to_currency" validate:"required"`
	Rate         float64   `json:"rate" validate:"required"`
	RateType     string    `json:"rate_type" validate:"required"`
	Source       string    `json:"source" validate:"required"`
	ValidDate    time.Time `json:"valid_date" validate:"required"`
}

func (s *TradeService) CreateExchangeRate(userID uuid.UUID, req CreateExchangeRateRequest) (*models.ExchangeRate, error) {
	user, err := s.userRepo.GetUser(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Deactivate previous rates for the same currency pair and type
	existingRates, err := s.tradeRepo.GetExchangeRatesByCompany(user.CompanyID, req.FromCurrency, req.ToCurrency, req.RateType)
	if err == nil {
		for _, rate := range existingRates {
			if rate.ValidDate.Before(req.ValidDate) {
				rate.IsActive = false
				s.tradeRepo.UpdateExchangeRate(&rate)
			}
		}
	}

	exchangeRate := &models.ExchangeRate{
		CompanyID:    user.CompanyID,
		FromCurrency: req.FromCurrency,
		ToCurrency:   req.ToCurrency,
		Rate:         req.Rate,
		RateType:     req.RateType,
		Source:       req.Source,
		ValidDate:    req.ValidDate,
		IsActive:     true,
		CreatedBy:    userID,
	}

	if err := s.tradeRepo.CreateExchangeRate(exchangeRate); err != nil {
		return nil, fmt.Errorf("failed to create exchange rate: %w", err)
	}

	return s.tradeRepo.GetExchangeRate(exchangeRate.ID)
}

func (s *TradeService) GetLatestExchangeRate(companyID uuid.UUID, fromCurrency, toCurrency, rateType string) (*models.ExchangeRate, error) {
	return s.tradeRepo.GetLatestExchangeRate(companyID, fromCurrency, toCurrency, rateType)
}

func (s *TradeService) GetExchangeRatesByCompany(companyID uuid.UUID, fromCurrency, toCurrency, rateType string) ([]models.ExchangeRate, error) {
	return s.tradeRepo.GetExchangeRatesByCompany(companyID, fromCurrency, toCurrency, rateType)
}

// Analytics and Statistics Methods
func (s *TradeService) GetTradeStatistics(companyID uuid.UUID, startDate, endDate time.Time) (map[string]interface{}, error) {
	return s.tradeRepo.GetTradeStatistics(companyID, startDate, endDate)
}

func (s *TradeService) GetShipmentsByCountry(companyID uuid.UUID, startDate, endDate time.Time) ([]map[string]interface{}, error) {
	return s.tradeRepo.GetShipmentsByCountry(companyID, startDate, endDate)
}

func (s *TradeService) GetTopTradingPartners(companyID uuid.UUID, startDate, endDate time.Time, limit int) ([]map[string]interface{}, error) {
	return s.tradeRepo.GetTopTradingPartners(companyID, startDate, endDate, limit)
}

// Utility methods
func (s *TradeService) CalculateTariffDuty(companyID uuid.UUID, hsCode, countryCode string, tradeType string, value float64) (*map[string]interface{}, error) {
	// Get tariff code
	tariffCode, err := s.tradeRepo.GetTariffCodeByHSCode(companyID, hsCode)
	if err != nil {
		return nil, fmt.Errorf("tariff code not found: %w", err)
	}

	// Get applicable tariff rates
	rates, err := s.tradeRepo.GetTariffRatesByTariffCode(tariffCode.ID, countryCode, tradeType, "")
	if err != nil || len(rates) == 0 {
		// Use base rate if no specific rate found
		duty := value * (tariffCode.BaseRate / 100)
		vat := value * (tariffCode.VAT / 100)
		exciseTax := value * (tariffCode.ExciseTax / 100)
		
		return &map[string]interface{}{
			"tariff_code": tariffCode,
			"base_rate":   tariffCode.BaseRate,
			"duty":        duty,
			"vat":         vat,
			"excise_tax":  exciseTax,
			"total":       duty + vat + exciseTax,
		}, nil
	}

	// Use the first (best) rate found
	rate := rates[0]
	var duty float64

	switch rate.RateType {
	case "ad_valorem":
		duty = value * (rate.Rate / 100)
	case "specific":
		duty = rate.Rate // Assume rate is per unit, would need quantity
	case "compound":
		duty = value * (rate.Rate / 100) // Simplified
	}

	// Apply minimum/maximum duty if specified
	if rate.MinimumDuty > 0 && duty < rate.MinimumDuty {
		duty = rate.MinimumDuty
	}
	if rate.MaximumDuty > 0 && duty > rate.MaximumDuty {
		duty = rate.MaximumDuty
	}

	vat := value * (tariffCode.VAT / 100)
	exciseTax := value * (tariffCode.ExciseTax / 100)

	return &map[string]interface{}{
		"tariff_code":  tariffCode,
		"tariff_rate":  rate,
		"applied_rate": rate.Rate,
		"duty":         duty,
		"vat":          vat,
		"excise_tax":   exciseTax,
		"total":        duty + vat + exciseTax,
	}, nil
}

func (s *TradeService) ConvertCurrency(companyID uuid.UUID, amount float64, fromCurrency, toCurrency string) (*map[string]interface{}, error) {
	if fromCurrency == toCurrency {
		return &map[string]interface{}{
			"original_amount":    amount,
			"original_currency":  fromCurrency,
			"converted_amount":   amount,
			"converted_currency": toCurrency,
			"rate":              1.0,
			"rate_date":         time.Now(),
		}, nil
	}

	// Get latest exchange rate
	rate, err := s.tradeRepo.GetLatestExchangeRate(companyID, fromCurrency, toCurrency, "mid")
	if err != nil {
		return nil, fmt.Errorf("exchange rate not found: %w", err)
	}

	convertedAmount := amount * rate.Rate

	return &map[string]interface{}{
		"original_amount":    amount,
		"original_currency":  fromCurrency,
		"converted_amount":   convertedAmount,
		"converted_currency": toCurrency,
		"rate":              rate.Rate,
		"rate_date":         rate.ValidDate,
		"rate_source":       rate.Source,
	}, nil
}

// Document management helper
func (s *TradeService) GetTradeDocumentsByShipment(shipmentID uuid.UUID) ([]models.TradeDocument, error) {
	return s.tradeRepo.GetTradeDocumentsByShipment(shipmentID)
}