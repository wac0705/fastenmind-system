'use client'

import { useState, useEffect } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import { useMutation, useQuery } from '@tanstack/react-query'
import DashboardLayout from '@/components/layout/DashboardLayout'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Badge } from '@/components/ui/badge'
import { Switch } from '@/components/ui/switch'
import { useToast } from '@/components/ui/use-toast'
import {
  ArrowLeft,
  Save,
  Play,
  Eye,
  FileText,
  Database,
  Filter,
  Layout,
  BarChart3,
  LineChart,
  PieChart,
  TableIcon,
  Plus,
  Trash2,
  Settings,
  Code,
  Calendar,
  Clock,
  Users,
  ChevronDown,
  ChevronUp,
  Move,
  Copy,
  AlertCircle,
  CheckCircle
} from 'lucide-react'
import reportService from '@/services/report.service'
import {
  DndContext,
  closestCenter,
  KeyboardSensor,
  PointerSensor,
  useSensor,
  useSensors,
  DragEndEvent,
} from '@dnd-kit/core'
import {
  arrayMove,
  SortableContext,
  sortableKeyboardCoordinates,
  verticalListSortingStrategy,
} from '@dnd-kit/sortable'
import { useSortable } from '@dnd-kit/sortable'
import { CSS } from '@dnd-kit/utilities'

// Component types
const COMPONENT_TYPES = [
  { id: 'text', name: '文字', icon: FileText, category: 'basic' },
  { id: 'table', name: '表格', icon: TableIcon, category: 'basic' },
  { id: 'chart_bar', name: '長條圖', icon: BarChart3, category: 'chart' },
  { id: 'chart_line', name: '折線圖', icon: LineChart, category: 'chart' },
  { id: 'chart_pie', name: '圓餅圖', icon: PieChart, category: 'chart' },
  { id: 'kpi', name: 'KPI 卡片', icon: Layout, category: 'advanced' },
  { id: 'filter', name: '篩選器', icon: Filter, category: 'advanced' },
]

// Sortable component
function SortableComponent({ component, onEdit, onRemove, onDuplicate }: any) {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id: component.id })

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
  }

  const ComponentIcon = COMPONENT_TYPES.find(t => t.id === component.type)?.icon || FileText

  return (
    <div
      ref={setNodeRef}
      style={style}
      className={`border rounded-lg p-4 bg-white ${isDragging ? 'shadow-lg' : 'shadow'}`}
    >
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <div {...attributes} {...listeners} className="cursor-move">
            <Move className="h-4 w-4 text-gray-400" />
          </div>
          <ComponentIcon className="h-4 w-4 text-gray-500" />
          <span className="font-medium">{component.name || '未命名元件'}</span>
        </div>
        <div className="flex items-center gap-1">
          <Button
            variant="ghost"
            size="sm"
            onClick={() => onEdit(component)}
          >
            <Settings className="h-4 w-4" />
          </Button>
          <Button
            variant="ghost"
            size="sm"
            onClick={() => onDuplicate(component)}
          >
            <Copy className="h-4 w-4" />
          </Button>
          <Button
            variant="ghost"
            size="sm"
            onClick={() => onRemove(component.id)}
          >
            <Trash2 className="h-4 w-4" />
          </Button>
        </div>
      </div>
      {component.config?.description && (
        <p className="text-sm text-gray-500 mt-2">{component.config.description}</p>
      )}
    </div>
  )
}

