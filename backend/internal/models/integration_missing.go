package models

import (
	"time"
	"github.com/google/uuid"
)

// Webhook represents a webhook configuration
type Webhook struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	IntegrationID   uuid.UUID  `json:"integration_id"`
	Name            string     `json:"name"`
	URL             string     `json:"url"`
	Method          string     `json:"method"`
	Headers         JSONB      `json:"headers" gorm:"type:jsonb"`
	EventTypes      JSONB      `json:"event_types" gorm:"type:jsonb"`
	IsActive        bool       `json:"is_active"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}