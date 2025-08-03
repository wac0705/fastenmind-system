package messaging

import (
	"context"
	"time"
)

// MessageType 定義訊息類型
type MessageType string

const (
	// 訂單相關事件
	OrderCreated   MessageType = "order.created"
	OrderUpdated   MessageType = "order.updated"
	OrderCancelled MessageType = "order.cancelled"
	OrderCompleted MessageType = "order.completed"

	// 詢價相關事件
	InquiryCreated  MessageType = "inquiry.created"
	InquiryAssigned MessageType = "inquiry.assigned"
	InquiryQuoted   MessageType = "inquiry.quoted"

	// 庫存相關事件
	InventoryLow       MessageType = "inventory.low"
	InventoryRestocked MessageType = "inventory.restocked"

	// 客戶相關事件
	CustomerCreated MessageType = "customer.created"
	CustomerUpdated MessageType = "customer.updated"
	CreditChanged   MessageType = "customer.credit_changed"

	// 系統事件
	SystemAlert   MessageType = "system.alert"
	SystemMetrics MessageType = "system.metrics"
)

// Message 訊息介面
type Message interface {
	GetID() string
	GetType() MessageType
	GetTimestamp() time.Time
	GetPayload() interface{}
	GetHeaders() map[string]string
}

// BaseMessage 基礎訊息結構
type BaseMessage struct {
	ID        string              `json:"id"`
	Type      MessageType         `json:"type"`
	Timestamp time.Time           `json:"timestamp"`
	Payload   interface{}         `json:"payload"`
	Headers   map[string]string   `json:"headers"`
}

func (m *BaseMessage) GetID() string              { return m.ID }
func (m *BaseMessage) GetType() MessageType       { return m.Type }
func (m *BaseMessage) GetTimestamp() time.Time    { return m.Timestamp }
func (m *BaseMessage) GetPayload() interface{}    { return m.Payload }
func (m *BaseMessage) GetHeaders() map[string]string { return m.Headers }

// Publisher 發布者介面
type Publisher interface {
	Publish(ctx context.Context, topic string, message Message) error
	PublishBatch(ctx context.Context, topic string, messages []Message) error
}

// Subscriber 訂閱者介面
type Subscriber interface {
	Subscribe(ctx context.Context, topic string, handler MessageHandler) error
	Unsubscribe(ctx context.Context, topic string) error
}

// MessageHandler 訊息處理函數
type MessageHandler func(ctx context.Context, message Message) error

// MessageBroker 訊息代理介面
type MessageBroker interface {
	Publisher
	Subscriber
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

// DeadLetterQueue 死信佇列介面
type DeadLetterQueue interface {
	Push(message Message, err error) error
	Pop() (Message, error)
	List(limit int) ([]Message, error)
	Retry(messageID string) error
}

// MessageStore 訊息存儲介面
type MessageStore interface {
	Save(message Message) error
	Get(id string) (Message, error)
	List(filter MessageFilter) ([]Message, error)
	Delete(id string) error
}

// MessageFilter 訊息過濾器
type MessageFilter struct {
	Type      MessageType
	StartTime time.Time
	EndTime   time.Time
	Limit     int
	Offset    int
}

// RetryPolicy 重試策略
type RetryPolicy struct {
	MaxRetries     int
	InitialDelay   time.Duration
	MaxDelay       time.Duration
	Multiplier     float64
	MaxElapsedTime time.Duration
}

// DefaultRetryPolicy 預設重試策略
var DefaultRetryPolicy = RetryPolicy{
	MaxRetries:     3,
	InitialDelay:   1 * time.Second,
	MaxDelay:       30 * time.Second,
	Multiplier:     2.0,
	MaxElapsedTime: 5 * time.Minute,
}