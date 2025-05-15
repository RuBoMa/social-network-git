'use client';
import { usePathname } from 'next/navigation';
import { useEffect, useState } from 'react';
import Link from 'next/link';

export default function GroupBar() {
  const pathname = usePathname();
  const showGroupbar = pathname !== '/login' && pathname !== '/signup';
  const [groups, setGroups] = useState([]);

  useEffect(() => {
    async function fetchGroups() {
      try {
        const res = await fetch('http://localhost:8080/api/all-groups', {
          credentials: 'include',
        });

        if (res.ok) {
          const data = await res.json();
          setGroups(data);
        } else {
          console.error('Failed to fetch groups');
        }
      } catch (err) {
        console.error('Error fetching groups:', err);
      }
    }

    fetchGroups();
  }, [groups]);

  if (!showGroupbar) return null;

  return (
    <div className="w-1/6 bg-gray-200 p-4">
      <h2 className="text-lg font-bold mb-4">Groups</h2>
      <ul className="space-y-2">
        {groups?.length > 0 ? (
          groups.map((group) => (
            <li key={group.group_id}>
              <Link
                href={`/group?group_id=${group.group_id}`}
                className="text-blue-600 hover:underline"
                title={group.group_desc}
              >
                {group.group_title}
              </Link>
            </li>
          ))
        ) : (
          <li className="text-gray-500 text-sm">No groups available.</li>
        )}
      </ul>
    </div>
  );
}
