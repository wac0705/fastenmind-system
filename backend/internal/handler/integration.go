package handler

import (
	"github.com/fastenmind/fastener-api/internal/service"
)

type IntegrationHandler struct {
	integrationService service.IntegrationService
}

func NewIntegrationHandler(integrationService service.IntegrationService) *IntegrationHandler {
	return &IntegrationHandler{
		integrationService: integrationService,
	}
}

// All method implementations are in integration_stub.go until IntegrationService is ready