package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ExecutionTracker tracks agent executions
type ExecutionTracker interface {
	// StartExecution starts tracking a new execution
	StartExecution(ctx context.Context, agentType AgentType, agentName string, input AgentInput, userID uuid.UUID) (*AgentExecution, error)
	
	// UpdateExecution updates an execution status
	UpdateExecution(ctx context.Context, executionID uuid.UUID, status AgentStatus, output *AgentOutput, err error) error
	
	// GetExecution retrieves an execution by ID
	GetExecution(ctx context.Context, executionID uuid.UUID) (*AgentExecution, error)
	
	// ListExecutions lists executions with filters
	ListExecutions(ctx context.Context, filter ExecutionFilter) ([]*AgentExecution, error)
	
	// GetExecutionMetrics gets execution metrics
	GetExecutionMetrics(ctx context.Context, filter MetricsFilter) (*ExecutionMetrics, error)
	
	// StartChainExecution starts tracking a chain execution
	StartChainExecution(ctx context.Context, chainID uuid.UUID, userID uuid.UUID) (*AgentChainExecution, error)
	
	// UpdateChainExecution updates a chain execution
	UpdateChainExecution(ctx context.Context, executionID uuid.UUID, currentStep int, status AgentStatus, err error) error
}

// ExecutionFilter filters for listing executions
type ExecutionFilter struct {
	UserID     *uuid.UUID
	AgentType  *AgentType
	Status     *AgentStatus
	DateFrom   *time.Time
	DateTo     *time.Time
	ParentID   *uuid.UUID
	Limit      int
	Offset     int
}

// MetricsFilter filters for metrics
type MetricsFilter struct {
	AgentType *AgentType
	DateFrom  time.Time
	DateTo    time.Time
	GroupBy   string // hour, day, week, month
}

// ExecutionMetrics contains execution metrics
type ExecutionMetrics struct {
	TotalExecutions   int                       `json:"total_executions"`
	SuccessfulCount   int                       `json:"successful_count"`
	FailedCount       int                       `json:"failed_count"`
	AverageDuration   time.Duration             `json:"average_duration"`
	ByAgentType       map[AgentType]AgentMetrics `json:"by_agent_type"`
	ByTimeGroup       []TimeGroupMetrics        `json:"by_time_group"`
}

// AgentMetrics contains metrics for a specific agent type
type AgentMetrics struct {
	Count           int           `json:"count"`
	SuccessRate     float64       `json:"success_rate"`
	AverageDuration time.Duration `json:"average_duration"`
}

// TimeGroupMetrics contains metrics for a time period
type TimeGroupMetrics struct {
	Time            time.Time `json:"time"`
	Count           int       `json:"count"`
	SuccessfulCount int       `json:"successful_count"`
	FailedCount     int       `json:"failed_count"`
}

// DBExecutionTracker implements ExecutionTracker using database
type DBExecutionTracker struct {
	db *gorm.DB
	mu sync.RWMutex
}

// NewDBExecutionTracker creates a new database execution tracker
func NewDBExecutionTracker(db *gorm.DB) *DBExecutionTracker {
	// Auto migrate tables
	db.AutoMigrate(&AgentExecutionRecord{}, &AgentChainRecord{}, &AgentChainExecutionRecord{})
	
	return &DBExecutionTracker{
		db: db,
	}
}

