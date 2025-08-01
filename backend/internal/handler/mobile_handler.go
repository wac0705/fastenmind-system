package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/fastenmind/fastener-api/internal/models"
	"github.com/fastenmind/fastener-api/internal/service"
	"github.com/fastenmind/fastener-api/pkg/response"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type MobileHandler struct {
	mobileService service.MobileService
}

func NewMobileHandler(mobileService service.MobileService) *MobileHandler {
	return &MobileHandler{
		mobileService: mobileService,
	}
}

// Device Management endpoints
func (h *MobileHandler) RegisterDevice(c echo.Context) error {
	var device models.MobileDevice
	if err := c.Bind(&device); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body", err)
	}

	if err := h.mobileService.RegisterDevice(c.Request().Context(), &device); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to register device", err)
	}

	return response.Success(c, device, "Device registered successfully")
}

func (h *MobileHandler) UpdateDevice(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid device ID", err)
	}

	var device models.MobileDevice
	if err := c.Bind(&device); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body", err)
	}

	device.ID = id
	if err := h.mobileService.UpdateDevice(c.Request().Context(), &device); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to update device", err)
	}

	return response.Success(c, device, "Device updated successfully")
}

func (h *MobileHandler) GetDevice(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid device ID", err)
	}

	device, err := h.mobileService.GetDevice(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, http.StatusNotFound, "Device not found", err)
	}

	return response.Success(c, device, "Device retrieved successfully")
}

func (h *MobileHandler) GetDeviceByToken(c echo.Context) error {
	deviceToken := c.Param("token")

	device, err := h.mobileService.GetDeviceByToken(c.Request().Context(), deviceToken)
	if err != nil {
		return response.Error(c, http.StatusNotFound, "Device not found", err)
	}

	return response.Success(c, device, "Device retrieved successfully")
}

func (h *MobileHandler) ListUserDevices(c echo.Context) error {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid user ID", err)
	}

	params := getMobileListParams(c)
	devices, err := h.mobileService.ListUserDevices(c.Request().Context(), userID, params)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to list user devices", err)
	}

	return response.Success(c, map[string]interface{}{
		"data":  devices,
		"total": len(devices),
	}, "User devices retrieved successfully")
}

func (h *MobileHandler) ListCompanyDevices(c echo.Context) error {
	companyID := *getCompanyIDFromContext(c)
	params := getMobileListParams(c)

	devices, total, err := h.mobileService.ListCompanyDevices(c.Request().Context(), companyID, params)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to list company devices", err)
	}

	return response.Success(c, map[string]interface{}{
		"data":  devices,
		"total": total,
	}, "Company devices retrieved successfully")
}

func (h *MobileHandler) DeactivateDevice(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid device ID", err)
	}

	if err := h.mobileService.DeactivateDevice(c.Request().Context(), id); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to deactivate device", err)
	}

	return response.Success(c, nil, "Device deactivated successfully")
}

func (h *MobileHandler) UpdateDeviceLastSeen(c echo.Context) error {
	deviceToken := c.Param("token")

	if err := h.mobileService.UpdateDeviceLastSeen(c.Request().Context(), deviceToken); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to update device last seen", err)
	}

	return response.Success(c, nil, "Device last seen updated successfully")
}

// Push Notification endpoints
func (h *MobileHandler) SendPushNotification(c echo.Context) error {
	var notification models.PushNotification
	if err := c.Bind(&notification); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body", err)
	}

	if err := h.mobileService.SendPushNotification(c.Request().Context(), &notification); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to send push notification", err)
	}

	return response.Success(c, notification, "Push notification sent successfully")
}

func (h *MobileHandler) SendBulkPushNotifications(c echo.Context) error {
	var req struct {
		Notifications []models.PushNotification `json:"notifications"`
	}
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body", err)
	}

	if err := h.mobileService.SendBulkPushNotifications(c.Request().Context(), req.Notifications); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to send bulk push notifications", err)
	}

	return response.Success(c, map[string]interface{}{
		"count": len(req.Notifications),
	}, "Bulk push notifications sent successfully")
}

