'use client';
import { useEffect, useState } from 'react';
import EmojiPicker from 'emoji-picker-react';
import { sendMessage, addMessageHandler } from './ws';
import Author from './Author';
import GroupAvatar from './GroupAvatar';

export default function ChatWindow({ user, group, onClose, isGroupChat }) {
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState('');
  const [showEmoji, setShowEmoji] = useState(false);

  //user from localStorage
  const myUserId = Number(localStorage.getItem('user_id'));

  function handleEmojiClick(emojiData) {
    setInput(input + emojiData.emoji);
    setShowEmoji(false);
  }

  function handleSend() {
    if (!input.trim()) return;
    if (isGroupChat) {
      sendMessage({
        type: 'message',
        content: input,
        group_id: group.group_id,
      });
    } else {
      sendMessage({
        type: 'message',
        content: input,
        receiver: { user_id: user.user_id },
      });
    }
    setInput('');
  }

  useEffect(() => {
    const removeHandler = addMessageHandler((data) => {
      if (data.type === 'message') {
        // group: does the message belong to this group chat
        if (isGroupChat && data.group_id === group.group_id) {
          setMessages((msgs) => [...msgs, data]);
        }
        // private chat: does the message belong to this user
        if (
          !isGroupChat &&
          (
            (data.sender.user_id === user.user_id && data.receiver.user_id === myUserId) ||
            (data.sender.user_id === myUserId && data.receiver.user_id === user.user_id)
          )
        ) {
          setMessages((msgs) => [...msgs, data]);
        }
      }
    });
    return () => removeHandler && removeHandler();
  }, [user, group, isGroupChat]);

  return (
    <div className="fixed bottom-4 right-4 z-50 border border-gray-300 rounded-lg shadow-lg">
      <div className="bg-white w-80 h-[40vh] p-3 flex flex-col rounded-lg shadow-lg overflow-hidden">
        <header className="flex justify-between items-center p-2 border-b border-gray-300">
          {isGroupChat ? (
            <GroupAvatar group={group} disableLink={true} size="sm" />
          ) : (
            <Author author={user} disableLink={true} size="sm" />
          )}
          <button onClick={onClose} className="text-xl leading-none">&times;</button>
        </header>

        <div className="flex-1 p-2 flex flex-col space-y-2 overflow-y-auto">
          {messages.map((msg, idx) => (
            <div
              key={msg.id || idx}
              className={`flex flex-col ${msg.sender?.user_id === myUserId ? 'items-end' : 'items-start'}`}
            >
              <div className="flex items-center space-x-2">
                <span className="text-sm font-semibold">
                  {isGroupChat
                    ? (msg.sender?.nickname || msg.sender?.first_name || 'User')
                    : (msg.sender?.user_id === myUserId ? 'Me' : user.nickname || user.first_name)
                  }
                </span>
                <span className="text-xs text-gray-500">{msg.timestamp || ''}</span>
              </div>
              <div
                className={`mt-1 inline-block px-3 py-2 rounded-lg max-w-[50%]
                  ${msg.sender?.user_id === myUserId ? 'rounded-br-none' : 'rounded-bl-none'}`}
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