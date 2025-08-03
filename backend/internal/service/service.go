package service

import (
	"github.com/fastenmind/fastener-api/internal/config"
	"github.com/fastenmind/fastener-api/internal/repository"
	"gorm.io/gorm"
)

// Services holds all service instances
type Services struct {
	Account            AccountService
	Auth               AuthService
	Company            CompanyService
	Customer           CustomerService
	Inquiry            InquiryService
	Process            ProcessService
	Equipment          EquipmentService
	AssignmentRule     AssignmentRuleService
	EngineerAssignment *EngineerAssignmentService
	ProcessCost        *ProcessCostService
	Tariff             TariffService
	Compliance         ComplianceService
	N8N                N8NService
	Quote              QuoteService
	Order              OrderService
	Inventory          InventoryService
	Trade              TradeService
	Advanced           AdvancedService
	Integration        IntegrationService
	Report             ReportService
}

// NewServices creates new service instances
func NewServices(repos *repository.Repositories, cfg *config.Config, db *gorm.DB) *Services {
	n8nService := NewN8NService(repos.N8N)
	pdfGenerator := NewPDFGenerator()
	
	exchangeRateRepo := repository.NewExchangeRateRepository(db)
	
	return &Services{
		Account:            NewAccountService(repos.Account, cfg),
		Auth:               NewAuthService(repos.Account, cfg),
		Company:            NewCompanyService(repos.Company),
		Customer:           NewCustomerService(repos.Customer),
		Inquiry:            NewInquiryService(repos.Inquiry, repos.Account, repos.AssignmentRule),
		Process:            NewProcessService(repos.Process),
		Equipment:          NewEquipmentService(repos.Equipment),
		AssignmentRule:     NewAssignmentRuleService(repos.AssignmentRule),
		EngineerAssignment: NewEngineerAssignmentService(repos.EngineerAssignment, repos.Inquiry, repos.Account),
		ProcessCost:        NewProcessCostService(repos.ProcessCost, repos.Material, repos.Equipment.(*repository.EquipmentRepository), exchangeRateRepo),
		Tariff:             NewTariffService(repos.Tariff),
		Compliance:         NewComplianceService(repos.Compliance),
		N8N:                n8nService,
		Quote:              NewQuoteService(repos.Quote, repos.Inquiry, repos.Customer, n8nService, pdfGenerator),
		Order:              NewOrderService(repos.Order, repos.Quote, repos.Customer, n8nService),
		Inventory:          NewInventoryService(repos.Inventory, repos.Order, n8nService),
		Trade:              NewTradeService(repos.Trade),
		Advanced:           NewAdvancedService(),
		Integration:        NewIntegrationService(),
		Report:             NewReportService(repos.Report, repos.Company, repos.User),
	}
}