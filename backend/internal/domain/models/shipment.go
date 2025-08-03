package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ShipmentStatus represents the status of a shipment
type ShipmentStatus string

const (
	ShipmentStatusPending    ShipmentStatus = "pending"
	ShipmentStatusInTransit  ShipmentStatus = "in_transit"
	ShipmentStatusDelivered  ShipmentStatus = "delivered"
	ShipmentStatusCancelled  ShipmentStatus = "cancelled"
)

// Shipment represents a shipment record
type Shipment struct {
	BaseModel
	ShipmentNo      string          `gorm:"type:varchar(50);uniqueIndex;not null" json:"shipment_no"`
	OrderID         uuid.UUID       `gorm:"type:uuid;not null" json:"order_id"`
	Order           Order           `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	CompanyID       uuid.UUID       `gorm:"type:uuid;not null" json:"company_id"`
	Company         Company         `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Status          ShipmentStatus  `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	ShipmentDate    time.Time       `json:"shipment_date"`
	EstimatedArrival time.Time      `json:"estimated_arrival"`
	ActualArrival   *time.Time      `json:"actual_arrival"`
	Carrier         string          `gorm:"type:varchar(100)" json:"carrier"`
	TrackingNo      string          `gorm:"type:varchar(100)" json:"tracking_no"`
	FreightCost     decimal.Decimal `gorm:"type:decimal(15,2)" json:"freight_cost"`
	Currency        string          `gorm:"type:varchar(3)" json:"currency"`
	FromAddress     string          `gorm:"type:text" json:"from_address"`
	ToAddress       string          `gorm:"type:text" json:"to_address"`
	Notes           string          `gorm:"type:text" json:"notes"`
	
	// Relationships
	Items []ShipmentItem `gorm:"foreignKey:ShipmentID" json:"items,omitempty"`
}

func (Shipment) TableName() string {
	return "shipments"
}

// ShipmentItem represents items in a shipment
type ShipmentItem struct {
	BaseModel
	ShipmentID   uuid.UUID       `gorm:"type:uuid;not null" json:"shipment_id"`
	Shipment     Shipment        `gorm:"foreignKey:ShipmentID" json:"shipment,omitempty"`
	OrderItemID  uuid.UUID       `gorm:"type:uuid;not null" json:"order_item_id"`
	OrderItem    OrderItem       `gorm:"foreignKey:OrderItemID" json:"order_item,omitempty"`
	Quantity     decimal.Decimal `gorm:"type:decimal(15,3);not null" json:"quantity"`
	PackageCount int             `json:"package_count"`
	Weight       decimal.Decimal `gorm:"type:decimal(10,3)" json:"weight"`
	WeightUnit   string          `gorm:"type:varchar(10)" json:"weight_unit"`
	Notes        string          `gorm:"type:text" json:"notes"`
}

func (ShipmentItem) TableName() string {
	return "shipment_items"
}