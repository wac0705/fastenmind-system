import api from './api'

// Types
export interface Integration {
  id: string
  company_id: string
  name: string
  type: 'api' | 'webhook' | 'ftp' | 'sftp' | 'email' | 'database'
  provider: string
  status: 'active' | 'inactive' | 'error' | 'testing'
  configuration: string
  credentials: string
  api_version: string
  base_url: string
  auth_type: 'none' | 'api_key' | 'oauth2' | 'basic_auth' | 'token'
  headers: string
  rate_limit_rpm: number
  timeout_seconds: number
  retry_attempts: number
  is_active: boolean
  last_sync_at?: string
  last_error_at?: string
  last_error: string
  sync_count: number
  error_count: number
  success_rate: number
  avg_response_time: number
  created_at: string
  updated_at: string
  created_by: string
  
  // Relations
  company?: any
  creator?: any
  mappings?: IntegrationMapping[]
  webhooks?: Webhook[]
  sync_jobs?: DataSyncJob[]
  logs?: IntegrationLog[]
}

export interface IntegrationMapping {
  id: string
  company_id: string
  integration_id: string
  name: string
  direction: 'inbound' | 'outbound' | 'bidirectional'
  source_table: string
  target_table: string
  source_endpoint: string
  target_endpoint: string
  field_mappings: string
  transformations: string
  filters: string
  sync_frequency: 'realtime' | 'hourly' | 'daily' | 'weekly' | 'manual'
  is_active: boolean
  last_sync_at?: string
  next_sync_at?: string
  created_at: string
  updated_at: string
  created_by: string
  
  // Relations
  company?: any
  integration?: Integration
  creator?: any
}

export interface Webhook {
  id: string
  company_id: string
  integration_id?: string
  name: string
  url: string
  method: string
  headers: string
  auth_type: 'none' | 'api_key' | 'basic_auth' | 'bearer_token'
  auth_config: string
  events: string
  payload_format: 'json' | 'xml' | 'form'
  payload_template: string
  is_active: boolean
  retry_attempts: number
  retry_interval: number
  timeout_seconds: number
  last_triggered_at?: string
  trigger_count: number
  success_count: number
  failure_count: number
  created_at: string
  updated_at: string
  created_by: string
  
  // Relations
  company?: any
  integration?: Integration
  creator?: any
  deliveries?: WebhookDelivery[]
}

export interface WebhookDelivery {
  id: string
  webhook_id: string
  company_id: string
  event_type: string
  event_data: string
  request_url: string
  request_method: string
  request_headers: string
  request_body: string
  response_code: number
  response_headers: string
  response_body: string
  status: 'pending' | 'success' | 'failed' | 'retrying'
  attempt_count: number
  next_retry_at?: string
  completed_at?: string
  error_message: string
  response_time: number
  created_at: string
  updated_at: string
  
  // Relations
  webhook?: Webhook
  company?: any
}

export interface DataSyncJob {
  id: string
  company_id: string
  integration_id: string
  mapping_id?: string
  name: string
  type: 'full_sync' | 'incremental_sync' | 'delta_sync'
  direction: 'import' | 'export' | 'bidirectional'
  status: 'pending' | 'running' | 'completed' | 'failed' | 'cancelled'
  priority: 'low' | 'normal' | 'high' | 'urgent'
  scheduled_at?: string
  started_at?: string
  completed_at?: string
  duration: number
  total_records: number
  processed_records: number
  success_records: number
  error_records: number
  skipped_records: number
  progress: number
  configuration: string
  result: string
  error_log: string
  created_at: string
  updated_at: string
  created_by: string
  
  // Relations
  company?: any
  integration?: Integration
  mapping?: IntegrationMapping
  creator?: any
}

export interface IntegrationLog {
  id: string
  company_id: string
  integration_id: string
  sync_job_id?: string
  level: 'debug' | 'info' | 'warning' | 'error' | 'critical'
  category: 'api_call' | 'data_sync' | 'webhook' | 'auth' | 'config'
  message: string
  details: string
  request_data: string
  response_data: string
  error_code: string
  error_message: string
  duration: number
  user_id?: string
  ip_address: string
  user_agent: string
  trace_id: string
  created_at: string
  
  // Relations
  company?: any
  integration?: Integration
  sync_job?: DataSyncJob
  user?: any
}

export interface ApiKey {
  id: string
  company_id: string
  user_id: string
  name: string
  description: string
  key_hash: string
  key_prefix: string
  permissions: string
  scopes: string
  rate_limit: number
  is_active: boolean
  expires_at?: string
  last_used_at?: string
  usage_count: number
  created_at: string
  updated_at: string
  created_by: string
  
  // Relations
  company?: any
  user?: any
  creator?: any
}

