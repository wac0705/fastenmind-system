import api from './api'

// Types
export interface ProductionOrder {
  id: string
  company_id: string
  order_no: string
  status: 'planned' | 'released' | 'in_progress' | 'quality_check' | 'completed' | 'cancelled'
  priority: 'low' | 'medium' | 'high' | 'urgent'
  sales_order_id?: string
  customer_id?: string
  inventory_id: string
  product_name: string
  product_spec: string
  planned_quantity: number
  produced_quantity: number
  qualified_quantity: number
  defect_quantity: number
  unit: string
  planned_start_date: string
  planned_end_date: string
  actual_start_date?: string
  actual_end_date?: string
  route_id?: string
  current_station_id?: string
  completed_stations: number
  total_stations: number
  estimated_cost: number
  actual_cost: number
  material_cost: number
  labor_cost: number
  overhead_cost: number
  currency: string
  notes: string
  created_at: string
  updated_at: string
  created_by: string
  
  // Relations
  sales_order?: any
  customer?: any
  inventory?: any
  route?: ProductionRoute
  current_station?: WorkStation
  creator?: any
}

export interface ProductionRoute {
  id: string
  company_id: string
  route_no: string
  name: string
  description: string
  status: 'active' | 'inactive'
  inventory_id?: string
  product_category: string
  total_stations: number
  estimated_duration: number
  estimated_cost: number
  version: number
  is_active: boolean
  created_at: string
  updated_at: string
  created_by: string
  
  // Relations
  inventory?: any
  creator?: any
  operations?: RouteOperation[]
}

export interface RouteOperation {
  id: string
  route_id: string
  operation_no: number
  work_station_id: string
  name: string
  description: string
  instructions: string
  setup_time: number
  process_time: number
  teardown_time: number
  qc_required: boolean
  qc_instructions: string
  next_operation_id?: string
  created_at: string
  updated_at: string
  
  // Relations
  work_station?: WorkStation
  next_operation?: RouteOperation
}

export interface WorkStation {
  id: string
  company_id: string
  station_no: string
  name: string
  type: 'machine' | 'manual' | 'inspection' | 'assembly'
  status: 'available' | 'busy' | 'maintenance' | 'breakdown'
  capacity: number
  utilization_rate: number
  location: string
  department: string
  hourly_cost: number
  model: string
  manufacturer: string
  serial_number: string
  purchase_date?: string
  last_maintenance?: string
  next_maintenance?: string
  maintenance_notes: string
  description: string
  notes: string
  created_at: string
  updated_at: string
  created_by: string
  
  // Relations
  creator?: any
}

export interface ProductionTask {
  id: string
  production_order_id: string
  route_operation_id: string
  work_station_id: string
  task_no: number
  name: string
  status: 'pending' | 'in_progress' | 'completed' | 'on_hold' | 'cancelled'
  assigned_to?: string
  assigned_at?: string
  planned_quantity: number
  completed_quantity: number
  qualified_quantity: number
  defect_quantity: number
  planned_start_time: string
  planned_end_time: string
  actual_start_time?: string
  actual_end_time?: string
  qc_status: 'not_required' | 'pending' | 'passed' | 'failed'
  qc_notes: string
  qc_by?: string
  qc_at?: string
  notes: string
  issues: string
  created_at: string
  updated_at: string
  
  // Relations
  production_order?: ProductionOrder
  route_operation?: RouteOperation
  work_station?: WorkStation
  assigned_user?: any
  qc_user?: any
}

export interface ProductionMaterial {
  id: string
  production_order_id: string
  inventory_id: string
  planned_quantity: number
  issued_quantity: number
  consumed_quantity: number
  returned_quantity: number
  unit: string
  unit_cost: number
  total_cost: number
  status: 'planned' | 'issued' | 'consumed' | 'returned'
  issued_at?: string
  consumed_at?: string
  created_at: string
  updated_at: string
  
  // Relations
  inventory?: any
}

export interface QualityInspection {
  id: string
  company_id: string
  inspection_no: string
  type: 'incoming' | 'in_process' | 'final' | 'customer_return'
  status: 'pending' | 'in_progress' | 'passed' | 'failed' | 'on_hold'
  production_order_id?: string
  production_task_id?: string
  inventory_id?: string
  inspected_quantity: number
  qualified_quantity: number
  defect_quantity: number
  unit: string
  defect_types: string
  defect_reasons: string
  critical_defects: number
  major_defects: number
  minor_defects: number
  inspector_id: string
  inspected_at: string
  approved_by?: string
  approved_at?: string
  inspection_notes: string
  corrective_action: string
  attachment_path: string
  created_at: string
  updated_at: string
  
