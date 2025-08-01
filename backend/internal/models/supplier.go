package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Supplier represents a supplier/vendor
type Supplier struct {
	ID                uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID         uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	SupplierNo        string     `gorm:"not null;unique" json:"supplier_no"`
	Name              string     `gorm:"not null" json:"name"`
	NameEn            string     `json:"name_en"`
	Type              string     `gorm:"not null" json:"type"` // manufacturer, distributor, service_provider, raw_material
	Status            string     `gorm:"default:'active'" json:"status"` // active, inactive, suspended, blacklisted
	
	// Contact Information
	ContactPerson     string     `json:"contact_person"`
	ContactTitle      string     `json:"contact_title"`
	Phone             string     `json:"phone"`
	Mobile            string     `json:"mobile"`
	Email             string     `json:"email"`
	Website           string     `json:"website"`
	
	// Address
	Country           string     `json:"country"`
	State             string     `json:"state"`
	City              string     `json:"city"`
	Address           string     `json:"address"`
	PostalCode        string     `json:"postal_code"`
	
	// Business Information
	TaxNumber         string     `json:"tax_number"`
	BusinessLicense   string     `json:"business_license"`
	Industry          string     `json:"industry"`
	Established       *time.Time `json:"established"`
	Employees         int        `json:"employees"`
	AnnualRevenue     float64    `json:"annual_revenue"`
	Currency          string     `gorm:"default:'USD'" json:"currency"`
	
	// Payment Terms
	PaymentTerms      string     `json:"payment_terms"`
	PaymentMethod     string     `json:"payment_method"`
	CreditLimit       float64    `json:"credit_limit"`
	CreditDays        int        `json:"credit_days"`
	
	// Performance Metrics
	QualityRating     float64    `json:"quality_rating"`     // 0-100
	DeliveryRating    float64    `json:"delivery_rating"`    // 0-100
	ServiceRating     float64    `json:"service_rating"`     // 0-100
	OverallRating     float64    `json:"overall_rating"`     // 0-100
	TotalOrders       int        `json:"total_orders"`
	OnTimeDeliveries  int        `json:"on_time_deliveries"`
	DefectiveItems    int        `json:"defective_items"`
	
	// Certifications
	ISO9001           bool       `json:"iso_9001"`
	ISO14001          bool       `json:"iso_14001"`
	TS16949           bool       `json:"ts_16949"`
	OHSAS18001        bool       `json:"ohsas_18001"`
	CustomCert        string     `json:"custom_cert"`
	CertExpiry        *time.Time `json:"cert_expiry"`
	
	// Risk Assessment
	RiskLevel         string     `json:"risk_level"`         // low, medium, high, critical
	RiskFactors       string     `json:"risk_factors"`       // JSON array
	LastAuditDate     *time.Time `json:"last_audit_date"`
	NextAuditDate     *time.Time `json:"next_audit_date"`
	
	// Financial Information
	CreditRating      string     `json:"credit_rating"`      // AAA, AA, A, BBB, BB, B, CCC, CC, C, D
	FinancialHealth   string     `json:"financial_health"`   // excellent, good, fair, poor
	InsuranceCoverage float64    `json:"insurance_coverage"`
	
	// Additional Information
	Description       string     `json:"description"`
	Notes             string     `json:"notes"`
	Tags              string     `json:"tags"`               // JSON array
	
	// Timestamps
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	CreatedBy         uuid.UUID  `gorm:"type:uuid;not null" json:"created_by"`
	
	// Relations
	Company           *Company           `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Creator           *User              `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Contacts          []SupplierContact  `gorm:"foreignKey:SupplierID" json:"contacts,omitempty"`
	Products          []SupplierProduct  `gorm:"foreignKey:SupplierID" json:"products,omitempty"`
}

func (s *Supplier) BeforeCreate(tx *gorm.DB) error {
	s.ID = uuid.New()
	return nil
}

