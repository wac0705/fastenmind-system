package events

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// EventType represents the type of domain event
type EventType string

// Event interface that all domain events must implement
type Event interface {
	GetID() uuid.UUID
	GetType() EventType
	GetAggregateID() uuid.UUID
	GetAggregateType() string
	GetTimestamp() time.Time
	GetVersion() int
	GetData() interface{}
}

// BaseEvent contains common fields for all events
type BaseEvent struct {
	ID            uuid.UUID   `json:"id"`
	Type          EventType   `json:"type"`
	AggregateID   uuid.UUID   `json:"aggregate_id"`
	AggregateType string      `json:"aggregate_type"`
	Timestamp     time.Time   `json:"timestamp"`
	Version       int         `json:"version"`
	UserID        *uuid.UUID  `json:"user_id,omitempty"`
	CorrelationID string      `json:"correlation_id,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// NewBaseEvent creates a new base event
func NewBaseEvent(eventType EventType, aggregateID uuid.UUID, aggregateType string) BaseEvent {
	return BaseEvent{
		ID:            uuid.New(),
		Type:          eventType,
		AggregateID:   aggregateID,
		AggregateType: aggregateType,
		Timestamp:     time.Now().UTC(),
		Version:       1,
		Metadata:      make(map[string]interface{}),
	}
}

// GetID returns the event ID
func (e BaseEvent) GetID() uuid.UUID {
	return e.ID
}

// GetType returns the event type
func (e BaseEvent) GetType() EventType {
	return e.Type
}

// GetAggregateID returns the aggregate ID
func (e BaseEvent) GetAggregateID() uuid.UUID {
	return e.AggregateID
}

// GetAggregateType returns the aggregate type
func (e BaseEvent) GetAggregateType() string {
	return e.AggregateType
}

// GetTimestamp returns the event timestamp
func (e BaseEvent) GetTimestamp() time.Time {
	return e.Timestamp
}

// GetVersion returns the event version
func (e BaseEvent) GetVersion() int {
	return e.Version
}

// EventStore interface for persisting and retrieving events
type EventStore interface {
	// Save saves a new event to the store
	Save(event Event) error
	
	// SaveBatch saves multiple events in a single transaction
	SaveBatch(events []Event) error
	
	// GetByAggregateID retrieves all events for a specific aggregate
	GetByAggregateID(aggregateID uuid.UUID) ([]Event, error)
	
	// GetByAggregateIDAfterVersion retrieves events after a specific version
	GetByAggregateIDAfterVersion(aggregateID uuid.UUID, version int) ([]Event, error)
	
	// GetByEventType retrieves events of a specific type
	GetByEventType(eventType EventType, limit int) ([]Event, error)
	
	// GetSnapshot retrieves the latest snapshot for an aggregate
	GetSnapshot(aggregateID uuid.UUID) (*EventSnapshot, error)
	
	// SaveSnapshot saves a snapshot of an aggregate
	SaveSnapshot(snapshot *EventSnapshot) error
}

// EventSnapshot represents a snapshot of an aggregate at a point in time
type EventSnapshot struct {
	ID            uuid.UUID       `json:"id"`
	AggregateID   uuid.UUID       `json:"aggregate_id"`
	AggregateType string          `json:"aggregate_type"`
	Version       int             `json:"version"`
	Data          json.RawMessage `json:"data"`
	CreatedAt     time.Time       `json:"created_at"`
}

// EventPublisher interface for publishing events
type EventPublisher interface {
	// Publish publishes a single event
	Publish(event Event) error
	
	// PublishBatch publishes multiple events
	PublishBatch(events []Event) error
}

// EventHandler interface for handling events
type EventHandler interface {
	// Handle processes an event
	Handle(event Event) error
	
	// CanHandle returns true if the handler can process the event type
	CanHandle(eventType EventType) bool
}

// EventBus interface for event distribution
type EventBus interface {
	// Subscribe registers a handler for specific event types
	Subscribe(eventType EventType, handler EventHandler) error
	
	// Unsubscribe removes a handler for specific event types
	Unsubscribe(eventType EventType, handler EventHandler) error
	
	// Publish publishes an event to all subscribed handlers
	Publish(event Event) error
}