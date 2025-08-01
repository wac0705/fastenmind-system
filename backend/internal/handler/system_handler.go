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

type SystemHandler struct {
	systemService service.SystemService
}

func NewSystemHandler(systemService service.SystemService) *SystemHandler {
	return &SystemHandler{
		systemService: systemService,
	}
}

// System Config endpoints
func (h *SystemHandler) CreateSystemConfig(c echo.Context) error {
	var config models.SystemConfig
	if err := c.Bind(&config); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body", err)
	}

	userID := getUserIDFromContext(c)
	if err := h.systemService.CreateSystemConfig(c.Request().Context(), &config, userID); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to create system config", err)
	}

	return response.Success(c, config, "System config created successfully")
}

func (h *SystemHandler) UpdateSystemConfig(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid config ID", err)
	}

	var config models.SystemConfig
	if err := c.Bind(&config); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body", err)
	}

	config.ID = id
	userID := getUserIDFromContext(c)
	if err := h.systemService.UpdateSystemConfig(c.Request().Context(), &config, userID); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to update system config", err)
	}

	return response.Success(c, config, "System config updated successfully")
}

func (h *SystemHandler) GetSystemConfig(c echo.Context) error {
	key := c.Param("key")
	companyID := getCompanyIDFromContext(c)

	config, err := h.systemService.GetSystemConfig(c.Request().Context(), key, companyID)
	if err != nil {
		return response.Error(c, http.StatusNotFound, "System config not found", err)
	}

	return response.Success(c, config, "System config retrieved successfully")
}

func (h *SystemHandler) ListSystemConfigs(c echo.Context) error {
	companyID := getCompanyIDFromContext(c)
	params := getListParams(c)

	configs, err := h.systemService.ListSystemConfigs(c.Request().Context(), companyID, params)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to list system configs", err)
	}

	return response.Success(c, map[string]interface{}{
		"data":  configs,
		"total": len(configs),
	}, "System configs retrieved successfully")
}

func (h *SystemHandler) DeleteSystemConfig(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid config ID", err)
	}

	if err := h.systemService.DeleteSystemConfig(c.Request().Context(), id); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to delete system config", err)
	}

	return response.Success(c, nil, "System config deleted successfully")
}

func (h *SystemHandler) GetConfigValue(c echo.Context) error {
	key := c.Param("key")
	companyID := getCompanyIDFromContext(c)

	value, err := h.systemService.GetConfigValue(c.Request().Context(), key, companyID, nil)
	if err != nil {
		return response.Error(c, http.StatusNotFound, "Config value not found", err)
	}

	return response.Success(c, map[string]interface{}{
		"key":   key,
		"value": value,
	}, "Config value retrieved successfully")
}

func (h *SystemHandler) SetConfigValue(c echo.Context) error {
	key := c.Param("key")
	companyID := getCompanyIDFromContext(c)
	userID := getUserIDFromContext(c)

	var req struct {
		Value interface{} `json:"value"`
	}
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body", err)
	}

	if err := h.systemService.SetConfigValue(c.Request().Context(), key, req.Value, companyID, userID); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to set config value", err)
	}

	return response.Success(c, map[string]interface{}{
		"key":   key,
		"value": req.Value,
	}, "Config value set successfully")
}

// Role endpoints
func (h *SystemHandler) CreateRole(c echo.Context) error {
	var role models.Role
	if err := c.Bind(&role); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body", err)
	}

	userID := getUserIDFromContext(c)
	if err := h.systemService.CreateRole(c.Request().Context(), &role, userID); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to create role", err)
	}

	return response.Success(c, role, "Role created successfully")
}

func (h *SystemHandler) UpdateRole(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid role ID", err)
	}

	var role models.Role
	if err := c.Bind(&role); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body", err)
	}

	role.ID = id
	userID := getUserIDFromContext(c)
	if err := h.systemService.UpdateRole(c.Request().Context(), &role, userID); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to update role", err)
	}

	return response.Success(c, role, "Role updated successfully")
}

func (h *SystemHandler) GetRole(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid role ID", err)
	}

	role, err := h.systemService.GetRole(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, http.StatusNotFound, "Role not found", err)
	}

	return response.Success(c, role, "Role retrieved successfully")
}

