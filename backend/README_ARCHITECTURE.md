# FastenMind 架構升級說明

## 概述

本文檔說明 FastenMind 系統的最新架構升級，包括訊息佇列、CQRS 模式和分散式追蹤的實作。

## 1. 訊息佇列 (Message Queue)

### 技術選型
- **RabbitMQ**: 作為主要的訊息代理，支援多種訊息模式
- **訊息模式**: Topic Exchange，支援靈活的路由

### 主要組件

#### 訊息類型定義
```go
// pkg/messaging/types.go
- OrderCreated / OrderUpdated / OrderCancelled
- InquiryCreated / InquiryAssigned / InquiryQuoted
- InventoryLow / InventoryRestocked
- CustomerCreated / CustomerUpdated / CreditChanged
```

#### RabbitMQ 實作
```go
// pkg/messaging/rabbitmq/broker.go
- 自動重連機制
- 重試策略 (指數退避)
- 死信佇列處理
- 批量發布支援
```

### 使用範例
```go
// 發布訂單創建事件
err := messagingService.PublishOrderCreated(ctx, orderID, customerID, totalAmount)

// 訂閱訂單事件
err := messagingService.SubscribeToOrderEvents(ctx, orderHandler.Handle)
```

## 2. CQRS 模式 (Command Query Responsibility Segregation)

### 架構設計
- **命令端**: 處理寫入操作，確保資料一致性
- **查詢端**: 優化讀取效能，支援複雜查詢
- **事件溯源**: 記錄所有領域事件，支援審計和重播

### 主要組件

#### 命令 (Commands)
```go
// pkg/cqrs/command.go
- CreateOrderCommand
- UpdateInventoryCommand
- AssignEngineerCommand
- UpdateCustomerCreditCommand
- CalculateCostCommand
```

#### 查詢 (Queries)
```go
// pkg/cqrs/query.go
- GetOrderByIDQuery
- ListOrdersQuery
- GetInventoryStatusQuery
- GetCustomerStatisticsQuery
- GetEngineerWorkloadQuery
```

#### 事件存儲
```go
// pkg/cqrs/event_store.go
- SQLEventStore: 使用 PostgreSQL 存儲事件
- 支援事件重播和快照
- 事件版本控制
```

### 使用範例
```go
// 執行命令
cmd := cqrs.NewCreateOrderCommand(customerID, quoteID, items)
err := cqrsService.ExecuteCommand(ctx, cmd)

// 執行查詢
query := cqrs.NewGetOrderByIDQuery(orderID)
result, err := cqrsService.ExecuteQuery(ctx, query)
```

## 3. 分散式追蹤 (Distributed Tracing)

### 技術選型
- **OpenTelemetry**: 標準化的追蹤框架
- **Jaeger**: 追蹤資料的收集和視覺化

### 主要功能

#### 追蹤器配置
```go
// pkg/tracing/tracer.go
- 支援 Jaeger 和 OTLP 導出器
- 可配置的採樣率
- 自動傳播追蹤上下文
```

#### 中間件整合
```go
// pkg/tracing/middleware.go
- Echo HTTP 中間件
- GORM 資料庫追蹤
- RabbitMQ 訊息追蹤
- HTTP 客戶端追蹤
```

### 使用範例
```go
// 開始新的 span
ctx, span := tracing.StartSpan(ctx, "ProcessOrder",
    tracing.WithAttributes(map[string]interface{}{
        "order.id": orderID,
        "customer.id": customerID,
    }),
)
defer span.End()

// 記錄錯誤
if err != nil {
    tracing.RecordError(ctx, err, "Failed to process order")
}
```

## 4. 基礎設施配置

### Docker Compose
```yaml
# docker-compose.infrastructure.yml
- PostgreSQL: 主資料庫
- Redis: 快取和 Session 存儲
- RabbitMQ: 訊息佇列
- Jaeger: 分散式追蹤
- Prometheus: 監控指標
- Grafana: 視覺化儀表板
- N8N: 工作流自動化
```

### 環境變數配置
```bash
# 訊息佇列
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
RABBITMQ_EXCHANGE=fastenmind
RABBITMQ_MAX_RETRIES=3

# 追蹤
TRACING_ENABLED=true
TRACING_SERVICE_NAME=fastenmind-api
TRACING_ENDPOINT=http://localhost:14268/api/traces
TRACING_SAMPLING_RATE=1.0

# CQRS
EVENT_STORE_TYPE=sql
EVENT_STORE_SNAPSHOT_FREQ=10
SNAPSHOT_ENABLED=true
```

## 5. 效能優化

### 訊息處理
- 並行消費者支援
- 批量處理優化
- 訊息壓縮選項

### 查詢優化
- 讀寫分離架構
- 查詢結果快取
- 投影預計算

### 追蹤優化
- 動態採樣率調整
- 批量導出減少開銷
- 異步處理避免阻塞

## 6. 監控和觀測性

### 關鍵指標
- 訊息吞吐量和延遲
- 命令/查詢執行時間
- 錯誤率和重試次數
- 系統資源使用率

### 儀表板
- Grafana 預設儀表板
- 業務指標視覺化
- 系統健康狀態監控

## 7. 開發指南

### 新增事件處理器
```go
type MyEventHandler struct{}

func (h *MyEventHandler) Handle(ctx context.Context, message messaging.Message) error {
    // 處理邏輯
    return nil
}

// 註冊處理器
messagingService.Subscribe(ctx, "my.events", handler)
```

### 新增 CQRS 命令
```go
// 1. 定義命令
type MyCommand struct {
    cqrs.BaseCommand
    // 欄位定義
}

// 2. 實作處理器
type MyCommandHandler struct{}

func (h *MyCommandHandler) Handle(ctx context.Context, cmd cqrs.Command) error {
    // 處理邏輯
    return nil
}

// 3. 註冊處理器
commandBus.Register("MyCommand", handler)
```

## 8. 測試和除錯

### 本地測試環境
```bash
# 啟動基礎設施
docker-compose -f docker-compose.infrastructure.yml up -d

# 執行訊息佇列範例
go run examples/messaging/main.go

# 執行 CQRS 範例
go run examples/cqrs/main.go
```

### 追蹤除錯
- Jaeger UI: http://localhost:16686
- RabbitMQ Management: http://localhost:15672
- Grafana: http://localhost:3000

## 9. 部署考量

### 高可用性
- RabbitMQ 叢集配置
- PostgreSQL 主從複製
- 服務多實例部署

### 安全性
- TLS 加密通訊
- 認證和授權機制
- 敏感資料加密

## 10. 未來規劃

### 短期目標
- 完善事件重播機制
- 增加更多業務事件
- 優化查詢效能

### 長期目標
- 支援 Kafka 作為訊息佇列
- 實作 Saga 模式
- 多租戶完全隔離