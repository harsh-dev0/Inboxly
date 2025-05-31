import React from 'react';
import type { ChatMessage } from '../../types';
import { formatTime } from '../../lib/utils';
import { useAuth } from '../../contexts/AuthContext';

interface MessageBubbleProps {
  message: ChatMessage;
}

const MessageBubble: React.FC<MessageBubbleProps> = ({ message }) => {
  const { user } = useAuth();
  const isOwnMessage = user?.id === message.user_id;

  return (
    <div className={`flex mb-4 animate-slide-up ${isOwnMessage ? 'justify-end' : 'justify-start'}`}>
      <div className={`max-w-xs lg:max-w-md px-4 py-2 rounded-2xl ${
        isOwnMessage 
          ? 'bg-primary-600 text-white rounded-br-sm' 
          : 'bg-white text-gray-800 border border-gray-200 rounded-bl-sm'
      }`}>
        {!isOwnMessage && (
          <div className="text-xs font-medium text-primary-600 mb-1">
            {message.username}
          </div>
        )}
        <div className="break-words">
          {message.content}
        </div>
        <div className={`text-xs mt-1 ${
          isOwnMessage ? 'text-primary-100' : 'text-gray-500'
        }`}>
          {formatTime(message.timestamp)}
        </div>
      </div>
    </div>
  );
};

export default MessageBubble;