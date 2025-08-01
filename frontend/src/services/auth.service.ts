import api, { tokenManager } from './api'

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  access_token: string
  refresh_token: string
  expires_in: number
  user: {
    id: string
    username: string
    email: string
    full_name: string
    role: string
    company_id: string
  }
}

export interface User {
  id: string
  username: string
  email: string
  full_name: string
  role: string
  company_id: string
  phone_number?: string
  is_active: boolean
  is_email_verified: boolean
  last_login_at?: string
}

class AuthService {
  async login(credentials: LoginRequest): Promise<LoginResponse> {
    const response = await api.post<LoginResponse>('/auth/login', credentials)
    const { access_token, refresh_token, user } = response.data
    
    // Store tokens
    tokenManager.setTokens(access_token, refresh_token)
    
    // Store user info
    if (typeof window !== 'undefined') {
      localStorage.setItem('user', JSON.stringify(user))
    }
    
    return response.data
  }
  
  async logout(): Promise<void> {
    // Clear tokens and user data
    tokenManager.clearTokens()
    if (typeof window !== 'undefined') {
      localStorage.removeItem('user')
    }
  }
  
  async refreshToken(): Promise<LoginResponse> {
    const refreshToken = tokenManager.getRefreshToken()
    if (!refreshToken) {
      throw new Error('No refresh token available')
    }
    
    const response = await api.post<LoginResponse>('/auth/refresh', {
      refresh_token: refreshToken,
    })
    
    const { access_token, refresh_token: newRefreshToken } = response.data
    tokenManager.setTokens(access_token, newRefreshToken)
    
    return response.data
  }
  
  getCurrentUser(): User | null {
    if (typeof window !== 'undefined') {
      const userStr = localStorage.getItem('user')
      if (userStr) {
        try {
          return JSON.parse(userStr) as User
        } catch {
          return null
        }
      }
    }
    return null
  }
  
  isAuthenticated(): boolean {
    return !!tokenManager.getAccessToken()
  }
  
  hasRole(role: string): boolean {
    const user = this.getCurrentUser()
    return user?.role === role
  }
  
  hasAnyRole(roles: string[]): boolean {
    const user = this.getCurrentUser()
    return user ? roles.includes(user.role) : false
  }
}

export const authService = new AuthService()
export default authService