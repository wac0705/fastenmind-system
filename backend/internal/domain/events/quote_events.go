package events

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Quote event types
const (
	QuoteCreated   EventType = "quote.created"
	QuoteSubmitted EventType = "quote.submitted"
	QuoteApproved  EventType = "quote.approved"
	QuoteRejected  EventType = "quote.rejected"
	QuoteExpired   EventType = "quote.expired"
	QuoteRevised   EventType = "quote.revised"
	QuoteOrdered   EventType = "quote.ordered"
)

// QuoteCreatedEvent is emitted when a new quote is created
type QuoteCreatedEvent struct {
	BaseEvent
	QuoteNo      string          `json:"quote_no"`
	CustomerID   uuid.UUID       `json:"customer_id"`
	InquiryID    *uuid.UUID      `json:"inquiry_id,omitempty"`
	TotalAmount  decimal.Decimal `json:"total_amount"`
	Currency     string          `json:"currency"`
	ValidityDays int             `json:"validity_days"`
	CreatedBy    uuid.UUID       `json:"created_by"`
}

// NewQuoteCreatedEvent creates a new quote created event
func NewQuoteCreatedEvent(quoteID uuid.UUID, quoteNo string, customerID uuid.UUID, totalAmount decimal.Decimal, currency string) *QuoteCreatedEvent {
	return &QuoteCreatedEvent{
		BaseEvent:   NewBaseEvent(QuoteCreated, quoteID, "Quote"),
		QuoteNo:     quoteNo,
		CustomerID:  customerID,
		TotalAmount: totalAmount,
		Currency:    currency,
	}
}

// GetData returns the event data
func (e *QuoteCreatedEvent) GetData() interface{} {
	return e
}

// QuoteSubmittedEvent is emitted when a quote is submitted to customer
type QuoteSubmittedEvent struct {
	BaseEvent
	QuoteNo      string    `json:"quote_no"`
	SubmittedBy  uuid.UUID `json:"submitted_by"`
	SubmittedAt  time.Time `json:"submitted_at"`
	SubmitMethod string    `json:"submit_method"` // email, portal, etc.
}

// NewQuoteSubmittedEvent creates a new quote submitted event
func NewQuoteSubmittedEvent(quoteID uuid.UUID, quoteNo string, submittedBy uuid.UUID) *QuoteSubmittedEvent {
	return &QuoteSubmittedEvent{
		BaseEvent:   NewBaseEvent(QuoteSubmitted, quoteID, "Quote"),
		QuoteNo:     quoteNo,
		SubmittedBy: submittedBy,
		SubmittedAt: time.Now().UTC(),
	}
}

// GetData returns the event data
func (e *QuoteSubmittedEvent) GetData() interface{} {
	return e
}

// QuoteApprovedEvent is emitted when a quote is approved
type QuoteApprovedEvent struct {
	BaseEvent
	QuoteNo        string    `json:"quote_no"`
	ApprovedBy     uuid.UUID `json:"approved_by"`
	ApprovedAt     time.Time `json:"approved_at"`
	ApprovalNotes  string    `json:"approval_notes,omitempty"`
}

// NewQuoteApprovedEvent creates a new quote approved event
func NewQuoteApprovedEvent(quoteID uuid.UUID, quoteNo string, approvedBy uuid.UUID) *QuoteApprovedEvent {
	return &QuoteApprovedEvent{
		BaseEvent:  NewBaseEvent(QuoteApproved, quoteID, "Quote"),
		QuoteNo:    quoteNo,
		ApprovedBy: approvedBy,
		ApprovedAt: time.Now().UTC(),
	}
}

// GetData returns the event data
func (e *QuoteApprovedEvent) GetData() interface{} {
	return e
}

// QuoteRejectedEvent is emitted when a quote is rejected
type QuoteRejectedEvent struct {
	BaseEvent
	QuoteNo         string    `json:"quote_no"`
	RejectedBy      uuid.UUID `json:"rejected_by"`
	RejectedAt      time.Time `json:"rejected_at"`
	RejectionReason string    `json:"rejection_reason"`
}

// NewQuoteRejectedEvent creates a new quote rejected event
func NewQuoteRejectedEvent(quoteID uuid.UUID, quoteNo string, rejectedBy uuid.UUID, reason string) *QuoteRejectedEvent {
	return &QuoteRejectedEvent{
		BaseEvent:       NewBaseEvent(QuoteRejected, quoteID, "Quote"),
		QuoteNo:         quoteNo,
		RejectedBy:      rejectedBy,
		RejectedAt:      time.Now().UTC(),
		RejectionReason: reason,
	}
}

// GetData returns the event data
func (e *QuoteRejectedEvent) GetData() interface{} {
	return e
}

// QuoteRevisedEvent is emitted when a quote is revised
type QuoteRevisedEvent struct {
	BaseEvent
	QuoteNo        string                 `json:"quote_no"`
	RevisionNo     int                    `json:"revision_no"`
	RevisedBy      uuid.UUID              `json:"revised_by"`
	RevisedAt      time.Time              `json:"revised_at"`
	RevisionReason string                 `json:"revision_reason"`
	Changes        map[string]interface{} `json:"changes"` // What changed
}

// NewQuoteRevisedEvent creates a new quote revised event
func NewQuoteRevisedEvent(quoteID uuid.UUID, quoteNo string, revisionNo int, revisedBy uuid.UUID, reason string) *QuoteRevisedEvent {
	return &QuoteRevisedEvent{
		BaseEvent:      NewBaseEvent(QuoteRevised, quoteID, "Quote"),
		QuoteNo:        quoteNo,
		RevisionNo:     revisionNo,
		RevisedBy:      revisedBy,
		RevisedAt:      time.Now().UTC(),
		RevisionReason: reason,
		Changes:        make(map[string]interface{}),
	}
}

// GetData returns the event data
func (e *QuoteRevisedEvent) GetData() interface{} {
	return e
}

// QuoteOrderedEvent is emitted when a quote is converted to an order
type QuoteOrderedEvent struct {
	BaseEvent
	QuoteNo        string    `json:"quote_no"`
	OrderID        uuid.UUID `json:"order_id"`
	OrderNo        string    `json:"order_no"`
	CustomerPONo   string    `json:"customer_po_no"`
	OrderedAt      time.Time `json:"ordered_at"`
}

// NewQuoteOrderedEvent creates a new quote ordered event
func NewQuoteOrderedEvent(quoteID uuid.UUID, quoteNo string, orderID uuid.UUID, orderNo string) *QuoteOrderedEvent {
	return &QuoteOrderedEvent{
		BaseEvent: NewBaseEvent(QuoteOrdered, quoteID, "Quote"),
		QuoteNo:   quoteNo,
		OrderID:   orderID,
		OrderNo:   orderNo,
		OrderedAt: time.Now().UTC(),
	}
}

// GetData returns the event data
func (e *QuoteOrderedEvent) GetData() interface{} {
	return e
}