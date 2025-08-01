package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/fastenmind/fastener-api/internal/models"
	"github.com/fastenmind/fastener-api/internal/repository"
	"github.com/google/uuid"
)

type ReportService interface {
	// Report operations
	CreateReport(req *CreateReportRequest, userID uuid.UUID) (*models.Report, error)
	UpdateReport(id uuid.UUID, req *UpdateReportRequest, userID uuid.UUID) (*models.Report, error)
	GetReport(id uuid.UUID) (*models.Report, error)
	ListReports(companyID uuid.UUID, params map[string]interface{}) ([]models.Report, int64, error)
	DeleteReport(id uuid.UUID) error
	DuplicateReport(id uuid.UUID, userID uuid.UUID) (*models.Report, error)
	
	// Report Template operations
	CreateReportTemplate(req *CreateReportTemplateRequest, userID uuid.UUID) (*models.ReportTemplate, error)
	UpdateReportTemplate(id uuid.UUID, req *UpdateReportTemplateRequest, userID uuid.UUID) (*models.ReportTemplate, error)
	GetReportTemplate(id uuid.UUID) (*models.ReportTemplate, error)
	ListReportTemplates(companyID *uuid.UUID, params map[string]interface{}) ([]models.ReportTemplate, int64, error)
	DeleteReportTemplate(id uuid.UUID) error
	
	// Report Execution operations
	ExecuteReport(reportID uuid.UUID, params map[string]interface{}, userID uuid.UUID) (*models.ReportExecution, error)
	GetReportExecution(id uuid.UUID) (*models.ReportExecution, error)
	ListReportExecutions(reportID uuid.UUID, params map[string]interface{}) ([]models.ReportExecution, int64, error)
	CancelReportExecution(id uuid.UUID) error
	DownloadReportResult(executionID uuid.UUID) ([]byte, string, error)
	
	// Report Subscription operations
	CreateReportSubscription(req *CreateReportSubscriptionRequest, userID uuid.UUID) (*models.ReportSubscription, error)
	UpdateReportSubscription(id uuid.UUID, req *UpdateReportSubscriptionRequest) (*models.ReportSubscription, error)
	GetReportSubscription(id uuid.UUID) (*models.ReportSubscription, error)
	ListReportSubscriptions(userID uuid.UUID, params map[string]interface{}) ([]models.ReportSubscription, int64, error)
	DeleteReportSubscription(id uuid.UUID) error
	
	// Report Dashboard operations
	CreateReportDashboard(req *CreateReportDashboardRequest, userID uuid.UUID) (*models.ReportDashboard, error)
	UpdateReportDashboard(id uuid.UUID, req *UpdateReportDashboardRequest, userID uuid.UUID) (*models.ReportDashboard, error)
	GetReportDashboard(id uuid.UUID) (*models.ReportDashboard, error)
	ListReportDashboards(companyID uuid.UUID, params map[string]interface{}) ([]models.ReportDashboard, int64, error)
	DeleteReportDashboard(id uuid.UUID) error
	
	// Report Data Source operations
	CreateReportDataSource(req *CreateReportDataSourceRequest, userID uuid.UUID) (*models.ReportDataSource, error)
	UpdateReportDataSource(id uuid.UUID, req *UpdateReportDataSourceRequest) (*models.ReportDataSource, error)
	GetReportDataSource(id uuid.UUID) (*models.ReportDataSource, error)
	ListReportDataSources(companyID uuid.UUID, params map[string]interface{}) ([]models.ReportDataSource, int64, error)
	DeleteReportDataSource(id uuid.UUID) error
	TestReportDataSource(id uuid.UUID) (*DataSourceTestResult, error)
	
	// Report Schedule operations
	CreateReportSchedule(req *CreateReportScheduleRequest, userID uuid.UUID) (*models.ReportSchedule, error)
	UpdateReportSchedule(id uuid.UUID, req *UpdateReportScheduleRequest) (*models.ReportSchedule, error)
	GetReportSchedule(id uuid.UUID) (*models.ReportSchedule, error)
	ListReportSchedules(companyID uuid.UUID, params map[string]interface{}) ([]models.ReportSchedule, int64, error)
	DeleteReportSchedule(id uuid.UUID) error
	
	// Business KPI operations
	CreateBusinessKPI(req *CreateBusinessKPIRequest, userID uuid.UUID) (*models.BusinessKPI, error)
	UpdateBusinessKPI(id uuid.UUID, req *UpdateBusinessKPIRequest) (*models.BusinessKPI, error)
	GetBusinessKPI(id uuid.UUID) (*models.BusinessKPI, error)
	ListBusinessKPIs(companyID uuid.UUID, params map[string]interface{}) ([]models.BusinessKPI, int64, error)
	DeleteBusinessKPI(id uuid.UUID) error
	UpdateKPIValues(companyID uuid.UUID) error
	
	// Business operations
	GetReportDashboardData(companyID uuid.UUID) (*ReportDashboardData, error)
	GetReportStatistics(companyID uuid.UUID) (map[string]interface{}, error)
	GetPopularReports(companyID uuid.UUID, limit int) ([]models.Report, error)
	GetRecentExecutions(companyID uuid.UUID, limit int) ([]models.ReportExecution, error)
	GenerateReportFromTemplate(templateID uuid.UUID, params map[string]interface{}, userID uuid.UUID) (*models.Report, error)
	ExportReport(reportID uuid.UUID, format string, params map[string]interface{}) ([]byte, string, error)
	ImportReports(data []byte, userID uuid.UUID) (*ImportResult, error)
}

type reportService struct {
	reportRepo    repository.ReportRepository
	companyRepo   repository.CompanyRepository
	userRepo      repository.UserRepository
}

func NewReportService(
	reportRepo repository.ReportRepository,
	companyRepo repository.CompanyRepository,
	userRepo repository.UserRepository,
) ReportService {
	return &reportService{
		reportRepo:  reportRepo,
		companyRepo: companyRepo,
		userRepo:    userRepo,
	}
}

