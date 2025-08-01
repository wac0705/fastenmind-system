import api from './api'

// Types
export interface MobileDevice {
  id: string
  user_id: string
  company_id: string
  device_token: string
  platform: 'ios' | 'android' | 'web'
  device_type: 'phone' | 'tablet' | 'desktop'
  device_model: string
  os_version: string
  app_version: string
  device_name: string
  device_id: string
  bundle_id: string
  time_zone: string
  language: string
  country: string
  push_enabled: boolean
  badge_count: number
  notification_types: string
  is_active: boolean
  last_seen: string
  last_location: string
  is_jailbroken: boolean
  security_level: 'low' | 'normal' | 'high'
  registered_at: string
  updated_at: string
  
  // Relations
  user?: any
  company?: any
  push_notifications?: PushNotification[]
}

export interface PushNotification {
  id: string
  device_id: string
  user_id: string
  company_id: string
  title: string
  body: string
  type: string
  category: string
  data: string
  resource_id?: string
  resource_type: string
  priority: 'low' | 'normal' | 'high' | 'urgent'
  badge: number
  sound: string
  icon: string
  image: string
  actions: string
  status: 'pending' | 'sent' | 'delivered' | 'failed' | 'clicked'
  sent_at?: string
  delivered_at?: string
  clicked_at?: string
  provider_response: string
  error_message: string
  retry_count: number
  max_retries: number
  expires_at?: string
  ttl: number
  created_at: string
  updated_at: string
  scheduled_at?: string
  
  // Relations
  device?: MobileDevice
  user?: any
  company?: any
}

export interface MobileSession {
  id: string
  device_id: string
  user_id: string
  company_id: string
  session_token: string
  refresh_token: string
  start_time: string
  end_time?: string
  duration: number
  start_location: string
  end_location: string
  screen_views: number
  api_requests: number
  data_transferred: number
  app_state: 'active' | 'background' | 'terminated'
  last_activity: string
  network_type: string
  network_provider: string
  created_at: string
  updated_at: string
  
  // Relations
  device?: MobileDevice
  user?: any
  company?: any
}

export interface MobileAnalytics {
  id: string
  device_id: string
  user_id: string
  company_id: string
  event_type: string
  event_name: string
  event_category: string
  event_data: string
  screen_name: string
  screen_class: string
  event_timestamp: string
  session_id: string
  duration: number
  interaction_type: string
  element_id: string
  element_type: string
  load_time: number
  response_time: number
  error_message: string
  app_version: string
  os_version: string
  network_type: string
  battery_level: number
  memory_usage: number
  latitude: number
  longitude: number
  location_accuracy: number
  created_at: string
  
  // Relations
  device?: MobileDevice
  user?: any
  company?: any
  session?: MobileSession
}

export interface MobileAppVersion {
  id: string
  company_id?: string
  version: string
  build_number: string
  platform: 'ios' | 'android' | 'web'
  release_type: 'alpha' | 'beta' | 'production'
  release_notes: string
  release_notes_en: string
  status: 'draft' | 'review' | 'approved' | 'released' | 'deprecated'
  is_force_update: boolean
  is_active: boolean
  download_url: string
  file_size: number
  checksum: string
  min_os_version: string
  required_features: string
  install_count: number
  update_count: number
  crash_count: number
  rating_average: number
  rating_count: number
  rollout_percent: number
  rollout_regions: string
  released_at?: string
  created_at: string
  updated_at: string
  created_by: string
  
  // Relations
  company?: any
  creator?: any
}

export interface MobileOfflineData {
  id: string
  device_id: string
  user_id: string
  company_id: string
  data_type: string
  resource_id: string
  operation: 'create' | 'update' | 'delete'
  data_payload: string
  data_checksum: string
  status: 'pending' | 'syncing' | 'synced' | 'failed' | 'conflict'
  priority: number
  conflict_data: string
  resolved_by?: string
  resolution_method: string
  last_sync_attempt?: string
  synced_at?: string
  error_message: string
  retry_count: number
  max_retries: number
  created_at: string
  updated_at: string
  expires_at?: string
  
  // Relations
  device?: MobileDevice
  user?: any
  company?: any
  resolved_by_user?: any
}

