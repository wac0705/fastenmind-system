package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/fastenmind/fastener-api/internal/models"
	"github.com/fastenmind/fastener-api/internal/repository"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
)

type N8NService interface {
	// Connection
	TestConnection() (bool, string, error)
	GetAvailableWorkflows() ([]N8NWorkflowInfo, error)
	
	// Workflow Management
	CreateWorkflow(companyID, userID uuid.UUID, req CreateWorkflowRequest) (*models.N8NWorkflow, error)
	UpdateWorkflow(id uuid.UUID, req UpdateWorkflowRequest) (*models.N8NWorkflow, error)
	DeleteWorkflow(id uuid.UUID) error
	GetWorkflow(id uuid.UUID) (*models.N8NWorkflow, error)
	ListWorkflows(companyID uuid.UUID, params map[string]interface{}) ([]models.N8NWorkflow, int64, error)
	ToggleWorkflow(id uuid.UUID, active bool) (*models.N8NWorkflow, error)
	
	// Execution
	TriggerWorkflow(companyID, userID uuid.UUID, req TriggerWorkflowRequest) (*models.N8NExecution, error)
	GetExecutions(companyID uuid.UUID, params map[string]interface{}) ([]models.N8NExecution, int64, error)
	GetExecution(id uuid.UUID) (*models.N8NExecution, error)
	CancelExecution(id uuid.UUID) error
	
	// Webhooks
	RegisterWebhook(companyID uuid.UUID, req RegisterWebhookRequest) (*models.N8NWebhook, error)
	UpdateWebhook(id uuid.UUID, req UpdateWebhookRequest) (*models.N8NWebhook, error)
	DeleteWebhook(id uuid.UUID) error
	ListWebhooks(companyID uuid.UUID) ([]models.N8NWebhook, error)
	
	// Scheduled Tasks
	CreateScheduledTask(companyID uuid.UUID, req CreateScheduledTaskRequest) (*models.N8NScheduledTask, error)
	UpdateScheduledTask(id uuid.UUID, req UpdateScheduledTaskRequest) (*models.N8NScheduledTask, error)
	DeleteScheduledTask(id uuid.UUID) error
	ListScheduledTasks(companyID uuid.UUID) ([]models.N8NScheduledTask, error)
	
	// Event Processing
	LogEvent(companyID, userID uuid.UUID, eventType, entityType string, entityID uuid.UUID, data map[string]interface{}) error
	ProcessPendingEvents(companyID uuid.UUID) error
}

type CreateWorkflowRequest struct {
	Name          string                 `json:"name" validate:"required"`
	Description   string                 `json:"description"`
	WorkflowID    string                 `json:"workflow_id" validate:"required"`
	TriggerType   string                 `json:"trigger_type" validate:"required,oneof=webhook schedule manual event"`
	TriggerConfig map[string]interface{} `json:"trigger_config"`
}

type UpdateWorkflowRequest struct {
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	TriggerConfig map[string]interface{} `json:"trigger_config"`
	IsActive      *bool                  `json:"is_active"`
}

type TriggerWorkflowRequest struct {
	WorkflowID         string                 `json:"workflow_id" validate:"required"`
	Data               map[string]interface{} `json:"data"`
	WaitForCompletion  bool                   `json:"wait_for_completion"`
}

type RegisterWebhookRequest struct {
	Name        string            `json:"name" validate:"required"`
	WorkflowID  string            `json:"workflow_id" validate:"required"`
	EventTypes  []string          `json:"event_types" validate:"required"`
	TargetURL   string            `json:"target_url"`
	Headers     map[string]string `json:"headers"`
}

type UpdateWebhookRequest struct {
	Name       string            `json:"name"`
	EventTypes []string          `json:"event_types"`
	TargetURL  string            `json:"target_url"`
	Headers    map[string]string `json:"headers"`
	IsActive   *bool             `json:"is_active"`
}

type CreateScheduledTaskRequest struct {
	Name           string                 `json:"name" validate:"required"`
	WorkflowID     string                 `json:"workflow_id" validate:"required"`
	CronExpression string                 `json:"cron_expression" validate:"required"`
	Timezone       string                 `json:"timezone"`
	Data           map[string]interface{} `json:"data"`
}

type UpdateScheduledTaskRequest struct {
	Name           string                 `json:"name"`
	CronExpression string                 `json:"cron_expression"`
	Timezone       string                 `json:"timezone"`
	Data           map[string]interface{} `json:"data"`
	IsActive       *bool                  `json:"is_active"`
}

