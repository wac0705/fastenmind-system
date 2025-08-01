package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AIAssistant AI 助手
type AIAssistant struct {
	ID                uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID         uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	UserID            uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	Name              string     `gorm:"not null" json:"name"`
	Type              string     `gorm:"not null" json:"type"`         // chat, recommendation, analysis, automation
	Model             string     `gorm:"not null" json:"model"`        // gpt-4, claude-3, gemini-pro
	Status            string     `gorm:"not null" json:"status"`       // active, inactive, training, error
	Configuration     string     `json:"configuration"`                // JSON configuration
	SystemPrompt      string     `json:"system_prompt"`
	Temperature       float64    `json:"temperature"`
	MaxTokens         int        `json:"max_tokens"`
	TopP              float64    `json:"top_p"`
	FrequencyPenalty  float64    `json:"frequency_penalty"`
	PresencePenalty   float64    `json:"presence_penalty"`
	IsActive          bool       `gorm:"default:true" json:"is_active"`
	UsageCount        int64      `gorm:"default:0" json:"usage_count"`
	TokensUsed        int64      `gorm:"default:0" json:"tokens_used"`
	CostAccumulated   float64    `gorm:"default:0" json:"cost_accumulated"`
	LastUsed          *time.Time `json:"last_used"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	CreatedBy         uuid.UUID  `gorm:"type:uuid" json:"created_by"`

	// Relations
	Company   *Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	User      *User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Creator   *User    `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Sessions  []AIConversationSession `gorm:"foreignKey:AssistantID" json:"sessions,omitempty"`
}

// AIConversationSession AI 對話會話
type AIConversationSession struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	AssistantID  uuid.UUID  `gorm:"type:uuid;not null" json:"assistant_id"`
	UserID       uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	CompanyID    uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	Title        string     `json:"title"`
	Context      string     `json:"context"`                          // JSON context data
	Status       string     `gorm:"not null" json:"status"`           // active, completed, archived
	MessageCount int        `gorm:"default:0" json:"message_count"`
	TokensUsed   int64      `gorm:"default:0" json:"tokens_used"`
	Cost         float64    `gorm:"default:0" json:"cost"`
	StartTime    time.Time  `json:"start_time"`
	EndTime      *time.Time `json:"end_time"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`

	// Relations
	Assistant *AIAssistant `gorm:"foreignKey:AssistantID" json:"assistant,omitempty"`
	User      *User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Company   *Company     `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Messages  []AIMessage  `gorm:"foreignKey:SessionID" json:"messages,omitempty"`
}

// AIMessage AI 對話訊息
type AIMessage struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	SessionID  uuid.UUID `gorm:"type:uuid;not null" json:"session_id"`
	Role       string    `gorm:"not null" json:"role"`           // user, assistant, system
	Content    string    `gorm:"type:text" json:"content"`
	TokenCount int       `json:"token_count"`
	Cost       float64   `json:"cost"`
	ModelUsed  string    `json:"model_used"`
	ResponseTime int64   `json:"response_time"`                  // milliseconds
	Metadata   string    `json:"metadata"`                       // JSON metadata
	CreatedAt  time.Time `json:"created_at"`

	// Relations
	Session *AIConversationSession `gorm:"foreignKey:SessionID" json:"session,omitempty"`
}

// SmartRecommendation 智能推薦
type SmartRecommendation struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID    uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	UserID       uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	Type         string     `gorm:"not null" json:"type"`         // product, customer, supplier, pricing, process
	Category     string     `json:"category"`
	Title        string     `gorm:"not null" json:"title"`
	Description  string     `json:"description"`
	Data         string     `json:"data"`                         // JSON recommendation data
	Score        float64    `json:"score"`                        // confidence score 0-1
	Priority     string     `gorm:"not null" json:"priority"`     // low, medium, high, urgent
	Status       string     `gorm:"not null" json:"status"`       // pending, viewed, accepted, rejected, implemented
	Source       string     `json:"source"`                       // ai, algorithm, manual, user_behavior
	SourceData   string     `json:"source_data"`                  // JSON source analysis data
	ResourceID   *uuid.UUID `gorm:"type:uuid" json:"resource_id"` // related resource ID
	ResourceType string     `json:"resource_type"`                // inquiry, quote, order, product, customer
	ExpiresAt    *time.Time `json:"expires_at"`
	ViewedAt     *time.Time `json:"viewed_at"`
	ActionedAt   *time.Time `json:"actioned_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`

	// Relations
	Company *Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	User    *User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// AdvancedSearch 高級搜索配置
