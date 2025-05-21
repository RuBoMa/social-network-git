'use client'
import { useState } from 'react'
import { useSearchParams } from 'next/navigation'
import { useEffect } from 'react'
import ImageIcon from '../components/AddImageIcon'
import ImageUploadPreview from '../components/ImageUploadPreview'
import Author from '../components/Author'


export default function CreatePost({ onSuccess }) {
  const searchParams = useSearchParams()
  const groupID = searchParams.get('group_id')
  
  const [postTitle, setPostTitle] = useState('')
  const [postContent, setContent] = useState('')
  const [privacy, setPostPrivacy] = useState(groupID ? 'followers' : 'public');
  const [postImage, setPostImage] = useState(null)
  const [followers, setFollowers] = useState([])
  const [selectedUsers, setSelectedUsers] = useState([])
  const [searchTerm, setSearchTerm] = useState('')

  async function handlePost(e) {
    e.preventDefault()

    const formData = new FormData()
    formData.append('post_title', postTitle)
    formData.append('post_content', postContent)
    formData.append('privacy', privacy)
    if (postImage) formData.append('post_image', postImage)
    if (groupID) formData.append('group_id', groupID)
    if (privacy === 'custom') {
      const userIds = selectedUsers.map(user => user.user_id);
      formData.append('custom_users', JSON.stringify(userIds));
    }

    const res = await fetch('http://localhost:8080/api/create-post', {
      method: 'POST',
      credentials: 'include',
      body: formData,
    })

    if (res.ok) {
      setContent('')
      setPostTitle('')
      setPostImage(null)
      setSelectedUsers([])
      setSearchTerm('')
      onSuccess && onSuccess() // trigger reload if provided
    } else {
      alert('Failed to post')
    }
  }

  
  useEffect(() => {
    if (privacy === 'custom' && Array.isArray(followers) && followers.length === 0) {
      fetch('http://localhost:8080/api/followers', {
        credentials: 'include'
      })
      .then(res => res.json())
      .then(data => {
        console.log('Fetched followers:', data)
        setFollowers(data)
      })
      .catch(err => console.error('Failed to fetch followers', err))
      }
    }, [privacy])

    //remove users from selectedUsers if they are not in followers
    async function fetchFollowers() {
      const res = await fetch('http://localhost:8080/api/followers', { credentials: 'include' });
      const data = await res.json();
      setFollowers(data);
    }

    useEffect(() => {
      if (privacy === 'custom' && Array.isArray(followers) && followers.length === 0) {
        fetchFollowers()
      }
    }, [privacy])

    async function handleUnfollow(userId) {
      // ... unfollow query
      if (unfollowSuccess) {
        fetchFollowers()  // updating followers
        // if needed, update selectedUsers
      }
    }
    

    console.log('Followers:', followers)

    const filteredFollowers = (followers || []).filter(user => {
      const search = searchTerm.toLowerCase()
      return (
        (user.first_name && user.first_name.toLowerCase().includes(search)) ||
        (user.last_name && user.last_name.toLowerCase().includes(search)) ||
        (user.nickname && user.nickname.toLowerCase().includes(search))
      ) && !selectedUsers.some(selected => selected.user_id === user.user_id)
    })

    function addUser(user) {
      if (!selectedUsers.some(u => u.user_id === user.user_id)) {
        setSelectedUsers(prev => [...prev, user])
      }
      setSearchTerm('') // empty the search
    }

    function removeUser(userId) {
      setSelectedUsers(prev => prev.filter(u => u.user_id !== userId))
    }

  return (
    <form className="relative max-w-full mx-0 mt-1 p-2 bg-white rounded-xl shadow-md" onSubmit={handlePost}>
      {/* title for the post */}
      <label className="block mb-4">
        <input
          type="text"
          value={postTitle}
          onChange={(e) => setPostTitle(e.target.value)}
          placeholder="Title"
          required
          className="mt-1 block w-full border border-gray-300 rounded p-2"
        />
      </label>
  
      {/* textarea for the post */}
      <label className="block mb-4">
        <textarea
          value={postContent}
          onChange={(e) => setContent(e.target.value)}
          placeholder="What's on your mind?"
          required
          className="mt-1 block w-full border border-gray-300 rounded p-2"
        />
      </label>
      {/* image upload, privacy options and button in one div */}
      <div className="flex flex-wrap items-start justify-between gap-6 mb-4 pt-4">

      {/* image upload */}
        <label className="inline-flex items-center space-x-2 cursor-pointer">
          <input
            type="file"
            accept="image/*"
            onChange={(e) => setPostImage(e.target.files[0])}
            className="hidden"
          />
          <ImageIcon />
        </label>
  
        {postImage && (
          <ImageUploadPreview imageFile={postImage} setImageFile={setPostImage} />
        )}
  
      {/* privacy options */}
  {!groupID && (
        <div className="mb-4">
          {['public', 'followers', 'custom'].map(option => (
            <label key={option} className="inline-flex items-center mr-4 cursor-pointer">
              <input
                type="radio"
                name="privacy"
                value={option}
                checked={privacy === option}
                onChange={(e) => setPostPrivacy(e.target.value)}
                className="form-radio text-blue-600"
              />
              <span className="ml-1 capitalize">{option}</span>
            </label>
          ))}
        </div>
      )}
      <button
        type="submit"
        className="w-auto bg-blue-600 text-white px-4 py-1 rounded hover:bg-blue-700 transition"
      >
        Submit
          </button>
      </div>

      {/* if privacy is custom, then tag-input */}
      {privacy === 'custom' && (
        <div className="mb-4">
          {/* Tag container */}
          <div className="flex flex-wrap gap-2 mb-2">
            {selectedUsers.map(user => (
              <div
                key={user.user_id}
                className="flex items-center bg-blue-200 text-blue-800 rounded px-2 py-1 text-sm"
              >
                 <Author author={user} disableLink={true} size="s" />
                <button
                  type="button"
                  onClick={() => removeUser(user.user_id)}
                  className="ml-1 font-bold"
                  aria-label="Remove user"
                >
                  &times;
                </button>
              </div>
            ))}
          </div>

          {/* Search followers */}
          <input
            type="text"
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            placeholder="Search followers..."
            className="w-full border border-gray-300 rounded p-2"
          />

          {/* Dropdown with filtered followers */}
          {searchTerm && filteredFollowers.length > 0 && (
            <ul className="border border-gray-300 rounded max-h-30 overflow-auto mt-1 bg-white shadow-sm z-10 relative">
              {filteredFollowers.map(user => (
                <li
                  key={user.user_id}
                  onClick={() => addUser(user)}
                  className="cursor-pointer px-3 py-1 hover:bg-blue-100"
                >
                   <Author author={user} disableLink={true} size="s" />
                </li>
              ))}
            </ul>
          )}

          {searchTerm && filteredFollowers.length === 0 && (
            <p className="text-sm text-gray-500 mt-1">No users found.</p>
          )}
        </div>
      )}

        </form>
      )
    }
    
  