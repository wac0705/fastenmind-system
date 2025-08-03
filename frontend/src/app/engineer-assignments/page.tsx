'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from '@/components/ui/dialog';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { toast } from '@/components/ui/use-toast';
import { format } from 'date-fns';
import { zhTW } from 'date-fns/locale';
import { 
  Users, 
  FileText, 
  Clock, 
  AlertCircle, 
  CheckCircle2, 
  UserPlus,
  RefreshCw,
  BarChart3,
  Calendar
} from 'lucide-react';

interface EngineerAssignment {
  id: string;
  inquiry_id: string;
  inquiry_number: string;
  customer_name: string;
  product_name: string;
  engineer_id: string;
  engineer_name: string;
  assigned_by: string;
  assigned_at: string;
  status: 'pending' | 'in_progress' | 'completed' | 'cancelled';
  priority: 'low' | 'normal' | 'high' | 'urgent';
  due_date?: string;
  completed_at?: string;
  notes?: string;
}

interface EngineerAvailability {
  engineer_id: string;
  engineer_name: string;
  department: string;
  expertise: string[];
  current_load: number;
  max_load: number;
  is_available: boolean;
  expertise_match: number;
}

interface EngineerWorkload {
  engineer_id: string;
  engineer_name: string;
  total_assignments: number;
  pending: number;
  in_progress: number;
  completed: number;
  completed_on_time: number;
  overdue: number;
  avg_completion_time: number;
}

