package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/fastenmind/fastener-api/internal/models"
	"github.com/fastenmind/fastener-api/internal/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AssignmentService struct {
	db         *gorm.DB
	assignRepo *repositories.AssignmentRepository
	inquiryRepo *repositories.InquiryRepository
}

func NewAssignmentService(db *gorm.DB) *AssignmentService {
	return &AssignmentService{
		db:          db,
		assignRepo:  repositories.NewAssignmentRepository(db),
		inquiryRepo: repositories.NewInquiryRepository(db),
	}
}

// AutoAssignInquiry 自動分派詢價單
func (s *AssignmentService) AutoAssignInquiry(inquiryID uuid.UUID) (*models.AssignmentHistory, error) {
	// 獲取詢價單資訊
	inquiry, err := s.inquiryRepo.GetByID(inquiryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get inquiry: %w", err)
	}

	if inquiry.AssignedEngineerID != nil {
		return nil, errors.New("inquiry already assigned")
	}

	// 獲取所有啟用的分派規則
	rules, err := s.assignRepo.GetActiveRules()
	if err != nil {
		return nil, fmt.Errorf("failed to get assignment rules: %w", err)
	}

	// 根據優先級排序規則
	sort.Slice(rules, func(i, j int) bool {
		return rules[i].Priority > rules[j].Priority
	})

	// 嘗試根據每個規則分派
	for _, rule := range rules {
		engineer, err := s.findEngineerByRule(inquiry, rule)
		if err != nil {
			continue // 嘗試下一個規則
		}

		if engineer != nil {
			// 執行分派
			return s.assignToEngineer(inquiry, engineer.ID, nil, "auto", fmt.Sprintf("Assigned by rule: %s", rule.RuleName), &rule.ID)
		}
	}

	// 如果沒有規則匹配，使用預設的負載平衡
	engineer, err := s.findEngineerByLoadBalance(inquiry)
	if err != nil {
		return nil, fmt.Errorf("failed to find engineer by load balance: %w", err)
	}

	if engineer == nil {
		return nil, errors.New("no available engineer found")
	}

	return s.assignToEngineer(inquiry, engineer.ID, nil, "auto", "Assigned by default load balance", nil)
}

// ManualAssignInquiry 手動分派詢價單
func (s *AssignmentService) ManualAssignInquiry(req models.AssignmentRequest, assignedBy uuid.UUID) (*models.AssignmentHistory, error) {
	// 獲取詢價單資訊
	inquiry, err := s.inquiryRepo.GetByID(req.InquiryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get inquiry: %w", err)
	}

	// 檢查工程師是否有能力處理
	capable, err := s.isEngineerCapable(req.EngineerID, inquiry.ProductCategory)
	if err != nil {
		return nil, fmt.Errorf("failed to check engineer capability: %w", err)
	}

	if !capable {
		return nil, errors.New("engineer not capable of handling this product category")
	}

	var fromEngineer *uuid.UUID
	if inquiry.AssignedEngineerID != nil {
		fromEngineer = inquiry.AssignedEngineerID
	}

	return s.assignToEngineer(inquiry, req.EngineerID, fromEngineer, req.AssignmentType, req.Reason, nil)
}

// SelfSelectInquiry 工程師自選詢價單
func (s *AssignmentService) SelfSelectInquiry(inquiryID, engineerID uuid.UUID) (*models.AssignmentHistory, error) {
	// 獲取詢價單資訊
	inquiry, err := s.inquiryRepo.GetByID(inquiryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get inquiry: %w", err)
	}

	if inquiry.AssignedEngineerID != nil {
		return nil, errors.New("inquiry already assigned")
	}

	// 檢查工程師是否有能力處理
	capable, err := s.isEngineerCapable(engineerID, inquiry.ProductCategory)
	if err != nil {
		return nil, fmt.Errorf("failed to check engineer capability: %w", err)
	}

	if !capable {
		return nil, errors.New("not capable of handling this product category")
	}

	// 檢查工程師工作量
	workload, err := s.assignRepo.GetEngineerWorkload(engineerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get engineer workload: %w", err)
	}

	preference, err := s.assignRepo.GetEngineerPreference(engineerID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to get engineer preference: %w", err)
	}

	maxAssignments := 10 // 預設值
	if preference != nil {
		maxAssignments = preference.MaxDailyAssignments
	}

	if workload.CurrentInquiries >= maxAssignments {
		return nil, errors.New("maximum daily assignments reached")
	}

	return s.assignToEngineer(inquiry, engineerID, nil, "self_select", "Engineer self-selected", nil)
}

