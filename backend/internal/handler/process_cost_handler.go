package handler

import (
	"net/http"
	"strconv"

	"github.com/fastenmind/fastener-api/internal/models"
	"github.com/fastenmind/fastener-api/internal/service"
	"github.com/fastenmind/fastener-api/pkg/response"

	"github.com/labstack/echo/v4"
)

type ProcessCostHandler struct {
	processCostService *service.ProcessCostService
}

func NewProcessCostHandler(processCostService *service.ProcessCostService) *ProcessCostHandler {
	return &ProcessCostHandler{
		processCostService: processCostService,
	}
}

// GetCostTemplates 獲取成本模板列表
func (h *ProcessCostHandler) GetCostTemplates(c echo.Context) error {
	companyID := c.Get("company_id").(string)
	processType := c.QueryParam("process_type")
	category := c.QueryParam("category")
	
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	
	templates, total, err := h.processCostService.GetCostTemplates(companyID, processType, category, page, limit)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to get cost templates", err)
	}
	
	return response.SuccessWithPagination(c, "Cost templates retrieved successfully", templates, int(total), page, limit)
}

// CreateCostTemplate 創建成本模板
func (h *ProcessCostHandler) CreateCostTemplate(c echo.Context) error {
	var req models.ProcessCostTemplateNew
	
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request", err)
	}
	
	if err := c.Validate(req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Validation failed", err)
	}
	
	companyID := c.Get("company_id").(string)
	userID := c.Get("user_id").(string)
	
	template, err := h.processCostService.CreateCostTemplate(&req, companyID, userID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to create cost template", err)
	}
	
	return response.Success(c, template, "Cost template created successfully")
}

// UpdateCostTemplate 更新成本模板
func (h *ProcessCostHandler) UpdateCostTemplate(c echo.Context) error {
	templateID := c.Param("id")
	
	var req models.ProcessCostTemplateNew
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request", err)
	}
	
	companyID := c.Get("company_id").(string)
	userID := c.Get("user_id").(string)
	
	template, err := h.processCostService.UpdateCostTemplate(templateID, &req, companyID, userID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to update cost template", err)
	}
	
	return response.Success(c, template, "Cost template updated successfully")
}

// DeleteCostTemplate 刪除成本模板
func (h *ProcessCostHandler) DeleteCostTemplate(c echo.Context) error {
	templateID := c.Param("id")
	companyID := c.Get("company_id").(string)
	
	if err := h.processCostService.DeleteCostTemplate(templateID, companyID); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to delete cost template", err)
	}
	
	return response.Success(c, nil, "Cost template deleted successfully")
}

// CalculateProcessCost 計算製程成本
func (h *ProcessCostHandler) CalculateProcessCost(c echo.Context) error {
	var req models.ProcessCostCalculationRequestNew
	
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request", err)
	}
	
	if err := c.Validate(req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Validation failed", err)
	}
	
	companyID := c.Get("company_id").(string)
	
	result, err := h.processCostService.CalculateProcessCost(&req, companyID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to calculate process cost", err)
	}
	
	return response.Success(c, result, "Process cost calculated successfully")
}

// GetCostHistory 獲取成本計算歷史
func (h *ProcessCostHandler) GetCostHistory(c echo.Context) error {
	companyID := c.Get("company_id").(string)
	inquiryID := c.QueryParam("inquiry_id")
	productID := c.QueryParam("product_id")
	
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	
	history, total, err := h.processCostService.GetCostHistory(companyID, inquiryID, productID, page, limit)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to get cost history", err)
	}
	
	return response.SuccessWithPagination(c, "Cost history retrieved successfully", history, int(total), page, limit)
}

// GetMaterialCosts 獲取材料成本
func (h *ProcessCostHandler) GetMaterialCosts(c echo.Context) error {
	companyID := c.Get("company_id").(string)
	materialType := c.QueryParam("type")
	
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	
	materials, total, err := h.processCostService.GetMaterialCosts(companyID, materialType, page, limit)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to get material costs", err)
	}
	
	return response.SuccessWithPagination(c, "Material costs retrieved successfully", materials, int(total), page, limit)
}

