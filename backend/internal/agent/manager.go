package agent

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

// Manager manages agent executions
type Manager interface {
	// RegisterAgent registers an agent
	RegisterAgent(agent Agent) error
	
	// GetAgent retrieves an agent by type
	GetAgent(agentType AgentType) (Agent, error)
	
	// ExecuteAgent executes a single agent
	ExecuteAgent(ctx context.Context, agentType AgentType, input AgentInput, userID uuid.UUID) (*AgentExecution, error)
	
	// ExecuteChain executes an agent chain
	ExecuteChain(ctx context.Context, chainID uuid.UUID, initialInput AgentInput, userID uuid.UUID) (*AgentChainExecution, error)
	
	// CreateChain creates a new agent chain
	CreateChain(chain *AgentChain) error
	
	// GetChain retrieves a chain by ID
	GetChain(chainID uuid.UUID) (*AgentChain, error)
	
	// ListAgents lists all registered agents
	ListAgents() []Agent
}

// DefaultManager is the default implementation of Manager
type DefaultManager struct {
	agents   map[AgentType]Agent
	chains   map[uuid.UUID]*AgentChain
	tracker  ExecutionTracker
	mu       sync.RWMutex
}

// NewManager creates a new agent manager
func NewManager(tracker ExecutionTracker) *DefaultManager {
	manager := &DefaultManager{
		agents:  make(map[AgentType]Agent),
		chains:  make(map[uuid.UUID]*AgentChain),
		tracker: tracker,
	}
	
	// Register default agents
	manager.RegisterAgent(NewProductManagerAgent())
	manager.RegisterAgent(NewUIDesignerAgent())
	manager.RegisterAgent(NewEngineerAgent())
	manager.RegisterAgent(NewN8NAgent())
	
	return manager
}

// RegisterAgent registers an agent
func (m *DefaultManager) RegisterAgent(agent Agent) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if agent == nil {
		return fmt.Errorf("agent cannot be nil")
	}
	
	m.agents[agent.GetType()] = agent
	return nil
}

// GetAgent retrieves an agent by type
func (m *DefaultManager) GetAgent(agentType AgentType) (Agent, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	agent, exists := m.agents[agentType]
	if !exists {
		return nil, fmt.Errorf("agent not found: %s", agentType)
	}
	
	return agent, nil
}

// ExecuteAgent executes a single agent
func (m *DefaultManager) ExecuteAgent(ctx context.Context, agentType AgentType, input AgentInput, userID uuid.UUID) (*AgentExecution, error) {
	// Get agent
	agent, err := m.GetAgent(agentType)
	if err != nil {
		return nil, err
	}
	
	// Start tracking execution
	execution, err := m.tracker.StartExecution(ctx, agentType, agent.GetName(), input, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to start execution tracking: %w", err)
	}
	
	// Execute agent
	output, err := agent.Execute(ctx, input)
	
	// Update execution status
	status := AgentStatusCompleted
	if err != nil {
		status = AgentStatusFailed
	}
	
	if updateErr := m.tracker.UpdateExecution(ctx, execution.ID, status, &output, err); updateErr != nil {
		// Log error but don't fail the execution
		fmt.Printf("Failed to update execution tracking: %v\n", updateErr)
	}
	
	// Update execution object
	execution.Status = status
	execution.Output = &output
	if err != nil {
		execution.Error = err.Error()
	}
	
	return execution, err
}

// ExecuteChain executes an agent chain
func (m *DefaultManager) ExecuteChain(ctx context.Context, chainID uuid.UUID, initialInput AgentInput, userID uuid.UUID) (*AgentChainExecution, error) {
	// Get chain
	chain, err := m.GetChain(chainID)
	if err != nil {
		return nil, err
	}
	
	// Start chain execution tracking
	chainExecution, err := m.tracker.StartChainExecution(ctx, chainID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to start chain execution tracking: %w", err)
	}
	
	// Execute each step
	currentInput := initialInput
	currentContext := make(map[string]interface{})
	
	for i, step := range chain.Steps {
		// Update current step
		if err := m.tracker.UpdateChainExecution(ctx, chainExecution.ID, i, AgentStatusRunning, nil); err != nil {
			fmt.Printf("Failed to update chain execution: %v\n", err)
		}
		
		// Prepare step input
		stepInput := AgentInput{
			Task:              currentInput.Task,
			Context:           currentContext,
			Parameters:        step.Parameters,
			ParentExecutionID: &chainExecution.ID,
		}
		
		// Apply input mapping
		for target, source := range step.InputMapping {
			if value, exists := currentContext[source]; exists {
				stepInput.Parameters[target] = value
			}
		}
		
		// Execute step
		stepExecution, err := m.ExecuteAgent(ctx, step.AgentType, stepInput, userID)
		if err != nil {
			if !step.ContinueOnFail {
				// Update chain status to failed
				m.tracker.UpdateChainExecution(ctx, chainExecution.ID, i, AgentStatusFailed, err)
				return chainExecution, fmt.Errorf("step %d failed: %w", i, err)
			}
			// Continue despite failure
			fmt.Printf("Step %d failed but continuing: %v\n", i, err)
		}
		
		// Add step execution to chain
		chainExecution.StepExecutions = append(chainExecution.StepExecutions, *stepExecution)
		
		// Update context with step output
		if stepExecution.Output != nil && stepExecution.Output.Result != nil {
			if resultMap, ok := stepExecution.Output.Result.(map[string]interface{}); ok {
				for k, v := range resultMap {
					currentContext[fmt.Sprintf("step_%d_%s", i, k)] = v
				}
			}
		}
		
		// Use output files as input for next step if needed
		if stepExecution.Output != nil && len(stepExecution.Output.Files) > 0 {
			currentInput.Files = make([]FileInput, len(stepExecution.Output.Files))
			for j, file := range stepExecution.Output.Files {
				currentInput.Files[j] = FileInput{
					Name:        file.Name,
					Content:     file.Content,
					ContentType: file.ContentType,
				}
			}
		}
	}
	
	// Update chain execution as completed
	if err := m.tracker.UpdateChainExecution(ctx, chainExecution.ID, len(chain.Steps), AgentStatusCompleted, nil); err != nil {
		fmt.Printf("Failed to update chain execution: %v\n", err)
	}
	
	chainExecution.Status = AgentStatusCompleted
	chainExecution.CurrentStep = len(chain.Steps)
	
	return chainExecution, nil
}

