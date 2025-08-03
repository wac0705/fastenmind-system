package agent

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// EngineerAgent implements the Engineer AI agent
type EngineerAgent struct {
	name        string
	description string
}

// NewEngineerAgent creates a new Engineer agent
func NewEngineerAgent() *EngineerAgent {
	return &EngineerAgent{
		name:        "Engineer Agent",
		description: "Generates frontend pages, API endpoints, and backend logic based on design specs",
	}
}

// GetType returns the agent type
func (a *EngineerAgent) GetType() AgentType {
	return AgentTypeEngineer
}

// GetName returns the agent name
func (a *EngineerAgent) GetName() string {
	return a.name
}

// GetDescription returns the agent description
func (a *EngineerAgent) GetDescription() string {
	return a.description
}

// Validate validates the input
func (a *EngineerAgent) Validate(input AgentInput) error {
	if input.Task == "" {
		return fmt.Errorf("task is required")
	}
	
	// Check if design spec is provided
	designProvided := false
	if _, ok := input.Context["design_spec"]; ok {
		designProvided = true
	}
	if _, ok := input.Parameters["design_spec"]; ok {
		designProvided = true
	}
	if len(input.Files) > 0 {
		for _, file := range input.Files {
			if strings.Contains(file.Name, "DESIGN") {
				designProvided = true
				break
			}
		}
	}
	
	if !designProvided {
		return fmt.Errorf("design specification is required")
	}
	
	return nil
}

// Execute runs the Engineer agent task
func (a *EngineerAgent) Execute(ctx context.Context, input AgentInput) (AgentOutput, error) {
	// Validate input
	if err := a.Validate(input); err != nil {
		return AgentOutput{
			Success: false,
			Error:   err.Error(),
		}, err
	}
	
	// Extract module name
	moduleName := a.extractModuleName(input)
	
	// Generate React component
	componentCode := a.generateReactComponent(moduleName, input.Task)
	
	// Generate API handler
	apiCode := a.generateAPIHandler(moduleName)
	
	// Generate service layer
	serviceCode := a.generateServiceCode(moduleName)
	
	// Generate repository layer
	repositoryCode := a.generateRepositoryCode(moduleName)
	
	// Create output files
	files := []FileOutput{
		{
			Name:        fmt.Sprintf("%sList.tsx", moduleName),
			Path:        fmt.Sprintf("/frontend/pages/%s/list.tsx", strings.ToLower(moduleName)),
			Content:     componentCode,
			ContentType: "text/typescript",
		},
		{
			Name:        fmt.Sprintf("%s_handler.go", strings.ToLower(moduleName)),
			Path:        fmt.Sprintf("/backend/internal/handler/%s_handler.go", strings.ToLower(moduleName)),
			Content:     apiCode,
			ContentType: "text/x-go",
		},
		{
			Name:        fmt.Sprintf("%s_service.go", strings.ToLower(moduleName)),
			Path:        fmt.Sprintf("/backend/internal/service/%s_service.go", strings.ToLower(moduleName)),
			Content:     serviceCode,
			ContentType: "text/x-go",
		},
		{
			Name:        fmt.Sprintf("%s_repository.go", strings.ToLower(moduleName)),
			Path:        fmt.Sprintf("/backend/internal/repository/%s_repository.go", strings.ToLower(moduleName)),
			Content:     repositoryCode,
			ContentType: "text/x-go",
		},
	}
	
	// Generate next steps
	nextSteps := []string{
		"Review generated code for completeness",
		"Run linting and formatting tools",
		"Write unit tests for the components",
		"Integrate with existing routing",
		"Test API endpoints with Postman/Swagger",
		"Deploy to development environment",
	}
	
	return AgentOutput{
		Success: true,
		Result: map[string]interface{}{
			"module_name":     moduleName,
			"files_generated": len(files),
			"frontend_tech":   "React/Next.js + TypeScript",
			"backend_tech":    "Go + Echo Framework",
			"api_endpoints": []string{
				fmt.Sprintf("GET /api/v1/%s", strings.ToLower(moduleName)),
				fmt.Sprintf("GET /api/v1/%s/:id", strings.ToLower(moduleName)),
				fmt.Sprintf("POST /api/v1/%s", strings.ToLower(moduleName)),
				fmt.Sprintf("PUT /api/v1/%s/:id", strings.ToLower(moduleName)),
				fmt.Sprintf("DELETE /api/v1/%s/:id", strings.ToLower(moduleName)),
			},
		},
		Files:     files,
		NextSteps: nextSteps,
		Metadata: map[string]interface{}{
			"generated_at": time.Now(),
			"version":      "1.0",
		},
	}, nil
}

