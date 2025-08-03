'use client'

import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import DashboardLayout from '@/components/layout/DashboardLayout'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
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
  Settings,
  Users,
  Shield,
  Database,
  Activity,
  Bell,
  HardDrive,
  Zap,
  CheckCircle,
  AlertCircle,
  XCircle,
  Clock,
  Wrench,
  Download,
  Upload,
  RefreshCw,
  Save,
  Eye,
  Edit,
  Trash2,
  Plus,
  Search,
  Filter,
  Server,
  Cpu,
  MemoryStick,
  Network,
  Calendar,
  Play,
  Pause,
  Square
} from 'lucide-react'
import systemService from '@/services/system.service'
import { format } from 'date-fns'
import { zhTW } from 'date-fns/locale'

export default function SystemManagementPage() {
  const [activeTab, setActiveTab] = useState('overview')
  const [searchQuery, setSearchQuery] = useState('')
  const [categoryFilter, setCategoryFilter] = useState('')
  const [page, setPage] = useState(1)
  const pageSize = 20

  // Fetch system statistics
  const { data: systemStats, isLoading: isLoadingStats } = useQuery({
    queryKey: ['system-statistics'],
    queryFn: () => systemService.getUserStatistics(),
  })

  // Fetch system health
  const { data: systemHealth, isLoading: isLoadingHealth, refetch: refetchHealth } = useQuery({
    queryKey: ['system-health'],
    queryFn: () => systemService.getSystemHealth(),
  })

  // Fetch system configs
  const { data: systemConfigs, isLoading: isLoadingConfigs } = useQuery({
    queryKey: ['system-configs', page, searchQuery, categoryFilter],
    queryFn: () => systemService.listSystemConfigs({
      page,
      page_size: pageSize,
      search: searchQuery || undefined,
      category: categoryFilter || undefined,
    }),
  })

  // Fetch roles
  const { data: roles, isLoading: isLoadingRoles } = useQuery({
    queryKey: ['system-roles'],
    queryFn: () => systemService.listRoles({ page: 1, page_size: 100 }),
  })

  // Fetch backup records
  const { data: backupRecords } = useQuery({
    queryKey: ['backup-records'],
    queryFn: () => ({ data: [] }),
  })

  // Fetch system tasks
  const { data: systemTasks } = useQuery({
    queryKey: ['system-tasks'],
    queryFn: () => ({ data: [] }),
  })

  const getHealthStatusBadge = (status: string) => {
    const statusConfig: Record<string, { label: string; variant: any; icon: any }> = {
      healthy: { label: '健康', variant: 'success', icon: CheckCircle },
      warning: { label: '警告', variant: 'warning', icon: AlertCircle },
      critical: { label: '嚴重', variant: 'destructive', icon: XCircle },
      down: { label: '離線', variant: 'secondary', icon: XCircle },
    }

    const config = statusConfig[status] || { label: status, variant: 'default', icon: AlertCircle }
    const Icon = config.icon
    
    return (
      <Badge variant={config.variant as any} className="flex items-center gap-1">
        <Icon className="h-3 w-3" />
        {config.label}
      </Badge>
    )
  }

  const getTaskStatusBadge = (status: string) => {
    const statusConfig: Record<string, { label: string; variant: any; icon: any }> = {
      pending: { label: '等待中', variant: 'secondary', icon: Clock },
      running: { label: '執行中', variant: 'warning', icon: RefreshCw },
      completed: { label: '已完成', variant: 'success', icon: CheckCircle },
      failed: { label: '失敗', variant: 'destructive', icon: XCircle },
      cancelled: { label: '已取消', variant: 'secondary', icon: Square },
    }

    const config = statusConfig[status] || { label: status, variant: 'default', icon: AlertCircle }
    const Icon = config.icon
    
    return (
      <Badge variant={config.variant as any} className="flex items-center gap-1">
        <Icon className="h-3 w-3" />
        {config.label}
      </Badge>
    )
  }

  const handleCheckSystemHealth = async () => {
    try {
      await systemService.getSystemHealth()
      refetchHealth()
    } catch (error) {
      console.error('Failed to check system health:', error)
    }
  }

  const handlePerformBackup = async () => {
    try {
      console.log('Backup performed')
      // Show success message and refresh data
    } catch (error) {
      console.error('Failed to perform backup:', error)
    }
  }

  if (isLoadingStats || isLoadingHealth) {
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
            <h1 className="text-3xl font-bold text-gray-900">系統管理</h1>
            <p className="mt-2 text-gray-600">系統設定、監控、備份與維護管理</p>
          </div>
          <div className="flex gap-2">
            <Button variant="outline" onClick={handleCheckSystemHealth}>
              <Activity className="mr-2 h-4 w-4" />
              健康檢查
            </Button>
            <Button variant="outline" onClick={handlePerformBackup}>
              <HardDrive className="mr-2 h-4 w-4" />
              執行備份
            </Button>
          </div>
        </div>

        <Tabs value={activeTab} onValueChange={setActiveTab}>
          <TabsList>
            <TabsTrigger value="overview">系統概覽</TabsTrigger>
            <TabsTrigger value="config">系統設定</TabsTrigger>
            <TabsTrigger value="roles">角色權限</TabsTrigger>
            <TabsTrigger value="health">系統監控</TabsTrigger>
            <TabsTrigger value="backup">備份管理</TabsTrigger>
            <TabsTrigger value="tasks">任務管理</TabsTrigger>
            <TabsTrigger value="audit">稽核日誌</TabsTrigger>
          </TabsList>

          <TabsContent value="overview" className="space-y-6">
            {/* System Statistics */}
            {systemStats && (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                <Card>
                  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">系統使用者</CardTitle>
                    <Users className="h-4 w-4 text-muted-foreground" />
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold">{systemStats.total_users}</div>
                    <p className="text-xs text-muted-foreground">
                      活躍會話: {(systemStats as any).active_sessions || 0}
                    </p>
                  </CardContent>
                </Card>
                <Card>
                  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">未讀通知</CardTitle>
                    <Bell className="h-4 w-4 text-muted-foreground" />
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold">{(systemStats as any).unread_notifications || 0}</div>
                    <p className="text-xs text-muted-foreground">
                      系統通知數量
                    </p>
                  </CardContent>
                </Card>
                <Card>
                  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">待處理任務</CardTitle>
                    <Zap className="h-4 w-4 text-muted-foreground" />
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold">{(systemStats as any).pending_tasks || 0}</div>
                    <p className="text-xs text-muted-foreground">
                      背景處理任務
                    </p>
                  </CardContent>
                </Card>
                <Card>
                  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">系統狀態</CardTitle>
                    <Activity className="h-4 w-4 text-muted-foreground" />
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold text-green-600">正常</div>
                    <p className="text-xs text-muted-foreground">
                      系統運行狀態
                    </p>
                  </CardContent>
                </Card>
              </div>
            )}

            {/* System Health Overview */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Server className="h-5 w-5" />
                  系統健康狀態
                </CardTitle>
              </CardHeader>
              <CardContent>
                {systemHealth && systemHealth.length > 0 ? (
                  <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {systemHealth.map((health: any) => (
                      <div key={health.id} className="p-4 border rounded-lg">
                        <div className="flex justify-between items-center mb-2">
                          <h4 className="font-medium capitalize">{health.component}</h4>
                          {getHealthStatusBadge(health.status)}
                        </div>
                        <div className="space-y-1 text-sm text-gray-600">
                          <div className="flex justify-between">
                            <span>回應時間:</span>
                            <span>{health.response_time.toFixed(1)}ms</span>
                          </div>
                          {health.cpu_usage > 0 && (
                            <div className="flex justify-between">
                              <span>CPU 使用率:</span>
                              <span>{health.cpu_usage.toFixed(1)}%</span>
                            </div>
                          )}
                          {health.memory_usage > 0 && (
                            <div className="flex justify-between">
                              <span>記憶體使用率:</span>
                              <span>{health.memory_usage.toFixed(1)}%</span>
                            </div>
                          )}
                        </div>
                        <p className="text-xs text-gray-500 mt-2">{health.message}</p>
                      </div>
                    ))}
                  </div>
                ) : (
                  <div className="text-center py-8">
                    <Server className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                    <p className="text-gray-500">暫無健康監控資料</p>
                  </div>
                )}
              </CardContent>
            </Card>

            {/* Recent Activities */}
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              {/* Recent Backups */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <HardDrive className="h-5 w-5" />
                    最近備份
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  {backupRecords?.data.length === 0 ? (
                    <p className="text-center text-gray-500 py-4">暫無備份記錄</p>
                  ) : (
                    <div className="space-y-3">
                      {backupRecords?.data.slice(0, 5).map((backup: any) => (
                        <div key={backup.id} className="flex items-center justify-between p-3 border rounded-lg">
                          <div>
                            <p className="font-medium">{backup.backup_name}</p>
                            <p className="text-sm text-gray-500">{backup.type}</p>
                          </div>
                          <div className="text-right">
                            <p className="text-sm">{getTaskStatusBadge(backup.status)}</p>
                            <p className="text-xs text-gray-500">
                              {format(new Date(backup.started_at), 'MM/dd HH:mm', { locale: zhTW })}
                            </p>
                          </div>
                        </div>
                      ))}
                    </div>
                  )}
                </CardContent>
              </Card>

              {/* Recent Tasks */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Zap className="h-5 w-5" />
                    最近任務
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  {systemTasks?.data.length === 0 ? (
                    <p className="text-center text-gray-500 py-4">暫無系統任務</p>
                  ) : (
                    <div className="space-y-3">
                      {systemTasks?.data.slice(0, 5).map((task: any) => (
                        <div key={task.id} className="flex items-center justify-between p-3 border rounded-lg">
                          <div>
                            <p className="font-medium">{task.name}</p>
                            <p className="text-sm text-gray-500">{task.type}</p>
                          </div>
                          <div className="text-right">
                            <p className="text-sm">{getTaskStatusBadge(task.status)}</p>
                            <p className="text-xs text-gray-500">
                              {format(new Date(task.created_at), 'MM/dd HH:mm', { locale: zhTW })}
                            </p>
                          </div>
                        </div>
                      ))}
                    </div>
                  )}
                </CardContent>
              </Card>
            </div>
          </TabsContent>

          <TabsContent value="config">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <Settings className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">系統設定功能開發中</p>
                  <p className="text-sm text-gray-400 mt-2">包含郵件設定、安全設定、整合設定等</p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="roles">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <Shield className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">角色權限管理功能開發中</p>
                  <p className="text-sm text-gray-400 mt-2">管理使用者角色與權限設定</p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="health">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center justify-between">
                  <span className="flex items-center gap-2">
                    <Activity className="h-5 w-5" />
                    系統監控
                  </span>
                  <Button variant="outline" size="sm" onClick={handleCheckSystemHealth}>
                    <RefreshCw className="h-4 w-4 mr-2" />
                    重新檢查
                  </Button>
                </CardTitle>
              </CardHeader>
              <CardContent>
                {systemHealth && systemHealth.length > 0 ? (
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead>組件</TableHead>
                        <TableHead>狀態</TableHead>
                        <TableHead>回應時間</TableHead>
                        <TableHead>CPU 使用率</TableHead>
                        <TableHead>記憶體使用率</TableHead>
                        <TableHead>檢查時間</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {systemHealth.map((health: any) => (
                        <TableRow key={health.id}>
                          <TableCell className="font-medium capitalize">{health.component}</TableCell>
                          <TableCell>{getHealthStatusBadge(health.status)}</TableCell>
                          <TableCell>{health.response_time.toFixed(1)}ms</TableCell>
                          <TableCell>{health.cpu_usage > 0 ? `${health.cpu_usage.toFixed(1)}%` : '-'}</TableCell>
                          <TableCell>{health.memory_usage > 0 ? `${health.memory_usage.toFixed(1)}%` : '-'}</TableCell>
                          <TableCell>
                            {format(new Date(health.checked_at), 'MM/dd HH:mm:ss', { locale: zhTW })}
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                ) : (
                  <div className="text-center py-8">
                    <Activity className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                    <p className="text-gray-500">暫無監控資料</p>
                  </div>
                )}
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="backup">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <HardDrive className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">備份管理功能開發中</p>
                  <p className="text-sm text-gray-400 mt-2">系統備份、還原與備份策略管理</p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="tasks">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <Zap className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">任務管理功能開發中</p>
                  <p className="text-sm text-gray-400 mt-2">背景任務監控與排程管理</p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="audit">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <Database className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">稽核日誌功能開發中</p>
                  <p className="text-sm text-gray-400 mt-2">系統操作記錄與安全稽核</p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </div>
    </DashboardLayout>
  )
}