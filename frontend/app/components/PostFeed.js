'use client'

import { useEffect, useState } from 'react'
import { useSearchParams } from 'next/navigation'
import Link from 'next/link'
import Author from '../components/Author'

export function PostFeed({ reloadTrigger }) {
  const searchParams = useSearchParams()
  const groupID = searchParams?.get('group_id')

  const [posts, setPosts] = useState(null)

  useEffect(() => {
    async function fetchPosts() {
      const url = groupID
        ? `http://localhost:8080/api/feed?group_id=${groupID}`
        : 'http://localhost:8080/api/feed'
      console.log('Fetching posts from:', url)
      const res = await fetch(url, {
        credentials: 'include',
        method: 'GET',
      })

      console.log('Response status:', res) // Log the response status

      if (res.ok) {
        const data = await res.json()
        console.log(data)
        setPosts(data)
      } else {
        console.error('Failed to load posts')
      }
    }

    fetchPosts()
  }, [groupID, reloadTrigger])

    // Don't try to render until posts are loaded
  if (posts === undefined) {
    return <p>Loading feed...</p>
  }

  return (
    <div>
      {Array.isArray(posts) && posts.length > 0 ? (
        posts.map((post, i) => (
                <Link
              key={post.post_id}
              href={`/post?post_id=${post.post_id}`}
              className="block mb-4 p-4 rounded shadow border border-gray-200 bg-white hover:bg-gray-100 transition cursor-pointer"
            >
              <div className="flex justify-between items-center mb-2">
                <Author author={post.author} size="s" />
                <div className="text-right">
                  <p className="text-xs text-gray-500">
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
                  <p className="text-xs text-gray-500">
                    {Array.isArray(post.comments)
                      ? `${post.comments.length} comment${post.comments.length === 1 ? '' : 's'}`
                      : '0 comments'}
                  </p>
                </div>
              </div>
              <h3 className="text-lg font-semibold break-words text-gray-700 pt-3">
                {post.post_title}
              </h3>
              <p className="text-gray-700">
                {post.post_content.length > 50
                  ? post.post_content.slice(0, 50) + '...'
                  : post.post_content}
              </p>
            </Link>
        ))
      ) : (
          <p className="text-gray-500">No posts to show.</p>
      )}
    </div>
  )
}
