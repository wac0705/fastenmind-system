package repositories

import (
	"github.com/fastenmind/fastener-api/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AssignmentRepository struct {
	db *gorm.DB
}

func NewAssignmentRepository(db *gorm.DB) *AssignmentRepository {
	return &AssignmentRepository{db: db}
}

// GetActiveRules 獲取所有啟用的分派規則
func (r *AssignmentRepository) GetActiveRules() ([]models.AssignmentRule, error) {
	var rules []models.AssignmentRule
	err := r.db.Where("is_active = ?", true).
		Order("priority DESC").
		Find(&rules).Error
	return rules, err
}

// GetAllRules 獲取所有分派規則
func (r *AssignmentRepository) GetAllRules() ([]models.AssignmentRule, error) {
	var rules []models.AssignmentRule
	err := r.db.Order("priority DESC, created_at DESC").Find(&rules).Error
	return rules, err
}

// GetRuleByID 根據ID獲取規則
func (r *AssignmentRepository) GetRuleByID(id uuid.UUID) (*models.AssignmentRule, error) {
	var rule models.AssignmentRule
	err := r.db.First(&rule, "id = ?", id).Error
	return &rule, err
}

// CreateRule 創建規則
func (r *AssignmentRepository) CreateRule(rule *models.AssignmentRule) error {
	return r.db.Create(rule).Error
}

// UpdateRule 更新規則
func (r *AssignmentRepository) UpdateRule(id uuid.UUID, rule models.AssignmentRule) error {
	return r.db.Model(&models.AssignmentRule{}).
		Where("id = ?", id).
		Updates(rule).Error
}

// DeleteRule 刪除規則
func (r *AssignmentRepository) DeleteRule(id uuid.UUID) error {
	return r.db.Delete(&models.AssignmentRule{}, "id = ?", id).Error
}

// GetCapableEngineers 獲取有能力處理特定產品類別的工程師
func (r *AssignmentRepository) GetCapableEngineers(productCategory string) ([]models.Account, error) {
	var engineers []models.Account
	err := r.db.Joins("JOIN engineer_capabilities ON engineer_capabilities.engineer_id = accounts.id").
		Where("engineer_capabilities.product_category = ? AND engineer_capabilities.is_active = ? AND accounts.is_active = ? AND accounts.role = ?", 
			productCategory, true, true, "engineer").
		Find(&engineers).Error
	return engineers, err
}

// GetEngineersBySkillLevel 根據技能等級獲取工程師
func (r *AssignmentRepository) GetEngineersBySkillLevel(productCategory string, minSkillLevel int) ([]models.Account, error) {
	var engineers []models.Account
	err := r.db.Joins("JOIN engineer_capabilities ON engineer_capabilities.engineer_id = accounts.id").
		Where("engineer_capabilities.product_category = ? AND engineer_capabilities.skill_level >= ? AND engineer_capabilities.is_active = ? AND accounts.is_active = ? AND accounts.role = ?", 
			productCategory, minSkillLevel, true, true, "engineer").
		Find(&engineers).Error
	return engineers, err
}

// GetEngineerWorkload 獲取工程師工作負載
func (r *AssignmentRepository) GetEngineerWorkload(engineerID uuid.UUID) (*models.EngineerWorkload, error) {
	var workload models.EngineerWorkload
	err := r.db.FirstOrCreate(&workload, models.EngineerWorkload{EngineerID: engineerID}).Error
	return &workload, err
}

// GetAllEngineersWorkload 獲取所有工程師工作負載統計
func (r *AssignmentRepository) GetAllEngineersWorkload() ([]models.EngineerWorkloadStats, error) {
	var stats []models.EngineerWorkloadStats
	
	query := `
		SELECT 
			ew.engineer_id,
			a.full_name as engineer_name,
			ew.current_inquiries,
			ew.completed_today,
			ew.completed_this_week,
			ew.completed_this_month,
			ew.average_completion_hours,
			ARRAY_AGG(DISTINCT ec.product_category) as skill_categories
		FROM engineer_workload ew
		JOIN accounts a ON a.id = ew.engineer_id
		LEFT JOIN engineer_capabilities ec ON ec.engineer_id = ew.engineer_id AND ec.is_active = true
		WHERE a.is_active = true AND a.role = 'engineer'
		GROUP BY ew.engineer_id, a.full_name, ew.current_inquiries, ew.completed_today, 
				 ew.completed_this_week, ew.completed_this_month, ew.average_completion_hours
		ORDER BY ew.current_inquiries ASC
	`
	
	err := r.db.Raw(query).Scan(&stats).Error
	return stats, err
}

