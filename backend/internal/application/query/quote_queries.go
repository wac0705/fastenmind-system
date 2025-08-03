package query

import (
	"time"

	"github.com/google/uuid"
	"github.com/fastenmind/fastener-api/internal/domain/valueobject"
)

// GetQuoteByIDQuery 根據ID獲取報價查詢
type GetQuoteByIDQuery struct {
	BaseQuery
	QuoteID uuid.UUID `json:"quote_id"`
}

// NewGetQuoteByIDQuery 創建根據ID獲取報價查詢
func NewGetQuoteByIDQuery(quoteID uuid.UUID) GetQuoteByIDQuery {
	return GetQuoteByIDQuery{
		BaseQuery: NewBaseQuery("GetQuoteByID"),
		QuoteID:   quoteID,
	}
}

// GetQuoteByNumberQuery 根據編號獲取報價查詢
type GetQuoteByNumberQuery struct {
	BaseQuery
	QuoteNumber string `json:"quote_number"`
}

// NewGetQuoteByNumberQuery 創建根據編號獲取報價查詢
func NewGetQuoteByNumberQuery(quoteNumber string) GetQuoteByNumberQuery {
	return GetQuoteByNumberQuery{
		BaseQuery:   NewBaseQuery("GetQuoteByNumber"),
		QuoteNumber: quoteNumber,
	}
}

// ListQuotesQuery 列出報價查詢
type ListQuotesQuery struct {
	BaseQuery
	CustomerID *uuid.UUID              `json:"customer_id,omitempty"`
	CompanyID  *uuid.UUID              `json:"company_id,omitempty"`
	Status     *valueobject.QuoteStatus `json:"status,omitempty"`
	DateFrom   *time.Time              `json:"date_from,omitempty"`
	DateTo     *time.Time              `json:"date_to,omitempty"`
	PageRequest
}

// NewListQuotesQuery 創建列出報價查詢
func NewListQuotesQuery() ListQuotesQuery {
	return ListQuotesQuery{
		BaseQuery:   NewBaseQuery("ListQuotes"),
		PageRequest: NewPageRequest(1, 20),
	}
}

// SearchQuotesQuery 搜索報價查詢
type SearchQuotesQuery struct {
	BaseQuery
	Keyword      string      `json:"keyword"`
	CustomerName string      `json:"customer_name,omitempty"`
	ProductName  string      `json:"product_name,omitempty"`
	MinAmount    *float64    `json:"min_amount,omitempty"`
	MaxAmount    *float64    `json:"max_amount,omitempty"`
	CompanyID    *uuid.UUID  `json:"company_id,omitempty"`
	PageRequest
}

// NewSearchQuotesQuery 創建搜索報價查詢
func NewSearchQuotesQuery(keyword string) SearchQuotesQuery {
	return SearchQuotesQuery{
		BaseQuery:   NewBaseQuery("SearchQuotes"),
		Keyword:     keyword,
		PageRequest: NewPageRequest(1, 20),
	}
}

// GetExpiringQuotesQuery 獲取即將過期的報價查詢
type GetExpiringQuotesQuery struct {
	BaseQuery
	WithinDays int        `json:"within_days"`
	CompanyID  *uuid.UUID `json:"company_id,omitempty"`
	PageRequest
}

// NewGetExpiringQuotesQuery 創建獲取即將過期的報價查詢
func NewGetExpiringQuotesQuery(withinDays int) GetExpiringQuotesQuery {
	return GetExpiringQuotesQuery{
		BaseQuery:   NewBaseQuery("GetExpiringQuotes"),
		WithinDays:  withinDays,
		PageRequest: NewPageRequest(1, 50),
	}
}

// GetQuoteStatisticsQuery 獲取報價統計查詢
type GetQuoteStatisticsQuery struct {
	BaseQuery
	CompanyID  uuid.UUID  `json:"company_id"`
	DateFrom   time.Time  `json:"date_from"`
	DateTo     time.Time  `json:"date_to"`
}

// NewGetQuoteStatisticsQuery 創建獲取報價統計查詢
func NewGetQuoteStatisticsQuery(companyID uuid.UUID, dateFrom, dateTo time.Time) GetQuoteStatisticsQuery {
	return GetQuoteStatisticsQuery{
		BaseQuery: NewBaseQuery("GetQuoteStatistics"),
		CompanyID: companyID,
		DateFrom:  dateFrom,
		DateTo:    dateTo,
	}
}

// DTO 定義

// QuoteDTO 報價數據傳輸對象
type QuoteDTO struct {
	ID             uuid.UUID                `json:"id"`
	QuoteNumber    string                   `json:"quote_number"`
	CustomerID     uuid.UUID                `json:"customer_id"`
	CustomerName   string                   `json:"customer_name"`
	CompanyID      uuid.UUID                `json:"company_id"`
	CompanyName    string                   `json:"company_name"`
	Status         valueobject.QuoteStatus  `json:"status"`
	ValidUntil     time.Time                `json:"valid_until"`
	Items          []QuoteItemDTO           `json:"items"`
	Terms          QuoteTermsDTO            `json:"terms"`
	PricingSummary PricingSummaryDTO        `json:"pricing_summary"`
	CreatedAt      time.Time                `json:"created_at"`
	UpdatedAt      time.Time                `json:"updated_at"`
	Version        int                      `json:"version"`
}

