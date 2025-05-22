'use client'
import { useState, useEffect } from 'react'
import { useSearchParams } from 'next/navigation'

export default function EventPage() {
    const searchParams = useSearchParams()
    const eventId = searchParams.get('event_id')

    const [event, setEvent] = useState(null)
    const [error, setError] = useState(null)

    
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
          console.log('Response status:', res) // Log the response status
    
    
          const data = await res.json()
          if (res.ok) {
            console.log('Response is OK') // Log if the response is OK
            console.log('Fetched event:', data) // Log the fetched posts
            setEvent(data)
          } else {
            console.error('Failed to load posts')
            setError(data.message || 'Failed to load posts')

          }
        }
    
        fetchEvent()
      }, [eventId])

      if (error) {
        return <div className="text-red-600">{error}</div>
      }

      if (!event) {
        return <div>Loading event...</div>
      }
        return (
            <div className="max-w-xl mx-auto p-4 bg-white rounded shadow">
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
            </div>
        )

    }