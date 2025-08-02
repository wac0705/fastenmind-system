package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/fastenmind/fastener-api/internal/config"
	"github.com/fastenmind/fastener-api/internal/handler"
	"github.com/fastenmind/fastener-api/internal/middleware"
	"github.com/fastenmind/fastener-api/internal/repository"
	"github.com/fastenmind/fastener-api/internal/service"
	"github.com/fastenmind/fastener-api/pkg/database"
	"github.com/fastenmind/fastener-api/pkg/logger"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize configuration
	cfg := config.New()

	// Initialize logger
	log := logger.New(cfg.Server.Environment)

	// Initialize database wrapper
	dbWrapper, err := database.NewWrapper(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer dbWrapper.Close()

	// Initialize Echo
	e := echo.New()
	e.HideBanner = true
	e.Logger = log

	// Middleware
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.RequestID())
	e.Use(middleware.CORS(cfg.CORS))
	e.Use(middleware.Security())

	// Initialize repositories
	repos := repository.NewRepositories(dbWrapper.GormDB)

	// Initialize services
	services := service.NewServices(repos, cfg)

	// Initialize handlers
	h := handler.NewHandlers(services)

	// Routes
	setupRoutes(e, h, cfg)

	// Start server
	go func() {
		addr := fmt.Sprintf(":%s", cfg.Server.Port)
		log.Info("Starting server on " + addr)
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	if err := e.Shutdown(ctx); err != nil {
		log.Fatal("Failed to shutdown server:", err)
	}
	
	log.Info("Server stopped")
}

