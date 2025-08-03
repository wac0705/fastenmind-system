package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TariffCode 關稅代碼
type TariffCode struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID       uuid.UUID `gorm:"type:uuid;not null" json:"company_id"`
	HSCode          string    `gorm:"not null;index" json:"hs_code"`           // Harmonized System Code
	Description     string    `gorm:"not null" json:"description"`
	DescriptionEN   string    `json:"description_en"`
	Category        string    `json:"category"`
	Unit            string    `json:"unit"`                                    // kg, piece, m2, etc.
	BaseRate        float64   `json:"base_rate"`                               // Base tariff rate
	PreferentialRate float64  `json:"preferential_rate"`                       // Preferential rate
	VAT             float64   `json:"vat"`                                     // VAT rate
	ExciseTax       float64   `json:"excise_tax"`                              // Excise tax rate
	ImportRestriction string  `json:"import_restriction"`                      // JSON restriction details
	ExportRestriction string  `json:"export_restriction"`                      // JSON restriction details
	RequiredCerts   string    `json:"required_certs"`                          // JSON required certificates
	IsActive        bool      `gorm:"default:true" json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	CreatedBy       uuid.UUID `gorm:"type:uuid" json:"created_by"`

	// Relations
	Company *Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Creator *User    `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

// TariffRate 關稅稅率
type TariffRate struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID     uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	TariffCodeID  uuid.UUID  `gorm:"type:uuid;not null" json:"tariff_code_id"`
	CountryCode   string     `gorm:"not null" json:"country_code"`            // ISO country code
	CountryName   string     `gorm:"not null" json:"country_name"`
	Rate          float64    `gorm:"not null" json:"rate"`                    // Tariff rate percentage
	RateType      string     `gorm:"not null" json:"rate_type"`               // ad_valorem, specific, compound
	MinimumDuty   float64    `json:"minimum_duty"`
	MaximumDuty   float64    `json:"maximum_duty"`
	Currency      string     `gorm:"not null;default:'USD'" json:"currency"`
	TradeType     string     `gorm:"not null" json:"trade_type"`              // import, export
	AgreementType string     `json:"agreement_type"`                          // mfn, fta, gsp, etc.
	ValidFrom     time.Time  `gorm:"not null" json:"valid_from"`
	ValidTo       *time.Time `json:"valid_to"`
	IsActive      bool       `gorm:"default:true" json:"is_active"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	CreatedBy     uuid.UUID  `gorm:"type:uuid" json:"created_by"`

	// Relations
	Company     *Company    `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	TariffCode  *TariffCode `gorm:"foreignKey:TariffCodeID" json:"tariff_code,omitempty"`
	Creator     *User       `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

// TradeDocument 貿易文件
type TradeDocument struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID     uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	DocumentType  string     `gorm:"not null" json:"document_type"`           // invoice, packing_list, bl, co, etc.
	DocumentNo    string     `gorm:"not null" json:"document_no"`
	Title         string     `gorm:"not null" json:"title"`
	Description   string     `json:"description"`
	FilePath      string     `json:"file_path"`
	FileSize      int64      `json:"file_size"`
	FileType      string     `json:"file_type"`
	Version       int        `gorm:"default:1" json:"version"`
	Status        string     `gorm:"not null" json:"status"`                  // draft, submitted, approved, rejected, expired
	IsRequired    bool       `gorm:"default:false" json:"is_required"`
	ValidFrom     *time.Time `json:"valid_from"`
	ValidTo       *time.Time `json:"valid_to"`
	IssuedBy      string     `json:"issued_by"`
	IssuedAt      *time.Time `json:"issued_at"`
	ApprovedBy    *uuid.UUID `gorm:"type:uuid" json:"approved_by"`
	ApprovedAt    *time.Time `json:"approved_at"`
	RejectedReason string    `json:"rejected_reason"`
	Metadata      string     `json:"metadata"`                                // JSON metadata
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	CreatedBy     uuid.UUID  `gorm:"type:uuid" json:"created_by"`

	// Relations
	Company   *Company  `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Creator   *User     `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Approver  *User     `gorm:"foreignKey:ApprovedBy" json:"approver,omitempty"`
	Shipments []Shipment `gorm:"many2many:shipment_documents;" json:"shipments,omitempty"`
}

// Shipment 運輸管理
type Shipment struct {
	ID                uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID         uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	ShipmentNo        string     `gorm:"not null;unique" json:"shipment_no"`
	OrderID           *uuid.UUID `gorm:"type:uuid" json:"order_id"`
	Type              string     `gorm:"not null" json:"type"`                 // import, export
	Status            string     `gorm:"not null" json:"status"`               // pending, in_transit, customs, delivered, cancelled
	Method            string     `gorm:"not null" json:"method"`               // sea, air, land, express
	CarrierName       string     `json:"carrier_name"`
	TrackingNo        string     `json:"tracking_no"`
	ContainerNo       string     `json:"container_no"`
	ContainerType     string     `json:"container_type"`                       // 20ft, 40ft, 40hc, etc.
	OriginCountry     string     `gorm:"not null" json:"origin_country"`
	OriginPort        string     `json:"origin_port"`
	OriginAddress     string     `json:"origin_address"`
	DestCountry       string     `gorm:"not null" json:"dest_country"`
	DestPort          string     `json:"dest_port"`
	DestAddress       string     `json:"dest_address"`
	EstimatedDeparture *time.Time `json:"estimated_departure"`
	ActualDeparture   *time.Time `json:"actual_departure"`
	EstimatedArrival  *time.Time `json:"estimated_arrival"`
	ActualArrival     *time.Time `json:"actual_arrival"`
	GrossWeight       float64    `json:"gross_weight"`                         // in kg
	NetWeight         float64    `json:"net_weight"`                           // in kg
	Volume            float64    `json:"volume"`                               // in m3
	PackageCount      int        `json:"package_count"`
	PackageType       string     `json:"package_type"`                         // carton, pallet, etc.
	InsuranceValue    float64    `json:"insurance_value"`
	InsuranceCurrency string     `json:"insurance_currency"`
	FreightCost       float64    `json:"freight_cost"`
	FreightCurrency   string     `json:"freight_currency"`
	CustomsValue      float64    `json:"customs_value"`
	CustomsCurrency   string     `json:"customs_currency"`
	TotalDuty         float64    `json:"total_duty"`
	TotalTax          float64    `json:"total_tax"`
	SpecialInstructions string   `json:"special_instructions"`
	InternalNotes     string     `json:"internal_notes"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	CreatedBy         uuid.UUID  `gorm:"type:uuid" json:"created_by"`

	// Relations
	Company   *Company        `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Order     *Order          `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Creator   *User           `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Documents []TradeDocument `gorm:"many2many:shipment_documents;" json:"documents,omitempty"`
	Items     []ShipmentItem  `gorm:"foreignKey:ShipmentID" json:"items,omitempty"`
	Events    []ShipmentEvent `gorm:"foreignKey:ShipmentID" json:"events,omitempty"`
}

// ShipmentItem 運輸項目
type ShipmentItem struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID     uuid.UUID `gorm:"type:uuid;not null" json:"company_id"`
	ShipmentID    uuid.UUID `gorm:"type:uuid;not null" json:"shipment_id"`
	ProductID     *uuid.UUID `gorm:"type:uuid" json:"product_id"`
	HSCode        string    `json:"hs_code"`
	ProductName   string    `gorm:"not null" json:"product_name"`
	Description   string    `json:"description"`
	Quantity      float64   `gorm:"not null" json:"quantity"`
	Unit          string    `gorm:"not null" json:"unit"`
	UnitWeight    float64   `json:"unit_weight"`                              // in kg
	UnitValue     float64   `json:"unit_value"`
	Currency      string    `gorm:"default:'USD'" json:"currency"`
	TotalWeight   float64   `json:"total_weight"`                             // quantity * unit_weight
	TotalValue    float64   `json:"total_value"`                              // quantity * unit_value
	CountryOrigin string    `json:"country_origin"`
	Manufacturer  string    `json:"manufacturer"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Relations
	Company  *Company  `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Shipment *Shipment `gorm:"foreignKey:ShipmentID" json:"shipment,omitempty"`
	Product  *Product  `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

