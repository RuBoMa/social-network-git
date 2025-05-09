// app/components/ProfilePicture.js

export default function ProfilePicture() {
    return (
      <div className="flex justify-center">
        <img
          src="/path-to-profile-pic.jpg"  // Replace with actual path to profile picture
          alt="Profiil"
          className="w-32 h-32 rounded-full border-4 border-gray-300"
        />
      </div>
    );
  }
  