package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/fastenmind/fastener-api/internal/service"
)

type AdvancedHandler struct {
	advancedService *service.AdvancedService
}

func NewAdvancedHandler(advancedService *service.AdvancedService) *AdvancedHandler {
	return &AdvancedHandler{
		advancedService: advancedService,
	}
}

// AI Assistant Endpoints
func (h *AdvancedHandler) CreateAIAssistant(c echo.Context) error {
	userID := getUserIDFromContext(c)
	
	var req service.CreateAIAssistantRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	assistant, err := h.advancedService.CreateAIAssistant(userID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, assistant)
}

func (h *AdvancedHandler) GetAIAssistant(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid assistant ID")
	}

	assistant, err := h.advancedService.GetAIAssistant(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Assistant not found")
	}

	return c.JSON(http.StatusOK, assistant)
}

func (h *AdvancedHandler) ListAIAssistants(c echo.Context) error {
	companyID := getCompanyIDFromContext(c)
	
	assistantType := c.QueryParam("type")
	isActiveStr := c.QueryParam("is_active")
	
	var isActive *bool
	if isActiveStr != "" {
		activeVal := isActiveStr == "true"
		isActive = &activeVal
	}

	assistants, err := h.advancedService.GetAIAssistantsByCompany(companyID, assistantType, isActive)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": assistants,
	})
}

func (h *AdvancedHandler) UpdateAIAssistant(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid assistant ID")
	}

	userID := getUserIDFromContext(c)
	
	var req service.UpdateAIAssistantRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	assistant, err := h.advancedService.UpdateAIAssistant(id, userID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, assistant)
}

func (h *AdvancedHandler) DeleteAIAssistant(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid assistant ID")
	}

	err = h.advancedService.DeleteAIAssistant(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

// AI Conversation Endpoints
func (h *AdvancedHandler) StartConversation(c echo.Context) error {
	userID := getUserIDFromContext(c)
	
	var req struct {
		AssistantID uuid.UUID              `json:"assistant_id" validate:"required"`
		Title       string                 `json:"title"`
		Context     map[string]interface{} `json:"context"`
	}
	
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	session, err := h.advancedService.StartConversation(userID, req.AssistantID, req.Title, req.Context)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, session)
}

func (h *AdvancedHandler) SendMessage(c echo.Context) error {
	sessionID, err := uuid.Parse(c.Param("session_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid session ID")
	}

	userID := getUserIDFromContext(c)
	
	var req struct {
		Content string `json:"content" validate:"required"`
	}
	
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	response, err := h.advancedService.SendMessage(sessionID, userID, req.Content)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, response)
}

func (h *AdvancedHandler) GetConversationHistory(c echo.Context) error {
	sessionID, err := uuid.Parse(c.Param("session_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid session ID")
	}

	userID := getUserIDFromContext(c)
	
	limitStr := c.QueryParam("limit")
	limit := 50 // default
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	messages, err := h.advancedService.GetConversationHistory(sessionID, userID, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": messages,
	})
}

func (h *AdvancedHandler) EndConversation(c echo.Context) error {
	sessionID, err := uuid.Parse(c.Param("session_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid session ID")
	}

	userID := getUserIDFromContext(c)

	err = h.advancedService.EndConversation(sessionID, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

// Smart Recommendation Endpoints
func (h *AdvancedHandler) CreateRecommendation(c echo.Context) error {
	userID := getUserIDFromContext(c)
	
	var req service.CreateRecommendationRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	recommendation, err := h.advancedService.CreateRecommendation(userID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, recommendation)
}

func (h *AdvancedHandler) ListRecommendations(c echo.Context) error {
	userID := getUserIDFromContext(c)
	
	recommendationType := c.QueryParam("type")
	status := c.QueryParam("status")
	limitStr := c.QueryParam("limit")
	
	limit := 20 // default
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	recommendations, err := h.advancedService.GetRecommendationsByUser(userID, recommendationType, status, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": recommendations,
	})
}

func (h *AdvancedHandler) UpdateRecommendationStatus(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid recommendation ID")
	}

	userID := getUserIDFromContext(c)
	
	var req struct {
		Status string `json:"status" validate:"required"`
	}
	
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	recommendation, err := h.advancedService.UpdateRecommendationStatus(id, userID, req.Status)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, recommendation)
}

// Advanced Search Endpoints
func (h *AdvancedHandler) CreateAdvancedSearch(c echo.Context) error {
	userID := getUserIDFromContext(c)
	
	var req service.CreateAdvancedSearchRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	search, err := h.advancedService.CreateAdvancedSearch(userID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, search)
}

func (h *AdvancedHandler) ListAdvancedSearches(c echo.Context) error {
	userID := getUserIDFromContext(c)
	
	searchType := c.QueryParam("search_type")
	isPublicStr := c.QueryParam("is_public")
	
	var isPublic *bool
	if isPublicStr != "" {
		publicVal := isPublicStr == "true"
		isPublic = &publicVal
	}

	searches, err := h.advancedService.GetAdvancedSearchesByUser(userID, searchType, isPublic)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": searches,
	})
}

func (h *AdvancedHandler) ExecuteAdvancedSearch(c echo.Context) error {
	searchID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid search ID")
	}

	userID := getUserIDFromContext(c)

	result, err := h.advancedService.ExecuteAdvancedSearch(searchID, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

// Batch Operation Endpoints
func (h *AdvancedHandler) CreateBatchOperation(c echo.Context) error {
	userID := getUserIDFromContext(c)
	
	var req service.CreateBatchOperationRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	operation, err := h.advancedService.CreateBatchOperation(userID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, operation)
}

func (h *AdvancedHandler) ListBatchOperations(c echo.Context) error {
	userID := getUserIDFromContext(c)
	
	status := c.QueryParam("status")
	limitStr := c.QueryParam("limit")
	
	limit := 20 // default
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	operations, err := h.advancedService.GetBatchOperationsByUser(userID, status, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": operations,
	})
}

