import api from './api'

// Types
export interface AIAssistant {
  id: string
  company_id: string
  user_id: string
  name: string
  type: 'chat' | 'recommendation' | 'analysis' | 'automation'
  model: string
  status: 'active' | 'inactive' | 'training' | 'error'
  configuration: string
  system_prompt: string
  temperature: number
  max_tokens: number
  top_p: number
  frequency_penalty: number
  presence_penalty: number
  is_active: boolean
  usage_count: number
  tokens_used: number
  cost_accumulated: number
  last_used?: string
  created_at: string
  updated_at: string
  created_by: string
  
  // Relations
  company?: any
  user?: any
  creator?: any
  sessions?: AIConversationSession[]
}

export interface AIConversationSession {
  id: string
  assistant_id: string
  user_id: string
  company_id: string
  title: string
  context: string
  status: 'active' | 'completed' | 'archived'
  message_count: number
  tokens_used: number
  cost: number
  start_time: string
  end_time?: string
  created_at: string
  updated_at: string
  
  // Relations
  assistant?: AIAssistant
  user?: any
  company?: any
  messages?: AIMessage[]
}

export interface AIMessage {
  id: string
  session_id: string
  role: 'user' | 'assistant' | 'system'
  content: string
  token_count: number
  cost: number
  model_used: string
  response_time: number
  metadata: string
  created_at: string
  
  // Relations
  session?: AIConversationSession
}

export interface SmartRecommendation {
  id: string
  company_id: string
  user_id: string
  type: 'product' | 'customer' | 'supplier' | 'pricing' | 'process'
  category: string
  title: string
  description: string
  data: string
  score: number
  priority: 'low' | 'medium' | 'high' | 'urgent'
  status: 'pending' | 'viewed' | 'accepted' | 'rejected' | 'implemented'
  source: 'ai' | 'algorithm' | 'manual' | 'user_behavior'
  source_data: string
  resource_id?: string
  resource_type: string
  expires_at?: string
  viewed_at?: string
  actioned_at?: string
  created_at: string
  updated_at: string
  
  // Relations
  company?: any
  user?: any
}

export interface AdvancedSearch {
  id: string
  company_id: string
  user_id: string
  name: string
  description: string
  search_type: 'table' | 'global' | 'cross_reference'
  table_name: string
  filters: string
  sort_config: string
  columns: string
  is_public: boolean
  usage_count: number
  last_used?: string
  created_at: string
  updated_at: string
  
  // Relations
  company?: any
  user?: any
}

export interface BatchOperation {
  id: string
  company_id: string
  user_id: string
  operation_type: 'update' | 'delete' | 'export' | 'import' | 'send_email' | 'change_status'
  target_table: string
  target_ids: string
  parameters: string
  status: 'pending' | 'running' | 'completed' | 'failed' | 'cancelled'
  progress: number
  total_items: number
  processed_items: number
  success_count: number
  error_count: number
  error_log: string
  result: string
  started_at?: string
  completed_at?: string
  created_at: string
  updated_at: string
  
  // Relations
  company?: any
  user?: any
}

export interface CustomField {
  id: string
  company_id: string
  table_name: string
  field_name: string
  field_label: string
  field_type: 'text' | 'number' | 'date' | 'boolean' | 'select' | 'multi_select' | 'file'
  default_value: string
  options: string
  validation: string
  is_required: boolean
  is_searchable: boolean
  is_active: boolean
  display_order: number
  created_at: string
  updated_at: string
  created_by: string
  
  // Relations
  company?: any
  creator?: any
}

export interface CustomFieldValue {
  id: string
  company_id: string
  field_id: string
  resource_id: string
  resource_type: string
  value: string
  created_at: string
  updated_at: string
  
  // Relations
  company?: any
  field?: CustomField
}

export interface SecurityEvent {
  id: string
  company_id: string
  user_id?: string
  event_type: 'login' | 'logout' | 'failed_login' | 'data_access' | 'permission_change' | 'suspicious_activity'
  severity: 'low' | 'medium' | 'high' | 'critical'
  description: string
  ip_address: string
  user_agent: string
  location: string
  device_info: string
  resource_type: string
  resource_id?: string
  action_details: string
  risk_score: number
  status: 'new' | 'investigating' | 'resolved' | 'false_positive'
  resolved_by?: string
  resolved_at?: string
  created_at: string
  
