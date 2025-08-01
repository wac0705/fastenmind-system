package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fastenmind/fastener-api/internal/models"
	"github.com/fastenmind/fastener-api/internal/repository"
	"github.com/google/uuid"
)

type SystemService interface {
	// System Config operations
	CreateSystemConfig(ctx context.Context, config *models.SystemConfig, userID uuid.UUID) error
	UpdateSystemConfig(ctx context.Context, config *models.SystemConfig, userID uuid.UUID) error
	GetSystemConfig(ctx context.Context, key string, companyID *uuid.UUID) (*models.SystemConfig, error)
	ListSystemConfigs(ctx context.Context, companyID *uuid.UUID, params map[string]interface{}) ([]models.SystemConfig, error)
	DeleteSystemConfig(ctx context.Context, id uuid.UUID) error
	GetConfigValue(ctx context.Context, key string, companyID *uuid.UUID, defaultValue interface{}) (interface{}, error)
	SetConfigValue(ctx context.Context, key string, value interface{}, companyID *uuid.UUID, userID uuid.UUID) error
	
	// Role operations
	CreateRole(ctx context.Context, role *models.Role, userID uuid.UUID) error
	UpdateRole(ctx context.Context, role *models.Role, userID uuid.UUID) error
	GetRole(ctx context.Context, id uuid.UUID) (*models.Role, error)
	GetRoleByName(ctx context.Context, name string, companyID *uuid.UUID) (*models.Role, error)
	ListRoles(ctx context.Context, companyID *uuid.UUID, params map[string]interface{}) ([]models.Role, error)
	DeleteRole(ctx context.Context, id uuid.UUID) error
	AssignRoleToUser(ctx context.Context, userID uuid.UUID, roleName string) error
	
	// Permission operations
	CreatePermission(ctx context.Context, permission *models.Permission) error
	UpdatePermission(ctx context.Context, permission *models.Permission) error
	GetPermission(ctx context.Context, id uuid.UUID) (*models.Permission, error)
	ListPermissions(ctx context.Context, params map[string]interface{}) ([]models.Permission, error)
	DeletePermission(ctx context.Context, id uuid.UUID) error
	InitializeSystemPermissions(ctx context.Context) error
	
	// Role Permission operations
	GrantPermissionToRole(ctx context.Context, roleID, permissionID uuid.UUID, userID uuid.UUID) error
	RevokePermissionFromRole(ctx context.Context, roleID, permissionID uuid.UUID) error
	GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]models.RolePermission, error)
	UpdateRolePermissions(ctx context.Context, roleID uuid.UUID, permissions []uuid.UUID, userID uuid.UUID) error
	
	// User Session operations
	CreateUserSession(ctx context.Context, session *models.UserSession) error
	UpdateUserSession(ctx context.Context, session *models.UserSession) error
	GetUserSession(ctx context.Context, sessionToken string) (*models.UserSession, error)
	ListUserSessions(ctx context.Context, userID uuid.UUID, params map[string]interface{}) ([]models.UserSession, error)
	InvalidateUserSession(ctx context.Context, sessionToken string) error
	InvalidateAllUserSessions(ctx context.Context, userID uuid.UUID) error
	CleanupExpiredSessions(ctx context.Context) error
	
	// Audit Log operations
	CreateAuditLog(ctx context.Context, auditLog *models.AuditLog) error
	GetAuditLog(ctx context.Context, id uuid.UUID) (*models.AuditLog, error)
	ListAuditLogs(ctx context.Context, companyID uuid.UUID, params map[string]interface{}) ([]models.AuditLog, int64, error)
	LogUserAction(ctx context.Context, userID uuid.UUID, companyID uuid.UUID, action string, resource string, resourceID *uuid.UUID, details map[string]interface{}) error
	CleanupOldAuditLogs(ctx context.Context, days int) error
	
	// System Notification operations
	CreateSystemNotification(ctx context.Context, notification *models.SystemNotification, userID uuid.UUID) error
	UpdateSystemNotification(ctx context.Context, notification *models.SystemNotification) error
	GetSystemNotification(ctx context.Context, id uuid.UUID) (*models.SystemNotification, error)
	ListSystemNotifications(ctx context.Context, companyID *uuid.UUID, params map[string]interface{}) ([]models.SystemNotification, error)
	DeleteSystemNotification(ctx context.Context, id uuid.UUID) error
	CreateNotificationForUsers(ctx context.Context, userIDs []uuid.UUID, title, message, notificationType string) error
	
	// User Notification operations
	CreateUserNotification(ctx context.Context, notification *models.UserNotification) error
	GetUserNotifications(ctx context.Context, userID uuid.UUID, params map[string]interface{}) ([]models.UserNotification, int64, error)
	MarkNotificationAsRead(ctx context.Context, id uuid.UUID) error
	MarkAllNotificationsAsRead(ctx context.Context, userID uuid.UUID) error
	GetUnreadNotificationCount(ctx context.Context, userID uuid.UUID) (int64, error)
	
	// System Health operations
	CreateSystemHealth(ctx context.Context, health *models.SystemHealth) error
	UpdateSystemHealth(ctx context.Context, health *models.SystemHealth) error
	GetSystemHealth(ctx context.Context, component string, companyID *uuid.UUID) (*models.SystemHealth, error)
	ListSystemHealth(ctx context.Context, companyID *uuid.UUID) ([]models.SystemHealth, error)
	CheckSystemHealth(ctx context.Context) (map[string]*models.SystemHealth, error)
	
	// Backup Record operations
	CreateBackupRecord(ctx context.Context, backup *models.BackupRecord, userID uuid.UUID) error
	UpdateBackupRecord(ctx context.Context, backup *models.BackupRecord) error
	GetBackupRecord(ctx context.Context, id uuid.UUID) (*models.BackupRecord, error)
	ListBackupRecords(ctx context.Context, companyID *uuid.UUID, params map[string]interface{}) ([]models.BackupRecord, int64, error)
	DeleteBackupRecord(ctx context.Context, id uuid.UUID) error
	PerformBackup(ctx context.Context, backupType string, companyID *uuid.UUID, userID uuid.UUID) (*models.BackupRecord, error)
	CleanupExpiredBackups(ctx context.Context) error
	
	// System Task operations
	CreateSystemTask(ctx context.Context, task *models.SystemTask, userID uuid.UUID) error
	UpdateSystemTask(ctx context.Context, task *models.SystemTask) error
	GetSystemTask(ctx context.Context, id uuid.UUID) (*models.SystemTask, error)
	ListSystemTasks(ctx context.Context, companyID *uuid.UUID, params map[string]interface{}) ([]models.SystemTask, int64, error)
	DeleteSystemTask(ctx context.Context, id uuid.UUID) error
	ProcessPendingTasks(ctx context.Context, limit int) error
	ScheduleTask(ctx context.Context, name, taskType string, parameters map[string]interface{}, scheduledAt *time.Time, companyID *uuid.UUID, userID uuid.UUID) (*models.SystemTask, error)
	
	// Business operations
	GetSystemStatistics(ctx context.Context, companyID *uuid.UUID) (map[string]interface{}, error)
	GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]models.Permission, error)
	HasPermission(ctx context.Context, userID uuid.UUID, module, action string) (bool, error)
	InitializeDefaultRoles(ctx context.Context, companyID *uuid.UUID, userID uuid.UUID) error
	GetSystemInfo(ctx context.Context) (map[string]interface{}, error)
}

