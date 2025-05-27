'use client';
import { useEffect, useState } from 'react';
import EmojiPicker from 'emoji-picker-react';
import { sendMessage, addMessageHandler } from './ws';
import Author from './Author';

export default function ChatWindow({ chatPartner, onClose }) {
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState('');
  const [showEmoji, setShowEmoji] = useState(false);

  function handleEmojiClick(emojiData) {
  setInput(input + emojiData.emoji);
  setShowEmoji(false);
  }

  useEffect(() => {
    console.log("ChatWindow mounted for user:", chatPartner);

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
        
        // const isIncoming = data.sender.user_id !== user.user_id;
        // If want to show nickname, fetch own information from local storage
        let nickname = "Me";
        // Checking if the message is from chat partner
        if (data.sender.user_id === chatPartner.user_id) {
          nickname = chatPartner.nickname || chatPartner.first_name;
        }
        
        const incomingMsg = {
          id: Date.now(),
          senderId: data.receiver.user_id,
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
    }
  });


  // Cleanup handler when component unmounts
  return () => {
    console.log("ChatWindow unmounting, removing message handler");
    if (removeHandler) removeHandler();
    };
  }, [chatPartner.user_id]);

  function handleSend() {
    if (!input.trim()) return;

    sendMessage({
      type: 'message',
      content: input,
      receiver: {
        user_id: chatPartner.user_id,
      },
    });

    console.log('Message sent to server:', {
      type: 'message',
      content: input,
      receiver: { user_id: chatPartner.user_id },
    });

    setInput('');
  }


  return (
    <div className="fixed bottom-4 right-4 z-50 border border-gray-300 rounded-lg shadow-lg">
      <div className="bg-white w-80 h-[40vh] p-3 flex flex-col rounded-lg shadow-lg overflow-hidden">
        <header className="flex justify-between items-center p-2 border-b border-gray-300">
         <Author author={chatPartner} disableLink={true} size="sm" />
          <button onClick={onClose} className="text-xl leading-none">&times;</button>
        </header>

        <div className="flex-1 p-2 flex flex-col space-y-2 overflow-y-auto">
          {messages.map((msg) => (
            <div
              key={msg.id}
              className={`flex flex-col ${msg.senderId === chatPartner.user_id ? 'items-end' : 'items-start'}`}
            >
              <div className="flex items-center space-x-2">
                <span className="text-sm font-semibold">{msg.senderName}</span>
                <span className="text-xs text-gray-500">{msg.timestamp}</span>
              </div>
              <div
                className={`mt-1 inline-block bg-gray-200 px-3 py-2 rounded-lg max-w-[50%]
                  ${msg.senderId === chatPartner.user_id ? 'rounded-br-none' : 'rounded-bl-none'}`}
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
          >😊</button>
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