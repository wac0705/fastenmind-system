package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// QuoteVersion 報價單版本
type QuoteVersion struct {
	ID            uuid.UUID    `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	QuoteID       uuid.UUID    `json:"quote_id" gorm:"type:uuid;not null"`
	Quote         *Quote       `json:"quote" gorm:"foreignKey:QuoteID"`
	VersionNumber int          `json:"version_number" gorm:"not null"`
	VersionNotes  string       `json:"version_notes" gorm:"type:text"`
	ChangeSummary string       `json:"change_summary" gorm:"type:text"`
	IsCurrent     bool         `json:"is_current" gorm:"default:false"`
	Items         []QuoteItem  `json:"items" gorm:"foreignKey:QuoteVersionID"`
	Terms         []QuoteTerm  `json:"terms" gorm:"foreignKey:QuoteVersionID"`
	
	// Cost fields (for compatibility with service)
	MaterialCost   float64      `json:"material_cost"`
	ProcessCost    float64      `json:"process_cost"`
	SurfaceCost    float64      `json:"surface_cost"`
	HeatTreatCost  float64      `json:"heat_treat_cost"`
	PackagingCost  float64      `json:"packaging_cost"`
	ShippingCost   float64      `json:"shipping_cost"`
	TariffCost     float64      `json:"tariff_cost"`
	OverheadRate   float64      `json:"overhead_rate"`
	ProfitRate     float64      `json:"profit_rate"`
	TotalCost      float64      `json:"total_cost"`
	UnitPrice      float64      `json:"unit_price"`
	
	CreatedBy     uuid.UUID    `json:"created_by" gorm:"type:uuid;not null"`
	Creator       Account      `json:"creator" gorm:"foreignKey:CreatedBy"`
	CreatedAt     time.Time    `json:"created_at" gorm:"autoCreateTime"`
}

// QuoteItem 報價單項目明細
type QuoteItem struct {
	ID                uuid.UUID         `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	QuoteVersionID    uuid.UUID         `json:"quote_version_id" gorm:"type:uuid;not null"`
	ItemNo            int               `json:"item_no" gorm:"not null"`
	ProductName       string            `json:"product_name" gorm:"type:varchar(200);not null"`
	ProductSpecs      string            `json:"product_specs" gorm:"type:text"`
	Quantity          int               `json:"quantity" gorm:"not null"`
	Unit              string            `json:"unit" gorm:"type:varchar(20);not null"`
	UnitPrice         float64           `json:"unit_price" gorm:"type:decimal(15,4);not null"`
	TotalPrice        float64           `json:"total_price" gorm:"type:decimal(15,4);not null"`
	CostCalculationID *uuid.UUID        `json:"cost_calculation_id" gorm:"type:uuid"`
	CostCalculation   *CostCalculation  `json:"cost_calculation" gorm:"foreignKey:CostCalculationID"`
	MarginPercentage  float64           `json:"margin_percentage" gorm:"type:decimal(5,2)"`
	Notes             string            `json:"notes" gorm:"type:text"`
	CreatedAt         time.Time         `json:"created_at" gorm:"autoCreateTime"`
}

