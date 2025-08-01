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
  FileText,
  Plus,
  Eye,
  Edit,
  Download,
  Calendar,
  DollarSign,
  Users,
  Building2
} from 'lucide-react'
import financeService from '@/services/finance.service'
import Pagination from '@/components/common/Pagination'
import { format } from 'date-fns'
import { zhTW } from 'date-fns/locale'

export default function InvoicesPage() {
  const router = useRouter()
  const [page, setPage] = useState(1)
  const [search, setSearch] = useState('')
  const [typeFilter, setTypeFilter] = useState<string>('')
  const [statusFilter, setStatusFilter] = useState<string>('')

  // Fetch invoices
  const { data: invoicesData, isLoading } = useQuery({
    queryKey: ['invoices', page, search, typeFilter, statusFilter],
    queryFn: () => financeService.listInvoices({
      page,
      page_size: 20,
      search,
      type: typeFilter,
      status: statusFilter,
    }),
  })

  const getInvoiceStatusBadge = (status: string) => {
    const statusConfig: Record<string, { label: string; variant: any }> = {
      draft: { label: '草稿', variant: 'secondary' },
      issued: { label: '已開立', variant: 'info' },
      sent: { label: '已寄送', variant: 'warning' },
      paid: { label: '已付款', variant: 'success' },
      partial_paid: { label: '部分付款', variant: 'warning' },
      overdue: { label: '逾期', variant: 'destructive' },
      cancelled: { label: '已取消', variant: 'secondary' },
    }

    const config = statusConfig[status] || { label: status, variant: 'default' }
    return <Badge variant={config.variant as any}>{config.label}</Badge>
  }

  const getInvoiceTypeBadge = (type: string) => {
    const typeConfig: Record<string, { label: string; variant: any }> = {
      sales: { label: '銷售', variant: 'success' },
      purchase: { label: '採購', variant: 'info' },
      credit_note: { label: '貸項憑單', variant: 'warning' },
      debit_note: { label: '借項憑單', variant: 'destructive' },
    }

    const config = typeConfig[type] || { label: type, variant: 'default' }
    return <Badge variant={config.variant as any}>{config.label}</Badge>
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
            <h1 className="text-3xl font-bold text-gray-900">發票管理</h1>
            <p className="mt-2 text-gray-600">管理所有銷售與採購發票</p>
          </div>
          <div className="flex gap-2">
            <Button variant="outline">
              <Download className="mr-2 h-4 w-4" />
              匯出
            </Button>
            <Button onClick={() => router.push('/finance/invoices/new')}>
              <Plus className="mr-2 h-4 w-4" />
              新增發票
            </Button>
          </div>
        </div>

        {/* Filters */}
        <Card>
          <CardHeader>
            <CardTitle>發票列表</CardTitle>
            <CardDescription>查看和管理所有發票</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="flex flex-col sm:flex-row gap-4 mb-6">
              <div className="relative flex-1">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
                <Input
                  placeholder="搜尋發票編號、客戶或供應商..."
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
                  <SelectItem value="sales">銷售</SelectItem>
                  <SelectItem value="purchase">採購</SelectItem>
                  <SelectItem value="credit_note">貸項憑單</SelectItem>
                  <SelectItem value="debit_note">借項憑單</SelectItem>
                </SelectContent>
              </Select>
              <Select value={statusFilter} onValueChange={(value) => { setStatusFilter(value); setPage(1); }}>
                <SelectTrigger className="w-full sm:w-[150px]">
                  <SelectValue placeholder="狀態" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="">全部狀態</SelectItem>
                  <SelectItem value="draft">草稿</SelectItem>
                  <SelectItem value="issued">已開立</SelectItem>
                  <SelectItem value="sent">已寄送</SelectItem>
                  <SelectItem value="paid">已付款</SelectItem>
                  <SelectItem value="partial_paid">部分付款</SelectItem>
                  <SelectItem value="overdue">逾期</SelectItem>
                  <SelectItem value="cancelled">已取消</SelectItem>
                </SelectContent>
              </Select>
            </div>

            {/* Invoices Table */}
            <div className="rounded-md border">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>發票編號</TableHead>
                    <TableHead>類型</TableHead>
                    <TableHead>客戶/供應商</TableHead>
                    <TableHead>開立日期</TableHead>
                    <TableHead>到期日期</TableHead>
                    <TableHead className="text-right">金額</TableHead>
                    <TableHead className="text-right">已付金額</TableHead>
                    <TableHead>狀態</TableHead>
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
                  ) : invoicesData?.data.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={9} className="text-center py-8 text-gray-500">
                        沒有找到發票記錄
                      </TableCell>
                    </TableRow>
                  ) : (
                    invoicesData?.data.map((invoice) => (
                      <TableRow key={invoice.id} className="hover:bg-gray-50">
                        <TableCell className="font-medium">{invoice.invoice_no}</TableCell>
                        <TableCell>{getInvoiceTypeBadge(invoice.type)}</TableCell>
                        <TableCell>
                          <div className="flex items-center gap-2">
                            {invoice.type === 'sales' ? (
                              <Users className="h-4 w-4 text-green-500" />
                            ) : (
                              <Building2 className="h-4 w-4 text-blue-500" />
                            )}
                            <span className="truncate">
                              {invoice.customer?.name || invoice.supplier?.name || '-'}
                            </span>
                          </div>
                        </TableCell>
                        <TableCell>
                          <div className="flex items-center gap-1">
                            <Calendar className="h-3 w-3 text-gray-400" />
                            <span className="text-sm">
                              {format(new Date(invoice.issue_date), 'yyyy/MM/dd', { locale: zhTW })}
                            </span>
                          </div>
                        </TableCell>
                        <TableCell>
                          <div className="flex items-center gap-1">
                            <Calendar className="h-3 w-3 text-gray-400" />
                            <span className="text-sm">
                              {format(new Date(invoice.due_date), 'yyyy/MM/dd', { locale: zhTW })}
                            </span>
                          </div>
                        </TableCell>
                        <TableCell className="text-right">
                          <div className="flex items-center justify-end gap-1">
                            <span className="text-xs text-gray-500">{invoice.currency}</span>
                            <span className="font-medium">{invoice.total_amount.toLocaleString()}</span>
                          </div>
                        </TableCell>
                        <TableCell className="text-right">
                          {invoice.paid_amount > 0 ? (
                            <div className="flex items-center justify-end gap-1">
                              <span className="text-xs text-gray-500">{invoice.currency}</span>
                              <span className="font-medium text-green-600">
                                {invoice.paid_amount.toLocaleString()}
                              </span>
                            </div>
                          ) : (
                            <span className="text-gray-400">-</span>
                          )}
                        </TableCell>
                        <TableCell>{getInvoiceStatusBadge(invoice.status)}</TableCell>
                        <TableCell>
                          <div className="flex items-center gap-1">
                            <Button
                              variant="ghost"
                              size="icon"
                              onClick={() => router.push(`/finance/invoices/${invoice.id}`)}
                            >
                              <Eye className="h-4 w-4" />
                            </Button>
                            <Button
                              variant="ghost"
                              size="icon"
                              onClick={() => router.push(`/finance/invoices/${invoice.id}/edit`)}
                            >
                              <Edit className="h-4 w-4" />
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
            {invoicesData && invoicesData.pagination.total > 0 && (
              <div className="mt-4">
                <Pagination
                  currentPage={page}
                  totalPages={invoicesData.pagination.total_pages}
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