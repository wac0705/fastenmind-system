package service

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/fastenmind/fastener-api/internal/models"
	"github.com/fastenmind/fastener-api/internal/repository"
	"github.com/google/uuid"
)

type QuoteService interface {
	List(companyID uuid.UUID, params map[string]interface{}) ([]models.Quote, int64, error)
	Get(id uuid.UUID) (*models.Quote, error)
	Create(companyID, userID uuid.UUID, req CreateQuoteRequest) (*models.Quote, error)
	Update(id, userID uuid.UUID, req UpdateQuoteRequest) (*models.Quote, error)
	Delete(id uuid.UUID) error
	SubmitForReview(id, userID uuid.UUID) (*models.Quote, error)
	Review(id, reviewerID uuid.UUID, req ReviewQuoteRequest) (*models.Quote, error)
	Send(id, userID uuid.UUID, req SendQuoteRequest) (*models.Quote, error)
	GetCostBreakdown(id uuid.UUID) (*CostBreakdown, error)
	GeneratePDF(id uuid.UUID) ([]byte, string, error)
	GetVersions(id uuid.UUID) ([]models.QuoteVersion, error)
	Duplicate(id, userID uuid.UUID) (*models.Quote, error)
}

type CreateQuoteRequest struct {
	InquiryID      uuid.UUID `json:"inquiry_id" validate:"required"`
	MaterialCost   float64   `json:"material_cost" validate:"min=0"`
	ProcessCost    float64   `json:"process_cost" validate:"min=0"`
	SurfaceCost    float64   `json:"surface_cost" validate:"min=0"`
	HeatTreatCost  float64   `json:"heat_treat_cost" validate:"min=0"`
	PackagingCost  float64   `json:"packaging_cost" validate:"min=0"`
	ShippingCost   float64   `json:"shipping_cost" validate:"min=0"`
	TariffCost     float64   `json:"tariff_cost" validate:"min=0"`
	OverheadRate   float64   `json:"overhead_rate" validate:"min=0,max=100"`
	ProfitRate     float64   `json:"profit_rate" validate:"min=0,max=100"`
	Currency       string    `json:"currency"`
	ValidUntil     string    `json:"valid_until" validate:"required"`
	DeliveryDays   int       `json:"delivery_days" validate:"min=1"`
	PaymentTerms   string    `json:"payment_terms" validate:"required"`
	Notes          string    `json:"notes"`
}

type UpdateQuoteRequest struct {
	MaterialCost   float64 `json:"material_cost"`
	ProcessCost    float64 `json:"process_cost"`
	SurfaceCost    float64 `json:"surface_cost"`
	HeatTreatCost  float64 `json:"heat_treat_cost"`
	PackagingCost  float64 `json:"packaging_cost"`
	ShippingCost   float64 `json:"shipping_cost"`
	TariffCost     float64 `json:"tariff_cost"`
	OverheadRate   float64 `json:"overhead_rate"`
	ProfitRate     float64 `json:"profit_rate"`
	DeliveryDays   int     `json:"delivery_days"`
	PaymentTerms   string  `json:"payment_terms"`
	Notes          string  `json:"notes"`
}

type ReviewQuoteRequest struct {
	Action   string `json:"action" validate:"required,oneof=approve reject"`
	Comments string `json:"comments"`
}

type SendQuoteRequest struct {
	Message string `json:"message"`
}

type CostBreakdown struct {
	MaterialCost    float64 `json:"material_cost"`
	ProcessCost     float64 `json:"process_cost"`
	SurfaceCost     float64 `json:"surface_cost"`
	HeatTreatCost   float64 `json:"heat_treat_cost"`
	PackagingCost   float64 `json:"packaging_cost"`
	ShippingCost    float64 `json:"shipping_cost"`
	TariffCost      float64 `json:"tariff_cost"`
	Subtotal        float64 `json:"subtotal"`
	OverheadAmount  float64 `json:"overhead_amount"`
	ProfitAmount    float64 `json:"profit_amount"`
	TotalCost       float64 `json:"total_cost"`
	ProcessDetails  []ProcessDetail `json:"process_details,omitempty"`
}

