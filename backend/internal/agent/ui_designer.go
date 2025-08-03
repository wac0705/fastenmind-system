package agent

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// UIDesignerAgent implements the UI Designer AI agent
type UIDesignerAgent struct {
	name        string
	description string
}

// NewUIDesignerAgent creates a new UI Designer agent
func NewUIDesignerAgent() *UIDesignerAgent {
	return &UIDesignerAgent{
		name:        "UI Designer Agent",
		description: "Creates wireframes, UI components, and design specifications based on PRD",
	}
}

// GetType returns the agent type
func (a *UIDesignerAgent) GetType() AgentType {
	return AgentTypeUIDesigner
}

// GetName returns the agent name
func (a *UIDesignerAgent) GetName() string {
	return a.name
}

// GetDescription returns the agent description
func (a *UIDesignerAgent) GetDescription() string {
	return a.description
}

// Validate validates the input
func (a *UIDesignerAgent) Validate(input AgentInput) error {
	if input.Task == "" {
		return fmt.Errorf("task is required")
	}
	
	// Check if PRD content is provided
	prdProvided := false
	if _, ok := input.Context["prd_content"]; ok {
		prdProvided = true
	}
	if _, ok := input.Parameters["prd_content"]; ok {
		prdProvided = true
	}
	if len(input.Files) > 0 {
		prdProvided = true
	}
	
	if !prdProvided {
		return fmt.Errorf("PRD content is required (via context, parameters, or files)")
	}
	
	return nil
}

// Execute runs the UI Designer agent task
func (a *UIDesignerAgent) Execute(ctx context.Context, input AgentInput) (AgentOutput, error) {
	// Validate input
	if err := a.Validate(input); err != nil {
		return AgentOutput{
			Success: false,
			Error:   err.Error(),
		}, err
	}
	
	// Extract PRD content
	prdContent := a.extractPRDContent(input)
	
	// Generate design specification
	designSpec := a.generateDesignSpec(input.Task, prdContent)
	
	// Generate component specifications
	componentSpec := a.generateComponentSpec(input.Task)
	
	// Generate wireframe description
	wireframe := a.generateWireframe(input.Task)
	
	// Create output files
	files := []FileOutput{
		{
			Name:        "DESIGN_SPEC.md",
			Path:        "/docs/design/DESIGN_SPEC.md",
			Content:     designSpec,
			ContentType: "text/markdown",
		},
		{
			Name:        "COMPONENT_SPEC.md",
			Path:        "/docs/design/COMPONENT_SPEC.md",
			Content:     componentSpec,
			ContentType: "text/markdown",
		},
		{
			Name:        "WIREFRAME.md",
			Path:        "/docs/design/WIREFRAME.md",
			Content:     wireframe,
			ContentType: "text/markdown",
		},
	}
	
	// Generate next steps
	nextSteps := []string{
		"Review design specifications with stakeholders",
		"Create high-fidelity mockups if needed",
		"Call Engineer Agent to implement the UI",
		"Prepare style guide and design tokens",
	}
	
	return AgentOutput{
		Success: true,
		Result: map[string]interface{}{
			"design_spec":     designSpec,
			"component_count": a.countComponents(componentSpec),
			"page_count":      a.countPages(wireframe),
			"ui_framework":    "Next.js + Tailwind CSS",
		},
		Files:     files,
		NextSteps: nextSteps,
		Metadata: map[string]interface{}{
			"generated_at": time.Now(),
			"version":      "1.0",
			"framework":    "React/Next.js",
		},
	}, nil
}

// extractPRDContent extracts PRD content from various sources
func (a *UIDesignerAgent) extractPRDContent(input AgentInput) string {
	// Check context first
	if prd, ok := input.Context["prd_content"].(string); ok {
		return prd
	}
	
	// Check parameters
	if prd, ok := input.Parameters["prd_content"].(string); ok {
		return prd
	}
	
	// Check files
	for _, file := range input.Files {
		if strings.Contains(file.Name, "PRD") {
			return file.Content
		}
	}
	
	return ""
}

