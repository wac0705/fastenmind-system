package repositories

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"fastenmind-system/internal/models"
)

type AdvancedRepository struct {
	db *gorm.DB
}

func NewAdvancedRepository(db *gorm.DB) *AdvancedRepository {
	return &AdvancedRepository{db: db}
}

// AI Assistant Methods
func (r *AdvancedRepository) CreateAIAssistant(assistant *models.AIAssistant) error {
	return r.db.Create(assistant).Error
}

func (r *AdvancedRepository) GetAIAssistant(id uuid.UUID) (*models.AIAssistant, error) {
	var assistant models.AIAssistant
	err := r.db.Preload("Company").Preload("User").Preload("Creator").
		First(&assistant, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &assistant, nil
}

func (r *AdvancedRepository) GetAIAssistantsByCompany(companyID uuid.UUID, assistantType string, isActive *bool) ([]models.AIAssistant, error) {
	var assistants []models.AIAssistant
	query := r.db.Where("company_id = ?", companyID)
	
	if assistantType != "" {
		query = query.Where("type = ?", assistantType)
	}
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}
	
	err := query.Preload("Company").Preload("User").Preload("Creator").
		Order("created_at DESC").Find(&assistants).Error
	return assistants, err
}

func (r *AdvancedRepository) UpdateAIAssistant(assistant *models.AIAssistant) error {
	return r.db.Save(assistant).Error
}

func (r *AdvancedRepository) DeleteAIAssistant(id uuid.UUID) error {
	return r.db.Delete(&models.AIAssistant{}, "id = ?", id).Error
}

func (r *AdvancedRepository) IncrementAIAssistantUsage(id uuid.UUID, tokensUsed int64, cost float64) error {
	now := time.Now()
	return r.db.Model(&models.AIAssistant{}).Where("id = ?", id).Updates(map[string]interface{}{
		"usage_count":      gorm.Expr("usage_count + ?", 1),
		"tokens_used":      gorm.Expr("tokens_used + ?", tokensUsed),
		"cost_accumulated": gorm.Expr("cost_accumulated + ?", cost),
		"last_used":        now,
		"updated_at":       now,
	}).Error
}

// AI Conversation Session Methods
func (r *AdvancedRepository) CreateConversationSession(session *models.AIConversationSession) error {
	return r.db.Create(session).Error
}

