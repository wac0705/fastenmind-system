package repository

import (
	"github.com/fastenmind/fastener-api/internal/models"
	"github.com/google/uuid"
)

// Missing N8NRepository methods
func (r *n8nRepositoryGorm) DeleteWebhook(id uuid.UUID) error { return ErrNotImplemented }
func (r *n8nRepositoryGorm) GetWebhook(id uuid.UUID) (*models.N8NWebhook, error) { return nil, ErrNotImplemented }
func (r *n8nRepositoryGorm) GetWorkflow(id uuid.UUID) (*models.N8NWorkflow, error) { return nil, ErrNotImplemented }
func (r *n8nRepositoryGorm) ListWorkflows(companyID uuid.UUID, params map[string]interface{}) ([]models.N8NWorkflow, int64, error) { return nil, 0, ErrNotImplemented }
func (r *n8nRepositoryGorm) GetWorkflowByN8NID(workflowID string) (*models.N8NWorkflow, error) { return nil, ErrNotImplemented }
func (r *n8nRepositoryGorm) GetExecution(id uuid.UUID) (*models.N8NExecution, error) { return nil, ErrNotImplemented }
func (r *n8nRepositoryGorm) GetExecutionByN8NID(executionID string) (*models.N8NExecution, error) { return nil, ErrNotImplemented }
func (r *n8nRepositoryGorm) ListExecutions(companyID uuid.UUID, params map[string]interface{}) ([]models.N8NExecution, int64, error) { return nil, 0, ErrNotImplemented }
func (r *n8nRepositoryGorm) UpdateWebhook(webhook *models.N8NWebhook) error { return ErrNotImplemented }
func (r *n8nRepositoryGorm) ListWebhooks(companyID uuid.UUID) ([]models.N8NWebhook, error) { return nil, ErrNotImplemented }
func (r *n8nRepositoryGorm) GetWebhooksByEventType(companyID uuid.UUID, eventType string) ([]models.N8NWebhook, error) { return nil, ErrNotImplemented }
func (r *n8nRepositoryGorm) CreateScheduledTask(task *models.N8NScheduledTask) error { return ErrNotImplemented }
func (r *n8nRepositoryGorm) UpdateScheduledTask(task *models.N8NScheduledTask) error { return ErrNotImplemented }
func (r *n8nRepositoryGorm) DeleteScheduledTask(id uuid.UUID) error { return ErrNotImplemented }
func (r *n8nRepositoryGorm) GetScheduledTask(id uuid.UUID) (*models.N8NScheduledTask, error) { return nil, ErrNotImplemented }
func (r *n8nRepositoryGorm) ListScheduledTasks(companyID uuid.UUID) ([]models.N8NScheduledTask, error) { return nil, ErrNotImplemented }
func (r *n8nRepositoryGorm) GetDueScheduledTasks() ([]models.N8NScheduledTask, error) { return nil, ErrNotImplemented }
func (r *n8nRepositoryGorm) GetUnprocessedEvents(companyID uuid.UUID) ([]models.N8NEventLog, error) { return nil, ErrNotImplemented }
func (r *n8nRepositoryGorm) MarkEventProcessed(id uuid.UUID, workflowIDs []string) error { return ErrNotImplemented }

// Missing QuoteRepository methods
func (r *quoteRepositoryGorm) Update(quote *models.Quote) error { return ErrNotImplemented }
func (r *quoteRepositoryGorm) Delete(id uuid.UUID) error { return ErrNotImplemented }
func (r *quoteRepositoryGorm) Get(id uuid.UUID) (*models.Quote, error) { return nil, ErrNotImplemented }
func (r *quoteRepositoryGorm) GetWithDetails(id uuid.UUID) (*models.Quote, error) { return nil, ErrNotImplemented }
func (r *quoteRepositoryGorm) List(companyID uuid.UUID, params map[string]interface{}) ([]models.Quote, int64, error) { return nil, 0, ErrNotImplemented }
func (r *quoteRepositoryGorm) CreateVersion(version *models.QuoteVersion) error { return ErrNotImplemented }
func (r *quoteRepositoryGorm) GetVersions(quoteID uuid.UUID) ([]models.QuoteVersion, error) { return nil, ErrNotImplemented }
func (r *quoteRepositoryGorm) LogActivity(activity *models.QuoteActivity) error { return ErrNotImplemented }
func (r *quoteRepositoryGorm) GetActivities(quoteID uuid.UUID) ([]models.QuoteActivity, error) { return nil, ErrNotImplemented }

