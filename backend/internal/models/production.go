package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProductionOrder represents a production order
type ProductionOrder struct {
	ID                uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID         uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	OrderNo           string     `gorm:"not null;unique" json:"order_no"`
	Status            string     `gorm:"not null" json:"status"` // planned, released, in_progress, quality_check, completed, cancelled
	Priority          string     `gorm:"not null" json:"priority"` // low, medium, high, urgent
	
	// Reference
	SalesOrderID      *uuid.UUID `gorm:"type:uuid" json:"sales_order_id"`
	CustomerID        *uuid.UUID `gorm:"type:uuid" json:"customer_id"`
	InventoryID       uuid.UUID  `gorm:"type:uuid;not null" json:"inventory_id"`
	
	// Production Details
	ProductName       string     `gorm:"not null" json:"product_name"`
	ProductSpec       string     `json:"product_spec"`
	PlannedQuantity   float64    `gorm:"not null" json:"planned_quantity"`
	ProducedQuantity  float64    `json:"produced_quantity"`
	QualifiedQuantity float64    `json:"qualified_quantity"`
	DefectQuantity    float64    `json:"defect_quantity"`
	Unit              string     `gorm:"not null" json:"unit"`
	
	// Scheduling
	PlannedStartDate  time.Time  `json:"planned_start_date"`
	PlannedEndDate    time.Time  `json:"planned_end_date"`
	ActualStartDate   *time.Time `json:"actual_start_date"`
	ActualEndDate     *time.Time `json:"actual_end_date"`
	
	// Process
	RouteID           *uuid.UUID `gorm:"type:uuid" json:"route_id"`
	CurrentStationID  *uuid.UUID `gorm:"type:uuid" json:"current_station_id"`
	CompletedStations int        `json:"completed_stations"`
	TotalStations     int        `json:"total_stations"`
	
	// Cost
	EstimatedCost     float64    `json:"estimated_cost"`
	ActualCost        float64    `json:"actual_cost"`
	MaterialCost      float64    `json:"material_cost"`
	LaborCost         float64    `json:"labor_cost"`
	OverheadCost      float64    `json:"overhead_cost"`
	Currency          string     `gorm:"default:'USD'" json:"currency"`
	
	// Notes
	Notes             string     `json:"notes"`
	
	// Timestamps
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	CreatedBy         uuid.UUID  `gorm:"type:uuid;not null" json:"created_by"`
	
	// Relations
	Company           *Company           `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	SalesOrder        *Order             `gorm:"foreignKey:SalesOrderID" json:"sales_order,omitempty"`
	Customer          *Customer          `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Inventory         *Inventory         `gorm:"foreignKey:InventoryID" json:"inventory,omitempty"`
	Route             *ProductionRoute   `gorm:"foreignKey:RouteID" json:"route,omitempty"`
	CurrentStation    *WorkStation       `gorm:"foreignKey:CurrentStationID" json:"current_station,omitempty"`
	Creator           *User              `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

func (po *ProductionOrder) BeforeCreate(tx *gorm.DB) error {
	po.ID = uuid.New()
	return nil
}

// ProductionRoute represents a production routing
type ProductionRoute struct {
	ID                uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID         uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	RouteNo           string     `gorm:"not null;unique" json:"route_no"`
	Name              string     `gorm:"not null" json:"name"`
	Description       string     `json:"description"`
	Status            string     `gorm:"default:'active'" json:"status"` // active, inactive
	
	// Product Info
	InventoryID       *uuid.UUID `gorm:"type:uuid" json:"inventory_id"`
	ProductCategory   string     `json:"product_category"`
	
	// Routing Details
	TotalStations     int        `json:"total_stations"`
	EstimatedDuration float64    `json:"estimated_duration"` // in hours
	EstimatedCost     float64    `json:"estimated_cost"`
	
	// Version Control
	Version           int        `gorm:"default:1" json:"version"`
	IsActive          bool       `gorm:"default:true" json:"is_active"`
	
	// Timestamps
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	CreatedBy         uuid.UUID  `gorm:"type:uuid;not null" json:"created_by"`
	
	// Relations
	Company           *Company           `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Inventory         *Inventory         `gorm:"foreignKey:InventoryID" json:"inventory,omitempty"`
	Creator           *User              `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Operations        []RouteOperation   `gorm:"foreignKey:RouteID" json:"operations,omitempty"`
}

func (pr *ProductionRoute) BeforeCreate(tx *gorm.DB) error {
	pr.ID = uuid.New()
	return nil
}

// RouteOperation represents an operation in a production route
type RouteOperation struct {
	ID                uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	RouteID           uuid.UUID  `gorm:"type:uuid;not null" json:"route_id"`
	OperationNo       int        `gorm:"not null" json:"operation_no"`
	WorkStationID     uuid.UUID  `gorm:"type:uuid;not null" json:"work_station_id"`
	
	// Operation Details
	Name              string     `gorm:"not null" json:"name"`
	Description       string     `json:"description"`
	Instructions      string     `json:"instructions"`
	
	// Time
	SetupTime         float64    `json:"setup_time"`     // in minutes
	ProcessTime       float64    `json:"process_time"`   // in minutes per unit
	TeardownTime      float64    `json:"teardown_time"`  // in minutes
	
	// Quality Control
	QCRequired        bool       `gorm:"default:false" json:"qc_required"`
	QCInstructions    string     `json:"qc_instructions"`
	
	// Next Operation
	NextOperationID   *uuid.UUID `gorm:"type:uuid" json:"next_operation_id"`
	
	// Timestamps
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	
	// Relations
	Route             *ProductionRoute `gorm:"foreignKey:RouteID" json:"route,omitempty"`
	WorkStation       *WorkStation     `gorm:"foreignKey:WorkStationID" json:"work_station,omitempty"`
	NextOperation     *RouteOperation  `gorm:"foreignKey:NextOperationID" json:"next_operation,omitempty"`
}

func (ro *RouteOperation) BeforeCreate(tx *gorm.DB) error {
	ro.ID = uuid.New()
	return nil
}

// WorkStation represents a work station or machine
type WorkStation struct {
	ID                uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID         uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	StationNo         string     `gorm:"not null;unique" json:"station_no"`
	Name              string     `gorm:"not null" json:"name"`
	Type              string     `gorm:"not null" json:"type"` // machine, manual, inspection, assembly
	Status            string     `gorm:"default:'available'" json:"status"` // available, busy, maintenance, breakdown
	
	// Capacity
	Capacity          float64    `json:"capacity"`          // units per hour
	UtilizationRate   float64    `json:"utilization_rate"` // percentage
	
	// Location
	Location          string     `json:"location"`
	Department        string     `json:"department"`
	
	// Cost
	HourlyCost        float64    `json:"hourly_cost"`
	
	// Equipment Info
	Model             string     `json:"model"`
	Manufacturer      string     `json:"manufacturer"`
	SerialNumber      string     `json:"serial_number"`
	PurchaseDate      *time.Time `json:"purchase_date"`
	
	// Maintenance
	LastMaintenance   *time.Time `json:"last_maintenance"`
	NextMaintenance   *time.Time `json:"next_maintenance"`
	MaintenanceNotes  string     `json:"maintenance_notes"`
	
	// Notes
	Description       string     `json:"description"`
	Notes             string     `json:"notes"`
	
	// Timestamps
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	CreatedBy         uuid.UUID  `gorm:"type:uuid;not null" json:"created_by"`
	
	// Relations
	Company           *Company   `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Creator           *User      `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

func (ws *WorkStation) BeforeCreate(tx *gorm.DB) error {
	ws.ID = uuid.New()
	return nil
}

// ProductionTask represents a specific task in production
type ProductionTask struct {
	ID                uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	ProductionOrderID uuid.UUID  `gorm:"type:uuid;not null" json:"production_order_id"`
	RouteOperationID  uuid.UUID  `gorm:"type:uuid;not null" json:"route_operation_id"`
	WorkStationID     uuid.UUID  `gorm:"type:uuid;not null" json:"work_station_id"`
	
	// Task Details
	TaskNo            int        `gorm:"not null" json:"task_no"`
	Name              string     `gorm:"not null" json:"name"`
	Status            string     `gorm:"not null" json:"status"` // pending, in_progress, completed, on_hold, cancelled
	
	// Assignment
	AssignedTo        *uuid.UUID `gorm:"type:uuid" json:"assigned_to"`
	AssignedAt        *time.Time `json:"assigned_at"`
	
	// Quantity
	PlannedQuantity   float64    `gorm:"not null" json:"planned_quantity"`
	CompletedQuantity float64    `json:"completed_quantity"`
	QualifiedQuantity float64    `json:"qualified_quantity"`
	DefectQuantity    float64    `json:"defect_quantity"`
	
	// Time Tracking
	PlannedStartTime  time.Time  `json:"planned_start_time"`
	PlannedEndTime    time.Time  `json:"planned_end_time"`
	ActualStartTime   *time.Time `json:"actual_start_time"`
	ActualEndTime     *time.Time `json:"actual_end_time"`
	
	// Quality Control
	QCStatus          string     `json:"qc_status"` // not_required, pending, passed, failed
	QCNotes           string     `json:"qc_notes"`
	QCBy              *uuid.UUID `gorm:"type:uuid" json:"qc_by"`
	QCAt              *time.Time `json:"qc_at"`
	
	// Notes
	Notes             string     `json:"notes"`
	Issues            string     `json:"issues"`
	
	// Timestamps
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	
	// Relations
	ProductionOrder   *ProductionOrder `gorm:"foreignKey:ProductionOrderID" json:"production_order,omitempty"`
	RouteOperation    *RouteOperation  `gorm:"foreignKey:RouteOperationID" json:"route_operation,omitempty"`
	WorkStation       *WorkStation     `gorm:"foreignKey:WorkStationID" json:"work_station,omitempty"`
	AssignedUser      *User            `gorm:"foreignKey:AssignedTo" json:"assigned_user,omitempty"`
	QCUser            *User            `gorm:"foreignKey:QCBy" json:"qc_user,omitempty"`
}

func (pt *ProductionTask) BeforeCreate(tx *gorm.DB) error {
	pt.ID = uuid.New()
	return nil
}

// ProductionMaterial represents materials used in production
type ProductionMaterial struct {
	ID                uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	ProductionOrderID uuid.UUID  `gorm:"type:uuid;not null" json:"production_order_id"`
	InventoryID       uuid.UUID  `gorm:"type:uuid;not null" json:"inventory_id"`
	
	// Material Details
	PlannedQuantity   float64    `gorm:"not null" json:"planned_quantity"`
	IssuedQuantity    float64    `json:"issued_quantity"`
	ConsumedQuantity  float64    `json:"consumed_quantity"`
	ReturnedQuantity  float64    `json:"returned_quantity"`
	Unit              string     `gorm:"not null" json:"unit"`
	
	// Cost
	UnitCost          float64    `json:"unit_cost"`
	TotalCost         float64    `json:"total_cost"`
	
	// Status
	Status            string     `gorm:"default:'planned'" json:"status"` // planned, issued, consumed, returned
	
	// Timestamps
	IssuedAt          *time.Time `json:"issued_at"`
	ConsumedAt        *time.Time `json:"consumed_at"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	
	// Relations
	ProductionOrder   *ProductionOrder `gorm:"foreignKey:ProductionOrderID" json:"production_order,omitempty"`
	Inventory         *Inventory       `gorm:"foreignKey:InventoryID" json:"inventory,omitempty"`
}

