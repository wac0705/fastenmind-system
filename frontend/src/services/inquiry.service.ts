import api from './api'
import { AxiosResponse } from 'axios'
import n8nService from './n8n.service'

export interface Inquiry {
  id: string
  inquiry_no: string
  company_id: string
  customer_id: string
  sales_id: string
  status: string
  product_category: string
  product_name: string
  drawing_files: string[]
  quantity: number
  unit: string
  required_date: string
  incoterm: string
  destination_port?: string
  destination_address?: string
  payment_terms?: string
  special_requirements?: string
  assigned_engineer_id?: string
  assigned_at?: string
  quote_id?: string
  quoted_at?: string
  created_at: string
  updated_at: string
  
  // Relations
  customer?: Customer
  sales?: User
  assigned_engineer?: User
}

export interface Customer {
  id: string
  customer_code: string
  name: string
  name_en?: string
  country: string
  contact_person?: string
  contact_email?: string
  contact_phone?: string
}

export interface User {
  id: string
  full_name: string
  email: string
  role: string
}

export interface CreateInquiryRequest {
  customer_id: string
  product_category: string
  product_name: string
  drawing_files: string[]
  quantity: number
  unit: string
  required_date: string
  incoterm: string
  destination_port?: string
  destination_address?: string
  payment_terms?: string
  special_requirements?: string
}

export interface InquiryListParams {
  page?: number
  page_size?: number
  status?: string
  assigned_engineer_id?: string
  customer_id?: string
  search?: string
}

export interface PaginatedResponse<T> {
  data: T[]
  pagination: {
    page: number
    page_size: number
    total: number
  }
}

class InquiryService {
  async list(params?: InquiryListParams): Promise<PaginatedResponse<Inquiry>> {
    const response = await api.get<PaginatedResponse<Inquiry>>('/inquiries', { params })
    return response.data
  }

  async get(id: string): Promise<Inquiry> {
    const response = await api.get<Inquiry>(`/inquiries/${id}`)
    return response.data
  }

  async create(data: CreateInquiryRequest): Promise<Inquiry> {
    const response = await api.post<Inquiry>('/inquiries', data)
    
    // Trigger N8N workflow for new inquiry
    try {
      await n8nService.notifyInquiryCreated(response.data.id)
    } catch (error) {
      console.error('Failed to trigger N8N workflow:', error)
    }
    
    return response.data
  }

  async update(id: string, data: Partial<CreateInquiryRequest>): Promise<Inquiry> {
    const response = await api.put<Inquiry>(`/inquiries/${id}`, data)
    return response.data
  }

  async delete(id: string): Promise<void> {
    await api.delete(`/inquiries/${id}`)
  }

  async assignEngineer(id: string, engineerId: string, notes?: string): Promise<Inquiry> {
    const response = await api.post<Inquiry>(`/inquiries/${id}/assign`, {
      engineer_id: engineerId,
      notes,
    })
    return response.data
  }

  async createQuote(id: string, quoteData: any): Promise<any> {
    const response = await api.post(`/inquiries/${id}/quote`, quoteData)
    return response.data
  }

  async uploadDrawing(file: File): Promise<{ url: string }> {
    const formData = new FormData()
    formData.append('file', file)
    
    const response = await api.post<{ url: string }>('/files/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })
    
    return response.data
  }

  async getStatusOptions(): Promise<string[]> {
    return [
      'draft',
      'pending',
      'assigned',
      'in_progress',
      'under_review',
      'approved',
      'quoted',
      'rejected',
      'cancelled',
    ]
  }

  async getIncotermOptions(): Promise<string[]> {
    return ['EXW', 'FCA', 'FOB', 'CFR', 'CIF', 'CPT', 'CIP', 'DAP', 'DPU', 'DDP']
  }
}

export const inquiryService = new InquiryService()
export default inquiryService