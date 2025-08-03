package command

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// TestCommand 測試命令
type TestCommand struct {
	BaseCommand
	Value string
}

func (c TestCommand) Validate() error {
	if c.Value == "" {
		return errors.New("value is required")
	}
	return nil
}

func TestInMemoryCommandBus_Register(t *testing.T) {
	logger := zap.NewNop()
	bus := NewInMemoryCommandBus(logger)

	// Test successful registration
	t.Run("Successful registration", func(t *testing.T) {
		handler := func(ctx context.Context, cmd TestCommand) error {
			return nil
		}
		
		err := bus.Register("TestCommand", handler)
		assert.NoError(t, err)
	})

	// Test duplicate registration
	t.Run("Duplicate registration", func(t *testing.T) {
		handler := func(ctx context.Context, cmd TestCommand) error {
			return nil
		}
		
		bus.Register("DuplicateCommand", handler)
		err := bus.Register("DuplicateCommand", handler)
		
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already registered")
	})

	// Test invalid handler type
	t.Run("Invalid handler type", func(t *testing.T) {
		err := bus.Register("InvalidCommand", "not a function")
		
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "handler must be a function")
	})

	// Test invalid handler signature
	t.Run("Invalid handler signature", func(t *testing.T) {
		handler := func(cmd TestCommand) error {
			return nil
		}
		
		err := bus.Register("InvalidSignature", handler)
		
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "handler must have exactly 2 parameters")
	})
}

func TestInMemoryCommandBus_Send(t *testing.T) {
	logger := zap.NewNop()
	bus := NewInMemoryCommandBus(logger)

	// Test successful command execution
	t.Run("Successful execution", func(t *testing.T) {
		executed := false
		handler := func(ctx context.Context, cmd TestCommand) error {
			executed = true
			assert.Equal(t, "test value", cmd.Value)
			return nil
		}
		
		bus.Register("TestCommand", handler)
		
		cmd := TestCommand{
			BaseCommand: NewBaseCommand("TestCommand"),
			Value:       "test value",
		}
		
		err := bus.Send(context.Background(), cmd)
		
		assert.NoError(t, err)
		assert.True(t, executed)
	})

	// Test command validation failure
	t.Run("Validation failure", func(t *testing.T) {
		handler := func(ctx context.Context, cmd TestCommand) error {
			return nil
		}
		
		bus.Register("ValidationTest", handler)
		
		cmd := TestCommand{
			BaseCommand: NewBaseCommand("ValidationTest"),
			Value:       "", // Invalid
		}
		
		err := bus.Send(context.Background(), cmd)
		
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "value is required")
	})

	// Test handler not found
	t.Run("Handler not found", func(t *testing.T) {
		cmd := TestCommand{
			BaseCommand: NewBaseCommand("UnregisteredCommand"),
			Value:       "test",
		}
		
		err := bus.Send(context.Background(), cmd)
		
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no handler registered")
	})

	// Test handler error
	t.Run("Handler error", func(t *testing.T) {
		handler := func(ctx context.Context, cmd TestCommand) error {
			return errors.New("handler error")
		}
		
		bus.Register("ErrorCommand", handler)
		
		cmd := TestCommand{
			BaseCommand: NewBaseCommand("ErrorCommand"),
			Value:       "test",
		}
		
		err := bus.Send(context.Background(), cmd)
		
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "handler error")
	})
}

func TestInMemoryCommandBus_SendAsync(t *testing.T) {
	logger := zap.NewNop()
	bus := NewInMemoryCommandBus(logger)

	// Test async execution
	t.Run("Async execution", func(t *testing.T) {
		handler := func(ctx context.Context, cmd TestCommand) error {
			time.Sleep(50 * time.Millisecond)
			return nil
		}
		
		bus.Register("AsyncCommand", handler)
		
		cmd := TestCommand{
			BaseCommand: NewBaseCommand("AsyncCommand"),
			Value:       "test",
		}
		
		errChan := bus.SendAsync(context.Background(), cmd)
		
		select {
		case err := <-errChan:
			assert.NoError(t, err)
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timeout waiting for async command")
		}
	})
}

func TestValidationMiddleware(t *testing.T) {
	logger := zap.NewNop()
	middleware := ValidationMiddleware(logger)
	
	// Test with valid command
	t.Run("Valid command", func(t *testing.T) {
		executed := false
		next := HandlerFunc[Command](func(ctx context.Context, cmd Command) error {
			executed = true
			return nil
		})
		
		handler := middleware(next)
		
		cmd := TestCommand{
			BaseCommand: NewBaseCommand("TestCommand"),
			Value:       "valid",
		}
		
		err := handler.Handle(context.Background(), cmd)
		
		assert.NoError(t, err)
		assert.True(t, executed)
	})
	
	// Test with invalid command
	t.Run("Invalid command", func(t *testing.T) {
		executed := false
		next := HandlerFunc[Command](func(ctx context.Context, cmd Command) error {
			executed = true
			return nil
		})
		
		handler := middleware(next)
		
		cmd := TestCommand{
			BaseCommand: NewBaseCommand("TestCommand"),
			Value:       "", // Invalid
		}
		
		err := handler.Handle(context.Background(), cmd)
		
		assert.Error(t, err)
		assert.False(t, executed)
		assert.Contains(t, err.Error(), "validation failed")
	})
}

func TestMetricsMiddleware(t *testing.T) {
	var recordedType string
	var recordedSuccess bool
	var recordedDuration float64
	
	recordMetric := func(commandType string, success bool, duration float64) {
		recordedType = commandType
		recordedSuccess = success
		recordedDuration = duration
	}
	
	middleware := MetricsMiddleware(recordMetric)
	
	// Test successful command
	t.Run("Successful command", func(t *testing.T) {
		next := HandlerFunc[Command](func(ctx context.Context, cmd Command) error {
			time.Sleep(10 * time.Millisecond)
			return nil
		})
		
		handler := middleware(next)
		
		cmd := TestCommand{
			BaseCommand: NewBaseCommand("MetricsTest"),
			Value:       "test",
		}
		
		err := handler.Handle(context.Background(), cmd)
		
		assert.NoError(t, err)
		assert.Equal(t, "MetricsTest", recordedType)
		assert.True(t, recordedSuccess)
		assert.Greater(t, recordedDuration, 0.0)
	})
	
	// Test failed command
	t.Run("Failed command", func(t *testing.T) {
		next := HandlerFunc[Command](func(ctx context.Context, cmd Command) error {
			return errors.New("command failed")
		})
		
		handler := middleware(next)
		
		cmd := TestCommand{
			BaseCommand: NewBaseCommand("FailedMetricsTest"),
			Value:       "test",
		}
		
		err := handler.Handle(context.Background(), cmd)
		
		assert.Error(t, err)
		assert.Equal(t, "FailedMetricsTest", recordedType)
		assert.False(t, recordedSuccess)
	})
}