type N8NWorkflowInfo struct {
	ID   string   `json:"id"`
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

type n8nService struct {
	repo         repository.N8NRepository
	n8nURL       string
	n8nAPIKey    string
	httpClient   *http.Client
	cronParser   cron.Parser
}

func NewN8NService(repo repository.N8NRepository) N8NService {
	n8nURL := os.Getenv("N8N_URL")
	if n8nURL == "" {
		n8nURL = "http://localhost:5678"
	}
	
	return &n8nService{
		repo:       repo,
		n8nURL:     n8nURL,
		n8nAPIKey:  os.Getenv("N8N_API_KEY"),
		httpClient: &http.Client{Timeout: 30 * time.Second},
		cronParser: cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow),
	}
}

// Connection
func (s *n8nService) TestConnection() (bool, string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/workflows", s.n8nURL), nil)
	if err != nil {
		return false, "", err
	}
	
	if s.n8nAPIKey != "" {
		req.Header.Set("X-N8N-API-KEY", s.n8nAPIKey)
	}
	
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return false, "", fmt.Errorf("failed to connect to N8N: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return false, "", fmt.Errorf("N8N returned status %d", resp.StatusCode)
	}
	
	// Try to get version from headers or response
	version := resp.Header.Get("X-N8N-Version")
	if version == "" {
		version = "unknown"
	}
	
	return true, version, nil
}

