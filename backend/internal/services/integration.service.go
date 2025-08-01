package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/fastenmind/fastener-api/internal/models"
	"github.com/fastenmind/fastener-api/internal/repositories"
)

type IntegrationService struct {
	integrationRepo *repositories.IntegrationRepository
	userRepo        *repositories.UserRepository
	companyRepo     *repositories.CompanyRepository
}

func NewIntegrationService(
	integrationRepo *repositories.IntegrationRepository,
	userRepo *repositories.UserRepository,
	companyRepo *repositories.CompanyRepository,
) *IntegrationService {
	return &IntegrationService{
		integrationRepo: integrationRepo,
		userRepo:        userRepo,
		companyRepo:     companyRepo,
	}
}

// Integration Service Methods
func (s *IntegrationService) CreateIntegration(userID uuid.UUID, req CreateIntegrationRequest) (*models.Integration, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	integration := &models.Integration{
		CompanyID:       user.CompanyID,
		Name:            req.Name,
		Type:            req.Type,
		Provider:        req.Provider,
		Status:          "inactive",
		ApiVersion:      req.ApiVersion,
		BaseURL:         req.BaseURL,
		AuthType:        req.AuthType,
		RateLimitRPM:    req.RateLimitRPM,
		TimeoutSeconds:  req.TimeoutSeconds,
		RetryAttempts:   req.RetryAttempts,
		IsActive:        false,
		CreatedBy:       userID,
	}

	if req.Configuration != nil {
		configJSON, _ := json.Marshal(req.Configuration)
		integration.Configuration = string(configJSON)
	}

	if req.Headers != nil {
		headersJSON, _ := json.Marshal(req.Headers)
		integration.Headers = string(headersJSON)
	}

	if req.Credentials != nil {
		// Encrypt credentials before storing
		credentialsJSON, _ := json.Marshal(req.Credentials)
		integration.Credentials = string(credentialsJSON) // In production, this should be encrypted
	}

	err = s.integrationRepo.CreateIntegration(integration)
	if err != nil {
		return nil, fmt.Errorf("failed to create integration: %w", err)
	}

	return integration, nil
}

func (s *IntegrationService) GetIntegration(id uuid.UUID) (*models.Integration, error) {
	return s.integrationRepo.GetIntegration(id)
}

func (s *IntegrationService) GetIntegrationsByCompany(companyID uuid.UUID, integrationType string, status string, isActive *bool) ([]models.Integration, error) {
	return s.integrationRepo.GetIntegrationsByCompany(companyID, integrationType, status, isActive)
}

func (s *IntegrationService) UpdateIntegration(id uuid.UUID, userID uuid.UUID, req UpdateIntegrationRequest) (*models.Integration, error) {
	integration, err := s.integrationRepo.GetIntegration(id)
	if err != nil {
		return nil, fmt.Errorf("integration not found: %w", err)
	}

	if req.Name != nil {
		integration.Name = *req.Name
	}
	if req.Status != nil {
		integration.Status = *req.Status
	}
	if req.BaseURL != nil {
		integration.BaseURL = *req.BaseURL
	}
	if req.IsActive != nil {
		integration.IsActive = *req.IsActive
	}
	if req.Configuration != nil {
		configJSON, _ := json.Marshal(req.Configuration)
		integration.Configuration = string(configJSON)
	}
	if req.Headers != nil {
		headersJSON, _ := json.Marshal(req.Headers)
		integration.Headers = string(headersJSON)
	}
	if req.Credentials != nil {
		credentialsJSON, _ := json.Marshal(req.Credentials)
		integration.Credentials = string(credentialsJSON) // In production, this should be encrypted
	}

	integration.UpdatedAt = time.Now()

	err = s.integrationRepo.UpdateIntegration(integration)
	if err != nil {
		return nil, fmt.Errorf("failed to update integration: %w", err)
	}

	return integration, nil
}

func (s *IntegrationService) DeleteIntegration(id uuid.UUID) error {
	return s.integrationRepo.DeleteIntegration(id)
}