export default function ReportDesignerPage() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const { toast } = useToast()
  const templateId = searchParams.get('template')

  const [activeTab, setActiveTab] = useState('design')
  const [components, setComponents] = useState<any[]>([])
  const [selectedComponent, setSelectedComponent] = useState<any>(null)
  const [previewMode, setPreviewMode] = useState(false)

  const [formData, setFormData] = useState({
    name: '',
    report_no: '',
    description: '',
    category: 'sales',
    type: 'summary',
    status: 'active',
    schedule_config: {
      enabled: false,
      frequency: 'daily',
      time: '08:00',
      recipients: [],
    },
    permissions: {
      view_users: [],
      edit_users: [],
      is_public: false,
    },
  })

  const sensors = useSensors(
    useSensor(PointerSensor),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates,
    })
  )

  // Load template if specified
  useEffect(() => {
    if (templateId) {
      loadTemplate(templateId)
    }
  }, [templateId])

  const loadTemplate = async (id: string) => {
    try {
      const template = await reportService.getReportTemplate(id)
      if ((template as any).config) {
        setComponents((template as any).config.components || [])
        setFormData({
          ...formData,
          name: template.name,
          description: template.description,
          category: template.category,
          type: template.type,
        })
      }
    } catch (error) {
      toast({
        title: '載入範本失敗',
        variant: 'destructive',
      })
    }
  }

  // Create report mutation
  const createMutation = useMutation({
    mutationFn: (data: any) => reportService.createReport(data),
    onSuccess: (data) => {
      toast({ title: '報表建立成功' })
      router.push(`/reports/${data.id}`)
    },
    onError: (error: any) => {
      toast({
        title: '建立失敗',
        description: error.response?.data?.message || '建立報表時發生錯誤',
        variant: 'destructive',
      })
    },
  })

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event

    if (active.id !== over?.id) {
      setComponents((items) => {
        const oldIndex = items.findIndex((item) => item.id === active.id)
        const newIndex = items.findIndex((item) => item.id === over?.id)
        return arrayMove(items, oldIndex, newIndex)
      })
    }
  }

  const addComponent = (type: string) => {
    const newComponent = {
      id: `component_${Date.now()}`,
      type,
      name: COMPONENT_TYPES.find(t => t.id === type)?.name || '新元件',
      config: {},
    }
    setComponents([...components, newComponent])
  }

  const removeComponent = (id: string) => {
    setComponents(components.filter(c => c.id !== id))
  }

  const duplicateComponent = (component: any) => {
    const newComponent = {
      ...component,
      id: `component_${Date.now()}`,
      name: `${component.name} (複製)`,
    }
    const index = components.findIndex(c => c.id === component.id)
    const newComponents = [...components]
    newComponents.splice(index + 1, 0, newComponent)
    setComponents(newComponents)
  }

  const editComponent = (component: any) => {
    setSelectedComponent(component)
    // Open component editor
  }

  const updateComponent = (id: string, updates: any) => {
    setComponents(components.map(c => 
      c.id === id ? { ...c, ...updates } : c
    ))
  }

  const handleSave = () => {
    const reportData = {
      ...formData,
      config: {
        components,
        layout: 'vertical',
        settings: {},
      },
    }
    createMutation.mutate(reportData)
  }

  const handlePreview = () => {
    setPreviewMode(!previewMode)
  }

  return (
    <DashboardLayout>
      <div className="h-full flex flex-col">
        {/* Header */}
        <div className="flex items-center justify-between px-6 py-4 border-b">
          <div className="flex items-center gap-4">
            <Button variant="ghost" size="icon" onClick={() => router.back()}>
              <ArrowLeft className="h-4 w-4" />
            </Button>
            <div>
              <h1 className="text-2xl font-bold">新增報表</h1>
              <p className="text-sm text-gray-600">設計您的自訂報表</p>
            </div>
          </div>
          <div className="flex items-center gap-2">
            <Button variant="outline" onClick={handlePreview}>
              <Eye className="mr-2 h-4 w-4" />
              {previewMode ? '編輯' : '預覽'}
            </Button>
            <Button onClick={handleSave}>
              <Save className="mr-2 h-4 w-4" />
              儲存報表
            </Button>
          </div>
        </div>

        {/* Main Content */}
        <div className="flex-1 flex overflow-hidden">
          {/* Sidebar */}
          {!previewMode && (
            <div className="w-80 border-r bg-gray-50 overflow-y-auto">
              <Tabs value={activeTab} onValueChange={setActiveTab}>
                <TabsList className="w-full">
                  <TabsTrigger value="design" className="flex-1">設計</TabsTrigger>
                  <TabsTrigger value="data" className="flex-1">資料</TabsTrigger>
                  <TabsTrigger value="settings" className="flex-1">設定</TabsTrigger>
                </TabsList>

                <TabsContent value="design" className="p-4 space-y-4">
                  <div>
                    <h3 className="font-medium mb-2">基本元件</h3>
                    <div className="grid grid-cols-2 gap-2">
                      {COMPONENT_TYPES.filter(c => c.category === 'basic').map((component) => {
                        const Icon = component.icon
                        return (
                          <Button
                            key={component.id}
                            variant="outline"
                            className="h-20 flex flex-col gap-2"
                            onClick={() => addComponent(component.id)}
                          >
                            <Icon className="h-6 w-6" />
                            <span className="text-xs">{component.name}</span>
                          </Button>
                        )
                      })}
                    </div>
                  </div>

                  <div>
                    <h3 className="font-medium mb-2">圖表元件</h3>
                    <div className="grid grid-cols-2 gap-2">
                      {COMPONENT_TYPES.filter(c => c.category === 'chart').map((component) => {
                        const Icon = component.icon
                        return (
                          <Button
                            key={component.id}
                            variant="outline"
                            className="h-20 flex flex-col gap-2"
                            onClick={() => addComponent(component.id)}
                          >
                            <Icon className="h-6 w-6" />
                            <span className="text-xs">{component.name}</span>
                          </Button>
                        )
                      })}
                    </div>
                  </div>

                  <div>
                    <h3 className="font-medium mb-2">進階元件</h3>
                    <div className="grid grid-cols-2 gap-2">
                      {COMPONENT_TYPES.filter(c => c.category === 'advanced').map((component) => {
                        const Icon = component.icon
                        return (
                          <Button
                            key={component.id}
                            variant="outline"
                            className="h-20 flex flex-col gap-2"
                            onClick={() => addComponent(component.id)}
                          >
                            <Icon className="h-6 w-6" />
                            <span className="text-xs">{component.name}</span>
                          </Button>
                        )
                      })}
                    </div>
                  </div>
                </TabsContent>

                <TabsContent value="data" className="p-4 space-y-4">
                  <Card>
                    <CardHeader>
                      <CardTitle className="text-base">資料來源</CardTitle>
                    </CardHeader>
                    <CardContent>
                      <Select>
                        <SelectTrigger>
                          <SelectValue placeholder="選擇資料來源" />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="sales">銷售資料</SelectItem>
                          <SelectItem value="inventory">庫存資料</SelectItem>
                          <SelectItem value="finance">財務資料</SelectItem>
                          <SelectItem value="custom">自訂查詢</SelectItem>
                        </SelectContent>
                      </Select>
                    </CardContent>
                  </Card>

                  <Card>
                    <CardHeader>
                      <CardTitle className="text-base">篩選條件</CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-2">
                      <Button variant="outline" className="w-full">
                        <Plus className="mr-2 h-4 w-4" />
                        新增篩選條件
                      </Button>
                    </CardContent>
                  </Card>

                  <Card>
                    <CardHeader>
                      <CardTitle className="text-base">參數設定</CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-2">
                      <Button variant="outline" className="w-full">
                        <Plus className="mr-2 h-4 w-4" />
                        新增參數
                      </Button>
                    </CardContent>
                  </Card>
                </TabsContent>

                <TabsContent value="settings" className="p-4 space-y-4">
                  <div className="space-y-4">
                    <div className="grid gap-2">
                      <Label htmlFor="name">報表名稱 *</Label>
                      <Input
                        id="name"
                        value={formData.name}
                        onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                        placeholder="例如：月度銷售報表"
                      />
                    </div>

                    <div className="grid gap-2">
                      <Label htmlFor="report_no">報表編號</Label>
                      <Input
                        id="report_no"
                        value={formData.report_no}
                        onChange={(e) => setFormData({ ...formData, report_no: e.target.value })}
                        placeholder="系統將自動產生"
                      />
                    </div>

                    <div className="grid gap-2">
                      <Label htmlFor="description">描述</Label>
                      <Textarea
                        id="description"
                        value={formData.description}
                        onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                        placeholder="描述此報表的用途"
                        rows={3}
                      />
                    </div>

                    <div className="grid grid-cols-2 gap-4">
                      <div className="grid gap-2">
                        <Label htmlFor="category">分類</Label>
                        <Select
                          value={formData.category}
                          onValueChange={(value) => setFormData({ ...formData, category: value })}
                        >
                          <SelectTrigger>
                            <SelectValue />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectItem value="sales">銷售</SelectItem>
                            <SelectItem value="finance">財務</SelectItem>
                            <SelectItem value="production">生產</SelectItem>
                            <SelectItem value="inventory">庫存</SelectItem>
                            <SelectItem value="supplier">供應商</SelectItem>
                            <SelectItem value="customer">客戶</SelectItem>
                            <SelectItem value="system">系統</SelectItem>
                          </SelectContent>
                        </Select>
                      </div>
                      <div className="grid gap-2">
                        <Label htmlFor="type">類型</Label>
                        <Select
                          value={formData.type}
                          onValueChange={(value) => setFormData({ ...formData, type: value })}
                        >
                          <SelectTrigger>
                            <SelectValue />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectItem value="summary">摘要</SelectItem>
                            <SelectItem value="detail">詳細</SelectItem>
                            <SelectItem value="trend">趨勢</SelectItem>
                            <SelectItem value="comparison">比較</SelectItem>
                            <SelectItem value="dashboard">儀表板</SelectItem>
                          </SelectContent>
                        </Select>
                      </div>
                    </div>

                    <Card>
                      <CardHeader>
                        <CardTitle className="text-base">排程設定</CardTitle>
                      </CardHeader>
                      <CardContent className="space-y-4">
                        <div className="flex items-center space-x-2">
                          <Switch
                            checked={formData.schedule_config.enabled}
                            onCheckedChange={(checked) => 
                              setFormData({
                                ...formData,
                                schedule_config: { ...formData.schedule_config, enabled: checked }
                              })
                            }
                          />
                          <Label>啟用排程</Label>
                        </div>
                        
                        {formData.schedule_config.enabled && (
                          <>
                            <div className="grid gap-2">
                              <Label>頻率</Label>
                              <Select
                                value={formData.schedule_config.frequency}
                                onValueChange={(value) => 
                                  setFormData({
                                    ...formData,
                                    schedule_config: { ...formData.schedule_config, frequency: value }
                                  })
                                }
                              >
                                <SelectTrigger>
                                  <SelectValue />
                                </SelectTrigger>
                                <SelectContent>
                                  <SelectItem value="daily">每日</SelectItem>
                                  <SelectItem value="weekly">每週</SelectItem>
                                  <SelectItem value="monthly">每月</SelectItem>
                                </SelectContent>
                              </Select>
                            </div>
                            
                            <div className="grid gap-2">
                              <Label>執行時間</Label>
                              <Input
                                type="time"
                                value={formData.schedule_config.time}
                                onChange={(e) => 
                                  setFormData({
                                    ...formData,
                                    schedule_config: { ...formData.schedule_config, time: e.target.value }
                                  })
                                }
                              />
                            </div>
                          </>
                        )}
                      </CardContent>
                    </Card>

                    <Card>
                      <CardHeader>
                        <CardTitle className="text-base">權限設定</CardTitle>
                      </CardHeader>
                      <CardContent className="space-y-4">
                        <div className="flex items-center space-x-2">
                          <Switch
                            checked={formData.permissions.is_public}
                            onCheckedChange={(checked) => 
                              setFormData({
                                ...formData,
                                permissions: { ...formData.permissions, is_public: checked }
                              })
                            }
                          />
                          <Label>公開報表</Label>
                        </div>
                      </CardContent>
                    </Card>
                  </div>
                </TabsContent>
              </Tabs>
            </div>
          )}

          {/* Canvas */}
          <div className="flex-1 overflow-auto p-6 bg-gray-100">
            <div className="max-w-5xl mx-auto">
              {components.length === 0 ? (
                <Card className="border-dashed">
                  <CardContent className="flex flex-col items-center justify-center py-12">
                    <Layout className="h-12 w-12 text-gray-300 mb-4" />
                    <p className="text-gray-500 mb-4">
                      {previewMode ? '此報表尚無內容' : '從左側拖曳元件開始設計報表'}
                    </p>
                    {!previewMode && (
                      <Button variant="outline" onClick={() => setActiveTab('design')}>
                        開始設計
                      </Button>
                    )}
                  </CardContent>
                </Card>
              ) : previewMode ? (
                <div className="space-y-4">
                  {components.map((component) => (
                    <Card key={component.id}>
                      <CardHeader>
                        <CardTitle>{component.name}</CardTitle>
                      </CardHeader>
                      <CardContent>
                        <div className="h-48 bg-gray-100 rounded flex items-center justify-center">
                          <p className="text-gray-500">預覽內容</p>
                        </div>
                      </CardContent>
                    </Card>
                  ))}
                </div>
              ) : (
                <DndContext
                  sensors={sensors}
                  collisionDetection={closestCenter}
                  onDragEnd={handleDragEnd}
                >
                  <SortableContext
                    items={components}
                    strategy={verticalListSortingStrategy}
                  >
                    <div className="space-y-4">
                      {components.map((component) => (
                        <SortableComponent
                          key={component.id}
                          component={component}
                          onEdit={editComponent}
                          onRemove={removeComponent}
                          onDuplicate={duplicateComponent}
                        />
                      ))}
                    </div>
                  </SortableContext>
                </DndContext>
              )}
            </div>
          </div>
        </div>
      </div>
    </DashboardLayout>
  )
}