func (s *n8nService) GetAvailableWorkflows() ([]N8NWorkflowInfo, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/workflows", s.n8nURL), nil)
	if err != nil {
		return nil, err
	}
	
	if s.n8nAPIKey != "" {
		req.Header.Set("X-N8N-API-KEY", s.n8nAPIKey)
	}
	
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("N8N returned status %d", resp.StatusCode)
	}
	
	var result struct {
		Data []struct {
			ID     string            `json:"id"`
			Name   string            `json:"name"`
			Tags   []string          `json:"tags"`
			Active bool              `json:"active"`
			Nodes  []json.RawMessage `json:"nodes"`
		} `json:"data"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	
	workflows := make([]N8NWorkflowInfo, 0, len(result.Data))
	for _, w := range result.Data {
		if w.Active {
			workflows = append(workflows, N8NWorkflowInfo{
				ID:   w.ID,
				Name: w.Name,
				Tags: w.Tags,
			})
		}
	}
	
	return workflows, nil
}

// Workflow Management
func (s *n8nService) CreateWorkflow(companyID, userID uuid.UUID, req CreateWorkflowRequest) (*models.N8NWorkflow, error) {
	// Validate cron expression if schedule trigger
	if req.TriggerType == "schedule" {
		if cronExpr, ok := req.TriggerConfig["cron"].(string); ok {
			if _, err := s.cronParser.Parse(cronExpr); err != nil {
				return nil, fmt.Errorf("invalid cron expression: %w", err)
			}
		} else {
			return nil, errors.New("cron expression required for schedule trigger")
		}
	}
	
	workflow := &models.N8NWorkflow{
		CompanyID:     companyID,
		Name:          req.Name,
		Description:   req.Description,
		WorkflowID:    req.WorkflowID,
		TriggerType:   req.TriggerType,
		TriggerConfig: models.JSONB(req.TriggerConfig),
		IsActive:      true,
		CreatedBy:     userID,
	}
	
	if err := s.repo.CreateWorkflow(workflow); err != nil {
		return nil, err
	}
	
	// If event-based trigger, register webhook
	if req.TriggerType == "event" && req.TriggerConfig["event"] != nil {
		eventType := req.TriggerConfig["event"].(string)
		webhook := &models.N8NWebhook{
			CompanyID:  companyID,
			Name:       fmt.Sprintf("%s Webhook", req.Name),
			WorkflowID: req.WorkflowID,
			EventTypes: []string{eventType},
			IsActive:   true,
		}
		s.repo.CreateWebhook(webhook)
	}
	
	return workflow, nil
}

func (s *n8nService) UpdateWorkflow(id uuid.UUID, req UpdateWorkflowRequest) (*models.N8NWorkflow, error) {
	workflow, err := s.repo.GetWorkflow(id)
	if err != nil {
		return nil, err
	}
	
	if req.Name != "" {
		workflow.Name = req.Name
	}
	if req.Description != "" {
		workflow.Description = req.Description
	}
	if req.TriggerConfig != nil {
		workflow.TriggerConfig = models.JSONB(req.TriggerConfig)
	}
	if req.IsActive != nil {
		workflow.IsActive = *req.IsActive
	}
	
	if err := s.repo.UpdateWorkflow(workflow); err != nil {
		return nil, err
	}
	
	return workflow, nil
}

func (s *n8nService) DeleteWorkflow(id uuid.UUID) error {
	return s.repo.DeleteWorkflow(id)
}

func (s *n8nService) GetWorkflow(id uuid.UUID) (*models.N8NWorkflow, error) {
	return s.repo.GetWorkflow(id)
}

func (s *n8nService) ListWorkflows(companyID uuid.UUID, params map[string]interface{}) ([]models.N8NWorkflow, int64, error) {
	return s.repo.ListWorkflows(companyID, params)
}

func (s *n8nService) ToggleWorkflow(id uuid.UUID, active bool) (*models.N8NWorkflow, error) {
	workflow, err := s.repo.GetWorkflow(id)
	if err != nil {
		return nil, err
	}
	
	workflow.IsActive = active
	if err := s.repo.UpdateWorkflow(workflow); err != nil {
		return nil, err
	}
	
	return workflow, nil
}

// Execution
func (s *n8nService) TriggerWorkflow(companyID, userID uuid.UUID, req TriggerWorkflowRequest) (*models.N8NExecution, error) {
	// Create execution record
	execution := &models.N8NExecution{
		CompanyID:   companyID,
		WorkflowID:  req.WorkflowID,
		ExecutionID: uuid.New().String(), // Temporary ID
		Status:      "running",
		StartedAt:   time.Now(),
		InputData:   models.JSONB(req.Data),
		TriggeredBy: userID,
	}
	
	if err := s.repo.CreateExecution(execution); err != nil {
		return nil, err
	}
	
	// Trigger N8N workflow
	go s.executeN8NWorkflow(execution, req)
	
	return execution, nil
}

func (s *n8nService) executeN8NWorkflow(execution *models.N8NExecution, req TriggerWorkflowRequest) {
	// Build N8N webhook URL
	webhookURL := fmt.Sprintf("%s/webhook/%s", s.n8nURL, req.WorkflowID)
	
	// Prepare request body
	body, _ := json.Marshal(req.Data)
	
	// Create HTTP request
	httpReq, err := http.NewRequest("POST", webhookURL, bytes.NewReader(body))
	if err != nil {
		s.updateExecutionStatus(execution, "error", err.Error(), nil)
		return
	}
	
	httpReq.Header.Set("Content-Type", "application/json")
	
	// Execute request
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		s.updateExecutionStatus(execution, "error", err.Error(), nil)
		return
	}
	defer resp.Body.Close()
	
	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		s.updateExecutionStatus(execution, "error", err.Error(), nil)
		return
	}
	
	// Parse response
	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err == nil {
		if executionID, ok := result["executionId"].(string); ok {
			execution.ExecutionID = executionID
		}
	}
	
	// Update execution status
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		s.updateExecutionStatus(execution, "success", "", result)
	} else {
		s.updateExecutionStatus(execution, "error", fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(respBody)), nil)
	}
}

func (s *n8nService) updateExecutionStatus(execution *models.N8NExecution, status, errorMsg string, output map[string]interface{}) {
	now := time.Now()
	execution.Status = status
	execution.FinishedAt = &now
	execution.ErrorMessage = errorMsg
	if output != nil {
		execution.OutputData = models.JSONB(output)
	}
	
	s.repo.UpdateExecution(execution)
}

func (s *n8nService) GetExecutions(companyID uuid.UUID, params map[string]interface{}) ([]models.N8NExecution, int64, error) {
	return s.repo.ListExecutions(companyID, params)
}

func (s *n8nService) GetExecution(id uuid.UUID) (*models.N8NExecution, error) {
	return s.repo.GetExecution(id)
}

func (s *n8nService) CancelExecution(id uuid.UUID) error {
	execution, err := s.repo.GetExecution(id)
	if err != nil {
		return err
	}
	
	if execution.Status != "running" {
		return errors.New("execution is not running")
	}
	
	// In a real implementation, you would call N8N API to cancel the execution
	// For now, just update the status
	execution.Status = "canceled"
	now := time.Now()
	execution.FinishedAt = &now
	
	return s.repo.UpdateExecution(execution)
}

// Webhooks
func (s *n8nService) RegisterWebhook(companyID uuid.UUID, req RegisterWebhookRequest) (*models.N8NWebhook, error) {
	webhook := &models.N8NWebhook{
		CompanyID:  companyID,
		Name:       req.Name,
		WorkflowID: req.WorkflowID,
		EventTypes: req.EventTypes,
		TargetURL:  req.TargetURL,
		Headers:    models.JSONB(req.Headers),
		IsActive:   true,
	}
	
	if err := s.repo.CreateWebhook(webhook); err != nil {
		return nil, err
	}
	
	return webhook, nil
}

func (s *n8nService) UpdateWebhook(id uuid.UUID, req UpdateWebhookRequest) (*models.N8NWebhook, error) {
	webhook, err := s.repo.GetWebhook(id)
	if err != nil {
		return nil, err
	}
	
	if req.Name != "" {
		webhook.Name = req.Name
	}
	if req.EventTypes != nil {
		webhook.EventTypes = req.EventTypes
	}
	if req.TargetURL != "" {
		webhook.TargetURL = req.TargetURL
	}
	if req.Headers != nil {
		webhook.Headers = models.JSONB(req.Headers)
	}
	if req.IsActive != nil {
		webhook.IsActive = *req.IsActive
	}
	
	if err := s.repo.UpdateWebhook(webhook); err != nil {
		return nil, err
	}
	
	return webhook, nil
}

func (s *n8nService) DeleteWebhook(id uuid.UUID) error {
	return s.repo.DeleteWebhook(id)
}

func (s *n8nService) ListWebhooks(companyID uuid.UUID) ([]models.N8NWebhook, error) {
	return s.repo.ListWebhooks(companyID)
}

// Scheduled Tasks
func (s *n8nService) CreateScheduledTask(companyID uuid.UUID, req CreateScheduledTaskRequest) (*models.N8NScheduledTask, error) {
	// Validate cron expression
	schedule, err := s.cronParser.Parse(req.CronExpression)
	if err != nil {
		return nil, fmt.Errorf("invalid cron expression: %w", err)
	}
	
	timezone := req.Timezone
	if timezone == "" {
		timezone = "UTC"
	}
	
	nextRun := schedule.Next(time.Now())
	task := &models.N8NScheduledTask{
		CompanyID:      companyID,
		Name:           req.Name,
		WorkflowID:     req.WorkflowID,
		CronExpression: req.CronExpression,
		Timezone:       timezone,
		Data:           models.JSONB(req.Data),
		IsActive:       true,
		NextRunAt:      &nextRun,
	}
	
	if err := s.repo.CreateScheduledTask(task); err != nil {
		return nil, err
	}
	
	return task, nil
}

func (s *n8nService) UpdateScheduledTask(id uuid.UUID, req UpdateScheduledTaskRequest) (*models.N8NScheduledTask, error) {
	task, err := s.repo.GetScheduledTask(id)
	if err != nil {
		return nil, err
	}
	
	if req.Name != "" {
		task.Name = req.Name
	}
	if req.CronExpression != "" {
		// Validate new cron expression
		schedule, err := s.cronParser.Parse(req.CronExpression)
		if err != nil {
			return nil, fmt.Errorf("invalid cron expression: %w", err)
		}
		task.CronExpression = req.CronExpression
		nextRun := schedule.Next(time.Now())
		task.NextRunAt = &nextRun
	}
	if req.Timezone != "" {
		task.Timezone = req.Timezone
	}
	if req.Data != nil {
		task.Data = models.JSONB(req.Data)
	}
	if req.IsActive != nil {
		task.IsActive = *req.IsActive
	}
	
	if err := s.repo.UpdateScheduledTask(task); err != nil {
		return nil, err
	}
	
	return task, nil
}

func (s *n8nService) DeleteScheduledTask(id uuid.UUID) error {
	return s.repo.DeleteScheduledTask(id)
}

func (s *n8nService) ListScheduledTasks(companyID uuid.UUID) ([]models.N8NScheduledTask, error) {
	return s.repo.ListScheduledTasks(companyID)
}

// Event Processing
func (s *n8nService) LogEvent(companyID, userID uuid.UUID, eventType, entityType string, entityID uuid.UUID, data map[string]interface{}) error {
	event := &models.N8NEventLog{
		CompanyID:   companyID,
		EventType:   eventType,
		EntityType:  entityType,
		EntityID:    entityID,
		EventData:   models.JSONB(data),
		TriggeredBy: userID,
	}
	
	if err := s.repo.CreateEventLog(event); err != nil {
		return err
	}
	
	// Process event asynchronously
	go s.processEvent(companyID, event)
	
	return nil
}

func (s *n8nService) processEvent(companyID uuid.UUID, event *models.N8NEventLog) {
	// Find webhooks for this event type
	webhooks, err := s.repo.GetWebhooksByEventType(companyID, event.EventType)
	if err != nil {
		return
	}
	
	triggeredWorkflows := make([]string, 0)
	
	for _, webhook := range webhooks {
		// Trigger workflow
		_, err := s.TriggerWorkflow(companyID, event.TriggeredBy, TriggerWorkflowRequest{
			WorkflowID: webhook.WorkflowID,
			Data: map[string]interface{}{
				"event":      event.EventType,
				"entity":     event.EntityType,
				"entity_id":  event.EntityID,
				"event_data": event.EventData,
				"timestamp":  event.CreatedAt,
			},
		})
		
		if err == nil {
			triggeredWorkflows = append(triggeredWorkflows, webhook.WorkflowID)
		}
	}
	
	// Mark event as processed
	s.repo.MarkEventProcessed(event.ID, triggeredWorkflows)
}

func (s *n8nService) ProcessPendingEvents(companyID uuid.UUID) error {
	events, err := s.repo.GetUnprocessedEvents(companyID)
	if err != nil {
		return err
	}
	
	for _, event := range events {
		s.processEvent(companyID, &event)
	}
	
	return nil
}