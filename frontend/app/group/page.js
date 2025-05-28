'use client'
import { useState, useEffect } from 'react'
import { useSearchParams } from 'next/navigation'
import CreatePost from '../components/CreatePost'
import CreateEvent from '../components/Group/CreateEvent'
import { PostFeed } from '../components/PostFeed'
import JoinGroupButton from '../components/Group/JoinGroupButton'
import GroupInvitation from '../components/Group/GroupInvitation'
import ErrorMessage from '../components/ErrorMessage'
import ResponseButton from '../components/Group/ResponseButton'
import Link from 'next/link'
import Author from '../components/Author'
import GroupRequest from '../components/Group/GroupRequest'
import ChatWindow from '../components/ChatWindow'
import GroupAvatar from '../components/GroupAvatar'

export default function GroupPage() {
  const searchParams = useSearchParams();
  const groupId = searchParams.get('group_id');

  const [group, setGroup]         = useState(null);
  const [reloadPosts, setReloadPosts]   = useState(false);
  const [reloadGroup, setReloadGroup]   = useState(false);
  const [showEventForm, setShowEventForm] = useState(false);
  const [showGroupChat, setShowGroupChat] = useState(false);
  const [showMembers, setShowMembers] = useState(false);

  useEffect(() => {
    async function fetchGroup() {
      if (!groupId) return;

      const res = await fetch(
        `http://localhost:8080/api/group?group_id=${groupId}`,
        { credentials: 'include', method: 'GET', headers: { Accept: 'application/json' } }
      );
      const data = await res.json();
      if (res.ok) {
        setGroup(data);
        console.log('Fetched group:', data);
      } else {
        ErrorMessage(data.message || 'Failed to load group');
      }
    }

    fetchGroup();
  }, [groupId, reloadGroup]);

  if (!group) return <div>Loading group...</div>;

  return (
    <div className="flex flex-col items-center p-4 w-full overflow-x-hidden">
      <div className="w-full bg-white pb-4 border-b border-gray-300">
        <div className="flex justify-between">
            <GroupAvatar group={group} disableLink={true} size="xl" />
             <button
            className="text-sm text-black-600 focus:outline-none hover:underline cursor-pointer transition"
            onClick={() => setShowMembers(true)}
            title="Show group members"
          >
            Members: {group.group_members?.length || 0}
          </button>
        </div>
        <p className="text-gray-700 p-2 mb-2 italic">{group.group_desc}</p>
        <div className="flex justify-between text-sm text-gray-600 ">
          <div className="flex items-center gap-1">
            Created by <Author author={group.group_creator} />
          </div>
          <div className="flex items-center gap-1">
            Created at{' '}
            <span>{new Date(group.group_created_at).toLocaleDateString()}</span>
          </div>
        </div>
      </div>

        {showMembers && (
        <div
          className="fixed inset-0 flex items-center justify-center z-50"
          style={{ background: 'rgba(0,0,0,0.2)' }}
          onClick={() => setShowMembers(false)}
        >
          <div
            className="bg-white rounded-lg shadow-lg p-6 max-w-sm w-full relative"
            onClick={e => e.stopPropagation()}
          >
            <button
              className="absolute top-2 right-2 text-gray-500 hover:text-gray-700"
              onClick={() => setShowMembers(false)}
            >
              ✕
            </button>
            <h4 className="text-lg font-semibold mb-4">Group Members</h4>
            {group.group_members && group.group_members.length > 0 ? (
              <ul>
                {group.group_members.map(member => (
                  <li key={member.user_id} className="mb-2">
                    <span className="flex items-center space-x-2">
                      <Author author={member} size="sm" />
                    </span>
                  </li>
                ))}
              </ul>
            ) : (
              <p>No members yet.</p>
            )}
          </div>
        </div>
      )}

      {group.is_member ? (
      <div className="w-full">
        {group.group_requests && group.group_requests.length > 0 && (
          <div className="mt-4">
            <h3 className="text-mg font-semibold">Join requests</h3>
            <GroupRequest
              requests={group.group_requests}
              groupId={group.group_id}
              onResponse={() => setReloadGroup(prev => !prev)}
            />
          </div>
        )}
        <div className="mt-4">
        <h3 className="text-mg font-semibold">Invite Users to Join</h3>
        <GroupInvitation groupId={group.group_id} />
        </div>
        <div className="mt-4">
          <div className="flex justify-end mb-2">
            <button
              onClick={() => setShowEventForm(prev => !prev)}
              className="text-blue-500 cursor-pointer text-sm"
              type="button"
            >
              {showEventForm ? 'Create Post +' : 'Create Event +'}
            </button>
          </div>

          {showEventForm ? (
            <CreateEvent
              onClose={() => setShowEventForm(false)}
              onSuccess={() => {
                setShowEventForm(false);
                setReloadGroup(prev => !prev);
              }}
            />
          ) : (
            <CreatePost onSuccess={() => setReloadPosts(prev => !prev)} />
          )}
        </div>
              {/* Group Chat Section */}
           {!group.has_chat && (
              <button
                className="bg-blue-600 text-white px-4 py-2 rounded mb-2"
                onClick={() => setShowGroupChat(true)}
              >
                Open Group Chat
              </button>
            )}
            {showGroupChat && (
              <ChatWindow
                group={group}
                onClose={() => setShowGroupChat(false)}
                isGroupChat={true}
              />
            )}

        <div className="my-6 w-full">
          <h2 className="text-xl font-semibold mb-3 text-gray-800">Upcoming Events</h2>
          {(!group.group_events || group.group_events.length === 0) ? (
              <p className="text-gray-500">No upcoming events for this group.</p>
          ) : (
            <div className="overflow-x-auto w-full">
              <div className="flex space-x-4 pb-4">
                {group.group_events.map(event => (
                  <Link
                    key={event.event_id}
                    href={`/event?event_id=${event.event_id}`}
                    className="w-[280px] flex-shrink-0 bg-white rounded shadow hover:bg-gray-50 transition-colors duration-200 ease-in-out border border-gray-200 overflow-hidden"
                  >
                    <div className="p-4">
                      <h3
                        className="text-lg font-bold text-blue-600 mb-2 truncate"
                        title={event.title}
                      >
                        {event.title}
                      </h3>
                      <p className="text-sm text-gray-500 mb-1">
                        {new Date(event.event_time || event.event_date)
                          .toLocaleString([], {
                            year: 'numeric',
                            month: 'long',
                            day: 'numeric',
                            hour: '2-digit',
                            minute: '2-digit',
                            hour12: false,
                            timeZone: 'UTC'
                          })}
                      </p>
                      <p
                        className="text-sm text-gray-700 line-clamp-3"
                        title={event.description}
                      >
                        {event.description}
                      </p>
                    </div>
                  </Link>
                ))}
              </div>
            </div>
          )}
        </div>


        <h2 className="text-xl font-semibold mb-3 ">Group Posts</h2>
        <PostFeed reloadTrigger={reloadPosts} />
        </div>
      ) : (
        <div className="mb-4">
          {group.request_status === "" && (
            <JoinGroupButton
              groupId={group.group_id}
              onJoin={() =>
                setGroup(prev => ({ ...prev, request_status: 'requested' }))
              }
            />
          )}
          {group.request_status === 'requested' && (
            <p className="text-yellow-500 font-semibold">
              Request sent, waiting for approval
            </p>
          )}
          {group.request_status === 'invited' && (
            <>
              <ResponseButton
                groupId={group.group_id}
                requestId={group.request_id}
                status="accepted"
                onResponse={() => setReloadGroup(true)}
              />
              <ResponseButton
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
  );
}