'use client';

import { useEffect, useState } from 'react';
import { usePathname } from 'next/navigation';

export default function Header() {
    const pathname = usePathname();
    const [user, setUser] = useState(null);

    // Determine if the navbar should be shown
    const showNavbar = pathname !== '/login' && pathname !== '/signup';

    // Fetch user data only if the navbar should be shown
    useEffect(() => {
        if (!showNavbar) return; // Skip fetching user data if the navbar is hidden

        const fetchUser = async () => {
        try {
            const res = await fetch('http://localhost:8080/api/profile', {
            method: 'GET',
            credentials: 'include', // Include cookies for authentication
            headers: {
                'Content-Type': 'application/json',
            },
            });
            if (res.ok) {
            const data = await res.json();
            setUser(data); // Save user data
            console.log('User data:', data);
            } else if (res.status === 401) {
            window.location.href = '/login'; // Redirect to login if unauthorized
            } else {
            console.error('Failed to fetch user data');
            }
        } catch (err) {
            console.error('Error fetching user data:', err);
        }
        };

        fetchUser();
    }, [showNavbar]);

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
                src={user?.user.avatar_path ? `http://localhost:8080${user.user.avatar_path}` : '/avatar.png'}
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