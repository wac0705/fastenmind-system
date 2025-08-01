package handler

import (
	"net/http"
	"strconv"
	"time"
	
	"github.com/fastenmind/fastener-api/internal/middleware"
	"github.com/fastenmind/fastener-api/internal/models"
	"github.com/fastenmind/fastener-api/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type FinanceHandler struct {
	financeService service.FinanceService
}

func NewFinanceHandler(financeService service.FinanceService) *FinanceHandler {
	return &FinanceHandler{
		financeService: financeService,
	}
}

func (h *FinanceHandler) Register(e *echo.Echo) {
	api := e.Group("/api/v1")
	
	// Invoice routes
	api.POST("/invoices", h.CreateInvoice, middleware.Auth())
	api.PUT("/invoices/:id", h.UpdateInvoice, middleware.Auth())
	api.GET("/invoices/:id", h.GetInvoice, middleware.Auth())
	api.GET("/invoices", h.ListInvoices, middleware.Auth())
	api.POST("/invoices/generate/:orderId", h.GenerateInvoiceFromOrder, middleware.Auth())
	api.GET("/invoices/:id/items", h.GetInvoiceItems, middleware.Auth())
	
	// Payment routes
	api.POST("/payments", h.ProcessPayment, middleware.Auth())
	api.GET("/payments/:id", h.GetPayment, middleware.Auth())
	api.GET("/payments", h.ListPayments, middleware.Auth())
	api.GET("/invoices/:invoiceId/payments", h.GetPaymentsByInvoice, middleware.Auth())
	
	// Expense routes
	api.POST("/expenses", h.CreateExpense, middleware.Auth())
	api.PUT("/expenses/:id", h.UpdateExpense, middleware.Auth())
	api.GET("/expenses/:id", h.GetExpense, middleware.Auth())
	api.GET("/expenses", h.ListExpenses, middleware.Auth())
	api.POST("/expenses/:id/approve", h.ApproveExpense, middleware.Auth())
	api.POST("/expenses/:id/reject", h.RejectExpense, middleware.Auth())
	
	// AR/AP routes
	api.GET("/accounts-receivable", h.GetAccountReceivables, middleware.Auth())
	api.GET("/accounts-payable", h.GetAccountPayables, middleware.Auth())
	api.GET("/ar-summary", h.GetARSummary, middleware.Auth())
	api.GET("/ap-summary", h.GetAPSummary, middleware.Auth())
	
	// Bank account routes
	api.POST("/bank-accounts", h.CreateBankAccount, middleware.Auth())
	api.PUT("/bank-accounts/:id", h.UpdateBankAccount, middleware.Auth())
	api.GET("/bank-accounts/:id", h.GetBankAccount, middleware.Auth())
	api.GET("/bank-accounts", h.ListBankAccounts, middleware.Auth())
	
	// Financial period routes
	api.POST("/financial-periods", h.CreateFinancialPeriod, middleware.Auth())
	api.POST("/financial-periods/:id/close", h.CloseFinancialPeriod, middleware.Auth())
	api.GET("/financial-periods/current", h.GetCurrentPeriod, middleware.Auth())
	api.GET("/financial-periods", h.ListFinancialPeriods, middleware.Auth())
	
	// Report routes
	api.GET("/finance/dashboard", h.GetFinancialDashboard, middleware.Auth())
	api.GET("/finance/cash-flow", h.GetCashFlowReport, middleware.Auth())
	api.GET("/finance/aging-report", h.GetAgingReport, middleware.Auth())
}

// Invoice handlers
func (h *FinanceHandler) CreateInvoice(c echo.Context) error {
	var req models.Invoice
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	
	user := c.Get("user").(*models.User)
	req.CompanyID = user.CompanyID
	req.CreatedBy = user.ID
	
	if err := h.financeService.CreateInvoice(&req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusCreated, req)
}

func (h *FinanceHandler) UpdateInvoice(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid invoice ID")
	}
	
	var req models.Invoice
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	
	req.ID = id
	
	if err := h.financeService.UpdateInvoice(&req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, req)
}

func (h *FinanceHandler) GetInvoice(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid invoice ID")
	}
	
	invoice, err := h.financeService.GetInvoice(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Invoice not found")
	}
	
	return c.JSON(http.StatusOK, invoice)
}