// QuoteApproval 報價單審核記錄
type QuoteApproval struct {
	ID                 uuid.UUID     `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	QuoteID            uuid.UUID     `json:"quote_id" gorm:"type:uuid;not null"`
	Quote              *Quote        `json:"quote" gorm:"foreignKey:QuoteID"`
	QuoteVersionID     uuid.UUID     `json:"quote_version_id" gorm:"type:uuid;not null"`
	QuoteVersion       *QuoteVersion `json:"quote_version" gorm:"foreignKey:QuoteVersionID"`
	ApprovalLevel      int           `json:"approval_level" gorm:"not null"` // 1: 初審, 2: 複審, 3: 終審
	ApproverRole       string        `json:"approver_role" gorm:"type:varchar(20);not null"`
	RequiredApproverID *uuid.UUID    `json:"required_approver_id" gorm:"type:uuid"`
	RequiredApprover   *Account      `json:"required_approver" gorm:"foreignKey:RequiredApproverID"`
	ActualApproverID   *uuid.UUID    `json:"actual_approver_id" gorm:"type:uuid"`
	ActualApprover     *Account      `json:"actual_approver" gorm:"foreignKey:ActualApproverID"`
	ApprovalStatus     string        `json:"approval_status" gorm:"type:varchar(20);not null;default:'pending'"`
	ApprovalNotes      string        `json:"approval_notes" gorm:"type:text"`
	ApprovedAt         *time.Time    `json:"approved_at"`
	CreatedAt          time.Time     `json:"created_at" gorm:"autoCreateTime"`
}

// QuoteTermsTemplate 報價單條款模板
type QuoteTermsTemplate struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TemplateName string    `json:"template_name" gorm:"type:varchar(100);not null"`
	TemplateType string    `json:"template_type" gorm:"type:varchar(50);not null"`
	Content      string    `json:"content" gorm:"type:text;not null"`
	Language     string    `json:"language" gorm:"type:varchar(10);default:'zh-TW'"`
	IsDefault    bool      `json:"is_default" gorm:"default:false"`
	IsActive     bool      `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// QuoteTerm 報價單條款
type QuoteTerm struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	QuoteVersionID uuid.UUID `json:"quote_version_id" gorm:"type:uuid;not null"`
	TermType       string    `json:"term_type" gorm:"type:varchar(50);not null"`
	TermContent    string    `json:"term_content" gorm:"type:text;not null"`
	SortOrder      int       `json:"sort_order" gorm:"default:0"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// QuoteAttachment 報價單附件
type QuoteAttachment struct {
	ID             uuid.UUID     `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	QuoteID        uuid.UUID     `json:"quote_id" gorm:"type:uuid;not null"`
	Quote          *Quote        `json:"quote" gorm:"foreignKey:QuoteID"`
	QuoteVersionID *uuid.UUID    `json:"quote_version_id" gorm:"type:uuid"`
	QuoteVersion   *QuoteVersion `json:"quote_version" gorm:"foreignKey:QuoteVersionID"`
	FileName       string        `json:"file_name" gorm:"type:varchar(255);not null"`
	FilePath       string        `json:"file_path" gorm:"type:varchar(500);not null"`
	FileSize       int           `json:"file_size"`
	FileType       string        `json:"file_type" gorm:"type:varchar(50)"`
	Description    string        `json:"description" gorm:"type:text"`
	UploadedBy     uuid.UUID     `json:"uploaded_by" gorm:"type:uuid;not null"`
	Uploader       Account       `json:"uploader" gorm:"foreignKey:UploadedBy"`
	UploadedAt     time.Time     `json:"uploaded_at" gorm:"autoCreateTime"`
}

// QuoteActivityLog 報價單活動日誌
type QuoteActivityLog struct {
	ID                  uuid.UUID        `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	QuoteID             uuid.UUID        `json:"quote_id" gorm:"type:uuid;not null"`
	Quote               *Quote           `json:"quote" gorm:"foreignKey:QuoteID"`
	QuoteVersionID      *uuid.UUID       `json:"quote_version_id" gorm:"type:uuid"`
	QuoteVersion        *QuoteVersion    `json:"quote_version" gorm:"foreignKey:QuoteVersionID"`
	ActivityType        string           `json:"activity_type" gorm:"type:varchar(50);not null"`
	ActivityDescription string           `json:"activity_description" gorm:"type:text"`
	ActivityData        datatypes.JSON   `json:"activity_data" gorm:"type:jsonb"`
	PerformedBy         uuid.UUID        `json:"performed_by" gorm:"type:uuid;not null"`
	Performer           Account          `json:"performer" gorm:"foreignKey:PerformedBy"`
	PerformedAt         time.Time        `json:"performed_at" gorm:"autoCreateTime"`
}

// QuoteSendLog 報價單發送記錄
type QuoteSendLog struct {
	ID             uuid.UUID               `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	QuoteID        uuid.UUID               `json:"quote_id" gorm:"type:uuid;not null"`
	Quote          *Quote                  `json:"quote" gorm:"foreignKey:QuoteID"`
	QuoteVersionID uuid.UUID               `json:"quote_version_id" gorm:"type:uuid;not null"`
	QuoteVersion   *QuoteVersion           `json:"quote_version" gorm:"foreignKey:QuoteVersionID"`
	SendMethod     string                  `json:"send_method" gorm:"type:varchar(20);not null"`
	RecipientEmail string                  `json:"recipient_email" gorm:"type:varchar(255)"`
	RecipientName  string                  `json:"recipient_name" gorm:"type:varchar(100)"`
	CcEmails       datatypes.JSON          `json:"cc_emails" gorm:"type:jsonb"`
	Subject        string                  `json:"subject" gorm:"type:varchar(500)"`
	Message        string                  `json:"message" gorm:"type:text"`
	Attachments    datatypes.JSON          `json:"attachments" gorm:"type:jsonb"`
	SendStatus     string                  `json:"send_status" gorm:"type:varchar(20);not null;default:'pending'"`
	SentAt         *time.Time              `json:"sent_at"`
	ErrorMessage   string                  `json:"error_message" gorm:"type:text"`
	CreatedBy      uuid.UUID               `json:"created_by" gorm:"type:uuid;not null"`
	Creator        Account                 `json:"creator" gorm:"foreignKey:CreatedBy"`
	CreatedAt      time.Time               `json:"created_at" gorm:"autoCreateTime"`
}

