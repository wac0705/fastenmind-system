package repository

import (
	"context"
	"fmt"

	"github.com/fastenmind/fastener-api/internal/domain/models"
	"github.com/fastenmind/fastener-api/internal/infrastructure/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BaseRepository provides common database operations with read/write separation
type BaseRepository struct {
	rw *database.ReadWriteDB
}

// NewBaseRepository creates a new base repository
func NewBaseRepository(rw *database.ReadWriteDB) *BaseRepository {
	return &BaseRepository{rw: rw}
}

// GetDB returns appropriate database connection based on options
func (r *BaseRepository) GetDB(opts ...database.QueryOptions) *gorm.DB {
	var queryOpts database.QueryOptions
	if len(opts) > 0 {
		queryOpts = opts[0]
	}

	// Determine which DB to use
	var db *gorm.DB
	if queryOpts.UseWriteDB || queryOpts.LockForUpdate || queryOpts.LockForShare {
		db = r.rw.Write()
	} else if queryOpts.ConsistentRead {
		db = r.rw.ReadPreferPrimary()
	} else {
		db = r.rw.Read()
	}

	// Apply query options
	return database.ApplyQueryOptions(db, queryOpts)
}

// Create creates a new record using write DB
func (r *BaseRepository) Create(ctx context.Context, value interface{}) error {
	return r.rw.Write().WithContext(ctx).Create(value).Error
}

// Update updates a record using write DB
func (r *BaseRepository) Update(ctx context.Context, value interface{}) error {
	return r.rw.Write().WithContext(ctx).Save(value).Error
}

// Delete soft deletes a record using write DB
func (r *BaseRepository) Delete(ctx context.Context, value interface{}) error {
	return r.rw.Write().WithContext(ctx).Delete(value).Error
}

// FindByID finds a record by ID using read DB
func (r *BaseRepository) FindByID(ctx context.Context, id uuid.UUID, dest interface{}, opts ...database.QueryOptions) error {
	db := r.GetDB(opts...).WithContext(ctx)
	err := db.First(dest, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}

// FindAll finds all records with pagination using read DB
func (r *BaseRepository) FindAll(ctx context.Context, dest interface{}, offset, limit int, opts ...database.QueryOptions) error {
	db := r.GetDB(opts...).WithContext(ctx)
	return db.Offset(offset).Limit(limit).Find(dest).Error
}

// Count counts records using read DB
func (r *BaseRepository) Count(ctx context.Context, model interface{}, opts ...database.QueryOptions) (int64, error) {
	var count int64
	db := r.GetDB(opts...).WithContext(ctx)
	err := db.Model(model).Count(&count).Error
	return count, err
}

// Transaction executes a function within a transaction on write DB
func (r *BaseRepository) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return r.rw.Write().WithContext(ctx).Transaction(fn)
}

// Exists checks if a record exists using read DB
func (r *BaseRepository) Exists(ctx context.Context, model interface{}, conditions ...interface{}) (bool, error) {
	var count int64
	db := r.GetDB().WithContext(ctx)
	err := db.Model(model).Where(conditions[0], conditions[1:]...).Count(&count).Error
	return count > 0, err
}

// InquiryRepositoryRW implements InquiryRepository with read/write separation
type InquiryRepositoryRW struct {
	*BaseRepository
}

// NewInquiryRepositoryRW creates a new inquiry repository with R/W separation
func NewInquiryRepositoryRW(rw *database.ReadWriteDB) *InquiryRepositoryRW {
	return &InquiryRepositoryRW{
		BaseRepository: NewBaseRepository(rw),
	}
}

// Create creates a new inquiry
func (r *InquiryRepositoryRW) Create(ctx context.Context, inquiry *models.Inquiry) error {
	return r.BaseRepository.Create(ctx, inquiry)
}

// GetByID retrieves an inquiry by ID
func (r *InquiryRepositoryRW) GetByID(ctx context.Context, id uuid.UUID) (*models.Inquiry, error) {
	var inquiry models.Inquiry
	err := r.FindByID(ctx, id, &inquiry)
	if err != nil {
		return nil, err
	}
	return &inquiry, nil
}

// Update updates an inquiry
func (r *InquiryRepositoryRW) Update(ctx context.Context, inquiry *models.Inquiry) error {
	return r.BaseRepository.Update(ctx, inquiry)
}

// List lists inquiries with pagination
func (r *InquiryRepositoryRW) List(ctx context.Context, params ListInquiriesParams) ([]*models.Inquiry, int64, error) {
	db := r.GetDB().WithContext(ctx)
	
	// Apply filters
	if params.CompanyID != uuid.Nil {
		db = db.Where("company_id = ?", params.CompanyID)
	}
	if params.CustomerID != nil && *params.CustomerID != uuid.Nil {
		db = db.Where("customer_id = ?", *params.CustomerID)
	}
	if params.Status != nil {
		db = db.Where("status = ?", *params.Status)
	}
	if params.AssignedEngineerID != nil && *params.AssignedEngineerID != uuid.Nil {
		db = db.Where("assigned_engineer_id = ?", *params.AssignedEngineerID)
	}
	if params.StartDate != nil {
		db = db.Where("created_at >= ?", *params.StartDate)
	}
	if params.EndDate != nil {
		db = db.Where("created_at <= ?", *params.EndDate)
	}
	
	// Count total
	var total int64
	if err := db.Model(&models.Inquiry{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Fetch with pagination
	var inquiries []*models.Inquiry
	err := db.Offset(params.Offset).
		Limit(params.Limit).
		Order("created_at DESC").
		Find(&inquiries).Error
		
	return inquiries, total, err
}

// GetByInquiryNo retrieves an inquiry by inquiry number
func (r *InquiryRepositoryRW) GetByInquiryNo(ctx context.Context, inquiryNo string) (*models.Inquiry, error) {
	var inquiry models.Inquiry
	db := r.GetDB().WithContext(ctx)
	err := db.Where("inquiry_no = ?", inquiryNo).First(&inquiry).Error
	if err == gorm.ErrRecordNotFound {
		return nil, ErrNotFound
	}
	return &inquiry, err
}

// ExistsByInquiryNo checks if an inquiry exists by inquiry number
func (r *InquiryRepositoryRW) ExistsByInquiryNo(ctx context.Context, inquiryNo string) (bool, error) {
	return r.Exists(ctx, &models.Inquiry{}, "inquiry_no = ?", inquiryNo)
}

// QuoteRepositoryRW implements QuoteRepository with read/write separation
type QuoteRepositoryRW struct {
	*BaseRepository
}

// NewQuoteRepositoryRW creates a new quote repository with R/W separation
func NewQuoteRepositoryRW(rw *database.ReadWriteDB) *QuoteRepositoryRW {
	return &QuoteRepositoryRW{
		BaseRepository: NewBaseRepository(rw),
	}
}

// Create creates a new quote
func (r *QuoteRepositoryRW) Create(ctx context.Context, quote *models.Quote) error {
	return r.BaseRepository.Create(ctx, quote)
}

// GetByID retrieves a quote by ID
func (r *QuoteRepositoryRW) GetByID(ctx context.Context, id uuid.UUID) (*models.Quote, error) {
	var quote models.Quote
	// Use consistent read for quotes to ensure data consistency
	err := r.FindByID(ctx, id, &quote, database.QueryOptions{ConsistentRead: true})
	if err != nil {
		return nil, err
	}
	return &quote, nil
}

// Update updates a quote
func (r *QuoteRepositoryRW) Update(ctx context.Context, quote *models.Quote) error {
	return r.BaseRepository.Update(ctx, quote)
}

// CreateWithItems creates a quote with items in a transaction
func (r *QuoteRepositoryRW) CreateWithItems(ctx context.Context, quote *models.Quote, items []models.QuoteItem) error {
	return r.Transaction(ctx, func(tx *gorm.DB) error {
		// Create quote
		if err := tx.Create(quote).Error; err != nil {
			return err
		}
		
		// Create items
		for i := range items {
			items[i].QuoteID = quote.ID
			if err := tx.Create(&items[i]).Error; err != nil {
				return err
			}
		}
		
		// Update quote with items
		quote.Items = items
		
		return nil
	})
}

// List lists quotes with pagination
func (r *QuoteRepositoryRW) List(ctx context.Context, params ListQuotesParams) ([]*models.Quote, int64, error) {
	db := r.GetDB().WithContext(ctx)
	
	// Apply filters
	if params.CompanyID != uuid.Nil {
		db = db.Where("company_id = ?", params.CompanyID)
	}
	if params.CustomerID != nil && *params.CustomerID != uuid.Nil {
		db = db.Where("customer_id = ?", *params.CustomerID)
	}
	if params.Status != nil {
		db = db.Where("status = ?", *params.Status)
	}
	if params.PreparedBy != nil && *params.PreparedBy != uuid.Nil {
		db = db.Where("prepared_by = ?", *params.PreparedBy)
	}
	if params.StartDate != nil {
		db = db.Where("created_at >= ?", *params.StartDate)
	}
	if params.EndDate != nil {
		db = db.Where("created_at <= ?", *params.EndDate)
	}
	
	// Count total
	var total int64
	if err := db.Model(&models.Quote{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Fetch with pagination
	var quotes []*models.Quote
	err := db.Offset(params.Offset).
		Limit(params.Limit).
		Order("created_at DESC").
		Preload("Items").
		Find(&quotes).Error
		
	return quotes, total, err
}

// GetByQuoteNo retrieves a quote by quote number
func (r *QuoteRepositoryRW) GetByQuoteNo(ctx context.Context, quoteNo string) (*models.Quote, error) {
	var quote models.Quote
	db := r.GetDB(database.QueryOptions{ConsistentRead: true}).WithContext(ctx)
	err := db.Where("quote_no = ?", quoteNo).Preload("Items").First(&quote).Error
	if err == gorm.ErrRecordNotFound {
		return nil, ErrNotFound
	}
	return &quote, err
}

// ExistsByQuoteNo checks if a quote exists by quote number
func (r *QuoteRepositoryRW) ExistsByQuoteNo(ctx context.Context, quoteNo string) (bool, error) {
	return r.Exists(ctx, &models.Quote{}, "quote_no = ?", quoteNo)
}

// OrderRepositoryRW implements OrderRepository with read/write separation
type OrderRepositoryRW struct {
	*BaseRepository
}

// NewOrderRepositoryRW creates a new order repository with R/W separation
func NewOrderRepositoryRW(rw *database.ReadWriteDB) *OrderRepositoryRW {
	return &OrderRepositoryRW{
		BaseRepository: NewBaseRepository(rw),
	}
}

// Create creates a new order
func (r *OrderRepositoryRW) Create(ctx context.Context, order *models.Order) error {
	return r.BaseRepository.Create(ctx, order)
}

// GetByID retrieves an order by ID
func (r *OrderRepositoryRW) GetByID(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	var order models.Order
	// Use consistent read for orders
	err := r.FindByID(ctx, id, &order, database.QueryOptions{ConsistentRead: true})
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// Update updates an order
func (r *OrderRepositoryRW) Update(ctx context.Context, order *models.Order) error {
	return r.BaseRepository.Update(ctx, order)
}

// CreateWithItems creates an order with items in a transaction
func (r *OrderRepositoryRW) CreateWithItems(ctx context.Context, order *models.Order, items []models.OrderItem) error {
	return r.Transaction(ctx, func(tx *gorm.DB) error {
		// Create order
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		
		// Create items
		for i := range items {
			items[i].OrderID = order.ID
			if err := tx.Create(&items[i]).Error; err != nil {
				return err
			}
		}
		
		// Update order with items
		order.Items = items
		
		return nil
	})
}

// UpdateStatus updates order status
func (r *OrderRepositoryRW) UpdateStatus(ctx context.Context, orderID uuid.UUID, status models.OrderStatus, updatedBy uuid.UUID) error {
	return r.rw.Write().WithContext(ctx).
		Model(&models.Order{}).
		Where("id = ?", orderID).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_by": updatedBy,
		}).Error
}

// List lists orders with pagination
func (r *OrderRepositoryRW) List(ctx context.Context, params ListOrdersParams) ([]*models.Order, int64, error) {
	db := r.GetDB().WithContext(ctx)
	
	// Apply filters
	if params.CompanyID != uuid.Nil {
		db = db.Where("company_id = ?", params.CompanyID)
	}
	if params.CustomerID != nil && *params.CustomerID != uuid.Nil {
		db = db.Where("customer_id = ?", *params.CustomerID)
	}
	if params.Status != nil {
		db = db.Where("status = ?", *params.Status)
	}
	if params.StartDate != nil {
		db = db.Where("created_at >= ?", *params.StartDate)
	}
	if params.EndDate != nil {
		db = db.Where("created_at <= ?", *params.EndDate)
	}
	
	// Count total
	var total int64
	if err := db.Model(&models.Order{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Fetch with pagination
	var orders []*models.Order
	err := db.Offset(params.Offset).
		Limit(params.Limit).
		Order("created_at DESC").
		Preload("Items").
		Find(&orders).Error
		
	return orders, total, err
}

// GetByOrderNo retrieves an order by order number
func (r *OrderRepositoryRW) GetByOrderNo(ctx context.Context, orderNo string) (*models.Order, error) {
	var order models.Order
	db := r.GetDB(database.QueryOptions{ConsistentRead: true}).WithContext(ctx)
	err := db.Where("order_no = ?", orderNo).Preload("Items").First(&order).Error
	if err == gorm.ErrRecordNotFound {
		return nil, ErrNotFound
	}
	return &order, err
}

// ExistsByOrderNo checks if an order exists by order number
func (r *OrderRepositoryRW) ExistsByOrderNo(ctx context.Context, orderNo string) (bool, error) {
	return r.Exists(ctx, &models.Order{}, "order_no = ?", orderNo)
}

// BatchOperations provides batch operations with read/write separation
type BatchOperations struct {
	rw *database.ReadWriteDB
}

// NewBatchOperations creates a new batch operations instance
func NewBatchOperations(rw *database.ReadWriteDB) *BatchOperations {
	return &BatchOperations{rw: rw}
}

// BatchCreate creates multiple records in a transaction
func (b *BatchOperations) BatchCreate(ctx context.Context, records interface{}) error {
	return b.rw.Write().WithContext(ctx).CreateInBatches(records, 100).Error
}

// BatchUpdate updates multiple records
func (b *BatchOperations) BatchUpdate(ctx context.Context, model interface{}, updates map[string]interface{}, conditions ...interface{}) error {
	db := b.rw.Write().WithContext(ctx).Model(model)
	if len(conditions) > 0 {
		db = db.Where(conditions[0], conditions[1:]...)
	}
	return db.Updates(updates).Error
}

// BatchDelete soft deletes multiple records
func (b *BatchOperations) BatchDelete(ctx context.Context, model interface{}, conditions ...interface{}) error {
	db := b.rw.Write().WithContext(ctx)
	if len(conditions) > 0 {
		db = db.Where(conditions[0], conditions[1:]...)
	}
	return db.Delete(model).Error
}