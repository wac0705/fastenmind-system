import api from './api'

// Types
export interface Supplier {
  id: string
  company_id: string
  supplier_no: string
  name: string
  name_en: string
  type: 'manufacturer' | 'distributor' | 'service_provider' | 'raw_material'
  status: 'active' | 'inactive' | 'suspended' | 'blacklisted'
  
  // Contact Information
  contact_person: string
  contact_title: string
  phone: string
  mobile: string
  email: string
  website: string
  
  // Address
  country: string
  state: string
  city: string
  address: string
  postal_code: string
  
  // Business Information
  tax_number: string
  business_license: string
  industry: string
  established?: string
  employees: number
  annual_revenue: number
  currency: string
  
  // Payment Terms
  payment_terms: string
  payment_method: string
  credit_limit: number
  credit_days: number
  
  // Performance Metrics
  quality_rating: number
  delivery_rating: number
  service_rating: number
  overall_rating: number
  total_orders: number
  on_time_deliveries: number
  defective_items: number
  
  // Certifications
  iso_9001: boolean
  iso_14001: boolean
  ts_16949: boolean
  ohsas_18001: boolean
  custom_cert: string
  cert_expiry?: string
  
  // Risk Assessment
  risk_level: 'low' | 'medium' | 'high' | 'critical'
  risk_factors: string
  last_audit_date?: string
  next_audit_date?: string
  
  // Financial Information
  credit_rating: string
  financial_health: 'excellent' | 'good' | 'fair' | 'poor'
  insurance_coverage: number
  
  // Additional Information
  description: string
  notes: string
  tags: string
  
  // Timestamps
  created_at: string
  updated_at: string
  created_by: string
  
  // Relations
  company?: any
  creator?: any
  contacts?: SupplierContact[]
  products?: SupplierProduct[]
}

export interface SupplierContact {
  id: string
  supplier_id: string
  name: string
  title: string
  department: string
  phone: string
  mobile: string
  email: string
  is_primary: boolean
  is_active: boolean
  responsibilities: string
  languages: string
  created_at: string
  updated_at: string
  supplier?: Supplier
}

export interface SupplierProduct {
  id: string
  supplier_id: string
  inventory_id?: string
  product_name: string
  product_code: string
  supplier_part_no: string
  category: string
  specification: string
  unit: string
  unit_price: number
  currency: string
  min_order_qty: number
  max_order_qty: number
  price_breaks: string
  lead_time_days: number
  quality_grade: string
  certification_req: boolean
  certificates: string
  status: 'active' | 'inactive' | 'discontinued'
  is_preferred: boolean
  last_purchase_date?: string
  last_purchase_price: number
  total_purchased: number
  quality_issues: number
  delivery_issues: number
  created_at: string
  updated_at: string
  supplier?: Supplier
  inventory?: any
}

export interface PurchaseOrder {
  id: string
  company_id: string
  order_no: string
  status: 'draft' | 'sent' | 'confirmed' | 'partial_received' | 'received' | 'cancelled'
  supplier_id: string
  order_date: string
  required_date: string
  promised_date?: string
  sub_total: number
  tax_rate: number
  tax_amount: number
  shipping_cost: number
  total_amount: number
  currency: string
  exchange_rate: number
  payment_terms: string
  payment_method: string
  shipping_address: string
  shipping_method: string
  tracking_number: string
  notes: string
  internal_notes: string
  created_at: string
  updated_at: string
  created_by: string
  approved_by?: string
  approved_at?: string
  
  // Relations
  company?: any
  supplier?: Supplier
  creator?: any
  approver?: any
  items?: PurchaseOrderItem[]
}

export interface PurchaseOrderItem {
  id: string
  purchase_order_id: string
  supplier_product_id?: string
  inventory_id?: string
  product_name: string
  product_code: string
  supplier_part_no: string
  specification: string
  ordered_quantity: number
  received_quantity: number
  unit: string
  unit_price: number
  total_price: number
  status: 'pending' | 'partial_received' | 'received' | 'cancelled'
  quality_requirement: string
  inspection_required: boolean
  created_at: string
  updated_at: string
  
  // Relations
  purchase_order?: PurchaseOrder
  supplier_product?: SupplierProduct
  inventory?: any
}

export interface SupplierEvaluation {
  id: string
  company_id: string
  supplier_id: string
  evaluation_no: string
  start_date: string
  end_date: string
  evaluation_type: 'monthly' | 'quarterly' | 'annual' | 'ad_hoc'
  quality_score: number
  delivery_score: number
  service_score: number
  cost_score: number
  technical_score: number
  overall_score: number
  total_orders: number
  on_time_deliveries: number
  quality_defects: number
  service_issues: number
  cost_savings: number
  strengths: string
  weaknesses: string
  recommendations: string
  action_items: string
  status: 'draft' | 'completed' | 'approved'
  evaluated_by: string
  evaluated_at: string
  approved_by?: string
  approved_at?: string
  created_at: string
  updated_at: string
  