export interface ExternalSystem {
  id: string
  company_id: string
  name: string
  system_type: 'erp' | 'crm' | 'accounting' | 'warehouse' | 'shipping'
  vendor: string
  version: string
  base_url: string
  database_config: string
  api_config: string
  ftp_config: string
  sftp_config: string
  status: 'active' | 'inactive' | 'testing' | 'error'
  is_active: boolean
  last_test_at?: string
  last_test_result: string
  created_at: string
  updated_at: string
  created_by: string
  
  // Relations
  company?: any
  creator?: any
  integrations?: Integration[]
}

export interface IntegrationTemplate {
  id: string
  company_id?: string
  name: string
  description: string
  category: string
  provider: string
  version: string
  configuration: string
  mappings: string
  requirements: string
  is_public: boolean
  is_active: boolean
  usage_count: number
  rating: number
  created_at: string
  updated_at: string
  created_by: string
  
  // Relations
  company?: any
  creator?: any
}

// Request types
export interface CreateIntegrationRequest {
  name: string
  type: string
  provider: string
  api_version?: string
  base_url?: string
  auth_type?: string
  configuration?: Record<string, any>
  credentials?: Record<string, any>
  headers?: Record<string, string>
  rate_limit_rpm?: number
  timeout_seconds?: number
  retry_attempts?: number
}

export interface UpdateIntegrationRequest {
  name?: string
  status?: string
  base_url?: string
  is_active?: boolean
  configuration?: Record<string, any>
  credentials?: Record<string, any>
  headers?: Record<string, string>
}

export interface CreateIntegrationMappingRequest {
  integration_id: string
  name: string
  direction: string
  source_table?: string
  target_table?: string
  source_endpoint?: string
  target_endpoint?: string
  field_mappings?: Record<string, any>
  transformations?: Record<string, any>
  filters?: Record<string, any>
  sync_frequency: string
}

export interface UpdateIntegrationMappingRequest {
  name?: string
  sync_frequency?: string
  is_active?: boolean
  field_mappings?: Record<string, any>
  transformations?: Record<string, any>
  filters?: Record<string, any>
}

export interface CreateWebhookRequest {
  integration_id?: string
  name: string
  url: string
  method?: string
  headers?: Record<string, string>
  auth_type?: string
  auth_config?: Record<string, any>
  events?: string[]
  payload_format?: string
  payload_template?: string
  retry_attempts?: number
  retry_interval?: number
  timeout_seconds?: number
}

export interface UpdateWebhookRequest {
  name?: string
  url?: string
  is_active?: boolean
  headers?: Record<string, string>
  events?: string[]
}

export interface TriggerWebhookRequest {
  event_type: string
  event_data?: Record<string, any>
}

export interface CreateDataSyncJobRequest {
  integration_id: string
  mapping_id?: string
  name: string
  type: string
  direction: string
  priority?: string
  scheduled_at?: string
  configuration?: Record<string, any>
}

export interface CreateApiKeyRequest {
  name: string
  description?: string
  permissions?: string[]
  scopes?: string[]
  rate_limit?: number
  expires_at?: string
}

export interface CreateApiKeyResponse {
  id: string
  key: string
  key_id: string
  name: string
  expires_at?: string
}

export interface CreateExternalSystemRequest {
  name: string
  system_type: string
  vendor?: string
  version?: string
  base_url?: string
  database_config?: Record<string, any>
  api_config?: Record<string, any>
  ftp_config?: Record<string, any>
  sftp_config?: Record<string, any>
}

export interface CreateFromTemplateRequest {
  name: string
  configuration?: Record<string, any>
}

export interface IntegrationTestResult {
  success: boolean
  response_time: number
  message: string
  details: Record<string, any>
}

export interface SystemTestResult {
  success: boolean
  response_time: number
  message: string
  details: Record<string, any>
}

export interface ListParams {
  type?: string
  status?: string
  is_active?: boolean
  integration_id?: string
  user_id?: string
  system_type?: string
  category?: string
  provider?: string
  is_public?: boolean
  direction?: string
  sync_frequency?: string
  level?: string
  priority?: string
  limit?: number
  start_time?: string
  end_time?: string
  days?: number
}

export interface PaginatedResponse<T> {
  data: T[]
  total?: number
}

class IntegrationService {
  // Integration Methods
  async createIntegration(data: CreateIntegrationRequest): Promise<Integration> {
    const response = await api.post('/integrations', data)
    return response.data
  }

  async getIntegration(id: string): Promise<Integration> {
    const response = await api.get(`/integrations/${id}`)
    return response.data
  }

  async listIntegrations(params: ListParams = {}): Promise<PaginatedResponse<Integration>> {
    const response = await api.get('/integrations', { params })
    return response.data
  }

