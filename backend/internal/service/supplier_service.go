package service

import (
	"fmt"
	"time"

	"github.com/fastenmind/fastener-api/internal/models"
	"github.com/fastenmind/fastener-api/internal/repository"
	"github.com/google/uuid"
)

type SupplierService interface {
	// Supplier operations
	CreateSupplier(req *CreateSupplierRequest, userID uuid.UUID) (*models.Supplier, error)
	UpdateSupplier(id uuid.UUID, req *UpdateSupplierRequest, userID uuid.UUID) (*models.Supplier, error)
	GetSupplier(id uuid.UUID) (*models.Supplier, error)
	ListSuppliers(companyID uuid.UUID, params map[string]interface{}) ([]models.Supplier, int64, error)
	
	// Supplier Contact operations
	AddSupplierContact(supplierID uuid.UUID, req *CreateSupplierContactRequest) (*models.SupplierContact, error)
	UpdateSupplierContact(id uuid.UUID, req *UpdateSupplierContactRequest) (*models.SupplierContact, error)
	GetSupplierContacts(supplierID uuid.UUID) ([]models.SupplierContact, error)
	DeleteSupplierContact(id uuid.UUID) error
	
	// Supplier Product operations
	AddSupplierProduct(supplierID uuid.UUID, req *CreateSupplierProductRequest) (*models.SupplierProduct, error)
	UpdateSupplierProduct(id uuid.UUID, req *UpdateSupplierProductRequest) (*models.SupplierProduct, error)
	GetSupplierProducts(supplierID uuid.UUID, params map[string]interface{}) ([]models.SupplierProduct, error)
	DeleteSupplierProduct(id uuid.UUID) error
	
	// Purchase Order operations
	CreatePurchaseOrder(req *CreatePurchaseOrderRequest, userID uuid.UUID) (*models.PurchaseOrder, error)
	UpdatePurchaseOrder(id uuid.UUID, req *UpdatePurchaseOrderRequest, userID uuid.UUID) (*models.PurchaseOrder, error)
	GetPurchaseOrder(id uuid.UUID) (*models.PurchaseOrder, error)
	ListPurchaseOrders(companyID uuid.UUID, params map[string]interface{}) ([]models.PurchaseOrder, int64, error)
	ApprovePurchaseOrder(id uuid.UUID, userID uuid.UUID) error
	SendPurchaseOrder(id uuid.UUID) error
	ReceivePurchaseOrder(id uuid.UUID, items []PurchaseOrderReceiptItem) error
	
	// Purchase Order Item operations
	AddPurchaseOrderItem(purchaseOrderID uuid.UUID, req *CreatePurchaseOrderItemRequest) (*models.PurchaseOrderItem, error)
	UpdatePurchaseOrderItem(id uuid.UUID, req *UpdatePurchaseOrderItemRequest) (*models.PurchaseOrderItem, error)
	GetPurchaseOrderItems(purchaseOrderID uuid.UUID) ([]models.PurchaseOrderItem, error)
	DeletePurchaseOrderItem(id uuid.UUID) error
	
	// Supplier Evaluation operations
	CreateSupplierEvaluation(req *CreateSupplierEvaluationRequest, userID uuid.UUID) (*models.SupplierEvaluation, error)
	UpdateSupplierEvaluation(id uuid.UUID, req *UpdateSupplierEvaluationRequest, userID uuid.UUID) (*models.SupplierEvaluation, error)
	GetSupplierEvaluation(id uuid.UUID) (*models.SupplierEvaluation, error)
	ListSupplierEvaluations(companyID uuid.UUID, params map[string]interface{}) ([]models.SupplierEvaluation, int64, error)
	ApproveSupplierEvaluation(id uuid.UUID, userID uuid.UUID) error
	
	// Business operations
	UpdateSupplierPerformance(supplierID uuid.UUID) error
	CalculateSupplierRisk(supplierID uuid.UUID) (string, error)
	GetSupplierDashboard(companyID uuid.UUID) (*SupplierDashboard, error)
}

type supplierService struct {
	supplierRepo  repository.SupplierRepository
	inventoryRepo repository.InventoryRepository
}

func NewSupplierService(supplierRepo repository.SupplierRepository, inventoryRepo repository.InventoryRepository) SupplierService {
	return &supplierService{
		supplierRepo:  supplierRepo,
		inventoryRepo: inventoryRepo,
	}
}

// Request structs
type CreateSupplierRequest struct {
	Name              string     `json:"name" validate:"required"`
	NameEn            string     `json:"name_en"`
	Type              string     `json:"type" validate:"required"`
	ContactPerson     string     `json:"contact_person"`
	ContactTitle      string     `json:"contact_title"`
	Phone             string     `json:"phone"`
	Mobile            string     `json:"mobile"`
	Email             string     `json:"email"`
	Website           string     `json:"website"`
	Country           string     `json:"country"`
	State             string     `json:"state"`
	City              string     `json:"city"`
	Address           string     `json:"address"`
	PostalCode        string     `json:"postal_code"`
	TaxNumber         string     `json:"tax_number"`
	BusinessLicense   string     `json:"business_license"`
	Industry          string     `json:"industry"`
	Established       *time.Time `json:"established"`
	Employees         int        `json:"employees"`
	AnnualRevenue     float64    `json:"annual_revenue"`
	Currency          string     `json:"currency"`
	PaymentTerms      string     `json:"payment_terms"`
	PaymentMethod     string     `json:"payment_method"`
	CreditLimit       float64    `json:"credit_limit"`
	CreditDays        int        `json:"credit_days"`
	ISO9001           bool       `json:"iso_9001"`
	ISO14001          bool       `json:"iso_14001"`
	TS16949           bool       `json:"ts_16949"`
	OHSAS18001        bool       `json:"ohsas_18001"`
	CustomCert        string     `json:"custom_cert"`
	CertExpiry        *time.Time `json:"cert_expiry"`
	CreditRating      string     `json:"credit_rating"`
	FinancialHealth   string     `json:"financial_health"`
	InsuranceCoverage float64    `json:"insurance_coverage"`
	Description       string     `json:"description"`
	Notes             string     `json:"notes"`
	Tags              string     `json:"tags"`
}

