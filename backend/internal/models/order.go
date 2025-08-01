package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Order represents a purchase order
type Order struct {
	ID             uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	OrderNo        string     `gorm:"not null;unique" json:"order_no"`
	QuoteID        uuid.UUID  `gorm:"type:uuid;not null" json:"quote_id"`
	CompanyID      uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	CustomerID     uuid.UUID  `gorm:"type:uuid;not null" json:"customer_id"`
	SalesID        uuid.UUID  `gorm:"type:uuid;not null" json:"sales_id"`
	Status         string     `gorm:"not null" json:"status"` // pending, confirmed, in_production, quality_check, ready_to_ship, shipped, delivered, completed, cancelled
	
	// Order Details
	PONumber       string     `json:"po_number"`       // Customer's PO number
	Quantity       int        `json:"quantity"`
	UnitPrice      float64    `json:"unit_price"`
	TotalAmount    float64    `json:"total_amount"`
	Currency       string     `gorm:"default:'USD'" json:"currency"`
	
	// Delivery Info
	DeliveryMethod string     `json:"delivery_method"` // EXW, FOB, CIF, etc.
	DeliveryDate   time.Time  `json:"delivery_date"`
	ShippingAddress string    `json:"shipping_address"`
	
	// Payment Info
	PaymentTerms   string     `json:"payment_terms"`
	PaymentStatus  string     `json:"payment_status"` // pending, partial, paid
	DownPayment    float64    `json:"down_payment"`
	PaidAmount     float64    `json:"paid_amount"`
	
	// Workflow dates
	ConfirmedAt    *time.Time `json:"confirmed_at,omitempty"`
	InProductionAt *time.Time `json:"in_production_at,omitempty"`
	QualityCheckAt *time.Time `json:"quality_check_at,omitempty"`
	ReadyToShipAt  *time.Time `json:"ready_to_ship_at,omitempty"`
	ShippedAt      *time.Time `json:"shipped_at,omitempty"`
	DeliveredAt    *time.Time `json:"delivered_at,omitempty"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`
	CancelledAt    *time.Time `json:"cancelled_at,omitempty"`
	
	// Additional Info
	Notes          string     `json:"notes,omitempty"`
	InternalNotes  string     `json:"internal_notes,omitempty"`
	
	// Timestamps
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	
	// Relations
	Quote          *Quote     `gorm:"foreignKey:QuoteID" json:"quote,omitempty"`
	Company        *Company   `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Customer       *Customer  `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Sales          *User      `gorm:"foreignKey:SalesID" json:"sales,omitempty"`
}

func (o *Order) BeforeCreate(tx *gorm.DB) error {
	o.ID = uuid.New()
	return nil
}

// OrderItem represents items in an order
type OrderItem struct {
	ID             uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	OrderID        uuid.UUID  `gorm:"type:uuid;not null" json:"order_id"`
	PartNo         string     `gorm:"not null" json:"part_no"`
	Description    string     `json:"description"`
	Quantity       int        `gorm:"not null" json:"quantity"`
	UnitPrice      float64    `gorm:"not null" json:"unit_price"`
	TotalPrice     float64    `gorm:"not null" json:"total_price"`
	
	// Production Info
	Material       string     `json:"material"`
	SurfaceTreatment string   `json:"surface_treatment"`
	HeatTreatment  string     `json:"heat_treatment"`
	Specifications string     `json:"specifications"`
	
	// Timestamps
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	
	// Relations
	Order          *Order     `gorm:"foreignKey:OrderID" json:"order,omitempty"`
}

func (oi *OrderItem) BeforeCreate(tx *gorm.DB) error {
	oi.ID = uuid.New()
	return nil
}

// OrderActivity represents order activity log
type OrderActivity struct {
	ID             uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	OrderID        uuid.UUID  `gorm:"type:uuid;not null" json:"order_id"`
	UserID         uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	Action         string     `gorm:"not null" json:"action"`
	Description    string     `json:"description"`
	Metadata       string     `json:"metadata,omitempty"` // JSON string for additional data
	CreatedAt      time.Time  `json:"created_at"`
	
	// Relations
	Order          *Order     `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	User           *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (oa *OrderActivity) BeforeCreate(tx *gorm.DB) error {
	oa.ID = uuid.New()
	return nil
}

// OrderDocument represents documents attached to an order
type OrderDocument struct {
	ID             uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	OrderID        uuid.UUID  `gorm:"type:uuid;not null" json:"order_id"`
	DocumentType   string     `gorm:"not null" json:"document_type"` // po, invoice, packing_list, bl, certificate, etc.
	FileName       string     `gorm:"not null" json:"file_name"`
	FilePath       string     `gorm:"not null" json:"file_path"`
	FileSize       int64      `json:"file_size"`
	UploadedBy     uuid.UUID  `gorm:"type:uuid;not null" json:"uploaded_by"`
	CreatedAt      time.Time  `json:"created_at"`
	
	// Relations
	Order          *Order     `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Uploader       *User      `gorm:"foreignKey:UploadedBy" json:"uploader,omitempty"`
}

func (od *OrderDocument) BeforeCreate(tx *gorm.DB) error {
	od.ID = uuid.New()
	return nil
}