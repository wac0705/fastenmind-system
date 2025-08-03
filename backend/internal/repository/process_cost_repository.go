package repository

import (
	"time"
	"github.com/google/uuid"
	"github.com/fastenmind/fastener-api/internal/models"
	"gorm.io/gorm"
)

type ProcessCostRepository struct {
	db *gorm.DB
}

func NewProcessCostRepository(db *gorm.DB) *ProcessCostRepository {
	return &ProcessCostRepository{db: db}
}

// GetTemplates 獲取成本模板
func (r *ProcessCostRepository) GetTemplates(companyID, processType, category string, offset, limit int) ([]models.ProcessCostTemplate, int64, error) {
	var templates []models.ProcessCostTemplate
	var total int64
	
	query := r.db.Model(&models.ProcessCostTemplate{}).
		Where("company_id = ? AND deleted_at IS NULL", companyID)
	
	if processType != "" {
		query = query.Where("process_type = ?", processType)
	}
	
	if category != "" {
		query = query.Where("category = ?", category)
	}
	
	// 計算總數
	query.Count(&total)
	
	// 獲取分頁數據
	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&templates).Error
	
	return templates, total, err
}

// GetTemplateByID 根據ID獲取模板
func (r *ProcessCostRepository) GetTemplateByID(id, companyID string) (*models.ProcessCostTemplateNew, error) {
	var template models.ProcessCostTemplateNew
	err := r.db.Where("id = ? AND company_id = ? AND deleted_at IS NULL", id, companyID).
		First(&template).Error
	return &template, err
}

// CreateTemplate 創建模板
func (r *ProcessCostRepository) CreateTemplate(template *models.ProcessCostTemplate) error {
	return r.db.Create(template).Error
}

// UpdateTemplate 更新模板
func (r *ProcessCostRepository) UpdateTemplate(template *models.ProcessCostTemplate) error {
	return r.db.Save(template).Error
}

// DeleteTemplate 刪除模板
func (r *ProcessCostRepository) DeleteTemplate(id, companyID string) error {
	return r.db.Model(&models.ProcessCostTemplate{}).
		Where("id = ? AND company_id = ?", id, companyID).
		Update("deleted_at", gorm.DeletedAt{}).Error
}

// SaveCalculationHistory 保存計算歷史
func (r *ProcessCostRepository) SaveCalculationHistory(history *models.CostCalculationHistory) error {
	return r.db.Create(history).Error
}

