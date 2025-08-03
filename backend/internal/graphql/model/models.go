package model

import (
	"time"
	"github.com/google/uuid"
)

// Account represents a user account
type Account struct {
	ID        uuid.UUID `json:"id"`
	CompanyID uuid.UUID `json:"company_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
}

// Inquiry represents an inquiry
type Inquiry struct {
	ID                   uuid.UUID `json:"id"`
	CustomerID           uuid.UUID `json:"customer_id"`
	CompanyID            uuid.UUID `json:"company_id"`
	SalesID              *uuid.UUID `json:"sales_id"`
	EngineerID           *uuid.UUID `json:"engineer_id"`
	ProductCategory      string    `json:"product_category"`
	ProductName          string    `json:"product_name"`
	DrawingFiles         []string  `json:"drawing_files"`
	Quantity             int       `json:"quantity"`
	Unit                 string    `json:"unit"`
	RequiredDate         *time.Time `json:"required_date"`
	Incoterm             string    `json:"incoterm"`
	DestinationPort      string    `json:"destination_port"`
	DestinationAddress   string    `json:"destination_address"`
	PaymentTerms         string    `json:"payment_terms"`
	SpecialRequirements  string    `json:"special_requirements"`
	Status               InquiryStatus `json:"status"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// InquiryStatus represents the status of an inquiry
type InquiryStatus string

const (
	InquiryStatusPending  InquiryStatus = "PENDING"
	InquiryStatusAssigned InquiryStatus = "ASSIGNED"
	InquiryStatusQuoted   InquiryStatus = "QUOTED"
	InquiryStatusRejected InquiryStatus = "REJECTED"
	InquiryStatusClosed   InquiryStatus = "CLOSED"
)

// Input types for mutations
type CreateInquiryInput struct {
	CustomerID           uuid.UUID  `json:"customer_id"`
	ProductCategory      string     `json:"product_category"`
	ProductName          string     `json:"product_name"`
	DrawingFiles         []string   `json:"drawing_files"`
	Quantity             int        `json:"quantity"`
	Unit                 string     `json:"unit"`
	RequiredDate         *time.Time `json:"required_date"`
	Incoterm             string     `json:"incoterm"`
	DestinationPort      string     `json:"destination_port"`
	DestinationAddress   string     `json:"destination_address"`
	PaymentTerms         string     `json:"payment_terms"`
	SpecialRequirements  string     `json:"special_requirements"`
}

type UpdateInquiryInput struct {
	ProductName          *string    `json:"product_name"`
	Quantity             *int       `json:"quantity"`
	RequiredDate         *time.Time `json:"required_date"`
	SpecialRequirements  *string    `json:"special_requirements"`
}

// Filter types for queries
type InquiryFilter struct {
	CompanyID    *uuid.UUID     `json:"company_id"`
	CustomerID   *uuid.UUID     `json:"customer_id"`
	SalesID      *uuid.UUID     `json:"sales_id"`
	EngineerID   *uuid.UUID     `json:"engineer_id"`
	Status       *InquiryStatus `json:"status"`
	DateFrom     *time.Time     `json:"date_from"`
	DateTo       *time.Time     `json:"date_to"`
	SearchTerm   *string        `json:"search_term"`
}

// Pagination types
type PageInput struct {
	First  *int    `json:"first"`
	After  *string `json:"after"`
	Before *string `json:"before"`
	Last   *int    `json:"last"`
}

type PageInfo struct {
	HasNextPage     bool   `json:"has_next_page"`
	HasPreviousPage bool   `json:"has_previous_page"`
	StartCursor     string `json:"start_cursor"`
	EndCursor       string `json:"end_cursor"`
}

// Connection types for GraphQL relay-style pagination
type InquiryEdge struct {
	Node   *Inquiry `json:"node"`
	Cursor string   `json:"cursor"`
}

type InquiryConnection struct {
	Edges      []*InquiryEdge `json:"edges"`
	PageInfo   *PageInfo      `json:"page_info"`
	TotalCount int            `json:"total_count"`
}

// Statistics types
type InquiryStatistics struct {
	TotalCount          int     `json:"total_count"`
	PendingCount        int     `json:"pending_count"`
	AssignedCount       int     `json:"assigned_count"`
	QuotedCount         int     `json:"quoted_count"`
	RejectedCount       int     `json:"rejected_count"`
	ConversionRate      float64 `json:"conversion_rate"`
	AverageResponseTime float64 `json:"average_response_time"`
}