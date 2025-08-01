package handler

import (
	"net/http"
	"strconv"

	"github.com/fastenmind/fastener-api/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type InventoryHandler struct {
	service service.InventoryService
}

func NewInventoryHandler(service service.InventoryService) *InventoryHandler {
	return &InventoryHandler{service: service}
}

// List godoc
// @Summary List inventory items
// @Description Get a list of inventory items with optional filters
// @Tags Inventory
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param category query string false "Filter by category"
// @Param warehouse_id query string false "Filter by warehouse"
// @Param status query string false "Filter by status"
// @Param low_stock query bool false "Filter low stock items"
// @Param search query string false "Search term"
// @Success 200 {object} map[string]interface{}
// @Router /api/inventory [get]
func (h *InventoryHandler) List(c echo.Context) error {
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
	if category := c.QueryParam("category"); category != "" {
		params["category"] = category
	}
	if warehouseID := c.QueryParam("warehouse_id"); warehouseID != "" {
		params["warehouse_id"] = warehouseID
	}
	if status := c.QueryParam("status"); status != "" {
		params["status"] = status
	}
	if lowStock := c.QueryParam("low_stock"); lowStock == "true" {
		params["low_stock"] = true
	}
	if search := c.QueryParam("search"); search != "" {
		params["search"] = search
	}
	
	companyID := getCompanyIDFromContext(c)
	
	items, total, err := h.service.List(companyID, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": items,
		"pagination": map[string]interface{}{
			"total":     total,
			"page":      params["page"],
			"page_size": params["page_size"],
		},
	})
}

// Get godoc
// @Summary Get inventory item details
// @Description Get detailed information about a specific inventory item
// @Tags Inventory
// @Accept json
// @Produce json
// @Param id path string true "Inventory ID"
// @Success 200 {object} models.Inventory
// @Router /api/inventory/{id} [get]
func (h *InventoryHandler) Get(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inventory ID")
	}
	
	item, err := h.service.Get(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Inventory item not found")
	}
	
	return c.JSON(http.StatusOK, item)
}

// GetBySKU godoc
// @Summary Get inventory item by SKU
// @Description Get inventory item by SKU
// @Tags Inventory
// @Accept json
// @Produce json
// @Param sku path string true "SKU"
// @Success 200 {object} models.Inventory
// @Router /api/inventory/sku/{sku} [get]
func (h *InventoryHandler) GetBySKU(c echo.Context) error {
	sku := c.Param("sku")
	
	item, err := h.service.GetBySKU(sku)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Inventory item not found")
	}
	
	return c.JSON(http.StatusOK, item)
}

// Create godoc
// @Summary Create inventory item
// @Description Create a new inventory item
// @Tags Inventory
// @Accept json
// @Produce json
// @Param request body service.CreateInventoryRequest true "Create inventory request"
// @Success 201 {object} models.Inventory
// @Router /api/inventory [post]
func (h *InventoryHandler) Create(c echo.Context) error {
	var req service.CreateInventoryRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	companyID := getCompanyIDFromContext(c)
	
	item, err := h.service.Create(companyID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusCreated, item)
}

// Update godoc
// @Summary Update inventory item
// @Description Update inventory item details
// @Tags Inventory
// @Accept json
// @Produce json
// @Param id path string true "Inventory ID"
// @Param request body service.UpdateInventoryRequest true "Update inventory request"
// @Success 200 {object} models.Inventory
// @Router /api/inventory/{id} [put]
func (h *InventoryHandler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inventory ID")
	}
	
	var req service.UpdateInventoryRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	item, err := h.service.Update(id, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusOK, item)
}

// Delete godoc
// @Summary Delete inventory item
// @Description Delete an inventory item (soft delete)
// @Tags Inventory
// @Accept json
// @Produce json
// @Param id path string true "Inventory ID"
// @Success 204
// @Router /api/inventory/{id} [delete]
func (h *InventoryHandler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inventory ID")
	}
	
	if err := h.service.Delete(id); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	return c.NoContent(http.StatusNoContent)
}

