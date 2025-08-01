import api from './api'

// Types
export interface Report {
  id: string
  company_id: string
  report_no: string
  name: string
  name_en: string
  category: 'sales' | 'finance' | 'production' | 'inventory' | 'supplier' | 'customer' | 'system'
  type: 'summary' | 'detail' | 'trend' | 'comparison' | 'dashboard'
  status: 'active' | 'inactive' | 'archived'
  data_source: string
  filters: string
  columns: string
  sorting: string
  grouping: string
  aggregation: string
  chart_config: string
  is_scheduled: boolean
  schedule: string
  recipients: string
  file_format: string
  template_id?: string
  layout: string
  styling: string
  is_public: boolean
  shared_with: string
  cache_enabled: boolean
  cache_ttl: number
  query_timeout: number
  description: string
  tags: string
  version: number
  view_count: number
  last_viewed?: string
  execute_count: number
  last_executed?: string
  avg_exec_time: number
  created_at: string
  updated_at: string
  created_by: string
  updated_by?: string
  
  // Relations
  company?: any
  template?: ReportTemplate
  creator?: any
  updated_by_user?: any
  executions?: ReportExecution[]
}

export interface ReportTemplate {
  id: string
  company_id?: string
  name: string
  name_en: string
  category: string
  type: string
  is_system_template: boolean
  data_source: string
  filters: string
  columns: string
  sorting: string
  grouping: string
  aggregation: string
  chart_config: string
  layout: string
  styling: string
  description: string
  preview: string
  tags: string
  industry: string
  language: string
  usage_count: number
  rating: number
  rating_count: number
  created_at: string
  updated_at: string
  created_by?: string
  
  // Relations
  company?: any
  creator?: any
}

export interface ReportExecution {
  id: string
  report_id: string
  status: 'pending' | 'running' | 'completed' | 'failed' | 'cancelled'
  parameters: string
  file_path: string
  file_size: number
  file_format: string
  row_count: number
  error_message: string
  execution_time: number
  is_scheduled: boolean
  scheduled_at?: string
  trigger_type: 'manual' | 'scheduled' | 'api'
  executed_by: string
  ip_address: string
  user_agent: string
  started_at: string
  completed_at?: string
  created_at: string
  
  // Relations
  report?: Report
  executed_by_user?: any
}

export interface ReportSubscription {
  id: string
  report_id: string
  user_id: string
  is_active: boolean
  email: string
  schedule: string
  file_format: string
  parameters: string
  delivery_method: string
  delivery_config: string
  last_delivered?: string
  delivery_count: number
  failure_count: number
  last_error: string
  created_at: string
  updated_at: string
  
  // Relations
  report?: Report
  user?: any
}

export interface ReportDashboard {
  id: string
  company_id: string
  name: string
  name_en: string
  layout: string
  theme: 'light' | 'dark' | 'auto'
  refresh_rate: number
  widgets: string
  filters: string
  is_public: boolean
  shared_with: string
  description: string
  tags: string
  is_default: boolean
  view_count: number
  last_viewed?: string
  created_at: string
  updated_at: string
  created_by: string
  updated_by?: string
  
  // Relations
  company?: any
  creator?: any
  updated_by_user?: any
}

export interface ReportDataSource {
  id: string
  company_id: string
  name: string
  type: 'database' | 'api' | 'file' | 'internal'
  connection_string: string
  credentials: string
  settings: string
  description: string
  schema: string
  sample_data: string
  status: 'active' | 'inactive' | 'error'
  last_tested?: string
  last_error: string
  test_result: string
  usage_count: number
  last_used?: string
  created_at: string
  updated_at: string
  created_by: string
  
  // Relations
  company?: any
  creator?: any
}

export interface ReportSchedule {
  id: string
  report_id: string
  name: string
  cron_expression: string
  timezone: string
  is_active: boolean
  parameters: string
  file_format: string
  recipients: string
  next_run?: string
  last_run?: string
  last_status: string
  run_count: number
  failure_count: number
  last_error: string
  created_at: string
  updated_at: string
  created_by: string
  
  // Relations
  report?: Report
  creator?: any
}