// generateDesignSpec generates the design specification
func (a *UIDesignerAgent) generateDesignSpec(task, prdContent string) string {
	spec := fmt.Sprintf(`# 🎨 FastenMind Design Specification

## Overview
This document outlines the UI/UX design for: %s

## Design Principles
1. **Clarity**: Clear visual hierarchy and intuitive navigation
2. **Efficiency**: Minimize clicks and optimize workflows
3. **Consistency**: Uniform design patterns across all screens
4. **Accessibility**: WCAG 2.1 AA compliance
5. **Responsiveness**: Mobile-first approach

## Color Palette
| Usage | Color | Hex Code |
|-------|-------|----------|
| Primary | Blue | #3B82F6 |
| Secondary | Gray | #6B7280 |
| Success | Green | #10B981 |
| Warning | Yellow | #F59E0B |
| Error | Red | #EF4444 |
| Background | White | #FFFFFF |
| Text Primary | Dark Gray | #1F2937 |

## Typography
- **Font Family**: Inter, system-ui, sans-serif
- **Headings**: 
  - H1: 36px/44px, font-weight: 700
  - H2: 30px/36px, font-weight: 600
  - H3: 24px/32px, font-weight: 600
- **Body**: 16px/24px, font-weight: 400
- **Small**: 14px/20px, font-weight: 400

## Layout Structure
### Grid System
- 12-column grid
- Gutter: 24px
- Margin: 24px (desktop), 16px (mobile)
- Max container width: 1280px

### Breakpoints
- Mobile: < 640px
- Tablet: 640px - 1024px
- Desktop: > 1024px

## Component Library
Based on the PRD requirements, the following components are needed:
- Navigation (Header, Sidebar, Breadcrumbs)
- Forms (Input, Select, Checkbox, Radio, DatePicker)
- Buttons (Primary, Secondary, Tertiary, Icon)
- Tables (Data Grid, Sorting, Filtering, Pagination)
- Cards (Info Card, Stats Card, Action Card)
- Modals (Confirmation, Form Modal, Info Modal)
- Notifications (Toast, Alert, Badge)

## Interaction Patterns
1. **Loading States**: Skeleton screens for better perceived performance
2. **Error Handling**: Inline validation with clear error messages
3. **Success Feedback**: Toast notifications for successful actions
4. **Empty States**: Helpful messages and action prompts
5. **Hover Effects**: Subtle transitions (0.2s ease)

## Accessibility Requirements
- Keyboard navigation support
- Screen reader compatibility
- Color contrast ratio: minimum 4.5:1
- Focus indicators on interactive elements
- ARIA labels for complex components

## Animation Guidelines
- Duration: 200-300ms for micro-interactions
- Easing: ease-in-out for smooth transitions
- Purpose: Enhance usability, not decoration
- Performance: Use CSS transforms, avoid layout shifts

`, task)

	// Add PRD-specific sections if available
	if prdContent != "" {
		spec += `
## PRD-Based UI Requirements
Based on the product requirements, specific UI considerations include:
- Role-based views and permissions
- Multi-step workflows with progress indicators
- Real-time data updates where applicable
- Export functionality with format options
`
	}

	return spec
}

