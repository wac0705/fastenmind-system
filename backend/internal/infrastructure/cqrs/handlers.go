package cqrs

import (
	"context"
	"fmt"
	
	"github.com/fastenmind/fastener-api/internal/models"
	"github.com/fastenmind/fastener-api/internal/repository"
	"github.com/fastenmind/fastener-api/internal/service"
	"github.com/fastenmind/fastener-api/pkg/cqrs"
	"github.com/google/uuid"
)

// CreateOrderCommandHandler 創建訂單命令處理器
type CreateOrderCommandHandler struct {
	orderService    service.OrderService
	messagingService *MessagingService
	eventStore      cqrs.EventStore
}

// NewCreateOrderCommandHandler 創建命令處理器
func NewCreateOrderCommandHandler(
	orderService service.OrderService,
	messagingService *MessagingService,
	eventStore cqrs.EventStore,
) *CreateOrderCommandHandler {
	return &CreateOrderCommandHandler{
		orderService:     orderService,
		messagingService: messagingService,
		eventStore:       eventStore,
	}
}

// Handle 處理創建訂單命令
func (h *CreateOrderCommandHandler) Handle(ctx context.Context, command cqrs.Command) error {
	cmd, ok := command.(*cqrs.CreateOrderCommand)
	if !ok {
		return fmt.Errorf("invalid command type")
	}
	
	// 創建訂單
	order := &models.Order{
		ID:         uuid.New(),
		CustomerID: uuid.MustParse(cmd.CustomerID),
		QuoteID:    uuid.MustParse(cmd.QuoteID),
		Status:     "pending",
		Currency:   cmd.PaymentInfo.Currency,
		Notes:      cmd.Notes,
	}
	
	// 轉換訂單項目
	for _, item := range cmd.Items {
		order.Items = append(order.Items, models.OrderItem{
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			Discount:    item.Discount,
			Tax:         item.Tax,
			Total:       item.Total,
		})
		order.TotalAmount += item.Total
	}
	
	// 保存訂單
	if err := h.orderService.Create(ctx, order); err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}
	
	// 創建事件
	event := cqrs.NewOrderCreatedEvent(
		order.ID.String(),
		order.CustomerID.String(),
		order.TotalAmount,
		order.Currency,
	)
	
	// 保存事件到事件存儲
	if err := h.eventStore.Save(ctx, []cqrs.Event{event}); err != nil {
		return fmt.Errorf("failed to save event: %w", err)
	}
	
	// 發布事件到訊息佇列
	if err := h.messagingService.PublishOrderCreated(
		ctx,
		order.ID.String(),
		order.CustomerID.String(),
		order.TotalAmount,
	); err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}
	
	return nil
}

// UpdateInventoryCommandHandler 更新庫存命令處理器
type UpdateInventoryCommandHandler struct {
	inventoryService service.InventoryService
	eventStore       cqrs.EventStore
}

// NewUpdateInventoryCommandHandler 創建命令處理器
func NewUpdateInventoryCommandHandler(
	inventoryService service.InventoryService,
	eventStore cqrs.EventStore,
) *UpdateInventoryCommandHandler {
	return &UpdateInventoryCommandHandler{
		inventoryService: inventoryService,
		eventStore:       eventStore,
	}
}

// Handle 處理更新庫存命令
func (h *UpdateInventoryCommandHandler) Handle(ctx context.Context, command cqrs.Command) error {
	cmd, ok := command.(*cqrs.UpdateInventoryCommand)
	if !ok {
		return fmt.Errorf("invalid command type")
	}
	
	// 根據操作類型更新庫存
	var err error
	switch cmd.Type {
	case "add":
		err = h.inventoryService.AddStock(ctx, cmd.ProductID, cmd.WarehouseID, cmd.Quantity)
	case "subtract":
		err = h.inventoryService.SubtractStock(ctx, cmd.ProductID, cmd.WarehouseID, cmd.Quantity)
	case "reserve":
		err = h.inventoryService.ReserveStock(ctx, cmd.ProductID, cmd.WarehouseID, cmd.Quantity, cmd.Reference)
	case "release":
		err = h.inventoryService.ReleaseStock(ctx, cmd.ProductID, cmd.WarehouseID, cmd.Quantity, cmd.Reference)
	default:
		return fmt.Errorf("unknown operation type: %s", cmd.Type)
	}
	
	if err != nil {
		return fmt.Errorf("failed to update inventory: %w", err)
	}
	
	// 創建事件
	event := cqrs.NewInventoryUpdatedEvent(
		cmd.ProductID,
		cmd.WarehouseID,
		cmd.Quantity,
		cmd.Type,
	)
	
	// 保存事件
	if err := h.eventStore.Save(ctx, []cqrs.Event{event}); err != nil {
		return fmt.Errorf("failed to save event: %w", err)
	}
	
	return nil
}

// Query Handlers

// GetOrderByIDQueryHandler 獲取訂單查詢處理器
type GetOrderByIDQueryHandler struct {
	orderRepo repository.OrderRepository
}