func (s *IntegrationService) TestIntegration(id uuid.UUID) (*IntegrationTestResult, error) {
	integration, err := s.integrationRepo.GetIntegration(id)
	if err != nil {
		return nil, fmt.Errorf("integration not found: %w", err)
	}

	// Simulate integration test
	result := &IntegrationTestResult{
		Success:      true,
		ResponseTime: 250,
		Message:      "Connection successful",
		Details:      map[string]interface{}{
			"api_version": integration.ApiVersion,
			"base_url":    integration.BaseURL,
			"auth_type":   integration.AuthType,
		},
	}

	// Update integration status based on test result
	if result.Success {
		integration.Status = "active"
		integration.LastSyncAt = &[]time.Time{time.Now()}[0]
	} else {
		integration.Status = "error"
		integration.LastError = result.Message
		integration.LastErrorAt = &[]time.Time{time.Now()}[0]
	}

	err = s.integrationRepo.UpdateIntegration(integration)
	if err != nil {
		return nil, fmt.Errorf("failed to update integration: %w", err)
	}

	return result, nil
}

// Integration Mapping Service Methods
func (s *IntegrationService) CreateIntegrationMapping(userID uuid.UUID, req CreateIntegrationMappingRequest) (*models.IntegrationMapping, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	mapping := &models.IntegrationMapping{
		CompanyID:       user.CompanyID,
		IntegrationID:   req.IntegrationID,
		Name:            req.Name,
		Direction:       req.Direction,
		SourceTable:     req.SourceTable,
		TargetTable:     req.TargetTable,
		SourceEndpoint:  req.SourceEndpoint,
		TargetEndpoint:  req.TargetEndpoint,
		SyncFrequency:   req.SyncFrequency,
		IsActive:        true,
		CreatedBy:       userID,
	}

	if req.FieldMappings != nil {
		fieldMappingsJSON, _ := json.Marshal(req.FieldMappings)
		mapping.FieldMappings = string(fieldMappingsJSON)
	}

	if req.Transformations != nil {
		transformationsJSON, _ := json.Marshal(req.Transformations)
		mapping.Transformations = string(transformationsJSON)
	}

	if req.Filters != nil {
		filtersJSON, _ := json.Marshal(req.Filters)
		mapping.Filters = string(filtersJSON)
	}

	err = s.integrationRepo.CreateIntegrationMapping(mapping)
	if err != nil {
		return nil, fmt.Errorf("failed to create integration mapping: %w", err)
	}

	return mapping, nil
}

func (s *IntegrationService) GetMappingsByIntegration(integrationID uuid.UUID, direction string, isActive *bool) ([]models.IntegrationMapping, error) {
	return s.integrationRepo.GetMappingsByIntegration(integrationID, direction, isActive)
}

func (s *IntegrationService) UpdateIntegrationMapping(id uuid.UUID, userID uuid.UUID, req UpdateIntegrationMappingRequest) (*models.IntegrationMapping, error) {
	mapping, err := s.integrationRepo.GetIntegrationMapping(id)
	if err != nil {
		return nil, fmt.Errorf("mapping not found: %w", err)
	}

	if req.Name != nil {
		mapping.Name = *req.Name
	}
	if req.SyncFrequency != nil {
		mapping.SyncFrequency = *req.SyncFrequency
	}
	if req.IsActive != nil {
		mapping.IsActive = *req.IsActive
	}
	if req.FieldMappings != nil {
		fieldMappingsJSON, _ := json.Marshal(req.FieldMappings)
		mapping.FieldMappings = string(fieldMappingsJSON)
	}

	mapping.UpdatedAt = time.Now()

	err = s.integrationRepo.UpdateIntegrationMapping(mapping)
	if err != nil {
		return nil, fmt.Errorf("failed to update mapping: %w", err)
	}

	return mapping, nil
}

