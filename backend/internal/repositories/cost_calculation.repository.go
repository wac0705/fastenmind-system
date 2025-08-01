package repositories

import (
	"time"

	"github.com/fastenmind/fastener-api/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CostCalculationRepository struct {
	db *gorm.DB
}

func NewCostCalculationRepository(db *gorm.DB) *CostCalculationRepository {
	return &CostCalculationRepository{db: db}
}

// GetCurrentCostParameters 獲取當前有效的成本參數
func (r *CostCalculationRepository) GetCurrentCostParameters() (map[string]models.CostParameter, error) {
	var params []models.CostParameter
	now := time.Now()
	
	err := r.db.Where("effective_date <= ? AND (end_date IS NULL OR end_date >= ?)", now, now).
		Find(&params).Error
	if err != nil {
		return nil, err
	}

	// 轉換為 map 方便查詢
	paramMap := make(map[string]models.CostParameter)
	for _, param := range params {
		// 如果有多個同類型參數，使用最新的
		if existing, exists := paramMap[param.ParameterType]; !exists || param.EffectiveDate.After(existing.EffectiveDate) {
			paramMap[param.ParameterType] = param
		}
	}

	return paramMap, nil
}

// GetProcessRoutes 獲取產品類別的製程路線
func (r *CostCalculationRepository) GetProcessRoutes(productCategory string) ([]models.ProductProcessRoute, error) {
	var routes []models.ProductProcessRoute
	query := r.db.Where("is_active = ?", true)
	
	if productCategory != "" {
		query = query.Where("product_category = ?", productCategory)
	}
	
	err := query.Order("is_default DESC, route_name").Find(&routes).Error
	return routes, err
}

// GetProcessRouteByID 根據ID獲取完整的製程路線（包含明細）
func (r *CostCalculationRepository) GetProcessRouteByID(id uuid.UUID) (*models.ProductProcessRoute, error) {
	var route models.ProductProcessRoute
	err := r.db.Preload("RouteDetails", func(db *gorm.DB) *gorm.DB {
		return db.Order("sequence")
	}).
		Preload("RouteDetails.ProcessStep").
		Preload("RouteDetails.Equipment").
		First(&route, "id = ?", id).Error
	return &route, err
}

// GetAllProcessSteps 獲取所有製程步驟
func (r *CostCalculationRepository) GetAllProcessSteps() ([]models.ProcessStep, error) {
	var steps []models.ProcessStep
	err := r.db.Preload("ProcessCategory").
		Preload("DefaultEquipment").
		Where("is_active = ?", true).
		Order("sort_order, name").
		Find(&steps).Error
	return steps, err
}

// GetProcessStepByID 根據ID獲取製程步驟
func (r *CostCalculationRepository) GetProcessStepByID(id uuid.UUID) (*models.ProcessStep, error) {
	var step models.ProcessStep
	err := r.db.Preload("ProcessCategory").
		Preload("DefaultEquipment").
		First(&step, "id = ?", id).Error
	return &step, err
}

// GetEquipmentList 獲取設備列表
func (r *CostCalculationRepository) GetEquipmentList(categoryID *uuid.UUID) ([]models.Equipment, error) {
	var equipment []models.Equipment
	query := r.db.Where("is_active = ?", true)
	
	if categoryID != nil {
		query = query.Where("process_category_id = ?", *categoryID)
	}
	
	err := query.Preload("ProcessCategory").
		Order("name").
		Find(&equipment).Error
	return equipment, err
}

// GetEquipmentByID 根據ID獲取設備
func (r *CostCalculationRepository) GetEquipmentByID(id uuid.UUID) (*models.Equipment, error) {
	var equipment models.Equipment
	err := r.db.Preload("ProcessCategory").
		First(&equipment, "id = ?", id).Error
	return &equipment, err
}

// GetCalculationByID 根據ID獲取成本計算（包含明細）
func (r *CostCalculationRepository) GetCalculationByID(id uuid.UUID) (*models.CostCalculation, error) {
	var calc models.CostCalculation
	err := r.db.Preload("Inquiry").
		Preload("Route").
		Preload("Route.RouteDetails", func(db *gorm.DB) *gorm.DB {
			return db.Order("sequence")
		}).
		Preload("CalculatedByUser").
		Preload("ApprovedByUser").
		Preload("Details", func(db *gorm.DB) *gorm.DB {
			return db.Order("sequence")
		}).
		Preload("Details.ProcessStep").
		Preload("Details.Equipment").
		First(&calc, "id = ?", id).Error
	return &calc, err
}

