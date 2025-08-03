package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/fastenmind/fastener-api/internal/models"
	"github.com/fastenmind/fastener-api/internal/service"
)

// TradeHandler handles trade-related HTTP requests
type TradeHandler struct {
	tradeService service.TradeService
}

// NewTradeHandler creates a new trade handler
func NewTradeHandler(tradeService service.TradeService) *TradeHandler {
	return &TradeHandler{
		tradeService: tradeService,
	}
}

// Tariff Code Management

// CreateTariffCode creates a new tariff code
func (h *TradeHandler) CreateTariffCode(c echo.Context) error {
	var code models.TariffCode
	if err := c.Bind(&code); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Get company ID from context (assuming it's set by auth middleware)
	code.CompanyID = c.Get("company_id").(uuid.UUID)

	if err := h.tradeService.CreateTariffCode(c.Request().Context(), &code); err != nil {
		if err == service.ErrDuplicateHSCode {
			return c.JSON(http.StatusConflict, map[string]string{"error": "HS code already exists"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create tariff code"})
	}

	return c.JSON(http.StatusCreated, code)
}

// GetTariffCode retrieves a tariff code by ID
func (h *TradeHandler) GetTariffCode(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid tariff code ID"})
	}

	code, err := h.tradeService.GetTariffCode(c.Request().Context(), id)
	if err != nil {
		if err == service.ErrTariffCodeNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Tariff code not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get tariff code"})
	}

	return c.JSON(http.StatusOK, code)
}

// ListTariffCodes lists tariff codes with pagination and filtering
func (h *TradeHandler) ListTariffCodes(c echo.Context) error {
	params := make(map[string]interface{})
	
	// Get company ID from context
	params["company_id"] = c.Get("company_id").(uuid.UUID)
	
	// Parse query parameters
	if search := c.QueryParam("search"); search != "" {
		params["search"] = search
	}
	if category := c.QueryParam("category"); category != "" {
		params["category"] = category
	}
	if isActive := c.QueryParam("is_active"); isActive != "" {
		if active, err := strconv.ParseBool(isActive); err == nil {
			params["is_active"] = active
		}
	}
	
	// Pagination
	page := 1
	pageSize := 20
	if p := c.QueryParam("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if ps := c.QueryParam("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 {
			pageSize = parsed
		}
	}
	params["page"] = page
	params["page_size"] = pageSize

	codes, total, err := h.tradeService.ListTariffCodes(c.Request().Context(), params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to list tariff codes"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": codes,
		"pagination": map[string]interface{}{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// UpdateTariffCode updates a tariff code
func (h *TradeHandler) UpdateTariffCode(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid tariff code ID"})
	}

	var code models.TariffCode
	if err := c.Bind(&code); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	code.ID = id
	code.CompanyID = c.Get("company_id").(uuid.UUID)

	if err := h.tradeService.UpdateTariffCode(c.Request().Context(), &code); err != nil {
		if err == service.ErrTariffCodeNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Tariff code not found"})
		}
		if err == service.ErrDuplicateHSCode {
			return c.JSON(http.StatusConflict, map[string]string{"error": "HS code already exists"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update tariff code"})
	}

	return c.JSON(http.StatusOK, code)
}

// DeleteTariffCode deletes a tariff code
func (h *TradeHandler) DeleteTariffCode(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid tariff code ID"})
	}

	if err := h.tradeService.DeleteTariffCode(c.Request().Context(), id); err != nil {
		if err == service.ErrTariffCodeNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Tariff code not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete tariff code"})
	}

	return c.JSON(http.StatusNoContent, nil)
}

// Tariff Rate Management

// CreateTariffRate creates a new tariff rate
func (h *TradeHandler) CreateTariffRate(c echo.Context) error {
	var rate models.TariffRate
	if err := c.Bind(&rate); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	rate.CompanyID = c.Get("company_id").(uuid.UUID)

	if err := h.tradeService.CreateTariffRate(c.Request().Context(), &rate); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create tariff rate"})
	}

	return c.JSON(http.StatusCreated, rate)
}

// GetTariffRatesByTariffCode gets tariff rates by tariff code ID
func (h *TradeHandler) GetTariffRatesByTariffCode(c echo.Context) error {
	tariffCodeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid tariff code ID"})
	}

	params := map[string]interface{}{
		"tariff_code_id": tariffCodeID,
	}

	rates, err := h.tradeService.ListTariffRates(c.Request().Context(), params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get tariff rates"})
	}

	return c.JSON(http.StatusOK, rates)
}

// ListTariffRates lists all tariff rates with filtering
func (h *TradeHandler) ListTariffRates(c echo.Context) error {
	params := make(map[string]interface{})
	
	// Parse query parameters
	if countryCode := c.QueryParam("country_code"); countryCode != "" {
		params["country_code"] = countryCode
	}
	if tradeType := c.QueryParam("trade_type"); tradeType != "" {
		params["trade_type"] = tradeType
	}
	if isActive := c.QueryParam("is_active"); isActive != "" {
		if active, err := strconv.ParseBool(isActive); err == nil {
			params["is_active"] = active
		}
	}

	rates, err := h.tradeService.ListTariffRates(c.Request().Context(), params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to list tariff rates"})
	}

	return c.JSON(http.StatusOK, rates)
}

// Shipment Management

// ListShipments lists shipments with pagination and filtering
func (h *TradeHandler) ListShipments(c echo.Context) error {
	params := make(map[string]interface{})
	
	// Get company ID from context
	params["company_id"] = c.Get("company_id").(uuid.UUID)
	
	// Parse query parameters
	if status := c.QueryParam("status"); status != "" {
		params["status"] = status
	}
	if transportMode := c.QueryParam("transport_mode"); transportMode != "" {
		params["transport_mode"] = transportMode
	}
	if search := c.QueryParam("search"); search != "" {
		params["search"] = search
	}
	
	// Date range
	if fromDate := c.QueryParam("from_date"); fromDate != "" {
		if t, err := time.Parse("2006-01-02", fromDate); err == nil {
			params["from_date"] = t
		}
	}
	if toDate := c.QueryParam("to_date"); toDate != "" {
		if t, err := time.Parse("2006-01-02", toDate); err == nil {
			params["to_date"] = t
		}
	}
	
	// Pagination
	page := 1
	pageSize := 20
	if p := c.QueryParam("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if ps := c.QueryParam("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 {
			pageSize = parsed
		}
	}
	params["page"] = page
	params["page_size"] = pageSize

	shipments, total, err := h.tradeService.ListShipments(c.Request().Context(), params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to list shipments"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": shipments,
		"pagination": map[string]interface{}{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// CreateShipment creates a new shipment
func (h *TradeHandler) CreateShipment(c echo.Context) error {
	var shipment models.Shipment
	if err := c.Bind(&shipment); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	shipment.CompanyID = c.Get("company_id").(uuid.UUID)

	if err := h.tradeService.CreateShipment(c.Request().Context(), &shipment); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create shipment"})
	}

	return c.JSON(http.StatusCreated, shipment)
}

// GetShipment retrieves a shipment by ID
func (h *TradeHandler) GetShipment(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid shipment ID"})
	}

	shipment, err := h.tradeService.GetShipment(c.Request().Context(), id)
	if err != nil {
		if err == service.ErrShipmentNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Shipment not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get shipment"})
	}

	return c.JSON(http.StatusOK, shipment)
}

// UpdateShipment updates a shipment
func (h *TradeHandler) UpdateShipment(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid shipment ID"})
	}

	var shipment models.Shipment
	if err := c.Bind(&shipment); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	shipment.ID = id
	shipment.CompanyID = c.Get("company_id").(uuid.UUID)

	if err := h.tradeService.UpdateShipment(c.Request().Context(), &shipment); err != nil {
		if err == service.ErrShipmentNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Shipment not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update shipment"})
	}

	return c.JSON(http.StatusOK, shipment)
}

// GetTradeDocumentsByShipment gets documents for a shipment
func (h *TradeHandler) GetTradeDocumentsByShipment(c echo.Context) error {
	shipmentID := c.Param("id")

	params := map[string]interface{}{
		"resource_type": "shipment",
		"resource_id":   shipmentID,
	}

	docs, err := h.tradeService.ListTradeDocuments(c.Request().Context(), params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get documents"})
	}

	return c.JSON(http.StatusOK, docs)
}

// CreateShipmentEvent creates a new shipment event
func (h *TradeHandler) CreateShipmentEvent(c echo.Context) error {
	shipmentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid shipment ID"})
	}

	var event models.ShipmentEvent
	if err := c.Bind(&event); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	event.ShipmentID = shipmentID
	event.CompanyID = c.Get("company_id").(uuid.UUID)

	if err := h.tradeService.CreateShipmentEvent(c.Request().Context(), &event); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create shipment event"})
	}

	return c.JSON(http.StatusCreated, event)
}

// GetShipmentEvents gets events for a shipment
func (h *TradeHandler) GetShipmentEvents(c echo.Context) error {
	shipmentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid shipment ID"})
	}

	events, err := h.tradeService.GetShipmentEvents(c.Request().Context(), shipmentID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get shipment events"})
	}

	return c.JSON(http.StatusOK, events)
}

// Letter of Credit Management

// ListLetterOfCredits lists letters of credit with pagination and filtering
func (h *TradeHandler) ListLetterOfCredits(c echo.Context) error {
	params := make(map[string]interface{})
	
	// Get company ID from context
	params["company_id"] = c.Get("company_id").(uuid.UUID)
	
	// Parse query parameters
	if status := c.QueryParam("status"); status != "" {
		params["status"] = status
	}
	if lcType := c.QueryParam("type"); lcType != "" {
		params["type"] = lcType
	}
	if search := c.QueryParam("search"); search != "" {
		params["search"] = search
	}
	
	// Pagination
	page := 1
	pageSize := 20
	if p := c.QueryParam("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if ps := c.QueryParam("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 {
			pageSize = parsed
		}
	}
	params["page"] = page
	params["page_size"] = pageSize

	lcs, total, err := h.tradeService.ListLettersOfCredit(c.Request().Context(), params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to list letters of credit"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": lcs,
		"pagination": map[string]interface{}{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// CreateLetterOfCredit creates a new letter of credit
func (h *TradeHandler) CreateLetterOfCredit(c echo.Context) error {
	var lc models.LetterOfCredit
	if err := c.Bind(&lc); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	lc.CompanyID = c.Get("company_id").(uuid.UUID)

	if err := h.tradeService.CreateLetterOfCredit(c.Request().Context(), &lc); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create letter of credit"})
	}

	return c.JSON(http.StatusCreated, lc)
}

// GetLetterOfCredit retrieves a letter of credit by ID
func (h *TradeHandler) GetLetterOfCredit(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid letter of credit ID"})
	}

	lc, err := h.tradeService.GetLetterOfCredit(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get letter of credit"})
	}

	return c.JSON(http.StatusOK, lc)
}

// GetExpiringLetterOfCredits gets expiring letters of credit
func (h *TradeHandler) GetExpiringLetterOfCredits(c echo.Context) error {
	days := 30
	if d := c.QueryParam("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil && parsed > 0 {
			days = parsed
		}
	}

	params := map[string]interface{}{
		"company_id": c.Get("company_id").(uuid.UUID),
		"status":     "issued",
		"expiry_before": time.Now().AddDate(0, 0, days),
	}

	lcs, _, err := h.tradeService.ListLettersOfCredit(c.Request().Context(), params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get expiring letters of credit"})
	}

	return c.JSON(http.StatusOK, lcs)
}

// CreateLCUtilization creates a new LC utilization
func (h *TradeHandler) CreateLCUtilization(c echo.Context) error {
	lcID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid letter of credit ID"})
	}

	var util models.LCUtilization
	if err := c.Bind(&util); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	util.LCID = lcID
	util.CompanyID = c.Get("company_id").(uuid.UUID)

	// TODO: Implement LC utilization in service
	return c.JSON(http.StatusNotImplemented, map[string]string{"error": "Not implemented"})
}