// GetCalculationHistory 獲取計算歷史
func (r *ProcessCostRepository) GetCalculationHistory(companyID, inquiryID, productID string, offset, limit int) ([]models.CostCalculationHistory, int64, error) {
	var histories []models.CostCalculationHistory
	var total int64
	
	query := r.db.Model(&models.CostCalculationHistory{}).
		Where("company_id = ?", companyID)
	
	if inquiryID != "" {
		query = query.Where("inquiry_id = ?", inquiryID)
	}
	
	if productID != "" {
		query = query.Where("product_id = ?", productID)
	}
	
	// 計算總數
	query.Count(&total)
	
	// 獲取分頁數據
	err := query.
		Preload("CalculatedByUser").
		Order("calculated_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&histories).Error
	
	return histories, total, err
}

// GetProcessingRate 獲取加工費率
func (r *ProcessCostRepository) GetProcessingRate(processType, equipmentID, companyID string) (*models.ProcessingRate, error) {
	var rate models.ProcessingRate
	
	query := r.db.Where("company_id = ? AND process_type = ? AND is_active = ?", 
		companyID, processType, true)
	
	if equipmentID != "" {
		query = query.Where("equipment_id = ?", equipmentID)
	}
	
	err := query.First(&rate).Error
	return &rate, err
}

// GetProcessingRates 獲取多個加工費率
func (r *ProcessCostRepository) GetProcessingRates(companyID, processType, equipmentID string) ([]models.ProcessingRate, error) {
	var rates []models.ProcessingRate
	
	query := r.db.Where("company_id = ? AND is_active = ?", companyID, true)
	
	if processType != "" {
		query = query.Where("process_type = ?", processType)
	}
	
	if equipmentID != "" {
		query = query.Where("equipment_id = ?", equipmentID)
	}
	
	err := query.Find(&rates).Error
	return rates, err
}

// GetProcessingRateByID 根據ID獲取加工費率
func (r *ProcessCostRepository) GetProcessingRateByID(id, companyID string) (*models.ProcessingRate, error) {
	var rate models.ProcessingRate
	err := r.db.Where("id = ? AND company_id = ?", id, companyID).First(&rate).Error
	return &rate, err
}

// UpdateProcessingRate 更新加工費率
func (r *ProcessCostRepository) UpdateProcessingRate(rate *models.ProcessingRate) error {
	return r.db.Save(rate).Error
}

// GetSurfaceTreatmentRate 獲取表面處理費率
func (r *ProcessCostRepository) GetSurfaceTreatmentRate(treatmentType, companyID string) (*models.SurfaceTreatmentRate, error) {
	var rate models.SurfaceTreatmentRate
	err := r.db.Where("company_id = ? AND treatment_type = ? AND is_active = ?", 
		companyID, treatmentType, true).First(&rate).Error
	return &rate, err
}

// GetCostDataForAnalysis 獲取成本分析數據
func (r *ProcessCostRepository) GetCostDataForAnalysis(companyID, analysisType, startDate, endDate string) ([]models.CostData, error) {
	var data []models.CostData
	
	query := r.db.Table("cost_calculation_histories").
		Select(`
			calculated_at as date,
			result->>'$.material_cost' as material_cost,
			result->>'$.processing_cost' as processing_cost,
			result->>'$.surface_treatment_cost' as surface_treatment_cost,
			result->>'$.packaging_cost' as packaging_cost,
			result->>'$.overhead_cost' as overhead_cost,
			result->>'$.total_cost' as total_cost
		`).
		Where("company_id = ?", companyID)
	
	if startDate != "" {
		query = query.Where("calculated_at >= ?", startDate)
	}
	
	if endDate != "" {
		query = query.Where("calculated_at <= ?", endDate)
	}
	
	err := query.Scan(&data).Error
	return data, err
}

// GetReportData 獲取報告數據
func (r *ProcessCostRepository) GetReportData(companyID, reportType, startDate, endDate string) (interface{}, error) {
	// 根據報告類型返回不同的數據結構
	switch reportType {
	case "summary":
		return r.getSummaryReportData(companyID, startDate, endDate)
	case "detailed":
		return r.getDetailedReportData(companyID, startDate, endDate)
	case "comparison":
		return r.getComparisonReportData(companyID, startDate, endDate)
	default:
		return r.getSummaryReportData(companyID, startDate, endDate)
	}
}

// GetSettings 獲取成本設定
func (r *ProcessCostRepository) GetSettings(companyID string) (*models.CostSettings, error) {
	var settings models.CostSettings
	err := r.db.Where("company_id = ?", companyID).First(&settings).Error
	if err == gorm.ErrRecordNotFound {
		// 返回默認設定
		return &models.CostSettings{
			ID:           uuid.New().String(),
			CompanyID:    companyID,
			SettingType:  "default",
			SettingName:  "default_settings",
			SettingValue: 15.0, // Default overhead rate
			Unit:         "percentage",
			Description:  "Default cost settings",
			IsActive:     true,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}, nil
	}
	return &settings, err
}

// CreateSettings 創建成本設定
func (r *ProcessCostRepository) CreateSettings(settings *models.CostSettings) error {
	return r.db.Create(settings).Error
}

// UpdateSettings 更新成本設定
func (r *ProcessCostRepository) UpdateSettings(settings *models.CostSettings) error {
	return r.db.Save(settings).Error
}

// 私有輔助函數

func (r *ProcessCostRepository) getSummaryReportData(companyID, startDate, endDate string) (interface{}, error) {
	type SummaryData struct {
		TotalCost         float64 `json:"total_cost"`
		TotalOrders       int     `json:"total_orders"`
		AverageCost       float64 `json:"average_cost"`
		MaterialCostRatio float64 `json:"material_cost_ratio"`
		ProcessCostRatio  float64 `json:"process_cost_ratio"`
		TopMaterials      []struct {
			MaterialName string  `json:"material_name"`
			TotalCost    float64 `json:"total_cost"`
			Percentage   float64 `json:"percentage"`
		} `json:"top_materials"`
		TopProcesses []struct {
			ProcessType string  `json:"process_type"`
			TotalCost   float64 `json:"total_cost"`
			Percentage  float64 `json:"percentage"`
		} `json:"top_processes"`
	}
	
	var summary SummaryData
	
	// 獲取總成本和訂單數
	r.db.Table("cost_calculation_histories").
		Select("SUM(result->>'$.total_cost') as total_cost, COUNT(*) as total_orders").
		Where("company_id = ? AND calculated_at BETWEEN ? AND ?", companyID, startDate, endDate).
		Scan(&summary)
	
	if summary.TotalOrders > 0 {
		summary.AverageCost = summary.TotalCost / float64(summary.TotalOrders)
	}
	
	// 計算成本比例等其他數據...
	
	return summary, nil
}

func (r *ProcessCostRepository) getDetailedReportData(companyID, startDate, endDate string) (interface{}, error) {
	var histories []models.CostCalculationHistory
	
	err := r.db.Where("company_id = ? AND calculated_at BETWEEN ? AND ?", 
		companyID, startDate, endDate).
		Order("calculated_at DESC").
		Find(&histories).Error
	
	return histories, err
}

func (r *ProcessCostRepository) getComparisonReportData(companyID, startDate, endDate string) (interface{}, error) {
	// 實現比較報告數據獲取
	return nil, nil
}