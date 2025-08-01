import api from './api'

// Types
export interface Invoice {
  id: string
  company_id: string
  invoice_no: string
  type: 'sales' | 'purchase' | 'credit_note' | 'debit_note'
  status: 'draft' | 'issued' | 'sent' | 'paid' | 'partial_paid' | 'overdue' | 'cancelled'
  order_id?: string
  customer_id?: string
  supplier_id?: string
  issue_date: string
  due_date: string
  payment_date?: string
  sub_total: number
  tax_rate: number
  tax_amount: number
  discount_rate: number
  discount_amount: number
  total_amount: number
  paid_amount: number
  balance_amount: number
  currency: string
  exchange_rate: number
  payment_terms: string
  payment_method: string
  bank_account: string
  notes: string
  internal_notes: string
  created_at: string
  updated_at: string
  created_by: string
  
  // Relations
  order?: any
  customer?: any
  supplier?: any
  creator?: any
}

export interface InvoiceItem {
  id: string
  invoice_id: string
  description: string
  quantity: number
  unit: string
  unit_price: number
  total_price: number
  tax_rate: number
  tax_amount: number
  order_item_id?: string
  inventory_id?: string
  created_at: string
  updated_at: string
  
  // Relations
  order_item?: any
  inventory?: any
}

export interface Payment {
  id: string
  company_id: string
  payment_no: string
  type: 'incoming' | 'outgoing'
  status: 'pending' | 'completed' | 'failed' | 'cancelled'
  invoice_id?: string
  customer_id?: string
  supplier_id?: string
  payment_date: string
  amount: number
  currency: string
  exchange_rate: number
  payment_method: 'cash' | 'check' | 'bank_transfer' | 'credit_card'
  bank_name: string
  bank_account: string
  transaction_no: string
  check_no: string
  notes: string
  attachment_path: string
  created_at: string
  updated_at: string
  created_by: string
  
  // Relations
  invoice?: Invoice
  customer?: any
  supplier?: any
  creator?: any
}

export interface Expense {
  id: string
  company_id: string
  expense_no: string
  category: string
  sub_category: string
  status: 'draft' | 'submitted' | 'approved' | 'paid' | 'rejected'
  supplier_id?: string
  department_id?: string
  project_id?: string
  expense_date: string
  description: string
  amount: number
  tax_amount: number
  total_amount: number
  currency: string
  payment_method: string
  payment_status: 'unpaid' | 'paid'
  paid_date?: string
  paid_by?: string
  submitted_by: string
  submitted_at: string
  approved_by?: string
  approved_at?: string
  receipt_no: string
  attachment_path: string
  notes: string
  created_at: string
  updated_at: string
  
  // Relations
  supplier?: any
  submitter?: any
  approver?: any
  payer?: any
}

export interface AccountReceivable {
  id: string
  company_id: string
  customer_id: string
  invoice_id: string
  invoice_amount: number
  paid_amount: number
  balance_amount: number
  currency: string
  invoice_date: string
  due_date: string
  last_payment_date?: string
  days_overdue: number
  aging_category: 'current' | '30days' | '60days' | '90days' | 'over90days'
  status: 'open' | 'partial' | 'paid' | 'written_off'
  collection_status: 'normal' | 'warning' | 'critical'
  created_at: string
  updated_at: string
  
  // Relations
  customer?: any
  invoice?: Invoice
}

export interface BankAccount {
  id: string
  company_id: string
  account_name: string
  account_no: string
  bank_name: string
  bank_code: string
  branch_name: string
  swift_code: string
  account_type: 'checking' | 'savings' | 'credit'
  currency: string
  current_balance: number
  available_balance: number
  status: 'active' | 'inactive' | 'closed'
  is_default: boolean
  notes: string
  created_at: string
  updated_at: string
}

export interface ARAPSummary {
  total_amount: number
  paid_amount: number
  balance_amount: number
  current: number
  days_30: number
  days_60: number
  days_90: number
  over_90: number
  open_items: number
  currency: string
}

export interface FinancialDashboard {
  revenue: number
  expenses: number
  profit: number
  cash_balance: number
  ar_summary: ARAPSummary
  ap_summary: ARAPSummary
  currency: string
  date: string
}

export interface CreateInvoiceRequest {
  type: string
  order_id?: string
  customer_id?: string
  supplier_id?: string
  issue_date: string
  due_date: string
  sub_total: number
  tax_rate: number
  tax_amount: number
  discount_rate: number
  discount_amount: number
  currency: string
  exchange_rate: number
  payment_terms: string
  payment_method: string
  bank_account: string
  notes: string
  internal_notes: string
}

export interface CreatePaymentRequest {
  type: string
  invoice_id?: string
  customer_id?: string
  supplier_id?: string
  payment_date: string
  amount: number
  currency: string
  exchange_rate: number
  payment_method: string
  bank_name: string
  bank_account: string
  transaction_no: string
  check_no: string
  notes: string
}