export interface BusinessKPI {
  id: string
  company_id: string
  name: string
  category: 'sales' | 'finance' | 'production' | 'inventory' | 'quality'
  formula: string
  data_sources: string
  filters: string
  unit: string
  target_value: number
  target_type: string
  threshold_high: number
  threshold_low: number
  display_format: string
  chart_type: string
  color_scheme: string
  description: string
  frequency: 'realtime' | 'hourly' | 'daily' | 'weekly' | 'monthly'
  is_active: boolean
  current_value: number
  previous_value: number
  trend: 'up' | 'down' | 'stable'
  last_updated?: string
  created_at: string
  updated_at: string
  created_by: string
  
  // Relations
  company?: any
  creator?: any
}

export interface ReportDashboardData {
  total_reports: number
  total_executions: number
  total_templates: number
  total_dashboards: number
  total_data_sources: number
  total_kpis: number
  scheduled_reports: number
  active_subscriptions: number
  reports_by_category: CategoryCount[]
  executions_by_status: StatusCount[]
  popular_reports: Report[]
  recent_executions: ReportExecution[]
  scheduled_reports_list: ReportSchedule[]
  system_health: Record<string, any>
}

export interface CategoryCount {
  category: string
  count: number
}

export interface StatusCount {
  status: string
  count: number
}

export interface DataSourceTestResult {
  success: boolean
  message: string
  connection_time: number
  sample_data: Record<string, any>[]
  schema: Record<string, any>
  error: string
}

export interface ImportResult {
  success: number
  failed: number
  errors: string[]
  imported_reports: string[]
}

export interface ReportColumn {
  name: string
  label: string
  data_type: string
  format: string
  width: number
  visible: boolean
  sortable: boolean
  filterable: boolean
  aggregable: boolean
}

export interface ReportSort {
  column: string
  direction: 'asc' | 'desc'
}

export interface DashboardWidget {
  id: string
  type: 'chart' | 'table' | 'kpi' | 'text'
  title: string
  position: Record<string, number>
  report_id?: string
  kpi_id?: string
  config: Record<string, any>
  data_source: Record<string, any>
}

export interface CreateReportRequest {
  name: string
  name_en?: string
  category: string
  type: string
  data_source?: Record<string, any>
  filters?: Record<string, any>
  columns?: ReportColumn[]
  sorting?: ReportSort[]
  grouping?: string[]
  aggregation?: Record<string, any>
  chart_config?: Record<string, any>
  template_id?: string
  layout?: Record<string, any>
  styling?: Record<string, any>
  is_public?: boolean
  shared_with?: string[]
  cache_enabled?: boolean
  cache_ttl?: number
  query_timeout?: number
  description?: string
  tags?: string[]
}

export interface UpdateReportRequest {
  name?: string
  name_en?: string
  category?: string
  type?: string
  status?: string
  data_source?: Record<string, any>
  filters?: Record<string, any>
  columns?: ReportColumn[]
  sorting?: ReportSort[]
  grouping?: string[]
  aggregation?: Record<string, any>
  chart_config?: Record<string, any>
  layout?: Record<string, any>
  styling?: Record<string, any>
  is_public?: boolean
  shared_with?: string[]
  cache_enabled?: boolean
  cache_ttl?: number
  query_timeout?: number
  description?: string
  tags?: string[]
}

export interface CreateReportTemplateRequest {
  name: string
  name_en?: string
  category: string
  type: string
  is_system_template?: boolean
  data_source?: Record<string, any>
  filters?: Record<string, any>
  columns?: ReportColumn[]
  sorting?: ReportSort[]
  grouping?: string[]
  aggregation?: Record<string, any>
  chart_config?: Record<string, any>
  layout?: Record<string, any>
  styling?: Record<string, any>
  description?: string
  preview?: string
  tags?: string[]
  industry?: string
  language?: string
}

export interface UpdateReportTemplateRequest {
  name?: string
  name_en?: string
  category?: string
  type?: string
  data_source?: Record<string, any>
  filters?: Record<string, any>
  columns?: ReportColumn[]
  sorting?: ReportSort[]
  grouping?: string[]
  aggregation?: Record<string, any>
  chart_config?: Record<string, any>
  layout?: Record<string, any>
  styling?: Record<string, any>
  description?: string
  preview?: string
  tags?: string[]
  industry?: string
  language?: string
}