type ProcessDetail struct {
	Name string  `json:"name"`
	Cost float64 `json:"cost"`
}

type quoteService struct {
	quoteRepo      repository.QuoteRepository
	inquiryRepo    repository.InquiryRepository
	customerRepo   repository.CustomerRepository
	n8nService     N8NService
	pdfGenerator   PDFGenerator
}

func NewQuoteService(
	quoteRepo repository.QuoteRepository,
	inquiryRepo repository.InquiryRepository,
	customerRepo repository.CustomerRepository,
	n8nService N8NService,
	pdfGenerator PDFGenerator,
) QuoteService {
	return &quoteService{
		quoteRepo:    quoteRepo,
		inquiryRepo:  inquiryRepo,
		customerRepo: customerRepo,
		n8nService:   n8nService,
		pdfGenerator: pdfGenerator,
	}
}

func (s *quoteService) List(companyID uuid.UUID, params map[string]interface{}) ([]models.Quote, int64, error) {
	return s.quoteRepo.List(companyID, params)
}

func (s *quoteService) Get(id uuid.UUID) (*models.Quote, error) {
	return s.quoteRepo.GetWithDetails(id)
}

func (s *quoteService) Create(companyID, userID uuid.UUID, req CreateQuoteRequest) (*models.Quote, error) {
	// Get inquiry details
	inquiry, err := s.inquiryRepo.Get(req.InquiryID)
	if err != nil {
		return nil, errors.New("inquiry not found")
	}
	
	if inquiry.Status != "assigned" {
		return nil, errors.New("inquiry must be assigned before creating quote")
	}
	
	// Calculate costs
	subtotal := req.MaterialCost + req.ProcessCost + req.SurfaceCost + 
		req.HeatTreatCost + req.PackagingCost + req.ShippingCost + req.TariffCost
	
	overheadAmount := subtotal * (req.OverheadRate / 100)
	subtotalWithOverhead := subtotal + overheadAmount
	profitAmount := subtotalWithOverhead * (req.ProfitRate / 100)
	totalCost := subtotalWithOverhead + profitAmount
	unitPrice := totalCost / float64(inquiry.Quantity)
	
	// Parse ValidUntil date
	validUntil, err := time.Parse("2006-01-02", req.ValidUntil)
	if err != nil {
		return nil, errors.New("invalid date format for valid_until, expected YYYY-MM-DD")
	}
	
	// Create quote
	quote := &models.Quote{
		CompanyID:      companyID,
		InquiryID:      req.InquiryID,
		CustomerID:     inquiry.CustomerID,
		SalesID:        inquiry.SalesID,
		EngineerID:     userID,
		Status:         "draft",
		MaterialCost:   req.MaterialCost,
		ProcessCost:    req.ProcessCost,
		SurfaceCost:    req.SurfaceCost,
		HeatTreatCost:  req.HeatTreatCost,
		PackagingCost:  req.PackagingCost,
		ShippingCost:   req.ShippingCost,
		TariffCost:     req.TariffCost,
		OverheadRate:   req.OverheadRate,
		ProfitRate:     req.ProfitRate,
		TotalCost:      totalCost,
		UnitPrice:      unitPrice,
		Currency:       req.Currency,
		ValidUntil:     validUntil,
		DeliveryDays:   req.DeliveryDays,
		PaymentTerms:   req.PaymentTerms,
		Notes:          req.Notes,
	}
	
	if quote.Currency == "" {
		quote.Currency = "USD"
	}
	
	// Generate quote number
	quote.QuoteNo = s.generateQuoteNo(companyID)
	
	if err := s.quoteRepo.Create(quote); err != nil {
		return nil, err
	}
	
	// Create version 1
	s.createVersion(quote, userID, "Initial version")
	
	// Update inquiry status
	inquiry.Status = "quoted"
	inquiry.QuoteID = &quote.ID
	now := time.Now()
	inquiry.QuotedAt = &now
	s.inquiryRepo.Update(inquiry)
	
	// Log activity
	s.logActivity(quote.ID, userID, "created", "Quote created")
	
	// Trigger N8N workflow
	go s.n8nService.LogEvent(companyID, userID, "quote.created", "quote", quote.ID, map[string]interface{}{
		"quote_no":    quote.QuoteNo,
		"customer_id": quote.CustomerID,
		"total_cost":  quote.TotalCost,
	})
	
	return s.Get(quote.ID)
}