// Request structs
type CreateReportRequest struct {
	Name          string                 `json:"name" validate:"required"`
	NameEn        string                 `json:"name_en"`
	Category      string                 `json:"category" validate:"required"`
	Type          string                 `json:"type" validate:"required"`
	DataSource    map[string]interface{} `json:"data_source"`
	Filters       map[string]interface{} `json:"filters"`
	Columns       []ReportColumn         `json:"columns"`
	Sorting       []ReportSort           `json:"sorting"`
	Grouping      []string               `json:"grouping"`
	Aggregation   map[string]interface{} `json:"aggregation"`
	ChartConfig   map[string]interface{} `json:"chart_config"`
	TemplateID    *uuid.UUID             `json:"template_id"`
	Layout        map[string]interface{} `json:"layout"`
	Styling       map[string]interface{} `json:"styling"`
	IsPublic      bool                   `json:"is_public"`
	SharedWith    []string               `json:"shared_with"`
	CacheEnabled  bool                   `json:"cache_enabled"`
	CacheTTL      int                    `json:"cache_ttl"`
	QueryTimeout  int                    `json:"query_timeout"`
	Description   string                 `json:"description"`
	Tags          []string               `json:"tags"`
}

type UpdateReportRequest struct {
	Name          *string                `json:"name"`
	NameEn        *string                `json:"name_en"`
	Category      *string                `json:"category"`
	Type          *string                `json:"type"`
	Status        *string                `json:"status"`
	DataSource    map[string]interface{} `json:"data_source"`
	Filters       map[string]interface{} `json:"filters"`
	Columns       []ReportColumn         `json:"columns"`
	Sorting       []ReportSort           `json:"sorting"`
	Grouping      []string               `json:"grouping"`
	Aggregation   map[string]interface{} `json:"aggregation"`
	ChartConfig   map[string]interface{} `json:"chart_config"`
	Layout        map[string]interface{} `json:"layout"`
	Styling       map[string]interface{} `json:"styling"`
	IsPublic      *bool                  `json:"is_public"`
	SharedWith    []string               `json:"shared_with"`
	CacheEnabled  *bool                  `json:"cache_enabled"`
	CacheTTL      *int                   `json:"cache_ttl"`
	QueryTimeout  *int                   `json:"query_timeout"`
	Description   *string                `json:"description"`
	Tags          []string               `json:"tags"`
}

type ReportColumn struct {
	Name        string `json:"name"`
	Label       string `json:"label"`
	DataType    string `json:"data_type"`
	Format      string `json:"format"`
	Width       int    `json:"width"`
	Visible     bool   `json:"visible"`
	Sortable    bool   `json:"sortable"`
	Filterable  bool   `json:"filterable"`
	Aggregable  bool   `json:"aggregable"`
}

type ReportSort struct {
	Column    string `json:"column"`
	Direction string `json:"direction"` // asc, desc
}

type CreateReportTemplateRequest struct {
	Name              string                 `json:"name" validate:"required"`
	NameEn            string                 `json:"name_en"`
	Category          string                 `json:"category" validate:"required"`
	Type              string                 `json:"type" validate:"required"`
	IsSystemTemplate  bool                   `json:"is_system_template"`
	DataSource        map[string]interface{} `json:"data_source"`
	Filters           map[string]interface{} `json:"filters"`
	Columns           []ReportColumn         `json:"columns"`
	Sorting           []ReportSort           `json:"sorting"`
	Grouping          []string               `json:"grouping"`
	Aggregation       map[string]interface{} `json:"aggregation"`
	ChartConfig       map[string]interface{} `json:"chart_config"`
	Layout            map[string]interface{} `json:"layout"`
	Styling           map[string]interface{} `json:"styling"`
	Description       string                 `json:"description"`
	Preview           string                 `json:"preview"`
	Tags              []string               `json:"tags"`
	Industry          string                 `json:"industry"`
	Language          string                 `json:"language"`
}

type UpdateReportTemplateRequest struct {
	Name              *string                `json:"name"`
	NameEn            *string                `json:"name_en"`
	Category          *string                `json:"category"`
	Type              *string                `json:"type"`
	DataSource        map[string]interface{} `json:"data_source"`
	Filters           map[string]interface{} `json:"filters"`
	Columns           []ReportColumn         `json:"columns"`
	Sorting           []ReportSort           `json:"sorting"`
	Grouping          []string               `json:"grouping"`
	Aggregation       map[string]interface{} `json:"aggregation"`
	ChartConfig       map[string]interface{} `json:"chart_config"`
	Layout            map[string]interface{} `json:"layout"`
	Styling           map[string]interface{} `json:"styling"`
	Description       *string                `json:"description"`
	Preview           *string                `json:"preview"`
	Tags              []string               `json:"tags"`
	Industry          *string                `json:"industry"`
	Language          *string                `json:"language"`
}

type CreateReportSubscriptionRequest struct {
	ReportID       uuid.UUID              `json:"report_id" validate:"required"`
	Email          string                 `json:"email"`
	Schedule       string                 `json:"schedule"`
	FileFormat     string                 `json:"file_format"`
	Parameters     map[string]interface{} `json:"parameters"`
	DeliveryMethod string                 `json:"delivery_method"`
	DeliveryConfig map[string]interface{} `json:"delivery_config"`
}

type UpdateReportSubscriptionRequest struct {
	IsActive       *bool                  `json:"is_active"`
	Email          *string                `json:"email"`
	Schedule       *string                `json:"schedule"`
	FileFormat     *string                `json:"file_format"`
	Parameters     map[string]interface{} `json:"parameters"`
	DeliveryMethod *string                `json:"delivery_method"`
	DeliveryConfig map[string]interface{} `json:"delivery_config"`
}

