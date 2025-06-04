'use client';

import { useEffect } from 'react';
import { useRouter, usePathname, useSearchParams } from 'next/navigation';
import { Geist, Geist_Mono } from 'next/font/google';
import './globals.css';
import Header from './components/Header';
import GroupBar from './components/Group/GroupBar';
import ChatBar from './components/ChatBar';
import { UserProvider } from './context/UserContext';

const geistSans = Geist({
  variable: '--font-geist-sans',
  subsets: ['latin'],
});

const geistMono = Geist_Mono({
  variable: '--font-geist-mono',
  subsets: ['latin'],
});

export default function RootLayout({ children }) {
  const router = useRouter();
  useEffect(() => {
    // Check if the user is authenticated
    const token = localStorage.getItem('token');
    if (!token) {
      // Redirect to login if no token is found
      router.replace('/login');
    }
  }, [router]);
  useEffect(() => {
    // catch 401 Unauthorized responses globally for all fetches
    const _fetch = window.fetch;
    window.fetch = async (...args) => {
      const res = await _fetch(...args);
      if (res.status === 401) {
        // clear client state
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        // redirect to login
        router.replace('/login');
      }
      return res;
    };
    return () => {
      window.fetch = _fetch;
    };
  }, [router]);


  useEffect(() => {
    const prevSearch = window.location.search;

    const interval = setInterval(() => {
      if (window.location.search !== prevSearch) {
        window.location.reload(); // Full page reload
      }
    }, 100); // Check every 100ms (can adjust)

    return () => clearInterval(interval);
  }, []);

  return (
    <html lang="en">
      <body className={`${geistSans.variable} ${geistMono.variable} antialiased`}>
        <UserProvider>
          <Header />
          <div className="flex min-h-screen">
            <GroupBar />
            <main className="flex-1 overflow-y-auto p-4">
              <div className="max-w-7xl mx-auto">{children}</div>
            </main>
            <ChatBar />
          </div>
        </UserProvider>
      </body>
    </html>
  );
}