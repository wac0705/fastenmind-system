package gateway

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/fastenmind/fastener-api/internal/middleware"
	"github.com/labstack/echo/v4"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
)

// APIGateway provides API gateway functionality
type APIGateway struct {
	routes          map[string]*Route
	services        map[string]*Service
	loadBalancers   map[string]LoadBalancer
	circuitBreakers map[string]*gobreaker.CircuitBreaker
	rateLimiters    map[string]*rate.Limiter
	mu              sync.RWMutex
	config          *Config
}

// Config holds API gateway configuration
type Config struct {
	// Gateway settings
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	MaxHeaderBytes  int
	
	// Security settings
	EnableCORS      bool
	EnableCSRF      bool
	EnableRateLimit bool
	RateLimitRPS    int
	
	// Circuit breaker settings
	CBMaxRequests   uint32
	CBInterval      time.Duration
	CBTimeout       time.Duration
	CBFailureRatio  float64
	
	// Retry settings
	RetryAttempts   int
	RetryDelay      time.Duration
	
	// Logging
	EnableAccessLog bool
	EnableDebugLog  bool
}

// Route represents an API route
type Route struct {
	Path            string
	Service         string
	Methods         []string
	Rewrite         string
	StripPrefix     bool
	Authentication  bool
	RateLimitRPS    int
	Timeout         time.Duration
	CircuitBreaker  bool
	CacheEnabled    bool
	CacheDuration   time.Duration
	Transformations []Transformation
}

// Service represents a backend service
type Service struct {
	Name      string
	Endpoints []string
	Health    string
	Timeout   time.Duration
	Retries   int
}

// Transformation represents a request/response transformation
type Transformation struct {
	Type      string // "request" or "response"
	Operation string // "add_header", "remove_header", "modify_body"
	Config    map[string]interface{}
}

// NewAPIGateway creates a new API gateway
func NewAPIGateway(config *Config) *APIGateway {
	return &APIGateway{
		routes:          make(map[string]*Route),
		services:        make(map[string]*Service),
		loadBalancers:   make(map[string]LoadBalancer),
		circuitBreakers: make(map[string]*gobreaker.CircuitBreaker),
		rateLimiters:    make(map[string]*rate.Limiter),
		config:          config,
	}
}

// RegisterRoute registers a new route
func (g *APIGateway) RegisterRoute(route *Route) {
	g.mu.Lock()
	defer g.mu.Unlock()
	
	g.routes[route.Path] = route
	
	// Initialize rate limiter if needed
	if route.RateLimitRPS > 0 {
		g.rateLimiters[route.Path] = rate.NewLimiter(rate.Limit(route.RateLimitRPS), route.RateLimitRPS)
	}
	
	// Initialize circuit breaker if needed
	if route.CircuitBreaker {
		g.circuitBreakers[route.Path] = gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:        route.Path,
			MaxRequests: g.config.CBMaxRequests,
			Interval:    g.config.CBInterval,
			Timeout:     g.config.CBTimeout,
			ReadyToTrip: func(counts gobreaker.Counts) bool {
				failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
				return counts.Requests >= 3 && failureRatio >= g.config.CBFailureRatio
			},
		})
	}
}

// RegisterService registers a backend service
func (g *APIGateway) RegisterService(service *Service) {
	g.mu.Lock()
	defer g.mu.Unlock()
	
	g.services[service.Name] = service
	g.loadBalancers[service.Name] = NewRoundRobinLoadBalancer(service.Endpoints)
}

// SetupRoutes sets up Echo routes
func (g *APIGateway) SetupRoutes(e *echo.Echo) {
	// Add middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	
	if g.config.EnableCORS {
		e.Use(middleware.CORS())
	}
	
	// Setup gateway routes
	for path, route := range g.routes {
		for _, method := range route.Methods {
			e.Add(method, path, g.handleRequest(route))
		}
	}
	
	// Health check endpoint
	e.GET("/health", g.healthCheck)
	
	// Metrics endpoint
	e.GET("/metrics", g.metrics)
}

