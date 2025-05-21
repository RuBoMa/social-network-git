'use client'
import { useState } from 'react'
import { useSearchParams } from 'next/navigation'
import ImageIcon from '../components/AddImageIcon'
import ImageUploadPreview from '../components/ImageUploadPreview'

export default function CreatePost({ onSuccess }) {
  const searchParams = useSearchParams()
  const groupID = searchParams.get('group_id')
  
  const [postTitle, setPostTitle] = useState('')
  const [postContent, setContent] = useState('')
  const [privacy, setPostPrivacy] = useState(groupID ? 'followers' : 'public');
  const [postImage, setPostImage] = useState(null)

  async function handlePost(e) {
    e.preventDefault()

    const formData = new FormData()
    formData.append('post_title', postTitle)
    formData.append('post_content', postContent)
    formData.append('privacy', privacy)
    if (postImage) formData.append('post_image', postImage)
    if (groupID) formData.append('group_id', groupID)

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

  return (
    <form className="max-w-full mx-0 mt-1 p-2 p-4 rounded shadow border-gray-200 border" onSubmit={handlePost}>
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

      <div className="flex items-center justify-between gap-6">
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
          <div>
        <div className="flex gap-4 text-sm">
          {["public", "followers", "custom"].map((option) => (
            <label key={option} className="inline-flex items-center">
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
        </div>
        )}
        <button
          type="submit"
          className="bg-blue-500 hover:bg-blue-600 text-white px-2 py-1 rounded text-sm"
        >
          Submit
        </button>
      </div>
    </form>
  )
}
