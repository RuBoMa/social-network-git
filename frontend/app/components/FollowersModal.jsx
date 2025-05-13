'use client'
export default function FollowersModal({ followers }) {
  return (
    <div className="absolute top-full left-0 mt-2 w-56 bg-white border rounded-lg shadow-lg max-h-64 overflow-y-auto z-50">
      <ul className="p-2 space-y-2">
        {followers.length > 0 ? (
          followers.map(f => (
            <li key={f.id} className="flex items-center space-x-3">
              {f.avatar_url && (
                <img
                  src={f.avatar_url}
                  alt="avatar"
                  className="w-6 h-6 rounded-full object-cover"
                />
              )}
              <span>{f.nickname || `${f.first_name} ${f.last_name}`}</span>
            </li>
          ))
        ) : (
          <li className="text-gray-500 text-center">No followers yet.</li>
        )}
      </ul>
    </div>
  )
}