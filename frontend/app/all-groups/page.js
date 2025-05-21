'use client';
import { useEffect, useState } from 'react';
import Link from 'next/link';
import ErrorMessage from '../components/ErrorMessage';
import CreateGroup from '../components/CreateGroup';

export default function GroupsPage() {
const [groups, setGroups] = useState([]);
const [sortType, setSortType] = useState('alphabetical'); // or 'created'
const [showModal, setShowModal] = useState(false);


  useEffect(() => {
    async function fetchGroups() {
      try {
        const res = await fetch('http://localhost:8080/api/all-groups', {
          credentials: 'include',
          method: 'GET',
        });

        const data = await res.json();
        if (res.ok) {
          setGroups(data);
          console.log('Fetched groups:', data); // Log the fetched groups
        } 
      } catch (err) {
        console.error('Error fetching groups:', err);
        ErrorMessage(data.message || 'Failed to load groups');
      }
    }

    fetchGroups();
  }, []);

  if (!groups) {
    return <div>Loading groups...</div>
  };

  return (
     <div className="p-4 border rounded mb-6">
        <div className="flex items-center justify-between mb-4">
            <h2 className="text-lg font-bold">All Groups</h2>
            <button
            onClick={() => setShowModal(true)}
            className="px-3 py-1 bg-green-600 text-white rounded"
            >
            + Create Group
            </button>
        </div>

        {showModal && <CreateGroup onClose={() => setShowModal(false)} />}
        <select
            value={sortType}
            onChange={(e) => setSortType(e.target.value)}
            disabled={showModal}   // disable while modal is open
            className={`mb-4 p-1 rounded border ${showModal ? 'cursor-not-allowed opacity-50' : ''}`}
            >
            <option value="alphabetical">Sort: A-Z</option>
            <option value="created">Sort: Newest</option>
        </select>

      <ul className="space-y-2">
        {groups?.length > 0 ? (
            [...groups]
            .sort((a, b) =>
              sortType === 'alphabetical'
                ? a.group_name.localeCompare(b.group_name)
                : new Date(b.group_created_at) - new Date(a.group_created_at)
            )
            .map((group) => (
            <li key={group.group_id}>
            <Link href={`/group?group_id=${group.group_id}`} className="text-blue-600 hover:underline">
                {group.group_name}
            </Link>
            <p className="text-sm text-gray-700">{group.group_desc.slice(0, 100)}...</p>
            <p className="text-xs text-gray-500">
                Created by {group.group_creator.nickname || `${group.group_creator.first_name} ${group.group_creator.last_name}` } on {new Date(group.group_created_at).toLocaleDateString()}
            </p>
            </li>
          ))
        ) : (
          <li className="text-gray-500 text-sm">No groups available.</li>
        )}
      </ul>
    </div>
  );
}
