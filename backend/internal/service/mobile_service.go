package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/fastenmind/fastener-api/internal/models"
	"github.com/fastenmind/fastener-api/internal/repository"
	"github.com/google/uuid"
)

type MobileService interface {
	// Device Management
	RegisterDevice(ctx context.Context, device *models.MobileDevice) error
	UpdateDevice(ctx context.Context, device *models.MobileDevice) error
	GetDevice(ctx context.Context, id uuid.UUID) (*models.MobileDevice, error)
	GetDeviceByToken(ctx context.Context, deviceToken string) (*models.MobileDevice, error)
	ListUserDevices(ctx context.Context, userID uuid.UUID, params map[string]interface{}) ([]models.MobileDevice, error)
	ListCompanyDevices(ctx context.Context, companyID uuid.UUID, params map[string]interface{}) ([]models.MobileDevice, int64, error)
	DeactivateDevice(ctx context.Context, id uuid.UUID) error
	UpdateDeviceLastSeen(ctx context.Context, deviceToken string) error
	
	// Push Notifications
	SendPushNotification(ctx context.Context, notification *models.PushNotification) error
	SendBulkPushNotifications(ctx context.Context, notifications []models.PushNotification) error
	GetPushNotification(ctx context.Context, id uuid.UUID) (*models.PushNotification, error)
	ListPushNotifications(ctx context.Context, params map[string]interface{}) ([]models.PushNotification, int64, error)
	ListDeviceNotifications(ctx context.Context, deviceID uuid.UUID, params map[string]interface{}) ([]models.PushNotification, int64, error)
	ProcessPendingNotifications(ctx context.Context, limit int) error
	MarkNotificationDelivered(ctx context.Context, id uuid.UUID) error
	MarkNotificationClicked(ctx context.Context, id uuid.UUID) error
	SendNotificationToUsers(ctx context.Context, userIDs []uuid.UUID, title, body, notificationType string, data map[string]interface{}) error
	CleanupExpiredNotifications(ctx context.Context) error
	
	// Session Management
	CreateMobileSession(ctx context.Context, deviceToken string, userID uuid.UUID) (*models.MobileSession, error)
	UpdateMobileSession(ctx context.Context, sessionToken string, updates map[string]interface{}) error
	GetMobileSession(ctx context.Context, sessionToken string) (*models.MobileSession, error)
	ListMobileSessions(ctx context.Context, params map[string]interface{}) ([]models.MobileSession, int64, error)
	EndMobileSession(ctx context.Context, sessionToken string) error
	ValidateSession(ctx context.Context, sessionToken string) (*models.MobileSession, error)
	
	// Analytics
	TrackEvent(ctx context.Context, event *models.MobileAnalytics) error
	TrackScreenView(ctx context.Context, deviceID, userID uuid.UUID, screenName, screenClass string, duration int64) error
	TrackUserInteraction(ctx context.Context, deviceID, userID uuid.UUID, interactionType, elementID string, data map[string]interface{}) error
	GetAnalyticsSummary(ctx context.Context, companyID uuid.UUID, startDate, endDate time.Time) (map[string]interface{}, error)
	ListAnalyticsEvents(ctx context.Context, params map[string]interface{}) ([]models.MobileAnalytics, int64, error)
	CleanupOldAnalytics(ctx context.Context, days int) error
	
	// App Version Management
	CreateAppVersion(ctx context.Context, version *models.MobileAppVersion, userID uuid.UUID) error
	UpdateAppVersion(ctx context.Context, version *models.MobileAppVersion) error
	GetLatestAppVersion(ctx context.Context, platform string, companyID *uuid.UUID) (*models.MobileAppVersion, error)
	ListAppVersions(ctx context.Context, params map[string]interface{}) ([]models.MobileAppVersion, int64, error)
	CheckForUpdates(ctx context.Context, deviceToken, currentVersion string) (*models.MobileAppVersion, bool, error)
	
	// Offline Data Sync
	CreateOfflineData(ctx context.Context, data *models.MobileOfflineData) error
	SyncOfflineData(ctx context.Context, deviceID uuid.UUID, limit int) error
	ListPendingOfflineData(ctx context.Context, deviceID uuid.UUID, limit int) ([]models.MobileOfflineData, error)
	MarkOfflineDataSynced(ctx context.Context, id uuid.UUID) error
	ResolveDataConflict(ctx context.Context, id uuid.UUID, resolution string, resolvedBy uuid.UUID) error
	
	// Configuration Management
	GetMobileConfig(ctx context.Context, key, platform string, companyID *uuid.UUID) (interface{}, error)
	GetMobileConfigs(ctx context.Context, platform string, companyID *uuid.UUID) (map[string]interface{}, error)
	SetMobileConfig(ctx context.Context, config *models.MobileConfiguration, userID uuid.UUID) error
	
	// Business Operations
	GetMobileStatistics(ctx context.Context, companyID *uuid.UUID) (map[string]interface{}, error)
	GetDeviceUsageStats(ctx context.Context, deviceID uuid.UUID, days int) (map[string]interface{}, error)
	GetNotificationStats(ctx context.Context, companyID uuid.UUID, days int) (map[string]interface{}, error)
	GenerateMobileDashboard(ctx context.Context, companyID uuid.UUID) (map[string]interface{}, error)
}

