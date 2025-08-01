import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'

// Public routes that don't require authentication
const publicRoutes = ['/', '/login', '/about']

// Role-based route protection
const roleRoutes: Record<string, string[]> = {
  admin: ['/admin'],
  manager: ['/manager', '/admin'],
  engineer: ['/engineer'],
  sales: ['/sales'],
}

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl
  
  // Check if it's a public route
  if (publicRoutes.includes(pathname)) {
    return NextResponse.next()
  }
  
  // Get token from cookie (Next.js will automatically handle this in production)
  const token = request.cookies.get('access_token')?.value
  
  // If no token, redirect to login
  if (!token) {
    const loginUrl = new URL('/login', request.url)
    loginUrl.searchParams.set('from', pathname)
    return NextResponse.redirect(loginUrl)
  }
  
  // TODO: In production, decode JWT and check role permissions
  // For now, we'll allow authenticated users to proceed
  
  return NextResponse.next()
}

export const config = {
  matcher: [
    /*
     * Match all request paths except for the ones starting with:
     * - api (API routes)
     * - _next/static (static files)
     * - _next/image (image optimization files)
     * - favicon.ico (favicon file)
     */
    '/((?!api|_next/static|_next/image|favicon.ico).*)',
  ],
}