package repository

import (
	"github.com/fastenmind/fastener-api/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuoteRepository interface {
	Create(quote *models.Quote) error
	Update(quote *models.Quote) error
	Delete(id uuid.UUID) error
	Get(id uuid.UUID) (*models.Quote, error)
	GetWithDetails(id uuid.UUID) (*models.Quote, error)
	List(companyID uuid.UUID, params map[string]interface{}) ([]models.Quote, int64, error)
	
	// Version management
	CreateVersion(version *models.QuoteVersion) error
	GetVersions(quoteID uuid.UUID) ([]models.QuoteVersion, error)
	
	// Activity log
	LogActivity(activity *models.QuoteActivity) error
	GetActivities(quoteID uuid.UUID) ([]models.QuoteActivity, error)
}

type quoteRepository struct {
	db *gorm.DB
}

func NewQuoteRepository(db interface{}) QuoteRepository {
	gormDB, ok := db.(*gorm.DB)
	if !ok {
		panic("invalid database type, expected *gorm.DB")
	}
	return &quoteRepository{db: gormDB}
}

func (r *quoteRepository) Create(quote *models.Quote) error {
	return r.db.Create(quote).Error
}

func (r *quoteRepository) Update(quote *models.Quote) error {
	return r.db.Save(quote).Error
}

func (r *quoteRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Quote{}, id).Error
}

func (r *quoteRepository) Get(id uuid.UUID) (*models.Quote, error) {
	var quote models.Quote
	if err := r.db.First(&quote, id).Error; err != nil {
		return nil, err
	}
	return &quote, nil
}

func (r *quoteRepository) GetWithDetails(id uuid.UUID) (*models.Quote, error) {
	var quote models.Quote
	err := r.db.Preload("Inquiry").
		Preload("Inquiry.Customer").
		Preload("Customer").
		Preload("Sales").
		Preload("Engineer").
		Preload("Reviewer").
		Preload("SentBy").
		First(&quote, id).Error
		
	if err != nil {
		return nil, err
	}
	return &quote, nil
}

func (r *quoteRepository) List(companyID uuid.UUID, params map[string]interface{}) ([]models.Quote, int64, error) {
	var quotes []models.Quote
	var total int64
	
	query := r.db.Model(&models.Quote{}).Where("company_id = ?", companyID)
	
	// Apply filters
	if status, ok := params["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	
	if customerID, ok := params["customer_id"].(string); ok && customerID != "" {
		query = query.Where("customer_id = ?", customerID)
	}
	
	if engineerID, ok := params["engineer_id"].(string); ok && engineerID != "" {
		query = query.Where("engineer_id = ?", engineerID)
	}
	
	if search, ok := params["search"].(string); ok && search != "" {
		query = query.Where("quote_no LIKE ?", "%"+search+"%")
	}
	
	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Apply pagination
	if page, ok := params["page"].(int); ok && page > 0 {
		pageSize := 20
		if ps, ok := params["page_size"].(int); ok && ps > 0 {
			pageSize = ps
		}
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}
	
	// Load with relations
	if err := query.
		Preload("Inquiry").
		Preload("Customer").
		Preload("Sales").
		Preload("Engineer").
		Order("created_at DESC").
		Find(&quotes).Error; err != nil {
		return nil, 0, err
	}
	
	return quotes, total, nil
}

func (r *quoteRepository) CreateVersion(version *models.QuoteVersion) error {
	return r.db.Create(version).Error
}

func (r *quoteRepository) GetVersions(quoteID uuid.UUID) ([]models.QuoteVersion, error) {
	var versions []models.QuoteVersion
	err := r.db.Where("quote_id = ?", quoteID).
		Preload("Creator").
		Order("version_number DESC").
		Find(&versions).Error
	return versions, err
}

func (r *quoteRepository) LogActivity(activity *models.QuoteActivity) error {
	return r.db.Create(activity).Error
}

func (r *quoteRepository) GetActivities(quoteID uuid.UUID) ([]models.QuoteActivity, error) {
	var activities []models.QuoteActivity
	err := r.db.Where("quote_id = ?", quoteID).
		Preload("User").
		Order("created_at DESC").
		Find(&activities).Error
	return activities, err
}