//this page has to have the chat, group list and create post possibility and to view other posts as we have on the main page

// app/page.js
import Chat from "./components/Chat";
import GroupsSidebar from "./components/GroupsSidebar";
import PostFeed from "./components/PostFeed";
import CreatePost from "./components/CreatePost";


export default function HomePage() {
  return (
    <div style={{ display: "flex" }}>
      <aside style={{ width: "20%" }}>
        <GroupsSidebar />
      </aside>

      <main style={{ width: "60%", padding: "0 20px" }}>
        <CreatePost />
        <PostFeed />
      </main>

      <aside style={{ width: "20%" }}>
        <Chat />
      </aside>
    </div>
  );
}

  