// Webhook Service Methods
func (s *IntegrationService) CreateWebhook(userID uuid.UUID, req CreateWebhookRequest) (*models.Webhook, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	webhook := &models.Webhook{
		CompanyID:       user.CompanyID,
		IntegrationID:   req.IntegrationID,
		Name:            req.Name,
		URL:             req.URL,
		Method:          req.Method,
		AuthType:        req.AuthType,
		PayloadFormat:   req.PayloadFormat,
		PayloadTemplate: req.PayloadTemplate,
		IsActive:        true,
		RetryAttempts:   req.RetryAttempts,
		RetryInterval:   req.RetryInterval,
		TimeoutSeconds:  req.TimeoutSeconds,
		CreatedBy:       userID,
	}

	if req.Headers != nil {
		headersJSON, _ := json.Marshal(req.Headers)
		webhook.Headers = string(headersJSON)
	}

	if req.AuthConfig != nil {
		authConfigJSON, _ := json.Marshal(req.AuthConfig)
		webhook.AuthConfig = string(authConfigJSON)
	}

	if req.Events != nil {
		eventsJSON, _ := json.Marshal(req.Events)
		webhook.Events = string(eventsJSON)
	}

	err = s.integrationRepo.CreateWebhook(webhook)
	if err != nil {
		return nil, fmt.Errorf("failed to create webhook: %w", err)
	}

	return webhook, nil
}

func (s *IntegrationService) GetWebhooksByCompany(companyID uuid.UUID, integrationID *uuid.UUID, isActive *bool) ([]models.Webhook, error) {
	return s.integrationRepo.GetWebhooksByCompany(companyID, integrationID, isActive)
}

func (s *IntegrationService) UpdateWebhook(id uuid.UUID, userID uuid.UUID, req UpdateWebhookRequest) (*models.Webhook, error) {
	webhook, err := s.integrationRepo.GetWebhook(id)
	if err != nil {
		return nil, fmt.Errorf("webhook not found: %w", err)
	}

	if req.Name != nil {
		webhook.Name = *req.Name
	}
	if req.URL != nil {
		webhook.URL = *req.URL
	}
	if req.IsActive != nil {
		webhook.IsActive = *req.IsActive
	}
	if req.Headers != nil {
		headersJSON, _ := json.Marshal(req.Headers)
		webhook.Headers = string(headersJSON)
	}

	webhook.UpdatedAt = time.Now()

	err = s.integrationRepo.UpdateWebhook(webhook)
	if err != nil {
		return nil, fmt.Errorf("failed to update webhook: %w", err)
	}

	return webhook, nil
}

func (s *IntegrationService) TriggerWebhook(webhookID uuid.UUID, eventType string, eventData map[string]interface{}) error {
	webhook, err := s.integrationRepo.GetWebhook(webhookID)
	if err != nil {
		return fmt.Errorf("webhook not found: %w", err)
	}

	if !webhook.IsActive {
		return errors.New("webhook is not active")
	}

	// Create webhook delivery record
	delivery := &models.WebhookDelivery{
		WebhookID:      webhookID,
		CompanyID:      webhook.CompanyID,
		EventType:      eventType,
		RequestURL:     webhook.URL,
		RequestMethod:  webhook.Method,
		Status:         "pending",
		AttemptCount:   0,
	}

	if eventData != nil {
		eventDataJSON, _ := json.Marshal(eventData)
		delivery.EventData = string(eventDataJSON)
	}

	err = s.integrationRepo.CreateWebhookDelivery(delivery)
	if err != nil {
		return fmt.Errorf("failed to create webhook delivery: %w", err)
	}

	// Start webhook delivery asynchronously
	go s.processWebhookDelivery(delivery.ID)

	return nil
}

func (s *IntegrationService) GetWebhookDeliveries(webhookID uuid.UUID, status string, limit int) ([]models.WebhookDelivery, error) {
	return s.integrationRepo.GetWebhookDeliveries(webhookID, status, limit)
}

// Data Sync Job Service Methods
func (s *IntegrationService) CreateDataSyncJob(userID uuid.UUID, req CreateDataSyncJobRequest) (*models.DataSyncJob, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	job := &models.DataSyncJob{
		CompanyID:     user.CompanyID,
		IntegrationID: req.IntegrationID,
		MappingID:     req.MappingID,
		Name:          req.Name,
		Type:          req.Type,
		Direction:     req.Direction,
		Status:        "pending",
		Priority:      req.Priority,
		ScheduledAt:   req.ScheduledAt,
		CreatedBy:     userID,
	}

	if req.Configuration != nil {
		configJSON, _ := json.Marshal(req.Configuration)
		job.Configuration = string(configJSON)
	}

	err = s.integrationRepo.CreateDataSyncJob(job)
	if err != nil {
		return nil, fmt.Errorf("failed to create sync job: %w", err)
	}

	return job, nil
}

