'use client'

import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import DashboardLayout from '@/components/layout/DashboardLayout'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
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
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Separator } from '@/components/ui/separator'
import { Switch } from '@/components/ui/switch'
import { useToast } from '@/components/ui/use-toast'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { 
  Users,
  UserPlus,
  Search,
  Filter,
  Download,
  Upload,
  Edit,
  Trash2,
  Key,
  Shield,
  Mail,
  Phone,
  Building,
  MapPin,
  Calendar,
  Clock,
  CheckCircle,
  XCircle,
  AlertCircle,
  Lock,
  Unlock,
  RefreshCw,
  Eye,
  EyeOff,
  Settings,
  UserCheck,
  UserX,
  Activity,
  Globe,
  Smartphone,
  Monitor,
  History,
  Ban,
  MoreVertical
} from 'lucide-react'
import { format } from 'date-fns'
import { zhTW } from 'date-fns/locale'
import systemService from '@/services/system.service'

export default function UsersManagementPage() {
  const queryClient = useQueryClient()
  const { toast } = useToast()
  const [searchQuery, setSearchQuery] = useState('')
  const [roleFilter, setRoleFilter] = useState('')
  const [statusFilter, setStatusFilter] = useState('')
  const [departmentFilter, setDepartmentFilter] = useState('')
  const [page, setPage] = useState(1)
  const pageSize = 20

  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false)
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false)
  const [isPasswordDialogOpen, setIsPasswordDialogOpen] = useState(false)
  const [selectedUser, setSelectedUser] = useState<any>(null)
  const [activeTab, setActiveTab] = useState('list')

  const [formData, setFormData] = useState({
    username: '',
    email: '',
    full_name: '',
    full_name_en: '',
    employee_no: '',
    department: '',
    position: '',
    phone: '',
    mobile: '',
    role_ids: [] as string[],
    is_active: true,
    company_ids: [] as string[],
    language: 'zh-TW',
    timezone: 'Asia/Taipei',
    notification_settings: {
      email: true,
      sms: false,
      push: true,
    },
  })

  const [passwordData, setPasswordData] = useState({
    new_password: '',
    confirm_password: '',
  })

  // Fetch users
  const { data: usersData, isLoading, refetch } = useQuery({
    queryKey: ['users', page, searchQuery, roleFilter, statusFilter, departmentFilter],
    queryFn: () => systemService.listUsers({
      page,
      page_size: pageSize,
      search: searchQuery || undefined,
      role_id: roleFilter || undefined,
      is_active: statusFilter === '' ? undefined : statusFilter === 'active',
      department: departmentFilter || undefined,
    }),
  })

  // Fetch roles
  const { data: rolesData } = useQuery({
    queryKey: ['roles'],
    queryFn: () => systemService.listRoles(),
  })

  // Fetch user statistics
  const { data: userStats } = useQuery({
    queryKey: ['user-stats'],
    queryFn: () => systemService.getUserStatistics(),
  })

  // Fetch online users
  const { data: onlineUsers } = useQuery({
    queryKey: ['online-users'],
    queryFn: () => systemService.getOnlineUsers(),
    refetchInterval: 30000, // Refresh every 30 seconds
  })

  // Create user mutation
  const createMutation = useMutation({
    mutationFn: (data: any) => systemService.createUser(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] })
      toast({ title: '使用者建立成功' })
      setIsCreateDialogOpen(false)
      resetForm()
    },
    onError: (error: any) => {
      toast({
        title: '建立失敗',
        description: error.response?.data?.message || '建立使用者時發生錯誤',
        variant: 'destructive',
      })
    },
  })

  // Update user mutation
  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: any }) => 
      systemService.updateUser(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] })
      toast({ title: '使用者更新成功' })
      setIsEditDialogOpen(false)
      setSelectedUser(null)
    },
    onError: (error: any) => {
      toast({
        title: '更新失敗',
        description: error.response?.data?.message || '更新使用者時發生錯誤',
        variant: 'destructive',
      })
    },
  })

  // Delete user mutation
  const deleteMutation = useMutation({
    mutationFn: (id: string) => systemService.deleteUser(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] })
      toast({ title: '使用者刪除成功' })
    },
    onError: (error: any) => {
      toast({
        title: '刪除失敗',
        description: error.response?.data?.message || '刪除使用者時發生錯誤',
        variant: 'destructive',
      })
    },
  })

  // Reset password mutation
  const resetPasswordMutation = useMutation({
    mutationFn: ({ id, password }: { id: string; password: string }) =>
      systemService.resetUserPassword(id, password),
    onSuccess: () => {
      toast({ title: '密碼重設成功' })
      setIsPasswordDialogOpen(false)
      setPasswordData({ new_password: '', confirm_password: '' })
    },
    onError: (error: any) => {
      toast({
        title: '重設失敗',
        description: error.response?.data?.message || '重設密碼時發生錯誤',
        variant: 'destructive',
      })
    },
  })

  const resetForm = () => {
    setFormData({
      username: '',
      email: '',
      full_name: '',
      full_name_en: '',
      employee_no: '',
      department: '',
      position: '',
      phone: '',
      mobile: '',
      role_ids: [],
      is_active: true,
      company_ids: [],
      language: 'zh-TW',
      timezone: 'Asia/Taipei',
      notification_settings: {
        email: true,
        sms: false,
        push: true,
      },
    })
  }

  const handleSearch = () => {
    setPage(1)
    refetch()
  }

  const handleClearFilters = () => {
    setSearchQuery('')
    setRoleFilter('')
    setStatusFilter('')
    setDepartmentFilter('')
    setPage(1)
    refetch()
  }

  const handleCreate = () => {
    createMutation.mutate(formData)
  }

  const handleUpdate = () => {
    if (selectedUser) {
      updateMutation.mutate({
        id: selectedUser.id,
        data: formData,
      })
    }
  }

  const handleDelete = (user: any) => {
    if (confirm(`確定要刪除使用者「${user.full_name}」嗎？`)) {
      deleteMutation.mutate(user.id)
    }
  }

  const handleEdit = (user: any) => {
    setSelectedUser(user)
    setFormData({
      username: user.username,
      email: user.email,
      full_name: user.full_name,
      full_name_en: user.full_name_en || '',
      employee_no: user.employee_no || '',
      department: user.department || '',
      position: user.position || '',
      phone: user.phone || '',
      mobile: user.mobile || '',
      role_ids: user.roles?.map((r: any) => r.id) || [],
      is_active: user.is_active,
      company_ids: user.companies?.map((c: any) => c.id) || [],
      language: user.language || 'zh-TW',
      timezone: user.timezone || 'Asia/Taipei',
      notification_settings: user.notification_settings || {
        email: true,
        sms: false,
        push: true,
      },
    })
    setIsEditDialogOpen(true)
  }

  const handleResetPassword = () => {
    if (passwordData.new_password !== passwordData.confirm_password) {
      toast({
        title: '密碼不符',
        description: '新密碼與確認密碼不相同',
        variant: 'destructive',
      })
      return
    }

    if (selectedUser) {
      resetPasswordMutation.mutate({
        id: selectedUser.id,
        password: passwordData.new_password,
      })
    }
  }

  const handleToggleUserStatus = async (user: any) => {
    try {
      await systemService.updateUser(user.id, { is_active: !user.is_active })
      refetch()
      toast({ 
        title: user.is_active ? '使用者已停用' : '使用者已啟用' 
      })
    } catch (error) {
      toast({
        title: '操作失敗',
        variant: 'destructive',
      })
    }
  }

  const handleExport = async () => {
    try {
      const blob = await systemService.exportUsers({
        format: 'csv',
        role_id: roleFilter || undefined,
        is_active: statusFilter === '' ? undefined : statusFilter === 'active',
        department: departmentFilter || undefined,
      })
      const url = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `users_${format(new Date(), 'yyyyMMdd_HHmmss')}.csv`
      document.body.appendChild(a)
      a.click()
      window.URL.revokeObjectURL(url)
      document.body.removeChild(a)
      toast({ title: '匯出成功' })
    } catch (error) {
      toast({
        title: '匯出失敗',
        variant: 'destructive',
      })
    }
  }

  const getStatusBadge = (isActive: boolean) => {
    if (isActive) {
      return (
        <Badge variant="success" className="flex items-center gap-1">
          <CheckCircle className="h-3 w-3" />
          啟用
        </Badge>
      )
    }
    return (
      <Badge variant="secondary" className="flex items-center gap-1">
        <XCircle className="h-3 w-3" />
        停用
      </Badge>
    )
  }

  const getOnlineStatusBadge = (userId: string) => {
    const isOnline = onlineUsers?.some((u: any) => u.id === userId)
    if (isOnline) {
      return (
        <div className="flex items-center gap-1">
          <div className="w-2 h-2 bg-green-500 rounded-full animate-pulse" />
          <span className="text-xs text-green-600">線上</span>
        </div>
      )
    }
    return null
  }

  return (
    <DashboardLayout>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex justify-between items-center">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">使用者管理</h1>
            <p className="mt-2 text-gray-600">管理系統使用者、角色權限與存取控制</p>
          </div>
          <div className="flex gap-2">
            <Button variant="outline" onClick={handleExport}>
              <Download className="mr-2 h-4 w-4" />
              匯出
            </Button>
            <Button onClick={() => setIsCreateDialogOpen(true)}>
              <UserPlus className="mr-2 h-4 w-4" />
              新增使用者
            </Button>
          </div>
        </div>

        {/* Statistics */}
        {userStats && (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-4">
            <Card>
              <CardHeader className="pb-3">
                <CardTitle className="text-sm flex items-center gap-2">
                  <Users className="h-4 w-4" />
                  總使用者數
                </CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-2xl font-bold">{userStats.total_users}</p>
                <p className="text-sm text-gray-500">
                  系統管理員: {userStats.admin_users}
                </p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="pb-3">
                <CardTitle className="text-sm flex items-center gap-2">
                  <UserCheck className="h-4 w-4 text-green-600" />
                  啟用使用者
                </CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-2xl font-bold">{userStats.active_users}</p>
                <p className="text-sm text-gray-500">
                  {((userStats.active_users / userStats.total_users) * 100).toFixed(1)}% 啟用率
                </p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="pb-3">
                <CardTitle className="text-sm flex items-center gap-2">
                  <Activity className="h-4 w-4 text-blue-600" />
                  線上使用者
                </CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-2xl font-bold">{onlineUsers?.length || 0}</p>
                <p className="text-sm text-gray-500">當前線上</p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="pb-3">
                <CardTitle className="text-sm flex items-center gap-2">
                  <Calendar className="h-4 w-4" />
                  本月新增
                </CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-2xl font-bold">{userStats.new_users_this_month}</p>
                <p className="text-sm text-gray-500">新使用者</p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="pb-3">
                <CardTitle className="text-sm flex items-center gap-2">
                  <Clock className="h-4 w-4" />
                  平均使用時長
                </CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-2xl font-bold">{userStats.avg_session_duration}</p>
                <p className="text-sm text-gray-500">分鐘/天</p>
              </CardContent>
            </Card>
          </div>
        )}

        <Tabs value={activeTab} onValueChange={setActiveTab}>
          <TabsList>
            <TabsTrigger value="list">使用者列表</TabsTrigger>
            <TabsTrigger value="roles">角色管理</TabsTrigger>
            <TabsTrigger value="permissions">權限設定</TabsTrigger>
            <TabsTrigger value="sessions">線上工作階段</TabsTrigger>
            <TabsTrigger value="logs">操作記錄</TabsTrigger>
          </TabsList>

          <TabsContent value="list" className="space-y-4">
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
                        placeholder="姓名、帳號或信箱"
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                        className="pl-10"
                        onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
                      />
                    </div>
                  </div>
                  <div className="space-y-2">
                    <label className="text-sm font-medium">角色</label>
                    <Select value={roleFilter} onValueChange={setRoleFilter}>
                      <SelectTrigger>
                        <SelectValue placeholder="選擇角色" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="">全部</SelectItem>
                        {rolesData?.data.map((role: any) => (
                          <SelectItem key={role.id} value={role.id}>
                            {role.name}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>
                  <div className="space-y-2">
                    <label className="text-sm font-medium">狀態</label>
                    <Select value={statusFilter} onValueChange={setStatusFilter}>
                      <SelectTrigger>
                        <SelectValue placeholder="選擇狀態" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="">全部</SelectItem>
                        <SelectItem value="active">啟用</SelectItem>
                        <SelectItem value="inactive">停用</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                  <div className="space-y-2">
                    <label className="text-sm font-medium">部門</label>
                    <Input
                      placeholder="輸入部門名稱"
                      value={departmentFilter}
                      onChange={(e) => setDepartmentFilter(e.target.value)}
                      onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
                    />
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
                </div>
              </CardContent>
            </Card>

            {/* Users Table */}
            <Card>
              <CardContent className="p-0">
                {isLoading ? (
                  <div className="text-center py-8">
                    <RefreshCw className="h-8 w-8 animate-spin mx-auto mb-4 text-gray-400" />
                    <p className="text-gray-500">載入中...</p>
                  </div>
                ) : usersData?.data.length === 0 ? (
                  <div className="text-center py-8">
                    <Users className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                    <p className="text-gray-500">暫無使用者資料</p>
                  </div>
                ) : (
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead>使用者</TableHead>
                        <TableHead>帳號/信箱</TableHead>
                        <TableHead>角色</TableHead>
                        <TableHead>部門/職位</TableHead>
                        <TableHead>狀態</TableHead>
                        <TableHead>最後登入</TableHead>
                        <TableHead>操作</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {usersData?.data.map((user: any) => (
                        <TableRow key={user.id}>
                          <TableCell>
                            <div className="flex items-center gap-3">
                              <Avatar>
                                <AvatarImage src={user.avatar} alt={user.full_name} />
                                <AvatarFallback>
                                  {user.full_name.substring(0, 2)}
                                </AvatarFallback>
                              </Avatar>
                              <div>
                                <p className="font-medium">{user.full_name}</p>
                                {user.full_name_en && (
                                  <p className="text-sm text-gray-500">{user.full_name_en}</p>
                                )}
                                {getOnlineStatusBadge(user.id)}
                              </div>
                            </div>
                          </TableCell>
                          <TableCell>
                            <div>
                              <p>{user.username}</p>
                              <p className="text-sm text-gray-500">{user.email}</p>
                            </div>
                          </TableCell>
                          <TableCell>
                            <div className="flex flex-wrap gap-1">
                              {user.roles?.map((role: any) => (
                                <Badge key={role.id} variant="outline">
                                  <Shield className="h-3 w-3 mr-1" />
                                  {role.name}
                                </Badge>
                              ))}
                            </div>
                          </TableCell>
                          <TableCell>
                            <div>
                              <p>{user.department || '-'}</p>
                              <p className="text-sm text-gray-500">{user.position || '-'}</p>
                            </div>
                          </TableCell>
                          <TableCell>{getStatusBadge(user.is_active)}</TableCell>
                          <TableCell>
                            {user.last_login ? (
                              <div>
                                <p>{format(new Date(user.last_login), 'yyyy/MM/dd', { locale: zhTW })}</p>
                                <p className="text-sm text-gray-500">
                                  {format(new Date(user.last_login), 'HH:mm', { locale: zhTW })}
                                </p>
                              </div>
                            ) : (
                              <span className="text-gray-400">尚未登入</span>
                            )}
                          </TableCell>
                          <TableCell>
                            <div className="flex items-center gap-2">
                              <Button
                                variant="ghost"
                                size="sm"
                                onClick={() => handleEdit(user)}
                              >
                                <Edit className="h-4 w-4" />
                              </Button>
                              <Button
                                variant="ghost"
                                size="sm"
                                onClick={() => {
                                  setSelectedUser(user)
                                  setIsPasswordDialogOpen(true)
                                }}
                              >
                                <Key className="h-4 w-4" />
                              </Button>
                              <Button
                                variant="ghost"
                                size="sm"
                                onClick={() => handleToggleUserStatus(user)}
                              >
                                {user.is_active ? (
                                  <Lock className="h-4 w-4" />
                                ) : (
                                  <Unlock className="h-4 w-4" />
                                )}
                              </Button>
                              <Button
                                variant="ghost"
                                size="sm"
                                onClick={() => handleDelete(user)}
                              >
                                <Trash2 className="h-4 w-4" />
                              </Button>
                            </div>
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                )}

                {/* Pagination */}
                {usersData && usersData.total > pageSize && (
                  <div className="flex items-center justify-between p-4 border-t">
                    <p className="text-sm text-gray-500">
                      顯示 {(page - 1) * pageSize + 1} 到 {Math.min(page * pageSize, usersData.total)} 項，
                      共 {usersData.total} 項
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
                        disabled={page * pageSize >= usersData.total}
                      >
                        下一頁
                      </Button>
                    </div>
                  </div>
                )}
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="roles">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <Shield className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">角色管理功能開發中</p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="permissions">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <Key className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">權限設定功能開發中</p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="sessions">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <Monitor className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">工作階段管理功能開發中</p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="logs">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <History className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">操作記錄功能開發中</p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>

        {/* Create User Dialog */}
        <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
          <DialogContent className="sm:max-w-[600px]">
            <DialogHeader>
              <DialogTitle>新增使用者</DialogTitle>
              <DialogDescription>
                建立新的系統使用者帳號
              </DialogDescription>
            </DialogHeader>
            <div className="grid gap-4 py-4">
              <div className="grid grid-cols-2 gap-4">
                <div className="grid gap-2">
                  <Label htmlFor="username">帳號 *</Label>
                  <Input
                    id="username"
                    value={formData.username}
                    onChange={(e) => setFormData({ ...formData, username: e.target.value })}
                    placeholder="輸入使用者帳號"
                  />
                </div>
                <div className="grid gap-2">
                  <Label htmlFor="email">電子郵件 *</Label>
                  <Input
                    id="email"
                    type="email"
                    value={formData.email}
                    onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                    placeholder="user@example.com"
                  />
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div className="grid gap-2">
                  <Label htmlFor="full_name">中文姓名 *</Label>
                  <Input
                    id="full_name"
                    value={formData.full_name}
                    onChange={(e) => setFormData({ ...formData, full_name: e.target.value })}
                    placeholder="輸入中文姓名"
                  />
                </div>
                <div className="grid gap-2">
                  <Label htmlFor="full_name_en">英文姓名</Label>
                  <Input
                    id="full_name_en"
                    value={formData.full_name_en}
                    onChange={(e) => setFormData({ ...formData, full_name_en: e.target.value })}
                    placeholder="輸入英文姓名"
                  />
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div className="grid gap-2">
                  <Label htmlFor="employee_no">員工編號</Label>
                  <Input
                    id="employee_no"
                    value={formData.employee_no}
                    onChange={(e) => setFormData({ ...formData, employee_no: e.target.value })}
                    placeholder="輸入員工編號"
                  />
                </div>
                <div className="grid gap-2">
                  <Label htmlFor="department">部門</Label>
                  <Input
                    id="department"
                    value={formData.department}
                    onChange={(e) => setFormData({ ...formData, department: e.target.value })}
                    placeholder="輸入部門名稱"
                  />
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div className="grid gap-2">
                  <Label htmlFor="position">職位</Label>
                  <Input
                    id="position"
                    value={formData.position}
                    onChange={(e) => setFormData({ ...formData, position: e.target.value })}
                    placeholder="輸入職位名稱"
                  />
                </div>
                <div className="grid gap-2">
                  <Label htmlFor="mobile">手機</Label>
                  <Input
                    id="mobile"
                    value={formData.mobile}
                    onChange={(e) => setFormData({ ...formData, mobile: e.target.value })}
                    placeholder="輸入手機號碼"
                  />
                </div>
              </div>

              <div className="grid gap-2">
                <Label>角色</Label>
                <div className="flex flex-wrap gap-2">
                  {rolesData?.data.map((role: any) => (
                    <label key={role.id} className="flex items-center gap-2">
                      <input
                        type="checkbox"
                        checked={formData.role_ids.includes(role.id)}
                        onChange={(e) => {
                          if (e.target.checked) {
                            setFormData({
                              ...formData,
                              role_ids: [...formData.role_ids, role.id],
                            })
                          } else {
                            setFormData({
                              ...formData,
                              role_ids: formData.role_ids.filter(id => id !== role.id),
                            })
                          }
                        }}
                        className="rounded border-gray-300"
                      />
                      <span className="text-sm">{role.name}</span>
                    </label>
                  ))}
                </div>
              </div>

              <div className="flex items-center justify-between">
                <Label htmlFor="is_active">帳號狀態</Label>
                <Switch
                  id="is_active"
                  checked={formData.is_active}
                  onCheckedChange={(checked) => setFormData({ ...formData, is_active: checked })}
                />
              </div>
            </div>
            <DialogFooter>
              <Button variant="outline" onClick={() => setIsCreateDialogOpen(false)}>
                取消
              </Button>
              <Button onClick={handleCreate}>
                建立使用者
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>

        {/* Edit User Dialog */}
        <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
          <DialogContent className="sm:max-w-[600px]">
            <DialogHeader>
              <DialogTitle>編輯使用者</DialogTitle>
              <DialogDescription>
                更新使用者資料
              </DialogDescription>
            </DialogHeader>
            <div className="grid gap-4 py-4">
              <div className="grid grid-cols-2 gap-4">
                <div className="grid gap-2">
                  <Label htmlFor="edit-username">帳號</Label>
                  <Input
                    id="edit-username"
                    value={formData.username}
                    disabled
                    className="bg-gray-50"
                  />
                </div>
                <div className="grid gap-2">
                  <Label htmlFor="edit-email">電子郵件 *</Label>
                  <Input
                    id="edit-email"
                    type="email"
                    value={formData.email}
                    onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                  />
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div className="grid gap-2">
                  <Label htmlFor="edit-full_name">中文姓名 *</Label>
                  <Input
                    id="edit-full_name"
                    value={formData.full_name}
                    onChange={(e) => setFormData({ ...formData, full_name: e.target.value })}
                  />
                </div>
                <div className="grid gap-2">
                  <Label htmlFor="edit-full_name_en">英文姓名</Label>
                  <Input
                    id="edit-full_name_en"
                    value={formData.full_name_en}
                    onChange={(e) => setFormData({ ...formData, full_name_en: e.target.value })}
                  />
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div className="grid gap-2">
                  <Label htmlFor="edit-employee_no">員工編號</Label>
                  <Input
                    id="edit-employee_no"
                    value={formData.employee_no}
                    onChange={(e) => setFormData({ ...formData, employee_no: e.target.value })}
                  />
                </div>
                <div className="grid gap-2">
                  <Label htmlFor="edit-department">部門</Label>
                  <Input
                    id="edit-department"
                    value={formData.department}
                    onChange={(e) => setFormData({ ...formData, department: e.target.value })}
                  />
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div className="grid gap-2">
                  <Label htmlFor="edit-position">職位</Label>
                  <Input
                    id="edit-position"
                    value={formData.position}
                    onChange={(e) => setFormData({ ...formData, position: e.target.value })}
                  />
                </div>
                <div className="grid gap-2">
                  <Label htmlFor="edit-mobile">手機</Label>
                  <Input
                    id="edit-mobile"
                    value={formData.mobile}
                    onChange={(e) => setFormData({ ...formData, mobile: e.target.value })}
                  />
                </div>
              </div>

              <div className="grid gap-2">
                <Label>角色</Label>
                <div className="flex flex-wrap gap-2">
                  {rolesData?.data.map((role: any) => (
                    <label key={role.id} className="flex items-center gap-2">
                      <input
                        type="checkbox"
                        checked={formData.role_ids.includes(role.id)}
                        onChange={(e) => {
                          if (e.target.checked) {
                            setFormData({
                              ...formData,
                              role_ids: [...formData.role_ids, role.id],
                            })
                          } else {
                            setFormData({
                              ...formData,
                              role_ids: formData.role_ids.filter(id => id !== role.id),
                            })
                          }
                        }}
                        className="rounded border-gray-300"
                      />
                      <span className="text-sm">{role.name}</span>
                    </label>
                  ))}
                </div>
              </div>

              <div className="flex items-center justify-between">
                <Label htmlFor="edit-is_active">帳號狀態</Label>
                <Switch
                  id="edit-is_active"
                  checked={formData.is_active}
                  onCheckedChange={(checked) => setFormData({ ...formData, is_active: checked })}
                />
              </div>
            </div>
            <DialogFooter>
              <Button variant="outline" onClick={() => setIsEditDialogOpen(false)}>
                取消
              </Button>
              <Button onClick={handleUpdate}>
                更新使用者
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>

        {/* Reset Password Dialog */}
        <Dialog open={isPasswordDialogOpen} onOpenChange={setIsPasswordDialogOpen}>
          <DialogContent className="sm:max-w-[400px]">
            <DialogHeader>
              <DialogTitle>重設密碼</DialogTitle>
              <DialogDescription>
                為使用者「{selectedUser?.full_name}」重設密碼
              </DialogDescription>
            </DialogHeader>
            <div className="grid gap-4 py-4">
              <div className="grid gap-2">
                <Label htmlFor="new_password">新密碼</Label>
                <Input
                  id="new_password"
                  type="password"
                  value={passwordData.new_password}
                  onChange={(e) => setPasswordData({ ...passwordData, new_password: e.target.value })}
                  placeholder="輸入新密碼"
                />
              </div>
              <div className="grid gap-2">
                <Label htmlFor="confirm_password">確認密碼</Label>
                <Input
                  id="confirm_password"
                  type="password"
                  value={passwordData.confirm_password}
                  onChange={(e) => setPasswordData({ ...passwordData, confirm_password: e.target.value })}
                  placeholder="再次輸入新密碼"
                />
              </div>
            </div>
            <DialogFooter>
              <Button variant="outline" onClick={() => setIsPasswordDialogOpen(false)}>
                取消
              </Button>
              <Button onClick={handleResetPassword}>
                重設密碼
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>
    </DashboardLayout>
  )
}