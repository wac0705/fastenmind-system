'use client'

import { useState, useEffect } from 'react'
import { useRouter, useParams } from 'next/navigation'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import DashboardLayout from '@/components/layout/DashboardLayout'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Separator } from '@/components/ui/separator'
import { Badge } from '@/components/ui/badge'
import { useToast } from '@/components/ui/use-toast'
import { 
  ArrowLeft, 
  Save, 
  Calculator, 
  Plus, 
  Trash2, 
  AlertCircle,
  Package,
  DollarSign,
  FileText,
  User,
  Building2,
  Calendar,
  Truck,
  CreditCard
} from 'lucide-react'
import { format } from 'date-fns'
import { zhTW } from 'date-fns/locale'
import { useAuthStore } from '@/store/auth.store'
import quoteService from '@/services/quote.service'
import { customerService } from '@/services/customer.service'
import { Quote, UpdateQuoteRequest, QuoteItemRequest, QuoteTermRequest } from '@/types/quote'

export default function EditQuotePage() {
  const router = useRouter()
  const params = useParams()
  const queryClient = useQueryClient()
  const { toast } = useToast()
  const { user } = useAuthStore()
  const quoteId = params.id as string

  const [activeSection, setActiveSection] = useState('basic')
  const [newItem, setNewItem] = useState<QuoteItemRequest>({
    product_name: '',
    product_specs: '',
    quantity: 1,
    unit: 'PCS',
    unit_price: 0,
    notes: ''
  })
  const [newTerm, setNewTerm] = useState<QuoteTermRequest>({
    term_type: '付款條件',
    term_content: '',
    sort_order: 1
  })
  const [items, setItems] = useState<QuoteItemRequest[]>([])
  const [terms, setTerms] = useState<QuoteTermRequest[]>([])
  const [basicInfo, setBasicInfo] = useState({
    validity_days: 30,
    payment_terms: 'T/T 30 days',
    delivery_terms: 'FOB',
    remarks: ''
  })

  // Fetch quote details
  const { data: quote, isLoading } = useQuery({
    queryKey: ['quote', quoteId],
    queryFn: () => quoteService.getQuote(quoteId),
  })

  // Fetch quote versions for items and terms
  const { data: versions } = useQuery({
    queryKey: ['quote-versions', quoteId],
    queryFn: () => quoteService.getQuoteVersions(quoteId),
    enabled: !!quote,
  })

  // Fetch customers
  const { data: customers } = useQuery({
    queryKey: ['customers'],
    queryFn: () => customerService.list({ page_size: 100 }),
  })

  // Initialize form data when quote is loaded
  useEffect(() => {
    if (quote) {
      setBasicInfo({
        validity_days: quote.validity_days || 30,
        payment_terms: quote.payment_terms || 'T/T 30 days',
        delivery_terms: quote.delivery_terms || 'FOB',
        remarks: quote.remarks || ''
      })
    }
  }, [quote])

  // Initialize items and terms when versions are loaded
  useEffect(() => {
    if (versions && versions.length > 0 && versions[0].items) {
      const currentVersion = versions[0]
      setItems(currentVersion.items?.map(item => ({
        product_name: item.product_name,
        product_specs: item.product_specs || '',
        quantity: item.quantity,
        unit: item.unit,
        unit_price: item.unit_price,
        notes: item.notes || ''
      })) || [])
      
      if (currentVersion.terms) {
        setTerms(currentVersion.terms.map(term => ({
          term_type: term.term_type,
          term_content: term.term_content,
          sort_order: term.sort_order
        })))
      }
    }
  }, [versions])

  // Update quote mutation
  const updateQuoteMutation = useMutation({
    mutationFn: (data: UpdateQuoteRequest) => quoteService.updateQuote(quoteId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['quote', quoteId] })
      queryClient.invalidateQueries({ queryKey: ['quote-versions', quoteId] })
      toast({ title: '報價單更新成功' })
      router.push(`/quotes/${quoteId}`)
    },
    onError: (error: any) => {
      toast({
        title: '更新失敗',
        description: error.response?.data?.message || '更新報價單時發生錯誤',
        variant: 'destructive',
      })
    },
  })

  const handleAddItem = () => {
    if (!newItem.product_name || newItem.quantity <= 0 || newItem.unit_price <= 0) {
      toast({
        title: '請填寫完整的產品資訊',
        variant: 'destructive'
      })
      return
    }

    setItems(prev => [...prev, newItem])
    setNewItem({
      product_name: '',
      product_specs: '',
      quantity: 1,
      unit: 'PCS',
      unit_price: 0,
      notes: ''
    })
  }

  const handleRemoveItem = (index: number) => {
    setItems(prev => prev.filter((_, i) => i !== index))
  }

  const handleAddTerm = () => {
    if (!newTerm.term_content) {
      toast({
        title: '請填寫條款內容',
        variant: 'destructive'
      })
      return
    }

    setTerms(prev => [...prev, { ...newTerm, sort_order: prev.length + 1 }])
    setNewTerm({
      term_type: '付款條件',
      term_content: '',
      sort_order: 1
    })
  }

  const handleRemoveTerm = (index: number) => {
    setTerms(prev => prev.filter((_, i) => i !== index))
  }

  const handleSubmit = () => {
    if (items.length === 0) {
      toast({
        title: '請至少添加一個報價項目',
        variant: 'destructive'
      })
      return
    }

    const updateData: UpdateQuoteRequest = {
      create_new_version: true,
      version_notes: '編輯更新',
      validity_days: basicInfo.validity_days,
      payment_terms: basicInfo.payment_terms,
      delivery_terms: basicInfo.delivery_terms,
      remarks: basicInfo.remarks,
      items,
      terms
    }

    updateQuoteMutation.mutate(updateData)
  }

  const getTotalAmount = () => {
    return items.reduce((sum, item) => sum + item.quantity * item.unit_price, 0)
  }

  if (isLoading) {
    return (
      <DashboardLayout>
        <div className="flex items-center justify-center h-64">
          <div className="text-center">載入中...</div>
        </div>
      </DashboardLayout>
    )
  }

  if (!quote) {
    return (
      <DashboardLayout>
        <div className="flex items-center justify-center h-64">
          <div className="text-center">找不到報價單</div>
        </div>
      </DashboardLayout>
    )
  }

  // Check permissions
  const canEdit = (quote.status === 'draft' || quote.status === 'rejected') && 
                  (quote.created_by === user?.id || ['admin', 'manager'].includes(user?.role || ''))

  if (!canEdit) {
    return (
      <DashboardLayout>
        <div className="flex items-center justify-center h-64">
          <div className="text-center text-red-600">無權限編輯此報價單</div>
        </div>
      </DashboardLayout>
    )
  }

  return (
    <DashboardLayout>
      <div className="max-w-7xl mx-auto space-y-6">
        {/* Header */}
        <div className="flex justify-between items-start">
          <div className="flex items-center gap-4">
            <Button
              variant="ghost"
              size="icon"
              onClick={() => router.push(`/quotes/${quoteId}`)}
            >
              <ArrowLeft className="h-4 w-4" />
            </Button>
            <div>
              <h1 className="text-3xl font-bold text-gray-900">編輯報價單</h1>
              <p className="mt-1 text-gray-600">
                報價單號：{quote.quote_no} | 狀態：
                <Badge variant={quote.status === 'draft' ? 'secondary' : 'destructive'} className="ml-1">
                  {quote.status === 'draft' ? '草稿' : '已拒絕'}
                </Badge>
              </p>
            </div>
          </div>
          <div className="flex gap-2">
            <Button variant="outline" onClick={() => router.push(`/quotes/${quoteId}`)}>取消</Button>
            <Button onClick={handleSubmit} disabled={updateQuoteMutation.isPending}>
              <Save className="mr-2 h-4 w-4" />
              {updateQuoteMutation.isPending ? '更新中...' : '儲存變更'}
            </Button>
          </div>
        </div>

        {/* Navigation Tabs */}
        <div className="border-b border-gray-200">
          <nav className="-mb-px flex space-x-8">
            {[
              { id: 'basic', name: '基本資訊', icon: FileText },
              { id: 'items', name: '報價項目', icon: Package },
              { id: 'terms', name: '條款設定', icon: FileText },
              { id: 'summary', name: '摘要預覽', icon: Calculator },
            ].map((tab) => {
              const Icon = tab.icon
              return (
                <button
                  key={tab.id}
                  onClick={() => setActiveSection(tab.id)}
                  className={`${
                    activeSection === tab.id
                      ? 'border-blue-500 text-blue-600'
                      : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                  } whitespace-nowrap py-2 px-1 border-b-2 font-medium text-sm flex items-center gap-2`}
                >
                  <Icon className="h-4 w-4" />
                  {tab.name}
                </button>
              )
            })}
          </nav>
        </div>

        {/* Basic Info Section */}
        {activeSection === 'basic' && (
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <FileText className="h-5 w-5" />
                基本資訊
              </CardTitle>
              <CardDescription>編輯報價單的基本資訊</CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              {/* Customer Info (Read-only) */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <Label className="text-gray-500">客戶</Label>
                  <div className="flex items-center gap-2 mt-1 p-3 bg-gray-50 rounded-md">
                    <Building2 className="h-4 w-4 text-gray-400" />
                    <span className="font-medium">{quote.customer?.name}</span>
                  </div>
                </div>
                <div>
                  <Label className="text-gray-500">建立者</Label>
                  <div className="flex items-center gap-2 mt-1 p-3 bg-gray-50 rounded-md">
                    <User className="h-4 w-4 text-gray-400" />
                    <span>{quote.created_by_user?.name}</span>
                  </div>
                </div>
              </div>

              {/* Editable Fields */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <Label htmlFor="validity_days">有效天數</Label>
                  <div className="flex items-center gap-2 mt-1">
                    <Calendar className="h-4 w-4 text-gray-400" />
                    <Input
                      id="validity_days"
                      type="number"
                      value={basicInfo.validity_days}
                      onChange={(e) => setBasicInfo(prev => ({ ...prev, validity_days: parseInt(e.target.value) || 30 }))}
                      min={1}
                    />
                  </div>
                </div>
                <div>
                  <Label htmlFor="payment_terms">付款條件</Label>
                  <div className="flex items-center gap-2 mt-1">
                    <CreditCard className="h-4 w-4 text-gray-400" />
                    <Input
                      id="payment_terms"
                      value={basicInfo.payment_terms}
                      onChange={(e) => setBasicInfo(prev => ({ ...prev, payment_terms: e.target.value }))}
                    />
                  </div>
                </div>
                <div>
                  <Label htmlFor="delivery_terms">交貨條件</Label>
                  <div className="flex items-center gap-2 mt-1">
                    <Truck className="h-4 w-4 text-gray-400" />
                    <Input
                      id="delivery_terms"
                      value={basicInfo.delivery_terms}
                      onChange={(e) => setBasicInfo(prev => ({ ...prev, delivery_terms: e.target.value }))}
                    />
                  </div>
                </div>
              </div>

              <div>
                <Label htmlFor="remarks">備註</Label>
                <Textarea
                  id="remarks"
                  value={basicInfo.remarks}
                  onChange={(e) => setBasicInfo(prev => ({ ...prev, remarks: e.target.value }))}
                  rows={3}
                  className="mt-1"
                />
              </div>
            </CardContent>
          </Card>
        )}

        {/* Items Section */}
        {activeSection === 'items' && (
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Package className="h-5 w-5" />
                報價項目
              </CardTitle>
              <CardDescription>編輯報價項目和價格</CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              {/* Add New Item Form */}
              <div className="bg-gray-50 p-4 rounded-lg">
                <h4 className="font-medium mb-3">新增項目</h4>
                <div className="grid grid-cols-1 md:grid-cols-6 gap-3">
                  <div className="md:col-span-2">
                    <Input
                      placeholder="產品名稱"
                      value={newItem.product_name}
                      onChange={(e) => setNewItem(prev => ({ ...prev, product_name: e.target.value }))}
                    />
                  </div>
                  <div className="md:col-span-2">
                    <Input
                      placeholder="規格"
                      value={newItem.product_specs}
                      onChange={(e) => setNewItem(prev => ({ ...prev, product_specs: e.target.value }))}
                    />
                  </div>
                  <div>
                    <Input
                      type="number"
                      placeholder="數量"
                      value={newItem.quantity}
                      onChange={(e) => setNewItem(prev => ({ ...prev, quantity: parseInt(e.target.value) || 1 }))}
                      min={1}
                    />
                  </div>
                  <div>
                    <Input
                      placeholder="單位"
                      value={newItem.unit}
                      onChange={(e) => setNewItem(prev => ({ ...prev, unit: e.target.value }))}
                    />
                  </div>
                  <div>
                    <Input
                      type="number"
                      placeholder="單價"
                      value={newItem.unit_price}
                      onChange={(e) => setNewItem(prev => ({ ...prev, unit_price: parseFloat(e.target.value) || 0 }))}
                      step="0.0001"
                      min={0}
                    />
                  </div>
                  <div className="md:col-span-5">
                    <Input
                      placeholder="備註"
                      value={newItem.notes}
                      onChange={(e) => setNewItem(prev => ({ ...prev, notes: e.target.value }))}
                    />
                  </div>
                  <div>
                    <Button onClick={handleAddItem} className="w-full">
                      <Plus className="mr-2 h-4 w-4" />
                      新增
                    </Button>
                  </div>
                </div>
              </div>

              {/* Items List */}
              {items.length > 0 ? (
                <div className="overflow-x-auto">
                  <table className="min-w-full divide-y divide-gray-200">
                    <thead className="bg-gray-50">
                      <tr>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                          產品名稱
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                          規格
                        </th>
                        <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                          數量
                        </th>
                        <th className="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider">
                          單位
                        </th>
                        <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                          單價
                        </th>
                        <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                          總價
                        </th>
                        <th className="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider">
                          操作
                        </th>
                      </tr>
                    </thead>
                    <tbody className="bg-white divide-y divide-gray-200">
                      {items.map((item, index) => (
                        <tr key={index}>
                          <td className="px-6 py-4 text-sm text-gray-900">
                            {item.product_name}
                          </td>
                          <td className="px-6 py-4 text-sm text-gray-500">
                            {item.product_specs || '-'}
                          </td>
                          <td className="px-6 py-4 text-sm text-gray-900 text-right">
                            {item.quantity.toLocaleString()}
                          </td>
                          <td className="px-6 py-4 text-sm text-gray-900 text-center">
                            {item.unit}
                          </td>
                          <td className="px-6 py-4 text-sm text-gray-900 text-right">
                            ${item.unit_price.toFixed(4)}
                          </td>
                          <td className="px-6 py-4 text-sm font-medium text-gray-900 text-right">
                            ${(item.quantity * item.unit_price).toFixed(2)}
                          </td>
                          <td className="px-6 py-4 text-center">
                            <Button
                              variant="ghost"
                              size="icon"
                              onClick={() => handleRemoveItem(index)}
                              className="text-red-600 hover:text-red-900"
                            >
                              <Trash2 className="h-4 w-4" />
                            </Button>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                    <tfoot className="bg-gray-50">
                      <tr>
                        <td colSpan={5} className="px-6 py-4 text-right text-sm font-medium text-gray-900">
                          總計
                        </td>
                        <td className="px-6 py-4 text-right text-sm font-bold text-gray-900">
                          ${getTotalAmount().toFixed(2)}
                        </td>
                        <td></td>
                      </tr>
                    </tfoot>
                  </table>
                </div>
              ) : (
                <div className="text-center py-8 text-gray-500">
                  尚未添加報價項目
                </div>
              )}
            </CardContent>
          </Card>
        )}

        {/* Terms Section */}
        {activeSection === 'terms' && (
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <FileText className="h-5 w-5" />
                條款設定
              </CardTitle>
              <CardDescription>設定報價條款和條件</CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              {/* Add New Term Form */}
              <div className="bg-gray-50 p-4 rounded-lg">
                <h4 className="font-medium mb-3">新增條款</h4>
                <div className="grid grid-cols-1 md:grid-cols-4 gap-3">
                  <div>
                    <select
                      value={newTerm.term_type}
                      onChange={(e) => setNewTerm(prev => ({ ...prev, term_type: e.target.value }))}
                      className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                    >
                      <option value="付款條件">付款條件</option>
                      <option value="交貨條件">交貨條件</option>
                      <option value="品質要求">品質要求</option>
                      <option value="保固條款">保固條款</option>
                      <option value="其他條款">其他條款</option>
                    </select>
                  </div>
                  <div className="md:col-span-2">
                    <Textarea
                      placeholder="條款內容"
                      value={newTerm.term_content}
                      onChange={(e) => setNewTerm(prev => ({ ...prev, term_content: e.target.value }))}
                      rows={2}
                    />
                  </div>
                  <div>
                    <Button onClick={handleAddTerm} className="w-full h-full">
                      <Plus className="mr-2 h-4 w-4" />
                      新增
                    </Button>
                  </div>
                </div>
              </div>

              {/* Terms List */}
              {terms.length > 0 ? (
                <div className="space-y-3">
                  {terms.map((term, index) => (
                    <div key={index} className="border rounded-lg p-4 flex justify-between items-start">
                      <div className="flex-1">
                        <h5 className="font-medium text-sm text-gray-700 mb-1">{term.term_type}</h5>
                        <p className="text-sm text-gray-600 whitespace-pre-wrap">{term.term_content}</p>
                      </div>
                      <Button
                        variant="ghost"
                        size="icon"
                        onClick={() => handleRemoveTerm(index)}
                        className="ml-4 text-red-600 hover:text-red-900"
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    </div>
                  ))}
                </div>
              ) : (
                <div className="text-center py-8 text-gray-500">
                  尚未添加條款
                </div>
              )}
            </CardContent>
          </Card>
        )}

        {/* Summary Section */}
        {activeSection === 'summary' && (
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Calculator className="h-5 w-5" />
                摘要預覽
              </CardTitle>
              <CardDescription>檢視報價單完整內容</CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              {/* Basic Info Summary */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div className="space-y-4">
                  <div>
                    <Label className="text-gray-500">客戶</Label>
                    <p className="font-medium">{quote.customer?.name}</p>
                  </div>
                  <div>
                    <Label className="text-gray-500">有效期限</Label>
                    <p>{basicInfo.validity_days} 天</p>
                  </div>
                  <div>
                    <Label className="text-gray-500">付款條件</Label>
                    <p>{basicInfo.payment_terms}</p>
                  </div>
                </div>
                <div className="space-y-4">
                  <div>
                    <Label className="text-gray-500">交貨條件</Label>
                    <p>{basicInfo.delivery_terms}</p>
                  </div>
                  <div>
                    <Label className="text-gray-500">總金額</Label>
                    <p className="text-2xl font-bold text-green-600">
                      ${getTotalAmount().toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
                    </p>
                  </div>
                </div>
              </div>

              <Separator />

              {/* Items Summary */}
              <div>
                <h4 className="font-medium mb-3">報價項目 ({items.length})</h4>
                {items.length > 0 ? (
                  <div className="space-y-2">
                    {items.map((item, index) => (
                      <div key={index} className="flex justify-between items-center p-3 bg-gray-50 rounded-md">
                        <div>
                          <p className="font-medium">{item.product_name}</p>
                          <p className="text-sm text-gray-600">
                            {item.quantity.toLocaleString()} {item.unit} × ${item.unit_price.toFixed(4)}
                          </p>
                        </div>
                        <p className="font-medium">${(item.quantity * item.unit_price).toFixed(2)}</p>
                      </div>
                    ))}
                  </div>
                ) : (
                  <p className="text-gray-500">無項目</p>
                )}
              </div>

              <Separator />

              {/* Terms Summary */}
              <div>
                <h4 className="font-medium mb-3">條款設定 ({terms.length})</h4>
                {terms.length > 0 ? (
                  <div className="space-y-3">
                    {terms.map((term, index) => (
                      <div key={index} className="p-3 bg-gray-50 rounded-md">
                        <h5 className="font-medium text-sm text-gray-700 mb-1">{term.term_type}</h5>
                        <p className="text-sm text-gray-600 whitespace-pre-wrap">{term.term_content}</p>
                      </div>
                    ))}
                  </div>
                ) : (
                  <p className="text-gray-500">無條款</p>
                )}
              </div>

              {/* Action Buttons */}
              <div className="pt-6 border-t">
                <div className="flex justify-end gap-4">
                  <Button variant="outline" onClick={() => router.push(`/quotes/${quoteId}`)}>取消</Button>
                  <Button onClick={handleSubmit} disabled={updateQuoteMutation.isPending}>
                    <Save className="mr-2 h-4 w-4" />
                    {updateQuoteMutation.isPending ? '更新中...' : '儲存變更'}
                  </Button>
                </div>
              </div>
            </CardContent>
          </Card>
        )}
      </div>
    </DashboardLayout>
  )
}