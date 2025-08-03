package aggregate

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/fastenmind/fastener-api/internal/domain/entity"
	"github.com/fastenmind/fastener-api/internal/domain/valueobject"
)

func TestNewQuoteAggregate(t *testing.T) {
	customerID := uuid.New()
	companyID := uuid.New()

	// Test successful creation
	t.Run("Successful creation", func(t *testing.T) {
		quote, err := NewQuoteAggregate(customerID, companyID)
		
		assert.NoError(t, err)
		assert.NotNil(t, quote)
		assert.NotEqual(t, uuid.Nil, quote.ID)
		assert.Equal(t, customerID, quote.CustomerID)
		assert.Equal(t, companyID, quote.CompanyID)
		assert.Equal(t, valueobject.QuoteStatusDraft, quote.Status)
		assert.NotEmpty(t, quote.QuoteNumber)
		assert.Equal(t, 1, quote.Version)
		assert.Len(t, quote.GetDomainEvents(), 1)
	})

	// Test with nil customer ID
	t.Run("Nil customer ID", func(t *testing.T) {
		quote, err := NewQuoteAggregate(uuid.Nil, companyID)
		
		assert.Error(t, err)
		assert.Nil(t, quote)
		assert.Equal(t, "customer ID is required", err.Error())
	})

	// Test with nil company ID
	t.Run("Nil company ID", func(t *testing.T) {
		quote, err := NewQuoteAggregate(customerID, uuid.Nil)
		
		assert.Error(t, err)
		assert.Nil(t, quote)
		assert.Equal(t, "company ID is required", err.Error())
	})
}

func TestQuoteAggregate_AddItem(t *testing.T) {
	quote := createTestQuote(t)
	
	// Test successful addition
	t.Run("Successful addition", func(t *testing.T) {
		item := createTestQuoteItem(t)
		
		err := quote.AddItem(item)
		
		assert.NoError(t, err)
		assert.Len(t, quote.Items, 1)
		assert.Equal(t, item.ID, quote.Items[0].ID)
		assert.Greater(t, quote.PricingSummary.Total, 0.0)
		assert.Len(t, quote.GetDomainEvents(), 2) // Create + AddItem
	})

	// Test adding duplicate item
	t.Run("Duplicate item", func(t *testing.T) {
		item := createTestQuoteItem(t)
		quote.AddItem(item)
		
		// Try to add same product with same specification
		duplicateItem := item
		duplicateItem.ID = uuid.New()
		
		err := quote.AddItem(duplicateItem)
		
		assert.Error(t, err)
		assert.Equal(t, "duplicate quote item", err.Error())
	})

	// Test adding to non-draft quote
	t.Run("Non-draft quote", func(t *testing.T) {
		submittedQuote := createTestQuote(t)
		submittedQuote.Status = valueobject.QuoteStatusPending
		
		item := createTestQuoteItem(t)
		err := submittedQuote.AddItem(item)
		
		assert.Error(t, err)
		assert.Equal(t, "only draft quotes can be modified", err.Error())
	})
}

func TestQuoteAggregate_Submit(t *testing.T) {
	// Test successful submission
	t.Run("Successful submission", func(t *testing.T) {
		quote := createTestQuote(t)
		item := createTestQuoteItem(t)
		quote.AddItem(item)
		
		err := quote.Submit()
		
		assert.NoError(t, err)
		assert.Equal(t, valueobject.QuoteStatusPending, quote.Status)
		assert.Equal(t, 2, quote.Version)
	})

	// Test submission without items
	t.Run("No items", func(t *testing.T) {
		quote := createTestQuote(t)
		
		err := quote.Submit()
		
		assert.Error(t, err)
		assert.Equal(t, "quote must have at least one item", err.Error())
	})

	// Test submission of non-draft quote
	t.Run("Non-draft quote", func(t *testing.T) {
		quote := createTestQuote(t)
		quote.Status = valueobject.QuoteStatusPending
		
		err := quote.Submit()
		
		assert.Error(t, err)
		assert.Equal(t, "only draft quotes can be submitted", err.Error())
	})
}

