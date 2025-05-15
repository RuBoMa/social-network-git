'use client';
import { usePathname } from 'next/navigation';

export default function ChatBar() {
    const pathname = usePathname();
    const showChatbar = pathname !== '/login' && pathname !== '/signup';

    if (!showChatbar) return null; // Don't render if `showChatbar` is false

    return (
        <div className="w-1/6 bg-gray-200 p-4">
        <h2 className="text-lg font-bold mb-4">Chats</h2>
        <ul className="space-y-2">
            <li>
            <a href="/chat/user1" className="text-blue-600 hover:underline">
                User 1
            </a>
            </li>
            <li>
            <a href="/chat/user2" className="text-blue-600 hover:underline">
                User 2
            </a>
            </li>
            <li>
            <a href="/chat/user3" className="text-blue-600 hover:underline">
                User 3
            </a>
            </li>
            {/* Add more users as needed */}
        </ul>
        </div>
    );
    }