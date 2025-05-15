'use client'
import { useState, useEffect } from 'react'
import { useSearchParams } from 'next/navigation'
import { useRouter } from 'next/navigation'

export default function PostPage() {
    const searchParams = useSearchParams()
    const router = useRouter()
    const postId = searchParams.get('post_id') // this is your query param
    const [post, setPost] = useState(null)
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
      }, [postId])

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
                <h3>{post.post_title}</h3>
                <div className="flex items-center gap-2 mb-2">
                {post.author?.avatar_path && (
                 <img
                src= {`http://localhost:8080${post.author.avatar_path}`}
                alt="Author"
                className="w-10 h-10 rounded-full"
                />
                )}
                <span className="font-medium">
                {post.author?.nickname || post.author?.first_name || "Unknown Author"}
                </span>
                </div>
                <p>{post.post_content}</p>
                {post.post_image && (
                  <img src={post.post_image} alt="Post visual" style={{ maxWidth: '100%' }} />
                )}
                <p><small>{new Date(post.created_at).toLocaleString()}</small></p>
                <div className="mt-6">
                <h4 className="text-md font-semibold mb-2">Comments</h4>
                {post.comments && post.comments.length > 0 ? (
        post.comments.map((comment, i) => (
          <div key={i} className="mb-4 p-3 border rounded bg-gray-50">
            <div className="flex items-center mb-2">
              {comment.comment_author?.avatar_path && (
                <img
                  src={`http://localhost:8080${comment.author.avatar_path}`}
                  alt="Author"
                  className="w-8 h-8 rounded-full mr-2"
                />
              )}
              <span className="font-semibold">{comment.comment_content?.nickname || comment.comment_author?.first_name}</span>
            </div>
            <p>{comment.comment_content}</p>
            <p className="text-xs text-gray-500">
              {new Date(comment.created_at).toLocaleString()}
            </p>
          </div>
        ))
      ) : (
        <p className="text-sm text-gray-600">No comments yet.</p>
      )}
    </div>

    {/* Add Comment Form */}
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
  </div>
)
}

