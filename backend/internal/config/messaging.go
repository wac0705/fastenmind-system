package config

import (
	"os"
	"strconv"
)

// MessagingConfig 訊息配置
type MessagingConfig struct {
	RabbitMQ RabbitMQConfig `json:"rabbitmq" yaml:"rabbitmq"`
	Kafka    KafkaConfig    `json:"kafka" yaml:"kafka"`
}

// RabbitMQConfig RabbitMQ 配置
type RabbitMQConfig struct {
	URL            string `json:"url" yaml:"url" env:"RABBITMQ_URL" envDefault:"amqp://guest:guest@localhost:5672/"`
	Exchange       string `json:"exchange" yaml:"exchange" env:"RABBITMQ_EXCHANGE" envDefault:"fastenmind"`
	ExchangeType   string `json:"exchange_type" yaml:"exchange_type" env:"RABBITMQ_EXCHANGE_TYPE" envDefault:"topic"`
	Durable        bool   `json:"durable" yaml:"durable" env:"RABBITMQ_DURABLE" envDefault:"true"`
	AutoDelete     bool   `json:"auto_delete" yaml:"auto_delete" env:"RABBITMQ_AUTO_DELETE" envDefault:"false"`
	MaxRetries     int    `json:"max_retries" yaml:"max_retries" env:"RABBITMQ_MAX_RETRIES" envDefault:"3"`
	RetryDelay     string `json:"retry_delay" yaml:"retry_delay" env:"RABBITMQ_RETRY_DELAY" envDefault:"5s"`
	MaxConnections int    `json:"max_connections" yaml:"max_connections" env:"RABBITMQ_MAX_CONNECTIONS" envDefault:"10"`
}

// KafkaConfig Kafka 配置
type KafkaConfig struct {
	Brokers       []string `json:"brokers" yaml:"brokers" env:"KAFKA_BROKERS" envDefault:"localhost:9092"`
	ConsumerGroup string   `json:"consumer_group" yaml:"consumer_group" env:"KAFKA_CONSUMER_GROUP" envDefault:"fastenmind"`
	ProducerRetries int    `json:"producer_retries" yaml:"producer_retries" env:"KAFKA_PRODUCER_RETRIES" envDefault:"3"`
	CompressionType string `json:"compression_type" yaml:"compression_type" env:"KAFKA_COMPRESSION_TYPE" envDefault:"snappy"`
}

// TracingConfig 追蹤配置
type TracingConfig struct {
	Enabled        bool    `json:"enabled" yaml:"enabled" env:"TRACING_ENABLED" envDefault:"true"`
	ServiceName    string  `json:"service_name" yaml:"service_name" env:"TRACING_SERVICE_NAME" envDefault:"fastenmind-api"`
	ServiceVersion string  `json:"service_version" yaml:"service_version" env:"TRACING_SERVICE_VERSION" envDefault:"1.0.0"`
	Environment    string  `json:"environment" yaml:"environment" env:"TRACING_ENVIRONMENT" envDefault:"development"`
	ExporterType   string  `json:"exporter_type" yaml:"exporter_type" env:"TRACING_EXPORTER_TYPE" envDefault:"jaeger"`
	Endpoint       string  `json:"endpoint" yaml:"endpoint" env:"TRACING_ENDPOINT" envDefault:"http://localhost:14268/api/traces"`
	SamplingRate   float64 `json:"sampling_rate" yaml:"sampling_rate" env:"TRACING_SAMPLING_RATE" envDefault:"1.0"`
}

// CQRSConfig CQRS 配置
type CQRSConfig struct {
	EventStore EventStoreConfig `json:"event_store" yaml:"event_store"`
	Snapshot   SnapshotConfig   `json:"snapshot" yaml:"snapshot"`
}

// EventStoreConfig 事件存儲配置
type EventStoreConfig struct {
	Type              string `json:"type" yaml:"type" env:"EVENT_STORE_TYPE" envDefault:"sql"`
	ConnectionString  string `json:"connection_string" yaml:"connection_string" env:"EVENT_STORE_CONNECTION"`
	MaxEventsPerLoad  int    `json:"max_events_per_load" yaml:"max_events_per_load" env:"EVENT_STORE_MAX_EVENTS" envDefault:"1000"`
	SnapshotFrequency int    `json:"snapshot_frequency" yaml:"snapshot_frequency" env:"EVENT_STORE_SNAPSHOT_FREQ" envDefault:"10"`
}

// SnapshotConfig 快照配置
type SnapshotConfig struct {
	Enabled       bool   `json:"enabled" yaml:"enabled" env:"SNAPSHOT_ENABLED" envDefault:"true"`
	StorageType   string `json:"storage_type" yaml:"storage_type" env:"SNAPSHOT_STORAGE_TYPE" envDefault:"redis"`
	RetentionDays int    `json:"retention_days" yaml:"retention_days" env:"SNAPSHOT_RETENTION_DAYS" envDefault:"30"`
}

// Update main config
func (c *Config) LoadMessagingConfig() {
	c.Messaging = MessagingConfig{
		RabbitMQ: RabbitMQConfig{
			URL:            getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
			Exchange:       getEnv("RABBITMQ_EXCHANGE", "fastenmind"),
			ExchangeType:   getEnv("RABBITMQ_EXCHANGE_TYPE", "topic"),
			Durable:        getEnvAsBool("RABBITMQ_DURABLE", true),
			AutoDelete:     getEnvAsBool("RABBITMQ_AUTO_DELETE", false),
			MaxRetries:     getEnvAsInt("RABBITMQ_MAX_RETRIES", 3),
			RetryDelay:     getEnv("RABBITMQ_RETRY_DELAY", "5s"),
			MaxConnections: getEnvAsInt("RABBITMQ_MAX_CONNECTIONS", 10),
		},
	}
	
	c.Tracing = TracingConfig{
		Enabled:        getEnvAsBool("TRACING_ENABLED", true),
		ServiceName:    getEnv("TRACING_SERVICE_NAME", "fastenmind-api"),
		ServiceVersion: getEnv("TRACING_SERVICE_VERSION", "1.0.0"),
		Environment:    getEnv("TRACING_ENVIRONMENT", "development"),
		ExporterType:   getEnv("TRACING_EXPORTER_TYPE", "jaeger"),
		Endpoint:       getEnv("TRACING_ENDPOINT", "http://localhost:14268/api/traces"),
		SamplingRate:   getEnvAsFloat("TRACING_SAMPLING_RATE", 1.0),
	}
	
	c.CQRS = CQRSConfig{
		EventStore: EventStoreConfig{
			Type:              getEnv("EVENT_STORE_TYPE", "sql"),
			MaxEventsPerLoad:  getEnvAsInt("EVENT_STORE_MAX_EVENTS", 1000),
			SnapshotFrequency: getEnvAsInt("EVENT_STORE_SNAPSHOT_FREQ", 10),
		},
		Snapshot: SnapshotConfig{
			Enabled:       getEnvAsBool("SNAPSHOT_ENABLED", true),
			StorageType:   getEnv("SNAPSHOT_STORAGE_TYPE", "redis"),
			RetentionDays: getEnvAsInt("SNAPSHOT_RETENTION_DAYS", 30),
		},
	}
}

// Helper functions
func getEnvAsFloat(key string, defaultVal float64) float64 {
	if value, exists := os.LookupEnv(key); exists {
		if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			return floatVal
		}
	}
	return defaultVal
}