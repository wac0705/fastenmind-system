'use client'

import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
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
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Textarea } from '@/components/ui/textarea'
import { Label } from '@/components/ui/label'
import { 
  Search,
  Receipt,
  Plus,
  Eye,
  Download,
  Calendar,
  DollarSign,
  CheckCircle,
  XCircle,
  Clock,
  AlertCircle,
  User,
  Building2,
  FileText,
  Trash2
} from 'lucide-react'
import financeService from '@/services/finance.service'
import Pagination from '@/components/common/Pagination'
import { format } from 'date-fns'
import { zhTW } from 'date-fns/locale'
import { useToast } from '@/components/ui/use-toast'
import { useAuthStore } from '@/store/auth.store'

export default function ExpensesPage() {
  const router = useRouter()
  const queryClient = useQueryClient()
  const { toast } = useToast()
  const { user } = useAuthStore()
  
  const [page, setPage] = useState(1)
  const [search, setSearch] = useState('')
  const [statusFilter, setStatusFilter] = useState<string>('')
  const [categoryFilter, setCategoryFilter] = useState<string>('')
  const [isApprovalDialogOpen, setIsApprovalDialogOpen] = useState(false)
  const [selectedExpense, setSelectedExpense] = useState<any>(null)
  const [approvalAction, setApprovalAction] = useState<'approve' | 'reject'>('approve')
  const [approvalNotes, setApprovalNotes] = useState('')

  // Fetch expenses
  const { data: expensesData, isLoading } = useQuery({
    queryKey: ['expenses', page, search, statusFilter, categoryFilter],
    queryFn: () => financeService.listExpenses({
      page,
      page_size: 20,
      search,
      status: statusFilter,
      category: categoryFilter,
    }),
  })

  // Approve expense mutation
  const approveExpenseMutation = useMutation({
    mutationFn: (id: string) => financeService.approveExpense(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['expenses'] })
      toast({ title: '費用已核准' })
      setIsApprovalDialogOpen(false)
      setApprovalNotes('')
    },
    onError: (error: any) => {
      toast({
        title: '核准失敗',
        description: error.response?.data?.message || '核准費用時發生錯誤',
        variant: 'destructive',
      })
    },
  })

  // Reject expense mutation
  const rejectExpenseMutation = useMutation({
    mutationFn: ({ id, reason }: { id: string; reason: string }) => financeService.rejectExpense(id, reason),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['expenses'] })
      toast({ title: '費用已拒絕' })
      setIsApprovalDialogOpen(false)
      setApprovalNotes('')
    },
    onError: (error: any) => {
      toast({
        title: '拒絕失敗',
        description: error.response?.data?.message || '拒絕費用時發生錯誤',
        variant: 'destructive',
      })
    },
  })

  // Calculate statistics
  const stats = expensesData?.data.reduce((acc, expense) => {
    acc.total++
    acc.totalAmount += expense.total_amount
    
    if (expense.status === 'submitted') acc.pending++
    else if (expense.status === 'approved') acc.approved++
    else if (expense.status === 'rejected') acc.rejected++
    else if (expense.status === 'paid') acc.paid++
    
    return acc
  }, {
    total: 0,
    totalAmount: 0,
    pending: 0,
    approved: 0,
    rejected: 0,
    paid: 0,
  })

  const getExpenseStatusBadge = (status: string) => {
    const statusConfig: Record<string, { label: string; variant: any; icon: any }> = {
      draft: { label: '草稿', variant: 'secondary', icon: FileText },
      submitted: { label: '待審核', variant: 'warning', icon: Clock },
      approved: { label: '已核准', variant: 'success', icon: CheckCircle },
      paid: { label: '已付款', variant: 'success', icon: CheckCircle },
      rejected: { label: '已拒絕', variant: 'destructive', icon: XCircle },
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

  const handleApproval = (expense: any, action: 'approve' | 'reject') => {
    setSelectedExpense(expense)
    setApprovalAction(action)
    setIsApprovalDialogOpen(true)
  }

  const handleApprovalSubmit = () => {
    if (!selectedExpense) return

    if (approvalAction === 'approve') {
      approveExpenseMutation.mutate(selectedExpense.id)
    } else {
      if (!approvalNotes.trim()) {
        toast({
          title: '請填寫拒絕原因',
          variant: 'destructive'
        })
        return
      }
      rejectExpenseMutation.mutate({ id: selectedExpense.id, reason: approvalNotes })
    }
  }

  const handleSearch = (value: string) => {
    setSearch(value)
    setPage(1)
  }

  const canApprove = (expense: any) => {
    return ['admin', 'finance', 'manager'].includes(user?.role || '') && 
           expense.status === 'submitted' &&
           expense.submitted_by !== user?.id
  }

  return (
    <DashboardLayout>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex justify-between items-center">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">費用管理</h1>
            <p className="mt-2 text-gray-600">管理員工費用申請與審核</p>
          </div>
          <div className="flex gap-2">
            <Button variant="outline">
              <Download className="mr-2 h-4 w-4" />
              匯出資料
            </Button>
            <Button onClick={() => router.push('/finance/expenses/new')}>
              <Plus className="mr-2 h-4 w-4" />
              新增費用
            </Button>
          </div>
        </div>

        {/* Statistics Cards */}
        {stats && (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-6">
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">總費用</CardTitle>
                <DollarSign className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">
                  ${stats.totalAmount.toLocaleString()}
                </div>
                <p className="text-xs text-muted-foreground">
                  {stats.total} 筆費用
                </p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">待審核</CardTitle>
                <Clock className="h-4 w-4 text-yellow-600" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold text-yellow-600">{stats.pending}</div>
                <p className="text-xs text-muted-foreground">需要審核</p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">已核准</CardTitle>
                <CheckCircle className="h-4 w-4 text-green-600" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold text-green-600">{stats.approved}</div>
                <p className="text-xs text-muted-foreground">等待付款</p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">已付款</CardTitle>
                <CheckCircle className="h-4 w-4 text-blue-600" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold text-blue-600">{stats.paid}</div>
                <p className="text-xs text-muted-foreground">完成處理</p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">已拒絕</CardTitle>
                <XCircle className="h-4 w-4 text-red-600" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold text-red-600">{stats.rejected}</div>
                <p className="text-xs text-muted-foreground">不符要求</p>
              </CardContent>
            </Card>
          </div>
        )}

        {/* Filters */}
        <Card>
          <CardHeader>
            <CardTitle>費用列表</CardTitle>
            <CardDescription>查看和管理所有費用申請</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="flex flex-col sm:flex-row gap-4 mb-6">
              <div className="relative flex-1">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
                <Input
                  placeholder="搜尋費用編號、描述..."
                  value={search}
                  onChange={(e) => handleSearch(e.target.value)}
                  className="pl-10"
                />
              </div>
              <Select value={statusFilter} onValueChange={(value) => { setStatusFilter(value); setPage(1); }}>
                <SelectTrigger className="w-full sm:w-[150px]">
                  <SelectValue placeholder="狀態" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="">全部狀態</SelectItem>
                  <SelectItem value="draft">草稿</SelectItem>
                  <SelectItem value="submitted">待審核</SelectItem>
                  <SelectItem value="approved">已核准</SelectItem>
                  <SelectItem value="paid">已付款</SelectItem>
                  <SelectItem value="rejected">已拒絕</SelectItem>
                </SelectContent>
              </Select>
              <Select value={categoryFilter} onValueChange={(value) => { setCategoryFilter(value); setPage(1); }}>
                <SelectTrigger className="w-full sm:w-[180px]">
                  <SelectValue placeholder="費用類別" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="">全部類別</SelectItem>
                  <SelectItem value="travel">差旅費</SelectItem>
                  <SelectItem value="meal">餐費</SelectItem>
                  <SelectItem value="transportation">交通費</SelectItem>
                  <SelectItem value="office_supplies">辦公用品</SelectItem>
                  <SelectItem value="communication">通訊費</SelectItem>
                  <SelectItem value="training">培訓費</SelectItem>
                  <SelectItem value="marketing">行銷費用</SelectItem>
                  <SelectItem value="equipment">設備費用</SelectItem>
                  <SelectItem value="other">其他</SelectItem>
                </SelectContent>
              </Select>
            </div>

            {/* Expenses Table */}
            <div className="rounded-md border">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>費用編號</TableHead>
                    <TableHead>類別</TableHead>
                    <TableHead>描述</TableHead>
                    <TableHead>申請人</TableHead>
                    <TableHead>費用日期</TableHead>
                    <TableHead className="text-right">金額</TableHead>
                    <TableHead>狀態</TableHead>
                    <TableHead>操作</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {isLoading ? (
                    <TableRow>
                      <TableCell colSpan={8} className="text-center py-8">
                        載入中...
                      </TableCell>
                    </TableRow>
                  ) : expensesData?.data.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={8} className="text-center py-8 text-gray-500">
                        沒有找到費用記錄
                      </TableCell>
                    </TableRow>
                  ) : (
                    expensesData?.data.map((expense) => (
                      <TableRow key={expense.id} className="hover:bg-gray-50">
                        <TableCell className="font-medium">{expense.expense_no}</TableCell>
                        <TableCell>
                          <Badge variant="outline">{expense.category}</Badge>
                        </TableCell>
                        <TableCell className="max-w-xs truncate">{expense.description}</TableCell>
                        <TableCell>
                          <div className="flex items-center gap-2">
                            <User className="h-4 w-4 text-gray-400" />
                            <span>{expense.submitter?.full_name}</span>
                          </div>
                        </TableCell>
                        <TableCell>
                          <div className="flex items-center gap-1">
                            <Calendar className="h-3 w-3 text-gray-400" />
                            <span className="text-sm">
                              {format(new Date(expense.expense_date), 'yyyy/MM/dd', { locale: zhTW })}
                            </span>
                          </div>
                        </TableCell>
                        <TableCell className="text-right">
                          <div className="flex items-center justify-end gap-1">
                            <span className="text-xs text-gray-500">{expense.currency}</span>
                            <span className="font-medium">{expense.total_amount.toLocaleString()}</span>
                          </div>
                        </TableCell>
                        <TableCell>{getExpenseStatusBadge(expense.status)}</TableCell>
                        <TableCell>
                          <div className="flex items-center gap-1">
                            <Button
                              variant="ghost"
                              size="icon"
                              onClick={() => router.push(`/finance/expenses/${expense.id}`)}
                            >
                              <Eye className="h-4 w-4" />
                            </Button>
                            {canApprove(expense) && (
                              <>
                                <Button
                                  variant="ghost"
                                  size="icon"
                                  onClick={() => handleApproval(expense, 'approve')}
                                  className="text-green-600 hover:text-green-900"
                                >
                                  <CheckCircle className="h-4 w-4" />
                                </Button>
                                <Button
                                  variant="ghost"
                                  size="icon"
                                  onClick={() => handleApproval(expense, 'reject')}
                                  className="text-red-600 hover:text-red-900"
                                >
                                  <XCircle className="h-4 w-4" />
                                </Button>
                              </>
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
            {expensesData && expensesData.pagination.total > 0 && (
              <div className="mt-4">
                <Pagination
                  currentPage={page}
                  totalPages={expensesData.pagination.total_pages}
                  onPageChange={setPage}
                />
              </div>
            )}
          </CardContent>
        </Card>

        {/* Approval Dialog */}
        <Dialog open={isApprovalDialogOpen} onOpenChange={setIsApprovalDialogOpen}>
          <DialogContent className="sm:max-w-[500px]">
            <DialogHeader>
              <DialogTitle>
                {approvalAction === 'approve' ? '核准費用' : '拒絕費用'}
              </DialogTitle>
              <DialogDescription>
                費用編號：{selectedExpense?.expense_no}
                <br />
                金額：{selectedExpense?.currency} {selectedExpense?.total_amount.toLocaleString()}
              </DialogDescription>
            </DialogHeader>
            <div className="grid gap-4 py-4">
              {approvalAction === 'reject' && (
                <div className="grid gap-2">
                  <Label htmlFor="approval-notes">拒絕原因 *</Label>
                  <Textarea
                    id="approval-notes"
                    value={approvalNotes}
                    onChange={(e) => setApprovalNotes(e.target.value)}
                    placeholder="請說明拒絕的原因..."
                    rows={4}
                  />
                </div>
              )}
              {approvalAction === 'approve' && (
                <div className="p-4 bg-green-50 rounded-md">
                  <p className="text-sm text-green-800">
                    確認核准此費用申請？核准後將可進行付款作業。
                  </p>
                </div>
              )}
            </div>
            <DialogFooter>
              <Button variant="outline" onClick={() => setIsApprovalDialogOpen(false)}>
                取消
              </Button>
              <Button 
                variant={approvalAction === 'approve' ? 'default' : 'destructive'}
                onClick={handleApprovalSubmit}
                disabled={approveExpenseMutation.isPending || rejectExpenseMutation.isPending}
              >
                {approveExpenseMutation.isPending || rejectExpenseMutation.isPending 
                  ? '處理中...' 
                  : (approvalAction === 'approve' ? '確認核准' : '確認拒絕')}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>
    </DashboardLayout>
  )
}