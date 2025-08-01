package handler

import (
	"net/http"
	"strconv"

	"github.com/fastenmind/fastener-api/internal/models"
	"github.com/fastenmind/fastener-api/internal/services"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AssignmentHandler struct {
	assignmentService *services.AssignmentService
}

func NewAssignmentHandler(assignmentService *services.AssignmentService) *AssignmentHandler {
	return &AssignmentHandler{
		assignmentService: assignmentService,
	}
}

// AutoAssign 自動分派詢價單
// @Summary 自動分派詢價單
// @Description 根據規則自動分派詢價單給工程師
// @Tags Assignment
// @Accept json
// @Produce json
// @Param request body models.AutoAssignmentRequest true "自動分派請求"
// @Success 200 {object} models.AssignmentHistory
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/assignments/auto [post]
func (h *AssignmentHandler) AutoAssign(c echo.Context) error {
	var req models.AutoAssignmentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	history, err := h.assignmentService.AutoAssignInquiry(req.InquiryID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, history)
}

// ManualAssign 手動分派詢價單
// @Summary 手動分派詢價單
// @Description 手動指定工程師處理詢價單
// @Tags Assignment
// @Accept json
// @Produce json
// @Param request body models.AssignmentRequest true "手動分派請求"
// @Success 200 {object} models.AssignmentHistory
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/assignments/manual [post]
func (h *AssignmentHandler) ManualAssign(c echo.Context) error {
	var req models.AssignmentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// 從 JWT 中獲取當前用戶 ID
	userID := c.Get("user_id").(uuid.UUID)

	history, err := h.assignmentService.ManualAssignInquiry(req, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, history)
}

// SelfSelect 工程師自選詢價單
// @Summary 工程師自選詢價單
// @Description 工程師主動選擇處理某個詢價單
// @Tags Assignment
// @Accept json
// @Produce json
// @Param inquiry_id path string true "詢價單ID"
// @Success 200 {object} models.AssignmentHistory
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/assignments/self-select/{inquiry_id} [post]
func (h *AssignmentHandler) SelfSelect(c echo.Context) error {
	inquiryIDStr := c.Param("inquiry_id")
	inquiryID, err := uuid.Parse(inquiryIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid inquiry ID"})
	}

	// 從 JWT 中獲取當前工程師 ID
	engineerID := c.Get("user_id").(uuid.UUID)

	history, err := h.assignmentService.SelfSelectInquiry(inquiryID, engineerID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, history)
}

// GetWorkloadStats 獲取工程師工作量統計
// @Summary 獲取工程師工作量統計
// @Description 獲取所有工程師的工作量統計資訊
// @Tags Assignment
// @Produce json
// @Success 200 {array} models.EngineerWorkloadStats
// @Failure 500 {object} map[string]string
// @Router /api/assignments/workload-stats [get]
func (h *AssignmentHandler) GetWorkloadStats(c echo.Context) error {
	stats, err := h.assignmentService.GetEngineerWorkloadStats()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, stats)
}

// GetAssignmentHistory 獲取分派歷史
// @Summary 獲取分派歷史
// @Description 獲取詢價單或工程師的分派歷史記錄
// @Tags Assignment
// @Produce json
// @Param inquiry_id query string false "詢價單ID"
// @Param engineer_id query string false "工程師ID"
// @Param limit query int false "限制數量"
// @Success 200 {array} models.AssignmentHistory
// @Failure 500 {object} map[string]string
// @Router /api/assignments/history [get]
func (h *AssignmentHandler) GetAssignmentHistory(c echo.Context) error {
	var inquiryID, engineerID *uuid.UUID
	limit := 50 // 預設限制

	if id := c.QueryParam("inquiry_id"); id != "" {
		parsed, err := uuid.Parse(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid inquiry ID"})
		}
		inquiryID = &parsed
	}

	if id := c.QueryParam("engineer_id"); id != "" {
		parsed, err := uuid.Parse(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid engineer ID"})
		}
		engineerID = &parsed
	}

	if l := c.QueryParam("limit"); l != "" {
		// Parse limit from string
		var err error
		limit, err = strconv.Atoi(l)
		if err != nil {
			limit = 50
		}
	}

	history, err := h.assignmentService.GetAssignmentHistory(inquiryID, engineerID, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, history)
}

