package handler

import (
	"net/http"
	"strconv"
	
	"github.com/fastenmind/fastener-api/internal/middleware"
	"github.com/fastenmind/fastener-api/internal/models"
	"github.com/fastenmind/fastener-api/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ProductionHandler struct {
	productionService service.ProductionService
}

func NewProductionHandler(productionService service.ProductionService) *ProductionHandler {
	return &ProductionHandler{
		productionService: productionService,
	}
}

func (h *ProductionHandler) Register(e *echo.Echo) {
	api := e.Group("/api/v1")
	
	// Production Order routes
	api.POST("/production-orders", h.CreateProductionOrder, middleware.Auth())
	api.PUT("/production-orders/:id", h.UpdateProductionOrder, middleware.Auth())
	api.GET("/production-orders/:id", h.GetProductionOrder, middleware.Auth())
	api.GET("/production-orders", h.ListProductionOrders, middleware.Auth())
	api.POST("/production-orders/:id/release", h.ReleaseProductionOrder, middleware.Auth())
	api.POST("/production-orders/:id/start", h.StartProductionOrder, middleware.Auth())
	api.POST("/production-orders/:id/complete", h.CompleteProductionOrder, middleware.Auth())
	api.POST("/production-orders/:id/cancel", h.CancelProductionOrder, middleware.Auth())
	
	// Production Route routes
	api.POST("/production-routes", h.CreateProductionRoute, middleware.Auth())
	api.PUT("/production-routes/:id", h.UpdateProductionRoute, middleware.Auth())
	api.GET("/production-routes/:id", h.GetProductionRoute, middleware.Auth())
	api.GET("/production-routes", h.ListProductionRoutes, middleware.Auth())
	
	// Route Operation routes
	api.POST("/production-routes/:routeId/operations", h.CreateRouteOperation, middleware.Auth())
	api.PUT("/route-operations/:id", h.UpdateRouteOperation, middleware.Auth())
	api.GET("/production-routes/:routeId/operations", h.GetRouteOperations, middleware.Auth())
	api.DELETE("/route-operations/:id", h.DeleteRouteOperation, middleware.Auth())
	
	// Work Station routes
	api.POST("/work-stations", h.CreateWorkStation, middleware.Auth())
	api.PUT("/work-stations/:id", h.UpdateWorkStation, middleware.Auth())
	api.GET("/work-stations/:id", h.GetWorkStation, middleware.Auth())
	api.GET("/work-stations", h.ListWorkStations, middleware.Auth())
	
	// Production Task routes
	api.POST("/production-tasks", h.CreateProductionTask, middleware.Auth())
	api.PUT("/production-tasks/:id", h.UpdateProductionTask, middleware.Auth())
	api.GET("/production-tasks/:id", h.GetProductionTask, middleware.Auth())
	api.GET("/production-tasks", h.ListProductionTasks, middleware.Auth())
	api.POST("/production-tasks/:id/assign", h.AssignTask, middleware.Auth())
	api.POST("/production-tasks/:id/start", h.StartTask, middleware.Auth())
	api.POST("/production-tasks/:id/complete", h.CompleteTask, middleware.Auth())
	
	// Production Material routes
	api.POST("/production-materials", h.CreateProductionMaterial, middleware.Auth())
	api.POST("/production-orders/:id/issue-materials", h.IssueMaterials, middleware.Auth())
	api.GET("/production-orders/:id/materials", h.GetProductionMaterials, middleware.Auth())
	
	// Quality Inspection routes
	api.POST("/quality-inspections", h.CreateQualityInspection, middleware.Auth())
	api.PUT("/quality-inspections/:id", h.UpdateQualityInspection, middleware.Auth())
	api.GET("/quality-inspections/:id", h.GetQualityInspection, middleware.Auth())
	api.GET("/quality-inspections", h.ListQualityInspections, middleware.Auth())
	api.POST("/quality-inspections/:id/approve", h.ApproveInspection, middleware.Auth())
	api.POST("/quality-inspections/:id/reject", h.RejectInspection, middleware.Auth())
	
	// Dashboard routes
	api.GET("/production/dashboard", h.GetProductionDashboard, middleware.Auth())
	api.GET("/production/stats", h.GetProductionStats, middleware.Auth())
}

// Production Order handlers
func (h *ProductionHandler) CreateProductionOrder(c echo.Context) error {
	var req models.ProductionOrder
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	
	user := c.Get("user").(*models.User)
	req.CompanyID = user.CompanyID
	req.CreatedBy = user.ID
	
	if err := h.productionService.CreateProductionOrder(&req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusCreated, req)
}

