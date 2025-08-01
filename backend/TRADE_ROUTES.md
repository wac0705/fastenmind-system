# Trade Module Routes

To integrate the Trade module routes into the main server, add the following routes to the `setupRoutes` function in `cmd/server/main.go`:

```go
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

// Compliance
protected.POST("/trade/compliance/check", h.Trade.RunComplianceCheck)
protected.GET("/trade/compliance/checks", h.Trade.GetComplianceChecksByResource)
protected.GET("/trade/compliance/failed-checks", h.Trade.GetFailedComplianceChecks)

// Exchange rates
protected.GET("/trade/exchange-rates", h.Trade.ListExchangeRates)
protected.POST("/trade/exchange-rates", h.Trade.CreateExchangeRate)
protected.GET("/trade/exchange-rates/latest", h.Trade.GetLatestExchangeRate)

// Analytics
protected.GET("/trade/analytics/statistics", h.Trade.GetTradeStatistics)
protected.GET("/trade/analytics/shipments-by-country", h.Trade.GetShipmentsByCountry)
protected.GET("/trade/analytics/top-trading-partners", h.Trade.GetTopTradingPartners)

// Utilities
protected.POST("/trade/utils/calculate-tariff-duty", h.Trade.CalculateTariffDuty)
protected.POST("/trade/utils/convert-currency", h.Trade.ConvertCurrency)
```

## Handler Integration

Add the Trade handler to the Handlers struct in the handler package:

```go
type Handlers struct {
    // ... existing handlers
    Trade *TradeHandler
}

// In NewHandlers function:
Trade: NewTradeHandler(services.Trade),
```

## Service Integration

Add the Trade service to the Services struct in the service package:

```go
type Services struct {
    // ... existing services
    Trade *TradeService
}

// In NewServices function:
Trade: NewTradeService(repos.Trade, repos.User),
```

## Repository Integration

Add the Trade repository to the Repositories struct in the repository package:

```go
type Repositories struct {
    // ... existing repositories
    Trade *TradeRepository
}

// In NewRepositories function:
Trade: NewTradeRepository(db),
```