type UpdateSupplierRequest struct {
	Name              *string    `json:"name"`
	NameEn            *string    `json:"name_en"`
	Type              *string    `json:"type"`
	Status            *string    `json:"status"`
	ContactPerson     *string    `json:"contact_person"`
	ContactTitle      *string    `json:"contact_title"`
	Phone             *string    `json:"phone"`
	Mobile            *string    `json:"mobile"`
	Email             *string    `json:"email"`
	Website           *string    `json:"website"`
	Country           *string    `json:"country"`
	State             *string    `json:"state"`
	City              *string    `json:"city"`
	Address           *string    `json:"address"`
	PostalCode        *string    `json:"postal_code"`
	TaxNumber         *string    `json:"tax_number"`
	BusinessLicense   *string    `json:"business_license"`
	Industry          *string    `json:"industry"`
	Established       *time.Time `json:"established"`
	Employees         *int       `json:"employees"`
	AnnualRevenue     *float64   `json:"annual_revenue"`
	Currency          *string    `json:"currency"`
	PaymentTerms      *string    `json:"payment_terms"`
	PaymentMethod     *string    `json:"payment_method"`
	CreditLimit       *float64   `json:"credit_limit"`
	CreditDays        *int       `json:"credit_days"`
	QualityRating     *float64   `json:"quality_rating"`
	DeliveryRating    *float64   `json:"delivery_rating"`
	ServiceRating     *float64   `json:"service_rating"`
	OverallRating     *float64   `json:"overall_rating"`
	ISO9001           *bool      `json:"iso_9001"`
	ISO14001          *bool      `json:"iso_14001"`
	TS16949           *bool      `json:"ts_16949"`
	OHSAS18001        *bool      `json:"ohsas_18001"`
	CustomCert        *string    `json:"custom_cert"`
	CertExpiry        *time.Time `json:"cert_expiry"`
	RiskLevel         *string    `json:"risk_level"`
	RiskFactors       *string    `json:"risk_factors"`
	LastAuditDate     *time.Time `json:"last_audit_date"`
	NextAuditDate     *time.Time `json:"next_audit_date"`
	CreditRating      *string    `json:"credit_rating"`
	FinancialHealth   *string    `json:"financial_health"`
	InsuranceCoverage *float64   `json:"insurance_coverage"`
	Description       *string    `json:"description"`
	Notes             *string    `json:"notes"`
	Tags              *string    `json:"tags"`
}

type CreateSupplierContactRequest struct {
	Name             string `json:"name" validate:"required"`
	Title            string `json:"title"`
	Department       string `json:"department"`
	Phone            string `json:"phone"`
	Mobile           string `json:"mobile"`
	Email            string `json:"email"`
	IsPrimary        bool   `json:"is_primary"`
	Responsibilities string `json:"responsibilities"`
	Languages        string `json:"languages"`
}

type UpdateSupplierContactRequest struct {
	Name             *string `json:"name"`
	Title            *string `json:"title"`
	Department       *string `json:"department"`
	Phone            *string `json:"phone"`
	Mobile           *string `json:"mobile"`
	Email            *string `json:"email"`
	IsPrimary        *bool   `json:"is_primary"`
	IsActive         *bool   `json:"is_active"`
	Responsibilities *string `json:"responsibilities"`
	Languages        *string `json:"languages"`
}

type CreateSupplierProductRequest struct {
	InventoryID       *uuid.UUID `json:"inventory_id"`
	ProductName       string     `json:"product_name" validate:"required"`
	ProductCode       string     `json:"product_code"`
	SupplierPartNo    string     `json:"supplier_part_no"`
	Category          string     `json:"category"`
	Specification     string     `json:"specification"`
	Unit              string     `json:"unit" validate:"required"`
	UnitPrice         float64    `json:"unit_price"`
	Currency          string     `json:"currency"`
	MinOrderQty       float64    `json:"min_order_qty"`
	MaxOrderQty       float64    `json:"max_order_qty"`
	PriceBreaks       string     `json:"price_breaks"`
	LeadTimeDays      int        `json:"lead_time_days"`
	QualityGrade      string     `json:"quality_grade"`
	CertificationReq  bool       `json:"certification_req"`
	Certificates      string     `json:"certificates"`
	IsPreferred       bool       `json:"is_preferred"`
}

type UpdateSupplierProductRequest struct {
	InventoryID         *uuid.UUID `json:"inventory_id"`
	ProductName         *string    `json:"product_name"`
	ProductCode         *string    `json:"product_code"`
	SupplierPartNo      *string    `json:"supplier_part_no"`
	Category            *string    `json:"category"`
	Specification       *string    `json:"specification"`
	Unit                *string    `json:"unit"`
	UnitPrice           *float64   `json:"unit_price"`
	Currency            *string    `json:"currency"`
	MinOrderQty         *float64   `json:"min_order_qty"`
	MaxOrderQty         *float64   `json:"max_order_qty"`
	PriceBreaks         *string    `json:"price_breaks"`
	LeadTimeDays        *int       `json:"lead_time_days"`
	QualityGrade        *string    `json:"quality_grade"`
	CertificationReq    *bool      `json:"certification_req"`
	Certificates        *string    `json:"certificates"`
	Status              *string    `json:"status"`
	IsPreferred         *bool      `json:"is_preferred"`
	LastPurchaseDate    *time.Time `json:"last_purchase_date"`
	LastPurchasePrice   *float64   `json:"last_purchase_price"`
	TotalPurchased      *float64   `json:"total_purchased"`
	QualityIssues       *int       `json:"quality_issues"`
	DeliveryIssues      *int       `json:"delivery_issues"`
}

type CreatePurchaseOrderRequest struct {
	SupplierID      uuid.UUID                          `json:"supplier_id" validate:"required"`
	OrderDate       time.Time                          `json:"order_date"`
	RequiredDate    time.Time                          `json:"required_date"`
	PaymentTerms    string                             `json:"payment_terms"`
	PaymentMethod   string                             `json:"payment_method"`
	ShippingAddress string                             `json:"shipping_address"`
	ShippingMethod  string                             `json:"shipping_method"`
	Currency        string                             `json:"currency"`
	ExchangeRate    float64                            `json:"exchange_rate"`
	TaxRate         float64                            `json:"tax_rate"`
	ShippingCost    float64                            `json:"shipping_cost"`
	Notes           string                             `json:"notes"`
	InternalNotes   string                             `json:"internal_notes"`
	Items           []CreatePurchaseOrderItemRequest   `json:"items" validate:"required,min=1,dive"`
}

type UpdatePurchaseOrderRequest struct {
	Status          *string    `json:"status"`
	RequiredDate    *time.Time `json:"required_date"`
	PromisedDate    *time.Time `json:"promised_date"`
	PaymentTerms    *string    `json:"payment_terms"`
	PaymentMethod   *string    `json:"payment_method"`
	ShippingAddress *string    `json:"shipping_address"`
	ShippingMethod  *string    `json:"shipping_method"`
	TrackingNumber  *string    `json:"tracking_number"`
	Currency        *string    `json:"currency"`
	ExchangeRate    *float64   `json:"exchange_rate"`
	TaxRate         *float64   `json:"tax_rate"`
	ShippingCost    *float64   `json:"shipping_cost"`
	Notes           *string    `json:"notes"`
	InternalNotes   *string    `json:"internal_notes"`
}

type CreatePurchaseOrderItemRequest struct {
	SupplierProductID    *uuid.UUID `json:"supplier_product_id"`
	InventoryID          *uuid.UUID `json:"inventory_id"`
	ProductName          string     `json:"product_name" validate:"required"`
	ProductCode          string     `json:"product_code"`
	SupplierPartNo       string     `json:"supplier_part_no"`
	Specification        string     `json:"specification"`
	OrderedQuantity      float64    `json:"ordered_quantity" validate:"required,gt=0"`
	Unit                 string     `json:"unit" validate:"required"`
	UnitPrice            float64    `json:"unit_price" validate:"required,gte=0"`
	QualityRequirement   string     `json:"quality_requirement"`
	InspectionRequired   bool       `json:"inspection_required"`
}

