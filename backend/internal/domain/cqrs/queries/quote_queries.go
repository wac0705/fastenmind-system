package queries

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// GetQuoteByIDQuery retrieves a quote by ID
type GetQuoteByIDQuery struct {
	QuoteID uuid.UUID `json:"quote_id"`
}

// GetName returns the query name
func (q GetQuoteByIDQuery) GetName() string {
	return "GetQuoteByID"
}

// Validate validates the query
func (q GetQuoteByIDQuery) Validate() error {
	if q.QuoteID == uuid.Nil {
		return errors.New("quote ID is required")
	}
	return nil
}

// ListQuotesQuery retrieves quotes with filters
type ListQuotesQuery struct {
	CompanyID    *uuid.UUID `json:"company_id,omitempty"`
	CustomerID   *uuid.UUID `json:"customer_id,omitempty"`
	SalesID      *uuid.UUID `json:"sales_id,omitempty"`
	EngineerID   *uuid.UUID `json:"engineer_id,omitempty"`
	Status       *string    `json:"status,omitempty"`
	Currency     *string    `json:"currency,omitempty"`
	DateFrom     *time.Time `json:"date_from,omitempty"`
	DateTo       *time.Time `json:"date_to,omitempty"`
	AmountMin    *float64   `json:"amount_min,omitempty"`
	AmountMax    *float64   `json:"amount_max,omitempty"`
	SearchTerm   string     `json:"search_term,omitempty"`
	Page         int        `json:"page"`
	PageSize     int        `json:"page_size"`
	SortBy       string     `json:"sort_by"`
	SortOrder    string     `json:"sort_order"`
}

// GetName returns the query name
func (q ListQuotesQuery) GetName() string {
	return "ListQuotes"
}

// Validate validates the query
func (q ListQuotesQuery) Validate() error {
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
	
	if q.AmountMin != nil && q.AmountMax != nil && *q.AmountMin > *q.AmountMax {
		return errors.New("amount min must be less than amount max")
	}
	
	return nil
}

// GetQuoteRevisionsQuery retrieves all revisions of a quote
type GetQuoteRevisionsQuery struct {
	QuoteID uuid.UUID `json:"quote_id"`
}

// GetName returns the query name
func (q GetQuoteRevisionsQuery) GetName() string {
	return "GetQuoteRevisions"
}

// Validate validates the query
func (q GetQuoteRevisionsQuery) Validate() error {
	if q.QuoteID == uuid.Nil {
		return errors.New("quote ID is required")
	}
	return nil
}

// GetExpiringQuotesQuery retrieves quotes that are about to expire
type GetExpiringQuotesQuery struct {
	CompanyID    *uuid.UUID `json:"company_id,omitempty"`
	DaysToExpire int        `json:"days_to_expire"`
	Limit        int        `json:"limit"`
}

// GetName returns the query name
func (q GetExpiringQuotesQuery) GetName() string {
	return "GetExpiringQuotes"
}

// Validate validates the query
func (q GetExpiringQuotesQuery) Validate() error {
	if q.DaysToExpire <= 0 {
		q.DaysToExpire = 7 // Default to 7 days
	}
	if q.Limit <= 0 || q.Limit > 100 {
		q.Limit = 20
	}
	return nil
}

// GetQuoteStatisticsQuery retrieves quote statistics
type GetQuoteStatisticsQuery struct {
	CompanyID    *uuid.UUID `json:"company_id,omitempty"`
	CustomerID   *uuid.UUID `json:"customer_id,omitempty"`
	SalesID      *uuid.UUID `json:"sales_id,omitempty"`
	DateFrom     time.Time  `json:"date_from"`
	DateTo       time.Time  `json:"date_to"`
	GroupBy      string     `json:"group_by"` // daily, weekly, monthly
	IncludeItems bool       `json:"include_items"`
}

// GetName returns the query name
func (q GetQuoteStatisticsQuery) GetName() string {
	return "GetQuoteStatistics"
}

// Validate validates the query
func (q GetQuoteStatisticsQuery) Validate() error {
	if q.DateFrom.IsZero() || q.DateTo.IsZero() {
		return errors.New("date range is required")
	}
	if q.DateFrom.After(q.DateTo) {
		return errors.New("date from must be before date to")
	}
	if q.GroupBy == "" {
		q.GroupBy = "monthly"
	}
	if q.GroupBy != "daily" && q.GroupBy != "weekly" && q.GroupBy != "monthly" {
		return errors.New("invalid group by value")
	}
	return nil
}

// GetQuoteConversionRateQuery gets the conversion rate of quotes to orders
type GetQuoteConversionRateQuery struct {
	CompanyID  *uuid.UUID `json:"company_id,omitempty"`
	CustomerID *uuid.UUID `json:"customer_id,omitempty"`
	DateFrom   time.Time  `json:"date_from"`
	DateTo     time.Time  `json:"date_to"`
}

// GetName returns the query name
func (q GetQuoteConversionRateQuery) GetName() string {
	return "GetQuoteConversionRate"
}

// Validate validates the query
func (q GetQuoteConversionRateQuery) Validate() error {
	if q.DateFrom.IsZero() || q.DateTo.IsZero() {
		return errors.New("date range is required")
	}
	if q.DateFrom.After(q.DateTo) {
		return errors.New("date from must be before date to")
	}
	return nil
}