'use client'

import { useState, useEffect } from 'react'
import { useQuery } from '@tanstack/react-query'
import DashboardLayout from '@/components/layout/DashboardLayout'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
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
  Smartphone,
  Tablet,
  Monitor,
  Bell,
  Activity,
  Download,
  Wifi,
  WifiOff,
  Battery,
  Cpu,
  HardDrive,
  Users,
  TrendingUp,
  TrendingDown,
  PlayCircle,
  PauseCircle,
  CheckCircle,
  XCircle,
  AlertCircle,
  Clock,
  RefreshCw,
  Search,
  Filter,
  Settings,
  Eye,
  Send,
  BarChart3,
  Globe,
  Shield,
  Zap
} from 'lucide-react'
import mobileService from '@/services/mobile.service'
import { format } from 'date-fns'
import { zhTW } from 'date-fns/locale'

export default function MobilePage() {
  const [activeTab, setActiveTab] = useState('overview')
  const [searchQuery, setSearchQuery] = useState('')
  const [platformFilter, setPlatformFilter] = useState('')
  const [statusFilter, setStatusFilter] = useState('')
  const [page, setPage] = useState(1)
  const pageSize = 20

  // PWA 相關狀態
  const [isOnline, setIsOnline] = useState(true)
  const [installPrompt, setInstallPrompt] = useState<any>(null)
  const [isInstalled, setIsInstalled] = useState(false)
  const [notificationPermission, setNotificationPermission] = useState<NotificationPermission>('default')

  // Fetch mobile statistics
  const { data: mobileStats, isLoading: isLoadingStats } = useQuery({
    queryKey: ['mobile-statistics'],
    queryFn: () => mobileService.getMobileStatistics(),
  })

  // Fetch company devices
  const { data: devicesData, isLoading: isLoadingDevices } = useQuery({
    queryKey: ['company-devices', page, searchQuery, platformFilter, statusFilter],
    queryFn: () => mobileService.listCompanyDevices({
      page,
      page_size: pageSize,
      search: searchQuery || undefined,
      platform: platformFilter || undefined,
      is_active: statusFilter === 'active' ? true : statusFilter === 'inactive' ? false : undefined,
    }),
  })

  // Fetch push notifications
  const { data: notificationsData } = useQuery({
    queryKey: ['push-notifications'],
    queryFn: () => mobileService.listPushNotifications({ page: 1, page_size: 10 }),
  })

  // Fetch notification stats
  const { data: notificationStats } = useQuery({
    queryKey: ['notification-stats'],
    queryFn: () => mobileService.getNotificationStats(30),
  })

  // PWA 初始化
  useEffect(() => {
    // 檢查網路狀態
    const updateOnlineStatus = () => setIsOnline(typeof navigator !== 'undefined' ? navigator.onLine : true)
    window.addEventListener('online', updateOnlineStatus)
    window.addEventListener('offline', updateOnlineStatus)

    // 檢查是否已安裝為 PWA
    const checkInstallStatus = () => {
      setIsInstalled(window.matchMedia('(display-mode: standalone)').matches)
    }
    checkInstallStatus()

    // 檢查通知權限
    if ('Notification' in window) {
      setNotificationPermission(Notification.permission)
    }

    // 註冊 Service Worker
    mobileService.registerServiceWorker()

    // 監聽安裝提示
    const handleInstallPrompt = (e: any) => {
      e.preventDefault()
      setInstallPrompt(e)
    }
    window.addEventListener('beforeinstallprompt', handleInstallPrompt)

    return () => {
      window.removeEventListener('online', updateOnlineStatus)
      window.removeEventListener('offline', updateOnlineStatus)
      window.removeEventListener('beforeinstallprompt', handleInstallPrompt)
    }
  }, [])

  const getPlatformIcon = (platform: string) => {
    switch (platform.toLowerCase()) {
      case 'ios':
        return <Smartphone className="h-4 w-4" />
      case 'android':
        return <Smartphone className="h-4 w-4" />
      case 'web':
        return <Monitor className="h-4 w-4" />
      default:
        return <Monitor className="h-4 w-4" />
    }
  }

  const getDeviceTypeIcon = (deviceType: string) => {
    switch (deviceType.toLowerCase()) {
      case 'phone':
        return <Smartphone className="h-4 w-4" />
      case 'tablet':
        return <Tablet className="h-4 w-4" />
      case 'desktop':
        return <Monitor className="h-4 w-4" />
      default:
        return <Monitor className="h-4 w-4" />
    }
  }

  const getStatusBadge = (isActive: boolean, lastSeen: string) => {
    const lastSeenDate = new Date(lastSeen)
    const now = new Date()
    const diffMinutes = (now.getTime() - lastSeenDate.getTime()) / (1000 * 60)

    if (!isActive) {
      return <Badge variant="secondary">已停用</Badge>
    } else if (diffMinutes < 5) {
      return <Badge variant="success">線上</Badge>
    } else if (diffMinutes < 60) {
      return <Badge variant="warning">最近活躍</Badge>
    } else {
      return <Badge variant="secondary">離線</Badge>
    }
  }

  const getNotificationStatusBadge = (status: string) => {
    const statusConfig: Record<string, { label: string; variant: any; icon: any }> = {
      pending: { label: '等待中', variant: 'secondary', icon: Clock },
      sent: { label: '已發送', variant: 'info', icon: Send },
      delivered: { label: '已送達', variant: 'success', icon: CheckCircle },
      clicked: { label: '已點擊', variant: 'success', icon: Eye },
      failed: { label: '發送失敗', variant: 'destructive', icon: XCircle },
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

  const handleInstallApp = async () => {
    if (installPrompt) {
      installPrompt.prompt()
      const choiceResult = await installPrompt.userChoice
      if (choiceResult.outcome === 'accepted') {
        setInstallPrompt(null)
        setIsInstalled(true)
      }
    }
  }

  const handleRequestNotificationPermission = async () => {
    const permission = await mobileService.requestNotificationPermission()
    setNotificationPermission(permission)
  }

  const handleSendTestNotification = async () => {
    try {
      await mobileService.sendNotificationToUsers({
        user_ids: ['current-user'], // 這裡應該使用實際的用戶ID
        title: 'FastenMind 測試通知',
        body: '這是一個測試推播通知',
        type: 'test'
      })
    } catch (error) {
      console.error('Failed to send test notification:', error)
    }
  }

  if (isLoadingStats) {
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
            <h1 className="text-3xl font-bold text-gray-900">行動應用程式</h1>
            <p className="mt-2 text-gray-600">行動裝置管理、推播通知與 PWA 功能</p>
          </div>
          <div className="flex items-center gap-4">
            {/* 網路狀態指示器 */}
            <div className="flex items-center gap-2">
              {isOnline ? (
                <>
                  <Wifi className="h-4 w-4 text-green-500" />
                  <span className="text-sm text-green-600">線上</span>
                </>
              ) : (
                <>
                  <WifiOff className="h-4 w-4 text-red-500" />
                  <span className="text-sm text-red-600">離線</span>
                </>
              )}
            </div>
            
            {/* PWA 安裝按鈕 */}
            {!isInstalled && installPrompt && (
              <Button variant="outline" onClick={handleInstallApp}>
                <Download className="mr-2 h-4 w-4" />
                安裝應用程式
              </Button>
            )}
            
            {/* 通知權限按鈕 */}
            {notificationPermission === 'default' && (
              <Button variant="outline" onClick={handleRequestNotificationPermission}>
                <Bell className="mr-2 h-4 w-4" />
                啟用通知
              </Button>
            )}
            
            {/* 測試通知按鈕 */}
            {notificationPermission === 'granted' && (
              <Button variant="outline" onClick={handleSendTestNotification}>
                <Send className="mr-2 h-4 w-4" />
                測試通知
              </Button>
            )}
          </div>
        </div>

        <Tabs value={activeTab} onValueChange={setActiveTab}>
          <TabsList>
            <TabsTrigger value="overview">概覽</TabsTrigger>
            <TabsTrigger value="devices">裝置管理</TabsTrigger>
            <TabsTrigger value="notifications">推播通知</TabsTrigger>
            <TabsTrigger value="analytics">使用分析</TabsTrigger>
            <TabsTrigger value="config">設定管理</TabsTrigger>
          </TabsList>

          <TabsContent value="overview" className="space-y-6">
            {/* Mobile Statistics */}
            {mobileStats && (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                <Card>
                  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">總裝置數</CardTitle>
                    <Smartphone className="h-4 w-4 text-muted-foreground" />
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold">{mobileStats.total_devices}</div>
                    <p className="text-xs text-muted-foreground">
                      活躍裝置: {mobileStats.active_devices}
                    </p>
                  </CardContent>
                </Card>
                <Card>
                  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">推播通知</CardTitle>
                    <Bell className="h-4 w-4 text-muted-foreground" />
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold">{mobileStats.total_notifications}</div>
                    <p className="text-xs text-muted-foreground">
                      本月通知數量
                    </p>
                  </CardContent>
                </Card>
                <Card>
                  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">PWA 狀態</CardTitle>
                    <Globe className="h-4 w-4 text-muted-foreground" />
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold">
                      {isInstalled ? '已安裝' : '網頁版'}
                    </div>
                    <p className="text-xs text-muted-foreground">
                      {isOnline ? '線上模式' : '離線模式'}
                    </p>
                  </CardContent>
                </Card>
                <Card>
                  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">通知權限</CardTitle>
                    <Shield className="h-4 w-4 text-muted-foreground" />
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold">
                      {notificationPermission === 'granted' ? '已授權' : 
                       notificationPermission === 'denied' ? '已拒絕' : '未設定'}
                    </div>
                    <p className="text-xs text-muted-foreground">
                      推播通知狀態
                    </p>
                  </CardContent>
                </Card>
              </div>
            )}

            {/* Platform Distribution */}
            {mobileStats && (
              <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                <Card>
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                      <BarChart3 className="h-5 w-5" />
                      平台分布
                    </CardTitle>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    {mobileStats.devices_by_platform.map((item) => (
                      <div key={item.platform} className="flex justify-between items-center">
                        <div className="flex items-center gap-2">
                          {getPlatformIcon(item.platform)}
                          <span className="text-sm capitalize">{item.platform}</span>
                        </div>
                        <div className="flex items-center gap-2">
                          <div className="w-24 bg-gray-200 rounded-full h-2">
                            <div 
                              className="bg-blue-500 h-2 rounded-full" 
                              style={{ width: `${(item.count / mobileStats.total_devices) * 100}%` }}
                            ></div>
                          </div>
                          <span className="text-sm font-medium">{item.count}</span>
                        </div>
                      </div>
                    ))}
                  </CardContent>
                </Card>

                <Card>
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                      <Activity className="h-5 w-5" />
                      通知狀態統計
                    </CardTitle>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    {mobileStats.notifications_by_status.map((item) => (
                      <div key={item.status} className="flex justify-between items-center">
                        <span className="text-sm">{getNotificationStatusBadge(item.status)}</span>
                        <div className="flex items-center gap-2">
                          <div className="w-24 bg-gray-200 rounded-full h-2">
                            <div 
                              className="bg-green-500 h-2 rounded-full" 
                              style={{ width: `${(item.count / mobileStats.total_notifications) * 100}%` }}
                            ></div>
                          </div>
                          <span className="text-sm font-medium">{item.count}</span>
                        </div>
                      </div>
                    ))}
                  </CardContent>
                </Card>
              </div>
            )}
          </TabsContent>

          <TabsContent value="devices" className="space-y-6">
            {/* Filters */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Filter className="h-5 w-5" />
                  篩選條件
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                  <div className="space-y-2">
                    <label className="text-sm font-medium">搜尋</label>
                    <div className="relative">
                      <Search className="absolute left-3 top-3 h-4 w-4 text-gray-400" />
                      <Input
                        placeholder="裝置名稱或型號"
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                        className="pl-10"
                      />
                    </div>
                  </div>
                  <div className="space-y-2">
                    <label className="text-sm font-medium">平台</label>
                    <Select value={platformFilter} onValueChange={setPlatformFilter}>
                      <SelectTrigger>
                        <SelectValue placeholder="選擇平台" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="">全部</SelectItem>
                        <SelectItem value="ios">iOS</SelectItem>
                        <SelectItem value="android">Android</SelectItem>
                        <SelectItem value="web">Web</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                  <div className="space-y-2">
                    <label className="text-sm font-medium">狀態</label>
                    <Select value={statusFilter} onValueChange={setStatusFilter}>
                      <SelectTrigger>
                        <SelectValue placeholder="選擇狀態" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="">全部</SelectItem>
                        <SelectItem value="active">啟用</SelectItem>
                        <SelectItem value="inactive">停用</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Devices Table */}
            <Card>
              <CardHeader>
                <CardTitle>裝置列表</CardTitle>
                <CardDescription>
                  共 {devicesData?.total || 0} 個裝置
                </CardDescription>
              </CardHeader>
              <CardContent>
                {isLoadingDevices ? (
                  <div className="text-center py-8">載入中...</div>
                ) : devicesData?.data.length === 0 ? (
                  <div className="text-center py-8">
                    <Smartphone className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                    <p className="text-gray-500">暫無裝置資料</p>
                  </div>
                ) : (
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead>裝置資訊</TableHead>
                        <TableHead>平台</TableHead>
                        <TableHead>類型</TableHead>
                        <TableHead>狀態</TableHead>
                        <TableHead>最後上線</TableHead>
                        <TableHead>用戶</TableHead>
                        <TableHead>操作</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {devicesData?.data.map((device) => (
                        <TableRow key={device.id}>
                          <TableCell>
                            <div>
                              <p className="font-medium">{device.device_name}</p>
                              <p className="text-sm text-gray-500">{device.device_model}</p>
                              <p className="text-xs text-gray-400">{device.os_version}</p>
                            </div>
                          </TableCell>
                          <TableCell>
                            <div className="flex items-center gap-2">
                              {getPlatformIcon(device.platform)}
                              <span className="capitalize">{device.platform}</span>
                            </div>
                          </TableCell>
                          <TableCell>
                            <div className="flex items-center gap-2">
                              {getDeviceTypeIcon(device.device_type)}
                              <span className="capitalize">{device.device_type}</span>
                            </div>
                          </TableCell>
                          <TableCell>
                            {getStatusBadge(device.is_active, device.last_seen)}
                          </TableCell>
                          <TableCell>
                            <div className="text-sm">
                              <p>{format(new Date(device.last_seen), 'MM/dd HH:mm', { locale: zhTW })}</p>
                              <p className="text-gray-500">
                                {device.time_zone}
                              </p>
                            </div>
                          </TableCell>
                          <TableCell>
                            <div className="text-sm">
                              <p className="font-medium">{device.user?.full_name}</p>
                              <p className="text-gray-500">{device.user?.email}</p>
                            </div>
                          </TableCell>
                          <TableCell>
                            <div className="flex items-center gap-2">
                              <Button variant="ghost" size="sm" title="查看詳情">
                                <Eye className="h-4 w-4" />
                              </Button>
                              <Button variant="ghost" size="sm" title="設定">
                                <Settings className="h-4 w-4" />
                              </Button>
                            </div>
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                )}
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="notifications">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <Bell className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">推播通知管理功能開發中</p>
                  <p className="text-sm text-gray-400 mt-2">包含通知發送、統計分析與模板管理</p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="analytics">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <BarChart3 className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">使用分析功能開發中</p>
                  <p className="text-sm text-gray-400 mt-2">包含用戶行為分析、效能監控與使用統計</p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="config">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <Settings className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">設定管理功能開發中</p>
                  <p className="text-sm text-gray-400 mt-2">包含應用程式設定、功能開關與個人化配置</p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </div>
    </DashboardLayout>
  )
}