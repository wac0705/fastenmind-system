'use client'

import { useState } from 'react'
import { useParams, useRouter } from 'next/navigation'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import DashboardLayout from '@/components/layout/DashboardLayout'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Separator } from '@/components/ui/separator'
import { Textarea } from '@/components/ui/textarea'
import { Label } from '@/components/ui/label'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { 
  ArrowLeft,
  Package, 
  Calendar,
  Building2,
  User,
  CheckCircle,
  XCircle,
  Clock,
  Truck,
  Factory,
  AlertCircle,
  DollarSign,
  FileText,
  Edit,
  History,
  Upload,
  Download,
  Trash2,
  CreditCard,
  MapPin
} from 'lucide-react'
import { format } from 'date-fns'
import { zhTW } from 'date-fns/locale'
import { useToast } from '@/components/ui/use-toast'
import orderService, { Order, OrderItem, OrderDocument, OrderActivity } from '@/services/order.service'
import { useAuthStore } from '@/store/auth.store'

export default function OrderDetailPage() {
  const params = useParams()
  const router = useRouter()
  const queryClient = useQueryClient()
  const { toast } = useToast()
  const { user } = useAuthStore()
  const orderId = params.id as string

  const [activeTab, setActiveTab] = useState('details')
  const [isStatusDialogOpen, setIsStatusDialogOpen] = useState(false)
  const [selectedStatus, setSelectedStatus] = useState('')
  const [statusNotes, setStatusNotes] = useState('')
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false)

  // Fetch order details
  const { data: order, isLoading } = useQuery({
    queryKey: ['order', orderId],
    queryFn: () => orderService.get(orderId),
  })

  // Fetch order items
  const { data: items } = useQuery({
    queryKey: ['order-items', orderId],
    queryFn: () => orderService.getItems(orderId),
    enabled: !!order,
  })

  // Fetch documents
  const { data: documents } = useQuery({
    queryKey: ['order-documents', orderId],
    queryFn: () => orderService.getDocuments(orderId),
    enabled: !!order,
  })

  // Fetch activities
  const { data: activities } = useQuery({
    queryKey: ['order-activities', orderId],
    queryFn: () => orderService.getActivities(orderId),
    enabled: !!order,
  })

  // Update status mutation
  const updateStatusMutation = useMutation({
    mutationFn: (data: { status: string; notes: string }) => 
      orderService.updateStatus(orderId, data.status, data.notes),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['order', orderId] })
      queryClient.invalidateQueries({ queryKey: ['order-activities', orderId] })
      toast({ title: '訂單狀態已更新' })
      setIsStatusDialogOpen(false)
      setStatusNotes('')
    },
    onError: (error: any) => {
      toast({
        title: '更新失敗',
        description: error.response?.data?.message || '更新訂單狀態時發生錯誤',
        variant: 'destructive',
      })
    },
  })

  const getStatusBadge = (status: string) => {
    const statusConfig: Record<string, { label: string; variant: any; icon: any }> = {
      pending: { label: '待確認', variant: 'secondary', icon: Clock },
      confirmed: { label: '已確認', variant: 'info', icon: CheckCircle },
      in_production: { label: '生產中', variant: 'warning', icon: Factory },
      quality_check: { label: '品檢中', variant: 'warning', icon: AlertCircle },
      ready_to_ship: { label: '待出貨', variant: 'info', icon: Package },
      shipped: { label: '已出貨', variant: 'info', icon: Truck },
      delivered: { label: '已送達', variant: 'success', icon: CheckCircle },
      completed: { label: '已完成', variant: 'success', icon: CheckCircle },
      cancelled: { label: '已取消', variant: 'destructive', icon: XCircle },
    }

    const config = statusConfig[status] || { label: status, variant: 'default', icon: Package }
    const Icon = config.icon
    
    return (
      <Badge variant={config.variant as any} className="flex items-center gap-1">
        <Icon className="h-3 w-3" />
        {config.label}
      </Badge>
    )
  }

  const getPaymentStatusBadge = (status: string) => {
    const statusConfig: Record<string, { label: string; variant: any; icon: any }> = {
      pending: { label: '未付款', variant: 'secondary', icon: Clock },
      partial: { label: '部分付款', variant: 'warning', icon: CreditCard },
      paid: { label: '已付款', variant: 'success', icon: CheckCircle },
    }

    const config = statusConfig[status] || { label: status, variant: 'default', icon: CreditCard }
    const Icon = config.icon
    
    return (
      <Badge variant={config.variant as any} className="flex items-center gap-1">
        <Icon className="h-3 w-3" />
        {config.label}
      </Badge>
    )
  }

  const getNextStatuses = (currentStatus: string) => {
    const statusFlow: Record<string, string[]> = {
      pending: ['confirmed', 'cancelled'],
      confirmed: ['in_production', 'cancelled'],
      in_production: ['quality_check', 'cancelled'],
      quality_check: ['ready_to_ship', 'in_production', 'cancelled'],
      ready_to_ship: ['shipped', 'cancelled'],
      shipped: ['delivered', 'cancelled'],
      delivered: ['completed'],
      completed: [],
      cancelled: [],
    }

    return statusFlow[currentStatus] || []
  }

  const handleStatusUpdate = () => {
    if (!selectedStatus) return
    updateStatusMutation.mutate({ status: selectedStatus, notes: statusNotes })
  }

  const canUpdateStatus = ['admin', 'manager', 'engineer'].includes(user?.role || '')

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

  const nextStatuses = getNextStatuses(order.status)

  return (
    <DashboardLayout>
      <div className="max-w-7xl mx-auto space-y-6">
        {/* Header */}
        <div className="flex justify-between items-start">
          <div className="flex items-center gap-4">
            <Button
              variant="ghost"
              size="icon"
              onClick={() => router.push('/orders')}
            >
              <ArrowLeft className="h-4 w-4" />
            </Button>
            <div>
              <h1 className="text-3xl font-bold text-gray-900">訂單詳情</h1>
              <p className="mt-1 text-gray-600">
                訂單號：{order.order_no} | PO 號碼：{order.po_number}
              </p>
            </div>
          </div>
          <div className="flex gap-2">
            {canUpdateStatus && (
              <Button variant="outline" onClick={() => router.push(`/orders/${orderId}/edit`)}>
                <Edit className="mr-2 h-4 w-4" />
                編輯訂單
              </Button>
            )}
            {canUpdateStatus && nextStatuses.length > 0 && (
              <Button onClick={() => setIsStatusDialogOpen(true)}>
                更新狀態
              </Button>
            )}
          </div>
        </div>

        <Tabs value={activeTab} onValueChange={setActiveTab}>
          <TabsList>
            <TabsTrigger value="details">訂單詳情</TabsTrigger>
            <TabsTrigger value="items">產品明細</TabsTrigger>
            <TabsTrigger value="documents">相關文件</TabsTrigger>
            <TabsTrigger value="timeline">時間軸</TabsTrigger>
          </TabsList>

          <TabsContent value="details" className="space-y-6">
            {/* Basic Info */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center justify-between">
                  基本資訊
                  {getStatusBadge(order.status)}
                </CardTitle>
              </CardHeader>
              <CardContent className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div className="space-y-4">
                  <div>
                    <Label className="text-gray-500">客戶</Label>
                    <div className="flex items-center gap-2 mt-1">
                      <Building2 className="h-4 w-4 text-gray-400" />
                      <span className="font-medium">{order.customer?.name}</span>
                    </div>
                  </div>
                  <div>
                    <Label className="text-gray-500">業務人員</Label>
                    <div className="flex items-center gap-2 mt-1">
                      <User className="h-4 w-4 text-gray-400" />
                      <span>{order.sales?.full_name}</span>
                    </div>
                  </div>
                  <div>
                    <Label className="text-gray-500">報價單號</Label>
                    <div className="flex items-center gap-2 mt-1">
                      <FileText className="h-4 w-4 text-gray-400" />
                      <Button
                        variant="link"
                        className="p-0 h-auto"
                        onClick={() => router.push(`/quotes/${order.quote_id}`)}
                      >
                        {order.quote?.quote_no}
                      </Button>
                    </div>
                  </div>
                </div>
                <div className="space-y-4">
                  <div>
                    <Label className="text-gray-500">交貨日期</Label>
                    <div className="flex items-center gap-2 mt-1">
                      <Calendar className="h-4 w-4 text-gray-400" />
                      <span>{format(new Date(order.delivery_date), 'yyyy/MM/dd', { locale: zhTW })}</span>
                    </div>
                  </div>
                  <div>
                    <Label className="text-gray-500">交貨方式</Label>
                    <div className="flex items-center gap-2 mt-1">
                      <Truck className="h-4 w-4 text-gray-400" />
                      <span>{order.delivery_method}</span>
                    </div>
                  </div>
                  <div>
                    <Label className="text-gray-500">建立日期</Label>
                    <div className="flex items-center gap-2 mt-1">
                      <Calendar className="h-4 w-4 text-gray-400" />
                      <span>{format(new Date(order.created_at), 'yyyy/MM/dd HH:mm', { locale: zhTW })}</span>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Financial Info */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center justify-between">
                  財務資訊
                  {getPaymentStatusBadge(order.payment_status)}
                </CardTitle>
              </CardHeader>
              <CardContent className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div className="space-y-4">
                  <div>
                    <Label className="text-gray-500">訂單總額</Label>
                    <div className="flex items-center gap-2 mt-1">
                      <DollarSign className="h-4 w-4 text-gray-400" />
                      <span className="text-xl font-semibold">
                        {order.currency} {order.total_amount.toLocaleString()}
                      </span>
                    </div>
                  </div>
                  <div>
                    <Label className="text-gray-500">訂購數量</Label>
                    <div className="flex items-center gap-2 mt-1">
                      <Package className="h-4 w-4 text-gray-400" />
                      <span>{order.quantity.toLocaleString()} pcs</span>
                    </div>
                  </div>
                  <div>
                    <Label className="text-gray-500">單價</Label>
                    <div className="flex items-center gap-2 mt-1">
                      <DollarSign className="h-4 w-4 text-gray-400" />
                      <span>{order.currency} {order.unit_price.toFixed(4)}</span>
                    </div>
                  </div>
                </div>
                <div className="space-y-4">
                  <div>
                    <Label className="text-gray-500">付款條件</Label>
                    <div className="flex items-center gap-2 mt-1">
                      <CreditCard className="h-4 w-4 text-gray-400" />
                      <span>{order.payment_terms}</span>
                    </div>
                  </div>
                  <div>
                    <Label className="text-gray-500">頭期款</Label>
                    <div className="flex items-center gap-2 mt-1">
                      <DollarSign className="h-4 w-4 text-gray-400" />
                      <span>{order.currency} {order.down_payment.toLocaleString()}</span>
                    </div>
                  </div>
                  <div>
                    <Label className="text-gray-500">已付金額</Label>
                    <div className="flex items-center gap-2 mt-1">
                      <DollarSign className="h-4 w-4 text-gray-400" />
                      <span>{order.currency} {order.paid_amount.toLocaleString()}</span>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Shipping Info */}
            <Card>
              <CardHeader>
                <CardTitle>配送資訊</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-2">
                  <Label className="text-gray-500">配送地址</Label>
                  <div className="flex items-start gap-2">
                    <MapPin className="h-4 w-4 text-gray-400 mt-1" />
                    <p className="whitespace-pre-wrap">{order.shipping_address}</p>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Notes */}
            {(order.notes || order.internal_notes) && (
              <Card>
                <CardHeader>
                  <CardTitle>備註</CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  {order.notes && (
                    <div>
                      <Label className="text-gray-500">客戶備註</Label>
                      <p className="mt-1 whitespace-pre-wrap">{order.notes}</p>
                    </div>
                  )}
                  {order.internal_notes && (
                    <div>
                      <Label className="text-gray-500">內部備註</Label>
                      <p className="mt-1 whitespace-pre-wrap">{order.internal_notes}</p>
                    </div>
                  )}
                </CardContent>
              </Card>
            )}
          </TabsContent>

          <TabsContent value="items" className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle>產品明細</CardTitle>
                <CardDescription>訂單包含的產品項目</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="rounded-md border">
                  <table className="w-full">
                    <thead>
                      <tr className="border-b bg-gray-50">
                        <th className="px-4 py-3 text-left text-sm font-medium text-gray-900">料號</th>
                        <th className="px-4 py-3 text-left text-sm font-medium text-gray-900">描述</th>
                        <th className="px-4 py-3 text-left text-sm font-medium text-gray-900">材質</th>
                        <th className="px-4 py-3 text-left text-sm font-medium text-gray-900">表面處理</th>
                        <th className="px-4 py-3 text-right text-sm font-medium text-gray-900">數量</th>
                        <th className="px-4 py-3 text-right text-sm font-medium text-gray-900">單價</th>
                        <th className="px-4 py-3 text-right text-sm font-medium text-gray-900">總價</th>
                      </tr>
                    </thead>
                    <tbody>
                      {items?.map((item) => (
                        <tr key={item.id} className="border-b">
                          <td className="px-4 py-3 text-sm font-medium">{item.part_no}</td>
                          <td className="px-4 py-3 text-sm">{item.description || '-'}</td>
                          <td className="px-4 py-3 text-sm">{item.material || '-'}</td>
                          <td className="px-4 py-3 text-sm">{item.surface_treatment || '-'}</td>
                          <td className="px-4 py-3 text-sm text-right">{item.quantity.toLocaleString()}</td>
                          <td className="px-4 py-3 text-sm text-right">${item.unit_price.toFixed(4)}</td>
                          <td className="px-4 py-3 text-sm text-right font-medium">
                            ${item.total_price.toLocaleString()}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                    <tfoot>
                      <tr className="bg-gray-50">
                        <td colSpan={6} className="px-4 py-3 text-sm font-medium text-right">
                          總計
                        </td>
                        <td className="px-4 py-3 text-sm font-bold text-right">
                          {order.currency} {order.total_amount.toLocaleString()}
                        </td>
                      </tr>
                    </tfoot>
                  </table>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="documents" className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle>相關文件</CardTitle>
                <CardDescription>訂單相關的所有文件</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {documents?.length === 0 ? (
                    <p className="text-center text-gray-500 py-8">尚無相關文件</p>
                  ) : (
                    documents?.map((doc) => (
                      <div key={doc.id} className="flex items-center justify-between p-4 border rounded-lg">
                        <div className="flex items-center gap-3">
                          <FileText className="h-8 w-8 text-gray-400" />
                          <div>
                            <p className="font-medium">{doc.file_name}</p>
                            <p className="text-sm text-gray-500">
                              {doc.document_type} • {(doc.file_size / 1024).toFixed(2)} KB • 
                              上傳者: {doc.uploader?.full_name}
                            </p>
                          </div>
                        </div>
                        <div className="flex gap-2">
                          <Button variant="ghost" size="sm">
                            <Download className="h-4 w-4" />
                          </Button>
                          <Button variant="ghost" size="sm">
                            <Trash2 className="h-4 w-4 text-red-500" />
                          </Button>
                        </div>
                      </div>
                    ))
                  )}
                </div>
                <div className="mt-4">
                  <Button variant="outline" className="w-full">
                    <Upload className="mr-2 h-4 w-4" />
                    上傳文件
                  </Button>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="timeline" className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle>活動時間軸</CardTitle>
                <CardDescription>訂單的所有活動記錄</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="relative space-y-8">
                  {activities?.map((activity, index) => (
                    <div key={activity.id} className="relative flex gap-4">
                      {index !== activities.length - 1 && (
                        <div className="absolute left-4 top-8 bottom-0 w-0.5 bg-gray-200" />
                      )}
                      <div className="relative flex h-8 w-8 items-center justify-center rounded-full bg-gray-100">
                        <History className="h-4 w-4 text-gray-600" />
                      </div>
                      <div className="flex-1 space-y-1">
                        <p className="text-sm font-medium">{activity.description}</p>
                        <p className="text-xs text-gray-500">
                          {activity.user?.full_name} • {format(new Date(activity.created_at), 'yyyy/MM/dd HH:mm', { locale: zhTW })}
                        </p>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>

        {/* Status Update Dialog */}
        <Dialog open={isStatusDialogOpen} onOpenChange={setIsStatusDialogOpen}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>更新訂單狀態</DialogTitle>
              <DialogDescription>
                選擇新的訂單狀態並添加備註
              </DialogDescription>
            </DialogHeader>
            <div className="space-y-4">
              <div>
                <Label>新狀態</Label>
                <Select value={selectedStatus} onValueChange={setSelectedStatus}>
                  <SelectTrigger>
                    <SelectValue placeholder="選擇狀態" />
                  </SelectTrigger>
                  <SelectContent>
                    {nextStatuses.map((status) => (
                      <SelectItem key={status} value={status}>
                        {getStatusBadge(status)}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <div>
                <Label>備註</Label>
                <Textarea
                  placeholder="輸入狀態更新備註..."
                  value={statusNotes}
                  onChange={(e) => setStatusNotes(e.target.value)}
                  rows={3}
                />
              </div>
            </div>
            <DialogFooter>
              <Button variant="outline" onClick={() => setIsStatusDialogOpen(false)}>
                取消
              </Button>
              <Button onClick={handleStatusUpdate} disabled={!selectedStatus}>
                更新狀態
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>
    </DashboardLayout>
  )
}