'use client'

import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { useRouter } from 'next/navigation'
import DashboardLayout from '@/components/layout/DashboardLayout'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Progress } from '@/components/ui/progress'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { 
  Users,
  Shield,
  Settings,
  History,
  Bell,
  Activity,
  Database,
  HardDrive,
  Cpu,
  MemoryStick,
  Wifi,
  Globe,
  Lock,
  Key,
  FileText,
  AlertTriangle,
  CheckCircle,
  XCircle,
  RefreshCw,
  Download,
  Upload,
  Zap,
  Clock,
  Calendar,
  TrendingUp,
  TrendingDown,
  BarChart3,
  PieChart,
  Server,
  Cloud,
  GitBranch
} from 'lucide-react'
import systemService from '@/services/system.service'
import { format } from 'date-fns'
import { zhTW } from 'date-fns/locale'

export default function SystemManagementPage() {
  const router = useRouter()
  const [activeTab, setActiveTab] = useState('overview')

  // Fetch system info
  const { data: systemInfo, isLoading: isLoadingInfo } = useQuery({
    queryKey: ['system-info'],
    queryFn: () => systemService.getSystemInfo(),
  })

  // Fetch system health
  const { data: systemHealth, isLoading: isLoadingHealth } = useQuery({
    queryKey: ['system-health'],
    queryFn: () => systemService.getSystemHealth(),
    refetchInterval: 30000, // Refresh every 30 seconds
  })

  // Fetch system metrics
  const { data: systemMetrics } = useQuery({
    queryKey: ['system-metrics'],
    queryFn: () => systemService.getSystemMetrics(),
    refetchInterval: 10000, // Refresh every 10 seconds
  })

  // Fetch user statistics
  const { data: userStats } = useQuery({
    queryKey: ['user-stats'],
    queryFn: () => systemService.getUserStatistics(),
  })

  // Fetch cache statistics
  const { data: cacheStats } = useQuery({
    queryKey: ['cache-stats'],
    queryFn: () => systemService.getCacheStatistics(),
  })

  // Fetch license info
  const { data: licenseInfo } = useQuery({
    queryKey: ['license-info'],
    queryFn: () => systemService.getLicenseInfo(),
  })

  const getHealthStatus = (status: string) => {
    switch (status) {
      case 'healthy':
        return { icon: CheckCircle, color: 'text-green-600', label: '健康' }
      case 'warning':
        return { icon: AlertTriangle, color: 'text-yellow-600', label: '警告' }
      case 'error':
        return { icon: XCircle, color: 'text-red-600', label: '錯誤' }
      default:
        return { icon: Activity, color: 'text-gray-600', label: '未知' }
    }
  }

  const handleClearCache = async () => {
    try {
      await systemService.clearCache()
      // Refresh cache stats
    } catch (error) {
      console.error('Failed to clear cache:', error)
    }
  }

  const handleRunDiagnostics = async () => {
    try {
      const result = await systemService.runSystemDiagnostics()
      console.log('Diagnostics result:', result)
    } catch (error) {
      console.error('Failed to run diagnostics:', error)
    }
  }

  const navigationCards = [
    {
      title: '使用者管理',
      description: '管理系統使用者帳號與權限',
      icon: Users,
      href: '/system/users',
      stats: userStats?.total_users || 0,
      statsLabel: '使用者',
    },
    {
      title: '角色權限',
      description: '設定角色與權限控制',
      icon: Shield,
      href: '/system/roles',
      stats: '12',
      statsLabel: '角色',
    },
    {
      title: '系統設定',
      description: '管理系統參數與配置',
      icon: Settings,
      href: '/system/settings',
      stats: '48',
      statsLabel: '設定項',
    },
    {
      title: '操作記錄',
      description: '查看系統操作日誌',
      icon: History,
      href: '/system/audit-logs',
      stats: '1.2k',
      statsLabel: '記錄',
    },
    {
      title: '系統通知',
      description: '管理系統通知與提醒',
      icon: Bell,
      href: '/system/notifications',
      stats: '8',
      statsLabel: '待處理',
    },
    {
      title: '系統健康',
      description: '監控系統運行狀態',
      icon: Activity,
      href: '/system/health',
      stats: systemHealth?.status || 'healthy',
      statsLabel: '狀態',
    },
  ]

  if (isLoadingInfo || isLoadingHealth) {
    return (
      <DashboardLayout>
        <div className="flex items-center justify-center h-64">
          <div className="text-center">
            <RefreshCw className="h-8 w-8 animate-spin mx-auto mb-4 text-gray-400" />
            <p className="text-gray-500">載入中...</p>
          </div>
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
            <p className="mt-2 text-gray-600">監控系統狀態、管理使用者權限與系統設定</p>
          </div>
          <div className="flex gap-2">
            <Button variant="outline" onClick={handleRunDiagnostics}>
              <Zap className="mr-2 h-4 w-4" />
              系統診斷
            </Button>
            <Button variant="outline" onClick={handleClearCache}>
              <RefreshCw className="mr-2 h-4 w-4" />
              清除快取
            </Button>
          </div>
        </div>

        <Tabs value={activeTab} onValueChange={setActiveTab}>
          <TabsList>
            <TabsTrigger value="overview">總覽</TabsTrigger>
            <TabsTrigger value="performance">效能監控</TabsTrigger>
            <TabsTrigger value="storage">儲存空間</TabsTrigger>
            <TabsTrigger value="network">網路狀態</TabsTrigger>
            <TabsTrigger value="license">授權資訊</TabsTrigger>
          </TabsList>

          <TabsContent value="overview" className="space-y-6">
            {/* System Status */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
              <Card>
                <CardHeader className="pb-3">
                  <CardTitle className="text-sm flex items-center gap-2">
                    <Server className="h-4 w-4" />
                    系統狀態
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="flex items-center gap-2">
                    {systemHealth && (
                      <>
                        {(() => {
                          const Icon = getHealthStatus(systemHealth.status).icon;
                          return <Icon className={`h-5 w-5 ${getHealthStatus(systemHealth.status).color}`} />;
                        })()}
                        <span className={`font-medium ${getHealthStatus(systemHealth.status).color}`}>
                          {getHealthStatus(systemHealth.status).label}
                        </span>
                      </>
                    )}
                  </div>
                  <p className="text-sm text-gray-500 mt-1">
                    執行時間: {systemInfo?.uptime || '0'} 小時
                  </p>
                </CardContent>
              </Card>

              <Card>
                <CardHeader className="pb-3">
                  <CardTitle className="text-sm flex items-center gap-2">
                    <Cpu className="h-4 w-4" />
                    CPU 使用率
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="flex items-center justify-between">
                    <span className="text-2xl font-bold">
                      {systemMetrics?.cpu_usage?.toFixed(1) || 0}%
                    </span>
                    {systemMetrics?.cpu_trend > 0 ? (
                      <TrendingUp className="h-5 w-5 text-red-600" />
                    ) : (
                      <TrendingDown className="h-5 w-5 text-green-600" />
                    )}
                  </div>
                  <Progress value={systemMetrics?.cpu_usage || 0} className="mt-2" />
                </CardContent>
              </Card>

              <Card>
                <CardHeader className="pb-3">
                  <CardTitle className="text-sm flex items-center gap-2">
                    <MemoryStick className="h-4 w-4" />
                    記憶體使用
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="flex items-center justify-between">
                    <span className="text-2xl font-bold">
                      {systemMetrics?.memory_usage?.toFixed(1) || 0}%
                    </span>
                    <span className="text-sm text-gray-500">
                      {systemMetrics?.memory_used || 0}GB / {systemMetrics?.memory_total || 0}GB
                    </span>
                  </div>
                  <Progress value={systemMetrics?.memory_usage || 0} className="mt-2" />
                </CardContent>
              </Card>

              <Card>
                <CardHeader className="pb-3">
                  <CardTitle className="text-sm flex items-center gap-2">
                    <HardDrive className="h-4 w-4" />
                    儲存空間
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="flex items-center justify-between">
                    <span className="text-2xl font-bold">
                      {systemMetrics?.disk_usage?.toFixed(1) || 0}%
                    </span>
                    <span className="text-sm text-gray-500">
                      {systemMetrics?.disk_used || 0}GB / {systemMetrics?.disk_total || 0}GB
                    </span>
                  </div>
                  <Progress value={systemMetrics?.disk_usage || 0} className="mt-2" />
                </CardContent>
              </Card>
            </div>

            {/* Navigation Cards */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {navigationCards.map((card) => {
                const Icon = card.icon
                return (
                  <Card
                    key={card.href}
                    className="hover:shadow-lg transition-shadow cursor-pointer"
                    onClick={() => router.push(card.href)}
                  >
                    <CardHeader>
                      <div className="flex items-center justify-between">
                        <Icon className="h-8 w-8 text-gray-400" />
                        <div className="text-right">
                          <p className="text-2xl font-bold">{card.stats}</p>
                          <p className="text-sm text-gray-500">{card.statsLabel}</p>
                        </div>
                      </div>
                      <CardTitle className="mt-4">{card.title}</CardTitle>
                      <CardDescription>{card.description}</CardDescription>
                    </CardHeader>
                  </Card>
                )
              })}
            </div>

            {/* System Information */}
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Server className="h-5 w-5" />
                    系統資訊
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-3">
                  <div className="flex justify-between">
                    <span className="text-gray-500">作業系統</span>
                    <span className="font-medium">{systemInfo?.os || 'Unknown'}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-500">系統版本</span>
                    <span className="font-medium">{systemInfo?.version || '1.0.0'}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-500">主機名稱</span>
                    <span className="font-medium">{systemInfo?.hostname || 'localhost'}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-500">系統架構</span>
                    <span className="font-medium">{systemInfo?.architecture || 'x64'}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-500">CPU 核心數</span>
                    <span className="font-medium">{systemInfo?.cpu_cores || 0}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-500">總記憶體</span>
                    <span className="font-medium">{systemInfo?.total_memory || 0} GB</span>
                  </div>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Database className="h-5 w-5" />
                    資料庫狀態
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-3">
                  <div className="flex justify-between">
                    <span className="text-gray-500">資料庫類型</span>
                    <span className="font-medium">{systemInfo?.database_type || 'PostgreSQL'}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-500">資料庫版本</span>
                    <span className="font-medium">{systemInfo?.database_version || '14.0'}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-500">連線數</span>
                    <span className="font-medium">{systemMetrics?.db_connections || 0}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-500">資料庫大小</span>
                    <span className="font-medium">{systemMetrics?.db_size || 0} GB</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-500">快取命中率</span>
                    <span className="font-medium">{cacheStats?.hit_rate?.toFixed(1) || 0}%</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-500">平均查詢時間</span>
                    <span className="font-medium">{systemMetrics?.avg_query_time || 0} ms</span>
                  </div>
                </CardContent>
              </Card>
            </div>
          </TabsContent>

          <TabsContent value="performance">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <BarChart3 className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">效能監控功能開發中</p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="storage">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <HardDrive className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">儲存空間管理功能開發中</p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="network">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <Wifi className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">網路狀態監控功能開發中</p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="license">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Key className="h-5 w-5" />
                  授權資訊
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                {licenseInfo ? (
                  <>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div>
                        <label className="text-sm text-gray-500">授權類型</label>
                        <p className="font-medium">{licenseInfo.license_type || 'Enterprise'}</p>
                      </div>
                      <div>
                        <label className="text-sm text-gray-500">授權狀態</label>
                        <Badge variant={licenseInfo.is_valid ? 'success' : 'destructive'}>
                          {licenseInfo.is_valid ? '有效' : '無效'}
                        </Badge>
                      </div>
                      <div>
                        <label className="text-sm text-gray-500">公司名稱</label>
                        <p className="font-medium">{licenseInfo.company_name || '-'}</p>
                      </div>
                      <div>
                        <label className="text-sm text-gray-500">授權使用者數</label>
                        <p className="font-medium">{licenseInfo.max_users || '無限制'}</p>
                      </div>
                      <div>
                        <label className="text-sm text-gray-500">授權開始日期</label>
                        <p className="font-medium">
                          {licenseInfo.start_date
                            ? format(new Date(licenseInfo.start_date), 'yyyy/MM/dd', { locale: zhTW })
                            : '-'}
                        </p>
                      </div>
                      <div>
                        <label className="text-sm text-gray-500">授權到期日期</label>
                        <p className="font-medium">
                          {licenseInfo.expiry_date
                            ? format(new Date(licenseInfo.expiry_date), 'yyyy/MM/dd', { locale: zhTW })
                            : '永久'}
                        </p>
                      </div>
                    </div>
                    <div className="pt-4 border-t">
                      <label className="text-sm text-gray-500">授權功能</label>
                      <div className="flex flex-wrap gap-2 mt-2">
                        {licenseInfo.features?.map((feature: string) => (
                          <Badge key={feature} variant="outline">
                            {feature}
                          </Badge>
                        ))}
                      </div>
                    </div>
                  </>
                ) : (
                  <div className="text-center py-8">
                    <Lock className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                    <p className="text-gray-500">無授權資訊</p>
                  </div>
                )}
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </div>
    </DashboardLayout>
  )
}