func (h *SystemHandler) ListRoles(c echo.Context) error {
	companyID := getCompanyIDFromContext(c)
	params := getListParams(c)

	roles, err := h.systemService.ListRoles(c.Request().Context(), companyID, params)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to list roles", err)
	}

	return response.Success(c, map[string]interface{}{
		"data":  roles,
		"total": len(roles),
	}, "Roles retrieved successfully")
}

func (h *SystemHandler) DeleteRole(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid role ID", err)
	}

	if err := h.systemService.DeleteRole(c.Request().Context(), id); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to delete role", err)
	}

	return response.Success(c, nil, "Role deleted successfully")
}

// Permission endpoints
func (h *SystemHandler) CreatePermission(c echo.Context) error {
	var permission models.Permission
	if err := c.Bind(&permission); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body", err)
	}

	if err := h.systemService.CreatePermission(c.Request().Context(), &permission); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to create permission", err)
	}

	return response.Success(c, permission, "Permission created successfully")
}

func (h *SystemHandler) ListPermissions(c echo.Context) error {
	params := getListParams(c)

	permissions, err := h.systemService.ListPermissions(c.Request().Context(), params)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to list permissions", err)
	}

	return response.Success(c, map[string]interface{}{
		"data":  permissions,
		"total": len(permissions),
	}, "Permissions retrieved successfully")
}

func (h *SystemHandler) InitializeSystemPermissions(c echo.Context) error {
	if err := h.systemService.InitializeSystemPermissions(c.Request().Context()); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to initialize system permissions", err)
	}

	return response.Success(c, nil, "System permissions initialized successfully")
}

// Role Permission endpoints
func (h *SystemHandler) GetRolePermissions(c echo.Context) error {
	roleID, err := uuid.Parse(c.Param("roleId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid role ID", err)
	}

	permissions, err := h.systemService.GetRolePermissions(c.Request().Context(), roleID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to get role permissions", err)
	}

	return response.Success(c, permissions, "Role permissions retrieved successfully")
}

func (h *SystemHandler) UpdateRolePermissions(c echo.Context) error {
	roleID, err := uuid.Parse(c.Param("roleId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid role ID", err)
	}

	var req struct {
		Permissions []uuid.UUID `json:"permissions"`
	}
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body", err)
	}

	userID := getUserIDFromContext(c)
	if err := h.systemService.UpdateRolePermissions(c.Request().Context(), roleID, req.Permissions, userID); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to update role permissions", err)
	}

	return response.Success(c, nil, "Role permissions updated successfully")
}

// User Session endpoints
func (h *SystemHandler) ListUserSessions(c echo.Context) error {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid user ID", err)
	}

	params := getListParams(c)
	sessions, err := h.systemService.ListUserSessions(c.Request().Context(), userID, params)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to list user sessions", err)
	}

	return response.Success(c, map[string]interface{}{
		"data":  sessions,
		"total": len(sessions),
	}, "User sessions retrieved successfully")
}

func (h *SystemHandler) InvalidateUserSession(c echo.Context) error {
	sessionToken := c.Param("token")

	if err := h.systemService.InvalidateUserSession(c.Request().Context(), sessionToken); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to invalidate session", err)
	}

	return response.Success(c, nil, "Session invalidated successfully")
}

func (h *SystemHandler) InvalidateAllUserSessions(c echo.Context) error {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid user ID", err)
	}

	if err := h.systemService.InvalidateAllUserSessions(c.Request().Context(), userID); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to invalidate all sessions", err)
	}

	return response.Success(c, nil, "All user sessions invalidated successfully")
}

// Audit Log endpoints
func (h *SystemHandler) ListAuditLogs(c echo.Context) error {
	companyID := getCompanyIDFromContext(c)
	params := getListParams(c)

	logs, total, err := h.systemService.ListAuditLogs(c.Request().Context(), companyID, params)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to list audit logs", err)
	}

	return response.Success(c, map[string]interface{}{
		"data":  logs,
		"total": total,
	}, "Audit logs retrieved successfully")
}

func (h *SystemHandler) GetAuditLog(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid audit log ID", err)
	}

	log, err := h.systemService.GetAuditLog(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, http.StatusNotFound, "Audit log not found", err)
	}

	return response.Success(c, log, "Audit log retrieved successfully")
}

// System Notification endpoints
func (h *SystemHandler) CreateSystemNotification(c echo.Context) error {
	var notification models.SystemNotification
	if err := c.Bind(&notification); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body", err)
	}

	userID := getUserIDFromContext(c)
	if err := h.systemService.CreateSystemNotification(c.Request().Context(), &notification, userID); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to create system notification", err)
	}

	return response.Success(c, notification, "System notification created successfully")
}

