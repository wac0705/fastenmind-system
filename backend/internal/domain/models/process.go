package models

import (
	"time"
	
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Process represents a manufacturing process
type Process struct {
	BaseModel
	CompanyID    uuid.UUID       `gorm:"type:uuid;not null" json:"company_id"`
	Company      Company         `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	ProcessCode  string          `gorm:"type:varchar(50);not null" json:"process_code"`
	Name         string          `gorm:"type:varchar(100);not null" json:"name"`
	NameEn       string          `gorm:"type:varchar(100)" json:"name_en"`
	Category     string          `gorm:"type:varchar(50)" json:"category"`
	Description  string          `gorm:"type:text" json:"description"`
	CostPerHour  decimal.Decimal `gorm:"type:decimal(15,4)" json:"cost_per_hour"`
	SetupCost    decimal.Decimal `gorm:"type:decimal(15,4)" json:"setup_cost"`
	MinBatchSize int             `json:"min_batch_size"`
	MaxBatchSize int             `json:"max_batch_size"`
	IsActive     bool            `gorm:"default:true" json:"is_active"`
	
	// Relationships
	Equipment       []Equipment      `gorm:"many2many:process_equipment;" json:"equipment,omitempty"`
	ProductProcesses []ProductProcess `gorm:"foreignKey:ProcessID" json:"product_processes,omitempty"`
}

func (Process) TableName() string {
	return "processes"
}

// Equipment represents manufacturing equipment
type Equipment struct {
	BaseModel
	CompanyID      uuid.UUID       `gorm:"type:uuid;not null" json:"company_id"`
	Company        Company         `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	EquipmentCode  string          `gorm:"type:varchar(50);not null" json:"equipment_code"`
	Name           string          `gorm:"type:varchar(100);not null" json:"name"`
	Model          string          `gorm:"type:varchar(100)" json:"model"`
	Manufacturer   string          `gorm:"type:varchar(100)" json:"manufacturer"`
	Capacity       string          `gorm:"type:varchar(100)" json:"capacity"`
	CostPerHour    decimal.Decimal `gorm:"type:decimal(15,4)" json:"cost_per_hour"`
	MaintenanceCost decimal.Decimal `gorm:"type:decimal(15,4)" json:"maintenance_cost"`
	Location       string          `gorm:"type:varchar(100)" json:"location"`
	Status         string          `gorm:"type:varchar(20);default:'operational'" json:"status"`
	PurchaseDate   *time.Time      `json:"purchase_date"`
	IsActive       bool            `gorm:"default:true" json:"is_active"`
	
	// Relationships
	Processes []Process `gorm:"many2many:process_equipment;" json:"processes,omitempty"`
}

func (Equipment) TableName() string {
	return "equipment"
}