type UpdatePurchaseOrderItemRequest struct {
	SupplierProductID    *uuid.UUID `json:"supplier_product_id"`
	InventoryID          *uuid.UUID `json:"inventory_id"`
	ProductName          *string    `json:"product_name"`
	ProductCode          *string    `json:"product_code"`
	SupplierPartNo       *string    `json:"supplier_part_no"`
	Specification        *string    `json:"specification"`
	OrderedQuantity      *float64   `json:"ordered_quantity"`
	ReceivedQuantity     *float64   `json:"received_quantity"`
	Unit                 *string    `json:"unit"`
	UnitPrice            *float64   `json:"unit_price"`
	Status               *string    `json:"status"`
	QualityRequirement   *string    `json:"quality_requirement"`
	InspectionRequired   *bool      `json:"inspection_required"`
}

type PurchaseOrderReceiptItem struct {
	ItemID           uuid.UUID `json:"item_id" validate:"required"`
	ReceivedQuantity float64   `json:"received_quantity" validate:"required,gt=0"`
	QualityPassed    bool      `json:"quality_passed"`
	InspectionNotes  string    `json:"inspection_notes"`
}

type CreateSupplierEvaluationRequest struct {
	SupplierID        uuid.UUID `json:"supplier_id" validate:"required"`
	StartDate         time.Time `json:"start_date"`
	EndDate           time.Time `json:"end_date"`
	EvaluationType    string    `json:"evaluation_type" validate:"required"`
	QualityScore      float64   `json:"quality_score" validate:"min=0,max=100"`
	DeliveryScore     float64   `json:"delivery_score" validate:"min=0,max=100"`
	ServiceScore      float64   `json:"service_score" validate:"min=0,max=100"`
	CostScore         float64   `json:"cost_score" validate:"min=0,max=100"`
	TechnicalScore    float64   `json:"technical_score" validate:"min=0,max=100"`
	TotalOrders       int       `json:"total_orders"`
	OnTimeDeliveries  int       `json:"on_time_deliveries"`
	QualityDefects    int       `json:"quality_defects"`
	ServiceIssues     int       `json:"service_issues"`
	CostSavings       float64   `json:"cost_savings"`
	Strengths         string    `json:"strengths"`
	Weaknesses        string    `json:"weaknesses"`
	Recommendations   string    `json:"recommendations"`
	ActionItems       string    `json:"action_items"`
	EvaluatedAt       time.Time `json:"evaluated_at"`
}

type UpdateSupplierEvaluationRequest struct {
	StartDate         *time.Time `json:"start_date"`
	EndDate           *time.Time `json:"end_date"`
	EvaluationType    *string    `json:"evaluation_type"`
	QualityScore      *float64   `json:"quality_score"`
	DeliveryScore     *float64   `json:"delivery_score"`
	ServiceScore      *float64   `json:"service_score"`
	CostScore         *float64   `json:"cost_score"`
	TechnicalScore    *float64   `json:"technical_score"`
	TotalOrders       *int       `json:"total_orders"`
	OnTimeDeliveries  *int       `json:"on_time_deliveries"`
	QualityDefects    *int       `json:"quality_defects"`
	ServiceIssues     *int       `json:"service_issues"`
	CostSavings       *float64   `json:"cost_savings"`
	Strengths         *string    `json:"strengths"`
	Weaknesses        *string    `json:"weaknesses"`
	Recommendations   *string    `json:"recommendations"`
	ActionItems       *string    `json:"action_items"`
	Status            *string    `json:"status"`
}

type SupplierDashboard struct {
	TotalSuppliers       int     `json:"total_suppliers"`
	ActiveSuppliers      int     `json:"active_suppliers"`
	SuspendedSuppliers   int     `json:"suspended_suppliers"`
	BlacklistedSuppliers int     `json:"blacklisted_suppliers"`
	TotalPurchaseOrders  int     `json:"total_purchase_orders"`
	DraftOrders          int     `json:"draft_orders"`
	SentOrders           int     `json:"sent_orders"`
	ConfirmedOrders      int     `json:"confirmed_orders"`
	ReceivedOrders       int     `json:"received_orders"`
	CancelledOrders      int     `json:"cancelled_orders"`
	TotalPurchaseValue   float64 `json:"total_purchase_value"`
	PendingValue         float64 `json:"pending_value"`
	ReceivedValue        float64 `json:"received_value"`
	HighRiskSuppliers    int     `json:"high_risk_suppliers"`
	CriticalRiskSuppliers int    `json:"critical_risk_suppliers"`
	AverageQualityRating  float64 `json:"average_quality_rating"`
	AverageDeliveryRating float64 `json:"average_delivery_rating"`
	AverageServiceRating  float64 `json:"average_service_rating"`
	TotalEvaluations     int     `json:"total_evaluations"`
	PendingEvaluations   int     `json:"pending_evaluations"`
	CompletedEvaluations int     `json:"completed_evaluations"`
}

// Supplier operations
func (s *supplierService) CreateSupplier(req *CreateSupplierRequest, userID uuid.UUID) (*models.Supplier, error) {
	supplier := &models.Supplier{
		CompanyID:         userID, // This should be retrieved from user context
		Name:              req.Name,
		NameEn:            req.NameEn,
		Type:              req.Type,
		Status:            "active",
		ContactPerson:     req.ContactPerson,
		ContactTitle:      req.ContactTitle,
		Phone:             req.Phone,
		Mobile:            req.Mobile,
		Email:             req.Email,
		Website:           req.Website,
		Country:           req.Country,
		State:             req.State,
		City:              req.City,
		Address:           req.Address,
		PostalCode:        req.PostalCode,
		TaxNumber:         req.TaxNumber,
		BusinessLicense:   req.BusinessLicense,
		Industry:          req.Industry,
		Established:       req.Established,
		Employees:         req.Employees,
		AnnualRevenue:     req.AnnualRevenue,
		Currency:          req.Currency,
		PaymentTerms:      req.PaymentTerms,
		PaymentMethod:     req.PaymentMethod,
		CreditLimit:       req.CreditLimit,
		CreditDays:        req.CreditDays,
		ISO9001:           req.ISO9001,
		ISO14001:          req.ISO14001,
		TS16949:           req.TS16949,
		OHSAS18001:        req.OHSAS18001,
		CustomCert:        req.CustomCert,
		CertExpiry:        req.CertExpiry,
		RiskLevel:         "medium",
		CreditRating:      req.CreditRating,
		FinancialHealth:   req.FinancialHealth,
		InsuranceCoverage: req.InsuranceCoverage,
		Description:       req.Description,
		Notes:             req.Notes,
		Tags:              req.Tags,
		CreatedBy:         userID,
	}

	// Generate supplier number
	timestamp := time.Now().Unix()
	supplier.SupplierNo = fmt.Sprintf("SUP%d", timestamp)

	if err := s.supplierRepo.CreateSupplier(supplier); err != nil {
		return nil, fmt.Errorf("failed to create supplier: %w", err)
	}

	return supplier, nil
}

func (s *supplierService) UpdateSupplier(id uuid.UUID, req *UpdateSupplierRequest, userID uuid.UUID) (*models.Supplier, error) {
	supplier, err := s.supplierRepo.GetSupplier(id)
	if err != nil {
		return nil, fmt.Errorf("supplier not found: %w", err)
	}

	// Update fields if provided
	updateSupplierFields(supplier, req)

	if err := s.supplierRepo.UpdateSupplier(supplier); err != nil {
		return nil, fmt.Errorf("failed to update supplier: %w", err)
	}

	return supplier, nil
}

