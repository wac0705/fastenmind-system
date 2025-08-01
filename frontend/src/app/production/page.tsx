'use client'

import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { useRouter } from 'next/navigation'
import DashboardLayout from '@/components/layout/DashboardLayout'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { 
  Factory,
  Package,
  Settings,
  Clock,
  CheckCircle,
  XCircle,
  AlertCircle,
  TrendingUp,
  TrendingDown,
  Plus,
  Eye,
  Play,
  Pause,
  Square,
  Users,
  Wrench,
  BarChart3,
  Calendar,
  Zap
} from 'lucide-react'
import productionService from '@/services/production.service'
import { format } from 'date-fns'
import { zhTW } from 'date-fns/locale'

export default function ProductionPage() {
  const router = useRouter()
  const [activeTab, setActiveTab] = useState('dashboard')

  // Fetch production dashboard
  const { data: dashboard, isLoading: isLoadingDashboard } = useQuery({
    queryKey: ['production-dashboard'],
    queryFn: () => productionService.getProductionDashboard(),
  })

  // Fetch production stats
  const { data: stats } = useQuery({
    queryKey: ['production-stats'],
    queryFn: () => productionService.getProductionStats(),
  })

  // Fetch recent production orders
  const { data: recentOrders } = useQuery({
    queryKey: ['recent-production-orders'],
    queryFn: () => productionService.listProductionOrders({ page: 1, page_size: 5 }),
  })

  // Fetch my tasks
  const { data: myTasks } = useQuery({
    queryKey: ['my-production-tasks'],
    queryFn: () => productionService.listProductionTasks({ 
      status: 'in_progress,assigned',
      page: 1, 
      page_size: 5 
    }),
  })

  // Fetch pending inspections
  const { data: pendingInspections } = useQuery({
    queryKey: ['pending-inspections'],
    queryFn: () => productionService.listQualityInspections({ 
      status: 'pending',
      page: 1, 
      page_size: 5 
    }),
  })

  const getOrderStatusBadge = (status: string) => {
    const statusConfig: Record<string, { label: string; variant: any; icon: any }> = {
      planned: { label: '計劃中', variant: 'secondary', icon: Calendar },
      released: { label: '已發布', variant: 'info', icon: Play },
      in_progress: { label: '生產中', variant: 'warning', icon: Factory },
      quality_check: { label: '品檢中', variant: 'warning', icon: CheckCircle },
      completed: { label: '已完成', variant: 'success', icon: CheckCircle },
      cancelled: { label: '已取消', variant: 'destructive', icon: XCircle },
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

  const getPriorityBadge = (priority: string) => {
    const priorityConfig: Record<string, { label: string; variant: any }> = {
      low: { label: '低', variant: 'secondary' },
      medium: { label: '中', variant: 'info' },
      high: { label: '高', variant: 'warning' },
      urgent: { label: '緊急', variant: 'destructive' },
    }

    const config = priorityConfig[priority] || { label: priority, variant: 'default' }
    return <Badge variant={config.variant as any}>{config.label}</Badge>
  }

  const getTaskStatusBadge = (status: string) => {
    const statusConfig: Record<string, { label: string; variant: any }> = {
      pending: { label: '待處理', variant: 'secondary' },
      assigned: { label: '已指派', variant: 'info' },
      in_progress: { label: '進行中', variant: 'warning' },
      completed: { label: '已完成', variant: 'success' },
      on_hold: { label: '暫停', variant: 'warning' },
      cancelled: { label: '已取消', variant: 'destructive' },
    }

    const config = statusConfig[status] || { label: status, variant: 'default' }
    return <Badge variant={config.variant as any}>{config.label}</Badge>
  }

  const getStationStatusBadge = (status: string) => {
    const statusConfig: Record<string, { label: string; variant: any; icon: any }> = {
      available: { label: '可用', variant: 'success', icon: CheckCircle },
      busy: { label: '忙碌', variant: 'warning', icon: Factory },
      maintenance: { label: '維護中', variant: 'info', icon: Wrench },
      breakdown: { label: '故障', variant: 'destructive', icon: XCircle },
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

  if (isLoadingDashboard) {
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
            <h1 className="text-3xl font-bold text-gray-900">生產管理</h1>
            <p className="mt-2 text-gray-600">管理生產訂單、工藝路線、工作站與品質檢驗</p>
          </div>
          <div className="flex gap-2">
            <Button variant="outline" onClick={() => router.push('/production/reports')}>
              <BarChart3 className="mr-2 h-4 w-4" />
              生產報表
            </Button>
            <Button onClick={() => router.push('/production/orders/new')}>
              <Plus className="mr-2 h-4 w-4" />
              新增生產單
            </Button>
          </div>
        </div>

        <Tabs value={activeTab} onValueChange={setActiveTab}>
          <TabsList>
            <TabsTrigger value="dashboard">生產概覽</TabsTrigger>
            <TabsTrigger value="orders">生產訂單</TabsTrigger>
            <TabsTrigger value="tasks">生產任務</TabsTrigger>
            <TabsTrigger value="stations">工作站</TabsTrigger>
            <TabsTrigger value="quality">品質管理</TabsTrigger>
          </TabsList>

          <TabsContent value="dashboard" className="space-y-6">
            {/* Production KPIs */}
            {dashboard && (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                <Card>
                  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">總生產單</CardTitle>
                    <Package className="h-4 w-4 text-muted-foreground" />
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold">{dashboard.total_orders}</div>
                    <p className="text-xs text-muted-foreground">
                      進行中: {dashboard.in_progress_orders} | 已完成: {dashboard.completed_orders}
                    </p>
                  </CardContent>
                </Card>
                <Card>
                  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">生產效率</CardTitle>
                    <TrendingUp className="h-4 w-4 text-muted-foreground" />
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold">
                      {dashboard.production_efficiency.toFixed(1)}%
                    </div>
                    <p className="text-xs text-muted-foreground">
                      實際產量 / 計劃產量
                    </p>
                  </CardContent>
                </Card>
                <Card>
                  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">品質合格率</CardTitle>
                    <CheckCircle className="h-4 w-4 text-muted-foreground" />
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold">
                      {dashboard.quality_rate.toFixed(1)}%
                    </div>
                    <p className="text-xs text-muted-foreground">
                      合格數量 / 生產數量
                    </p>
                  </CardContent>
                </Card>
                <Card>
                  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">不良品數量</CardTitle>
                    <XCircle className="h-4 w-4 text-muted-foreground" />
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold text-red-600">
                      {dashboard.total_defect_quantity.toLocaleString()}
                    </div>
                    <p className="text-xs text-muted-foreground">需要改善的產品</p>
                  </CardContent>
                </Card>
              </div>
            )}

            {/* Production & Station Status */}
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              {/* Production Status */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Factory className="h-5 w-5" />
                    生產狀態分佈
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  {dashboard && (
                    <div className="space-y-3">
                      <div className="flex justify-between items-center">
                        <span className="text-sm">計劃中</span>
                        <div className="flex items-center gap-2">
                          <div className="w-24 bg-gray-200 rounded-full h-2">
                            <div 
                              className="bg-gray-400 h-2 rounded-full" 
                              style={{ width: `${(dashboard.planned_orders / dashboard.total_orders) * 100}%` }}
                            ></div>
                          </div>
                          <span className="text-sm font-medium">{dashboard.planned_orders}</span>
                        </div>
                      </div>
                      <div className="flex justify-between items-center">
                        <span className="text-sm">已發布</span>
                        <div className="flex items-center gap-2">
                          <div className="w-24 bg-gray-200 rounded-full h-2">
                            <div 
                              className="bg-blue-500 h-2 rounded-full" 
                              style={{ width: `${(dashboard.released_orders / dashboard.total_orders) * 100}%` }}
                            ></div>
                          </div>
                          <span className="text-sm font-medium">{dashboard.released_orders}</span>
                        </div>
                      </div>
                      <div className="flex justify-between items-center">
                        <span className="text-sm">生產中</span>
                        <div className="flex items-center gap-2">
                          <div className="w-24 bg-gray-200 rounded-full h-2">
                            <div 
                              className="bg-yellow-500 h-2 rounded-full" 
                              style={{ width: `${(dashboard.in_progress_orders / dashboard.total_orders) * 100}%` }}
                            ></div>
                          </div>
                          <span className="text-sm font-medium">{dashboard.in_progress_orders}</span>
                        </div>
                      </div>
                      <div className="flex justify-between items-center">
                        <span className="text-sm">已完成</span>
                        <div className="flex items-center gap-2">
                          <div className="w-24 bg-gray-200 rounded-full h-2">
                            <div 
                              className="bg-green-500 h-2 rounded-full" 
                              style={{ width: `${(dashboard.completed_orders / dashboard.total_orders) * 100}%` }}
                            ></div>
                          </div>
                          <span className="text-sm font-medium">{dashboard.completed_orders}</span>
                        </div>
                      </div>
                    </div>
                  )}
                </CardContent>
              </Card>

              {/* Work Station Status */}
              {stats && (
                <Card>
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                      <Settings className="h-5 w-5" />
                      工作站狀態
                    </CardTitle>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    <div className="grid grid-cols-2 gap-4">
                      <div className="text-center">
                        <div className="text-2xl font-bold text-green-600">{stats.available_stations}</div>
                        <p className="text-sm text-gray-500">可用</p>
                      </div>
                      <div className="text-center">
                        <div className="text-2xl font-bold text-yellow-600">{stats.busy_stations}</div>
                        <p className="text-sm text-gray-500">忙碌</p>
                      </div>
                      <div className="text-center">
                        <div className="text-2xl font-bold text-blue-600">{stats.maintenance_stations}</div>
                        <p className="text-sm text-gray-500">維護中</p>
                      </div>
                      <div className="text-center">
                        <div className="text-2xl font-bold text-red-600">{stats.breakdown_stations}</div>
                        <p className="text-sm text-gray-500">故障</p>
                      </div>
                    </div>
                    <div className="pt-4 border-t">
                      <div className="flex justify-between text-sm">
                        <span>總工作站數</span>
                        <span className="font-medium">{stats.total_work_stations}</span>
                      </div>
                      <div className="flex justify-between text-sm mt-1">
                        <span>設備利用率</span>
                        <span className="font-medium">
                          {((stats.busy_stations / stats.total_work_stations) * 100).toFixed(1)}%
                        </span>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              )}
            </div>

            {/* Recent Activities */}
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              {/* Recent Production Orders */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center justify-between">
                    <span className="flex items-center gap-2">
                      <Package className="h-5 w-5" />
                      最近生產單
                    </span>
                    <Button variant="ghost" size="sm" onClick={() => router.push('/production/orders')}>
                      <Eye className="h-4 w-4" />
                    </Button>
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  {recentOrders?.data.length === 0 ? (
                    <p className="text-center text-gray-500 py-4">暫無生產單記錄</p>
                  ) : (
                    <div className="space-y-3">
                      {recentOrders?.data.map((order) => (
                        <div key={order.id} className="flex items-center justify-between p-3 border rounded-lg">
                          <div>
                            <p className="font-medium">{order.order_no}</p>
                            <p className="text-sm text-gray-500">{order.product_name}</p>
                            <div className="flex items-center gap-2 mt-1">
                              {getOrderStatusBadge(order.status)}
                              {getPriorityBadge(order.priority)}
                            </div>
                          </div>
                          <div className="text-right">
                            <p className="font-medium">
                              {order.planned_quantity} {order.unit}
                            </p>
                            <p className="text-sm text-gray-500">
                              {format(new Date(order.planned_end_date), 'MM/dd', { locale: zhTW })}
                            </p>
                          </div>
                        </div>
                      ))}
                    </div>
                  )}
                </CardContent>
              </Card>

              {/* My Tasks */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center justify-between">
                    <span className="flex items-center gap-2">
                      <Users className="h-5 w-5" />
                      我的任務
                    </span>
                    <Button variant="ghost" size="sm" onClick={() => router.push('/production/tasks')}>
                      <Eye className="h-4 w-4" />
                    </Button>
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  {myTasks?.length === 0 ? (
                    <p className="text-center text-gray-500 py-4">暫無指派任務</p>
                  ) : (
                    <div className="space-y-3">
                      {myTasks?.map((task) => (
                        <div key={task.id} className="flex items-center justify-between p-3 border rounded-lg">
                          <div>
                            <p className="font-medium">{task.name}</p>
                            <p className="text-sm text-gray-500">
                              {task.production_order?.product_name}
                            </p>
                            <div className="mt-1">
                              {getTaskStatusBadge(task.status)}
                            </div>
                          </div>
                          <div className="text-right">
                            <p className="font-medium">
                              {task.planned_quantity} {task.production_order?.unit}
                            </p>
                            <p className="text-sm text-gray-500">
                              {task.work_station?.name}
                            </p>
                          </div>
                        </div>
                      ))}
                    </div>
                  )}
                </CardContent>
              </Card>
            </div>

            {/* Pending Quality Inspections */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center justify-between">
                  <span className="flex items-center gap-2">
                    <CheckCircle className="h-5 w-5" />
                    待檢驗項目
                  </span>
                  <Button variant="ghost" size="sm" onClick={() => router.push('/production/quality')}>
                    <Eye className="h-4 w-4" />
                  </Button>
                </CardTitle>
                <CardDescription>需要進行品質檢驗的項目</CardDescription>
              </CardHeader>
              <CardContent>
                {pendingInspections?.data.length === 0 ? (
                  <p className="text-center text-gray-500 py-4">暫無待檢驗項目</p>
                ) : (
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead>檢驗編號</TableHead>
                        <TableHead>類型</TableHead>
                        <TableHead>產品</TableHead>
                        <TableHead>數量</TableHead>
                        <TableHead>檢驗員</TableHead>
                        <TableHead>狀態</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {pendingInspections?.data.map((inspection) => (
                        <TableRow key={inspection.id}>
                          <TableCell className="font-medium">{inspection.inspection_no}</TableCell>
                          <TableCell>
                            <Badge variant="info">
                              {inspection.type === 'incoming' && '來料檢驗'}
                              {inspection.type === 'in_process' && '製程檢驗'}
                              {inspection.type === 'final' && '最終檢驗'}
                              {inspection.type === 'customer_return' && '客戶退貨'}
                            </Badge>
                          </TableCell>
                          <TableCell>{inspection.production_order?.product_name || '-'}</TableCell>
                          <TableCell>
                            {inspection.inspected_quantity} {inspection.unit}
                          </TableCell>
                          <TableCell>{inspection.inspector?.full_name}</TableCell>
                          <TableCell>
                            <Badge variant="warning">待檢驗</Badge>
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                )}
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="orders">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <Package className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">生產訂單管理功能開發中</p>
                  <Button className="mt-4" onClick={() => router.push('/production/orders')}>
                    前往生產訂單管理
                  </Button>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="tasks">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <Users className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">生產任務管理功能開發中</p>
                  <Button className="mt-4" onClick={() => router.push('/production/tasks')}>
                    前往生產任務管理
                  </Button>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="stations">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <Settings className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">工作站管理功能開發中</p>
                  <Button className="mt-4" onClick={() => router.push('/production/stations')}>
                    前往工作站管理
                  </Button>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="quality">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <CheckCircle className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">品質管理功能開發中</p>
                  <Button className="mt-4" onClick={() => router.push('/production/quality')}>
                    前往品質管理
                  </Button>
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </div>
    </DashboardLayout>
  )
}