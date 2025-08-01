package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MobileDevice represents registered mobile devices
type MobileDevice struct {
	ID               uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	UserID           uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	CompanyID        uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	DeviceToken      string     `gorm:"not null;unique" json:"device_token"`
	Platform         string     `gorm:"not null" json:"platform"`         // ios, android, web
	DeviceType       string     `gorm:"not null" json:"device_type"`       // phone, tablet, desktop
	DeviceModel      string     `json:"device_model"`
	OSVersion        string     `json:"os_version"`
	AppVersion       string     `json:"app_version"`
	
	// Device Information
	DeviceName       string     `json:"device_name"`
	DeviceID         string     `json:"device_id"`                         // unique device identifier
	BundleID         string     `json:"bundle_id"`                         // app bundle identifier
	TimeZone         string     `json:"time_zone"`
	Language         string     `json:"language"`
	Country          string     `json:"country"`
	
	// Push Notification Settings
	PushEnabled      bool       `gorm:"default:true" json:"push_enabled"`
	BadgeCount       int        `json:"badge_count"`
	NotificationTypes string    `json:"notification_types"`               // JSON array of enabled notification types
	
	// Device Status
	IsActive         bool       `gorm:"default:true" json:"is_active"`
	LastSeen         time.Time  `json:"last_seen"`
	LastLocation     string     `json:"last_location"`                     // JSON location data
	
	// Security
	IsJailbroken     bool       `gorm:"default:false" json:"is_jailbroken"`
	SecurityLevel    string     `gorm:"default:'normal'" json:"security_level"` // low, normal, high
	
	// Timestamps
	RegisteredAt     time.Time  `json:"registered_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	
	// Relations
	User             *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Company          *Company   `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	PushNotifications []PushNotification `gorm:"foreignKey:DeviceID" json:"push_notifications,omitempty"`
}

func (md *MobileDevice) BeforeCreate(tx *gorm.DB) error {
	md.ID = uuid.New()
	return nil
}

// PushNotification represents push notifications sent to mobile devices
type PushNotification struct {
	ID               uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	DeviceID         uuid.UUID  `gorm:"type:uuid;not null" json:"device_id"`
	UserID           uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	CompanyID        uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	
	// Notification Content
	Title            string     `gorm:"not null" json:"title"`
	Body             string     `gorm:"not null" json:"body"`
	Type             string     `gorm:"not null" json:"type"`              // inquiry, quote, order, system, etc.
	Category         string     `json:"category"`                          // notification category for grouping
	
	// Notification Data
	Data             string     `json:"data"`                              // JSON payload data
	ResourceID       *uuid.UUID `gorm:"type:uuid" json:"resource_id"`
	ResourceType     string     `json:"resource_type"`                     // inquiry, quote, order, etc.
	
	// Notification Settings
	Priority         string     `gorm:"default:'normal'" json:"priority"`  // low, normal, high, urgent
	Badge            int        `json:"badge"`                             // badge count increment
	Sound            string     `json:"sound"`                             // notification sound
	Icon             string     `json:"icon"`                              // notification icon
	Image            string     `json:"image"`                             // notification image URL
	
	// Action Buttons
	Actions          string     `json:"actions"`                           // JSON array of action buttons
	
	// Delivery Status
	Status           string     `gorm:"not null" json:"status"`            // pending, sent, delivered, failed, clicked
	SentAt           *time.Time `json:"sent_at"`
	DeliveredAt      *time.Time `json:"delivered_at"`
	ClickedAt        *time.Time `json:"clicked_at"`
	
	// Delivery Information
	ProviderResponse string     `json:"provider_response"`                 // FCM/APNS response
	ErrorMessage     string     `json:"error_message"`
	RetryCount       int        `json:"retry_count"`
	MaxRetries       int        `gorm:"default:3" json:"max_retries"`
	
	// Expiration
	ExpiresAt        *time.Time `json:"expires_at"`
	TTL              int        `gorm:"default:86400" json:"ttl"`          // time to live in seconds
	
	// Timestamps
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	ScheduledAt      *time.Time `json:"scheduled_at"`
	
	// Relations
	Device           *MobileDevice `gorm:"foreignKey:DeviceID" json:"device,omitempty"`
	User             *User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Company          *Company      `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
}

func (pn *PushNotification) BeforeCreate(tx *gorm.DB) error {
	pn.ID = uuid.New()
	return nil
}

// MobileSession represents mobile app sessions
type MobileSession struct {
	ID               uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	DeviceID         uuid.UUID  `gorm:"type:uuid;not null" json:"device_id"`
	UserID           uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	CompanyID        uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	
	// Session Information
	SessionToken     string     `gorm:"not null;unique" json:"session_token"`
	RefreshToken     string     `gorm:"not null;unique" json:"refresh_token"`
	
	// Session Details
	StartTime        time.Time  `json:"start_time"`
	EndTime          *time.Time `json:"end_time"`
	Duration         int64      `json:"duration"`                          // session duration in seconds
	
	// Location Information
	StartLocation    string     `json:"start_location"`                    // JSON location data
	EndLocation      string     `json:"end_location"`                      // JSON location data
	
	// Usage Statistics
	ScreenViews      int        `json:"screen_views"`
	APIRequests      int        `json:"api_requests"`
	DataTransferred  int64      `json:"data_transferred"`                  // bytes
	
	// App State
	AppState         string     `json:"app_state"`                         // active, background, terminated
	LastActivity     time.Time  `json:"last_activity"`
	
	// Network Information
	NetworkType      string     `json:"network_type"`                      // wifi, cellular, offline
	NetworkProvider  string     `json:"network_provider"`
	
	// Timestamps
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	
	// Relations
	Device           *MobileDevice `gorm:"foreignKey:DeviceID" json:"device,omitempty"`
	User             *User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Company          *Company      `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
}

