'use client';

import { usePathname } from 'next/navigation';
import { useUser } from '../context/UserContext';

export default function Header() {
  const pathname = usePathname();
  const { user } = useUser();

  // Determine if the navbar should be shown
  const showNavbar = pathname !== '/login' && pathname !== '/signup';

  // Render nothing if the navbar shouldn't be shown
  if (!showNavbar) return null;

  return (
      <header className="flex items-center justify-between p-4 bg-gray-100">
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
              src={user?.avatar_path ? `http://localhost:8080${user.avatar_path}` : '/avatar.png'}
              alt="Profile"
              className="w-10 h-10 rounded-full object-cover"
          />
          </a>
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