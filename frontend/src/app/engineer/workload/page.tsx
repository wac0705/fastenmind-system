'use client';

import { useQuery } from '@tanstack/react-query';
import DashboardLayout from '@/components/layout/DashboardLayout';
import assignmentService from '@/services/assignment.service';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Progress } from '@/components/ui/progress';
import { Loader2, Users, Clock, CheckCircle, TrendingUp } from 'lucide-react';

export default function EngineerWorkloadPage() {
  // 獲取工作量統計
  const { data: workloadStats = [], isLoading } = useQuery({
    queryKey: ['engineer-workload'],
    queryFn: () => assignmentService.getWorkloadStats(),
  });

  // 計算總體統計
  const totalStats = workloadStats.reduce(
    (acc, curr) => ({
      totalInquiries: acc.totalInquiries + curr.current_inquiries,
      completedToday: acc.completedToday + curr.completed_today,
      completedThisWeek: acc.completedThisWeek + curr.completed_this_week,
      completedThisMonth: acc.completedThisMonth + curr.completed_this_month,
    }),
    { totalInquiries: 0, completedToday: 0, completedThisWeek: 0, completedThisMonth: 0 }
  );

  const getWorkloadLevel = (currentInquiries: number) => {
    if (currentInquiries === 0) return { label: '空閒', color: 'text-gray-500', bgColor: 'bg-gray-100' };
    if (currentInquiries <= 3) return { label: '正常', color: 'text-green-600', bgColor: 'bg-green-100' };
    if (currentInquiries <= 6) return { label: '忙碌', color: 'text-yellow-600', bgColor: 'bg-yellow-100' };
    return { label: '超載', color: 'text-red-600', bgColor: 'bg-red-100' };
  };

  const getSkillBadgeVariant = (category: string): any => {
    const variants: Record<string, any> = {
      'screws': 'default',
      'nuts': 'secondary',
      'washers': 'outline',
      'bolts': 'default',
      'special': 'destructive',
      'custom': 'destructive',
    };
    return variants[category] || 'default';
  };

  if (isLoading) {
    return (
      <DashboardLayout>
        <div className="flex items-center justify-center h-96">
          <Loader2 className="h-8 w-8 animate-spin" />
        </div>
      </DashboardLayout>
    );
  }

  return (
    <DashboardLayout>
      <div className="space-y-6">
        {/* Header */}
        <div>
          <h1 className="text-3xl font-bold text-gray-900">工程師工作量統計</h1>
          <p className="mt-2 text-gray-600">查看所有工程師的工作負載和完成情況</p>
        </div>

        {/* 總體統計卡片 */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">進行中任務</CardTitle>
              <Users className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{totalStats.totalInquiries}</div>
              <p className="text-xs text-muted-foreground">所有工程師總計</p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">今日完成</CardTitle>
              <CheckCircle className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{totalStats.completedToday}</div>
              <p className="text-xs text-muted-foreground">詢價單處理數</p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">本週完成</CardTitle>
              <TrendingUp className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{totalStats.completedThisWeek}</div>
              <p className="text-xs text-muted-foreground">週累計處理數</p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">本月完成</CardTitle>
              <Clock className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{totalStats.completedThisMonth}</div>
              <p className="text-xs text-muted-foreground">月累計處理數</p>
            </CardContent>
          </Card>
        </div>

        {/* 工程師工作量列表 */}
        <Card>
          <CardHeader>
            <CardTitle>工程師工作量明細</CardTitle>
          </CardHeader>
          <CardContent>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>工程師</TableHead>
                  <TableHead>專長領域</TableHead>
                  <TableHead>進行中</TableHead>
                  <TableHead>工作負載</TableHead>
                  <TableHead>今日完成</TableHead>
                  <TableHead>本週完成</TableHead>
                  <TableHead>本月完成</TableHead>
                  <TableHead>平均處理時間</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {workloadStats.map((engineer) => {
                  const workloadLevel = getWorkloadLevel(engineer.current_inquiries);
                  const workloadPercentage = Math.min((engineer.current_inquiries / 10) * 100, 100);
                  
                  return (
                    <TableRow key={engineer.engineer_id}>
                      <TableCell className="font-medium">{engineer.engineer_name}</TableCell>
                      <TableCell>
                        <div className="flex flex-wrap gap-1">
                          {engineer.skill_categories.map((category) => (
                            <Badge 
                              key={category} 
                              variant={getSkillBadgeVariant(category)}
                              className="text-xs"
                            >
                              {category}
                            </Badge>
                          ))}
                        </div>
                      </TableCell>
                      <TableCell>
                        <span className="font-semibold">{engineer.current_inquiries}</span>
                      </TableCell>
                      <TableCell>
                        <div className="space-y-1">
                          <div className="flex items-center gap-2">
                            <Progress value={workloadPercentage} className="w-20 h-2" />
                            <span className={`text-xs font-medium ${workloadLevel.color}`}>
                              {workloadLevel.label}
                            </span>
                          </div>
                        </div>
                      </TableCell>
                      <TableCell>{engineer.completed_today}</TableCell>
                      <TableCell>{engineer.completed_this_week}</TableCell>
                      <TableCell>{engineer.completed_this_month}</TableCell>
                      <TableCell>
                        {engineer.average_completion_hours > 0 
                          ? `${engineer.average_completion_hours.toFixed(1)} 小時`
                          : '-'
                        }
                      </TableCell>
                    </TableRow>
                  );
                })}
              </TableBody>
            </Table>
          </CardContent>
        </Card>

        {/* 工作負載分佈圖表 */}
        <Card>
          <CardHeader>
            <CardTitle>工作負載分佈</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {workloadStats.map((engineer) => {
                const percentage = Math.min((engineer.current_inquiries / 10) * 100, 100);
                const workloadLevel = getWorkloadLevel(engineer.current_inquiries);
                
                return (
                  <div key={engineer.engineer_id} className="space-y-2">
                    <div className="flex justify-between text-sm">
                      <span className="font-medium">{engineer.engineer_name}</span>
                      <span className={workloadLevel.color}>
                        {engineer.current_inquiries} 件進行中
                      </span>
                    </div>
                    <Progress value={percentage} className="h-2" />
                  </div>
                );
              })}
            </div>
          </CardContent>
        </Card>
      </div>
    </DashboardLayout>
  );
}