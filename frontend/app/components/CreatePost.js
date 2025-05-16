'use client'
import { useState } from 'react'
import { useSearchParams } from 'next/navigation'

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
    <form className="max-w-full mx-0 mt-1 p-2 bg-white rounded-xl shadow-md" onSubmit={handlePost}>
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
          <svg xmlns="http://www.w3.org/2000/svg" fill="none"
            viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor"
            className="w-6 h-6 text-black-600 hover:text-blue-800">
            <path strokeLinecap="round" strokeLinejoin="round"
              d="m2.25 15.75 5.159-5.159a2.25 2.25 0 0 1 3.182 0l5.159 5.159m-1.5-1.5
                 1.409-1.409a2.25 2.25 0 0 1 3.182 0l2.909 2.909m-18 3.75h16.5a1.5 1.5
                 0 0 0 1.5-1.5V6a1.5 1.5 0 0 0-1.5-1.5H3.75A1.5 1.5 0 0 0 2.25 6v12a1.5
                 1.5 0 0 0 1.5 1.5Zm10.5-11.25h.008v.008h-.008V8.25Zm.375 0a.375.375
                 0 1 1-.75 0 .375.375 0 0 1 .75 0Z" />
          </svg>
        </label>

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
          className="w-auto bg-blue-600 text-sm text-white mx-2 py-0.5 px-3 rounded hover:bg-blue-700 transition whitespace-nowrap"
        >
          Submit
        </button>
      </div>
    </form>
  )
}
