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
  const [isPublic, setIsPublic] = useState(false) // New state for profile privacy

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
    formData.append('is_public', isPublic) // Append the privacy setting


    const res = await fetch('http://localhost:8080/api/signup', {
      method: 'POST',
      contentType: 'multipart/form-data',
      credentials: 'include',
      body: formData
    })
    if (res.ok) router.push('/login')
    else alert('Signup failed')
  }

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-50">
      <div className="p-8 max-w-md w-full bg-white rounded-lg shadow-lg">
        <h1 className="text-2xl mb-4 text-center">Create an account</h1>
        <form onSubmit={handleSubmit} encType="multipart/form-data">
          <label className="block mb-2">
            Email <span className="text-red-500">*</span>
            <input
              type="email"
              value={email}
              onChange={e => setEmail(e.target.value)}
              required
              className="block w-full mt-1 mb-2 border rounded px-2 py-1"
            />
          </label>

          <label className="block mb-2">
            Password <span className="text-red-500">*</span>
            <input
              type="password"
              value={password}
              onChange={e => setPassword(e.target.value)}
              required
              className="block w-full mt-1 mb-2 border rounded px-2 py-1"
            />
          </label>

          <label className="block mb-2">
            First Name <span className="text-red-500">*</span>
            <input
              value={firstName}
              onChange={e => setFirstName(e.target.value)}
              required
              className="block w-full mt-1 mb-2 border rounded px-2 py-1"
            />
          </label>

          <label className="block mb-2">
            Last Name <span className="text-red-500">*</span>
            <input
              value={lastName}
              onChange={e => setLastName(e.target.value)}
              required
              className="block w-full mt-1 mb-2 border rounded px-2 py-1"
            />
          </label>

          <label className="block mb-4">
            Date of Birth <span className="text-red-500">*</span>
            <input
              type="date"
              value={dob}
              onChange={e => setDob(e.target.value)}
              required
              className="block w-full mt-1 border rounded px-2 py-1"
            />
          </label>

          <label className="block mb-4">
            Avatar
            <div className="flex items-center mt-1">
              <label
                htmlFor="avatar"
                className="cursor-pointer bg-gray-200 text-gray-600 px-1 py-1 text-sm rounded border hover:bg-gray-300"
              >
                Choose File
              </label>
              <span className="ml-3 text-sm text-gray-600">
                {avatar?.name || 'No file chosen'}
              </span>
            </div>
            <input
              id="avatar"
              type="file"
              accept="image/*"
              onChange={e => setAvatar(e.target.files[0])}
              className="hidden"
            />
          </label>

          <label className="block mb-2">
            Nickname
            <input
              value={nickname}
              onChange={e => setNickname(e.target.value)}
              className="block w-full mt-1 mb-2 border rounded px-2 py-1"
            />
          </label>

          <label className="block mb-4">
            About Me
            <textarea
              value={aboutMe}
              onChange={e => setAboutMe(e.target.value)}
              className="block w-full mt-1 border rounded px-2 py-1"
            />
          </label>

          <label className="block mb-4">
            Profile Privacy
            <div className="flex items-center mt-1">
              <div
                className={`relative w-12 h-6 bg-gray-300 rounded-full cursor-pointer transition-colors ${
                  isPublic ? 'bg-blue-500' : 'bg-gray-300'
                }`}
                onClick={() => setIsPublic(!isPublic)}
              >
                <div
                  className={`absolute top-0.5 left-0.5 w-5 h-5 bg-white rounded-full shadow-md transition-transform ${
                    isPublic ? 'translate-x-6' : 'translate-x-0'
                  }`}
                ></div>
              </div>
              <span className="ml-3">{isPublic ? 'Public' : 'Private'}</span>
            </div>
          </label>

          <button
            type="submit"
            className="w-full py-2 bg-blue-500 text-white rounded disabled:bg-blue-300 disabled:cursor-not-allowed"
          >
            Sign up
          </button>
        </form>
        <p className="mt-4 text-center">
          Already have an account?{' '}
          <Link href="/login" className="text-blue-500 underline">
            Sign in
          </Link>
        </p>
      </div>
    </div>
  )
}