func (s *supplierService) GetSupplier(id uuid.UUID) (*models.Supplier, error) {
	return s.supplierRepo.GetSupplier(id)
}

func (s *supplierService) ListSuppliers(companyID uuid.UUID, params map[string]interface{}) ([]models.Supplier, int64, error) {
	return s.supplierRepo.ListSuppliers(companyID, params)
}

// Supplier Contact operations
func (s *supplierService) AddSupplierContact(supplierID uuid.UUID, req *CreateSupplierContactRequest) (*models.SupplierContact, error) {
	contact := &models.SupplierContact{
		SupplierID:       supplierID,
		Name:             req.Name,
		Title:            req.Title,
		Department:       req.Department,
		Phone:            req.Phone,
		Mobile:           req.Mobile,
		Email:            req.Email,
		IsPrimary:        req.IsPrimary,
		IsActive:         true,
		Responsibilities: req.Responsibilities,
		Languages:        req.Languages,
	}

	if err := s.supplierRepo.CreateSupplierContact(contact); err != nil {
		return nil, fmt.Errorf("failed to create supplier contact: %w", err)
	}

	return contact, nil
}

func (s *supplierService) UpdateSupplierContact(id uuid.UUID, req *UpdateSupplierContactRequest) (*models.SupplierContact, error) {
	contacts, err := s.supplierRepo.GetSupplierContacts(uuid.Nil) // Need to fix this
	if err != nil {
		return nil, fmt.Errorf("contact not found: %w", err)
	}
	
	var contact *models.SupplierContact
	for _, c := range contacts {
		if c.ID == id {
			contact = &c
			break
		}
	}
	
	if contact == nil {
		return nil, fmt.Errorf("contact not found")
	}

	// Update fields if provided
	if req.Name != nil {
		contact.Name = *req.Name
	}
	if req.Title != nil {
		contact.Title = *req.Title
	}
	if req.Department != nil {
		contact.Department = *req.Department
	}
	if req.Phone != nil {
		contact.Phone = *req.Phone
	}
	if req.Mobile != nil {
		contact.Mobile = *req.Mobile
	}
	if req.Email != nil {
		contact.Email = *req.Email
	}
	if req.IsPrimary != nil {
		contact.IsPrimary = *req.IsPrimary
	}
	if req.IsActive != nil {
		contact.IsActive = *req.IsActive
	}
	if req.Responsibilities != nil {
		contact.Responsibilities = *req.Responsibilities
	}
	if req.Languages != nil {
		contact.Languages = *req.Languages
	}

	if err := s.supplierRepo.UpdateSupplierContact(contact); err != nil {
		return nil, fmt.Errorf("failed to update supplier contact: %w", err)
	}

	return contact, nil
}

func (s *supplierService) GetSupplierContacts(supplierID uuid.UUID) ([]models.SupplierContact, error) {
	return s.supplierRepo.GetSupplierContacts(supplierID)
}

func (s *supplierService) DeleteSupplierContact(id uuid.UUID) error {
	return s.supplierRepo.DeleteSupplierContact(id)
}

// Supplier Product operations
func (s *supplierService) AddSupplierProduct(supplierID uuid.UUID, req *CreateSupplierProductRequest) (*models.SupplierProduct, error) {
	product := &models.SupplierProduct{
		SupplierID:       supplierID,
		InventoryID:      req.InventoryID,
		ProductName:      req.ProductName,
		ProductCode:      req.ProductCode,
		SupplierPartNo:   req.SupplierPartNo,
		Category:         req.Category,
		Specification:    req.Specification,
		Unit:             req.Unit,
		UnitPrice:        req.UnitPrice,
		Currency:         req.Currency,
		MinOrderQty:      req.MinOrderQty,
		MaxOrderQty:      req.MaxOrderQty,
		PriceBreaks:      req.PriceBreaks,
		LeadTimeDays:     req.LeadTimeDays,
		QualityGrade:     req.QualityGrade,
		CertificationReq: req.CertificationReq,
		Certificates:     req.Certificates,
		Status:           "active",
		IsPreferred:      req.IsPreferred,
	}

	if err := s.supplierRepo.CreateSupplierProduct(product); err != nil {
		return nil, fmt.Errorf("failed to create supplier product: %w", err)
	}

	return product, nil
}

func (s *supplierService) UpdateSupplierProduct(id uuid.UUID, req *UpdateSupplierProductRequest) (*models.SupplierProduct, error) {
	products, err := s.supplierRepo.GetSupplierProducts(uuid.Nil, map[string]interface{}{}) // Need to fix this
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}
	
	var product *models.SupplierProduct
	for _, p := range products {
		if p.ID == id {
			product = &p
			break
		}
	}
	
	if product == nil {
		return nil, fmt.Errorf("product not found")
	}

	// Update fields if provided
	updateSupplierProductFields(product, req)

	if err := s.supplierRepo.UpdateSupplierProduct(product); err != nil {
		return nil, fmt.Errorf("failed to update supplier product: %w", err)
	}

	return product, nil
}

func (s *supplierService) GetSupplierProducts(supplierID uuid.UUID, params map[string]interface{}) ([]models.SupplierProduct, error) {
	return s.supplierRepo.GetSupplierProducts(supplierID, params)
}

func (s *supplierService) DeleteSupplierProduct(id uuid.UUID) error {
	return s.supplierRepo.DeleteSupplierProduct(id)
}

// Purchase Order operations
func (s *supplierService) CreatePurchaseOrder(req *CreatePurchaseOrderRequest, userID uuid.UUID) (*models.PurchaseOrder, error) {
	order := &models.PurchaseOrder{
		CompanyID:       userID, // This should be retrieved from user context
		SupplierID:      req.SupplierID,
		OrderDate:       req.OrderDate,
		RequiredDate:    req.RequiredDate,
		Status:          "draft",
		PaymentTerms:    req.PaymentTerms,
		PaymentMethod:   req.PaymentMethod,
		ShippingAddress: req.ShippingAddress,
		ShippingMethod:  req.ShippingMethod,
		Currency:        req.Currency,
		ExchangeRate:    req.ExchangeRate,
		TaxRate:         req.TaxRate,
		ShippingCost:    req.ShippingCost,
		Notes:           req.Notes,
		InternalNotes:   req.InternalNotes,
		CreatedBy:       userID,
	}

	// Generate order number
	timestamp := time.Now().Unix()
	order.OrderNo = fmt.Sprintf("PO%d", timestamp)

	// Calculate totals
	var subTotal float64
	for _, itemReq := range req.Items {
		subTotal += itemReq.OrderedQuantity * itemReq.UnitPrice
	}

	order.SubTotal = subTotal
	order.TaxAmount = subTotal * req.TaxRate / 100
	order.TotalAmount = subTotal + order.TaxAmount + req.ShippingCost

	if err := s.supplierRepo.CreatePurchaseOrder(order); err != nil {
		return nil, fmt.Errorf("failed to create purchase order: %w", err)
	}

	// Create order items
	for _, itemReq := range req.Items {
		item := &models.PurchaseOrderItem{
			PurchaseOrderID:      order.ID,
			SupplierProductID:    itemReq.SupplierProductID,
			InventoryID:          itemReq.InventoryID,
			ProductName:          itemReq.ProductName,
			ProductCode:          itemReq.ProductCode,
			SupplierPartNo:       itemReq.SupplierPartNo,
			Specification:        itemReq.Specification,
			OrderedQuantity:      itemReq.OrderedQuantity,
			Unit:                 itemReq.Unit,
			UnitPrice:            itemReq.UnitPrice,
			TotalPrice:           itemReq.OrderedQuantity * itemReq.UnitPrice,
			Status:               "pending",
			QualityRequirement:   itemReq.QualityRequirement,
			InspectionRequired:   itemReq.InspectionRequired,
		}

		if err := s.supplierRepo.CreatePurchaseOrderItem(item); err != nil {
			return nil, fmt.Errorf("failed to create purchase order item: %w", err)
		}
	}

	return order, nil
}

