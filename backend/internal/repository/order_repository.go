package repository

import (
	"github.com/fastenmind/fastener-api/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *models.Order) error
	Update(order *models.Order) error
	UpdateOrder(order *models.Order) error  // Alias for Update
	Delete(id uuid.UUID) error
	Get(id uuid.UUID) (*models.Order, error)
	GetOrder(id uuid.UUID) (*models.Order, error)  // Alias for Get
	GetWithDetails(id uuid.UUID) (*models.Order, error)
	List(companyID uuid.UUID, params map[string]interface{}) ([]models.Order, int64, error)
	GetByOrderNo(orderNo string) (*models.Order, error)
	
	// Order items
	CreateItem(item *models.OrderItem) error
	UpdateItem(item *models.OrderItem) error
	DeleteItem(id uuid.UUID) error
	GetItems(orderID uuid.UUID) ([]models.OrderItem, error)
	GetOrderItems(orderID uuid.UUID) ([]models.OrderItem, error)  // Alias for GetItems
	
	// Activity log
	LogActivity(activity *models.OrderActivity) error
	GetActivities(orderID uuid.UUID) ([]models.OrderActivity, error)
	
	// Documents
	AddDocument(doc *models.OrderDocument) error
	RemoveDocument(id uuid.UUID) error
	GetDocuments(orderID uuid.UUID) ([]models.OrderDocument, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db interface{}) OrderRepository {
	gormDB, ok := db.(*gorm.DB)
	if !ok {
		panic("invalid database type, expected *gorm.DB")
	}
	return &orderRepository{db: gormDB}
}

func (r *orderRepository) Create(order *models.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) Update(order *models.Order) error {
	return r.db.Save(order).Error
}

func (r *orderRepository) UpdateOrder(order *models.Order) error {
	return r.Update(order)
}

func (r *orderRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Order{}, id).Error
}

func (r *orderRepository) Get(id uuid.UUID) (*models.Order, error) {
	var order models.Order
	if err := r.db.First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) GetOrder(id uuid.UUID) (*models.Order, error) {
	return r.Get(id)
}

func (r *orderRepository) GetWithDetails(id uuid.UUID) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("Quote").
		Preload("Quote.Inquiry").
		Preload("Customer").
		Preload("Sales").
		First(&order, id).Error
		
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) List(companyID uuid.UUID, params map[string]interface{}) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64
	
	query := r.db.Model(&models.Order{}).Where("company_id = ?", companyID)
	
	// Apply filters
	if status, ok := params["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	
	if customerID, ok := params["customer_id"].(string); ok && customerID != "" {
		query = query.Where("customer_id = ?", customerID)
	}
	
	if salesID, ok := params["sales_id"].(string); ok && salesID != "" {
		query = query.Where("sales_id = ?", salesID)
	}
	
	if paymentStatus, ok := params["payment_status"].(string); ok && paymentStatus != "" {
		query = query.Where("payment_status = ?", paymentStatus)
	}
	
	if search, ok := params["search"].(string); ok && search != "" {
		query = query.Where("order_no LIKE ? OR po_number LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	
	// Date range filters
	if startDate, ok := params["start_date"].(string); ok && startDate != "" {
		query = query.Where("created_at >= ?", startDate)
	}
	
	if endDate, ok := params["end_date"].(string); ok && endDate != "" {
		query = query.Where("created_at <= ?", endDate)
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
		Preload("Quote").
		Preload("Customer").
		Preload("Sales").
		Find(&orders).Error; err != nil {
		return nil, 0, err
	}
	
	return orders, total, nil
}

func (r *orderRepository) GetByOrderNo(orderNo string) (*models.Order, error) {
	var order models.Order
	if err := r.db.Where("order_no = ?", orderNo).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// Order items
func (r *orderRepository) CreateItem(item *models.OrderItem) error {
	return r.db.Create(item).Error
}

func (r *orderRepository) UpdateItem(item *models.OrderItem) error {
	return r.db.Save(item).Error
}

func (r *orderRepository) DeleteItem(id uuid.UUID) error {
	return r.db.Delete(&models.OrderItem{}, id).Error
}

func (r *orderRepository) GetItems(orderID uuid.UUID) ([]models.OrderItem, error) {
	var items []models.OrderItem
	err := r.db.Where("order_id = ?", orderID).Find(&items).Error
	return items, err
}

func (r *orderRepository) GetOrderItems(orderID uuid.UUID) ([]models.OrderItem, error) {
	return r.GetItems(orderID)
}

// Activity log
func (r *orderRepository) LogActivity(activity *models.OrderActivity) error {
	return r.db.Create(activity).Error
}

func (r *orderRepository) GetActivities(orderID uuid.UUID) ([]models.OrderActivity, error) {
	var activities []models.OrderActivity
	err := r.db.Where("order_id = ?", orderID).
		Preload("User").
		Order("created_at DESC").
		Find(&activities).Error
	return activities, err
}

// Documents
func (r *orderRepository) AddDocument(doc *models.OrderDocument) error {
	return r.db.Create(doc).Error
}

func (r *orderRepository) RemoveDocument(id uuid.UUID) error {
	return r.db.Delete(&models.OrderDocument{}, id).Error
}

func (r *orderRepository) GetDocuments(orderID uuid.UUID) ([]models.OrderDocument, error) {
	var docs []models.OrderDocument
	err := r.db.Where("order_id = ?", orderID).
		Preload("Uploader").
		Order("created_at DESC").
		Find(&docs).Error
	return docs, err
}