// GetLCUtilizations gets utilizations for a letter of credit
func (h *TradeHandler) GetLCUtilizations(c echo.Context) error {
	lcID := c.Param("id")
	
	// TODO: Implement LC utilization retrieval in service
	return c.JSON(http.StatusNotImplemented, map[string]string{"error": "Not implemented", "lc_id": lcID})
}

// Compliance Management

// RunComplianceCheck runs a compliance check
func (h *TradeHandler) RunComplianceCheck(c echo.Context) error {
	var req struct {
		ResourceType string `json:"resource_type"`
		ResourceID   string `json:"resource_id"`
	}
	
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	checkResult := map[string]interface{}{
		"status": "pending",
		"message": "Compliance check initiated",
		"company_id": c.Get("company_id").(uuid.UUID).String(),
	}
	
	check := &models.TradeComplianceCheck{
		ResourceType: req.ResourceType,
		ResourceID:   req.ResourceID,
		CheckType:    "manual",
		Status:       "pending",
		Result:       models.JSONB(checkResult),
		CheckedAt:    time.Now(),
	}

	if err := h.tradeService.CreateComplianceCheck(c.Request().Context(), check); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create compliance check"})
	}

	return c.JSON(http.StatusCreated, check)
}

// GetComplianceChecksByResource gets compliance checks for a resource
func (h *TradeHandler) GetComplianceChecksByResource(c echo.Context) error {
	resourceType := c.QueryParam("resource_type")
	resourceID := c.QueryParam("resource_id")

	if resourceType == "" || resourceID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "resource_type and resource_id are required"})
	}

	checks, err := h.tradeService.GetComplianceChecks(c.Request().Context(), resourceType, resourceID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get compliance checks"})
	}

	return c.JSON(http.StatusOK, checks)
}