// AdjustStock godoc
// @Summary Adjust stock
// @Description Adjust inventory stock level
// @Tags Inventory
// @Accept json
// @Produce json
// @Param id path string true "Inventory ID"
// @Param request body service.StockAdjustmentRequest true "Stock adjustment request"
// @Success 200 {object} models.StockMovement
// @Router /api/inventory/{id}/adjust [post]
func (h *InventoryHandler) AdjustStock(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inventory ID")
	}
	
	var req service.StockAdjustmentRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	userID := getUserIDFromContext(c)
	
	movement, err := h.service.AdjustStock(id, userID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusOK, movement)
}

// TransferStock godoc
// @Summary Transfer stock
// @Description Transfer stock between warehouses
// @Tags Inventory
// @Accept json
// @Produce json
// @Param request body service.StockTransferRequest true "Stock transfer request"
// @Success 200 {object} models.StockMovement
// @Router /api/inventory/transfer [post]
func (h *InventoryHandler) TransferStock(c echo.Context) error {
	var req service.StockTransferRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	userID := getUserIDFromContext(c)
	
	movement, err := h.service.TransferStock(userID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusOK, movement)
}

// GetMovements godoc
// @Summary Get stock movements
// @Description Get stock movement history for an inventory item
// @Tags Inventory
// @Accept json
// @Produce json
// @Param id path string true "Inventory ID"
// @Param movement_type query string false "Filter by movement type"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {array} models.StockMovement
// @Router /api/inventory/{id}/movements [get]
func (h *InventoryHandler) GetMovements(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inventory ID")
	}
	
	params := make(map[string]interface{})
	if movementType := c.QueryParam("movement_type"); movementType != "" {
		params["movement_type"] = movementType
	}
	if startDate := c.QueryParam("start_date"); startDate != "" {
		params["start_date"] = startDate
	}
	if endDate := c.QueryParam("end_date"); endDate != "" {
		params["end_date"] = endDate
	}
	
	movements, err := h.service.GetMovements(id, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, movements)
}

// GetStats godoc
// @Summary Get inventory statistics
// @Description Get inventory statistics for the company
// @Tags Inventory
// @Accept json
// @Produce json
// @Success 200 {object} service.InventoryStats
// @Router /api/inventory/stats [get]
func (h *InventoryHandler) GetStats(c echo.Context) error {
	companyID := getCompanyIDFromContext(c)
	
	stats, err := h.service.GetInventoryStats(companyID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, stats)
}

// GetLowStock godoc
// @Summary Get low stock items
// @Description Get inventory items with low stock
// @Tags Inventory
// @Accept json
// @Produce json
// @Success 200 {array} models.Inventory
// @Router /api/inventory/low-stock [get]
func (h *InventoryHandler) GetLowStock(c echo.Context) error {
	companyID := getCompanyIDFromContext(c)
	
	items, err := h.service.GetLowStockItems(companyID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, items)
}

// GetValuation godoc
// @Summary Get stock valuation
// @Description Get inventory valuation report
// @Tags Inventory
// @Accept json
// @Produce json
// @Success 200 {object} service.StockValuation
// @Router /api/inventory/valuation [get]
func (h *InventoryHandler) GetValuation(c echo.Context) error {
	companyID := getCompanyIDFromContext(c)
	
	valuation, err := h.service.GetStockValuation(companyID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, valuation)
}

// Warehouse handlers

// ListWarehouses godoc
// @Summary List warehouses
// @Description Get a list of warehouses
// @Tags Warehouses
// @Accept json
// @Produce json
// @Success 200 {array} models.Warehouse
// @Router /api/warehouses [get]
func (h *InventoryHandler) ListWarehouses(c echo.Context) error {
	companyID := getCompanyIDFromContext(c)
	
	warehouses, err := h.service.ListWarehouses(companyID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, warehouses)
}

