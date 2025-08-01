package repositories

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/fastenmind/fastener-api/internal/model"
)

type CompanyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) *CompanyRepository {
	return &CompanyRepository{db: db}
}

func (r *CompanyRepository) GetCompanyByID(id uuid.UUID) (*model.Company, error) {
	var company model.Company
	err := r.db.First(&company, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &company, nil
}

func (r *CompanyRepository) CreateCompany(company *model.Company) error {
	return r.db.Create(company).Error
}

func (r *CompanyRepository) UpdateCompany(company *model.Company) error {
	return r.db.Save(company).Error
}

func (r *CompanyRepository) GetAllCompanies() ([]model.Company, error) {
	var companies []model.Company
	err := r.db.Find(&companies).Error
	return companies, err
}