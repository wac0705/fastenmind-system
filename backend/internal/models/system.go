package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SystemConfig represents system configuration settings
type SystemConfig struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID     *uuid.UUID `gorm:"type:uuid" json:"company_id"` // null for global settings
	Category      string     `gorm:"not null" json:"category"`    // general, security, email, storage, integration, etc.
	Key           string     `gorm:"not null" json:"key"`
	Value         string     `json:"value"`
	DataType      string     `gorm:"not null" json:"data_type"`   // string, number, boolean, json, encrypted
	DefaultValue  string     `json:"default_value"`
	Description   string     `json:"description"`
	IsEditable    bool       `gorm:"default:true" json:"is_editable"`
	IsRequired    bool       `gorm:"default:false" json:"is_required"`
	ValidationRule string    `json:"validation_rule"`            // JSON validation rules
	
	// Display Settings
	DisplayName   string     `json:"display_name"`
	DisplayOrder  int        `json:"display_order"`
	GroupName     string     `json:"group_name"`
	InputType     string     `gorm:"default:'text'" json:"input_type"` // text, number, boolean, select, textarea, password, file
	Options       string     `json:"options"`                    // JSON options for select inputs
	
	// Timestamps
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	UpdatedBy     *uuid.UUID `gorm:"type:uuid" json:"updated_by"`
	
	// Relations
	Company       *Company   `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	UpdatedByUser *User      `gorm:"foreignKey:UpdatedBy" json:"updated_by_user,omitempty"`
}

func (sc *SystemConfig) BeforeCreate(tx *gorm.DB) error {
	sc.ID = uuid.New()
	return nil
}

// Role represents user roles in the system
type Role struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID     *uuid.UUID `gorm:"type:uuid" json:"company_id"` // null for system roles
	Name          string     `gorm:"not null" json:"name"`
	DisplayName   string     `gorm:"not null" json:"display_name"`
	Description   string     `json:"description"`
	Level         int        `gorm:"not null" json:"level"`       // 1=super_admin, 2=admin, 3=manager, 4=user
	IsSystemRole  bool       `gorm:"default:false" json:"is_system_role"`
	IsActive      bool       `gorm:"default:true" json:"is_active"`
	
	// Role Properties
	Color         string     `json:"color"`                       // hex color for UI display
	Icon          string     `json:"icon"`                        // icon name for UI display
	MaxUsers      int        `json:"max_users"`                   // maximum users with this role (0 = unlimited)
	
	// Timestamps
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	CreatedBy     *uuid.UUID `gorm:"type:uuid" json:"created_by"`
	
	// Relations
	Company       *Company        `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Creator       *User           `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Permissions   []RolePermission `gorm:"foreignKey:RoleID" json:"permissions,omitempty"`
	Users         []User          `gorm:"foreignKey:Role" json:"users,omitempty"`
}

func (r *Role) BeforeCreate(tx *gorm.DB) error {
	r.ID = uuid.New()
	return nil
}

// Permission represents system permissions
type Permission struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	Module        string     `gorm:"not null" json:"module"`      // inquiry, quote, order, inventory, etc.
	Action        string     `gorm:"not null" json:"action"`      // create, read, update, delete, approve, export, etc.
	Resource      string     `json:"resource"`                    // specific resource within module
	Name          string     `gorm:"not null" json:"name"`        // unique permission name
	DisplayName   string     `gorm:"not null" json:"display_name"`
	Description   string     `json:"description"`
	Category      string     `json:"category"`                    // core, advanced, admin
	IsSystemPerm  bool       `gorm:"default:true" json:"is_system_perm"`
	
	// Permission Properties
	RequiredLevel int        `gorm:"default:1" json:"required_level"` // minimum role level required
	Conditions    string     `json:"conditions"`                  // JSON conditions for dynamic permissions
	
	// Timestamps
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	
	// Relations
	RolePermissions []RolePermission `gorm:"foreignKey:PermissionID" json:"role_permissions,omitempty"`
}

func (p *Permission) BeforeCreate(tx *gorm.DB) error {
	p.ID = uuid.New()
	return nil
}

// RolePermission represents the mapping between roles and permissions
type RolePermission struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	RoleID       uuid.UUID  `gorm:"type:uuid;not null" json:"role_id"`
	PermissionID uuid.UUID  `gorm:"type:uuid;not null" json:"permission_id"`
	IsGranted    bool       `gorm:"default:true" json:"is_granted"`
	Constraints  string     `json:"constraints"`                 // JSON constraints for conditional permissions
	
	// Timestamps
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	GrantedBy    *uuid.UUID `gorm:"type:uuid" json:"granted_by"`
	
	// Relations
	Role         *Role       `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Permission   *Permission `gorm:"foreignKey:PermissionID" json:"permission,omitempty"`
	GrantedByUser *User      `gorm:"foreignKey:GrantedBy" json:"granted_by_user,omitempty"`
}

