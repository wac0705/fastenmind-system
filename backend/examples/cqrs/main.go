package main

import (
	"context"
	"fmt"
	"log"
	"time"
	
	"github.com/fastenmind/fastener-api/internal/config"
	infra "github.com/fastenmind/fastener-api/internal/infrastructure/cqrs"
	"github.com/fastenmind/fastener-api/internal/models"
	"github.com/fastenmind/fastener-api/pkg/cqrs"
	"github.com/fastenmind/fastener-api/pkg/database"
	"github.com/fastenmind/fastener-api/pkg/tracing"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func main() {
	// 載入配置
	cfg := config.New()
	
	// 初始化追蹤
	tracer, err := tracing.NewTracer(tracing.Config{
		ServiceName:    "cqrs-example",
		ServiceVersion: "1.0.0",
		Environment:    "development",
		ExporterType:   "jaeger",
		Endpoint:       "http://localhost:14268/api/traces",
		SamplingRate:   1.0,
		Enabled:        true,
	})
	if err != nil {
		log.Fatal("Failed to initialize tracer:", err)
	}
	defer tracer.Close(context.Background())
	
	// 初始化資料庫
	dbWrapper, err := database.NewWrapper(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer dbWrapper.Close()
	
	// 執行遷移
	if err := migrateDatabase(dbWrapper.GormDB); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	
	// 創建 CQRS 組件
	commandBus := cqrs.NewSimpleCommandBus()
	queryBus := cqrs.NewSimpleQueryBus()
	eventBus := cqrs.NewSimpleEventBus()
	eventStore := cqrs.NewSQLEventStore(dbWrapper.GormDB)
	
	// 註冊命令處理器
	registerCommandHandlers(commandBus)
	
	// 註冊查詢處理器
	registerQueryHandlers(queryBus)
	
	// 註冊事件處理器
	registerEventHandlers(eventBus)
	
	// 創建 CQRS 服務
	cqrsService := infra.NewCQRSService(commandBus, queryBus, eventBus, eventStore)
	
	// 執行示例場景
	ctx := context.Background()
	runExampleScenario(ctx, cqrsService)
}

// migrateDatabase 執行資料庫遷移
func migrateDatabase(db *gorm.DB) error {
	return db.AutoMigrate(
		&cqrs.EventRecord{},
		&models.Order{},
		&models.Customer{},
		&models.Inventory{},
	)
}

// registerCommandHandlers 註冊命令處理器
func registerCommandHandlers(bus cqrs.CommandBus) {
	// 註冊創建訂單命令處理器
	bus.Register("CreateOrder", &MockCreateOrderHandler{})
	
	// 註冊更新庫存命令處理器
	bus.Register("UpdateInventory", &MockUpdateInventoryHandler{})
	
	// 註冊其他命令處理器...
}

// registerQueryHandlers 註冊查詢處理器
func registerQueryHandlers(bus cqrs.QueryBus) {
	// 註冊獲取訂單查詢處理器
	bus.Register("GetOrderByID", &MockGetOrderHandler{})
	
	// 註冊客戶統計查詢處理器
	bus.Register("GetCustomerStatistics", &MockCustomerStatsHandler{})
	
	// 註冊其他查詢處理器...
}

// registerEventHandlers 註冊事件處理器
func registerEventHandlers(bus cqrs.EventBus) {
	// 訂單創建事件處理器
	bus.Subscribe("OrderCreated", &OrderCreatedEventHandler{})
	
	// 庫存更新事件處理器
	bus.Subscribe("InventoryUpdated", &InventoryUpdatedEventHandler{})
}

// runExampleScenario 執行示例場景
func runExampleScenario(ctx context.Context, service *infra.CQRSService) {
	// 開始追蹤
	ctx, span := tracing.StartSpan(ctx, "ExampleScenario")
	defer span.End()
	
	fmt.Println("=== CQRS Example Scenario ===")
	
	// 1. 創建訂單命令
	fmt.Println("\n1. Creating order...")
	createOrderCmd := cqrs.NewCreateOrderCommand(
		"CUST-001",
		"QUOTE-001",
		[]cqrs.OrderItem{
			{
				ProductID:   "PROD-001",
				ProductName: "Hex Bolt M10x50",
				Quantity:    100,
				UnitPrice:   0.5,
				Total:       50.0,
			},
		},
	)
	
	if err := service.ExecuteCommand(ctx, createOrderCmd); err != nil {
		log.Printf("Failed to create order: %v", err)
	} else {
		fmt.Println("Order created successfully!")
	}
	
	// 2. 更新庫存命令
	fmt.Println("\n2. Updating inventory...")
	updateInventoryCmd := cqrs.NewUpdateInventoryCommand(
		"PROD-001",
		"WH-001",
		100,
		"subtract",
	)
	
	if err := service.ExecuteCommand(ctx, updateInventoryCmd); err != nil {
		log.Printf("Failed to update inventory: %v", err)
	} else {
		fmt.Println("Inventory updated successfully!")
	}
	
	// 3. 查詢訂單
	fmt.Println("\n3. Querying order...")
	getOrderQuery := cqrs.NewGetOrderByIDQuery("ORDER-001")
	
	result, err := service.ExecuteQuery(ctx, getOrderQuery)
	if err != nil {
		log.Printf("Failed to get order: %v", err)
	} else {
		order := result.(*cqrs.OrderResult)
		fmt.Printf("Order found: ID=%s, Customer=%s, Total=%.2f\n",
			order.ID, order.CustomerName, order.TotalAmount)
	}
	
	// 4. 查詢客戶統計
	fmt.Println("\n4. Querying customer statistics...")
	statsQuery := cqrs.NewGetCustomerStatisticsQuery("CUST-001")
	
	result, err = service.ExecuteQuery(ctx, statsQuery)
	if err != nil {
		log.Printf("Failed to get statistics: %v", err)
	} else {
		stats := result.(*cqrs.CustomerStatisticsResult)
		fmt.Printf("Customer Stats: Orders=%d, Revenue=%.2f, Avg Order=%.2f\n",
			stats.TotalOrders, stats.TotalRevenue, stats.AverageOrderValue)
	}
	
	// 5. 獲取事件歷史
	fmt.Println("\n5. Getting event history...")
	events, err := service.GetEvents(ctx, "ORDER-001", 0)
	if err != nil {
		log.Printf("Failed to get events: %v", err)
	} else {
		fmt.Printf("Found %d events for ORDER-001\n", len(events))
		for _, event := range events {
			fmt.Printf("  - Event: %s at %s\n", event.GetEventType(), event.GetTimestamp())
		}
	}
}

// Mock Handlers

// MockCreateOrderHandler 模擬創建訂單處理器
type MockCreateOrderHandler struct{}

func (h *MockCreateOrderHandler) Handle(ctx context.Context, command cqrs.Command) error {
	cmd := command.(*cqrs.CreateOrderCommand)
	
	// 模擬創建訂單
	log.Printf("Creating order for customer %s with %d items", cmd.CustomerID, len(cmd.Items))
	
	// 模擬延遲
	time.Sleep(50 * time.Millisecond)
	
	return nil
}

// MockUpdateInventoryHandler 模擬更新庫存處理器
type MockUpdateInventoryHandler struct{}

func (h *MockUpdateInventoryHandler) Handle(ctx context.Context, command cqrs.Command) error {
	cmd := command.(*cqrs.UpdateInventoryCommand)
	
	// 模擬更新庫存
	log.Printf("Updating inventory: Product=%s, Warehouse=%s, Quantity=%.2f, Type=%s",
		cmd.ProductID, cmd.WarehouseID, cmd.Quantity, cmd.Type)
	
	// 模擬延遲
	time.Sleep(30 * time.Millisecond)
	
	return nil
}

// MockGetOrderHandler 模擬獲取訂單處理器
type MockGetOrderHandler struct{}

func (h *MockGetOrderHandler) Handle(ctx context.Context, query cqrs.Query) (cqrs.QueryResult, error) {
	q := query.(*cqrs.GetOrderByIDQuery)
	
	// 模擬查詢訂單
	log.Printf("Getting order: %s", q.OrderID)
	
	// 返回模擬結果
	return &cqrs.OrderResult{
		ID:           q.OrderID,
		OrderNo:      "ORD-2024-001",
		CustomerID:   "CUST-001",
		CustomerName: "ABC Company",
		Status:       "confirmed",
		TotalAmount:  1500.00,
		Currency:     "USD",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

// MockCustomerStatsHandler 模擬客戶統計處理器
type MockCustomerStatsHandler struct{}

func (h *MockCustomerStatsHandler) Handle(ctx context.Context, query cqrs.Query) (cqrs.QueryResult, error) {
	q := query.(*cqrs.GetCustomerStatisticsQuery)
	
	// 模擬查詢統計
	log.Printf("Getting statistics for customer: %s", q.CustomerID)
	
	// 返回模擬結果
	return &cqrs.CustomerStatisticsResult{
		CustomerID:        q.CustomerID,
		Period:            q.Period,
		TotalOrders:       25,
		TotalQuotes:       40,
		TotalRevenue:      15000.00,
		AverageOrderValue: 600.00,
		ConversionRate:    0.625,
	}, nil
}

// Event Handlers

// OrderCreatedEventHandler 訂單創建事件處理器
type OrderCreatedEventHandler struct{}

func (h *OrderCreatedEventHandler) Handle(ctx context.Context, event cqrs.Event) error {
	// 開始追蹤
	ctx, span := tracing.StartSpan(ctx, "HandleOrderCreatedEvent",
		tracing.WithAttributes(map[string]interface{}{
			"event.id":   event.GetID(),
			"event.type": event.GetEventType(),
		}),
	)
	defer span.End()
	
	log.Printf("[EventHandler] Order created: AggregateID=%s", event.GetAggregateID())
	
	// 這裡可以執行副作用，如發送郵件、更新讀模型等
	
	return nil
}

// InventoryUpdatedEventHandler 庫存更新事件處理器
type InventoryUpdatedEventHandler struct{}

func (h *InventoryUpdatedEventHandler) Handle(ctx context.Context, event cqrs.Event) error {
	// 開始追蹤
	ctx, span := tracing.StartSpan(ctx, "HandleInventoryUpdatedEvent",
		tracing.WithAttributes(map[string]interface{}{
			"event.id":   event.GetID(),
			"event.type": event.GetEventType(),
		}),
	)
	defer span.End()
	
	log.Printf("[EventHandler] Inventory updated: AggregateID=%s", event.GetAggregateID())
	
	// 檢查是否需要重新訂購
	// 更新庫存報表等
	
	return nil
}