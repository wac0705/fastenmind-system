package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/fastenmind/fastener-api/internal/infrastructure/cache"
	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"
)

// RateLimiterConfig defines the rate limiter configuration
type RateLimiterConfig struct {
	// Skipper defines a function to skip middleware
	Skipper func(c echo.Context) bool
	
	// Store for distributed rate limiting (Redis)
	Store *cache.RedisCache
	
	// Rate limit per second
	Rate int
	
	// Burst size
	Burst int
	
	// Key generator
	KeyGenerator func(c echo.Context) string
	
	// Error handler
	ErrorHandler func(c echo.Context, err error) error
	
	// Use Redis for distributed rate limiting
	UseRedis bool
}

// DefaultRateLimiterConfig returns default rate limiter configuration
var DefaultRateLimiterConfig = RateLimiterConfig{
	Skipper: func(c echo.Context) bool {
		return false
	},
	Rate:  10,
	Burst: 20,
	KeyGenerator: func(c echo.Context) string {
		// Default to IP-based rate limiting
		return c.RealIP()
	},
	ErrorHandler: func(c echo.Context, err error) error {
		return echo.NewHTTPError(http.StatusTooManyRequests, "rate limit exceeded")
	},
	UseRedis: false,
}

// memoryStore holds in-memory rate limiters
type memoryStore struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
}

var globalMemoryStore = &memoryStore{
	limiters: make(map[string]*rate.Limiter),
}

// RateLimiter returns a rate limiting middleware
func RateLimiter(config RateLimiterConfig) echo.MiddlewareFunc {
	// Apply defaults
	if config.Skipper == nil {
		config.Skipper = DefaultRateLimiterConfig.Skipper
	}
	if config.Rate == 0 {
		config.Rate = DefaultRateLimiterConfig.Rate
	}
	if config.Burst == 0 {
		config.Burst = DefaultRateLimiterConfig.Burst
	}
	if config.KeyGenerator == nil {
		config.KeyGenerator = DefaultRateLimiterConfig.KeyGenerator
	}
	if config.ErrorHandler == nil {
		config.ErrorHandler = DefaultRateLimiterConfig.ErrorHandler
	}
	
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}
			
			key := config.KeyGenerator(c)
			if key == "" {
				return next(c)
			}
			
			// Use Redis for distributed rate limiting
			if config.UseRedis && config.Store != nil {
				allowed, err := checkRedisRateLimit(c, config.Store, key, config.Rate, config.Burst)
				if err != nil {
					return config.ErrorHandler(c, err)
				}
				if !allowed {
					return config.ErrorHandler(c, fmt.Errorf("rate limit exceeded"))
				}
			} else {
				// Use in-memory rate limiting
				limiter := getOrCreateLimiter(key, config.Rate, config.Burst)
				if !limiter.Allow() {
					return config.ErrorHandler(c, fmt.Errorf("rate limit exceeded"))
				}
			}
			
			// Set rate limit headers
			setRateLimitHeaders(c, config.Rate, config.Burst)
			
			return next(c)
		}
	}
}

// getOrCreateLimiter gets or creates a rate limiter for the given key
func getOrCreateLimiter(key string, rate, burst int) *rate.Limiter {
	globalMemoryStore.mu.RLock()
	limiter, exists := globalMemoryStore.limiters[key]
	globalMemoryStore.mu.RUnlock()
	
	if !exists {
		globalMemoryStore.mu.Lock()
		limiter = rate.NewLimiter(rate.Limit(rate), burst)
		globalMemoryStore.limiters[key] = limiter
		globalMemoryStore.mu.Unlock()
	}
	
	return limiter
}

// checkRedisRateLimit checks rate limit using Redis
func checkRedisRateLimit(c echo.Context, store *cache.RedisCache, key string, rateLimit, burst int) (bool, error) {
	ctx := c.Request().Context()
	
	// Use a sliding window algorithm
	now := time.Now().Unix()
	windowStart := now - 60 // 1 minute window
	
	// Redis key for rate limiting
	redisKey := fmt.Sprintf("ratelimit:%s:%d", key, now/60)
	
	// Get current count
	countStr, err := store.Get(ctx, redisKey)
	if err != nil {
		return false, err
	}
	
	var count int
	if countStr != nil {
		count, _ = strconv.Atoi(string(countStr))
	}
	
	// Check if limit exceeded
	if count >= rateLimit*60 { // Convert per-second to per-minute
		return false, nil
	}
	
	// Increment counter
	count++
	err = store.Set(ctx, redisKey, []byte(strconv.Itoa(count)), 2*time.Minute)
	
	return true, err
}