func (s *quoteService) Update(id, userID uuid.UUID, req UpdateQuoteRequest) (*models.Quote, error) {
	quote, err := s.quoteRepo.Get(id)
	if err != nil {
		return nil, err
	}
	
	if quote.Status != "draft" {
		return nil, errors.New("can only update draft quotes")
	}
	
	// Update fields
	quote.MaterialCost = req.MaterialCost
	quote.ProcessCost = req.ProcessCost
	quote.SurfaceCost = req.SurfaceCost
	quote.HeatTreatCost = req.HeatTreatCost
	quote.PackagingCost = req.PackagingCost
	quote.ShippingCost = req.ShippingCost
	quote.TariffCost = req.TariffCost
	quote.OverheadRate = req.OverheadRate
	quote.ProfitRate = req.ProfitRate
	quote.DeliveryDays = req.DeliveryDays
	quote.PaymentTerms = req.PaymentTerms
	quote.Notes = req.Notes
	
	// Recalculate
	subtotal := quote.MaterialCost + quote.ProcessCost + quote.SurfaceCost + 
		quote.HeatTreatCost + quote.PackagingCost + quote.ShippingCost + quote.TariffCost
	
	overheadAmount := subtotal * (quote.OverheadRate / 100)
	subtotalWithOverhead := subtotal + overheadAmount
	profitAmount := subtotalWithOverhead * (quote.ProfitRate / 100)
	quote.TotalCost = subtotalWithOverhead + profitAmount
	
	// Get inquiry for quantity
	inquiry, _ := s.inquiryRepo.Get(quote.InquiryID)
	if inquiry != nil {
		quote.UnitPrice = quote.TotalCost / float64(inquiry.Quantity)
	}
	
	if err := s.quoteRepo.Update(quote); err != nil {
		return nil, err
	}
	
	// Create new version
	s.createVersion(quote, userID, "Updated costs")
	
	// Log activity
	s.logActivity(quote.ID, userID, "updated", "Quote updated")
	
	return s.Get(quote.ID)
}

func (s *quoteService) Delete(id uuid.UUID) error {
	quote, err := s.quoteRepo.Get(id)
	if err != nil {
		return err
	}
	
	if quote.Status != "draft" {
		return errors.New("can only delete draft quotes")
	}
	
	return s.quoteRepo.Delete(id)
}

func (s *quoteService) SubmitForReview(id, userID uuid.UUID) (*models.Quote, error) {
	quote, err := s.quoteRepo.Get(id)
	if err != nil {
		return nil, err
	}
	
	if quote.Status != "draft" {
		return nil, errors.New("quote must be in draft status")
	}
	
	quote.Status = "pending_review"
	now := time.Now()
	quote.SubmittedAt = &now
	
	if err := s.quoteRepo.Update(quote); err != nil {
		return nil, err
	}
	
	// Log activity
	s.logActivity(quote.ID, userID, "submitted", "Quote submitted for review")
	
	// Trigger workflow
	go s.n8nService.LogEvent(quote.CompanyID, userID, "quote.submitted_for_review", "quote", quote.ID, map[string]interface{}{
		"quote_no": quote.QuoteNo,
	})
	
	return s.Get(quote.ID)
}

