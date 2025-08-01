package repositories

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/fastenmind/fastener-api/internal/models"
)

type IntegrationRepository struct {
	db *gorm.DB
}

func NewIntegrationRepository(db *gorm.DB) *IntegrationRepository {
	return &IntegrationRepository{db: db}
}

// Integration Methods
func (r *IntegrationRepository) CreateIntegration(integration *models.Integration) error {
	return r.db.Create(integration).Error
}

func (r *IntegrationRepository) GetIntegration(id uuid.UUID) (*models.Integration, error) {
	var integration models.Integration
	err := r.db.Preload("Company").Preload("Creator").Preload("Mappings").Preload("Webhooks").
		First(&integration, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &integration, nil
}

func (r *IntegrationRepository) GetIntegrationsByCompany(companyID uuid.UUID, integrationType string, status string, isActive *bool) ([]models.Integration, error) {
	var integrations []models.Integration
	query := r.db.Where("company_id = ?", companyID)
	
	if integrationType != "" {
		query = query.Where("type = ?", integrationType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}
	
	err := query.Preload("Company").Preload("Creator").
		Order("created_at DESC").Find(&integrations).Error
	return integrations, err
}

func (r *IntegrationRepository) UpdateIntegration(integration *models.Integration) error {
	return r.db.Save(integration).Error
}

func (r *IntegrationRepository) DeleteIntegration(id uuid.UUID) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete related records
	if err := tx.Delete(&models.IntegrationMapping{}, "integration_id = ?", id).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Delete(&models.Webhook{}, "integration_id = ?", id).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Delete(&models.DataSyncJob{}, "integration_id = ?", id).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Delete(&models.IntegrationLog{}, "integration_id = ?", id).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(&models.Integration{}, "id = ?", id).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *IntegrationRepository) UpdateIntegrationStats(id uuid.UUID, isSuccess bool, responseTime int64) error {
	updates := map[string]interface{}{
		"sync_count":     gorm.Expr("sync_count + ?", 1),
		"last_sync_at":   time.Now(),
		"updated_at":     time.Now(),
	}

	if isSuccess {
		updates["last_error"] = ""
		updates["last_error_at"] = nil
	} else {
		updates["error_count"] = gorm.Expr("error_count + ?", 1)
		updates["last_error_at"] = time.Now()
	}

	if responseTime > 0 {
		// Calculate running average
		updates["avg_response_time"] = gorm.Expr(
			"(avg_response_time * sync_count + ?) / (sync_count + 1)",
			responseTime,
		)
	}

	return r.db.Model(&models.Integration{}).Where("id = ?", id).Updates(updates).Error
}

func (r *IntegrationRepository) UpdateIntegrationSuccessRate(id uuid.UUID) error {
	return r.db.Model(&models.Integration{}).Where("id = ?", id).Update(
		"success_rate", 
		gorm.Expr("CASE WHEN sync_count > 0 THEN ((sync_count - error_count) * 100.0 / sync_count) ELSE 0.0 END"),
	).Error
}

// Integration Mapping Methods
func (r *IntegrationRepository) CreateIntegrationMapping(mapping *models.IntegrationMapping) error {
	return r.db.Create(mapping).Error
}

func (r *IntegrationRepository) GetIntegrationMapping(id uuid.UUID) (*models.IntegrationMapping, error) {
	var mapping models.IntegrationMapping
	err := r.db.Preload("Company").Preload("Integration").Preload("Creator").
		First(&mapping, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &mapping, nil
}

func (r *IntegrationRepository) GetMappingsByIntegration(integrationID uuid.UUID, direction string, isActive *bool) ([]models.IntegrationMapping, error) {
	var mappings []models.IntegrationMapping
	query := r.db.Where("integration_id = ?", integrationID)
	
	if direction != "" {
		query = query.Where("direction = ?", direction)
	}
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}
	
	err := query.Preload("Company").Preload("Integration").Preload("Creator").
		Order("created_at DESC").Find(&mappings).Error
	return mappings, err
}

func (r *IntegrationRepository) UpdateIntegrationMapping(mapping *models.IntegrationMapping) error {
	return r.db.Save(mapping).Error
}

func (r *IntegrationRepository) DeleteIntegrationMapping(id uuid.UUID) error {
	return r.db.Delete(&models.IntegrationMapping{}, "id = ?", id).Error
}

func (r *IntegrationRepository) UpdateMappingLastSync(id uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&models.IntegrationMapping{}).Where("id = ?", id).Updates(map[string]interface{}{
		"last_sync_at": now,
		"updated_at":   now,
	}).Error
}

