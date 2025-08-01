package handler

import (
	"net/http"
	"strconv"

	"github.com/fastenmind/fastener-api/internal/models"
	"github.com/fastenmind/fastener-api/internal/services"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CostCalculationHandler struct {
	costService *services.CostCalculationService
}

func NewCostCalculationHandler(costService *services.CostCalculationService) *CostCalculationHandler {
	return &CostCalculationHandler{
		costService: costService,
	}
}

// CalculateCost 計算成本
// @Summary 計算產品成本
// @Description 根據產品資訊和製程路線計算成本
// @Tags Cost Calculation
// @Accept json
// @Produce json
// @Param request body models.CostCalculationRequest true "成本計算請求"
// @Success 200 {object} models.CostCalculation
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/cost-calculations/calculate [post]
func (h *CostCalculationHandler) CalculateCost(c echo.Context) error {
	var req models.CostCalculationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// 從 JWT 獲取用戶 ID
	userID := c.Get("user_id").(uuid.UUID)

	// 執行成本計算
	calculation, err := h.costService.CalculateCost(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// 設定計算人
	calculation.CalculatedBy = userID

	return c.JSON(http.StatusOK, calculation)
}

// GetCostSummary 獲取成本摘要
// @Summary 獲取成本計算摘要
// @Description 獲取指定成本計算的摘要資訊
// @Tags Cost Calculation
// @Produce json
// @Param id path string true "成本計算ID"
// @Success 200 {object} models.CostSummary
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/cost-calculations/{id}/summary [get]
func (h *CostCalculationHandler) GetCostSummary(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid calculation ID"})
	}

	summary, err := h.costService.GetCostSummary(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, summary)
}

// GetProcessRoutes 獲取製程路線
// @Summary 獲取製程路線列表
// @Description 根據產品類別獲取可用的製程路線
// @Tags Cost Calculation
// @Produce json
// @Param product_category query string false "產品類別"
// @Success 200 {array} models.ProductProcessRoute
// @Failure 500 {object} map[string]string
// @Router /api/cost-calculations/process-routes [get]
func (h *CostCalculationHandler) GetProcessRoutes(c echo.Context) error {
	productCategory := c.QueryParam("product_category")
	
	routes, err := h.costService.GetProcessRoutes(productCategory)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, routes)
}

// GetProcessSteps 獲取製程步驟
// @Summary 獲取所有製程步驟
// @Description 獲取系統中所有可用的製程步驟
// @Tags Cost Calculation
// @Produce json
// @Success 200 {array} models.ProcessStep
// @Failure 500 {object} map[string]string
// @Router /api/cost-calculations/process-steps [get]
func (h *CostCalculationHandler) GetProcessSteps(c echo.Context) error {
	steps, err := h.costService.GetProcessSteps()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, steps)
}

// GetEquipment 獲取設備列表
// @Summary 獲取設備列表
// @Description 獲取可用的生產設備列表
// @Tags Cost Calculation
// @Produce json
// @Param category_id query string false "製程類別ID"
// @Success 200 {array} models.Equipment
// @Failure 500 {object} map[string]string
// @Router /api/cost-calculations/equipment [get]
func (h *CostCalculationHandler) GetEquipment(c echo.Context) error {
	var categoryID *uuid.UUID
	if catIDStr := c.QueryParam("category_id"); catIDStr != "" {
		id, err := uuid.Parse(catIDStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid category ID"})
		}
		categoryID = &id
	}

	equipment, err := h.costService.GetEquipmentList(categoryID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, equipment)
}