export interface MobileConfiguration {
  id: string
  company_id?: string
  platform: 'ios' | 'android' | 'web' | 'all'
  config_key: string
  config_value: string
  config_type: 'string' | 'number' | 'boolean' | 'object' | 'array'
  category: string
  description: string
  default_value: string
  validation_rules: string
  min_version: string
  max_version: string
  is_enabled: boolean
  rollout_percent: number
  target_devices: string
  cache_ttl: number
  is_secure: boolean
  created_at: string
  updated_at: string
  created_by: string
  
  // Relations
  company?: any
  creator?: any
}

export interface MobileStatistics {
  devices_by_platform: { platform: string; count: number }[]
  notifications_by_status: { status: string; count: number }[]
  total_devices: number
  active_devices: number
  total_notifications: number
}

export interface DeviceUsageStats {
  session_stats: {
    total_sessions: number
    avg_duration: number
    total_screen_views: number
    total_api_requests: number
  }
  event_counts: { event_type: string; count: number }[]
}

export interface NotificationStats {
  notifications_by_type: { type: string; count: number }[]
  delivery_stats: {
    total_sent: number
    total_delivered: number
    total_clicked: number
    total_failed: number
    delivery_rate: number
    click_rate: number
  }
}

export interface AnalyticsSummary {
  event_counts: { event_type: string; count: number }[]
  active_users: number
  session_stats: {
    total_sessions: number
    avg_duration: number
    total_screen_views: number
  }
}

export interface RegisterDeviceRequest {
  device_token: string
  platform: string
  device_type: string
  device_model?: string
  os_version?: string
  app_version?: string
  device_name?: string
  device_id?: string
  bundle_id?: string
  time_zone?: string
  language?: string
  country?: string
}

export interface SendPushNotificationRequest {
  device_id?: string
  user_id?: string
  title: string
  body: string
  type: string
  category?: string
  data?: Record<string, any>
  resource_id?: string
  resource_type?: string
  priority?: string
  badge?: number
  sound?: string
  icon?: string
  image?: string
  actions?: Array<{ action: string; title: string; icon?: string }>
  scheduled_at?: string
}

export interface SendBulkNotificationRequest {
  notifications: SendPushNotificationRequest[]
}

export interface SendNotificationToUsersRequest {
  user_ids: string[]
  title: string
  body: string
  type: string
  data?: Record<string, any>
}

export interface CreateMobileSessionRequest {
  device_token: string
  user_id: string
}

export interface UpdateMobileSessionRequest {
  app_state?: string
  network_type?: string
  screen_views?: number
  api_requests?: number
  data_transferred?: number
}

export interface TrackEventRequest {
  device_id: string
  user_id: string
  event_type: string
  event_name: string
  event_category?: string
  event_data?: Record<string, any>
  screen_name?: string
  screen_class?: string
  session_id?: string
  duration?: number
  interaction_type?: string
  element_id?: string
  element_type?: string
  load_time?: number
  response_time?: number
  error_message?: string
  app_version?: string
  os_version?: string
  network_type?: string
  battery_level?: number
  memory_usage?: number
  latitude?: number
  longitude?: number
  location_accuracy?: number
}

export interface TrackScreenViewRequest {
  device_id: string
  user_id: string
  screen_name: string
  screen_class: string
  duration: number
}

export interface TrackUserInteractionRequest {
  device_id: string
  user_id: string
  interaction_type: string
  element_id: string
  data?: Record<string, any>
}

export interface CreateAppVersionRequest {
  version: string
  build_number: string
  platform: string
  release_type: string
  release_notes: string
  release_notes_en?: string
  download_url?: string
  file_size?: number
  checksum?: string
  min_os_version?: string
  required_features?: string[]
  rollout_percent?: number
  rollout_regions?: string[]
}

export interface CreateOfflineDataRequest {
  device_id: string
  user_id: string
  data_type: string
  resource_id: string
  operation: string
  data_payload: Record<string, any>
  priority?: number
}

export interface CreateMobileConfigRequest {
  platform: string
  config_key: string
  config_value: any
  config_type: string
  category?: string
  description?: string
  default_value?: any
  validation_rules?: Record<string, any>
  min_version?: string
  max_version?: string
  is_enabled?: boolean
  rollout_percent?: number
  target_devices?: Record<string, any>
  cache_ttl?: number
  is_secure?: boolean
}