// QuoteItemDTO 報價項目數據傳輸對象
type QuoteItemDTO struct {
	ID            uuid.UUID    `json:"id"`
	ProductID     uuid.UUID    `json:"product_id"`
	ProductName   string       `json:"product_name"`
	ProductSKU    string       `json:"product_sku"`
	Specification string       `json:"specification"`
	Material      MaterialDTO  `json:"material"`
	Quantity      int          `json:"quantity"`
	UnitPrice     float64      `json:"unit_price"`
	TaxRate       float64      `json:"tax_rate"`
	DiscountRate  float64      `json:"discount_rate"`
	TotalPrice    float64      `json:"total_price"`
	TaxAmount     float64      `json:"tax_amount"`
	DiscountAmount float64     `json:"discount_amount"`
	FinalPrice    float64      `json:"final_price"`
	LeadTime      string       `json:"lead_time"`
	Notes         string       `json:"notes"`
}

// MaterialDTO 材料數據傳輸對象
type MaterialDTO struct {
	Type        string `json:"type"`
	Grade       string `json:"grade"`
	Standard    string `json:"standard"`
	Finish      string `json:"finish"`
	Description string `json:"description"`
}

// QuoteTermsDTO 報價條款數據傳輸對象
type QuoteTermsDTO struct {
	PaymentTerms       PaymentTermsDTO  `json:"payment_terms"`
	DeliveryTerms      DeliveryTermsDTO `json:"delivery_terms"`
	WarrantyTerms      WarrantyTermsDTO `json:"warranty_terms"`
	Currency           string           `json:"currency"`
	DiscountPercentage float64          `json:"discount_percentage"`
	Notes              string           `json:"notes"`
}

// PaymentTermsDTO 付款條款數據傳輸對象
type PaymentTermsDTO struct {
	Type           string  `json:"type"`
	NetDays        int     `json:"net_days"`
	DepositPercent float64 `json:"deposit_percent"`
	Description    string  `json:"description"`
}

// DeliveryTermsDTO 交貨條款數據傳輸對象
type DeliveryTermsDTO struct {
	Incoterm     string `json:"incoterm"`
	LeadTimeDays int    `json:"lead_time_days"`
	Location     string `json:"location"`
	Description  string `json:"description"`
}

// WarrantyTermsDTO 保固條款數據傳輸對象
type WarrantyTermsDTO struct {
	Duration    string   `json:"duration"`
	Type        string   `json:"type"`
	Coverage    string   `json:"coverage"`
	Exclusions  []string `json:"exclusions"`
	Description string   `json:"description"`
}

// PricingSummaryDTO 定價摘要數據傳輸對象
type PricingSummaryDTO struct {
	Subtotal      float64 `json:"subtotal"`
	TotalTax      float64 `json:"total_tax"`
	TotalDiscount float64 `json:"total_discount"`
	Total         float64 `json:"total"`
	Currency      string  `json:"currency"`
}

// QuoteStatisticsDTO 報價統計數據傳輸對象
type QuoteStatisticsDTO struct {
	TotalQuotes       int                        `json:"total_quotes"`
	QuotesByStatus    map[string]int             `json:"quotes_by_status"`
	TotalValue        float64                    `json:"total_value"`
	AverageValue      float64                    `json:"average_value"`
	ConversionRate    float64                    `json:"conversion_rate"`
	TopCustomers      []CustomerStatDTO          `json:"top_customers"`
	TopProducts       []ProductStatDTO           `json:"top_products"`
	MonthlyTrend      []MonthlyTrendDTO          `json:"monthly_trend"`
}

// CustomerStatDTO 客戶統計數據傳輸對象
type CustomerStatDTO struct {
	CustomerID   uuid.UUID `json:"customer_id"`
	CustomerName string    `json:"customer_name"`
	QuoteCount   int       `json:"quote_count"`
	TotalValue   float64   `json:"total_value"`
}

// ProductStatDTO 產品統計數據傳輸對象
type ProductStatDTO struct {
	ProductID   uuid.UUID `json:"product_id"`
	ProductName string    `json:"product_name"`
	ProductSKU  string    `json:"product_sku"`
	QuoteCount  int       `json:"quote_count"`
	Quantity    int       `json:"quantity"`
	TotalValue  float64   `json:"total_value"`
}

// MonthlyTrendDTO 月度趨勢數據傳輸對象
type MonthlyTrendDTO struct {
	Month      string  `json:"month"`
	QuoteCount int     `json:"quote_count"`
	TotalValue float64 `json:"total_value"`
}