export default function EngineerAssignmentsPage() {
  const router = useRouter();
  const [assignments, setAssignments] = useState<EngineerAssignment[]>([]);
  const [engineers, setEngineers] = useState<EngineerAvailability[]>([]);
  const [workloads, setWorkloads] = useState<EngineerWorkload[]>([]);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState('assignments');
  const [filterStatus, setFilterStatus] = useState<string>('all');
  const [filterEngineer, setFilterEngineer] = useState<string>('all');
  const [searchTerm, setSearchTerm] = useState('');
  
  // 分派對話框
  const [showAssignDialog, setShowAssignDialog] = useState(false);
  const [selectedInquiry, setSelectedInquiry] = useState<any>(null);
  const [assignmentForm, setAssignmentForm] = useState({
    engineer_id: '',
    priority: 'normal',
    due_date: '',
    notes: ''
  });
  
  // 重新分派對話框
  const [showReassignDialog, setShowReassignDialog] = useState(false);
  const [selectedAssignment, setSelectedAssignment] = useState<EngineerAssignment | null>(null);
  const [reassignForm, setReassignForm] = useState({
    new_engineer_id: '',
    reason: ''
  });

  useEffect(() => {
    fetchData();
  }, [activeTab]);

  const fetchData = async () => {
    setLoading(true);
    try {
      if (activeTab === 'assignments') {
        await fetchAssignments();
      } else if (activeTab === 'workload') {
        await fetchWorkloads();
      } else if (activeTab === 'stats') {
        await fetchStats();
      }
    } catch (error) {
      console.error('Failed to fetch data:', error);
      toast({
        title: '錯誤',
        description: '載入資料失敗',
        variant: 'destructive'
      });
    } finally {
      setLoading(false);
    }
  };

  const fetchAssignments = async () => {
    const response = await fetch('/api/v1/engineer-assignments');
    if (response.ok) {
      const data = await response.json();
      setAssignments(data.data || []);
    }
  };

  const fetchWorkloads = async () => {
    const response = await fetch('/api/v1/engineer-assignments/workload');
    if (response.ok) {
      const data = await response.json();
      setWorkloads(data.data || []);
    }
  };

  const fetchStats = async () => {
    // 統計數據將在這裡獲取
  };

  const fetchAvailableEngineers = async (inquiryId: string) => {
    const response = await fetch(`/api/v1/engineer-assignments/available?inquiry_id=${inquiryId}`);
    if (response.ok) {
      const data = await response.json();
      setEngineers(data.data || []);
    }
  };

  const handleAssign = async () => {
    try {
      const response = await fetch('/api/v1/engineer-assignments/assign', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          inquiry_id: selectedInquiry.id,
          ...assignmentForm
        })
      });

      if (response.ok) {
        toast({
          title: '成功',
          description: '工程師分派成功'
        });
        setShowAssignDialog(false);
        fetchAssignments();
      } else {
        throw new Error('分派失敗');
      }
    } catch (error) {
      toast({
        title: '錯誤',
        description: '分派工程師失敗',
        variant: 'destructive'
      });
    }
  };

  const handleReassign = async () => {
    if (!selectedAssignment) return;

    try {
      const response = await fetch(`/api/v1/engineer-assignments/${selectedAssignment.id}/reassign`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(reassignForm)
      });

      if (response.ok) {
        toast({
          title: '成功',
          description: '重新分派成功'
        });
        setShowReassignDialog(false);
        fetchAssignments();
      } else {
        throw new Error('重新分派失敗');
      }
    } catch (error) {
      toast({
        title: '錯誤',
        description: '重新分派失敗',
        variant: 'destructive'
      });
    }
  };

  const handleAutoAssign = async (inquiryId: string) => {
    try {
      const response = await fetch('/api/v1/engineer-assignments/auto-assign', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          inquiry_id: inquiryId,
          rules: {
            consider_workload: true,
            consider_expertise: true,
            consider_availability: true,
            max_assignments: 10
          }
        })
      });

      if (response.ok) {
        toast({
          title: '成功',
          description: '自動分派成功'
        });
        fetchAssignments();
      } else {
        throw new Error('自動分派失敗');
      }
    } catch (error) {
      toast({
        title: '錯誤',
        description: '自動分派失敗',
        variant: 'destructive'
      });
    }
  };

  const getStatusBadge = (status: string) => {
    const statusConfig = {
      pending: { label: '待處理', variant: 'default' as const },
      in_progress: { label: '處理中', variant: 'secondary' as const },
      completed: { label: '已完成', variant: 'success' as const },
      cancelled: { label: '已取消', variant: 'destructive' as const }
    };
    
    const config = statusConfig[status as keyof typeof statusConfig] || statusConfig.pending;
    return <Badge variant={config.variant}>{config.label}</Badge>;
  };

  const getPriorityBadge = (priority: string) => {
    const priorityConfig = {
      low: { label: '低', className: 'bg-gray-100 text-gray-800' },
      normal: { label: '一般', className: 'bg-blue-100 text-blue-800' },
      high: { label: '高', className: 'bg-orange-100 text-orange-800' },
      urgent: { label: '緊急', className: 'bg-red-100 text-red-800' }
    };
    
    const config = priorityConfig[priority as keyof typeof priorityConfig] || priorityConfig.normal;
    return <Badge className={config.className}>{config.label}</Badge>;
  };

  const filteredAssignments = assignments.filter(assignment => {
    if (filterStatus !== 'all' && assignment.status !== filterStatus) return false;
    if (filterEngineer !== 'all' && assignment.engineer_id !== filterEngineer) return false;
    if (searchTerm && !assignment.inquiry_number.toLowerCase().includes(searchTerm.toLowerCase()) &&
        !assignment.customer_name.toLowerCase().includes(searchTerm.toLowerCase())) return false;
    return true;
  });

  return (
    <div className="container mx-auto p-6 space-y-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold">工程師分派管理</h1>
        <Button onClick={() => router.push('/inquiries')} variant="outline">
          <FileText className="mr-2 h-4 w-4" />
          返回詢價單
        </Button>
      </div>

      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-3">
          <TabsTrigger value="assignments">分派管理</TabsTrigger>
          <TabsTrigger value="workload">工作負載</TabsTrigger>
          <TabsTrigger value="stats">統計分析</TabsTrigger>
        </TabsList>

        <TabsContent value="assignments" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>分派列表</CardTitle>
              <div className="flex gap-4 mt-4">
                <Input
                  placeholder="搜尋詢價單號或客戶..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className="max-w-sm"
                />
                <Select value={filterStatus} onValueChange={setFilterStatus}>
                  <SelectTrigger className="w-[180px]">
                    <SelectValue placeholder="篩選狀態" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">所有狀態</SelectItem>
                    <SelectItem value="pending">待處理</SelectItem>
                    <SelectItem value="in_progress">處理中</SelectItem>
                    <SelectItem value="completed">已完成</SelectItem>
                    <SelectItem value="cancelled">已取消</SelectItem>
                  </SelectContent>
                </Select>
                <Select value={filterEngineer} onValueChange={setFilterEngineer}>
                  <SelectTrigger className="w-[180px]">
                    <SelectValue placeholder="篩選工程師" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">所有工程師</SelectItem>
                    {/* 動態載入工程師列表 */}
                  </SelectContent>
                </Select>
              </div>
            </CardHeader>
            <CardContent>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>詢價單號</TableHead>
                    <TableHead>客戶</TableHead>
                    <TableHead>產品</TableHead>
                    <TableHead>工程師</TableHead>
                    <TableHead>狀態</TableHead>
                    <TableHead>優先級</TableHead>
                    <TableHead>分派時間</TableHead>
                    <TableHead>截止日期</TableHead>
                    <TableHead>操作</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {filteredAssignments.map((assignment) => (
                    <TableRow key={assignment.id}>
                      <TableCell className="font-medium">{assignment.inquiry_number}</TableCell>
                      <TableCell>{assignment.customer_name}</TableCell>
                      <TableCell>{assignment.product_name}</TableCell>
                      <TableCell>{assignment.engineer_name}</TableCell>
                      <TableCell>{getStatusBadge(assignment.status)}</TableCell>
                      <TableCell>{getPriorityBadge(assignment.priority)}</TableCell>
                      <TableCell>
                        {format(new Date(assignment.assigned_at), 'yyyy-MM-dd HH:mm', { locale: zhTW })}
                      </TableCell>
                      <TableCell>
                        {assignment.due_date ? 
                          format(new Date(assignment.due_date), 'yyyy-MM-dd', { locale: zhTW }) : 
                          '-'
                        }
                      </TableCell>
                      <TableCell>
                        <Button
                          size="sm"
                          variant="outline"
                          onClick={() => {
                            setSelectedAssignment(assignment);
                            setShowReassignDialog(true);
                          }}
                          disabled={assignment.status === 'completed' || assignment.status === 'cancelled'}
                        >
                          <RefreshCw className="h-4 w-4" />
                        </Button>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="workload" className="space-y-4">
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            {workloads.map((workload) => (
              <Card key={workload.engineer_id}>
                <CardHeader>
                  <CardTitle className="text-lg">{workload.engineer_name}</CardTitle>
                </CardHeader>
                <CardContent className="space-y-2">
                  <div className="flex justify-between">
                    <span className="text-sm text-muted-foreground">總分派數</span>
                    <span className="font-medium">{workload.total_assignments}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-sm text-muted-foreground">待處理</span>
                    <Badge variant="default">{workload.pending}</Badge>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-sm text-muted-foreground">處理中</span>
                    <Badge variant="secondary">{workload.in_progress}</Badge>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-sm text-muted-foreground">已完成</span>
                    <Badge variant="success">{workload.completed}</Badge>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-sm text-muted-foreground">準時完成</span>
                    <span className="text-green-600">{workload.completed_on_time}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-sm text-muted-foreground">逾期</span>
                    <span className="text-red-600">{workload.overdue}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-sm text-muted-foreground">平均完成時間</span>
                    <span>{workload.avg_completion_time.toFixed(1)} 小時</span>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </TabsContent>

        <TabsContent value="stats" className="space-y-4">
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">總分派數</CardTitle>
                <FileText className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">0</div>
                <p className="text-xs text-muted-foreground">本月新增 0 筆</p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">平均處理時間</CardTitle>
                <Clock className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">0 小時</div>
                <p className="text-xs text-muted-foreground">較上月 +0%</p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">準時完成率</CardTitle>
                <CheckCircle2 className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">0%</div>
                <p className="text-xs text-muted-foreground">目標 95%</p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">活躍工程師</CardTitle>
                <Users className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">0</div>
                <p className="text-xs text-muted-foreground">共 0 位工程師</p>
              </CardContent>
            </Card>
          </div>
          
          {/* 這裡可以加入圖表展示 */}
        </TabsContent>
      </Tabs>

      {/* 分派對話框 */}
      <Dialog open={showAssignDialog} onOpenChange={setShowAssignDialog}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>分派工程師</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <div>
              <Label>選擇工程師</Label>
              <Select
                value={assignmentForm.engineer_id}
                onValueChange={(value) => setAssignmentForm({...assignmentForm, engineer_id: value})}
              >
                <SelectTrigger>
                  <SelectValue placeholder="選擇工程師" />
                </SelectTrigger>
                <SelectContent>
                  {engineers.map((engineer) => (
                    <SelectItem key={engineer.engineer_id} value={engineer.engineer_id}>
                      <div className="flex justify-between items-center w-full">
                        <span>{engineer.engineer_name}</span>
                        <div className="flex gap-2">
                          {engineer.is_available ? 
                            <Badge variant="success">可用</Badge> : 
                            <Badge variant="destructive">忙碌</Badge>
                          }
                          <span className="text-xs text-muted-foreground">
                            {engineer.current_load}/{engineer.max_load}
                          </span>
                        </div>
                      </div>
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div>
              <Label>優先級</Label>
              <Select
                value={assignmentForm.priority}
                onValueChange={(value) => setAssignmentForm({...assignmentForm, priority: value})}
              >
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="low">低</SelectItem>
                  <SelectItem value="normal">一般</SelectItem>
                  <SelectItem value="high">高</SelectItem>
                  <SelectItem value="urgent">緊急</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div>
              <Label>截止日期</Label>
              <Input
                type="date"
                value={assignmentForm.due_date}
                onChange={(e) => setAssignmentForm({...assignmentForm, due_date: e.target.value})}
              />
            </div>
            <div>
              <Label>備註</Label>
              <Textarea
                value={assignmentForm.notes}
                onChange={(e) => setAssignmentForm({...assignmentForm, notes: e.target.value})}
                placeholder="輸入備註..."
              />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setShowAssignDialog(false)}>
              取消
            </Button>
            <Button onClick={handleAssign}>
              確認分派
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* 重新分派對話框 */}
      <Dialog open={showReassignDialog} onOpenChange={setShowReassignDialog}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>重新分派工程師</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <div>
              <Label>當前工程師</Label>
              <Input value={selectedAssignment?.engineer_name || ''} disabled />
            </div>
            <div>
              <Label>新工程師</Label>
              <Select
                value={reassignForm.new_engineer_id}
                onValueChange={(value) => setReassignForm({...reassignForm, new_engineer_id: value})}
              >
                <SelectTrigger>
                  <SelectValue placeholder="選擇新工程師" />
                </SelectTrigger>
                <SelectContent>
                  {engineers.map((engineer) => (
                    <SelectItem 
                      key={engineer.engineer_id} 
                      value={engineer.engineer_id}
                      disabled={engineer.engineer_id === selectedAssignment?.engineer_id}
                    >
                      <div className="flex justify-between items-center w-full">
                        <span>{engineer.engineer_name}</span>
                        <div className="flex gap-2">
                          {engineer.is_available ? 
                            <Badge variant="success">可用</Badge> : 
                            <Badge variant="destructive">忙碌</Badge>
                          }
                          <span className="text-xs text-muted-foreground">
                            {engineer.current_load}/{engineer.max_load}
                          </span>
                        </div>
                      </div>
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div>
              <Label>重新分派原因</Label>
              <Textarea
                value={reassignForm.reason}
                onChange={(e) => setReassignForm({...reassignForm, reason: e.target.value})}
                placeholder="請輸入重新分派的原因..."
                required
              />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setShowReassignDialog(false)}>
              取消
            </Button>
            <Button onClick={handleReassign}>
              確認重新分派
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}