// GetCalculationsByInquiry 根據詢價單ID獲取成本計算列表
func (r *CostCalculationRepository) GetCalculationsByInquiry(inquiryID uuid.UUID) ([]models.CostCalculation, error) {
	var calculations []models.CostCalculation
	err := r.db.Where("inquiry_id = ?", inquiryID).
		Preload("CalculatedByUser").
		Preload("ApprovedByUser").
		Order("created_at DESC").
		Find(&calculations).Error
	return calculations, err
}

// GetCalculations 獲取成本計算列表（分頁）
func (r *CostCalculationRepository) GetCalculations(page, pageSize int, status string) ([]models.CostCalculation, int64, error) {
	var calculations []models.CostCalculation
	var total int64

	query := r.db.Model(&models.CostCalculation{})
	
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 計算總數
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 獲取分頁資料
	offset := (page - 1) * pageSize
	err := query.Preload("Inquiry").
		Preload("CalculatedByUser").
		Preload("ApprovedByUser").
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&calculations).Error

	return calculations, total, err
}

// CreateProcessCategory 創建製程類別
func (r *CostCalculationRepository) CreateProcessCategory(category *models.ProcessCategory) error {
	return r.db.Create(category).Error
}

// UpdateProcessCategory 更新製程類別
func (r *CostCalculationRepository) UpdateProcessCategory(id uuid.UUID, category models.ProcessCategory) error {
	return r.db.Model(&models.ProcessCategory{}).
		Where("id = ?", id).
		Updates(category).Error
}

// CreateEquipment 創建設備
func (r *CostCalculationRepository) CreateEquipment(equipment *models.Equipment) error {
	return r.db.Create(equipment).Error
}

// UpdateEquipment 更新設備
func (r *CostCalculationRepository) UpdateEquipment(id uuid.UUID, equipment models.Equipment) error {
	return r.db.Model(&models.Equipment{}).
		Where("id = ?", id).
		Updates(equipment).Error
}

// CreateProcessStep 創建製程步驟
func (r *CostCalculationRepository) CreateProcessStep(step *models.ProcessStep) error {
	return r.db.Create(step).Error
}

// UpdateProcessStep 更新製程步驟
func (r *CostCalculationRepository) UpdateProcessStep(id uuid.UUID, step models.ProcessStep) error {
	return r.db.Model(&models.ProcessStep{}).
		Where("id = ?", id).
		Updates(step).Error
}

// CreateProcessRoute 創建製程路線
func (r *CostCalculationRepository) CreateProcessRoute(route *models.ProductProcessRoute) error {
	return r.db.Create(route).Error
}

// UpdateProcessRoute 更新製程路線
func (r *CostCalculationRepository) UpdateProcessRoute(id uuid.UUID, route models.ProductProcessRoute) error {
	return r.db.Model(&models.ProductProcessRoute{}).
		Where("id = ?", id).
		Updates(route).Error
}

// DeleteProcessRoute 刪除製程路線
func (r *CostCalculationRepository) DeleteProcessRoute(id uuid.UUID) error {
	// 先刪除明細
	if err := r.db.Where("route_id = ?", id).Delete(&models.ProcessRouteDetail{}).Error; err != nil {
		return err
	}
	// 再刪除主檔
	return r.db.Delete(&models.ProductProcessRoute{}, "id = ?", id).Error
}

// UpdateCostParameter 更新成本參數
func (r *CostCalculationRepository) UpdateCostParameter(param models.CostParameter) error {
	// 將舊參數設定結束日期
	yesterday := time.Now().AddDate(0, 0, -1)
	if err := r.db.Model(&models.CostParameter{}).
		Where("parameter_type = ? AND end_date IS NULL", param.ParameterType).
		Update("end_date", yesterday).Error; err != nil {
		return err
	}
	
	// 創建新參數
	param.EffectiveDate = time.Now()
	return r.db.Create(&param).Error
}

// GetCostParameters 獲取成本參數歷史
func (r *CostCalculationRepository) GetCostParameters() ([]models.CostParameter, error) {
	var params []models.CostParameter
	err := r.db.Order("parameter_type, effective_date DESC").Find(&params).Error
	return params, err
}