// Webhook Methods
func (r *IntegrationRepository) CreateWebhook(webhook *models.Webhook) error {
	return r.db.Create(webhook).Error
}

func (r *IntegrationRepository) GetWebhook(id uuid.UUID) (*models.Webhook, error) {
	var webhook models.Webhook
	err := r.db.Preload("Company").Preload("Integration").Preload("Creator").
		First(&webhook, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &webhook, nil
}

func (r *IntegrationRepository) GetWebhooksByCompany(companyID uuid.UUID, integrationID *uuid.UUID, isActive *bool) ([]models.Webhook, error) {
	var webhooks []models.Webhook
	query := r.db.Where("company_id = ?", companyID)
	
	if integrationID != nil {
		query = query.Where("integration_id = ?", *integrationID)
	}
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}
	
	err := query.Preload("Company").Preload("Integration").Preload("Creator").
		Order("created_at DESC").Find(&webhooks).Error
	return webhooks, err
}

func (r *IntegrationRepository) UpdateWebhook(webhook *models.Webhook) error {
	return r.db.Save(webhook).Error
}

func (r *IntegrationRepository) DeleteWebhook(id uuid.UUID) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Delete(&models.WebhookDelivery{}, "webhook_id = ?", id).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(&models.Webhook{}, "id = ?", id).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *IntegrationRepository) UpdateWebhookStats(id uuid.UUID, isSuccess bool) error {
	now := time.Now()
	updates := map[string]interface{}{
		"trigger_count":      gorm.Expr("trigger_count + ?", 1),
		"last_triggered_at":  now,
		"updated_at":         now,
	}

	if isSuccess {
		updates["success_count"] = gorm.Expr("success_count + ?", 1)
	} else {
		updates["failure_count"] = gorm.Expr("failure_count + ?", 1)
	}

	return r.db.Model(&models.Webhook{}).Where("id = ?", id).Updates(updates).Error
}

// Webhook Delivery Methods
func (r *IntegrationRepository) CreateWebhookDelivery(delivery *models.WebhookDelivery) error {
	return r.db.Create(delivery).Error
}

