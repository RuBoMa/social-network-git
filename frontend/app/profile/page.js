// Replace with this when backend is implemented!:

// import Profile from '../../components/Profile'

// export const metadata = {
//   title: 'Your Profile',
// }

// export default async function ProfilePage() {
//   // call backend and fetch user data, example url!
//   const res = await fetch('http://localhost:8080/api/users/me', {
//     credentials: 'include',
//     cache: 'no-store'
//   })
//   if (!res.ok) {
//     throw new Error('Failed to fetch profile')
//   }
//   const user = await res.json()

//   // render and return current userdata
//   return <Profile user={user} />
// }

import Profile from '../../components/Profile'

export const metadata = {
  title: 'Your Profile (Mock)'
}

export default function ProfilePage() {
  // ==== MOCK DATA ====
  const user = {
    // now served from public/
    avatar_url: '/sky.jpg',
    first_name: 'Jane',
    last_name: 'Doe',
    nickname: 'jdoe',
    followers_count: 4, // maybe don't need, calculate on length of followers?
    followers: [
      { id: 1, avatar_url: '/sky.jpg', first_name: 'Alice', last_name: 'Smith', nickname: 'alice' },
      { id: 2, avatar_url: '/sky.jpg', first_name: 'Bob', last_name: 'Brown', nickname: 'bobby' },
      { id: 3, avatar_url: '/sky.jpg', first_name: 'Carol', last_name: 'Jones', nickname: '' },
      { id: 4, avatar_url: '/sky.jpg', first_name: 'Kira', last_name: 'Kirasson', nickname: 'kiki' }
    ]
  }
  return <Profile user={user} />
}
