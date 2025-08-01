package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"fastenmind-system/internal/models"
	"fastenmind-system/internal/repositories"
)

type AdvancedService struct {
	advancedRepo   *repositories.AdvancedRepository
	userRepo       *repositories.UserRepository
	companyRepo    *repositories.CompanyRepository
}

func NewAdvancedService(
	advancedRepo *repositories.AdvancedRepository,
	userRepo *repositories.UserRepository,
	companyRepo *repositories.CompanyRepository,
) *AdvancedService {
	return &AdvancedService{
		advancedRepo: advancedRepo,
		userRepo:     userRepo,
		companyRepo:  companyRepo,
	}
}

// AI Assistant Service Methods
func (s *AdvancedService) CreateAIAssistant(userID uuid.UUID, req CreateAIAssistantRequest) (*models.AIAssistant, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	assistant := &models.AIAssistant{
		CompanyID:        user.CompanyID,
		UserID:           userID,
		Name:             req.Name,
		Type:             req.Type,
		Model:            req.Model,
		Status:           "active",
		SystemPrompt:     req.SystemPrompt,
		Temperature:      req.Temperature,
		MaxTokens:        req.MaxTokens,
		TopP:             req.TopP,
		FrequencyPenalty: req.FrequencyPenalty,
		PresencePenalty:  req.PresencePenalty,
		IsActive:         true,
		CreatedBy:        userID,
	}

	if req.Configuration != nil {
		configJSON, _ := json.Marshal(req.Configuration)
		assistant.Configuration = string(configJSON)
	}

	err = s.advancedRepo.CreateAIAssistant(assistant)
	if err != nil {
		return nil, fmt.Errorf("failed to create AI assistant: %w", err)
	}

	return assistant, nil
}

func (s *AdvancedService) GetAIAssistant(id uuid.UUID) (*models.AIAssistant, error) {
	return s.advancedRepo.GetAIAssistant(id)
}

func (s *AdvancedService) GetAIAssistantsByCompany(companyID uuid.UUID, assistantType string, isActive *bool) ([]models.AIAssistant, error) {
	return s.advancedRepo.GetAIAssistantsByCompany(companyID, assistantType, isActive)
}

func (s *AdvancedService) UpdateAIAssistant(id uuid.UUID, userID uuid.UUID, req UpdateAIAssistantRequest) (*models.AIAssistant, error) {
	assistant, err := s.advancedRepo.GetAIAssistant(id)
	if err != nil {
		return nil, fmt.Errorf("assistant not found: %w", err)
	}

	if req.Name != nil {
		assistant.Name = *req.Name
	}
	if req.SystemPrompt != nil {
		assistant.SystemPrompt = *req.SystemPrompt
	}
	if req.Temperature != nil {
		assistant.Temperature = *req.Temperature
	}
	if req.MaxTokens != nil {
		assistant.MaxTokens = *req.MaxTokens
	}
	if req.IsActive != nil {
		assistant.IsActive = *req.IsActive
	}
	if req.Configuration != nil {
		configJSON, _ := json.Marshal(req.Configuration)
		assistant.Configuration = string(configJSON)
	}

	assistant.UpdatedAt = time.Now()

	err = s.advancedRepo.UpdateAIAssistant(assistant)
	if err != nil {
		return nil, fmt.Errorf("failed to update AI assistant: %w", err)
	}

	return assistant, nil
}

func (s *AdvancedService) DeleteAIAssistant(id uuid.UUID) error {
	return s.advancedRepo.DeleteAIAssistant(id)
}