// SupplierContact represents additional contacts for a supplier
type SupplierContact struct {
	ID                uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	SupplierID        uuid.UUID  `gorm:"type:uuid;not null" json:"supplier_id"`
	Name              string     `gorm:"not null" json:"name"`
	Title             string     `json:"title"`
	Department        string     `json:"department"`
	Phone             string     `json:"phone"`
	Mobile            string     `json:"mobile"`
	Email             string     `json:"email"`
	IsPrimary         bool       `gorm:"default:false" json:"is_primary"`
	IsActive          bool       `gorm:"default:true" json:"is_active"`
	
	// Specialties
	Responsibilities  string     `json:"responsibilities"`   // JSON array
	Languages         string     `json:"languages"`          // JSON array
	
	// Timestamps
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	
	// Relations
	Supplier          *Supplier  `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
}

func (sc *SupplierContact) BeforeCreate(tx *gorm.DB) error {
	sc.ID = uuid.New()
	return nil
}

// SupplierProduct represents products/services offered by a supplier
type SupplierProduct struct {
	ID                uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	SupplierID        uuid.UUID  `gorm:"type:uuid;not null" json:"supplier_id"`
	InventoryID       *uuid.UUID `gorm:"type:uuid" json:"inventory_id"`
	
	// Product Information
	ProductName       string     `gorm:"not null" json:"product_name"`
	ProductCode       string     `json:"product_code"`
	SupplierPartNo    string     `json:"supplier_part_no"`
	Category          string     `json:"category"`
	Specification     string     `json:"specification"`
	Unit              string     `gorm:"not null" json:"unit"`
	
	// Pricing
	UnitPrice         float64    `json:"unit_price"`
	Currency          string     `gorm:"default:'USD'" json:"currency"`
	MinOrderQty       float64    `json:"min_order_qty"`
	MaxOrderQty       float64    `json:"max_order_qty"`
	PriceBreaks       string     `json:"price_breaks"`       // JSON array of quantity breaks
	
	// Lead Time
	LeadTimeDays      int        `json:"lead_time_days"`
	
	// Quality
	QualityGrade      string     `json:"quality_grade"`
	CertificationReq  bool       `json:"certification_req"`
	Certificates      string     `json:"certificates"`       // JSON array
	
	// Status
	Status            string     `gorm:"default:'active'" json:"status"` // active, inactive, discontinued
	IsPreferred       bool       `gorm:"default:false" json:"is_preferred"`
	LastPurchaseDate  *time.Time `json:"last_purchase_date"`
	LastPurchasePrice float64    `json:"last_purchase_price"`
	
	// Performance
	TotalPurchased    float64    `json:"total_purchased"`
	QualityIssues     int        `json:"quality_issues"`
	DeliveryIssues    int        `json:"delivery_issues"`
	
	// Timestamps
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	
	// Relations
	Supplier          *Supplier  `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	Inventory         *Inventory `gorm:"foreignKey:InventoryID" json:"inventory,omitempty"`
}

func (sp *SupplierProduct) BeforeCreate(tx *gorm.DB) error {
	sp.ID = uuid.New()
	return nil
}