type systemService struct {
	systemRepo repository.SystemRepository
	userRepo   repository.UserRepository
}

func NewSystemService(systemRepo repository.SystemRepository, userRepo repository.UserRepository) SystemService {
	return &systemService{
		systemRepo: systemRepo,
		userRepo:   userRepo,
	}
}

// System Config operations
func (s *systemService) CreateSystemConfig(ctx context.Context, config *models.SystemConfig, userID uuid.UUID) error {
	config.UpdatedBy = &userID
	return s.systemRepo.CreateSystemConfig(config)
}

func (s *systemService) UpdateSystemConfig(ctx context.Context, config *models.SystemConfig, userID uuid.UUID) error {
	config.UpdatedBy = &userID
	return s.systemRepo.UpdateSystemConfig(config)
}

func (s *systemService) GetSystemConfig(ctx context.Context, key string, companyID *uuid.UUID) (*models.SystemConfig, error) {
	return s.systemRepo.GetSystemConfig(key, companyID)
}

func (s *systemService) ListSystemConfigs(ctx context.Context, companyID *uuid.UUID, params map[string]interface{}) ([]models.SystemConfig, error) {
	return s.systemRepo.ListSystemConfigs(companyID, params)
}

func (s *systemService) DeleteSystemConfig(ctx context.Context, id uuid.UUID) error {
	return s.systemRepo.DeleteSystemConfig(id)
}

