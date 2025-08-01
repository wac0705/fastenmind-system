import api from '@/lib/axios'
import { ApiResponse, PaginatedResponse } from '@/types/api'

export interface Order {
  id: string
  order_no: string
  quote_id: string
  company_id: string
  customer_id: string
  sales_id: string
  status: string
  po_number: string
  quantity: number
  unit_price: number
  total_amount: number
  currency: string
  delivery_method: string
  delivery_date: string
  shipping_address: string
  payment_terms: string
  payment_status: string
  down_payment: number
  paid_amount: number
  notes?: string
  internal_notes?: string
  confirmed_at?: string
  in_production_at?: string
  quality_check_at?: string
  ready_to_ship_at?: string
  shipped_at?: string
  delivered_at?: string
  completed_at?: string
  cancelled_at?: string
  created_at: string
  updated_at: string
  
  // Relations
  quote?: any
  customer?: any
  sales?: any
}

export interface OrderItem {
  id: string
  order_id: string
  part_no: string
  description?: string
  quantity: number
  unit_price: number
  total_price: number
  material?: string
  surface_treatment?: string
  heat_treatment?: string
  specifications?: string
}

export interface OrderDocument {
  id: string
  order_id: string
  document_type: string
  file_name: string
  file_path: string
  file_size: number
  uploaded_by: string
  created_at: string
  uploader?: any
}

export interface OrderActivity {
  id: string
  order_id: string
  user_id: string
  action: string
  description: string
  metadata?: string
  created_at: string
  user?: any
}

export interface CreateOrderRequest {
  quote_id: string
  po_number: string
  quantity: number
  delivery_method: string
  delivery_date: string
  shipping_address: string
  payment_terms: string
  down_payment?: number
  notes?: string
}

export interface UpdateOrderRequest {
  po_number?: string
  delivery_method?: string
  delivery_date?: string
  shipping_address?: string
  payment_terms?: string
  notes?: string
  internal_notes?: string
}

export interface OrderItemRequest {
  part_no: string
  description?: string
  quantity: number
  unit_price: number
  material?: string
  surface_treatment?: string
  heat_treatment?: string
  specifications?: string
}

export interface OrderStats {
  total_orders: number
  pending_orders: number
  in_production: number
  completed_orders: number
  total_revenue: number
  avg_order_value: number
}

class OrderService {
  async list(params?: {
    page?: number
    page_size?: number
    status?: string
    customer_id?: string
    sales_id?: string
    payment_status?: string
    search?: string
    start_date?: string
    end_date?: string
  }): Promise<PaginatedResponse<Order>> {
    const { data } = await api.get('/api/orders', { params })
    return data
  }

  async get(id: string): Promise<Order> {
    const { data } = await api.get(`/api/orders/${id}`)
    return data
  }

  async createFromQuote(data: CreateOrderRequest): Promise<Order> {
    const response = await api.post('/api/orders', data)
    return response.data
  }

  async update(id: string, data: UpdateOrderRequest): Promise<Order> {
    const response = await api.put(`/api/orders/${id}`, data)
    return response.data
  }

  async updateStatus(id: string, status: string, notes?: string): Promise<Order> {
    const response = await api.put(`/api/orders/${id}/status`, { status, notes })
    return response.data
  }

  async delete(id: string): Promise<void> {
    await api.delete(`/api/orders/${id}`)
  }

  // Order items
  async getItems(orderId: string): Promise<OrderItem[]> {
    const { data } = await api.get(`/api/orders/${orderId}/items`)
    return data
  }

  async updateItems(orderId: string, items: OrderItemRequest[]): Promise<void> {
    await api.put(`/api/orders/${orderId}/items`, items)
  }

  // Documents
  async getDocuments(orderId: string): Promise<OrderDocument[]> {
    const { data } = await api.get(`/api/orders/${orderId}/documents`)
    return data
  }

  async addDocument(orderId: string, document: {
    document_type: string
    file_name: string
    file_path: string
    file_size: number
  }): Promise<OrderDocument> {
    const { data } = await api.post(`/api/orders/${orderId}/documents`, document)
    return data
  }

  async removeDocument(orderId: string, docId: string): Promise<void> {
    await api.delete(`/api/orders/${orderId}/documents/${docId}`)
  }

  // Activities
  async getActivities(orderId: string): Promise<OrderActivity[]> {
    const { data } = await api.get(`/api/orders/${orderId}/activities`)
    return data
  }

  // Stats
  async getStats(): Promise<OrderStats> {
    const { data } = await api.get('/api/orders/stats')
    return data
  }
}

export const orderService = new OrderService()
export default orderService