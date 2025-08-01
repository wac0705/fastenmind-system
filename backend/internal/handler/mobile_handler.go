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
		return response.Error(c, http.StatusBadRequest, "Invalid request body")
	}

	if err := h.mobileService.RegisterDevice(c.Request().Context(), &device); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to register device")
	}

	return response.SuccessWithMessage(c, "Device registered successfully", device)
}

func (h *MobileHandler) UpdateDevice(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid device ID")
	}

	var device models.MobileDevice
	if err := c.Bind(&device); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body")
	}

	device.ID = id
	if err := h.mobileService.UpdateDevice(c.Request().Context(), &device); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to update device")
	}

	return response.SuccessWithMessage(c, "Device updated successfully", device)
}

func (h *MobileHandler) GetDevice(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid device ID")
	}

	device, err := h.mobileService.GetDevice(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, http.StatusNotFound, "Device not found")
	}

	return response.SuccessWithMessage(c, "Device retrieved successfully", device)
}

func (h *MobileHandler) GetDeviceByToken(c echo.Context) error {
	deviceToken := c.Param("token")

	device, err := h.mobileService.GetDeviceByToken(c.Request().Context(), deviceToken)
	if err != nil {
		return response.Error(c, http.StatusNotFound, "Device not found")
	}

	return response.SuccessWithMessage(c, "Device retrieved successfully", device)
}

func (h *MobileHandler) ListUserDevices(c echo.Context) error {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid user ID")
	}

	params := getMobileListParams(c)
	devices, err := h.mobileService.ListUserDevices(c.Request().Context(), userID, params)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to list user devices")
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
		return response.Error(c, http.StatusInternalServerError, "Failed to list company devices")
	}

	return response.Success(c, map[string]interface{}{
		"data":  devices,
		"total": total,
	}, "Company devices retrieved successfully")
}

func (h *MobileHandler) DeactivateDevice(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid device ID")
	}

	if err := h.mobileService.DeactivateDevice(c.Request().Context(), id); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to deactivate device")
	}

	return response.SuccessWithMessage(c, "Device deactivated successfully", nil)
}

func (h *MobileHandler) UpdateDeviceLastSeen(c echo.Context) error {
	deviceToken := c.Param("token")

	if err := h.mobileService.UpdateDeviceLastSeen(c.Request().Context(), deviceToken); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to update device last seen")
	}

	return response.SuccessWithMessage(c, "Device last seen updated successfully", nil)
}

// Push Notification endpoints
func (h *MobileHandler) SendPushNotification(c echo.Context) error {
	var notification models.PushNotification
	if err := c.Bind(&notification); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body")
	}

	if err := h.mobileService.SendPushNotification(c.Request().Context(), &notification); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to send push notification")
	}

	return response.SuccessWithMessage(c, "Push notification sent successfully", notification)
}

func (h *MobileHandler) SendBulkPushNotifications(c echo.Context) error {
	var req struct {
		Notifications []models.PushNotification `json:"notifications"`
	}
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body")
	}

	if err := h.mobileService.SendBulkPushNotifications(c.Request().Context(), req.Notifications); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to send bulk push notifications")
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
		return response.Error(c, http.StatusInternalServerError, "Failed to list push notifications")
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
		return response.Error(c, http.StatusBadRequest, "Invalid device ID")
	}

	params := getMobileListParams(c)
	notifications, total, err := h.mobileService.ListDeviceNotifications(c.Request().Context(), deviceID, params)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to list device notifications")
	}

	return response.Success(c, map[string]interface{}{
		"data":  notifications,
		"total": total,
	}, "Device notifications retrieved successfully")
}

func (h *MobileHandler) MarkNotificationDelivered(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid notification ID")
	}

	if err := h.mobileService.MarkNotificationDelivered(c.Request().Context(), id); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to mark notification as delivered")
	}

	return response.SuccessWithMessage(c, "Notification marked as delivered successfully", nil)
}

func (h *MobileHandler) MarkNotificationClicked(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid notification ID")
	}

	if err := h.mobileService.MarkNotificationClicked(c.Request().Context(), id); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to mark notification as clicked")
	}

	return response.SuccessWithMessage(c, "Notification marked as clicked successfully", nil)
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
		return response.Error(c, http.StatusBadRequest, "Invalid request body")
	}

	if err := h.mobileService.SendNotificationToUsers(c.Request().Context(), req.UserIDs, req.Title, req.Body, req.Type, req.Data); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to send notifications to users")
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
		return response.Error(c, http.StatusInternalServerError, "Failed to process pending notifications")
	}

	return response.SuccessWithMessage(c, "Pending notifications processed successfully", nil)
}