  // Relations
  company?: any
  user?: any
  resolved_by_user?: any
}

export interface PerformanceMetric {
  id: string
  company_id: string
  metric_type: 'api_response_time' | 'database_query_time' | 'page_load_time' | 'memory_usage' | 'cpu_usage'
  metric_name: string
  value: number
  unit: string
  context: string
  endpoint: string
  method: string
  status_code: number
  user_id?: string
  session_id: string
  trace_id: string
  created_at: string
  
  // Relations
  company?: any
  user?: any
}

export interface BackupRecord {
  id: string
  company_id: string
  backup_type: 'full' | 'incremental' | 'differential'
  status: 'running' | 'completed' | 'failed' | 'cancelled'
  file_size: number
  file_path: string
  checksum: string
  tables: string
  compression_type: string
  encryption_type: string
  start_time: string
  end_time?: string
  duration: number
  error_message: string
  created_by: string
  created_at: string
  
  // Relations
  company?: any
  creator?: any
}

export interface SystemLanguage {
  id: string
  company_id?: string
  language_code: string
  language_name: string
  native_name: string
  is_active: boolean
  is_default: boolean
  rtl: boolean
  date_format: string
  time_format: string
  number_format: string
  currency_format: string
  translation_progress: number
  created_at: string
  updated_at: string
  
  // Relations
  company?: any
  translations?: Translation[]
}

export interface Translation {
  id: string
  language_id: string
  translation_key: string
  translation: string
  context: string
  is_approved: boolean
  created_at: string
  updated_at: string
  created_by: string
  approved_by?: string
  
  // Relations
  language?: SystemLanguage
  creator?: any
  approver?: any
}

// Request types
export interface CreateAIAssistantRequest {
  name: string
  type: string
  model: string
  system_prompt?: string
  temperature?: number
  max_tokens?: number
  top_p?: number
  frequency_penalty?: number
  presence_penalty?: number
  configuration?: Record<string, any>
}

export interface UpdateAIAssistantRequest {
  name?: string
  system_prompt?: string
  temperature?: number
  max_tokens?: number
  is_active?: boolean
  configuration?: Record<string, any>
}

export interface StartConversationRequest {
  assistant_id: string
  title?: string
  context?: Record<string, any>
}

export interface SendMessageRequest {
  content: string
}

export interface CreateRecommendationRequest {
  type: string
  category?: string
  title: string
  description?: string
  data?: Record<string, any>
  score: number
  priority: string
  source: string
  source_data?: Record<string, any>
  resource_id?: string
  resource_type?: string
  expires_at?: string
}

export interface CreateAdvancedSearchRequest {
  name: string
  description?: string
  search_type: string
  table_name?: string
  filters?: Record<string, any>
  sort_config?: Record<string, any>
  columns?: string[]
  is_public?: boolean
}

export interface CreateBatchOperationRequest {
  operation_type: string
  target_table: string
  target_ids: string[]
  parameters?: Record<string, any>
}

export interface CreateCustomFieldRequest {
  table_name: string
  field_name: string
  field_label: string
  field_type: string
  default_value?: string
  options?: string[]
  validation?: Record<string, any>
  is_required?: boolean
  is_searchable?: boolean
  display_order?: number
}

export interface SetCustomFieldValueRequest {
  field_id: string
  resource_id: string
  resource_type: string
  value: string
}

export interface CreateSecurityEventRequest {
  company_id: string
  user_id?: string
  event_type: string
  severity: string
  description?: string
  ip_address?: string
  user_agent?: string
  location?: string
  device_info?: Record<string, any>
  resource_type?: string
  resource_id?: string
  action_details?: Record<string, any>
  risk_score: number
}

export interface RecordPerformanceMetricRequest {
  company_id: string
  metric_type: string
  metric_name: string
  value: number
  unit?: string
  context?: Record<string, any>
  endpoint?: string
  method?: string
  status_code?: number
  user_id?: string
  session_id?: string
  trace_id?: string
}

export interface CreateBackupRecordRequest {
  backup_type: string
  tables?: string[]
  compression_type?: string
  encryption_type?: string
}

