export interface User {
    id: number;
    username: string;
    email: string;
    created_at: string;
  }
  
  export interface Message {
    id: number;
    user_id: number;
    username: string;
    content: string;
    created_at: string;
  }
  
  export interface LoginRequest {
    username: string;
    password: string;
  }
  
  export interface RegisterRequest {
    username: string;
    email: string;
    password: string;
  }
  
  export interface LoginResponse {
    token: string;
    user: User;
  }
  
  export interface WSMessage {
    type: string;
    payload: any;
  }
  
  export interface ChatMessage {
    type: string;
    user_id: number;
    username: string;
    content: string;
    timestamp: string;
  }
  
  export interface AuthContextType {
    user: User | null;
    token: string | null;
    login: (username: string, password: string) => Promise<void>;
    register: (username: string, email: string, password: string) => Promise<void>;
    logout: () => void;
    loading: boolean;
  }