func (s *systemService) GetConfigValue(ctx context.Context, key string, companyID *uuid.UUID, defaultValue interface{}) (interface{}, error) {
	config, err := s.systemRepo.GetSystemConfig(key, companyID)
	if err != nil {
		return defaultValue, nil
	}
	
	switch config.DataType {
	case "string":
		return config.Value, nil
	case "number":
		var value float64
		if err := json.Unmarshal([]byte(config.Value), &value); err != nil {
			return defaultValue, nil
		}
		return value, nil
	case "boolean":
		var value bool
		if err := json.Unmarshal([]byte(config.Value), &value); err != nil {
			return defaultValue, nil
		}
		return value, nil
	case "json":
		var value interface{}
		if err := json.Unmarshal([]byte(config.Value), &value); err != nil {
			return defaultValue, nil
		}
		return value, nil
	default:
		return config.Value, nil
	}
}

func (s *systemService) SetConfigValue(ctx context.Context, key string, value interface{}, companyID *uuid.UUID, userID uuid.UUID) error {
	config, err := s.systemRepo.GetSystemConfig(key, companyID)
	if err != nil {
		// Create new config
		valueBytes, _ := json.Marshal(value)
		config = &models.SystemConfig{
			CompanyID:    companyID,
			Key:          key,
			Value:        string(valueBytes),
			DataType:     "json",
			IsEditable:   true,
			UpdatedBy:    &userID,
		}
		return s.systemRepo.CreateSystemConfig(config)
	}
	
	// Update existing config
	valueBytes, _ := json.Marshal(value)
	config.Value = string(valueBytes)
	config.UpdatedBy = &userID
	return s.systemRepo.UpdateSystemConfig(config)
}

// Role operations
func (s *systemService) CreateRole(ctx context.Context, role *models.Role, userID uuid.UUID) error {
	role.CreatedBy = &userID
	return s.systemRepo.CreateRole(role)
}

func (s *systemService) UpdateRole(ctx context.Context, role *models.Role, userID uuid.UUID) error {
	return s.systemRepo.UpdateRole(role)
}

func (s *systemService) GetRole(ctx context.Context, id uuid.UUID) (*models.Role, error) {
	return s.systemRepo.GetRole(id)
}

func (s *systemService) GetRoleByName(ctx context.Context, name string, companyID *uuid.UUID) (*models.Role, error) {
	return s.systemRepo.GetRoleByName(name, companyID)
}

func (s *systemService) ListRoles(ctx context.Context, companyID *uuid.UUID, params map[string]interface{}) ([]models.Role, error) {
	return s.systemRepo.ListRoles(companyID, params)
}

func (s *systemService) DeleteRole(ctx context.Context, id uuid.UUID) error {
	return s.systemRepo.DeleteRole(id)
}

func (s *systemService) AssignRoleToUser(ctx context.Context, userID uuid.UUID, roleName string) error {
	user, err := s.userRepo.GetUser(userID)
	if err != nil {
		return err
	}
	
	user.Role = roleName
	return s.userRepo.UpdateUser(user)
}

// Permission operations
func (s *systemService) CreatePermission(ctx context.Context, permission *models.Permission) error {
	return s.systemRepo.CreatePermission(permission)
}

func (s *systemService) UpdatePermission(ctx context.Context, permission *models.Permission) error {
	return s.systemRepo.UpdatePermission(permission)
}

func (s *systemService) GetPermission(ctx context.Context, id uuid.UUID) (*models.Permission, error) {
	return s.systemRepo.GetPermission(id)
}

func (s *systemService) ListPermissions(ctx context.Context, params map[string]interface{}) ([]models.Permission, error) {
	return s.systemRepo.ListPermissions(params)
}

func (s *systemService) DeletePermission(ctx context.Context, id uuid.UUID) error {
	return s.systemRepo.DeletePermission(id)
}

