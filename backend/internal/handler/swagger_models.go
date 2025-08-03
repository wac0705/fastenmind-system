package handler

import (
	"time"

	"github.com/google/uuid"
)

// Swagger request/response models for documentation

// LoginRequest represents login credentials
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required,min=6" example:"password123"`
}

// LoginResponse represents successful login response
type LoginResponse struct {
	Success      bool   `json:"success" example:"true"`
	Message      string `json:"message" example:"Login successful"`
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIs..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIs..."`
	ExpiresIn    int64  `json:"expires_in" example:"3600"`
	User         UserResponse `json:"user"`
}

// UserResponse represents user information
type UserResponse struct {
	ID        uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email     string    `json:"email" example:"user@example.com"`
	Name      string    `json:"name" example:"John Doe"`
	Role      string    `json:"role" example:"admin"`
	CompanyID uuid.UUID `json:"company_id" example:"550e8400-e29b-41d4-a716-446655440001"`
}

// RegisterRequest represents user registration data
type RegisterRequest struct {
	Email           string `json:"email" validate:"required,email" example:"user@example.com"`
	Password        string `json:"password" validate:"required,min=6" example:"password123"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password" example:"password123"`
	Name            string `json:"name" validate:"required" example:"John Doe"`
	CompanyName     string `json:"company_name" validate:"required" example:"FastenMind Corp"`
}

// RegisterResponse represents successful registration response
type RegisterResponse struct {
	Success bool         `json:"success" example:"true"`
	Message string       `json:"message" example:"Registration successful"`
	User    UserResponse `json:"user"`
}

// RefreshRequest represents token refresh request
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required" example:"eyJhbGciOiJIUzI1NiIs..."`
}

// RefreshResponse represents token refresh response
type RefreshResponse struct {
	Success      bool   `json:"success" example:"true"`
	Message      string `json:"message" example:"Token refreshed successfully"`
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIs..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIs..."`
	ExpiresIn    int64  `json:"expires_in" example:"3600"`
}

