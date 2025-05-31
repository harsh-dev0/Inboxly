import React from 'react';
import { MessageCircle } from 'lucide-react';

const Loading: React.FC = () => {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-primary-50 to-primary-100">
      <div className="text-center">
        <div className="bg-primary-600 p-4 rounded-full mx-auto mb-4 animate-pulse">
          <MessageCircle className="w-8 h-8 text-white" />
        </div>
        <div className="text-lg font-medium text-gray-700">Loading...</div>
        <div className="text-sm text-gray-500 mt-1">Please wait</div>
      </div>
    </div>
  );
};

export default Loading;