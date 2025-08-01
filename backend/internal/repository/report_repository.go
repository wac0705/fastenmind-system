package repository

import (
	"github.com/fastenmind/fastener-api/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReportRepository interface {
	// Report operations
	CreateReport(report *models.Report) error
	UpdateReport(report *models.Report) error
	GetReport(id uuid.UUID) (*models.Report, error)
	ListReports(companyID uuid.UUID, params map[string]interface{}) ([]models.Report, int64, error)
	DeleteReport(id uuid.UUID) error
	
	// Report Template operations
	CreateReportTemplate(template *models.ReportTemplate) error
	UpdateReportTemplate(template *models.ReportTemplate) error
	GetReportTemplate(id uuid.UUID) (*models.ReportTemplate, error)
	ListReportTemplates(companyID *uuid.UUID, params map[string]interface{}) ([]models.ReportTemplate, int64, error)
	DeleteReportTemplate(id uuid.UUID) error
	
	// Report Execution operations
	CreateReportExecution(execution *models.ReportExecution) error
	UpdateReportExecution(execution *models.ReportExecution) error
	GetReportExecution(id uuid.UUID) (*models.ReportExecution, error)
	ListReportExecutions(reportID uuid.UUID, params map[string]interface{}) ([]models.ReportExecution, int64, error)
	DeleteReportExecution(id uuid.UUID) error
	
	// Report Subscription operations
	CreateReportSubscription(subscription *models.ReportSubscription) error
	UpdateReportSubscription(subscription *models.ReportSubscription) error
	GetReportSubscription(id uuid.UUID) (*models.ReportSubscription, error)
	ListReportSubscriptions(userID uuid.UUID, params map[string]interface{}) ([]models.ReportSubscription, int64, error)
	DeleteReportSubscription(id uuid.UUID) error
	
	// Report Dashboard operations
	CreateReportDashboard(dashboard *models.ReportDashboard) error
	UpdateReportDashboard(dashboard *models.ReportDashboard) error
	GetReportDashboard(id uuid.UUID) (*models.ReportDashboard, error)
	ListReportDashboards(companyID uuid.UUID, params map[string]interface{}) ([]models.ReportDashboard, int64, error)
	DeleteReportDashboard(id uuid.UUID) error
	
	// Report Data Source operations
	CreateReportDataSource(dataSource *models.ReportDataSource) error
	UpdateReportDataSource(dataSource *models.ReportDataSource) error
	GetReportDataSource(id uuid.UUID) (*models.ReportDataSource, error)
	ListReportDataSources(companyID uuid.UUID, params map[string]interface{}) ([]models.ReportDataSource, int64, error)
	DeleteReportDataSource(id uuid.UUID) error
	
	// Report Schedule operations
	CreateReportSchedule(schedule *models.ReportSchedule) error
	UpdateReportSchedule(schedule *models.ReportSchedule) error
	GetReportSchedule(id uuid.UUID) (*models.ReportSchedule, error)
	ListReportSchedules(companyID uuid.UUID, params map[string]interface{}) ([]models.ReportSchedule, int64, error)
	DeleteReportSchedule(id uuid.UUID) error
	
	// Business KPI operations
	CreateBusinessKPI(kpi *models.BusinessKPI) error
	UpdateBusinessKPI(kpi *models.BusinessKPI) error
	GetBusinessKPI(id uuid.UUID) (*models.BusinessKPI, error)
	ListBusinessKPIs(companyID uuid.UUID, params map[string]interface{}) ([]models.BusinessKPI, int64, error)
	DeleteBusinessKPI(id uuid.UUID) error
	
	// Business operations
	GetReportStatistics(companyID uuid.UUID) (map[string]interface{}, error)
	GetPopularReports(companyID uuid.UUID, limit int) ([]models.Report, error)
	GetRecentExecutions(companyID uuid.UUID, limit int) ([]models.ReportExecution, error)
	GetScheduledReports(companyID uuid.UUID) ([]models.ReportSchedule, error)
}

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db interface{}) ReportRepository {
	gormDB, ok := db.(*gorm.DB)
	if !ok {
		panic("invalid database type, expected *gorm.DB")
	}
	return &reportRepository{db: gormDB}
}