  // Relations
  production_order?: ProductionOrder
  production_task?: ProductionTask
  inventory?: any
  inspector?: any
  approver?: any
}

export interface ProductionDashboard {
  total_orders: number
  planned_orders: number
  released_orders: number
  in_progress_orders: number
  completed_orders: number
  cancelled_orders: number
  total_planned_quantity: number
  total_produced_quantity: number
  total_qualified_quantity: number
  total_defect_quantity: number
  production_efficiency: number
  quality_rate: number
}

export interface ProductionStats {
  total_work_stations: number
  available_stations: number
  busy_stations: number
  maintenance_stations: number
  breakdown_stations: number
  total_inspections: number
  passed_inspections: number
  failed_inspections: number
  pending_inspections: number
}

export interface CreateProductionOrderRequest {
  sales_order_id?: string
  customer_id?: string
  inventory_id: string
  product_name: string
  product_spec: string
  planned_quantity: number
  unit: string
  planned_start_date: string
  planned_end_date: string
  route_id?: string
  priority: string
  notes: string
}

export interface CreateProductionRouteRequest {
  name: string
  description: string
  inventory_id?: string
  product_category: string
  estimated_duration: number
  estimated_cost: number
}

export interface CreateWorkStationRequest {
  name: string
  type: string
  capacity: number
  location: string
  department: string
  hourly_cost: number
  model: string
  manufacturer: string
  serial_number: string
  purchase_date?: string
  description: string
  notes: string
}

export interface CreateQualityInspectionRequest {
  type: string
  production_order_id?: string
  production_task_id?: string
  inventory_id?: string
  inspected_quantity: number
  qualified_quantity: number
  defect_quantity: number
  unit: string
  defect_types: string
  defect_reasons: string
  critical_defects: number
  major_defects: number
  minor_defects: number
  inspected_at: string
  inspection_notes: string
  corrective_action: string
}

export interface ListParams {
  page?: number
  page_size?: number
  search?: string
  status?: string
  priority?: string
  customer_id?: string
  inventory_id?: string
  start_date?: string
  end_date?: string
  type?: string
  department?: string
  assigned_to?: string
  work_station_id?: string
  production_order_id?: string
  inspector_id?: string
  product_category?: string
}

export interface PaginatedResponse<T> {
  data: T[]
  pagination: {
    page: number
    page_size: number
    total: number
    total_pages: number
  }
}

class ProductionService {
  // Production Order operations
  async createProductionOrder(data: CreateProductionOrderRequest): Promise<ProductionOrder> {
    const response = await api.post('/production-orders', data)
    return response.data
  }

  async updateProductionOrder(id: string, data: Partial<ProductionOrder>): Promise<ProductionOrder> {
    const response = await api.put(`/production-orders/${id}`, data)
    return response.data
  }

  async getProductionOrder(id: string): Promise<ProductionOrder> {
    const response = await api.get(`/production-orders/${id}`)
    return response.data
  }

  async listProductionOrders(params: ListParams = {}): Promise<PaginatedResponse<ProductionOrder>> {
    const response = await api.get('/production-orders', { params })
    return response.data
  }

  async releaseProductionOrder(id: string): Promise<void> {
    await api.post(`/production-orders/${id}/release`)
  }

  async startProductionOrder(id: string): Promise<void> {
    await api.post(`/production-orders/${id}/start`)
  }

  async completeProductionOrder(id: string): Promise<void> {
    await api.post(`/production-orders/${id}/complete`)
  }

  async cancelProductionOrder(id: string, reason: string): Promise<void> {
    await api.post(`/production-orders/${id}/cancel`, { reason })
  }

  // Production Route operations
  async createProductionRoute(data: CreateProductionRouteRequest): Promise<ProductionRoute> {
    const response = await api.post('/production-routes', data)
    return response.data
  }

  async updateProductionRoute(id: string, data: Partial<ProductionRoute>): Promise<ProductionRoute> {
    const response = await api.put(`/production-routes/${id}`, data)
    return response.data
  }

  async getProductionRoute(id: string): Promise<ProductionRoute> {
    const response = await api.get(`/production-routes/${id}`)
    return response.data
  }

  async listProductionRoutes(params: ListParams = {}): Promise<ProductionRoute[]> {
    const response = await api.get('/production-routes', { params })
    return response.data
  }

