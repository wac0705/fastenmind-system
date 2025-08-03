package service

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/fastenmind/fastener-api/internal/model"
	"github.com/fastenmind/fastener-api/internal/models"
	"github.com/fastenmind/fastener-api/internal/repository"

	"github.com/google/uuid"
)

type EngineerAssignmentService struct {
	assignmentRepo *repository.EngineerAssignmentRepository
	inquiryRepo    repository.InquiryRepository
	accountRepo    repository.AccountRepository
}

func NewEngineerAssignmentService(
	assignmentRepo *repository.EngineerAssignmentRepository,
	inquiryRepo repository.InquiryRepository,
	accountRepo repository.AccountRepository,
) *EngineerAssignmentService {
	return &EngineerAssignmentService{
		assignmentRepo: assignmentRepo,
		inquiryRepo:    inquiryRepo,
		accountRepo:    accountRepo,
	}
}

// GetAvailableEngineers 獲取可用的工程師列表
func (s *EngineerAssignmentService) GetAvailableEngineers(companyID, inquiryID string) ([]models.EngineerAvailability, error) {
	// 獲取所有工程師
	pagination := &model.Pagination{Page: 1, PageSize: 100}
	accounts, err := s.accountRepo.List(context.Background(), uuid.MustParse(companyID), pagination)
	if err != nil {
		return nil, err
	}
	
	// 過濾出工程師
	var engineers []*model.Account
	for _, acc := range accounts {
		if acc.Role == "engineer" {
			engineers = append(engineers, acc)
		}
	}
	
	// 獲取工程師當前的工作負載
	workloads, err := s.assignmentRepo.GetEngineersWorkload(companyID)
	if err != nil {
		return nil, err
	}
	
	// 如果提供了詢價單ID，獲取詢價單詳情以匹配專長
	var inquiry *model.Inquiry
	if inquiryID != "" {
		id, err := uuid.Parse(inquiryID)
		if err != nil {
			return nil, err
		}
		inquiry, err = s.inquiryRepo.Get(id)
		if err != nil {
			return nil, err
		}
	}
	
	// 構建可用工程師列表
	availableEngineers := make([]models.EngineerAvailability, 0)
	for _, engineer := range engineers {
		maxAssignments := engineer.MaxAssignments
		if maxAssignments == 0 {
			maxAssignments = 10 // 默認最大分派數
		}
		
		currentLoad := workloads[engineer.ID.String()]
		availability := models.EngineerAvailability{
			EngineerID:   engineer.ID.String(),
			EngineerName: engineer.FullName,
			Department:   engineer.Department,
			Expertise:    engineer.Expertise,
			CurrentLoad:  currentLoad,
			MaxLoad:      maxAssignments,
			IsAvailable:  currentLoad < maxAssignments,
		}
		
		// 計算專長匹配度
		if inquiry != nil && len(engineer.Expertise) > 0 {
			availability.ExpertiseMatch = s.calculateExpertiseMatch(engineer.Expertise, inquiry.ProductCategory)
		}
		
		availableEngineers = append(availableEngineers, availability)
	}
	
	// 按照可用性和匹配度排序
	sort.Slice(availableEngineers, func(i, j int) bool {
		if availableEngineers[i].IsAvailable != availableEngineers[j].IsAvailable {
			return availableEngineers[i].IsAvailable
		}
		if availableEngineers[i].ExpertiseMatch != availableEngineers[j].ExpertiseMatch {
			return availableEngineers[i].ExpertiseMatch > availableEngineers[j].ExpertiseMatch
		}
		return availableEngineers[i].CurrentLoad < availableEngineers[j].CurrentLoad
	})
	
	return availableEngineers, nil
}

