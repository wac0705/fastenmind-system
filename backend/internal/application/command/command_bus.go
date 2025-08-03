package command

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"go.uber.org/zap"
)

// InMemoryCommandBus 內存命令總線實現
type InMemoryCommandBus struct {
	handlers    map[string]interface{}
	middlewares []Middleware
	logger      *zap.Logger
	mu          sync.RWMutex
}

// NewInMemoryCommandBus 創建內存命令總線
func NewInMemoryCommandBus(logger *zap.Logger) *InMemoryCommandBus {
	return &InMemoryCommandBus{
		handlers:    make(map[string]interface{}),
		middlewares: make([]Middleware, 0),
		logger:      logger,
	}
}

// Register 註冊命令處理器
func (b *InMemoryCommandBus) Register(commandType string, handler interface{}) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	if _, exists := b.handlers[commandType]; exists {
		return fmt.Errorf("handler for command type %s already registered", commandType)
	}
	
	// 驗證處理器類型
	handlerType := reflect.TypeOf(handler)
	if handlerType.Kind() != reflect.Func {
		return errors.New("handler must be a function")
	}
	
	// 驗證函數簽名
	if handlerType.NumIn() != 2 {
		return errors.New("handler must have exactly 2 parameters: context and command")
	}
	
	if handlerType.In(0) != reflect.TypeOf((*context.Context)(nil)).Elem() {
		return errors.New("handler's first parameter must be context.Context")
	}
	
	b.handlers[commandType] = handler
	b.logger.Info("Command handler registered", zap.String("command_type", commandType))
	
	return nil
}

// Send 發送命令
func (b *InMemoryCommandBus) Send(ctx context.Context, cmd Command) error {
	// 驗證命令
	if err := cmd.Validate(); err != nil {
		return fmt.Errorf("command validation failed: %w", err)
	}
	
	b.mu.RLock()
	handler, exists := b.handlers[cmd.GetCommandType()]
	b.mu.RUnlock()
	
	if !exists {
		return fmt.Errorf("no handler registered for command type: %s", cmd.GetCommandType())
	}
	
	// 記錄命令執行
	b.logger.Info("Executing command",
		zap.String("command_id", cmd.GetCommandID().String()),
		zap.String("command_type", cmd.GetCommandType()))
	
	// 執行處理器
	handlerValue := reflect.ValueOf(handler)
	cmdValue := reflect.ValueOf(cmd)
	ctxValue := reflect.ValueOf(ctx)
	
	// 調用處理器
	results := handlerValue.Call([]reflect.Value{ctxValue, cmdValue})
	
	// 檢查錯誤
	if len(results) > 0 {
		lastResult := results[len(results)-1]
		if !lastResult.IsNil() && lastResult.Type().Implements(reflect.TypeOf((*error)(nil)).Elem()) {
			err := lastResult.Interface().(error)
			b.logger.Error("Command execution failed",
				zap.String("command_id", cmd.GetCommandID().String()),
				zap.String("command_type", cmd.GetCommandType()),
				zap.Error(err))
			return err
		}
	}
	
	b.logger.Info("Command executed successfully",
		zap.String("command_id", cmd.GetCommandID().String()),
		zap.String("command_type", cmd.GetCommandType()))
	
	return nil
}

// SendAsync 異步發送命令
func (b *InMemoryCommandBus) SendAsync(ctx context.Context, cmd Command) <-chan error {
	errChan := make(chan error, 1)
	
	go func() {
		defer close(errChan)
		errChan <- b.Send(ctx, cmd)
	}()
	
	return errChan
}

// Use 添加中間件
func (b *InMemoryCommandBus) Use(middleware Middleware) {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	b.middlewares = append(b.middlewares, middleware)
}

// ValidationMiddleware 驗證中間件
func ValidationMiddleware(logger *zap.Logger) Middleware {
	return func(next Handler[Command]) Handler[Command] {
		return HandlerFunc[Command](func(ctx context.Context, cmd Command) error {
			logger.Debug("Validating command",
				zap.String("command_type", cmd.GetCommandType()))
			
			if err := cmd.Validate(); err != nil {
				logger.Error("Command validation failed",
					zap.String("command_type", cmd.GetCommandType()),
					zap.Error(err))
				return fmt.Errorf("validation failed: %w", err)
			}
			
			return next.Handle(ctx, cmd)
		})
	}
}

// LoggingMiddleware 日誌中間件
func LoggingMiddleware(logger *zap.Logger) Middleware {
	return func(next Handler[Command]) Handler[Command] {
		return HandlerFunc[Command](func(ctx context.Context, cmd Command) error {
			logger.Info("Command received",
				zap.String("command_id", cmd.GetCommandID().String()),
				zap.String("command_type", cmd.GetCommandType()),
				zap.Time("timestamp", cmd.GetTimestamp()))
			
			err := next.Handle(ctx, cmd)
			
			if err != nil {
				logger.Error("Command failed",
					zap.String("command_id", cmd.GetCommandID().String()),
					zap.String("command_type", cmd.GetCommandType()),
					zap.Error(err))
			} else {
				logger.Info("Command completed",
					zap.String("command_id", cmd.GetCommandID().String()),
					zap.String("command_type", cmd.GetCommandType()))
			}
			
			return err
		})
	}
}

// MetricsMiddleware 指標中間件
func MetricsMiddleware(recordMetric func(commandType string, success bool, duration float64)) Middleware {
	return func(next Handler[Command]) Handler[Command] {
		return HandlerFunc[Command](func(ctx context.Context, cmd Command) error {
			start := time.Now()
			
			err := next.Handle(ctx, cmd)
			
			duration := time.Since(start).Seconds()
			success := err == nil
			
			recordMetric(cmd.GetCommandType(), success, duration)
			
			return err
		})
	}
}