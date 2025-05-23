'use client'
import { useRouter } from 'next/navigation'

export default function BackButton({ label = "Back" }) {
  const router = useRouter()
  return (
    <button
      onClick={() => router.back()}
      className="mb-4 px-3 py-1 rounded bg-gray-200 hover:bg-gray-300 text-sm font-semibold"
    >
      {label}
    </button>
  )
}