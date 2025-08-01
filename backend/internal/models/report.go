package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Report represents a business report
type Report struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID     uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	ReportNo      string     `gorm:"not null;unique" json:"report_no"`
	Name          string     `gorm:"not null" json:"name"`
	NameEn        string     `json:"name_en"`
	Category      string     `gorm:"not null" json:"category"` // sales, finance, production, inventory, supplier, customer, system
	Type          string     `gorm:"not null" json:"type"`     // summary, detail, trend, comparison, dashboard
	Status        string     `gorm:"default:'active'" json:"status"` // active, inactive, archived
	
	// Report Configuration
	DataSource    string     `json:"data_source"`    // JSON configuration for data sources
	Filters       string     `json:"filters"`        // JSON configuration for filters
	Columns       string     `json:"columns"`        // JSON configuration for columns
	Sorting       string     `json:"sorting"`        // JSON configuration for sorting
	Grouping      string     `json:"grouping"`       // JSON configuration for grouping
	Aggregation   string     `json:"aggregation"`    // JSON configuration for aggregation
	ChartConfig   string     `json:"chart_config"`   // JSON configuration for charts
	
	// Schedule Configuration
	IsScheduled   bool       `gorm:"default:false" json:"is_scheduled"`
	Schedule      string     `json:"schedule"`       // cron expression
	Recipients    string     `json:"recipients"`     // JSON array of email recipients
	FileFormat    string     `gorm:"default:'pdf'" json:"file_format"` // pdf, excel, csv, json
	
	// Template Configuration
	TemplateID    *uuid.UUID `gorm:"type:uuid" json:"template_id"`
	Layout        string     `json:"layout"`         // JSON configuration for layout
	Styling       string     `json:"styling"`        // JSON configuration for styling
	
	// Access Control
	IsPublic      bool       `gorm:"default:false" json:"is_public"`
	SharedWith    string     `json:"shared_with"`    // JSON array of user/role IDs
	
	// Performance Settings
	CacheEnabled  bool       `gorm:"default:true" json:"cache_enabled"`
	CacheTTL      int        `gorm:"default:300" json:"cache_ttl"` // seconds
	QueryTimeout  int        `gorm:"default:30" json:"query_timeout"` // seconds
	
	// Metadata
	Description   string     `json:"description"`
	Tags          string     `json:"tags"`           // JSON array
	Version       int        `gorm:"default:1" json:"version"`
	
	// Usage Statistics
	ViewCount     int        `json:"view_count"`
	LastViewed    *time.Time `json:"last_viewed"`
	ExecuteCount  int        `json:"execute_count"`
	LastExecuted  *time.Time `json:"last_executed"`
	AvgExecTime   float64    `json:"avg_exec_time"` // milliseconds
	
	// Timestamps
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	CreatedBy     uuid.UUID  `gorm:"type:uuid;not null" json:"created_by"`
	UpdatedBy     *uuid.UUID `gorm:"type:uuid" json:"updated_by"`
	
	// Relations
	Company       *Company        `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Template      *ReportTemplate `gorm:"foreignKey:TemplateID" json:"template,omitempty"`
	Creator       *User           `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	UpdatedByUser *User           `gorm:"foreignKey:UpdatedBy" json:"updated_by_user,omitempty"`
	Executions    []ReportExecution `gorm:"foreignKey:ReportID" json:"executions,omitempty"`
}

func (r *Report) BeforeCreate(tx *gorm.DB) error {
	r.ID = uuid.New()
	return nil
}

// ReportTemplate represents a predefined report template
type ReportTemplate struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID     *uuid.UUID `gorm:"type:uuid" json:"company_id"` // null for system templates
	Name          string     `gorm:"not null" json:"name"`
	NameEn        string     `json:"name_en"`
	Category      string     `gorm:"not null" json:"category"`
	Type          string     `gorm:"not null" json:"type"`
	IsSystemTemplate bool    `gorm:"default:false" json:"is_system_template"`
	
	// Template Configuration
	DataSource    string     `json:"data_source"`
	Filters       string     `json:"filters"`
	Columns       string     `json:"columns"`
	Sorting       string     `json:"sorting"`
	Grouping      string     `json:"grouping"`
	Aggregation   string     `json:"aggregation"`
	ChartConfig   string     `json:"chart_config"`
	Layout        string     `json:"layout"`
	Styling       string     `json:"styling"`
	
	// Template Properties
	Description   string     `json:"description"`
	Preview       string     `json:"preview"`        // base64 encoded preview image
	Tags          string     `json:"tags"`
	Industry      string     `json:"industry"`       // fastener, manufacturing, general
	Language      string     `gorm:"default:'zh-TW'" json:"language"`
	
	// Usage Statistics
	UsageCount    int        `json:"usage_count"`
	Rating        float64    `json:"rating"`         // 0-5
	RatingCount   int        `json:"rating_count"`
	
	// Timestamps
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	CreatedBy     *uuid.UUID `gorm:"type:uuid" json:"created_by"` // null for system templates
	
	// Relations
	Company       *Company   `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Creator       *User      `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

