package repository

import (
	"time"
	"github.com/fastenmind/fastener-api/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductionRepository interface {
	// Production Order operations
	CreateProductionOrder(order *models.ProductionOrder) error
	UpdateProductionOrder(order *models.ProductionOrder) error
	GetProductionOrder(id uuid.UUID) (*models.ProductionOrder, error)
	ListProductionOrders(companyID uuid.UUID, params map[string]interface{}) ([]models.ProductionOrder, int64, error)
	
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
	GetTasksByProductionOrder(productionOrderID uuid.UUID) ([]models.ProductionTask, error)
	
	// Production Material operations
	CreateProductionMaterial(material *models.ProductionMaterial) error
	UpdateProductionMaterial(material *models.ProductionMaterial) error
	GetProductionMaterials(productionOrderID uuid.UUID) ([]models.ProductionMaterial, error)
	
	// Quality Inspection operations
	CreateQualityInspection(inspection *models.QualityInspection) error
	UpdateQualityInspection(inspection *models.QualityInspection) error
	GetQualityInspection(id uuid.UUID) (*models.QualityInspection, error)
	ListQualityInspections(companyID uuid.UUID, params map[string]interface{}) ([]models.QualityInspection, int64, error)
}

type productionRepository struct {
	db *gorm.DB
}

func NewProductionRepository(db interface{}) ProductionRepository {
	gormDB, ok := db.(*gorm.DB)
	if !ok {
		panic("invalid database type, expected *gorm.DB")
	}
	return &productionRepository{db: gormDB}
}

// Production Order operations
func (r *productionRepository) CreateProductionOrder(order *models.ProductionOrder) error {
	return r.db.Create(order).Error
}

func (r *productionRepository) UpdateProductionOrder(order *models.ProductionOrder) error {
	return r.db.Save(order).Error
}

