'use client'

import { useState, useEffect } from 'react'
import { useQuery } from '@tanstack/react-query'
import DashboardLayout from '@/components/layout/DashboardLayout'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Progress } from '@/components/ui/progress'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { 
  Link2,
  Webhook,
  Database,
  RefreshCw as Sync,
  Key,
  Settings,
  Play,
  Pause,
  AlertTriangle,
  CheckCircle,
  XCircle,
  Clock,
  Zap,
  Activity,
  BarChart3,
  Plus,
  Search,
  Filter,
  Download,
  Upload,
  RefreshCw,
  Eye,
  Edit,
  Trash2,
  TestTube,
  Globe,
  Server,
  Code,
  Package,
  Truck,
  FileText,
  TrendingUp,
  TrendingDown,
  Users,
  Calendar,
  History,
  Shield,
  Network,
  Cpu
} from 'lucide-react'
import integrationService from '@/services/integration.service'
import { format } from 'date-fns'
import { zhTW } from 'date-fns/locale'

export default function IntegrationsPage() {
  const [activeTab, setActiveTab] = useState('overview')
  const [searchQuery, setSearchQuery] = useState('')
  const [statusFilter, setStatusFilter] = useState('')
  const [typeFilter, setTypeFilter] = useState('')

  // Fetch integrations data
  const { data: integrations, isLoading: isLoadingIntegrations } = useQuery({
    queryKey: ['integrations', searchQuery, statusFilter, typeFilter],
    queryFn: () => integrationService.listIntegrations({
      status: statusFilter || undefined,
      type: typeFilter || undefined,
    }),
  })

  // Fetch webhooks data
  const { data: webhooks } = useQuery({
    queryKey: ['webhooks'],
    queryFn: () => integrationService.listWebhooks({ limit: 10 }),
  })

  // Fetch sync jobs data
  const { data: syncJobs } = useQuery({
    queryKey: ['sync-jobs'],
    queryFn: () => integrationService.listDataSyncJobs('all', { limit: 10 }),
  })

  // Fetch API keys data
  const { data: apiKeys } = useQuery({
    queryKey: ['api-keys'],
    queryFn: () => integrationService.listApiKeys({ limit: 10 }),
  })

  // Fetch external systems data
  const { data: externalSystems } = useQuery({
    queryKey: ['external-systems'],
    queryFn: () => integrationService.listExternalSystems({ limit: 10 }),
  })

  // Fetch analytics data
  const { data: integrationStats } = useQuery({
    queryKey: ['integration-stats'],
    queryFn: () => integrationService.getIntegrationStats(),
  })

  const { data: integrationsByType } = useQuery({
    queryKey: ['integrations-by-type'],
    queryFn: () => integrationService.getIntegrationsByType(),
  })

  const { data: syncTrends } = useQuery({
    queryKey: ['sync-trends'],
    queryFn: () => integrationService.getSyncJobTrends({ days: 30 }),
  })

  const getStatusBadge = (status: string, type: 'integration' | 'webhook' | 'sync' | 'system') => {
    const statusConfig: Record<string, Record<string, { label: string; variant: any; icon: any }>> = {
      integration: {
        active: { label: '啟用', variant: 'success', icon: CheckCircle },
        inactive: { label: '停用', variant: 'secondary', icon: Pause },
        error: { label: '錯誤', variant: 'destructive', icon: XCircle },
        testing: { label: '測試中', variant: 'warning', icon: TestTube },
      },
      webhook: {
        active: { label: '啟用', variant: 'success', icon: CheckCircle },
        inactive: { label: '停用', variant: 'secondary', icon: Pause },
      },
      sync: {
        pending: { label: '等待中', variant: 'secondary', icon: Clock },
        running: { label: '執行中', variant: 'info', icon: Play },
        completed: { label: '已完成', variant: 'success', icon: CheckCircle },
        failed: { label: '失敗', variant: 'destructive', icon: XCircle },
        cancelled: { label: '已取消', variant: 'secondary', icon: XCircle },
      },
      system: {
        active: { label: '連線中', variant: 'success', icon: CheckCircle },
        inactive: { label: '未連線', variant: 'secondary', icon: Pause },
        testing: { label: '測試中', variant: 'warning', icon: TestTube },
        error: { label: '連線錯誤', variant: 'destructive', icon: XCircle },
      },
    }

    const config = statusConfig[type]?.[status] || { label: status, variant: 'default', icon: AlertTriangle }
    const Icon = config.icon
    
    return (
      <Badge variant={config.variant as any} className="flex items-center gap-1">
        <Icon className="h-3 w-3" />
        {config.label}
      </Badge>
    )
  }

  const getTypeIcon = (type: string) => {
    const typeIcons: Record<string, any> = {
      api: Globe,
      webhook: Webhook,
      ftp: Server,
      sftp: Shield,
      email: FileText,
      database: Database,
      erp: Settings,
      crm: Users,
      accounting: BarChart3,
      warehouse: Package,
      shipping: Truck,
    }
    
    const Icon = typeIcons[type.toLowerCase()] || Network
    return <Icon className="h-4 w-4" />
  }

  const getPriorityBadge = (priority: string) => {
    const priorityConfig: Record<string, { label: string; variant: any; className: string }> = {
      low: { label: '低', variant: 'secondary', className: 'text-gray-600' },
      normal: { label: '一般', variant: 'info', className: 'text-blue-600' },
      high: { label: '高', variant: 'warning', className: 'text-orange-600' },
      urgent: { label: '緊急', variant: 'destructive', className: 'text-red-600 animate-pulse' },
    }

    const config = priorityConfig[priority] || { label: priority, variant: 'default', className: '' }
    
    return (
      <Badge variant={config.variant as any} className={config.className}>
        {config.label}
      </Badge>
    )
  }

  if (isLoadingIntegrations) {
    return (
      <DashboardLayout>
        <div className="flex items-center justify-center h-64">
          <div className="text-center">載入中...</div>
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
            <h1 className="text-3xl font-bold text-gray-900">整合功能</h1>
            <p className="mt-2 text-gray-600">API 整合、Webhook、數據同步與第三方系統連接</p>
          </div>
          <div className="flex items-center gap-4">
            <Button variant="outline">
              <Download className="mr-2 h-4 w-4" />
              匯出配置
            </Button>
            <Button variant="outline">
              <Upload className="mr-2 h-4 w-4" />
              匯入配置
            </Button>
            <Button>
              <Plus className="mr-2 h-4 w-4" />
              新增整合
            </Button>
          </div>
        </div>

        <Tabs value={activeTab} onValueChange={setActiveTab}>
          <TabsList className="grid w-full grid-cols-8">
            <TabsTrigger value="overview">概覽</TabsTrigger>
            <TabsTrigger value="integrations">整合管理</TabsTrigger>
            <TabsTrigger value="webhooks">Webhook</TabsTrigger>
            <TabsTrigger value="sync-jobs">同步任務</TabsTrigger>
            <TabsTrigger value="api-keys">API 金鑰</TabsTrigger>
            <TabsTrigger value="external-systems">外部系統</TabsTrigger>
            <TabsTrigger value="templates">整合模板</TabsTrigger>
            <TabsTrigger value="logs">記錄檔</TabsTrigger>
          </TabsList>

          <TabsContent value="overview" className="space-y-6">
            {/* Statistics Cards */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">總整合數</CardTitle>
                  <Link2 className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">{integrationStats?.total_integrations || 0}</div>
                  <p className="text-xs text-muted-foreground">
                    啟用: {integrationStats?.active_integrations || 0}
                  </p>
                </CardContent>
              </Card>
              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">同步任務</CardTitle>
                  <Sync className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">{integrationStats?.total_sync_jobs || 0}</div>
                  <p className="text-xs text-muted-foreground">
                    成功率: {integrationStats?.avg_success_rate?.toFixed(1) || 0}%
                  </p>
                </CardContent>
              </Card>
              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">Webhook 傳送</CardTitle>
                  <Webhook className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">{integrationStats?.webhook_deliveries || 0}</div>
                  <p className="text-xs text-muted-foreground">
                    本月傳送次數
                  </p>
                </CardContent>
              </Card>
              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">平均響應時間</CardTitle>
                  <Activity className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">{integrationStats?.avg_response_time?.toFixed(0) || 0}ms</div>
                  <p className="text-xs text-muted-foreground">
                    API 響應時間
                  </p>
                </CardContent>
              </Card>
            </div>

            {/* Charts and Recent Activity */}
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              {/* Integration Types */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <BarChart3 className="h-5 w-5" />
                    整合類型分布
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  {integrationsByType?.data && integrationsByType.data.length > 0 ? (
                    integrationsByType.data.map((item) => (
                      <div key={item.type} className="flex justify-between items-center">
                        <div className="flex items-center gap-2">
                          {getTypeIcon(item.type)}
                          <span className="text-sm capitalize">{item.type}</span>
                        </div>
                        <div className="flex items-center gap-2">
                          <div className="w-24 bg-gray-200 rounded-full h-2">
                            <div 
                              className="bg-blue-500 h-2 rounded-full" 
                              style={{ width: `${(item.count / (integrationStats?.total_integrations || 1)) * 100}%` }}
                            ></div>
                          </div>
                          <span className="text-sm font-medium">{item.count}</span>
                        </div>
                      </div>
                    ))
                  ) : (
                    <div className="text-center py-8 text-gray-500">
                      <BarChart3 className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                      <p>暫無數據</p>
                    </div>
                  )}
                </CardContent>
              </Card>

              {/* Recent Sync Jobs */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <History className="h-5 w-5" />
                    最近同步任務
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  {syncJobs?.data && syncJobs.data.length > 0 ? (
                    syncJobs.data.slice(0, 5).map((job) => (
                      <div key={job.id} className="flex justify-between items-start">
                        <div className="flex-1">
                          <p className="font-medium">{job.name}</p>
                          <p className="text-sm text-gray-500">{job.type} • {job.direction}</p>
                          <div className="flex items-center gap-2 mt-2">
                            {getPriorityBadge(job.priority)}
                            {getStatusBadge(job.status, 'sync')}
                          </div>
                          {job.status === 'running' && (
                            <Progress value={job.progress} className="h-2 mt-2" />
                          )}
                        </div>
                        <div className="text-right text-sm ml-4">
                          <p className="font-medium">
                            {job.processed_records}/{job.total_records}
                          </p>
                          <p className="text-gray-500">
                            {format(new Date(job.created_at), 'MM/dd HH:mm', { locale: zhTW })}
                          </p>
                        </div>
                      </div>
                    ))
                  ) : (
                    <div className="text-center py-8 text-gray-500">
                      <Sync className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                      <p>暫無同步任務</p>
                    </div>
                  )}
                </CardContent>
              </Card>
            </div>

            {/* System Health */}
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
              {/* Active Integrations */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Link2 className="h-5 w-5" />
                    啟用整合
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  {integrations?.data && integrations.data.length > 0 ? (
                    integrations.data.filter(i => i.is_active).slice(0, 5).map((integration) => (
                      <div key={integration.id} className="flex justify-between items-center">
                        <div className="flex items-center gap-2">
                          {getTypeIcon(integration.type)}
                          <div>
                            <p className="text-sm font-medium">{integration.name}</p>
                            <p className="text-xs text-gray-500">{integration.provider}</p>
                          </div>
                        </div>
                        <div className="text-right">
                          {getStatusBadge(integration.status, 'integration')}
                          <p className="text-xs text-gray-500 mt-1">
                            {integration.success_rate.toFixed(1)}% 成功率
                          </p>
                        </div>
                      </div>
                    ))
                  ) : (
                    <div className="text-center py-4 text-gray-500">
                      <p>暫無啟用的整合</p>
                    </div>
                  )}
                </CardContent>
              </Card>

              {/* Active Webhooks */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Webhook className="h-5 w-5" />
                    活躍 Webhook
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  {webhooks?.data && webhooks.data.length > 0 ? (
                    webhooks.data.filter(w => w.is_active).slice(0, 5).map((webhook) => (
                      <div key={webhook.id} className="flex justify-between items-center">
                        <div>
                          <p className="text-sm font-medium">{webhook.name}</p>
                          <p className="text-xs text-gray-500">{webhook.url}</p>
                        </div>
                        <div className="text-right">
                          <p className="text-xs font-medium">{webhook.trigger_count} 次觸發</p>
                          <p className="text-xs text-gray-500">
                            成功: {webhook.success_count}
                          </p>
                        </div>
                      </div>
                    ))
                  ) : (
                    <div className="text-center py-4 text-gray-500">
                      <p>暫無活躍的 Webhook</p>
                    </div>
                  )}
                </CardContent>
              </Card>

              {/* External Systems */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Server className="h-5 w-5" />
                    外部系統
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  {externalSystems?.data && externalSystems.data.length > 0 ? (
                    externalSystems.data.slice(0, 5).map((system) => (
                      <div key={system.id} className="flex justify-between items-center">
                        <div>
                          <p className="text-sm font-medium">{system.name}</p>
                          <p className="text-xs text-gray-500">{system.vendor} {system.version}</p>
                        </div>
                        <div className="text-right">
                          {getStatusBadge(system.status, 'system')}
                          <p className="text-xs text-gray-500 mt-1">
                            {system.last_test_at ? format(new Date(system.last_test_at), 'MM/dd HH:mm', { locale: zhTW }) : '未測試'}
                          </p>
                        </div>
                      </div>
                    ))
                  ) : (
                    <div className="text-center py-4 text-gray-500">
                      <p>暫無外部系統</p>
                    </div>
                  )}
                </CardContent>
              </Card>
            </div>
          </TabsContent>

          <TabsContent value="integrations">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <Link2 className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">整合管理功能開發中</p>
                  <p className="text-sm text-gray-400 mt-2">包含 API 整合配置、連線測試與狀態監控</p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="webhooks">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <Webhook className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">Webhook 管理功能開發中</p>
                  <p className="text-sm text-gray-400 mt-2">包含 Webhook 配置、事件觸發與傳送記錄</p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="sync-jobs">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <Sync className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">同步任務功能開發中</p>
                  <p className="text-sm text-gray-400 mt-2">包含數據同步配置、任務調度與進度監控</p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="api-keys">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <Key className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">API 金鑰管理功能開發中</p>
                  <p className="text-sm text-gray-400 mt-2">包含金鑰生成、權限控制與使用統計</p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="external-systems">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <Server className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">外部系統管理功能開發中</p>
                  <p className="text-sm text-gray-400 mt-2">包含系統註冊、連線測試與狀態監控</p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="templates">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <FileText className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">整合模板功能開發中</p>
                  <p className="text-sm text-gray-400 mt-2">包含預設模板、自訂配置與快速部署</p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="logs">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <History className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">整合記錄檔功能開發中</p>
                  <p className="text-sm text-gray-400 mt-2">包含操作記錄、錯誤追蹤與效能分析</p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </div>
    </DashboardLayout>
  )
}