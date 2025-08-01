package models

import (
	"time"

	"github.com/google/uuid"
)

// ProcessCategory 製程類別
type ProcessCategory struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Code        string    `json:"code" gorm:"type:varchar(20);unique;not null"`
	Name        string    `json:"name" gorm:"type:varchar(100);not null"`
	NameEN      string    `json:"name_en" gorm:"type:varchar(100)"`
	Description string    `json:"description" gorm:"type:text"`
	SortOrder   int       `json:"sort_order" gorm:"default:0"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// Equipment 設備主檔
type Equipment struct {
	ID                    uuid.UUID        `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Code                  string           `json:"code" gorm:"type:varchar(50);unique;not null"`
	Name                  string           `json:"name" gorm:"type:varchar(100);not null"`
	NameEN                string           `json:"name_en" gorm:"type:varchar(100)"`
	ProcessCategoryID     uuid.UUID        `json:"process_category_id" gorm:"type:uuid"`
	ProcessCategory       *ProcessCategory `json:"process_category" gorm:"foreignKey:ProcessCategoryID"`
	Specs                 string           `json:"specs" gorm:"type:text"`
	CapacityPerHour       float64          `json:"capacity_per_hour" gorm:"type:decimal(10,2)"`
	PowerConsumption      float64          `json:"power_consumption" gorm:"type:decimal(10,2)"`
	DepreciationYears     int              `json:"depreciation_years" gorm:"default:10"`
	PurchaseCost          float64          `json:"purchase_cost" gorm:"type:decimal(15,2)"`
	MaintenanceCostPerYear float64         `json:"maintenance_cost_per_year" gorm:"type:decimal(15,2)"`
	Location              string           `json:"location" gorm:"type:varchar(100)"`
	IsActive              bool             `json:"is_active" gorm:"default:true"`
	CreatedAt             time.Time        `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt             time.Time        `json:"updated_at" gorm:"autoUpdateTime"`
}

// ProcessStep 製程步驟
type ProcessStep struct {
	ID                 uuid.UUID        `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Code               string           `json:"code" gorm:"type:varchar(50);unique;not null"`
	Name               string           `json:"name" gorm:"type:varchar(100);not null"`
	NameEN             string           `json:"name_en" gorm:"type:varchar(100)"`
	ProcessCategoryID  uuid.UUID        `json:"process_category_id" gorm:"type:uuid"`
	ProcessCategory    *ProcessCategory `json:"process_category" gorm:"foreignKey:ProcessCategoryID"`
	DefaultEquipmentID *uuid.UUID       `json:"default_equipment_id" gorm:"type:uuid"`
	DefaultEquipment   *Equipment       `json:"default_equipment" gorm:"foreignKey:DefaultEquipmentID"`
	SetupTimeMinutes   float64          `json:"setup_time_minutes" gorm:"type:decimal(10,2);default:0"`
	CycleTimeSeconds   float64          `json:"cycle_time_seconds" gorm:"type:decimal(10,2)"`
	LaborRequired      int              `json:"labor_required" gorm:"default:1"`
	Description        string           `json:"description" gorm:"type:text"`
	SortOrder          int              `json:"sort_order" gorm:"default:0"`
	IsActive           bool             `json:"is_active" gorm:"default:true"`
	CreatedAt          time.Time        `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt          time.Time        `json:"updated_at" gorm:"autoUpdateTime"`
}

// CostParameter 成本參數
type CostParameter struct {
	ID            uuid.UUID  `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ParameterType string     `json:"parameter_type" gorm:"type:varchar(50);not null"`
	ParameterName string     `json:"parameter_name" gorm:"type:varchar(100);not null"`
	Value         float64    `json:"value" gorm:"type:decimal(15,4);not null"`
	Unit          string     `json:"unit" gorm:"type:varchar(20)"`
	EffectiveDate time.Time  `json:"effective_date" gorm:"type:date;not null"`
	EndDate       *time.Time `json:"end_date" gorm:"type:date"`
	Description   string     `json:"description" gorm:"type:text"`
	CreatedAt     time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

