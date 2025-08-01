'use client'

import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { useRouter } from 'next/navigation'
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
  BarChart3,
  FileText,
  Plus,
  Search,
  Filter,
  Download,
  Upload,
  Eye,
  Edit,
  Trash2,
  Play,
  Pause,
  Square,
  Copy,
  Calendar,
  Clock,
  Users,
  TrendingUp,
  TrendingDown,
  Activity,
  Database,
  Settings,
  Zap,
  CheckCircle,
  XCircle,
  AlertCircle,
  RefreshCw
} from 'lucide-react'
import reportService from '@/services/report.service'
import { format } from 'date-fns'
import { zhTW } from 'date-fns/locale'

export default function ReportsPage() {
  const router = useRouter()
  const [activeTab, setActiveTab] = useState('dashboard')
  const [searchQuery, setSearchQuery] = useState('')
  const [categoryFilter, setCategoryFilter] = useState('')
  const [typeFilter, setTypeFilter] = useState('')
  const [statusFilter, setStatusFilter] = useState('')
  const [page, setPage] = useState(1)
  const pageSize = 20

  // Fetch report dashboard data
  const { data: dashboardData, isLoading: isLoadingDashboard } = useQuery({
    queryKey: ['report-dashboard'],
    queryFn: () => reportService.getReportDashboardData(),
  })

  // Fetch reports list
  const { data: reportsData, isLoading: isLoadingReports, refetch: refetchReports } = useQuery({
    queryKey: ['reports', page, searchQuery, categoryFilter, typeFilter, statusFilter],
    queryFn: () => reportService.listReports({
      page,
      page_size: pageSize,
      search: searchQuery || undefined,
      category: categoryFilter || undefined,
      type: typeFilter || undefined,
      status: statusFilter || undefined,
    }),
  })

  // Fetch report templates
  const { data: templatesData } = useQuery({
    queryKey: ['report-templates'],
    queryFn: () => reportService.listReportTemplates({ page: 1, page_size: 10 }),
  })

  // Fetch recent executions
  const { data: recentExecutions } = useQuery({
    queryKey: ['recent-executions'],
    queryFn: () => reportService.getRecentExecutions(10),
  })

  const getStatusBadge = (status: string) => {
    const statusConfig: Record<string, { label: string; variant: any; icon: any }> = {
      active: { label: '啟用', variant: 'success', icon: CheckCircle },
      inactive: { label: '停用', variant: 'secondary', icon: Pause },
      archived: { label: '已封存', variant: 'warning', icon: Clock },
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

  const getCategoryBadge = (category: string) => {
    const categoryConfig: Record<string, { label: string; variant: any }> = {
      sales: { label: '銷售', variant: 'info' },
      finance: { label: '財務', variant: 'success' },
      production: { label: '生產', variant: 'warning' },
      inventory: { label: '庫存', variant: 'secondary' },
      supplier: { label: '供應商', variant: 'info' },
      customer: { label: '客戶', variant: 'success' },
      system: { label: '系統', variant: 'secondary' },
    }

    const config = categoryConfig[category] || { label: category, variant: 'default' }
    return <Badge variant={config.variant as any}>{config.label}</Badge>
  }

  const getTypeBadge = (type: string) => {
    const typeConfig: Record<string, { label: string; variant: any }> = {
      summary: { label: '摘要', variant: 'info' },
      detail: { label: '詳細', variant: 'secondary' },
      trend: { label: '趨勢', variant: 'warning' },
      comparison: { label: '比較', variant: 'success' },
      dashboard: { label: '儀表板', variant: 'info' },
    }

    const config = typeConfig[type] || { label: type, variant: 'default' }
    return <Badge variant={config.variant as any}>{config.label}</Badge>
  }

  const getExecutionStatusBadge = (status: string) => {
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

  const handleSearch = () => {
    setPage(1)
    refetchReports()
  }

  const handleClearFilters = () => {
    setSearchQuery('')
    setCategoryFilter('')
    setTypeFilter('')
    setStatusFilter('')
    setPage(1)
    refetchReports()
  }

  const handleExecuteReport = async (reportId: string) => {
    try {
      await reportService.executeReport(reportId)
      // Show success message and refresh data
      refetchReports()
    } catch (error) {
      console.error('Failed to execute report:', error)
    }
  }

  const handleDuplicateReport = async (reportId: string) => {
    try {
      await reportService.duplicateReport(reportId)
      refetchReports()
    } catch (error) {
      console.error('Failed to duplicate report:', error)
    }
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
            <h1 className="text-3xl font-bold text-gray-900">報表中心</h1>
            <p className="mt-2 text-gray-600">管理報表、範本與儀表板，提供全面的業務分析</p>
          </div>
          <div className="flex gap-2">
            <Button variant="outline" onClick={() => router.push('/reports/templates')}>
              <FileText className="mr-2 h-4 w-4" />
              範本庫
            </Button>
            <Button variant="outline" onClick={() => router.push('/reports/dashboards')}>
              <BarChart3 className="mr-2 h-4 w-4" />
              儀表板
            </Button>
            <Button onClick={() => router.push('/reports/new')}>
              <Plus className="mr-2 h-4 w-4" />
              新增報表
            </Button>
          </div>
        </div>

        <Tabs value={activeTab} onValueChange={setActiveTab}>
          <TabsList>
            <TabsTrigger value="dashboard">報表概覽</TabsTrigger>
            <TabsTrigger value="reports">我的報表</TabsTrigger>
            <TabsTrigger value="templates">範本庫</TabsTrigger>
            <TabsTrigger value="executions">執行記錄</TabsTrigger>
            <TabsTrigger value="settings">系統設定</TabsTrigger>
          </TabsList>

          <TabsContent value="dashboard" className="space-y-6">
            {/* Dashboard KPIs */}
            {dashboardData && (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                <Card>
                  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">總報表數</CardTitle>
                    <FileText className="h-4 w-4 text-muted-foreground" />
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold">{dashboardData.total_reports}</div>
                    <p className="text-xs text-muted-foreground">
                      範本: {dashboardData.total_templates} | 儀表板: {dashboardData.total_dashboards}
                    </p>
                  </CardContent>
                </Card>
                <Card>
                  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">執行次數</CardTitle>
                    <Activity className="h-4 w-4 text-muted-foreground" />
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold">{dashboardData.total_executions}</div>
                    <p className="text-xs text-muted-foreground">
                      排程報表: {dashboardData.scheduled_reports}
                    </p>
                  </CardContent>
                </Card>
                <Card>
                  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">資料來源</CardTitle>
                    <Database className="h-4 w-4 text-muted-foreground" />
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold">{dashboardData.total_data_sources}</div>
                    <p className="text-xs text-muted-foreground">
                      KPI指標: {dashboardData.total_kpis}
                    </p>
                  </CardContent>
                </Card>
                <Card>
                  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">系統狀態</CardTitle>
                    <Zap className="h-4 w-4 text-muted-foreground" />
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold text-green-600">健康</div>
                    <p className="text-xs text-muted-foreground">
                      CPU: {dashboardData.system_health?.cpu_usage?.toFixed(1)}%
                    </p>
                  </CardContent>
                </Card>
              </div>
            )}

            {/* Charts and Analytics */}
            {dashboardData && (
              <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                <Card>
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                      <BarChart3 className="h-5 w-5" />
                      報表分類統計
                    </CardTitle>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    {dashboardData.reports_by_category.map((item) => (
                      <div key={item.category} className="flex justify-between items-center">
                        <span className="text-sm">{getCategoryBadge(item.category)}</span>
                        <div className="flex items-center gap-2">
                          <div className="w-24 bg-gray-200 rounded-full h-2">
                            <div 
                              className="bg-blue-500 h-2 rounded-full" 
                              style={{ width: `${(item.count / dashboardData.total_reports) * 100}%` }}
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
                      執行狀態統計
                    </CardTitle>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    {dashboardData.executions_by_status.map((item) => (
                      <div key={item.status} className="flex justify-between items-center">
                        <span className="text-sm">{getExecutionStatusBadge(item.status)}</span>
                        <div className="flex items-center gap-2">
                          <div className="w-24 bg-gray-200 rounded-full h-2">
                            <div 
                              className="bg-green-500 h-2 rounded-full" 
                              style={{ width: `${(item.count / dashboardData.total_executions) * 100}%` }}
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

            {/* Recent Activities */}
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              {/* Popular Reports */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center justify-between">
                    <span className="flex items-center gap-2">
                      <TrendingUp className="h-5 w-5" />
                      熱門報表
                    </span>
                    <Button variant="ghost" size="sm" onClick={() => setActiveTab('reports')}>
                      <Eye className="h-4 w-4" />
                    </Button>
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  {dashboardData?.popular_reports.length === 0 ? (
                    <p className="text-center text-gray-500 py-4">暫無熱門報表</p>
                  ) : (
                    <div className="space-y-3">
                      {dashboardData?.popular_reports.map((report) => (
                        <div key={report.id} className="flex items-center justify-between p-3 border rounded-lg">
                          <div>
                            <p className="font-medium">{report.name}</p>
                            <div className="flex items-center gap-2 mt-1">
                              {getCategoryBadge(report.category)}
                              {getTypeBadge(report.type)}
                            </div>
                          </div>
                          <div className="text-right">
                            <p className="text-sm text-gray-500">
                              執行 {report.execute_count} 次
                            </p>
                            <p className="text-sm text-gray-500">
                              觀看 {report.view_count} 次
                            </p>
                          </div>
                        </div>
                      ))}
                    </div>
                  )}
                </CardContent>
              </Card>

              {/* Recent Executions */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center justify-between">
                    <span className="flex items-center gap-2">
                      <Clock className="h-5 w-5" />
                      最近執行
                    </span>
                    <Button variant="ghost" size="sm" onClick={() => setActiveTab('executions')}>
                      <Eye className="h-4 w-4" />
                    </Button>
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  {recentExecutions?.length === 0 ? (
                    <p className="text-center text-gray-500 py-4">暫無執行記錄</p>
                  ) : (
                    <div className="space-y-3">
                      {recentExecutions?.map((execution) => (
                        <div key={execution.id} className="flex items-center justify-between p-3 border rounded-lg">
                          <div>
                            <p className="font-medium">{execution.report?.name}</p>
                            <p className="text-sm text-gray-500">
                              {execution.executed_by_user?.full_name}
                            </p>
                            <div className="mt-1">
                              {getExecutionStatusBadge(execution.status)}
                            </div>
                          </div>
                          <div className="text-right">
                            <p className="text-sm text-gray-500">
                              {format(new Date(execution.started_at), 'MM/dd HH:mm', { locale: zhTW })}
                            </p>
                            {execution.execution_time > 0 && (
                              <p className="text-sm text-gray-500">
                                {(execution.execution_time / 1000).toFixed(1)}s
                              </p>
                            )}
                          </div>
                        </div>
                      ))}
                    </div>
                  )}
                </CardContent>
              </Card>
            </div>
          </TabsContent>

          <TabsContent value="reports" className="space-y-6">
            {/* Filters */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Filter className="h-5 w-5" />
                  篩選條件
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                  <div className="space-y-2">
                    <label className="text-sm font-medium">搜尋</label>
                    <div className="relative">
                      <Search className="absolute left-3 top-3 h-4 w-4 text-gray-400" />
                      <Input
                        placeholder="報表名稱或描述"
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                        className="pl-10"
                        onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
                      />
                    </div>
                  </div>
                  <div className="space-y-2">
                    <label className="text-sm font-medium">分類</label>
                    <Select value={categoryFilter} onValueChange={setCategoryFilter}>
                      <SelectTrigger>
                        <SelectValue placeholder="選擇分類" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="">全部</SelectItem>
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
                  <div className="space-y-2">
                    <label className="text-sm font-medium">類型</label>
                    <Select value={typeFilter} onValueChange={setTypeFilter}>
                      <SelectTrigger>
                        <SelectValue placeholder="選擇類型" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="">全部</SelectItem>
                        <SelectItem value="summary">摘要</SelectItem>
                        <SelectItem value="detail">詳細</SelectItem>
                        <SelectItem value="trend">趨勢</SelectItem>
                        <SelectItem value="comparison">比較</SelectItem>
                        <SelectItem value="dashboard">儀表板</SelectItem>
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
                        <SelectItem value="archived">已封存</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                </div>
                <div className="flex gap-2 mt-4">
                  <Button onClick={handleSearch}>
                    <Search className="mr-2 h-4 w-4" />
                    搜尋
                  </Button>
                  <Button variant="outline" onClick={handleClearFilters}>
                    清除篩選
                  </Button>
                  <Button variant="outline">
                    <Download className="mr-2 h-4 w-4" />
                    匯出
                  </Button>
                  <Button variant="outline">
                    <Upload className="mr-2 h-4 w-4" />
                    匯入
                  </Button>
                </div>
              </CardContent>
            </Card>

            {/* Reports Table */}
            <Card>
              <CardHeader>
                <CardTitle>我的報表</CardTitle>
                <CardDescription>
                  共 {reportsData?.total || 0} 個報表
                </CardDescription>
              </CardHeader>
              <CardContent>
                {isLoadingReports ? (
                  <div className="text-center py-8">載入中...</div>
                ) : reportsData?.data.length === 0 ? (
                  <div className="text-center py-8">
                    <FileText className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                    <p className="text-gray-500">暫無報表資料</p>
                  </div>
                ) : (
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead>報表名稱</TableHead>
                        <TableHead>分類</TableHead>
                        <TableHead>類型</TableHead>
                        <TableHead>狀態</TableHead>
                        <TableHead>最後執行</TableHead>
                        <TableHead>統計資料</TableHead>
                        <TableHead>操作</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {reportsData?.data.map((report) => (
                        <TableRow key={report.id} className="cursor-pointer hover:bg-gray-50">
                          <TableCell>
                            <div>
                              <p className="font-medium">{report.name}</p>
                              <p className="text-sm text-gray-500">{report.report_no}</p>
                              {report.description && (
                                <p className="text-sm text-gray-400 truncate max-w-xs">
                                  {report.description}
                                </p>
                              )}
                            </div>
                          </TableCell>
                          <TableCell>{getCategoryBadge(report.category)}</TableCell>
                          <TableCell>{getTypeBadge(report.type)}</TableCell>
                          <TableCell>{getStatusBadge(report.status)}</TableCell>
                          <TableCell>
                            {report.last_executed ? (
                              <div className="text-sm">
                                <p>{format(new Date(report.last_executed), 'MM/dd HH:mm', { locale: zhTW })}</p>
                                {report.avg_exec_time > 0 && (
                                  <p className="text-gray-500">
                                    平均 {(report.avg_exec_time / 1000).toFixed(1)}s
                                  </p>
                                )}
                              </div>
                            ) : (
                              <span className="text-gray-400">未執行過</span>
                            )}
                          </TableCell>
                          <TableCell>
                            <div className="text-sm">
                              <p>執行 {report.execute_count} 次</p>
                              <p className="text-gray-500">觀看 {report.view_count} 次</p>
                            </div>
                          </TableCell>
                          <TableCell>
                            <div className="flex items-center gap-2">
                              <Button 
                                variant="ghost" 
                                size="sm"
                                onClick={() => handleExecuteReport(report.id)}
                                title="執行報表"
                              >
                                <Play className="h-4 w-4" />
                              </Button>
                              <Button 
                                variant="ghost" 
                                size="sm"
                                onClick={() => router.push(`/reports/${report.id}`)}
                                title="檢視報表"
                              >
                                <Eye className="h-4 w-4" />
                              </Button>
                              <Button 
                                variant="ghost" 
                                size="sm"
                                onClick={() => router.push(`/reports/${report.id}/edit`)}
                                title="編輯報表"
                              >
                                <Edit className="h-4 w-4" />
                              </Button>
                              <Button 
                                variant="ghost" 
                                size="sm"
                                onClick={() => handleDuplicateReport(report.id)}
                                title="複製報表"
                              >
                                <Copy className="h-4 w-4" />
                              </Button>
                            </div>
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                )}

                {/* Pagination */}
                {reportsData && reportsData.total > pageSize && (
                  <div className="flex items-center justify-between mt-4">
                    <p className="text-sm text-gray-500">
                      顯示 {(page - 1) * pageSize + 1} 到 {Math.min(page * pageSize, reportsData.total)} 項，
                      共 {reportsData.total} 項
                    </p>
                    <div className="flex gap-2">
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => setPage(page - 1)}
                        disabled={page === 1}
                      >
                        上一頁
                      </Button>
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => setPage(page + 1)}
                        disabled={page * pageSize >= reportsData.total}
                      >
                        下一頁
                      </Button>
                    </div>
                  </div>
                )}
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="templates">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <FileText className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">範本庫功能開發中</p>
                  <Button className="mt-4" onClick={() => router.push('/reports/templates')}>
                    前往範本庫
                  </Button>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="executions">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <Activity className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">執行記錄功能開發中</p>
                  <Button className="mt-4" onClick={() => router.push('/reports/executions')}>
                    前往執行記錄
                  </Button>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="settings">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <Settings className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">系統設定功能開發中</p>
                  <Button className="mt-4" onClick={() => router.push('/reports/settings')}>
                    前往系統設定
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