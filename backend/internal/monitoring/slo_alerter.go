package monitoring

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

// WebhookSLOAlerter implements SLOAlerter using webhooks
type WebhookSLOAlerter struct {
	webhookURL string
	httpClient *http.Client
	config     *AlerterConfig
	cooldowns  map[string]time.Time
	mu         sync.RWMutex
}

// AlerterConfig holds alerter configuration
type AlerterConfig struct {
	Timeout      time.Duration
	RetryCount   int
	RetryDelay   time.Duration
	CooldownTime time.Duration
}

// NewWebhookSLOAlerter creates a new webhook-based SLO alerter
func NewWebhookSLOAlerter(webhookURL string, config *AlerterConfig) *WebhookSLOAlerter {
	if config == nil {
		config = &AlerterConfig{
			Timeout:      30 * time.Second,
			RetryCount:   3,
			RetryDelay:   5 * time.Second,
			CooldownTime: 15 * time.Minute,
		}
	}

	return &WebhookSLOAlerter{
		webhookURL: webhookURL,
		httpClient: &http.Client{Timeout: config.Timeout},
		config:     config,
		cooldowns:  make(map[string]time.Time),
	}
}

// SendAlert sends an SLO alert via webhook
func (w *WebhookSLOAlerter) SendAlert(ctx context.Context, alert *SLOAlert) error {
	// Check cooldown
	alertKey := fmt.Sprintf("%s:%s", alert.SLOName, alert.RuleName)
	if w.isInCooldown(alertKey) {
		return nil // Skip alert due to cooldown
	}

	// Convert alert to webhook payload
	payload := &WebhookPayload{
		AlertID:     alert.ID.String(),
		Type:        "slo_alert",
		Title:       fmt.Sprintf("SLO Alert: %s", alert.SLOName),
		Message:     alert.Message,
		Severity:    string(alert.Severity),
		Timestamp:   alert.Timestamp,
		SLOName:     alert.SLOName,
		RuleName:    alert.RuleName,
		Value:       alert.Value,
		Threshold:   alert.Threshold,
		Tags:        alert.Tags,
		Metadata:    alert.Metadata,
	}

	// Send with retries
	err := w.sendWithRetry(ctx, payload)
	if err == nil {
		w.setCooldown(alertKey)
	}

	return err
}

