'use client'
import { useState, useEffect, Suspense } from 'react'
import { useRouter } from 'next/navigation'
import { useRef } from 'react';
import Link from 'next/link'
import Author from '../components/Author'
import ChatWindow from '../components/ChatWindow'
import ErrorMessage from '../components/ErrorMessage'

function ProfileContent() {
  const router = useRouter()
  const [user, setUser] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [followers, setFollowers] = useState([]);
  const [following, setFollowing] = useState([]);
  const [showFollowingList, setShowFollowingList] = useState(false);
  const [showFollowersList, setShowFollowersList] = useState(false);
  const [chatUserId, setChatUserId] = useState(null);
  const [showUnfollowPopup, setShowUnfollowPopup] = useState(false);
  const [showPrivacyPopup, setShowPrivacyPopup] = useState(false);
  const [privacyPopupText, setPrivacyPopupText] = useState('');
  const unfollowTimeoutRef = useRef(null);
  const privacyTimeoutRef = useRef(null);

  const [userId, setUserId] = useState(null)

  const [currentUser, setCurrentUser] = useState(null);

    useEffect(() => {
    const storedUser = localStorage.getItem("user");
    if (storedUser) {
      try {
        setCurrentUser(JSON.parse(storedUser));
      } catch (e) {
        setCurrentUser(null);
      }
    } else {
      setCurrentUser(null);
    }
  }, []);

  useEffect(() => {
    // Only run in browser
    const params = new URLSearchParams(window.location.search)
             const id = params.get('user_id')
      if (id !== null && id !== undefined && id !== "") {
        setUserId(id)
      } else {
        setError('No user_id provided in URL')
      }
  }, [])

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

      const data = await res.json();
      if (res.ok) {
        setUser(data); // Save user data
      } else if (res.status === 401) {
        router.push('/login'); // Redirect to login if unauthorized
      } else {
        setError(data.message || 'Failed to fetch profile data');
      }
    } catch (err) {
      console.error('Error fetching profile:', err);
      setError(err.message || 'An error occurred while fetching profile data');
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

      const data = await res.json();
      if (res.ok) {
        setFollowers(Array.isArray(data) ? data : []);
      } else if (res.status === 401) {
        router.push('/login');
      } else {
        setError(data.message || 'Failed to fetch followers');
      }
    } catch (err) {
      console.error('Error fetching followers:', err);
      setError(err.message || 'An error occurred while fetching followers');
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

      const data = await res.json();
      if (res.ok) {
        setFollowing(Array.isArray(data) ? data : []);
      } else if (res.status === 401) {
        router.push('/login');
      } else {
        setError(data.message || 'Failed to fetch following');
      }
    } catch (err) {
      console.error('Error fetching following:', err);
      setError(err.message || 'An error occurred while fetching following');
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

  const handlePrivacy = async (isPublic) => {
    // Optimistically update UI for animation
    setUser(prev => ({
      ...prev,
      user: {
        ...prev.user,
        is_public: isPublic,
      },
    }));

    try {
      const res = await fetch(`http://localhost:8080/api/privacy`, {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          user_id: Number(userId),
          is_public: isPublic,
        }),
      });

      if (res.ok) {
        
    

        // Re-fetch profile and followers/follow_requests to update UI
        setTimeout(async () => {
              if (privacyTimeoutRef.current) clearTimeout(privacyTimeoutRef.current);
        setPrivacyPopupText(isPublic ? "Your profile is now public." : "Your profile is now private.");
        setShowPrivacyPopup(true);
        privacyTimeoutRef.current = setTimeout(() => {
          setShowPrivacyPopup(false);
          privacyTimeoutRef.current = null;
        }, 2000);
        await fetchProfile();
        await fetchFollowers();
      }, 250);
      
      } else {
        console.log('Failed to update privacy settings');
      }
    } catch (err) {
      setUser(prev => ({
        ...prev,
        user: {
          ...prev.user,
          is_public: !isPublic,
        },
      }));
      console.error('Error updating privacy settings:', err);
    }
  };
  
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
         if (status === 'unfollow') {
        // Clear any previous timeout
        if (unfollowTimeoutRef.current) {
          clearTimeout(unfollowTimeoutRef.current);
        }
        setShowUnfollowPopup(true);
        unfollowTimeoutRef.current = setTimeout(() => {
          setShowUnfollowPopup(false);
          unfollowTimeoutRef.current = null;
        }, 2000);
      } else {
        // If you follow again, hide the popup immediately
        if (unfollowTimeoutRef.current) {
          clearTimeout(unfollowTimeoutRef.current);
          unfollowTimeoutRef.current = null;
        }
        setShowUnfollowPopup(false);
      }
          // if the profile is private and the user is following
        if (status === 'follow' && user && user.user && !user.user.is_public) {
          setShowUnfollowPopup(false);
          setUser(prev => ({
            ...prev,
            has_requested: true
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

  // Handle error state
    if (error) {
      return (
        <div className="p-4">
          <ErrorMessage message={error} />
        </div>
      )
    }

  // Handle loading state
  if (loading) {
    return <div>Loading...</div>
  }

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
{ /*unfollow and privacy popups*/ }
      {showUnfollowPopup && (
        <div
          className="fixed top-8 left-1/2 transform -translate-x-1/2 bg-gray-800 text-white px-4 py-2 rounded shadow-lg z-50
            transition-all duration-500 ease-out
            animate-popup"
          style={{
          animation: 'popup 0.8s cubic-bezier(0.33, 1, 0.68, 1)'
        }}
        >
          You have unfollowed {user.user.nickname || user.user.first_name}.
        </div>
      )}

      {showPrivacyPopup && (
        <div
          className="fixed top-20 left-1/2 transform -translate-x-1/2 bg-blue-700 text-white px-4 py-2 rounded shadow-lg z-50
            transition-all duration-500 ease-out animate-popup"
          style={{
            animation: 'popup 0.8s cubic-bezier(0.33, 1, 0.68, 1)'
          }}
        >
          {privacyPopupText}
        </div>
      )}

      <div className="p-8 max-w-md w-full bg-white rounded-lg shadow-lg">
        <h1 className="text-2xl mb-4 text-center">Profile</h1>
      {user.is_own_profile && (
        <div className="flex items-center justify-center mb-4">
          <span className={`mr-2 text-sm ${user.user.is_public ? 'text-gray-400' : 'text-blue-600'}`}>
            Private
          </span>
          <button
            type="button"
            className={`relative inline-flex h-6 w-10 items-center rounded-full transition-colors focus:outline-none ${
              user.user.is_public ? 'bg-green-500' : 'bg-blue-600'
            }`}
            onClick={() => handlePrivacy(!user.user.is_public)}
            aria-pressed={user.user.is_public}
          >
            <span
              className={`inline-block h-5 w-5 transform rounded-full bg-white shadow transition-transform duration-200 ${
                user.user.is_public ? 'translate-x-4' : 'translate-x-1'
              }`}
            />
          </button>
          <span className={`ml-2 text-sm ${user.user.is_public ? 'text-green-600' : 'text-gray-400'}`}>
            Public
          </span>
        </div>
      )}
        <div className="mb-4">
          <img
            src={user.user.avatar_path ? `http://localhost:8080${user.user.avatar_path}` : '/avatar.png'}// Fallback to a default avatar if none exists
            alt="Profile"
            className="w-32 h-32 rounded-full mx-auto object-cover"
            />
          <h2 className="text-xl text-center mt-4">{user.user.nickname || user.user.first_name}</h2>
          { (user.user.is_public || user.is_follower || user.is_own_profile) &&
          (<p className="text-center text-gray-500">{user.user.first_name} {user.user.last_name}</p>)
          }
          { (user.user.is_public || user.is_follower || user.is_own_profile) &&
          (<p className="text-center text-gray-500">{user.user.date_of_birth && (
            new Date(user.user.date_of_birth).toLocaleDateString('fi-FI', {
              day: '2-digit',
              month: '2-digit',
              year: 'numeric',
            })
          )}
          </p>)
          }
          { (user.user.is_public || user.is_follower || user.is_own_profile) && 
          (<p className="text-center text-gray-600">{user.user.email}</p>)
          }
        </div> 
        {!user.is_own_profile && ( user.is_follower ? (
          <button
          className="bg-red-500 text-white px-4 py-2 rounded hover:bg-red-600"
          onClick={() => handleFollow('unfollow')}
          >
            Unfollow
          </button>
        ) : user.has_requested ? (
          <p className="text-yellow-600 font-semibold text-center">
              Follow request sent. Waiting for approval.
          </p>

        ) : (
        <button
            className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600"
            onClick={() => handleFollow('follow')}
            disabled={loading || user.has_requested}
        >
            Follow
          </button>
        ))}
        {/* Chat button */}
        {(!user.is_own_profile && user.show_chat_button) && (
          <button
            onClick={() => setChatUserId(user.user.user_id)}
            className="bg-blue-600 text-white px-4 py-2 rounded"
          >
            Chat
          </button>
        )}

        {chatUserId && (
          <ChatWindow
            chatPartner={user.user}
            isGroupChat={false}
            onClose={() => setChatUserId(null)}
            currentUser={currentUser ? currentUser.user_id : null}
            group={null}
          />
        )}

        <div className="mb-4">
          <h3 className="text-lg font-semibold">About me</h3>
          <p>{user.user.about_me || 'No bio available.'}</p>
        </div>

        <div className="mb-4">
          <h3 className="text-lg font-semibold">Followers: <button
            className="text-blue-600 cursor-pointer p-0 bg-transparent border-none"
            onClick={() => setShowFollowersList(true)}
            disabled={followers.length === 0 || (!user.is_own_profile && !user.user.is_public && !user.is_follower)}
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
                  <ul>
                    {followers.map(f => (
                      <li key={f.user_id} className="mb-2">
                        <span className="flex items-center space-x-2 hover:underline">
                          <Author author={f} size="sm" />
                        </span>
                      </li>
                    ))}
                  </ul>
              </div>
            </div>
          )}
        </div>

        <div className="mb-4">
          <h3 className="text-lg font-semibold">Following: <button
            className="text-blue-600 cursor-pointer p-0 bg-transparent border-none"
            onClick={() => setShowFollowingList(true)}
            disabled={following.length === 0 || (!user.is_own_profile && !user.user.is_public && !user.is_follower)}
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
                <ul>
                  {following.map(f => (
                    <li key={f.user_id} className="mb-2">
                      <span className="flex items-center space-x-2 hover:underline">
                        <Author author={f} size="sm"/>
                      </span>
                    </li>
                  ))}
                </ul>
            </div>
          </div>
        )}
        </div>

        <div className="mb-4">
          <h3 className="text-lg font-semibold mb-4">My posts</h3>
          <ul>
            {user.posts ? (
              user.posts.map((post, index) => (
                <li key={index} className="block mb-4 p-4 rounded shadow border border-gray-200 bg-white hover:bg-gray-100 transition cursor-pointer">
                  <Link href={`/post?post_id=${post.post_id}`}>
                    <div className="flex justify-between items-center mb-2">
                      <h4 className="text-lg font-semibold">{post.post_title || 'Untitled Post'}</h4>
                      <div className="text-right">
                        <p className="text-xs text-gray-500">
                          {new Date(post.created_at).toLocaleString('en-GB', {
                            day: '2-digit',
                            month: '2-digit',
                            year: 'numeric',
                            hour: '2-digit',
                            minute: '2-digit',
                            hour12: false,
                            timeZone: 'UTC',
                          })}
                        </p>
                        <p className="text-xs text-gray-500">
                          {Array.isArray(post.comments)
                            ? `${post.comments.length} comment${post.comments.length === 1 ? '' : 's'}`
                            : '0 comments'}
                        </p>
                      </div>
                    </div>
                    <p className="text-gray-700">
                        {post.post_content.length > 50
                          ? post.post_content.slice(0, 50) + '...'
                          : post.post_content}
                    </p>
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

export default function ProfilePage() {
  return (
    <Suspense fallback={<div className="flex items-center justify-center min-h-screen bg-gray-50">Loading profile...</div>}>
      <ProfileContent />
    </Suspense>
  )
}