type AdvancedSearch struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID   uuid.UUID `gorm:"type:uuid;not null" json:"company_id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	SearchType  string    `gorm:"not null" json:"search_type"`   // table, global, cross_reference
	TableName   string    `json:"table_name"`
	Filters     string    `json:"filters"`                       // JSON filter configuration
	SortConfig  string    `json:"sort_config"`                   // JSON sort configuration
	Columns     string    `json:"columns"`                       // JSON column configuration
	IsPublic    bool      `gorm:"default:false" json:"is_public"`
	UsageCount  int64     `gorm:"default:0" json:"usage_count"`
	LastUsed    *time.Time `json:"last_used"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// Relations
	Company *Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	User    *User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// BatchOperation 批量操作
type BatchOperation struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID     uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	UserID        uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	OperationType string     `gorm:"not null" json:"operation_type"` // update, delete, export, import, send_email, change_status
	TargetTable   string     `gorm:"not null" json:"target_table"`
	TargetIDs     string     `json:"target_ids"`                     // JSON array of target IDs
	Parameters    string     `json:"parameters"`                     // JSON operation parameters
	Status        string     `gorm:"not null" json:"status"`         // pending, running, completed, failed, cancelled
	Progress      int        `gorm:"default:0" json:"progress"`      // 0-100
	TotalItems    int        `json:"total_items"`
	ProcessedItems int       `gorm:"default:0" json:"processed_items"`
	SuccessCount  int        `gorm:"default:0" json:"success_count"`
	ErrorCount    int        `gorm:"default:0" json:"error_count"`
	ErrorLog      string     `json:"error_log"`                      // JSON error details
	Result        string     `json:"result"`                         // JSON operation result
	StartedAt     *time.Time `json:"started_at"`
	CompletedAt   *time.Time `json:"completed_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`

	// Relations
	Company *Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	User    *User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// CustomField 自訂欄位