func (rt *ReportTemplate) BeforeCreate(tx *gorm.DB) error {
	rt.ID = uuid.New()
	return nil
}

// ReportExecution represents an execution instance of a report
type ReportExecution struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	ReportID      uuid.UUID  `gorm:"type:uuid;not null" json:"report_id"`
	Status        string     `gorm:"not null" json:"status"` // pending, running, completed, failed, cancelled
	
	// Execution Parameters
	Parameters    string     `json:"parameters"`     // JSON parameters used in execution
	FilePath      string     `json:"file_path"`      // path to generated file
	FileSize      int64      `json:"file_size"`      // bytes
	FileFormat    string     `json:"file_format"`
	
	// Execution Results
	RowCount      int        `json:"row_count"`
	ErrorMessage  string     `json:"error_message"`
	ExecutionTime float64    `json:"execution_time"` // milliseconds
	
	// Schedule Information
	IsScheduled   bool       `gorm:"default:false" json:"is_scheduled"`
	ScheduledAt   *time.Time `json:"scheduled_at"`
	TriggerType   string     `json:"trigger_type"`   // manual, scheduled, api
	
	// Access Information
	ExecutedBy    uuid.UUID  `gorm:"type:uuid;not null" json:"executed_by"`
	IPAddress     string     `json:"ip_address"`
	UserAgent     string     `json:"user_agent"`
	
	// Timestamps
	StartedAt     time.Time  `json:"started_at"`
	CompletedAt   *time.Time `json:"completed_at"`
	CreatedAt     time.Time  `json:"created_at"`
	
	// Relations
	Report        *Report    `gorm:"foreignKey:ReportID" json:"report,omitempty"`
	ExecutedByUser *User     `gorm:"foreignKey:ExecutedBy" json:"executed_by_user,omitempty"`
}

func (re *ReportExecution) BeforeCreate(tx *gorm.DB) error {
	re.ID = uuid.New()
	return nil
}

// ReportSubscription represents user subscriptions to scheduled reports
type ReportSubscription struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	ReportID      uuid.UUID  `gorm:"type:uuid;not null" json:"report_id"`
	UserID        uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	
	// Subscription Settings
	IsActive      bool       `gorm:"default:true" json:"is_active"`
	Email         string     `json:"email"`          // override user email
	Schedule      string     `json:"schedule"`       // override report schedule
	FileFormat    string     `json:"file_format"`    // override report format
	Parameters    string     `json:"parameters"`     // JSON parameters
	
	// Delivery Settings
	DeliveryMethod string    `gorm:"default:'email'" json:"delivery_method"` // email, webhook, ftp
	DeliveryConfig string    `json:"delivery_config"` // JSON configuration
	
	// Subscription Status
	LastDelivered *time.Time `json:"last_delivered"`
	DeliveryCount int        `json:"delivery_count"`
	FailureCount  int        `json:"failure_count"`
	LastError     string     `json:"last_error"`
	
	// Timestamps
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	
	// Relations
	Report        *Report    `gorm:"foreignKey:ReportID" json:"report,omitempty"`
	User          *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (rs *ReportSubscription) BeforeCreate(tx *gorm.DB) error {
	rs.ID = uuid.New()
	return nil
}

// ReportDashboard represents a collection of reports in dashboard format
type ReportDashboard struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID     uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	Name          string     `gorm:"not null" json:"name"`
	NameEn        string     `json:"name_en"`
	
	// Dashboard Configuration
	Layout        string     `json:"layout"`         // JSON configuration
	Theme         string     `gorm:"default:'light'" json:"theme"` // light, dark, auto
	RefreshRate   int        `gorm:"default:300" json:"refresh_rate"` // seconds
	
	// Dashboard Content
	Widgets       string     `json:"widgets"`        // JSON array of widget configurations
	Filters       string     `json:"filters"`        // JSON global filters
	
	// Access Control
	IsPublic      bool       `gorm:"default:false" json:"is_public"`
	SharedWith    string     `json:"shared_with"`    // JSON array of user/role IDs
	
	// Dashboard Properties
	Description   string     `json:"description"`
	Tags          string     `json:"tags"`
	IsDefault     bool       `gorm:"default:false" json:"is_default"`
	
	// Usage Statistics
	ViewCount     int        `json:"view_count"`
	LastViewed    *time.Time `json:"last_viewed"`
	
	// Timestamps
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	CreatedBy     uuid.UUID  `gorm:"type:uuid;not null" json:"created_by"`
	UpdatedBy     *uuid.UUID `gorm:"type:uuid" json:"updated_by"`
	
	// Relations
	Company       *Company   `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Creator       *User      `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	UpdatedByUser *User      `gorm:"foreignKey:UpdatedBy" json:"updated_by_user,omitempty"`
}

