'use client';

import { useEffect } from 'react';
import { usePathname } from 'next/navigation';
import { useUser } from '../context/UserContext';
import Link from 'next/link';
import NotificationsDropdown from './Notification';
import HomeIcon from '../../public/home.png';
import Image from 'next/image';
import initWebSocket from './ws';


export default function Header() {
 
  const pathname = usePathname();

  const { user } = useUser();

  // Determine if the navbar should be shown
  const showNavbar = pathname !== '/login' && pathname !== '/signup';

  // Render nothing if the navbar shouldn't be shown
  if (!showNavbar) return null;

  useEffect(() => {
  initWebSocket((data) => {
    console.log("🟡 New WS message:", data);
    // Handle the message
  });
}, []);


  return (
      <header className="flex items-center justify-between p-4 bg-gray-100">
      {/* Left Section */}
      <div className="flex items-center space-x-4 p-2">
          <Link href="/" className="flex items-center">
              <Image
              src={HomeIcon}
              alt="Home"
              width={30}
              height={30}
              className="mr-2"
              />
          </Link>
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
          <NotificationsDropdown />
          <button
          onClick={handleLogout}
          className="text-red-600 hover:underline"
          >
          Logout
          </button>
            <Link href={`/profile?user_id=${user?.user_id}`} >
              <img
              src={user?.avatar_path ? `http://localhost:8080${user.avatar_path}` : '/avatar.png'}
              alt="Profile"
              className="w-10 h-10 rounded-full object-cover"/>
            </Link>
      </div>
      </header>
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
        localStorage.removeItem('user'); // Clear storage
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