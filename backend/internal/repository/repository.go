package repository

import "gorm.io/gorm"

// Repositories holds all repository instances
type Repositories struct {
	Account        AccountRepository
	Company        CompanyRepository
	Customer       CustomerRepository
	Inquiry        InquiryRepository
	Process        ProcessRepository
	Equipment      EquipmentRepository
	AssignmentRule AssignmentRuleRepository
	Tariff         TariffRepository
	Compliance     ComplianceRepository
	N8N            N8NRepository
	Quote          QuoteRepository
	Order          OrderRepository
	Inventory      InventoryRepository
	Trade          TradeRepository
	Advanced       AdvancedRepository
	Integration    IntegrationRepository
}

// NewRepositories creates new repository instances
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Account:        NewAccountRepositoryGorm(db),
		Company:        NewCompanyRepository(db),
		Customer:       NewCustomerRepository(db),
		Inquiry:        NewInquiryRepository(db),
		Process:        NewProcessRepository(db),
		Equipment:      NewEquipmentRepository(db),
		AssignmentRule: NewAssignmentRuleRepository(db),
		Tariff:         NewTariffRepository(db),
		Compliance:     NewComplianceRepository(db),
		N8N:            NewN8NRepository(db),
		Quote:          NewQuoteRepository(db),
		Order:          NewOrderRepository(db),
		Inventory:      NewInventoryRepository(db),
		Trade:          NewTradeRepository(db),
		Advanced:       NewAdvancedRepository(db),
		Integration:    NewIntegrationRepository(db),
	}
}