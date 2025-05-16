'use client'
import { useState } from 'react'
import CreatePost from './components/CreatePost'
import { PostFeed } from './components/PostFeed'

export default function mainPage() {
  const [reloadPage, setReloadPage] = useState(false)

  return (
    <div className="max-w-5xl mx-auto px-4">
      <CreatePost onSuccess={() => setReloadPage(prev => !prev)} />
      <PostFeed reloadTrigger={reloadPage} />
    </div>
  )
}
