'use client'
import { useState, useEffect } from 'react'
import { useSearchParams } from 'next/navigation'
import Link from 'next/link'
import Author from '../components/Author'
import ImageIcon from '../components/AddImageIcon'
import ImageUploadPreview from '../components/ImageUploadPreview'
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
            console.log('Response is OK') // Log if the response is OK
            console.log('Fetched post:', data) // Log the fetched posts
            console.log('Fetched comments:', data.comments) // Log the fetched comments
            setPost(data)
          } else {
            console.error('Failed to load posts')
            setError(data.message || 'Failed to load posts')

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

            setCommentInput('') // clear input
            setCommentImage(null) // clear image
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
      return <ErrorMessage message={error} />
    }
    // Don't try to render until post is loaded
    if (!post) {
      return <p>Loading post...</p>
    }

      return (
        <div>
          {/* div for the whole post */}
          <div className="p-4 rounded mb-6 shadow-md">
          {/* div for the post author and timestamp */}
          <div className="flex items-center justify-between mb-2">
            {/* Author */}
            <Author author={post.author} size="lg" />
            {/* Timestamp */}
            <p className="text-xs text-gray-500 mb-2">
              {new Date(post.created_at).toLocaleString('en-GB', {
                day: '2-digit',
                month: '2-digit',
                year: 'numeric',
                hour: '2-digit',
                minute: '2-digit',
              })}
            </p>
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

              {/* Comments Section */}
              <div className="mt-6"> {/* <- this creates the gap */}

                {post.comments && post.comments.length > 0 ? (
                post.comments.map((comment, i) => (
                <div key={i} className="mb-4 p-3 rounded bg-gray-50">
                  
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

                  {/* Image if exists */}
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