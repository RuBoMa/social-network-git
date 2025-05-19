'use client'
import { useState } from 'react'
import { useSearchParams } from 'next/navigation'
import { useEffect } from 'react'
import ImageIcon from '../components/AddImageIcon'
import ImageUploadPreview from '../components/ImageUploadPreview'
import FollowersModal from './followersModal'

export default function CreatePost({ onSuccess }) {
  const searchParams = useSearchParams()
  const groupID = searchParams.get('group_id')
  
  const [postTitle, setPostTitle] = useState('')
  const [postContent, setContent] = useState('')
  const [privacy, setPostPrivacy] = useState(groupID ? 'followers' : 'public');
  const [postImage, setPostImage] = useState(null)
  const [followers, setFollowers] = useState([])
  const [selectedUsers, setSelectedUsers] = useState([])
  const [showModal, setShowModal] = useState(false)

  async function handlePost(e) {
    e.preventDefault()

    const formData = new FormData()
    formData.append('post_title', postTitle)
    formData.append('post_content', postContent)
    formData.append('privacy', privacy)
    if (postImage) formData.append('post_image', postImage)
    if (groupID) formData.append('group_id', groupID)
    if (privacy === 'followers') formData.append('followers', JSON.stringify(followers))
    if (privacy === 'custom') formData.append('custom_users', JSON.stringify(selectedUsers))

    const res = await fetch('http://localhost:8080/api/create-post', {
      method: 'POST',
      credentials: 'include',
      body: formData,
    })

    if (res.ok) {
      setContent('')
      setPostTitle('')
      setPostImage(null)
      onSuccess && onSuccess() // trigger reload if provided
    } else {
      alert('Failed to post')
    }
  }

  useEffect(() => {
    if (privacy === 'custom') {
      fetch('http://localhost:8080/api/followers', {
        credentials: 'include'
      })
        .then(res => res.json())
        .then(data => {
          console.log('Fetched followers:', data)
          setFollowers(data)})
        .catch(err => console.error('Failed to fetch followers', err))
    }
  }, [privacy])

  return (
    <form className="relative max-w-full mx-0 mt-1 p-2 bg-white rounded-xl shadow-md" onSubmit={handlePost}>
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
  
      <label className="block mb-4">
        <textarea
          value={postContent}
          onChange={(e) => setContent(e.target.value)}
          placeholder="What's on your mind?"
          required
          className="mt-1 block w-full border border-gray-300 rounded p-2"
        />
      </label>
  
      <div className="flex items-center justify-between gap-6 mb-4">
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
  
        {!groupID && (
          <div className="relative inline-block">
            <div className="inline-flex items-center gap-4 text-sm">
            {["public", "followers", "custom"].map((option) => (
              <label key={option} className="inline-flex items-center">
              <input
                type="radio"
                name="privacy"
                value={option}
                checked={privacy === option}
                onChange={(e) => {
                  setPostPrivacy(e.target.value)
                  if (e.target.value === 'custom') {
                    setShowModal(true)  // opening the modal
                }
              }}
              className="form-radio text-blue-600"
            />
              <span className="ml-1 capitalize">{option}</span>
            </label>
            ))}
            </div>
          </div>
        )}
      </div>
  
      {/* Modal popup */}
      {showModal && (
        <FollowersModal
          followers={followers}
          selectedUsers={selectedUsers}
          setSelectedUsers={setSelectedUsers}
          onClose={() => setShowModal(false)}
        />
      )}
  
      <button
        type="submit"
        className="w-auto bg-blue-600 text-sm text-white mx-2 py-0.5 px-3 rounded hover:bg-blue-700 transition whitespace-nowrap"
      >
        Submit
      </button>
    </form>
  )
}