func (h *FinanceHandler) ListInvoices(c echo.Context) error {
	user := c.Get("user").(*models.User)
	params := make(map[string]interface{})
	
	// Parse query parameters
	if invoiceType := c.QueryParam("type"); invoiceType != "" {
		params["type"] = invoiceType
	}
	if status := c.QueryParam("status"); status != "" {
		params["status"] = status
	}
	if customerID := c.QueryParam("customer_id"); customerID != "" {
		params["customer_id"] = customerID
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
	if overdue := c.QueryParam("overdue"); overdue == "true" {
		params["overdue"] = true
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
	
	invoices, total, err := h.financeService.ListInvoices(user.CompanyID, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": invoices,
		"pagination": map[string]interface{}{
			"page":       page,
			"page_size":  pageSize,
			"total":      total,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

func (h *FinanceHandler) GenerateInvoiceFromOrder(c echo.Context) error {
	orderID, err := uuid.Parse(c.Param("orderId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid order ID")
	}
	
	user := c.Get("user").(*models.User)
	
	invoice, err := h.financeService.GenerateInvoiceFromOrder(orderID, user.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusCreated, invoice)
}

func (h *FinanceHandler) GetInvoiceItems(c echo.Context) error {
	invoiceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid invoice ID")
	}
	
	items, err := h.financeService.GetInvoiceItems(invoiceID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, items)
}

// Payment handlers
func (h *FinanceHandler) ProcessPayment(c echo.Context) error {
	var req models.Payment
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	
	user := c.Get("user").(*models.User)
	req.CompanyID = user.CompanyID
	req.CreatedBy = user.ID
	
	if err := h.financeService.ProcessPayment(&req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusCreated, req)
}

func (h *FinanceHandler) GetPayment(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid payment ID")
	}
	
	payment, err := h.financeService.GetPayment(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Payment not found")
	}
	
	return c.JSON(http.StatusOK, payment)
}

func (h *FinanceHandler) ListPayments(c echo.Context) error {
	user := c.Get("user").(*models.User)
	params := make(map[string]interface{})
	
	// Parse query parameters
	if paymentType := c.QueryParam("type"); paymentType != "" {
		params["type"] = paymentType
	}
	if status := c.QueryParam("status"); status != "" {
		params["status"] = status
	}
	if paymentMethod := c.QueryParam("payment_method"); paymentMethod != "" {
		params["payment_method"] = paymentMethod
	}
	if startDate := c.QueryParam("start_date"); startDate != "" {
		params["start_date"] = startDate
	}
	if endDate := c.QueryParam("end_date"); endDate != "" {
		params["end_date"] = endDate
	}
	
	// Pagination
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 {
		page = 1
	}
	params["page"] = page
	
	payments, total, err := h.financeService.ListPayments(user.CompanyID, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": payments,
		"pagination": map[string]interface{}{
			"page":       page,
			"page_size":  20,
			"total":      total,
			"total_pages": (total + 19) / 20,
		},
	})
}

func (h *FinanceHandler) GetPaymentsByInvoice(c echo.Context) error {
	invoiceID, err := uuid.Parse(c.Param("invoiceId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid invoice ID")
	}
	
	payments, err := h.financeService.GetPaymentsByInvoice(invoiceID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, payments)
}

// Expense handlers
func (h *FinanceHandler) CreateExpense(c echo.Context) error {
	var req models.Expense
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	
	user := c.Get("user").(*models.User)
	req.CompanyID = user.CompanyID
	req.SubmittedBy = user.ID
	
	if err := h.financeService.CreateExpense(&req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusCreated, req)
}

func (h *FinanceHandler) UpdateExpense(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid expense ID")
	}
	
	var req models.Expense
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	
	req.ID = id
	
	if err := h.financeService.UpdateExpense(&req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, req)
}

func (h *FinanceHandler) GetExpense(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid expense ID")
	}
	
	expense, err := h.financeService.GetExpense(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Expense not found")
	}
	
	return c.JSON(http.StatusOK, expense)
}

func (h *FinanceHandler) ListExpenses(c echo.Context) error {
	user := c.Get("user").(*models.User)
	params := make(map[string]interface{})
	
	// Parse query parameters
	if category := c.QueryParam("category"); category != "" {
		params["category"] = category
	}
	if status := c.QueryParam("status"); status != "" {
		params["status"] = status
	}
	if paymentStatus := c.QueryParam("payment_status"); paymentStatus != "" {
		params["payment_status"] = paymentStatus
	}
	if startDate := c.QueryParam("start_date"); startDate != "" {
		params["start_date"] = startDate
	}
	if endDate := c.QueryParam("end_date"); endDate != "" {
		params["end_date"] = endDate
	}
	
	// Pagination
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 {
		page = 1
	}
	params["page"] = page
	
	expenses, total, err := h.financeService.ListExpenses(user.CompanyID, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": expenses,
		"pagination": map[string]interface{}{
			"page":       page,
			"page_size":  20,
			"total":      total,
			"total_pages": (total + 19) / 20,
		},
	})
}

func (h *FinanceHandler) ApproveExpense(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid expense ID")
	}
	
	user := c.Get("user").(*models.User)
	
	if err := h.financeService.ApproveExpense(id, user.ID); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Expense approved successfully",
	})
}

func (h *FinanceHandler) RejectExpense(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid expense ID")
	}
	
	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	
	user := c.Get("user").(*models.User)
	
	if err := h.financeService.RejectExpense(id, user.ID, req.Reason); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Expense rejected successfully",
	})
}

