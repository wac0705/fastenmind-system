package models

import (
	"time"
)

// EngineerAssignment 工程師分派記錄
type EngineerAssignment struct {
	ID          string     `json:"id" gorm:"primaryKey"`
	CompanyID   string     `json:"company_id" gorm:"index"`
	InquiryID   string     `json:"inquiry_id" gorm:"index"`
	EngineerID  string     `json:"engineer_id" gorm:"index"`
	AssignedBy  string     `json:"assigned_by"`
	AssignedAt  time.Time  `json:"assigned_at"`
	Status      string     `json:"status"` // pending, in_progress, completed, cancelled
	Priority    string     `json:"priority"` // low, normal, high, urgent
	DueDate     *time.Time `json:"due_date"`
	CompletedAt *time.Time `json:"completed_at"`
	Notes       string     `json:"notes"`
	CreatedBy   string     `json:"created_by"`
	UpdatedBy   string     `json:"updated_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at" gorm:"index"`
	
	// Relations
	Company  Company `json:"company,omitempty" gorm:"foreignKey:CompanyID"`
	Inquiry  Inquiry `json:"inquiry,omitempty" gorm:"foreignKey:InquiryID"`
	Engineer Account `json:"engineer,omitempty" gorm:"foreignKey:EngineerID"`
	Assigner Account `json:"assigner,omitempty" gorm:"foreignKey:AssignedBy"`
}

// EngineerAssignmentHistory 工程師分派歷史記錄
type EngineerAssignmentHistory struct {
	ID           string    `json:"id" gorm:"primaryKey"`
	AssignmentID string    `json:"assignment_id" gorm:"index"`
	Action       string    `json:"action"` // assigned, reassigned, status_changed
	FromEngineer string    `json:"from_engineer,omitempty"`
	ToEngineer   string    `json:"to_engineer,omitempty"`
	FromStatus   string    `json:"from_status,omitempty"`
	ToStatus     string    `json:"to_status,omitempty"`
	ActionBy     string    `json:"action_by"`
	ActionAt     time.Time `json:"action_at"`
	Reason       string    `json:"reason"`
	
	// Relations
	Assignment   EngineerAssignment `json:"assignment,omitempty" gorm:"foreignKey:AssignmentID"`
	FromAccount  *Account          `json:"from_account,omitempty" gorm:"foreignKey:FromEngineer"`
	ToAccount    *Account          `json:"to_account,omitempty" gorm:"foreignKey:ToEngineer"`
	ActionByUser Account           `json:"action_by_user,omitempty" gorm:"foreignKey:ActionBy"`
}

// EngineerAvailability 工程師可用性資訊
type EngineerAvailability struct {
	EngineerID     string   `json:"engineer_id"`
	EngineerName   string   `json:"engineer_name"`
	Department     string   `json:"department"`
	Expertise      []string `json:"expertise"`
	CurrentLoad    int      `json:"current_load"`
	MaxLoad        int      `json:"max_load"`
	IsAvailable    bool     `json:"is_available"`
	ExpertiseMatch float64  `json:"expertise_match"` // 0-1 專長匹配度
}

// EngineerWorkloadSummary 工程師工作負載摘要
type EngineerWorkloadSummary struct {
	EngineerID        string  `json:"engineer_id"`
	EngineerName      string  `json:"engineer_name"`
	TotalAssignments  int     `json:"total_assignments"`
	Pending           int     `json:"pending"`
	InProgress        int     `json:"in_progress"`
	Completed         int     `json:"completed"`
	CompletedOnTime   int     `json:"completed_on_time"`
	Overdue           int     `json:"overdue"`
	AvgCompletionTime float64 `json:"avg_completion_time"` // 小時
}

// AssignmentStats 分派統計數據
type AssignmentStats struct {
	Period            string             `json:"period"`
	TotalAssignments  int64              `json:"total_assignments"`
	StatusBreakdown   map[string]int64   `json:"status_breakdown"`
	AvgCompletionTime float64            `json:"avg_completion_time"` // 小時
	OnTimeRate        float64            `json:"on_time_rate"`        // 0-1
	EngineerStats     []EngineerStat     `json:"engineer_stats"`
	TimeSeries        []TimeSeriesData   `json:"time_series"`
}

// EngineerStat 單個工程師統計
type EngineerStat struct {
	EngineerID        string  `json:"engineer_id"`
	EngineerName      string  `json:"engineer_name"`
	TotalAssigned     int64   `json:"total_assigned"`
	Completed         int64   `json:"completed"`
	AvgCompletionTime float64 `json:"avg_completion_time"`
	CompletionRate    float64 `json:"completion_rate"`
}

// TimeSeriesData 時間序列數據
type TimeSeriesData struct {
	Period    string `json:"period"`
	Assigned  int64  `json:"assigned"`
	Completed int64  `json:"completed"`
}