// CustomerResponse represents customer data
type CustomerResponse struct {
	ID                 uuid.UUID      `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	CustomerCode       string         `json:"customer_code" example:"CUST001"`
	Name               string         `json:"name" example:"ABC Corporation"`
	NameEN             string         `json:"name_en" example:"ABC Corporation"`
	ShortName          string         `json:"short_name" example:"ABC"`
	Country            string         `json:"country" example:"US"`
	Currency           string         `json:"currency" example:"USD"`
	PaymentTerms       int            `json:"payment_terms" example:"30"`
	CreditLimit        float64        `json:"credit_limit" example:"100000.00"`
	TaxID              string         `json:"tax_id" example:"12-3456789"`
	ContactPerson      string         `json:"contact_person" example:"John Smith"`
	ContactPhone       string         `json:"contact_phone" example:"+1-555-123-4567"`
	ContactEmail       string         `json:"contact_email" example:"john@abc.com"`
	ShippingAddress    AddressInfo    `json:"shipping_address"`
	BillingAddress     AddressInfo    `json:"billing_address"`
	Status             string         `json:"status" example:"active"`
	CreatedAt          time.Time      `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt          time.Time      `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

// AddressInfo represents address information
type AddressInfo struct {
	Street     string `json:"street" example:"123 Main St"`
	City       string `json:"city" example:"New York"`
	State      string `json:"state" example:"NY"`
	Country    string `json:"country" example:"US"`
	PostalCode string `json:"postal_code" example:"10001"`
}

// CreateCustomerRequest represents customer creation data
type CreateCustomerRequest struct {
	CustomerCode    string      `json:"customer_code" validate:"required" example:"CUST001"`
	Name            string      `json:"name" validate:"required" example:"ABC Corporation"`
	NameEN          string      `json:"name_en" example:"ABC Corporation"`
	ShortName       string      `json:"short_name" example:"ABC"`
	Country         string      `json:"country" validate:"required,country_code" example:"US"`
	Currency        string      `json:"currency" validate:"required,currency" example:"USD"`
	PaymentTerms    int         `json:"payment_terms" validate:"required,min=0" example:"30"`
	CreditLimit     float64     `json:"credit_limit" validate:"min=0" example:"100000.00"`
	TaxID           string      `json:"tax_id" example:"12-3456789"`
	ContactPerson   string      `json:"contact_person" validate:"required" example:"John Smith"`
	ContactPhone    string      `json:"contact_phone" validate:"required,phone" example:"+1-555-123-4567"`
	ContactEmail    string      `json:"contact_email" validate:"required,email" example:"john@abc.com"`
	ShippingAddress AddressInfo `json:"shipping_address" validate:"required"`
	BillingAddress  AddressInfo `json:"billing_address"`
}

// UpdateCustomerRequest represents customer update data
type UpdateCustomerRequest struct {
	Name            *string      `json:"name,omitempty" example:"ABC Corporation Updated"`
	NameEN          *string      `json:"name_en,omitempty" example:"ABC Corporation Updated"`
	ShortName       *string      `json:"short_name,omitempty" example:"ABC"`
	Country         *string      `json:"country,omitempty" validate:"omitempty,country_code" example:"CA"`
	Currency        *string      `json:"currency,omitempty" validate:"omitempty,currency" example:"CAD"`
	PaymentTerms    *int         `json:"payment_terms,omitempty" validate:"omitempty,min=0" example:"45"`
	CreditLimit     *float64     `json:"credit_limit,omitempty" validate:"omitempty,min=0" example:"150000.00"`
	ContactPerson   *string      `json:"contact_person,omitempty" example:"Jane Doe"`
	ContactPhone    *string      `json:"contact_phone,omitempty" validate:"omitempty,phone" example:"+1-555-987-6543"`
	ContactEmail    *string      `json:"contact_email,omitempty" validate:"omitempty,email" example:"jane@abc.com"`
	ShippingAddress *AddressInfo `json:"shipping_address,omitempty"`
	BillingAddress  *AddressInfo `json:"billing_address,omitempty"`
	Status          *string      `json:"status,omitempty" validate:"omitempty,oneof=active inactive" example:"active"`
}

// SupplierResponse represents supplier data
type SupplierResponse struct {
	ID            uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	SupplierCode  string    `json:"supplier_code" example:"SUP001"`
	Name          string    `json:"name" example:"XYZ Suppliers"`
	Type          string    `json:"type" example:"material"`
	Country       string    `json:"country" example:"CN"`
	Currency      string    `json:"currency" example:"USD"`
	PaymentTerms  int       `json:"payment_terms" example:"60"`
	ContactPerson string    `json:"contact_person" example:"Li Wei"`
	ContactPhone  string    `json:"contact_phone" example:"+86-21-12345678"`
	ContactEmail  string    `json:"contact_email" example:"liwei@xyz.com"`
	Address       AddressInfo `json:"address"`
	Status        string    `json:"status" example:"active"`
	CreatedAt     time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
}

// ProductResponse represents product data
type ProductResponse struct {
	ID           uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	ProductCode  string    `json:"product_code" example:"PROD001"`
	CustomerID   uuid.UUID `json:"customer_id" example:"550e8400-e29b-41d4-a716-446655440001"`
	Name         string    `json:"name" example:"Hex Bolt M10x50"`
	NameEN       string    `json:"name_en" example:"Hex Bolt M10x50"`
	Category     string    `json:"category" example:"bolt"`
	Specifications map[string]interface{} `json:"specifications"`
	DrawingNo    string    `json:"drawing_no" example:"DWG-001"`
	Weight       float64   `json:"weight" example:"0.125"`
	WeightUnit   string    `json:"weight_unit" example:"kg"`
	MaterialGrade string   `json:"material_grade" example:"8.8"`
	Surface      string    `json:"surface" example:"zinc_plated"`
	HSCode       string    `json:"hs_code" example:"73181500"`
	Status       string    `json:"status" example:"active"`
	CreatedAt    time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
}

// CreateQuoteRequest represents quote creation data
type CreateQuoteRequest struct {
	CustomerID   uuid.UUID                 `json:"customer_id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	ValidDays    int                      `json:"valid_days" validate:"required,min=1" example:"30"`
	Currency     string                   `json:"currency" validate:"required,currency" example:"USD"`
	ExchangeRate float64                  `json:"exchange_rate" validate:"required,positive" example:"1.0"`
	Items        []CreateQuoteItemRequest `json:"items" validate:"required,min=1"`
}

