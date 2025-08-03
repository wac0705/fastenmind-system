package config

import (
	"os"
	"strconv"
	"time"
)

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host        string
	Port        string
	Password    string
	DB          int
	MaxRetries  int
	PoolSize    int
	TTL         time.Duration
	EnableCache bool
}

// LoadRedisConfig loads Redis configuration from environment variables
func LoadRedisConfig() *RedisConfig {
	db, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	maxRetries, _ := strconv.Atoi(getEnv("REDIS_MAX_RETRIES", "3"))
	poolSize, _ := strconv.Atoi(getEnv("REDIS_POOL_SIZE", "10"))
	ttlMinutes, _ := strconv.Atoi(getEnv("REDIS_TTL_MINUTES", "30"))
	enableCache, _ := strconv.ParseBool(getEnv("ENABLE_CACHE", "true"))

	return &RedisConfig{
		Host:        getEnv("REDIS_HOST", "localhost"),
		Port:        getEnv("REDIS_PORT", "6379"),
		Password:    getEnv("REDIS_PASSWORD", ""),
		DB:          db,
		MaxRetries:  maxRetries,
		PoolSize:    poolSize,
		TTL:         time.Duration(ttlMinutes) * time.Minute,
		EnableCache: enableCache,
	}
}

// GetAddr returns the Redis address in host:port format
func (c *RedisConfig) GetAddr() string {
	return c.Host + ":" + c.Port
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}