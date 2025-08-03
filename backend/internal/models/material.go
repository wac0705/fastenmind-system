package models

import (
	"time"
)

// MaterialCategory 材料類別
type MaterialCategory struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	CompanyID   string    `json:"company_id" gorm:"index"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ParentID    string    `json:"parent_id,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedBy   string    `json:"created_by"`
	UpdatedBy   string    `json:"updated_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// MaterialPriceHistory 材料價格歷史
type MaterialPriceHistory struct {
	ID           string    `json:"id" gorm:"primaryKey"`
	MaterialID   string    `json:"material_id" gorm:"index"`
	OldPrice     float64   `json:"old_price"`
	NewPrice     float64   `json:"new_price"`
	Price        float64   `json:"price"`
	Currency     string    `json:"currency"`
	SupplierID   string    `json:"supplier_id,omitempty"`
	ValidFrom    time.Time `json:"valid_from"`
	ValidTo      *time.Time `json:"valid_to"`
	Reason       string    `json:"reason"`
	ChangedBy    string    `json:"changed_by"`
	CreatedBy    string    `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
}

// MaterialPriceUpdate 材料價格更新請求
type MaterialPriceUpdate struct {
	MaterialID string    `json:"material_id" validate:"required"`
	Price      float64   `json:"price" validate:"required,min=0"`
	NewPrice   float64   `json:"new_price" validate:"required,min=0"`
	Currency   string    `json:"currency" validate:"required"`
	SupplierID string    `json:"supplier_id,omitempty"`
	ValidFrom  time.Time `json:"valid_from" validate:"required"`
	Reason     string    `json:"reason"`
}

// MaterialStatistics 材料統計
type MaterialStatistics struct {
	TotalMaterials      int64              `json:"total_materials"`
	ActiveMaterials     int64              `json:"active_materials"`
	LowStockMaterials   int64              `json:"low_stock_materials"`
	OutOfStockMaterials int64              `json:"out_of_stock_materials"`
	TotalInventoryValue float64            `json:"total_inventory_value"`
	LastUpdated         time.Time          `json:"last_updated"`
	ByType              map[string]int64   `json:"by_type"`
	PriceTrend          []PriceTrendData   `json:"price_trend"`
}

// PriceTrendData 價格趨勢數據
type PriceTrendData struct {
	Date       time.Time `json:"date"`
	AvgPrice   float64   `json:"avg_price"`
	MinPrice   float64   `json:"min_price"`
	MaxPrice   float64   `json:"max_price"`
	ItemCount  int       `json:"item_count"`
}

// Table names
func (MaterialCategory) TableName() string {
	return "material_categories"
}

func (MaterialPriceHistory) TableName() string {
	return "material_price_histories"
}