// extractModuleName extracts module name from input
func (a *EngineerAgent) extractModuleName(input AgentInput) string {
	// Check parameters first
	if name, ok := input.Parameters["module_name"].(string); ok && name != "" {
		return name
	}
	
	// Try to extract from task
	task := strings.ToLower(input.Task)
	if strings.Contains(task, "customer") {
		return "Customer"
	} else if strings.Contains(task, "product") {
		return "Product"
	} else if strings.Contains(task, "order") {
		return "Order"
	} else if strings.Contains(task, "inquiry") {
		return "Inquiry"
	} else if strings.Contains(task, "quote") {
		return "Quote"
	}
	
	return "Module"
}

// generateReactComponent generates a React component
func (a *EngineerAgent) generateReactComponent(moduleName, task string) string {
	lowerName := strings.ToLower(moduleName)
	return fmt.Sprintf(`import React, { useState, useEffect } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { toast } from 'react-hot-toast';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Plus, Edit, Trash2, Search } from 'lucide-react';
import { api } from '@/lib/api';

// Types
interface %s {
  id: string;
  code: string;
  name: string;
  description?: string;
  status: 'active' | 'inactive';
  createdAt: string;
  updatedAt: string;
}

// Schema
const %sSchema = z.object({
  code: z.string().min(1, 'Code is required'),
  name: z.string().min(1, 'Name is required'),
  description: z.string().optional(),
  status: z.enum(['active', 'inactive']),
});

type %sFormData = z.infer<typeof %sSchema>;

// API functions
const fetch%sList = async (): Promise<%s[]> => {
  const response = await api.get('/%s');
  return response.data;
};

const create%s = async (data: %sFormData): Promise<%s> => {
  const response = await api.post('/%s', data);
  return response.data;
};

const update%s = async ({ id, data }: { id: string; data: %sFormData }): Promise<%s> => {
  const response = await api.put('/%s/'+id, data);
  return response.data;
};

const delete%s = async (id: string): Promise<void> => {
  await api.delete('/%s/'+id);
};

export default function %sList() {
  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [editingItem, setEditingItem] = useState<%s | null>(null);
  const [searchTerm, setSearchTerm] = useState('');
  const queryClient = useQueryClient();

  // Queries
  const { data: items = [], isLoading } = useQuery({
    queryKey: ['%s'],
    queryFn: fetch%sList,
  });

  // Mutations
  const createMutation = useMutation({
    mutationFn: create%s,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['%s'] });
      toast.success('%s created successfully');
      setIsCreateOpen(false);
    },
    onError: () => {
      toast.error('Failed to create %s');
    },
  });

  const updateMutation = useMutation({
    mutationFn: update%s,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['%s'] });
      toast.success('%s updated successfully');
      setEditingItem(null);
    },
    onError: () => {
      toast.error('Failed to update %s');
    },
  });

  const deleteMutation = useMutation({
    mutationFn: delete%s,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['%s'] });
      toast.success('%s deleted successfully');
    },
    onError: () => {
      toast.error('Failed to delete %s');
    },
  });

  // Form
  const form = useForm<%sFormData>({
    resolver: zodResolver(%sSchema),
    defaultValues: {
      code: '',
      name: '',
      description: '',
      status: 'active',
    },
  });

  // Reset form when dialog opens/closes
  useEffect(() => {
    if (editingItem) {
      form.reset({
        code: editingItem.code,
        name: editingItem.name,
        description: editingItem.description || '',
        status: editingItem.status,
      });
    } else {
      form.reset();
    }
  }, [editingItem, form]);

  // Filter items based on search
  const filteredItems = items.filter(
    (item) =>
      item.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      item.code.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const handleSubmit = (data: %sFormData) => {
    if (editingItem) {
      updateMutation.mutate({ id: editingItem.id, data });
    } else {
      createMutation.mutate(data);
    }
  };

  const handleDelete = (id: string) => {
    if (confirm('Are you sure you want to delete this item?')) {
      deleteMutation.mutate(id);
    }
  };

  return (
    <div className="container mx-auto py-6">
      <Card>
        <CardHeader>
          <div className="flex justify-between items-center">
            <CardTitle>%s Management</CardTitle>
            <Button onClick={() => setIsCreateOpen(true)}>
              <Plus className="mr-2 h-4 w-4" />
              Add New
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          <div className="mb-4">
            <div className="relative">
              <Search className="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="Search by name or code..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="pl-8"
              />
            </div>
          </div>

          {isLoading ? (
            <div className="text-center py-4">Loading...</div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Code</TableHead>
                  <TableHead>Name</TableHead>
                  <TableHead>Description</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {filteredItems.map((item) => (
                  <TableRow key={item.id}>
                    <TableCell>{item.code}</TableCell>
                    <TableCell>{item.name}</TableCell>
                    <TableCell>{item.description}</TableCell>
                    <TableCell>
                      <span
                        className={` + "`" + `px-2 py-1 rounded-full text-xs ${
                          item.status === 'active'
                            ? 'bg-green-100 text-green-800'
                            : 'bg-gray-100 text-gray-800'
                        }` + "`" + `}
                      >
                        {item.status}
                      </span>
                    </TableCell>
                    <TableCell>
                      <div className="flex gap-2">
                        <Button
                          size="sm"
                          variant="outline"
                          onClick={() => setEditingItem(item)}
                        >
                          <Edit className="h-4 w-4" />
                        </Button>
                        <Button
                          size="sm"
                          variant="outline"
                          onClick={() => handleDelete(item.id)}
                        >
                          <Trash2 className="h-4 w-4" />
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>

      <Dialog
        open={isCreateOpen || !!editingItem}
        onOpenChange={(open) => {
          if (!open) {
            setIsCreateOpen(false);
            setEditingItem(null);
          }
        }}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>
              {editingItem ? 'Edit' : 'Create'} %s
            </DialogTitle>
          </DialogHeader>
          <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-4">
            <div>
              <Input
                {...form.register('code')}
                placeholder="Code"
                error={form.formState.errors.code?.message}
              />
            </div>
            <div>
              <Input
                {...form.register('name')}
                placeholder="Name"
                error={form.formState.errors.name?.message}
              />
            </div>
            <div>
              <Input
                {...form.register('description')}
                placeholder="Description (optional)"
              />
            </div>
            <div>
              <select
                {...form.register('status')}
                className="w-full p-2 border rounded"
              >
                <option value="active">Active</option>
                <option value="inactive">Inactive</option>
              </select>
            </div>
            <div className="flex justify-end gap-2">
              <Button
                type="button"
                variant="outline"
                onClick={() => {
                  setIsCreateOpen(false);
                  setEditingItem(null);
                }}
              >
                Cancel
              </Button>
              <Button type="submit">
                {editingItem ? 'Update' : 'Create'}
              </Button>
            </div>
          </form>
        </DialogContent>
      </Dialog>
    </div>
  );
}
`,
		moduleName,
		lowerName,
		moduleName,
		lowerName,
		moduleName,
		moduleName,
		lowerName,
		moduleName,
		moduleName,
		moduleName,
		lowerName,
		moduleName,
		moduleName,
		moduleName,
		lowerName,
		moduleName,
		lowerName,
		moduleName,
		moduleName,
		lowerName,
		moduleName,
		moduleName,
		lowerName,
		moduleName,
		lowerName,
		moduleName,
		moduleName,
		lowerName,
		moduleName,
		lowerName,
		moduleName,
		moduleName,
		lowerName,
		moduleName,
		lowerName,
		moduleName,
		moduleName,
		lowerName,
		moduleName,
		lowerName,
		moduleName,
		moduleName,
		lowerName,
		moduleName,
		moduleName,
		moduleName,
	)
}