// AgentExecutionRecord is the database model for agent executions
type AgentExecutionRecord struct {
	ID                uuid.UUID       `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	AgentType         string          `gorm:"type:varchar(50);not null;index"`
	AgentName         string          `gorm:"type:varchar(100);not null"`
	Status            string          `gorm:"type:varchar(20);not null;index"`
	Input             json.RawMessage `gorm:"type:jsonb;not null"`
	Output            json.RawMessage `gorm:"type:jsonb"`
	StartedAt         time.Time       `gorm:"not null;index"`
	CompletedAt       *time.Time      `gorm:"index"`
	DurationMs        *int64          
	Error             string          `gorm:"type:text"`
	UserID            uuid.UUID       `gorm:"type:uuid;not null;index"`
	ParentExecutionID *uuid.UUID      `gorm:"type:uuid;index"`
	Metadata          json.RawMessage `gorm:"type:jsonb"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (AgentExecutionRecord) TableName() string {
	return "agent_executions"
}

// AgentChainRecord is the database model for agent chains
type AgentChainRecord struct {
	ID          uuid.UUID       `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string          `gorm:"type:varchar(100);not null"`
	Description string          `gorm:"type:text"`
	Steps       json.RawMessage `gorm:"type:jsonb;not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (AgentChainRecord) TableName() string {
	return "agent_chains"
}

// AgentChainExecutionRecord is the database model for chain executions
type AgentChainExecutionRecord struct {
	ID             uuid.UUID       `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ChainID        uuid.UUID       `gorm:"type:uuid;not null;index"`
	Status         string          `gorm:"type:varchar(20);not null;index"`
	CurrentStep    int             `gorm:"not null"`
	StepExecutions json.RawMessage `gorm:"type:jsonb"`
	StartedAt      time.Time       `gorm:"not null;index"`
	CompletedAt    *time.Time
	Error          string          `gorm:"type:text"`
	UserID         uuid.UUID       `gorm:"type:uuid;not null;index"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (AgentChainExecutionRecord) TableName() string {
	return "agent_chain_executions"
}

// StartExecution starts tracking a new execution
func (t *DBExecutionTracker) StartExecution(ctx context.Context, agentType AgentType, agentName string, input AgentInput, userID uuid.UUID) (*AgentExecution, error) {
	execution := &AgentExecution{
		ID:                uuid.New(),
		AgentType:         agentType,
		AgentName:         agentName,
		Status:            AgentStatusRunning,
		Input:             input,
		StartedAt:         time.Now(),
		UserID:            userID,
		ParentExecutionID: input.ParentExecutionID,
		Metadata:          make(map[string]interface{}),
	}
	
	// Convert to record
	inputJSON, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input: %w", err)
	}
	
	metadataJSON, err := json.Marshal(execution.Metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metadata: %w", err)
	}
	
	record := &AgentExecutionRecord{
		ID:                execution.ID,
		AgentType:         string(agentType),
		AgentName:         agentName,
		Status:            string(AgentStatusRunning),
		Input:             inputJSON,
		StartedAt:         execution.StartedAt,
		UserID:            userID,
		ParentExecutionID: input.ParentExecutionID,
		Metadata:          metadataJSON,
	}
	
	if err := t.db.WithContext(ctx).Create(record).Error; err != nil {
		return nil, fmt.Errorf("failed to create execution record: %w", err)
	}
	
	return execution, nil
}

// UpdateExecution updates an execution status
func (t *DBExecutionTracker) UpdateExecution(ctx context.Context, executionID uuid.UUID, status AgentStatus, output *AgentOutput, err error) error {
	updates := map[string]interface{}{
		"status": string(status),
	}
	
	completedAt := time.Now()
	updates["completed_at"] = completedAt
	
	// Calculate duration
	var record AgentExecutionRecord
	if err := t.db.WithContext(ctx).Where("id = ?", executionID).First(&record).Error; err != nil {
		return fmt.Errorf("failed to find execution: %w", err)
	}
	
	duration := completedAt.Sub(record.StartedAt)
	durationMs := duration.Milliseconds()
	updates["duration_ms"] = durationMs
	
	if output != nil {
		outputJSON, err := json.Marshal(output)
		if err != nil {
			return fmt.Errorf("failed to marshal output: %w", err)
		}
		updates["output"] = outputJSON
	}
	
	if err != nil {
		updates["error"] = err.Error()
		updates["status"] = string(AgentStatusFailed)
	}
	
	if err := t.db.WithContext(ctx).Model(&AgentExecutionRecord{}).
		Where("id = ?", executionID).
		Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update execution: %w", err)
	}
	
	return nil
}