// CreateQuoteItemRequest represents quotation item data
type CreateQuoteItemRequest struct {
	ProductID    uuid.UUID `json:"product_id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Quantity     float64   `json:"quantity" validate:"required,positive" example:"10000"`
	QuantityUnit string    `json:"quantity_unit" validate:"required" example:"pcs"`
	TargetPrice  float64   `json:"target_price" validate:"positive" example:"0.50"`
	Notes        string    `json:"notes" example:"Rush order"`
}

// QuoteResponse represents quotation data
type QuoteResponse struct {
	ID              uuid.UUID            `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	QuoteNo     string              `json:"quotation_no" example:"QT-2024-0001"`
	CustomerID      uuid.UUID           `json:"customer_id" example:"550e8400-e29b-41d4-a716-446655440001"`
	Customer        CustomerResponse    `json:"customer"`
	ValidUntil      time.Time           `json:"valid_until" example:"2024-02-01T00:00:00Z"`
	Currency        string              `json:"currency" example:"USD"`
	ExchangeRate    float64             `json:"exchange_rate" example:"1.0"`
	TotalAmount     float64             `json:"total_amount" example:"5000.00"`
	Status          string              `json:"status" example:"draft"`
	Items           []QuoteItemResponse `json:"items"`
	CreatedAt       time.Time           `json:"created_at" example:"2024-01-01T00:00:00Z"`
	SubmittedAt     *time.Time          `json:"submitted_at,omitempty" example:"2024-01-02T00:00:00Z"`
	ApprovedAt      *time.Time          `json:"approved_at,omitempty" example:"2024-01-03T00:00:00Z"`
}

// QuoteItemResponse represents quotation item data
type QuoteItemResponse struct {
	ID           uuid.UUID       `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	ProductID    uuid.UUID       `json:"product_id" example:"550e8400-e29b-41d4-a716-446655440001"`
	Product      ProductResponse `json:"product"`
	Quantity     float64         `json:"quantity" example:"10000"`
	QuantityUnit string          `json:"quantity_unit" example:"pcs"`
	UnitPrice    float64         `json:"unit_price" example:"0.45"`
	TotalPrice   float64         `json:"total_price" example:"4500.00"`
	TargetPrice  float64         `json:"target_price" example:"0.50"`
	Margin       float64         `json:"margin" example:"10.0"`
	Notes        string          `json:"notes" example:"Rush order"`
}

// QuoteCalculationResponse represents quotation calculation result
type QuoteCalculationResponse struct {
	Success bool                         `json:"success" example:"true"`
	Message string                       `json:"message" example:"Calculation completed"`
	Items   []CalculatedQuoteItem    `json:"items"`
	Summary QuoteCalculationSummary  `json:"summary"`
}

// CalculatedQuoteItem represents calculated item details
type CalculatedQuoteItem struct {
	ProductID      uuid.UUID `json:"product_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	MaterialCost   float64   `json:"material_cost" example:"0.20"`
	ProcessingCost float64   `json:"processing_cost" example:"0.15"`
	OtherCost      float64   `json:"other_cost" example:"0.05"`
	TotalCost      float64   `json:"total_cost" example:"0.40"`
	SuggestedPrice float64   `json:"suggested_price" example:"0.45"`
	Margin         float64   `json:"margin" example:"12.5"`
}

// QuoteCalculationSummary represents calculation summary
type QuoteCalculationSummary struct {
	TotalCost      float64 `json:"total_cost" example:"4000.00"`
	TotalPrice     float64 `json:"total_price" example:"4500.00"`
	AverageMargin  float64 `json:"average_margin" example:"12.5"`
}

// ApprovalRequest represents approval/rejection request
type ApprovalRequest struct {
	Action  string `json:"action" validate:"required,oneof=approve reject" example:"approve"`
	Comment string `json:"comment" example:"Approved with conditions"`
}

