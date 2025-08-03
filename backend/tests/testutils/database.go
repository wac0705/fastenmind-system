package testutils

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/fastenmind/fastener-api/internal/model"
	"github.com/fastenmind/fastener-api/pkg/database"
)

// TestDB provides utilities for test database operations
type TestDB struct {
	DB *gorm.DB
}

// NewTestDB creates a new test database connection
func NewTestDB(t *testing.T) *TestDB {
	config := &database.Config{
		Host:     getEnvOrDefault("TEST_DB_HOST", "localhost"),
		Port:     5432,
		User:     getEnvOrDefault("TEST_DB_USER", "postgres"),
		Password: getEnvOrDefault("TEST_DB_PASSWORD", "password"),
		DBName:   getEnvOrDefault("TEST_DB_NAME", "fastenmind_test"),
		SSLMode:  "disable",
	}

	db, err := database.Connect(config)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	return &TestDB{DB: db}
}

// SetupTables creates all necessary tables for testing
func (tdb *TestDB) SetupTables(t *testing.T) {
	err := tdb.DB.AutoMigrate(
		&model.Company{},
		&model.Account{},
		&model.Customer{},
		&model.Inquiry{},
		&model.Quote{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate test tables: %v", err)
	}
}

// CleanupTables removes all test data
func (tdb *TestDB) CleanupTables() {
	tables := []string{
		"quotes",
		"inquiries", 
		"customers",
		"accounts",
		"companies",
	}

	// 使用 GORM 的安全方法來清理資料
	for _, table := range tables {
		tdb.DB.Exec("TRUNCATE TABLE " + table + " RESTART IDENTITY CASCADE")
	}
}

// CreateTestCompany creates a test company
func (tdb *TestDB) CreateTestCompany() *model.Company {
	company := &model.Company{
		Code:    fmt.Sprintf("TEST%s", uuid.New().String()[:8]),
		Name:    "Test Company",
		Country: "US",
		Type:    "headquarters",
	}
	
	err := tdb.DB.Create(company).Error
	if err != nil {
		log.Fatalf("Failed to create test company: %v", err)
	}
	
	return company
}

// CreateTestAccount creates a test account
func (tdb *TestDB) CreateTestAccount(companyID uuid.UUID, role string) *model.Account {
	account := &model.Account{
		CompanyID:    companyID,
		Username:     fmt.Sprintf("test%s", uuid.New().String()[:8]),
		Email:        fmt.Sprintf("test%s@example.com", uuid.New().String()[:8]),
		PasswordHash: "$2a$10$XYZ...", // Pre-hashed password for "password123"
		FullName:     "Test User",
		Role:         role,
		IsActive:     true,
	}
	
	err := tdb.DB.Create(account).Error
	if err != nil {
		log.Fatalf("Failed to create test account: %v", err)
	}
	
	return account
}

// CreateTestCustomer creates a test customer
func (tdb *TestDB) CreateTestCustomer(companyID uuid.UUID) *model.Customer {
	customer := &model.Customer{
		CompanyID:    companyID,
		CustomerCode: fmt.Sprintf("CUST%s", uuid.New().String()[:8]),
		Name:         "Test Customer",
		Country:      "US",
		Currency:     "USD",
		IsActive:     true,
	}
	
	err := tdb.DB.Create(customer).Error
	if err != nil {
		log.Fatalf("Failed to create test customer: %v", err)
	}
	
	return customer
}

// CreateTestInquiry creates a test inquiry
func (tdb *TestDB) CreateTestInquiry(companyID, customerID, salesID uuid.UUID) *model.Inquiry {
	inquiry := &model.Inquiry{
		InquiryNo:       fmt.Sprintf("INQ%s", uuid.New().String()[:8]),
		CompanyID:       companyID,
		CustomerID:      customerID,
		SalesID:         salesID,
		Status:          "pending",
		ProductCategory: "Bolts",
		ProductName:     "Hex Bolt M8x20",
		Quantity:        10000,
		Unit:            "pcs",
		Incoterm:        "FOB",
		PaymentTerms:    "T/T 30 days",
	}
	
	err := tdb.DB.Create(inquiry).Error
	if err != nil {
		log.Fatalf("Failed to create test inquiry: %v", err)
	}
	
	return inquiry
}

// CreateTestQuote creates a test quote
func (tdb *TestDB) CreateTestQuote(inquiryID, companyID, customerID, engineerID uuid.UUID) *model.Quote {
	quote := &model.Quote{
		QuoteNo:       fmt.Sprintf("QUO%s", uuid.New().String()[:8]),
		InquiryID:     inquiryID,
		CompanyID:     companyID,
		CustomerID:    customerID,
		EngineerID:    engineerID,
		Status:        "draft",
		MaterialCost:  1000.00,
		ProcessCost:   500.00,
		TotalCost:     2000.00,
		UnitPrice:     0.20,
		Currency:      "USD",
		DeliveryDays:  30,
		PaymentTerms:  "T/T 30 days",
	}
	
	err := tdb.DB.Create(quote).Error
	if err != nil {
		log.Fatalf("Failed to create test quote: %v", err)
	}
	
	return quote
}

// Close closes the test database connection
func (tdb *TestDB) Close() {
	sqlDB, err := tdb.DB.DB()
	if err == nil {
		sqlDB.Close()
	}
}

// WithTransaction executes a function within a database transaction
func (tdb *TestDB) WithTransaction(fn func(*gorm.DB) error) error {
	return tdb.DB.Transaction(fn)
}

// getEnvOrDefault returns environment variable value or default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}