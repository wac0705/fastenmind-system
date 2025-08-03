package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

// RedisCache implements a cache using Redis
type RedisCache struct {
	client *redis.Client
	prefix string
	ttl    time.Duration
}

// NewRedisCache creates a new Redis cache
func NewRedisCache(addr, password string, db int, prefix string, ttl time.Duration) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	
	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}
	
	return &RedisCache{
		client: client,
		prefix: prefix,
		ttl:    ttl,
	}, nil
}

// Close closes the Redis connection
func (c *RedisCache) Close() error {
	return c.client.Close()
}

// buildKey builds a cache key with prefix
func (c *RedisCache) buildKey(key string) string {
	return fmt.Sprintf("%s:%s", c.prefix, key)
}

// Get retrieves a value from cache
func (c *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	result, err := c.client.Get(ctx, c.buildKey(key)).Bytes()
	if err == redis.Nil {
		return nil, nil // Key not found
	}
	return result, err
}

// GetJSON retrieves and unmarshals a JSON value from cache
func (c *RedisCache) GetJSON(ctx context.Context, key string, dest interface{}) error {
	data, err := c.Get(ctx, key)
	if err != nil {
		return err
	}
	if data == nil {
		return nil // Key not found
	}
	return json.Unmarshal(data, dest)
}

// Set stores a value in cache
func (c *RedisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	if ttl == 0 {
		ttl = c.ttl
	}
	return c.client.Set(ctx, c.buildKey(key), value, ttl).Err()
}

// SetJSON marshals and stores a JSON value in cache
func (c *RedisCache) SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.Set(ctx, key, data, ttl)
}

// Delete removes a value from cache
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, c.buildKey(key)).Err()
}

// DeletePattern removes all keys matching a pattern
func (c *RedisCache) DeletePattern(ctx context.Context, pattern string) error {
	iter := c.client.Scan(ctx, 0, c.buildKey(pattern), 0).Iterator()
	var keys []string
	
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	
	if err := iter.Err(); err != nil {
		return err
	}
	
	if len(keys) > 0 {
		return c.client.Del(ctx, keys...).Err()
	}
	
	return nil
}

// Exists checks if a key exists in cache
func (c *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	result, err := c.client.Exists(ctx, c.buildKey(key)).Result()
	return result > 0, err
}

// Expire sets TTL for a key
func (c *RedisCache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return c.client.Expire(ctx, c.buildKey(key), ttl).Err()
}

// CacheService provides high-level caching operations
type CacheService struct {
	cache *RedisCache
}

// NewCacheService creates a new cache service
func NewCacheService(cache *RedisCache) *CacheService {
	return &CacheService{
		cache: cache,
	}
}

// GetInquiry retrieves an inquiry from cache
func (s *CacheService) GetInquiry(ctx context.Context, inquiryID uuid.UUID) (interface{}, error) {
	var inquiry interface{}
	key := fmt.Sprintf("inquiry:%s", inquiryID.String())
	err := s.cache.GetJSON(ctx, key, &inquiry)
	return inquiry, err
}

// SetInquiry stores an inquiry in cache
func (s *CacheService) SetInquiry(ctx context.Context, inquiryID uuid.UUID, inquiry interface{}) error {
	key := fmt.Sprintf("inquiry:%s", inquiryID.String())
	return s.cache.SetJSON(ctx, key, inquiry, 30*time.Minute)
}

// GetQuote retrieves a quote from cache
func (s *CacheService) GetQuote(ctx context.Context, quoteID uuid.UUID) (interface{}, error) {
	var quote interface{}
	key := fmt.Sprintf("quote:%s", quoteID.String())
	err := s.cache.GetJSON(ctx, key, &quote)
	return quote, err
}

// SetQuote stores a quote in cache
func (s *CacheService) SetQuote(ctx context.Context, quoteID uuid.UUID, quote interface{}) error {
	key := fmt.Sprintf("quote:%s", quoteID.String())
	return s.cache.SetJSON(ctx, key, quote, 30*time.Minute)
}

// GetOrder retrieves an order from cache
func (s *CacheService) GetOrder(ctx context.Context, orderID uuid.UUID) (interface{}, error) {
	var order interface{}
	key := fmt.Sprintf("order:%s", orderID.String())
	err := s.cache.GetJSON(ctx, key, &order)
	return order, err
}

// SetOrder stores an order in cache
func (s *CacheService) SetOrder(ctx context.Context, orderID uuid.UUID, order interface{}) error {
	key := fmt.Sprintf("order:%s", orderID.String())
	return s.cache.SetJSON(ctx, key, order, 30*time.Minute)
}

// InvalidateInquiry removes an inquiry from cache
func (s *CacheService) InvalidateInquiry(ctx context.Context, inquiryID uuid.UUID) error {
	key := fmt.Sprintf("inquiry:%s", inquiryID.String())
	return s.cache.Delete(ctx, key)
}

// InvalidateQuote removes a quote from cache
func (s *CacheService) InvalidateQuote(ctx context.Context, quoteID uuid.UUID) error {
	key := fmt.Sprintf("quote:%s", quoteID.String())
	return s.cache.Delete(ctx, key)
}

// InvalidateOrder removes an order from cache
func (s *CacheService) InvalidateOrder(ctx context.Context, orderID uuid.UUID) error {
	key := fmt.Sprintf("order:%s", orderID.String())
	return s.cache.Delete(ctx, key)
}

// GetUserPermissions retrieves user permissions from cache
func (s *CacheService) GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]string, error) {
	var permissions []string
	key := fmt.Sprintf("user:permissions:%s", userID.String())
	err := s.cache.GetJSON(ctx, key, &permissions)
	return permissions, err
}

// SetUserPermissions stores user permissions in cache
func (s *CacheService) SetUserPermissions(ctx context.Context, userID uuid.UUID, permissions []string) error {
	key := fmt.Sprintf("user:permissions:%s", userID.String())
	return s.cache.SetJSON(ctx, key, permissions, 1*time.Hour)
}

// CacheTTL defines cache TTL for different types
var CacheTTL = struct {
	Short  time.Duration
	Medium time.Duration
	Long   time.Duration
}{
	Short:  5 * time.Minute,
	Medium: 30 * time.Minute,
	Long:   2 * time.Hour,
}