func (s *supplierService) UpdatePurchaseOrder(id uuid.UUID, req *UpdatePurchaseOrderRequest, userID uuid.UUID) (*models.PurchaseOrder, error) {
	order, err := s.supplierRepo.GetPurchaseOrder(id)
	if err != nil {
		return nil, fmt.Errorf("purchase order not found: %w", err)
	}

	// Update fields if provided
	updatePurchaseOrderFields(order, req)

	if err := s.supplierRepo.UpdatePurchaseOrder(order); err != nil {
		return nil, fmt.Errorf("failed to update purchase order: %w", err)
	}

	return order, nil
}

func (s *supplierService) GetPurchaseOrder(id uuid.UUID) (*models.PurchaseOrder, error) {
	return s.supplierRepo.GetPurchaseOrder(id)
}

func (s *supplierService) ListPurchaseOrders(companyID uuid.UUID, params map[string]interface{}) ([]models.PurchaseOrder, int64, error) {
	return s.supplierRepo.ListPurchaseOrders(companyID, params)
}

func (s *supplierService) ApprovePurchaseOrder(id uuid.UUID, userID uuid.UUID) error {
	order, err := s.supplierRepo.GetPurchaseOrder(id)
	if err != nil {
		return fmt.Errorf("purchase order not found: %w", err)
	}

	if order.Status != "draft" {
		return fmt.Errorf("purchase order is not in draft status")
	}

	now := time.Now()
	order.Status = "sent"
	order.ApprovedBy = &userID
	order.ApprovedAt = &now

	return s.supplierRepo.UpdatePurchaseOrder(order)
}

func (s *supplierService) SendPurchaseOrder(id uuid.UUID) error {
	order, err := s.supplierRepo.GetPurchaseOrder(id)
	if err != nil {
		return fmt.Errorf("purchase order not found: %w", err)
	}

	if order.Status != "sent" {
		return fmt.Errorf("purchase order is not approved")
	}

	order.Status = "confirmed"

	return s.supplierRepo.UpdatePurchaseOrder(order)
}

func (s *supplierService) ReceivePurchaseOrder(id uuid.UUID, items []PurchaseOrderReceiptItem) error {
	order, err := s.supplierRepo.GetPurchaseOrder(id)
	if err != nil {
		return fmt.Errorf("purchase order not found: %w", err)
	}

	// Update item received quantities
	for _, receiptItem := range items {
		orderItems, err := s.supplierRepo.GetPurchaseOrderItems(id)
		if err != nil {
			return fmt.Errorf("failed to get purchase order items: %w", err)
		}

		for _, orderItem := range orderItems {
			if orderItem.ID == receiptItem.ItemID {
				orderItem.ReceivedQuantity = receiptItem.ReceivedQuantity
				
				if orderItem.ReceivedQuantity >= orderItem.OrderedQuantity {
					orderItem.Status = "received"
				} else if orderItem.ReceivedQuantity > 0 {
					orderItem.Status = "partial_received"
				}

				if err := s.supplierRepo.UpdatePurchaseOrderItem(&orderItem); err != nil {
					return fmt.Errorf("failed to update purchase order item: %w", err)
				}
				break
			}
		}
	}

	// Update order status
	allReceived := true
	partialReceived := false
	orderItems, _ := s.supplierRepo.GetPurchaseOrderItems(id)
	
	for _, item := range orderItems {
		if item.ReceivedQuantity == 0 {
			allReceived = false
		} else if item.ReceivedQuantity < item.OrderedQuantity {
			allReceived = false
			partialReceived = true
		} else if item.ReceivedQuantity > 0 {
			partialReceived = true
		}
	}

	if allReceived {
		order.Status = "received"
	} else if partialReceived {
		order.Status = "partial_received"
	}

	return s.supplierRepo.UpdatePurchaseOrder(order)
}

// Purchase Order Item operations
func (s *supplierService) AddPurchaseOrderItem(purchaseOrderID uuid.UUID, req *CreatePurchaseOrderItemRequest) (*models.PurchaseOrderItem, error) {
	item := &models.PurchaseOrderItem{
		PurchaseOrderID:    purchaseOrderID,
		SupplierProductID:  req.SupplierProductID,
		InventoryID:        req.InventoryID,
		ProductName:        req.ProductName,
		ProductCode:        req.ProductCode,
		SupplierPartNo:     req.SupplierPartNo,
		Specification:      req.Specification,
		OrderedQuantity:    req.OrderedQuantity,
		Unit:               req.Unit,
		UnitPrice:          req.UnitPrice,
		TotalPrice:         req.OrderedQuantity * req.UnitPrice,
		Status:             "pending",
		QualityRequirement: req.QualityRequirement,
		InspectionRequired: req.InspectionRequired,
	}

	if err := s.supplierRepo.CreatePurchaseOrderItem(item); err != nil {
		return nil, fmt.Errorf("failed to create purchase order item: %w", err)
	}

	return item, nil
}

func (s *supplierService) UpdatePurchaseOrderItem(id uuid.UUID, req *UpdatePurchaseOrderItemRequest) (*models.PurchaseOrderItem, error) {
	items, err := s.supplierRepo.GetPurchaseOrderItems(uuid.Nil) // Need to fix this
	if err != nil {
		return nil, fmt.Errorf("item not found: %w", err)
	}
	
	var item *models.PurchaseOrderItem
	for _, i := range items {
		if i.ID == id {
			item = &i
			break
		}
	}
	
	if item == nil {
		return nil, fmt.Errorf("item not found")
	}

	// Update fields if provided
	updatePurchaseOrderItemFields(item, req)

	if err := s.supplierRepo.UpdatePurchaseOrderItem(item); err != nil {
		return nil, fmt.Errorf("failed to update purchase order item: %w", err)
	}

	return item, nil
}

func (s *supplierService) GetPurchaseOrderItems(purchaseOrderID uuid.UUID) ([]models.PurchaseOrderItem, error) {
	return s.supplierRepo.GetPurchaseOrderItems(purchaseOrderID)
}

func (s *supplierService) DeletePurchaseOrderItem(id uuid.UUID) error {
	return s.supplierRepo.DeletePurchaseOrderItem(id)
}