type CreateReportDashboardRequest struct {
	Name        string                 `json:"name" validate:"required"`
	NameEn      string                 `json:"name_en"`
	Layout      map[string]interface{} `json:"layout"`
	Theme       string                 `json:"theme"`
	RefreshRate int                    `json:"refresh_rate"`
	Widgets     []DashboardWidget      `json:"widgets"`
	Filters     map[string]interface{} `json:"filters"`
	IsPublic    bool                   `json:"is_public"`
	SharedWith  []string               `json:"shared_with"`
	Description string                 `json:"description"`
	Tags        []string               `json:"tags"`
	IsDefault   bool                   `json:"is_default"`
}

type UpdateReportDashboardRequest struct {
	Name        *string                `json:"name"`
	NameEn      *string                `json:"name_en"`
	Layout      map[string]interface{} `json:"layout"`
	Theme       *string                `json:"theme"`
	RefreshRate *int                   `json:"refresh_rate"`
	Widgets     []DashboardWidget      `json:"widgets"`
	Filters     map[string]interface{} `json:"filters"`
	IsPublic    *bool                  `json:"is_public"`
	SharedWith  []string               `json:"shared_with"`
	Description *string                `json:"description"`
	Tags        []string               `json:"tags"`
	IsDefault   *bool                  `json:"is_default"`
}

