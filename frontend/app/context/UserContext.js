'use client';

import { createContext, useContext, useState, useEffect } from 'react';

const UserContext = createContext();

export function UserProvider({ children }) {
    const [user, setUserState] = useState(null);

    // Update both state and localStorage
    const setUser = (newUser) => {
    setUserState(newUser);
    if (newUser) {
        localStorage.setItem('user', JSON.stringify(newUser));
    } else {
        localStorage.removeItem('user');
    }
    };

    // Rehydrate on mount
    useEffect(() => {
        const stored = localStorage.getItem('user');
        if (stored) {
            setUserState(JSON.parse(stored));
        }
    }, []);

    return (
        <UserContext.Provider value={{ user, setUser }}>
            {children}
        </UserContext.Provider>
        );
    }

    export function useUser() {
        const context = useContext(UserContext);
        if (!context) {
            throw new Error('useUser must be used within a UserProvider');
        }
    return context;
}
