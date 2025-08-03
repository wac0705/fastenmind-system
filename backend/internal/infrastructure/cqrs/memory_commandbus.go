package cqrs

import (
	"context"
	"fmt"
	"sync"

	"github.com/fastenmind/fastener-api/internal/domain/cqrs"
)

// MemoryCommandBus implements CommandBus using in-memory dispatch
type MemoryCommandBus struct {
	handlers map[string]cqrs.CommandHandler
	mu       sync.RWMutex
}

// NewMemoryCommandBus creates a new in-memory command bus
func NewMemoryCommandBus() *MemoryCommandBus {
	return &MemoryCommandBus{
		handlers: make(map[string]cqrs.CommandHandler),
	}
}

// Register registers a handler for a command type
func (b *MemoryCommandBus) Register(commandName string, handler cqrs.CommandHandler) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	if handler == nil {
		return fmt.Errorf("handler cannot be nil")
	}
	
	if _, exists := b.handlers[commandName]; exists {
		return fmt.Errorf("handler already registered for command: %s", commandName)
	}
	
	b.handlers[commandName] = handler
	return nil
}

// Send sends a command to its handler
func (b *MemoryCommandBus) Send(ctx context.Context, command cqrs.Command) error {
	if command == nil {
		return fmt.Errorf("command cannot be nil")
	}
	
	// Validate command
	if err := command.Validate(); err != nil {
		return fmt.Errorf("command validation failed: %w", err)
	}
	
	b.mu.RLock()
	handler, exists := b.handlers[command.GetName()]
	b.mu.RUnlock()
	
	if !exists {
		return fmt.Errorf("no handler registered for command: %s", command.GetName())
	}
	
	// Execute handler
	return handler.Handle(ctx, command)
}

// SendAsync sends a command asynchronously
func (b *MemoryCommandBus) SendAsync(ctx context.Context, command cqrs.Command) error {
	if command == nil {
		return fmt.Errorf("command cannot be nil")
	}
	
	// Validate command
	if err := command.Validate(); err != nil {
		return fmt.Errorf("command validation failed: %w", err)
	}
	
	b.mu.RLock()
	handler, exists := b.handlers[command.GetName()]
	b.mu.RUnlock()
	
	if !exists {
		return fmt.Errorf("no handler registered for command: %s", command.GetName())
	}
	
	// Execute handler asynchronously
	go func() {
		if err := handler.Handle(ctx, command); err != nil {
			// In production, this should be logged or sent to an error handler
			fmt.Printf("async command execution failed: %v\n", err)
		}
	}()
	
	return nil
}

// MemoryQueryBus implements QueryBus using in-memory dispatch
type MemoryQueryBus struct {
	handlers map[string]cqrs.QueryHandler
	mu       sync.RWMutex
}

// NewMemoryQueryBus creates a new in-memory query bus
func NewMemoryQueryBus() *MemoryQueryBus {
	return &MemoryQueryBus{
		handlers: make(map[string]cqrs.QueryHandler),
	}
}

// Register registers a handler for a query type
func (b *MemoryQueryBus) Register(queryName string, handler cqrs.QueryHandler) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	if handler == nil {
		return fmt.Errorf("handler cannot be nil")
	}
	
	if _, exists := b.handlers[queryName]; exists {
		return fmt.Errorf("handler already registered for query: %s", queryName)
	}
	
	b.handlers[queryName] = handler
	return nil
}

// Send sends a query to its handler
func (b *MemoryQueryBus) Send(ctx context.Context, query cqrs.Query) (cqrs.QueryResult, error) {
	if query == nil {
		return nil, fmt.Errorf("query cannot be nil")
	}
	
	// Validate query
	if err := query.Validate(); err != nil {
		return nil, fmt.Errorf("query validation failed: %w", err)
	}
	
	b.mu.RLock()
	handler, exists := b.handlers[query.GetName()]
	b.mu.RUnlock()
	
	if !exists {
		return nil, fmt.Errorf("no handler registered for query: %s", query.GetName())
	}
	
	// Execute handler
	return handler.Handle(ctx, query)
}