// Session Management endpoints
func (h *MobileHandler) CreateMobileSession(c echo.Context) error {
	var req struct {
		DeviceToken string    `json:"device_token"`
		UserID      uuid.UUID `json:"user_id"`
	}
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body")
	}

	session, err := h.mobileService.CreateMobileSession(c.Request().Context(), req.DeviceToken, req.UserID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to create mobile session")
	}

	return response.SuccessWithMessage(c, "Mobile session created successfully", session)
}

func (h *MobileHandler) UpdateMobileSession(c echo.Context) error {
	sessionToken := c.Param("token")

	var updates map[string]interface{}
	if err := c.Bind(&updates); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body")
	}

	if err := h.mobileService.UpdateMobileSession(c.Request().Context(), sessionToken, updates); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to update mobile session")
	}

	return response.SuccessWithMessage(c, "Mobile session updated successfully", nil)
}

func (h *MobileHandler) GetMobileSession(c echo.Context) error {
	sessionToken := c.Param("token")

	session, err := h.mobileService.GetMobileSession(c.Request().Context(), sessionToken)
	if err != nil {
		return response.Error(c, http.StatusNotFound, "Mobile session not found")
	}

	return response.SuccessWithMessage(c, "Mobile session retrieved successfully", session)
}

func (h *MobileHandler) EndMobileSession(c echo.Context) error {
	sessionToken := c.Param("token")

	if err := h.mobileService.EndMobileSession(c.Request().Context(), sessionToken); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to end mobile session")
	}

	return response.SuccessWithMessage(c, "Mobile session ended successfully", nil)
}

func (h *MobileHandler) ValidateSession(c echo.Context) error {
	sessionToken := c.Param("token")

	session, err := h.mobileService.ValidateSession(c.Request().Context(), sessionToken)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "Invalid session")
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
		return response.Error(c, http.StatusBadRequest, "Invalid request body")
	}

	if err := h.mobileService.TrackEvent(c.Request().Context(), &event); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to track event")
	}

	return response.SuccessWithMessage(c, "Event tracked successfully", nil)
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
		return response.Error(c, http.StatusBadRequest, "Invalid request body")
	}

	if err := h.mobileService.TrackScreenView(c.Request().Context(), req.DeviceID, req.UserID, req.ScreenName, req.ScreenClass, req.Duration); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to track screen view")
	}

	return response.SuccessWithMessage(c, "Screen view tracked successfully", nil)
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
		return response.Error(c, http.StatusBadRequest, "Invalid request body")
	}

	if err := h.mobileService.TrackUserInteraction(c.Request().Context(), req.DeviceID, req.UserID, req.InteractionType, req.ElementID, req.Data); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to track user interaction")
	}

	return response.SuccessWithMessage(c, "User interaction tracked successfully", nil)
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
			return response.Error(c, http.StatusBadRequest, "Invalid start date format")
		}
	} else {
		startDate = time.Now().AddDate(0, 0, -30) // Last 30 days
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return response.Error(c, http.StatusBadRequest, "Invalid end date format")
		}
	} else {
		endDate = time.Now()
	}

	summary, err := h.mobileService.GetAnalyticsSummary(c.Request().Context(), companyID, startDate, endDate)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to get analytics summary")
	}

	return response.SuccessWithMessage(c, "Analytics summary retrieved successfully", summary)
}

// App Version Management endpoints
func (h *MobileHandler) CreateAppVersion(c echo.Context) error {
	var version models.MobileAppVersion
	if err := c.Bind(&version); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body")
	}

	userID := getUserIDFromContext(c)
	if err := h.mobileService.CreateAppVersion(c.Request().Context(), &version, userID); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to create app version")
	}

	return response.SuccessWithMessage(c, "App version created successfully", version)
}

func (h *MobileHandler) UpdateAppVersion(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid version ID")
	}

	var version models.MobileAppVersion
	if err := c.Bind(&version); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body")
	}

	version.ID = id
	if err := h.mobileService.UpdateAppVersion(c.Request().Context(), &version); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to update app version")
	}

	return response.SuccessWithMessage(c, "App version updated successfully", version)
}

func (h *MobileHandler) GetLatestAppVersion(c echo.Context) error {
	platform := c.QueryParam("platform")
	if platform == "" {
		return response.Error(c, http.StatusBadRequest, "Platform parameter is required")
	}

	companyID := getCompanyIDFromContext(c)
	version, err := h.mobileService.GetLatestAppVersion(c.Request().Context(), platform, companyID)
	if err != nil {
		return response.Error(c, http.StatusNotFound, "No app version found")
	}

	return response.SuccessWithMessage(c, "Latest app version retrieved successfully", version)
}

