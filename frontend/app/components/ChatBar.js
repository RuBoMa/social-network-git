'use client';
import { usePathname } from 'next/navigation';
import { useEffect, useState } from 'react';
import ChatWindow from './ChatWindow';
import Author from './Author';
import GroupAvatar from './GroupAvatar';
import { addMessageHandler } from './ws';

export default function ChatBar() {
  const pathname = usePathname();
  const showChatbar = pathname !== '/login' && pathname !== '/signup';
  const [users, setUsers] = useState([]);
  const [openUser, setOpenUser] = useState(null); // track currently opened chat
  const [currentUserId, setCurrentUserId] = useState(null);
  const [groups, setGroups] = useState([]);
  const [openGroup, setOpenGroup] = useState(null);


    useEffect(() => {
      async function fetchGroups() {
        try {
            const res = await fetch('http://localhost:8080/api/my-groups', { credentials: 'include' });
            if (!res.ok) throw new Error('Failed to fetch groups');
            const data = await res.json();
            setGroups(Array.isArray(data) ? data : data.groups || []);
          } catch (err) {
            console.error('Error fetching groups:', err);
        }
      }
      fetchGroups();
    }, []);
  const [unreadChats, setUnreadChats] = useState({});

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

  // Listen for custom event to open chat from ProfilePage
  useEffect(() => {
  function handleOpenChat(e) {
    const userId = e.detail.user_id;
    let user = users.find(u => u.user_id === userId);

    if (!user) {
      // Fetch user info if not in chat list
      fetch(`http://localhost:8080/api/profile?user_id=${userId}`, {
        credentials: 'include'
      })
        .then(res => res.json())
        .then(data => {
          const userData = data.user || data;
          setOpenUser({ ...userData, user_id: userId });
        })
        .catch(() => setOpenUser({ user_id: userId }));
    } else {
      setOpenUser(user);
    }
  }
    window.addEventListener('open-chat', handleOpenChat);
    return () => window.removeEventListener('open-chat', handleOpenChat);
  }, [users]);

  if (!showChatbar) return null;

  const filteredUsers = users.filter(u => u.user_id !== currentUserId);

  useEffect(() => {
    const storedUser = localStorage.getItem('user');
    if (storedUser) {
      try {
        const parsedUser = JSON.parse(storedUser);
        setCurrentUserId(parsedUser.user_id);
      } catch (error) {
        console.error('Error parsing user from localStorage:', error);
      }
    } else {
      console.error('No currentUser found in localStorage');
    }
  }, []);

  useEffect(() => {
    const handler = (data) => {
      // Ensure the notification is only for the receiver
      if (data.type === 'message' && data.receiver?.user_id) {
        setUnreadChats((prev) => {
          const updated = {
            ...prev,
            [data.sender.user_id]: (prev[data.sender.user_id] || 0) + 1,
          };
          return updated;
        });
      }
    };

    const removeHandler = addMessageHandler(handler);
    return () => removeHandler();
  }, [currentUserId]);


  return (
    <>
      {openUser && <ChatWindow chatPartner={openUser} onClose={() => setOpenUser(null)} />}
    {openGroup && <ChatWindow group={openGroup} onClose={() => setOpenGroup(null)} />}
      <div className="w-1/6 bg-gray-200 p-4 overflow-y-auto">
        <h2 className="text-lg font-bold mb-4">Chats</h2>
        <ul className="space-y-2">
          {filteredUsers.map(user => (
            <li key={user.user_id}>
              <button
                onClick={() => {
                  setOpenUser(user);
                  setUnreadChats(prev => ({
                    ...prev,
                    [user.user_id]: 0
                  }));
                }}
                className="flex items-center space-x-2 w-full text-left"

              >
                <Author author={user} disableLink={true} size="sm" />

                {/* 🔴 Unread dot */}
                {unreadChats[user.user_id] > 0 && (
                  <span className="ml-auto h-2 w-2 bg-red-600 rounded-full"></span>

                )}
              </button>
            </li>
          ))}
        </ul>
        <ul className="space-y-2">
          {groups.map(group => (
            <li key={group.group_id}>
              <button
                onClick={() => {
                  setOpenUser(null);
                  setOpenGroup(group);
                }}
                className="flex items-center space-x-2 w-full text-left"
              >
                <GroupAvatar group={group} disableLink={true} size="sm" />
              </button>
            </li>
          ))}
        </ul>
      </div>
    </>
  );
}