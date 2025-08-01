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
import { Label } from '@/components/ui/label'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
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
  TrendingUp,
  TrendingDown,
  AlertCircle,
  Edit,
  History,
  Warehouse as WarehouseIcon,
  DollarSign,
  RefreshCw,
  ArrowUpDown,
  BarChart3,
  FileText,
  Truck
} from 'lucide-react'
import { format } from 'date-fns'
import { zhTW } from 'date-fns/locale'
import { useToast } from '@/components/ui/use-toast'
import inventoryService, { Inventory, StockMovement, StockAdjustmentRequest, StockTransferRequest } from '@/services/inventory.service'
import { useAuthStore } from '@/store/auth.store'

export default function InventoryDetailPage() {
  const params = useParams()
  const router = useRouter()
  const queryClient = useQueryClient()
  const { toast } = useToast()
  const { user } = useAuthStore()
  const inventoryId = params.id as string

  const [activeTab, setActiveTab] = useState('details')
  const [isAdjustDialogOpen, setIsAdjustDialogOpen] = useState(false)
  const [isTransferDialogOpen, setIsTransferDialogOpen] = useState(false)
  const [adjustmentData, setAdjustmentData] = useState<StockAdjustmentRequest>({
    quantity: 0,
    reason: '',
    notes: '',
  })
  const [transferData, setTransferData] = useState<StockTransferRequest>({
    inventory_id: inventoryId,
    quantity: 0,
    from_warehouse_id: '',
    to_warehouse_id: '',
    notes: '',
  })

  // Fetch inventory details
  const { data: inventory, isLoading } = useQuery({
    queryKey: ['inventory', inventoryId],
    queryFn: () => inventoryService.get(inventoryId),
  })

  // Fetch stock movements
  const { data: movements } = useQuery({
    queryKey: ['stock-movements', inventoryId],
    queryFn: () => inventoryService.getMovements(inventoryId),
    enabled: !!inventory,
  })

  // Fetch warehouses
  const { data: warehouses } = useQuery({
    queryKey: ['warehouses'],
    queryFn: () => inventoryService.listWarehouses(),
  })

  // Stock adjustment mutation
  const adjustStockMutation = useMutation({
    mutationFn: (data: StockAdjustmentRequest) => 
      inventoryService.adjustStock(inventoryId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['inventory', inventoryId] })
      queryClient.invalidateQueries({ queryKey: ['stock-movements', inventoryId] })
      toast({ title: '庫存調整成功' })
      setIsAdjustDialogOpen(false)
      setAdjustmentData({ quantity: 0, reason: '', notes: '' })
    },
    onError: (error: any) => {
      toast({
        title: '調整失敗',
        description: error.response?.data?.message || '庫存調整時發生錯誤',
        variant: 'destructive',
      })
    },
  })

  // Stock transfer mutation
  const transferStockMutation = useMutation({
    mutationFn: (data: StockTransferRequest) => 
      inventoryService.transferStock(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['inventory', inventoryId] })
      queryClient.invalidateQueries({ queryKey: ['stock-movements', inventoryId] })
      toast({ title: '庫存轉移成功' })
      setIsTransferDialogOpen(false)
      setTransferData({
        inventory_id: inventoryId,
        quantity: 0,
        from_warehouse_id: '',
        to_warehouse_id: '',
        notes: '',
      })
    },
    onError: (error: any) => {
      toast({
        title: '轉移失敗',
        description: error.response?.data?.message || '庫存轉移時發生錯誤',
        variant: 'destructive',
      })
    },
  })

  const getStockStatusBadge = (item: Inventory) => {
    if (item.current_stock <= 0) {
      return <Badge variant="destructive">缺貨</Badge>
    } else if (item.current_stock <= item.min_stock) {
      return <Badge variant="warning">低庫存</Badge>
    } else if (item.max_stock > 0 && item.current_stock >= item.max_stock) {
      return <Badge variant="info">超量</Badge>
    }
    return <Badge variant="success">正常</Badge>
  }

  const getMovementTypeBadge = (type: string, reason: string) => {
    const config: Record<string, { label: string; variant: any; icon: any }> = {
      in: { label: '入庫', variant: 'success', icon: TrendingUp },
      out: { label: '出庫', variant: 'destructive', icon: TrendingDown },
      adjustment: { label: '調整', variant: 'warning', icon: RefreshCw },
      transfer: { label: '轉移', variant: 'info', icon: ArrowUpDown },
    }

    const typeConfig = config[type] || { label: type, variant: 'default', icon: Package }
    const Icon = typeConfig.icon
    
    return (
      <Badge variant={typeConfig.variant as any} className="flex items-center gap-1">
        <Icon className="h-3 w-3" />
        {typeConfig.label}
      </Badge>
    )
  }

  const handleAdjustSubmit = () => {
    if (!adjustmentData.quantity || !adjustmentData.reason) {
      toast({
        title: '錯誤',
        description: '請填寫必要欄位',
        variant: 'destructive',
      })
      return
    }
    adjustStockMutation.mutate(adjustmentData)
  }

  const handleTransferSubmit = () => {
    if (!transferData.quantity || !transferData.from_warehouse_id || !transferData.to_warehouse_id) {
      toast({
        title: '錯誤',
        description: '請填寫必要欄位',
        variant: 'destructive',
      })
      return
    }
    transferStockMutation.mutate(transferData)
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

  if (!inventory) {
    return (
      <DashboardLayout>
        <div className="flex items-center justify-center h-64">
          <div className="text-center">找不到庫存品項</div>
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
              onClick={() => router.push('/inventory')}
            >
              <ArrowLeft className="h-4 w-4" />
            </Button>
            <div>
              <h1 className="text-3xl font-bold text-gray-900">庫存詳情</h1>
              <p className="mt-1 text-gray-600">
                SKU: {inventory.sku} | {inventory.name}
              </p>
            </div>
          </div>
          <div className="flex gap-2">
            <Button variant="outline" onClick={() => router.push(`/inventory/${inventoryId}/edit`)}>
              <Edit className="mr-2 h-4 w-4" />
              編輯
            </Button>
            <Button variant="outline" onClick={() => setIsAdjustDialogOpen(true)}>
              <RefreshCw className="mr-2 h-4 w-4" />
              調整庫存
            </Button>
            <Button variant="outline" onClick={() => setIsTransferDialogOpen(true)}>
              <ArrowUpDown className="mr-2 h-4 w-4" />
              轉移庫存
            </Button>
          </div>
        </div>

        <Tabs value={activeTab} onValueChange={setActiveTab}>
          <TabsList>
            <TabsTrigger value="details">基本資訊</TabsTrigger>
            <TabsTrigger value="stock">庫存狀態</TabsTrigger>
            <TabsTrigger value="movements">異動記錄</TabsTrigger>
            <TabsTrigger value="analytics">分析報表</TabsTrigger>
          </TabsList>

          <TabsContent value="details" className="space-y-6">
            {/* Basic Info */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center justify-between">
                  基本資訊
                  {getStockStatusBadge(inventory)}
                </CardTitle>
              </CardHeader>
              <CardContent className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div className="space-y-4">
                  <div>
                    <Label className="text-gray-500">SKU</Label>
                    <p className="font-medium">{inventory.sku}</p>
                  </div>
                  <div>
                    <Label className="text-gray-500">料號</Label>
                    <p className="font-medium">{inventory.part_no}</p>
                  </div>
                  <div>
                    <Label className="text-gray-500">品名</Label>
                    <p className="font-medium">{inventory.name}</p>
                  </div>
                  <div>
                    <Label className="text-gray-500">描述</Label>
                    <p>{inventory.description || '-'}</p>
                  </div>
                </div>
                <div className="space-y-4">
                  <div>
                    <Label className="text-gray-500">類別</Label>
                    <p className="font-medium">
                      {inventory.category === 'raw_material' && '原材料'}
                      {inventory.category === 'semi_finished' && '半成品'}
                      {inventory.category === 'finished_goods' && '成品'}
                    </p>
                  </div>
                  <div>
                    <Label className="text-gray-500">材質</Label>
                    <p>{inventory.material || '-'}</p>
                  </div>
                  <div>
                    <Label className="text-gray-500">規格</Label>
                    <p>{inventory.specification || '-'}</p>
                  </div>
                  <div>
                    <Label className="text-gray-500">單位</Label>
                    <p>{inventory.unit}</p>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Supplier Info */}
            {inventory.primary_supplier && (
              <Card>
                <CardHeader>
                  <CardTitle>供應商資訊</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                      <Label className="text-gray-500">主要供應商</Label>
                      <div className="flex items-center gap-2 mt-1">
                        <Building2 className="h-4 w-4 text-gray-400" />
                        <span>{inventory.primary_supplier.name}</span>
                      </div>
                    </div>
                    <div>
                      <Label className="text-gray-500">交貨期</Label>
                      <div className="flex items-center gap-2 mt-1">
                        <Truck className="h-4 w-4 text-gray-400" />
                        <span>{inventory.lead_time_days} 天</span>
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>
            )}
          </TabsContent>

          <TabsContent value="stock" className="space-y-6">
            {/* Stock Levels */}
            <Card>
              <CardHeader>
                <CardTitle>庫存水位</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                  <div>
                    <Label className="text-gray-500">現有庫存</Label>
                    <p className="text-2xl font-bold">{inventory.current_stock.toLocaleString()} {inventory.unit}</p>
                  </div>
                  <div>
                    <Label className="text-gray-500">可用庫存</Label>
                    <p className="text-2xl font-bold text-green-600">{inventory.available_stock.toLocaleString()} {inventory.unit}</p>
                  </div>
                  <div>
                    <Label className="text-gray-500">保留庫存</Label>
                    <p className="text-2xl font-bold text-orange-600">{inventory.reserved_stock.toLocaleString()} {inventory.unit}</p>
                  </div>
                </div>
                <Separator className="my-4" />
                <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
                  <div>
                    <Label className="text-gray-500">最低庫存</Label>
                    <p className="font-medium">{inventory.min_stock} {inventory.unit}</p>
                  </div>
                  <div>
                    <Label className="text-gray-500">最高庫存</Label>
                    <p className="font-medium">{inventory.max_stock || '-'} {inventory.unit}</p>
                  </div>
                  <div>
                    <Label className="text-gray-500">再訂購點</Label>
                    <p className="font-medium">{inventory.reorder_point} {inventory.unit}</p>
                  </div>
                  <div>
                    <Label className="text-gray-500">訂購量</Label>
                    <p className="font-medium">{inventory.reorder_quantity} {inventory.unit}</p>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Location Info */}
            <Card>
              <CardHeader>
                <CardTitle>儲位資訊</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <Label className="text-gray-500">倉庫</Label>
                    <div className="flex items-center gap-2 mt-1">
                      <WarehouseIcon className="h-4 w-4 text-gray-400" />
                      <span>{inventory.warehouse?.name || '-'}</span>
                    </div>
                  </div>
                  <div>
                    <Label className="text-gray-500">儲位</Label>
                    <p>{inventory.location || '-'}</p>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Cost Info */}
            <Card>
              <CardHeader>
                <CardTitle>成本資訊</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                  <div>
                    <Label className="text-gray-500">標準成本</Label>
                    <div className="flex items-center gap-1 mt-1">
                      <span className="text-xs text-gray-500">{inventory.currency}</span>
                      <span className="text-lg font-medium">{inventory.standard_cost.toFixed(2)}</span>
                    </div>
                  </div>
                  <div>
                    <Label className="text-gray-500">平均成本</Label>
                    <div className="flex items-center gap-1 mt-1">
                      <span className="text-xs text-gray-500">{inventory.currency}</span>
                      <span className="text-lg font-medium">{inventory.average_cost.toFixed(2)}</span>
                    </div>
                  </div>
                  <div>
                    <Label className="text-gray-500">最近採購價</Label>
                    <div className="flex items-center gap-1 mt-1">
                      <span className="text-xs text-gray-500">{inventory.currency}</span>
                      <span className="text-lg font-medium">{inventory.last_purchase_price.toFixed(2)}</span>
                    </div>
                  </div>
                </div>
                <Separator className="my-4" />
                <div>
                  <Label className="text-gray-500">庫存價值</Label>
                  <p className="text-2xl font-bold text-primary">
                    {inventory.currency} {(inventory.current_stock * inventory.average_cost).toLocaleString()}
                  </p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="movements" className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle>異動記錄</CardTitle>
                <CardDescription>最近的庫存異動記錄</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {movements?.length === 0 ? (
                    <p className="text-center text-gray-500 py-8">尚無異動記錄</p>
                  ) : (
                    movements?.map((movement) => (
                      <div key={movement.id} className="flex items-start justify-between p-4 border rounded-lg">
                        <div className="flex gap-4">
                          <div className="mt-1">
                            {getMovementTypeBadge(movement.movement_type, movement.reason)}
                          </div>
                          <div className="space-y-1">
                            <p className="font-medium">
                              {movement.reason === 'purchase' && '採購入庫'}
                              {movement.reason === 'sales' && '銷售出庫'}
                              {movement.reason === 'production' && '生產領用'}
                              {movement.reason === 'adjustment' && '庫存調整'}
                              {movement.reason === 'return' && '退貨'}
                              {movement.reason === 'damage' && '損壞'}
                              {movement.reason === 'transfer' && '庫存轉移'}
                              {movement.reason === 'initial' && '初始庫存'}
                              {movement.reason === 'found' && '盤盈'}
                              {movement.reason === 'loss' && '盤虧'}
                            </p>
                            <p className="text-sm text-gray-500">
                              數量: {movement.quantity > 0 ? '+' : ''}{movement.quantity} {inventory.unit} | 
                              庫存: {movement.before_quantity} → {movement.after_quantity}
                            </p>
                            {movement.reference_no && (
                              <p className="text-sm text-gray-500">
                                參考單號: {movement.reference_no}
                              </p>
                            )}
                            {movement.notes && (
                              <p className="text-sm text-gray-500">備註: {movement.notes}</p>
                            )}
                          </div>
                        </div>
                        <div className="text-right">
                          <p className="text-sm text-gray-500">
                            {format(new Date(movement.created_at), 'yyyy/MM/dd HH:mm', { locale: zhTW })}
                          </p>
                          <p className="text-sm text-gray-500">
                            {movement.creator?.full_name || '系統'}
                          </p>
                        </div>
                      </div>
                    ))
                  )}
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="analytics" className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle>庫存分析</CardTitle>
                <CardDescription>庫存趨勢與使用分析</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="text-center text-gray-500 py-12">
                  <BarChart3 className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p>分析報表功能開發中</p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>

        {/* Stock Adjustment Dialog */}
        <Dialog open={isAdjustDialogOpen} onOpenChange={setIsAdjustDialogOpen}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>調整庫存</DialogTitle>
              <DialogDescription>
                調整 {inventory.name} 的庫存數量
              </DialogDescription>
            </DialogHeader>
            <div className="space-y-4">
              <div>
                <Label>調整數量 *</Label>
                <Input
                  type="number"
                  value={adjustmentData.quantity}
                  onChange={(e) => setAdjustmentData({ ...adjustmentData, quantity: parseFloat(e.target.value) || 0 })}
                  placeholder="正數為增加，負數為減少"
                />
              </div>
              <div>
                <Label>調整原因 *</Label>
                <Select
                  value={adjustmentData.reason}
                  onValueChange={(value) => setAdjustmentData({ ...adjustmentData, reason: value })}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="選擇原因" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="damage">損壞</SelectItem>
                    <SelectItem value="loss">遺失</SelectItem>
                    <SelectItem value="found">盤盈</SelectItem>
                    <SelectItem value="correction">修正</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div>
                <Label>備註</Label>
                <Textarea
                  value={adjustmentData.notes}
                  onChange={(e) => setAdjustmentData({ ...adjustmentData, notes: e.target.value })}
                  placeholder="輸入調整說明..."
                  rows={3}
                />
              </div>
            </div>
            <DialogFooter>
              <Button variant="outline" onClick={() => setIsAdjustDialogOpen(false)}>
                取消
              </Button>
              <Button onClick={handleAdjustSubmit}>
                確認調整
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>

        {/* Stock Transfer Dialog */}
        <Dialog open={isTransferDialogOpen} onOpenChange={setIsTransferDialogOpen}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>轉移庫存</DialogTitle>
              <DialogDescription>
                在不同倉庫間轉移 {inventory.name}
              </DialogDescription>
            </DialogHeader>
            <div className="space-y-4">
              <div>
                <Label>轉移數量 *</Label>
                <Input
                  type="number"
                  value={transferData.quantity}
                  onChange={(e) => setTransferData({ ...transferData, quantity: parseFloat(e.target.value) || 0 })}
                  placeholder="輸入轉移數量"
                  min="0"
                />
              </div>
              <div>
                <Label>來源倉庫 *</Label>
                <Select
                  value={transferData.from_warehouse_id}
                  onValueChange={(value) => setTransferData({ ...transferData, from_warehouse_id: value })}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="選擇來源倉庫" />
                  </SelectTrigger>
                  <SelectContent>
                    {warehouses?.map((warehouse) => (
                      <SelectItem key={warehouse.id} value={warehouse.id}>
                        {warehouse.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <div>
                <Label>目標倉庫 *</Label>
                <Select
                  value={transferData.to_warehouse_id}
                  onValueChange={(value) => setTransferData({ ...transferData, to_warehouse_id: value })}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="選擇目標倉庫" />
                  </SelectTrigger>
                  <SelectContent>
                    {warehouses?.filter(w => w.id !== transferData.from_warehouse_id).map((warehouse) => (
                      <SelectItem key={warehouse.id} value={warehouse.id}>
                        {warehouse.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <div>
                <Label>備註</Label>
                <Textarea
                  value={transferData.notes}
                  onChange={(e) => setTransferData({ ...transferData, notes: e.target.value })}
                  placeholder="輸入轉移說明..."
                  rows={3}
                />
              </div>
            </div>
            <DialogFooter>
              <Button variant="outline" onClick={() => setIsTransferDialogOpen(false)}>
                取消
              </Button>
              <Button onClick={handleTransferSubmit}>
                確認轉移
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>
    </DashboardLayout>
  )
}