func (rp *RolePermission) BeforeCreate(tx *gorm.DB) error {
	rp.ID = uuid.New()
	return nil
}

// UserSession represents active user sessions
type UserSession struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	UserID        uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	SessionToken  string     `gorm:"not null;unique" json:"session_token"`
	RefreshToken  string     `gorm:"not null;unique" json:"refresh_token"`
	
	// Session Information
	IPAddress     string     `json:"ip_address"`
	UserAgent     string     `json:"user_agent"`
	DeviceInfo    string     `json:"device_info"`                // JSON device information
	Location      string     `json:"location"`                   // JSON location information
	
	// Session Status
	IsActive      bool       `gorm:"default:true" json:"is_active"`
	LastActivity  time.Time  `json:"last_activity"`
	ExpiresAt     time.Time  `json:"expires_at"`
	
	// Security Information
	LoginMethod   string     `json:"login_method"`               // password, sso, api_key, etc.
	TwoFactorUsed bool       `gorm:"default:false" json:"two_factor_used"`
	RiskScore     float64    `json:"risk_score"`                 // 0-100, higher = more risky
	
	// Timestamps
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	
	// Relations
	User          *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (us *UserSession) BeforeCreate(tx *gorm.DB) error {
	us.ID = uuid.New()
	return nil
}

