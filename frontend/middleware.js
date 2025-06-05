import { NextResponse } from 'next/server'

export async function middleware(req) {
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

}

export const config = {
  matcher: ['/((?!_next/static|favicon.ico).*)']
}