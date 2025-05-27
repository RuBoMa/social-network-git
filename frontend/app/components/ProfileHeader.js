'use client'
export default function ProfileHeader({ user }) {
  return (
    <div className="flex items-center gap-4">
      <img src={user.image} alt="Profile" className="w-20 h-20 rounded-full" />
      <div>
        <h2 className="text-xl font-bold">{user.name}</h2>
        <p>{user.followers.length} followers</p>
      </div>
    </div>
  );
}