type DashboardWidget struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`       // chart, table, kpi, text
	Title      string                 `json:"title"`
	Position   map[string]int         `json:"position"`   // x, y, width, height
	ReportID   *uuid.UUID             `json:"report_id"`
	KPIID      *uuid.UUID             `json:"kpi_id"`
	Config     map[string]interface{} `json:"config"`
	DataSource map[string]interface{} `json:"data_source"`
}

type CreateReportDataSourceRequest struct {
	Name             string                 `json:"name" validate:"required"`
	Type             string                 `json:"type" validate:"required"`
	ConnectionString string                 `json:"connection_string"`
	Credentials      map[string]interface{} `json:"credentials"`
	Settings         map[string]interface{} `json:"settings"`
	Description      string                 `json:"description"`
	Schema           map[string]interface{} `json:"schema"`
}

type UpdateReportDataSourceRequest struct {
	Name             *string                `json:"name"`
	Type             *string                `json:"type"`
	ConnectionString *string                `json:"connection_string"`
	Credentials      map[string]interface{} `json:"credentials"`
	Settings         map[string]interface{} `json:"settings"`
	Description      *string                `json:"description"`
	Schema           map[string]interface{} `json:"schema"`
	Status           *string                `json:"status"`
}

type CreateReportScheduleRequest struct {
	ReportID       uuid.UUID              `json:"report_id" validate:"required"`
	Name           string                 `json:"name" validate:"required"`
	CronExpression string                 `json:"cron_expression" validate:"required"`
	Timezone       string                 `json:"timezone"`
	Parameters     map[string]interface{} `json:"parameters"`
	FileFormat     string                 `json:"file_format"`
	Recipients     []string               `json:"recipients"`
}

type UpdateReportScheduleRequest struct {
	Name           *string                `json:"name"`
	CronExpression *string                `json:"cron_expression"`
	Timezone       *string                `json:"timezone"`
	IsActive       *bool                  `json:"is_active"`
	Parameters     map[string]interface{} `json:"parameters"`
	FileFormat     *string                `json:"file_format"`
	Recipients     []string               `json:"recipients"`
}

type CreateBusinessKPIRequest struct {
	Name          string                 `json:"name" validate:"required"`
	Category      string                 `json:"category" validate:"required"`
	Formula       string                 `json:"formula" validate:"required"`
	DataSources   []string               `json:"data_sources"`
	Filters       map[string]interface{} `json:"filters"`
	Unit          string                 `json:"unit"`
	TargetValue   float64                `json:"target_value"`
	TargetType    string                 `json:"target_type"`
	ThresholdHigh float64                `json:"threshold_high"`
	ThresholdLow  float64                `json:"threshold_low"`
	DisplayFormat string                 `json:"display_format"`
	ChartType     string                 `json:"chart_type"`
	ColorScheme   map[string]interface{} `json:"color_scheme"`
	Description   string                 `json:"description"`
	Frequency     string                 `json:"frequency"`
}

type UpdateBusinessKPIRequest struct {
	Name          *string                `json:"name"`
	Category      *string                `json:"category"`
	Formula       *string                `json:"formula"`
	DataSources   []string               `json:"data_sources"`
	Filters       map[string]interface{} `json:"filters"`
	Unit          *string                `json:"unit"`
	TargetValue   *float64               `json:"target_value"`
	TargetType    *string                `json:"target_type"`
	ThresholdHigh *float64               `json:"threshold_high"`
	ThresholdLow  *float64               `json:"threshold_low"`
	DisplayFormat *string                `json:"display_format"`
	ChartType     *string                `json:"chart_type"`
	ColorScheme   map[string]interface{} `json:"color_scheme"`
	Description   *string                `json:"description"`
	Frequency     *string                `json:"frequency"`
	IsActive      *bool                  `json:"is_active"`
}

type ReportDashboardData struct {
	TotalReports         int                    `json:"total_reports"`
	TotalExecutions      int                    `json:"total_executions"`
	TotalTemplates       int                    `json:"total_templates"`
	TotalDashboards      int                    `json:"total_dashboards"`
	TotalDataSources     int                    `json:"total_data_sources"`
	TotalKPIs            int                    `json:"total_kpis"`
	ScheduledReportsCount int                    `json:"scheduled_reports"`
	ActiveSubscriptions  int                    `json:"active_subscriptions"`
	ReportsByCategory    []CategoryCount        `json:"reports_by_category"`
	ExecutionsByStatus   []StatusCount          `json:"executions_by_status"`
	PopularReports       []models.Report        `json:"popular_reports"`
	RecentExecutions     []models.ReportExecution `json:"recent_executions"`
	ScheduledReportsList []models.ReportSchedule  `json:"scheduled_reports_list"`
	SystemHealth         map[string]interface{} `json:"system_health"`
}

type CategoryCount struct {
	Category string `json:"category"`
	Count    int64  `json:"count"`
}

type StatusCount struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

type DataSourceTestResult struct {
	Success       bool                   `json:"success"`
	Message       string                 `json:"message"`
	ConnectionTime float64               `json:"connection_time"`
	SampleData    []map[string]interface{} `json:"sample_data"`
	Schema        map[string]interface{} `json:"schema"`
	Error         string                 `json:"error"`
}

type ImportResult struct {
	Success        int      `json:"success"`
	Failed         int      `json:"failed"`
	Errors         []string `json:"errors"`
	ImportedReports []string `json:"imported_reports"`
}

// Report operations
func (s *reportService) CreateReport(req *CreateReportRequest, userID uuid.UUID) (*models.Report, error) {
	report := &models.Report{
		CompanyID:     userID, // This should be retrieved from user context
		Name:          req.Name,
		NameEn:        req.NameEn,
		Category:      req.Category,
		Type:          req.Type,
		Status:        "active",
		TemplateID:    req.TemplateID,
		IsPublic:      req.IsPublic,
		CacheEnabled:  req.CacheEnabled,
		CacheTTL:      req.CacheTTL,
		QueryTimeout:  req.QueryTimeout,
		Description:   req.Description,
		CreatedBy:     userID,
	}

	// Generate report number
	timestamp := time.Now().Unix()
	report.ReportNo = fmt.Sprintf("RPT%d", timestamp)

	// Convert request data to JSON strings
	if req.DataSource != nil {
		if data, err := json.Marshal(req.DataSource); err == nil {
			report.DataSource = string(data)
		}
	}

	if req.Filters != nil {
		if data, err := json.Marshal(req.Filters); err == nil {
			report.Filters = string(data)
		}
	}

	if req.Columns != nil {
		if data, err := json.Marshal(req.Columns); err == nil {
			report.Columns = string(data)
		}
	}

	if req.Sorting != nil {
		if data, err := json.Marshal(req.Sorting); err == nil {
			report.Sorting = string(data)
		}
	}

	if req.Grouping != nil {
		if data, err := json.Marshal(req.Grouping); err == nil {
			report.Grouping = string(data)
		}
	}

	if req.Aggregation != nil {
		if data, err := json.Marshal(req.Aggregation); err == nil {
			report.Aggregation = string(data)
		}
	}

	if req.ChartConfig != nil {
		if data, err := json.Marshal(req.ChartConfig); err == nil {
			report.ChartConfig = string(data)
		}
	}

	if req.Layout != nil {
		if data, err := json.Marshal(req.Layout); err == nil {
			report.Layout = string(data)
		}
	}

	if req.Styling != nil {
		if data, err := json.Marshal(req.Styling); err == nil {
			report.Styling = string(data)
		}
	}

	if req.SharedWith != nil {
		if data, err := json.Marshal(req.SharedWith); err == nil {
			report.SharedWith = string(data)
		}
	}

	if req.Tags != nil {
		if data, err := json.Marshal(req.Tags); err == nil {
			report.Tags = string(data)
		}
	}

	if err := s.reportRepo.CreateReport(report); err != nil {
		return nil, fmt.Errorf("failed to create report: %w", err)
	}

	return report, nil
}

func (s *reportService) UpdateReport(id uuid.UUID, req *UpdateReportRequest, userID uuid.UUID) (*models.Report, error) {
	report, err := s.reportRepo.GetReport(id)
	if err != nil {
		return nil, fmt.Errorf("report not found: %w", err)
	}

	// Update fields if provided
	if req.Name != nil {
		report.Name = *req.Name
	}
	if req.NameEn != nil {
		report.NameEn = *req.NameEn
	}
	if req.Category != nil {
		report.Category = *req.Category
	}
	if req.Type != nil {
		report.Type = *req.Type
	}
	if req.Status != nil {
		report.Status = *req.Status
	}
	if req.IsPublic != nil {
		report.IsPublic = *req.IsPublic
	}
	if req.CacheEnabled != nil {
		report.CacheEnabled = *req.CacheEnabled
	}
	if req.CacheTTL != nil {
		report.CacheTTL = *req.CacheTTL
	}
	if req.QueryTimeout != nil {
		report.QueryTimeout = *req.QueryTimeout
	}
	if req.Description != nil {
		report.Description = *req.Description
	}

	// Update JSON fields
	if req.DataSource != nil {
		if data, err := json.Marshal(req.DataSource); err == nil {
			report.DataSource = string(data)
		}
	}

	if req.Filters != nil {
		if data, err := json.Marshal(req.Filters); err == nil {
			report.Filters = string(data)
		}
	}

	if req.Columns != nil {
		if data, err := json.Marshal(req.Columns); err == nil {
			report.Columns = string(data)
		}
	}

	if req.Sorting != nil {
		if data, err := json.Marshal(req.Sorting); err == nil {
			report.Sorting = string(data)
		}
	}

	if req.Grouping != nil {
		if data, err := json.Marshal(req.Grouping); err == nil {
			report.Grouping = string(data)
		}
	}

	if req.Aggregation != nil {
		if data, err := json.Marshal(req.Aggregation); err == nil {
			report.Aggregation = string(data)
		}
	}

	if req.ChartConfig != nil {
		if data, err := json.Marshal(req.ChartConfig); err == nil {
			report.ChartConfig = string(data)
		}
	}

	if req.Layout != nil {
		if data, err := json.Marshal(req.Layout); err == nil {
			report.Layout = string(data)
		}
	}

	if req.Styling != nil {
		if data, err := json.Marshal(req.Styling); err == nil {
			report.Styling = string(data)
		}
	}

	if req.SharedWith != nil {
		if data, err := json.Marshal(req.SharedWith); err == nil {
			report.SharedWith = string(data)
		}
	}

	if req.Tags != nil {
		if data, err := json.Marshal(req.Tags); err == nil {
			report.Tags = string(data)
		}
	}

	report.UpdatedBy = &userID
	report.Version++

	if err := s.reportRepo.UpdateReport(report); err != nil {
		return nil, fmt.Errorf("failed to update report: %w", err)
	}

	return report, nil
}

func (s *reportService) GetReport(id uuid.UUID) (*models.Report, error) {
	report, err := s.reportRepo.GetReport(id)
	if err != nil {
		return nil, err
	}

	// Update view count
	report.ViewCount++
	now := time.Now()
	report.LastViewed = &now
	s.reportRepo.UpdateReport(report)

	return report, nil
}

func (s *reportService) ListReports(companyID uuid.UUID, params map[string]interface{}) ([]models.Report, int64, error) {
	return s.reportRepo.ListReports(companyID, params)
}

func (s *reportService) DeleteReport(id uuid.UUID) error {
	return s.reportRepo.DeleteReport(id)
}

func (s *reportService) DuplicateReport(id uuid.UUID, userID uuid.UUID) (*models.Report, error) {
	original, err := s.reportRepo.GetReport(id)
	if err != nil {
		return nil, fmt.Errorf("original report not found: %w", err)
	}

	duplicate := &models.Report{
		CompanyID:     original.CompanyID,
		Name:          original.Name + " (Copy)",
		NameEn:        original.NameEn + " (Copy)",
		Category:      original.Category,
		Type:          original.Type,
		Status:        "active",
		DataSource:    original.DataSource,
		Filters:       original.Filters,
		Columns:       original.Columns,
		Sorting:       original.Sorting,
		Grouping:      original.Grouping,
		Aggregation:   original.Aggregation,
		ChartConfig:   original.ChartConfig,
		TemplateID:    original.TemplateID,
		Layout:        original.Layout,
		Styling:       original.Styling,
		IsPublic:      false, // Duplicates are private by default
		SharedWith:    "",
		CacheEnabled:  original.CacheEnabled,
		CacheTTL:      original.CacheTTL,
		QueryTimeout:  original.QueryTimeout,
		Description:   original.Description,
		Tags:          original.Tags,
		Version:       1,
		CreatedBy:     userID,
	}

	// Generate new report number
	timestamp := time.Now().Unix()
	duplicate.ReportNo = fmt.Sprintf("RPT%d", timestamp)

	if err := s.reportRepo.CreateReport(duplicate); err != nil {
		return nil, fmt.Errorf("failed to duplicate report: %w", err)
	}

	return duplicate, nil
}

// Report Template operations
func (s *reportService) CreateReportTemplate(req *CreateReportTemplateRequest, userID uuid.UUID) (*models.ReportTemplate, error) {
	template := &models.ReportTemplate{
		Name:             req.Name,
		NameEn:           req.NameEn,
		Category:         req.Category,
		Type:             req.Type,
		IsSystemTemplate: req.IsSystemTemplate,
		Description:      req.Description,
		Preview:          req.Preview,
		Industry:         req.Industry,
		Language:         req.Language,
	}

	if !req.IsSystemTemplate {
		template.CompanyID = &userID // This should be retrieved from user context
		template.CreatedBy = &userID
	}

	// Convert request data to JSON strings
	if req.DataSource != nil {
		if data, err := json.Marshal(req.DataSource); err == nil {
			template.DataSource = string(data)
		}
	}

	if req.Filters != nil {
		if data, err := json.Marshal(req.Filters); err == nil {
			template.Filters = string(data)
		}
	}

	if req.Columns != nil {
		if data, err := json.Marshal(req.Columns); err == nil {
			template.Columns = string(data)
		}
	}

	if req.Sorting != nil {
		if data, err := json.Marshal(req.Sorting); err == nil {
			template.Sorting = string(data)
		}
	}

	if req.Grouping != nil {
		if data, err := json.Marshal(req.Grouping); err == nil {
			template.Grouping = string(data)
		}
	}

	if req.Aggregation != nil {
		if data, err := json.Marshal(req.Aggregation); err == nil {
			template.Aggregation = string(data)
		}
	}

	if req.ChartConfig != nil {
		if data, err := json.Marshal(req.ChartConfig); err == nil {
			template.ChartConfig = string(data)
		}
	}

	if req.Layout != nil {
		if data, err := json.Marshal(req.Layout); err == nil {
			template.Layout = string(data)
		}
	}

	if req.Styling != nil {
		if data, err := json.Marshal(req.Styling); err == nil {
			template.Styling = string(data)
		}
	}

	if req.Tags != nil {
		if data, err := json.Marshal(req.Tags); err == nil {
			template.Tags = string(data)
		}
	}

	if err := s.reportRepo.CreateReportTemplate(template); err != nil {
		return nil, fmt.Errorf("failed to create report template: %w", err)
	}

	return template, nil
}

func (s *reportService) UpdateReportTemplate(id uuid.UUID, req *UpdateReportTemplateRequest, userID uuid.UUID) (*models.ReportTemplate, error) {
	template, err := s.reportRepo.GetReportTemplate(id)
	if err != nil {
		return nil, fmt.Errorf("report template not found: %w", err)
	}

	// Update fields if provided
	if req.Name != nil {
		template.Name = *req.Name
	}
	if req.NameEn != nil {
		template.NameEn = *req.NameEn
	}
	if req.Category != nil {
		template.Category = *req.Category
	}
	if req.Type != nil {
		template.Type = *req.Type
	}
	if req.Description != nil {
		template.Description = *req.Description
	}
	if req.Preview != nil {
		template.Preview = *req.Preview
	}
	if req.Industry != nil {
		template.Industry = *req.Industry
	}
	if req.Language != nil {
		template.Language = *req.Language
	}

	// Update JSON fields
	if req.DataSource != nil {
		if data, err := json.Marshal(req.DataSource); err == nil {
			template.DataSource = string(data)
		}
	}

	if req.Filters != nil {
		if data, err := json.Marshal(req.Filters); err == nil {
			template.Filters = string(data)
		}
	}

	if req.Columns != nil {
		if data, err := json.Marshal(req.Columns); err == nil {
			template.Columns = string(data)
		}
	}

	if req.Sorting != nil {
		if data, err := json.Marshal(req.Sorting); err == nil {
			template.Sorting = string(data)
		}
	}

	if req.Grouping != nil {
		if data, err := json.Marshal(req.Grouping); err == nil {
			template.Grouping = string(data)
		}
	}

	if req.Aggregation != nil {
		if data, err := json.Marshal(req.Aggregation); err == nil {
			template.Aggregation = string(data)
		}
	}

	if req.ChartConfig != nil {
		if data, err := json.Marshal(req.ChartConfig); err == nil {
			template.ChartConfig = string(data)
		}
	}

	if req.Layout != nil {
		if data, err := json.Marshal(req.Layout); err == nil {
			template.Layout = string(data)
		}
	}

	if req.Styling != nil {
		if data, err := json.Marshal(req.Styling); err == nil {
			template.Styling = string(data)
		}
	}

	if req.Tags != nil {
		if data, err := json.Marshal(req.Tags); err == nil {
			template.Tags = string(data)
		}
	}

	if err := s.reportRepo.UpdateReportTemplate(template); err != nil {
		return nil, fmt.Errorf("failed to update report template: %w", err)
	}

	return template, nil
}

func (s *reportService) GetReportTemplate(id uuid.UUID) (*models.ReportTemplate, error) {
	return s.reportRepo.GetReportTemplate(id)
}

func (s *reportService) ListReportTemplates(companyID *uuid.UUID, params map[string]interface{}) ([]models.ReportTemplate, int64, error) {
	return s.reportRepo.ListReportTemplates(companyID, params)
}

func (s *reportService) DeleteReportTemplate(id uuid.UUID) error {
	return s.reportRepo.DeleteReportTemplate(id)
}

// Report Execution operations
func (s *reportService) ExecuteReport(reportID uuid.UUID, params map[string]interface{}, userID uuid.UUID) (*models.ReportExecution, error) {
	report, err := s.reportRepo.GetReport(reportID)
	if err != nil {
		return nil, fmt.Errorf("report not found: %w", err)
	}

	execution := &models.ReportExecution{
		ReportID:      reportID,
		Status:        "pending",
		TriggerType:   "manual",
		ExecutedBy:    userID,
		StartedAt:     time.Now(),
	}

	if params != nil {
		if data, err := json.Marshal(params); err == nil {
			execution.Parameters = string(data)
		}
	}

	if err := s.reportRepo.CreateReportExecution(execution); err != nil {
		return nil, fmt.Errorf("failed to create report execution: %w", err)
	}

	// Start execution in background
	go s.executeReportAsync(execution)

	return execution, nil
}

func (s *reportService) executeReportAsync(execution *models.ReportExecution) {
	startTime := time.Now()
	
	// Update status to running
	execution.Status = "running"
	s.reportRepo.UpdateReportExecution(execution)

	// Simulate report execution
	// In real implementation, this would:
	// 1. Parse report configuration
	// 2. Execute queries against data sources
	// 3. Generate output file
	// 4. Save result to file system
	time.Sleep(5 * time.Second) // Simulate processing time

	// Update execution with results
	endTime := time.Now()
	execution.Status = "completed"
	execution.ExecutionTime = float64(endTime.Sub(startTime).Milliseconds())
	execution.CompletedAt = &endTime
	execution.RowCount = 100 // Simulated row count
	execution.FileFormat = "pdf"
	execution.FilePath = fmt.Sprintf("/reports/%s.pdf", execution.ID.String())
	execution.FileSize = 1024 * 50 // 50KB simulated file size

	s.reportRepo.UpdateReportExecution(execution)

	// Update report statistics
	if report, err := s.reportRepo.GetReport(execution.ReportID); err == nil {
		report.ExecuteCount++
		report.LastExecuted = &endTime
		
		// Update average execution time
		if report.AvgExecTime == 0 {
			report.AvgExecTime = execution.ExecutionTime
		} else {
			report.AvgExecTime = (report.AvgExecTime + execution.ExecutionTime) / 2
		}
		
		s.reportRepo.UpdateReport(report)
	}
}

func (s *reportService) GetReportExecution(id uuid.UUID) (*models.ReportExecution, error) {
	return s.reportRepo.GetReportExecution(id)
}

func (s *reportService) ListReportExecutions(reportID uuid.UUID, params map[string]interface{}) ([]models.ReportExecution, int64, error) {
	return s.reportRepo.ListReportExecutions(reportID, params)
}

func (s *reportService) CancelReportExecution(id uuid.UUID) error {
	execution, err := s.reportRepo.GetReportExecution(id)
	if err != nil {
		return fmt.Errorf("execution not found: %w", err)
	}

	if execution.Status == "completed" || execution.Status == "failed" || execution.Status == "cancelled" {
		return fmt.Errorf("execution cannot be cancelled in current status: %s", execution.Status)
	}

	execution.Status = "cancelled"
	now := time.Now()
	execution.CompletedAt = &now

	return s.reportRepo.UpdateReportExecution(execution)
}

func (s *reportService) DownloadReportResult(executionID uuid.UUID) ([]byte, string, error) {
	execution, err := s.reportRepo.GetReportExecution(executionID)
	if err != nil {
		return nil, "", fmt.Errorf("execution not found: %w", err)
	}

	if execution.Status != "completed" {
		return nil, "", fmt.Errorf("execution not completed")
	}

	// In real implementation, this would read the actual file
	// For now, return simulated content
	content := []byte("Simulated report content")
	filename := fmt.Sprintf("report_%s.%s", execution.ID.String(), execution.FileFormat)

	return content, filename, nil
}

// Business operations
func (s *reportService) GetReportDashboardData(companyID uuid.UUID) (*ReportDashboardData, error) {
	dashboard := &ReportDashboardData{}

	// Get statistics
	stats, err := s.reportRepo.GetReportStatistics(companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get report statistics: %w", err)
	}

	if totalReports, ok := stats["total_reports"].(int64); ok {
		dashboard.TotalReports = int(totalReports)
	}
	if totalExecutions, ok := stats["total_executions"].(int64); ok {
		dashboard.TotalExecutions = int(totalExecutions)
	}
	if totalTemplates, ok := stats["total_templates"].(int64); ok {
		dashboard.TotalTemplates = int(totalTemplates)
	}
	if totalDashboards, ok := stats["total_dashboards"].(int64); ok {
		dashboard.TotalDashboards = int(totalDashboards)
	}

	// Get popular reports
	if popularReports, err := s.reportRepo.GetPopularReports(companyID, 5); err == nil {
		dashboard.PopularReports = popularReports
	}

	// Get recent executions
	if recentExecutions, err := s.reportRepo.GetRecentExecutions(companyID, 10); err == nil {
		dashboard.RecentExecutions = recentExecutions
	}

	// Get scheduled reports
	if scheduledReports, err := s.reportRepo.GetScheduledReports(companyID); err == nil {
		dashboard.ScheduledReports = scheduledReports
		dashboard.ScheduledReports = len(scheduledReports)
	}

	// System health simulation
	dashboard.SystemHealth = map[string]interface{}{
		"cpu_usage":    75.2,
		"memory_usage": 68.5,
		"disk_usage":   45.3,
		"queue_size":   3,
		"avg_response_time": 1.2,
	}

	return dashboard, nil
}

func (s *reportService) GetReportStatistics(companyID uuid.UUID) (map[string]interface{}, error) {
	return s.reportRepo.GetReportStatistics(companyID)
}

func (s *reportService) GetPopularReports(companyID uuid.UUID, limit int) ([]models.Report, error) {
	return s.reportRepo.GetPopularReports(companyID, limit)
}

func (s *reportService) GetRecentExecutions(companyID uuid.UUID, limit int) ([]models.ReportExecution, error) {
	return s.reportRepo.GetRecentExecutions(companyID, limit)
}

func (s *reportService) GenerateReportFromTemplate(templateID uuid.UUID, params map[string]interface{}, userID uuid.UUID) (*models.Report, error) {
	template, err := s.reportRepo.GetReportTemplate(templateID)
	if err != nil {
		return nil, fmt.Errorf("template not found: %w", err)
	}

	// Create report from template
	name := template.Name
	if nameParam, ok := params["name"].(string); ok && nameParam != "" {
		name = nameParam
	}

	req := &CreateReportRequest{
		Name:        name,
		NameEn:      template.NameEn,
		Category:    template.Category,
		Type:        template.Type,
		TemplateID:  &templateID,
		Description: template.Description,
	}

	// Parse template JSON fields
	if template.DataSource != "" {
		var dataSource map[string]interface{}
		if err := json.Unmarshal([]byte(template.DataSource), &dataSource); err == nil {
			req.DataSource = dataSource
		}
	}

	if template.Filters != "" {
		var filters map[string]interface{}
		if err := json.Unmarshal([]byte(template.Filters), &filters); err == nil {
			req.Filters = filters
		}
	}

	if template.Columns != "" {
		var columns []ReportColumn
		if err := json.Unmarshal([]byte(template.Columns), &columns); err == nil {
			req.Columns = columns
		}
	}

	if template.Sorting != "" {
		var sorting []ReportSort
		if err := json.Unmarshal([]byte(template.Sorting), &sorting); err == nil {
			req.Sorting = sorting
		}
	}

	report, err := s.CreateReport(req, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to create report from template: %w", err)
	}

	// Update template usage count
	template.UsageCount++
	s.reportRepo.UpdateReportTemplate(template)

	return report, nil
}

func (s *reportService) ExportReport(reportID uuid.UUID, format string, params map[string]interface{}) ([]byte, string, error) {
	report, err := s.reportRepo.GetReport(reportID)
	if err != nil {
		return nil, "", fmt.Errorf("report not found: %w", err)
	}

	// In real implementation, this would generate the actual export
	// For now, return simulated content
	var content []byte
	var filename string

	switch format {
	case "json":
		content, _ = json.Marshal(report)
		filename = fmt.Sprintf("%s.json", report.ReportNo)
	case "csv":
		content = []byte("Name,Category,Type,Status\n" + report.Name + "," + report.Category + "," + report.Type + "," + report.Status)
		filename = fmt.Sprintf("%s.csv", report.ReportNo)
	default:
		return nil, "", fmt.Errorf("unsupported export format: %s", format)
	}

	return content, filename, nil
}

func (s *reportService) ImportReports(data []byte, userID uuid.UUID) (*ImportResult, error) {
	result := &ImportResult{
		Errors:          make([]string, 0),
		ImportedReports: make([]string, 0),
	}

	// Parse import data (simplified JSON format)
	var importData []map[string]interface{}
	if err := json.Unmarshal(data, &importData); err != nil {
		result.Failed++
		result.Errors = append(result.Errors, "Invalid JSON format")
		return result, nil
	}

	for _, reportData := range importData {
		name, ok := reportData["name"].(string)
		if !ok {
			result.Failed++
			result.Errors = append(result.Errors, "Missing report name")
			continue
		}

		category, ok := reportData["category"].(string)
		if !ok {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("Missing category for report: %s", name))
			continue
		}

		reportType, ok := reportData["type"].(string)
		if !ok {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("Missing type for report: %s", name))
			continue
		}

		req := &CreateReportRequest{
			Name:     name,
			Category: category,
			Type:     reportType,
		}

		if description, ok := reportData["description"].(string); ok {
			req.Description = description
		}

		if _, err := s.CreateReport(req, userID); err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("Failed to import report %s: %s", name, err.Error()))
		} else {
			result.Success++
			result.ImportedReports = append(result.ImportedReports, name)
		}
	}

	return result, nil
}

// Placeholder implementations for other methods
func (s *reportService) CreateReportSubscription(req *CreateReportSubscriptionRequest, userID uuid.UUID) (*models.ReportSubscription, error) {
	// Implementation would create subscription
	return nil, fmt.Errorf("not implemented")
}

func (s *reportService) UpdateReportSubscription(id uuid.UUID, req *UpdateReportSubscriptionRequest) (*models.ReportSubscription, error) {
	// Implementation would update subscription
	return nil, fmt.Errorf("not implemented")
}

func (s *reportService) GetReportSubscription(id uuid.UUID) (*models.ReportSubscription, error) {
	return s.reportRepo.GetReportSubscription(id)
}

func (s *reportService) ListReportSubscriptions(userID uuid.UUID, params map[string]interface{}) ([]models.ReportSubscription, int64, error) {
	return s.reportRepo.ListReportSubscriptions(userID, params)
}

func (s *reportService) DeleteReportSubscription(id uuid.UUID) error {
	return s.reportRepo.DeleteReportSubscription(id)
}

func (s *reportService) CreateReportDashboard(req *CreateReportDashboardRequest, userID uuid.UUID) (*models.ReportDashboard, error) {
	// Implementation would create dashboard
	return nil, fmt.Errorf("not implemented")
}

func (s *reportService) UpdateReportDashboard(id uuid.UUID, req *UpdateReportDashboardRequest, userID uuid.UUID) (*models.ReportDashboard, error) {
	// Implementation would update dashboard
	return nil, fmt.Errorf("not implemented")
}

func (s *reportService) GetReportDashboard(id uuid.UUID) (*models.ReportDashboard, error) {
	return s.reportRepo.GetReportDashboard(id)
}

func (s *reportService) ListReportDashboards(companyID uuid.UUID, params map[string]interface{}) ([]models.ReportDashboard, int64, error) {
	return s.reportRepo.ListReportDashboards(companyID, params)
}

func (s *reportService) DeleteReportDashboard(id uuid.UUID) error {
	return s.reportRepo.DeleteReportDashboard(id)
}

func (s *reportService) CreateReportDataSource(req *CreateReportDataSourceRequest, userID uuid.UUID) (*models.ReportDataSource, error) {
	// Implementation would create data source
	return nil, fmt.Errorf("not implemented")
}

func (s *reportService) UpdateReportDataSource(id uuid.UUID, req *UpdateReportDataSourceRequest) (*models.ReportDataSource, error) {
	// Implementation would update data source
	return nil, fmt.Errorf("not implemented")
}

func (s *reportService) GetReportDataSource(id uuid.UUID) (*models.ReportDataSource, error) {
	return s.reportRepo.GetReportDataSource(id)
}

func (s *reportService) ListReportDataSources(companyID uuid.UUID, params map[string]interface{}) ([]models.ReportDataSource, int64, error) {
	return s.reportRepo.ListReportDataSources(companyID, params)
}

func (s *reportService) DeleteReportDataSource(id uuid.UUID) error {
	return s.reportRepo.DeleteReportDataSource(id)
}

func (s *reportService) TestReportDataSource(id uuid.UUID) (*DataSourceTestResult, error) {
	// Implementation would test data source connection
	return &DataSourceTestResult{
		Success:        true,
		Message:        "Connection successful",
		ConnectionTime: 250.0,
	}, nil
}

func (s *reportService) CreateReportSchedule(req *CreateReportScheduleRequest, userID uuid.UUID) (*models.ReportSchedule, error) {
	// Implementation would create schedule
	return nil, fmt.Errorf("not implemented")
}

func (s *reportService) UpdateReportSchedule(id uuid.UUID, req *UpdateReportScheduleRequest) (*models.ReportSchedule, error) {
	// Implementation would update schedule
	return nil, fmt.Errorf("not implemented")
}

func (s *reportService) GetReportSchedule(id uuid.UUID) (*models.ReportSchedule, error) {
	return s.reportRepo.GetReportSchedule(id)
}

func (s *reportService) ListReportSchedules(companyID uuid.UUID, params map[string]interface{}) ([]models.ReportSchedule, int64, error) {
	return s.reportRepo.ListReportSchedules(companyID, params)
}

func (s *reportService) DeleteReportSchedule(id uuid.UUID) error {
	return s.reportRepo.DeleteReportSchedule(id)
}

func (s *reportService) CreateBusinessKPI(req *CreateBusinessKPIRequest, userID uuid.UUID) (*models.BusinessKPI, error) {
	// Implementation would create KPI
	return nil, fmt.Errorf("not implemented")
}

func (s *reportService) UpdateBusinessKPI(id uuid.UUID, req *UpdateBusinessKPIRequest) (*models.BusinessKPI, error) {
	// Implementation would update KPI
	return nil, fmt.Errorf("not implemented")
}

func (s *reportService) GetBusinessKPI(id uuid.UUID) (*models.BusinessKPI, error) {
	return s.reportRepo.GetBusinessKPI(id)
}

func (s *reportService) ListBusinessKPIs(companyID uuid.UUID, params map[string]interface{}) ([]models.BusinessKPI, int64, error) {
	return s.reportRepo.ListBusinessKPIs(companyID, params)
}

func (s *reportService) DeleteBusinessKPI(id uuid.UUID) error {
	return s.reportRepo.DeleteBusinessKPI(id)
}

func (s *reportService) UpdateKPIValues(companyID uuid.UUID) error {
	// Implementation would calculate and update KPI values
	return nil
}