// AI Conversation Service Methods
func (s *AdvancedService) StartConversation(userID uuid.UUID, assistantID uuid.UUID, title string, context map[string]interface{}) (*models.AIConversationSession, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	assistant, err := s.advancedRepo.GetAIAssistant(assistantID)
	if err != nil {
		return nil, fmt.Errorf("assistant not found: %w", err)
	}

	if !assistant.IsActive {
		return nil, errors.New("assistant is not active")
	}

	session := &models.AIConversationSession{
		AssistantID: assistantID,
		UserID:      userID,
		CompanyID:   user.CompanyID,
		Title:       title,
		Status:      "active",
		StartTime:   time.Now(),
	}

	if context != nil {
		contextJSON, _ := json.Marshal(context)
		session.Context = string(contextJSON)
	}

	err = s.advancedRepo.CreateConversationSession(session)
	if err != nil {
		return nil, fmt.Errorf("failed to create conversation session: %w", err)
	}

	return session, nil
}

func (s *AdvancedService) SendMessage(sessionID uuid.UUID, userID uuid.UUID, content string) (*AIMessageResponse, error) {
	session, err := s.advancedRepo.GetConversationSession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	if session.UserID != userID {
		return nil, errors.New("unauthorized access to session")
	}

	if session.Status != "active" {
		return nil, errors.New("session is not active")
	}

	// Create user message
	userMessage := &models.AIMessage{
		SessionID: sessionID,
		Role:      "user",
		Content:   content,
		CreatedAt: time.Now(),
	}

	err = s.advancedRepo.CreateAIMessage(userMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to create user message: %w", err)
	}

	// Here you would integrate with AI service (OpenAI, Claude, etc.)
	// For now, we'll simulate a response
	response := s.simulateAIResponse(content, session.Assistant)

	// Create assistant message
	assistantMessage := &models.AIMessage{
		SessionID:    sessionID,
		Role:         "assistant",
		Content:      response.Content,
		TokenCount:   response.TokenCount,
		Cost:         response.Cost,
		ModelUsed:    session.Assistant.Model,
		ResponseTime: response.ResponseTime,
		CreatedAt:    time.Now(),
	}

	err = s.advancedRepo.CreateAIMessage(assistantMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to create assistant message: %w", err)
	}

	// Update session and assistant statistics
	err = s.advancedRepo.UpdateSessionStats(sessionID, int64(response.TokenCount), response.Cost)
	if err != nil {
		return nil, fmt.Errorf("failed to update session stats: %w", err)
	}

	err = s.advancedRepo.IncrementAIAssistantUsage(session.AssistantID, int64(response.TokenCount), response.Cost)
	if err != nil {
		return nil, fmt.Errorf("failed to update assistant usage: %w", err)
	}

	return response, nil
}

func (s *AdvancedService) GetConversationHistory(sessionID uuid.UUID, userID uuid.UUID, limit int) ([]models.AIMessage, error) {
	session, err := s.advancedRepo.GetConversationSession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	if session.UserID != userID {
		return nil, errors.New("unauthorized access to session")
	}

	return s.advancedRepo.GetMessagesBySession(sessionID, limit)
}