// ShipmentEvent 運輸事件
type ShipmentEvent struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID   uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	ShipmentID  uuid.UUID  `gorm:"type:uuid;not null" json:"shipment_id"`
	EventType   string     `gorm:"not null" json:"event_type"`               // departure, arrival, customs_clearance, delivery, delay, etc.
	Status      string     `gorm:"not null" json:"status"`                   // completed, in_progress, pending, cancelled
	Location    string     `json:"location"`
	Description string     `json:"description"`
	Longitude   float64    `json:"longitude"`
	Latitude    float64    `json:"latitude"`
	EventTime   time.Time  `gorm:"not null" json:"event_time"`
	RecordedAt  time.Time  `gorm:"not null" json:"recorded_at"`
	Source      string     `json:"source"`                                   // manual, api, tracking, etc.
	CreatedAt   time.Time  `json:"created_at"`
	CreatedBy   uuid.UUID  `gorm:"type:uuid" json:"created_by"`

	// Relations
	Company  *Company  `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Shipment *Shipment `gorm:"foreignKey:ShipmentID" json:"shipment,omitempty"`
	Creator  *User     `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

// LetterOfCredit 信用狀
type LetterOfCredit struct {
	ID                uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID         uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	LCNumber          string     `gorm:"not null;unique" json:"lc_number"`
	Type              string     `gorm:"not null" json:"type"`                 // sight, usance, revolving, standby
	Status            string     `gorm:"not null" json:"status"`               // draft, issued, advised, confirmed, utilized, expired
	Amount            float64    `gorm:"not null" json:"amount"`
	Currency          string     `gorm:"not null" json:"currency"`
	ApplicantName     string     `gorm:"not null" json:"applicant_name"`       // Buyer
	ApplicantAddress  string     `json:"applicant_address"`
	BeneficiaryName   string     `gorm:"not null" json:"beneficiary_name"`     // Seller
	BeneficiaryAddress string    `json:"beneficiary_address"`
	IssuingBank       string     `gorm:"not null" json:"issuing_bank"`
	AdvisingBank      string     `json:"advising_bank"`
	ConfirmingBank    string     `json:"confirming_bank"`
	IssueDate         time.Time  `gorm:"not null" json:"issue_date"`
	ExpiryDate        time.Time  `gorm:"not null" json:"expiry_date"`
	LastShipmentDate  *time.Time `json:"last_shipment_date"`
	PartialShipment   bool       `gorm:"default:false" json:"partial_shipment"`
	Transhipment      bool       `gorm:"default:false" json:"transhipment"`
	PortOfLoading     string     `json:"port_of_loading"`
	PortOfDischarge   string     `json:"port_of_discharge"`
	Description       string     `json:"description"`                          // Goods description
	Documents         string     `json:"documents"`                            // JSON required documents
	Terms             string     `json:"terms"`                                // JSON terms and conditions
	UtilizedAmount    float64    `gorm:"default:0" json:"utilized_amount"`
	AvailableAmount   float64    `json:"available_amount"`                     // amount - utilized_amount
	Amendments        string     `json:"amendments"`                           // JSON amendments history
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	CreatedBy         uuid.UUID  `gorm:"type:uuid" json:"created_by"`

	// Relations
	Company    *Company          `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Creator    *User             `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Shipments  []Shipment        `gorm:"foreignKey:LCID" json:"shipments,omitempty"`
	Utilizations []LCUtilization `gorm:"foreignKey:LCID" json:"utilizations,omitempty"`
}

