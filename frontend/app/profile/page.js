'use client'
import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { useSearchParams } from 'next/navigation'
import { useUser } from '../context/UserContext';
import Link from 'next/link'

export default function ProfilePage() {
  const searchParams = useSearchParams()
  const router = useRouter()
  const userId = searchParams.get('user_id') // this is your query param
  const [user, setUser] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const status = 'follow'


  // Fetch profile data
  useEffect(() => {
    const fetchProfile = async () => {
      setLoading(true);
      setError(null);
      try {
        const res = await fetch(`http://localhost:8080/api/profile?user_id=${userId}`,{
          method: 'GET',
          credentials: 'include', // Include cookies for authentication
          headers: {
            'Content-Type': 'application/json',
          }
        });

        if (res.ok) {
          const data = await res.json();
          setUser(data); // Save user data
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
    if (userId) {
        fetchProfile();
      }
  }, [userId]); // Fetch profile data when userId changes

  const localUser = JSON.parse(localStorage.getItem('user')); // Parse the string into an object

  const handleFollow = async (status) => {
  const handleFollow = async (status) => {
    try {
      console.log("userId", userId)
      const res = await fetch(`http://localhost:8080/api/request`, {
        method: 'POST',
        credentials: 'include', // Include cookies for authentication
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          receiver: {
            user_id: Number(userId), },
          sender: {
            user_id: Number(localUser.user_id) },
            status: status,
            user_id: Number(localUser.user_id) },
            status: status,
        }),
      });

      if (res.ok) {
      const data = await res.json();
      console.log(`${status === 'follow' ? 'Followed' : 'Unfollowed'} user:`, data);

      // Update the UI to reflect the new follow state
      setUser((prevUser) => ({
        ...prevUser,
        is_follower: status === 'follow', // Update is_follower based on the action
      }));
      } else {
        console.error(`Failed to ${status} user`);
      }
      } catch (err) {
        console.error(`Error trying to ${status} user:`, err);
      }
  };

  // Handle loading state
  if (loading) {
    return <div>Loading...</div>
  }

  // Handle error state
  if (error) {
    return <div>Error: {error}</div>
  }

  console.log("post data", user)
  // Render profile page
  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-50">
      <div className="p-8 max-w-md w-full bg-white rounded-lg shadow-lg">
        <h1 className="text-2xl mb-4 text-center">Profile</h1>
        <div className="mb-4">
          <img
            src={user.user.avatar_path ? `http://localhost:8080${user.user.avatar_path}` : '/avatar.png'}// Fallback to a default avatar if none exists
            alt="Profile"
            className="w-32 h-32 rounded-full mx-auto object-cover"
          />
          <h2 className="text-xl text-center mt-4">{user.user.nickname || `${user.user.first_name} ${user.user.last_name}`}</h2>
          <p className="text-center text-gray-600">{user.user.email}</p>
        </div> 
        {user.is_own_profile ? (
          <p className="text-center text-gray-500">This is your profile.</p>
        ) : user.is_follower ? (
          <button
            className="bg-red-500 text-white px-4 py-2 rounded hover:bg-red-600"
            onClick={() => handleFollow('unfollow')}
          >
            Unfollow
          </button>
        ) : (
          <button
            className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600"
            onClick={() => handleFollow('follow')}
          >
            Follow
          </button>
        )}

        <div className="mb-4">
          <h3 className="text-lg font-semibold">About me</h3>
          <p>{user.user.about_me || 'No bio available.'}</p>
        </div>

        <div className="mb-4">
          <h3 className="text-lg font-semibold">Followers</h3>
          <ul>
            {user.user.followers && user.user.followers.length > 0 ? (
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
            {user.user.following && user.user.following.length > 0 ? (
              user.following.map((following, index) => (
                <li key={index}>{following}</li>
              ))
            ) : (
              <p>Not following anyone yet.</p>
            )}
          </ul>
        </div>

        <div className="mb-4">
          <h3 className="text-lg font-semibold">Groups</h3>
          <ul>
            {user.user.groups && user.user.groups.length > 0 ? (
              user.groups.map((group, index) => (
                <li key={index}>{group}</li>
              ))
            ) : (
              <p>No groups yet.</p>
            )}
          </ul>
        </div>

        <div className="mb-4">
          <h3 className="text-lg font-semibold">Posts</h3>
          <ul>
            {user.posts ? (
              user.posts.map((post, index) => (
                <li key={index} className="mb-4 border border-gray-300 rounded-lg p-4 bg-white shadow-sm">
                  <h4 className="text-lg font-semibold">{post.post_title || 'Untitled Post'}</h4>
                  <p>{post.post_content || 'No content available.'}</p>
                  {post.post_image && (
                    <img
                      src={`http://localhost:8080${post.post_image}`}
                      alt="Post visual"
                      style={{ maxWidth: '100%' }}
                    />
                  )}
                  <p className="text-sm text-gray-500">
                    {post.created_at ? new Date(post.created_at).toLocaleString() : 'Unknown Date'}
                  </p>
                </li>
              ))
            ) : (
              <p>No posts yet.</p>
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