func (s *AdvancedService) EndConversation(sessionID uuid.UUID, userID uuid.UUID) error {
	session, err := s.advancedRepo.GetConversationSession(sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	if session.UserID != userID {
		return errors.New("unauthorized access to session")
	}

	return s.advancedRepo.EndConversationSession(sessionID)
}

// Smart Recommendation Service Methods
func (s *AdvancedService) CreateRecommendation(userID uuid.UUID, req CreateRecommendationRequest) (*models.SmartRecommendation, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	recommendation := &models.SmartRecommendation{
		CompanyID:    user.CompanyID,
		UserID:       userID,
		Type:         req.Type,
		Category:     req.Category,
		Title:        req.Title,
		Description:  req.Description,
		Score:        req.Score,
		Priority:     req.Priority,
		Status:       "pending",
		Source:       req.Source,
		ResourceID:   req.ResourceID,
		ResourceType: req.ResourceType,
		ExpiresAt:    req.ExpiresAt,
	}

	if req.Data != nil {
		dataJSON, _ := json.Marshal(req.Data)
		recommendation.Data = string(dataJSON)
	}

	if req.SourceData != nil {
		sourceDataJSON, _ := json.Marshal(req.SourceData)
		recommendation.SourceData = string(sourceDataJSON)
	}

	err = s.advancedRepo.CreateRecommendation(recommendation)
	if err != nil {
		return nil, fmt.Errorf("failed to create recommendation: %w", err)
	}

	return recommendation, nil
}

func (s *AdvancedService) GetRecommendationsByUser(userID uuid.UUID, recommendationType string, status string, limit int) ([]models.SmartRecommendation, error) {
	return s.advancedRepo.GetRecommendationsByUser(userID, recommendationType, status, limit)
}

func (s *AdvancedService) UpdateRecommendationStatus(id uuid.UUID, userID uuid.UUID, status string) (*models.SmartRecommendation, error) {
	recommendation, err := s.advancedRepo.GetRecommendation(id)
	if err != nil {
		return nil, fmt.Errorf("recommendation not found: %w", err)
	}

	if recommendation.UserID != userID {
		return nil, errors.New("unauthorized access to recommendation")
	}

	validStatuses := []string{"pending", "viewed", "accepted", "rejected", "implemented"}
	if !contains(validStatuses, status) {
		return nil, errors.New("invalid status")
	}

	err = s.advancedRepo.UpdateRecommendationStatus(id, status)
	if err != nil {
		return nil, fmt.Errorf("failed to update recommendation status: %w", err)
	}

	return s.advancedRepo.GetRecommendation(id)
}

// Advanced Search Service Methods
func (s *AdvancedService) CreateAdvancedSearch(userID uuid.UUID, req CreateAdvancedSearchRequest) (*models.AdvancedSearch, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	search := &models.AdvancedSearch{
		CompanyID:    user.CompanyID,
		UserID:       userID,
		Name:         req.Name,
		Description:  req.Description,
		SearchType:   req.SearchType,
		TableName:    req.TableName,
		IsPublic:     req.IsPublic,
	}

	if req.Filters != nil {
		filtersJSON, _ := json.Marshal(req.Filters)
		search.Filters = string(filtersJSON)
	}

	if req.SortConfig != nil {
		sortConfigJSON, _ := json.Marshal(req.SortConfig)
		search.SortConfig = string(sortConfigJSON)
	}

	if req.Columns != nil {
		columnsJSON, _ := json.Marshal(req.Columns)
		search.Columns = string(columnsJSON)
	}

	err = s.advancedRepo.CreateAdvancedSearch(search)
	if err != nil {
		return nil, fmt.Errorf("failed to create advanced search: %w", err)
	}

	return search, nil
}

func (s *AdvancedService) GetAdvancedSearchesByUser(userID uuid.UUID, searchType string, isPublic *bool) ([]models.AdvancedSearch, error) {
	return s.advancedRepo.GetAdvancedSearchesByUser(userID, searchType, isPublic)
}

func (s *AdvancedService) ExecuteAdvancedSearch(searchID uuid.UUID, userID uuid.UUID) (*AdvancedSearchResult, error) {
	search, err := s.advancedRepo.GetAdvancedSearch(searchID)
	if err != nil {
		return nil, fmt.Errorf("search not found: %w", err)
	}

	if search.UserID != userID && !search.IsPublic {
		return nil, errors.New("unauthorized access to search")
	}

	// Increment usage count
	err = s.advancedRepo.IncrementSearchUsage(searchID)
	if err != nil {
		return nil, fmt.Errorf("failed to update search usage: %w", err)
	}

	// Execute the search (implementation would depend on the search engine)
	result := &AdvancedSearchResult{
		SearchID:    searchID,
		ExecutedAt:  time.Now(),
		TotalCount:  0,
		Results:     []map[string]interface{}{},
		ExecutionTime: 100, // milliseconds
	}

	return result, nil
}

// Batch Operation Service Methods
func (s *AdvancedService) CreateBatchOperation(userID uuid.UUID, req CreateBatchOperationRequest) (*models.BatchOperation, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	operation := &models.BatchOperation{
		CompanyID:     user.CompanyID,
		UserID:        userID,
		OperationType: req.OperationType,
		TargetTable:   req.TargetTable,
		Status:        "pending",
		TotalItems:    len(req.TargetIDs),
	}

	if len(req.TargetIDs) > 0 {
		targetIDsJSON, _ := json.Marshal(req.TargetIDs)
		operation.TargetIDs = string(targetIDsJSON)
	}

	if req.Parameters != nil {
		parametersJSON, _ := json.Marshal(req.Parameters)
		operation.Parameters = string(parametersJSON)
	}

	err = s.advancedRepo.CreateBatchOperation(operation)
	if err != nil {
		return nil, fmt.Errorf("failed to create batch operation: %w", err)
	}

	// Start processing asynchronously
	go s.processBatchOperation(operation.ID)

	return operation, nil
}

func (s *AdvancedService) GetBatchOperationsByUser(userID uuid.UUID, status string, limit int) ([]models.BatchOperation, error) {
	return s.advancedRepo.GetBatchOperationsByUser(userID, status, limit)
}

func (s *AdvancedService) GetBatchOperation(id uuid.UUID, userID uuid.UUID) (*models.BatchOperation, error) {
	operation, err := s.advancedRepo.GetBatchOperation(id)
	if err != nil {
		return nil, fmt.Errorf("operation not found: %w", err)
	}

	if operation.UserID != userID {
		return nil, errors.New("unauthorized access to operation")
	}

	return operation, nil
}

// Custom Field Service Methods
func (s *AdvancedService) CreateCustomField(userID uuid.UUID, req CreateCustomFieldRequest) (*models.CustomField, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	field := &models.CustomField{
		CompanyID:    user.CompanyID,
		TableName:    req.TableName,
		FieldName:    req.FieldName,
		FieldLabel:   req.FieldLabel,
		FieldType:    req.FieldType,
		DefaultValue: req.DefaultValue,
		IsRequired:   req.IsRequired,
		IsSearchable: req.IsSearchable,
		IsActive:     true,
		DisplayOrder: req.DisplayOrder,
		CreatedBy:    userID,
	}

	if req.Options != nil {
		optionsJSON, _ := json.Marshal(req.Options)
		field.Options = string(optionsJSON)
	}

	if req.Validation != nil {
		validationJSON, _ := json.Marshal(req.Validation)
		field.Validation = string(validationJSON)
	}

	err = s.advancedRepo.CreateCustomField(field)
	if err != nil {
		return nil, fmt.Errorf("failed to create custom field: %w", err)
	}

	return field, nil
}

func (s *AdvancedService) GetCustomFieldsByTable(companyID uuid.UUID, tableName string, isActive *bool) ([]models.CustomField, error) {
	return s.advancedRepo.GetCustomFieldsByTable(companyID, tableName, isActive)
}

func (s *AdvancedService) SetCustomFieldValue(userID uuid.UUID, req SetCustomFieldValueRequest) error {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	field, err := s.advancedRepo.GetCustomField(req.FieldID)
	if err != nil {
		return fmt.Errorf("field not found: %w", err)
	}

	if field.CompanyID != user.CompanyID {
		return errors.New("unauthorized access to field")
	}

	value := &models.CustomFieldValue{
		CompanyID:    user.CompanyID,
		FieldID:      req.FieldID,
		ResourceID:   req.ResourceID,
		ResourceType: req.ResourceType,
		Value:        req.Value,
	}

	return s.advancedRepo.SetCustomFieldValue(value)
}

func (s *AdvancedService) GetCustomFieldValues(resourceID uuid.UUID, resourceType string) ([]models.CustomFieldValue, error) {
	return s.advancedRepo.GetCustomFieldValues(resourceID, resourceType)
}

// Security Event Service Methods
func (s *AdvancedService) CreateSecurityEvent(req CreateSecurityEventRequest) (*models.SecurityEvent, error) {
	event := &models.SecurityEvent{
		CompanyID:     req.CompanyID,
		UserID:        req.UserID,
		EventType:     req.EventType,
		Severity:      req.Severity,
		Description:   req.Description,
		IPAddress:     req.IPAddress,
		UserAgent:     req.UserAgent,
		Location:      req.Location,
		ResourceType:  req.ResourceType,
		ResourceID:    req.ResourceID,
		RiskScore:     req.RiskScore,
		Status:        "new",
	}

	if req.DeviceInfo != nil {
		deviceInfoJSON, _ := json.Marshal(req.DeviceInfo)
		event.DeviceInfo = string(deviceInfoJSON)
	}

	if req.ActionDetails != nil {
		actionDetailsJSON, _ := json.Marshal(req.ActionDetails)
		event.ActionDetails = string(actionDetailsJSON)
	}

	err := s.advancedRepo.CreateSecurityEvent(event)
	if err != nil {
		return nil, fmt.Errorf("failed to create security event: %w", err)
	}

	return event, nil
}

func (s *AdvancedService) GetSecurityEventsByCompany(companyID uuid.UUID, eventType string, severity string, status string, limit int) ([]models.SecurityEvent, error) {
	return s.advancedRepo.GetSecurityEventsByCompany(companyID, eventType, severity, status, limit)
}

// Performance Metric Service Methods
func (s *AdvancedService) RecordPerformanceMetric(req RecordPerformanceMetricRequest) error {
	metric := &models.PerformanceMetric{
		CompanyID:  req.CompanyID,
		MetricType: req.MetricType,
		MetricName: req.MetricName,
		Value:      req.Value,
		Unit:       req.Unit,
		Endpoint:   req.Endpoint,
		Method:     req.Method,
		StatusCode: req.StatusCode,
		UserID:     req.UserID,
		SessionID:  req.SessionID,
		TraceID:    req.TraceID,
	}

	if req.Context != nil {
		contextJSON, _ := json.Marshal(req.Context)
		metric.Context = string(contextJSON)
	}

	return s.advancedRepo.CreatePerformanceMetric(metric)
}

func (s *AdvancedService) GetPerformanceStats(companyID uuid.UUID, metricType string, startTime time.Time, endTime time.Time) (map[string]interface{}, error) {
	return s.advancedRepo.GetPerformanceStats(companyID, metricType, startTime, endTime)
}

// Backup Service Methods
func (s *AdvancedService) CreateBackupRecord(userID uuid.UUID, req CreateBackupRecordRequest) (*models.BackupRecord, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	record := &models.BackupRecord{
		CompanyID:       user.CompanyID,
		BackupType:      req.BackupType,
		Status:          "running",
		StartTime:       time.Now(),
		CompressionType: req.CompressionType,
		EncryptionType:  req.EncryptionType,
		CreatedBy:       userID,
	}

	if req.Tables != nil {
		tablesJSON, _ := json.Marshal(req.Tables)
		record.Tables = string(tablesJSON)
	}

	err = s.advancedRepo.CreateBackupRecord(record)
	if err != nil {
		return nil, fmt.Errorf("failed to create backup record: %w", err)
	}

	// Start backup process asynchronously
	go s.processBackup(record.ID)

	return record, nil
}

func (s *AdvancedService) GetBackupRecordsByCompany(companyID uuid.UUID, backupType string, status string, limit int) ([]models.BackupRecord, error) {
	return s.advancedRepo.GetBackupRecordsByCompany(companyID, backupType, status, limit)
}

// Multi-language Support Service Methods
func (s *AdvancedService) GetSystemLanguages(companyID *uuid.UUID, isActive *bool) ([]models.SystemLanguage, error) {
	return s.advancedRepo.GetSystemLanguages(companyID, isActive)
}

func (s *AdvancedService) GetTranslationsByLanguage(languageCode string, onlyApproved bool) (map[string]string, error) {
	language, err := s.advancedRepo.GetSystemLanguageByCode(languageCode)
	if err != nil {
		return nil, fmt.Errorf("language not found: %w", err)
	}

	translations, err := s.advancedRepo.GetTranslationsByLanguage(language.ID, onlyApproved)
	if err != nil {
		return nil, fmt.Errorf("failed to get translations: %w", err)
	}

	result := make(map[string]string)
	for _, translation := range translations {
		result[translation.TranslationKey] = translation.Translation
	}

	return result, nil
}

// Private helper methods
func (s *AdvancedService) simulateAIResponse(content string, assistant *models.AIAssistant) *AIMessageResponse {
	// This is a simulation - in reality, you would call OpenAI, Claude, etc.
	responseText := fmt.Sprintf("I understand you said: %s. How can I help you further?", content)
	
	return &AIMessageResponse{
		Content:      responseText,
		TokenCount:   len(strings.Split(content, " ")) + len(strings.Split(responseText, " ")),
		Cost:         0.002, // $0.002 per response
		ResponseTime: 1500,  // 1.5 seconds
	}
}

func (s *AdvancedService) processBatchOperation(operationID uuid.UUID) {
	operation, err := s.advancedRepo.GetBatchOperation(operationID)
	if err != nil {
		return
	}

	err = s.advancedRepo.StartBatchOperation(operationID)
	if err != nil {
		return
	}

	// Simulate batch processing
	for i := 0; i < operation.TotalItems; i++ {
		time.Sleep(100 * time.Millisecond) // Simulate processing time
		
		progress := int((float64(i+1) / float64(operation.TotalItems)) * 100)
		s.advancedRepo.UpdateBatchOperationProgress(operationID, progress, i+1, i+1, 0, "")
	}
}

func (s *AdvancedService) processBackup(backupID uuid.UUID) {
	// Simulate backup process
	time.Sleep(5 * time.Second)
	
	endTime := time.Now()
	s.advancedRepo.UpdateBackupRecord(&models.BackupRecord{
		ID:          backupID,
		Status:      "completed",
		FileSize:    1024 * 1024 * 100, // 100MB
		FilePath:    "/backups/backup_" + backupID.String() + ".sql.gz",
		Checksum:    "abc123def456",
		EndTime:     &endTime,
		Duration:    5,
		UpdatedAt:   time.Now(),
	})
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Request/Response types
type CreateAIAssistantRequest struct {
	Name             string                 `json:"name" validate:"required"`
	Type             string                 `json:"type" validate:"required"`
	Model            string                 `json:"model" validate:"required"`
	SystemPrompt     string                 `json:"system_prompt"`
	Temperature      float64                `json:"temperature"`
	MaxTokens        int                    `json:"max_tokens"`
	TopP             float64                `json:"top_p"`
	FrequencyPenalty float64                `json:"frequency_penalty"`
	PresencePenalty  float64                `json:"presence_penalty"`
	Configuration    map[string]interface{} `json:"configuration"`
}

type UpdateAIAssistantRequest struct {
	Name             *string                 `json:"name"`
	SystemPrompt     *string                 `json:"system_prompt"`
	Temperature      *float64                `json:"temperature"`
	MaxTokens        *int                    `json:"max_tokens"`
	IsActive         *bool                   `json:"is_active"`
	Configuration    *map[string]interface{} `json:"configuration"`
}

type AIMessageResponse struct {
	Content      string `json:"content"`
	TokenCount   int    `json:"token_count"`
	Cost         float64 `json:"cost"`
	ResponseTime int64  `json:"response_time"`
}

type CreateRecommendationRequest struct {
	Type         string                 `json:"type" validate:"required"`
	Category     string                 `json:"category"`
	Title        string                 `json:"title" validate:"required"`
	Description  string                 `json:"description"`
	Data         map[string]interface{} `json:"data"`
	Score        float64                `json:"score"`
	Priority     string                 `json:"priority" validate:"required"`
	Source       string                 `json:"source"`
	SourceData   map[string]interface{} `json:"source_data"`
	ResourceID   *uuid.UUID             `json:"resource_id"`
	ResourceType string                 `json:"resource_type"`
	ExpiresAt    *time.Time             `json:"expires_at"`
}

type CreateAdvancedSearchRequest struct {
	Name        string                 `json:"name" validate:"required"`
	Description string                 `json:"description"`
	SearchType  string                 `json:"search_type" validate:"required"`
	TableName   string                 `json:"table_name"`
	Filters     map[string]interface{} `json:"filters"`
	SortConfig  map[string]interface{} `json:"sort_config"`
	Columns     []string               `json:"columns"`
	IsPublic    bool                   `json:"is_public"`
}

type AdvancedSearchResult struct {
	SearchID      uuid.UUID                `json:"search_id"`
	ExecutedAt    time.Time                `json:"executed_at"`
	TotalCount    int                      `json:"total_count"`
	Results       []map[string]interface{} `json:"results"`
	ExecutionTime int64                    `json:"execution_time"`
}

type CreateBatchOperationRequest struct {
	OperationType string                 `json:"operation_type" validate:"required"`
	TargetTable   string                 `json:"target_table" validate:"required"`
	TargetIDs     []uuid.UUID            `json:"target_ids" validate:"required"`
	Parameters    map[string]interface{} `json:"parameters"`
}

type CreateCustomFieldRequest struct {
	TableName    string                 `json:"table_name" validate:"required"`
	FieldName    string                 `json:"field_name" validate:"required"`
	FieldLabel   string                 `json:"field_label" validate:"required"`
	FieldType    string                 `json:"field_type" validate:"required"`
	DefaultValue string                 `json:"default_value"`
	Options      []string               `json:"options"`
	Validation   map[string]interface{} `json:"validation"`
	IsRequired   bool                   `json:"is_required"`
	IsSearchable bool                   `json:"is_searchable"`
	DisplayOrder int                    `json:"display_order"`
}

type SetCustomFieldValueRequest struct {
	FieldID      uuid.UUID `json:"field_id" validate:"required"`
	ResourceID   uuid.UUID `json:"resource_id" validate:"required"`
	ResourceType string    `json:"resource_type" validate:"required"`
	Value        string    `json:"value"`
}

type CreateSecurityEventRequest struct {
	CompanyID     uuid.UUID              `json:"company_id" validate:"required"`
	UserID        *uuid.UUID             `json:"user_id"`
	EventType     string                 `json:"event_type" validate:"required"`
	Severity      string                 `json:"severity" validate:"required"`
	Description   string                 `json:"description"`
	IPAddress     string                 `json:"ip_address"`
	UserAgent     string                 `json:"user_agent"`
	Location      string                 `json:"location"`
	DeviceInfo    map[string]interface{} `json:"device_info"`
	ResourceType  string                 `json:"resource_type"`
	ResourceID    *uuid.UUID             `json:"resource_id"`
	ActionDetails map[string]interface{} `json:"action_details"`
	RiskScore     float64                `json:"risk_score"`
}

type RecordPerformanceMetricRequest struct {
	CompanyID  uuid.UUID              `json:"company_id" validate:"required"`
	MetricType string                 `json:"metric_type" validate:"required"`
	MetricName string                 `json:"metric_name" validate:"required"`
	Value      float64                `json:"value" validate:"required"`
	Unit       string                 `json:"unit"`
	Context    map[string]interface{} `json:"context"`
	Endpoint   string                 `json:"endpoint"`
	Method     string                 `json:"method"`
	StatusCode int                    `json:"status_code"`
	UserID     *uuid.UUID             `json:"user_id"`
	SessionID  string                 `json:"session_id"`
	TraceID    string                 `json:"trace_id"`
}

type CreateBackupRecordRequest struct {
	BackupType      string   `json:"backup_type" validate:"required"`
	Tables          []string `json:"tables"`
	CompressionType string   `json:"compression_type"`
	EncryptionType  string   `json:"encryption_type"`
}