func (h *MobileHandler) ListPushNotifications(c echo.Context) error {
	params := getMobileListParams(c)
	if companyID := getCompanyIDFromContext(c); companyID != nil {
		params["company_id"] = *companyID
	}

	notifications, total, err := h.mobileService.ListPushNotifications(c.Request().Context(), params)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to list push notifications", err)
	}

	return response.Success(c, map[string]interface{}{
		"data":  notifications,
		"total": total,
	}, "Push notifications retrieved successfully")
}

func (h *MobileHandler) ListDeviceNotifications(c echo.Context) error {
	deviceIDStr := c.Param("deviceId")
	deviceID, err := uuid.Parse(deviceIDStr)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid device ID", err)
	}

	params := getMobileListParams(c)
	notifications, total, err := h.mobileService.ListDeviceNotifications(c.Request().Context(), deviceID, params)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to list device notifications", err)
	}

	return response.Success(c, map[string]interface{}{
		"data":  notifications,
		"total": total,
	}, "Device notifications retrieved successfully")
}

func (h *MobileHandler) MarkNotificationDelivered(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid notification ID", err)
	}

	if err := h.mobileService.MarkNotificationDelivered(c.Request().Context(), id); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to mark notification as delivered", err)
	}

	return response.Success(c, nil, "Notification marked as delivered successfully")
}

func (h *MobileHandler) MarkNotificationClicked(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid notification ID", err)
	}

	if err := h.mobileService.MarkNotificationClicked(c.Request().Context(), id); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to mark notification as clicked", err)
	}

	return response.Success(c, nil, "Notification marked as clicked successfully")
}

func (h *MobileHandler) SendNotificationToUsers(c echo.Context) error {
	var req struct {
		UserIDs []uuid.UUID            `json:"user_ids"`
		Title   string                 `json:"title"`
		Body    string                 `json:"body"`
		Type    string                 `json:"type"`
		Data    map[string]interface{} `json:"data"`
	}
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body", err)
	}

	if err := h.mobileService.SendNotificationToUsers(c.Request().Context(), req.UserIDs, req.Title, req.Body, req.Type, req.Data); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to send notifications to users", err)
	}

	return response.Success(c, map[string]interface{}{
		"user_count": len(req.UserIDs),
	}, "Notifications sent to users successfully")
}

func (h *MobileHandler) ProcessPendingNotifications(c echo.Context) error {
	limit := 50
	if limitStr := c.QueryParam("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	if err := h.mobileService.ProcessPendingNotifications(c.Request().Context(), limit); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to process pending notifications", err)
	}

	return response.Success(c, nil, "Pending notifications processed successfully")
}

// Session Management endpoints
func (h *MobileHandler) CreateMobileSession(c echo.Context) error {
	var req struct {
		DeviceToken string    `json:"device_token"`
		UserID      uuid.UUID `json:"user_id"`
	}
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body", err)
	}

	session, err := h.mobileService.CreateMobileSession(c.Request().Context(), req.DeviceToken, req.UserID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to create mobile session", err)
	}

	return response.Success(c, session, "Mobile session created successfully")
}

func (h *MobileHandler) UpdateMobileSession(c echo.Context) error {
	sessionToken := c.Param("token")

	var updates map[string]interface{}
	if err := c.Bind(&updates); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body", err)
	}

	if err := h.mobileService.UpdateMobileSession(c.Request().Context(), sessionToken, updates); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to update mobile session", err)
	}

	return response.Success(c, nil, "Mobile session updated successfully")
}

func (h *MobileHandler) GetMobileSession(c echo.Context) error {
	sessionToken := c.Param("token")

	session, err := h.mobileService.GetMobileSession(c.Request().Context(), sessionToken)
	if err != nil {
		return response.Error(c, http.StatusNotFound, "Mobile session not found", err)
	}

	return response.Success(c, session, "Mobile session retrieved successfully")
}

