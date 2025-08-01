package repository

import (
	"fmt"
	"github.com/fastenmind/fastener-api/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InventoryRepository interface {
	// Inventory CRUD
	Create(inventory *models.Inventory) error
	Update(inventory *models.Inventory) error
	Delete(id uuid.UUID) error
	Get(id uuid.UUID) (*models.Inventory, error)
	GetBySKU(sku string) (*models.Inventory, error)
	List(companyID uuid.UUID, params map[string]interface{}) ([]models.Inventory, int64, error)
	
	// Stock operations
	UpdateStock(id uuid.UUID, quantity float64, isReserved bool) error
	GetLowStockItems(companyID uuid.UUID) ([]models.Inventory, error)
	GetOverstockItems(companyID uuid.UUID) ([]models.Inventory, error)
	
	// Warehouse operations
	CreateWarehouse(warehouse *models.Warehouse) error
	UpdateWarehouse(warehouse *models.Warehouse) error
	GetWarehouse(id uuid.UUID) (*models.Warehouse, error)
	ListWarehouses(companyID uuid.UUID) ([]models.Warehouse, error)
	
	// Stock movements
	CreateMovement(movement *models.StockMovement) error
	GetMovements(inventoryID uuid.UUID, params map[string]interface{}) ([]models.StockMovement, error)
	GetMovementsByReference(refType string, refID uuid.UUID) ([]models.StockMovement, error)
	
	// Stock alerts
	CreateAlert(alert *models.StockAlert) error
	UpdateAlert(alert *models.StockAlert) error
	GetActiveAlerts(companyID uuid.UUID) ([]models.StockAlert, error)
	GetAlertsByInventory(inventoryID uuid.UUID) ([]models.StockAlert, error)
	
	// Stock take
	CreateStockTake(stockTake *models.StockTake) error
	UpdateStockTake(stockTake *models.StockTake) error
	GetStockTake(id uuid.UUID) (*models.StockTake, error)
	ListStockTakes(companyID uuid.UUID, params map[string]interface{}) ([]models.StockTake, error)
	CreateStockTakeItem(item *models.StockTakeItem) error
	UpdateStockTakeItem(item *models.StockTakeItem) error
	GetStockTakeItems(stockTakeID uuid.UUID) ([]models.StockTakeItem, error)
}

type inventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db interface{}) InventoryRepository {
	gormDB, ok := db.(*gorm.DB)
	if !ok {
		panic("invalid database type, expected *gorm.DB")
	}
	return &inventoryRepository{db: gormDB}
}

// Inventory CRUD
func (r *inventoryRepository) Create(inventory *models.Inventory) error {
	return r.db.Create(inventory).Error
}

func (r *inventoryRepository) Update(inventory *models.Inventory) error {
	return r.db.Save(inventory).Error
}

func (r *inventoryRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Inventory{}, id).Error
}

func (r *inventoryRepository) Get(id uuid.UUID) (*models.Inventory, error) {
	var inventory models.Inventory
	err := r.db.Preload("Warehouse").
		Preload("PrimarySupplier").
		First(&inventory, id).Error
	if err != nil {
		return nil, err
	}
	return &inventory, nil
}

func (r *inventoryRepository) GetBySKU(sku string) (*models.Inventory, error) {
	var inventory models.Inventory
	err := r.db.Where("sku = ?", sku).First(&inventory).Error
	if err != nil {
		return nil, err
	}
	return &inventory, nil
}