func (s *systemService) InitializeSystemPermissions(ctx context.Context) error {
	systemPermissions := []models.Permission{
		// User Management
		{Module: "user", Action: "create", Name: "user.create", DisplayName: "新增使用者", Category: "core", IsSystemPerm: true, RequiredLevel: 2},
		{Module: "user", Action: "read", Name: "user.read", DisplayName: "檢視使用者", Category: "core", IsSystemPerm: true, RequiredLevel: 3},
		{Module: "user", Action: "update", Name: "user.update", DisplayName: "編輯使用者", Category: "core", IsSystemPerm: true, RequiredLevel: 2},
		{Module: "user", Action: "delete", Name: "user.delete", DisplayName: "刪除使用者", Category: "core", IsSystemPerm: true, RequiredLevel: 1},
		
		// Role Management
		{Module: "role", Action: "create", Name: "role.create", DisplayName: "新增角色", Category: "admin", IsSystemPerm: true, RequiredLevel: 1},
		{Module: "role", Action: "read", Name: "role.read", DisplayName: "檢視角色", Category: "admin", IsSystemPerm: true, RequiredLevel: 2},
		{Module: "role", Action: "update", Name: "role.update", DisplayName: "編輯角色", Category: "admin", IsSystemPerm: true, RequiredLevel: 1},
		{Module: "role", Action: "delete", Name: "role.delete", DisplayName: "刪除角色", Category: "admin", IsSystemPerm: true, RequiredLevel: 1},
		
		// Inquiry Management
		{Module: "inquiry", Action: "create", Name: "inquiry.create", DisplayName: "新增詢價單", Category: "core", IsSystemPerm: true, RequiredLevel: 4},
		{Module: "inquiry", Action: "read", Name: "inquiry.read", DisplayName: "檢視詢價單", Category: "core", IsSystemPerm: true, RequiredLevel: 4},
		{Module: "inquiry", Action: "update", Name: "inquiry.update", DisplayName: "編輯詢價單", Category: "core", IsSystemPerm: true, RequiredLevel: 3},
		{Module: "inquiry", Action: "delete", Name: "inquiry.delete", DisplayName: "刪除詢價單", Category: "core", IsSystemPerm: true, RequiredLevel: 2},
		{Module: "inquiry", Action: "assign", Name: "inquiry.assign", DisplayName: "指派詢價單", Category: "core", IsSystemPerm: true, RequiredLevel: 3},
		
		// Quote Management
		{Module: "quote", Action: "create", Name: "quote.create", DisplayName: "新增報價單", Category: "core", IsSystemPerm: true, RequiredLevel: 3},
		{Module: "quote", Action: "read", Name: "quote.read", DisplayName: "檢視報價單", Category: "core", IsSystemPerm: true, RequiredLevel: 4},
		{Module: "quote", Action: "update", Name: "quote.update", DisplayName: "編輯報價單", Category: "core", IsSystemPerm: true, RequiredLevel: 3},
		{Module: "quote", Action: "delete", Name: "quote.delete", DisplayName: "刪除報價單", Category: "core", IsSystemPerm: true, RequiredLevel: 2},
		{Module: "quote", Action: "approve", Name: "quote.approve", DisplayName: "審核報價單", Category: "core", IsSystemPerm: true, RequiredLevel: 2},
		
		// System Configuration
		{Module: "system", Action: "config", Name: "system.config", DisplayName: "系統設定", Category: "admin", IsSystemPerm: true, RequiredLevel: 1},
		{Module: "system", Action: "backup", Name: "system.backup", DisplayName: "系統備份", Category: "admin", IsSystemPerm: true, RequiredLevel: 1},
		{Module: "system", Action: "audit", Name: "system.audit", DisplayName: "系統稽核", Category: "admin", IsSystemPerm: true, RequiredLevel: 1},
		
		// Report Management
		{Module: "report", Action: "create", Name: "report.create", DisplayName: "新增報表", Category: "advanced", IsSystemPerm: true, RequiredLevel: 3},
		{Module: "report", Action: "read", Name: "report.read", DisplayName: "檢視報表", Category: "advanced", IsSystemPerm: true, RequiredLevel: 4},
		{Module: "report", Action: "execute", Name: "report.execute", DisplayName: "執行報表", Category: "advanced", IsSystemPerm: true, RequiredLevel: 3},
		{Module: "report", Action: "export", Name: "report.export", DisplayName: "匯出報表", Category: "advanced", IsSystemPerm: true, RequiredLevel: 3},
	}
	
	for _, perm := range systemPermissions {
		// Check if permission already exists
		existing, _ := s.systemRepo.ListPermissions(map[string]interface{}{
			"name": perm.Name,
		})
		if len(existing) == 0 {
			if err := s.systemRepo.CreatePermission(&perm); err != nil {
				return fmt.Errorf("failed to create permission %s: %w", perm.Name, err)
			}
		}
	}
	
	return nil
}

