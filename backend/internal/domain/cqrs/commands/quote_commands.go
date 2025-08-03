package commands

import (
	"errors"

	"github.com/fastenmind/fastener-api/internal/domain/cqrs"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// CreateQuoteCommand creates a new quote
type CreateQuoteCommand struct {
	cqrs.BaseCommand
	CompanyID    uuid.UUID               `json:"company_id"`
	CustomerID   uuid.UUID               `json:"customer_id"`
	InquiryID    *uuid.UUID              `json:"inquiry_id,omitempty"`
	Currency     string                  `json:"currency"`
	ValidityDays int                     `json:"validity_days"`
	Incoterm     string                  `json:"incoterm"`
	PaymentTerms string                  `json:"payment_terms"`
	Items        []CreateQuoteItemCommand `json:"items"`
	Notes        string                  `json:"notes,omitempty"`
}

// CreateQuoteItemCommand represents an item in the quote
type CreateQuoteItemCommand struct {
	ProductID       *uuid.UUID      `json:"product_id,omitempty"`
	Description     string          `json:"description"`
	Specifications  map[string]interface{} `json:"specifications,omitempty"`
	Quantity        decimal.Decimal `json:"quantity"`
	Unit            string          `json:"unit"`
	UnitPrice       decimal.Decimal `json:"unit_price"`
	DiscountPercent decimal.Decimal `json:"discount_percent,omitempty"`
	LeadTimeDays    int             `json:"lead_time_days"`
}

// NewCreateQuoteCommand creates a new create quote command
func NewCreateQuoteCommand(userID, companyID, customerID uuid.UUID) *CreateQuoteCommand {
	return &CreateQuoteCommand{
		BaseCommand:  cqrs.NewBaseCommand("CreateQuote", userID),
		CompanyID:    companyID,
		CustomerID:   customerID,
		ValidityDays: 30, // Default validity
		Items:        make([]CreateQuoteItemCommand, 0),
	}
}

// Validate validates the command
func (c *CreateQuoteCommand) Validate() error {
	if c.CompanyID == uuid.Nil {
		return errors.New("company ID is required")
	}
	if c.CustomerID == uuid.Nil {
		return errors.New("customer ID is required")
	}
	if c.Currency == "" {
		return errors.New("currency is required")
	}
	if c.Incoterm == "" {
		return errors.New("incoterm is required")
	}
	if len(c.Items) == 0 {
		return errors.New("at least one item is required")
	}
	
	// Validate items
	for i, item := range c.Items {
		if item.Description == "" {
			return errors.New("item description is required")
		}
		if item.Quantity.LessThanOrEqual(decimal.Zero) {
			return errors.New("item quantity must be greater than zero")
		}
		if item.Unit == "" {
			return errors.New("item unit is required")
		}
		if item.UnitPrice.LessThan(decimal.Zero) {
			return errors.New("item unit price cannot be negative")
		}
		if item.DiscountPercent.LessThan(decimal.Zero) || item.DiscountPercent.GreaterThan(decimal.NewFromInt(100)) {
			return errors.New("item discount percent must be between 0 and 100")
		}
		if item.LeadTimeDays < 0 {
			return errors.New("item lead time cannot be negative")
		}
		_ = i // Use index to avoid unused variable warning
	}
	
	return nil
}

// SubmitQuoteCommand submits a quote to customer
type SubmitQuoteCommand struct {
	cqrs.BaseCommand
	QuoteID      uuid.UUID `json:"quote_id"`
	SubmitMethod string    `json:"submit_method"` // email, portal, etc.
	Message      string    `json:"message,omitempty"`
}

// NewSubmitQuoteCommand creates a new submit quote command
func NewSubmitQuoteCommand(userID, quoteID uuid.UUID) *SubmitQuoteCommand {
	return &SubmitQuoteCommand{
		BaseCommand:  cqrs.NewBaseCommand("SubmitQuote", userID),
		QuoteID:      quoteID,
		SubmitMethod: "email", // Default method
	}
}

// Validate validates the command
func (c *SubmitQuoteCommand) Validate() error {
	if c.QuoteID == uuid.Nil {
		return errors.New("quote ID is required")
	}
	if c.SubmitMethod == "" {
		return errors.New("submit method is required")
	}
	return nil
}

// ApproveQuoteCommand approves a quote
type ApproveQuoteCommand struct {
	cqrs.BaseCommand
	QuoteID       uuid.UUID `json:"quote_id"`
	ApprovalNotes string    `json:"approval_notes,omitempty"`
}

// NewApproveQuoteCommand creates a new approve quote command
func NewApproveQuoteCommand(userID, quoteID uuid.UUID) *ApproveQuoteCommand {
	return &ApproveQuoteCommand{
		BaseCommand: cqrs.NewBaseCommand("ApproveQuote", userID),
		QuoteID:     quoteID,
	}
}

// Validate validates the command
func (c *ApproveQuoteCommand) Validate() error {
	if c.QuoteID == uuid.Nil {
		return errors.New("quote ID is required")
	}
	return nil
}

// RejectQuoteCommand rejects a quote
type RejectQuoteCommand struct {
	cqrs.BaseCommand
	QuoteID         uuid.UUID `json:"quote_id"`
	RejectionReason string    `json:"rejection_reason"`
}

// NewRejectQuoteCommand creates a new reject quote command
func NewRejectQuoteCommand(userID, quoteID uuid.UUID, reason string) *RejectQuoteCommand {
	return &RejectQuoteCommand{
		BaseCommand:     cqrs.NewBaseCommand("RejectQuote", userID),
		QuoteID:         quoteID,
		RejectionReason: reason,
	}
}

// Validate validates the command
func (c *RejectQuoteCommand) Validate() error {
	if c.QuoteID == uuid.Nil {
		return errors.New("quote ID is required")
	}
	if c.RejectionReason == "" {
		return errors.New("rejection reason is required")
	}
	return nil
}

// ReviseQuoteCommand creates a new revision of a quote
type ReviseQuoteCommand struct {
	cqrs.BaseCommand
	QuoteID        uuid.UUID                `json:"quote_id"`
	RevisionReason string                   `json:"revision_reason"`
	Items          []CreateQuoteItemCommand `json:"items,omitempty"`
	ValidityDays   *int                     `json:"validity_days,omitempty"`
	PaymentTerms   *string                  `json:"payment_terms,omitempty"`
	Notes          *string                  `json:"notes,omitempty"`
}

// NewReviseQuoteCommand creates a new revise quote command
func NewReviseQuoteCommand(userID, quoteID uuid.UUID, reason string) *ReviseQuoteCommand {
	return &ReviseQuoteCommand{
		BaseCommand:    cqrs.NewBaseCommand("ReviseQuote", userID),
		QuoteID:        quoteID,
		RevisionReason: reason,
	}
}

// Validate validates the command
func (c *ReviseQuoteCommand) Validate() error {
	if c.QuoteID == uuid.Nil {
		return errors.New("quote ID is required")
	}
	if c.RevisionReason == "" {
		return errors.New("revision reason is required")
	}
	
	// Validate items if provided
	for _, item := range c.Items {
		if item.Description == "" {
			return errors.New("item description is required")
		}
		if item.Quantity.LessThanOrEqual(decimal.Zero) {
			return errors.New("item quantity must be greater than zero")
		}
		if item.UnitPrice.LessThan(decimal.Zero) {
			return errors.New("item unit price cannot be negative")
		}
	}
	
	return nil
}