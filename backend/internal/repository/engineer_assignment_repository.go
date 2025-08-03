package repository

import (
	"fmt"
	"time"

	"github.com/fastenmind/fastener-api/internal/models"
	"gorm.io/gorm"
)

type EngineerAssignmentRepository struct {
	db *gorm.DB
}

func NewEngineerAssignmentRepository(db *gorm.DB) *EngineerAssignmentRepository {
	return &EngineerAssignmentRepository{db: db}
}

// Create 創建分派記錄
func (r *EngineerAssignmentRepository) Create(assignment *models.EngineerAssignment) error {
	return r.db.Create(assignment).Error
}

// Update 更新分派記錄
func (r *EngineerAssignmentRepository) Update(assignment *models.EngineerAssignment) error {
	return r.db.Save(assignment).Error
}

// GetByID 根據ID獲取分派記錄
func (r *EngineerAssignmentRepository) GetByID(id, companyID string) (*models.EngineerAssignment, error) {
	var assignment models.EngineerAssignment
	err := r.db.Where("id = ? AND company_id = ?", id, companyID).First(&assignment).Error
	if err != nil {
		return nil, err
	}
	return &assignment, nil
}

// GetActiveAssignment 獲取詢價單的活躍分派
func (r *EngineerAssignmentRepository) GetActiveAssignment(inquiryID, companyID string) (*models.EngineerAssignment, error) {
	var assignment models.EngineerAssignment
	err := r.db.Where("inquiry_id = ? AND company_id = ? AND status IN ?", 
		inquiryID, companyID, []string{"pending", "in_progress"}).First(&assignment).Error
	if err != nil {
		return nil, err
	}
	return &assignment, nil
}

// GetEngineerActiveAssignments 獲取工程師的活躍分派數量
func (r *EngineerAssignmentRepository) GetEngineerActiveAssignments(engineerID, companyID string) (int, error) {
	var count int64
	err := r.db.Model(&models.EngineerAssignment{}).
		Where("engineer_id = ? AND company_id = ? AND status IN ?", 
			engineerID, companyID, []string{"pending", "in_progress"}).
		Count(&count).Error
	return int(count), err
}

// GetEngineersWorkload 獲取所有工程師的工作負載
func (r *EngineerAssignmentRepository) GetEngineersWorkload(companyID string) (map[string]int, error) {
	type result struct {
		EngineerID string
		Count      int
	}
	
	var results []result
	err := r.db.Model(&models.EngineerAssignment{}).
		Select("engineer_id, COUNT(*) as count").
		Where("company_id = ? AND status IN ?", companyID, []string{"pending", "in_progress"}).
		Group("engineer_id").
		Scan(&results).Error
	
	if err != nil {
		return nil, err
	}
	
	workloads := make(map[string]int)
	for _, r := range results {
		workloads[r.EngineerID] = r.Count
	}
	
	return workloads, nil
}

// GetEngineerAssignmentsByPeriod 獲取工程師在指定期間的分派
func (r *EngineerAssignmentRepository) GetEngineerAssignmentsByPeriod(engineerID, companyID, startDate, endDate string) ([]models.EngineerAssignment, error) {
	var assignments []models.EngineerAssignment
	
	query := r.db.Where("engineer_id = ? AND company_id = ?", engineerID, companyID)
	
	if startDate != "" {
		query = query.Where("assigned_at >= ?", startDate)
	}
	
	if endDate != "" {
		query = query.Where("assigned_at <= ?", endDate)
	}
	
	err := query.Find(&assignments).Error
	return assignments, err
}

// CreateHistory 創建分派歷史記錄
func (r *EngineerAssignmentRepository) CreateHistory(history *models.EngineerAssignmentHistory) error {
	return r.db.Create(history).Error
}

