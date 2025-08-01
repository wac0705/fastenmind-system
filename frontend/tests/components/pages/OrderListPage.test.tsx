import { render, screen, fireEvent, waitFor } from '../../utils/test-utils'
import { mockServices, mockOrder } from '../../utils/test-utils'
import OrderListPage from '@/app/orders/page'

// Mock the order service
jest.mock('@/services/order.service', () => mockServices.order)

// Mock next/navigation
const mockPush = jest.fn()
jest.mock('next/navigation', () => ({
  useRouter: () => ({
    push: mockPush,
  }),
  useSearchParams: () => ({
    get: jest.fn(),
  }),
}))

describe('OrderListPage', () => {
  beforeEach(() => {
    jest.clearAllMocks()
    mockPush.mockClear()
  })

  it('renders order list correctly', async () => {
    const mockOrders = [
      { ...mockOrder, id: '1', order_no: 'ORD-2024-001', status: 'confirmed' },
      { ...mockOrder, id: '2', order_no: 'ORD-2024-002', status: 'in_production' },
    ]
    
    mockServices.order.getOrders.mockResolvedValue({ data: mockOrders })
    
    render(<OrderListPage />, {
      initialAuthState: {
        user: { role: 'engineer' },
        isAuthenticated: true,
      }
    })
    
    // Check page title
    expect(screen.getByText('訂單管理')).toBeInTheDocument()
    
    // Wait for orders to load
    await waitFor(() => {
      expect(screen.getByText('ORD-2024-001')).toBeInTheDocument()
      expect(screen.getByText('ORD-2024-002')).toBeInTheDocument()
    })
  })

  it('filters orders by status', async () => {
    const allOrders = [
      { ...mockOrder, id: '1', status: 'confirmed' },
      { ...mockOrder, id: '2', status: 'in_production' },
    ]
    const confirmedOrders = [
      { ...mockOrder, id: '1', status: 'confirmed' },
    ]
    
    mockServices.order.getOrders
      .mockResolvedValueOnce({ data: allOrders })
      .mockResolvedValueOnce({ data: confirmedOrders })
    
    render(<OrderListPage />, {
      initialAuthState: {
        user: { role: 'engineer' },
        isAuthenticated: true,
      }
    })
    
    // Wait for initial load
    await waitFor(() => {
      expect(screen.getByText('ORD-2024-001')).toBeInTheDocument()
    })
    
    // Click status filter
    const statusFilter = screen.getByTestId('status-filter')
    fireEvent.click(statusFilter)
    
    // Select confirmed status
    const confirmedOption = screen.getByText('已確認')
    fireEvent.click(confirmedOption)
    
    // Wait for filtered results
    await waitFor(() => {
      expect(mockServices.order.getOrders).toHaveBeenCalledWith({ status: 'confirmed' })
    })
  })

  it('searches orders by order number', async () => {
    const searchResults = [
      { ...mockOrder, id: '1', order_no: 'ORD-2024-001' },
    ]
    
    mockServices.order.getOrders
      .mockResolvedValueOnce({ data: [mockOrder] })
      .mockResolvedValueOnce({ data: searchResults })
    
    render(<OrderListPage />, {
      initialAuthState: {
        user: { role: 'engineer' },
        isAuthenticated: true,
      }
    })
    
    // Wait for initial load
    await waitFor(() => {
      expect(screen.getByDisplayValue('')).toBeInTheDocument()
    })
    
    // Enter search term
    const searchInput = screen.getByPlaceholderText('搜索訂單號、PO號碼或客戶名稱')
    fireEvent.change(searchInput, { target: { value: 'ORD-2024-001' } })
    fireEvent.keyDown(searchInput, { key: 'Enter', code: 'Enter' })
    
    // Wait for search results
    await waitFor(() => {
      expect(mockServices.order.getOrders).toHaveBeenCalledWith({ search: 'ORD-2024-001' })
    })
  })

  it('navigates to order detail when row is clicked', async () => {
    mockServices.order.getOrders.mockResolvedValue({ data: [mockOrder] })
    
    render(<OrderListPage />, {
      initialAuthState: {
        user: { role: 'engineer' },
        isAuthenticated: true,
      }
    })
    
    // Wait for orders to load
    await waitFor(() => {
      expect(screen.getByText(mockOrder.order_no)).toBeInTheDocument()
    })
    
    // Click on order row
    const orderRow = screen.getByTestId(`order-row-${mockOrder.id}`)
    fireEvent.click(orderRow)
    
    // Should navigate to order detail
    expect(mockPush).toHaveBeenCalledWith(`/orders/${mockOrder.id}`)
  })

  it('shows create order button for authorized roles', () => {
    mockServices.order.getOrders.mockResolvedValue({ data: [] })
    
    render(<OrderListPage />, {
      initialAuthState: {
        user: { role: 'sales' },
        isAuthenticated: true,
      }
    })
    
    // Should show create button for sales role
    expect(screen.getByText('建立訂單')).toBeInTheDocument()
  })

  it('hides create order button for unauthorized roles', () => {
    mockServices.order.getOrders.mockResolvedValue({ data: [] })
    
    render(<OrderListPage />, {
      initialAuthState: {
        user: { role: 'viewer' },
        isAuthentified: true,
      }
    })
    
    // Should not show create button for viewer role
    expect(screen.queryByText('建立訂單')).not.toBeInTheDocument()
  })

  it('exports orders to CSV', async () => {
    const mockOrders = [mockOrder]
    mockServices.order.getOrders.mockResolvedValue({ data: mockOrders })
    
    // Mock URL.createObjectURL and document.createElement
    global.URL.createObjectURL = jest.fn(() => 'blob:url')
    const mockClick = jest.fn()
    const mockLink = {
      href: '',
      download: '',
      click: mockClick,
    }
    jest.spyOn(document, 'createElement').mockReturnValue(mockLink as any)
    
    render(<OrderListPage />, {
      initialAuthState: {
        user: { role: 'manager' },
        isAuthenticated: true,
      }
    })
    
    // Wait for orders to load
    await waitFor(() => {
      expect(screen.getByText(mockOrder.order_no)).toBeInTheDocument()
    })
    
    // Click export button
    const exportButton = screen.getByTestId('export-button')
    fireEvent.click(exportButton)
    
    // Should trigger CSV download
    expect(mockClick).toHaveBeenCalled()
    
    // Cleanup
    global.URL.createObjectURL.mockRestore()
  })

  it('handles loading state', () => {
    // Mock a pending promise
    mockServices.order.getOrders.mockReturnValue(new Promise(() => {}))
    
    render(<OrderListPage />, {
      initialAuthState: {
        user: { role: 'engineer' },
        isAuthenticated: true,
      }
    })
    
    // Should show loading indicator
    expect(screen.getByText('載入中...')).toBeInTheDocument()
  })

  it('handles error state', async () => {
    const errorMessage = 'Failed to load orders'
    mockServices.order.getOrders.mockRejectedValue(new Error(errorMessage))
    
    render(<OrderListPage />, {
      initialAuthState: {
        user: { role: 'engineer' },
        isAuthenticated: true,
      }
    })
    
    // Wait for error to appear
    await waitFor(() => {
      expect(screen.getByText(errorMessage)).toBeInTheDocument()
    })
  })

  it('shows status badges correctly', async () => {
    const mockOrders = [
      { ...mockOrder, id: '1', status: 'pending' },
      { ...mockOrder, id: '2', status: 'confirmed' },
      { ...mockOrder, id: '3', status: 'in_production' },
      { ...mockOrder, id: '4', status: 'shipped' },
      { ...mockOrder, id: '5', status: 'completed' },
      { ...mockOrder, id: '6', status: 'cancelled' },
    ]
    
    mockServices.order.getOrders.mockResolvedValue({ data: mockOrders })
    
    render(<OrderListPage />, {
      initialAuthState: {
        user: { role: 'engineer' },
        isAuthenticated: true,
      }
    })
    
    // Wait for orders to load
    await waitFor(() => {
      expect(screen.getByText('待確認')).toBeInTheDocument()
      expect(screen.getByText('已確認')).toBeInTheDocument()
      expect(screen.getByText('生產中')).toBeInTheDocument()
      expect(screen.getByText('已出貨')).toBeInTheDocument()
      expect(screen.getByText('已完成')).toBeInTheDocument()
      expect(screen.getByText('已取消')).toBeInTheDocument()
    })
  })

  it('shows payment status badges correctly', async () => {
    const mockOrders = [
      { ...mockOrder, id: '1', payment_status: 'pending' },
      { ...mockOrder, id: '2', payment_status: 'partial' },
      { ...mockOrder, id: '3', payment_status: 'paid' },
    ]
    
    mockServices.order.getOrders.mockResolvedValue({ data: mockOrders })
    
    render(<OrderListPage />, {
      initialAuthState: {
        user: { role: 'engineer' },
        isAuthenticated: true,
      }
    })
    
    // Wait for orders to load
    await waitFor(() => {
      expect(screen.getByText('未付款')).toBeInTheDocument()
      expect(screen.getByText('部分付款')).toBeInTheDocument()
      expect(screen.getByText('已付款')).toBeInTheDocument()
    })
  })

  it('handles pagination', async () => {
    const page1Orders = Array.from({ length: 10 }, (_, i) => ({
      ...mockOrder,
      id: `${i + 1}`,
      order_no: `ORD-2024-${String(i + 1).padStart(3, '0')}`,
    }))
    
    mockServices.order.getOrders.mockResolvedValue({ 
      data: page1Orders,
      pagination: { total: 25, page: 1, limit: 10 }
    })
    
    render(<OrderListPage />, {
      initialAuthState: {
        user: { role: 'engineer' },
        isAuthenticated: true,
      }
    })
    
    // Wait for orders to load
    await waitFor(() => {
      expect(screen.getByText('ORD-2024-001')).toBeInTheDocument()
    })
    
    // Should show pagination controls
    expect(screen.getByText('1')).toBeInTheDocument() // Current page
    expect(screen.getByText('2')).toBeInTheDocument() // Next page
    expect(screen.getByText('3')).toBeInTheDocument() // Page 3
  })
})