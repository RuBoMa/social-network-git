'use client'


export default function FollowersModal({ followers, selectedUsers, setSelectedUsers, onClose }) {
  return (
    <div className="absolute mt-2 bg-white rounded shadow-lg flex flex-col z-50 w-72 max-h-80 overflow-auto p-4">
        <h3 className="text-md font-semibold mb-2">Select followers:</h3>


        {Array.isArray(followers) && followers.length > 0 ? (
          <div className="grid grid-cols-2 gap-2">
            {followers.map((follower) => (
              <label key={follower.user_id} className="flex items-center gap-2">
                <input
                  type="checkbox"
                  checked={selectedUsers.includes(follower.user_id)}
                  onChange={() => {
                    setSelectedUsers((prev) =>
                      prev.includes(follower.user_id)
                        ? prev.filter((id) => id !== follower.user_id)
                        : [...prev, follower.user_id]
                    );
                  }}
                />
                <span>{follower.nickname || follower.first_name || 'Unnamed'}</span>
              </label>
            ))}
          </div>
        ) : (
            <p>No followers found.</p>
        )}

        <button
          onClick={onClose}
          className="absolute top-2 right-2 text-gray-500 hover:text-black text-xl font-bold"
          aria-label="Close"
          >
          x
        </button>
      </div>
  );
}