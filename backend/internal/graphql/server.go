package graphql

import (
	"context"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/fastenmind/fastener-api/internal/domain/cqrs"
	"github.com/fastenmind/fastener-api/internal/domain/events"
	"github.com/fastenmind/fastener-api/internal/graphql/resolver"
	"github.com/fastenmind/fastener-api/internal/infrastructure/cache"
	"github.com/fastenmind/fastener-api/internal/service"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

// Config holds GraphQL server configuration
type Config struct {
	Services    *service.Services
	CommandBus  cqrs.CommandBus
	QueryBus    cqrs.QueryBus
	EventBus    events.EventBus
	Cache       *cache.CacheService
	EnablePlayground bool
}

// NewServer creates a new GraphQL server
func NewServer(cfg Config) *handler.Server {
	// Create resolver with dependencies
	_ = resolver.NewResolver(
		cfg.Services,
		cfg.CommandBus,
		cfg.QueryBus,
		cfg.EventBus,
		cfg.Cache,
	)
	
	// TODO: Implement actual GraphQL schema generation
	// For now, create a minimal server to avoid compilation errors
	srv := handler.New(nil)

	// Configure server
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	// WebSocket support for subscriptions
	srv.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Configure origin checking for production
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		InitFunc: func(ctx context.Context, initPayload transport.InitPayload) (context.Context, error) {
			// Initialize WebSocket connection
			// Extract auth token from connection params if needed
			return ctx, nil
		},
		KeepAlivePingInterval: 10,
	})

	// Add extensions
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})

	// TODO: Error handling will be implemented with schema
	// srv.SetErrorPresenter(...)
	// srv.SetRecoverFunc(...)

	// Query complexity limiting
	srv.Use(extension.FixedComplexityLimit(1000))

	return srv
}

// RegisterRoutes registers GraphQL routes with Echo
func RegisterRoutes(e *echo.Echo, cfg Config) {
	srv := NewServer(cfg)

	// GraphQL endpoint
	e.POST("/graphql", echo.WrapHandler(srv))
	e.GET("/graphql", echo.WrapHandler(srv))

	// GraphQL Playground
	if cfg.EnablePlayground {
		e.GET("/playground", echo.WrapHandler(playground.Handler("GraphQL Playground", "/graphql")))
	}

	// Subscription endpoint
	e.Any("/graphql/ws", echo.WrapHandler(srv))
}

// AuthMiddleware adds authentication context to GraphQL requests
func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Extract user from Echo context
		user := c.Get("user")
		if user != nil {
			// Add user to request context for GraphQL resolvers
			ctx := context.WithValue(c.Request().Context(), "user", user)
			c.SetRequest(c.Request().WithContext(ctx))
		}
		
		return next(c)
	}
}

// DataLoaderMiddleware adds data loaders to context for N+1 query prevention
func DataLoaderMiddleware(services *service.Services) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Create data loaders
			loaders := NewDataLoaders(services)
			
			// Add to context
			ctx := context.WithValue(c.Request().Context(), "loaders", loaders)
			c.SetRequest(c.Request().WithContext(ctx))
			
			return next(c)
		}
	}
}

// DataLoaders holds all data loaders
type DataLoaders struct {
	CompanyByID  *CompanyLoader
	CustomerByID *CustomerLoader
	AccountByID  *AccountLoader
}

// NewDataLoaders creates new data loaders
func NewDataLoaders(services *service.Services) *DataLoaders {
	return &DataLoaders{
		CompanyByID:  NewCompanyLoader(services),
		CustomerByID: NewCustomerLoader(services),
		AccountByID:  NewAccountLoader(services),
	}
}

// GetLoaders extracts data loaders from context
func GetLoaders(ctx context.Context) *DataLoaders {
	return ctx.Value("loaders").(*DataLoaders)
}