import api from './api'

// Types
export interface TariffCode {
  id: string
  company_id: string
  hs_code: string
  description: string
  description_en: string
  category: string
  unit: string
  base_rate: number
  preferential_rate: number
  vat: number
  excise_tax: number
  import_restriction: string
  export_restriction: string
  required_certs: string
  is_active: boolean
  created_at: string
  updated_at: string
  created_by: string
  
  // Relations
  company?: any
  creator?: any
}

export interface TariffRate {
  id: string
  company_id: string
  tariff_code_id: string
  country_code: string
  country_name: string
  rate: number
  rate_type: 'ad_valorem' | 'specific' | 'compound'
  minimum_duty: number
  maximum_duty: number
  currency: string
  trade_type: 'import' | 'export'
  agreement_type: string
  valid_from: string
  valid_to?: string
  is_active: boolean
  created_at: string
  updated_at: string
  created_by: string
  
  // Relations
  company?: any
  tariff_code?: TariffCode
  creator?: any
}

export interface TradeDocument {
  id: string
  company_id: string
  document_type: string
  document_no: string
  title: string
  description: string
  file_path: string
  file_size: number
  file_type: string
  version: number
  status: 'draft' | 'submitted' | 'approved' | 'rejected' | 'expired'
  is_required: boolean
  valid_from?: string
  valid_to?: string
  issued_by: string
  issued_at?: string
  approved_by?: string
  approved_at?: string
  rejected_reason: string
  metadata: string
  created_at: string
  updated_at: string
  created_by: string
  
  // Relations
  company?: any
  creator?: any
  approver?: any
  shipments?: Shipment[]
}

export interface Shipment {
  id: string
  company_id: string
  shipment_no: string
  order_id?: string
  type: 'import' | 'export'
  status: 'pending' | 'in_transit' | 'customs' | 'delivered' | 'cancelled'
  method: 'sea' | 'air' | 'land' | 'express'
  carrier_name: string
  tracking_no: string
  container_no: string
  container_type: string
  origin_country: string
  origin_port: string
  origin_address: string
  dest_country: string
  dest_port: string
  dest_address: string
  estimated_departure?: string
  actual_departure?: string
  estimated_arrival?: string
  actual_arrival?: string
  gross_weight: number
  net_weight: number
  volume: number
  package_count: number
  package_type: string
  insurance_value: number
  insurance_currency: string
  freight_cost: number
  freight_currency: string
  customs_value: number
  customs_currency: string
  total_duty: number
  total_tax: number
  special_instructions: string
  internal_notes: string
  created_at: string
  updated_at: string
  created_by: string
  
  // Relations
  company?: any
  order?: any
  creator?: any
  documents?: TradeDocument[]
  items?: ShipmentItem[]
  events?: ShipmentEvent[]
}

export interface ShipmentItem {
  id: string
  company_id: string
  shipment_id: string
  product_id?: string
  hs_code: string
  product_name: string
  description: string
  quantity: number
  unit: string
  unit_weight: number
  unit_value: number
  currency: string
  total_weight: number
  total_value: number
  country_origin: string
  manufacturer: string
  created_at: string
  updated_at: string
  
  // Relations
  company?: any
  shipment?: Shipment
  product?: any
}

export interface ShipmentEvent {
  id: string
  company_id: string
  shipment_id: string
  event_type: string
  status: 'completed' | 'in_progress' | 'pending' | 'cancelled'
  location: string
  description: string
  longitude: number
  latitude: number
  event_time: string
  recorded_at: string
  source: string
  created_at: string
  created_by: string
  
  // Relations
  company?: any
  shipment?: Shipment
  creator?: any
}

export interface LetterOfCredit {
  id: string
  company_id: string
  lc_number: string
  type: 'sight' | 'usance' | 'revolving' | 'standby'
  status: 'draft' | 'issued' | 'advised' | 'confirmed' | 'utilized' | 'expired'
  amount: number
  currency: string
  applicant_name: string
  applicant_address: string
  beneficiary_name: string
  beneficiary_address: string
  issuing_bank: string
  advising_bank: string
  confirming_bank: string
  issue_date: string
  expiry_date: string
  last_shipment_date?: string
  partial_shipment: boolean
  transhipment: boolean
  port_of_loading: string
  port_of_discharge: string
  description: string
  documents: string
  terms: string
  utilized_amount: number
  available_amount: number
  amendments: string
  created_at: string
  updated_at: string
  created_by: string
  
