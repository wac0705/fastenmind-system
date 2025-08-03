package aggregate

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/fastenmind/fastener-api/internal/domain/valueobject"
)

// CustomerAggregate 客戶聚合根
type CustomerAggregate struct {
	ID           uuid.UUID
	Code         string
	Name         string
	CompanyID    uuid.UUID
	ContactInfo  valueobject.ContactInfo
	Address      valueobject.Address
	CreditStatus valueobject.CreditStatus
	CreditLimit  float64
	PaymentTerms valueobject.PaymentTerms
	Currency     valueobject.Currency
	TaxID        string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	
	// 統計信息
	TotalQuotes     int
	TotalOrders     int
	TotalRevenue    float64
	LastOrderDate   *time.Time
	OutstandingDebt float64
}

// NewCustomerAggregate 創建新客戶
func NewCustomerAggregate(code, name string, companyID uuid.UUID) (*CustomerAggregate, error) {
	if code == "" {
		return nil, errors.New("customer code is required")
	}
	
	if name == "" {
		return nil, errors.New("customer name is required")
	}
	
	if companyID == uuid.Nil {
		return nil, errors.New("company ID is required")
	}
	
	now := time.Now()
	return &CustomerAggregate{
		ID:           uuid.New(),
		Code:         code,
		Name:         name,
		CompanyID:    companyID,
		CreditStatus: valueobject.CreditStatusGood,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// CanCreateQuote 檢查是否可以創建報價
func (c *CustomerAggregate) CanCreateQuote() bool {
	if !c.IsActive {
		return false
	}
	
	if c.CreditStatus == valueobject.CreditStatusBlocked {
		return false
	}
	
	if c.CreditStatus == valueobject.CreditStatusWarning && c.OutstandingDebt > c.CreditLimit {
		return false
	}
	
	return true
}

// UpdateCreditStatus 更新信用狀態
func (c *CustomerAggregate) UpdateCreditStatus(status valueobject.CreditStatus) error {
	if err := status.Validate(); err != nil {
		return err
	}
	
	c.CreditStatus = status
	c.UpdatedAt = time.Now()
	
	return nil
}

// UpdateCreditLimit 更新信用額度
func (c *CustomerAggregate) UpdateCreditLimit(limit float64) error {
	if limit < 0 {
		return errors.New("credit limit cannot be negative")
	}
	
	c.CreditLimit = limit
	c.UpdatedAt = time.Now()
	
	return nil
}

// UpdateContactInfo 更新聯絡信息
func (c *CustomerAggregate) UpdateContactInfo(info valueobject.ContactInfo) error {
	if err := info.Validate(); err != nil {
		return err
	}
	
	c.ContactInfo = info
	c.UpdatedAt = time.Now()
	
	return nil
}

// UpdateAddress 更新地址
func (c *CustomerAggregate) UpdateAddress(address valueobject.Address) error {
	if err := address.Validate(); err != nil {
		return err
	}
	
	c.Address = address
	c.UpdatedAt = time.Now()
	
	return nil
}

// Activate 啟用客戶
func (c *CustomerAggregate) Activate() error {
	if c.IsActive {
		return errors.New("customer is already active")
	}
	
	c.IsActive = true
	c.UpdatedAt = time.Now()
	
	return nil
}

// Deactivate 停用客戶
func (c *CustomerAggregate) Deactivate() error {
	if !c.IsActive {
		return errors.New("customer is already inactive")
	}
	
	if c.OutstandingDebt > 0 {
		return errors.New("cannot deactivate customer with outstanding debt")
	}
	
	c.IsActive = false
	c.UpdatedAt = time.Now()
	
	return nil
}

// UpdateStatistics 更新統計信息
func (c *CustomerAggregate) UpdateStatistics(stats CustomerStatistics) {
	c.TotalQuotes = stats.TotalQuotes
	c.TotalOrders = stats.TotalOrders
	c.TotalRevenue = stats.TotalRevenue
	c.LastOrderDate = stats.LastOrderDate
	c.OutstandingDebt = stats.OutstandingDebt
	c.UpdatedAt = time.Now()
}

// CustomerStatistics 客戶統計信息
type CustomerStatistics struct {
	TotalQuotes     int
	TotalOrders     int
	TotalRevenue    float64
	LastOrderDate   *time.Time
	OutstandingDebt float64
}