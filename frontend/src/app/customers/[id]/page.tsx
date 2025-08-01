'use client'

import { useState } from 'react'
import { useParams, useRouter } from 'next/navigation'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import DashboardLayout from '@/components/layout/DashboardLayout'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { 
  Table, 
  TableBody, 
  TableCell, 
  TableHead, 
  TableHeader, 
  TableRow 
} from '@/components/ui/table'
import { 
  ArrowLeft,
  Building2,
  Mail,
  Phone,
  MapPin,
  CreditCard,
  Calendar,
  Edit,
  FileText,
  Package,
  TrendingUp,
  AlertCircle,
  DollarSign,
  Clock,
  CheckCircle,
  XCircle
} from 'lucide-react'
import { useToast } from '@/components/ui/use-toast'
import customerService from '@/services/customer.service'
import Link from 'next/link'

export default function CustomerDetailPage() {
  const params = useParams()
  const router = useRouter()
  const { toast } = useToast()
  const customerId = params.id as string

  // Fetch customer details
  const { data: customer, isLoading } = useQuery({
    queryKey: ['customer', customerId],
    queryFn: () => customerService.get(customerId),
  })

  // Fetch related data
  const { data: inquiries } = useQuery({
    queryKey: ['customer-inquiries', customerId],
    queryFn: () => customerService.getInquiries(customerId),
    enabled: !!customerId,
  })

  const { data: quotes } = useQuery({
    queryKey: ['customer-quotes', customerId],
    queryFn: () => customerService.getQuotes(customerId),
    enabled: !!customerId,
  })

  const { data: orders } = useQuery({
    queryKey: ['customer-orders', customerId],
    queryFn: () => customerService.getOrders(customerId),
    enabled: !!customerId,
  })

  const { data: creditHistory } = useQuery({
    queryKey: ['customer-credit-history', customerId],
    queryFn: () => customerService.getCreditHistory(customerId),
    enabled: !!customerId,
  })

  if (isLoading) {
    return (
      <DashboardLayout>
        <div className="flex items-center justify-center h-96">
          <div className="text-center">載入中...</div>
        </div>
      </DashboardLayout>
    )
  }

  if (!customer) {
    return (
      <DashboardLayout>
        <div className="flex items-center justify-center h-96">
          <div className="text-center">
            <AlertCircle className="h-12 w-12 text-gray-400 mx-auto mb-4" />
            <p className="text-gray-500">找不到客戶資料</p>
          </div>
        </div>
      </DashboardLayout>
    )
  }

  const getCountryFlag = (countryCode: string) => {
    const flags: Record<string, string> = {
      TW: '🇹🇼',
      CN: '🇨🇳',
      US: '🇺🇸',
      DE: '🇩🇪',
      JP: '🇯🇵',
      KR: '🇰🇷',
      VN: '🇻🇳',
      TH: '🇹🇭',
      MY: '🇲🇾',
      SG: '🇸🇬',
    }
    return flags[countryCode] || '🌍'
  }

  const getStatusBadge = (status: string) => {
    const statusConfig: Record<string, { variant: any; label: string }> = {
      draft: { variant: 'secondary', label: '草稿' },
      submitted: { variant: 'default', label: '已提交' },
      approved: { variant: 'success', label: '已核准' },
      rejected: { variant: 'destructive', label: '已拒絕' },
      pending: { variant: 'warning', label: '待處理' },
      processing: { variant: 'default', label: '處理中' },
      completed: { variant: 'success', label: '已完成' },
      cancelled: { variant: 'secondary', label: '已取消' },
    }
    const config = statusConfig[status] || { variant: 'default', label: status }
    return <Badge variant={config.variant}>{config.label}</Badge>
  }

  return (
    <DashboardLayout>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-4">
            <Button variant="ghost" size="icon" onClick={() => router.back()}>
              <ArrowLeft className="h-4 w-4" />
            </Button>
            <div>
              <h1 className="text-3xl font-bold text-gray-900 flex items-center gap-2">
                <Building2 className="h-8 w-8" />
                {customer.name}
              </h1>
              <div className="flex items-center gap-4 mt-2 text-gray-600">
                <span>{getCountryFlag(customer.country)} {customer.country}</span>
                <span>客戶代碼: {customer.customer_code}</span>
                <Badge variant={customer.is_active ? 'success' : 'secondary'}>
                  {customer.is_active ? '啟用' : '停用'}
                </Badge>
              </div>
            </div>
          </div>
          <Link href={`/customers?edit=${customer.id}`}>
            <Button>
              <Edit className="mr-2 h-4 w-4" />
              編輯資料
            </Button>
          </Link>
        </div>

        {/* Overview Cards */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">信用額度</CardTitle>
              <CreditCard className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">
                {customer.currency} {customer.credit_limit?.toLocaleString() || '0'}
              </div>
              <p className="text-xs text-muted-foreground">
                已使用 {customer.credit_used?.toLocaleString() || '0'}
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">總詢價數</CardTitle>
              <FileText className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{inquiries?.length || 0}</div>
              <p className="text-xs text-muted-foreground">
                本月新增 {inquiries?.filter(i => new Date(i.created_at).getMonth() === new Date().getMonth()).length || 0}
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">總訂單數</CardTitle>
              <Package className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{orders?.length || 0}</div>
              <p className="text-xs text-muted-foreground">
                進行中 {orders?.filter(o => o.status === 'processing').length || 0}
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">總交易額</CardTitle>
              <TrendingUp className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">
                {customer.currency} {customer.total_revenue?.toLocaleString() || '0'}
              </div>
              <p className="text-xs text-muted-foreground">
                本年度
              </p>
            </CardContent>
          </Card>
        </div>

        {/* Main Content */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Customer Information */}
          <div className="lg:col-span-1 space-y-6">
            <Card>
              <CardHeader>
                <CardTitle>基本資料</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div>
                  <p className="text-sm text-gray-500">公司全名</p>
                  <p className="font-medium">{customer.name}</p>
                  {customer.name_en && (
                    <p className="text-sm text-gray-600">{customer.name_en}</p>
                  )}
                </div>
                {customer.short_name && (
                  <div>
                    <p className="text-sm text-gray-500">簡稱</p>
                    <p className="font-medium">{customer.short_name}</p>
                  </div>
                )}
                {customer.tax_id && (
                  <div>
                    <p className="text-sm text-gray-500">統一編號/稅號</p>
                    <p className="font-medium">{customer.tax_id}</p>
                  </div>
                )}
                <div>
                  <p className="text-sm text-gray-500">付款條件</p>
                  <p className="font-medium">{customer.payment_terms || 'T/T 30 days'}</p>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>聯絡資訊</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                {customer.contact_person && (
                  <div>
                    <p className="text-sm text-gray-500">聯絡人</p>
                    <p className="font-medium">{customer.contact_person}</p>
                  </div>
                )}
                {customer.contact_email && (
                  <div className="flex items-center gap-2">
                    <Mail className="h-4 w-4 text-gray-400" />
                    <a href={`mailto:${customer.contact_email}`} className="text-blue-600 hover:underline">
                      {customer.contact_email}
                    </a>
                  </div>
                )}
                {customer.contact_phone && (
                  <div className="flex items-center gap-2">
                    <Phone className="h-4 w-4 text-gray-400" />
                    <a href={`tel:${customer.contact_phone}`} className="text-blue-600 hover:underline">
                      {customer.contact_phone}
                    </a>
                  </div>
                )}
                {customer.address && (
                  <div>
                    <p className="text-sm text-gray-500 flex items-center gap-1 mb-1">
                      <MapPin className="h-4 w-4" />
                      公司地址
                    </p>
                    <p className="font-medium">{customer.address}</p>
                  </div>
                )}
                {customer.shipping_address && (
                  <div>
                    <p className="text-sm text-gray-500 flex items-center gap-1 mb-1">
                      <MapPin className="h-4 w-4" />
                      送貨地址
                    </p>
                    <p className="font-medium">{customer.shipping_address}</p>
                  </div>
                )}
              </CardContent>
            </Card>
          </div>

          {/* Tabs for related data */}
          <div className="lg:col-span-2">
            <Tabs defaultValue="inquiries" className="w-full">
              <TabsList className="grid w-full grid-cols-4">
                <TabsTrigger value="inquiries">詢價單</TabsTrigger>
                <TabsTrigger value="quotes">報價單</TabsTrigger>
                <TabsTrigger value="orders">訂單</TabsTrigger>
                <TabsTrigger value="credit">信用記錄</TabsTrigger>
              </TabsList>

              <TabsContent value="inquiries" className="space-y-4">
                <Card>
                  <CardHeader>
                    <CardTitle>詢價記錄</CardTitle>
                    <CardDescription>最近的詢價單列表</CardDescription>
                  </CardHeader>
                  <CardContent>
                    <Table>
                      <TableHeader>
                        <TableRow>
                          <TableHead>詢價單號</TableHead>
                          <TableHead>日期</TableHead>
                          <TableHead>項目數</TableHead>
                          <TableHead>狀態</TableHead>
                        </TableRow>
                      </TableHeader>
                      <TableBody>
                        {inquiries?.map((inquiry) => (
                          <TableRow key={inquiry.id}>
                            <TableCell>
                              <Link href={`/inquiries/${inquiry.id}`} className="text-blue-600 hover:underline">
                                {inquiry.inquiry_no}
                              </Link>
                            </TableCell>
                            <TableCell>{new Date(inquiry.created_at).toLocaleDateString()}</TableCell>
                            <TableCell>{inquiry.item_count || 0}</TableCell>
                            <TableCell>{getStatusBadge(inquiry.status)}</TableCell>
                          </TableRow>
                        ))}
                      </TableBody>
                    </Table>
                  </CardContent>
                </Card>
              </TabsContent>

              <TabsContent value="quotes" className="space-y-4">
                <Card>
                  <CardHeader>
                    <CardTitle>報價記錄</CardTitle>
                    <CardDescription>最近的報價單列表</CardDescription>
                  </CardHeader>
                  <CardContent>
                    <Table>
                      <TableHeader>
                        <TableRow>
                          <TableHead>報價單號</TableHead>
                          <TableHead>日期</TableHead>
                          <TableHead>金額</TableHead>
                          <TableHead>狀態</TableHead>
                        </TableRow>
                      </TableHeader>
                      <TableBody>
                        {quotes?.map((quote) => (
                          <TableRow key={quote.id}>
                            <TableCell>
                              <Link href={`/quotes/${quote.id}`} className="text-blue-600 hover:underline">
                                {quote.quote_no}
                              </Link>
                            </TableCell>
                            <TableCell>{new Date(quote.created_at).toLocaleDateString()}</TableCell>
                            <TableCell>{quote.currency} {quote.total_amount?.toLocaleString()}</TableCell>
                            <TableCell>{getStatusBadge(quote.status)}</TableCell>
                          </TableRow>
                        ))}
                      </TableBody>
                    </Table>
                  </CardContent>
                </Card>
              </TabsContent>

              <TabsContent value="orders" className="space-y-4">
                <Card>
                  <CardHeader>
                    <CardTitle>訂單記錄</CardTitle>
                    <CardDescription>最近的訂單列表</CardDescription>
                  </CardHeader>
                  <CardContent>
                    <Table>
                      <TableHeader>
                        <TableRow>
                          <TableHead>訂單號</TableHead>
                          <TableHead>日期</TableHead>
                          <TableHead>金額</TableHead>
                          <TableHead>狀態</TableHead>
                        </TableRow>
                      </TableHeader>
                      <TableBody>
                        {orders?.map((order) => (
                          <TableRow key={order.id}>
                            <TableCell>
                              <Link href={`/orders/${order.id}`} className="text-blue-600 hover:underline">
                                {order.order_no}
                              </Link>
                            </TableCell>
                            <TableCell>{new Date(order.created_at).toLocaleDateString()}</TableCell>
                            <TableCell>{order.currency} {order.total_amount?.toLocaleString()}</TableCell>
                            <TableCell>{getStatusBadge(order.status)}</TableCell>
                          </TableRow>
                        ))}
                      </TableBody>
                    </Table>
                  </CardContent>
                </Card>
              </TabsContent>

              <TabsContent value="credit" className="space-y-4">
                <Card>
                  <CardHeader>
                    <CardTitle>信用記錄</CardTitle>
                    <CardDescription>信用額度使用與付款記錄</CardDescription>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-4">
                      <div className="grid grid-cols-2 gap-4 p-4 bg-gray-50 rounded">
                        <div>
                          <p className="text-sm text-gray-500">信用額度</p>
                          <p className="text-xl font-bold">{customer.currency} {customer.credit_limit?.toLocaleString() || '0'}</p>
                        </div>
                        <div>
                          <p className="text-sm text-gray-500">已使用額度</p>
                          <p className="text-xl font-bold">{customer.currency} {customer.credit_used?.toLocaleString() || '0'}</p>
                        </div>
                      </div>
                      
                      <Table>
                        <TableHeader>
                          <TableRow>
                            <TableHead>日期</TableHead>
                            <TableHead>類型</TableHead>
                            <TableHead>金額</TableHead>
                            <TableHead>狀態</TableHead>
                          </TableRow>
                        </TableHeader>
                        <TableBody>
                          {creditHistory?.map((record) => (
                            <TableRow key={record.id}>
                              <TableCell>{new Date(record.date).toLocaleDateString()}</TableCell>
                              <TableCell>{record.type === 'payment' ? '付款' : '使用'}</TableCell>
                              <TableCell>
                                <span className={record.type === 'payment' ? 'text-green-600' : 'text-red-600'}>
                                  {record.type === 'payment' ? '+' : '-'}{customer.currency} {record.amount?.toLocaleString()}
                                </span>
                              </TableCell>
                              <TableCell>
                                {record.status === 'completed' ? (
                                  <CheckCircle className="h-4 w-4 text-green-600" />
                                ) : record.status === 'pending' ? (
                                  <Clock className="h-4 w-4 text-yellow-600" />
                                ) : (
                                  <XCircle className="h-4 w-4 text-red-600" />
                                )}
                              </TableCell>
                            </TableRow>
                          ))}
                        </TableBody>
                      </Table>
                    </div>
                  </CardContent>
                </Card>
              </TabsContent>
            </Tabs>
          </div>
        </div>
      </div>
    </DashboardLayout>
  )
}