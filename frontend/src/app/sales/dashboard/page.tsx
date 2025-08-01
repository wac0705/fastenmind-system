'use client'

import { useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { useQuery } from '@tanstack/react-query'
import DashboardLayout from '@/components/layout/DashboardLayout'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { 
  FileText, 
  Package, 
  DollarSign, 
  TrendingUp,
  Users,
  Clock,
  CheckCircle,
  AlertCircle,
  Plus,
  ArrowRight
} from 'lucide-react'
import { useAuthStore } from '@/store/auth.store'
import inquiryService from '@/services/inquiry.service'
import quoteService from '@/services/quote.service'
import orderService from '@/services/order.service'
import { format } from 'date-fns'
import { zhTW } from 'date-fns/locale'

export default function SalesDashboardPage() {
  const router = useRouter()
  const { user } = useAuthStore()

  // Check if user is sales
  useEffect(() => {
    if (user && user.role !== 'sales') {
      router.push('/dashboard')
    }
  }, [user, router])

  // Fetch statistics
  const { data: inquiryStats } = useQuery({
    queryKey: ['inquiry-stats'],
    queryFn: () => inquiryService.list({ page: 1, page_size: 1 }),
  })

  const { data: quoteStats } = useQuery({
    queryKey: ['quote-stats'],
    queryFn: () => quoteService.list({ page: 1, page_size: 1 }),
  })

  const { data: orderStats } = useQuery({
    queryKey: ['order-stats'],
    queryFn: () => orderService.getStats(),
  })

  // Fetch recent items
  const { data: recentInquiries } = useQuery({
    queryKey: ['recent-inquiries'],
    queryFn: () => inquiryService.list({ 
      page: 1, 
      page_size: 5,
      sales_id: user?.id 
    }),
  })

  const { data: recentOrders } = useQuery({
    queryKey: ['recent-orders'],
    queryFn: () => orderService.list({ 
      page: 1, 
      page_size: 5,
      sales_id: user?.id 
    }),
  })

  return (
    <DashboardLayout>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex justify-between items-center">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">業務儀表板</h1>
            <p className="mt-2 text-gray-600">歡迎回來，{user?.full_name}</p>
          </div>
          <Button onClick={() => router.push('/inquiries/new')}>
            <Plus className="mr-2 h-4 w-4" />
            新增詢價
          </Button>
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">待處理詢價</CardTitle>
              <FileText className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{inquiryStats?.pagination?.total || 0}</div>
              <p className="text-xs text-muted-foreground">需要跟進的詢價單</p>
            </CardContent>
          </Card>
          
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">進行中報價</CardTitle>
              <Clock className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{quoteStats?.pagination?.total || 0}</div>
              <p className="text-xs text-muted-foreground">等待客戶回覆</p>
            </CardContent>
          </Card>
          
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">本月訂單</CardTitle>
              <Package className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{orderStats?.total_orders || 0}</div>
              <p className="text-xs text-muted-foreground">本月新增訂單數</p>
            </CardContent>
          </Card>
          
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">本月業績</CardTitle>
              <DollarSign className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">${orderStats?.total_revenue?.toLocaleString() || 0}</div>
              <p className="text-xs text-muted-foreground">本月總營收</p>
            </CardContent>
          </Card>
        </div>

        {/* Recent Activities */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Recent Inquiries */}
          <Card>
            <CardHeader>
              <CardTitle>最新詢價</CardTitle>
              <CardDescription>最近收到的客戶詢價</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {recentInquiries?.data?.length === 0 ? (
                  <p className="text-center text-gray-500 py-4">暫無詢價記錄</p>
                ) : (
                  recentInquiries?.data?.map((inquiry) => (
                    <div key={inquiry.id} className="flex items-center justify-between">
                      <div className="space-y-1">
                        <p className="font-medium">{inquiry.inquiry_no}</p>
                        <p className="text-sm text-gray-500">
                          {inquiry.customer?.name} • {inquiry.part_no}
                        </p>
                      </div>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => router.push(`/inquiries/${inquiry.id}`)}
                      >
                        <ArrowRight className="h-4 w-4" />
                      </Button>
                    </div>
                  ))
                )}
              </div>
            </CardContent>
          </Card>

          {/* Recent Orders */}
          <Card>
            <CardHeader>
              <CardTitle>最新訂單</CardTitle>
              <CardDescription>最近建立的訂單</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {recentOrders?.data?.length === 0 ? (
                  <p className="text-center text-gray-500 py-4">暫無訂單記錄</p>
                ) : (
                  recentOrders?.data?.map((order) => (
                    <div key={order.id} className="flex items-center justify-between">
                      <div className="space-y-1">
                        <p className="font-medium">{order.order_no}</p>
                        <p className="text-sm text-gray-500">
                          {order.customer?.name} • ${order.total_amount.toLocaleString()}
                        </p>
                      </div>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => router.push(`/orders/${order.id}`)}
                      >
                        <ArrowRight className="h-4 w-4" />
                      </Button>
                    </div>
                  ))
                )}
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Quick Actions */}
        <Card>
          <CardHeader>
            <CardTitle>快速操作</CardTitle>
            <CardDescription>常用功能快捷入口</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <Button
                variant="outline"
                className="h-auto flex-col py-4"
                onClick={() => router.push('/inquiries/new')}
              >
                <FileText className="h-8 w-8 mb-2" />
                <span>新增詢價</span>
              </Button>
              <Button
                variant="outline"
                className="h-auto flex-col py-4"
                onClick={() => router.push('/customers')}
              >
                <Users className="h-8 w-8 mb-2" />
                <span>客戶管理</span>
              </Button>
              <Button
                variant="outline"
                className="h-auto flex-col py-4"
                onClick={() => router.push('/orders')}
              >
                <Package className="h-8 w-8 mb-2" />
                <span>訂單列表</span>
              </Button>
            </div>
          </CardContent>
        </Card>
      </div>
    </DashboardLayout>
  )
}