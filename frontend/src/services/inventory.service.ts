import api from '@/lib/axios'
import { ApiResponse, PaginatedResponse } from '@/types/api'

export interface Inventory {
  id: string
  company_id: string
  sku: string
  part_no: string
  name: string
  description?: string
  category: string
  material?: string
  specification?: string
  surface_treatment?: string
  heat_treatment?: string
  unit: string
  current_stock: number
  available_stock: number
  reserved_stock: number
  min_stock: number
  max_stock: number
  reorder_point: number
  reorder_quantity: number
  warehouse_id?: string
  location?: string
  last_purchase_price: number
  average_cost: number
  standard_cost: number
  currency: string
  primary_supplier_id?: string
  lead_time_days: number
  status: string
  is_active: boolean
  created_at: string
  updated_at: string
  last_stock_check_at?: string
  
  // Relations
  warehouse?: Warehouse
  primary_supplier?: any
}

export interface Warehouse {
  id: string
  company_id: string
  code: string
  name: string
  type: string
  address?: string
  manager?: string
  phone?: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface StockMovement {
  id: string
  company_id: string
  inventory_id: string
  movement_type: string
  reason: string
  quantity: number
  unit_cost: number
  total_cost: number
  reference_type?: string
  reference_id?: string
  reference_no?: string
  from_warehouse_id?: string
  to_warehouse_id?: string
  from_location?: string
  to_location?: string
  before_quantity: number
  after_quantity: number
  batch_no?: string
  serial_no?: string
  expiry_date?: string
  notes?: string
  created_by: string
  approved_by?: string
  approved_at?: string
  created_at: string
  
  // Relations
  inventory?: Inventory
  from_warehouse?: Warehouse
  to_warehouse?: Warehouse
  creator?: any
  approver?: any
}

export interface StockAlert {
  id: string
  company_id: string
  inventory_id: string
  alert_type: string
  status: string
  priority: string
  current_level: number
  threshold_level: number
  message: string
  acknowledged_by?: string
  acknowledged_at?: string
  resolved_by?: string
  resolved_at?: string
  resolution?: string
  created_at: string
  updated_at: string
  
  // Relations
  inventory?: Inventory
}

export interface StockTake {
  id: string
  company_id: string
  reference_no: string
  warehouse_id: string
  status: string
  type: string
  scheduled_date: string
  started_at?: string
  completed_at?: string
  created_by: string
  assigned_to: string
  reviewed_by?: string
  total_items: number
  counted_items: number
  variance_items: number
  total_variance: number
  notes?: string
  created_at: string
  updated_at: string
  
  // Relations
  warehouse?: Warehouse
  creator?: any
  assignee?: any
  reviewer?: any
}

export interface CreateInventoryRequest {
  sku: string
  part_no: string
  name: string
  description?: string
  category: string
  material?: string
  specification?: string
  surface_treatment?: string
  heat_treatment?: string
  unit?: string
  initial_stock?: number
  min_stock?: number
  max_stock?: number
  reorder_point?: number
  reorder_quantity?: number
  warehouse_id: string
  location?: string
  standard_cost?: number
  primary_supplier_id?: string
  lead_time_days?: number
}

export interface UpdateInventoryRequest {
  name?: string
  description?: string
  category?: string
  material?: string
  specification?: string
  surface_treatment?: string
  heat_treatment?: string
  min_stock?: number
  max_stock?: number
  reorder_point?: number
  reorder_quantity?: number
  location?: string
  standard_cost?: number
  lead_time_days?: number
  status?: string
}

export interface StockAdjustmentRequest {
  quantity: number
  reason: string
  notes?: string
  batch_no?: string
  warehouse_id?: string
}

export interface StockTransferRequest {
  inventory_id: string
  quantity: number
  from_warehouse_id: string
  to_warehouse_id: string
  from_location?: string
  to_location?: string
  notes?: string
}

export interface CreateWarehouseRequest {
  code: string
  name: string
  type: string
  address?: string
  manager?: string
  phone?: string
}

export interface InventoryStats {
  total_items: number
  total_value: number
  low_stock_items: number
  out_of_stock_items: number
  overstock_items: number
  active_alerts: number
}

export interface StockValuation {
  total_value: number
  by_category: Record<string, number>
  by_warehouse: Record<string, number>
  top_value_items: {
    inventory_id: string
    sku: string
    name: string
    quantity: number
    unit_cost: number
    total_value: number
  }[]
}

class InventoryService {
  // Inventory management
  async list(params?: {
    page?: number
    page_size?: number
    category?: string
    warehouse_id?: string
    status?: string
    low_stock?: boolean
    search?: string
  }): Promise<PaginatedResponse<Inventory>> {
    const { data } = await api.get('/api/inventory', { params })
    return data
  }

  async get(id: string): Promise<Inventory> {
    const { data } = await api.get(`/api/inventory/${id}`)
    return data
  }

  async getBySKU(sku: string): Promise<Inventory> {
    const { data } = await api.get(`/api/inventory/sku/${sku}`)
    return data
  }

  async create(data: CreateInventoryRequest): Promise<Inventory> {
    const response = await api.post('/api/inventory', data)
    return response.data
  }

  async update(id: string, data: UpdateInventoryRequest): Promise<Inventory> {
    const response = await api.put(`/api/inventory/${id}`, data)
    return response.data
  }

  async delete(id: string): Promise<void> {
    await api.delete(`/api/inventory/${id}`)
  }

  // Stock operations
  async adjustStock(id: string, data: StockAdjustmentRequest): Promise<StockMovement> {
    const response = await api.post(`/api/inventory/${id}/adjust`, data)
    return response.data
  }

  async transferStock(data: StockTransferRequest): Promise<StockMovement> {
    const response = await api.post('/api/inventory/transfer', data)
    return response.data
  }

  async getMovements(inventoryId: string, params?: {
    movement_type?: string
    start_date?: string
    end_date?: string
  }): Promise<StockMovement[]> {
    const { data } = await api.get(`/api/inventory/${inventoryId}/movements`, { params })
    return data
  }

  // Reports
  async getStats(): Promise<InventoryStats> {
    const { data } = await api.get('/api/inventory/stats')
    return data
  }

  async getLowStockItems(): Promise<Inventory[]> {
    const { data } = await api.get('/api/inventory/low-stock')
    return data
  }

  async getValuation(): Promise<StockValuation> {
    const { data } = await api.get('/api/inventory/valuation')
    return data
  }

  // Alerts
  async getAlerts(): Promise<StockAlert[]> {
    const { data } = await api.get('/api/inventory/alerts')
    return data
  }

  // Warehouses
  async listWarehouses(): Promise<Warehouse[]> {
    const { data } = await api.get('/api/warehouses')
    return data
  }

  async createWarehouse(data: CreateWarehouseRequest): Promise<Warehouse> {
    const response = await api.post('/api/warehouses', data)
    return response.data
  }

  // Stock takes
  async listStockTakes(params?: {
    status?: string
    warehouse_id?: string
  }): Promise<StockTake[]> {
    const { data } = await api.get('/api/stock-takes', { params })
    return data
  }

  async createStockTake(data: {
    warehouse_id: string
    type: string
    scheduled_date: string
    assigned_to: string
    notes?: string
  }): Promise<StockTake> {
    const response = await api.post('/api/stock-takes', data)
    return response.data
  }
}

export const inventoryService = new InventoryService()
export default inventoryService