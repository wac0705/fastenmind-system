'use client'

import { useState, useEffect } from 'react'
import { useQuery } from '@tanstack/react-query'
import DashboardLayout from '@/components/layout/DashboardLayout'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Progress } from '@/components/ui/progress'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { 
  Brain,
  Lightbulb,
  Search,
  Layers,
  Settings2,
  Shield,
  Activity,
  Database,
  Globe,
  Zap,
  TrendingUp,
  Clock,
  AlertTriangle,
  CheckCircle,
  XCircle,
  Play,
  Eye,
  Download,
  RefreshCw,
  BarChart3,
  Users,
  FileText,
  Cpu,
  HardDrive,
  Network,
  Lock,
  Sparkles,
  MessageSquare,
  Target,
  Filter,
  Package,
  Wrench
} from 'lucide-react'
import advancedService from '@/services/advanced.service'
import { format } from 'date-fns'
import { zhTW } from 'date-fns/locale'

export default function AdvancedPage() {
  const [activeTab, setActiveTab] = useState('overview')

  // Fetch analytics data
  const { data: aiUsageStats, isLoading: isLoadingAI } = useQuery({
    queryKey: ['ai-usage-stats'],
    queryFn: () => advancedService.getAIUsageStats(),
  })

  const { data: recommendationStats, isLoading: isLoadingRecommendations } = useQuery({
    queryKey: ['recommendation-stats'],
    queryFn: () => advancedService.getRecommendationStats(),
  })

  const { data: securityEventStats, isLoading: isLoadingSecurity } = useQuery({
    queryKey: ['security-event-stats'],
    queryFn: () => advancedService.getSecurityEventStats(),
  })

  // Fetch recent items
  const { data: aiAssistants } = useQuery({
    queryKey: ['ai-assistants'],
    queryFn: () => advancedService.listAIAssistants({ limit: 5 }),
  })

  const { data: recommendations } = useQuery({
    queryKey: ['recommendations'],
    queryFn: () => advancedService.listRecommendations({ limit: 10 }),
  })

  const { data: batchOperations } = useQuery({
    queryKey: ['batch-operations'],
    queryFn: () => advancedService.listBatchOperations({ limit: 5 }),
  })

  const { data: securityEvents } = useQuery({
    queryKey: ['security-events'],
    queryFn: () => advancedService.listSecurityEvents({ limit: 10 }),
  })

  const { data: backups } = useQuery({
    queryKey: ['backups'],
    queryFn: () => advancedService.listBackups({ limit: 5 }),
  })

  const getStatusBadge = (status: string, type: 'operation' | 'security' | 'backup' | 'recommendation') => {
    const statusConfig: Record<string, Record<string, { label: string; variant: any; icon: any }>> = {
      operation: {
        pending: { label: '等待中', variant: 'secondary', icon: Clock },
        running: { label: '執行中', variant: 'info', icon: Play },
        completed: { label: '已完成', variant: 'success', icon: CheckCircle },
        failed: { label: '失敗', variant: 'destructive', icon: XCircle },
        cancelled: { label: '已取消', variant: 'secondary', icon: XCircle },
      },
      security: {
        new: { label: '新事件', variant: 'destructive', icon: AlertTriangle },
        investigating: { label: '調查中', variant: 'warning', icon: Search },
        resolved: { label: '已解決', variant: 'success', icon: CheckCircle },
        false_positive: { label: '誤報', variant: 'secondary', icon: XCircle },
      },
      backup: {
        running: { label: '備份中', variant: 'info', icon: Play },
        completed: { label: '已完成', variant: 'success', icon: CheckCircle },
        failed: { label: '失敗', variant: 'destructive', icon: XCircle },
        cancelled: { label: '已取消', variant: 'secondary', icon: XCircle },
      },
      recommendation: {
        pending: { label: '待處理', variant: 'secondary', icon: Clock },
        viewed: { label: '已查看', variant: 'info', icon: Eye },
        accepted: { label: '已接受', variant: 'success', icon: CheckCircle },
        rejected: { label: '已拒絕', variant: 'destructive', icon: XCircle },
        implemented: { label: '已實施', variant: 'success', icon: CheckCircle },
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

  const getPriorityBadge = (priority: string) => {
    const priorityConfig: Record<string, { label: string; variant: any; className: string }> = {
      low: { label: '低', variant: 'secondary', className: 'text-gray-600' },
      medium: { label: '中', variant: 'warning', className: 'text-yellow-600' },
      high: { label: '高', variant: 'destructive', className: 'text-orange-600' },
      urgent: { label: '緊急', variant: 'destructive', className: 'text-red-600 animate-pulse' },
    }

    const config = priorityConfig[priority] || { label: priority, variant: 'default', className: '' }
    
    return (
      <Badge variant={config.variant as any} className={config.className}>
        {config.label}
      </Badge>
    )
  }

  const getSeverityBadge = (severity: string) => {
    const severityConfig: Record<string, { label: string; variant: any; className: string }> = {
      low: { label: '低', variant: 'success', className: 'text-green-600' },
      medium: { label: '中', variant: 'warning', className: 'text-yellow-600' },
      high: { label: '高', variant: 'destructive', className: 'text-orange-600' },
      critical: { label: '嚴重', variant: 'destructive', className: 'text-red-600 font-bold' },
    }

    const config = severityConfig[severity] || { label: severity, variant: 'default', className: '' }
    
    return (
      <Badge variant={config.variant as any} className={config.className}>
        {config.label}
      </Badge>
    )
  }

  if (isLoadingAI || isLoadingRecommendations || isLoadingSecurity) {
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
            <h1 className="text-3xl font-bold text-gray-900">進階功能</h1>
            <p className="mt-2 text-gray-600">AI 助手、智能推薦、高級搜索與系統優化</p>
          </div>
          <div className="flex items-center gap-4">
            <Button variant="outline">
              <RefreshCw className="mr-2 h-4 w-4" />
              重新載入
            </Button>
            <Button>
              <Settings2 className="mr-2 h-4 w-4" />
              系統設定
            </Button>
          </div>
        </div>

        <Tabs value={activeTab} onValueChange={setActiveTab}>
          <TabsList className="grid w-full grid-cols-8">
            <TabsTrigger value="overview">概覽</TabsTrigger>
            <TabsTrigger value="ai-assistant">AI 助手</TabsTrigger>
            <TabsTrigger value="recommendations">智能推薦</TabsTrigger>
            <TabsTrigger value="search">高級搜索</TabsTrigger>
            <TabsTrigger value="batch-ops">批量操作</TabsTrigger>
            <TabsTrigger value="security">安全監控</TabsTrigger>
            <TabsTrigger value="performance">效能監控</TabsTrigger>
            <TabsTrigger value="system">系統管理</TabsTrigger>
          </TabsList>

          <TabsContent value="overview" className="space-y-6">
            {/* Statistics Cards */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">AI 助手使用次數</CardTitle>
                  <Brain className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">{aiUsageStats?.total_sessions || 0}</div>
                  <p className="text-xs text-muted-foreground">
                    活躍助手: {aiUsageStats?.active_assistants || 0}
                  </p>
                </CardContent>
              </Card>
              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">智能推薦</CardTitle>
                  <Lightbulb className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">{recommendations?.data?.length || 0}</div>
                  <p className="text-xs text-muted-foreground">
                    待處理: {recommendations?.data?.filter(r => r.status === 'pending').length || 0}
                  </p>
                </CardContent>
              </Card>
              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">安全事件</CardTitle>
                  <Shield className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">{securityEvents?.data?.length || 0}</div>
                  <p className="text-xs text-muted-foreground">
                    高風險: {securityEvents?.data?.filter(e => e.severity === 'high' || e.severity === 'critical').length || 0}
                  </p>
                </CardContent>
              </Card>
              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">批量操作</CardTitle>
                  <Layers className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">{batchOperations?.data?.length || 0}</div>
                  <p className="text-xs text-muted-foreground">
                    執行中: {batchOperations?.data?.filter(op => op.status === 'running').length || 0}
                  </p>
                </CardContent>
              </Card>
            </div>

            {/* Recent Activity */}
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              {/* Recent AI Conversations */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <MessageSquare className="h-5 w-5" />
                    最近 AI 對話
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  {aiAssistants?.data && aiAssistants.data.length > 0 ? (
                    <div className="space-y-4">
                      {aiAssistants.data.map((assistant) => (
                        <div key={assistant.id} className="flex justify-between items-center">
                          <div>
                            <p className="font-medium">{assistant.name}</p>
                            <p className="text-sm text-gray-500">
                              {assistant.type} • 使用 {assistant.usage_count} 次
                            </p>
                          </div>
                          <div className="text-right text-sm">
                            <p className="font-medium">${assistant.cost_accumulated.toFixed(4)}</p>
                            <p className="text-gray-500">
                              {assistant.last_used ? format(new Date(assistant.last_used), 'MM/dd HH:mm', { locale: zhTW }) : '未使用'}
                            </p>
                          </div>
                        </div>
                      ))}
                    </div>
                  ) : (
                    <div className="text-center py-8 text-gray-500">
                      <Brain className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                      <p>暫無 AI 助手</p>
                    </div>
                  )}
                </CardContent>
              </Card>

              {/* Recent Recommendations */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Target className="h-5 w-5" />
                    最新推薦
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  {recommendations?.data && recommendations.data.length > 0 ? (
                    <div className="space-y-4">
                      {recommendations.data.slice(0, 5).map((recommendation) => (
                        <div key={recommendation.id} className="flex justify-between items-start">
                          <div className="flex-1">
                            <p className="font-medium">{recommendation.title}</p>
                            <p className="text-sm text-gray-500 mt-1">{recommendation.description}</p>
                            <div className="flex items-center gap-2 mt-2">
                              {getPriorityBadge(recommendation.priority)}
                              {getStatusBadge(recommendation.status, 'recommendation')}
                            </div>
                          </div>
                          <div className="text-right text-sm ml-4">
                            <p className="font-medium">信心度: {(recommendation.score * 100).toFixed(0)}%</p>
                            <p className="text-gray-500">
                              {format(new Date(recommendation.created_at), 'MM/dd HH:mm', { locale: zhTW })}
                            </p>
                          </div>
                        </div>
                      ))}
                    </div>
                  ) : (
                    <div className="text-center py-8 text-gray-500">
                      <Lightbulb className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                      <p>暫無推薦</p>
                    </div>
                  )}
                </CardContent>
              </Card>
            </div>

            {/* System Health */}
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
              {/* Batch Operations Status */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Layers className="h-5 w-5" />
                    批量操作狀態
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  {batchOperations?.data && batchOperations.data.length > 0 ? (
                    batchOperations.data.map((operation) => (
                      <div key={operation.id} className="space-y-2">
                        <div className="flex justify-between items-center">
                          <p className="text-sm font-medium">{operation.operation_type}</p>
                          {getStatusBadge(operation.status, 'operation')}
                        </div>
                        <Progress value={operation.progress} className="h-2" />
                        <p className="text-xs text-gray-500">
                          {operation.processed_items}/{operation.total_items} 項目
                        </p>
                      </div>
                    ))
                  ) : (
                    <div className="text-center py-4 text-gray-500">
                      <p>暫無批量操作</p>
                    </div>
                  )}
                </CardContent>
              </Card>

              {/* Security Events */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Shield className="h-5 w-5" />
                    安全事件
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  {securityEvents?.data && securityEvents.data.length > 0 ? (
                    securityEvents.data.slice(0, 3).map((event) => (
                      <div key={event.id} className="flex justify-between items-start">
                        <div>
                          <p className="text-sm font-medium">{event.event_type}</p>
                          <p className="text-xs text-gray-500">{event.description}</p>
                        </div>
                        <div className="text-right">
                          {getSeverityBadge(event.severity)}
                          <p className="text-xs text-gray-500 mt-1">
                            {format(new Date(event.created_at), 'MM/dd HH:mm', { locale: zhTW })}
                          </p>
                        </div>
                      </div>
                    ))
                  ) : (
                    <div className="text-center py-4 text-gray-500">
                      <p>暫無安全事件</p>
                    </div>
                  )}
                </CardContent>
              </Card>

              {/* Backup Status */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Database className="h-5 w-5" />
                    備份狀態
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  {backups?.data && backups.data.length > 0 ? (
                    backups.data.map((backup) => (
                      <div key={backup.id} className="flex justify-between items-center">
                        <div>
                          <p className="text-sm font-medium">{backup.backup_type}</p>
                          <p className="text-xs text-gray-500">
                            {(backup.file_size / (1024 * 1024)).toFixed(1)} MB
                          </p>
                        </div>
                        <div className="text-right">
                          {getStatusBadge(backup.status, 'backup')}
                          <p className="text-xs text-gray-500 mt-1">
                            {format(new Date(backup.created_at), 'MM/dd HH:mm', { locale: zhTW })}
                          </p>
                        </div>
                      </div>
                    ))
                  ) : (
                    <div className="text-center py-4 text-gray-500">
                      <p>暫無備份記錄</p>
                    </div>
                  )}
                </CardContent>
              </Card>
            </div>
          </TabsContent>

          <TabsContent value="ai-assistant">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <Brain className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">AI 助手管理功能開發中</p>
                  <p className="text-sm text-gray-400 mt-2">包含對話管理、模型配置與使用分析</p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="recommendations">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <Lightbulb className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">智能推薦系統開發中</p>
                  <p className="text-sm text-gray-400 mt-2">包含推薦引擎、個人化建議與效果追蹤</p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="search">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <Search className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">高級搜索功能開發中</p>
                  <p className="text-sm text-gray-400 mt-2">包含複雜查詢、搜索模板與結果導出</p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="batch-ops">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <Layers className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">批量操作功能開發中</p>
                  <p className="text-sm text-gray-400 mt-2">包含批量更新、刪除、導入導出與進度追蹤</p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="security">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <Shield className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">安全監控功能開發中</p>
                  <p className="text-sm text-gray-400 mt-2">包含威脅檢測、異常行為分析與風險評估</p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="performance">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <Activity className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">效能監控功能開發中</p>
                  <p className="text-sm text-gray-400 mt-2">包含系統監控、效能分析與優化建議</p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="system">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <Wrench className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">系統管理功能開發中</p>
                  <p className="text-sm text-gray-400 mt-2">包含自訂欄位、多語言支援、備份恢復與系統配置</p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </div>
    </DashboardLayout>
  )
}