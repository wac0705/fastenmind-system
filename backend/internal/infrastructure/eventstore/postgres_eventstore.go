package eventstore

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/fastenmind/fastener-api/internal/domain/events"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// EventRecord represents how events are stored in the database
type EventRecord struct {
	ID            uuid.UUID       `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	EventType     string          `gorm:"type:varchar(100);not null;index"`
	AggregateID   uuid.UUID       `gorm:"type:uuid;not null;index"`
	AggregateType string          `gorm:"type:varchar(50);not null;index"`
	EventVersion  int             `gorm:"not null"`
	EventData     json.RawMessage `gorm:"type:jsonb;not null"`
	Metadata      json.RawMessage `gorm:"type:jsonb"`
	UserID        *uuid.UUID      `gorm:"type:uuid"`
	CreatedAt     time.Time       `gorm:"not null;index"`
}

func (EventRecord) TableName() string {
	return "event_store"
}

// SnapshotRecord represents how snapshots are stored
type SnapshotRecord struct {
	ID            uuid.UUID       `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	AggregateID   uuid.UUID       `gorm:"type:uuid;not null;uniqueIndex"`
	AggregateType string          `gorm:"type:varchar(50);not null"`
	Version       int             `gorm:"not null"`
	Data          json.RawMessage `gorm:"type:jsonb;not null"`
	CreatedAt     time.Time       `gorm:"not null"`
}

func (SnapshotRecord) TableName() string {
	return "event_snapshots"
}

// PostgresEventStore implements EventStore using PostgreSQL
type PostgresEventStore struct {
	db *gorm.DB
}

// NewPostgresEventStore creates a new PostgreSQL event store
func NewPostgresEventStore(db *gorm.DB) *PostgresEventStore {
	// Auto migrate tables
	db.AutoMigrate(&EventRecord{}, &SnapshotRecord{})
	
	return &PostgresEventStore{
		db: db,
	}
}

// Save saves a new event to the store
func (s *PostgresEventStore) Save(event events.Event) error {
	eventData, err := json.Marshal(event.GetData())
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}
	
	metadata := make(map[string]interface{})
	if baseEvent, ok := event.(events.BaseEvent); ok {
		metadata = baseEvent.Metadata
	}
	
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}
	
	record := EventRecord{
		ID:            event.GetID(),
		EventType:     string(event.GetType()),
		AggregateID:   event.GetAggregateID(),
		AggregateType: event.GetAggregateType(),
		EventVersion:  event.GetVersion(),
		EventData:     eventData,
		Metadata:      metadataJSON,
		CreatedAt:     event.GetTimestamp(),
	}
	
	if baseEvent, ok := event.(events.BaseEvent); ok && baseEvent.UserID != nil {
		record.UserID = baseEvent.UserID
	}
	
	return s.db.Create(&record).Error
}

