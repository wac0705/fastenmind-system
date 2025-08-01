'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
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
import { ArrowLeft, Save, Package, Warehouse, DollarSign, Truck } from 'lucide-react'
import { useToast } from '@/components/ui/use-toast'
import inventoryService, { CreateInventoryRequest } from '@/services/inventory.service'
import { Separator } from '@/components/ui/separator'

export default function NewInventoryPage() {
  const router = useRouter()
  const { toast } = useToast()

  const [formData, setFormData] = useState<CreateInventoryRequest>({
    sku: '',
    part_no: '',
    name: '',
    description: '',
    category: 'raw_material',
    material: '',
    specification: '',
    surface_treatment: '',
    heat_treatment: '',
    unit: 'PCS',
    initial_stock: 0,
    min_stock: 0,
    max_stock: 0,
    reorder_point: 0,
    reorder_quantity: 0,
    warehouse_id: '',
    location: '',
    standard_cost: 0,
    primary_supplier_id: '',
    lead_time_days: 0,
  })

  // Fetch warehouses
  const { data: warehouses } = useQuery({
    queryKey: ['warehouses'],
    queryFn: () => inventoryService.listWarehouses(),
  })

  // Create inventory mutation
  const createInventoryMutation = useMutation({
    mutationFn: (data: CreateInventoryRequest) => inventoryService.create(data),
    onSuccess: (inventory) => {
      toast({
        title: '成功',
        description: '庫存品項已建立',
      })
      router.push(`/inventory/${inventory.id}`)
    },
    onError: (error: any) => {
      toast({
        title: '錯誤',
        description: error.response?.data?.message || '建立庫存品項時發生錯誤',
        variant: 'destructive',
      })
    },
  })

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target
    setFormData(prev => ({
      ...prev,
      [name]: ['initial_stock', 'min_stock', 'max_stock', 'reorder_point', 'reorder_quantity', 'standard_cost', 'lead_time_days'].includes(name)
        ? parseFloat(value) || 0
        : value,
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
    
    if (!formData.sku || !formData.part_no || !formData.name || !formData.warehouse_id) {
      toast({
        title: '錯誤',
        description: '請填寫所有必要欄位',
        variant: 'destructive',
      })
      return
    }

    createInventoryMutation.mutate(formData)
  }

  // Set default warehouse if only one exists
  useEffect(() => {
    if (warehouses && warehouses.length === 1 && !formData.warehouse_id) {
      setFormData(prev => ({ ...prev, warehouse_id: warehouses[0].id }))
    }
  }, [warehouses])

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
            <h1 className="text-3xl font-bold text-gray-900">新增庫存品項</h1>
            <p className="mt-1 text-gray-600">建立新的庫存管理品項</p>
          </div>
        </div>

        <div className="space-y-6">
          {/* Basic Info */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Package className="h-5 w-5" />
                基本資訊
              </CardTitle>
              <CardDescription>
                設定品項的基本識別資訊
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="sku">SKU (庫存單位) *</Label>
                  <Input
                    id="sku"
                    name="sku"
                    value={formData.sku}
                    onChange={handleInputChange}
                    placeholder="輸入唯一的 SKU 編號"
                    required
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="part_no">料號 *</Label>
                  <Input
                    id="part_no"
                    name="part_no"
                    value={formData.part_no}
                    onChange={handleInputChange}
                    placeholder="輸入產品料號"
                    required
                  />
                </div>
              </div>
              <div className="space-y-2">
                <Label htmlFor="name">品名 *</Label>
                <Input
                  id="name"
                  name="name"
                  value={formData.name}
                  onChange={handleInputChange}
                  placeholder="輸入產品名稱"
                  required
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="description">描述</Label>
                <Textarea
                  id="description"
                  name="description"
                  value={formData.description}
                  onChange={handleInputChange}
                  placeholder="輸入產品描述..."
                  rows={3}
                />
              </div>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="category">類別 *</Label>
                  <Select
                    value={formData.category}
                    onValueChange={(value) => handleSelectChange('category', value)}
                  >
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="raw_material">原材料</SelectItem>
                      <SelectItem value="semi_finished">半成品</SelectItem>
                      <SelectItem value="finished_goods">成品</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                <div className="space-y-2">
                  <Label htmlFor="unit">單位</Label>
                  <Select
                    value={formData.unit}
                    onValueChange={(value) => handleSelectChange('unit', value)}
                  >
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="PCS">PCS (個)</SelectItem>
                      <SelectItem value="KG">KG (公斤)</SelectItem>
                      <SelectItem value="M">M (公尺)</SelectItem>
                      <SelectItem value="L">L (公升)</SelectItem>
                      <SelectItem value="BOX">BOX (箱)</SelectItem>
                      <SelectItem value="SET">SET (套)</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Specifications */}
          <Card>
            <CardHeader>
              <CardTitle>規格資訊</CardTitle>
              <CardDescription>
                設定產品的材質與處理規格
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="material">材質</Label>
                  <Input
                    id="material"
                    name="material"
                    value={formData.material}
                    onChange={handleInputChange}
                    placeholder="例如: SUS304, 45# 鋼"
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="specification">規格</Label>
                  <Input
                    id="specification"
                    name="specification"
                    value={formData.specification}
                    onChange={handleInputChange}
                    placeholder="例如: M8x50"
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="surface_treatment">表面處理</Label>
                  <Input
                    id="surface_treatment"
                    name="surface_treatment"
                    value={formData.surface_treatment}
                    onChange={handleInputChange}
                    placeholder="例如: 鍍鋅, 發黑"
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="heat_treatment">熱處理</Label>
                  <Input
                    id="heat_treatment"
                    name="heat_treatment"
                    value={formData.heat_treatment}
                    onChange={handleInputChange}
                    placeholder="例如: 調質, 淬火"
                  />
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Warehouse & Location */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Warehouse className="h-5 w-5" />
                倉儲資訊
              </CardTitle>
              <CardDescription>
                設定存放倉庫與儲位
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="warehouse_id">倉庫 *</Label>
                  <Select
                    value={formData.warehouse_id}
                    onValueChange={(value) => handleSelectChange('warehouse_id', value)}
                  >
                    <SelectTrigger>
                      <SelectValue placeholder="選擇倉庫" />
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
                <div className="space-y-2">
                  <Label htmlFor="location">儲位</Label>
                  <Input
                    id="location"
                    name="location"
                    value={formData.location}
                    onChange={handleInputChange}
                    placeholder="例如: A1-2-3"
                  />
                </div>
              </div>
              <div className="space-y-2">
                <Label htmlFor="initial_stock">初始庫存</Label>
                <Input
                  id="initial_stock"
                  name="initial_stock"
                  type="number"
                  value={formData.initial_stock}
                  onChange={handleInputChange}
                  min="0"
                  step="0.01"
                />
              </div>
            </CardContent>
          </Card>

          {/* Stock Control */}
          <Card>
            <CardHeader>
              <CardTitle>庫存控制</CardTitle>
              <CardDescription>
                設定庫存水位與補貨參數
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="min_stock">最低庫存</Label>
                  <Input
                    id="min_stock"
                    name="min_stock"
                    type="number"
                    value={formData.min_stock}
                    onChange={handleInputChange}
                    min="0"
                    step="0.01"
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="max_stock">最高庫存</Label>
                  <Input
                    id="max_stock"
                    name="max_stock"
                    type="number"
                    value={formData.max_stock}
                    onChange={handleInputChange}
                    min="0"
                    step="0.01"
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="reorder_point">再訂購點</Label>
                  <Input
                    id="reorder_point"
                    name="reorder_point"
                    type="number"
                    value={formData.reorder_point}
                    onChange={handleInputChange}
                    min="0"
                    step="0.01"
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="reorder_quantity">訂購批量</Label>
                  <Input
                    id="reorder_quantity"
                    name="reorder_quantity"
                    type="number"
                    value={formData.reorder_quantity}
                    onChange={handleInputChange}
                    min="0"
                    step="0.01"
                  />
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Cost & Supplier */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <DollarSign className="h-5 w-5" />
                成本與供應商
              </CardTitle>
              <CardDescription>
                設定成本資訊與供應商
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="standard_cost">標準成本</Label>
                  <Input
                    id="standard_cost"
                    name="standard_cost"
                    type="number"
                    value={formData.standard_cost}
                    onChange={handleInputChange}
                    min="0"
                    step="0.01"
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="lead_time_days">交貨期 (天)</Label>
                  <Input
                    id="lead_time_days"
                    name="lead_time_days"
                    type="number"
                    value={formData.lead_time_days}
                    onChange={handleInputChange}
                    min="0"
                  />
                </div>
              </div>
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
              disabled={createInventoryMutation.isPending}
            >
              <Save className="mr-2 h-4 w-4" />
              {createInventoryMutation.isPending ? '建立中...' : '建立品項'}
            </Button>
          </div>
        </div>
      </form>
    </DashboardLayout>
  )
}