// UpdateMaterialCost 更新材料成本
func (h *ProcessCostHandler) UpdateMaterialCost(c echo.Context) error {
	materialID := c.Param("id")
	
	var req models.MaterialCostNew
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request", err)
	}
	
	companyID := c.Get("company_id").(string)
	userID := c.Get("user_id").(string)
	
	material, err := h.processCostService.UpdateMaterialCost(materialID, &req, companyID, userID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to update material cost", err)
	}
	
	return response.Success(c, material, "Material cost updated successfully")
}

// GetProcessingRates 獲取加工費率
func (h *ProcessCostHandler) GetProcessingRates(c echo.Context) error {
	companyID := c.Get("company_id").(string)
	processType := c.QueryParam("process_type")
	equipmentID := c.QueryParam("equipment_id")
	
	rates, err := h.processCostService.GetProcessingRates(companyID, processType, equipmentID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to get processing rates", err)
	}
	
	return response.Success(c, rates, "Processing rates retrieved successfully")
}

// UpdateProcessingRate 更新加工費率
func (h *ProcessCostHandler) UpdateProcessingRate(c echo.Context) error {
	rateID := c.Param("id")
	
	var req models.ProcessingRate
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request", err)
	}
	
	companyID := c.Get("company_id").(string)
	userID := c.Get("user_id").(string)
	
	rate, err := h.processCostService.UpdateProcessingRate(rateID, &req, companyID, userID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to update processing rate", err)
	}
	
	return response.Success(c, rate, "Processing rate updated successfully")
}

// BatchCalculateCost 批量計算成本
func (h *ProcessCostHandler) BatchCalculateCost(c echo.Context) error {
	var req models.BatchCostCalculationRequest
	
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request", err)
	}
	
	if err := c.Validate(req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Validation failed", err)
	}
	
	companyID := c.Get("company_id").(string)
	userID := c.Get("user_id").(string)
	
	results, err := h.processCostService.BatchCalculateCost(&req, companyID, userID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to batch calculate costs", err)
	}
	
	return response.Success(c, results, "Batch cost calculation completed")
}

// GetCostAnalysis 獲取成本分析
func (h *ProcessCostHandler) GetCostAnalysis(c echo.Context) error {
	companyID := c.Get("company_id").(string)
	analysisType := c.QueryParam("type") // material, process, total
	period := c.QueryParam("period")     // daily, weekly, monthly, yearly
	startDate := c.QueryParam("start_date")
	endDate := c.QueryParam("end_date")
	
	analysis, err := h.processCostService.GetCostAnalysis(companyID, analysisType, period, startDate, endDate)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to get cost analysis", err)
	}
	
	return response.Success(c, analysis, "Cost analysis retrieved successfully")
}

// ExportCostReport 導出成本報告
func (h *ProcessCostHandler) ExportCostReport(c echo.Context) error {
	companyID := c.Get("company_id").(string)
	format := c.QueryParam("format") // excel, pdf, csv
	reportType := c.QueryParam("type")
	startDate := c.QueryParam("start_date")
	endDate := c.QueryParam("end_date")
	
	if format == "" {
		format = "excel"
	}
	
	reportData, fileName, err := h.processCostService.ExportCostReport(companyID, format, reportType, startDate, endDate)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to export cost report", err)
	}
	
	// 設置響應頭
	c.Response().Header().Set("Content-Type", getContentType(format))
	c.Response().Header().Set("Content-Disposition", "attachment; filename="+fileName)
	
	return c.Blob(http.StatusOK, getContentType(format), reportData)
}

// GetCostSettings 獲取成本設定
func (h *ProcessCostHandler) GetCostSettings(c echo.Context) error {
	companyID := c.Get("company_id").(string)
	
	settings, err := h.processCostService.GetCostSettings(companyID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to get cost settings", err)
	}
	
	return response.Success(c, settings, "Cost settings retrieved successfully")
}

// UpdateCostSettings 更新成本設定
func (h *ProcessCostHandler) UpdateCostSettings(c echo.Context) error {
	var req models.CostSettings
	
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request", err)
	}
	
	companyID := c.Get("company_id").(string)
	userID := c.Get("user_id").(string)
	
	settings, err := h.processCostService.UpdateCostSettings(&req, companyID, userID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to update cost settings", err)
	}
	
	return response.Success(c, settings, "Cost settings updated successfully")
}

func getContentType(format string) string {
	switch format {
	case "pdf":
		return "application/pdf"
	case "csv":
		return "text/csv"
	default:
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	}
}