// setRateLimitHeaders sets rate limit headers
func setRateLimitHeaders(c echo.Context, rateLimit, burst int) {
	c.Response().Header().Set("X-RateLimit-Limit", strconv.Itoa(rateLimit))
	c.Response().Header().Set("X-RateLimit-Burst", strconv.Itoa(burst))
	c.Response().Header().Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(time.Minute).Unix(), 10))
}

// IPRateLimiter creates an IP-based rate limiter
func IPRateLimiter(rate, burst int) echo.MiddlewareFunc {
	return RateLimiter(RateLimiterConfig{
		Rate:  rate,
		Burst: burst,
		KeyGenerator: func(c echo.Context) string {
			return c.RealIP()
		},
	})
}

// UserRateLimiter creates a user-based rate limiter
func UserRateLimiter(rate, burst int) echo.MiddlewareFunc {
	return RateLimiter(RateLimiterConfig{
		Rate:  rate,
		Burst: burst,
		KeyGenerator: func(c echo.Context) string {
			// Extract user ID from JWT token
			user := c.Get("user")
			if user != nil {
				if userMap, ok := user.(map[string]interface{}); ok {
					if userID, ok := userMap["id"].(string); ok {
						return "user:" + userID
					}
				}
			}
			// Fallback to IP
			return c.RealIP()
		},
	})
}

// EndpointRateLimiter creates an endpoint-specific rate limiter
func EndpointRateLimiter(rate, burst int) echo.MiddlewareFunc {
	return RateLimiter(RateLimiterConfig{
		Rate:  rate,
		Burst: burst,
		KeyGenerator: func(c echo.Context) string {
			// Combine IP and endpoint
			return fmt.Sprintf("%s:%s:%s", c.RealIP(), c.Request().Method, c.Path())
		},
	})
}

// AdaptiveRateLimiter adjusts rate limits based on user behavior
type AdaptiveRateLimiter struct {
	baseRate  int
	baseBurst int
	store     *cache.RedisCache
}

// NewAdaptiveRateLimiter creates a new adaptive rate limiter
func NewAdaptiveRateLimiter(baseRate, baseBurst int, store *cache.RedisCache) *AdaptiveRateLimiter {
	return &AdaptiveRateLimiter{
		baseRate:  baseRate,
		baseBurst: baseBurst,
		store:     store,
	}
}

// Middleware returns the adaptive rate limiter middleware
func (a *AdaptiveRateLimiter) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			key := c.RealIP()
			
			// Get user reputation score
			reputation := a.getUserReputation(c, key)
			
			// Adjust rate limits based on reputation
			rate := a.baseRate
			burst := a.baseBurst
			
			if reputation > 0.8 {
				// Good reputation: increase limits
				rate = int(float64(rate) * 1.5)
				burst = int(float64(burst) * 1.5)
			} else if reputation < 0.3 {
				// Poor reputation: decrease limits
				rate = int(float64(rate) * 0.5)
				burst = int(float64(burst) * 0.5)
			}
			
			// Apply rate limiting
			limiter := getOrCreateLimiter(key, rate, burst)
			if !limiter.Allow() {
				return echo.NewHTTPError(http.StatusTooManyRequests, "rate limit exceeded")
			}
			
			// Update reputation based on behavior
			go a.updateReputation(c, key)
			
			return next(c)
		}
	}
}

// getUserReputation gets the user's reputation score
func (a *AdaptiveRateLimiter) getUserReputation(c echo.Context, key string) float64 {
	ctx := c.Request().Context()
	reputationKey := fmt.Sprintf("reputation:%s", key)
	
	data, err := a.store.Get(ctx, reputationKey)
	if err != nil || data == nil {
		return 0.5 // Default reputation
	}
	
	reputation, _ := strconv.ParseFloat(string(data), 64)
	return reputation
}

// updateReputation updates the user's reputation based on behavior
func (a *AdaptiveRateLimiter) updateReputation(c echo.Context, key string) {
	// Simple reputation algorithm
	// Good behavior: successful requests, low error rate
	// Bad behavior: many 4xx errors, rapid requests
	
	status := c.Response().Status
	reputation := a.getUserReputation(c, key)
	
	if status >= 200 && status < 300 {
		// Successful request: improve reputation
		reputation = reputation*0.9 + 0.1
	} else if status >= 400 && status < 500 {
		// Client error: decrease reputation
		reputation = reputation*0.9 - 0.1
	}
	
	// Clamp reputation between 0 and 1
	if reputation < 0 {
		reputation = 0
	} else if reputation > 1 {
		reputation = 1
	}
	
	// Store updated reputation
	ctx := c.Request().Context()
	reputationKey := fmt.Sprintf("reputation:%s", key)
	a.store.Set(ctx, reputationKey, []byte(fmt.Sprintf("%.2f", reputation)), 24*time.Hour)
}