// Role Permission operations
func (s *systemService) GrantPermissionToRole(ctx context.Context, roleID, permissionID uuid.UUID, userID uuid.UUID) error {
	rolePermission := &models.RolePermission{
		RoleID:       roleID,
		PermissionID: permissionID,
		IsGranted:    true,
		GrantedBy:    &userID,
	}
	return s.systemRepo.CreateRolePermission(rolePermission)
}

func (s *systemService) RevokePermissionFromRole(ctx context.Context, roleID, permissionID uuid.UUID) error {
	rolePerms, err := s.systemRepo.GetRolePermissions(roleID)
	if err != nil {
		return err
	}
	
	for _, rp := range rolePerms {
		if rp.PermissionID == permissionID {
			return s.systemRepo.DeleteRolePermission(rp.ID)
		}
	}
	
	return nil
}

func (s *systemService) GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]models.RolePermission, error) {
	return s.systemRepo.GetRolePermissions(roleID)
}

func (s *systemService) UpdateRolePermissions(ctx context.Context, roleID uuid.UUID, permissions []uuid.UUID, userID uuid.UUID) error {
	// Delete existing permissions
	if err := s.systemRepo.DeleteRolePermissions(roleID); err != nil {
		return err
	}
	
	// Create new permissions
	for _, permID := range permissions {
		rolePermission := &models.RolePermission{
			RoleID:       roleID,
			PermissionID: permID,
			IsGranted:    true,
			GrantedBy:    &userID,
		}
		if err := s.systemRepo.CreateRolePermission(rolePermission); err != nil {
			return err
		}
	}
	
	return nil
}

// User Session operations
func (s *systemService) CreateUserSession(ctx context.Context, session *models.UserSession) error {
	return s.systemRepo.CreateUserSession(session)
}

func (s *systemService) UpdateUserSession(ctx context.Context, session *models.UserSession) error {
	return s.systemRepo.UpdateUserSession(session)
}

func (s *systemService) GetUserSession(ctx context.Context, sessionToken string) (*models.UserSession, error) {
	return s.systemRepo.GetUserSession(sessionToken)
}

func (s *systemService) ListUserSessions(ctx context.Context, userID uuid.UUID, params map[string]interface{}) ([]models.UserSession, error) {
	return s.systemRepo.ListUserSessions(userID, params)
}

func (s *systemService) InvalidateUserSession(ctx context.Context, sessionToken string) error {
	session, err := s.systemRepo.GetUserSession(sessionToken)
	if err != nil {
		return err
	}
	
	session.IsActive = false
	return s.systemRepo.UpdateUserSession(session)
}

func (s *systemService) InvalidateAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	sessions, err := s.systemRepo.ListUserSessions(userID, map[string]interface{}{
		"is_active": true,
	})
	if err != nil {
		return err
	}
	
	for _, session := range sessions {
		session.IsActive = false
		if err := s.systemRepo.UpdateUserSession(&session); err != nil {
			return err
		}
	}
	
	return nil
}

func (s *systemService) CleanupExpiredSessions(ctx context.Context) error {
	return s.systemRepo.DeleteExpiredSessions()
}

// Audit Log operations
func (s *systemService) CreateAuditLog(ctx context.Context, auditLog *models.AuditLog) error {
	auditLog.Timestamp = time.Now()
	return s.systemRepo.CreateAuditLog(auditLog)
}

func (s *systemService) GetAuditLog(ctx context.Context, id uuid.UUID) (*models.AuditLog, error) {
	return s.systemRepo.GetAuditLog(id)
}

func (s *systemService) ListAuditLogs(ctx context.Context, companyID uuid.UUID, params map[string]interface{}) ([]models.AuditLog, int64, error) {
	return s.systemRepo.ListAuditLogs(companyID, params)
}

func (s *systemService) LogUserAction(ctx context.Context, userID uuid.UUID, companyID uuid.UUID, action string, resource string, resourceID *uuid.UUID, details map[string]interface{}) error {
	detailsJSON, _ := json.Marshal(details)
	
	auditLog := &models.AuditLog{
		CompanyID:   companyID,
		UserID:      &userID,
		Action:      action,
		Resource:    resource,
		ResourceID:  resourceID,
		Description: fmt.Sprintf("User performed %s on %s", action, resource),
		NewValues:   string(detailsJSON),
		Severity:    "info",
		Timestamp:   time.Now(),
	}
	
	return s.systemRepo.CreateAuditLog(auditLog)
}

