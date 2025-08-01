package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type WebhookService struct {
	client *http.Client
}

func NewWebhookService() *WebhookService {
	return &WebhookService{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *WebhookService) SendWebhook(url string, payload interface{}) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	return nil
}

func (s *WebhookService) NotifyQuoteCreated(webhookURL string, quoteID string) error {
	payload := map[string]interface{}{
		"event": "quote.created",
		"data": map[string]string{
			"quote_id": quoteID,
		},
		"timestamp": time.Now().UTC(),
	}
	return s.SendWebhook(webhookURL, payload)
}

func (s *WebhookService) TriggerQuoteCreated(quoteID string) error {
	// In a real implementation, you would get webhook URLs from configuration
	// For now, we'll just log the event
	fmt.Printf("Quote created event triggered for quote ID: %s\n", quoteID)
	return nil
}

func (s *WebhookService) TriggerQuoteSubmittedForApproval(quoteID string) error {
	// Log the event
	fmt.Printf("Quote submitted for approval event triggered for quote ID: %s\n", quoteID)
	return nil
}