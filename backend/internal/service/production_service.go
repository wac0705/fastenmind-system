package service

import (
	"fmt"
	"time"
	"github.com/fastenmind/fastener-api/internal/models"
	"github.com/fastenmind/fastener-api/internal/repository"
	"github.com/google/uuid"
)

type ProductionService interface {
	// Production Order operations
	CreateProductionOrder(order *models.ProductionOrder) error
	UpdateProductionOrder(order *models.ProductionOrder) error
	GetProductionOrder(id uuid.UUID) (*models.ProductionOrder, error)
	ListProductionOrders(companyID uuid.UUID, params map[string]interface{}) ([]models.ProductionOrder, int64, error)
	ReleaseProductionOrder(id uuid.UUID, userID uuid.UUID) error
	StartProductionOrder(id uuid.UUID, userID uuid.UUID) error
	CompleteProductionOrder(id uuid.UUID, userID uuid.UUID) error
	CancelProductionOrder(id uuid.UUID, userID uuid.UUID, reason string) error
	
	// Production Route operations
	CreateProductionRoute(route *models.ProductionRoute) error
	UpdateProductionRoute(route *models.ProductionRoute) error
	GetProductionRoute(id uuid.UUID) (*models.ProductionRoute, error)
	ListProductionRoutes(companyID uuid.UUID, params map[string]interface{}) ([]models.ProductionRoute, error)
	
	// Route Operation operations
	CreateRouteOperation(operation *models.RouteOperation) error
	UpdateRouteOperation(operation *models.RouteOperation) error
	GetRouteOperations(routeID uuid.UUID) ([]models.RouteOperation, error)
	DeleteRouteOperation(id uuid.UUID) error
	
	// Work Station operations
	CreateWorkStation(station *models.WorkStation) error
	UpdateWorkStation(station *models.WorkStation) error
	GetWorkStation(id uuid.UUID) (*models.WorkStation, error)
	ListWorkStations(companyID uuid.UUID, params map[string]interface{}) ([]models.WorkStation, error)
	
	// Production Task operations
	CreateProductionTask(task *models.ProductionTask) error
	UpdateProductionTask(task *models.ProductionTask) error
	GetProductionTask(id uuid.UUID) (*models.ProductionTask, error)
	ListProductionTasks(params map[string]interface{}) ([]models.ProductionTask, error)
	AssignTask(taskID uuid.UUID, userID uuid.UUID, assignedBy uuid.UUID) error
	StartTask(taskID uuid.UUID, userID uuid.UUID) error
	CompleteTask(taskID uuid.UUID, userID uuid.UUID, completedQuantity float64, qualifiedQuantity float64, notes string) error
	
	// Production Material operations
	CreateProductionMaterial(material *models.ProductionMaterial) error
	IssueMaterials(productionOrderID uuid.UUID, userID uuid.UUID) error
	GetProductionMaterials(productionOrderID uuid.UUID) ([]models.ProductionMaterial, error)
	
	// Quality Inspection operations
	CreateQualityInspection(inspection *models.QualityInspection) error
	UpdateQualityInspection(inspection *models.QualityInspection) error
	GetQualityInspection(id uuid.UUID) (*models.QualityInspection, error)
	ListQualityInspections(companyID uuid.UUID, params map[string]interface{}) ([]models.QualityInspection, int64, error)
	ApproveInspection(id uuid.UUID, approverID uuid.UUID) error
	RejectInspection(id uuid.UUID, approverID uuid.UUID, reason string) error
	
	// Dashboard and Reports
	GetProductionDashboard(companyID uuid.UUID) (*ProductionDashboard, error)
	GetProductionStats(companyID uuid.UUID) (*ProductionStats, error)
}

type productionService struct {
	productionRepo repository.ProductionRepository
	inventoryRepo  repository.InventoryRepository
	orderRepo      repository.OrderRepository
}

func NewProductionService(
	productionRepo repository.ProductionRepository,
	inventoryRepo repository.InventoryRepository,
	orderRepo repository.OrderRepository,
) ProductionService {
	return &productionService{
		productionRepo: productionRepo,
		inventoryRepo:  inventoryRepo,
		orderRepo:      orderRepo,
	}
}

