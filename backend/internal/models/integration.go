package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Integration 整合配置
type Integration struct {
	ID                uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID         uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	Name              string     `gorm:"not null" json:"name"`
	Type              string     `gorm:"not null" json:"type"`              // api, webhook, ftp, sftp, email, database
	Provider          string     `gorm:"not null" json:"provider"`          // custom, erp, crm, accounting, shipping, payment
	Status            string     `gorm:"not null" json:"status"`            // active, inactive, error, testing
	Configuration     string     `json:"configuration"`                     // JSON configuration
	Credentials       string     `json:"credentials"`                       // Encrypted credentials JSON
	ApiVersion        string     `json:"api_version"`
	BaseURL           string     `json:"base_url"`
	AuthType          string     `json:"auth_type"`                         // none, api_key, oauth2, basic_auth, token
	Headers           string     `json:"headers"`                           // JSON headers
	RateLimitRPM      int        `json:"rate_limit_rpm"`                    // Requests per minute
	TimeoutSeconds    int        `json:"timeout_seconds"`
	RetryAttempts     int        `json:"retry_attempts"`
	IsActive          bool       `gorm:"default:true" json:"is_active"`
	LastSyncAt        *time.Time `json:"last_sync_at"`
	LastErrorAt       *time.Time `json:"last_error_at"`
	LastError         string     `json:"last_error"`
	SyncCount         int64      `gorm:"default:0" json:"sync_count"`
	ErrorCount        int64      `gorm:"default:0" json:"error_count"`
	SuccessRate       float64    `json:"success_rate"`
	AvgResponseTime   int64      `json:"avg_response_time"`                 // milliseconds
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	CreatedBy         uuid.UUID  `gorm:"type:uuid" json:"created_by"`

	// Relations
	Company        *Company               `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Creator        *User                  `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Mappings       []IntegrationMapping   `gorm:"foreignKey:IntegrationID" json:"mappings,omitempty"`
	Webhooks       []Webhook              `gorm:"foreignKey:IntegrationID" json:"webhooks,omitempty"`
	SyncJobs       []DataSyncJob          `gorm:"foreignKey:IntegrationID" json:"sync_jobs,omitempty"`
	Logs           []IntegrationLog       `gorm:"foreignKey:IntegrationID" json:"logs,omitempty"`
}

// IntegrationMapping 資料映射配置
type IntegrationMapping struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID       uuid.UUID `gorm:"type:uuid;not null" json:"company_id"`
	IntegrationID   uuid.UUID `gorm:"type:uuid;not null" json:"integration_id"`
	Name            string    `gorm:"not null" json:"name"`
	Direction       string    `gorm:"not null" json:"direction"`       // inbound, outbound, bidirectional
	SourceTable     string    `json:"source_table"`
	TargetTable     string    `json:"target_table"`
	SourceEndpoint  string    `json:"source_endpoint"`
	TargetEndpoint  string    `json:"target_endpoint"`
	FieldMappings   string    `json:"field_mappings"`                  // JSON field mapping rules
	Transformations string    `json:"transformations"`                 // JSON transformation rules
	Filters         string    `json:"filters"`                         // JSON filter conditions
	SyncFrequency   string    `gorm:"not null" json:"sync_frequency"`  // realtime, hourly, daily, weekly, manual
	IsActive        bool      `gorm:"default:true" json:"is_active"`
	LastSyncAt      *time.Time `json:"last_sync_at"`
	NextSyncAt      *time.Time `json:"next_sync_at"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	CreatedBy       uuid.UUID `gorm:"type:uuid" json:"created_by"`

	// Relations
	Company     *Company    `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Integration *Integration `gorm:"foreignKey:IntegrationID" json:"integration,omitempty"`
	Creator     *User       `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

// Webhook Webhook 配置
type Webhook struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID     uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	IntegrationID *uuid.UUID `gorm:"type:uuid" json:"integration_id"`
	Name          string     `gorm:"not null" json:"name"`
	URL           string     `gorm:"not null" json:"url"`
	Method        string     `gorm:"not null;default:'POST'" json:"method"`
	Headers       string     `json:"headers"`                           // JSON headers
	AuthType      string     `json:"auth_type"`                         // none, api_key, basic_auth, bearer_token
	AuthConfig    string     `json:"auth_config"`                       // JSON auth configuration
	Events        string     `json:"events"`                            // JSON array of event types
	PayloadFormat string     `gorm:"default:'json'" json:"payload_format"` // json, xml, form
	PayloadTemplate string   `json:"payload_template"`                  // Custom payload template
	IsActive      bool       `gorm:"default:true" json:"is_active"`
	RetryAttempts int        `gorm:"default:3" json:"retry_attempts"`
	RetryInterval int        `gorm:"default:60" json:"retry_interval"`  // seconds
	TimeoutSeconds int       `gorm:"default:30" json:"timeout_seconds"`
	LastTriggeredAt *time.Time `json:"last_triggered_at"`
	TriggerCount  int64      `gorm:"default:0" json:"trigger_count"`
	SuccessCount  int64      `gorm:"default:0" json:"success_count"`
	FailureCount  int64      `gorm:"default:0" json:"failure_count"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	CreatedBy     uuid.UUID  `gorm:"type:uuid" json:"created_by"`

	// Relations
	Company     *Company     `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Integration *Integration `gorm:"foreignKey:IntegrationID" json:"integration,omitempty"`
	Creator     *User        `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Deliveries  []WebhookDelivery `gorm:"foreignKey:WebhookID" json:"deliveries,omitempty"`
}