  // Relations
  company?: any
  supplier?: Supplier
  evaluator?: any
  approver?: any
}

export interface SupplierDashboard {
  total_suppliers: number
  active_suppliers: number
  suspended_suppliers: number
  blacklisted_suppliers: number
  total_purchase_orders: number
  draft_orders: number
  sent_orders: number
  confirmed_orders: number
  received_orders: number
  cancelled_orders: number
  total_purchase_value: number
  pending_value: number
  received_value: number
  high_risk_suppliers: number
  critical_risk_suppliers: number
  average_quality_rating: number
  average_delivery_rating: number
  average_service_rating: number
  total_evaluations: number
  pending_evaluations: number
  completed_evaluations: number
}

export interface CreateSupplierRequest {
  name: string
  name_en?: string
  type: string
  contact_person?: string
  contact_title?: string
  phone?: string
  mobile?: string
  email?: string
  website?: string
  country?: string
  state?: string
  city?: string
  address?: string
  postal_code?: string
  tax_number?: string
  business_license?: string
  industry?: string
  established?: string
  employees?: number
  annual_revenue?: number
  currency?: string
  payment_terms?: string
  payment_method?: string
  credit_limit?: number
  credit_days?: number
  iso_9001?: boolean
  iso_14001?: boolean
  ts_16949?: boolean
  ohsas_18001?: boolean
  custom_cert?: string
  cert_expiry?: string
  credit_rating?: string
  financial_health?: string
  insurance_coverage?: number
  description?: string
  notes?: string
  tags?: string
}

export interface UpdateSupplierRequest {
  name?: string
  name_en?: string
  type?: string
  status?: string
  contact_person?: string
  contact_title?: string
  phone?: string
  mobile?: string
  email?: string
  website?: string
  country?: string
  state?: string
  city?: string
  address?: string
  postal_code?: string
  tax_number?: string
  business_license?: string
  industry?: string
  established?: string
  employees?: number
  annual_revenue?: number
  currency?: string
  payment_terms?: string
  payment_method?: string
  credit_limit?: number
  credit_days?: number
  quality_rating?: number
  delivery_rating?: number
  service_rating?: number
  overall_rating?: number
  iso_9001?: boolean
  iso_14001?: boolean
  ts_16949?: boolean
  ohsas_18001?: boolean
  custom_cert?: string
  cert_expiry?: string
  risk_level?: string
  risk_factors?: string
  last_audit_date?: string
  next_audit_date?: string
  credit_rating?: string
  financial_health?: string
  insurance_coverage?: number
  description?: string
  notes?: string
  tags?: string
}

export interface CreateSupplierContactRequest {
  name: string
  title?: string
  department?: string
  phone?: string
  mobile?: string
  email?: string
  is_primary?: boolean
  responsibilities?: string
  languages?: string
}

export interface UpdateSupplierContactRequest {
  name?: string
  title?: string
  department?: string
  phone?: string
  mobile?: string
  email?: string
  is_primary?: boolean
  is_active?: boolean
  responsibilities?: string
  languages?: string
}

export interface CreateSupplierProductRequest {
  inventory_id?: string
  product_name: string
  product_code?: string
  supplier_part_no?: string
  category?: string
  specification?: string
  unit: string
  unit_price?: number
  currency?: string
  min_order_qty?: number
  max_order_qty?: number
  price_breaks?: string
  lead_time_days?: number
  quality_grade?: string
  certification_req?: boolean
  certificates?: string
  is_preferred?: boolean
}

export interface UpdateSupplierProductRequest {
  inventory_id?: string
  product_name?: string
  product_code?: string
  supplier_part_no?: string
  category?: string
  specification?: string
  unit?: string
  unit_price?: number
  currency?: string
  min_order_qty?: number
  max_order_qty?: number
  price_breaks?: string
  lead_time_days?: number
  quality_grade?: string
  certification_req?: boolean
  certificates?: string
  status?: string
  is_preferred?: boolean
  last_purchase_date?: string
  last_purchase_price?: number
  total_purchased?: number
  quality_issues?: number
  delivery_issues?: number
}

export interface CreatePurchaseOrderRequest {
  supplier_id: string
  order_date: string
  required_date: string
  payment_terms?: string
  payment_method?: string
  shipping_address?: string
  shipping_method?: string
  currency?: string
  exchange_rate?: number
  tax_rate?: number
  shipping_cost?: number
  notes?: string
  internal_notes?: string
  items: CreatePurchaseOrderItemRequest[]
}

