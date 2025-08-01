package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"github.com/fastenmind/fastener-api/internal/handler"
	"github.com/fastenmind/fastener-api/internal/model"
	"github.com/fastenmind/fastener-api/internal/repository"
	"github.com/fastenmind/fastener-api/internal/service"
	"github.com/fastenmind/fastener-api/pkg/database"
)

type QuoteIntegrationTestSuite struct {
	suite.Suite
	db           *gorm.DB
	echo         *echo.Echo
	quoteHandler *handler.QuoteHandler
	authHandler  *handler.AuthHandler
	cleanup      func()
	
	// Test data
	company  *model.Company
	customer *model.Customer
	sales    *model.Account
	engineer *model.Account
	inquiry  *model.Inquiry
}

func (suite *QuoteIntegrationTestSuite) SetupSuite() {
	// Setup test database connection
	config := &database.Config{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "password",
		DBName:   "fastenmind_test",
		SSLMode:  "disable",
	}
	
	db, err := database.Connect(config)
	suite.Require().NoError(err)
	suite.db = db
	
	// Auto-migrate tables
	err = db.AutoMigrate(
		&model.Company{},
		&model.Account{},
		&model.Customer{},
		&model.Inquiry{},
		&model.Quote{},
	)
	suite.Require().NoError(err)
	
	// Setup services and handlers
	accountRepo := repository.NewAccountRepository(db)
	authService := service.NewAuthService(accountRepo)
	suite.authHandler = handler.NewAuthHandler(authService)
	
	quoteRepo := repository.NewQuoteRepository(db)
	inquiryRepo := repository.NewInquiryRepository(db)
	quoteService := service.NewQuoteService(quoteRepo, inquiryRepo)
	suite.quoteHandler = handler.NewQuoteHandler(quoteService)
	
	// Setup Echo
	suite.echo = echo.New()
	
	// Setup cleanup function
	suite.cleanup = func() {
		// Clean up test data in correct order due to foreign key constraints
		db.Exec("DELETE FROM quotes")
		db.Exec("DELETE FROM inquiries")
		db.Exec("DELETE FROM customers")
		db.Exec("DELETE FROM accounts")
		db.Exec("DELETE FROM companies")
	}
}

func (suite *QuoteIntegrationTestSuite) TearDownSuite() {
	if suite.cleanup != nil {
		suite.cleanup()
	}
}

func (suite *QuoteIntegrationTestSuite) SetupTest() {
	// Clean up before each test
	suite.cleanup()
	
	// Create test data
	suite.setupTestData()
}

func (suite *QuoteIntegrationTestSuite) setupTestData() {
	// Create test company
	suite.company = &model.Company{
		Code:    "TEST001",
		Name:    "Test Company",
		Country: "US",
		Type:    "headquarters",
	}
	err := suite.db.Create(suite.company).Error
	suite.Require().NoError(err)
	
	// Create test customer
	suite.customer = &model.Customer{
		CompanyID:    suite.company.ID,
		CustomerCode: "CUST001",
		Name:         "Test Customer",
		Country:      "US",
		Currency:     "USD",
		IsActive:     true,
	}
	err = suite.db.Create(suite.customer).Error
	suite.Require().NoError(err)
	
	// Create test sales account
	suite.sales = &model.Account{
		CompanyID:    suite.company.ID,
		Username:     "sales1",
		Email:        "sales@test.com",
		PasswordHash: "$2a$10$XYZ...", // Pre-hashed password
		FullName:     "Sales User",
		Role:         "sales",
		IsActive:     true,
	}
	err = suite.db.Create(suite.sales).Error
	suite.Require().NoError(err)
	
	// Create test engineer account
	suite.engineer = &model.Account{
		CompanyID:    suite.company.ID,
		Username:     "engineer1",
		Email:        "engineer@test.com",
		PasswordHash: "$2a$10$XYZ...", // Pre-hashed password
		FullName:     "Engineer User",
		Role:         "engineer",
		IsActive:     true,
	}
	err = suite.db.Create(suite.engineer).Error
	suite.Require().NoError(err)
	
	// Create test inquiry
	suite.inquiry = &model.Inquiry{
		InquiryNo:       "INQ-2024-001",
		CompanyID:       suite.company.ID,
		CustomerID:      suite.customer.ID,
		SalesID:         suite.sales.ID,
		Status:          "assigned",
		ProductCategory: "Bolts",
		ProductName:     "Hex Bolt M8x20",
		Quantity:        10000,
		Unit:            "pcs",
		RequiredDate:    time.Now().AddDate(0, 1, 0),
		Incoterm:        "FOB",
		PaymentTerms:    "T/T 30 days",
		AssignedEngineerID: &suite.engineer.ID,
	}
	err = suite.db.Create(suite.inquiry).Error
	suite.Require().NoError(err)
}

