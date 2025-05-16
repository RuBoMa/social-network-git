'use client'
import { useState, useEffect } from 'react'
import { useSearchParams } from 'next/navigation'
import CreatePost from '../components/CreatePost'
import { PostFeed } from '../components/PostFeed'

export default function PostPage() {
    const searchParams = useSearchParams()
    const groupId = searchParams.get('group_id') // this is your query param
    const [group, setGroup] = useState(null)
    const [reloadPosts, setReloadPosts] = useState(false)
    
      useEffect(() => {
        async function fetchGroup() {
          if (!groupId) return

          const res = await fetch(`http://localhost:8080/api/group?group_id=${groupId}`, {
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
            console.log('Fetched group:', data) // Log the fetched group
            setGroup(data)
          } else {
            console.error('Failed to load group')
          }
        }
    
        fetchGroup()
      }, [groupId])

      if (!group) {
        return <div>Loading group...</div>
      }

        return (
        <div className="flex flex-col items-center p-4">
            <h1 className="text-2xl font-bold mb-2">{group.group_name}</h1>
            <p className="text-gray-700 mb-2 italic">{group.group_desc}</p>

            <div className="text-sm text-gray-600 mb-4">
            <p>
                Created by <span className="font-medium">{group.group_creator.user_id}</span> on{" "}
                {new Date(group.group_created_at).toLocaleDateString()}
            </p>
            <p>Members: {group.group_members?.length}</p>
            </div>

            { /* IF GROUP MEMBER */ }
           {group.is_member ? (
        <div className="w-full">
        <CreatePost onSuccess={() => setReloadPosts(prev => !prev)} />
        <h2 className="text-xl font-semibold my-4">Group Posts</h2>
        <PostFeed reloadTrigger={reloadPosts} />
        </div>
      ) : (
        <p className="text-red-500 font-semibold">Join the group to see all posts.</p>
      )}
    </div>
  )
}