package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Inventory represents inventory items
type Inventory struct {
	ID                 uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID          uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	SKU                string     `gorm:"not null;unique" json:"sku"`         // Stock Keeping Unit
	PartNo             string     `gorm:"not null" json:"part_no"`
	Name               string     `gorm:"not null" json:"name"`
	Description        string     `json:"description"`
	Category           string     `json:"category"`                            // raw_material, semi_finished, finished_goods
	
	// Specifications
	Material           string     `json:"material"`
	Specification      string     `json:"specification"`
	SurfaceTreatment   string     `json:"surface_treatment"`
	HeatTreatment      string     `json:"heat_treatment"`
	Unit               string     `gorm:"default:'PCS'" json:"unit"`          // PCS, KG, M, etc.
	
	// Stock Levels
	CurrentStock       float64    `json:"current_stock"`
	AvailableStock     float64    `json:"available_stock"`                     // Current - Reserved
	ReservedStock      float64    `json:"reserved_stock"`                      // Reserved for orders
	MinStock           float64    `json:"min_stock"`                           // Minimum stock level
	MaxStock           float64    `json:"max_stock"`                           // Maximum stock level
	ReorderPoint       float64    `json:"reorder_point"`                       // Reorder trigger point
	ReorderQuantity    float64    `json:"reorder_quantity"`                    // Standard reorder quantity
	
	// Location
	WarehouseID        *uuid.UUID `gorm:"type:uuid" json:"warehouse_id"`
	Location           string     `json:"location"`                            // Rack/Shelf location
	
	// Cost Information
	LastPurchasePrice  float64    `json:"last_purchase_price"`
	AverageCost        float64    `json:"average_cost"`
	StandardCost       float64    `json:"standard_cost"`
	Currency           string     `gorm:"default:'USD'" json:"currency"`
	
	// Supplier Info
	PrimarySupplierID  *uuid.UUID `gorm:"type:uuid" json:"primary_supplier_id"`
	LeadTimeDays       int        `json:"lead_time_days"`
	
	// Status
	Status             string     `gorm:"default:'active'" json:"status"`       // active, inactive, discontinued
	IsActive           bool       `gorm:"default:true" json:"is_active"`
	
	// Timestamps
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	LastStockCheckAt   *time.Time `json:"last_stock_check_at"`
	
	// Relations
	Company            *Company   `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Warehouse          *Warehouse `gorm:"foreignKey:WarehouseID" json:"warehouse,omitempty"`
	PrimarySupplier    *Supplier  `gorm:"foreignKey:PrimarySupplierID" json:"primary_supplier,omitempty"`
}

func (i *Inventory) BeforeCreate(tx *gorm.DB) error {
	i.ID = uuid.New()
	return nil
}

// Warehouse represents storage locations
type Warehouse struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID   uuid.UUID `gorm:"type:uuid;not null" json:"company_id"`
	Code        string    `gorm:"not null;unique" json:"code"`
	Name        string    `gorm:"not null" json:"name"`
	Type        string    `json:"type"`         // main, branch, consignment
	Address     string    `json:"address"`
	Manager     string    `json:"manager"`
	Phone       string    `json:"phone"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Relations
	Company     *Company  `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
}

func (w *Warehouse) BeforeCreate(tx *gorm.DB) error {
	w.ID = uuid.New()
	return nil
}

// StockMovement represents inventory transactions
type StockMovement struct {
	ID               uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID        uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	InventoryID      uuid.UUID  `gorm:"type:uuid;not null" json:"inventory_id"`
	MovementType     string     `gorm:"not null" json:"movement_type"`      // in, out, adjustment, transfer
	Reason           string     `gorm:"not null" json:"reason"`             // purchase, sales, production, adjustment, return, damage
	Quantity         float64    `gorm:"not null" json:"quantity"`           // Positive for in, negative for out
	UnitCost         float64    `json:"unit_cost"`
	TotalCost        float64    `json:"total_cost"`
	
	// Reference
	ReferenceType    string     `json:"reference_type"`                     // order, purchase_order, production, adjustment
	ReferenceID      *uuid.UUID `gorm:"type:uuid" json:"reference_id"`
	ReferenceNo      string     `json:"reference_no"`
	
	// Location
	FromWarehouseID  *uuid.UUID `gorm:"type:uuid" json:"from_warehouse_id"`
	ToWarehouseID    *uuid.UUID `gorm:"type:uuid" json:"to_warehouse_id"`
	FromLocation     string     `json:"from_location"`
	ToLocation       string     `json:"to_location"`
	
	// Stock Levels (snapshot after movement)
	BeforeQuantity   float64    `json:"before_quantity"`
	AfterQuantity    float64    `json:"after_quantity"`
	
	// Additional Info
	BatchNo          string     `json:"batch_no"`
	SerialNo         string     `json:"serial_no"`
	ExpiryDate       *time.Time `json:"expiry_date"`
	Notes            string     `json:"notes"`
	
	// User Info
	CreatedBy        uuid.UUID  `gorm:"type:uuid;not null" json:"created_by"`
	ApprovedBy       *uuid.UUID `gorm:"type:uuid" json:"approved_by"`
	ApprovedAt       *time.Time `json:"approved_at"`
	
	// Timestamps
	CreatedAt        time.Time  `json:"created_at"`
	
	// Relations
	Company          *Company   `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Inventory        *Inventory `gorm:"foreignKey:InventoryID" json:"inventory,omitempty"`
	FromWarehouse    *Warehouse `gorm:"foreignKey:FromWarehouseID" json:"from_warehouse,omitempty"`
	ToWarehouse      *Warehouse `gorm:"foreignKey:ToWarehouseID" json:"to_warehouse,omitempty"`
	Creator          *User      `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Approver         *User      `gorm:"foreignKey:ApprovedBy" json:"approver,omitempty"`
}

func (sm *StockMovement) BeforeCreate(tx *gorm.DB) error {
	sm.ID = uuid.New()
	return nil
}

// StockAlert represents inventory alerts
type StockAlert struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID    uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	InventoryID  uuid.UUID  `gorm:"type:uuid;not null" json:"inventory_id"`
	AlertType    string     `gorm:"not null" json:"alert_type"`         // low_stock, overstock, expiring, reorder
	Status       string     `gorm:"default:'active'" json:"status"`      // active, acknowledged, resolved
	Priority     string     `json:"priority"`                            // high, medium, low
	
	// Alert Details
	CurrentLevel float64    `json:"current_level"`
	ThresholdLevel float64  `json:"threshold_level"`
	Message      string     `json:"message"`
	
	// Resolution
	AcknowledgedBy *uuid.UUID `gorm:"type:uuid" json:"acknowledged_by"`
	AcknowledgedAt *time.Time `json:"acknowledged_at"`
	ResolvedBy     *uuid.UUID `gorm:"type:uuid" json:"resolved_by"`
	ResolvedAt     *time.Time `json:"resolved_at"`
	Resolution     string     `json:"resolution"`
	
	// Timestamps
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	
	// Relations
	Company        *Company   `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Inventory      *Inventory `gorm:"foreignKey:InventoryID" json:"inventory,omitempty"`
	Acknowledger   *User      `gorm:"foreignKey:AcknowledgedBy" json:"acknowledger,omitempty"`
	Resolver       *User      `gorm:"foreignKey:ResolvedBy" json:"resolver,omitempty"`
}