// WebhookDelivery Webhook 發送記錄
type WebhookDelivery struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	WebhookID    uuid.UUID  `gorm:"type:uuid;not null" json:"webhook_id"`
	CompanyID    uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	EventType    string     `gorm:"not null" json:"event_type"`
	EventData    string     `json:"event_data"`                        // JSON event data
	RequestURL   string     `json:"request_url"`
	RequestMethod string    `json:"request_method"`
	RequestHeaders string   `json:"request_headers"`                   // JSON headers
	RequestBody  string     `json:"request_body"`
	ResponseCode int        `json:"response_code"`
	ResponseHeaders string  `json:"response_headers"`                  // JSON headers
	ResponseBody string     `json:"response_body"`
	Status       string     `gorm:"not null" json:"status"`            // pending, success, failed, retrying
	AttemptCount int        `gorm:"default:1" json:"attempt_count"`
	NextRetryAt  *time.Time `json:"next_retry_at"`
	CompletedAt  *time.Time `json:"completed_at"`
	ErrorMessage string     `json:"error_message"`
	ResponseTime int64      `json:"response_time"`                     // milliseconds
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`

	// Relations
	Webhook *Webhook `gorm:"foreignKey:WebhookID" json:"webhook,omitempty"`
	Company *Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
}

// DataSyncJob 數據同步任務
type DataSyncJob struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID     uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	IntegrationID uuid.UUID  `gorm:"type:uuid;not null" json:"integration_id"`
	MappingID     *uuid.UUID `gorm:"type:uuid" json:"mapping_id"`
	Name          string     `gorm:"not null" json:"name"`
	Type          string     `gorm:"not null" json:"type"`              // full_sync, incremental_sync, delta_sync
	Direction     string     `gorm:"not null" json:"direction"`         // import, export, bidirectional
	Status        string     `gorm:"not null" json:"status"`            // pending, running, completed, failed, cancelled
	Priority      string     `gorm:"not null" json:"priority"`          // low, normal, high, urgent
	ScheduledAt   *time.Time `json:"scheduled_at"`
	StartedAt     *time.Time `json:"started_at"`
	CompletedAt   *time.Time `json:"completed_at"`
	Duration      int64      `json:"duration"`                          // seconds
	TotalRecords  int64      `json:"total_records"`
	ProcessedRecords int64   `json:"processed_records"`
	SuccessRecords int64     `json:"success_records"`
	ErrorRecords  int64      `json:"error_records"`
	SkippedRecords int64     `json:"skipped_records"`
	Progress      int        `json:"progress"`                          // 0-100
	Configuration string     `json:"configuration"`                     // JSON sync configuration
	Result        string     `json:"result"`                            // JSON sync result
	ErrorLog      string     `json:"error_log"`                         // JSON error details
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	CreatedBy     uuid.UUID  `gorm:"type:uuid" json:"created_by"`

	// Relations
	Company     *Company           `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Integration *Integration       `gorm:"foreignKey:IntegrationID" json:"integration,omitempty"`
	Mapping     *IntegrationMapping `gorm:"foreignKey:MappingID" json:"mapping,omitempty"`
	Creator     *User              `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

// IntegrationLog 整合日誌
type IntegrationLog struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID     uuid.UUID `gorm:"type:uuid;not null" json:"company_id"`
	IntegrationID uuid.UUID `gorm:"type:uuid;not null" json:"integration_id"`
	SyncJobID     *uuid.UUID `gorm:"type:uuid" json:"sync_job_id"`
	Level         string    `gorm:"not null" json:"level"`             // debug, info, warning, error, critical
	Category      string    `gorm:"not null" json:"category"`          // api_call, data_sync, webhook, auth, config
	Message       string    `gorm:"not null" json:"message"`
	Details       string    `json:"details"`                           // JSON additional details
	RequestData   string    `json:"request_data"`                      // JSON request data
	ResponseData  string    `json:"response_data"`                     // JSON response data
	ErrorCode     string    `json:"error_code"`
	ErrorMessage  string    `json:"error_message"`
	Duration      int64     `json:"duration"`                          // milliseconds
	UserID        *uuid.UUID `gorm:"type:uuid" json:"user_id"`
	IPAddress     string    `json:"ip_address"`
	UserAgent     string    `json:"user_agent"`
	TraceID       string    `json:"trace_id"`
	CreatedAt     time.Time `json:"created_at"`

	// Relations
	Company     *Company     `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Integration *Integration `gorm:"foreignKey:IntegrationID" json:"integration,omitempty"`
	SyncJob     *DataSyncJob `gorm:"foreignKey:SyncJobID" json:"sync_job,omitempty"`
	User        *User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// ApiKey API 金鑰管理