// Report operations
func (r *reportRepository) CreateReport(report *models.Report) error {
	return r.db.Create(report).Error
}

func (r *reportRepository) UpdateReport(report *models.Report) error {
	return r.db.Save(report).Error
}

func (r *reportRepository) GetReport(id uuid.UUID) (*models.Report, error) {
	var report models.Report
	err := r.db.Preload("Company").
		Preload("Template").
		Preload("Creator").
		Preload("UpdatedByUser").
		First(&report, id).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *reportRepository) ListReports(companyID uuid.UUID, params map[string]interface{}) ([]models.Report, int64, error) {
	var reports []models.Report
	var total int64

	query := r.db.Model(&models.Report{}).Where("company_id = ?", companyID)

	// Apply filters
	if category, ok := params["category"].(string); ok && category != "" {
		query = query.Where("category = ?", category)
	}

	if reportType, ok := params["type"].(string); ok && reportType != "" {
		query = query.Where("type = ?", reportType)
	}

	if status, ok := params["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}

	if isScheduled, ok := params["is_scheduled"].(bool); ok {
		query = query.Where("is_scheduled = ?", isScheduled)
	}

	if isPublic, ok := params["is_public"].(bool); ok {
		query = query.Where("is_public = ?", isPublic)
	}

	if search, ok := params["search"].(string); ok && search != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", 
			"%"+search+"%", "%"+search+"%")
	}

	if createdBy, ok := params["created_by"].(string); ok && createdBy != "" {
		query = query.Where("created_by = ?", createdBy)
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

	// Apply sorting
	sortBy := "created_at"
	if sb, ok := params["sort_by"].(string); ok && sb != "" {
		sortBy = sb
	}
	sortOrder := "DESC"
	if so, ok := params["sort_order"].(string); ok && so != "" {
		sortOrder = so
	}
	query = query.Order(sortBy + " " + sortOrder)

	// Load with relations
	if err := query.
		Preload("Creator").
		Preload("Template").
		Find(&reports).Error; err != nil {
		return nil, 0, err
	}

	return reports, total, nil
}

func (r *reportRepository) DeleteReport(id uuid.UUID) error {
	return r.db.Delete(&models.Report{}, id).Error
}

// Report Template operations
func (r *reportRepository) CreateReportTemplate(template *models.ReportTemplate) error {
	return r.db.Create(template).Error
}

func (r *reportRepository) UpdateReportTemplate(template *models.ReportTemplate) error {
	return r.db.Save(template).Error
}

func (r *reportRepository) GetReportTemplate(id uuid.UUID) (*models.ReportTemplate, error) {
	var template models.ReportTemplate
	err := r.db.Preload("Company").
		Preload("Creator").
		First(&template, id).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *reportRepository) ListReportTemplates(companyID *uuid.UUID, params map[string]interface{}) ([]models.ReportTemplate, int64, error) {
	var templates []models.ReportTemplate
	var total int64

	query := r.db.Model(&models.ReportTemplate{})
	
	// Filter by company or system templates
	if companyID != nil {
		query = query.Where("company_id = ? OR is_system_template = ?", *companyID, true)
	} else {
		query = query.Where("is_system_template = ?", true)
	}

	// Apply filters
	if category, ok := params["category"].(string); ok && category != "" {
		query = query.Where("category = ?", category)
	}

	if templateType, ok := params["type"].(string); ok && templateType != "" {
		query = query.Where("type = ?", templateType)
	}

	if isSystemTemplate, ok := params["is_system_template"].(bool); ok {
		query = query.Where("is_system_template = ?", isSystemTemplate)
	}

	if industry, ok := params["industry"].(string); ok && industry != "" {
		query = query.Where("industry = ?", industry)
	}

	if language, ok := params["language"].(string); ok && language != "" {
		query = query.Where("language = ?", language)
	}

	if search, ok := params["search"].(string); ok && search != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", 
			"%"+search+"%", "%"+search+"%")
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

	// Apply sorting
	sortBy := "usage_count"
	if sb, ok := params["sort_by"].(string); ok && sb != "" {
		sortBy = sb
	}
	sortOrder := "DESC"
	if so, ok := params["sort_order"].(string); ok && so != "" {
		sortOrder = so
	}
	query = query.Order(sortBy + " " + sortOrder)

	// Load with relations
	if err := query.
		Preload("Creator").
		Find(&templates).Error; err != nil {
		return nil, 0, err
	}

	return templates, total, nil
}

