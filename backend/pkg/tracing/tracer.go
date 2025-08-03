package tracing

import (
	"context"
	"fmt"
	"io"
	
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
)

// Config 追蹤配置
type Config struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	ExporterType   string // jaeger, otlp
	Endpoint       string
	SamplingRate   float64
	Enabled        bool
}

// Tracer 追蹤器包裝
type Tracer struct {
	tracer   trace.Tracer
	provider *sdktrace.TracerProvider
	config   Config
}

// NewTracer 創建新的追蹤器
func NewTracer(config Config) (*Tracer, error) {
	if !config.Enabled {
		return &Tracer{
			tracer: otel.Tracer(config.ServiceName),
			config: config,
		}, nil
	}
	
	// 創建資源
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(config.ServiceName),
			semconv.ServiceVersion(config.ServiceVersion),
			attribute.String("environment", config.Environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}
	
	// 創建導出器
	var exporter sdktrace.SpanExporter
	switch config.ExporterType {
	case "jaeger":
		exporter, err = createJaegerExporter(config)
	case "otlp":
		exporter, err = createOTLPExporter(config)
	default:
		return nil, fmt.Errorf("unsupported exporter type: %s", config.ExporterType)
	}
	
	if err != nil {
		return nil, fmt.Errorf("failed to create exporter: %w", err)
	}
	
	// 創建追蹤提供者
	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(config.SamplingRate)),
	)
	
	// 設置全局追蹤提供者
	otel.SetTracerProvider(provider)
	
	// 設置全局傳播器
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
	
	return &Tracer{
		tracer:   provider.Tracer(config.ServiceName),
		provider: provider,
		config:   config,
	}, nil
}

// Start 開始一個新的 span
func (t *Tracer) Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, spanName, opts...)
}

// StartWithAttributes 開始一個帶屬性的新 span
func (t *Tracer) StartWithAttributes(ctx context.Context, spanName string, attrs map[string]interface{}) (context.Context, trace.Span) {
	var attributes []attribute.KeyValue
	for k, v := range attrs {
		attributes = append(attributes, attributeFromValue(k, v))
	}
	
	return t.tracer.Start(ctx, spanName, trace.WithAttributes(attributes...))
}

// Close 關閉追蹤器
func (t *Tracer) Close(ctx context.Context) error {
	if t.provider != nil {
		return t.provider.Shutdown(ctx)
	}
	return nil
}

// createJaegerExporter 創建 Jaeger 導出器
func createJaegerExporter(config Config) (sdktrace.SpanExporter, error) {
	return jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(config.Endpoint)))
}

// createOTLPExporter 創建 OTLP 導出器
func createOTLPExporter(config Config) (sdktrace.SpanExporter, error) {
	client := otlptracegrpc.NewClient(
		otlptracegrpc.WithEndpoint(config.Endpoint),
		otlptracegrpc.WithInsecure(),
	)
	
	return otlptrace.New(context.Background(), client)
}

// attributeFromValue 從值創建屬性
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

// SpanOption span 選項
type SpanOption func(*spanOptions)

type spanOptions struct {
	attributes map[string]interface{}
	kind       trace.SpanKind
}

// WithAttributes 設置 span 屬性
func WithAttributes(attrs map[string]interface{}) SpanOption {
	return func(o *spanOptions) {
		o.attributes = attrs
	}
}

// WithSpanKind 設置 span 類型
func WithSpanKind(kind trace.SpanKind) SpanOption {
	return func(o *spanOptions) {
		o.kind = kind
	}
}

// StartSpan 開始一個新的 span（輔助函數）
func StartSpan(ctx context.Context, name string, opts ...SpanOption) (context.Context, trace.Span) {
	options := &spanOptions{
		kind: trace.SpanKindInternal,
	}
	
	for _, opt := range opts {
		opt(options)
	}
	
	tracer := otel.Tracer("fastenmind")
	
	var traceOpts []trace.SpanStartOption
	traceOpts = append(traceOpts, trace.WithSpanKind(options.kind))
	
	if options.attributes != nil {
		var attrs []attribute.KeyValue
		for k, v := range options.attributes {
			attrs = append(attrs, attributeFromValue(k, v))
		}
		traceOpts = append(traceOpts, trace.WithAttributes(attrs...))
	}
	
	return tracer.Start(ctx, name, traceOpts...)
}

// RecordError 記錄錯誤到當前 span
func RecordError(ctx context.Context, err error, description string) {
	span := trace.SpanFromContext(ctx)
	if span != nil && err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String("error.description", description),
		))
	}
}

// SetStatus 設置 span 狀態
func SetStatus(ctx context.Context, code codes.Code, description string) {
	span := trace.SpanFromContext(ctx)
	if span != nil {
		span.SetStatus(code, description)
	}
}

// AddEvent 添加事件到當前 span
func AddEvent(ctx context.Context, name string, attrs map[string]interface{}) {
	span := trace.SpanFromContext(ctx)
	if span != nil {
		var attributes []attribute.KeyValue
		for k, v := range attrs {
			attributes = append(attributes, attributeFromValue(k, v))
		}
		span.AddEvent(name, trace.WithAttributes(attributes...))
	}
}

// ExtractTraceID 從上下文中提取追蹤 ID
func ExtractTraceID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if span != nil {
		return span.SpanContext().TraceID().String()
	}
	return ""
}

// InjectTraceContext 注入追蹤上下文到載體
func InjectTraceContext(ctx context.Context, carrier propagation.TextMapCarrier) {
	otel.GetTextMapPropagator().Inject(ctx, carrier)
}

// ExtractTraceContext 從載體中提取追蹤上下文
func ExtractTraceContext(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
	return otel.GetTextMapPropagator().Extract(ctx, carrier)
}

// TracerProvider 追蹤提供者介面
type TracerProvider interface {
	Tracer(instrumentationName string, opts ...trace.TracerOption) trace.Tracer
	Shutdown(ctx context.Context) error
}

// NoopTracerProvider 無操作追蹤提供者
type NoopTracerProvider struct{}

func (n NoopTracerProvider) Tracer(instrumentationName string, opts ...trace.TracerOption) trace.Tracer {
	return otel.Tracer(instrumentationName, opts...)
}

func (n NoopTracerProvider) Shutdown(ctx context.Context) error {
	return nil
}

// SpanProcessor 自定義 span 處理器
type SpanProcessor struct {
	exporter sdktrace.SpanExporter
}

// NewSpanProcessor 創建新的 span 處理器
func NewSpanProcessor(exporter sdktrace.SpanExporter) *SpanProcessor {
	return &SpanProcessor{
		exporter: exporter,
	}
}

// OnStart 在 span 開始時調用
func (p *SpanProcessor) OnStart(parent context.Context, s sdktrace.ReadWriteSpan) {}

// OnEnd 在 span 結束時調用
func (p *SpanProcessor) OnEnd(s sdktrace.ReadOnlySpan) {
	// 可以在這裡添加自定義邏輯
}

// Shutdown 關閉處理器
func (p *SpanProcessor) Shutdown(ctx context.Context) error {
	if closer, ok := p.exporter.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

// ForceFlush 強制刷新
func (p *SpanProcessor) ForceFlush(ctx context.Context) error {
	return nil
}