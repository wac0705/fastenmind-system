package handler

import (
	"net/http"
	"strconv"

	"github.com/fastenmind/fastener-api/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type QuoteHandler struct {
	service service.QuoteService
}

func NewQuoteHandler(service service.QuoteService) *QuoteHandler {
	return &QuoteHandler{service: service}
}

// List godoc
// @Summary List quotes
// @Description Get a list of quotes with optional filters
// @Tags Quotes
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param status query string false "Filter by status"
// @Param customer_id query string false "Filter by customer"
// @Param engineer_id query string false "Filter by engineer"
// @Param search query string false "Search term"
// @Success 200 {object} map[string]interface{}
// @Router /api/quotes [get]
func (h *QuoteHandler) List(c echo.Context) error {
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
	if engineerID := c.QueryParam("engineer_id"); engineerID != "" {
		params["engineer_id"] = engineerID
	}
	if search := c.QueryParam("search"); search != "" {
		params["search"] = search
	}
	
	companyID := getCompanyIDFromContext(c)
	userID := getUserIDFromContext(c)
	
	quotes, total, err := h.service.List(companyID, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": quotes,
		"pagination": map[string]interface{}{
			"total":     total,
			"page":      params["page"],
			"page_size": params["page_size"],
		},
	})
}

// Get godoc
// @Summary Get quote details
// @Description Get detailed information about a specific quote
// @Tags Quotes
// @Accept json
// @Produce json
// @Param id path string true "Quote ID"
// @Success 200 {object} models.Quote
// @Router /api/quotes/{id} [get]
func (h *QuoteHandler) Get(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid quote ID")
	}
	
	quote, err := h.service.Get(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Quote not found")
	}
	
	return c.JSON(http.StatusOK, quote)
}

// Create godoc
// @Summary Create quote
// @Description Create a new quote for an inquiry
// @Tags Quotes
// @Accept json
// @Produce json
// @Param request body service.CreateQuoteRequest true "Create quote request"
// @Success 201 {object} models.Quote
// @Router /api/quotes [post]
func (h *QuoteHandler) Create(c echo.Context) error {
	var req service.CreateQuoteRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	companyID := getCompanyIDFromContext(c)
	userID := getUserIDFromContext(c)
	
	quote, err := h.service.Create(companyID, userID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusCreated, quote)
}

// Update godoc
// @Summary Update quote
// @Description Update quote details
// @Tags Quotes
// @Accept json
// @Produce json
// @Param id path string true "Quote ID"
// @Param request body service.UpdateQuoteRequest true "Update quote request"
// @Success 200 {object} models.Quote
// @Router /api/quotes/{id} [put]
func (h *QuoteHandler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid quote ID")
	}
	
	var req service.UpdateQuoteRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	companyID := getCompanyIDFromContext(c)
	userID := getUserIDFromContext(c)
	
	quote, err := h.service.Update(id, userID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusOK, quote)
}

// Delete godoc
// @Summary Delete quote
// @Description Delete a quote (soft delete)
// @Tags Quotes
// @Accept json
// @Produce json
// @Param id path string true "Quote ID"
// @Success 204
// @Router /api/quotes/{id} [delete]
func (h *QuoteHandler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid quote ID")
	}
	
	if err := h.service.Delete(id); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	return c.NoContent(http.StatusNoContent)
}

// SubmitForReview godoc
// @Summary Submit quote for review
// @Description Submit a draft quote for review
// @Tags Quotes
// @Accept json
// @Produce json
// @Param id path string true "Quote ID"
// @Success 200 {object} models.Quote
// @Router /api/quotes/{id}/submit-review [post]
func (h *QuoteHandler) SubmitForReview(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid quote ID")
	}
	
	companyID := getCompanyIDFromContext(c)
	userID := getUserIDFromContext(c)
	
	quote, err := h.service.SubmitForReview(id, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusOK, quote)
}