func (s *quoteService) Review(id, reviewerID uuid.UUID, req ReviewQuoteRequest) (*models.Quote, error) {
	quote, err := s.quoteRepo.Get(id)
	if err != nil {
		return nil, err
	}
	
	if quote.Status != "pending_review" && quote.Status != "under_review" {
		return nil, errors.New("quote is not pending review")
	}
	
	now := time.Now()
	quote.ReviewedAt = &now
	quote.ReviewerID = &reviewerID
	quote.ReviewComments = req.Comments
	
	if req.Action == "approve" {
		quote.Status = "approved"
	} else {
		quote.Status = "rejected"
	}
	
	if err := s.quoteRepo.Update(quote); err != nil {
		return nil, err
	}
	
	// Log activity
	action := fmt.Sprintf("%s quote", req.Action)
	s.logActivity(quote.ID, reviewerID, req.Action, action)
	
	// Trigger workflow
	eventType := fmt.Sprintf("quote.%s", req.Action)
	go s.n8nService.LogEvent(quote.CompanyID, reviewerID, eventType, "quote", quote.ID, map[string]interface{}{
		"quote_no": quote.QuoteNo,
		"comments": req.Comments,
	})
	
	return s.Get(quote.ID)
}

func (s *quoteService) Send(id, userID uuid.UUID, req SendQuoteRequest) (*models.Quote, error) {
	quote, err := s.quoteRepo.GetWithDetails(id)
	if err != nil {
		return nil, err
	}
	
	if quote.Status != "approved" {
		return nil, errors.New("quote must be approved before sending")
	}
	
	// Generate PDF
	_, _, err = s.GeneratePDF(id)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}
	
	// TODO: Send email with PDF attachment
	// For now, just update status
	
	quote.Status = "sent"
	now := time.Now()
	quote.SentAt = &now
	quote.SentByID = &userID
	
	if err := s.quoteRepo.Update(quote); err != nil {
		return nil, err
	}
	
	// Log activity
	s.logActivity(quote.ID, userID, "sent", "Quote sent to customer")
	
	// Trigger workflow
	go s.n8nService.LogEvent(quote.CompanyID, userID, "quote.sent", "quote", quote.ID, map[string]interface{}{
		"quote_no":      quote.QuoteNo,
		"customer_id":   quote.CustomerID,
		"email_message": req.Message,
	})
	
	return s.Get(quote.ID)
}

func (s *quoteService) GetCostBreakdown(id uuid.UUID) (*CostBreakdown, error) {
	quote, err := s.quoteRepo.Get(id)
	if err != nil {
		return nil, err
	}
	
	subtotal := quote.MaterialCost + quote.ProcessCost + quote.SurfaceCost + 
		quote.HeatTreatCost + quote.PackagingCost + quote.ShippingCost + quote.TariffCost
	
	overheadAmount := subtotal * (quote.OverheadRate / 100)
	subtotalWithOverhead := subtotal + overheadAmount
	profitAmount := subtotalWithOverhead * (quote.ProfitRate / 100)
	
	breakdown := &CostBreakdown{
		MaterialCost:   quote.MaterialCost,
		ProcessCost:    quote.ProcessCost,
		SurfaceCost:    quote.SurfaceCost,
		HeatTreatCost:  quote.HeatTreatCost,
		PackagingCost:  quote.PackagingCost,
		ShippingCost:   quote.ShippingCost,
		TariffCost:     quote.TariffCost,
		Subtotal:       subtotal,
		OverheadAmount: overheadAmount,
		ProfitAmount:   profitAmount,
		TotalCost:      quote.TotalCost,
	}
	
	// TODO: Add process details if available
	
	return breakdown, nil
}

func (s *quoteService) GeneratePDF(id uuid.UUID) ([]byte, string, error) {
	quote, err := s.quoteRepo.GetWithDetails(id)
	if err != nil {
		return nil, "", err
	}
	
	// Generate PDF using PDF generator
	pdfData, err := s.pdfGenerator.GenerateQuotePDF(quote)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate PDF: %w", err)
	}
	
	filename := fmt.Sprintf("Quote_%s.pdf", quote.QuoteNo)
	
	return pdfData, filename, nil
}

