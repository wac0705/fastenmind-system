package handler

import (
	"net/http"
	"strconv"

	"github.com/fastenmind/fastener-api/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ReportHandler struct {
	reportService service.ReportService
}

func NewReportHandler(reportService service.ReportService) *ReportHandler {
	return &ReportHandler{
		reportService: reportService,
	}
}

// Report operations
func (h *ReportHandler) CreateReport(c echo.Context) error {
	var req service.CreateReportRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	report, err := h.reportService.CreateReport(&req, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, report)
}

func (h *ReportHandler) UpdateReport(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid report ID"})
	}

	var req service.UpdateReportRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	report, err := h.reportService.UpdateReport(id, &req, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, report)
}

func (h *ReportHandler) GetReport(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid report ID"})
	}

	report, err := h.reportService.GetReport(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Report not found"})
	}

	return c.JSON(http.StatusOK, report)
}

func (h *ReportHandler) ListReports(c echo.Context) error {
	companyID, err := getCompanyIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	params := make(map[string]interface{})

	if page := c.QueryParam("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil {
			params["page"] = p
		}
	}

	if pageSize := c.QueryParam("page_size"); pageSize != "" {
		if ps, err := strconv.Atoi(pageSize); err == nil {
			params["page_size"] = ps
		}
	}

	if search := c.QueryParam("search"); search != "" {
		params["search"] = search
	}

	if category := c.QueryParam("category"); category != "" {
		params["category"] = category
	}

	if reportType := c.QueryParam("type"); reportType != "" {
		params["type"] = reportType
	}

	if status := c.QueryParam("status"); status != "" {
		params["status"] = status
	}

	if isScheduled := c.QueryParam("is_scheduled"); isScheduled != "" {
		if scheduled, err := strconv.ParseBool(isScheduled); err == nil {
			params["is_scheduled"] = scheduled
		}
	}

	if isPublic := c.QueryParam("is_public"); isPublic != "" {
		if public, err := strconv.ParseBool(isPublic); err == nil {
			params["is_public"] = public
		}
	}

	if createdBy := c.QueryParam("created_by"); createdBy != "" {
		params["created_by"] = createdBy
	}

	if sortBy := c.QueryParam("sort_by"); sortBy != "" {
		params["sort_by"] = sortBy
	}

	if sortOrder := c.QueryParam("sort_order"); sortOrder != "" {
		params["sort_order"] = sortOrder
	}

	reports, total, err := h.reportService.ListReports(companyID, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":  reports,
		"total": total,
	})
}

func (h *ReportHandler) DeleteReport(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid report ID"})
	}

	if err := h.reportService.DeleteReport(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *ReportHandler) DuplicateReport(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid report ID"})
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	report, err := h.reportService.DuplicateReport(id, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, report)
}

// Report Template operations
func (h *ReportHandler) CreateReportTemplate(c echo.Context) error {
	var req service.CreateReportTemplateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	template, err := h.reportService.CreateReportTemplate(&req, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, template)
}

func (h *ReportHandler) UpdateReportTemplate(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid template ID"})
	}

	var req service.UpdateReportTemplateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	template, err := h.reportService.UpdateReportTemplate(id, &req, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, template)
}

func (h *ReportHandler) GetReportTemplate(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid template ID"})
	}

	template, err := h.reportService.GetReportTemplate(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Template not found"})
	}

	return c.JSON(http.StatusOK, template)
}

