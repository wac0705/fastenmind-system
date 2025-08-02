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

// Cost Calculation Models (CostCalculation is already defined in process_cost.go)

// Order Models (OrderItem is already defined in order.go)

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

// ShipmentEvent is already defined in trade.go

// TradeDocument is already defined in trade.go

// LetterOfCredit is already defined in trade.go

// LCUtilization is already defined in trade.go

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

// ExchangeRate is already defined in trade.go

// Advanced Models - only define missing ones
// Note: Most advanced models are already defined in advanced.go

// AIConversation is not defined in advanced.go
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

// Recommendation is not defined in advanced.go
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

// PerformanceStats is not defined in advanced.go
type PerformanceStats struct {
	ServiceName      string    `json:"service_name"`
	TotalRequests    int       `json:"total_requests"`
	SuccessfulRequests int     `json:"successful_requests"`
	FailedRequests   int       `json:"failed_requests"`
	AverageDuration  float64   `json:"average_duration"`
	MinDuration      int       `json:"min_duration"`
	MaxDuration      int       `json:"max_duration"`
}

// Integration Models
// The following models are already defined in integration.go:
// - IntegrationMapping
// - WebhookDelivery
// - DataSyncJob
// - ApiKey
// - ExternalSystem
// - IntegrationTemplate