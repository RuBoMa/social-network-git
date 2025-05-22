'use client';

import Link from 'next/link';

export default function EventSection({ events }) {
  if (!events || events.length === 0) {
    return (
      <div className="my-6">
        <h3 className="text-xl font-semibold mb-3 text-gray-800">Upcoming Events</h3>
        <div className="p-4 bg-gray-50 rounded-lg shadow-sm">
          <p className="text-gray-500">No upcoming events for this group.</p>
        </div>
      </div>
    );
  }

  return (
    <div className="my-6 max-w-4xl mx-auto">
      <h3 className="text-xl font-semibold mb-3 text-gray-800">Upcoming Events</h3>
      <div className="flex overflow-x-auto space-x-4 pb-4">
        {events.map(event => (
          <Link 
            key={event.event_id || event.id} 
            href={`/event?event_id=${event.event_id || event.id}`}
            className="block min-w-[280px] max-w-[320px] bg-white rounded-xl shadow-lg hover:shadow-xl transition-shadow duration-300 ease-in-out border border-gray-200 overflow-hidden"
          >
            <div className="p-5">
              <h4 className="text-lg font-bold text-blue-600 mb-2 truncate" title={event.title}>
                {event.title}
              </h4>
              <p className="text-sm text-gray-500 mb-1">
                {new Date(event.event_time || event.event_date).toLocaleString([], { 
                  year: 'numeric', 
                  month: 'long', 
                  day: 'numeric', 
                  hour: '2-digit', 
                  minute: '2-digit', 
                  hour12: false,
                  timezone: 'UTC'
                })}
              </p>
              <p className="text-sm text-gray-700 line-clamp-3" title={event.description}>
                {event.description}
              </p>
            </div>
          </Link>
        ))}
      </div>
    </div>
  );
}