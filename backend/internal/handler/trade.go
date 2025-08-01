package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/fastenmind/fastener-api/internal/service"
)

type TradeHandler struct {
	tradeService *service.TradeService
}

func NewTradeHandler(tradeService *service.TradeService) *TradeHandler {
	return &TradeHandler{
		tradeService: tradeService,
	}
}

// TariffCode Handlers
func (h *TradeHandler) CreateTariffCode(c echo.Context) error {
	userID := getUserIDFromContext(c)

	var req service.CreateTariffCodeRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	tariffCode, err := h.tradeService.CreateTariffCode(userID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, tariffCode)
}

func (h *TradeHandler) GetTariffCode(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid tariff code ID")
	}

	tariffCode, err := h.tradeService.GetTariffCode(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Tariff code not found")
	}

	return c.JSON(http.StatusOK, tariffCode)
}

func (h *TradeHandler) ListTariffCodes(c echo.Context) error {
	companyID := getCompanyIDFromContext(c)

	hsCode := c.QueryParam("hs_code")
	category := c.QueryParam("category")
	isActiveStr := c.QueryParam("is_active")

	var isActive *bool
	if isActiveStr != "" {
		activeVal := isActiveStr == "true"
		isActive = &activeVal
	}

	tariffCodes, err := h.tradeService.GetTariffCodesByCompany(companyID, hsCode, category, isActive)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": tariffCodes,
	})
}

func (h *TradeHandler) UpdateTariffCode(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid tariff code ID")
	}

	userID := getUserIDFromContext(c)

	var req service.UpdateTariffCodeRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	tariffCode, err := h.tradeService.UpdateTariffCode(id, userID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, tariffCode)
}

func (h *TradeHandler) DeleteTariffCode(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid tariff code ID")
	}

	err = h.tradeService.DeleteTariffCode(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

// TariffRate Handlers
func (h *TradeHandler) CreateTariffRate(c echo.Context) error {
	userID := getUserIDFromContext(c)

	var req service.CreateTariffRateRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	tariffRate, err := h.tradeService.CreateTariffRate(userID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, tariffRate)
}

func (h *TradeHandler) GetTariffRatesByTariffCode(c echo.Context) error {
	tariffCodeID, err := uuid.Parse(c.Param("tariff_code_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid tariff code ID")
	}

	countryCode := c.QueryParam("country_code")
	tradeType := c.QueryParam("trade_type")
	agreementType := c.QueryParam("agreement_type")

	tariffRates, err := h.tradeService.GetTariffRatesByTariffCode(tariffCodeID, countryCode, tradeType, agreementType)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": tariffRates,
	})
}

func (h *TradeHandler) ListTariffRates(c echo.Context) error {
	companyID := getCompanyIDFromContext(c)

	countryCode := c.QueryParam("country_code")
	tradeType := c.QueryParam("trade_type")

	tariffRates, err := h.tradeService.GetTariffRatesByCompany(companyID, countryCode, tradeType)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": tariffRates,
	})
}

// Shipment Handlers
func (h *TradeHandler) CreateShipment(c echo.Context) error {
	userID := getUserIDFromContext(c)

	var req service.CreateShipmentRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	shipment, err := h.tradeService.CreateShipment(userID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, shipment)
}

func (h *TradeHandler) GetShipment(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid shipment ID")
	}

	shipment, err := h.tradeService.GetShipment(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Shipment not found")
	}

	return c.JSON(http.StatusOK, shipment)
}

func (h *TradeHandler) ListShipments(c echo.Context) error {
	companyID := getCompanyIDFromContext(c)

	shipmentType := c.QueryParam("type")
	status := c.QueryParam("status")
	method := c.QueryParam("method")

	shipments, err := h.tradeService.GetShipmentsByCompany(companyID, shipmentType, status, method)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": shipments,
	})
}

func (h *TradeHandler) UpdateShipment(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid shipment ID")
	}

	userID := getUserIDFromContext(c)

	var req service.UpdateShipmentRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	shipment, err := h.tradeService.UpdateShipment(id, userID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, shipment)
}

