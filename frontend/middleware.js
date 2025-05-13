import { NextResponse } from 'next/server'

export function middleware(req) {
  const { pathname } = req.nextUrl
  // allow next internals and the login page
  if (
    pathname.startsWith('/_next') ||
    pathname.startsWith('/static') ||
    pathname === '/login' ||
    pathname === '/signup'
  ) {
    return NextResponse.next()
  }

  // if no session cookie → redirect to /login
  const hasSession = !!req.cookies.get('session_id')
  if (!hasSession) {
    const url = req.nextUrl.clone()
    url.pathname = '/login'
    return NextResponse.redirect(url)
  }

  return NextResponse.next()
}

export const config = {
  matcher: ['/((?!_next/static|favicon.ico).*)']
}