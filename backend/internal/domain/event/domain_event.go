package event

import (
	"time"

	"github.com/google/uuid"
)

// DomainEvent 領域事件接口
type DomainEvent interface {
	GetAggregateID() uuid.UUID
	GetEventType() string
	GetOccurredAt() time.Time
	GetEventID() uuid.UUID
	GetEventVersion() int
}

// BaseDomainEvent 基礎領域事件
type BaseDomainEvent struct {
	EventID       uuid.UUID `json:"event_id"`
	AggregateID   uuid.UUID `json:"aggregate_id"`
	EventType     string    `json:"event_type"`
	OccurredAt    time.Time `json:"occurred_at"`
	EventVersion  int       `json:"event_version"`
}

// GetEventID 獲取事件ID
func (e BaseDomainEvent) GetEventID() uuid.UUID {
	return e.EventID
}

// GetAggregateID 獲取聚合根ID
func (e BaseDomainEvent) GetAggregateID() uuid.UUID {
	return e.AggregateID
}

// GetEventType 獲取事件類型
func (e BaseDomainEvent) GetEventType() string {
	return e.EventType
}

// GetOccurredAt 獲取發生時間
func (e BaseDomainEvent) GetOccurredAt() time.Time {
	return e.OccurredAt
}

// GetEventVersion 獲取事件版本
func (e BaseDomainEvent) GetEventVersion() int {
	return e.EventVersion
}

// NewBaseDomainEvent 創建基礎領域事件
func NewBaseDomainEvent(aggregateID uuid.UUID, eventType string) BaseDomainEvent {
	return BaseDomainEvent{
		EventID:      uuid.New(),
		AggregateID:  aggregateID,
		EventType:    eventType,
		OccurredAt:   time.Now(),
		EventVersion: 1,
	}
}