  // Relations
  company?: any
  creator?: any
  shipments?: Shipment[]
  utilizations?: LCUtilization[]
}

export interface LCUtilization {
  id: string
  company_id: string
  lc_id: string
  shipment_id?: string
  amount: number
  currency: string
  description: string
  documents_ref: string
  status: 'pending' | 'accepted' | 'rejected'
  utilized_at: string
  created_at: string
  created_by: string
  
  // Relations
  company?: any
  lc?: LetterOfCredit
  shipment?: Shipment
  creator?: any
}

export interface TradeCompliance {
  id: string
  company_id: string
  compliance_type: string
  rule_name: string
  description: string
  country_code: string
  product_codes: string
  entity_list: string
  rule_details: string
  severity: 'low' | 'medium' | 'high' | 'critical'
  status: 'active' | 'suspended' | 'expired'
  valid_from: string
  valid_to?: string
  source: string
  last_updated: string
  updated_by?: string
  created_at: string
  updated_at: string
  created_by: string
  
  // Relations
  company?: any
  creator?: any
  updater?: any
  checks?: ComplianceCheck[]
}

export interface ComplianceCheck {
  id: string
  company_id: string
  compliance_id: string
  resource_type: string
  resource_id: string
  check_type: 'automatic' | 'manual'
  result: 'passed' | 'failed' | 'warning' | 'pending'
  score: number
  issues: string
  recommendations: string
  checked_at: string
  checked_by?: string
  resolved_at?: string
  resolved_by?: string
  notes: string
  created_at: string
  updated_at: string
  
  // Relations
  company?: any
  compliance?: TradeCompliance
  checker?: any
  resolver?: any
}

export interface ExchangeRate {
  id: string
  company_id: string
  from_currency: string
  to_currency: string
  rate: number
  rate_type: 'buy' | 'sell' | 'mid' | 'official'
  source: string
  valid_date: string
  is_active: boolean
  created_at: string
  updated_at: string
  created_by: string
  
  // Relations
  company?: any
  creator?: any
}

export interface TradeRegulation {
  id: string
  company_id?: string
  regulation_type: string
  country_code: string
  regulation_code: string
  title: string
  description: string
  requirements: string
  penalties: string
  documents_needed: string
  processing_time: number
  fees: string
  applicable_hs: string
  effective_date: string
  expiry_date?: string
  status: 'active' | 'suspended' | 'repealed'
  official_url: string
  last_review_date?: string
  reviewed_by?: string
  created_at: string
  updated_at: string
  created_by: string
  
  // Relations
  company?: any
  creator?: any
  reviewer?: any
}

export interface TradeAgreement {
  id: string
  company_id?: string
  agreement_type: string
  agreement_code: string
  name: string
  description: string
  countries: string
  benefits: string
  tariff_reductions: string
  quota_limits: string
  rules_of_origin: string
  effective_date: string
  expiry_date?: string
  status: 'active' | 'suspended' | 'expired' | 'negotiating'
  official_url: string
  created_at: string
  updated_at: string
  created_by: string
  
  // Relations
  company?: any
  creator?: any
}

// Request types
export interface CreateTariffCodeRequest {
  hs_code: string
  description: string
  description_en?: string
  category?: string
  unit?: string
  base_rate?: number
  preferential_rate?: number
  vat?: number
  excise_tax?: number
  import_restriction?: Record<string, any>
  export_restriction?: Record<string, any>
  required_certs?: string[]
}

export interface UpdateTariffCodeRequest {
  description?: string
  description_en?: string
  category?: string
  unit?: string
  base_rate?: number
  preferential_rate?: number
  vat?: number
  excise_tax?: number
  import_restriction?: Record<string, any>
  export_restriction?: Record<string, any>
  required_certs?: string[]
  is_active?: boolean
}

export interface CreateTariffRateRequest {
  tariff_code_id: string
  country_code: string
  country_name: string
  rate: number
  rate_type: string
  minimum_duty?: number
  maximum_duty?: number
  currency?: string
  trade_type: string
  agreement_type?: string
  valid_from: string
  valid_to?: string
}