  // Route Operation operations
  async createRouteOperation(routeId: string, data: Partial<RouteOperation>): Promise<RouteOperation> {
    const response = await api.post(`/production-routes/${routeId}/operations`, data)
    return response.data
  }

  async updateRouteOperation(id: string, data: Partial<RouteOperation>): Promise<RouteOperation> {
    const response = await api.put(`/route-operations/${id}`, data)
    return response.data
  }

  async getRouteOperations(routeId: string): Promise<RouteOperation[]> {
    const response = await api.get(`/production-routes/${routeId}/operations`)
    return response.data
  }

  async deleteRouteOperation(id: string): Promise<void> {
    await api.delete(`/route-operations/${id}`)
  }

  // Work Station operations
  async createWorkStation(data: CreateWorkStationRequest): Promise<WorkStation> {
    const response = await api.post('/work-stations', data)
    return response.data
  }

  async updateWorkStation(id: string, data: Partial<WorkStation>): Promise<WorkStation> {
    const response = await api.put(`/work-stations/${id}`, data)
    return response.data
  }

  async getWorkStation(id: string): Promise<WorkStation> {
    const response = await api.get(`/work-stations/${id}`)
    return response.data
  }

  async listWorkStations(params: ListParams = {}): Promise<WorkStation[]> {
    const response = await api.get('/work-stations', { params })
    return response.data
  }

  // Production Task operations
  async createProductionTask(data: Partial<ProductionTask>): Promise<ProductionTask> {
    const response = await api.post('/production-tasks', data)
    return response.data
  }

  async updateProductionTask(id: string, data: Partial<ProductionTask>): Promise<ProductionTask> {
    const response = await api.put(`/production-tasks/${id}`, data)
    return response.data
  }

  async getProductionTask(id: string): Promise<ProductionTask> {
    const response = await api.get(`/production-tasks/${id}`)
    return response.data
  }

  async listProductionTasks(params: ListParams = {}): Promise<ProductionTask[]> {
    const response = await api.get('/production-tasks', { params })
    return response.data
  }

  async assignTask(id: string, userId: string): Promise<void> {
    await api.post(`/production-tasks/${id}/assign`, { user_id: userId })
  }

  async startTask(id: string): Promise<void> {
    await api.post(`/production-tasks/${id}/start`)
  }

  async completeTask(id: string, completedQuantity: number, qualifiedQuantity: number, notes: string): Promise<void> {
    await api.post(`/production-tasks/${id}/complete`, {
      completed_quantity: completedQuantity,
      qualified_quantity: qualifiedQuantity,
      notes
    })
  }

  // Production Material operations
  async createProductionMaterial(data: Partial<ProductionMaterial>): Promise<ProductionMaterial> {
    const response = await api.post('/production-materials', data)
    return response.data
  }

  async issueMaterials(productionOrderId: string): Promise<void> {
    await api.post(`/production-orders/${productionOrderId}/issue-materials`)
  }

  async getProductionMaterials(productionOrderId: string): Promise<ProductionMaterial[]> {
    const response = await api.get(`/production-orders/${productionOrderId}/materials`)
    return response.data
  }

  // Quality Inspection operations
  async createQualityInspection(data: CreateQualityInspectionRequest): Promise<QualityInspection> {
    const response = await api.post('/quality-inspections', data)
    return response.data
  }

  async updateQualityInspection(id: string, data: Partial<QualityInspection>): Promise<QualityInspection> {
    const response = await api.put(`/quality-inspections/${id}`, data)
    return response.data
  }

  async getQualityInspection(id: string): Promise<QualityInspection> {
    const response = await api.get(`/quality-inspections/${id}`)
    return response.data
  }

  async listQualityInspections(params: ListParams = {}): Promise<PaginatedResponse<QualityInspection>> {
    const response = await api.get('/quality-inspections', { params })
    return response.data
  }

  async approveInspection(id: string): Promise<void> {
    await api.post(`/quality-inspections/${id}/approve`)
  }

  async rejectInspection(id: string, reason: string): Promise<void> {
    await api.post(`/quality-inspections/${id}/reject`, { reason })
  }

  // Dashboard operations
  async getProductionDashboard(): Promise<ProductionDashboard> {
    const response = await api.get('/production/dashboard')
    return response.data
  }

  async getProductionStats(): Promise<ProductionStats> {
    const response = await api.get('/production/stats')
    return response.data
  }
}

export default new ProductionService()