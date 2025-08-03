package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/fastenmind/fastener-api/internal/domain/aggregate"
	"github.com/fastenmind/fastener-api/internal/domain/valueobject"
)

// QuoteRepository 報價倉儲接口
type QuoteRepository interface {
	// Save 保存報價
	Save(ctx context.Context, quote *aggregate.QuoteAggregate) error
	
	// FindByID 根據ID查詢報價
	FindByID(ctx context.Context, id uuid.UUID) (*aggregate.QuoteAggregate, error)
	
	// FindByNumber 根據報價單號查詢
	FindByNumber(ctx context.Context, quoteNumber valueobject.QuoteNumber) (*aggregate.QuoteAggregate, error)
	
	// FindByCustomer 查詢客戶的報價
	FindByCustomer(ctx context.Context, customerID uuid.UUID, spec QuerySpecification) ([]*aggregate.QuoteAggregate, error)
	
	// FindByCompany 查詢公司的報價
	FindByCompany(ctx context.Context, companyID uuid.UUID, spec QuerySpecification) ([]*aggregate.QuoteAggregate, error)
	
	// FindByStatus 根據狀態查詢報價
	FindByStatus(ctx context.Context, status valueobject.QuoteStatus, spec QuerySpecification) ([]*aggregate.QuoteAggregate, error)
	
	// FindExpiring 查詢即將過期的報價
	FindExpiring(ctx context.Context, withinDays int) ([]*aggregate.QuoteAggregate, error)
	
	// Update 更新報價
	Update(ctx context.Context, quote *aggregate.QuoteAggregate) error
	
	// Delete 刪除報價
	Delete(ctx context.Context, id uuid.UUID) error
	
	// Count 計數
	Count(ctx context.Context, spec QuerySpecification) (int64, error)
	
	// FindBySpecification 根據規格查詢
	FindBySpecification(ctx context.Context, spec Specification) ([]*aggregate.QuoteAggregate, error)
}

// QuerySpecification 查詢規格
type QuerySpecification struct {
	// 分頁
	Limit  int
	Offset int
	
	// 排序
	OrderBy   string
	OrderDesc bool
	
	// 過濾條件
	Filters map[string]interface{}
	
	// 日期範圍
	DateFrom *time.Time
	DateTo   *time.Time
	
	// 包含已刪除
	IncludeDeleted bool
}

// Specification 規格模式接口
type Specification interface {
	IsSatisfiedBy(quote *aggregate.QuoteAggregate) bool
	And(spec Specification) Specification
	Or(spec Specification) Specification
	Not() Specification
}

// BaseSpecification 基礎規格實現
type BaseSpecification struct {
	predicate func(*aggregate.QuoteAggregate) bool
}

// IsSatisfiedBy 檢查是否滿足規格
func (s BaseSpecification) IsSatisfiedBy(quote *aggregate.QuoteAggregate) bool {
	return s.predicate(quote)
}

// And 與操作
func (s BaseSpecification) And(spec Specification) Specification {
	return BaseSpecification{
		predicate: func(quote *aggregate.QuoteAggregate) bool {
			return s.IsSatisfiedBy(quote) && spec.IsSatisfiedBy(quote)
		},
	}
}

// Or 或操作
func (s BaseSpecification) Or(spec Specification) Specification {
	return BaseSpecification{
		predicate: func(quote *aggregate.QuoteAggregate) bool {
			return s.IsSatisfiedBy(quote) || spec.IsSatisfiedBy(quote)
		},
	}
}

// Not 非操作
func (s BaseSpecification) Not() Specification {
	return BaseSpecification{
		predicate: func(quote *aggregate.QuoteAggregate) bool {
			return !s.IsSatisfiedBy(quote)
		},
	}
}

// 預定義規格

// CustomerQuoteSpec 客戶報價規格
func CustomerQuoteSpec(customerID uuid.UUID) Specification {
	return BaseSpecification{
		predicate: func(quote *aggregate.QuoteAggregate) bool {
			return quote.CustomerID == customerID
		},
	}
}

// CompanyQuoteSpec 公司報價規格
func CompanyQuoteSpec(companyID uuid.UUID) Specification {
	return BaseSpecification{
		predicate: func(quote *aggregate.QuoteAggregate) bool {
			return quote.CompanyID == companyID
		},
	}
}

// StatusQuoteSpec 狀態報價規格
func StatusQuoteSpec(status valueobject.QuoteStatus) Specification {
	return BaseSpecification{
		predicate: func(quote *aggregate.QuoteAggregate) bool {
			return quote.Status == status
		},
	}
}

// ActiveQuoteSpec 活躍報價規格
func ActiveQuoteSpec() Specification {
	return BaseSpecification{
		predicate: func(quote *aggregate.QuoteAggregate) bool {
			return quote.Status.IsActive()
		},
	}
}

// ExpiredQuoteSpec 過期報價規格
func ExpiredQuoteSpec() Specification {
	return BaseSpecification{
		predicate: func(quote *aggregate.QuoteAggregate) bool {
			return quote.Status == valueobject.QuoteStatusExpired || time.Now().After(quote.ValidUntil)
		},
	}
}

// DateRangeQuoteSpec 日期範圍報價規格
func DateRangeQuoteSpec(from, to time.Time) Specification {
	return BaseSpecification{
		predicate: func(quote *aggregate.QuoteAggregate) bool {
			return quote.CreatedAt.After(from) && quote.CreatedAt.Before(to)
		},
	}
}

// HighValueQuoteSpec 高價值報價規格
func HighValueQuoteSpec(threshold float64) Specification {
	return BaseSpecification{
		predicate: func(quote *aggregate.QuoteAggregate) bool {
			return quote.PricingSummary.Total > threshold
		},
	}
}