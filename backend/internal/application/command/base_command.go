package command

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Command 命令接口
type Command interface {
	GetCommandID() uuid.UUID
	GetCommandType() string
	GetTimestamp() time.Time
	Validate() error
}

// BaseCommand 基礎命令
type BaseCommand struct {
	CommandID   uuid.UUID `json:"command_id"`
	CommandType string    `json:"command_type"`
	Timestamp   time.Time `json:"timestamp"`
}

// GetCommandID 獲取命令ID
func (c BaseCommand) GetCommandID() uuid.UUID {
	return c.CommandID
}

// GetCommandType 獲取命令類型
func (c BaseCommand) GetCommandType() string {
	return c.CommandType
}

// GetTimestamp 獲取時間戳
func (c BaseCommand) GetTimestamp() time.Time {
	return c.Timestamp
}

// NewBaseCommand 創建基礎命令
func NewBaseCommand(commandType string) BaseCommand {
	return BaseCommand{
		CommandID:   uuid.New(),
		CommandType: commandType,
		Timestamp:   time.Now(),
	}
}

// Handler 命令處理器接口
type Handler[T Command] interface {
	Handle(ctx context.Context, command T) error
}

// HandlerFunc 命令處理器函數類型
type HandlerFunc[T Command] func(ctx context.Context, command T) error

// Handle 處理命令
func (f HandlerFunc[T]) Handle(ctx context.Context, command T) error {
	return f(ctx, command)
}

// Bus 命令總線接口
type Bus interface {
	// Register 註冊命令處理器
	Register(commandType string, handler interface{}) error
	
	// Send 發送命令
	Send(ctx context.Context, command Command) error
	
	// SendAsync 異步發送命令
	SendAsync(ctx context.Context, command Command) <-chan error
}

// Middleware 命令中間件
type Middleware func(Handler[Command]) Handler[Command]

// Result 命令結果
type Result struct {
	Success bool
	Data    interface{}
	Error   error
}

// NewSuccessResult 創建成功結果
func NewSuccessResult(data interface{}) Result {
	return Result{
		Success: true,
		Data:    data,
	}
}

// NewErrorResult 創建錯誤結果
func NewErrorResult(err error) Result {
	return Result{
		Success: false,
		Error:   err,
	}
}