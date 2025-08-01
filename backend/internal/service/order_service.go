package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/fastenmind/fastener-api/internal/models"
	"github.com/fastenmind/fastener-api/internal/repository"
	"github.com/google/uuid"
)

type OrderService interface {
	List(companyID uuid.UUID, params map[string]interface{}) ([]models.Order, int64, error)
	Get(id uuid.UUID) (*models.Order, error)
	CreateFromQuote(companyID, userID, quoteID uuid.UUID, req CreateOrderRequest) (*models.Order, error)
	Update(id, userID uuid.UUID, req UpdateOrderRequest) (*models.Order, error)
	UpdateStatus(id, userID uuid.UUID, status string, notes string) (*models.Order, error)
	Delete(id uuid.UUID) error
	
	// Order items
	GetItems(orderID uuid.UUID) ([]models.OrderItem, error)
	UpdateItems(orderID uuid.UUID, items []OrderItemRequest) error
	
	// Documents
	AddDocument(orderID, userID uuid.UUID, req AddDocumentRequest) (*models.OrderDocument, error)
	RemoveDocument(docID uuid.UUID) error
	GetDocuments(orderID uuid.UUID) ([]models.OrderDocument, error)
	
	// Activities
	GetActivities(orderID uuid.UUID) ([]models.OrderActivity, error)
	
	// Reports
	GetOrderStats(companyID uuid.UUID, params map[string]interface{}) (*OrderStats, error)
}

type CreateOrderRequest struct {
	QuoteID         uuid.UUID `json:"quote_id" validate:"required"`
	PONumber        string    `json:"po_number" validate:"required"`
	Quantity        int       `json:"quantity" validate:"required,min=1"`
	DeliveryMethod  string    `json:"delivery_method" validate:"required"`
	DeliveryDate    string    `json:"delivery_date" validate:"required"`
	ShippingAddress string    `json:"shipping_address" validate:"required"`
	PaymentTerms    string    `json:"payment_terms" validate:"required"`
	DownPayment     float64   `json:"down_payment"`
	Notes           string    `json:"notes"`
}

type UpdateOrderRequest struct {
	PONumber        string  `json:"po_number"`
	DeliveryMethod  string  `json:"delivery_method"`
	DeliveryDate    string  `json:"delivery_date"`
	ShippingAddress string  `json:"shipping_address"`
	PaymentTerms    string  `json:"payment_terms"`
	Notes           string  `json:"notes"`
	InternalNotes   string  `json:"internal_notes"`
}

type OrderItemRequest struct {
	PartNo          string  `json:"part_no" validate:"required"`
	Description     string  `json:"description"`
	Quantity        int     `json:"quantity" validate:"required,min=1"`
	UnitPrice       float64 `json:"unit_price" validate:"required,min=0"`
	Material        string  `json:"material"`
	SurfaceTreatment string `json:"surface_treatment"`
	HeatTreatment   string  `json:"heat_treatment"`
	Specifications  string  `json:"specifications"`
}

type AddDocumentRequest struct {
	DocumentType string `json:"document_type" validate:"required"`
	FileName     string `json:"file_name" validate:"required"`
	FilePath     string `json:"file_path" validate:"required"`
	FileSize     int64  `json:"file_size"`
}

type OrderStats struct {
	TotalOrders     int     `json:"total_orders"`
	PendingOrders   int     `json:"pending_orders"`
	InProduction    int     `json:"in_production"`
	CompletedOrders int     `json:"completed_orders"`
	TotalRevenue    float64 `json:"total_revenue"`
	AvgOrderValue   float64 `json:"avg_order_value"`
}

type orderService struct {
	orderRepo    repository.OrderRepository
	quoteRepo    repository.QuoteRepository
	customerRepo repository.CustomerRepository
	n8nService   N8NService
}

func NewOrderService(
	orderRepo repository.OrderRepository,
	quoteRepo repository.QuoteRepository,
	customerRepo repository.CustomerRepository,
	n8nService N8NService,
) OrderService {
	return &orderService{
		orderRepo:    orderRepo,
		quoteRepo:    quoteRepo,
		customerRepo: customerRepo,
		n8nService:   n8nService,
	}
}

func (s *orderService) List(companyID uuid.UUID, params map[string]interface{}) ([]models.Order, int64, error) {
	return s.orderRepo.List(companyID, params)
}

func (s *orderService) Get(id uuid.UUID) (*models.Order, error) {
	return s.orderRepo.GetWithDetails(id)
}