func (h *ProductionHandler) UpdateProductionOrder(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid production order ID")
	}
	
	var req models.ProductionOrder
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	
	req.ID = id
	
	if err := h.productionService.UpdateProductionOrder(&req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, req)
}

func (h *ProductionHandler) GetProductionOrder(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid production order ID")
	}
	
	order, err := h.productionService.GetProductionOrder(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Production order not found")
	}
	
	return c.JSON(http.StatusOK, order)
}

func (h *ProductionHandler) ListProductionOrders(c echo.Context) error {
	user := c.Get("user").(*models.User)
	params := make(map[string]interface{})
	
	// Parse query parameters
	if status := c.QueryParam("status"); status != "" {
		params["status"] = status
	}
	if priority := c.QueryParam("priority"); priority != "" {
		params["priority"] = priority
	}
	if customerID := c.QueryParam("customer_id"); customerID != "" {
		params["customer_id"] = customerID
	}
	if inventoryID := c.QueryParam("inventory_id"); inventoryID != "" {
		params["inventory_id"] = inventoryID
	}
	if startDate := c.QueryParam("start_date"); startDate != "" {
		params["start_date"] = startDate
	}
	if endDate := c.QueryParam("end_date"); endDate != "" {
		params["end_date"] = endDate
	}
	if search := c.QueryParam("search"); search != "" {
		params["search"] = search
	}
	
	// Pagination
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 {
		page = 1
	}
	params["page"] = page
	
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 {
		pageSize = 20
	}
	params["page_size"] = pageSize
	
	orders, total, err := h.productionService.ListProductionOrders(user.CompanyID, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": orders,
		"pagination": map[string]interface{}{
			"page":       page,
			"page_size":  pageSize,
			"total":      total,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

func (h *ProductionHandler) ReleaseProductionOrder(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid production order ID")
	}
	
	user := c.Get("user").(*models.User)
	
	if err := h.productionService.ReleaseProductionOrder(id, user.ID); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Production order released successfully",
	})
}

func (h *ProductionHandler) StartProductionOrder(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid production order ID")
	}
	
	user := c.Get("user").(*models.User)
	
	if err := h.productionService.StartProductionOrder(id, user.ID); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Production order started successfully",
	})
}

func (h *ProductionHandler) CompleteProductionOrder(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid production order ID")
	}
	
	user := c.Get("user").(*models.User)
	
	if err := h.productionService.CompleteProductionOrder(id, user.ID); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Production order completed successfully",
	})
}

func (h *ProductionHandler) CancelProductionOrder(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid production order ID")
	}
	
	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	
	user := c.Get("user").(*models.User)
	
	if err := h.productionService.CancelProductionOrder(id, user.ID, req.Reason); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Production order cancelled successfully",
	})
}

// Production Route handlers
func (h *ProductionHandler) CreateProductionRoute(c echo.Context) error {
	var req models.ProductionRoute
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	
	user := c.Get("user").(*models.User)
	req.CompanyID = user.CompanyID
	req.CreatedBy = user.ID
	
	if err := h.productionService.CreateProductionRoute(&req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusCreated, req)
}

func (h *ProductionHandler) UpdateProductionRoute(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid route ID")
	}
	
	var req models.ProductionRoute
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	
	req.ID = id
	
	if err := h.productionService.UpdateProductionRoute(&req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, req)
}

func (h *ProductionHandler) GetProductionRoute(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid route ID")
	}
	
	route, err := h.productionService.GetProductionRoute(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Production route not found")
	}
	
	return c.JSON(http.StatusOK, route)
}

func (h *ProductionHandler) ListProductionRoutes(c echo.Context) error {
	user := c.Get("user").(*models.User)
	params := make(map[string]interface{})
	
	// Parse query parameters
	if status := c.QueryParam("status"); status != "" {
		params["status"] = status
	}
	if inventoryID := c.QueryParam("inventory_id"); inventoryID != "" {
		params["inventory_id"] = inventoryID
	}
	if category := c.QueryParam("product_category"); category != "" {
		params["product_category"] = category
	}
	if search := c.QueryParam("search"); search != "" {
		params["search"] = search
	}
	
	routes, err := h.productionService.ListProductionRoutes(user.CompanyID, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, routes)
}

