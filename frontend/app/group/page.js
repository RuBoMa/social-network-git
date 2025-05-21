'use client'
import { useState, useEffect } from 'react'
import { useSearchParams } from 'next/navigation'
import CreatePost from '../components/CreatePost'
import { PostFeed } from '../components/PostFeed'
import JoinGroupButton from '../components/JoinGroupButton'
import GroupInvitation from '../components/Group/GroupInvitation'
import ErrorMessage from '../components/ErrorMessage'
import InviteResponseButton from '../components/Group/InviteResponseButton'

export default function GroupPage() {
    const searchParams = useSearchParams()
    const groupId = searchParams.get('group_id')
    const [group, setGroup] = useState(null)
    const [reloadPosts, setReloadPosts] = useState(false)
    const [reloadGroup, setReloadGroup] = useState(false);

    
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
    
          const data = await res.json()
          if (res.ok) {
            console.log('Fetched group:', data) // Log the fetched group
            setGroup(data)
          } else {
            ErrorMessage(data.message || 'Failed to load group')
          }
        }
    
        fetchGroup()
      }, [groupId, reloadGroup])

      if (!group) {
        return <div>Loading group...</div>
      }

      return (
        <div className="flex flex-col items-center p-4">
          <div className="w-full bg-white p-4 border border-gray-200 rounded shadow">
            <h1 className="text-2xl font-bold mb-2">{group.group_name}</h1>
            <p className="text-gray-700 mb-2 italic">{group.group_desc}</p>

            <div className="text-sm text-gray-600">
            <p>
                Created by <span className="font-medium">{group.group_creator.nickname || `${group.group_creator.first_name} ${group.group_creator.last_name}` }</span>
            </p>
            <p>
                Created at <span>{" "}
                {new Date(group.group_created_at).toLocaleDateString()}
                </span>
            </p>
            <p>Members: {group.group_members?.length}</p>
            </div>
            </div>
           {group.is_member ? (
        <div className="w-full">
            <GroupInvitation groupId={group.group_id} />
            <CreatePost onSuccess={() => setReloadPosts(prev => !prev)} />
            <h2 className="text-xl font-semibold my-4">Group Posts</h2>
            <PostFeed reloadTrigger={reloadPosts} />
            </div>
          ) : (
              <div className="mb-4">
                {group.request_status === "" && (
                  <JoinGroupButton
                    groupId={group.group_id}
                    onJoin={() => setGroup(prev => ({ ...prev, request_status: 'requested' }))}
                  />
                )}
                {group.request_status === 'requested' && (
                  <p className="text-yellow-500 font-semibold">Request sent, waiting for approval</p>
                )}
                {group.request_status === 'invited' && (
                <>
                  <InviteResponseButton
                    groupId={group.group_id}
                    requestId={group.request_id}
                    status="accepted"
                    onResponse={() => setReloadGroup(true)}  // toggle to trigger reload
                  />
                  <InviteResponseButton
                    groupId={group.group_id}
                    requestId={group.request_id}
                    status="rejected"
                    onResponse={() => setReloadGroup(true)}
                  />
                </>
              )}
              </div>
          )}
        </div>
      )
}