// WebhookPayload represents the webhook payload structure
type WebhookPayload struct {
	AlertID     string                 `json:"alert_id"`
	Type        string                 `json:"type"`
	Title       string                 `json:"title"`
	Message     string                 `json:"message"`
	Severity    string                 `json:"severity"`
	Timestamp   time.Time              `json:"timestamp"`
	SLOName     string                 `json:"slo_name"`
	RuleName    string                 `json:"rule_name"`
	Value       float64                `json:"value"`
	Threshold   float64                `json:"threshold"`
	Tags        map[string]string      `json:"tags"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// sendWithRetry sends webhook with retry logic
func (w *WebhookSLOAlerter) sendWithRetry(ctx context.Context, payload *WebhookPayload) error {
	var lastErr error

	for attempt := 0; attempt <= w.config.RetryCount; attempt++ {
		if attempt > 0 {
			select {
			case <-time.After(w.config.RetryDelay):
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		err := w.sendWebhook(ctx, payload)
		if err == nil {
			return nil
		}

		lastErr = err
	}

	return fmt.Errorf("failed to send alert after %d attempts: %w", w.config.RetryCount+1, lastErr)
}

// sendWebhook sends the actual webhook request
func (w *WebhookSLOAlerter) sendWebhook(ctx context.Context, payload *WebhookPayload) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", w.webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "FastenMind-SLO-Alerter/1.0")

	resp, err := w.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	return nil
}

// isInCooldown checks if an alert is in cooldown period
func (w *WebhookSLOAlerter) isInCooldown(alertKey string) bool {
	w.mu.RLock()
	defer w.mu.RUnlock()

	cooldownUntil, exists := w.cooldowns[alertKey]
	if !exists {
		return false
	}

	return time.Now().Before(cooldownUntil)
}

// setCooldown sets cooldown for an alert
func (w *WebhookSLOAlerter) setCooldown(alertKey string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.cooldowns[alertKey] = time.Now().Add(w.config.CooldownTime)
}

// MultiChannelSLOAlerter sends alerts to multiple channels
type MultiChannelSLOAlerter struct {
	alerters []SLOAlerter
	config   *MultiChannelConfig
}

// MultiChannelConfig holds multi-channel alerter configuration
type MultiChannelConfig struct {
	FailFast        bool          // Stop on first failure
	RequireAll      bool          // Require all channels to succeed
	ChannelTimeout  time.Duration // Timeout per channel
}

// NewMultiChannelSLOAlerter creates a new multi-channel alerter
func NewMultiChannelSLOAlerter(alerters []SLOAlerter, config *MultiChannelConfig) *MultiChannelSLOAlerter {
	if config == nil {
		config = &MultiChannelConfig{
			FailFast:       false,
			RequireAll:     false,
			ChannelTimeout: 30 * time.Second,
		}
	}

	return &MultiChannelSLOAlerter{
		alerters: alerters,
		config:   config,
	}
}

// SendAlert sends alert to all configured channels
func (m *MultiChannelSLOAlerter) SendAlert(ctx context.Context, alert *SLOAlert) error {
	if len(m.alerters) == 0 {
		return fmt.Errorf("no alerters configured")
	}

	type result struct {
		index int
		err   error
	}

	results := make(chan result, len(m.alerters))
	var wg sync.WaitGroup

	// Send to all channels concurrently
	for i, alerter := range m.alerters {
		wg.Add(1)
		go func(index int, a SLOAlerter) {
			defer wg.Done()

			// Create timeout context for this channel
			channelCtx, cancel := context.WithTimeout(ctx, m.config.ChannelTimeout)
			defer cancel()

			err := a.SendAlert(channelCtx, alert)
			results <- result{index: index, err: err}
		}(i, alerter)
	}

	// Wait for all to complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	var errors []error
	successCount := 0

	for res := range results {
		if res.err != nil {
			errors = append(errors, fmt.Errorf("channel %d: %w", res.index, res.err))
			if m.config.FailFast {
				return res.err
			}
		} else {
			successCount++
		}
	}

	// Check success criteria
	if m.config.RequireAll && len(errors) > 0 {
		return fmt.Errorf("not all channels succeeded: %v", errors)
	}

	if successCount == 0 {
		return fmt.Errorf("all channels failed: %v", errors)
	}

	return nil
}

// SlackSLOAlerter implements SLOAlerter for Slack
type SlackSLOAlerter struct {
	webhookURL string
	httpClient *http.Client
	config     *AlerterConfig
	cooldowns  map[string]time.Time
	mu         sync.RWMutex
}

// SlackMessage represents a Slack message payload
type SlackMessage struct {
	Text        string            `json:"text"`
	Username    string            `json:"username,omitempty"`
	Channel     string            `json:"channel,omitempty"`
	IconEmoji   string            `json:"icon_emoji,omitempty"`
	Attachments []SlackAttachment `json:"attachments,omitempty"`
}

// SlackAttachment represents a Slack message attachment
type SlackAttachment struct {
	Color     string       `json:"color,omitempty"`
	Title     string       `json:"title,omitempty"`
	Text      string       `json:"text,omitempty"`
	Fields    []SlackField `json:"fields,omitempty"`
	Timestamp int64        `json:"ts,omitempty"`
}

// SlackField represents a Slack attachment field
type SlackField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

// NewSlackSLOAlerter creates a new Slack-based SLO alerter
func NewSlackSLOAlerter(webhookURL string, config *AlerterConfig) *SlackSLOAlerter {
	if config == nil {
		config = &AlerterConfig{
			Timeout:      30 * time.Second,
			RetryCount:   3,
			RetryDelay:   5 * time.Second,
			CooldownTime: 15 * time.Minute,
		}
	}

	return &SlackSLOAlerter{
		webhookURL: webhookURL,
		httpClient: &http.Client{Timeout: config.Timeout},
		config:     config,
		cooldowns:  make(map[string]time.Time),
	}
}

// SendAlert sends an SLO alert to Slack
func (s *SlackSLOAlerter) SendAlert(ctx context.Context, alert *SLOAlert) error {
	// Check cooldown
	alertKey := fmt.Sprintf("%s:%s", alert.SLOName, alert.RuleName)
	if s.isInCooldown(alertKey) {
		return nil
	}

	// Create Slack message
	message := s.formatSlackMessage(alert)

	// Send with retries
	err := s.sendWithRetry(ctx, message)
	if err == nil {
		s.setCooldown(alertKey)
	}

	return err
}

// formatSlackMessage formats the alert as a Slack message
func (s *SlackSLOAlerter) formatSlackMessage(alert *SLOAlert) *SlackMessage {
	color := s.getColorForSeverity(alert.Severity)
	emoji := s.getEmojiForSeverity(alert.Severity)

	fields := []SlackField{
		{
			Title: "SLO",
			Value: alert.SLOName,
			Short: true,
		},
		{
			Title: "Rule",
			Value: alert.RuleName,
			Short: true,
		},
		{
			Title: "Current Value",
			Value: fmt.Sprintf("%.2f", alert.Value),
			Short: true,
		},
		{
			Title: "Threshold",
			Value: fmt.Sprintf("%.2f", alert.Threshold),
			Short: true,
		},
	}

	// Add tags as fields
	for key, value := range alert.Tags {
		fields = append(fields, SlackField{
			Title: key,
			Value: value,
			Short: true,
		})
	}

	return &SlackMessage{
		Text:     fmt.Sprintf("%s SLO Alert: %s", emoji, alert.SLOName),
		Username: "FastenMind SLO Monitor",
		Attachments: []SlackAttachment{
			{
				Color:     color,
				Title:     fmt.Sprintf("%s - %s", alert.SLOName, alert.RuleName),
				Text:      alert.Message,
				Fields:    fields,
				Timestamp: alert.Timestamp.Unix(),
			},
		},
	}
}

// getColorForSeverity returns Slack color for severity
func (s *SlackSLOAlerter) getColorForSeverity(severity AlertSeverity) string {
	switch severity {
	case AlertSeverityCritical:
		return "danger"
	case AlertSeverityWarning:
		return "warning"
	case AlertSeverityInfo:
		return "good"
	default:
		return "#439FE0"
	}
}

// getEmojiForSeverity returns emoji for severity
func (s *SlackSLOAlerter) getEmojiForSeverity(severity AlertSeverity) string {
	switch severity {
	case AlertSeverityCritical:
		return "🚨"
	case AlertSeverityWarning:
		return "⚠️"
	case AlertSeverityInfo:
		return "ℹ️"
	default:
		return "📊"
	}
}

// sendWithRetry sends Slack message with retry logic
func (s *SlackSLOAlerter) sendWithRetry(ctx context.Context, message *SlackMessage) error {
	var lastErr error

	for attempt := 0; attempt <= s.config.RetryCount; attempt++ {
		if attempt > 0 {
			select {
			case <-time.After(s.config.RetryDelay):
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		err := s.sendSlackMessage(ctx, message)
		if err == nil {
			return nil
		}

		lastErr = err
	}

	return fmt.Errorf("failed to send Slack alert after %d attempts: %w", s.config.RetryCount+1, lastErr)
}

// sendSlackMessage sends the actual Slack message
func (s *SlackSLOAlerter) sendSlackMessage(ctx context.Context, message *SlackMessage) error {
	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("Slack returned status %d", resp.StatusCode)
	}

	return nil
}

// isInCooldown and setCooldown methods (same as WebhookSLOAlerter)
func (s *SlackSLOAlerter) isInCooldown(alertKey string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cooldownUntil, exists := s.cooldowns[alertKey]
	if !exists {
		return false
	}

	return time.Now().Before(cooldownUntil)
}

func (s *SlackSLOAlerter) setCooldown(alertKey string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cooldowns[alertKey] = time.Now().Add(s.config.CooldownTime)
}

// EmailSLOAlerter implements SLOAlerter for email notifications
type EmailSLOAlerter struct {
	smtpConfig *SMTPConfig
	config     *AlerterConfig
	cooldowns  map[string]time.Time
	mu         sync.RWMutex
}

// SMTPConfig holds SMTP configuration
type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	To       []string
}

// NewEmailSLOAlerter creates a new email-based SLO alerter
func NewEmailSLOAlerter(smtpConfig *SMTPConfig, config *AlerterConfig) *EmailSLOAlerter {
	if config == nil {
		config = &AlerterConfig{
			Timeout:      30 * time.Second,
			RetryCount:   3,
			RetryDelay:   5 * time.Second,
			CooldownTime: 15 * time.Minute,
		}
	}

	return &EmailSLOAlerter{
		smtpConfig: smtpConfig,
		config:     config,
		cooldowns:  make(map[string]time.Time),
	}
}

// SendAlert sends an SLO alert via email
func (e *EmailSLOAlerter) SendAlert(ctx context.Context, alert *SLOAlert) error {
	// Check cooldown
	alertKey := fmt.Sprintf("%s:%s", alert.SLOName, alert.RuleName)
	if e.isInCooldown(alertKey) {
		return nil
	}

	// Format email
	subject := fmt.Sprintf("[%s] SLO Alert: %s", alert.Severity, alert.SLOName)
	body := e.formatEmailBody(alert)

	// Send email (implementation would depend on email library)
	err := e.sendEmail(ctx, subject, body)
	if err == nil {
		e.setCooldown(alertKey)
	}

	return err
}

// formatEmailBody formats the alert as an email body
func (e *EmailSLOAlerter) formatEmailBody(alert *SLOAlert) string {
	return fmt.Sprintf(`
SLO Alert: %s

Details:
- SLO: %s
- Rule: %s
- Severity: %s
- Message: %s
- Current Value: %.2f
- Threshold: %.2f
- Timestamp: %s

Tags:
%s

This alert was generated by the FastenMind SLO monitoring system.
`, alert.SLOName, alert.SLOName, alert.RuleName, alert.Severity, alert.Message,
		alert.Value, alert.Threshold, alert.Timestamp.Format(time.RFC3339),
		e.formatTags(alert.Tags))
}

// formatTags formats tags for email display
func (e *EmailSLOAlerter) formatTags(tags map[string]string) string {
	if len(tags) == 0 {
		return "None"
	}

	var result string
	for key, value := range tags {
		result += fmt.Sprintf("- %s: %s\n", key, value)
	}
	return result
}

// sendEmail sends the actual email (placeholder implementation)
func (e *EmailSLOAlerter) sendEmail(ctx context.Context, subject, body string) error {
	// This would implement actual SMTP sending
	// For now, return nil to indicate success
	return nil
}

// isInCooldown and setCooldown methods (same as other alerters)
func (e *EmailSLOAlerter) isInCooldown(alertKey string) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()

	cooldownUntil, exists := e.cooldowns[alertKey]
	if !exists {
		return false
	}

	return time.Now().Before(cooldownUntil)
}

func (e *EmailSLOAlerter) setCooldown(alertKey string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.cooldowns[alertKey] = time.Now().Add(e.config.CooldownTime)
}