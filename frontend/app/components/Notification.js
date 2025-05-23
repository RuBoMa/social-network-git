'use client'
import { useState, useRef, useEffect } from 'react'
import BellIcon from '../../public/bell.png'
import Image from 'next/image'

// Fetch from database (websocket?)

// should fetch notification, but only once when page is loaded
export default function NotificationsDropdown() {
  const [open, setOpen] = useState(false)
  const containerRef = useRef() // reference to the dropdown container
  const [notifications, setNotifications] = useState(null)

  useEffect(() => {
    async function fetchNotifications() {
      try {
        const res = await fetch('http://localhost:8080/api/notifications', {
          credentials: 'include',
        })

        if (res.ok) {

          const data = await res.json()
          setNotifications(data)
          console.log('Fetched notifications:', data) // Log the fetched notifications
        }
      } catch (err) {
        console.error('Error fetching notifications:', err)
      }
    }

    fetchNotifications()
  }, []) // empty dependency array to run only once

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
            {notifications?.length > 0 ? (
              notifications.map((notification) => (
                <div
                  key={notification.notification_id}
                  className="px-4 py-2 border-b border-gray-300"
                >
                  <p className="text-sm text-gray-700">
                    {notification.message}
                  </p>
                  <p className="text-xs text-gray-500">
                    {new Date(notification.created_at).toLocaleString(
                      'en-US',
                      {
                        year: 'numeric', 
                        month: 'long', 
                        day: 'numeric', 
                        hour: '2-digit', 
                        minute: '2-digit', 
                        hour12: false,
                        timeZone: 'UTC',
                      }
                    )}
                  </p>
                </div>
              ))
            ) : (
              <div className="px-4 py-2 text-sm text-gray-500">
                No notifications to show.
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  )
}