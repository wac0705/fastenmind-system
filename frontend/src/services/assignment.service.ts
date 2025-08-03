import api from './api'

export interface EngineerCapability {
  id: string;
  engineer_id: string;
  engineer?: {
    id: string;
    full_name: string;
    email: string;
  };
  product_category: string;
  process_type: string;
  skill_level: number;
  max_concurrent_inquiries: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface AssignmentRule {
  id: string;
  rule_name: string;
  rule_type: 'auto' | 'rotation' | 'load_balance' | 'skill_based';
  priority: number;
  conditions: {
    product_categories?: string[];
    customer_types?: string[];
    min_skill_level?: number;
    auto_assign?: boolean;
    amount_range?: {
      min?: number;
      max?: number;
    };
  };
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface EngineerWorkloadStats {
  engineer_id: string;
  engineer_name: string;
  current_inquiries: number;
  completed_today: number;
  completed_this_week: number;
  completed_this_month: number;
  average_completion_hours: number;
  skill_categories: string[];
}

export interface AssignmentHistory {
  id: string;
  inquiry_id: string;
  inquiry?: {
    inquiry_no: string;
    product_name: string;
    customer_name: string;
  };
  assigned_from?: string;
  assigned_from_user?: {
    full_name: string;
  };
  assigned_to: string;
  assigned_to_user: {
    full_name: string;
  };
  assigned_by?: string;
  assigned_by_user?: {
    full_name: string;
  };
  assignment_type: 'auto' | 'manual' | 'reassign' | 'self_select';
  assignment_reason: string;
  rule_id?: string;
  rule?: AssignmentRule;
  assigned_at: string;
}

export interface EngineerPreference {
  id: string;
  engineer_id: string;
  preferred_categories: string[];
  preferred_customers: string[];
  max_daily_assignments: number;
  auto_accept_enabled: boolean;
  notification_enabled: boolean;
  created_at: string;
  updated_at: string;
}

export interface AssignmentRequest {
  inquiry_id: string;
  engineer_id: string;
  reason?: string;
  assignment_type: 'manual' | 'reassign';
}

export interface AutoAssignmentRequest {
  inquiry_id: string;
}

class AssignmentService {
  // 自動分派
  async autoAssign(inquiryId: string): Promise<AssignmentHistory> {
    const response = await api.post('/assignments/auto', {
      inquiry_id: inquiryId
    });
    return response.data;
  }

  // 手動分派
  async manualAssign(request: AssignmentRequest): Promise<AssignmentHistory> {
    const response = await api.post('/assignments/manual', request);
    return response.data;
  }

  // 分派工程師
  async assignEngineer(inquiryId: string, engineerId: string, notes?: string): Promise<AssignmentHistory> {
    return this.manualAssign({
      inquiry_id: inquiryId,
      engineer_id: engineerId,
      reason: notes,
      assignment_type: 'manual'
    });
  }

  // 工程師自選
  async selfSelect(inquiryId: string): Promise<AssignmentHistory> {
    const response = await api.post(`/assignments/self-select/${inquiryId}`);
    return response.data;
  }

  // 獲取工作量統計
  async getWorkloadStats(): Promise<EngineerWorkloadStats[]> {
    const response = await api.get('/assignments/workload-stats');
    return response.data;
  }

  // 獲取分派歷史
  async getAssignmentHistory(params?: {
    inquiry_id?: string;
    engineer_id?: string;
    limit?: number;
  }): Promise<AssignmentHistory[]> {
    const response = await api.get('/assignments/history', { params });
    return response.data;
  }

  // 獲取分派規則
  async getAssignmentRules(): Promise<AssignmentRule[]> {
    const response = await api.get('/assignments/rules');
    return response.data;
  }

  // 更新分派規則
  async updateAssignmentRule(id: string, rule: Partial<AssignmentRule>): Promise<void> {
    await api.put(`/assignments/rules/${id}`, rule);
  }

  // 獲取工程師能力
  async getEngineerCapabilities(engineerId: string): Promise<EngineerCapability[]> {
    const response = await api.get(`/assignments/engineers/${engineerId}/capabilities`);
    return response.data;
  }

  // 更新工程師能力
  async updateEngineerCapability(capability: Partial<EngineerCapability>): Promise<void> {
    await api.put('/assignments/capabilities', capability);
  }

  // 更新工程師偏好
  async updateEngineerPreference(preference: Partial<EngineerPreference>): Promise<void> {
    await api.put('/assignments/preferences', preference);
  }

  // 獲取可用工程師列表（根據產品類別）
  async getAvailableEngineers(productCategory: string): Promise<EngineerWorkloadStats[]> {
    const stats = await this.getWorkloadStats();
    // 過濾出有能力處理該產品類別的工程師
    return stats.filter(engineer => 
      engineer.skill_categories.includes(productCategory)
    ).sort((a, b) => a.current_inquiries - b.current_inquiries);
  }

  // 建議工程師（基於規則和工作量）
  async suggestEngineer(inquiryId: string): Promise<{
    suggested_engineer: EngineerWorkloadStats | null;
    matching_rules: AssignmentRule[];
    reason: string;
  }> {
    // 這個邏輯可以在前端實現，或呼叫後端專門的建議 API
    try {
      const response = await api.get(`/inquiries/${inquiryId}`);
      const inquiry = response.data;
      
      const availableEngineers = await this.getAvailableEngineers(inquiry.product_category);
      const rules = await this.getAssignmentRules();
      
      // 找出匹配的規則
      const matchingRules = rules.filter(rule => 
        rule.is_active && 
        rule.conditions.product_categories?.includes(inquiry.product_category)
      );
      
      if (availableEngineers.length > 0) {
        return {
          suggested_engineer: availableEngineers[0],
          matching_rules: matchingRules,
          reason: '根據工作量和技能匹配建議'
        };
      }
      
      return {
        suggested_engineer: null,
        matching_rules: matchingRules,
        reason: '沒有找到合適的工程師'
      };
    } catch (error) {
      console.error('Error suggesting engineer:', error);
      throw error;
    }
  }
}

export const assignmentService = new AssignmentService()
export default assignmentService