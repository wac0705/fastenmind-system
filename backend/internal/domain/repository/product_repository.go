package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/fastenmind/fastener-api/internal/domain/entity"
)

// ProductRepository 產品倉儲接口
type ProductRepository interface {
	// FindByID 根據ID查詢產品
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Product, error)
	
	// FindBySKU 根據SKU查詢產品
	FindBySKU(ctx context.Context, sku string) (*entity.Product, error)
	
	// Save 保存產品
	Save(ctx context.Context, product *entity.Product) error
	
	// Update 更新產品
	Update(ctx context.Context, product *entity.Product) error
	
	// Delete 刪除產品
	Delete(ctx context.Context, id uuid.UUID) error
	
	// FindByCategory 根據分類查詢產品
	FindByCategory(ctx context.Context, categoryID uuid.UUID) ([]*entity.Product, error)
	
	// Search 搜索產品
	Search(ctx context.Context, criteria ProductSearchCriteria) ([]*entity.Product, error)
	
	// CheckAvailability 檢查產品可用性
	CheckAvailability(ctx context.Context, productID uuid.UUID, quantity int) (bool, error)
}

// ProductSearchCriteria 產品搜索條件
type ProductSearchCriteria struct {
	Name         string
	SKU          string
	CategoryID   *uuid.UUID
	MaterialType string
	MinPrice     *float64
	MaxPrice     *float64
	InStock      *bool
	Limit        int
	Offset       int
}