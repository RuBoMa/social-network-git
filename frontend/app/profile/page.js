'use client'
import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
// import Link from 'next/link'
// import React, { useState, useEffect, use } from 'react';

// export default function ProfilePage() {
//   const router = useRouter()
//   const [user, setUser] = useState(null)

//   const mockUser = {
//     username: 'john_doe',
//     email: 'john@example.com',
//     profilePicture: 'https://example.com/avatar.jpg',
//     followers: ['follower1', 'follower2', 'follower3'],
//     following: ['user1', 'user2']
//   }


//    // Fetch profile data
//    useEffect(() => {
//     const fetchProfile = async () => {
//       try {
//         const res = await fetch('http://localhost:8080/api/profile', {
//           method: 'GET',
//           credentials: 'include',
//         })

//         if (res.ok) {
//           const data = await res.json()
//           setUser(data)
//         } else {
//           //if backend is not available, use mock data
//           setUser(mockUser)
//         }
//       } catch (error) {
//         // if error occurs, use mock data
//         setUser(mockUser)
//       }
//     }

//     fetchProfile()
//   }, [])

//   // if the user has not logged in
//   if (!user) {
//     return <div>Loading...</div>
//   }

//   return (
//     <div className="flex items-center justify-center min-h-screen bg-gray-50">
//       <div className="p-8 max-w-md w-full bg-white rounded-lg shadow-lg">
//         <h1 className="text-2xl mb-4 text-center">Profile</h1>
//         <div className="mb-4">
//           <img
//             src={user.profilePicture} // Assuming the profile picture URL is in user.profilePicture
//             alt="Profile"
//             className="w-32 h-32 rounded-full mx-auto"
//           />
//           <h2 className="text-xl text-center mt-4">{user.username}</h2>
//           <p className="text-center text-gray-600">{user.email}</p>
//         </div>

//         <div className="mb-4">
//           <h3 className="text-lg font-semibold">Followers</h3>
//           <ul>
//             {user.followers.map((follower, index) => (
//               <li key={index}>{follower}</li>
//             ))}
//           </ul>
//         </div>

//         <div className="mb-4">
//           <h3 className="text-lg font-semibold">Following</h3>
//           <ul>
//             {user.following.map((following, index) => (
//               <li key={index}>{following}</li>
//             ))}
//           </ul>
//         </div>

//         <Link
//           href="/edit-profile"
//           className="block mt-4 text-blue-500 text-center"
//         >
//           Edit Profile
//         </Link>
//       </div>
//     </div>
//   )
// }

// Replace with this when backend is implemented!:

// import Profile from '../../components/Profile'

// export const metadata = {
//   title: 'Your Profile',
// }

// export default async function ProfilePage() {
//   // call backend and fetch user data, example url!
//   const res = await fetch('http://localhost:8080/api/users/me', {
//     credentials: 'include',
//     cache: 'no-store'
//   })
//   if (!res.ok) {
//     throw new Error('Failed to fetch profile')
//   }
//   const user = await res.json()

//   // render and return current userdata
//   return <Profile user={user} />
// }

import Profile from '../components/Profile'


export default function ProfilePage() {
  console.log("ProfilePage component rendered")
  const router = useRouter()

  useEffect(() => {
    // Check if the user is logged in
    const loggedIn = localStorage.getItem('loggedIn')
    console.log("Logged in status:", loggedIn)
    if (!loggedIn) {
      console.log("User not logged in, redirecting to login page")
      // Redirect to login page if not logged in
      router.push('/login')
    }
  }
  , [])



  // ==== MOCK DATA ====
  const user = {
    // now served from public/
    avatar_url: '/sky.jpg',
    first_name: 'Jane',
    last_name: 'Doe',
    nickname: 'jdoe',
    followers_count: 4, // maybe don't need, calculate on length of followers?
    followers: [
      { id: 1, avatar_url: '/sky.jpg', first_name: 'Alice', last_name: 'Smith', nickname: 'alice' },
      { id: 2, avatar_url: '/sky.jpg', first_name: 'Bob', last_name: 'Brown', nickname: 'bobby' },
      { id: 3, avatar_url: '/sky.jpg', first_name: 'Carol', last_name: 'Jones', nickname: '' },
      { id: 4, avatar_url: '/sky.jpg', first_name: 'Kira', last_name: 'Kirasson', nickname: 'kiki' }
    ]
  }
  return <Profile user={user} />
}
