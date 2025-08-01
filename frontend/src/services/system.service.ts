import api from './api'

// Types
export interface User {
  id: string
  username: string
  email: string
  full_name: string
  full_name_en?: string
  avatar?: string
  employee_no?: string
  department?: string
  position?: string
  phone?: string
  mobile?: string
  is_active: boolean
  is_superuser: boolean
  language: string
  timezone: string
  theme: string
  notification_settings?: any
  last_login?: string
  last_password_change?: string
  login_attempts: number
  locked_until?: string
  created_at: string
  updated_at: string
  companies?: Company[]
  roles?: Role[]
}

export interface Role {
  id: string
  company_id?: string
  name: string
  name_en?: string
  code: string
  description?: string
  is_system_role: boolean
  is_active: boolean
  permissions?: Permission[]
  created_at: string
  updated_at: string
  company?: Company
}

export interface Permission {
  id: string
  name: string
  name_en?: string
  code: string
  resource: string
  action: string
  description?: string
  is_system_permission: boolean
  created_at: string
  updated_at: string
}

export interface Company {
  id: string
  code: string
  name: string
  name_en?: string
  description?: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface SystemConfig {
  id: string
  company_id?: string
  category: string
  key: string
  value: any
  value_type: string
  description?: string
  is_encrypted: boolean
  is_public: boolean
  validation_rules?: any
  created_at: string
  updated_at: string
  created_by?: string
  updated_by?: string
}

export interface AuditLog {
  id: string
  company_id?: string
  user_id: string
  user_name: string
  action: string
  resource: string
  resource_id?: string
  old_value?: any
  new_value?: any
  ip_address: string
  user_agent: string
  session_id?: string
  result: 'success' | 'failure'
  error_message?: string
  created_at: string
  user?: User
  company?: Company
}

export interface SystemNotification {
  id: string
  user_id?: string
  company_id?: string
  type: string
  level: 'info' | 'warning' | 'error' | 'success'
  title: string
  title_en?: string
  content: string
  content_en?: string
  data?: any
  is_read: boolean
  read_at?: string
  action_url?: string
  expires_at?: string
  created_at: string
  created_by?: string
  user?: User
  company?: Company
}

export interface UserSession {
  id: string
  user_id: string
  token: string
  refresh_token?: string
  ip_address: string
  user_agent: string
  device_info?: any
  location?: string
  is_active: boolean
  last_activity?: string
  expires_at: string
  created_at: string
  user?: User
}

export interface UserStatistics {
  total_users: number
  active_users: number
  inactive_users: number
  admin_users: number
  new_users_this_month: number
  avg_session_duration: number
  users_by_role: { role: string; count: number }[]
  users_by_department: { department: string; count: number }[]
}

export interface OnlineUser {
  id: string
  username: string
  full_name: string
  avatar?: string
  last_activity: string
  ip_address: string
  location?: string
  device_type: string
}

export interface CreateUserRequest {
  username: string
  email: string
  password: string
  full_name: string
  full_name_en?: string
  employee_no?: string
  department?: string
  position?: string
  phone?: string
  mobile?: string
  role_ids?: string[]
  company_ids?: string[]
  is_active?: boolean
  language?: string
  timezone?: string
  notification_settings?: any
}

export interface UpdateUserRequest {
  email?: string
  full_name?: string
  full_name_en?: string
  employee_no?: string
  department?: string
  position?: string
  phone?: string
  mobile?: string
  role_ids?: string[]
  company_ids?: string[]
  is_active?: boolean
  language?: string
  timezone?: string
  theme?: string
  notification_settings?: any
}

export interface CreateRoleRequest {
  name: string
  name_en?: string
  code: string
  description?: string
  permission_ids?: string[]
  is_active?: boolean
}

export interface UpdateRoleRequest {
  name?: string
  name_en?: string
  description?: string
  permission_ids?: string[]
  is_active?: boolean
}

export interface CreateSystemConfigRequest {
  category: string
  key: string
  value: any
  value_type?: string
  description?: string
  is_encrypted?: boolean
  is_public?: boolean
  validation_rules?: any
}

export interface UpdateSystemConfigRequest {
  value?: any
  description?: string
  is_encrypted?: boolean
  is_public?: boolean
  validation_rules?: any
}

export interface CreateSystemNotificationRequest {
  user_id?: string
  company_id?: string
  type: string
  level: string
  title: string
  title_en?: string
  content: string
  content_en?: string
  data?: any
  action_url?: string
  expires_at?: string
}

export interface ListParams {
  page?: number
  page_size?: number
  search?: string
  sort_by?: string
  sort_order?: string
  company_id?: string
  is_active?: boolean
  role_id?: string
  department?: string
  is_system_role?: boolean
  resource?: string
  action?: string
  category?: string
  is_public?: boolean
  user_id?: string
  start_date?: string
  end_date?: string
  result?: string
  type?: string
  level?: string
  is_read?: boolean
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
}

class SystemService {
  // User Management
  async createUser(data: CreateUserRequest): Promise<User> {
    const response = await api.post('/system/users', data)
    return response.data
  }