type mobileService struct {
	mobileRepo repository.MobileRepository
	userRepo   repository.UserRepository
}

func NewMobileService(mobileRepo repository.MobileRepository, userRepo repository.UserRepository) MobileService {
	return &mobileService{
		mobileRepo: mobileRepo,
		userRepo:   userRepo,
	}
}

// Device Management
func (s *mobileService) RegisterDevice(ctx context.Context, device *models.MobileDevice) error {
	// Check if device already exists
	existingDevice, _ := s.mobileRepo.GetDeviceByToken(device.DeviceToken)
	if existingDevice != nil {
		// Update existing device
		existingDevice.DeviceModel = device.DeviceModel
		existingDevice.OSVersion = device.OSVersion
		existingDevice.AppVersion = device.AppVersion
		existingDevice.DeviceName = device.DeviceName
		existingDevice.TimeZone = device.TimeZone
		existingDevice.Language = device.Language
		existingDevice.Country = device.Country
		existingDevice.IsActive = true
		existingDevice.LastSeen = time.Now()
		return s.mobileRepo.UpdateDevice(existingDevice)
	}
	
	// Register new device
	device.RegisteredAt = time.Now()
	device.LastSeen = time.Now()
	device.IsActive = true
	device.PushEnabled = true
	device.BadgeCount = 0
	device.SecurityLevel = "normal"
	
	return s.mobileRepo.RegisterDevice(device)
}

func (s *mobileService) UpdateDevice(ctx context.Context, device *models.MobileDevice) error {
	return s.mobileRepo.UpdateDevice(device)
}

func (s *mobileService) GetDevice(ctx context.Context, id uuid.UUID) (*models.MobileDevice, error) {
	return s.mobileRepo.GetDevice(id)
}

func (s *mobileService) GetDeviceByToken(ctx context.Context, deviceToken string) (*models.MobileDevice, error) {
	return s.mobileRepo.GetDeviceByToken(deviceToken)
}

func (s *mobileService) ListUserDevices(ctx context.Context, userID uuid.UUID, params map[string]interface{}) ([]models.MobileDevice, error) {
	return s.mobileRepo.ListUserDevices(userID, params)
}

func (s *mobileService) ListCompanyDevices(ctx context.Context, companyID uuid.UUID, params map[string]interface{}) ([]models.MobileDevice, int64, error) {
	return s.mobileRepo.ListCompanyDevices(companyID, params)
}

func (s *mobileService) DeactivateDevice(ctx context.Context, id uuid.UUID) error {
	return s.mobileRepo.DeactivateDevice(id)
}

func (s *mobileService) UpdateDeviceLastSeen(ctx context.Context, deviceToken string) error {
	return s.mobileRepo.UpdateDeviceLastSeen(deviceToken, time.Now())
}

// Push Notifications
func (s *mobileService) SendPushNotification(ctx context.Context, notification *models.PushNotification) error {
	// Set default values
	if notification.Priority == "" {
		notification.Priority = "normal"
	}
	if notification.Status == "" {
		notification.Status = "pending"
	}
	if notification.MaxRetries == 0 {
		notification.MaxRetries = 3
	}
	if notification.TTL == 0 {
		notification.TTL = 86400 // 24 hours
	}
	
	// Set expiration
	if notification.ExpiresAt == nil {
		expiresAt := time.Now().Add(time.Duration(notification.TTL) * time.Second)
		notification.ExpiresAt = &expiresAt
	}
	
	// Create notification record
	if err := s.mobileRepo.CreatePushNotification(notification); err != nil {
		return err
	}
	
	// Process notification immediately if not scheduled
	if notification.ScheduledAt == nil || notification.ScheduledAt.Before(time.Now()) {
		go s.processPushNotification(notification)
	}
	
	return nil
}

