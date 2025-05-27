'use client'
import React from 'react'

export default function ImageUploadPreview({ imageFile, setImageFile }) {
  if (!imageFile) return null

  const handleRemove = () => {
    setImageFile(null)
  }

  return (
    <div className="mt-2 flex items-center space-x-2">
      <button
        type="button"
        onClick={handleRemove}
        className="text-red-500 hover:text-red-700"
        aria-label="Remove image"
      >
        ✕
      </button>
      <img
        src={URL.createObjectURL(imageFile)}
        alt="Preview"
        className="h-10 w-10 rounded object-cover"
        onLoad={() => URL.revokeObjectURL(imageFile)}
      />
    </div>
  )
}
