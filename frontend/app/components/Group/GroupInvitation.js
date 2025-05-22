'use client'
import { useEffect, useState } from 'react'
import Author from '../Author'

export default function GroupInvitation({ groupId }) {
  const [users, setUsers] = useState([])

  useEffect(() => {
    async function fetchAvailableUsers() {
      const res = await fetch(`http://localhost:8080/api/group/invite?group_id=${groupId}`, {
        credentials: 'include',
        method: 'GET',})

      const data = await res.json()
      if (res.ok) {
        console.log('Fetched available users:', data)
        setUsers(data)
      } else {
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
        receiver: { user_id: userId },
        status: 'invited',
      }),
    })

    const data = await res.json()
    if (res.ok) {
      console.log(`Invited user ${userId}`)
      setUsers(prev => prev.filter(u => u.user_id !== userId))
    } else {
      console.error('Failed to invite user:', data.message)
    }
  }

  return (
    <div className="w-full border-b border-gray-300 pb-4 pt-4">
      {users?.length > 0 ? (
        <div>
          <ul className="max-h-34 overflow-y-auto border border-gray-200 rounded p-2 space-y-2 shadow">
            {users.map(user => (
              <li
                key={user.user_id}
                className="flex justify-between items-center border-b border-gray-200 pb-2 last:border-b-0 last:pb-0"
              >
                <Author author={user} size="sm" />

                <button
                  onClick={() => inviteUser(user.user_id)}
                  className="bg-blue-500 hover:bg-blue-600 text-white px-2 py-1 rounded text-sm"
                >
                  Invite
                </button>
              </li>
            ))}
          </ul>
        </div>
      ) : (
        <p className="text-gray-500">No available users to invite.</p>
      )}
    </div>
  )
}