export interface ListParams {
  type?: string
  status?: string
  is_active?: boolean
  is_public?: boolean
  search_type?: string
  table_name?: string
  resource_type?: string
  event_type?: string
  severity?: string
  metric_type?: string
  backup_type?: string
  language_code?: string
  only_approved?: boolean
  limit?: number
  start_time?: string
  end_time?: string
}

export interface PaginatedResponse<T> {
  data: T[]
  total?: number
}

export interface AIMessageResponse {
  content: string
  token_count: number
  cost: number
  response_time: number
}

export interface AdvancedSearchResult {
  search_id: string
  executed_at: string
  total_count: number
  results: Record<string, any>[]
  execution_time: number
}

class AdvancedService {
  // AI Assistant Methods
  async createAIAssistant(data: CreateAIAssistantRequest): Promise<AIAssistant> {
    const response = await api.post('/advanced/ai-assistants', data)
    return response.data
  }

  async getAIAssistant(id: string): Promise<AIAssistant> {
    const response = await api.get(`/advanced/ai-assistants/${id}`)
    return response.data
  }

  async listAIAssistants(params: ListParams = {}): Promise<PaginatedResponse<AIAssistant>> {
    const response = await api.get('/advanced/ai-assistants', { params })
    return response.data
  }

  async updateAIAssistant(id: string, data: UpdateAIAssistantRequest): Promise<AIAssistant> {
    const response = await api.put(`/advanced/ai-assistants/${id}`, data)
    return response.data
  }

  async deleteAIAssistant(id: string): Promise<void> {
    await api.delete(`/advanced/ai-assistants/${id}`)
  }

  // AI Conversation Methods
  async startConversation(data: StartConversationRequest): Promise<AIConversationSession> {
    const response = await api.post('/advanced/conversations', data)
    return response.data
  }

  async sendMessage(sessionId: string, data: SendMessageRequest): Promise<AIMessageResponse> {
    const response = await api.post(`/advanced/conversations/${sessionId}/messages`, data)
    return response.data
  }

  async getConversationHistory(sessionId: string, limit?: number): Promise<PaginatedResponse<AIMessage>> {
    const response = await api.get(`/advanced/conversations/${sessionId}/messages`, {
      params: { limit }
    })
    return response.data
  }

  async endConversation(sessionId: string): Promise<void> {
    await api.delete(`/advanced/conversations/${sessionId}`)
  }

  // Smart Recommendation Methods
  async createRecommendation(data: CreateRecommendationRequest): Promise<SmartRecommendation> {
    const response = await api.post('/advanced/recommendations', data)
    return response.data
  }

  async listRecommendations(params: ListParams = {}): Promise<PaginatedResponse<SmartRecommendation>> {
    const response = await api.get('/advanced/recommendations', { params })
    return response.data
  }

  async updateRecommendationStatus(id: string, status: string): Promise<SmartRecommendation> {
    const response = await api.put(`/advanced/recommendations/${id}/status`, { status })
    return response.data
  }

  // Advanced Search Methods
  async createAdvancedSearch(data: CreateAdvancedSearchRequest): Promise<AdvancedSearch> {
    const response = await api.post('/advanced/searches', data)
    return response.data
  }

  async listAdvancedSearches(params: ListParams = {}): Promise<PaginatedResponse<AdvancedSearch>> {
    const response = await api.get('/advanced/searches', { params })
    return response.data
  }

  async executeAdvancedSearch(id: string): Promise<AdvancedSearchResult> {
    const response = await api.post(`/advanced/searches/${id}/execute`)
    return response.data
  }

  // Batch Operation Methods
  async createBatchOperation(data: CreateBatchOperationRequest): Promise<BatchOperation> {
    const response = await api.post('/advanced/batch-operations', data)
    return response.data
  }

  async listBatchOperations(params: ListParams = {}): Promise<PaginatedResponse<BatchOperation>> {
    const response = await api.get('/advanced/batch-operations', { params })
    return response.data
  }

  async getBatchOperation(id: string): Promise<BatchOperation> {
    const response = await api.get(`/advanced/batch-operations/${id}`)
    return response.data
  }

  // Custom Field Methods
  async createCustomField(data: CreateCustomFieldRequest): Promise<CustomField> {
    const response = await api.post('/advanced/custom-fields', data)
    return response.data
  }

