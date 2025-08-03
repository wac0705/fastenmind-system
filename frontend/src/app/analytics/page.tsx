'use client';

import { useState, useEffect } from 'react';
import { useQuery } from '@tanstack/react-query';
import {
  LineChart, Line, BarChart, Bar, PieChart, Pie, Cell,
  XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer,
  AreaChart, Area
} from 'recharts';
import {
  TrendingUp, TrendingDown, DollarSign, Package,
  Users, FileText, Clock, CheckCircle, XCircle,
  Calendar, Filter, Download
} from 'lucide-react';
import { format, subDays, startOfMonth, endOfMonth } from 'date-fns';
import { zhTW } from 'date-fns/locale';
import DashboardLayout from '@/components/layout/DashboardLayout';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { reportService } from '@/services/report.service';
import LoadingSpinner from '@/components/LoadingSpinner';

const COLORS = ['#2563eb', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6', '#6366f1'];

export default function AnalyticsPage() {
  const [dateRange, setDateRange] = useState('30days');
  const [activeTab, setActiveTab] = useState('overview');

  // Fetch dashboard data
  const { data: dashboardData, isLoading } = useQuery({
    queryKey: ['analytics-dashboard', dateRange],
    queryFn: () => reportService.getDashboardData(dateRange),
  });

  // Fetch quote trends
  const { data: quoteTrends } = useQuery({
    queryKey: ['quote-trends', dateRange],
    queryFn: () => reportService.getQuoteTrends(dateRange),
    enabled: activeTab === 'quotes',
  });

  // Fetch customer analytics
  const { data: customerAnalytics } = useQuery({
    queryKey: ['customer-analytics', dateRange],
    queryFn: () => reportService.getCustomerAnalytics(dateRange),
    enabled: activeTab === 'customers',
  });

  // Fetch process analytics
  const { data: processAnalytics } = useQuery({
    queryKey: ['process-analytics', dateRange],
    queryFn: () => reportService.getProcessAnalytics(dateRange),
    enabled: activeTab === 'processes',
  });

  const handleExportReport = async (reportType: string) => {
    try {
      const blob = await reportService.exportReport(reportType, dateRange);
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `${reportType}-report-${format(new Date(), 'yyyy-MM-dd')}.xlsx`;
      a.click();
      window.URL.revokeObjectURL(url);
    } catch (error) {
      console.error('Export failed:', error);
    }
  };

  if (isLoading) {
    return (
      <DashboardLayout>
        <div className="flex items-center justify-center h-96">
          <LoadingSpinner />
        </div>
      </DashboardLayout>
    );
  }

  const kpiData = dashboardData?.kpis || {
    totalRevenue: 0,
    totalQuotes: 0,
    conversionRate: 0,
    avgQuoteValue: 0,
    activeCustomers: 0,
    pendingApprovals: 0,
  };

  return (
    <DashboardLayout>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex justify-between items-start">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">分析儀表板</h1>
            <p className="mt-2 text-gray-600">業務數據分析與洞察</p>
          </div>
          <div className="flex gap-3">
            <Select value={dateRange} onValueChange={setDateRange}>
              <SelectTrigger className="w-[180px]">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="7days">過去 7 天</SelectItem>
                <SelectItem value="30days">過去 30 天</SelectItem>
                <SelectItem value="90days">過去 90 天</SelectItem>
                <SelectItem value="thisMonth">本月</SelectItem>
                <SelectItem value="lastMonth">上月</SelectItem>
                <SelectItem value="thisYear">今年</SelectItem>
              </SelectContent>
            </Select>
            <Button variant="outline" onClick={() => window.location.reload()}>
              <Filter className="mr-2 h-4 w-4" />
              重新整理
            </Button>
          </div>
        </div>

        {/* KPI Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-6 gap-4">
          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-600">總營收</p>
                  <p className="text-2xl font-bold mt-1">
                    ${kpiData.totalRevenue.toLocaleString()}
                  </p>
                  <p className="text-xs text-green-600 mt-1 flex items-center">
                    <TrendingUp className="h-3 w-3 mr-1" />
                    +12.5%
                  </p>
                </div>
                <DollarSign className="h-8 w-8 text-green-600" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-600">報價單數</p>
                  <p className="text-2xl font-bold mt-1">{kpiData.totalQuotes}</p>
                  <p className="text-xs text-green-600 mt-1 flex items-center">
                    <TrendingUp className="h-3 w-3 mr-1" />
                    +8.3%
                  </p>
                </div>
                <FileText className="h-8 w-8 text-blue-600" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-600">成交率</p>
                  <p className="text-2xl font-bold mt-1">
                    {kpiData.conversionRate.toFixed(1)}%
                  </p>
                  <p className="text-xs text-red-600 mt-1 flex items-center">
                    <TrendingDown className="h-3 w-3 mr-1" />
                    -2.1%
                  </p>
                </div>
                <CheckCircle className="h-8 w-8 text-green-600" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-600">平均報價額</p>
                  <p className="text-2xl font-bold mt-1">
                    ${kpiData.avgQuoteValue.toFixed(0)}
                  </p>
                  <p className="text-xs text-green-600 mt-1 flex items-center">
                    <TrendingUp className="h-3 w-3 mr-1" />
                    +5.7%
                  </p>
                </div>
                <Package className="h-8 w-8 text-purple-600" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-600">活躍客戶</p>
                  <p className="text-2xl font-bold mt-1">{kpiData.activeCustomers}</p>
                  <p className="text-xs text-green-600 mt-1 flex items-center">
                    <TrendingUp className="h-3 w-3 mr-1" />
                    +15 新增
                  </p>
                </div>
                <Users className="h-8 w-8 text-indigo-600" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-600">待審核</p>
                  <p className="text-2xl font-bold mt-1">{kpiData.pendingApprovals}</p>
                  <p className="text-xs text-gray-600 mt-1">需要處理</p>
                </div>
                <Clock className="h-8 w-8 text-yellow-600" />
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Analytics Tabs */}
        <Tabs value={activeTab} onValueChange={setActiveTab}>
          <TabsList className="grid w-full grid-cols-4">
            <TabsTrigger value="overview">總覽</TabsTrigger>
            <TabsTrigger value="quotes">報價分析</TabsTrigger>
            <TabsTrigger value="customers">客戶分析</TabsTrigger>
            <TabsTrigger value="processes">製程分析</TabsTrigger>
          </TabsList>

          <TabsContent value="overview" className="space-y-6">
            {/* Revenue Trend */}
            <Card>
              <CardHeader>
                <div className="flex justify-between items-center">
                  <div>
                    <CardTitle>營收趨勢</CardTitle>
                    <CardDescription>每日營收變化</CardDescription>
                  </div>
                  <Button size="sm" variant="outline" onClick={() => handleExportReport('revenue')}>
                    <Download className="mr-2 h-4 w-4" />
                    匯出
                  </Button>
                </div>
              </CardHeader>
              <CardContent>
                <div className="h-[300px]">
                  <ResponsiveContainer width="100%" height="100%">
                    <AreaChart data={dashboardData?.revenueTrend || []}>
                      <CartesianGrid strokeDasharray="3 3" />
                      <XAxis dataKey="date" />
                      <YAxis />
                      <Tooltip />
                      <Area 
                        type="monotone" 
                        dataKey="revenue" 
                        stroke="#2563eb" 
                        fill="#3b82f6"
                        fillOpacity={0.3}
                      />
                    </AreaChart>
                  </ResponsiveContainer>
                </div>
              </CardContent>
            </Card>

            {/* Quote Status Distribution */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <Card>
                <CardHeader>
                  <CardTitle>報價單狀態分佈</CardTitle>
                  <CardDescription>各狀態報價單數量</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="h-[300px]">
                    <ResponsiveContainer width="100%" height="100%">
                      <PieChart>
                        <Pie
                          data={dashboardData?.quoteStatusDistribution || []}
                          cx="50%"
                          cy="50%"
                          labelLine={false}
                          label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
                          outerRadius={80}
                          fill="#8884d8"
                          dataKey="value"
                        >
                          {(dashboardData?.quoteStatusDistribution || []).map((entry: any, index: number) => (
                            <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                          ))}
                        </Pie>
                        <Tooltip />
                      </PieChart>
                    </ResponsiveContainer>
                  </div>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle>部門績效</CardTitle>
                  <CardDescription>各部門本月表現</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="h-[300px]">
                    <ResponsiveContainer width="100%" height="100%">
                      <BarChart data={dashboardData?.departmentPerformance || []}>
                        <CartesianGrid strokeDasharray="3 3" />
                        <XAxis dataKey="department" />
                        <YAxis />
                        <Tooltip />
                        <Bar dataKey="quotes" fill="#2563eb" />
                        <Bar dataKey="revenue" fill="#10b981" />
                      </BarChart>
                    </ResponsiveContainer>
                  </div>
                </CardContent>
              </Card>
            </div>
          </TabsContent>

          <TabsContent value="quotes" className="space-y-6">
            {/* Quote Trends */}
            <Card>
              <CardHeader>
                <CardTitle>報價趨勢分析</CardTitle>
                <CardDescription>報價數量與金額趨勢</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="h-[400px]">
                  <ResponsiveContainer width="100%" height="100%">
                    <LineChart data={quoteTrends?.daily || []}>
                      <CartesianGrid strokeDasharray="3 3" />
                      <XAxis dataKey="date" />
                      <YAxis yAxisId="left" />
                      <YAxis yAxisId="right" orientation="right" />
                      <Tooltip />
                      <Legend />
                      <Line 
                        yAxisId="left"
                        type="monotone" 
                        dataKey="count" 
                        stroke="#2563eb" 
                        name="報價數量"
                      />
                      <Line 
                        yAxisId="right"
                        type="monotone" 
                        dataKey="value" 
                        stroke="#10b981" 
                        name="報價金額"
                      />
                    </LineChart>
                  </ResponsiveContainer>
                </div>
              </CardContent>
            </Card>

            {/* Top Products */}
            <Card>
              <CardHeader>
                <CardTitle>熱門產品</CardTitle>
                <CardDescription>報價次數最多的產品</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {quoteTrends?.topProducts?.map((product: any, index: number) => (
                    <div key={index} className="flex items-center justify-between">
                      <div className="flex items-center gap-3">
                        <div className={`w-3 h-3 rounded-full bg-${COLORS[index]}`} />
                        <span className="font-medium">{product.name}</span>
                      </div>
                      <div className="text-right">
                        <p className="font-semibold">{product.count} 次</p>
                        <p className="text-sm text-gray-600">${product.totalValue.toLocaleString()}</p>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="customers" className="space-y-6">
            {/* Customer Growth */}
            <Card>
              <CardHeader>
                <CardTitle>客戶成長趨勢</CardTitle>
                <CardDescription>新客戶與活躍客戶趨勢</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="h-[400px]">
                  <ResponsiveContainer width="100%" height="100%">
                    <AreaChart data={customerAnalytics?.growth || []}>
                      <CartesianGrid strokeDasharray="3 3" />
                      <XAxis dataKey="date" />
                      <YAxis />
                      <Tooltip />
                      <Legend />
                      <Area 
                        type="monotone" 
                        dataKey="newCustomers" 
                        stackId="1"
                        stroke="#2563eb" 
                        fill="#3b82f6"
                        name="新客戶"
                      />
                      <Area 
                        type="monotone" 
                        dataKey="activeCustomers" 
                        stackId="1"
                        stroke="#10b981" 
                        fill="#34d399"
                        name="活躍客戶"
                      />
                    </AreaChart>
                  </ResponsiveContainer>
                </div>
              </CardContent>
            </Card>

            {/* Top Customers */}
            <Card>
              <CardHeader>
                <CardTitle>重要客戶</CardTitle>
                <CardDescription>營收貢獻最高的客戶</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {customerAnalytics?.topCustomers?.map((customer: any, index: number) => (
                    <div key={index} className="flex items-center justify-between p-4 border rounded-lg">
                      <div>
                        <p className="font-medium">{customer.name}</p>
                        <p className="text-sm text-gray-600">{customer.quoteCount} 個報價</p>
                      </div>
                      <div className="text-right">
                        <p className="font-semibold">${customer.revenue.toLocaleString()}</p>
                        <p className="text-sm text-gray-600">{customer.percentage.toFixed(1)}% 佔比</p>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="processes" className="space-y-6">
            {/* Process Utilization */}
            <Card>
              <CardHeader>
                <CardTitle>製程使用率</CardTitle>
                <CardDescription>各製程設備使用情況</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="h-[400px]">
                  <ResponsiveContainer width="100%" height="100%">
                    <BarChart data={processAnalytics?.utilization || []} layout="horizontal">
                      <CartesianGrid strokeDasharray="3 3" />
                      <XAxis type="number" />
                      <YAxis dataKey="process" type="category" />
                      <Tooltip />
                      <Bar dataKey="utilization" fill="#2563eb">
                        {(processAnalytics?.utilization || []).map((entry: any, index: number) => (
                          <Cell 
                            key={`cell-${index}`} 
                            fill={entry.utilization > 80 ? '#ef4444' : entry.utilization > 60 ? '#f59e0b' : '#10b981'} 
                          />
                        ))}
                      </Bar>
                    </BarChart>
                  </ResponsiveContainer>
                </div>
              </CardContent>
            </Card>

            {/* Cost Analysis */}
            <Card>
              <CardHeader>
                <CardTitle>成本分析</CardTitle>
                <CardDescription>各類成本佔比</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="h-[300px]">
                  <ResponsiveContainer width="100%" height="100%">
                    <PieChart>
                      <Pie
                        data={processAnalytics?.costBreakdown || []}
                        cx="50%"
                        cy="50%"
                        innerRadius={60}
                        outerRadius={80}
                        fill="#8884d8"
                        paddingAngle={5}
                        dataKey="value"
                      >
                        {(processAnalytics?.costBreakdown || []).map((entry: any, index: number) => (
                          <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                        ))}
                      </Pie>
                      <Tooltip />
                      <Legend />
                    </PieChart>
                  </ResponsiveContainer>
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </div>
    </DashboardLayout>
  );
}