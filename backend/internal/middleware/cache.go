package middleware

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/fastenmind/fastener-api/internal/infrastructure/cache"
	"github.com/labstack/echo/v4"
)

// CacheConfig defines cache middleware configuration
type CacheConfig struct {
	// Skipper defines a function to skip middleware
	Skipper func(c echo.Context) bool
	
	// Cache service
	Cache *cache.RedisCache
	
	// TTL for cached responses
	TTL time.Duration
	
	// Methods to cache (default: GET)
	Methods []string
	
	// Key generator function
	KeyGenerator func(c echo.Context) string
}

// DefaultCacheConfig returns default cache configuration
var DefaultCacheConfig = CacheConfig{
	Skipper: func(c echo.Context) bool {
		return false
	},
	TTL:     5 * time.Minute,
	Methods: []string{http.MethodGet},
	KeyGenerator: func(c echo.Context) string {
		// Generate cache key from method, path, and query string
		method := c.Request().Method
		path := c.Request().URL.Path
		query := c.Request().URL.RawQuery
		
		key := fmt.Sprintf("%s:%s", method, path)
		if query != "" {
			key += "?" + query
		}
		
		// Hash the key if it's too long
		if len(key) > 200 {
			h := md5.New()
			h.Write([]byte(key))
			key = hex.EncodeToString(h.Sum(nil))
		}
		
		return key
	},
}

// CachedResponse represents a cached HTTP response
type CachedResponse struct {
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
	Body    []byte            `json:"body"`
}

// CacheMiddleware returns a cache middleware
func CacheMiddleware(config CacheConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultCacheConfig.Skipper
	}
	if config.Cache == nil {
		panic("cache middleware: Cache is required")
	}
	if len(config.Methods) == 0 {
		config.Methods = DefaultCacheConfig.Methods
	}
	if config.KeyGenerator == nil {
		config.KeyGenerator = DefaultCacheConfig.KeyGenerator
	}
	if config.TTL == 0 {
		config.TTL = DefaultCacheConfig.TTL
	}
	
	// Convert methods to map for faster lookup
	methodsMap := make(map[string]bool)
	for _, method := range config.Methods {
		methodsMap[method] = true
	}
	
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Skip if configured
			if config.Skipper(c) {
				return next(c)
			}
			
			// Check if method should be cached
			if !methodsMap[c.Request().Method] {
				return next(c)
			}
			
			// Generate cache key
			key := config.KeyGenerator(c)
			if key == "" {
				return next(c)
			}
			
			// Add cache headers
			c.Response().Header().Set("X-Cache-Key", key)
			
			// Try to get from cache
			ctx := c.Request().Context()
			var cached CachedResponse
			
			err := config.Cache.GetJSON(ctx, key, &cached)
			if err == nil && cached.Status > 0 {
				// Cache hit
				c.Response().Header().Set("X-Cache-Status", "HIT")
				
				// Set headers
				for k, v := range cached.Headers {
					c.Response().Header().Set(k, v)
				}
				
				// Write response
				return c.Blob(cached.Status, "application/json", cached.Body)
			}
			
			// Cache miss
			c.Response().Header().Set("X-Cache-Status", "MISS")
			
			// Capture response
			rec := &responseRecorder{
				ResponseWriter: c.Response().Writer,
				status:         http.StatusOK,
				body:           new(bytes.Buffer),
			}
			c.Response().Writer = rec
			
			// Process request
			err = next(c)
			if err != nil {
				return err
			}
			
			// Cache successful responses only
			if rec.status >= 200 && rec.status < 300 {
				cached := CachedResponse{
					Status:  rec.status,
					Headers: make(map[string]string),
					Body:    rec.body.Bytes(),
				}
				
				// Copy cacheable headers
				for k, v := range rec.Header() {
					if isCacheableHeader(k) && len(v) > 0 {
						cached.Headers[k] = v[0]
					}
				}
				
				// Store in cache
				go func() {
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					defer cancel()
					
					if err := config.Cache.SetJSON(ctx, key, cached, config.TTL); err != nil {
						// Log error but don't fail the request
						fmt.Printf("cache middleware: failed to cache response: %v\n", err)
					}
				}()
			}
			
			return nil
		}
	}
}

// responseRecorder captures the response for caching
type responseRecorder struct {
	http.ResponseWriter
	status int
	body   *bytes.Buffer
}

func (r *responseRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// isCacheableHeader checks if a header should be cached
func isCacheableHeader(header string) bool {
	// Headers to exclude from cache
	excludedHeaders := map[string]bool{
		"Set-Cookie":     true,
		"Authorization":  true,
		"X-Cache-Status": true,
		"X-Cache-Key":    true,
	}
	
	header = strings.ToLower(header)
	for excluded := range excludedHeaders {
		if strings.ToLower(excluded) == header {
			return false
		}
	}
	
	return true
}

// CacheInvalidationMiddleware provides cache invalidation on write operations
func CacheInvalidationMiddleware(cache *cache.RedisCache) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Process request
			err := next(c)
			
			// Invalidate cache on successful write operations
			if err == nil && isWriteMethod(c.Request().Method) && isSuccessStatus(c.Response().Status) {
				go invalidateRelatedCache(c, cache)
			}
			
			return err
		}
	}
}

// isWriteMethod checks if the HTTP method is a write operation
func isWriteMethod(method string) bool {
	writeMethods := map[string]bool{
		http.MethodPost:   true,
		http.MethodPut:    true,
		http.MethodPatch:  true,
		http.MethodDelete: true,
	}
	return writeMethods[method]
}

// isSuccessStatus checks if the status code indicates success
func isSuccessStatus(status int) bool {
	return status >= 200 && status < 300
}

// invalidateRelatedCache invalidates cache entries related to the current request
func invalidateRelatedCache(c echo.Context, cache *cache.RedisCache) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	path := c.Request().URL.Path
	
	// Extract resource type from path (e.g., /api/v1/inquiries/123 -> inquiries)
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 3 {
		resource := parts[2]
		
		// Invalidate list endpoints
		patterns := []string{
			fmt.Sprintf("GET:/api/v1/%s", resource),
			fmt.Sprintf("GET:/api/v1/%s*", resource),
		}
		
		for _, pattern := range patterns {
			if err := cache.DeletePattern(ctx, pattern); err != nil {
				fmt.Printf("cache invalidation: failed to delete pattern %s: %v\n", pattern, err)
			}
		}
	}
}

// ConditionalCacheMiddleware caches responses based on custom conditions
func ConditionalCacheMiddleware(cache *cache.RedisCache, condition func(c echo.Context) bool) echo.MiddlewareFunc {
	config := CacheConfig{
		Cache: cache,
		Skipper: func(c echo.Context) bool {
			return !condition(c)
		},
	}
	return CacheMiddleware(config)
}