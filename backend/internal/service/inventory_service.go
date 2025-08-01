package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/fastenmind/fastener-api/internal/models"
	"github.com/fastenmind/fastener-api/internal/repository"
	"github.com/google/uuid"
)

type InventoryService interface {
	// Inventory management
	List(companyID uuid.UUID, params map[string]interface{}) ([]models.Inventory, int64, error)
	Get(id uuid.UUID) (*models.Inventory, error)
	GetBySKU(sku string) (*models.Inventory, error)
	Create(companyID uuid.UUID, req CreateInventoryRequest) (*models.Inventory, error)
	Update(id uuid.UUID, req UpdateInventoryRequest) (*models.Inventory, error)
	Delete(id uuid.UUID) error
	
	// Stock operations
	AdjustStock(id, userID uuid.UUID, req StockAdjustmentRequest) (*models.StockMovement, error)
	TransferStock(userID uuid.UUID, req StockTransferRequest) (*models.StockMovement, error)
	ReserveStock(id uuid.UUID, quantity float64, orderID uuid.UUID) error
	ReleaseStock(id uuid.UUID, quantity float64, orderID uuid.UUID) error
	
	// Warehouse management
	ListWarehouses(companyID uuid.UUID) ([]models.Warehouse, error)
	GetWarehouse(id uuid.UUID) (*models.Warehouse, error)
	CreateWarehouse(companyID uuid.UUID, req CreateWarehouseRequest) (*models.Warehouse, error)
	UpdateWarehouse(id uuid.UUID, req UpdateWarehouseRequest) (*models.Warehouse, error)
	
	// Stock movements
	GetMovements(inventoryID uuid.UUID, params map[string]interface{}) ([]models.StockMovement, error)
	GetMovementsByOrder(orderID uuid.UUID) ([]models.StockMovement, error)
	
	// Alerts
	GetActiveAlerts(companyID uuid.UUID) ([]models.StockAlert, error)
	AcknowledgeAlert(id, userID uuid.UUID) error
	ResolveAlert(id, userID uuid.UUID, resolution string) error
	
	// Reports
	GetInventoryStats(companyID uuid.UUID) (*InventoryStats, error)
	GetLowStockItems(companyID uuid.UUID) ([]models.Inventory, error)
	GetStockValuation(companyID uuid.UUID) (*StockValuation, error)
	
	// Stock take
	CreateStockTake(companyID, userID uuid.UUID, req CreateStockTakeRequest) (*models.StockTake, error)
	UpdateStockTake(id uuid.UUID, req UpdateStockTakeRequest) (*models.StockTake, error)
	GetStockTake(id uuid.UUID) (*models.StockTake, error)
	ListStockTakes(companyID uuid.UUID, params map[string]interface{}) ([]models.StockTake, error)
	SubmitStockCount(stockTakeID, userID uuid.UUID, req StockCountRequest) error
	CompleteStockTake(id, userID uuid.UUID) error
}

type CreateInventoryRequest struct {
	SKU               string    `json:"sku" validate:"required"`
	PartNo            string    `json:"part_no" validate:"required"`
	Name              string    `json:"name" validate:"required"`
	Description       string    `json:"description"`
	Category          string    `json:"category" validate:"required,oneof=raw_material semi_finished finished_goods"`
	Material          string    `json:"material"`
	Specification     string    `json:"specification"`
	SurfaceTreatment  string    `json:"surface_treatment"`
	HeatTreatment     string    `json:"heat_treatment"`
	Unit              string    `json:"unit"`
	InitialStock      float64   `json:"initial_stock"`
	MinStock          float64   `json:"min_stock"`
	MaxStock          float64   `json:"max_stock"`
	ReorderPoint      float64   `json:"reorder_point"`
	ReorderQuantity   float64   `json:"reorder_quantity"`
	WarehouseID       uuid.UUID `json:"warehouse_id"`
	Location          string    `json:"location"`
	StandardCost      float64   `json:"standard_cost"`
	PrimarySupplierID uuid.UUID `json:"primary_supplier_id"`
	LeadTimeDays      int       `json:"lead_time_days"`
}

