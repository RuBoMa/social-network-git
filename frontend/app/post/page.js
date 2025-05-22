'use client'
import { useState, useEffect } from 'react'
import { useSearchParams } from 'next/navigation'
import Author from '../components/Author'
import ImageIcon from '../components/AddImageIcon'
import ImageUploadPreview from '../components/ImageUploadPreview'
import ErrorMessage from '../components/ErrorMessage'

export default function PostPage() {
    const searchParams = useSearchParams()
    const postId = searchParams.get('post_id')
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
              'Accept': 'application/json'
            }
          })
          console.log('Response status:', res)
    
          const data = await res.json()
          if (res.ok) {
            console.log('Response is OK')
            console.log('Fetched post:', data)
            console.log('Fetched comments:', data.comments)
            setPost(data)
            setError(null);
          } else {
            console.error('Failed to load post:', data.message)
            setError(data.message || 'Failed to load post')
          }
        }
    
        fetchPost()
      }, [postId, reloadPost])

      const handleCommentSubmit = async (e) => {

        e.preventDefault()
      
        if (!commentInput.trim()) return
      
        try {
          const formData = new FormData()
          formData.append('post_id', Number(postId))
          formData.append('comment_content', commentInput)
          if (commentImage) formData.append('comment_image', commentImage)

            
          const res = await fetch('http://localhost:8080/api/comment', {
            method: 'POST',
            credentials: 'include',
            body: formData,
          })
          
          if (res.ok) {
            const bodyText = await res.text()
            console.log('Comment posted:', bodyText)

            const response = JSON.parse(bodyText)
            console.log('Comment posted:', response)

            setCommentInput('')
            setCommentImage(null)
            setReloadPost(prev => !prev) // trigger re-fetch
          } else {
            const errorText = await res.text()
            console.error('Failed to post comment:', errorText)
          }
        } catch (error) {
          console.error('Error submitting comment:', error)
        }
    }
    
    if (error) {
      return (
        <div className="p-4">
          <ErrorMessage message={error} />
        </div>
      )
    }
    
    // Don't try to render until post is loaded
    if (!post) {
      return <p>Loading post...</p>
    }

      return (
        <div>
          <div className="p-4 rounded mb-6 shadow-md">
          <div className="flex items-center justify-between mb-2">
            <Author author={post.author} size="lg" />
            <p className="text-xs text-gray-500 mb-2">
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
            </div>

            <h3 className="text-lg font-bold mb-2">{post.post_title}</h3>
            <p className="mb-4">{post.post_content}</p>
            {post.post_image && (
              <img
                src={`http://localhost:8080${post.post_image}`}
                alt="Post visual"
                className="max-w-full rounded"
              />
            )}
          </div>

          <div className="mt-6">
            <h4 className="text-md font-semibold mb-2">Comments</h4>

              <form onSubmit={handleCommentSubmit} className="mt-6 p-4 rounded">
                <div className="relative w-full">
                <textarea
                  value={commentInput}
                  onChange={(e) => setCommentInput(e.target.value)}
                  placeholder="Write your comment..."
                  className="w-full p-2 border rounded mb-2"
                  rows="3"
                  required
                />

                <label className="absolute bottom-4 right-2 inline-flex items-center space-x-2 cursor-pointer">
                  <input
                    type="file"
                    accept="image/*"
                    onChange={(e) => setCommentImage(e.target.files[0])}
                    className="hidden"
                />
                  <ImageIcon />
                </label>
                </div>

                  <ImageUploadPreview imageFile={commentImage} setImageFile={setCommentImage} />
                  
                <button
                  type="submit"
                  className="bg-blue-500 mt-2 text-white px-3 py-1 rounded hover:bg-blue-600"
                >
                  Post Comment
                </button>
              </form>

              <div className="mt-6">

                {post.comments && post.comments.length > 0 ? (
                post.comments.map((comment, i) => (
                <div key={i} className="mb-4 p-3 rounded bg-gray-50">
                  
                  <div className="flex items-center justify-between mb-2">
                <div className="flex items-center mb-2">
                      <Author author={comment.comment_author} size="xs" />
                </div>

                  <p className="text-sm text-gray-500 mb-2">
                  {new Date(comment.created_at).toLocaleString('en-GB', {
                    day: '2-digit',
                    month: '2-digit',
                    year: 'numeric',
                    hour: '2-digit',
                    minute: '2-digit',
                    hour12: false,
                    timeZone: 'UTC',
                  })}
                  </p>
                </div>

                  <p>{comment.comment_content}</p>

                  {comment.comment_image && (
                    <img
                      src={`http://localhost:8080${comment.comment_image}`}
                      alt="Comment attachment"
                      className="max-w-full rounded mt-2"
                    />
                  )}

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