// Production Order operations
func (s *productionService) CreateProductionOrder(order *models.ProductionOrder) error {
	// Generate order number
	order.OrderNo = s.generateOrderNo(order.CompanyID)
	
	// Set initial status
	if order.Status == "" {
		order.Status = "planned"
	}
	
	// Set default priority
	if order.Priority == "" {
		order.Priority = "medium"
	}
	
	// Initialize quantities
	order.ProducedQuantity = 0
	order.QualifiedQuantity = 0
	order.DefectQuantity = 0
	
	// Calculate estimated cost if route is provided
	if order.RouteID != nil {
		route, err := s.productionRepo.GetProductionRoute(*order.RouteID)
		if err == nil {
			order.EstimatedCost = route.EstimatedCost * order.PlannedQuantity
		}
	}
	
	if err := s.productionRepo.CreateProductionOrder(order); err != nil {
		return err
	}
	
	// Create production tasks if route is provided
	if order.RouteID != nil {
		if err := s.createProductionTasks(order); err != nil {
			return err
		}
	}
	
	// Create material requirements
	if err := s.createMaterialRequirements(order); err != nil {
		return err
	}
	
	return nil
}

func (s *productionService) UpdateProductionOrder(order *models.ProductionOrder) error {
	return s.productionRepo.UpdateProductionOrder(order)
}

func (s *productionService) GetProductionOrder(id uuid.UUID) (*models.ProductionOrder, error) {
	return s.productionRepo.GetProductionOrder(id)
}

func (s *productionService) ListProductionOrders(companyID uuid.UUID, params map[string]interface{}) ([]models.ProductionOrder, int64, error) {
	return s.productionRepo.ListProductionOrders(companyID, params)
}

func (s *productionService) ReleaseProductionOrder(id uuid.UUID, userID uuid.UUID) error {
	order, err := s.productionRepo.GetProductionOrder(id)
	if err != nil {
		return err
	}
	
	if order.Status != "planned" {
		return fmt.Errorf("production order must be in planned status to release")
	}
	
	// Check material availability
	materials, err := s.productionRepo.GetProductionMaterials(id)
	if err != nil {
		return err
	}
	
	for _, material := range materials {
		inventory, err := s.inventoryRepo.Get(material.InventoryID)
		if err != nil {
			return err
		}
		
		if inventory.AvailableStock < material.PlannedQuantity {
			return fmt.Errorf("insufficient material: %s (required: %.2f, available: %.2f)",
				inventory.Name, material.PlannedQuantity, inventory.AvailableStock)
		}
	}
	
	order.Status = "released"
	return s.productionRepo.UpdateProductionOrder(order)
}

func (s *productionService) StartProductionOrder(id uuid.UUID, userID uuid.UUID) error {
	order, err := s.productionRepo.GetProductionOrder(id)
	if err != nil {
		return err
	}
	
	if order.Status != "released" {
		return fmt.Errorf("production order must be released to start")
	}
	
	now := time.Now()
	order.Status = "in_progress"
	order.ActualStartDate = &now
	
	// Issue materials
	if err := s.IssueMaterials(id, userID); err != nil {
		return err
	}
	
	return s.productionRepo.UpdateProductionOrder(order)
}

func (s *productionService) CompleteProductionOrder(id uuid.UUID, userID uuid.UUID) error {
	order, err := s.productionRepo.GetProductionOrder(id)
	if err != nil {
		return err
	}
	
	if order.Status != "in_progress" && order.Status != "quality_check" {
		return fmt.Errorf("production order must be in progress or quality check to complete")
	}
	
	// Check if all tasks are completed
	tasks, err := s.productionRepo.GetTasksByProductionOrder(id)
	if err != nil {
		return err
	}
	
	for _, task := range tasks {
		if task.Status != "completed" {
			return fmt.Errorf("all tasks must be completed before completing production order")
		}
	}
	
	now := time.Now()
	order.Status = "completed"
	order.ActualEndDate = &now
	
	// Update inventory with produced quantity
	if order.QualifiedQuantity > 0 {
		if err := s.updateInventoryAfterProduction(order); err != nil {
			return err
		}
	}
	
	return s.productionRepo.UpdateProductionOrder(order)
}

func (s *productionService) CancelProductionOrder(id uuid.UUID, userID uuid.UUID, reason string) error {
	order, err := s.productionRepo.GetProductionOrder(id)
	if err != nil {
		return err
	}
	
	if order.Status == "completed" || order.Status == "cancelled" {
		return fmt.Errorf("cannot cancel completed or already cancelled production order")
	}
	
	order.Status = "cancelled"
	order.Notes = fmt.Sprintf("Cancelled: %s", reason)
	
	// Return issued materials to inventory
	if err := s.returnIssuedMaterials(id); err != nil {
		return err
	}
	
	return s.productionRepo.UpdateProductionOrder(order)
}