func (pm *ProductionMaterial) BeforeCreate(tx *gorm.DB) error {
	pm.ID = uuid.New()
	return nil
}

// QualityInspection represents quality inspection records
type QualityInspection struct {
	ID                uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID         uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	InspectionNo      string     `gorm:"not null;unique" json:"inspection_no"`
	Type              string     `gorm:"not null" json:"type"` // incoming, in_process, final, customer_return
	Status            string     `gorm:"not null" json:"status"` // pending, in_progress, passed, failed, on_hold
	
	// Reference
	ProductionOrderID *uuid.UUID `gorm:"type:uuid" json:"production_order_id"`
	ProductionTaskID  *uuid.UUID `gorm:"type:uuid" json:"production_task_id"`
	InventoryID       *uuid.UUID `gorm:"type:uuid" json:"inventory_id"`
	
	// Inspection Details
	InspectedQuantity float64    `gorm:"not null" json:"inspected_quantity"`
	QualifiedQuantity float64    `json:"qualified_quantity"`
	DefectQuantity    float64    `json:"defect_quantity"`
	Unit              string     `gorm:"not null" json:"unit"`
	
	// Defect Analysis
	DefectTypes       string     `json:"defect_types"`     // JSON array of defect types
	DefectReasons     string     `json:"defect_reasons"`   // JSON array of defect reasons
	CriticalDefects   int        `json:"critical_defects"`
	MajorDefects      int        `json:"major_defects"`
	MinorDefects      int        `json:"minor_defects"`
	
	// Inspector
	InspectorID       uuid.UUID  `gorm:"type:uuid;not null" json:"inspector_id"`
	InspectedAt       time.Time  `json:"inspected_at"`
	
	// Approval
	ApprovedBy        *uuid.UUID `gorm:"type:uuid" json:"approved_by"`
	ApprovedAt        *time.Time `json:"approved_at"`
	
	// Notes
	InspectionNotes   string     `json:"inspection_notes"`
	CorrectiveAction  string     `json:"corrective_action"`
	
	// Attachments
	AttachmentPath    string     `json:"attachment_path"`
	
	// Timestamps
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	
	// Relations
	Company           *Company         `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	ProductionOrder   *ProductionOrder `gorm:"foreignKey:ProductionOrderID" json:"production_order,omitempty"`
	ProductionTask    *ProductionTask  `gorm:"foreignKey:ProductionTaskID" json:"production_task,omitempty"`
	Inventory         *Inventory       `gorm:"foreignKey:InventoryID" json:"inventory,omitempty"`
	Inspector         *User            `gorm:"foreignKey:InspectorID" json:"inspector,omitempty"`
	Approver          *User            `gorm:"foreignKey:ApprovedBy" json:"approver,omitempty"`
}

func (qi *QualityInspection) BeforeCreate(tx *gorm.DB) error {
	qi.ID = uuid.New()
	return nil
}