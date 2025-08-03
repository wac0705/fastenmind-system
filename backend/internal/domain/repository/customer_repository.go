package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/fastenmind/fastener-api/internal/domain/aggregate"
)

// CustomerRepository 客戶倉儲接口
type CustomerRepository interface {
	// FindByID 根據ID查詢客戶
	FindByID(ctx context.Context, id uuid.UUID) (*aggregate.CustomerAggregate, error)
	
	// FindByCode 根據客戶代碼查詢
	FindByCode(ctx context.Context, code string) (*aggregate.CustomerAggregate, error)
	
	// Save 保存客戶
	Save(ctx context.Context, customer *aggregate.CustomerAggregate) error
	
	// Update 更新客戶
	Update(ctx context.Context, customer *aggregate.CustomerAggregate) error
	
	// Delete 刪除客戶
	Delete(ctx context.Context, id uuid.UUID) error
	
	// FindByCompany 查詢公司的客戶
	FindByCompany(ctx context.Context, companyID uuid.UUID) ([]*aggregate.CustomerAggregate, error)
	
	// Search 搜索客戶
	Search(ctx context.Context, criteria CustomerSearchCriteria) ([]*aggregate.CustomerAggregate, error)
}

// CustomerSearchCriteria 客戶搜索條件
type CustomerSearchCriteria struct {
	CompanyID    *uuid.UUID
	Name         string
	Code         string
	Country      string
	CreditStatus string
	Limit        int
	Offset       int
}