// ExchangeRateResponse represents exchange rate data
type ExchangeRateResponse struct {
	ID           uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	BaseCurrency string    `json:"base_currency" example:"USD"`
	Currency     string    `json:"currency" example:"EUR"`
	Rate         float64   `json:"rate" example:"0.85"`
	ValidFrom    time.Time `json:"valid_from" example:"2024-01-01T00:00:00Z"`
	ValidTo      time.Time `json:"valid_to" example:"2024-12-31T23:59:59Z"`
	UpdatedAt    time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

// UpdateExchangeRateRequest represents exchange rate update data
type UpdateExchangeRateRequest struct {
	BaseCurrency string    `json:"base_currency" validate:"required,currency" example:"USD"`
	Currency     string    `json:"currency" validate:"required,currency" example:"EUR"`
	Rate         float64   `json:"rate" validate:"required,positive" example:"0.85"`
	ValidFrom    time.Time `json:"valid_from" validate:"required" example:"2024-01-01T00:00:00Z"`
	ValidTo      time.Time `json:"valid_to" validate:"required" example:"2024-12-31T23:59:59Z"`
}

// SalesReportResponse represents sales report data
type SalesReportResponse struct {
	Period      ReportPeriod        `json:"period"`
	Summary     SalesReportSummary  `json:"summary"`
	Details     []SalesReportDetail `json:"details"`
	GeneratedAt time.Time          `json:"generated_at" example:"2024-01-15T10:00:00Z"`
}

// ReportPeriod represents report period
type ReportPeriod struct {
	StartDate time.Time `json:"start_date" example:"2024-01-01T00:00:00Z"`
	EndDate   time.Time `json:"end_date" example:"2024-01-31T23:59:59Z"`
}

// SalesReportSummary represents sales summary
type SalesReportSummary struct {
	TotalQuotes int     `json:"total_quotations" example:"50"`
	TotalAmount     float64 `json:"total_amount" example:"250000.00"`
	AverageAmount   float64 `json:"average_amount" example:"5000.00"`
	TopCustomer     string  `json:"top_customer" example:"ABC Corporation"`
	TopProduct      string  `json:"top_product" example:"Hex Bolt M10x50"`
}

// SalesReportDetail represents sales detail
type SalesReportDetail struct {
	Date         time.Time `json:"date" example:"2024-01-15T00:00:00Z"`
	QuoteNo  string    `json:"quotation_no" example:"QT-2024-0001"`
	CustomerName string    `json:"customer_name" example:"ABC Corporation"`
	Amount       float64   `json:"amount" example:"5000.00"`
	Status       string    `json:"status" example:"approved"`
}

// CostAnalysisResponse represents cost analysis data
type CostAnalysisResponse struct {
	Product         ProductResponse         `json:"product"`
	CostBreakdown   CostBreakdown          `json:"cost_breakdown"`
	PriceHistory    []PriceHistoryItem     `json:"price_history"`
	Recommendations []CostRecommendation   `json:"recommendations"`
}

// CostBreakdown represents cost components
type CostBreakdown struct {
	MaterialCost      float64 `json:"material_cost" example:"0.20"`
	MaterialCostRatio float64 `json:"material_cost_ratio" example:"50.0"`
	ProcessingCost    float64 `json:"processing_cost" example:"0.15"`
	ProcessingCostRatio float64 `json:"processing_cost_ratio" example:"37.5"`
	OtherCost         float64 `json:"other_cost" example:"0.05"`
	OtherCostRatio    float64 `json:"other_cost_ratio" example:"12.5"`
	TotalCost         float64 `json:"total_cost" example:"0.40"`
}

// PriceHistoryItem represents historical price data
type PriceHistoryItem struct {
	Date      time.Time `json:"date" example:"2024-01-01T00:00:00Z"`
	UnitPrice float64   `json:"unit_price" example:"0.45"`
	Quantity  float64   `json:"quantity" example:"10000"`
	Customer  string    `json:"customer" example:"ABC Corporation"`
}

// CostRecommendation represents cost optimization recommendation
type CostRecommendation struct {
	Type        string  `json:"type" example:"material_substitution"`
	Description string  `json:"description" example:"Consider using alternative material grade"`
	Potential   float64 `json:"potential_savings" example:"0.02"`
	Impact      string  `json:"impact" example:"5% cost reduction"`
}