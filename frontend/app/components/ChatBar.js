'use client';
import { usePathname } from 'next/navigation';
import { useEffect, useState } from 'react';

export default function ChatBar() {
    const pathname = usePathname();
    const showChatbar = pathname !== '/login' && pathname !== '/signup';
    const [users, setUsers] = useState([]);

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
        <div className="w-1/6 bg-gray-200 p-4 overflow-y-auto max-h-screen">
            <h2 className="text-lg font-bold mb-4">Chats</h2>
            <ul className="space-y-2">
                {users.map((user) => (
                    <li key={user.id} className="flex items-center space-x-2">
                        <img
                            src={user.avatar_path ? `http://localhost:8080${user.avatar_path}` : '/avatar.png'}
                            alt={user.username}
                            className="w-8 h-8 rounded-full"
                        />
                        <a
                            href={`/chat/${user.id}`}
                            className="text-blue-600 hover:underline"
                        >
                            {user.first_name} {user.last_name}
                        </a>
                    </li>
                ))}
            </ul>
        </div>
    );
}