// NewGetOrderByIDQueryHandler 創建查詢處理器
func NewGetOrderByIDQueryHandler(orderRepo repository.OrderRepository) *GetOrderByIDQueryHandler {
	return &GetOrderByIDQueryHandler{
		orderRepo: orderRepo,
	}
}

// Handle 處理查詢
func (h *GetOrderByIDQueryHandler) Handle(ctx context.Context, query cqrs.Query) (cqrs.QueryResult, error) {
	q, ok := query.(*cqrs.GetOrderByIDQuery)
	if !ok {
		return nil, fmt.Errorf("invalid query type")
	}
	
	// 獲取訂單
	orderID, err := uuid.Parse(q.OrderID)
	if err != nil {
		return nil, fmt.Errorf("invalid order ID: %w", err)
	}
	
	order, err := h.orderRepo.GetByID(orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	
	// 轉換為查詢結果
	result := &cqrs.OrderResult{
		ID:           order.ID.String(),
		OrderNo:      order.OrderNo,
		CustomerID:   order.CustomerID.String(),
		CustomerName: order.Customer.Name,
		Status:       order.Status,
		TotalAmount:  order.TotalAmount,
		Currency:     order.Currency,
		CreatedAt:    order.CreatedAt,
		UpdatedAt:    order.UpdatedAt,
	}
	
	// 包含訂單項目
	if q.IncludeItems {
		for _, item := range order.Items {
			result.Items = append(result.Items, cqrs.OrderItem{
				ProductID:   item.ProductID,
				ProductName: item.ProductName,
				Quantity:    item.Quantity,
				UnitPrice:   item.UnitPrice,
				Discount:    item.Discount,
				Tax:         item.Tax,
				Total:       item.Total,
			})
		}
	}
	
	return result, nil
}

// GetCustomerStatisticsQueryHandler 獲取客戶統計查詢處理器
type GetCustomerStatisticsQueryHandler struct {
	customerService service.CustomerService
}

// NewGetCustomerStatisticsQueryHandler 創建查詢處理器
func NewGetCustomerStatisticsQueryHandler(customerService service.CustomerService) *GetCustomerStatisticsQueryHandler {
	return &GetCustomerStatisticsQueryHandler{
		customerService: customerService,
	}
}

// Handle 處理查詢
func (h *GetCustomerStatisticsQueryHandler) Handle(ctx context.Context, query cqrs.Query) (cqrs.QueryResult, error) {
	q, ok := query.(*cqrs.GetCustomerStatisticsQuery)
	if !ok {
		return nil, fmt.Errorf("invalid query type")
	}
	
	// 獲取統計數據
	stats, err := h.customerService.GetStatistics(ctx, q.CustomerID, q.StartDate, q.EndDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get statistics: %w", err)
	}
	
	// 轉換為查詢結果
	result := &cqrs.CustomerStatisticsResult{
		CustomerID:        q.CustomerID,
		Period:            q.Period,
		TotalOrders:       stats.TotalOrders,
		TotalQuotes:       stats.TotalQuotes,
		TotalRevenue:      stats.TotalRevenue,
		AverageOrderValue: stats.AverageOrderValue,
		ConversionRate:    stats.ConversionRate,
	}
	
	// 轉換熱門產品
	for _, product := range stats.TopProducts {
		result.TopProducts = append(result.TopProducts, cqrs.ProductSummary{
			ProductID:   product.ProductID,
			ProductName: product.ProductName,
			Quantity:    product.Quantity,
			Revenue:     product.Revenue,
		})
	}
	
	// 轉換趨勢數據
	for _, point := range stats.TrendData {
		result.TrendData = append(result.TrendData, cqrs.TrendDataPoint{
			Date:    point.Date,
			Orders:  point.Orders,
			Revenue: point.Revenue,
		})
	}
	
	return result, nil
}

// CQRSService CQRS 服務
type CQRSService struct {
	commandBus cqrs.CommandBus
	queryBus   cqrs.QueryBus
	eventBus   cqrs.EventBus
	eventStore cqrs.EventStore
}

// NewCQRSService 創建 CQRS 服務
func NewCQRSService(
	commandBus cqrs.CommandBus,
	queryBus cqrs.QueryBus,
	eventBus cqrs.EventBus,
	eventStore cqrs.EventStore,
) *CQRSService {
	return &CQRSService{
		commandBus: commandBus,
		queryBus:   queryBus,
		eventBus:   eventBus,
		eventStore: eventStore,
	}
}

// ExecuteCommand 執行命令
func (s *CQRSService) ExecuteCommand(ctx context.Context, command cqrs.Command) error {
	return s.commandBus.Send(ctx, command)
}

// ExecuteQuery 執行查詢
func (s *CQRSService) ExecuteQuery(ctx context.Context, query cqrs.Query) (cqrs.QueryResult, error) {
	return s.queryBus.Send(ctx, query)
}

// PublishEvent 發布事件
func (s *CQRSService) PublishEvent(ctx context.Context, event cqrs.Event) error {
	return s.eventBus.Publish(ctx, event)
}

// GetEvents 獲取事件
func (s *CQRSService) GetEvents(ctx context.Context, aggregateID string, fromVersion int) ([]cqrs.Event, error) {
	return s.eventStore.GetEvents(ctx, aggregateID, fromVersion)
}