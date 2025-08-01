package repository

import (
	"github.com/fastenmind/fastener-api/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type N8NRepository interface {
	// Workflows
	CreateWorkflow(workflow *models.N8NWorkflow) error
	UpdateWorkflow(workflow *models.N8NWorkflow) error
	DeleteWorkflow(id uuid.UUID) error
	GetWorkflow(id uuid.UUID) (*models.N8NWorkflow, error)
	ListWorkflows(companyID uuid.UUID, params map[string]interface{}) ([]models.N8NWorkflow, int64, error)
	GetWorkflowByN8NID(workflowID string) (*models.N8NWorkflow, error)
	
	// Executions
	CreateExecution(execution *models.N8NExecution) error
	UpdateExecution(execution *models.N8NExecution) error
	GetExecution(id uuid.UUID) (*models.N8NExecution, error)
	GetExecutionByN8NID(executionID string) (*models.N8NExecution, error)
	ListExecutions(companyID uuid.UUID, params map[string]interface{}) ([]models.N8NExecution, int64, error)
	
	// Webhooks
	CreateWebhook(webhook *models.N8NWebhook) error
	UpdateWebhook(webhook *models.N8NWebhook) error
	DeleteWebhook(id uuid.UUID) error
	GetWebhook(id uuid.UUID) (*models.N8NWebhook, error)
	ListWebhooks(companyID uuid.UUID) ([]models.N8NWebhook, error)
	GetWebhooksByEventType(companyID uuid.UUID, eventType string) ([]models.N8NWebhook, error)
	
	// Scheduled Tasks
	CreateScheduledTask(task *models.N8NScheduledTask) error
	UpdateScheduledTask(task *models.N8NScheduledTask) error
	DeleteScheduledTask(id uuid.UUID) error
	GetScheduledTask(id uuid.UUID) (*models.N8NScheduledTask, error)
	ListScheduledTasks(companyID uuid.UUID) ([]models.N8NScheduledTask, error)
	GetDueScheduledTasks() ([]models.N8NScheduledTask, error)
	
	// Event Logs
	CreateEventLog(event *models.N8NEventLog) error
	GetUnprocessedEvents(companyID uuid.UUID) ([]models.N8NEventLog, error)
	MarkEventProcessed(id uuid.UUID, workflowIDs []string) error
}

type n8nRepository struct {
	db *gorm.DB
}

func NewN8NRepository(db interface{}) N8NRepository {
	gormDB, ok := db.(*gorm.DB)
	if !ok {
		panic("invalid database type, expected *gorm.DB")
	}
	return &n8nRepository{db: gormDB}
}

// Workflows
func (r *n8nRepository) CreateWorkflow(workflow *models.N8NWorkflow) error {
	return r.db.Create(workflow).Error
}

func (r *n8nRepository) UpdateWorkflow(workflow *models.N8NWorkflow) error {
	return r.db.Save(workflow).Error
}

func (r *n8nRepository) DeleteWorkflow(id uuid.UUID) error {
	return r.db.Delete(&models.N8NWorkflow{}, id).Error
}

func (r *n8nRepository) GetWorkflow(id uuid.UUID) (*models.N8NWorkflow, error) {
	var workflow models.N8NWorkflow
	if err := r.db.Preload("Creator").First(&workflow, id).Error; err != nil {
		return nil, err
	}
	return &workflow, nil
}