func (h *MobileHandler) EndMobileSession(c echo.Context) error {
	sessionToken := c.Param("token")

	if err := h.mobileService.EndMobileSession(c.Request().Context(), sessionToken); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to end mobile session", err)
	}

	return response.Success(c, nil, "Mobile session ended successfully")
}

func (h *MobileHandler) ValidateSession(c echo.Context) error {
	sessionToken := c.Param("token")

	session, err := h.mobileService.ValidateSession(c.Request().Context(), sessionToken)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "Invalid session", err)
	}

	return response.Success(c, map[string]interface{}{
		"valid":   true,
		"session": session,
	}, "Session validated successfully")
}

// Analytics endpoints
func (h *MobileHandler) TrackEvent(c echo.Context) error {
	var event models.MobileAnalytics
	if err := c.Bind(&event); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body", err)
	}

	if err := h.mobileService.TrackEvent(c.Request().Context(), &event); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to track event", err)
	}

	return response.Success(c, nil, "Event tracked successfully")
}

func (h *MobileHandler) TrackScreenView(c echo.Context) error {
	var req struct {
		DeviceID    uuid.UUID `json:"device_id"`
		UserID      uuid.UUID `json:"user_id"`
		ScreenName  string    `json:"screen_name"`
		ScreenClass string    `json:"screen_class"`
		Duration    int64     `json:"duration"`
	}
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body", err)
	}

	if err := h.mobileService.TrackScreenView(c.Request().Context(), req.DeviceID, req.UserID, req.ScreenName, req.ScreenClass, req.Duration); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to track screen view", err)
	}

	return response.Success(c, nil, "Screen view tracked successfully")
}

func (h *MobileHandler) TrackUserInteraction(c echo.Context) error {
	var req struct {
		DeviceID        uuid.UUID              `json:"device_id"`
		UserID          uuid.UUID              `json:"user_id"`
		InteractionType string                 `json:"interaction_type"`
		ElementID       string                 `json:"element_id"`
		Data            map[string]interface{} `json:"data"`
	}
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body", err)
	}

	if err := h.mobileService.TrackUserInteraction(c.Request().Context(), req.DeviceID, req.UserID, req.InteractionType, req.ElementID, req.Data); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to track user interaction", err)
	}

	return response.Success(c, nil, "User interaction tracked successfully")
}

func (h *MobileHandler) GetAnalyticsSummary(c echo.Context) error {
	companyID := *getCompanyIDFromContext(c)

	startDateStr := c.QueryParam("start_date")
	endDateStr := c.QueryParam("end_date")

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			return response.Error(c, http.StatusBadRequest, "Invalid start date format", err)
		}
	} else {
		startDate = time.Now().AddDate(0, 0, -30) // Last 30 days
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return response.Error(c, http.StatusBadRequest, "Invalid end date format", err)
		}
	} else {
		endDate = time.Now()
	}

	summary, err := h.mobileService.GetAnalyticsSummary(c.Request().Context(), companyID, startDate, endDate)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to get analytics summary", err)
	}

	return response.Success(c, summary, "Analytics summary retrieved successfully")
}

// App Version Management endpoints
func (h *MobileHandler) CreateAppVersion(c echo.Context) error {
	var version models.MobileAppVersion
	if err := c.Bind(&version); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body", err)
	}

	userID := getUserIDFromContext(c)
	if err := h.mobileService.CreateAppVersion(c.Request().Context(), &version, userID); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to create app version", err)
	}

	return response.Success(c, version, "App version created successfully")
}

func (h *MobileHandler) UpdateAppVersion(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid version ID", err)
	}

	var version models.MobileAppVersion
	if err := c.Bind(&version); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body", err)
	}

	version.ID = id
	if err := h.mobileService.UpdateAppVersion(c.Request().Context(), &version); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to update app version", err)
	}

	return response.Success(c, version, "App version updated successfully")
}