export interface UpdatePurchaseOrderRequest {
  status?: string
  required_date?: string
  promised_date?: string
  payment_terms?: string
  payment_method?: string
  shipping_address?: string
  shipping_method?: string
  tracking_number?: string
  currency?: string
  exchange_rate?: number
  tax_rate?: number
  shipping_cost?: number
  notes?: string
  internal_notes?: string
}

export interface CreatePurchaseOrderItemRequest {
  supplier_product_id?: string
  inventory_id?: string
  product_name: string
  product_code?: string
  supplier_part_no?: string
  specification?: string
  ordered_quantity: number
  unit: string
  unit_price: number
  quality_requirement?: string
  inspection_required?: boolean
}

export interface UpdatePurchaseOrderItemRequest {
  supplier_product_id?: string
  inventory_id?: string
  product_name?: string
  product_code?: string
  supplier_part_no?: string
  specification?: string
  ordered_quantity?: number
  received_quantity?: number
  unit?: string
  unit_price?: number
  status?: string
  quality_requirement?: string
  inspection_required?: boolean
}

export interface PurchaseOrderReceiptItem {
  item_id: string
  received_quantity: number
  quality_passed?: boolean
  inspection_notes?: string
}

export interface CreateSupplierEvaluationRequest {
  supplier_id: string
  start_date: string
  end_date: string
  evaluation_type: string
  quality_score: number
  delivery_score: number
  service_score: number
  cost_score: number
  technical_score: number
  total_orders?: number
  on_time_deliveries?: number
  quality_defects?: number
  service_issues?: number
  cost_savings?: number
  strengths?: string
  weaknesses?: string
  recommendations?: string
  action_items?: string
  evaluated_at: string
}

export interface UpdateSupplierEvaluationRequest {
  start_date?: string
  end_date?: string
  evaluation_type?: string
  quality_score?: number
  delivery_score?: number
  service_score?: number
  cost_score?: number
  technical_score?: number
  total_orders?: number
  on_time_deliveries?: number
  quality_defects?: number
  service_issues?: number
  cost_savings?: number
  strengths?: string
  weaknesses?: string
  recommendations?: string
  action_items?: string
  status?: string
}

export interface ListParams {
  page?: number
  page_size?: number
  search?: string
  type?: string
  status?: string
  country?: string
  risk_level?: string
  sort_by?: string
  sort_order?: string
  supplier_id?: string
  start_date?: string
  end_date?: string
  evaluation_type?: string
  evaluated_by?: string
  category?: string
  is_preferred?: boolean
  inspector_id?: string
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
}

class SupplierService {
  // Supplier operations
  async createSupplier(data: CreateSupplierRequest): Promise<Supplier> {
    const response = await api.post('/suppliers', data)
    return response.data
  }

  async updateSupplier(id: string, data: UpdateSupplierRequest): Promise<Supplier> {
    const response = await api.put(`/suppliers/${id}`, data)
    return response.data
  }

  async getSupplier(id: string): Promise<Supplier> {
    const response = await api.get(`/suppliers/${id}`)
    return response.data
  }

  async listSuppliers(params: ListParams = {}): Promise<PaginatedResponse<Supplier>> {
    const response = await api.get('/suppliers', { params })
    return response.data
  }

  // Supplier Contact operations
  async addSupplierContact(supplierId: string, data: CreateSupplierContactRequest): Promise<SupplierContact> {
    const response = await api.post(`/suppliers/${supplierId}/contacts`, data)
    return response.data
  }

  async updateSupplierContact(id: string, data: UpdateSupplierContactRequest): Promise<SupplierContact> {
    const response = await api.put(`/supplier-contacts/${id}`, data)
    return response.data
  }

  async getSupplierContacts(supplierId: string): Promise<SupplierContact[]> {
    const response = await api.get(`/suppliers/${supplierId}/contacts`)
    return response.data
  }

  async deleteSupplierContact(id: string): Promise<void> {
    await api.delete(`/supplier-contacts/${id}`)
  }

  // Supplier Product operations
  async addSupplierProduct(supplierId: string, data: CreateSupplierProductRequest): Promise<SupplierProduct> {
    const response = await api.post(`/suppliers/${supplierId}/products`, data)
    return response.data
  }

  async updateSupplierProduct(id: string, data: UpdateSupplierProductRequest): Promise<SupplierProduct> {
    const response = await api.put(`/supplier-products/${id}`, data)
    return response.data
  }

  async getSupplierProducts(supplierId: string, params: ListParams = {}): Promise<SupplierProduct[]> {
    const response = await api.get(`/suppliers/${supplierId}/products`, { params })
    return response.data
  }

  async deleteSupplierProduct(id: string): Promise<void> {
    await api.delete(`/supplier-products/${id}`)
  }

  // Purchase Order operations
  async createPurchaseOrder(data: CreatePurchaseOrderRequest): Promise<PurchaseOrder> {
    const response = await api.post('/purchase-orders', data)
    return response.data
  }

