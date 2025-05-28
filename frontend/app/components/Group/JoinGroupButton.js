'use client'
import { useState } from 'react'

export default function JoinGroupButton({ groupId, onJoin }) {
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)

  async function handleJoin() {
    if (!groupId) {
      console.error('Group ID is required to join a group.')
      return
    }
    const user = JSON.parse(localStorage.getItem('user'))
    const userID = user?.user_id

    setLoading(true)
    setError(null)
    try {
      const res = await fetch(`http://localhost:8080/api/request`, {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ joining_user: {user_id: Number(userID) }, group: {group_id: groupId}, status: "requested" }),
      })

      const data = await res.json()
      if (res.ok) {
        console.log('Join request sent:', data)
      }

      if (onJoin) onJoin()
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="mt-4">
     <p className="text-red-500 font-semibold">Join the group to see all posts.</p>
      <button
        onClick={handleJoin}
        disabled={loading}
        className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
      >
        {loading ? 'Joining...' : 'Request to Join'}
      </button>
      {error && <p className="text-red-500 text-sm mt-2">{error}</p>}
    </div>
  )
}
