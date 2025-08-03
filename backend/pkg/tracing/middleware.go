package tracing

import (
	"context"
	"fmt"
	"net/http"
	"time"
	
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

// EchoMiddleware Echo 框架的追蹤中間件
func EchoMiddleware(serviceName string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tracer := otel.Tracer(serviceName)
			
			// 從請求頭中提取追蹤上下文
			ctx := otel.GetTextMapPropagator().Extract(
				c.Request().Context(),
				propagation.HeaderCarrier(c.Request().Header),
			)
			
			// 開始新的 span
			spanName := fmt.Sprintf("%s %s", c.Request().Method, c.Path())
			ctx, span := tracer.Start(ctx, spanName,
				trace.WithSpanKind(trace.SpanKindServer),
				trace.WithAttributes(
					attribute.String("http.method", c.Request().Method),
					attribute.String("http.url", c.Request().URL.String()),
					attribute.String("http.target", c.Path()),
					attribute.String("http.host", c.Request().Host),
					attribute.String("http.scheme", c.Scheme()),
					attribute.String("http.user_agent", c.Request().UserAgent()),
					attribute.String("http.remote_addr", c.Request().RemoteAddr),
				),
			)
			defer span.End()
			
			// 設置請求上下文
			c.SetRequest(c.Request().WithContext(ctx))
			
			// 處理請求
			err := next(c)
			
			// 記錄響應信息
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				
				// 如果是 HTTP 錯誤，提取狀態碼
				if he, ok := err.(*echo.HTTPError); ok {
					span.SetAttributes(attribute.Int("http.status_code", he.Code))
				}
			} else {
				span.SetAttributes(attribute.Int("http.status_code", c.Response().Status))
				if c.Response().Status >= 400 {
					span.SetStatus(codes.Error, http.StatusText(c.Response().Status))
				} else {
					span.SetStatus(codes.Ok, "")
				}
			}
			
			// 注入追蹤上下文到響應頭
			otel.GetTextMapPropagator().Inject(
				ctx,
				propagation.HeaderCarrier(c.Response().Header()),
			)
			
			return err
		}
	}
}

// GormPlugin GORM 的追蹤插件
type GormPlugin struct {
	serviceName string
}

// NewGormPlugin 創建新的 GORM 追蹤插件
func NewGormPlugin(serviceName string) *GormPlugin {
	return &GormPlugin{
		serviceName: serviceName,
	}
}

// Name 插件名稱
func (p *GormPlugin) Name() string {
	return "opentelemetry"
}

// Initialize 初始化插件
func (p *GormPlugin) Initialize(db *gorm.DB) error {
	// 註冊回調
	p.registerCallbacks(db)
	return nil
}

// registerCallbacks 註冊 GORM 回調
func (p *GormPlugin) registerCallbacks(db *gorm.DB) {
	// Create 回調
	db.Callback().Create().Before("gorm:create").Register("otel:before_create", p.beforeCreate)
	db.Callback().Create().After("gorm:create").Register("otel:after_create", p.afterCreate)
	
	// Query 回調
	db.Callback().Query().Before("gorm:query").Register("otel:before_query", p.beforeQuery)
	db.Callback().Query().After("gorm:query").Register("otel:after_query", p.afterQuery)
	
	// Update 回調
	db.Callback().Update().Before("gorm:update").Register("otel:before_update", p.beforeUpdate)
	db.Callback().Update().After("gorm:update").Register("otel:after_update", p.afterUpdate)
	
	// Delete 回調
	db.Callback().Delete().Before("gorm:delete").Register("otel:before_delete", p.beforeDelete)
	db.Callback().Delete().After("gorm:delete").Register("otel:after_delete", p.afterDelete)
}

// beforeCreate 創建前回調
func (p *GormPlugin) beforeCreate(db *gorm.DB) {
	p.startSpan(db, "gorm.create")
}

// afterCreate 創建後回調
func (p *GormPlugin) afterCreate(db *gorm.DB) {
	p.endSpan(db)
}

// beforeQuery 查詢前回調
func (p *GormPlugin) beforeQuery(db *gorm.DB) {
	p.startSpan(db, "gorm.query")
}

// afterQuery 查詢後回調
func (p *GormPlugin) afterQuery(db *gorm.DB) {
	p.endSpan(db)
}

// beforeUpdate 更新前回調
func (p *GormPlugin) beforeUpdate(db *gorm.DB) {
	p.startSpan(db, "gorm.update")
}