export interface CreateReportDashboardRequest {
  name: string
  name_en?: string
  layout?: Record<string, any>
  theme?: string
  refresh_rate?: number
  widgets?: DashboardWidget[]
  filters?: Record<string, any>
  is_public?: boolean
  shared_with?: string[]
  description?: string
  tags?: string[]
  is_default?: boolean
}

export interface UpdateReportDashboardRequest {
  name?: string
  name_en?: string
  layout?: Record<string, any>
  theme?: string
  refresh_rate?: number
  widgets?: DashboardWidget[]
  filters?: Record<string, any>
  is_public?: boolean
  shared_with?: string[]
  description?: string
  tags?: string[]
  is_default?: boolean
}

export interface ListParams {
  page?: number
  page_size?: number
  search?: string
  category?: string
  type?: string
  status?: string
  is_scheduled?: boolean
  is_public?: boolean
  created_by?: string
  sort_by?: string
  sort_order?: string
  is_system_template?: boolean
  industry?: string
  language?: string
  trigger_type?: string
  executed_by?: string
  start_date?: string
  end_date?: string
  is_active?: boolean
  delivery_method?: string
  report_id?: string
  theme?: string
  is_default?: boolean
  frequency?: string
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
}

class ReportService {
  // Report operations
  async createReport(data: CreateReportRequest): Promise<Report> {
    const response = await api.post('/reports', data)
    return response.data
  }

  async updateReport(id: string, data: UpdateReportRequest): Promise<Report> {
    const response = await api.put(`/reports/${id}`, data)
    return response.data
  }

  async getReport(id: string): Promise<Report> {
    const response = await api.get(`/reports/${id}`)
    return response.data
  }

  async listReports(params: ListParams = {}): Promise<PaginatedResponse<Report>> {
    const response = await api.get('/reports', { params })
    return response.data
  }

  async deleteReport(id: string): Promise<void> {
    await api.delete(`/reports/${id}`)
  }

  async duplicateReport(id: string): Promise<Report> {
    const response = await api.post(`/reports/${id}/duplicate`)
    return response.data
  }

  // Report Template operations
  async createReportTemplate(data: CreateReportTemplateRequest): Promise<ReportTemplate> {
    const response = await api.post('/report-templates', data)
    return response.data
  }

  async updateReportTemplate(id: string, data: UpdateReportTemplateRequest): Promise<ReportTemplate> {
    const response = await api.put(`/report-templates/${id}`, data)
    return response.data
  }

  async getReportTemplate(id: string): Promise<ReportTemplate> {
    const response = await api.get(`/report-templates/${id}`)
    return response.data
  }

  async listReportTemplates(params: ListParams = {}): Promise<PaginatedResponse<ReportTemplate>> {
    const response = await api.get('/report-templates', { params })
    return response.data
  }

  async deleteReportTemplate(id: string): Promise<void> {
    await api.delete(`/report-templates/${id}`)
  }

  // Report Execution operations
  async executeReport(id: string, params: Record<string, any> = {}): Promise<ReportExecution> {
    const response = await api.post(`/reports/${id}/execute`, params)
    return response.data
  }

  async getReportExecution(id: string): Promise<ReportExecution> {
    const response = await api.get(`/report-executions/${id}`)
    return response.data
  }

  async listReportExecutions(reportId: string, params: ListParams = {}): Promise<PaginatedResponse<ReportExecution>> {
    const response = await api.get(`/reports/${reportId}/executions`, { params })
    return response.data
  }

  async cancelReportExecution(id: string): Promise<void> {
    await api.post(`/report-executions/${id}/cancel`)
  }

  async downloadReportResult(id: string): Promise<Blob> {
    const response = await api.get(`/report-executions/${id}/download`, {
      responseType: 'blob'
    })
    return response.data
  }

  // Report Subscription operations
  async getReportSubscription(id: string): Promise<ReportSubscription> {
    const response = await api.get(`/report-subscriptions/${id}`)
    return response.data
  }