func (s *quoteService) GetVersions(id uuid.UUID) ([]models.QuoteVersion, error) {
	return s.quoteRepo.GetVersions(id)
}

func (s *quoteService) Duplicate(id, userID uuid.UUID) (*models.Quote, error) {
	originalQuote, err := s.quoteRepo.Get(id)
	if err != nil {
		return nil, err
	}
	
	// Create new quote with same values
	newQuote := &models.Quote{
		CompanyID:      originalQuote.CompanyID,
		InquiryID:      originalQuote.InquiryID,
		CustomerID:     originalQuote.CustomerID,
		SalesID:        originalQuote.SalesID,
		EngineerID:     userID,
		Status:         "draft",
		MaterialCost:   originalQuote.MaterialCost,
		ProcessCost:    originalQuote.ProcessCost,
		SurfaceCost:    originalQuote.SurfaceCost,
		HeatTreatCost:  originalQuote.HeatTreatCost,
		PackagingCost:  originalQuote.PackagingCost,
		ShippingCost:   originalQuote.ShippingCost,
		TariffCost:     originalQuote.TariffCost,
		OverheadRate:   originalQuote.OverheadRate,
		ProfitRate:     originalQuote.ProfitRate,
		TotalCost:      originalQuote.TotalCost,
		UnitPrice:      originalQuote.UnitPrice,
		Currency:       originalQuote.Currency,
		ValidUntil:     originalQuote.ValidUntil,
		DeliveryDays:   originalQuote.DeliveryDays,
		PaymentTerms:   originalQuote.PaymentTerms,
		Notes:          originalQuote.Notes,
	}
	
	// Generate new quote number
	newQuote.QuoteNo = s.generateQuoteNo(originalQuote.CompanyID)
	
	if err := s.quoteRepo.Create(newQuote); err != nil {
		return nil, err
	}
	
	// Create version 1
	s.createVersion(newQuote, userID, fmt.Sprintf("Duplicated from %s", originalQuote.QuoteNo))
	
	// Log activity
	s.logActivity(newQuote.ID, userID, "created", fmt.Sprintf("Quote duplicated from %s", originalQuote.QuoteNo))
	
	return s.Get(newQuote.ID)
}

func (s *quoteService) generateQuoteNo(companyID uuid.UUID) string {
	// TODO: Implement proper quote number generation
	return fmt.Sprintf("Q-%s", time.Now().Format("20060102-150405"))
}

func (s *quoteService) createVersion(quote *models.Quote, userID uuid.UUID, summary string) {
	// Get current version number
	versions, _ := s.quoteRepo.GetVersions(quote.ID)
	versionNumber := len(versions) + 1
	
	version := &models.QuoteVersion{
		QuoteID:        quote.ID,
		VersionNumber:  versionNumber,
		MaterialCost:   quote.MaterialCost,
		ProcessCost:    quote.ProcessCost,
		SurfaceCost:    quote.SurfaceCost,
		HeatTreatCost:  quote.HeatTreatCost,
		PackagingCost:  quote.PackagingCost,
		ShippingCost:   quote.ShippingCost,
		TariffCost:     quote.TariffCost,
		OverheadRate:   quote.OverheadRate,
		ProfitRate:     quote.ProfitRate,
		TotalCost:      quote.TotalCost,
		UnitPrice:      quote.UnitPrice,
		ChangeSummary:  summary,
		CreatedBy:      userID,
	}
	
	s.quoteRepo.CreateVersion(version)
}

func (s *quoteService) logActivity(quoteID, userID uuid.UUID, action, description string) {
	activity := &models.QuoteActivity{
		QuoteID:     quoteID,
		UserID:      userID,
		Action:      action,
		Description: description,
	}
	
	s.quoteRepo.LogActivity(activity)
}