  async updateUser(id: string, data: UpdateUserRequest): Promise<User> {
    const response = await api.put(`/system/users/${id}`, data)
    return response.data
  }

  async getUser(id: string): Promise<User> {
    const response = await api.get(`/system/users/${id}`)
    return response.data
  }

  async listUsers(params: ListParams = {}): Promise<PaginatedResponse<User>> {
    const response = await api.get('/system/users', { params })
    return response.data
  }

  async deleteUser(id: string): Promise<void> {
    await api.delete(`/system/users/${id}`)
  }

  async resetUserPassword(id: string, password: string): Promise<void> {
    await api.post(`/system/users/${id}/reset-password`, { password })
  }

  async unlockUser(id: string): Promise<void> {
    await api.post(`/system/users/${id}/unlock`)
  }

  async getUserStatistics(): Promise<UserStatistics> {
    const response = await api.get('/system/users/statistics')
    return response.data
  }

  async getOnlineUsers(): Promise<OnlineUser[]> {
    const response = await api.get('/system/users/online')
    return response.data
  }

  async exportUsers(params: any = {}): Promise<Blob> {
    const response = await api.get('/system/users/export', {
      params,
      responseType: 'blob'
    })
    return response.data
  }

  // Role Management
  async createRole(data: CreateRoleRequest): Promise<Role> {
    const response = await api.post('/system/roles', data)
    return response.data
  }

  async updateRole(id: string, data: UpdateRoleRequest): Promise<Role> {
    const response = await api.put(`/system/roles/${id}`, data)
    return response.data
  }

  async getRole(id: string): Promise<Role> {
    const response = await api.get(`/system/roles/${id}`)
    return response.data
  }

  async listRoles(params: ListParams = {}): Promise<PaginatedResponse<Role>> {
    const response = await api.get('/system/roles', { params })
    return response.data
  }

  async deleteRole(id: string): Promise<void> {
    await api.delete(`/system/roles/${id}`)
  }

  async duplicateRole(id: string): Promise<Role> {
    const response = await api.post(`/system/roles/${id}/duplicate`)
    return response.data
  }

  async assignPermissionsToRole(roleId: string, permissionIds: string[]): Promise<void> {
    await api.post(`/system/roles/${roleId}/permissions`, { permission_ids: permissionIds })
  }

  // Permission Management
  async getPermission(id: string): Promise<Permission> {
    const response = await api.get(`/system/permissions/${id}`)
    return response.data
  }

  async listPermissions(params: ListParams = {}): Promise<PaginatedResponse<Permission>> {
    const response = await api.get('/system/permissions', { params })
    return response.data
  }

  async getPermissionsByResource(resource: string): Promise<Permission[]> {
    const response = await api.get(`/system/permissions/by-resource/${resource}`)
    return response.data
  }

  async checkPermission(resource: string, action: string): Promise<boolean> {
    const response = await api.post('/system/permissions/check', { resource, action })
    return response.data.has_permission
  }

  // System Configuration
  async createSystemConfig(data: CreateSystemConfigRequest): Promise<SystemConfig> {
    const response = await api.post('/system/configs', data)
    return response.data
  }

  async updateSystemConfig(id: string, data: UpdateSystemConfigRequest): Promise<SystemConfig> {
    const response = await api.put(`/system/configs/${id}`, data)
    return response.data
  }

  async getSystemConfig(id: string): Promise<SystemConfig> {
    const response = await api.get(`/system/configs/${id}`)
    return response.data
  }

  async listSystemConfigs(params: ListParams = {}): Promise<PaginatedResponse<SystemConfig>> {
    const response = await api.get('/system/configs', { params })
    return response.data
  }

  async deleteSystemConfig(id: string): Promise<void> {
    await api.delete(`/system/configs/${id}`)
  }

  async getSystemConfigByKey(category: string, key: string): Promise<SystemConfig> {
    const response = await api.get(`/system/configs/${category}/${key}`)
    return response.data
  }