func (r *productionRepository) GetProductionOrder(id uuid.UUID) (*models.ProductionOrder, error) {
	var order models.ProductionOrder
	err := r.db.Preload("SalesOrder").
		Preload("Customer").
		Preload("Inventory").
		Preload("Route").
		Preload("CurrentStation").
		Preload("Creator").
		First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *productionRepository) ListProductionOrders(companyID uuid.UUID, params map[string]interface{}) ([]models.ProductionOrder, int64, error) {
	var orders []models.ProductionOrder
	var total int64

	query := r.db.Model(&models.ProductionOrder{}).Where("company_id = ?", companyID)

	// Apply filters
	if status, ok := params["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}

	if priority, ok := params["priority"].(string); ok && priority != "" {
		query = query.Where("priority = ?", priority)
	}

	if customerID, ok := params["customer_id"].(string); ok && customerID != "" {
		query = query.Where("customer_id = ?", customerID)
	}

	if inventoryID, ok := params["inventory_id"].(string); ok && inventoryID != "" {
		query = query.Where("inventory_id = ?", inventoryID)
	}

	// Date range filters
	if startDate, ok := params["start_date"].(string); ok && startDate != "" {
		query = query.Where("planned_start_date >= ?", startDate)
	}

	if endDate, ok := params["end_date"].(string); ok && endDate != "" {
		query = query.Where("planned_end_date <= ?", endDate)
	}

	if search, ok := params["search"].(string); ok && search != "" {
		query = query.Where("order_no LIKE ? OR product_name LIKE ?", "%"+search+"%", "%"+search+"%")
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
		Preload("SalesOrder").
		Preload("Customer").
		Preload("Inventory").
		Preload("Route").
		Preload("CurrentStation").
		Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// Production Route operations
func (r *productionRepository) CreateProductionRoute(route *models.ProductionRoute) error {
	return r.db.Create(route).Error
}

func (r *productionRepository) UpdateProductionRoute(route *models.ProductionRoute) error {
	return r.db.Save(route).Error
}

func (r *productionRepository) GetProductionRoute(id uuid.UUID) (*models.ProductionRoute, error) {
	var route models.ProductionRoute
	err := r.db.Preload("Inventory").
		Preload("Creator").
		Preload("Operations").
		Preload("Operations.WorkStation").
		First(&route, id).Error
	if err != nil {
		return nil, err
	}
	return &route, nil
}

func (r *productionRepository) ListProductionRoutes(companyID uuid.UUID, params map[string]interface{}) ([]models.ProductionRoute, error) {
	var routes []models.ProductionRoute
	query := r.db.Where("company_id = ?", companyID)

	// Apply filters
	if status, ok := params["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}

	if inventoryID, ok := params["inventory_id"].(string); ok && inventoryID != "" {
		query = query.Where("inventory_id = ?", inventoryID)
	}

	if category, ok := params["product_category"].(string); ok && category != "" {
		query = query.Where("product_category = ?", category)
	}

	if search, ok := params["search"].(string); ok && search != "" {
		query = query.Where("route_no LIKE ? OR name LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Apply sorting
	query = query.Order("name ASC")

	// Load with relations
	if err := query.
		Preload("Inventory").
		Preload("Creator").
		Find(&routes).Error; err != nil {
		return nil, err
	}

	return routes, nil
}

// Route Operation operations
func (r *productionRepository) CreateRouteOperation(operation *models.RouteOperation) error {
	return r.db.Create(operation).Error
}

func (r *productionRepository) UpdateRouteOperation(operation *models.RouteOperation) error {
	return r.db.Save(operation).Error
}

func (r *productionRepository) GetRouteOperations(routeID uuid.UUID) ([]models.RouteOperation, error) {
	var operations []models.RouteOperation
	err := r.db.Where("route_id = ?", routeID).
		Preload("WorkStation").
		Order("operation_no ASC").
		Find(&operations).Error
	return operations, err
}

func (r *productionRepository) DeleteRouteOperation(id uuid.UUID) error {
	return r.db.Delete(&models.RouteOperation{}, id).Error
}

// Work Station operations
func (r *productionRepository) CreateWorkStation(station *models.WorkStation) error {
	return r.db.Create(station).Error
}

func (r *productionRepository) UpdateWorkStation(station *models.WorkStation) error {
	return r.db.Save(station).Error
}

func (r *productionRepository) GetWorkStation(id uuid.UUID) (*models.WorkStation, error) {
	var station models.WorkStation
	err := r.db.Preload("Creator").First(&station, id).Error
	if err != nil {
		return nil, err
	}
	return &station, nil
}

func (r *productionRepository) ListWorkStations(companyID uuid.UUID, params map[string]interface{}) ([]models.WorkStation, error) {
	var stations []models.WorkStation
	query := r.db.Where("company_id = ?", companyID)

	// Apply filters
	if stationType, ok := params["type"].(string); ok && stationType != "" {
		query = query.Where("type = ?", stationType)
	}

	if status, ok := params["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}

	if department, ok := params["department"].(string); ok && department != "" {
		query = query.Where("department = ?", department)
	}

	if search, ok := params["search"].(string); ok && search != "" {
		query = query.Where("station_no LIKE ? OR name LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Apply sorting
	query = query.Order("station_no ASC")

	// Load with relations
	if err := query.Preload("Creator").Find(&stations).Error; err != nil {
		return nil, err
	}

	return stations, nil
}

// Production Task operations
func (r *productionRepository) CreateProductionTask(task *models.ProductionTask) error {
	return r.db.Create(task).Error
}

func (r *productionRepository) UpdateProductionTask(task *models.ProductionTask) error {
	return r.db.Save(task).Error
}

func (r *productionRepository) GetProductionTask(id uuid.UUID) (*models.ProductionTask, error) {
	var task models.ProductionTask
	err := r.db.Preload("ProductionOrder").
		Preload("RouteOperation").
		Preload("RouteOperation.WorkStation").
		Preload("WorkStation").
		Preload("AssignedUser").
		Preload("QCUser").
		First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *productionRepository) ListProductionTasks(params map[string]interface{}) ([]models.ProductionTask, error) {
	var tasks []models.ProductionTask
	query := r.db.Model(&models.ProductionTask{})

	// Apply filters
	if status, ok := params["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}

	if assignedTo, ok := params["assigned_to"].(string); ok && assignedTo != "" {
		query = query.Where("assigned_to = ?", assignedTo)
	}

	if workStationID, ok := params["work_station_id"].(string); ok && workStationID != "" {
		query = query.Where("work_station_id = ?", workStationID)
	}

	if productionOrderID, ok := params["production_order_id"].(string); ok && productionOrderID != "" {
		query = query.Where("production_order_id = ?", productionOrderID)
	}

	// Date range filters
	if startDate, ok := params["start_date"].(string); ok && startDate != "" {
		query = query.Where("planned_start_time >= ?", startDate)
	}

	if endDate, ok := params["end_date"].(string); ok && endDate != "" {
		query = query.Where("planned_end_time <= ?", endDate)
	}

	// Apply sorting
	query = query.Order("planned_start_time ASC")

	// Load with relations
	if err := query.
		Preload("ProductionOrder").
		Preload("RouteOperation").
		Preload("WorkStation").
		Preload("AssignedUser").
		Find(&tasks).Error; err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *productionRepository) GetTasksByProductionOrder(productionOrderID uuid.UUID) ([]models.ProductionTask, error) {
	var tasks []models.ProductionTask
	err := r.db.Where("production_order_id = ?", productionOrderID).
		Preload("RouteOperation").
		Preload("WorkStation").
		Preload("AssignedUser").
		Order("task_no ASC").
		Find(&tasks).Error
	return tasks, err
}

// Production Material operations
func (r *productionRepository) CreateProductionMaterial(material *models.ProductionMaterial) error {
	return r.db.Create(material).Error
}

func (r *productionRepository) UpdateProductionMaterial(material *models.ProductionMaterial) error {
	return r.db.Save(material).Error
}

func (r *productionRepository) GetProductionMaterials(productionOrderID uuid.UUID) ([]models.ProductionMaterial, error) {
	var materials []models.ProductionMaterial
	err := r.db.Where("production_order_id = ?", productionOrderID).
		Preload("Inventory").
		Find(&materials).Error
	return materials, err
}

// Quality Inspection operations
func (r *productionRepository) CreateQualityInspection(inspection *models.QualityInspection) error {
	return r.db.Create(inspection).Error
}

func (r *productionRepository) UpdateQualityInspection(inspection *models.QualityInspection) error {
	return r.db.Save(inspection).Error
}

func (r *productionRepository) GetQualityInspection(id uuid.UUID) (*models.QualityInspection, error) {
	var inspection models.QualityInspection
	err := r.db.Preload("ProductionOrder").
		Preload("ProductionTask").
		Preload("Inventory").
		Preload("Inspector").
		Preload("Approver").
		First(&inspection, id).Error
	if err != nil {
		return nil, err
	}
	return &inspection, nil
}

func (r *productionRepository) ListQualityInspections(companyID uuid.UUID, params map[string]interface{}) ([]models.QualityInspection, int64, error) {
	var inspections []models.QualityInspection
	var total int64

	query := r.db.Model(&models.QualityInspection{}).Where("company_id = ?", companyID)

	// Apply filters
	if inspectionType, ok := params["type"].(string); ok && inspectionType != "" {
		query = query.Where("type = ?", inspectionType)
	}

	if status, ok := params["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}

	if productionOrderID, ok := params["production_order_id"].(string); ok && productionOrderID != "" {
		query = query.Where("production_order_id = ?", productionOrderID)
	}

	if inspectorID, ok := params["inspector_id"].(string); ok && inspectorID != "" {
		query = query.Where("inspector_id = ?", inspectorID)
	}

	// Date range filters
	if startDate, ok := params["start_date"].(string); ok && startDate != "" {
		query = query.Where("inspected_at >= ?", startDate)
	}

	if endDate, ok := params["end_date"].(string); ok && endDate != "" {
		query = query.Where("inspected_at <= ?", endDate)
	}

	if search, ok := params["search"].(string); ok && search != "" {
		query = query.Where("inspection_no LIKE ?", "%"+search+"%")
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
	query = query.Order("inspected_at DESC")

	// Load with relations
	if err := query.
		Preload("ProductionOrder").
		Preload("ProductionTask").
		Preload("Inventory").
		Preload("Inspector").
		Find(&inspections).Error; err != nil {
		return nil, 0, err
	}

	return inspections, total, nil
}