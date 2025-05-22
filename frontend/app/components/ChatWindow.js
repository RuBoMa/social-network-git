'use client';
import { useEffect, useState } from 'react';
import { initWebSocket, sendMessage, closeWebSocket } from './ws';

export default function ChatWindow({ user, onClose }) {
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState('');

  useEffect(() => {
    initWebSocket(user.user_id, (msg) => {
      if (msg.type === 'messageBE' && msg.sender_id === user.user_id) {
        setMessages((prev) => [...prev, {
          id: Date.now(),
          senderId: msg.sender_id,
          senderName: msg.sender_name,
          timestamp: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit', hour12: false }),
          content: msg.content
        }]);
      }
  
      // Handle history message from backend
      if (msg.type === 'chatBE' && Array.isArray(msg.history)) {
        const historyMsgs = msg.history.map((m) => ({
          id: m.id || Date.now(),  
          senderId: m.sender?.user_id,
          senderName: m.sender?.nickname || 'Unknown',
          timestamp: new Date(m.created_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit', hour12: false }),
          content: m.content,
        }));
  
        setMessages(historyMsgs.reverse()); // Reverse to show oldest first
      }
    });
  
    // After websocket connected, request chat history
    sendMessage({
      type: 'chatBE',
      sender: { user_id: user.user_id },
      receiver: { user_id: user.user_id }, // Replace with the chat partner's id here
      group_id: 0,
    });
  
    return () => {
      closeWebSocket();
    };
  }, [user.user_id]);
  

  function handleSend() {
    if (!input.trim()) return;

    const timeString = new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit', hour12: false });

    const localMsg = {
      id: Date.now(),
      senderId: 'me',
      senderName: 'Me',
      timestamp: timeString,
      content: input
    };
    setMessages((msgs) => [...msgs, localMsg]);

    sendMessage({
      type: 'messageBE',
      receiver_id: user.user_id,
      content: input,
    });

    setInput('');
  }

  return (
    <div className="fixed bottom-4 right-4 z-50 border border-gray-300 rounded-lg shadow-lg">
      <div className="bg-white w-80 h-[40vh] flex flex-col rounded-lg shadow-lg overflow-hidden">
        <header className="flex justify-between items-center p-2 border-b border-gray-300">
          <h2 className="font-bold">{user.nickname}</h2>
          <button onClick={onClose} className="text-xl leading-none">&times;</button>
        </header>

        <div className="flex-1 p-2 flex flex-col space-y-2 overflow-y-auto">
          {messages.map((msg) => (
            <div
              key={msg.id}
              className={`flex flex-col ${msg.senderId === 'me' ? 'items-end' : 'items-start'}`}
            >
              <div className="flex items-center space-x-2">
                <span className="text-sm font-semibold">{msg.senderName}</span>
                <span className="text-xs text-gray-500">{msg.timestamp}</span>
              </div>
              <div
                className={`mt-1 inline-block bg-gray-200 px-3 py-2 rounded-lg max-w-[50%]
                  ${msg.senderId === 'me' ? 'rounded-br-none' : 'rounded-bl-none'}`}
              >
                {msg.content}
              </div>
            </div>
          ))}
        </div>

        <footer className="p-4 border-t border-gray-300 flex">
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
