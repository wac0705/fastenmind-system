'use client'

import { useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { useAuthStore } from '@/store/auth.store'
import authService from '@/services/auth.service'

export default function DashboardPage() {
  const router = useRouter()
  const { setUser, user } = useAuthStore()

  useEffect(() => {
    // Check if user is authenticated
    const currentUser = authService.getCurrentUser()
    if (!currentUser) {
      router.push('/login')
      return
    }

    setUser(currentUser)

    // Redirect based on role
    switch (currentUser.role) {
      case 'admin':
        router.push('/admin/dashboard')
        break
      case 'manager':
        router.push('/manager/dashboard')
        break
      case 'engineer':
        router.push('/engineer/dashboard')
        break
      case 'sales':
        router.push('/sales/dashboard')
        break
      default:
        router.push('/login')
    }
  }, [router, setUser])

  return (
    <div className="flex items-center justify-center min-h-screen">
      <div className="text-center">
        <h1 className="text-2xl font-semibold">正在載入...</h1>
      </div>
    </div>
  )
}