type CustomField struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID    uuid.UUID `gorm:"type:uuid;not null" json:"company_id"`
	TableName    string    `gorm:"not null" json:"table_name"`
	FieldName    string    `gorm:"not null" json:"field_name"`
	FieldLabel   string    `gorm:"not null" json:"field_label"`
	FieldType    string    `gorm:"not null" json:"field_type"`    // text, number, date, boolean, select, multi_select, file
	DefaultValue string    `json:"default_value"`
	Options      string    `json:"options"`                       // JSON options for select types
	Validation   string    `json:"validation"`                    // JSON validation rules
	IsRequired   bool      `gorm:"default:false" json:"is_required"`
	IsSearchable bool      `gorm:"default:true" json:"is_searchable"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	DisplayOrder int       `json:"display_order"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedBy    uuid.UUID `gorm:"type:uuid" json:"created_by"`

	// Relations
	Company *Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Creator *User    `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

// CustomFieldValue 自訂欄位值
type CustomFieldValue struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID   uuid.UUID `gorm:"type:uuid;not null" json:"company_id"`
	FieldID     uuid.UUID `gorm:"type:uuid;not null" json:"field_id"`
	ResourceID  uuid.UUID `gorm:"type:uuid;not null" json:"resource_id"`
	ResourceType string   `gorm:"not null" json:"resource_type"`
	Value       string    `json:"value"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relations
	Company *Company    `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Field   *CustomField `gorm:"foreignKey:FieldID" json:"field,omitempty"`
}

// SecurityEvent 安全事件
type SecurityEvent struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID     uuid.UUID `gorm:"type:uuid;not null" json:"company_id"`
	UserID        *uuid.UUID `gorm:"type:uuid" json:"user_id"`
	EventType     string    `gorm:"not null" json:"event_type"`     // login, logout, failed_login, data_access, permission_change, suspicious_activity
	Severity      string    `gorm:"not null" json:"severity"`       // low, medium, high, critical
	Description   string    `json:"description"`
	IPAddress     string    `json:"ip_address"`
	UserAgent     string    `json:"user_agent"`
	Location      string    `json:"location"`                       // Geographic location
	DeviceInfo    string    `json:"device_info"`                    // JSON device information
	ResourceType  string    `json:"resource_type"`
	ResourceID    *uuid.UUID `gorm:"type:uuid" json:"resource_id"`
	ActionDetails string    `json:"action_details"`                 // JSON action details
	RiskScore     float64   `json:"risk_score"`                     // 0-100
	Status        string    `gorm:"not null" json:"status"`         // new, investigating, resolved, false_positive
	ResolvedBy    *uuid.UUID `gorm:"type:uuid" json:"resolved_by"`
	ResolvedAt    *time.Time `json:"resolved_at"`
	CreatedAt     time.Time  `json:"created_at"`

	// Relations
	Company    *Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	User       *User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ResolvedByUser *User `gorm:"foreignKey:ResolvedBy" json:"resolved_by_user,omitempty"`
}

// PerformanceMetric 效能指標
type PerformanceMetric struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID   uuid.UUID `gorm:"type:uuid;not null" json:"company_id"`
	MetricType  string    `gorm:"not null" json:"metric_type"`    // api_response_time, database_query_time, page_load_time, memory_usage, cpu_usage
	MetricName  string    `gorm:"not null" json:"metric_name"`
	Value       float64   `json:"value"`
	Unit        string    `json:"unit"`                           // ms, seconds, mb, gb, percentage
	Context     string    `json:"context"`                        // JSON context information
	Endpoint    string    `json:"endpoint"`
	Method      string    `json:"method"`
	StatusCode  int       `json:"status_code"`
	UserID      *uuid.UUID `gorm:"type:uuid" json:"user_id"`
	SessionID   string    `json:"session_id"`
	TraceID     string    `json:"trace_id"`
	CreatedAt   time.Time `json:"created_at"`

	// Relations
	Company *Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	User    *User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// BackupRecord 備份記錄
type BackupRecord struct {
	ID              uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID       uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	BackupType      string     `gorm:"not null" json:"backup_type"`      // full, incremental, differential
	Status          string     `gorm:"not null" json:"status"`           // running, completed, failed, cancelled
	FileSize        int64      `json:"file_size"`                        // bytes
	FilePath        string     `json:"file_path"`
	Checksum        string     `json:"checksum"`
	Tables          string     `json:"tables"`                           // JSON array of backed up tables
	CompressionType string     `json:"compression_type"`                 // gzip, zip, none
	EncryptionType  string     `json:"encryption_type"`                  // aes256, none
	StartTime       time.Time  `json:"start_time"`
	EndTime         *time.Time `json:"end_time"`
	Duration        int64      `json:"duration"`                         // seconds
	ErrorMessage    string     `json:"error_message"`
	CreatedBy       uuid.UUID  `gorm:"type:uuid" json:"created_by"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`

	// Relations
	Company *Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Creator *User    `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

