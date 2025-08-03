package tracing

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Config holds tracing configuration
type Config struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	
	// Exporter configuration
	ExporterType string // "jaeger", "otlp"
	Endpoint     string
	
	// Sampling
	SampleRate float64
	
	// Resource attributes
	Attributes map[string]string
}

// InitTracer initializes the OpenTelemetry tracer
func InitTracer(ctx context.Context, config Config) (func(context.Context) error, error) {
	// Create resource
	res, err := createResource(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}
	
	// Create exporter
	var exporter sdktrace.SpanExporter
	switch config.ExporterType {
	case "jaeger":
		exporter, err = createJaegerExporter(config)
	case "otlp":
		exporter, err = createOTLPExporter(ctx, config)
	default:
		return nil, fmt.Errorf("unsupported exporter type: %s", config.ExporterType)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create exporter: %w", err)
	}
	
	// Create tracer provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(config.SampleRate)),
	)
	
	// Set global tracer provider
	otel.SetTracerProvider(tp)
	
	// Set global propagator
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
	
	// Return shutdown function
	return tp.Shutdown, nil
}

// createResource creates the resource for the tracer
func createResource(config Config) (*resource.Resource, error) {
	// Default attributes
	attrs := []attribute.KeyValue{
		semconv.ServiceNameKey.String(config.ServiceName),
		semconv.ServiceVersionKey.String(config.ServiceVersion),
		semconv.DeploymentEnvironmentKey.String(config.Environment),
	}
	
	// Add custom attributes
	for k, v := range config.Attributes {
		attrs = append(attrs, attribute.String(k, v))
	}
	
	return resource.New(context.Background(),
		resource.WithAttributes(attrs...),
		resource.WithHost(),
		resource.WithContainer(),
	)
}

// createJaegerExporter creates a Jaeger exporter
func createJaegerExporter(config Config) (sdktrace.SpanExporter, error) {
	return jaeger.New(jaeger.WithCollectorEndpoint(
		jaeger.WithEndpoint(config.Endpoint),
	))
}

// createOTLPExporter creates an OTLP exporter
func createOTLPExporter(ctx context.Context, config Config) (sdktrace.SpanExporter, error) {
	conn, err := grpc.DialContext(ctx, config.Endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection: %w", err)
	}
	
	return otlptrace.New(ctx, otlptracegrpc.NewClient(
		otlptracegrpc.WithGRPCConn(conn),
	))
}

// Tracer is a wrapper around OpenTelemetry tracer
type Tracer struct {
	tracer trace.Tracer
}

// NewTracer creates a new tracer
func NewTracer(name string) *Tracer {
	return &Tracer{
		tracer: otel.Tracer(name),
	}
}

// Start starts a new span
func (t *Tracer) Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, spanName, opts...)
}

// StartWithAttributes starts a new span with attributes
func (t *Tracer) StartWithAttributes(ctx context.Context, spanName string, attrs map[string]interface{}) (context.Context, trace.Span) {
	var attributes []attribute.KeyValue
	for k, v := range attrs {
		attributes = append(attributes, attributeFromValue(k, v))
	}
	
	return t.tracer.Start(ctx, spanName, trace.WithAttributes(attributes...))
}

