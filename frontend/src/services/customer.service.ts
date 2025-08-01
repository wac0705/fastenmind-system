import api from './api'

export interface Customer {
  id: string
  company_id: string
  customer_code: string
  name: string
  name_en?: string
  short_name?: string
  country: string
  tax_id?: string
  address?: string
  shipping_address?: string
  contact_person?: string
  contact_phone?: string
  contact_email?: string
  payment_terms?: string
  credit_limit?: number
  credit_used?: number
  total_revenue?: number
  currency: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface CreateCustomerRequest {
  customer_code: string
  name: string
  name_en?: string
  short_name?: string
  country: string
  tax_id?: string
  address?: string
  shipping_address?: string
  contact_person?: string
  contact_phone?: string
  contact_email?: string
  payment_terms?: string
  credit_limit?: number
  currency?: string
  is_active?: boolean
}

export interface CustomerListParams {
  page?: number
  page_size?: number
  search?: string
  country?: string
  is_active?: boolean
}

export interface CustomerTransaction {
  id: string
  customer_id: string
  type: 'inquiry' | 'quote' | 'order' | 'payment'
  reference_no: string
  amount: number
  currency: string
  date: string
  status: string
  description: string
}

export interface CustomerStats {
  total_inquiries: number
  total_quotes: number
  total_orders: number
  total_revenue: number
  outstanding_amount: number
  average_order_value: number
  last_transaction_date?: string
}

class CustomerService {
  async list(params?: CustomerListParams): Promise<{ data: Customer[]; pagination: any }> {
    const response = await api.get('/customers', { params })
    return response.data
  }

  async get(id: string): Promise<Customer> {
    const response = await api.get(`/customers/${id}`)
    return response.data
  }

  async create(data: CreateCustomerRequest): Promise<Customer> {
    const response = await api.post('/customers', data)
    return response.data
  }

  async update(id: string, data: Partial<CreateCustomerRequest>): Promise<Customer> {
    const response = await api.put(`/customers/${id}`, data)
    return response.data
  }

  async delete(id: string): Promise<void> {
    await api.delete(`/customers/${id}`)
  }

  async getStats(id: string): Promise<CustomerStats> {
    const response = await api.get(`/customers/${id}/stats`)
    return response.data
  }

  async getTransactions(id: string, params?: {
    page?: number
    page_size?: number
    type?: string
    from_date?: string
    to_date?: string
  }): Promise<{ data: CustomerTransaction[]; pagination: any }> {
    const response = await api.get(`/customers/${id}/transactions`, { params })
    return response.data
  }

  async checkCreditLimit(id: string, amount: number): Promise<{
    available_credit: number
    is_within_limit: boolean
    message?: string
  }> {
    const response = await api.post(`/customers/${id}/check-credit`, { amount })
    return response.data
  }

  async exportCustomers(params?: CustomerListParams): Promise<Blob> {
    const response = await api.get('/customers/export', {
      params,
      responseType: 'blob'
    })
    return response.data
  }

  async importCustomers(file: File): Promise<{
    success: number
    failed: number
    errors: Array<{ row: number; error: string }>
  }> {
    const formData = new FormData()
    formData.append('file', file)
    
    const response = await api.post('/customers/import', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })
    
    return response.data
  }

  async getCountryOptions(): Promise<Array<{ code: string; name: string }>> {
    return [
      { code: 'TW', name: '台灣 Taiwan' },
      { code: 'CN', name: '中國 China' },
      { code: 'US', name: '美國 United States' },
      { code: 'DE', name: '德國 Germany' },
      { code: 'JP', name: '日本 Japan' },
      { code: 'KR', name: '韓國 South Korea' },
      { code: 'VN', name: '越南 Vietnam' },
      { code: 'TH', name: '泰國 Thailand' },
      { code: 'MY', name: '馬來西亞 Malaysia' },
      { code: 'SG', name: '新加坡 Singapore' },
    ]
  }

  async getCurrencyOptions(): Promise<string[]> {
    return ['USD', 'EUR', 'TWD', 'CNY', 'JPY', 'KRW', 'SGD', 'MYR', 'THB', 'VND']
  }

  // Get customer related data for detail page
  async getInquiries(customerId: string): Promise<any[]> {
    const response = await api.get(`/customers/${customerId}/inquiries`)
    return response.data.data || []
  }

  async getQuotes(customerId: string): Promise<any[]> {
    const response = await api.get(`/customers/${customerId}/quotes`)
    return response.data.data || []
  }

  async getOrders(customerId: string): Promise<any[]> {
    const response = await api.get(`/customers/${customerId}/orders`)
    return response.data.data || []
  }

  async getCreditHistory(customerId: string): Promise<any[]> {
    const response = await api.get(`/customers/${customerId}/credit-history`)
    return response.data.data || []
  }
}

export const customerService = new CustomerService()
export default customerService