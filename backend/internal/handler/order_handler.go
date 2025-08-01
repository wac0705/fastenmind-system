package handler

import (
	"net/http"
	"strconv"

	"github.com/fastenmind/fastener-api/internal/middleware"
	"github.com/fastenmind/fastener-api/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type OrderHandler struct {
	service service.OrderService
}

func NewOrderHandler(service service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

// List godoc
// @Summary List orders
// @Description Get a list of orders with optional filters
// @Tags Orders
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param status query string false "Filter by status"
// @Param customer_id query string false "Filter by customer"
// @Param sales_id query string false "Filter by sales person"
// @Param payment_status query string false "Filter by payment status"
// @Param search query string false "Search term"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{}
// @Router /api/orders [get]
func (h *OrderHandler) List(c echo.Context) error {
	params := make(map[string]interface{})
	
	if page := c.QueryParam("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil {
			params["page"] = p
		}
	}
	if pageSize := c.QueryParam("page_size"); pageSize != "" {
		if ps, err := strconv.Atoi(pageSize); err == nil {
			params["page_size"] = ps
		}
	}
	if status := c.QueryParam("status"); status != "" {
		params["status"] = status
	}
	if customerID := c.QueryParam("customer_id"); customerID != "" {
		params["customer_id"] = customerID
	}
	if salesID := c.QueryParam("sales_id"); salesID != "" {
		params["sales_id"] = salesID
	}
	if paymentStatus := c.QueryParam("payment_status"); paymentStatus != "" {
		params["payment_status"] = paymentStatus
	}
	if search := c.QueryParam("search"); search != "" {
		params["search"] = search
	}
	if startDate := c.QueryParam("start_date"); startDate != "" {
		params["start_date"] = startDate
	}
	if endDate := c.QueryParam("end_date"); endDate != "" {
		params["end_date"] = endDate
	}
	
	userClaims := c.Get("user").(*middleware.Claims)
	
	orders, total, err := h.service.List(userClaims.CompanyID, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": orders,
		"pagination": map[string]interface{}{
			"total":     total,
			"page":      params["page"],
			"page_size": params["page_size"],
		},
	})
}

// Get godoc
// @Summary Get order details
// @Description Get detailed information about a specific order
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} models.Order
// @Router /api/orders/{id} [get]
func (h *OrderHandler) Get(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid order ID")
	}
	
	order, err := h.service.Get(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Order not found")
	}
	
	return c.JSON(http.StatusOK, order)
}

// CreateFromQuote godoc
// @Summary Create order from quote
// @Description Create a new order from an accepted quote
// @Tags Orders
// @Accept json
// @Produce json
// @Param request body service.CreateOrderRequest true "Create order request"
// @Success 201 {object} models.Order
// @Router /api/orders [post]
func (h *OrderHandler) CreateFromQuote(c echo.Context) error {
	var req service.CreateOrderRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	userClaims := c.Get("user").(*middleware.Claims)
	
	order, err := h.service.CreateFromQuote(userClaims.CompanyID, userClaims.UserID, req.QuoteID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusCreated, order)
}

// Update godoc
// @Summary Update order
// @Description Update order details
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param request body service.UpdateOrderRequest true "Update order request"
// @Success 200 {object} models.Order
// @Router /api/orders/{id} [put]
func (h *OrderHandler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid order ID")
	}
	
	var req service.UpdateOrderRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	userClaims := c.Get("user").(*middleware.Claims)
	
	order, err := h.service.Update(id, userClaims.UserID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusOK, order)
}