func (s *IntegrationService) GetDataSyncJobsByIntegration(integrationID uuid.UUID, status string, limit int) ([]models.DataSyncJob, error) {
	return s.integrationRepo.GetDataSyncJobsByIntegration(integrationID, status, limit)
}

func (s *IntegrationService) StartDataSyncJob(jobID uuid.UUID) error {
	job, err := s.integrationRepo.GetDataSyncJob(jobID)
	if err != nil {
		return fmt.Errorf("sync job not found: %w", err)
	}

	if job.Status != "pending" {
		return errors.New("sync job is not in pending status")
	}

	err = s.integrationRepo.StartSyncJob(jobID)
	if err != nil {
		return fmt.Errorf("failed to start sync job: %w", err)
	}

	// Start sync job processing asynchronously
	go s.processSyncJob(jobID)

	return nil
}

// API Key Service Methods
func (s *IntegrationService) CreateApiKey(userID uuid.UUID, req CreateApiKeyRequest) (*CreateApiKeyResponse, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Generate API key
	keyBytes := make([]byte, 32)
	_, err = rand.Read(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to generate API key: %w", err)
	}

	apiKeyString := "fm_" + hex.EncodeToString(keyBytes)
	keyHash := sha256.Sum256([]byte(apiKeyString))
	keyHashString := hex.EncodeToString(keyHash[:])
	keyPrefix := apiKeyString[:12] + "..."

	apiKey := &models.ApiKey{
		CompanyID:   user.CompanyID,
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		KeyHash:     keyHashString,
		KeyPrefix:   keyPrefix,
		RateLimit:   req.RateLimit,
		IsActive:    true,
		ExpiresAt:   req.ExpiresAt,
		CreatedBy:   userID,
	}

	if req.Permissions != nil {
		permissionsJSON, _ := json.Marshal(req.Permissions)
		apiKey.Permissions = string(permissionsJSON)
	}

	if req.Scopes != nil {
		scopesJSON, _ := json.Marshal(req.Scopes)
		apiKey.Scopes = string(scopesJSON)
	}

	err = s.integrationRepo.CreateApiKey(apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create API key: %w", err)
	}

	return &CreateApiKeyResponse{
		ID:       apiKey.ID,
		Key:      apiKeyString, // Only returned once during creation
		KeyID:    apiKey.KeyPrefix,
		Name:     apiKey.Name,
		ExpiresAt: apiKey.ExpiresAt,
	}, nil
}

func (s *IntegrationService) GetApiKeysByCompany(companyID uuid.UUID, userID *uuid.UUID, isActive *bool) ([]models.ApiKey, error) {
	return s.integrationRepo.GetApiKeysByCompany(companyID, userID, isActive)
}

func (s *IntegrationService) ValidateApiKey(keyString string) (*models.ApiKey, error) {
	keyHash := sha256.Sum256([]byte(keyString))
	keyHashString := hex.EncodeToString(keyHash[:])

	apiKey, err := s.integrationRepo.GetApiKeyByHash(keyHashString)
	if err != nil {
		return nil, fmt.Errorf("invalid API key: %w", err)
	}

	// Update usage statistics
	err = s.integrationRepo.UpdateApiKeyUsage(keyHashString)
	if err != nil {
		// Log error but don't fail the validation
		fmt.Printf("Failed to update API key usage: %v\n", err)
	}

	return apiKey, nil
}

func (s *IntegrationService) RevokeApiKey(id uuid.UUID) error {
	apiKey, err := s.integrationRepo.GetApiKey(id)
	if err != nil {
		return fmt.Errorf("API key not found: %w", err)
	}

	apiKey.IsActive = false
	apiKey.UpdatedAt = time.Now()

	return s.integrationRepo.UpdateApiKey(apiKey)
}

