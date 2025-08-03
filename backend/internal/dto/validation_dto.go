package dto

import (
	"time"
)

// CreateOrderRequest 創建訂單請求
type CreateOrderRequest struct {
	CustomerID   string       `json:"customer_id" validate:"required,uuid"`
	QuoteID      string       `json:"quote_id" validate:"required,uuid"`
	Items        []OrderItem  `json:"items" validate:"required,min=1,dive"`
	ShippingInfo ShippingInfo `json:"shipping_info" validate:"required"`
	PaymentTerms string       `json:"payment_terms" validate:"required,min=3,max=100"`
	Currency     string       `json:"currency" validate:"required,currency"`
	Notes        string       `json:"notes" validate:"max=500"`
}

// OrderItem 訂單項目
type OrderItem struct {
	ProductID   string  `json:"product_id" validate:"required,uuid"`
	ProductName string  `json:"product_name" validate:"required,min=1,max=200"`
	Quantity    int     `json:"quantity" validate:"required,min=1"`
	UnitPrice   float64 `json:"unit_price" validate:"required,price"`
	Discount    float64 `json:"discount" validate:"percentage"`
}

// ShippingInfo 配送資訊
type ShippingInfo struct {
	Method      string    `json:"method" validate:"required,oneof=air sea land express"`
	Address     string    `json:"address" validate:"required,min=10,max=500"`
	Contact     string    `json:"contact" validate:"required,min=2,max=100"`
	Phone       string    `json:"phone" validate:"required,phone"`
	DeliveryDate time.Time `json:"delivery_date" validate:"required,future_date"`
}

// CreateCustomerRequest 創建客戶請求
type CreateCustomerRequest struct {
	CustomerCode  string  `json:"customer_code" validate:"required,min=3,max=20,alphanum"`
	Name          string  `json:"name" validate:"required,min=2,max=200"`
	Country       string  `json:"country" validate:"required,country"`
	Currency      string  `json:"currency" validate:"required,currency"`
	Address       string  `json:"address" validate:"required,min=10,max=500"`
	Phone         string  `json:"phone" validate:"required,phone"`
	Email         string  `json:"email" validate:"required,email"`
	ContactPerson string  `json:"contact_person" validate:"required,min=2,max=100"`
	CreditLimit   float64 `json:"credit_limit" validate:"min=0"`
	PaymentTerms  string  `json:"payment_terms" validate:"required,min=3,max=50"`
}

// UpdateCustomerRequest 更新客戶請求
type UpdateCustomerRequest struct {
	Name          string  `json:"name" validate:"omitempty,min=2,max=200"`
	Country       string  `json:"country" validate:"omitempty,country"`
	Currency      string  `json:"currency" validate:"omitempty,currency"`
	Address       string  `json:"address" validate:"omitempty,min=10,max=500"`
	Phone         string  `json:"phone" validate:"omitempty,phone"`
	Email         string  `json:"email" validate:"omitempty,email"`
	ContactPerson string  `json:"contact_person" validate:"omitempty,min=2,max=100"`
	CreditLimit   float64 `json:"credit_limit" validate:"omitempty,min=0"`
	PaymentTerms  string  `json:"payment_terms" validate:"omitempty,min=3,max=50"`
}

// CreateInquiryRequest 創建詢價請求
type CreateInquiryRequest struct {
	CustomerID      string                 `json:"customer_id" validate:"required,uuid"`
	ProductCategory string                 `json:"product_category" validate:"required,min=2,max=50"`
	ProductName     string                 `json:"product_name" validate:"required,min=2,max=200"`
	Quantity        int                    `json:"quantity" validate:"required,min=1"`
	Unit            string                 `json:"unit" validate:"required,oneof=pcs kg m ton"`
	Specifications  map[string]interface{} `json:"specifications"`
	RequiredDate    time.Time              `json:"required_date" validate:"required,future_date"`
	Incoterm        string                 `json:"incoterm" validate:"required,oneof=EXW FCA CPT CIP DAT DAP DDP FOB CFR CIF"`
	DeliveryPort    string                 `json:"delivery_port" validate:"max=100"`
	PaymentTerms    string                 `json:"payment_terms" validate:"required,min=3,max=100"`
	Remarks         string                 `json:"remarks" validate:"max=1000"`
}

