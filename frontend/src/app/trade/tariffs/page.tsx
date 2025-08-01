'use client'

import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useRouter } from 'next/navigation'
import DashboardLayout from '@/components/layout/DashboardLayout'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'
import { 
  Calculator,
  Globe,
  Plus,
  Search,
  Filter,
  Download,
  Upload,
  Edit,
  Trash2,
  Eye,
  FileText,
  DollarSign,
  Percent,
  AlertTriangle,
  CheckCircle,
  Info,
  RefreshCw,
  TrendingUp,
  TrendingDown,
  Package,
  Shield,
  Calendar,
  Flag
} from 'lucide-react'
import { toast } from '@/components/ui/use-toast'
import tradeService from '@/services/trade.service'
import { format } from 'date-fns'
import { zhTW } from 'date-fns/locale'

export default function TariffsPage() {
  const router = useRouter()
  const queryClient = useQueryClient()
  const [activeTab, setActiveTab] = useState('codes')
  const [searchQuery, setSearchQuery] = useState('')
  const [categoryFilter, setCategoryFilter] = useState('')
  const [isCreateCodeDialogOpen, setIsCreateCodeDialogOpen] = useState(false)
  const [isCreateRateDialogOpen, setIsCreateRateDialogOpen] = useState(false)
  const [isCalculatorDialogOpen, setIsCalculatorDialogOpen] = useState(false)
  const [selectedTariffCode, setSelectedTariffCode] = useState<any>(null)

  // Form state for creating tariff code
  const [codeFormData, setCodeFormData] = useState({
    hs_code: '',
    description: '',
    description_en: '',
    category: '',
    unit: 'kg',
    base_rate: 0,
    preferential_rate: 0,
    vat: 0,
    excise_tax: 0,
    required_certs: [] as string[],
  })

  // Form state for creating tariff rate
  const [rateFormData, setRateFormData] = useState({
    tariff_code_id: '',
    country_code: '',
    country_name: '',
    rate: 0,
    rate_type: 'ad_valorem',
    minimum_duty: 0,
    maximum_duty: 0,
    currency: 'USD',
    trade_type: 'import',
    agreement_type: 'mfn',
    valid_from: format(new Date(), 'yyyy-MM-dd'),
    valid_to: '',
  })

  // Calculator form state
  const [calculatorForm, setCalculatorForm] = useState({
    hs_code: '',
    country_code: '',
    trade_type: 'import',
    value: 0,
  })

  const [calculationResult, setCalculationResult] = useState<any>(null)

  // Fetch tariff codes
  const { data: tariffCodes, isLoading: isLoadingCodes, refetch: refetchCodes } = useQuery({
    queryKey: ['tariff-codes', searchQuery, categoryFilter],
    queryFn: () => tradeService.listTariffCodes({
      hs_code: searchQuery || undefined,
      category: categoryFilter || undefined,
    }),
  })

  // Fetch tariff rates
  const { data: tariffRates, isLoading: isLoadingRates, refetch: refetchRates } = useQuery({
    queryKey: ['tariff-rates'],
    queryFn: () => tradeService.listTariffRates(),
  })

  // Create tariff code mutation
  const createTariffCodeMutation = useMutation({
    mutationFn: (data: any) => tradeService.createTariffCode(data),
    onSuccess: () => {
      toast({
        title: '成功',
        description: '關稅代碼已成功建立',
      })
      queryClient.invalidateQueries({ queryKey: ['tariff-codes'] })
      setIsCreateCodeDialogOpen(false)
      resetCodeForm()
    },
    onError: (error: any) => {
      toast({
        title: '錯誤',
        description: error.response?.data?.message || '建立關稅代碼失敗',
        variant: 'destructive',
      })
    },
  })

  // Create tariff rate mutation
  const createTariffRateMutation = useMutation({
    mutationFn: (data: any) => tradeService.createTariffRate(data),
    onSuccess: () => {
      toast({
        title: '成功',
        description: '關稅稅率已成功建立',
      })
      queryClient.invalidateQueries({ queryKey: ['tariff-rates'] })
      queryClient.invalidateQueries({ queryKey: ['tariff-code-rates'] })
      setIsCreateRateDialogOpen(false)
      resetRateForm()
    },
    onError: (error: any) => {
      toast({
        title: '錯誤',
        description: error.response?.data?.message || '建立關稅稅率失敗',
        variant: 'destructive',
      })
    },
  })

  // Update tariff code mutation
  const updateTariffCodeMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: any }) => 
      tradeService.updateTariffCode(id, data),
    onSuccess: () => {
      toast({
        title: '成功',
        description: '關稅代碼已更新',
      })
      queryClient.invalidateQueries({ queryKey: ['tariff-codes'] })
    },
    onError: (error: any) => {
      toast({
        title: '錯誤',
        description: error.response?.data?.message || '更新關稅代碼失敗',
        variant: 'destructive',
      })
    },
  })

  // Delete tariff code mutation
  const deleteTariffCodeMutation = useMutation({
    mutationFn: (id: string) => tradeService.deleteTariffCode(id),
    onSuccess: () => {
      toast({
        title: '成功',
        description: '關稅代碼已刪除',
      })
      queryClient.invalidateQueries({ queryKey: ['tariff-codes'] })
    },
    onError: (error: any) => {
      toast({
        title: '錯誤',
        description: error.response?.data?.message || '刪除關稅代碼失敗',
        variant: 'destructive',
      })
    },
  })

  // Calculate tariff duty mutation
  const calculateDutyMutation = useMutation({
    mutationFn: (data: any) => tradeService.calculateTariffDuty(data),
    onSuccess: (result) => {
      setCalculationResult(result)
    },
    onError: (error: any) => {
      toast({
        title: '錯誤',
        description: error.response?.data?.message || '計算關稅失敗',
        variant: 'destructive',
      })
    },
  })

  const resetCodeForm = () => {
    setCodeFormData({
      hs_code: '',
      description: '',
      description_en: '',
      category: '',
      unit: 'kg',
      base_rate: 0,
      preferential_rate: 0,
      vat: 0,
      excise_tax: 0,
      required_certs: [],
    })
  }

  const resetRateForm = () => {
    setRateFormData({
      tariff_code_id: '',
      country_code: '',
      country_name: '',
      rate: 0,
      rate_type: 'ad_valorem',
      minimum_duty: 0,
      maximum_duty: 0,
      currency: 'USD',
      trade_type: 'import',
      agreement_type: 'mfn',
      valid_from: format(new Date(), 'yyyy-MM-dd'),
      valid_to: '',
    })
  }

  const handleCreateTariffCode = () => {
    if (!codeFormData.hs_code || !codeFormData.description) {
      toast({
        title: '錯誤',
        description: '請填寫必要欄位',
        variant: 'destructive',
      })
      return
    }

    createTariffCodeMutation.mutate(codeFormData)
  }

  const handleCreateTariffRate = () => {
    if (!rateFormData.tariff_code_id || !rateFormData.country_code || !rateFormData.country_name) {
      toast({
        title: '錯誤',
        description: '請填寫必要欄位',
        variant: 'destructive',
      })
      return
    }

    createTariffRateMutation.mutate({
      ...rateFormData,
      valid_from: new Date(rateFormData.valid_from).toISOString(),
      valid_to: rateFormData.valid_to ? new Date(rateFormData.valid_to).toISOString() : undefined,
    })
  }

  const handleCalculateDuty = () => {
    if (!calculatorForm.hs_code || !calculatorForm.country_code || calculatorForm.value <= 0) {
      toast({
        title: '錯誤',
        description: '請填寫所有必要欄位',
        variant: 'destructive',
      })
      return
    }

    calculateDutyMutation.mutate(calculatorForm)
  }

  const handleToggleActive = (id: string, isActive: boolean) => {
    updateTariffCodeMutation.mutate({
      id,
      data: { is_active: !isActive },
    })
  }

  const getRateTypeBadge = (rateType: string) => {
    const typeConfig: Record<string, { label: string; variant: any }> = {
      ad_valorem: { label: '從價稅', variant: 'default' },
      specific: { label: '從量稅', variant: 'secondary' },
      compound: { label: '複合稅', variant: 'outline' },
    }

    const config = typeConfig[rateType] || { label: rateType, variant: 'default' }
    
    return (
      <Badge variant={config.variant as any}>
        {config.label}
      </Badge>
    )
  }

  const getAgreementTypeBadge = (agreementType: string) => {
    const typeConfig: Record<string, { label: string; variant: any }> = {
      mfn: { label: '最惠國', variant: 'default' },
      fta: { label: '自由貿易協定', variant: 'success' },
      gsp: { label: '普惠制', variant: 'info' },
      cptpp: { label: 'CPTPP', variant: 'secondary' },
    }

    const config = typeConfig[agreementType] || { label: agreementType, variant: 'default' }
    
    return (
      <Badge variant={config.variant as any}>
        {config.label}
      </Badge>
    )
  }

  const filteredCodes = tariffCodes?.data?.filter((code) => {
    if (searchQuery) {
      const query = searchQuery.toLowerCase()
      return (
        code.hs_code.toLowerCase().includes(query) ||
        code.description.toLowerCase().includes(query) ||
        code.description_en?.toLowerCase().includes(query)
      )
    }
    return true
  }) || []

  return (
    <DashboardLayout>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex justify-between items-center">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">關稅管理</h1>
            <p className="mt-2 text-gray-600">管理 HS 代碼、設定關稅稅率、計算進出口關稅</p>
          </div>
          <div className="flex items-center gap-4">
            <Dialog open={isCalculatorDialogOpen} onOpenChange={setIsCalculatorDialogOpen}>
              <DialogTrigger asChild>
                <Button variant="outline">
                  <Calculator className="mr-2 h-4 w-4" />
                  關稅計算器
                </Button>
              </DialogTrigger>
              <DialogContent className="max-w-2xl">
                <DialogHeader>
                  <DialogTitle>關稅計算器</DialogTitle>
                  <DialogDescription>
                    輸入商品資訊以計算應繳關稅
                  </DialogDescription>
                </DialogHeader>
                <div className="grid gap-4 py-4">
                  <div className="grid grid-cols-2 gap-4">
                    <div className="space-y-2">
                      <Label htmlFor="calc-hs-code">HS 代碼*</Label>
                      <Input
                        id="calc-hs-code"
                        value={calculatorForm.hs_code}
                        onChange={(e) => setCalculatorForm({ ...calculatorForm, hs_code: e.target.value })}
                        placeholder="輸入 HS 代碼"
                      />
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="calc-country">國家代碼*</Label>
                      <Input
                        id="calc-country"
                        value={calculatorForm.country_code}
                        onChange={(e) => setCalculatorForm({ ...calculatorForm, country_code: e.target.value })}
                        placeholder="如: US, CN, JP"
                      />
                    </div>
                  </div>
                  <div className="grid grid-cols-2 gap-4">
                    <div className="space-y-2">
                      <Label htmlFor="calc-type">貿易類型*</Label>
                      <Select value={calculatorForm.trade_type} onValueChange={(value) => setCalculatorForm({ ...calculatorForm, trade_type: value })}>
                        <SelectTrigger>
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="import">進口</SelectItem>
                          <SelectItem value="export">出口</SelectItem>
                        </SelectContent>
                      </Select>
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="calc-value">貨物價值 (USD)*</Label>
                      <Input
                        id="calc-value"
                        type="number"
                        value={calculatorForm.value}
                        onChange={(e) => setCalculatorForm({ ...calculatorForm, value: parseFloat(e.target.value) || 0 })}
                        placeholder="0.00"
                      />
                    </div>
                  </div>

                  {calculationResult && (
                    <div className="mt-6 p-4 bg-gray-50 rounded-lg space-y-4">
                      <h4 className="font-medium flex items-center gap-2">
                        <Calculator className="h-4 w-4" />
                        計算結果
                      </h4>
                      <div className="grid grid-cols-2 gap-4">
                        <div>
                          <p className="text-sm text-gray-500">適用稅率</p>
                          <p className="font-medium">{calculationResult.applied_rate}%</p>
                        </div>
                        <div>
                          <p className="text-sm text-gray-500">關稅金額</p>
                          <p className="font-medium">{tradeService.formatCurrency(calculationResult.duty)}</p>
                        </div>
                        <div>
                          <p className="text-sm text-gray-500">增值稅 (VAT)</p>
                          <p className="font-medium">{tradeService.formatCurrency(calculationResult.vat)}</p>
                        </div>
                        <div>
                          <p className="text-sm text-gray-500">消費稅</p>
                          <p className="font-medium">{tradeService.formatCurrency(calculationResult.excise_tax)}</p>
                        </div>
                      </div>
                      <div className="border-t pt-4">
                        <div className="flex justify-between items-center">
                          <span className="font-medium">總計應繳稅額</span>
                          <span className="text-xl font-bold text-blue-600">
                            {tradeService.formatCurrency(calculationResult.total)}
                          </span>
                        </div>
                      </div>
                      {calculationResult.tariff_code && (
                        <div className="text-sm text-gray-600">
                          <p>HS 代碼: {calculationResult.tariff_code.hs_code}</p>
                          <p>商品描述: {calculationResult.tariff_code.description}</p>
                        </div>
                      )}
                    </div>
                  )}
                </div>
                <DialogFooter>
                  <Button variant="outline" onClick={() => {
                    setIsCalculatorDialogOpen(false)
                    setCalculationResult(null)
                  }}>
                    關閉
                  </Button>
                  <Button onClick={handleCalculateDuty} disabled={calculateDutyMutation.isPending}>
                    {calculateDutyMutation.isPending ? '計算中...' : '計算關稅'}
                  </Button>
                </DialogFooter>
              </DialogContent>
            </Dialog>
            <Button variant="outline">
              <Download className="mr-2 h-4 w-4" />
              匯出稅率表
            </Button>
            <Dialog open={isCreateCodeDialogOpen} onOpenChange={setIsCreateCodeDialogOpen}>
              <DialogTrigger asChild>
                <Button>
                  <Plus className="mr-2 h-4 w-4" />
                  新增 HS 代碼
                </Button>
              </DialogTrigger>
              <DialogContent className="max-w-2xl">
                <DialogHeader>
                  <DialogTitle>新增 HS 代碼</DialogTitle>
                  <DialogDescription>
                    新增關稅代碼及基本稅率設定
                  </DialogDescription>
                </DialogHeader>
                <div className="grid gap-4 py-4">
                  <div className="grid grid-cols-2 gap-4">
                    <div className="space-y-2">
                      <Label htmlFor="hs_code">HS 代碼*</Label>
                      <Input
                        id="hs_code"
                        value={codeFormData.hs_code}
                        onChange={(e) => setCodeFormData({ ...codeFormData, hs_code: e.target.value })}
                        placeholder="如: 7318.15.00"
                      />
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="category">分類</Label>
                      <Input
                        id="category"
                        value={codeFormData.category}
                        onChange={(e) => setCodeFormData({ ...codeFormData, category: e.target.value })}
                        placeholder="如: 螺絲螺栓"
                      />
                    </div>
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="description">中文描述*</Label>
                    <Input
                      id="description"
                      value={codeFormData.description}
                      onChange={(e) => setCodeFormData({ ...codeFormData, description: e.target.value })}
                      placeholder="輸入商品中文描述"
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="description_en">英文描述</Label>
                    <Input
                      id="description_en"
                      value={codeFormData.description_en}
                      onChange={(e) => setCodeFormData({ ...codeFormData, description_en: e.target.value })}
                      placeholder="輸入商品英文描述"
                    />
                  </div>
                  <div className="grid grid-cols-2 gap-4">
                    <div className="space-y-2">
                      <Label htmlFor="unit">計量單位</Label>
                      <Select value={codeFormData.unit} onValueChange={(value) => setCodeFormData({ ...codeFormData, unit: value })}>
                        <SelectTrigger>
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="kg">公斤 (kg)</SelectItem>
                          <SelectItem value="piece">件 (piece)</SelectItem>
                          <SelectItem value="m">公尺 (m)</SelectItem>
                          <SelectItem value="m2">平方公尺 (m²)</SelectItem>
                          <SelectItem value="m3">立方公尺 (m³)</SelectItem>
                          <SelectItem value="liter">公升 (L)</SelectItem>
                        </SelectContent>
                      </Select>
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="base_rate">基本稅率 (%)</Label>
                      <Input
                        id="base_rate"
                        type="number"
                        value={codeFormData.base_rate}
                        onChange={(e) => setCodeFormData({ ...codeFormData, base_rate: parseFloat(e.target.value) || 0 })}
                        placeholder="0.0"
                      />
                    </div>
                  </div>
                  <div className="grid grid-cols-3 gap-4">
                    <div className="space-y-2">
                      <Label htmlFor="preferential_rate">優惠稅率 (%)</Label>
                      <Input
                        id="preferential_rate"
                        type="number"
                        value={codeFormData.preferential_rate}
                        onChange={(e) => setCodeFormData({ ...codeFormData, preferential_rate: parseFloat(e.target.value) || 0 })}
                        placeholder="0.0"
                      />
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="vat">增值稅率 (%)</Label>
                      <Input
                        id="vat"
                        type="number"
                        value={codeFormData.vat}
                        onChange={(e) => setCodeFormData({ ...codeFormData, vat: parseFloat(e.target.value) || 0 })}
                        placeholder="0.0"
                      />
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="excise_tax">消費稅率 (%)</Label>
                      <Input
                        id="excise_tax"
                        type="number"
                        value={codeFormData.excise_tax}
                        onChange={(e) => setCodeFormData({ ...codeFormData, excise_tax: parseFloat(e.target.value) || 0 })}
                        placeholder="0.0"
                      />
                    </div>
                  </div>
                </div>
                <DialogFooter>
                  <Button variant="outline" onClick={() => setIsCreateCodeDialogOpen(false)}>
                    取消
                  </Button>
                  <Button onClick={handleCreateTariffCode} disabled={createTariffCodeMutation.isPending}>
                    {createTariffCodeMutation.isPending ? '建立中...' : '建立'}
                  </Button>
                </DialogFooter>
              </DialogContent>
            </Dialog>
          </div>
        </div>

        <Tabs value={activeTab} onValueChange={setActiveTab}>
          <TabsList>
            <TabsTrigger value="codes">HS 代碼管理</TabsTrigger>
            <TabsTrigger value="rates">國家稅率設定</TabsTrigger>
            <TabsTrigger value="agreements">貿易協定</TabsTrigger>
            <TabsTrigger value="statistics">統計分析</TabsTrigger>
          </TabsList>

          <TabsContent value="codes" className="space-y-4">
            {/* Search and Filter */}
            <Card>
              <CardContent className="p-6">
                <div className="flex items-center gap-4">
                  <div className="flex-1">
                    <div className="relative">
                      <Search className="absolute left-3 top-2.5 h-4 w-4 text-gray-400" />
                      <Input
                        placeholder="搜尋 HS 代碼或商品描述..."
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                        className="pl-10"
                      />
                    </div>
                  </div>
                  <Select value={categoryFilter} onValueChange={setCategoryFilter}>
                    <SelectTrigger className="w-48">
                      <SelectValue placeholder="選擇分類" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="">全部分類</SelectItem>
                      <SelectItem value="螺絲螺栓">螺絲螺栓</SelectItem>
                      <SelectItem value="螺帽">螺帽</SelectItem>
                      <SelectItem value="墊圈">墊圈</SelectItem>
                      <SelectItem value="鉚釘">鉚釘</SelectItem>
                      <SelectItem value="其他">其他</SelectItem>
                    </SelectContent>
                  </Select>
                  <Button variant="outline" size="icon" onClick={() => refetchCodes()}>
                    <RefreshCw className="h-4 w-4" />
                  </Button>
                </div>
              </CardContent>
            </Card>

            {/* HS Codes Table */}
            <Card>
              <CardContent className="p-0">
                {isLoadingCodes ? (
                  <div className="flex items-center justify-center h-64">
                    <div className="text-center">載入中...</div>
                  </div>
                ) : filteredCodes.length === 0 ? (
                  <div className="flex flex-col items-center justify-center h-64 text-gray-500">
                    <Package className="h-12 w-12 mb-4 text-gray-300" />
                    <p>沒有找到 HS 代碼</p>
                  </div>
                ) : (
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead>HS 代碼</TableHead>
                        <TableHead>商品描述</TableHead>
                        <TableHead>分類</TableHead>
                        <TableHead className="text-center">基本稅率</TableHead>
                        <TableHead className="text-center">優惠稅率</TableHead>
                        <TableHead className="text-center">增值稅</TableHead>
                        <TableHead className="text-center">狀態</TableHead>
                        <TableHead className="text-right">操作</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {filteredCodes.map((code) => (
                        <TableRow key={code.id}>
                          <TableCell className="font-medium">
                            <div className="flex items-center gap-2">
                              <Globe className="h-4 w-4 text-gray-400" />
                              {code.hs_code}
                            </div>
                          </TableCell>
                          <TableCell>
                            <div>
                              <p className="font-medium">{code.description}</p>
                              {code.description_en && (
                                <p className="text-sm text-gray-500">{code.description_en}</p>
                              )}
                            </div>
                          </TableCell>
                          <TableCell>
                            {code.category && (
                              <Badge variant="outline">{code.category}</Badge>
                            )}
                          </TableCell>
                          <TableCell className="text-center">
                            <div className="flex items-center justify-center gap-1">
                              <Percent className="h-3 w-3" />
                              {code.base_rate}%
                            </div>
                          </TableCell>
                          <TableCell className="text-center">
                            <div className="flex items-center justify-center gap-1">
                              <Percent className="h-3 w-3" />
                              {code.preferential_rate}%
                            </div>
                          </TableCell>
                          <TableCell className="text-center">
                            <div className="flex items-center justify-center gap-1">
                              <Percent className="h-3 w-3" />
                              {code.vat}%
                            </div>
                          </TableCell>
                          <TableCell className="text-center">
                            <Badge variant={code.is_active ? 'success' as any : 'secondary' as any}>
                              {code.is_active ? '啟用' : '停用'}
                            </Badge>
                          </TableCell>
                          <TableCell className="text-right">
                            <div className="flex items-center justify-end gap-2">
                              <Button
                                variant="ghost"
                                size="icon"
                                onClick={() => {
                                  setSelectedTariffCode(code)
                                  setRateFormData({ ...rateFormData, tariff_code_id: code.id })
                                  setIsCreateRateDialogOpen(true)
                                }}
                              >
                                <Plus className="h-4 w-4" />
                              </Button>
                              <Button
                                variant="ghost"
                                size="icon"
                                onClick={() => handleToggleActive(code.id, code.is_active)}
                              >
                                {code.is_active ? (
                                  <Shield className="h-4 w-4" />
                                ) : (
                                  <Shield className="h-4 w-4 text-gray-400" />
                                )}
                              </Button>
                              <Button
                                variant="ghost"
                                size="icon"
                                onClick={() => deleteTariffCodeMutation.mutate(code.id)}
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
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="rates" className="space-y-4">
            <div className="flex justify-end">
              <Dialog open={isCreateRateDialogOpen} onOpenChange={setIsCreateRateDialogOpen}>
                <DialogTrigger asChild>
                  <Button>
                    <Plus className="mr-2 h-4 w-4" />
                    新增稅率
                  </Button>
                </DialogTrigger>
                <DialogContent className="max-w-2xl">
                  <DialogHeader>
                    <DialogTitle>新增國家稅率</DialogTitle>
                    <DialogDescription>
                      為特定國家設定關稅稅率
                      {selectedTariffCode && (
                        <span className="block mt-2 font-medium">
                          HS 代碼: {selectedTariffCode.hs_code} - {selectedTariffCode.description}
                        </span>
                      )}
                    </DialogDescription>
                  </DialogHeader>
                  <div className="grid gap-4 py-4">
                    {!selectedTariffCode && (
                      <div className="space-y-2">
                        <Label htmlFor="rate-tariff-code">選擇 HS 代碼*</Label>
                        <Select value={rateFormData.tariff_code_id} onValueChange={(value) => setRateFormData({ ...rateFormData, tariff_code_id: value })}>
                          <SelectTrigger>
                            <SelectValue placeholder="選擇 HS 代碼" />
                          </SelectTrigger>
                          <SelectContent>
                            {tariffCodes?.data?.map((code) => (
                              <SelectItem key={code.id} value={code.id}>
                                {code.hs_code} - {code.description}
                              </SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                      </div>
                    )}
                    <div className="grid grid-cols-2 gap-4">
                      <div className="space-y-2">
                        <Label htmlFor="country_code">國家代碼*</Label>
                        <Input
                          id="country_code"
                          value={rateFormData.country_code}
                          onChange={(e) => setRateFormData({ ...rateFormData, country_code: e.target.value })}
                          placeholder="如: US, CN, JP"
                        />
                      </div>
                      <div className="space-y-2">
                        <Label htmlFor="country_name">國家名稱*</Label>
                        <Input
                          id="country_name"
                          value={rateFormData.country_name}
                          onChange={(e) => setRateFormData({ ...rateFormData, country_name: e.target.value })}
                          placeholder="如: 美國, 中國, 日本"
                        />
                      </div>
                    </div>
                    <div className="grid grid-cols-3 gap-4">
                      <div className="space-y-2">
                        <Label htmlFor="rate">稅率 (%)*</Label>
                        <Input
                          id="rate"
                          type="number"
                          value={rateFormData.rate}
                          onChange={(e) => setRateFormData({ ...rateFormData, rate: parseFloat(e.target.value) || 0 })}
                          placeholder="0.0"
                        />
                      </div>
                      <div className="space-y-2">
                        <Label htmlFor="rate_type">稅率類型*</Label>
                        <Select value={rateFormData.rate_type} onValueChange={(value) => setRateFormData({ ...rateFormData, rate_type: value })}>
                          <SelectTrigger>
                            <SelectValue />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectItem value="ad_valorem">從價稅</SelectItem>
                            <SelectItem value="specific">從量稅</SelectItem>
                            <SelectItem value="compound">複合稅</SelectItem>
                          </SelectContent>
                        </Select>
                      </div>
                      <div className="space-y-2">
                        <Label htmlFor="trade_type">貿易類型*</Label>
                        <Select value={rateFormData.trade_type} onValueChange={(value) => setRateFormData({ ...rateFormData, trade_type: value })}>
                          <SelectTrigger>
                            <SelectValue />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectItem value="import">進口</SelectItem>
                            <SelectItem value="export">出口</SelectItem>
                          </SelectContent>
                        </Select>
                      </div>
                    </div>
                    <div className="grid grid-cols-3 gap-4">
                      <div className="space-y-2">
                        <Label htmlFor="minimum_duty">最低稅額</Label>
                        <Input
                          id="minimum_duty"
                          type="number"
                          value={rateFormData.minimum_duty}
                          onChange={(e) => setRateFormData({ ...rateFormData, minimum_duty: parseFloat(e.target.value) || 0 })}
                          placeholder="0.00"
                        />
                      </div>
                      <div className="space-y-2">
                        <Label htmlFor="maximum_duty">最高稅額</Label>
                        <Input
                          id="maximum_duty"
                          type="number"
                          value={rateFormData.maximum_duty}
                          onChange={(e) => setRateFormData({ ...rateFormData, maximum_duty: parseFloat(e.target.value) || 0 })}
                          placeholder="0.00"
                        />
                      </div>
                      <div className="space-y-2">
                        <Label htmlFor="agreement_type">協定類型</Label>
                        <Select value={rateFormData.agreement_type} onValueChange={(value) => setRateFormData({ ...rateFormData, agreement_type: value })}>
                          <SelectTrigger>
                            <SelectValue />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectItem value="mfn">最惠國待遇</SelectItem>
                            <SelectItem value="fta">自由貿易協定</SelectItem>
                            <SelectItem value="gsp">普惠制</SelectItem>
                            <SelectItem value="cptpp">CPTPP</SelectItem>
                          </SelectContent>
                        </Select>
                      </div>
                    </div>
                    <div className="grid grid-cols-2 gap-4">
                      <div className="space-y-2">
                        <Label htmlFor="valid_from">生效日期*</Label>
                        <Input
                          id="valid_from"
                          type="date"
                          value={rateFormData.valid_from}
                          onChange={(e) => setRateFormData({ ...rateFormData, valid_from: e.target.value })}
                        />
                      </div>
                      <div className="space-y-2">
                        <Label htmlFor="valid_to">失效日期</Label>
                        <Input
                          id="valid_to"
                          type="date"
                          value={rateFormData.valid_to}
                          onChange={(e) => setRateFormData({ ...rateFormData, valid_to: e.target.value })}
                        />
                      </div>
                    </div>
                  </div>
                  <DialogFooter>
                    <Button variant="outline" onClick={() => {
                      setIsCreateRateDialogOpen(false)
                      setSelectedTariffCode(null)
                      resetRateForm()
                    }}>
                      取消
                    </Button>
                    <Button onClick={handleCreateTariffRate} disabled={createTariffRateMutation.isPending}>
                      {createTariffRateMutation.isPending ? '建立中...' : '建立'}
                    </Button>
                  </DialogFooter>
                </DialogContent>
              </Dialog>
            </div>

            <Card>
              <CardContent className="p-0">
                {isLoadingRates ? (
                  <div className="flex items-center justify-center h-64">
                    <div className="text-center">載入中...</div>
                  </div>
                ) : !tariffRates?.data || tariffRates.data.length === 0 ? (
                  <div className="flex flex-col items-center justify-center h-64 text-gray-500">
                    <Globe className="h-12 w-12 mb-4 text-gray-300" />
                    <p>尚未設定國家稅率</p>
                  </div>
                ) : (
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead>HS 代碼</TableHead>
                        <TableHead>國家</TableHead>
                        <TableHead>貿易類型</TableHead>
                        <TableHead className="text-center">稅率</TableHead>
                        <TableHead className="text-center">稅率類型</TableHead>
                        <TableHead>協定類型</TableHead>
                        <TableHead>有效期間</TableHead>
                        <TableHead className="text-center">狀態</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {tariffRates.data.map((rate) => (
                        <TableRow key={rate.id}>
                          <TableCell>
                            <div>
                              <p className="font-medium">{rate.tariff_code?.hs_code}</p>
                              <p className="text-sm text-gray-500">{rate.tariff_code?.description}</p>
                            </div>
                          </TableCell>
                          <TableCell>
                            <div className="flex items-center gap-2">
                              <Flag className="h-4 w-4" />
                              <span>{rate.country_name} ({rate.country_code})</span>
                            </div>
                          </TableCell>
                          <TableCell>
                            <Badge variant={rate.trade_type === 'import' ? 'default' as any : 'secondary' as any}>
                              {rate.trade_type === 'import' ? '進口' : '出口'}
                            </Badge>
                          </TableCell>
                          <TableCell className="text-center">
                            <div className="flex items-center justify-center gap-1">
                              <Percent className="h-3 w-3" />
                              {rate.rate}%
                            </div>
                          </TableCell>
                          <TableCell className="text-center">
                            {getRateTypeBadge(rate.rate_type)}
                          </TableCell>
                          <TableCell>
                            {getAgreementTypeBadge(rate.agreement_type)}
                          </TableCell>
                          <TableCell>
                            <div className="text-sm">
                              <p>{format(new Date(rate.valid_from), 'yyyy/MM/dd', { locale: zhTW })}</p>
                              {rate.valid_to && (
                                <p className="text-gray-500">
                                  至 {format(new Date(rate.valid_to), 'yyyy/MM/dd', { locale: zhTW })}
                                </p>
                              )}
                            </div>
                          </TableCell>
                          <TableCell className="text-center">
                            <Badge variant={rate.is_active ? 'success' as any : 'secondary' as any}>
                              {rate.is_active ? '有效' : '無效'}
                            </Badge>
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                )}
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="agreements">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <FileText className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">貿易協定管理功能開發中</p>
                  <p className="text-sm text-gray-400 mt-2">包含 FTA、CPTPP、RCEP 等貿易協定管理</p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="statistics">
            <Card>
              <CardContent className="p-6">
                <div className="text-center py-12">
                  <BarChart3 className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500">關稅統計分析功能開發中</p>
                  <p className="text-sm text-gray-400 mt-2">包含稅收分析、趨勢報表與預測功能</p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </div>
    </DashboardLayout>
  )
}