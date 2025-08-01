'use client'

import { useState } from 'react'
import { useParams, useRouter } from 'next/navigation'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import DashboardLayout from '@/components/layout/DashboardLayout'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Separator } from '@/components/ui/separator'
import { Label } from '@/components/ui/label'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { 
  ArrowLeft,
  FileText, 
  DollarSign, 
  Calendar,
  Building2,
  User,
  CreditCard,
  Truck,
  Receipt,
  Download,
  Edit,
  Send,
  CheckCircle,
  XCircle,
  Users,
  AlertCircle
} from 'lucide-react'
import { format } from 'date-fns'
import { zhTW } from 'date-fns/locale'
import { useToast } from '@/components/ui/use-toast'
import financeService, { Invoice, InvoiceItem, Payment } from '@/services/finance.service'
import { useAuthStore } from '@/store/auth.store'

export default function InvoiceDetailPage() {
  const params = useParams()
  const router = useRouter()
  const queryClient = useQueryClient()
  const { toast } = useToast()
  const { user } = useAuthStore()
  const invoiceId = params.id as string

  const [activeTab, setActiveTab] = useState('details')
  const [isPaymentDialogOpen, setIsPaymentDialogOpen] = useState(false)
  const [paymentAmount, setPaymentAmount] = useState(0)
  const [paymentMethod, setPaymentMethod] = useState('bank_transfer')

  // Fetch invoice details
  const { data: invoice, isLoading } = useQuery({
    queryKey: ['invoice', invoiceId],
    queryFn: () => financeService.getInvoice(invoiceId),
  })

  // Fetch invoice items
  const { data: items } = useQuery({
    queryKey: ['invoice-items', invoiceId],
    queryFn: () => financeService.getInvoiceItems(invoiceId),
    enabled: !!invoice,
  })

  // Fetch payments
  const { data: payments } = useQuery({
    queryKey: ['invoice-payments', invoiceId],
    queryFn: () => financeService.getPaymentsByInvoice(invoiceId),
    enabled: !!invoice,
  })

  // Process payment mutation
  const processPaymentMutation = useMutation({
    mutationFn: (data: any) => financeService.processPayment(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['invoice', invoiceId] })
      queryClient.invalidateQueries({ queryKey: ['invoice-payments', invoiceId] })
      toast({ title: '付款記錄已新增' })
      setIsPaymentDialogOpen(false)
      setPaymentAmount(0)
    },
    onError: (error: any) => {
      toast({
        title: '新增付款失敗',
        description: error.response?.data?.message || '新增付款記錄時發生錯誤',
        variant: 'destructive',
      })
    },
  })

  const getInvoiceStatusBadge = (status: string) => {
    const statusConfig: Record<string, { label: string; variant: any; icon: any }> = {
      draft: { label: '草稿', variant: 'secondary', icon: FileText },
      issued: { label: '已開立', variant: 'info', icon: CheckCircle },
      sent: { label: '已寄送', variant: 'warning', icon: Send },
      paid: { label: '已付款', variant: 'success', icon: CheckCircle },
      partial_paid: { label: '部分付款', variant: 'warning', icon: AlertCircle },
      overdue: { label: '逾期', variant: 'destructive', icon: XCircle },
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

  const getInvoiceTypeBadge = (type: string) => {
    const typeConfig: Record<string, { label: string; variant: any; icon: any }> = {
      sales: { label: '銷售發票', variant: 'success', icon: DollarSign },
      purchase: { label: '採購發票', variant: 'info', icon: Receipt },
      credit_note: { label: '貸項憑單', variant: 'warning', icon: FileText },
      debit_note: { label: '借項憑單', variant: 'destructive', icon: FileText },
    }

    const config = typeConfig[type] || { label: type, variant: 'default', icon: FileText }
    const Icon = config.icon
    
    return (
      <Badge variant={config.variant as any} className="flex items-center gap-1">
        <Icon className="h-3 w-3" />
        {config.label}
      </Badge>
    )
  }

  const handleProcessPayment = () => {
    if (!paymentAmount || paymentAmount <= 0) {
      toast({
        title: '請輸入有效的付款金額',
        variant: 'destructive'
      })
      return
    }

    const paymentData = {
      type: invoice?.type === 'sales' ? 'incoming' : 'outgoing',
      invoice_id: invoiceId,
      customer_id: invoice?.customer_id,
      supplier_id: invoice?.supplier_id,
      payment_date: format(new Date(), 'yyyy-MM-dd'),
      amount: paymentAmount,
      currency: invoice?.currency,
      exchange_rate: invoice?.exchange_rate || 1,
      payment_method: paymentMethod,
      notes: `${invoice?.type === 'sales' ? '收款' : '付款'} - ${invoice?.invoice_no}`,
    }

    processPaymentMutation.mutate(paymentData)
  }

  const canEdit = ['admin', 'finance', 'manager'].includes(user?.role || '') && 
                 ['draft', 'issued'].includes(invoice?.status || '')
  const canProcessPayment = ['admin', 'finance', 'manager'].includes(user?.role || '') && 
                           ['issued', 'sent', 'partial_paid'].includes(invoice?.status || '')

  if (isLoading) {
    return (
      <DashboardLayout>
        <div className="flex items-center justify-center h-64">
          <div className="text-center">載入中...</div>
        </div>
      </DashboardLayout>
    )
  }

  if (!invoice) {
    return (
      <DashboardLayout>
        <div className="flex items-center justify-center h-64">
          <div className="text-center">找不到發票</div>
        </div>
      </DashboardLayout>
    )
  }

  return (
    <DashboardLayout>
      <div className="max-w-7xl mx-auto space-y-6">
        {/* Header */}
        <div className="flex justify-between items-start">
          <div className="flex items-center gap-4">
            <Button
              variant="ghost"
              size="icon"
              onClick={() => router.push('/finance/invoices')}
            >
              <ArrowLeft className="h-4 w-4" />
            </Button>
            <div>
              <h1 className="text-3xl font-bold text-gray-900">發票詳情</h1>
              <p className="mt-1 text-gray-600">
                發票編號：{invoice.invoice_no}
              </p>
            </div>
          </div>
          <div className="flex gap-2">
            {canEdit && (
              <Button variant="outline" onClick={() => router.push(`/finance/invoices/${invoiceId}/edit`)}>
                <Edit className="mr-2 h-4 w-4" />
                編輯
              </Button>
            )}
            {canProcessPayment && (
              <Button onClick={() => setIsPaymentDialogOpen(true)}>
                <CreditCard className="mr-2 h-4 w-4" />
                {invoice.type === 'sales' ? '記錄收款' : '記錄付款'}
              </Button>
            )}
            <Button variant="outline">
              <Download className="mr-2 h-4 w-4" />
              下載 PDF
            </Button>
          </div>
        </div>

        <Tabs value={activeTab} onValueChange={setActiveTab}>
          <TabsList>
            <TabsTrigger value="details">發票詳情</TabsTrigger>
            <TabsTrigger value="items">項目明細</TabsTrigger>
            <TabsTrigger value="payments">付款記錄</TabsTrigger>
          </TabsList>

          <TabsContent value="details" className="space-y-6">
            {/* Basic Info */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center justify-between">
                  基本資訊
                  <div className="flex gap-2">
                    {getInvoiceTypeBadge(invoice.type)}
                    {getInvoiceStatusBadge(invoice.status)}
                  </div>
                </CardTitle>
              </CardHeader>
              <CardContent className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div className="space-y-4">
                  <div>
                    <Label className="text-gray-500">
                      {invoice.type === 'sales' ? '客戶' : '供應商'}
                    </Label>
                    <div className="flex items-center gap-2 mt-1">
                      {invoice.type === 'sales' ? (
                        <Users className="h-4 w-4 text-green-500" />
                      ) : (
                        <Building2 className="h-4 w-4 text-blue-500" />
                      )}
                      <span className="font-medium">
                        {invoice.customer?.name || invoice.supplier?.name}
                      </span>
                    </div>
                  </div>
                  <div>
                    <Label className="text-gray-500">開立日期</Label>
                    <div className="flex items-center gap-2 mt-1">
                      <Calendar className="h-4 w-4 text-gray-400" />
                      <span>{format(new Date(invoice.issue_date), 'yyyy/MM/dd', { locale: zhTW })}</span>
                    </div>
                  </div>
                  <div>
                    <Label className="text-gray-500">到期日期</Label>
                    <div className="flex items-center gap-2 mt-1">
                      <Calendar className="h-4 w-4 text-gray-400" />
                      <span>{format(new Date(invoice.due_date), 'yyyy/MM/dd', { locale: zhTW })}</span>
                    </div>
                  </div>
                </div>
                <div className="space-y-4">
                  <div>
                    <Label className="text-gray-500">付款條件</Label>
                    <div className="flex items-center gap-2 mt-1">
                      <CreditCard className="h-4 w-4 text-gray-400" />
                      <span>{invoice.payment_terms}</span>
                    </div>
                  </div>
                  <div>
                    <Label className="text-gray-500">付款方式</Label>
                    <div className="flex items-center gap-2 mt-1">
                      <CreditCard className="h-4 w-4 text-gray-400" />
                      <span>{invoice.payment_method}</span>
                    </div>
                  </div>
                  <div>
                    <Label className="text-gray-500">建立者</Label>
                    <div className="flex items-center gap-2 mt-1">
                      <User className="h-4 w-4 text-gray-400" />
                      <span>{invoice.creator?.name || '-'}</span>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Financial Summary */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <DollarSign className="h-5 w-5" />
                  金額資訊
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
                  <div className="text-center p-4 bg-gray-50 rounded-lg">
                    <p className="text-sm text-gray-500">小計</p>
                    <p className="text-xl font-bold mt-1">
                      {invoice.currency} {invoice.sub_total.toLocaleString()}
                    </p>
                  </div>
                  <div className="text-center p-4 bg-blue-50 rounded-lg">
                    <p className="text-sm text-gray-500">稅額</p>
                    <p className="text-xl font-bold mt-1 text-blue-600">
                      {invoice.currency} {invoice.tax_amount.toLocaleString()}
                    </p>
                  </div>
                  <div className="text-center p-4 bg-green-50 rounded-lg">
                    <p className="text-sm text-gray-500">總金額</p>
                    <p className="text-2xl font-bold mt-1 text-green-600">
                      {invoice.currency} {invoice.total_amount.toLocaleString()}
                    </p>
                  </div>
                  <div className="text-center p-4 bg-orange-50 rounded-lg">
                    <p className="text-sm text-gray-500">未付金額</p>
                    <p className="text-xl font-bold mt-1 text-orange-600">
                      {invoice.currency} {invoice.balance_amount.toLocaleString()}
                    </p>
                  </div>
                </div>

                {invoice.notes && (
                  <div className="mt-6">
                    <Label className="text-gray-500">備註</Label>
                    <p className="mt-1 text-sm text-gray-600 whitespace-pre-wrap">{invoice.notes}</p>
                  </div>
                )}
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="items" className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle>項目明細</CardTitle>
                <CardDescription>發票包含的項目詳情</CardDescription>
              </CardHeader>
              <CardContent>
                {items && items.length > 0 ? (
                  <div className="overflow-x-auto">
                    <table className="min-w-full divide-y divide-gray-200">
                      <thead className="bg-gray-50">
                        <tr>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">描述</th>
                          <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">數量</th>
                          <th className="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase">單位</th>
                          <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">單價</th>
                          <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">小計</th>
                          <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">稅額</th>
                          <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">總計</th>
                        </tr>
                      </thead>
                      <tbody className="bg-white divide-y divide-gray-200">
                        {items.map((item) => (
                          <tr key={item.id}>
                            <td className="px-6 py-4 text-sm text-gray-900">{item.description}</td>
                            <td className="px-6 py-4 text-sm text-gray-900 text-right">{item.quantity.toLocaleString()}</td>
                            <td className="px-6 py-4 text-sm text-gray-900 text-center">{item.unit}</td>
                            <td className="px-6 py-4 text-sm text-gray-900 text-right">
                              {invoice.currency} {item.unit_price.toFixed(4)}
                            </td>
                            <td className="px-6 py-4 text-sm text-gray-900 text-right">
                              {invoice.currency} {(item.quantity * item.unit_price).toFixed(2)}
                            </td>
                            <td className="px-6 py-4 text-sm text-gray-900 text-right">
                              {invoice.currency} {item.tax_amount.toFixed(2)}
                            </td>
                            <td className="px-6 py-4 text-sm font-medium text-gray-900 text-right">
                              {invoice.currency} {item.total_price.toFixed(2)}
                            </td>
                          </tr>
                        ))}
                      </tbody>
                      <tfoot className="bg-gray-50">
                        <tr>
                          <td colSpan={6} className="px-6 py-4 text-right text-sm font-medium text-gray-900">
                            總計
                          </td>
                          <td className="px-6 py-4 text-right text-sm font-bold text-gray-900">
                            {invoice.currency} {invoice.total_amount.toFixed(2)}
                          </td>
                        </tr>
                      </tfoot>
                    </table>
                  </div>
                ) : (
                  <div className="text-center py-8 text-gray-500">
                    暫無項目明細
                  </div>
                )}
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="payments" className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle>付款記錄</CardTitle>
                <CardDescription>此發票的付款歷程</CardDescription>
              </CardHeader>
              <CardContent>
                {payments && payments.length > 0 ? (
                  <div className="space-y-4">
                    {payments.map((payment) => (
                      <div key={payment.id} className="border rounded-lg p-4">
                        <div className="flex justify-between items-start">
                          <div className="space-y-1">
                            <p className="font-medium">{payment.payment_no}</p>
                            <p className="text-sm text-gray-500">
                              {format(new Date(payment.payment_date), 'yyyy/MM/dd HH:mm', { locale: zhTW })}
                            </p>
                            <p className="text-sm text-gray-500">
                              {payment.payment_method === 'bank_transfer' && '銀行轉帳'}
                              {payment.payment_method === 'cash' && '現金'}
                              {payment.payment_method === 'check' && '支票'}
                              {payment.payment_method === 'credit_card' && '信用卡'}
                            </p>
                          </div>
                          <div className="text-right">
                            <p className="text-lg font-semibold text-green-600">
                              {payment.currency} {payment.amount.toLocaleString()}
                            </p>
                            <Badge variant="success">已完成</Badge>
                          </div>
                        </div>
                        {payment.notes && (
                          <p className="text-sm text-gray-600 mt-2">{payment.notes}</p>
                        )}
                      </div>
                    ))}
                  </div>
                ) : (
                  <div className="text-center py-8 text-gray-500">
                    暫無付款記錄
                  </div>
                )}
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>

        {/* Payment Dialog */}
        <Dialog open={isPaymentDialogOpen} onOpenChange={setIsPaymentDialogOpen}>
          <DialogContent className="sm:max-w-[500px]">
            <DialogHeader>
              <DialogTitle>
                {invoice.type === 'sales' ? '記錄收款' : '記錄付款'}
              </DialogTitle>
              <DialogDescription>
                為發票 {invoice.invoice_no} {invoice.type === 'sales' ? '記錄收款' : '記錄付款'}資訊
              </DialogDescription>
            </DialogHeader>
            <div className="grid gap-4 py-4">
              <div className="grid gap-2">
                <Label htmlFor="payment-amount">付款金額</Label>
                <Input
                  id="payment-amount"
                  type="number"
                  value={paymentAmount}
                  onChange={(e) => setPaymentAmount(parseFloat(e.target.value) || 0)}
                  placeholder={`最大金額: ${invoice.balance_amount}`}
                  max={invoice.balance_amount}
                  step="0.01"
                />
                <p className="text-sm text-gray-500">
                  未付金額: {invoice.currency} {invoice.balance_amount.toLocaleString()}
                </p>
              </div>
              <div className="grid gap-2">
                <Label htmlFor="payment-method">付款方式</Label>
                <select
                  id="payment-method"
                  value={paymentMethod}
                  onChange={(e) => setPaymentMethod(e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  <option value="bank_transfer">銀行轉帳</option>
                  <option value="cash">現金</option>
                  <option value="check">支票</option>
                  <option value="credit_card">信用卡</option>
                </select>
              </div>
            </div>
            <DialogFooter>
              <Button variant="outline" onClick={() => setIsPaymentDialogOpen(false)}>
                取消
              </Button>
              <Button 
                onClick={handleProcessPayment}
                disabled={processPaymentMutation.isPending || paymentAmount <= 0}
              >
                {processPaymentMutation.isPending ? '處理中...' : '確認記錄'}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>
    </DashboardLayout>
  )
}