func (s *orderService) CreateFromQuote(companyID, userID, quoteID uuid.UUID, req CreateOrderRequest) (*models.Order, error) {
	// Get quote details
	quote, err := s.quoteRepo.GetWithDetails(quoteID)
	if err != nil {
		return nil, errors.New("quote not found")
	}
	
	if quote.Status != "sent" && quote.Status != "accepted" {
		return nil, errors.New("quote must be sent or accepted to create order")
	}
	
	// Parse delivery date
	deliveryDate, err := time.Parse("2006-01-02", req.DeliveryDate)
	if err != nil {
		return nil, errors.New("invalid delivery date format")
	}
	
	// Calculate total amount
	totalAmount := quote.UnitPrice * float64(req.Quantity)
	
	// Create order
	order := &models.Order{
		CompanyID:       companyID,
		QuoteID:         quoteID,
		CustomerID:      quote.CustomerID,
		SalesID:         quote.SalesID,
		Status:          "pending",
		PONumber:        req.PONumber,
		Quantity:        req.Quantity,
		UnitPrice:       quote.UnitPrice,
		TotalAmount:     totalAmount,
		Currency:        quote.Currency,
		DeliveryMethod:  req.DeliveryMethod,
		DeliveryDate:    deliveryDate,
		ShippingAddress: req.ShippingAddress,
		PaymentTerms:    req.PaymentTerms,
		PaymentStatus:   "pending",
		DownPayment:     req.DownPayment,
		PaidAmount:      0,
		Notes:           req.Notes,
	}
	
	// Generate order number
	order.OrderNo = s.generateOrderNo(companyID)
	
	if err := s.orderRepo.Create(order); err != nil {
		return nil, err
	}
	
	// Create order item from inquiry
	if quote.Inquiry != nil {
		item := &models.OrderItem{
			OrderID:          order.ID,
			PartNo:           quote.Inquiry.PartNo,
			Description:      quote.Inquiry.Description,
			Quantity:         float64(req.Quantity),
			UnitPrice:        quote.UnitPrice,
			TotalPrice:       totalAmount,
			Material:         quote.Inquiry.Material,
			SurfaceTreatment: quote.Inquiry.SurfaceTreatment,
			HeatTreatment:    quote.Inquiry.HeatTreatment,
			Specifications:   quote.Inquiry.Specifications,
		}
		s.orderRepo.CreateItem(item)
	}
	
	// Update quote status
	quote.Status = "accepted"
	s.quoteRepo.Update(quote)
	
	// Log activity
	s.logActivity(order.ID, userID, "created", fmt.Sprintf("Order created from quote %s", quote.QuoteNo))
	
	// Trigger N8N workflow
	go s.n8nService.LogEvent(companyID, userID, "order.created", "order", order.ID, map[string]interface{}{
		"order_no":     order.OrderNo,
		"quote_no":     quote.QuoteNo,
		"customer_id":  order.CustomerID,
		"total_amount": order.TotalAmount,
		"po_number":    order.PONumber,
	})
	
	return s.Get(order.ID)
}

func (s *orderService) Update(id, userID uuid.UUID, req UpdateOrderRequest) (*models.Order, error) {
	order, err := s.orderRepo.Get(id)
	if err != nil {
		return nil, err
	}
	
	// Update fields
	if req.PONumber != "" {
		order.PONumber = req.PONumber
	}
	if req.DeliveryMethod != "" {
		order.DeliveryMethod = req.DeliveryMethod
	}
	if req.DeliveryDate != "" {
		deliveryDate, err := time.Parse("2006-01-02", req.DeliveryDate)
		if err != nil {
			return nil, errors.New("invalid delivery date format")
		}
		order.DeliveryDate = deliveryDate
	}
	if req.ShippingAddress != "" {
		order.ShippingAddress = req.ShippingAddress
	}
	if req.PaymentTerms != "" {
		order.PaymentTerms = req.PaymentTerms
	}
	order.Notes = req.Notes
	order.InternalNotes = req.InternalNotes
	
	if err := s.orderRepo.Update(order); err != nil {
		return nil, err
	}
	
	// Log activity
	s.logActivity(order.ID, userID, "updated", "Order details updated")
	
	return s.Get(order.ID)
}

func (s *orderService) UpdateStatus(id, userID uuid.UUID, status string, notes string) (*models.Order, error) {
	order, err := s.orderRepo.Get(id)
	if err != nil {
		return nil, err
	}
	
	// Validate status transition
	if !s.isValidStatusTransition(order.Status, status) {
		return nil, fmt.Errorf("invalid status transition from %s to %s", order.Status, status)
	}
	
	// Update status and timestamps
	previousStatus := order.Status
	order.Status = status
	now := time.Now()
	
	switch status {
	case "confirmed":
		order.ConfirmedAt = &now
	case "in_production":
		order.InProductionAt = &now
	case "quality_check":
		order.QualityCheckAt = &now
	case "ready_to_ship":
		order.ReadyToShipAt = &now
	case "shipped":
		order.ShippedAt = &now
	case "delivered":
		order.DeliveredAt = &now
	case "completed":
		order.CompletedAt = &now
	case "cancelled":
		order.CancelledAt = &now
	}
	
	if err := s.orderRepo.Update(order); err != nil {
		return nil, err
	}
	
	// Log activity
	description := fmt.Sprintf("Status changed from %s to %s", previousStatus, status)
	if notes != "" {
		description += ": " + notes
	}
	s.logActivity(order.ID, userID, "status_changed", description)
	
	// Trigger N8N workflow
	go s.n8nService.LogEvent(order.CompanyID, userID, fmt.Sprintf("order.%s", status), "order", order.ID, map[string]interface{}{
		"order_no":        order.OrderNo,
		"previous_status": previousStatus,
		"new_status":      status,
		"notes":           notes,
	})
	
	return s.Get(order.ID)
}