// handleRequest handles incoming requests
func (g *APIGateway) handleRequest(route *Route) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Rate limiting
		if route.RateLimitRPS > 0 {
			limiter := g.rateLimiters[route.Path]
			if !limiter.Allow() {
				return c.JSON(http.StatusTooManyRequests, map[string]string{
					"error": "Rate limit exceeded",
				})
			}
		}
		
		// Get service
		service, exists := g.services[route.Service]
		if !exists {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{
				"error": "Service not available",
			})
		}
		
		// Get backend endpoint
		lb := g.loadBalancers[route.Service]
		endpoint := lb.Next()
		if endpoint == "" {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{
				"error": "No healthy endpoints available",
			})
		}
		
		// Apply request transformations
		g.applyTransformations(c, route.Transformations, "request")
		
		// Circuit breaker
		if route.CircuitBreaker {
			cb := g.circuitBreakers[route.Path]
			result, err := cb.Execute(func() (interface{}, error) {
				return g.proxyRequest(c, endpoint, route)
			})
			
			if err != nil {
				if err == gobreaker.ErrOpenState {
					return c.JSON(http.StatusServiceUnavailable, map[string]string{
						"error": "Circuit breaker is open",
					})
				}
				return err
			}
			
			return result.(error)
		}
		
		// Proxy request
		return g.proxyRequest(c, endpoint, route)
	}
}

// proxyRequest proxies the request to backend service
func (g *APIGateway) proxyRequest(c echo.Context, endpoint string, route *Route) error {
	// Parse backend URL
	target, err := url.Parse(endpoint)
	if err != nil {
		return err
	}
	
	// Create reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(target)
	
	// Customize director
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		
		// Rewrite path if needed
		if route.Rewrite != "" {
			req.URL.Path = route.Rewrite
		} else if route.StripPrefix {
			req.URL.Path = strings.TrimPrefix(req.URL.Path, route.Path)
		}
		
		// Add gateway headers
		req.Header.Set("X-Gateway-Route", route.Path)
		req.Header.Set("X-Forwarded-For", c.RealIP())
		req.Header.Set("X-Request-ID", c.Request().Header.Get("X-Request-ID"))
	}
	
	// Customize error handler
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		c.JSON(http.StatusBadGateway, map[string]string{
			"error": "Backend service error",
		})
	}
	
	// Set timeout
	ctx := c.Request().Context()
	if route.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, route.Timeout)
		defer cancel()
		c.Request() = c.Request().WithContext(ctx)
	}
	
	// Execute proxy
	proxy.ServeHTTP(c.Response(), c.Request())
	
	// Apply response transformations
	g.applyTransformations(c, route.Transformations, "response")
	
	return nil
}

// applyTransformations applies request/response transformations
func (g *APIGateway) applyTransformations(c echo.Context, transformations []Transformation, transformType string) {
	for _, t := range transformations {
		if t.Type != transformType {
			continue
		}
		
		switch t.Operation {
		case "add_header":
			key := t.Config["key"].(string)
			value := t.Config["value"].(string)
			if transformType == "request" {
				c.Request().Header.Set(key, value)
			} else {
				c.Response().Header().Set(key, value)
			}
			
		case "remove_header":
			key := t.Config["key"].(string)
			if transformType == "request" {
				c.Request().Header.Del(key)
			} else {
				c.Response().Header().Del(key)
			}
			
		case "modify_body":
			// Implementation for body modification
			// This would require buffering the body
		}
	}
}

// healthCheck performs health check on all services
func (g *APIGateway) healthCheck(c echo.Context) error {
	g.mu.RLock()
	defer g.mu.RUnlock()
	
	health := make(map[string]interface{})
	health["gateway"] = "healthy"
	health["timestamp"] = time.Now()
	
	services := make(map[string]string)
	for name, service := range g.services {
		lb := g.loadBalancers[name]
		if lb.HealthyEndpoints() > 0 {
			services[name] = "healthy"
		} else {
			services[name] = "unhealthy"
		}
	}
	health["services"] = services
	
	return c.JSON(http.StatusOK, health)
}

