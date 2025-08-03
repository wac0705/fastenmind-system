package rabbitmq

import (
	"context"
	"fmt"
	"time"

	"github.com/fastenmind/fastener-api/pkg/concurrent"
	"github.com/fastenmind/fastener-api/pkg/messaging"
)

// RabbitMQService 實現服務介面的 RabbitMQ 代理
type RabbitMQService struct {
	broker      *RabbitMQBroker
	name        string
	status      concurrent.ServiceStatus
	stopChannel chan struct{}
}

// NewRabbitMQService 創建新的 RabbitMQ 服務
func NewRabbitMQService(config Config) *RabbitMQService {
	broker := NewRabbitMQBroker(config)
	
	return &RabbitMQService{
		broker:      broker,
		name:        "rabbitmq-broker",
		status:      concurrent.StatusStopped,
		stopChannel: make(chan struct{}),
	}
}

// Start 啟動服務
func (s *RabbitMQService) Start(ctx context.Context) error {
	if s.status == concurrent.StatusRunning {
		return nil
	}
	
	s.status = concurrent.StatusStarting
	
	// 啟動代理
	if err := s.broker.Start(ctx); err != nil {
		s.status = concurrent.StatusError
		return fmt.Errorf("failed to start RabbitMQ broker: %w", err)
	}
	
	s.status = concurrent.StatusRunning
	
	// 啟動健康檢查
	go s.healthCheck(ctx)
	
	return nil
}

// Stop 停止服務
func (s *RabbitMQService) Stop(ctx context.Context) error {
	if s.status != concurrent.StatusRunning {
		return nil
	}
	
	s.status = concurrent.StatusStopping
	
	// 發送停止信號
	close(s.stopChannel)
	
	// 停止代理
	if err := s.broker.Stop(ctx); err != nil {
		s.status = concurrent.StatusError
		return fmt.Errorf("failed to stop RabbitMQ broker: %w", err)
	}
	
	s.status = concurrent.StatusStopped
	return nil
}

// Name 返回服務名稱
func (s *RabbitMQService) Name() string {
	return s.name
}

// Status 返回服務狀態
func (s *RabbitMQService) Status() concurrent.ServiceStatus {
	return s.status
}

// GetBroker 獲取底層代理
func (s *RabbitMQService) GetBroker() messaging.MessageBroker {
	return s.broker
}

// healthCheck 健康檢查
func (s *RabbitMQService) healthCheck(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-s.stopChannel:
			return
		case <-ticker.C:
			// 檢查連接狀態
			s.broker.connMu.RLock()
			conn := s.broker.conn
			s.broker.connMu.RUnlock()
			
			if conn != nil && !conn.IsClosed() {
				s.broker.logger.Debug("RabbitMQ health check passed")
			} else {
				s.broker.logger.Error("RabbitMQ connection is closed", nil)
				s.status = concurrent.StatusError
			}
		}
	}
}

// MessageBrokerService 訊息代理服務管理器
type MessageBrokerService struct {
	service    concurrent.Service
	broker     messaging.MessageBroker
	registry   *concurrent.ServiceRegistry
}

// NewMessageBrokerService 創建訊息代理服務管理器
func NewMessageBrokerService(config Config, registry *concurrent.ServiceRegistry) (*MessageBrokerService, error) {
	// 創建 RabbitMQ 服務
	service := NewRabbitMQService(config)
	
	// 註冊服務
	if err := registry.Register(service); err != nil {
		return nil, fmt.Errorf("failed to register RabbitMQ service: %w", err)
	}
	
	return &MessageBrokerService{
		service:  service,
		broker:   service.GetBroker(),
		registry: registry,
	}, nil
}

// Start 啟動服務
func (m *MessageBrokerService) Start(ctx context.Context) error {
	return m.registry.Start(ctx, m.service.Name())
}

// Stop 停止服務
func (m *MessageBrokerService) Stop(ctx context.Context) error {
	return m.registry.Stop(ctx, m.service.Name())
}

// GetBroker 獲取訊息代理
func (m *MessageBrokerService) GetBroker() messaging.MessageBroker {
	return m.broker
}

// GetStatus 獲取服務狀態
func (m *MessageBrokerService) GetStatus() concurrent.ServiceStatus {
	return m.service.Status()
}

// WaitForReady 等待服務就緒
func (m *MessageBrokerService) WaitForReady(ctx context.Context, timeout time.Duration) error {
	return m.registry.WaitForService(ctx, m.service.Name(), concurrent.StatusRunning, timeout)
}