func TestQuoteAggregate_Approve(t *testing.T) {
	approverID := uuid.New()

	// Test successful approval
	t.Run("Successful approval", func(t *testing.T) {
		quote := createSubmittedQuote(t)
		
		err := quote.Approve(approverID)
		
		assert.NoError(t, err)
		assert.Equal(t, valueobject.QuoteStatusApproved, quote.Status)
		assert.Equal(t, 3, quote.Version)
	})

	// Test approval of non-pending quote
	t.Run("Non-pending quote", func(t *testing.T) {
		quote := createTestQuote(t)
		
		err := quote.Approve(approverID)
		
		assert.Error(t, err)
		assert.Equal(t, "only pending quotes can be approved", err.Error())
	})

	// Test approval of expired quote
	t.Run("Expired quote", func(t *testing.T) {
		quote := createSubmittedQuote(t)
		quote.ValidUntil = time.Now().Add(-24 * time.Hour)
		
		err := quote.Approve(approverID)
		
		assert.Error(t, err)
		assert.Equal(t, "quote has expired", err.Error())
	})
}

func TestQuoteAggregate_ExtendValidity(t *testing.T) {
	// Test successful extension
	t.Run("Successful extension", func(t *testing.T) {
		quote := createTestQuote(t)
		newValidUntil := time.Now().AddDate(0, 2, 0)
		
		err := quote.ExtendValidity(newValidUntil)
		
		assert.NoError(t, err)
		assert.Equal(t, newValidUntil.Truncate(time.Second), quote.ValidUntil.Truncate(time.Second))
	})

	// Test extension with past date
	t.Run("Past date", func(t *testing.T) {
		quote := createTestQuote(t)
		pastDate := time.Now().Add(-24 * time.Hour)
		
		err := quote.ExtendValidity(pastDate)
		
		assert.Error(t, err)
		assert.Equal(t, "new validity date must be in the future", err.Error())
	})

	// Test extension of expired quote
	t.Run("Expired quote", func(t *testing.T) {
		quote := createTestQuote(t)
		quote.Status = valueobject.QuoteStatusExpired
		
		newValidUntil := time.Now().AddDate(0, 2, 0)
		err := quote.ExtendValidity(newValidUntil)
		
		assert.Error(t, err)
		assert.Equal(t, "expired quotes cannot be extended", err.Error())
	})
}

func TestQuoteAggregate_Clone(t *testing.T) {
	// Test successful clone
	t.Run("Successful clone", func(t *testing.T) {
		originalQuote := createTestQuote(t)
		item := createTestQuoteItem(t)
		originalQuote.AddItem(item)
		
		clonedQuote, err := originalQuote.Clone()
		
		assert.NoError(t, err)
		assert.NotNil(t, clonedQuote)
		assert.NotEqual(t, originalQuote.ID, clonedQuote.ID)
		assert.NotEqual(t, originalQuote.QuoteNumber, clonedQuote.QuoteNumber)
		assert.Equal(t, originalQuote.CustomerID, clonedQuote.CustomerID)
		assert.Equal(t, originalQuote.CompanyID, clonedQuote.CompanyID)
		assert.Len(t, clonedQuote.Items, 1)
		assert.NotEqual(t, originalQuote.Items[0].ID, clonedQuote.Items[0].ID)
		assert.Equal(t, valueobject.QuoteStatusDraft, clonedQuote.Status)
	})
}

// Helper functions

func createTestQuote(t *testing.T) *QuoteAggregate {
	customerID := uuid.New()
	companyID := uuid.New()
	
	quote, err := NewQuoteAggregate(customerID, companyID)
	assert.NoError(t, err)
	
	// Clear initial events for easier testing
	quote.ClearDomainEvents()
	
	return quote
}

func createTestQuoteItem(t *testing.T) entity.QuoteItem {
	productID := uuid.New()
	
	item, err := entity.NewQuoteItem(productID, "Test Product", 10, 100.00)
	assert.NoError(t, err)
	
	item.Material = entity.Material{
		Type:  entity.MaterialTypeSteel,
		Grade: "Grade 8",
	}
	
	return *item
}

func createSubmittedQuote(t *testing.T) *QuoteAggregate {
	quote := createTestQuote(t)
	item := createTestQuoteItem(t)
	
	quote.AddItem(item)
	quote.Submit()
	quote.ClearDomainEvents()
	
	return quote
}