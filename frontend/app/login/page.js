'use client'
import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { useUser } from '../context/UserContext'
import Link from 'next/link'
import initWebSocket from '../components/ws'

export default function LoginPage() {
  const router = useRouter()
  const { setUser } = useUser(); // Access the setUser function from context
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')

  async function handleSubmit(e) {
    e.preventDefault()

    localStorage.removeItem('token'); // Clear any existing token
    localStorage.removeItem('user'); // Clear any existing user data

    const res = await fetch('http://localhost:8080/api/login', {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json', },
      body: JSON.stringify({ email, password })
    })
    
    // backend will send "token", need to save to localstorage
    if (res.ok) {
      const userData = await res.json(); // Fetch user data from the response
      console.log('User data from backend:', userData);
      setUser(userData); // Save user data
      localStorage.setItem('token', userData.token); // Save token to local storage
      localStorage.setItem('user', JSON.stringify(userData)); // Save user data to local storage
      initWebSocket(); // Initialize WebSocket connection
      router.replace('/'); // Redirect to the home page
    } else {
      alert('Login failed')
    }
  }

  return (
    <div className="flex items-center justify-center min-h-screen">
      <div className="p-8 max-w-md w-full bg-white rounded-lg shadow-lg">
        <h1 className="text-2xl mb-4 text-center">Welcome</h1>
        <form onSubmit={handleSubmit}>
          <label className="block mb-2">
            Email
            <input
              type="email"
              value={email}
              onChange={e => setEmail(e.target.value)}
              required
              className="block w-full mt-1 mb-2 border rounded px-2 py-1"
            />
          </label>
          <label className="block mb-4">
            Password
            <input
              type="password"
              value={password}
              onChange={e => setPassword(e.target.value)}
              required
              className="block w-full mt-1 border rounded px-2 py-1"
            />
          </label>
          <button
            type="submit"
            className="bg-sky-600/60 hover:bg-sky-700/60 text-white font-bold rounded text-md my-2 mr-2 p-2 px-4 w-full cursor-pointer"
          >
            Sign in
          </button>
        </form>
        <p className="mt-4 text-center">
          Don’t have an account?{' '}
          <Link href="/signup" className="text-sky-800/60 underline">
            Sign up
          </Link>
        </p>
      </div>
    </div>
  )
}