package eventbus

import (
	"fmt"
	"sync"

	"github.com/fastenmind/fastener-api/internal/domain/events"
)

// MemoryEventBus implements EventBus using in-memory pub/sub
type MemoryEventBus struct {
	handlers map[events.EventType][]events.EventHandler
	mu       sync.RWMutex
}

// NewMemoryEventBus creates a new in-memory event bus
func NewMemoryEventBus() *MemoryEventBus {
	return &MemoryEventBus{
		handlers: make(map[events.EventType][]events.EventHandler),
	}
}

// Subscribe registers a handler for specific event types
func (b *MemoryEventBus) Subscribe(eventType events.EventType, handler events.EventHandler) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	if handler == nil {
		return fmt.Errorf("handler cannot be nil")
	}
	
	b.handlers[eventType] = append(b.handlers[eventType], handler)
	return nil
}

// Unsubscribe removes a handler for specific event types
func (b *MemoryEventBus) Unsubscribe(eventType events.EventType, handler events.EventHandler) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	handlers, exists := b.handlers[eventType]
	if !exists {
		return nil
	}
	
	// Find and remove the handler
	for i, h := range handlers {
		if h == handler {
			b.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}
	
	// Clean up empty handler lists
	if len(b.handlers[eventType]) == 0 {
		delete(b.handlers, eventType)
	}
	
	return nil
}

// Publish publishes an event to all subscribed handlers
func (b *MemoryEventBus) Publish(event events.Event) error {
	b.mu.RLock()
	handlers := b.handlers[event.GetType()]
	b.mu.RUnlock()
	
	// Execute handlers concurrently
	var wg sync.WaitGroup
	errors := make(chan error, len(handlers))
	
	for _, handler := range handlers {
		wg.Add(1)
		go func(h events.EventHandler) {
			defer wg.Done()
			
			if h.CanHandle(event.GetType()) {
				if err := h.Handle(event); err != nil {
					errors <- fmt.Errorf("handler error for event %s: %w", event.GetType(), err)
				}
			}
		}(handler)
	}
	
	// Wait for all handlers to complete
	wg.Wait()
	close(errors)
	
	// Collect any errors
	var errs []error
	for err := range errors {
		if err != nil {
			errs = append(errs, err)
		}
	}
	
	if len(errs) > 0 {
		return fmt.Errorf("event bus errors: %v", errs)
	}
	
	return nil
}

// AsyncEventBus implements EventBus with async message processing
type AsyncEventBus struct {
	*MemoryEventBus
	eventQueue chan events.Event
	workers    int
	wg         sync.WaitGroup
	quit       chan struct{}
}

// NewAsyncEventBus creates a new async event bus
func NewAsyncEventBus(queueSize, workers int) *AsyncEventBus {
	bus := &AsyncEventBus{
		MemoryEventBus: NewMemoryEventBus(),
		eventQueue:     make(chan events.Event, queueSize),
		workers:        workers,
		quit:           make(chan struct{}),
	}
	
	bus.start()
	return bus
}

// start starts the worker goroutines
func (b *AsyncEventBus) start() {
	for i := 0; i < b.workers; i++ {
		b.wg.Add(1)
		go b.worker()
	}
}

// worker processes events from the queue
func (b *AsyncEventBus) worker() {
	defer b.wg.Done()
	
	for {
		select {
		case event := <-b.eventQueue:
			b.MemoryEventBus.Publish(event)
		case <-b.quit:
			return
		}
	}
}

// Publish queues an event for async processing
func (b *AsyncEventBus) Publish(event events.Event) error {
	select {
	case b.eventQueue <- event:
		return nil
	default:
		return fmt.Errorf("event queue is full")
	}
}

// Stop gracefully shuts down the async event bus
func (b *AsyncEventBus) Stop() {
	close(b.quit)
	b.wg.Wait()
	close(b.eventQueue)
}