// metrics returns gateway metrics
func (g *APIGateway) metrics(c echo.Context) error {
	metrics := make(map[string]interface{})
	
	// Circuit breaker states
	cbStates := make(map[string]string)
	for path, cb := range g.circuitBreakers {
		state := "closed"
		if cb.State() == gobreaker.StateOpen {
			state = "open"
		} else if cb.State() == gobreaker.StateHalfOpen {
			state = "half-open"
		}
		cbStates[path] = state
	}
	metrics["circuit_breakers"] = cbStates
	
	// Load balancer stats
	lbStats := make(map[string]interface{})
	for name, lb := range g.loadBalancers {
		lbStats[name] = map[string]interface{}{
			"total_endpoints":   lb.TotalEndpoints(),
			"healthy_endpoints": lb.HealthyEndpoints(),
		}
	}
	metrics["load_balancers"] = lbStats
	
	return c.JSON(http.StatusOK, metrics)
}

// LoadBalancer interface
type LoadBalancer interface {
	Next() string
	MarkHealthy(endpoint string)
	MarkUnhealthy(endpoint string)
	TotalEndpoints() int
	HealthyEndpoints() int
}

// RoundRobinLoadBalancer implements round-robin load balancing
type RoundRobinLoadBalancer struct {
	endpoints []string
	current   int
	healthy   map[string]bool
	mu        sync.Mutex
}

// NewRoundRobinLoadBalancer creates a new round-robin load balancer
func NewRoundRobinLoadBalancer(endpoints []string) *RoundRobinLoadBalancer {
	healthy := make(map[string]bool)
	for _, ep := range endpoints {
		healthy[ep] = true
	}
	
	return &RoundRobinLoadBalancer{
		endpoints: endpoints,
		healthy:   healthy,
	}
}

// Next returns the next healthy endpoint
func (lb *RoundRobinLoadBalancer) Next() string {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	
	if len(lb.endpoints) == 0 {
		return ""
	}
	
	// Find next healthy endpoint
	start := lb.current
	for {
		lb.current = (lb.current + 1) % len(lb.endpoints)
		endpoint := lb.endpoints[lb.current]
		
		if lb.healthy[endpoint] {
			return endpoint
		}
		
		if lb.current == start {
			// No healthy endpoints
			return ""
		}
	}
}

// MarkHealthy marks an endpoint as healthy
func (lb *RoundRobinLoadBalancer) MarkHealthy(endpoint string) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	lb.healthy[endpoint] = true
}

// MarkUnhealthy marks an endpoint as unhealthy
func (lb *RoundRobinLoadBalancer) MarkUnhealthy(endpoint string) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	lb.healthy[endpoint] = false
}

// TotalEndpoints returns total number of endpoints
func (lb *RoundRobinLoadBalancer) TotalEndpoints() int {
	return len(lb.endpoints)
}

// HealthyEndpoints returns number of healthy endpoints
func (lb *RoundRobinLoadBalancer) HealthyEndpoints() int {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	
	count := 0
	for _, healthy := range lb.healthy {
		if healthy {
			count++
		}
	}
	return count
}

// RequestTransformer provides request transformation capabilities
type RequestTransformer struct {
	rules []TransformRule
}

type TransformRule struct {
	Match     MatchCondition
	Transform TransformAction
}

type MatchCondition struct {
	Path    string
	Method  string
	Headers map[string]string
}

type TransformAction struct {
	AddHeaders    map[string]string
	RemoveHeaders []string
	RewritePath   string
	AddQuery      map[string]string
}

// ResponseCache provides response caching
type ResponseCache struct {
	cache sync.Map
	ttl   time.Duration
}

type CachedResponse struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
	Timestamp  time.Time
}

// NewResponseCache creates a new response cache
func NewResponseCache(ttl time.Duration) *ResponseCache {
	rc := &ResponseCache{ttl: ttl}
	go rc.cleanup()
	return rc
}

func (rc *ResponseCache) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		now := time.Now()
		rc.cache.Range(func(key, value interface{}) bool {
			cached := value.(*CachedResponse)
			if now.Sub(cached.Timestamp) > rc.ttl {
				rc.cache.Delete(key)
			}
			return true
		})
	}
}