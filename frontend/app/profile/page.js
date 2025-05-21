'use client'
import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { useSearchParams } from 'next/navigation'
import Link from 'next/link'

export default function ProfilePage() {
  const searchParams = useSearchParams()
  const router = useRouter()
  const userId = searchParams.get('user_id') // this is your query param
  const [user, setUser] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [following, setFollowing] = useState([]);
  const [followers, setFollowers] = useState([]);
  const { user: localUser } = useUser();

 // Fetch profile data
  const fetchProfile = async () => {
    setLoading(true);
    setError(null);
    try {
      const res = await fetch(`http://localhost:8080/api/profile?user_id=${userId}`, {
        method: 'GET',
        credentials: 'include', // Include cookies for authentication
        headers: {
          'Content-Type': 'application/json',
        },
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

  const fetchFollowing = async () => {
    try {
      const res = await fetch(`http://localhost:8080/api/followers?user_id=${userId}`, {
        method: 'GET',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
      });

      if (res.ok) {
        const data = await res.json();
        setFollowing(data); // Save following data for the viewed profile
        console.log('Fetched following:', data); // Log the fetched following
      } else if (res.status === 401) {
        router.push('/login');
      } else {
        setError('Failed to fetch following');
      }
    } catch (err) {
      console.error('Error fetching following:', err);
      setError('An error occurred while fetching following');
    }
  };

  const fetchFollowers = async () => {
    try {
      const res = await fetch(`http://localhost:8080/api/followers?user_id=${userId}`, {
        method: 'GET',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
      });

      if (res.ok) {
        const data = await res.json();
        setFollowers(Array.isArray(data) ? data : []);
      } else if (res.status === 401) {
        router.push('/login');
      } else {
        setError('Failed to fetch followers');
      }
    } catch (err) {
      console.error('Error fetching followers:', err);
      setError('An error occurred while fetching followers');
    }
  };

  useEffect(() => {
    if (userId) {
      fetchProfile();
      fetchFollowing();
      fetchFollowers(); 
    }
  }, [userId]); // Fetch profile data when userId changes

  const handleFollow = async (status) => {
    try {
      // Disable the button while processing the request
      setLoading(true);

      const res = await fetch(`http://localhost:8080/api/request`, {
        method: 'POST',
        credentials: 'include', // Include cookies for authentication
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          receiver: {
            user_id: Number(userId), // Profile owner's ID
          },
          status: status, // "follow" or "unfollow"
        }),
      });

      if (res.ok) {
        // Reuse fetchProfile to fetch the updated profile data
        await fetchProfile();
      } else {
        console.error(`Failed to ${status} user`);
      }
    } catch (err) {
      console.error(`Error trying to ${status} user:`, err);
    } finally {
      // Re-enable the button
      setLoading(false);
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
          <h3 className="text-lg font-semibold">My followers</h3>
          {Array.isArray(followers) && followers.length > 0 ? (
            <ul>
              {followers.map(f => (
                <li key={f.user_id} className="flex items-center space-x-2">
                  <img
                    src={f.avatar_path ? `http://localhost:8080${f.avatar_path}` : '/avatar.png'}
                    alt={f.nickname || `${f.first_name} ${f.last_name}`}
                    className="w-6 h-6 rounded-full object-cover"
                  />
                  <span>{f.nickname || `${f.first_name} ${f.last_name}`}</span>
                </li>
              ))}
            </ul>
          ) : (
            <p>No followers yet.</p>
          )}
        </div>

        <div className="mb-4">
          <h3 className="text-lg font-semibold">Followers</h3>
          <p>{user.followers_count > 0 ? `${user.followers_count} followers` : 'No followers yet.'}</p>
        </div>


        <div className="mb-4">
          <h3 className="text-lg font-semibold">Users following</h3>
          {Array.isArray(following) && following.length > 0 ? (
            <ul>
              {following.map(f => (
                <li key={f.user_id} className="flex items-center space-x-2">
                  <img
                    src={f.avatar_path ? `http://localhost:8080${f.avatar_path}` : '/avatar.png'}
                    alt={f.nickname || `${f.first_name} ${f.last_name}`}
                    className="w-6 h-6 rounded-full object-cover"
                  />
                  <span>{f.nickname || `${f.first_name} ${f.last_name}`}</span>
                </li>
              ))}
            </ul>
          ) : (
            <p>No followers yet.</p>
          )}
        </div>

         <div className="mb-4">
          <h3 className="text-lg font-semibold">Following</h3>
          <p>{user.following_count > 0 ? `${user.following_count} following` : 'Not following anyone yet.'}</p>
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