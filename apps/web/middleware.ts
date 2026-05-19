import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'

// Paths that unauthenticated users can access
const PUBLIC_PATHS = [
  '/login',
  '/register',
  '/forgot-password',
  '/reset-password',
  '/verify-email',
  '/auth/callback'
]

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl

  const hasSession =
    Boolean(request.cookies.get('meetext_access')?.value) ||
    Boolean(request.cookies.get('meetext_refresh')?.value)

  const isPublicPath = PUBLIC_PATHS.some((p) => pathname === p || pathname.startsWith(p + '/'))

  // If user is authenticated and trying to access an auth page, redirect to dashboard
  if (hasSession && isPublicPath) {
    return NextResponse.redirect(new URL('/dashboard', request.url))
  }

  // If user is NOT authenticated and trying to access a protected page, redirect to login
  if (!hasSession && !isPublicPath) {
    const loginUrl = new URL('/login', request.url)
    // Avoid setting 'next' to '/' to keep URL clean, or let it be.
    if (pathname !== '/') {
      loginUrl.searchParams.set('next', pathname)
    }
    return NextResponse.redirect(loginUrl)
  }

  // Handle root route redirection specifically
  if (pathname === '/') {
    return NextResponse.redirect(new URL('/dashboard', request.url))
  }

  return NextResponse.next()
}

export const config = {
  // Only run on page routes — skip _next, static files, api proxy, favicon
  matcher: [
    '/((?!_next/static|_next/image|favicon.ico|api/).*)',
  ],
}
