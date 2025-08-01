'use client'

import { useState, useEffect } from 'react'
import { useParams, useRouter } from 'next/navigation'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
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
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Separator } from '@/components/ui/separator'
import { 
  ArrowLeft,
  Save, 
  Package,
  Calendar,
  DollarSign,
  Truck,
  CreditCard,
  Plus,
  Trash2,
  FileText,
  Building2,
  User,
  AlertCircle
} from 'lucide-react'
import { format } from 'date-fns'
import { useToast } from '@/components/ui/use-toast'
import orderService, { UpdateOrderRequest, OrderItemRequest } from '@/services/order.service'
import { useAuthStore } from '@/store/auth.store'

export default function EditOrderPage() {
  const params = useParams()
  const router = useRouter()
  const queryClient = useQueryClient()
  const { toast } = useToast()
  const { user } = useAuthStore()
  const orderId = params.id as string

  const [activeTab, setActiveTab] = useState('basic')
  const [orderData, setOrderData] = useState<UpdateOrderRequest>({
    po_number: '',
    delivery_method: '',
    delivery_date: '',
    shipping_address: '',
    payment_terms: '',
    notes: '',
    internal_notes: '',
  })
  const [items, setItems] = useState<OrderItemRequest[]>([])
  const [newItem, setNewItem] = useState<OrderItemRequest>({
    part_no: '',
    description: '',
    quantity: 1,
    unit_price: 0,
    material: '',
    surface_treatment: '',
    heat_treatment: '',
    specifications: '',
  })

  // Fetch order details
  const { data: order, isLoading } = useQuery({
    queryKey: ['order', orderId],
    queryFn: () => orderService.get(orderId),
  })

  // Fetch order items
  const { data: orderItems } = useQuery({
    queryKey: ['order-items', orderId],
    queryFn: () => orderService.getItems(orderId),
    enabled: !!order,
  })

  // Initialize form data when order is loaded
  useEffect(() => {
    if (order) {
      setOrderData({
        po_number: order.po_number,
        delivery_method: order.delivery_method,
        delivery_date: format(new Date(order.delivery_date), 'yyyy-MM-dd'),
        shipping_address: order.shipping_address,
        payment_terms: order.payment_terms,
        notes: order.notes || '',
        internal_notes: order.internal_notes || '',
      })
    }
  }, [order])

  // Initialize items when order items are loaded
  useEffect(() => {
    if (orderItems) {
      setItems(orderItems.map(item => ({
        part_no: item.part_no,
        description: item.description || '',
        quantity: item.quantity,
        unit_price: item.unit_price,
        material: item.material || '',
        surface_treatment: item.surface_treatment || '',
        heat_treatment: item.heat_treatment || '',
        specifications: item.specifications || '',
      })))
    }
  }, [orderItems])

  // Update order mutation
  const updateOrderMutation = useMutation({
    mutationFn: (data: UpdateOrderRequest) => orderService.update(orderId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['order', orderId] })
      toast({ title: '訂單更新成功' })
      router.push(`/orders/${orderId}`)
    },
    onError: (error: any) => {
      toast({
        title: '更新失敗',
        description: error.response?.data?.message || '更新訂單時發生錯誤',
        variant: 'destructive',
      })
    },
  })

  // Update items mutation
  const updateItemsMutation = useMutation({
    mutationFn: (items: OrderItemRequest[]) => orderService.updateItems(orderId, items),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['order-items', orderId] })
      toast({ title: '產品明細更新成功' })
    },
    onError: (error: any) => {
      toast({
        title: '更新失敗',
        description: error.response?.data?.message || '更新產品明細時發生錯誤',
        variant: 'destructive',
      })
    },
  })

  const handleInputChange = (field: keyof UpdateOrderRequest, value: string) => {
    setOrderData(prev => ({
      ...prev,
      [field]: value,
    }))
  }

  const handleAddItem = () => {
    if (!newItem.part_no || newItem.quantity <= 0 || newItem.unit_price <= 0) {
      toast({
        title: '請填寫完整的產品資訊',
        variant: 'destructive'
      })
      return
    }

    setItems(prev => [...prev, newItem])
    setNewItem({
      part_no: '',
      description: '',
      quantity: 1,
      unit_price: 0,
      material: '',
      surface_treatment: '',
      heat_treatment: '',
      specifications: '',
    })
  }

  const handleRemoveItem = (index: number) => {
    setItems(prev => prev.filter((_, i) => i !== index))
  }

  const handleUpdateItem = (index: number, field: keyof OrderItemRequest, value: string | number) => {
    setItems(prev => prev.map((item, i) => 
      i === index ? { ...item, [field]: value } : item
    ))
  }

  const handleSubmit = () => {
    updateOrderMutation.mutate(orderData)
  }

  const handleSaveItems = () => {
    updateItemsMutation.mutate(items)
  }

  const getTotalAmount = () => {
    return items.reduce((sum, item) => sum + item.quantity * item.unit_price, 0)
  }

  const canEdit = ['admin', 'manager', 'sales'].includes(user?.role || '') && 
                 ['pending', 'confirmed'].includes(order?.status || '')

  if (isLoading) {
    return (
      <DashboardLayout>
        <div className="flex items-center justify-center h-64">
          <div className="text-center">載入中...</div>
        </div>
      </DashboardLayout>
    )
  }

  if (!order) {
    return (
      <DashboardLayout>
        <div className="flex items-center justify-center h-64">
          <div className="text-center">找不到訂單</div>
        </div>
      </DashboardLayout>
    )
  }

  if (!canEdit) {
    return (
      <DashboardLayout>
        <div className="flex items-center justify-center h-64">
          <div className="text-center text-red-600">無權限編輯此訂單</div>
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
              onClick={() => router.push(`/orders/${orderId}`)}
            >
              <ArrowLeft className="h-4 w-4" />
            </Button>
            <div>
              <h1 className="text-3xl font-bold text-gray-900">編輯訂單</h1>
              <p className="mt-1 text-gray-600">
                訂單號：{order.order_no} | PO 號碼：{order.po_number}
              </p>
            </div>
          </div>
          <div className="flex gap-2">
            <Button variant="outline" onClick={() => router.push(`/orders/${orderId}`)}>取消</Button>
            <Button onClick={handleSubmit} disabled={updateOrderMutation.isPending}>
              <Save className="mr-2 h-4 w-4" />
              {updateOrderMutation.isPending ? '更新中...' : '儲存變更'}
            </Button>
          </div>
        </div>

        {/* Navigation Tabs */}
        <div className="border-b border-gray-200">
          <nav className="-mb-px flex space-x-8">
            {[
              { id: 'basic', name: '基本資訊', icon: FileText },
              { id: 'items', name: '產品明細', icon: Package },
              { id: 'shipping', name: '配送資訊', icon: Truck },
              { id: 'payment', name: '付款資訊', icon: CreditCard },
            ].map((tab) => {
              const Icon = tab.icon
              return (
                <button
                  key={tab.id}
                  onClick={() => setActiveTab(tab.id)}
                  className={`${
                    activeTab === tab.id
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

        {/* Basic Info Tab */}
        {activeTab === 'basic' && (
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <FileText className="h-5 w-5" />
                基本資訊
              </CardTitle>
              <CardDescription>編輯訂單的基本資訊</CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              {/* Read-only info */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <Label className="text-gray-500">客戶</Label>
                  <div className="flex items-center gap-2 mt-1 p-3 bg-gray-50 rounded-md">
                    <Building2 className="h-4 w-4 text-gray-400" />
                    <span className="font-medium">{order.customer?.name}</span>
                  </div>
                </div>
                <div>
                  <Label className="text-gray-500">業務人員</Label>
                  <div className="flex items-center gap-2 mt-1 p-3 bg-gray-50 rounded-md">
                    <User className="h-4 w-4 text-gray-400" />
                    <span>{order.sales?.full_name}</span>
                  </div>
                </div>
              </div>

              {/* Editable fields */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <Label htmlFor="po_number">PO 號碼</Label>
                  <Input
                    id="po_number"
                    value={orderData.po_number}
                    onChange={(e) => handleInputChange('po_number', e.target.value)}
                    className="mt-1"
                  />
                </div>
                <div>
                  <Label htmlFor="delivery_date">交貨日期</Label>
                  <Input
                    id="delivery_date"
                    type="date"
                    value={orderData.delivery_date}
                    onChange={(e) => handleInputChange('delivery_date', e.target.value)}
                    className="mt-1"
                  />
                </div>
              </div>

              <div>
                <Label htmlFor="notes">備註</Label>
                <Textarea
                  id="notes"
                  value={orderData.notes}
                  onChange={(e) => handleInputChange('notes', e.target.value)}
                  rows={3}
                  className="mt-1"
                />
              </div>

              <div>
                <Label htmlFor="internal_notes">內部備註</Label>
                <Textarea
                  id="internal_notes"
                  value={orderData.internal_notes}
                  onChange={(e) => handleInputChange('internal_notes', e.target.value)}
                  rows={3}
                  className="mt-1"
                />
              </div>
            </CardContent>
          </Card>
        )}

        {/* Items Tab */}
        {activeTab === 'items' && (
          <Card>
            <CardHeader className="flex flex-row items-center justify-between">
              <div>
                <CardTitle className="flex items-center gap-2">
                  <Package className="h-5 w-5" />
                  產品明細
                </CardTitle>
                <CardDescription>編輯訂單產品項目</CardDescription>
              </div>
              <Button onClick={handleSaveItems} disabled={updateItemsMutation.isPending}>
                <Save className="mr-2 h-4 w-4" />
                儲存明細
              </Button>
            </CardHeader>
            <CardContent className="space-y-6">
              {/* Add New Item Form */}
              <div className="bg-gray-50 p-4 rounded-lg">
                <h4 className="font-medium mb-3">新增產品項目</h4>
                <div className="grid grid-cols-1 md:grid-cols-8 gap-3">
                  <div className="md:col-span-2">
                    <Input
                      placeholder="料號"
                      value={newItem.part_no}
                      onChange={(e) => setNewItem(prev => ({ ...prev, part_no: e.target.value }))}
                    />
                  </div>
                  <div className="md:col-span-2">
                    <Input
                      placeholder="描述"
                      value={newItem.description}
                      onChange={(e) => setNewItem(prev => ({ ...prev, description: e.target.value }))}
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
                      type="number"
                      placeholder="單價"
                      value={newItem.unit_price}
                      onChange={(e) => setNewItem(prev => ({ ...prev, unit_price: parseFloat(e.target.value) || 0 }))}
                      step="0.0001"
                      min={0}
                    />
                  </div>
                  <div>
                    <Input
                      placeholder="材質"
                      value={newItem.material}
                      onChange={(e) => setNewItem(prev => ({ ...prev, material: e.target.value }))}
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
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">料號</th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">描述</th>
                        <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">數量</th>
                        <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">單價</th>
                        <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">總價</th>
                        <th className="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase">操作</th>
                      </tr>
                    </thead>
                    <tbody className="bg-white divide-y divide-gray-200">
                      {items.map((item, index) => (
                        <tr key={index}>
                          <td className="px-6 py-4">
                            <Input
                              value={item.part_no}
                              onChange={(e) => handleUpdateItem(index, 'part_no', e.target.value)}
                              className="text-sm"
                            />
                          </td>
                          <td className="px-6 py-4">
                            <Input
                              value={item.description}
                              onChange={(e) => handleUpdateItem(index, 'description', e.target.value)}
                              className="text-sm"
                            />
                          </td>
                          <td className="px-6 py-4">
                            <Input
                              type="number"
                              value={item.quantity}
                              onChange={(e) => handleUpdateItem(index, 'quantity', parseInt(e.target.value) || 1)}
                              className="text-sm text-right"
                              min={1}
                            />
                          </td>
                          <td className="px-6 py-4">
                            <Input
                              type="number"
                              value={item.unit_price}
                              onChange={(e) => handleUpdateItem(index, 'unit_price', parseFloat(e.target.value) || 0)}
                              className="text-sm text-right"
                              step="0.0001"
                              min={0}
                            />
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
                        <td colSpan={4} className="px-6 py-4 text-right text-sm font-medium text-gray-900">
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
                  尚未添加產品項目
                </div>
              )}
            </CardContent>
          </Card>
        )}

        {/* Shipping Tab */}
        {activeTab === 'shipping' && (
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Truck className="h-5 w-5" />
                配送資訊
              </CardTitle>
              <CardDescription>編輯配送方式和地址</CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              <div>
                <Label htmlFor="delivery_method">交貨方式</Label>
                <Select
                  value={orderData.delivery_method}
                  onValueChange={(value) => handleInputChange('delivery_method', value)}
                >
                  <SelectTrigger className="mt-1">
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

              <div>
                <Label htmlFor="shipping_address">配送地址</Label>
                <Textarea
                  id="shipping_address"
                  value={orderData.shipping_address}
                  onChange={(e) => handleInputChange('shipping_address', e.target.value)}
                  rows={4}
                  className="mt-1"
                  placeholder="輸入完整的配送地址"
                />
              </div>
            </CardContent>
          </Card>
        )}

        {/* Payment Tab */}
        {activeTab === 'payment' && (
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <CreditCard className="h-5 w-5" />
                付款資訊
              </CardTitle>
              <CardDescription>編輯付款條件</CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              <div>
                <Label htmlFor="payment_terms">付款條件</Label>
                <Select
                  value={orderData.payment_terms}
                  onValueChange={(value) => handleInputChange('payment_terms', value)}
                >
                  <SelectTrigger className="mt-1">
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

              <div className="bg-blue-50 p-4 rounded-lg">
                <div className="flex items-center gap-2 mb-2">
                  <AlertCircle className="h-4 w-4 text-blue-600" />
                  <span className="font-medium text-blue-900">付款資訊</span>
                </div>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
                  <div>
                    <span className="text-gray-600">總金額：</span>
                    <span className="font-medium">{order.currency} {order.total_amount.toLocaleString()}</span>
                  </div>
                  <div>
                    <span className="text-gray-600">頭期款：</span>
                    <span className="font-medium">{order.currency} {order.down_payment.toLocaleString()}</span>
                  </div>
                  <div>
                    <span className="text-gray-600">已付金額：</span>
                    <span className="font-medium">{order.currency} {order.paid_amount.toLocaleString()}</span>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        )}
      </div>
    </DashboardLayout>
  )
}