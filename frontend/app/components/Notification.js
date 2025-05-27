'use client'
import { useState, useRef, useEffect } from 'react'
import BellIcon from '../../public/bell.png'
import Image from 'next/image'
import Link from 'next/link'
import { sendMessage, addMessageHandler } from './ws'

export default function NotificationsDropdown() {
  const [open, setOpen] = useState(false)
  const containerRef = useRef()
  const [notifications, setNotifications] = useState([])

  useEffect(() => {
    if (!localStorage.getItem('token')) {
      return;
    }
    async function fetchNotifications() {
      const res = await fetch('http://localhost:8080/api/notifications', {
        credentials: 'include',
      })
      if (res.ok) setNotifications(await res.json())
    }
    fetchNotifications()
  }, [])

    // inside NotificationsDropdown
  useEffect(() => {
    if (!localStorage.getItem('token')) return;

    const removeHandler = addMessageHandler((data) => {
      if (data.type === 'group_invite' || data.type === 'follow_request' || data.type === 'join_request' || data.type === 'new_event') {
        setNotifications(prev => [data, ...prev]);
      } else if (data.type === 'mark_notification_read') {
        setNotifications(prev =>
          prev.map(n =>
            n.notification_id === data.notification_id ? { ...n, is_read: true } : n
          )
        );
      }
    });

    return () => {
      if (removeHandler) removeHandler();
    };
  }, []);
  
  async function markAsRead(id) {
    sendMessage({
      type: 'mark_notification_read',
      notification_id: id,
    })
  }

  return (
    <div className="relative cursor-pointer">
      <button onClick={() => setOpen(o => !o)} className="p-1" aria-label="Notifications">
        <Image src={BellIcon} alt="Notifications" width={22} height={22}/>
          {notifications?.length > 0 && (
          <span className="absolute top-0 right-0 flex items-center justify-center h-4 w-4 rounded-full bg-red-600 text-white text-xs font-semibold">
            {notifications.length}
          </span>
          )}
      </button>

      {open && (
        <div ref={containerRef} className="absolute right-0 mt-2 w-64 bg-white rounded shadow-lg z-50 overflow-hidden border border-gray-200">
          <h4 className="px-4 py-2 border-b border-gray-300 font-semibold">Notifications</h4>
          <div className="max-h-64 overflow-y-auto">
            {Array.isArray(notifications) && notifications.length > 0
              ? notifications.map(notification => {
                  const readStyle = notification.is_read ? 'bg-gray-100' : 'hover:bg-gray-50'
                  let href = '#'
                  let displayMessage = 'New Notification';

                  if (notification.type === 'group_invite') {
                    displayMessage = `${notification.request.sender.nickname} invited you to join "${notification.request.group.group_name}".`;
                    href = `/group?group_id=${notification.request.group.group_id}`;
                  } else if (notification.type === 'follow_request') {
                    displayMessage = `${notification.request.sender.nickname} sent you a follow request.`;
                    // not implemented visually yet
                    // backend as well?
                  } else if (notification.type === 'join_request') {
                    displayMessage = `${notification.request.sender.nickname} requested to join "${notification.request.group.group_name}".`;
                    href = `/group?group_id=${notification.request.group.group_id}`;
                  } else if (notification.type === 'new_event') {
                    displayMessage = `A new event "${notification.event.title}" has been created in "${notification.event.group.group_name}".`;
                    href = `/event?event_id=${notification.event.event_id}`;
                  }

                  return (
                    <Link
                      key={notification.notification_id}
                      href={href}
                      className={`block px-4 py-2 border-b border-gray-300 ${readStyle}`}
                      onClick={() => markAsRead(notification.notification_id)}
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
                          }
                        )}
                      </p>
                    </Link>
                  )
                }) : (
                  <div className="px-4 py-2 text-sm text-gray-500">
                    No notifications.
                  </div>
                )
            }
          </div>
        </div>
      )}
    </div>
  )
}