export interface CreateShipmentRequest {
  order_id?: string
  shipment_no: string
  type: string
  method: string
  carrier_name?: string
  tracking_no?: string
  container_no?: string
  container_type?: string
  origin_country: string
  origin_port?: string
  origin_address?: string
  dest_country: string
  dest_port?: string
  dest_address?: string
  estimated_departure?: string
  estimated_arrival?: string
  gross_weight?: number
  net_weight?: number
  volume?: number
  package_count?: number
  package_type?: string
  insurance_value?: number
  insurance_currency?: string
  freight_cost?: number
  freight_currency?: string
  customs_value?: number
  customs_currency?: string
  special_instructions?: string
  internal_notes?: string
  items?: CreateShipmentItemRequest[]
}

export interface CreateShipmentItemRequest {
  product_id?: string
  hs_code?: string
  product_name: string
  description?: string
  quantity: number
  unit: string
  unit_weight?: number
  unit_value?: number
  currency?: string
  country_origin?: string
  manufacturer?: string
}

export interface UpdateShipmentRequest {
  status?: string
  tracking_no?: string
  actual_departure?: string
  actual_arrival?: string
  total_duty?: number
  total_tax?: number
  special_instructions?: string
  internal_notes?: string
}

export interface CreateShipmentEventRequest {
  event_type: string
  status: string
  location?: string
  description?: string
  longitude?: number
  latitude?: number
  event_time: string
  source?: string
}

export interface CreateLetterOfCreditRequest {
  lc_number: string
  type: string
  amount: number
  currency: string
  applicant_name: string
  applicant_address?: string
  beneficiary_name: string
  beneficiary_address?: string
  issuing_bank: string
  advising_bank?: string
  confirming_bank?: string
  issue_date: string
  expiry_date: string
  last_shipment_date?: string
  partial_shipment?: boolean
  transhipment?: boolean
  port_of_loading?: string
  port_of_discharge?: string
  description?: string
  documents?: string[]
  terms?: Record<string, any>
}

export interface CreateLCUtilizationRequest {
  lc_id: string
  shipment_id?: string
  amount: number
  currency: string
  description?: string
  documents_ref?: Record<string, any>
  utilized_at: string
}

export interface CreateComplianceCheckRequest {
  compliance_id: string
  resource_type: string
  resource_id: string
  check_type: string
}

export interface CreateExchangeRateRequest {
  from_currency: string
  to_currency: string
  rate: number
  rate_type: string
  source: string
  valid_date: string
}

export interface ListParams {
  hs_code?: string
  category?: string
  is_active?: boolean
  country_code?: string
  trade_type?: string
  agreement_type?: string
  type?: string
  status?: string
  method?: string
  resource_type?: string
  resource_id?: string
  result?: string
  from?: string
  to?: string
  rate_type?: string
  start_date?: string
  end_date?: string
  days?: number
  limit?: number
}

export interface PaginatedResponse<T> {
  data: T[]
  total?: number
}

export interface TariffDutyCalculation {
  tariff_code: TariffCode
  tariff_rate?: TariffRate
  applied_rate: number
  duty: number
  vat: number
  excise_tax: number
  total: number
}

export interface CurrencyConversion {
  original_amount: number
  original_currency: string
  converted_amount: number
  converted_currency: string
  rate: number
  rate_date: string
  rate_source: string
}

export interface TradeStatistics {
  shipments: {
    total_shipments: number
    import_shipments: number
    export_shipments: number
    in_transit_shipments: number
    delivered_shipments: number
    total_value: number
    total_weight: number
  }
  letter_of_credits: {
    total_lcs: number
    active_lcs: number
    utilized_lcs: number
    expiring_lcs: number
    total_amount: number
    utilized_amount: number
  }
  compliance: {
    total_compliances: number
    active_compliances: number
    failed_checks: number
    warning_checks: number
    passed_checks: number
  }
}

class TradeService {
  // TariffCode Methods
  async createTariffCode(data: CreateTariffCodeRequest): Promise<TariffCode> {
    const response = await api.post('/trade/tariff-codes', data)
    return response.data
  }

  async getTariffCode(id: string): Promise<TariffCode> {
    const response = await api.get(`/trade/tariff-codes/${id}`)
    return response.data
  }

