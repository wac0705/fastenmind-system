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
  Filter,
  Package,
  Calendar,
  DollarSign,
  TrendingUp,
  Clock,
  CheckCircle,
  XCircle,
  Truck,
  Factory,
  AlertCircle,
  Eye,
  Edit,
  Plus,
  Download,
  MoreHorizontal,
  FileText
} from 'lucide-react'
import { format } from 'date-fns'
import { zhTW } from 'date-fns/locale'
import orderService, { Order, OrderStats } from '@/services/order.service'
import { customerService } from '@/services/customer.service'
import Pagination from '@/components/common/Pagination'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { useAuthStore } from '@/store/auth.store'

export default function OrdersPage() {
  const router = useRouter()
  const { user } = useAuthStore()
  const [page, setPage] = useState(1)
  const [search, setSearch] = useState('')
  const [statusFilter, setStatusFilter] = useState<string>('')
  const [customerFilter, setCustomerFilter] = useState<string>('')
  const [paymentStatusFilter, setPaymentStatusFilter] = useState<string>('')
  const [dateRange, setDateRange] = useState('')

  // Fetch orders
  const { data: ordersData, isLoading: isLoadingOrders } = useQuery({
    queryKey: ['orders', page, search, statusFilter, customerFilter, paymentStatusFilter],
    queryFn: () => orderService.list({
      page,
      page_size: 20,
      search,
      status: statusFilter,
      customer_id: customerFilter,
      payment_status: paymentStatusFilter,
    }),
  })

  // Fetch stats
  const { data: stats } = useQuery({
    queryKey: ['order-stats'],
    queryFn: () => orderService.getStats(),
  })

  // Fetch customers for filter
  const { data: customersData } = useQuery({
    queryKey: ['customers-filter'],
    queryFn: () => customerService.list({ page: 1, page_size: 100 }),
  })

  const getStatusBadge = (status: string) => {
    const statusConfig: Record<string, { label: string; variant: any; icon: any }> = {
      pending: { label: '待確認', variant: 'secondary', icon: Clock },
      confirmed: { label: '已確認', variant: 'info', icon: CheckCircle },
      in_production: { label: '生產中', variant: 'warning', icon: Factory },
      quality_check: { label: '品檢中', variant: 'warning', icon: AlertCircle },
      ready_to_ship: { label: '待出貨', variant: 'info', icon: Package },
      shipped: { label: '已出貨', variant: 'info', icon: Truck },
      delivered: { label: '已送達', variant: 'success', icon: CheckCircle },
      completed: { label: '已完成', variant: 'success', icon: CheckCircle },
      cancelled: { label: '已取消', variant: 'destructive', icon: XCircle },
    }

    const config = statusConfig[status] || { label: status, variant: 'default', icon: Package }
    const Icon = config.icon
    
    return (
      <Badge variant={config.variant as any} className="flex items-center gap-1">
        <Icon className="h-3 w-3" />
        {config.label}
      </Badge>
    )
  }

  const getPaymentStatusBadge = (status: string) => {
    const statusConfig: Record<string, { label: string; variant: any }> = {
      pending: { label: '未付款', variant: 'secondary' },
      partial: { label: '部分付款', variant: 'warning' },
      paid: { label: '已付款', variant: 'success' },
    }

    const config = statusConfig[status] || { label: status, variant: 'default' }
    
    return <Badge variant={config.variant as any}>{config.label}</Badge>
  }

  const handleSearch = (value: string) => {
    setSearch(value)
    setPage(1)
  }

  const handleExport = () => {
    // TODO: Implement export functionality
    console.log('Export orders')
  }

  const canEdit = (order: Order) => {
    return ['admin', 'manager', 'sales'].includes(user?.role || '') && 
           ['pending', 'confirmed'].includes(order.status)
  }

  const canDelete = (order: Order) => {
    return ['admin', 'manager'].includes(user?.role || '') && 
           order.status === 'pending'
  }

  return (
    <DashboardLayout>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex justify-between items-start">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">訂單管理</h1>
            <p className="mt-2 text-gray-600">管理所有客戶訂單與出貨狀態</p>
          </div>
          <div className="flex gap-3">
            <Button variant="outline" onClick={handleExport}>
              <Download className="mr-2 h-4 w-4" />
              匯出資料
            </Button>
            <Button onClick={() => router.push('/orders/new')}>
              <Plus className="mr-2 h-4 w-4" />
              新增訂單
            </Button>
          </div>
        </div>

        {/* Stats Cards */}
        {stats && (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">總訂單數</CardTitle>
                <Package className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{stats.total_orders}</div>
                <p className="text-xs text-muted-foreground">所有訂單</p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">生產中</CardTitle>
                <Factory className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{stats.in_production}</div>
                <p className="text-xs text-muted-foreground">正在生產的訂單</p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">總營收</CardTitle>
                <DollarSign className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">${stats.total_revenue.toLocaleString()}</div>
                <p className="text-xs text-muted-foreground">所有訂單金額</p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">平均訂單金額</CardTitle>
                <TrendingUp className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">${stats.avg_order_value.toLocaleString()}</div>
                <p className="text-xs text-muted-foreground">每筆訂單平均值</p>
              </CardContent>
            </Card>
          </div>
        )}

        {/* Filters */}
        <Card>
          <CardHeader>
            <CardTitle>訂單列表</CardTitle>
            <CardDescription>查看和管理所有訂單</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="flex flex-col sm:flex-row gap-4 mb-6">
              <div className="relative flex-1">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
                <Input
                  placeholder="搜尋訂單號或 PO 號碼..."
                  value={search}
                  onChange={(e) => handleSearch(e.target.value)}
                  className="pl-10"
                />
              </div>
              <Select value={statusFilter} onValueChange={(value) => { setStatusFilter(value); setPage(1); }}>
                <SelectTrigger className="w-full sm:w-[180px]">
                  <SelectValue placeholder="訂單狀態" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="">全部狀態</SelectItem>
                  <SelectItem value="pending">待確認</SelectItem>
                  <SelectItem value="confirmed">已確認</SelectItem>
                  <SelectItem value="in_production">生產中</SelectItem>
                  <SelectItem value="quality_check">品檢中</SelectItem>
                  <SelectItem value="ready_to_ship">待出貨</SelectItem>
                  <SelectItem value="shipped">已出貨</SelectItem>
                  <SelectItem value="delivered">已送達</SelectItem>
                  <SelectItem value="completed">已完成</SelectItem>
                  <SelectItem value="cancelled">已取消</SelectItem>
                </SelectContent>
              </Select>
              <Select value={customerFilter} onValueChange={(value) => { setCustomerFilter(value); setPage(1); }}>
                <SelectTrigger className="w-full sm:w-[200px]">
                  <SelectValue placeholder="選擇客戶" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="">全部客戶</SelectItem>
                  {customersData?.data.map((customer) => (
                    <SelectItem key={customer.id} value={customer.id}>
                      {customer.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              <Select value={paymentStatusFilter} onValueChange={(value) => { setPaymentStatusFilter(value); setPage(1); }}>
                <SelectTrigger className="w-full sm:w-[180px]">
                  <SelectValue placeholder="付款狀態" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="">全部狀態</SelectItem>
                  <SelectItem value="pending">未付款</SelectItem>
                  <SelectItem value="partial">部分付款</SelectItem>
                  <SelectItem value="paid">已付款</SelectItem>
                </SelectContent>
              </Select>
              <Select value={dateRange} onValueChange={(value) => { setDateRange(value); setPage(1); }}>
                <SelectTrigger className="w-full sm:w-[150px]">
                  <SelectValue placeholder="時間範圍" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="">全部時間</SelectItem>
                  <SelectItem value="today">今天</SelectItem>
                  <SelectItem value="week">本週</SelectItem>
                  <SelectItem value="month">本月</SelectItem>
                  <SelectItem value="quarter">本季</SelectItem>
                  <SelectItem value="year">今年</SelectItem>
                </SelectContent>
              </Select>
            </div>

            {/* Orders Table */}
            <div className="rounded-md border">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>訂單號</TableHead>
                    <TableHead>PO 號碼</TableHead>
                    <TableHead>客戶</TableHead>
                    <TableHead>數量</TableHead>
                    <TableHead>金額</TableHead>
                    <TableHead>狀態</TableHead>
                    <TableHead>付款狀態</TableHead>
                    <TableHead>交貨日期</TableHead>
                    <TableHead>建立日期</TableHead>
                    <TableHead className="text-right">操作</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {isLoadingOrders ? (
                    <TableRow>
                      <TableCell colSpan={10} className="text-center py-8">
                        載入中...
                      </TableCell>
                    </TableRow>
                  ) : ordersData?.data.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={10} className="text-center py-8 text-gray-500">
                        沒有找到訂單
                      </TableCell>
                    </TableRow>
                  ) : (
                    ordersData?.data.map((order) => (
                      <TableRow key={order.id} className="hover:bg-gray-50">
                        <TableCell className="font-medium">{order.order_no}</TableCell>
                        <TableCell>{order.po_number}</TableCell>
                        <TableCell>{order.customer?.name || '-'}</TableCell>
                        <TableCell>{order.quantity.toLocaleString()}</TableCell>
                        <TableCell>
                          <div className="flex items-center gap-1">
                            <span className="text-xs text-gray-500">{order.currency}</span>
                            <span className="font-medium">{order.total_amount.toLocaleString()}</span>
                          </div>
                        </TableCell>
                        <TableCell>{getStatusBadge(order.status)}</TableCell>
                        <TableCell>{getPaymentStatusBadge(order.payment_status)}</TableCell>
                        <TableCell>
                          <div className="flex items-center gap-1">
                            <Calendar className="h-3 w-3 text-gray-400" />
                            <span className="text-sm">
                              {format(new Date(order.delivery_date), 'yyyy/MM/dd', { locale: zhTW })}
                            </span>
                          </div>
                        </TableCell>
                        <TableCell className="text-sm text-gray-500">
                          {format(new Date(order.created_at), 'yyyy/MM/dd', { locale: zhTW })}
                        </TableCell>
                        <TableCell className="text-right">
                          <DropdownMenu>
                            <DropdownMenuTrigger asChild>
                              <Button variant="ghost" size="icon">
                                <MoreHorizontal className="h-4 w-4" />
                              </Button>
                            </DropdownMenuTrigger>
                            <DropdownMenuContent align="end">
                              <DropdownMenuLabel>操作</DropdownMenuLabel>
                              <DropdownMenuItem onClick={() => router.push(`/orders/${order.id}`)}>
                                <Eye className="mr-2 h-4 w-4" />
                                檢視詳情
                              </DropdownMenuItem>
                              {canEdit(order) && (
                                <DropdownMenuItem onClick={() => router.push(`/orders/${order.id}/edit`)}>
                                  <Edit className="mr-2 h-4 w-4" />
                                  編輯訂單
                                </DropdownMenuItem>
                              )}
                              <DropdownMenuItem onClick={() => router.push(`/quotes/${order.quote_id}`)}>
                                <FileText className="mr-2 h-4 w-4" />
                                查看報價單
                              </DropdownMenuItem>
                              <DropdownMenuSeparator />
                              <DropdownMenuItem>
                                <Download className="mr-2 h-4 w-4" />
                                下載 PDF
                              </DropdownMenuItem>
                            </DropdownMenuContent>
                          </DropdownMenu>
                        </TableCell>
                      </TableRow>
                    ))
                  )}
                </TableBody>
              </Table>
            </div>

            {/* Pagination */}
            {ordersData && ordersData.pagination.total > 0 && (
              <div className="mt-4">
                <Pagination
                  currentPage={page}
                  totalPages={Math.ceil(ordersData.pagination.total / 20)}
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