// Review godoc
// @Summary Review quote
// @Description Review a quote (approve or reject)
// @Tags Quotes
// @Accept json
// @Produce json
// @Param id path string true "Quote ID"
// @Param request body service.ReviewQuoteRequest true "Review request"
// @Success 200 {object} models.Quote
// @Router /api/quotes/{id}/review [post]
func (h *QuoteHandler) Review(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid quote ID")
	}
	
	var req service.ReviewQuoteRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	companyID := getCompanyIDFromContext(c)
	userID := getUserIDFromContext(c)
	
	quote, err := h.service.Review(id, userID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusOK, quote)
}

// Send godoc
// @Summary Send quote to customer
// @Description Send approved quote to customer via email
// @Tags Quotes
// @Accept json
// @Produce json
// @Param id path string true "Quote ID"
// @Param request body service.SendQuoteRequest true "Send request"
// @Success 200 {object} models.Quote
// @Router /api/quotes/{id}/send [post]
func (h *QuoteHandler) Send(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid quote ID")
	}
	
	var req service.SendQuoteRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	companyID := getCompanyIDFromContext(c)
	userID := getUserIDFromContext(c)
	
	quote, err := h.service.Send(id, userID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusOK, quote)
}

// GetCostBreakdown godoc
// @Summary Get cost breakdown
// @Description Get detailed cost breakdown for a quote
// @Tags Quotes
// @Accept json
// @Produce json
// @Param id path string true "Quote ID"
// @Success 200 {object} service.CostBreakdown
// @Router /api/quotes/{id}/cost-breakdown [get]
func (h *QuoteHandler) GetCostBreakdown(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid quote ID")
	}
	
	breakdown, err := h.service.GetCostBreakdown(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Quote not found")
	}
	
	return c.JSON(http.StatusOK, breakdown)
}

// GeneratePDF godoc
// @Summary Generate quote PDF
// @Description Generate PDF document for a quote
// @Tags Quotes
// @Accept json
// @Produce application/pdf
// @Param id path string true "Quote ID"
// @Success 200 {file} binary
// @Router /api/quotes/{id}/pdf [get]
func (h *QuoteHandler) GeneratePDF(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid quote ID")
	}
	
	pdfData, filename, err := h.service.GeneratePDF(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	c.Response().Header().Set("Content-Type", "application/pdf")
	c.Response().Header().Set("Content-Disposition", "attachment; filename="+filename)
	
	return c.Blob(http.StatusOK, "application/pdf", pdfData)
}

// GetVersions godoc
// @Summary Get quote versions
// @Description Get version history for a quote
// @Tags Quotes
// @Accept json
// @Produce json
// @Param id path string true "Quote ID"
// @Success 200 {array} models.QuoteVersion
// @Router /api/quotes/{id}/versions [get]
func (h *QuoteHandler) GetVersions(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid quote ID")
	}
	
	versions, err := h.service.GetVersions(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, versions)
}

// Duplicate godoc
// @Summary Duplicate quote
// @Description Create a copy of an existing quote
// @Tags Quotes
// @Accept json
// @Produce json
// @Param id path string true "Quote ID"
// @Success 201 {object} models.Quote
// @Router /api/quotes/{id}/duplicate [post]
func (h *QuoteHandler) Duplicate(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid quote ID")
	}
	
	companyID := getCompanyIDFromContext(c)
	userID := getUserIDFromContext(c)
	
	quote, err := h.service.Duplicate(id, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusCreated, quote)
}

// RegisterRoutes registers all quote routes
func (h *QuoteHandler) RegisterRoutes(e *echo.Echo, authMiddleware echo.MiddlewareFunc) {
	quotes := e.Group("/api/quotes", authMiddleware)
	
	quotes.GET("", h.List)
	quotes.POST("", h.Create)
	quotes.GET("/:id", h.Get)
	quotes.PUT("/:id", h.Update)
	quotes.DELETE("/:id", h.Delete)
	quotes.POST("/:id/submit-review", h.SubmitForReview)
	quotes.POST("/:id/review", h.Review)
	quotes.POST("/:id/send", h.Send)
	quotes.GET("/:id/cost-breakdown", h.GetCostBreakdown)
	quotes.GET("/:id/pdf", h.GeneratePDF)
	quotes.GET("/:id/versions", h.GetVersions)
	quotes.POST("/:id/duplicate", h.Duplicate)
}