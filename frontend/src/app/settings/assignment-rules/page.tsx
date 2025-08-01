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
  DialogTrigger,
} from '@/components/ui/dialog'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Plus, Edit, Trash2, Users, Package, ArrowUpDown } from 'lucide-react'
import { useToast } from '@/components/ui/use-toast'
import assignmentService, { AssignmentRule } from '@/services/assignment.service'

export default function AssignmentRulesPage() {
  const queryClient = useQueryClient()
  const { toast } = useToast()
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false)
  const [editingRule, setEditingRule] = useState<AssignmentRule | null>(null)
  const [formData, setFormData] = useState<Partial<AssignmentRule>>({
    rule_name: '',
    rule_type: 'load_balance',
    priority: 100,
    conditions: {
      product_categories: [],
      auto_assign: true,
    },
    is_active: true,
  })

  // Fetch assignment rules
  const { data: rules = [], isLoading } = useQuery({
    queryKey: ['assignment-rules'],
    queryFn: () => assignmentService.getAssignmentRules(),
  })

  // Fetch workload stats for available engineers
  const { data: workloadStats = [] } = useQuery({
    queryKey: ['engineer-workload'],
    queryFn: () => assignmentService.getWorkloadStats(),
  })

  // Fetch customers (mock data for now)
  const customers = [
    { id: 'cust-001', name: '台灣機械股份有限公司', code: 'CUST-001' },
    { id: 'cust-002', name: 'Global Auto Parts GmbH', code: 'CUST-002' },
    { id: 'cust-003', name: 'American Tools Inc.', code: 'CUST-003' },
  ]

  // Product categories
  const productCategories = ['screws', 'nuts', 'washers', 'bolts', 'special', 'custom']

  // Create/Update mutation
  const saveMutation = useMutation({
    mutationFn: async (data: Partial<AssignmentRule>) => {
      if (editingRule?.id) {
        return assignmentService.updateAssignmentRule(editingRule.id, data)
      }
      // For new rules, we'll need to implement a create API
      throw new Error('Create API not implemented yet')
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['assignment-rules'] })
      toast({ title: editingRule ? '規則更新成功' : '規則建立成功' })
      setIsCreateDialogOpen(false)
      setEditingRule(null)
      resetForm()
    },
    onError: (error: any) => {
      toast({
        title: editingRule ? '更新失敗' : '建立失敗',
        description: error.message || '操作規則時發生錯誤',
        variant: 'destructive',
      })
    },
  })


  const resetForm = () => {
    setFormData({
      rule_name: '',
      rule_type: 'load_balance',
      priority: 100,
      conditions: {
        product_categories: [],
        auto_assign: true,
      },
      is_active: true,
    })
  }

  const handleSubmit = () => {
    saveMutation.mutate(formData)
  }

  const handleEdit = (rule: AssignmentRule) => {
    setEditingRule(rule)
    setFormData({
      rule_name: rule.rule_name,
      rule_type: rule.rule_type,
      priority: rule.priority,
      conditions: rule.conditions,
      is_active: rule.is_active,
    })
    setIsCreateDialogOpen(true)
  }


  const getRuleTypeIcon = (type: string) => {
    switch (type) {
      case 'by_customer':
        return <Users className="h-4 w-4" />
      case 'by_product':
        return <Package className="h-4 w-4" />
      case 'by_both':
        return <ArrowUpDown className="h-4 w-4" />
      default:
        return null
    }
  }

  const getRuleTypeLabel = (type: string) => {
    switch (type) {
      case 'auto':
        return '自動分派'
      case 'rotation':
        return '輪流分派'
      case 'load_balance':
        return '負載平衡'
      case 'skill_based':
        return '技能匹配'
      default:
        return type
    }
  }

  return (
    <DashboardLayout>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex justify-between items-center">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">工程師分派規則</h1>
            <p className="mt-2 text-gray-600">設定詢價單自動分派給工程師的規則</p>
          </div>
          <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
            <DialogTrigger asChild>
              <Button onClick={() => { resetForm(); setEditingRule(null); }}>
                <Plus className="mr-2 h-4 w-4" />
                新增規則
              </Button>
            </DialogTrigger>
            <DialogContent className="sm:max-w-[500px]">
              <DialogHeader>
                <DialogTitle>{editingRule ? '編輯規則' : '新增分派規則'}</DialogTitle>
                <DialogDescription>
                  設定詢價單自動分派的條件和指定工程師
                </DialogDescription>
              </DialogHeader>
              <div className="grid gap-4 py-4">
                <div className="grid gap-2">
                  <Label htmlFor="rule_name">規則名稱</Label>
                  <Input
                    id="rule_name"
                    value={formData.rule_name || ''}
                    onChange={(e) => setFormData({ ...formData, rule_name: e.target.value })}
                    placeholder="例如：標準件自動分派"
                  />
                </div>

                <div className="grid gap-2">
                  <Label htmlFor="rule_type">規則類型</Label>
                  <Select
                    value={formData.rule_type}
                    onValueChange={(value: any) => setFormData({ ...formData, rule_type: value })}
                  >
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="auto">自動分派</SelectItem>
                      <SelectItem value="rotation">輪流分派</SelectItem>
                      <SelectItem value="load_balance">負載平衡</SelectItem>
                      <SelectItem value="skill_based">技能匹配</SelectItem>
                    </SelectContent>
                  </Select>
                </div>

                <div className="grid gap-2">
                  <Label>產品類別條件</Label>
                  <div className="space-y-2">
                    {productCategories.map((category) => (
                      <div key={category} className="flex items-center space-x-2">
                        <input
                          type="checkbox"
                          id={`cat-${category}`}
                          checked={formData.conditions?.product_categories?.includes(category) || false}
                          onChange={(e) => {
                            const categories = formData.conditions?.product_categories || []
                            if (e.target.checked) {
                              setFormData({
                                ...formData,
                                conditions: {
                                  ...formData.conditions,
                                  product_categories: [...categories, category]
                                }
                              })
                            } else {
                              setFormData({
                                ...formData,
                                conditions: {
                                  ...formData.conditions,
                                  product_categories: categories.filter(c => c !== category)
                                }
                              })
                            }
                          }}
                        />
                        <Label htmlFor={`cat-${category}`}>{category}</Label>
                      </div>
                    ))}
                  </div>
                </div>

                {formData.rule_type === 'skill_based' && (
                  <div className="grid gap-2">
                    <Label htmlFor="min_skill_level">最低技能等級</Label>
                    <Input
                      id="min_skill_level"
                      type="number"
                      min="1"
                      max="5"
                      value={formData.conditions?.min_skill_level || 1}
                      onChange={(e) => setFormData({
                        ...formData,
                        conditions: {
                          ...formData.conditions,
                          min_skill_level: parseInt(e.target.value)
                        }
                      })}
                    />
                    <p className="text-sm text-gray-500">1 (初級) - 5 (專家)</p>
                  </div>
                )}

                <div className="grid gap-2">
                  <Label htmlFor="priority">優先順序</Label>
                  <Input
                    id="priority"
                    type="number"
                    value={formData.priority}
                    onChange={(e) => setFormData({ ...formData, priority: parseInt(e.target.value) })}
                    min="1"
                    max="999"
                  />
                  <p className="text-sm text-gray-500">數字越小優先級越高</p>
                </div>

                <div className="flex items-center space-x-2">
                  <Switch
                    id="is_active"
                    checked={formData.is_active}
                    onCheckedChange={(checked) => setFormData({ ...formData, is_active: checked })}
                  />
                  <Label htmlFor="is_active">啟用規則</Label>
                </div>
              </div>
              <DialogFooter>
                <Button variant="outline" onClick={() => setIsCreateDialogOpen(false)}>
                  取消
                </Button>
                <Button onClick={handleSubmit}>
                  {editingRule ? '更新' : '建立'}
                </Button>
              </DialogFooter>
            </DialogContent>
          </Dialog>
        </div>

        {/* Rules List */}
        <Card>
          <CardHeader>
            <CardTitle>分派規則列表</CardTitle>
            <CardDescription>
              管理詢價單自動分派的規則，規則會依優先順序執行
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
                      <TableHead>規則名稱</TableHead>
                      <TableHead>類型</TableHead>
                      <TableHead>條件</TableHead>
                      <TableHead>優先順序</TableHead>
                      <TableHead>狀態</TableHead>
                      <TableHead className="text-right">操作</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {rules.map((rule) => (
                      <TableRow key={rule.id}>
                        <TableCell className="font-medium">{rule.rule_name}</TableCell>
                        <TableCell>
                          <div className="flex items-center gap-2">
                            {getRuleTypeIcon(rule.rule_type)}
                            <Badge variant="outline">{getRuleTypeLabel(rule.rule_type)}</Badge>
                          </div>
                        </TableCell>
                        <TableCell>
                          <div className="text-sm space-y-1">
                            {rule.conditions.product_categories && rule.conditions.product_categories.length > 0 && (
                              <div>
                                <span className="font-medium">產品：</span>
                                {rule.conditions.product_categories.join(', ')}
                              </div>
                            )}
                            {rule.conditions.min_skill_level && (
                              <div>
                                <span className="font-medium">最低技能：</span>
                                Lv.{rule.conditions.min_skill_level}
                              </div>
                            )}
                            {rule.conditions.auto_assign && (
                              <Badge variant="secondary" className="text-xs">自動指派</Badge>
                            )}
                          </div>
                        </TableCell>
                        <TableCell>{rule.priority}</TableCell>
                        <TableCell>
                          <Switch
                            checked={rule.is_active}
                            onCheckedChange={async (checked) => {
                              try {
                                await assignmentService.updateAssignmentRule(rule.id, {
                                  is_active: checked
                                })
                                queryClient.invalidateQueries({ queryKey: ['assignment-rules'] })
                                toast({ title: `規則已${checked ? '啟用' : '停用'}` })
                              } catch (error) {
                                toast({ 
                                  title: '操作失敗',
                                  variant: 'destructive'
                                })
                              }
                            }}
                          />
                        </TableCell>
                        <TableCell className="text-right">
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => handleEdit(rule)}
                          >
                            <Edit className="h-4 w-4" />
                          </Button>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Usage Tips */}
        <Card>
          <CardHeader>
            <CardTitle>使用說明</CardTitle>
          </CardHeader>
          <CardContent className="space-y-2 text-sm text-gray-600">
            <p>• 系統會依據優先順序（數字越小越優先）執行規則</p>
            <p>• 當有多條規則符合時，會選擇優先順序最高的規則</p>
            <p>• 技能匹配規則會根據工程師的技能等級分派</p>
            <p>• 負載平衡會選擇當前工作量最少的工程師</p>
            <p>• 停用的規則不會被執行</p>
            <p>• 如果沒有符合的規則，詢價單會進入待分派狀態</p>
          </CardContent>
        </Card>
      </div>
    </DashboardLayout>
  )
}