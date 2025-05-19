'use client';
import { usePathname } from 'next/navigation';
import { useEffect, useState } from 'react';
import ChatWindow from './ChatWindow';

export default function ChatBar() {
  const pathname = usePathname();
  const showChatbar = pathname !== '/login' && pathname !== '/signup';
  const [users, setUsers] = useState([]);
  const [openUser, setOpenUser] = useState(null); // keep track of current open chatwindow

    useEffect(() => {
        async function fetchUsers() {
            try {
                const res = await fetch('http://localhost:8080/api/users');
                if (!res.ok) throw new Error('Failed to fetch users');
                const data = await res.json();
                setUsers(data);
            } catch (err) {
                console.error('Error fetching users:', err);
            }
        }

        fetchUsers();
    }, []);

  if (!showChatbar) return null;

  return (
    <>
      {openUser && <ChatWindow user={openUser} onClose={() => setOpenUser(null)} />}
      <div className="w-1/6 bg-gray-200 p-4 overflow-y-auto max-h-screen">
        <h2 className="text-lg font-bold mb-4">Chats</h2>
        <ul className="space-y-2">
          {users.map(user => (
            <li key={user.user_id}>
              <button
                onClick={() => setOpenUser(user)}
                className="flex items-center space-x-2 w-full text-left"
              >
                <img
                  src={user.avatar_path ? `http://localhost:8080${user.avatar_path}` : '/avatar.png'}
                  alt={user.username}
                  className="w-8 h-8 rounded-full"
                />
                <span className="text-blue-600 hover:underline">
                  {user.first_name} {user.last_name}
                </span>
              </button>
            </li>
          ))}
        </ul>
      </div>
    </>
  );
}

