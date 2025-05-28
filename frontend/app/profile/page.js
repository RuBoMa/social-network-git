'use client'
import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { useSearchParams } from 'next/navigation'
import Link from 'next/link'
import Author from '../components/Author'

export default function ProfilePage() {
  const searchParams = useSearchParams()
  const router = useRouter()
  const userId = searchParams.get('user_id') // this is your query param
  const [user, setUser] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [followers, setFollowers] = useState([]);
  const [following, setFollowing] = useState([]);
  const [showFollowingList, setShowFollowingList] = useState(false);
  const [showFollowersList, setShowFollowersList] = useState(false);

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

  const fetchFollowing = async () => {
    try {
      const res = await fetch(`http://localhost:8080/api/following?user_id=${userId}`, {
        method: 'GET',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
      });

      if (res.ok) {
        const data = await res.json();
        setFollowing(Array.isArray(data) ? data : []);
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

  useEffect(() => {
    if (userId) {
      fetchProfile();
      fetchFollowers(); 
      fetchFollowing(); // Fetch following data
      setShowFollowingList(false); // <-- Close the popup when userId changes
      setShowFollowersList(false); // <-- Close the popup when userId changes

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
          // if the profile is private and the user is following, set request_status to 'requested'
        if (status === 'follow' && user && user.user && !user.user.is_public) {
          setUser(prev => ({
            ...prev,
            request_status: 'requested'
          }));
          return; // Exit early if the profile is private
        }
        await fetchProfile();
        await fetchFollowers(); // Fetch updated followers list
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

    const handleFollowRequest = async (requestId, status) => {
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
          request_id: requestId, // Use request_id if available
          status: status, // "accepted" or "declined"
        }),
      });

      if (res.ok) {
        // Reuse fetchProfile to fetch the updated profile data
        await fetchProfile();
        await fetchFollowers(); // Fetch updated followers list
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
      <div className="flex flex-col items-center w-full max-w-md">
      {user.follow_requests && user.follow_requests.length > 0 && (
        <div className="mb-6">
        <h3 className="text-lg font-semibold mb-2">Pending Follow Requests</h3>
        <ul>
        {user.follow_requests.map(req => (
          <li key={req.request_id} className="flex items-center gap-2 mb-2">
            <Author author={req.sender} size="sm" />
            <button
              className="bg-green-500 text-white px-2 py-1 rounded"
              onClick={() => handleFollowRequest(req.request_id, 'accepted')}
            >
              Accept
            </button>
            <button
              className="bg-red-500 text-white px-2 py-1 rounded"
              onClick={() => handleFollowRequest(req.request_id, 'declined')}
            >
             Decline
            </button>
          </li>
       ))}
      </ul>
    </div>
)}
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
        ) : user.request_status === 'requested' ? (
          <p className="text-yellow-600 font-semibold text-center">
              Follow request sent. Waiting for approval.
          </p>

        ) : (
        <button
            className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600"
            onClick={() => handleFollow('follow')}
            disabled={loading || user.request_status === 'requested'}
        >
            Follow
          </button>
        )}

        <div className="mb-4">
          <h3 className="text-lg font-semibold">About me</h3>
          <p>{user.user.about_me || 'No bio available.'}</p>
        </div>

        <div className="mb-4">
          <h3 className="text-lg font-semibold">Followers: <button
            className="text-blue-600 cursor-pointer p-0 bg-transparent border-none"
            onClick={() => setShowFollowersList(true)}
            disabled={followers.length === 0}
            >
            {followers.length > 0 ? `${followers.length}` : '0'}
          </button></h3>
          
          {showFollowersList && (
            <div
            className="fixed inset-0 flex items-center justify-center z-50"
            style={{ background: 'rgba(0,0,0,0.2)' }}
            onClick={() => setShowFollowersList(false)}
            >
              <div
                className="bg-white rounded-lg shadow-lg p-6 max-w-sm w-full relative"
                onClick={e => e.stopPropagation()}
                >
                <button
                  className="absolute top-2 right-2 text-gray-500 hover:text-gray-700"
                  onClick={() => setShowFollowersList(false)}
                  >
                  ✕
                </button>
                <h4 className="text-lg font-semibold mb-4">Followers</h4>
                {followers.length > 0 ? (
                  <ul>
                    {followers.map(f => (
                      <li key={f.user_id} className="mb-2">
                        <Link href={`/profile?user_id=${f.user_id}`}>
                          <span className="flex items-center space-x-2 hover:underline">
                           <Author author={f} size="sm" />
                          </span>
                        </Link>
                      </li>
                    ))}
                  </ul>
                ) : (
                  <p>No followers yet.</p>
                )}
              </div>
            </div>
          )}
        </div>


        <div className="mb-4">
          <h3 className="text-lg font-semibold">Following: <button
            className="text-blue-600 cursor-pointer p-0 bg-transparent border-none"
            onClick={() => setShowFollowingList(true)}
            disabled={following.length === 0}
            >
            {user.following_count > 0 ? `${user.following_count}` : '0'}
          </button></h3>
          
          {showFollowingList && (
            <div
            className="fixed inset-0 flex items-center justify-center z-50"
            style={{ background: 'rgba(0,0,0,0.2)' }}
            onClick={() => setShowFollowingList(false)} // Close on overlay click
            >
            <div
              className="bg-white rounded-lg shadow-lg p-6 max-w-sm w-full relative"
              onClick={e => e.stopPropagation()} // Prevent closing when clicking inside the modal
              >
              <button
                className="absolute top-2 right-2 text-gray-500 hover:text-gray-700"
                onClick={() => setShowFollowingList(false)}
                >
                ✕
              </button>
              
              <h4 className="text-lg font-semibold mb-4">Following</h4>
              {following.length > 0 ? (
                <ul>
                  {following.map(f => (
                    <li key={f.user_id} className="mb-2">
                      <Link href={`/profile?user_id=${f.user_id}`}>
                        <span className="flex items-center space-x-2 hover:underline">
                         <Author author={f} size="sm"/>
                        </span>
                      </Link>
                    </li>
                  ))}
                </ul>
              ) : (
                <p>No following yet.</p>
              )}
            </div>
          </div>
        )}
        </div>

        <div className="mb-4">
          <h3 className="text-lg font-semibold mb-4">My posts</h3>
          <ul>
            {user.posts ? (
              user.posts.map((post, index) => (
                <li key={index} className="mb-4 border border-gray-300 rounded-lg p-4 bg-white shadow-sm hover:bg-gray-100 transition">
                  <Link href={`/post?post_id=${post.post_id}`}>
                    <div className="cursor-pointer">
                      <h4 className="text-lg font-semibold">{post.post_title || 'Untitled Post'}</h4>
                        <p className="text-gray-700">
                            {post.post_content.length > 50
                              ? post.post_content.slice(0, 50) + '...'
                              : post.post_content}
                        </p>
                      <p className="text-sm text-gray-500">
                        {post.created_at ? new Date(post.created_at).toLocaleString() : 'Unknown Date'}
                      </p>
                    </div>
                  </Link>
                </li>
              ))
            ) : (
              <p>No posts yet.</p>
            )}
          </ul>
          </div>
        </div>

      </div>
    </div>
  )
}