  async updateSystemConfigByKey(category: string, key: string, value: any): Promise<SystemConfig> {
    const response = await api.put(`/system/configs/${category}/${key}`, { value })
    return response.data
  }

  async exportSystemConfigs(): Promise<Blob> {
    const response = await api.get('/system/configs/export', {
      responseType: 'blob'
    })
    return response.data
  }

  async importSystemConfigs(file: File): Promise<any> {
    const formData = new FormData()
    formData.append('file', file)
    
    const response = await api.post('/system/configs/import', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
    
    return response.data
  }

  // Audit Log
  async listAuditLogs(params: ListParams = {}): Promise<PaginatedResponse<AuditLog>> {
    const response = await api.get('/system/audit-logs', { params })
    return response.data
  }

  async getAuditLog(id: string): Promise<AuditLog> {
    const response = await api.get(`/system/audit-logs/${id}`)
    return response.data
  }

  async exportAuditLogs(params: any = {}): Promise<Blob> {
    const response = await api.get('/system/audit-logs/export', {
      params,
      responseType: 'blob'
    })
    return response.data
  }

  // System Notifications
  async createSystemNotification(data: CreateSystemNotificationRequest): Promise<SystemNotification> {
    const response = await api.post('/system/notifications', data)
    return response.data
  }

  async getSystemNotification(id: string): Promise<SystemNotification> {
    const response = await api.get(`/system/notifications/${id}`)
    return response.data
  }

  async listSystemNotifications(params: ListParams = {}): Promise<PaginatedResponse<SystemNotification>> {
    const response = await api.get('/system/notifications', { params })
    return response.data
  }

  async markNotificationAsRead(id: string): Promise<void> {
    await api.post(`/system/notifications/${id}/read`)
  }

  async markAllNotificationsAsRead(): Promise<void> {
    await api.post('/system/notifications/read-all')
  }

  async deleteSystemNotification(id: string): Promise<void> {
    await api.delete(`/system/notifications/${id}`)
  }

  async getUnreadNotificationCount(): Promise<number> {
    const response = await api.get('/system/notifications/unread-count')
    return response.data.count
  }

  // User Sessions
  async listUserSessions(params: ListParams = {}): Promise<PaginatedResponse<UserSession>> {
    const response = await api.get('/system/sessions', { params })
    return response.data
  }

  async getUserSession(id: string): Promise<UserSession> {
    const response = await api.get(`/system/sessions/${id}`)
    return response.data
  }

  async terminateUserSession(id: string): Promise<void> {
    await api.post(`/system/sessions/${id}/terminate`)
  }

  async terminateAllUserSessions(userId: string): Promise<void> {
    await api.post(`/system/users/${userId}/terminate-sessions`)
  }

  // System Information
  async getSystemInfo(): Promise<any> {
    const response = await api.get('/system/info')
    return response.data
  }

  async getSystemHealth(): Promise<any> {
    const response = await api.get('/system/health')
    return response.data
  }

  async getSystemMetrics(): Promise<any> {
    const response = await api.get('/system/metrics')
    return response.data
  }

  async runSystemDiagnostics(): Promise<any> {
    const response = await api.post('/system/diagnostics')
    return response.data
  }

  // Backup and Restore
  async createBackup(options: any = {}): Promise<any> {
    const response = await api.post('/system/backup', options)
    return response.data
  }

  async listBackups(): Promise<any[]> {
    const response = await api.get('/system/backups')
    return response.data
  }

  async restoreBackup(backupId: string): Promise<any> {
    const response = await api.post(`/system/backups/${backupId}/restore`)
    return response.data
  }

  async deleteBackup(backupId: string): Promise<void> {
    await api.delete(`/system/backups/${backupId}`)
  }

  // Cache Management
  async clearCache(cacheType?: string): Promise<void> {
    await api.post('/system/cache/clear', { cache_type: cacheType })
  }

  async getCacheStatistics(): Promise<any> {
    const response = await api.get('/system/cache/statistics')
    return response.data
  }

  // License Management
  async getLicenseInfo(): Promise<any> {
    const response = await api.get('/system/license')
    return response.data
  }

  async updateLicense(licenseKey: string): Promise<any> {
    const response = await api.post('/system/license', { license_key: licenseKey })
    return response.data
  }

  async validateLicense(): Promise<boolean> {
    const response = await api.post('/system/license/validate')
    return response.data.is_valid
  }
}

export default new SystemService()