func (s *mobileService) SendBulkPushNotifications(ctx context.Context, notifications []models.PushNotification) error {
	for _, notification := range notifications {
		if err := s.SendPushNotification(ctx, &notification); err != nil {
			log.Printf("Failed to send push notification to device %s: %v", notification.DeviceID, err)
			continue
		}
	}
	return nil
}

func (s *mobileService) processPushNotification(notification *models.PushNotification) {
	// Get device information
	device, err := s.mobileRepo.GetDevice(notification.DeviceID)
	if err != nil {
		s.mobileRepo.MarkNotificationFailed(notification.ID, "Device not found")
		return
	}
	
	if !device.IsActive || !device.PushEnabled {
		s.mobileRepo.MarkNotificationFailed(notification.ID, "Device inactive or push disabled")
		return
	}
	
	// Simulate sending to FCM/APNS
	// In real implementation, this would call the actual push notification service
	time.Sleep(100 * time.Millisecond) // Simulate network delay
	
	// Mark as sent (simulate success)
	providerResponse := fmt.Sprintf("Success: Message sent to %s", device.Platform)
	if err := s.mobileRepo.MarkNotificationSent(notification.ID, time.Now(), providerResponse); err != nil {
		log.Printf("Failed to mark notification as sent: %v", err)
	}
	
	// Simulate delivery confirmation after a short delay
	go func() {
		time.Sleep(2 * time.Second)
		s.mobileRepo.MarkNotificationDelivered(notification.ID, time.Now())
	}()
}

func (s *mobileService) GetPushNotification(ctx context.Context, id uuid.UUID) (*models.PushNotification, error) {
	return s.mobileRepo.GetPushNotification(id)
}

func (s *mobileService) ListPushNotifications(ctx context.Context, params map[string]interface{}) ([]models.PushNotification, int64, error) {
	return s.mobileRepo.ListPushNotifications(params)
}

func (s *mobileService) ListDeviceNotifications(ctx context.Context, deviceID uuid.UUID, params map[string]interface{}) ([]models.PushNotification, int64, error) {
	return s.mobileRepo.ListDeviceNotifications(deviceID, params)
}

func (s *mobileService) ProcessPendingNotifications(ctx context.Context, limit int) error {
	notifications, err := s.mobileRepo.GetPendingNotifications(limit)
	if err != nil {
		return err
	}
	
	for _, notification := range notifications {
		go s.processPushNotification(&notification)
	}
	
	return nil
}

func (s *mobileService) MarkNotificationDelivered(ctx context.Context, id uuid.UUID) error {
	return s.mobileRepo.MarkNotificationDelivered(id, time.Now())
}

func (s *mobileService) MarkNotificationClicked(ctx context.Context, id uuid.UUID) error {
	return s.mobileRepo.MarkNotificationClicked(id, time.Now())
}

func (s *mobileService) SendNotificationToUsers(ctx context.Context, userIDs []uuid.UUID, title, body, notificationType string, data map[string]interface{}) error {
	for _, userID := range userIDs {
		// Get user's devices
		devices, err := s.mobileRepo.ListUserDevices(userID, map[string]interface{}{
			"is_active": true,
		})
		if err != nil {
			continue
		}
		
		// Get user info
		user, err := s.userRepo.GetUser(userID)
		if err != nil {
			continue
		}
		
		// Send notification to each device
		for _, device := range devices {
			dataJSON, _ := json.Marshal(data)
			
			notification := &models.PushNotification{
				DeviceID:  device.ID,
				UserID:    userID,
				CompanyID: user.CompanyID,
				Title:     title,
				Body:      body,
				Type:      notificationType,
				Data:      string(dataJSON),
				Priority:  "normal",
				Badge:     1,
			}
			
			if err := s.SendPushNotification(ctx, notification); err != nil {
				log.Printf("Failed to send notification to user %s device %s: %v", userID, device.ID, err)
			}
		}
	}
	
	return nil
}

func (s *mobileService) CleanupExpiredNotifications(ctx context.Context) error {
	return s.mobileRepo.DeleteExpiredNotifications()
}

