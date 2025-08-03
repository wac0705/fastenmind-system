package command

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/fastenmind/fastener-api/internal/domain/entity"
	"github.com/fastenmind/fastener-api/internal/domain/valueobject"
)

// CreateQuoteCommand 創建報價命令
type CreateQuoteCommand struct {
	BaseCommand
	CustomerID uuid.UUID                    `json:"customer_id"`
	CompanyID  uuid.UUID                    `json:"company_id"`
	Terms      *valueobject.QuoteTerms      `json:"terms,omitempty"`
	Items      []CreateQuoteItemCommand     `json:"items"`
}

// CreateQuoteItemCommand 創建報價項目命令
type CreateQuoteItemCommand struct {
	ProductID     uuid.UUID       `json:"product_id"`
	Quantity      int             `json:"quantity"`
	Specification string          `json:"specification"`
	Material      entity.Material `json:"material"`
}

// NewCreateQuoteCommand 創建新的創建報價命令
func NewCreateQuoteCommand(customerID, companyID uuid.UUID) CreateQuoteCommand {
	return CreateQuoteCommand{
		BaseCommand: NewBaseCommand("CreateQuote"),
		CustomerID:  customerID,
		CompanyID:   companyID,
		Items:       make([]CreateQuoteItemCommand, 0),
	}
}

// Validate 驗證命令
func (c CreateQuoteCommand) Validate() error {
	if c.CustomerID == uuid.Nil {
		return errors.New("customer ID is required")
	}
	
	if c.CompanyID == uuid.Nil {
		return errors.New("company ID is required")
	}
	
	if len(c.Items) == 0 {
		return errors.New("at least one item is required")
	}
	
	for _, item := range c.Items {
		if item.ProductID == uuid.Nil {
			return errors.New("product ID is required for all items")
		}
		
		if item.Quantity <= 0 {
			return errors.New("quantity must be positive for all items")
		}
	}
	
	if c.Terms != nil {
		if err := c.Terms.Validate(); err != nil {
			return err
		}
	}
	
	return nil
}

// UpdateQuoteItemsCommand 更新報價項目命令
type UpdateQuoteItemsCommand struct {
	BaseCommand
	QuoteID uuid.UUID        `json:"quote_id"`
	Updates []ItemUpdate     `json:"updates"`
}

// ItemUpdate 項目更新
type ItemUpdate struct {
	Action    string    `json:"action"` // add, update, remove
	ItemID    uuid.UUID `json:"item_id,omitempty"`
	ProductID uuid.UUID `json:"product_id,omitempty"`
	Quantity  int       `json:"quantity,omitempty"`
	UnitPrice float64   `json:"unit_price,omitempty"`
}

// NewUpdateQuoteItemsCommand 創建更新報價項目命令
func NewUpdateQuoteItemsCommand(quoteID uuid.UUID) UpdateQuoteItemsCommand {
	return UpdateQuoteItemsCommand{
		BaseCommand: NewBaseCommand("UpdateQuoteItems"),
		QuoteID:     quoteID,
		Updates:     make([]ItemUpdate, 0),
	}
}

// Validate 驗證命令
func (c UpdateQuoteItemsCommand) Validate() error {
	if c.QuoteID == uuid.Nil {
		return errors.New("quote ID is required")
	}
	
	if len(c.Updates) == 0 {
		return errors.New("at least one update is required")
	}
	
	for _, update := range c.Updates {
		switch update.Action {
		case "add":
			if update.ProductID == uuid.Nil {
				return errors.New("product ID is required for add action")
			}
			if update.Quantity <= 0 {
				return errors.New("quantity must be positive for add action")
			}
		case "update":
			if update.ItemID == uuid.Nil {
				return errors.New("item ID is required for update action")
			}
		case "remove":
			if update.ItemID == uuid.Nil {
				return errors.New("item ID is required for remove action")
			}
		default:
			return errors.New("invalid action: " + update.Action)
		}
	}
	
	return nil
}

