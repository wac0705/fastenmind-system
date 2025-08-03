package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/fastenmind/fastener-api/pkg/messaging"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

// RabbitMQBroker RabbitMQ 訊息代理實現
type RabbitMQBroker struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	connMu       sync.RWMutex  // 保護連接和通道
	url          string
	exchange     string
	exchangeType string
	durable      bool
	autoDelete   bool
	
	subscribers map[string]*subscriber
	subMu       sync.RWMutex  // 保護訂閱者映射
	
	reconnectDelay time.Duration
	maxReconnect   int
	
	logger      Logger
	retryPolicy messaging.RetryPolicy
}

// Logger 日誌介面
type Logger interface {
	Info(msg string, fields ...interface{})
	Error(msg string, err error, fields ...interface{})
	Debug(msg string, fields ...interface{})
}

// Config RabbitMQ 配置
type Config struct {
	URL            string
	Exchange       string
	ExchangeType   string
	Durable        bool
	AutoDelete     bool
	ReconnectDelay time.Duration
	MaxReconnect   int
	RetryPolicy    messaging.RetryPolicy
	Logger         Logger
}

// subscriber 訂閱者結構
type subscriber struct {
	queue      string
	routingKey string
	handler    messaging.MessageHandler
	consumer   string
	cancelFunc context.CancelFunc
}

// NewRabbitMQBroker 創建 RabbitMQ 訊息代理
func NewRabbitMQBroker(config Config) *RabbitMQBroker {
	if config.ExchangeType == "" {
		config.ExchangeType = "topic"
	}
	if config.ReconnectDelay == 0 {
		config.ReconnectDelay = 5 * time.Second
	}
	if config.MaxReconnect == 0 {
		config.MaxReconnect = 5
	}
	
	return &RabbitMQBroker{
		url:            config.URL,
		exchange:       config.Exchange,
		exchangeType:   config.ExchangeType,
		durable:        config.Durable,
		autoDelete:     config.AutoDelete,
		subscribers:    make(map[string]*subscriber),
		reconnectDelay: config.ReconnectDelay,
		maxReconnect:   config.MaxReconnect,
		logger:         config.Logger,
		retryPolicy:    config.RetryPolicy,
	}
}

// Start 啟動訊息代理
func (b *RabbitMQBroker) Start(ctx context.Context) error {
	if err := b.connect(); err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	
	// 宣告交換器
	if err := b.channel.ExchangeDeclare(
		b.exchange,
		b.exchangeType,
		b.durable,
		b.autoDelete,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}
	
	// 啟動連接監控
	go b.handleReconnect(ctx)
	
	return nil
}

// Stop 停止訊息代理
func (b *RabbitMQBroker) Stop(ctx context.Context) error {
	// 停止所有訂閱者
	b.subMu.Lock()
	for _, sub := range b.subscribers {
		if sub.cancelFunc != nil {
			sub.cancelFunc()
		}
	}
	b.subMu.Unlock()
	
	// 關閉通道和連接
	b.connMu.Lock()
	defer b.connMu.Unlock()
	
	if b.channel != nil {
		b.channel.Close()
	}
	if b.conn != nil {
		b.conn.Close()
	}
	
	return nil
}

// Publish 發布訊息
func (b *RabbitMQBroker) Publish(ctx context.Context, topic string, message messaging.Message) error {
	// 序列化訊息
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}
	
	// 設置訊息屬性
	publishing := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Body:         body,
		Timestamp:    message.GetTimestamp(),
		MessageId:    message.GetID(),
		Headers:      b.convertHeaders(message.GetHeaders()),
	}
	
	// 發布訊息（使用讀鎖保護通道）
	b.connMu.RLock()
	defer b.connMu.RUnlock()
	
	if b.channel == nil {
		return fmt.Errorf("channel is not initialized")
	}
	
	err = b.channel.Publish(
		b.exchange,
		topic,
		false,
		false,
		publishing,
	)
	
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}
	
	b.logger.Debug("Message published", "topic", topic, "messageId", message.GetID())
	
	return nil
}

// PublishBatch 批量發布訊息
func (b *RabbitMQBroker) PublishBatch(ctx context.Context, topic string, messages []messaging.Message) error {
	for _, msg := range messages {
		if err := b.Publish(ctx, topic, msg); err != nil {
			return err
		}
	}
	return nil
}

