'use client';
import { useState } from 'react';
import { useSearchParams } from 'next/navigation';

// needs to be updated to work with backend!

export default function CreateEvent({ onClose, onSuccess }) {
  const searchParams = useSearchParams();
  const groupId = Number(searchParams.get('group_id'));

  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [date, setDate] = useState('');

  async function handleNewEvent(e) {
    e.preventDefault();
    if (!name.trim() || !description.trim() || !date) return;

    const res = await fetch('http://localhost:8080/api/create-event', {
      method:      'POST',
      credentials: 'include',
      headers:     { 'Content-Type': 'application/json' },
      body:        JSON.stringify({
        title:      name,
        description,
        event_date: date,
        group:      { group_id: groupId }
      })
    });

    if (res.ok) {
      onSuccess  && onSuccess();
      onClose    && onClose();
    } else {
      const err = await res.json();
      console.error('CreateEvent failed:', err.message);
    }
  }

  return (
    <form
      onSubmit={handleNewEvent}
      className="max-w-full mx-0 mt-1 p-2 p-4 rounded shadow border border-gray-200"
    >
      <label className="block mb-4">
        <input
          type="text"
          placeholder="Event Name"
          value={name}
          onChange={e => setName(e.target.value)}
          required
          maxLength={50}
          className="mt-1 block w-full border border-gray-300 rounded p-2"
        />
      </label>

      <label className="block mb-4">
        <textarea
          placeholder="Event Description"
          value={description}
          onChange={e => setDescription(e.target.value)}
          required
          maxLength={800}
          className="mt-1 block w-full border border-gray-300 rounded p-2 h-24 resize-none"
        />
      </label>

      <label className="block mb-4">
        <input
          type="datetime-local"
          value={date}
          onChange={e => setDate(e.target.value)}
          required
          className="mt-1 block w-full border border-gray-300 rounded p-2"
        />
      </label>

      <div className="flex justify-end">
        <button
          type="submit"
          className="bg-blue-500 hover:bg-blue-600 text-white px-2 py-1 rounded text-sm"
        >
          Submit
        </button>
      </div>
    </form>
  );
}
