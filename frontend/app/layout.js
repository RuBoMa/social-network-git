'use client'
import { Geist, Geist_Mono } from "next/font/google";
import { usePathname } from "next/navigation";
import "./globals.css";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export default function RootLayout({ children }) {
  const pathname = usePathname();

  // Show the navbar on all pages during development
  const isDevelopment = process.env.NODE_ENV === "development";
  const showNavbar = isDevelopment || (pathname !== "/login" && pathname !== "/register");

  return (
    <html lang="en">
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
      >
        {showNavbar && (
          <nav className="flex items-center justify-between p-4 bg-gray-100">
            {/* Left Section */}
            <div className="flex items-center space-x-4">
              <a href="/" className="text-blue-600 hover:underline">Home</a>
            </div>

            {/* Center Section */}
            <div className="flex-grow flex justify-center">
              <input
                type="text"
                placeholder="Search..."
                className="w-1/2 p-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>

            {/* Right Section */}
            <div className="flex items-center space-x-4">
              <a href="/notifications" className="text-blue-600 hover:underline">Notifications</a>
              <a href="/profile" className="hover:underline">
                <img
                  src="/path-to-profile-picture.jpg" // Replace with the actual path to the profile picture
                  alt="Profile"
                  className="w-10 h-10 rounded-full object-cover"
                />
              </a>
            </div>
          </nav>
        )}
        <div className="flex flex-row min-h-screen">
          {/* Left Sidebar */}
          <aside className="w-1/6 bg-gray-200 p-4">
            <h2 className="text-lg font-bold mb-4">Groups</h2>
            <ul className="space-y-2">
              <li><a href="/group1" className="text-blue-600 hover:underline">Group 1</a></li>
              <li><a href="/group2" className="text-blue-600 hover:underline">Group 2</a></li>
              <li><a href="/group3" className="text-blue-600 hover:underline">Group 3</a></li>
              {/* Add more groups as needed */}
            </ul>
          </aside>

          {/* Main Content */}
          <main className="flex-1 p-4">
            {children}
          </main>

          {/* Right Sidebar */}
          <aside className="w-1/6 bg-gray-200 p-4">
            <h2 className="text-lg font-bold mb-4">Chats</h2>
            <ul className="space-y-2">
              <li><a href="/chat/user1" className="text-blue-600 hover:underline">User 1</a></li>
              <li><a href="/chat/user2" className="text-blue-600 hover:underline">User 2</a></li>
              <li><a href="/chat/user3" className="text-blue-600 hover:underline">User 3</a></li>
              {/* Add more users as needed */}
            </ul>
          </aside>
        </div>
      </body>
    </html>
  );
}