// generateAPIHandler generates API handler code
func (a *EngineerAgent) generateAPIHandler(moduleName string) string {
	lowerName := strings.ToLower(moduleName)
	return fmt.Sprintf(`package handler

import (
	"net/http"

	"github.com/fastenmind/fastener-api/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// %sHandler handles %s-related requests
type %sHandler struct {
	service *service.%sService
}

// New%sHandler creates a new %s handler
func New%sHandler(service *service.%sService) *%sHandler {
	return &%sHandler{
		service: service,
	}
}

// List handles GET /%s
func (h *%sHandler) List(c echo.Context) error {
	ctx := c.Request().Context()
	
	// Parse query parameters
	params := service.List%sParams{
		Page:     getIntParam(c, "page", 1),
		PageSize: getIntParam(c, "page_size", 20),
		Search:   c.QueryParam("search"),
		Status:   c.QueryParam("status"),
		SortBy:   c.QueryParam("sort_by"),
		SortDir:  c.QueryParam("sort_dir"),
	}
	
	// Get company ID from context
	user := getUserFromContext(c)
	if user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}
	params.CompanyID = user.CompanyID
	
	// Call service
	result, err := h.service.List(ctx, params)
	if err != nil {
		return handleError(c, err)
	}
	
	return c.JSON(http.StatusOK, result)
}

// Get handles GET /%s/:id
func (h *%sHandler) Get(c echo.Context) error {
	ctx := c.Request().Context()
	
	// Parse ID
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}
	
	// Get company ID from context
	user := getUserFromContext(c)
	if user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}
	
	// Call service
	result, err := h.service.GetByID(ctx, id, user.CompanyID)
	if err != nil {
		return handleError(c, err)
	}
	
	return c.JSON(http.StatusOK, result)
}

// Create handles POST /%s
func (h *%sHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()
	
	// Parse request body
	var req service.Create%sRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	
	// Validate
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	// Get user from context
	user := getUserFromContext(c)
	if user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}
	req.CompanyID = user.CompanyID
	req.CreatedBy = user.ID
	
	// Call service
	result, err := h.service.Create(ctx, req)
	if err != nil {
		return handleError(c, err)
	}
	
	return c.JSON(http.StatusCreated, result)
}

// Update handles PUT /%s/:id
func (h *%sHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	
	// Parse ID
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}
	
	// Parse request body
	var req service.Update%sRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	
	// Validate
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	// Get user from context
	user := getUserFromContext(c)
	if user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}
	req.UpdatedBy = user.ID
	
	// Call service
	result, err := h.service.Update(ctx, id, user.CompanyID, req)
	if err != nil {
		return handleError(c, err)
	}
	
	return c.JSON(http.StatusOK, result)
}

// Delete handles DELETE /%s/:id
func (h *%sHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	
	// Parse ID
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}
	
	// Get user from context
	user := getUserFromContext(c)
	if user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}
	
	// Call service
	if err := h.service.Delete(ctx, id, user.CompanyID, user.ID); err != nil {
		return handleError(c, err)
	}
	
	return c.NoContent(http.StatusNoContent)
}

// RegisterRoutes registers all routes
func (h *%sHandler) RegisterRoutes(e *echo.Echo, middleware ...echo.MiddlewareFunc) {
	g := e.Group("/api/v1/%s", middleware...)
	
	g.GET("", h.List)
	g.GET("/:id", h.Get)
	g.POST("", h.Create)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
}
`,
		moduleName, lowerName,
		moduleName,
		moduleName,
		moduleName, lowerName,
		moduleName, moduleName,
		moduleName,
		moduleName,
		lowerName,
		moduleName,
		moduleName,
		lowerName,
		moduleName,
		lowerName,
		moduleName,
		moduleName,
		lowerName,
		moduleName,
		moduleName,
		lowerName,
		moduleName,
		moduleName,
		lowerName,
		moduleName,
		lowerName,
	)
}