// GetHistory 獲取分派歷史
func (r *EngineerAssignmentRepository) GetHistory(companyID, inquiryID, engineerID string, offset, limit int) ([]models.EngineerAssignmentHistory, int64, error) {
	var histories []models.EngineerAssignmentHistory
	var total int64
	
	query := r.db.Model(&models.EngineerAssignmentHistory{}).
		Joins("JOIN engineer_assignments ON assignment_histories.assignment_id = engineer_assignments.id").
		Where("engineer_assignments.company_id = ?", companyID)
	
	if inquiryID != "" {
		query = query.Where("engineer_assignments.inquiry_id = ?", inquiryID)
	}
	
	if engineerID != "" {
		query = query.Where("(assignment_histories.from_engineer = ? OR assignment_histories.to_engineer = ?)", 
			engineerID, engineerID)
	}
	
	// 計算總數
	query.Count(&total)
	
	// 獲取分頁數據
	err := query.
		Select("assignment_histories.*").
		Order("assignment_histories.action_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&histories).Error
	
	return histories, total, err
}

// GetStats 獲取分派統計數據
func (r *EngineerAssignmentRepository) GetStats(companyID, period string) (*models.AssignmentStats, error) {
	stats := &models.AssignmentStats{
		Period: period,
	}
	
	// 根據期間計算時間範圍
	var startDate time.Time
	now := time.Now()
	
	switch period {
	case "daily":
		startDate = now.AddDate(0, 0, -30) // 最近30天
	case "weekly":
		startDate = now.AddDate(0, 0, -7*12) // 最近12週
	case "monthly":
		startDate = now.AddDate(0, -12, 0) // 最近12個月
	default:
		startDate = now.AddDate(0, -12, 0)
	}
	
	// 總分派數
	r.db.Model(&models.EngineerAssignment{}).
		Where("company_id = ? AND assigned_at >= ?", companyID, startDate).
		Count(&stats.TotalAssignments)
	
	// 各狀態統計
	var statusStats []struct {
		Status string
		Count  int64
	}
	
	r.db.Model(&models.EngineerAssignment{}).
		Select("status, COUNT(*) as count").
		Where("company_id = ? AND assigned_at >= ?", companyID, startDate).
		Group("status").
		Scan(&statusStats)
	
	stats.StatusBreakdown = make(map[string]int64)
	for _, s := range statusStats {
		stats.StatusBreakdown[s.Status] = s.Count
	}
	
	// 平均完成時間
	var avgCompletion struct {
		AvgHours float64
	}
	
	r.db.Model(&models.EngineerAssignment{}).
		Select("AVG(EXTRACT(EPOCH FROM (completed_at - assigned_at))/3600) as avg_hours").
		Where("company_id = ? AND status = ? AND completed_at IS NOT NULL AND assigned_at >= ?", 
			companyID, "completed", startDate).
		Scan(&avgCompletion)
	
	stats.AvgCompletionTime = avgCompletion.AvgHours
	
	// 準時完成率
	var onTimeStats struct {
		Total   int64
		OnTime  int64
	}
	
	r.db.Model(&models.EngineerAssignment{}).
		Where("company_id = ? AND status = ? AND due_date IS NOT NULL AND completed_at IS NOT NULL AND assigned_at >= ?",
			companyID, "completed", startDate).
		Count(&onTimeStats.Total)
	
	r.db.Model(&models.EngineerAssignment{}).
		Where("company_id = ? AND status = ? AND due_date IS NOT NULL AND completed_at IS NOT NULL AND completed_at <= due_date AND assigned_at >= ?",
			companyID, "completed", startDate).
		Count(&onTimeStats.OnTime)
	
	if onTimeStats.Total > 0 {
		stats.OnTimeRate = float64(onTimeStats.OnTime) / float64(onTimeStats.Total)
	}
	
	// 各工程師統計
	type engineerStat struct {
		EngineerID       string
		EngineerName     string
		TotalAssigned    int64
		Completed        int64
		AvgCompletionTime float64
	}
	
	var engineerStats []engineerStat
	
	r.db.Table("engineer_assignments").
		Select(`
			engineer_assignments.engineer_id,
			accounts.name as engineer_name,
			COUNT(*) as total_assigned,
			COUNT(CASE WHEN engineer_assignments.status = 'completed' THEN 1 END) as completed,
			AVG(CASE 
				WHEN engineer_assignments.status = 'completed' AND completed_at IS NOT NULL 
				THEN EXTRACT(EPOCH FROM (completed_at - assigned_at))/3600 
			END) as avg_completion_time
		`).
		Joins("JOIN accounts ON engineer_assignments.engineer_id = accounts.id").
		Where("engineer_assignments.company_id = ? AND engineer_assignments.assigned_at >= ?", companyID, startDate).
		Group("engineer_assignments.engineer_id, accounts.name").
		Scan(&engineerStats)
	
	stats.EngineerStats = make([]models.EngineerStat, len(engineerStats))
	for i, es := range engineerStats {
		stats.EngineerStats[i] = models.EngineerStat{
			EngineerID:        es.EngineerID,
			EngineerName:      es.EngineerName,
			TotalAssigned:     es.TotalAssigned,
			Completed:         es.Completed,
			AvgCompletionTime: es.AvgCompletionTime,
			CompletionRate:    float64(es.Completed) / float64(es.TotalAssigned),
		}
	}
	
	// 時間序列數據
	var timeSeriesFormat string
	switch period {
	case "daily":
		timeSeriesFormat = "YYYY-MM-DD"
	case "weekly":
		timeSeriesFormat = "YYYY-IW"
	case "monthly":
		timeSeriesFormat = "YYYY-MM"
	}
	
	type timeSeries struct {
		Period    string
		Assigned  int64
		Completed int64
	}
	
	var timeSeriesData []timeSeries
	
	r.db.Table("engineer_assignments").
		Select(fmt.Sprintf(`
			TO_CHAR(assigned_at, '%s') as period,
			COUNT(*) as assigned,
			COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed
		`, timeSeriesFormat)).
		Where("company_id = ? AND assigned_at >= ?", companyID, startDate).
		Group("period").
		Order("period").
		Scan(&timeSeriesData)
	
	stats.TimeSeries = make([]models.TimeSeriesData, len(timeSeriesData))
	for i, ts := range timeSeriesData {
		stats.TimeSeries[i] = models.TimeSeriesData{
			Period:    ts.Period,
			Assigned:  ts.Assigned,
			Completed: ts.Completed,
		}
	}
	
	return stats, nil
}