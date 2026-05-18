import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'

const PUBLIC_PATHS = ['/login', '/register', '/forgot-password', '/reset-password', '/verify-email']

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl

  // Always allow public auth pages — prevents redirect loop
  if (PUBLIC_PATHS.some((p) => pathname === p || pathname.startsWith(p + '/'))) {
    return NextResponse.next()
  }

  // Check token cookie (set by lib/api.ts on login)
  const token = request.cookies.get('meetext_token')?.value

  if (!token) {
    const loginUrl = new URL('/login', request.url)
    return NextResponse.redirect(loginUrl)
  }

  return NextResponse.next()
}

export const config = {
  // Only run on page routes — skip _next, static files, api proxy, favicon
  matcher: [
    '/((?!_next/static|_next/image|favicon.ico|api/).*)',
  ],
}