// AuditLog represents system audit trails
type AuditLog struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID     uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	UserID        *uuid.UUID `gorm:"type:uuid" json:"user_id"`
	SessionID     *uuid.UUID `gorm:"type:uuid" json:"session_id"`
	
	// Action Information
	Action        string     `gorm:"not null" json:"action"`     // create, update, delete, login, logout, etc.
	Resource      string     `gorm:"not null" json:"resource"`   // table name or resource type
	ResourceID    *uuid.UUID `gorm:"type:uuid" json:"resource_id"`
	Module        string     `json:"module"`                     // inquiry, quote, order, etc.
	
	// Change Information
	OldValues     string     `json:"old_values"`                 // JSON of old values
	NewValues     string     `json:"new_values"`                 // JSON of new values
	Changes       string     `json:"changes"`                    // JSON summary of changes
	
	// Request Information
	IPAddress     string     `json:"ip_address"`
	UserAgent     string     `json:"user_agent"`
	RequestMethod string     `json:"request_method"`
	RequestPath   string     `json:"request_path"`
	RequestBody   string     `json:"request_body"`
	ResponseCode  int        `json:"response_code"`
	
	// Additional Information
	Description   string     `json:"description"`
	Severity      string     `gorm:"default:'info'" json:"severity"` // info, warning, error, critical
	Tags          string     `json:"tags"`                       // JSON array of tags
	
	// Timestamps
	Timestamp     time.Time  `json:"timestamp"`
	CreatedAt     time.Time  `json:"created_at"`
	
	// Relations
	Company       *Company    `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	User          *User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Session       *UserSession `gorm:"foreignKey:SessionID" json:"session,omitempty"`
}

func (al *AuditLog) BeforeCreate(tx *gorm.DB) error {
	al.ID = uuid.New()
	return nil
}

// SystemNotification represents system-wide notifications
type SystemNotification struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID     *uuid.UUID `gorm:"type:uuid" json:"company_id"` // null for global notifications
	Title         string     `gorm:"not null" json:"title"`
	Message       string     `gorm:"not null" json:"message"`
	Type          string     `gorm:"not null" json:"type"`        // info, warning, error, success, maintenance
	Priority      string     `gorm:"default:'normal'" json:"priority"` // low, normal, high, urgent
	
	// Targeting
	TargetRoles   string     `json:"target_roles"`               // JSON array of role names
	TargetUsers   string     `json:"target_users"`               // JSON array of user IDs
	
	// Notification Properties
	IsActive      bool       `gorm:"default:true" json:"is_active"`
	IsPersistent  bool       `gorm:"default:false" json:"is_persistent"` // stays until dismissed
	AutoDismiss   int        `json:"auto_dismiss"`               // seconds until auto dismiss (0 = no auto dismiss)
	
	// Display Settings
	Icon          string     `json:"icon"`
	Color         string     `json:"color"`
	ActionLabel   string     `json:"action_label"`
	ActionURL     string     `json:"action_url"`
	
	// Scheduling
	ShowFrom      *time.Time `json:"show_from"`
	ShowUntil     *time.Time `json:"show_until"`
	
	// Statistics
	ViewCount     int        `json:"view_count"`
	ClickCount    int        `json:"click_count"`
	DismissCount  int        `json:"dismiss_count"`
	
	// Timestamps
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	CreatedBy     uuid.UUID  `gorm:"type:uuid;not null" json:"created_by"`
	
	// Relations
	Company       *Company   `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Creator       *User      `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

func (sn *SystemNotification) BeforeCreate(tx *gorm.DB) error {
	sn.ID = uuid.New()
	return nil
}

// UserNotification represents individual user notifications
type UserNotification struct {
	ID                   uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	UserID               uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	SystemNotificationID *uuid.UUID `gorm:"type:uuid" json:"system_notification_id"`
	
	// Notification Content
	Title         string     `json:"title"`
	Message       string     `json:"message"`
	Type          string     `json:"type"`
	Priority      string     `json:"priority"`
	
	// Notification Properties
	IsRead        bool       `gorm:"default:false" json:"is_read"`
	IsDismissed   bool       `gorm:"default:false" json:"is_dismissed"`
	IsActionable  bool       `gorm:"default:false" json:"is_actionable"`
	
	// Action Information
	ActionLabel   string     `json:"action_label"`
	ActionURL     string     `json:"action_url"`
	ActionData    string     `json:"action_data"`                // JSON action data
	
	// Tracking
	ReadAt        *time.Time `json:"read_at"`
	DismissedAt   *time.Time `json:"dismissed_at"`
	ClickedAt     *time.Time `json:"clicked_at"`
	
	// Timestamps
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	
	// Relations
	User               *User               `gorm:"foreignKey:UserID" json:"user,omitempty"`
	SystemNotification *SystemNotification `gorm:"foreignKey:SystemNotificationID" json:"system_notification,omitempty"`
}

func (un *UserNotification) BeforeCreate(tx *gorm.DB) error {
	un.ID = uuid.New()
	return nil
}

// SystemHealth represents system health monitoring
type SystemHealth struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID     *uuid.UUID `gorm:"type:uuid" json:"company_id"` // null for global health
	Component     string     `gorm:"not null" json:"component"`   // database, redis, email, storage, api, etc.
	Status        string     `gorm:"not null" json:"status"`      // healthy, warning, critical, down
	
	// Health Metrics
	ResponseTime  float64    `json:"response_time"`              // milliseconds
	ErrorRate     float64    `json:"error_rate"`                 // percentage
	Uptime        float64    `json:"uptime"`                     // percentage
	
	// Resource Usage
	CPUUsage      float64    `json:"cpu_usage"`                  // percentage
	MemoryUsage   float64    `json:"memory_usage"`               // percentage
	DiskUsage     float64    `json:"disk_usage"`                 // percentage
	NetworkIn     int64      `json:"network_in"`                 // bytes
	NetworkOut    int64      `json:"network_out"`                // bytes
	
	// Health Details
	Message       string     `json:"message"`
	Details       string     `json:"details"`                    // JSON additional details
	Metrics       string     `json:"metrics"`                    // JSON metrics data
	
	// Check Information
	CheckedAt     time.Time  `json:"checked_at"`
	CheckDuration float64    `json:"check_duration"`             // milliseconds
	
	// Timestamps
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	
	// Relations
	Company       *Company   `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
}

