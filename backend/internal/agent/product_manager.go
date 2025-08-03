package agent

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// ProductManagerAgent implements the Product Manager AI agent
type ProductManagerAgent struct {
	name        string
	description string
}

// NewProductManagerAgent creates a new Product Manager agent
func NewProductManagerAgent() *ProductManagerAgent {
	return &ProductManagerAgent{
		name:        "Product Manager Agent",
		description: "Analyzes requirements, creates PRDs, defines workflows and role permissions",
	}
}

// GetType returns the agent type
func (a *ProductManagerAgent) GetType() AgentType {
	return AgentTypeProductManager
}

// GetName returns the agent name
func (a *ProductManagerAgent) GetName() string {
	return a.name
}

// GetDescription returns the agent description
func (a *ProductManagerAgent) GetDescription() string {
	return a.description
}

// Validate validates the input
func (a *ProductManagerAgent) Validate(input AgentInput) error {
	if input.Task == "" {
		return fmt.Errorf("task is required")
	}
	
	// Check for required parameters
	if moduleName, ok := input.Parameters["module_name"].(string); !ok || moduleName == "" {
		return fmt.Errorf("module_name parameter is required")
	}
	
	return nil
}

// Execute runs the Product Manager agent task
func (a *ProductManagerAgent) Execute(ctx context.Context, input AgentInput) (AgentOutput, error) {
	// Validate input
	if err := a.Validate(input); err != nil {
		return AgentOutput{
			Success: false,
			Error:   err.Error(),
		}, err
	}
	
	moduleName := input.Parameters["module_name"].(string)
	targetRoles := input.Parameters["target_roles"]
	mainPurpose := input.Parameters["main_purpose"]
	specialLogic := input.Parameters["special_logic"]
	
	// Generate PRD
	prd := a.generatePRD(moduleName, input.Task, targetRoles, mainPurpose, specialLogic)
	
	// Create output files
	files := []FileOutput{
		{
			Name:        fmt.Sprintf("PRD_%s.md", strings.ReplaceAll(moduleName, " ", "_")),
			Path:        fmt.Sprintf("/docs/prd/PRD_%s.md", strings.ReplaceAll(moduleName, " ", "_")),
			Content:     prd,
			ContentType: "text/markdown",
		},
	}
	
	// Generate next steps
	nextSteps := []string{
		"Review and validate the PRD with stakeholders",
		"Call UI Designer Agent to create wireframes based on this PRD",
		"Prepare detailed API specifications",
		"Define data models and database schema",
	}
	
	return AgentOutput{
		Success: true,
		Result: map[string]interface{}{
			"module_name": moduleName,
			"prd_summary": a.extractSummary(prd),
			"features":    a.extractFeatures(prd),
			"roles":       a.extractRoles(prd),
		},
		Files:     files,
		NextSteps: nextSteps,
		Metadata: map[string]interface{}{
			"generated_at": time.Now(),
			"version":      "1.0",
		},
	}, nil
}