// ShipmentEvent Handlers
func (h *TradeHandler) CreateShipmentEvent(c echo.Context) error {
	shipmentID, err := uuid.Parse(c.Param("shipment_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid shipment ID")
	}

	userID := getUserIDFromContext(c)

	var req service.CreateShipmentEventRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	event, err := h.tradeService.CreateShipmentEvent(shipmentID, userID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, event)
}

func (h *TradeHandler) GetShipmentEvents(c echo.Context) error {
	shipmentID, err := uuid.Parse(c.Param("shipment_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid shipment ID")
	}

	events, err := h.tradeService.GetShipmentEventsByShipment(shipmentID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": events,
	})
}

// LetterOfCredit Handlers
func (h *TradeHandler) CreateLetterOfCredit(c echo.Context) error {
	userID := getUserIDFromContext(c)

	var req service.CreateLetterOfCreditRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	lc, err := h.tradeService.CreateLetterOfCredit(userID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, lc)
}

func (h *TradeHandler) GetLetterOfCredit(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid LC ID")
	}

	lc, err := h.tradeService.GetLetterOfCredit(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Letter of Credit not found")
	}

	return c.JSON(http.StatusOK, lc)
}

func (h *TradeHandler) ListLetterOfCredits(c echo.Context) error {
	companyID := getCompanyIDFromContext(c)

	lcType := c.QueryParam("type")
	status := c.QueryParam("status")

	lcs, err := h.tradeService.GetLetterOfCreditsByCompany(companyID, lcType, status)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": lcs,
	})
}

func (h *TradeHandler) GetExpiringLetterOfCredits(c echo.Context) error {
	companyID := getCompanyIDFromContext(c)

	daysStr := c.QueryParam("days")
	days := 30 // default
	if daysStr != "" {
		if parsedDays, err := strconv.Atoi(daysStr); err == nil {
			days = parsedDays
		}
	}

	lcs, err := h.tradeService.GetExpiringLetterOfCredits(companyID, days)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": lcs,
	})
}

// LCUtilization Handlers
func (h *TradeHandler) CreateLCUtilization(c echo.Context) error {
	userID := getUserIDFromContext(c)

	var req service.CreateLCUtilizationRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	utilization, err := h.tradeService.CreateLCUtilization(userID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, utilization)
}

func (h *TradeHandler) GetLCUtilizations(c echo.Context) error {
	lcID, err := uuid.Parse(c.Param("lc_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid LC ID")
	}

	status := c.QueryParam("status")

	utilizations, err := h.tradeService.GetLCUtilizationsByLC(lcID, status)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": utilizations,
	})
}

// Compliance Handlers
func (h *TradeHandler) RunComplianceCheck(c echo.Context) error {
	userID := getUserIDFromContext(c)

	var req service.CreateComplianceCheckRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	check, err := h.tradeService.RunComplianceCheck(userID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, check)
}

func (h *TradeHandler) GetComplianceChecksByResource(c echo.Context) error {
	companyID := getCompanyIDFromContext(c)

	resourceType := c.QueryParam("resource_type")
	if resourceType == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Resource type is required")
	}

	resourceIDStr := c.QueryParam("resource_id")
	if resourceIDStr == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Resource ID is required")
	}

	resourceID, err := uuid.Parse(resourceIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid resource ID")
	}

	checks, err := h.tradeService.GetComplianceChecksByResource(companyID, resourceType, resourceID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": checks,
	})
}

func (h *TradeHandler) GetFailedComplianceChecks(c echo.Context) error {
	companyID := getCompanyIDFromContext(c)

	checks, err := h.tradeService.GetFailedComplianceChecks(companyID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": checks,
	})
}

// ExchangeRate Handlers
func (h *TradeHandler) CreateExchangeRate(c echo.Context) error {
	userID := getUserIDFromContext(c)

	var req service.CreateExchangeRateRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	rate, err := h.tradeService.CreateExchangeRate(userID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, rate)
}

