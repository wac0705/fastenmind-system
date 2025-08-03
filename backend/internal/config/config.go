package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Redis    RedisConfig
	CORS     CORSConfig
	Upload   UploadConfig
	Email    EmailConfig
	Messaging MessagingConfig
	Tracing   TracingConfig
	CQRS      CQRSConfig
}

type ServerConfig struct {
	Port        string
	Environment string
}

type DatabaseConfig struct {
	// Primary (write) database
	Primary DBConnectionConfig `json:"primary"`
	
	// Read replicas
	Replicas []DBConnectionConfig `json:"replicas"`
	
	// Connection pool settings
	MaxOpenConns    int `json:"max_open_conns"`
	MaxIdleConns    int `json:"max_idle_conns"`
	ConnMaxLifetime int `json:"conn_max_lifetime_minutes"`
	
	// Read/Write separation settings
	ReadWriteSeparation bool `json:"read_write_separation"`
	PreferPrimaryReads  bool `json:"prefer_primary_reads"`
}

// DBConnectionConfig holds configuration for a single database connection
type DBConnectionConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
	SSLMode  string `json:"ssl_mode"`
	
	// Weight for load balancing (higher = more traffic)
	Weight int `json:"weight"`
}

type JWTConfig struct {
	SecretKey              string
	AccessTokenExpiration  string
	RefreshTokenExpiration string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

type UploadConfig struct {
	MaxSize int64
	Path    string
}

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	FromAddress  string
}

func New() *Config {
	cfg := &Config{
		Server: ServerConfig{
			Port:        getEnv("SERVER_PORT", "8080"),
			Environment: getEnv("SERVER_ENV", "development"),
		},
		Database: DatabaseConfig{
			Primary: DBConnectionConfig{
				Host:     getEnv("DB_PRIMARY_HOST", "localhost"),
				Port:     getEnvAsInt("DB_PRIMARY_PORT", 5432),
				Database: getEnv("DB_PRIMARY_NAME", "fastenmind_db"),
				Username: getEnv("DB_PRIMARY_USER", "fastenmind"),
				Password: getEnv("DB_PRIMARY_PASSWORD", "fastenmind123"),
				SSLMode:  getEnv("DB_PRIMARY_SSL_MODE", "disable"),
				Weight:   1,
			},
			MaxOpenConns:        getEnvAsInt("DB_MAX_OPEN_CONNS", 100),
			MaxIdleConns:        getEnvAsInt("DB_MAX_IDLE_CONNS", 10),
			ConnMaxLifetime:     getEnvAsInt("DB_CONN_MAX_LIFETIME_MINUTES", 60),
			ReadWriteSeparation: getEnvAsBool("DB_READ_WRITE_SEPARATION", false),
			PreferPrimaryReads:  getEnvAsBool("DB_PREFER_PRIMARY_READS", false),
		},
		JWT: JWTConfig{
			SecretKey:              getEnv("JWT_SECRET_KEY", "your-secret-key"),
			AccessTokenExpiration:  getEnv("JWT_ACCESS_TOKEN_EXPIRE", "15m"),
			RefreshTokenExpiration: getEnv("JWT_REFRESH_TOKEN_EXPIRE", "7d"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		CORS: CORSConfig{
			AllowedOrigins: strings.Split(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"), ","),
			AllowedMethods: strings.Split(getEnv("CORS_ALLOWED_METHODS", "GET,POST,PUT,DELETE,OPTIONS"), ","),
			AllowedHeaders: strings.Split(getEnv("CORS_ALLOWED_HEADERS", "Accept,Authorization,Content-Type,X-CSRF-Token"), ","),
		},
		Upload: UploadConfig{
			MaxSize: getEnvAsInt64("MAX_UPLOAD_SIZE", 10485760), // 10MB
			Path:    getEnv("UPLOAD_PATH", "./uploads"),
		},
		Email: EmailConfig{
			SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
			SMTPPort:     getEnvAsInt("SMTP_PORT", 587),
			SMTPUser:     getEnv("SMTP_USER", ""),
			SMTPPassword: getEnv("SMTP_PASSWORD", ""),
			FromAddress:  getEnv("SMTP_FROM", "noreply@fastenmind.com"),
		},
	}
	
	// Load messaging and other configs
	cfg.LoadMessagingConfig()
	
	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// DSN returns the database connection string
func (c DBConnectionConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.Username, c.Password, c.Database, c.SSLMode,
	)
}