'use client'

import { useState } from 'react'
import { useQuery, useMutation } from '@tanstack/react-query'
import DashboardLayout from '@/components/layout/DashboardLayout'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { 
  Calculator, 
  Plus, 
  Trash2, 
  Settings,
  Clock,
  DollarSign,
  Zap,
  Factory,
  Users,
  FileText,
  Loader2
} from 'lucide-react'
import costCalculationService, { 
  ProcessStep, 
  Equipment, 
  ProductProcessRoute, 
  CostCalculationRequest,
  CostSummary 
} from '@/services/cost-calculation.service'
import { useToast } from '@/components/ui/use-toast'

interface CustomProcessStep {
  id: string
  process_step_id: string
  equipment_id: string
  setup_time?: number
  cycle_time?: number
  process_step?: ProcessStep
  equipment?: Equipment
}

export default function CostCalculatorPage() {
  const { toast } = useToast()
  const [customSteps, setCustomSteps] = useState<CustomProcessStep[]>([])
  const [productName, setProductName] = useState('')
  const [productCategory, setProductCategory] = useState('')
  const [materialType, setMaterialType] = useState('')
  const [quantity, setQuantity] = useState(1000)
  const [materialCost, setMaterialCost] = useState(0)
  const [marginPercentage, setMarginPercentage] = useState(30)
  const [selectedRouteId, setSelectedRouteId] = useState<string>('')
  const [useCustomRoute, setUseCustomRoute] = useState(false)
  const [lastCalculation, setLastCalculation] = useState<CostSummary | null>(null)

  // Fetch process steps
  const { data: processSteps = [] } = useQuery({
    queryKey: ['process-steps'],
    queryFn: () => costCalculationService.getProcessSteps(),
  })

  // Fetch equipment
  const { data: equipment = [] } = useQuery({
    queryKey: ['equipment'],
    queryFn: () => costCalculationService.getEquipment(),
  })

  // Fetch process routes
  const { data: processRoutes = [] } = useQuery({
    queryKey: ['process-routes', productCategory],
    queryFn: () => costCalculationService.getProcessRoutes(productCategory),
    enabled: !!productCategory,
  })

  const addCustomStep = () => {
    const newStep: CustomProcessStep = {
      id: `step-${Date.now()}`,
      process_step_id: '',
      equipment_id: '',
    }
    setCustomSteps([...customSteps, newStep])
  }

  const removeCustomStep = (id: string) => {
    setCustomSteps(customSteps.filter(step => step.id !== id))
  }

  const updateCustomStep = (id: string, updates: Partial<CustomProcessStep>) => {
    setCustomSteps(customSteps.map(step => 
      step.id === id ? { ...step, ...updates } : step
    ))
  }

  // 成本計算 mutation
  const calculateMutation = useMutation({
    mutationFn: async () => {
      const request: CostCalculationRequest = {
        product_name: productName,
        product_category: productCategory,
        material_type: materialType,
        quantity: quantity,
        material_cost: materialCost,
        margin_percentage: marginPercentage,
      }

      if (useCustomRoute && customSteps.length > 0) {
        request.custom_route = customSteps
          .filter(step => step.process_step_id && step.equipment_id)
          .map(step => ({
            process_step_id: step.process_step_id,
            equipment_id: step.equipment_id,
            setup_time: step.setup_time,
            cycle_time: step.cycle_time,
          }))
      } else if (selectedRouteId) {
        request.route_id = selectedRouteId
      }

      const calculation = await costCalculationService.calculateCost(request)
      const summary = await costCalculationService.getCostSummary(calculation.id)
      return summary
    },
    onSuccess: (summary) => {
      setLastCalculation(summary)
      toast({
        title: '計算完成',
        description: `總成本: $${summary.total_cost.toFixed(2)}, 單位成本: $${summary.unit_cost.toFixed(4)}`,
      })
    },
    onError: (error) => {
      toast({
        title: '計算失敗',
        description: error.message || '計算成本時發生錯誤',
        variant: 'destructive',
      })
    },
  })

  const getEquipmentForStep = (stepId: string) => {
    const step = processSteps.find(s => s.id === stepId)
    if (!step) return equipment
    return equipment.filter(eq => eq.process_category_id === step.process_category_id)
  }

  const productCategories = ['screws', 'nuts', 'washers', 'bolts', 'special', 'custom']
  const materialTypes = ['碳鋼', '不鏽鋼', '黃銅', '銅', '塑膐']

  return (
    <DashboardLayout>
      <div className="space-y-6">
        {/* Header */}
        <div>
          <h1 className="text-3xl font-bold text-gray-900 flex items-center gap-2">
            <Calculator className="h-8 w-8" />
            製程成本計算器
          </h1>
          <p className="mt-2 text-gray-600">
            選擇製程路線或自定義製程，計算產品的製造成本
          </p>
        </div>

        {/* Product Information */}
        <Card>
          <CardHeader>
            <CardTitle>產品資訊</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              <div>
                <Label htmlFor="product-name">產品名稱</Label>
                <Input
                  id="product-name"
                  value={productName}
                  onChange={(e) => setProductName(e.target.value)}
                  placeholder="例如：M10x30 六角螺栓"
                />
              </div>
              <div>
                <Label htmlFor="product-category">產品類別</Label>
                <Select value={productCategory} onValueChange={setProductCategory}>
                  <SelectTrigger>
                    <SelectValue placeholder="選擇類別" />
                  </SelectTrigger>
                  <SelectContent>
                    {productCategories.map((cat) => (
                      <SelectItem key={cat} value={cat}>{cat}</SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <div>
                <Label htmlFor="material-type">材質</Label>
                <Select value={materialType} onValueChange={setMaterialType}>
                  <SelectTrigger>
                    <SelectValue placeholder="選擇材質" />
                  </SelectTrigger>
                  <SelectContent>
                    {materialTypes.map((type) => (
                      <SelectItem key={type} value={type}>{type}</SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <div>
                <Label htmlFor="quantity">數量</Label>
                <Input
                  id="quantity"
                  type="number"
                  value={quantity}
                  onChange={(e) => setQuantity(parseInt(e.target.value) || 1)}
                  min="1"
                />
              </div>
              <div>
                <Label htmlFor="material-cost">材料成本 (USD)</Label>
                <Input
                  id="material-cost"
                  type="number"
                  step="0.01"
                  value={materialCost}
                  onChange={(e) => setMaterialCost(parseFloat(e.target.value) || 0)}
                  min="0"
                />
              </div>
              <div>
                <Label htmlFor="margin">毛利率 (%)</Label>
                <Input
                  id="margin"
                  type="number"
                  value={marginPercentage}
                  onChange={(e) => setMarginPercentage(parseFloat(e.target.value) || 0)}
                  min="0"
                  max="100"
                />
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Process Route Selection */}
        <Tabs defaultValue="standard" onValueChange={(value) => setUseCustomRoute(value === 'custom')}>
          <TabsList className="grid w-full grid-cols-2">
            <TabsTrigger value="standard">標準製程路線</TabsTrigger>
            <TabsTrigger value="custom">自定義製程</TabsTrigger>
          </TabsList>
          
          <TabsContent value="standard">
            <Card>
              <CardHeader>
                <CardTitle>選擇製程路線</CardTitle>
                <CardDescription>選擇預先定義的標準製程路線</CardDescription>
              </CardHeader>
              <CardContent>
                {processRoutes.length === 0 ? (
                  <p className="text-sm text-gray-500">請先選擇產品類別</p>
                ) : (
                  <div className="space-y-2">
                    {processRoutes.map((route) => (
                      <div key={route.id} className="flex items-center space-x-2 p-3 border rounded-lg">
                        <input
                          type="radio"
                          id={`route-${route.id}`}
                          name="route"
                          value={route.id}
                          checked={selectedRouteId === route.id}
                          onChange={(e) => setSelectedRouteId(e.target.value)}
                        />
                        <label htmlFor={`route-${route.id}`} className="flex-1 cursor-pointer">
                          <div className="font-medium">{route.route_name}</div>
                          <div className="text-sm text-gray-500">
                            {route.route_details.length} 個製程步驟
                            {route.is_default && <Badge variant="secondary" className="ml-2">預設</Badge>}
                          </div>
                        </label>
                      </div>
                    ))}
                  </div>
                )}
              </CardContent>
            </Card>
          </TabsContent>
          
          <TabsContent value="custom">
            <Card>
              <CardHeader className="flex flex-row items-center justify-between">
                <div>
                  <CardTitle>自定義製程步驟</CardTitle>
                  <CardDescription>按順序添加製程步驟</CardDescription>
                </div>
                <Button onClick={addCustomStep} size="sm">
                  <Plus className="mr-2 h-4 w-4" />
                  添加步驟
                </Button>
              </CardHeader>
              <CardContent>
            <div className="space-y-4">
              {customSteps.length === 0 ? (
                <div className="text-center py-8 text-gray-500">
                  尚未添加製程步驟，點擊上方按鈕開始
                </div>
              ) : (
                customSteps.map((step, index) => (
                  <Card key={step.id} className="p-4">
                    <div className="space-y-4">
                      <div className="flex items-center justify-between">
                        <h4 className="font-medium">步驟 {index + 1}</h4>
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => removeCustomStep(step.id)}
                        >
                          <Trash2 className="h-4 w-4" />
                        </Button>
                      </div>

                      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <div>
                          <Label>製程步驟</Label>
                          <Select
                            value={step.process_step_id}
                            onValueChange={(value) => updateCustomStep(step.id, { 
                              process_step_id: value,
                              equipment_id: '' // Reset equipment when step changes
                            })}
                          >
                            <SelectTrigger>
                              <SelectValue placeholder="選擇製程" />
                            </SelectTrigger>
                            <SelectContent>
                              {processSteps.map((ps) => (
                                <SelectItem key={ps.id} value={ps.id}>
                                  {ps.name} {ps.process_category && `(${ps.process_category.name})`}
                                </SelectItem>
                              ))}
                            </SelectContent>
                          </Select>
                        </div>

                        <div>
                          <Label>設備</Label>
                          <Select
                            value={step.equipment_id}
                            onValueChange={(value) => updateCustomStep(step.id, { equipment_id: value })}
                            disabled={!step.process_step_id}
                          >
                            <SelectTrigger>
                              <SelectValue placeholder="選擇設備" />
                            </SelectTrigger>
                            <SelectContent>
                              {getEquipmentForStep(step.process_step_id).map((eq) => (
                                <SelectItem key={eq.id} value={eq.id}>
                                  {eq.name}
                                </SelectItem>
                              ))}
                            </SelectContent>
                          </Select>
                        </div>
                      </div>

                      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <div>
                          <Label>設置時間 (分鐘) - 選填</Label>
                          <Input
                            type="number"
                            value={step.setup_time || ''}
                            onChange={(e) => updateCustomStep(step.id, { 
                              setup_time: parseFloat(e.target.value) || undefined 
                            })}
                            placeholder="使用預設值"
                            min="0"
                          />
                        </div>

                        <div>
                          <Label>單件週期時間 (秒) - 選填</Label>
                          <Input
                            type="number"
                            value={step.cycle_time || ''}
                            onChange={(e) => updateCustomStep(step.id, { 
                              cycle_time: parseFloat(e.target.value) || undefined 
                            })}
                            placeholder="使用預設值"
                            min="0"
                          />
                        </div>
                      </div>
                    </div>
                  </Card>
                ))
              )}
            </div>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>

        {/* Calculate Button */}
        <div className="flex justify-center">
          <Button 
            onClick={() => calculateMutation.mutate()} 
            disabled={calculateMutation.isPending || (!selectedRouteId && customSteps.length === 0) || !productName || !productCategory}
            size="lg"
          >
            {calculateMutation.isPending ? (
              <><Loader2 className="mr-2 h-5 w-5 animate-spin" /> 計算中...</>
            ) : (
              <><Calculator className="mr-2 h-5 w-5" /> 計算成本</>
            )}
          </Button>
        </div>

        {/* Calculation Results */}
        {lastCalculation && (
          <Card>
            <CardHeader>
              <CardTitle>成本計算結果</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
                <Card>
                  <CardContent className="pt-6">
                    <div className="text-2xl font-bold">${lastCalculation.material_cost.toFixed(2)}</div>
                    <p className="text-sm text-muted-foreground">材料成本</p>
                  </CardContent>
                </Card>
                <Card>
                  <CardContent className="pt-6">
                    <div className="text-2xl font-bold">${lastCalculation.process_cost.toFixed(2)}</div>
                    <p className="text-sm text-muted-foreground">加工成本</p>
                  </CardContent>
                </Card>
                <Card>
                  <CardContent className="pt-6">
                    <div className="text-2xl font-bold">${lastCalculation.overhead_cost.toFixed(2)}</div>
                    <p className="text-sm text-muted-foreground">管理費用</p>
                  </CardContent>
                </Card>
                <Card>
                  <CardContent className="pt-6">
                    <div className="text-2xl font-bold text-green-600">${lastCalculation.total_cost.toFixed(2)}</div>
                    <p className="text-sm text-muted-foreground">總成本</p>
                  </CardContent>
                </Card>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
                <Card>
                  <CardContent className="pt-6">
                    <div className="flex justify-between items-center">
                      <div>
                        <p className="text-sm text-muted-foreground">單位成本</p>
                        <p className="text-2xl font-bold">${lastCalculation.unit_cost.toFixed(4)}</p>
                      </div>
                      <DollarSign className="h-8 w-8 text-gray-400" />
                    </div>
                  </CardContent>
                </Card>
                <Card>
                  <CardContent className="pt-6">
                    <div className="flex justify-between items-center">
                      <div>
                        <p className="text-sm text-muted-foreground">建議售價</p>
                        <p className="text-2xl font-bold text-blue-600">${lastCalculation.suggested_price.toFixed(2)}</p>
                        <p className="text-xs text-gray-500">毛利率: {lastCalculation.margin_percentage}%</p>
                      </div>
                      <FileText className="h-8 w-8 text-gray-400" />
                    </div>
                  </CardContent>
                </Card>
              </div>

              {/* Process Breakdown */}
              <div>
                <h3 className="text-lg font-semibold mb-4">製程成本明細</h3>
                <div className="space-y-2">
                  {lastCalculation.process_breakdown.map((process, index) => (
                    <div key={index} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                      <div className="flex-1">
                        <div className="font-medium">{process.process_name}</div>
                        {process.equipment_name && (
                          <div className="text-sm text-gray-500">設備: {process.equipment_name}</div>
                        )}
                      </div>
                      <div className="text-right">
                        <div className="font-medium">${process.total_cost.toFixed(2)}</div>
                        <div className="text-xs text-gray-500">
                          工時: {process.total_time_hours.toFixed(2)} 小時
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            </CardContent>
          </Card>
        )}
      </div>
    </DashboardLayout>
  )
}