'use client';

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
  )
}