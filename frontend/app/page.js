"use client";
import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import CreatePost from "./components/CreatePost";
import { PostFeed } from "./components/PostFeed";

export default function mainPage() {
  const [reloadPage, setReloadPage] = useState(false);
  const [isAuthenticated, setIsAuthenticated] = useState(null);
  const router = useRouter();

  useEffect(() => {
    const token = localStorage.getItem("token");
    if (!token) {
      router.replace("/login"); // Redirect to login if not authenticated
    } else {
      setIsAuthenticated(true); // User is authenticated
    }
  }, [router]);

  if (isAuthenticated === null) {
    // Show a loading spinner or blank screen while checking authentication
    return <div>Loading...</div>;
  }

  return (
    <div className="mx-auto px-4 space-y-4">
      <CreatePost onSuccess={() => setReloadPage((prev) => !prev)} />
      <PostFeed reloadTrigger={reloadPage} />
    </div>
  );
}
