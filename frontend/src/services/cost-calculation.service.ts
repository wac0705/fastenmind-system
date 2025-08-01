import api from './api';

export interface ProcessCategory {
  id: string;
  code: string;
  name: string;
  name_en: string;
  description?: string;
  sort_order: number;
  is_active: boolean;
}

export interface Equipment {
  id: string;
  code: string;
  name: string;
  name_en?: string;
  process_category_id: string;
  process_category?: ProcessCategory;
  specs?: string;
  capacity_per_hour: number;
  power_consumption: number;
  depreciation_years: number;
  purchase_cost: number;
  maintenance_cost_per_year: number;
  location?: string;
  is_active: boolean;
}

export interface ProcessStep {
  id: string;
  code: string;
  name: string;
  name_en?: string;
  process_category_id: string;
  process_category?: ProcessCategory;
  default_equipment_id?: string;
  default_equipment?: Equipment;
  setup_time_minutes: number;
  cycle_time_seconds: number;
  labor_required: number;
  description?: string;
  sort_order: number;
  is_active: boolean;
}

export interface ProcessRouteDetail {
  id: string;
  route_id: string;
  sequence: number;
  process_step_id: string;
  process_step?: ProcessStep;
  equipment_id?: string;
  equipment?: Equipment;
  setup_time_override?: number;
  cycle_time_override?: number;
  yield_rate: number;
  notes?: string;
}

export interface ProductProcessRoute {
  id: string;
  product_category: string;
  material_type?: string;
  size_range?: string;
  route_name: string;
  is_default: boolean;
  is_active: boolean;
  route_details: ProcessRouteDetail[];
}

export interface CostParameter {
  id: string;
  parameter_type: string;
  parameter_name: string;
  value: number;
  unit?: string;
  effective_date: string;
  end_date?: string;
  description?: string;
}

export interface CostCalculation {
  id: string;
  inquiry_id?: string;
  quote_id?: string;
  calculation_no: string;
  product_name: string;
  quantity: number;
  material_cost: number;
  process_cost: number;
  overhead_cost: number;
  total_cost: number;
  unit_cost: number;
  margin_percentage: number;
  selling_price: number;
  route_id?: string;
  route?: ProductProcessRoute;
  calculated_by: string;
  calculated_by_user?: {
    id: string;
    full_name: string;
  };
  calculated_at: string;
  approved_by?: string;
  approved_by_user?: {
    id: string;
    full_name: string;
  };
  approved_at?: string;
  status: 'draft' | 'submitted' | 'approved' | 'rejected';
  details?: CostCalculationDetail[];
  created_at: string;
  updated_at: string;
}

export interface CostCalculationDetail {
  id: string;
  calculation_id: string;
  sequence: number;
  process_step_id: string;
  process_step?: ProcessStep;
  equipment_id?: string;
  equipment?: Equipment;
  setup_time: number;
  cycle_time: number;
  total_time_hours: number;
  labor_cost: number;
  equipment_cost: number;
  electricity_cost: number;
  other_cost: number;
  subtotal_cost: number;
  yield_loss_cost: number;
  notes?: string;
}

export interface CostCalculationRequest {
  inquiry_id?: string;
  product_name: string;
  product_category: string;
  material_type?: string;
  size_range?: string;
  quantity: number;
  material_cost?: number;
  route_id?: string;
  custom_route?: {
    process_step_id: string;
    equipment_id: string;
    setup_time?: number;
    cycle_time?: number;
  }[];
  margin_percentage?: number;
}

export interface CostSummary {
  material_cost: number;
  process_cost: number;
  overhead_cost: number;
  total_cost: number;
  unit_cost: number;
  suggested_price: number;
  margin_percentage: number;
  process_breakdown: {
    process_name: string;
    equipment_name?: string;
    total_time_hours: number;
    labor_cost: number;
    equipment_cost: number;
    electricity_cost: number;
    total_cost: number;
  }[];
}

class CostCalculationService {
  // 計算成本
  async calculateCost(request: CostCalculationRequest): Promise<CostCalculation> {
    const response = await api.post('/cost-calculations/calculate', request);
    return response.data;
  }

  // 獲取成本摘要
  async getCostSummary(calculationId: string): Promise<CostSummary> {
    const response = await api.get(`/cost-calculations/${calculationId}/summary`);
    return response.data;
  }

  // 獲取製程路線
  async getProcessRoutes(productCategory?: string): Promise<ProductProcessRoute[]> {
    const params = productCategory ? { product_category: productCategory } : {};
    const response = await api.get('/cost-calculations/process-routes', { params });
    return response.data;
  }

  // 獲取製程步驟
  async getProcessSteps(): Promise<ProcessStep[]> {
    const response = await api.get('/cost-calculations/process-steps');
    return response.data;
  }

  // 獲取設備列表
  async getEquipment(categoryId?: string): Promise<Equipment[]> {
    const params = categoryId ? { category_id: categoryId } : {};
    const response = await api.get('/cost-calculations/equipment', { params });
    return response.data;
  }

  // 獲取成本計算列表
  async getCalculations(params?: {
    page?: number;
    page_size?: number;
    status?: string;
  }): Promise<{ data: CostCalculation[]; pagination: any }> {
    const response = await api.get('/cost-calculations', { params });
    return response.data;
  }

  // 獲取單個成本計算
  async getCalculation(id: string): Promise<CostCalculation> {
    const response = await api.get(`/cost-calculations/${id}`);
    return response.data;
  }

  // 審核成本計算
  async approveCalculation(id: string): Promise<void> {
    await api.post(`/cost-calculations/${id}/approve`);
  }

  // 獲取成本參數
  async getCostParameters(): Promise<CostParameter[]> {
    const response = await api.get('/cost-calculations/parameters');
    return response.data;
  }

  // 更新成本參數
  async updateCostParameter(parameter: Partial<CostParameter>): Promise<void> {
    await api.put('/cost-calculations/parameters', parameter);
  }

  // 建議製程路線
  async suggestRoute(productCategory: string, materialType?: string, sizeRange?: string): Promise<ProductProcessRoute | null> {
    const routes = await this.getProcessRoutes(productCategory);
    
    // 優先找預設路線
    const defaultRoute = routes.find(r => r.is_default);
    if (defaultRoute) return defaultRoute;

    // 根據材料和尺寸匹配
    if (materialType || sizeRange) {
      const matchedRoute = routes.find(r => 
        (!materialType || r.material_type === materialType) &&
        (!sizeRange || r.size_range === sizeRange)
      );
      if (matchedRoute) return matchedRoute;
    }

    // 返回第一個可用路線
    return routes.length > 0 ? routes[0] : null;
  }
}

export default new CostCalculationService();