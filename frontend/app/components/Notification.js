'use client'
import { useState, useRef, useEffect } from 'react'
import BellIcon from '../../public/bell.png'
import Image from 'next/image'

// Fetch from database (websocket?)

const fetchNotifications = async () => {
  const res = await fetch('http://localhost:8080/api/notifications', {
    credentials: 'include',
    method: 'GET',
    headers: {
      Accept: 'application/json',
    },
  })
  const data = await res.json()
  if (res.ok) {
    return data
  } else {
    console.error('Failed to load notifications')
    return []
  }
}
const mockNotifications = [
  { id: 1, message: 'Alice commented on your post', time: '2m ago' },
  { id: 2, message: 'Bob liked your post', time: '10m ago' },
  { id: 3, message: 'Carol invited you to “Study” Group', time: '1h ago' },
  { id: 4, message: 'Erik invited you to “Fika at my place!” Group', time: '3h ago' },
  { id: 5, message: 'Carol wants to follow you', time: '1d ago' },
]

export default function NotificationsDropdown() {
  const [open, setOpen] = useState(false)
  const containerRef = useRef() // reference to the dropdown container

  // close on outside click
  useEffect(() => {
    function onClick(e) {
      if (containerRef.current && !containerRef.current.contains(e.target)) {
        setOpen(false)
      }
    }
    document.addEventListener('mousedown', onClick)
    return () => document.removeEventListener('mousedown', onClick)
  }, [])

  return (
    <div className="relative cursor-pointer">
      <button
        onClick={() => setOpen(o => !o)}
        className="p-1"
        aria-label="Notifications"
      >
        <Image
          src={BellIcon}
          alt="Notifications"
          width={22}
          height={22}
        />
      </button>

      {open && (
        <div
          ref={containerRef} // clip dropdown to container
          className="
            absolute right-0 mt-2 w-64
            bg-white rounded shadow-lg z-50
            overflow-hidden
            border border-gray-200
            transition-transform transform
          "
        >
          <h4 className="px-4 py-2 border-b border-gray-300 font-semibold">
            Notifications
          </h4>
          <div className="max-h-64 overflow-y-auto">
            <ul>
              {mockNotifications.map(n => (
                <li key={n.id} className="px-4 py-2 hover:bg-gray-100">
                  <p className="text-sm">{n.message}</p>
                  <p className="text-xs text-gray-400">{n.time}</p>
                </li>
              ))}
              {mockNotifications.length === 0 && (
                <li className="px-4 py-2 text-sm text-gray-500">
                  No notifications
                </li>
              )}
            </ul>
          </div>
        </div>
      )}
    </div>
  )
}