func (h *MobileHandler) GetLatestAppVersion(c echo.Context) error {
	platform := c.QueryParam("platform")
	if platform == "" {
		return response.Error(c, http.StatusBadRequest, "Platform parameter is required", nil)
	}

	companyID := getCompanyIDFromContext(c)
	version, err := h.mobileService.GetLatestAppVersion(c.Request().Context(), platform, companyID)
	if err != nil {
		return response.Error(c, http.StatusNotFound, "No app version found", err)
	}

	return response.Success(c, version, "Latest app version retrieved successfully")
}

func (h *MobileHandler) CheckForUpdates(c echo.Context) error {
	deviceToken := c.QueryParam("device_token")
	currentVersion := c.QueryParam("current_version")

	if deviceToken == "" || currentVersion == "" {
		return response.Error(c, http.StatusBadRequest, "Device token and current version are required", nil)
	}

	latestVersion, hasUpdate, err := h.mobileService.CheckForUpdates(c.Request().Context(), deviceToken, currentVersion)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to check for updates", err)
	}

	return response.Success(c, map[string]interface{}{
		"has_update":     hasUpdate,
		"latest_version": latestVersion,
	}, "Update check completed successfully")
}

// Offline Data Sync endpoints
func (h *MobileHandler) CreateOfflineData(c echo.Context) error {
	var data models.MobileOfflineData
	if err := c.Bind(&data); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body", err)
	}

	if err := h.mobileService.CreateOfflineData(c.Request().Context(), &data); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to create offline data", err)
	}

	return response.Success(c, data, "Offline data created successfully")
}

func (h *MobileHandler) SyncOfflineData(c echo.Context) error {
	deviceIDStr := c.Param("deviceId")
	deviceID, err := uuid.Parse(deviceIDStr)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid device ID", err)
	}

	limit := 10
	if limitStr := c.QueryParam("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	if err := h.mobileService.SyncOfflineData(c.Request().Context(), deviceID, limit); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to sync offline data", err)
	}

	return response.Success(c, nil, "Offline data synced successfully")
}

func (h *MobileHandler) ListPendingOfflineData(c echo.Context) error {
	deviceIDStr := c.Param("deviceId")
	deviceID, err := uuid.Parse(deviceIDStr)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid device ID", err)
	}

	limit := 50
	if limitStr := c.QueryParam("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	data, err := h.mobileService.ListPendingOfflineData(c.Request().Context(), deviceID, limit)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to list pending offline data", err)
	}

	return response.Success(c, map[string]interface{}{
		"data":  data,
		"total": len(data),
	}, "Pending offline data retrieved successfully")
}

// Configuration Management endpoints
func (h *MobileHandler) GetMobileConfig(c echo.Context) error {
	key := c.Param("key")
	platform := c.QueryParam("platform")
	if platform == "" {
		return response.Error(c, http.StatusBadRequest, "Platform parameter is required", nil)
	}

	companyID := getCompanyIDFromContext(c)
	value, err := h.mobileService.GetMobileConfig(c.Request().Context(), key, platform, companyID)
	if err != nil {
		return response.Error(c, http.StatusNotFound, "Configuration not found", err)
	}

	return response.Success(c, map[string]interface{}{
		"key":   key,
		"value": value,
	}, "Mobile configuration retrieved successfully")
}

func (h *MobileHandler) GetMobileConfigs(c echo.Context) error {
	platform := c.QueryParam("platform")
	if platform == "" {
		return response.Error(c, http.StatusBadRequest, "Platform parameter is required", nil)
	}

	companyID := getCompanyIDFromContext(c)
	configs, err := h.mobileService.GetMobileConfigs(c.Request().Context(), platform, companyID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to get mobile configurations", err)
	}

	return response.Success(c, configs, "Mobile configurations retrieved successfully")
}

