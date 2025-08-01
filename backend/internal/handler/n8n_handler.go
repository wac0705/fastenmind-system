package handler

import (
	"net/http"
	"strconv"

	"github.com/fastenmind/fastener-api/internal/middleware"
	"github.com/fastenmind/fastener-api/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type N8NHandler struct {
	service service.N8NService
}

func NewN8NHandler(service service.N8NService) *N8NHandler {
	return &N8NHandler{service: service}
}

// TestConnection godoc
// @Summary Test N8N connection
// @Description Test connection to N8N instance
// @Tags N8N
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/n8n/test-connection [get]
func (h *N8NHandler) TestConnection(c echo.Context) error {
	connected, version, err := h.service.TestConnection()
	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"connected": false,
			"message":   err.Error(),
		})
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"connected": connected,
		"version":   version,
	})
}

// GetAvailableWorkflows godoc
// @Summary Get available N8N workflows
// @Description Get list of available workflows from N8N
// @Tags N8N
// @Accept json
// @Produce json
// @Success 200 {array} service.N8NWorkflowInfo
// @Router /api/n8n/available-workflows [get]
func (h *N8NHandler) GetAvailableWorkflows(c echo.Context) error {
	workflows, err := h.service.GetAvailableWorkflows()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, workflows)
}

// CreateWorkflow godoc
// @Summary Create workflow
// @Description Create a new workflow configuration
// @Tags N8N
// @Accept json
// @Produce json
// @Param request body service.CreateWorkflowRequest true "Create workflow request"
// @Success 201 {object} models.N8NWorkflow
// @Router /api/n8n/workflows [post]
func (h *N8NHandler) CreateWorkflow(c echo.Context) error {
	var req service.CreateWorkflowRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	userClaims := c.Get("user").(*middleware.Claims)
	
	workflow, err := h.service.CreateWorkflow(userClaims.CompanyID, userClaims.UserID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusCreated, workflow)
}

// UpdateWorkflow godoc
// @Summary Update workflow
// @Description Update workflow configuration
// @Tags N8N
// @Accept json
// @Produce json
// @Param id path string true "Workflow ID"
// @Param request body service.UpdateWorkflowRequest true "Update workflow request"
// @Success 200 {object} models.N8NWorkflow
// @Router /api/n8n/workflows/{id} [put]
func (h *N8NHandler) UpdateWorkflow(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid workflow ID")
	}
	
	var req service.UpdateWorkflowRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	workflow, err := h.service.UpdateWorkflow(id, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusOK, workflow)
}

// DeleteWorkflow godoc
// @Summary Delete workflow
// @Description Delete a workflow configuration
// @Tags N8N
// @Accept json
// @Produce json
// @Param id path string true "Workflow ID"
// @Success 204
// @Router /api/n8n/workflows/{id} [delete]
func (h *N8NHandler) DeleteWorkflow(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid workflow ID")
	}
	
	if err := h.service.DeleteWorkflow(id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.NoContent(http.StatusNoContent)
}

// GetWorkflow godoc
// @Summary Get workflow
// @Description Get workflow details
// @Tags N8N
// @Accept json
// @Produce json
// @Param id path string true "Workflow ID"
// @Success 200 {object} models.N8NWorkflow
// @Router /api/n8n/workflows/{id} [get]
func (h *N8NHandler) GetWorkflow(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid workflow ID")
	}
	
	workflow, err := h.service.GetWorkflow(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Workflow not found")
	}
	
	return c.JSON(http.StatusOK, workflow)
}

// ListWorkflows godoc
// @Summary List workflows
// @Description List all workflows for the company
// @Tags N8N
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param is_active query bool false "Filter by active status"
// @Success 200 {object} map[string]interface{}
// @Router /api/n8n/workflows [get]
func (h *N8NHandler) ListWorkflows(c echo.Context) error {
	userClaims := c.Get("user").(*middleware.Claims)
	
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
	if isActive := c.QueryParam("is_active"); isActive != "" {
		if active, err := strconv.ParseBool(isActive); err == nil {
			params["is_active"] = active
		}
	}
	
	workflows, total, err := h.service.ListWorkflows(userClaims.CompanyID, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": workflows,
		"pagination": map[string]interface{}{
			"total":     total,
			"page":      params["page"],
			"page_size": params["page_size"],
		},
	})
}