// SubmitQuoteCommand 提交報價命令
type SubmitQuoteCommand struct {
	BaseCommand
	QuoteID uuid.UUID `json:"quote_id"`
}

// NewSubmitQuoteCommand 創建提交報價命令
func NewSubmitQuoteCommand(quoteID uuid.UUID) SubmitQuoteCommand {
	return SubmitQuoteCommand{
		BaseCommand: NewBaseCommand("SubmitQuote"),
		QuoteID:     quoteID,
	}
}

// Validate 驗證命令
func (c SubmitQuoteCommand) Validate() error {
	if c.QuoteID == uuid.Nil {
		return errors.New("quote ID is required")
	}
	return nil
}

// ApproveQuoteCommand 批准報價命令
type ApproveQuoteCommand struct {
	BaseCommand
	QuoteID    uuid.UUID `json:"quote_id"`
	ApproverID uuid.UUID `json:"approver_id"`
}

// NewApproveQuoteCommand 創建批准報價命令
func NewApproveQuoteCommand(quoteID, approverID uuid.UUID) ApproveQuoteCommand {
	return ApproveQuoteCommand{
		BaseCommand: NewBaseCommand("ApproveQuote"),
		QuoteID:     quoteID,
		ApproverID:  approverID,
	}
}

// Validate 驗證命令
func (c ApproveQuoteCommand) Validate() error {
	if c.QuoteID == uuid.Nil {
		return errors.New("quote ID is required")
	}
	if c.ApproverID == uuid.Nil {
		return errors.New("approver ID is required")
	}
	return nil
}

// RejectQuoteCommand 拒絕報價命令
type RejectQuoteCommand struct {
	BaseCommand
	QuoteID uuid.UUID `json:"quote_id"`
	Reason  string    `json:"reason"`
}

// NewRejectQuoteCommand 創建拒絕報價命令
func NewRejectQuoteCommand(quoteID uuid.UUID, reason string) RejectQuoteCommand {
	return RejectQuoteCommand{
		BaseCommand: NewBaseCommand("RejectQuote"),
		QuoteID:     quoteID,
		Reason:      reason,
	}
}

// Validate 驗證命令
func (c RejectQuoteCommand) Validate() error {
	if c.QuoteID == uuid.Nil {
		return errors.New("quote ID is required")
	}
	if c.Reason == "" {
		return errors.New("rejection reason is required")
	}
	return nil
}

// ExtendQuoteValidityCommand 延長報價有效期命令
type ExtendQuoteValidityCommand struct {
	BaseCommand
	QuoteID       uuid.UUID `json:"quote_id"`
	NewValidUntil time.Time `json:"new_valid_until"`
}

// NewExtendQuoteValidityCommand 創建延長報價有效期命令
func NewExtendQuoteValidityCommand(quoteID uuid.UUID, newValidUntil time.Time) ExtendQuoteValidityCommand {
	return ExtendQuoteValidityCommand{
		BaseCommand:   NewBaseCommand("ExtendQuoteValidity"),
		QuoteID:       quoteID,
		NewValidUntil: newValidUntil,
	}
}

// Validate 驗證命令
func (c ExtendQuoteValidityCommand) Validate() error {
	if c.QuoteID == uuid.Nil {
		return errors.New("quote ID is required")
	}
	if c.NewValidUntil.Before(time.Now()) {
		return errors.New("new validity date must be in the future")
	}
	return nil
}

// CloneQuoteCommand 複製報價命令
type CloneQuoteCommand struct {
	BaseCommand
	QuoteID uuid.UUID `json:"quote_id"`
}

// NewCloneQuoteCommand 創建複製報價命令
func NewCloneQuoteCommand(quoteID uuid.UUID) CloneQuoteCommand {
	return CloneQuoteCommand{
		BaseCommand: NewBaseCommand("CloneQuote"),
		QuoteID:     quoteID,
	}
}

// Validate 驗證命令
func (c CloneQuoteCommand) Validate() error {
	if c.QuoteID == uuid.Nil {
		return errors.New("quote ID is required")
	}
	return nil
}