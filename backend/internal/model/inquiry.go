package model

import (
	"time"

	"github.com/google/uuid"
)

// Inquiry represents an inquiry request
type Inquiry struct {
	Base
	InquiryNo          string     `json:"inquiry_no" db:"inquiry_no"`
	CompanyID          uuid.UUID  `json:"company_id" db:"company_id"`
	CustomerID         uuid.UUID  `json:"customer_id" db:"customer_id"`
	SalesID            uuid.UUID  `json:"sales_id" db:"sales_id"`
	Status             string     `json:"status" db:"status"`
	ProductCategory    string     `json:"product_category" db:"product_category"`
	ProductName        string     `json:"product_name" db:"product_name"`
	DrawingFiles       []string   `json:"drawing_files" db:"drawing_files"`
	Quantity           int        `json:"quantity" db:"quantity"`
	Unit               string     `json:"unit" db:"unit"`
	RequiredDate       time.Time  `json:"required_date" db:"required_date"`
	Incoterm           string     `json:"incoterm" db:"incoterm"`
	DestinationPort    *string    `json:"destination_port,omitempty" db:"destination_port"`
	DestinationAddress *string    `json:"destination_address,omitempty" db:"destination_address"`
	PaymentTerms       *string    `json:"payment_terms,omitempty" db:"payment_terms"`
	SpecialRequirements *string   `json:"special_requirements,omitempty" db:"special_requirements"`
	
	// Assignment
	AssignedEngineerID *uuid.UUID `json:"assigned_engineer_id,omitempty" db:"assigned_engineer_id"`
	AssignedAt         *time.Time `json:"assigned_at,omitempty" db:"assigned_at"`
	
	// Quote
	QuoteID            *uuid.UUID `json:"quote_id,omitempty" db:"quote_id"`
	QuotedAt           *time.Time `json:"quoted_at,omitempty" db:"quoted_at"`
	
	// Relations
	Company          *Company  `json:"company,omitempty"`
	Customer         *Customer `json:"customer,omitempty"`
	Sales            *Account  `json:"sales,omitempty"`
	AssignedEngineer *Account  `json:"assigned_engineer,omitempty"`
	Quote            *Quote    `json:"quote,omitempty"`
}

// InquiryStatus constants
const (
	InquiryStatusDraft          = "draft"
	InquiryStatusPending        = "pending"
	InquiryStatusAssigned       = "assigned"
	InquiryStatusInProgress     = "in_progress"
	InquiryStatusUnderReview    = "under_review"
	InquiryStatusApproved       = "approved"
	InquiryStatusQuoted         = "quoted"
	InquiryStatusRejected       = "rejected"
	InquiryStatusCancelled      = "cancelled"
)

// CreateInquiryRequest represents inquiry creation request
type CreateInquiryRequest struct {
	CustomerID          uuid.UUID `json:"customer_id" validate:"required"`
	ProductCategory     string    `json:"product_category" validate:"required"`
	ProductName         string    `json:"product_name" validate:"required"`
	DrawingFiles        []string  `json:"drawing_files" validate:"required,min=1"`
	Quantity            int       `json:"quantity" validate:"required,min=1"`
	Unit                string    `json:"unit" validate:"required"`
	RequiredDate        time.Time `json:"required_date" validate:"required"`
	Incoterm            string    `json:"incoterm" validate:"required,oneof=EXW FCA FOB CFR CIF CPT CIP DAP DPU DDP"`
	DestinationPort     *string   `json:"destination_port,omitempty"`
	DestinationAddress  *string   `json:"destination_address,omitempty"`
	PaymentTerms        *string   `json:"payment_terms,omitempty"`
	SpecialRequirements *string   `json:"special_requirements,omitempty"`
}

// AssignEngineerRequest represents engineer assignment request
type AssignEngineerRequest struct {
	EngineerID uuid.UUID `json:"engineer_id" validate:"required"`
	Notes      *string   `json:"notes,omitempty"`
}