// LCUtilization 信用狀使用記錄
type LCUtilization struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID    uuid.UUID `gorm:"type:uuid;not null" json:"company_id"`
	LCID         uuid.UUID `gorm:"type:uuid;not null" json:"lc_id"`
	ShipmentID   *uuid.UUID `gorm:"type:uuid" json:"shipment_id"`
	Amount       float64   `gorm:"not null" json:"amount"`
	Currency     string    `gorm:"not null" json:"currency"`
	Description  string    `json:"description"`
	DocumentsRef string    `json:"documents_ref"`                            // JSON document references
	Status       string    `gorm:"not null" json:"status"`                   // pending, accepted, rejected
	UtilizedAt   time.Time `gorm:"not null" json:"utilized_at"`
	CreatedAt    time.Time `json:"created_at"`
	CreatedBy    uuid.UUID `gorm:"type:uuid" json:"created_by"`

	// Relations
	Company *Company       `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	LC      *LetterOfCredit `gorm:"foreignKey:LCID" json:"lc,omitempty"`
	Shipment *Shipment     `gorm:"foreignKey:ShipmentID" json:"shipment,omitempty"`
	Creator *User          `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

// TradeCompliance 貿易合規
type TradeCompliance struct {
	ID             uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID      uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	ComplianceType string     `gorm:"not null" json:"compliance_type"`        // export_control, sanctions, anti_dumping, etc.
	RuleName       string     `gorm:"not null" json:"rule_name"`
	Description    string     `json:"description"`
	CountryCode    string     `json:"country_code"`
	ProductCodes   string     `json:"product_codes"`                          // JSON array of HS codes
	EntityList     string     `json:"entity_list"`                            // JSON list of restricted entities
	RuleDetails    string     `json:"rule_details"`                           // JSON rule configuration
	Severity       string     `gorm:"not null" json:"severity"`               // low, medium, high, critical
	Status         string     `gorm:"not null" json:"status"`                 // active, suspended, expired
	ValidFrom      time.Time  `gorm:"not null" json:"valid_from"`
	ValidTo        *time.Time `json:"valid_to"`
	Source         string     `json:"source"`                                 // regulatory_body, internal_policy
	LastUpdated    time.Time  `json:"last_updated"`
	UpdatedBy      *uuid.UUID `gorm:"type:uuid" json:"updated_by"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	CreatedBy      uuid.UUID  `gorm:"type:uuid" json:"created_by"`

	// Relations
	Company  *Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Creator  *User    `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Updater  *User    `gorm:"foreignKey:UpdatedBy" json:"updater,omitempty"`
	Checks   []ComplianceCheck `gorm:"foreignKey:ComplianceID" json:"checks,omitempty"`
}

