package cqrs

import (
	"context"
	"time"
	
	"github.com/google/uuid"
)

// Query 查詢介面
type Query interface {
	GetID() string
	GetName() string
	GetTimestamp() time.Time
}

// BaseQuery 基礎查詢
type BaseQuery struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
}

func (q *BaseQuery) GetID() string        { return q.ID }
func (q *BaseQuery) GetName() string      { return q.Name }
func (q *BaseQuery) GetTimestamp() time.Time { return q.Timestamp }

// QueryResult 查詢結果介面
type QueryResult interface{}

// QueryHandler 查詢處理器介面
type QueryHandler interface {
	Handle(ctx context.Context, query Query) (QueryResult, error)
}

// QueryBus 查詢匯流排介面
type QueryBus interface {
	Register(queryName string, handler QueryHandler) error
	Send(ctx context.Context, query Query) (QueryResult, error)
}

// GetOrderByIDQuery 根據ID獲取訂單查詢
type GetOrderByIDQuery struct {
	BaseQuery
	OrderID      string `json:"order_id" validate:"required"`
	IncludeItems bool   `json:"include_items"`
}

func NewGetOrderByIDQuery(orderID string) *GetOrderByIDQuery {
	return &GetOrderByIDQuery{
		BaseQuery: BaseQuery{
			ID:        uuid.New().String(),
			Name:      "GetOrderByID",
			Timestamp: time.Now(),
		},
		OrderID:      orderID,
		IncludeItems: true,
	}
}

// ListOrdersQuery 列出訂單查詢
type ListOrdersQuery struct {
	BaseQuery
	CustomerID   string     `json:"customer_id"`
	Status       string     `json:"status"`
	DateFrom     time.Time  `json:"date_from"`
	DateTo       time.Time  `json:"date_to"`
	Pagination   Pagination `json:"pagination"`
	SortBy       string     `json:"sort_by"`
	SortOrder    string     `json:"sort_order"`
}

func NewListOrdersQuery() *ListOrdersQuery {
	return &ListOrdersQuery{
		BaseQuery: BaseQuery{
			ID:        uuid.New().String(),
			Name:      "ListOrders",
			Timestamp: time.Now(),
		},
		Pagination: Pagination{
			Page:  1,
			Limit: 20,
		},
		SortBy:    "created_at",
		SortOrder: "desc",
	}
}

// GetInventoryStatusQuery 獲取庫存狀態查詢
type GetInventoryStatusQuery struct {
	BaseQuery
	ProductID    string   `json:"product_id"`
	WarehouseIDs []string `json:"warehouse_ids"`
	IncludeReserved bool  `json:"include_reserved"`
}

func NewGetInventoryStatusQuery(productID string) *GetInventoryStatusQuery {
	return &GetInventoryStatusQuery{
		BaseQuery: BaseQuery{
			ID:        uuid.New().String(),
			Name:      "GetInventoryStatus",
			Timestamp: time.Now(),
		},
		ProductID: productID,
		IncludeReserved: true,
	}
}

// GetCustomerStatisticsQuery 獲取客戶統計查詢
type GetCustomerStatisticsQuery struct {
	BaseQuery
	CustomerID    string    `json:"customer_id" validate:"required"`
	Period        string    `json:"period"` // daily, weekly, monthly, yearly
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	IncludeOrders bool      `json:"include_orders"`
	IncludeQuotes bool      `json:"include_quotes"`
}

func NewGetCustomerStatisticsQuery(customerID string) *GetCustomerStatisticsQuery {
	return &GetCustomerStatisticsQuery{
		BaseQuery: BaseQuery{
			ID:        uuid.New().String(),
			Name:      "GetCustomerStatistics",
			Timestamp: time.Now(),
		},
		CustomerID:    customerID,
		Period:        "monthly",
		IncludeOrders: true,
		IncludeQuotes: true,
	}
}

// GetEngineerWorkloadQuery 獲取工程師工作負載查詢
type GetEngineerWorkloadQuery struct {
	BaseQuery
	EngineerIDs  []string   `json:"engineer_ids"`
	DateFrom     time.Time  `json:"date_from"`
	DateTo       time.Time  `json:"date_to"`
	GroupBy      string     `json:"group_by"` // day, week, month
}

func NewGetEngineerWorkloadQuery() *GetEngineerWorkloadQuery {
	return &GetEngineerWorkloadQuery{
		BaseQuery: BaseQuery{
			ID:        uuid.New().String(),
			Name:      "GetEngineerWorkload",
			Timestamp: time.Now(),
		},
		DateFrom: time.Now().AddDate(0, -1, 0),
		DateTo:   time.Now(),
		GroupBy:  "week",
	}
}

// SearchProductsQuery 搜尋產品查詢
type SearchProductsQuery struct {
	BaseQuery
	Keyword      string            `json:"keyword"`
	Category     string            `json:"category"`
	PriceRange   PriceRange        `json:"price_range"`
	Attributes   map[string]string `json:"attributes"`
	InStock      *bool             `json:"in_stock"`
	Pagination   Pagination        `json:"pagination"`
}

func NewSearchProductsQuery(keyword string) *SearchProductsQuery {
	return &SearchProductsQuery{
		BaseQuery: BaseQuery{
			ID:        uuid.New().String(),
			Name:      "SearchProducts",
			Timestamp: time.Now(),
		},
		Keyword: keyword,
		Pagination: Pagination{
			Page:  1,
			Limit: 20,
		},
	}
}