  async updateIntegration(id: string, data: UpdateIntegrationRequest): Promise<Integration> {
    const response = await api.put(`/integrations/${id}`, data)
    return response.data
  }

  async deleteIntegration(id: string): Promise<void> {
    await api.delete(`/integrations/${id}`)
  }

  async testIntegration(id: string): Promise<IntegrationTestResult> {
    const response = await api.post(`/integrations/${id}/test`)
    return response.data
  }

  // Integration Mapping Methods
  async createIntegrationMapping(data: CreateIntegrationMappingRequest): Promise<IntegrationMapping> {
    const response = await api.post('/integrations/mappings', data)
    return response.data
  }

  async listIntegrationMappings(integrationId: string, params: ListParams = {}): Promise<PaginatedResponse<IntegrationMapping>> {
    const response = await api.get(`/integrations/${integrationId}/mappings`, { params })
    return response.data
  }

  async updateIntegrationMapping(id: string, data: UpdateIntegrationMappingRequest): Promise<IntegrationMapping> {
    const response = await api.put(`/integrations/mappings/${id}`, data)
    return response.data
  }

  async deleteIntegrationMapping(id: string): Promise<void> {
    await api.delete(`/integrations/mappings/${id}`)
  }

  // Webhook Methods
  async createWebhook(data: CreateWebhookRequest): Promise<Webhook> {
    const response = await api.post('/integrations/webhooks', data)
    return response.data
  }

  async listWebhooks(params: ListParams = {}): Promise<PaginatedResponse<Webhook>> {
    const response = await api.get('/integrations/webhooks', { params })
    return response.data
  }

  async updateWebhook(id: string, data: UpdateWebhookRequest): Promise<Webhook> {
    const response = await api.put(`/integrations/webhooks/${id}`, data)
    return response.data
  }

  async deleteWebhook(id: string): Promise<void> {
    await api.delete(`/integrations/webhooks/${id}`)
  }

  async triggerWebhook(id: string, data: TriggerWebhookRequest): Promise<{ message: string }> {
    const response = await api.post(`/integrations/webhooks/${id}/trigger`, data)
    return response.data
  }

  async getWebhookDeliveries(webhookId: string, params: ListParams = {}): Promise<PaginatedResponse<WebhookDelivery>> {
    const response = await api.get(`/integrations/webhooks/${webhookId}/deliveries`, { params })
    return response.data
  }

  // Data Sync Job Methods
  async createDataSyncJob(data: CreateDataSyncJobRequest): Promise<DataSyncJob> {
    const response = await api.post('/integrations/sync-jobs', data)
    return response.data
  }

  async listDataSyncJobs(integrationId: string, params: ListParams = {}): Promise<PaginatedResponse<DataSyncJob>> {
    const response = await api.get(`/integrations/${integrationId}/sync-jobs`, { params })
    return response.data
  }

  async getDataSyncJob(id: string): Promise<DataSyncJob> {
    const response = await api.get(`/integrations/sync-jobs/${id}`)
    return response.data
  }

  async startDataSyncJob(jobId: string): Promise<{ message: string }> {
    const response = await api.post(`/integrations/sync-jobs/${jobId}/start`)
    return response.data
  }

  // API Key Methods
  async createApiKey(data: CreateApiKeyRequest): Promise<CreateApiKeyResponse> {
    const response = await api.post('/integrations/api-keys', data)
    return response.data
  }

  async listApiKeys(params: ListParams = {}): Promise<PaginatedResponse<ApiKey>> {
    const response = await api.get('/integrations/api-keys', { params })
    return response.data
  }

  async revokeApiKey(id: string): Promise<{ message: string }> {
    const response = await api.delete(`/integrations/api-keys/${id}`)
    return response.data
  }

  // External System Methods
  async createExternalSystem(data: CreateExternalSystemRequest): Promise<ExternalSystem> {
    const response = await api.post('/integrations/external-systems', data)
    return response.data
  }

  async listExternalSystems(params: ListParams = {}): Promise<PaginatedResponse<ExternalSystem>> {
    const response = await api.get('/integrations/external-systems', { params })
    return response.data
  }

  async getExternalSystem(id: string): Promise<ExternalSystem> {
    const response = await api.get(`/integrations/external-systems/${id}`)
    return response.data
  }

  async updateExternalSystem(id: string, data: Partial<CreateExternalSystemRequest>): Promise<ExternalSystem> {
    const response = await api.put(`/integrations/external-systems/${id}`, data)
    return response.data
  }

  async deleteExternalSystem(id: string): Promise<void> {
    await api.delete(`/integrations/external-systems/${id}`)
  }

