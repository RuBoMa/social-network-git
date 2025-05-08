'use client'
import { useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'

export default function SignupPage() {
  const router = useRouter()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [firstName, setFirstName] = useState('')
  const [lastName, setLastName] = useState('')
  const [dob, setDob] = useState('')
  const [avatar, setAvatar] = useState(null)
  const [nickname, setNickname] = useState('')
  const [aboutMe, setAboutMe] = useState('')

  async function handleSubmit(e) {
    e.preventDefault()
    const formData = new FormData()
    formData.append('email', email)
    formData.append('password', password)
    formData.append('first_name', firstName)
    formData.append('last_name', lastName)
    formData.append('date_of_birth', dob)
    if (avatar) formData.append('avatar', avatar)
    if (nickname) formData.append('nickname', nickname)
    if (aboutMe) formData.append('about_me', aboutMe)

    const res = await fetch('http://localhost:8080/api/signup', {
      method: 'POST',
      credentials: 'include',
      body: formData
    })
    if (res.ok) router.push('/login')
    else alert('Signup failed')
  }

  return (
    <div className="p-8 max-w-md mx-auto">
      <h1 className="text-2xl mb-4">Create an account</h1>
      <form onSubmit={handleSubmit} encType="multipart/form-data">
        <input
          type="email"
          value={email}
          onChange={e => setEmail(e.target.value)}
          placeholder="Email"
          required
          className="block w-full mb-2"
        />
        <input
          type="password"
          value={password}
          onChange={e => setPassword(e.target.value)}
          placeholder="Password"
          required
          className="block w-full mb-2"
        />
        <input
          value={firstName}
          onChange={e => setFirstName(e.target.value)}
          placeholder="First Name"
          required
          className="block w-full mb-2"
        />
        <input
          value={lastName}
          onChange={e => setLastName(e.target.value)}
          placeholder="Last Name"
          required
          className="block w-full mb-2"
        />
        <input
          type="date"
          value={dob}
          onChange={e => setDob(e.target.value)}
          placeholder="Date of Birth"
          required
          className="block w-full mb-2"
        />
        <label className="block mb-2">
          Avatar (optional):
          <input
            type="file"
            accept="image/*"
            onChange={e => setAvatar(e.target.files[0])}
            className="block mt-1"
          />
        </label>
        <input
          value={nickname}
          onChange={e => setNickname(e.target.value)}
          placeholder="Nickname (optional)"
          className="block w-full mb-2"
        />
        <textarea
          value={aboutMe}
          onChange={e => setAboutMe(e.target.value)}
          placeholder="About Me (optional)"
          className="block w-full mb-4"
        />
        <button
          type="submit"
          className="bg-green-500 text-white px-4 py-2"
        >
          Sign Up
        </button>
      </form>
      <p className="mt-4">
        Already have an account?{' '}
        <Link href="/login" className="text-blue-500 underline">
          Log in
        </Link>
      </p>
    </div>
  )
}