// ComplianceCheck 合規檢查
type ComplianceCheck struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID    uuid.UUID `gorm:"type:uuid;not null" json:"company_id"`
	ComplianceID uuid.UUID `gorm:"type:uuid;not null" json:"compliance_id"`
	ResourceType string    `gorm:"not null" json:"resource_type"`            // shipment, order, customer, product
	ResourceID   uuid.UUID `gorm:"type:uuid;not null" json:"resource_id"`
	CheckType    string    `gorm:"not null" json:"check_type"`               // automatic, manual
	Result       string    `gorm:"not null" json:"result"`                   // passed, failed, warning, pending
	Score        float64   `json:"score"`                                    // Compliance score 0-100
	Issues       string    `json:"issues"`                                   // JSON array of issues found
	Recommendations string `json:"recommendations"`                          // JSON recommendations
	CheckedAt    time.Time `gorm:"not null" json:"checked_at"`
	CheckedBy    *uuid.UUID `gorm:"type:uuid" json:"checked_by"`
	ResolvedAt   *time.Time `json:"resolved_at"`
	ResolvedBy   *uuid.UUID `gorm:"type:uuid" json:"resolved_by"`
	Notes        string    `json:"notes"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relations
	Company    *Company         `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Compliance *TradeCompliance `gorm:"foreignKey:ComplianceID" json:"compliance,omitempty"`
	Checker    *User            `gorm:"foreignKey:CheckedBy" json:"checker,omitempty"`
	Resolver   *User            `gorm:"foreignKey:ResolvedBy" json:"resolver,omitempty"`
}

// TradeExchangeRate 貿易匯率管理
type TradeExchangeRate struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID    uuid.UUID `gorm:"type:uuid;not null" json:"company_id"`
	FromCurrency string    `gorm:"not null" json:"from_currency"`
	ToCurrency   string    `gorm:"not null" json:"to_currency"`
	Rate         float64   `gorm:"not null" json:"rate"`
	RateType     string    `gorm:"not null" json:"rate_type"`               // buy, sell, mid, official
	Source       string    `gorm:"not null" json:"source"`                 // bank, api, manual, central_bank
	ValidDate    time.Time `gorm:"not null" json:"valid_date"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedBy    uuid.UUID `gorm:"type:uuid" json:"created_by"`

	// Relations
	Company *Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Creator *User    `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

// TradeRegulation 貿易法規
type TradeRegulation struct {
	ID              uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID       *uuid.UUID `gorm:"type:uuid" json:"company_id"`          // NULL for global regulations
	RegulationType  string     `gorm:"not null" json:"regulation_type"`      // import, export, transit, customs
	CountryCode     string     `gorm:"not null" json:"country_code"`
	RegulationCode  string     `gorm:"not null" json:"regulation_code"`
	Title           string     `gorm:"not null" json:"title"`
	Description     string     `json:"description"`
	Requirements    string     `json:"requirements"`                         // JSON requirements
	Penalties       string     `json:"penalties"`                            // JSON penalty information
	DocumentsNeeded string     `json:"documents_needed"`                     // JSON required documents
	ProcessingTime  int        `json:"processing_time"`                      // days
	Fees            string     `json:"fees"`                                 // JSON fee structure
	ApplicableHS    string     `json:"applicable_hs"`                        // JSON HS code patterns
	EffectiveDate   time.Time  `gorm:"not null" json:"effective_date"`
	ExpiryDate      *time.Time `json:"expiry_date"`
	Status          string     `gorm:"not null" json:"status"`               // active, suspended, repealed
	OfficialURL     string     `json:"official_url"`
	LastReviewDate  *time.Time `json:"last_review_date"`
	ReviewedBy      *uuid.UUID `gorm:"type:uuid" json:"reviewed_by"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	CreatedBy       uuid.UUID  `gorm:"type:uuid" json:"created_by"`

	// Relations
	Company  *Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Creator  *User    `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Reviewer *User    `gorm:"foreignKey:ReviewedBy" json:"reviewer,omitempty"`
}

