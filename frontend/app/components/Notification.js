'use client'
import { useState, useRef, useEffect } from 'react'
import BellIcon from '../../public/bell.png'
import Image from 'next/image'

// check for type to know what to display
// type = new_event -> Event -> for group members
// type = group_invite -> Request -> if somebody sends you a group invite
// type = follow_request -> Request -> follower requests
// type = join_request -> Request -> join group requests -> will only show for group creator

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
              notifications.map((notification) => {
                let displayMessage = 'You have a new notification.'; // Default message

                // Helper to safely get a user's name
                const getUserName = (user) => {
                  if (!user) return 'Someone';
                  return user.nickname || `${user.first_name || ''} ${user.last_name || ''}`.trim() || 'Someone';
                };

                // Helper to safely get a group's name
                const getGroupName = (group) => {
                  return group?.group_name || 'a group';
                };

                // Log the whole notification object to see its top-level structure
                // console.log('Processing notification:', notification);

                if (notification.Type === 'new_event' && notification.Event) {
                  const eventTitle = notification.Event.title || 'an event';
                  displayMessage = `New event created: "${eventTitle}"`;
                  if (notification.Event.description) {
                    displayMessage += ` - ${notification.Event.description.substring(0, 30)}${notification.Event.description.length > 30 ? '...' : ''}`;
                  }
                } else if (notification.Type === 'group_invite') {
                  // Specifically log for 'group_invite'
                  console.log('GROUP INVITE DATA:', notification);
                  console.log('Request Object:', notification.Request);
                  if (notification.Request) {
                    console.log('Sender Object:', notification.Request.Sender);
                    console.log('Group Object:', notification.Request.Group);
                  }

                  // Original condition
                  if (notification.Request?.Sender && notification.Request?.Group) {
                    const senderName = getUserName(notification.Request.Sender);
                    const groupName = getGroupName(notification.Request.Group);
                    displayMessage = `${senderName} sent you an invite to join ${groupName}.`;
                  } else {
                    console.log('Group invite condition not met: Request, Sender, or Group might be missing/undefined.');
                  }
                } else if (notification.Type === 'follow_request' && notification.Request?.Sender) {
                  const senderName = getUserName(notification.Request.Sender);
                  displayMessage = `${senderName} sent you a follow request.`;
                } else if (notification.Type === 'join_request' && notification.Request?.Sender && notification.Request?.Group) {
                  const senderName = getUserName(notification.Request.Sender);
                  const groupName = getGroupName(notification.Request.Group);
                  // Assuming the current user is the creator/admin of the group this request is for
                  displayMessage = `${senderName} wants to join your group ${groupName}.`;
                } else if (notification.message) {
                  // Fallback to a pre-existing message if the backend provides one and other types don't match
                  displayMessage = notification.message;
                }

                return (
                  <div
                    key={notification.notification_id}
                    className="px-4 py-2 border-b border-gray-300 hover:bg-gray-50" // Added hover effect
                  >
                    <p className="text-sm text-gray-700">
                      {displayMessage}
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
                          // timeZone: 'UTC', // Consider if UTC is always desired or local time
                        }
                      )}
                    </p>
                  </div>
                );
              })
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