  async testExternalSystem(id: string): Promise<SystemTestResult> {
    const response = await api.post(`/integrations/external-systems/${id}/test`)
    return response.data
  }

  // Integration Template Methods
  async listIntegrationTemplates(params: ListParams = {}): Promise<PaginatedResponse<IntegrationTemplate>> {
    const response = await api.get('/integrations/templates', { params })
    return response.data
  }

  async getIntegrationTemplate(id: string): Promise<IntegrationTemplate> {
    const response = await api.get(`/integrations/templates/${id}`)
    return response.data
  }

  async createIntegrationFromTemplate(templateId: string, data: CreateFromTemplateRequest): Promise<Integration> {
    const response = await api.post(`/integrations/templates/${templateId}/create`, data)
    return response.data
  }

  // Analytics Methods
  async getIntegrationStats(params: ListParams = {}): Promise<Record<string, any>> {
    const response = await api.get('/integrations/stats', { params })
    return response.data
  }

  async getIntegrationsByType(): Promise<{ data: Array<{ type: string; count: number }> }> {
    const response = await api.get('/integrations/analytics/by-type')
    return response.data
  }

  async getSyncJobTrends(params: ListParams = {}): Promise<{ data: Array<{ date: string; success_count: number; failure_count: number }> }> {
    const response = await api.get('/integrations/analytics/sync-trends', { params })
    return response.data
  }

  // Integration Log Methods
  async getIntegrationLogs(integrationId: string, params: ListParams = {}): Promise<PaginatedResponse<IntegrationLog>> {
    const response = await api.get(`/integrations/${integrationId}/logs`, { params })
    return response.data
  }

  // Utility Methods
  async validateMapping(data: { field_mappings?: Record<string, any>; transformations?: Record<string, any> }): Promise<{
    valid: boolean
    errors: string[]
    warnings: string[]
  }> {
    const response = await api.post('/integrations/validate-mapping', data)
    return response.data
  }

  async previewDataTransformation(data: {
    sample_data: Array<Record<string, any>>
    field_mappings?: Record<string, any>
    transformations?: Record<string, any>
  }): Promise<{
    preview: Array<{
      original: Record<string, any>
      transformed: Record<string, any>
    }>
  }> {
    const response = await api.post('/integrations/preview-transformation', data)
    return response.data
  }

  async exportIntegrationConfig(id: string, format: 'json' | 'yaml' = 'json'): Promise<Blob> {
    const response = await api.get(`/integrations/${id}/export`, {
      params: { format },
      responseType: 'blob'
    })
    return response.data
  }

  async importIntegrationConfig(file: File): Promise<Integration> {
    const formData = new FormData()
    formData.append('file', file)
    
    const response = await api.post('/integrations/import', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
    return response.data
  }

  async generateApiClient(integrationId: string, language: 'javascript' | 'python' | 'php' | 'curl' = 'javascript'): Promise<{
    code: string
    filename: string
  }> {
    const response = await api.get(`/integrations/${integrationId}/generate-client`, {
      params: { language }
    })
    return response.data
  }

  async testWebhookUrl(url: string, method: string = 'POST', headers?: Record<string, string>): Promise<{
    success: boolean
    response_time: number
    status_code: number
    message: string
  }> {
    const response = await api.post('/integrations/test-webhook-url', {
      url,
      method,
      headers
    })
    return response.data
  }

  async scheduleDataSync(integrationId: string, data: {
    mapping_id?: string
    type: string
    direction: string
    scheduled_at: string
    configuration?: Record<string, any>
  }): Promise<DataSyncJob> {
    const response = await api.post(`/integrations/${integrationId}/schedule-sync`, data)
    return response.data
  }

  async cancelDataSyncJob(jobId: string): Promise<{ message: string }> {
    const response = await api.post(`/integrations/sync-jobs/${jobId}/cancel`)
    return response.data
  }

  async retryFailedSyncJob(jobId: string): Promise<{ message: string }> {
    const response = await api.post(`/integrations/sync-jobs/${jobId}/retry`)
    return response.data
  }

  async getIntegrationHealth(id: string): Promise<{
    status: 'healthy' | 'warning' | 'critical'
    last_check: string
    issues: Array<{ type: string; message: string; severity: string }>
    recommendations: string[]
  }> {
    const response = await api.get(`/integrations/${id}/health`)
    return response.data
  }

  async runIntegrationDiagnostics(id: string): Promise<{
    connectivity: boolean
    authentication: boolean
    permissions: boolean
    rate_limits: boolean
    data_quality: boolean
    issues: Array<{ component: string; message: string; resolution: string }>
  }> {
    const response = await api.post(`/integrations/${id}/diagnostics`)
    return response.data
  }
}

export default new IntegrationService()