func (rd *ReportDashboard) BeforeCreate(tx *gorm.DB) error {
	rd.ID = uuid.New()
	return nil
}

// ReportDataSource represents external data sources for reports
type ReportDataSource struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID     uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	Name          string     `gorm:"not null" json:"name"`
	Type          string     `gorm:"not null" json:"type"` // database, api, file, internal
	
	// Connection Configuration
	ConnectionString string  `json:"connection_string"` // encrypted
	Credentials     string   `json:"credentials"`       // encrypted JSON
	Settings        string   `json:"settings"`          // JSON configuration
	
	// Data Source Properties
	Description     string   `json:"description"`
	Schema          string   `json:"schema"`            // JSON schema definition
	SampleData      string   `json:"sample_data"`       // JSON sample data
	
	// Connection Status
	Status          string   `gorm:"default:'active'" json:"status"` // active, inactive, error
	LastTested      *time.Time `json:"last_tested"`
	LastError       string   `json:"last_error"`
	TestResult      string   `json:"test_result"`       // JSON test results
	
	// Usage Statistics
	UsageCount      int      `json:"usage_count"`
	LastUsed        *time.Time `json:"last_used"`
	
	// Timestamps
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	CreatedBy       uuid.UUID  `gorm:"type:uuid;not null" json:"created_by"`
	
	// Relations
	Company         *Company   `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Creator         *User      `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

func (rds *ReportDataSource) BeforeCreate(tx *gorm.DB) error {
	rds.ID = uuid.New()
	return nil
}

// ReportSchedule represents scheduled report executions
type ReportSchedule struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	ReportID      uuid.UUID  `gorm:"type:uuid;not null" json:"report_id"`
	Name          string     `gorm:"not null" json:"name"`
	
	// Schedule Configuration
	CronExpression string    `gorm:"not null" json:"cron_expression"`
	Timezone      string     `gorm:"default:'Asia/Taipei'" json:"timezone"`
	IsActive      bool       `gorm:"default:true" json:"is_active"`
	
	// Execution Settings
	Parameters    string     `json:"parameters"`     // JSON parameters
	FileFormat    string     `gorm:"default:'pdf'" json:"file_format"`
	Recipients    string     `json:"recipients"`     // JSON array
	
	// Schedule Status
	NextRun       *time.Time `json:"next_run"`
	LastRun       *time.Time `json:"last_run"`
	LastStatus    string     `json:"last_status"`    // success, failed, cancelled
	RunCount      int        `json:"run_count"`
	FailureCount  int        `json:"failure_count"`
	LastError     string     `json:"last_error"`
	
	// Timestamps
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	CreatedBy     uuid.UUID  `gorm:"type:uuid;not null" json:"created_by"`
	
	// Relations
	Report        *Report    `gorm:"foreignKey:ReportID" json:"report,omitempty"`
	Creator       *User      `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

func (rs *ReportSchedule) BeforeCreate(tx *gorm.DB) error {
	rs.ID = uuid.New()
	return nil
}

// BusinessKPI represents key performance indicators for business metrics
type BusinessKPI struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID     uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	Name          string     `gorm:"not null" json:"name"`
	Category      string     `gorm:"not null" json:"category"` // sales, finance, production, inventory, quality
	
	// KPI Configuration
	Formula       string     `gorm:"not null" json:"formula"`      // calculation formula
	DataSources   string     `json:"data_sources"`     // JSON array of data sources
	Filters       string     `json:"filters"`          // JSON filters
	Unit          string     `json:"unit"`             // %, $, pcs, etc.
	
	// Target Settings
	TargetValue   float64    `json:"target_value"`
	TargetType    string     `json:"target_type"`      // minimum, maximum, exact, range
	ThresholdHigh float64    `json:"threshold_high"`   // warning threshold
	ThresholdLow  float64    `json:"threshold_low"`    // critical threshold
	
	// Display Settings
	DisplayFormat string     `json:"display_format"`   // number, percentage, currency
	ChartType     string     `json:"chart_type"`       // line, bar, gauge, sparkline
	ColorScheme   string     `json:"color_scheme"`     // JSON color configuration
	
	// KPI Properties
	Description   string     `json:"description"`
	Frequency     string     `gorm:"default:'daily'" json:"frequency"` // realtime, hourly, daily, weekly, monthly
	IsActive      bool       `gorm:"default:true" json:"is_active"`
	
	// Current Values
	CurrentValue  float64    `json:"current_value"`
	PreviousValue float64    `json:"previous_value"`
	Trend         string     `json:"trend"`            // up, down, stable
	LastUpdated   *time.Time `json:"last_updated"`
	
	// Timestamps
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	CreatedBy     uuid.UUID  `gorm:"type:uuid;not null" json:"created_by"`
	
	// Relations
	Company       *Company   `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Creator       *User      `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

func (bk *BusinessKPI) BeforeCreate(tx *gorm.DB) error {
	bk.ID = uuid.New()
	return nil
}