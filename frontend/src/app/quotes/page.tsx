'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { useQuery } from '@tanstack/react-query'
import DashboardLayout from '@/components/layout/DashboardLayout'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { 
  Table, 
  TableBody, 
  TableCell, 
  TableHead, 
  TableHeader, 
  TableRow 
} from '@/components/ui/table'
import { 
  Search, 
  Eye, 
  Edit, 
  FileText, 
  Send,
  Copy,
  CheckCircle,
  XCircle,
  Clock,
  DollarSign
} from 'lucide-react'
import quoteService from '@/services/quote.service'
import { Quote } from '@/types/quote'
import { useAuthStore } from '@/store/auth.store'
import { format } from 'date-fns'
import { zhTW } from 'date-fns/locale'

export default function QuotesPage() {
  const router = useRouter()
  const { user } = useAuthStore()
  const [searchTerm, setSearchTerm] = useState('')
  const [currentPage, setCurrentPage] = useState(1)
  const [statusFilter, setStatusFilter] = useState('all')

  const { data, isLoading } = useQuery({
    queryKey: ['quotes', currentPage, searchTerm, statusFilter],
    queryFn: () => quoteService.getQuotes({
      page: currentPage,
      page_size: 10,
      status: statusFilter === 'all' ? undefined : statusFilter,
    }),
  })

  const getStatusBadge = (status: string) => {
    const statusConfig: Record<string, { label: string; variant: any; icon: any }> = {
      draft: { label: '草稿', variant: 'secondary', icon: FileText },
      pending_approval: { label: '待審核', variant: 'warning', icon: Clock },
      approved: { label: '已核准', variant: 'success', icon: CheckCircle },
      sent: { label: '已發送', variant: 'info', icon: Send },
      accepted: { label: '已接受', variant: 'success', icon: CheckCircle },
      rejected: { label: '已拒絕', variant: 'destructive', icon: XCircle },
      expired: { label: '已過期', variant: 'secondary', icon: Clock },
      cancelled: { label: '已取消', variant: 'secondary', icon: XCircle },
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

  const handleView = (id: string) => {
    router.push(`/quotes/${id}`)
  }

  const handleEdit = (id: string) => {
    router.push(`/quotes/${id}/edit`)
  }

  const handleDuplicate = async (id: string) => {
    try {
      // TODO: Implement duplicate functionality
      console.log('Duplicate quote:', id)
    } catch (error) {
      console.error('複製失敗:', error)
    }
  }

  const statusOptions = [
    { value: 'all', label: '全部狀態' },
    { value: 'draft', label: '草稿' },
    { value: 'pending_approval', label: '待審核' },
    { value: 'approved', label: '已核准' },
    { value: 'sent', label: '已發送' },
  ]

  // Summary statistics
  const summaryStats = {
    total: data?.pagination.total || 0,
    draft: 0,
    pending_approval: 0,
    approved: 0,
    sent: 0,
  }

  if (data?.data) {
    data.data.forEach(quote => {
      if (quote.status === 'draft') summaryStats.draft++
      else if (quote.status === 'pending_approval') summaryStats.pending_approval++
      else if (quote.status === 'approved') summaryStats.approved++
      else if (quote.status === 'sent') summaryStats.sent++
    })
  }

  return (
    <DashboardLayout>
      <div className="space-y-6">
        {/* Header */}
        <div>
          <h1 className="text-3xl font-bold text-gray-900">報價單管理</h1>
          <p className="mt-2 text-gray-600">管理所有報價單和審核流程</p>
        </div>

        {/* Summary Cards */}
        <div className="grid grid-cols-1 md:grid-cols-5 gap-4">
          <Card className="cursor-pointer hover:shadow-md transition-shadow" onClick={() => setStatusFilter('all')}>
            <CardContent className="p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-600">總報價單</p>
                  <p className="text-2xl font-bold">{summaryStats.total}</p>
                </div>
                <DollarSign className="h-8 w-8 text-gray-400" />
              </div>
            </CardContent>
          </Card>

          <Card className="cursor-pointer hover:shadow-md transition-shadow" onClick={() => setStatusFilter('draft')}>
            <CardContent className="p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-600">草稿</p>
                  <p className="text-2xl font-bold">{summaryStats.draft}</p>
                </div>
                <FileText className="h-8 w-8 text-gray-400" />
              </div>
            </CardContent>
          </Card>

          <Card className="cursor-pointer hover:shadow-md transition-shadow" onClick={() => setStatusFilter('pending_approval')}>
            <CardContent className="p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-600">待審核</p>
                  <p className="text-2xl font-bold text-yellow-600">{summaryStats.pending_approval}</p>
                </div>
                <Clock className="h-8 w-8 text-yellow-400" />
              </div>
            </CardContent>
          </Card>

          <Card className="cursor-pointer hover:shadow-md transition-shadow" onClick={() => setStatusFilter('approved')}>
            <CardContent className="p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-600">已核准</p>
                  <p className="text-2xl font-bold text-green-600">{summaryStats.approved}</p>
                </div>
                <CheckCircle className="h-8 w-8 text-green-400" />
              </div>
            </CardContent>
          </Card>

          <Card className="cursor-pointer hover:shadow-md transition-shadow" onClick={() => setStatusFilter('sent')}>
            <CardContent className="p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-600">已發送</p>
                  <p className="text-2xl font-bold text-blue-600">{summaryStats.sent}</p>
                </div>
                <Send className="h-8 w-8 text-blue-400" />
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Search and Filters */}
        <Card>
          <CardContent className="p-6">
            <div className="flex gap-4">
              <div className="flex-1 relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-4 w-4" />
                <Input
                  placeholder="搜尋報價單號、客戶名稱..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className="pl-10"
                />
              </div>
              <select
                className="px-4 py-2 border rounded-md"
                value={statusFilter}
                onChange={(e) => setStatusFilter(e.target.value)}
              >
                {statusOptions.map(option => (
                  <option key={option.value} value={option.value}>
                    {option.label}
                  </option>
                ))}
              </select>
            </div>
          </CardContent>
        </Card>

        {/* Quotes Table */}
        <Card>
          <CardHeader>
            <CardTitle>報價單列表</CardTitle>
            <CardDescription>
              {statusFilter !== 'all' && `篩選: ${statusOptions.find(opt => opt.value === statusFilter)?.label}`}
            </CardDescription>
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <div className="text-center py-8">載入中...</div>
            ) : (
              <div className="overflow-x-auto">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>報價單號</TableHead>
                      <TableHead>詢價單號</TableHead>
                      <TableHead>客戶</TableHead>
                      <TableHead>總金額</TableHead>
                      <TableHead>單價</TableHead>
                      <TableHead>有效期限</TableHead>
                      <TableHead>狀態</TableHead>
                      <TableHead>工程師</TableHead>
                      <TableHead className="text-right">操作</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {data?.data.map((quote) => (
                      <TableRow key={quote.id}>
                        <TableCell className="font-medium">
                          {quote.quote_no}
                        </TableCell>
                        <TableCell>{quote.inquiry?.inquiry_no || '-'}</TableCell>
                        <TableCell>{quote.customer?.name || '-'}</TableCell>
                        <TableCell>
                          <span className="font-medium">
                            ${quote.total_amount.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
                          </span>
                        </TableCell>
                        <TableCell>
                          ${(quote.total_amount / (quote.inquiry?.quantity || 1)).toFixed(4)}
                        </TableCell>
                        <TableCell>
                          {quote.valid_until ? format(new Date(quote.valid_until), 'yyyy/MM/dd', {
                            locale: zhTW,
                          }) : '-'}
                        </TableCell>
                        <TableCell>{getStatusBadge(quote.status)}</TableCell>
                        <TableCell>{quote.created_by_user?.name || '-'}</TableCell>
                        <TableCell className="text-right">
                          <div className="flex justify-end gap-2">
                            <Button
                              variant="ghost"
                              size="icon"
                              onClick={() => handleView(quote.id)}
                            >
                              <Eye className="h-4 w-4" />
                            </Button>
                            {(quote.status === 'draft' || quote.status === 'rejected') && (
                              <Button
                                variant="ghost"
                                size="icon"
                                onClick={() => handleEdit(quote.id)}
                              >
                                <Edit className="h-4 w-4" />
                              </Button>
                            )}
                            <Button
                              variant="ghost"
                              size="icon"
                              onClick={() => handleDuplicate(quote.id)}
                              title="複製報價單"
                            >
                              <Copy className="h-4 w-4" />
                            </Button>
                          </div>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Pagination */}
        {data && data.pagination.total > data.pagination.page_size && (
          <div className="flex justify-center gap-2">
            <Button
              variant="outline"
              disabled={currentPage === 1}
              onClick={() => setCurrentPage(currentPage - 1)}
            >
              上一頁
            </Button>
            <span className="flex items-center px-4">
              第 {currentPage} 頁，共 {Math.ceil(data.pagination.total / data.pagination.page_size)} 頁
            </span>
            <Button
              variant="outline"
              disabled={currentPage * data.pagination.page_size >= data.pagination.total}
              onClick={() => setCurrentPage(currentPage + 1)}
            >
              下一頁
            </Button>
          </div>
        )}
      </div>
    </DashboardLayout>
  )
}