// AR/AP handlers
func (h *FinanceHandler) GetAccountReceivables(c echo.Context) error {
	user := c.Get("user").(*models.User)
	params := make(map[string]interface{})
	
	// Parse query parameters
	if status := c.QueryParam("status"); status != "" {
		params["status"] = status
	}
	if agingCategory := c.QueryParam("aging_category"); agingCategory != "" {
		params["aging_category"] = agingCategory
	}
	if customerID := c.QueryParam("customer_id"); customerID != "" {
		params["customer_id"] = customerID
	}
	
	ars, err := h.financeService.GetAccountReceivables(user.CompanyID, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, ars)
}

func (h *FinanceHandler) GetAccountPayables(c echo.Context) error {
	user := c.Get("user").(*models.User)
	params := make(map[string]interface{})
	
	// Parse query parameters
	if status := c.QueryParam("status"); status != "" {
		params["status"] = status
	}
	if priority := c.QueryParam("payment_priority"); priority != "" {
		params["payment_priority"] = priority
	}
	if supplierID := c.QueryParam("supplier_id"); supplierID != "" {
		params["supplier_id"] = supplierID
	}
	
	aps, err := h.financeService.GetAccountPayables(user.CompanyID, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, aps)
}

func (h *FinanceHandler) GetARSummary(c echo.Context) error {
	user := c.Get("user").(*models.User)
	
	summary, err := h.financeService.GetARSummary(user.CompanyID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, summary)
}

func (h *FinanceHandler) GetAPSummary(c echo.Context) error {
	user := c.Get("user").(*models.User)
	
	summary, err := h.financeService.GetAPSummary(user.CompanyID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, summary)
}

// Bank account handlers
func (h *FinanceHandler) CreateBankAccount(c echo.Context) error {
	var req models.BankAccount
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	
	user := c.Get("user").(*models.User)
	req.CompanyID = user.CompanyID
	
	if err := h.financeService.CreateBankAccount(&req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusCreated, req)
}

func (h *FinanceHandler) UpdateBankAccount(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid bank account ID")
	}
	
	var req models.BankAccount
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	
	req.ID = id
	
	if err := h.financeService.UpdateBankAccount(&req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, req)
}

func (h *FinanceHandler) GetBankAccount(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid bank account ID")
	}
	
	account, err := h.financeService.GetBankAccount(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Bank account not found")
	}
	
	return c.JSON(http.StatusOK, account)
}

func (h *FinanceHandler) ListBankAccounts(c echo.Context) error {
	user := c.Get("user").(*models.User)
	
	accounts, err := h.financeService.ListBankAccounts(user.CompanyID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, accounts)
}

// Financial period handlers
func (h *FinanceHandler) CreateFinancialPeriod(c echo.Context) error {
	var req models.FinancialPeriod
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	
	user := c.Get("user").(*models.User)
	req.CompanyID = user.CompanyID
	
	if err := h.financeService.CreateFinancialPeriod(&req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusCreated, req)
}

func (h *FinanceHandler) CloseFinancialPeriod(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid period ID")
	}
	
	user := c.Get("user").(*models.User)
	
	if err := h.financeService.CloseFinancialPeriod(id, user.ID); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Financial period closed successfully",
	})
}

func (h *FinanceHandler) GetCurrentPeriod(c echo.Context) error {
	user := c.Get("user").(*models.User)
	
	period, err := h.financeService.GetCurrentPeriod(user.CompanyID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "No current financial period found")
	}
	
	return c.JSON(http.StatusOK, period)
}

func (h *FinanceHandler) ListFinancialPeriods(c echo.Context) error {
	user := c.Get("user").(*models.User)
	
	periods, err := h.financeService.ListFinancialPeriods(user.CompanyID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, periods)
}

// Report handlers
func (h *FinanceHandler) GetFinancialDashboard(c echo.Context) error {
	user := c.Get("user").(*models.User)
	
	dashboard, err := h.financeService.GetFinancialDashboard(user.CompanyID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, dashboard)
}

func (h *FinanceHandler) GetCashFlowReport(c echo.Context) error {
	user := c.Get("user").(*models.User)
	
	// Parse date parameters
	startDateStr := c.QueryParam("start_date")
	endDateStr := c.QueryParam("end_date")
	
	if startDateStr == "" || endDateStr == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "start_date and end_date are required")
	}
	
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid start_date format")
	}
	
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid end_date format")
	}
	
	report, err := h.financeService.GetCashFlowReport(user.CompanyID, startDate, endDate)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, report)
}

func (h *FinanceHandler) GetAgingReport(c echo.Context) error {
	user := c.Get("user").(*models.User)
	
	reportType := c.QueryParam("type")
	if reportType != "receivable" && reportType != "payable" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid report type. Must be 'receivable' or 'payable'")
	}
	
	report, err := h.financeService.GetAgingReport(user.CompanyID, reportType)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, report)
}