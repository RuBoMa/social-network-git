'use client'
import { useEffect, useState } from 'react'

// export function PostFeed() {
//   const [posts, setPosts] = useState([])

//   useEffect(() => {
//     async function fetchPosts() {
//       const res = await fetch('http://localhost:8080/api/feed', {
//         credentials: 'include', 
//         headers: {
//           'Accept': 'application/json' //telling the server we want JSON
//         }
//       })

//       console.log('Response status:', res) // Log the response status

//       if (res.ok) {
//         const data = await res.json()
//         setPosts(data)
//       } else {
//         console.error('Failed to load posts')
//       }
//     }

//     fetchPosts()
//   }, [])

//   return (
//     <div>
//       {posts.map((post, i) => (
//         <div key={i} className="post">
//           <p><strong>{post.author}</strong>: {post.content}</p>
//         </div>
//       ))}
//     </div>
//   )
// }

export function CreatePost() {
  //state to store different post data
  const [postTitle, setPostTitle] = useState('');
  const [postContent, setContent] = useState('')
  const [privacy, setPostPrivacy] = useState('public');
  const [postImage, setPostImage] = useState(null);

  //handling the post submission
  async function handlePost(e) {
    e.preventDefault() //preventing page reload on form submit

    //creating a FormData object to send the post data
    const formData = new FormData();
    formData.append('post_title', postTitle);
    formData.append('post_content', postContent);
    formData.append('privacy', privacy);

    //if an image is selected, append it to the form data
    if (postImage) {
      formData.append('image', postImage); 
    }

    //sending the form data to the backend API
    const res = await fetch('http://localhost:8080/api/create-post', {
      method: 'POST',
      credentials: 'include', 
      body: formData,
    })

    //if the response is ok, clear the form fields
    if (res.ok) {
      setContent('')
      setPostTitle('')
      setPostImage(null)
      // refresh post list or show success
    } else {
      alert('Failed to post')
    }
  }

  return (
    <form className="max-w-md mx-auto p-6 bg-white rounded shadow" onSubmit={handlePost}>
      <label className="block mb-4">
      <span className="block text-sm font-medium text-gray-700">Title:</span>
        <input
          type="text"
          value={postTitle}
          onChange={(e) => setPostTitle(e.target.value)}
          required
          className="mt-1 block w-full border border-gray-300 rounded p-2"
        />
      </label>

      <label className="block mb-4">
      <span className="block text-sm font-medium text-gray-700">Content:</span>
        <textarea
          value={postContent}
          onChange={(e) => setContent(e.target.value)}
          required
          className="mt-1 block w-full border border-gray-300 rounded p-2"
        />
      </label>
    
      <label className="block mb-4">
  <span className="block text-sm font-medium text-gray-700 mb-2">Privacy:</span>

  <div className="flex flex-col space-y-2">
    <label className="inline-flex items-center">
      <input
        type="radio"
        name="privacy"
        value="public"
        checked={privacy === 'public'}
        onChange={(e) => setPostPrivacy(e.target.value)}
        className="form-radio text-blue-600"
      />
      <span className="ml-2">Public</span>
    </label>

    <label className="inline-flex items-center">
      <input
        type="radio"
        name="privacy"
        value="followers"
        checked={privacy === 'followers'}
        onChange={(e) => setPostPrivacy(e.target.value)}
        className="form-radio text-blue-600"
      />
      <span className="ml-2">Followers</span>
    </label>

    <label className="inline-flex items-center">
      <input
        type="radio"
        name="privacy"
        value="custom"
        checked={privacy === 'custom'}
        onChange={(e) => setPostPrivacy(e.target.value)}
        className="form-radio text-blue-600"
      />
      <span className="ml-2">Custom</span>
    </label>
  </div>
</label>
      <label className="block mb-4">
      <span className="block text-sm font-medium text-gray-700">Image:</span>
        <input
          type="file"
          accept="image/*"
          onChange={(e) => setPostImage(e.target.files[0])}
          className="mt-1 block w-full"
        />
      </label>
      <button type="submit" className="w-full bg-blue-600 text-white py-2 px-4 rounded hover:bg-blue-700">Create Post</button>
    </form>
  );
}

export default function mainPage() {
  return (
    <div>
      <CreatePost />
      {/* <PostFeed /> */}
    </div>
  )
}