func (h *SystemHandler) ListSystemNotifications(c echo.Context) error {
	companyID := getCompanyIDFromContext(c)
	params := getListParams(c)

	notifications, err := h.systemService.ListSystemNotifications(c.Request().Context(), companyID, params)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to list system notifications", err)
	}

	return response.Success(c, map[string]interface{}{
		"data":  notifications,
		"total": len(notifications),
	}, "System notifications retrieved successfully")
}

// User Notification endpoints
func (h *SystemHandler) GetUserNotifications(c echo.Context) error {
	userID := getUserIDFromContext(c)
	params := getListParams(c)

	notifications, total, err := h.systemService.GetUserNotifications(c.Request().Context(), userID, params)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to get user notifications", err)
	}

	return response.Success(c, map[string]interface{}{
		"data":  notifications,
		"total": total,
	}, "User notifications retrieved successfully")
}

func (h *SystemHandler) MarkNotificationAsRead(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid notification ID", err)
	}

	if err := h.systemService.MarkNotificationAsRead(c.Request().Context(), id); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to mark notification as read", err)
	}

	return response.Success(c, nil, "Notification marked as read successfully")
}

func (h *SystemHandler) MarkAllNotificationsAsRead(c echo.Context) error {
	userID := getUserIDFromContext(c)

	if err := h.systemService.MarkAllNotificationsAsRead(c.Request().Context(), userID); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to mark all notifications as read", err)
	}

	return response.Success(c, nil, "All notifications marked as read successfully")
}

func (h *SystemHandler) GetUnreadNotificationCount(c echo.Context) error {
	userID := getUserIDFromContext(c)

	count, err := h.systemService.GetUnreadNotificationCount(c.Request().Context(), userID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to get unread notification count", err)
	}

	return response.Success(c, map[string]interface{}{
		"count": count,
	}, "Unread notification count retrieved successfully")
}

// System Health endpoints
func (h *SystemHandler) ListSystemHealth(c echo.Context) error {
	companyID := getCompanyIDFromContext(c)

	healthList, err := h.systemService.ListSystemHealth(c.Request().Context(), companyID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to list system health", err)
	}

	return response.Success(c, healthList, "System health retrieved successfully")
}

func (h *SystemHandler) CheckSystemHealth(c echo.Context) error {
	healthChecks, err := h.systemService.CheckSystemHealth(c.Request().Context())
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to check system health", err)
	}

	return response.Success(c, healthChecks, "System health check completed successfully")
}

// Backup Record endpoints
func (h *SystemHandler) ListBackupRecords(c echo.Context) error {
	companyID := getCompanyIDFromContext(c)
	params := getListParams(c)

	backups, total, err := h.systemService.ListBackupRecords(c.Request().Context(), companyID, params)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to list backup records", err)
	}

	return response.Success(c, map[string]interface{}{
		"data":  backups,
		"total": total,
	}, "Backup records retrieved successfully")
}

func (h *SystemHandler) PerformBackup(c echo.Context) error {
	var req struct {
		Type string `json:"type"`
	}
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body", err)
	}

	companyID := getCompanyIDFromContext(c)
	userID := getUserIDFromContext(c)

	backup, err := h.systemService.PerformBackup(c.Request().Context(), req.Type, companyID, userID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to perform backup", err)
	}

	return response.Success(c, backup, "Backup started successfully")
}

// System Task endpoints
func (h *SystemHandler) ListSystemTasks(c echo.Context) error {
	companyID := getCompanyIDFromContext(c)
	params := getListParams(c)

	tasks, total, err := h.systemService.ListSystemTasks(c.Request().Context(), companyID, params)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to list system tasks", err)
	}

	return response.Success(c, map[string]interface{}{
		"data":  tasks,
		"total": total,
	}, "System tasks retrieved successfully")
}

func (h *SystemHandler) ScheduleTask(c echo.Context) error {
	var req struct {
		Name        string                 `json:"name"`
		Type        string                 `json:"type"`
		Parameters  map[string]interface{} `json:"parameters"`
		ScheduledAt *time.Time             `json:"scheduled_at"`
	}
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body", err)
	}

	companyID := getCompanyIDFromContext(c)
	userID := getUserIDFromContext(c)

	task, err := h.systemService.ScheduleTask(c.Request().Context(), req.Name, req.Type, req.Parameters, req.ScheduledAt, companyID, userID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to schedule task", err)
	}

	return response.Success(c, task, "Task scheduled successfully")
}

