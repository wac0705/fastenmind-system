package handler

// @title FastenMind API
// @version 1.0
// @description FastenMind Fastener Manufacturing ERP System API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.fastenmind.com/support
// @contact.email support@fastenmind.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// Auth API Documentation

// @Summary User Login
// @Description Authenticate user and receive JWT tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} swagger.ErrorResponse
// @Failure 401 {object} swagger.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) LoginDocs() {}

// @Summary User Registration
// @Description Register a new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration information"
// @Success 201 {object} RegisterResponse
// @Failure 400 {object} swagger.ErrorResponse
// @Failure 409 {object} swagger.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) RegisterDocs() {}

// @Summary Refresh Token
// @Description Refresh access token using refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RefreshRequest true "Refresh token"
// @Success 200 {object} RefreshResponse
// @Failure 400 {object} swagger.ErrorResponse
// @Failure 401 {object} swagger.ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshTokenDocs() {}

// @Summary User Logout
// @Description Logout user and invalidate tokens
// @Tags Auth
// @Security Bearer
// @Accept json
// @Produce json
// @Success 200 {object} swagger.Response
// @Failure 401 {object} swagger.ErrorResponse
// @Router /auth/logout [post]
func (h *AuthHandler) LogoutDocs() {}

// Customer API Documentation

// @Summary List Customers
// @Description Get paginated list of customers with optional filters
// @Tags Customers
// @Security Bearer
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param search query string false "Search by name, code, or contact"
// @Param country query string false "Filter by country code"
// @Param currency query string false "Filter by currency"
// @Param status query string false "Filter by status"
// @Success 200 {object} swagger.PaginationResponse{data=[]CustomerResponse}
// @Failure 401 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Router /customers [get]
func (h *CustomerHandler) ListDocs() {}

// @Summary Get Customer
// @Description Get customer details by ID
// @Tags Customers
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Success 200 {object} swagger.Response{data=CustomerResponse}
// @Failure 401 {object} swagger.ErrorResponse
// @Failure 404 {object} swagger.ErrorResponse
// @Router /customers/{id} [get]
func (h *CustomerHandler) GetDocs() {}

// @Summary Create Customer
// @Description Create a new customer
// @Tags Customers
// @Security Bearer
// @Accept json
// @Produce json
// @Param request body CreateCustomerRequest true "Customer information"
// @Success 201 {object} swagger.Response{data=CustomerResponse}
// @Failure 400 {object} swagger.ErrorResponse
// @Failure 401 {object} swagger.ErrorResponse
// @Failure 409 {object} swagger.ErrorResponse
// @Router /customers [post]
func (h *CustomerHandler) CreateDocs() {}

// @Summary Update Customer
// @Description Update customer information
// @Tags Customers
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Param request body UpdateCustomerRequest true "Customer update information"
// @Success 200 {object} swagger.Response{data=CustomerResponse}
// @Failure 400 {object} swagger.ErrorResponse
// @Failure 401 {object} swagger.ErrorResponse
// @Failure 404 {object} swagger.ErrorResponse
// @Router /customers/{id} [put]
func (h *CustomerHandler) UpdateDocs() {}

// @Summary Delete Customer
// @Description Delete a customer (soft delete)
// @Tags Customers
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Success 200 {object} swagger.Response
// @Failure 401 {object} swagger.ErrorResponse
// @Failure 404 {object} swagger.ErrorResponse
// @Router /customers/{id} [delete]
func (h *CustomerHandler) DeleteDocs() {}

// Supplier API Documentation

// @Summary List Suppliers
// @Description Get paginated list of suppliers with optional filters
// @Tags Suppliers
// @Security Bearer
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param search query string false "Search by name or code"
// @Param type query string false "Filter by supplier type"
// @Param status query string false "Filter by status"
// @Success 200 {object} swagger.PaginationResponse{data=[]SupplierResponse}
// @Failure 401 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Router /suppliers [get]
func (h *SupplierHandler) ListDocs() {}

// Quote API Documentation