// generateComponentSpec generates component specifications
func (a *UIDesignerAgent) generateComponentSpec(task string) string {
	return fmt.Sprintf(`# 🧩 Component Specifications

## Core Components

### 1. DataTable Component
\`\`\`typescript
interface DataTableProps {
  columns: ColumnDef[]
  data: any[]
  pagination?: boolean
  sorting?: boolean
  filtering?: boolean
  selection?: boolean
  actions?: ActionDef[]
}
\`\`\`

**Features:**
- Server-side pagination
- Multi-column sorting
- Advanced filtering
- Row selection
- Bulk actions
- Export functionality

### 2. FormBuilder Component
\`\`\`typescript
interface FormBuilderProps {
  schema: FormSchema
  initialValues?: Record<string, any>
  onSubmit: (values: any) => Promise<void>
  validation?: ValidationRules
}
\`\`\`

**Features:**
- Dynamic field rendering
- Conditional logic
- Real-time validation
- File upload support
- Auto-save drafts

### 3. StatusWorkflow Component
\`\`\`typescript
interface StatusWorkflowProps {
  currentStatus: string
  availableTransitions: Transition[]
  onTransition: (newStatus: string) => Promise<void>
  timeline?: TimelineEvent[]
}
\`\`\`

**Features:**
- Visual status indicator
- Action buttons for transitions
- Timeline view
- Permission-based actions

### 4. SearchFilter Component
\`\`\`typescript
interface SearchFilterProps {
  fields: FilterField[]
  onSearch: (filters: FilterValues) => void
  savedFilters?: SavedFilter[]
  allowSave?: boolean
}
\`\`\`

**Features:**
- Multi-field search
- Date range picker
- Saved filter sets
- Quick filters
- Clear all option

### 5. Dashboard Widget
\`\`\`typescript
interface DashboardWidgetProps {
  title: string
  type: 'chart' | 'stat' | 'list' | 'progress'
  data: any
  refreshInterval?: number
  actions?: WidgetAction[]
}
\`\`\`

**Features:**
- Multiple visualization types
- Auto-refresh
- Drill-down capability
- Export widget data

## Shared UI Patterns

### Loading States
- Skeleton screens for initial load
- Spinner for actions
- Progress bar for long operations

### Error States
- Inline field errors
- Form submission errors
- Network error handling
- Fallback UI components

### Empty States
- Informative messages
- Action prompts
- Illustrations where appropriate

### Success States
- Toast notifications
- Success pages
- Confirmation modals

## Design Tokens
\`\`\`javascript
export const tokens = {
  colors: {
    primary: '#3B82F6',
    secondary: '#6B7280',
    // ... rest of palette
  },
  spacing: {
    xs: '4px',
    sm: '8px',
    md: '16px',
    lg: '24px',
    xl: '32px',
  },
  borderRadius: {
    sm: '4px',
    md: '8px',
    lg: '12px',
    full: '9999px',
  },
  shadows: {
    sm: '0 1px 2px rgba(0, 0, 0, 0.05)',
    md: '0 4px 6px rgba(0, 0, 0, 0.1)',
    lg: '0 10px 15px rgba(0, 0, 0, 0.1)',
  },
}
\`\`\`
`}

// generateWireframe generates wireframe descriptions
func (a *UIDesignerAgent) generateWireframe(task string) string {
	return fmt.Sprintf(`# 📐 Wireframe Specifications

## Page Layouts

### 1. List View Page
\`\`\`
+--------------------------------------------------+
| [Logo] Navigation Menu            [User] [Logout] |
+--------------------------------------------------+
| Breadcrumb > Path > Current Page                  |
+--------------------------------------------------+
| Page Title                    [+ New] [⚙ Actions] |
+--------------------------------------------------+
| [🔍 Search.....................] [Filter] [Sort]  |
+--------------------------------------------------+
| □ | Column 1 | Column 2 | Column 3 | Actions     |
| □ | Data     | Data     | Data     | [👁] [✏] [🗑] |
| □ | Data     | Data     | Data     | [👁] [✏] [🗑] |
| □ | Data     | Data     | Data     | [👁] [✏] [🗑] |
+--------------------------------------------------+
| Showing 1-10 of 100 | [<] [1] 2 3 ... 10 [>]    |
+--------------------------------------------------+
\`\`\`

### 2. Detail View Page
\`\`\`
+--------------------------------------------------+
| [Logo] Navigation Menu            [User] [Logout] |
+--------------------------------------------------+
| [<-] Back to List | Record Name                  |
+--------------------------------------------------+
| Tab 1 | Tab 2 | Tab 3                   [Edit]   |
+--------------------------------------------------+
| Section Header                                    |
| Field Label    | Field Value                      |
| Field Label    | Field Value                      |
|                                                   |
| Section Header                                    |
| Field Label    | Field Value                      |
| Field Label    | Field Value                      |
|                                                   |
| Activity Timeline                                 |
| • Action performed by User (2 hours ago)          |
| • Status changed by User (1 day ago)              |
+--------------------------------------------------+
\`\`\`

### 3. Form Page (Create/Edit)
\`\`\`
+--------------------------------------------------+
| [Logo] Navigation Menu            [User] [Logout] |
+--------------------------------------------------+
| [X] Create New Record                             |
+--------------------------------------------------+
| Basic Information                                 |
| Field Label *                                     |
| [Input Field.............................]        |
|                                                   |
| Field Label                                       |
| [Dropdown Selection         v]                    |
|                                                   |
| Field Label                                       |
| [📅 Date Picker................]                  |
|                                                   |
| Additional Details                                |
| Field Label                                       |
| [Text Area..............................]         |
| [......................................]         |
|                                                   |
| Attachments                                       |
| [📎 Drop files here or click to upload]           |
|                                                   |
| [Cancel]                    [Save Draft] [Submit] |
+--------------------------------------------------+
\`\`\`

### 4. Dashboard Page
\`\`\`
+--------------------------------------------------+
| [Logo] Navigation Menu            [User] [Logout] |
+--------------------------------------------------+
| Dashboard                        [Date Range ▼]   |
+--------------------------------------------------+
| +----------------+ +----------------+ +---------+ |
| | Total Records  | | Pending Items  | | Success | |
| | 1,234          | | 56             | | 95.2%   | |
| | ↑ 12.3%        | | ↓ 8.5%         | | ↑ 2.1%  | |
| +----------------+ +----------------+ +---------+ |
|                                                   |
| +------------------------+ +---------------------+ |
| | 📊 Chart Title         | | Recent Activities   | |
| | [Chart Area]           | | • Item 1 (2m ago)   | |
| | [.............]        | | • Item 2 (1h ago)   | |
| | [.............]        | | • Item 3 (3h ago)   | |
| +------------------------+ | [View All]          | |
|                           +---------------------+ |
+--------------------------------------------------+
\`\`\`

## Mobile Responsive Layouts

### Mobile List View
\`\`\`
+----------------------+
| ☰ Menu  Logo    [👤] |
+----------------------+
| Page Title           |
| [🔍 Search...]  [+]  |
+----------------------+
| Card Item 1          |
| Subtitle info    [>] |
+----------------------+
| Card Item 2          |
| Subtitle info    [>] |
+----------------------+
| [Load More]          |
+----------------------+
\`\`\`

### Mobile Form View
\`\`\`
+----------------------+
| [<] Create New       |
+----------------------+
| Field Label *        |
| [Input............]  |
|                      |
| Field Label          |
| [Select........v]    |
|                      |
| [Save] [Cancel]      |
+----------------------+
\`\`\`

## Interaction Flows

### Create New Record Flow
1. User clicks [+ New] button
2. Modal or new page opens with form
3. User fills required fields
4. Validation occurs on blur
5. User clicks [Submit]
6. Loading state shows
7. Success: Redirect to detail view
8. Error: Show inline errors

### Edit Record Flow
1. User clicks [Edit] on list or detail view
2. Form opens in edit mode
3. Current values pre-filled
4. User modifies fields
5. Auto-save draft (optional)
6. User clicks [Save]
7. Optimistic update with rollback on error

### Delete Record Flow
1. User clicks [Delete]
2. Confirmation modal appears
3. User confirms deletion
4. Soft delete performed
5. Success toast notification
6. Record removed from view

### Bulk Operations Flow
1. User selects multiple records
2. Bulk action menu appears
3. User selects action
4. Confirmation if destructive
5. Progress indicator shows
6. Results summary displayed
`}

// Helper methods

func (a *UIDesignerAgent) countComponents(componentSpec string) int {
	// Simple count of interface definitions
	return strings.Count(componentSpec, "interface")
}

func (a *UIDesignerAgent) countPages(wireframe string) int {
	// Count page layout sections
	return strings.Count(wireframe, "### ") - strings.Count(wireframe, "### Mobile")
}