package commands

import (
	"errors"
	"time"

	"github.com/fastenmind/fastener-api/internal/domain/cqrs"
	"github.com/google/uuid"
)

// CreateInquiryCommand creates a new inquiry
type CreateInquiryCommand struct {
	cqrs.BaseCommand
	CompanyID           uuid.UUID `json:"company_id"`
	CustomerID          uuid.UUID `json:"customer_id"`
	ProductCategory     string    `json:"product_category"`
	ProductName         string    `json:"product_name"`
	DrawingFiles        []string  `json:"drawing_files"`
	Quantity            int       `json:"quantity"`
	Unit                string    `json:"unit"`
	RequiredDate        time.Time `json:"required_date"`
	Incoterm            string    `json:"incoterm"`
	DestinationPort     string    `json:"destination_port"`
	DestinationAddress  string    `json:"destination_address"`
	PaymentTerms        string    `json:"payment_terms"`
	SpecialRequirements string    `json:"special_requirements"`
}

// NewCreateInquiryCommand creates a new create inquiry command
func NewCreateInquiryCommand(userID, companyID, customerID uuid.UUID) *CreateInquiryCommand {
	return &CreateInquiryCommand{
		BaseCommand: cqrs.NewBaseCommand("CreateInquiry", userID),
		CompanyID:   companyID,
		CustomerID:  customerID,
	}
}

// Validate validates the command
func (c *CreateInquiryCommand) Validate() error {
	if c.CompanyID == uuid.Nil {
		return errors.New("company ID is required")
	}
	if c.CustomerID == uuid.Nil {
		return errors.New("customer ID is required")
	}
	if c.ProductCategory == "" {
		return errors.New("product category is required")
	}
	if c.ProductName == "" {
		return errors.New("product name is required")
	}
	if c.Quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}
	if c.Unit == "" {
		return errors.New("unit is required")
	}
	if c.RequiredDate.IsZero() {
		return errors.New("required date is required")
	}
	if c.Incoterm == "" {
		return errors.New("incoterm is required")
	}
	return nil
}

// AssignInquiryCommand assigns an inquiry to an engineer
type AssignInquiryCommand struct {
	cqrs.BaseCommand
	InquiryID      uuid.UUID `json:"inquiry_id"`
	EngineerID     uuid.UUID `json:"engineer_id"`
	AssignmentNote string    `json:"assignment_note"`
}

// NewAssignInquiryCommand creates a new assign inquiry command
func NewAssignInquiryCommand(userID, inquiryID, engineerID uuid.UUID) *AssignInquiryCommand {
	return &AssignInquiryCommand{
		BaseCommand: cqrs.NewBaseCommand("AssignInquiry", userID),
		InquiryID:   inquiryID,
		EngineerID:  engineerID,
	}
}

// Validate validates the command
func (c *AssignInquiryCommand) Validate() error {
	if c.InquiryID == uuid.Nil {
		return errors.New("inquiry ID is required")
	}
	if c.EngineerID == uuid.Nil {
		return errors.New("engineer ID is required")
	}
	return nil
}

// RejectInquiryCommand rejects an inquiry
type RejectInquiryCommand struct {
	cqrs.BaseCommand
	InquiryID       uuid.UUID `json:"inquiry_id"`
	RejectionReason string    `json:"rejection_reason"`
}

// NewRejectInquiryCommand creates a new reject inquiry command
func NewRejectInquiryCommand(userID, inquiryID uuid.UUID, reason string) *RejectInquiryCommand {
	return &RejectInquiryCommand{
		BaseCommand:     cqrs.NewBaseCommand("RejectInquiry", userID),
		InquiryID:       inquiryID,
		RejectionReason: reason,
	}
}

// Validate validates the command
func (c *RejectInquiryCommand) Validate() error {
	if c.InquiryID == uuid.Nil {
		return errors.New("inquiry ID is required")
	}
	if c.RejectionReason == "" {
		return errors.New("rejection reason is required")
	}
	return nil
}

// UpdateInquiryCommand updates an inquiry
type UpdateInquiryCommand struct {
	cqrs.BaseCommand
	InquiryID           uuid.UUID  `json:"inquiry_id"`
	ProductName         *string    `json:"product_name,omitempty"`
	Quantity            *int       `json:"quantity,omitempty"`
	RequiredDate        *time.Time `json:"required_date,omitempty"`
	SpecialRequirements *string    `json:"special_requirements,omitempty"`
}

// NewUpdateInquiryCommand creates a new update inquiry command
func NewUpdateInquiryCommand(userID, inquiryID uuid.UUID) *UpdateInquiryCommand {
	return &UpdateInquiryCommand{
		BaseCommand: cqrs.NewBaseCommand("UpdateInquiry", userID),
		InquiryID:   inquiryID,
	}
}

// Validate validates the command
func (c *UpdateInquiryCommand) Validate() error {
	if c.InquiryID == uuid.Nil {
		return errors.New("inquiry ID is required")
	}
	
	// Check if at least one field is being updated
	if c.ProductName == nil && c.Quantity == nil && 
	   c.RequiredDate == nil && c.SpecialRequirements == nil {
		return errors.New("at least one field must be updated")
	}
	
	// Validate quantity if provided
	if c.Quantity != nil && *c.Quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}
	
	return nil
}