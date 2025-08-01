package repository

import (
	"time"

	"github.com/fastenmind/fastener-api/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MobileRepository interface {
	// Mobile Device operations
	RegisterDevice(device *models.MobileDevice) error
	UpdateDevice(device *models.MobileDevice) error
	GetDevice(id uuid.UUID) (*models.MobileDevice, error)
	GetDeviceByToken(deviceToken string) (*models.MobileDevice, error)
	ListUserDevices(userID uuid.UUID, params map[string]interface{}) ([]models.MobileDevice, error)
	ListCompanyDevices(companyID uuid.UUID, params map[string]interface{}) ([]models.MobileDevice, int64, error)
	DeactivateDevice(id uuid.UUID) error
	DeleteDevice(id uuid.UUID) error
	UpdateDeviceLastSeen(deviceToken string, lastSeen time.Time) error
	
	// Push Notification operations
	CreatePushNotification(notification *models.PushNotification) error
	UpdatePushNotification(notification *models.PushNotification) error
	GetPushNotification(id uuid.UUID) (*models.PushNotification, error)
	ListPushNotifications(params map[string]interface{}) ([]models.PushNotification, int64, error)
	ListDeviceNotifications(deviceID uuid.UUID, params map[string]interface{}) ([]models.PushNotification, int64, error)
	GetPendingNotifications(limit int) ([]models.PushNotification, error)
	MarkNotificationSent(id uuid.UUID, sentAt time.Time, providerResponse string) error
	MarkNotificationDelivered(id uuid.UUID, deliveredAt time.Time) error
	MarkNotificationClicked(id uuid.UUID, clickedAt time.Time) error
	MarkNotificationFailed(id uuid.UUID, errorMessage string) error
	DeleteExpiredNotifications() error
	
	// Mobile Session operations
	CreateMobileSession(session *models.MobileSession) error
	UpdateMobileSession(session *models.MobileSession) error
	GetMobileSession(sessionToken string) (*models.MobileSession, error)
	ListMobileSessions(params map[string]interface{}) ([]models.MobileSession, int64, error)
	EndMobileSession(sessionToken string, endTime time.Time) error
	CleanupExpiredSessions() error
	
	// Mobile Analytics operations
	CreateAnalyticsEvent(event *models.MobileAnalytics) error
	ListAnalyticsEvents(params map[string]interface{}) ([]models.MobileAnalytics, int64, error)
	GetAnalyticsSummary(companyID uuid.UUID, startDate, endDate time.Time) (map[string]interface{}, error)
	CleanupOldAnalytics(days int) error
	
	// Mobile App Version operations
	CreateAppVersion(version *models.MobileAppVersion) error
	UpdateAppVersion(version *models.MobileAppVersion) error
	GetAppVersion(id uuid.UUID) (*models.MobileAppVersion, error)
	GetLatestAppVersion(platform string, companyID *uuid.UUID) (*models.MobileAppVersion, error)
	ListAppVersions(params map[string]interface{}) ([]models.MobileAppVersion, int64, error)
	DeleteAppVersion(id uuid.UUID) error
	
	// Mobile Offline Data operations
	CreateOfflineData(data *models.MobileOfflineData) error
	UpdateOfflineData(data *models.MobileOfflineData) error
	GetOfflineData(id uuid.UUID) (*models.MobileOfflineData, error)
	ListPendingOfflineData(deviceID uuid.UUID, limit int) ([]models.MobileOfflineData, error)
	ListOfflineData(params map[string]interface{}) ([]models.MobileOfflineData, int64, error)
	MarkOfflineDataSynced(id uuid.UUID, syncedAt time.Time) error
	MarkOfflineDataFailed(id uuid.UUID, errorMessage string) error
	DeleteSyncedOfflineData(olderThan time.Time) error
	
	// Mobile Configuration operations
	CreateMobileConfig(config *models.MobileConfiguration) error
	UpdateMobileConfig(config *models.MobileConfiguration) error
	GetMobileConfig(id uuid.UUID) (*models.MobileConfiguration, error)
	GetMobileConfigByKey(key string, platform string, companyID *uuid.UUID) (*models.MobileConfiguration, error)
	ListMobileConfigs(params map[string]interface{}) ([]models.MobileConfiguration, error)
	DeleteMobileConfig(id uuid.UUID) error
	
	// Business operations
	GetMobileStatistics(companyID *uuid.UUID) (map[string]interface{}, error)
	GetDeviceUsageStats(deviceID uuid.UUID, days int) (map[string]interface{}, error)
	GetNotificationStats(companyID uuid.UUID, days int) (map[string]interface{}, error)
}