func (h *ReportHandler) ListReportTemplates(c echo.Context) error {
	var companyID *uuid.UUID
	if cid, err := getCompanyIDFromContext(c); err == nil {
		companyID = &cid
	}

	params := make(map[string]interface{})

	if page := c.QueryParam("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil {
			params["page"] = p
		}
	}

	if pageSize := c.QueryParam("page_size"); pageSize != "" {
		if ps, err := strconv.Atoi(pageSize); err == nil {
			params["page_size"] = ps
		}
	}

	if search := c.QueryParam("search"); search != "" {
		params["search"] = search
	}

	if category := c.QueryParam("category"); category != "" {
		params["category"] = category
	}

	if templateType := c.QueryParam("type"); templateType != "" {
		params["type"] = templateType
	}

	if isSystemTemplate := c.QueryParam("is_system_template"); isSystemTemplate != "" {
		if system, err := strconv.ParseBool(isSystemTemplate); err == nil {
			params["is_system_template"] = system
		}
	}

	if industry := c.QueryParam("industry"); industry != "" {
		params["industry"] = industry
	}

	if language := c.QueryParam("language"); language != "" {
		params["language"] = language
	}

	if sortBy := c.QueryParam("sort_by"); sortBy != "" {
		params["sort_by"] = sortBy
	}

	if sortOrder := c.QueryParam("sort_order"); sortOrder != "" {
		params["sort_order"] = sortOrder
	}

	templates, total, err := h.reportService.ListReportTemplates(companyID, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":  templates,
		"total": total,
	})
}

func (h *ReportHandler) DeleteReportTemplate(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid template ID"})
	}

	if err := h.reportService.DeleteReportTemplate(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// Report Execution operations
func (h *ReportHandler) ExecuteReport(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid report ID"})
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	var params map[string]interface{}
	if err := c.Bind(&params); err != nil {
		params = make(map[string]interface{})
	}

	execution, err := h.reportService.ExecuteReport(id, params, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusAccepted, execution)
}

func (h *ReportHandler) GetReportExecution(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid execution ID"})
	}

	execution, err := h.reportService.GetReportExecution(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Execution not found"})
	}

	return c.JSON(http.StatusOK, execution)
}

func (h *ReportHandler) ListReportExecutions(c echo.Context) error {
	reportID, err := uuid.Parse(c.Param("report_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid report ID"})
	}

	params := make(map[string]interface{})

	if page := c.QueryParam("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil {
			params["page"] = p
		}
	}

	if pageSize := c.QueryParam("page_size"); pageSize != "" {
		if ps, err := strconv.Atoi(pageSize); err == nil {
			params["page_size"] = ps
		}
	}

	if status := c.QueryParam("status"); status != "" {
		params["status"] = status
	}

	if isScheduled := c.QueryParam("is_scheduled"); isScheduled != "" {
		if scheduled, err := strconv.ParseBool(isScheduled); err == nil {
			params["is_scheduled"] = scheduled
		}
	}

	if triggerType := c.QueryParam("trigger_type"); triggerType != "" {
		params["trigger_type"] = triggerType
	}

	if executedBy := c.QueryParam("executed_by"); executedBy != "" {
		params["executed_by"] = executedBy
	}

	if startDate := c.QueryParam("start_date"); startDate != "" {
		params["start_date"] = startDate
	}

	if endDate := c.QueryParam("end_date"); endDate != "" {
		params["end_date"] = endDate
	}

	executions, total, err := h.reportService.ListReportExecutions(reportID, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":  executions,
		"total": total,
	})
}

func (h *ReportHandler) CancelReportExecution(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid execution ID"})
	}

	if err := h.reportService.CancelReportExecution(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Execution cancelled successfully"})
}

func (h *ReportHandler) DownloadReportResult(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid execution ID"})
	}

	content, filename, err := h.reportService.DownloadReportResult(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	c.Response().Header().Set("Content-Disposition", "attachment; filename="+filename)
	return c.Blob(http.StatusOK, "application/octet-stream", content)
}

// Business operations
func (h *ReportHandler) GetReportDashboard(c echo.Context) error {
	companyID, err := getCompanyIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	dashboard, err := h.reportService.GetReportDashboardData(companyID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, dashboard)
}

func (h *ReportHandler) GetReportStatistics(c echo.Context) error {
	companyID, err := getCompanyIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	stats, err := h.reportService.GetReportStatistics(companyID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, stats)
}