  async listReportSubscriptions(params: ListParams = {}): Promise<PaginatedResponse<ReportSubscription>> {
    const response = await api.get('/report-subscriptions', { params })
    return response.data
  }

  async deleteReportSubscription(id: string): Promise<void> {
    await api.delete(`/report-subscriptions/${id}`)
  }

  // Report Dashboard operations
  async getReportDashboard(id: string): Promise<ReportDashboard> {
    const response = await api.get(`/report-dashboards/${id}`)
    return response.data
  }

  async listReportDashboards(params: ListParams = {}): Promise<PaginatedResponse<ReportDashboard>> {
    const response = await api.get('/report-dashboards', { params })
    return response.data
  }

  async deleteReportDashboard(id: string): Promise<void> {
    await api.delete(`/report-dashboards/${id}`)
  }

  // Report Data Source operations
  async getReportDataSource(id: string): Promise<ReportDataSource> {
    const response = await api.get(`/report-data-sources/${id}`)
    return response.data
  }

  async listReportDataSources(params: ListParams = {}): Promise<PaginatedResponse<ReportDataSource>> {
    const response = await api.get('/report-data-sources', { params })
    return response.data
  }

  async deleteReportDataSource(id: string): Promise<void> {
    await api.delete(`/report-data-sources/${id}`)
  }

  async testReportDataSource(id: string): Promise<DataSourceTestResult> {
    const response = await api.post(`/report-data-sources/${id}/test`)
    return response.data
  }

  // Report Schedule operations
  async getReportSchedule(id: string): Promise<ReportSchedule> {
    const response = await api.get(`/report-schedules/${id}`)
    return response.data
  }

  async listReportSchedules(params: ListParams = {}): Promise<PaginatedResponse<ReportSchedule>> {
    const response = await api.get('/report-schedules', { params })
    return response.data
  }

  async deleteReportSchedule(id: string): Promise<void> {
    await api.delete(`/report-schedules/${id}`)
  }

  // Business KPI operations
  async getBusinessKPI(id: string): Promise<BusinessKPI> {
    const response = await api.get(`/business-kpis/${id}`)
    return response.data
  }

  async listBusinessKPIs(params: ListParams = {}): Promise<PaginatedResponse<BusinessKPI>> {
    const response = await api.get('/business-kpis', { params })
    return response.data
  }

  async deleteBusinessKPI(id: string): Promise<void> {
    await api.delete(`/business-kpis/${id}`)
  }

  async updateKPIValues(): Promise<void> {
    await api.post('/business-kpis/update-values')
  }

  // Business operations
  async getReportDashboardData(): Promise<ReportDashboardData> {
    const response = await api.get('/reports/dashboard')
    return response.data
  }

  async getReportStatistics(): Promise<Record<string, any>> {
    const response = await api.get('/reports/statistics')
    return response.data
  }

  async getPopularReports(limit: number = 10): Promise<Report[]> {
    const response = await api.get('/reports/popular', { params: { limit } })
    return response.data
  }

  async getRecentExecutions(limit: number = 10): Promise<ReportExecution[]> {
    const response = await api.get('/reports/recent-executions', { params: { limit } })
    return response.data
  }

  async generateReportFromTemplate(templateId: string, params: Record<string, any> = {}): Promise<Report> {
    const response = await api.post(`/report-templates/${templateId}/generate`, params)
    return response.data
  }

  async exportReport(id: string, format: string = 'json', params: Record<string, any> = {}): Promise<Blob> {
    const response = await api.get(`/reports/${id}/export`, {
      params: { format, ...params },
      responseType: 'blob'
    })
    return response.data
  }

