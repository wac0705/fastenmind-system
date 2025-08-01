package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/fastenmind/fastener-api/internal/model"
)

// InquiryRepository handles inquiry-related data operations
type InquiryRepository interface {
	Create(inquiry *model.Inquiry) error
	Get(id uuid.UUID) (*model.Inquiry, error)
	Update(inquiry *model.Inquiry) error
	Delete(id uuid.UUID) error
	List(companyID uuid.UUID, params map[string]interface{}) ([]model.Inquiry, int64, error)
}

type inquiryRepository struct {
	db *gorm.DB
}

// NewInquiryRepository creates a new inquiry repository
func NewInquiryRepository(db interface{}) InquiryRepository {
	gormDB, ok := db.(*gorm.DB)
	if !ok {
		return &inquiryRepository{}
	}
	return &inquiryRepository{db: gormDB}
}

func (r *inquiryRepository) Create(inquiry *model.Inquiry) error {
	return r.db.Create(inquiry).Error
}

func (r *inquiryRepository) Get(id uuid.UUID) (*model.Inquiry, error) {
	var inquiry model.Inquiry
	if err := r.db.First(&inquiry, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &inquiry, nil
}

func (r *inquiryRepository) Update(inquiry *model.Inquiry) error {
	return r.db.Save(inquiry).Error
}

func (r *inquiryRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Inquiry{}, id).Error
}

func (r *inquiryRepository) List(companyID uuid.UUID, params map[string]interface{}) ([]model.Inquiry, int64, error) {
	var inquiries []model.Inquiry
	var total int64
	
	query := r.db.Model(&model.Inquiry{}).Where("company_id = ?", companyID)
	
	// Count total
	query.Count(&total)
	
	// Apply pagination if provided
	if limit, ok := params["limit"].(int); ok && limit > 0 {
		query = query.Limit(limit)
	}
	if offset, ok := params["offset"].(int); ok && offset > 0 {
		query = query.Offset(offset)
	}
	
	err := query.Find(&inquiries).Error
	return inquiries, total, err
}