func (suite *QuoteIntegrationTestSuite) TestCreateQuote_Success() {
	createRequest := handler.CreateQuoteRequest{
		InquiryID:     suite.inquiry.ID,
		MaterialCost:  1000.00,
		ProcessCost:   500.00,
		SurfaceCost:   200.00,
		PackagingCost: 100.00,
		ShippingCost:  300.00,
		OverheadRate:  0.15,
		ProfitRate:    0.20,
		Currency:      "USD",
		ValidUntil:    time.Now().AddDate(0, 1, 0),
		DeliveryDays:  30,
		PaymentTerms:  "T/T 30 days",
		Notes:         "Test quote",
	}
	
	requestBody, _ := json.Marshal(createRequest)
	req := httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewBuffer(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	
	// Add auth context
	req.Header.Set("Authorization", "Bearer mock-token")
	
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)
	
	// Mock JWT claims
	c.Set("user_id", suite.engineer.ID.String())
	c.Set("company_id", suite.company.ID.String())
	c.Set("role", "engineer")
	
	err := suite.quoteHandler.CreateQuote(c)
	suite.NoError(err)
	suite.Equal(http.StatusCreated, rec.Code)
	
	var response model.Quote
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal(suite.inquiry.ID, response.InquiryID)
	suite.Equal(suite.engineer.ID, response.EngineerID)
	suite.Equal("draft", response.Status)
	suite.True(response.TotalCost > 0)
	suite.True(response.UnitPrice > 0)
	
	// Verify quote was saved to database
	var savedQuote model.Quote
	err = suite.db.First(&savedQuote, response.ID).Error
	suite.NoError(err)
	suite.Equal(response.ID, savedQuote.ID)
}

