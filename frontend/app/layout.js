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

  const showNavbar = pathname !== "/login" && pathname !== "/signup";
  const showGroupbar = pathname !== "/login" && pathname !== "/signup";
  const showChatbar = pathname !== "/login" && pathname !== "/signup";

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
              <button
                onClick={handleLogout}
                className="text-red-600 hover:underline"
              >
                Logout
              </button>
              <a href="/profile" className="hover:underline">
                <img
                  src="/avatar.png" // Replace with the actual path to the profile picture
                  alt="Profile"
                  className="w-10 h-10 rounded-full object-cover"
                />
              </a>
            </div>
          </nav>
        )}
    
        <div className="flex min-h-screen">
          {/* Left Sidebar */}
          {showGroupbar && (
          <div className="w-1/6 bg-gray-200 p-4">
            <h2 className="text-lg font-bold mb-4">Groups</h2>
            <ul className="space-y-2">
              <li><a href="/group1" className="text-blue-600 hover:underline">Group 1</a></li>
              <li><a href="/group2" className="text-blue-600 hover:underline">Group 2</a></li>
              <li><a href="/group3" className="text-blue-600 hover:underline">Group 3</a></li>
              {/* Add more groups as needed */}
            </ul>
          </div>
          )}
          {/* Main Content */}
          <div className="flex-1 p-4">
            {children}
          </div>

          {/* Right Sidebar */}
          {showChatbar && (
          <div className="w-1/6 bg-gray-200 p-4">
            <h2 className="text-lg font-bold mb-4">Chats</h2>
            <ul className="space-y-2">
              <li><a href="/chat/user1" className="text-blue-600 hover:underline">User 1</a></li>
              <li><a href="/chat/user2" className="text-blue-600 hover:underline">User 2</a></li>
              <li><a href="/chat/user3" className="text-blue-600 hover:underline">User 3</a></li>
              {/* Add more users as needed */}
            </ul>
          </div>
          )}
        </div>
      </body>
    </html>
  );
}

function handleLogout() {
  fetch('http://localhost:8080/api/logout', {
    method: 'POST',
    credentials: 'include', // Include cookies for authentication
    headers: {
      'Content-Type': 'application/json', // Ensure the correct Content-Type
    },
  })
    .then((res) => {
      if (res.ok) {
        window.location.href = '/login';
      } else {
        console.error('Logout failed with status:', res.status);
        alert('Logout failed');
      }
    })
    .catch((err) => {
      console.error('Error logging out:', err);
      alert('An error occurred while logging out');
    });
}