type mobileRepository struct {
	db *gorm.DB
}

func NewMobileRepository(db interface{}) MobileRepository {
	gormDB, ok := db.(*gorm.DB)
	if !ok {
		panic("invalid database type, expected *gorm.DB")
	}
	return &mobileRepository{db: gormDB}
}

// Mobile Device operations
func (r *mobileRepository) RegisterDevice(device *models.MobileDevice) error {
	return r.db.Create(device).Error
}

func (r *mobileRepository) UpdateDevice(device *models.MobileDevice) error {
	return r.db.Save(device).Error
}

func (r *mobileRepository) GetDevice(id uuid.UUID) (*models.MobileDevice, error) {
	var device models.MobileDevice
	err := r.db.Preload("User").
		Preload("Company").
		First(&device, id).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

func (r *mobileRepository) GetDeviceByToken(deviceToken string) (*models.MobileDevice, error) {
	var device models.MobileDevice
	err := r.db.Where("device_token = ? AND is_active = ?", deviceToken, true).
		Preload("User").
		Preload("Company").
		First(&device).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

func (r *mobileRepository) ListUserDevices(userID uuid.UUID, params map[string]interface{}) ([]models.MobileDevice, error) {
	var devices []models.MobileDevice
	
	query := r.db.Where("user_id = ?", userID)
	
	// Apply filters
	if platform, ok := params["platform"].(string); ok && platform != "" {
		query = query.Where("platform = ?", platform)
	}
	
	if isActive, ok := params["is_active"].(bool); ok {
		query = query.Where("is_active = ?", isActive)
	}
	
	err := query.Order("last_seen DESC").Find(&devices).Error
	return devices, err
}

func (r *mobileRepository) ListCompanyDevices(companyID uuid.UUID, params map[string]interface{}) ([]models.MobileDevice, int64, error) {
	var devices []models.MobileDevice
	var total int64
	
	query := r.db.Model(&models.MobileDevice{}).Where("company_id = ?", companyID)
	
	// Apply filters
	if platform, ok := params["platform"].(string); ok && platform != "" {
		query = query.Where("platform = ?", platform)
	}
	
	if isActive, ok := params["is_active"].(bool); ok {
		query = query.Where("is_active = ?", isActive)
	}
	
	if search, ok := params["search"].(string); ok && search != "" {
		query = query.Where("device_name LIKE ? OR device_model LIKE ?", 
			"%"+search+"%", "%"+search+"%")
	}
	
	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Apply pagination
	if page, ok := params["page"].(int); ok && page > 0 {
		pageSize := 20
		if ps, ok := params["page_size"].(int); ok && ps > 0 {
			pageSize = ps
		}
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}
	
	// Load with relations
	err := query.Order("last_seen DESC").
		Preload("User").
		Find(&devices).Error
	
	return devices, total, err
}

func (r *mobileRepository) DeactivateDevice(id uuid.UUID) error {
	return r.db.Model(&models.MobileDevice{}).
		Where("id = ?", id).
		Update("is_active", false).Error
}

func (r *mobileRepository) DeleteDevice(id uuid.UUID) error {
	return r.db.Delete(&models.MobileDevice{}, id).Error
}

func (r *mobileRepository) UpdateDeviceLastSeen(deviceToken string, lastSeen time.Time) error {
	return r.db.Model(&models.MobileDevice{}).
		Where("device_token = ?", deviceToken).
		Update("last_seen", lastSeen).Error
}

// Push Notification operations
func (r *mobileRepository) CreatePushNotification(notification *models.PushNotification) error {
	return r.db.Create(notification).Error
}

func (r *mobileRepository) UpdatePushNotification(notification *models.PushNotification) error {
	return r.db.Save(notification).Error
}

func (r *mobileRepository) GetPushNotification(id uuid.UUID) (*models.PushNotification, error) {
	var notification models.PushNotification
	err := r.db.Preload("Device").
		Preload("User").
		Preload("Company").
		First(&notification, id).Error
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

func (r *mobileRepository) ListPushNotifications(params map[string]interface{}) ([]models.PushNotification, int64, error) {
	var notifications []models.PushNotification
	var total int64
	
	query := r.db.Model(&models.PushNotification{})
	
	// Apply filters
	if companyID, ok := params["company_id"].(uuid.UUID); ok {
		query = query.Where("company_id = ?", companyID)
	}
	
	if userID, ok := params["user_id"].(uuid.UUID); ok {
		query = query.Where("user_id = ?", userID)
	}
	
	if deviceID, ok := params["device_id"].(uuid.UUID); ok {
		query = query.Where("device_id = ?", deviceID)
	}
	
	if status, ok := params["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	
	if notificationType, ok := params["type"].(string); ok && notificationType != "" {
		query = query.Where("type = ?", notificationType)
	}
	
	if startDate, ok := params["start_date"].(string); ok && startDate != "" {
		query = query.Where("created_at >= ?", startDate)
	}
	
	if endDate, ok := params["end_date"].(string); ok && endDate != "" {
		query = query.Where("created_at <= ?", endDate)
	}
	
	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Apply pagination
	if page, ok := params["page"].(int); ok && page > 0 {
		pageSize := 20
		if ps, ok := params["page_size"].(int); ok && ps > 0 {
			pageSize = ps
		}
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}
	
	// Load with relations
	err := query.Order("created_at DESC").
		Preload("Device").
		Preload("User").
		Find(&notifications).Error
	
	return notifications, total, err
}

func (r *mobileRepository) ListDeviceNotifications(deviceID uuid.UUID, params map[string]interface{}) ([]models.PushNotification, int64, error) {
	var notifications []models.PushNotification
	var total int64
	
	query := r.db.Model(&models.PushNotification{}).Where("device_id = ?", deviceID)
	
	// Apply filters
	if status, ok := params["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	
	if notificationType, ok := params["type"].(string); ok && notificationType != "" {
		query = query.Where("type = ?", notificationType)
	}
	
	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Apply pagination
	if page, ok := params["page"].(int); ok && page > 0 {
		pageSize := 20
		if ps, ok := params["page_size"].(int); ok && ps > 0 {
			pageSize = ps
		}
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}
	
	err := query.Order("created_at DESC").Find(&notifications).Error
	return notifications, total, err
}

func (r *mobileRepository) GetPendingNotifications(limit int) ([]models.PushNotification, error) {
	var notifications []models.PushNotification
	err := r.db.Where("status = ? AND (scheduled_at IS NULL OR scheduled_at <= ?)", 
		"pending", time.Now()).
		Order("priority DESC, created_at ASC").
		Limit(limit).
		Preload("Device").
		Find(&notifications).Error
	return notifications, err
}

func (r *mobileRepository) MarkNotificationSent(id uuid.UUID, sentAt time.Time, providerResponse string) error {
	return r.db.Model(&models.PushNotification{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":            "sent",
			"sent_at":           sentAt,
			"provider_response": providerResponse,
		}).Error
}

func (r *mobileRepository) MarkNotificationDelivered(id uuid.UUID, deliveredAt time.Time) error {
	return r.db.Model(&models.PushNotification{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       "delivered",
			"delivered_at": deliveredAt,
		}).Error
}

func (r *mobileRepository) MarkNotificationClicked(id uuid.UUID, clickedAt time.Time) error {
	return r.db.Model(&models.PushNotification{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     "clicked",
			"clicked_at": clickedAt,
		}).Error
}

func (r *mobileRepository) MarkNotificationFailed(id uuid.UUID, errorMessage string) error {
	return r.db.Model(&models.PushNotification{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":        "failed",
			"error_message": errorMessage,
		}).Error
}

func (r *mobileRepository) DeleteExpiredNotifications() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&models.PushNotification{}).Error
}

// Mobile Session operations
func (r *mobileRepository) CreateMobileSession(session *models.MobileSession) error {
	return r.db.Create(session).Error
}

func (r *mobileRepository) UpdateMobileSession(session *models.MobileSession) error {
	return r.db.Save(session).Error
}

func (r *mobileRepository) GetMobileSession(sessionToken string) (*models.MobileSession, error) {
	var session models.MobileSession
	err := r.db.Where("session_token = ?", sessionToken).
		Preload("Device").
		Preload("User").
		First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *mobileRepository) ListMobileSessions(params map[string]interface{}) ([]models.MobileSession, int64, error) {
	var sessions []models.MobileSession
	var total int64
	
	query := r.db.Model(&models.MobileSession{})
	
	// Apply filters
	if companyID, ok := params["company_id"].(uuid.UUID); ok {
		query = query.Where("company_id = ?", companyID)
	}
	
	if userID, ok := params["user_id"].(uuid.UUID); ok {
		query = query.Where("user_id = ?", userID)
	}
	
	if deviceID, ok := params["device_id"].(uuid.UUID); ok {
		query = query.Where("device_id = ?", deviceID)
	}
	
	if startDate, ok := params["start_date"].(string); ok && startDate != "" {
		query = query.Where("start_time >= ?", startDate)
	}
	
	if endDate, ok := params["end_date"].(string); ok && endDate != "" {
		query = query.Where("start_time <= ?", endDate)
	}
	
	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Apply pagination
	if page, ok := params["page"].(int); ok && page > 0 {
		pageSize := 20
		if ps, ok := params["page_size"].(int); ok && ps > 0 {
			pageSize = ps
		}
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}
	
	err := query.Order("start_time DESC").
		Preload("Device").
		Preload("User").
		Find(&sessions).Error
	
	return sessions, total, err
}

func (r *mobileRepository) EndMobileSession(sessionToken string, endTime time.Time) error {
	session, err := r.GetMobileSession(sessionToken)
	if err != nil {
		return err
	}
	
	duration := endTime.Sub(session.StartTime).Milliseconds()
	
	return r.db.Model(&models.MobileSession{}).
		Where("session_token = ?", sessionToken).
		Updates(map[string]interface{}{
			"end_time": endTime,
			"duration": duration,
		}).Error
}

func (r *mobileRepository) CleanupExpiredSessions() error {
	// Sessions older than 30 days
	cutoffDate := time.Now().AddDate(0, 0, -30)
	return r.db.Where("start_time < ?", cutoffDate).Delete(&models.MobileSession{}).Error
}

// Mobile Analytics operations
func (r *mobileRepository) CreateAnalyticsEvent(event *models.MobileAnalytics) error {
	return r.db.Create(event).Error
}

func (r *mobileRepository) ListAnalyticsEvents(params map[string]interface{}) ([]models.MobileAnalytics, int64, error) {
	var events []models.MobileAnalytics
	var total int64
	
	query := r.db.Model(&models.MobileAnalytics{})
	
	// Apply filters
	if companyID, ok := params["company_id"].(uuid.UUID); ok {
		query = query.Where("company_id = ?", companyID)
	}
	
	if userID, ok := params["user_id"].(uuid.UUID); ok {
		query = query.Where("user_id = ?", userID)
	}
	
	if deviceID, ok := params["device_id"].(uuid.UUID); ok {
		query = query.Where("device_id = ?", deviceID)
	}
	
	if eventType, ok := params["event_type"].(string); ok && eventType != "" {
		query = query.Where("event_type = ?", eventType)
	}
	
	if startDate, ok := params["start_date"].(string); ok && startDate != "" {
		query = query.Where("event_timestamp >= ?", startDate)
	}
	
	if endDate, ok := params["end_date"].(string); ok && endDate != "" {
		query = query.Where("event_timestamp <= ?", endDate)
	}
	
	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Apply pagination
	if page, ok := params["page"].(int); ok && page > 0 {
		pageSize := 50
		if ps, ok := params["page_size"].(int); ok && ps > 0 {
			pageSize = ps
		}
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}
	
	err := query.Order("event_timestamp DESC").Find(&events).Error
	return events, total, err
}

func (r *mobileRepository) GetAnalyticsSummary(companyID uuid.UUID, startDate, endDate time.Time) (map[string]interface{}, error) {
	summary := make(map[string]interface{})
	
	// Event counts by type
	var eventCounts []struct {
		EventType string `json:"event_type"`
		Count     int64  `json:"count"`
	}
	if err := r.db.Model(&models.MobileAnalytics{}).
		Where("company_id = ? AND event_timestamp BETWEEN ? AND ?", companyID, startDate, endDate).
		Select("event_type, COUNT(*) as count").
		Group("event_type").
		Find(&eventCounts).Error; err != nil {
		return nil, err
	}
	summary["event_counts"] = eventCounts
	
	// Active users
	var activeUsers int64
	if err := r.db.Model(&models.MobileAnalytics{}).
		Where("company_id = ? AND event_timestamp BETWEEN ? AND ?", companyID, startDate, endDate).
		Distinct("user_id").
		Count(&activeUsers).Error; err != nil {
		return nil, err
	}
	summary["active_users"] = activeUsers
	
	// Session statistics
	var sessionStats struct {
		TotalSessions   int64 `json:"total_sessions"`
		AvgDuration     int64 `json:"avg_duration"`
		TotalScreenViews int64 `json:"total_screen_views"`
	}
	if err := r.db.Model(&models.MobileSession{}).
		Where("company_id = ? AND start_time BETWEEN ? AND ?", companyID, startDate, endDate).
		Select("COUNT(*) as total_sessions, AVG(duration) as avg_duration, SUM(screen_views) as total_screen_views").
		Find(&sessionStats).Error; err != nil {
		return nil, err
	}
	summary["session_stats"] = sessionStats
	
	return summary, nil
}

func (r *mobileRepository) CleanupOldAnalytics(days int) error {
	cutoffDate := time.Now().AddDate(0, 0, -days)
	return r.db.Where("event_timestamp < ?", cutoffDate).Delete(&models.MobileAnalytics{}).Error
}

// Mobile App Version operations
func (r *mobileRepository) CreateAppVersion(version *models.MobileAppVersion) error {
	return r.db.Create(version).Error
}

func (r *mobileRepository) UpdateAppVersion(version *models.MobileAppVersion) error {
	return r.db.Save(version).Error
}

func (r *mobileRepository) GetAppVersion(id uuid.UUID) (*models.MobileAppVersion, error) {
	var version models.MobileAppVersion
	err := r.db.Preload("Company").
		Preload("Creator").
		First(&version, id).Error
	if err != nil {
		return nil, err
	}
	return &version, nil
}

func (r *mobileRepository) GetLatestAppVersion(platform string, companyID *uuid.UUID) (*models.MobileAppVersion, error) {
	var version models.MobileAppVersion
	query := r.db.Where("platform = ? AND status = ? AND is_active = ?", platform, "released", true)
	
	if companyID != nil {
		query = query.Where("company_id = ? OR company_id IS NULL", *companyID)
	} else {
		query = query.Where("company_id IS NULL")
	}
	
	err := query.Order("created_at DESC").First(&version).Error
	if err != nil {
		return nil, err
	}
	return &version, nil
}

func (r *mobileRepository) ListAppVersions(params map[string]interface{}) ([]models.MobileAppVersion, int64, error) {
	var versions []models.MobileAppVersion
	var total int64
	
	query := r.db.Model(&models.MobileAppVersion{})
	
	// Apply filters
	if companyID, ok := params["company_id"].(*uuid.UUID); ok && companyID != nil {
		query = query.Where("company_id = ? OR company_id IS NULL", *companyID)
	} else {
		query = query.Where("company_id IS NULL")
	}
	
	if platform, ok := params["platform"].(string); ok && platform != "" {
		query = query.Where("platform = ?", platform)
	}
	
	if status, ok := params["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	
	if releaseType, ok := params["release_type"].(string); ok && releaseType != "" {
		query = query.Where("release_type = ?", releaseType)
	}
	
	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Apply pagination
	if page, ok := params["page"].(int); ok && page > 0 {
		pageSize := 20
		if ps, ok := params["page_size"].(int); ok && ps > 0 {
			pageSize = ps
		}
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}
	
	err := query.Order("created_at DESC").
		Preload("Creator").
		Find(&versions).Error
	
	return versions, total, err
}

func (r *mobileRepository) DeleteAppVersion(id uuid.UUID) error {
	return r.db.Delete(&models.MobileAppVersion{}, id).Error
}

// Mobile Offline Data operations
func (r *mobileRepository) CreateOfflineData(data *models.MobileOfflineData) error {
	return r.db.Create(data).Error
}

func (r *mobileRepository) UpdateOfflineData(data *models.MobileOfflineData) error {
	return r.db.Save(data).Error
}

func (r *mobileRepository) GetOfflineData(id uuid.UUID) (*models.MobileOfflineData, error) {
	var data models.MobileOfflineData
	err := r.db.Preload("Device").
		Preload("User").
		Preload("ResolvedByUser").
		First(&data, id).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *mobileRepository) ListPendingOfflineData(deviceID uuid.UUID, limit int) ([]models.MobileOfflineData, error) {
	var data []models.MobileOfflineData
	err := r.db.Where("device_id = ? AND status IN ?", deviceID, []string{"pending", "failed"}).
		Order("priority DESC, created_at ASC").
		Limit(limit).
		Find(&data).Error
	return data, err
}

func (r *mobileRepository) ListOfflineData(params map[string]interface{}) ([]models.MobileOfflineData, int64, error) {
	var data []models.MobileOfflineData
	var total int64
	
	query := r.db.Model(&models.MobileOfflineData{})
	
	// Apply filters
	if companyID, ok := params["company_id"].(uuid.UUID); ok {
		query = query.Where("company_id = ?", companyID)
	}
	
	if deviceID, ok := params["device_id"].(uuid.UUID); ok {
		query = query.Where("device_id = ?", deviceID)
	}
	
	if status, ok := params["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	
	if dataType, ok := params["data_type"].(string); ok && dataType != "" {
		query = query.Where("data_type = ?", dataType)
	}
	
	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Apply pagination
	if page, ok := params["page"].(int); ok && page > 0 {
		pageSize := 20
		if ps, ok := params["page_size"].(int); ok && ps > 0 {
			pageSize = ps
		}
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}
	
	err := query.Order("created_at DESC").
		Preload("Device").
		Preload("User").
		Find(&data).Error
	
	return data, total, err
}

func (r *mobileRepository) MarkOfflineDataSynced(id uuid.UUID, syncedAt time.Time) error {
	return r.db.Model(&models.MobileOfflineData{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":    "synced",
			"synced_at": syncedAt,
		}).Error
}

func (r *mobileRepository) MarkOfflineDataFailed(id uuid.UUID, errorMessage string) error {
	return r.db.Model(&models.MobileOfflineData{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":        "failed",
			"error_message": errorMessage,
		}).Error
}

func (r *mobileRepository) DeleteSyncedOfflineData(olderThan time.Time) error {
	return r.db.Where("status = ? AND synced_at < ?", "synced", olderThan).
		Delete(&models.MobileOfflineData{}).Error
}

// Mobile Configuration operations
func (r *mobileRepository) CreateMobileConfig(config *models.MobileConfiguration) error {
	return r.db.Create(config).Error
}

func (r *mobileRepository) UpdateMobileConfig(config *models.MobileConfiguration) error {
	return r.db.Save(config).Error
}

func (r *mobileRepository) GetMobileConfig(id uuid.UUID) (*models.MobileConfiguration, error) {
	var config models.MobileConfiguration
	err := r.db.Preload("Company").
		Preload("Creator").
		First(&config, id).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *mobileRepository) GetMobileConfigByKey(key string, platform string, companyID *uuid.UUID) (*models.MobileConfiguration, error) {
	var config models.MobileConfiguration
	query := r.db.Where("config_key = ? AND (platform = ? OR platform = 'all') AND is_enabled = ?", 
		key, platform, true)
	
	if companyID != nil {
		query = query.Where("company_id = ? OR company_id IS NULL", *companyID).
			Order("company_id DESC") // Prioritize company-specific configs
	} else {
		query = query.Where("company_id IS NULL")
	}
	
	err := query.First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *mobileRepository) ListMobileConfigs(params map[string]interface{}) ([]models.MobileConfiguration, error) {
	var configs []models.MobileConfiguration
	
	query := r.db.Model(&models.MobileConfiguration{})
	
	// Apply filters
	if companyID, ok := params["company_id"].(*uuid.UUID); ok && companyID != nil {
		query = query.Where("company_id = ? OR company_id IS NULL", *companyID)
	} else {
		query = query.Where("company_id IS NULL")
	}
	
	if platform, ok := params["platform"].(string); ok && platform != "" {
		query = query.Where("platform = ? OR platform = 'all'", platform)
	}
	
	if category, ok := params["category"].(string); ok && category != "" {
		query = query.Where("category = ?", category)
	}
	
	if isEnabled, ok := params["is_enabled"].(bool); ok {
		query = query.Where("is_enabled = ?", isEnabled)
	}
	
	err := query.Order("category ASC, config_key ASC").
		Preload("Creator").
		Find(&configs).Error
	
	return configs, err
}

func (r *mobileRepository) DeleteMobileConfig(id uuid.UUID) error {
	return r.db.Delete(&models.MobileConfiguration{}, id).Error
}

// Business operations
func (r *mobileRepository) GetMobileStatistics(companyID *uuid.UUID) (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	
	// Device statistics
	var deviceStats []struct {
		Platform string `json:"platform"`
		Count    int64  `json:"count"`
	}
	deviceQuery := r.db.Model(&models.MobileDevice{}).Where("is_active = ?", true)
	if companyID != nil {
		deviceQuery = deviceQuery.Where("company_id = ?", *companyID)
	}
	if err := deviceQuery.Select("platform, COUNT(*) as count").
		Group("platform").
		Find(&deviceStats).Error; err != nil {
		return nil, err
	}
	stats["devices_by_platform"] = deviceStats
	
	// Push notification statistics
	var notificationStats []struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}
	notificationQuery := r.db.Model(&models.PushNotification{}).
		Where("created_at >= ?", time.Now().AddDate(0, 0, -30))
	if companyID != nil {
		notificationQuery = notificationQuery.Where("company_id = ?", *companyID)
	}
	if err := notificationQuery.Select("status, COUNT(*) as count").
		Group("status").
		Find(&notificationStats).Error; err != nil {
		return nil, err
	}
	stats["notifications_by_status"] = notificationStats
	
	// Total counts
	var totalDevices, activeDevices, totalNotifications int64
	
	deviceCountQuery := r.db.Model(&models.MobileDevice{})
	if companyID != nil {
		deviceCountQuery = deviceCountQuery.Where("company_id = ?", *companyID)
	}
	deviceCountQuery.Count(&totalDevices)
	
	activeDeviceQuery := r.db.Model(&models.MobileDevice{}).
		Where("is_active = ? AND last_seen >= ?", true, time.Now().AddDate(0, 0, -7))
	if companyID != nil {
		activeDeviceQuery = activeDeviceQuery.Where("company_id = ?", *companyID)
	}
	activeDeviceQuery.Count(&activeDevices)
	
	notificationCountQuery := r.db.Model(&models.PushNotification{}).
		Where("created_at >= ?", time.Now().AddDate(0, 0, -30))
	if companyID != nil {
		notificationCountQuery = notificationCountQuery.Where("company_id = ?", *companyID)
	}
	notificationCountQuery.Count(&totalNotifications)
	
	stats["total_devices"] = totalDevices
	stats["active_devices"] = activeDevices
	stats["total_notifications"] = totalNotifications
	
	return stats, nil
}

func (r *mobileRepository) GetDeviceUsageStats(deviceID uuid.UUID, days int) (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	startDate := time.Now().AddDate(0, 0, -days)
	
	// Session statistics
	var sessionStats struct {
		TotalSessions   int64 `json:"total_sessions"`
		AvgDuration     int64 `json:"avg_duration"`
		TotalScreenViews int64 `json:"total_screen_views"`
		TotalAPIRequests int64 `json:"total_api_requests"`
	}
	if err := r.db.Model(&models.MobileSession{}).
		Where("device_id = ? AND start_time >= ?", deviceID, startDate).
		Select("COUNT(*) as total_sessions, AVG(duration) as avg_duration, SUM(screen_views) as total_screen_views, SUM(api_requests) as total_api_requests").
		Find(&sessionStats).Error; err != nil {
		return nil, err
	}
	stats["session_stats"] = sessionStats
	
	// Analytics events
	var eventCounts []struct {
		EventType string `json:"event_type"`
		Count     int64  `json:"count"`
	}
	if err := r.db.Model(&models.MobileAnalytics{}).
		Where("device_id = ? AND event_timestamp >= ?", deviceID, startDate).
		Select("event_type, COUNT(*) as count").
		Group("event_type").
		Find(&eventCounts).Error; err != nil {
		return nil, err
	}
	stats["event_counts"] = eventCounts
	
	return stats, nil
}

func (r *mobileRepository) GetNotificationStats(companyID uuid.UUID, days int) (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	startDate := time.Now().AddDate(0, 0, -days)
	
	// Notification counts by type
	var typeCounts []struct {
		Type  string `json:"type"`
		Count int64  `json:"count"`
	}
	if err := r.db.Model(&models.PushNotification{}).
		Where("company_id = ? AND created_at >= ?", companyID, startDate).
		Select("type, COUNT(*) as count").
		Group("type").
		Find(&typeCounts).Error; err != nil {
		return nil, err
	}
	stats["notifications_by_type"] = typeCounts
	
	// Delivery rates
	var deliveryStats struct {
		TotalSent      int64   `json:"total_sent"`
		TotalDelivered int64   `json:"total_delivered"`
		TotalClicked   int64   `json:"total_clicked"`
		TotalFailed    int64   `json:"total_failed"`
		DeliveryRate   float64 `json:"delivery_rate"`
		ClickRate      float64 `json:"click_rate"`
	}
	
	r.db.Model(&models.PushNotification{}).
		Where("company_id = ? AND created_at >= ? AND status IN ?", 
			companyID, startDate, []string{"sent", "delivered", "clicked"}).
		Count(&deliveryStats.TotalSent)
	
	r.db.Model(&models.PushNotification{}).
		Where("company_id = ? AND created_at >= ? AND status IN ?", 
			companyID, startDate, []string{"delivered", "clicked"}).
		Count(&deliveryStats.TotalDelivered)
	
	r.db.Model(&models.PushNotification{}).
		Where("company_id = ? AND created_at >= ? AND status = ?", 
			companyID, startDate, "clicked").
		Count(&deliveryStats.TotalClicked)
	
	r.db.Model(&models.PushNotification{}).
		Where("company_id = ? AND created_at >= ? AND status = ?", 
			companyID, startDate, "failed").
		Count(&deliveryStats.TotalFailed)
	
	if deliveryStats.TotalSent > 0 {
		deliveryStats.DeliveryRate = float64(deliveryStats.TotalDelivered) / float64(deliveryStats.TotalSent) * 100
		deliveryStats.ClickRate = float64(deliveryStats.TotalClicked) / float64(deliveryStats.TotalSent) * 100
	}
	
	stats["delivery_stats"] = deliveryStats
	
	return stats, nil
}