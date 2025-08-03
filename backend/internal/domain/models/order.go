package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusDraft      OrderStatus = "draft"
	OrderStatusConfirmed  OrderStatus = "confirmed"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
)

// Order represents a customer order
type Order struct {
	BaseModel
	OrderNo            string          `gorm:"type:varchar(50);uniqueIndex;not null" json:"order_no"`
	CompanyID          uuid.UUID       `gorm:"type:uuid;not null" json:"company_id"`
	Company            Company         `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	CustomerID         uuid.UUID       `gorm:"type:uuid;not null" json:"customer_id"`
	Customer           Customer        `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	QuoteID            *uuid.UUID      `gorm:"type:uuid" json:"quote_id"`
	Quote              *Quote          `gorm:"foreignKey:QuoteID" json:"quote,omitempty"`
	CustomerPONo       string          `gorm:"type:varchar(50)" json:"customer_po_no"`
	CustomerPODate     time.Time       `json:"customer_po_date"`
	SalesID            uuid.UUID       `gorm:"type:uuid;not null" json:"sales_id"`
	Sales              Account         `gorm:"foreignKey:SalesID" json:"sales,omitempty"`
	Status             OrderStatus     `gorm:"type:varchar(20);not null;default:'draft'" json:"status"`
	OrderDate          time.Time       `json:"order_date"`
	RequiredDate       time.Time       `json:"required_date"`
	ShippedDate        *time.Time      `json:"shipped_date"`
	DeliveredDate      *time.Time      `json:"delivered_date"`
	Currency           string          `gorm:"type:varchar(3);not null" json:"currency"`
	ExchangeRate       decimal.Decimal `gorm:"type:decimal(10,6)" json:"exchange_rate"`
	Incoterm           string          `gorm:"type:varchar(10);not null" json:"incoterm"`
	PaymentTerms       string          `gorm:"type:varchar(100)" json:"payment_terms"`
	DeliveryAddress    string          `gorm:"type:text" json:"delivery_address"`
	BillingAddress     string          `gorm:"type:text" json:"billing_address"`
	SubTotal           decimal.Decimal `gorm:"type:decimal(15,2);not null" json:"sub_total"`
	DiscountPercent    decimal.Decimal `gorm:"type:decimal(5,2)" json:"discount_percent"`
	DiscountAmount     decimal.Decimal `gorm:"type:decimal(15,2)" json:"discount_amount"`
	TaxPercent         decimal.Decimal `gorm:"type:decimal(5,2)" json:"tax_percent"`
	TaxAmount          decimal.Decimal `gorm:"type:decimal(15,2)" json:"tax_amount"`
	ShippingCost       decimal.Decimal `gorm:"type:decimal(15,2)" json:"shipping_cost"`
	TotalAmount        decimal.Decimal `gorm:"type:decimal(15,2);not null" json:"total_amount"`
	Notes              string          `gorm:"type:text" json:"notes"`
	InternalNotes      string          `gorm:"type:text" json:"internal_notes"`
	
	// Relationships
	Items      []OrderItem    `gorm:"foreignKey:OrderID" json:"items,omitempty"`
	Shipments  []Shipment     `gorm:"foreignKey:OrderID" json:"shipments,omitempty"`
	Invoices   []Invoice      `gorm:"foreignKey:OrderID" json:"invoices,omitempty"`
}

func (Order) TableName() string {
	return "orders"
}

// OrderItem represents an item in an order
type OrderItem struct {
	BaseModel
	OrderID          uuid.UUID       `gorm:"type:uuid;not null" json:"order_id"`
	Order            Order           `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	ProductID        *uuid.UUID      `gorm:"type:uuid" json:"product_id"`
	Product          *Product        `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	ItemNo           int             `gorm:"not null" json:"item_no"`
	Description      string          `gorm:"type:text;not null" json:"description"`
	Specifications   string          `gorm:"type:jsonb" json:"specifications"`
	Quantity         decimal.Decimal `gorm:"type:decimal(15,3);not null" json:"quantity"`
	Unit             string          `gorm:"type:varchar(20);not null" json:"unit"`
	UnitPrice        decimal.Decimal `gorm:"type:decimal(15,4);not null" json:"unit_price"`
	DiscountPercent  decimal.Decimal `gorm:"type:decimal(5,2)" json:"discount_percent"`
	Amount           decimal.Decimal `gorm:"type:decimal(15,2);not null" json:"amount"`
	ShippedQuantity  decimal.Decimal `gorm:"type:decimal(15,3)" json:"shipped_quantity"`
	DeliveredQuantity decimal.Decimal `gorm:"type:decimal(15,3)" json:"delivered_quantity"`
	Notes            string          `gorm:"type:text" json:"notes"`
}

func (OrderItem) TableName() string {
	return "order_items"
}