func (h *MobileHandler) CheckForUpdates(c echo.Context) error {
	deviceToken := c.QueryParam("device_token")
	currentVersion := c.QueryParam("current_version")

	if deviceToken == "" || currentVersion == "" {
		return response.Error(c, http.StatusBadRequest, "Device token and current version are required")
	}

	latestVersion, hasUpdate, err := h.mobileService.CheckForUpdates(c.Request().Context(), deviceToken, currentVersion)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to check for updates")
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
		return response.Error(c, http.StatusBadRequest, "Invalid request body")
	}

	if err := h.mobileService.CreateOfflineData(c.Request().Context(), &data); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to create offline data")
	}

	return response.SuccessWithMessage(c, "Offline data created successfully", data)
}

func (h *MobileHandler) SyncOfflineData(c echo.Context) error {
	deviceIDStr := c.Param("deviceId")
	deviceID, err := uuid.Parse(deviceIDStr)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid device ID")
	}

	limit := 10
	if limitStr := c.QueryParam("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	if err := h.mobileService.SyncOfflineData(c.Request().Context(), deviceID, limit); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to sync offline data")
	}

	return response.SuccessWithMessage(c, "Offline data synced successfully", nil)
}

func (h *MobileHandler) ListPendingOfflineData(c echo.Context) error {
	deviceIDStr := c.Param("deviceId")
	deviceID, err := uuid.Parse(deviceIDStr)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid device ID")
	}

	limit := 50
	if limitStr := c.QueryParam("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	data, err := h.mobileService.ListPendingOfflineData(c.Request().Context(), deviceID, limit)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to list pending offline data")
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
		return response.Error(c, http.StatusBadRequest, "Platform parameter is required")
	}

	companyID := getCompanyIDFromContext(c)
	value, err := h.mobileService.GetMobileConfig(c.Request().Context(), key, platform, companyID)
	if err != nil {
		return response.Error(c, http.StatusNotFound, "Configuration not found")
	}

	return response.Success(c, map[string]interface{}{
		"key":   key,
		"value": value,
	}, "Mobile configuration retrieved successfully")
}

func (h *MobileHandler) GetMobileConfigs(c echo.Context) error {
	platform := c.QueryParam("platform")
	if platform == "" {
		return response.Error(c, http.StatusBadRequest, "Platform parameter is required")
	}

	companyID := getCompanyIDFromContext(c)
	configs, err := h.mobileService.GetMobileConfigs(c.Request().Context(), platform, companyID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to get mobile configurations")
	}

	return response.SuccessWithMessage(c, "Mobile configurations retrieved successfully", configs)
}

func (h *MobileHandler) SetMobileConfig(c echo.Context) error {
	var config models.MobileConfiguration
	if err := c.Bind(&config); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body")
	}

	userID := getUserIDFromContext(c)
	if err := h.mobileService.SetMobileConfig(c.Request().Context(), &config, userID); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to set mobile configuration")
	}

	return response.SuccessWithMessage(c, "Mobile configuration set successfully", config)
}

// Business Operations endpoints
func (h *MobileHandler) GetMobileStatistics(c echo.Context) error {
	companyID := getCompanyIDFromContext(c)

	stats, err := h.mobileService.GetMobileStatistics(c.Request().Context(), companyID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to get mobile statistics")
	}

	return response.SuccessWithMessage(c, "Mobile statistics retrieved successfully", stats)
}

func (h *MobileHandler) GetDeviceUsageStats(c echo.Context) error {
	deviceIDStr := c.Param("deviceId")
	deviceID, err := uuid.Parse(deviceIDStr)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid device ID")
	}

	days := 30
	if daysStr := c.QueryParam("days"); daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil {
			days = d
		}
	}

	stats, err := h.mobileService.GetDeviceUsageStats(c.Request().Context(), deviceID, days)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to get device usage stats")
	}

	return response.SuccessWithMessage(c, "Device usage statistics retrieved successfully", stats)
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
		return response.Error(c, http.StatusInternalServerError, "Failed to get notification stats")
	}

	return response.SuccessWithMessage(c, "Notification statistics retrieved successfully", stats)
}

func (h *MobileHandler) GenerateMobileDashboard(c echo.Context) error {
	companyID := *getCompanyIDFromContext(c)

	dashboard, err := h.mobileService.GenerateMobileDashboard(c.Request().Context(), companyID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to generate mobile dashboard")
	}

	return response.SuccessWithMessage(c, "Mobile dashboard generated successfully", dashboard)
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