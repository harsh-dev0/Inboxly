import React from 'react';
import { LogOut, MessageCircle, Users } from 'lucide-react';
import { useAuth } from '../../contexts/AuthContext';
import Button from '../ui/Button';

interface ChatHeaderProps {
  onlineUsers?: number;
}

const ChatHeader: React.FC<ChatHeaderProps> = ({ onlineUsers = 0 }) => {
  const { user, logout } = useAuth();

  return (
    <div className="bg-white border-b border-gray-200 px-4 py-3 flex items-center justify-between">
      <div className="flex items-center space-x-3">
        <div className="bg-primary-600 p-2 rounded-full">
          <MessageCircle className="w-6 h-6 text-white" />
        </div>
        <div>
          <h1 className="text-lg font-semibold text-gray-900">Chat Room</h1>
          <div className="flex items-center space-x-1 text-sm text-gray-500">
            <Users className="w-4 h-4" />
            <span>{onlineUsers} online</span>
          </div>
        </div>
      </div>

      <div className="flex items-center space-x-3">
        <div className="text-right">
          <div className="text-sm font-medium text-gray-900">{user?.username}</div>
          <div className="text-xs text-gray-500">Online</div>
        </div>
        <Button
          variant="ghost"
          size="sm"
          onClick={logout}
          className="text-gray-600 hover:text-red-600"
        >
          <LogOut className="w-4 h-4" />
        </Button>
      </div>
    </div>
  );
};

export default ChatHeader;