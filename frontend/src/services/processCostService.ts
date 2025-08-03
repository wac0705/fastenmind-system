import apiClient from './api';

export interface ProcessCostTemplate {
  id: string;
  company_id: string;
  name: string;
  description: string;
  process_type: string;
  category: string;
  parameters: any;
  total_cost: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface MaterialCost {
  id: string;
  company_id: string;
  name: string;
  type: string;
  specification: string;
  unit_price: number;
  currency: string;
  unit: string;
  density: number;
  supplier: string;
  supplier_id?: string;
  lead_time?: number;
  min_order_quantity?: number;
  price_valid_until?: string;
  price_change_reason?: string;
}

export interface ProcessingRate {
  id: string;
  company_id: string;
  process_type: string;
  equipment_id?: string;
  equipment_name?: string;
  hourly_rate: number;
  setup_cost: number;
  minimum_charge: number;
  currency: string;
  is_active: boolean;
}

export interface ProcessCostCalculationRequest {
  inquiry_id?: string;
  product_id?: string;
  product_spec: {
    length: number;
    width: number;
    height: number;
    diameter?: number;
    thickness?: number;
    weight?: number;
    complexity: 'low' | 'medium' | 'high';
    custom_specs?: Record<string, number>;
  };
  material_id: string;
  material_utilization?: number;
  quantity: number;
  processes: Array<{
    process_type: string;
    equipment_id?: string;
    parameters?: Record<string, any>;
    sequence: number;
  }>;
  surface_treatment?: string;
  overhead_rate?: number;
  profit_margin?: number;
  base_currency?: string;
  target_currency?: string;
}

export interface ProcessCostResult {
  request_id: string;
  material_cost: number;
  processing_cost: number;
  surface_treatment_cost: number;
  packaging_cost: number;
  overhead_cost: number;
  total_cost: number;
  unit_cost: number;
  profit_amount: number;
  final_price: number;
  final_unit_price: number;
  currency: string;
  exchange_rate?: number;
  converted_price?: number;
  converted_unit_price?: number;
  details: Array<{
    type: string;
    description: string;
    unit_cost: number;
    quantity: number;
    total_cost: number;
  }>;
  calculated_at: string;
}

export interface CostSettings {
  id: string;
  company_id: string;
  default_overhead_rate: number;
  default_profit_margin: number;
  material_markup: number;
  labor_cost_multiplier: number;
  auto_update_prices: boolean;
  price_update_frequency: string;
}

class ProcessCostService {
  // 成本模板
  async getCostTemplates(params?: {
    process_type?: string;
    category?: string;
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
    
    const response = await apiClient.get(`/process-costs/templates?${queryParams.toString()}`);
    return response.data;
  }

  async createCostTemplate(template: Omit<ProcessCostTemplate, 'id' | 'created_at' | 'updated_at'>) {
    const response = await apiClient.post('/process-costs/templates', template);
    return response.data;
  }

  async updateCostTemplate(id: string, template: Partial<ProcessCostTemplate>) {
    const response = await apiClient.put(`/process-costs/templates/${id}`, template);
    return response.data;
  }

  async deleteCostTemplate(id: string) {
    const response = await apiClient.delete(`/process-costs/templates/${id}`);
    return response.data;
  }

  // 成本計算
  async calculateProcessCost(request: ProcessCostCalculationRequest) {
    const response = await apiClient.post('/process-costs/calculate', request);
    return response.data;
  }

  async batchCalculateCost(items: ProcessCostCalculationRequest[]) {
    const response = await apiClient.post('/process-costs/batch-calculate', { items });
    return response.data;
  }

  // 成本歷史
  async getCostHistory(params?: {
    inquiry_id?: string;
    product_id?: string;
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
    
    const response = await apiClient.get(`/process-costs/history?${queryParams.toString()}`);
    return response.data;
  }

  // 材料成本
  async getMaterialCosts(params?: {
    type?: string;
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
    
    const response = await apiClient.get(`/process-costs/materials?${queryParams.toString()}`);
    return response.data;
  }

  async updateMaterialCost(id: string, material: Partial<MaterialCost>) {
    const response = await apiClient.put(`/process-costs/materials/${id}`, material);
    return response.data;
  }

  async createMaterial(material: Omit<MaterialCost, 'id'>) {
    const response = await apiClient.post('/process-costs/materials', material);
    return response.data;
  }

  // 加工費率
  async getProcessingRates(params?: {
    process_type?: string;
    equipment_id?: string;
  }) {
    const queryParams = new URLSearchParams();
    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined) {
          queryParams.append(key, value);
        }
      });
    }
    
    const response = await apiClient.get(`/process-costs/processing-rates?${queryParams.toString()}`);
    return response.data;
  }

  async updateProcessingRate(id: string, rate: Partial<ProcessingRate>) {
    const response = await apiClient.put(`/process-costs/processing-rates/${id}`, rate);
    return response.data;
  }

  async createProcessingRate(rate: Omit<ProcessingRate, 'id'>) {
    const response = await apiClient.post('/process-costs/processing-rates', rate);
    return response.data;
  }

  // 成本分析
  async getCostAnalysis(params: {
    type?: string;
    period?: string;
    start_date?: string;
    end_date?: string;
  }) {
    const queryParams = new URLSearchParams();
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined) {
        queryParams.append(key, value);
      }
    });
    
    const response = await apiClient.get(`/process-costs/analysis?${queryParams.toString()}`);
    return response.data;
  }

  // 成本報告導出
  async exportCostReport(params: {
    format: 'excel' | 'pdf' | 'csv';
    type?: string;
    start_date?: string;
    end_date?: string;
  }) {
    const queryParams = new URLSearchParams();
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined) {
        queryParams.append(key, value);
      }
    });
    
    const response = await apiClient.get(`/process-costs/export?${queryParams.toString()}`, {
      responseType: 'blob'
    });
    
    // 創建下載連結
    const url = window.URL.createObjectURL(new Blob([response.data]));
    const link = document.createElement('a');
    link.href = url;
    link.setAttribute('download', `cost_report_${Date.now()}.${params.format}`);
    document.body.appendChild(link);
    link.click();
    link.remove();
    
    return response.data;
  }

  // 成本設定
  async getCostSettings() {
    const response = await apiClient.get('/process-costs/settings');
    return response.data;
  }

  async updateCostSettings(settings: Partial<CostSettings>) {
    const response = await apiClient.put('/process-costs/settings', settings);
    return response.data;
  }

  // 搜尋材料
  async searchMaterials(keyword: string) {
    const response = await apiClient.get(`/process-costs/materials/search?keyword=${encodeURIComponent(keyword)}`);
    return response.data;
  }

  // 獲取表面處理選項
  async getSurfaceTreatmentOptions() {
    const response = await apiClient.get('/process-costs/surface-treatments');
    return response.data;
  }

  // 批量更新材料價格
  async batchUpdateMaterialPrices(updates: Array<{
    material_id: string;
    new_price: number;
    reason: string;
  }>) {
    const response = await apiClient.post('/process-costs/materials/batch-update-prices', { updates });
    return response.data;
  }

  // 複製成本模板
  async duplicateCostTemplate(templateId: string, newName: string) {
    const response = await apiClient.post(`/process-costs/templates/${templateId}/duplicate`, { name: newName });
    return response.data;
  }
}

export const processCostService = new ProcessCostService();