// afterUpdate 更新後回調
func (p *GormPlugin) afterUpdate(db *gorm.DB) {
	p.endSpan(db)
}

// beforeDelete 刪除前回調
func (p *GormPlugin) beforeDelete(db *gorm.DB) {
	p.startSpan(db, "gorm.delete")
}

// afterDelete 刪除後回調
func (p *GormPlugin) afterDelete(db *gorm.DB) {
	p.endSpan(db)
}

// startSpan 開始 span
func (p *GormPlugin) startSpan(db *gorm.DB, operation string) {
	tracer := otel.Tracer(p.serviceName)
	
	ctx := db.Statement.Context
	if ctx == nil {
		return
	}
	
	// 開始新的 span
	ctx, span := tracer.Start(ctx, operation,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String("db.system", "mysql"),
			attribute.String("db.operation", operation),
			attribute.String("db.statement", db.Statement.SQL.String()),
			attribute.String("db.table", db.Statement.Table),
		),
	)
	
	// 存儲 span 和開始時間
	db.Set("otel:span", span)
	db.Set("otel:start_time", time.Now())
	db.Statement.Context = ctx
}

// endSpan 結束 span
func (p *GormPlugin) endSpan(db *gorm.DB) {
	// 獲取 span
	val, ok := db.Get("otel:span")
	if !ok {
		return
	}
	
	span, ok := val.(trace.Span)
	if !ok {
		return
	}
	
	// 獲取開始時間
	if startTime, ok := db.Get("otel:start_time"); ok {
		if st, ok := startTime.(time.Time); ok {
			span.SetAttributes(attribute.Int64("db.duration_ms", time.Since(st).Milliseconds()))
		}
	}
	
	// 記錄錯誤
	if db.Error != nil && db.Error != gorm.ErrRecordNotFound {
		span.RecordError(db.Error)
		span.SetStatus(codes.Error, db.Error.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}
	
	// 記錄受影響的行數
	span.SetAttributes(attribute.Int64("db.rows_affected", db.RowsAffected))
	
	// 結束 span
	span.End()
}

// HTTPClient 帶追蹤的 HTTP 客戶端
type HTTPClient struct {
	client      *http.Client
	serviceName string
}

// NewHTTPClient 創建帶追蹤的 HTTP 客戶端
func NewHTTPClient(serviceName string) *HTTPClient {
	return &HTTPClient{
		client:      &http.Client{},
		serviceName: serviceName,
	}
}

// Do 執行 HTTP 請求
func (c *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	tracer := otel.Tracer(c.serviceName)
	
	// 開始新的 span
	ctx, span := tracer.Start(req.Context(), fmt.Sprintf("HTTP %s", req.Method),
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String("http.method", req.Method),
			attribute.String("http.url", req.URL.String()),
			attribute.String("http.host", req.Host),
		),
	)
	defer span.End()
	
	// 注入追蹤上下文到請求頭
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	
	// 執行請求
	req = req.WithContext(ctx)
	resp, err := c.client.Do(req)
	
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return resp, err
	}
	
	// 記錄響應信息
	span.SetAttributes(attribute.Int("http.status_code", resp.StatusCode))
	
	if resp.StatusCode >= 400 {
		span.SetStatus(codes.Error, http.StatusText(resp.StatusCode))
	} else {
		span.SetStatus(codes.Ok, "")
	}
	
	return resp, nil
}

// MessageMiddleware 訊息中間件追蹤
func MessageMiddleware(serviceName string) func(next func(ctx context.Context, message interface{}) error) func(ctx context.Context, message interface{}) error {
	return func(next func(ctx context.Context, message interface{}) error) func(ctx context.Context, message interface{}) error {
		return func(ctx context.Context, message interface{}) error {
			tracer := otel.Tracer(serviceName)
			
			// 提取訊息類型
			messageType := fmt.Sprintf("%T", message)
			
			// 開始新的 span
			ctx, span := tracer.Start(ctx, fmt.Sprintf("Message %s", messageType),
				trace.WithSpanKind(trace.SpanKindConsumer),
				trace.WithAttributes(
					attribute.String("messaging.system", "rabbitmq"),
					attribute.String("messaging.destination", messageType),
					attribute.String("messaging.operation", "process"),
				),
			)
			defer span.End()
			
			// 處理訊息
			err := next(ctx, message)
			
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			} else {
				span.SetStatus(codes.Ok, "")
			}
			
			return err
		}
	}
}