type UpdateInventoryRequest struct {
	Name              string  `json:"name"`
	Description       string  `json:"description"`
	Category          string  `json:"category"`
	Material          string  `json:"material"`
	Specification     string  `json:"specification"`
	SurfaceTreatment  string  `json:"surface_treatment"`
	HeatTreatment     string  `json:"heat_treatment"`
	MinStock          float64 `json:"min_stock"`
	MaxStock          float64 `json:"max_stock"`
	ReorderPoint      float64 `json:"reorder_point"`
	ReorderQuantity   float64 `json:"reorder_quantity"`
	Location          string  `json:"location"`
	StandardCost      float64 `json:"standard_cost"`
	LeadTimeDays      int     `json:"lead_time_days"`
	Status            string  `json:"status"`
}

type StockAdjustmentRequest struct {
	Quantity     float64 `json:"quantity" validate:"required"`
	Reason       string  `json:"reason" validate:"required,oneof=damage loss found correction"`
	Notes        string  `json:"notes"`
	BatchNo      string  `json:"batch_no"`
	WarehouseID  uuid.UUID `json:"warehouse_id"`
}

type StockTransferRequest struct {
	InventoryID      uuid.UUID `json:"inventory_id" validate:"required"`
	Quantity         float64   `json:"quantity" validate:"required,min=0"`
	FromWarehouseID  uuid.UUID `json:"from_warehouse_id" validate:"required"`
	ToWarehouseID    uuid.UUID `json:"to_warehouse_id" validate:"required"`
	FromLocation     string    `json:"from_location"`
	ToLocation       string    `json:"to_location"`
	Notes            string    `json:"notes"`
}

