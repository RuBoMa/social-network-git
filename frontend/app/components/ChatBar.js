'use client';
import { usePathname } from 'next/navigation';
import { useEffect, useState } from 'react';
import ChatWindow from './ChatWindow';
import Author from './Author';
import { addMessageHandler } from './ws';

export default function ChatBar() {
  const pathname = usePathname();
  const showChatbar = pathname !== '/login' && pathname !== '/signup';
  const [users, setUsers] = useState([]);
  const [openUser, setOpenUser] = useState(null); // track currently opened chat
  const [currentUserId, setCurrentUserId] = useState(null);

  // Get user_id from localStorage instead of fetching it
  useEffect(() => {
    const storedId = localStorage.getItem("user_id");

    if (storedId) {
      setCurrentUserId(parseInt(storedId));
    } else {
      console.warn("No user_id found in localStorage");
    }
  }, []);

  // Fetch sorted users once currentUserId is set
  useEffect(() => {
    if (!currentUserId) return;

    async function fetchSortedUsers() {
      try {
        const res = await fetch('http://localhost:8080/api/users/sorted', {
          credentials: 'include',
        });

        if (!res.ok) {
          console.error('Failed to fetch sorted users:', res.status);
          return;
        }

        const data = await res.json();
        const sortedUsers = Array.isArray(data) ? data : data.users || [];
        setUsers(sortedUsers);
      } catch (err) {
        console.error('Error fetching users:', err);
      }
    }

    fetchSortedUsers();
  }, [currentUserId]);

  // WebSocket updates to user list
  useEffect(() => {
    if (!currentUserId) return;

    const removeHandler = addMessageHandler((data) => {
      console.log("ChatBar received message:", data);
      if (data.type === 'sorted_users') {
        console.log('Updating chat list with sorted users:', data.users);
        setUsers(data.users || []);
      }
    });

    return () => {
      if (removeHandler) removeHandler();
    };
  }, [currentUserId]);

  if (!showChatbar) return null;

  return (
    <>
      {openUser && <ChatWindow chatPartner={openUser} onClose={() => setOpenUser(null)} />}
      <div className="w-1/6 bg-gray-200 p-4 overflow-y-auto">
        <h2 className="text-lg font-bold mb-4">Chats</h2>
        {users.length === 0 ? (
          <p className="text-gray-500 text-sm">No conversations yet</p>
        ) : (
          <ul className="space-y-2">
            {users.map(user => (
              <li key={user.user_id}>
                <button
                  onClick={() => setOpenUser(user)}
                  className="flex items-center space-x-2 w-full text-left hover:bg-gray-300 p-2 rounded"
                >
                  <Author author={user} disableLink={true} size="sm" />
                </button>
              </li>
            ))}
          </ul>
        )}
      </div>
    </>
  );
}