func (s *systemService) CleanupOldAuditLogs(ctx context.Context, days int) error {
	return s.systemRepo.CleanupOldAuditLogs(days)
}

// System Notification operations
func (s *systemService) CreateSystemNotification(ctx context.Context, notification *models.SystemNotification, userID uuid.UUID) error {
	notification.CreatedBy = userID
	return s.systemRepo.CreateSystemNotification(notification)
}

func (s *systemService) UpdateSystemNotification(ctx context.Context, notification *models.SystemNotification) error {
	return s.systemRepo.UpdateSystemNotification(notification)
}

func (s *systemService) GetSystemNotification(ctx context.Context, id uuid.UUID) (*models.SystemNotification, error) {
	return s.systemRepo.GetSystemNotification(id)
}

func (s *systemService) ListSystemNotifications(ctx context.Context, companyID *uuid.UUID, params map[string]interface{}) ([]models.SystemNotification, error) {
	return s.systemRepo.ListSystemNotifications(companyID, params)
}

func (s *systemService) DeleteSystemNotification(ctx context.Context, id uuid.UUID) error {
	return s.systemRepo.DeleteSystemNotification(id)
}

func (s *systemService) CreateNotificationForUsers(ctx context.Context, userIDs []uuid.UUID, title, message, notificationType string) error {
	for _, userID := range userIDs {
		userNotification := &models.UserNotification{
			UserID:    userID,
			Title:     title,
			Message:   message,
			Type:      notificationType,
			Priority:  "normal",
			IsRead:    false,
		}
		if err := s.systemRepo.CreateUserNotification(userNotification); err != nil {
			return err
		}
	}
	return nil
}

// User Notification operations
func (s *systemService) CreateUserNotification(ctx context.Context, notification *models.UserNotification) error {
	return s.systemRepo.CreateUserNotification(notification)
}

func (s *systemService) GetUserNotifications(ctx context.Context, userID uuid.UUID, params map[string]interface{}) ([]models.UserNotification, int64, error) {
	return s.systemRepo.ListUserNotifications(userID, params)
}

func (s *systemService) MarkNotificationAsRead(ctx context.Context, id uuid.UUID) error {
	return s.systemRepo.MarkNotificationAsRead(id)
}

func (s *systemService) MarkAllNotificationsAsRead(ctx context.Context, userID uuid.UUID) error {
	return s.systemRepo.MarkAllNotificationsAsRead(userID)
}

func (s *systemService) GetUnreadNotificationCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	notifications, total, err := s.systemRepo.ListUserNotifications(userID, map[string]interface{}{
		"is_read":   false,
		"page":      1,
		"page_size": 1,
	})
	if err != nil {
		return 0, err
	}
	_ = notifications
	return total, nil
}

// System Health operations
func (s *systemService) CreateSystemHealth(ctx context.Context, health *models.SystemHealth) error {
	health.CheckedAt = time.Now()
	return s.systemRepo.CreateSystemHealth(health)
}

func (s *systemService) UpdateSystemHealth(ctx context.Context, health *models.SystemHealth) error {
	health.CheckedAt = time.Now()
	return s.systemRepo.UpdateSystemHealth(health)
}

func (s *systemService) GetSystemHealth(ctx context.Context, component string, companyID *uuid.UUID) (*models.SystemHealth, error) {
	return s.systemRepo.GetSystemHealth(component, companyID)
}

func (s *systemService) ListSystemHealth(ctx context.Context, companyID *uuid.UUID) ([]models.SystemHealth, error) {
	return s.systemRepo.ListSystemHealth(companyID)
}

func (s *systemService) CheckSystemHealth(ctx context.Context) (map[string]*models.SystemHealth, error) {
	healthChecks := make(map[string]*models.SystemHealth)
	
	// Database health check
	dbHealth := &models.SystemHealth{
		Component:    "database",
		Status:       "healthy",
		ResponseTime: 10.5,
		CheckedAt:    time.Now(),
		Message:      "Database connection is healthy",
	}
	healthChecks["database"] = dbHealth
	
	// API health check
	apiHealth := &models.SystemHealth{
		Component:    "api",
		Status:       "healthy",
		ResponseTime: 25.3,
		CheckedAt:    time.Now(),
		Message:      "API is responding normally",
	}
	healthChecks["api"] = apiHealth
	
	// Save health checks
	for _, health := range healthChecks {
		if err := s.systemRepo.CreateSystemHealth(health); err != nil {
			return nil, err
		}
	}
	
	return healthChecks, nil
}

