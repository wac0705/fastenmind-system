'use client'

import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { useRouter } from 'next/navigation'
import DashboardLayout from '@/components/layout/DashboardLayout'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { 
  Users,
  Plus,
  Search,
  Filter,
  Download,
  Upload,
  Eye,
  Edit,
  Trash2,
  Star,
  TrendingUp,
  TrendingDown,
  AlertTriangle,
  Shield,
  Building2,
  Phone,
  Mail,
  MapPin,
  Calendar,
  Package,
  FileText,
  BarChart3,
  Settings
} from 'lucide-react'
import supplierService from '@/services/supplier.service'
import { format } from 'date-fns'
import { zhTW } from 'date-fns/locale'

export default function SuppliersPage() {
  const router = useRouter()
  const [activeTab, setActiveTab] = useState('dashboard')
  const [searchQuery, setSearchQuery] = useState('')
  const [statusFilter, setStatusFilter] = useState('')
  const [typeFilter, setTypeFilter] = useState('')
  const [riskFilter, setRiskFilter] = useState('')
  const [page, setPage] = useState(1)
  const pageSize = 20

  // Fetch supplier dashboard
  const { data: dashboard, isLoading: isLoadingDashboard } = useQuery({
    queryKey: ['supplier-dashboard'],
    queryFn: () => supplierService.getSupplierDashboard(),
  })

  // Fetch suppliers list
  const { data: suppliersData, isLoading: isLoadingSuppliers, refetch: refetchSuppliers } = useQuery({
    queryKey: ['suppliers', page, searchQuery, statusFilter, typeFilter, riskFilter],
    queryFn: () => supplierService.listSuppliers({
      page,
      page_size: pageSize,
      search: searchQuery || undefined,
      status: statusFilter || undefined,
      type: typeFilter || undefined,
      risk_level: riskFilter || undefined,
    }),
  })

  // Fetch recent purchase orders
  const { data: recentPurchaseOrders } = useQuery({
    queryKey: ['recent-purchase-orders'],
    queryFn: () => supplierService.listPurchaseOrders({ page: 1, page_size: 5 }),
  })

  // Fetch pending evaluations
  const { data: pendingEvaluations } = useQuery({
    queryKey: ['pending-evaluations'],
    queryFn: () => supplierService.listSupplierEvaluations({ 
      status: 'draft,completed',
      page: 1, 
      page_size: 5 
    }),
  })

  const getStatusBadge = (status: string) => {
    const statusConfig: Record<string, { label: string; variant: any }> = {
      active: { label: '活躍', variant: 'success' },
      inactive: { label: '非活躍', variant: 'secondary' },
      suspended: { label: '暫停', variant: 'warning' },
      blacklisted: { label: '黑名單', variant: 'destructive' },
    }

    const config = statusConfig[status] || { label: status, variant: 'default' }
    return <Badge variant={config.variant as any}>{config.label}</Badge>
  }

  const getTypeBadge = (type: string) => {
    const typeConfig: Record<string, { label: string; variant: any }> = {
      manufacturer: { label: '製造商', variant: 'info' },
      distributor: { label: '經銷商', variant: 'warning' },
      service_provider: { label: '服務商', variant: 'secondary' },
      raw_material: { label: '原料商', variant: 'success' },
    }

    const config = typeConfig[type] || { label: type, variant: 'default' }
    return <Badge variant={config.variant as any}>{config.label}</Badge>
  }

  const getRiskBadge = (riskLevel: string) => {
    const riskConfig: Record<string, { label: string; variant: any; icon: any }> = {
      low: { label: '低風險', variant: 'success', icon: Shield },
      medium: { label: '中風險', variant: 'warning', icon: AlertTriangle },
      high: { label: '高風險', variant: 'destructive', icon: AlertTriangle },
      critical: { label: '極高風險', variant: 'destructive', icon: AlertTriangle },
    }

    const config = riskConfig[riskLevel] || { label: riskLevel, variant: 'default', icon: Shield }
    const Icon = config.icon
    
    return (
      <Badge variant={config.variant as any} className="flex items-center gap-1">
        <Icon className="h-3 w-3" />
        {config.label}
      </Badge>
    )
  }

  const getRatingColor = (rating: number) => {
    if (rating >= 90) return 'text-green-600'
    if (rating >= 80) return 'text-yellow-600'
    if (rating >= 70) return 'text-orange-600'
    return 'text-red-600'
  }

  const handleSearch = () => {
    setPage(1)
    refetchSuppliers()
  }

  const handleClearFilters = () => {
    setSearchQuery('')
    setStatusFilter('')
    setTypeFilter('')
    setRiskFilter('')
    setPage(1)
    refetchSuppliers()
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
            <h1 className="text-3xl font-bold text-gray-900">供應商管理</h1>
            <p className="mt-2 text-gray-600">管理供應商資料、採購訂單與績效評估</p>
          </div>
          <div className="flex gap-2">
            <Button variant="outline" onClick={() => router.push('/suppliers/evaluations')}>
              <BarChart3 className="mr-2 h-4 w-4" />
              績效評估
            </Button>
            <Button variant="outline" onClick={() => router.push('/suppliers/purchase-orders')}>
              <FileText className="mr-2 h-4 w-4" />
              採購訂單
            </Button>
            <Button onClick={() => router.push('/suppliers/new')}>
              <Plus className="mr-2 h-4 w-4" />
              新增供應商
            </Button>
          </div>
        </div>

        <Tabs value={activeTab} onValueChange={setActiveTab}>
          <TabsList>
            <TabsTrigger value="dashboard">供應商概覽</TabsTrigger>
            <TabsTrigger value="suppliers">供應商清單</TabsTrigger>
            <TabsTrigger value="performance">績效分析</TabsTrigger>
            <TabsTrigger value="risk">風險管理</TabsTrigger>
          </TabsList>

          <TabsContent value="dashboard" className="space-y-6">
            {/* Dashboard KPIs */}
            {dashboard && (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                <Card>
                  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">總供應商數</CardTitle>
                    <Users className="h-4 w-4 text-muted-foreground" />
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold">{dashboard.total_suppliers}</div>
                    <p className="text-xs text-muted-foreground">
                      活躍: {dashboard.active_suppliers} | 暫停: {dashboard.suspended_suppliers}
                    </p>
                  </CardContent>
                </Card>
                <Card>
                  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">採購訂單</CardTitle>
                    <Package className="h-4 w-4 text-muted-foreground" />
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold">{dashboard.total_purchase_orders}</div>
                    <p className="text-xs text-muted-foreground">
                      進行中: {dashboard.confirmed_orders} | 已完成: {dashboard.received_orders}
                    </p>
                  </CardContent>
                </Card>
                <Card>
                  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">採購金額</CardTitle>
                    <TrendingUp className="h-4 w-4 text-muted-foreground" />
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold">
                      ${dashboard.total_purchase_value.toLocaleString()}
                    </div>
                    <p className="text-xs text-muted-foreground">
                      待交付: ${dashboard.pending_value.toLocaleString()}
                    </p>
                  </CardContent>
                </Card>
                <Card>
                  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">高風險供應商</CardTitle>
                    <AlertTriangle className="h-4 w-4 text-muted-foreground" />
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold text-red-600">
                      {dashboard.high_risk_suppliers + dashboard.critical_risk_suppliers}
                    </div>
                    <p className="text-xs text-muted-foreground">需要特別關注</p>
                  </CardContent>
                </Card>
              </div>
            )}

            {/* Performance Overview */}
            {dashboard && (
              <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                <Card>
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                      <BarChart3 className="h-5 w-5" />
                      供應商績效概況
                    </CardTitle>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    <div className="space-y-3">
                      <div className="flex justify-between items-center">
                        <span className="text-sm">品質評分</span>
                        <div className="flex items-center gap-2">
                          <div className="w-24 bg-gray-200 rounded-full h-2">
                            <div 
                              className="bg-green-500 h-2 rounded-full" 
                              style={{ width: `${dashboard.average_quality_rating}%` }}
                            ></div>
                          </div>
                          <span className="text-sm font-medium">{dashboard.average_quality_rating.toFixed(1)}</span>
                        </div>
                      </div>
                      <div className="flex justify-between items-center">
                        <span className="text-sm">交期評分</span>
                        <div className="flex items-center gap-2">
                          <div className="w-24 bg-gray-200 rounded-full h-2">
                            <div 
                              className="bg-blue-500 h-2 rounded-full" 
                              style={{ width: `${dashboard.average_delivery_rating}%` }}
                            ></div>
                          </div>
                          <span className="text-sm font-medium">{dashboard.average_delivery_rating.toFixed(1)}</span>
                        </div>
                      </div>
                      <div className="flex justify-between items-center">
                        <span className="text-sm">服務評分</span>
                        <div className="flex items-center gap-2">
                          <div className="w-24 bg-gray-200 rounded-full h-2">
                            <div 
                              className="bg-purple-500 h-2 rounded-full" 
                              style={{ width: `${dashboard.average_service_rating}%` }}
                            ></div>
                          </div>
                          <span className="text-sm font-medium">{dashboard.average_service_rating.toFixed(1)}</span>
                        </div>
                      </div>
                    </div>
                  </CardContent>
                </Card>

                <Card>
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                      <FileText className="h-5 w-5" />
                      評估狀態
                    </CardTitle>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    <div className="grid grid-cols-2 gap-4">
                      <div className="text-center">
                        <div className="text-2xl font-bold text-blue-600">{dashboard.total_evaluations}</div>
                        <p className="text-sm text-gray-500">總評估數</p>
                      </div>
                      <div className="text-center">
                        <div className="text-2xl font-bold text-yellow-600">{dashboard.pending_evaluations}</div>
                        <p className="text-sm text-gray-500">待完成</p>
                      </div>
                      <div className="text-center">
                        <div className="text-2xl font-bold text-green-600">{dashboard.completed_evaluations}</div>
                        <p className="text-sm text-gray-500">已完成</p>
                      </div>
                      <div className="text-center">
                        <div className="text-2xl font-bold text-red-600">
                          {dashboard.high_risk_suppliers + dashboard.critical_risk_suppliers}
                        </div>
                        <p className="text-sm text-gray-500">高風險</p>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              </div>
            )}

            {/* Recent Activities */}
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              {/* Recent Purchase Orders */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center justify-between">
                    <span className="flex items-center gap-2">
                      <Package className="h-5 w-5" />
                      最近採購訂單
                    </span>
                    <Button variant="ghost" size="sm" onClick={() => router.push('/suppliers/purchase-orders')}>
                      <Eye className="h-4 w-4" />
                    </Button>
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  {recentPurchaseOrders?.data.length === 0 ? (
                    <p className="text-center text-gray-500 py-4">暫無採購訂單記錄</p>
                  ) : (
                    <div className="space-y-3">
                      {recentPurchaseOrders?.data.map((order) => (
                        <div key={order.id} className="flex items-center justify-between p-3 border rounded-lg">
                          <div>
                            <p className="font-medium">{order.order_no}</p>
                            <p className="text-sm text-gray-500">{order.supplier?.name}</p>
                            <div className="mt-1">
                              {getStatusBadge(order.status)}
                            </div>
                          </div>
                          <div className="text-right">
                            <p className="font-medium">
                              ${order.total_amount.toLocaleString()}
                            </p>
                            <p className="text-sm text-gray-500">
                              {format(new Date(order.order_date), 'MM/dd', { locale: zhTW })}
                            </p>
                          </div>
                        </div>
                      ))}
                    </div>
                  )}
                </CardContent>
              </Card>

              {/* Pending Evaluations */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center justify-between">
                    <span className="flex items-center gap-2">
                      <BarChart3 className="h-5 w-5" />
                      待完成評估
                    </span>
                    <Button variant="ghost" size="sm" onClick={() => router.push('/suppliers/evaluations')}>
                      <Eye className="h-4 w-4" />
                    </Button>
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  {pendingEvaluations?.data.length === 0 ? (
                    <p className="text-center text-gray-500 py-4">暫無待完成評估</p>
                  ) : (
                    <div className="space-y-3">
                      {pendingEvaluations?.data.map((evaluation) => (
                        <div key={evaluation.id} className="flex items-center justify-between p-3 border rounded-lg">
                          <div>
                            <p className="font-medium">{evaluation.evaluation_no}</p>
                            <p className="text-sm text-gray-500">{evaluation.supplier?.name}</p>
                            <div className="mt-1">
                              <Badge variant={evaluation.status === 'draft' ? 'secondary' : 'warning'}>
                                {evaluation.status === 'draft' ? '草稿' : '待審核'}
                              </Badge>
                            </div>
                          </div>
                          <div className="text-right">
                            <p className="font-medium">
                              {evaluation.evaluation_type === 'monthly' && '月度評估'}
                              {evaluation.evaluation_type === 'quarterly' && '季度評估'}
                              {evaluation.evaluation_type === 'annual' && '年度評估'}
                            </p>
                            <p className="text-sm text-gray-500">
                              {format(new Date(evaluation.evaluated_at), 'MM/dd', { locale: zhTW })}
                            </p>
                          </div>
                        </div>
                      ))}
                    </div>
                  )}
                </CardContent>
              </Card>
            </div>
          </TabsContent>

          <TabsContent value="suppliers" className="space-y-6">
            {/* Filters */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Filter className="h-5 w-5" />
                  篩選條件
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                  <div className="space-y-2">
                    <label className="text-sm font-medium">搜尋</label>
                    <div className="relative">
                      <Search className="absolute left-3 top-3 h-4 w-4 text-gray-400" />
                      <Input
                        placeholder="供應商名稱或編號"
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                        className="pl-10"
                        onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
                      />
                    </div>
                  </div>
                  <div className="space-y-2">
                    <label className="text-sm font-medium">狀態</label>
                    <Select value={statusFilter} onValueChange={setStatusFilter}>
                      <SelectTrigger>
                        <SelectValue placeholder="選擇狀態" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="">全部</SelectItem>
                        <SelectItem value="active">活躍</SelectItem>
                        <SelectItem value="inactive">非活躍</SelectItem>
                        <SelectItem value="suspended">暫停</SelectItem>
                        <SelectItem value="blacklisted">黑名單</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                  <div className="space-y-2">
                    <label className="text-sm font-medium">類型</label>
                    <Select value={typeFilter} onValueChange={setTypeFilter}>
                      <SelectTrigger>
                        <SelectValue placeholder="選擇類型" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="">全部</SelectItem>
                        <SelectItem value="manufacturer">製造商</SelectItem>
                        <SelectItem value="distributor">經銷商</SelectItem>
                        <SelectItem value="service_provider">服務商</SelectItem>
                        <SelectItem value="raw_material">原料商</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                  <div className="space-y-2">
                    <label className="text-sm font-medium">風險等級</label>
                    <Select value={riskFilter} onValueChange={setRiskFilter}>
                      <SelectTrigger>
                        <SelectValue placeholder="選擇風險等級" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="">全部</SelectItem>
                        <SelectItem value="low">低風險</SelectItem>
                        <SelectItem value="medium">中風險</SelectItem>
                        <SelectItem value="high">高風險</SelectItem>
                        <SelectItem value="critical">極高風險</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                </div>
                <div className="flex gap-2 mt-4">
                  <Button onClick={handleSearch}>
                    <Search className="mr-2 h-4 w-4" />
                    搜尋
                  </Button>
                  <Button variant="outline" onClick={handleClearFilters}>
                    清除篩選
                  </Button>
                  <Button variant="outline">
                    <Download className="mr-2 h-4 w-4" />
                    匯出
                  </Button>
                  <Button variant="outline">
                    <Upload className="mr-2 h-4 w-4" />
                    匯入
                  </Button>
                </div>
              </CardContent>
            </Card>

            {/* Suppliers Table */}
            <Card>
              <CardHeader>
                <CardTitle>供應商清單</CardTitle>
                <CardDescription>
                  共 {suppliersData?.total || 0} 個供應商
                </CardDescription>
              </CardHeader>
              <CardContent>
                {isLoadingSuppliers ? (
                  <div className="text-center py-8">載入中...</div>
                ) : suppliersData?.data.length === 0 ? (
                  <div className="text-center py-8">
                    <Users className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                    <p className="text-gray-500">暫無供應商資料</p>
                  </div>
                ) : (
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead>供應商資訊</TableHead>
                        <TableHead>類型</TableHead>
                        <TableHead>狀態</TableHead>
                        <TableHead>風險等級</TableHead>
                        <TableHead>品質評分</TableHead>
                        <TableHead>交期評分</TableHead>
                        <TableHead>聯絡方式</TableHead>
                        <TableHead>操作</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {suppliersData?.data.map((supplier) => (
                        <TableRow key={supplier.id} className="cursor-pointer hover:bg-gray-50">
                          <TableCell>
                            <div>
                              <div className="flex items-center gap-2">
                                <p className="font-medium">{supplier.name}</p>
                                {supplier.overall_rating >= 90 && (
                                  <Star className="h-4 w-4 text-yellow-500 fill-current" />
                                )}
                              </div>
                              <p className="text-sm text-gray-500">{supplier.supplier_no}</p>
                              {supplier.name_en && (
                                <p className="text-sm text-gray-400">{supplier.name_en}</p>
                              )}
                            </div>
                          </TableCell>
                          <TableCell>{getTypeBadge(supplier.type)}</TableCell>
                          <TableCell>{getStatusBadge(supplier.status)}</TableCell>
                          <TableCell>{getRiskBadge(supplier.risk_level)}</TableCell>
                          <TableCell>
                            <span className={getRatingColor(supplier.quality_rating)}>
                              {supplier.quality_rating.toFixed(1)}
                            </span>
                          </TableCell>
                          <TableCell>
                            <span className={getRatingColor(supplier.delivery_rating)}>
                              {supplier.delivery_rating.toFixed(1)}
                            </span>
                          </TableCell>
                          <TableCell>
                            <div className="text-sm">
                              {supplier.phone && (
                                <div className="flex items-center gap-1">
                                  <Phone className="h-3 w-3" />
                                  {supplier.phone}
                                </div>
                              )}
                              {supplier.email && (
                                <div className="flex items-center gap-1">
                                  <Mail className="h-3 w-3" />
                                  {supplier.email}
                                </div>
                              )}
                            </div>
                          </TableCell>
                          <TableCell>
                            <div className="flex items-center gap-2">
                              <Button 
                                variant="ghost" 
                                size="sm"
                                onClick={() => router.push(`/suppliers/${supplier.id}`)}
                              >
                                <Eye className="h-4 w-4" />
                              </Button>
                              <Button 
                                variant="ghost" 
                                size="sm"
                                onClick={() => router.push(`/suppliers/${supplier.id}/edit`)}
                              >
                                <Edit className="h-4 w-4" />
                              </Button>
                            </div>
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                )}

                {/* Pagination */}
                {suppliersData && suppliersData.total > pageSize && (
                  <div className="flex items-center justify-between mt-4">
                    <p className="text-sm text-gray-500">
                      顯示 {(page - 1) * pageSize + 1} 到 {Math.min(page * pageSize, suppliersData.total)} 項，
                      共 {suppliersData.total} 項
                    </p>
                    <div className="flex gap-2">
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => setPage(page - 1)}
                        disabled={page === 1}
                      >
                        上一頁
                      </Button>
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => setPage(page + 1)}
                        disabled={page * pageSize >= suppliersData.total}
                      >
                        下一頁
                      </Button>
                    </div>
                  </div>
                )}
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="performance">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <BarChart3 className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">績效分析功能開發中</p>
                  <Button className="mt-4" onClick={() => router.push('/suppliers/performance')}>
                    前往績效分析
                  </Button>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="risk">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <AlertTriangle className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">風險管理功能開發中</p>
                  <Button className="mt-4" onClick={() => router.push('/suppliers/risk-management')}>
                    前往風險管理
                  </Button>
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </div>
    </DashboardLayout>
  )
}