  async listTariffCodes(params: ListParams = {}): Promise<PaginatedResponse<TariffCode>> {
    const response = await api.get('/trade/tariff-codes', { params })
    return response.data
  }

  async updateTariffCode(id: string, data: UpdateTariffCodeRequest): Promise<TariffCode> {
    const response = await api.put(`/trade/tariff-codes/${id}`, data)
    return response.data
  }

  async deleteTariffCode(id: string): Promise<void> {
    await api.delete(`/trade/tariff-codes/${id}`)
  }

  // TariffRate Methods
  async createTariffRate(data: CreateTariffRateRequest): Promise<TariffRate> {
    const response = await api.post('/trade/tariff-rates', data)
    return response.data
  }

  async getTariffRatesByTariffCode(tariffCodeId: string, params: ListParams = {}): Promise<PaginatedResponse<TariffRate>> {
    const response = await api.get(`/trade/tariff-codes/${tariffCodeId}/rates`, { params })
    return response.data
  }

  async listTariffRates(params: ListParams = {}): Promise<PaginatedResponse<TariffRate>> {
    const response = await api.get('/trade/tariff-rates', { params })
    return response.data
  }

  // Shipment Methods
  async createShipment(data: CreateShipmentRequest): Promise<Shipment> {
    const response = await api.post('/trade/shipments', data)
    return response.data
  }

  async getShipment(id: string): Promise<Shipment> {
    const response = await api.get(`/trade/shipments/${id}`)
    return response.data
  }

  async listShipments(params: ListParams = {}): Promise<PaginatedResponse<Shipment>> {
    const response = await api.get('/trade/shipments', { params })
    return response.data
  }

  async updateShipment(id: string, data: UpdateShipmentRequest): Promise<Shipment> {
    const response = await api.put(`/trade/shipments/${id}`, data)
    return response.data
  }

  // ShipmentEvent Methods
  async createShipmentEvent(shipmentId: string, data: CreateShipmentEventRequest): Promise<ShipmentEvent> {
    const response = await api.post(`/trade/shipments/${shipmentId}/events`, data)
    return response.data
  }

  async getShipmentEvents(shipmentId: string): Promise<PaginatedResponse<ShipmentEvent>> {
    const response = await api.get(`/trade/shipments/${shipmentId}/events`)
    return response.data
  }

  // LetterOfCredit Methods
  async createLetterOfCredit(data: CreateLetterOfCreditRequest): Promise<LetterOfCredit> {
    const response = await api.post('/trade/letter-of-credits', data)
    return response.data
  }

  async getLetterOfCredit(id: string): Promise<LetterOfCredit> {
    const response = await api.get(`/trade/letter-of-credits/${id}`)
    return response.data
  }

  async listLetterOfCredits(params: ListParams = {}): Promise<PaginatedResponse<LetterOfCredit>> {
    const response = await api.get('/trade/letter-of-credits', { params })
    return response.data
  }

  async getExpiringLetterOfCredits(days: number = 30): Promise<PaginatedResponse<LetterOfCredit>> {
    const response = await api.get('/trade/letter-of-credits/expiring', { params: { days } })
    return response.data
  }

  // LCUtilization Methods
  async createLCUtilization(data: CreateLCUtilizationRequest): Promise<LCUtilization> {
    const response = await api.post('/trade/lc-utilizations', data)
    return response.data
  }

  async getLCUtilizations(lcId: string, params: ListParams = {}): Promise<PaginatedResponse<LCUtilization>> {
    const response = await api.get(`/trade/letter-of-credits/${lcId}/utilizations`, { params })
    return response.data
  }

  // Compliance Methods
  async runComplianceCheck(data: CreateComplianceCheckRequest): Promise<ComplianceCheck> {
    const response = await api.post('/trade/compliance/check', data)
    return response.data
  }

  async getComplianceChecksByResource(params: ListParams): Promise<PaginatedResponse<ComplianceCheck>> {
    const response = await api.get('/trade/compliance/checks', { params })
    return response.data
  }

  async getFailedComplianceChecks(): Promise<PaginatedResponse<ComplianceCheck>> {
    const response = await api.get('/trade/compliance/failed-checks')
    return response.data
  }

