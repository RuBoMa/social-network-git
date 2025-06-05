'use client'
import { useState, useEffect } from 'react'
// import { useSearchParams } from 'next/navigation'
import Author from '../components/Author'
import BackButton from '../components/BackButton'

export default function EventPage() {
    // const searchParams = useSearchParams()
    // const eventId = searchParams.get('event_id')

    const [eventId, setEventId] = useState(null)

    useEffect(() => {
        // Only run in browser
        const params = new URLSearchParams(window.location.search)
        setEventId(params.get('event_id'))
    }, [])

    const [event, setEvent] = useState(null)
    const [error, setError] = useState(null)
    const [attendance, setAttendance] = useState(null)
    const [attendLoading, setAttendLoading] = useState(false)

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
            setAttendance(data.attendance || null)
        } else {
            console.error('Failed to load posts')
            setError(data.message || 'Failed to load posts')
        }
    }

    useEffect(() => {
        fetchEvent()
    }, [eventId])

    async function handleAttendance(response) {
        setAttendLoading(true);
        try {
            //query to the server to update the attendance
            const res = await fetch('http://localhost:8080/api/event-attendance', {
                method: 'POST',
                credentials: 'include',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    event: { event_id: event.event_id },
                    response, // 'going','not going'
                }),
            });

            console.log('Response status for attendance:', res) // Log the response status
            if (res.ok) {
                await fetchEvent(); //fetch the event again to get the updated attendance
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

    const goingUsers = (event.members_going || [])
    console.log('Going users:', goingUsers) // Log the going users

    return (
        <div className="p-4">
            <BackButton />
        <div className="p-4 flex justify-center items-center">
               <div className="max-w-xl p-10 p-4 bg-white rounded shadow flex flex-col items-center text-center">
            <div className="flex justify-between w-full mb-1 gap-4">
                    <img
                        src="/event_logo.png"
                        alt="Event logo"
                        className="w-24 h-24 object-contain"
                    />
                    {/* date box*/}
                    {event.event_date && (
                        <div className="flex flex-col items-center justify-center bg-sky-100 rounded-lg shadow px-5 py-2 mr-4 min-w-[70px]">
                        <span className="text-2xl font-bold text-sky-700">
                            {new Date(event.event_date).toLocaleString('en-GB', { day: '2-digit', timeZone: 'UTC' })}
                        </span>
                        <span className="uppercase text-sm text-sky-600 font-semibold">
                            {new Date(event.event_date).toLocaleString('en-GB', { month: 'short', timeZone: 'UTC' })}
                        </span>
                        <span className="text-xs text-sky-500">
                            {new Date(event.event_date).toLocaleString('en-GB', { year: 'numeric', timeZone: 'UTC' })}
                        </span>
                        <span className="text-sm text-sky-600 mt-2 font-semibold">
                            {new Date(event.event_date).toLocaleString('en-GB', { hour: '2-digit', minute: '2-digit', hour12: false, timeZone: 'UTC' })}
                        </span>
                        </div>
                    )}
                    </div>
                    <h1 className="text-xl font-bold">{event.title}</h1>
                     <p className="mb-2">{event.description}</p>

                    {/* Going and Not going buttons */}
                    <div className="flex gap-2 mt-4 justify-center">
                        <button
                            className={`px-3 py-1 rounded font-bold ${attendance === 'going' ? 'bg-sky-700 text-white' : 'bg-gray-200 text-gray-800'
                                }`}
                            disabled={attendance === 'going' || attendLoading}
                            onClick={() => handleAttendance('going')}
                        >
                            Going
                        </button>
                        <button
                            className={`px-3 py-1 rounded font-bold ${attendance === 'not going' ? 'bg-sky-700 text-white' : 'bg-gray-200 text-gray-800'
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
                                    <Author author={resp} size="sm" />
                                </li>
                            ))}
                        </ul>
                    </div>
            </div>
        </div>
        </div>
    )
}