// Route Operation handlers
func (h *ProductionHandler) CreateRouteOperation(c echo.Context) error {
	routeID, err := uuid.Parse(c.Param("routeId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid route ID")
	}
	
	var req models.RouteOperation
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	
	req.RouteID = routeID
	
	if err := h.productionService.CreateRouteOperation(&req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusCreated, req)
}

func (h *ProductionHandler) UpdateRouteOperation(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid operation ID")
	}
	
	var req models.RouteOperation
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	
	req.ID = id
	
	if err := h.productionService.UpdateRouteOperation(&req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, req)
}

func (h *ProductionHandler) GetRouteOperations(c echo.Context) error {
	routeID, err := uuid.Parse(c.Param("routeId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid route ID")
	}
	
	operations, err := h.productionService.GetRouteOperations(routeID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, operations)
}

func (h *ProductionHandler) DeleteRouteOperation(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid operation ID")
	}
	
	if err := h.productionService.DeleteRouteOperation(id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Route operation deleted successfully",
	})
}

// Work Station handlers
func (h *ProductionHandler) CreateWorkStation(c echo.Context) error {
	var req models.WorkStation
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	
	user := c.Get("user").(*models.User)
	req.CompanyID = user.CompanyID
	req.CreatedBy = user.ID
	
	if err := h.productionService.CreateWorkStation(&req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusCreated, req)
}

func (h *ProductionHandler) UpdateWorkStation(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid work station ID")
	}
	
	var req models.WorkStation
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	
	req.ID = id
	
	if err := h.productionService.UpdateWorkStation(&req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, req)
}

func (h *ProductionHandler) GetWorkStation(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid work station ID")
	}
	
	station, err := h.productionService.GetWorkStation(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Work station not found")
	}
	
	return c.JSON(http.StatusOK, station)
}

func (h *ProductionHandler) ListWorkStations(c echo.Context) error {
	user := c.Get("user").(*models.User)
	params := make(map[string]interface{})
	
	// Parse query parameters
	if stationType := c.QueryParam("type"); stationType != "" {
		params["type"] = stationType
	}
	if status := c.QueryParam("status"); status != "" {
		params["status"] = status
	}
	if department := c.QueryParam("department"); department != "" {
		params["department"] = department
	}
	if search := c.QueryParam("search"); search != "" {
		params["search"] = search
	}
	
	stations, err := h.productionService.ListWorkStations(user.CompanyID, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, stations)
}

// Production Task handlers
func (h *ProductionHandler) CreateProductionTask(c echo.Context) error {
	var req models.ProductionTask
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	
	if err := h.productionService.CreateProductionTask(&req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusCreated, req)
}

func (h *ProductionHandler) UpdateProductionTask(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid task ID")
	}
	
	var req models.ProductionTask
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	
	req.ID = id
	
	if err := h.productionService.UpdateProductionTask(&req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, req)
}

func (h *ProductionHandler) GetProductionTask(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid task ID")
	}
	
	task, err := h.productionService.GetProductionTask(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Production task not found")
	}
	
	return c.JSON(http.StatusOK, task)
}

func (h *ProductionHandler) ListProductionTasks(c echo.Context) error {
	params := make(map[string]interface{})
	
	// Parse query parameters
	if status := c.QueryParam("status"); status != "" {
		params["status"] = status
	}
	if assignedTo := c.QueryParam("assigned_to"); assignedTo != "" {
		params["assigned_to"] = assignedTo
	}
	if workStationID := c.QueryParam("work_station_id"); workStationID != "" {
		params["work_station_id"] = workStationID
	}
	if productionOrderID := c.QueryParam("production_order_id"); productionOrderID != "" {
		params["production_order_id"] = productionOrderID
	}
	if startDate := c.QueryParam("start_date"); startDate != "" {
		params["start_date"] = startDate
	}
	if endDate := c.QueryParam("end_date"); endDate != "" {
		params["end_date"] = endDate
	}
	
	tasks, err := h.productionService.ListProductionTasks(params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, tasks)
}

