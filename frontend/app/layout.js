'use client';

import { Geist, Geist_Mono } from 'next/font/google';
import './globals.css';
import Header from './components/Header';
import GroupBar from './components/GroupBar';
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
  const showGroupbar = true; // Replace with your logic to show/hide the group bar
  const showChatbar = true; // Replace with your logic to show/hide the chat bar

  return (
    <html lang="en">
      <body className={`${geistSans.variable} ${geistMono.variable} antialiased`}>
        <UserProvider>
        <Header />
        <div className="flex min-h-screen">
          <GroupBar showGroupbar={showGroupbar} />
          <div className="flex-1 p-4">{children}</div>
          <ChatBar showChatbar={showChatbar} />
        </div>
        </UserProvider>
      </body>
    </html>
  );
}