type ApiKey struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID   uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	UserID      uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	Name        string     `gorm:"not null" json:"name"`
	Description string     `json:"description"`
	KeyHash     string     `gorm:"not null;unique" json:"key_hash"`     // Hashed API key
	KeyPrefix   string     `gorm:"not null" json:"key_prefix"`          // First 8 chars for display
	Permissions string     `json:"permissions"`                         // JSON permissions array
	Scopes      string     `json:"scopes"`                              // JSON scopes array
	RateLimit   int        `json:"rate_limit"`                          // requests per minute
	IsActive    bool       `gorm:"default:true" json:"is_active"`
	ExpiresAt   *time.Time `json:"expires_at"`
	LastUsedAt  *time.Time `json:"last_used_at"`
	UsageCount  int64      `gorm:"default:0" json:"usage_count"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CreatedBy   uuid.UUID  `gorm:"type:uuid" json:"created_by"`

	// Relations
	Company *Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	User    *User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Creator *User    `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

// ExternalSystem 外部系統配置
type ExternalSystem struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID     uuid.UUID `gorm:"type:uuid;not null" json:"company_id"`
	Name          string    `gorm:"not null" json:"name"`
	SystemType    string    `gorm:"not null" json:"system_type"`       // erp, crm, accounting, warehouse, shipping
	Vendor        string    `json:"vendor"`                            // sap, oracle, salesforce, quickbooks
	Version       string    `json:"version"`
	BaseURL       string    `json:"base_url"`
	DatabaseConfig string   `json:"database_config"`                   // JSON database connection config
	ApiConfig     string    `json:"api_config"`                        // JSON API configuration
	FtpConfig     string    `json:"ftp_config"`                        // JSON FTP configuration
	SftpConfig    string    `json:"sftp_config"`                       // JSON SFTP configuration
	Status        string    `gorm:"not null" json:"status"`            // active, inactive, testing, error
	IsActive      bool      `gorm:"default:true" json:"is_active"`
	LastTestAt    *time.Time `json:"last_test_at"`
	LastTestResult string   `json:"last_test_result"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	CreatedBy     uuid.UUID `gorm:"type:uuid" json:"created_by"`

	// Relations
	Company      *Company      `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Creator      *User         `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Integrations []Integration `gorm:"foreignKey:Provider;references:SystemType" json:"integrations,omitempty"`
}

// DataTransformation 數據轉換規則
type DataTransformation struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID     uuid.UUID `gorm:"type:uuid;not null" json:"company_id"`
	MappingID     uuid.UUID `gorm:"type:uuid;not null" json:"mapping_id"`
	Name          string    `gorm:"not null" json:"name"`
	Description   string    `json:"description"`
	Type          string    `gorm:"not null" json:"type"`              // field_mapping, value_transformation, aggregation, calculation
	SourceField   string    `json:"source_field"`
	TargetField   string    `json:"target_field"`
	TransformRule string    `json:"transform_rule"`                    // JSON transformation rule
	Conditions    string    `json:"conditions"`                        // JSON condition rules
	DefaultValue  string    `json:"default_value"`
	IsRequired    bool      `gorm:"default:false" json:"is_required"`
	IsActive      bool      `gorm:"default:true" json:"is_active"`
	ExecutionOrder int      `json:"execution_order"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	CreatedBy     uuid.UUID `gorm:"type:uuid" json:"created_by"`

	// Relations
	Company *Company           `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Mapping *IntegrationMapping `gorm:"foreignKey:MappingID" json:"mapping,omitempty"`
	Creator *User              `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

