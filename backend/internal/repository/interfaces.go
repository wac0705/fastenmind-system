package repository

// Placeholder interfaces for compilation
type CompanyRepository interface{}
type CustomerRepository interface{}
// InquiryRepository is defined in inquiry_repository.go
type ProcessRepository interface{}
type EquipmentRepository interface{}
type AssignmentRuleRepository interface{}
// TariffRepository is defined in tariff_repository.go
type ComplianceRepository interface{}

// Placeholder constructors
func NewCompanyRepository(db interface{}) CompanyRepository { return nil }
func NewCustomerRepository(db interface{}) CustomerRepository { return nil }
// NewInquiryRepository is defined in inquiry_repository.go
func NewProcessRepository(db interface{}) ProcessRepository { return nil }
func NewEquipmentRepository(db interface{}) EquipmentRepository { return nil }
func NewAssignmentRuleRepository(db interface{}) AssignmentRuleRepository { return nil }
// NewTariffRepository is defined in tariff_repository.go
func NewComplianceRepository(db interface{}) ComplianceRepository { return nil }