func (r *IntegrationRepository) GetWebhookDelivery(id uuid.UUID) (*models.WebhookDelivery, error) {
	var delivery models.WebhookDelivery
	err := r.db.Preload("Webhook").Preload("Company").First(&delivery, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &delivery, nil
}

func (r *IntegrationRepository) GetWebhookDeliveries(webhookID uuid.UUID, status string, limit int) ([]models.WebhookDelivery, error) {
	var deliveries []models.WebhookDelivery
	query := r.db.Where("webhook_id = ?", webhookID)
	
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	err := query.Preload("Webhook").Preload("Company").
		Order("created_at DESC").Limit(limit).Find(&deliveries).Error
	return deliveries, err
}

func (r *IntegrationRepository) UpdateWebhookDelivery(delivery *models.WebhookDelivery) error {
	return r.db.Save(delivery).Error
}

func (r *IntegrationRepository) GetPendingWebhookDeliveries(limit int) ([]models.WebhookDelivery, error) {
	var deliveries []models.WebhookDelivery
	err := r.db.Where("status IN ? AND (next_retry_at IS NULL OR next_retry_at <= ?)", 
		[]string{"pending", "retrying"}, time.Now()).
		Preload("Webhook").Preload("Company").
		Order("created_at ASC").Limit(limit).Find(&deliveries).Error
	return deliveries, err
}

// Data Sync Job Methods
func (r *IntegrationRepository) CreateDataSyncJob(job *models.DataSyncJob) error {
	return r.db.Create(job).Error
}

func (r *IntegrationRepository) GetDataSyncJob(id uuid.UUID) (*models.DataSyncJob, error) {
	var job models.DataSyncJob
	err := r.db.Preload("Company").Preload("Integration").Preload("Mapping").Preload("Creator").
		First(&job, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *IntegrationRepository) GetDataSyncJobsByIntegration(integrationID uuid.UUID, status string, limit int) ([]models.DataSyncJob, error) {
	var jobs []models.DataSyncJob
	query := r.db.Where("integration_id = ?", integrationID)
	
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	err := query.Preload("Company").Preload("Integration").Preload("Mapping").Preload("Creator").
		Order("created_at DESC").Limit(limit).Find(&jobs).Error
	return jobs, err
}

func (r *IntegrationRepository) UpdateDataSyncJob(job *models.DataSyncJob) error {
	return r.db.Save(job).Error
}

func (r *IntegrationRepository) UpdateSyncJobProgress(id uuid.UUID, progress int, processedRecords int64, successRecords int64, errorRecords int64) error {
	updates := map[string]interface{}{
		"progress":         progress,
		"processed_records": processedRecords,
		"success_records":  successRecords,
		"error_records":    errorRecords,
		"updated_at":       time.Now(),
	}
	
	if progress >= 100 {
		updates["status"] = "completed"
		updates["completed_at"] = time.Now()
	}
	
	return r.db.Model(&models.DataSyncJob{}).Where("id = ?", id).Updates(updates).Error
}

func (r *IntegrationRepository) StartSyncJob(id uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&models.DataSyncJob{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":     "running",
		"started_at": now,
		"updated_at": now,
	}).Error
}

func (r *IntegrationRepository) GetPendingSyncJobs(limit int) ([]models.DataSyncJob, error) {
	var jobs []models.DataSyncJob
	err := r.db.Where("status = ? AND (scheduled_at IS NULL OR scheduled_at <= ?)", 
		"pending", time.Now()).
		Preload("Company").Preload("Integration").Preload("Mapping").
		Order("priority DESC, created_at ASC").Limit(limit).Find(&jobs).Error
	return jobs, err
}

// Integration Log Methods
func (r *IntegrationRepository) CreateIntegrationLog(log *models.IntegrationLog) error {
	return r.db.Create(log).Error
}

func (r *IntegrationRepository) GetIntegrationLogs(integrationID uuid.UUID, level string, category string, startTime time.Time, endTime time.Time, limit int) ([]models.IntegrationLog, error) {
	var logs []models.IntegrationLog
	query := r.db.Where("integration_id = ? AND created_at BETWEEN ? AND ?", 
		integrationID, startTime, endTime)
	
	if level != "" {
		query = query.Where("level = ?", level)
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}
	
	err := query.Preload("Company").Preload("Integration").Preload("SyncJob").Preload("User").
		Order("created_at DESC").Limit(limit).Find(&logs).Error
	return logs, err
}

func (r *IntegrationRepository) GetIntegrationLogsByCompany(companyID uuid.UUID, level string, startTime time.Time, endTime time.Time, limit int) ([]models.IntegrationLog, error) {
	var logs []models.IntegrationLog
	query := r.db.Where("company_id = ? AND created_at BETWEEN ? AND ?", 
		companyID, startTime, endTime)
	
	if level != "" {
		query = query.Where("level = ?", level)
	}
	
	err := query.Preload("Company").Preload("Integration").Preload("SyncJob").Preload("User").
		Order("created_at DESC").Limit(limit).Find(&logs).Error
	return logs, err
}

func (r *IntegrationRepository) CleanupOldLogs(retentionDays int) error {
	cutoffDate := time.Now().AddDate(0, 0, -retentionDays)
	return r.db.Delete(&models.IntegrationLog{}, "created_at < ?", cutoffDate).Error
}

// API Key Methods
func (r *IntegrationRepository) CreateApiKey(apiKey *models.ApiKey) error {
	return r.db.Create(apiKey).Error
}

func (r *IntegrationRepository) GetApiKey(id uuid.UUID) (*models.ApiKey, error) {
	var apiKey models.ApiKey
	err := r.db.Preload("Company").Preload("User").Preload("Creator").
		First(&apiKey, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &apiKey, nil
}

func (r *IntegrationRepository) GetApiKeyByHash(keyHash string) (*models.ApiKey, error) {
	var apiKey models.ApiKey
	err := r.db.Preload("Company").Preload("User").
		Where("key_hash = ? AND is_active = ?", keyHash, true).
		First(&apiKey).Error
	if err != nil {
		return nil, err
	}
	
	// Check expiration
	if apiKey.ExpiresAt != nil && apiKey.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("API key has expired")
	}
	
	return &apiKey, nil
}

func (r *IntegrationRepository) GetApiKeysByCompany(companyID uuid.UUID, userID *uuid.UUID, isActive *bool) ([]models.ApiKey, error) {
	var apiKeys []models.ApiKey
	query := r.db.Where("company_id = ?", companyID)
	
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}
	
	err := query.Preload("Company").Preload("User").Preload("Creator").
		Order("created_at DESC").Find(&apiKeys).Error
	return apiKeys, err
}

