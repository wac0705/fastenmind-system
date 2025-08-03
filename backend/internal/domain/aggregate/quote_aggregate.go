package aggregate

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/fastenmind/fastener-api/internal/domain/entity"
	"github.com/fastenmind/fastener-api/internal/domain/valueobject"
	"github.com/fastenmind/fastener-api/internal/domain/event"
)

// QuoteAggregate 報價聚合根
type QuoteAggregate struct {
	// 基本屬性
	ID           uuid.UUID
	QuoteNumber  valueobject.QuoteNumber
	CustomerID   uuid.UUID
	CompanyID    uuid.UUID
	Status       valueobject.QuoteStatus
	ValidUntil   time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	
	// 聚合內實體
	Items        []entity.QuoteItem
	Terms        valueobject.QuoteTerms
	PricingSummary valueobject.PricingSummary
	
	// 領域事件
	domainEvents []event.DomainEvent
	
	// 版本控制
	Version      int
}

// NewQuoteAggregate 創建新的報價聚合
func NewQuoteAggregate(customerID, companyID uuid.UUID) (*QuoteAggregate, error) {
	if customerID == uuid.Nil {
		return nil, errors.New("customer ID is required")
	}
	if companyID == uuid.Nil {
		return nil, errors.New("company ID is required")
	}
	
	now := time.Now()
	quote := &QuoteAggregate{
		ID:          uuid.New(),
		QuoteNumber: valueobject.GenerateQuoteNumber(companyID, now),
		CustomerID:  customerID,
		CompanyID:   companyID,
		Status:      valueobject.QuoteStatusDraft,
		ValidUntil:  now.AddDate(0, 1, 0), // 預設有效期一個月
		CreatedAt:   now,
		UpdatedAt:   now,
		Items:       make([]entity.QuoteItem, 0),
		domainEvents: make([]event.DomainEvent, 0),
		Version:     1,
	}
	
	// 發布領域事件
	quote.addDomainEvent(event.NewQuoteCreated(quote.ID, quote.QuoteNumber.String(), customerID, companyID))
	
	return quote, nil
}

// AddItem 添加報價項目
func (q *QuoteAggregate) AddItem(item entity.QuoteItem) error {
	// 業務規則驗證
	if err := q.validateCanModify(); err != nil {
		return err
	}
	
	if err := item.Validate(); err != nil {
		return err
	}
	
	// 檢查是否重複
	for _, existingItem := range q.Items {
		if existingItem.ProductID == item.ProductID && existingItem.Specification == item.Specification {
			return errors.New("duplicate quote item")
		}
	}
	
	q.Items = append(q.Items, item)
	q.recalculatePricing()
	q.UpdatedAt = time.Now()
	
	// 發布領域事件
	q.addDomainEvent(event.NewQuoteItemAdded(q.ID, item.ID, item.ProductID))
	
	return nil
}

// RemoveItem 移除報價項目
func (q *QuoteAggregate) RemoveItem(itemID uuid.UUID) error {
	if err := q.validateCanModify(); err != nil {
		return err
	}
	
	found := false
	newItems := make([]entity.QuoteItem, 0, len(q.Items)-1)
	for _, item := range q.Items {
		if item.ID != itemID {
			newItems = append(newItems, item)
		} else {
			found = true
		}
	}
	
	if !found {
		return errors.New("quote item not found")
	}
	
	q.Items = newItems
	q.recalculatePricing()
	q.UpdatedAt = time.Now()
	
	// 發布領域事件
	q.addDomainEvent(event.NewQuoteItemRemoved(q.ID, itemID))
	
	return nil
}

// UpdateItem 更新報價項目
func (q *QuoteAggregate) UpdateItem(itemID uuid.UUID, updates func(*entity.QuoteItem) error) error {
	if err := q.validateCanModify(); err != nil {
		return err
	}
	
	for i := range q.Items {
		if q.Items[i].ID == itemID {
			if err := updates(&q.Items[i]); err != nil {
				return err
			}
			
			if err := q.Items[i].Validate(); err != nil {
				return err
			}
			
			q.recalculatePricing()
			q.UpdatedAt = time.Now()
			
			// 發布領域事件
			q.addDomainEvent(event.NewQuoteItemUpdated(q.ID, itemID))
			
			return nil
		}
	}
	
	return errors.New("quote item not found")
}

// Submit 提交報價
func (q *QuoteAggregate) Submit() error {
	if q.Status != valueobject.QuoteStatusDraft {
		return errors.New("only draft quotes can be submitted")
	}
	
	if len(q.Items) == 0 {
		return errors.New("quote must have at least one item")
	}
	
	// 驗證所有項目都有定價
	for _, item := range q.Items {
		if item.UnitPrice <= 0 {
			return errors.New("all items must have pricing")
		}
	}
	
	q.Status = valueobject.QuoteStatusPending
	q.UpdatedAt = time.Now()
	q.Version++
	
	// 發布領域事件
	q.addDomainEvent(event.NewQuoteSubmitted(q.ID, q.QuoteNumber.String()))
	
	return nil
}

// Approve 批准報價
func (q *QuoteAggregate) Approve(approverID uuid.UUID) error {
	if q.Status != valueobject.QuoteStatusPending {
		return errors.New("only pending quotes can be approved")
	}
	
	if time.Now().After(q.ValidUntil) {
		return errors.New("quote has expired")
	}
	
	q.Status = valueobject.QuoteStatusApproved
	q.UpdatedAt = time.Now()
	q.Version++
	
	// 發布領域事件
	q.addDomainEvent(event.NewQuoteApproved(q.ID, q.QuoteNumber.String(), approverID))
	
	return nil
}