// UpdateStatus godoc
// @Summary Update order status
// @Description Update the status of an order
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param request body map[string]string true "Status update request"
// @Success 200 {object} models.Order
// @Router /api/orders/{id}/status [put]
func (h *OrderHandler) UpdateStatus(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid order ID")
	}
	
	var req struct {
		Status string `json:"status" validate:"required"`
		Notes  string `json:"notes"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	userClaims := c.Get("user").(*middleware.Claims)
	
	order, err := h.service.UpdateStatus(id, userClaims.UserID, req.Status, req.Notes)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusOK, order)
}

// Delete godoc
// @Summary Delete order
// @Description Delete an order (soft delete)
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 204
// @Router /api/orders/{id} [delete]
func (h *OrderHandler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid order ID")
	}
	
	if err := h.service.Delete(id); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	return c.NoContent(http.StatusNoContent)
}

// GetItems godoc
// @Summary Get order items
// @Description Get items for an order
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {array} models.OrderItem
// @Router /api/orders/{id}/items [get]
func (h *OrderHandler) GetItems(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid order ID")
	}
	
	items, err := h.service.GetItems(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, items)
}

// UpdateItems godoc
// @Summary Update order items
// @Description Update items for an order
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param request body []service.OrderItemRequest true "Order items"
// @Success 200 {object} map[string]string
// @Router /api/orders/{id}/items [put]
func (h *OrderHandler) UpdateItems(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid order ID")
	}
	
	var items []service.OrderItemRequest
	if err := c.Bind(&items); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	if err := h.service.UpdateItems(id, items); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusOK, map[string]string{"message": "Items updated successfully"})
}

// GetDocuments godoc
// @Summary Get order documents
// @Description Get documents for an order
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {array} models.OrderDocument
// @Router /api/orders/{id}/documents [get]
func (h *OrderHandler) GetDocuments(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid order ID")
	}
	
	documents, err := h.service.GetDocuments(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, documents)
}

// AddDocument godoc
// @Summary Add document to order
// @Description Add a document to an order
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param request body service.AddDocumentRequest true "Document details"
// @Success 201 {object} models.OrderDocument
// @Router /api/orders/{id}/documents [post]
func (h *OrderHandler) AddDocument(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid order ID")
	}
	
	var req service.AddDocumentRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	userClaims := c.Get("user").(*middleware.Claims)
	
	doc, err := h.service.AddDocument(id, userClaims.UserID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusCreated, doc)
}

// RemoveDocument godoc
// @Summary Remove document from order
// @Description Remove a document from an order
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param doc_id path string true "Document ID"
// @Success 204
// @Router /api/orders/{id}/documents/{doc_id} [delete]
func (h *OrderHandler) RemoveDocument(c echo.Context) error {
	docID, err := uuid.Parse(c.Param("doc_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid document ID")
	}
	
	if err := h.service.RemoveDocument(docID); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	return c.NoContent(http.StatusNoContent)
}

// GetActivities godoc
// @Summary Get order activities
// @Description Get activity history for an order
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {array} models.OrderActivity
// @Router /api/orders/{id}/activities [get]
func (h *OrderHandler) GetActivities(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid order ID")
	}
	
	activities, err := h.service.GetActivities(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, activities)
}

// GetStats godoc
// @Summary Get order statistics
// @Description Get order statistics for the company
// @Tags Orders
// @Accept json
// @Produce json
// @Success 200 {object} service.OrderStats
// @Router /api/orders/stats [get]
func (h *OrderHandler) GetStats(c echo.Context) error {
	userClaims := c.Get("user").(*middleware.Claims)
	
	params := make(map[string]interface{})
	// Add any filter params from query string if needed
	
	stats, err := h.service.GetOrderStats(userClaims.CompanyID, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, stats)
}

// RegisterRoutes registers all order routes
func (h *OrderHandler) RegisterRoutes(e *echo.Echo, authMiddleware echo.MiddlewareFunc) {
	orders := e.Group("/api/orders", authMiddleware)
	
	orders.GET("", h.List)
	orders.POST("", h.CreateFromQuote)
	orders.GET("/stats", h.GetStats)
	orders.GET("/:id", h.Get)
	orders.PUT("/:id", h.Update)
	orders.DELETE("/:id", h.Delete)
	orders.PUT("/:id/status", h.UpdateStatus)
	orders.GET("/:id/items", h.GetItems)
	orders.PUT("/:id/items", h.UpdateItems)
	orders.GET("/:id/documents", h.GetDocuments)
	orders.POST("/:id/documents", h.AddDocument)
	orders.DELETE("/:id/documents/:doc_id", h.RemoveDocument)
	orders.GET("/:id/activities", h.GetActivities)
}