// Session Management
func (s *mobileService) CreateMobileSession(ctx context.Context, deviceToken string, userID uuid.UUID) (*models.MobileSession, error) {
	device, err := s.mobileRepo.GetDeviceByToken(deviceToken)
	if err != nil {
		return nil, err
	}
	
	user, err := s.userRepo.GetUser(userID)
	if err != nil {
		return nil, err
	}
	
	// Generate session tokens
	sessionToken, err := generateToken(32)
	if err != nil {
		return nil, err
	}
	
	refreshToken, err := generateToken(32)
	if err != nil {
		return nil, err
	}
	
	session := &models.MobileSession{
		DeviceID:     device.ID,
		UserID:       userID,
		CompanyID:    user.CompanyID,
		SessionToken: sessionToken,
		RefreshToken: refreshToken,
		StartTime:    time.Now(),
		LastActivity: time.Now(),
		AppState:     "active",
		NetworkType:  "unknown",
	}
	
	if err := s.mobileRepo.CreateMobileSession(session); err != nil {
		return nil, err
	}
	
	// Update device last seen
	s.mobileRepo.UpdateDeviceLastSeen(deviceToken, time.Now())
	
	return session, nil
}

func (s *mobileService) UpdateMobileSession(ctx context.Context, sessionToken string, updates map[string]interface{}) error {
	session, err := s.mobileRepo.GetMobileSession(sessionToken)
	if err != nil {
		return err
	}
	
	// Update fields
	if appState, ok := updates["app_state"].(string); ok {
		session.AppState = appState
	}
	if networkType, ok := updates["network_type"].(string); ok {
		session.NetworkType = networkType
	}
	if screenViews, ok := updates["screen_views"].(int); ok {
		session.ScreenViews += screenViews
	}
	if apiRequests, ok := updates["api_requests"].(int); ok {
		session.APIRequests += apiRequests
	}
	if dataTransferred, ok := updates["data_transferred"].(int64); ok {
		session.DataTransferred += dataTransferred
	}
	
	session.LastActivity = time.Now()
	
	return s.mobileRepo.UpdateMobileSession(session)
}

func (s *mobileService) GetMobileSession(ctx context.Context, sessionToken string) (*models.MobileSession, error) {
	return s.mobileRepo.GetMobileSession(sessionToken)
}

func (s *mobileService) ListMobileSessions(ctx context.Context, params map[string]interface{}) ([]models.MobileSession, int64, error) {
	return s.mobileRepo.ListMobileSessions(params)
}

func (s *mobileService) EndMobileSession(ctx context.Context, sessionToken string) error {
	return s.mobileRepo.EndMobileSession(sessionToken, time.Now())
}

func (s *mobileService) ValidateSession(ctx context.Context, sessionToken string) (*models.MobileSession, error) {
	session, err := s.mobileRepo.GetMobileSession(sessionToken)
	if err != nil {
		return nil, err
	}
	
	// Check if session is still valid (within 24 hours of last activity)
	if time.Since(session.LastActivity) > 24*time.Hour {
		return nil, fmt.Errorf("session expired")
	}
	
	return session, nil
}

// Analytics
func (s *mobileService) TrackEvent(ctx context.Context, event *models.MobileAnalytics) error {
	event.EventTimestamp = time.Now()
	return s.mobileRepo.CreateAnalyticsEvent(event)
}

func (s *mobileService) TrackScreenView(ctx context.Context, deviceID, userID uuid.UUID, screenName, screenClass string, duration int64) error {
	user, err := s.userRepo.GetUser(userID)
	if err != nil {
		return err
	}
	
	event := &models.MobileAnalytics{
		DeviceID:       deviceID,
		UserID:         userID,
		CompanyID:      user.CompanyID,
		EventType:      "screen_view",
		EventName:      "screen_view",
		EventCategory:  "navigation",
		ScreenName:     screenName,
		ScreenClass:    screenClass,
		Duration:       duration,
		EventTimestamp: time.Now(),
	}
	
	return s.mobileRepo.CreateAnalyticsEvent(event)
}

func (s *mobileService) TrackUserInteraction(ctx context.Context, deviceID, userID uuid.UUID, interactionType, elementID string, data map[string]interface{}) error {
	user, err := s.userRepo.GetUser(userID)
	if err != nil {
		return err
	}
	
	eventData, _ := json.Marshal(data)
	
	event := &models.MobileAnalytics{
		DeviceID:        deviceID,
		UserID:          userID,
		CompanyID:       user.CompanyID,
		EventType:       "user_interaction",
		EventName:       interactionType,
		EventCategory:   "interaction",
		EventData:       string(eventData),
		InteractionType: interactionType,
		ElementID:       elementID,
		EventTimestamp:  time.Now(),
	}
	
	return s.mobileRepo.CreateAnalyticsEvent(event)
}