// Reject 拒絕報價
func (q *QuoteAggregate) Reject(reason string) error {
	if q.Status != valueobject.QuoteStatusPending {
		return errors.New("only pending quotes can be rejected")
	}
	
	q.Status = valueobject.QuoteStatusRejected
	q.UpdatedAt = time.Now()
	q.Version++
	
	// 發布領域事件
	q.addDomainEvent(event.NewQuoteRejected(q.ID, q.QuoteNumber.String(), reason))
	
	return nil
}

// Expire 設置報價過期
func (q *QuoteAggregate) Expire() error {
	if q.Status == valueobject.QuoteStatusExpired {
		return nil // 已經過期
	}
	
	if q.Status == valueobject.QuoteStatusApproved {
		return errors.New("approved quotes cannot be expired")
	}
	
	q.Status = valueobject.QuoteStatusExpired
	q.UpdatedAt = time.Now()
	q.Version++
	
	// 發布領域事件
	q.addDomainEvent(event.NewQuoteExpired(q.ID, q.QuoteNumber.String()))
	
	return nil
}

// UpdateTerms 更新報價條款
func (q *QuoteAggregate) UpdateTerms(terms valueobject.QuoteTerms) error {
	if err := q.validateCanModify(); err != nil {
		return err
	}
	
	if err := terms.Validate(); err != nil {
		return err
	}
	
	q.Terms = terms
	q.UpdatedAt = time.Now()
	
	// 發布領域事件
	q.addDomainEvent(event.NewQuoteTermsUpdated(q.ID))
	
	return nil
}

// ExtendValidity 延長有效期
func (q *QuoteAggregate) ExtendValidity(newValidUntil time.Time) error {
	if q.Status == valueobject.QuoteStatusExpired {
		return errors.New("expired quotes cannot be extended")
	}
	
	if q.Status == valueobject.QuoteStatusRejected {
		return errors.New("rejected quotes cannot be extended")
	}
	
	if newValidUntil.Before(time.Now()) {
		return errors.New("new validity date must be in the future")
	}
	
	if newValidUntil.Before(q.ValidUntil) {
		return errors.New("new validity date must be after current validity date")
	}
	
	q.ValidUntil = newValidUntil
	q.UpdatedAt = time.Now()
	
	// 發布領域事件
	q.addDomainEvent(event.NewQuoteValidityExtended(q.ID, newValidUntil))
	
	return nil
}

// Clone 複製報價（用於創建新版本）
func (q *QuoteAggregate) Clone() (*QuoteAggregate, error) {
	newQuote, err := NewQuoteAggregate(q.CustomerID, q.CompanyID)
	if err != nil {
		return nil, err
	}
	
	// 複製項目
	for _, item := range q.Items {
		clonedItem := item.Clone()
		newQuote.Items = append(newQuote.Items, clonedItem)
	}
	
	// 複製條款
	newQuote.Terms = q.Terms
	
	// 重新計算定價
	newQuote.recalculatePricing()
	
	// 發布領域事件
	newQuote.addDomainEvent(event.NewQuoteCloned(q.ID, newQuote.ID))
	
	return newQuote, nil
}

// GetDomainEvents 獲取領域事件
func (q *QuoteAggregate) GetDomainEvents() []event.DomainEvent {
	return q.domainEvents
}

// ClearDomainEvents 清除領域事件
func (q *QuoteAggregate) ClearDomainEvents() {
	q.domainEvents = []event.DomainEvent{}
}

// 私有方法

func (q *QuoteAggregate) validateCanModify() error {
	if q.Status != valueobject.QuoteStatusDraft {
		return errors.New("only draft quotes can be modified")
	}
	return nil
}

func (q *QuoteAggregate) recalculatePricing() {
	var subtotal, totalTax, totalDiscount float64
	
	for _, item := range q.Items {
		itemTotal := item.CalculateTotal()
		subtotal += itemTotal
		totalTax += item.CalculateTax()
		totalDiscount += item.CalculateDiscount()
	}
	
	// 應用報價級別的折扣
	if q.Terms.DiscountPercentage > 0 {
		additionalDiscount := subtotal * (q.Terms.DiscountPercentage / 100)
		totalDiscount += additionalDiscount
	}
	
	total := subtotal + totalTax - totalDiscount
	
	q.PricingSummary = valueobject.PricingSummary{
		Subtotal:      subtotal,
		TotalTax:      totalTax,
		TotalDiscount: totalDiscount,
		Total:         total,
		Currency:      q.Terms.Currency,
	}
}

func (q *QuoteAggregate) addDomainEvent(event event.DomainEvent) {
	q.domainEvents = append(q.domainEvents, event)
}

// Specification 模式實現

// QuoteSpecification 報價規格接口
type QuoteSpecification interface {
	IsSatisfiedBy(quote *QuoteAggregate) bool
}

// ActiveQuoteSpecification 活躍報價規格
type ActiveQuoteSpecification struct{}

func (s ActiveQuoteSpecification) IsSatisfiedBy(quote *QuoteAggregate) bool {
	return quote.Status == valueobject.QuoteStatusPending || 
		   quote.Status == valueobject.QuoteStatusApproved
}

// ExpiredQuoteSpecification 過期報價規格
type ExpiredQuoteSpecification struct{}

func (s ExpiredQuoteSpecification) IsSatisfiedBy(quote *QuoteAggregate) bool {
	return quote.Status == valueobject.QuoteStatusExpired || 
		   time.Now().After(quote.ValidUntil)
}

// HighValueQuoteSpecification 高價值報價規格
type HighValueQuoteSpecification struct {
	Threshold float64
}

func (s HighValueQuoteSpecification) IsSatisfiedBy(quote *QuoteAggregate) bool {
	return quote.PricingSummary.Total > s.Threshold
}