// generateServiceCode generates service layer code
func (a *EngineerAgent) generateServiceCode(moduleName string) string {
	lowerName := strings.ToLower(moduleName)
	return fmt.Sprintf(`package service

import (
	"context"
	"fmt"

	"github.com/fastenmind/fastener-api/internal/domain/models"
	"github.com/fastenmind/fastener-api/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// %sService handles %s business logic
type %sService struct {
	repo *repository.%sRepository
	db   *gorm.DB
}

// New%sService creates a new %s service
func New%sService(repo *repository.%sRepository, db *gorm.DB) *%sService {
	return &%sService{
		repo: repo,
		db:   db,
	}
}

// List%sParams contains parameters for listing %ss
type List%sParams struct {
	CompanyID uuid.UUID
	Page      int
	PageSize  int
	Search    string
	Status    string
	SortBy    string
	SortDir   string
}

// List%sResult contains the list result
type List%sResult struct {
	Items      []models.%s ` + "`" + `json:"items"` + "`" + `
	Total      int64        ` + "`" + `json:"total"` + "`" + `
	Page       int          ` + "`" + `json:"page"` + "`" + `
	PageSize   int          ` + "`" + `json:"page_size"` + "`" + `
	TotalPages int          ` + "`" + `json:"total_pages"` + "`" + `
}

// List lists %ss with pagination
func (s *%sService) List(ctx context.Context, params List%sParams) (*List%sResult, error) {
	// Set defaults
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 20
	}
	
	// Call repository
	items, total, err := s.repo.List(ctx, repository.List%sParams{
		CompanyID: params.CompanyID,
		Offset:    (params.Page - 1) * params.PageSize,
		Limit:     params.PageSize,
		Search:    params.Search,
		Status:    params.Status,
		SortBy:    params.SortBy,
		SortDir:   params.SortDir,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list %ss: %%w", err)
	}
	
	// Calculate total pages
	totalPages := int(total) / params.PageSize
	if int(total)%%params.PageSize > 0 {
		totalPages++
	}
	
	return &List%sResult{
		Items:      items,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetByID gets a %s by ID
func (s *%sService) GetByID(ctx context.Context, id, companyID uuid.UUID) (*models.%s, error) {
	item, err := s.repo.GetByID(ctx, id, companyID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get %s: %%w", err)
	}
	
	return item, nil
}

// Create%sRequest contains data for creating a %s
type Create%sRequest struct {
	Code        string    ` + "`" + `json:"code" validate:"required,max=50"` + "`" + `
	Name        string    ` + "`" + `json:"name" validate:"required,max=100"` + "`" + `
	Description string    ` + "`" + `json:"description,omitempty" validate:"max=500"` + "`" + `
	Status      string    ` + "`" + `json:"status" validate:"required,oneof=active inactive"` + "`" + `
	CompanyID   uuid.UUID ` + "`" + `json:"-"` + "`" + `
	CreatedBy   uuid.UUID ` + "`" + `json:"-"` + "`" + `
}

// Create creates a new %s
func (s *%sService) Create(ctx context.Context, req Create%sRequest) (*models.%s, error) {
	// Check for duplicate code
	exists, err := s.repo.ExistsByCode(ctx, req.CompanyID, req.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to check duplicate: %%w", err)
	}
	if exists {
		return nil, ErrDuplicateCode
	}
	
	// Create model
	item := &models.%s{
		CompanyID:   req.CompanyID,
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
		CreatedBy:   &req.CreatedBy,
	}
	
	// Save to database
	if err := s.repo.Create(ctx, item); err != nil {
		return nil, fmt.Errorf("failed to create %s: %%w", err)
	}
	
	return item, nil
}

// Update%sRequest contains data for updating a %s
type Update%sRequest struct {
	Name        string    ` + "`" + `json:"name,omitempty" validate:"max=100"` + "`" + `
	Description string    ` + "`" + `json:"description,omitempty" validate:"max=500"` + "`" + `
	Status      string    ` + "`" + `json:"status,omitempty" validate:"omitempty,oneof=active inactive"` + "`" + `
	UpdatedBy   uuid.UUID ` + "`" + `json:"-"` + "`" + `
}

// Update updates a %s
func (s *%sService) Update(ctx context.Context, id, companyID uuid.UUID, req Update%sRequest) (*models.%s, error) {
	// Get existing item
	item, err := s.repo.GetByID(ctx, id, companyID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get %s: %%w", err)
	}
	
	// Update fields
	if req.Name != "" {
		item.Name = req.Name
	}
	if req.Description != "" {
		item.Description = req.Description
	}
	if req.Status != "" {
		item.Status = req.Status
	}
	item.UpdatedBy = &req.UpdatedBy
	
	// Save changes
	if err := s.repo.Update(ctx, item); err != nil {
		return nil, fmt.Errorf("failed to update %s: %%w", err)
	}
	
	return item, nil
}

// Delete deletes a %s
func (s *%sService) Delete(ctx context.Context, id, companyID, deletedBy uuid.UUID) error {
	// Check if exists
	_, err := s.repo.GetByID(ctx, id, companyID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrNotFound
		}
		return fmt.Errorf("failed to get %s: %%w", err)
	}
	
	// Soft delete
	if err := s.repo.Delete(ctx, id, deletedBy); err != nil {
		return fmt.Errorf("failed to delete %s: %%w", err)
	}
	
	return nil
}
`,
		moduleName, lowerName,
		moduleName,
		moduleName,
		moduleName, lowerName,
		moduleName, moduleName,
		moduleName,
		moduleName,
		moduleName, lowerName,
		moduleName,
		moduleName,
		moduleName,
		moduleName,
		lowerName,
		moduleName,
		moduleName, moduleName,
		moduleName,
		lowerName,
		moduleName,
		lowerName,
		moduleName, moduleName,
		lowerName,
		moduleName, lowerName,
		moduleName,
		moduleName,
		moduleName, moduleName,
		moduleName,
		lowerName,
		moduleName,
		moduleName,
		lowerName,
		moduleName, lowerName,
		moduleName,
		moduleName,
		moduleName, moduleName,
		lowerName,
		lowerName,
		lowerName,
		moduleName,
		lowerName,
		lowerName,
	)
}