// GetAssignmentHistory 獲取分派歷史
func (r *AssignmentRepository) GetAssignmentHistory(inquiryID *uuid.UUID, engineerID *uuid.UUID, limit int) ([]models.AssignmentHistory, error) {
	var history []models.AssignmentHistory
	query := r.db.Preload("Inquiry").
		Preload("AssignedToUser").
		Preload("AssignedFromUser").
		Preload("AssignedByUser").
		Preload("Rule")

	if inquiryID != nil {
		query = query.Where("inquiry_id = ?", *inquiryID)
	}
	if engineerID != nil {
		query = query.Where("assigned_to = ?", *engineerID)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Order("assigned_at DESC").Find(&history).Error
	return history, err
}

// CreateAssignmentHistory 創建分派歷史記錄
func (r *AssignmentRepository) CreateAssignmentHistory(history *models.AssignmentHistory) error {
	return r.db.Create(history).Error
}

// GetEngineerCapabilities 獲取工程師能力列表
func (r *AssignmentRepository) GetEngineerCapabilities(engineerID uuid.UUID) ([]models.EngineerCapability, error) {
	var capabilities []models.EngineerCapability
	err := r.db.Where("engineer_id = ?", engineerID).Find(&capabilities).Error
	return capabilities, err
}

// CreateCapability 創建工程師能力
func (r *AssignmentRepository) CreateCapability(capability *models.EngineerCapability) error {
	return r.db.Create(capability).Error
}

// UpdateCapability 更新工程師能力
func (r *AssignmentRepository) UpdateCapability(capability models.EngineerCapability) error {
	return r.db.Model(&models.EngineerCapability{}).
		Where("engineer_id = ? AND product_category = ? AND process_type = ?", 
			capability.EngineerID, capability.ProductCategory, capability.ProcessType).
		Updates(capability).Error
}

// DeleteCapability 刪除工程師能力
func (r *AssignmentRepository) DeleteCapability(id uuid.UUID) error {
	return r.db.Delete(&models.EngineerCapability{}, "id = ?", id).Error
}

// GetEngineerPreference 獲取工程師偏好設定
func (r *AssignmentRepository) GetEngineerPreference(engineerID uuid.UUID) (*models.EngineerPreference, error) {
	var preference models.EngineerPreference
	err := r.db.Where("engineer_id = ?", engineerID).First(&preference).Error
	return &preference, err
}

// CreatePreference 創建工程師偏好
func (r *AssignmentRepository) CreatePreference(preference *models.EngineerPreference) error {
	return r.db.Create(preference).Error
}

// UpdatePreference 更新工程師偏好
func (r *AssignmentRepository) UpdatePreference(preference models.EngineerPreference) error {
	return r.db.Model(&models.EngineerPreference{}).
		Where("engineer_id = ?", preference.EngineerID).
		Updates(preference).Error
}

// ResetDailyWorkload 重置每日工作量（用於定時任務）
func (r *AssignmentRepository) ResetDailyWorkload() error {
	return r.db.Model(&models.EngineerWorkload{}).
		Update("completed_today", 0).Error
}

// ResetWeeklyWorkload 重置每週工作量
func (r *AssignmentRepository) ResetWeeklyWorkload() error {
	return r.db.Model(&models.EngineerWorkload{}).
		Update("completed_this_week", 0).Error
}

// ResetMonthlyWorkload 重置每月工作量
func (r *AssignmentRepository) ResetMonthlyWorkload() error {
	return r.db.Model(&models.EngineerWorkload{}).
		Update("completed_this_month", 0).Error
}

// InquiryRepository 詢價單 Repository（如果尚未存在）
type InquiryRepository struct {
	db *gorm.DB
}

func NewInquiryRepository(db *gorm.DB) *InquiryRepository {
	return &InquiryRepository{db: db}
}

// GetByID 根據ID獲取詢價單
func (r *InquiryRepository) GetByID(id uuid.UUID) (*models.Inquiry, error) {
	var inquiry models.Inquiry
	err := r.db.Preload("Customer").
		Preload("Sales").
		Preload("AssignedEngineer").
		First(&inquiry, "id = ?", id).Error
	return &inquiry, err
}