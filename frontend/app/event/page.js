'use client'
import { useState, useEffect } from 'react'
import { useSearchParams } from 'next/navigation'
import Author from '../components/Author'

export default function EventPage() {
    const searchParams = useSearchParams()
    const eventId = searchParams.get('event_id')

    const [event, setEvent] = useState(null)
    const [error, setError] = useState(null)
    const [attendance, setAttendance] = useState(null)
    const [attendLoading, setAttendLoading] = useState(false)

    useEffect(() => {
        async function fetchEvent() {
          if (!eventId) return

          const res = await fetch(`http://localhost:8080/api/event?event_id=${eventId}`, {
            credentials: 'include', 
            method: 'GET',
            headers: {
              'Accept': 'application/json' //telling the server we want JSON
            }
          })
          console.log('Response status for event:', res) // Log the response status
    
    
          const data = await res.json()
          if (res.ok) {
            console.log('Response is OK') // Log if the response is OK
            console.log('Fetched event:', data) // Log the fetched posts
            setEvent(data)
            setAttendance(data.user_response || null)
          } else {
            console.error('Failed to load posts')
            setError(data.message || 'Failed to load posts')

          }
        }
    
        fetchEvent()
      }, [eventId])

      async function handleAttendance(response) {
        setAttendLoading(true);
        try {
            const res = await fetch('http://localhost:8080/api/event-attendance', {
            method: 'POST',
            credentials: 'include',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                event_id: event.event_id,
                response, // 'going' või 'not going'
            }),
            });

        console.log('Response status for attendance:', res) // Log the response status
            if (res.ok) {
            setAttendance(response);

        const updatedRes = await fetch(`http://localhost:8080/api/event?event_id=${event.event_id}`, {
          credentials: 'include',
          headers: { 'Accept': 'application/json' },
        })
        if (updatedRes.ok) {
          const updatedData = await updatedRes.json()
          setEvent(updatedData.event)
        }
      }
    } finally {
        setAttendLoading(false);
    }
    }

      if (error) {
        return <div className="text-red-600">{error}</div>
      }

      if (!event) {
        return <div>Loading event...</div>
      }

    const goingUsers = (event.event_responses || []).filter(r => r.response === 'going')
    const notGoingUsers = (event.event_responses || []).filter(r => r.response === 'not_going')

        return (
        <div className="flex items-center justify-center bg-gray-100 p-4">
            <div className="max-w-xl p-10 p-4 bg-white rounded shadow">
                <div className="flex items-center justify-between mb-2">
                    <h1 className="text-2xl font-bold mb-2">{event.title}</h1>
                    <span className="text-2xl font-extrabold text-black mb-2">
                        {event.event_date
                        ? new Date(event.event_date).toLocaleString('en-GB', {
                        day: '2-digit',
                        month: 'long',
                        year: 'numeric',
                        hour: '2-digit',
                        minute: '2-digit',
                        hour12: false,
                        timeZone: 'UTC',
                        })
                        : 'No date provided'}
                    </span>
                </div>
            <p className="mb-2">{event.description}</p>

            {/* Going and Not going buttons */}
            <div className="flex gap-2 mt-4 justify-end">
                <button
                    className={`px-3 py-1 rounded font-bold ${
                    attendance === 'going' ? 'bg-green-600 text-white' : 'bg-gray-200 text-gray-800'
                }`}
                    disabled={attendance === 'going' || attendLoading}
                    onClick={() => handleAttendance('going')}
                    >
                    Going
                </button>
                <button
                    className={`px-3 py-1 rounded font-bold ${
                    attendance === 'not going' ? 'bg-red-600 text-white' : 'bg-gray-200 text-gray-800'
                }`}
                    disabled={attendance === 'not going' || attendLoading}
                    onClick={() => handleAttendance('not going')}
                    >
                    Not going
                </button>
                </div>

                {/* Displaying the attendace */}
                <div className="mt-6 flex gap-1">
                    <h2 className="text-sm font-semibold mb-1">Going ({goingUsers.length}):</h2>
                    <ul className="list-inside mb-4">
                    {goingUsers.length === 0 && <li>No one is going yet.</li>}
                    {goingUsers.map((resp) => (
                        <li key={resp.user_id}>
                         <Author author={user} size="sm" />
                        </li>
                    ))}
                    </ul>

                    <h2 className="text-sm font-semibold mb-1">Not going ({notGoingUsers.length}):</h2>
                    <ul className="list-inside">
                    {notGoingUsers.length === 0 && <li>No one has declined yet.</li>}
                    {notGoingUsers.map((resp) => (
                        <li key={resp.user_id}>
                         <Author author={user} size="sm" />
                        </li>
                    ))}
                    </ul>
                </div>
            </div>
        </div>
            
        )

    }