func (r *IntegrationRepository) UpdateApiKey(apiKey *models.ApiKey) error {
	return r.db.Save(apiKey).Error
}

func (r *IntegrationRepository) UpdateApiKeyUsage(keyHash string) error {
	now := time.Now()
	return r.db.Model(&models.ApiKey{}).Where("key_hash = ?", keyHash).Updates(map[string]interface{}{
		"usage_count":  gorm.Expr("usage_count + ?", 1),
		"last_used_at": now,
		"updated_at":   now,
	}).Error
}

func (r *IntegrationRepository) DeleteApiKey(id uuid.UUID) error {
	return r.db.Delete(&models.ApiKey{}, "id = ?", id).Error
}

// External System Methods
func (r *IntegrationRepository) CreateExternalSystem(system *models.ExternalSystem) error {
	return r.db.Create(system).Error
}

func (r *IntegrationRepository) GetExternalSystem(id uuid.UUID) (*models.ExternalSystem, error) {
	var system models.ExternalSystem
	err := r.db.Preload("Company").Preload("Creator").
		First(&system, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &system, nil
}

func (r *IntegrationRepository) GetExternalSystemsByCompany(companyID uuid.UUID, systemType string, status string) ([]models.ExternalSystem, error) {
	var systems []models.ExternalSystem
	query := r.db.Where("company_id = ?", companyID)
	
	if systemType != "" {
		query = query.Where("system_type = ?", systemType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	err := query.Preload("Company").Preload("Creator").
		Order("created_at DESC").Find(&systems).Error
	return systems, err
}

func (r *IntegrationRepository) UpdateExternalSystem(system *models.ExternalSystem) error {
	return r.db.Save(system).Error
}

func (r *IntegrationRepository) DeleteExternalSystem(id uuid.UUID) error {
	return r.db.Delete(&models.ExternalSystem{}, "id = ?", id).Error
}

// Data Transformation Methods
func (r *IntegrationRepository) CreateDataTransformation(transformation *models.DataTransformation) error {
	return r.db.Create(transformation).Error
}

func (r *IntegrationRepository) GetDataTransformation(id uuid.UUID) (*models.DataTransformation, error) {
	var transformation models.DataTransformation
	err := r.db.Preload("Company").Preload("Mapping").Preload("Creator").
		First(&transformation, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &transformation, nil
}

func (r *IntegrationRepository) GetTransformationsByMapping(mappingID uuid.UUID, isActive *bool) ([]models.DataTransformation, error) {
	var transformations []models.DataTransformation
	query := r.db.Where("mapping_id = ?", mappingID)
	
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}
	
	err := query.Preload("Company").Preload("Mapping").Preload("Creator").
		Order("execution_order ASC, created_at ASC").Find(&transformations).Error
	return transformations, err
}

func (r *IntegrationRepository) UpdateDataTransformation(transformation *models.DataTransformation) error {
	return r.db.Save(transformation).Error
}

func (r *IntegrationRepository) DeleteDataTransformation(id uuid.UUID) error {
	return r.db.Delete(&models.DataTransformation{}, "id = ?", id).Error
}

// Integration Template Methods
func (r *IntegrationRepository) CreateIntegrationTemplate(template *models.IntegrationTemplate) error {
	return r.db.Create(template).Error
}

func (r *IntegrationRepository) GetIntegrationTemplate(id uuid.UUID) (*models.IntegrationTemplate, error) {
	var template models.IntegrationTemplate
	err := r.db.Preload("Company").Preload("Creator").First(&template, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *IntegrationRepository) GetIntegrationTemplates(companyID *uuid.UUID, category string, provider string, isPublic *bool, isActive *bool) ([]models.IntegrationTemplate, error) {
	var templates []models.IntegrationTemplate
	query := r.db.Model(&models.IntegrationTemplate{})
	
	if companyID != nil {
		query = query.Where("company_id = ? OR is_public = ?", *companyID, true)
	} else {
		query = query.Where("is_public = ?", true)
	}
	
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if provider != "" {
		query = query.Where("provider = ?", provider)
	}
	if isPublic != nil {
		query = query.Where("is_public = ?", *isPublic)
	}
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}
	
	err := query.Preload("Company").Preload("Creator").
		Order("usage_count DESC, rating DESC, created_at DESC").Find(&templates).Error
	return templates, err
}

func (r *IntegrationRepository) UpdateIntegrationTemplate(template *models.IntegrationTemplate) error {
	return r.db.Save(template).Error
}

func (r *IntegrationRepository) IncrementTemplateUsage(id uuid.UUID) error {
	return r.db.Model(&models.IntegrationTemplate{}).Where("id = ?", id).Update(
		"usage_count", gorm.Expr("usage_count + ?", 1)).Error
}

func (r *IntegrationRepository) DeleteIntegrationTemplate(id uuid.UUID) error {
	return r.db.Delete(&models.IntegrationTemplate{}, "id = ?", id).Error
}

// Analytics Methods
func (r *IntegrationRepository) GetIntegrationStats(companyID uuid.UUID, startTime time.Time, endTime time.Time) (map[string]interface{}, error) {
	var stats struct {
		TotalIntegrations   int64
		ActiveIntegrations  int64
		TotalSyncJobs      int64
		SuccessfulSyncs    int64
		FailedSyncs        int64
		TotalWebhooks      int64
		WebhookDeliveries  int64
		AvgResponseTime    float64
		AvgSuccessRate     float64
	}
	
	// Get integration counts
	r.db.Model(&models.Integration{}).Where("company_id = ?", companyID).Count(&stats.TotalIntegrations)
	r.db.Model(&models.Integration{}).Where("company_id = ? AND is_active = ?", companyID, true).Count(&stats.ActiveIntegrations)
	
	// Get sync job stats
	r.db.Model(&models.DataSyncJob{}).Where("company_id = ? AND created_at BETWEEN ? AND ?", 
		companyID, startTime, endTime).Count(&stats.TotalSyncJobs)
	r.db.Model(&models.DataSyncJob{}).Where("company_id = ? AND status = ? AND created_at BETWEEN ? AND ?", 
		companyID, "completed", startTime, endTime).Count(&stats.SuccessfulSyncs)
	r.db.Model(&models.DataSyncJob{}).Where("company_id = ? AND status = ? AND created_at BETWEEN ? AND ?", 
		companyID, "failed", startTime, endTime).Count(&stats.FailedSyncs)
	
	// Get webhook stats
	r.db.Model(&models.Webhook{}).Where("company_id = ?", companyID).Count(&stats.TotalWebhooks)
	r.db.Model(&models.WebhookDelivery{}).Where("company_id = ? AND created_at BETWEEN ? AND ?", 
		companyID, startTime, endTime).Count(&stats.WebhookDeliveries)
	
	// Get average metrics
	r.db.Model(&models.Integration{}).Where("company_id = ? AND avg_response_time > 0", companyID).
		Select("AVG(avg_response_time)").Scan(&stats.AvgResponseTime)
	r.db.Model(&models.Integration{}).Where("company_id = ? AND success_rate > 0", companyID).
		Select("AVG(success_rate)").Scan(&stats.AvgSuccessRate)
	
	return map[string]interface{}{
		"total_integrations":   stats.TotalIntegrations,
		"active_integrations":  stats.ActiveIntegrations,
		"total_sync_jobs":      stats.TotalSyncJobs,
		"successful_syncs":     stats.SuccessfulSyncs,
		"failed_syncs":         stats.FailedSyncs,
		"total_webhooks":       stats.TotalWebhooks,
		"webhook_deliveries":   stats.WebhookDeliveries,
		"avg_response_time":    stats.AvgResponseTime,
		"avg_success_rate":     stats.AvgSuccessRate,
	}, nil
}

func (r *IntegrationRepository) GetIntegrationsByType(companyID uuid.UUID) ([]map[string]interface{}, error) {
	var results []struct {
		Type  string
		Count int64
	}
	
	err := r.db.Model(&models.Integration{}).
		Where("company_id = ?", companyID).
		Select("type, COUNT(*) as count").
		Group("type").
		Scan(&results).Error
	
	if err != nil {
		return nil, err
	}
	
	var data []map[string]interface{}
	for _, result := range results {
		data = append(data, map[string]interface{}{
			"type":  result.Type,
			"count": result.Count,
		})
	}
	
	return data, nil
}

func (r *IntegrationRepository) GetSyncJobTrends(companyID uuid.UUID, days int) ([]map[string]interface{}, error) {
	var results []struct {
		Date         string
		SuccessCount int64
		FailureCount int64
	}
	
	err := r.db.Raw(`
		SELECT 
			DATE(created_at) as date,
			SUM(CASE WHEN status = 'completed' THEN 1 ELSE 0 END) as success_count,
			SUM(CASE WHEN status = 'failed' THEN 1 ELSE 0 END) as failure_count
		FROM data_sync_jobs 
		WHERE company_id = ? AND created_at >= ?
		GROUP BY DATE(created_at)
		ORDER BY date
	`, companyID, time.Now().AddDate(0, 0, -days)).Scan(&results).Error
	
	if err != nil {
		return nil, err
	}
	
	var data []map[string]interface{}
	for _, result := range results {
		data = append(data, map[string]interface{}{
			"date":          result.Date,
			"success_count": result.SuccessCount,
			"failure_count": result.FailureCount,
		})
	}
	
	return data, nil
}

func (r *IntegrationRepository) GetTopErrorsByIntegration(companyID uuid.UUID, limit int) ([]map[string]interface{}, error) {
	var results []struct {
		IntegrationName string
		ErrorCount      int64
		LastError       string
	}
	
	err := r.db.Raw(`
		SELECT 
			i.name as integration_name,
			COUNT(il.id) as error_count,
			il.error_message as last_error
		FROM integration_logs il
		JOIN integrations i ON il.integration_id = i.id
		WHERE il.company_id = ? AND il.level = 'error'
		GROUP BY i.name, il.error_message
		ORDER BY error_count DESC
		LIMIT ?
	`, companyID, limit).Scan(&results).Error
	
	if err != nil {
		return nil, err
	}
	
	var data []map[string]interface{}
	for _, result := range results {
		data = append(data, map[string]interface{}{
			"integration_name": result.IntegrationName,
			"error_count":      result.ErrorCount,
			"last_error":       result.LastError,
		})
	}
	
	return data, nil
}