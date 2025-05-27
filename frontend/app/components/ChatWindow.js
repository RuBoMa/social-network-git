'use client';
import { useEffect, useState } from 'react';
import { sendMessage, addMessageHandler } from './ws';

export default function ChatWindow({ user, onClose }) {
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState('');

  useEffect(() => {
      console.log("ChatWindow mounted for user:", user);

      // Add message handler specifically for this chat
      const removeHandler = addMessageHandler((data) => {
          console.log("ChatWindow received message:", data);
          
          if (data.type === 'message') {
              // Check if this message belongs to the current chat
              const isChatWithOpenUser =
                  data.sender.user_id === user.user_id || data.receiver.user_id === user.user_id;

              console.log('Is chat with open user:', isChatWithOpenUser);
              console.log('Sender ID:', data.sender.user_id, 'Chat User ID:', user.user_id);
              console.log('Receiver ID:', data.receiver.user_id);

              if (isChatWithOpenUser) {
                  const timeString = new Date().toLocaleTimeString([], {
                      hour: '2-digit',
                      minute: '2-digit',
                      second: '2-digit',
                      hour12: false,
                  });

                  // const isIncoming = data.sender.user_id !== user.user_id;

                  const incomingMsg = {
                      id: Date.now(),
                      reveiverId: data.sender.user_id,
                      receiverName: data.receiver.user_id === user.user_id ? 'Me' : user.nickname,
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
          }
      });

      // Cleanup handler when component unmounts
      return () => {
          console.log("ChatWindow unmounting, removing message handler");
          if (removeHandler) removeHandler();
      };
  }, [user.user_id]);

  function handleSend() {
    if (!input.trim()) return;

    const timeString = new Date().toLocaleTimeString([], {
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',
        hour12: false
    });

    const localMsg = {
        id: Date.now(),
        senderId: user.user_id,
        senderName: user.nickname || 'Me',
        timestamp: timeString,
        content: input
    };

    setMessages((msgs) => [...msgs, localMsg]);

    sendMessage({
        type: 'message',
        content: input,
        receiver: {
            user_id: user.user_id,
        },
    });

    console.log('Message sent to server:', {
        type: 'message',
        content: input,
        receiver: { user_id: user.user_id },
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
              className={`flex flex-col ${msg.senderId === user.user_id ? 'items-end' : 'items-start'}`}
            >
              <div className="flex items-center space-x-2">
                <span className="text-sm font-semibold">{msg.senderName}</span>
                <span className="text-xs text-gray-500">{msg.timestamp}</span>
              </div>
              <div
                className={`mt-1 inline-block bg-gray-200 px-3 py-2 rounded-lg max-w-[50%]
                  ${msg.senderId === user.user_id ? 'rounded-br-none' : 'rounded-bl-none'}`}
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
