'use client';

import { useEffect } from 'react';
import { usePathname } from 'next/navigation';
import { useUser } from '../context/UserContext';
import Link from 'next/link'
import Image from 'next/image';
import NotificationsDropdown from './Notification';
import HomeIcon from '../../public/home.png';
import Logout from '../../public/logout.png';
import initWebSocket, { addMessageHandler, closeWebSocket } from './ws';
import SearchBar from './Searchbar';

export default function Header() {

  const pathname = usePathname();
  const { user, setUser } = useUser();
    // Determine if the navbar should be shown
  const showNavbar = pathname !== '/login' && pathname !== '/signup';

  useEffect(() => {
    if (!localStorage.getItem('token')) {
      return;
    }
initWebSocket(); // Initialize WebSocket connection
    const removeHandler = addMessageHandler((data) => {
    });

    return () => {
        removeHandler?.();
    };
}, []);


  // This needs to be after UseEffect to ensure user is set before checking (Hook issues)
  if (!showNavbar) return null;

function handleLogout() {
  if (!localStorage.getItem('token')) {
    return;
  }
  fetch('http://localhost:8080/api/logout', {
    method: 'POST',
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
    },
  })
    .then((res) => {
      if (res.ok) {
        localStorage.removeItem('user'); // Clear storage
        localStorage.removeItem('token'); // Clear token
        closeWebSocket();
        setUser(null); // Clear user context
        window.location.href = '/login';
      } else {
        console.log('Logout failed with status:', res.status);
        alert('Logout failed');
      }
    })
    .catch((err) => {
      console.error('Error logging out:', err);
      alert('An error occurred while logging out');
    });
}

  return (
  <header className="flex items-center justify-between p-4 bg-sky-600/40 shadow-lg">
      {/* Left Section */}
      <div className="flex items-center space-x-4 p-2">
        <Link href="/" className="flex items-center">
          <Image
            src={HomeIcon}
            alt="Home"
            width={30}
            height={30}
            className="ml-2 cursor-pointer"
          />
        </Link>
      </div>

      {/* Center Section */}
      <div className="flex-grow flex justify-center">
        <SearchBar />
      </div>

      {/* Right Section */}
      <div className="flex items-center space-x-4">
        <NotificationsDropdown />
        <button
          onClick={handleLogout}
        >
          <Image
            src={Logout}
            alt="Logout"
            width={24}
            height={24}
            className="ml-2 cursor-pointer mr-1"
          />

        </button>
        <Link href={`/profile?user_id=${user?.user_id}`} >
          <img
            src={user?.avatar_path ? `http://localhost:8080${user.avatar_path}` : '/avatar.png'}
            alt="Profile"
            className="w-8 h-8 rounded-full object-cover" />
        </Link>
      </div>
    </header>
  );
}
