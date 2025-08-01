'use client'

import { useState } from 'react'
import { useParams, useRouter } from 'next/navigation'
import { useQuery } from '@tanstack/react-query'
import DashboardLayout from '@/components/layout/DashboardLayout'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Progress } from '@/components/ui/progress'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { useToast } from '@/components/ui/use-toast'
import {
  ArrowLeft,
  Download,
  Share2,
  Printer,
  Mail,
  Copy,
  RefreshCw,
  CheckCircle,
  XCircle,
  AlertCircle,
  Clock,
  Calendar,
  User,
  Hash,
  FileText,
  FileJson,
  FileCsv,
  FileSpreadsheet,
  ChevronLeft,
  ChevronRight,
  ZoomIn,
  ZoomOut,
  Maximize2,
  Database,
  Timer,
  Activity,
  Info,
  BarChart3,
  LineChart,
  PieChart,
  TableIcon
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

export default function ReportExecutionPage() {
  const params = useParams()
  const router = useRouter()
  const { toast } = useToast()
  const executionId = params.id as string
  
  const [activeTab, setActiveTab] = useState('result')
  const [exportFormat, setExportFormat] = useState('pdf')
  const [zoom, setZoom] = useState(100)
  const [currentPage, setCurrentPage] = useState(1)

  // Fetch execution details
  const { data: execution, isLoading, refetch } = useQuery({
    queryKey: ['report-execution', executionId],
    queryFn: () => reportService.getReportExecution(executionId),
  })

  const getStatusBadge = (status: string) => {
    const statusConfig: Record<string, { label: string; variant: any; icon: any }> = {
      pending: { label: '等待中', variant: 'secondary', icon: Clock },
      running: { label: '執行中', variant: 'warning', icon: RefreshCw },
      completed: { label: '已完成', variant: 'success', icon: CheckCircle },
      failed: { label: '失敗', variant: 'destructive', icon: XCircle },
      cancelled: { label: '已取消', variant: 'secondary', icon: AlertCircle },
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

  const handleExport = async () => {
    try {
      const blob = await reportService.exportExecution(executionId, exportFormat)
      const url = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `${execution?.execution_no || 'execution'}.${exportFormat}`
      document.body.appendChild(a)
      a.click()
      window.URL.revokeObjectURL(url)
      document.body.removeChild(a)
      toast({ title: '匯出成功' })
    } catch (error) {
      toast({
        title: '匯出失敗',
        variant: 'destructive',
      })
    }
  }

  const handleShare = () => {
    const url = `${window.location.origin}/reports/executions/${executionId}/public`
    navigator.clipboard.writeText(url)
    toast({ title: '分享連結已複製' })
  }

  const handlePrint = () => {
    window.print()
  }

  const handleEmail = () => {
    // Implement email functionality
    toast({ title: '郵件功能開發中' })
  }

  const renderChart = (component: any) => {
    const chartOptions = {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: {
          position: 'bottom' as const,
        },
      },
    }

    // Sample data - would come from execution result
    const sampleData = {
      labels: ['一月', '二月', '三月', '四月', '五月', '六月'],
      datasets: [
        {
          label: '銷售額',
          data: [65, 59, 80, 81, 56, 55],
          backgroundColor: 'rgba(59, 130, 246, 0.5)',
          borderColor: 'rgb(59, 130, 246)',
          tension: 0.4,
        },
      ],
    }

    switch (component.type) {
      case 'chart_bar':
        return <Bar data={sampleData} options={chartOptions} />
      case 'chart_line':
        return <Line data={sampleData} options={chartOptions} />
      case 'chart_pie':
        return <Pie data={sampleData} options={chartOptions} />
      default:
        return null
    }
  }

  const renderTable = (data: any) => {
    // Sample table data - would come from execution result
    const headers = ['產品', '數量', '單價', '總額']
    const rows = [
      ['螺絲 M8x30', '1000', '$0.05', '$50.00'],
      ['螺帽 M8', '1000', '$0.03', '$30.00'],
      ['墊圈 M8', '2000', '$0.01', '$20.00'],
    ]

    return (
      <div className="overflow-x-auto">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              {headers.map((header, index) => (
                <th
                  key={index}
                  className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                >
                  {header}
                </th>
              ))}
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {rows.map((row, rowIndex) => (
              <tr key={rowIndex}>
                {row.map((cell, cellIndex) => (
                  <td key={cellIndex} className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                    {cell}
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    )
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

  if (!execution) {
    return (
      <DashboardLayout>
        <div className="flex items-center justify-center h-64">
          <div className="text-center">
            <FileText className="h-12 w-12 mx-auto mb-4 text-gray-300" />
            <p className="text-gray-500">找不到執行記錄</p>
            <Button className="mt-4" onClick={() => router.push('/reports')}>
              返回報表列表
            </Button>
          </div>
        </div>
      </DashboardLayout>
    )
  }

  // Calculate progress percentage
  const progress = execution.status === 'running' 
    ? Math.floor((Date.now() - new Date(execution.started_at).getTime()) / 1000 % 100)
    : execution.status === 'completed' ? 100 : 0

  return (
    <DashboardLayout>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-4">
            <Button variant="ghost" size="icon" onClick={() => router.back()}>
              <ArrowLeft className="h-4 w-4" />
            </Button>
            <div>
              <div className="flex items-center gap-2 mb-2">
                <h1 className="text-3xl font-bold">{execution.report?.name}</h1>
                {getStatusBadge(execution.status)}
              </div>
              <p className="text-gray-600">執行編號：{execution.execution_no}</p>
            </div>
          </div>
          <div className="flex items-center gap-2">
            <Button variant="outline" onClick={handleEmail}>
              <Mail className="mr-2 h-4 w-4" />
              郵件
            </Button>
            <Button variant="outline" onClick={handlePrint}>
              <Printer className="mr-2 h-4 w-4" />
              列印
            </Button>
            <Button variant="outline" onClick={handleShare}>
              <Share2 className="mr-2 h-4 w-4" />
              分享
            </Button>
            <Select value={exportFormat} onValueChange={setExportFormat}>
              <SelectTrigger className="w-32">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="pdf">PDF</SelectItem>
                <SelectItem value="excel">Excel</SelectItem>
                <SelectItem value="csv">CSV</SelectItem>
                <SelectItem value="json">JSON</SelectItem>
              </SelectContent>
            </Select>
            <Button onClick={handleExport}>
              <Download className="mr-2 h-4 w-4" />
              匯出
            </Button>
          </div>
        </div>

        {/* Status Progress */}
        {execution.status === 'running' && (
          <Card>
            <CardContent className="pt-6">
              <div className="space-y-2">
                <div className="flex justify-between text-sm">
                  <span>執行進度</span>
                  <span>{progress}%</span>
                </div>
                <Progress value={progress} />
                <p className="text-sm text-gray-500">報表正在執行中，請稍候...</p>
              </div>
            </CardContent>
          </Card>
        )}

        {/* Info Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="text-sm flex items-center gap-2">
                <User className="h-4 w-4" />
                執行者
              </CardTitle>
            </CardHeader>
            <CardContent>
              <p>{execution.executed_by_user?.full_name}</p>
            </CardContent>
          </Card>
          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="text-sm flex items-center gap-2">
                <Calendar className="h-4 w-4" />
                執行時間
              </CardTitle>
            </CardHeader>
            <CardContent>
              <p>{format(new Date(execution.started_at), 'yyyy/MM/dd HH:mm', { locale: zhTW })}</p>
            </CardContent>
          </Card>
          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="text-sm flex items-center gap-2">
                <Timer className="h-4 w-4" />
                執行耗時
              </CardTitle>
            </CardHeader>
            <CardContent>
              <p>
                {execution.execution_time 
                  ? `${(execution.execution_time / 1000).toFixed(1)} 秒`
                  : '計算中...'}
              </p>
            </CardContent>
          </Card>
          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="text-sm flex items-center gap-2">
                <Database className="h-4 w-4" />
                資料筆數
              </CardTitle>
            </CardHeader>
            <CardContent>
              <p>{execution.result_count || 0} 筆</p>
            </CardContent>
          </Card>
        </div>

        <Tabs value={activeTab} onValueChange={setActiveTab}>
          <TabsList>
            <TabsTrigger value="result">執行結果</TabsTrigger>
            <TabsTrigger value="parameters">執行參數</TabsTrigger>
            <TabsTrigger value="logs">執行日誌</TabsTrigger>
            <TabsTrigger value="errors">錯誤訊息</TabsTrigger>
          </TabsList>

          <TabsContent value="result" className="space-y-4">
            {execution.status === 'completed' ? (
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center justify-between">
                    <span>報表內容</span>
                    <div className="flex items-center gap-2">
                      <Button variant="outline" size="sm" onClick={() => setZoom(zoom - 10)} disabled={zoom <= 50}>
                        <ZoomOut className="h-4 w-4" />
                      </Button>
                      <span className="text-sm font-medium">{zoom}%</span>
                      <Button variant="outline" size="sm" onClick={() => setZoom(zoom + 10)} disabled={zoom >= 200}>
                        <ZoomIn className="h-4 w-4" />
                      </Button>
                      <Button variant="outline" size="sm">
                        <Maximize2 className="h-4 w-4" />
                      </Button>
                    </div>
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div 
                    className="border rounded-lg p-8 bg-white overflow-auto"
                    style={{ transform: `scale(${zoom / 100})`, transformOrigin: 'top left' }}
                  >
                    {/* Render report components */}
                    {execution.result?.components?.map((component: any, index: number) => (
                      <div key={index} className="mb-8">
                        {component.type === 'text' && (
                          <div>
                            <h3 className="text-lg font-semibold mb-2">{component.name}</h3>
                            <p className="text-gray-700">{component.content}</p>
                          </div>
                        )}
                        {component.type === 'table' && (
                          <div>
                            <h3 className="text-lg font-semibold mb-2">{component.name}</h3>
                            {renderTable(component.data)}
                          </div>
                        )}
                        {component.type.startsWith('chart_') && (
                          <div>
                            <h3 className="text-lg font-semibold mb-2">{component.name}</h3>
                            <div className="h-80">
                              {renderChart(component)}
                            </div>
                          </div>
                        )}
                        {component.type === 'kpi' && (
                          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                            {component.data?.map((kpi: any, kpiIndex: number) => (
                              <Card key={kpiIndex}>
                                <CardHeader className="pb-3">
                                  <CardTitle className="text-sm">{kpi.label}</CardTitle>
                                </CardHeader>
                                <CardContent>
                                  <p className="text-2xl font-bold">{kpi.value}</p>
                                  {kpi.change && (
                                    <p className={`text-sm ${kpi.change > 0 ? 'text-green-600' : 'text-red-600'}`}>
                                      {kpi.change > 0 ? '+' : ''}{kpi.change}%
                                    </p>
                                  )}
                                </CardContent>
                              </Card>
                            ))}
                          </div>
                        )}
                      </div>
                    ))}
                    
                    {/* Default content if no components */}
                    {!execution.result?.components && (
                      <div className="space-y-6">
                        <div>
                          <h2 className="text-2xl font-bold mb-4">銷售報表</h2>
                          <p className="text-gray-600 mb-6">
                            報表期間：{format(new Date(), 'yyyy年MM月', { locale: zhTW })}
                          </p>
                        </div>
                        
                        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                          <Card>
                            <CardHeader className="pb-3">
                              <CardTitle className="text-sm">總銷售額</CardTitle>
                            </CardHeader>
                            <CardContent>
                              <p className="text-2xl font-bold">$1,234,567</p>
                              <p className="text-sm text-green-600">+12.5%</p>
                            </CardContent>
                          </Card>
                          <Card>
                            <CardHeader className="pb-3">
                              <CardTitle className="text-sm">訂單數量</CardTitle>
                            </CardHeader>
                            <CardContent>
                              <p className="text-2xl font-bold">342</p>
                              <p className="text-sm text-green-600">+8.3%</p>
                            </CardContent>
                          </Card>
                          <Card>
                            <CardHeader className="pb-3">
                              <CardTitle className="text-sm">新客戶</CardTitle>
                            </CardHeader>
                            <CardContent>
                              <p className="text-2xl font-bold">28</p>
                              <p className="text-sm text-red-600">-5.2%</p>
                            </CardContent>
                          </Card>
                          <Card>
                            <CardHeader className="pb-3">
                              <CardTitle className="text-sm">平均單價</CardTitle>
                            </CardHeader>
                            <CardContent>
                              <p className="text-2xl font-bold">$3,612</p>
                              <p className="text-sm text-green-600">+3.8%</p>
                            </CardContent>
                          </Card>
                        </div>

                        <div className="h-80">
                          <h3 className="text-lg font-semibold mb-4">月度銷售趨勢</h3>
                          <Line
                            data={{
                              labels: ['一月', '二月', '三月', '四月', '五月', '六月'],
                              datasets: [
                                {
                                  label: '銷售額',
                                  data: [180000, 195000, 210000, 225000, 215000, 234567],
                                  borderColor: 'rgb(59, 130, 246)',
                                  backgroundColor: 'rgba(59, 130, 246, 0.1)',
                                  tension: 0.4,
                                },
                              ],
                            }}
                            options={{
                              responsive: true,
                              maintainAspectRatio: false,
                              plugins: {
                                legend: {
                                  display: false,
                                },
                              },
                            }}
                          />
                        </div>
                      </div>
                    )}
                  </div>
                  
                  {/* Pagination for multi-page reports */}
                  {execution.result?.total_pages > 1 && (
                    <div className="flex items-center justify-center gap-2 mt-4">
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => setCurrentPage(currentPage - 1)}
                        disabled={currentPage === 1}
                      >
                        <ChevronLeft className="h-4 w-4" />
                      </Button>
                      <span className="text-sm">
                        第 {currentPage} 頁，共 {execution.result.total_pages} 頁
                      </span>
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => setCurrentPage(currentPage + 1)}
                        disabled={currentPage === execution.result.total_pages}
                      >
                        <ChevronRight className="h-4 w-4" />
                      </Button>
                    </div>
                  )}
                </CardContent>
              </Card>
            ) : execution.status === 'failed' ? (
              <Card>
                <CardContent className="py-12">
                  <div className="text-center">
                    <XCircle className="h-12 w-12 mx-auto mb-4 text-red-500" />
                    <p className="text-lg font-medium text-gray-900">報表執行失敗</p>
                    <p className="text-gray-500 mt-2">{execution.error_message || '執行過程中發生錯誤'}</p>
                    <Button className="mt-4" onClick={() => router.push(`/reports/${execution.report_id}`)}>
                      返回報表
                    </Button>
                  </div>
                </CardContent>
              </Card>
            ) : (
              <Card>
                <CardContent className="py-12">
                  <div className="text-center">
                    <Activity className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                    <p className="text-gray-500">報表尚未執行完成</p>
                  </div>
                </CardContent>
              </Card>
            )}
          </TabsContent>

          <TabsContent value="parameters" className="space-y-4">
            <Card>
              <CardHeader>
                <CardTitle>執行參數</CardTitle>
                <CardDescription>本次執行使用的參數設定</CardDescription>
              </CardHeader>
              <CardContent>
                {execution.parameters ? (
                  <div className="space-y-4">
                    {Object.entries(execution.parameters).map(([key, value]) => (
                      <div key={key} className="flex justify-between py-2 border-b">
                        <span className="font-medium">{key}</span>
                        <span className="text-gray-600">{String(value)}</span>
                      </div>
                    ))}
                  </div>
                ) : (
                  <p className="text-center text-gray-500 py-4">無執行參數</p>
                )}
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="logs" className="space-y-4">
            <Card>
              <CardHeader>
                <CardTitle>執行日誌</CardTitle>
                <CardDescription>報表執行過程的詳細日誌</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="bg-gray-900 text-gray-100 p-4 rounded-lg font-mono text-sm overflow-auto max-h-96">
                  {execution.logs ? (
                    <pre>{execution.logs}</pre>
                  ) : (
                    <>
                      <p>[2024-01-15 10:30:00] 開始執行報表...</p>
                      <p>[2024-01-15 10:30:01] 連接資料庫成功</p>
                      <p>[2024-01-15 10:30:02] 開始查詢資料...</p>
                      <p>[2024-01-15 10:30:05] 查詢完成，共 342 筆記錄</p>
                      <p>[2024-01-15 10:30:06] 開始生成報表...</p>
                      <p>[2024-01-15 10:30:08] 報表生成完成</p>
                      <p>[2024-01-15 10:30:08] 執行成功</p>
                    </>
                  )}
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="errors" className="space-y-4">
            <Card>
              <CardHeader>
                <CardTitle>錯誤訊息</CardTitle>
                <CardDescription>執行過程中的錯誤和警告</CardDescription>
              </CardHeader>
              <CardContent>
                {execution.error_message ? (
                  <div className="space-y-4">
                    <div className="flex items-start gap-3 p-4 bg-red-50 border border-red-200 rounded-lg">
                      <XCircle className="h-5 w-5 text-red-500 mt-0.5" />
                      <div>
                        <p className="font-medium text-red-900">執行錯誤</p>
                        <p className="text-sm text-red-700 mt-1">{execution.error_message}</p>
                        {execution.error_details && (
                          <pre className="text-xs text-red-600 mt-2 overflow-auto">
                            {JSON.stringify(execution.error_details, null, 2)}
                          </pre>
                        )}
                      </div>
                    </div>
                  </div>
                ) : (
                  <div className="text-center py-8">
                    <CheckCircle className="h-12 w-12 mx-auto mb-4 text-green-500" />
                    <p className="text-gray-500">無錯誤訊息</p>
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