// generateRepositoryCode generates repository layer code
func (a *EngineerAgent) generateRepositoryCode(moduleName string) string {
	lowerName := strings.ToLower(moduleName)
	tableName := lowerName + "s"
	
	return fmt.Sprintf(`package repository

import (
	"context"
	"fmt"

	"github.com/fastenmind/fastener-api/internal/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// %sRepository handles %s data access
type %sRepository struct {
	db *gorm.DB
}

// New%sRepository creates a new %s repository
func New%sRepository(db *gorm.DB) *%sRepository {
	return &%sRepository{
		db: db,
	}
}

// List%sParams contains parameters for listing %ss
type List%sParams struct {
	CompanyID uuid.UUID
	Offset    int
	Limit     int
	Search    string
	Status    string
	SortBy    string
	SortDir   string
}

// List lists %ss with filters
func (r *%sRepository) List(ctx context.Context, params List%sParams) ([]models.%s, int64, error) {
	var items []models.%s
	var total int64
	
	// Base query
	query := r.db.WithContext(ctx).Model(&models.%s{}).
		Where("company_id = ? AND deleted_at IS NULL", params.CompanyID)
	
	// Apply filters
	if params.Search != "" {
		search := "%%" + params.Search + "%%"
		query = query.Where("(name ILIKE ? OR code ILIKE ?)", search, search)
	}
	
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	
	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count %ss: %%w", err)
	}
	
	// Apply sorting
	if params.SortBy != "" {
		order := params.SortBy
		if params.SortDir == "desc" {
			order += " DESC"
		}
		query = query.Order(order)
	} else {
		query = query.Order("created_at DESC")
	}
	
	// Apply pagination
	query = query.Offset(params.Offset).Limit(params.Limit)
	
	// Execute query
	if err := query.Find(&items).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list %ss: %%w", err)
	}
	
	return items, total, nil
}

// GetByID gets a %s by ID
func (r *%sRepository) GetByID(ctx context.Context, id, companyID uuid.UUID) (*models.%s, error) {
	var item models.%s
	
	err := r.db.WithContext(ctx).
		Where("id = ? AND company_id = ? AND deleted_at IS NULL", id, companyID).
		First(&item).Error
		
	if err != nil {
		return nil, err
	}
	
	return &item, nil
}

// GetByCode gets a %s by code
func (r *%sRepository) GetByCode(ctx context.Context, companyID uuid.UUID, code string) (*models.%s, error) {
	var item models.%s
	
	err := r.db.WithContext(ctx).
		Where("company_id = ? AND code = ? AND deleted_at IS NULL", companyID, code).
		First(&item).Error
		
	if err != nil {
		return nil, err
	}
	
	return &item, nil
}

// ExistsByCode checks if a %s with the given code exists
func (r *%sRepository) ExistsByCode(ctx context.Context, companyID uuid.UUID, code string) (bool, error) {
	var count int64
	
	err := r.db.WithContext(ctx).Model(&models.%s{}).
		Where("company_id = ? AND code = ? AND deleted_at IS NULL", companyID, code).
		Count(&count).Error
		
	if err != nil {
		return false, err
	}
	
	return count > 0, nil
}

// Create creates a new %s
func (r *%sRepository) Create(ctx context.Context, item *models.%s) error {
	return r.db.WithContext(ctx).Create(item).Error
}

// Update updates a %s
func (r *%sRepository) Update(ctx context.Context, item *models.%s) error {
	return r.db.WithContext(ctx).Save(item).Error
}

// Delete soft deletes a %s
func (r *%sRepository) Delete(ctx context.Context, id, deletedBy uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.%s{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"deleted_at": gorm.DeletedAt{Time: time.Now(), Valid: true},
			"updated_by": deletedBy,
		}).Error
}

// BulkCreate creates multiple %ss
func (r *%sRepository) BulkCreate(ctx context.Context, items []models.%s) error {
	return r.db.WithContext(ctx).CreateInBatches(items, 100).Error
}

// GetStats gets statistics for %ss
func (r *%sRepository) GetStats(ctx context.Context, companyID uuid.UUID) (map[string]interface{}, error) {
	var stats struct {
		Total    int64 ` + "`" + `json:"total"` + "`" + `
		Active   int64 ` + "`" + `json:"active"` + "`" + `
		Inactive int64 ` + "`" + `json:"inactive"` + "`" + `
	}
	
	// Total count
	r.db.WithContext(ctx).Model(&models.%s{}).
		Where("company_id = ? AND deleted_at IS NULL", companyID).
		Count(&stats.Total)
	
	// Active count
	r.db.WithContext(ctx).Model(&models.%s{}).
		Where("company_id = ? AND status = ? AND deleted_at IS NULL", companyID, "active").
		Count(&stats.Active)
	
	// Inactive count
	r.db.WithContext(ctx).Model(&models.%s{}).
		Where("company_id = ? AND status = ? AND deleted_at IS NULL", companyID, "inactive").
		Count(&stats.Inactive)
	
	return map[string]interface{}{
		"total":    stats.Total,
		"active":   stats.Active,
		"inactive": stats.Inactive,
	}, nil
}
`,
		moduleName, lowerName,
		moduleName,
		moduleName, lowerName,
		moduleName, moduleName,
		moduleName,
		moduleName, lowerName,
		moduleName,
		lowerName,
		moduleName,
		moduleName, moduleName,
		moduleName,
		moduleName,
		lowerName,
		lowerName,
		lowerName,
		moduleName, moduleName,
		moduleName,
		lowerName,
		moduleName, moduleName,
		moduleName,
		lowerName,
		moduleName,
		lowerName,
		moduleName, moduleName,
		lowerName,
		moduleName,
		lowerName,
		moduleName,
		moduleName,
		lowerName,
		moduleName,
		moduleName,
		lowerName,
		moduleName,
		moduleName,
		moduleName,
		lowerName,
	)
}

