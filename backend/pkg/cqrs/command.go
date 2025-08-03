package cqrs

import (
	"context"
	"time"
	
	"github.com/google/uuid"
)

// Command 命令介面
type Command interface {
	GetID() string
	GetName() string
	GetTimestamp() time.Time
	Validate() error
}

// BaseCommand 基礎命令
type BaseCommand struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
}

func (c *BaseCommand) GetID() string        { return c.ID }
func (c *BaseCommand) GetName() string      { return c.Name }
func (c *BaseCommand) GetTimestamp() time.Time { return c.Timestamp }

// CommandHandler 命令處理器介面
type CommandHandler interface {
	Handle(ctx context.Context, command Command) error
}

// CommandBus 命令匯流排介面
type CommandBus interface {
	Register(commandName string, handler CommandHandler) error
	Send(ctx context.Context, command Command) error
	SendAsync(ctx context.Context, command Command) <-chan error
}

// CommandStore 命令存儲介面
type CommandStore interface {
	Save(command Command) error
	Get(id string) (Command, error)
	List(filter CommandFilter) ([]Command, error)
}

// CommandFilter 命令過濾器
type CommandFilter struct {
	Name      string
	StartTime time.Time
	EndTime   time.Time
	Limit     int
	Offset    int
}

// CreateOrderCommand 創建訂單命令
type CreateOrderCommand struct {
	BaseCommand
	CustomerID   string                 `json:"customer_id" validate:"required"`
	QuoteID      string                 `json:"quote_id" validate:"required"`
	Items        []OrderItem            `json:"items" validate:"required,min=1"`
	ShippingInfo ShippingInfo          `json:"shipping_info" validate:"required"`
	PaymentInfo  PaymentInfo           `json:"payment_info" validate:"required"`
	Notes        string                 `json:"notes"`
	Metadata     map[string]interface{} `json:"metadata"`
}

func NewCreateOrderCommand(customerID, quoteID string, items []OrderItem) *CreateOrderCommand {
	return &CreateOrderCommand{
		BaseCommand: BaseCommand{
			ID:        uuid.New().String(),
			Name:      "CreateOrder",
			Timestamp: time.Now(),
		},
		CustomerID: customerID,
		QuoteID:    quoteID,
		Items:      items,
	}
}

func (c *CreateOrderCommand) Validate() error {
	// 實作驗證邏輯
	return nil
}

// UpdateInventoryCommand 更新庫存命令
type UpdateInventoryCommand struct {
	BaseCommand
	ProductID    string  `json:"product_id" validate:"required"`
	WarehouseID  string  `json:"warehouse_id" validate:"required"`
	Quantity     float64 `json:"quantity" validate:"required"`
	Type         string  `json:"type" validate:"required,oneof=add subtract reserve release"`
	Reference    string  `json:"reference"`
	Reason       string  `json:"reason"`
}

func NewUpdateInventoryCommand(productID, warehouseID string, quantity float64, opType string) *UpdateInventoryCommand {
	return &UpdateInventoryCommand{
		BaseCommand: BaseCommand{
			ID:        uuid.New().String(),
			Name:      "UpdateInventory",
			Timestamp: time.Now(),
		},
		ProductID:   productID,
		WarehouseID: warehouseID,
		Quantity:    quantity,
		Type:        opType,
	}
}

func (c *UpdateInventoryCommand) Validate() error {
	// 實作驗證邏輯
	return nil
}

// AssignEngineerCommand 分派工程師命令
type AssignEngineerCommand struct {
	BaseCommand
	InquiryID      string    `json:"inquiry_id" validate:"required"`
	EngineerID     string    `json:"engineer_id" validate:"required"`
	Priority       string    `json:"priority" validate:"required,oneof=low medium high urgent"`
	EstimatedHours float64   `json:"estimated_hours"`
	DueDate        time.Time `json:"due_date"`
	Notes          string    `json:"notes"`
}

func NewAssignEngineerCommand(inquiryID, engineerID string, priority string) *AssignEngineerCommand {
	return &AssignEngineerCommand{
		BaseCommand: BaseCommand{
			ID:        uuid.New().String(),
			Name:      "AssignEngineer",
			Timestamp: time.Now(),
		},
		InquiryID:  inquiryID,
		EngineerID: engineerID,
		Priority:   priority,
	}
}

