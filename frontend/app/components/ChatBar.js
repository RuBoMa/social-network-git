'use client';
import { usePathname } from 'next/navigation';
import { useEffect, useState } from 'react';
import ChatWindow from './ChatWindow';
import Author from './Author';

export default function ChatBar() {
  const pathname = usePathname();
  const showChatbar = pathname !== '/login' && pathname !== '/signup';
  const [users, setUsers] = useState([]);
  const [openUser, setOpenUser] = useState(null); // keep track of current open chatwindow
  const [currentUserId, setCurrentUserId] = useState(null);

    useEffect(() => {
        async function fetchUsers() {
            try {
                const res = await fetch('http://localhost:8080/api/users');
                if (!res.ok) throw new Error('Failed to fetch users');
                const data = await res.json();
                const userList = Array.isArray(data) ? data : data.users || [];
                setUsers(userList);
            } catch (err) {
                console.error('Error fetching users:', err);
            }
        }

        fetchUsers();
    }, []);

    console.log('Fetched users:', users);

  if (!showChatbar) return null;

//   useEffect(() => {
//   async function fetchCurrentUser() {
//     const res = await fetch('http://localhost:8080/api/me', { credentials: 'include' });
//     if (res.ok) {
//       const data = await res.json();
//       setCurrentUserId(data.user_id);
//     }
//   }
//   fetchCurrentUser();
// }, []);

  const filteredUsers = users.filter(u => u.user_id !== currentUserId);

  return (
    <>
    {openUser && <ChatWindow chatPartner={openUser} onClose={() => setOpenUser(null)} />}
      <div className="w-1/6 bg-gray-200 p-4 overflow-y-auto">
        <h2 className="text-lg font-bold mb-4">Chats</h2>
        <ul className="space-y-2">
          {filteredUsers.map(user => (
            <li key={user.user_id}>
              <button
                onClick={() => setOpenUser(user)}
                className="flex items-center space-x-2 w-full text-left"
              >
                <Author author={user} disableLink={true} size="sm" />
              </button>
            </li>
          ))}
        </ul>
      </div>
    </>
  );
}
