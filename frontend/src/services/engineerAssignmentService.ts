import apiClient from './api';

export interface EngineerAssignment {
  id: string;
  company_id: string;
  inquiry_id: string;
  engineer_id: string;
  assigned_by: string;
  assigned_at: string;
  status: 'pending' | 'in_progress' | 'completed' | 'cancelled';
  priority: 'low' | 'normal' | 'high' | 'urgent';
  due_date?: string;
  completed_at?: string;
  notes?: string;
  inquiry?: any;
  engineer?: any;
  assigner?: any;
}

export interface EngineerAvailability {
  engineer_id: string;
  engineer_name: string;
  department: string;
  expertise: string[];
  current_load: number;
  max_load: number;
  is_available: boolean;
  expertise_match: number;
}

export interface EngineerWorkload {
  engineer_id: string;
  engineer_name: string;
  total_assignments: number;
  pending: number;
  in_progress: number;
  completed: number;
  completed_on_time: number;
  overdue: number;
  avg_completion_time: number;
}

export interface AssignmentStats {
  period: string;
  total_assignments: number;
  status_breakdown: Record<string, number>;
  avg_completion_time: number;
  on_time_rate: number;
  engineer_stats: Array<{
    engineer_id: string;
    engineer_name: string;
    total_assigned: number;
    completed: number;
    avg_completion_time: number;
    completion_rate: number;
  }>;
  time_series: Array<{
    period: string;
    assigned: number;
    completed: number;
  }>;
}

export interface AssignmentHistory {
  id: string;
  assignment_id: string;
  action: string;
  from_engineer?: string;
  to_engineer?: string;
  from_status?: string;
  to_status?: string;
  action_by: string;
  action_at: string;
  reason?: string;
}

class EngineerAssignmentService {
  // 獲取分派列表
  async getAssignments(params?: {
    status?: string;
    engineer_id?: string;
    inquiry_id?: string;
    page?: number;
    limit?: number;
  }) {
    const queryParams = new URLSearchParams();
    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined) {
          queryParams.append(key, value.toString());
        }
      });
    }
    
    const response = await apiClient.get(`/engineer-assignments?${queryParams.toString()}`);
    return response.data;
  }

  // 獲取可用工程師
  async getAvailableEngineers(inquiryId?: string) {
    const url = inquiryId 
      ? `/engineer-assignments/available?inquiry_id=${inquiryId}`
      : '/engineer-assignments/available';
    
    const response = await apiClient.get(url);
    return response.data;
  }

  // 分派工程師
  async assignEngineer(data: {
    inquiry_id: string;
    engineer_id: string;
    priority?: string;
    due_date?: string;
    notes?: string;
  }) {
    const response = await apiClient.post('/engineer-assignments/assign', data);
    return response.data;
  }

  // 重新分派工程師
  async reassignEngineer(assignmentId: string, data: {
    new_engineer_id: string;
    reason: string;
  }) {
    const response = await apiClient.put(`/engineer-assignments/${assignmentId}/reassign`, data);
    return response.data;
  }

  // 自動分派工程師
  async autoAssignEngineer(data: {
    inquiry_id: string;
    rules?: {
      consider_workload?: boolean;
      consider_expertise?: boolean;
      consider_availability?: boolean;
      max_assignments?: number;
    };
  }) {
    const response = await apiClient.post('/engineer-assignments/auto-assign', data);
    return response.data;
  }

  // 獲取工程師工作負載
  async getEngineerWorkload(params?: {
    start_date?: string;
    end_date?: string;
  }) {
    const queryParams = new URLSearchParams();
    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined) {
          queryParams.append(key, value);
        }
      });
    }
    
    const response = await apiClient.get(`/engineer-assignments/workload?${queryParams.toString()}`);
    return response.data;
  }

  // 獲取分派歷史
  async getAssignmentHistory(params?: {
    inquiry_id?: string;
    engineer_id?: string;
    page?: number;
    limit?: number;
  }) {
    const queryParams = new URLSearchParams();
    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined) {
          queryParams.append(key, value.toString());
        }
      });
    }
    
    const response = await apiClient.get(`/engineer-assignments/history?${queryParams.toString()}`);
    return response.data;
  }

  // 獲取分派統計
  async getAssignmentStats(period: 'daily' | 'weekly' | 'monthly' = 'monthly') {
    const response = await apiClient.get(`/engineer-assignments/stats?period=${period}`);
    return response.data;
  }

  // 更新分派狀態
  async updateAssignmentStatus(assignmentId: string, data: {
    status: 'pending' | 'in_progress' | 'completed' | 'cancelled';
    notes?: string;
  }) {
    const response = await apiClient.put(`/engineer-assignments/${assignmentId}/status`, data);
    return response.data;
  }

  // 獲取單個分派詳情
  async getAssignment(assignmentId: string) {
    const response = await apiClient.get(`/engineer-assignments/${assignmentId}`);
    return response.data;
  }

  // 批量分派
  async batchAssign(data: {
    inquiry_ids: string[];
    engineer_id: string;
    priority?: string;
    notes?: string;
  }) {
    const response = await apiClient.post('/engineer-assignments/batch-assign', data);
    return response.data;
  }

  // 獲取工程師績效報告
  async getEngineerPerformance(engineerId: string, params?: {
    start_date?: string;
    end_date?: string;
  }) {
    const queryParams = new URLSearchParams();
    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined) {
          queryParams.append(key, value);
        }
      });
    }
    
    const response = await apiClient.get(`/engineer-assignments/engineers/${engineerId}/performance?${queryParams.toString()}`);
    return response.data;
  }

  // 導出分派報告
  async exportAssignmentReport(params: {
    format: 'excel' | 'pdf' | 'csv';
    start_date?: string;
    end_date?: string;
    engineer_id?: string;
    status?: string;
  }) {
    const queryParams = new URLSearchParams();
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined) {
        queryParams.append(key, value);
      }
    });
    
    const response = await apiClient.get(`/engineer-assignments/export?${queryParams.toString()}`, {
      responseType: 'blob'
    });
    
    // 創建下載連結
    const url = window.URL.createObjectURL(new Blob([response.data]));
    const link = document.createElement('a');
    link.href = url;
    link.setAttribute('download', `assignment_report_${Date.now()}.${params.format}`);
    document.body.appendChild(link);
    link.click();
    link.remove();
    
    return response.data;
  }
}

export const engineerAssignmentService = new EngineerAssignmentService();