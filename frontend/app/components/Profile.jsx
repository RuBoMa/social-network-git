// frontend/components/Profile.jsx (client)
'use client'
import { useState } from 'react'
import FollowersModal from './FollowersModal'

export default function Profile({ user }) {
  const [show, setShow] = useState(false)
  const name = user.nickname || `${user.first_name} ${user.last_name}`

  return (
    <>
      <div className="flex items-center">
        <img
          src={user.avatar_url}
          alt="Avatar"
          className="w-16 h-16 rounded-full object-cover"
        />
        <h2 className="ml-4 text-2xl font-semibold">{name}</h2>
      </div>

      <div className="relative inline-block">
        <button onClick={() => setShow(!show)}>
          {user.followers_count} followers
        </button>
        {show && <FollowersModal followers={user.followers} onClose={() => setShow(false)} />}
      </div>
    </>
  )
}