func (h *ProductionHandler) AssignTask(c echo.Context) error {
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid task ID")
	}
	
	var req struct {
		UserID string `json:"user_id"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}
	
	currentUser := c.Get("user").(*models.User)
	
	if err := h.productionService.AssignTask(taskID, userID, currentUser.ID); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Task assigned successfully",
	})
}

func (h *ProductionHandler) StartTask(c echo.Context) error {
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid task ID")
	}
	
	user := c.Get("user").(*models.User)
	
	if err := h.productionService.StartTask(taskID, user.ID); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Task started successfully",
	})
}

func (h *ProductionHandler) CompleteTask(c echo.Context) error {
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid task ID")
	}
	
	var req struct {
		CompletedQuantity float64 `json:"completed_quantity"`
		QualifiedQuantity float64 `json:"qualified_quantity"`
		Notes             string  `json:"notes"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	
	user := c.Get("user").(*models.User)
	
	if err := h.productionService.CompleteTask(taskID, user.ID, req.CompletedQuantity, req.QualifiedQuantity, req.Notes); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Task completed successfully",
	})
}

// Production Material handlers
func (h *ProductionHandler) CreateProductionMaterial(c echo.Context) error {
	var req models.ProductionMaterial
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	
	if err := h.productionService.CreateProductionMaterial(&req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusCreated, req)
}

func (h *ProductionHandler) IssueMaterials(c echo.Context) error {
	productionOrderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid production order ID")
	}
	
	user := c.Get("user").(*models.User)
	
	if err := h.productionService.IssueMaterials(productionOrderID, user.ID); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Materials issued successfully",
	})
}

func (h *ProductionHandler) GetProductionMaterials(c echo.Context) error {
	productionOrderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid production order ID")
	}
	
	materials, err := h.productionService.GetProductionMaterials(productionOrderID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, materials)
}

// Quality Inspection handlers
func (h *ProductionHandler) CreateQualityInspection(c echo.Context) error {
	var req models.QualityInspection
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	
	user := c.Get("user").(*models.User)
	req.CompanyID = user.CompanyID
	req.InspectorID = user.ID
	
	if err := h.productionService.CreateQualityInspection(&req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusCreated, req)
}

func (h *ProductionHandler) UpdateQualityInspection(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inspection ID")
	}
	
	var req models.QualityInspection
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	
	req.ID = id
	
	if err := h.productionService.UpdateQualityInspection(&req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, req)
}

func (h *ProductionHandler) GetQualityInspection(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inspection ID")
	}
	
	inspection, err := h.productionService.GetQualityInspection(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Quality inspection not found")
	}
	
	return c.JSON(http.StatusOK, inspection)
}

func (h *ProductionHandler) ListQualityInspections(c echo.Context) error {
	user := c.Get("user").(*models.User)
	params := make(map[string]interface{})
	
	// Parse query parameters
	if inspectionType := c.QueryParam("type"); inspectionType != "" {
		params["type"] = inspectionType
	}
	if status := c.QueryParam("status"); status != "" {
		params["status"] = status
	}
	if productionOrderID := c.QueryParam("production_order_id"); productionOrderID != "" {
		params["production_order_id"] = productionOrderID
	}
	if inspectorID := c.QueryParam("inspector_id"); inspectorID != "" {
		params["inspector_id"] = inspectorID
	}
	if startDate := c.QueryParam("start_date"); startDate != "" {
		params["start_date"] = startDate
	}
	if endDate := c.QueryParam("end_date"); endDate != "" {
		params["end_date"] = endDate
	}
	if search := c.QueryParam("search"); search != "" {
		params["search"] = search
	}
	
	// Pagination
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 {
		page = 1
	}
	params["page"] = page
	
	inspections, total, err := h.productionService.ListQualityInspections(user.CompanyID, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": inspections,
		"pagination": map[string]interface{}{
			"page":       page,
			"page_size":  20,
			"total":      total,
			"total_pages": (total + 19) / 20,
		},
	})
}

func (h *ProductionHandler) ApproveInspection(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inspection ID")
	}
	
	user := c.Get("user").(*models.User)
	
	if err := h.productionService.ApproveInspection(id, user.ID); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Inspection approved successfully",
	})
}

func (h *ProductionHandler) RejectInspection(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inspection ID")
	}
	
	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	
	user := c.Get("user").(*models.User)
	
	if err := h.productionService.RejectInspection(id, user.ID, req.Reason); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Inspection rejected successfully",
	})
}

// Dashboard handlers
func (h *ProductionHandler) GetProductionDashboard(c echo.Context) error {
	user := c.Get("user").(*models.User)
	
	dashboard, err := h.productionService.GetProductionDashboard(user.CompanyID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, dashboard)
}

func (h *ProductionHandler) GetProductionStats(c echo.Context) error {
	user := c.Get("user").(*models.User)
	
	stats, err := h.productionService.GetProductionStats(user.CompanyID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, stats)
}