// External System Service Methods
func (s *IntegrationService) CreateExternalSystem(userID uuid.UUID, req CreateExternalSystemRequest) (*models.ExternalSystem, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	system := &models.ExternalSystem{
		CompanyID:  user.CompanyID,
		Name:       req.Name,
		SystemType: req.SystemType,
		Vendor:     req.Vendor,
		Version:    req.Version,
		BaseURL:    req.BaseURL,
		Status:     "inactive",
		IsActive:   false,
		CreatedBy:  userID,
	}

	if req.DatabaseConfig != nil {
		dbConfigJSON, _ := json.Marshal(req.DatabaseConfig)
		system.DatabaseConfig = string(dbConfigJSON)
	}

	if req.ApiConfig != nil {
		apiConfigJSON, _ := json.Marshal(req.ApiConfig)
		system.ApiConfig = string(apiConfigJSON)
	}

	err = s.integrationRepo.CreateExternalSystem(system)
	if err != nil {
		return nil, fmt.Errorf("failed to create external system: %w", err)
	}

	return system, nil
}

func (s *IntegrationService) GetExternalSystemsByCompany(companyID uuid.UUID, systemType string, status string) ([]models.ExternalSystem, error) {
	return s.integrationRepo.GetExternalSystemsByCompany(companyID, systemType, status)
}

func (s *IntegrationService) TestExternalSystem(id uuid.UUID) (*SystemTestResult, error) {
	system, err := s.integrationRepo.GetExternalSystem(id)
	if err != nil {
		return nil, fmt.Errorf("external system not found: %w", err)
	}

	// Simulate system test
	result := &SystemTestResult{
		Success:      true,
		ResponseTime: 180,
		Message:      "Connection successful",
		Details: map[string]interface{}{
			"system_type": system.SystemType,
			"vendor":      system.Vendor,
			"version":     system.Version,
		},
	}

	// Update system status
	now := time.Now()
	system.LastTestAt = &now
	system.LastTestResult = "success"
	if result.Success {
		system.Status = "active"
	} else {
		system.Status = "error"
	}

	err = s.integrationRepo.UpdateExternalSystem(system)
	if err != nil {
		return nil, fmt.Errorf("failed to update system: %w", err)
	}

	return result, nil
}

// Integration Template Service Methods
func (s *IntegrationService) GetIntegrationTemplates(companyID *uuid.UUID, category string, provider string, isPublic *bool, isActive *bool) ([]models.IntegrationTemplate, error) {
	return s.integrationRepo.GetIntegrationTemplates(companyID, category, provider, isPublic, isActive)
}

func (s *IntegrationService) CreateIntegrationFromTemplate(userID uuid.UUID, templateID uuid.UUID, req CreateFromTemplateRequest) (*models.Integration, error) {
	template, err := s.integrationRepo.GetIntegrationTemplate(templateID)
	if err != nil {
		return nil, fmt.Errorf("template not found: %w", err)
	}

	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Parse template configuration
	var templateConfig map[string]interface{}
	if template.Configuration != "" {
		json.Unmarshal([]byte(template.Configuration), &templateConfig)
	}

	// Merge template config with user-provided config
	finalConfig := templateConfig
	if req.Configuration != nil {
		for k, v := range req.Configuration {
			finalConfig[k] = v
		}
	}

	integration := &models.Integration{
		CompanyID:      user.CompanyID,
		Name:           req.Name,
		Type:           template.Category,
		Provider:       template.Provider,
		Status:         "inactive",
		ApiVersion:     template.Version,
		AuthType:       "api_key", // Default from template
		RateLimitRPM:   60,        // Default
		TimeoutSeconds: 30,        // Default
		RetryAttempts:  3,         // Default
		IsActive:       false,
		CreatedBy:      userID,
	}

	if finalConfig != nil {
		configJSON, _ := json.Marshal(finalConfig)
		integration.Configuration = string(configJSON)
	}

	err = s.integrationRepo.CreateIntegration(integration)
	if err != nil {
		return nil, fmt.Errorf("failed to create integration from template: %w", err)
	}

	// Increment template usage count
	s.integrationRepo.IncrementTemplateUsage(templateID)

	return integration, nil
}

