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
  CreditCard,
  Plus,
  Eye,
  Download,
  Calendar,
  DollarSign,
  TrendingUp,
  TrendingDown,
  Clock,
  CheckCircle,
  XCircle,
  AlertCircle,
  Building2,
  Users
} from 'lucide-react'
import financeService from '@/services/finance.service'
import Pagination from '@/components/common/Pagination'
import { format } from 'date-fns'
import { zhTW } from 'date-fns/locale'

export default function PaymentsPage() {
  const router = useRouter()
  const [page, setPage] = useState(1)
  const [search, setSearch] = useState('')
  const [typeFilter, setTypeFilter] = useState<string>('')
  const [statusFilter, setStatusFilter] = useState<string>('')
  const [methodFilter, setMethodFilter] = useState<string>('')

  // Fetch payments
  const { data: paymentsData, isLoading } = useQuery({
    queryKey: ['payments', page, search, typeFilter, statusFilter, methodFilter],
    queryFn: () => financeService.listPayments({
      page,
      page_size: 20,
      search,
      type: typeFilter,
      status: statusFilter,
      payment_method: methodFilter,
    }),
  })

  // Calculate statistics
  const stats = paymentsData?.data.reduce((acc, payment) => {
    acc.total++
    if (payment.type === 'incoming') {
      acc.totalIncoming += payment.amount
      acc.countIncoming++
    } else {
      acc.totalOutgoing += payment.amount
      acc.countOutgoing++
    }
    
    if (payment.status === 'completed') {
      acc.completed++
    } else if (payment.status === 'pending') {
      acc.pending++
    }
    
    return acc
  }, {
    total: 0,
    totalIncoming: 0,
    totalOutgoing: 0,
    countIncoming: 0,
    countOutgoing: 0,
    completed: 0,
    pending: 0,
  })

  const getPaymentStatusBadge = (status: string) => {
    const statusConfig: Record<string, { label: string; variant: any; icon: any }> = {
      pending: { label: '處理中', variant: 'warning', icon: Clock },
      completed: { label: '已完成', variant: 'success', icon: CheckCircle },
      failed: { label: '失敗', variant: 'destructive', icon: XCircle },
      cancelled: { label: '已取消', variant: 'secondary', icon: XCircle },
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

  const getPaymentTypeBadge = (type: string) => {
    const typeConfig: Record<string, { label: string; variant: any; icon: any }> = {
      incoming: { label: '收款', variant: 'success', icon: TrendingUp },
      outgoing: { label: '付款', variant: 'info', icon: TrendingDown },
    }

    const config = typeConfig[type] || { label: type, variant: 'default', icon: CreditCard }
    const Icon = config.icon
    
    return (
      <Badge variant={config.variant as any} className="flex items-center gap-1">
        <Icon className="h-3 w-3" />
        {config.label}
      </Badge>
    )
  }

  const getPaymentMethodLabel = (method: string) => {
    const methodLabels: Record<string, string> = {
      cash: '現金',
      check: '支票',
      bank_transfer: '銀行轉帳',
      credit_card: '信用卡',
    }
    return methodLabels[method] || method
  }

  const handleSearch = (value: string) => {
    setSearch(value)
    setPage(1)
  }

  return (
    <DashboardLayout>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex justify-between items-center">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">付款管理</h1>
            <p className="mt-2 text-gray-600">管理所有收款與付款記錄</p>
          </div>
          <div className="flex gap-2">
            <Button variant="outline">
              <Download className="mr-2 h-4 w-4" />
              匯出資料
            </Button>
            <Button onClick={() => router.push('/finance/payments/new')}>
              <Plus className="mr-2 h-4 w-4" />
              新增付款
            </Button>
          </div>
        </div>

        {/* Statistics Cards */}
        {stats && (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">總收款金額</CardTitle>
                <TrendingUp className="h-4 w-4 text-green-600" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold text-green-600">
                  ${stats.totalIncoming.toLocaleString()}
                </div>
                <p className="text-xs text-muted-foreground">
                  {stats.countIncoming} 筆收款
                </p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">總付款金額</CardTitle>
                <TrendingDown className="h-4 w-4 text-red-600" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold text-red-600">
                  ${stats.totalOutgoing.toLocaleString()}
                </div>
                <p className="text-xs text-muted-foreground">
                  {stats.countOutgoing} 筆付款
                </p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">已完成</CardTitle>
                <CheckCircle className="h-4 w-4 text-green-600" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{stats.completed}</div>
                <p className="text-xs text-muted-foreground">
                  成功完成的付款
                </p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">處理中</CardTitle>
                <Clock className="h-4 w-4 text-yellow-600" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold text-yellow-600">{stats.pending}</div>
                <p className="text-xs text-muted-foreground">
                  等待處理的付款
                </p>
              </CardContent>
            </Card>
          </div>
        )}

        {/* Filters */}
        <Card>
          <CardHeader>
            <CardTitle>付款列表</CardTitle>
            <CardDescription>查看和管理所有付款記錄</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="flex flex-col sm:flex-row gap-4 mb-6">
              <div className="relative flex-1">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
                <Input
                  placeholder="搜尋付款編號、客戶或供應商..."
                  value={search}
                  onChange={(e) => handleSearch(e.target.value)}
                  className="pl-10"
                />
              </div>
              <Select value={typeFilter} onValueChange={(value) => { setTypeFilter(value); setPage(1); }}>
                <SelectTrigger className="w-full sm:w-[150px]">
                  <SelectValue placeholder="類型" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="">全部類型</SelectItem>
                  <SelectItem value="incoming">收款</SelectItem>
                  <SelectItem value="outgoing">付款</SelectItem>
                </SelectContent>
              </Select>
              <Select value={statusFilter} onValueChange={(value) => { setStatusFilter(value); setPage(1); }}>
                <SelectTrigger className="w-full sm:w-[150px]">
                  <SelectValue placeholder="狀態" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="">全部狀態</SelectItem>
                  <SelectItem value="pending">處理中</SelectItem>
                  <SelectItem value="completed">已完成</SelectItem>
                  <SelectItem value="failed">失敗</SelectItem>
                  <SelectItem value="cancelled">已取消</SelectItem>
                </SelectContent>
              </Select>
              <Select value={methodFilter} onValueChange={(value) => { setMethodFilter(value); setPage(1); }}>
                <SelectTrigger className="w-full sm:w-[150px]">
                  <SelectValue placeholder="付款方式" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="">全部方式</SelectItem>
                  <SelectItem value="cash">現金</SelectItem>
                  <SelectItem value="check">支票</SelectItem>
                  <SelectItem value="bank_transfer">銀行轉帳</SelectItem>
                  <SelectItem value="credit_card">信用卡</SelectItem>
                </SelectContent>
              </Select>
            </div>

            {/* Payments Table */}
            <div className="rounded-md border">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>付款編號</TableHead>
                    <TableHead>類型</TableHead>
                    <TableHead>對象</TableHead>
                    <TableHead>付款日期</TableHead>
                    <TableHead className="text-right">金額</TableHead>
                    <TableHead>付款方式</TableHead>
                    <TableHead>狀態</TableHead>
                    <TableHead>發票編號</TableHead>
                    <TableHead>操作</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {isLoading ? (
                    <TableRow>
                      <TableCell colSpan={9} className="text-center py-8">
                        載入中...
                      </TableCell>
                    </TableRow>
                  ) : paymentsData?.data.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={9} className="text-center py-8 text-gray-500">
                        沒有找到付款記錄
                      </TableCell>
                    </TableRow>
                  ) : (
                    paymentsData?.data.map((payment) => (
                      <TableRow key={payment.id} className="hover:bg-gray-50">
                        <TableCell className="font-medium">{payment.payment_no}</TableCell>
                        <TableCell>{getPaymentTypeBadge(payment.type)}</TableCell>
                        <TableCell>
                          <div className="flex items-center gap-2">
                            {payment.type === 'incoming' ? (
                              <Users className="h-4 w-4 text-green-500" />
                            ) : (
                              <Building2 className="h-4 w-4 text-blue-500" />
                            )}
                            <span className="truncate">
                              {payment.customer?.name || payment.supplier?.name || '-'}
                            </span>
                          </div>
                        </TableCell>
                        <TableCell>
                          <div className="flex items-center gap-1">
                            <Calendar className="h-3 w-3 text-gray-400" />
                            <span className="text-sm">
                              {format(new Date(payment.payment_date), 'yyyy/MM/dd', { locale: zhTW })}
                            </span>
                          </div>
                        </TableCell>
                        <TableCell className="text-right">
                          <div className="flex items-center justify-end gap-1">
                            <span className="text-xs text-gray-500">{payment.currency}</span>
                            <span className={`font-medium ${
                              payment.type === 'incoming' ? 'text-green-600' : 'text-red-600'
                            }`}>
                              {payment.amount.toLocaleString()}
                            </span>
                          </div>
                        </TableCell>
                        <TableCell>{getPaymentMethodLabel(payment.payment_method)}</TableCell>
                        <TableCell>{getPaymentStatusBadge(payment.status)}</TableCell>
                        <TableCell>
                          {payment.invoice ? (
                            <Button
                              variant="link"
                              size="sm"
                              className="p-0 h-auto"
                              onClick={() => router.push(`/finance/invoices/${payment.invoice?.id}`)}
                            >
                              {payment.invoice.invoice_no}
                            </Button>
                          ) : (
                            <span className="text-gray-400">-</span>
                          )}
                        </TableCell>
                        <TableCell>
                          <div className="flex items-center gap-1">
                            <Button
                              variant="ghost"
                              size="icon"
                              onClick={() => router.push(`/finance/payments/${payment.id}`)}
                            >
                              <Eye className="h-4 w-4" />
                            </Button>
                          </div>
                        </TableCell>
                      </TableRow>
                    ))
                  )}
                </TableBody>
              </Table>
            </div>

            {/* Pagination */}
            {paymentsData && paymentsData.pagination.total > 0 && (
              <div className="mt-4">
                <Pagination
                  currentPage={page}
                  totalPages={paymentsData.pagination.total_pages}
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