// Subscribe 訂閱主題
func (b *RabbitMQBroker) Subscribe(ctx context.Context, topic string, handler messaging.MessageHandler) error {
	b.subMu.Lock()
	defer b.subMu.Unlock()
	
	// 生成佇列名稱
	queueName := fmt.Sprintf("%s.%s.%s", b.exchange, topic, uuid.New().String())
	
	// 宣告佇列
	queue, err := b.channel.QueueDeclare(
		queueName,
		b.durable,
		b.autoDelete,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}
	
	// 綁定佇列到交換器
	err = b.channel.QueueBind(
		queue.Name,
		topic,
		b.exchange,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}
	
	// 創建消費者
	consumerTag := fmt.Sprintf("consumer-%s", uuid.New().String())
	msgs, err := b.channel.Consume(
		queue.Name,
		consumerTag,
		false, // 手動確認
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to consume messages: %w", err)
	}
	
	// 創建訂閱者上下文
	subCtx, cancel := context.WithCancel(ctx)
	
	// 保存訂閱者資訊
	sub := &subscriber{
		queue:      queue.Name,
		routingKey: topic,
		handler:    handler,
		consumer:   consumerTag,
		cancelFunc: cancel,
	}
	b.subscribers[topic] = sub
	
	// 啟動訊息處理
	go b.handleMessages(subCtx, msgs, handler)
	
	b.logger.Info("Subscribed to topic", "topic", topic, "queue", queue.Name)
	
	return nil
}

// Unsubscribe 取消訂閱
func (b *RabbitMQBroker) Unsubscribe(ctx context.Context, topic string) error {
	b.subMu.Lock()
	defer b.subMu.Unlock()
	
	sub, exists := b.subscribers[topic]
	if !exists {
		return fmt.Errorf("no subscription found for topic: %s", topic)
	}
	
	// 取消消費者
	if err := b.channel.Cancel(sub.consumer, false); err != nil {
		return fmt.Errorf("failed to cancel consumer: %w", err)
	}
	
	// 取消上下文
	if sub.cancelFunc != nil {
		sub.cancelFunc()
	}
	
	// 刪除佇列
	if _, err := b.channel.QueueDelete(sub.queue, false, false, false); err != nil {
		return fmt.Errorf("failed to delete queue: %w", err)
	}
	
	delete(b.subscribers, topic)
	
	b.logger.Info("Unsubscribed from topic", "topic", topic)
	
	return nil
}

// connect 建立連接
func (b *RabbitMQBroker) connect() error {
	conn, err := amqp.Dial(b.url)
	if err != nil {
		return err
	}
	
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return err
	}
	
	// 使用鎖保護連接和通道的更新
	b.connMu.Lock()
	b.conn = conn
	b.channel = ch
	b.connMu.Unlock()
	
	return nil
}

// handleReconnect 處理重連
func (b *RabbitMQBroker) handleReconnect(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case err := <-b.conn.NotifyClose(make(chan *amqp.Error)):
			if err != nil {
				b.logger.Error("Connection closed", err, "url", b.url)
				b.reconnect(ctx)
			}
		}
	}
}

// reconnect 重新連接
func (b *RabbitMQBroker) reconnect(ctx context.Context) {
	for i := 0; i < b.maxReconnect; i++ {
		select {
		case <-ctx.Done():
			return
		case <-time.After(b.reconnectDelay):
			b.logger.Info("Attempting to reconnect", "attempt", i+1)
			
			if err := b.connect(); err != nil {
				b.logger.Error("Reconnection failed", err, "attempt", i+1)
				continue
			}
			
			// 重新宣告交換器
			if err := b.channel.ExchangeDeclare(
				b.exchange,
				b.exchangeType,
				b.durable,
				b.autoDelete,
				false,
				false,
				nil,
			); err != nil {
				b.logger.Error("Failed to declare exchange", err)
				continue
			}
			
			// 重新訂閱
			if err := b.resubscribe(ctx); err != nil {
				b.logger.Error("Failed to resubscribe", err)
				continue
			}
			
			b.logger.Info("Reconnected successfully")
			return
		}
	}
}

// resubscribe 重新訂閱
func (b *RabbitMQBroker) resubscribe(ctx context.Context) error {
	b.subMu.Lock()
	defer b.subMu.Unlock()
	
	for topic, sub := range b.subscribers {
		if err := b.Subscribe(ctx, topic, sub.handler); err != nil {
			return err
		}
	}
	
	return nil
}

