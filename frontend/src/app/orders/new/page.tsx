'use client'

import { useState, useEffect } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import { useQuery, useMutation } from '@tanstack/react-query'
import DashboardLayout from '@/components/layout/DashboardLayout'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { ArrowLeft, Save, Package, FileText, Calendar, DollarSign, Truck, CreditCard } from 'lucide-react'
import { useToast } from '@/components/ui/use-toast'
import { format } from 'date-fns'
import { zhTW } from 'date-fns/locale'
import orderService, { CreateOrderRequest } from '@/services/order.service'
import quoteService from '@/services/quote.service'
import { Separator } from '@/components/ui/separator'

export default function NewOrderPage() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const { toast } = useToast()
  const quoteId = searchParams.get('quote_id')

  const [formData, setFormData] = useState<CreateOrderRequest>({
    quote_id: quoteId || '',
    po_number: '',
    quantity: 0,
    delivery_method: 'FOB',
    delivery_date: format(new Date(Date.now() + 30 * 24 * 60 * 60 * 1000), 'yyyy-MM-dd'),
    shipping_address: '',
    payment_terms: 'NET30',
    down_payment: 0,
    notes: '',
  })

  // Fetch quote details if quote_id is provided
  const { data: quote } = useQuery({
    queryKey: ['quote', quoteId],
    queryFn: () => quoteService.getQuote(quoteId as string),
    enabled: !!quoteId,
  })

  // Fetch all approved quotes for selection
  const { data: availableQuotes } = useQuery({
    queryKey: ['available-quotes'],
    queryFn: () => quoteService.getQuotes({ status: 'approved', page_size: 100 }),
    enabled: !quoteId,
  })

  useEffect(() => {
    if (quote) {
      setFormData(prev => ({
        ...prev,
        quantity: quote.inquiry?.quantity || 0,
        payment_terms: quote.payment_terms || 'NET30',
        shipping_address: quote.customer?.address || '',
      }))
    }
  }, [quote])

  // Create order mutation
  const createOrderMutation = useMutation({
    mutationFn: (data: CreateOrderRequest) => orderService.createFromQuote(data),
    onSuccess: (order) => {
      toast({
        title: '成功',
        description: '訂單已建立',
      })
      router.push(`/orders/${order.id}`)
    },
    onError: (error: any) => {
      toast({
        title: '錯誤',
        description: error.response?.data?.message || '建立訂單時發生錯誤',
        variant: 'destructive',
      })
    },
  })

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target
    setFormData(prev => ({
      ...prev,
      [name]: name === 'quantity' || name === 'down_payment' ? parseFloat(value) || 0 : value,
    }))
  }

  const handleSelectChange = (name: string, value: string) => {
    setFormData(prev => ({
      ...prev,
      [name]: value,
    }))
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!formData.quote_id) {
      toast({
        title: '錯誤',
        description: '請選擇報價單',
        variant: 'destructive',
      })
      return
    }

    if (!formData.po_number) {
      toast({
        title: '錯誤',
        description: '請輸入 PO 號碼',
        variant: 'destructive',
      })
      return
    }

    createOrderMutation.mutate(formData)
  }

  const calculateTotal = () => {
    if (quote && formData.quantity) {
      return quote.unit_price * formData.quantity
    }
    return 0
  }

  return (
    <DashboardLayout>
      <form onSubmit={handleSubmit} className="max-w-4xl mx-auto">
        <div className="flex items-center gap-4 mb-6">
          <Button
            type="button"
            variant="ghost"
            size="icon"
            onClick={() => router.back()}
          >
            <ArrowLeft className="h-4 w-4" />
          </Button>
          <div>
            <h1 className="text-3xl font-bold text-gray-900">建立新訂單</h1>
            <p className="mt-1 text-gray-600">從報價單建立新的訂單</p>
          </div>
        </div>

        <div className="space-y-6">
          {/* Quote Selection */}
          {!quoteId && (
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <FileText className="h-5 w-5" />
                  選擇報價單
                </CardTitle>
                <CardDescription>
                  請選擇一個已核准的報價單來建立訂單
                </CardDescription>
              </CardHeader>
              <CardContent>
                <Select
                  value={formData.quote_id}
                  onValueChange={(value) => setFormData(prev => ({ ...prev, quote_id: value }))}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="選擇報價單" />
                  </SelectTrigger>
                  <SelectContent>
                    {availableQuotes?.data.map((quote) => (
                      <SelectItem key={quote.id} value={quote.id}>
                        {quote.quote_no} - {quote.customer?.name} ({quote.currency} {quote.total_amount.toFixed(2)})
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </CardContent>
            </Card>
          )}

          {/* Quote Info */}
          {quote && (
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <FileText className="h-5 w-5" />
                  報價單資訊
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <Label className="text-gray-500">報價單號</Label>
                    <p className="font-medium">{quote.quote_no}</p>
                  </div>
                  <div>
                    <Label className="text-gray-500">客戶</Label>
                    <p className="font-medium">{quote.customer?.name}</p>
                  </div>
                  <div>
                    <Label className="text-gray-500">料號</Label>
                    <p className="font-medium">{quote.inquiry?.part_no}</p>
                  </div>
                  <div>
                    <Label className="text-gray-500">單價</Label>
                    <p className="font-medium">{quote.currency} {quote.unit_price.toFixed(4)}</p>
                  </div>
                </div>
              </CardContent>
            </Card>
          )}

          {/* Order Details */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Package className="h-5 w-5" />
                訂單資訊
              </CardTitle>
              <CardDescription>
                填寫訂單的基本資訊
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="po_number">PO 號碼 *</Label>
                  <Input
                    id="po_number"
                    name="po_number"
                    value={formData.po_number}
                    onChange={handleInputChange}
                    placeholder="輸入客戶 PO 號碼"
                    required
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="quantity">訂購數量 *</Label>
                  <Input
                    id="quantity"
                    name="quantity"
                    type="number"
                    value={formData.quantity}
                    onChange={handleInputChange}
                    placeholder="輸入訂購數量"
                    min="1"
                    required
                  />
                </div>
              </div>

              {quote && formData.quantity > 0 && (
                <div className="bg-gray-50 p-4 rounded-lg">
                  <div className="flex justify-between items-center">
                    <span className="text-gray-600">訂單總額</span>
                    <span className="text-xl font-semibold">
                      {quote.currency} {calculateTotal().toLocaleString()}
                    </span>
                  </div>
                </div>
              )}
            </CardContent>
          </Card>

          {/* Delivery Info */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Truck className="h-5 w-5" />
                交貨資訊
              </CardTitle>
              <CardDescription>
                設定交貨方式和地址
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="delivery_method">交貨方式 *</Label>
                  <Select
                    value={formData.delivery_method}
                    onValueChange={(value) => handleSelectChange('delivery_method', value)}
                  >
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="EXW">EXW - 工廠交貨</SelectItem>
                      <SelectItem value="FOB">FOB - 裝運港交貨</SelectItem>
                      <SelectItem value="CIF">CIF - 成本加保險費加運費</SelectItem>
                      <SelectItem value="DDP">DDP - 完稅後交貨</SelectItem>
                      <SelectItem value="FCA">FCA - 貨交承運人</SelectItem>
                      <SelectItem value="CPT">CPT - 運費付至</SelectItem>
                      <SelectItem value="CIP">CIP - 運費保險費付至</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                <div className="space-y-2">
                  <Label htmlFor="delivery_date">交貨日期 *</Label>
                  <Input
                    id="delivery_date"
                    name="delivery_date"
                    type="date"
                    value={formData.delivery_date}
                    onChange={handleInputChange}
                    required
                  />
                </div>
              </div>
              <div className="space-y-2">
                <Label htmlFor="shipping_address">配送地址 *</Label>
                <Textarea
                  id="shipping_address"
                  name="shipping_address"
                  value={formData.shipping_address}
                  onChange={handleInputChange}
                  placeholder="輸入完整的配送地址"
                  rows={3}
                  required
                />
              </div>
            </CardContent>
          </Card>

          {/* Payment Info */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <CreditCard className="h-5 w-5" />
                付款資訊
              </CardTitle>
              <CardDescription>
                設定付款條件
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="payment_terms">付款條件 *</Label>
                  <Select
                    value={formData.payment_terms}
                    onValueChange={(value) => handleSelectChange('payment_terms', value)}
                  >
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="COD">貨到付款</SelectItem>
                      <SelectItem value="NET30">月結 30 天</SelectItem>
                      <SelectItem value="NET60">月結 60 天</SelectItem>
                      <SelectItem value="NET90">月結 90 天</SelectItem>
                      <SelectItem value="T/T">電匯</SelectItem>
                      <SelectItem value="L/C">信用狀</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                <div className="space-y-2">
                  <Label htmlFor="down_payment">頭期款金額</Label>
                  <Input
                    id="down_payment"
                    name="down_payment"
                    type="number"
                    value={formData.down_payment}
                    onChange={handleInputChange}
                    placeholder="輸入頭期款金額"
                    min="0"
                    step="0.01"
                  />
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Notes */}
          <Card>
            <CardHeader>
              <CardTitle>備註</CardTitle>
            </CardHeader>
            <CardContent>
              <Textarea
                name="notes"
                value={formData.notes}
                onChange={handleInputChange}
                placeholder="輸入任何額外的備註或特殊要求..."
                rows={4}
              />
            </CardContent>
          </Card>

          {/* Actions */}
          <div className="flex justify-end gap-4">
            <Button
              type="button"
              variant="outline"
              onClick={() => router.back()}
            >
              取消
            </Button>
            <Button
              type="submit"
              disabled={createOrderMutation.isPending}
            >
              <Save className="mr-2 h-4 w-4" />
              {createOrderMutation.isPending ? '建立中...' : '建立訂單'}
            </Button>
          </div>
        </div>
      </form>
    </DashboardLayout>
  )
}