// GetExecution retrieves an execution by ID
func (t *DBExecutionTracker) GetExecution(ctx context.Context, executionID uuid.UUID) (*AgentExecution, error) {
	var record AgentExecutionRecord
	if err := t.db.WithContext(ctx).Where("id = ?", executionID).First(&record).Error; err != nil {
		return nil, fmt.Errorf("failed to find execution: %w", err)
	}
	
	return t.recordToExecution(&record)
}

// ListExecutions lists executions with filters
func (t *DBExecutionTracker) ListExecutions(ctx context.Context, filter ExecutionFilter) ([]*AgentExecution, error) {
	query := t.db.WithContext(ctx).Model(&AgentExecutionRecord{})
	
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.AgentType != nil {
		query = query.Where("agent_type = ?", string(*filter.AgentType))
	}
	if filter.Status != nil {
		query = query.Where("status = ?", string(*filter.Status))
	}
	if filter.DateFrom != nil {
		query = query.Where("started_at >= ?", *filter.DateFrom)
	}
	if filter.DateTo != nil {
		query = query.Where("started_at <= ?", *filter.DateTo)
	}
	if filter.ParentID != nil {
		query = query.Where("parent_execution_id = ?", *filter.ParentID)
	}
	
	query = query.Order("started_at DESC")
	
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}
	
	var records []AgentExecutionRecord
	if err := query.Find(&records).Error; err != nil {
		return nil, fmt.Errorf("failed to list executions: %w", err)
	}
	
	executions := make([]*AgentExecution, len(records))
	for i, record := range records {
		execution, err := t.recordToExecution(&record)
		if err != nil {
			return nil, err
		}
		executions[i] = execution
	}
	
	return executions, nil
}

// GetExecutionMetrics gets execution metrics
func (t *DBExecutionTracker) GetExecutionMetrics(ctx context.Context, filter MetricsFilter) (*ExecutionMetrics, error) {
	metrics := &ExecutionMetrics{
		ByAgentType: make(map[AgentType]AgentMetrics),
	}
	
	// Base query
	query := t.db.WithContext(ctx).Model(&AgentExecutionRecord{}).
		Where("started_at >= ? AND started_at <= ?", filter.DateFrom, filter.DateTo)
	
	if filter.AgentType != nil {
		query = query.Where("agent_type = ?", string(*filter.AgentType))
	}
	
	// Total count
	query.Count(&metrics.TotalExecutions)
	
	// Success/failed counts
	t.db.WithContext(ctx).Model(&AgentExecutionRecord{}).
		Where("started_at >= ? AND started_at <= ? AND status = ?", filter.DateFrom, filter.DateTo, string(AgentStatusCompleted)).
		Count(&metrics.SuccessfulCount)
	
	metrics.FailedCount = metrics.TotalExecutions - metrics.SuccessfulCount
	
	// Average duration
	var avgDuration float64
	t.db.WithContext(ctx).Model(&AgentExecutionRecord{}).
		Where("started_at >= ? AND started_at <= ? AND duration_ms IS NOT NULL", filter.DateFrom, filter.DateTo).
		Select("AVG(duration_ms)").
		Scan(&avgDuration)
	
	metrics.AverageDuration = time.Duration(avgDuration) * time.Millisecond
	
	// By agent type metrics
	var agentTypes []string
	t.db.WithContext(ctx).Model(&AgentExecutionRecord{}).
		Distinct("agent_type").
		Pluck("agent_type", &agentTypes)
	
	for _, at := range agentTypes {
		agentType := AgentType(at)
		var count int64
		var successCount int64
		var avgDur float64
		
		t.db.WithContext(ctx).Model(&AgentExecutionRecord{}).
			Where("started_at >= ? AND started_at <= ? AND agent_type = ?", filter.DateFrom, filter.DateTo, at).
			Count(&count)
		
		t.db.WithContext(ctx).Model(&AgentExecutionRecord{}).
			Where("started_at >= ? AND started_at <= ? AND agent_type = ? AND status = ?", 
				filter.DateFrom, filter.DateTo, at, string(AgentStatusCompleted)).
			Count(&successCount)
		
		t.db.WithContext(ctx).Model(&AgentExecutionRecord{}).
			Where("started_at >= ? AND started_at <= ? AND agent_type = ? AND duration_ms IS NOT NULL", 
				filter.DateFrom, filter.DateTo, at).
			Select("AVG(duration_ms)").
			Scan(&avgDur)
		
		metrics.ByAgentType[agentType] = AgentMetrics{
			Count:           int(count),
			SuccessRate:     float64(successCount) / float64(count),
			AverageDuration: time.Duration(avgDur) * time.Millisecond,
		}
	}
	
	return metrics, nil
}

