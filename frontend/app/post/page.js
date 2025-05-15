'use client'
import { useState, useEffect } from 'react'
import { useSearchParams } from 'next/navigation'
import Link from 'next/link'

export default function PostPage() {
    const searchParams = useSearchParams()
    const postId = searchParams.get('post_id') // this is your query param
    const [post, setPost] = useState(null)
    const [reloadPost, setReloadPost] = useState(false)
    const [commentInput, setCommentInput] = useState('')
    
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
    
    
          if (res.ok) {
            console.log('Response is OK') // Log if the response is OK
            const data = await res.json()
            console.log('Fetched post:', data) // Log the fetched posts
            setPost(data)
          } else {
            console.error('Failed to load posts')
          }
        }
    
        fetchPost()
      }, [postId, reloadPost])

      const handleCommentSubmit = async (e) => {

        e.preventDefault()
      
        if (!commentInput.trim()) return
      
        try {
          const res = await fetch('http://localhost:8080/api/comment', {
            method: 'POST',
            credentials: 'include',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({
              post_id: Number(postId), // make sure this is defined
              comment_content: commentInput,
            }),
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
                src={post.post_image}
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