// TradeAgreement 貿易協定
type TradeAgreement struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID     *uuid.UUID `gorm:"type:uuid" json:"company_id"`           // NULL for multilateral agreements
	AgreementType string     `gorm:"not null" json:"agreement_type"`        // fta, customs_union, trade_pact, bilateral
	AgreementCode string     `gorm:"not null;unique" json:"agreement_code"`
	Name          string     `gorm:"not null" json:"name"`
	Description   string     `json:"description"`
	Countries     string     `gorm:"not null" json:"countries"`             // JSON array of country codes
	Benefits      string     `json:"benefits"`                              // JSON benefits description
	TariffReductions string  `json:"tariff_reductions"`                     // JSON tariff reduction details
	QuotaLimits   string     `json:"quota_limits"`                          // JSON quota information
	RulesOfOrigin string     `json:"rules_of_origin"`                       // JSON origin rules
	EffectiveDate time.Time  `gorm:"not null" json:"effective_date"`
	ExpiryDate    *time.Time `json:"expiry_date"`
	Status        string     `gorm:"not null" json:"status"`                // active, suspended, expired, negotiating
	OfficialURL   string     `json:"official_url"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	CreatedBy     uuid.UUID  `gorm:"type:uuid" json:"created_by"`

	// Relations
	Company *Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Creator *User    `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

// BeforeCreate hooks
func (tc *TariffCode) BeforeCreate(tx *gorm.DB) error {
	if tc.ID == uuid.Nil {
		tc.ID = uuid.New()
	}
	return nil
}

func (tr *TariffRate) BeforeCreate(tx *gorm.DB) error {
	if tr.ID == uuid.Nil {
		tr.ID = uuid.New()
	}
	return nil
}

func (td *TradeDocument) BeforeCreate(tx *gorm.DB) error {
	if td.ID == uuid.Nil {
		td.ID = uuid.New()
	}
	return nil
}

func (s *Shipment) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

func (si *ShipmentItem) BeforeCreate(tx *gorm.DB) error {
	if si.ID == uuid.Nil {
		si.ID = uuid.New()
	}
	return nil
}

func (se *ShipmentEvent) BeforeCreate(tx *gorm.DB) error {
	if se.ID == uuid.Nil {
		se.ID = uuid.New()
	}
	return nil
}

func (lc *LetterOfCredit) BeforeCreate(tx *gorm.DB) error {
	if lc.ID == uuid.Nil {
		lc.ID = uuid.New()
	}
	return nil
}

func (lcu *LCUtilization) BeforeCreate(tx *gorm.DB) error {
	if lcu.ID == uuid.Nil {
		lcu.ID = uuid.New()
	}
	return nil
}

func (tc *TradeCompliance) BeforeCreate(tx *gorm.DB) error {
	if tc.ID == uuid.Nil {
		tc.ID = uuid.New()
	}
	return nil
}

func (cc *ComplianceCheck) BeforeCreate(tx *gorm.DB) error {
	if cc.ID == uuid.Nil {
		cc.ID = uuid.New()
	}
	return nil
}

func (er *TradeExchangeRate) BeforeCreate(tx *gorm.DB) error {
	if er.ID == uuid.Nil {
		er.ID = uuid.New()
	}
	return nil
}

func (tr *TradeRegulation) BeforeCreate(tx *gorm.DB) error {
	if tr.ID == uuid.Nil {
		tr.ID = uuid.New()
	}
	return nil
}

func (ta *TradeAgreement) BeforeCreate(tx *gorm.DB) error {
	if ta.ID == uuid.Nil {
		ta.ID = uuid.New()
	}
	return nil
}