'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import * as z from 'zod'
import { Eye, EyeOff, Loader2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { useToast } from '@/components/ui/use-toast'
import authService, { LoginRequest } from '@/services/auth.service'
import { useAuthStore } from '@/store/auth.store'

const loginSchema = z.object({
  username: z.string().min(1, '請輸入使用者名稱'),
  password: z.string().min(1, '請輸入密碼'),
})

type LoginFormData = z.infer<typeof loginSchema>

export default function LoginPage() {
  const router = useRouter()
  const { toast } = useToast()
  const setUser = useAuthStore((state) => state.setUser)
  const [showPassword, setShowPassword] = useState(false)
  const [isLoading, setIsLoading] = useState(false)

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      username: '',
      password: '',
    },
  })

  const onSubmit = async (data: LoginFormData) => {
    setIsLoading(true)
    try {
      const response = await authService.login(data as LoginRequest)
      setUser({
        ...response.user,
        is_active: true,
        is_email_verified: true
      })
      
      toast({
        title: '登入成功',
        description: `歡迎回來，${response.user.full_name}！`,
      })
      
      // Redirect based on role
      switch (response.user.role) {
        case 'admin':
          router.push('/admin/dashboard')
          break
        case 'manager':
          router.push('/manager/dashboard')
          break
        case 'engineer':
          router.push('/engineer/inquiries')
          break
        case 'sales':
          router.push('/sales/inquiries')
          break
        default:
          router.push('/dashboard')
      }
    } catch (error: any) {
      toast({
        title: '登入失敗',
        description: error.response?.data?.message || '使用者名稱或密碼錯誤',
        variant: 'destructive',
      })
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-gray-100 to-gray-200">
      <Card className="w-full max-w-md">
        <CardHeader className="space-y-1">
          <CardTitle className="text-2xl font-bold text-center">
            FastenMind 系統登入
          </CardTitle>
          <CardDescription className="text-center">
            請輸入您的帳號密碼登入系統
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="username">使用者名稱</Label>
              <Input
                id="username"
                type="text"
                placeholder="請輸入使用者名稱或 Email"
                {...register('username')}
                disabled={isLoading}
              />
              {errors.username && (
                <p className="text-sm text-red-500">{errors.username.message}</p>
              )}
            </div>
            
            <div className="space-y-2">
              <Label htmlFor="password">密碼</Label>
              <div className="relative">
                <Input
                  id="password"
                  type={showPassword ? 'text' : 'password'}
                  placeholder="請輸入密碼"
                  {...register('password')}
                  disabled={isLoading}
                />
                <Button
                  type="button"
                  variant="ghost"
                  size="sm"
                  className="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent"
                  onClick={() => setShowPassword(!showPassword)}
                  disabled={isLoading}
                >
                  {showPassword ? (
                    <EyeOff className="h-4 w-4 text-gray-500" />
                  ) : (
                    <Eye className="h-4 w-4 text-gray-500" />
                  )}
                </Button>
              </div>
              {errors.password && (
                <p className="text-sm text-red-500">{errors.password.message}</p>
              )}
            </div>
            
            <Button
              type="submit"
              className="w-full"
              disabled={isLoading}
            >
              {isLoading ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  登入中...
                </>
              ) : (
                '登入'
              )}
            </Button>
          </form>
          
          <div className="mt-6 text-center text-sm text-gray-600">
            <p>測試帳號：</p>
            <p>admin / password123 (系統管理員)</p>
            <p>engineer1 / password123 (工程師)</p>
            <p>sales1 / password123 (業務人員)</p>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}