// CreateQuoteRequest 創建報價請求
type CreateQuoteRequest struct {
	InquiryID    string        `json:"inquiry_id" validate:"required,uuid"`
	MaterialCost float64       `json:"material_cost" validate:"required,min=0"`
	ProcessCost  float64       `json:"process_cost" validate:"required,min=0"`
	OverheadCost float64       `json:"overhead_cost" validate:"required,min=0"`
	ProfitMargin float64       `json:"profit_margin" validate:"required,percentage"`
	Currency     string        `json:"currency" validate:"required,currency"`
	ValidDays    int           `json:"valid_days" validate:"required,min=1,max=365"`
	DeliveryDays int           `json:"delivery_days" validate:"required,min=1,max=365"`
	PaymentTerms string        `json:"payment_terms" validate:"required,min=3,max=100"`
	CostDetails  []CostDetail  `json:"cost_details" validate:"dive"`
	Remarks      string        `json:"remarks" validate:"max=1000"`
}

// CostDetail 成本明細
type CostDetail struct {
	Category    string  `json:"category" validate:"required,oneof=material processing overhead"`
	Description string  `json:"description" validate:"required,min=2,max=200"`
	Quantity    float64 `json:"quantity" validate:"required,min=0"`
	Unit        string  `json:"unit" validate:"required,min=1,max=20"`
	UnitCost    float64 `json:"unit_cost" validate:"required,min=0"`
	TotalCost   float64 `json:"total_cost" validate:"required,min=0"`
}

// ProcessCostCalculationRequest 製程成本計算請求
type ProcessCostCalculationRequest struct {
	ProductID           string                   `json:"product_id" validate:"omitempty,uuid"`
	ProductName         string                   `json:"product_name" validate:"required,min=2,max=200"`
	ProductSpec         map[string]interface{}   `json:"product_spec"`
	MaterialID          string                   `json:"material_id" validate:"required,uuid"`
	MaterialUtilization float64                  `json:"material_utilization" validate:"required,min=1,max=100"`
	Quantity            int                      `json:"quantity" validate:"required,min=1"`
	Processes           []map[string]interface{} `json:"processes"`
	SurfaceTreatment    string                   `json:"surface_treatment" validate:"max=100"`
	OverheadRate        float64                  `json:"overhead_rate" validate:"percentage"`
	ProfitMargin        float64                  `json:"profit_margin" validate:"percentage"`
	TargetCurrency      string                   `json:"target_currency" validate:"omitempty,currency"`
}

// LoginRequest 登入請求
type LoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}

// ChangePasswordRequest 更改密碼請求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6,max=100"`
	Confirm     string `json:"confirm" validate:"required,eqfield=NewPassword"`
}

// FileUploadRequest 檔案上傳請求
type FileUploadRequest struct {
	Category    string `form:"category" validate:"required,oneof=document image attachment"`
	Description string `form:"description" validate:"max=200"`
	Tags        string `form:"tags" validate:"max=100"`
}

// SearchRequest 搜尋請求
type SearchRequest struct {
	Keyword  string `query:"keyword" validate:"required,min=1,max=100"`
	Page     int    `query:"page" validate:"min=1"`
	PageSize int    `query:"page_size" validate:"min=1,max=100"`
	SortBy   string `query:"sort_by" validate:"omitempty,oneof=created_at updated_at name"`
	SortOrder string `query:"sort_order" validate:"omitempty,oneof=asc desc"`
}

// DateRangeRequest 日期範圍請求
type DateRangeRequest struct {
	StartDate string `query:"start_date" validate:"omitempty,datetime=2006-01-02"`
	EndDate   string `query:"end_date" validate:"omitempty,datetime=2006-01-02"`
}

// BatchOperationRequest 批次操作請求
type BatchOperationRequest struct {
	IDs       []string `json:"ids" validate:"required,min=1,dive,uuid"`
	Operation string   `json:"operation" validate:"required,oneof=delete activate deactivate approve reject"`
	Reason    string   `json:"reason" validate:"max=500"`
}