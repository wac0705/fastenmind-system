package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/fastenmind/fastener-api/internal/model"
	"github.com/fastenmind/fastener-api/internal/repository"
)

var (
	// ErrCustomerNotFound is returned when a customer is not found
	ErrCustomerNotFound = errors.New("customer not found")
	// ErrCustomerCodeExists is returned when a customer code already exists
	ErrCustomerCodeExists = errors.New("customer code already exists")
)

type customerService struct {
	repo repository.CustomerRepository
}

// NewCustomerService creates a new customer service
func NewCustomerService(repo repository.CustomerRepository) CustomerService {
	return &customerService{
		repo: repo,
	}
}

// List returns a paginated list of customers
func (s *customerService) List(ctx context.Context, filter model.CustomerFilter) ([]*model.Customer, *model.Pagination, error) {
	return s.repo.List(ctx, filter)
}

// Create creates a new customer
func (s *customerService) Create(ctx context.Context, req model.CreateCustomerRequest) (*model.Customer, error) {
	// Check if customer code already exists
	exists, err := s.repo.ExistsByCode(ctx, req.CompanyID, req.CustomerCode)
	if err != nil {
		return nil, fmt.Errorf("failed to check customer code: %w", err)
	}
	if exists {
		return nil, ErrCustomerCodeExists
	}

	customer := &model.Customer{
		Base: model.Base{
			ID: uuid.New(),
		},
		CompanyID:       req.CompanyID,
		CustomerCode:    req.CustomerCode,
		Name:            req.Name,
		NameEn:          req.NameEn,
		ShortName:       req.ShortName,
		Country:         req.Country,
		TaxID:           req.TaxID,
		Address:         req.Address,
		ShippingAddress: req.ShippingAddress,
		ContactPerson:   req.ContactPerson,
		ContactPhone:    req.ContactPhone,
		ContactEmail:    req.ContactEmail,
		PaymentTerms:    req.PaymentTerms,
		CreditLimit:     req.CreditLimit,
		Currency:        req.Currency,
		IsActive:        req.IsActive,
	}

	err = s.repo.Create(ctx, customer)
	if err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	return customer, nil
}

// GetByID returns a customer by ID
func (s *customerService) GetByID(ctx context.Context, id, companyID uuid.UUID) (*model.Customer, error) {
	customer, err := s.repo.GetByID(ctx, id, companyID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, ErrCustomerNotFound
		}
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	return customer, nil
}

// Update updates a customer
func (s *customerService) Update(ctx context.Context, req model.UpdateCustomerRequest) (*model.Customer, error) {
	// Get existing customer
	customer, err := s.repo.GetByID(ctx, req.ID, req.CompanyID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, ErrCustomerNotFound
		}
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	// Check if customer code changed and already exists
	if customer.CustomerCode != req.CustomerCode {
		exists, err := s.repo.ExistsByCode(ctx, req.CompanyID, req.CustomerCode)
		if err != nil {
			return nil, fmt.Errorf("failed to check customer code: %w", err)
		}
		if exists {
			return nil, ErrCustomerCodeExists
		}
	}

	// Update fields
	customer.CustomerCode = req.CustomerCode
	customer.Name = req.Name
	customer.NameEn = req.NameEn
	customer.ShortName = req.ShortName
	customer.Country = req.Country
	customer.TaxID = req.TaxID
	customer.Address = req.Address
	customer.ShippingAddress = req.ShippingAddress
	customer.ContactPerson = req.ContactPerson
	customer.ContactPhone = req.ContactPhone
	customer.ContactEmail = req.ContactEmail
	customer.PaymentTerms = req.PaymentTerms
	customer.CreditLimit = req.CreditLimit
	customer.Currency = req.Currency
	customer.IsActive = req.IsActive

	err = s.repo.Update(ctx, customer)
	if err != nil {
		return nil, fmt.Errorf("failed to update customer: %w", err)
	}

	return customer, nil
}

// Delete deletes a customer
func (s *customerService) Delete(ctx context.Context, id, companyID uuid.UUID) error {
	// Check if customer exists
	_, err := s.repo.GetByID(ctx, id, companyID)
	if err != nil {
		if err == repository.ErrNotFound {
			return ErrCustomerNotFound
		}
		return fmt.Errorf("failed to get customer: %w", err)
	}

	err = s.repo.Delete(ctx, id, companyID)
	if err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}

	return nil
}

// GetStatistics returns customer statistics
func (s *customerService) GetStatistics(ctx context.Context, customerID, companyID uuid.UUID) (*model.CustomerStatistics, error) {
	stats, err := s.repo.GetStatistics(ctx, customerID, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer statistics: %w", err)
	}

	return stats, nil
}

// Export exports customers in the specified format
func (s *customerService) Export(ctx context.Context, filter model.CustomerFilter, format string) ([]byte, string, error) {
	// Get all customers without pagination
	filter.Page = 0
	filter.PageSize = 0
	customers, _, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, "", fmt.Errorf("failed to list customers: %w", err)
	}

	switch format {
	case "csv":
		data, err := s.exportToCSV(customers)
		if err != nil {
			return nil, "", fmt.Errorf("failed to export to CSV: %w", err)
		}
		return data, "text/csv", nil
	case "json":
		data, err := s.exportToJSON(customers)
		if err != nil {
			return nil, "", fmt.Errorf("failed to export to JSON: %w", err)
		}
		return data, "application/json", nil
	default:
		return nil, "", fmt.Errorf("unsupported export format: %s", format)
	}
}

// GetCreditHistory returns customer credit history
func (s *customerService) GetCreditHistory(ctx context.Context, customerID, companyID uuid.UUID) ([]*model.CreditHistory, error) {
	history, err := s.repo.GetCreditHistory(ctx, customerID, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get credit history: %w", err)
	}

	return history, nil
}

// Helper function to export customers to CSV
func (s *customerService) exportToCSV(customers []*model.Customer) ([]byte, error) {
	// TODO: Implement CSV export
	return []byte("customer_code,name,country,contact_email\n"), nil
}

// Helper function to export customers to JSON
func (s *customerService) exportToJSON(customers []*model.Customer) ([]byte, error) {
	// TODO: Implement JSON export
	return []byte("[]"), nil
}