// GetFailedComplianceChecks gets failed compliance checks
func (h *TradeHandler) GetFailedComplianceChecks(c echo.Context) error {
	// TODO: Implement failed compliance checks retrieval
	return c.JSON(http.StatusNotImplemented, map[string]string{"error": "Not implemented"})
}

// Exchange Rate Management

// ListExchangeRates lists exchange rates
func (h *TradeHandler) ListExchangeRates(c echo.Context) error {
	// TODO: Implement exchange rate listing
	return c.JSON(http.StatusNotImplemented, map[string]string{"error": "Not implemented"})
}

// CreateExchangeRate creates a new exchange rate
func (h *TradeHandler) CreateExchangeRate(c echo.Context) error {
	var rate models.ExchangeRate
	if err := c.Bind(&rate); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := h.tradeService.CreateExchangeRate(c.Request().Context(), &rate); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create exchange rate"})
	}

	return c.JSON(http.StatusCreated, rate)
}

// GetLatestExchangeRate gets the latest exchange rate
func (h *TradeHandler) GetLatestExchangeRate(c echo.Context) error {
	from := c.QueryParam("from")
	to := c.QueryParam("to")

	if from == "" || to == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "from and to currencies are required"})
	}

	rate, err := h.tradeService.GetLatestExchangeRate(c.Request().Context(), from, to)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get exchange rate"})
	}

	return c.JSON(http.StatusOK, rate)
}