// AssignEngineer 分配工程師到詢價單
func (s *EngineerAssignmentService) AssignEngineer(companyID, inquiryID, engineerID, assignedBy, priority, dueDate, notes string) (*models.EngineerAssignment, error) {
	// 檢查詢價單是否存在
	id, err := uuid.Parse(inquiryID)
	if err != nil {
		return nil, err
	}
	inquiry, err := s.inquiryRepo.Get(id)
	if err != nil {
		return nil, err
	}
	
	// 檢查是否已經有分派
	existingAssignment, _ := s.assignmentRepo.GetActiveAssignment(inquiryID, companyID)
	if existingAssignment != nil {
		return nil, errors.New("inquiry already has an active assignment")
	}
	
	// 檢查工程師是否可用
	engineer, err := s.accountRepo.GetByID(context.Background(), uuid.MustParse(engineerID))
	if err != nil {
		return nil, err
	}
	
	if engineer.Role != "engineer" {
		return nil, errors.New("selected user is not an engineer")
	}
	
	// 檢查工程師工作負載
	workload, err := s.assignmentRepo.GetEngineerActiveAssignments(engineerID, companyID)
	if err != nil {
		return nil, err
	}
	
	if workload >= engineer.MaxAssignments {
		return nil, errors.New("engineer has reached maximum assignment limit")
	}
	
	// 創建分派記錄
	assignment := &models.EngineerAssignment{
		ID:         uuid.New().String(),
		CompanyID:  companyID,
		InquiryID:  inquiryID,
		EngineerID: engineerID,
		AssignedBy: assignedBy,
		AssignedAt: time.Now(),
		Status:     "pending",
		Priority:   priority,
		Notes:      notes,
		CreatedBy:  assignedBy,
		UpdatedBy:  assignedBy,
	}
	
	if dueDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", dueDate); err == nil {
			assignment.DueDate = &parsedDate
		}
	}
	
	// 保存分派記錄
	if err := s.assignmentRepo.Create(assignment); err != nil {
		return nil, err
	}
	
	// 更新詢價單狀態
	engineerUUID, err := uuid.Parse(engineerID)
	if err != nil {
		return nil, err
	}
	assignedByUUID, err := uuid.Parse(assignedBy)
	if err != nil {
		return nil, err
	}
	inquiry.Status = "assigned"
	inquiry.AssignedEngineerID = &engineerUUID
	inquiry.UpdatedBy = &assignedByUUID
	if err := s.inquiryRepo.Update(inquiry); err != nil {
		return nil, err
	}
	
	// 記錄分派歷史
	history := &models.EngineerAssignmentHistory{
		ID:           uuid.New().String(),
		AssignmentID: assignment.ID,
		Action:       "assigned",
		FromEngineer: "",
		ToEngineer:   engineerID,
		ActionBy:     assignedBy,
		ActionAt:     time.Now(),
		Reason:       notes,
	}
	
	if err := s.assignmentRepo.CreateHistory(history); err != nil {
		// 記錄失敗不影響主流程
		fmt.Printf("Failed to create assignment history: %v\n", err)
	}
	
	return assignment, nil
}

// ReassignEngineer 重新分配工程師
func (s *EngineerAssignmentService) ReassignEngineer(companyID, assignmentID, newEngineerID, reassignedBy, reason string) (*models.EngineerAssignment, error) {
	// 獲取現有分派記錄
	assignment, err := s.assignmentRepo.GetByID(assignmentID, companyID)
	if err != nil {
		return nil, err
	}
	
	if assignment.Status == "completed" || assignment.Status == "cancelled" {
		return nil, errors.New("cannot reassign completed or cancelled assignment")
	}
	
	// 檢查新工程師是否可用
	newEngineer, err := s.accountRepo.GetByID(context.Background(), uuid.MustParse(newEngineerID))
	if err != nil {
		return nil, err
	}
	
	if newEngineer.Role != "engineer" {
		return nil, errors.New("selected user is not an engineer")
	}
	
	// 記錄重新分派歷史
	history := &models.EngineerAssignmentHistory{
		ID:           uuid.New().String(),
		AssignmentID: assignment.ID,
		Action:       "reassigned",
		FromEngineer: assignment.EngineerID,
		ToEngineer:   newEngineerID,
		ActionBy:     reassignedBy,
		ActionAt:     time.Now(),
		Reason:       reason,
	}
	
	// 更新分派記錄
	assignment.EngineerID = newEngineerID
	assignment.UpdatedBy = reassignedBy
	assignment.UpdatedAt = time.Now()
	
	if err := s.assignmentRepo.Update(assignment); err != nil {
		return nil, err
	}
	
	// 更新詢價單
	inquiryID, err := uuid.Parse(assignment.InquiryID)
	if err != nil {
		return nil, err
	}
	inquiry, err := s.inquiryRepo.Get(inquiryID)
	if err != nil {
		return nil, err
	}
	
	newEngineerUUID, err := uuid.Parse(newEngineerID)
	if err != nil {
		return nil, err
	}
	reassignedByUUID, err := uuid.Parse(reassignedBy)
	if err != nil {
		return nil, err
	}
	inquiry.AssignedEngineerID = &newEngineerUUID
	inquiry.UpdatedBy = &reassignedByUUID
	if err := s.inquiryRepo.Update(inquiry); err != nil {
		return nil, err
	}
	
	// 保存歷史記錄
	if err := s.assignmentRepo.CreateHistory(history); err != nil {
		fmt.Printf("Failed to create reassignment history: %v\n", err)
	}
	
	return assignment, nil
}

