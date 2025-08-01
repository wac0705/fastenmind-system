package repository

import (
	"github.com/fastenmind/fastener-api/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SupplierRepository interface {
	// Supplier operations
	CreateSupplier(supplier *models.Supplier) error
	UpdateSupplier(supplier *models.Supplier) error
	GetSupplier(id uuid.UUID) (*models.Supplier, error)
	ListSuppliers(companyID uuid.UUID, params map[string]interface{}) ([]models.Supplier, int64, error)
	
	// Supplier Contact operations
	CreateSupplierContact(contact *models.SupplierContact) error
	UpdateSupplierContact(contact *models.SupplierContact) error
	GetSupplierContacts(supplierID uuid.UUID) ([]models.SupplierContact, error)
	DeleteSupplierContact(id uuid.UUID) error
	
	// Supplier Product operations
	CreateSupplierProduct(product *models.SupplierProduct) error
	UpdateSupplierProduct(product *models.SupplierProduct) error
	GetSupplierProducts(supplierID uuid.UUID, params map[string]interface{}) ([]models.SupplierProduct, error)
	DeleteSupplierProduct(id uuid.UUID) error
	
	// Purchase Order operations
	CreatePurchaseOrder(order *models.PurchaseOrder) error
	UpdatePurchaseOrder(order *models.PurchaseOrder) error
	GetPurchaseOrder(id uuid.UUID) (*models.PurchaseOrder, error)
	ListPurchaseOrders(companyID uuid.UUID, params map[string]interface{}) ([]models.PurchaseOrder, int64, error)
	
	// Purchase Order Item operations
	CreatePurchaseOrderItem(item *models.PurchaseOrderItem) error
	UpdatePurchaseOrderItem(item *models.PurchaseOrderItem) error
	GetPurchaseOrderItems(purchaseOrderID uuid.UUID) ([]models.PurchaseOrderItem, error)
	DeletePurchaseOrderItem(id uuid.UUID) error
	
	// Supplier Evaluation operations
	CreateSupplierEvaluation(evaluation *models.SupplierEvaluation) error
	UpdateSupplierEvaluation(evaluation *models.SupplierEvaluation) error
	GetSupplierEvaluation(id uuid.UUID) (*models.SupplierEvaluation, error)
	ListSupplierEvaluations(companyID uuid.UUID, params map[string]interface{}) ([]models.SupplierEvaluation, int64, error)
}

type supplierRepository struct {
	db *gorm.DB
}

func NewSupplierRepository(db interface{}) SupplierRepository {
	gormDB, ok := db.(*gorm.DB)
	if !ok {
		panic("invalid database type, expected *gorm.DB")
	}
	return &supplierRepository{db: gormDB}
}

// Supplier operations
func (r *supplierRepository) CreateSupplier(supplier *models.Supplier) error {
	return r.db.Create(supplier).Error
}

func (r *supplierRepository) UpdateSupplier(supplier *models.Supplier) error {
	return r.db.Save(supplier).Error
}

func (r *supplierRepository) GetSupplier(id uuid.UUID) (*models.Supplier, error) {
	var supplier models.Supplier
	err := r.db.Preload("Company").
		Preload("Creator").
		Preload("Contacts").
		Preload("Products").
		First(&supplier, id).Error
	if err != nil {
		return nil, err
	}
	return &supplier, nil
}

