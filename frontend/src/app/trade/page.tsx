'use client'

import { useState, useEffect } from 'react'
import { useQuery } from '@tanstack/react-query'
import DashboardLayout from '@/components/layout/DashboardLayout'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Progress } from '@/components/ui/progress'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { 
  Globe,
  Ship,
  FileText,
  CreditCard,
  TrendingUp,
  TrendingDown,
  Package,
  Truck,
  Plane,
  DollarSign,
  AlertTriangle,
  CheckCircle,
  Clock,
  Plus,
  Search,
  Filter,
  Download,
  Upload,
  RefreshCw,
  BarChart3,
  PieChart,
  MapPin,
  Calculator,
  Shield,
  Calendar,
  Users,
  Activity,
  Zap
} from 'lucide-react'
import tradeService from '@/services/trade.service'
import { format } from 'date-fns'
import { zhTW } from 'date-fns/locale'

export default function TradePage() {
  const [activeTab, setActiveTab] = useState('overview')
  const [selectedDateRange, setSelectedDateRange] = useState('30')

  // Fetch trade statistics
  const { data: tradeStats, isLoading: isLoadingStats } = useQuery({
    queryKey: ['trade-statistics', selectedDateRange],
    queryFn: () => tradeService.getTradeStatistics({
      start_date: format(new Date(Date.now() - parseInt(selectedDateRange) * 24 * 60 * 60 * 1000), 'yyyy-MM-dd'),
      end_date: format(new Date(), 'yyyy-MM-dd'),
    }),
  })

  // Fetch recent shipments
  const { data: recentShipments } = useQuery({
    queryKey: ['recent-shipments'],
    queryFn: () => tradeService.listShipments({ limit: 5 }),
  })

  // Fetch expiring LCs
  const { data: expiringLCs } = useQuery({
    queryKey: ['expiring-lcs'],
    queryFn: () => tradeService.getExpiringLetterOfCredits(30),
  })

  // Fetch failed compliance checks
  const { data: failedChecks } = useQuery({
    queryKey: ['failed-compliance-checks'],
    queryFn: () => tradeService.getFailedComplianceChecks(),
  })

  // Fetch top trading partners
  const { data: tradingPartners } = useQuery({
    queryKey: ['top-trading-partners', selectedDateRange],
    queryFn: () => tradeService.getTopTradingPartners({
      start_date: format(new Date(Date.now() - parseInt(selectedDateRange) * 24 * 60 * 60 * 1000), 'yyyy-MM-dd'),
      end_date: format(new Date(), 'yyyy-MM-dd'),
      limit: 5,
    }),
  })

  // Fetch shipments by country
  const { data: shipmentsByCountry } = useQuery({
    queryKey: ['shipments-by-country', selectedDateRange],
    queryFn: () => tradeService.getShipmentsByCountry({
      start_date: format(new Date(Date.now() - parseInt(selectedDateRange) * 24 * 60 * 60 * 1000), 'yyyy-MM-dd'),
      end_date: format(new Date(), 'yyyy-MM-dd'),
    }),
  })

  const getShipmentStatusBadge = (status: string) => {
    const statusConfig: Record<string, { label: string; variant: any; icon: any }> = {
      pending: { label: '待處理', variant: 'secondary', icon: Clock },
      in_transit: { label: '運輸中', variant: 'info', icon: Truck },
      customs: { label: '清關中', variant: 'warning', icon: FileText },
      delivered: { label: '已送達', variant: 'success', icon: CheckCircle },
      cancelled: { label: '已取消', variant: 'destructive', icon: AlertTriangle },
    }

    const config = statusConfig[status] || { label: status, variant: 'default', icon: AlertTriangle }
    const Icon = config.icon
    
    return (
      <Badge variant={config.variant as any} className="flex items-center gap-1">
        <Icon className="h-3 w-3" />
        {config.label}
      </Badge>
    )
  }

  const getLCStatusBadge = (status: string) => {
    const statusConfig: Record<string, { label: string; variant: any }> = {
      draft: { label: '草稿', variant: 'secondary' },
      issued: { label: '已開立', variant: 'info' },
      advised: { label: '已通知', variant: 'warning' },
      confirmed: { label: '已確認', variant: 'success' },
      utilized: { label: '已使用', variant: 'default' },
      expired: { label: '已過期', variant: 'destructive' },
    }

    const config = statusConfig[status] || { label: status, variant: 'default' }
    
    return (
      <Badge variant={config.variant as any}>
        {config.label}
      </Badge>
    )
  }

  const getComplianceResultBadge = (result: string) => {
    const resultConfig: Record<string, { label: string; variant: any; icon: any }> = {
      passed: { label: '通過', variant: 'success', icon: CheckCircle },
      failed: { label: '未通過', variant: 'destructive', icon: AlertTriangle },
      warning: { label: '警告', variant: 'warning', icon: AlertTriangle },
      pending: { label: '待檢查', variant: 'secondary', icon: Clock },
    }

    const config = resultConfig[result] || { label: result, variant: 'default', icon: AlertTriangle }
    const Icon = config.icon
    
    return (
      <Badge variant={config.variant as any} className="flex items-center gap-1">
        <Icon className="h-3 w-3" />
        {config.label}
      </Badge>
    )
  }

  const getMethodIcon = (method: string) => {
    const methodIcons: Record<string, any> = {
      sea: Ship,
      air: Plane,
      land: Truck,
      express: Package,
    }
    
    const Icon = methodIcons[method] || Truck
    return <Icon className="h-4 w-4" />
  }

  if (isLoadingStats) {
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
            <h1 className="text-3xl font-bold text-gray-900">國際貿易管理</h1>
            <p className="mt-2 text-gray-600">關稅管理、運輸追蹤、信用狀與合規檢查</p>
          </div>
          <div className="flex items-center gap-4">
            <select
              value={selectedDateRange}
              onChange={(e) => setSelectedDateRange(e.target.value)}
              className="px-3 py-2 border border-gray-300 rounded-md text-sm"
            >
              <option value="7">過去 7 天</option>
              <option value="30">過去 30 天</option>
              <option value="90">過去 90 天</option>
              <option value="365">過去一年</option>
            </select>
            <Button variant="outline">
              <Download className="mr-2 h-4 w-4" />
              匯出報表
            </Button>
            <Button>
              <Plus className="mr-2 h-4 w-4" />
              新增運輸
            </Button>
          </div>
        </div>

        <Tabs value={activeTab} onValueChange={setActiveTab}>
          <TabsList className="grid w-full grid-cols-6">
            <TabsTrigger value="overview">概覽</TabsTrigger>
            <TabsTrigger value="shipments">運輸管理</TabsTrigger>
            <TabsTrigger value="tariffs">關稅管理</TabsTrigger>
            <TabsTrigger value="letter-of-credits">信用狀</TabsTrigger>
            <TabsTrigger value="compliance">合規檢查</TabsTrigger>
            <TabsTrigger value="analytics">分析報表</TabsTrigger>
          </TabsList>

          <TabsContent value="overview" className="space-y-6">
            {/* Statistics Cards */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">總運輸量</CardTitle>
                  <Ship className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">{tradeStats?.shipments?.total_shipments || 0}</div>
                  <p className="text-xs text-muted-foreground">
                    進口: {tradeStats?.shipments?.import_shipments || 0} | 出口: {tradeStats?.shipments?.export_shipments || 0}
                  </p>
                </CardContent>
              </Card>
              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">貿易總額</CardTitle>
                  <DollarSign className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">
                    {tradeService.formatCurrency(tradeStats?.shipments?.total_value || 0)}
                  </div>
                  <p className="text-xs text-muted-foreground">
                    總重量: {tradeService.formatWeight(tradeStats?.shipments?.total_weight || 0)}
                  </p>
                </CardContent>
              </Card>
              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">信用狀</CardTitle>
                  <CreditCard className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">{tradeStats?.letter_of_credits?.total_lcs || 0}</div>
                  <p className="text-xs text-muted-foreground">
                    已使用: {tradeService.formatCurrency(tradeStats?.letter_of_credits?.utilized_amount || 0)}
                  </p>
                </CardContent>
              </Card>
              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">合規檢查</CardTitle>
                  <Shield className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">{tradeStats?.compliance?.passed_checks || 0}</div>
                  <p className="text-xs text-red-600">
                    失敗: {tradeStats?.compliance?.failed_checks || 0} | 警告: {tradeStats?.compliance?.warning_checks || 0}
                  </p>
                </CardContent>
              </Card>
            </div>

            {/* Status Overview */}
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              {/* Recent Shipments */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Ship className="h-5 w-5" />
                    最近運輸
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  {recentShipments?.data && recentShipments.data.length > 0 ? (
                    recentShipments.data.map((shipment) => (
                      <div key={shipment.id} className="flex justify-between items-start">
                        <div className="flex-1">
                          <p className="font-medium">{shipment.shipment_no}</p>
                          <p className="text-sm text-gray-500">
                            {getMethodIcon(shipment.method)} {shipment.origin_country} → {shipment.dest_country}
                          </p>
                          <div className="flex items-center gap-2 mt-2">
                            {getShipmentStatusBadge(shipment.status)}
                            <Badge variant="outline" className="text-xs">
                              {tradeService.getShipmentTypeIcon(shipment.type)} {shipment.type}
                            </Badge>
                          </div>
                        </div>
                        <div className="text-right text-sm ml-4">
                          <p className="font-medium">
                            {tradeService.formatCurrency(shipment.customs_value, shipment.customs_currency)}
                          </p>
                          <p className="text-gray-500">
                            {format(new Date(shipment.created_at), 'MM/dd HH:mm', { locale: zhTW })}
                          </p>
                        </div>
                      </div>
                    ))
                  ) : (
                    <div className="text-center py-8 text-gray-500">
                      <Ship className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                      <p>暫無運輸記錄</p>
                    </div>
                  )}
                </CardContent>
              </Card>

              {/* Expiring Letter of Credits */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <AlertTriangle className="h-5 w-5 text-orange-500" />
                    即將到期信用狀
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  {expiringLCs?.data && expiringLCs.data.length > 0 ? (
                    expiringLCs.data.map((lc) => (
                      <div key={lc.id} className="flex justify-between items-start">
                        <div className="flex-1">
                          <p className="font-medium">{lc.lc_number}</p>
                          <p className="text-sm text-gray-500">{lc.beneficiary_name}</p>
                          <div className="flex items-center gap-2 mt-2">
                            {getLCStatusBadge(lc.status)}
                            <Badge variant="outline" className="text-xs">
                              {lc.type}
                            </Badge>
                          </div>
                        </div>
                        <div className="text-right text-sm ml-4">
                          <p className="font-medium">
                            {tradeService.formatCurrency(lc.available_amount, lc.currency)}
                          </p>
                          <p className="text-red-600">
                            {format(new Date(lc.expiry_date), 'yyyy/MM/dd', { locale: zhTW })}
                          </p>
                        </div>
                      </div>
                    ))
                  ) : (
                    <div className="text-center py-8 text-gray-500">
                      <CreditCard className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                      <p>暫無即將到期的信用狀</p>
                    </div>
                  )}
                </CardContent>
              </Card>
            </div>

            {/* Failed Compliance Checks */}
            {failedChecks?.data && failedChecks.data.length > 0 && (
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <AlertTriangle className="h-5 w-5 text-red-500" />
                    待處理合規問題
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    {failedChecks.data.slice(0, 5).map((check) => (
                      <div key={check.id} className="flex justify-between items-start p-4 border rounded-lg">
                        <div className="flex-1">
                          <p className="font-medium">{check.resource_type}</p>
                          <p className="text-sm text-gray-500 mt-1">
                            {check.issues ? JSON.parse(check.issues).join(', ') : '合規檢查失敗'}
                          </p>
                          <div className="flex items-center gap-2 mt-2">
                            {getComplianceResultBadge(check.result)}
                            <span className="text-xs text-gray-500">
                              得分: {check.score}/100
                            </span>
                          </div>
                        </div>
                        <div className="text-right text-sm ml-4">
                          <p className="text-gray-500">
                            {format(new Date(check.checked_at), 'MM/dd HH:mm', { locale: zhTW })}
                          </p>
                          <Button size="sm" variant="outline" className="mt-2">
                            處理
                          </Button>
                        </div>
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>
            )}

            {/* Trading Partners */}
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Globe className="h-5 w-5" />
                    主要貿易夥伴
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  {tradingPartners?.data && tradingPartners.data.length > 0 ? (
                    tradingPartners.data.map((partner, index) => (
                      <div key={`${partner.country}-${partner.type}-${index}`} className="flex justify-between items-center">
                        <div className="flex items-center gap-2">
                          <div className="w-8 h-8 bg-blue-100 rounded-full flex items-center justify-center">
                            <span className="text-xs font-medium">{index + 1}</span>
                          </div>
                          <div>
                            <p className="font-medium">{partner.country}</p>
                            <p className="text-sm text-gray-500">
                              {partner.type === 'import' ? '進口' : '出口'} • {partner.shipment_count} 次運輸
                            </p>
                          </div>
                        </div>
                        <div className="text-right">
                          <p className="font-medium">
                            {tradeService.formatCurrency(partner.total_value)}
                          </p>
                          <div className="w-20 bg-gray-200 rounded-full h-2 mt-1">
                            <div
                              className="bg-blue-500 h-2 rounded-full"
                              style={{
                                width: `${Math.min((partner.total_value / (tradingPartners.data[0]?.total_value || 1)) * 100, 100)}%`
                              }}
                            ></div>
                          </div>
                        </div>
                      </div>
                    ))
                  ) : (
                    <div className="text-center py-8 text-gray-500">
                      <Globe className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                      <p>暫無貿易夥伴數據</p>
                    </div>
                  )}
                </CardContent>
              </Card>

              {/* Shipments by Country */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <MapPin className="h-5 w-5" />
                    運輸路線統計
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  {shipmentsByCountry?.data && shipmentsByCountry.data.length > 0 ? (
                    shipmentsByCountry.data.slice(0, 5).map((route, index) => (
                      <div key={`${route.origin_country}-${route.dest_country}-${index}`} className="flex justify-between items-center">
                        <div className="flex-1">
                          <p className="font-medium">
                            {route.origin_country} → {route.dest_country}
                          </p>
                          <p className="text-sm text-gray-500">
                            {route.shipment_count} 次運輸
                          </p>
                        </div>
                        <div className="text-right">
                          <p className="font-medium">
                            {tradeService.formatCurrency(route.total_value)}
                          </p>
                        </div>
                      </div>
                    ))
                  ) : (
                    <div className="text-center py-8 text-gray-500">
                      <MapPin className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                      <p>暫無運輸路線數據</p>
                    </div>
                  )}
                </CardContent>
              </Card>
            </div>
          </TabsContent>

          <TabsContent value="shipments">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <Ship className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">運輸管理功能開發中</p>
                  <p className="text-sm text-gray-400 mt-2">包含運輸單建立、追蹤、事件記錄與文件管理</p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="tariffs">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <Calculator className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">關稅管理功能開發中</p>
                  <p className="text-sm text-gray-400 mt-2">包含 HS 代碼管理、稅率設定與關稅計算</p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="letter-of-credits">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <CreditCard className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">信用狀管理功能開發中</p>
                  <p className="text-sm text-gray-400 mt-2">包含信用狀開立、使用記錄與到期提醒</p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="compliance">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <Shield className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">合規檢查功能開發中</p>
                  <p className="text-sm text-gray-400 mt-2">包含貿易法規檢查、制裁清單比對與風險評估</p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="analytics">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <BarChart3 className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">分析報表功能開發中</p>
                  <p className="text-sm text-gray-400 mt-2">包含貿易趨勢分析、成本分析與合規報告</p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </div>
    </DashboardLayout>
  )
}