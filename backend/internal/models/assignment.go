package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// EngineerCapability 工程師能力矩陣
type EngineerCapability struct {
	ID                   uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	EngineerID           uuid.UUID `json:"engineer_id" gorm:"type:uuid;not null"`
	Engineer             Account   `json:"engineer" gorm:"foreignKey:EngineerID"`
	ProductCategory      string    `json:"product_category" gorm:"type:varchar(50);not null"`
	ProcessType          string    `json:"process_type" gorm:"type:varchar(50);not null"`
	SkillLevel           int       `json:"skill_level" gorm:"not null;check:skill_level >= 1 AND skill_level <= 5"`
	MaxConcurrentInquiries int     `json:"max_concurrent_inquiries" gorm:"default:5"`
	IsActive             bool      `json:"is_active" gorm:"default:true"`
	CreatedAt            time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt            time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// AssignmentRule 分派規則
type AssignmentRule struct {
	ID         uuid.UUID      `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	RuleName   string         `json:"rule_name" gorm:"type:varchar(100);not null"`
	RuleType   string         `json:"rule_type" gorm:"type:varchar(20);not null;check:rule_type IN ('auto','rotation','load_balance','skill_based')"`
	Priority   int            `json:"priority" gorm:"default:0"`
	Conditions datatypes.JSON `json:"conditions" gorm:"type:jsonb;not null"`
	IsActive   bool           `json:"is_active" gorm:"default:true"`
	CreatedAt  time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedBy  *uuid.UUID     `json:"created_by" gorm:"type:uuid"`
	UpdatedBy  *uuid.UUID     `json:"updated_by" gorm:"type:uuid"`
}

// EngineerWorkload 工程師工作負載
type EngineerWorkload struct {
	ID                     uuid.UUID  `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	EngineerID             uuid.UUID  `json:"engineer_id" gorm:"type:uuid;not null;uniqueIndex"`
	Engineer               Account    `json:"engineer" gorm:"foreignKey:EngineerID"`
	CurrentInquiries       int        `json:"current_inquiries" gorm:"default:0"`
	CompletedToday         int        `json:"completed_today" gorm:"default:0"`
	CompletedThisWeek      int        `json:"completed_this_week" gorm:"default:0"`
	CompletedThisMonth     int        `json:"completed_this_month" gorm:"default:0"`
	AverageCompletionHours float64    `json:"average_completion_hours" gorm:"type:decimal(10,2)"`
	LastAssignedAt         *time.Time `json:"last_assigned_at"`
	UpdatedAt              time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

// AssignmentHistory 分派歷史記錄
type AssignmentHistory struct {
	ID               uuid.UUID  `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	InquiryID        uuid.UUID  `json:"inquiry_id" gorm:"type:uuid;not null"`
	Inquiry          Inquiry    `json:"inquiry" gorm:"foreignKey:InquiryID"`
	AssignedFrom     *uuid.UUID `json:"assigned_from" gorm:"type:uuid"`
	AssignedFromUser *Account   `json:"assigned_from_user" gorm:"foreignKey:AssignedFrom"`
	AssignedTo       uuid.UUID  `json:"assigned_to" gorm:"type:uuid;not null"`
	AssignedToUser   Account    `json:"assigned_to_user" gorm:"foreignKey:AssignedTo"`
	AssignedBy       *uuid.UUID `json:"assigned_by" gorm:"type:uuid"`
	AssignedByUser   *Account   `json:"assigned_by_user" gorm:"foreignKey:AssignedBy"`
	AssignmentType   string     `json:"assignment_type" gorm:"type:varchar(20);not null;check:assignment_type IN ('auto','manual','reassign','self_select')"`
	AssignmentReason string     `json:"assignment_reason" gorm:"type:text"`
	RuleID           *uuid.UUID `json:"rule_id" gorm:"type:uuid"`
	Rule             *AssignmentRule `json:"rule" gorm:"foreignKey:RuleID"`
	AssignedAt       time.Time  `json:"assigned_at" gorm:"autoCreateTime"`
}

// EngineerPreference 工程師偏好設定
type EngineerPreference struct {
	ID                   uuid.UUID      `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	EngineerID           uuid.UUID      `json:"engineer_id" gorm:"type:uuid;not null;uniqueIndex"`
	Engineer             Account        `json:"engineer" gorm:"foreignKey:EngineerID"`
	PreferredCategories  datatypes.JSON `json:"preferred_categories" gorm:"type:jsonb"`
	PreferredCustomers   datatypes.JSON `json:"preferred_customers" gorm:"type:jsonb"`
	MaxDailyAssignments  int            `json:"max_daily_assignments" gorm:"default:10"`
	AutoAcceptEnabled    bool           `json:"auto_accept_enabled" gorm:"default:true"`
	NotificationEnabled  bool           `json:"notification_enabled" gorm:"default:true"`
	CreatedAt            time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt            time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (EngineerCapability) TableName() string {
	return "engineer_capabilities"
}

func (AssignmentRule) TableName() string {
	return "assignment_rules"
}

func (EngineerWorkload) TableName() string {
	return "engineer_workload"
}

func (AssignmentHistory) TableName() string {
	return "assignment_history"
}

func (EngineerPreference) TableName() string {
	return "engineer_preferences"
}

// AssignmentRequest 分派請求
type AssignmentRequest struct {
	InquiryID      uuid.UUID `json:"inquiry_id" binding:"required"`
	EngineerID     uuid.UUID `json:"engineer_id" binding:"required"`
	Reason         string    `json:"reason"`
	AssignmentType string    `json:"assignment_type" binding:"required,oneof=manual reassign"`
}

// AutoAssignmentRequest 自動分派請求
type AutoAssignmentRequest struct {
	InquiryID uuid.UUID `json:"inquiry_id" binding:"required"`
}

// EngineerWorkloadStats 工程師工作量統計
type EngineerWorkloadStats struct {
	EngineerID             uuid.UUID `json:"engineer_id"`
	EngineerName           string    `json:"engineer_name"`
	CurrentInquiries       int       `json:"current_inquiries"`
	CompletedToday         int       `json:"completed_today"`
	CompletedThisWeek      int       `json:"completed_this_week"`
	CompletedThisMonth     int       `json:"completed_this_month"`
	AverageCompletionHours float64   `json:"average_completion_hours"`
	SkillCategories        []string  `json:"skill_categories"`
}

// RuleCondition 規則條件結構
type RuleCondition struct {
	ProductCategories []string `json:"product_categories,omitempty"`
	CustomerTypes     []string `json:"customer_types,omitempty"`
	MinSkillLevel     int      `json:"min_skill_level,omitempty"`
	AutoAssign        bool     `json:"auto_assign,omitempty"`
	AmountRange       *struct {
		Min float64 `json:"min,omitempty"`
		Max float64 `json:"max,omitempty"`
	} `json:"amount_range,omitempty"`
}