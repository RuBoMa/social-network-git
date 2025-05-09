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
          <nav className="flex flex-row space-x-4 p-4 bg-gray-100">
            <a href="/" className="text-blue-600 hover:underline">Home</a>
            <a href="/notifications" className="text-blue-600 hover:underline">Notifications</a>
            <a href="/profile" className="hover:underline">
              <img
                src="/path-to-profile-picture.jpg" // Replace with the actual path to the profile picture
                alt="Profile"
                className="w-10 h-10 rounded-full object-cover"
              />
            </a>
          </nav>
        )}
        {children}
      </body>
    </html>
  );
}