package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/fastenmind/fastener-api/internal/service"
	"github.com/fastenmind/fastener-api/pkg/response"

	"github.com/labstack/echo/v4"
)

type EngineerAssignmentHandler struct {
	assignmentService *service.EngineerAssignmentService
}

func NewEngineerAssignmentHandler(assignmentService *service.EngineerAssignmentService) *EngineerAssignmentHandler {
	return &EngineerAssignmentHandler{
		assignmentService: assignmentService,
	}
}

// GetAvailableEngineers 獲取可用的工程師列表
func (h *EngineerAssignmentHandler) GetAvailableEngineers(c echo.Context) error {
	companyID := c.Get("company_id").(string)
	inquiryID := c.QueryParam("inquiry_id")
	
	engineers, err := h.assignmentService.GetAvailableEngineers(companyID, inquiryID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to get available engineers", err)
	}
	
	return response.Success(c, engineers, "Available engineers retrieved successfully")
}

// AssignEngineer 分配工程師到詢價單
func (h *EngineerAssignmentHandler) AssignEngineer(c echo.Context) error {
	var req struct {
		InquiryID  string `json:"inquiry_id" validate:"required"`
		EngineerID string `json:"engineer_id" validate:"required"`
		Priority   string `json:"priority"`
		DueDate    string `json:"due_date"`
		Notes      string `json:"notes"`
	}
	
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request", err)
	}
	
	if err := c.Validate(req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Validation failed", err)
	}
	
	userID := c.Get("user_id").(string)
	companyID := c.Get("company_id").(string)
	
	assignment, err := h.assignmentService.AssignEngineer(companyID, req.InquiryID, req.EngineerID, userID, req.Priority, req.DueDate, req.Notes)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to assign engineer", err)
	}
	
	return response.Success(c, assignment, "Engineer assigned successfully")
}

// ReassignEngineer 重新分配工程師
func (h *EngineerAssignmentHandler) ReassignEngineer(c echo.Context) error {
	assignmentID := c.Param("id")
	
	var req struct {
		NewEngineerID string `json:"new_engineer_id" validate:"required"`
		Reason        string `json:"reason" validate:"required"`
	}
	
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request", err)
	}
	
	if err := c.Validate(req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Validation failed", err)
	}
	
	userID := c.Get("user_id").(string)
	companyID := c.Get("company_id").(string)
	
	assignment, err := h.assignmentService.ReassignEngineer(companyID, assignmentID, req.NewEngineerID, userID, req.Reason)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to reassign engineer", err)
	}
	
	return response.Success(c, assignment, "Engineer reassigned successfully")
}

// GetAssignmentHistory 獲取分派歷史記錄
func (h *EngineerAssignmentHandler) GetAssignmentHistory(c echo.Context) error {
	inquiryID := c.QueryParam("inquiry_id")
	engineerID := c.QueryParam("engineer_id")
	companyID := c.Get("company_id").(string)
	
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	
	history, total, err := h.assignmentService.GetAssignmentHistory(companyID, inquiryID, engineerID, page, limit)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to get assignment history", err)
	}
	
	return response.SuccessWithPagination(c, "Assignment history retrieved successfully", history, int(total), page, limit)
}

// GetEngineerWorkload 獲取工程師工作負載
func (h *EngineerAssignmentHandler) GetEngineerWorkload(c echo.Context) error {
	companyID := c.Get("company_id").(string)
	startDate := c.QueryParam("start_date")
	endDate := c.QueryParam("end_date")
	
	if startDate == "" {
		startDate = time.Now().Format("2006-01-02")
	}
	if endDate == "" {
		endDate = time.Now().AddDate(0, 1, 0).Format("2006-01-02")
	}
	
	workload, err := h.assignmentService.GetEngineerWorkload(companyID, startDate, endDate)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to get engineer workload", err)
	}
	
	return response.Success(c, workload, "Engineer workload retrieved successfully")
}

// UpdateAssignmentStatus 更新分派狀態
func (h *EngineerAssignmentHandler) UpdateAssignmentStatus(c echo.Context) error {
	assignmentID := c.Param("id")
	
	var req struct {
		Status string `json:"status" validate:"required,oneof=pending in_progress completed cancelled"`
		Notes  string `json:"notes"`
	}
	
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request", err)
	}
	
	if err := c.Validate(req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Validation failed", err)
	}
	
	userID := c.Get("user_id").(string)
	companyID := c.Get("company_id").(string)
	
	assignment, err := h.assignmentService.UpdateAssignmentStatus(companyID, assignmentID, req.Status, userID, req.Notes)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to update assignment status", err)
	}
	
	return response.Success(c, assignment, "Assignment status updated successfully")
}

// GetAssignmentStats 獲取分派統計數據
func (h *EngineerAssignmentHandler) GetAssignmentStats(c echo.Context) error {
	companyID := c.Get("company_id").(string)
	period := c.QueryParam("period") // daily, weekly, monthly
	if period == "" {
		period = "monthly"
	}
	
	stats, err := h.assignmentService.GetAssignmentStats(companyID, period)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to get assignment stats", err)
	}
	
	return response.Success(c, stats, "Assignment stats retrieved successfully")
}

// AutoAssignEngineer 自動分配工程師
func (h *EngineerAssignmentHandler) AutoAssignEngineer(c echo.Context) error {
	var req struct {
		InquiryID string `json:"inquiry_id" validate:"required"`
		Rules     struct {
			ConsiderWorkload    bool `json:"consider_workload"`
			ConsiderExpertise   bool `json:"consider_expertise"`
			ConsiderAvailability bool `json:"consider_availability"`
			MaxAssignments      int  `json:"max_assignments"`
		} `json:"rules"`
	}
	
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request", err)
	}
	
	if err := c.Validate(req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Validation failed", err)
	}
	
	userID := c.Get("user_id").(string)
	companyID := c.Get("company_id").(string)
	
	assignment, err := h.assignmentService.AutoAssignEngineer(companyID, req.InquiryID, userID, req.Rules)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to auto-assign engineer", err)
	}
	
	return response.Success(c, assignment, "Engineer auto-assigned successfully")
}