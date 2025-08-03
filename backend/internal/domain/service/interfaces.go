package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/fastenmind/fastener-api/internal/domain/event"
)

// PricingService 定價服務接口
type PricingService interface {
	// CalculatePrice 計算價格
	CalculatePrice(ctx context.Context, request PricingRequest) (*PricingResponse, error)
	
	// GetCustomerPricing 獲取客戶專屬定價
	GetCustomerPricing(ctx context.Context, customerID, productID uuid.UUID) (*CustomerPricing, error)
	
	// CalculateBulkDiscount 計算批量折扣
	CalculateBulkDiscount(ctx context.Context, productID uuid.UUID, quantity int) float64
}

// EventPublisher 事件發布器接口
type EventPublisher interface {
	// Publish 發布事件
	Publish(event event.DomainEvent) error
	
	// PublishAsync 異步發布事件
	PublishAsync(event event.DomainEvent)
}

// CustomerPricing 客戶定價
type CustomerPricing struct {
	BasePrice        float64
	DiscountRate     float64
	SpecialTerms     string
	MinOrderQuantity int
}