import api from './api'

export interface N8NWorkflow {
  id: string
  name: string
  description?: string
  workflow_id: string // N8N workflow ID
  trigger_type: 'webhook' | 'schedule' | 'manual' | 'event'
  trigger_config: Record<string, any>
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface N8NExecution {
  id: string
  workflow_id: string
  execution_id: string // N8N execution ID
  status: 'running' | 'success' | 'error' | 'canceled'
  started_at: string
  finished_at?: string
  error_message?: string
  input_data?: Record<string, any>
  output_data?: Record<string, any>
}

export interface WorkflowTrigger {
  workflow_id: string
  data: Record<string, any>
  wait_for_completion?: boolean
}

export interface WebhookRegistration {
  name: string
  workflow_id: string
  event_types: string[]
  target_url?: string
  headers?: Record<string, string>
  is_active: boolean
}

export interface ScheduledTask {
  name: string
  workflow_id: string
  cron_expression: string
  timezone?: string
  data?: Record<string, any>
  is_active: boolean
}

// Predefined workflow templates
export const WORKFLOW_TEMPLATES = {
  INQUIRY_NOTIFICATION: {
    name: '詢價單通知',
    description: '新詢價單建立時通知相關人員',
    trigger_type: 'event',
    event: 'inquiry.created',
  },
  QUOTE_REMINDER: {
    name: '報價提醒',
    description: '報價到期前提醒業務人員',
    trigger_type: 'schedule',
    schedule: '0 9 * * *', // Daily at 9 AM
  },
  COST_UPDATE_ALERT: {
    name: '成本更新提醒',
    description: '定期提醒更新材料和製程成本',
    trigger_type: 'schedule',
    schedule: '0 9 1 * *', // Monthly on 1st at 9 AM
  },
  EXCHANGE_RATE_UPDATE: {
    name: '匯率更新',
    description: '每日更新匯率資料',
    trigger_type: 'schedule',
    schedule: '0 8 * * 1-5', // Weekdays at 8 AM
  },
  CUSTOMER_CREDIT_CHECK: {
    name: '客戶信用檢查',
    description: '檢查客戶信用額度並通知',
    trigger_type: 'event',
    event: 'quote.created',
  },
  APPROVAL_WORKFLOW: {
    name: '審核流程',
    description: '報價審核自動化流程',
    trigger_type: 'event',
    event: 'quote.submitted_for_review',
  },
}

class N8NService {
  // Workflow Management
  async listWorkflows(params?: {
    page?: number
    page_size?: number
    is_active?: boolean
  }): Promise<{ data: N8NWorkflow[]; pagination: any }> {
    const response = await api.get('/n8n/workflows', { params })
    return response.data
  }

  async getWorkflow(id: string): Promise<N8NWorkflow> {
    const response = await api.get(`/n8n/workflows/${id}`)
    return response.data
  }

  async createWorkflow(data: {
    name: string
    description?: string
    workflow_id: string
    trigger_type: string
    trigger_config: Record<string, any>
  }): Promise<N8NWorkflow> {
    const response = await api.post('/n8n/workflows', data)
    return response.data
  }

  async updateWorkflow(id: string, data: Partial<N8NWorkflow>): Promise<N8NWorkflow> {
    const response = await api.put(`/n8n/workflows/${id}`, data)
    return response.data
  }

  async deleteWorkflow(id: string): Promise<void> {
    await api.delete(`/n8n/workflows/${id}`)
  }

  async toggleWorkflow(id: string, active: boolean): Promise<N8NWorkflow> {
    const response = await api.patch(`/n8n/workflows/${id}/toggle`, { is_active: active })
    return response.data
  }

  // Workflow Execution
  async triggerWorkflow(data: WorkflowTrigger): Promise<N8NExecution> {
    const response = await api.post('/n8n/trigger', data)
    return response.data
  }