func (h *AdvancedHandler) GetBatchOperation(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid operation ID")
	}

	userID := getUserIDFromContext(c)

	operation, err := h.advancedService.GetBatchOperation(id, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, operation)
}

// Custom Field Endpoints
func (h *AdvancedHandler) CreateCustomField(c echo.Context) error {
	userID := getUserIDFromContext(c)
	
	var req service.CreateCustomFieldRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	field, err := h.advancedService.CreateCustomField(userID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, field)
}

func (h *AdvancedHandler) ListCustomFields(c echo.Context) error {
	companyID := getCompanyIDFromContext(c)
	tableName := c.QueryParam("table_name")
	
	if tableName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "table_name parameter is required")
	}
	
	isActiveStr := c.QueryParam("is_active")
	var isActive *bool
	if isActiveStr != "" {
		activeVal := isActiveStr == "true"
		isActive = &activeVal
	}

	fields, err := h.advancedService.GetCustomFieldsByTable(companyID, tableName, isActive)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": fields,
	})
}

func (h *AdvancedHandler) SetCustomFieldValue(c echo.Context) error {
	userID := getUserIDFromContext(c)
	
	var req service.SetCustomFieldValueRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err := h.advancedService.SetCustomFieldValue(userID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *AdvancedHandler) GetCustomFieldValues(c echo.Context) error {
	resourceID, err := uuid.Parse(c.Param("resource_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid resource ID")
	}

	resourceType := c.QueryParam("resource_type")
	if resourceType == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "resource_type parameter is required")
	}

	values, err := h.advancedService.GetCustomFieldValues(resourceID, resourceType)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": values,
	})
}

// Security Event Endpoints
func (h *AdvancedHandler) CreateSecurityEvent(c echo.Context) error {
	var req service.CreateSecurityEventRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	event, err := h.advancedService.CreateSecurityEvent(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, event)
}

func (h *AdvancedHandler) ListSecurityEvents(c echo.Context) error {
	companyID := getCompanyIDFromContext(c)
	
	eventType := c.QueryParam("event_type")
	severity := c.QueryParam("severity")
	status := c.QueryParam("status")
	limitStr := c.QueryParam("limit")
	
	limit := 50 // default
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	events, err := h.advancedService.GetSecurityEventsByCompany(companyID, eventType, severity, status, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": events,
	})
}

// Performance Metric Endpoints
func (h *AdvancedHandler) RecordPerformanceMetric(c echo.Context) error {
	var req service.RecordPerformanceMetricRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err := h.advancedService.RecordPerformanceMetric(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusCreated)
}

func (h *AdvancedHandler) GetPerformanceStats(c echo.Context) error {
	companyID := getCompanyIDFromContext(c)
	
	metricType := c.QueryParam("metric_type")
	if metricType == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "metric_type parameter is required")
	}
	
	startTimeStr := c.QueryParam("start_time")
	endTimeStr := c.QueryParam("end_time")
	
	var startTime, endTime time.Time
	var err error
	
	if startTimeStr != "" {
		startTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid start_time format")
		}
	} else {
		startTime = time.Now().AddDate(0, 0, -7) // Last 7 days
	}
	
	if endTimeStr != "" {
		endTime, err = time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid end_time format")
		}
	} else {
		endTime = time.Now()
	}

	stats, err := h.advancedService.GetPerformanceStats(companyID, metricType, startTime, endTime)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, stats)
}

// Backup Endpoints
func (h *AdvancedHandler) CreateBackup(c echo.Context) error {
	userID := getUserIDFromContext(c)
	
	var req service.CreateBackupRecordRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	record, err := h.advancedService.CreateBackupRecord(userID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, record)
}

func (h *AdvancedHandler) ListBackups(c echo.Context) error {
	companyID := getCompanyIDFromContext(c)
	
	backupType := c.QueryParam("backup_type")
	status := c.QueryParam("status")
	limitStr := c.QueryParam("limit")
	
	limit := 20 // default
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	records, err := h.advancedService.GetBackupRecordsByCompany(companyID, backupType, status, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": records,
	})
}

// Multi-language Support Endpoints
func (h *AdvancedHandler) ListSystemLanguages(c echo.Context) error {
	companyID := getCompanyIDFromContext(c)
	
	isActiveStr := c.QueryParam("is_active")
	var isActive *bool
	if isActiveStr != "" {
		activeVal := isActiveStr == "true"
		isActive = &activeVal
	}

	languages, err := h.advancedService.GetSystemLanguages(&companyID, isActive)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": languages,
	})
}

func (h *AdvancedHandler) GetTranslations(c echo.Context) error {
	languageCode := c.Param("language_code")
	if languageCode == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Language code is required")
	}
	
	onlyApprovedStr := c.QueryParam("only_approved")
	onlyApproved := onlyApprovedStr == "true"

	translations, err := h.advancedService.GetTranslationsByLanguage(languageCode, onlyApproved)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, translations)
}

// Helper functions
func getUserIDFromContext(c echo.Context) uuid.UUID {
	userID, _ := c.Get("user_id").(uuid.UUID)
	return userID
}

func getCompanyIDFromContext(c echo.Context) uuid.UUID {
	companyID, _ := c.Get("company_id").(uuid.UUID)
	return companyID
}