func (r *reportRepository) DeleteReportTemplate(id uuid.UUID) error {
	return r.db.Delete(&models.ReportTemplate{}, id).Error
}

// Report Execution operations
func (r *reportRepository) CreateReportExecution(execution *models.ReportExecution) error {
	return r.db.Create(execution).Error
}

func (r *reportRepository) UpdateReportExecution(execution *models.ReportExecution) error {
	return r.db.Save(execution).Error
}

func (r *reportRepository) GetReportExecution(id uuid.UUID) (*models.ReportExecution, error) {
	var execution models.ReportExecution
	err := r.db.Preload("Report").
		Preload("ExecutedByUser").
		First(&execution, id).Error
	if err != nil {
		return nil, err
	}
	return &execution, nil
}

func (r *reportRepository) ListReportExecutions(reportID uuid.UUID, params map[string]interface{}) ([]models.ReportExecution, int64, error) {
	var executions []models.ReportExecution
	var total int64

	query := r.db.Model(&models.ReportExecution{}).Where("report_id = ?", reportID)

	// Apply filters
	if status, ok := params["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}

	if isScheduled, ok := params["is_scheduled"].(bool); ok {
		query = query.Where("is_scheduled = ?", isScheduled)
	}

	if triggerType, ok := params["trigger_type"].(string); ok && triggerType != "" {
		query = query.Where("trigger_type = ?", triggerType)
	}

	if executedBy, ok := params["executed_by"].(string); ok && executedBy != "" {
		query = query.Where("executed_by = ?", executedBy)
	}

	if startDate, ok := params["start_date"].(string); ok && startDate != "" {
		query = query.Where("started_at >= ?", startDate)
	}

	if endDate, ok := params["end_date"].(string); ok && endDate != "" {
		query = query.Where("started_at <= ?", endDate)
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

	// Apply sorting
	query = query.Order("started_at DESC")

	// Load with relations
	if err := query.
		Preload("ExecutedByUser").
		Find(&executions).Error; err != nil {
		return nil, 0, err
	}

	return executions, total, nil
}

func (r *reportRepository) DeleteReportExecution(id uuid.UUID) error {
	return r.db.Delete(&models.ReportExecution{}, id).Error
}

// Report Subscription operations
func (r *reportRepository) CreateReportSubscription(subscription *models.ReportSubscription) error {
	return r.db.Create(subscription).Error
}

func (r *reportRepository) UpdateReportSubscription(subscription *models.ReportSubscription) error {
	return r.db.Save(subscription).Error
}

func (r *reportRepository) GetReportSubscription(id uuid.UUID) (*models.ReportSubscription, error) {
	var subscription models.ReportSubscription
	err := r.db.Preload("Report").
		Preload("User").
		First(&subscription, id).Error
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

func (r *reportRepository) ListReportSubscriptions(userID uuid.UUID, params map[string]interface{}) ([]models.ReportSubscription, int64, error) {
	var subscriptions []models.ReportSubscription
	var total int64

	query := r.db.Model(&models.ReportSubscription{}).Where("user_id = ?", userID)

	// Apply filters
	if isActive, ok := params["is_active"].(bool); ok {
		query = query.Where("is_active = ?", isActive)
	}

	if deliveryMethod, ok := params["delivery_method"].(string); ok && deliveryMethod != "" {
		query = query.Where("delivery_method = ?", deliveryMethod)
	}

	if reportID, ok := params["report_id"].(string); ok && reportID != "" {
		query = query.Where("report_id = ?", reportID)
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

	// Apply sorting
	query = query.Order("created_at DESC")

	// Load with relations
	if err := query.
		Preload("Report").
		Find(&subscriptions).Error; err != nil {
		return nil, 0, err
	}

	return subscriptions, total, nil
}

func (r *reportRepository) DeleteReportSubscription(id uuid.UUID) error {
	return r.db.Delete(&models.ReportSubscription{}, id).Error
}

// Report Dashboard operations
func (r *reportRepository) CreateReportDashboard(dashboard *models.ReportDashboard) error {
	return r.db.Create(dashboard).Error
}

func (r *reportRepository) UpdateReportDashboard(dashboard *models.ReportDashboard) error {
	return r.db.Save(dashboard).Error
}

func (r *reportRepository) GetReportDashboard(id uuid.UUID) (*models.ReportDashboard, error) {
	var dashboard models.ReportDashboard
	err := r.db.Preload("Company").
		Preload("Creator").
		Preload("UpdatedByUser").
		First(&dashboard, id).Error
	if err != nil {
		return nil, err
	}
	return &dashboard, nil
}

func (r *reportRepository) ListReportDashboards(companyID uuid.UUID, params map[string]interface{}) ([]models.ReportDashboard, int64, error) {
	var dashboards []models.ReportDashboard
	var total int64

	query := r.db.Model(&models.ReportDashboard{}).Where("company_id = ?", companyID)

	// Apply filters
	if isPublic, ok := params["is_public"].(bool); ok {
		query = query.Where("is_public = ?", isPublic)
	}

	if isDefault, ok := params["is_default"].(bool); ok {
		query = query.Where("is_default = ?", isDefault)
	}

	if theme, ok := params["theme"].(string); ok && theme != "" {
		query = query.Where("theme = ?", theme)
	}

	if search, ok := params["search"].(string); ok && search != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", 
			"%"+search+"%", "%"+search+"%")
	}

	if createdBy, ok := params["created_by"].(string); ok && createdBy != "" {
		query = query.Where("created_by = ?", createdBy)
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

	// Apply sorting
	sortBy := "created_at"
	if sb, ok := params["sort_by"].(string); ok && sb != "" {
		sortBy = sb
	}
	sortOrder := "DESC"
	if so, ok := params["sort_order"].(string); ok && so != "" {
		sortOrder = so
	}
	query = query.Order(sortBy + " " + sortOrder)

	// Load with relations
	if err := query.
		Preload("Creator").
		Find(&dashboards).Error; err != nil {
		return nil, 0, err
	}

	return dashboards, total, nil
}

func (r *reportRepository) DeleteReportDashboard(id uuid.UUID) error {
	return r.db.Delete(&models.ReportDashboard{}, id).Error
}

// Report Data Source operations
func (r *reportRepository) CreateReportDataSource(dataSource *models.ReportDataSource) error {
	return r.db.Create(dataSource).Error
}

func (r *reportRepository) UpdateReportDataSource(dataSource *models.ReportDataSource) error {
	return r.db.Save(dataSource).Error
}

func (r *reportRepository) GetReportDataSource(id uuid.UUID) (*models.ReportDataSource, error) {
	var dataSource models.ReportDataSource
	err := r.db.Preload("Company").
		Preload("Creator").
		First(&dataSource, id).Error
	if err != nil {
		return nil, err
	}
	return &dataSource, nil
}

func (r *reportRepository) ListReportDataSources(companyID uuid.UUID, params map[string]interface{}) ([]models.ReportDataSource, int64, error) {
	var dataSources []models.ReportDataSource
	var total int64

	query := r.db.Model(&models.ReportDataSource{}).Where("company_id = ?", companyID)

	// Apply filters
	if dataSourceType, ok := params["type"].(string); ok && dataSourceType != "" {
		query = query.Where("type = ?", dataSourceType)
	}

	if status, ok := params["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}

	if search, ok := params["search"].(string); ok && search != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", 
			"%"+search+"%", "%"+search+"%")
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

	// Apply sorting
	query = query.Order("created_at DESC")

	// Load with relations
	if err := query.
		Preload("Creator").
		Find(&dataSources).Error; err != nil {
		return nil, 0, err
	}

	return dataSources, total, nil
}

func (r *reportRepository) DeleteReportDataSource(id uuid.UUID) error {
	return r.db.Delete(&models.ReportDataSource{}, id).Error
}

// Report Schedule operations
func (r *reportRepository) CreateReportSchedule(schedule *models.ReportSchedule) error {
	return r.db.Create(schedule).Error
}

func (r *reportRepository) UpdateReportSchedule(schedule *models.ReportSchedule) error {
	return r.db.Save(schedule).Error
}

func (r *reportRepository) GetReportSchedule(id uuid.UUID) (*models.ReportSchedule, error) {
	var schedule models.ReportSchedule
	err := r.db.Preload("Report").
		Preload("Creator").
		First(&schedule, id).Error
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (r *reportRepository) ListReportSchedules(companyID uuid.UUID, params map[string]interface{}) ([]models.ReportSchedule, int64, error) {
	var schedules []models.ReportSchedule
	var total int64

	query := r.db.Model(&models.ReportSchedule{}).
		Joins("JOIN reports ON reports.id = report_schedules.report_id").
		Where("reports.company_id = ?", companyID)

	// Apply filters
	if isActive, ok := params["is_active"].(bool); ok {
		query = query.Where("report_schedules.is_active = ?", isActive)
	}

	if reportID, ok := params["report_id"].(string); ok && reportID != "" {
		query = query.Where("report_schedules.report_id = ?", reportID)
	}

	if lastStatus, ok := params["last_status"].(string); ok && lastStatus != "" {
		query = query.Where("report_schedules.last_status = ?", lastStatus)
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

	// Apply sorting
	query = query.Order("report_schedules.next_run ASC")

	// Load with relations
	if err := query.
		Preload("Report").
		Preload("Creator").
		Find(&schedules).Error; err != nil {
		return nil, 0, err
	}

	return schedules, total, nil
}

func (r *reportRepository) DeleteReportSchedule(id uuid.UUID) error {
	return r.db.Delete(&models.ReportSchedule{}, id).Error
}

// Business KPI operations
func (r *reportRepository) CreateBusinessKPI(kpi *models.BusinessKPI) error {
	return r.db.Create(kpi).Error
}

func (r *reportRepository) UpdateBusinessKPI(kpi *models.BusinessKPI) error {
	return r.db.Save(kpi).Error
}

func (r *reportRepository) GetBusinessKPI(id uuid.UUID) (*models.BusinessKPI, error) {
	var kpi models.BusinessKPI
	err := r.db.Preload("Company").
		Preload("Creator").
		First(&kpi, id).Error
	if err != nil {
		return nil, err
	}
	return &kpi, nil
}

func (r *reportRepository) ListBusinessKPIs(companyID uuid.UUID, params map[string]interface{}) ([]models.BusinessKPI, int64, error) {
	var kpis []models.BusinessKPI
	var total int64

	query := r.db.Model(&models.BusinessKPI{}).Where("company_id = ?", companyID)

	// Apply filters
	if category, ok := params["category"].(string); ok && category != "" {
		query = query.Where("category = ?", category)
	}

	if isActive, ok := params["is_active"].(bool); ok {
		query = query.Where("is_active = ?", isActive)
	}

	if frequency, ok := params["frequency"].(string); ok && frequency != "" {
		query = query.Where("frequency = ?", frequency)
	}

	if search, ok := params["search"].(string); ok && search != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", 
			"%"+search+"%", "%"+search+"%")
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

	// Apply sorting
	sortBy := "created_at"
	if sb, ok := params["sort_by"].(string); ok && sb != "" {
		sortBy = sb
	}
	sortOrder := "DESC"
	if so, ok := params["sort_order"].(string); ok && so != "" {
		sortOrder = so
	}
	query = query.Order(sortBy + " " + sortOrder)

	// Load with relations
	if err := query.
		Preload("Creator").
		Find(&kpis).Error; err != nil {
		return nil, 0, err
	}

	return kpis, total, nil
}

func (r *reportRepository) DeleteBusinessKPI(id uuid.UUID) error {
	return r.db.Delete(&models.BusinessKPI{}, id).Error
}

// Business operations
func (r *reportRepository) GetReportStatistics(companyID uuid.UUID) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Count reports by category
	var reportStats []struct {
		Category string `json:"category"`
		Count    int64  `json:"count"`
	}
	if err := r.db.Model(&models.Report{}).
		Where("company_id = ?", companyID).
		Select("category, COUNT(*) as count").
		Group("category").
		Find(&reportStats).Error; err != nil {
		return nil, err
	}
	stats["reports_by_category"] = reportStats

	// Count executions by status
	var executionStats []struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}
	if err := r.db.Model(&models.ReportExecution{}).
		Joins("JOIN reports ON reports.id = report_executions.report_id").
		Where("reports.company_id = ?", companyID).
		Select("report_executions.status, COUNT(*) as count").
		Group("report_executions.status").
		Find(&executionStats).Error; err != nil {
		return nil, err
	}
	stats["executions_by_status"] = executionStats

	// Total counts
	var totalReports, totalExecutions, totalTemplates, totalDashboards int64
	r.db.Model(&models.Report{}).Where("company_id = ?", companyID).Count(&totalReports)
	r.db.Model(&models.ReportExecution{}).
		Joins("JOIN reports ON reports.id = report_executions.report_id").
		Where("reports.company_id = ?", companyID).Count(&totalExecutions)
	r.db.Model(&models.ReportTemplate{}).
		Where("company_id = ? OR is_system_template = ?", companyID, true).Count(&totalTemplates)
	r.db.Model(&models.ReportDashboard{}).Where("company_id = ?", companyID).Count(&totalDashboards)

	stats["total_reports"] = totalReports
	stats["total_executions"] = totalExecutions
	stats["total_templates"] = totalTemplates
	stats["total_dashboards"] = totalDashboards

	return stats, nil
}

func (r *reportRepository) GetPopularReports(companyID uuid.UUID, limit int) ([]models.Report, error) {
	var reports []models.Report
	err := r.db.Where("company_id = ?", companyID).
		Order("view_count DESC").
		Limit(limit).
		Preload("Creator").
		Find(&reports).Error
	return reports, err
}

func (r *reportRepository) GetRecentExecutions(companyID uuid.UUID, limit int) ([]models.ReportExecution, error) {
	var executions []models.ReportExecution
	err := r.db.Joins("JOIN reports ON reports.id = report_executions.report_id").
		Where("reports.company_id = ?", companyID).
		Order("report_executions.started_at DESC").
		Limit(limit).
		Preload("Report").
		Preload("ExecutedByUser").
		Find(&executions).Error
	return executions, err
}

func (r *reportRepository) GetScheduledReports(companyID uuid.UUID) ([]models.ReportSchedule, error) {
	var schedules []models.ReportSchedule
	err := r.db.Joins("JOIN reports ON reports.id = report_schedules.report_id").
		Where("reports.company_id = ? AND report_schedules.is_active = ?", companyID, true).
		Order("report_schedules.next_run ASC").
		Preload("Report").
		Find(&schedules).Error
	return schedules, err
}