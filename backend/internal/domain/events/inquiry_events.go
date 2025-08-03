package events

import (
	"time"

	"github.com/google/uuid"
)

// Inquiry event types
const (
	InquiryCreated   EventType = "inquiry.created"
	InquiryAssigned  EventType = "inquiry.assigned"
	InquiryQuoted    EventType = "inquiry.quoted"
	InquiryRejected  EventType = "inquiry.rejected"
	InquiryCancelled EventType = "inquiry.cancelled"
	InquiryUpdated   EventType = "inquiry.updated"
)

// InquiryCreatedEvent is emitted when a new inquiry is created
type InquiryCreatedEvent struct {
	BaseEvent
	InquiryNo           string    `json:"inquiry_no"`
	CustomerID          uuid.UUID `json:"customer_id"`
	SalesID             uuid.UUID `json:"sales_id"`
	ProductCategory     string    `json:"product_category"`
	ProductName         string    `json:"product_name"`
	Quantity            int       `json:"quantity"`
	RequiredDate        time.Time `json:"required_date"`
}

// NewInquiryCreatedEvent creates a new inquiry created event
func NewInquiryCreatedEvent(inquiryID uuid.UUID, inquiryNo string, customerID, salesID uuid.UUID) *InquiryCreatedEvent {
	return &InquiryCreatedEvent{
		BaseEvent:   NewBaseEvent(InquiryCreated, inquiryID, "Inquiry"),
		InquiryNo:   inquiryNo,
		CustomerID:  customerID,
		SalesID:     salesID,
	}
}

// GetData returns the event data
func (e *InquiryCreatedEvent) GetData() interface{} {
	return e
}

// InquiryAssignedEvent is emitted when an inquiry is assigned to an engineer
type InquiryAssignedEvent struct {
	BaseEvent
	InquiryNo    string    `json:"inquiry_no"`
	EngineerID   uuid.UUID `json:"engineer_id"`
	AssignedBy   uuid.UUID `json:"assigned_by"`
	AssignedAt   time.Time `json:"assigned_at"`
	AssignReason string    `json:"assign_reason,omitempty"`
}

// NewInquiryAssignedEvent creates a new inquiry assigned event
func NewInquiryAssignedEvent(inquiryID uuid.UUID, inquiryNo string, engineerID, assignedBy uuid.UUID) *InquiryAssignedEvent {
	return &InquiryAssignedEvent{
		BaseEvent:  NewBaseEvent(InquiryAssigned, inquiryID, "Inquiry"),
		InquiryNo:  inquiryNo,
		EngineerID: engineerID,
		AssignedBy: assignedBy,
		AssignedAt: time.Now().UTC(),
	}
}

// GetData returns the event data
func (e *InquiryAssignedEvent) GetData() interface{} {
	return e
}

// InquiryQuotedEvent is emitted when a quote is created for an inquiry
type InquiryQuotedEvent struct {
	BaseEvent
	InquiryNo    string    `json:"inquiry_no"`
	QuoteID      uuid.UUID `json:"quote_id"`
	QuoteNo      string    `json:"quote_no"`
	TotalAmount  float64   `json:"total_amount"`
	Currency     string    `json:"currency"`
	ValidityDays int       `json:"validity_days"`
	QuotedBy     uuid.UUID `json:"quoted_by"`
	QuotedAt     time.Time `json:"quoted_at"`
}

// NewInquiryQuotedEvent creates a new inquiry quoted event
func NewInquiryQuotedEvent(inquiryID uuid.UUID, inquiryNo string, quoteID uuid.UUID, quoteNo string) *InquiryQuotedEvent {
	return &InquiryQuotedEvent{
		BaseEvent: NewBaseEvent(InquiryQuoted, inquiryID, "Inquiry"),
		InquiryNo: inquiryNo,
		QuoteID:   quoteID,
		QuoteNo:   quoteNo,
		QuotedAt:  time.Now().UTC(),
	}
}

// GetData returns the event data
func (e *InquiryQuotedEvent) GetData() interface{} {
	return e
}

// InquiryRejectedEvent is emitted when an inquiry is rejected
type InquiryRejectedEvent struct {
	BaseEvent
	InquiryNo      string    `json:"inquiry_no"`
	RejectedBy     uuid.UUID `json:"rejected_by"`
	RejectedAt     time.Time `json:"rejected_at"`
	RejectionReason string   `json:"rejection_reason"`
}

// NewInquiryRejectedEvent creates a new inquiry rejected event
func NewInquiryRejectedEvent(inquiryID uuid.UUID, inquiryNo string, rejectedBy uuid.UUID, reason string) *InquiryRejectedEvent {
	return &InquiryRejectedEvent{
		BaseEvent:       NewBaseEvent(InquiryRejected, inquiryID, "Inquiry"),
		InquiryNo:       inquiryNo,
		RejectedBy:      rejectedBy,
		RejectedAt:      time.Now().UTC(),
		RejectionReason: reason,
	}
}

// GetData returns the event data
func (e *InquiryRejectedEvent) GetData() interface{} {
	return e
}