// IntegrationTemplate 整合模板
type IntegrationTemplate struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID     *uuid.UUID `gorm:"type:uuid" json:"company_id"`        // NULL for system templates
	Name          string     `gorm:"not null" json:"name"`
	Description   string     `json:"description"`
	Category      string     `gorm:"not null" json:"category"`           // erp, crm, accounting, shipping, payment
	Provider      string     `gorm:"not null" json:"provider"`
	Version       string     `json:"version"`
	Configuration string     `json:"configuration"`                      // JSON template configuration
	Mappings      string     `json:"mappings"`                           // JSON default mappings
	Requirements  string     `json:"requirements"`                       // JSON requirements
	IsPublic      bool       `gorm:"default:false" json:"is_public"`
	IsActive      bool       `gorm:"default:true" json:"is_active"`
	UsageCount    int64      `gorm:"default:0" json:"usage_count"`
	Rating        float64    `json:"rating"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	CreatedBy     uuid.UUID  `gorm:"type:uuid" json:"created_by"`

	// Relations
	Company *Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Creator *User    `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

// BeforeCreate hooks
func (i *Integration) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}

func (im *IntegrationMapping) BeforeCreate(tx *gorm.DB) error {
	if im.ID == uuid.Nil {
		im.ID = uuid.New()
	}
	return nil
}

func (w *Webhook) BeforeCreate(tx *gorm.DB) error {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return nil
}

func (wd *WebhookDelivery) BeforeCreate(tx *gorm.DB) error {
	if wd.ID == uuid.Nil {
		wd.ID = uuid.New()
	}
	return nil
}

func (dsj *DataSyncJob) BeforeCreate(tx *gorm.DB) error {
	if dsj.ID == uuid.Nil {
		dsj.ID = uuid.New()
	}
	return nil
}

func (il *IntegrationLog) BeforeCreate(tx *gorm.DB) error {
	if il.ID == uuid.Nil {
		il.ID = uuid.New()
	}
	return nil
}

func (ak *ApiKey) BeforeCreate(tx *gorm.DB) error {
	if ak.ID == uuid.Nil {
		ak.ID = uuid.New()
	}
	return nil
}

func (es *ExternalSystem) BeforeCreate(tx *gorm.DB) error {
	if es.ID == uuid.Nil {
		es.ID = uuid.New()
	}
	return nil
}

func (dt *DataTransformation) BeforeCreate(tx *gorm.DB) error {
	if dt.ID == uuid.Nil {
		dt.ID = uuid.New()
	}
	return nil
}

func (it *IntegrationTemplate) BeforeCreate(tx *gorm.DB) error {
	if it.ID == uuid.Nil {
		it.ID = uuid.New()
	}
	return nil
}