'use client'

import { useState, useEffect } from 'react'
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
import { Plus, Search, Eye, Edit, Trash2, FileText } from 'lucide-react'
import inquiryService, { Inquiry } from '@/services/inquiry.service'
import { useAuthStore } from '@/store/auth.store'
import { format } from 'date-fns'
import { zhTW } from 'date-fns/locale'

export default function InquiriesPage() {
  const router = useRouter()
  const { user } = useAuthStore()
  const [searchTerm, setSearchTerm] = useState('')
  const [currentPage, setCurrentPage] = useState(1)

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['inquiries', currentPage, searchTerm],
    queryFn: () => inquiryService.list({
      page: currentPage,
      page_size: 10,
      search: searchTerm,
    }),
  })

  const getStatusBadge = (status: string) => {
    const statusConfig: Record<string, { label: string; variant: any }> = {
      draft: { label: '草稿', variant: 'secondary' },
      pending: { label: '待處理', variant: 'warning' },
      assigned: { label: '已分派', variant: 'info' },
      in_progress: { label: '處理中', variant: 'default' },
      under_review: { label: '審核中', variant: 'warning' },
      approved: { label: '已核准', variant: 'success' },
      quoted: { label: '已報價', variant: 'success' },
      rejected: { label: '已拒絕', variant: 'destructive' },
      cancelled: { label: '已取消', variant: 'secondary' },
    }

    const config = statusConfig[status] || { label: status, variant: 'default' }
    
    return (
      <Badge variant={config.variant as any}>
        {config.label}
      </Badge>
    )
  }

  const handleView = (id: string) => {
    router.push(`/inquiries/${id}`)
  }

  const handleEdit = (id: string) => {
    router.push(`/inquiries/${id}/edit`)
  }

  const handleDelete = async (id: string) => {
    if (confirm('確定要刪除這筆詢價單嗎？')) {
      try {
        await inquiryService.delete(id)
        refetch()
      } catch (error) {
        console.error('刪除失敗:', error)
      }
    }
  }

  return (
    <DashboardLayout>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex justify-between items-center">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">詢價單管理</h1>
            <p className="mt-2 text-gray-600">管理所有的詢價單和報價請求</p>
          </div>
          {user?.role === 'sales' && (
            <Button onClick={() => router.push('/inquiries/new')}>
              <Plus className="mr-2 h-4 w-4" />
              新增詢價單
            </Button>
          )}
        </div>

        {/* Search and Filters */}
        <Card>
          <CardContent className="p-6">
            <div className="flex gap-4">
              <div className="flex-1 relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-4 w-4" />
                <Input
                  placeholder="搜尋詢價單號、客戶名稱或產品..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className="pl-10"
                />
              </div>
              <Button variant="outline">篩選</Button>
            </div>
          </CardContent>
        </Card>

        {/* Inquiries Table */}
        <Card>
          <CardHeader>
            <CardTitle>詢價單列表</CardTitle>
            <CardDescription>
              共 {data?.pagination.total || 0} 筆詢價單
            </CardDescription>
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <div className="text-center py-8">載入中...</div>
            ) : error ? (
              <div className="text-center py-8 text-red-500">載入失敗</div>
            ) : (
              <div className="overflow-x-auto">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>詢價單號</TableHead>
                      <TableHead>客戶</TableHead>
                      <TableHead>產品</TableHead>
                      <TableHead>數量</TableHead>
                      <TableHead>交期</TableHead>
                      <TableHead>狀態</TableHead>
                      <TableHead>負責工程師</TableHead>
                      <TableHead className="text-right">操作</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {data?.data.map((inquiry) => (
                      <TableRow key={inquiry.id}>
                        <TableCell className="font-medium">
                          {inquiry.inquiry_no}
                        </TableCell>
                        <TableCell>{inquiry.customer?.name || '-'}</TableCell>
                        <TableCell>
                          <div>
                            <p className="font-medium">{inquiry.product_name}</p>
                            <p className="text-sm text-gray-500">
                              {inquiry.product_category}
                            </p>
                          </div>
                        </TableCell>
                        <TableCell>
                          {inquiry.quantity.toLocaleString()} {inquiry.unit}
                        </TableCell>
                        <TableCell>
                          {format(new Date(inquiry.required_date), 'yyyy/MM/dd', {
                            locale: zhTW,
                          })}
                        </TableCell>
                        <TableCell>{getStatusBadge(inquiry.status)}</TableCell>
                        <TableCell>
                          {inquiry.assigned_engineer?.full_name || '-'}
                        </TableCell>
                        <TableCell className="text-right">
                          <div className="flex justify-end gap-2">
                            <Button
                              variant="ghost"
                              size="icon"
                              onClick={() => handleView(inquiry.id)}
                            >
                              <Eye className="h-4 w-4" />
                            </Button>
                            {(user?.role === 'sales' || user?.role === 'admin') && (
                              <>
                                <Button
                                  variant="ghost"
                                  size="icon"
                                  onClick={() => handleEdit(inquiry.id)}
                                >
                                  <Edit className="h-4 w-4" />
                                </Button>
                                <Button
                                  variant="ghost"
                                  size="icon"
                                  onClick={() => handleDelete(inquiry.id)}
                                >
                                  <Trash2 className="h-4 w-4" />
                                </Button>
                              </>
                            )}
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