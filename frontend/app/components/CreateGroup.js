import { useState } from "react";

export default function CreateGroup({ onClose }) {
    const [name, setName] = useState("");
    const [description, setDescription] = useState("");

    async function handleNewGroup(e) {
    e.preventDefault();
    if (!name.trim() || !description.trim()) return; // simple validation

    const res = await fetch('http://localhost:8080/api/create-group', {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        group_name: name,
        group_desc: description
        }),
    })

    const data = await res.json();
    if (res.ok) {
        console.log('Group created successfully:', data);
        setName('');
        setDescription('');
        const groupId = data.message;
        window.location.href = `/group?group_id=${groupId}`;
    } else {
        console.error('Failed to create group:', data.message);
        ErrorMessage(data.message || 'Failed to create group');
    }

  };

  return (
      <div className="bg-white p-8 w-[90%] max-w-full">
        <h2 className="text-xl font-bold mb-4">Create Group</h2>
        <form onSubmit={handleNewGroup}>
          <input
            type="text"
            className="block w-full mb-3 border p-2 rounded"
            placeholder="Group Name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            required
          />
          <textarea
            className="block w-full mb-3 border p-2 rounded"
            placeholder="Group Description"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            required
          />
          <div className="flex justify-end space-x-2">
            <button
              type="button"
              className="px-4 py-2 border rounded"
              onClick={onClose}
            >
              Cancel
            </button>
            <button
              type="submit"
              className="px-4 py-2 bg-blue-600 text-white rounded"
            >
              Create Group
            </button>
          </div>
        </form>
      </div>
  );
}
