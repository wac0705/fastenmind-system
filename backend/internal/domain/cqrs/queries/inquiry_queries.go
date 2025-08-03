package queries

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// GetInquiryByIDQuery retrieves an inquiry by ID
type GetInquiryByIDQuery struct {
	InquiryID uuid.UUID `json:"inquiry_id"`
}

// GetName returns the query name
func (q GetInquiryByIDQuery) GetName() string {
	return "GetInquiryByID"
}

// Validate validates the query
func (q GetInquiryByIDQuery) Validate() error {
	if q.InquiryID == uuid.Nil {
		return errors.New("inquiry ID is required")
	}
	return nil
}

// ListInquiriesQuery retrieves inquiries with filters
type ListInquiriesQuery struct {
	CompanyID    *uuid.UUID `json:"company_id,omitempty"`
	CustomerID   *uuid.UUID `json:"customer_id,omitempty"`
	SalesID      *uuid.UUID `json:"sales_id,omitempty"`
	EngineerID   *uuid.UUID `json:"engineer_id,omitempty"`
	Status       *string    `json:"status,omitempty"`
	DateFrom     *time.Time `json:"date_from,omitempty"`
	DateTo       *time.Time `json:"date_to,omitempty"`
	SearchTerm   string     `json:"search_term,omitempty"`
	Page         int        `json:"page"`
	PageSize     int        `json:"page_size"`
	SortBy       string     `json:"sort_by"`
	SortOrder    string     `json:"sort_order"`
}

// GetName returns the query name
func (q ListInquiriesQuery) GetName() string {
	return "ListInquiries"
}

// Validate validates the query
func (q ListInquiriesQuery) Validate() error {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 || q.PageSize > 100 {
		q.PageSize = 20
	}
	if q.SortBy == "" {
		q.SortBy = "created_at"
	}
	if q.SortOrder != "asc" && q.SortOrder != "desc" {
		q.SortOrder = "desc"
	}
	return nil
}

// GetInquiryStatisticsQuery retrieves inquiry statistics
type GetInquiryStatisticsQuery struct {
	CompanyID  *uuid.UUID `json:"company_id,omitempty"`
	CustomerID *uuid.UUID `json:"customer_id,omitempty"`
	DateFrom   time.Time  `json:"date_from"`
	DateTo     time.Time  `json:"date_to"`
	GroupBy    string     `json:"group_by"` // daily, weekly, monthly
}

// GetName returns the query name
func (q GetInquiryStatisticsQuery) GetName() string {
	return "GetInquiryStatistics"
}

// Validate validates the query
func (q GetInquiryStatisticsQuery) Validate() error {
	if q.DateFrom.IsZero() || q.DateTo.IsZero() {
		return errors.New("date range is required")
	}
	if q.DateFrom.After(q.DateTo) {
		return errors.New("date from must be before date to")
	}
	if q.GroupBy == "" {
		q.GroupBy = "daily"
	}
	if q.GroupBy != "daily" && q.GroupBy != "weekly" && q.GroupBy != "monthly" {
		return errors.New("invalid group by value")
	}
	return nil
}

// SearchInquiriesQuery performs full-text search on inquiries
type SearchInquiriesQuery struct {
	SearchTerm string     `json:"search_term"`
	CompanyID  *uuid.UUID `json:"company_id,omitempty"`
	Limit      int        `json:"limit"`
}

// GetName returns the query name
func (q SearchInquiriesQuery) GetName() string {
	return "SearchInquiries"
}

// Validate validates the query
func (q SearchInquiriesQuery) Validate() error {
	if q.SearchTerm == "" {
		return errors.New("search term is required")
	}
	if len(q.SearchTerm) < 3 {
		return errors.New("search term must be at least 3 characters")
	}
	if q.Limit <= 0 || q.Limit > 50 {
		q.Limit = 10
	}
	return nil
}

// GetPendingInquiriesQuery retrieves pending inquiries for assignment
type GetPendingInquiriesQuery struct {
	CompanyID    *uuid.UUID `json:"company_id,omitempty"`
	OlderThanHours int      `json:"older_than_hours,omitempty"`
	Limit        int        `json:"limit"`
}

// GetName returns the query name
func (q GetPendingInquiriesQuery) GetName() string {
	return "GetPendingInquiries"
}

// Validate validates the query
func (q GetPendingInquiriesQuery) Validate() error {
	if q.Limit <= 0 || q.Limit > 100 {
		q.Limit = 20
	}
	return nil
}