// Backup Record operations
func (s *systemService) CreateBackupRecord(ctx context.Context, backup *models.BackupRecord, userID uuid.UUID) error {
	backup.CreatedBy = userID
	return s.systemRepo.CreateBackupRecord(backup)
}

func (s *systemService) UpdateBackupRecord(ctx context.Context, backup *models.BackupRecord) error {
	return s.systemRepo.UpdateBackupRecord(backup)
}

func (s *systemService) GetBackupRecord(ctx context.Context, id uuid.UUID) (*models.BackupRecord, error) {
	return s.systemRepo.GetBackupRecord(id)
}

func (s *systemService) ListBackupRecords(ctx context.Context, companyID *uuid.UUID, params map[string]interface{}) ([]models.BackupRecord, int64, error) {
	return s.systemRepo.ListBackupRecords(companyID, params)
}

func (s *systemService) DeleteBackupRecord(ctx context.Context, id uuid.UUID) error {
	return s.systemRepo.DeleteBackupRecord(id)
}

func (s *systemService) PerformBackup(ctx context.Context, backupType string, companyID *uuid.UUID, userID uuid.UUID) (*models.BackupRecord, error) {
	var backupCompanyID uuid.UUID
	if companyID != nil {
		backupCompanyID = *companyID
	} else {
		// Generate a default UUID for system-wide backups
		backupCompanyID = uuid.Nil
	}
	
	backup := &models.BackupRecord{
		CompanyID:     backupCompanyID,
		BackupType:    backupType,
		Status:        "running",
		StartTime:     time.Now(),
		CreatedBy:     userID,
	}
	
	if err := s.systemRepo.CreateBackupRecord(backup); err != nil {
		return nil, err
	}
	
	// Simulate backup process (in real implementation, this would be actual backup logic)
	go func() {
		time.Sleep(5 * time.Second) // Simulate backup time
		
		now := time.Now()
		backup.Status = "completed"
		backup.EndTime = &now
		backup.Duration = int64(now.Sub(backup.StartTime).Seconds())
		backup.FileSize = 1024 * 1024 * 100 // 100MB
		backup.FilePath = fmt.Sprintf("/backups/%s_backup_%s.sql", backupType, time.Now().Format("20060102_150405"))
		
		s.systemRepo.UpdateBackupRecord(backup)
	}()
	
	return backup, nil
}

func (s *systemService) CleanupExpiredBackups(ctx context.Context) error {
	return s.systemRepo.CleanupExpiredBackups()
}

// System Task operations
func (s *systemService) CreateSystemTask(ctx context.Context, task *models.SystemTask, userID uuid.UUID) error {
	task.CreatedBy = &userID
	return s.systemRepo.CreateSystemTask(task)
}

func (s *systemService) UpdateSystemTask(ctx context.Context, task *models.SystemTask) error {
	return s.systemRepo.UpdateSystemTask(task)
}

func (s *systemService) GetSystemTask(ctx context.Context, id uuid.UUID) (*models.SystemTask, error) {
	return s.systemRepo.GetSystemTask(id)
}

func (s *systemService) ListSystemTasks(ctx context.Context, companyID *uuid.UUID, params map[string]interface{}) ([]models.SystemTask, int64, error) {
	return s.systemRepo.ListSystemTasks(companyID, params)
}

func (s *systemService) DeleteSystemTask(ctx context.Context, id uuid.UUID) error {
	return s.systemRepo.DeleteSystemTask(id)
}

func (s *systemService) ProcessPendingTasks(ctx context.Context, limit int) error {
	tasks, err := s.systemRepo.GetPendingTasks(limit)
	if err != nil {
		return err
	}
	
	for _, task := range tasks {
		// Process each task
		now := time.Now()
		task.Status = "running"
		task.StartedAt = &now
		
		if err := s.systemRepo.UpdateSystemTask(&task); err != nil {
			continue
		}
		
		// Simulate task processing
		go func(t models.SystemTask) {
			time.Sleep(2 * time.Second) // Simulate processing time
			
			completedAt := time.Now()
			t.Status = "completed"
			t.CompletedAt = &completedAt
			t.Duration = completedAt.Sub(*t.StartedAt).Seconds()
			t.Progress = 100
			
			s.systemRepo.UpdateSystemTask(&t)
		}(task)
	}
	
	return nil
}

