'use client';
import { useEffect, useState, useRef } from 'react';
import EmojiPicker from 'emoji-picker-react';
import { sendMessage, addMessageHandler } from './ws';
import Author from './Author';

export default function ChatWindow({ chatPartner, onClose }) {
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState('');
  const [showEmoji, setShowEmoji] = useState(false);

  const PAGE_SIZE = 10;
  const [page, setPage] = useState(0);
  const messagesRef = useRef(null);

  function handleEmojiClick(emojiData) {
    setInput(input + emojiData.emoji);
    setShowEmoji(false);
  }

  useEffect(() => {
    console.log("ChatWindow mounted for user:", chatPartner);
    setPage(0); // Reset page when chat partner changes
    setMessages([]); // Clear messages when chat partner changes

    sendMessage({
      type: 'chat',
      receiver: { user_id: chatPartner.user_id },
      page: 0,
      page_size: PAGE_SIZE,
    });

    // Add message handler specifically for this chat
    const removeHandler = addMessageHandler((data) => {
      console.log("ChatWindow received message:", data);

      if (data.type === 'message') {
        // Check if this message belongs to the current chat
        const isChatWithOpenUser =
          data.sender.user_id === chatPartner.user_id || data.receiver.user_id === chatPartner.user_id;

        console.log('Is chat with open user:', isChatWithOpenUser);
        console.log('Sender ID:', data.sender.user_id, 'Chat User ID:', chatPartner.user_id);
        console.log('Receiver ID:', data.receiver.user_id);

        if (isChatWithOpenUser) {
          const timeString = new Date().toLocaleTimeString([], {
            hour: '2-digit',
            minute: '2-digit',
            second: '2-digit',
            hour12: false,
          });

          // If want to show nickname, fetch own information from local storage
          let nickname = "Me";
          // Checking if the message is from chat partner
          if (data.sender.user_id === chatPartner.user_id) {
            nickname = chatPartner.nickname || chatPartner.first_name;
          }

          const incomingMsg = {
            id: Date.now(),
            senderId: data.sender.user_id,
            senderName: nickname,
            timestamp: timeString,
            content: data.content,
          };

          console.log('Adding message to chat:', incomingMsg);

          setMessages((msgs) => {
            const newMessages = [...msgs, incomingMsg];
            console.log('Updated messages array:', newMessages);
            return newMessages;
          });
        } else {
          console.log('Message filtered out - not for this chat');
        }
      } else if (data.type === 'chat') {
        console.log("Number of history messages received:", data.history.length);

        console.log('Chat history received:', data);

        const formattedMessages = data.history
          .sort((a, b) => new Date(a.created_at) - new Date(b.created_at)) // Sort messages by creation date
          .map((msg, index) => ({
            id: msg.id || `${msg.sender.user_id}-${msg.created_at}-${index}`,
            senderId: msg.sender.user_id,
            senderName: msg.sender.user_id === chatPartner.user_id
              ? (chatPartner.nickname || chatPartner.first_name)
              : "Me",
            timestamp: new Date(msg.created_at).toLocaleTimeString([], {
              hour: '2-digit',
              minute: '2-digit',
              second: '2-digit',
              hour12: false,
            }),
            content: msg.content,
          }));

          if (page === 0) {
            // First page: replace messages and scroll to bottom
            setMessages(formattedMessages);
          } else {
            // Older pages: prepend messages and maintain scroll position
            const el = messagesRef.current;
            const oldScrollHeight = el?.scrollHeight || 0;
  
            setMessages((prev) => [...formattedMessages, ...prev]);
  
            setTimeout(() => {
              if (el) {
                const newScrollHeight = el.scrollHeight;
                el.scrollTop = newScrollHeight - oldScrollHeight;
              }
            }, 0);
          }
        }
      });


    // Cleanup handler when component unmounts
    return () => {
      console.log("ChatWindow unmounting, removing message handler");
      if (removeHandler) removeHandler();
    };
  }, [chatPartner.user_id]);
  useEffect(() => {
    if (page === 0 && messagesRef.current) {
      messagesRef.current.scrollTop = messagesRef.current.scrollHeight;
    }
  }, [messages, page]);

  // Scroll event to load older messages when scrolled to top
  useEffect(() => {
    const el = messagesRef.current;
    if (!el) return;

    function onScroll() {
      if (el.scrollTop === 0 && messages.length >= PAGE_SIZE * (page + 1)) {
        setPage((p) => p + 1);
        // Request older messages for next page
        sendMessage({
          type: 'chat',
          receiver: { user_id: chatPartner.user_id },
          page: page + 1,
          pageSize: PAGE_SIZE,
        });
      }
    }

    el.addEventListener('scroll', onScroll);
    return () => el.removeEventListener('scroll', onScroll);
  }, [messages, page, chatPartner.user_id]);

  function handleSend() {
    if (!input.trim()) return;

    sendMessage({
      type: 'message',
      content: input,
      receiver: {
        user_id: chatPartner.user_id,
      },
    });
    setInput('');
  }


  return (
    <div className="fixed bottom-4 right-4 z-50 border border-gray-300 rounded-lg shadow-lg">
      <div className="bg-white w-80 h-[40vh] p-3 flex flex-col rounded-lg shadow-lg overflow-hidden">
        <header className="flex justify-between items-center p-2 border-b border-gray-300">
          <Author author={chatPartner} disableLink={true} size="sm" />
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
                msg.senderId === chatPartner.user_id
                  ? 'items-start'
                  : 'items-end'
              }`}
            >
              <div className="flex items-center space-x-2">
                <span className="text-sm font-semibold">{msg.senderName}</span>
                <span className="text-xs text-gray-500">{msg.timestamp}</span>
              </div>
              <div
                className={`mt-1 inline-block bg-gray-200 px-3 py-2 rounded-lg max-w-[50%]
                  ${
                    msg.senderId === chatPartner.user_id
                      ? 'rounded-br-none'
                      : 'rounded-bl-none'
                  }`}
              >
                {msg.content}
              </div>
            </div>
          ))}
        </div>

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
            onChange={(e) => setInput(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === 'Enter') {
                e.preventDefault();
                handleSend();
              }
            }}
            className="flex-1 border border-gray-300 rounded px-2 py-1 mr-2"
            placeholder="Type a message…"
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