func (r *inventoryRepository) List(companyID uuid.UUID, params map[string]interface{}) ([]models.Inventory, int64, error) {
	var items []models.Inventory
	var total int64
	
	query := r.db.Model(&models.Inventory{}).Where("company_id = ?", companyID)
	
	// Apply filters
	if category, ok := params["category"].(string); ok && category != "" {
		query = query.Where("category = ?", category)
	}
	
	if warehouseID, ok := params["warehouse_id"].(string); ok && warehouseID != "" {
		query = query.Where("warehouse_id = ?", warehouseID)
	}
	
	if status, ok := params["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	
	if lowStock, ok := params["low_stock"].(bool); ok && lowStock {
		query = query.Where("current_stock <= min_stock")
	}
	
	if search, ok := params["search"].(string); ok && search != "" {
		query = query.Where("sku LIKE ? OR part_no LIKE ? OR name LIKE ?", 
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
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
	
	// Apply sorting
	sortBy := "created_at"
	if sb, ok := params["sort_by"].(string); ok && sb != "" {
		sortBy = sb
	}
	sortOrder := "DESC"
	if so, ok := params["sort_order"].(string); ok && so != "" {
		sortOrder = so
	}
	query = query.Order(sortBy + " " + sortOrder)
	
	// Load with relations
	if err := query.
		Preload("Warehouse").
		Preload("PrimarySupplier").
		Find(&items).Error; err != nil {
		return nil, 0, err
	}
	
	return items, total, nil
}

// Stock operations
func (r *inventoryRepository) UpdateStock(id uuid.UUID, quantity float64, isReserved bool) error {
	var inventory models.Inventory
	if err := r.db.First(&inventory, id).Error; err != nil {
		return err
	}
	
	if isReserved {
		inventory.ReservedStock += quantity
		inventory.AvailableStock = inventory.CurrentStock - inventory.ReservedStock
	} else {
		inventory.CurrentStock += quantity
		inventory.AvailableStock = inventory.CurrentStock - inventory.ReservedStock
	}
	
	return r.db.Save(&inventory).Error
}

func (r *inventoryRepository) GetLowStockItems(companyID uuid.UUID) ([]models.Inventory, error) {
	var items []models.Inventory
	err := r.db.Where("company_id = ? AND current_stock <= min_stock AND is_active = ?", 
		companyID, true).
		Preload("Warehouse").
		Find(&items).Error
	return items, err
}

func (r *inventoryRepository) GetOverstockItems(companyID uuid.UUID) ([]models.Inventory, error) {
	var items []models.Inventory
	err := r.db.Where("company_id = ? AND current_stock >= max_stock AND max_stock > 0 AND is_active = ?", 
		companyID, true).
		Preload("Warehouse").
		Find(&items).Error
	return items, err
}

// Warehouse operations
func (r *inventoryRepository) CreateWarehouse(warehouse *models.Warehouse) error {
	return r.db.Create(warehouse).Error
}

func (r *inventoryRepository) UpdateWarehouse(warehouse *models.Warehouse) error {
	return r.db.Save(warehouse).Error
}

func (r *inventoryRepository) GetWarehouse(id uuid.UUID) (*models.Warehouse, error) {
	var warehouse models.Warehouse
	err := r.db.First(&warehouse, id).Error
	if err != nil {
		return nil, err
	}
	return &warehouse, nil
}

func (r *inventoryRepository) ListWarehouses(companyID uuid.UUID) ([]models.Warehouse, error) {
	var warehouses []models.Warehouse
	err := r.db.Where("company_id = ? AND is_active = ?", companyID, true).
		Order("name").
		Find(&warehouses).Error
	return warehouses, err
}

// Stock movements
func (r *inventoryRepository) CreateMovement(movement *models.StockMovement) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Create movement record
		if err := tx.Create(movement).Error; err != nil {
			return err
		}
		
		// Update inventory stock
		var inventory models.Inventory
		if err := tx.First(&inventory, movement.InventoryID).Error; err != nil {
			return err
		}
		
		// Update before/after quantities
		movement.BeforeQuantity = inventory.CurrentStock
		inventory.CurrentStock += movement.Quantity
		inventory.AvailableStock = inventory.CurrentStock - inventory.ReservedStock
		movement.AfterQuantity = inventory.CurrentStock
		
		// Save updated inventory
		if err := tx.Save(&inventory).Error; err != nil {
			return err
		}
		
		// Update movement with before/after quantities
		if err := tx.Save(movement).Error; err != nil {
			return err
		}
		
		// Check for alerts
		if inventory.CurrentStock <= inventory.MinStock {
			alert := &models.StockAlert{
				CompanyID:      inventory.CompanyID,
				InventoryID:    inventory.ID,
				AlertType:      "low_stock",
				Priority:       "high",
				CurrentLevel:   inventory.CurrentStock,
				ThresholdLevel: inventory.MinStock,
				Message:        fmt.Sprintf("Low stock alert: %s (SKU: %s) - Current: %.2f, Min: %.2f", 
					inventory.Name, inventory.SKU, inventory.CurrentStock, inventory.MinStock),
			}
			tx.Create(alert)
		}
		
		return nil
	})
}