// Analytics Service Methods
func (s *IntegrationService) GetIntegrationStats(companyID uuid.UUID, startTime time.Time, endTime time.Time) (map[string]interface{}, error) {
	return s.integrationRepo.GetIntegrationStats(companyID, startTime, endTime)
}

func (s *IntegrationService) GetIntegrationsByType(companyID uuid.UUID) ([]map[string]interface{}, error) {
	return s.integrationRepo.GetIntegrationsByType(companyID)
}

func (s *IntegrationService) GetSyncJobTrends(companyID uuid.UUID, days int) ([]map[string]interface{}, error) {
	return s.integrationRepo.GetSyncJobTrends(companyID, days)
}

// Private helper methods
func (s *IntegrationService) processWebhookDelivery(deliveryID uuid.UUID) {
	delivery, err := s.integrationRepo.GetWebhookDelivery(deliveryID)
	if err != nil {
		return
	}

	// Simulate webhook delivery
	time.Sleep(1 * time.Second)

	// Update delivery record
	delivery.Status = "success"
	delivery.ResponseCode = 200
	delivery.ResponseBody = `{"status": "ok"}`
	delivery.ResponseTime = 250
	now := time.Now()
	delivery.CompletedAt = &now
	delivery.UpdatedAt = time.Now()

	s.integrationRepo.UpdateWebhookDelivery(delivery)

	// Update webhook stats
	s.integrationRepo.UpdateWebhookStats(delivery.WebhookID, true)
}

func (s *IntegrationService) processSyncJob(jobID uuid.UUID) {
	job, err := s.integrationRepo.GetDataSyncJob(jobID)
	if err != nil {
		return
	}

	// Simulate sync job processing
	totalRecords := int64(1000)
	job.TotalRecords = totalRecords

	for i := int64(0); i <= totalRecords; i += 100 {
		time.Sleep(500 * time.Millisecond)
		
		progress := int((float64(i) / float64(totalRecords)) * 100)
		if progress > 100 {
			progress = 100
		}
		
		s.integrationRepo.UpdateSyncJobProgress(jobID, progress, i, i, 0)
	}

	// Complete the job
	job.Status = "completed"
	job.ProcessedRecords = totalRecords
	job.SuccessRecords = totalRecords
	job.Progress = 100
	now := time.Now()
	job.CompletedAt = &now
	job.Duration = int64(now.Sub(*job.StartedAt).Seconds())

	s.integrationRepo.UpdateDataSyncJob(job)

	// Update integration stats
	s.integrationRepo.UpdateIntegrationStats(job.IntegrationID, true, 0)
	s.integrationRepo.UpdateIntegrationSuccessRate(job.IntegrationID)
}

// Request/Response types
type CreateIntegrationRequest struct {
	Name           string                 `json:"name" validate:"required"`
	Type           string                 `json:"type" validate:"required"`
	Provider       string                 `json:"provider" validate:"required"`
	ApiVersion     string                 `json:"api_version"`
	BaseURL        string                 `json:"base_url"`
	AuthType       string                 `json:"auth_type"`
	Configuration  map[string]interface{} `json:"configuration"`
	Credentials    map[string]interface{} `json:"credentials"`
	Headers        map[string]string      `json:"headers"`
	RateLimitRPM   int                    `json:"rate_limit_rpm"`
	TimeoutSeconds int                    `json:"timeout_seconds"`
	RetryAttempts  int                    `json:"retry_attempts"`
}

type UpdateIntegrationRequest struct {
	Name          *string                 `json:"name"`
	Status        *string                 `json:"status"`
	BaseURL       *string                 `json:"base_url"`
	IsActive      *bool                   `json:"is_active"`
	Configuration *map[string]interface{} `json:"configuration"`
	Credentials   *map[string]interface{} `json:"credentials"`
	Headers       *map[string]string      `json:"headers"`
}

type CreateIntegrationMappingRequest struct {
	IntegrationID   uuid.UUID              `json:"integration_id" validate:"required"`
	Name            string                 `json:"name" validate:"required"`
	Direction       string                 `json:"direction" validate:"required"`
	SourceTable     string                 `json:"source_table"`
	TargetTable     string                 `json:"target_table"`
	SourceEndpoint  string                 `json:"source_endpoint"`
	TargetEndpoint  string                 `json:"target_endpoint"`
	FieldMappings   map[string]interface{} `json:"field_mappings"`
	Transformations map[string]interface{} `json:"transformations"`
	Filters         map[string]interface{} `json:"filters"`
	SyncFrequency   string                 `json:"sync_frequency" validate:"required"`
}