// GetCostAnalysisQuery 獲取成本分析查詢
type GetCostAnalysisQuery struct {
	BaseQuery
	ProductIDs    []string  `json:"product_ids"`
	Period        string    `json:"period"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	GroupBy       string    `json:"group_by"`
	IncludeDetail bool      `json:"include_detail"`
}

func NewGetCostAnalysisQuery() *GetCostAnalysisQuery {
	return &GetCostAnalysisQuery{
		BaseQuery: BaseQuery{
			ID:        uuid.New().String(),
			Name:      "GetCostAnalysis",
			Timestamp: time.Now(),
		},
		Period:        "monthly",
		IncludeDetail: true,
	}
}

// GetReportDataQuery 獲取報表數據查詢
type GetReportDataQuery struct {
	BaseQuery
	ReportType   string                 `json:"report_type" validate:"required"`
	Parameters   map[string]interface{} `json:"parameters"`
	Format       string                 `json:"format"` // json, csv, excel
	TimeZone     string                 `json:"timezone"`
}

func NewGetReportDataQuery(reportType string) *GetReportDataQuery {
	return &GetReportDataQuery{
		BaseQuery: BaseQuery{
			ID:        uuid.New().String(),
			Name:      "GetReportData",
			Timestamp: time.Now(),
		},
		ReportType: reportType,
		Format:     "json",
		TimeZone:   "UTC",
		Parameters: make(map[string]interface{}),
	}
}

// Supporting types

// Pagination 分頁
type Pagination struct {
	Page   int `json:"page"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

// PriceRange 價格範圍
type PriceRange struct {
	Min      float64 `json:"min"`
	Max      float64 `json:"max"`
	Currency string  `json:"currency"`
}

// Query Results

// OrderResult 訂單查詢結果
type OrderResult struct {
	ID           string       `json:"id"`
	OrderNo      string       `json:"order_no"`
	CustomerID   string       `json:"customer_id"`
	CustomerName string       `json:"customer_name"`
	Status       string       `json:"status"`
	Items        []OrderItem  `json:"items,omitempty"`
	TotalAmount  float64      `json:"total_amount"`
	Currency     string       `json:"currency"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

// OrderListResult 訂單列表查詢結果
type OrderListResult struct {
	Orders     []OrderResult `json:"orders"`
	Pagination Pagination    `json:"pagination"`
}

// InventoryStatusResult 庫存狀態查詢結果
type InventoryStatusResult struct {
	ProductID    string                    `json:"product_id"`
	ProductName  string                    `json:"product_name"`
	TotalStock   float64                   `json:"total_stock"`
	Available    float64                   `json:"available"`
	Reserved     float64                   `json:"reserved"`
	ByWarehouse  []WarehouseInventory      `json:"by_warehouse"`
	LastUpdated  time.Time                 `json:"last_updated"`
}

// WarehouseInventory 倉庫庫存
type WarehouseInventory struct {
	WarehouseID   string  `json:"warehouse_id"`
	WarehouseName string  `json:"warehouse_name"`
	Stock         float64 `json:"stock"`
	Available     float64 `json:"available"`
	Reserved      float64 `json:"reserved"`
}

// CustomerStatisticsResult 客戶統計查詢結果
type CustomerStatisticsResult struct {
	CustomerID      string                 `json:"customer_id"`
	Period          string                 `json:"period"`
	TotalOrders     int                    `json:"total_orders"`
	TotalQuotes     int                    `json:"total_quotes"`
	TotalRevenue    float64                `json:"total_revenue"`
	AverageOrderValue float64              `json:"average_order_value"`
	ConversionRate  float64                `json:"conversion_rate"`
	TopProducts     []ProductSummary       `json:"top_products"`
	TrendData       []TrendDataPoint       `json:"trend_data"`
}

// ProductSummary 產品摘要
type ProductSummary struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Revenue     float64 `json:"revenue"`
}

// TrendDataPoint 趨勢數據點
type TrendDataPoint struct {
	Date    time.Time `json:"date"`
	Orders  int       `json:"orders"`
	Revenue float64   `json:"revenue"`
}

// EngineerWorkloadResult 工程師工作負載查詢結果
type EngineerWorkloadResult struct {
	Engineers []EngineerLoad `json:"engineers"`
	Summary   WorkloadSummary `json:"summary"`
}

// EngineerLoad 工程師負載
type EngineerLoad struct {
	EngineerID       string           `json:"engineer_id"`
	EngineerName     string           `json:"engineer_name"`
	TotalAssignments int              `json:"total_assignments"`
	CompletedTasks   int              `json:"completed_tasks"`
	PendingTasks     int              `json:"pending_tasks"`
	UtilizationRate  float64          `json:"utilization_rate"`
	WorkloadByPeriod []WorkloadPeriod `json:"workload_by_period"`
}

// WorkloadPeriod 工作負載期間
type WorkloadPeriod struct {
	Period      string  `json:"period"`
	Assignments int     `json:"assignments"`
	Hours       float64 `json:"hours"`
}

// WorkloadSummary 工作負載摘要
type WorkloadSummary struct {
	TotalEngineers   int     `json:"total_engineers"`
	AverageWorkload  float64 `json:"average_workload"`
	PeakPeriod       string  `json:"peak_period"`
	UnderUtilized    int     `json:"under_utilized"`
	OverUtilized     int     `json:"over_utilized"`
}