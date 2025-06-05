'use client'
import { useState, useEffect } from 'react'
import Author from '../components/Author'
import ImageIcon from '../components/AddImageIcon'
import ImageUploadPreview from '../components/ImageUploadPreview'
import ErrorMessage from '../components/ErrorMessage'
import BackButton from '../components/BackButton'

export default function PostPage() {
    const [post, setPost] = useState(null)
    const [reloadPost, setReloadPost] = useState(false)
    const [commentInput, setCommentInput] = useState('')
    const [commentImage, setCommentImage] = useState(null)
    const [error, setError] = useState(null)
    const [imageError, setImageError] = useState(null)
    
    const [postId, setPostId] = useState(null)

    useEffect(() => {
      const params = new URLSearchParams(window.location.search)
          const id = params.get('post_id')
      if (id !== null && id !== undefined && id !== "") {
        setPostId(id)
      } else {
        setError('No post_id provided in URL')
      }
    }, [])

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
    
          const data = await res.json()
          if (res.ok) {
            console.log('Response is OK')
            console.log('Fetched post:', data)
            console.log('Fetched comments:', data.comments)
            setPost(data)
            setError(null);
          } else {
            console.log('Failed to load post:', data.message)
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
            console.log('Failed to post comment:', errorText)
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
        <div className='m-4'>
          <BackButton href="/feed" className="mb-4" />
          <div className="p-4 rounded mb-6 bg-white rounded shadow border border-gray-200">
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
            <p className="mb-4 break-words">{post.post_content}</p>
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

              <form onSubmit={handleCommentSubmit} className="relative max-w-full mx-0 mt-1 p-4 bg-white rounded shadow border border-gray-200">
                <div className="block">
                <textarea
                  value={commentInput}
                  onChange={(e) => setCommentInput(e.target.value)}
                  placeholder="Write your comment..."
                  className="mt-1 block w-full border border-gray-300 rounded p-2"
                  rows="3"
                  maxLength={200}
                  required
                />

                {imageError && (
                  <p className="text-red-600 text-sm mb-2">{imageError}</p>
                )}

                  <ImageUploadPreview imageFile={commentImage} setImageFile={setCommentImage} />

                  <div className="flex items-center justify-between mt-4">
                    <label className="flex items-center cursor-pointer">
                      <input
                        type="file"
                        accept="image/*"
                        onChange={(e) => {
                          const file = e.target.files[0]
                          if (file && file.size > 1 * 1024 * 1024) {
                            setImageError("This image is too big. Has to be under 1MB")
                            setCommentImage(null)
                          } else {
                            setImageError(null)
                            setCommentImage(file)
                          }
                        }}
                        className="hidden"
                      />
                      <ImageIcon />
                    </label>
                    <button
                      type="submit"
                      className="bg-sky-600/60 hover:bg-sky-700/60 text-white font-bold rounded text-md p-2 px-4"
                    >
                      Submit
                    </button>
                  </div>
                </div>
              </form>

              <div className="mt-6">

                {post.comments && post.comments.length > 0 ? (
                post.comments.map((comment, i) => (
                <div key={i} className="relative max-w-full my-4 mt-1 p-4 bg-white rounded shadow border border-gray-200">
                  
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
