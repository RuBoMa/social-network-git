'use client'
import { useRouter } from 'next/navigation'

export default function BackButton({ label = "Back" }) {
  const router = useRouter()
  return (
    <button
      onClick={() => router.back()}
      className="bg-sky-600/60 hover:bg-sky-700/60 text-white font-bold rounded text-md mb-4 p-2 px-4 cursor-pointer"
    >
      {label}
    </button>
  )
}