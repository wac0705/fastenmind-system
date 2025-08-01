package handler

import "github.com/fastenmind/fastener-api/internal/service"

// Handlers holds all handler instances
type Handlers struct {
	Auth           *AuthHandler
	Account        *AccountHandler
	Company        *CompanyHandler
	Customer       *CustomerHandler
	Inquiry        *InquiryHandler
	Process        *ProcessHandler
	Equipment      *EquipmentHandler
	AssignmentRule *AssignmentRuleHandler
	Tariff         *TariffHandler
	Compliance     *ComplianceHandler
	N8N            *N8NHandler
	Quote          *QuoteHandler
	Order          *OrderHandler
	Inventory      *InventoryHandler
	Trade          *TradeHandler
	Advanced       *AdvancedHandler
	Integration    *IntegrationHandler
}

// NewHandlers creates new handler instances
func NewHandlers(services *service.Services) *Handlers {
	return &Handlers{
		Auth:           NewAuthHandler(services.Auth),
		Account:        NewAccountHandler(services.Account),
		Company:        NewCompanyHandler(services.Company),
		Customer:       NewCustomerHandler(services.Customer),
		Inquiry:        NewInquiryHandler(services.Inquiry),
		Process:        NewProcessHandler(services.Process),
		Equipment:      NewEquipmentHandler(services.Equipment),
		AssignmentRule: NewAssignmentRuleHandler(services.AssignmentRule),
		Tariff:         NewTariffHandler(services.Tariff),
		Compliance:     NewComplianceHandler(services.Compliance),
		N8N:            NewN8NHandler(services.N8N),
		Quote:          NewQuoteHandler(services.Quote),
		Order:          NewOrderHandler(services.Order),
		Inventory:      NewInventoryHandler(services.Inventory),
		Trade:          NewTradeHandler(services.Trade),
		Advanced:       NewAdvancedHandler(services.Advanced),
		Integration:    NewIntegrationHandler(services.Integration),
	}
}