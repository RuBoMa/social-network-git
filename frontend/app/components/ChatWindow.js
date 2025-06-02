"use client";
import { useEffect, useState, useRef } from "react";
import EmojiPicker from "emoji-picker-react";
import { sendMessage, addMessageHandler } from "./ws";
import Author from "./Author";
import GroupAvatar from "./GroupAvatar";

export default function ChatWindow({
  chatPartner,
  group,
  onClose,
  isGroupChat,
  currentUser,
}) {
  const [messages, setMessages] = useState([]);
  const [isTyping, setIsTyping] = useState(false);
  const [typingUser, setTypingUser] = useState(null);
  const [input, setInput] = useState("");
  const [showEmoji, setShowEmoji] = useState(false);
  const typingTimeoutRef = useRef(null);
  const lastSentTypingRef = useRef(0);
  const messagesRef = useRef(null);

  // const currentUser = JSON.parse(localStorage.getItem("user"));

  if (isGroupChat) {
    if (!group || !group.group_id) return null;
  } else {
    if (!chatPartner || !chatPartner.user_id) return null;
  }

  function handleEmojiClick(emojiData) {
    setInput(input + emojiData.emoji);
    setShowEmoji(false);
  }

  useEffect(() => {
    console.log("ChatWindow mounted for user:", chatPartner);
    setMessages([]); // Reset messages when chat partner changes

    if (isGroupChat) {
      // Fetch group chat history
      setMessages([]); // Reset messages for group chat
      sendMessage({
        type: "chat",
        group_id: group.group_id,
      });
    } else {
      // Fetch private chat history
      setMessages([]); // Reset messages for private chat
      sendMessage({
        type: "chat",
        receiver: { user_id: chatPartner.user_id },
      });
    }

    const removeHandler = addMessageHandler((data) => {
      console.log("Received data:", data);
      if (data.type === "message") {
        console.log("Processing filtered message:", data);
        const timeString = new Date().toLocaleTimeString([], {
          year: "numeric",
          month: "2-digit",
          day: "2-digit",
          hour: "2-digit",
          minute: "2-digit",
          second: "2-digit",
          hour12: false,
        });

        let nickname =
          data.sender.nickname || data.sender.first_name || "Unknown User";

        const incomingMsg = {
          id: Date.now(),
          senderId: data.sender.user_id,
          senderName: nickname,
          timestamp: timeString,
          content: data.content,
        };

        setMessages((msgs) => [...msgs, incomingMsg]);
      } else if (data.type === "chat") {
        const formattedMessages = data.history.map((msg, index) => ({
          id: msg.id || `${msg.sender.user_id}-${msg.created_at}-${index}`,
          senderId: msg.sender.user_id,
          senderName:
            msg.sender.nickname || msg.sender.first_name || "Unknown User",
          timestamp: new Date(msg.created_at).toLocaleTimeString([], {
            year: "numeric",
            month: "2-digit",
            day: "2-digit",
            hour: "2-digit",
            minute: "2-digit",
            second: "2-digit",
            hour12: false,
          }),
          content: msg.content,
        }));

        setMessages(formattedMessages);
      } else if (data.type === "typing") {
        if (isGroupChat) {
          // Handle typing indicator for group chats
          if (data.sender.user_id !== currentUser.user_id) {
            setIsTyping(true);
            setTypingUser(data.sender);
            console.log("Typing user details:", data.sender);
          }
        } else {
          // Handle typing indicator for private chats
          if (data.sender.user_id === chatPartner.user_id) {
            setIsTyping(true);
            setTypingUser(data.sender);
            console.log("Typing user details:", data.sender);
          }
        }
      } else if (data.type === "stop_typing") {
        if (isGroupChat) {
          // Handle stop typing indicator for group chats
          if (data.sender.user_id !== currentUser.user_id) {
            setIsTyping(false);
            setTypingUser(null);
            console.log("Typing indicator OFF for group:", group.group_id);
          }
        } else {
          // Handle stop typing indicator for private chats
          if (data.sender.user_id === chatPartner.user_id) {
            setIsTyping(false);
            setTypingUser(null);
            console.log("Typing indicator OFF for:", chatPartner);
          }
        }
      }
    });

    return () => {
      if (removeHandler) removeHandler();
    };
  }, [isGroupChat, group?.group_id, chatPartner?.user_id]);

  useEffect(() => {
    if (messagesRef.current) {
      messagesRef.current.scrollTop = messagesRef.current.scrollHeight;
    }
  }, [messages]);

  function handleSend() {
    if (!input.trim()) return;
    console.log("Sending message:", input);
    console.log("Chat partner:", chatPartner);
    console.log("Group:", group);
    console.log("Is group chat:", isGroupChat);

    if (isGroupChat) {
      sendMessage({
        type: "message",
        content: input,
        group_id: group.group_id,
      });
    } else {
      sendMessage({
        type: "message",
        content: input,
        receiver: {
          user_id: chatPartner.user_id,
        },
      });
    }
    if (isGroupChat) {
      sendMessage({
        type: "stopTypingBE",
        group_id: group.group_id,
      });
    } else {
      sendMessage({
        type: "stopTypingBE",
        receiver: { user_id: chatPartner.user_id },
      });
    }

    setInput("");
  }

  return (
    <div className="fixed bottom-4 right-4 z-50 border border-gray-300 rounded-lg shadow-lg">
      <div className="bg-white w-80 h-[40vh] p-3 flex flex-col rounded-lg shadow-lg overflow-hidden">
        <header className="flex justify-between items-center p-2 border-b border-gray-300">
          {isGroupChat ? (
            <GroupAvatar group={group} disableLink={true} size="sm" />
          ) : (
            <Author author={chatPartner} disableLink={true} size="sm" />
          )}
          <button onClick={onClose} className="text-xl leading-none">
            &times;
          </button>
        </header>

        <div
          ref={messagesRef}
          className="flex-1 p-2 flex flex-col space-y-2 overflow-y-auto"
        >
          {messages.map((msg) => (
            <div
              key={msg.id}
              className={`flex flex-col ${
                msg.senderId === currentUser ? "items-end" : "items-start"
              }`}
            >
              <div className="flex items-center space-x-2">
                <span className="text-sm font-semibold">
                  {msg.senderName || "Unknown"}
                </span>
                <span className="text-xs text-gray-500">
                  {msg.timestamp || ""}
                </span>
              </div>
              <div
                className={`mt-1 inline-block bg-gray-200 px-3 py-2 rounded-lg max-w-[50%]
                  ${
                    msg.senderId === currentUser
                      ? "rounded-br-none"
                      : "rounded-bl-none"
                  }`}
              >
                {msg.content}
              </div>
            </div>
          ))}
        </div>
        {isTyping && (
          <div className="text-gray-500 text-sm italic mt-2 flex items-center">
            <span>
              {isGroupChat
                ? `${
                    typingUser?.nickname ||
                    typingUser?.first_name ||
                    typingUser?.user_id ||
                    "Someone"
                  } is typing...`
                : `${
                    chatPartner?.nickname || chatPartner?.first_name || "User"
                  } is typing`}
            </span>
            <span className="ml-1 flex space-x-1">
              <span className="w-1.5 h-1.5 bg-gray-500 rounded-full animate-dots"></span>
              <span className="w-1.5 h-1.5 bg-gray-500 rounded-full animate-dots delay-150"></span>
              <span className="w-1.5 h-1.5 bg-gray-500 rounded-full animate-dots delay-300"></span>
            </span>
          </div>
        )}

        <footer className="p-4 border-t border-gray-300 flex items-center">
          <button
            type="button"
            className="mr-2 text-xl"
            onClick={() => setShowEmoji((v) => !v)}
          >
            😊
          </button>
          {showEmoji && (
            <div className="absolute bottom-16 right-4 z-50">
              <EmojiPicker onEmojiClick={handleEmojiClick} />
            </div>
          )}
          <input
            type="text"
            value={input}
            onChange={(e) => {
              setInput(e.target.value);

              const now = Date.now();
              if (now - lastSentTypingRef.current > 1000) {
                if (isGroupChat) {
                  sendMessage({
                    type: "typingBE",
                    group_id: group.group_id, // Send typing event to group
                  });
                } else {
                  sendMessage({
                    type: "typingBE",
                    receiver: { user_id: chatPartner.user_id }, // Send typing event to private chat
                  });
                }
                lastSentTypingRef.current = now;
              }

              if (typingTimeoutRef.current)
                clearTimeout(typingTimeoutRef.current);
              typingTimeoutRef.current = setTimeout(() => {
                if (isGroupChat) {
                  sendMessage({
                    type: "stopTypingBE",
                    group_id: group.group_id, // Stop typing event for group
                  });
                } else {
                  sendMessage({
                    type: "stopTypingBE",
                    receiver: { user_id: chatPartner.user_id }, // Stop typing event for private chat
                  });
                }
              }, 1500); // Typing indicator stays for 1.5 seconds after typing stops
            }}
            onKeyDown={(e) => {
              if (e.key === "Enter") {
                e.preventDefault();
                handleSend();
              }
            }}
            className="flex-1 border border-gray-300 rounded px-2 py-1 mr-2"
            placeholder="Type a message…"
            maxLength={200}
          />
          <button
            onClick={handleSend}
            className="bg-blue-600 text-white px-3 py-1 rounded"
          >
            Send
          </button>
        </footer>
      </div>
    </div>
  );
}
