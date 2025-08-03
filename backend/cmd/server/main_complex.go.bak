package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fastenmind/fastener-api/internal/config"
	"github.com/fastenmind/fastener-api/internal/handler"
	"github.com/fastenmind/fastener-api/internal/middleware"
	"github.com/fastenmind/fastener-api/internal/repository"
	"github.com/fastenmind/fastener-api/internal/service"
	"github.com/fastenmind/fastener-api/pkg/concurrent"
	"github.com/fastenmind/fastener-api/pkg/database"
	"github.com/fastenmind/fastener-api/pkg/logger"
	"github.com/fastenmind/fastener-api/pkg/resources"
	"github.com/fastenmind/fastener-api/pkg/swagger"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	
	. "github.com/fastenmind/fastener-api/internal/handler" // Import CleanupHandler
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

	// Initialize global resource manager
	rm := resources.GetGlobalResourceManager()
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := resources.ShutdownGlobalResourceManager(shutdownCtx); err != nil {
			log.Error("Failed to shutdown resource manager:", err)
		}
	}()

	// Initialize service registry
	serviceRegistry := concurrent.NewServiceRegistry()

	// Initialize database wrapper
	dbWrapper, err := database.NewWrapper(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	
	// Register database with resource manager
	if err := rm.Databases().RegisterDB("main", dbWrapper.GormDB); err != nil {
		log.Fatal("Failed to register database:", err)
	}

	// Initialize Echo
	e := echo.New()
	e.HideBanner = true
	e.Logger = log

	// Setup error handling
	middleware.SetupErrorHandling(e)
	
	// Middleware
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.RequestID())
	e.Use(middleware.CORS(cfg.CORS))
	e.Use(middleware.Security())
	e.Use(middleware.ValidationMiddleware())
	e.Use(middleware.ResponseCleanup())
	e.Use(CleanupHandler())

	// Initialize repositories
	repos := repository.NewRepositories(dbWrapper.GormDB)

	// Initialize services
	services := service.NewServices(repos, cfg, dbWrapper.GormDB)

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
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	// Create shutdown context with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create a channel to track shutdown completion
	shutdownComplete := make(chan struct{})

	go func() {
		// Shutdown HTTP server
		if err := e.Shutdown(shutdownCtx); err != nil {
			log.Error("Failed to shutdown server:", err)
		}

		// Stop all services
		if err := serviceRegistry.StopAll(shutdownCtx); err != nil {
			log.Error("Failed to stop services:", err)
		}

		// Close database connections
		dbWrapper.Close()

		close(shutdownComplete)
	}()

	// Wait for shutdown to complete or timeout
	select {
	case <-shutdownComplete:
		log.Info("Server stopped gracefully")
	case <-shutdownCtx.Done():
		log.Error("Server shutdown timeout exceeded, forcing exit")
	}
}

