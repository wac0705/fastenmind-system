'use client'

import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import DashboardLayout from '@/components/layout/DashboardLayout'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { 
  Calculator, 
  Globe, 
  Package, 
  AlertCircle,
  TrendingUp,
  FileText,
  Search,
  Info,
  DollarSign
} from 'lucide-react'
import { useToast } from '@/components/ui/use-toast'
import tariffService, { 
  TariffCalculationRequest, 
  TariffCalculationResult 
} from '@/services/tariff.service'

export default function TariffCalculatorPage() {
  const { toast } = useToast()
  const [calculationData, setCalculationData] = useState<TariffCalculationRequest>({
    hs_code: '',
    from_country: 'TW',
    to_country: 'US',
    product_value: 0,
    quantity: 1000,
    unit: 'PCS',
    weight_kg: 0,
    currency: 'USD',
    incoterm: 'FOB',
    preferential_treatment: false,
  })

  const [calculationResult, setCalculationResult] = useState<TariffCalculationResult | null>(null)
  const [isCalculating, setIsCalculating] = useState(false)

  // Fetch country options
  const { data: countries = [] } = useQuery({
    queryKey: ['countries'],
    queryFn: () => tariffService.getCountryList(),
  })

  // Fetch common HS codes
  const { data: commonHSCodes = [] } = useQuery({
    queryKey: ['fastener-hs-codes'],
    queryFn: () => tariffService.getFastenerHSCodes(),
  })

  const handleCalculate = async () => {
    if (!calculationData.hs_code || !calculationData.product_value) {
      toast({
        title: '請填寫必要欄位',
        description: 'HS Code 和產品價值為必填',
        variant: 'destructive',
      })
      return
    }

    setIsCalculating(true)
    try {
      const result = await tariffService.calculateTariff(calculationData)
      setCalculationResult(result)
      toast({ title: '關稅計算完成' })
    } catch (error: any) {
      toast({
        title: '計算失敗',
        description: error.response?.data?.message || '關稅計算時發生錯誤',
        variant: 'destructive',
      })
    } finally {
      setIsCalculating(false)
    }
  }

  const incoterms = [
    { value: 'EXW', label: 'EXW - 工廠交貨' },
    { value: 'FCA', label: 'FCA - 貨交承運人' },
    { value: 'CPT', label: 'CPT - 運費付至' },
    { value: 'CIP', label: 'CIP - 運費保險費付至' },
    { value: 'DAP', label: 'DAP - 目的地交貨' },
    { value: 'DPU', label: 'DPU - 目的地卸貨後交貨' },
    { value: 'DDP', label: 'DDP - 完稅後交貨' },
    { value: 'FAS', label: 'FAS - 船邊交貨' },
    { value: 'FOB', label: 'FOB - 船上交貨' },
    { value: 'CFR', label: 'CFR - 成本加運費' },
    { value: 'CIF', label: 'CIF - 成本保險費加運費' },
  ]

  const units = ['PCS', 'KGS', 'SET', 'BOX', 'CTN']

  return (
    <DashboardLayout>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex justify-between items-center">
          <div>
            <h1 className="text-3xl font-bold text-gray-900 flex items-center gap-2">
              <Calculator className="h-8 w-8" />
              關稅計算器
            </h1>
            <p className="mt-2 text-gray-600">計算進出口關稅與貿易成本</p>
          </div>
        </div>

        <Tabs defaultValue="calculator" className="space-y-4">
          <TabsList>
            <TabsTrigger value="calculator">關稅計算</TabsTrigger>
            <TabsTrigger value="hs-codes">HS Code 查詢</TabsTrigger>
            <TabsTrigger value="agreements">貿易協定</TabsTrigger>
          </TabsList>

          <TabsContent value="calculator" className="space-y-6">
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              {/* Input Form */}
              <Card>
                <CardHeader>
                  <CardTitle>計算參數</CardTitle>
                  <CardDescription>輸入產品與貿易資訊</CardDescription>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="grid grid-cols-2 gap-4">
                    <div className="col-span-2">
                      <Label htmlFor="hs_code">HS Code *</Label>
                      <div className="flex gap-2">
                        <Input
                          id="hs_code"
                          value={calculationData.hs_code}
                          onChange={(e) => setCalculationData({ ...calculationData, hs_code: e.target.value })}
                          placeholder="例如: 7318.15"
                        />
                        <Select
                          value={calculationData.hs_code}
                          onValueChange={(value) => setCalculationData({ ...calculationData, hs_code: value })}
                        >
                          <SelectTrigger className="w-[200px]">
                            <SelectValue placeholder="常用 HS Code" />
                          </SelectTrigger>
                          <SelectContent>
                            {commonHSCodes.map((hs) => (
                              <SelectItem key={hs.code} value={hs.code}>
                                {hs.code}
                              </SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                      </div>
                      {calculationData.hs_code && (
                        <p className="text-sm text-gray-500 mt-1">
                          {commonHSCodes.find(hs => hs.code === calculationData.hs_code)?.description}
                        </p>
                      )}
                    </div>

                    <div>
                      <Label htmlFor="from_country">出口國 *</Label>
                      <Select
                        value={calculationData.from_country}
                        onValueChange={(value) => setCalculationData({ ...calculationData, from_country: value })}
                      >
                        <SelectTrigger>
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          {countries.map((country) => (
                            <SelectItem key={country.code} value={country.code}>
                              {country.name}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </div>

                    <div>
                      <Label htmlFor="to_country">進口國 *</Label>
                      <Select
                        value={calculationData.to_country}
                        onValueChange={(value) => setCalculationData({ ...calculationData, to_country: value })}
                      >
                        <SelectTrigger>
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          {countries.map((country) => (
                            <SelectItem key={country.code} value={country.code}>
                              {country.name}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </div>

                    <div>
                      <Label htmlFor="product_value">產品價值 *</Label>
                      <Input
                        id="product_value"
                        type="number"
                        value={calculationData.product_value}
                        onChange={(e) => setCalculationData({ ...calculationData, product_value: parseFloat(e.target.value) || 0 })}
                        min="0"
                        step="0.01"
                      />
                    </div>

                    <div>
                      <Label htmlFor="currency">幣別</Label>
                      <Select
                        value={calculationData.currency}
                        onValueChange={(value) => setCalculationData({ ...calculationData, currency: value })}
                      >
                        <SelectTrigger>
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="USD">USD</SelectItem>
                          <SelectItem value="EUR">EUR</SelectItem>
                          <SelectItem value="TWD">TWD</SelectItem>
                          <SelectItem value="CNY">CNY</SelectItem>
                          <SelectItem value="JPY">JPY</SelectItem>
                        </SelectContent>
                      </Select>
                    </div>

                    <div>
                      <Label htmlFor="quantity">數量</Label>
                      <Input
                        id="quantity"
                        type="number"
                        value={calculationData.quantity}
                        onChange={(e) => setCalculationData({ ...calculationData, quantity: parseFloat(e.target.value) || 0 })}
                        min="0"
                      />
                    </div>

                    <div>
                      <Label htmlFor="unit">單位</Label>
                      <Select
                        value={calculationData.unit}
                        onValueChange={(value) => setCalculationData({ ...calculationData, unit: value })}
                      >
                        <SelectTrigger>
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          {units.map((unit) => (
                            <SelectItem key={unit} value={unit}>
                              {unit}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </div>

                    <div>
                      <Label htmlFor="weight_kg">重量 (KG)</Label>
                      <Input
                        id="weight_kg"
                        type="number"
                        value={calculationData.weight_kg}
                        onChange={(e) => setCalculationData({ ...calculationData, weight_kg: parseFloat(e.target.value) || 0 })}
                        min="0"
                        step="0.01"
                      />
                    </div>

                    <div>
                      <Label htmlFor="incoterm">國貿條件</Label>
                      <Select
                        value={calculationData.incoterm}
                        onValueChange={(value) => setCalculationData({ ...calculationData, incoterm: value })}
                      >
                        <SelectTrigger>
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          {incoterms.map((term) => (
                            <SelectItem key={term.value} value={term.value}>
                              {term.label}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </div>
                  </div>

                  <div className="flex items-center space-x-2">
                    <input
                      type="checkbox"
                      id="preferential"
                      checked={calculationData.preferential_treatment}
                      onChange={(e) => setCalculationData({ ...calculationData, preferential_treatment: e.target.checked })}
                      className="rounded border-gray-300"
                    />
                    <Label htmlFor="preferential">申請優惠關稅待遇</Label>
                  </div>

                  <Button 
                    onClick={handleCalculate} 
                    className="w-full"
                    disabled={isCalculating}
                  >
                    {isCalculating ? '計算中...' : '計算關稅'}
                  </Button>
                </CardContent>
              </Card>

              {/* Result Display */}
              <Card>
                <CardHeader>
                  <CardTitle>計算結果</CardTitle>
                  <CardDescription>關稅明細與成本分析</CardDescription>
                </CardHeader>
                <CardContent>
                  {calculationResult ? (
                    <div className="space-y-4">
                      {/* Summary */}
                      <div className="bg-blue-50 p-4 rounded-lg">
                        <div className="flex items-center justify-between mb-2">
                          <span className="text-sm text-gray-600">總關稅</span>
                          <span className="text-2xl font-bold text-blue-600">
                            {calculationResult.currency} {calculationResult.calculated_tariff.toFixed(2)}
                          </span>
                        </div>
                        <div className="flex items-center justify-between">
                          <span className="text-sm text-gray-600">有效稅率</span>
                          <span className="text-lg font-medium">
                            {(calculationResult.effective_rate * 100).toFixed(2)}%
                          </span>
                        </div>
                      </div>

                      {/* Details */}
                      <div className="space-y-3">
                        <div className="flex justify-between py-2 border-b">
                          <span className="text-gray-600">HS Code</span>
                          <span className="font-medium">{calculationResult.hs_code}</span>
                        </div>
                        <div className="flex justify-between py-2 border-b">
                          <span className="text-gray-600">路線</span>
                          <span className="font-medium">
                            {calculationResult.from_country} → {calculationResult.to_country}
                          </span>
                        </div>
                        <div className="flex justify-between py-2 border-b">
                          <span className="text-gray-600">產品價值</span>
                          <span className="font-medium">
                            {calculationResult.currency} {calculationResult.product_value.toFixed(2)}
                          </span>
                        </div>

                        {/* Calculation Breakdown */}
                        <div className="bg-gray-50 p-3 rounded-lg space-y-2">
                          <p className="font-medium text-sm mb-2">計算明細</p>
                          <div className="flex justify-between text-sm">
                            <span className="text-gray-600">基礎價值</span>
                            <span>{calculationResult.currency} {calculationResult.calculation_details.base_value.toFixed(2)}</span>
                          </div>
                          {calculationResult.calculation_details.ad_valorem_duty && (
                            <div className="flex justify-between text-sm">
                              <span className="text-gray-600">從價稅</span>
                              <span>{calculationResult.currency} {calculationResult.calculation_details.ad_valorem_duty.toFixed(2)}</span>
                            </div>
                          )}
                          {calculationResult.calculation_details.specific_duty && (
                            <div className="flex justify-between text-sm">
                              <span className="text-gray-600">從量稅</span>
                              <span>{calculationResult.currency} {calculationResult.calculation_details.specific_duty.toFixed(2)}</span>
                            </div>
                          )}
                          <div className="flex justify-between text-sm font-medium pt-2 border-t">
                            <span>總關稅</span>
                            <span>{calculationResult.currency} {calculationResult.calculation_details.total_duty.toFixed(2)}</span>
                          </div>
                        </div>

                        {calculationResult.calculation_details.preferential_applied && (
                          <Alert className="bg-green-50 border-green-200">
                            <AlertCircle className="h-4 w-4 text-green-600" />
                            <AlertTitle className="text-green-800">已套用優惠稅率</AlertTitle>
                            <AlertDescription className="text-green-700">
                              節省: {calculationResult.currency} {calculationResult.calculation_details.preferential_savings?.toFixed(2) || '0.00'}
                            </AlertDescription>
                          </Alert>
                        )}

                        {calculationResult.warnings && calculationResult.warnings.length > 0 && (
                          <Alert>
                            <AlertCircle className="h-4 w-4" />
                            <AlertTitle>注意事項</AlertTitle>
                            <AlertDescription>
                              <ul className="list-disc list-inside">
                                {calculationResult.warnings.map((warning, index) => (
                                  <li key={index}>{warning}</li>
                                ))}
                              </ul>
                            </AlertDescription>
                          </Alert>
                        )}
                      </div>
                    </div>
                  ) : (
                    <div className="text-center py-12 text-gray-500">
                      <Calculator className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                      <p>請輸入參數開始計算</p>
                    </div>
                  )}
                </CardContent>
              </Card>
            </div>

            {/* Info Cards */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <Card>
                <CardHeader className="pb-3">
                  <CardTitle className="text-sm flex items-center gap-2">
                    <Info className="h-4 w-4" />
                    關稅類型
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <p className="text-sm text-gray-600">
                    <strong>從價稅 (Ad Valorem):</strong> 按貨物價值的百分比計算
                  </p>
                  <p className="text-sm text-gray-600 mt-2">
                    <strong>從量稅 (Specific):</strong> 按貨物數量或重量計算
                  </p>
                  <p className="text-sm text-gray-600 mt-2">
                    <strong>混合稅 (Compound):</strong> 結合從價稅和從量稅
                  </p>
                </CardContent>
              </Card>

              <Card>
                <CardHeader className="pb-3">
                  <CardTitle className="text-sm flex items-center gap-2">
                    <Globe className="h-4 w-4" />
                    優惠待遇
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <p className="text-sm text-gray-600">
                    符合自由貿易協定（FTA）或優惠貿易安排的產品可享受降低或免除關稅。
                    需要提供原產地證明文件。
                  </p>
                </CardContent>
              </Card>

              <Card>
                <CardHeader className="pb-3">
                  <CardTitle className="text-sm flex items-center gap-2">
                    <FileText className="h-4 w-4" />
                    注意事項
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <p className="text-sm text-gray-600">
                    關稅計算僅供參考，實際關稅可能因海關估價、匯率變動等因素而有所不同。
                    建議諮詢專業報關行。
                  </p>
                </CardContent>
              </Card>
            </div>
          </TabsContent>

          <TabsContent value="hs-codes" className="space-y-4">
            <Card>
              <CardHeader>
                <CardTitle>HS Code 查詢</CardTitle>
                <CardDescription>搜尋和瀏覽協調制度編碼</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div className="relative">
                    <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-4 w-4" />
                    <Input
                      placeholder="搜尋 HS Code 或產品描述..."
                      className="pl-10"
                    />
                  </div>

                  <div className="space-y-2">
                    <h3 className="font-medium">常用緊固件 HS Code</h3>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-2">
                      {commonHSCodes.map((hs) => (
                        <div
                          key={hs.code}
                          className="p-3 border rounded-lg hover:bg-gray-50 cursor-pointer"
                          onClick={() => setCalculationData({ ...calculationData, hs_code: hs.code })}
                        >
                          <div className="flex justify-between items-center">
                            <div>
                              <p className="font-medium">{hs.code}</p>
                              <p className="text-sm text-gray-600">{hs.description}</p>
                            </div>
                            <Badge variant="outline">選擇</Badge>
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="agreements" className="space-y-4">
            <Card>
              <CardHeader>
                <CardTitle>貿易協定查詢</CardTitle>
                <CardDescription>查看各國間的優惠貿易安排</CardDescription>
              </CardHeader>
              <CardContent>
                <Alert>
                  <Info className="h-4 w-4" />
                  <AlertTitle>功能開發中</AlertTitle>
                  <AlertDescription>
                    貿易協定查詢功能正在開發中，敬請期待。
                  </AlertDescription>
                </Alert>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </div>
    </DashboardLayout>
  )
}