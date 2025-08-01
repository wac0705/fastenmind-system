'use client'

import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { useRouter } from 'next/navigation'
import DashboardLayout from '@/components/layout/DashboardLayout'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { 
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Badge } from '@/components/ui/badge'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { 
  Search,
  Filter,
  Package,
  Plus,
  AlertTriangle,
  TrendingDown,
  Archive,
  Warehouse,
  DollarSign,
  BarChart3,
  AlertCircle
} from 'lucide-react'
import inventoryService, { Inventory, InventoryStats } from '@/services/inventory.service'
import Pagination from '@/components/common/Pagination'

export default function InventoryPage() {
  const router = useRouter()
  const [page, setPage] = useState(1)
  const [search, setSearch] = useState('')
  const [categoryFilter, setCategoryFilter] = useState<string>('')
  const [warehouseFilter, setWarehouseFilter] = useState<string>('')
  const [statusFilter, setStatusFilter] = useState<string>('')
  const [lowStockOnly, setLowStockOnly] = useState(false)

  // Fetch inventory items
  const { data: inventoryData, isLoading: isLoadingInventory } = useQuery({
    queryKey: ['inventory', page, search, categoryFilter, warehouseFilter, statusFilter, lowStockOnly],
    queryFn: () => inventoryService.list({
      page,
      page_size: 20,
      search,
      category: categoryFilter,
      warehouse_id: warehouseFilter,
      status: statusFilter,
      low_stock: lowStockOnly,
    }),
  })

  // Fetch stats
  const { data: stats } = useQuery({
    queryKey: ['inventory-stats'],
    queryFn: () => inventoryService.getStats(),
  })

  // Fetch warehouses for filter
  const { data: warehouses } = useQuery({
    queryKey: ['warehouses'],
    queryFn: () => inventoryService.listWarehouses(),
  })

  const getStockStatusBadge = (item: Inventory) => {
    if (item.current_stock <= 0) {
      return <Badge variant="destructive">缺貨</Badge>
    } else if (item.current_stock <= item.min_stock) {
      return <Badge variant="warning">低庫存</Badge>
    } else if (item.max_stock > 0 && item.current_stock >= item.max_stock) {
      return <Badge variant="info">超量</Badge>
    }
    return <Badge variant="success">正常</Badge>
  }

  const getCategoryBadge = (category: string) => {
    const categoryConfig: Record<string, { label: string; variant: any }> = {
      raw_material: { label: '原材料', variant: 'secondary' },
      semi_finished: { label: '半成品', variant: 'warning' },
      finished_goods: { label: '成品', variant: 'success' },
    }

    const config = categoryConfig[category] || { label: category, variant: 'default' }
    return <Badge variant={config.variant as any}>{config.label}</Badge>
  }

  const handleSearch = (value: string) => {
    setSearch(value)
    setPage(1)
  }

  return (
    <DashboardLayout>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex justify-between items-center">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">庫存管理</h1>
            <p className="mt-2 text-gray-600">管理產品庫存與倉儲</p>
          </div>
          <div className="flex gap-2">
            <Button variant="outline" onClick={() => router.push('/inventory/stock-take')}>
              <BarChart3 className="mr-2 h-4 w-4" />
              盤點作業
            </Button>
            <Button onClick={() => router.push('/inventory/new')}>
              <Plus className="mr-2 h-4 w-4" />
              新增品項
            </Button>
          </div>
        </div>

        {/* Stats Cards */}
        {stats && (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">總品項數</CardTitle>
                <Package className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{stats.total_items}</div>
                <p className="text-xs text-muted-foreground">管理中的庫存品項</p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">庫存總值</CardTitle>
                <DollarSign className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">${stats.total_value.toLocaleString()}</div>
                <p className="text-xs text-muted-foreground">所有庫存價值</p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">低庫存</CardTitle>
                <TrendingDown className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold text-orange-600">{stats.low_stock_items}</div>
                <p className="text-xs text-muted-foreground">需要補貨的品項</p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">庫存警示</CardTitle>
                <AlertCircle className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold text-red-600">{stats.active_alerts}</div>
                <p className="text-xs text-muted-foreground">待處理警示</p>
              </CardContent>
            </Card>
          </div>
        )}

        {/* Filters */}
        <Card>
          <CardHeader>
            <CardTitle>庫存列表</CardTitle>
            <CardDescription>查看和管理所有庫存品項</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="flex flex-col sm:flex-row gap-4 mb-6">
              <div className="relative flex-1">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
                <Input
                  placeholder="搜尋 SKU、料號或品名..."
                  value={search}
                  onChange={(e) => handleSearch(e.target.value)}
                  className="pl-10"
                />
              </div>
              <Select value={categoryFilter} onValueChange={(value) => { setCategoryFilter(value); setPage(1); }}>
                <SelectTrigger className="w-full sm:w-[180px]">
                  <SelectValue placeholder="類別" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="">全部類別</SelectItem>
                  <SelectItem value="raw_material">原材料</SelectItem>
                  <SelectItem value="semi_finished">半成品</SelectItem>
                  <SelectItem value="finished_goods">成品</SelectItem>
                </SelectContent>
              </Select>
              <Select value={warehouseFilter} onValueChange={(value) => { setWarehouseFilter(value); setPage(1); }}>
                <SelectTrigger className="w-full sm:w-[200px]">
                  <SelectValue placeholder="倉庫" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="">全部倉庫</SelectItem>
                  {warehouses?.map((warehouse) => (
                    <SelectItem key={warehouse.id} value={warehouse.id}>
                      {warehouse.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              <Button
                variant={lowStockOnly ? "default" : "outline"}
                onClick={() => { setLowStockOnly(!lowStockOnly); setPage(1); }}
                className="whitespace-nowrap"
              >
                <AlertTriangle className="mr-2 h-4 w-4" />
                低庫存
              </Button>
            </div>

            {/* Inventory Table */}
            <div className="rounded-md border">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>SKU</TableHead>
                    <TableHead>品名</TableHead>
                    <TableHead>料號</TableHead>
                    <TableHead>類別</TableHead>
                    <TableHead>倉庫</TableHead>
                    <TableHead className="text-right">現有庫存</TableHead>
                    <TableHead className="text-right">可用庫存</TableHead>
                    <TableHead className="text-right">單位成本</TableHead>
                    <TableHead>狀態</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {isLoadingInventory ? (
                    <TableRow>
                      <TableCell colSpan={9} className="text-center py-8">
                        載入中...
                      </TableCell>
                    </TableRow>
                  ) : inventoryData?.data.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={9} className="text-center py-8 text-gray-500">
                        沒有找到庫存品項
                      </TableCell>
                    </TableRow>
                  ) : (
                    inventoryData?.data.map((item) => (
                      <TableRow
                        key={item.id}
                        className="cursor-pointer hover:bg-gray-50"
                        onClick={() => router.push(`/inventory/${item.id}`)}
                      >
                        <TableCell className="font-medium">{item.sku}</TableCell>
                        <TableCell>{item.name}</TableCell>
                        <TableCell>{item.part_no}</TableCell>
                        <TableCell>{getCategoryBadge(item.category)}</TableCell>
                        <TableCell>
                          <div className="flex items-center gap-1">
                            <Warehouse className="h-3 w-3 text-gray-400" />
                            <span className="text-sm">{item.warehouse?.name || '-'}</span>
                          </div>
                        </TableCell>
                        <TableCell className="text-right">
                          <span className={item.current_stock <= item.min_stock ? 'text-orange-600 font-medium' : ''}>
                            {item.current_stock.toLocaleString()} {item.unit}
                          </span>
                        </TableCell>
                        <TableCell className="text-right">
                          {item.available_stock.toLocaleString()} {item.unit}
                        </TableCell>
                        <TableCell className="text-right">
                          <div className="flex items-center justify-end gap-1">
                            <span className="text-xs text-gray-500">{item.currency}</span>
                            <span>{item.average_cost.toFixed(2)}</span>
                          </div>
                        </TableCell>
                        <TableCell>{getStockStatusBadge(item)}</TableCell>
                      </TableRow>
                    ))
                  )}
                </TableBody>
              </Table>
            </div>

            {/* Pagination */}
            {inventoryData && inventoryData.pagination.total > 0 && (
              <div className="mt-4">
                <Pagination
                  currentPage={page}
                  totalPages={Math.ceil(inventoryData.pagination.total / 20)}
                  onPageChange={setPage}
                />
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </DashboardLayout>
  )
}