// GetEngineerWorkloadStats 獲取工程師工作量統計
func (s *AssignmentService) GetEngineerWorkloadStats() ([]models.EngineerWorkloadStats, error) {
	return s.assignRepo.GetAllEngineersWorkload()
}

// GetAssignmentHistory 獲取分派歷史
func (s *AssignmentService) GetAssignmentHistory(inquiryID *uuid.UUID, engineerID *uuid.UUID, limit int) ([]models.AssignmentHistory, error) {
	return s.assignRepo.GetAssignmentHistory(inquiryID, engineerID, limit)
}

// UpdateAssignmentRule 更新分派規則
func (s *AssignmentService) UpdateAssignmentRule(ruleID uuid.UUID, rule models.AssignmentRule) error {
	return s.assignRepo.UpdateRule(ruleID, rule)
}

// GetAssignmentRules 獲取所有分派規則
func (s *AssignmentService) GetAssignmentRules() ([]models.AssignmentRule, error) {
	return s.assignRepo.GetAllRules()
}

// UpdateEngineerCapability 更新工程師能力
func (s *AssignmentService) UpdateEngineerCapability(capability models.EngineerCapability) error {
	return s.assignRepo.UpdateCapability(capability)
}

// GetEngineerCapabilities 獲取工程師能力列表
func (s *AssignmentService) GetEngineerCapabilities(engineerID uuid.UUID) ([]models.EngineerCapability, error) {
	return s.assignRepo.GetEngineerCapabilities(engineerID)
}

// UpdateEngineerPreference 更新工程師偏好設定
func (s *AssignmentService) UpdateEngineerPreference(preference models.EngineerPreference) error {
	return s.assignRepo.UpdatePreference(preference)
}

// Private methods

func (s *AssignmentService) findEngineerByRule(inquiry *models.Inquiry, rule models.AssignmentRule) (*models.Account, error) {
	var condition models.RuleCondition
	if err := json.Unmarshal(rule.Conditions, &condition); err != nil {
		return nil, fmt.Errorf("failed to parse rule conditions: %w", err)
	}

	// 檢查產品類別是否匹配
	if len(condition.ProductCategories) > 0 {
		matched := false
		for _, cat := range condition.ProductCategories {
			if cat == inquiry.ProductCategory {
				matched = true
				break
			}
		}
		if !matched {
			return nil, errors.New("product category not matched")
		}
	}

	// 根據規則類型選擇工程師
	switch rule.RuleType {
	case "load_balance":
		return s.findEngineerByLoadBalance(inquiry)
	case "skill_based":
		return s.findEngineerBySkill(inquiry, condition.MinSkillLevel)
	case "rotation":
		return s.findEngineerByRotation(inquiry)
	default:
		return s.findEngineerByLoadBalance(inquiry)
	}
}

