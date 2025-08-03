package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
	
	"github.com/fastenmind/fastener-api/internal/config"
	infra "github.com/fastenmind/fastener-api/internal/infrastructure/messaging"
	"github.com/fastenmind/fastener-api/pkg/messaging"
	"github.com/fastenmind/fastener-api/pkg/messaging/rabbitmq"
	"github.com/fastenmind/fastener-api/pkg/tracing"
)

func main() {
	// 載入配置
	cfg := config.New()
	
	// 初始化追蹤
	tracer, err := tracing.NewTracer(tracing.Config{
		ServiceName:    "messaging-example",
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
	
	// 創建 RabbitMQ 配置
	rabbitConfig := rabbitmq.Config{
		URL:          cfg.Messaging.RabbitMQ.URL,
		Exchange:     cfg.Messaging.RabbitMQ.Exchange,
		ExchangeType: cfg.Messaging.RabbitMQ.ExchangeType,
		Durable:      cfg.Messaging.RabbitMQ.Durable,
		AutoDelete:   cfg.Messaging.RabbitMQ.AutoDelete,
		RetryPolicy: messaging.RetryPolicy{
			MaxRetries:     3,
			InitialDelay:   1 * time.Second,
			MaxDelay:       30 * time.Second,
			Multiplier:     2.0,
			MaxElapsedTime: 5 * time.Minute,
		},
		Logger: &SimpleLogger{},
	}
	
	// 創建訊息服務
	messagingService, err := infra.NewMessagingService(rabbitConfig)
	if err != nil {
		log.Fatal("Failed to create messaging service:", err)
	}
	
	// 啟動訊息服務
	ctx := context.Background()
	if err := messagingService.Start(ctx); err != nil {
		log.Fatal("Failed to start messaging service:", err)
	}
	defer messagingService.Stop(ctx)
	
	// 訂閱訂單事件
	orderHandler := &OrderEventHandler{}
	if err := messagingService.SubscribeToOrderEvents(ctx, orderHandler.Handle); err != nil {
		log.Fatal("Failed to subscribe to order events:", err)
	}
	
	// 訂閱庫存警報
	inventoryHandler := &InventoryAlertHandler{}
	if err := messagingService.SubscribeToInventoryAlerts(ctx, inventoryHandler.Handle); err != nil {
		log.Fatal("Failed to subscribe to inventory alerts:", err)
	}
	
	// 發布測試事件
	go publishTestEvents(ctx, messagingService)
	
	// 等待中斷信號
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	<-sigCh
	
	fmt.Println("Shutting down...")
}

// publishTestEvents 發布測試事件
func publishTestEvents(ctx context.Context, service *infra.MessagingService) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	
	orderCount := 0
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			orderCount++
			
			// 發布訂單創建事件
			orderID := fmt.Sprintf("ORDER-%d", orderCount)
			customerID := fmt.Sprintf("CUST-%d", orderCount%10)
			amount := float64(100 + orderCount*50)
			
			if err := service.PublishOrderCreated(ctx, orderID, customerID, amount); err != nil {
				log.Printf("Failed to publish order created event: %v", err)
			} else {
				log.Printf("Published order created event: %s", orderID)
			}
			
			// 每 3 個訂單發布一個庫存警報
			if orderCount%3 == 0 {
				productID := fmt.Sprintf("PROD-%d", orderCount%5)
				warehouseID := "WH-001"
				currentStock := float64(orderCount % 10)
				minStock := 20.0
				
				if err := service.PublishInventoryLow(ctx, productID, warehouseID, currentStock, minStock); err != nil {
					log.Printf("Failed to publish inventory low event: %v", err)
				} else {
					log.Printf("Published inventory low event: %s", productID)
				}
			}
		}
	}
}

// OrderEventHandler 訂單事件處理器
type OrderEventHandler struct{}

func (h *OrderEventHandler) Handle(ctx context.Context, message messaging.Message) error {
	// 開始追蹤
	ctx, span := tracing.StartSpan(ctx, "HandleOrderEvent",
		tracing.WithAttributes(map[string]interface{}{
			"message.id":   message.GetID(),
			"message.type": string(message.GetType()),
		}),
	)
	defer span.End()
	
	payload := message.GetPayload().(map[string]interface{})
	
	log.Printf("[OrderHandler] Received event: Type=%s, OrderID=%s, CustomerID=%s, Amount=%.2f",
		message.GetType(),
		payload["order_id"],
		payload["customer_id"],
		payload["total_amount"],
	)
	
	// 模擬處理時間
	time.Sleep(100 * time.Millisecond)
	
	// 記錄事件到追蹤
	tracing.AddEvent(ctx, "order_processed", map[string]interface{}{
		"order_id": payload["order_id"],
		"status":   "success",
	})
	
	return nil
}

// InventoryAlertHandler 庫存警報處理器
type InventoryAlertHandler struct{}

func (h *InventoryAlertHandler) Handle(ctx context.Context, message messaging.Message) error {
	// 開始追蹤
	ctx, span := tracing.StartSpan(ctx, "HandleInventoryAlert",
		tracing.WithAttributes(map[string]interface{}{
			"message.id":   message.GetID(),
			"message.type": string(message.GetType()),
		}),
	)
	defer span.End()
	
	payload := message.GetPayload().(map[string]interface{})
	
	log.Printf("[InventoryHandler] ALERT: Product=%s, Warehouse=%s, Current=%.2f, Min=%.2f",
		payload["product_id"],
		payload["warehouse_id"],
		payload["current_stock"],
		payload["min_stock"],
	)
	
	// 模擬發送通知
	time.Sleep(200 * time.Millisecond)
	
	// 記錄警報處理
	tracing.AddEvent(ctx, "alert_processed", map[string]interface{}{
		"product_id": payload["product_id"],
		"action":     "notification_sent",
	})
	
	return nil
}

// SimpleLogger 簡單日誌實現
type SimpleLogger struct{}

func (l *SimpleLogger) Info(msg string, fields ...interface{}) {
	log.Printf("[INFO] %s %v", msg, fields)
}

func (l *SimpleLogger) Error(msg string, err error, fields ...interface{}) {
	log.Printf("[ERROR] %s: %v %v", msg, err, fields)
}

func (l *SimpleLogger) Debug(msg string, fields ...interface{}) {
	log.Printf("[DEBUG] %s %v", msg, fields)
}