func (suite *QuoteIntegrationTestSuite) TestGetQuote_Success() {
	// Create a test quote first
	quote := &model.Quote{
		QuoteNo:       "QUO-2024-001",
		InquiryID:     suite.inquiry.ID,
		CompanyID:     suite.company.ID,
		CustomerID:    suite.customer.ID,
		EngineerID:    suite.engineer.ID,
		Status:        "draft",
		MaterialCost:  1000.00,
		ProcessCost:   500.00,
		TotalCost:     2000.00,
		UnitPrice:     0.20,
		Currency:      "USD",
		ValidUntil:    time.Now().AddDate(0, 1, 0),
		DeliveryDays:  30,
		PaymentTerms:  "T/T 30 days",
	}
	err := suite.db.Create(quote).Error
	suite.Require().NoError(err)
	
	req := httptest.NewRequest(http.MethodGet, "/quotes/"+quote.ID.String(), nil)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(quote.ID.String())
	
	// Mock JWT claims
	c.Set("user_id", suite.engineer.ID.String())
	c.Set("company_id", suite.company.ID.String())
	c.Set("role", "engineer")
	
	err = suite.quoteHandler.GetQuote(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
	
	var response model.Quote
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal(quote.ID, response.ID)
	suite.Equal(quote.QuoteNo, response.QuoteNo)
	suite.Equal(quote.Status, response.Status)
}

func (suite *QuoteIntegrationTestSuite) TestListQuotes_Success() {
	// Create multiple test quotes
	quotes := []*model.Quote{
		{
			QuoteNo:      "QUO-2024-001",
			InquiryID:    suite.inquiry.ID,
			CompanyID:    suite.company.ID,
			CustomerID:   suite.customer.ID,
			EngineerID:   suite.engineer.ID,
			Status:       "draft",
			MaterialCost: 1000.00,
			TotalCost:    2000.00,
			Currency:     "USD",
		},
		{
			QuoteNo:      "QUO-2024-002",
			InquiryID:    suite.inquiry.ID,
			CompanyID:    suite.company.ID,
			CustomerID:   suite.customer.ID,
			EngineerID:   suite.engineer.ID,
			Status:       "approved",
			MaterialCost: 1500.00,
			TotalCost:    3000.00,
			Currency:     "USD",
		},
	}
	
	for _, quote := range quotes {
		err := suite.db.Create(quote).Error
		suite.Require().NoError(err)
	}
	
	req := httptest.NewRequest(http.MethodGet, "/quotes", nil)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)
	
	// Mock JWT claims
	c.Set("user_id", suite.engineer.ID.String())
	c.Set("company_id", suite.company.ID.String())
	c.Set("role", "engineer")
	
	err := suite.quoteHandler.ListQuotes(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
	
	var response []model.Quote
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Len(response, 2)
}

func (suite *QuoteIntegrationTestSuite) TestUpdateQuoteStatus_Success() {
	// Create a test quote first
	quote := &model.Quote{
		QuoteNo:    "QUO-2024-001",
		InquiryID:  suite.inquiry.ID,
		CompanyID:  suite.company.ID,
		CustomerID: suite.customer.ID,
		EngineerID: suite.engineer.ID,
		Status:     "draft",
		TotalCost:  2000.00,
		Currency:   "USD",
	}
	err := suite.db.Create(quote).Error
	suite.Require().NoError(err)
	
	updateRequest := handler.UpdateStatusRequest{
		Status: "pending_approval",
		Notes:  "Ready for review",
	}
	
	requestBody, _ := json.Marshal(updateRequest)
	req := httptest.NewRequest(http.MethodPut, "/quotes/"+quote.ID.String()+"/status", bytes.NewBuffer(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(quote.ID.String())
	
	// Mock JWT claims
	c.Set("user_id", suite.engineer.ID.String())
	c.Set("company_id", suite.company.ID.String())
	c.Set("role", "engineer")
	
	err = suite.quoteHandler.UpdateQuoteStatus(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
	
	// Verify status was updated in database
	var updatedQuote model.Quote
	err = suite.db.First(&updatedQuote, quote.ID).Error
	suite.NoError(err)
	suite.Equal("pending_approval", updatedQuote.Status)
}

func (suite *QuoteIntegrationTestSuite) TestDeleteQuote_Success() {
	// Create a test quote first
	quote := &model.Quote{
		QuoteNo:    "QUO-2024-001",
		InquiryID:  suite.inquiry.ID,
		CompanyID:  suite.company.ID,
		CustomerID: suite.customer.ID,
		EngineerID: suite.engineer.ID,
		Status:     "draft",
		TotalCost:  2000.00,
		Currency:   "USD",
	}
	err := suite.db.Create(quote).Error
	suite.Require().NoError(err)
	
	req := httptest.NewRequest(http.MethodDelete, "/quotes/"+quote.ID.String(), nil)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(quote.ID.String())
	
	// Mock JWT claims
	c.Set("user_id", suite.engineer.ID.String())
	c.Set("company_id", suite.company.ID.String())
	c.Set("role", "engineer")
	
	err = suite.quoteHandler.DeleteQuote(c)
	suite.NoError(err)
	suite.Equal(http.StatusNoContent, rec.Code)
	
	// Verify quote was soft deleted
	var deletedQuote model.Quote
	err = suite.db.First(&deletedQuote, quote.ID).Error
	suite.Error(err) // Should not find the quote
	suite.Equal(gorm.ErrRecordNotFound, err)
}

func (suite *QuoteIntegrationTestSuite) TestQuotePermissions() {
	// Create a test quote
	quote := &model.Quote{
		QuoteNo:    "QUO-2024-001",
		InquiryID:  suite.inquiry.ID,
		CompanyID:  suite.company.ID,
		CustomerID: suite.customer.ID,
		EngineerID: suite.engineer.ID,
		Status:     "draft",
		TotalCost:  2000.00,
		Currency:   "USD",
	}
	err := suite.db.Create(quote).Error
	suite.Require().NoError(err)
	
	// Test with different roles
	testCases := []struct {
		role           string
		shouldAccess   bool
		expectedStatus int
	}{
		{"admin", true, http.StatusOK},
		{"manager", true, http.StatusOK},
		{"engineer", true, http.StatusOK},
		{"sales", false, http.StatusForbidden},
		{"viewer", false, http.StatusForbidden},
	}
	
	for _, tc := range testCases {
		req := httptest.NewRequest(http.MethodGet, "/quotes/"+quote.ID.String(), nil)
		rec := httptest.NewRecorder()
		c := suite.echo.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(quote.ID.String())
		
		// Mock JWT claims with different roles
		c.Set("user_id", suite.engineer.ID.String())
		c.Set("company_id", suite.company.ID.String())
		c.Set("role", tc.role)
		
		err = suite.quoteHandler.GetQuote(c)
		
		if tc.shouldAccess {
			suite.NoError(err)
			suite.Equal(tc.expectedStatus, rec.Code)
		} else {
			suite.Error(err)
		}
	}
}

func TestQuoteIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(QuoteIntegrationTestSuite))
}