'use client'

import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { useRouter } from 'next/navigation'
import DashboardLayout from '@/components/layout/DashboardLayout'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Progress } from '@/components/ui/progress'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import {
  BarChart3,
  LineChart,
  PieChart,
  TrendingUp,
  TrendingDown,
  ArrowUpRight,
  ArrowDownRight,
  DollarSign,
  Package,
  Users,
  ShoppingCart,
  Target,
  Activity,
  Calendar,
  Clock,
  RefreshCw,
  Settings,
  Plus,
  Download,
  Filter,
  Layout,
  Grid3x3,
  List,
  Eye,
  Edit,
  Star,
  StarOff,
  AlertCircle,
  CheckCircle,
  XCircle
} from 'lucide-react'
import reportService from '@/services/report.service'
import { format } from 'date-fns'
import { zhTW } from 'date-fns/locale'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  ArcElement,
  Title,
  Tooltip,
  Legend,
  ChartOptions
} from 'chart.js'
import { Line, Bar, Pie, Doughnut } from 'react-chartjs-2'

// Register ChartJS components
ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  ArcElement,
  Title,
  Tooltip,
  Legend
)

export default function ReportDashboardsPage() {
  const router = useRouter()
  const [selectedPeriod, setSelectedPeriod] = useState('month')
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid')
  const [activeTab, setActiveTab] = useState('overview')

  // Fetch dashboard data
  const { data: dashboardData, isLoading, refetch } = useQuery({
    queryKey: ['kpi-dashboard', selectedPeriod],
    queryFn: () => reportService.getKPIDashboard(selectedPeriod),
  })

  // Fetch business KPIs
  const { data: kpisData } = useQuery({
    queryKey: ['business-kpis'],
    queryFn: () => reportService.getBusinessKPIs(),
  })

  const getKPITrend = (value: number) => {
    if (value > 0) {
      return {
        icon: TrendingUp,
        color: 'text-green-600',
        bgColor: 'bg-green-50',
        arrow: ArrowUpRight,
      }
    } else if (value < 0) {
      return {
        icon: TrendingDown,
        color: 'text-red-600',
        bgColor: 'bg-red-50',
        arrow: ArrowDownRight,
      }
    }
    return {
      icon: Activity,
      color: 'text-gray-600',
      bgColor: 'bg-gray-50',
      arrow: null,
    }
  }

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'achieved':
        return { icon: CheckCircle, color: 'text-green-600' }
      case 'at_risk':
        return { icon: AlertCircle, color: 'text-yellow-600' }
      case 'missed':
        return { icon: XCircle, color: 'text-red-600' }
      default:
        return { icon: Activity, color: 'text-gray-600' }
    }
  }

  // Chart options
  const chartOptions: ChartOptions<any> = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        position: 'bottom' as const,
      },
    },
  }

  // Sample chart data
  const salesChartData = {
    labels: ['一月', '二月', '三月', '四月', '五月', '六月'],
    datasets: [
      {
        label: '銷售額',
        data: [65, 59, 80, 81, 56, 55],
        backgroundColor: 'rgba(59, 130, 246, 0.5)',
        borderColor: 'rgb(59, 130, 246)',
        tension: 0.4,
      },
      {
        label: '目標',
        data: [70, 70, 70, 70, 70, 70],
        backgroundColor: 'rgba(156, 163, 175, 0.5)',
        borderColor: 'rgb(156, 163, 175)',
        borderDash: [5, 5],
      },
    ],
  }

  const categoryChartData = {
    labels: ['緊固件', '五金工具', '建材', '電子零件', '其他'],
    datasets: [
      {
        data: [30, 25, 20, 15, 10],
        backgroundColor: [
          'rgba(59, 130, 246, 0.8)',
          'rgba(34, 197, 94, 0.8)',
          'rgba(251, 146, 60, 0.8)',
          'rgba(168, 85, 247, 0.8)',
          'rgba(156, 163, 175, 0.8)',
        ],
      },
    ],
  }

  if (isLoading) {
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
            <h1 className="text-3xl font-bold text-gray-900">KPI 監控儀表板</h1>
            <p className="mt-2 text-gray-600">即時監控關鍵績效指標，掌握業務狀況</p>
          </div>
          <div className="flex items-center gap-2">
            <Select value={selectedPeriod} onValueChange={setSelectedPeriod}>
              <SelectTrigger className="w-32">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="day">今日</SelectItem>
                <SelectItem value="week">本週</SelectItem>
                <SelectItem value="month">本月</SelectItem>
                <SelectItem value="quarter">本季</SelectItem>
                <SelectItem value="year">本年</SelectItem>
              </SelectContent>
            </Select>
            <Button variant="outline" onClick={() => refetch()}>
              <RefreshCw className="mr-2 h-4 w-4" />
              重新整理
            </Button>
            <Button variant="outline" onClick={() => setViewMode(viewMode === 'grid' ? 'list' : 'grid')}>
              {viewMode === 'grid' ? <List className="h-4 w-4" /> : <Grid3x3 className="h-4 w-4" />}
            </Button>
            <Button variant="outline">
              <Download className="mr-2 h-4 w-4" />
              匯出
            </Button>
            <Button onClick={() => router.push('/reports/dashboards/new')}>
              <Plus className="mr-2 h-4 w-4" />
              新增儀表板
            </Button>
          </div>
        </div>

        <Tabs value={activeTab} onValueChange={setActiveTab}>
          <TabsList>
            <TabsTrigger value="overview">總覽</TabsTrigger>
            <TabsTrigger value="sales">銷售分析</TabsTrigger>
            <TabsTrigger value="operations">營運指標</TabsTrigger>
            <TabsTrigger value="finance">財務概況</TabsTrigger>
            <TabsTrigger value="custom">自訂儀表板</TabsTrigger>
          </TabsList>

          <TabsContent value="overview" className="space-y-6">
            {/* Key Metrics */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">總營收</CardTitle>
                  <DollarSign className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">
                    {dashboardData?.revenue?.currency} {dashboardData?.revenue?.amount.toLocaleString()}
                  </div>
                  <div className="flex items-center mt-2">
                    {dashboardData?.revenue?.change !== 0 && (
                      <>
                        {getKPITrend(dashboardData?.revenue?.change).arrow && (
                          <ArrowUpRight className={`h-4 w-4 ${getKPITrend(dashboardData?.revenue?.change).color}`} />
                        )}
                        <span className={`text-sm ${getKPITrend(dashboardData?.revenue?.change).color}`}>
                          {Math.abs(dashboardData?.revenue?.change)}%
                        </span>
                      </>
                    )}
                    <span className="text-sm text-muted-foreground ml-2">較上期</span>
                  </div>
                  <Progress value={dashboardData?.revenue?.progress} className="mt-2" />
                </CardContent>
              </Card>

              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">訂單數量</CardTitle>
                  <ShoppingCart className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">{dashboardData?.orders?.count.toLocaleString()}</div>
                  <div className="flex items-center mt-2">
                    {dashboardData?.orders?.change !== 0 && (
                      <>
                        {getKPITrend(dashboardData?.orders?.change).arrow && (
                          <ArrowUpRight className={`h-4 w-4 ${getKPITrend(dashboardData?.orders?.change).color}`} />
                        )}
                        <span className={`text-sm ${getKPITrend(dashboardData?.orders?.change).color}`}>
                          {Math.abs(dashboardData?.orders?.change)}%
                        </span>
                      </>
                    )}
                    <span className="text-sm text-muted-foreground ml-2">較上期</span>
                  </div>
                  <Progress value={dashboardData?.orders?.progress} className="mt-2" />
                </CardContent>
              </Card>

              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">新客戶</CardTitle>
                  <Users className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">{dashboardData?.newCustomers?.count}</div>
                  <div className="flex items-center mt-2">
                    {dashboardData?.newCustomers?.change !== 0 && (
                      <>
                        {getKPITrend(dashboardData?.newCustomers?.change).arrow && (
                          <ArrowUpRight className={`h-4 w-4 ${getKPITrend(dashboardData?.newCustomers?.change).color}`} />
                        )}
                        <span className={`text-sm ${getKPITrend(dashboardData?.newCustomers?.change).color}`}>
                          {Math.abs(dashboardData?.newCustomers?.change)}%
                        </span>
                      </>
                    )}
                    <span className="text-sm text-muted-foreground ml-2">較上期</span>
                  </div>
                  <Progress value={dashboardData?.newCustomers?.progress} className="mt-2" />
                </CardContent>
              </Card>

              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">庫存周轉率</CardTitle>
                  <Package className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">{dashboardData?.inventoryTurnover?.rate}x</div>
                  <div className="flex items-center mt-2">
                    {dashboardData?.inventoryTurnover?.change !== 0 && (
                      <>
                        {getKPITrend(dashboardData?.inventoryTurnover?.change).arrow && (
                          <ArrowUpRight className={`h-4 w-4 ${getKPITrend(dashboardData?.inventoryTurnover?.change).color}`} />
                        )}
                        <span className={`text-sm ${getKPITrend(dashboardData?.inventoryTurnover?.change).color}`}>
                          {Math.abs(dashboardData?.inventoryTurnover?.change)}%
                        </span>
                      </>
                    )}
                    <span className="text-sm text-muted-foreground ml-2">較上期</span>
                  </div>
                  <Progress value={dashboardData?.inventoryTurnover?.progress} className="mt-2" />
                </CardContent>
              </Card>
            </div>

            {/* Charts Row */}
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <LineChart className="h-5 w-5" />
                    銷售趨勢
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="h-80">
                    <Line data={salesChartData} options={chartOptions} />
                  </div>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <PieChart className="h-5 w-5" />
                    產品類別分布
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="h-80">
                    <Doughnut data={categoryChartData} options={chartOptions} />
                  </div>
                </CardContent>
              </Card>
            </div>

            {/* KPI Goals */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center justify-between">
                  <span className="flex items-center gap-2">
                    <Target className="h-5 w-5" />
                    KPI 目標達成情況
                  </span>
                  <Button variant="outline" size="sm" onClick={() => router.push('/reports/kpis')}>
                    管理 KPI
                  </Button>
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {kpisData?.map((kpi: any) => {
                    const StatusIcon = getStatusIcon(kpi.status).icon
                    const statusColor = getStatusIcon(kpi.status).color
                    
                    return (
                      <div key={kpi.id} className="flex items-center justify-between p-4 border rounded-lg">
                        <div className="flex-1">
                          <div className="flex items-center gap-2">
                            <StatusIcon className={`h-4 w-4 ${statusColor}`} />
                            <h4 className="font-medium">{kpi.name}</h4>
                            <Badge variant="outline">{kpi.category}</Badge>
                          </div>
                          <p className="text-sm text-gray-500 mt-1">{kpi.description}</p>
                          <div className="flex items-center gap-4 mt-2">
                            <span className="text-sm">
                              目標: <strong>{kpi.target_value} {kpi.unit}</strong>
                            </span>
                            <span className="text-sm">
                              實際: <strong>{kpi.actual_value} {kpi.unit}</strong>
                            </span>
                            <span className="text-sm">
                              達成率: <strong>{kpi.achievement_rate}%</strong>
                            </span>
                          </div>
                        </div>
                        <div className="w-32">
                          <Progress value={kpi.achievement_rate} className="h-2" />
                        </div>
                      </div>
                    )
                  })}
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="sales" className="space-y-6">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <BarChart3 className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">銷售分析儀表板開發中</p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="operations" className="space-y-6">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <Activity className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">營運指標儀表板開發中</p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="finance" className="space-y-6">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <DollarSign className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">財務概況儀表板開發中</p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="custom" className="space-y-6">
            {/* Custom Dashboards Grid */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {dashboardData?.customDashboards?.map((dashboard: any) => (
                <Card key={dashboard.id} className="hover:shadow-lg transition-shadow cursor-pointer">
                  <CardHeader>
                    <div className="flex justify-between items-start">
                      <div>
                        <CardTitle className="text-lg">{dashboard.name}</CardTitle>
                        <CardDescription>{dashboard.description}</CardDescription>
                      </div>
                      <div className="flex items-center gap-1">
                        {dashboard.is_favorite ? (
                          <Star className="h-4 w-4 text-yellow-500 fill-yellow-500" />
                        ) : (
                          <StarOff className="h-4 w-4 text-gray-400" />
                        )}
                      </div>
                    </div>
                  </CardHeader>
                  <CardContent>
                    <div className="grid grid-cols-2 gap-4 mb-4">
                      <div className="text-center p-3 bg-gray-50 rounded">
                        <Layout className="h-6 w-6 mx-auto mb-1 text-gray-400" />
                        <p className="text-xs text-gray-500">元件數</p>
                        <p className="font-semibold">{dashboard.widget_count}</p>
                      </div>
                      <div className="text-center p-3 bg-gray-50 rounded">
                        <Clock className="h-6 w-6 mx-auto mb-1 text-gray-400" />
                        <p className="text-xs text-gray-500">更新頻率</p>
                        <p className="font-semibold">{dashboard.refresh_interval}</p>
                      </div>
                    </div>
                    <div className="flex gap-2">
                      <Button
                        size="sm"
                        className="flex-1"
                        onClick={() => router.push(`/reports/dashboards/${dashboard.id}`)}
                      >
                        <Eye className="mr-2 h-4 w-4" />
                        檢視
                      </Button>
                      <Button
                        size="sm"
                        variant="outline"
                        onClick={() => router.push(`/reports/dashboards/${dashboard.id}/edit`)}
                      >
                        <Edit className="h-4 w-4" />
                      </Button>
                    </div>
                  </CardContent>
                </Card>
              ))}
              
              {/* Add New Dashboard Card */}
              <Card className="border-dashed hover:shadow-lg transition-shadow cursor-pointer">
                <CardContent className="flex flex-col items-center justify-center h-full py-12">
                  <Plus className="h-12 w-12 text-gray-300 mb-4" />
                  <p className="text-gray-500 mb-4">建立新儀表板</p>
                  <Button onClick={() => router.push('/reports/dashboards/new')}>
                    開始建立
                  </Button>
                </CardContent>
              </Card>
            </div>
          </TabsContent>
        </Tabs>
      </div>
    </DashboardLayout>
  )
}