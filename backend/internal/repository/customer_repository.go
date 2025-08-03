package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/fastenmind/fastener-api/internal/model"
)

var (
	// ErrNotFound is returned when a record is not found
	ErrNotFound = errors.New("record not found")
)

type customerRepository struct {
	db *gorm.DB
}

// NewCustomerRepository creates a new customer repository
func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{
		db: db,
	}
}

// List returns a paginated list of customers
func (r *customerRepository) List(ctx context.Context, filter model.CustomerFilter) ([]*model.Customer, *model.Pagination, error) {
	var customers []*model.Customer
	
	query := r.db.WithContext(ctx).Model(&model.Customer{}).
		Where("company_id = ?", filter.CompanyID)
	
	// Apply filters
	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		query = query.Where("(customer_code ILIKE ? OR name ILIKE ? OR name_en ILIKE ? OR contact_person ILIKE ?)",
			searchPattern, searchPattern, searchPattern, searchPattern)
	}
	
	if filter.Country != "" {
		query = query.Where("country = ?", filter.Country)
	}
	
	if filter.IsActive {
		query = query.Where("is_active = ?", true)
	}
	
	// Count total records
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to count customers: %w", err)
	}
	
	// Apply pagination if requested
	if filter.Page > 0 && filter.PageSize > 0 {
		offset := (filter.Page - 1) * filter.PageSize
		query = query.Offset(offset).Limit(filter.PageSize)
	}
	
	// Fetch customers
	if err := query.Order("created_at DESC").Find(&customers).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to list customers: %w", err)
	}
	
	// Build pagination
	var pagination *model.Pagination
	if filter.Page > 0 && filter.PageSize > 0 {
		totalPages := int((total + int64(filter.PageSize) - 1) / int64(filter.PageSize))
		pagination = &model.Pagination{
			Page:       filter.Page,
			PageSize:   filter.PageSize,
			Total:      int(total),
			TotalPages: totalPages,
		}
	}
	
	return customers, pagination, nil
}

// Create creates a new customer
func (r *customerRepository) Create(ctx context.Context, customer *model.Customer) error {
	if err := r.db.WithContext(ctx).Create(customer).Error; err != nil {
		return fmt.Errorf("failed to create customer: %w", err)
	}
	return nil
}

// GetByID returns a customer by ID
func (r *customerRepository) GetByID(ctx context.Context, id, companyID uuid.UUID) (*model.Customer, error) {
	var customer model.Customer
	
	err := r.db.WithContext(ctx).
		Where("id = ? AND company_id = ?", id, companyID).
		First(&customer).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}
	
	return &customer, nil
}

// Update updates a customer
func (r *customerRepository) Update(ctx context.Context, customer *model.Customer) error {
	result := r.db.WithContext(ctx).
		Model(customer).
		Where("id = ? AND company_id = ?", customer.ID, customer.CompanyID).
		Updates(customer)
	
	if result.Error != nil {
		return fmt.Errorf("failed to update customer: %w", result.Error)
	}
	
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	
	return nil
}

// Delete deletes a customer
func (r *customerRepository) Delete(ctx context.Context, id, companyID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND company_id = ?", id, companyID).
		Delete(&model.Customer{})
	
	if result.Error != nil {
		return fmt.Errorf("failed to delete customer: %w", result.Error)
	}
	
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	
	return nil
}

// ExistsByCode checks if a customer with the given code exists
func (r *customerRepository) ExistsByCode(ctx context.Context, companyID uuid.UUID, code string) (bool, error) {
	var count int64
	
	err := r.db.WithContext(ctx).
		Model(&model.Customer{}).
		Where("company_id = ? AND customer_code = ?", companyID, code).
		Count(&count).Error
	
	if err != nil {
		return false, fmt.Errorf("failed to check customer code: %w", err)
	}
	
	return count > 0, nil
}

// GetStatistics returns customer statistics
func (r *customerRepository) GetStatistics(ctx context.Context, customerID, companyID uuid.UUID) (*model.CustomerStatistics, error) {
	stats := &model.CustomerStatistics{}
	
	// Get customer credit limit
	var customer model.Customer
	if err := r.db.WithContext(ctx).
		Select("credit_limit").
		Where("id = ? AND company_id = ?", customerID, companyID).
		First(&customer).Error; err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}
	
	if customer.CreditLimit != nil {
		stats.CreditLimit = *customer.CreditLimit
	}
	
	// Count inquiries
	var inquiryCount int64
	if err := r.db.WithContext(ctx).
		Table("inquiries").
		Where("customer_id = ? AND company_id = ?", customerID, companyID).
		Count(&inquiryCount).Error; err != nil {
		return nil, fmt.Errorf("failed to count inquiries: %w", err)
	}
	stats.TotalInquiries = int(inquiryCount)
	
	// Count quotes
	var quoteCount int64
	if err := r.db.WithContext(ctx).
		Table("quotes").
		Where("customer_id = ? AND company_id = ?", customerID, companyID).
		Count(&quoteCount).Error; err != nil {
		return nil, fmt.Errorf("failed to count quotes: %w", err)
	}
	stats.TotalQuotes = int(quoteCount)
	
	// Count orders and calculate revenue
	type OrderStats struct {
		Count         int64   `gorm:"column:count"`
		TotalRevenue  float64 `gorm:"column:total_revenue"`
		LastOrderDate *string `gorm:"column:last_order_date"`
	}
	
	var orderStats OrderStats
	if err := r.db.WithContext(ctx).
		Table("orders").
		Select("COUNT(*) as count, COALESCE(SUM(total_amount), 0) as total_revenue, MAX(created_at) as last_order_date").
		Where("customer_id = ? AND company_id = ? AND status NOT IN ('cancelled', 'draft')", customerID, companyID).
		Scan(&orderStats).Error; err != nil {
		return nil, fmt.Errorf("failed to get order stats: %w", err)
	}
	
	stats.TotalOrders = int(orderStats.Count)
	stats.TotalRevenue = orderStats.TotalRevenue
	
	if stats.TotalOrders > 0 {
		stats.AverageOrderValue = stats.TotalRevenue / float64(stats.TotalOrders)
	}
	
	// Calculate credit used (outstanding orders)
	var creditUsed float64
	if err := r.db.WithContext(ctx).
		Table("orders").
		Select("COALESCE(SUM(total_amount - paid_amount), 0)").
		Where("customer_id = ? AND company_id = ? AND status IN ('confirmed', 'in_production', 'shipped')", customerID, companyID).
		Scan(&creditUsed).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate credit used: %w", err)
	}
	
	stats.CreditUsed = creditUsed
	stats.CreditAvailable = stats.CreditLimit - stats.CreditUsed
	
	return stats, nil
}

// GetCreditHistory returns customer credit history
func (r *customerRepository) GetCreditHistory(ctx context.Context, customerID, companyID uuid.UUID) ([]*model.CreditHistory, error) {
	var history []*model.CreditHistory
	
	// TODO: Implement credit history tracking
	// This would typically involve a separate credit_history table
	// For now, return empty array
	
	return history, nil
}