  async updatePurchaseOrder(id: string, data: UpdatePurchaseOrderRequest): Promise<PurchaseOrder> {
    const response = await api.put(`/purchase-orders/${id}`, data)
    return response.data
  }

  async getPurchaseOrder(id: string): Promise<PurchaseOrder> {
    const response = await api.get(`/purchase-orders/${id}`)
    return response.data
  }

  async listPurchaseOrders(params: ListParams = {}): Promise<PaginatedResponse<PurchaseOrder>> {
    const response = await api.get('/purchase-orders', { params })
    return response.data
  }

  async approvePurchaseOrder(id: string): Promise<void> {
    await api.post(`/purchase-orders/${id}/approve`)
  }

  async sendPurchaseOrder(id: string): Promise<void> {
    await api.post(`/purchase-orders/${id}/send`)
  }

  async receivePurchaseOrder(id: string, items: PurchaseOrderReceiptItem[]): Promise<void> {
    await api.post(`/purchase-orders/${id}/receive`, { items })
  }

  // Purchase Order Item operations
  async addPurchaseOrderItem(purchaseOrderId: string, data: CreatePurchaseOrderItemRequest): Promise<PurchaseOrderItem> {
    const response = await api.post(`/purchase-orders/${purchaseOrderId}/items`, data)
    return response.data
  }

  async updatePurchaseOrderItem(id: string, data: UpdatePurchaseOrderItemRequest): Promise<PurchaseOrderItem> {
    const response = await api.put(`/purchase-order-items/${id}`, data)
    return response.data
  }

  async getPurchaseOrderItems(purchaseOrderId: string): Promise<PurchaseOrderItem[]> {
    const response = await api.get(`/purchase-orders/${purchaseOrderId}/items`)
    return response.data
  }

  async deletePurchaseOrderItem(id: string): Promise<void> {
    await api.delete(`/purchase-order-items/${id}`)
  }

  // Supplier Evaluation operations
  async createSupplierEvaluation(data: CreateSupplierEvaluationRequest): Promise<SupplierEvaluation> {
    const response = await api.post('/supplier-evaluations', data)
    return response.data
  }

  async updateSupplierEvaluation(id: string, data: UpdateSupplierEvaluationRequest): Promise<SupplierEvaluation> {
    const response = await api.put(`/supplier-evaluations/${id}`, data)
    return response.data
  }

  async getSupplierEvaluation(id: string): Promise<SupplierEvaluation> {
    const response = await api.get(`/supplier-evaluations/${id}`)
    return response.data
  }

  async listSupplierEvaluations(params: ListParams = {}): Promise<PaginatedResponse<SupplierEvaluation>> {
    const response = await api.get('/supplier-evaluations', { params })
    return response.data
  }

  async approveSupplierEvaluation(id: string): Promise<void> {
    await api.post(`/supplier-evaluations/${id}/approve`)
  }

  // Business operations
  async updateSupplierPerformance(supplierId: string): Promise<void> {
    await api.post(`/suppliers/${supplierId}/update-performance`)
  }

  async calculateSupplierRisk(supplierId: string): Promise<{ supplier_id: string; risk_level: string }> {
    const response = await api.post(`/suppliers/${supplierId}/calculate-risk`)
    return response.data
  }

  async getSupplierDashboard(): Promise<SupplierDashboard> {
    const response = await api.get('/suppliers/dashboard')
    return response.data
  }

  // Utility functions
  async exportSuppliers(params: ListParams = {}): Promise<Blob> {
    const response = await api.get('/suppliers/export', { 
      params,
      responseType: 'blob'
    })
    return response.data
  }

  async importSuppliers(file: File): Promise<{ success: number; errors: string[] }> {
    const formData = new FormData()
    formData.append('file', file)
    
    const response = await api.post('/suppliers/import', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
    
    return response.data
  }

  async validateSupplierData(data: CreateSupplierRequest): Promise<{ valid: boolean; errors: string[] }> {
    const response = await api.post('/suppliers/validate', data)
    return response.data
  }

  async getSupplierTypes(): Promise<string[]> {
    const response = await api.get('/suppliers/types')
    return response.data
  }

  async getSupplierStatuses(): Promise<string[]> {
    const response = await api.get('/suppliers/statuses')
    return response.data
  }

  async getRiskLevels(): Promise<string[]> {
    const response = await api.get('/suppliers/risk-levels')
    return response.data
  }

  async getEvaluationTypes(): Promise<string[]> {
    const response = await api.get('/supplier-evaluations/types')
    return response.data
  }

  async getPurchaseOrderStatuses(): Promise<string[]> {
    const response = await api.get('/purchase-orders/statuses')
    return response.data
  }
}

export default new SupplierService()