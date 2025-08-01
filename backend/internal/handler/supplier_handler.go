package handler

import (
	"net/http"
	"strconv"

	"github.com/fastenmind/fastener-api/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type SupplierHandler struct {
	supplierService service.SupplierService
}

func NewSupplierHandler(supplierService service.SupplierService) *SupplierHandler {
	return &SupplierHandler{
		supplierService: supplierService,
	}
}

// Supplier operations
func (h *SupplierHandler) CreateSupplier(c echo.Context) error {
	var req service.CreateSupplierRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	supplier, err := h.supplierService.CreateSupplier(&req, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, supplier)
}

func (h *SupplierHandler) UpdateSupplier(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid supplier ID"})
	}

	var req service.UpdateSupplierRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	supplier, err := h.supplierService.UpdateSupplier(id, &req, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, supplier)
}

func (h *SupplierHandler) GetSupplier(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid supplier ID"})
	}

	supplier, err := h.supplierService.GetSupplier(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Supplier not found"})
	}

	return c.JSON(http.StatusOK, supplier)
}

func (h *SupplierHandler) ListSuppliers(c echo.Context) error {
	companyID, err := getCompanyIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

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

	if search := c.QueryParam("search"); search != "" {
		params["search"] = search
	}

	if supplierType := c.QueryParam("type"); supplierType != "" {
		params["type"] = supplierType
	}

	if status := c.QueryParam("status"); status != "" {
		params["status"] = status
	}

	if country := c.QueryParam("country"); country != "" {
		params["country"] = country
	}

	if riskLevel := c.QueryParam("risk_level"); riskLevel != "" {
		params["risk_level"] = riskLevel
	}

	if sortBy := c.QueryParam("sort_by"); sortBy != "" {
		params["sort_by"] = sortBy
	}

	if sortOrder := c.QueryParam("sort_order"); sortOrder != "" {
		params["sort_order"] = sortOrder
	}

	suppliers, total, err := h.supplierService.ListSuppliers(companyID, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":  suppliers,
		"total": total,
	})
}

// Supplier Contact operations
func (h *SupplierHandler) AddSupplierContact(c echo.Context) error {
	supplierID, err := uuid.Parse(c.Param("supplier_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid supplier ID"})
	}

	var req service.CreateSupplierContactRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	contact, err := h.supplierService.AddSupplierContact(supplierID, &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, contact)
}

func (h *SupplierHandler) UpdateSupplierContact(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid contact ID"})
	}

	var req service.UpdateSupplierContactRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	contact, err := h.supplierService.UpdateSupplierContact(id, &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, contact)
}

func (h *SupplierHandler) GetSupplierContacts(c echo.Context) error {
	supplierID, err := uuid.Parse(c.Param("supplier_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid supplier ID"})
	}

	contacts, err := h.supplierService.GetSupplierContacts(supplierID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, contacts)
}

func (h *SupplierHandler) DeleteSupplierContact(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid contact ID"})
	}

	if err := h.supplierService.DeleteSupplierContact(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// Supplier Product operations
func (h *SupplierHandler) AddSupplierProduct(c echo.Context) error {
	supplierID, err := uuid.Parse(c.Param("supplier_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid supplier ID"})
	}

	var req service.CreateSupplierProductRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	product, err := h.supplierService.AddSupplierProduct(supplierID, &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, product)
}

func (h *SupplierHandler) UpdateSupplierProduct(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
	}

	var req service.UpdateSupplierProductRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	product, err := h.supplierService.UpdateSupplierProduct(id, &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, product)
}

func (h *SupplierHandler) GetSupplierProducts(c echo.Context) error {
	supplierID, err := uuid.Parse(c.Param("supplier_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid supplier ID"})
	}

	params := make(map[string]interface{})

	if status := c.QueryParam("status"); status != "" {
		params["status"] = status
	}

	if category := c.QueryParam("category"); category != "" {
		params["category"] = category
	}

	if preferred := c.QueryParam("is_preferred"); preferred != "" {
		if p, err := strconv.ParseBool(preferred); err == nil {
			params["is_preferred"] = p
		}
	}

	if search := c.QueryParam("search"); search != "" {
		params["search"] = search
	}

	products, err := h.supplierService.GetSupplierProducts(supplierID, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, products)
}

func (h *SupplierHandler) DeleteSupplierProduct(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
	}

	if err := h.supplierService.DeleteSupplierProduct(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// Purchase Order operations
func (h *SupplierHandler) CreatePurchaseOrder(c echo.Context) error {
	var req service.CreatePurchaseOrderRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	order, err := h.supplierService.CreatePurchaseOrder(&req, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, order)
}

func (h *SupplierHandler) UpdatePurchaseOrder(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid purchase order ID"})
	}

	var req service.UpdatePurchaseOrderRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	order, err := h.supplierService.UpdatePurchaseOrder(id, &req, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, order)
}