// CreateWarehouse godoc
// @Summary Create warehouse
// @Description Create a new warehouse
// @Tags Warehouses
// @Accept json
// @Produce json
// @Param request body service.CreateWarehouseRequest true "Create warehouse request"
// @Success 201 {object} models.Warehouse
// @Router /api/warehouses [post]
func (h *InventoryHandler) CreateWarehouse(c echo.Context) error {
	var req service.CreateWarehouseRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	companyID := getCompanyIDFromContext(c)
	
	warehouse, err := h.service.CreateWarehouse(companyID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusCreated, warehouse)
}

// Alert handlers

// GetAlerts godoc
// @Summary Get active alerts
// @Description Get active inventory alerts
// @Tags Inventory
// @Accept json
// @Produce json
// @Success 200 {array} models.StockAlert
// @Router /api/inventory/alerts [get]
func (h *InventoryHandler) GetAlerts(c echo.Context) error {
	companyID := getCompanyIDFromContext(c)
	
	alerts, err := h.service.GetActiveAlerts(companyID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, alerts)
}

// Stock take handlers

// ListStockTakes godoc
// @Summary List stock takes
// @Description Get a list of stock takes
// @Tags Stock Takes
// @Accept json
// @Produce json
// @Param status query string false "Filter by status"
// @Param warehouse_id query string false "Filter by warehouse"
// @Success 200 {array} models.StockTake
// @Router /api/stock-takes [get]
func (h *InventoryHandler) ListStockTakes(c echo.Context) error {
	params := make(map[string]interface{})
	
	if status := c.QueryParam("status"); status != "" {
		params["status"] = status
	}
	if warehouseID := c.QueryParam("warehouse_id"); warehouseID != "" {
		params["warehouse_id"] = warehouseID
	}
	
	companyID := getCompanyIDFromContext(c)
	
	stockTakes, err := h.service.ListStockTakes(companyID, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, stockTakes)
}

// CreateStockTake godoc
// @Summary Create stock take
// @Description Create a new stock take
// @Tags Stock Takes
// @Accept json
// @Produce json
// @Param request body service.CreateStockTakeRequest true "Create stock take request"
// @Success 201 {object} models.StockTake
// @Router /api/stock-takes [post]
func (h *InventoryHandler) CreateStockTake(c echo.Context) error {
	var req service.CreateStockTakeRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	companyID := getCompanyIDFromContext(c)
	userID := getUserIDFromContext(c)
	
	stockTake, err := h.service.CreateStockTake(companyID, userID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusCreated, stockTake)
}

// RegisterRoutes registers all inventory routes
func (h *InventoryHandler) RegisterRoutes(e *echo.Echo, authMiddleware echo.MiddlewareFunc) {
	inventory := e.Group("/api/inventory", authMiddleware)
	
	inventory.GET("", h.List)
	inventory.POST("", h.Create)
	inventory.GET("/stats", h.GetStats)
	inventory.GET("/low-stock", h.GetLowStock)
	inventory.GET("/valuation", h.GetValuation)
	inventory.GET("/alerts", h.GetAlerts)
	inventory.POST("/transfer", h.TransferStock)
	inventory.GET("/sku/:sku", h.GetBySKU)
	inventory.GET("/:id", h.Get)
	inventory.PUT("/:id", h.Update)
	inventory.DELETE("/:id", h.Delete)
	inventory.POST("/:id/adjust", h.AdjustStock)
	inventory.GET("/:id/movements", h.GetMovements)
	
	// Warehouse routes
	warehouses := e.Group("/api/warehouses", authMiddleware)
	warehouses.GET("", h.ListWarehouses)
	warehouses.POST("", h.CreateWarehouse)
	
	// Stock take routes
	stockTakes := e.Group("/api/stock-takes", authMiddleware)
	stockTakes.GET("", h.ListStockTakes)
	stockTakes.POST("", h.CreateStockTake)
}