// Production Route operations
func (s *productionService) CreateProductionRoute(route *models.ProductionRoute) error {
	// Generate route number
	route.RouteNo = s.generateRouteNo(route.CompanyID)
	
	return s.productionRepo.CreateProductionRoute(route)
}

func (s *productionService) UpdateProductionRoute(route *models.ProductionRoute) error {
	return s.productionRepo.UpdateProductionRoute(route)
}

func (s *productionService) GetProductionRoute(id uuid.UUID) (*models.ProductionRoute, error) {
	return s.productionRepo.GetProductionRoute(id)
}

func (s *productionService) ListProductionRoutes(companyID uuid.UUID, params map[string]interface{}) ([]models.ProductionRoute, error) {
	return s.productionRepo.ListProductionRoutes(companyID, params)
}

// Route Operation operations
func (s *productionService) CreateRouteOperation(operation *models.RouteOperation) error {
	return s.productionRepo.CreateRouteOperation(operation)
}

func (s *productionService) UpdateRouteOperation(operation *models.RouteOperation) error {
	return s.productionRepo.UpdateRouteOperation(operation)
}

func (s *productionService) GetRouteOperations(routeID uuid.UUID) ([]models.RouteOperation, error) {
	return s.productionRepo.GetRouteOperations(routeID)
}

func (s *productionService) DeleteRouteOperation(id uuid.UUID) error {
	return s.productionRepo.DeleteRouteOperation(id)
}

// Work Station operations
func (s *productionService) CreateWorkStation(station *models.WorkStation) error {
	// Generate station number
	station.StationNo = s.generateStationNo(station.CompanyID)
	
	// Set default status
	if station.Status == "" {
		station.Status = "available"
	}
	
	return s.productionRepo.CreateWorkStation(station)
}

func (s *productionService) UpdateWorkStation(station *models.WorkStation) error {
	return s.productionRepo.UpdateWorkStation(station)
}

func (s *productionService) GetWorkStation(id uuid.UUID) (*models.WorkStation, error) {
	return s.productionRepo.GetWorkStation(id)
}

func (s *productionService) ListWorkStations(companyID uuid.UUID, params map[string]interface{}) ([]models.WorkStation, error) {
	return s.productionRepo.ListWorkStations(companyID, params)
}

// Production Task operations
func (s *productionService) CreateProductionTask(task *models.ProductionTask) error {
	return s.productionRepo.CreateProductionTask(task)
}

func (s *productionService) UpdateProductionTask(task *models.ProductionTask) error {
	return s.productionRepo.UpdateProductionTask(task)
}

func (s *productionService) GetProductionTask(id uuid.UUID) (*models.ProductionTask, error) {
	return s.productionRepo.GetProductionTask(id)
}

func (s *productionService) ListProductionTasks(params map[string]interface{}) ([]models.ProductionTask, error) {
	return s.productionRepo.ListProductionTasks(params)
}

func (s *productionService) AssignTask(taskID uuid.UUID, userID uuid.UUID, assignedBy uuid.UUID) error {
	task, err := s.productionRepo.GetProductionTask(taskID)
	if err != nil {
		return err
	}
	
	if task.Status != "pending" {
		return fmt.Errorf("task must be pending to assign")
	}
	
	now := time.Now()
	task.AssignedTo = &userID
	task.AssignedAt = &now
	task.Status = "assigned"
	
	return s.productionRepo.UpdateProductionTask(task)
}

func (s *productionService) StartTask(taskID uuid.UUID, userID uuid.UUID) error {
	task, err := s.productionRepo.GetProductionTask(taskID)
	if err != nil {
		return err
	}
	
	if task.Status != "assigned" && task.Status != "pending" {
		return fmt.Errorf("task must be assigned or pending to start")
	}
	
	if task.AssignedTo != nil && *task.AssignedTo != userID {
		return fmt.Errorf("task is assigned to another user")
	}
	
	now := time.Now()
	task.Status = "in_progress"
	task.ActualStartTime = &now
	task.AssignedTo = &userID
	
	return s.productionRepo.UpdateProductionTask(task)
}