  async getExecutions(workflowId?: string, params?: {
    page?: number
    page_size?: number
    status?: string
    from_date?: string
    to_date?: string
  }): Promise<{ data: N8NExecution[]; pagination: any }> {
    const url = workflowId 
      ? `/n8n/workflows/${workflowId}/executions`
      : '/n8n/executions'
    const response = await api.get(url, { params })
    return response.data
  }

  async getExecution(id: string): Promise<N8NExecution> {
    const response = await api.get(`/n8n/executions/${id}`)
    return response.data
  }

  async cancelExecution(id: string): Promise<void> {
    await api.post(`/n8n/executions/${id}/cancel`)
  }

  // Webhook Management
  async listWebhooks(): Promise<WebhookRegistration[]> {
    const response = await api.get('/n8n/webhooks')
    return response.data
  }

  async registerWebhook(data: WebhookRegistration): Promise<WebhookRegistration> {
    const response = await api.post('/n8n/webhooks', data)
    return response.data
  }

  async updateWebhook(id: string, data: Partial<WebhookRegistration>): Promise<WebhookRegistration> {
    const response = await api.put(`/n8n/webhooks/${id}`, data)
    return response.data
  }

  async deleteWebhook(id: string): Promise<void> {
    await api.delete(`/n8n/webhooks/${id}`)
  }

  // Scheduled Tasks
  async listScheduledTasks(): Promise<ScheduledTask[]> {
    const response = await api.get('/n8n/scheduled-tasks')
    return response.data
  }

  async createScheduledTask(data: ScheduledTask): Promise<ScheduledTask> {
    const response = await api.post('/n8n/scheduled-tasks', data)
    return response.data
  }

  async updateScheduledTask(id: string, data: Partial<ScheduledTask>): Promise<ScheduledTask> {
    const response = await api.put(`/n8n/scheduled-tasks/${id}`, data)
    return response.data
  }

  async deleteScheduledTask(id: string): Promise<void> {
    await api.delete(`/n8n/scheduled-tasks/${id}`)
  }

  // Integration Helpers
  async testConnection(): Promise<{ connected: boolean; version?: string; message?: string }> {
    const response = await api.get('/n8n/test-connection')
    return response.data
  }

  async getAvailableWorkflows(): Promise<Array<{ id: string; name: string; tags: string[] }>> {
    const response = await api.get('/n8n/available-workflows')
    return response.data
  }

  // Business-specific triggers
  async notifyInquiryCreated(inquiryId: string): Promise<void> {
    await this.triggerWorkflow({
      workflow_id: 'inquiry_notification',
      data: {
        event: 'inquiry.created',
        inquiry_id: inquiryId,
        timestamp: new Date().toISOString(),
      },
    })
  }

  async notifyQuoteReview(quoteId: string, reviewerId: string): Promise<void> {
    await this.triggerWorkflow({
      workflow_id: 'quote_review',
      data: {
        event: 'quote.review_requested',
        quote_id: quoteId,
        reviewer_id: reviewerId,
        timestamp: new Date().toISOString(),
      },
    })
  }

  async notifyQuoteApproved(quoteId: string, approverId: string): Promise<void> {
    await this.triggerWorkflow({
      workflow_id: 'quote_approved',
      data: {
        event: 'quote.approved',
        quote_id: quoteId,
        approver_id: approverId,
        timestamp: new Date().toISOString(),
      },
    })
  }

  async checkCreditLimit(customerId: string, amount: number): Promise<void> {
    await this.triggerWorkflow({
      workflow_id: 'credit_check',
      data: {
        event: 'credit.check_requested',
        customer_id: customerId,
        amount: amount,
        timestamp: new Date().toISOString(),
      },
    })
  }

  async scheduleFollowUp(entityType: 'inquiry' | 'quote', entityId: string, followUpDate: string): Promise<void> {
    await this.triggerWorkflow({
      workflow_id: 'follow_up_reminder',
      data: {
        event: 'follow_up.scheduled',
        entity_type: entityType,
        entity_id: entityId,
        follow_up_date: followUpDate,
        timestamp: new Date().toISOString(),
      },
    })
  }
}

export const n8nService = new N8NService()
export default n8nService