export interface ListParams {
  page?: number
  page_size?: number
  search?: string
  platform?: string
  status?: string
  type?: string
  is_active?: boolean
  user_id?: string
  device_id?: string
  start_date?: string
  end_date?: string
  event_type?: string
  data_type?: string
  release_type?: string
  category?: string
  is_enabled?: boolean
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
}

class MobileService {
  // Device Management
  async registerDevice(data: RegisterDeviceRequest): Promise<MobileDevice> {
    const response = await api.post('/mobile/devices/register', data)
    return response.data
  }

  async updateDevice(id: string, data: Partial<MobileDevice>): Promise<MobileDevice> {
    const response = await api.put(`/mobile/devices/${id}`, data)
    return response.data
  }

  async getDevice(id: string): Promise<MobileDevice> {
    const response = await api.get(`/mobile/devices/${id}`)
    return response.data
  }

  async getDeviceByToken(token: string): Promise<MobileDevice> {
    const response = await api.get(`/mobile/devices/token/${token}`)
    return response.data
  }

  async listUserDevices(userId: string, params: ListParams = {}): Promise<PaginatedResponse<MobileDevice>> {
    const response = await api.get(`/mobile/users/${userId}/devices`, { params })
    return response.data
  }

  async listCompanyDevices(params: ListParams = {}): Promise<PaginatedResponse<MobileDevice>> {
    const response = await api.get('/mobile/devices', { params })
    return response.data
  }

  async deactivateDevice(id: string): Promise<void> {
    await api.put(`/mobile/devices/${id}/deactivate`)
  }

  async updateDeviceLastSeen(token: string): Promise<void> {
    await api.put(`/mobile/devices/token/${token}/last-seen`)
  }

  // Push Notifications
  async sendPushNotification(data: SendPushNotificationRequest): Promise<PushNotification> {
    const response = await api.post('/mobile/notifications/send', data)
    return response.data
  }

  async sendBulkPushNotifications(data: SendBulkNotificationRequest): Promise<{ count: number }> {
    const response = await api.post('/mobile/notifications/send-bulk', data)
    return response.data
  }

  async listPushNotifications(params: ListParams = {}): Promise<PaginatedResponse<PushNotification>> {
    const response = await api.get('/mobile/notifications', { params })
    return response.data
  }

  async listDeviceNotifications(deviceId: string, params: ListParams = {}): Promise<PaginatedResponse<PushNotification>> {
    const response = await api.get(`/mobile/devices/${deviceId}/notifications`, { params })
    return response.data
  }

  async markNotificationDelivered(id: string): Promise<void> {
    await api.put(`/mobile/notifications/${id}/delivered`)
  }

  async markNotificationClicked(id: string): Promise<void> {
    await api.put(`/mobile/notifications/${id}/clicked`)
  }

  async sendNotificationToUsers(data: SendNotificationToUsersRequest): Promise<{ user_count: number }> {
    const response = await api.post('/mobile/notifications/send-to-users', data)
    return response.data
  }

  async processPendingNotifications(limit: number = 50): Promise<void> {
    await api.post('/mobile/notifications/process-pending', { limit })
  }

  // Session Management
  async createMobileSession(data: CreateMobileSessionRequest): Promise<MobileSession> {
    const response = await api.post('/mobile/sessions', data)
    return response.data
  }

  async updateMobileSession(token: string, data: UpdateMobileSessionRequest): Promise<void> {
    await api.put(`/mobile/sessions/${token}`, data)
  }

  async getMobileSession(token: string): Promise<MobileSession> {
    const response = await api.get(`/mobile/sessions/${token}`)
    return response.data
  }

  async endMobileSession(token: string): Promise<void> {
    await api.delete(`/mobile/sessions/${token}`)
  }

  async validateSession(token: string): Promise<{ valid: boolean; session: MobileSession }> {
    const response = await api.get(`/mobile/sessions/${token}/validate`)
    return response.data
  }

  // Analytics
  async trackEvent(data: TrackEventRequest): Promise<void> {
    await api.post('/mobile/analytics/track', data)
  }

  async trackScreenView(data: TrackScreenViewRequest): Promise<void> {
    await api.post('/mobile/analytics/screen-view', data)
  }