func (r *n8nRepository) ListWorkflows(companyID uuid.UUID, params map[string]interface{}) ([]models.N8NWorkflow, int64, error) {
	var workflows []models.N8NWorkflow
	var total int64
	
	query := r.db.Model(&models.N8NWorkflow{}).Where("company_id = ?", companyID)
	
	if isActive, ok := params["is_active"].(bool); ok {
		query = query.Where("is_active = ?", isActive)
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
	
	if err := query.Preload("Creator").Order("created_at DESC").Find(&workflows).Error; err != nil {
		return nil, 0, err
	}
	
	return workflows, total, nil
}

func (r *n8nRepository) GetWorkflowByN8NID(workflowID string) (*models.N8NWorkflow, error) {
	var workflow models.N8NWorkflow
	if err := r.db.Where("workflow_id = ?", workflowID).First(&workflow).Error; err != nil {
		return nil, err
	}
	return &workflow, nil
}

// Executions
func (r *n8nRepository) CreateExecution(execution *models.N8NExecution) error {
	return r.db.Create(execution).Error
}

func (r *n8nRepository) UpdateExecution(execution *models.N8NExecution) error {
	return r.db.Save(execution).Error
}

func (r *n8nRepository) GetExecution(id uuid.UUID) (*models.N8NExecution, error) {
	var execution models.N8NExecution
	if err := r.db.Preload("Trigger").First(&execution, id).Error; err != nil {
		return nil, err
	}
	return &execution, nil
}

func (r *n8nRepository) GetExecutionByN8NID(executionID string) (*models.N8NExecution, error) {
	var execution models.N8NExecution
	if err := r.db.Where("execution_id = ?", executionID).First(&execution).Error; err != nil {
		return nil, err
	}
	return &execution, nil
}

func (r *n8nRepository) ListExecutions(companyID uuid.UUID, params map[string]interface{}) ([]models.N8NExecution, int64, error) {
	var executions []models.N8NExecution
	var total int64
	
	query := r.db.Model(&models.N8NExecution{}).Where("company_id = ?", companyID)
	
	if workflowID, ok := params["workflow_id"].(string); ok && workflowID != "" {
		query = query.Where("workflow_id = ?", workflowID)
	}
	
	if status, ok := params["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	
	if fromDate, ok := params["from_date"].(string); ok && fromDate != "" {
		query = query.Where("started_at >= ?", fromDate)
	}
	
	if toDate, ok := params["to_date"].(string); ok && toDate != "" {
		query = query.Where("started_at <= ?", toDate)
	}
	
	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Apply pagination
	if page, ok := params["page"].(int); ok && page > 0 {
		pageSize := 50
		if ps, ok := params["page_size"].(int); ok && ps > 0 {
			pageSize = ps
		}
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}
	
	if err := query.Preload("Trigger").Order("started_at DESC").Find(&executions).Error; err != nil {
		return nil, 0, err
	}
	
	return executions, total, nil
}

// Webhooks
func (r *n8nRepository) CreateWebhook(webhook *models.N8NWebhook) error {
	return r.db.Create(webhook).Error
}

func (r *n8nRepository) UpdateWebhook(webhook *models.N8NWebhook) error {
	return r.db.Save(webhook).Error
}

func (r *n8nRepository) DeleteWebhook(id uuid.UUID) error {
	return r.db.Delete(&models.N8NWebhook{}, id).Error
}

func (r *n8nRepository) GetWebhook(id uuid.UUID) (*models.N8NWebhook, error) {
	var webhook models.N8NWebhook
	if err := r.db.First(&webhook, id).Error; err != nil {
		return nil, err
	}
	return &webhook, nil
}

func (r *n8nRepository) ListWebhooks(companyID uuid.UUID) ([]models.N8NWebhook, error) {
	var webhooks []models.N8NWebhook
	if err := r.db.Where("company_id = ? AND is_active = ?", companyID, true).Find(&webhooks).Error; err != nil {
		return nil, err
	}
	return webhooks, nil
}

func (r *n8nRepository) GetWebhooksByEventType(companyID uuid.UUID, eventType string) ([]models.N8NWebhook, error) {
	var webhooks []models.N8NWebhook
	if err := r.db.Where("company_id = ? AND is_active = ? AND ? = ANY(event_types)", 
		companyID, true, eventType).Find(&webhooks).Error; err != nil {
		return nil, err
	}
	return webhooks, nil
}

// Scheduled Tasks
func (r *n8nRepository) CreateScheduledTask(task *models.N8NScheduledTask) error {
	return r.db.Create(task).Error
}

func (r *n8nRepository) UpdateScheduledTask(task *models.N8NScheduledTask) error {
	return r.db.Save(task).Error
}

func (r *n8nRepository) DeleteScheduledTask(id uuid.UUID) error {
	return r.db.Delete(&models.N8NScheduledTask{}, id).Error
}

func (r *n8nRepository) GetScheduledTask(id uuid.UUID) (*models.N8NScheduledTask, error) {
	var task models.N8NScheduledTask
	if err := r.db.First(&task, id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *n8nRepository) ListScheduledTasks(companyID uuid.UUID) ([]models.N8NScheduledTask, error) {
	var tasks []models.N8NScheduledTask
	if err := r.db.Where("company_id = ?", companyID).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *n8nRepository) GetDueScheduledTasks() ([]models.N8NScheduledTask, error) {
	var tasks []models.N8NScheduledTask
	if err := r.db.Where("is_active = ? AND (next_run_at IS NULL OR next_run_at <= NOW())", 
		true).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

// Event Logs
func (r *n8nRepository) CreateEventLog(event *models.N8NEventLog) error {
	return r.db.Create(event).Error
}

func (r *n8nRepository) GetUnprocessedEvents(companyID uuid.UUID) ([]models.N8NEventLog, error) {
	var events []models.N8NEventLog
	if err := r.db.Where("company_id = ? AND processed_at IS NULL", 
		companyID).Order("created_at ASC").Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

func (r *n8nRepository) MarkEventProcessed(id uuid.UUID, workflowIDs []string) error {
	now := gorm.NowFunc()
	return r.db.Model(&models.N8NEventLog{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"processed_at": now,
			"workflow_ids": workflowIDs,
		}).Error
}