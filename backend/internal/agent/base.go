package agent

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// AgentType represents the type of agent
type AgentType string

const (
	AgentTypeProductManager AgentType = "product_manager"
	AgentTypeUIDesigner     AgentType = "ui_designer"
	AgentTypeEngineer       AgentType = "engineer"
	AgentTypeN8N            AgentType = "n8n_automation"
)

// AgentStatus represents the status of an agent execution
type AgentStatus string

const (
	AgentStatusPending   AgentStatus = "pending"
	AgentStatusRunning   AgentStatus = "running"
	AgentStatusCompleted AgentStatus = "completed"
	AgentStatusFailed    AgentStatus = "failed"
	AgentStatusCancelled AgentStatus = "cancelled"
)

// Agent interface that all agents must implement
type Agent interface {
	// GetType returns the agent type
	GetType() AgentType
	
	// GetName returns the agent name
	GetName() string
	
	// GetDescription returns the agent description
	GetDescription() string
	
	// Execute runs the agent task
	Execute(ctx context.Context, input AgentInput) (AgentOutput, error)
	
	// Validate validates the input
	Validate(input AgentInput) error
}

// AgentInput represents input to an agent
type AgentInput struct {
	// Task description
	Task string `json:"task"`
	
	// Context from previous agents
	Context map[string]interface{} `json:"context"`
	
	// Parameters specific to the agent
	Parameters map[string]interface{} `json:"parameters"`
	
	// Files or documents
	Files []FileInput `json:"files,omitempty"`
	
	// Parent execution ID for tracking
	ParentExecutionID *uuid.UUID `json:"parent_execution_id,omitempty"`
}

// FileInput represents a file input
type FileInput struct {
	Name        string `json:"name"`
	Content     string `json:"content"`
	ContentType string `json:"content_type"`
}

// AgentOutput represents output from an agent
type AgentOutput struct {
	// Success indicates if the task was successful
	Success bool `json:"success"`
	
	// Result contains the main output
	Result interface{} `json:"result"`
	
	// Files generated
	Files []FileOutput `json:"files,omitempty"`
	
	// Metadata about the execution
	Metadata map[string]interface{} `json:"metadata"`
	
	// Error message if failed
	Error string `json:"error,omitempty"`
	
	// Next steps or recommendations
	NextSteps []string `json:"next_steps,omitempty"`
}

// FileOutput represents a file output
type FileOutput struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Content     string `json:"content"`
	ContentType string `json:"content_type"`
}

// AgentExecution tracks an agent execution
type AgentExecution struct {
	ID                uuid.UUID              `json:"id"`
	AgentType         AgentType              `json:"agent_type"`
	AgentName         string                 `json:"agent_name"`
	Status            AgentStatus            `json:"status"`
	Input             AgentInput             `json:"input"`
	Output            *AgentOutput           `json:"output,omitempty"`
	StartedAt         time.Time              `json:"started_at"`
	CompletedAt       *time.Time             `json:"completed_at,omitempty"`
	Duration          *time.Duration         `json:"duration,omitempty"`
	Error             string                 `json:"error,omitempty"`
	UserID            uuid.UUID              `json:"user_id"`
	ParentExecutionID *uuid.UUID             `json:"parent_execution_id,omitempty"`
	Metadata          map[string]interface{} `json:"metadata"`
}

// AgentChain represents a chain of agents to execute
type AgentChain struct {
	ID          uuid.UUID        `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Steps       []AgentChainStep `json:"steps"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

// AgentChainStep represents a step in an agent chain
type AgentChainStep struct {
	Order          int                    `json:"order"`
	AgentType      AgentType              `json:"agent_type"`
	Name           string                 `json:"name"`
	Description    string                 `json:"description"`
	InputMapping   map[string]string      `json:"input_mapping"`
	Parameters     map[string]interface{} `json:"parameters"`
	ContinueOnFail bool                   `json:"continue_on_fail"`
}

// AgentChainExecution tracks execution of an agent chain
type AgentChainExecution struct {
	ID              uuid.UUID         `json:"id"`
	ChainID         uuid.UUID         `json:"chain_id"`
	Status          AgentStatus       `json:"status"`
	CurrentStep     int               `json:"current_step"`
	StepExecutions  []AgentExecution  `json:"step_executions"`
	StartedAt       time.Time         `json:"started_at"`
	CompletedAt     *time.Time        `json:"completed_at,omitempty"`
	Error           string            `json:"error,omitempty"`
	UserID          uuid.UUID         `json:"user_id"`
}