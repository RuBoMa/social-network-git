'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';

export default function SearchBar() {
    const [query, setQuery] = useState('');
    const [results, setResults] = useState(null);

    const clearSearch = () => {
        setQuery('');
        setResults(null);
    };

    useEffect(() => {
        if (query.length === 0) {
            setResults(null);
            return;
        }

        const delayDebounce = setTimeout(() => {
                fetch(`http://localhost:8080/api/search?q=${query}`, {
                method: 'GET',            
                credentials: 'include',   
                })
                .then(res => res.json())
                .then((data) => {
                    console.log("Search results:", data); 
                    setResults(data);
                })
                .catch(err => {
                    console.error('Search error:', err);
                    setResults(null);
                });
        }, 300);


        return () => clearTimeout(delayDebounce);
    }, [query]);

    const isEmptyResult = results &&
        (!results.users?.length && !results.groups?.length && !results.posts?.length && !results.events?.length);


    return (
        <div className="relative w-full max-w-md">
            <input
                type="text"
                value={query}
                onChange={(e) => setQuery(e.target.value)}
                placeholder="Search..."
                maxLength={50}
                className="w-full p-2 border border-gray-400 bg-gray-100 rounded-md focus:outline-none focus:ring-1 focus:ring-gray-600"
            />

            {results && (
                <div className="absolute top-full mt-2 w-full bg-white border border-gray-300 rounded shadow-lg z-50">
                    <div className="p-2 max-h-64 overflow-y-auto">
                        {results.users?.length > 0 && (
                            <div>
                                <p className="font-bold">Users</p>
                                {results.users.map((u) => (
                                    <Link key={u.user_id} href={`/profile?user_id=${u.user_id}`}>
                                        <span
                                            onClick={clearSearch}

                                            className="block hover:bg-gray-100 p-1"
                                        >
                                            {u.nickname || [u.first_name, u.last_name].filter(Boolean).join(" ") || "Unknown user"}
                                        </span>
                                    </Link>
                                ))}
                            </div>
                        )}
                        {results.groups?.length > 0 && (
                            <div>
                                <p className="font-bold mt-2">Groups</p>
                                {results.groups.map((g) => (
                                    <Link key={g.group_id} href={`/group?group_id=${g.group_id}`}>
                                        <span
                                            onClick={clearSearch}
                                            className="block hover:bg-gray-100 p-1"
                                        >
                                            {g.group_name}
                                        </span>
                                    </Link>
                                ))}
                            </div>
                        )}
                        {results.posts?.length > 0 && (
                            <div>
                                <p className="font-bold mt-2">Posts</p>
                                {results.posts.map((p) => (
                                    <Link key={p.post_id} href={`/post?post_id=${p.post_id}`}>
                                        <span
                                            onClick={clearSearch}
                                            className="block hover:bg-gray-100 p-1"
                                        >
                                            {p.post_title}
                                        </span>
                                    </Link>
                                ))}
                            </div>
                        )}
                        {results.events?.length > 0 && (
                            <div>
                                <p className="font-bold mt-2">Events</p>
                                {results.events.map((e) => (
                                    <Link key={e.event_id} href={`/event?event_id=${e.event_id}`}>
                                        <span
                                            onClick={clearSearch}
                                            className="block hover:bg-gray-100 p-1"
                                        >
                                            {e.title}
                                        </span>
                                    </Link>
                                ))}
                            </div>
                        )}
                        {isEmptyResult && (
                            <p className="text-gray-500 font-bold mt-2">No results found</p>
                        )}
                    </div>
                </div>
            )}
        </div>
    );
}