// Missing OrderRepository methods
func (r *orderRepositoryGorm) Delete(id uuid.UUID) error { return ErrNotImplemented }
func (r *orderRepositoryGorm) Update(order *models.Order) error { return ErrNotImplemented }
func (r *orderRepositoryGorm) Create(order *models.Order) error { return ErrNotImplemented }
func (r *orderRepositoryGorm) Get(id uuid.UUID) (*models.Order, error) { return nil, ErrNotImplemented }
func (r *orderRepositoryGorm) GetOrder(id uuid.UUID) (*models.Order, error) { return nil, ErrNotImplemented }
func (r *orderRepositoryGorm) GetWithDetails(id uuid.UUID) (*models.Order, error) { return nil, ErrNotImplemented }
func (r *orderRepositoryGorm) List(companyID uuid.UUID, params map[string]interface{}) ([]models.Order, int64, error) { return nil, 0, ErrNotImplemented }
func (r *orderRepositoryGorm) GetByOrderNo(orderNo string) (*models.Order, error) { return nil, ErrNotImplemented }
func (r *orderRepositoryGorm) CreateItem(item *models.OrderItem) error { return ErrNotImplemented }
func (r *orderRepositoryGorm) UpdateItem(item *models.OrderItem) error { return ErrNotImplemented }
func (r *orderRepositoryGorm) DeleteItem(id uuid.UUID) error { return ErrNotImplemented }
func (r *orderRepositoryGorm) GetItems(orderID uuid.UUID) ([]models.OrderItem, error) { return nil, ErrNotImplemented }
// GetOrderItems is already implemented in stubs_gorm.go
func (r *orderRepositoryGorm) LogActivity(activity *models.OrderActivity) error { return ErrNotImplemented }
func (r *orderRepositoryGorm) GetActivities(orderID uuid.UUID) ([]models.OrderActivity, error) { return nil, ErrNotImplemented }
func (r *orderRepositoryGorm) AddDocument(doc *models.OrderDocument) error { return ErrNotImplemented }
func (r *orderRepositoryGorm) RemoveDocument(id uuid.UUID) error { return ErrNotImplemented }
func (r *orderRepositoryGorm) GetDocuments(orderID uuid.UUID) ([]models.OrderDocument, error) { return nil, ErrNotImplemented }

// Missing InventoryRepository methods
func (r *inventoryRepositoryGorm) Create(inventory *models.Inventory) error { return ErrNotImplemented }
func (r *inventoryRepositoryGorm) Update(inventory *models.Inventory) error { return ErrNotImplemented }
func (r *inventoryRepositoryGorm) Delete(id uuid.UUID) error { return ErrNotImplemented }
func (r *inventoryRepositoryGorm) Get(id uuid.UUID) (*models.Inventory, error) { return nil, ErrNotImplemented }
func (r *inventoryRepositoryGorm) GetBySKU(sku string) (*models.Inventory, error) { return nil, ErrNotImplemented }
func (r *inventoryRepositoryGorm) List(companyID uuid.UUID, params map[string]interface{}) ([]models.Inventory, int64, error) { return nil, 0, ErrNotImplemented }
func (r *inventoryRepositoryGorm) UpdateStock(id uuid.UUID, quantity float64, isReserved bool) error { return ErrNotImplemented }
func (r *inventoryRepositoryGorm) GetOverstockItems(companyID uuid.UUID) ([]models.Inventory, error) { return nil, ErrNotImplemented }
func (r *inventoryRepositoryGorm) CreateWarehouse(warehouse *models.Warehouse) error { return ErrNotImplemented }
func (r *inventoryRepositoryGorm) UpdateWarehouse(warehouse *models.Warehouse) error { return ErrNotImplemented }
func (r *inventoryRepositoryGorm) GetWarehouse(id uuid.UUID) (*models.Warehouse, error) { return nil, ErrNotImplemented }
func (r *inventoryRepositoryGorm) ListWarehouses(companyID uuid.UUID) ([]models.Warehouse, error) { return nil, ErrNotImplemented }
func (r *inventoryRepositoryGorm) CreateMovement(movement *models.StockMovement) error { return ErrNotImplemented }
func (r *inventoryRepositoryGorm) GetMovements(inventoryID uuid.UUID, params map[string]interface{}) ([]models.StockMovement, error) { return nil, ErrNotImplemented }
func (r *inventoryRepositoryGorm) GetMovementsByReference(refType string, refID uuid.UUID) ([]models.StockMovement, error) { return nil, ErrNotImplemented }
func (r *inventoryRepositoryGorm) CreateAlert(alert *models.StockAlert) error { return ErrNotImplemented }
func (r *inventoryRepositoryGorm) UpdateAlert(alert *models.StockAlert) error { return ErrNotImplemented }
func (r *inventoryRepositoryGorm) GetActiveAlerts(companyID uuid.UUID) ([]models.StockAlert, error) { return nil, ErrNotImplemented }
func (r *inventoryRepositoryGorm) GetAlertsByInventory(inventoryID uuid.UUID) ([]models.StockAlert, error) { return nil, ErrNotImplemented }
func (r *inventoryRepositoryGorm) CreateStockTake(stockTake *models.StockTake) error { return ErrNotImplemented }
func (r *inventoryRepositoryGorm) UpdateStockTake(stockTake *models.StockTake) error { return ErrNotImplemented }
func (r *inventoryRepositoryGorm) GetStockTake(id uuid.UUID) (*models.StockTake, error) { return nil, ErrNotImplemented }
func (r *inventoryRepositoryGorm) ListStockTakes(companyID uuid.UUID, params map[string]interface{}) ([]models.StockTake, error) { return nil, ErrNotImplemented }
func (r *inventoryRepositoryGorm) CreateStockTakeItem(item *models.StockTakeItem) error { return ErrNotImplemented }
func (r *inventoryRepositoryGorm) UpdateStockTakeItem(item *models.StockTakeItem) error { return ErrNotImplemented }
func (r *inventoryRepositoryGorm) GetStockTakeItems(stockTakeID uuid.UUID) ([]models.StockTakeItem, error) { return nil, ErrNotImplemented }