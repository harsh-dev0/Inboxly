import { useEffect, useRef, useState, useCallback } from 'react';
import type { ChatMessage, WSMessage } from '../types';
import { useAuth } from '../contexts/AuthContext';

const isSecure = window.location.protocol === 'https:';
const wsProtocol = isSecure ? 'wss://' : 'ws://';
const httpProtocol = isSecure ? 'https://' : 'http://';
const WS_BASE = import.meta.env.VITE_WS_URL || `${wsProtocol}${window.location.hostname}:8080`;

interface HistoryMessage {
  user_id: number;
  username: string;
  content: string;
  created_at: string;
}

export const useWebSocket = () => {
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [isConnected, setIsConnected] = useState(false);
  const [onlineUsers, setOnlineUsers] = useState(0);
  const [connectionError, setConnectionError] = useState<string | null>(null);
  const wsRef = useRef<WebSocket | null>(null);
  const { token, user } = useAuth();
  const reconnectTimeoutRef = useRef<number | null>(null);
  const connectAttemptsRef = useRef(0);
  const messageIdsRef = useRef<Set<string>>(new Set());
  const pendingSentMessagesRef = useRef<Set<string>>(new Set());
  
  const clearChat = useCallback(() => {
    setMessages([]);
    messageIdsRef.current.clear();
    pendingSentMessagesRef.current.clear();
  }, []);

  const fetchChatHistory = async () => {
    try {
      // Use HTTPS when the page is served over HTTPS, HTTP otherwise
      const httpBaseUrl = WS_BASE.replace(wsProtocol, httpProtocol);
      
      const response = await fetch(`${httpBaseUrl}/api/chat/messages`, {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });
      
      if (response.ok) {
        const data = await response.json();
        if (data && data.data && Array.isArray(data.data)) {
          const historyMessages = data.data.map((msg: HistoryMessage) => ({
            type: 'chat_message',
            user_id: msg.user_id,
            username: msg.username,
            content: msg.content,
            timestamp: msg.created_at
          }));
          
          setMessages([]);
          messageIdsRef.current.clear();
          pendingSentMessagesRef.current.clear();
          
          setMessages(historyMessages);
          historyMessages.forEach((msg: ChatMessage) => {
            const msgId = `${msg.user_id}-${msg.content}-${msg.timestamp}`;
            messageIdsRef.current.add(msgId);
          });
        }
      }
    } catch (error) {
      console.error('Failed to fetch chat history:', error);
    }
  };
  
  const addUniqueMessage = useCallback((message: ChatMessage): boolean => {
    const messageId = `${message.user_id}-${message.content}-${message.timestamp}`;
    
    if (pendingSentMessagesRef.current.has(messageId)) {
      pendingSentMessagesRef.current.delete(messageId);
      return false;
    }
    
    if (!messageIdsRef.current.has(messageId)) {
      messageIdsRef.current.add(messageId);
      setMessages(prev => [...prev, message]);
      return true;
    }
    return false;
  }, []);

  const connect = useCallback(() => {
    if (!token || !user) {
      setConnectionError('No authentication token available');
      return;
    }

    if (wsRef.current && wsRef.current.readyState !== WebSocket.CLOSED) {
      return;
    }

    try {
      connectAttemptsRef.current += 1;
      
      if (connectAttemptsRef.current > 5) {
        setTimeout(() => {
          connectAttemptsRef.current = 0;
        }, 10000);
        return;
      }
      
      const wsUrl = `${WS_BASE}/api/chat/ws?token=${encodeURIComponent(token)}`;
      const ws = new WebSocket(wsUrl);
      wsRef.current = ws;
      setConnectionError(null);

      ws.onopen = () => {
        setIsConnected(true);
        setConnectionError(null);
        connectAttemptsRef.current = 0;
        
        fetchChatHistory();
      };

      ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          
          switch (data.type) {
            case 'message':
            case 'chat_message':
              if (data.payload) {
                const chatMessage: ChatMessage = {
                  type: data.type,
                  user_id: data.payload.user_id || user?.id || 0,
                  username: data.payload.username || user?.username || 'Unknown',
                  content: data.payload.content,
                  timestamp: data.payload.timestamp || new Date().toISOString()
                };
                addUniqueMessage(chatMessage);
              } else {
                const chatMessage: ChatMessage = {
                  type: data.type,
                  user_id: data.user_id || user?.id || 0,
                  username: data.username || user?.username || 'Unknown',
                  content: data.content,
                  timestamp: data.timestamp || new Date().toISOString()
                };
                addUniqueMessage(chatMessage);
              }
              break;
            case 'user_joined':
              if (data.payload && data.payload.username) {
                const systemMessage: ChatMessage = {
                  type: 'system',
                  user_id: 0,
                  username: 'System',
                  content: `${data.payload.username} joined the chat`,
                  timestamp: new Date().toISOString()
                };
                addUniqueMessage(systemMessage);
              }
              break;
            case 'user_left':
              if (data.payload && data.payload.username) {
                const systemMessage: ChatMessage = {
                  type: 'system',
                  user_id: 0,
                  username: 'System',
                  content: `${data.payload.username} left the chat`,
                  timestamp: new Date().toISOString()
                };
                addUniqueMessage(systemMessage);
              }
              break;
            case 'online_count':
              if (typeof data.payload === 'number') {
                setOnlineUsers(data.payload);
              }
              break;
          }
        } catch (error) {
          console.error('Error parsing WebSocket message:', error);
        }
      };

      ws.onclose = () => {
        if (isConnected) {
          setIsConnected(false);
        }
        
        if (reconnectTimeoutRef.current !== null) {
          clearTimeout(reconnectTimeoutRef.current);
        }
        
        if (connectAttemptsRef.current < 5) {
          const reconnectDelay = Math.min(1000 * Math.pow(1.5, connectAttemptsRef.current), 10000);
          reconnectTimeoutRef.current = window.setTimeout(() => {
            if (token) {
              connect();
            }
          }, reconnectDelay);
        }
      };

      ws.onerror = (error) => {
        console.error('WebSocket error:', error);
        setConnectionError('Failed to connect to chat server');
        if (isConnected) {
          setIsConnected(false);
        }
      };
    } catch (error) {
      console.error('Connection error:', error);
      setConnectionError('Failed to create WebSocket connection');
    }
  }, [token, user, addUniqueMessage]);

  const disconnect = useCallback(() => {
    connectAttemptsRef.current = 0;
    
    if (reconnectTimeoutRef.current !== null) {
      clearTimeout(reconnectTimeoutRef.current);
      reconnectTimeoutRef.current = null;
    }
    
    if (wsRef.current) {
      const ws = wsRef.current;
      wsRef.current = null;
      
      try {
        ws.close();
      } catch (e) {
        console.error('Error closing WebSocket:', e);
      }
    }
    
    setIsConnected(false);
  }, []);

  const sendMessage = useCallback((content: string) => {
    if (!content.trim()) return;
    
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      const message: WSMessage = {
        type: 'chat_message',
        payload: { content }
      };
      wsRef.current.send(JSON.stringify(message));
    } else {
      connect();
    }
  }, [user, connect]);

  useEffect(() => {
    if (token && user) {
      connect();
    }

    return () => {
      disconnect();
    };
  }, [token, user, connect, disconnect]);

  return {
    messages,
    isConnected,
    onlineUsers,
    connectionError,
    sendMessage,
    connect,
    disconnect,
    clearChat
  };
};