// WithSpan executes a function within a span
func (t *Tracer) WithSpan(ctx context.Context, spanName string, fn func(context.Context) error) error {
	ctx, span := t.Start(ctx, spanName)
	defer span.End()
	
	err := fn(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	
	return err
}

// attributeFromValue creates an attribute from a value
func attributeFromValue(key string, value interface{}) attribute.KeyValue {
	switch v := value.(type) {
	case string:
		return attribute.String(key, v)
	case int:
		return attribute.Int(key, v)
	case int64:
		return attribute.Int64(key, v)
	case float64:
		return attribute.Float64(key, v)
	case bool:
		return attribute.Bool(key, v)
	default:
		return attribute.String(key, fmt.Sprintf("%v", v))
	}
}

// SpanFromContext returns the span from context
func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

// ContextWithSpan returns a new context with the span
func ContextWithSpan(ctx context.Context, span trace.Span) context.Context {
	return trace.ContextWithSpan(ctx, span)
}

// RecordError records an error on the current span
func RecordError(ctx context.Context, err error, opts ...trace.EventOption) {
	span := SpanFromContext(ctx)
	span.RecordError(err, opts...)
	span.SetStatus(codes.Error, err.Error())
}

// AddEvent adds an event to the current span
func AddEvent(ctx context.Context, name string, attrs ...attribute.KeyValue) {
	span := SpanFromContext(ctx)
	span.AddEvent(name, trace.WithAttributes(attrs...))
}

// SetAttributes sets attributes on the current span
func SetAttributes(ctx context.Context, attrs ...attribute.KeyValue) {
	span := SpanFromContext(ctx)
	span.SetAttributes(attrs...)
}

// SetStatus sets the status of the current span
func SetStatus(ctx context.Context, code codes.Code, description string) {
	span := SpanFromContext(ctx)
	span.SetStatus(code, description)
}

// TraceID returns the trace ID from the context
func TraceID(ctx context.Context) string {
	span := SpanFromContext(ctx)
	if span.SpanContext().HasTraceID() {
		return span.SpanContext().TraceID().String()
	}
	return ""
}

// SpanID returns the span ID from the context
func SpanID(ctx context.Context) string {
	span := SpanFromContext(ctx)
	if span.SpanContext().HasSpanID() {
		return span.SpanContext().SpanID().String()
	}
	return ""
}

// HTTPCarrier adapts http.Header to propagation.TextMapCarrier
type HTTPCarrier struct {
	Headers map[string][]string
}

// Get returns the value associated with the passed key
func (c HTTPCarrier) Get(key string) string {
	vals := c.Headers[key]
	if len(vals) > 0 {
		return vals[0]
	}
	return ""
}

// Set stores the key-value pair
func (c HTTPCarrier) Set(key string, value string) {
	c.Headers[key] = []string{value}
}

// Keys lists the keys stored in this carrier
func (c HTTPCarrier) Keys() []string {
	keys := make([]string, 0, len(c.Headers))
	for k := range c.Headers {
		keys = append(keys, k)
	}
	return keys
}

// Extract extracts trace context from carrier
func Extract(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
	return otel.GetTextMapPropagator().Extract(ctx, carrier)
}

// Inject injects trace context into carrier
func Inject(ctx context.Context, carrier propagation.TextMapCarrier) {
	otel.GetTextMapPropagator().Inject(ctx, carrier)
}

// Metrics for tracing
type TracingMetrics struct {
	SpansStarted   int64
	SpansFinished  int64
	SpansErrored   int64
	TraceDuration  time.Duration
}

// MetricsCollector collects tracing metrics
type MetricsCollector struct {
	metrics map[string]*TracingMetrics
	mu      sync.RWMutex
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		metrics: make(map[string]*TracingMetrics),
	}
}

// RecordSpan records span metrics
func (c *MetricsCollector) RecordSpan(operation string, duration time.Duration, hasError bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	if _, exists := c.metrics[operation]; !exists {
		c.metrics[operation] = &TracingMetrics{}
	}
	
	c.metrics[operation].SpansStarted++
	c.metrics[operation].SpansFinished++
	c.metrics[operation].TraceDuration += duration
	
	if hasError {
		c.metrics[operation].SpansErrored++
	}
}

// GetMetrics returns metrics for an operation
func (c *MetricsCollector) GetMetrics(operation string) *TracingMetrics {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	if metrics, exists := c.metrics[operation]; exists {
		return metrics
	}
	return nil
}

// GetAllMetrics returns all metrics
func (c *MetricsCollector) GetAllMetrics() map[string]*TracingMetrics {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	result := make(map[string]*TracingMetrics)
	for k, v := range c.metrics {
		result[k] = v
	}
	return result
}