'use client';
import { useEffect, useState } from 'react';

export default function ChatWindow({ user, onClose }) {
  // Holds the chat history
  const [messages, setMessages] = useState([
    // get messages ---
  ]);
  const [input, setInput] = useState('');

  useEffect(() => {
    // load user messages ---
  }, [user.user_id]);

  // update handeSend to actually send the message, user id, user name etc
  function handleSend() {
    if (!input.trim()) return;
    const now = new Date();
    const timeString = now.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit', hour12: false });
    setMessages(msgs => [
      ...msgs,
      {
        id: Date.now(),
        senderId: 'me',
        senderName: 'Me',
        timestamp: timeString,
        content: input
      }
    ]);
    setInput('');
    // send message to server ---
  }

return (
    <div className="fixed bottom-4 right-4 z-50 border border-gray-300 rounded-lg shadow-lg">
      <div className="bg-white w-80 h-[40vh] flex flex-col rounded-lg shadow-lg overflow-hidden">
        <header className="flex justify-between items-center p-2 border-b border-gray-300">
          <h2 className="font-bold">
            {user.first_name} {user.last_name}
          </h2>
          <button onClick={onClose} className="text-xl leading-none">&times;</button>
        </header>

        <div className="flex-1 p-2 flex flex-col space-y-2 overflow-y-auto">
          {messages.map(msg => (
            <div
              key={msg.id}
              // check for actual user, if sender of reciever and update alignment
              className={`flex flex-col ${msg.senderId === 'me' ? 'items-end' : 'items-start'}`}
            >
              <div className="flex items-center space-x-2">
                <span className="text-sm font-semibold">{msg.senderName}</span>
                <span className="text-xs text-gray-500">{msg.timestamp}</span>
              </div>
              <div
              // same thing here, check for actual user id
                className={`mt-1 inline-block bg-gray-200 px-3 py-2 rounded-lg max-w-[50%]
                  ${msg.senderId === 'me'
                    ? 'rounded-br-none' // all rounded except bottom-right
                    : 'rounded-bl-none' // all rounded except bottom-left
                  }`}
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
            onChange={e => setInput(e.target.value)}
            onKeyDown={e => {
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