func (sh *SystemHealth) BeforeCreate(tx *gorm.DB) error {
	sh.ID = uuid.New()
	return nil
}

// BackupRecord represents system backup records
// Commented out - duplicate definition exists in advanced.go
/*
type BackupRecord struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID     *uuid.UUID `gorm:"type:uuid" json:"company_id"` // null for full system backup
	Type          string     `gorm:"not null" json:"type"`        // full, incremental, differential
	Status        string     `gorm:"not null" json:"status"`      // running, completed, failed, cancelled
	
	// Backup Information
	BackupName    string     `json:"backup_name"`
	Description   string     `json:"description"`
	FilePath      string     `json:"file_path"`
	FileSize      int64      `json:"file_size"`                  // bytes
	Checksum      string     `json:"checksum"`
	
	// Backup Scope
	Tables        string     `json:"tables"`                     // JSON array of table names
	DataOnly      bool       `gorm:"default:false" json:"data_only"`
	SchemaOnly    bool       `gorm:"default:false" json:"schema_only"`
	
	// Timing Information
	StartedAt     time.Time  `json:"started_at"`
	CompletedAt   *time.Time `json:"completed_at"`
	Duration      float64    `json:"duration"`                   // seconds
	
	// Retention
	RetentionDays int        `gorm:"default:30" json:"retention_days"`
	ExpiresAt     time.Time  `json:"expires_at"`
	
	// Error Information
	ErrorMessage  string     `json:"error_message"`
	ErrorDetails  string     `json:"error_details"`
	
	// Timestamps
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	CreatedBy     *uuid.UUID `gorm:"type:uuid" json:"created_by"`
	
	// Relations
	Company       *Company   `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Creator       *User      `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}
*/

// Commented out - using BackupRecord from advanced.go
/*
func (br *BackupRecord) BeforeCreate(tx *gorm.DB) error {
	br.ID = uuid.New()
	return nil
}
*/

// SystemTask represents background system tasks
type SystemTask struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID     *uuid.UUID `gorm:"type:uuid" json:"company_id"`
	Name          string     `gorm:"not null" json:"name"`
	Type          string     `gorm:"not null" json:"type"`        // backup, cleanup, migration, report, etc.
	Status        string     `gorm:"not null" json:"status"`      // pending, running, completed, failed, cancelled
	Priority      string     `gorm:"default:'normal'" json:"priority"` // low, normal, high, urgent
	
	// Task Information
	Description   string     `json:"description"`
	Parameters    string     `json:"parameters"`                 // JSON task parameters
	Result        string     `json:"result"`                     // JSON task result
	
	// Progress Information
	Progress      float64    `json:"progress"`                   // 0-100 percentage
	CurrentStep   string     `json:"current_step"`
	TotalSteps    int        `json:"total_steps"`
	
	// Timing Information
	ScheduledAt   *time.Time `json:"scheduled_at"`
	StartedAt     *time.Time `json:"started_at"`
	CompletedAt   *time.Time `json:"completed_at"`
	Duration      float64    `json:"duration"`                   // seconds
	
	// Retry Information
	MaxRetries    int        `gorm:"default:3" json:"max_retries"`
	RetryCount    int        `json:"retry_count"`
	NextRetryAt   *time.Time `json:"next_retry_at"`
	
	// Error Information
	ErrorMessage  string     `json:"error_message"`
	ErrorDetails  string     `json:"error_details"`
	
	// Timestamps
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	CreatedBy     *uuid.UUID `gorm:"type:uuid" json:"created_by"`
	
	// Relations
	Company       *Company   `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Creator       *User      `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

func (st *SystemTask) BeforeCreate(tx *gorm.DB) error {
	st.ID = uuid.New()
	return nil
}