// Statistics and Analytics

// GetTradeStatistics gets trade statistics
func (h *TradeHandler) GetTradeStatistics(c echo.Context) error {
	// TODO: Implement trade statistics
	return c.JSON(http.StatusNotImplemented, map[string]string{"error": "Not implemented"})
}

// GetShipmentsByCountry gets shipments grouped by country
func (h *TradeHandler) GetShipmentsByCountry(c echo.Context) error {
	// TODO: Implement shipments by country
	return c.JSON(http.StatusNotImplemented, map[string]string{"error": "Not implemented"})
}

// GetTopTradingPartners gets top trading partners
func (h *TradeHandler) GetTopTradingPartners(c echo.Context) error {
	// TODO: Implement top trading partners
	return c.JSON(http.StatusNotImplemented, map[string]string{"error": "Not implemented"})
}

// Calculation Tools

// CalculateTariffDuty calculates tariff duty
func (h *TradeHandler) CalculateTariffDuty(c echo.Context) error {
	var params models.TariffCalculationParams
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	params.CompanyID = c.Get("company_id").(uuid.UUID)

	result, err := h.tradeService.CalculateTariff(c.Request().Context(), params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to calculate tariff: " + err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

// ConvertCurrency converts between currencies
func (h *TradeHandler) ConvertCurrency(c echo.Context) error {
	var req struct {
		Amount float64 `json:"amount"`
		From   string  `json:"from"`
		To     string  `json:"to"`
	}
	
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	rate, err := h.tradeService.GetLatestExchangeRate(c.Request().Context(), req.From, req.To)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get exchange rate"})
	}

	convertedAmount := req.Amount * rate.Rate

	return c.JSON(http.StatusOK, map[string]interface{}{
		"amount":           req.Amount,
		"from":             req.From,
		"to":               req.To,
		"rate":             rate.Rate,
		"converted_amount": convertedAmount,
		"rate_date":        rate.EffectiveDate,
	})
}