func setupRoutes(e *echo.Echo, h *handler.Handlers, cfg *config.Config) {
	// API Group
	api := e.Group("/api/v1")

	// Health check
	api.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "ok",
			"time":   time.Now(),
		})
	})

	// Public routes
	public := api.Group("")
	{
		public.POST("/auth/login", h.Auth.Login)
		public.POST("/auth/refresh", h.Auth.RefreshToken)
		public.POST("/auth/register", h.Auth.Register)
	}

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.JWT(cfg.JWT.SecretKey))
	{
		// Account routes
		protected.GET("/accounts", h.Account.List)
		protected.GET("/accounts/:id", h.Account.Get)
		protected.PUT("/accounts/:id", h.Account.Update)
		protected.DELETE("/accounts/:id", h.Account.Delete)
		protected.PUT("/accounts/:id/password", h.Account.ChangePassword)

		// Company routes
		protected.GET("/companies", h.Company.List)
		protected.POST("/companies", h.Company.Create)
		protected.GET("/companies/:id", h.Company.Get)
		protected.PUT("/companies/:id", h.Company.Update)
		protected.DELETE("/companies/:id", h.Company.Delete)

		// Customer routes
		protected.GET("/customers", h.Customer.List)
		protected.POST("/customers", h.Customer.Create)
		protected.GET("/customers/:id", h.Customer.Get)
		protected.PUT("/customers/:id", h.Customer.Update)
		protected.DELETE("/customers/:id", h.Customer.Delete)

		// Inquiry routes
		protected.GET("/inquiries", h.Inquiry.List)
		protected.POST("/inquiries", h.Inquiry.Create)
		protected.GET("/inquiries/:id", h.Inquiry.Get)
		protected.PUT("/inquiries/:id", h.Inquiry.Update)
		protected.DELETE("/inquiries/:id", h.Inquiry.Delete)
		protected.POST("/inquiries/:id/assign", h.Inquiry.AssignEngineer)
		protected.POST("/inquiries/:id/quote", h.Inquiry.CreateQuote)

		// Process routes
		protected.GET("/processes", h.Process.List)
		protected.POST("/processes", h.Process.Create)
		protected.GET("/processes/:id", h.Process.Get)
		protected.PUT("/processes/:id", h.Process.Update)
		protected.DELETE("/processes/:id", h.Process.Delete)

		// Equipment routes
		protected.GET("/equipment", h.Equipment.List)
		protected.POST("/equipment", h.Equipment.Create)
		protected.GET("/equipment/:id", h.Equipment.Get)
		protected.PUT("/equipment/:id", h.Equipment.Update)
		protected.DELETE("/equipment/:id", h.Equipment.Delete)

		// Assignment rules routes
		protected.GET("/assignment-rules", h.AssignmentRule.List)
		protected.POST("/assignment-rules", h.AssignmentRule.Create)
		protected.GET("/assignment-rules/:id", h.AssignmentRule.Get)
		protected.PUT("/assignment-rules/:id", h.AssignmentRule.Update)
		protected.DELETE("/assignment-rules/:id", h.AssignmentRule.Delete)

		// Tariff routes
		h.Tariff.RegisterRoutes(e, middleware.JWT(cfg.JWT.SecretKey))
		
		// N8N routes
		h.N8N.RegisterRoutes(e, middleware.JWT(cfg.JWT.SecretKey))
		
		// Quote routes
		h.Quote.RegisterRoutes(e, middleware.JWT(cfg.JWT.SecretKey))
		
		// Order routes
		h.Order.RegisterRoutes(e, middleware.JWT(cfg.JWT.SecretKey))
		
		// Inventory routes
		h.Inventory.RegisterRoutes(e, middleware.JWT(cfg.JWT.SecretKey))

		// Compliance routes
		protected.POST("/compliance/check", h.Compliance.Check)
		protected.GET("/compliance/rules", h.Compliance.GetRules)
		protected.GET("/compliance/documents", h.Compliance.GetDocumentRequirements)
		protected.POST("/compliance/validate-documents", h.Compliance.ValidateDocuments)

		// Trade routes
		protected.GET("/trade/tariff-codes", h.Trade.ListTariffCodes)
		protected.POST("/trade/tariff-codes", h.Trade.CreateTariffCode)
		protected.GET("/trade/tariff-codes/:id", h.Trade.GetTariffCode)
		protected.PUT("/trade/tariff-codes/:id", h.Trade.UpdateTariffCode)
		protected.DELETE("/trade/tariff-codes/:id", h.Trade.DeleteTariffCode)

		// Tariff rates
		protected.POST("/trade/tariff-rates", h.Trade.CreateTariffRate)
		protected.GET("/trade/tariff-codes/:tariff_code_id/rates", h.Trade.GetTariffRatesByTariffCode)
		protected.GET("/trade/tariff-rates", h.Trade.ListTariffRates)

		// Shipments
		protected.GET("/trade/shipments", h.Trade.ListShipments)
		protected.POST("/trade/shipments", h.Trade.CreateShipment)
		protected.GET("/trade/shipments/:id", h.Trade.GetShipment)
		protected.PUT("/trade/shipments/:id", h.Trade.UpdateShipment)
		protected.GET("/trade/shipments/:shipment_id/documents", h.Trade.GetTradeDocumentsByShipment)

		// Shipment events
		protected.POST("/trade/shipments/:shipment_id/events", h.Trade.CreateShipmentEvent)
		protected.GET("/trade/shipments/:shipment_id/events", h.Trade.GetShipmentEvents)

		// Letter of Credits
		protected.GET("/trade/letter-of-credits", h.Trade.ListLetterOfCredits)
		protected.POST("/trade/letter-of-credits", h.Trade.CreateLetterOfCredit)
		protected.GET("/trade/letter-of-credits/:id", h.Trade.GetLetterOfCredit)
		protected.GET("/trade/letter-of-credits/expiring", h.Trade.GetExpiringLetterOfCredits)

		// LC Utilizations
		protected.POST("/trade/lc-utilizations", h.Trade.CreateLCUtilization)
		protected.GET("/trade/letter-of-credits/:lc_id/utilizations", h.Trade.GetLCUtilizations)

		// Trade Compliance
		protected.POST("/trade/compliance/check", h.Trade.RunComplianceCheck)
		protected.GET("/trade/compliance/checks", h.Trade.GetComplianceChecksByResource)
		protected.GET("/trade/compliance/failed-checks", h.Trade.GetFailedComplianceChecks)

		// Exchange rates
		protected.GET("/trade/exchange-rates", h.Trade.ListExchangeRates)
		protected.POST("/trade/exchange-rates", h.Trade.CreateExchangeRate)
		protected.GET("/trade/exchange-rates/latest", h.Trade.GetLatestExchangeRate)

		// Trade Analytics
		protected.GET("/trade/analytics/statistics", h.Trade.GetTradeStatistics)
		protected.GET("/trade/analytics/shipments-by-country", h.Trade.GetShipmentsByCountry)
		protected.GET("/trade/analytics/top-trading-partners", h.Trade.GetTopTradingPartners)

		// Trade Utilities
		protected.POST("/trade/utils/calculate-tariff-duty", h.Trade.CalculateTariffDuty)
		protected.POST("/trade/utils/convert-currency", h.Trade.ConvertCurrency)

		// Advanced Feature routes
		// AI Assistant
		protected.GET("/advanced/ai-assistants", h.Advanced.ListAIAssistants)
		protected.POST("/advanced/ai-assistants", h.Advanced.CreateAIAssistant)
		protected.GET("/advanced/ai-assistants/:id", h.Advanced.GetAIAssistant)
		protected.PUT("/advanced/ai-assistants/:id", h.Advanced.UpdateAIAssistant)
		protected.DELETE("/advanced/ai-assistants/:id", h.Advanced.DeleteAIAssistant)

		// AI Conversations
		protected.POST("/advanced/conversations", h.Advanced.StartConversation)
		protected.POST("/advanced/conversations/:session_id/messages", h.Advanced.SendMessage)
		protected.GET("/advanced/conversations/:session_id/messages", h.Advanced.GetConversationHistory)
		protected.POST("/advanced/conversations/:session_id/end", h.Advanced.EndConversation)

		// Recommendations
		protected.GET("/advanced/recommendations", h.Advanced.ListRecommendations)
		protected.POST("/advanced/recommendations", h.Advanced.CreateRecommendation)
		protected.PUT("/advanced/recommendations/:id/status", h.Advanced.UpdateRecommendationStatus)

		// Advanced Search
		protected.GET("/advanced/searches", h.Advanced.ListAdvancedSearches)
		protected.POST("/advanced/searches", h.Advanced.CreateAdvancedSearch)
		protected.POST("/advanced/searches/:id/execute", h.Advanced.ExecuteAdvancedSearch)

		// Batch Operations
		protected.GET("/advanced/batch-operations", h.Advanced.ListBatchOperations)
		protected.POST("/advanced/batch-operations", h.Advanced.CreateBatchOperation)
		protected.GET("/advanced/batch-operations/:id", h.Advanced.GetBatchOperation)

		// Custom Fields
		protected.GET("/advanced/custom-fields", h.Advanced.ListCustomFields)
		protected.POST("/advanced/custom-fields", h.Advanced.CreateCustomField)
		protected.POST("/advanced/custom-field-values", h.Advanced.SetCustomFieldValue)
		protected.GET("/advanced/custom-field-values/:resource_id", h.Advanced.GetCustomFieldValues)

		// Security Events
		protected.GET("/advanced/security-events", h.Advanced.ListSecurityEvents)
		protected.POST("/advanced/security-events", h.Advanced.CreateSecurityEvent)

		// Performance Metrics
		protected.POST("/advanced/performance-metrics", h.Advanced.RecordPerformanceMetric)
		protected.GET("/advanced/performance-stats", h.Advanced.GetPerformanceStats)

		// Backups
		protected.GET("/advanced/backups", h.Advanced.ListBackups)
		protected.POST("/advanced/backups", h.Advanced.CreateBackup)

		// Multi-language
		protected.GET("/advanced/languages", h.Advanced.ListSystemLanguages)
		protected.GET("/advanced/translations/:language_code", h.Advanced.GetTranslations)

		// Integration routes
		// Integrations
		protected.GET("/integrations", h.Integration.ListIntegrations)
		protected.POST("/integrations", h.Integration.CreateIntegration)
		protected.GET("/integrations/:id", h.Integration.GetIntegration)
		protected.PUT("/integrations/:id", h.Integration.UpdateIntegration)
		protected.DELETE("/integrations/:id", h.Integration.DeleteIntegration)
		protected.POST("/integrations/:id/test", h.Integration.TestIntegration)

		// Integration Mappings
		protected.GET("/integrations/:integration_id/mappings", h.Integration.ListIntegrationMappings)
		protected.POST("/integrations/mappings", h.Integration.CreateIntegrationMapping)
		protected.PUT("/integrations/mappings/:id", h.Integration.UpdateIntegrationMapping)

		// Webhooks
		protected.GET("/integrations/webhooks", h.Integration.ListWebhooks)
		protected.POST("/integrations/webhooks", h.Integration.CreateWebhook)
		protected.PUT("/integrations/webhooks/:id", h.Integration.UpdateWebhook)
		protected.POST("/integrations/webhooks/:id/trigger", h.Integration.TriggerWebhook)
		protected.GET("/integrations/webhooks/:webhook_id/deliveries", h.Integration.GetWebhookDeliveries)

		// Data Sync Jobs
		protected.GET("/integrations/:integration_id/sync-jobs", h.Integration.ListDataSyncJobs)
		protected.POST("/integrations/sync-jobs", h.Integration.CreateDataSyncJob)
		protected.POST("/integrations/sync-jobs/:job_id/start", h.Integration.StartDataSyncJob)

		// API Keys
		protected.GET("/integrations/api-keys", h.Integration.ListApiKeys)
		protected.POST("/integrations/api-keys", h.Integration.CreateApiKey)
		protected.DELETE("/integrations/api-keys/:id", h.Integration.RevokeApiKey)

		// External Systems
		protected.GET("/integrations/external-systems", h.Integration.ListExternalSystems)
		protected.POST("/integrations/external-systems", h.Integration.CreateExternalSystem)
		protected.POST("/integrations/external-systems/:id/test", h.Integration.TestExternalSystem)

		// Integration Templates
		protected.GET("/integrations/templates", h.Integration.ListIntegrationTemplates)
		protected.POST("/integrations/templates/:template_id/create", h.Integration.CreateIntegrationFromTemplate)

		// Integration Analytics
		protected.GET("/integrations/analytics/stats", h.Integration.GetIntegrationStats)
		protected.GET("/integrations/analytics/by-type", h.Integration.GetIntegrationsByType)
		protected.GET("/integrations/analytics/sync-trends", h.Integration.GetSyncJobTrends)

		// Integration Utilities
		protected.POST("/integrations/utils/validate-mapping", h.Integration.ValidateMapping)
		protected.POST("/integrations/utils/preview-transformation", h.Integration.PreviewDataTransformation)
	}
}