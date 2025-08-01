'use client'

import { useState } from 'react'
import { useParams, useRouter } from 'next/navigation'
import { useQuery, useMutation } from '@tanstack/react-query'
import DashboardLayout from '@/components/layout/DashboardLayout'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Separator } from '@/components/ui/separator'
import { Progress } from '@/components/ui/progress'
import { useToast } from '@/components/ui/use-toast'
import {
  ArrowLeft,
  Play,
  Pause,
  Square,
  Download,
  Share2,
  Edit,
  Trash2,
  Clock,
  Calendar,
  RefreshCw,
  FileText,
  Database,
  Filter,
  Settings,
  Eye,
  EyeOff,
  CheckCircle,
  XCircle,
  AlertCircle,
  Info,
  Mail,
  MessageSquare,
  Printer,
  Copy,
  BarChart3,
  LineChart,
  PieChart,
  TableIcon,
  Layout,
  Code,
  FileJson,
  FileCsv,
  FileSpreadsheet,
  FileImage,
  User,
  Building,
  Hash
} from 'lucide-react'
import reportService from '@/services/report.service'
import { format } from 'date-fns'
import { zhTW } from 'date-fns/locale'

export default function ReportDetailPage() {
  const params = useParams()
  const router = useRouter()
  const { toast } = useToast()
  const reportId = params.id as string
  
  const [activeTab, setActiveTab] = useState('preview')
  const [exportFormat, setExportFormat] = useState('pdf')
  const [isExecuting, setIsExecuting] = useState(false)
  const [executionProgress, setExecutionProgress] = useState(0)

  // Fetch report details
  const { data: report, isLoading, refetch } = useQuery({
    queryKey: ['report', reportId],
    queryFn: () => reportService.getReport(reportId),
  })

  // Fetch report executions
  const { data: executions } = useQuery({
    queryKey: ['report-executions', reportId],
    queryFn: () => reportService.getReportExecutions(reportId, { page: 1, page_size: 10 }),
  })

  // Execute report mutation
  const executeMutation = useMutation({
    mutationFn: (params: any) => reportService.executeReport(reportId, params),
    onSuccess: (data) => {
      toast({ title: '報表執行成功' })
      refetch()
      setIsExecuting(false)
      setExecutionProgress(100)
      // Redirect to execution result
      router.push(`/reports/executions/${data.id}`)
    },
    onError: (error: any) => {
      toast({
        title: '執行失敗',
        description: error.response?.data?.message || '執行報表時發生錯誤',
        variant: 'destructive',
      })
      setIsExecuting(false)
      setExecutionProgress(0)
    },
  })

  // Delete report mutation
  const deleteMutation = useMutation({
    mutationFn: () => reportService.deleteReport(reportId),
    onSuccess: () => {
      toast({ title: '報表刪除成功' })
      router.push('/reports')
    },
    onError: (error: any) => {
      toast({
        title: '刪除失敗',
        description: error.response?.data?.message || '刪除報表時發生錯誤',
        variant: 'destructive',
      })
    },
  })

  const handleExecute = () => {
    setIsExecuting(true)
    setExecutionProgress(0)
    
    // Simulate progress
    const interval = setInterval(() => {
      setExecutionProgress((prev) => {
        if (prev >= 90) {
          clearInterval(interval)
          return 90
        }
        return prev + 10
      })
    }, 500)

    executeMutation.mutate({})
  }

  const handleExport = async () => {
    try {
      const blob = await reportService.exportReport(reportId, exportFormat)
      const url = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `${report?.report_no || 'report'}.${exportFormat}`
      document.body.appendChild(a)
      a.click()
      window.URL.revokeObjectURL(url)
      document.body.removeChild(a)
      toast({ title: '報表匯出成功' })
    } catch (error) {
      toast({
        title: '匯出失敗',
        variant: 'destructive',
      })
    }
  }

  const handleShare = () => {
    // Implement share functionality
    toast({ title: '分享連結已複製' })
  }

  const handleDelete = () => {
    if (confirm(`確定要刪除報表「${report?.name}」嗎？此操作無法復原。`)) {
      deleteMutation.mutate()
    }
  }

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

  const getCategoryIcon = (category: string) => {
    const icons: Record<string, any> = {
      sales: BarChart3,
      finance: LineChart,
      production: PieChart,
      inventory: TableIcon,
      supplier: Building,
      customer: User,
      system: Database,
    }
    return icons[category] || FileText
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

  if (!report) {
    return (
      <DashboardLayout>
        <div className="flex items-center justify-center h-64">
          <div className="text-center">
            <FileText className="h-12 w-12 mx-auto mb-4 text-gray-300" />
            <p className="text-gray-500">找不到報表</p>
            <Button className="mt-4" onClick={() => router.push('/reports')}>
              返回報表列表
            </Button>
          </div>
        </div>
      </DashboardLayout>
    )
  }

  const CategoryIcon = getCategoryIcon(report.category)

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
                <CategoryIcon className="h-6 w-6 text-gray-500" />
                <h1 className="text-3xl font-bold">{report.name}</h1>
                {getStatusBadge(report.status)}
              </div>
              <p className="text-gray-600">{report.description}</p>
            </div>
          </div>
          <div className="flex items-center gap-2">
            <Button variant="outline" onClick={handleShare}>
              <Share2 className="mr-2 h-4 w-4" />
              分享
            </Button>
            <Button variant="outline" onClick={() => router.push(`/reports/${reportId}/edit`)}>
              <Edit className="mr-2 h-4 w-4" />
              編輯
            </Button>
            <Button variant="outline" onClick={handleDelete}>
              <Trash2 className="mr-2 h-4 w-4" />
              刪除
            </Button>
            <Button onClick={handleExecute} disabled={isExecuting}>
              {isExecuting ? (
                <>
                  <RefreshCw className="mr-2 h-4 w-4 animate-spin" />
                  執行中...
                </>
              ) : (
                <>
                  <Play className="mr-2 h-4 w-4" />
                  執行報表
                </>
              )}
            </Button>
          </div>
        </div>

        {/* Execution Progress */}
        {isExecuting && (
          <Card>
            <CardContent className="pt-6">
              <div className="space-y-2">
                <div className="flex justify-between text-sm">
                  <span>執行進度</span>
                  <span>{executionProgress}%</span>
                </div>
                <Progress value={executionProgress} />
                <p className="text-sm text-gray-500">正在執行報表，請稍候...</p>
              </div>
            </CardContent>
          </Card>
        )}

        {/* Info Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="text-sm flex items-center gap-2">
                <Hash className="h-4 w-4" />
                報表編號
              </CardTitle>
            </CardHeader>
            <CardContent>
              <p className="font-mono">{report.report_no}</p>
            </CardContent>
          </Card>
          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="text-sm flex items-center gap-2">
                <Calendar className="h-4 w-4" />
                建立時間
              </CardTitle>
            </CardHeader>
            <CardContent>
              <p>{format(new Date(report.created_at), 'yyyy/MM/dd HH:mm', { locale: zhTW })}</p>
            </CardContent>
          </Card>
          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="text-sm flex items-center gap-2">
                <Clock className="h-4 w-4" />
                最後執行
              </CardTitle>
            </CardHeader>
            <CardContent>
              <p>
                {report.last_executed
                  ? format(new Date(report.last_executed), 'yyyy/MM/dd HH:mm', { locale: zhTW })
                  : '尚未執行'}
              </p>
            </CardContent>
          </Card>
          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="text-sm flex items-center gap-2">
                <Activity className="h-4 w-4" />
                執行統計
              </CardTitle>
            </CardHeader>
            <CardContent>
              <p>
                {report.execute_count} 次 / 平均 {(report.avg_exec_time / 1000).toFixed(1)}s
              </p>
            </CardContent>
          </Card>
        </div>

        <Tabs value={activeTab} onValueChange={setActiveTab}>
          <TabsList>
            <TabsTrigger value="preview">預覽</TabsTrigger>
            <TabsTrigger value="parameters">參數設定</TabsTrigger>
            <TabsTrigger value="executions">執行記錄</TabsTrigger>
            <TabsTrigger value="schedule">排程設定</TabsTrigger>
            <TabsTrigger value="permissions">權限管理</TabsTrigger>
          </TabsList>

          <TabsContent value="preview" className="space-y-4">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center justify-between">
                  <span>報表預覽</span>
                  <div className="flex items-center gap-2">
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
                    <Button variant="outline" onClick={handleExport}>
                      <Download className="mr-2 h-4 w-4" />
                      匯出
                    </Button>
                  </div>
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="border rounded-lg p-8 bg-gray-50 min-h-[400px]">
                  <div className="text-center text-gray-500">
                    <Layout className="h-12 w-12 mx-auto mb-4" />
                    <p>報表預覽將在執行後顯示</p>
                    <Button className="mt-4" onClick={handleExecute}>
                      立即執行
                    </Button>
                  </div>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="parameters" className="space-y-4">
            <Card>
              <CardHeader>
                <CardTitle>報表參數</CardTitle>
                <CardDescription>設定報表執行時的參數</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid gap-4">
                  <div className="grid gap-2">
                    <Label htmlFor="date_range">日期範圍</Label>
                    <div className="grid grid-cols-2 gap-2">
                      <Input type="date" placeholder="開始日期" />
                      <Input type="date" placeholder="結束日期" />
                    </div>
                  </div>
                  <div className="grid gap-2">
                    <Label htmlFor="company">公司</Label>
                    <Select>
                      <SelectTrigger>
                        <SelectValue placeholder="選擇公司" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="all">全部公司</SelectItem>
                        <SelectItem value="company1">公司 A</SelectItem>
                        <SelectItem value="company2">公司 B</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                  <div className="grid gap-2">
                    <Label htmlFor="department">部門</Label>
                    <Select>
                      <SelectTrigger>
                        <SelectValue placeholder="選擇部門" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="all">全部部門</SelectItem>
                        <SelectItem value="sales">業務部</SelectItem>
                        <SelectItem value="finance">財務部</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                </div>
                <Separator />
                <div className="flex justify-end gap-2">
                  <Button variant="outline">重設預設值</Button>
                  <Button>儲存參數</Button>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="executions" className="space-y-4">
            <Card>
              <CardHeader>
                <CardTitle>執行記錄</CardTitle>
                <CardDescription>最近 10 次執行記錄</CardDescription>
              </CardHeader>
              <CardContent>
                {executions?.data.length === 0 ? (
                  <div className="text-center py-8">
                    <Clock className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                    <p className="text-gray-500">暫無執行記錄</p>
                  </div>
                ) : (
                  <div className="space-y-4">
                    {executions?.data.map((execution) => (
                      <div
                        key={execution.id}
                        className="flex items-center justify-between p-4 border rounded-lg hover:bg-gray-50 cursor-pointer"
                        onClick={() => router.push(`/reports/executions/${execution.id}`)}
                      >
                        <div className="flex-1">
                          <div className="flex items-center gap-2">
                            <p className="font-medium">
                              執行 #{execution.execution_no}
                            </p>
                            {getExecutionStatusBadge(execution.status)}
                          </div>
                          <div className="flex items-center gap-4 mt-1 text-sm text-gray-500">
                            <span>執行者：{execution.executed_by_user?.full_name}</span>
                            <span>
                              開始時間：
                              {format(new Date(execution.started_at), 'MM/dd HH:mm', { locale: zhTW })}
                            </span>
                            {execution.execution_time > 0 && (
                              <span>耗時：{(execution.execution_time / 1000).toFixed(1)}s</span>
                            )}
                          </div>
                        </div>
                        <Button variant="ghost" size="sm">
                          <Eye className="h-4 w-4" />
                        </Button>
                      </div>
                    ))}
                  </div>
                )}
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="schedule" className="space-y-4">
            <Card>
              <CardHeader>
                <CardTitle>排程設定</CardTitle>
                <CardDescription>設定報表自動執行排程</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex items-center justify-between">
                  <div>
                    <h3 className="font-medium">啟用排程</h3>
                    <p className="text-sm text-gray-500">自動依照排程執行報表</p>
                  </div>
                  <input
                    type="checkbox"
                    checked={report.schedule_config?.enabled || false}
                    className="toggle"
                  />
                </div>
                <Separator />
                {report.schedule_config?.enabled && (
                  <>
                    <div className="grid gap-4">
                      <div className="grid gap-2">
                        <Label>執行頻率</Label>
                        <Select defaultValue={report.schedule_config?.frequency || 'daily'}>
                          <SelectTrigger>
                            <SelectValue />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectItem value="hourly">每小時</SelectItem>
                            <SelectItem value="daily">每日</SelectItem>
                            <SelectItem value="weekly">每週</SelectItem>
                            <SelectItem value="monthly">每月</SelectItem>
                          </SelectContent>
                        </Select>
                      </div>
                      <div className="grid gap-2">
                        <Label>執行時間</Label>
                        <Input
                          type="time"
                          defaultValue={report.schedule_config?.time || '08:00'}
                        />
                      </div>
                      <div className="grid gap-2">
                        <Label>收件人</Label>
                        <Input
                          placeholder="輸入電子郵件，多個收件人用逗號分隔"
                          defaultValue={report.schedule_config?.recipients?.join(', ') || ''}
                        />
                      </div>
                    </div>
                    <div className="flex justify-end">
                      <Button>儲存排程</Button>
                    </div>
                  </>
                )}
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="permissions" className="space-y-4">
            <Card>
              <CardHeader>
                <CardTitle>權限管理</CardTitle>
                <CardDescription>管理報表的檢視與編輯權限</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex items-center justify-between">
                  <div>
                    <h3 className="font-medium">公開報表</h3>
                    <p className="text-sm text-gray-500">允許所有使用者檢視此報表</p>
                  </div>
                  <input
                    type="checkbox"
                    checked={report.permissions?.is_public || false}
                    className="toggle"
                  />
                </div>
                <Separator />
                <div className="space-y-4">
                  <div>
                    <Label>檢視權限使用者</Label>
                    <p className="text-sm text-gray-500 mb-2">這些使用者可以檢視報表</p>
                    <Input placeholder="輸入使用者名稱或電子郵件" />
                    <div className="mt-2 space-y-2">
                      {report.permissions?.view_users?.map((user, index) => (
                        <div key={index} className="flex items-center justify-between p-2 border rounded">
                          <span>{user}</span>
                          <Button variant="ghost" size="sm">
                            <Trash2 className="h-4 w-4" />
                          </Button>
                        </div>
                      ))}
                    </div>
                  </div>
                  <div>
                    <Label>編輯權限使用者</Label>
                    <p className="text-sm text-gray-500 mb-2">這些使用者可以編輯報表</p>
                    <Input placeholder="輸入使用者名稱或電子郵件" />
                    <div className="mt-2 space-y-2">
                      {report.permissions?.edit_users?.map((user, index) => (
                        <div key={index} className="flex items-center justify-between p-2 border rounded">
                          <span>{user}</span>
                          <Button variant="ghost" size="sm">
                            <Trash2 className="h-4 w-4" />
                          </Button>
                        </div>
                      ))}
                    </div>
                  </div>
                </div>
                <div className="flex justify-end">
                  <Button>儲存權限設定</Button>
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </div>
    </DashboardLayout>
  )
}