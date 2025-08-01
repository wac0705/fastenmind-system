'use client'

import { useEffect, useState } from 'react'
import DashboardLayout from '@/components/layout/DashboardLayout'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { FileText, Clock, CheckCircle, AlertCircle } from 'lucide-react'
import { useAuthStore } from '@/store/auth.store'

interface DashboardStats {
  pendingInquiries: number
  inProgressInquiries: number
  completedToday: number
  urgentInquiries: number
}

export default function EngineerDashboard() {
  const { user } = useAuthStore()
  const [stats, setStats] = useState<DashboardStats>({
    pendingInquiries: 5,
    inProgressInquiries: 3,
    completedToday: 2,
    urgentInquiries: 1,
  })

  return (
    <DashboardLayout>
      <div className="space-y-6">
        {/* Header */}
        <div>
          <h1 className="text-3xl font-bold text-gray-900">
            工程師工作台
          </h1>
          <p className="mt-2 text-gray-600">
            歡迎回來，{user?.full_name}！以下是您的工作概況。
          </p>
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">
                待處理詢價
              </CardTitle>
              <FileText className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{stats.pendingInquiries}</div>
              <p className="text-xs text-muted-foreground">
                需要您處理的新詢價單
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">
                處理中
              </CardTitle>
              <Clock className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{stats.inProgressInquiries}</div>
              <p className="text-xs text-muted-foreground">
                正在處理的詢價單
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">
                今日完成
              </CardTitle>
              <CheckCircle className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{stats.completedToday}</div>
              <p className="text-xs text-muted-foreground">
                今天已完成的報價
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">
                緊急詢價
              </CardTitle>
              <AlertCircle className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-red-600">
                {stats.urgentInquiries}
              </div>
              <p className="text-xs text-muted-foreground">
                需要優先處理
              </p>
            </CardContent>
          </Card>
        </div>

        {/* Recent Inquiries */}
        <Card>
          <CardHeader>
            <CardTitle>最新詢價單</CardTitle>
            <CardDescription>
              最近分派給您的詢價單
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <div className="flex items-center justify-between p-4 border rounded-lg">
                <div>
                  <p className="font-medium">INQ-2024-001</p>
                  <p className="text-sm text-gray-600">
                    M8x30 六角螺栓 - 10,000 pcs
                  </p>
                  <p className="text-xs text-gray-500">
                    客戶：台灣機械股份有限公司
                  </p>
                </div>
                <div className="text-right">
                  <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800">
                    待處理
                  </span>
                  <p className="text-xs text-gray-500 mt-1">
                    2 小時前
                  </p>
                </div>
              </div>

              <div className="flex items-center justify-between p-4 border rounded-lg">
                <div>
                  <p className="font-medium">INQ-2024-002</p>
                  <p className="text-sm text-gray-600">
                    不鏽鋼華司 - 50,000 pcs
                  </p>
                  <p className="text-xs text-gray-500">
                    客戶：Global Auto Parts GmbH
                  </p>
                </div>
                <div className="text-right">
                  <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                    處理中
                  </span>
                  <p className="text-xs text-gray-500 mt-1">
                    5 小時前
                  </p>
                </div>
              </div>

              <div className="flex items-center justify-between p-4 border rounded-lg">
                <div>
                  <p className="font-medium">INQ-2024-003</p>
                  <p className="text-sm text-gray-600">
                    自攻螺絲 - 100,000 pcs
                  </p>
                  <p className="text-xs text-gray-500">
                    客戶：American Tools Inc.
                  </p>
                </div>
                <div className="text-right">
                  <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800">
                    緊急
                  </span>
                  <p className="text-xs text-gray-500 mt-1">
                    1 天前
                  </p>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Quick Actions */}
        <Card>
          <CardHeader>
            <CardTitle>快速操作</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              <button className="p-4 text-center border rounded-lg hover:bg-gray-50 transition-colors">
                <FileText className="h-8 w-8 mx-auto mb-2 text-blue-600" />
                <p className="text-sm font-medium">查看詢價單</p>
              </button>
              <button className="p-4 text-center border rounded-lg hover:bg-gray-50 transition-colors">
                <Clock className="h-8 w-8 mx-auto mb-2 text-orange-600" />
                <p className="text-sm font-medium">處理中項目</p>
              </button>
              <button className="p-4 text-center border rounded-lg hover:bg-gray-50 transition-colors">
                <CheckCircle className="h-8 w-8 mx-auto mb-2 text-green-600" />
                <p className="text-sm font-medium">完成報價</p>
              </button>
              <button className="p-4 text-center border rounded-lg hover:bg-gray-50 transition-colors">
                <AlertCircle className="h-8 w-8 mx-auto mb-2 text-red-600" />
                <p className="text-sm font-medium">緊急任務</p>
              </button>
            </div>
          </CardContent>
        </Card>
      </div>
    </DashboardLayout>
  )
}