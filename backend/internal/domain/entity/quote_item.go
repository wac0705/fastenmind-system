package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// QuoteItem 報價項目實體
type QuoteItem struct {
	ID            uuid.UUID
	ProductID     uuid.UUID
	ProductName   string
	Specification string
	Material      Material
	Quantity      int
	UnitPrice     float64
	TaxRate       float64
	DiscountRate  float64
	LeadTime      time.Duration
	Notes         string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	
	// 計算欄位
	totalPrice    float64
	taxAmount     float64
	discountAmount float64
}

// NewQuoteItem 創建新的報價項目
func NewQuoteItem(productID uuid.UUID, productName string, quantity int, unitPrice float64) (*QuoteItem, error) {
	if productID == uuid.Nil {
		return nil, errors.New("product ID is required")
	}
	
	if productName == "" {
		return nil, errors.New("product name is required")
	}
	
	if quantity <= 0 {
		return nil, errors.New("quantity must be positive")
	}
	
	if unitPrice < 0 {
		return nil, errors.New("unit price cannot be negative")
	}
	
	now := time.Now()
	item := &QuoteItem{
		ID:          uuid.New(),
		ProductID:   productID,
		ProductName: productName,
		Quantity:    quantity,
		UnitPrice:   unitPrice,
		TaxRate:     0,
		DiscountRate: 0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	
	item.recalculate()
	
	return item, nil
}

// Validate 驗證報價項目
func (i *QuoteItem) Validate() error {
	if i.ID == uuid.Nil {
		return errors.New("item ID is required")
	}
	
	if i.ProductID == uuid.Nil {
		return errors.New("product ID is required")
	}
	
	if i.ProductName == "" {
		return errors.New("product name is required")
	}
	
	if i.Quantity <= 0 {
		return errors.New("quantity must be positive")
	}
	
	if i.UnitPrice < 0 {
		return errors.New("unit price cannot be negative")
	}
	
	if i.TaxRate < 0 || i.TaxRate > 100 {
		return errors.New("tax rate must be between 0 and 100")
	}
	
	if i.DiscountRate < 0 || i.DiscountRate > 100 {
		return errors.New("discount rate must be between 0 and 100")
	}
	
	if err := i.Material.Validate(); err != nil {
		return err
	}
	
	return nil
}

// UpdateQuantity 更新數量
func (i *QuoteItem) UpdateQuantity(quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}
	
	i.Quantity = quantity
	i.UpdatedAt = time.Now()
	i.recalculate()
	
	return nil
}

// UpdateUnitPrice 更新單價
func (i *QuoteItem) UpdateUnitPrice(unitPrice float64) error {
	if unitPrice < 0 {
		return errors.New("unit price cannot be negative")
	}
	
	i.UnitPrice = unitPrice
	i.UpdatedAt = time.Now()
	i.recalculate()
	
	return nil
}

// SetTaxRate 設置稅率
func (i *QuoteItem) SetTaxRate(rate float64) error {
	if rate < 0 || rate > 100 {
		return errors.New("tax rate must be between 0 and 100")
	}
	
	i.TaxRate = rate
	i.UpdatedAt = time.Now()
	i.recalculate()
	
	return nil
}

// SetDiscountRate 設置折扣率
func (i *QuoteItem) SetDiscountRate(rate float64) error {
	if rate < 0 || rate > 100 {
		return errors.New("discount rate must be between 0 and 100")
	}
	
	i.DiscountRate = rate
	i.UpdatedAt = time.Now()
	i.recalculate()
	
	return nil
}

// SetMaterial 設置材料
func (i *QuoteItem) SetMaterial(material Material) error {
	if err := material.Validate(); err != nil {
		return err
	}
	
	i.Material = material
	i.UpdatedAt = time.Now()
	
	return nil
}

// SetLeadTime 設置交期
func (i *QuoteItem) SetLeadTime(leadTime time.Duration) error {
	if leadTime < 0 {
		return errors.New("lead time cannot be negative")
	}
	
	i.LeadTime = leadTime
	i.UpdatedAt = time.Now()
	
	return nil
}

// CalculateTotal 計算總價
func (i *QuoteItem) CalculateTotal() float64 {
	return i.totalPrice
}

// CalculateTax 計算稅額
func (i *QuoteItem) CalculateTax() float64 {
	return i.taxAmount
}

// CalculateDiscount 計算折扣金額
func (i *QuoteItem) CalculateDiscount() float64 {
	return i.discountAmount
}

// GetNetPrice 獲取淨價（扣除折扣後）
func (i *QuoteItem) GetNetPrice() float64 {
	return i.totalPrice - i.discountAmount
}

// GetFinalPrice 獲取最終價格（含稅）
func (i *QuoteItem) GetFinalPrice() float64 {
	return i.totalPrice + i.taxAmount - i.discountAmount
}

// Clone 複製報價項目
func (i *QuoteItem) Clone() QuoteItem {
	return QuoteItem{
		ID:            uuid.New(),
		ProductID:     i.ProductID,
		ProductName:   i.ProductName,
		Specification: i.Specification,
		Material:      i.Material,
		Quantity:      i.Quantity,
		UnitPrice:     i.UnitPrice,
		TaxRate:       i.TaxRate,
		DiscountRate:  i.DiscountRate,
		LeadTime:      i.LeadTime,
		Notes:         i.Notes,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		totalPrice:    i.totalPrice,
		taxAmount:     i.taxAmount,
		discountAmount: i.discountAmount,
	}
}

// 私有方法

func (i *QuoteItem) recalculate() {
	// 計算基礎總價
	i.totalPrice = float64(i.Quantity) * i.UnitPrice
	
	// 計算折扣金額
	i.discountAmount = i.totalPrice * (i.DiscountRate / 100)
	
	// 計算稅額（基於折扣後金額）
	netAmount := i.totalPrice - i.discountAmount
	i.taxAmount = netAmount * (i.TaxRate / 100)
}

// Material 材料
type Material struct {
	Type        MaterialType
	Grade       string
	Standard    string
	Finish      string
	Description string
}

// MaterialType 材料類型
type MaterialType string

const (
	MaterialTypeSteel         MaterialType = "STEEL"
	MaterialTypeStainlessSteel MaterialType = "STAINLESS_STEEL"
	MaterialTypeBrass         MaterialType = "BRASS"
	MaterialTypeAluminum      MaterialType = "ALUMINUM"
	MaterialTypeTitanium      MaterialType = "TITANIUM"
	MaterialTypePlastic       MaterialType = "PLASTIC"
	MaterialTypeOther         MaterialType = "OTHER"
)

// Validate 驗證材料
func (m Material) Validate() error {
	switch m.Type {
	case MaterialTypeSteel, MaterialTypeStainlessSteel, MaterialTypeBrass,
		 MaterialTypeAluminum, MaterialTypeTitanium, MaterialTypePlastic, MaterialTypeOther:
		// Valid
	default:
		return errors.New("invalid material type")
	}
	
	if m.Grade == "" {
		return errors.New("material grade is required")
	}
	
	return nil
}

// ProductSpecification 產品規格
type ProductSpecification struct {
	Diameter     float64
	Length       float64
	ThreadPitch  float64
	HeadType     string
	DriveType    string
	Coating      string
	Attributes   map[string]string
}

// Validate 驗證產品規格
func (p ProductSpecification) Validate() error {
	if p.Diameter <= 0 {
		return errors.New("diameter must be positive")
	}
	
	if p.Length <= 0 {
		return errors.New("length must be positive")
	}
	
	return nil
}