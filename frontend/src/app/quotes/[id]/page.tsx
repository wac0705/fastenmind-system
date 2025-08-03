'use client'

import { useState } from 'react'
import { useParams, useRouter } from 'next/navigation'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import DashboardLayout from '@/components/layout/DashboardLayout'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Separator } from '@/components/ui/separator'
import { Textarea } from '@/components/ui/textarea'
import { Label } from '@/components/ui/label'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { 
  ArrowLeft,
  FileText, 
  DollarSign, 
  Package, 
  Calendar,
  Building2,
  User,
  CheckCircle,
  XCircle,
  Clock,
  Send,
  Download,
  Edit,
  History,
  AlertCircle,
  Truck,
  Globe,
  CreditCard,
  FileDown,
  Mail,
  ShoppingCart
} from 'lucide-react'
import { format } from 'date-fns'
import { zhTW } from 'date-fns/locale'
import { useToast } from '@/components/ui/use-toast'
import quoteService from '@/services/quote.service'
import { Quote, QuoteVersion, QuoteActivityLog } from '@/types/quote'
import { useAuthStore } from '@/store/auth.store'

export default function QuoteDetailPage() {
  const params = useParams()
  const router = useRouter()
  const queryClient = useQueryClient()
  const { toast } = useToast()
  const { user } = useAuthStore()
  const quoteId = params.id as string

  const [activeTab, setActiveTab] = useState('details')
  const [isReviewDialogOpen, setIsReviewDialogOpen] = useState(false)
  const [reviewAction, setReviewAction] = useState<'approve' | 'reject' | ''>('')
  const [reviewComments, setReviewComments] = useState('')
  const [isSendDialogOpen, setIsSendDialogOpen] = useState(false)
  const [emailMessage, setEmailMessage] = useState('')

  // Fetch quote details
  const { data: quote, isLoading } = useQuery({
    queryKey: ['quote', quoteId],
    queryFn: () => quoteService.getQuote(quoteId),
  })

  // Fetch quote versions
  const { data: versions } = useQuery({
    queryKey: ['quote-versions', quoteId],
    queryFn: () => quoteService.getQuoteVersions(quoteId),
    enabled: !!quote,
  })

  // Fetch activity logs
  const { data: activityLogs } = useQuery({
    queryKey: ['quote-activities', quoteId],
    queryFn: () => quoteService.getQuoteActivityLogs(quoteId),
    enabled: !!quote,
  })

  // Submit for approval mutation
  const submitForApprovalMutation = useMutation({
    mutationFn: () => quoteService.submitForApproval(quoteId, { notes: '' }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['quote', quoteId] })
      toast({ title: '已送出審核' })
    },
    onError: (error: any) => {
      toast({
        title: '送審失敗',
        description: error.response?.data?.message || '送出審核時發生錯誤',
        variant: 'destructive',
      })
    },
  })

  // Approve mutation
  const approveMutation = useMutation({
    mutationFn: (data: { approved: boolean; notes: string }) => 
      quoteService.approveQuote(quoteId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['quote', quoteId] })
      toast({ title: reviewAction === 'approve' ? '已核准報價' : '已拒絕報價' })
      setIsReviewDialogOpen(false)
      setReviewComments('')
    },
    onError: (error: any) => {
      toast({
        title: '審核失敗',
        description: error.response?.data?.message || '審核時發生錯誤',
        variant: 'destructive',
      })
    },
  })

  // Send quote mutation
  const sendQuoteMutation = useMutation({
    mutationFn: (data: { recipient_email: string; message?: string }) => quoteService.sendQuote(quoteId, {
      recipient_email: data.recipient_email,
      recipient_name: quote?.customer?.name,
      subject: `報價單 ${quote?.quote_no}`,
      message: data.message,
      attach_pdf: true
    }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['quote', quoteId] })
      toast({ title: '報價單已發送' })
      setIsSendDialogOpen(false)
      setEmailMessage('')
    },
    onError: (error: any) => {
      toast({
        title: '發送失敗',
        description: error.response?.data?.message || '發送報價單時發生錯誤',
        variant: 'destructive',
      })
    },
  })

  // Generate PDF function (placeholder)
  const handleGeneratePDF = () => {
    // TODO: Implement PDF generation
    toast({ 
      title: 'PDF 功能尚未實作', 
      description: '此功能將在後續版本中提供',
      variant: 'destructive'
    })
  }

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

  const canEdit = (quote?.status === 'draft' || quote?.status === 'rejected') && quote?.created_by === user?.id
  const canSubmitApproval = quote?.status === 'draft' && quote?.created_by === user?.id
  const canApprove = quote?.status === 'pending_approval' && 
    ['admin', 'manager', 'engineer_lead', 'sales_manager', 'general_manager'].includes(user?.role || '')
  const canSend = quote?.status === 'approved' && ['admin', 'sales'].includes(user?.role || '')

  const handleReview = (action: 'approve' | 'reject') => {
    setReviewAction(action)
    setIsReviewDialogOpen(true)
  }

  const handleReviewSubmit = () => {
    approveMutation.mutate({ approved: reviewAction === 'approve', notes: reviewComments })
  }

  const handleSend = () => {
    setIsSendDialogOpen(true)
  }

  const handleSendSubmit = () => {
    const recipientEmail = quote?.customer?.email || ''
    sendQuoteMutation.mutate({ recipient_email: recipientEmail, message: emailMessage })
  }

  if (isLoading) {
    return (
      <DashboardLayout>
        <div className="flex items-center justify-center h-64">
          <div className="text-center">載入中...</div>
        </div>
      </DashboardLayout>
    )
  }

  if (!quote) {
    return (
      <DashboardLayout>
        <div className="flex items-center justify-center h-64">
          <div className="text-center">找不到報價單</div>
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
              onClick={() => router.push('/quotes')}
            >
              <ArrowLeft className="h-4 w-4" />
            </Button>
            <div>
              <h1 className="text-3xl font-bold text-gray-900">報價單詳情</h1>
              <p className="mt-1 text-gray-600">
                報價單號：{quote.quote_no} | 詢價單號：{quote.inquiry?.inquiry_no}
              </p>
            </div>
          </div>
          <div className="flex gap-2">
            {canEdit && (
              <Button variant="outline" onClick={() => router.push(`/quotes/${quoteId}/edit`)}>
                <Edit className="mr-2 h-4 w-4" />
                編輯
              </Button>
            )}
            {canSubmitApproval && (
              <Button onClick={() => submitForApprovalMutation.mutate()}>
                送出審核
              </Button>
            )}
            {canApprove && (
              <>
                <Button className="bg-green-600 hover:bg-green-700 text-white" onClick={() => handleReview('approve')}>
                  <CheckCircle className="mr-2 h-4 w-4" />
                  核准
                </Button>
                <Button variant="destructive" onClick={() => handleReview('reject')}>
                  <XCircle className="mr-2 h-4 w-4" />
                  拒絕
                </Button>
              </>
            )}
            {canSend && (
              <Button onClick={handleSend}>
                <Send className="mr-2 h-4 w-4" />
                發送報價
              </Button>
            )}
            <Button variant="outline" onClick={handleGeneratePDF}>
              <Download className="mr-2 h-4 w-4" />
              下載 PDF
            </Button>
            {(quote.status === 'sent' || quote.status === 'approved') && ['admin', 'sales'].includes(user?.role || '') && (
              <Button onClick={() => router.push(`/orders/new?quote_id=${quoteId}`)}>
                <ShoppingCart className="mr-2 h-4 w-4" />
                建立訂單
              </Button>
            )}
          </div>
        </div>

        {/* Status Alert */}
        {quote.status === 'expired' && (
          <Alert>
            <AlertCircle className="h-4 w-4" />
            <AlertTitle>報價已過期</AlertTitle>
            <AlertDescription>
              此報價單已於 {quote.valid_until ? format(new Date(quote.valid_until), 'yyyy/MM/dd', { locale: zhTW }) : '未知日期'} 過期
            </AlertDescription>
          </Alert>
        )}

        <Tabs value={activeTab} onValueChange={setActiveTab}>
          <TabsList>
            <TabsTrigger value="details">報價詳情</TabsTrigger>
            <TabsTrigger value="cost">成本明細</TabsTrigger>
            <TabsTrigger value="history">版本歷史</TabsTrigger>
            <TabsTrigger value="timeline">時間軸</TabsTrigger>
          </TabsList>

          <TabsContent value="details" className="space-y-6">
            {/* Basic Info */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center justify-between">
                  基本資訊
                  {getStatusBadge(quote.status)}
                </CardTitle>
              </CardHeader>
              <CardContent className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div className="space-y-4">
                  <div>
                    <Label className="text-gray-500">客戶</Label>
                    <div className="flex items-center gap-2 mt-1">
                      <Building2 className="h-4 w-4 text-gray-400" />
                      <span className="font-medium">{quote.customer?.name}</span>
                    </div>
                  </div>
                  <div>
                    <Label className="text-gray-500">建立者</Label>
                    <div className="flex items-center gap-2 mt-1">
                      <User className="h-4 w-4 text-gray-400" />
                      <span>{quote.created_by_user?.name || '-'}</span>
                    </div>
                  </div>
                </div>
                <div className="space-y-4">
                  <div>
                    <Label className="text-gray-500">有效期限</Label>
                    <div className="flex items-center gap-2 mt-1">
                      <Calendar className="h-4 w-4 text-gray-400" />
                      <span>{quote.valid_until ? format(new Date(quote.valid_until), 'yyyy/MM/dd', { locale: zhTW }) : '未設定'}</span>
                    </div>
                  </div>
                  <div>
                    <Label className="text-gray-500">交貨天數</Label>
                    <div className="flex items-center gap-2 mt-1">
                      <Truck className="h-4 w-4 text-gray-400" />
                      <span>{quote.delivery_terms || '未設定'}</span>
                    </div>
                  </div>
                  <div>
                    <Label className="text-gray-500">付款條件</Label>
                    <div className="flex items-center gap-2 mt-1">
                      <CreditCard className="h-4 w-4 text-gray-400" />
                      <span>{quote.payment_terms}</span>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Product Info */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Package className="h-5 w-5" />
                  產品資訊
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div>
                    <Label className="text-gray-500">產品名稱</Label>
                    <p className="mt-1 font-medium">{quote.items?.[0]?.product_name || '未設定'}</p>
                  </div>
                  <div>
                    <Label className="text-gray-500">產品類別</Label>
                    <p className="mt-1">未設定</p>
                  </div>
                  <div>
                    <Label className="text-gray-500">數量</Label>
                    <p className="mt-1">{quote.items?.[0]?.quantity.toLocaleString() || '0'} {quote.items?.[0]?.unit || 'PCS'}</p>
                  </div>
                  <div>
                    <Label className="text-gray-500">交貨條件</Label>
                    <div className="flex items-center gap-2 mt-1">
                      <Globe className="h-4 w-4 text-gray-400" />
                      <span>{quote.delivery_terms || '未設定'}</span>
                    </div>
                  </div>
                </div>
                {quote.remarks && (
                  <div>
                    <Label className="text-gray-500">備註</Label>
                    <p className="mt-1 text-sm text-gray-600">{quote.remarks}</p>
                  </div>
                )}
              </CardContent>
            </Card>

            {/* Pricing */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <DollarSign className="h-5 w-5" />
                  報價資訊
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                    <div className="text-center p-4 bg-gray-50 rounded-lg">
                      <p className="text-sm text-gray-500">總金額</p>
                      <p className="text-2xl font-bold mt-1">
                        ${quote.total_amount.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
                      </p>
                    </div>
                    <div className="text-center p-4 bg-gray-50 rounded-lg">
                      <p className="text-sm text-gray-500">有效天數</p>
                      <p className="text-2xl font-bold mt-1">
                        {quote.validity_days || 30} 天
                      </p>
                    </div>
                    <div className="text-center p-4 bg-blue-50 rounded-lg">
                      <p className="text-sm text-gray-500">付款條件</p>
                      <p className="text-xl font-bold mt-1 text-blue-600">
                        {quote.payment_terms || 'T/T 30 days'}
                      </p>
                    </div>
                    <div className="text-center p-4 bg-green-50 rounded-lg">
                      <p className="text-sm text-gray-500">交貨條件</p>
                      <p className="text-xl font-bold mt-1 text-green-600">
                        {quote.delivery_terms || 'FOB'}
                      </p>
                    </div>
                  </div>
                  
                  {quote.remarks && (
                    <div className="mt-4">
                      <Label className="text-gray-500">備註</Label>
                      <p className="mt-1 text-sm text-gray-600 whitespace-pre-wrap">{quote.remarks}</p>
                    </div>
                  )}
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="cost" className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle>報價項目明細</CardTitle>
                <CardDescription>當前版本的報價項目</CardDescription>
              </CardHeader>
              <CardContent>
                {versions && versions.length > 0 && versions[0].items ? (
                  <div className="space-y-6">
                    {/* Quote Items */}
                    <div className="overflow-x-auto">
                      <table className="min-w-full divide-y divide-gray-200">
                        <thead className="bg-gray-50">
                          <tr>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                              項次
                            </th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                              產品名稱
                            </th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                              規格
                            </th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                              數量
                            </th>
                            <th className="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider">
                              單位
                            </th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                              單價
                            </th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                              總價
                            </th>
                          </tr>
                        </thead>
                        <tbody className="bg-white divide-y divide-gray-200">
                          {versions[0].items.map((item) => (
                            <tr key={item.id}>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                                {item.item_no}
                              </td>
                              <td className="px-6 py-4 text-sm text-gray-900">
                                {item.product_name}
                              </td>
                              <td className="px-6 py-4 text-sm text-gray-500">
                                {item.product_specs || '-'}
                              </td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 text-right">
                                {item.quantity.toLocaleString()}
                              </td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 text-center">
                                {item.unit}
                              </td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 text-right">
                                ${item.unit_price.toFixed(4)}
                              </td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900 text-right">
                                ${item.total_price.toFixed(2)}
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
                              ${quote.total_amount.toFixed(2)}
                            </td>
                          </tr>
                        </tfoot>
                      </table>
                    </div>

                    {/* Terms */}
                    {versions[0].terms && versions[0].terms.length > 0 && (
                      <div className="mt-6">
                        <h4 className="font-medium mb-3">報價條款</h4>
                        <div className="space-y-3">
                          {versions[0].terms.sort((a, b) => a.sort_order - b.sort_order).map((term) => (
                            <div key={term.id} className="border rounded-lg p-4">
                              <h5 className="font-medium text-sm text-gray-700 mb-2">{term.term_type}</h5>
                              <p className="text-sm text-gray-600 whitespace-pre-wrap">{term.term_content}</p>
                            </div>
                          ))}
                        </div>
                      </div>
                    )}
                  </div>
                ) : (
                  <div className="text-center py-8 text-gray-500">
                    暫無報價項目資料
                  </div>
                )}
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="history" className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <History className="h-5 w-5" />
                  版本歷史
                </CardTitle>
                <CardDescription>報價單修改記錄</CardDescription>
              </CardHeader>
              <CardContent>
                {versions && versions.length > 0 ? (
                  <div className="space-y-4">
                    {versions.map((version: QuoteVersion, index: number) => (
                      <div key={version.id} className="border-l-2 border-gray-200 pl-4 pb-4">
                        <div className="flex items-center gap-2 mb-2">
                          <div className={`w-3 h-3 rounded-full -ml-[22px] ${
                            index === 0 ? 'bg-blue-500' : 'bg-gray-300'
                          }`} />
                          <span className="text-sm text-gray-500">
                            版本 {version.version_number}
                          </span>
                          {index === 0 && (
                            <Badge variant="outline" className="text-xs">當前版本</Badge>
                          )}
                        </div>
                        <div className="ml-2">
                          <p className="text-sm">
                            由 <span className="font-medium">{version.creator?.name || version.created_by}</span> 於{' '}
                            {format(new Date(version.created_at), 'yyyy/MM/dd HH:mm', { locale: zhTW })} 建立
                          </p>
                          {version.version_notes && (
                            <p className="text-sm text-gray-600 mt-1">{version.version_notes}</p>
                          )}
                          <div className="mt-2 text-sm text-gray-500">
                            版本 {version.version_number} {version.is_current ? '(當前版本)' : ''}
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                ) : (
                  <div className="text-center py-8 text-gray-500">
                    暫無版本記錄
                  </div>
                )}
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="timeline" className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle>活動時間軸</CardTitle>
                <CardDescription>報價單處理歷程</CardDescription>
              </CardHeader>
              <CardContent>
                {activityLogs && activityLogs.length > 0 ? (
                  <div className="space-y-4">
                    {activityLogs.map((log: QuoteActivityLog, index: number) => {
                      const isLast = index === activityLogs.length - 1;
                      const getActivityColor = () => {
                        switch (log.activity_type) {
                          case 'created': return 'bg-green-500';
                          case 'updated': return 'bg-blue-500';
                          case 'submitted': return 'bg-yellow-500';
                          case 'approved': return 'bg-green-500';
                          case 'rejected': return 'bg-red-500';
                          case 'sent': return 'bg-blue-500';
                          default: return 'bg-gray-500';
                        }
                      };
                      
                      return (
                        <div key={log.id} className={`border-l-2 border-gray-200 pl-4 ${!isLast ? 'pb-4' : ''}`}>
                          <div className="flex items-center gap-2 mb-2">
                            <div className={`w-3 h-3 rounded-full -ml-[22px] ${getActivityColor()}`} />
                            <span className="text-sm text-gray-500">
                              {format(new Date(log.performed_at), 'yyyy/MM/dd HH:mm', { locale: zhTW })}
                            </span>
                          </div>
                          <div className="ml-2">
                            <p className="text-sm">
                              <span className="font-medium">{log.performer?.name || '系統'}</span>{' '}
                              {log.activity_type === 'created' && '建立報價單'}
                              {log.activity_type === 'updated' && '更新報價單'}
                              {log.activity_type === 'submitted' && '送出審核'}
                              {log.activity_type === 'approved' && '核准報價單'}
                              {log.activity_type === 'rejected' && '拒絕報價單'}
                              {log.activity_type === 'sent' && '發送報價單'}
                            </p>
                            {log.activity_description && (
                              <p className="text-sm text-gray-600 mt-1">{log.activity_description}</p>
                            )}
                          </div>
                        </div>
                      );
                    })}
                  </div>
                ) : (
                  <div className="text-center py-8 text-gray-500">
                    暫無活動記錄
                  </div>
                )}
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>

        {/* Review Dialog */}
        <Dialog open={isReviewDialogOpen} onOpenChange={setIsReviewDialogOpen}>
          <DialogContent className="sm:max-w-[500px]">
            <DialogHeader>
              <DialogTitle>
                {reviewAction === 'approve' ? '核准報價單' : '拒絕報價單'}
              </DialogTitle>
              <DialogDescription>
                請提供審核意見
              </DialogDescription>
            </DialogHeader>
            <div className="grid gap-4 py-4">
              <div className="grid gap-2">
                <Label htmlFor="review-comments">審核意見</Label>
                <Textarea
                  id="review-comments"
                  value={reviewComments}
                  onChange={(e) => setReviewComments(e.target.value)}
                  placeholder={reviewAction === 'approve' ? '核准原因或建議...' : '拒絕原因或改善建議...'}
                  rows={4}
                />
              </div>
            </div>
            <DialogFooter>
              <Button variant="outline" onClick={() => setIsReviewDialogOpen(false)}>
                取消
              </Button>
              <Button 
                variant={reviewAction === 'approve' ? 'default' : 'destructive'}
                onClick={handleReviewSubmit}
              >
                確認{reviewAction === 'approve' ? '核准' : '拒絕'}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>

        {/* Send Dialog */}
        <Dialog open={isSendDialogOpen} onOpenChange={setIsSendDialogOpen}>
          <DialogContent className="sm:max-w-[600px]">
            <DialogHeader>
              <DialogTitle className="flex items-center gap-2">
                <Mail className="h-5 w-5" />
                發送報價單
              </DialogTitle>
              <DialogDescription>
                報價單將發送給：{quote.customer?.email || '未設定'}
              </DialogDescription>
            </DialogHeader>
            <div className="grid gap-4 py-4">
              <div className="grid gap-2">
                <Label htmlFor="email-message">郵件訊息（選填）</Label>
                <Textarea
                  id="email-message"
                  value={emailMessage}
                  onChange={(e) => setEmailMessage(e.target.value)}
                  placeholder="您可以在此輸入額外的訊息..."
                  rows={4}
                />
              </div>
              <Alert>
                <AlertCircle className="h-4 w-4" />
                <AlertDescription>
                  系統將自動附上報價單 PDF 檔案，並包含基本報價資訊。
                </AlertDescription>
              </Alert>
            </div>
            <DialogFooter>
              <Button variant="outline" onClick={() => setIsSendDialogOpen(false)}>
                取消
              </Button>
              <Button onClick={handleSendSubmit}>
                <Send className="mr-2 h-4 w-4" />
                發送
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>
    </DashboardLayout>
  )
}