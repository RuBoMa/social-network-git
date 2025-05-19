'use client'
import React from 'react'


export default function Author({ author, size = "md" }) {

    if (!author) return null
  
    const sizeClasses = {
      xs: "w-4 h-4 text-xs",
      sm: "w-6 h-6 text-sm",
      md: "w-8 h-8 text-base",
      lg: "w-10 h-10 text-lg",
    }

    const imageUrl = author.avatar_path && author.avatar_path.trim() !== ""
    ? `http://localhost:8080${author.avatar_path}`
    : '/avatar.png'
  
    return (
      <div className="flex items-center gap-2">
          <img
            src={imageUrl}
            
            alt="Author"
            className={`rounded-full ${sizeClasses[size] || sizeClasses.md}`}
          />

        <span className="font-medium">{author.nickname || author.first_name || 'Unknown Author'}</span>
      </div>
    )
  }