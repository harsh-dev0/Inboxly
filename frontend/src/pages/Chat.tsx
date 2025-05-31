import React, { useEffect, useRef } from 'react';
import { useWebSocket } from '../hooks/useWebSocket';
import ChatHeader from '../components/chats/ChatHeader';
import MessageBubble from '../components/chats/MessageBubble';
import MessageInput from '../components/chats/MessageInput';
import { Wifi, WifiOff, Trash2 } from 'lucide-react';

const Chat: React.FC = () => {
  const { messages, isConnected, onlineUsers, sendMessage, clearChat } = useWebSocket();
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const handleSendMessage = (content: string) => {
    sendMessage(content);
  };

  const handleClearChat = () => {
    if (window.confirm('Are you sure you want to clear the chat? This only affects your view.')) {
      clearChat();
    }
  };

  return (
    <div className="flex flex-col h-screen bg-gray-50">
      <ChatHeader onlineUsers={onlineUsers} />
      
      {/* Connection Status and Actions */}
      <div className="px-4 py-2 flex items-center justify-between">
        <div className={`text-sm flex items-center space-x-2 ${
          isConnected ? 'text-green-700' : 'text-red-700'
        }`}>
          {isConnected ? (
            <>
              <Wifi className="w-4 h-4" />
              <span>Connected ({onlineUsers} online)</span>
            </>
          ) : (
            <>
              <WifiOff className="w-4 h-4" />
              <span>Disconnected - Attempting to reconnect...</span>
            </>
          )}
        </div>
        
        <button 
          onClick={handleClearChat}
          className="text-gray-500 hover:text-red-500 p-1 rounded-full hover:bg-gray-100 transition"
          title="Clear chat"
        >
          <Trash2 className="w-4 h-4" />
        </button>
      </div>

      {/* Messages Container */}
      <div className="flex-1 overflow-y-auto px-4 py-4 space-y-1">
        {messages.length === 0 ? (
          <div className="flex items-center justify-center h-full">
            <div className="text-center text-gray-500">
              <div className="text-lg font-medium mb-2">Welcome to the chat!</div>
              <div className="text-sm">Start a conversation by sending a message below.</div>
            </div>
          </div>
        ) : (
          messages.map((message, index) => (
            <MessageBubble
              key={`${message.user_id}-${message.timestamp}-${index}`}
              message={message}
            />
          ))
        )}
        <div ref={messagesEndRef} />
      </div>

      <MessageInput 
        onSendMessage={handleSendMessage}
        disabled={!isConnected}
      />
    </div>
  );
};

export default Chat;