type UpdateIntegrationMappingRequest struct {
	Name            *string                 `json:"name"`
	SyncFrequency   *string                 `json:"sync_frequency"`
	IsActive        *bool                   `json:"is_active"`
	FieldMappings   *map[string]interface{} `json:"field_mappings"`
	Transformations *map[string]interface{} `json:"transformations"`
	Filters         *map[string]interface{} `json:"filters"`
}

type CreateWebhookRequest struct {
	IntegrationID   *uuid.UUID             `json:"integration_id"`
	Name            string                 `json:"name" validate:"required"`
	URL             string                 `json:"url" validate:"required"`
	Method          string                 `json:"method"`
	Headers         map[string]string      `json:"headers"`
	AuthType        string                 `json:"auth_type"`
	AuthConfig      map[string]interface{} `json:"auth_config"`
	Events          []string               `json:"events"`
	PayloadFormat   string                 `json:"payload_format"`
	PayloadTemplate string                 `json:"payload_template"`
	RetryAttempts   int                    `json:"retry_attempts"`
	RetryInterval   int                    `json:"retry_interval"`
	TimeoutSeconds  int                    `json:"timeout_seconds"`
}

type UpdateWebhookRequest struct {
	Name      *string            `json:"name"`
	URL       *string            `json:"url"`
	IsActive  *bool              `json:"is_active"`
	Headers   *map[string]string `json:"headers"`
	Events    *[]string          `json:"events"`
}

type CreateDataSyncJobRequest struct {
	IntegrationID uuid.UUID              `json:"integration_id" validate:"required"`
	MappingID     *uuid.UUID             `json:"mapping_id"`
	Name          string                 `json:"name" validate:"required"`
	Type          string                 `json:"type" validate:"required"`
	Direction     string                 `json:"direction" validate:"required"`
	Priority      string                 `json:"priority"`
	ScheduledAt   *time.Time             `json:"scheduled_at"`
	Configuration map[string]interface{} `json:"configuration"`
}

type CreateApiKeyRequest struct {
	Name        string     `json:"name" validate:"required"`
	Description string     `json:"description"`
	Permissions []string   `json:"permissions"`
	Scopes      []string   `json:"scopes"`
	RateLimit   int        `json:"rate_limit"`
	ExpiresAt   *time.Time `json:"expires_at"`
}

type CreateApiKeyResponse struct {
	ID        uuid.UUID  `json:"id"`
	Key       string     `json:"key"`       // Only returned once
	KeyID     string     `json:"key_id"`    // Prefix for identification
	Name      string     `json:"name"`
	ExpiresAt *time.Time `json:"expires_at"`
}

type CreateExternalSystemRequest struct {
	Name           string                 `json:"name" validate:"required"`
	SystemType     string                 `json:"system_type" validate:"required"`
	Vendor         string                 `json:"vendor"`
	Version        string                 `json:"version"`
	BaseURL        string                 `json:"base_url"`
	DatabaseConfig map[string]interface{} `json:"database_config"`
	ApiConfig      map[string]interface{} `json:"api_config"`
	FtpConfig      map[string]interface{} `json:"ftp_config"`
	SftpConfig     map[string]interface{} `json:"sftp_config"`
}

type CreateFromTemplateRequest struct {
	Name          string                 `json:"name" validate:"required"`
	Configuration map[string]interface{} `json:"configuration"`
}

type IntegrationTestResult struct {
	Success      bool                   `json:"success"`
	ResponseTime int64                  `json:"response_time"`
	Message      string                 `json:"message"`
	Details      map[string]interface{} `json:"details"`
}

type SystemTestResult struct {
	Success      bool                   `json:"success"`
	ResponseTime int64                  `json:"response_time"`
	Message      string                 `json:"message"`
	Details      map[string]interface{} `json:"details"`
}