'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';

export default function SearchBar() {
    const [query, setQuery] = useState('');
    const [results, setResults] = useState(null);

    useEffect(() => {
        if (query.length === 0) {
            setResults(null);
            return;
        }

        const delayDebounce = setTimeout(() => {
            fetch(`/api/search?q=${encodeURIComponent(query)}`)
                .then(res => res.json())
                .then(data => setResults(data))
                .catch(err => {
                    console.error('Search error:', err);
                    setResults(null);
                });
        }, 300);

        return () => clearTimeout(delayDebounce);
    }, [query]);

    return (
        <div className="relative w-full max-w-md">
            <input
                type="text"
                value={query}
                onChange={(e) => setQuery(e.target.value)}
                placeholder="Search..."
                className="w-full p-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />

            {results && (
                <div className="absolute top-full mt-2 w-full bg-white border border-gray-300 rounded shadow-lg z-50">
                    <div className="p-2 max-h-64 overflow-y-auto">
                        {results.Users?.length > 0 && (
                            <div>
                                <p className="font-bold">Users</p>
                                {results.Users.map((u) => (
                                    <Link key={u.id} href={`/profile?user_id=${u.id}`} className="block hover:bg-gray-100 p-1">
                                        {u.username}
                                    </Link>
                                ))}
                            </div>
                        )}
                        {results.Groups?.length > 0 && (
                            <div>
                                <p className="font-bold mt-2">Groups</p>
                                {results.Groups.map((g) => (
                                    <Link key={g.id} href={`/group?group_id=${g.id}`} className="block hover:bg-gray-100 p-1">
                                        {g.name}
                                    </Link>
                                ))}
                            </div>
                        )}
                        {results.Posts?.length > 0 && (
                            <div>
                                <p className="font-bold mt-2">Posts</p>
                                {results.Posts.map((p) => (
                                    <Link key={p.id} href={`/post?post_id=${p.id}`} className="block hover:bg-gray-100 p-1">
                                        {p.title}
                                    </Link>
                                ))}
                            </div>
                        )}
                        {results.Events?.length > 0 && (
                            <div>
                                <p className="font-bold mt-2">Events</p>
                                {results.Events.map((e) => (
                                    <Link key={e.id} href={`/event?event_id=${e.id}`} className="block hover:bg-gray-100 p-1">
                                        {e.name}
                                    </Link>
                                ))}
                            </div>
                        )}
                        {results.Users?.length === 0 &&
                            results.Groups?.length === 0 &&
                            results.Posts?.length === 0 &&
                            results.Events?.length === 0 && <p>No results found</p>}
                    </div>
                </div>
            )}
        </div>
    );
}
