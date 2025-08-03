package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/fastenmind/fastener-api/internal/domain/valueobject"
)

// Product 產品實體
type Product struct {
	ID             uuid.UUID
	SKU            string
	Name           string
	Description    string
	CategoryID     uuid.UUID
	CategoryName   string
	Specifications ProductSpecification
	Material       Material
	BasePrice      float64
	Currency       valueobject.Currency
	UnitOfMeasure  string
	MinOrderQty    int
	LeadTime       time.Duration
	Weight         float64
	Dimensions     Dimensions
	IsActive       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
	
	// 庫存信息
	StockQuantity  int
	ReservedQty    int
	AvailableQty   int
	ReorderPoint   int
	ReorderQty     int
}

// NewProduct 創建新產品
func NewProduct(sku, name string, categoryID uuid.UUID, basePrice float64) (*Product, error) {
	if sku == "" {
		return nil, errors.New("SKU is required")
	}
	
	if name == "" {
		return nil, errors.New("product name is required")
	}
	
	if categoryID == uuid.Nil {
		return nil, errors.New("category ID is required")
	}
	
	if basePrice < 0 {
		return nil, errors.New("base price cannot be negative")
	}
	
	now := time.Now()
	return &Product{
		ID:            uuid.New(),
		SKU:           sku,
		Name:          name,
		CategoryID:    categoryID,
		BasePrice:     basePrice,
		Currency:      valueobject.CurrencyUSD,
		UnitOfMeasure: "PCS",
		MinOrderQty:   1,
		IsActive:      true,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

// Validate 驗證產品
func (p *Product) Validate() error {
	if p.ID == uuid.Nil {
		return errors.New("product ID is required")
	}
	
	if p.SKU == "" {
		return errors.New("SKU is required")
	}
	
	if p.Name == "" {
		return errors.New("product name is required")
	}
	
	if p.CategoryID == uuid.Nil {
		return errors.New("category ID is required")
	}
	
	if p.BasePrice < 0 {
		return errors.New("base price cannot be negative")
	}
	
	if p.MinOrderQty <= 0 {
		return errors.New("minimum order quantity must be positive")
	}
	
	if err := p.Currency.Validate(); err != nil {
		return err
	}
	
	if err := p.Material.Validate(); err != nil {
		return err
	}
	
	return nil
}

// IsAvailable 檢查產品是否可用
func (p *Product) IsAvailable(quantity int) bool {
	if !p.IsActive {
		return false
	}
	
	if quantity < p.MinOrderQty {
		return false
	}
	
	return p.AvailableQty >= quantity
}

// Reserve 預留庫存
func (p *Product) Reserve(quantity int) error {
	if quantity <= 0 {
		return errors.New("reserve quantity must be positive")
	}
	
	if p.AvailableQty < quantity {
		return errors.New("insufficient available quantity")
	}
	
	p.ReservedQty += quantity
	p.AvailableQty -= quantity
	p.UpdatedAt = time.Now()
	
	return nil
}

// Release 釋放預留庫存
func (p *Product) Release(quantity int) error {
	if quantity <= 0 {
		return errors.New("release quantity must be positive")
	}
	
	if p.ReservedQty < quantity {
		return errors.New("insufficient reserved quantity")
	}
	
	p.ReservedQty -= quantity
	p.AvailableQty += quantity
	p.UpdatedAt = time.Now()
	
	return nil
}

// UpdateStock 更新庫存
func (p *Product) UpdateStock(quantity int) error {
	if quantity < 0 {
		return errors.New("stock quantity cannot be negative")
	}
	
	p.StockQuantity = quantity
	p.AvailableQty = quantity - p.ReservedQty
	p.UpdatedAt = time.Now()
	
	return nil
}

// NeedsReorder 檢查是否需要重新訂購
func (p *Product) NeedsReorder() bool {
	return p.StockQuantity <= p.ReorderPoint
}

// UpdatePrice 更新價格
func (p *Product) UpdatePrice(price float64) error {
	if price < 0 {
		return errors.New("price cannot be negative")
	}
	
	p.BasePrice = price
	p.UpdatedAt = time.Now()
	
	return nil
}

// Activate 啟用產品
func (p *Product) Activate() error {
	if p.IsActive {
		return errors.New("product is already active")
	}
	
	p.IsActive = true
	p.UpdatedAt = time.Now()
	
	return nil
}

// Deactivate 停用產品
func (p *Product) Deactivate() error {
	if !p.IsActive {
		return errors.New("product is already inactive")
	}
	
	if p.ReservedQty > 0 {
		return errors.New("cannot deactivate product with reserved stock")
	}
	
	p.IsActive = false
	p.UpdatedAt = time.Now()
	
	return nil
}

// Dimensions 尺寸
type Dimensions struct {
	Length float64
	Width  float64
	Height float64
	Unit   string
}

// Volume 計算體積
func (d Dimensions) Volume() float64 {
	return d.Length * d.Width * d.Height
}

// Validate 驗證尺寸
func (d Dimensions) Validate() error {
	if d.Length <= 0 || d.Width <= 0 || d.Height <= 0 {
		return errors.New("dimensions must be positive")
	}
	
	if d.Unit == "" {
		return errors.New("dimension unit is required")
	}
	
	return nil
}