  async importReports(file: File): Promise<ImportResult> {
    const formData = new FormData()
    formData.append('file', file)
    
    const response = await api.post('/reports/import', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
    
    return response.data
  }

  // Utility functions
  async getReportCategories(): Promise<string[]> {
    const response = await api.get('/reports/categories')
    return response.data
  }

  async getReportTypes(): Promise<string[]> {
    const response = await api.get('/reports/types')
    return response.data
  }

  async getChartTypes(): Promise<string[]> {
    const response = await api.get('/reports/chart-types')
    return response.data
  }

  async getDataSourceTypes(): Promise<string[]> {
    const response = await api.get('/report-data-sources/types')
    return response.data
  }

  async getKPICategories(): Promise<string[]> {
    const response = await api.get('/business-kpis/categories')
    return response.data
  }

  async validateReportConfiguration(config: Record<string, any>): Promise<{ valid: boolean; errors: string[] }> {
    const response = await api.post('/reports/validate', config)
    return response.data
  }

  async previewReport(config: Record<string, any>): Promise<Record<string, any>> {
    const response = await api.post('/reports/preview', config)
    return response.data
  }

  // Additional methods for report execution and management
  async getKPIDashboard(period: string): Promise<any> {
    const response = await api.get('/reports/kpi-dashboard', { params: { period } })
    return response.data
  }

  async getBusinessKPIs(): Promise<any[]> {
    const response = await api.get('/business-kpis')
    return response.data
  }

  async duplicateReportTemplate(id: string): Promise<ReportTemplate> {
    const response = await api.post(`/report-templates/${id}/duplicate`)
    return response.data
  }

  async getReportExecutions(reportId: string, params: ListParams = {}): Promise<PaginatedResponse<ReportExecution>> {
    return this.listReportExecutions(reportId, params)
  }

  async exportExecution(id: string, format: string = 'pdf'): Promise<Blob> {
    const response = await api.get(`/report-executions/${id}/export`, {
      params: { format },
      responseType: 'blob'
    })
    return response.data
  }

  async listExecutions(params: ListParams = {}): Promise<PaginatedResponse<ReportExecution>> {
    const response = await api.get('/report-executions', { params })
    return response.data
  }

  async getExecutionStats(): Promise<any> {
    const response = await api.get('/report-executions/stats')
    return response.data
  }

  async rerunExecution(id: string): Promise<ReportExecution> {
    const response = await api.post(`/report-executions/${id}/rerun`)
    return response.data
  }

  async cancelExecution(id: string): Promise<void> {
    await api.post(`/report-executions/${id}/cancel`)
  }

  async exportExecutions(params: Record<string, any> = {}): Promise<Blob> {
    const response = await api.get('/report-executions/export', {
      params,
      responseType: 'blob'
    })
    return response.data
  }

  // Analytics Dashboard Methods
  async getDashboardData(dateRange: string): Promise<any> {
    const response = await api.get('/reports/dashboard', {
      params: { date_range: dateRange }
    })
    return response.data
  }

  async getQuoteTrends(dateRange: string): Promise<any> {
    const response = await api.get('/reports/quote-trends', {
      params: { date_range: dateRange }
    })
    return response.data
  }

  async getCustomerAnalytics(dateRange: string): Promise<any> {
    const response = await api.get('/reports/customer-analytics', {
      params: { date_range: dateRange }
    })
    return response.data
  }

  async getProcessAnalytics(dateRange: string): Promise<any> {
    const response = await api.get('/reports/process-analytics', {
      params: { date_range: dateRange }
    })
    return response.data
  }

  async exportReport(reportType: string, dateRange: string): Promise<Blob> {
    const response = await api.get(`/reports/export/${reportType}`, {
      params: { date_range: dateRange },
      responseType: 'blob'
    })
    return response.data
  }

  async getEngineerPerformance(engineerId?: string, dateRange?: string): Promise<any> {
    const response = await api.get('/reports/engineer-performance', {
      params: { engineer_id: engineerId, date_range: dateRange }
    })
    return response.data
  }

  async getSalesReport(dateRange: string): Promise<any> {
    const response = await api.get('/reports/sales', {
      params: { date_range: dateRange }
    })
    return response.data
  }

  async getCostAnalysisReport(dateRange: string): Promise<any> {
    const response = await api.get('/reports/cost-analysis', {
      params: { date_range: dateRange }
    })
    return response.data
  }

  async generateCustomReport(config: {
    reportType: string;
    filters: Record<string, any>;
    columns: string[];
    groupBy?: string;
    orderBy?: string;
  }): Promise<any> {
    const response = await api.post('/reports/custom', config)
    return response.data
  }
}

export const reportService = new ReportService()
export default reportService