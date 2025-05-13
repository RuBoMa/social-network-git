'use client'
import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import React, { useState, useEffect, use } from 'react';

export default function ProfilePage() {
  const router = useRouter()
  const [user, setUser] = useState(null)


   // Fetch profile data
   useEffect(() => {
    const fetchProfile = async () => {
      try {
        const res = await fetch('http://localhost:8080/api/profile', {
          method: 'GET',
          credentials: 'include',
        })

        if (res.ok) {
          const data = await res.json()
          setUser(data)
        } else {
          //if backend is not available, use mock data
          setUser(mockUser)
        }
      } catch (error) {
        // if error occurs, use mock data
        setUser(mockUser)
      }
    }

    fetchProfile()
  }, [])

  // if the user has not logged in
  if (!user) {
    return <div>Loading...</div>
  }

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-50">
      <div className="p-8 max-w-md w-full bg-white rounded-lg shadow-lg">
        <h1 className="text-2xl mb-4 text-center">Profile</h1>
        <div className="mb-4">
          <img
            src={user.profilePicture} // Assuming the profile picture URL is in user.profilePicture
            alt="Profile"
            className="w-32 h-32 rounded-full mx-auto"
          />
          <h2 className="text-xl text-center mt-4">{user.username}</h2>
          <p className="text-center text-gray-600">{user.email}</p>
        </div>

        <div className="mb-4">
          <h3 className="text-lg font-semibold">Followers</h3>
          <ul>
            {user.followers.map((follower, index) => (
              <li key={index}>{follower}</li>
            ))}
          </ul>
        </div>

        <div className="mb-4">
          <h3 className="text-lg font-semibold">Following</h3>
          <ul>
            {user.following.map((following, index) => (
              <li key={index}>{following}</li>
            ))}
          </ul>
        </div>

        <Link
          href="/edit-profile"
          className="block mt-4 text-blue-500 text-center"
        >
          Edit Profile
        </Link>
      </div>
    </div>
  )
}

