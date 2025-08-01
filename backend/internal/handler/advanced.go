package handler

import (
	"github.com/fastenmind/fastener-api/internal/service"
)

type AdvancedHandler struct {
	advancedService service.AdvancedService
}

func NewAdvancedHandler(advancedService service.AdvancedService) *AdvancedHandler {
	return &AdvancedHandler{
		advancedService: advancedService,
	}
}

// All method implementations are in advanced_stub.go until AdvancedService is ready