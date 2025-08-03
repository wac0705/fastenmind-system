package cqrs

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
	
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Event 事件介面
type Event interface {
	GetID() string
	GetAggregateID() string
	GetAggregateType() string
	GetEventType() string
	GetEventVersion() int
	GetTimestamp() time.Time
	GetData() interface{}
	GetMetadata() map[string]interface{}
}

// BaseEvent 基礎事件
type BaseEvent struct {
	ID            string                 `json:"id"`
	AggregateID   string                 `json:"aggregate_id"`
	AggregateType string                 `json:"aggregate_type"`
	EventType     string                 `json:"event_type"`
	EventVersion  int                    `json:"event_version"`
	Timestamp     time.Time              `json:"timestamp"`
	Data          interface{}            `json:"data"`
	Metadata      map[string]interface{} `json:"metadata"`
}

func (e *BaseEvent) GetID() string              { return e.ID }
func (e *BaseEvent) GetAggregateID() string     { return e.AggregateID }
func (e *BaseEvent) GetAggregateType() string   { return e.AggregateType }
func (e *BaseEvent) GetEventType() string       { return e.EventType }
func (e *BaseEvent) GetEventVersion() int       { return e.EventVersion }
func (e *BaseEvent) GetTimestamp() time.Time    { return e.Timestamp }
func (e *BaseEvent) GetData() interface{}       { return e.Data }
func (e *BaseEvent) GetMetadata() map[string]interface{} { return e.Metadata }

// EventStore 事件存儲介面
type EventStore interface {
	Save(ctx context.Context, events []Event) error
	GetEvents(ctx context.Context, aggregateID string, fromVersion int) ([]Event, error)
	GetEventsByType(ctx context.Context, eventType string, limit int) ([]Event, error)
	GetEventStream(ctx context.Context, streamName string, fromPosition int) ([]Event, error)
}

// EventPublisher 事件發布者介面
type EventPublisher interface {
	Publish(ctx context.Context, events []Event) error
}

// EventHandler 事件處理器介面
type EventHandler interface {
	Handle(ctx context.Context, event Event) error
}

// EventBus 事件匯流排介面
type EventBus interface {
	Subscribe(eventType string, handler EventHandler) error
	Publish(ctx context.Context, event Event) error
}

// EventSourcingRepository 事件溯源儲存庫介面
type EventSourcingRepository interface {
	Save(ctx context.Context, aggregate AggregateRoot) error
	GetByID(ctx context.Context, id string) (AggregateRoot, error)
}

// AggregateRoot 聚合根介面
type AggregateRoot interface {
	GetID() string
	GetVersion() int
	GetUncommittedEvents() []Event
	MarkEventsAsCommitted()
	LoadFromHistory(events []Event)
}

// SQLEventStore SQL事件存儲實現
type SQLEventStore struct {
	db *gorm.DB
}

// EventRecord 事件記錄
type EventRecord struct {
	ID             string    `gorm:"primaryKey"`
	AggregateID    string    `gorm:"index"`
	AggregateType  string    `gorm:"index"`
	EventType      string    `gorm:"index"`
	EventVersion   int       `gorm:"index"`
	EventData      string    `gorm:"type:text"`
	EventMetadata  string    `gorm:"type:text"`
	CreatedAt      time.Time `gorm:"index"`
}

func (EventRecord) TableName() string {
	return "event_store"
}

// NewSQLEventStore 創建SQL事件存儲
func NewSQLEventStore(db *gorm.DB) *SQLEventStore {
	return &SQLEventStore{db: db}
}

// Save 保存事件
func (s *SQLEventStore) Save(ctx context.Context, events []Event) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		for _, event := range events {
			// 序列化事件數據
			dataJSON, err := json.Marshal(event.GetData())
			if err != nil {
				return fmt.Errorf("failed to marshal event data: %w", err)
			}
			
			// 序列化元數據
			metadataJSON, err := json.Marshal(event.GetMetadata())
			if err != nil {
				return fmt.Errorf("failed to marshal event metadata: %w", err)
			}
			
			record := EventRecord{
				ID:            event.GetID(),
				AggregateID:   event.GetAggregateID(),
				AggregateType: event.GetAggregateType(),
				EventType:     event.GetEventType(),
				EventVersion:  event.GetEventVersion(),
				EventData:     string(dataJSON),
				EventMetadata: string(metadataJSON),
				CreatedAt:     event.GetTimestamp(),
			}
			
			if err := tx.Create(&record).Error; err != nil {
				return fmt.Errorf("failed to save event: %w", err)
			}
		}
		return nil
	})
}

// GetEvents 獲取聚合的事件
func (s *SQLEventStore) GetEvents(ctx context.Context, aggregateID string, fromVersion int) ([]Event, error) {
	var records []EventRecord
	
	err := s.db.Where("aggregate_id = ? AND event_version >= ?", aggregateID, fromVersion).
		Order("event_version ASC").
		Find(&records).Error
		
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}
	
	events := make([]Event, len(records))
	for i, record := range records {
		event, err := s.recordToEvent(record)
		if err != nil {
			return nil, err
		}
		events[i] = event
	}
	
	return events, nil
}

// GetEventsByType 根據類型獲取事件
func (s *SQLEventStore) GetEventsByType(ctx context.Context, eventType string, limit int) ([]Event, error) {
	var records []EventRecord
	
	query := s.db.Where("event_type = ?", eventType).Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	
	err := query.Find(&records).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get events by type: %w", err)
	}
	
	events := make([]Event, len(records))
	for i, record := range records {
		event, err := s.recordToEvent(record)
		if err != nil {
			return nil, err
		}
		events[i] = event
	}
	
	return events, nil
}