// handleMessages 處理訊息
func (b *RabbitMQBroker) handleMessages(ctx context.Context, msgs <-chan amqp.Delivery, handler messaging.MessageHandler) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-msgs:
			go b.processMessage(ctx, msg, handler)
		}
	}
}

// processMessage 處理單個訊息
func (b *RabbitMQBroker) processMessage(ctx context.Context, delivery amqp.Delivery, handler messaging.MessageHandler) {
	// 反序列化訊息
	var message messaging.BaseMessage
	if err := json.Unmarshal(delivery.Body, &message); err != nil {
		b.logger.Error("Failed to unmarshal message", err, "messageId", delivery.MessageId)
		delivery.Nack(false, false)
		return
	}
	
	// 處理訊息
	if err := handler(ctx, &message); err != nil {
		b.logger.Error("Failed to handle message", err, "messageId", message.ID)
		
		// 重試邏輯
		retryCount := b.getRetryCount(delivery.Headers)
		if retryCount < b.retryPolicy.MaxRetries {
			b.retryMessage(ctx, delivery, retryCount+1)
		} else {
			// 發送到死信佇列
			b.sendToDeadLetter(ctx, &message, err)
			delivery.Ack(false)
		}
		return
	}
	
	// 確認訊息
	delivery.Ack(false)
}

// getRetryCount 獲取重試次數
func (b *RabbitMQBroker) getRetryCount(headers amqp.Table) int {
	if headers == nil {
		return 0
	}
	
	if count, ok := headers["x-retry-count"].(int); ok {
		return count
	}
	
	return 0
}

// retryMessage 重試訊息
func (b *RabbitMQBroker) retryMessage(ctx context.Context, delivery amqp.Delivery, retryCount int) {
	// 計算延遲時間
	delay := b.calculateRetryDelay(retryCount)
	
	// 更新重試次數
	if delivery.Headers == nil {
		delivery.Headers = make(amqp.Table)
	}
	delivery.Headers["x-retry-count"] = retryCount
	
	// 延遲後重新發布
	go func() {
		select {
		case <-ctx.Done():
			return
		case <-time.After(delay):
			publishing := amqp.Publishing{
				DeliveryMode: delivery.DeliveryMode,
				ContentType:  delivery.ContentType,
				Body:         delivery.Body,
				Timestamp:    delivery.Timestamp,
				MessageId:    delivery.MessageId,
				Headers:      delivery.Headers,
			}
			
			if err := b.channel.Publish(
				b.exchange,
				delivery.RoutingKey,
				false,
				false,
				publishing,
			); err != nil {
				b.logger.Error("Failed to retry message", err, "messageId", delivery.MessageId)
			}
		}
	}()
	
	// 確認原始訊息
	delivery.Ack(false)
}

// calculateRetryDelay 計算重試延遲
func (b *RabbitMQBroker) calculateRetryDelay(retryCount int) time.Duration {
	delay := b.retryPolicy.InitialDelay
	
	for i := 1; i < retryCount; i++ {
		delay = time.Duration(float64(delay) * b.retryPolicy.Multiplier)
		if delay > b.retryPolicy.MaxDelay {
			delay = b.retryPolicy.MaxDelay
			break
		}
	}
	
	return delay
}

// sendToDeadLetter 發送到死信佇列
func (b *RabbitMQBroker) sendToDeadLetter(ctx context.Context, message messaging.Message, err error) {
	// 創建死信訊息
	deadLetterMsg := &messaging.BaseMessage{
		ID:        uuid.New().String(),
		Type:      "dead_letter",
		Timestamp: time.Now(),
		Payload: map[string]interface{}{
			"original_message": message,
			"error":           err.Error(),
			"timestamp":       time.Now(),
		},
		Headers: map[string]string{
			"original_id":   message.GetID(),
			"original_type": string(message.GetType()),
		},
	}
	
	// 發布到死信佇列
	if err := b.Publish(ctx, "dead_letter", deadLetterMsg); err != nil {
		b.logger.Error("Failed to send to dead letter queue", err, "originalMessageId", message.GetID())
	}
}

// convertHeaders 轉換標頭
func (b *RabbitMQBroker) convertHeaders(headers map[string]string) amqp.Table {
	table := make(amqp.Table)
	for k, v := range headers {
		table[k] = v
	}
	return table
}