'use client'

import { useState } from 'react'
import { useParams, useRouter } from 'next/navigation'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import DashboardLayout from '@/components/layout/DashboardLayout'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { useToast } from '@/components/ui/use-toast'
import { 
  ArrowLeft, 
  FileText, 
  Calendar, 
  Package, 
  User, 
  Building, 
  Globe,
  DollarSign,
  UserPlus,
  CheckCircle,
  AlertCircle
} from 'lucide-react'
import { format } from 'date-fns'
import { zhTW } from 'date-fns/locale'
import inquiryService from '@/services/inquiry.service'
import assignmentService from '@/services/assignment.service'
import { useAuthStore } from '@/store/auth.store'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'

export default function InquiryDetailPage() {
  const params = useParams()
  const router = useRouter()
  const queryClient = useQueryClient()
  const { toast } = useToast()
  const { user } = useAuthStore()
  const inquiryId = params.id as string

  const [isAssignDialogOpen, setIsAssignDialogOpen] = useState(false)
  const [selectedEngineerId, setSelectedEngineerId] = useState('')
  const [assignmentNotes, setAssignmentNotes] = useState('')

  // Fetch inquiry details
  const { data: inquiry, isLoading } = useQuery({
    queryKey: ['inquiry', inquiryId],
    queryFn: () => inquiryService.get(inquiryId),
  })

  // Fetch engineer suggestion
  const { data: suggestion } = useQuery({
    queryKey: ['engineer-suggestion', inquiryId],
    queryFn: () => assignmentService.suggestEngineer(inquiryId),
    enabled: !inquiry?.assigned_engineer_id && (user?.role === 'admin' || user?.role === 'manager'),
  })

  // Fetch available engineers
  const { data: engineers = [] } = useQuery({
    queryKey: ['available-engineers', inquiryId],
    queryFn: () => assignmentService.getAvailableEngineers(inquiryId),
  })

  // Assign engineer mutation
  const assignMutation = useMutation({
    mutationFn: ({ engineerId, notes }: { engineerId: string; notes?: string }) =>
      assignmentService.assignEngineer(inquiryId, engineerId, notes),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['inquiry', inquiryId] })
      toast({ title: '工程師分派成功' })
      setIsAssignDialogOpen(false)
      setSelectedEngineerId('')
      setAssignmentNotes('')
    },
    onError: (error: any) => {
      toast({
        title: '分派失敗',
        description: error.response?.data?.message || '分派工程師時發生錯誤',
        variant: 'destructive',
      })
    },
  })

  const handleAssign = () => {
    if (!selectedEngineerId) {
      toast({
        title: '請選擇工程師',
        variant: 'destructive',
      })
      return
    }
    assignMutation.mutate({ engineerId: selectedEngineerId, notes: assignmentNotes })
  }

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

  if (isLoading) {
    return (
      <DashboardLayout>
        <div className="flex items-center justify-center h-64">
          <div className="text-center">載入中...</div>
        </div>
      </DashboardLayout>
    )
  }

  if (!inquiry) {
    return (
      <DashboardLayout>
        <div className="flex items-center justify-center h-64">
          <div className="text-center">找不到詢價單</div>
        </div>
      </DashboardLayout>
    )
  }

  return (
    <DashboardLayout>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-4">
            <Button
              variant="ghost"
              size="icon"
              onClick={() => router.push('/inquiries')}
            >
              <ArrowLeft className="h-4 w-4" />
            </Button>
            <div>
              <h1 className="text-3xl font-bold text-gray-900">
                詢價單詳情
              </h1>
              <p className="mt-1 text-gray-600">
                詢價單號：{inquiry.inquiry_no}
              </p>
            </div>
          </div>
          <div className="flex items-center gap-2">
            {getStatusBadge(inquiry.status)}
            {(user?.role === 'admin' || user?.role === 'manager') && 
             !inquiry.assigned_engineer_id && (
              <Button onClick={() => setIsAssignDialogOpen(true)}>
                <UserPlus className="mr-2 h-4 w-4" />
                分派工程師
              </Button>
            )}
          </div>
        </div>

        {/* Engineer Suggestion Alert */}
        {suggestion?.suggested_engineer && !inquiry.assigned_engineer_id && (
          <Card className="border-blue-200 bg-blue-50">
            <CardHeader className="pb-3">
              <div className="flex items-center gap-2">
                <AlertCircle className="h-5 w-5 text-blue-600" />
                <CardTitle className="text-lg text-blue-900">系統建議分派</CardTitle>
              </div>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-blue-800">
                根據分派規則，建議將此詢價單分派給：
                <strong className="ml-1">{suggestion.suggested_engineer.engineer_name}</strong>
              </p>
              <p className="text-sm text-blue-700 mt-1">
                原因：{suggestion.reason}
              </p>
            </CardContent>
          </Card>
        )}

        {/* Main Content */}
        <Tabs defaultValue="details" className="space-y-4">
          <TabsList>
            <TabsTrigger value="details">詢價詳情</TabsTrigger>
            <TabsTrigger value="customer">客戶資訊</TabsTrigger>
            <TabsTrigger value="trade">交易條件</TabsTrigger>
            <TabsTrigger value="files">相關檔案</TabsTrigger>
            {inquiry.quote_id && <TabsTrigger value="quote">報價資訊</TabsTrigger>}
          </TabsList>

          {/* Details Tab */}
          <TabsContent value="details" className="space-y-4">
            <Card>
              <CardHeader>
                <CardTitle>產品資訊</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <Label className="text-sm text-gray-600">產品分類</Label>
                    <p className="font-medium">{inquiry.product_category}</p>
                  </div>
                  <div>
                    <Label className="text-sm text-gray-600">產品名稱</Label>
                    <p className="font-medium">{inquiry.product_name}</p>
                  </div>
                  <div>
                    <Label className="text-sm text-gray-600">數量</Label>
                    <p className="font-medium">
                      {inquiry.quantity.toLocaleString()} {inquiry.unit}
                    </p>
                  </div>
                  <div>
                    <Label className="text-sm text-gray-600">要求交期</Label>
                    <p className="font-medium">
                      {format(new Date(inquiry.required_date), 'yyyy年MM月dd日', {
                        locale: zhTW,
                      })}
                    </p>
                  </div>
                </div>
                {inquiry.special_requirements && (
                  <div>
                    <Label className="text-sm text-gray-600">特殊要求</Label>
                    <p className="mt-1 whitespace-pre-wrap">{inquiry.special_requirements}</p>
                  </div>
                )}
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>處理資訊</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <Label className="text-sm text-gray-600">業務人員</Label>
                    <p className="font-medium">{inquiry.sales?.full_name || '-'}</p>
                  </div>
                  <div>
                    <Label className="text-sm text-gray-600">負責工程師</Label>
                    <p className="font-medium">
                      {inquiry.assigned_engineer?.full_name || '未分派'}
                    </p>
                  </div>
                  <div>
                    <Label className="text-sm text-gray-600">建立時間</Label>
                    <p className="font-medium">
                      {format(new Date(inquiry.created_at), 'yyyy/MM/dd HH:mm', {
                        locale: zhTW,
                      })}
                    </p>
                  </div>
                  {inquiry.assigned_at && (
                    <div>
                      <Label className="text-sm text-gray-600">分派時間</Label>
                      <p className="font-medium">
                        {format(new Date(inquiry.assigned_at), 'yyyy/MM/dd HH:mm', {
                          locale: zhTW,
                        })}
                      </p>
                    </div>
                  )}
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          {/* Customer Tab */}
          <TabsContent value="customer">
            <Card>
              <CardHeader>
                <CardTitle>客戶資訊</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <Label className="text-sm text-gray-600">客戶名稱</Label>
                    <p className="font-medium">{inquiry.customer?.name}</p>
                  </div>
                  <div>
                    <Label className="text-sm text-gray-600">客戶代碼</Label>
                    <p className="font-medium">{inquiry.customer?.customer_code}</p>
                  </div>
                  <div>
                    <Label className="text-sm text-gray-600">國家</Label>
                    <p className="font-medium">{inquiry.customer?.country}</p>
                  </div>
                  <div>
                    <Label className="text-sm text-gray-600">聯絡人</Label>
                    <p className="font-medium">{inquiry.customer?.contact_person || '-'}</p>
                  </div>
                  <div>
                    <Label className="text-sm text-gray-600">聯絡電話</Label>
                    <p className="font-medium">{inquiry.customer?.contact_phone || '-'}</p>
                  </div>
                  <div>
                    <Label className="text-sm text-gray-600">聯絡信箱</Label>
                    <p className="font-medium">{inquiry.customer?.contact_email || '-'}</p>
                  </div>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          {/* Trade Terms Tab */}
          <TabsContent value="trade">
            <Card>
              <CardHeader>
                <CardTitle>交易條件</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <Label className="text-sm text-gray-600">國際貿易條件</Label>
                    <p className="font-medium">{inquiry.incoterm}</p>
                  </div>
                  <div>
                    <Label className="text-sm text-gray-600">目的港/地址</Label>
                    <p className="font-medium">
                      {inquiry.destination_port || inquiry.destination_address || '-'}
                    </p>
                  </div>
                  <div>
                    <Label className="text-sm text-gray-600">付款條件</Label>
                    <p className="font-medium">{inquiry.payment_terms || '-'}</p>
                  </div>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          {/* Files Tab */}
          <TabsContent value="files">
            <Card>
              <CardHeader>
                <CardTitle>圖紙檔案</CardTitle>
                <CardDescription>
                  共 {inquiry.drawing_files.length} 個檔案
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-2">
                  {inquiry.drawing_files.map((file, index) => (
                    <div
                      key={index}
                      className="flex items-center justify-between p-3 bg-gray-50 rounded-lg"
                    >
                      <div className="flex items-center gap-3">
                        <FileText className="h-8 w-8 text-gray-400" />
                        <div>
                          <p className="font-medium">檔案 {index + 1}</p>
                          <p className="text-sm text-gray-500">點擊下載</p>
                        </div>
                      </div>
                      <Button variant="outline" size="sm">
                        下載
                      </Button>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>

        {/* Assign Engineer Dialog */}
        <Dialog open={isAssignDialogOpen} onOpenChange={setIsAssignDialogOpen}>
          <DialogContent className="sm:max-w-[500px]">
            <DialogHeader>
              <DialogTitle>分派工程師</DialogTitle>
              <DialogDescription>
                選擇一位工程師負責處理此詢價單
              </DialogDescription>
            </DialogHeader>
            <div className="grid gap-4 py-4">
              <div className="grid gap-2">
                <Label htmlFor="engineer">選擇工程師</Label>
                <Select
                  value={selectedEngineerId}
                  onValueChange={setSelectedEngineerId}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="選擇工程師" />
                  </SelectTrigger>
                  <SelectContent>
                    {engineers.map((engineer) => (
                      <SelectItem key={engineer.engineer_id} value={engineer.engineer_id}>
                        <div className="flex items-center justify-between w-full">
                          <span>{engineer.engineer_name}</span>
                          {engineer.engineer_id === suggestion?.suggested_engineer?.engineer_id && (
                            <Badge variant="info" className="ml-2">建議</Badge>
                          )}
                        </div>
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <div className="grid gap-2">
                <Label htmlFor="notes">備註（選填）</Label>
                <Textarea
                  id="notes"
                  value={assignmentNotes}
                  onChange={(e) => setAssignmentNotes(e.target.value)}
                  placeholder="輸入分派備註..."
                  rows={3}
                />
              </div>
            </div>
            <DialogFooter>
              <Button variant="outline" onClick={() => setIsAssignDialogOpen(false)}>
                取消
              </Button>
              <Button onClick={handleAssign}>
                確認分派
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>
    </DashboardLayout>
  )
}