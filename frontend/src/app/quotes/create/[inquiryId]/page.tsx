'use client'

import { useState, useEffect } from 'react'
import { useParams, useRouter } from 'next/navigation'
import { useQuery, useMutation } from '@tanstack/react-query'
import DashboardLayout from '@/components/layout/DashboardLayout'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Separator } from '@/components/ui/separator'
import { useToast } from '@/components/ui/use-toast'
import { 
  ArrowLeft,
  Calculator,
  DollarSign,
  Package,
  Truck,
  Globe,
  Wrench,
  Plus,
  Loader2,
  AlertCircle
} from 'lucide-react'
import { format, addDays } from 'date-fns'
import { zhTW } from 'date-fns/locale'
import inquiryService from '@/services/inquiry.service'
import quoteService from '@/services/quote.service'
import processService, { ProcessCostCalculation } from '@/services/process.service'
import { Alert, AlertDescription } from '@/components/ui/alert'

interface CostSection {
  material_cost: number
  process_cost: number
  surface_cost: number
  heat_treat_cost: number
  packaging_cost: number
  shipping_cost: number
  tariff_cost: number
}

export default function CreateQuotePage() {
  const params = useParams()
  const router = useRouter()
  const { toast } = useToast()
  const inquiryId = params.inquiryId as string

  const [costs, setCosts] = useState<CostSection>({
    material_cost: 0,
    process_cost: 0,
    surface_cost: 0,
    heat_treat_cost: 0,
    packaging_cost: 0,
    shipping_cost: 0,
    tariff_cost: 0,
  })
  
  const [overheadRate, setOverheadRate] = useState(15) // 15%
  const [profitRate, setProfitRate] = useState(20) // 20%
  const [deliveryDays, setDeliveryDays] = useState(30)
  const [paymentTerms, setPaymentTerms] = useState('T/T 30 days')
  const [notes, setNotes] = useState('')
  const [processCalculations, setProcessCalculations] = useState<ProcessCostCalculation[]>([])

  // Fetch inquiry details
  const { data: inquiry } = useQuery({
    queryKey: ['inquiry', inquiryId],
    queryFn: () => inquiryService.get(inquiryId),
  })

  // Create quote mutation
  const createQuoteMutation = useMutation({
    mutationFn: (data: any) => quoteService.create(data),
    onSuccess: (quote) => {
      toast({
        title: '報價單建立成功',
        description: `報價單號：${quote.quote_no}`,
      })
      router.push(`/quotes/${quote.id}`)
    },
    onError: (error: any) => {
      toast({
        title: '建立失敗',
        description: error.response?.data?.message || '建立報價單時發生錯誤',
        variant: 'destructive',
      })
    },
  })

  // Calculate totals
  const calculateTotals = () => {
    const subtotal = Object.values(costs).reduce((sum, cost) => sum + cost, 0)
    const overheadAmount = subtotal * (overheadRate / 100)
    const subtotalWithOverhead = subtotal + overheadAmount
    const profitAmount = subtotalWithOverhead * (profitRate / 100)
    const totalCost = subtotalWithOverhead + profitAmount
    const unitPrice = inquiry ? totalCost / inquiry.quantity : 0

    return {
      subtotal,
      overheadAmount,
      profitAmount,
      totalCost,
      unitPrice,
    }
  }

  const totals = calculateTotals()

  const handleSubmit = () => {
    if (!inquiry) return

    const validUntil = addDays(new Date(), 30) // Quote valid for 30 days

    createQuoteMutation.mutate({
      inquiry_id: inquiryId,
      ...costs,
      overhead_rate: overheadRate,
      profit_rate: profitRate,
      currency: 'USD',
      valid_until: format(validUntil, 'yyyy-MM-dd'),
      delivery_days: deliveryDays,
      payment_terms: paymentTerms,
      notes,
    })
  }

  const handleProcessCostCalculate = () => {
    // Navigate to cost calculator with inquiry context
    router.push(`/cost-calculator?inquiry=${inquiryId}`)
  }

  if (!inquiry) {
    return (
      <DashboardLayout>
        <div className="flex items-center justify-center h-64">
          <div className="text-center">載入中...</div>
        </div>
      </DashboardLayout>
    )
  }

  return (
    <DashboardLayout>
      <div className="max-w-6xl mx-auto space-y-6">
        {/* Header */}
        <div className="flex items-center gap-4">
          <Button
            variant="ghost"
            size="icon"
            onClick={() => router.push('/inquiries')}
          >
            <ArrowLeft className="h-4 w-4" />
          </Button>
          <div>
            <h1 className="text-3xl font-bold text-gray-900">建立報價單</h1>
            <p className="mt-1 text-gray-600">
              詢價單號：{inquiry.inquiry_no} | 客戶：{inquiry.customer?.name}
            </p>
          </div>
        </div>

        {/* Product Info Alert */}
        <Alert>
          <AlertCircle className="h-4 w-4" />
          <AlertDescription>
            <strong>產品：</strong>{inquiry.product_name} | 
            <strong> 數量：</strong>{inquiry.quantity.toLocaleString()} {inquiry.unit} | 
            <strong> 交期：</strong>{format(new Date(inquiry.required_date), 'yyyy/MM/dd', { locale: zhTW })}
          </AlertDescription>
        </Alert>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Cost Input Section */}
          <div className="lg:col-span-2 space-y-6">
            {/* Material Cost */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Package className="h-5 w-5" />
                  材料成本
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="grid gap-4">
                  <div>
                    <Label htmlFor="material_cost">材料成本總額</Label>
                    <Input
                      id="material_cost"
                      type="number"
                      value={costs.material_cost}
                      onChange={(e) => setCosts({ ...costs, material_cost: parseFloat(e.target.value) || 0 })}
                      step="0.01"
                      min="0"
                    />
                    <p className="text-sm text-gray-500 mt-1">
                      包含原材料採購成本
                    </p>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Manufacturing Cost */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Wrench className="h-5 w-5" />
                  製造成本
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div>
                  <div className="flex justify-between items-center mb-2">
                    <Label htmlFor="process_cost">製程成本</Label>
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={handleProcessCostCalculate}
                    >
                      <Calculator className="mr-2 h-4 w-4" />
                      成本計算器
                    </Button>
                  </div>
                  <Input
                    id="process_cost"
                    type="number"
                    value={costs.process_cost}
                    onChange={(e) => setCosts({ ...costs, process_cost: parseFloat(e.target.value) || 0 })}
                    step="0.01"
                    min="0"
                  />
                </div>
                
                <div>
                  <Label htmlFor="surface_cost">表面處理成本</Label>
                  <Input
                    id="surface_cost"
                    type="number"
                    value={costs.surface_cost}
                    onChange={(e) => setCosts({ ...costs, surface_cost: parseFloat(e.target.value) || 0 })}
                    step="0.01"
                    min="0"
                  />
                </div>

                <div>
                  <Label htmlFor="heat_treat_cost">熱處理成本</Label>
                  <Input
                    id="heat_treat_cost"
                    type="number"
                    value={costs.heat_treat_cost}
                    onChange={(e) => setCosts({ ...costs, heat_treat_cost: parseFloat(e.target.value) || 0 })}
                    step="0.01"
                    min="0"
                  />
                </div>
              </CardContent>
            </Card>

            {/* Logistics & Trade Cost */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Truck className="h-5 w-5" />
                  物流與貿易成本
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div>
                  <Label htmlFor="packaging_cost">包裝成本</Label>
                  <Input
                    id="packaging_cost"
                    type="number"
                    value={costs.packaging_cost}
                    onChange={(e) => setCosts({ ...costs, packaging_cost: parseFloat(e.target.value) || 0 })}
                    step="0.01"
                    min="0"
                  />
                </div>

                <div>
                  <Label htmlFor="shipping_cost">運輸成本</Label>
                  <Input
                    id="shipping_cost"
                    type="number"
                    value={costs.shipping_cost}
                    onChange={(e) => setCosts({ ...costs, shipping_cost: parseFloat(e.target.value) || 0 })}
                    step="0.01"
                    min="0"
                  />
                  <p className="text-sm text-gray-500 mt-1">
                    {inquiry.incoterm} - {inquiry.destination_port || inquiry.destination_address}
                  </p>
                </div>

                <div>
                  <Label htmlFor="tariff_cost">關稅及相關費用</Label>
                  <div className="flex gap-2">
                    <Input
                      id="tariff_cost"
                      type="number"
                      value={costs.tariff_cost}
                      onChange={(e) => setCosts({ ...costs, tariff_cost: parseFloat(e.target.value) || 0 })}
                      step="0.01"
                      min="0"
                    />
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => router.push(`/tariffs?inquiry=${inquiryId}`)}
                    >
                      <Globe className="mr-2 h-4 w-4" />
                      計算關稅
                    </Button>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Quote Settings */}
            <Card>
              <CardHeader>
                <CardTitle>報價設定</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <Label htmlFor="overhead_rate">管理費率 (%)</Label>
                    <Input
                      id="overhead_rate"
                      type="number"
                      value={overheadRate}
                      onChange={(e) => setOverheadRate(parseFloat(e.target.value) || 0)}
                      step="0.1"
                      min="0"
                      max="100"
                    />
                  </div>
                  <div>
                    <Label htmlFor="profit_rate">利潤率 (%)</Label>
                    <Input
                      id="profit_rate"
                      type="number"
                      value={profitRate}
                      onChange={(e) => setProfitRate(parseFloat(e.target.value) || 0)}
                      step="0.1"
                      min="0"
                      max="100"
                    />
                  </div>
                </div>

                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <Label htmlFor="delivery_days">交貨天數</Label>
                    <Input
                      id="delivery_days"
                      type="number"
                      value={deliveryDays}
                      onChange={(e) => setDeliveryDays(parseInt(e.target.value) || 0)}
                      min="1"
                    />
                  </div>
                  <div>
                    <Label htmlFor="payment_terms">付款條件</Label>
                    <Input
                      id="payment_terms"
                      value={paymentTerms}
                      onChange={(e) => setPaymentTerms(e.target.value)}
                      placeholder="例如：T/T 30 days"
                    />
                  </div>
                </div>

                <div>
                  <Label htmlFor="notes">備註</Label>
                  <Textarea
                    id="notes"
                    value={notes}
                    onChange={(e) => setNotes(e.target.value)}
                    placeholder="輸入報價相關備註..."
                    rows={3}
                  />
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Cost Summary Section */}
          <div className="space-y-6">
            <Card className="sticky top-6">
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <DollarSign className="h-5 w-5" />
                  成本摘要
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {/* Cost Items */}
                  <div className="space-y-2 text-sm">
                    <div className="flex justify-between">
                      <span>材料成本</span>
                      <span>${costs.material_cost.toFixed(2)}</span>
                    </div>
                    <div className="flex justify-between">
                      <span>製程成本</span>
                      <span>${costs.process_cost.toFixed(2)}</span>
                    </div>
                    <div className="flex justify-between">
                      <span>表面處理</span>
                      <span>${costs.surface_cost.toFixed(2)}</span>
                    </div>
                    <div className="flex justify-between">
                      <span>熱處理</span>
                      <span>${costs.heat_treat_cost.toFixed(2)}</span>
                    </div>
                    <div className="flex justify-between">
                      <span>包裝成本</span>
                      <span>${costs.packaging_cost.toFixed(2)}</span>
                    </div>
                    <div className="flex justify-between">
                      <span>運輸成本</span>
                      <span>${costs.shipping_cost.toFixed(2)}</span>
                    </div>
                    <div className="flex justify-between">
                      <span>關稅費用</span>
                      <span>${costs.tariff_cost.toFixed(2)}</span>
                    </div>
                  </div>

                  <Separator />

                  {/* Subtotal */}
                  <div className="flex justify-between font-medium">
                    <span>小計</span>
                    <span>${totals.subtotal.toFixed(2)}</span>
                  </div>

                  {/* Overhead */}
                  <div className="flex justify-between text-sm">
                    <span>管理費 ({overheadRate}%)</span>
                    <span>${totals.overheadAmount.toFixed(2)}</span>
                  </div>

                  {/* Profit */}
                  <div className="flex justify-between text-sm">
                    <span>利潤 ({profitRate}%)</span>
                    <span>${totals.profitAmount.toFixed(2)}</span>
                  </div>

                  <Separator />

                  {/* Total */}
                  <div className="space-y-2">
                    <div className="flex justify-between text-lg font-bold">
                      <span>總成本</span>
                      <span className="text-green-600">
                        ${totals.totalCost.toFixed(2)}
                      </span>
                    </div>
                    <div className="flex justify-between">
                      <span>單位價格</span>
                      <span className="font-medium">
                        ${totals.unitPrice.toFixed(4)}
                      </span>
                    </div>
                  </div>
                </div>

                {/* Actions */}
                <div className="mt-6 space-y-2">
                  <Button
                    className="w-full"
                    onClick={handleSubmit}
                    disabled={createQuoteMutation.isPending}
                  >
                    {createQuoteMutation.isPending ? (
                      <>
                        <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                        建立中...
                      </>
                    ) : (
                      '建立報價單'
                    )}
                  </Button>
                  <Button
                    variant="outline"
                    className="w-full"
                    onClick={() => router.push('/inquiries')}
                  >
                    取消
                  </Button>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </DashboardLayout>
  )
}