// ProductProcessRoute 產品製程路線
type ProductProcessRoute struct {
	ID              uuid.UUID              `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ProductCategory string                 `json:"product_category" gorm:"type:varchar(50);not null"`
	MaterialType    string                 `json:"material_type" gorm:"type:varchar(50)"`
	SizeRange       string                 `json:"size_range" gorm:"type:varchar(50)"`
	RouteName       string                 `json:"route_name" gorm:"type:varchar(100);not null"`
	IsDefault       bool                   `json:"is_default" gorm:"default:false"`
	IsActive        bool                   `json:"is_active" gorm:"default:true"`
	RouteDetails    []ProcessRouteDetail   `json:"route_details" gorm:"foreignKey:RouteID"`
	CreatedAt       time.Time              `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time              `json:"updated_at" gorm:"autoUpdateTime"`
}

// ProcessRouteDetail 製程路線明細
type ProcessRouteDetail struct {
	ID                uuid.UUID    `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	RouteID           uuid.UUID    `json:"route_id" gorm:"type:uuid;not null"`
	Sequence          int          `json:"sequence" gorm:"not null"`
	ProcessStepID     uuid.UUID    `json:"process_step_id" gorm:"type:uuid;not null"`
	ProcessStep       *ProcessStep `json:"process_step" gorm:"foreignKey:ProcessStepID"`
	EquipmentID       *uuid.UUID   `json:"equipment_id" gorm:"type:uuid"`
	Equipment         *Equipment   `json:"equipment" gorm:"foreignKey:EquipmentID"`
	SetupTimeOverride *float64     `json:"setup_time_override" gorm:"type:decimal(10,2)"`
	CycleTimeOverride *float64     `json:"cycle_time_override" gorm:"type:decimal(10,2)"`
	YieldRate         float64      `json:"yield_rate" gorm:"type:decimal(5,2);default:98.00"`
	Notes             string       `json:"notes" gorm:"type:text"`
	CreatedAt         time.Time    `json:"created_at" gorm:"autoCreateTime"`
}

// CostCalculation 成本計算記錄
type CostCalculation struct {
	ID                uuid.UUID                `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	InquiryID         *uuid.UUID               `json:"inquiry_id" gorm:"type:uuid"`
	Inquiry           *Inquiry                 `json:"inquiry" gorm:"foreignKey:InquiryID"`
	QuoteID           *uuid.UUID               `json:"quote_id" gorm:"type:uuid"`
	CalculationNo     string                   `json:"calculation_no" gorm:"type:varchar(50);unique;not null"`
	ProductName       string                   `json:"product_name" gorm:"type:varchar(200);not null"`
	Quantity          int                      `json:"quantity" gorm:"not null"`
	MaterialCost      float64                  `json:"material_cost" gorm:"type:decimal(15,4)"`
	ProcessCost       float64                  `json:"process_cost" gorm:"type:decimal(15,4)"`
	OverheadCost      float64                  `json:"overhead_cost" gorm:"type:decimal(15,4)"`
	TotalCost         float64                  `json:"total_cost" gorm:"type:decimal(15,4)"`
	UnitCost          float64                  `json:"unit_cost" gorm:"type:decimal(15,6)"`
	MarginPercentage  float64                  `json:"margin_percentage" gorm:"type:decimal(5,2)"`
	SellingPrice      float64                  `json:"selling_price" gorm:"type:decimal(15,4)"`
	RouteID           *uuid.UUID               `json:"route_id" gorm:"type:uuid"`
	Route             *ProductProcessRoute     `json:"route" gorm:"foreignKey:RouteID"`
	CalculatedBy      uuid.UUID                `json:"calculated_by" gorm:"type:uuid"`
	CalculatedByUser  Account                  `json:"calculated_by_user" gorm:"foreignKey:CalculatedBy"`
	CalculatedAt      time.Time                `json:"calculated_at" gorm:"autoCreateTime"`
	ApprovedBy        *uuid.UUID               `json:"approved_by" gorm:"type:uuid"`
	ApprovedByUser    *Account                 `json:"approved_by_user" gorm:"foreignKey:ApprovedBy"`
	ApprovedAt        *time.Time               `json:"approved_at"`
	Status            string                   `json:"status" gorm:"type:varchar(20);default:'draft'"`
	Details           []CostCalculationDetail  `json:"details" gorm:"foreignKey:CalculationID"`
	CreatedAt         time.Time                `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time                `json:"updated_at" gorm:"autoUpdateTime"`
}

// CostCalculationDetail 成本計算明細
type CostCalculationDetail struct {
	ID               uuid.UUID    `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CalculationID    uuid.UUID    `json:"calculation_id" gorm:"type:uuid;not null"`
	Sequence         int          `json:"sequence" gorm:"not null"`
	ProcessStepID    uuid.UUID    `json:"process_step_id" gorm:"type:uuid;not null"`
	ProcessStep      *ProcessStep `json:"process_step" gorm:"foreignKey:ProcessStepID"`
	EquipmentID      *uuid.UUID   `json:"equipment_id" gorm:"type:uuid"`
	Equipment        *Equipment   `json:"equipment" gorm:"foreignKey:EquipmentID"`
	SetupTime        float64      `json:"setup_time" gorm:"type:decimal(10,2)"`
	CycleTime        float64      `json:"cycle_time" gorm:"type:decimal(10,2)"`
	TotalTimeHours   float64      `json:"total_time_hours" gorm:"type:decimal(10,4)"`
	LaborCost        float64      `json:"labor_cost" gorm:"type:decimal(15,4)"`
	EquipmentCost    float64      `json:"equipment_cost" gorm:"type:decimal(15,4)"`
	ElectricityCost  float64      `json:"electricity_cost" gorm:"type:decimal(15,4)"`
	OtherCost        float64      `json:"other_cost" gorm:"type:decimal(15,4)"`
	SubtotalCost     float64      `json:"subtotal_cost" gorm:"type:decimal(15,4)"`
	YieldLossCost    float64      `json:"yield_loss_cost" gorm:"type:decimal(15,4)"`
	Notes            string       `json:"notes" gorm:"type:text"`
	CreatedAt        time.Time    `json:"created_at" gorm:"autoCreateTime"`
}

