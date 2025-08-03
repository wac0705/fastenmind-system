package valueobject

import "errors"

// PricingSummary 定價摘要值對象
type PricingSummary struct {
	Subtotal      float64
	TotalTax      float64
	TotalDiscount float64
	Total         float64
	Currency      Currency
}

// Validate 驗證定價摘要
func (p PricingSummary) Validate() error {
	if p.Subtotal < 0 {
		return errors.New("subtotal cannot be negative")
	}
	
	if p.TotalTax < 0 {
		return errors.New("total tax cannot be negative")
	}
	
	if p.TotalDiscount < 0 {
		return errors.New("total discount cannot be negative")
	}
	
	if p.Total < 0 {
		return errors.New("total cannot be negative")
	}
	
	// 驗證計算邏輯
	calculatedTotal := p.Subtotal + p.TotalTax - p.TotalDiscount
	tolerance := 0.01 // 允許的誤差範圍
	
	if abs(calculatedTotal-p.Total) > tolerance {
		return errors.New("total does not match calculation")
	}
	
	if err := p.Currency.Validate(); err != nil {
		return err
	}
	
	return nil
}

// Money 貨幣金額值對象
type Money struct {
	Amount   float64
	Currency Currency
}

// NewMoney 創建貨幣金額
func NewMoney(amount float64, currency Currency) (Money, error) {
	if amount < 0 {
		return Money{}, errors.New("amount cannot be negative")
	}
	
	if err := currency.Validate(); err != nil {
		return Money{}, err
	}
	
	return Money{
		Amount:   amount,
		Currency: currency,
	}, nil
}

// Add 加法
func (m Money) Add(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, errors.New("cannot add money with different currencies")
	}
	
	return Money{
		Amount:   m.Amount + other.Amount,
		Currency: m.Currency,
	}, nil
}

// Subtract 減法
func (m Money) Subtract(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, errors.New("cannot subtract money with different currencies")
	}
	
	if m.Amount < other.Amount {
		return Money{}, errors.New("insufficient amount")
	}
	
	return Money{
		Amount:   m.Amount - other.Amount,
		Currency: m.Currency,
	}, nil
}

// Multiply 乘法
func (m Money) Multiply(factor float64) (Money, error) {
	if factor < 0 {
		return Money{}, errors.New("factor cannot be negative")
	}
	
	return Money{
		Amount:   m.Amount * factor,
		Currency: m.Currency,
	}, nil
}

// Percentage 計算百分比
type Percentage struct {
	value float64
}

// NewPercentage 創建百分比
func NewPercentage(value float64) (Percentage, error) {
	if value < 0 || value > 100 {
		return Percentage{}, errors.New("percentage must be between 0 and 100")
	}
	
	return Percentage{value: value}, nil
}

// Value 獲取百分比值
func (p Percentage) Value() float64 {
	return p.value
}

// ToDecimal 轉換為小數
func (p Percentage) ToDecimal() float64 {
	return p.value / 100
}

// Apply 應用百分比到金額
func (p Percentage) Apply(amount float64) float64 {
	return amount * p.ToDecimal()
}

// 輔助函數
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}