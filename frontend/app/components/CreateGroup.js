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

  }

  return (
    <form
      onSubmit={handleNewGroup}
      className="max-w-full mx-auto mt-1 mb-4 p-4 rounded shadow border border-gray-200 bg-white"
    >
      <h2 className="text-xl font-bold mb-4">Create Group</h2>

      <label className="block mb-4">
        <input
          type="text"
          placeholder="Group Name"
          value={name}
          onChange={e => setName(e.target.value)}
          required
          className="mt-1 block w-full border border-gray-300 rounded p-2"
        />
      </label>

      <label className="block mb-6">
        <textarea
          placeholder="Group Description"
          value={description}
          onChange={e => setDescription(e.target.value)}
          required
          className="mt-1 block w-full border border-gray-300 rounded p-2 h-24 resize-none"
        />
      </label>

      <div className="flex justify-end space-x-2">
        <button
          type="submit"
          className="px-3 py-1 bg-blue-600 text-white rounded hover:bg-blue-700 text-sm"
        >
          Submit
        </button>
      </div>
    </form>
  );
}