func (h *SystemHandler) ProcessPendingTasks(c echo.Context) error {
	limit := 10
	if limitStr := c.QueryParam("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	if err := h.systemService.ProcessPendingTasks(c.Request().Context(), limit); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to process pending tasks", err)
	}

	return response.Success(c, nil, "Pending tasks processed successfully")
}

// Business operations
func (h *SystemHandler) GetSystemStatistics(c echo.Context) error {
	companyID := getCompanyIDFromContext(c)

	stats, err := h.systemService.GetSystemStatistics(c.Request().Context(), companyID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to get system statistics", err)
	}

	return response.Success(c, stats, "System statistics retrieved successfully")
}

func (h *SystemHandler) GetUserPermissions(c echo.Context) error {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		userID = getUserIDFromContext(c) // Use current user if no userID provided
	}

	permissions, err := h.systemService.GetUserPermissions(c.Request().Context(), userID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to get user permissions", err)
	}

	return response.Success(c, permissions, "User permissions retrieved successfully")
}

func (h *SystemHandler) CheckUserPermission(c echo.Context) error {
	module := c.QueryParam("module")
	action := c.QueryParam("action")
	userID := getUserIDFromContext(c)

	if module == "" || action == "" {
		return response.Error(c, http.StatusBadRequest, "Module and action are required", nil)
	}

	hasPermission, err := h.systemService.HasPermission(c.Request().Context(), userID, module, action)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to check permission", err)
	}

	return response.Success(c, map[string]interface{}{
		"has_permission": hasPermission,
		"module":         module,
		"action":         action,
	}, "Permission check completed successfully")
}

func (h *SystemHandler) InitializeDefaultRoles(c echo.Context) error {
	companyID := getCompanyIDFromContext(c)
	userID := getUserIDFromContext(c)

	if err := h.systemService.InitializeDefaultRoles(c.Request().Context(), companyID, userID); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to initialize default roles", err)
	}

	return response.Success(c, nil, "Default roles initialized successfully")
}

func (h *SystemHandler) GetSystemInfo(c echo.Context) error {
	info, err := h.systemService.GetSystemInfo(c.Request().Context())
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to get system info", err)
	}

	return response.Success(c, info, "System info retrieved successfully")
}

// Note: getUserIDFromContext and getCompanyIDFromContext are defined in common.go

func getListParams(c echo.Context) map[string]interface{} {
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
	
	if category := c.QueryParam("category"); category != "" {
		params["category"] = category
	}
	
	if status := c.QueryParam("status"); status != "" {
		params["status"] = status
	}
	
	if isActive := c.QueryParam("is_active"); isActive != "" {
		if active, err := strconv.ParseBool(isActive); err == nil {
			params["is_active"] = active
		}
	}
	
	if level := c.QueryParam("level"); level != "" {
		if l, err := strconv.Atoi(level); err == nil {
			params["level"] = l
		}
	}
	
	if module := c.QueryParam("module"); module != "" {
		params["module"] = module
	}
	
	if action := c.QueryParam("action"); action != "" {
		params["action"] = action
	}
	
	if severity := c.QueryParam("severity"); severity != "" {
		params["severity"] = severity
	}
	
	if startDate := c.QueryParam("start_date"); startDate != "" {
		params["start_date"] = startDate
	}
	
	if endDate := c.QueryParam("end_date"); endDate != "" {
		params["end_date"] = endDate
	}
	
	if userID := c.QueryParam("user_id"); userID != "" {
		params["user_id"] = userID
	}
	
	if ipAddress := c.QueryParam("ip_address"); ipAddress != "" {
		params["ip_address"] = ipAddress
	}
	
	if notificationType := c.QueryParam("type"); notificationType != "" {
		params["type"] = notificationType
	}
	
	if priority := c.QueryParam("priority"); priority != "" {
		params["priority"] = priority
	}
	
	if isRead := c.QueryParam("is_read"); isRead != "" {
		if read, err := strconv.ParseBool(isRead); err == nil {
			params["is_read"] = read
		}
	}
	
	if isDismissed := c.QueryParam("is_dismissed"); isDismissed != "" {
		if dismissed, err := strconv.ParseBool(isDismissed); err == nil {
			params["is_dismissed"] = dismissed
		}
	}
	
	return params
}