// TableName overrides
func (ProcessCategory) TableName() string       { return "process_categories" }
func (Equipment) TableName() string              { return "equipment" }
func (ProcessStep) TableName() string            { return "process_steps" }
func (CostParameter) TableName() string          { return "cost_parameters" }
func (ProductProcessRoute) TableName() string    { return "product_process_routes" }
func (ProcessRouteDetail) TableName() string     { return "process_route_details" }
func (CostCalculation) TableName() string        { return "cost_calculations" }
func (CostCalculationDetail) TableName() string  { return "cost_calculation_details" }

// Request/Response structures

// CostCalculationRequest 成本計算請求
type CostCalculationRequest struct {
	InquiryID       *uuid.UUID `json:"inquiry_id"`
	ProductName     string     `json:"product_name" binding:"required"`
	ProductCategory string     `json:"product_category" binding:"required"`
	MaterialType    string     `json:"material_type"`
	SizeRange       string     `json:"size_range"`
	Quantity        int        `json:"quantity" binding:"required,min=1"`
	MaterialCost    float64    `json:"material_cost"`
	RouteID         *uuid.UUID `json:"route_id"`
	CustomRoute     []struct {
		ProcessStepID uuid.UUID `json:"process_step_id"`
		EquipmentID   uuid.UUID `json:"equipment_id"`
		SetupTime     float64   `json:"setup_time"`
		CycleTime     float64   `json:"cycle_time"`
	} `json:"custom_route"`
	MarginPercentage float64 `json:"margin_percentage"`
}

// CostSummary 成本摘要
type CostSummary struct {
	MaterialCost      float64 `json:"material_cost"`
	ProcessCost       float64 `json:"process_cost"`
	OverheadCost      float64 `json:"overhead_cost"`
	TotalCost         float64 `json:"total_cost"`
	UnitCost          float64 `json:"unit_cost"`
	SuggestedPrice    float64 `json:"suggested_price"`
	MarginPercentage  float64 `json:"margin_percentage"`
	ProcessBreakdown  []ProcessCostBreakdown `json:"process_breakdown"`
}

// ProcessCostBreakdown 製程成本明細
type ProcessCostBreakdown struct {
	ProcessName     string  `json:"process_name"`
	EquipmentName   string  `json:"equipment_name"`
	TotalTimeHours  float64 `json:"total_time_hours"`
	LaborCost       float64 `json:"labor_cost"`
	EquipmentCost   float64 `json:"equipment_cost"`
	ElectricityCost float64 `json:"electricity_cost"`
	TotalCost       float64 `json:"total_cost"`
}