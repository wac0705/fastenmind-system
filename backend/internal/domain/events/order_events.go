package events

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Order event types
const (
	OrderCreated    EventType = "order.created"
	OrderConfirmed  EventType = "order.confirmed"
	OrderProcessing EventType = "order.processing"
	OrderShipped    EventType = "order.shipped"
	OrderDelivered  EventType = "order.delivered"
	OrderCancelled  EventType = "order.cancelled"
	OrderModified   EventType = "order.modified"
	OrderInvoiced   EventType = "order.invoiced"
)

// OrderCreatedEvent is emitted when a new order is created
type OrderCreatedEvent struct {
	BaseEvent
	OrderNo        string          `json:"order_no"`
	CustomerID     uuid.UUID       `json:"customer_id"`
	QuoteID        *uuid.UUID      `json:"quote_id,omitempty"`
	CustomerPONo   string          `json:"customer_po_no"`
	TotalAmount    decimal.Decimal `json:"total_amount"`
	Currency       string          `json:"currency"`
	RequiredDate   time.Time       `json:"required_date"`
	CreatedBy      uuid.UUID       `json:"created_by"`
}

// NewOrderCreatedEvent creates a new order created event
func NewOrderCreatedEvent(orderID uuid.UUID, orderNo string, customerID uuid.UUID, totalAmount decimal.Decimal) *OrderCreatedEvent {
	return &OrderCreatedEvent{
		BaseEvent:   NewBaseEvent(OrderCreated, orderID, "Order"),
		OrderNo:     orderNo,
		CustomerID:  customerID,
		TotalAmount: totalAmount,
	}
}

// GetData returns the event data
func (e *OrderCreatedEvent) GetData() interface{} {
	return e
}

// OrderConfirmedEvent is emitted when an order is confirmed
type OrderConfirmedEvent struct {
	BaseEvent
	OrderNo       string    `json:"order_no"`
	ConfirmedBy   uuid.UUID `json:"confirmed_by"`
	ConfirmedAt   time.Time `json:"confirmed_at"`
	EstimatedDate time.Time `json:"estimated_date"`
}

// NewOrderConfirmedEvent creates a new order confirmed event
func NewOrderConfirmedEvent(orderID uuid.UUID, orderNo string, confirmedBy uuid.UUID) *OrderConfirmedEvent {
	return &OrderConfirmedEvent{
		BaseEvent:   NewBaseEvent(OrderConfirmed, orderID, "Order"),
		OrderNo:     orderNo,
		ConfirmedBy: confirmedBy,
		ConfirmedAt: time.Now().UTC(),
	}
}

// GetData returns the event data
func (e *OrderConfirmedEvent) GetData() interface{} {
	return e
}

// OrderShippedEvent is emitted when an order is shipped
type OrderShippedEvent struct {
	BaseEvent
	OrderNo      string          `json:"order_no"`
	ShipmentID   uuid.UUID       `json:"shipment_id"`
	ShipmentNo   string          `json:"shipment_no"`
	Carrier      string          `json:"carrier"`
	TrackingNo   string          `json:"tracking_no"`
	ShippedItems []ShippedItem   `json:"shipped_items"`
	ShippedAt    time.Time       `json:"shipped_at"`
}

// ShippedItem represents an item that was shipped
type ShippedItem struct {
	OrderItemID uuid.UUID       `json:"order_item_id"`
	ProductID   uuid.UUID       `json:"product_id"`
	Quantity    decimal.Decimal `json:"quantity"`
}

// NewOrderShippedEvent creates a new order shipped event
func NewOrderShippedEvent(orderID uuid.UUID, orderNo string, shipmentID uuid.UUID, shipmentNo string) *OrderShippedEvent {
	return &OrderShippedEvent{
		BaseEvent:  NewBaseEvent(OrderShipped, orderID, "Order"),
		OrderNo:    orderNo,
		ShipmentID: shipmentID,
		ShipmentNo: shipmentNo,
		ShippedAt:  time.Now().UTC(),
	}
}

// GetData returns the event data
func (e *OrderShippedEvent) GetData() interface{} {
	return e
}

// OrderDeliveredEvent is emitted when an order is delivered
type OrderDeliveredEvent struct {
	BaseEvent
	OrderNo       string    `json:"order_no"`
	DeliveredAt   time.Time `json:"delivered_at"`
	ReceivedBy    string    `json:"received_by"`
	DeliveryNotes string    `json:"delivery_notes,omitempty"`
}

// NewOrderDeliveredEvent creates a new order delivered event
func NewOrderDeliveredEvent(orderID uuid.UUID, orderNo string, receivedBy string) *OrderDeliveredEvent {
	return &OrderDeliveredEvent{
		BaseEvent:   NewBaseEvent(OrderDelivered, orderID, "Order"),
		OrderNo:     orderNo,
		DeliveredAt: time.Now().UTC(),
		ReceivedBy:  receivedBy,
	}
}

// GetData returns the event data
func (e *OrderDeliveredEvent) GetData() interface{} {
	return e
}

// OrderCancelledEvent is emitted when an order is cancelled
type OrderCancelledEvent struct {
	BaseEvent
	OrderNo            string    `json:"order_no"`
	CancelledBy        uuid.UUID `json:"cancelled_by"`
	CancelledAt        time.Time `json:"cancelled_at"`
	CancellationReason string    `json:"cancellation_reason"`
	RefundAmount       decimal.Decimal `json:"refund_amount,omitempty"`
}

// NewOrderCancelledEvent creates a new order cancelled event
func NewOrderCancelledEvent(orderID uuid.UUID, orderNo string, cancelledBy uuid.UUID, reason string) *OrderCancelledEvent {
	return &OrderCancelledEvent{
		BaseEvent:          NewBaseEvent(OrderCancelled, orderID, "Order"),
		OrderNo:            orderNo,
		CancelledBy:        cancelledBy,
		CancelledAt:        time.Now().UTC(),
		CancellationReason: reason,
	}
}

// GetData returns the event data
func (e *OrderCancelledEvent) GetData() interface{} {
	return e
}

// OrderInvoicedEvent is emitted when an invoice is created for an order
type OrderInvoicedEvent struct {
	BaseEvent
	OrderNo      string          `json:"order_no"`
	InvoiceID    uuid.UUID       `json:"invoice_id"`
	InvoiceNo    string          `json:"invoice_no"`
	InvoiceAmount decimal.Decimal `json:"invoice_amount"`
	Currency     string          `json:"currency"`
	DueDate      time.Time       `json:"due_date"`
	InvoicedAt   time.Time       `json:"invoiced_at"`
}

// NewOrderInvoicedEvent creates a new order invoiced event
func NewOrderInvoicedEvent(orderID uuid.UUID, orderNo string, invoiceID uuid.UUID, invoiceNo string) *OrderInvoicedEvent {
	return &OrderInvoicedEvent{
		BaseEvent:  NewBaseEvent(OrderInvoiced, orderID, "Order"),
		OrderNo:    orderNo,
		InvoiceID:  invoiceID,
		InvoiceNo:  invoiceNo,
		InvoicedAt: time.Now().UTC(),
	}
}

// GetData returns the event data
func (e *OrderInvoicedEvent) GetData() interface{} {
	return e
}