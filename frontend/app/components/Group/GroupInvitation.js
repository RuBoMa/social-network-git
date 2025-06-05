'use client'
import { useEffect, useState } from 'react'
import Author from '../Author'

export default function GroupInvitation({ groupId }) {
  const [users, setUsers] = useState([])
  const [open, setOpen] = useState(false)

  useEffect(() => {
    async function fetchAvailableUsers() {
      const res = await fetch(`http://localhost:8080/api/group/invite?group_id=${groupId}`, {
        credentials: 'include',
        method: 'GET',
      })
      const data = await res.json()
      if (res.ok) setUsers(data)
      else {
        console.error('Failed to fetch available users:', data.message)
        setUsers([])
      }
    }
    if (groupId) fetchAvailableUsers()
  }, [groupId])

  async function inviteUser(userId) {
    const res = await fetch('http://localhost:8080/api/request', {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        group: { group_id: groupId },
        joining_user: { user_id: userId },
        status: 'invited',
      }),
    })
    const data = await res.json()
    if (res.ok) setUsers(prev => prev.filter(u => u.user_id !== userId))
    else console.error('Failed to invite user:', data.message)
  }

  return (
    <div className="w-full border-b border-gray-300 pb-4">

      <button
        onClick={() => setOpen(o => !o)}
        className="flex space-x-2 items-center w-full text-left text-black-600 font-semibold focus:outline-none"
      >
        <span>Invite users to join</span>
        <span className={`transform transition-transform duration-200 ${open ? 'rotate-90' : 'rotate-0'}`}>
          ▸
        </span>
      </button>

      <div
        className={`
          overflow-y-auto overflow-x-hidden
          transition-[max-height,opacity] duration-300 ease-in-out
          ${open ? 'max-h-50 opacity-100' : 'max-h-0 opacity-0'}
        `}
      >
        {users?.length > 0 ? (
          <ul className="mt-2 border border-gray-200 rounded p-2 space-y-2 shadow">
            {users.map(user => (
              <li key={user.user_id} className="flex justify-between items-center border-b border-gray-200 pb-2 last:border-b-0">
                <Author author={user} size="sm" />
                <button
                  onClick={() => inviteUser(user.user_id)}
                  className="bg-sky-600/60 hover:bg-sky-800/60 text-white px-2 py-1 rounded text-sm"
                >
                  Invite
                </button>
              </li>
            ))}
          </ul>
        ) : (
          <p className="mt-2 text-gray-500">No available users to invite.</p>
        )}
      </div>
    </div>
  )
}