  // ExchangeRate Methods
  async createExchangeRate(data: CreateExchangeRateRequest): Promise<ExchangeRate> {
    const response = await api.post('/trade/exchange-rates', data)
    return response.data
  }

  async getLatestExchangeRate(from: string, to: string, type: string = 'mid'): Promise<ExchangeRate> {
    const response = await api.get('/trade/exchange-rates/latest', { 
      params: { from, to, type } 
    })
    return response.data
  }

  async listExchangeRates(params: ListParams = {}): Promise<PaginatedResponse<ExchangeRate>> {
    const response = await api.get('/trade/exchange-rates', { params })
    return response.data
  }

  // Analytics Methods
  async getTradeStatistics(params: ListParams = {}): Promise<TradeStatistics> {
    const response = await api.get('/trade/analytics/statistics', { params })
    return response.data
  }

  async getShipmentsByCountry(params: ListParams = {}): Promise<{ data: Array<{ origin_country: string; dest_country: string; shipment_count: number; total_value: number }> }> {
    const response = await api.get('/trade/analytics/shipments-by-country', { params })
    return response.data
  }

  async getTopTradingPartners(params: ListParams = {}): Promise<{ data: Array<{ country: string; type: string; shipment_count: number; total_value: number }> }> {
    const response = await api.get('/trade/analytics/top-trading-partners', { params })
    return response.data
  }

  // Utility Methods
  async calculateTariffDuty(data: {
    hs_code: string
    country_code: string
    trade_type: string
    value: number
  }): Promise<TariffDutyCalculation> {
    const response = await api.post('/trade/utils/calculate-tariff-duty', data)
    return response.data
  }

  async convertCurrency(data: {
    amount: number
    from_currency: string
    to_currency: string
  }): Promise<CurrencyConversion> {
    const response = await api.post('/trade/utils/convert-currency', data)
    return response.data
  }

  async getTradeDocumentsByShipment(shipmentId: string): Promise<PaginatedResponse<TradeDocument>> {
    const response = await api.get(`/trade/shipments/${shipmentId}/documents`)
    return response.data
  }

  // Helper Methods
  getShipmentStatusColor(status: string): string {
    const statusColors: Record<string, string> = {
      pending: 'text-yellow-600 bg-yellow-100',
      in_transit: 'text-blue-600 bg-blue-100',
      customs: 'text-orange-600 bg-orange-100',
      delivered: 'text-green-600 bg-green-100',
      cancelled: 'text-red-600 bg-red-100',
    }
    return statusColors[status] || 'text-gray-600 bg-gray-100'
  }

  getShipmentTypeIcon(type: string): string {
    const typeIcons: Record<string, string> = {
      import: '📥',
      export: '📤',
    }
    return typeIcons[type] || '📦'
  }

  getShipmentMethodIcon(method: string): string {
    const methodIcons: Record<string, string> = {
      sea: '🚢',
      air: '✈️',
      land: '🚛',
      express: '📦',
    }
    return methodIcons[method] || '🚚'
  }

  getLCStatusColor(status: string): string {
    const statusColors: Record<string, string> = {
      draft: 'text-gray-600 bg-gray-100',
      issued: 'text-blue-600 bg-blue-100',
      advised: 'text-yellow-600 bg-yellow-100',
      confirmed: 'text-green-600 bg-green-100',
      utilized: 'text-purple-600 bg-purple-100',
      expired: 'text-red-600 bg-red-100',
    }
    return statusColors[status] || 'text-gray-600 bg-gray-100'
  }

  getComplianceResultColor(result: string): string {
    const resultColors: Record<string, string> = {
      passed: 'text-green-600 bg-green-100',
      failed: 'text-red-600 bg-red-100',
      warning: 'text-yellow-600 bg-yellow-100',
      pending: 'text-blue-600 bg-blue-100',
    }
    return resultColors[result] || 'text-gray-600 bg-gray-100'
  }

  formatCurrency(amount: number, currency: string = 'USD'): string {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: currency,
    }).format(amount)
  }

  formatWeight(weight: number, unit: string = 'kg'): string {
    return `${weight.toLocaleString()} ${unit}`
  }

  formatVolume(volume: number, unit: string = 'm³'): string {
    return `${volume.toLocaleString()} ${unit}`
  }
}

export default new TradeService()