// @Summary Create Quote
// @Description Create a new quote request
// @Tags Quotes
// @Security Bearer
// @Accept json
// @Produce json
// @Param request body CreateQuoteRequest true "Quote information"
// @Success 201 {object} swagger.Response{data=QuoteResponse}
// @Failure 400 {object} swagger.ErrorResponse
// @Failure 401 {object} swagger.ErrorResponse
// @Router /quotes [post]
func (h *QuoteHandler) CreateDocs() {}

// @Summary Calculate Quote
// @Description Calculate quote pricing
// @Tags Quotes
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path string true "Quote ID"
// @Success 200 {object} swagger.Response{data=QuoteCalculationResponse}
// @Failure 400 {object} swagger.ErrorResponse
// @Failure 401 {object} swagger.ErrorResponse
// @Failure 404 {object} swagger.ErrorResponse
// @Router /quotes/{id}/calculate [post]
func (h *QuoteHandler) CalculateDocs() {}

// @Summary Submit Quote
// @Description Submit quote for approval
// @Tags Quotes
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path string true "Quote ID"
// @Success 200 {object} swagger.Response{data=QuoteResponse}
// @Failure 400 {object} swagger.ErrorResponse
// @Failure 401 {object} swagger.ErrorResponse
// @Failure 404 {object} swagger.ErrorResponse
// @Router /quotes/{id}/submit [post]
func (h *QuoteHandler) SubmitDocs() {}

// @Summary Approve Quote
// @Description Approve a submitted quote
// @Tags Quotes
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path string true "Quote ID"
// @Param request body ApprovalRequest true "Approval information"
// @Success 200 {object} swagger.Response{data=QuoteResponse}
// @Failure 400 {object} swagger.ErrorResponse
// @Failure 401 {object} swagger.ErrorResponse
// @Failure 403 {object} swagger.ErrorResponse
// @Failure 404 {object} swagger.ErrorResponse
// @Router /quotes/{id}/approve [post]
func (h *QuoteHandler) ApproveDocs() {}

// Exchange Rate API Documentation

// @Summary List Exchange Rates
// @Description Get current exchange rates
// @Tags ExchangeRates
// @Security Bearer
// @Accept json
// @Produce json
// @Param base_currency query string false "Base currency" default(USD)
// @Success 200 {object} swagger.Response{data=[]ExchangeRateResponse}
// @Failure 401 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Router /exchange-rates [get]
func (h *ExchangeRateHandler) ListDocs() {}

// @Summary Update Exchange Rate
// @Description Update exchange rate for a currency pair
// @Tags ExchangeRates
// @Security Bearer
// @Accept json
// @Produce json
// @Param request body UpdateExchangeRateRequest true "Exchange rate information"
// @Success 200 {object} swagger.Response{data=ExchangeRateResponse}
// @Failure 400 {object} swagger.ErrorResponse
// @Failure 401 {object} swagger.ErrorResponse
// @Failure 403 {object} swagger.ErrorResponse
// @Router /exchange-rates [put]
func (h *ExchangeRateHandler) UpdateDocs() {}

// Report API Documentation

// @Summary Generate Sales Report
// @Description Generate sales report for specified period
// @Tags Reports
// @Security Bearer
// @Accept json
// @Produce json
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Param customer_id query string false "Filter by customer ID"
// @Param format query string false "Output format (json/excel/pdf)" default(json)
// @Success 200 {object} swagger.Response{data=SalesReportResponse}
// @Failure 400 {object} swagger.ErrorResponse
// @Failure 401 {object} swagger.ErrorResponse
// @Router /reports/sales [get]
func (h *ReportHandler) SalesReportDocs() {}

// @Summary Generate Cost Analysis Report
// @Description Generate cost analysis report
// @Tags Reports
// @Security Bearer
// @Accept json
// @Produce json
// @Param product_id query string false "Filter by product ID"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} swagger.Response{data=CostAnalysisResponse}
// @Failure 400 {object} swagger.ErrorResponse
// @Failure 401 {object} swagger.ErrorResponse
// @Router /reports/cost-analysis [get]
func (h *ReportHandler) CostAnalysisReportDocs() {}