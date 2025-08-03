'use client'

import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { useRouter } from 'next/navigation'
import DashboardLayout from '@/components/layout/DashboardLayout'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
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
  DollarSign,
  TrendingUp,
  TrendingDown,
  AlertCircle,
  Receipt,
  CreditCard,
  Building2,
  FileText,
  Clock,
  Plus,
  Eye,
  Calendar,
  Users,
  PieChart,
  BarChart3,
  ArrowUpCircle,
  ArrowDownCircle
} from 'lucide-react'
import financeService, { FinancialDashboard } from '@/services/finance.service'
import { format } from 'date-fns'
import { zhTW } from 'date-fns/locale'

export default function FinancePage() {
  const router = useRouter()
  const [activeTab, setActiveTab] = useState('dashboard')

  // Fetch financial dashboard
  const { data: dashboard, isLoading: isLoadingDashboard } = useQuery({
    queryKey: ['financial-dashboard'],
    queryFn: () => financeService.getFinancialDashboard(),
  })

  // Fetch recent invoices
  const { data: recentInvoices } = useQuery({
    queryKey: ['recent-invoices'],
    queryFn: () => financeService.listInvoices({ page: 1, page_size: 5 }),
  })

  // Fetch recent payments
  const { data: recentPayments } = useQuery({
    queryKey: ['recent-payments'],
    queryFn: () => financeService.listPayments({ page: 1, page_size: 5 }),
  })

  // Fetch pending expenses
  const { data: pendingExpenses } = useQuery({
    queryKey: ['pending-expenses'],
    queryFn: () => financeService.listExpenses({ status: 'submitted', page: 1, page_size: 5 }),
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

  const getPaymentStatusBadge = (status: string) => {
    const statusConfig: Record<string, { label: string; variant: any }> = {
      pending: { label: '處理中', variant: 'warning' },
      completed: { label: '已完成', variant: 'success' },
      failed: { label: '失敗', variant: 'destructive' },
      cancelled: { label: '已取消', variant: 'secondary' },
    }

    const config = statusConfig[status] || { label: status, variant: 'default' }
    return <Badge variant={config.variant as any}>{config.label}</Badge>
  }

  if (isLoadingDashboard) {
    return (
      <DashboardLayout>
        <div className="flex items-center justify-center h-64">
          <div className="text-center">載入中...</div>
        </div>
      </DashboardLayout>
    )
  }

  return (
    <DashboardLayout>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex justify-between items-center">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">財務管理</h1>
            <p className="mt-2 text-gray-600">管理發票、付款、費用與財務報表</p>
          </div>
          <div className="flex gap-2">
            <Button variant="outline" onClick={() => router.push('/finance/reports')}>
              <BarChart3 className="mr-2 h-4 w-4" />
              財務報表
            </Button>
            <Button onClick={() => router.push('/finance/invoices/new')}>
              <Plus className="mr-2 h-4 w-4" />
              新增發票
            </Button>
          </div>
        </div>

        <Tabs value={activeTab} onValueChange={setActiveTab}>
          <TabsList>
            <TabsTrigger value="dashboard">財務概覽</TabsTrigger>
            <TabsTrigger value="invoices">發票管理</TabsTrigger>
            <TabsTrigger value="payments">付款記錄</TabsTrigger>
            <TabsTrigger value="expenses">費用管理</TabsTrigger>
            <TabsTrigger value="receivables">應收帳款</TabsTrigger>
          </TabsList>

          <TabsContent value="dashboard" className="space-y-6">
            {/* Financial KPIs */}
            {dashboard && (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                <Card>
                  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">本月營收</CardTitle>
                    <TrendingUp className="h-4 w-4 text-muted-foreground" />
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold">
                      {dashboard.currency} {dashboard.revenue.toLocaleString()}
                    </div>
                    <p className="text-xs text-muted-foreground">銷售發票總額</p>
                  </CardContent>
                </Card>
                <Card>
                  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">本月費用</CardTitle>
                    <TrendingDown className="h-4 w-4 text-muted-foreground" />
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold">
                      {dashboard.currency} {dashboard.expenses.toLocaleString()}
                    </div>
                    <p className="text-xs text-muted-foreground">已批准費用總額</p>
                  </CardContent>
                </Card>
                <Card>
                  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">本月淨利</CardTitle>
                    <DollarSign className="h-4 w-4 text-muted-foreground" />
                  </CardHeader>
                  <CardContent>
                    <div className={`text-2xl font-bold ${dashboard.profit >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                      {dashboard.currency} {dashboard.profit.toLocaleString()}
                    </div>
                    <p className="text-xs text-muted-foreground">營收 - 費用</p>
                  </CardContent>
                </Card>
                <Card>
                  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">銀行餘額</CardTitle>
                    <Building2 className="h-4 w-4 text-muted-foreground" />
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold">
                      {dashboard.currency} {dashboard.cash_balance.toLocaleString()}
                    </div>
                    <p className="text-xs text-muted-foreground">所有銀行帳戶</p>
                  </CardContent>
                </Card>
              </div>
            )}

            {/* AR/AP Summary */}
            {dashboard && (
              <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                <Card>
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                      <ArrowUpCircle className="h-5 w-5 text-green-600" />
                      應收帳款
                    </CardTitle>
                    <CardDescription>客戶欠款情況</CardDescription>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    <div className="flex justify-between items-center">
                      <span className="text-sm text-gray-500">總金額</span>
                      <span className="font-medium">
                        {dashboard.ar_summary.currency} {dashboard.ar_summary.total_amount.toLocaleString()}
                      </span>
                    </div>
                    <div className="flex justify-between items-center">
                      <span className="text-sm text-gray-500">已收金額</span>
                      <span className="font-medium text-green-600">
                        {dashboard.ar_summary.currency} {dashboard.ar_summary.paid_amount.toLocaleString()}
                      </span>
                    </div>
                    <div className="flex justify-between items-center">
                      <span className="text-sm text-gray-500">未收金額</span>
                      <span className="font-medium text-orange-600">
                        {dashboard.ar_summary.currency} {dashboard.ar_summary.balance_amount.toLocaleString()}
                      </span>
                    </div>
                    <div className="space-y-2">
                      <div className="flex justify-between text-sm">
                        <span>當期: {dashboard.ar_summary.currency} {dashboard.ar_summary.current.toLocaleString()}</span>
                        <span>30天: {dashboard.ar_summary.currency} {dashboard.ar_summary.days_30.toLocaleString()}</span>
                      </div>
                      <div className="flex justify-between text-sm">
                        <span>60天: {dashboard.ar_summary.currency} {dashboard.ar_summary.days_60.toLocaleString()}</span>
                        <span>90天+: {dashboard.ar_summary.currency} {dashboard.ar_summary.over_90.toLocaleString()}</span>
                      </div>
                    </div>
                  </CardContent>
                </Card>

                <Card>
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                      <ArrowDownCircle className="h-5 w-5 text-red-600" />
                      應付帳款
                    </CardTitle>
                    <CardDescription>供應商付款情況</CardDescription>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    <div className="flex justify-between items-center">
                      <span className="text-sm text-gray-500">總金額</span>
                      <span className="font-medium">
                        {dashboard.ap_summary.currency} {dashboard.ap_summary.total_amount.toLocaleString()}
                      </span>
                    </div>
                    <div className="flex justify-between items-center">
                      <span className="text-sm text-gray-500">已付金額</span>
                      <span className="font-medium text-green-600">
                        {dashboard.ap_summary.currency} {dashboard.ap_summary.paid_amount.toLocaleString()}
                      </span>
                    </div>
                    <div className="flex justify-between items-center">
                      <span className="text-sm text-gray-500">未付金額</span>
                      <span className="font-medium text-red-600">
                        {dashboard.ap_summary.currency} {dashboard.ap_summary.balance_amount.toLocaleString()}
                      </span>
                    </div>
                    <div className="text-center">
                      <span className="text-sm text-gray-500">
                        {dashboard.ap_summary.open_items} 張未付發票
                      </span>
                    </div>
                  </CardContent>
                </Card>
              </div>
            )}

            {/* Recent Activities */}
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              {/* Recent Invoices */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center justify-between">
                    <span className="flex items-center gap-2">
                      <FileText className="h-5 w-5" />
                      最近發票
                    </span>
                    <Button variant="ghost" size="sm" onClick={() => router.push('/finance/invoices')}>
                      <Eye className="h-4 w-4" />
                    </Button>
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  {recentInvoices?.data.length === 0 ? (
                    <p className="text-center text-gray-500 py-4">暫無發票記錄</p>
                  ) : (
                    <div className="space-y-3">
                      {recentInvoices?.data.map((invoice) => (
                        <div key={invoice.id} className="flex items-center justify-between p-3 border rounded-lg">
                          <div>
                            <p className="font-medium">{invoice.invoice_no}</p>
                            <p className="text-sm text-gray-500">
                              {invoice.customer?.name || invoice.supplier?.name}
                            </p>
                          </div>
                          <div className="text-right">
                            <p className="font-medium">
                              {invoice.currency} {invoice.total_amount.toLocaleString()}
                            </p>
                            {getInvoiceStatusBadge(invoice.status)}
                          </div>
                        </div>
                      ))}
                    </div>
                  )}
                </CardContent>
              </Card>

              {/* Recent Payments */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center justify-between">
                    <span className="flex items-center gap-2">
                      <CreditCard className="h-5 w-5" />
                      最近付款
                    </span>
                    <Button variant="ghost" size="sm" onClick={() => router.push('/finance/payments')}>
                      <Eye className="h-4 w-4" />
                    </Button>
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  {recentPayments?.data.length === 0 ? (
                    <p className="text-center text-gray-500 py-4">暫無付款記錄</p>
                  ) : (
                    <div className="space-y-3">
                      {recentPayments?.data.map((payment) => (
                        <div key={payment.id} className="flex items-center justify-between p-3 border rounded-lg">
                          <div>
                            <p className="font-medium">{payment.payment_no}</p>
                            <p className="text-sm text-gray-500">
                              {format(new Date(payment.payment_date), 'yyyy/MM/dd', { locale: zhTW })}
                            </p>
                          </div>
                          <div className="text-right">
                            <p className="font-medium">
                              {payment.currency} {payment.amount.toLocaleString()}
                            </p>
                            {getPaymentStatusBadge(payment.status)}
                          </div>
                        </div>
                      ))}
                    </div>
                  )}
                </CardContent>
              </Card>
            </div>

            {/* Pending Expenses */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center justify-between">
                  <span className="flex items-center gap-2">
                    <Receipt className="h-5 w-5" />
                    待審核費用
                  </span>
                  <Button variant="ghost" size="sm" onClick={() => router.push('/finance/expenses')}>
                    <Eye className="h-4 w-4" />
                  </Button>
                </CardTitle>
                <CardDescription>需要審核的費用申請</CardDescription>
              </CardHeader>
              <CardContent>
                {pendingExpenses?.data.length === 0 ? (
                  <p className="text-center text-gray-500 py-4">暫無待審核費用</p>
                ) : (
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead>費用編號</TableHead>
                        <TableHead>類別</TableHead>
                        <TableHead>描述</TableHead>
                        <TableHead>申請人</TableHead>
                        <TableHead className="text-right">金額</TableHead>
                        <TableHead>狀態</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {pendingExpenses?.data.map((expense) => (
                        <TableRow key={expense.id}>
                          <TableCell className="font-medium">{expense.expense_no}</TableCell>
                          <TableCell>{expense.category}</TableCell>
                          <TableCell>{expense.description}</TableCell>
                          <TableCell>{expense.submitter?.full_name}</TableCell>
                          <TableCell className="text-right">
                            {expense.currency} {expense.total_amount.toLocaleString()}
                          </TableCell>
                          <TableCell>
                            <Badge variant="warning">待審核</Badge>
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                )}
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="invoices">
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center justify-between">
                    <span className="flex items-center gap-2">
                      <FileText className="h-5 w-5" />
                      發票概覽
                    </span>
                    <Button variant="ghost" size="sm" onClick={() => router.push('/finance/invoices')}>
                      <Eye className="h-4 w-4" />
                    </Button>
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    <div className="flex justify-between items-center">
                      <span className="text-sm text-gray-500">本月銷售發票</span>
                      <span className="font-medium">
                        {dashboard?.currency} {dashboard?.revenue.toLocaleString()}
                      </span>
                    </div>
                    <div className="flex justify-between items-center">
                      <span className="text-sm text-gray-500">待收款發票</span>
                      <span className="font-medium text-orange-600">
                        {dashboard?.ar_summary.open_items} 張
                      </span>
                    </div>
                    <div className="flex justify-between items-center">
                      <span className="text-sm text-gray-500">逾期發票</span>
                      <span className="font-medium text-red-600">
                        {dashboard?.currency} {((dashboard?.ar_summary?.days_30 || 0) + (dashboard?.ar_summary?.days_60 || 0) + (dashboard?.ar_summary?.over_90 || 0)).toLocaleString()}
                      </span>
                    </div>
                  </div>
                  <Button className="w-full mt-4" onClick={() => router.push('/finance/invoices')}>
                    管理發票
                  </Button>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <DollarSign className="h-5 w-5" />
                    快速動作
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-3">
                  <Button 
                    variant="outline" 
                    className="w-full justify-start"
                    onClick={() => router.push('/finance/invoices/new')}
                  >
                    <Plus className="mr-2 h-4 w-4" />
                    建立新發票
                  </Button>
                  <Button 
                    variant="outline" 
                    className="w-full justify-start"
                    onClick={() => router.push('/finance/payments/new')}
                  >
                    <CreditCard className="mr-2 h-4 w-4" />
                    記錄付款
                  </Button>
                  <Button 
                    variant="outline" 
                    className="w-full justify-start"
                    onClick={() => router.push('/finance/expenses/new')}
                  >
                    <Receipt className="mr-2 h-4 w-4" />
                    新增費用
                  </Button>
                </CardContent>
              </Card>
            </div>
          </TabsContent>

          <TabsContent value="payments">
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center justify-between">
                    <span className="flex items-center gap-2">
                      <CreditCard className="h-5 w-5" />
                      付款統計
                    </span>
                    <Button variant="ghost" size="sm" onClick={() => router.push('/finance/payments')}>
                      <Eye className="h-4 w-4" />
                    </Button>
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    <div className="flex justify-between items-center">
                      <span className="text-sm text-gray-500">本月收款</span>
                      <span className="font-medium text-green-600">
                        {dashboard?.currency} {dashboard?.ar_summary.paid_amount.toLocaleString()}
                      </span>
                    </div>
                    <div className="flex justify-between items-center">
                      <span className="text-sm text-gray-500">本月支出</span>
                      <span className="font-medium text-red-600">
                        {dashboard?.currency} {dashboard?.expenses.toLocaleString()}
                      </span>
                    </div>
                    <div className="flex justify-between items-center">
                      <span className="text-sm text-gray-500">銀行餘額</span>
                      <span className="font-medium">
                        {dashboard?.currency} {dashboard?.cash_balance.toLocaleString()}
                      </span>
                    </div>
                  </div>
                  <Button className="w-full mt-4" onClick={() => router.push('/finance/payments')}>
                    查看付款記錄
                  </Button>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Building2 className="h-5 w-5" />
                    銀行帳戶
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-3">
                    <div className="p-3 border rounded-lg">
                      <div className="flex justify-between items-center">
                        <span className="text-sm font-medium">主要帳戶</span>
                        <span className="text-sm text-green-600 font-medium">
                          {dashboard?.currency} {dashboard?.cash_balance.toLocaleString()}
                        </span>
                      </div>
                      <p className="text-xs text-gray-500 mt-1">可用餘額</p>
                    </div>
                  </div>
                  <Button 
                    variant="outline" 
                    className="w-full mt-4"
                    onClick={() => router.push('/finance/bank-accounts')}
                  >
                    管理銀行帳戶
                  </Button>
                </CardContent>
              </Card>
            </div>
          </TabsContent>

          <TabsContent value="expenses">
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center justify-between">
                    <span className="flex items-center gap-2">
                      <Receipt className="h-5 w-5" />
                      費用統計
                    </span>
                    <Button variant="ghost" size="sm" onClick={() => router.push('/finance/expenses')}>
                      <Eye className="h-4 w-4" />
                    </Button>
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    <div className="flex justify-between items-center">
                      <span className="text-sm text-gray-500">本月總費用</span>
                      <span className="font-medium">
                        {dashboard?.currency} {dashboard?.expenses.toLocaleString()}
                      </span>
                    </div>
                    <div className="flex justify-between items-center">
                      <span className="text-sm text-gray-500">待審核費用</span>
                      <span className="font-medium text-yellow-600">
                        {pendingExpenses?.data.length || 0} 筆
                      </span>
                    </div>
                    <div className="flex justify-between items-center">
                      <span className="text-sm text-gray-500">費用類別</span>
                      <span className="text-sm text-gray-500">差旅、辦公、設備</span>
                    </div>
                  </div>
                  <Button className="w-full mt-4" onClick={() => router.push('/finance/expenses')}>
                    管理費用申請
                  </Button>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Clock className="h-5 w-5" />
                    待處理費用
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  {(pendingExpenses?.data?.length || 0) > 0 ? (
                    <div className="space-y-3">
                      {pendingExpenses?.data?.slice(0, 3).map((expense) => (
                        <div key={expense.id} className="p-3 border rounded-lg">
                          <div className="flex justify-between items-start">
                            <div>
                              <p className="text-sm font-medium">{expense.expense_no}</p>
                              <p className="text-xs text-gray-500">{expense.description}</p>
                            </div>
                            <span className="text-sm font-medium">
                              {expense.currency} {expense.total_amount.toLocaleString()}
                            </span>
                          </div>
                        </div>
                      ))}
                    </div>
                  ) : (
                    <p className="text-center text-gray-500 py-4">暫無待審核費用</p>
                  )}
                  <Button 
                    variant="outline" 
                    className="w-full mt-4"
                    onClick={() => router.push('/finance/expenses?status=submitted')}
                  >
                    查看所有待審核
                  </Button>
                </CardContent>
              </Card>
            </div>
          </TabsContent>

          <TabsContent value="receivables">
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center justify-between">
                    <span className="flex items-center gap-2">
                      <ArrowUpCircle className="h-5 w-5 text-green-600" />
                      應收帳款分析
                    </span>
                    <Button variant="ghost" size="sm" onClick={() => router.push('/finance/receivables')}>
                      <Eye className="h-4 w-4" />
                    </Button>
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    <div className="flex justify-between items-center">
                      <span className="text-sm text-gray-500">總應收金額</span>
                      <span className="font-medium">
                        {dashboard?.ar_summary.currency} {dashboard?.ar_summary.balance_amount.toLocaleString()}
                      </span>
                    </div>
                    <div className="flex justify-between items-center">
                      <span className="text-sm text-gray-500">當期未收</span>
                      <span className="font-medium text-green-600">
                        {dashboard?.ar_summary.currency} {dashboard?.ar_summary.current.toLocaleString()}
                      </span>
                    </div>
                    <div className="flex justify-between items-center">
                      <span className="text-sm text-gray-500">逾期金額</span>
                      <span className="font-medium text-red-600">
                        {dashboard?.ar_summary?.currency} {((dashboard?.ar_summary?.days_30 || 0) + (dashboard?.ar_summary?.days_60 || 0) + (dashboard?.ar_summary?.over_90 || 0)).toLocaleString()}
                      </span>
                    </div>
                  </div>
                  <Button className="w-full mt-4" onClick={() => router.push('/finance/receivables')}>
                    查看應收帳款
                  </Button>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <AlertCircle className="h-5 w-5 text-orange-600" />
                    催收提醒
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-3">
                    <div className="p-3 border-l-4 border-yellow-400 bg-yellow-50 rounded">
                      <p className="text-sm font-medium">30天內逾期</p>
                      <p className="text-xs text-gray-600">
                        {dashboard?.ar_summary.currency} {dashboard?.ar_summary.days_30.toLocaleString()}
                      </p>
                    </div>
                    <div className="p-3 border-l-4 border-orange-400 bg-orange-50 rounded">
                      <p className="text-sm font-medium">60天內逾期</p>
                      <p className="text-xs text-gray-600">
                        {dashboard?.ar_summary.currency} {dashboard?.ar_summary.days_60.toLocaleString()}
                      </p>
                    </div>
                    <div className="p-3 border-l-4 border-red-400 bg-red-50 rounded">
                      <p className="text-sm font-medium">90天以上逾期</p>
                      <p className="text-xs text-gray-600">
                        {dashboard?.ar_summary.currency} {dashboard?.ar_summary.over_90.toLocaleString()}
                      </p>
                    </div>
                  </div>
                  <Button 
                    variant="outline" 
                    className="w-full mt-4"
                    onClick={() => router.push('/finance/receivables?aging=over90days')}
                  >
                    重點催收清單
                  </Button>
                </CardContent>
              </Card>
            </div>
          </TabsContent>
        </Tabs>
      </div>
    </DashboardLayout>
  )
}