func (sa *StockAlert) BeforeCreate(tx *gorm.DB) error {
	sa.ID = uuid.New()
	return nil
}

// StockTake represents periodic stock counting
type StockTake struct {
	ID             uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CompanyID      uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	ReferenceNo    string     `gorm:"not null;unique" json:"reference_no"`
	WarehouseID    uuid.UUID  `gorm:"type:uuid;not null" json:"warehouse_id"`
	Status         string     `gorm:"default:'draft'" json:"status"`       // draft, in_progress, completed, cancelled
	Type           string     `json:"type"`                                // full, cycle, spot
	
	// Dates
	ScheduledDate  time.Time  `json:"scheduled_date"`
	StartedAt      *time.Time `json:"started_at"`
	CompletedAt    *time.Time `json:"completed_at"`
	
	// Users
	CreatedBy      uuid.UUID  `gorm:"type:uuid;not null" json:"created_by"`
	AssignedTo     uuid.UUID  `gorm:"type:uuid" json:"assigned_to"`
	ReviewedBy     *uuid.UUID `gorm:"type:uuid" json:"reviewed_by"`
	
	// Summary
	TotalItems     int        `json:"total_items"`
	CountedItems   int        `json:"counted_items"`
	VarianceItems  int        `json:"variance_items"`
	TotalVariance  float64    `json:"total_variance"`
	
	Notes          string     `json:"notes"`
	
	// Timestamps
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	
	// Relations
	Company        *Company   `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Warehouse      *Warehouse `gorm:"foreignKey:WarehouseID" json:"warehouse,omitempty"`
	Creator        *User      `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Assignee       *User      `gorm:"foreignKey:AssignedTo" json:"assignee,omitempty"`
	Reviewer       *User      `gorm:"foreignKey:ReviewedBy" json:"reviewer,omitempty"`
}

func (st *StockTake) BeforeCreate(tx *gorm.DB) error {
	st.ID = uuid.New()
	return nil
}

// StockTakeItem represents individual items in a stock take
type StockTakeItem struct {
	ID              uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	StockTakeID     uuid.UUID  `gorm:"type:uuid;not null" json:"stock_take_id"`
	InventoryID     uuid.UUID  `gorm:"type:uuid;not null" json:"inventory_id"`
	
	// Quantities
	SystemQuantity  float64    `json:"system_quantity"`
	CountedQuantity float64    `json:"counted_quantity"`
	Variance        float64    `json:"variance"`
	
	// Status
	Status          string     `gorm:"default:'pending'" json:"status"`     // pending, counted, verified
	CountedBy       *uuid.UUID `gorm:"type:uuid" json:"counted_by"`
	CountedAt       *time.Time `json:"counted_at"`
	VerifiedBy      *uuid.UUID `gorm:"type:uuid" json:"verified_by"`
	VerifiedAt      *time.Time `json:"verified_at"`
	
	Notes           string     `json:"notes"`
	
	// Timestamps
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	
	// Relations
	StockTake       *StockTake `gorm:"foreignKey:StockTakeID" json:"stock_take,omitempty"`
	Inventory       *Inventory `gorm:"foreignKey:InventoryID" json:"inventory,omitempty"`
	Counter         *User      `gorm:"foreignKey:CountedBy" json:"counter,omitempty"`
	Verifier        *User      `gorm:"foreignKey:VerifiedBy" json:"verifier,omitempty"`
}

func (sti *StockTakeItem) BeforeCreate(tx *gorm.DB) error {
	sti.ID = uuid.New()
	return nil
}