package repository

import (
	"github.com/fastenmind/fastener-api/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SystemRepository interface {
	// System Config operations
	CreateSystemConfig(config *models.SystemConfig) error
	UpdateSystemConfig(config *models.SystemConfig) error
	GetSystemConfig(key string, companyID *uuid.UUID) (*models.SystemConfig, error)
	ListSystemConfigs(companyID *uuid.UUID, params map[string]interface{}) ([]models.SystemConfig, error)
	DeleteSystemConfig(id uuid.UUID) error
	
	// Role operations
	CreateRole(role *models.Role) error
	UpdateRole(role *models.Role) error
	GetRole(id uuid.UUID) (*models.Role, error)
	GetRoleByName(name string, companyID *uuid.UUID) (*models.Role, error)
	ListRoles(companyID *uuid.UUID, params map[string]interface{}) ([]models.Role, error)
	DeleteRole(id uuid.UUID) error
	
	// Permission operations
	CreatePermission(permission *models.Permission) error
	UpdatePermission(permission *models.Permission) error
	GetPermission(id uuid.UUID) (*models.Permission, error)
	ListPermissions(params map[string]interface{}) ([]models.Permission, error)
	DeletePermission(id uuid.UUID) error
	
	// Role Permission operations
	CreateRolePermission(rolePermission *models.RolePermission) error
	UpdateRolePermission(rolePermission *models.RolePermission) error
	GetRolePermissions(roleID uuid.UUID) ([]models.RolePermission, error)
	DeleteRolePermission(id uuid.UUID) error
	DeleteRolePermissions(roleID uuid.UUID) error
	
	// User Session operations
	CreateUserSession(session *models.UserSession) error
	UpdateUserSession(session *models.UserSession) error
	GetUserSession(sessionToken string) (*models.UserSession, error)
	ListUserSessions(userID uuid.UUID, params map[string]interface{}) ([]models.UserSession, error)
	DeleteUserSession(id uuid.UUID) error
	DeleteExpiredSessions() error
	
	// Audit Log operations
	CreateAuditLog(auditLog *models.AuditLog) error
	GetAuditLog(id uuid.UUID) (*models.AuditLog, error)
	ListAuditLogs(companyID uuid.UUID, params map[string]interface{}) ([]models.AuditLog, int64, error)
	DeleteAuditLog(id uuid.UUID) error
	CleanupOldAuditLogs(days int) error
	
	// System Notification operations
	CreateSystemNotification(notification *models.SystemNotification) error
	UpdateSystemNotification(notification *models.SystemNotification) error
	GetSystemNotification(id uuid.UUID) (*models.SystemNotification, error)
	ListSystemNotifications(companyID *uuid.UUID, params map[string]interface{}) ([]models.SystemNotification, error)
	DeleteSystemNotification(id uuid.UUID) error
	
	// User Notification operations
	CreateUserNotification(notification *models.UserNotification) error
	UpdateUserNotification(notification *models.UserNotification) error
	GetUserNotification(id uuid.UUID) (*models.UserNotification, error)
	ListUserNotifications(userID uuid.UUID, params map[string]interface{}) ([]models.UserNotification, int64, error)
	DeleteUserNotification(id uuid.UUID) error
	MarkNotificationAsRead(id uuid.UUID) error
	MarkAllNotificationsAsRead(userID uuid.UUID) error
	
	// System Health operations
	CreateSystemHealth(health *models.SystemHealth) error
	UpdateSystemHealth(health *models.SystemHealth) error
	GetSystemHealth(component string, companyID *uuid.UUID) (*models.SystemHealth, error)
	ListSystemHealth(companyID *uuid.UUID) ([]models.SystemHealth, error)
	DeleteSystemHealth(id uuid.UUID) error
	
	// Backup Record operations
	CreateBackupRecord(backup *models.BackupRecord) error
	UpdateBackupRecord(backup *models.BackupRecord) error
	GetBackupRecord(id uuid.UUID) (*models.BackupRecord, error)
	ListBackupRecords(companyID *uuid.UUID, params map[string]interface{}) ([]models.BackupRecord, int64, error)
	DeleteBackupRecord(id uuid.UUID) error
	CleanupExpiredBackups() error
	
	// System Task operations
	CreateSystemTask(task *models.SystemTask) error
	UpdateSystemTask(task *models.SystemTask) error
	GetSystemTask(id uuid.UUID) (*models.SystemTask, error)
	ListSystemTasks(companyID *uuid.UUID, params map[string]interface{}) ([]models.SystemTask, int64, error)
	DeleteSystemTask(id uuid.UUID) error
	GetPendingTasks(limit int) ([]models.SystemTask, error)
	
	// Business operations
	GetSystemStatistics(companyID *uuid.UUID) (map[string]interface{}, error)
	GetUserPermissions(userID uuid.UUID) ([]models.Permission, error)
	HasPermission(userID uuid.UUID, module, action string) (bool, error)
}

