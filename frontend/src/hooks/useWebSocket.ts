import { useEffect, useRef, useState } from 'react';
import type { ChatMessage, WSMessage } from '../types';
import { useAuth } from '../contexts/AuthContext';

const WS_URL = import.meta.env.VITE_WS_URL || 'ws://localhost:8080';

export const useWebSocket = () => {
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [isConnected, setIsConnected] = useState(false);
  const [onlineUsers, setOnlineUsers] = useState(0);
  const wsRef = useRef<WebSocket | null>(null);
  const { token, user } = useAuth();

  const connect = () => {
    if (!token) return;

    const ws = new WebSocket(`${WS_URL}/api/chat/ws?token=${token}`);
    wsRef.current = ws;

    ws.onopen = () => {
      console.log('WebSocket connected');
      setIsConnected(true);
    };

    ws.onmessage = (event) => {
      try {
        const data: WSMessage = JSON.parse(event.data);
        
        switch (data.type) {
          case 'message':
            setMessages(prev => [...prev, data.payload as ChatMessage]);
            break;
          case 'user_joined':
            console.log(`${data.payload.username} joined`);
            break;
          case 'user_left':
            console.log(`${data.payload.username} left`);
            break;
          case 'online_count':
            setOnlineUsers(data.payload);
            break;
          default:
            console.log('Unknown message type:', data.type);
        }
      } catch (error) {
        console.error('Error parsing WebSocket message:', error);
      }
    };

    ws.onclose = () => {
      console.log('WebSocket disconnected');
      setIsConnected(false);
      // Attempt to reconnect after 3 seconds
      setTimeout(() => {
        if (token) {
          connect();
        }
      }, 3000);
    };

    ws.onerror = (error) => {
      console.error('WebSocket error:', error);
      setIsConnected(false);
    };
  };

  const disconnect = () => {
    if (wsRef.current) {
      wsRef.current.close();
      wsRef.current = null;
    }
    setIsConnected(false);
  };

  const sendMessage = (content: string) => {
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      const message: WSMessage = {
        type: 'message',
        payload: { content }
      };
      wsRef.current.send(JSON.stringify(message));
    }
  };

  useEffect(() => {
    if (token && user) {
      connect();
    }

    return () => {
      disconnect();
    };
  }, [token, user]);

  return {
    messages,
    isConnected,
    onlineUsers,
    sendMessage,
    connect,
    disconnect
  };
};