// generatePRD generates a Product Requirements Document
func (a *ProductManagerAgent) generatePRD(moduleName, task string, targetRoles, mainPurpose, specialLogic interface{}) string {
	prd := fmt.Sprintf(`# 📘 FastenMind PRD Document - %s

## 1️⃣ Feature Overview
%s

### Main Purpose
%s

### Target Users
%s

## 2️⃣ User Scenarios and Roles
| Role | Permissions | Description |
|------|------------|-------------|
| Admin | Create, Read, Update, Delete | Full system access |
| Manager | Read, Update, Approve | Can review and approve |
| Engineer | Read, Update | Can view and update technical details |
| Sales | Create, Read | Can create new records and view |

## 3️⃣ Feature List
| ID | Feature Name | Description |
|----|--------------|-------------|
| F01 | Data Entry | Create and manage %s records |
| F02 | Search & Filter | Advanced search with multiple criteria |
| F03 | Approval Workflow | Multi-level approval process |
| F04 | Export & Reports | Generate reports in multiple formats |

## 4️⃣ Data Model (Draft)
### Main Entity: %s
| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Unique identifier |
| code | string | %s code |
| name | string | Display name |
| status | enum | Current status |
| created_at | timestamp | Creation time |
| updated_at | timestamp | Last update time |

## 5️⃣ Business Process Flow
1. User initiates request → System validates input
2. Data saved to draft status → Notification sent
3. Reviewer checks data → Approves or rejects
4. If approved → Status updated → Process continues
5. If rejected → Return to requester with comments

## 6️⃣ Permission Rules
- Each role can only access data within their scope
- Cross-company data requires special permissions
- Delete operations require manager approval
- Audit trail for all modifications

## 7️⃣ Validation Rules
- Required fields must be completed
- Business logic validation on submission
- Duplicate checking on key fields
- Data format validation

## 8️⃣ Success Criteria
- All CRUD operations working correctly
- Role-based access control enforced
- Search returns results within 2 seconds
- Export supports Excel, PDF, CSV formats

## 9️⃣ Technical Considerations
- RESTful API design
- Real-time updates via WebSocket
- Optimistic locking for concurrent edits
- Comprehensive error handling

## 🔄 Integration Points
- Email notifications on status changes
- Webhook support for external systems
- API endpoints for third-party integration
- Event-driven architecture for scalability

`,
		moduleName,
		task,
		a.formatValue(mainPurpose, "Enable efficient management and processing"),
		a.formatValue(targetRoles, "Business users and system administrators"),
		strings.ToLower(moduleName),
		strings.ToLower(moduleName),
		moduleName,
	)

	// Add special logic section if provided
	if specialLogic != nil && specialLogic != "" {
		prd += fmt.Sprintf(`
## 🎯 Special Business Logic
%s
`, a.formatValue(specialLogic, ""))
	}

	prd += `
## 📋 Appendix
- API endpoint specifications to be defined
- UI/UX guidelines to follow company standards
- Performance requirements: < 3s page load
- Security requirements: OAuth 2.0, encrypted data at rest
`

	return prd
}

// Helper methods

func (a *ProductManagerAgent) formatValue(value interface{}, defaultValue string) string {
	if value == nil {
		return defaultValue
	}
	
	switch v := value.(type) {
	case string:
		if v == "" {
			return defaultValue
		}
		return v
	case []string:
		if len(v) == 0 {
			return defaultValue
		}
		return strings.Join(v, ", ")
	default:
		return fmt.Sprintf("%v", v)
	}
}

func (a *ProductManagerAgent) extractSummary(prd string) string {
	lines := strings.Split(prd, "\n")
	for i, line := range lines {
		if strings.Contains(line, "Feature Overview") && i+1 < len(lines) {
			return strings.TrimSpace(lines[i+1])
		}
	}
	return "Feature requirements document generated"
}

func (a *ProductManagerAgent) extractFeatures(prd string) []string {
	features := []string{}
	lines := strings.Split(prd, "\n")
	inFeatureSection := false
	
	for _, line := range lines {
		if strings.Contains(line, "Feature List") {
			inFeatureSection = true
			continue
		}
		if inFeatureSection && strings.HasPrefix(line, "| F") {
			parts := strings.Split(line, "|")
			if len(parts) >= 3 {
				feature := strings.TrimSpace(parts[2])
				if feature != "" && feature != "Feature Name" {
					features = append(features, feature)
				}
			}
		}
		if inFeatureSection && line == "" {
			break
		}
	}
	
	return features
}

func (a *ProductManagerAgent) extractRoles(prd string) []string {
	roles := []string{}
	lines := strings.Split(prd, "\n")
	inRoleSection := false
	
	for _, line := range lines {
		if strings.Contains(line, "User Scenarios and Roles") {
			inRoleSection = true
			continue
		}
		if inRoleSection && strings.HasPrefix(line, "|") && !strings.Contains(line, "Role") {
			parts := strings.Split(line, "|")
			if len(parts) >= 2 {
				role := strings.TrimSpace(parts[1])
				if role != "" && role != "------" {
					roles = append(roles, role)
				}
			}
		}
		if inRoleSection && line == "" {
			break
		}
	}
	
	return roles
}