// Supplier Evaluation operations
func (s *supplierService) CreateSupplierEvaluation(req *CreateSupplierEvaluationRequest, userID uuid.UUID) (*models.SupplierEvaluation, error) {
	evaluation := &models.SupplierEvaluation{
		CompanyID:        userID, // This should be retrieved from user context
		SupplierID:       req.SupplierID,
		StartDate:        req.StartDate,
		EndDate:          req.EndDate,
		EvaluationType:   req.EvaluationType,
		QualityScore:     req.QualityScore,
		DeliveryScore:    req.DeliveryScore,
		ServiceScore:     req.ServiceScore,
		CostScore:        req.CostScore,
		TechnicalScore:   req.TechnicalScore,
		TotalOrders:      req.TotalOrders,
		OnTimeDeliveries: req.OnTimeDeliveries,
		QualityDefects:   req.QualityDefects,
		ServiceIssues:    req.ServiceIssues,
		CostSavings:      req.CostSavings,
		Strengths:        req.Strengths,
		Weaknesses:       req.Weaknesses,
		Recommendations:  req.Recommendations,
		ActionItems:      req.ActionItems,
		Status:           "draft",
		EvaluatedBy:      userID,
		EvaluatedAt:      req.EvaluatedAt,
	}

	// Calculate overall score
	evaluation.OverallScore = (req.QualityScore + req.DeliveryScore + req.ServiceScore + req.CostScore + req.TechnicalScore) / 5

	// Generate evaluation number
	timestamp := time.Now().Unix()
	evaluation.EvaluationNo = fmt.Sprintf("EVL%d", timestamp)

	if err := s.supplierRepo.CreateSupplierEvaluation(evaluation); err != nil {
		return nil, fmt.Errorf("failed to create supplier evaluation: %w", err)
	}

	return evaluation, nil
}

func (s *supplierService) UpdateSupplierEvaluation(id uuid.UUID, req *UpdateSupplierEvaluationRequest, userID uuid.UUID) (*models.SupplierEvaluation, error) {
	evaluation, err := s.supplierRepo.GetSupplierEvaluation(id)
	if err != nil {
		return nil, fmt.Errorf("supplier evaluation not found: %w", err)
	}

	// Update fields if provided
	updateSupplierEvaluationFields(evaluation, req)

	// Recalculate overall score if individual scores updated
	if req.QualityScore != nil || req.DeliveryScore != nil || req.ServiceScore != nil || req.CostScore != nil || req.TechnicalScore != nil {
		evaluation.OverallScore = (evaluation.QualityScore + evaluation.DeliveryScore + evaluation.ServiceScore + evaluation.CostScore + evaluation.TechnicalScore) / 5
	}

	if err := s.supplierRepo.UpdateSupplierEvaluation(evaluation); err != nil {
		return nil, fmt.Errorf("failed to update supplier evaluation: %w", err)
	}

	return evaluation, nil
}

func (s *supplierService) GetSupplierEvaluation(id uuid.UUID) (*models.SupplierEvaluation, error) {
	return s.supplierRepo.GetSupplierEvaluation(id)
}

func (s *supplierService) ListSupplierEvaluations(companyID uuid.UUID, params map[string]interface{}) ([]models.SupplierEvaluation, int64, error) {
	return s.supplierRepo.ListSupplierEvaluations(companyID, params)
}

func (s *supplierService) ApproveSupplierEvaluation(id uuid.UUID, userID uuid.UUID) error {
	evaluation, err := s.supplierRepo.GetSupplierEvaluation(id)
	if err != nil {
		return fmt.Errorf("supplier evaluation not found: %w", err)
	}

	if evaluation.Status != "completed" {
		return fmt.Errorf("supplier evaluation is not completed")
	}

	now := time.Now()
	evaluation.Status = "approved"
	evaluation.ApprovedBy = &userID
	evaluation.ApprovedAt = &now

	if err := s.supplierRepo.UpdateSupplierEvaluation(evaluation); err != nil {
		return fmt.Errorf("failed to approve supplier evaluation: %w", err)
	}

	// Update supplier performance metrics
	return s.UpdateSupplierPerformance(evaluation.SupplierID)
}

// Business operations
func (s *supplierService) UpdateSupplierPerformance(supplierID uuid.UUID) error {
	supplier, err := s.supplierRepo.GetSupplier(supplierID)
	if err != nil {
		return fmt.Errorf("supplier not found: %w", err)
	}

	// Get recent evaluations
	params := map[string]interface{}{
		"supplier_id": supplierID.String(),
		"status":      "approved",
		"page_size":   10,
	}
	
	evaluations, _, err := s.supplierRepo.ListSupplierEvaluations(supplier.CompanyID, params)
	if err != nil {
		return fmt.Errorf("failed to get supplier evaluations: %w", err)
	}

	if len(evaluations) == 0 {
		return nil // No evaluations to process
	}

	// Calculate average ratings
	var totalQuality, totalDelivery, totalService, totalOverall float64
	for _, eval := range evaluations {
		totalQuality += eval.QualityScore
		totalDelivery += eval.DeliveryScore
		totalService += eval.ServiceScore
		totalOverall += eval.OverallScore
	}

	count := float64(len(evaluations))
	supplier.QualityRating = totalQuality / count
	supplier.DeliveryRating = totalDelivery / count
	supplier.ServiceRating = totalService / count
	supplier.OverallRating = totalOverall / count

	// Update other performance metrics from purchase orders
	poParams := map[string]interface{}{
		"supplier_id": supplierID.String(),
		"status":      "received",
		"page_size":   100,
	}
	
	orders, _, err := s.supplierRepo.ListPurchaseOrders(supplier.CompanyID, poParams)
	if err == nil {
		supplier.TotalOrders = len(orders)
		
		onTimeCount := 0
		for _, order := range orders {
			if order.ActualEndDate != nil && order.RequiredDate.After(*order.ActualEndDate) {
				onTimeCount++
			}
		}
		supplier.OnTimeDeliveries = onTimeCount
	}

	return s.supplierRepo.UpdateSupplier(supplier)
}

func (s *supplierService) CalculateSupplierRisk(supplierID uuid.UUID) (string, error) {
	supplier, err := s.supplierRepo.GetSupplier(supplierID)
	if err != nil {
		return "", fmt.Errorf("supplier not found: %w", err)
	}

	riskScore := 0

	// Financial health risk
	switch supplier.FinancialHealth {
	case "poor":
		riskScore += 30
	case "fair":
		riskScore += 20
	case "good":
		riskScore += 10
	case "excellent":
		riskScore += 0
	}

	// Credit rating risk
	switch supplier.CreditRating {
	case "D", "C", "CC", "CCC":
		riskScore += 30
	case "B", "BB", "BBB":
		riskScore += 20
	case "A", "AA":
		riskScore += 10
	case "AAA":
		riskScore += 0
	}

	// Performance risk
	if supplier.QualityRating < 70 {
		riskScore += 20
	} else if supplier.QualityRating < 80 {
		riskScore += 10
	}

	if supplier.DeliveryRating < 70 {
		riskScore += 20
	} else if supplier.DeliveryRating < 80 {
		riskScore += 10
	}

	// Certification risk
	if !supplier.ISO9001 && !supplier.ISO14001 && !supplier.TS16949 {
		riskScore += 15
	}

	// Determine risk level
	var riskLevel string
	switch {
	case riskScore >= 70:
		riskLevel = "critical"
	case riskScore >= 50:
		riskLevel = "high"
	case riskScore >= 30:
		riskLevel = "medium"
	default:
		riskLevel = "low"
	}

	// Update supplier risk level
	supplier.RiskLevel = riskLevel
	if err := s.supplierRepo.UpdateSupplier(supplier); err != nil {
		return riskLevel, fmt.Errorf("failed to update supplier risk level: %w", err)
	}

	return riskLevel, nil
}