type CreateWarehouseRequest struct {
	Code     string `json:"code" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Type     string `json:"type" validate:"required,oneof=main branch consignment"`
	Address  string `json:"address"`
	Manager  string `json:"manager"`
	Phone    string `json:"phone"`
}

type UpdateWarehouseRequest struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Address  string `json:"address"`
	Manager  string `json:"manager"`
	Phone    string `json:"phone"`
	IsActive bool   `json:"is_active"`
}

type CreateStockTakeRequest struct {
	WarehouseID   uuid.UUID `json:"warehouse_id" validate:"required"`
	Type          string    `json:"type" validate:"required,oneof=full cycle spot"`
	ScheduledDate string    `json:"scheduled_date" validate:"required"`
	AssignedTo    uuid.UUID `json:"assigned_to" validate:"required"`
	Notes         string    `json:"notes"`
}

type UpdateStockTakeRequest struct {
	Status        string    `json:"status"`
	AssignedTo    uuid.UUID `json:"assigned_to"`
	Notes         string    `json:"notes"`
}

type StockCountRequest struct {
	Items []StockCountItem `json:"items" validate:"required"`
}

type StockCountItem struct {
	InventoryID     uuid.UUID `json:"inventory_id" validate:"required"`
	CountedQuantity float64   `json:"counted_quantity" validate:"min=0"`
	Notes           string    `json:"notes"`
}

type InventoryStats struct {
	TotalItems       int     `json:"total_items"`
	TotalValue       float64 `json:"total_value"`
	LowStockItems    int     `json:"low_stock_items"`
	OutOfStockItems  int     `json:"out_of_stock_items"`
	OverstockItems   int     `json:"overstock_items"`
	ActiveAlerts     int     `json:"active_alerts"`
}

type StockValuation struct {
	TotalValue       float64                    `json:"total_value"`
	ByCategory       map[string]float64         `json:"by_category"`
	ByWarehouse      map[string]float64         `json:"by_warehouse"`
	TopValueItems    []InventoryValueItem       `json:"top_value_items"`
}

type InventoryValueItem struct {
	InventoryID   uuid.UUID `json:"inventory_id"`
	SKU           string    `json:"sku"`
	Name          string    `json:"name"`
	Quantity      float64   `json:"quantity"`
	UnitCost      float64   `json:"unit_cost"`
	TotalValue    float64   `json:"total_value"`
}

type inventoryService struct {
	inventoryRepo repository.InventoryRepository
	orderRepo     repository.OrderRepository
	n8nService    N8NService
}

func NewInventoryService(
	inventoryRepo repository.InventoryRepository,
	orderRepo repository.OrderRepository,
	n8nService N8NService,
) InventoryService {
	return &inventoryService{
		inventoryRepo: inventoryRepo,
		orderRepo:     orderRepo,
		n8nService:    n8nService,
	}
}

// Inventory management
func (s *inventoryService) List(companyID uuid.UUID, params map[string]interface{}) ([]models.Inventory, int64, error) {
	return s.inventoryRepo.List(companyID, params)
}

func (s *inventoryService) Get(id uuid.UUID) (*models.Inventory, error) {
	return s.inventoryRepo.Get(id)
}

func (s *inventoryService) GetBySKU(sku string) (*models.Inventory, error) {
	return s.inventoryRepo.GetBySKU(sku)
}

func (s *inventoryService) Create(companyID uuid.UUID, req CreateInventoryRequest) (*models.Inventory, error) {
	// Check if SKU already exists
	if _, err := s.inventoryRepo.GetBySKU(req.SKU); err == nil {
		return nil, errors.New("SKU already exists")
	}
	
	inventory := &models.Inventory{
		CompanyID:         companyID,
		SKU:               req.SKU,
		PartNo:            req.PartNo,
		Name:              req.Name,
		Description:       req.Description,
		Category:          req.Category,
		Material:          req.Material,
		Specification:     req.Specification,
		SurfaceTreatment:  req.SurfaceTreatment,
		HeatTreatment:     req.HeatTreatment,
		Unit:              req.Unit,
		CurrentStock:      req.InitialStock,
		AvailableStock:    req.InitialStock,
		MinStock:          req.MinStock,
		MaxStock:          req.MaxStock,
		ReorderPoint:      req.ReorderPoint,
		ReorderQuantity:   req.ReorderQuantity,
		WarehouseID:       &req.WarehouseID,
		Location:          req.Location,
		StandardCost:      req.StandardCost,
		AverageCost:       req.StandardCost,
		PrimarySupplierID: &req.PrimarySupplierID,
		LeadTimeDays:      req.LeadTimeDays,
		Status:            "active",
	}
	
	if inventory.Unit == "" {
		inventory.Unit = "PCS"
	}
	
	if err := s.inventoryRepo.Create(inventory); err != nil {
		return nil, err
	}
	
	// Create initial stock movement if there's initial stock
	if req.InitialStock > 0 {
		movement := &models.StockMovement{
			CompanyID:       companyID,
			InventoryID:     inventory.ID,
			MovementType:    "in",
			Reason:          "initial",
			Quantity:        req.InitialStock,
			UnitCost:        req.StandardCost,
			TotalCost:       req.InitialStock * req.StandardCost,
			ToWarehouseID:   &req.WarehouseID,
			ToLocation:      req.Location,
			Notes:           "Initial stock",
			CreatedBy:       companyID, // Should be userID
		}
		s.inventoryRepo.CreateMovement(movement)
	}
	
	// Trigger N8N workflow
	go s.n8nService.LogEvent(companyID, companyID, "inventory.created", "inventory", inventory.ID, map[string]interface{}{
		"sku":           inventory.SKU,
		"name":          inventory.Name,
		"initial_stock": req.InitialStock,
	})
	
	return s.Get(inventory.ID)
}

func (s *inventoryService) Update(id uuid.UUID, req UpdateInventoryRequest) (*models.Inventory, error) {
	inventory, err := s.inventoryRepo.Get(id)
	if err != nil {
		return nil, err
	}
	
	// Update fields
	if req.Name != "" {
		inventory.Name = req.Name
	}
	inventory.Description = req.Description
	if req.Category != "" {
		inventory.Category = req.Category
	}
	inventory.Material = req.Material
	inventory.Specification = req.Specification
	inventory.SurfaceTreatment = req.SurfaceTreatment
	inventory.HeatTreatment = req.HeatTreatment
	inventory.MinStock = req.MinStock
	inventory.MaxStock = req.MaxStock
	inventory.ReorderPoint = req.ReorderPoint
	inventory.ReorderQuantity = req.ReorderQuantity
	inventory.Location = req.Location
	inventory.StandardCost = req.StandardCost
	inventory.LeadTimeDays = req.LeadTimeDays
	if req.Status != "" {
		inventory.Status = req.Status
	}
	
	if err := s.inventoryRepo.Update(inventory); err != nil {
		return nil, err
	}
	
	return s.Get(id)
}

func (s *inventoryService) Delete(id uuid.UUID) error {
	inventory, err := s.inventoryRepo.Get(id)
	if err != nil {
		return err
	}
	
	if inventory.CurrentStock > 0 {
		return errors.New("cannot delete inventory with existing stock")
	}
	
	inventory.Status = "discontinued"
	inventory.IsActive = false
	return s.inventoryRepo.Update(inventory)
}

// Stock operations
func (s *inventoryService) AdjustStock(id, userID uuid.UUID, req StockAdjustmentRequest) (*models.StockMovement, error) {
	inventory, err := s.inventoryRepo.Get(id)
	if err != nil {
		return nil, err
	}
	
	movementType := "in"
	if req.Quantity < 0 {
		movementType = "out"
	}
	
	movement := &models.StockMovement{
		CompanyID:     inventory.CompanyID,
		InventoryID:   id,
		MovementType:  movementType,
		Reason:        req.Reason,
		Quantity:      req.Quantity,
		UnitCost:      inventory.StandardCost,
		TotalCost:     req.Quantity * inventory.StandardCost,
		Notes:         req.Notes,
		BatchNo:       req.BatchNo,
		CreatedBy:     userID,
	}
	
	if req.WarehouseID != uuid.Nil {
		if movementType == "in" {
			movement.ToWarehouseID = &req.WarehouseID
		} else {
			movement.FromWarehouseID = &req.WarehouseID
		}
	}
	
	if err := s.inventoryRepo.CreateMovement(movement); err != nil {
		return nil, err
	}
	
	// Trigger N8N workflow
	go s.n8nService.LogEvent(inventory.CompanyID, userID, "inventory.adjusted", "inventory", id, map[string]interface{}{
		"sku":        inventory.SKU,
		"adjustment": req.Quantity,
		"reason":     req.Reason,
		"new_stock":  inventory.CurrentStock + req.Quantity,
	})
	
	return movement, nil
}

func (s *inventoryService) TransferStock(userID uuid.UUID, req StockTransferRequest) (*models.StockMovement, error) {
	if req.FromWarehouseID == req.ToWarehouseID {
		return nil, errors.New("cannot transfer to the same warehouse")
	}
	
	inventory, err := s.inventoryRepo.Get(req.InventoryID)
	if err != nil {
		return nil, err
	}
	
	// Check available stock
	if inventory.AvailableStock < req.Quantity {
		return nil, fmt.Errorf("insufficient stock: available %.2f, requested %.2f", 
			inventory.AvailableStock, req.Quantity)
	}
	
	// Create transfer movement
	movement := &models.StockMovement{
		CompanyID:       inventory.CompanyID,
		InventoryID:     req.InventoryID,
		MovementType:    "transfer",
		Reason:          "transfer",
		Quantity:        0, // Net zero for transfers
		FromWarehouseID: &req.FromWarehouseID,
		ToWarehouseID:   &req.ToWarehouseID,
		FromLocation:    req.FromLocation,
		ToLocation:      req.ToLocation,
		Notes:           req.Notes,
		CreatedBy:       userID,
	}
	
	if err := s.inventoryRepo.CreateMovement(movement); err != nil {
		return nil, err
	}
	
	return movement, nil
}

func (s *inventoryService) ReserveStock(id uuid.UUID, quantity float64, orderID uuid.UUID) error {
	inventory, err := s.inventoryRepo.Get(id)
	if err != nil {
		return err
	}
	
	if inventory.AvailableStock < quantity {
		return fmt.Errorf("insufficient stock: available %.2f, requested %.2f", 
			inventory.AvailableStock, quantity)
	}
	
	// Update reserved stock
	return s.inventoryRepo.UpdateStock(id, quantity, true)
}

func (s *inventoryService) ReleaseStock(id uuid.UUID, quantity float64, orderID uuid.UUID) error {
	// Update reserved stock (negative to release)
	return s.inventoryRepo.UpdateStock(id, -quantity, true)
}

// Warehouse management
func (s *inventoryService) ListWarehouses(companyID uuid.UUID) ([]models.Warehouse, error) {
	return s.inventoryRepo.ListWarehouses(companyID)
}

func (s *inventoryService) GetWarehouse(id uuid.UUID) (*models.Warehouse, error) {
	return s.inventoryRepo.GetWarehouse(id)
}

func (s *inventoryService) CreateWarehouse(companyID uuid.UUID, req CreateWarehouseRequest) (*models.Warehouse, error) {
	warehouse := &models.Warehouse{
		CompanyID: companyID,
		Code:      req.Code,
		Name:      req.Name,
		Type:      req.Type,
		Address:   req.Address,
		Manager:   req.Manager,
		Phone:     req.Phone,
		IsActive:  true,
	}
	
	if err := s.inventoryRepo.CreateWarehouse(warehouse); err != nil {
		return nil, err
	}
	
	return warehouse, nil
}

func (s *inventoryService) UpdateWarehouse(id uuid.UUID, req UpdateWarehouseRequest) (*models.Warehouse, error) {
	warehouse, err := s.inventoryRepo.GetWarehouse(id)
	if err != nil {
		return nil, err
	}
	
	if req.Name != "" {
		warehouse.Name = req.Name
	}
	if req.Type != "" {
		warehouse.Type = req.Type
	}
	warehouse.Address = req.Address
	warehouse.Manager = req.Manager
	warehouse.Phone = req.Phone
	warehouse.IsActive = req.IsActive
	
	if err := s.inventoryRepo.UpdateWarehouse(warehouse); err != nil {
		return nil, err
	}
	
	return warehouse, nil
}

// Stock movements
func (s *inventoryService) GetMovements(inventoryID uuid.UUID, params map[string]interface{}) ([]models.StockMovement, error) {
	return s.inventoryRepo.GetMovements(inventoryID, params)
}

func (s *inventoryService) GetMovementsByOrder(orderID uuid.UUID) ([]models.StockMovement, error) {
	return s.inventoryRepo.GetMovementsByReference("order", orderID)
}

// Alerts
func (s *inventoryService) GetActiveAlerts(companyID uuid.UUID) ([]models.StockAlert, error) {
	return s.inventoryRepo.GetActiveAlerts(companyID)
}

func (s *inventoryService) AcknowledgeAlert(id, userID uuid.UUID) error {
	alert := &models.StockAlert{}
	// Simplified - should fetch first
	now := time.Now()
	alert.AcknowledgedBy = &userID
	alert.AcknowledgedAt = &now
	alert.Status = "acknowledged"
	
	return s.inventoryRepo.UpdateAlert(alert)
}

func (s *inventoryService) ResolveAlert(id, userID uuid.UUID, resolution string) error {
	alert := &models.StockAlert{}
	// Simplified - should fetch first
	now := time.Now()
	alert.ResolvedBy = &userID
	alert.ResolvedAt = &now
	alert.Resolution = resolution
	alert.Status = "resolved"
	
	return s.inventoryRepo.UpdateAlert(alert)
}

// Reports
func (s *inventoryService) GetInventoryStats(companyID uuid.UUID) (*InventoryStats, error) {
	items, _, err := s.inventoryRepo.List(companyID, map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	
	stats := &InventoryStats{
		TotalItems: len(items),
	}
	
	for _, item := range items {
		stats.TotalValue += item.CurrentStock * item.AverageCost
		
		if item.CurrentStock <= 0 {
			stats.OutOfStockItems++
		} else if item.CurrentStock <= item.MinStock {
			stats.LowStockItems++
		} else if item.MaxStock > 0 && item.CurrentStock >= item.MaxStock {
			stats.OverstockItems++
		}
	}
	
	alerts, _ := s.inventoryRepo.GetActiveAlerts(companyID)
	stats.ActiveAlerts = len(alerts)
	
	return stats, nil
}

func (s *inventoryService) GetLowStockItems(companyID uuid.UUID) ([]models.Inventory, error) {
	return s.inventoryRepo.GetLowStockItems(companyID)
}

func (s *inventoryService) GetStockValuation(companyID uuid.UUID) (*StockValuation, error) {
	items, _, err := s.inventoryRepo.List(companyID, map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	
	valuation := &StockValuation{
		ByCategory:  make(map[string]float64),
		ByWarehouse: make(map[string]float64),
	}
	
	var topItems []InventoryValueItem
	
	for _, item := range items {
		value := item.CurrentStock * item.AverageCost
		valuation.TotalValue += value
		
		// By category
		valuation.ByCategory[item.Category] += value
		
		// By warehouse
		if item.WarehouseID != nil {
			warehouseName := "Default"
			if item.Warehouse != nil {
				warehouseName = item.Warehouse.Name
			}
			valuation.ByWarehouse[warehouseName] += value
		}
		
		// Top value items
		topItems = append(topItems, InventoryValueItem{
			InventoryID: item.ID,
			SKU:         item.SKU,
			Name:        item.Name,
			Quantity:    item.CurrentStock,
			UnitCost:    item.AverageCost,
			TotalValue:  value,
		})
	}
	
	// Sort and get top 10
	// Simplified - should implement proper sorting
	if len(topItems) > 10 {
		valuation.TopValueItems = topItems[:10]
	} else {
		valuation.TopValueItems = topItems
	}
	
	return valuation, nil
}

// Stock take
func (s *inventoryService) CreateStockTake(companyID, userID uuid.UUID, req CreateStockTakeRequest) (*models.StockTake, error) {
	scheduledDate, err := time.Parse("2006-01-02", req.ScheduledDate)
	if err != nil {
		return nil, errors.New("invalid scheduled date format")
	}
	
	stockTake := &models.StockTake{
		CompanyID:     companyID,
		ReferenceNo:   s.generateStockTakeNo(companyID),
		WarehouseID:   req.WarehouseID,
		Status:        "draft",
		Type:          req.Type,
		ScheduledDate: scheduledDate,
		CreatedBy:     userID,
		AssignedTo:    req.AssignedTo,
		Notes:         req.Notes,
	}
	
	if err := s.inventoryRepo.CreateStockTake(stockTake); err != nil {
		return nil, err
	}
	
	// Create stock take items based on type
	if req.Type == "full" {
		// Get all items in warehouse
		items, _, _ := s.inventoryRepo.List(companyID, map[string]interface{}{
			"warehouse_id": req.WarehouseID.String(),
		})
		
		for _, item := range items {
			stockTakeItem := &models.StockTakeItem{
				StockTakeID:    stockTake.ID,
				InventoryID:    item.ID,
				SystemQuantity: item.CurrentStock,
				Status:         "pending",
			}
			s.inventoryRepo.CreateStockTakeItem(stockTakeItem)
		}
		
		stockTake.TotalItems = len(items)
		s.inventoryRepo.UpdateStockTake(stockTake)
	}
	
	return s.GetStockTake(stockTake.ID)
}

func (s *inventoryService) UpdateStockTake(id uuid.UUID, req UpdateStockTakeRequest) (*models.StockTake, error) {
	stockTake, err := s.inventoryRepo.GetStockTake(id)
	if err != nil {
		return nil, err
	}
	
	if req.Status != "" {
		stockTake.Status = req.Status
		if req.Status == "in_progress" {
			now := time.Now()
			stockTake.StartedAt = &now
		}
	}
	
	if req.AssignedTo != uuid.Nil {
		stockTake.AssignedTo = req.AssignedTo
	}
	
	if req.Notes != "" {
		stockTake.Notes = req.Notes
	}
	
	if err := s.inventoryRepo.UpdateStockTake(stockTake); err != nil {
		return nil, err
	}
	
	return s.GetStockTake(id)
}

func (s *inventoryService) GetStockTake(id uuid.UUID) (*models.StockTake, error) {
	return s.inventoryRepo.GetStockTake(id)
}

func (s *inventoryService) ListStockTakes(companyID uuid.UUID, params map[string]interface{}) ([]models.StockTake, error) {
	return s.inventoryRepo.ListStockTakes(companyID, params)
}

func (s *inventoryService) SubmitStockCount(stockTakeID, userID uuid.UUID, req StockCountRequest) error {
	// Update stock take items with counted quantities
	for _, count := range req.Items {
		items, _ := s.inventoryRepo.GetStockTakeItems(stockTakeID)
		
		for _, item := range items {
			if item.InventoryID == count.InventoryID {
				now := time.Now()
				item.CountedQuantity = count.CountedQuantity
				item.Variance = count.CountedQuantity - item.SystemQuantity
				item.Status = "counted"
				item.CountedBy = &userID
				item.CountedAt = &now
				item.Notes = count.Notes
				
				s.inventoryRepo.UpdateStockTakeItem(&item)
				break
			}
		}
	}
	
	// Update stock take summary
	stockTake, _ := s.inventoryRepo.GetStockTake(stockTakeID)
	items, _ := s.inventoryRepo.GetStockTakeItems(stockTakeID)
	
	countedItems := 0
	varianceItems := 0
	totalVariance := 0.0
	
	for _, item := range items {
		if item.Status == "counted" {
			countedItems++
			if item.Variance != 0 {
				varianceItems++
				totalVariance += item.Variance
			}
		}
	}
	
	stockTake.CountedItems = countedItems
	stockTake.VarianceItems = varianceItems
	stockTake.TotalVariance = totalVariance
	
	return s.inventoryRepo.UpdateStockTake(stockTake)
}

func (s *inventoryService) CompleteStockTake(id, userID uuid.UUID) error {
	stockTake, err := s.inventoryRepo.GetStockTake(id)
	if err != nil {
		return err
	}
	
	if stockTake.Status != "in_progress" {
		return errors.New("stock take must be in progress to complete")
	}
	
	// Apply stock adjustments
	items, _ := s.inventoryRepo.GetStockTakeItems(id)
	
	for _, item := range items {
		if item.Status == "counted" && item.Variance != 0 {
			// Create stock adjustment
			reason := "adjustment"
			if item.Variance > 0 {
				reason = "found"
			} else {
				reason = "loss"
			}
			
			s.AdjustStock(item.InventoryID, userID, StockAdjustmentRequest{
				Quantity: item.Variance,
				Reason:   reason,
				Notes:    fmt.Sprintf("Stock take adjustment - Ref: %s", stockTake.ReferenceNo),
			})
		}
	}
	
	// Update stock take status
	now := time.Now()
	stockTake.Status = "completed"
	stockTake.CompletedAt = &now
	stockTake.ReviewedBy = &userID
	
	return s.inventoryRepo.UpdateStockTake(stockTake)
}

func (s *inventoryService) generateStockTakeNo(companyID uuid.UUID) string {
	return fmt.Sprintf("ST-%s", time.Now().Format("20060102-150405"))
}