import { render, screen, fireEvent, waitFor } from '../utils/test-utils'
import { mockServices, mockApiResponses } from '../utils/test-utils'
import LoginPage from '@/app/login/page'

// Mock the auth service
jest.mock('@/services/auth.service', () => mockServices.auth)

// Mock next/navigation
const mockPush = jest.fn()
jest.mock('next/navigation', () => ({
  useRouter: () => ({
    push: mockPush,
  }),
}))

describe('LoginPage', () => {
  beforeEach(() => {
    jest.clearAllMocks()
    mockPush.mockClear()
  })

  it('renders login form correctly', () => {
    render(<LoginPage />)
    
    expect(screen.getByText('登入 FastenMind')).toBeInTheDocument()
    expect(screen.getByLabelText('帳號')).toBeInTheDocument()
    expect(screen.getByLabelText('密碼')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: '登入' })).toBeInTheDocument()
  })

  it('validates required fields', async () => {
    render(<LoginPage />)
    
    const loginButton = screen.getByRole('button', { name: '登入' })
    fireEvent.click(loginButton)
    
    await waitFor(() => {
      expect(screen.getByText('請輸入帳號')).toBeInTheDocument()
      expect(screen.getByText('請輸入密碼')).toBeInTheDocument()
    })
  })

  it('validates minimum password length', async () => {
    render(<LoginPage />)
    
    const usernameInput = screen.getByLabelText('帳號')
    const passwordInput = screen.getByLabelText('密碼')
    const loginButton = screen.getByRole('button', { name: '登入' })
    
    fireEvent.change(usernameInput, { target: { value: 'testuser' } })
    fireEvent.change(passwordInput, { target: { value: '123' } })
    fireEvent.click(loginButton)
    
    await waitFor(() => {
      expect(screen.getByText('密碼至少需要6個字元')).toBeInTheDocument()
    })
  })

  it('submits form with valid credentials', async () => {
    mockServices.auth.login.mockResolvedValue(mockApiResponses.login)
    
    render(<LoginPage />)
    
    const usernameInput = screen.getByLabelText('帳號')
    const passwordInput = screen.getByLabelText('密碼')
    const loginButton = screen.getByRole('button', { name: '登入' })
    
    fireEvent.change(usernameInput, { target: { value: 'testuser' } })
    fireEvent.change(passwordInput, { target: { value: 'password123' } })
    fireEvent.click(loginButton)
    
    await waitFor(() => {
      expect(mockServices.auth.login).toHaveBeenCalledWith({
        username: 'testuser',
        password: 'password123',
      })
    })
  })

  it('handles login success and redirects', async () => {
    mockServices.auth.login.mockResolvedValue(mockApiResponses.login)
    
    render(<LoginPage />)
    
    const usernameInput = screen.getByLabelText('帳號')
    const passwordInput = screen.getByLabelText('密碼')
    const loginButton = screen.getByRole('button', { name: '登入' })
    
    fireEvent.change(usernameInput, { target: { value: 'testuser' } })
    fireEvent.change(passwordInput, { target: { value: 'password123' } })
    fireEvent.click(loginButton)
    
    await waitFor(() => {
      expect(mockPush).toHaveBeenCalledWith('/dashboard')
    })
  })

  it('handles login error', async () => {
    const errorMessage = 'Invalid credentials'
    mockServices.auth.login.mockRejectedValue(new Error(errorMessage))
    
    render(<LoginPage />)
    
    const usernameInput = screen.getByLabelText('帳號')
    const passwordInput = screen.getByLabelText('密碼')
    const loginButton = screen.getByRole('button', { name: '登入' })
    
    fireEvent.change(usernameInput, { target: { value: 'testuser' } })
    fireEvent.change(passwordInput, { target: { value: 'wrongpassword' } })
    fireEvent.click(loginButton)
    
    await waitFor(() => {
      expect(screen.getByText(errorMessage)).toBeInTheDocument()
    })
  })

  it('shows loading state during login', async () => {
    // Create a promise that we can control
    let resolveLogin: (value: any) => void
    const loginPromise = new Promise((resolve) => {
      resolveLogin = resolve
    })
    
    mockServices.auth.login.mockReturnValue(loginPromise)
    
    render(<LoginPage />)
    
    const usernameInput = screen.getByLabelText('帳號')
    const passwordInput = screen.getByLabelText('密碼')
    const loginButton = screen.getByRole('button', { name: '登入' })
    
    fireEvent.change(usernameInput, { target: { value: 'testuser' } })
    fireEvent.change(passwordInput, { target: { value: 'password123' } })
    fireEvent.click(loginButton)
    
    // Should show loading state
    expect(screen.getByText('登入中...')).toBeInTheDocument()
    expect(loginButton).toBeDisabled()
    
    // Resolve the login
    resolveLogin!(mockApiResponses.login)
    
    await waitFor(() => {
      expect(screen.queryByText('登入中...')).not.toBeInTheDocument()
    })
  })

  it('toggles password visibility', () => {
    render(<LoginPage />)
    
    const passwordInput = screen.getByLabelText('密碼')
    const toggleButton = screen.getByRole('button', { name: '顯示密碼' })
    
    // Initially password should be hidden
    expect(passwordInput).toHaveAttribute('type', 'password')
    
    // Click to show password
    fireEvent.click(toggleButton)
    expect(passwordInput).toHaveAttribute('type', 'text')
    
    // Click to hide password again
    fireEvent.click(toggleButton)
    expect(passwordInput).toHaveAttribute('type', 'password')
  })

  it('handles Enter key submission', async () => {
    mockServices.auth.login.mockResolvedValue(mockApiResponses.login)
    
    render(<LoginPage />)
    
    const usernameInput = screen.getByLabelText('帳號')
    const passwordInput = screen.getByLabelText('密碼')
    
    fireEvent.change(usernameInput, { target: { value: 'testuser' } })
    fireEvent.change(passwordInput, { target: { value: 'password123' } })
    
    // Press Enter in password field
    fireEvent.keyDown(passwordInput, { key: 'Enter', code: 'Enter' })
    
    await waitFor(() => {
      expect(mockServices.auth.login).toHaveBeenCalledWith({
        username: 'testuser',
        password: 'password123',
      })
    })
  })

  it('focuses on username input on mount', () => {
    render(<LoginPage />)
    
    const usernameInput = screen.getByLabelText('帳號')
    expect(document.activeElement).toBe(usernameInput)
  })
})