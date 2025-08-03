package middleware

import (
	"fmt"
	"time"

	"github.com/fastenmind/fastener-api/internal/infrastructure/tracing"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// TracingConfig defines the tracing middleware configuration
type TracingConfig struct {
	// Tracer to use
	Tracer *tracing.Tracer
	
	// Skipper defines a function to skip middleware
	Skipper func(c echo.Context) bool
	
	// SpanNameFormatter formats the span name
	SpanNameFormatter func(c echo.Context) string
	
	// AttributesExtractor extracts additional attributes
	AttributesExtractor func(c echo.Context) []attribute.KeyValue
}

// DefaultTracingConfig returns default tracing configuration
var DefaultTracingConfig = TracingConfig{
	Skipper: func(c echo.Context) bool {
		// Skip health check endpoints
		return c.Path() == "/health" || c.Path() == "/metrics"
	},
	SpanNameFormatter: func(c echo.Context) string {
		return fmt.Sprintf("%s %s", c.Request().Method, c.Path())
	},
	AttributesExtractor: func(c echo.Context) []attribute.KeyValue {
		return nil
	},
}

// Tracing returns a tracing middleware
func Tracing(config TracingConfig) echo.MiddlewareFunc {
	// Apply defaults
	if config.Tracer == nil {
		config.Tracer = tracing.NewTracer("echo-server")
	}
	if config.Skipper == nil {
		config.Skipper = DefaultTracingConfig.Skipper
	}
	if config.SpanNameFormatter == nil {
		config.SpanNameFormatter = DefaultTracingConfig.SpanNameFormatter
	}
	if config.AttributesExtractor == nil {
		config.AttributesExtractor = DefaultTracingConfig.AttributesExtractor
	}
	
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}
			
			// Extract trace context from headers
			ctx := c.Request().Context()
			carrier := tracing.HTTPCarrier{Headers: c.Request().Header}
			ctx = tracing.Extract(ctx, carrier)
			
			// Start span
			spanName := config.SpanNameFormatter(c)
			ctx, span := config.Tracer.Start(ctx, spanName,
				oteltrace.WithSpanKind(oteltrace.SpanKindServer),
			)
			defer span.End()
			
			// Set standard attributes
			attrs := []attribute.KeyValue{
				semconv.HTTPMethodKey.String(c.Request().Method),
				semconv.HTTPTargetKey.String(c.Request().URL.Path),
				semconv.HTTPRouteKey.String(c.Path()),
				semconv.HTTPSchemeKey.String(c.Scheme()),
				semconv.HTTPHostKey.String(c.Request().Host),
				semconv.HTTPUserAgentKey.String(c.Request().UserAgent()),
				semconv.HTTPRequestContentLengthKey.Int64(c.Request().ContentLength),
				semconv.NetPeerIPKey.String(c.RealIP()),
			}
			
			// Add custom attributes
			if customAttrs := config.AttributesExtractor(c); customAttrs != nil {
				attrs = append(attrs, customAttrs...)
			}
			
			span.SetAttributes(attrs...)
			
			// Update request context
			c.SetRequest(c.Request().WithContext(ctx))
			
			// Add trace ID to response headers
			if span.SpanContext().HasTraceID() {
				c.Response().Header().Set("X-Trace-ID", span.SpanContext().TraceID().String())
			}
			
			// Execute handler
			start := time.Now()
			err := next(c)
			duration := time.Since(start)
			
			// Set response attributes
			span.SetAttributes(
				semconv.HTTPStatusCodeKey.Int(c.Response().Status),
				attribute.Int64("http.response.size", c.Response().Size),
				attribute.Int64("duration.ms", duration.Milliseconds()),
			)
			
			// Set span status based on HTTP status
			if err != nil {
				span.RecordError(err)
				span.SetStatus(oteltrace.Status{
					Code:        oteltrace.StatusError,
					Description: err.Error(),
				})
			} else if c.Response().Status >= 400 {
				span.SetStatus(oteltrace.Status{
					Code:        oteltrace.StatusError,
					Description: fmt.Sprintf("HTTP %d", c.Response().Status),
				})
			}
			
			return err
		}
	}
}

