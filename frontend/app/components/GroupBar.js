'use client';
import { usePathname } from 'next/navigation';

export default function GroupBar() {
    const pathname = usePathname();
    const showGroupbar = pathname !== '/login' && pathname !== '/signup';

    if (!showGroupbar) return null; // Don't render if `showGroupbar` is false

    return (
    <div className="w-1/6 bg-gray-200 p-4">
        <h2 className="text-lg font-bold mb-4">Groups</h2>
        <ul className="space-y-2">
        <li>
            <a href="/group1" className="text-blue-600 hover:underline">
            Group 1
            </a>
        </li>
        <li>
            <a href="/group2" className="text-blue-600 hover:underline">
            Group 2
            </a>
        </li>
        <li>
            <a href="/group3" className="text-blue-600 hover:underline">
            Group 3
            </a>
        </li>
        {/* Add more groups as needed */}
        </ul>
    </div>
    );
}