func (r *inventoryRepository) GetMovements(inventoryID uuid.UUID, params map[string]interface{}) ([]models.StockMovement, error) {
	var movements []models.StockMovement
	query := r.db.Where("inventory_id = ?", inventoryID)
	
	if movementType, ok := params["movement_type"].(string); ok && movementType != "" {
		query = query.Where("movement_type = ?", movementType)
	}
	
	if startDate, ok := params["start_date"].(string); ok && startDate != "" {
		query = query.Where("created_at >= ?", startDate)
	}
	
	if endDate, ok := params["end_date"].(string); ok && endDate != "" {
		query = query.Where("created_at <= ?", endDate)
	}
	
	err := query.
		Preload("Inventory").
		Preload("FromWarehouse").
		Preload("ToWarehouse").
		Preload("Creator").
		Order("created_at DESC").
		Find(&movements).Error
		
	return movements, err
}

func (r *inventoryRepository) GetMovementsByReference(refType string, refID uuid.UUID) ([]models.StockMovement, error) {
	var movements []models.StockMovement
	err := r.db.Where("reference_type = ? AND reference_id = ?", refType, refID).
		Preload("Inventory").
		Order("created_at DESC").
		Find(&movements).Error
	return movements, err
}

// Stock alerts
func (r *inventoryRepository) CreateAlert(alert *models.StockAlert) error {
	return r.db.Create(alert).Error
}

func (r *inventoryRepository) UpdateAlert(alert *models.StockAlert) error {
	return r.db.Save(alert).Error
}

func (r *inventoryRepository) GetActiveAlerts(companyID uuid.UUID) ([]models.StockAlert, error) {
	var alerts []models.StockAlert
	err := r.db.Where("company_id = ? AND status = ?", companyID, "active").
		Preload("Inventory").
		Order("priority DESC, created_at DESC").
		Find(&alerts).Error
	return alerts, err
}

func (r *inventoryRepository) GetAlertsByInventory(inventoryID uuid.UUID) ([]models.StockAlert, error) {
	var alerts []models.StockAlert
	err := r.db.Where("inventory_id = ?", inventoryID).
		Order("created_at DESC").
		Find(&alerts).Error
	return alerts, err
}

// Stock take
func (r *inventoryRepository) CreateStockTake(stockTake *models.StockTake) error {
	return r.db.Create(stockTake).Error
}

func (r *inventoryRepository) UpdateStockTake(stockTake *models.StockTake) error {
	return r.db.Save(stockTake).Error
}

func (r *inventoryRepository) GetStockTake(id uuid.UUID) (*models.StockTake, error) {
	var stockTake models.StockTake
	err := r.db.Preload("Warehouse").
		Preload("Creator").
		Preload("Assignee").
		Preload("Reviewer").
		First(&stockTake, id).Error
	if err != nil {
		return nil, err
	}
	return &stockTake, nil
}

func (r *inventoryRepository) ListStockTakes(companyID uuid.UUID, params map[string]interface{}) ([]models.StockTake, error) {
	var stockTakes []models.StockTake
	query := r.db.Where("company_id = ?", companyID)
	
	if status, ok := params["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	
	if warehouseID, ok := params["warehouse_id"].(string); ok && warehouseID != "" {
		query = query.Where("warehouse_id = ?", warehouseID)
	}
	
	err := query.
		Preload("Warehouse").
		Order("created_at DESC").
		Find(&stockTakes).Error
		
	return stockTakes, err
}

func (r *inventoryRepository) CreateStockTakeItem(item *models.StockTakeItem) error {
	return r.db.Create(item).Error
}

func (r *inventoryRepository) UpdateStockTakeItem(item *models.StockTakeItem) error {
	return r.db.Save(item).Error
}

func (r *inventoryRepository) GetStockTakeItems(stockTakeID uuid.UUID) ([]models.StockTakeItem, error) {
	var items []models.StockTakeItem
	err := r.db.Where("stock_take_id = ?", stockTakeID).
		Preload("Inventory").
		Preload("Counter").
		Preload("Verifier").
		Find(&items).Error
	return items, err
}