export interface CreateExpenseRequest {
  category: string
  sub_category: string
  supplier_id?: string
  department_id?: string
  project_id?: string
  expense_date: string
  description: string
  amount: number
  tax_amount: number
  currency: string
  payment_method: string
  receipt_no: string
  notes: string
}

export interface CreateBankAccountRequest {
  account_name: string
  account_no: string
  bank_name: string
  bank_code: string
  branch_name: string
  swift_code: string
  account_type: string
  currency: string
  current_balance: number
  available_balance: number
  is_default: boolean
  notes: string
}

export interface ListParams {
  page?: number
  page_size?: number
  search?: string
  type?: string
  status?: string
  customer_id?: string
  supplier_id?: string
  start_date?: string
  end_date?: string
  overdue?: boolean
  payment_method?: string
  category?: string
  payment_status?: string
  aging_category?: string
  payment_priority?: string
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

class FinanceService {
  // Invoice operations
  async createInvoice(data: CreateInvoiceRequest): Promise<Invoice> {
    const response = await api.post('/invoices', data)
    return response.data
  }

  async updateInvoice(id: string, data: Partial<Invoice>): Promise<Invoice> {
    const response = await api.put(`/invoices/${id}`, data)
    return response.data
  }

  async getInvoice(id: string): Promise<Invoice> {
    const response = await api.get(`/invoices/${id}`)
    return response.data
  }

  async listInvoices(params: ListParams = {}): Promise<PaginatedResponse<Invoice>> {
    const response = await api.get('/invoices', { params })
    return response.data
  }

  async generateInvoiceFromOrder(orderId: string): Promise<Invoice> {
    const response = await api.post(`/invoices/generate/${orderId}`)
    return response.data
  }

  async getInvoiceItems(invoiceId: string): Promise<InvoiceItem[]> {
    const response = await api.get(`/invoices/${invoiceId}/items`)
    return response.data
  }

  // Payment operations
  async processPayment(data: CreatePaymentRequest): Promise<Payment> {
    const response = await api.post('/payments', data)
    return response.data
  }

  async getPayment(id: string): Promise<Payment> {
    const response = await api.get(`/payments/${id}`)
    return response.data
  }

  async listPayments(params: ListParams = {}): Promise<PaginatedResponse<Payment>> {
    const response = await api.get('/payments', { params })
    return response.data
  }

  async getPaymentsByInvoice(invoiceId: string): Promise<Payment[]> {
    const response = await api.get(`/invoices/${invoiceId}/payments`)
    return response.data
  }

  // Expense operations
  async createExpense(data: CreateExpenseRequest): Promise<Expense> {
    const response = await api.post('/expenses', data)
    return response.data
  }

  async updateExpense(id: string, data: Partial<Expense>): Promise<Expense> {
    const response = await api.put(`/expenses/${id}`, data)
    return response.data
  }

  async getExpense(id: string): Promise<Expense> {
    const response = await api.get(`/expenses/${id}`)
    return response.data
  }

  async listExpenses(params: ListParams = {}): Promise<PaginatedResponse<Expense>> {
    const response = await api.get('/expenses', { params })
    return response.data
  }

  async approveExpense(id: string): Promise<void> {
    await api.post(`/expenses/${id}/approve`)
  }

  async rejectExpense(id: string, reason: string): Promise<void> {
    await api.post(`/expenses/${id}/reject`, { reason })
  }

  // AR/AP operations
  async getAccountReceivables(params: ListParams = {}): Promise<AccountReceivable[]> {
    const response = await api.get('/accounts-receivable', { params })
    return response.data
  }

  async getAccountPayables(params: ListParams = {}): Promise<any[]> {
    const response = await api.get('/accounts-payable', { params })
    return response.data
  }

  async getARSummary(): Promise<ARAPSummary> {
    const response = await api.get('/ar-summary')
    return response.data
  }

  async getAPSummary(): Promise<ARAPSummary> {
    const response = await api.get('/ap-summary')
    return response.data
  }

  // Bank account operations
  async createBankAccount(data: CreateBankAccountRequest): Promise<BankAccount> {
    const response = await api.post('/bank-accounts', data)
    return response.data
  }

  async updateBankAccount(id: string, data: Partial<BankAccount>): Promise<BankAccount> {
    const response = await api.put(`/bank-accounts/${id}`, data)
    return response.data
  }

  async getBankAccount(id: string): Promise<BankAccount> {
    const response = await api.get(`/bank-accounts/${id}`)
    return response.data
  }

  async listBankAccounts(): Promise<BankAccount[]> {
    const response = await api.get('/bank-accounts')
    return response.data
  }

  // Report operations
  async getFinancialDashboard(): Promise<FinancialDashboard> {
    const response = await api.get('/finance/dashboard')
    return response.data
  }

  async getCashFlowReport(startDate: string, endDate: string): Promise<any> {
    const response = await api.get('/finance/cash-flow', {
      params: { start_date: startDate, end_date: endDate }
    })
    return response.data
  }

  async getAgingReport(type: 'receivable' | 'payable'): Promise<any> {
    const response = await api.get('/finance/aging-report', {
      params: { type }
    })
    return response.data
  }
}

export default new FinanceService()