type systemRepository struct {
	db *gorm.DB
}

func NewSystemRepository(db interface{}) SystemRepository {
	gormDB, ok := db.(*gorm.DB)
	if !ok {
		panic("invalid database type, expected *gorm.DB")
	}
	return &systemRepository{db: gormDB}
}

// System Config operations
func (r *systemRepository) CreateSystemConfig(config *models.SystemConfig) error {
	return r.db.Create(config).Error
}

func (r *systemRepository) UpdateSystemConfig(config *models.SystemConfig) error {
	return r.db.Save(config).Error
}

func (r *systemRepository) GetSystemConfig(key string, companyID *uuid.UUID) (*models.SystemConfig, error) {
	var config models.SystemConfig
	query := r.db.Where("key = ?", key)
	
	if companyID != nil {
		query = query.Where("company_id = ? OR company_id IS NULL", *companyID).Order("company_id DESC")
	} else {
		query = query.Where("company_id IS NULL")
	}
	
	err := query.First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *systemRepository) ListSystemConfigs(companyID *uuid.UUID, params map[string]interface{}) ([]models.SystemConfig, error) {
	var configs []models.SystemConfig
	
	query := r.db.Model(&models.SystemConfig{})
	
	if companyID != nil {
		query = query.Where("company_id = ? OR company_id IS NULL", *companyID)
	} else {
		query = query.Where("company_id IS NULL")
	}
	
	// Apply filters
	if category, ok := params["category"].(string); ok && category != "" {
		query = query.Where("category = ?", category)
	}
	
	if editable, ok := params["is_editable"].(bool); ok {
		query = query.Where("is_editable = ?", editable)
	}
	
	if search, ok := params["search"].(string); ok && search != "" {
		query = query.Where("key LIKE ? OR display_name LIKE ? OR description LIKE ?", 
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	
	err := query.Order("category ASC, display_order ASC, key ASC").Find(&configs).Error
	return configs, err
}

func (r *systemRepository) DeleteSystemConfig(id uuid.UUID) error {
	return r.db.Delete(&models.SystemConfig{}, id).Error
}

// Role operations
func (r *systemRepository) CreateRole(role *models.Role) error {
	return r.db.Create(role).Error
}

func (r *systemRepository) UpdateRole(role *models.Role) error {
	return r.db.Save(role).Error
}

func (r *systemRepository) GetRole(id uuid.UUID) (*models.Role, error) {
	var role models.Role
	err := r.db.Preload("Company").
		Preload("Creator").
		Preload("Permissions").
		Preload("Permissions.Permission").
		First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *systemRepository) GetRoleByName(name string, companyID *uuid.UUID) (*models.Role, error) {
	var role models.Role
	query := r.db.Where("name = ?", name)
	
	if companyID != nil {
		query = query.Where("company_id = ? OR is_system_role = ?", *companyID, true)
	} else {
		query = query.Where("is_system_role = ?", true)
	}
	
	err := query.Preload("Permissions").
		Preload("Permissions.Permission").
		First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *systemRepository) ListRoles(companyID *uuid.UUID, params map[string]interface{}) ([]models.Role, error) {
	var roles []models.Role
	
	query := r.db.Model(&models.Role{})
	
	if companyID != nil {
		query = query.Where("company_id = ? OR is_system_role = ?", *companyID, true)
	} else {
		query = query.Where("is_system_role = ?", true)
	}
	
	// Apply filters
	if isActive, ok := params["is_active"].(bool); ok {
		query = query.Where("is_active = ?", isActive)
	}
	
	if isSystemRole, ok := params["is_system_role"].(bool); ok {
		query = query.Where("is_system_role = ?", isSystemRole)
	}
	
	if level, ok := params["level"].(int); ok && level > 0 {
		query = query.Where("level = ?", level)
	}
	
	if search, ok := params["search"].(string); ok && search != "" {
		query = query.Where("name LIKE ? OR display_name LIKE ? OR description LIKE ?", 
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	
	err := query.Order("level ASC, name ASC").
		Preload("Creator").
		Find(&roles).Error
	return roles, err
}

func (r *systemRepository) DeleteRole(id uuid.UUID) error {
	return r.db.Delete(&models.Role{}, id).Error
}

// Permission operations
func (r *systemRepository) CreatePermission(permission *models.Permission) error {
	return r.db.Create(permission).Error
}

func (r *systemRepository) UpdatePermission(permission *models.Permission) error {
	return r.db.Save(permission).Error
}

func (r *systemRepository) GetPermission(id uuid.UUID) (*models.Permission, error) {
	var permission models.Permission
	err := r.db.First(&permission, id).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *systemRepository) ListPermissions(params map[string]interface{}) ([]models.Permission, error) {
	var permissions []models.Permission
	
	query := r.db.Model(&models.Permission{})
	
	// Apply filters
	if module, ok := params["module"].(string); ok && module != "" {
		query = query.Where("module = ?", module)
	}
	
	if action, ok := params["action"].(string); ok && action != "" {
		query = query.Where("action = ?", action)
	}
	
	if category, ok := params["category"].(string); ok && category != "" {
		query = query.Where("category = ?", category)
	}
	
	if isSystemPerm, ok := params["is_system_perm"].(bool); ok {
		query = query.Where("is_system_perm = ?", isSystemPerm)
	}
	
	if search, ok := params["search"].(string); ok && search != "" {
		query = query.Where("name LIKE ? OR display_name LIKE ? OR description LIKE ?", 
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	
	err := query.Order("module ASC, action ASC, name ASC").Find(&permissions).Error
	return permissions, err
}

func (r *systemRepository) DeletePermission(id uuid.UUID) error {
	return r.db.Delete(&models.Permission{}, id).Error
}

// Role Permission operations
func (r *systemRepository) CreateRolePermission(rolePermission *models.RolePermission) error {
	return r.db.Create(rolePermission).Error
}

func (r *systemRepository) UpdateRolePermission(rolePermission *models.RolePermission) error {
	return r.db.Save(rolePermission).Error
}

func (r *systemRepository) GetRolePermissions(roleID uuid.UUID) ([]models.RolePermission, error) {
	var rolePermissions []models.RolePermission
	err := r.db.Where("role_id = ?", roleID).
		Preload("Permission").
		Find(&rolePermissions).Error
	return rolePermissions, err
}

func (r *systemRepository) DeleteRolePermission(id uuid.UUID) error {
	return r.db.Delete(&models.RolePermission{}, id).Error
}

func (r *systemRepository) DeleteRolePermissions(roleID uuid.UUID) error {
	return r.db.Where("role_id = ?", roleID).Delete(&models.RolePermission{}).Error
}

// User Session operations
func (r *systemRepository) CreateUserSession(session *models.UserSession) error {
	return r.db.Create(session).Error
}

func (r *systemRepository) UpdateUserSession(session *models.UserSession) error {
	return r.db.Save(session).Error
}

func (r *systemRepository) GetUserSession(sessionToken string) (*models.UserSession, error) {
	var session models.UserSession
	err := r.db.Where("session_token = ? AND is_active = ? AND expires_at > ?", 
		sessionToken, true, time.Now()).
		Preload("User").
		First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *systemRepository) ListUserSessions(userID uuid.UUID, params map[string]interface{}) ([]models.UserSession, error) {
	var sessions []models.UserSession
	
	query := r.db.Where("user_id = ?", userID)
	
	// Apply filters
	if isActive, ok := params["is_active"].(bool); ok {
		query = query.Where("is_active = ?", isActive)
	}
	
	if loginMethod, ok := params["login_method"].(string); ok && loginMethod != "" {
		query = query.Where("login_method = ?", loginMethod)
	}
	
	err := query.Order("created_at DESC").Find(&sessions).Error
	return sessions, err
}

func (r *systemRepository) DeleteUserSession(id uuid.UUID) error {
	return r.db.Delete(&models.UserSession{}, id).Error
}

func (r *systemRepository) DeleteExpiredSessions() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&models.UserSession{}).Error
}

// Audit Log operations
func (r *systemRepository) CreateAuditLog(auditLog *models.AuditLog) error {
	return r.db.Create(auditLog).Error
}

func (r *systemRepository) GetAuditLog(id uuid.UUID) (*models.AuditLog, error) {
	var auditLog models.AuditLog
	err := r.db.Preload("Company").
		Preload("User").
		Preload("Session").
		First(&auditLog, id).Error
	if err != nil {
		return nil, err
	}
	return &auditLog, nil
}

func (r *systemRepository) ListAuditLogs(companyID uuid.UUID, params map[string]interface{}) ([]models.AuditLog, int64, error) {
	var auditLogs []models.AuditLog
	var total int64
	
	query := r.db.Model(&models.AuditLog{}).Where("company_id = ?", companyID)
	
	// Apply filters
	if action, ok := params["action"].(string); ok && action != "" {
		query = query.Where("action = ?", action)
	}
	
	if resource, ok := params["resource"].(string); ok && resource != "" {
		query = query.Where("resource = ?", resource)
	}
	
	if module, ok := params["module"].(string); ok && module != "" {
		query = query.Where("module = ?", module)
	}
	
	if userID, ok := params["user_id"].(string); ok && userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	
	if severity, ok := params["severity"].(string); ok && severity != "" {
		query = query.Where("severity = ?", severity)
	}
	
	if startDate, ok := params["start_date"].(string); ok && startDate != "" {
		query = query.Where("timestamp >= ?", startDate)
	}
	
	if endDate, ok := params["end_date"].(string); ok && endDate != "" {
		query = query.Where("timestamp <= ?", endDate)
	}
	
	if ipAddress, ok := params["ip_address"].(string); ok && ipAddress != "" {
		query = query.Where("ip_address = ?", ipAddress)
	}
	
	if search, ok := params["search"].(string); ok && search != "" {
		query = query.Where("description LIKE ? OR request_path LIKE ?", 
			"%"+search+"%", "%"+search+"%")
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
	
	// Apply sorting
	query = query.Order("timestamp DESC")
	
	// Load with relations
	if err := query.
		Preload("User").
		Find(&auditLogs).Error; err != nil {
		return nil, 0, err
	}
	
	return auditLogs, total, nil
}

func (r *systemRepository) DeleteAuditLog(id uuid.UUID) error {
	return r.db.Delete(&models.AuditLog{}, id).Error
}

func (r *systemRepository) CleanupOldAuditLogs(days int) error {
	cutoffDate := time.Now().AddDate(0, 0, -days)
	return r.db.Where("timestamp < ?", cutoffDate).Delete(&models.AuditLog{}).Error
}

// System Notification operations
func (r *systemRepository) CreateSystemNotification(notification *models.SystemNotification) error {
	return r.db.Create(notification).Error
}

func (r *systemRepository) UpdateSystemNotification(notification *models.SystemNotification) error {
	return r.db.Save(notification).Error
}

func (r *systemRepository) GetSystemNotification(id uuid.UUID) (*models.SystemNotification, error) {
	var notification models.SystemNotification
	err := r.db.Preload("Company").
		Preload("Creator").
		First(&notification, id).Error
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

func (r *systemRepository) ListSystemNotifications(companyID *uuid.UUID, params map[string]interface{}) ([]models.SystemNotification, error) {
	var notifications []models.SystemNotification
	
	query := r.db.Model(&models.SystemNotification{})
	
	if companyID != nil {
		query = query.Where("company_id = ? OR company_id IS NULL", *companyID)
	} else {
		query = query.Where("company_id IS NULL")
	}
	
	// Apply filters
	if isActive, ok := params["is_active"].(bool); ok {
		query = query.Where("is_active = ?", isActive)
	}
	
	if notificationType, ok := params["type"].(string); ok && notificationType != "" {
		query = query.Where("type = ?", notificationType)
	}
	
	if priority, ok := params["priority"].(string); ok && priority != "" {
		query = query.Where("priority = ?", priority)
	}
	
	// Filter by current time for scheduled notifications
	now := time.Now()
	query = query.Where("(show_from IS NULL OR show_from <= ?) AND (show_until IS NULL OR show_until >= ?)", now, now)
	
	err := query.Order("priority DESC, created_at DESC").
		Preload("Creator").
		Find(&notifications).Error
	return notifications, err
}

func (r *systemRepository) DeleteSystemNotification(id uuid.UUID) error {
	return r.db.Delete(&models.SystemNotification{}, id).Error
}

// User Notification operations
func (r *systemRepository) CreateUserNotification(notification *models.UserNotification) error {
	return r.db.Create(notification).Error
}

func (r *systemRepository) UpdateUserNotification(notification *models.UserNotification) error {
	return r.db.Save(notification).Error
}

func (r *systemRepository) GetUserNotification(id uuid.UUID) (*models.UserNotification, error) {
	var notification models.UserNotification
	err := r.db.Preload("User").
		Preload("SystemNotification").
		First(&notification, id).Error
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

func (r *systemRepository) ListUserNotifications(userID uuid.UUID, params map[string]interface{}) ([]models.UserNotification, int64, error) {
	var notifications []models.UserNotification
	var total int64
	
	query := r.db.Model(&models.UserNotification{}).Where("user_id = ?", userID)
	
	// Apply filters
	if isRead, ok := params["is_read"].(bool); ok {
		query = query.Where("is_read = ?", isRead)
	}
	
	if isDismissed, ok := params["is_dismissed"].(bool); ok {
		query = query.Where("is_dismissed = ?", isDismissed)
	}
	
	if notificationType, ok := params["type"].(string); ok && notificationType != "" {
		query = query.Where("type = ?", notificationType)
	}
	
	if priority, ok := params["priority"].(string); ok && priority != "" {
		query = query.Where("priority = ?", priority)
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
	
	// Apply sorting
	query = query.Order("created_at DESC")
	
	// Load with relations
	if err := query.
		Preload("SystemNotification").
		Find(&notifications).Error; err != nil {
		return nil, 0, err
	}
	
	return notifications, total, nil
}

func (r *systemRepository) DeleteUserNotification(id uuid.UUID) error {
	return r.db.Delete(&models.UserNotification{}, id).Error
}

func (r *systemRepository) MarkNotificationAsRead(id uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&models.UserNotification{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		}).Error
}

func (r *systemRepository) MarkAllNotificationsAsRead(userID uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&models.UserNotification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		}).Error
}

// System Health operations
func (r *systemRepository) CreateSystemHealth(health *models.SystemHealth) error {
	return r.db.Create(health).Error
}

func (r *systemRepository) UpdateSystemHealth(health *models.SystemHealth) error {
	return r.db.Save(health).Error
}

func (r *systemRepository) GetSystemHealth(component string, companyID *uuid.UUID) (*models.SystemHealth, error) {
	var health models.SystemHealth
	query := r.db.Where("component = ?", component)
	
	if companyID != nil {
		query = query.Where("company_id = ?", *companyID)
	} else {
		query = query.Where("company_id IS NULL")
	}
	
	err := query.Order("checked_at DESC").First(&health).Error
	if err != nil {
		return nil, err
	}
	return &health, nil
}

func (r *systemRepository) ListSystemHealth(companyID *uuid.UUID) ([]models.SystemHealth, error) {
	var healthList []models.SystemHealth
	
	query := r.db.Model(&models.SystemHealth{})
	
	if companyID != nil {
		query = query.Where("company_id = ? OR company_id IS NULL", *companyID)
	} else {
		query = query.Where("company_id IS NULL")
	}
	
	err := query.Order("component ASC, checked_at DESC").Find(&healthList).Error
	return healthList, err
}

func (r *systemRepository) DeleteSystemHealth(id uuid.UUID) error {
	return r.db.Delete(&models.SystemHealth{}, id).Error
}

// Backup Record operations
func (r *systemRepository) CreateBackupRecord(backup *models.BackupRecord) error {
	return r.db.Create(backup).Error
}

func (r *systemRepository) UpdateBackupRecord(backup *models.BackupRecord) error {
	return r.db.Save(backup).Error
}

func (r *systemRepository) GetBackupRecord(id uuid.UUID) (*models.BackupRecord, error) {
	var backup models.BackupRecord
	err := r.db.Preload("Company").
		Preload("Creator").
		First(&backup, id).Error
	if err != nil {
		return nil, err
	}
	return &backup, nil
}

func (r *systemRepository) ListBackupRecords(companyID *uuid.UUID, params map[string]interface{}) ([]models.BackupRecord, int64, error) {
	var backups []models.BackupRecord
	var total int64
	
	query := r.db.Model(&models.BackupRecord{})
	
	if companyID != nil {
		query = query.Where("company_id = ? OR company_id IS NULL", *companyID)
	} else {
		query = query.Where("company_id IS NULL")
	}
	
	// Apply filters
	if backupType, ok := params["type"].(string); ok && backupType != "" {
		query = query.Where("type = ?", backupType)
	}
	
	if status, ok := params["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	
	if startDate, ok := params["start_date"].(string); ok && startDate != "" {
		query = query.Where("started_at >= ?", startDate)
	}
	
	if endDate, ok := params["end_date"].(string); ok && endDate != "" {
		query = query.Where("started_at <= ?", endDate)
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
	
	// Apply sorting
	query = query.Order("started_at DESC")
	
	// Load with relations
	if err := query.
		Preload("Creator").
		Find(&backups).Error; err != nil {
		return nil, 0, err
	}
	
	return backups, total, nil
}

func (r *systemRepository) DeleteBackupRecord(id uuid.UUID) error {
	return r.db.Delete(&models.BackupRecord{}, id).Error
}

func (r *systemRepository) CleanupExpiredBackups() error {
	now := time.Now()
	return r.db.Where("expires_at < ?", now).Delete(&models.BackupRecord{}).Error
}

// System Task operations
func (r *systemRepository) CreateSystemTask(task *models.SystemTask) error {
	return r.db.Create(task).Error
}

func (r *systemRepository) UpdateSystemTask(task *models.SystemTask) error {
	return r.db.Save(task).Error
}

func (r *systemRepository) GetSystemTask(id uuid.UUID) (*models.SystemTask, error) {
	var task models.SystemTask
	err := r.db.Preload("Company").
		Preload("Creator").
		First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *systemRepository) ListSystemTasks(companyID *uuid.UUID, params map[string]interface{}) ([]models.SystemTask, int64, error) {
	var tasks []models.SystemTask
	var total int64
	
	query := r.db.Model(&models.SystemTask{})
	
	if companyID != nil {
		query = query.Where("company_id = ? OR company_id IS NULL", *companyID)
	} else {
		query = query.Where("company_id IS NULL")
	}
	
	// Apply filters
	if taskType, ok := params["type"].(string); ok && taskType != "" {
		query = query.Where("type = ?", taskType)
	}
	
	if status, ok := params["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	
	if priority, ok := params["priority"].(string); ok && priority != "" {
		query = query.Where("priority = ?", priority)
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
	
	// Apply sorting
	query = query.Order("priority DESC, created_at DESC")
	
	// Load with relations
	if err := query.
		Preload("Creator").
		Find(&tasks).Error; err != nil {
		return nil, 0, err
	}
	
	return tasks, total, nil
}

func (r *systemRepository) DeleteSystemTask(id uuid.UUID) error {
	return r.db.Delete(&models.SystemTask{}, id).Error
}

func (r *systemRepository) GetPendingTasks(limit int) ([]models.SystemTask, error) {
	var tasks []models.SystemTask
	err := r.db.Where("status = ? AND (scheduled_at IS NULL OR scheduled_at <= ?)", 
		"pending", time.Now()).
		Order("priority DESC, created_at ASC").
		Limit(limit).
		Find(&tasks).Error
	return tasks, err
}

// Business operations
func (r *systemRepository) GetSystemStatistics(companyID *uuid.UUID) (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	
	// Count users by role
	var userStats []struct {
		Role  string `json:"role"`
		Count int64  `json:"count"`
	}
	userQuery := r.db.Model(&models.User{})
	if companyID != nil {
		userQuery = userQuery.Where("company_id = ?", *companyID)
	}
	if err := userQuery.Select("role, COUNT(*) as count").
		Group("role").
		Find(&userStats).Error; err != nil {
		return nil, err
	}
	stats["users_by_role"] = userStats
	
	// Count audit logs by action
	var auditStats []struct {
		Action string `json:"action"`
		Count  int64  `json:"count"`
	}
	auditQuery := r.db.Model(&models.AuditLog{})
	if companyID != nil {
		auditQuery = auditQuery.Where("company_id = ?", *companyID)
	}
	if err := auditQuery.Where("timestamp >= ?", time.Now().AddDate(0, 0, -30)).
		Select("action, COUNT(*) as count").
		Group("action").
		Find(&auditStats).Error; err != nil {
		return nil, err
	}
	stats["audit_logs_by_action"] = auditStats
	
	// Count total records
	var totalUsers, totalSessions, totalNotifications, totalTasks int64
	
	userCountQuery := r.db.Model(&models.User{})
	if companyID != nil {
		userCountQuery = userCountQuery.Where("company_id = ?", *companyID)
	}
	userCountQuery.Count(&totalUsers)
	
	sessionCountQuery := r.db.Model(&models.UserSession{}).Joins("JOIN users ON users.id = user_sessions.user_id")
	if companyID != nil {
		sessionCountQuery = sessionCountQuery.Where("users.company_id = ?", *companyID)
	}
	sessionCountQuery.Where("user_sessions.is_active = ?", true).Count(&totalSessions)
	
	notificationCountQuery := r.db.Model(&models.UserNotification{}).Joins("JOIN users ON users.id = user_notifications.user_id")
	if companyID != nil {
		notificationCountQuery = notificationCountQuery.Where("users.company_id = ?", *companyID)
	}
	notificationCountQuery.Where("user_notifications.is_read = ?", false).Count(&totalNotifications)
	
	taskCountQuery := r.db.Model(&models.SystemTask{})
	if companyID != nil {
		taskCountQuery = taskCountQuery.Where("company_id = ? OR company_id IS NULL", *companyID)
	}
	taskCountQuery.Where("status IN ?", []string{"pending", "running"}).Count(&totalTasks)
	
	stats["total_users"] = totalUsers
	stats["active_sessions"] = totalSessions
	stats["unread_notifications"] = totalNotifications
	stats["pending_tasks"] = totalTasks
	
	return stats, nil
}

func (r *systemRepository) GetUserPermissions(userID uuid.UUID) ([]models.Permission, error) {
	var permissions []models.Permission
	
	err := r.db.Table("permissions").
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("JOIN roles ON role_permissions.role_id = roles.id").
		Joins("JOIN users ON users.role = roles.name").
		Where("users.id = ? AND role_permissions.is_granted = ? AND roles.is_active = ?", 
			userID, true, true).
		Find(&permissions).Error
	
	return permissions, err
}

func (r *systemRepository) HasPermission(userID uuid.UUID, module, action string) (bool, error) {
	var count int64
	
	err := r.db.Table("permissions").
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("JOIN roles ON role_permissions.role_id = roles.id").
		Joins("JOIN users ON users.role = roles.name").
		Where("users.id = ? AND permissions.module = ? AND permissions.action = ? AND role_permissions.is_granted = ? AND roles.is_active = ?", 
			userID, module, action, true, true).
		Count(&count).Error
	
	return count > 0, err
}