func (c *AssignEngineerCommand) Validate() error {
	// 實作驗證邏輯
	return nil
}

// UpdateCustomerCreditCommand 更新客戶信用額度命令
type UpdateCustomerCreditCommand struct {
	BaseCommand
	CustomerID    string  `json:"customer_id" validate:"required"`
	CreditLimit   float64 `json:"credit_limit" validate:"min=0"`
	PaymentTerms  string  `json:"payment_terms"`
	ApprovedBy    string  `json:"approved_by" validate:"required"`
	Reason        string  `json:"reason" validate:"required"`
	EffectiveDate time.Time `json:"effective_date"`
}

func NewUpdateCustomerCreditCommand(customerID string, creditLimit float64, approvedBy string) *UpdateCustomerCreditCommand {
	return &UpdateCustomerCreditCommand{
		BaseCommand: BaseCommand{
			ID:        uuid.New().String(),
			Name:      "UpdateCustomerCredit",
			Timestamp: time.Now(),
		},
		CustomerID:    customerID,
		CreditLimit:   creditLimit,
		ApprovedBy:    approvedBy,
		EffectiveDate: time.Now(),
	}
}

func (c *UpdateCustomerCreditCommand) Validate() error {
	// 實作驗證邏輯
	return nil
}

// CalculateCostCommand 計算成本命令
type CalculateCostCommand struct {
	BaseCommand
	ProductID      string                 `json:"product_id" validate:"required"`
	Quantity       int                    `json:"quantity" validate:"required,min=1"`
	Specifications map[string]interface{} `json:"specifications"`
	Processes      []ProcessStep          `json:"processes"`
	MaterialID     string                 `json:"material_id"`
	Currency       string                 `json:"currency"`
	IncludeOptions CalculationOptions     `json:"include_options"`
}

func NewCalculateCostCommand(productID string, quantity int) *CalculateCostCommand {
	return &CalculateCostCommand{
		BaseCommand: BaseCommand{
			ID:        uuid.New().String(),
			Name:      "CalculateCost",
			Timestamp: time.Now(),
		},
		ProductID: productID,
		Quantity:  quantity,
	}
}

func (c *CalculateCostCommand) Validate() error {
	// 實作驗證邏輯
	return nil
}

// Supporting types

// OrderItem 訂單項目
type OrderItem struct {
	ProductID    string  `json:"product_id"`
	ProductName  string  `json:"product_name"`
	Quantity     int     `json:"quantity"`
	UnitPrice    float64 `json:"unit_price"`
	Discount     float64 `json:"discount"`
	Tax          float64 `json:"tax"`
	Total        float64 `json:"total"`
}

// ShippingInfo 運送資訊
type ShippingInfo struct {
	Method       string    `json:"method"`
	Address      Address   `json:"address"`
	Cost         float64   `json:"cost"`
	EstimatedDate time.Time `json:"estimated_date"`
	Carrier      string    `json:"carrier"`
	TrackingNo   string    `json:"tracking_no"`
}

// PaymentInfo 付款資訊
type PaymentInfo struct {
	Method        string    `json:"method"`
	Terms         string    `json:"terms"`
	DueDate       time.Time `json:"due_date"`
	Currency      string    `json:"currency"`
	ExchangeRate  float64   `json:"exchange_rate"`
}

// Address 地址
type Address struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	PostalCode string `json:"postal_code"`
	Contact    string `json:"contact"`
	Phone      string `json:"phone"`
}

// ProcessStep 製程步驟
type ProcessStep struct {
	StepID      string                 `json:"step_id"`
	ProcessType string                 `json:"process_type"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// CalculationOptions 計算選項
type CalculationOptions struct {
	IncludeMaterial   bool `json:"include_material"`
	IncludeLabor      bool `json:"include_labor"`
	IncludeOverhead   bool `json:"include_overhead"`
	IncludeProfit     bool `json:"include_profit"`
	ProfitMargin      float64 `json:"profit_margin"`
}