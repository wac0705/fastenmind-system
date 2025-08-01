package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// N8NWorkflow represents a configured workflow in the system
type N8NWorkflow struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID     uuid.UUID      `gorm:"type:uuid;not null" json:"company_id"`
	Name          string         `gorm:"not null" json:"name"`
	Description   string         `json:"description,omitempty"`
	WorkflowID    string         `gorm:"not null" json:"workflow_id"` // N8N workflow ID
	TriggerType   string         `gorm:"not null" json:"trigger_type"` // webhook, schedule, manual, event
	TriggerConfig JSONB          `gorm:"type:jsonb" json:"trigger_config"`
	IsActive      bool           `gorm:"default:true" json:"is_active"`
	CreatedBy     uuid.UUID      `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	
	// Relations
	Company   *Company       `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Creator   *User          `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

func (w *N8NWorkflow) BeforeCreate(tx *gorm.DB) error {
	w.ID = uuid.New()
	return nil
}

// N8NExecution represents a workflow execution record
type N8NExecution struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID    uuid.UUID      `gorm:"type:uuid;not null" json:"company_id"`
	WorkflowID   string         `gorm:"not null;index" json:"workflow_id"`
	ExecutionID  string         `gorm:"not null;unique" json:"execution_id"` // N8N execution ID
	Status       string         `gorm:"not null" json:"status"` // running, success, error, canceled
	StartedAt    time.Time      `json:"started_at"`
	FinishedAt   *time.Time     `json:"finished_at,omitempty"`
	ErrorMessage string         `json:"error_message,omitempty"`
	InputData    JSONB          `gorm:"type:jsonb" json:"input_data,omitempty"`
	OutputData   JSONB          `gorm:"type:jsonb" json:"output_data,omitempty"`
	TriggeredBy  uuid.UUID      `gorm:"type:uuid" json:"triggered_by"`
	
	// Relations
	Company     *Company       `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Trigger     *User          `gorm:"foreignKey:TriggeredBy" json:"trigger,omitempty"`
}

func (e *N8NExecution) BeforeCreate(tx *gorm.DB) error {
	e.ID = uuid.New()
	return nil
}

// N8NWebhook represents a webhook registration
type N8NWebhook struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID   uuid.UUID      `gorm:"type:uuid;not null" json:"company_id"`
	Name        string         `gorm:"not null" json:"name"`
	WorkflowID  string         `gorm:"not null" json:"workflow_id"`
	EventTypes  []string       `gorm:"type:text[]" json:"event_types"`
	TargetURL   string         `json:"target_url,omitempty"`
	Headers     JSONB          `gorm:"type:jsonb" json:"headers,omitempty"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	
	// Relations
	Company     *Company       `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
}

func (w *N8NWebhook) BeforeCreate(tx *gorm.DB) error {
	w.ID = uuid.New()
	return nil
}

// N8NScheduledTask represents a scheduled workflow task
type N8NScheduledTask struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID      uuid.UUID      `gorm:"type:uuid;not null" json:"company_id"`
	Name           string         `gorm:"not null" json:"name"`
	WorkflowID     string         `gorm:"not null" json:"workflow_id"`
	CronExpression string         `gorm:"not null" json:"cron_expression"`
	Timezone       string         `gorm:"default:'UTC'" json:"timezone"`
	Data           JSONB          `gorm:"type:jsonb" json:"data,omitempty"`
	IsActive       bool           `gorm:"default:true" json:"is_active"`
	LastRunAt      *time.Time     `json:"last_run_at,omitempty"`
	NextRunAt      *time.Time     `json:"next_run_at,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	
	// Relations
	Company        *Company       `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
}

func (s *N8NScheduledTask) BeforeCreate(tx *gorm.DB) error {
	s.ID = uuid.New()
	return nil
}

// N8NEventLog represents system events that can trigger workflows
type N8NEventLog struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID    uuid.UUID      `gorm:"type:uuid;not null" json:"company_id"`
	EventType    string         `gorm:"not null;index" json:"event_type"`
	EntityType   string         `gorm:"not null" json:"entity_type"`
	EntityID     uuid.UUID      `gorm:"type:uuid;not null" json:"entity_id"`
	EventData    JSONB          `gorm:"type:jsonb" json:"event_data"`
	TriggeredBy  uuid.UUID      `gorm:"type:uuid;not null" json:"triggered_by"`
	WorkflowIDs  []string       `gorm:"type:text[]" json:"workflow_ids,omitempty"`
	ProcessedAt  *time.Time     `json:"processed_at,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	
	// Relations
	Company      *Company       `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	User         *User          `gorm:"foreignKey:TriggeredBy" json:"user,omitempty"`
}

func (e *N8NEventLog) BeforeCreate(tx *gorm.DB) error {
	e.ID = uuid.New()
	return nil
}

// JSONB is a custom type for PostgreSQL JSONB fields
type JSONB map[string]interface{}

// Value implements the driver.Valuer interface
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return gorm.ErrInvalidData
	}
	
	return json.Unmarshal(bytes, j)
}

// ConvertMapToJSONB converts map[string]string to JSONB
func ConvertMapToJSONB(m map[string]string) JSONB {
	result := make(JSONB)
	for k, v := range m {
		result[k] = v
	}
	return result
}