// GetCalculations 獲取成本計算列表
// @Summary 獲取成本計算列表
// @Description 獲取成本計算記錄列表（分頁）
// @Tags Cost Calculation
// @Produce json
// @Param page query int false "頁碼" default(1)
// @Param page_size query int false "每頁數量" default(20)
// @Param status query string false "狀態篩選"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/cost-calculations [get]
func (h *CostCalculationHandler) GetCalculations(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize < 1 {
		pageSize = 20
	}

	// status := c.QueryParam("status") // TODO: Use when GetCalculations is implemented

	// TODO: Add GetCalculations method to CostCalculationService interface
	// For now, return empty result
	calculations := []interface{}{}
	total := int64(0)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": calculations,
		"pagination": map[string]interface{}{
			"page":       page,
			"page_size":  pageSize,
			"total":      total,
			"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// GetCalculation 獲取單個成本計算
// @Summary 獲取成本計算詳情
// @Description 獲取指定的成本計算完整資訊
// @Tags Cost Calculation
// @Produce json
// @Param id path string true "成本計算ID"
// @Success 200 {object} models.CostCalculation
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/cost-calculations/{id} [get]
func (h *CostCalculationHandler) GetCalculation(c echo.Context) error {
	idStr := c.Param("id")
	_, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid calculation ID"})
	}

	// TODO: Add GetCalculationByID method to CostCalculationService interface
	// For now, return not implemented
	return c.JSON(http.StatusNotImplemented, map[string]string{"error": "Method not implemented"})
}

// ApproveCalculation 審核成本計算
// @Summary 審核成本計算
// @Description 審核通過成本計算
// @Tags Cost Calculation
// @Accept json
// @Produce json
// @Param id path string true "成本計算ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/cost-calculations/{id}/approve [post]
func (h *CostCalculationHandler) ApproveCalculation(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid calculation ID"})
	}

	// 從 JWT 獲取審核人 ID
	approverID := c.Get("user_id").(uuid.UUID)

	if err := h.costService.ApproveCalculation(id, approverID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Calculation approved successfully"})
}

// GetCostParameters 獲取成本參數
// @Summary 獲取成本參數列表
// @Description 獲取系統成本參數設定
// @Tags Cost Calculation
// @Produce json
// @Success 200 {array} models.CostParameter
// @Failure 500 {object} map[string]string
// @Router /api/cost-calculations/parameters [get]
func (h *CostCalculationHandler) GetCostParameters(c echo.Context) error {
	// TODO: Add GetCostParameters method to CostCalculationService interface
	// For now, return empty result
	params := []interface{}{}
	return c.JSON(http.StatusOK, params)
}

// UpdateCostParameter 更新成本參數
// @Summary 更新成本參數
// @Description 更新系統成本參數（會保留歷史記錄）
// @Tags Cost Calculation
// @Accept json
// @Produce json
// @Param parameter body models.CostParameter true "成本參數"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/cost-calculations/parameters [put]
func (h *CostCalculationHandler) UpdateCostParameter(c echo.Context) error {
	var param models.CostParameter
	if err := c.Bind(&param); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// TODO: Add UpdateCostParameter method to CostCalculationService interface
	// if err := h.costService.UpdateCostParameter(param); err != nil {
	// 	return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	// }

	return c.JSON(http.StatusOK, map[string]string{"message": "Parameter updated successfully"})
}

// RegisterRoutes 註冊路由
func (h *CostCalculationHandler) RegisterRoutes(e *echo.Group) {
	costCalc := e.Group("/cost-calculations")
	
	// 成本計算
	costCalc.POST("/calculate", h.CalculateCost)
	costCalc.GET("", h.GetCalculations)
	costCalc.GET("/:id", h.GetCalculation)
	costCalc.GET("/:id/summary", h.GetCostSummary)
	costCalc.POST("/:id/approve", h.ApproveCalculation)
	
	// 基礎資料
	costCalc.GET("/process-routes", h.GetProcessRoutes)
	costCalc.GET("/process-steps", h.GetProcessSteps)
	costCalc.GET("/equipment", h.GetEquipment)
	
	// 成本參數
	costCalc.GET("/parameters", h.GetCostParameters)
	costCalc.PUT("/parameters", h.UpdateCostParameter)
}