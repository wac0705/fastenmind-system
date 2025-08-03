'use client'

import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { useRouter } from 'next/navigation'
import DashboardLayout from '@/components/layout/DashboardLayout'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
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
  Search,
  Filter,
  Download,
  RefreshCw,
  Eye,
  Trash2,
  Play,
  Square,
  Clock,
  CheckCircle,
  XCircle,
  AlertCircle,
  Calendar,
  User,
  FileText,
  Activity,
  Timer,
  Database,
  TrendingUp,
  TrendingDown
} from 'lucide-react'
import reportService from '@/services/report.service'
import { format } from 'date-fns'
import { zhTW } from 'date-fns/locale'

export default function ReportExecutionsPage() {
  const router = useRouter()
  const [searchQuery, setSearchQuery] = useState('')
  const [reportFilter, setReportFilter] = useState('')
  const [statusFilter, setStatusFilter] = useState('')
  const [userFilter, setUserFilter] = useState('')
  const [dateRange, setDateRange] = useState({ start: '', end: '' })
  const [page, setPage] = useState(1)
  const pageSize = 20

  // Fetch executions list
  const { data: executionsData, isLoading, refetch } = useQuery({
    queryKey: ['report-executions', page, searchQuery, reportFilter, statusFilter, userFilter, dateRange],
    queryFn: () => reportService.listExecutions({
      page,
      page_size: pageSize,
      search: searchQuery || undefined,
      report_id: reportFilter || undefined,
      status: statusFilter || undefined,
      executed_by: userFilter || undefined,
      start_date: dateRange.start || undefined,
      end_date: dateRange.end || undefined,
    }),
  })

  // Fetch execution statistics
  const { data: stats } = useQuery({
    queryKey: ['execution-stats'],
    queryFn: () => reportService.getExecutionStats(),
  })

  const getStatusBadge = (status: string) => {
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
    refetch()
  }

  const handleClearFilters = () => {
    setSearchQuery('')
    setReportFilter('')
    setStatusFilter('')
    setUserFilter('')
    setDateRange({ start: '', end: '' })
    setPage(1)
    refetch()
  }

  const handleViewExecution = (executionId: string) => {
    router.push(`/reports/executions/${executionId}`)
  }

  const handleRerunExecution = async (executionId: string) => {
    try {
      await reportService.rerunExecution(executionId)
      refetch()
    } catch (error) {
      console.error('Failed to rerun execution:', error)
    }
  }

  const handleCancelExecution = async (executionId: string) => {
    try {
      await reportService.cancelExecution(executionId)
      refetch()
    } catch (error) {
      console.error('Failed to cancel execution:', error)
    }
  }

  const handleExportExecutions = async () => {
    try {
      const blob = await reportService.exportExecutions({
        format: 'csv',
        status: statusFilter || undefined,
        start_date: dateRange.start || undefined,
        end_date: dateRange.end || undefined,
      })
      const url = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `executions_${format(new Date(), 'yyyyMMdd_HHmmss')}.csv`
      document.body.appendChild(a)
      a.click()
      window.URL.revokeObjectURL(url)
      document.body.removeChild(a)
    } catch (error) {
      console.error('Failed to export executions:', error)
    }
  }

  return (
    <DashboardLayout>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex justify-between items-center">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">執行記錄</h1>
            <p className="mt-2 text-gray-600">查看所有報表執行記錄與狀態</p>
          </div>
          <div className="flex gap-2">
            <Button variant="outline" onClick={() => router.push('/reports')}>
              <FileText className="mr-2 h-4 w-4" />
              報表中心
            </Button>
            <Button variant="outline" onClick={() => refetch()}>
              <RefreshCw className="mr-2 h-4 w-4" />
              重新整理
            </Button>
          </div>
        </div>

        {/* Statistics */}
        {stats && (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-4">
            <Card>
              <CardHeader className="pb-3">
                <CardTitle className="text-sm flex items-center gap-2">
                  <Activity className="h-4 w-4" />
                  總執行次數
                </CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-2xl font-bold">{stats.total_executions}</p>
                <div className="flex items-center gap-1 mt-1">
                  {stats.execution_trend > 0 ? (
                    <>
                      <TrendingUp className="h-4 w-4 text-green-600" />
                      <span className="text-sm text-green-600">+{stats.execution_trend}%</span>
                    </>
                  ) : (
                    <>
                      <TrendingDown className="h-4 w-4 text-red-600" />
                      <span className="text-sm text-red-600">{stats.execution_trend}%</span>
                    </>
                  )}
                </div>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="pb-3">
                <CardTitle className="text-sm flex items-center gap-2">
                  <CheckCircle className="h-4 w-4 text-green-600" />
                  成功執行
                </CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-2xl font-bold">{stats.successful_executions}</p>
                <p className="text-sm text-gray-500">{stats.success_rate.toFixed(1)}% 成功率</p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="pb-3">
                <CardTitle className="text-sm flex items-center gap-2">
                  <XCircle className="h-4 w-4 text-red-600" />
                  失敗執行
                </CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-2xl font-bold">{stats.failed_executions}</p>
                <p className="text-sm text-gray-500">{stats.failure_rate.toFixed(1)}% 失敗率</p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="pb-3">
                <CardTitle className="text-sm flex items-center gap-2">
                  <RefreshCw className="h-4 w-4 text-yellow-600" />
                  執行中
                </CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-2xl font-bold">{stats.running_executions}</p>
                <p className="text-sm text-gray-500">當前執行中</p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="pb-3">
                <CardTitle className="text-sm flex items-center gap-2">
                  <Timer className="h-4 w-4" />
                  平均耗時
                </CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-2xl font-bold">{(stats.avg_execution_time / 1000).toFixed(1)}s</p>
                <p className="text-sm text-gray-500">平均執行時間</p>
              </CardContent>
            </Card>
          </div>
        )}

        {/* Filters */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Filter className="h-5 w-5" />
              篩選條件
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              <div className="space-y-2">
                <label className="text-sm font-medium">搜尋</label>
                <div className="relative">
                  <Search className="absolute left-3 top-3 h-4 w-4 text-gray-400" />
                  <Input
                    placeholder="執行編號或報表名稱"
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    className="pl-10"
                    onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
                  />
                </div>
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">狀態</label>
                <Select value={statusFilter} onValueChange={setStatusFilter}>
                  <SelectTrigger>
                    <SelectValue placeholder="選擇狀態" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="">全部</SelectItem>
                    <SelectItem value="pending">等待中</SelectItem>
                    <SelectItem value="running">執行中</SelectItem>
                    <SelectItem value="completed">已完成</SelectItem>
                    <SelectItem value="failed">失敗</SelectItem>
                    <SelectItem value="cancelled">已取消</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">日期範圍</label>
                <div className="grid grid-cols-2 gap-2">
                  <Input
                    type="date"
                    value={dateRange.start}
                    onChange={(e) => setDateRange({ ...dateRange, start: e.target.value })}
                  />
                  <Input
                    type="date"
                    value={dateRange.end}
                    onChange={(e) => setDateRange({ ...dateRange, end: e.target.value })}
                  />
                </div>
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
              <Button variant="outline" onClick={handleExportExecutions}>
                <Download className="mr-2 h-4 w-4" />
                匯出
              </Button>
            </div>
          </CardContent>
        </Card>

        {/* Executions Table */}
        <Card>
          <CardHeader>
            <CardTitle>執行記錄列表</CardTitle>
            <CardDescription>
              共 {executionsData?.total || 0} 筆執行記錄
            </CardDescription>
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <div className="text-center py-8">
                <RefreshCw className="h-8 w-8 animate-spin mx-auto mb-4 text-gray-400" />
                <p className="text-gray-500">載入中...</p>
              </div>
            ) : executionsData?.data.length === 0 ? (
              <div className="text-center py-8">
                <Activity className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                <p className="text-gray-500">暫無執行記錄</p>
              </div>
            ) : (
              <>
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>執行編號</TableHead>
                      <TableHead>報表名稱</TableHead>
                      <TableHead>執行者</TableHead>
                      <TableHead>狀態</TableHead>
                      <TableHead>開始時間</TableHead>
                      <TableHead>耗時</TableHead>
                      <TableHead>資料筆數</TableHead>
                      <TableHead>操作</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {executionsData?.data.map((execution) => (
                      <TableRow key={execution.id} className="cursor-pointer hover:bg-gray-50">
                        <TableCell>
                          <p className="font-mono">{execution.id.slice(-8)}</p>
                        </TableCell>
                        <TableCell>
                          <div>
                            <p className="font-medium">{execution.report?.name}</p>
                            <p className="text-sm text-gray-500">{execution.report?.report_no}</p>
                          </div>
                        </TableCell>
                        <TableCell>
                          <div>
                            <p>{execution.executed_by_user?.full_name}</p>
                            <p className="text-sm text-gray-500">{execution.executed_by_user?.email}</p>
                          </div>
                        </TableCell>
                        <TableCell>{getStatusBadge(execution.status)}</TableCell>
                        <TableCell>
                          <div>
                            <p>{format(new Date(execution.started_at), 'yyyy/MM/dd', { locale: zhTW })}</p>
                            <p className="text-sm text-gray-500">
                              {format(new Date(execution.started_at), 'HH:mm:ss', { locale: zhTW })}
                            </p>
                          </div>
                        </TableCell>
                        <TableCell>
                          {execution.execution_time ? (
                            <div>
                              <p>{(execution.execution_time / 1000).toFixed(1)}s</p>
                              {execution.execution_time > 10000 && (
                                <p className="text-sm text-yellow-600">較慢</p>
                              )}
                            </div>
                          ) : execution.status === 'running' ? (
                            <span className="text-gray-500">執行中...</span>
                          ) : (
                            <span className="text-gray-400">-</span>
                          )}
                        </TableCell>
                        <TableCell>
                          {(execution as any).result_count !== null ? (
                            <div className="flex items-center gap-1">
                              <Database className="h-4 w-4 text-gray-400" />
                              <span>{((execution as any).result_count || 0).toLocaleString()}</span>
                            </div>
                          ) : (
                            <span className="text-gray-400">-</span>
                          )}
                        </TableCell>
                        <TableCell>
                          <div className="flex items-center gap-2">
                            <Button 
                              variant="ghost" 
                              size="sm"
                              onClick={() => handleViewExecution(execution.id)}
                              title="檢視詳情"
                            >
                              <Eye className="h-4 w-4" />
                            </Button>
                            {execution.status === 'completed' && (
                              <Button 
                                variant="ghost" 
                                size="sm"
                                onClick={() => handleRerunExecution(execution.id)}
                                title="重新執行"
                              >
                                <RefreshCw className="h-4 w-4" />
                              </Button>
                            )}
                            {execution.status === 'running' && (
                              <Button 
                                variant="ghost" 
                                size="sm"
                                onClick={() => handleCancelExecution(execution.id)}
                                title="取消執行"
                              >
                                <Square className="h-4 w-4" />
                              </Button>
                            )}
                          </div>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>

                {/* Pagination */}
                {executionsData && executionsData.total > pageSize && (
                  <div className="flex items-center justify-between mt-4">
                    <p className="text-sm text-gray-500">
                      顯示 {(page - 1) * pageSize + 1} 到 {Math.min(page * pageSize, executionsData.total)} 項，
                      共 {executionsData.total} 項
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
                        disabled={page * pageSize >= executionsData.total}
                      >
                        下一頁
                      </Button>
                    </div>
                  </div>
                )}
              </>
            )}
          </CardContent>
        </Card>
      </div>
    </DashboardLayout>
  )
}