// PurchaseOrder represents a purchase order
type PurchaseOrder struct {
	ID                uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID         uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	OrderNo           string     `gorm:"not null;unique" json:"order_no"`
	Status            string     `gorm:"not null" json:"status"` // draft, sent, confirmed, partial_received, received, cancelled
	
	// Supplier Information
	SupplierID        uuid.UUID  `gorm:"type:uuid;not null" json:"supplier_id"`
	
	// Order Details
	OrderDate         time.Time  `json:"order_date"`
	RequiredDate      time.Time  `json:"required_date"`
	PromisedDate      *time.Time `json:"promised_date"`
	
	// Financial
	SubTotal          float64    `json:"sub_total"`
	TaxRate           float64    `json:"tax_rate"`
	TaxAmount         float64    `json:"tax_amount"`
	ShippingCost      float64    `json:"shipping_cost"`
	TotalAmount       float64    `json:"total_amount"`
	Currency          string     `gorm:"default:'USD'" json:"currency"`
	ExchangeRate      float64    `gorm:"default:1" json:"exchange_rate"`
	
	// Payment
	PaymentTerms      string     `json:"payment_terms"`
	PaymentMethod     string     `json:"payment_method"`
	
	// Shipping
	ShippingAddress   string     `json:"shipping_address"`
	ShippingMethod    string     `json:"shipping_method"`
	TrackingNumber    string     `json:"tracking_number"`
	
	// Additional Information
	Notes             string     `json:"notes"`
	InternalNotes     string     `json:"internal_notes"`
	
	// Timestamps
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	CreatedBy         uuid.UUID  `gorm:"type:uuid;not null" json:"created_by"`
	ApprovedBy        *uuid.UUID `gorm:"type:uuid" json:"approved_by"`
	ApprovedAt        *time.Time `json:"approved_at"`
	
	// Relations
	Company           *Company              `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Supplier          *Supplier             `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	Creator           *User                 `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Approver          *User                 `gorm:"foreignKey:ApprovedBy" json:"approver,omitempty"`
	Items             []PurchaseOrderItem   `gorm:"foreignKey:PurchaseOrderID" json:"items,omitempty"`
}

func (po *PurchaseOrder) BeforeCreate(tx *gorm.DB) error {
	po.ID = uuid.New()
	return nil
}

// PurchaseOrderItem represents an item in a purchase order
type PurchaseOrderItem struct {
	ID                uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	PurchaseOrderID   uuid.UUID  `gorm:"type:uuid;not null" json:"purchase_order_id"`
	SupplierProductID *uuid.UUID `gorm:"type:uuid" json:"supplier_product_id"`
	InventoryID       *uuid.UUID `gorm:"type:uuid" json:"inventory_id"`
	
	// Item Details
	ProductName       string     `gorm:"not null" json:"product_name"`
	ProductCode       string     `json:"product_code"`
	SupplierPartNo    string     `json:"supplier_part_no"`
	Specification     string     `json:"specification"`
	
	// Quantity & Pricing
	OrderedQuantity   float64    `gorm:"not null" json:"ordered_quantity"`
	ReceivedQuantity  float64    `json:"received_quantity"`
	Unit              string     `gorm:"not null" json:"unit"`
	UnitPrice         float64    `gorm:"not null" json:"unit_price"`
	TotalPrice        float64    `gorm:"not null" json:"total_price"`
	
	// Status
	Status            string     `gorm:"default:'pending'" json:"status"` // pending, partial_received, received, cancelled
	
	// Quality
	QualityRequirement string    `json:"quality_requirement"`
	InspectionRequired bool      `gorm:"default:false" json:"inspection_required"`
	
	// Timestamps
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	
	// Relations
	PurchaseOrder     *PurchaseOrder    `gorm:"foreignKey:PurchaseOrderID" json:"purchase_order,omitempty"`
	SupplierProduct   *SupplierProduct  `gorm:"foreignKey:SupplierProductID" json:"supplier_product,omitempty"`
	Inventory         *Inventory        `gorm:"foreignKey:InventoryID" json:"inventory,omitempty"`
}

func (poi *PurchaseOrderItem) BeforeCreate(tx *gorm.DB) error {
	poi.ID = uuid.New()
	return nil
}

// SupplierEvaluation represents supplier performance evaluation
type SupplierEvaluation struct {
	ID                uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID         uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	SupplierID        uuid.UUID  `gorm:"type:uuid;not null" json:"supplier_id"`
	EvaluationNo      string     `gorm:"not null;unique" json:"evaluation_no"`
	
	// Evaluation Period
	StartDate         time.Time  `json:"start_date"`
	EndDate           time.Time  `json:"end_date"`
	EvaluationType    string     `json:"evaluation_type"`    // monthly, quarterly, annual, ad_hoc
	
	// Scores (0-100)
	QualityScore      float64    `json:"quality_score"`
	DeliveryScore     float64    `json:"delivery_score"`
	ServiceScore      float64    `json:"service_score"`
	CostScore         float64    `json:"cost_score"`
	TechnicalScore    float64    `json:"technical_score"`
	OverallScore      float64    `json:"overall_score"`
	
	// Metrics
	TotalOrders       int        `json:"total_orders"`
	OnTimeDeliveries  int        `json:"on_time_deliveries"`
	QualityDefects    int        `json:"quality_defects"`
	ServiceIssues     int        `json:"service_issues"`
	CostSavings       float64    `json:"cost_savings"`
	
	// Evaluation Details
	Strengths         string     `json:"strengths"`
	Weaknesses        string     `json:"weaknesses"`
	Recommendations   string     `json:"recommendations"`
	ActionItems       string     `json:"action_items"`       // JSON array
	
	// Status
	Status            string     `gorm:"default:'draft'" json:"status"` // draft, completed, approved
	
	// Evaluator
	EvaluatedBy       uuid.UUID  `gorm:"type:uuid;not null" json:"evaluated_by"`
	EvaluatedAt       time.Time  `json:"evaluated_at"`
	ApprovedBy        *uuid.UUID `gorm:"type:uuid" json:"approved_by"`
	ApprovedAt        *time.Time `json:"approved_at"`
	
	// Timestamps
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	
	// Relations
	Company           *Company   `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Supplier          *Supplier  `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	Evaluator         *User      `gorm:"foreignKey:EvaluatedBy" json:"evaluator,omitempty"`
	Approver          *User      `gorm:"foreignKey:ApprovedBy" json:"approver,omitempty"`
}

func (se *SupplierEvaluation) BeforeCreate(tx *gorm.DB) error {
	se.ID = uuid.New()
	return nil
}