// GetEventStream 獲取事件流
func (s *SQLEventStore) GetEventStream(ctx context.Context, streamName string, fromPosition int) ([]Event, error) {
	var records []EventRecord
	
	err := s.db.Where("aggregate_type = ?", streamName).
		Offset(fromPosition).
		Order("created_at ASC").
		Find(&records).Error
		
	if err != nil {
		return nil, fmt.Errorf("failed to get event stream: %w", err)
	}
	
	events := make([]Event, len(records))
	for i, record := range records {
		event, err := s.recordToEvent(record)
		if err != nil {
			return nil, err
		}
		events[i] = event
	}
	
	return events, nil
}

// recordToEvent 將記錄轉換為事件
func (s *SQLEventStore) recordToEvent(record EventRecord) (Event, error) {
	var data interface{}
	if err := json.Unmarshal([]byte(record.EventData), &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal event data: %w", err)
	}
	
	var metadata map[string]interface{}
	if err := json.Unmarshal([]byte(record.EventMetadata), &metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal event metadata: %w", err)
	}
	
	return &BaseEvent{
		ID:            record.ID,
		AggregateID:   record.AggregateID,
		AggregateType: record.AggregateType,
		EventType:     record.EventType,
		EventVersion:  record.EventVersion,
		Timestamp:     record.CreatedAt,
		Data:          data,
		Metadata:      metadata,
	}, nil
}

// SimpleEventBus 簡單事件匯流排實現
type SimpleEventBus struct {
	handlers map[string][]EventHandler
	mu       sync.RWMutex
}

// NewSimpleEventBus 創建簡單事件匯流排
func NewSimpleEventBus() *SimpleEventBus {
	return &SimpleEventBus{
		handlers: make(map[string][]EventHandler),
	}
}

// Subscribe 訂閱事件
func (b *SimpleEventBus) Subscribe(eventType string, handler EventHandler) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	b.handlers[eventType] = append(b.handlers[eventType], handler)
	return nil
}

// Publish 發布事件
func (b *SimpleEventBus) Publish(ctx context.Context, event Event) error {
	b.mu.RLock()
	handlers := b.handlers[event.GetEventType()]
	b.mu.RUnlock()
	
	for _, handler := range handlers {
		// 異步處理事件
		go func(h EventHandler) {
			if err := h.Handle(ctx, event); err != nil {
				// 記錄錯誤
				fmt.Printf("Error handling event %s: %v\n", event.GetID(), err)
			}
		}(handler)
	}
	
	return nil
}

// Domain Events

// OrderCreatedEvent 訂單創建事件
type OrderCreatedEvent struct {
	BaseEvent
	OrderID     string    `json:"order_id"`
	CustomerID  string    `json:"customer_id"`
	TotalAmount float64   `json:"total_amount"`
	Currency    string    `json:"currency"`
	Items       []OrderItem `json:"items"`
}

func NewOrderCreatedEvent(orderID, customerID string, totalAmount float64, currency string) *OrderCreatedEvent {
	return &OrderCreatedEvent{
		BaseEvent: BaseEvent{
			ID:            uuid.New().String(),
			AggregateID:   orderID,
			AggregateType: "Order",
			EventType:     "OrderCreated",
			EventVersion:  1,
			Timestamp:     time.Now(),
			Metadata:      make(map[string]interface{}),
		},
		OrderID:     orderID,
		CustomerID:  customerID,
		TotalAmount: totalAmount,
		Currency:    currency,
	}
}

// InventoryUpdatedEvent 庫存更新事件
type InventoryUpdatedEvent struct {
	BaseEvent
	ProductID   string  `json:"product_id"`
	WarehouseID string  `json:"warehouse_id"`
	Quantity    float64 `json:"quantity"`
	Type        string  `json:"type"`
	Reference   string  `json:"reference"`
}

func NewInventoryUpdatedEvent(productID, warehouseID string, quantity float64, updateType string) *InventoryUpdatedEvent {
	return &InventoryUpdatedEvent{
		BaseEvent: BaseEvent{
			ID:            uuid.New().String(),
			AggregateID:   productID,
			AggregateType: "Inventory",
			EventType:     "InventoryUpdated",
			EventVersion:  1,
			Timestamp:     time.Now(),
			Metadata:      make(map[string]interface{}),
		},
		ProductID:   productID,
		WarehouseID: warehouseID,
		Quantity:    quantity,
		Type:        updateType,
	}
}

// CustomerCreditUpdatedEvent 客戶信用額度更新事件
type CustomerCreditUpdatedEvent struct {
	BaseEvent
	CustomerID   string    `json:"customer_id"`
	OldLimit     float64   `json:"old_limit"`
	NewLimit     float64   `json:"new_limit"`
	ApprovedBy   string    `json:"approved_by"`
	Reason       string    `json:"reason"`
}

func NewCustomerCreditUpdatedEvent(customerID string, oldLimit, newLimit float64, approvedBy string) *CustomerCreditUpdatedEvent {
	return &CustomerCreditUpdatedEvent{
		BaseEvent: BaseEvent{
			ID:            uuid.New().String(),
			AggregateID:   customerID,
			AggregateType: "Customer",
			EventType:     "CustomerCreditUpdated",
			EventVersion:  1,
			Timestamp:     time.Now(),
			Metadata:      make(map[string]interface{}),
		},
		CustomerID: customerID,
		OldLimit:   oldLimit,
		NewLimit:   newLimit,
		ApprovedBy: approvedBy,
	}
}