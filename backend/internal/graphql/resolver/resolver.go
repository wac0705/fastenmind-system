package resolver

import (
	"github.com/fastenmind/fastener-api/internal/domain/cqrs"
	"github.com/fastenmind/fastener-api/internal/domain/events"
	"github.com/fastenmind/fastener-api/internal/infrastructure/cache"
	"github.com/fastenmind/fastener-api/internal/service"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver is the resolver root
type Resolver struct {
	services    *service.Services
	commandBus  cqrs.CommandBus
	queryBus    cqrs.QueryBus
	eventBus    events.EventBus
	cache       *cache.CacheService
}

// NewResolver creates a new resolver
func NewResolver(
	services *service.Services,
	commandBus cqrs.CommandBus,
	queryBus cqrs.QueryBus,
	eventBus events.EventBus,
	cache *cache.CacheService,
) *Resolver {
	return &Resolver{
		services:   services,
		commandBus: commandBus,
		queryBus:   queryBus,
		eventBus:   eventBus,
		cache:      cache,
	}
}