// ToggleWorkflow godoc
// @Summary Toggle workflow active status
// @Description Enable or disable a workflow
// @Tags N8N
// @Accept json
// @Produce json
// @Param id path string true "Workflow ID"
// @Param request body map[string]bool true "Active status"
// @Success 200 {object} models.N8NWorkflow
// @Router /api/n8n/workflows/{id}/toggle [patch]
func (h *N8NHandler) ToggleWorkflow(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid workflow ID")
	}
	
	var req struct {
		IsActive bool `json:"is_active"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	workflow, err := h.service.ToggleWorkflow(id, req.IsActive)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusOK, workflow)
}

// TriggerWorkflow godoc
// @Summary Trigger workflow
// @Description Manually trigger a workflow execution
// @Tags N8N
// @Accept json
// @Produce json
// @Param request body service.TriggerWorkflowRequest true "Trigger request"
// @Success 200 {object} models.N8NExecution
// @Router /api/n8n/trigger [post]
func (h *N8NHandler) TriggerWorkflow(c echo.Context) error {
	var req service.TriggerWorkflowRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	userClaims := c.Get("user").(*middleware.Claims)
	
	execution, err := h.service.TriggerWorkflow(userClaims.CompanyID, userClaims.UserID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusOK, execution)
}

// ListExecutions godoc
// @Summary List workflow executions
// @Description List workflow execution history
// @Tags N8N
// @Accept json
// @Produce json
// @Param workflow_id query string false "Filter by workflow ID"
// @Param status query string false "Filter by status"
// @Param from_date query string false "Filter from date"
// @Param to_date query string false "Filter to date"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} map[string]interface{}
// @Router /api/n8n/executions [get]
func (h *N8NHandler) ListExecutions(c echo.Context) error {
	userClaims := c.Get("user").(*middleware.Claims)
	
	params := make(map[string]interface{})
	if workflowID := c.QueryParam("workflow_id"); workflowID != "" {
		params["workflow_id"] = workflowID
	}
	if status := c.QueryParam("status"); status != "" {
		params["status"] = status
	}
	if fromDate := c.QueryParam("from_date"); fromDate != "" {
		params["from_date"] = fromDate
	}
	if toDate := c.QueryParam("to_date"); toDate != "" {
		params["to_date"] = toDate
	}
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
	
	executions, total, err := h.service.GetExecutions(userClaims.CompanyID, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": executions,
		"pagination": map[string]interface{}{
			"total":     total,
			"page":      params["page"],
			"page_size": params["page_size"],
		},
	})
}

// GetExecution godoc
// @Summary Get execution details
// @Description Get details of a specific execution
// @Tags N8N
// @Accept json
// @Produce json
// @Param id path string true "Execution ID"
// @Success 200 {object} models.N8NExecution
// @Router /api/n8n/executions/{id} [get]
func (h *N8NHandler) GetExecution(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid execution ID")
	}
	
	execution, err := h.service.GetExecution(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Execution not found")
	}
	
	return c.JSON(http.StatusOK, execution)
}

// CancelExecution godoc
// @Summary Cancel execution
// @Description Cancel a running execution
// @Tags N8N
// @Accept json
// @Produce json
// @Param id path string true "Execution ID"
// @Success 204
// @Router /api/n8n/executions/{id}/cancel [post]
func (h *N8NHandler) CancelExecution(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid execution ID")
	}
	
	if err := h.service.CancelExecution(id); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	return c.NoContent(http.StatusNoContent)
}

// ListWebhooks godoc
// @Summary List webhooks
// @Description List all webhooks for the company
// @Tags N8N
// @Accept json
// @Produce json
// @Success 200 {array} models.N8NWebhook
// @Router /api/n8n/webhooks [get]
func (h *N8NHandler) ListWebhooks(c echo.Context) error {
	userClaims := c.Get("user").(*middleware.Claims)
	
	webhooks, err := h.service.ListWebhooks(userClaims.CompanyID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, webhooks)
}

// RegisterWebhook godoc
// @Summary Register webhook
// @Description Register a new webhook
// @Tags N8N
// @Accept json
// @Produce json
// @Param request body service.RegisterWebhookRequest true "Register webhook request"
// @Success 201 {object} models.N8NWebhook
// @Router /api/n8n/webhooks [post]
func (h *N8NHandler) RegisterWebhook(c echo.Context) error {
	var req service.RegisterWebhookRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	userClaims := c.Get("user").(*middleware.Claims)
	
	webhook, err := h.service.RegisterWebhook(userClaims.CompanyID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusCreated, webhook)
}

// UpdateWebhook godoc
// @Summary Update webhook
// @Description Update webhook configuration
// @Tags N8N
// @Accept json
// @Produce json
// @Param id path string true "Webhook ID"
// @Param request body service.UpdateWebhookRequest true "Update webhook request"
// @Success 200 {object} models.N8NWebhook
// @Router /api/n8n/webhooks/{id} [put]
func (h *N8NHandler) UpdateWebhook(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid webhook ID")
	}
	
	var req service.UpdateWebhookRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	webhook, err := h.service.UpdateWebhook(id, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusOK, webhook)
}

// DeleteWebhook godoc
// @Summary Delete webhook
// @Description Delete a webhook
// @Tags N8N
// @Accept json
// @Produce json
// @Param id path string true "Webhook ID"
// @Success 204
// @Router /api/n8n/webhooks/{id} [delete]
func (h *N8NHandler) DeleteWebhook(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid webhook ID")
	}
	
	if err := h.service.DeleteWebhook(id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.NoContent(http.StatusNoContent)
}

// ListScheduledTasks godoc
// @Summary List scheduled tasks
// @Description List all scheduled tasks for the company
// @Tags N8N
// @Accept json
// @Produce json
// @Success 200 {array} models.N8NScheduledTask
// @Router /api/n8n/scheduled-tasks [get]
func (h *N8NHandler) ListScheduledTasks(c echo.Context) error {
	userClaims := c.Get("user").(*middleware.Claims)
	
	tasks, err := h.service.ListScheduledTasks(userClaims.CompanyID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, tasks)
}

// CreateScheduledTask godoc
// @Summary Create scheduled task
// @Description Create a new scheduled task
// @Tags N8N
// @Accept json
// @Produce json
// @Param request body service.CreateScheduledTaskRequest true "Create scheduled task request"
// @Success 201 {object} models.N8NScheduledTask
// @Router /api/n8n/scheduled-tasks [post]
func (h *N8NHandler) CreateScheduledTask(c echo.Context) error {
	var req service.CreateScheduledTaskRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	userClaims := c.Get("user").(*middleware.Claims)
	
	task, err := h.service.CreateScheduledTask(userClaims.CompanyID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusCreated, task)
}

// UpdateScheduledTask godoc
// @Summary Update scheduled task
// @Description Update scheduled task configuration
// @Tags N8N
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Param request body service.UpdateScheduledTaskRequest true "Update scheduled task request"
// @Success 200 {object} models.N8NScheduledTask
// @Router /api/n8n/scheduled-tasks/{id} [put]
func (h *N8NHandler) UpdateScheduledTask(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid task ID")
	}
	
	var req service.UpdateScheduledTaskRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	task, err := h.service.UpdateScheduledTask(id, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusOK, task)
}

// DeleteScheduledTask godoc
// @Summary Delete scheduled task
// @Description Delete a scheduled task
// @Tags N8N
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Success 204
// @Router /api/n8n/scheduled-tasks/{id} [delete]
func (h *N8NHandler) DeleteScheduledTask(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid task ID")
	}
	
	if err := h.service.DeleteScheduledTask(id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.NoContent(http.StatusNoContent)
}

// RegisterRoutes registers all N8N routes
func (h *N8NHandler) RegisterRoutes(e *echo.Echo, authMiddleware echo.MiddlewareFunc) {
	n8n := e.Group("/api/n8n", authMiddleware)
	
	// Connection
	n8n.GET("/test-connection", h.TestConnection)
	n8n.GET("/available-workflows", h.GetAvailableWorkflows)
	
	// Workflows
	n8n.GET("/workflows", h.ListWorkflows)
	n8n.POST("/workflows", h.CreateWorkflow)
	n8n.GET("/workflows/:id", h.GetWorkflow)
	n8n.PUT("/workflows/:id", h.UpdateWorkflow)
	n8n.DELETE("/workflows/:id", h.DeleteWorkflow)
	n8n.PATCH("/workflows/:id/toggle", h.ToggleWorkflow)
	
	// Execution
	n8n.POST("/trigger", h.TriggerWorkflow)
	n8n.GET("/executions", h.ListExecutions)
	n8n.GET("/executions/:id", h.GetExecution)
	n8n.POST("/executions/:id/cancel", h.CancelExecution)
	
	// Webhooks
	n8n.GET("/webhooks", h.ListWebhooks)
	n8n.POST("/webhooks", h.RegisterWebhook)
	n8n.PUT("/webhooks/:id", h.UpdateWebhook)
	n8n.DELETE("/webhooks/:id", h.DeleteWebhook)
	
	// Scheduled Tasks
	n8n.GET("/scheduled-tasks", h.ListScheduledTasks)
	n8n.POST("/scheduled-tasks", h.CreateScheduledTask)
	n8n.PUT("/scheduled-tasks/:id", h.UpdateScheduledTask)
	n8n.DELETE("/scheduled-tasks/:id", h.DeleteScheduledTask)
}