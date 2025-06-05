'use client';
import { useEffect, useState } from 'react';
import Link from 'next/link';
import ErrorMessage from '../components/ErrorMessage';
import CreateGroup from '../components/Group/CreateGroup';

export default function GroupsPage() {
  const [groups, setGroups] = useState([]);
  const [sortType, setSortType] = useState('alphabetical');

  useEffect(() => {
    async function fetchGroups() {
      try {
        const res = await fetch('http://localhost:8080/api/all-groups', {
          credentials: 'include',
          method: 'GET',
        });
        const data = await res.json();
        if (res.ok) {
          setGroups(Array.isArray(data) ? data : []);
        } else { 
          ErrorMessage(data.message || 'Failed to load groups');
        }
      } catch (err) {
        console.error('Error fetching groups:', err);
        ErrorMessage('Failed to load groups');
      }
    }
    fetchGroups();
  }, []);
  
  return (
    <div className="p-4">
    <div className="mb-6">
      <CreateGroup />
    </div>
    <div className="p-4 border border-gray-200 rounded shadow space-y-4">
      <div>
        <h2 className="text-lg font-bold mb-2">All Groups</h2>
        <select
          value={sortType}
          onChange={e => setSortType(e.target.value)}
          className="mb-4 p-1 rounded border"
        >
          <option value="alphabetical">Sort: A-Z</option>
          <option value="created">Sort: Newest</option>
        </select>

        <ul className="space-y-2">
          {groups.length > 0 ? (
            [...groups]
              .sort((a, b) =>
                sortType === 'alphabetical'
                  ? a.group_name.localeCompare(b.group_name)
                  : new Date(b.group_created_at) - new Date(a.group_created_at)
              )
              .map(group => (
                <li key={group.group_id}>
                  <Link
                    href={`/group?group_id=${group.group_id}`}
                    className="text-sky-600/80 hover:text-sky-800 transition cursor-pointer"
                  >
                    {group.group_name}
                  </Link>
                  <p className="text-sm text-gray-700 mt-1 mb-1">
                    {group.group_desc.slice(0, 100)}…
                  </p>
                  <p className="text-xs text-gray-500">
                    Created by{' '}
                    {group.group_creator.nickname ||
                      `${group.group_creator.first_name} ${group.group_creator.last_name}`}{' '}
                    on {new Date(group.group_created_at).toLocaleDateString()}
                  </p>
                </li>
              ))
          ) : (
            <li className="text-gray-500 text-sm">No groups available.</li>
          )}
        </ul>
      </div>
    </div>
    </div>
  );
}
