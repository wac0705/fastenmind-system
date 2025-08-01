import api from './api'

export interface ProcessCategory {
  id: string
  category_code: string
  category_name: string
  category_name_en?: string
  process_order: number
  is_required: boolean
  created_at: string
  updated_at: string
}

export interface Equipment {
  id: string
  equipment_code: string
  equipment_name: string
  equipment_name_en?: string
  process_category_id: string
  floor_area: number // 佔地面積 (平方米)
  power_consumption: number // 耗電功率 (kW)
  max_capacity: number // 最大產能 (件/小時)
  depreciation_years: number
  purchase_cost: number
  maintenance_cost_annual: number
  is_active: boolean
  created_at: string
  updated_at: string
  
  // Relations
  process_category?: ProcessCategory
}

export interface ProcessCostConfig {
  id: string
  company_id: string
  process_category_id: string
  equipment_id: string
  labor_cost_per_hour: number
  electricity_rate: number
  land_cost_per_sqm: number
  overhead_rate: number
  effective_date: string
  expiry_date?: string
  is_active: boolean
  created_at: string
  updated_at: string
  
  // Relations
  process_category?: ProcessCategory
  equipment?: Equipment
}

export interface ProcessCostCalculation {
  process_category: string
  equipment_name: string
  quantity: number
  unit_time: number // 每件加工時間（秒）
  labor_cost: number
  equipment_cost: number
  electricity_cost: number
  land_cost: number
  maintenance_cost: number
  overhead_cost: number
  total_cost: number
  unit_cost: number
}

export interface CreateProcessCostConfigRequest {
  process_category_id: string
  equipment_id: string
  labor_cost_per_hour: number
  electricity_rate: number
  land_cost_per_sqm: number
  overhead_rate?: number
  effective_date: string
  is_active?: boolean
}

class ProcessService {
  // Process Categories
  async listCategories(): Promise<ProcessCategory[]> {
    const response = await api.get('/process-categories')
    return response.data
  }

  async getCategory(id: string): Promise<ProcessCategory> {
    const response = await api.get(`/process-categories/${id}`)
    return response.data
  }

  // Equipment
  async listEquipment(params?: {
    process_category_id?: string
    is_active?: boolean
  }): Promise<Equipment[]> {
    const response = await api.get('/equipment', { params })
    return response.data
  }

  async getEquipment(id: string): Promise<Equipment> {
    const response = await api.get(`/equipment/${id}`)
    return response.data
  }

  async createEquipment(data: Partial<Equipment>): Promise<Equipment> {
    const response = await api.post('/equipment', data)
    return response.data
  }

  async updateEquipment(id: string, data: Partial<Equipment>): Promise<Equipment> {
    const response = await api.put(`/equipment/${id}`, data)
    return response.data
  }

  // Process Cost Config
  async listCostConfigs(params?: {
    process_category_id?: string
    equipment_id?: string
    is_active?: boolean
  }): Promise<ProcessCostConfig[]> {
    const response = await api.get('/process-cost-configs', { params })
    return response.data
  }

  async getCostConfig(id: string): Promise<ProcessCostConfig> {
    const response = await api.get(`/process-cost-configs/${id}`)
    return response.data
  }

  async createCostConfig(data: CreateProcessCostConfigRequest): Promise<ProcessCostConfig> {
    const response = await api.post('/process-cost-configs', data)
    return response.data
  }

  async updateCostConfig(id: string, data: Partial<CreateProcessCostConfigRequest>): Promise<ProcessCostConfig> {
    const response = await api.put(`/process-cost-configs/${id}`, data)
    return response.data
  }

  // Cost Calculation
  async calculateProcessCost(params: {
    equipment_id: string
    quantity: number
    unit_time: number // seconds per unit
  }): Promise<ProcessCostCalculation> {
    const response = await api.post('/process-cost/calculate', params)
    return response.data
  }

  async calculateTotalCost(processes: Array<{
    equipment_id: string
    quantity: number
    unit_time: number
  }>): Promise<{
    processes: ProcessCostCalculation[]
    total_cost: number
    unit_cost: number
  }> {
    const response = await api.post('/process-cost/calculate-total', { processes })
    return response.data
  }
}

export const processService = new ProcessService()
export default processService