// DatabaseTracing wraps database operations with tracing
func DatabaseTracing(tracer *tracing.Tracer) func(operationName string, fn func() error) error {
	return func(operationName string, fn func() error) error {
		// This would typically be integrated with GORM hooks
		// For demonstration purposes
		return tracer.WithSpan(nil, fmt.Sprintf("db.%s", operationName), func(ctx context.Context) error {
			tracing.SetAttributes(ctx,
				attribute.String("db.system", "postgresql"),
				attribute.String("db.operation", operationName),
			)
			return fn()
		})
	}
}

// CacheTracing wraps cache operations with tracing
func CacheTracing(tracer *tracing.Tracer) func(operationName string, key string, fn func() error) error {
	return func(operationName string, key string, fn func() error) error {
		return tracer.WithSpan(nil, fmt.Sprintf("cache.%s", operationName), func(ctx context.Context) error {
			tracing.SetAttributes(ctx,
				attribute.String("cache.system", "redis"),
				attribute.String("cache.operation", operationName),
				attribute.String("cache.key", key),
			)
			
			start := time.Now()
			err := fn()
			duration := time.Since(start)
			
			tracing.SetAttributes(ctx,
				attribute.Int64("cache.duration.ms", duration.Milliseconds()),
				attribute.Bool("cache.hit", err == nil),
			)
			
			return err
		})
	}
}

// HTTPClientTracing adds tracing to HTTP client requests
func HTTPClientTracing(tracer *tracing.Tracer) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// This would be used for outgoing HTTP requests
			// Inject trace context into outgoing headers
			carrier := tracing.HTTPCarrier{Headers: make(map[string][]string)}
			tracing.Inject(c.Request().Context(), carrier)
			
			// Add headers to outgoing request
			for k, v := range carrier.Headers {
				if len(v) > 0 {
					c.Request().Header.Set(k, v[0])
				}
			}
			
			return next(c)
		}
	}
}

// AsyncOperationTracing traces asynchronous operations
func AsyncOperationTracing(tracer *tracing.Tracer, operationName string, fn func(ctx context.Context) error) {
	go func() {
		ctx, span := tracer.Start(context.Background(), operationName,
			oteltrace.WithSpanKind(oteltrace.SpanKindInternal),
		)
		defer span.End()
		
		if err := fn(ctx); err != nil {
			span.RecordError(err)
			span.SetStatus(oteltrace.Status{
				Code:        oteltrace.StatusError,
				Description: err.Error(),
			})
		}
	}()
}

// TraceGRPCClient adds tracing to gRPC client calls
func TraceGRPCClient(ctx context.Context, tracer *tracing.Tracer, method string, fn func(context.Context) error) error {
	ctx, span := tracer.Start(ctx, method,
		oteltrace.WithSpanKind(oteltrace.SpanKindClient),
		oteltrace.WithAttributes(
			attribute.String("rpc.system", "grpc"),
			attribute.String("rpc.method", method),
		),
	)
	defer span.End()
	
	err := fn(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(oteltrace.Status{
			Code:        oteltrace.StatusError,
			Description: err.Error(),
		})
	}
	
	return err
}

// TracingStats provides middleware to collect tracing statistics
type TracingStats struct {
	collector *tracing.MetricsCollector
}

// NewTracingStats creates a new tracing stats middleware
func NewTracingStats() *TracingStats {
	return &TracingStats{
		collector: tracing.NewMetricsCollector(),
	}
}

// Middleware returns the tracing stats middleware
func (ts *TracingStats) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			operation := fmt.Sprintf("%s %s", c.Request().Method, c.Path())
			start := time.Now()
			
			err := next(c)
			
			duration := time.Since(start)
			ts.collector.RecordSpan(operation, duration, err != nil || c.Response().Status >= 400)
			
			return err
		}
	}
}

// GetStats returns the collected statistics
func (ts *TracingStats) GetStats() map[string]*tracing.TracingMetrics {
	return ts.collector.GetAllMetrics()
}