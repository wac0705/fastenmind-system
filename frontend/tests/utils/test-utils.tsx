import React, { ReactElement } from 'react'
import { render, RenderOptions } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { useAuthStore } from '@/store/auth.store'

// Mock data
export const mockUser = {
  id: '123e4567-e89b-12d3-a456-426614174000',
  username: 'testuser',
  email: 'test@example.com',
  full_name: 'Test User',
  role: 'engineer',
  company_id: '123e4567-e89b-12d3-a456-426614174001',
  is_active: true,
}

export const mockCompany = {
  id: '123e4567-e89b-12d3-a456-426614174001',
  code: 'TEST001',
  name: 'Test Company',
  type: 'headquarters',
  country: 'US',
}

export const mockQuote = {
  id: '123e4567-e89b-12d3-a456-426614174002',
  quote_no: 'QUO-2024-001',
  inquiry_id: '123e4567-e89b-12d3-a456-426614174003',
  status: 'draft',
  total_amount: 10000.00,
  currency: 'USD',
  valid_until: '2024-12-31',
  delivery_days: 30,
  payment_terms: 'T/T 30 days',
  created_at: '2024-01-01T00:00:00Z',
  customer: {
    id: '123e4567-e89b-12d3-a456-426614174004',
    name: 'Test Customer',
    email: 'customer@example.com',
  },
}

export const mockInquiry = {
  id: '123e4567-e89b-12d3-a456-426614174003',
  inquiry_no: 'INQ-2024-001',
  status: 'assigned',
  product_name: 'Hex Bolt M8x20',
  product_category: 'Bolts',
  quantity: 10000,
  unit: 'pcs',
  required_date: '2024-12-31',
  customer: mockQuote.customer,
}

export const mockOrder = {
  id: '123e4567-e89b-12d3-a456-426614174005',
  order_no: 'ORD-2024-001',
  po_number: 'PO-2024-001',
  status: 'confirmed',
  total_amount: 10000.00,
  currency: 'USD',
  delivery_date: '2024-12-31',
  payment_status: 'partial',
  customer: mockQuote.customer,
  created_at: '2024-01-01T00:00:00Z',
}

// Custom render function with providers
interface CustomRenderOptions extends Omit<RenderOptions, 'wrapper'> {
  queryClient?: QueryClient
  initialAuthState?: {
    user?: typeof mockUser | null
    isAuthenticated?: boolean
    token?: string | null
  }
}

function AllTheProviders({ 
  children, 
  queryClient,
  initialAuthState 
}: { 
  children: React.ReactNode
  queryClient?: QueryClient
  initialAuthState?: CustomRenderOptions['initialAuthState']
}) {
  const client = queryClient || new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
      },
      mutations: {
        retry: false,
      },
    },
  })

  // Set initial auth state if provided
  React.useEffect(() => {
    if (initialAuthState) {
      const { setUser, setAuthenticated, setToken } = useAuthStore.getState()
      if (initialAuthState.user !== undefined) {
        setUser(initialAuthState.user)
      }
      if (initialAuthState.isAuthenticated !== undefined) {
        setAuthenticated(initialAuthState.isAuthenticated)
      }
      if (initialAuthState.token !== undefined) {
        setToken(initialAuthState.token)
      }
    }
  }, [initialAuthState])

  return (
    <QueryClientProvider client={client}>
      {children}
    </QueryClientProvider>
  )
}

const customRender = (
  ui: ReactElement,
  options?: CustomRenderOptions
) => {
  const { queryClient, initialAuthState, ...renderOptions } = options || {}
  
  return render(ui, {
    wrapper: ({ children }) => (
      <AllTheProviders 
        queryClient={queryClient} 
        initialAuthState={initialAuthState}
      >
        {children}
      </AllTheProviders>
    ),
    ...renderOptions,
  })
}

// Mock API responses
export const mockApiResponses = {
  // Auth responses
  login: {
    user: mockUser,
    access_token: 'mock-access-token',
    refresh_token: 'mock-refresh-token',
    expires_in: 3600,
  },
  
  // Quote responses
  quotes: [mockQuote],
  quote: mockQuote,
  
  // Inquiry responses
  inquiries: [mockInquiry],
  inquiry: mockInquiry,
  
  // Order responses
  orders: [mockOrder],
  order: mockOrder,
}

// Mock service functions
export const mockServices = {
  auth: {
    login: jest.fn(() => Promise.resolve(mockApiResponses.login)),
    logout: jest.fn(() => Promise.resolve()),
    refreshToken: jest.fn(() => Promise.resolve(mockApiResponses.login)),
    getProfile: jest.fn(() => Promise.resolve(mockUser)),
  },
  
  quote: {
    getQuotes: jest.fn(() => Promise.resolve(mockApiResponses.quotes)),
    getQuote: jest.fn(() => Promise.resolve(mockApiResponses.quote)),
    createQuote: jest.fn(() => Promise.resolve(mockApiResponses.quote)),
    updateQuote: jest.fn(() => Promise.resolve(mockApiResponses.quote)),
    deleteQuote: jest.fn(() => Promise.resolve()),
  },
  
  inquiry: {
    getInquiries: jest.fn(() => Promise.resolve(mockApiResponses.inquiries)),
    getInquiry: jest.fn(() => Promise.resolve(mockApiResponses.inquiry)),
    createInquiry: jest.fn(() => Promise.resolve(mockApiResponses.inquiry)),
    updateInquiry: jest.fn(() => Promise.resolve(mockApiResponses.inquiry)),
  },
  
  order: {
    getOrders: jest.fn(() => Promise.resolve(mockApiResponses.orders)),
    getOrder: jest.fn(() => Promise.resolve(mockApiResponses.order)),
    createOrder: jest.fn(() => Promise.resolve(mockApiResponses.order)),
    updateOrder: jest.fn(() => Promise.resolve(mockApiResponses.order)),
  },
}

// Utility functions for testing
export const waitForLoadingToFinish = () => {
  return new Promise((resolve) => {
    setTimeout(resolve, 0)
  })
}

export const createMockQueryClient = () => {
  return new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
        gcTime: Infinity,
      },
      mutations: {
        retry: false,
      },
    },
  })
}

export const mockConsoleError = () => {
  const spy = jest.spyOn(console, 'error').mockImplementation(() => {})
  
  return {
    restore: () => spy.mockRestore(),
    calls: spy.mock.calls,
  }
}

export const mockConsoleWarn = () => {
  const spy = jest.spyOn(console, 'warn').mockImplementation(() => {})
  
  return {
    restore: () => spy.mockRestore(),
    calls: spy.mock.calls,
  }
}

// Re-export everything
export * from '@testing-library/react'
export { customRender as render }