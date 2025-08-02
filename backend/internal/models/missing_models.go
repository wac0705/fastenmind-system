package models

import (
	"time"
	"github.com/google/uuid"
)

// Compliance Models
type ComplianceRule struct {
	ID            uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	ProductType   string     `json:"product_type"`
	ExportCountry string     `json:"export_country"`
	ImportCountry string     `json:"import_country"`
	IsActive      bool       `json:"is_active"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type DocumentRequirement struct {
	ID            uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	DocumentType  string     `json:"document_type"`
	ProductType   string     `json:"product_type"`
	ExportCountry string     `json:"export_country"`
	ImportCountry string     `json:"import_country"`
	IsRequired    bool       `json:"is_required"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type ComplianceCheckResult struct {
	ID            uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	InquiryID     uuid.UUID  `json:"inquiry_id"`
	RuleID        uuid.UUID  `json:"rule_id"`
	Status        string     `json:"status"`
	Message       string     `json:"message"`
	CheckedAt     time.Time  `json:"checked_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// N8N Models
type N8NTrigger struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	WorkflowID      uuid.UUID  `json:"workflow_id"`
	TriggerType     string     `json:"trigger_type"`
	EventType       string     `json:"event_type"`
	Configuration   JSONB      `json:"configuration" gorm:"type:jsonb"`
	IsActive        bool       `json:"is_active"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type N8NFieldMapping struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	WorkflowID      uuid.UUID  `json:"workflow_id"`
	SourceField     string     `json:"source_field"`
	TargetField     string     `json:"target_field"`
	TransformRule   string     `json:"transform_rule"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// Cost Calculation Models
type CostCalculation struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProductName     string     `json:"product_name"`
	ProductSpecs    string     `json:"product_specs"`
	Quantity        int        `json:"quantity"`
	MaterialCost    float64    `json:"material_cost"`
	ProcessingCost  float64    `json:"processing_cost"`
	OverheadCost    float64    `json:"overhead_cost"`
	TotalCost       float64    `json:"total_cost"`
	CalculatedAt    time.Time  `json:"calculated_at"`
	CalculatedBy    uuid.UUID  `json:"calculated_by"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// Order Models
type OrderItem struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrderID         uuid.UUID  `json:"order_id"`
	ProductName     string     `json:"product_name"`
	ProductSpecs    string     `json:"product_specs"`
	Quantity        int        `json:"quantity"`
	UnitPrice       float64    `json:"unit_price"`
	TotalPrice      float64    `json:"total_price"`
	Status          string     `json:"status"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type ProductionSchedule struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrderID         uuid.UUID  `json:"order_id"`
	OrderItemID     uuid.UUID  `json:"order_item_id"`
	PlannedStart    time.Time  `json:"planned_start"`
	PlannedEnd      time.Time  `json:"planned_end"`
	ActualStart     *time.Time `json:"actual_start"`
	ActualEnd       *time.Time `json:"actual_end"`
	Status          string     `json:"status"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// Trade Models
type TradeTariffCode struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	HSCode          string     `json:"hs_code"`
	Description     string     `json:"description"`
	Category        string     `json:"category"`
	CompanyID       uuid.UUID  `json:"company_id"`
	IsActive        bool       `json:"is_active"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type TradeTariffRate struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TariffCodeID    uuid.UUID  `json:"tariff_code_id"`
	OriginCountry   string     `json:"origin_country"`
	DestCountry     string     `json:"dest_country"`
	Rate            float64    `json:"rate"`
	EffectiveFrom   time.Time  `json:"effective_from"`
	EffectiveTo     *time.Time `json:"effective_to"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type TradeShipment struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ShipmentNo      string     `json:"shipment_no"`
	OrderID         uuid.UUID  `json:"order_id"`
	Status          string     `json:"status"`
	ShippedDate     *time.Time `json:"shipped_date"`
	EstimatedArrival *time.Time `json:"estimated_arrival"`
	ActualArrival   *time.Time `json:"actual_arrival"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type ShipmentEvent struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ShipmentID      uuid.UUID  `json:"shipment_id"`
	EventType       string     `json:"event_type"`
	EventDate       time.Time  `json:"event_date"`
	Location        string     `json:"location"`
	Description     string     `json:"description"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type TradeDocument struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	DocumentType    string     `json:"document_type"`
	DocumentNo      string     `json:"document_no"`
	ShipmentID      uuid.UUID  `json:"shipment_id"`
	FilePath        string     `json:"file_path"`
	Status          string     `json:"status"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type LetterOfCredit struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	LCNumber        string     `json:"lc_number"`
	CustomerID      uuid.UUID  `json:"customer_id"`
	Amount          float64    `json:"amount"`
	Currency        string     `json:"currency"`
	IssuedDate      time.Time  `json:"issued_date"`
	ExpiryDate      time.Time  `json:"expiry_date"`
	Status          string     `json:"status"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type LCUtilization struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	LCID            uuid.UUID  `json:"lc_id"`
	OrderID         uuid.UUID  `json:"order_id"`
	Amount          float64    `json:"amount"`
	UtilizedDate    time.Time  `json:"utilized_date"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type TradeComplianceCheck struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ResourceType    string     `json:"resource_type"`
	ResourceID      string     `json:"resource_id"`
	CheckType       string     `json:"check_type"`
	Status          string     `json:"status"`
	Result          JSONB      `json:"result" gorm:"type:jsonb"`
	CheckedAt       time.Time  `json:"checked_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type ExchangeRate struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	FromCurrency    string     `json:"from_currency"`
	ToCurrency      string     `json:"to_currency"`
	Rate            float64    `json:"rate"`
	RateDate        time.Time  `json:"rate_date"`
	Source          string     `json:"source"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// Advanced Models
type AIAssistant struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name            string     `json:"name"`
	Type            string     `json:"type"`
	Configuration   JSONB      `json:"configuration" gorm:"type:jsonb"`
	CompanyID       uuid.UUID  `json:"company_id"`
	IsActive        bool       `json:"is_active"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type AIConversation struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	SessionID       string     `json:"session_id"`
	AssistantID     uuid.UUID  `json:"assistant_id"`
	UserID          uuid.UUID  `json:"user_id"`
	StartedAt       time.Time  `json:"started_at"`
	EndedAt         *time.Time `json:"ended_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type AIMessage struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ConversationID  uuid.UUID  `json:"conversation_id"`
	Role            string     `json:"role"`
	Content         string     `json:"content"`
	Metadata        JSONB      `json:"metadata" gorm:"type:jsonb"`
	CreatedAt       time.Time  `json:"created_at"`
}

type Recommendation struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID          uuid.UUID  `json:"user_id"`
	Type            string     `json:"type"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	ActionData      JSONB      `json:"action_data" gorm:"type:jsonb"`
	Status          string     `json:"status"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type AdvancedSearch struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name            string     `json:"name"`
	SearchType      string     `json:"search_type"`
	Criteria        JSONB      `json:"criteria" gorm:"type:jsonb"`
	UserID          uuid.UUID  `json:"user_id"`
	IsSaved         bool       `json:"is_saved"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type BatchOperation struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OperationType   string     `json:"operation_type"`
	TotalItems      int        `json:"total_items"`
	ProcessedItems  int        `json:"processed_items"`
	Status          string     `json:"status"`
	Parameters      JSONB      `json:"parameters" gorm:"type:jsonb"`
	Result          JSONB      `json:"result" gorm:"type:jsonb"`
	CompanyID       uuid.UUID  `json:"company_id"`
	CreatedBy       uuid.UUID  `json:"created_by"`
	StartedAt       time.Time  `json:"started_at"`
	CompletedAt     *time.Time `json:"completed_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type CustomField struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	FieldName       string     `json:"field_name"`
	FieldType       string     `json:"field_type"`
	EntityType      string     `json:"entity_type"`
	ValidationRules JSONB      `json:"validation_rules" gorm:"type:jsonb"`
	CompanyID       uuid.UUID  `json:"company_id"`
	IsRequired      bool       `json:"is_required"`
	IsActive        bool       `json:"is_active"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type CustomFieldValue struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	FieldID         uuid.UUID  `json:"field_id"`
	ResourceType    string     `json:"resource_type"`
	ResourceID      string     `json:"resource_id"`
	Value           string     `json:"value"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type SecurityEvent struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	EventType       string     `json:"event_type"`
	Severity        string     `json:"severity"`
	UserID          *uuid.UUID `json:"user_id"`
	IPAddress       string     `json:"ip_address"`
	UserAgent       string     `json:"user_agent"`
	Details         JSONB      `json:"details" gorm:"type:jsonb"`
	CompanyID       uuid.UUID  `json:"company_id"`
	OccurredAt      time.Time  `json:"occurred_at"`
	CreatedAt       time.Time  `json:"created_at"`
}

type PerformanceMetric struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ServiceName     string     `json:"service_name"`
	OperationName   string     `json:"operation_name"`
	Duration        int        `json:"duration"`
	Success         bool       `json:"success"`
	ErrorMessage    *string    `json:"error_message"`
	Metadata        JSONB      `json:"metadata" gorm:"type:jsonb"`
	RecordedAt      time.Time  `json:"recorded_at"`
	CreatedAt       time.Time  `json:"created_at"`
}

type PerformanceStats struct {
	ServiceName      string    `json:"service_name"`
	TotalRequests    int       `json:"total_requests"`
	SuccessfulRequests int     `json:"successful_requests"`
	FailedRequests   int       `json:"failed_requests"`
	AverageDuration  float64   `json:"average_duration"`
	MinDuration      int       `json:"min_duration"`
	MaxDuration      int       `json:"max_duration"`
}

type BackupRecord struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	BackupType      string     `json:"backup_type"`
	BackupSize      int64      `json:"backup_size"`
	BackupPath      string     `json:"backup_path"`
	Status          string     `json:"status"`
	CompanyID       uuid.UUID  `json:"company_id"`
	CreatedBy       uuid.UUID  `json:"created_by"`
	StartedAt       time.Time  `json:"started_at"`
	CompletedAt     *time.Time `json:"completed_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type SystemLanguage struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	LanguageCode    string     `json:"language_code"`
	LanguageName    string     `json:"language_name"`
	IsActive        bool       `json:"is_active"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// Integration Models
type IntegrationMapping struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	IntegrationID   uuid.UUID  `json:"integration_id"`
	SourceField     string     `json:"source_field"`
	TargetField     string     `json:"target_field"`
	TransformRule   string     `json:"transform_rule"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type WebhookDelivery struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	WebhookID       uuid.UUID  `json:"webhook_id"`
	Payload         JSONB      `json:"payload" gorm:"type:jsonb"`
	ResponseStatus  int        `json:"response_status"`
	ResponseBody    string     `json:"response_body"`
	DeliveredAt     time.Time  `json:"delivered_at"`
	CreatedAt       time.Time  `json:"created_at"`
}

type DataSyncJob struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	IntegrationID   uuid.UUID  `json:"integration_id"`
	JobType         string     `json:"job_type"`
	Status          string     `json:"status"`
	StartedAt       *time.Time `json:"started_at"`
	CompletedAt     *time.Time `json:"completed_at"`
	RecordsProcessed int       `json:"records_processed"`
	ErrorMessage    *string    `json:"error_message"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type ApiKey struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name            string     `json:"name"`
	KeyHash         string     `json:"key_hash"`
	CompanyID       uuid.UUID  `json:"company_id"`
	Permissions     JSONB      `json:"permissions" gorm:"type:jsonb"`
	ExpiresAt       *time.Time `json:"expires_at"`
	LastUsedAt      *time.Time `json:"last_used_at"`
	IsActive        bool       `json:"is_active"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type ExternalSystem struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	SystemName      string     `json:"system_name"`
	SystemType      string     `json:"system_type"`
	ConnectionConfig JSONB     `json:"connection_config" gorm:"type:jsonb"`
	CompanyID       uuid.UUID  `json:"company_id"`
	IsActive        bool       `json:"is_active"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type IntegrationTemplate struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TemplateName    string     `json:"template_name"`
	IntegrationType string     `json:"integration_type"`
	Configuration   JSONB      `json:"configuration" gorm:"type:jsonb"`
	IsPublic        bool       `json:"is_public"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}