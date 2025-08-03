'use client'

import { ReactNode, useState } from 'react'
import Link from 'next/link'
import { useRouter, usePathname } from 'next/navigation'
import { 
  Home, 
  Package, 
  Users, 
  FileText, 
  Settings, 
  LogOut,
  Menu,
  X,
  ChevronDown,
  Calculator,
  Wrench,
  DollarSign,
  Globe,
  Workflow,
  Factory,
  CheckCircle,
  BarChart3,
  Truck,
  Smartphone,
  Zap,
  Brain,
  Lightbulb,
  Search,
  Layers,
  Shield,
  Activity,
  Link2,
  Webhook,
  RefreshCw as Sync,
  Key,
  Server,
  History,
  CreditCard,
  Bell
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { useAuthStore } from '@/store/auth.store'
import authService from '@/services/auth.service'
import { cn } from '@/lib/utils'

interface DashboardLayoutProps {
  children: ReactNode
}

interface MenuItem {
  title: string
  href: string
  icon: any
  roles?: string[]
  children?: MenuItem[]
}

const menuItems: MenuItem[] = [
  {
    title: '首頁',
    href: '/dashboard',
    icon: Home,
  },
  {
    title: '詢價管理',
    href: '/inquiries',
    icon: FileText,
    roles: ['admin', 'manager', 'engineer', 'sales'],
    children: [
      { title: '詢價單列表', href: '/inquiries', icon: FileText },
      { title: '新增詢價單', href: '/inquiries/new', icon: FileText, roles: ['sales'] },
      { title: '待處理詢價', href: '/inquiries/pending', icon: FileText, roles: ['engineer'] },
    ],
  },
  {
    title: '報價管理',
    href: '/quotes',
    icon: DollarSign,
    roles: ['admin', 'manager', 'engineer'],
    children: [
      { title: '報價單列表', href: '/quotes', icon: DollarSign },
      { title: '待審核報價', href: '/quotes/review', icon: DollarSign, roles: ['manager'] },
    ],
  },
  {
    title: '工程管理',
    href: '/engineering',
    icon: Wrench,
    roles: ['admin', 'engineer'],
    children: [
      { title: '製程管理', href: '/processes', icon: Wrench },
      { title: '設備管理', href: '/equipment', icon: Wrench },
      { title: '成本計算', href: '/cost-calculator', icon: Calculator },
    ],
  },
  {
    title: '國際貿易',
    href: '/trade',
    icon: Globe,
    roles: ['admin', 'manager', 'engineer'],
    children: [
      { title: '貿易概覽', href: '/trade', icon: Globe },
      { title: '運輸管理', href: '/trade/shipments', icon: Truck },
      { title: '關稅管理', href: '/trade/tariffs', icon: Calculator },
      { title: '信用狀', href: '/trade/letter-of-credits', icon: CreditCard },
      { title: '合規檢查', href: '/trade/compliance', icon: Shield },
      { title: '匯率管理', href: '/trade/exchange-rates', icon: DollarSign },
      { title: '貿易分析', href: '/trade/analytics', icon: BarChart3 },
    ],
  },
  {
    title: '訂單管理',
    href: '/orders',
    icon: Package,
    roles: ['admin', 'manager', 'sales'],
    children: [
      { title: '訂單列表', href: '/orders', icon: Package },
      { title: '新增訂單', href: '/orders/new', icon: Package, roles: ['sales'] },
    ],
  },
  {
    title: '庫存管理',
    href: '/inventory',
    icon: Package,
    roles: ['admin', 'manager', 'engineer'],
    children: [
      { title: '庫存列表', href: '/inventory', icon: Package },
      { title: '新增品項', href: '/inventory/new', icon: Package },
      { title: '盤點作業', href: '/inventory/stock-take', icon: Package },
    ],
  },
  {
    title: '財務管理',
    href: '/finance',
    icon: DollarSign,
    roles: ['admin', 'manager'],
    children: [
      { title: '財務概覽', href: '/finance', icon: DollarSign },
      { title: '發票管理', href: '/finance/invoices', icon: FileText },
      { title: '付款記錄', href: '/finance/payments', icon: DollarSign },
      { title: '費用管理', href: '/finance/expenses', icon: DollarSign },
      { title: '應收帳款', href: '/finance/receivables', icon: DollarSign },
      { title: '財務報表', href: '/finance/reports', icon: DollarSign },
    ],
  },
  {
    title: '生產管理',
    href: '/production',
    icon: Factory,
    roles: ['admin', 'manager', 'engineer'],
    children: [
      { title: '生產概覽', href: '/production', icon: Factory },
      { title: '生產訂單', href: '/production/orders', icon: Package },
      { title: '生產任務', href: '/production/tasks', icon: Users },
      { title: '工作站', href: '/production/stations', icon: Settings },
      { title: '工藝路線', href: '/production/routes', icon: Workflow },
      { title: '品質管理', href: '/production/quality', icon: CheckCircle },
      { title: '生產報表', href: '/production/reports', icon: BarChart3 },
    ],
  },
  {
    title: '供應商管理',
    href: '/suppliers',
    icon: Truck,
    roles: ['admin', 'manager', 'engineer'],
    children: [
      { title: '供應商清單', href: '/suppliers', icon: Truck },
      { title: '新增供應商', href: '/suppliers/new', icon: Truck, roles: ['admin', 'manager'] },
      { title: '採購訂單', href: '/suppliers/purchase-orders', icon: FileText },
      { title: '績效評估', href: '/suppliers/evaluations', icon: BarChart3 },
      { title: '風險管理', href: '/suppliers/risk-management', icon: Settings },
    ],
  },
  {
    title: '客戶管理',
    href: '/customers',
    icon: Users,
    roles: ['admin', 'manager', 'sales'],
  },
  {
    title: '報表中心',
    href: '/reports',
    icon: BarChart3,
    roles: ['admin', 'manager', 'engineer'],
    children: [
      { title: '報表概覽', href: '/reports', icon: BarChart3 },
      { title: '新增報表', href: '/reports/new', icon: FileText, roles: ['admin', 'manager'] },
      { title: '範本庫', href: '/reports/templates', icon: FileText },
      { title: '儀表板', href: '/reports/dashboards', icon: BarChart3 },
      { title: '執行記錄', href: '/reports/executions', icon: Settings },
      { title: 'KPI 管理', href: '/reports/kpis', icon: BarChart3 },
      { title: '資料來源', href: '/reports/data-sources', icon: Settings },
    ],
  },
  {
    title: '自動化管理',
    href: '/workflows',
    icon: Workflow,
    roles: ['admin', 'manager'],
    children: [
      { title: '工作流程', href: '/workflows', icon: Workflow },
      { title: 'N8N 整合', href: '/workflows?tab=templates', icon: Workflow },
    ],
  },
  {
    title: '行動應用',
    href: '/mobile',
    icon: Smartphone,
    roles: ['admin', 'manager'],
  },
  {
    title: '整合功能',
    href: '/integrations',
    icon: Link2,
    roles: ['admin', 'manager'],
    children: [
      { title: '整合概覽', href: '/integrations', icon: Link2 },
      { title: '整合管理', href: '/integrations/manage', icon: Settings },
      { title: 'Webhook', href: '/integrations/webhooks', icon: Webhook },
      { title: '同步任務', href: '/integrations/sync-jobs', icon: Sync },
      { title: 'API 金鑰', href: '/integrations/api-keys', icon: Key },
      { title: '外部系統', href: '/integrations/external-systems', icon: Server },
      { title: '整合模板', href: '/integrations/templates', icon: FileText },
      { title: '記錄檔', href: '/integrations/logs', icon: History },
    ],
  },
  {
    title: '進階功能',
    href: '/advanced',
    icon: Zap,
    roles: ['admin', 'manager'],
    children: [
      { title: '功能概覽', href: '/advanced', icon: Zap },
      { title: 'AI 助手', href: '/advanced/ai-assistant', icon: Brain },
      { title: '智能推薦', href: '/advanced/recommendations', icon: Lightbulb },
      { title: '高級搜索', href: '/advanced/search', icon: Search },
      { title: '批量操作', href: '/advanced/batch-operations', icon: Layers },
      { title: '安全監控', href: '/advanced/security', icon: Shield },
      { title: '效能監控', href: '/advanced/performance', icon: Activity },
      { title: '系統管理', href: '/advanced/system', icon: Settings },
    ],
  },
  {
    title: '系統管理',
    href: '/system',
    icon: Settings,
    roles: ['admin'],
    children: [
      { title: '使用者管理', href: '/system/users', icon: Users },
      { title: '角色權限', href: '/system/roles', icon: Shield },
      { title: '系統設定', href: '/system/settings', icon: Settings },
      { title: '操作記錄', href: '/system/audit-logs', icon: History },
      { title: '系統通知', href: '/system/notifications', icon: Bell },
      { title: '系統健康', href: '/system/health', icon: Activity },
    ],
  },
  {
    title: '系統設定',
    href: '/settings',
    icon: Settings,
    roles: ['admin'],
    children: [
      { title: '公司管理', href: '/settings/companies', icon: Settings },
      { title: '使用者管理', href: '/settings/users', icon: Users },
      { title: '分派規則', href: '/settings/assignment-rules', icon: Settings },
      { title: '系統管理', href: '/settings/system', icon: Settings },
    ],
  },
]

export default function DashboardLayout({ children }: DashboardLayoutProps) {
  const router = useRouter()
  const pathname = usePathname()
  const { user, clearAuth } = useAuthStore()
  const [sidebarOpen, setSidebarOpen] = useState(true)
  const [expandedMenus, setExpandedMenus] = useState<string[]>([])

  const handleLogout = async () => {
    await authService.logout()
    clearAuth()
    router.push('/login')
  }

  const toggleMenu = (title: string) => {
    setExpandedMenus((prev) =>
      prev.includes(title)
        ? prev.filter((t) => t !== title)
        : [...prev, title]
    )
  }

  const filterMenuByRole = (items: MenuItem[]): MenuItem[] => {
    return items.filter((item) => {
      if (!item.roles || item.roles.length === 0) return true
      return user && item.roles.includes(user.role)
    })
  }

  const isActiveRoute = (href: string) => {
    return pathname === href || pathname.startsWith(href + '/')
  }

  return (
    <div className="flex h-screen bg-gray-100">
      {/* Sidebar */}
      <aside
        className={cn(
          'bg-white shadow-md transition-all duration-300',
          sidebarOpen ? 'w-64' : 'w-16'
        )}
      >
        <div className="flex h-full flex-col">
          {/* Logo */}
          <div className="flex h-16 items-center justify-between px-4 border-b">
            {sidebarOpen && (
              <h1 className="text-xl font-bold text-gray-800">FastenMind</h1>
            )}
            <Button
              variant="ghost"
              size="icon"
              onClick={() => setSidebarOpen(!sidebarOpen)}
            >
              {sidebarOpen ? <X size={20} /> : <Menu size={20} />}
            </Button>
          </div>

          {/* Navigation */}
          <nav className="flex-1 overflow-y-auto p-4">
            <ul className="space-y-2">
              {filterMenuByRole(menuItems).map((item) => (
                <li key={item.href}>
                  {item.children ? (
                    <div>
                      <button
                        onClick={() => toggleMenu(item.title)}
                        className={cn(
                          'flex w-full items-center justify-between rounded-lg px-3 py-2 text-sm font-medium transition-colors',
                          isActiveRoute(item.href)
                            ? 'bg-gray-100 text-gray-900'
                            : 'text-gray-700 hover:bg-gray-100'
                        )}
                      >
                        <div className="flex items-center">
                          <item.icon className="h-5 w-5" />
                          {sidebarOpen && (
                            <span className="ml-3">{item.title}</span>
                          )}
                        </div>
                        {sidebarOpen && (
                          <ChevronDown
                            className={cn(
                              'h-4 w-4 transition-transform',
                              expandedMenus.includes(item.title) && 'rotate-180'
                            )}
                          />
                        )}
                      </button>
                      {sidebarOpen && expandedMenus.includes(item.title) && (
                        <ul className="mt-2 ml-8 space-y-1">
                          {filterMenuByRole(item.children).map((child) => (
                            <li key={child.href}>
                              <Link
                                href={child.href}
                                className={cn(
                                  'flex items-center rounded-lg px-3 py-2 text-sm transition-colors',
                                  isActiveRoute(child.href)
                                    ? 'bg-gray-100 text-gray-900 font-medium'
                                    : 'text-gray-600 hover:bg-gray-50'
                                )}
                              >
                                {child.title}
                              </Link>
                            </li>
                          ))}
                        </ul>
                      )}
                    </div>
                  ) : (
                    <Link
                      href={item.href}
                      className={cn(
                        'flex items-center rounded-lg px-3 py-2 text-sm font-medium transition-colors',
                        isActiveRoute(item.href)
                          ? 'bg-gray-100 text-gray-900'
                          : 'text-gray-700 hover:bg-gray-100'
                      )}
                    >
                      <item.icon className="h-5 w-5" />
                      {sidebarOpen && <span className="ml-3">{item.title}</span>}
                    </Link>
                  )}
                </li>
              ))}
            </ul>
          </nav>

          {/* User info & Logout */}
          <div className="border-t p-4">
            <div className="flex items-center">
              <div className="flex-1">
                {sidebarOpen && user && (
                  <div>
                    <p className="text-sm font-medium text-gray-900">
                      {user.full_name}
                    </p>
                    <p className="text-xs text-gray-500">{user.role}</p>
                  </div>
                )}
              </div>
              <Button
                variant="ghost"
                size="icon"
                onClick={handleLogout}
                title="登出"
              >
                <LogOut className="h-5 w-5" />
              </Button>
            </div>
          </div>
        </div>
      </aside>

      {/* Main content */}
      <main className="flex-1 overflow-y-auto">
        <div className="p-8">{children}</div>
      </main>
    </div>
  )
}