// GetAssignmentRules 獲取分派規則
// @Summary 獲取分派規則
// @Description 獲取所有分派規則
// @Tags Assignment
// @Produce json
// @Success 200 {array} models.AssignmentRule
// @Failure 500 {object} map[string]string
// @Router /api/assignments/rules [get]
func (h *AssignmentHandler) GetAssignmentRules(c echo.Context) error {
	rules, err := h.assignmentService.GetAssignmentRules()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, rules)
}

// UpdateAssignmentRule 更新分派規則
// @Summary 更新分派規則
// @Description 更新指定的分派規則
// @Tags Assignment
// @Accept json
// @Produce json
// @Param id path string true "規則ID"
// @Param rule body models.AssignmentRule true "規則資料"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/assignments/rules/{id} [put]
func (h *AssignmentHandler) UpdateAssignmentRule(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid rule ID"})
	}

	var rule models.AssignmentRule
	if err := c.Bind(&rule); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := h.assignmentService.UpdateAssignmentRule(id, rule); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Rule updated successfully"})
}

// GetEngineerCapabilities 獲取工程師能力
// @Summary 獲取工程師能力
// @Description 獲取指定工程師的能力列表
// @Tags Assignment
// @Produce json
// @Param engineer_id path string true "工程師ID"
// @Success 200 {array} models.EngineerCapability
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/assignments/engineers/{engineer_id}/capabilities [get]
func (h *AssignmentHandler) GetEngineerCapabilities(c echo.Context) error {
	engineerIDStr := c.Param("engineer_id")
	engineerID, err := uuid.Parse(engineerIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid engineer ID"})
	}

	capabilities, err := h.assignmentService.GetEngineerCapabilities(engineerID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, capabilities)
}

// UpdateEngineerCapability 更新工程師能力
// @Summary 更新工程師能力
// @Description 更新工程師的能力設定
// @Tags Assignment
// @Accept json
// @Produce json
// @Param capability body models.EngineerCapability true "能力資料"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/assignments/capabilities [put]
func (h *AssignmentHandler) UpdateEngineerCapability(c echo.Context) error {
	var capability models.EngineerCapability
	if err := c.Bind(&capability); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := h.assignmentService.UpdateEngineerCapability(capability); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Capability updated successfully"})
}

// UpdateEngineerPreference 更新工程師偏好
// @Summary 更新工程師偏好
// @Description 更新工程師的偏好設定
// @Tags Assignment
// @Accept json
// @Produce json
// @Param preference body models.EngineerPreference true "偏好資料"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/assignments/preferences [put]
func (h *AssignmentHandler) UpdateEngineerPreference(c echo.Context) error {
	var preference models.EngineerPreference
	if err := c.Bind(&preference); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// 確保只能更新自己的偏好
	userID := c.Get("user_id").(uuid.UUID)
	preference.EngineerID = userID

	if err := h.assignmentService.UpdateEngineerPreference(preference); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Preference updated successfully"})
}

// RegisterRoutes 註冊路由
func (h *AssignmentHandler) RegisterRoutes(e *echo.Group) {
	assignments := e.Group("/assignments")
	
	// 分派操作
	assignments.POST("/auto", h.AutoAssign)
	assignments.POST("/manual", h.ManualAssign)
	assignments.POST("/self-select/:inquiry_id", h.SelfSelect)
	
	// 查詢操作
	assignments.GET("/workload-stats", h.GetWorkloadStats)
	assignments.GET("/history", h.GetAssignmentHistory)
	
	// 規則管理
	assignments.GET("/rules", h.GetAssignmentRules)
	assignments.PUT("/rules/:id", h.UpdateAssignmentRule)
	
	// 工程師能力管理
	assignments.GET("/engineers/:engineer_id/capabilities", h.GetEngineerCapabilities)
	assignments.PUT("/capabilities", h.UpdateEngineerCapability)
	
	// 工程師偏好管理
	assignments.PUT("/preferences", h.UpdateEngineerPreference)
}