// GetAssignmentHistory 獲取分派歷史
func (s *EngineerAssignmentService) GetAssignmentHistory(companyID, inquiryID, engineerID string, page, limit int) ([]models.EngineerAssignmentHistory, int64, error) {
	offset := (page - 1) * limit
	return s.assignmentRepo.GetHistory(companyID, inquiryID, engineerID, offset, limit)
}

// GetEngineerWorkload 獲取工程師工作負載
func (s *EngineerAssignmentService) GetEngineerWorkload(companyID, startDate, endDate string) ([]models.EngineerWorkloadSummary, error) {
	pagination := &model.Pagination{Page: 1, PageSize: 100}
	accounts, err := s.accountRepo.List(context.Background(), uuid.MustParse(companyID), pagination)
	if err != nil {
		return nil, err
	}
	
	// 過濾出工程師
	var engineers []*model.Account
	for _, acc := range accounts {
		if acc.Role == "engineer" {
			engineers = append(engineers, acc)
		}
	}
	
	workloads := make([]models.EngineerWorkloadSummary, 0)
	
	for _, engineer := range engineers {
		// 獲取該期間內的分派數據
		assignments, err := s.assignmentRepo.GetEngineerAssignmentsByPeriod(engineer.ID.String(), companyID, startDate, endDate)
		if err != nil {
			return nil, err
		}
		
		// 統計各狀態的分派數量
		statusCount := make(map[string]int)
		var totalAssignments int
		var completedOnTime int
		var overdue int
		
		for _, assignment := range assignments {
			statusCount[assignment.Status]++
			totalAssignments++
			
			if assignment.Status == "completed" {
				if assignment.DueDate != nil && assignment.CompletedAt != nil {
					if assignment.CompletedAt.Before(*assignment.DueDate) || assignment.CompletedAt.Equal(*assignment.DueDate) {
						completedOnTime++
					}
				}
			}
			
			if assignment.Status == "pending" || assignment.Status == "in_progress" {
				if assignment.DueDate != nil && time.Now().After(*assignment.DueDate) {
					overdue++
				}
			}
		}
		
		workload := models.EngineerWorkloadSummary{
			EngineerID:       engineer.ID.String(),
			EngineerName:     engineer.FullName,
			TotalAssignments: totalAssignments,
			Pending:          statusCount["pending"],
			InProgress:       statusCount["in_progress"],
			Completed:        statusCount["completed"],
			CompletedOnTime:  completedOnTime,
			Overdue:          overdue,
			AvgCompletionTime: s.calculateAvgCompletionTime(assignments),
		}
		
		workloads = append(workloads, workload)
	}
	
	return workloads, nil
}

// UpdateAssignmentStatus 更新分派狀態
func (s *EngineerAssignmentService) UpdateAssignmentStatus(companyID, assignmentID, status, updatedBy, notes string) (*models.EngineerAssignment, error) {
	assignment, err := s.assignmentRepo.GetByID(assignmentID, companyID)
	if err != nil {
		return nil, err
	}
	
	// 記錄狀態變更歷史
	history := &models.EngineerAssignmentHistory{
		ID:           uuid.New().String(),
		AssignmentID: assignment.ID,
		Action:       "status_changed",
		FromStatus:   assignment.Status,
		ToStatus:     status,
		ActionBy:     updatedBy,
		ActionAt:     time.Now(),
		Reason:       notes,
	}
	
	// 更新狀態
	assignment.Status = status
	assignment.UpdatedBy = updatedBy
	assignment.UpdatedAt = time.Now()
	
	if status == "completed" {
		now := time.Now()
		assignment.CompletedAt = &now
	}
	
	if err := s.assignmentRepo.Update(assignment); err != nil {
		return nil, err
	}
	
	// 如果完成或取消，更新詢價單狀態
	if status == "completed" || status == "cancelled" {
		inquiryID, err := uuid.Parse(assignment.InquiryID)
		if err == nil {
			inquiry, err := s.inquiryRepo.Get(inquiryID)
			if err == nil {
				if status == "completed" {
					inquiry.Status = "quoted"
				} else {
					inquiry.Status = "cancelled"
				}
				updatedByUUID, _ := uuid.Parse(updatedBy)
				inquiry.UpdatedBy = &updatedByUUID
				s.inquiryRepo.Update(inquiry)
			}
		}
	}
	
	// 保存歷史記錄
	if err := s.assignmentRepo.CreateHistory(history); err != nil {
		fmt.Printf("Failed to create status change history: %v\n", err)
	}
	
	return assignment, nil
}

