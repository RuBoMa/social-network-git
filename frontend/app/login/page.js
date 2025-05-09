'use client'
import { useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'

export default function LoginPage() {
  const router = useRouter()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')

  async function handleSubmit(e) {
    e.preventDefault()
    const res = await fetch('http://localhost:8080/api/login', {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password })
    })
    if (res.ok) router.push('/')
    else alert('Login failed')
  }

  return (
    <div className="p-8">
      <h1 className="text-2xl mb-4">Login</h1>
      <form onSubmit={handleSubmit}>
        <input
          type="email"
          value={email}
          onChange={e => setEmail(e.target.value)}
          placeholder="Email"
          required
          className="block w-full mb-2"
        />
        <input
          value={password}
          onChange={e => setPassword(e.target.value)}
          name="password"
          type="password"
          placeholder="Password"
          className="block mb-4"
        />
        <button
          type="submit"
          className="bg-blue-500 text-white px-4 py-2"
        >
          Sign In
        </button>
      </form>
      <p className="mt-4">
        Don’t have an account?{' '}
        <Link href="/signup" className="text-blue-500 underline">
          Sign up
        </Link>
      </p>
    </div>
  )
}