func (s *supplierService) GetSupplierDashboard(companyID uuid.UUID) (*SupplierDashboard, error) {
	dashboard := &SupplierDashboard{}

	// Get supplier counts
	suppliers, _, err := s.supplierRepo.ListSuppliers(companyID, map[string]interface{}{})
	if err != nil {
		return nil, fmt.Errorf("failed to get suppliers: %w", err)
	}

	dashboard.TotalSuppliers = len(suppliers)
	for _, supplier := range suppliers {
		switch supplier.Status {
		case "active":
			dashboard.ActiveSuppliers++
		case "suspended":
			dashboard.SuspendedSuppliers++
		case "blacklisted":
			dashboard.BlacklistedSuppliers++
		}

		switch supplier.RiskLevel {
		case "high":
			dashboard.HighRiskSuppliers++
		case "critical":
			dashboard.CriticalRiskSuppliers++
		}

		// Calculate average ratings
		if supplier.QualityRating > 0 {
			dashboard.AverageQualityRating += supplier.QualityRating
		}
		if supplier.DeliveryRating > 0 {
			dashboard.AverageDeliveryRating += supplier.DeliveryRating
		}
		if supplier.ServiceRating > 0 {
			dashboard.AverageServiceRating += supplier.ServiceRating
		}
	}

	if dashboard.TotalSuppliers > 0 {
		dashboard.AverageQualityRating /= float64(dashboard.TotalSuppliers)
		dashboard.AverageDeliveryRating /= float64(dashboard.TotalSuppliers)
		dashboard.AverageServiceRating /= float64(dashboard.TotalSuppliers)
	}

	// Get purchase order counts and values
	orders, _, err := s.supplierRepo.ListPurchaseOrders(companyID, map[string]interface{}{})
	if err != nil {
		return nil, fmt.Errorf("failed to get purchase orders: %w", err)
	}

	dashboard.TotalPurchaseOrders = len(orders)
	for _, order := range orders {
		dashboard.TotalPurchaseValue += order.TotalAmount

		switch order.Status {
		case "draft":
			dashboard.DraftOrders++
		case "sent":
			dashboard.SentOrders++
		case "confirmed":
			dashboard.ConfirmedOrders++
			dashboard.PendingValue += order.TotalAmount
		case "received":
			dashboard.ReceivedOrders++
			dashboard.ReceivedValue += order.TotalAmount
		case "cancelled":
			dashboard.CancelledOrders++
		}
	}

	// Get evaluation counts
	evaluations, _, err := s.supplierRepo.ListSupplierEvaluations(companyID, map[string]interface{}{})
	if err != nil {
		return nil, fmt.Errorf("failed to get supplier evaluations: %w", err)
	}

	dashboard.TotalEvaluations = len(evaluations)
	for _, eval := range evaluations {
		switch eval.Status {
		case "draft", "completed":
			dashboard.PendingEvaluations++
		case "approved":
			dashboard.CompletedEvaluations++
		}
	}

	return dashboard, nil
}

// Helper functions
func updateSupplierFields(supplier *models.Supplier, req *UpdateSupplierRequest) {
	if req.Name != nil {
		supplier.Name = *req.Name
	}
	if req.NameEn != nil {
		supplier.NameEn = *req.NameEn
	}
	if req.Type != nil {
		supplier.Type = *req.Type
	}
	if req.Status != nil {
		supplier.Status = *req.Status
	}
	if req.ContactPerson != nil {
		supplier.ContactPerson = *req.ContactPerson
	}
	if req.ContactTitle != nil {
		supplier.ContactTitle = *req.ContactTitle
	}
	if req.Phone != nil {
		supplier.Phone = *req.Phone
	}
	if req.Mobile != nil {
		supplier.Mobile = *req.Mobile
	}
	if req.Email != nil {
		supplier.Email = *req.Email
	}
	if req.Website != nil {
		supplier.Website = *req.Website
	}
	if req.Country != nil {
		supplier.Country = *req.Country
	}
	if req.State != nil {
		supplier.State = *req.State
	}
	if req.City != nil {
		supplier.City = *req.City
	}
	if req.Address != nil {
		supplier.Address = *req.Address
	}
	if req.PostalCode != nil {
		supplier.PostalCode = *req.PostalCode
	}
	if req.TaxNumber != nil {
		supplier.TaxNumber = *req.TaxNumber
	}
	if req.BusinessLicense != nil {
		supplier.BusinessLicense = *req.BusinessLicense
	}
	if req.Industry != nil {
		supplier.Industry = *req.Industry
	}
	if req.Established != nil {
		supplier.Established = req.Established
	}
	if req.Employees != nil {
		supplier.Employees = *req.Employees
	}
	if req.AnnualRevenue != nil {
		supplier.AnnualRevenue = *req.AnnualRevenue
	}
	if req.Currency != nil {
		supplier.Currency = *req.Currency
	}
	if req.PaymentTerms != nil {
		supplier.PaymentTerms = *req.PaymentTerms
	}
	if req.PaymentMethod != nil {
		supplier.PaymentMethod = *req.PaymentMethod
	}
	if req.CreditLimit != nil {
		supplier.CreditLimit = *req.CreditLimit
	}
	if req.CreditDays != nil {
		supplier.CreditDays = *req.CreditDays
	}
	if req.QualityRating != nil {
		supplier.QualityRating = *req.QualityRating
	}
	if req.DeliveryRating != nil {
		supplier.DeliveryRating = *req.DeliveryRating
	}
	if req.ServiceRating != nil {
		supplier.ServiceRating = *req.ServiceRating
	}
	if req.OverallRating != nil {
		supplier.OverallRating = *req.OverallRating
	}
	if req.ISO9001 != nil {
		supplier.ISO9001 = *req.ISO9001
	}
	if req.ISO14001 != nil {
		supplier.ISO14001 = *req.ISO14001
	}
	if req.TS16949 != nil {
		supplier.TS16949 = *req.TS16949
	}
	if req.OHSAS18001 != nil {
		supplier.OHSAS18001 = *req.OHSAS18001
	}
	if req.CustomCert != nil {
		supplier.CustomCert = *req.CustomCert
	}
	if req.CertExpiry != nil {
		supplier.CertExpiry = req.CertExpiry
	}
	if req.RiskLevel != nil {
		supplier.RiskLevel = *req.RiskLevel
	}
	if req.RiskFactors != nil {
		supplier.RiskFactors = *req.RiskFactors
	}
	if req.LastAuditDate != nil {
		supplier.LastAuditDate = req.LastAuditDate
	}
	if req.NextAuditDate != nil {
		supplier.NextAuditDate = req.NextAuditDate
	}
	if req.CreditRating != nil {
		supplier.CreditRating = *req.CreditRating
	}
	if req.FinancialHealth != nil {
		supplier.FinancialHealth = *req.FinancialHealth
	}
	if req.InsuranceCoverage != nil {
		supplier.InsuranceCoverage = *req.InsuranceCoverage
	}
	if req.Description != nil {
		supplier.Description = *req.Description
	}
	if req.Notes != nil {
		supplier.Notes = *req.Notes
	}
	if req.Tags != nil {
		supplier.Tags = *req.Tags
	}
}