func (s *productionService) CompleteTask(taskID uuid.UUID, userID uuid.UUID, completedQuantity float64, qualifiedQuantity float64, notes string) error {
	task, err := s.productionRepo.GetProductionTask(taskID)
	if err != nil {
		return err
	}
	
	if task.Status != "in_progress" {
		return fmt.Errorf("task must be in progress to complete")
	}
	
	if task.AssignedTo == nil || *task.AssignedTo != userID {
		return fmt.Errorf("task is not assigned to this user")
	}
	
	now := time.Now()
	task.Status = "completed"
	task.ActualEndTime = &now
	task.CompletedQuantity = completedQuantity
	task.QualifiedQuantity = qualifiedQuantity
	task.DefectQuantity = completedQuantity - qualifiedQuantity
	task.Notes = notes
	
	// Set QC status based on operation requirements
	if task.RouteOperation != nil {
		operation, _ := s.productionRepo.GetRouteOperations(task.RouteOperationID)
		if len(operation) > 0 && operation[0].QCRequired {
			task.QCStatus = "pending"
		} else {
			task.QCStatus = "not_required"
		}
	}
	
	if err := s.productionRepo.UpdateProductionTask(task); err != nil {
		return err
	}
	
	// Update production order progress
	return s.updateProductionOrderProgress(task.ProductionOrderID)
}

// Production Material operations
func (s *productionService) CreateProductionMaterial(material *models.ProductionMaterial) error {
	return s.productionRepo.CreateProductionMaterial(material)
}

func (s *productionService) IssueMaterials(productionOrderID uuid.UUID, userID uuid.UUID) error {
	materials, err := s.productionRepo.GetProductionMaterials(productionOrderID)
	if err != nil {
		return err
	}
	
	now := time.Now()
	
	for _, material := range materials {
		if material.Status == "planned" {
			// Check inventory availability
			inventory, err := s.inventoryRepo.Get(material.InventoryID)
			if err != nil {
				return err
			}
			
			if inventory.AvailableStock < material.PlannedQuantity {
				return fmt.Errorf("insufficient inventory for %s", inventory.Name)
			}
			
			// Reserve inventory
			inventory.ReservedStock += material.PlannedQuantity
			inventory.AvailableStock -= material.PlannedQuantity
			
			if err := s.inventoryRepo.Update(inventory); err != nil {
				return err
			}
			
			// Update material status
			material.Status = "issued"
			material.IssuedQuantity = material.PlannedQuantity
			material.IssuedAt = &now
			
			if err := s.productionRepo.UpdateProductionMaterial(&material); err != nil {
				return err
			}
		}
	}
	
	return nil
}

func (s *productionService) GetProductionMaterials(productionOrderID uuid.UUID) ([]models.ProductionMaterial, error) {
	return s.productionRepo.GetProductionMaterials(productionOrderID)
}

// Quality Inspection operations
func (s *productionService) CreateQualityInspection(inspection *models.QualityInspection) error {
	// Generate inspection number
	inspection.InspectionNo = s.generateInspectionNo(inspection.CompanyID)
	
	if inspection.Status == "" {
		inspection.Status = "pending"
	}
	
	return s.productionRepo.CreateQualityInspection(inspection)
}

func (s *productionService) UpdateQualityInspection(inspection *models.QualityInspection) error {
	return s.productionRepo.UpdateQualityInspection(inspection)
}

func (s *productionService) GetQualityInspection(id uuid.UUID) (*models.QualityInspection, error) {
	return s.productionRepo.GetQualityInspection(id)
}

func (s *productionService) ListQualityInspections(companyID uuid.UUID, params map[string]interface{}) ([]models.QualityInspection, int64, error) {
	return s.productionRepo.ListQualityInspections(companyID, params)
}

func (s *productionService) ApproveInspection(id uuid.UUID, approverID uuid.UUID) error {
	inspection, err := s.productionRepo.GetQualityInspection(id)
	if err != nil {
		return err
	}
	
	if inspection.Status != "in_progress" {
		return fmt.Errorf("inspection must be in progress to approve")
	}
	
	now := time.Now()
	inspection.Status = "passed"
	inspection.ApprovedBy = &approverID
	inspection.ApprovedAt = &now
	
	return s.productionRepo.UpdateQualityInspection(inspection)
}

func (s *productionService) RejectInspection(id uuid.UUID, approverID uuid.UUID, reason string) error {
	inspection, err := s.productionRepo.GetQualityInspection(id)
	if err != nil {
		return err
	}
	
	if inspection.Status != "in_progress" {
		return fmt.Errorf("inspection must be in progress to reject")
	}
	
	now := time.Now()
	inspection.Status = "failed"
	inspection.ApprovedBy = &approverID
	inspection.ApprovedAt = &now
	inspection.CorrectiveAction = reason
	
	return s.productionRepo.UpdateQualityInspection(inspection)
}

