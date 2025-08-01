'use client'

import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { useRouter } from 'next/navigation'
import DashboardLayout from '@/components/layout/DashboardLayout'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Badge } from '@/components/ui/badge'
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
  ArrowUpCircle,
  Eye,
  Download,
  Calendar,
  DollarSign,
  Users,
  AlertTriangle,
  Clock,
  CheckCircle,
  XCircle,
  Building2,
  FileText,
  Phone,
  Mail
} from 'lucide-react'
import financeService from '@/services/finance.service'
import Pagination from '@/components/common/Pagination'
import { format } from 'date-fns'
import { zhTW } from 'date-fns/locale'

export default function ReceivablesPage() {
  const router = useRouter()
  const [page, setPage] = useState(1)
  const [search, setSearch] = useState('')
  const [statusFilter, setStatusFilter] = useState<string>('')
  const [agingFilter, setAgingFilter] = useState<string>('')
  const [collectionFilter, setCollectionFilter] = useState<string>('')

  // Fetch receivables
  const { data: receivablesData, isLoading } = useQuery({
    queryKey: ['receivables', page, search, statusFilter, agingFilter, collectionFilter],
    queryFn: () => financeService.getAccountReceivables({
      page,
      page_size: 20,
      search,
      status: statusFilter,
      aging_category: agingFilter,
      collection_status: collectionFilter,
    }),
  })

  // Fetch AR summary
  const { data: arSummary } = useQuery({
    queryKey: ['ar-summary'],
    queryFn: () => financeService.getARSummary(),
  })

  const getStatusBadge = (status: string) => {
    const statusConfig: Record<string, { label: string; variant: any; icon: any }> = {
      open: { label: '未收', variant: 'warning', icon: Clock },
      partial: { label: '部分收款', variant: 'warning', icon: AlertTriangle },
      paid: { label: '已收款', variant: 'success', icon: CheckCircle },
      written_off: { label: '沖銷', variant: 'destructive', icon: XCircle },
    }

    const config = statusConfig[status] || { label: status, variant: 'default', icon: FileText }
    const Icon = config.icon
    
    return (
      <Badge variant={config.variant as any} className="flex items-center gap-1">
        <Icon className="h-3 w-3" />
        {config.label}
      </Badge>
    )
  }

  const getAgingBadge = (category: string, daysOverdue: number) => {
    const agingConfig: Record<string, { label: string; variant: any }> = {
      current: { label: '當期', variant: 'success' },
      '30days': { label: '30天內', variant: 'warning' },
      '60days': { label: '60天內', variant: 'warning' },
      '90days': { label: '90天內', variant: 'destructive' },
      over90days: { label: '超過90天', variant: 'destructive' },
    }

    const config = agingConfig[category] || { label: category, variant: 'default' }
    
    return (
      <div className="flex flex-col items-center">
        <Badge variant={config.variant as any}>{config.label}</Badge>
        {daysOverdue > 0 && (
          <span className="text-xs text-red-600 mt-1">逾期 {daysOverdue} 天</span>
        )}
      </div>
    )
  }

  const getCollectionPriorityBadge = (status: string) => {
    const priorityConfig: Record<string, { label: string; variant: any; icon: any }> = {
      normal: { label: '正常', variant: 'success', icon: CheckCircle },
      warning: { label: '注意', variant: 'warning', icon: AlertTriangle },
      critical: { label: '緊急', variant: 'destructive', icon: XCircle },
    }

    const config = priorityConfig[status] || { label: status, variant: 'default', icon: Clock }
    const Icon = config.icon
    
    return (
      <Badge variant={config.variant as any} className="flex items-center gap-1">
        <Icon className="h-3 w-3" />
        {config.label}
      </Badge>
    )
  }

  const handleSearch = (value: string) => {
    setSearch(value)
    setPage(1)
  }

  const handleContactCustomer = (customer: any, method: 'phone' | 'email') => {
    if (method === 'phone' && customer.phone) {
      window.open(`tel:${customer.phone}`)
    } else if (method === 'email' && customer.email) {
      window.open(`mailto:${customer.email}`)
    }
  }

  return (
    <DashboardLayout>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex justify-between items-center">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">應收帳款管理</h1>
            <p className="mt-2 text-gray-600">管理客戶欠款與收款進度</p>
          </div>
          <div className="flex gap-2">
            <Button variant="outline">
              <Download className="mr-2 h-4 w-4" />
              匯出帳齡表
            </Button>
          </div>
        </div>

        {/* AR Summary Cards */}
        {arSummary && (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-6">
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">總應收金額</CardTitle>
                <DollarSign className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">
                  {arSummary.currency} {arSummary.balance_amount.toLocaleString()}
                </div>
                <p className="text-xs text-muted-foreground">
                  {arSummary.open_items} 張未收發票
                </p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">當期</CardTitle>
                <CheckCircle className="h-4 w-4 text-green-600" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold text-green-600">
                  {arSummary.currency} {arSummary.current.toLocaleString()}
                </div>
                <p className="text-xs text-muted-foreground">未逾期金額</p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">30天內</CardTitle>
                <Clock className="h-4 w-4 text-yellow-600" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold text-yellow-600">
                  {arSummary.currency} {arSummary.days_30.toLocaleString()}
                </div>
                <p className="text-xs text-muted-foreground">30天內逾期</p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">60天內</CardTitle>
                <AlertTriangle className="h-4 w-4 text-orange-600" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold text-orange-600">
                  {arSummary.currency} {arSummary.days_60.toLocaleString()}
                </div>
                <p className="text-xs text-muted-foreground">60天內逾期</p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">90天以上</CardTitle>
                <XCircle className="h-4 w-4 text-red-600" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold text-red-600">
                  {arSummary.currency} {arSummary.over_90.toLocaleString()}
                </div>
                <p className="text-xs text-muted-foreground">高風險欠款</p>
              </CardContent>
            </Card>
          </div>
        )}

        {/* Filters */}
        <Card>
          <CardHeader>
            <CardTitle>應收帳款明細</CardTitle>
            <CardDescription>查看所有客戶的應收帳款詳情</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="flex flex-col sm:flex-row gap-4 mb-6">
              <div className="relative flex-1">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
                <Input
                  placeholder="搜尋客戶名稱或發票編號..."
                  value={search}
                  onChange={(e) => handleSearch(e.target.value)}
                  className="pl-10"
                />
              </div>
              <Select value={statusFilter} onValueChange={(value) => { setStatusFilter(value); setPage(1); }}>
                <SelectTrigger className="w-full sm:w-[150px]">
                  <SelectValue placeholder="收款狀態" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="">全部狀態</SelectItem>
                  <SelectItem value="open">未收</SelectItem>
                  <SelectItem value="partial">部分收款</SelectItem>
                  <SelectItem value="paid">已收款</SelectItem>
                  <SelectItem value="written_off">沖銷</SelectItem>
                </SelectContent>
              </Select>
              <Select value={agingFilter} onValueChange={(value) => { setAgingFilter(value); setPage(1); }}>
                <SelectTrigger className="w-full sm:w-[150px]">
                  <SelectValue placeholder="帳齡分析" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="">全部帳齡</SelectItem>
                  <SelectItem value="current">當期</SelectItem>
                  <SelectItem value="30days">30天內</SelectItem>
                  <SelectItem value="60days">60天內</SelectItem>
                  <SelectItem value="90days">90天內</SelectItem>
                  <SelectItem value="over90days">超過90天</SelectItem>
                </SelectContent>
              </Select>
              <Select value={collectionFilter} onValueChange={(value) => { setCollectionFilter(value); setPage(1); }}>
                <SelectTrigger className="w-full sm:w-[150px]">
                  <SelectValue placeholder="催收優先級" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="">全部優先級</SelectItem>
                  <SelectItem value="normal">正常</SelectItem>
                  <SelectItem value="warning">注意</SelectItem>
                  <SelectItem value="critical">緊急</SelectItem>
                </SelectContent>
              </Select>
            </div>

            {/* Receivables Table */}
            <div className="rounded-md border">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>客戶</TableHead>
                    <TableHead>發票編號</TableHead>
                    <TableHead>發票日期</TableHead>
                    <TableHead>到期日期</TableHead>
                    <TableHead className="text-right">發票金額</TableHead>
                    <TableHead className="text-right">已收金額</TableHead>
                    <TableHead className="text-right">未收金額</TableHead>
                    <TableHead>帳齡</TableHead>
                    <TableHead>收款狀態</TableHead>
                    <TableHead>催收優先級</TableHead>
                    <TableHead>操作</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {isLoading ? (
                    <TableRow>
                      <TableCell colSpan={11} className="text-center py-8">
                        載入中...
                      </TableCell>
                    </TableRow>
                  ) : receivablesData?.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={11} className="text-center py-8 text-gray-500">
                        沒有找到應收帳款記錄
                      </TableCell>
                    </TableRow>
                  ) : (
                    receivablesData?.map((receivable) => (
                      <TableRow key={receivable.id} className="hover:bg-gray-50">
                        <TableCell>
                          <div className="flex items-center gap-2">
                            <Users className="h-4 w-4 text-blue-500" />
                            <div>
                              <p className="font-medium">{receivable.customer?.name}</p>
                              <p className="text-xs text-gray-500">{receivable.customer?.customer_code}</p>
                            </div>
                          </div>
                        </TableCell>
                        <TableCell>
                          <Button
                            variant="link"
                            size="sm"
                            className="p-0 h-auto font-medium"
                            onClick={() => router.push(`/finance/invoices/${receivable.invoice?.id}`)}
                          >
                            {receivable.invoice?.invoice_no}
                          </Button>
                        </TableCell>
                        <TableCell>
                          <div className="flex items-center gap-1">
                            <Calendar className="h-3 w-3 text-gray-400" />
                            <span className="text-sm">
                              {format(new Date(receivable.invoice_date), 'yyyy/MM/dd', { locale: zhTW })}
                            </span>
                          </div>
                        </TableCell>
                        <TableCell>
                          <div className="flex items-center gap-1">
                            <Calendar className="h-3 w-3 text-gray-400" />
                            <span className="text-sm">
                              {format(new Date(receivable.due_date), 'yyyy/MM/dd', { locale: zhTW })}
                            </span>
                            {receivable.days_overdue > 0 && (
                              <span className="text-xs text-red-600 ml-1">
                                (逾期 {receivable.days_overdue} 天)
                              </span>
                            )}
                          </div>
                        </TableCell>
                        <TableCell className="text-right">
                          <div className="flex items-center justify-end gap-1">
                            <span className="text-xs text-gray-500">{receivable.currency}</span>
                            <span className="font-medium">{receivable.invoice_amount.toLocaleString()}</span>
                          </div>
                        </TableCell>
                        <TableCell className="text-right">
                          {receivable.paid_amount > 0 ? (
                            <div className="flex items-center justify-end gap-1">
                              <span className="text-xs text-gray-500">{receivable.currency}</span>
                              <span className="font-medium text-green-600">
                                {receivable.paid_amount.toLocaleString()}
                              </span>
                            </div>
                          ) : (
                            <span className="text-gray-400">-</span>
                          )}
                        </TableCell>
                        <TableCell className="text-right">
                          <div className="flex items-center justify-end gap-1">
                            <span className="text-xs text-gray-500">{receivable.currency}</span>
                            <span className={`font-medium ${
                              receivable.balance_amount > 0 ? 'text-red-600' : 'text-green-600'
                            }`}>
                              {receivable.balance_amount.toLocaleString()}
                            </span>
                          </div>
                        </TableCell>
                        <TableCell>
                          {getAgingBadge(receivable.aging_category, receivable.days_overdue)}
                        </TableCell>
                        <TableCell>{getStatusBadge(receivable.status)}</TableCell>
                        <TableCell>{getCollectionPriorityBadge(receivable.collection_status)}</TableCell>
                        <TableCell>
                          <div className="flex items-center gap-1">
                            <Button
                              variant="ghost"
                              size="icon"
                              onClick={() => router.push(`/finance/invoices/${receivable.invoice?.id}`)}
                            >
                              <Eye className="h-4 w-4" />
                            </Button>
                            {receivable.customer?.phone && (
                              <Button
                                variant="ghost"
                                size="icon"
                                onClick={() => handleContactCustomer(receivable.customer, 'phone')}
                                title="撥打電話"
                              >
                                <Phone className="h-4 w-4" />
                              </Button>
                            )}
                            {receivable.customer?.email && (
                              <Button
                                variant="ghost"
                                size="icon"
                                onClick={() => handleContactCustomer(receivable.customer, 'email')}
                                title="發送郵件"
                              >
                                <Mail className="h-4 w-4" />
                              </Button>
                            )}
                          </div>
                        </TableCell>
                      </TableRow>
                    ))
                  )}
                </TableBody>
              </Table>
            </div>

            {/* Pagination */}
            {receivablesData && receivablesData.length > 0 && (
              <div className="mt-4">
                <Pagination
                  currentPage={page}
                  totalPages={Math.ceil(receivablesData.length / 20)}
                  onPageChange={setPage}
                />
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </DashboardLayout>
  )
}