func (s *orderService) Delete(id uuid.UUID) error {
	order, err := s.orderRepo.Get(id)
	if err != nil {
		return err
	}
	
	if order.Status != "pending" && order.Status != "cancelled" {
		return errors.New("can only delete pending or cancelled orders")
	}
	
	return s.orderRepo.Delete(id)
}

// Order items
func (s *orderService) GetItems(orderID uuid.UUID) ([]models.OrderItem, error) {
	return s.orderRepo.GetItems(orderID)
}

func (s *orderService) UpdateItems(orderID uuid.UUID, items []OrderItemRequest) error {
	// Delete existing items
	existingItems, _ := s.orderRepo.GetItems(orderID)
	for _, item := range existingItems {
		s.orderRepo.DeleteItem(item.ID)
	}
	
	// Create new items
	totalAmount := 0.0
	for _, req := range items {
		totalPrice := float64(req.Quantity) * req.UnitPrice
		totalAmount += totalPrice
		
		item := &models.OrderItem{
			OrderID:          orderID,
			PartNo:           req.PartNo,
			Description:      req.Description,
			Quantity:         float64(req.Quantity),
			UnitPrice:        req.UnitPrice,
			TotalPrice:       totalPrice,
			Material:         req.Material,
			SurfaceTreatment: req.SurfaceTreatment,
			HeatTreatment:    req.HeatTreatment,
			Specifications:   req.Specifications,
		}
		
		if err := s.orderRepo.CreateItem(item); err != nil {
			return err
		}
	}
	
	// Update order total
	order, _ := s.orderRepo.Get(orderID)
	if order != nil {
		order.TotalAmount = totalAmount
		s.orderRepo.Update(order)
	}
	
	return nil
}

// Documents
func (s *orderService) AddDocument(orderID, userID uuid.UUID, req AddDocumentRequest) (*models.OrderDocument, error) {
	doc := &models.OrderDocument{
		OrderID:      orderID,
		DocumentType: req.DocumentType,
		FileName:     req.FileName,
		FilePath:     req.FilePath,
		FileSize:     req.FileSize,
		UploadedBy:   userID,
	}
	
	if err := s.orderRepo.AddDocument(doc); err != nil {
		return nil, err
	}
	
	// Log activity
	s.logActivity(orderID, userID, "document_added", fmt.Sprintf("Added %s: %s", req.DocumentType, req.FileName))
	
	return doc, nil
}

func (s *orderService) RemoveDocument(docID uuid.UUID) error {
	return s.orderRepo.RemoveDocument(docID)
}

func (s *orderService) GetDocuments(orderID uuid.UUID) ([]models.OrderDocument, error) {
	return s.orderRepo.GetDocuments(orderID)
}

// Activities
func (s *orderService) GetActivities(orderID uuid.UUID) ([]models.OrderActivity, error) {
	return s.orderRepo.GetActivities(orderID)
}

// Reports
func (s *orderService) GetOrderStats(companyID uuid.UUID, params map[string]interface{}) (*OrderStats, error) {
	// Get all orders for stats
	allOrders, _, err := s.orderRepo.List(companyID, map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	
	stats := &OrderStats{
		TotalOrders: len(allOrders),
	}
	
	totalRevenue := 0.0
	for _, order := range allOrders {
		switch order.Status {
		case "pending":
			stats.PendingOrders++
		case "in_production":
			stats.InProduction++
		case "completed":
			stats.CompletedOrders++
		}
		
		if order.Status != "cancelled" {
			totalRevenue += order.TotalAmount
		}
	}
	
	stats.TotalRevenue = totalRevenue
	if stats.TotalOrders > 0 {
		stats.AvgOrderValue = totalRevenue / float64(stats.TotalOrders)
	}
	
	return stats, nil
}

func (s *orderService) generateOrderNo(companyID uuid.UUID) string {
	// TODO: Implement proper order number generation with sequence
	return fmt.Sprintf("PO-%s", time.Now().Format("20060102-150405"))
}

func (s *orderService) isValidStatusTransition(from, to string) bool {
	validTransitions := map[string][]string{
		"pending":       {"confirmed", "cancelled"},
		"confirmed":     {"in_production", "cancelled"},
		"in_production": {"quality_check", "cancelled"},
		"quality_check": {"ready_to_ship", "in_production", "cancelled"},
		"ready_to_ship": {"shipped", "cancelled"},
		"shipped":       {"delivered", "cancelled"},
		"delivered":     {"completed"},
		"completed":     {},
		"cancelled":     {},
	}
	
	allowedStatuses, ok := validTransitions[from]
	if !ok {
		return false
	}
	
	for _, status := range allowedStatuses {
		if status == to {
			return true
		}
	}
	
	return false
}

func (s *orderService) logActivity(orderID, userID uuid.UUID, action, description string) {
	activity := &models.OrderActivity{
		OrderID:     orderID,
		UserID:      userID,
		Action:      action,
		Description: description,
	}
	
	s.orderRepo.LogActivity(activity)
}