// NewN8NAgent creates a new N8N automation agent
func NewN8NAgent() *N8NAgent {
	return &N8NAgent{
		name:        "N8N Automation Agent",
		description: "Configures N8N workflows for automated tasks and notifications",
	}
}

// N8NAgent implements the N8N automation AI agent
type N8NAgent struct {
	name        string
	description string
}

// GetType returns the agent type
func (a *N8NAgent) GetType() AgentType {
	return AgentTypeN8N
}

// GetName returns the agent name
func (a *N8NAgent) GetName() string {
	return a.name
}

// GetDescription returns the agent description
func (a *N8NAgent) GetDescription() string {
	return a.description
}

// Validate validates the input
func (a *N8NAgent) Validate(input AgentInput) error {
	if input.Task == "" {
		return fmt.Errorf("task is required")
	}
	return nil
}

// Execute runs the N8N agent task
func (a *N8NAgent) Execute(ctx context.Context, input AgentInput) (AgentOutput, error) {
	// For now, just generate N8N workflow configuration
	workflowConfig := a.generateN8NWorkflow(input.Task)
	
	files := []FileOutput{
		{
			Name:        "n8n_workflow.json",
			Path:        "/n8n/workflows/generated_workflow.json",
			Content:     workflowConfig,
			ContentType: "application/json",
		},
	}
	
	return AgentOutput{
		Success: true,
		Result: map[string]interface{}{
			"workflow_generated": true,
			"webhook_url":        "https://n8n.example.com/webhook/xxx",
		},
		Files:     files,
		NextSteps: []string{"Import workflow to N8N", "Configure webhook endpoints", "Test automation flow"},
		Metadata: map[string]interface{}{
			"generated_at": time.Now(),
		},
	}, nil
}

// generateN8NWorkflow generates N8N workflow configuration
func (a *N8NAgent) generateN8NWorkflow(task string) string {
	return `{
  "name": "Generated Workflow",
  "nodes": [
    {
      "parameters": {
        "path": "webhook-endpoint",
        "responseMode": "onReceived",
        "options": {}
      },
      "name": "Webhook",
      "type": "n8n-nodes-base.webhook",
      "typeVersion": 1,
      "position": [250, 300]
    }
  ],
  "connections": {},
  "active": false,
  "settings": {},
  "id": 1
}`
}