package repository

import (
	"github.com/fastenmind/fastener-api/internal/models"
	"gorm.io/gorm"
	"time"
)

type ExchangeRateRepository struct {
	db *gorm.DB
}

func NewExchangeRateRepository(db *gorm.DB) *ExchangeRateRepository {
	return &ExchangeRateRepository{db: db}
}

// GetLatestRate 獲取最新匯率
func (r *ExchangeRateRepository) GetLatestRate(fromCurrency, toCurrency, companyID string) (*models.ExchangeRate, error) {
	var rate models.ExchangeRate
	
	err := r.db.Where("company_id = ? AND from_currency = ? AND to_currency = ? AND effective_date <= ?", 
		companyID, fromCurrency, toCurrency, time.Now()).
		Order("effective_date DESC").
		First(&rate).Error
	
	if err == gorm.ErrRecordNotFound {
		// 如果找不到匯率，返回默認值 1.0
		return &models.ExchangeRate{
			FromCurrency: fromCurrency,
			ToCurrency:   toCurrency,
			Rate:         1.0,
			EffectiveDate: time.Now(),
		}, nil
	}
	
	return &rate, err
}

// CreateRate 創建匯率
func (r *ExchangeRateRepository) CreateRate(rate *models.ExchangeRate) error {
	return r.db.Create(rate).Error
}

// UpdateRate 更新匯率
func (r *ExchangeRateRepository) UpdateRate(rate *models.ExchangeRate) error {
	return r.db.Save(rate).Error
}

// GetRateHistory 獲取匯率歷史
func (r *ExchangeRateRepository) GetRateHistory(fromCurrency, toCurrency, companyID string, limit int) ([]models.ExchangeRate, error) {
	var rates []models.ExchangeRate
	
	err := r.db.Where("company_id = ? AND from_currency = ? AND to_currency = ?", 
		companyID, fromCurrency, toCurrency).
		Order("effective_date DESC").
		Limit(limit).
		Find(&rates).Error
	
	return rates, err
}