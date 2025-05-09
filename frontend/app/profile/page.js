import ProfileInfo from "@/components/ProfileInfo";
import ProfilePicture from "@/components/ProfilePicture";
import UserPosts from "@/components/UserPosts";

export default function ProfilePage() {
    return (
      <div className="p-4 space-y-4">
        <ProfilePicture />
        <ProfileInfo />
        <UserPosts />
      </div>
    );
  }