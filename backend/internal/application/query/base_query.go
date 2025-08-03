package query

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Query 查詢接口
type Query interface {
	GetQueryID() uuid.UUID
	GetQueryType() string
	GetTimestamp() time.Time
}

// BaseQuery 基礎查詢
type BaseQuery struct {
	QueryID   uuid.UUID `json:"query_id"`
	QueryType string    `json:"query_type"`
	Timestamp time.Time `json:"timestamp"`
}

// GetQueryID 獲取查詢ID
func (q BaseQuery) GetQueryID() uuid.UUID {
	return q.QueryID
}

// GetQueryType 獲取查詢類型
func (q BaseQuery) GetQueryType() string {
	return q.QueryType
}

// GetTimestamp 獲取時間戳
func (q BaseQuery) GetTimestamp() time.Time {
	return q.Timestamp
}

// NewBaseQuery 創建基礎查詢
func NewBaseQuery(queryType string) BaseQuery {
	return BaseQuery{
		QueryID:   uuid.New(),
		QueryType: queryType,
		Timestamp: time.Now(),
	}
}

// Handler 查詢處理器接口
type Handler[TQuery Query, TResult any] interface {
	Handle(ctx context.Context, query TQuery) (TResult, error)
}

// HandlerFunc 查詢處理器函數類型
type HandlerFunc[TQuery Query, TResult any] func(ctx context.Context, query TQuery) (TResult, error)

// Handle 處理查詢
func (f HandlerFunc[TQuery, TResult]) Handle(ctx context.Context, query TQuery) (TResult, error) {
	return f(ctx, query)
}

// Bus 查詢總線接口
type Bus interface {
	// Register 註冊查詢處理器
	Register(queryType string, handler interface{}) error
	
	// Send 發送查詢
	Send(ctx context.Context, query Query) (interface{}, error)
}

// Middleware 查詢中間件
type Middleware func(next interface{}) interface{}

// PageRequest 分頁請求
type PageRequest struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Sort     string `json:"sort,omitempty"`
	Order    string `json:"order,omitempty"`
}

// NewPageRequest 創建分頁請求
func NewPageRequest(page, pageSize int) PageRequest {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return PageRequest{
		Page:     page,
		PageSize: pageSize,
		Order:    "asc",
	}
}

// GetOffset 獲取偏移量
func (p PageRequest) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

// GetLimit 獲取限制數量
func (p PageRequest) GetLimit() int {
	return p.PageSize
}

// PageResult 分頁結果
type PageResult[T any] struct {
	Items      []T   `json:"items"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// NewPageResult 創建分頁結果
func NewPageResult[T any](items []T, totalItems int64, page, pageSize int) PageResult[T] {
	totalPages := int(totalItems / int64(pageSize))
	if totalItems%int64(pageSize) > 0 {
		totalPages++
	}
	
	return PageResult[T]{
		Items:      items,
		TotalItems: totalItems,
		TotalPages: totalPages,
		Page:       page,
		PageSize:   pageSize,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}

// EmptyPageResult 創建空的分頁結果
func EmptyPageResult[T any](page, pageSize int) PageResult[T] {
	return PageResult[T]{
		Items:      []T{},
		TotalItems: 0,
		TotalPages: 0,
		Page:       page,
		PageSize:   pageSize,
		HasNext:    false,
		HasPrev:    false,
	}
}