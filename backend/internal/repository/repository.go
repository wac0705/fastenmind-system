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
		Company:        NewCompanyRepositoryGorm(db),
		Customer:       NewCustomerRepositoryGorm(db),
		Inquiry:        NewInquiryRepositoryGorm(db),
		Process:        NewProcessRepositoryGorm(db),
		Equipment:      NewEquipmentRepositoryGorm(db),
		AssignmentRule: NewAssignmentRuleRepositoryGorm(db),
		Tariff:         NewTariffRepository(db),
		Compliance:     NewComplianceRepositoryGorm(db),
		N8N:            NewN8NRepositoryGorm(db),
		Quote:          NewQuoteRepositoryGorm(db),
		Order:          NewOrderRepositoryGorm(db),
		Inventory:      NewInventoryRepositoryGorm(db),
		Trade:          NewTradeRepositoryGorm(db),
		Advanced:       NewAdvancedRepositoryGorm(db),
		Integration:    NewIntegrationRepositoryGorm(db),
	}
}