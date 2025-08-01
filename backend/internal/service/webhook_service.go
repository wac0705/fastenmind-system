package service

import (
	"fmt"
	"github.com/google/uuid"
)

// WebhookEvent represents different event types that can trigger N8N workflows
type WebhookEvent struct {
	EventType   string                 `json:"event_type"`
	EntityType  string                 `json:"entity_type"`
	EntityID    uuid.UUID              `json:"entity_id"`
	CompanyID   uuid.UUID              `json:"company_id"`
	UserID      uuid.UUID              `json:"user_id"`
	Timestamp   string                 `json:"timestamp"`
	Data        map[string]interface{} `json:"data"`
}

// WebhookService handles event triggers for N8N workflows
type WebhookService struct {
	n8nService N8NService
}

func NewWebhookService(n8nService N8NService) *WebhookService {
	return &WebhookService{
		n8nService: n8nService,
	}
}

// TriggerInquiryCreated triggers N8N workflow when an inquiry is created
func (s *WebhookService) TriggerInquiryCreated(inquiry interface{}, companyID, userID uuid.UUID) error {
	return s.n8nService.LogEvent(
		companyID,
		userID,
		"inquiry.created",
		"inquiry",
		inquiry.(interface{ GetID() uuid.UUID }).GetID(),
		map[string]interface{}{
			"inquiry": inquiry,
			"action":  "created",
		},
	)
}

// TriggerInquiryAssigned triggers N8N workflow when an inquiry is assigned
func (s *WebhookService) TriggerInquiryAssigned(inquiryID, engineerID, companyID, userID uuid.UUID) error {
	return s.n8nService.LogEvent(
		companyID,
		userID,
		"inquiry.assigned",
		"inquiry",
		inquiryID,
		map[string]interface{}{
			"inquiry_id":  inquiryID,
			"engineer_id": engineerID,
			"action":      "assigned",
		},
	)
}

// TriggerQuoteCreated triggers N8N workflow when a quote is created
func (s *WebhookService) TriggerQuoteCreated(quote interface{}, companyID, userID uuid.UUID) error {
	return s.n8nService.LogEvent(
		companyID,
		userID,
		"quote.created",
		"quote",
		quote.(interface{ GetID() uuid.UUID }).GetID(),
		map[string]interface{}{
			"quote":  quote,
			"action": "created",
		},
	)
}

// TriggerQuoteSubmittedForApproval triggers N8N workflow when a quote is submitted for approval
func (s *WebhookService) TriggerQuoteSubmittedForApproval(quoteID uuid.UUID, companyID, userID uuid.UUID) error {
	return s.n8nService.LogEvent(
		companyID,
		userID,
		"quote.submitted_for_approval",
		"quote",
		quoteID,
		map[string]interface{}{
			"quote_id": quoteID,
			"action":   "submitted_for_approval",
		},
	)
}

// TriggerQuoteApproved triggers N8N workflow when a quote is approved
func (s *WebhookService) TriggerQuoteApproved(quoteID uuid.UUID, approverID uuid.UUID, companyID, userID uuid.UUID) error {
	return s.n8nService.LogEvent(
		companyID,
		userID,
		"quote.approved",
		"quote",
		quoteID,
		map[string]interface{}{
			"quote_id":    quoteID,
			"approver_id": approverID,
			"action":      "approved",
		},
	)
}

// TriggerQuoteRejected triggers N8N workflow when a quote is rejected
func (s *WebhookService) TriggerQuoteRejected(quoteID uuid.UUID, reason string, companyID, userID uuid.UUID) error {
	return s.n8nService.LogEvent(
		companyID,
		userID,
		"quote.rejected",
		"quote",
		quoteID,
		map[string]interface{}{
			"quote_id": quoteID,
			"reason":   reason,
			"action":   "rejected",
		},
	)
}

// TriggerQuoteSent triggers N8N workflow when a quote is sent
func (s *WebhookService) TriggerQuoteSent(quoteID uuid.UUID, recipientEmail string, companyID, userID uuid.UUID) error {
	return s.n8nService.LogEvent(
		companyID,
		userID,
		"quote.sent",
		"quote",
		quoteID,
		map[string]interface{}{
			"quote_id":        quoteID,
			"recipient_email": recipientEmail,
			"action":          "sent",
		},
	)
}

// TriggerCostUpdateRequired triggers N8N workflow when cost update is required
func (s *WebhookService) TriggerCostUpdateRequired(entityType string, entityID uuid.UUID, reason string, companyID, userID uuid.UUID) error {
	return s.n8nService.LogEvent(
		companyID,
		userID,
		"cost.update_required",
		entityType,
		entityID,
		map[string]interface{}{
			"entity_type": entityType,
			"entity_id":   entityID,
			"reason":      reason,
			"action":      "update_required",
		},
	)
}

// TriggerExchangeRateUpdated triggers N8N workflow when exchange rate is updated
func (s *WebhookService) TriggerExchangeRateUpdated(currency string, oldRate, newRate float64, companyID, userID uuid.UUID) error {
	changePercent := ((newRate - oldRate) / oldRate) * 100
	
	return s.n8nService.LogEvent(
		companyID,
		userID,
		"exchange_rate.updated",
		"exchange_rate",
		uuid.New(), // No specific entity ID for exchange rates
		map[string]interface{}{
			"currency":       currency,
			"old_rate":       oldRate,
			"new_rate":       newRate,
			"change_percent": fmt.Sprintf("%.2f", changePercent),
			"action":         "updated",
		},
	)
}

// TriggerCustomerCreditLimitExceeded triggers N8N workflow when customer exceeds credit limit
func (s *WebhookService) TriggerCustomerCreditLimitExceeded(customerID uuid.UUID, currentCredit, creditLimit float64, companyID, userID uuid.UUID) error {
	return s.n8nService.LogEvent(
		companyID,
		userID,
		"customer.credit_limit_exceeded",
		"customer",
		customerID,
		map[string]interface{}{
			"customer_id":    customerID,
			"current_credit": currentCredit,
			"credit_limit":   creditLimit,
			"exceeded_by":    currentCredit - creditLimit,
			"action":         "credit_limit_exceeded",
		},
	)
}

// TriggerEngineerWorkloadHigh triggers N8N workflow when engineer workload is high
func (s *WebhookService) TriggerEngineerWorkloadHigh(engineerID uuid.UUID, currentWorkload, maxWorkload int, companyID, userID uuid.UUID) error {
	return s.n8nService.LogEvent(
		companyID,
		userID,
		"engineer.workload_high",
		"engineer",
		engineerID,
		map[string]interface{}{
			"engineer_id":      engineerID,
			"current_workload": currentWorkload,
			"max_workload":     maxWorkload,
			"utilization":      fmt.Sprintf("%.1f%%", (float64(currentWorkload)/float64(maxWorkload))*100),
			"action":           "workload_high",
		},
	)
}

// TriggerEquipmentMaintenanceDue triggers N8N workflow when equipment maintenance is due
func (s *WebhookService) TriggerEquipmentMaintenanceDue(equipmentID uuid.UUID, dueDate string, companyID, userID uuid.UUID) error {
	return s.n8nService.LogEvent(
		companyID,
		userID,
		"equipment.maintenance_due",
		"equipment",
		equipmentID,
		map[string]interface{}{
			"equipment_id": equipmentID,
			"due_date":     dueDate,
			"action":       "maintenance_due",
		},
	)
}