  async listCustomFields(params: ListParams): Promise<PaginatedResponse<CustomField>> {
    const response = await api.get('/advanced/custom-fields', { params })
    return response.data
  }

  async setCustomFieldValue(data: SetCustomFieldValueRequest): Promise<void> {
    await api.post('/advanced/custom-field-values', data)
  }

  async getCustomFieldValues(resourceId: string, resourceType: string): Promise<PaginatedResponse<CustomFieldValue>> {
    const response = await api.get(`/advanced/custom-field-values/${resourceId}`, {
      params: { resource_type: resourceType }
    })
    return response.data
  }

  // Security Event Methods
  async createSecurityEvent(data: CreateSecurityEventRequest): Promise<SecurityEvent> {
    const response = await api.post('/advanced/security-events', data)
    return response.data
  }

  async listSecurityEvents(params: ListParams = {}): Promise<PaginatedResponse<SecurityEvent>> {
    const response = await api.get('/advanced/security-events', { params })
    return response.data
  }

  // Performance Metric Methods
  async recordPerformanceMetric(data: RecordPerformanceMetricRequest): Promise<void> {
    await api.post('/advanced/performance-metrics', data)
  }

  async getPerformanceStats(params: ListParams): Promise<Record<string, any>> {
    const response = await api.get('/advanced/performance-stats', { params })
    return response.data
  }

  // Backup Methods
  async createBackup(data: CreateBackupRecordRequest): Promise<BackupRecord> {
    const response = await api.post('/advanced/backups', data)
    return response.data
  }

  async listBackups(params: ListParams = {}): Promise<PaginatedResponse<BackupRecord>> {
    const response = await api.get('/advanced/backups', { params })
    return response.data
  }

  // Multi-language Support Methods
  async listSystemLanguages(params: ListParams = {}): Promise<PaginatedResponse<SystemLanguage>> {
    const response = await api.get('/advanced/languages', { params })
    return response.data
  }

  async getTranslations(languageCode: string, onlyApproved: boolean = true): Promise<Record<string, string>> {
    const response = await api.get(`/advanced/translations/${languageCode}`, {
      params: { only_approved: onlyApproved }
    })
    return response.data
  }

  // Analytics Methods
  async getAIUsageStats(startTime?: string, endTime?: string): Promise<Record<string, any>> {
    const response = await api.get('/advanced/analytics/ai-usage', {
      params: { start_time: startTime, end_time: endTime }
    })
    return response.data
  }

  async getRecommendationStats(startTime?: string, endTime?: string): Promise<Record<string, any>> {
    const response = await api.get('/advanced/analytics/recommendations', {
      params: { start_time: startTime, end_time: endTime }
    })
    return response.data
  }

  async getSecurityEventStats(startTime?: string, endTime?: string): Promise<Record<string, any>> {
    const response = await api.get('/advanced/analytics/security-events', {
      params: { start_time: startTime, end_time: endTime }
    })
    return response.data
  }

  // Utility Methods
  async testAIConnection(assistantId: string): Promise<{ success: boolean; message: string }> {
    const response = await api.post(`/advanced/ai-assistants/${assistantId}/test`)
    return response.data
  }

  async exportAdvancedSearchResults(searchId: string, format: 'csv' | 'xlsx' | 'json' = 'csv'): Promise<Blob> {
    const response = await api.get(`/advanced/searches/${searchId}/export`, {
      params: { format },
      responseType: 'blob'
    })
    return response.data
  }

  async generateAdvancedReport(type: string, params: Record<string, any> = {}): Promise<Record<string, any>> {
    const response = await api.post('/advanced/reports/generate', { type, params })
    return response.data
  }

  async validateCustomFieldConfiguration(config: Record<string, any>): Promise<{ valid: boolean; errors: string[] }> {
    const response = await api.post('/advanced/custom-fields/validate', config)
    return response.data
  }

  async previewBatchOperation(data: CreateBatchOperationRequest): Promise<{ affected_count: number; preview: any[] }> {
    const response = await api.post('/advanced/batch-operations/preview', data)
    return response.data
  }

  async optimizeSystemPerformance(): Promise<{ recommendations: any[]; estimated_improvement: number }> {
    const response = await api.post('/advanced/system/optimize')
    return response.data
  }
}

export default new AdvancedService()