func (s *AssignmentService) findEngineerByLoadBalance(inquiry *models.Inquiry) (*models.Account, error) {
	// 獲取所有有能力處理此類產品的工程師
	engineers, err := s.assignRepo.GetCapableEngineers(inquiry.ProductCategory)
	if err != nil {
		return nil, err
	}

	if len(engineers) == 0 {
		return nil, errors.New("no capable engineers found")
	}

	// 獲取工作量並選擇負載最低的
	var selectedEngineer *models.Account
	minWorkload := int(^uint(0) >> 1) // max int

	for _, eng := range engineers {
		workload, err := s.assignRepo.GetEngineerWorkload(eng.ID)
		if err != nil {
			continue
		}

		if workload.CurrentInquiries < minWorkload {
			minWorkload = workload.CurrentInquiries
			selectedEngineer = &eng
		}
	}

	return selectedEngineer, nil
}

func (s *AssignmentService) findEngineerBySkill(inquiry *models.Inquiry, minSkillLevel int) (*models.Account, error) {
	// 獲取符合技能要求的工程師
	engineers, err := s.assignRepo.GetEngineersBySkillLevel(inquiry.ProductCategory, minSkillLevel)
	if err != nil {
		return nil, err
	}

	if len(engineers) == 0 {
		return nil, errors.New("no engineers with required skill level")
	}

	// 在符合技能的工程師中選擇負載最低的
	return s.selectByLowestWorkload(engineers)
}

func (s *AssignmentService) findEngineerByRotation(inquiry *models.Inquiry) (*models.Account, error) {
	// 獲取所有有能力的工程師
	engineers, err := s.assignRepo.GetCapableEngineers(inquiry.ProductCategory)
	if err != nil {
		return nil, err
	}

	if len(engineers) == 0 {
		return nil, errors.New("no capable engineers found")
	}

	// 根據最後分派時間選擇
	var selectedEngineer *models.Account
	var oldestAssignment time.Time

	for _, eng := range engineers {
		workload, err := s.assignRepo.GetEngineerWorkload(eng.ID)
		if err != nil {
			continue
		}

		if workload.LastAssignedAt == nil || workload.LastAssignedAt.Before(oldestAssignment) {
			if workload.LastAssignedAt != nil {
				oldestAssignment = *workload.LastAssignedAt
			}
			selectedEngineer = &eng
		}
	}

	return selectedEngineer, nil
}

func (s *AssignmentService) selectByLowestWorkload(engineers []models.Account) (*models.Account, error) {
	var selectedEngineer *models.Account
	minWorkload := int(^uint(0) >> 1) // max int

	for _, eng := range engineers {
		workload, err := s.assignRepo.GetEngineerWorkload(eng.ID)
		if err != nil {
			continue
		}

		if workload.CurrentInquiries < minWorkload {
			minWorkload = workload.CurrentInquiries
			selectedEngineer = &eng
		}
	}

	return selectedEngineer, nil
}

func (s *AssignmentService) isEngineerCapable(engineerID uuid.UUID, productCategory string) (bool, error) {
	capabilities, err := s.assignRepo.GetEngineerCapabilities(engineerID)
	if err != nil {
		return false, err
	}

	for _, cap := range capabilities {
		if cap.ProductCategory == productCategory && cap.IsActive {
			return true, nil
		}
	}

	return false, nil
}

func (s *AssignmentService) assignToEngineer(inquiry *models.Inquiry, engineerID uuid.UUID, fromEngineer *uuid.UUID, assignmentType, reason string, ruleID *uuid.UUID) (*models.AssignmentHistory, error) {
	// 開始事務
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新詢價單
	inquiry.AssignedEngineerID = &engineerID
	inquiry.AssignedAt = func() *time.Time { t := time.Now(); return &t }()
	inquiry.Status = "assigned"

	if err := tx.Save(inquiry).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update inquiry: %w", err)
	}

	// 建立分派歷史記錄
	history := &models.AssignmentHistory{
		InquiryID:        inquiry.ID,
		AssignedFrom:     fromEngineer,
		AssignedTo:       engineerID,
		AssignmentType:   assignmentType,
		AssignmentReason: reason,
		RuleID:           ruleID,
	}

	if err := tx.Create(history).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create assignment history: %w", err)
	}

	// 提交事務
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return history, nil
}