// Helper methods

func (t *DBExecutionTracker) recordToExecution(record *AgentExecutionRecord) (*AgentExecution, error) {
	execution := &AgentExecution{
		ID:                record.ID,
		AgentType:         AgentType(record.AgentType),
		AgentName:         record.AgentName,
		Status:            AgentStatus(record.Status),
		StartedAt:         record.StartedAt,
		CompletedAt:       record.CompletedAt,
		Error:             record.Error,
		UserID:            record.UserID,
		ParentExecutionID: record.ParentExecutionID,
	}
	
	if err := json.Unmarshal(record.Input, &execution.Input); err != nil {
		return nil, fmt.Errorf("failed to unmarshal input: %w", err)
	}
	
	if record.Output != nil {
		var output AgentOutput
		if err := json.Unmarshal(record.Output, &output); err != nil {
			return nil, fmt.Errorf("failed to unmarshal output: %w", err)
		}
		execution.Output = &output
	}
	
	if record.Metadata != nil {
		if err := json.Unmarshal(record.Metadata, &execution.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
	}
	
	if record.DurationMs != nil {
		duration := time.Duration(*record.DurationMs) * time.Millisecond
		execution.Duration = &duration
	}
	
	return execution, nil
}

// StartChainExecution starts tracking a chain execution
func (t *DBExecutionTracker) StartChainExecution(ctx context.Context, chainID uuid.UUID, userID uuid.UUID) (*AgentChainExecution, error) {
	execution := &AgentChainExecution{
		ID:             uuid.New(),
		ChainID:        chainID,
		Status:         AgentStatusRunning,
		CurrentStep:    0,
		StepExecutions: []AgentExecution{},
		StartedAt:      time.Now(),
		UserID:         userID,
	}
	
	stepExecutionsJSON, _ := json.Marshal(execution.StepExecutions)
	
	record := &AgentChainExecutionRecord{
		ID:             execution.ID,
		ChainID:        chainID,
		Status:         string(AgentStatusRunning),
		CurrentStep:    0,
		StepExecutions: stepExecutionsJSON,
		StartedAt:      execution.StartedAt,
		UserID:         userID,
	}
	
	if err := t.db.WithContext(ctx).Create(record).Error; err != nil {
		return nil, fmt.Errorf("failed to create chain execution record: %w", err)
	}
	
	return execution, nil
}

// UpdateChainExecution updates a chain execution
func (t *DBExecutionTracker) UpdateChainExecution(ctx context.Context, executionID uuid.UUID, currentStep int, status AgentStatus, err error) error {
	updates := map[string]interface{}{
		"status":       string(status),
		"current_step": currentStep,
	}
	
	if status == AgentStatusCompleted || status == AgentStatusFailed || status == AgentStatusCancelled {
		completedAt := time.Now()
		updates["completed_at"] = completedAt
	}
	
	if err != nil {
		updates["error"] = err.Error()
	}
	
	if err := t.db.WithContext(ctx).Model(&AgentChainExecutionRecord{}).
		Where("id = ?", executionID).
		Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update chain execution: %w", err)
	}
	
	return nil
}