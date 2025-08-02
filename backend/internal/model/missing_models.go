package model

import (
	"time"
	"github.com/google/uuid"
)

// Process represents a manufacturing process
type Process struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CompanyID   uuid.UUID  `json:"company_id" gorm:"type:uuid;not null"`
	Name        string     `json:"name" gorm:"type:varchar(255);not null"`
	Description string     `json:"description"`
	IsActive    bool       `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// Equipment represents manufacturing equipment
type Equipment struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CompanyID   uuid.UUID  `json:"company_id" gorm:"type:uuid;not null"`
	Name        string     `json:"name" gorm:"type:varchar(255);not null"`
	Code        string     `json:"code" gorm:"type:varchar(100);unique"`
	Description string     `json:"description"`
	Status      string     `json:"status" gorm:"type:varchar(50)"`
	IsActive    bool       `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// AssignmentRule represents an assignment rule
type AssignmentRule struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CompanyID   uuid.UUID  `json:"company_id" gorm:"type:uuid;not null"`
	Name        string     `json:"name" gorm:"type:varchar(255);not null"`
	Description string     `json:"description"`
	RuleType    string     `json:"rule_type" gorm:"type:varchar(50)"`
	Conditions  string     `json:"conditions" gorm:"type:jsonb"`
	Priority    int        `json:"priority"`
	IsActive    bool       `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

