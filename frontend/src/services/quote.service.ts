import api from './api';
import { 
  Quote, 
  QuoteVersion, 
  QuoteActivityLog, 
  QuoteTermsTemplate,
  QuoteTemplate,
  CreateQuoteRequest,
  UpdateQuoteRequest,
  SubmitApprovalRequest,
  ApproveQuoteRequest,
  SendQuoteRequest
} from '@/types/quote';

export const quoteService = {
  // 基本 CRUD
  async createQuote(data: CreateQuoteRequest) {
    const response = await api.post<Quote>('/quotes', data);
    return response.data;
  },

  async getQuotes(params?: {
    page?: number;
    page_size?: number;
    status?: string;
  }) {
    const response = await api.get<{
      data: Quote[];
      pagination: {
        page: number;
        page_size: number;
        total: number;
        total_page: number;
      };
    }>('/quotes', { params });
    return response.data;
  },

  async getQuote(id: string) {
    const response = await api.get<Quote>(`/quotes/${id}`);
    return response.data;
  },

  async updateQuote(id: string, data: UpdateQuoteRequest) {
    const response = await api.put<Quote>(`/quotes/${id}`, data);
    return response.data;
  },

  // 審核流程
  async submitForApproval(id: string, data: SubmitApprovalRequest = {}) {
    const response = await api.post<{ message: string }>(
      `/quotes/${id}/submit`,
      data
    );
    return response.data;
  },

  async approveQuote(id: string, data: ApproveQuoteRequest) {
    const response = await api.post<{ message: string }>(
      `/quotes/${id}/approve`,
      data
    );
    return response.data;
  },

  // 發送
  async sendQuote(id: string, data: SendQuoteRequest) {
    const response = await api.post<{ message: string }>(
      `/quotes/${id}/send`,
      data
    );
    return response.data;
  },

  // 版本管理
  async getQuoteVersions(id: string) {
    const response = await api.get<QuoteVersion[]>(
      `/quotes/${id}/versions`
    );
    return response.data;
  },

  async getQuoteVersion(id: string, versionId: string) {
    const response = await api.get<QuoteVersion>(
      `/quotes/${id}/versions/${versionId}`
    );
    return response.data;
  },

  // 活動日誌
  async getQuoteActivityLogs(id: string) {
    const response = await api.get<QuoteActivityLog[]>(
      `/quotes/${id}/activities`
    );
    return response.data;
  },

  // 模板
  async getTermsTemplates(type?: string) {
    const response = await api.get<QuoteTermsTemplate[]>(
      '/quotes/terms-templates',
      { params: { type } }
    );
    return response.data;
  },

  async getQuoteTemplates() {
    const response = await api.get<QuoteTemplate[]>(
      '/quotes/templates'
    );
    return response.data;
  },

  // 狀態選項
  async getStatusOptions(): Promise<string[]> {
    return [
      'draft',
      'pending_approval',
      'approved',
      'rejected',
      'sent',
      'expired'
    ];
  }
};

export default quoteService;