func (r *AdvancedRepository) GetConversationSession(id uuid.UUID) (*models.AIConversationSession, error) {
	var session models.AIConversationSession
	err := r.db.Preload("Assistant").Preload("User").Preload("Company").Preload("Messages").
		First(&session, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *AdvancedRepository) GetConversationSessionsByUser(userID uuid.UUID, assistantID *uuid.UUID, status string, limit int) ([]models.AIConversationSession, error) {
	var sessions []models.AIConversationSession
	query := r.db.Where("user_id = ?", userID)
	
	if assistantID != nil {
		query = query.Where("assistant_id = ?", *assistantID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	err := query.Preload("Assistant").Preload("User").Preload("Company").
		Order("created_at DESC").Limit(limit).Find(&sessions).Error
	return sessions, err
}

func (r *AdvancedRepository) UpdateConversationSession(session *models.AIConversationSession) error {
	return r.db.Save(session).Error
}

func (r *AdvancedRepository) EndConversationSession(id uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&models.AIConversationSession{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":     "completed",
		"end_time":   now,
		"updated_at": now,
	}).Error
}

// AI Message Methods
func (r *AdvancedRepository) CreateAIMessage(message *models.AIMessage) error {
	return r.db.Create(message).Error
}

func (r *AdvancedRepository) GetMessagesBySession(sessionID uuid.UUID, limit int) ([]models.AIMessage, error) {
	var messages []models.AIMessage
	query := r.db.Where("session_id = ?", sessionID).Order("created_at ASC")
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	
	err := query.Find(&messages).Error
	return messages, err
}

func (r *AdvancedRepository) UpdateSessionStats(sessionID uuid.UUID, tokensUsed int64, cost float64) error {
	return r.db.Model(&models.AIConversationSession{}).Where("id = ?", sessionID).Updates(map[string]interface{}{
		"message_count": gorm.Expr("message_count + ?", 1),
		"tokens_used":   gorm.Expr("tokens_used + ?", tokensUsed),
		"cost":          gorm.Expr("cost + ?", cost),
		"updated_at":    time.Now(),
	}).Error
}

// Smart Recommendation Methods
func (r *AdvancedRepository) CreateRecommendation(recommendation *models.SmartRecommendation) error {
	return r.db.Create(recommendation).Error
}

func (r *AdvancedRepository) GetRecommendation(id uuid.UUID) (*models.SmartRecommendation, error) {
	var recommendation models.SmartRecommendation
	err := r.db.Preload("Company").Preload("User").First(&recommendation, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &recommendation, nil
}

func (r *AdvancedRepository) GetRecommendationsByUser(userID uuid.UUID, recommendationType string, status string, limit int) ([]models.SmartRecommendation, error) {
	var recommendations []models.SmartRecommendation
	query := r.db.Where("user_id = ?", userID)
	
	if recommendationType != "" {
		query = query.Where("type = ?", recommendationType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	err := query.Preload("Company").Preload("User").
		Order("priority DESC, score DESC, created_at DESC").Limit(limit).Find(&recommendations).Error
	return recommendations, err
}

func (r *AdvancedRepository) UpdateRecommendationStatus(id uuid.UUID, status string) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}
	
	if status == "viewed" {
		updates["viewed_at"] = time.Now()
	} else if status == "accepted" || status == "rejected" || status == "implemented" {
		updates["actioned_at"] = time.Now()
	}
	
	return r.db.Model(&models.SmartRecommendation{}).Where("id = ?", id).Updates(updates).Error
}

func (r *AdvancedRepository) CleanupExpiredRecommendations() error {
	return r.db.Delete(&models.SmartRecommendation{}, "expires_at < ?", time.Now()).Error
}

// Advanced Search Methods
func (r *AdvancedRepository) CreateAdvancedSearch(search *models.AdvancedSearch) error {
	return r.db.Create(search).Error
}

func (r *AdvancedRepository) GetAdvancedSearch(id uuid.UUID) (*models.AdvancedSearch, error) {
	var search models.AdvancedSearch
	err := r.db.Preload("Company").Preload("User").First(&search, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &search, nil
}

func (r *AdvancedRepository) GetAdvancedSearchesByUser(userID uuid.UUID, searchType string, isPublic *bool) ([]models.AdvancedSearch, error) {
	var searches []models.AdvancedSearch
	query := r.db.Where("user_id = ? OR is_public = ?", userID, true)
	
	if searchType != "" {
		query = query.Where("search_type = ?", searchType)
	}
	if isPublic != nil {
		query = query.Where("is_public = ?", *isPublic)
	}
	
	err := query.Preload("Company").Preload("User").
		Order("usage_count DESC, created_at DESC").Find(&searches).Error
	return searches, err
}

func (r *AdvancedRepository) UpdateAdvancedSearch(search *models.AdvancedSearch) error {
	return r.db.Save(search).Error
}

func (r *AdvancedRepository) IncrementSearchUsage(id uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&models.AdvancedSearch{}).Where("id = ?", id).Updates(map[string]interface{}{
		"usage_count": gorm.Expr("usage_count + ?", 1),
		"last_used":   now,
		"updated_at":  now,
	}).Error
}

func (r *AdvancedRepository) DeleteAdvancedSearch(id uuid.UUID) error {
	return r.db.Delete(&models.AdvancedSearch{}, "id = ?", id).Error
}

// Batch Operation Methods
func (r *AdvancedRepository) CreateBatchOperation(operation *models.BatchOperation) error {
	return r.db.Create(operation).Error
}

func (r *AdvancedRepository) GetBatchOperation(id uuid.UUID) (*models.BatchOperation, error) {
	var operation models.BatchOperation
	err := r.db.Preload("Company").Preload("User").First(&operation, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &operation, nil
}

func (r *AdvancedRepository) GetBatchOperationsByUser(userID uuid.UUID, status string, limit int) ([]models.BatchOperation, error) {
	var operations []models.BatchOperation
	query := r.db.Where("user_id = ?", userID)
	
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	err := query.Preload("Company").Preload("User").
		Order("created_at DESC").Limit(limit).Find(&operations).Error
	return operations, err
}

func (r *AdvancedRepository) UpdateBatchOperation(operation *models.BatchOperation) error {
	return r.db.Save(operation).Error
}

func (r *AdvancedRepository) UpdateBatchOperationProgress(id uuid.UUID, progress int, processedItems int, successCount int, errorCount int, errorLog string) error {
	updates := map[string]interface{}{
		"progress":        progress,
		"processed_items": processedItems,
		"success_count":   successCount,
		"error_count":     errorCount,
		"updated_at":      time.Now(),
	}
	
	if errorLog != "" {
		updates["error_log"] = errorLog
	}
	
	if progress >= 100 {
		updates["status"] = "completed"
		updates["completed_at"] = time.Now()
	}
	
	return r.db.Model(&models.BatchOperation{}).Where("id = ?", id).Updates(updates).Error
}

func (r *AdvancedRepository) StartBatchOperation(id uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&models.BatchOperation{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":     "running",
		"started_at": now,
		"updated_at": now,
	}).Error
}

// Custom Field Methods
func (r *AdvancedRepository) CreateCustomField(field *models.CustomField) error {
	return r.db.Create(field).Error
}

func (r *AdvancedRepository) GetCustomField(id uuid.UUID) (*models.CustomField, error) {
	var field models.CustomField
	err := r.db.Preload("Company").Preload("Creator").First(&field, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &field, nil
}

func (r *AdvancedRepository) GetCustomFieldsByTable(companyID uuid.UUID, tableName string, isActive *bool) ([]models.CustomField, error) {
	var fields []models.CustomField
	query := r.db.Where("company_id = ? AND table_name = ?", companyID, tableName)
	
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}
	
	err := query.Preload("Company").Preload("Creator").
		Order("display_order ASC, created_at ASC").Find(&fields).Error
	return fields, err
}

func (r *AdvancedRepository) UpdateCustomField(field *models.CustomField) error {
	return r.db.Save(field).Error
}

func (r *AdvancedRepository) DeleteCustomField(id uuid.UUID) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Delete(&models.CustomFieldValue{}, "field_id = ?", id).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(&models.CustomField{}, "id = ?", id).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// Custom Field Value Methods
func (r *AdvancedRepository) SetCustomFieldValue(value *models.CustomFieldValue) error {
	var existing models.CustomFieldValue
	err := r.db.Where("field_id = ? AND resource_id = ? AND resource_type = ?", 
		value.FieldID, value.ResourceID, value.ResourceType).First(&existing).Error
		
	if err == nil {
		existing.Value = value.Value
		existing.UpdatedAt = time.Now()
		return r.db.Save(&existing).Error
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return r.db.Create(value).Error
	}
	return err
}

func (r *AdvancedRepository) GetCustomFieldValues(resourceID uuid.UUID, resourceType string) ([]models.CustomFieldValue, error) {
	var values []models.CustomFieldValue
	err := r.db.Preload("Field").Where("resource_id = ? AND resource_type = ?", 
		resourceID, resourceType).Find(&values).Error
	return values, err
}

func (r *AdvancedRepository) GetCustomFieldValuesForMultipleResources(resourceIDs []uuid.UUID, resourceType string) (map[uuid.UUID][]models.CustomFieldValue, error) {
	var values []models.CustomFieldValue
	err := r.db.Preload("Field").Where("resource_id IN ? AND resource_type = ?", 
		resourceIDs, resourceType).Find(&values).Error
	if err != nil {
		return nil, err
	}
	
	result := make(map[uuid.UUID][]models.CustomFieldValue)
	for _, value := range values {
		result[value.ResourceID] = append(result[value.ResourceID], value)
	}
	return result, nil
}

// Security Event Methods
func (r *AdvancedRepository) CreateSecurityEvent(event *models.SecurityEvent) error {
	return r.db.Create(event).Error
}

func (r *AdvancedRepository) GetSecurityEvent(id uuid.UUID) (*models.SecurityEvent, error) {
	var event models.SecurityEvent
	err := r.db.Preload("Company").Preload("User").Preload("ResolvedByUser").
		First(&event, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *AdvancedRepository) GetSecurityEventsByCompany(companyID uuid.UUID, eventType string, severity string, status string, limit int) ([]models.SecurityEvent, error) {
	var events []models.SecurityEvent
	query := r.db.Where("company_id = ?", companyID)
	
	if eventType != "" {
		query = query.Where("event_type = ?", eventType)
	}
	if severity != "" {
		query = query.Where("severity = ?", severity)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	err := query.Preload("Company").Preload("User").Preload("ResolvedByUser").
		Order("risk_score DESC, created_at DESC").Limit(limit).Find(&events).Error
	return events, err
}

func (r *AdvancedRepository) UpdateSecurityEventStatus(id uuid.UUID, status string, resolvedBy *uuid.UUID) error {
	updates := map[string]interface{}{
		"status": status,
	}
	
	if status == "resolved" && resolvedBy != nil {
		updates["resolved_by"] = *resolvedBy
		updates["resolved_at"] = time.Now()
	}
	
	return r.db.Model(&models.SecurityEvent{}).Where("id = ?", id).Updates(updates).Error
}

// Performance Metric Methods
func (r *AdvancedRepository) CreatePerformanceMetric(metric *models.PerformanceMetric) error {
	return r.db.Create(metric).Error
}

func (r *AdvancedRepository) GetPerformanceMetrics(companyID uuid.UUID, metricType string, startTime time.Time, endTime time.Time, limit int) ([]models.PerformanceMetric, error) {
	var metrics []models.PerformanceMetric
	query := r.db.Where("company_id = ? AND created_at BETWEEN ? AND ?", companyID, startTime, endTime)
	
	if metricType != "" {
		query = query.Where("metric_type = ?", metricType)
	}
	
	err := query.Order("created_at DESC").Limit(limit).Find(&metrics).Error
	return metrics, err
}

func (r *AdvancedRepository) GetPerformanceStats(companyID uuid.UUID, metricType string, startTime time.Time, endTime time.Time) (map[string]interface{}, error) {
	var result struct {
		Avg   float64
		Min   float64
		Max   float64
		Count int64
	}
	
	query := r.db.Model(&models.PerformanceMetric{}).
		Where("company_id = ? AND metric_type = ? AND created_at BETWEEN ? AND ?", 
			companyID, metricType, startTime, endTime).
		Select("AVG(value) as avg, MIN(value) as min, MAX(value) as max, COUNT(*) as count")
	
	err := query.Scan(&result).Error
	if err != nil {
		return nil, err
	}
	
	return map[string]interface{}{
		"average": result.Avg,
		"minimum": result.Min,
		"maximum": result.Max,
		"count":   result.Count,
	}, nil
}

// Backup Record Methods
func (r *AdvancedRepository) CreateBackupRecord(record *models.BackupRecord) error {
	return r.db.Create(record).Error
}

func (r *AdvancedRepository) GetBackupRecord(id uuid.UUID) (*models.BackupRecord, error) {
	var record models.BackupRecord
	err := r.db.Preload("Company").Preload("Creator").First(&record, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *AdvancedRepository) GetBackupRecordsByCompany(companyID uuid.UUID, backupType string, status string, limit int) ([]models.BackupRecord, error) {
	var records []models.BackupRecord
	query := r.db.Where("company_id = ?", companyID)
	
	if backupType != "" {
		query = query.Where("backup_type = ?", backupType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	err := query.Preload("Company").Preload("Creator").
		Order("created_at DESC").Limit(limit).Find(&records).Error
	return records, err
}

func (r *AdvancedRepository) UpdateBackupRecord(record *models.BackupRecord) error {
	return r.db.Save(record).Error
}

// System Language Methods
func (r *AdvancedRepository) CreateSystemLanguage(language *models.SystemLanguage) error {
	return r.db.Create(language).Error
}

func (r *AdvancedRepository) GetSystemLanguage(id uuid.UUID) (*models.SystemLanguage, error) {
	var language models.SystemLanguage
	err := r.db.Preload("Company").Preload("Translations").First(&language, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &language, nil
}

func (r *AdvancedRepository) GetSystemLanguageByCode(code string) (*models.SystemLanguage, error) {
	var language models.SystemLanguage
	err := r.db.First(&language, "language_code = ?", code).Error
	if err != nil {
		return nil, err
	}
	return &language, nil
}

func (r *AdvancedRepository) GetSystemLanguages(companyID *uuid.UUID, isActive *bool) ([]models.SystemLanguage, error) {
	var languages []models.SystemLanguage
	query := r.db.Model(&models.SystemLanguage{})
	
	if companyID != nil {
		query = query.Where("company_id = ? OR company_id IS NULL", *companyID)
	}
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}
	
	err := query.Order("is_default DESC, language_name ASC").Find(&languages).Error
	return languages, err
}

func (r *AdvancedRepository) UpdateSystemLanguage(language *models.SystemLanguage) error {
	return r.db.Save(language).Error
}

// Translation Methods
func (r *AdvancedRepository) CreateTranslation(translation *models.Translation) error {
	return r.db.Create(translation).Error
}

func (r *AdvancedRepository) GetTranslation(languageID uuid.UUID, key string) (*models.Translation, error) {
	var translation models.Translation
	err := r.db.Where("language_id = ? AND translation_key = ?", languageID, key).
		First(&translation).Error
	if err != nil {
		return nil, err
	}
	return &translation, nil
}

func (r *AdvancedRepository) GetTranslationsByLanguage(languageID uuid.UUID, onlyApproved bool) ([]models.Translation, error) {
	var translations []models.Translation
	query := r.db.Where("language_id = ?", languageID)
	
	if onlyApproved {
		query = query.Where("is_approved = ?", true)
	}
	
	err := query.Order("translation_key ASC").Find(&translations).Error
	return translations, err
}

func (r *AdvancedRepository) UpdateTranslation(translation *models.Translation) error {
	return r.db.Save(translation).Error
}

func (r *AdvancedRepository) ApproveTranslation(id uuid.UUID, approvedBy uuid.UUID) error {
	return r.db.Model(&models.Translation{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_approved": true,
		"approved_by": approvedBy,
		"updated_at":  time.Now(),
	}).Error
}

func (r *AdvancedRepository) BulkUpsertTranslations(translations []models.Translation) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, translation := range translations {
		var existing models.Translation
		err := tx.Where("language_id = ? AND translation_key = ?", 
			translation.LanguageID, translation.TranslationKey).First(&existing).Error
			
		if err == nil {
			existing.Translation = translation.Translation
			existing.Context = translation.Context
			existing.UpdatedAt = time.Now()
			if err := tx.Save(&existing).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := tx.Create(&translation).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// Analytics Methods
func (r *AdvancedRepository) GetAIUsageStats(companyID uuid.UUID, startTime time.Time, endTime time.Time) (map[string]interface{}, error) {
	var stats struct {
		TotalSessions    int64
		TotalMessages    int64
		TotalTokensUsed  int64
		TotalCost       float64
		ActiveAssistants int64
	}
	
	// Get session stats
	r.db.Model(&models.AIConversationSession{}).
		Where("company_id = ? AND created_at BETWEEN ? AND ?", companyID, startTime, endTime).
		Count(&stats.TotalSessions)
	
	// Get message stats
	r.db.Model(&models.AIMessage{}).
		Joins("JOIN ai_conversation_sessions ON ai_messages.session_id = ai_conversation_sessions.id").
		Where("ai_conversation_sessions.company_id = ? AND ai_messages.created_at BETWEEN ? AND ?", 
			companyID, startTime, endTime).
		Count(&stats.TotalMessages)
	
	// Get token and cost stats
	r.db.Model(&models.AIConversationSession{}).
		Where("company_id = ? AND created_at BETWEEN ? AND ?", companyID, startTime, endTime).
		Select("SUM(tokens_used) as total_tokens_used, SUM(cost) as total_cost").
		Scan(&stats)
	
	// Get active assistants
	r.db.Model(&models.AIAssistant{}).
		Where("company_id = ? AND is_active = ? AND last_used BETWEEN ? AND ?", 
			companyID, true, startTime, endTime).
		Count(&stats.ActiveAssistants)
	
	return map[string]interface{}{
		"total_sessions":     stats.TotalSessions,
		"total_messages":     stats.TotalMessages,
		"total_tokens_used":  stats.TotalTokensUsed,
		"total_cost":         stats.TotalCost,
		"active_assistants":  stats.ActiveAssistants,
	}, nil
}

func (r *AdvancedRepository) GetRecommendationStats(companyID uuid.UUID, startTime time.Time, endTime time.Time) (map[string]interface{}, error) {
	var stats []struct {
		Type   string
		Status string
		Count  int64
	}
	
	err := r.db.Model(&models.SmartRecommendation{}).
		Where("company_id = ? AND created_at BETWEEN ? AND ?", companyID, startTime, endTime).
		Select("type, status, COUNT(*) as count").
		Group("type, status").
		Scan(&stats).Error
	
	if err != nil {
		return nil, err
	}
	
	result := make(map[string]interface{})
	for _, stat := range stats {
		key := fmt.Sprintf("%s_%s", stat.Type, stat.Status)
		result[key] = stat.Count
	}
	
	return result, nil
}

func (r *AdvancedRepository) GetSecurityEventStats(companyID uuid.UUID, startTime time.Time, endTime time.Time) (map[string]interface{}, error) {
	var stats []struct {
		EventType string
		Severity  string
		Count     int64
	}
	
	err := r.db.Model(&models.SecurityEvent{}).
		Where("company_id = ? AND created_at BETWEEN ? AND ?", companyID, startTime, endTime).
		Select("event_type, severity, COUNT(*) as count").
		Group("event_type, severity").
		Scan(&stats).Error
	
	if err != nil {
		return nil, err
	}
	
	result := make(map[string]interface{})
	for _, stat := range stats {
		key := fmt.Sprintf("%s_%s", stat.EventType, stat.Severity)
		result[key] = stat.Count
	}
	
	return result, nil
}