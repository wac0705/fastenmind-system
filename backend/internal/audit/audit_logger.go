package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/fastenmind/fastener-api/internal/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// EventType represents the type of audit event
type EventType string

const (
	EventTypeCreate         EventType = "CREATE"
	EventTypeUpdate         EventType = "UPDATE"
	EventTypeDelete         EventType = "DELETE"
	EventTypeAccess         EventType = "ACCESS"
	EventTypeLogin          EventType = "LOGIN"
	EventTypeLogout         EventType = "LOGOUT"
	EventTypeExport         EventType = "EXPORT"
	EventTypePermission     EventType = "PERMISSION"
	EventTypeConfiguration  EventType = "CONFIGURATION"
	EventTypeSecurity       EventType = "SECURITY"
)

// AuditLog represents an audit log entry
type AuditLog struct {
	ID              uuid.UUID              `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Timestamp       time.Time              `json:"timestamp"`
	EventType       EventType              `json:"event_type"`
	UserID          uuid.UUID              `json:"user_id"`
	Username        string                 `json:"username"`
	IPAddress       string                 `json:"ip_address"`
	UserAgent       string                 `json:"user_agent"`
	ResourceType    string                 `json:"resource_type"`
	ResourceID      string                 `json:"resource_id"`
	Action          string                 `json:"action"`
	Result          string                 `json:"result"`
	ErrorMessage    string                 `json:"error_message,omitempty"`
	OldValues       json.RawMessage        `gorm:"type:jsonb" json:"old_values,omitempty"`
	NewValues       json.RawMessage        `gorm:"type:jsonb" json:"new_values,omitempty"`
	Metadata        json.RawMessage        `gorm:"type:jsonb" json:"metadata,omitempty"`
	CorrelationID   string                 `json:"correlation_id,omitempty"`
	SessionID       string                 `json:"session_id,omitempty"`
	CompanyID       uuid.UUID              `json:"company_id"`
	Severity        string                 `json:"severity"`
	Tags            []string               `gorm:"type:text[]" json:"tags,omitempty"`
}

// AuditLogger provides audit logging functionality
type AuditLogger struct {
	db              *gorm.DB
	queue           chan *AuditLog
	bufferSize      int
	flushInterval   time.Duration
	encryptionKey   []byte
	wg              sync.WaitGroup
	done            chan bool
	sensitiveFields map[string]bool
}

// NewAuditLogger creates a new audit logger
func NewAuditLogger(db *gorm.DB, encryptionKey []byte) *AuditLogger {
	logger := &AuditLogger{
		db:            db,
		queue:         make(chan *AuditLog, 10000),
		bufferSize:    100,
		flushInterval: 5 * time.Second,
		encryptionKey: encryptionKey,
		done:          make(chan bool),
		sensitiveFields: map[string]bool{
			"password":     true,
			"ssn":          true,
			"credit_card":  true,
			"bank_account": true,
			"api_key":      true,
			"secret":       true,
		},
	}

	// Start background worker
	logger.wg.Add(1)
	go logger.worker()

	return logger
}

// LogEvent logs an audit event
func (al *AuditLogger) LogEvent(ctx context.Context, event *AuditLog) error {
	// Add context information
	if correlationID := ctx.Value("correlation_id"); correlationID != nil {
		event.CorrelationID = correlationID.(string)
	}
	if sessionID := ctx.Value("session_id"); sessionID != nil {
		event.SessionID = sessionID.(string)
	}

	// Sanitize sensitive data
	al.sanitizeSensitiveData(event)

	// Add to queue
	select {
	case al.queue <- event:
		return nil
	default:
		// Queue full, log synchronously
		return al.db.Create(event).Error
	}
}

// LogChange logs a data change event
func (al *AuditLogger) LogChange(ctx context.Context, resourceType string, resourceID string, oldValue, newValue interface{}, userInfo UserInfo) error {
	oldJSON, _ := json.Marshal(oldValue)
	newJSON, _ := json.Marshal(newValue)

	event := &AuditLog{
		Timestamp:    time.Now().UTC(),
		EventType:    EventTypeUpdate,
		UserID:       userInfo.UserID,
		Username:     userInfo.Username,
		IPAddress:    userInfo.IPAddress,
		UserAgent:    userInfo.UserAgent,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Action:       "UPDATE",
		Result:       "SUCCESS",
		OldValues:    oldJSON,
		NewValues:    newJSON,
		CompanyID:    userInfo.CompanyID,
		Severity:     "INFO",
	}

	return al.LogEvent(ctx, event)
}

// LogAccess logs a data access event
func (al *AuditLogger) LogAccess(ctx context.Context, resourceType string, resourceID string, action string, userInfo UserInfo) error {
	event := &AuditLog{
		Timestamp:    time.Now().UTC(),
		EventType:    EventTypeAccess,
		UserID:       userInfo.UserID,
		Username:     userInfo.Username,
		IPAddress:    userInfo.IPAddress,
		UserAgent:    userInfo.UserAgent,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Action:       action,
		Result:       "SUCCESS",
		CompanyID:    userInfo.CompanyID,
		Severity:     "INFO",
	}

	return al.LogEvent(ctx, event)
}

// LogSecurity logs a security event
func (al *AuditLogger) LogSecurity(ctx context.Context, action string, result string, details map[string]interface{}, userInfo UserInfo) error {
	metadata, _ := json.Marshal(details)

	severity := "WARNING"
	if result == "FAILURE" || result == "BLOCKED" {
		severity = "HIGH"
	}

	event := &AuditLog{
		Timestamp:  time.Now().UTC(),
		EventType:  EventTypeSecurity,
		UserID:     userInfo.UserID,
		Username:   userInfo.Username,
		IPAddress:  userInfo.IPAddress,
		UserAgent:  userInfo.UserAgent,
		Action:     action,
		Result:     result,
		Metadata:   metadata,
		CompanyID:  userInfo.CompanyID,
		Severity:   severity,
		Tags:       []string{"security", action},
	}

	return al.LogEvent(ctx, event)
}

// worker processes audit logs in background
func (al *AuditLogger) worker() {
	defer al.wg.Done()

	buffer := make([]*AuditLog, 0, al.bufferSize)
	ticker := time.NewTicker(al.flushInterval)
	defer ticker.Stop()

	for {
		select {
		case event := <-al.queue:
			buffer = append(buffer, event)
			if len(buffer) >= al.bufferSize {
				al.flush(buffer)
				buffer = buffer[:0]
			}

		case <-ticker.C:
			if len(buffer) > 0 {
				al.flush(buffer)
				buffer = buffer[:0]
			}

		case <-al.done:
			// Flush remaining events
			if len(buffer) > 0 {
				al.flush(buffer)
			}
			return
		}
	}
}

// flush writes buffered audit logs to database
func (al *AuditLogger) flush(logs []*AuditLog) {
	if len(logs) == 0 {
		return
	}

	// Batch insert
	if err := al.db.CreateInBatches(logs, 100).Error; err != nil {
		// Log error (avoid circular dependency)
		fmt.Printf("Failed to flush audit logs: %v\n", err)
	}
}

// sanitizeSensitiveData removes or masks sensitive information
func (al *AuditLogger) sanitizeSensitiveData(event *AuditLog) {
	// Sanitize old values
	if event.OldValues != nil {
		sanitized := al.sanitizeJSON(event.OldValues)
		event.OldValues = sanitized
	}

	// Sanitize new values
	if event.NewValues != nil {
		sanitized := al.sanitizeJSON(event.NewValues)
		event.NewValues = sanitized
	}

	// Sanitize metadata
	if event.Metadata != nil {
		sanitized := al.sanitizeJSON(event.Metadata)
		event.Metadata = sanitized
	}
}

// sanitizeJSON sanitizes sensitive fields in JSON data
func (al *AuditLogger) sanitizeJSON(data json.RawMessage) json.RawMessage {
	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return data
	}

	al.sanitizeMap(obj)

	sanitized, _ := json.Marshal(obj)
	return sanitized
}

// sanitizeMap recursively sanitizes sensitive fields in a map
func (al *AuditLogger) sanitizeMap(m map[string]interface{}) {
	for key, value := range m {
		if al.sensitiveFields[key] {
			m[key] = "***REDACTED***"
			continue
		}

		switch v := value.(type) {
		case map[string]interface{}:
			al.sanitizeMap(v)
		case []interface{}:
			for i, item := range v {
				if itemMap, ok := item.(map[string]interface{}); ok {
					al.sanitizeMap(itemMap)
					v[i] = itemMap
				}
			}
		}
	}
}

// Stop gracefully stops the audit logger
func (al *AuditLogger) Stop() {
	close(al.done)
	al.wg.Wait()
	close(al.queue)
}

// UserInfo contains user information for audit logging
type UserInfo struct {
	UserID    uuid.UUID
	Username  string
	IPAddress string
	UserAgent string
	CompanyID uuid.UUID
}

// AuditQuery provides audit log querying capabilities
type AuditQuery struct {
	db *gorm.DB
}

// NewAuditQuery creates a new audit query service
func NewAuditQuery(db *gorm.DB) *AuditQuery {
	return &AuditQuery{db: db}
}

// QueryOptions defines options for querying audit logs
type QueryOptions struct {
	StartTime    *time.Time
	EndTime      *time.Time
	UserID       *uuid.UUID
	EventType    *EventType
	ResourceType *string
	ResourceID   *string
	CompanyID    *uuid.UUID
	Severity     *string
	Limit        int
	Offset       int
}

// Query retrieves audit logs based on criteria
func (aq *AuditQuery) Query(opts QueryOptions) ([]*AuditLog, int64, error) {
	query := aq.db.Model(&AuditLog{})

	// Apply filters
	if opts.StartTime != nil {
		query = query.Where("timestamp >= ?", *opts.StartTime)
	}
	if opts.EndTime != nil {
		query = query.Where("timestamp <= ?", *opts.EndTime)
	}
	if opts.UserID != nil {
		query = query.Where("user_id = ?", *opts.UserID)
	}
	if opts.EventType != nil {
		query = query.Where("event_type = ?", *opts.EventType)
	}
	if opts.ResourceType != nil {
		query = query.Where("resource_type = ?", *opts.ResourceType)
	}
	if opts.ResourceID != nil {
		query = query.Where("resource_id = ?", *opts.ResourceID)
	}
	if opts.CompanyID != nil {
		query = query.Where("company_id = ?", *opts.CompanyID)
	}
	if opts.Severity != nil {
		query = query.Where("severity = ?", *opts.Severity)
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch with pagination
	var logs []*AuditLog
	err := query.
		Order("timestamp DESC").
		Limit(opts.Limit).
		Offset(opts.Offset).
		Find(&logs).Error

	return logs, total, err
}

// GetUserActivity retrieves user activity summary
func (aq *AuditQuery) GetUserActivity(userID uuid.UUID, days int) (map[EventType]int, error) {
	startTime := time.Now().AddDate(0, 0, -days)
	
	var results []struct {
		EventType EventType
		Count     int
	}

	err := aq.db.Model(&AuditLog{}).
		Select("event_type, COUNT(*) as count").
		Where("user_id = ? AND timestamp >= ?", userID, startTime).
		Group("event_type").
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	activity := make(map[EventType]int)
	for _, r := range results {
		activity[r.EventType] = r.Count
	}

	return activity, nil
}

// ComplianceReport generates compliance reports
type ComplianceReport struct {
	StartDate        time.Time
	EndDate          time.Time
	TotalEvents      int64
	EventsByType     map[EventType]int64
	SecurityEvents   int64
	FailedOperations int64
	UniqueUsers      int64
	TopUsers         []UserActivity
	SensitiveAccess  []SensitiveAccessLog
}

type UserActivity struct {
	UserID   uuid.UUID
	Username string
	Events   int64
}

type SensitiveAccessLog struct {
	Timestamp    time.Time
	UserID       uuid.UUID
	Username     string
	ResourceType string
	ResourceID   string
	Action       string
}

// GenerateComplianceReport generates a compliance report
func (aq *AuditQuery) GenerateComplianceReport(startDate, endDate time.Time, companyID uuid.UUID) (*ComplianceReport, error) {
	report := &ComplianceReport{
		StartDate:    startDate,
		EndDate:      endDate,
		EventsByType: make(map[EventType]int64),
	}

	// Total events
	aq.db.Model(&AuditLog{}).
		Where("timestamp BETWEEN ? AND ? AND company_id = ?", startDate, endDate, companyID).
		Count(&report.TotalEvents)

	// Events by type
	var typeResults []struct {
		EventType EventType
		Count     int64
	}
	aq.db.Model(&AuditLog{}).
		Select("event_type, COUNT(*) as count").
		Where("timestamp BETWEEN ? AND ? AND company_id = ?", startDate, endDate, companyID).
		Group("event_type").
		Find(&typeResults)

	for _, r := range typeResults {
		report.EventsByType[r.EventType] = r.Count
	}

	// Security events
	aq.db.Model(&AuditLog{}).
		Where("timestamp BETWEEN ? AND ? AND company_id = ? AND event_type = ?", 
			startDate, endDate, companyID, EventTypeSecurity).
		Count(&report.SecurityEvents)

	// Failed operations
	aq.db.Model(&AuditLog{}).
		Where("timestamp BETWEEN ? AND ? AND company_id = ? AND result != ?", 
			startDate, endDate, companyID, "SUCCESS").
		Count(&report.FailedOperations)

	// Unique users
	aq.db.Model(&AuditLog{}).
		Where("timestamp BETWEEN ? AND ? AND company_id = ?", startDate, endDate, companyID).
		Distinct("user_id").
		Count(&report.UniqueUsers)

	// Top users
	var topUsers []UserActivity
	aq.db.Model(&AuditLog{}).
		Select("user_id, username, COUNT(*) as events").
		Where("timestamp BETWEEN ? AND ? AND company_id = ?", startDate, endDate, companyID).
		Group("user_id, username").
		Order("events DESC").
		Limit(10).
		Find(&topUsers)
	report.TopUsers = topUsers

	return report, nil
}