func (s *mobileService) GetAnalyticsSummary(ctx context.Context, companyID uuid.UUID, startDate, endDate time.Time) (map[string]interface{}, error) {
	return s.mobileRepo.GetAnalyticsSummary(companyID, startDate, endDate)
}

func (s *mobileService) ListAnalyticsEvents(ctx context.Context, params map[string]interface{}) ([]models.MobileAnalytics, int64, error) {
	return s.mobileRepo.ListAnalyticsEvents(params)
}

func (s *mobileService) CleanupOldAnalytics(ctx context.Context, days int) error {
	return s.mobileRepo.CleanupOldAnalytics(days)
}

// App Version Management
func (s *mobileService) CreateAppVersion(ctx context.Context, version *models.MobileAppVersion, userID uuid.UUID) error {
	version.CreatedBy = userID
	version.Status = "draft"
	return s.mobileRepo.CreateAppVersion(version)
}

func (s *mobileService) UpdateAppVersion(ctx context.Context, version *models.MobileAppVersion) error {
	return s.mobileRepo.UpdateAppVersion(version)
}

func (s *mobileService) GetLatestAppVersion(ctx context.Context, platform string, companyID *uuid.UUID) (*models.MobileAppVersion, error) {
	return s.mobileRepo.GetLatestAppVersion(platform, companyID)
}

func (s *mobileService) ListAppVersions(ctx context.Context, params map[string]interface{}) ([]models.MobileAppVersion, int64, error) {
	return s.mobileRepo.ListAppVersions(params)
}

func (s *mobileService) CheckForUpdates(ctx context.Context, deviceToken, currentVersion string) (*models.MobileAppVersion, bool, error) {
	device, err := s.mobileRepo.GetDeviceByToken(deviceToken)
	if err != nil {
		return nil, false, err
	}
	
	latestVersion, err := s.mobileRepo.GetLatestAppVersion(device.Platform, &device.CompanyID)
	if err != nil {
		return nil, false, err
	}
	
	// Simple version comparison (in production, use semantic versioning)
	hasUpdate := latestVersion.Version != currentVersion
	
	return latestVersion, hasUpdate, nil
}

// Offline Data Sync
func (s *mobileService) CreateOfflineData(ctx context.Context, data *models.MobileOfflineData) error {
	data.Status = "pending"
	data.Priority = 1 // Normal priority
	data.MaxRetries = 5
	return s.mobileRepo.CreateOfflineData(data)
}

func (s *mobileService) SyncOfflineData(ctx context.Context, deviceID uuid.UUID, limit int) error {
	pendingData, err := s.mobileRepo.ListPendingOfflineData(deviceID, limit)
	if err != nil {
		return err
	}
	
	for _, data := range pendingData {
		// Process sync operation
		if err := s.processSyncOperation(&data); err != nil {
			// Mark as failed and increment retry count
			data.RetryCount++
			if data.RetryCount >= data.MaxRetries {
				s.mobileRepo.MarkOfflineDataFailed(data.ID, err.Error())
			} else {
				data.ErrorMessage = err.Error()
				s.mobileRepo.UpdateOfflineData(&data)
			}
			continue
		}
		
		// Mark as synced
		s.mobileRepo.MarkOfflineDataSynced(data.ID, time.Now())
	}
	
	return nil
}

func (s *mobileService) processSyncOperation(data *models.MobileOfflineData) error {
	// Simulate sync operation based on data type and operation
	// In real implementation, this would handle the actual data synchronization
	switch data.DataType {
	case "inquiry":
		return s.syncInquiryData(data)
	case "quote":
		return s.syncQuoteData(data)
	case "order":
		return s.syncOrderData(data)
	default:
		return fmt.Errorf("unsupported data type: %s", data.DataType)
	}
}

func (s *mobileService) syncInquiryData(data *models.MobileOfflineData) error {
	// Implement inquiry sync logic
	time.Sleep(100 * time.Millisecond) // Simulate processing time
	return nil
}

func (s *mobileService) syncQuoteData(data *models.MobileOfflineData) error {
	// Implement quote sync logic
	time.Sleep(100 * time.Millisecond) // Simulate processing time
	return nil
}

func (s *mobileService) syncOrderData(data *models.MobileOfflineData) error {
	// Implement order sync logic
	time.Sleep(100 * time.Millisecond) // Simulate processing time
	return nil
}