// Dashboard and Reports
func (s *productionService) GetProductionDashboard(companyID uuid.UUID) (*ProductionDashboard, error) {
	dashboard := &ProductionDashboard{}
	
	// Get production orders statistics
	params := map[string]interface{}{}
	orders, _, _ := s.productionRepo.ListProductionOrders(companyID, params)
	
	for _, order := range orders {
		dashboard.TotalOrders++
		
		switch order.Status {
		case "planned":
			dashboard.PlannedOrders++
		case "released":
			dashboard.ReleasedOrders++
		case "in_progress":
			dashboard.InProgressOrders++
		case "completed":
			dashboard.CompletedOrders++
		case "cancelled":
			dashboard.CancelledOrders++
		}
		
		dashboard.TotalPlannedQuantity += order.PlannedQuantity
		dashboard.TotalProducedQuantity += order.ProducedQuantity
		dashboard.TotalQualifiedQuantity += order.QualifiedQuantity
		dashboard.TotalDefectQuantity += order.DefectQuantity
	}
	
	// Calculate efficiency
	if dashboard.TotalProducedQuantity > 0 {
		dashboard.QualityRate = (dashboard.TotalQualifiedQuantity / dashboard.TotalProducedQuantity) * 100
	}
	
	if dashboard.TotalPlannedQuantity > 0 {
		dashboard.ProductionEfficiency = (dashboard.TotalProducedQuantity / dashboard.TotalPlannedQuantity) * 100
	}
	
	return dashboard, nil
}

func (s *productionService) GetProductionStats(companyID uuid.UUID) (*ProductionStats, error) {
	stats := &ProductionStats{}
	
	// Get work stations count
	stations, _ := s.productionRepo.ListWorkStations(companyID, map[string]interface{}{})
	stats.TotalWorkStations = len(stations)
	
	for _, station := range stations {
		switch station.Status {
		case "available":
			stats.AvailableStations++
		case "busy":
			stats.BusyStations++
		case "maintenance":
			stats.MaintenanceStations++
		case "breakdown":
			stats.BreakdownStations++
		}
	}
	
	// Get quality inspections count
	inspections, total, _ := s.productionRepo.ListQualityInspections(companyID, map[string]interface{}{})
	stats.TotalInspections = int(total)
	
	for _, inspection := range inspections {
		switch inspection.Status {
		case "passed":
			stats.PassedInspections++
		case "failed":
			stats.FailedInspections++
		case "pending":
			stats.PendingInspections++
		}
	}
	
	return stats, nil
}

// Helper methods
func (s *productionService) generateOrderNo(companyID uuid.UUID) string {
	return fmt.Sprintf("PO-%s-%06d", time.Now().Format("200601"), time.Now().Unix()%1000000)
}

func (s *productionService) generateRouteNo(companyID uuid.UUID) string {
	return fmt.Sprintf("RT-%s-%06d", time.Now().Format("200601"), time.Now().Unix()%1000000)
}

func (s *productionService) generateStationNo(companyID uuid.UUID) string {
	return fmt.Sprintf("WS-%06d", time.Now().Unix()%1000000)
}

func (s *productionService) generateInspectionNo(companyID uuid.UUID) string {
	return fmt.Sprintf("QI-%s-%06d", time.Now().Format("200601"), time.Now().Unix()%1000000)
}

func (s *productionService) createProductionTasks(order *models.ProductionOrder) error {
	if order.RouteID == nil {
		return nil
	}
	
	operations, err := s.productionRepo.GetRouteOperations(*order.RouteID)
	if err != nil {
		return err
	}
	
	for i, operation := range operations {
		task := &models.ProductionTask{
			ProductionOrderID: order.ID,
			RouteOperationID:  operation.ID,
			WorkStationID:     operation.WorkStationID,
			TaskNo:            i + 1,
			Name:              operation.Name,
			Status:            "pending",
			PlannedQuantity:   order.PlannedQuantity,
			PlannedStartTime:  order.PlannedStartDate.Add(time.Duration(i*8) * time.Hour), // 8 hours per task
			PlannedEndTime:    order.PlannedStartDate.Add(time.Duration((i+1)*8) * time.Hour),
		}
		
		if operation.QCRequired {
			task.QCStatus = "not_required" // Will be set to pending when task is completed
		} else {
			task.QCStatus = "not_required"
		}
		
		if err := s.productionRepo.CreateProductionTask(task); err != nil {
			return err
		}
	}
	
	return nil
}