// SystemLanguage 系統語言
type SystemLanguage struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID       *uuid.UUID `gorm:"type:uuid" json:"company_id"`
	LanguageCode    string    `gorm:"not null;unique" json:"language_code"` // en, zh-TW, zh-CN, ja, ko
	LanguageName    string    `gorm:"not null" json:"language_name"`
	NativeName      string    `gorm:"not null" json:"native_name"`
	IsActive        bool      `gorm:"default:true" json:"is_active"`
	IsDefault       bool      `gorm:"default:false" json:"is_default"`
	RTL             bool      `gorm:"default:false" json:"rtl"`             // Right-to-left
	DateFormat      string    `json:"date_format"`
	TimeFormat      string    `json:"time_format"`
	NumberFormat    string    `json:"number_format"`
	CurrencyFormat  string    `json:"currency_format"`
	TranslationProgress float64 `json:"translation_progress"`             // 0-100
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	// Relations
	Company      *Company           `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Translations []Translation     `gorm:"foreignKey:LanguageID" json:"translations,omitempty"`
}

// Translation 翻譯
type Translation struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	LanguageID     uuid.UUID `gorm:"type:uuid;not null" json:"language_id"`
	TranslationKey string    `gorm:"not null" json:"translation_key"`
	Translation    string    `gorm:"type:text" json:"translation"`
	Context        string    `json:"context"`
	IsApproved     bool      `gorm:"default:false" json:"is_approved"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	CreatedBy      uuid.UUID `gorm:"type:uuid" json:"created_by"`
	ApprovedBy     *uuid.UUID `gorm:"type:uuid" json:"approved_by"`

	// Relations
	Language     *SystemLanguage `gorm:"foreignKey:LanguageID" json:"language,omitempty"`
	Creator      *User           `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Approver     *User           `gorm:"foreignKey:ApprovedBy" json:"approver,omitempty"`
}

// BeforeCreate hooks
func (ai *AIAssistant) BeforeCreate(tx *gorm.DB) error {
	if ai.ID == uuid.Nil {
		ai.ID = uuid.New()
	}
	return nil
}

func (acs *AIConversationSession) BeforeCreate(tx *gorm.DB) error {
	if acs.ID == uuid.Nil {
		acs.ID = uuid.New()
	}
	return nil
}

func (am *AIMessage) BeforeCreate(tx *gorm.DB) error {
	if am.ID == uuid.Nil {
		am.ID = uuid.New()
	}
	return nil
}

func (sr *SmartRecommendation) BeforeCreate(tx *gorm.DB) error {
	if sr.ID == uuid.Nil {
		sr.ID = uuid.New()
	}
	return nil
}

func (as *AdvancedSearch) BeforeCreate(tx *gorm.DB) error {
	if as.ID == uuid.Nil {
		as.ID = uuid.New()
	}
	return nil
}

func (bo *BatchOperation) BeforeCreate(tx *gorm.DB) error {
	if bo.ID == uuid.Nil {
		bo.ID = uuid.New()
	}
	return nil
}

func (cf *CustomField) BeforeCreate(tx *gorm.DB) error {
	if cf.ID == uuid.Nil {
		cf.ID = uuid.New()
	}
	return nil
}

func (cfv *CustomFieldValue) BeforeCreate(tx *gorm.DB) error {
	if cfv.ID == uuid.Nil {
		cfv.ID = uuid.New()
	}
	return nil
}

func (se *SecurityEvent) BeforeCreate(tx *gorm.DB) error {
	if se.ID == uuid.Nil {
		se.ID = uuid.New()
	}
	return nil
}

func (pm *PerformanceMetric) BeforeCreate(tx *gorm.DB) error {
	if pm.ID == uuid.Nil {
		pm.ID = uuid.New()
	}
	return nil
}

func (br *BackupRecord) BeforeCreate(tx *gorm.DB) error {
	if br.ID == uuid.Nil {
		br.ID = uuid.New()
	}
	return nil
}

func (sl *SystemLanguage) BeforeCreate(tx *gorm.DB) error {
	if sl.ID == uuid.Nil {
		sl.ID = uuid.New()
	}
	return nil
}

func (t *Translation) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}