// QuoteTemplate 報價單模板
type QuoteTemplate struct {
	ID              uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TemplateName    string    `json:"template_name" gorm:"type:varchar(100);not null"`
	Description     string    `json:"description" gorm:"type:text"`
	HeaderLogoPath  string    `json:"header_logo_path" gorm:"type:varchar(500)"`
	FooterContent   string    `json:"footer_content" gorm:"type:text"`
	TermsConditions string    `json:"terms_conditions" gorm:"type:text"`
	CssStyles       string    `json:"css_styles" gorm:"type:text"`
	IsDefault       bool      `json:"is_default" gorm:"default:false"`
	IsActive        bool      `json:"is_active" gorm:"default:true"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName overrides
func (QuoteVersion) TableName() string       { return "quote_versions" }
func (QuoteItem) TableName() string           { return "quote_items" }
func (QuoteApproval) TableName() string       { return "quote_approvals" }
func (QuoteTermsTemplate) TableName() string  { return "quote_terms_templates" }
func (QuoteTerm) TableName() string           { return "quote_terms" }
func (QuoteAttachment) TableName() string     { return "quote_attachments" }
func (QuoteActivityLog) TableName() string    { return "quote_activity_logs" }
func (QuoteSendLog) TableName() string        { return "quote_send_logs" }
func (QuoteTemplate) TableName() string       { return "quote_templates" }

// Request/Response structures

// CreateQuoteRequest 創建報價單請求
type CreateQuoteRequest struct {
	InquiryID        uuid.UUID           `json:"inquiry_id" binding:"required"`
	CustomerID       uuid.UUID           `json:"customer_id" binding:"required"`
	ValidityDays     int                 `json:"validity_days"`
	PaymentTerms     string              `json:"payment_terms"`
	DeliveryTerms    string              `json:"delivery_terms"`
	Remarks          string              `json:"remarks"`
	Items            []QuoteItemRequest  `json:"items" binding:"required,min=1"`
	Terms            []QuoteTermRequest  `json:"terms"`
	UseTemplate      bool                `json:"use_template"`
	TemplateID       *uuid.UUID          `json:"template_id"`
}

// QuoteItemRequest 報價單項目請求
type QuoteItemRequest struct {
	ProductName       string     `json:"product_name" binding:"required"`
	ProductSpecs      string     `json:"product_specs"`
	Quantity          int        `json:"quantity" binding:"required,min=1"`
	Unit              string     `json:"unit" binding:"required"`
	UnitPrice         float64    `json:"unit_price" binding:"required,min=0"`
	CostCalculationID *uuid.UUID `json:"cost_calculation_id"`
	Notes             string     `json:"notes"`
}

// QuoteTermRequest 報價單條款請求
type QuoteTermRequest struct {
	TermType    string `json:"term_type" binding:"required"`
	TermContent string `json:"term_content" binding:"required"`
	SortOrder   int    `json:"sort_order"`
}

// UpdateQuoteRequest 更新報價單請求
type UpdateQuoteRequest struct {
	CreateNewVersion bool                `json:"create_new_version"`
	VersionNotes     string              `json:"version_notes"`
	ValidityDays     int                 `json:"validity_days"`
	PaymentTerms     string              `json:"payment_terms"`
	DeliveryTerms    string              `json:"delivery_terms"`
	Remarks          string              `json:"remarks"`
	Items            []QuoteItemRequest  `json:"items"`
	Terms            []QuoteTermRequest  `json:"terms"`
}

// SubmitApprovalRequest 提交審核請求
type SubmitApprovalRequest struct {
	Notes string `json:"notes"`
}

// ApproveQuoteRequest 審核報價單請求
type ApproveQuoteRequest struct {
	Approved bool   `json:"approved" binding:"required"`
	Notes    string `json:"notes"`
}

// SendQuoteRequest 發送報價單請求
type SendQuoteRequest struct {
	RecipientEmail string   `json:"recipient_email" binding:"required,email"`
	RecipientName  string   `json:"recipient_name"`
	CcEmails       []string `json:"cc_emails"`
	Recipients     []string `json:"recipients"`      // All recipients
	Subject        string   `json:"subject"`
	Body           string   `json:"body"`            // Email body
	Message        string   `json:"message"`
	AttachPDF      bool     `json:"attach_pdf"`
	AttachmentIDs  []string `json:"attachment_ids"`
}