func (s *productionService) createMaterialRequirements(order *models.ProductionOrder) error {
	// This would typically come from a BOM (Bill of Materials)
	// For now, we'll create a placeholder implementation
	// In a real system, you would query a BOM table based on the inventory item
	
	return nil
}

func (s *productionService) updateProductionOrderProgress(productionOrderID uuid.UUID) error {
	order, err := s.productionRepo.GetProductionOrder(productionOrderID)
	if err != nil {
		return err
	}
	
	tasks, err := s.productionRepo.GetTasksByProductionOrder(productionOrderID)
	if err != nil {
		return err
	}
	
	completedTasks := 0
	totalProduced := 0.0
	totalQualified := 0.0
	totalDefects := 0.0
	
	for _, task := range tasks {
		if task.Status == "completed" {
			completedTasks++
			totalProduced += task.CompletedQuantity
			totalQualified += task.QualifiedQuantity
			totalDefects += task.DefectQuantity
		}
	}
	
	order.CompletedStations = completedTasks
	order.TotalStations = len(tasks)
	order.ProducedQuantity = totalProduced / float64(len(tasks)) // Average across tasks
	order.QualifiedQuantity = totalQualified / float64(len(tasks))
	order.DefectQuantity = totalDefects / float64(len(tasks))
	
	// Update status based on progress
	if completedTasks == len(tasks) && len(tasks) > 0 {
		order.Status = "quality_check" // Ready for final QC
	}
	
	return s.productionRepo.UpdateProductionOrder(order)
}

func (s *productionService) updateInventoryAfterProduction(order *models.ProductionOrder) error {
	inventory, err := s.inventoryRepo.Get(order.InventoryID)
	if err != nil {
		return err
	}
	
	// Add qualified quantity to inventory
	inventory.CurrentStock += order.QualifiedQuantity
	inventory.AvailableStock += order.QualifiedQuantity
	
	return s.inventoryRepo.Update(inventory)
}

func (s *productionService) returnIssuedMaterials(productionOrderID uuid.UUID) error {
	materials, err := s.productionRepo.GetProductionMaterials(productionOrderID)
	if err != nil {
		return err
	}
	
	for _, material := range materials {
		if material.Status == "issued" {
			inventory, err := s.inventoryRepo.Get(material.InventoryID)
			if err != nil {
				continue
			}
			
			// Return unused materials
			unusedQuantity := material.IssuedQuantity - material.ConsumedQuantity
			if unusedQuantity > 0 {
				inventory.ReservedStock -= unusedQuantity
				inventory.AvailableStock += unusedQuantity
				
				s.inventoryRepo.Update(inventory)
			}
			
			material.Status = "returned"
			material.ReturnedQuantity = unusedQuantity
			s.productionRepo.UpdateProductionMaterial(&material)
		}
	}
	
	return nil
}

// Report structs
type ProductionDashboard struct {
	TotalOrders              int     `json:"total_orders"`
	PlannedOrders            int     `json:"planned_orders"`
	ReleasedOrders           int     `json:"released_orders"`
	InProgressOrders         int     `json:"in_progress_orders"`
	CompletedOrders          int     `json:"completed_orders"`
	CancelledOrders          int     `json:"cancelled_orders"`
	TotalPlannedQuantity     float64 `json:"total_planned_quantity"`
	TotalProducedQuantity    float64 `json:"total_produced_quantity"`
	TotalQualifiedQuantity   float64 `json:"total_qualified_quantity"`
	TotalDefectQuantity      float64 `json:"total_defect_quantity"`
	ProductionEfficiency     float64 `json:"production_efficiency"`
	QualityRate              float64 `json:"quality_rate"`
}

type ProductionStats struct {
	TotalWorkStations    int `json:"total_work_stations"`
	AvailableStations    int `json:"available_stations"`
	BusyStations         int `json:"busy_stations"`
	MaintenanceStations  int `json:"maintenance_stations"`
	BreakdownStations    int `json:"breakdown_stations"`
	TotalInspections     int `json:"total_inspections"`
	PassedInspections    int `json:"passed_inspections"`
	FailedInspections    int `json:"failed_inspections"`
	PendingInspections   int `json:"pending_inspections"`
}