// SaveBatch saves multiple events in a single transaction
func (s *PostgresEventStore) SaveBatch(events []events.Event) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		for _, event := range events {
			eventData, err := json.Marshal(event.GetData())
			if err != nil {
				return fmt.Errorf("failed to marshal event data: %w", err)
			}
			
			metadata := make(map[string]interface{})
			if baseEvent, ok := event.(events.BaseEvent); ok {
				metadata = baseEvent.Metadata
			}
			
			metadataJSON, err := json.Marshal(metadata)
			if err != nil {
				return fmt.Errorf("failed to marshal metadata: %w", err)
			}
			
			record := EventRecord{
				ID:            event.GetID(),
				EventType:     string(event.GetType()),
				AggregateID:   event.GetAggregateID(),
				AggregateType: event.GetAggregateType(),
				EventVersion:  event.GetVersion(),
				EventData:     eventData,
				Metadata:      metadataJSON,
				CreatedAt:     event.GetTimestamp(),
			}
			
			if baseEvent, ok := event.(events.BaseEvent); ok && baseEvent.UserID != nil {
				record.UserID = baseEvent.UserID
			}
			
			if err := tx.Create(&record).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetByAggregateID retrieves all events for a specific aggregate
func (s *PostgresEventStore) GetByAggregateID(aggregateID uuid.UUID) ([]events.Event, error) {
	var records []EventRecord
	err := s.db.Where("aggregate_id = ?", aggregateID).
		Order("event_version ASC").
		Find(&records).Error
	if err != nil {
		return nil, err
	}
	
	return s.recordsToEvents(records)
}

// GetByAggregateIDAfterVersion retrieves events after a specific version
func (s *PostgresEventStore) GetByAggregateIDAfterVersion(aggregateID uuid.UUID, version int) ([]events.Event, error) {
	var records []EventRecord
	err := s.db.Where("aggregate_id = ? AND event_version > ?", aggregateID, version).
		Order("event_version ASC").
		Find(&records).Error
	if err != nil {
		return nil, err
	}
	
	return s.recordsToEvents(records)
}

// GetByEventType retrieves events of a specific type
func (s *PostgresEventStore) GetByEventType(eventType events.EventType, limit int) ([]events.Event, error) {
	var records []EventRecord
	query := s.db.Where("event_type = ?", string(eventType)).
		Order("created_at DESC")
		
	if limit > 0 {
		query = query.Limit(limit)
	}
	
	err := query.Find(&records).Error
	if err != nil {
		return nil, err
	}
	
	return s.recordsToEvents(records)
}

// GetSnapshot retrieves the latest snapshot for an aggregate
func (s *PostgresEventStore) GetSnapshot(aggregateID uuid.UUID) (*events.EventSnapshot, error) {
	var record SnapshotRecord
	err := s.db.Where("aggregate_id = ?", aggregateID).First(&record).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	
	return &events.EventSnapshot{
		ID:            record.ID,
		AggregateID:   record.AggregateID,
		AggregateType: record.AggregateType,
		Version:       record.Version,
		Data:          record.Data,
		CreatedAt:     record.CreatedAt,
	}, nil
}

// SaveSnapshot saves a snapshot of an aggregate
func (s *PostgresEventStore) SaveSnapshot(snapshot *events.EventSnapshot) error {
	record := SnapshotRecord{
		ID:            snapshot.ID,
		AggregateID:   snapshot.AggregateID,
		AggregateType: snapshot.AggregateType,
		Version:       snapshot.Version,
		Data:          snapshot.Data,
		CreatedAt:     snapshot.CreatedAt,
	}
	
	// Upsert - update if exists, create if not
	return s.db.Save(&record).Error
}

// recordsToEvents converts database records to domain events
func (s *PostgresEventStore) recordsToEvents(records []EventRecord) ([]events.Event, error) {
	result := make([]events.Event, 0, len(records))
	
	for _, record := range records {
		event, err := s.recordToEvent(record)
		if err != nil {
			return nil, fmt.Errorf("failed to convert record to event: %w", err)
		}
		result = append(result, event)
	}
	
	return result, nil
}

// recordToEvent converts a single database record to a domain event
func (s *PostgresEventStore) recordToEvent(record EventRecord) (events.Event, error) {
	// Create base event
	baseEvent := events.BaseEvent{
		ID:            record.ID,
		Type:          events.EventType(record.EventType),
		AggregateID:   record.AggregateID,
		AggregateType: record.AggregateType,
		Timestamp:     record.CreatedAt,
		Version:       record.EventVersion,
		UserID:        record.UserID,
	}
	
	// Unmarshal metadata
	if len(record.Metadata) > 0 {
		if err := json.Unmarshal(record.Metadata, &baseEvent.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
	}
	
	// Create specific event type based on EventType
	switch events.EventType(record.EventType) {
	case events.InquiryCreated:
		var event events.InquiryCreatedEvent
		event.BaseEvent = baseEvent
		if err := json.Unmarshal(record.EventData, &event); err != nil {
			return nil, err
		}
		return &event, nil
		
	case events.InquiryAssigned:
		var event events.InquiryAssignedEvent
		event.BaseEvent = baseEvent
		if err := json.Unmarshal(record.EventData, &event); err != nil {
			return nil, err
		}
		return &event, nil
		
	case events.QuoteCreated:
		var event events.QuoteCreatedEvent
		event.BaseEvent = baseEvent
		if err := json.Unmarshal(record.EventData, &event); err != nil {
			return nil, err
		}
		return &event, nil
		
	case events.OrderCreated:
		var event events.OrderCreatedEvent
		event.BaseEvent = baseEvent
		if err := json.Unmarshal(record.EventData, &event); err != nil {
			return nil, err
		}
		return &event, nil
		
	// Add more event types as needed
	default:
		return nil, fmt.Errorf("unknown event type: %s", record.EventType)
	}
}