package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/fastenmind/fastener-api/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TariffHandler struct {
	service service.TariffService
}

func NewTariffHandler(service service.TariffService) *TariffHandler {
	return &TariffHandler{service: service}
}

// SearchHSCodes godoc
// @Summary Search HS codes
// @Description Search harmonized system codes
// @Tags Tariff
// @Accept json
// @Produce json
// @Param search query string false "Search term"
// @Param category query string false "Category filter"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} map[string]interface{}
// @Router /api/tariffs/hs-codes [get]
func (h *TariffHandler) SearchHSCodes(c echo.Context) error {
	params := make(map[string]interface{})
	
	if search := c.QueryParam("search"); search != "" {
		params["search"] = search
	}
	if category := c.QueryParam("category"); category != "" {
		params["category"] = category
	}
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
	
	params["is_active"] = true
	
	codes, total, err := h.service.SearchHSCodes(params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": codes,
		"pagination": map[string]interface{}{
			"total":     total,
			"page":      params["page"],
			"page_size": params["page_size"],
		},
	})
}

// GetHSCode godoc
// @Summary Get HS code details
// @Description Get details of a specific HS code
// @Tags Tariff
// @Accept json
// @Produce json
// @Param code path string true "HS Code"
// @Success 200 {object} models.HSCode
// @Router /api/tariffs/hs-codes/{code} [get]
func (h *TariffHandler) GetHSCode(c echo.Context) error {
	code := c.Param("code")
	
	hsCode, err := h.service.GetHSCode(code)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "HS code not found")
	}
	
	return c.JSON(http.StatusOK, hsCode)
}

// GetTariffRates godoc
// @Summary Get tariff rates
// @Description Get tariff rates for specific parameters
// @Tags Tariff
// @Accept json
// @Produce json
// @Param hs_code query string true "HS Code"
// @Param from_country query string false "Export country"
// @Param to_country query string false "Import country"
// @Success 200 {array} models.TariffRate
// @Router /api/tariffs/rates [get]
func (h *TariffHandler) GetTariffRates(c echo.Context) error {
	// This would typically query the tariff rates
	// For now, return empty array
	return c.JSON(http.StatusOK, []interface{}{})
}

// CalculateTariff godoc
// @Summary Calculate tariff
// @Description Calculate tariff for given parameters
// @Tags Tariff
// @Accept json
// @Produce json
// @Param request body service.TariffCalculationRequest true "Calculation request"
// @Success 200 {object} service.TariffCalculationResult
// @Router /api/tariffs/calculate [post]
func (h *TariffHandler) CalculateTariff(c echo.Context) error {
	var req service.TariffCalculationRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	// Get user context
	companyID := getCompanyIDFromContext(c)
	userID := getUserIDFromContext(c)
	req.CompanyID = companyID
	req.UserID = userID
	
	result, err := h.service.CalculateTariff(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusOK, result)
}

// BatchCalculateTariff godoc
// @Summary Batch calculate tariffs
// @Description Calculate tariffs for multiple items
// @Tags Tariff
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "Batch calculation request"
// @Success 200 {array} service.TariffCalculationResult
// @Router /api/tariffs/batch-calculate [post]
func (h *TariffHandler) BatchCalculateTariff(c echo.Context) error {
	var reqBody struct {
		Items []service.TariffCalculationRequest `json:"items"`
	}
	
	if err := c.Bind(&reqBody); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	// Get user context
	companyID := getCompanyIDFromContext(c)
	userID := getUserIDFromContext(c)
	
	// Set user context for each item
	for i := range reqBody.Items {
		reqBody.Items[i].CompanyID = companyID
		reqBody.Items[i].UserID = userID
	}
	
	results, err := h.service.BatchCalculateTariff(reqBody.Items)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, results)
}

// GetTradeAgreements godoc
// @Summary Get trade agreements
// @Description Get trade agreements between countries
// @Tags Tariff
// @Accept json
// @Produce json
// @Param countries query string true "Comma-separated country codes"
// @Success 200 {array} models.TradeAgreement
// @Router /api/tariffs/trade-agreements [get]
func (h *TariffHandler) GetTradeAgreements(c echo.Context) error {
	countriesParam := c.QueryParam("countries")
	if countriesParam == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "countries parameter required")
	}
	
	countries := strings.Split(countriesParam, ",")
	
	agreements, err := h.service.GetTradeAgreements(countries)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, agreements)
}

// ValidateHSCode godoc
// @Summary Validate HS code
// @Description Validate HS code for a specific country
// @Tags Tariff
// @Accept json
// @Produce json
// @Param request body map[string]string true "Validation request"
// @Success 200 {object} map[string]interface{}
// @Router /api/tariffs/validate-hs-code [post]
func (h *TariffHandler) ValidateHSCode(c echo.Context) error {
	var req struct {
		Code    string `json:"code"`
		Country string `json:"country"`
	}
	
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	valid, message, suggestions, err := h.service.ValidateHSCode(req.Code, req.Country)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	response := map[string]interface{}{
		"valid":   valid,
		"message": message,
	}
	
	if suggestions != nil {
		response["suggested_codes"] = suggestions
	}
	
	return c.JSON(http.StatusOK, response)
}

// GetCommonHSCodes godoc
// @Summary Get common HS codes
// @Description Get commonly used HS codes
// @Tags Tariff
// @Accept json
// @Produce json
// @Param category query string false "Category filter"
// @Success 200 {array} models.HSCode
// @Router /api/tariffs/common-hs-codes [get]
func (h *TariffHandler) GetCommonHSCodes(c echo.Context) error {
	category := c.QueryParam("category")
	
	codes, err := h.service.GetCommonHSCodes(category)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, codes)
}

// RegisterRoutes registers all tariff routes
func (h *TariffHandler) RegisterRoutes(e *echo.Echo, authMiddleware echo.MiddlewareFunc) {
	tariff := e.Group("/api/tariffs", authMiddleware)
	
	tariff.GET("/hs-codes", h.SearchHSCodes)
	tariff.GET("/hs-codes/:code", h.GetHSCode)
	tariff.GET("/rates", h.GetTariffRates)
	tariff.POST("/calculate", h.CalculateTariff)
	tariff.POST("/batch-calculate", h.BatchCalculateTariff)
	tariff.GET("/trade-agreements", h.GetTradeAgreements)
	tariff.POST("/validate-hs-code", h.ValidateHSCode)
	tariff.GET("/common-hs-codes", h.GetCommonHSCodes)
}