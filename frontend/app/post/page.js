'use client'
import { useState, useEffect } from 'react'
import { useSearchParams } from 'next/navigation'
import Link from 'next/link'
import ErrorMessage from '../components/ErrorMessage'

export default function PostPage() {
    const searchParams = useSearchParams()
    const postId = searchParams.get('post_id') // this is your query param
    const [post, setPost] = useState(null)
    const [reloadPost, setReloadPost] = useState(false)
    const [commentInput, setCommentInput] = useState('')
    const [commentImage, setCommentImage] = useState(null)

    const [error, setError] = useState(null)

    useEffect(() => {
        async function fetchPost() {
            if (!postId) return

            const res = await fetch(`http://localhost:8080/api/post?post_id=${postId}`, {
                credentials: 'include',
                method: 'GET',
                headers: {
              'Accept': 'application/json' //telling the server we want JSON
            }
          })
          console.log('Response status:', res) // Log the response status
    
    
          const data = await res.json()
          if (res.ok) {
            console.log('Fetched post:', data) // Log the fetched posts
            setPost(data)
          } else {
              setError(data.message || 'Failed to load post')
              setPost(null)
          } 
        }
    
        fetchPost()
      }, [postId, reloadPost])

      const handleCommentSubmit = async (e) => {

        e.preventDefault()
      
        if (!commentInput.trim()) return
        console.log(postId)

        const formData = new FormData()
        formData.append('post_id', postId)
        formData.append('comment_content', commentInput)
        if (commentImage) {
          formData.append('comment_image', commentImage)
        }

      
        try {
          const res = await fetch('http://localhost:8080/api/comment', {
            method: 'POST',
            credentials: 'include',
            body: formData,
          })
          if (res.ok) {
            const response = await res.json()
            console.log('Comment posted:', response)
            setCommentInput('') // clear input
            setReloadPost(prev => !prev) // trigger re-fetch
          } else {
            console.error('Failed to post comment')
          }
        } catch (error) {
          console.error('Error submitting comment:', error)
        }
    }
      
    if (error) {
      return <ErrorMessage message={error} />
    }
    // Don't try to render until post is loaded
    if (!post) {
      return <p>Loading post...</p>
    }

      return (
        <div>
          {/* Post */}
          <div className="p-4 border rounded mb-6">

            {/* Timestamp */}
            <p className="text-sm text-gray-500 mb-2">
              {new Date(post.created_at).toLocaleString('en-GB', {
                day: '2-digit',
                month: '2-digit',
                year: 'numeric',
                hour: '2-digit',
                minute: '2-digit',
              })}
            </p>

            {/* Author */}
            <div className="flex items-center gap-2 mb-2">
              {post.author?.avatar_path && (
                <img
                  src={`http://localhost:8080${post.author.avatar_path}`}
                  alt="Author"
                  className="w-10 h-10 rounded-full"
                />
              )}
              <Link
              href={`/profile?user_id=${post.author?.user_id}`}
              className="font-semibold text-blue-600 hover:underline"
              >
                {post.author?.nickname ||
                  `${post.author?.first_name || 'Unknown'} ${post.author?.last_name || ''}`}
              </Link>
            </div>

            {/* Title */}
            <h3 className="text-lg font-bold mb-2">{post.post_title}</h3>

            {/* Content */}
            <p className="mb-4">{post.post_content}</p>

            {/* Image */}
            {post.post_image && (
              <img
                src={`http://localhost:8080${post.post_image}`}
                alt="Post visual"
                className="max-w-full rounded"
              />
            )}
          </div>

          {/* Comment Section */}
          <div className="mt-6">
            <h4 className="text-md font-semibold mb-2">Comments</h4>

              {/* Comment Form */}
              <form
                onSubmit={handleCommentSubmit}
                className="mt-6 p-4 border rounded"
              >
                <h4 className="text-md font-semibold mb-2">Add a Comment</h4>
                <textarea
                  value={commentInput}
                  onChange={(e) => setCommentInput(e.target.value)}
                  placeholder="Write your comment..."
                  className="w-full p-2 border rounded mb-2"
                  rows="3"
                  required
                />
                {/* Comment Image */}      
                <div className="flex items-center justify-between gap-6 mb-4">
                  <label className="inline-flex items-center space-x-2 cursor-pointer">
                    <input
                      type="file"
                      accept="image/*"
                      onChange={(e) => setCommentImage(e.target.files[0])}
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
                </div>
                <button
                  type="submit"
                  className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600"
                >
                  Post Comment
                </button>
              </form>

              {/* Comments Section */}
              <div className="mt-6"> {/* <- this creates the gap */}

                {post.comments && post.comments.length > 0 ? (
                post.comments.map((comment, i) => (
                <div key={i} className="mb-4 p-3 border rounded bg-gray-50">
                  
                  {/* Author */}
                  <div className="flex items-center mb-2">
                  {comment.comment_author?.avatar_path && (
                    <img
                      src={`http://localhost:8080${comment.comment_author.avatar_path}`}
                      alt="Author"
                      className="w-8 h-8 rounded-full mr-2"
                    />
                  )}
                  <Link
                  href={`/profile?user_id=${comment.comment_author.user_id}`} 
                  className="font-semibold text-blue-600 hover:underline"
                  >
                  {comment.comment_author?.nickname || 
                  `${comment.comment_author?.first_name || 'Unknown'} ${comment.comment_author?.last_name || ''}`}
                  </Link>
                  </div>

                  {/* Timestamp */}
                  <p className="text-sm text-gray-500 mb-2">
                  {new Date(comment.created_at).toLocaleString('en-GB', {
                    day: '2-digit',
                    month: '2-digit',
                    year: 'numeric',
                    hour: '2-digit',
                    minute: '2-digit',
                  })}
                  </p>

                  {/* Content */}
                  <p>{comment.comment_content}</p>

                </div>
                ))
                ) : (
                  <p className="text-sm text-gray-600">No comments yet.</p>
                )}
              </div>
          </div>
        </div>
      )
}

