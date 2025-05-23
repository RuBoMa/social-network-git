'use client'
import Author from '../Author'
import InviteResponseButton from './ResponseButton'

export default function GroupRequest({ requests, groupId, onResponse }) {
  return (
    <div className="w-full border-b border-gray-300 pb-4 pt-4">
      {requests.length > 0 ? (
        <ul className="max-h-34 overflow-y-auto border border-gray-200 rounded p-2 space-y-2 shadow">
          {users.map(user => (
            <li
              key={sender.user_id}
              className="flex justify-between items-center border-b border-gray-200 pb-2 last:border-b-0 last:pb-0"
            >
              <Author author={sender} size="sm" />
              <div className="flex space-x-2">
                <InviteResponseButton
                  groupId={groupId}
                  requestId={requests.request_id}
                  status="accepted"
                  onResponse={onResponse}
                />
                <InviteResponseButton
                  groupId={groupId}
                  requestId={requests.request_id}
                  status="rejected"
                  onResponse={onResponse}
                />
              </div>
            </li>
          ))}
        </ul>
      ) : (
        <p>No requests available</p>
      )}
    </div>
  )
}