  async trackUserInteraction(data: TrackUserInteractionRequest): Promise<void> {
    await api.post('/mobile/analytics/interaction', data)
  }

  async getAnalyticsSummary(startDate?: string, endDate?: string): Promise<AnalyticsSummary> {
    const response = await api.get('/mobile/analytics/summary', {
      params: { start_date: startDate, end_date: endDate }
    })
    return response.data
  }

  // App Version Management
  async createAppVersion(data: CreateAppVersionRequest): Promise<MobileAppVersion> {
    const response = await api.post('/mobile/app-versions', data)
    return response.data
  }

  async updateAppVersion(id: string, data: Partial<CreateAppVersionRequest>): Promise<MobileAppVersion> {
    const response = await api.put(`/mobile/app-versions/${id}`, data)
    return response.data
  }

  async getLatestAppVersion(platform: string): Promise<MobileAppVersion> {
    const response = await api.get('/mobile/app-versions/latest', {
      params: { platform }
    })
    return response.data
  }

  async checkForUpdates(deviceToken: string, currentVersion: string): Promise<{ has_update: boolean; latest_version: MobileAppVersion }> {
    const response = await api.get('/mobile/app-versions/check-updates', {
      params: { device_token: deviceToken, current_version: currentVersion }
    })
    return response.data
  }

  // Offline Data Sync
  async createOfflineData(data: CreateOfflineDataRequest): Promise<MobileOfflineData> {
    const response = await api.post('/mobile/offline-data', data)
    return response.data
  }

  async syncOfflineData(deviceId: string, limit: number = 10): Promise<void> {
    await api.post(`/mobile/devices/${deviceId}/sync`, { limit })
  }

  async listPendingOfflineData(deviceId: string, limit: number = 50): Promise<PaginatedResponse<MobileOfflineData>> {
    const response = await api.get(`/mobile/devices/${deviceId}/offline-data/pending`, {
      params: { limit }
    })
    return response.data
  }

  // Configuration Management
  async getMobileConfig(key: string, platform: string): Promise<{ key: string; value: any }> {
    const response = await api.get(`/mobile/config/${key}`, {
      params: { platform }
    })
    return response.data
  }

  async getMobileConfigs(platform: string): Promise<Record<string, any>> {
    const response = await api.get('/mobile/config', {
      params: { platform }
    })
    return response.data
  }

  async setMobileConfig(data: CreateMobileConfigRequest): Promise<MobileConfiguration> {
    const response = await api.post('/mobile/config', data)
    return response.data
  }

  // Business Operations
  async getMobileStatistics(): Promise<MobileStatistics> {
    const response = await api.get('/mobile/statistics')
    return response.data
  }

  async getDeviceUsageStats(deviceId: string, days: number = 30): Promise<DeviceUsageStats> {
    const response = await api.get(`/mobile/devices/${deviceId}/usage-stats`, {
      params: { days }
    })
    return response.data
  }

  async getNotificationStats(days: number = 30): Promise<NotificationStats> {
    const response = await api.get('/mobile/notifications/stats', {
      params: { days }
    })
    return response.data
  }

  async generateMobileDashboard(): Promise<Record<string, any>> {
    const response = await api.get('/mobile/dashboard')
    return response.data
  }

  // PWA Utilities
  async registerServiceWorker(): Promise<ServiceWorkerRegistration | null> {
    if ('serviceWorker' in navigator) {
      try {
        const registration = await navigator.serviceWorker.register('/sw.js')
        console.log('Service Worker registered:', registration)
        return registration
      } catch (error) {
        console.error('Service Worker registration failed:', error)
        return null
      }
    }
    return null
  }

  async requestNotificationPermission(): Promise<NotificationPermission> {
    if ('Notification' in window) {
      return await Notification.requestPermission()
    }
    return 'denied'
  }

  async subscribeToPushNotifications(registration: ServiceWorkerRegistration): Promise<PushSubscription | null> {
    try {
      const subscription = await registration.pushManager.subscribe({
        userVisibleOnly: true,
        applicationServerKey: this.urlBase64ToUint8Array(process.env.NEXT_PUBLIC_VAPID_PUBLIC_KEY || '')
      })
      
      // Send subscription to server
      await api.post('/mobile/push-subscriptions', {
        subscription: subscription.toJSON()
      })
      
      return subscription
    } catch (error) {
      console.error('Push subscription failed:', error)
      return null
    }
  }