// GetAssignmentStats 獲取分派統計數據
func (s *EngineerAssignmentService) GetAssignmentStats(companyID, period string) (*models.AssignmentStats, error) {
	return s.assignmentRepo.GetStats(companyID, period)
}

// AutoAssignEngineer 自動分配工程師
func (s *EngineerAssignmentService) AutoAssignEngineer(companyID, inquiryID, assignedBy string, rules interface{}) (*models.EngineerAssignment, error) {
	// 獲取詢價單詳情
	id, err := uuid.Parse(inquiryID)
	if err != nil {
		return nil, err
	}
	inquiry, err := s.inquiryRepo.Get(id)
	if err != nil {
		return nil, err
	}
	
	// 獲取可用工程師列表
	availableEngineers, err := s.GetAvailableEngineers(companyID, inquiryID)
	if err != nil {
		return nil, err
	}
	
	// 過濾出真正可用的工程師
	var candidates []models.EngineerAvailability
	for _, engineer := range availableEngineers {
		if engineer.IsAvailable {
			candidates = append(candidates, engineer)
		}
	}
	
	if len(candidates) == 0 {
		return nil, errors.New("no available engineers found")
	}
	
	// 選擇最佳工程師（基於專長匹配度和當前負載）
	bestEngineer := candidates[0]
	for _, candidate := range candidates[1:] {
		// 優先考慮專長匹配度
		if candidate.ExpertiseMatch > bestEngineer.ExpertiseMatch {
			bestEngineer = candidate
		} else if candidate.ExpertiseMatch == bestEngineer.ExpertiseMatch {
			// 專長相同時，選擇負載較低的
			if candidate.CurrentLoad < bestEngineer.CurrentLoad {
				bestEngineer = candidate
			}
		}
	}
	
	// 執行分派
	priority := "normal"
	// 根據所需日期判斷優先級
	if time.Now().After(inquiry.RequiredDate.AddDate(0, 0, -7)) {
		priority = "high"
	}
	
	notes := fmt.Sprintf("Auto-assigned based on expertise match: %.0f%%, current load: %d/%d",
		bestEngineer.ExpertiseMatch*100,
		bestEngineer.CurrentLoad,
		bestEngineer.MaxLoad)
	
	return s.AssignEngineer(companyID, inquiryID, bestEngineer.EngineerID, assignedBy, priority, "", notes)
}

// calculateExpertiseMatch 計算專長匹配度
func (s *EngineerAssignmentService) calculateExpertiseMatch(expertise []string, productCategory string) float64 {
	if len(expertise) == 0 || productCategory == "" {
		return 0
	}
	
	// 簡單的匹配算法：檢查專長中是否包含產品類別關鍵字
	for _, exp := range expertise {
		if exp == productCategory {
			return 1.0 // 完全匹配
		}
		// 可以加入更複雜的相似度計算
	}
	
	return 0.5 // 默認匹配度
}

// calculateAvgCompletionTime 計算平均完成時間
func (s *EngineerAssignmentService) calculateAvgCompletionTime(assignments []models.EngineerAssignment) float64 {
	var totalHours float64
	var completedCount int
	
	for _, assignment := range assignments {
		if assignment.Status == "completed" && assignment.CompletedAt != nil {
			duration := assignment.CompletedAt.Sub(assignment.AssignedAt)
			totalHours += duration.Hours()
			completedCount++
		}
	}
	
	if completedCount == 0 {
		return 0
	}
	
	return totalHours / float64(completedCount)
}