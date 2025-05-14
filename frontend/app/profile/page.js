'use client'
import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'

export default function ProfilePage() {
  const router = useRouter()
  const [user, setUser] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  // Fetch profile data
useEffect(() => {
  const fetchProfile = async () => {
    try {
      const res = await fetch('http://localhost:8080/api/profile', {
        method: 'GET',
        credentials: 'include', // Include cookies for authentication
      });

      if (res.ok) {
        const data = await res.json();
        setUser(data);
      } else if (res.status === 401) {
        router.push('/login'); // Redirect to login if unauthorized
      } else {
        setError('Failed to fetch profile data');
      }
    } catch (err) {
      console.error('Error fetching profile:', err);
      setError('An error occurred while fetching profile data');
    } finally {
      setLoading(false);
    }
  };

  fetchProfile();
}, [router]);

  // Handle loading state
  if (loading) {
    return <div>Loading...</div>
  }

  // Handle error state
  if (error) {
    return <div>Error: {error}</div>
  }

  // Render profile page
  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-50">
      <div className="p-8 max-w-md w-full bg-white rounded-lg shadow-lg">
        <h1 className="text-2xl mb-4 text-center">Profile</h1>
        <div className="mb-4">
          <img
            src={user.profilePicture || '/default-avatar.png'} // Fallback to a default avatar if none exists
            alt="Profile"
            className="w-32 h-32 rounded-full mx-auto"
          />
          <h2 className="text-xl text-center mt-4">{user.nickname || `${user.firstName} ${user.lastName}`}</h2>
          <p className="text-center text-gray-600">{user.email}</p>
        </div>

        <div className="mb-4">
          <h3 className="text-lg font-semibold">Followers</h3>
          <ul>
            {user.followers && user.followers.length > 0 ? (
              user.followers.map((follower, index) => (
                <li key={index}>{follower}</li>
              ))
            ) : (
              <p>No followers yet.</p>
            )}
          </ul>
        </div>

        <div className="mb-4">
          <h3 className="text-lg font-semibold">Following</h3>
          <ul>
            {user.following && user.following.length > 0 ? (
              user.following.map((following, index) => (
                <li key={index}>{following}</li>
              ))
            ) : (
              <p>Not following anyone yet.</p>
            )}
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