'use client';
import { useEffect, useState, useRef } from 'react';
import EmojiPicker from 'emoji-picker-react';
import { sendMessage, addMessageHandler } from './ws';
import Author from './Author';

export default function ChatWindow({ chatPartner, onClose }) {
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState('');
  const [showEmoji, setShowEmoji] = useState(false);
  const messagesRef = useRef(null);

  function handleEmojiClick(emojiData) {
    setInput(input + emojiData.emoji);
    setShowEmoji(false);
  }

  useEffect(() => {
    console.log("ChatWindow mounted for user:", chatPartner);
    setMessages([]); // Reset messages when chat partner changes

    // Request entire chat history
    sendMessage({
      type: 'chat',
      receiver: { user_id: chatPartner.user_id },
      page: 0,
      page_size: 1000, // High number to fetch all messages at once
    });

    const removeHandler = addMessageHandler((data) => {
      if (data.type === 'message') {
        const isChatWithOpenUser =
          data.sender.user_id === chatPartner.user_id || data.receiver.user_id === chatPartner.user_id;

        if (isChatWithOpenUser) {
          const timeString = new Date().toLocaleTimeString([], {
            hour: '2-digit',
            minute: '2-digit',
            second: '2-digit',
            hour12: false,
          });

          let nickname = "Me";
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

          setMessages((msgs) => [...msgs, incomingMsg]);
        }
      } else if (data.type === 'chat') {
        const formattedMessages = data.history
          .sort((a, b) => new Date(a.created_at) - new Date(b.created_at))
          .map((msg, index) => ({
            id: msg.id || `${msg.sender.user_id}-${msg.created_at}-${index}`,
            senderId: msg.sender.user_id,
            senderName:
              msg.sender.user_id === chatPartner.user_id
                ? chatPartner.nickname || chatPartner.first_name
                : 'Me',
            timestamp: new Date(msg.created_at).toLocaleTimeString([], {
              hour: '2-digit',
              minute: '2-digit',
              second: '2-digit',
              hour12: false,
            }),
            content: msg.content,
          }));

        setMessages(formattedMessages);
      }
    });

    return () => {
      if (removeHandler) removeHandler();
    };
  }, [chatPartner.user_id]);

  useEffect(() => {
    if (messagesRef.current) {
      messagesRef.current.scrollTop = messagesRef.current.scrollHeight;
    }
  }, [messages]);

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
                msg.senderId === chatPartner.user_id ? 'items-start' : 'items-end'
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