  private urlBase64ToUint8Array(base64String: string): Uint8Array {
    const padding = '='.repeat((4 - base64String.length % 4) % 4)
    const base64 = (base64String + padding)
      .replace(/-/g, '+')
      .replace(/_/g, '/')

    const rawData = window.atob(base64)
    const outputArray = new Uint8Array(rawData.length)

    for (let i = 0; i < rawData.length; ++i) {
      outputArray[i] = rawData.charCodeAt(i)
    }
    return outputArray
  }

  // Offline Data Management
  async storeOfflineData(key: string, data: any): Promise<void> {
    if ('indexedDB' in window) {
      // 使用 IndexedDB 存儲離線資料
      const request = indexedDB.open('FastenMindOffline', 1)
      
      request.onupgradeneeded = (event) => {
        const db = (event.target as IDBOpenDBRequest).result
        if (!db.objectStoreNames.contains('offlineData')) {
          db.createObjectStore('offlineData', { keyPath: 'key' })
        }
      }
      
      request.onsuccess = (event) => {
        const db = (event.target as IDBOpenDBRequest).result
        const transaction = db.transaction(['offlineData'], 'readwrite')
        const store = transaction.objectStore('offlineData')
        store.put({ key, data, timestamp: Date.now() })
      }
    } else {
      // 後備方案：使用 localStorage
      localStorage.setItem(`offline_${key}`, JSON.stringify({
        data,
        timestamp: Date.now()
      }))
    }
  }

  async getOfflineData(key: string): Promise<any> {
    if ('indexedDB' in window) {
      return new Promise((resolve, reject) => {
        const request = indexedDB.open('FastenMindOffline', 1)
        
        request.onsuccess = (event) => {
          const db = (event.target as IDBOpenDBRequest).result
          const transaction = db.transaction(['offlineData'], 'readonly')
          const store = transaction.objectStore('offlineData')
          const getRequest = store.get(key)
          
          getRequest.onsuccess = () => {
            resolve(getRequest.result?.data || null)
          }
          
          getRequest.onerror = () => {
            reject(getRequest.error)
          }
        }
        
        request.onerror = () => {
          reject(request.error)
        }
      })
    } else {
      // 後備方案：使用 localStorage
      const stored = localStorage.getItem(`offline_${key}`)
      return stored ? JSON.parse(stored).data : null
    }
  }

  // Device Information
  getDeviceInfo(): Partial<RegisterDeviceRequest> {
    const userAgent = navigator.userAgent
    const platform = this.getPlatform()
    
    return {
      platform,
      device_type: this.getDeviceType(),
      os_version: this.getOSVersion(),
      app_version: process.env.NEXT_PUBLIC_APP_VERSION || '1.0.0',
      device_name: this.getDeviceName(),
      time_zone: Intl.DateTimeFormat().resolvedOptions().timeZone,
      language: navigator.language,
      country: this.getCountry()
    }
  }

  private getPlatform(): string {
    const userAgent = navigator.userAgent.toLowerCase()
    if (/iphone|ipad|ipod/.test(userAgent)) return 'ios'
    if (/android/.test(userAgent)) return 'android'
    return 'web'
  }

  private getDeviceType(): string {
    const userAgent = navigator.userAgent.toLowerCase()
    if (/tablet|ipad/.test(userAgent)) return 'tablet'
    if (/mobile|iphone|android/.test(userAgent)) return 'phone'
    return 'desktop'
  }

  private getOSVersion(): string {
    const userAgent = navigator.userAgent
    const match = userAgent.match(/(iPhone OS|Android|Windows NT|Mac OS X) ([\d._]+)/)
    return match ? match[2].replace(/_/g, '.') : 'unknown'
  }

  private getDeviceName(): string {
    const userAgent = navigator.userAgent
    if (/iPhone/.test(userAgent)) return 'iPhone'
    if (/iPad/.test(userAgent)) return 'iPad'
    if (/Android/.test(userAgent)) return 'Android Device'
    return 'Web Browser'
  }

  private getCountry(): string {
    // 這裡可以使用 IP 地理位置服務或其他方法獲取國家代碼
    return 'TW' // 預設為台灣
  }
}

export default new MobileService()