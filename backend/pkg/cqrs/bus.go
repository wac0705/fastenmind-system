package cqrs

import (
	"context"
	"fmt"
	"sync"
)

// SimplCommandBus 簡單命令匯流排實現
type SimpleCommandBus struct {
	handlers map[string]CommandHandler
	mu       sync.RWMutex
}

// NewSimpleCommandBus 創建簡單命令匯流排
func NewSimpleCommandBus() *SimpleCommandBus {
	return &SimpleCommandBus{
		handlers: make(map[string]CommandHandler),
	}
}

// Register 註冊命令處理器
func (b *SimpleCommandBus) Register(commandName string, handler CommandHandler) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	if _, exists := b.handlers[commandName]; exists {
		return fmt.Errorf("handler already registered for command: %s", commandName)
	}
	
	b.handlers[commandName] = handler
	return nil
}

// Send 同步發送命令
func (b *SimpleCommandBus) Send(ctx context.Context, command Command) error {
	b.mu.RLock()
	handler, exists := b.handlers[command.GetName()]
	b.mu.RUnlock()
	
	if !exists {
		return fmt.Errorf("no handler registered for command: %s", command.GetName())
	}
	
	// 驗證命令
	if validator, ok := command.(interface{ Validate() error }); ok {
		if err := validator.Validate(); err != nil {
			return fmt.Errorf("command validation failed: %w", err)
		}
	}
	
	return handler.Handle(ctx, command)
}

// SendAsync 異步發送命令
func (b *SimpleCommandBus) SendAsync(ctx context.Context, command Command) <-chan error {
	errCh := make(chan error, 1)
	
	go func() {
		defer close(errCh)
		errCh <- b.Send(ctx, command)
	}()
	
	return errCh
}

// SimpleQueryBus 簡單查詢匯流排實現
type SimpleQueryBus struct {
	handlers map[string]QueryHandler
	mu       sync.RWMutex
}

// NewSimpleQueryBus 創建簡單查詢匯流排
func NewSimpleQueryBus() *SimpleQueryBus {
	return &SimpleQueryBus{
		handlers: make(map[string]QueryHandler),
	}
}

// Register 註冊查詢處理器
func (b *SimpleQueryBus) Register(queryName string, handler QueryHandler) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	if _, exists := b.handlers[queryName]; exists {
		return fmt.Errorf("handler already registered for query: %s", queryName)
	}
	
	b.handlers[queryName] = handler
	return nil
}

// Send 發送查詢
func (b *SimpleQueryBus) Send(ctx context.Context, query Query) (QueryResult, error) {
	b.mu.RLock()
	handler, exists := b.handlers[query.GetName()]
	b.mu.RUnlock()
	
	if !exists {
		return nil, fmt.Errorf("no handler registered for query: %s", query.GetName())
	}
	
	return handler.Handle(ctx, query)
}

// MessageBusCommandBus 基於訊息匯流排的命令匯流排
type MessageBusCommandBus struct {
	localBus  *SimpleCommandBus
	messageBus interface{} // 可以是 RabbitMQ, Kafka 等
	topic      string
}

// NewMessageBusCommandBus 創建基於訊息匯流排的命令匯流排
func NewMessageBusCommandBus(messageBus interface{}, topic string) *MessageBusCommandBus {
	return &MessageBusCommandBus{
		localBus:   NewSimpleCommandBus(),
		messageBus: messageBus,
		topic:      topic,
	}
}

// Register 註冊命令處理器
func (b *MessageBusCommandBus) Register(commandName string, handler CommandHandler) error {
	return b.localBus.Register(commandName, handler)
}

// Send 發送命令到訊息匯流排
func (b *MessageBusCommandBus) Send(ctx context.Context, command Command) error {
	// 這裡可以實現將命令發送到訊息匯流排的邏輯
	// 例如：序列化命令並發布到 RabbitMQ
	return b.localBus.Send(ctx, command)
}

// SendAsync 異步發送命令
func (b *MessageBusCommandBus) SendAsync(ctx context.Context, command Command) <-chan error {
	return b.localBus.SendAsync(ctx, command)
}