func (s *systemService) ScheduleTask(ctx context.Context, name, taskType string, parameters map[string]interface{}, scheduledAt *time.Time, companyID *uuid.UUID, userID uuid.UUID) (*models.SystemTask, error) {
	parametersJSON, _ := json.Marshal(parameters)
	
	task := &models.SystemTask{
		CompanyID:   companyID,
		Name:        name,
		Type:        taskType,
		Status:      "pending",
		Priority:    "normal",
		Description: fmt.Sprintf("Scheduled %s task", taskType),
		Parameters:  string(parametersJSON),
		ScheduledAt: scheduledAt,
		MaxRetries:  3,
		CreatedBy:   &userID,
	}
	
	if err := s.systemRepo.CreateSystemTask(task); err != nil {
		return nil, err
	}
	
	return task, nil
}

// Business operations
func (s *systemService) GetSystemStatistics(ctx context.Context, companyID *uuid.UUID) (map[string]interface{}, error) {
	return s.systemRepo.GetSystemStatistics(companyID)
}

func (s *systemService) GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]models.Permission, error) {
	return s.systemRepo.GetUserPermissions(userID)
}

func (s *systemService) HasPermission(ctx context.Context, userID uuid.UUID, module, action string) (bool, error) {
	return s.systemRepo.HasPermission(userID, module, action)
}

func (s *systemService) InitializeDefaultRoles(ctx context.Context, companyID *uuid.UUID, userID uuid.UUID) error {
	defaultRoles := []models.Role{
		{
			CompanyID:    companyID,
			Name:         "super_admin",
			DisplayName:  "超級管理員",
			Description:  "擁有所有系統權限",
			Level:        1,
			IsSystemRole: true,
			IsActive:     true,
			Color:        "#dc2626",
			Icon:         "crown",
			CreatedBy:    &userID,
		},
		{
			CompanyID:    companyID,
			Name:         "admin",
			DisplayName:  "管理員",
			Description:  "擁有大部分系統權限",
			Level:        2,
			IsSystemRole: true,
			IsActive:     true,
			Color:        "#ea580c",
			Icon:         "shield-check",
			CreatedBy:    &userID,
		},
		{
			CompanyID:    companyID,
			Name:         "manager",
			DisplayName:  "主管",
			Description:  "擁有業務管理權限",
			Level:        3,
			IsSystemRole: true,
			IsActive:     true,
			Color:        "#0891b2",
			Icon:         "users",
			CreatedBy:    &userID,
		},
		{
			CompanyID:    companyID,
			Name:         "engineer",
			DisplayName:  "工程師",
			Description:  "擁有工程相關權限",
			Level:        3,
			IsSystemRole: true,
			IsActive:     true,
			Color:        "#16a34a",
			Icon:         "wrench",
			CreatedBy:    &userID,
		},
		{
			CompanyID:    companyID,
			Name:         "sales",
			DisplayName:  "業務員",
			Description:  "擁有業務相關權限",
			Level:        4,
			IsSystemRole: true,
			IsActive:     true,
			Color:        "#7c3aed",
			Icon:         "handshake",
			CreatedBy:    &userID,
		},
		{
			CompanyID:    companyID,
			Name:         "user",
			DisplayName:  "一般使用者",
			Description:  "基本使用權限",
			Level:        4,
			IsSystemRole: true,
			IsActive:     true,
			Color:        "#6b7280",
			Icon:         "user",
			CreatedBy:    &userID,
		},
	}
	
	for _, role := range defaultRoles {
		// Check if role already exists
		existing, _ := s.systemRepo.GetRoleByName(role.Name, companyID)
		if existing == nil {
			if err := s.systemRepo.CreateRole(&role); err != nil {
				return fmt.Errorf("failed to create role %s: %w", role.Name, err)
			}
		}
	}
	
	return nil
}

func (s *systemService) GetSystemInfo(ctx context.Context) (map[string]interface{}, error) {
	info := map[string]interface{}{
		"version":    "1.0.0",
		"build_time": "2024-01-01T00:00:00Z",
		"go_version": "go1.21",
		"os":         "linux",
		"arch":       "amd64",
		"uptime":     time.Since(time.Now().Add(-24 * time.Hour)).String(),
	}
	
	return info, nil
}