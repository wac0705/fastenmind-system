package microservices

import (
	"context"
	"time"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// QuoteServiceServer defines the interface for quote service
type QuoteServiceServer interface {
	CreateQuote(ctx context.Context, req *CreateQuoteRequest) (*CreateQuoteResponse, error)
	GetQuote(ctx context.Context, req *GetQuoteRequest) (*Quote, error)
	ListQuotes(ctx context.Context, req *ListQuotesRequest) (*ListQuotesResponse, error)
	UpdateQuote(ctx context.Context, req *UpdateQuoteRequest) (*Quote, error)
	SubmitQuote(ctx context.Context, req *SubmitQuoteRequest) (*Empty, error)
	ApproveQuote(ctx context.Context, req *ApproveQuoteRequest) (*Empty, error)
	RejectQuote(ctx context.Context, req *RejectQuoteRequest) (*Empty, error)
	GetQuoteStatistics(ctx context.Context, req *GetQuoteStatisticsRequest) (*QuoteStatistics, error)
}

// UnimplementedQuoteServiceServer is used for forward compatibility
type UnimplementedQuoteServiceServer struct{}

func (UnimplementedQuoteServiceServer) CreateQuote(ctx context.Context, req *CreateQuoteRequest) (*CreateQuoteResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method CreateQuote not implemented")
}

func (UnimplementedQuoteServiceServer) GetQuote(ctx context.Context, req *GetQuoteRequest) (*Quote, error) {
	return nil, status.Error(codes.Unimplemented, "method GetQuote not implemented")
}

func (UnimplementedQuoteServiceServer) ListQuotes(ctx context.Context, req *ListQuotesRequest) (*ListQuotesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method ListQuotes not implemented")
}

func (UnimplementedQuoteServiceServer) UpdateQuote(ctx context.Context, req *UpdateQuoteRequest) (*Quote, error) {
	return nil, status.Error(codes.Unimplemented, "method UpdateQuote not implemented")
}

func (UnimplementedQuoteServiceServer) SubmitQuote(ctx context.Context, req *SubmitQuoteRequest) (*Empty, error) {
	return nil, status.Error(codes.Unimplemented, "method SubmitQuote not implemented")
}

func (UnimplementedQuoteServiceServer) ApproveQuote(ctx context.Context, req *ApproveQuoteRequest) (*Empty, error) {
	return nil, status.Error(codes.Unimplemented, "method ApproveQuote not implemented")
}

func (UnimplementedQuoteServiceServer) RejectQuote(ctx context.Context, req *RejectQuoteRequest) (*Empty, error) {
	return nil, status.Error(codes.Unimplemented, "method RejectQuote not implemented")
}

func (UnimplementedQuoteServiceServer) GetQuoteStatistics(ctx context.Context, req *GetQuoteStatisticsRequest) (*QuoteStatistics, error) {
	return nil, status.Error(codes.Unimplemented, "method GetQuoteStatistics not implemented")
}

// Request/Response types

type CreateQuoteRequest struct {
	CustomerId string      `json:"customer_id"`
	CompanyId  string      `json:"company_id"`
	Items      []QuoteItem `json:"items"`
}

type CreateQuoteResponse struct {
	Id          string `json:"id"`
	QuoteNumber string `json:"quote_number"`
}

type GetQuoteRequest struct {
	Id string `json:"id"`
}

type ListQuotesRequest struct {
	CustomerId string `json:"customer_id"`
	CompanyId  string `json:"company_id"`
	Page       int32  `json:"page"`
	PageSize   int32  `json:"page_size"`
	SortBy     string `json:"sort_by"`
	SortDesc   bool   `json:"sort_desc"`
}

type ListQuotesResponse struct {
	Quotes     []*Quote `json:"quotes"`
	TotalCount int64    `json:"total_count"`
	Page       int32    `json:"page"`
	PageSize   int32    `json:"page_size"`
}

type UpdateQuoteRequest struct {
	Id    string      `json:"id"`
	Items []QuoteItem `json:"items"`
}

type SubmitQuoteRequest struct {
	Id string `json:"id"`
}

type ApproveQuoteRequest struct {
	Id         string `json:"id"`
	ApproverId string `json:"approver_id"`
}

type RejectQuoteRequest struct {
	Id     string `json:"id"`
	Reason string `json:"reason"`
}

type GetQuoteStatisticsRequest struct {
	CompanyId string      `json:"company_id"`
	StartDate *Timestamp  `json:"start_date"`
	EndDate   *Timestamp  `json:"end_date"`
}

// Timestamp represents a timestamp for protobuf compatibility
type Timestamp struct {
	Seconds int64 `json:"seconds"`
	Nanos   int32 `json:"nanos"`
}

// AsTime converts Timestamp to time.Time
func (t *Timestamp) AsTime() time.Time {
	if t == nil {
		return time.Time{}
	}
	return time.Unix(t.Seconds, int64(t.Nanos))
}

// Data types

type Quote struct {
	Id             string          `json:"id"`
	QuoteNumber    string          `json:"quote_number"`
	CustomerId     string          `json:"customer_id"`
	CompanyId      string          `json:"company_id"`
	Status         QuoteStatus     `json:"status"`
	ValidUntil     time.Time       `json:"valid_until"`
	Items          []QuoteItem     `json:"items"`
	Terms          QuoteTerms      `json:"terms"`
	PricingSummary PricingSummary  `json:"pricing_summary"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	Version        int32           `json:"version"`
}

type QuoteItem struct {
	Id            string   `json:"id"`
	ProductId     string   `json:"product_id"`
	ProductName   string   `json:"product_name"`
	Specification string   `json:"specification"`
	Material      Material `json:"material"`
	Quantity      int32    `json:"quantity"`
	UnitPrice     float64  `json:"unit_price"`
	TaxRate       float64  `json:"tax_rate"`
	DiscountRate  float64  `json:"discount_rate"`
	TotalPrice    float64  `json:"total_price"`
	LeadTime      string   `json:"lead_time"`
	Notes         string   `json:"notes"`
}

type Material struct {
	Type        string `json:"type"`
	Grade       string `json:"grade"`
	Standard    string `json:"standard"`
	Finish      string `json:"finish"`
	Description string `json:"description"`
}

type QuoteTerms struct {
	PaymentTerms      PaymentTerms  `json:"payment_terms"`
	DeliveryTerms     DeliveryTerms `json:"delivery_terms"`
	WarrantyTerms     WarrantyTerms `json:"warranty_terms"`
	Currency          string        `json:"currency"`
	DiscountPercentage float64      `json:"discount_percentage"`
	Notes             string        `json:"notes"`
}

type PaymentTerms struct {
	Type           string  `json:"type"`
	NetDays        int32   `json:"net_days"`
	DepositPercent float64 `json:"deposit_percent"`
	Description    string  `json:"description"`
}

type DeliveryTerms struct {
	Incoterm     string `json:"incoterm"`
	LeadTimeDays int32  `json:"lead_time_days"`
	Location     string `json:"location"`
	Description  string `json:"description"`
}

type WarrantyTerms struct {
	Duration    string   `json:"duration"`
	Type        string   `json:"type"`
	Coverage    string   `json:"coverage"`
	Exclusions  []string `json:"exclusions"`
	Description string   `json:"description"`
}

type PricingSummary struct {
	Subtotal      float64 `json:"subtotal"`
	TotalTax      float64 `json:"total_tax"`
	TotalDiscount float64 `json:"total_discount"`
	Total         float64 `json:"total"`
	Currency      string  `json:"currency"`
}

type QuoteStatistics struct {
	TotalQuotes     int32                 `json:"total_quotes"`
	QuotesByStatus  map[string]int32      `json:"quotes_by_status"`
	TotalValue      float64               `json:"total_value"`
	AverageValue    float64               `json:"average_value"`
	ConversionRate  float64               `json:"conversion_rate"`
	TopCustomers    []CustomerStat        `json:"top_customers"`
	TopProducts     []ProductStat         `json:"top_products"`
}

type CustomerStat struct {
	CustomerId   string  `json:"customer_id"`
	CustomerName string  `json:"customer_name"`
	QuoteCount   int32   `json:"quote_count"`
	TotalValue   float64 `json:"total_value"`
}

type ProductStat struct {
	ProductId     string  `json:"product_id"`
	ProductName   string  `json:"product_name"`
	ProductSku    string  `json:"product_sku"`
	QuoteCount    int32   `json:"quote_count"`
	TotalQuantity int32   `json:"total_quantity"`
	TotalValue    float64 `json:"total_value"`
}

type QuoteStatus int32

const (
	QuoteStatus_QUOTE_STATUS_UNSPECIFIED QuoteStatus = 0
	QuoteStatus_QUOTE_STATUS_DRAFT       QuoteStatus = 1
	QuoteStatus_QUOTE_STATUS_PENDING     QuoteStatus = 2
	QuoteStatus_QUOTE_STATUS_APPROVED    QuoteStatus = 3
	QuoteStatus_QUOTE_STATUS_REJECTED    QuoteStatus = 4
	QuoteStatus_QUOTE_STATUS_EXPIRED     QuoteStatus = 5
	QuoteStatus_QUOTE_STATUS_CANCELLED   QuoteStatus = 6
)

type Empty struct {}