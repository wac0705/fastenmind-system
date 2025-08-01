'use client'

import { useState, useEffect } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useSearchParams } from 'next/navigation'
import DashboardLayout from '@/components/layout/DashboardLayout'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { Switch } from '@/components/ui/switch'
import { 
  Table, 
  TableBody, 
  TableCell, 
  TableHead, 
  TableHeader, 
  TableRow 
} from '@/components/ui/table'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Separator } from '@/components/ui/separator'
import { 
  Plus, 
  Edit, 
  Trash2, 
  Search,
  Building2,
  Mail,
  Phone,
  Globe,
  Download,
  Upload,
  DollarSign
} from 'lucide-react'
import { useToast } from '@/components/ui/use-toast'
import customerService, { Customer, CreateCustomerRequest } from '@/services/customer.service'
import Link from 'next/link'

export default function CustomersPage() {
  const queryClient = useQueryClient()
  const { toast } = useToast()
  const searchParams = useSearchParams()
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false)
  const [editingCustomer, setEditingCustomer] = useState<Customer | null>(null)
  const [searchTerm, setSearchTerm] = useState('')
  const [countryFilter, setCountryFilter] = useState('all')
  const [currentPage, setCurrentPage] = useState(1)
  
  const [formData, setFormData] = useState<CreateCustomerRequest>({
    customer_code: '',
    name: '',
    name_en: '',
    short_name: '',
    country: 'TW',
    tax_id: '',
    address: '',
    shipping_address: '',
    contact_person: '',
    contact_phone: '',
    contact_email: '',
    payment_terms: 'T/T 30 days',
    credit_limit: 0,
    currency: 'USD',
    is_active: true,
  })

  // Check for edit parameter in URL
  useEffect(() => {
    const editId = searchParams.get('edit')
    if (editId && data?.data) {
      const customerToEdit = data.data.find(c => c.id === editId)
      if (customerToEdit) {
        handleEdit(customerToEdit)
      }
    }
  }, [searchParams, data])

  // Fetch customers
  const { data, isLoading } = useQuery({
    queryKey: ['customers', currentPage, searchTerm, countryFilter],
    queryFn: () => customerService.list({
      page: currentPage,
      page_size: 10,
      search: searchTerm,
      country: countryFilter === 'all' ? undefined : countryFilter,
    }),
  })

  // Fetch country options
  const { data: countries = [] } = useQuery({
    queryKey: ['countries'],
    queryFn: () => customerService.getCountryOptions(),
  })

  // Fetch currency options
  const { data: currencies = [] } = useQuery({
    queryKey: ['currencies'],
    queryFn: () => customerService.getCurrencyOptions(),
  })

  // Create mutation
  const createMutation = useMutation({
    mutationFn: (data: CreateCustomerRequest) => customerService.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['customers'] })
      toast({ title: '客戶建立成功' })
      setIsCreateDialogOpen(false)
      resetForm()
    },
    onError: (error: any) => {
      toast({
        title: '建立失敗',
        description: error.response?.data?.message || '建立客戶時發生錯誤',
        variant: 'destructive',
      })
    },
  })

  // Update mutation
  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<CreateCustomerRequest> }) =>
      customerService.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['customers'] })
      toast({ title: '客戶更新成功' })
      setEditingCustomer(null)
      setIsCreateDialogOpen(false)
      resetForm()
    },
    onError: (error: any) => {
      toast({
        title: '更新失敗',
        description: error.response?.data?.message || '更新客戶時發生錯誤',
        variant: 'destructive',
      })
    },
  })

  // Delete mutation
  const deleteMutation = useMutation({
    mutationFn: (id: string) => customerService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['customers'] })
      toast({ title: '客戶刪除成功' })
    },
    onError: (error: any) => {
      toast({
        title: '刪除失敗',
        description: error.response?.data?.message || '刪除客戶時發生錯誤',
        variant: 'destructive',
      })
    },
  })

  const resetForm = () => {
    setFormData({
      customer_code: '',
      name: '',
      name_en: '',
      short_name: '',
      country: 'TW',
      tax_id: '',
      address: '',
      shipping_address: '',
      contact_person: '',
      contact_phone: '',
      contact_email: '',
      payment_terms: 'T/T 30 days',
      credit_limit: 0,
      currency: 'USD',
      is_active: true,
    })
  }

  const handleSubmit = () => {
    if (editingCustomer) {
      updateMutation.mutate({ id: editingCustomer.id, data: formData })
    } else {
      createMutation.mutate(formData)
    }
  }

  const handleEdit = (customer: Customer) => {
    setEditingCustomer(customer)
    setFormData({
      customer_code: customer.customer_code,
      name: customer.name,
      name_en: customer.name_en || '',
      short_name: customer.short_name || '',
      country: customer.country,
      tax_id: customer.tax_id || '',
      address: customer.address || '',
      shipping_address: customer.shipping_address || '',
      contact_person: customer.contact_person || '',
      contact_phone: customer.contact_phone || '',
      contact_email: customer.contact_email || '',
      payment_terms: customer.payment_terms || 'T/T 30 days',
      credit_limit: customer.credit_limit || 0,
      currency: customer.currency,
      is_active: customer.is_active,
    })
    setIsCreateDialogOpen(true)
  }

  const handleDelete = (id: string) => {
    if (confirm('確定要刪除這個客戶嗎？')) {
      deleteMutation.mutate(id)
    }
  }

  const handleExport = async () => {
    try {
      const blob = await customerService.exportCustomers({
        search: searchTerm,
        country: countryFilter === 'all' ? undefined : countryFilter,
      })
      const url = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `customers_${new Date().toISOString().split('T')[0].replace(/-/g, '')}.csv`
      a.click()
      window.URL.revokeObjectURL(url)
      toast({ title: '匯出成功' })
    } catch (error) {
      toast({
        title: '匯出失敗',
        variant: 'destructive',
      })
    }
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

  return (
    <DashboardLayout>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex justify-between items-center">
          <div>
            <h1 className="text-3xl font-bold text-gray-900 flex items-center gap-2">
              <Building2 className="h-8 w-8" />
              客戶管理
            </h1>
            <p className="mt-2 text-gray-600">管理客戶資料與信用額度</p>
          </div>
          <div className="flex gap-2">
            <Button variant="outline" onClick={handleExport}>
              <Download className="mr-2 h-4 w-4" />
              匯出
            </Button>
            <Button onClick={() => { resetForm(); setEditingCustomer(null); setIsCreateDialogOpen(true); }}>
              <Plus className="mr-2 h-4 w-4" />
              新增客戶
            </Button>
          </div>
        </div>

        {/* Search and Filters */}
        <Card>
          <CardContent className="p-6">
            <div className="flex gap-4">
              <div className="flex-1 relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-4 w-4" />
                <Input
                  placeholder="搜尋客戶名稱、代碼、聯絡人..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className="pl-10"
                />
              </div>
              <Select value={countryFilter} onValueChange={setCountryFilter}>
                <SelectTrigger className="w-[200px]">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">所有國家</SelectItem>
                  {countries.map((country) => (
                    <SelectItem key={country.code} value={country.code}>
                      {getCountryFlag(country.code)} {country.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </CardContent>
        </Card>

        {/* Customers Table */}
        <Card>
          <CardHeader>
            <CardTitle>客戶列表</CardTitle>
            <CardDescription>
              共 {data?.pagination.total || 0} 個客戶
            </CardDescription>
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <div className="text-center py-8">載入中...</div>
            ) : (
              <div className="overflow-x-auto">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>客戶代碼</TableHead>
                      <TableHead>客戶名稱</TableHead>
                      <TableHead>國家</TableHead>
                      <TableHead>聯絡人</TableHead>
                      <TableHead>信用額度</TableHead>
                      <TableHead>幣別</TableHead>
                      <TableHead>狀態</TableHead>
                      <TableHead className="text-right">操作</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {data?.data.map((customer) => (
                      <TableRow key={customer.id}>
                        <TableCell className="font-medium">
                          <Link href={`/customers/${customer.id}`} className="text-blue-600 hover:underline">
                            {customer.customer_code}
                          </Link>
                        </TableCell>
                        <TableCell>
                          <div>
                            <Link href={`/customers/${customer.id}`} className="font-medium hover:text-blue-600">
                              {customer.name}
                            </Link>
                            {customer.name_en && (
                              <p className="text-sm text-gray-500">{customer.name_en}</p>
                            )}
                          </div>
                        </TableCell>
                        <TableCell>
                          <span className="flex items-center gap-1">
                            {getCountryFlag(customer.country)}
                            {customer.country}
                          </span>
                        </TableCell>
                        <TableCell>
                          <div className="space-y-1">
                            {customer.contact_person && (
                              <p className="text-sm">{customer.contact_person}</p>
                            )}
                            {customer.contact_email && (
                              <p className="text-xs text-gray-500 flex items-center gap-1">
                                <Mail className="h-3 w-3" />
                                {customer.contact_email}
                              </p>
                            )}
                            {customer.contact_phone && (
                              <p className="text-xs text-gray-500 flex items-center gap-1">
                                <Phone className="h-3 w-3" />
                                {customer.contact_phone}
                              </p>
                            )}
                          </div>
                        </TableCell>
                        <TableCell>
                          {customer.credit_limit ? (
                            <span className="flex items-center gap-1">
                              <DollarSign className="h-3 w-3" />
                              {customer.credit_limit.toLocaleString()}
                            </span>
                          ) : (
                            '-'
                          )}
                        </TableCell>
                        <TableCell>{customer.currency}</TableCell>
                        <TableCell>
                          <Badge variant={customer.is_active ? 'success' : 'secondary'}>
                            {customer.is_active ? '啟用' : '停用'}
                          </Badge>
                        </TableCell>
                        <TableCell className="text-right">
                          <div className="flex justify-end gap-2">
                            <Button
                              variant="ghost"
                              size="icon"
                              onClick={() => handleEdit(customer)}
                            >
                              <Edit className="h-4 w-4" />
                            </Button>
                            <Button
                              variant="ghost"
                              size="icon"
                              onClick={() => handleDelete(customer.id)}
                            >
                              <Trash2 className="h-4 w-4" />
                            </Button>
                          </div>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Create/Edit Dialog */}
        <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
          <DialogContent className="sm:max-w-[700px] max-h-[90vh] overflow-y-auto">
            <DialogHeader>
              <DialogTitle>{editingCustomer ? '編輯客戶' : '新增客戶'}</DialogTitle>
              <DialogDescription>
                填寫客戶基本資料與聯絡資訊
              </DialogDescription>
            </DialogHeader>
            <div className="grid gap-4 py-4">
              <div className="grid grid-cols-2 gap-4">
                <div className="grid gap-2">
                  <Label htmlFor="customer_code">客戶代碼 *</Label>
                  <Input
                    id="customer_code"
                    value={formData.customer_code}
                    onChange={(e) => setFormData({ ...formData, customer_code: e.target.value })}
                    placeholder="例如: CUST-001"
                  />
                </div>
                <div className="grid gap-2">
                  <Label htmlFor="country">國家 *</Label>
                  <Select
                    value={formData.country}
                    onValueChange={(value) => setFormData({ ...formData, country: value })}
                  >
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      {countries.map((country) => (
                        <SelectItem key={country.code} value={country.code}>
                          {getCountryFlag(country.code)} {country.name}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div className="grid gap-2">
                  <Label htmlFor="name">客戶名稱 (中文) *</Label>
                  <Input
                    id="name"
                    value={formData.name}
                    onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                    placeholder="公司全名"
                  />
                </div>
                <div className="grid gap-2">
                  <Label htmlFor="name_en">客戶名稱 (英文)</Label>
                  <Input
                    id="name_en"
                    value={formData.name_en}
                    onChange={(e) => setFormData({ ...formData, name_en: e.target.value })}
                    placeholder="Company Name"
                  />
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div className="grid gap-2">
                  <Label htmlFor="short_name">簡稱</Label>
                  <Input
                    id="short_name"
                    value={formData.short_name}
                    onChange={(e) => setFormData({ ...formData, short_name: e.target.value })}
                  />
                </div>
                <div className="grid gap-2">
                  <Label htmlFor="tax_id">統一編號/稅號</Label>
                  <Input
                    id="tax_id"
                    value={formData.tax_id}
                    onChange={(e) => setFormData({ ...formData, tax_id: e.target.value })}
                  />
                </div>
              </div>

              <div className="grid gap-2">
                <Label htmlFor="address">公司地址</Label>
                <Input
                  id="address"
                  value={formData.address}
                  onChange={(e) => setFormData({ ...formData, address: e.target.value })}
                  placeholder="完整地址"
                />
              </div>

              <div className="grid gap-2">
                <Label htmlFor="shipping_address">送貨地址</Label>
                <Input
                  id="shipping_address"
                  value={formData.shipping_address}
                  onChange={(e) => setFormData({ ...formData, shipping_address: e.target.value })}
                  placeholder="如與公司地址不同"
                />
              </div>

              <Separator />

              <div className="grid grid-cols-3 gap-4">
                <div className="grid gap-2">
                  <Label htmlFor="contact_person">聯絡人</Label>
                  <Input
                    id="contact_person"
                    value={formData.contact_person}
                    onChange={(e) => setFormData({ ...formData, contact_person: e.target.value })}
                  />
                </div>
                <div className="grid gap-2">
                  <Label htmlFor="contact_phone">聯絡電話</Label>
                  <Input
                    id="contact_phone"
                    value={formData.contact_phone}
                    onChange={(e) => setFormData({ ...formData, contact_phone: e.target.value })}
                  />
                </div>
                <div className="grid gap-2">
                  <Label htmlFor="contact_email">聯絡信箱</Label>
                  <Input
                    id="contact_email"
                    type="email"
                    value={formData.contact_email}
                    onChange={(e) => setFormData({ ...formData, contact_email: e.target.value })}
                  />
                </div>
              </div>

              <Separator />

              <div className="grid grid-cols-3 gap-4">
                <div className="grid gap-2">
                  <Label htmlFor="payment_terms">付款條件</Label>
                  <Input
                    id="payment_terms"
                    value={formData.payment_terms}
                    onChange={(e) => setFormData({ ...formData, payment_terms: e.target.value })}
                    placeholder="例如: T/T 30 days"
                  />
                </div>
                <div className="grid gap-2">
                  <Label htmlFor="credit_limit">信用額度</Label>
                  <Input
                    id="credit_limit"
                    type="number"
                    value={formData.credit_limit}
                    onChange={(e) => setFormData({ ...formData, credit_limit: parseFloat(e.target.value) || 0 })}
                    min="0"
                    step="1000"
                  />
                </div>
                <div className="grid gap-2">
                  <Label htmlFor="currency">幣別</Label>
                  <Select
                    value={formData.currency}
                    onValueChange={(value) => setFormData({ ...formData, currency: value })}
                  >
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      {currencies.map((currency) => (
                        <SelectItem key={currency} value={currency}>
                          {currency}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
              </div>

              <div className="flex items-center space-x-2">
                <Switch
                  id="is_active"
                  checked={formData.is_active}
                  onCheckedChange={(checked) => setFormData({ ...formData, is_active: checked })}
                />
                <Label htmlFor="is_active">啟用客戶</Label>
              </div>
            </div>
            <DialogFooter>
              <Button variant="outline" onClick={() => setIsCreateDialogOpen(false)}>
                取消
              </Button>
              <Button onClick={handleSubmit}>
                {editingCustomer ? '更新' : '建立'}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>
    </DashboardLayout>
  )
}