func updateSupplierProductFields(product *models.SupplierProduct, req *UpdateSupplierProductRequest) {
	if req.InventoryID != nil {
		product.InventoryID = req.InventoryID
	}
	if req.ProductName != nil {
		product.ProductName = *req.ProductName
	}
	if req.ProductCode != nil {
		product.ProductCode = *req.ProductCode
	}
	if req.SupplierPartNo != nil {
		product.SupplierPartNo = *req.SupplierPartNo
	}
	if req.Category != nil {
		product.Category = *req.Category
	}
	if req.Specification != nil {
		product.Specification = *req.Specification
	}
	if req.Unit != nil {
		product.Unit = *req.Unit
	}
	if req.UnitPrice != nil {
		product.UnitPrice = *req.UnitPrice
	}
	if req.Currency != nil {
		product.Currency = *req.Currency
	}
	if req.MinOrderQty != nil {
		product.MinOrderQty = *req.MinOrderQty
	}
	if req.MaxOrderQty != nil {
		product.MaxOrderQty = *req.MaxOrderQty
	}
	if req.PriceBreaks != nil {
		product.PriceBreaks = *req.PriceBreaks
	}
	if req.LeadTimeDays != nil {
		product.LeadTimeDays = *req.LeadTimeDays
	}
	if req.QualityGrade != nil {
		product.QualityGrade = *req.QualityGrade
	}
	if req.CertificationReq != nil {
		product.CertificationReq = *req.CertificationReq
	}
	if req.Certificates != nil {
		product.Certificates = *req.Certificates
	}
	if req.Status != nil {
		product.Status = *req.Status
	}
	if req.IsPreferred != nil {
		product.IsPreferred = *req.IsPreferred
	}
	if req.LastPurchaseDate != nil {
		product.LastPurchaseDate = req.LastPurchaseDate
	}
	if req.LastPurchasePrice != nil {
		product.LastPurchasePrice = *req.LastPurchasePrice
	}
	if req.TotalPurchased != nil {
		product.TotalPurchased = *req.TotalPurchased
	}
	if req.QualityIssues != nil {
		product.QualityIssues = *req.QualityIssues
	}
	if req.DeliveryIssues != nil {
		product.DeliveryIssues = *req.DeliveryIssues
	}
}

func updatePurchaseOrderFields(order *models.PurchaseOrder, req *UpdatePurchaseOrderRequest) {
	if req.Status != nil {
		order.Status = *req.Status
	}
	if req.RequiredDate != nil {
		order.RequiredDate = *req.RequiredDate
	}
	if req.PromisedDate != nil {
		order.PromisedDate = req.PromisedDate
	}
	if req.PaymentTerms != nil {
		order.PaymentTerms = *req.PaymentTerms
	}
	if req.PaymentMethod != nil {
		order.PaymentMethod = *req.PaymentMethod
	}
	if req.ShippingAddress != nil {
		order.ShippingAddress = *req.ShippingAddress
	}
	if req.ShippingMethod != nil {
		order.ShippingMethod = *req.ShippingMethod
	}
	if req.TrackingNumber != nil {
		order.TrackingNumber = *req.TrackingNumber
	}
	if req.Currency != nil {
		order.Currency = *req.Currency
	}
	if req.ExchangeRate != nil {
		order.ExchangeRate = *req.ExchangeRate
	}
	if req.TaxRate != nil {
		order.TaxRate = *req.TaxRate
		order.TaxAmount = order.SubTotal * (*req.TaxRate) / 100
		order.TotalAmount = order.SubTotal + order.TaxAmount + order.ShippingCost
	}
	if req.ShippingCost != nil {
		order.ShippingCost = *req.ShippingCost
		order.TotalAmount = order.SubTotal + order.TaxAmount + (*req.ShippingCost)
	}
	if req.Notes != nil {
		order.Notes = *req.Notes
	}
	if req.InternalNotes != nil {
		order.InternalNotes = *req.InternalNotes
	}
}

func updatePurchaseOrderItemFields(item *models.PurchaseOrderItem, req *UpdatePurchaseOrderItemRequest) {
	if req.SupplierProductID != nil {
		item.SupplierProductID = req.SupplierProductID
	}
	if req.InventoryID != nil {
		item.InventoryID = req.InventoryID
	}
	if req.ProductName != nil {
		item.ProductName = *req.ProductName
	}
	if req.ProductCode != nil {
		item.ProductCode = *req.ProductCode
	}
	if req.SupplierPartNo != nil {
		item.SupplierPartNo = *req.SupplierPartNo
	}
	if req.Specification != nil {
		item.Specification = *req.Specification
	}
	if req.OrderedQuantity != nil {
		item.OrderedQuantity = *req.OrderedQuantity
		item.TotalPrice = item.OrderedQuantity * item.UnitPrice
	}
	if req.ReceivedQuantity != nil {
		item.ReceivedQuantity = *req.ReceivedQuantity
	}
	if req.Unit != nil {
		item.Unit = *req.Unit
	}
	if req.UnitPrice != nil {
		item.UnitPrice = *req.UnitPrice
		item.TotalPrice = item.OrderedQuantity * (*req.UnitPrice)
	}
	if req.Status != nil {
		item.Status = *req.Status
	}
	if req.QualityRequirement != nil {
		item.QualityRequirement = *req.QualityRequirement
	}
	if req.InspectionRequired != nil {
		item.InspectionRequired = *req.InspectionRequired
	}
}

func updateSupplierEvaluationFields(evaluation *models.SupplierEvaluation, req *UpdateSupplierEvaluationRequest) {
	if req.StartDate != nil {
		evaluation.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		evaluation.EndDate = *req.EndDate
	}
	if req.EvaluationType != nil {
		evaluation.EvaluationType = *req.EvaluationType
	}
	if req.QualityScore != nil {
		evaluation.QualityScore = *req.QualityScore
	}
	if req.DeliveryScore != nil {
		evaluation.DeliveryScore = *req.DeliveryScore
	}
	if req.ServiceScore != nil {
		evaluation.ServiceScore = *req.ServiceScore
	}
	if req.CostScore != nil {
		evaluation.CostScore = *req.CostScore
	}
	if req.TechnicalScore != nil {
		evaluation.TechnicalScore = *req.TechnicalScore
	}
	if req.TotalOrders != nil {
		evaluation.TotalOrders = *req.TotalOrders
	}
	if req.OnTimeDeliveries != nil {
		evaluation.OnTimeDeliveries = *req.OnTimeDeliveries
	}
	if req.QualityDefects != nil {
		evaluation.QualityDefects = *req.QualityDefects
	}
	if req.ServiceIssues != nil {
		evaluation.ServiceIssues = *req.ServiceIssues
	}
	if req.CostSavings != nil {
		evaluation.CostSavings = *req.CostSavings
	}
	if req.Strengths != nil {
		evaluation.Strengths = *req.Strengths
	}
	if req.Weaknesses != nil {
		evaluation.Weaknesses = *req.Weaknesses
	}
	if req.Recommendations != nil {
		evaluation.Recommendations = *req.Recommendations
	}
	if req.ActionItems != nil {
		evaluation.ActionItems = *req.ActionItems
	}
	if req.Status != nil {
		evaluation.Status = *req.Status
	}
}