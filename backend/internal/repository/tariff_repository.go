package repository

import (
	"time"

	"github.com/fastenmind/fastener-api/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TariffRepository interface {
	// HS Codes
	FindHSCodes(params map[string]interface{}) ([]models.HSCode, int64, error)
	GetHSCode(code string) (*models.HSCode, error)
	CreateHSCode(hsCode *models.HSCode) error
	UpdateHSCode(hsCode *models.HSCode) error
	
	// Tariff Rates
	FindTariffRates(hsCode, fromCountry, toCountry string) ([]models.TariffRate, error)
	GetTariffRate(id uuid.UUID) (*models.TariffRate, error)
	CreateTariffRate(rate *models.TariffRate) error
	UpdateTariffRate(rate *models.TariffRate) error
	GetEffectiveTariffRate(hsCode, fromCountry, toCountry string, date time.Time) (*models.TariffRate, error)
	
	// Trade Agreements
	FindTradeAgreements(countries []string) ([]models.TradeAgreement, error)
	
	// Calculations
	CreateCalculation(calc *models.TariffCalculation) error
	GetCalculationHistory(companyID uuid.UUID, limit int) ([]models.TariffCalculation, error)
}

type tariffRepository struct {
	db *gorm.DB
}

func NewTariffRepository(db interface{}) TariffRepository {
	// Type assert to *gorm.DB
	gormDB, ok := db.(*gorm.DB)
	if !ok {
		panic("invalid database type, expected *gorm.DB")
	}
	return &tariffRepository{db: gormDB}
}

func (r *tariffRepository) FindHSCodes(params map[string]interface{}) ([]models.HSCode, int64, error) {
	var hsCodes []models.HSCode
	var total int64
	
	query := r.db.Model(&models.HSCode{})
	
	if search, ok := params["search"].(string); ok && search != "" {
		query = query.Where("code LIKE ? OR description LIKE ? OR description_en LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	
	if category, ok := params["category"].(string); ok && category != "" {
		query = query.Where("category = ?", category)
	}
	
	if isActive, ok := params["is_active"].(bool); ok {
		query = query.Where("is_active = ?", isActive)
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
	
	if err := query.Order("code").Find(&hsCodes).Error; err != nil {
		return nil, 0, err
	}
	
	return hsCodes, total, nil
}

func (r *tariffRepository) GetHSCode(code string) (*models.HSCode, error) {
	var hsCode models.HSCode
	if err := r.db.Where("code = ?", code).First(&hsCode).Error; err != nil {
		return nil, err
	}
	return &hsCode, nil
}

func (r *tariffRepository) CreateHSCode(hsCode *models.HSCode) error {
	return r.db.Create(hsCode).Error
}

func (r *tariffRepository) UpdateHSCode(hsCode *models.HSCode) error {
	return r.db.Save(hsCode).Error
}

func (r *tariffRepository) FindTariffRates(hsCode, fromCountry, toCountry string) ([]models.TariffRate, error) {
	var rates []models.TariffRate
	query := r.db.Model(&models.TariffRate{})
	
	if hsCode != "" {
		query = query.Where("hs_code = ?", hsCode)
	}
	if fromCountry != "" {
		query = query.Where("from_country = ?", fromCountry)
	}
	if toCountry != "" {
		query = query.Where("to_country = ?", toCountry)
	}
	
	// Only get active rates
	now := time.Now()
	query = query.Where("effective_from <= ?", now).
		Where("(effective_to IS NULL OR effective_to >= ?)", now)
	
	if err := query.Order("effective_from DESC").Find(&rates).Error; err != nil {
		return nil, err
	}
	
	return rates, nil
}

func (r *tariffRepository) GetTariffRate(id uuid.UUID) (*models.TariffRate, error) {
	var rate models.TariffRate
	if err := r.db.First(&rate, id).Error; err != nil {
		return nil, err
	}
	return &rate, nil
}

func (r *tariffRepository) CreateTariffRate(rate *models.TariffRate) error {
	return r.db.Create(rate).Error
}

func (r *tariffRepository) UpdateTariffRate(rate *models.TariffRate) error {
	return r.db.Save(rate).Error
}

func (r *tariffRepository) GetEffectiveTariffRate(hsCode, fromCountry, toCountry string, date time.Time) (*models.TariffRate, error) {
	var rate models.TariffRate
	err := r.db.Where("hs_code = ? AND from_country = ? AND to_country = ?", hsCode, fromCountry, toCountry).
		Where("effective_from <= ?", date).
		Where("(effective_to IS NULL OR effective_to >= ?)", date).
		Order("effective_from DESC").
		First(&rate).Error
		
	if err != nil {
		return nil, err
	}
	return &rate, nil
}

func (r *tariffRepository) FindTradeAgreements(countries []string) ([]models.TradeAgreement, error) {
	var agreements []models.TradeAgreement
	
	if len(countries) == 0 {
		return agreements, nil
	}
	
	// Find agreements where all provided countries are members
	query := r.db.Model(&models.TradeAgreement{})
	for _, country := range countries {
		query = query.Where("? = ANY(member_countries)", country)
	}
	
	now := time.Now()
	query = query.Where("effective_date <= ?", now)
	
	if err := query.Find(&agreements).Error; err != nil {
		return nil, err
	}
	
	return agreements, nil
}

func (r *tariffRepository) CreateCalculation(calc *models.TariffCalculation) error {
	return r.db.Create(calc).Error
}

func (r *tariffRepository) GetCalculationHistory(companyID uuid.UUID, limit int) ([]models.TariffCalculation, error) {
	var calculations []models.TariffCalculation
	
	query := r.db.Where("company_id = ?", companyID).
		Preload("User").
		Preload("Quote").
		Order("created_at DESC")
		
	if limit > 0 {
		query = query.Limit(limit)
	}
	
	if err := query.Find(&calculations).Error; err != nil {
		return nil, err
	}
	
	return calculations, nil
}