func (h *MobileHandler) SetMobileConfig(c echo.Context) error {
	var config models.MobileConfiguration
	if err := c.Bind(&config); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body", err)
	}

	userID := getUserIDFromContext(c)
	if err := h.mobileService.SetMobileConfig(c.Request().Context(), &config, userID); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to set mobile configuration", err)
	}

	return response.Success(c, config, "Mobile configuration set successfully")
}

// Business Operations endpoints
func (h *MobileHandler) GetMobileStatistics(c echo.Context) error {
	companyID := getCompanyIDFromContext(c)

	stats, err := h.mobileService.GetMobileStatistics(c.Request().Context(), companyID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to get mobile statistics", err)
	}

	return response.Success(c, stats, "Mobile statistics retrieved successfully")
}

func (h *MobileHandler) GetDeviceUsageStats(c echo.Context) error {
	deviceIDStr := c.Param("deviceId")
	deviceID, err := uuid.Parse(deviceIDStr)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid device ID", err)
	}

	days := 30
	if daysStr := c.QueryParam("days"); daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil {
			days = d
		}
	}

	stats, err := h.mobileService.GetDeviceUsageStats(c.Request().Context(), deviceID, days)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to get device usage stats", err)
	}

	return response.Success(c, stats, "Device usage statistics retrieved successfully")
}

func (h *MobileHandler) GetNotificationStats(c echo.Context) error {
	companyID := *getCompanyIDFromContext(c)

	days := 30
	if daysStr := c.QueryParam("days"); daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil {
			days = d
		}
	}

	stats, err := h.mobileService.GetNotificationStats(c.Request().Context(), companyID, days)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to get notification stats", err)
	}

	return response.Success(c, stats, "Notification statistics retrieved successfully")
}

func (h *MobileHandler) GenerateMobileDashboard(c echo.Context) error {
	companyID := *getCompanyIDFromContext(c)

	dashboard, err := h.mobileService.GenerateMobileDashboard(c.Request().Context(), companyID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to generate mobile dashboard", err)
	}

	return response.Success(c, dashboard, "Mobile dashboard generated successfully")
}

// Helper functions
func getMobileListParams(c echo.Context) map[string]interface{} {
	params := make(map[string]interface{})

	if page := c.QueryParam("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil {
			params["page"] = p
		}
	}

	if pageSize := c.QueryParam("page_size"); pageSize != "" {
		if ps, err := strconv.Atoi(pageSize); err == nil {
			params["page_size"] = ps
		}
	}

	if search := c.QueryParam("search"); search != "" {
		params["search"] = search
	}

	if platform := c.QueryParam("platform"); platform != "" {
		params["platform"] = platform
	}

	if status := c.QueryParam("status"); status != "" {
		params["status"] = status
	}

	if notificationType := c.QueryParam("type"); notificationType != "" {
		params["type"] = notificationType
	}

	if isActive := c.QueryParam("is_active"); isActive != "" {
		if active, err := strconv.ParseBool(isActive); err == nil {
			params["is_active"] = active
		}
	}

	if userID := c.QueryParam("user_id"); userID != "" {
		params["user_id"] = userID
	}

	if deviceID := c.QueryParam("device_id"); deviceID != "" {
		params["device_id"] = deviceID
	}

	if startDate := c.QueryParam("start_date"); startDate != "" {
		params["start_date"] = startDate
	}

	if endDate := c.QueryParam("end_date"); endDate != "" {
		params["end_date"] = endDate
	}

	if eventType := c.QueryParam("event_type"); eventType != "" {
		params["event_type"] = eventType
	}

	if dataType := c.QueryParam("data_type"); dataType != "" {
		params["data_type"] = dataType
	}

	if releaseType := c.QueryParam("release_type"); releaseType != "" {
		params["release_type"] = releaseType
	}

	if category := c.QueryParam("category"); category != "" {
		params["category"] = category
	}

	if isEnabled := c.QueryParam("is_enabled"); isEnabled != "" {
		if enabled, err := strconv.ParseBool(isEnabled); err == nil {
			params["is_enabled"] = enabled
		}
	}

	return params
}