func (ms *MobileSession) BeforeCreate(tx *gorm.DB) error {
	ms.ID = uuid.New()
	return nil
}

// MobileAnalytics represents mobile app analytics data
type MobileAnalytics struct {
	ID               uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	DeviceID         uuid.UUID  `gorm:"type:uuid;not null" json:"device_id"`
	UserID           uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	CompanyID        uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	
	// Event Information
	EventType        string     `gorm:"not null" json:"event_type"`        // screen_view, button_click, api_call, etc.
	EventName        string     `gorm:"not null" json:"event_name"`
	EventCategory    string     `json:"event_category"`
	
	// Event Data
	EventData        string     `json:"event_data"`                        // JSON event data
	ScreenName       string     `json:"screen_name"`
	ScreenClass      string     `json:"screen_class"`
	
	// Timing Information
	EventTimestamp   time.Time  `json:"event_timestamp"`
	SessionID        uuid.UUID  `gorm:"type:uuid" json:"session_id"`
	Duration         int64      `json:"duration"`                          // event duration in milliseconds
	
	// User Interaction
	InteractionType  string     `json:"interaction_type"`                  // tap, swipe, scroll, etc.
	ElementID        string     `json:"element_id"`
	ElementType      string     `json:"element_type"`
	
	// Performance Data
	LoadTime         int64      `json:"load_time"`                         // milliseconds
	ResponseTime     int64      `json:"response_time"`                     // milliseconds
	ErrorMessage     string     `json:"error_message"`
	
	// Context Information
	AppVersion       string     `json:"app_version"`
	OSVersion        string     `json:"os_version"`
	NetworkType      string     `json:"network_type"`
	BatteryLevel     float64    `json:"battery_level"`
	MemoryUsage      int64      `json:"memory_usage"`                      // bytes
	
	// Location (if enabled)
	Latitude         float64    `json:"latitude"`
	Longitude        float64    `json:"longitude"`
	LocationAccuracy float64    `json:"location_accuracy"`
	
	// Timestamps
	CreatedAt        time.Time  `json:"created_at"`
	
	// Relations
	Device           *MobileDevice `gorm:"foreignKey:DeviceID" json:"device,omitempty"`
	User             *User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Company          *Company      `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Session          *MobileSession `gorm:"foreignKey:SessionID" json:"session,omitempty"`
}

func (ma *MobileAnalytics) BeforeCreate(tx *gorm.DB) error {
	ma.ID = uuid.New()
	return nil
}

// MobileAppVersion represents mobile app versions and updates
type MobileAppVersion struct {
	ID               uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID        *uuid.UUID `gorm:"type:uuid" json:"company_id"`
	
	// Version Information
	Version          string     `gorm:"not null" json:"version"`           // 1.0.0
	BuildNumber      string     `gorm:"not null" json:"build_number"`      // 100
	Platform         string     `gorm:"not null" json:"platform"`          // ios, android, web
	
	// Release Information
	ReleaseType      string     `gorm:"not null" json:"release_type"`      // alpha, beta, production
	ReleaseNotes     string     `json:"release_notes"`                     // markdown format
	ReleaseNotesEn   string     `json:"release_notes_en"`
	
	// Version Status
	Status           string     `gorm:"not null" json:"status"`            // draft, review, approved, released, deprecated
	IsForceUpdate    bool       `gorm:"default:false" json:"is_force_update"`
	IsActive         bool       `gorm:"default:true" json:"is_active"`
	
	// Download Information
	DownloadURL      string     `json:"download_url"`
	FileSize         int64      `json:"file_size"`                         // bytes
	Checksum         string     `json:"checksum"`                          // SHA256 checksum
	
	// Compatibility
	MinOSVersion     string     `json:"min_os_version"`
	RequiredFeatures string     `json:"required_features"`                 // JSON array of required features
	
	// Statistics
	InstallCount     int64      `json:"install_count"`
	UpdateCount      int64      `json:"update_count"`
	CrashCount       int64      `json:"crash_count"`
	RatingAverage    float64    `json:"rating_average"`
	RatingCount      int64      `json:"rating_count"`
	
	// Rollout
	RolloutPercent   int        `gorm:"default:100" json:"rollout_percent"` // 0-100
	RolloutRegions   string     `json:"rollout_regions"`                   // JSON array of regions
	
	// Timestamps
	ReleasedAt       *time.Time `json:"released_at"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	CreatedBy        uuid.UUID  `gorm:"type:uuid;not null" json:"created_by"`
	
	// Relations
	Company          *Company   `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Creator          *User      `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

func (mav *MobileAppVersion) BeforeCreate(tx *gorm.DB) error {
	mav.ID = uuid.New()
	return nil
}

// MobileOfflineData represents offline data synchronization
type MobileOfflineData struct {
	ID               uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	DeviceID         uuid.UUID  `gorm:"type:uuid;not null" json:"device_id"`
	UserID           uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	CompanyID        uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	
	// Data Information
	DataType         string     `gorm:"not null" json:"data_type"`         // inquiry, quote, order, etc.
	ResourceID       uuid.UUID  `gorm:"type:uuid;not null" json:"resource_id"`
	Operation        string     `gorm:"not null" json:"operation"`         // create, update, delete
	
	// Data Content
	DataPayload      string     `json:"data_payload"`                      // JSON data
	DataChecksum     string     `json:"data_checksum"`                     // data integrity check
	
	// Sync Status
	Status           string     `gorm:"not null" json:"status"`            // pending, syncing, synced, failed, conflict
	Priority         int        `gorm:"default:1" json:"priority"`         // 1=low, 2=normal, 3=high
	
	// Conflict Resolution
	ConflictData     string     `json:"conflict_data"`                     // JSON conflict information
	ResolvedBy       *uuid.UUID `gorm:"type:uuid" json:"resolved_by"`
	ResolutionMethod string     `json:"resolution_method"`                 // server_wins, client_wins, merge
	
	// Sync Information
	LastSyncAttempt  *time.Time `json:"last_sync_attempt"`
	SyncedAt         *time.Time `json:"synced_at"`
	ErrorMessage     string     `json:"error_message"`
	RetryCount       int        `json:"retry_count"`
	MaxRetries       int        `gorm:"default:5" json:"max_retries"`
	
	// Timestamps
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	ExpiresAt        *time.Time `json:"expires_at"`
	
	// Relations
	Device           *MobileDevice `gorm:"foreignKey:DeviceID" json:"device,omitempty"`
	User             *User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Company          *Company      `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	ResolvedByUser   *User         `gorm:"foreignKey:ResolvedBy" json:"resolved_by_user,omitempty"`
}

func (mod *MobileOfflineData) BeforeCreate(tx *gorm.DB) error {
	mod.ID = uuid.New()
	return nil
}

// MobileConfiguration represents mobile app configuration
type MobileConfiguration struct {
	ID               uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID        *uuid.UUID `gorm:"type:uuid" json:"company_id"`
	Platform         string     `gorm:"not null" json:"platform"`          // ios, android, web, all
	
	// Configuration Information
	ConfigKey        string     `gorm:"not null" json:"config_key"`
	ConfigValue      string     `json:"config_value"`                      // JSON configuration value
	ConfigType       string     `gorm:"not null" json:"config_type"`       // string, number, boolean, object, array
	
	// Configuration Metadata
	Category         string     `json:"category"`                          // ui, security, features, etc.
	Description      string     `json:"description"`
	DefaultValue     string     `json:"default_value"`
	
	// Validation
	ValidationRules  string     `json:"validation_rules"`                  // JSON validation rules
	MinVersion       string     `json:"min_version"`                       // minimum app version required
	MaxVersion       string     `json:"max_version"`                       // maximum app version supported
	
	// Rollout
	IsEnabled        bool       `gorm:"default:true" json:"is_enabled"`
	RolloutPercent   int        `gorm:"default:100" json:"rollout_percent"` // 0-100
	TargetDevices    string     `json:"target_devices"`                    // JSON array of device criteria
	
	// Cache Settings
	CacheTTL         int        `gorm:"default:3600" json:"cache_ttl"`     // cache time to live in seconds
	IsSecure         bool       `gorm:"default:false" json:"is_secure"`    // encrypt configuration value
	
	// Timestamps
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	CreatedBy        uuid.UUID  `gorm:"type:uuid;not null" json:"created_by"`
	
	// Relations
	Company          *Company   `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Creator          *User      `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

func (mc *MobileConfiguration) BeforeCreate(tx *gorm.DB) error {
	mc.ID = uuid.New()
	return nil
}