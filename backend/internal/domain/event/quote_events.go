package event

import (
	"time"

	"github.com/google/uuid"
)

// QuoteCreated 報價創建事件
type QuoteCreated struct {
	BaseDomainEvent
	QuoteNumber string    `json:"quote_number"`
	CustomerID  uuid.UUID `json:"customer_id"`
	CompanyID   uuid.UUID `json:"company_id"`
}

// NewQuoteCreated 創建報價創建事件
func NewQuoteCreated(quoteID uuid.UUID, quoteNumber string, customerID, companyID uuid.UUID) *QuoteCreated {
	return &QuoteCreated{
		BaseDomainEvent: NewBaseDomainEvent(quoteID, "QuoteCreated"),
		QuoteNumber:     quoteNumber,
		CustomerID:      customerID,
		CompanyID:       companyID,
	}
}

// QuoteItemAdded 報價項目添加事件
type QuoteItemAdded struct {
	BaseDomainEvent
	ItemID    uuid.UUID `json:"item_id"`
	ProductID uuid.UUID `json:"product_id"`
}

// NewQuoteItemAdded 創建報價項目添加事件
func NewQuoteItemAdded(quoteID, itemID, productID uuid.UUID) *QuoteItemAdded {
	return &QuoteItemAdded{
		BaseDomainEvent: NewBaseDomainEvent(quoteID, "QuoteItemAdded"),
		ItemID:          itemID,
		ProductID:       productID,
	}
}

// QuoteItemRemoved 報價項目移除事件
type QuoteItemRemoved struct {
	BaseDomainEvent
	ItemID uuid.UUID `json:"item_id"`
}

// NewQuoteItemRemoved 創建報價項目移除事件
func NewQuoteItemRemoved(quoteID, itemID uuid.UUID) *QuoteItemRemoved {
	return &QuoteItemRemoved{
		BaseDomainEvent: NewBaseDomainEvent(quoteID, "QuoteItemRemoved"),
		ItemID:          itemID,
	}
}

// QuoteItemUpdated 報價項目更新事件
type QuoteItemUpdated struct {
	BaseDomainEvent
	ItemID uuid.UUID `json:"item_id"`
}

// NewQuoteItemUpdated 創建報價項目更新事件
func NewQuoteItemUpdated(quoteID, itemID uuid.UUID) *QuoteItemUpdated {
	return &QuoteItemUpdated{
		BaseDomainEvent: NewBaseDomainEvent(quoteID, "QuoteItemUpdated"),
		ItemID:          itemID,
	}
}

// QuoteSubmitted 報價提交事件
type QuoteSubmitted struct {
	BaseDomainEvent
	QuoteNumber string `json:"quote_number"`
}

// NewQuoteSubmitted 創建報價提交事件
func NewQuoteSubmitted(quoteID uuid.UUID, quoteNumber string) *QuoteSubmitted {
	return &QuoteSubmitted{
		BaseDomainEvent: NewBaseDomainEvent(quoteID, "QuoteSubmitted"),
		QuoteNumber:     quoteNumber,
	}
}

// QuoteApproved 報價批准事件
type QuoteApproved struct {
	BaseDomainEvent
	QuoteNumber string    `json:"quote_number"`
	ApproverID  uuid.UUID `json:"approver_id"`
}

// NewQuoteApproved 創建報價批准事件
func NewQuoteApproved(quoteID uuid.UUID, quoteNumber string, approverID uuid.UUID) *QuoteApproved {
	return &QuoteApproved{
		BaseDomainEvent: NewBaseDomainEvent(quoteID, "QuoteApproved"),
		QuoteNumber:     quoteNumber,
		ApproverID:      approverID,
	}
}

// QuoteRejected 報價拒絕事件
type QuoteRejected struct {
	BaseDomainEvent
	QuoteNumber string `json:"quote_number"`
	Reason      string `json:"reason"`
}

// NewQuoteRejected 創建報價拒絕事件
func NewQuoteRejected(quoteID uuid.UUID, quoteNumber string, reason string) *QuoteRejected {
	return &QuoteRejected{
		BaseDomainEvent: NewBaseDomainEvent(quoteID, "QuoteRejected"),
		QuoteNumber:     quoteNumber,
		Reason:          reason,
	}
}

// QuoteExpired 報價過期事件
type QuoteExpired struct {
	BaseDomainEvent
	QuoteNumber string `json:"quote_number"`
}

// NewQuoteExpired 創建報價過期事件
func NewQuoteExpired(quoteID uuid.UUID, quoteNumber string) *QuoteExpired {
	return &QuoteExpired{
		BaseDomainEvent: NewBaseDomainEvent(quoteID, "QuoteExpired"),
		QuoteNumber:     quoteNumber,
	}
}

// QuoteTermsUpdated 報價條款更新事件
type QuoteTermsUpdated struct {
	BaseDomainEvent
}

// NewQuoteTermsUpdated 創建報價條款更新事件
func NewQuoteTermsUpdated(quoteID uuid.UUID) *QuoteTermsUpdated {
	return &QuoteTermsUpdated{
		BaseDomainEvent: NewBaseDomainEvent(quoteID, "QuoteTermsUpdated"),
	}
}

// QuoteValidityExtended 報價有效期延長事件
type QuoteValidityExtended struct {
	BaseDomainEvent
	NewValidUntil time.Time `json:"new_valid_until"`
}

// NewQuoteValidityExtended 創建報價有效期延長事件
func NewQuoteValidityExtended(quoteID uuid.UUID, newValidUntil time.Time) *QuoteValidityExtended {
	return &QuoteValidityExtended{
		BaseDomainEvent: NewBaseDomainEvent(quoteID, "QuoteValidityExtended"),
		NewValidUntil:   newValidUntil,
	}
}

// QuoteCloned 報價複製事件
type QuoteCloned struct {
	BaseDomainEvent
	OriginalQuoteID uuid.UUID `json:"original_quote_id"`
	NewQuoteID      uuid.UUID `json:"new_quote_id"`
}

// NewQuoteCloned 創建報價複製事件
func NewQuoteCloned(originalQuoteID, newQuoteID uuid.UUID) *QuoteCloned {
	return &QuoteCloned{
		BaseDomainEvent: NewBaseDomainEvent(newQuoteID, "QuoteCloned"),
		OriginalQuoteID: originalQuoteID,
		NewQuoteID:      newQuoteID,
	}
}