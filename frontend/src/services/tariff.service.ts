import api from './api'

export interface HSCode {
  id: string
  code: string
  description: string
  description_en?: string
  unit: string
  category: string
  parent_code?: string
  is_active: boolean
}

export interface TariffRate {
  id: string
  hs_code: string
  from_country: string
  to_country: string
  rate_type: 'ad_valorem' | 'specific' | 'compound'
  rate_value: number
  specific_rate?: number
  currency?: string
  unit?: string
  effective_from: string
  effective_to?: string
  preferential_rate?: number
  preferential_conditions?: string
  notes?: string
}

export interface TariffCalculationRequest {
  hs_code: string
  from_country: string
  to_country: string
  product_value: number
  quantity: number
  unit: string
  weight_kg?: number
  currency: string
  incoterm?: string
  preferential_treatment?: boolean
}

export interface TariffCalculationResult {
  hs_code: string
  from_country: string
  to_country: string
  product_value: number
  tariff_rate: TariffRate
  calculated_tariff: number
  effective_rate: number
  currency: string
  calculation_details: {
    base_value: number
    ad_valorem_duty?: number
    specific_duty?: number
    total_duty: number
    preferential_applied: boolean
    preferential_savings?: number
  }
  warnings?: string[]
}

export interface CountryTradeAgreement {
  id: string
  agreement_name: string
  member_countries: string[]
  effective_date: string
  preferential_rates: boolean
  certificate_required: boolean
  description?: string
}

class TariffService {
  async searchHSCodes(params: {
    search?: string
    category?: string
    page?: number
    page_size?: number
  }): Promise<{ data: HSCode[]; pagination: any }> {
    const response = await api.get('/tariffs/hs-codes', { params })
    return response.data
  }

  async getHSCode(code: string): Promise<HSCode> {
    const response = await api.get(`/tariffs/hs-codes/${code}`)
    return response.data
  }

  async getTariffRates(params: {
    hs_code: string
    from_country?: string
    to_country?: string
  }): Promise<TariffRate[]> {
    const response = await api.get('/tariffs/rates', { params })
    return response.data
  }

  async calculateTariff(data: TariffCalculationRequest): Promise<TariffCalculationResult> {
    const response = await api.post('/tariffs/calculate', data)
    return response.data
  }

  async batchCalculate(
    items: TariffCalculationRequest[]
  ): Promise<TariffCalculationResult[]> {
    const response = await api.post('/tariffs/batch-calculate', { items })
    return response.data
  }

  async getTradeAgreements(countries: string[]): Promise<CountryTradeAgreement[]> {
    const response = await api.get('/tariffs/trade-agreements', {
      params: { countries: countries.join(',') }
    })
    return response.data
  }

  async validateHSCode(code: string, country: string): Promise<{
    valid: boolean
    message?: string
    suggested_codes?: HSCode[]
  }> {
    const response = await api.post('/tariffs/validate-hs-code', { code, country })
    return response.data
  }

  async getCommonHSCodes(category?: string): Promise<HSCode[]> {
    const response = await api.get('/tariffs/common-hs-codes', { params: { category } })
    return response.data
  }

  // Static data for demo/testing
  async getCountryList(): Promise<Array<{ code: string; name: string; region: string }>> {
    return [
      { code: 'TW', name: '台灣 Taiwan', region: 'Asia' },
      { code: 'CN', name: '中國 China', region: 'Asia' },
      { code: 'US', name: '美國 United States', region: 'Americas' },
      { code: 'DE', name: '德國 Germany', region: 'Europe' },
      { code: 'JP', name: '日本 Japan', region: 'Asia' },
      { code: 'KR', name: '韓國 South Korea', region: 'Asia' },
      { code: 'VN', name: '越南 Vietnam', region: 'Asia' },
      { code: 'TH', name: '泰國 Thailand', region: 'Asia' },
      { code: 'MY', name: '馬來西亞 Malaysia', region: 'Asia' },
      { code: 'SG', name: '新加坡 Singapore', region: 'Asia' },
      { code: 'GB', name: '英國 United Kingdom', region: 'Europe' },
      { code: 'FR', name: '法國 France', region: 'Europe' },
      { code: 'IT', name: '義大利 Italy', region: 'Europe' },
      { code: 'MX', name: '墨西哥 Mexico', region: 'Americas' },
      { code: 'CA', name: '加拿大 Canada', region: 'Americas' },
      { code: 'AU', name: '澳洲 Australia', region: 'Oceania' },
      { code: 'IN', name: '印度 India', region: 'Asia' },
      { code: 'BR', name: '巴西 Brazil', region: 'Americas' },
    ]
  }

  // Common HS codes for fasteners
  async getFastenerHSCodes(): Promise<Array<{ code: string; description: string }>> {
    return [
      { code: '7318.11', description: '鐵或鋼製螺栓 (Coach screws)' },
      { code: '7318.12', description: '鐵或鋼製其他木螺釘 (Other wood screws)' },
      { code: '7318.13', description: '鐵或鋼製鉤頭螺釘及環頭螺釘 (Screw hooks and screw rings)' },
      { code: '7318.14', description: '鐵或鋼製自攻螺釘 (Self-tapping screws)' },
      { code: '7318.15', description: '鐵或鋼製其他螺釘及螺栓 (Other screws and bolts)' },
      { code: '7318.16', description: '鐵或鋼製螺帽 (Nuts)' },
      { code: '7318.19', description: '鐵或鋼製其他螺紋製品 (Other threaded articles)' },
      { code: '7318.21', description: '鐵或鋼製彈簧墊圈及其他防鬆墊圈 (Spring washers and other lock washers)' },
      { code: '7318.22', description: '鐵或鋼製其他墊圈 (Other washers)' },
      { code: '7318.23', description: '鐵或鋼製鉚釘 (Rivets)' },
      { code: '7318.24', description: '鐵或鋼製開口銷及開尾銷 (Cotters and cotter-pins)' },
      { code: '7318.29', description: '鐵或鋼製其他無螺紋製品 (Other non-threaded articles)' },
    ]
  }
}

export const tariffService = new TariffService()
export default tariffService