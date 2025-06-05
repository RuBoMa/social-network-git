'use client'
import { useState, useRef, useEffect } from 'react'
import { useUser } from '../context/UserContext'
import BellIcon from '../../public/bell.png'
import Image from 'next/image'
import Link from 'next/link'
import { sendMessage, addMessageHandler } from './ws'

export default function NotificationsDropdown() {
  const { user: currentUser } = useUser()
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
      if (res.ok) {
        const data = await res.json()
        console.log('Fetched notifications:', data)
        setNotifications(data)
      }
    }
    fetchNotifications()
  }, [])
  
  // Close dropdown when clicking outside
  useEffect(() => {
  if (!open) return;

  function handleClickOutside(event) {
    if (containerRef.current && !containerRef.current.contains(event.target)) {
      setOpen(false);
    }
  }

  document.addEventListener('mousedown', handleClickOutside);
  return () => {
    document.removeEventListener('mousedown', handleClickOutside);
  };
}, [open]);

  // inside NotificationsDropdown
  useEffect(() => {
    if (!localStorage.getItem('token')) return;

    const removeHandler = addMessageHandler((data) => {
      if (
        data.type === 'group_invite' || 
        data.type === 'follow_request' || 
        data.type === 'join_request' || 
        data.type === 'new_event' ||
        data.type === 'join_accepted'
      ) {
        console.log('Received notification:', data);
        setNotifications(prev => Array.isArray(prev) ? [data, ...prev] : [data]);
      } else if (data.type === 'mark_notification_read') {
        setNotifications(prev =>
          prev?.filter(n =>
            n.notification_id !== data.notification_id
          )
        );
        setOpen(false); // Close dropdown when marking as read
      }
    });

    return () => {
      if (removeHandler) removeHandler();
    };
  }, []);
  
  async function markAsRead(id) {
    console.log('Marking notification as read:', id)
    sendMessage({
      type: 'mark_notification_read',
      notification_id: id,
    })
  }

  return (
    <div className="relative">
      <button onClick={() => setOpen(o => !o)} className="p-1 cursor-pointer" aria-label="Notifications">
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
                  const sender = notification.request?.sender;
                  const senderName =
                    notification.type === 'join_request'
                    ?  notification.request?.joining_user?.nickname ||
                      [notification.request?.joining_user?.first_name, notification.request?.joining_user?.last_name].filter(Boolean).join(" ") || "Someone"
                    :  sender?.nickname && sender.nickname.trim() !== ""
                        ? sender.nickname
                        : [sender?.first_name, sender?.last_name].filter(Boolean).join(" ") || "Someone";

                  if (notification.type === 'group_invite') {
                    displayMessage = `${senderName} invited you to join "${notification.request.group.group_name}".`;
                    href = `/group?group_id=${notification.request.group.group_id}`;
                  } else if (notification.type === 'follow_request') {
                    displayMessage = `${senderName} sent you a follow request.`;
                    href = `/profile?user_id=${currentUser.user_id}`;
                  } else if (notification.type === 'join_request') {
                    displayMessage = `${senderName} requested to join "${notification.request.group.group_name}".`;
                    href = `/group?group_id=${notification.request.group.group_id}`;
                  } else if (notification.type === 'new_event') {
                    displayMessage = `A new event "${notification.event.title}" has been created in "${notification.event.group.group_name}".`;
                    href = `/event?event_id=${notification.event.event_id}`;
                  } else if (notification.type === 'join_accepted') {
                    displayMessage = `Your request to join "${notification.request.group.group_name}" has been accepted.`;
                    href = `/group?group_id=${notification.request.group.group_id}`;
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
                          'en-GB',
                          {
                            year: "numeric",
                            month: "2-digit",
                            day: "2-digit",
                            hour: "2-digit",
                            minute: "2-digit",
                            second: "2-digit",
                            hour12: false,
                            timeZone: "UTC",
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