'use client'

import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useRouter } from 'next/navigation'
import DashboardLayout from '@/components/layout/DashboardLayout'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
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
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { useToast } from '@/components/ui/use-toast'
import { 
  FileText,
  Plus,
  Search,
  Filter,
  Download,
  Upload,
  Eye,
  Edit,
  Copy,
  Trash2,
  Layout,
  Database,
  LineChart,
  PieChart,
  BarChart3,
  TableIcon,
  Calendar,
  Users,
  Package,
  DollarSign,
  TrendingUp,
  AlertCircle,
  CheckCircle,
  Star,
  StarOff,
  Lock,
  Unlock
} from 'lucide-react'
import reportService from '@/services/report.service'
import { format } from 'date-fns'
import { zhTW } from 'date-fns/locale'

export default function ReportTemplatesPage() {
  const router = useRouter()
  const queryClient = useQueryClient()
  const { toast } = useToast()
  const [searchQuery, setSearchQuery] = useState('')
  const [categoryFilter, setCategoryFilter] = useState('')
  const [typeFilter, setTypeFilter] = useState('')
  const [isPublicFilter, setIsPublicFilter] = useState<string>('')
  const [page, setPage] = useState(1)
  const pageSize = 20

  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false)
  const [selectedTemplate, setSelectedTemplate] = useState<any>(null)
  const [formData, setFormData] = useState({
    name: '',
    name_en: '',
    description: '',
    category: 'sales',
    type: 'summary',
    is_public: false,
    tags: '',
  })

  // Fetch templates
  const { data: templatesData, isLoading, refetch } = useQuery({
    queryKey: ['report-templates', page, searchQuery, categoryFilter, typeFilter, isPublicFilter],
    queryFn: () => reportService.listReportTemplates({
      page,
      page_size: pageSize,
      search: searchQuery || undefined,
      category: categoryFilter || undefined,
      type: typeFilter || undefined,
      is_public: isPublicFilter === '' ? undefined : isPublicFilter === 'true',
    }),
  })

  // Create template mutation
  const createMutation = useMutation({
    mutationFn: (data: any) => reportService.createReportTemplate(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['report-templates'] })
      toast({ title: '範本建立成功' })
      setIsCreateDialogOpen(false)
      resetForm()
    },
    onError: (error: any) => {
      toast({
        title: '建立失敗',
        description: error.response?.data?.message || '建立範本時發生錯誤',
        variant: 'destructive',
      })
    },
  })

  // Delete template mutation
  const deleteMutation = useMutation({
    mutationFn: (id: string) => reportService.deleteReportTemplate(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['report-templates'] })
      toast({ title: '範本刪除成功' })
    },
    onError: (error: any) => {
      toast({
        title: '刪除失敗',
        description: error.response?.data?.message || '刪除範本時發生錯誤',
        variant: 'destructive',
      })
    },
  })

  const resetForm = () => {
    setFormData({
      name: '',
      name_en: '',
      description: '',
      category: 'sales',
      type: 'summary',
      is_public: false,
      tags: '',
    })
  }

  const handleSearch = () => {
    setPage(1)
    refetch()
  }

  const handleClearFilters = () => {
    setSearchQuery('')
    setCategoryFilter('')
    setTypeFilter('')
    setIsPublicFilter('')
    setPage(1)
    refetch()
  }

  const handleCreateFromTemplate = (templateId: string) => {
    router.push(`/reports/new?template=${templateId}`)
  }

  const handlePreviewTemplate = (templateId: string) => {
    router.push(`/reports/templates/${templateId}/preview`)
  }

  const handleEditTemplate = (templateId: string) => {
    router.push(`/reports/templates/${templateId}/edit`)
  }

  const handleDuplicateTemplate = async (templateId: string) => {
    try {
      await reportService.duplicateReportTemplate(templateId)
      refetch()
      toast({ title: '範本複製成功' })
    } catch (error) {
      toast({
        title: '複製失敗',
        variant: 'destructive',
      })
    }
  }

  const handleDeleteTemplate = (template: any) => {
    if (confirm(`確定要刪除範本「${template.name}」嗎？`)) {
      deleteMutation.mutate(template.id)
    }
  }

  const handleSubmit = () => {
    const tagsArray = formData.tags.split(',').map(tag => tag.trim()).filter(tag => tag)
    createMutation.mutate({
      ...formData,
      tags: tagsArray,
    })
  }

  const getCategoryIcon = (category: string) => {
    const icons: Record<string, any> = {
      sales: DollarSign,
      finance: TrendingUp,
      production: Package,
      inventory: Package,
      supplier: Users,
      customer: Users,
      system: Database,
    }
    return icons[category] || FileText
  }

  const getTypeIcon = (type: string) => {
    const icons: Record<string, any> = {
      summary: TableIcon,
      detail: FileText,
      trend: LineChart,
      comparison: BarChart3,
      dashboard: Layout,
    }
    return icons[type] || FileText
  }

  const getCategoryBadge = (category: string) => {
    const categoryConfig: Record<string, { label: string; variant: any }> = {
      sales: { label: '銷售', variant: 'info' },
      finance: { label: '財務', variant: 'success' },
      production: { label: '生產', variant: 'warning' },
      inventory: { label: '庫存', variant: 'secondary' },
      supplier: { label: '供應商', variant: 'info' },
      customer: { label: '客戶', variant: 'success' },
      system: { label: '系統', variant: 'secondary' },
    }

    const config = categoryConfig[category] || { label: category, variant: 'default' }
    return <Badge variant={config.variant as any}>{config.label}</Badge>
  }

  const getTypeBadge = (type: string) => {
    const typeConfig: Record<string, { label: string; variant: any }> = {
      summary: { label: '摘要', variant: 'info' },
      detail: { label: '詳細', variant: 'secondary' },
      trend: { label: '趨勢', variant: 'warning' },
      comparison: { label: '比較', variant: 'success' },
      dashboard: { label: '儀表板', variant: 'info' },
    }

    const config = typeConfig[type] || { label: type, variant: 'default' }
    return <Badge variant={config.variant as any}>{config.label}</Badge>
  }

  return (
    <DashboardLayout>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex justify-between items-center">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">報表範本庫</h1>
            <p className="mt-2 text-gray-600">預先設計的報表範本，快速建立專業報表</p>
          </div>
          <div className="flex gap-2">
            <Button variant="outline" onClick={() => router.push('/reports')}>
              返回報表中心
            </Button>
            <Button onClick={() => setIsCreateDialogOpen(true)}>
              <Plus className="mr-2 h-4 w-4" />
              新增範本
            </Button>
          </div>
        </div>

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
                    placeholder="範本名稱或描述"
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    className="pl-10"
                    onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
                  />
                </div>
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">分類</label>
                <Select value={categoryFilter} onValueChange={setCategoryFilter}>
                  <SelectTrigger>
                    <SelectValue placeholder="選擇分類" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="">全部</SelectItem>
                    <SelectItem value="sales">銷售</SelectItem>
                    <SelectItem value="finance">財務</SelectItem>
                    <SelectItem value="production">生產</SelectItem>
                    <SelectItem value="inventory">庫存</SelectItem>
                    <SelectItem value="supplier">供應商</SelectItem>
                    <SelectItem value="customer">客戶</SelectItem>
                    <SelectItem value="system">系統</SelectItem>
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
                    <SelectItem value="summary">摘要</SelectItem>
                    <SelectItem value="detail">詳細</SelectItem>
                    <SelectItem value="trend">趨勢</SelectItem>
                    <SelectItem value="comparison">比較</SelectItem>
                    <SelectItem value="dashboard">儀表板</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">範圍</label>
                <Select value={isPublicFilter} onValueChange={setIsPublicFilter}>
                  <SelectTrigger>
                    <SelectValue placeholder="選擇範圍" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="">全部</SelectItem>
                    <SelectItem value="true">公開範本</SelectItem>
                    <SelectItem value="false">私有範本</SelectItem>
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
            </div>
          </CardContent>
        </Card>

        {/* Templates Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
          {isLoading ? (
            <div className="col-span-full text-center py-12">
              <div className="text-gray-500">載入中...</div>
            </div>
          ) : templatesData?.data.length === 0 ? (
            <div className="col-span-full text-center py-12">
              <FileText className="h-12 w-12 mx-auto mb-4 text-gray-300" />
              <p className="text-gray-500">暫無範本資料</p>
            </div>
          ) : (
            templatesData?.data.map((template) => {
              const CategoryIcon = getCategoryIcon(template.category)
              const TypeIcon = getTypeIcon(template.type)
              
              return (
                <Card key={template.id} className="hover:shadow-lg transition-shadow cursor-pointer">
                  <CardHeader className="pb-4">
                    <div className="flex justify-between items-start">
                      <div className="flex items-center gap-2">
                        <CategoryIcon className="h-5 w-5 text-gray-500" />
                        <TypeIcon className="h-5 w-5 text-gray-500" />
                      </div>
                      <div className="flex items-center gap-1">
                        {(template as any).is_favorite ? (
                          <Star className="h-4 w-4 text-yellow-500 fill-yellow-500" />
                        ) : (
                          <StarOff className="h-4 w-4 text-gray-400" />
                        )}
                        {(template as any).is_public ? (
                          <Unlock className="h-4 w-4 text-green-500" />
                        ) : (
                          <Lock className="h-4 w-4 text-gray-400" />
                        )}
                      </div>
                    </div>
                    <CardTitle className="text-lg mt-2">{template.name}</CardTitle>
                    {template.name_en && (
                      <p className="text-sm text-gray-500">{template.name_en}</p>
                    )}
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-3">
                      {template.description && (
                        <p className="text-sm text-gray-600 line-clamp-2">
                          {template.description}
                        </p>
                      )}
                      
                      <div className="flex flex-wrap gap-2">
                        {getCategoryBadge(template.category)}
                        {getTypeBadge(template.type)}
                      </div>

                      {template.tags && template.tags.length > 0 && (
                        <div className="flex flex-wrap gap-1">
                          {(Array.isArray(template.tags) ? template.tags : []).map((tag: string, index: number) => (
                            <Badge key={index} variant="outline" className="text-xs">
                              {tag}
                            </Badge>
                          ))}
                        </div>
                      )}

                      <div className="flex items-center justify-between text-xs text-gray-500 pt-2">
                        <span>使用 {template.usage_count || 0} 次</span>
                        <span>{format(new Date(template.created_at), 'yyyy/MM/dd')}</span>
                      </div>

                      <div className="flex gap-2 pt-2">
                        <Button
                          size="sm"
                          className="flex-1"
                          onClick={() => handleCreateFromTemplate(template.id)}
                        >
                          使用範本
                        </Button>
                        <Button
                          size="sm"
                          variant="outline"
                          onClick={() => handlePreviewTemplate(template.id)}
                        >
                          <Eye className="h-4 w-4" />
                        </Button>
                        {!(template as any).is_system && (
                          <>
                            <Button
                              size="sm"
                              variant="outline"
                              onClick={() => handleEditTemplate(template.id)}
                            >
                              <Edit className="h-4 w-4" />
                            </Button>
                            <Button
                              size="sm"
                              variant="outline"
                              onClick={() => handleDuplicateTemplate(template.id)}
                            >
                              <Copy className="h-4 w-4" />
                            </Button>
                            <Button
                              size="sm"
                              variant="outline"
                              onClick={() => handleDeleteTemplate(template)}
                            >
                              <Trash2 className="h-4 w-4" />
                            </Button>
                          </>
                        )}
                      </div>
                    </div>
                  </CardContent>
                </Card>
              )
            })
          )}
        </div>

        {/* Pagination */}
        {templatesData && templatesData.total > pageSize && (
          <Card>
            <CardContent className="py-4">
              <div className="flex items-center justify-between">
                <p className="text-sm text-gray-500">
                  顯示 {(page - 1) * pageSize + 1} 到 {Math.min(page * pageSize, templatesData.total)} 項，
                  共 {templatesData.total} 項
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
                    disabled={page * pageSize >= templatesData.total}
                  >
                    下一頁
                  </Button>
                </div>
              </div>
            </CardContent>
          </Card>
        )}

        {/* Create Template Dialog */}
        <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
          <DialogContent className="sm:max-w-[600px]">
            <DialogHeader>
              <DialogTitle>新增報表範本</DialogTitle>
              <DialogDescription>
                建立新的報表範本供其他使用者使用
              </DialogDescription>
            </DialogHeader>
            <div className="grid gap-4 py-4">
              <div className="grid grid-cols-2 gap-4">
                <div className="grid gap-2">
                  <Label htmlFor="name">範本名稱 (中文) *</Label>
                  <Input
                    id="name"
                    value={formData.name}
                    onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                    placeholder="例如：月度銷售報表"
                  />
                </div>
                <div className="grid gap-2">
                  <Label htmlFor="name_en">範本名稱 (英文)</Label>
                  <Input
                    id="name_en"
                    value={formData.name_en}
                    onChange={(e) => setFormData({ ...formData, name_en: e.target.value })}
                    placeholder="例如：Monthly Sales Report"
                  />
                </div>
              </div>

              <div className="grid gap-2">
                <Label htmlFor="description">描述</Label>
                <Textarea
                  id="description"
                  value={formData.description}
                  onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                  placeholder="描述此範本的用途與內容"
                  rows={3}
                />
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div className="grid gap-2">
                  <Label htmlFor="category">分類 *</Label>
                  <Select
                    value={formData.category}
                    onValueChange={(value) => setFormData({ ...formData, category: value })}
                  >
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="sales">銷售</SelectItem>
                      <SelectItem value="finance">財務</SelectItem>
                      <SelectItem value="production">生產</SelectItem>
                      <SelectItem value="inventory">庫存</SelectItem>
                      <SelectItem value="supplier">供應商</SelectItem>
                      <SelectItem value="customer">客戶</SelectItem>
                      <SelectItem value="system">系統</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                <div className="grid gap-2">
                  <Label htmlFor="type">類型 *</Label>
                  <Select
                    value={formData.type}
                    onValueChange={(value) => setFormData({ ...formData, type: value })}
                  >
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="summary">摘要</SelectItem>
                      <SelectItem value="detail">詳細</SelectItem>
                      <SelectItem value="trend">趨勢</SelectItem>
                      <SelectItem value="comparison">比較</SelectItem>
                      <SelectItem value="dashboard">儀表板</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </div>

              <div className="grid gap-2">
                <Label htmlFor="tags">標籤</Label>
                <Input
                  id="tags"
                  value={formData.tags}
                  onChange={(e) => setFormData({ ...formData, tags: e.target.value })}
                  placeholder="用逗號分隔多個標籤，例如：月報,銷售,業績"
                />
              </div>

              <div className="flex items-center space-x-2">
                <input
                  type="checkbox"
                  id="is_public"
                  checked={formData.is_public}
                  onChange={(e) => setFormData({ ...formData, is_public: e.target.checked })}
                  className="rounded border-gray-300"
                />
                <Label htmlFor="is_public">設為公開範本</Label>
              </div>
            </div>
            <DialogFooter>
              <Button variant="outline" onClick={() => setIsCreateDialogOpen(false)}>
                取消
              </Button>
              <Button onClick={handleSubmit}>
                建立範本
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>
    </DashboardLayout>
  )
}