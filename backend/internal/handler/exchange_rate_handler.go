package handler

import (
	"net/http"
	"time"

	"github.com/fastenmind/fastener-api/internal/models"
	"github.com/fastenmind/fastener-api/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ExchangeRateHandler struct {
	webhookService *service.WebhookService
	db             *gorm.DB
}

func NewExchangeRateHandler(db *gorm.DB, webhookService *service.WebhookService) *ExchangeRateHandler {
	return &ExchangeRateHandler{
		db:             db,
		webhookService: webhookService,
	}
}

// BatchUpdateExchangeRates handles batch update of exchange rates from N8N
// @Summary Batch update exchange rates
// @Description Update multiple exchange rates at once
// @Tags Exchange Rates
// @Accept json
// @Produce json
// @Param request body BatchUpdateExchangeRatesRequest true "Exchange rates data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/exchange-rates/batch-update [post]
func (h *ExchangeRateHandler) BatchUpdateExchangeRates(c echo.Context) error {
	var req BatchUpdateExchangeRatesRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	companyID := c.Get("company_id").(uuid.UUID)
	userID := c.Get("user_id").(uuid.UUID)

	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	updatedRates := []map[string]interface{}{}
	
	for _, rate := range req.Rates {
		// Check if rate exists
		var existingRate models.ExchangeRate
		err := tx.Where("currency = ? AND base_currency = ?", rate.Currency, rate.Base).
			First(&existingRate).Error
		
		if err == nil {
			// Update existing rate
			oldRate := existingRate.Rate
			existingRate.Rate = rate.Rate
			existingRate.UpdatedAt = time.Now()
			
			if err := tx.Save(&existingRate).Error; err != nil {
				tx.Rollback()
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update rate"})
			}
			
			// Trigger webhook if rate changed significantly (> 2%)
			changePercent := ((rate.Rate - oldRate) / oldRate) * 100
			if changePercent > 2 || changePercent < -2 {
				go h.webhookService.TriggerExchangeRateUpdated(
					rate.Currency, oldRate, rate.Rate, companyID, userID,
				)
			}
			
			updatedRates = append(updatedRates, map[string]interface{}{
				"currency": rate.Currency,
				"old_rate": oldRate,
				"new_rate": rate.Rate,
				"changed":  oldRate != rate.Rate,
			})
		} else {
			// Create new rate
			newRate := &models.ExchangeRate{
				Currency:     rate.Currency,
				BaseCurrency: rate.Base,
				Rate:         rate.Rate,
				ValidFrom:    time.Now(),
				CompanyID:    companyID,
				CreatedBy:    userID,
			}
			
			if err := tx.Create(newRate).Error; err != nil {
				tx.Rollback()
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create rate"})
			}
			
			updatedRates = append(updatedRates, map[string]interface{}{
				"currency": rate.Currency,
				"new_rate": rate.Rate,
				"created":  true,
			})
		}
	}

	if err := tx.Commit().Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to commit transaction"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":       "Exchange rates updated successfully",
		"updated_count": len(updatedRates),
		"rates":         updatedRates,
		"timestamp":     time.Now(),
	})
}

// GetOutdatedCosts returns costs that haven't been updated in the specified period
// @Summary Get outdated costs
// @Description Get costs that need updating
// @Tags Costs
// @Produce json
// @Param days query int false "Days threshold" default(30)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/cost-calculations/outdated [get]
func (h *ExchangeRateHandler) GetOutdatedCosts(c echo.Context) error {
	days := 30 // Default to 30 days
	if daysParam := c.QueryParam("days"); daysParam != "" {
		// Parse days parameter
	}

	threshold := time.Now().AddDate(0, 0, -days)
	
	var outdatedCosts []map[string]interface{}
	
	// Material costs
	var materials []models.Material
	h.db.Where("updated_at < ?", threshold).Find(&materials)
	for _, m := range materials {
		outdatedCosts = append(outdatedCosts, map[string]interface{}{
			"id":           m.ID,
			"name":         m.Name,
			"cost_type":    "material",
			"last_updated": m.UpdatedAt,
			"days_old":     int(time.Since(m.UpdatedAt).Hours() / 24),
		})
	}
	
	// Process costs
	var processCosts []models.ProcessCostConfig
	h.db.Where("updated_at < ?", threshold).Find(&processCosts)
	for _, p := range processCosts {
		outdatedCosts = append(outdatedCosts, map[string]interface{}{
			"id":           p.ID,
			"name":         p.ProcessID,
			"cost_type":    "process",
			"last_updated": p.UpdatedAt,
			"days_old":     int(time.Since(p.UpdatedAt).Hours() / 24),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":      outdatedCosts,
		"total":     len(outdatedCosts),
		"threshold": threshold,
	})
}

// RegisterRoutes registers all exchange rate routes
func (h *ExchangeRateHandler) RegisterRoutes(e *echo.Group) {
	e.POST("/exchange-rates/batch-update", h.BatchUpdateExchangeRates)
	e.GET("/cost-calculations/outdated", h.GetOutdatedCosts)
}

type BatchUpdateExchangeRatesRequest struct {
	Rates []ExchangeRateData `json:"rates" validate:"required"`
}

type ExchangeRateData struct {
	Currency string  `json:"currency" validate:"required"`
	Rate     float64 `json:"rate" validate:"required,min=0"`
	Base     string  `json:"base" validate:"required"`
	Date     string  `json:"date"`
}