func (h *ReportHandler) GetPopularReports(c echo.Context) error {
	companyID, err := getCompanyIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	limit := 10
	if l := c.QueryParam("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	reports, err := h.reportService.GetPopularReports(companyID, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, reports)
}

func (h *ReportHandler) GetRecentExecutions(c echo.Context) error {
	companyID, err := getCompanyIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	limit := 10
	if l := c.QueryParam("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	executions, err := h.reportService.GetRecentExecutions(companyID, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, executions)
}

func (h *ReportHandler) GenerateReportFromTemplate(c echo.Context) error {
	templateID, err := uuid.Parse(c.Param("template_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid template ID"})
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	var params map[string]interface{}
	if err := c.Bind(&params); err != nil {
		params = make(map[string]interface{})
	}

	report, err := h.reportService.GenerateReportFromTemplate(templateID, params, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, report)
}

func (h *ReportHandler) ExportReport(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid report ID"})
	}

	format := c.QueryParam("format")
	if format == "" {
		format = "json" // default format
	}

	var params map[string]interface{}
	if err := c.Bind(&params); err != nil {
		params = make(map[string]interface{})
	}

	content, filename, err := h.reportService.ExportReport(id, format, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	c.Response().Header().Set("Content-Disposition", "attachment; filename="+filename)
	
	var contentType string
	switch format {
	case "json":
		contentType = "application/json"
	case "csv":
		contentType = "text/csv"
	default:
		contentType = "application/octet-stream"
	}

	return c.Blob(http.StatusOK, contentType, content)
}

func (h *ReportHandler) ImportReports(c echo.Context) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	// Get file from multipart form
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No file provided"})
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to open file"})
	}
	defer src.Close()

	// Read file content
	data := make([]byte, file.Size)
	if _, err := src.Read(data); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to read file"})
	}

	result, err := h.reportService.ImportReports(data, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

// Report Subscription operations (placeholder implementations)
func (h *ReportHandler) CreateReportSubscription(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{"error": "Not implemented"})
}

func (h *ReportHandler) UpdateReportSubscription(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{"error": "Not implemented"})
}

func (h *ReportHandler) GetReportSubscription(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid subscription ID"})
	}

	subscription, err := h.reportService.GetReportSubscription(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Subscription not found"})
	}

	return c.JSON(http.StatusOK, subscription)
}

func (h *ReportHandler) ListReportSubscriptions(c echo.Context) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	params := make(map[string]interface{})
	
	if page := c.QueryParam("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil {
			params["page"] = p
		}
	}

	if pageSize := c.QueryParam("page_size"); pageSize != "" {
		if ps, err := strconv.Atoi(pageSize); err == nil {
			params["page_size"] = ps
		}
	}

	subscriptions, total, err := h.reportService.ListReportSubscriptions(userID, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":  subscriptions,
		"total": total,
	})
}

func (h *ReportHandler) DeleteReportSubscription(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid subscription ID"})
	}

	if err := h.reportService.DeleteReportSubscription(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// Helper functions (these should be defined elsewhere in the codebase)
func getUserIDFromContext(c echo.Context) (uuid.UUID, error) {
	userID := c.Get("user_id")
	if userID == nil {
		return uuid.Nil, echo.NewHTTPError(http.StatusUnauthorized, "User ID not found in context")
	}

	switch v := userID.(type) {
	case string:
		return uuid.Parse(v)
	case uuid.UUID:
		return v, nil
	default:
		return uuid.Nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid user ID format")
	}
}

func getCompanyIDFromContext(c echo.Context) (uuid.UUID, error) {
	companyID := c.Get("company_id")
	if companyID == nil {
		return uuid.Nil, echo.NewHTTPError(http.StatusUnauthorized, "Company ID not found in context")
	}

	switch v := companyID.(type) {
	case string:
		return uuid.Parse(v)
	case uuid.UUID:
		return v, nil
	default:
		return uuid.Nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid company ID format")
	}
}