func (s *mobileService) ListPendingOfflineData(ctx context.Context, deviceID uuid.UUID, limit int) ([]models.MobileOfflineData, error) {
	return s.mobileRepo.ListPendingOfflineData(deviceID, limit)
}

func (s *mobileService) MarkOfflineDataSynced(ctx context.Context, id uuid.UUID) error {
	return s.mobileRepo.MarkOfflineDataSynced(id, time.Now())
}

func (s *mobileService) ResolveDataConflict(ctx context.Context, id uuid.UUID, resolution string, resolvedBy uuid.UUID) error {
	data, err := s.mobileRepo.GetOfflineData(id)
	if err != nil {
		return err
	}
	
	data.Status = "synced"
	data.ResolvedBy = &resolvedBy
	data.ResolutionMethod = resolution
	data.SyncedAt = &[]time.Time{time.Now()}[0]
	
	return s.mobileRepo.UpdateOfflineData(data)
}

// Configuration Management
func (s *mobileService) GetMobileConfig(ctx context.Context, key, platform string, companyID *uuid.UUID) (interface{}, error) {
	config, err := s.mobileRepo.GetMobileConfigByKey(key, platform, companyID)
	if err != nil {
		return nil, err
	}
	
	// Parse config value based on type
	switch config.ConfigType {
	case "string":
		return config.ConfigValue, nil
	case "number":
		var value float64
		if err := json.Unmarshal([]byte(config.ConfigValue), &value); err != nil {
			return nil, err
		}
		return value, nil
	case "boolean":
		var value bool
		if err := json.Unmarshal([]byte(config.ConfigValue), &value); err != nil {
			return nil, err
		}
		return value, nil
	case "object", "array":
		var value interface{}
		if err := json.Unmarshal([]byte(config.ConfigValue), &value); err != nil {
			return nil, err
		}
		return value, nil
	default:
		return config.ConfigValue, nil
	}
}

func (s *mobileService) GetMobileConfigs(ctx context.Context, platform string, companyID *uuid.UUID) (map[string]interface{}, error) {
	configs, err := s.mobileRepo.ListMobileConfigs(map[string]interface{}{
		"platform":   platform,
		"company_id": companyID,
		"is_enabled": true,
	})
	if err != nil {
		return nil, err
	}
	
	configMap := make(map[string]interface{})
	for _, config := range configs {
		value, err := s.GetMobileConfig(ctx, config.ConfigKey, platform, companyID)
		if err == nil {
			configMap[config.ConfigKey] = value
		}
	}
	
	return configMap, nil
}

func (s *mobileService) SetMobileConfig(ctx context.Context, config *models.MobileConfiguration, userID uuid.UUID) error {
	config.CreatedBy = userID
	return s.mobileRepo.CreateMobileConfig(config)
}

// Business Operations
func (s *mobileService) GetMobileStatistics(ctx context.Context, companyID *uuid.UUID) (map[string]interface{}, error) {
	return s.mobileRepo.GetMobileStatistics(companyID)
}

func (s *mobileService) GetDeviceUsageStats(ctx context.Context, deviceID uuid.UUID, days int) (map[string]interface{}, error) {
	return s.mobileRepo.GetDeviceUsageStats(deviceID, days)
}

func (s *mobileService) GetNotificationStats(ctx context.Context, companyID uuid.UUID, days int) (map[string]interface{}, error) {
	return s.mobileRepo.GetNotificationStats(companyID, days)
}

func (s *mobileService) GenerateMobileDashboard(ctx context.Context, companyID uuid.UUID) (map[string]interface{}, error) {
	dashboard := make(map[string]interface{})
	
	// Get mobile statistics
	stats, err := s.mobileRepo.GetMobileStatistics(&companyID)
	if err != nil {
		return nil, err
	}
	dashboard["statistics"] = stats
	
	// Get notification statistics for last 30 days
	notificationStats, err := s.mobileRepo.GetNotificationStats(companyID, 30)
	if err != nil {
		return nil, err
	}
	dashboard["notification_stats"] = notificationStats
	
	// Get analytics summary for last 7 days
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -7)
	analyticsStats, err := s.mobileRepo.GetAnalyticsSummary(companyID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	dashboard["analytics_stats"] = analyticsStats
	
	// Add timestamp
	dashboard["generated_at"] = time.Now()
	
	return dashboard, nil
}

// Helper functions
func generateToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}