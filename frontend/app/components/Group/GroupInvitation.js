'use client'
import { useEffect, useState } from 'react'

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
    const user = JSON.parse(localStorage.getItem('user'))
    const myUserId = user?.user_id

    const res = await fetch('http://localhost:8080/api/request', {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        group: { group_id: groupId },
        sender: { user_id: myUserId },
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
    <div className="w-full max-w-md mt-6">
      <h3 className="text-lg font-semibold mb-2">Invite Users to Join</h3>
      {users?.length === 0 ? (
        <p className="text-gray-500">No available users to invite.</p>
      ) : (
        <ul className="max-h-64 overflow-y-auto border rounded p-2 space-y-2">
          {users?.map(user => (
            <li key={user.user_id} className="flex justify-between items-center border-b pb-1">
              <span>{user.nickname || `${user.first_name} ${user.last_name}`}</span>
              <button
                onClick={() => inviteUser(user.user_id)}
                className="bg-blue-500 hover:bg-blue-600 text-white px-2 py-1 rounded text-sm"
              >
                Invite
              </button>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
