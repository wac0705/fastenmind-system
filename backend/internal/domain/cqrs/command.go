package cqrs

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Command represents a write operation
type Command interface {
	GetID() uuid.UUID
	GetName() string
	GetTimestamp() time.Time
	Validate() error
}

// BaseCommand contains common fields for all commands
type BaseCommand struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
	UserID    uuid.UUID `json:"user_id"`
}

// NewBaseCommand creates a new base command
func NewBaseCommand(name string, userID uuid.UUID) BaseCommand {
	return BaseCommand{
		ID:        uuid.New(),
		Name:      name,
		Timestamp: time.Now().UTC(),
		UserID:    userID,
	}
}

// GetID returns the command ID
func (c BaseCommand) GetID() uuid.UUID {
	return c.ID
}

// GetName returns the command name
func (c BaseCommand) GetName() string {
	return c.Name
}

// GetTimestamp returns the command timestamp
func (c BaseCommand) GetTimestamp() time.Time {
	return c.Timestamp
}

// CommandHandler processes commands
type CommandHandler interface {
	Handle(ctx context.Context, command Command) error
}

// CommandBus dispatches commands to handlers
type CommandBus interface {
	// Register registers a handler for a command type
	Register(commandName string, handler CommandHandler) error
	
	// Send sends a command to its handler
	Send(ctx context.Context, command Command) error
	
	// SendAsync sends a command asynchronously
	SendAsync(ctx context.Context, command Command) error
}

// Query represents a read operation
type Query interface {
	GetName() string
	Validate() error
}

// QueryResult represents the result of a query
type QueryResult interface{}

// QueryHandler processes queries
type QueryHandler interface {
	Handle(ctx context.Context, query Query) (QueryResult, error)
}

// QueryBus dispatches queries to handlers
type QueryBus interface {
	// Register registers a handler for a query type
	Register(queryName string, handler QueryHandler) error
	
	// Send sends a query to its handler
	Send(ctx context.Context, query Query) (QueryResult, error)
}