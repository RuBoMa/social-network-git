"use client";
import { usePathname } from "next/navigation";
import { useEffect, useState } from "react";
import ChatWindow from "./ChatWindow";
import Author from "./Author";
import GroupAvatar from "./GroupAvatar";
import { addMessageHandler } from "./ws";
import { sendMessage } from "./ws";

export default function ChatBar() {
  const pathname = usePathname();
  const showChatbar = pathname !== "/login" && pathname !== "/signup";
  const [users, setUsers] = useState([]);
  const [openUser, setOpenUser] = useState(null); // keep track of current open chatwindow
  const [currentUser, setCurrentUser] = useState(null);
  const [groups, setGroups] = useState([]);
  const [openGroup, setOpenGroup] = useState(null);

  const [unreadChats, setUnreadChats] = useState({});
  const [unreadGroupChats, setUnreadGroupChats] = useState({});

  useEffect(() => {
    const storedUser = localStorage.getItem("user");
    if (storedUser) {
      try {
        const parsedUser = JSON.parse(storedUser);
        setCurrentUser(parsedUser.user_id);
      } catch (error) {
        console.error("Error parsing user from localStorage:", error);
      }
    } else {
      console.log("No currentUser found in localStorage");
    }
  }, []);

  useEffect(() => {
    const handler = (data) => {
      if (data.type === "interacted_users_response") {
        // console.log("Updated interacted users/groups:", data.users, data.groups); // debug
        setUsers(data.users || []);
        setGroups(data.groups || []);
        // console.log("updated users:", data.users, data.groups); // debug
      }

      if (data.type === "message") {
        if (data.sender?.user_id === currentUser) {
          return; // ignore own message notifications
        }
        setUsers((prevUsers) => {
          const userExists = prevUsers.find((u) => u.user_id === data.sender.user_id);
          return userExists ? prevUsers : [...prevUsers, data.sender];
        });
        if (data.group_id) {
          // Update unread messages for groups
          setUnreadGroupChats((prev) => {
            const updated = {
              ...prev,
              [data.group_id]: (prev[data.group_id] || 0) + 1,
            };
            return updated;
          });
        } else if (data.receiver?.user_id) {
          // Update unread messages for users
          setUnreadChats((prev) => {
            const updated = {
              ...prev,
              [data.sender.user_id]: (prev[data.sender.user_id] || 0) + 1,
            };
            return updated;
          });
        }
      }
    };

    const removeHandler = addMessageHandler(handler);
    return () => removeHandler();
  }, [currentUser]);

  if (!showChatbar) return null;

  return (
    <>
      {openUser && (
        <ChatWindow 
        chatPartner={openUser}
        isGroupChat={false}
        currentUser={currentUser}
        onClose={() => setOpenUser(null)}
        />
      )}
      {openGroup && (
        <ChatWindow
          group={openGroup}
          isGroupChat={true}
          currentUser={currentUser}
          onClose={() => setOpenGroup(null)}
        />
      )}
      <div className="w-1/6 bg-gray-200 p-4 overflow-y-auto">
        <h2 className="text-lg font-bold mb-4">Chats</h2>
        <ul className="space-y-2">
          {users.map((user) => (
            <li key={user.user_id}>
              <button
                onClick={() => {
                  setOpenUser(user);
                  setOpenGroup(null); // Close any open group chat
                  setUnreadChats((prev) => ({
                    ...prev,
                    [user.user_id]: 0,
                  }));
                }}
                className="flex items-center space-x-2 w-full text-left"
              >
                <Author author={user} disableLink={true} size="sm" />

                {/* 🔴 Unread dot */}
                {unreadChats[user.user_id] > 0 && (
                  <span className="ml-auto h-2 w-2 bg-red-600 rounded-full"></span>
                )}
              </button>
            </li>
          ))}
        </ul>
        <ul className="space-y-2">
          {groups.map((group) => (
            <li key={group.group_id}>
              <button
                onClick={() => {
                  setOpenUser(null); // Close any open private chat
                  setOpenGroup(group);
                  setUnreadGroupChats((prev) => ({
                    ...prev,
                    [group.group_id]: 0,
                  }));
                }}
                className="flex items-center space-x-2 w-full text-left"
              >
                <GroupAvatar group={group} disableLink={true} truncateName={true} size="sm" />

                {/* 🔴 Unread dot */}
                {unreadGroupChats[group.group_id] > 0 && (
                  <span className="ml-auto h-2 w-2 bg-red-600 rounded-full"></span>
                )}
              </button>
            </li>
          ))}
        </ul>
      </div>
    </>
  );
}