func (h *TradeHandler) GetLatestExchangeRate(c echo.Context) error {
	companyID := getCompanyIDFromContext(c)

	fromCurrency := c.QueryParam("from")
	toCurrency := c.QueryParam("to")
	rateType := c.QueryParam("type")

	if fromCurrency == "" || toCurrency == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "From and to currencies are required")
	}

	if rateType == "" {
		rateType = "mid"
	}

	rate, err := h.tradeService.GetLatestExchangeRate(companyID, fromCurrency, toCurrency, rateType)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Exchange rate not found")
	}

	return c.JSON(http.StatusOK, rate)
}

func (h *TradeHandler) ListExchangeRates(c echo.Context) error {
	companyID := getCompanyIDFromContext(c)

	fromCurrency := c.QueryParam("from")
	toCurrency := c.QueryParam("to")
	rateType := c.QueryParam("type")

	rates, err := h.tradeService.GetExchangeRatesByCompany(companyID, fromCurrency, toCurrency, rateType)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": rates,
	})
}

// Analytics Handlers
func (h *TradeHandler) GetTradeStatistics(c echo.Context) error {
	companyID := getCompanyIDFromContext(c)

	startDateStr := c.QueryParam("start_date")
	endDateStr := c.QueryParam("end_date")

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid start_date format")
		}
	} else {
		startDate = time.Now().AddDate(0, 0, -30) // Last 30 days
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid end_date format")
		}
	} else {
		endDate = time.Now()
	}

	stats, err := h.tradeService.GetTradeStatistics(companyID, startDate, endDate)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, stats)
}

func (h *TradeHandler) GetShipmentsByCountry(c echo.Context) error {
	companyID := getCompanyIDFromContext(c)

	startDateStr := c.QueryParam("start_date")
	endDateStr := c.QueryParam("end_date")

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid start_date format")
		}
	} else {
		startDate = time.Now().AddDate(0, 0, -30)
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid end_date format")
		}
	} else {
		endDate = time.Now()
	}

	data, err := h.tradeService.GetShipmentsByCountry(companyID, startDate, endDate)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": data,
	})
}

func (h *TradeHandler) GetTopTradingPartners(c echo.Context) error {
	companyID := getCompanyIDFromContext(c)

	startDateStr := c.QueryParam("start_date")
	endDateStr := c.QueryParam("end_date")
	limitStr := c.QueryParam("limit")

	var startDate, endDate time.Time
	var err error
	limit := 10 // default

	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid start_date format")
		}
	} else {
		startDate = time.Now().AddDate(0, 0, -30)
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid end_date format")
		}
	} else {
		endDate = time.Now()
	}

	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	data, err := h.tradeService.GetTopTradingPartners(companyID, startDate, endDate, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": data,
	})
}

// Utility Handlers
func (h *TradeHandler) CalculateTariffDuty(c echo.Context) error {
	companyID := getCompanyIDFromContext(c)

	var req struct {
		HSCode      string  `json:"hs_code" validate:"required"`
		CountryCode string  `json:"country_code" validate:"required"`
		TradeType   string  `json:"trade_type" validate:"required"`
		Value       float64 `json:"value" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	result, err := h.tradeService.CalculateTariffDuty(companyID, req.HSCode, req.CountryCode, req.TradeType, req.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

func (h *TradeHandler) ConvertCurrency(c echo.Context) error {
	companyID := getCompanyIDFromContext(c)

	var req struct {
		Amount       float64 `json:"amount" validate:"required"`
		FromCurrency string  `json:"from_currency" validate:"required"`
		ToCurrency   string  `json:"to_currency" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	result, err := h.tradeService.ConvertCurrency(companyID, req.Amount, req.FromCurrency, req.ToCurrency)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

func (h *TradeHandler) GetTradeDocumentsByShipment(c echo.Context) error {
	shipmentID, err := uuid.Parse(c.Param("shipment_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid shipment ID")
	}

	documents, err := h.tradeService.GetTradeDocumentsByShipment(shipmentID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": documents,
	})
}