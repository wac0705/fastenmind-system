'use client'

import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useRouter, useParams } from 'next/navigation'
import DashboardLayout from '@/components/layout/DashboardLayout'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Textarea } from '@/components/ui/textarea'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'
import { 
  Ship,
  Truck,
  Plane,
  Package,
  MapPin,
  Calendar,
  Clock,
  CheckCircle,
  AlertTriangle,
  FileText,
  Plus,
  Edit,
  Download,
  Upload,
  Navigation,
  Weight,
  Box,
  DollarSign,
  Globe,
  History,
  AlertCircle,
  ArrowLeft,
  BarChart3,
  User
} from 'lucide-react'
import { toast } from '@/components/ui/use-toast'
import tradeService from '@/services/trade.service'
import { format } from 'date-fns'
import { zhTW } from 'date-fns/locale'

export default function ShipmentDetailPage() {
  const router = useRouter()
  const params = useParams()
  const shipmentId = params.id as string
  const queryClient = useQueryClient()
  
  const [activeTab, setActiveTab] = useState('details')
  const [isEventDialogOpen, setIsEventDialogOpen] = useState(false)
  const [isUpdateDialogOpen, setIsUpdateDialogOpen] = useState(false)
  
  // Event form state
  const [eventForm, setEventForm] = useState({
    event_type: 'update',
    status: 'completed',
    location: '',
    description: '',
    event_time: format(new Date(), "yyyy-MM-dd'T'HH:mm"),
  })
  
  // Update form state
  const [updateForm, setUpdateForm] = useState({
    status: '',
    tracking_no: '',
    actual_departure: '',
    actual_arrival: '',
    total_duty: 0,
    total_tax: 0,
    special_instructions: '',
    internal_notes: '',
  })

  // Fetch shipment details
  const { data: shipment, isLoading } = useQuery({
    queryKey: ['shipment', shipmentId],
    queryFn: () => tradeService.getShipment(shipmentId),
  })

  // Fetch shipment events
  const { data: events } = useQuery({
    queryKey: ['shipment-events', shipmentId],
    queryFn: () => tradeService.getShipmentEvents(shipmentId),
  })

  // Fetch trade documents
  const { data: documents } = useQuery({
    queryKey: ['shipment-documents', shipmentId],
    queryFn: () => tradeService.getTradeDocumentsByShipment(shipmentId),
  })

  // Create event mutation
  const createEventMutation = useMutation({
    mutationFn: (data: any) => tradeService.createShipmentEvent(shipmentId, data),
    onSuccess: () => {
      toast({
        title: '成功',
        description: '運輸事件已記錄',
      })
      queryClient.invalidateQueries({ queryKey: ['shipment-events', shipmentId] })
      setIsEventDialogOpen(false)
      resetEventForm()
    },
    onError: (error: any) => {
      toast({
        title: '錯誤',
        description: error.response?.data?.message || '記錄事件失敗',
        variant: 'destructive',
      })
    },
  })

  // Update shipment mutation
  const updateShipmentMutation = useMutation({
    mutationFn: (data: any) => tradeService.updateShipment(shipmentId, data),
    onSuccess: () => {
      toast({
        title: '成功',
        description: '運輸單已更新',
      })
      queryClient.invalidateQueries({ queryKey: ['shipment', shipmentId] })
      setIsUpdateDialogOpen(false)
    },
    onError: (error: any) => {
      toast({
        title: '錯誤',
        description: error.response?.data?.message || '更新運輸單失敗',
        variant: 'destructive',
      })
    },
  })

  const resetEventForm = () => {
    setEventForm({
      event_type: 'update',
      status: 'completed',
      location: '',
      description: '',
      event_time: format(new Date(), "yyyy-MM-dd'T'HH:mm"),
    })
  }

  const handleCreateEvent = () => {
    if (!eventForm.event_type || !eventForm.status) {
      toast({
        title: '錯誤',
        description: '請填寫必要欄位',
        variant: 'destructive',
      })
      return
    }

    createEventMutation.mutate({
      ...eventForm,
      event_time: new Date(eventForm.event_time).toISOString(),
    })
  }

  const handleUpdateShipment = () => {
    const updateData: any = {}
    
    if (updateForm.status) updateData.status = updateForm.status
    if (updateForm.tracking_no) updateData.tracking_no = updateForm.tracking_no
    if (updateForm.actual_departure) updateData.actual_departure = new Date(updateForm.actual_departure).toISOString()
    if (updateForm.actual_arrival) updateData.actual_arrival = new Date(updateForm.actual_arrival).toISOString()
    if (updateForm.total_duty > 0) updateData.total_duty = updateForm.total_duty
    if (updateForm.total_tax > 0) updateData.total_tax = updateForm.total_tax
    if (updateForm.special_instructions) updateData.special_instructions = updateForm.special_instructions
    if (updateForm.internal_notes) updateData.internal_notes = updateForm.internal_notes

    updateShipmentMutation.mutate(updateData)
  }

  const getStatusBadge = (status: string) => {
    const statusConfig: Record<string, { label: string; variant: any; icon: any }> = {
      pending: { label: '待處理', variant: 'secondary', icon: Clock },
      in_transit: { label: '運輸中', variant: 'info', icon: Truck },
      customs: { label: '清關中', variant: 'warning', icon: FileText },
      delivered: { label: '已送達', variant: 'success', icon: CheckCircle },
      cancelled: { label: '已取消', variant: 'destructive', icon: AlertTriangle },
    }

    const config = statusConfig[status] || { label: status, variant: 'default', icon: AlertTriangle }
    const Icon = config.icon
    
    return (
      <Badge variant={config.variant as any} className="flex items-center gap-1">
        <Icon className="h-3 w-3" />
        {config.label}
      </Badge>
    )
  }

  const getEventIcon = (eventType: string) => {
    const eventIcons: Record<string, any> = {
      created: Plus,
      departure: Navigation,
      arrival: MapPin,
      customs_clearance: FileText,
      delivery: CheckCircle,
      delay: AlertTriangle,
      update: Edit,
      status_change: History,
    }
    
    const Icon = eventIcons[eventType] || AlertCircle
    return <Icon className="h-4 w-4" />
  }

  const getMethodIcon = (method: string) => {
    const methodIcons: Record<string, any> = {
      sea: Ship,
      air: Plane,
      land: Truck,
      express: Package,
    }
    
    const Icon = methodIcons[method] || Truck
    return <Icon className="h-5 w-5" />
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

  if (!shipment) {
    return (
      <DashboardLayout>
        <div className="flex items-center justify-center h-64">
          <div className="text-center">
            <AlertCircle className="h-12 w-12 mx-auto mb-4 text-gray-300" />
            <p className="text-gray-500">找不到運輸單</p>
            <Button variant="outline" onClick={() => router.push('/trade/shipments')} className="mt-4">
              返回列表
            </Button>
          </div>
        </div>
      </DashboardLayout>
    )
  }

  return (
    <DashboardLayout>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex justify-between items-start">
          <div className="flex items-center gap-4">
            <Button variant="ghost" size="icon" onClick={() => router.push('/trade/shipments')}>
              <ArrowLeft className="h-4 w-4" />
            </Button>
            <div>
              <h1 className="text-3xl font-bold text-gray-900 flex items-center gap-3">
                {getMethodIcon(shipment.method)}
                {shipment.shipment_no}
              </h1>
              <p className="mt-2 text-gray-600">
                {shipment.origin_country} → {shipment.dest_country} • {shipment.carrier_name || '未指定承運商'}
              </p>
            </div>
          </div>
          <div className="flex items-center gap-3">
            {getStatusBadge(shipment.status)}
            <Dialog open={isUpdateDialogOpen} onOpenChange={setIsUpdateDialogOpen}>
              <DialogTrigger asChild>
                <Button variant="outline">
                  <Edit className="mr-2 h-4 w-4" />
                  更新資訊
                </Button>
              </DialogTrigger>
              <DialogContent>
                <DialogHeader>
                  <DialogTitle>更新運輸單</DialogTitle>
                  <DialogDescription>
                    更新運輸單狀態和資訊
                  </DialogDescription>
                </DialogHeader>
                <div className="grid gap-4 py-4">
                  <div className="space-y-2">
                    <Label htmlFor="update-status">狀態</Label>
                    <Select value={updateForm.status} onValueChange={(value) => setUpdateForm({ ...updateForm, status: value })}>
                      <SelectTrigger>
                        <SelectValue placeholder="選擇狀態" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="pending">待處理</SelectItem>
                        <SelectItem value="in_transit">運輸中</SelectItem>
                        <SelectItem value="customs">清關中</SelectItem>
                        <SelectItem value="delivered">已送達</SelectItem>
                        <SelectItem value="cancelled">已取消</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="update-tracking">追蹤號碼</Label>
                    <Input
                      id="update-tracking"
                      value={updateForm.tracking_no}
                      onChange={(e) => setUpdateForm({ ...updateForm, tracking_no: e.target.value })}
                      placeholder="輸入追蹤號碼"
                    />
                  </div>
                  <div className="grid grid-cols-2 gap-4">
                    <div className="space-y-2">
                      <Label htmlFor="update-departure">實際出發時間</Label>
                      <Input
                        id="update-departure"
                        type="datetime-local"
                        value={updateForm.actual_departure}
                        onChange={(e) => setUpdateForm({ ...updateForm, actual_departure: e.target.value })}
                      />
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="update-arrival">實際到達時間</Label>
                      <Input
                        id="update-arrival"
                        type="datetime-local"
                        value={updateForm.actual_arrival}
                        onChange={(e) => setUpdateForm({ ...updateForm, actual_arrival: e.target.value })}
                      />
                    </div>
                  </div>
                  <div className="grid grid-cols-2 gap-4">
                    <div className="space-y-2">
                      <Label htmlFor="update-duty">關稅總額</Label>
                      <Input
                        id="update-duty"
                        type="number"
                        value={updateForm.total_duty}
                        onChange={(e) => setUpdateForm({ ...updateForm, total_duty: parseFloat(e.target.value) || 0 })}
                      />
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="update-tax">稅金總額</Label>
                      <Input
                        id="update-tax"
                        type="number"
                        value={updateForm.total_tax}
                        onChange={(e) => setUpdateForm({ ...updateForm, total_tax: parseFloat(e.target.value) || 0 })}
                      />
                    </div>
                  </div>
                </div>
                <DialogFooter>
                  <Button variant="outline" onClick={() => setIsUpdateDialogOpen(false)}>
                    取消
                  </Button>
                  <Button onClick={handleUpdateShipment} disabled={updateShipmentMutation.isPending}>
                    {updateShipmentMutation.isPending ? '更新中...' : '更新'}
                  </Button>
                </DialogFooter>
              </DialogContent>
            </Dialog>
            <Button>
              <Download className="mr-2 h-4 w-4" />
              匯出文件
            </Button>
          </div>
        </div>

        <Tabs value={activeTab} onValueChange={setActiveTab}>
          <TabsList>
            <TabsTrigger value="details">詳細資訊</TabsTrigger>
            <TabsTrigger value="items">運輸項目 ({shipment.items?.length || 0})</TabsTrigger>
            <TabsTrigger value="events">事件記錄 ({events?.data?.length || 0})</TabsTrigger>
            <TabsTrigger value="documents">相關文件 ({documents?.data?.length || 0})</TabsTrigger>
            <TabsTrigger value="costs">費用明細</TabsTrigger>
          </TabsList>

          <TabsContent value="details" className="space-y-6">
            {/* Basic Information */}
            <Card>
              <CardHeader>
                <CardTitle>基本資訊</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-2 gap-6">
                  <div>
                    <Label className="text-gray-500">運輸類型</Label>
                    <p className="font-medium">
                      {shipment.type === 'import' ? '進口' : '出口'}
                    </p>
                  </div>
                  <div>
                    <Label className="text-gray-500">運輸方式</Label>
                    <p className="font-medium flex items-center gap-2">
                      {getMethodIcon(shipment.method)}
                      {shipment.method === 'sea' ? '海運' : 
                       shipment.method === 'air' ? '空運' : 
                       shipment.method === 'land' ? '陸運' : '快遞'}
                    </p>
                  </div>
                  <div>
                    <Label className="text-gray-500">承運商</Label>
                    <p className="font-medium">{shipment.carrier_name || '-'}</p>
                  </div>
                  <div>
                    <Label className="text-gray-500">追蹤號碼</Label>
                    <p className="font-medium">{shipment.tracking_no || '-'}</p>
                  </div>
                  <div>
                    <Label className="text-gray-500">貨櫃號碼</Label>
                    <p className="font-medium">{shipment.container_no || '-'}</p>
                  </div>
                  <div>
                    <Label className="text-gray-500">貨櫃類型</Label>
                    <p className="font-medium">{shipment.container_type || '-'}</p>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Route Information */}
            <Card>
              <CardHeader>
                <CardTitle>路線資訊</CardTitle>
              </CardHeader>
              <CardContent className="space-y-6">
                <div className="grid grid-cols-2 gap-6">
                  <div className="space-y-4">
                    <h4 className="font-medium flex items-center gap-2">
                      <MapPin className="h-4 w-4" />
                      起運地
                    </h4>
                    <div className="space-y-2 pl-6">
                      <div>
                        <Label className="text-gray-500">國家</Label>
                        <p className="font-medium">{shipment.origin_country}</p>
                      </div>
                      <div>
                        <Label className="text-gray-500">港口</Label>
                        <p className="font-medium">{shipment.origin_port || '-'}</p>
                      </div>
                      <div>
                        <Label className="text-gray-500">地址</Label>
                        <p className="font-medium">{shipment.origin_address || '-'}</p>
                      </div>
                    </div>
                  </div>
                  <div className="space-y-4">
                    <h4 className="font-medium flex items-center gap-2">
                      <MapPin className="h-4 w-4" />
                      目的地
                    </h4>
                    <div className="space-y-2 pl-6">
                      <div>
                        <Label className="text-gray-500">國家</Label>
                        <p className="font-medium">{shipment.dest_country}</p>
                      </div>
                      <div>
                        <Label className="text-gray-500">港口</Label>
                        <p className="font-medium">{shipment.dest_port || '-'}</p>
                      </div>
                      <div>
                        <Label className="text-gray-500">地址</Label>
                        <p className="font-medium">{shipment.dest_address || '-'}</p>
                      </div>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Schedule Information */}
            <Card>
              <CardHeader>
                <CardTitle>時程資訊</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-2 gap-6">
                  <div>
                    <Label className="text-gray-500">預計出發時間</Label>
                    <p className="font-medium">
                      {shipment.estimated_departure 
                        ? format(new Date(shipment.estimated_departure), 'yyyy/MM/dd HH:mm', { locale: zhTW })
                        : '-'}
                    </p>
                  </div>
                  <div>
                    <Label className="text-gray-500">實際出發時間</Label>
                    <p className="font-medium">
                      {shipment.actual_departure 
                        ? format(new Date(shipment.actual_departure), 'yyyy/MM/dd HH:mm', { locale: zhTW })
                        : '-'}
                    </p>
                  </div>
                  <div>
                    <Label className="text-gray-500">預計到達時間</Label>
                    <p className="font-medium">
                      {shipment.estimated_arrival 
                        ? format(new Date(shipment.estimated_arrival), 'yyyy/MM/dd HH:mm', { locale: zhTW })
                        : '-'}
                    </p>
                  </div>
                  <div>
                    <Label className="text-gray-500">實際到達時間</Label>
                    <p className="font-medium">
                      {shipment.actual_arrival 
                        ? format(new Date(shipment.actual_arrival), 'yyyy/MM/dd HH:mm', { locale: zhTW })
                        : '-'}
                    </p>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Physical Information */}
            <Card>
              <CardHeader>
                <CardTitle>貨物資訊</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-3 gap-6">
                  <div>
                    <Label className="text-gray-500">毛重</Label>
                    <p className="font-medium flex items-center gap-2">
                      <Weight className="h-4 w-4" />
                      {tradeService.formatWeight(shipment.gross_weight)}
                    </p>
                  </div>
                  <div>
                    <Label className="text-gray-500">淨重</Label>
                    <p className="font-medium flex items-center gap-2">
                      <Weight className="h-4 w-4" />
                      {tradeService.formatWeight(shipment.net_weight)}
                    </p>
                  </div>
                  <div>
                    <Label className="text-gray-500">體積</Label>
                    <p className="font-medium flex items-center gap-2">
                      <Box className="h-4 w-4" />
                      {tradeService.formatVolume(shipment.volume)}
                    </p>
                  </div>
                  <div>
                    <Label className="text-gray-500">包裝數量</Label>
                    <p className="font-medium">{shipment.package_count} 件</p>
                  </div>
                  <div>
                    <Label className="text-gray-500">包裝類型</Label>
                    <p className="font-medium">{shipment.package_type || '-'}</p>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Notes */}
            {(shipment.special_instructions || shipment.internal_notes) && (
              <Card>
                <CardHeader>
                  <CardTitle>備註</CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  {shipment.special_instructions && (
                    <div>
                      <Label className="text-gray-500">特殊指示</Label>
                      <p className="font-medium mt-1">{shipment.special_instructions}</p>
                    </div>
                  )}
                  {shipment.internal_notes && (
                    <div>
                      <Label className="text-gray-500">內部備註</Label>
                      <p className="font-medium mt-1">{shipment.internal_notes}</p>
                    </div>
                  )}
                </CardContent>
              </Card>
            )}
          </TabsContent>

          <TabsContent value="items">
            <Card>
              <CardHeader>
                <CardTitle>運輸項目清單</CardTitle>
              </CardHeader>
              <CardContent>
                {shipment.items && shipment.items.length > 0 ? (
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead>產品名稱</TableHead>
                        <TableHead>HS 代碼</TableHead>
                        <TableHead className="text-right">數量</TableHead>
                        <TableHead className="text-right">單位重量</TableHead>
                        <TableHead className="text-right">總重量</TableHead>
                        <TableHead className="text-right">單價</TableHead>
                        <TableHead className="text-right">總價值</TableHead>
                        <TableHead>產地</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {shipment.items.map((item) => (
                        <TableRow key={item.id}>
                          <TableCell>
                            <div>
                              <p className="font-medium">{item.product_name}</p>
                              {item.description && (
                                <p className="text-sm text-gray-500">{item.description}</p>
                              )}
                            </div>
                          </TableCell>
                          <TableCell>{item.hs_code || '-'}</TableCell>
                          <TableCell className="text-right">
                            {item.quantity} {item.unit}
                          </TableCell>
                          <TableCell className="text-right">
                            {tradeService.formatWeight(item.unit_weight)}
                          </TableCell>
                          <TableCell className="text-right">
                            {tradeService.formatWeight(item.total_weight)}
                          </TableCell>
                          <TableCell className="text-right">
                            {tradeService.formatCurrency(item.unit_value, item.currency)}
                          </TableCell>
                          <TableCell className="text-right">
                            {tradeService.formatCurrency(item.total_value, item.currency)}
                          </TableCell>
                          <TableCell>{item.country_origin || '-'}</TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                ) : (
                  <div className="text-center py-8 text-gray-500">
                    <Package className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                    <p>尚未添加運輸項目</p>
                  </div>
                )}
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="events" className="space-y-4">
            <div className="flex justify-end">
              <Dialog open={isEventDialogOpen} onOpenChange={setIsEventDialogOpen}>
                <DialogTrigger asChild>
                  <Button>
                    <Plus className="mr-2 h-4 w-4" />
                    新增事件
                  </Button>
                </DialogTrigger>
                <DialogContent>
                  <DialogHeader>
                    <DialogTitle>新增運輸事件</DialogTitle>
                    <DialogDescription>
                      記錄運輸過程中的重要事件
                    </DialogDescription>
                  </DialogHeader>
                  <div className="grid gap-4 py-4">
                    <div className="space-y-2">
                      <Label htmlFor="event-type">事件類型</Label>
                      <Select value={eventForm.event_type} onValueChange={(value) => setEventForm({ ...eventForm, event_type: value })}>
                        <SelectTrigger>
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="departure">出發</SelectItem>
                          <SelectItem value="arrival">到達</SelectItem>
                          <SelectItem value="customs_clearance">清關</SelectItem>
                          <SelectItem value="delivery">交付</SelectItem>
                          <SelectItem value="delay">延誤</SelectItem>
                          <SelectItem value="update">更新</SelectItem>
                        </SelectContent>
                      </Select>
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="event-location">地點</Label>
                      <Input
                        id="event-location"
                        value={eventForm.location}
                        onChange={(e) => setEventForm({ ...eventForm, location: e.target.value })}
                        placeholder="輸入事件發生地點"
                      />
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="event-time">時間</Label>
                      <Input
                        id="event-time"
                        type="datetime-local"
                        value={eventForm.event_time}
                        onChange={(e) => setEventForm({ ...eventForm, event_time: e.target.value })}
                      />
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="event-description">描述</Label>
                      <Textarea
                        id="event-description"
                        value={eventForm.description}
                        onChange={(e) => setEventForm({ ...eventForm, description: e.target.value })}
                        placeholder="輸入事件描述"
                        rows={3}
                      />
                    </div>
                  </div>
                  <DialogFooter>
                    <Button variant="outline" onClick={() => setIsEventDialogOpen(false)}>
                      取消
                    </Button>
                    <Button onClick={handleCreateEvent} disabled={createEventMutation.isPending}>
                      {createEventMutation.isPending ? '新增中...' : '新增事件'}
                    </Button>
                  </DialogFooter>
                </DialogContent>
              </Dialog>
            </div>

            <Card>
              <CardContent className="p-0">
                {events?.data && events.data.length > 0 ? (
                  <div className="space-y-0">
                    {events.data.map((event, index) => (
                      <div 
                        key={event.id} 
                        className={`flex items-start gap-4 p-4 ${index !== events.data.length - 1 ? 'border-b' : ''}`}
                      >
                        <div className="flex-shrink-0">
                          <div className="w-8 h-8 rounded-full bg-gray-100 flex items-center justify-center">
                            {getEventIcon(event.event_type)}
                          </div>
                        </div>
                        <div className="flex-1">
                          <div className="flex items-start justify-between">
                            <div>
                              <p className="font-medium">
                                {event.event_type === 'departure' ? '貨物出發' :
                                 event.event_type === 'arrival' ? '貨物到達' :
                                 event.event_type === 'customs_clearance' ? '清關處理' :
                                 event.event_type === 'delivery' ? '貨物交付' :
                                 event.event_type === 'delay' ? '運輸延誤' :
                                 event.event_type === 'status_change' ? '狀態變更' :
                                 '運輸更新'}
                              </p>
                              {event.location && (
                                <p className="text-sm text-gray-500 flex items-center gap-1 mt-1">
                                  <MapPin className="h-3 w-3" />
                                  {event.location}
                                </p>
                              )}
                              {event.description && (
                                <p className="text-sm text-gray-600 mt-2">{event.description}</p>
                              )}
                            </div>
                            <div className="text-right text-sm">
                              <p className="font-medium">
                                {format(new Date(event.event_time), 'yyyy/MM/dd', { locale: zhTW })}
                              </p>
                              <p className="text-gray-500">
                                {format(new Date(event.event_time), 'HH:mm', { locale: zhTW })}
                              </p>
                            </div>
                          </div>
                          <div className="flex items-center gap-4 mt-2 text-xs text-gray-500">
                            <span className="flex items-center gap-1">
                              <User className="h-3 w-3" />
                              {event.creator?.full_name || '系統'}
                            </span>
                            <span className="flex items-center gap-1">
                              <Clock className="h-3 w-3" />
                              記錄於 {format(new Date(event.recorded_at), 'MM/dd HH:mm', { locale: zhTW })}
                            </span>
                            {event.source && (
                              <Badge variant="outline" className="text-xs">
                                {event.source}
                              </Badge>
                            )}
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                ) : (
                  <div className="text-center py-8 text-gray-500">
                    <History className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                    <p>尚無事件記錄</p>
                  </div>
                )}
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="documents">
            <Card>
              <CardHeader>
                <div className="flex justify-between items-center">
                  <CardTitle>相關文件</CardTitle>
                  <Button>
                    <Upload className="mr-2 h-4 w-4" />
                    上傳文件
                  </Button>
                </div>
              </CardHeader>
              <CardContent>
                {documents?.data && documents.data.length > 0 ? (
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead>文件類型</TableHead>
                        <TableHead>文件編號</TableHead>
                        <TableHead>標題</TableHead>
                        <TableHead>狀態</TableHead>
                        <TableHead>發行日期</TableHead>
                        <TableHead className="text-right">操作</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {documents.data.map((doc) => (
                        <TableRow key={doc.id}>
                          <TableCell>
                            <Badge variant="outline">
                              {doc.document_type}
                            </Badge>
                          </TableCell>
                          <TableCell>{doc.document_no}</TableCell>
                          <TableCell>{doc.title}</TableCell>
                          <TableCell>
                            <Badge variant={
                              doc.status === 'approved' ? 'success' as any :
                              doc.status === 'rejected' ? 'destructive' as any :
                              'secondary' as any
                            }>
                              {doc.status}
                            </Badge>
                          </TableCell>
                          <TableCell>
                            {doc.issued_at 
                              ? format(new Date(doc.issued_at), 'yyyy/MM/dd', { locale: zhTW })
                              : '-'}
                          </TableCell>
                          <TableCell className="text-right">
                            <Button variant="ghost" size="icon">
                              <Download className="h-4 w-4" />
                            </Button>
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                ) : (
                  <div className="text-center py-8 text-gray-500">
                    <FileText className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                    <p>尚未上傳相關文件</p>
                  </div>
                )}
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="costs">
            <Card>
              <CardHeader>
                <CardTitle>費用明細</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div className="grid grid-cols-2 gap-6">
                    <div>
                      <h4 className="font-medium mb-4">基本費用</h4>
                      <div className="space-y-3">
                        <div className="flex justify-between items-center">
                          <span className="text-gray-600">報關價值</span>
                          <span className="font-medium">
                            {tradeService.formatCurrency(shipment.customs_value, shipment.customs_currency)}
                          </span>
                        </div>
                        <div className="flex justify-between items-center">
                          <span className="text-gray-600">運費</span>
                          <span className="font-medium">
                            {tradeService.formatCurrency(shipment.freight_cost, shipment.freight_currency)}
                          </span>
                        </div>
                        <div className="flex justify-between items-center">
                          <span className="text-gray-600">保險費</span>
                          <span className="font-medium">
                            {tradeService.formatCurrency(shipment.insurance_value, shipment.insurance_currency)}
                          </span>
                        </div>
                      </div>
                    </div>
                    <div>
                      <h4 className="font-medium mb-4">稅費</h4>
                      <div className="space-y-3">
                        <div className="flex justify-between items-center">
                          <span className="text-gray-600">關稅</span>
                          <span className="font-medium">
                            {tradeService.formatCurrency(shipment.total_duty, shipment.customs_currency)}
                          </span>
                        </div>
                        <div className="flex justify-between items-center">
                          <span className="text-gray-600">其他稅費</span>
                          <span className="font-medium">
                            {tradeService.formatCurrency(shipment.total_tax, shipment.customs_currency)}
                          </span>
                        </div>
                      </div>
                    </div>
                  </div>
                  <div className="border-t pt-4">
                    <div className="flex justify-between items-center text-lg">
                      <span className="font-medium">總費用</span>
                      <span className="font-bold">
                        {tradeService.formatCurrency(
                          shipment.customs_value + 
                          shipment.freight_cost + 
                          shipment.insurance_value + 
                          shipment.total_duty + 
                          shipment.total_tax,
                          shipment.customs_currency
                        )}
                      </span>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </div>
    </DashboardLayout>
  )
}