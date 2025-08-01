'use client'

import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import DashboardLayout from '@/components/layout/DashboardLayout'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { Switch } from '@/components/ui/switch'
import { 
  Table, 
  TableBody, 
  TableCell, 
  TableHead, 
  TableHeader, 
  TableRow 
} from '@/components/ui/table'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { 
  Plus, 
  Edit, 
  Trash2, 
  Wrench,
  Zap,
  Home,
  DollarSign,
  Clock
} from 'lucide-react'
import { useToast } from '@/components/ui/use-toast'
import processService, { Equipment, ProcessCategory } from '@/services/process.service'

export default function EquipmentPage() {
  const queryClient = useQueryClient()
  const { toast } = useToast()
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false)
  const [editingEquipment, setEditingEquipment] = useState<Equipment | null>(null)
  const [formData, setFormData] = useState<Partial<Equipment>>({
    equipment_code: '',
    equipment_name: '',
    equipment_name_en: '',
    process_category_id: '',
    floor_area: 0,
    power_consumption: 0,
    max_capacity: 100,
    depreciation_years: 10,
    purchase_cost: 0,
    maintenance_cost_annual: 0,
    is_active: true,
  })

  // Fetch equipment list
  const { data: equipment = [], isLoading } = useQuery({
    queryKey: ['equipment'],
    queryFn: () => processService.listEquipment(),
  })

  // Fetch process categories
  const { data: categories = [] } = useQuery({
    queryKey: ['process-categories'],
    queryFn: () => processService.listCategories(),
  })

  // Create mutation
  const createMutation = useMutation({
    mutationFn: (data: Partial<Equipment>) => processService.createEquipment(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['equipment'] })
      toast({ title: '設備建立成功' })
      setIsCreateDialogOpen(false)
      resetForm()
    },
    onError: (error: any) => {
      toast({
        title: '建立失敗',
        description: error.response?.data?.message || '建立設備時發生錯誤',
        variant: 'destructive',
      })
    },
  })

  // Update mutation
  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<Equipment> }) =>
      processService.updateEquipment(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['equipment'] })
      toast({ title: '設備更新成功' })
      setEditingEquipment(null)
      setIsCreateDialogOpen(false)
      resetForm()
    },
    onError: (error: any) => {
      toast({
        title: '更新失敗',
        description: error.response?.data?.message || '更新設備時發生錯誤',
        variant: 'destructive',
      })
    },
  })

  const resetForm = () => {
    setFormData({
      equipment_code: '',
      equipment_name: '',
      equipment_name_en: '',
      process_category_id: '',
      floor_area: 0,
      power_consumption: 0,
      max_capacity: 100,
      depreciation_years: 10,
      purchase_cost: 0,
      maintenance_cost_annual: 0,
      is_active: true,
    })
  }

  const handleSubmit = () => {
    if (editingEquipment) {
      updateMutation.mutate({ id: editingEquipment.id, data: formData })
    } else {
      createMutation.mutate(formData)
    }
  }

  const handleEdit = (equipment: Equipment) => {
    setEditingEquipment(equipment)
    setFormData({
      equipment_code: equipment.equipment_code,
      equipment_name: equipment.equipment_name,
      equipment_name_en: equipment.equipment_name_en,
      process_category_id: equipment.process_category_id,
      floor_area: equipment.floor_area,
      power_consumption: equipment.power_consumption,
      max_capacity: equipment.max_capacity,
      depreciation_years: equipment.depreciation_years,
      purchase_cost: equipment.purchase_cost,
      maintenance_cost_annual: equipment.maintenance_cost_annual,
      is_active: equipment.is_active,
    })
    setIsCreateDialogOpen(true)
  }

  const getCategoryName = (categoryId: string) => {
    const category = categories.find(c => c.id === categoryId)
    return category?.category_name || '-'
  }

  return (
    <DashboardLayout>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex justify-between items-center">
          <div>
            <h1 className="text-3xl font-bold text-gray-900 flex items-center gap-2">
              <Wrench className="h-8 w-8" />
              設備管理
            </h1>
            <p className="mt-2 text-gray-600">管理生產設備資訊與成本參數</p>
          </div>
          <Button onClick={() => { resetForm(); setEditingEquipment(null); setIsCreateDialogOpen(true); }}>
            <Plus className="mr-2 h-4 w-4" />
            新增設備
          </Button>
        </div>

        {/* Equipment List */}
        <Card>
          <CardHeader>
            <CardTitle>設備列表</CardTitle>
            <CardDescription>
              共 {equipment.length} 台設備
            </CardDescription>
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <div className="text-center py-8">載入中...</div>
            ) : (
              <div className="overflow-x-auto">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>設備代碼</TableHead>
                      <TableHead>設備名稱</TableHead>
                      <TableHead>製程類別</TableHead>
                      <TableHead>最大產能</TableHead>
                      <TableHead>耗電功率</TableHead>
                      <TableHead>佔地面積</TableHead>
                      <TableHead>狀態</TableHead>
                      <TableHead className="text-right">操作</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {equipment.map((eq) => (
                      <TableRow key={eq.id}>
                        <TableCell className="font-medium">{eq.equipment_code}</TableCell>
                        <TableCell>
                          <div>
                            <p className="font-medium">{eq.equipment_name}</p>
                            {eq.equipment_name_en && (
                              <p className="text-sm text-gray-500">{eq.equipment_name_en}</p>
                            )}
                          </div>
                        </TableCell>
                        <TableCell>{getCategoryName(eq.process_category_id)}</TableCell>
                        <TableCell>{eq.max_capacity} 件/時</TableCell>
                        <TableCell>{eq.power_consumption} kW</TableCell>
                        <TableCell>{eq.floor_area} m²</TableCell>
                        <TableCell>
                          <Badge variant={eq.is_active ? 'success' : 'secondary'}>
                            {eq.is_active ? '啟用' : '停用'}
                          </Badge>
                        </TableCell>
                        <TableCell className="text-right">
                          <div className="flex justify-end gap-2">
                            <Button
                              variant="ghost"
                              size="icon"
                              onClick={() => handleEdit(eq)}
                            >
                              <Edit className="h-4 w-4" />
                            </Button>
                          </div>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Create/Edit Dialog */}
        <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
          <DialogContent className="sm:max-w-[600px]">
            <DialogHeader>
              <DialogTitle>{editingEquipment ? '編輯設備' : '新增設備'}</DialogTitle>
              <DialogDescription>
                填寫設備資訊與成本參數
              </DialogDescription>
            </DialogHeader>
            <div className="grid gap-4 py-4 max-h-[60vh] overflow-y-auto">
              <div className="grid grid-cols-2 gap-4">
                <div className="grid gap-2">
                  <Label htmlFor="equipment_code">設備代碼 *</Label>
                  <Input
                    id="equipment_code"
                    value={formData.equipment_code}
                    onChange={(e) => setFormData({ ...formData, equipment_code: e.target.value })}
                    placeholder="例如: FM-001"
                  />
                </div>
                <div className="grid gap-2">
                  <Label htmlFor="process_category_id">製程類別 *</Label>
                  <Select
                    value={formData.process_category_id}
                    onValueChange={(value) => setFormData({ ...formData, process_category_id: value })}
                  >
                    <SelectTrigger>
                      <SelectValue placeholder="選擇製程類別" />
                    </SelectTrigger>
                    <SelectContent>
                      {categories.map((category) => (
                        <SelectItem key={category.id} value={category.id}>
                          {category.category_name}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div className="grid gap-2">
                  <Label htmlFor="equipment_name">設備名稱 (中文) *</Label>
                  <Input
                    id="equipment_name"
                    value={formData.equipment_name}
                    onChange={(e) => setFormData({ ...formData, equipment_name: e.target.value })}
                    placeholder="例如: 六模六衝成型機"
                  />
                </div>
                <div className="grid gap-2">
                  <Label htmlFor="equipment_name_en">設備名稱 (英文)</Label>
                  <Input
                    id="equipment_name_en"
                    value={formData.equipment_name_en}
                    onChange={(e) => setFormData({ ...formData, equipment_name_en: e.target.value })}
                    placeholder="例如: 6-Die 6-Blow Former"
                  />
                </div>
              </div>

              <div className="grid grid-cols-3 gap-4">
                <div className="grid gap-2">
                  <Label htmlFor="floor_area">
                    <Home className="inline h-4 w-4 mr-1" />
                    佔地面積 (m²)
                  </Label>
                  <Input
                    id="floor_area"
                    type="number"
                    value={formData.floor_area}
                    onChange={(e) => setFormData({ ...formData, floor_area: parseFloat(e.target.value) || 0 })}
                    step="0.1"
                    min="0"
                  />
                </div>
                <div className="grid gap-2">
                  <Label htmlFor="power_consumption">
                    <Zap className="inline h-4 w-4 mr-1" />
                    耗電功率 (kW)
                  </Label>
                  <Input
                    id="power_consumption"
                    type="number"
                    value={formData.power_consumption}
                    onChange={(e) => setFormData({ ...formData, power_consumption: parseFloat(e.target.value) || 0 })}
                    step="0.1"
                    min="0"
                  />
                </div>
                <div className="grid gap-2">
                  <Label htmlFor="max_capacity">
                    <Clock className="inline h-4 w-4 mr-1" />
                    最大產能 (件/時)
                  </Label>
                  <Input
                    id="max_capacity"
                    type="number"
                    value={formData.max_capacity}
                    onChange={(e) => setFormData({ ...formData, max_capacity: parseInt(e.target.value) || 0 })}
                    min="1"
                  />
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div className="grid gap-2">
                  <Label htmlFor="purchase_cost">
                    <DollarSign className="inline h-4 w-4 mr-1" />
                    購置成本
                  </Label>
                  <Input
                    id="purchase_cost"
                    type="number"
                    value={formData.purchase_cost}
                    onChange={(e) => setFormData({ ...formData, purchase_cost: parseFloat(e.target.value) || 0 })}
                    step="1000"
                    min="0"
                  />
                </div>
                <div className="grid gap-2">
                  <Label htmlFor="depreciation_years">折舊年限</Label>
                  <Input
                    id="depreciation_years"
                    type="number"
                    value={formData.depreciation_years}
                    onChange={(e) => setFormData({ ...formData, depreciation_years: parseInt(e.target.value) || 10 })}
                    min="1"
                    max="30"
                  />
                </div>
              </div>

              <div className="grid gap-2">
                <Label htmlFor="maintenance_cost_annual">年度維護成本</Label>
                <Input
                  id="maintenance_cost_annual"
                  type="number"
                  value={formData.maintenance_cost_annual}
                  onChange={(e) => setFormData({ ...formData, maintenance_cost_annual: parseFloat(e.target.value) || 0 })}
                  step="1000"
                  min="0"
                />
              </div>

              <div className="flex items-center space-x-2">
                <Switch
                  id="is_active"
                  checked={formData.is_active}
                  onCheckedChange={(checked) => setFormData({ ...formData, is_active: checked })}
                />
                <Label htmlFor="is_active">啟用設備</Label>
              </div>
            </div>
            <DialogFooter>
              <Button variant="outline" onClick={() => setIsCreateDialogOpen(false)}>
                取消
              </Button>
              <Button onClick={handleSubmit}>
                {editingEquipment ? '更新' : '建立'}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>
    </DashboardLayout>
  )
}