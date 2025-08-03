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
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Textarea } from '@/components/ui/textarea'
import { 
  Workflow, 
  Plus, 
  Play, 
  Pause, 
  Clock, 
  CheckCircle, 
  XCircle,
  AlertCircle,
  RefreshCw,
  Webhook,
  Calendar,
  Settings,
  Activity,
  Zap,
  Link,
  ExternalLink
} from 'lucide-react'
import { useToast } from '@/components/ui/use-toast'
import n8nService, { N8NWorkflow, WORKFLOW_TEMPLATES } from '@/services/n8n.service'
import { format } from 'date-fns'
import { zhTW } from 'date-fns/locale'

export default function WorkflowsPage() {
  const queryClient = useQueryClient()
  const { toast } = useToast()
  const [activeTab, setActiveTab] = useState('workflows')
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false)
  const [selectedTemplate, setSelectedTemplate] = useState('')
  
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    workflow_id: '',
    trigger_type: 'webhook',
    trigger_config: {},
  })

  // Test N8N connection
  const { data: connectionStatus } = useQuery({
    queryKey: ['n8n-connection'],
    queryFn: () => n8nService.testConnection(),
    retry: 1,
  })

  // Fetch workflows
  const { data: workflows, isLoading: workflowsLoading } = useQuery({
    queryKey: ['workflows'],
    queryFn: () => n8nService.listWorkflows(),
    enabled: connectionStatus?.connected,
  })

  // Fetch executions
  const { data: executions, isLoading: executionsLoading } = useQuery({
    queryKey: ['executions'],
    queryFn: () => n8nService.getExecutions(),
    enabled: connectionStatus?.connected,
  })

  // Fetch available workflows from N8N
  const { data: availableWorkflows = [] } = useQuery({
    queryKey: ['available-workflows'],
    queryFn: () => n8nService.getAvailableWorkflows(),
    enabled: connectionStatus?.connected && isCreateDialogOpen,
  })

  // Create workflow mutation
  const createWorkflowMutation = useMutation({
    mutationFn: (data: any) => n8nService.createWorkflow(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['workflows'] })
      toast({ title: '工作流程建立成功' })
      setIsCreateDialogOpen(false)
      resetForm()
    },
    onError: (error: any) => {
      toast({
        title: '建立失敗',
        description: error.response?.data?.message || '建立工作流程時發生錯誤',
        variant: 'destructive',
      })
    },
  })

  // Toggle workflow mutation
  const toggleWorkflowMutation = useMutation({
    mutationFn: ({ id, active }: { id: string; active: boolean }) => 
      n8nService.toggleWorkflow(id, active),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['workflows'] })
      toast({ title: '工作流程狀態已更新' })
    },
    onError: (error: any) => {
      toast({
        title: '更新失敗',
        description: error.response?.data?.message || '更新狀態時發生錯誤',
        variant: 'destructive',
      })
    },
  })

  // Trigger workflow mutation
  const triggerWorkflowMutation = useMutation({
    mutationFn: (workflowId: string) => 
      n8nService.triggerWorkflow({ 
        workflow_id: workflowId, 
        data: { manual_trigger: true } 
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['executions'] })
      toast({ title: '工作流程已觸發' })
    },
    onError: (error: any) => {
      toast({
        title: '觸發失敗',
        description: error.response?.data?.message || '觸發工作流程時發生錯誤',
        variant: 'destructive',
      })
    },
  })

  const resetForm = () => {
    setFormData({
      name: '',
      description: '',
      workflow_id: '',
      trigger_type: 'webhook',
      trigger_config: {},
    })
    setSelectedTemplate('')
  }

  const handleTemplateSelect = (templateKey: string) => {
    const template = WORKFLOW_TEMPLATES[templateKey as keyof typeof WORKFLOW_TEMPLATES]
    if (template) {
      setFormData({
        name: template.name,
        description: template.description || '',
        workflow_id: '',
        trigger_type: template.trigger_type,
        trigger_config: template.trigger_type === 'event' 
          ? { event: (template as any).event }
          : template.trigger_type === 'schedule'
          ? { cron: (template as any).schedule }
          : {},
      })
    }
  }

  const handleSubmit = () => {
    if (!formData.name || !formData.workflow_id) {
      toast({
        title: '請填寫必要欄位',
        variant: 'destructive',
      })
      return
    }
    createWorkflowMutation.mutate(formData)
  }

  const getExecutionStatusBadge = (status: string) => {
    const statusConfig: Record<string, { label: string; variant: any; icon: any }> = {
      running: { label: '執行中', variant: 'warning', icon: RefreshCw },
      success: { label: '成功', variant: 'success', icon: CheckCircle },
      error: { label: '錯誤', variant: 'destructive', icon: XCircle },
      canceled: { label: '已取消', variant: 'secondary', icon: XCircle },
    }

    const config = statusConfig[status] || { label: status, variant: 'default', icon: Activity }
    const Icon = config.icon
    
    return (
      <Badge variant={config.variant as any} className="flex items-center gap-1">
        <Icon className="h-3 w-3" />
        {config.label}
      </Badge>
    )
  }

  const getTriggerIcon = (triggerType: string) => {
    switch (triggerType) {
      case 'webhook':
        return <Webhook className="h-4 w-4" />
      case 'schedule':
        return <Clock className="h-4 w-4" />
      case 'event':
        return <Zap className="h-4 w-4" />
      default:
        return <Activity className="h-4 w-4" />
    }
  }

  if (!connectionStatus?.connected) {
    return (
      <DashboardLayout>
        <div className="space-y-6">
          <div>
            <h1 className="text-3xl font-bold text-gray-900 flex items-center gap-2">
              <Workflow className="h-8 w-8" />
              工作流程自動化
            </h1>
            <p className="mt-2 text-gray-600">使用 N8N 建立自動化工作流程</p>
          </div>

          <Alert>
            <AlertCircle className="h-4 w-4" />
            <AlertTitle>N8N 連線失敗</AlertTitle>
            <AlertDescription>
              {connectionStatus?.message || '無法連接到 N8N 服務，請檢查設定或聯絡系統管理員。'}
            </AlertDescription>
          </Alert>
        </div>
      </DashboardLayout>
    )
  }

  return (
    <DashboardLayout>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex justify-between items-center">
          <div>
            <h1 className="text-3xl font-bold text-gray-900 flex items-center gap-2">
              <Workflow className="h-8 w-8" />
              工作流程自動化
            </h1>
            <p className="mt-2 text-gray-600">
              使用 N8N 建立自動化工作流程 
              {connectionStatus?.version && (
                <span className="text-sm text-gray-500 ml-2">
                  (N8N v{connectionStatus.version})
                </span>
              )}
            </p>
          </div>
          <div className="flex gap-2">
            <Button
              variant="outline"
              onClick={() => window.open(process.env.NEXT_PUBLIC_N8N_URL || 'http://localhost:5678', '_blank')}
            >
              <ExternalLink className="mr-2 h-4 w-4" />
              開啟 N8N
            </Button>
            <Button onClick={() => { resetForm(); setIsCreateDialogOpen(true); }}>
              <Plus className="mr-2 h-4 w-4" />
              新增工作流程
            </Button>
          </div>
        </div>

        <Tabs value={activeTab} onValueChange={setActiveTab}>
          <TabsList>
            <TabsTrigger value="workflows">工作流程</TabsTrigger>
            <TabsTrigger value="executions">執行記錄</TabsTrigger>
            <TabsTrigger value="templates">範本</TabsTrigger>
          </TabsList>

          <TabsContent value="workflows" className="space-y-4">
            {workflowsLoading ? (
              <div className="text-center py-8">載入中...</div>
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {workflows?.data.map((workflow) => (
                  <Card key={workflow.id}>
                    <CardHeader>
                      <div className="flex items-start justify-between">
                        <div className="space-y-1">
                          <CardTitle className="text-lg flex items-center gap-2">
                            {getTriggerIcon(workflow.trigger_type)}
                            {workflow.name}
                          </CardTitle>
                          <CardDescription>
                            {workflow.description || '無描述'}
                          </CardDescription>
                        </div>
                        <Switch
                          checked={workflow.is_active}
                          onCheckedChange={(checked) => 
                            toggleWorkflowMutation.mutate({ 
                              id: workflow.id, 
                              active: checked 
                            })
                          }
                        />
                      </div>
                    </CardHeader>
                    <CardContent>
                      <div className="space-y-2 text-sm">
                        <div className="flex justify-between">
                          <span className="text-gray-500">觸發方式</span>
                          <Badge variant="outline">{workflow.trigger_type}</Badge>
                        </div>
                        <div className="flex justify-between">
                          <span className="text-gray-500">工作流程 ID</span>
                          <code className="text-xs bg-gray-100 px-1 rounded">
                            {workflow.workflow_id}
                          </code>
                        </div>
                        <div className="flex justify-between">
                          <span className="text-gray-500">建立時間</span>
                          <span>
                            {format(new Date(workflow.created_at), 'yyyy/MM/dd', {
                              locale: zhTW,
                            })}
                          </span>
                        </div>
                      </div>
                      <div className="mt-4 flex gap-2">
                        <Button
                          size="sm"
                          variant="outline"
                          className="flex-1"
                          onClick={() => triggerWorkflowMutation.mutate(workflow.workflow_id)}
                          disabled={!workflow.is_active}
                        >
                          <Play className="mr-1 h-3 w-3" />
                          手動觸發
                        </Button>
                        <Button
                          size="sm"
                          variant="outline"
                          className="flex-1"
                        >
                          <Settings className="mr-1 h-3 w-3" />
                          設定
                        </Button>
                      </div>
                    </CardContent>
                  </Card>
                ))}
              </div>
            )}
          </TabsContent>

          <TabsContent value="executions" className="space-y-4">
            <Card>
              <CardHeader>
                <CardTitle>執行記錄</CardTitle>
                <CardDescription>查看工作流程執行歷史</CardDescription>
              </CardHeader>
              <CardContent>
                {executionsLoading ? (
                  <div className="text-center py-8">載入中...</div>
                ) : (
                  <div className="overflow-x-auto">
                    <table className="w-full">
                      <thead>
                        <tr className="border-b">
                          <th className="text-left p-2">工作流程</th>
                          <th className="text-left p-2">執行 ID</th>
                          <th className="text-left p-2">狀態</th>
                          <th className="text-left p-2">開始時間</th>
                          <th className="text-left p-2">結束時間</th>
                          <th className="text-left p-2">耗時</th>
                        </tr>
                      </thead>
                      <tbody>
                        {executions?.data.map((execution) => (
                          <tr key={execution.id} className="border-b hover:bg-gray-50">
                            <td className="p-2">{execution.workflow_id}</td>
                            <td className="p-2">
                              <code className="text-xs bg-gray-100 px-1 rounded">
                                {execution.execution_id}
                              </code>
                            </td>
                            <td className="p-2">
                              {getExecutionStatusBadge(execution.status)}
                            </td>
                            <td className="p-2 text-sm">
                              {format(new Date(execution.started_at), 'MM/dd HH:mm:ss')}
                            </td>
                            <td className="p-2 text-sm">
                              {execution.finished_at 
                                ? format(new Date(execution.finished_at), 'MM/dd HH:mm:ss')
                                : '-'
                              }
                            </td>
                            <td className="p-2 text-sm">
                              {execution.finished_at 
                                ? `${Math.round((new Date(execution.finished_at).getTime() - new Date(execution.started_at).getTime()) / 1000)}s`
                                : '-'
                              }
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                )}
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="templates" className="space-y-4">
            <Card>
              <CardHeader>
                <CardTitle>工作流程範本</CardTitle>
                <CardDescription>快速建立常用的自動化流程</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  {Object.entries(WORKFLOW_TEMPLATES).map(([key, template]) => (
                    <Card key={key} className="cursor-pointer hover:shadow-md transition-shadow">
                      <CardContent className="p-4">
                        <div className="flex items-start justify-between">
                          <div className="space-y-1 flex-1">
                            <h3 className="font-medium flex items-center gap-2">
                              {getTriggerIcon(template.trigger_type)}
                              {template.name}
                            </h3>
                            <p className="text-sm text-gray-600">
                              {template.description}
                            </p>
                            <div className="flex items-center gap-2 mt-2">
                              <Badge variant="outline" className="text-xs">
                                {template.trigger_type}
                              </Badge>
                              {'event' in template && (
                                <Badge variant="secondary" className="text-xs">
                                  {template.event}
                                </Badge>
                              )}
                              {'schedule' in template && (
                                <Badge variant="secondary" className="text-xs">
                                  {template.schedule}
                                </Badge>
                              )}
                            </div>
                          </div>
                          <Button
                            size="sm"
                            onClick={() => {
                              setSelectedTemplate(key)
                              handleTemplateSelect(key)
                              setIsCreateDialogOpen(true)
                            }}
                          >
                            使用
                          </Button>
                        </div>
                      </CardContent>
                    </Card>
                  ))}
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>

        {/* Create Workflow Dialog */}
        <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
          <DialogContent className="sm:max-w-[600px]">
            <DialogHeader>
              <DialogTitle>新增工作流程</DialogTitle>
              <DialogDescription>
                設定工作流程的基本資訊和觸發條件
              </DialogDescription>
            </DialogHeader>
            <div className="grid gap-4 py-4">
              <div className="grid gap-2">
                <Label htmlFor="name">名稱 *</Label>
                <Input
                  id="name"
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  placeholder="例如：每日匯率更新"
                />
              </div>

              <div className="grid gap-2">
                <Label htmlFor="description">描述</Label>
                <Textarea
                  id="description"
                  value={formData.description}
                  onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                  placeholder="描述此工作流程的用途..."
                  rows={3}
                />
              </div>

              <div className="grid gap-2">
                <Label htmlFor="workflow_id">N8N 工作流程 *</Label>
                <Select
                  value={formData.workflow_id}
                  onValueChange={(value) => setFormData({ ...formData, workflow_id: value })}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="選擇 N8N 工作流程" />
                  </SelectTrigger>
                  <SelectContent>
                    {availableWorkflows.map((workflow) => (
                      <SelectItem key={workflow.id} value={workflow.id}>
                        {workflow.name}
                        {workflow.tags.length > 0 && (
                          <span className="text-xs text-gray-500 ml-2">
                            ({workflow.tags.join(', ')})
                          </span>
                        )}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              <div className="grid gap-2">
                <Label htmlFor="trigger_type">觸發方式</Label>
                <Select
                  value={formData.trigger_type}
                  onValueChange={(value) => setFormData({ ...formData, trigger_type: value })}
                >
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="webhook">Webhook</SelectItem>
                    <SelectItem value="schedule">排程</SelectItem>
                    <SelectItem value="event">事件</SelectItem>
                    <SelectItem value="manual">手動</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              {formData.trigger_type === 'schedule' && (
                <div className="grid gap-2">
                  <Label htmlFor="cron">Cron 表達式</Label>
                  <Input
                    id="cron"
                    value={(formData.trigger_config as any).cron || ''}
                    onChange={(e) => setFormData({ 
                      ...formData, 
                      trigger_config: { ...formData.trigger_config, cron: e.target.value }
                    })}
                    placeholder="例如：0 9 * * * (每天早上9點)"
                  />
                </div>
              )}

              {formData.trigger_type === 'event' && (
                <div className="grid gap-2">
                  <Label htmlFor="event">事件類型</Label>
                  <Select
                    value={(formData.trigger_config as any).event || ''}
                    onValueChange={(value) => setFormData({ 
                      ...formData, 
                      trigger_config: { ...formData.trigger_config, event: value }
                    })}
                  >
                    <SelectTrigger>
                      <SelectValue placeholder="選擇事件類型" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="inquiry.created">詢價單建立</SelectItem>
                      <SelectItem value="quote.created">報價單建立</SelectItem>
                      <SelectItem value="quote.submitted_for_review">報價單送審</SelectItem>
                      <SelectItem value="quote.approved">報價單核准</SelectItem>
                      <SelectItem value="quote.sent">報價單發送</SelectItem>
                      <SelectItem value="customer.created">客戶建立</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              )}
            </div>
            <DialogFooter>
              <Button variant="outline" onClick={() => setIsCreateDialogOpen(false)}>
                取消
              </Button>
              <Button onClick={handleSubmit}>
                建立
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>
    </DashboardLayout>
  )
}