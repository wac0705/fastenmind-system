'use client'

import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useRouter } from 'next/navigation'
import DashboardLayout from '@/components/layout/DashboardLayout'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
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
  Plus,
  Search,
  Filter,
  Download,
  Upload,
  Eye,
  Edit,
  FileText,
  MapPin,
  Calendar,
  Clock,
  CheckCircle,
  AlertTriangle,
  BarChart3,
  Globe,
  DollarSign,
  Weight,
  Box,
  Navigation,
  AlertCircle,
  RefreshCw
} from 'lucide-react'
import { toast } from '@/components/ui/use-toast'
import tradeService from '@/services/trade.service'
import { format } from 'date-fns'
import { zhTW } from 'date-fns/locale'

export default function ShipmentsPage() {
  const router = useRouter()
  const queryClient = useQueryClient()
  const [activeTab, setActiveTab] = useState('all')
  const [searchQuery, setSearchQuery] = useState('')
  const [statusFilter, setStatusFilter] = useState('')
  const [methodFilter, setMethodFilter] = useState('')
  const [typeFilter, setTypeFilter] = useState('')
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false)
  const [selectedShipment, setSelectedShipment] = useState<any>(null)

  // Form state for creating shipment
  const [formData, setFormData] = useState({
    shipment_no: '',
    type: 'export',
    method: 'sea',
    carrier_name: '',
    tracking_no: '',
    container_no: '',
    container_type: '40ft',
    origin_country: 'TW',
    origin_port: '',
    origin_address: '',
    dest_country: '',
    dest_port: '',
    dest_address: '',
    gross_weight: 0,
    net_weight: 0,
    volume: 0,
    package_count: 0,
    package_type: 'carton',
    insurance_value: 0,
    insurance_currency: 'USD',
    freight_cost: 0,
    freight_currency: 'USD',
    customs_value: 0,
    customs_currency: 'USD',
    special_instructions: '',
    internal_notes: '',
  })

  // Fetch shipments
  const { data: shipments, isLoading, refetch } = useQuery({
    queryKey: ['shipments', statusFilter, methodFilter, typeFilter],
    queryFn: () => tradeService.listShipments({
      status: statusFilter || undefined,
      method: methodFilter || undefined,
      type: typeFilter || undefined,
    }),
  })

  // Create shipment mutation
  const createShipmentMutation = useMutation({
    mutationFn: (data: any) => tradeService.createShipment(data),
    onSuccess: () => {
      toast({
        title: '成功',
        description: '運輸單已成功建立',
      })
      queryClient.invalidateQueries({ queryKey: ['shipments'] })
      setIsCreateDialogOpen(false)
      resetForm()
    },
    onError: (error: any) => {
      toast({
        title: '錯誤',
        description: error.response?.data?.message || '建立運輸單失敗',
        variant: 'destructive',
      })
    },
  })

  // Update shipment mutation
  const updateShipmentMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: any }) => 
      tradeService.updateShipment(id, data),
    onSuccess: () => {
      toast({
        title: '成功',
        description: '運輸單已更新',
      })
      queryClient.invalidateQueries({ queryKey: ['shipments'] })
    },
    onError: (error: any) => {
      toast({
        title: '錯誤',
        description: error.response?.data?.message || '更新運輸單失敗',
        variant: 'destructive',
      })
    },
  })

  const resetForm = () => {
    setFormData({
      shipment_no: '',
      type: 'export',
      method: 'sea',
      carrier_name: '',
      tracking_no: '',
      container_no: '',
      container_type: '40ft',
      origin_country: 'TW',
      origin_port: '',
      origin_address: '',
      dest_country: '',
      dest_port: '',
      dest_address: '',
      gross_weight: 0,
      net_weight: 0,
      volume: 0,
      package_count: 0,
      package_type: 'carton',
      insurance_value: 0,
      insurance_currency: 'USD',
      freight_cost: 0,
      freight_currency: 'USD',
      customs_value: 0,
      customs_currency: 'USD',
      special_instructions: '',
      internal_notes: '',
    })
  }

  const handleCreateShipment = () => {
    if (!formData.shipment_no || !formData.dest_country) {
      toast({
        title: '錯誤',
        description: '請填寫必要欄位',
        variant: 'destructive',
      })
      return
    }

    createShipmentMutation.mutate(formData)
  }

  const handleUpdateStatus = (shipmentId: string, newStatus: string) => {
    updateShipmentMutation.mutate({
      id: shipmentId,
      data: { status: newStatus },
    })
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

  const getMethodIcon = (method: string) => {
    const methodIcons: Record<string, any> = {
      sea: Ship,
      air: Plane,
      land: Truck,
      express: Package,
    }
    
    const Icon = methodIcons[method] || Truck
    return <Icon className="h-4 w-4" />
  }

  const getMethodLabel = (method: string) => {
    const labels: Record<string, string> = {
      sea: '海運',
      air: '空運',
      land: '陸運',
      express: '快遞',
    }
    return labels[method] || method
  }

  const getTypeLabel = (type: string) => {
    const labels: Record<string, string> = {
      import: '進口',
      export: '出口',
    }
    return labels[type] || type
  }

  const filteredShipments = shipments?.data?.filter((shipment) => {
    if (searchQuery) {
      const query = searchQuery.toLowerCase()
      return (
        shipment.shipment_no.toLowerCase().includes(query) ||
        shipment.carrier_name?.toLowerCase().includes(query) ||
        shipment.tracking_no?.toLowerCase().includes(query) ||
        shipment.container_no?.toLowerCase().includes(query)
      )
    }
    return true
  }) || []

  // Filter by active tab
  const tabFilteredShipments = activeTab === 'all' 
    ? filteredShipments 
    : filteredShipments.filter(s => {
        if (activeTab === 'in_transit') return s.status === 'in_transit'
        if (activeTab === 'customs') return s.status === 'customs'
        if (activeTab === 'delivered') return s.status === 'delivered'
        return true
      })

  return (
    <DashboardLayout>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex justify-between items-center">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">運輸管理</h1>
            <p className="mt-2 text-gray-600">管理進出口運輸單、追蹤貨物狀態</p>
          </div>
          <div className="flex items-center gap-4">
            <Button variant="outline">
              <Upload className="mr-2 h-4 w-4" />
              匯入運輸單
            </Button>
            <Button variant="outline">
              <Download className="mr-2 h-4 w-4" />
              匯出列表
            </Button>
            <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
              <DialogTrigger asChild>
                <Button>
                  <Plus className="mr-2 h-4 w-4" />
                  新增運輸單
                </Button>
              </DialogTrigger>
              <DialogContent className="max-w-3xl max-h-[90vh] overflow-y-auto">
                <DialogHeader>
                  <DialogTitle>新增運輸單</DialogTitle>
                  <DialogDescription>
                    填寫運輸資訊以建立新的運輸單
                  </DialogDescription>
                </DialogHeader>
                <div className="grid gap-4 py-4">
                  <div className="grid grid-cols-2 gap-4">
                    <div className="space-y-2">
                      <Label htmlFor="shipment_no">運輸單號*</Label>
                      <Input
                        id="shipment_no"
                        value={formData.shipment_no}
                        onChange={(e) => setFormData({ ...formData, shipment_no: e.target.value })}
                        placeholder="SH-2024-001"
                      />
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="type">類型*</Label>
                      <Select value={formData.type} onValueChange={(value) => setFormData({ ...formData, type: value })}>
                        <SelectTrigger>
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="import">進口</SelectItem>
                          <SelectItem value="export">出口</SelectItem>
                        </SelectContent>
                      </Select>
                    </div>
                  </div>

                  <div className="grid grid-cols-2 gap-4">
                    <div className="space-y-2">
                      <Label htmlFor="method">運輸方式*</Label>
                      <Select value={formData.method} onValueChange={(value) => setFormData({ ...formData, method: value })}>
                        <SelectTrigger>
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="sea">海運</SelectItem>
                          <SelectItem value="air">空運</SelectItem>
                          <SelectItem value="land">陸運</SelectItem>
                          <SelectItem value="express">快遞</SelectItem>
                        </SelectContent>
                      </Select>
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="carrier_name">承運商</Label>
                      <Input
                        id="carrier_name"
                        value={formData.carrier_name}
                        onChange={(e) => setFormData({ ...formData, carrier_name: e.target.value })}
                        placeholder="輸入承運商名稱"
                      />
                    </div>
                  </div>

                  <div className="grid grid-cols-2 gap-4">
                    <div className="space-y-2">
                      <Label htmlFor="tracking_no">追蹤號碼</Label>
                      <Input
                        id="tracking_no"
                        value={formData.tracking_no}
                        onChange={(e) => setFormData({ ...formData, tracking_no: e.target.value })}
                        placeholder="輸入追蹤號碼"
                      />
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="container_no">貨櫃號碼</Label>
                      <Input
                        id="container_no"
                        value={formData.container_no}
                        onChange={(e) => setFormData({ ...formData, container_no: e.target.value })}
                        placeholder="輸入貨櫃號碼"
                      />
                    </div>
                  </div>

                  <div className="space-y-2">
                    <Label>起運地資訊</Label>
                    <div className="grid grid-cols-3 gap-4">
                      <Input
                        placeholder="國家代碼*"
                        value={formData.origin_country}
                        onChange={(e) => setFormData({ ...formData, origin_country: e.target.value })}
                      />
                      <Input
                        placeholder="港口"
                        value={formData.origin_port}
                        onChange={(e) => setFormData({ ...formData, origin_port: e.target.value })}
                      />
                      <Input
                        placeholder="地址"
                        value={formData.origin_address}
                        onChange={(e) => setFormData({ ...formData, origin_address: e.target.value })}
                      />
                    </div>
                  </div>

                  <div className="space-y-2">
                    <Label>目的地資訊</Label>
                    <div className="grid grid-cols-3 gap-4">
                      <Input
                        placeholder="國家代碼*"
                        value={formData.dest_country}
                        onChange={(e) => setFormData({ ...formData, dest_country: e.target.value })}
                      />
                      <Input
                        placeholder="港口"
                        value={formData.dest_port}
                        onChange={(e) => setFormData({ ...formData, dest_port: e.target.value })}
                      />
                      <Input
                        placeholder="地址"
                        value={formData.dest_address}
                        onChange={(e) => setFormData({ ...formData, dest_address: e.target.value })}
                      />
                    </div>
                  </div>

                  <div className="grid grid-cols-3 gap-4">
                    <div className="space-y-2">
                      <Label htmlFor="gross_weight">毛重 (kg)</Label>
                      <Input
                        id="gross_weight"
                        type="number"
                        value={formData.gross_weight}
                        onChange={(e) => setFormData({ ...formData, gross_weight: parseFloat(e.target.value) || 0 })}
                      />
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="net_weight">淨重 (kg)</Label>
                      <Input
                        id="net_weight"
                        type="number"
                        value={formData.net_weight}
                        onChange={(e) => setFormData({ ...formData, net_weight: parseFloat(e.target.value) || 0 })}
                      />
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="volume">體積 (m³)</Label>
                      <Input
                        id="volume"
                        type="number"
                        value={formData.volume}
                        onChange={(e) => setFormData({ ...formData, volume: parseFloat(e.target.value) || 0 })}
                      />
                    </div>
                  </div>

                  <div className="grid grid-cols-3 gap-4">
                    <div className="space-y-2">
                      <Label htmlFor="customs_value">報關價值</Label>
                      <Input
                        id="customs_value"
                        type="number"
                        value={formData.customs_value}
                        onChange={(e) => setFormData({ ...formData, customs_value: parseFloat(e.target.value) || 0 })}
                      />
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="freight_cost">運費</Label>
                      <Input
                        id="freight_cost"
                        type="number"
                        value={formData.freight_cost}
                        onChange={(e) => setFormData({ ...formData, freight_cost: parseFloat(e.target.value) || 0 })}
                      />
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="insurance_value">保險價值</Label>
                      <Input
                        id="insurance_value"
                        type="number"
                        value={formData.insurance_value}
                        onChange={(e) => setFormData({ ...formData, insurance_value: parseFloat(e.target.value) || 0 })}
                      />
                    </div>
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="special_instructions">特殊指示</Label>
                    <Input
                      id="special_instructions"
                      value={formData.special_instructions}
                      onChange={(e) => setFormData({ ...formData, special_instructions: e.target.value })}
                      placeholder="輸入特殊處理指示"
                    />
                  </div>
                </div>
                <DialogFooter>
                  <Button variant="outline" onClick={() => setIsCreateDialogOpen(false)}>
                    取消
                  </Button>
                  <Button onClick={handleCreateShipment} disabled={createShipmentMutation.isPending}>
                    {createShipmentMutation.isPending ? '建立中...' : '建立運輸單'}
                  </Button>
                </DialogFooter>
              </DialogContent>
            </Dialog>
          </div>
        </div>

        {/* Filters */}
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center gap-4">
              <div className="flex-1">
                <div className="relative">
                  <Search className="absolute left-3 top-2.5 h-4 w-4 text-gray-400" />
                  <Input
                    placeholder="搜尋運輸單號、承運商、追蹤號碼..."
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    className="pl-10"
                  />
                </div>
              </div>
              <Select value={typeFilter} onValueChange={setTypeFilter}>
                <SelectTrigger className="w-32">
                  <SelectValue placeholder="類型" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="">全部類型</SelectItem>
                  <SelectItem value="import">進口</SelectItem>
                  <SelectItem value="export">出口</SelectItem>
                </SelectContent>
              </Select>
              <Select value={methodFilter} onValueChange={setMethodFilter}>
                <SelectTrigger className="w-32">
                  <SelectValue placeholder="運輸方式" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="">全部方式</SelectItem>
                  <SelectItem value="sea">海運</SelectItem>
                  <SelectItem value="air">空運</SelectItem>
                  <SelectItem value="land">陸運</SelectItem>
                  <SelectItem value="express">快遞</SelectItem>
                </SelectContent>
              </Select>
              <Select value={statusFilter} onValueChange={setStatusFilter}>
                <SelectTrigger className="w-32">
                  <SelectValue placeholder="狀態" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="">全部狀態</SelectItem>
                  <SelectItem value="pending">待處理</SelectItem>
                  <SelectItem value="in_transit">運輸中</SelectItem>
                  <SelectItem value="customs">清關中</SelectItem>
                  <SelectItem value="delivered">已送達</SelectItem>
                  <SelectItem value="cancelled">已取消</SelectItem>
                </SelectContent>
              </Select>
              <Button variant="outline" size="icon" onClick={() => refetch()}>
                <RefreshCw className="h-4 w-4" />
              </Button>
            </div>
          </CardContent>
        </Card>

        {/* Tabs */}
        <Tabs value={activeTab} onValueChange={setActiveTab}>
          <TabsList>
            <TabsTrigger value="all">
              全部 ({filteredShipments.length})
            </TabsTrigger>
            <TabsTrigger value="in_transit">
              運輸中 ({filteredShipments.filter(s => s.status === 'in_transit').length})
            </TabsTrigger>
            <TabsTrigger value="customs">
              清關中 ({filteredShipments.filter(s => s.status === 'customs').length})
            </TabsTrigger>
            <TabsTrigger value="delivered">
              已送達 ({filteredShipments.filter(s => s.status === 'delivered').length})
            </TabsTrigger>
          </TabsList>

          <TabsContent value={activeTab} className="mt-6">
            <Card>
              <CardContent className="p-0">
                {isLoading ? (
                  <div className="flex items-center justify-center h-64">
                    <div className="text-center">載入中...</div>
                  </div>
                ) : tabFilteredShipments.length === 0 ? (
                  <div className="flex flex-col items-center justify-center h-64 text-gray-500">
                    <Ship className="h-12 w-12 mb-4 text-gray-300" />
                    <p>沒有找到運輸單</p>
                  </div>
                ) : (
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead>運輸單號</TableHead>
                        <TableHead>類型</TableHead>
                        <TableHead>運輸方式</TableHead>
                        <TableHead>路線</TableHead>
                        <TableHead>狀態</TableHead>
                        <TableHead className="text-right">報關價值</TableHead>
                        <TableHead className="text-right">重量/體積</TableHead>
                        <TableHead>建立時間</TableHead>
                        <TableHead className="text-right">操作</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {tabFilteredShipments.map((shipment) => (
                        <TableRow key={shipment.id}>
                          <TableCell className="font-medium">
                            <div>
                              <p>{shipment.shipment_no}</p>
                              {shipment.tracking_no && (
                                <p className="text-sm text-gray-500">追蹤: {shipment.tracking_no}</p>
                              )}
                            </div>
                          </TableCell>
                          <TableCell>
                            <Badge variant="outline" className="text-xs">
                              {tradeService.getShipmentTypeIcon(shipment.type)} {getTypeLabel(shipment.type)}
                            </Badge>
                          </TableCell>
                          <TableCell>
                            <div className="flex items-center gap-2">
                              {getMethodIcon(shipment.method)}
                              <span>{getMethodLabel(shipment.method)}</span>
                            </div>
                          </TableCell>
                          <TableCell>
                            <div className="text-sm">
                              <p>{shipment.origin_country} → {shipment.dest_country}</p>
                              {shipment.carrier_name && (
                                <p className="text-gray-500">{shipment.carrier_name}</p>
                              )}
                            </div>
                          </TableCell>
                          <TableCell>{getStatusBadge(shipment.status)}</TableCell>
                          <TableCell className="text-right">
                            {tradeService.formatCurrency(shipment.customs_value, shipment.customs_currency)}
                          </TableCell>
                          <TableCell className="text-right">
                            <div className="text-sm">
                              <p>{tradeService.formatWeight(shipment.gross_weight)}</p>
                              <p className="text-gray-500">{tradeService.formatVolume(shipment.volume)}</p>
                            </div>
                          </TableCell>
                          <TableCell>
                            <div className="text-sm">
                              <p>{format(new Date(shipment.created_at), 'yyyy/MM/dd', { locale: zhTW })}</p>
                              <p className="text-gray-500">{format(new Date(shipment.created_at), 'HH:mm', { locale: zhTW })}</p>
                            </div>
                          </TableCell>
                          <TableCell className="text-right">
                            <div className="flex items-center justify-end gap-2">
                              <Button
                                variant="ghost"
                                size="icon"
                                onClick={() => router.push(`/trade/shipments/${shipment.id}`)}
                              >
                                <Eye className="h-4 w-4" />
                              </Button>
                              {shipment.status === 'pending' && (
                                <Button
                                  variant="ghost"
                                  size="icon"
                                  onClick={() => handleUpdateStatus(shipment.id, 'in_transit')}
                                >
                                  <Truck className="h-4 w-4" />
                                </Button>
                              )}
                              {shipment.status === 'in_transit' && (
                                <Button
                                  variant="ghost"
                                  size="icon"
                                  onClick={() => handleUpdateStatus(shipment.id, 'customs')}
                                >
                                  <FileText className="h-4 w-4" />
                                </Button>
                              )}
                              {shipment.status === 'customs' && (
                                <Button
                                  variant="ghost"
                                  size="icon"
                                  onClick={() => handleUpdateStatus(shipment.id, 'delivered')}
                                >
                                  <CheckCircle className="h-4 w-4" />
                                </Button>
                              )}
                            </div>
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                )}
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </div>
    </DashboardLayout>
  )
}