// CreateChain creates a new agent chain
func (m *DefaultManager) CreateChain(chain *AgentChain) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if chain == nil {
		return fmt.Errorf("chain cannot be nil")
	}
	
	if chain.ID == uuid.Nil {
		chain.ID = uuid.New()
	}
	
	// Validate chain steps
	for i, step := range chain.Steps {
		if _, exists := m.agents[step.AgentType]; !exists {
			return fmt.Errorf("step %d: unknown agent type: %s", i, step.AgentType)
		}
	}
	
	m.chains[chain.ID] = chain
	return nil
}

// GetChain retrieves a chain by ID
func (m *DefaultManager) GetChain(chainID uuid.UUID) (*AgentChain, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	chain, exists := m.chains[chainID]
	if !exists {
		return nil, fmt.Errorf("chain not found: %s", chainID)
	}
	
	return chain, nil
}

// ListAgents lists all registered agents
func (m *DefaultManager) ListAgents() []Agent {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	agents := make([]Agent, 0, len(m.agents))
	for _, agent := range m.agents {
		agents = append(agents, agent)
	}
	
	return agents
}

// CreateStandardChains creates standard agent chains
func (m *DefaultManager) CreateStandardChains() error {
	// PRD to Implementation chain
	prdToImplementation := &AgentChain{
		ID:          uuid.New(),
		Name:        "PRD to Implementation",
		Description: "Complete flow from PRD to working code",
		Steps: []AgentChainStep{
			{
				Order:       1,
				AgentType:   AgentTypeProductManager,
				Name:        "Create PRD",
				Description: "Analyze requirements and create PRD",
				Parameters:  map[string]interface{}{},
				InputMapping: map[string]string{
					"task": "task",
				},
			},
			{
				Order:       2,
				AgentType:   AgentTypeUIDesigner,
				Name:        "Design UI",
				Description: "Create wireframes and UI specifications",
				Parameters:  map[string]interface{}{},
				InputMapping: map[string]string{
					"prd_content": "step_0_prd_content",
				},
			},
			{
				Order:       3,
				AgentType:   AgentTypeEngineer,
				Name:        "Implement Code",
				Description: "Generate frontend and backend code",
				Parameters:  map[string]interface{}{},
				InputMapping: map[string]string{
					"design_spec": "step_1_design_spec",
					"prd_content": "step_0_prd_content",
				},
			},
			{
				Order:          4,
				AgentType:      AgentTypeN8N,
				Name:           "Setup Automation",
				Description:    "Configure N8N workflows",
				Parameters:     map[string]interface{}{},
				ContinueOnFail: true, // Optional step
				InputMapping: map[string]string{
					"module_info": "step_2_module_info",
				},
			},
		},
	}
	
	if err := m.CreateChain(prdToImplementation); err != nil {
		return fmt.Errorf("failed to create PRD to Implementation chain: %w", err)
	}
	
	// Quick Implementation chain (skip PRD)
	quickImplementation := &AgentChain{
		ID:          uuid.New(),
		Name:        "Quick Implementation",
		Description: "Direct UI design and implementation",
		Steps: []AgentChainStep{
			{
				Order:       1,
				AgentType:   AgentTypeUIDesigner,
				Name:        "Design UI",
				Description: "Create wireframes and UI specifications",
				Parameters:  map[string]interface{}{},
			},
			{
				Order:       2,
				AgentType:   AgentTypeEngineer,
				Name:        "Implement Code",
				Description: "Generate frontend and backend code",
				Parameters:  map[string]interface{}{},
				InputMapping: map[string]string{
					"design_spec": "step_0_design_spec",
				},
			},
		},
	}
	
	if err := m.CreateChain(quickImplementation); err != nil {
		return fmt.Errorf("failed to create Quick Implementation chain: %w", err)
	}
	
	return nil
}