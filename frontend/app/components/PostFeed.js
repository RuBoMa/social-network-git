import { useEffect, useState } from 'react'
import Link from 'next/link'

export function PostFeed() {
  console.log('PostFeed component rendered')
  const [posts, setPosts] = useState([])

  useEffect(() => {
    async function fetchPosts() {
      const res = await fetch('http://localhost:8080/api/feed', {
        credentials: 'include', 
        method: 'GET',
        headers: {
          'Accept': 'application/json' //telling the server we want JSON
        }
      })
      console.log('Response status:', res) // Log the response status


      if (res.ok) {
        const data = await res.json()
        console.log('Fetched posts:', data) // Log the fetched posts
        setPosts(data)
      } else {
        console.error('Failed to load posts')
      }
    }

    fetchPosts()
  }, [])

  return (
    <div>
      {(posts || []).length > 0 ? 
      (posts.map((post, i) => (
          <div key={i} className="post">
            <h3 className="text-lg font-semibold text-blue-600 hover:underline">
                <Link href={`/post?post_id=${post.post_id}`}>
                {post.post_title}
                </Link>
            </h3>
            <p>{post.post_content}</p>
            {post.post_image && (
              <img src={post.post_image} alt="Post visual" style={{ maxWidth: '100%' }} />
            )}
            <p><small>{new Date(post.created_at).toLocaleString()}</small></p>
          </div>
        ))
      ) : (
        <p>No posts to show.</p>
      )}
    </div>
  )
}