func setupRoutes(e *echo.Echo, h *handler.Handlers, cfg *config.Config) {
	// Swagger documentation
	swaggerConfig := swagger.Config{
		Title:       "FastenMind API",
		Description: "FastenMind Fastener Manufacturing ERP System API",
		Version:     "1.0",
		Host:        "localhost:" + cfg.Server.Port,
		BasePath:    "/api/v1",
		Schemes:     []string{"http", "https"},
	}
	swagger.RegisterRoutes(e, swaggerConfig)

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
		protected.GET("/customers/:id/statistics", h.Customer.GetStatistics)
		protected.GET("/customers/:id/credit-history", h.Customer.GetCreditHistory)
		protected.GET("/customers/export", h.Customer.Export)

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
		
		// Engineer assignment routes
		protected.GET("/engineer-assignments/available", h.EngineerAssignment.GetAvailableEngineers)
		protected.POST("/engineer-assignments/assign", h.EngineerAssignment.AssignEngineer)
		protected.PUT("/engineer-assignments/:id/reassign", h.EngineerAssignment.ReassignEngineer)
		protected.GET("/engineer-assignments/history", h.EngineerAssignment.GetAssignmentHistory)
		protected.GET("/engineer-assignments/workload", h.EngineerAssignment.GetEngineerWorkload)
		protected.PUT("/engineer-assignments/:id/status", h.EngineerAssignment.UpdateAssignmentStatus)
		protected.GET("/engineer-assignments/stats", h.EngineerAssignment.GetAssignmentStats)
		protected.POST("/engineer-assignments/auto-assign", h.EngineerAssignment.AutoAssignEngineer)
		
		// Process cost routes
		protected.GET("/process-costs/templates", h.ProcessCost.GetCostTemplates)
		protected.POST("/process-costs/templates", h.ProcessCost.CreateCostTemplate)
		protected.PUT("/process-costs/templates/:id", h.ProcessCost.UpdateCostTemplate)
		protected.DELETE("/process-costs/templates/:id", h.ProcessCost.DeleteCostTemplate)
		protected.POST("/process-costs/calculate", h.ProcessCost.CalculateProcessCost)
		protected.GET("/process-costs/history", h.ProcessCost.GetCostHistory)
		protected.GET("/process-costs/materials", h.ProcessCost.GetMaterialCosts)
		protected.PUT("/process-costs/materials/:id", h.ProcessCost.UpdateMaterialCost)
		protected.GET("/process-costs/processing-rates", h.ProcessCost.GetProcessingRates)
		protected.PUT("/process-costs/processing-rates/:id", h.ProcessCost.UpdateProcessingRate)
		protected.POST("/process-costs/batch-calculate", h.ProcessCost.BatchCalculateCost)
		protected.GET("/process-costs/analysis", h.ProcessCost.GetCostAnalysis)
		protected.GET("/process-costs/export", h.ProcessCost.ExportCostReport)
		protected.GET("/process-costs/settings", h.ProcessCost.GetCostSettings)
		protected.PUT("/process-costs/settings", h.ProcessCost.UpdateCostSettings)

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

		// Report routes
		// Reports
		protected.GET("/reports", h.Report.ListReports)
		protected.POST("/reports", h.Report.CreateReport)
		protected.GET("/reports/:id", h.Report.GetReport)
		protected.PUT("/reports/:id", h.Report.UpdateReport)
		protected.DELETE("/reports/:id", h.Report.DeleteReport)
		protected.POST("/reports/:id/duplicate", h.Report.DuplicateReport)
		protected.POST("/reports/:id/execute", h.Report.ExecuteReport)
		protected.GET("/reports/:id/export", h.Report.ExportReport)
		
		// Report Templates
		protected.GET("/reports/templates", h.Report.ListReportTemplates)
		protected.POST("/reports/templates", h.Report.CreateReportTemplate)
		protected.GET("/reports/templates/:id", h.Report.GetReportTemplate)
		protected.PUT("/reports/templates/:id", h.Report.UpdateReportTemplate)
		protected.DELETE("/reports/templates/:id", h.Report.DeleteReportTemplate)
		protected.POST("/reports/templates/:template_id/generate", h.Report.GenerateReportFromTemplate)
		
		// Report Executions
		protected.GET("/reports/:report_id/executions", h.Report.ListReportExecutions)
		protected.GET("/reports/executions/:id", h.Report.GetReportExecution)
		protected.POST("/reports/executions/:id/cancel", h.Report.CancelReportExecution)
		protected.GET("/reports/executions/:id/download", h.Report.DownloadReportResult)
		
		// Report Subscriptions
		protected.GET("/reports/subscriptions", h.Report.ListReportSubscriptions)
		protected.POST("/reports/subscriptions", h.Report.CreateReportSubscription)
		protected.GET("/reports/subscriptions/:id", h.Report.GetReportSubscription)
		protected.PUT("/reports/subscriptions/:id", h.Report.UpdateReportSubscription)
		protected.DELETE("/reports/subscriptions/:id", h.Report.DeleteReportSubscription)
		
		// Report Business Operations
		protected.GET("/reports/dashboard", h.Report.GetReportDashboard)
		protected.GET("/reports/statistics", h.Report.GetReportStatistics)
		protected.GET("/reports/popular", h.Report.GetPopularReports)
		protected.GET("/reports/recent-executions", h.Report.GetRecentExecutions)
		protected.POST("/reports/import", h.Report.ImportReports)
	}
}