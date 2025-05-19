import { useState, useEffect } from 'react';

export default function TagInput({ followers, selectedUsers, setSelectedUsers }) {
  const [search, setSearch] = useState('');
  const [filteredFollowers, setFilteredFollowers] = useState([]);

  useEffect(() => {
    if (search.trim() === '') {
      setFilteredFollowers([]);
    } else {
      const lowerSearch = search.toLowerCase();
      setFilteredFollowers(
        followers.filter(
          (f) =>
            (f.nickname?.toLowerCase().includes(lowerSearch) ||
             f.first_name?.toLowerCase().includes(lowerSearch)) &&
            !selectedUsers.includes(f.user_id)
        )
      );
    }
  }, [search, followers, selectedUsers]);

  const addUser = (userId) => {
    setSelectedUsers([...selectedUsers, userId]);
    setSearch(''); // empty the search input after adding a user
  };

  const removeUser = (userId) => {
    setSelectedUsers(selectedUsers.filter((id) => id !== userId));
  };

  return (
    <div className="w-full max-w-md">
      {/* used followers as tags */}
      <div className="flex flex-wrap gap-2 mb-2">
        {selectedUsers.map((userId) => {
          const user = followers.find((f) => f.user_id === userId);
          if (!user) return null;
          const name = user.nickname || user.first_name || 'Unnamed';

          return (
            <div
              key={userId}
              className="flex items-center bg-blue-200 text-blue-900 rounded px-2 py-1 text-sm"
            >
              <span>{name}</span>
              <button
                type="button"
                className="ml-1 font-bold focus:outline-none"
                onClick={() => removeUser(userId)}
              >
                &times;
              </button>
            </div>
          );
        })}
      </div>

      {/* search */}
      <input
        type="text"
        value={search}
        onChange={(e) => setSearch(e.target.value)}
        placeholder="Search followers..."
        className="w-full border rounded px-2 py-1"
      />

      {/* search results */}
      {filteredFollowers.length > 0 && (
        <div className="border rounded mt-1 max-h-40 overflow-auto bg-white shadow-lg z-10 relative">
          {filteredFollowers.map((follower) => {
            const name = follower.nickname || follower.first_name || 'Unnamed';
            return (
              <div
                key={follower.user_id}
                onClick={() => addUser(follower.user_id)}
                className="cursor-pointer px-2 py-1 hover:bg-blue-100"
              >
                {name}
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}