func (r *supplierRepository) ListSuppliers(companyID uuid.UUID, params map[string]interface{}) ([]models.Supplier, int64, error) {
	var suppliers []models.Supplier
	var total int64

	query := r.db.Model(&models.Supplier{}).Where("company_id = ?", companyID)

	// Apply filters
	if supplierType, ok := params["type"].(string); ok && supplierType != "" {
		query = query.Where("type = ?", supplierType)
	}

	if status, ok := params["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}

	if country, ok := params["country"].(string); ok && country != "" {
		query = query.Where("country = ?", country)
	}

	if riskLevel, ok := params["risk_level"].(string); ok && riskLevel != "" {
		query = query.Where("risk_level = ?", riskLevel)
	}

	if search, ok := params["search"].(string); ok && search != "" {
		query = query.Where("supplier_no LIKE ? OR name LIKE ? OR name_en LIKE ?", 
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
		Preload("Creator").
		Find(&suppliers).Error; err != nil {
		return nil, 0, err
	}

	return suppliers, total, nil
}

// Supplier Contact operations
func (r *supplierRepository) CreateSupplierContact(contact *models.SupplierContact) error {
	return r.db.Create(contact).Error
}

func (r *supplierRepository) UpdateSupplierContact(contact *models.SupplierContact) error {
	return r.db.Save(contact).Error
}

func (r *supplierRepository) GetSupplierContacts(supplierID uuid.UUID) ([]models.SupplierContact, error) {
	var contacts []models.SupplierContact
	err := r.db.Where("supplier_id = ? AND is_active = ?", supplierID, true).
		Order("is_primary DESC, name ASC").
		Find(&contacts).Error
	return contacts, err
}

func (r *supplierRepository) DeleteSupplierContact(id uuid.UUID) error {
	return r.db.Delete(&models.SupplierContact{}, id).Error
}

// Supplier Product operations
func (r *supplierRepository) CreateSupplierProduct(product *models.SupplierProduct) error {
	return r.db.Create(product).Error
}

func (r *supplierRepository) UpdateSupplierProduct(product *models.SupplierProduct) error {
	return r.db.Save(product).Error
}

func (r *supplierRepository) GetSupplierProducts(supplierID uuid.UUID, params map[string]interface{}) ([]models.SupplierProduct, error) {
	var products []models.SupplierProduct
	query := r.db.Where("supplier_id = ?", supplierID)

	// Apply filters
	if status, ok := params["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}

	if category, ok := params["category"].(string); ok && category != "" {
		query = query.Where("category = ?", category)
	}

	if preferred, ok := params["is_preferred"].(bool); ok {
		query = query.Where("is_preferred = ?", preferred)
	}

	if search, ok := params["search"].(string); ok && search != "" {
		query = query.Where("product_name LIKE ? OR product_code LIKE ? OR supplier_part_no LIKE ?", 
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Load with relations
	err := query.
		Preload("Inventory").
		Order("is_preferred DESC, product_name ASC").
		Find(&products).Error

	return products, err
}

func (r *supplierRepository) DeleteSupplierProduct(id uuid.UUID) error {
	return r.db.Delete(&models.SupplierProduct{}, id).Error
}

// Purchase Order operations
func (r *supplierRepository) CreatePurchaseOrder(order *models.PurchaseOrder) error {
	return r.db.Create(order).Error
}

func (r *supplierRepository) UpdatePurchaseOrder(order *models.PurchaseOrder) error {
	return r.db.Save(order).Error
}

func (r *supplierRepository) GetPurchaseOrder(id uuid.UUID) (*models.PurchaseOrder, error) {
	var order models.PurchaseOrder
	err := r.db.Preload("Company").
		Preload("Supplier").
		Preload("Creator").
		Preload("Approver").
		Preload("Items").
		Preload("Items.SupplierProduct").
		Preload("Items.Inventory").
		First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *supplierRepository) ListPurchaseOrders(companyID uuid.UUID, params map[string]interface{}) ([]models.PurchaseOrder, int64, error) {
	var orders []models.PurchaseOrder
	var total int64

	query := r.db.Model(&models.PurchaseOrder{}).Where("company_id = ?", companyID)

	// Apply filters
	if status, ok := params["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}

	if supplierID, ok := params["supplier_id"].(string); ok && supplierID != "" {
		query = query.Where("supplier_id = ?", supplierID)
	}

	if startDate, ok := params["start_date"].(string); ok && startDate != "" {
		query = query.Where("order_date >= ?", startDate)
	}

	if endDate, ok := params["end_date"].(string); ok && endDate != "" {
		query = query.Where("order_date <= ?", endDate)
	}

	if search, ok := params["search"].(string); ok && search != "" {
		query = query.Where("order_no LIKE ?", "%"+search+"%")
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
	query = query.Order("order_date DESC")

	// Load with relations
	if err := query.
		Preload("Supplier").
		Preload("Creator").
		Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// Purchase Order Item operations
func (r *supplierRepository) CreatePurchaseOrderItem(item *models.PurchaseOrderItem) error {
	return r.db.Create(item).Error
}

func (r *supplierRepository) UpdatePurchaseOrderItem(item *models.PurchaseOrderItem) error {
	return r.db.Save(item).Error
}

func (r *supplierRepository) GetPurchaseOrderItems(purchaseOrderID uuid.UUID) ([]models.PurchaseOrderItem, error) {
	var items []models.PurchaseOrderItem
	err := r.db.Where("purchase_order_id = ?", purchaseOrderID).
		Preload("SupplierProduct").
		Preload("Inventory").
		Find(&items).Error
	return items, err
}

func (r *supplierRepository) DeletePurchaseOrderItem(id uuid.UUID) error {
	return r.db.Delete(&models.PurchaseOrderItem{}, id).Error
}

// Supplier Evaluation operations
func (r *supplierRepository) CreateSupplierEvaluation(evaluation *models.SupplierEvaluation) error {
	return r.db.Create(evaluation).Error
}

func (r *supplierRepository) UpdateSupplierEvaluation(evaluation *models.SupplierEvaluation) error {
	return r.db.Save(evaluation).Error
}

func (r *supplierRepository) GetSupplierEvaluation(id uuid.UUID) (*models.SupplierEvaluation, error) {
	var evaluation models.SupplierEvaluation
	err := r.db.Preload("Company").
		Preload("Supplier").
		Preload("Evaluator").
		Preload("Approver").
		First(&evaluation, id).Error
	if err != nil {
		return nil, err
	}
	return &evaluation, nil
}

func (r *supplierRepository) ListSupplierEvaluations(companyID uuid.UUID, params map[string]interface{}) ([]models.SupplierEvaluation, int64, error) {
	var evaluations []models.SupplierEvaluation
	var total int64

	query := r.db.Model(&models.SupplierEvaluation{}).Where("company_id = ?", companyID)

	// Apply filters
	if status, ok := params["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}

	if supplierID, ok := params["supplier_id"].(string); ok && supplierID != "" {
		query = query.Where("supplier_id = ?", supplierID)
	}

	if evaluationType, ok := params["evaluation_type"].(string); ok && evaluationType != "" {
		query = query.Where("evaluation_type = ?", evaluationType)
	}

	if evaluatedBy, ok := params["evaluated_by"].(string); ok && evaluatedBy != "" {
		query = query.Where("evaluated_by = ?", evaluatedBy)
	}

	if startDate, ok := params["start_date"].(string); ok && startDate != "" {
		query = query.Where("evaluated_at >= ?", startDate)
	}

	if endDate, ok := params["end_date"].(string); ok && endDate != "" {
		query = query.Where("evaluated_at <= ?", endDate)
	}

	if search, ok := params["search"].(string); ok && search != "" {
		query = query.Where("evaluation_no LIKE ?", "%"+search+"%")
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
	query = query.Order("evaluated_at DESC")

	// Load with relations
	if err := query.
		Preload("Supplier").
		Preload("Evaluator").
		Find(&evaluations).Error; err != nil {
		return nil, 0, err
	}

	return evaluations, total, nil
}