func (h *SupplierHandler) GetPurchaseOrder(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid purchase order ID"})
	}

	order, err := h.supplierService.GetPurchaseOrder(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Purchase order not found"})
	}

	return c.JSON(http.StatusOK, order)
}

func (h *SupplierHandler) ListPurchaseOrders(c echo.Context) error {
	companyID, err := getCompanyIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

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

	if supplierID := c.QueryParam("supplier_id"); supplierID != "" {
		params["supplier_id"] = supplierID
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

	orders, total, err := h.supplierService.ListPurchaseOrders(companyID, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":  orders,
		"total": total,
	})
}

func (h *SupplierHandler) ApprovePurchaseOrder(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid purchase order ID"})
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	if err := h.supplierService.ApprovePurchaseOrder(id, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Purchase order approved successfully"})
}

func (h *SupplierHandler) SendPurchaseOrder(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid purchase order ID"})
	}

	if err := h.supplierService.SendPurchaseOrder(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Purchase order sent successfully"})
}

func (h *SupplierHandler) ReceivePurchaseOrder(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid purchase order ID"})
	}

	var req struct {
		Items []service.PurchaseOrderReceiptItem `json:"items"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := h.supplierService.ReceivePurchaseOrder(id, req.Items); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Purchase order received successfully"})
}

// Purchase Order Item operations
func (h *SupplierHandler) AddPurchaseOrderItem(c echo.Context) error {
	purchaseOrderID, err := uuid.Parse(c.Param("purchase_order_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid purchase order ID"})
	}

	var req service.CreatePurchaseOrderItemRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	item, err := h.supplierService.AddPurchaseOrderItem(purchaseOrderID, &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, item)
}

func (h *SupplierHandler) UpdatePurchaseOrderItem(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid item ID"})
	}

	var req service.UpdatePurchaseOrderItemRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	item, err := h.supplierService.UpdatePurchaseOrderItem(id, &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, item)
}

func (h *SupplierHandler) GetPurchaseOrderItems(c echo.Context) error {
	purchaseOrderID, err := uuid.Parse(c.Param("purchase_order_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid purchase order ID"})
	}

	items, err := h.supplierService.GetPurchaseOrderItems(purchaseOrderID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, items)
}

func (h *SupplierHandler) DeletePurchaseOrderItem(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid item ID"})
	}

	if err := h.supplierService.DeletePurchaseOrderItem(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// Supplier Evaluation operations
func (h *SupplierHandler) CreateSupplierEvaluation(c echo.Context) error {
	var req service.CreateSupplierEvaluationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	evaluation, err := h.supplierService.CreateSupplierEvaluation(&req, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, evaluation)
}

func (h *SupplierHandler) UpdateSupplierEvaluation(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid evaluation ID"})
	}

	var req service.UpdateSupplierEvaluationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	evaluation, err := h.supplierService.UpdateSupplierEvaluation(id, &req, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, evaluation)
}

func (h *SupplierHandler) GetSupplierEvaluation(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid evaluation ID"})
	}

	evaluation, err := h.supplierService.GetSupplierEvaluation(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Evaluation not found"})
	}

	return c.JSON(http.StatusOK, evaluation)
}

func (h *SupplierHandler) ListSupplierEvaluations(c echo.Context) error {
	companyID, err := getCompanyIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

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

	if supplierID := c.QueryParam("supplier_id"); supplierID != "" {
		params["supplier_id"] = supplierID
	}

	if evaluationType := c.QueryParam("evaluation_type"); evaluationType != "" {
		params["evaluation_type"] = evaluationType
	}

	if evaluatedBy := c.QueryParam("evaluated_by"); evaluatedBy != "" {
		params["evaluated_by"] = evaluatedBy
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

	evaluations, total, err := h.supplierService.ListSupplierEvaluations(companyID, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":  evaluations,
		"total": total,
	})
}

func (h *SupplierHandler) ApproveSupplierEvaluation(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid evaluation ID"})
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	if err := h.supplierService.ApproveSupplierEvaluation(id, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Supplier evaluation approved successfully"})
}

// Business operations
func (h *SupplierHandler) UpdateSupplierPerformance(c echo.Context) error {
	supplierID, err := uuid.Parse(c.Param("supplier_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid supplier ID"})
	}

	if err := h.supplierService.UpdateSupplierPerformance(supplierID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Supplier performance updated successfully"})
}

func (h *SupplierHandler) CalculateSupplierRisk(c echo.Context) error {
	supplierID, err := uuid.Parse(c.Param("supplier_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid supplier ID"})
	}

	riskLevel, err := h.supplierService.CalculateSupplierRisk(supplierID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"supplier_id": supplierID,
		"risk_level":  riskLevel,
	})
}

func (h *SupplierHandler) GetSupplierDashboard(c echo.Context) error {
	companyID, err := getCompanyIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	dashboard, err := h.supplierService.GetSupplierDashboard(companyID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, dashboard)
}

// Note: getUserIDFromContext and getCompanyIDFromContext are defined in common.go