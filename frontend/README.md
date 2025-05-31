# Chat App Frontend

A modern, minimal React frontend for the real-time chat application built with Vite, TypeScript, and Tailwind CSS.

## Features

- 🔐 **Authentication** - Login and registration with JWT tokens
- 💬 **Real-time Chat** - WebSocket-based messaging
- 🎨 **Modern UI** - Clean, minimal design with Tailwind CSS
- 📱 **Responsive** - Mobile-first design
- ⚡ **Fast** - Built with Vite for optimal performance
- 🛡️ **Type Safe** - Full TypeScript support

## Tech Stack

- **React 18** - UI library
- **TypeScript** - Type safety
- **Vite** - Build tool and dev server
- **Tailwind CSS** - Styling
- **React Router** - Navigation
- **Axios** - HTTP client
- **WebSocket** - Real-time communication
- **Lucide React** - Icons

## Getting Started

### Prerequisites

- Node.js 18+ 
- npm or yarn
- Backend server running on port 8080

### Installation

1. **Clone and navigate to the frontend directory**
   ```bash
   cd frontend
   ```

2. **Install dependencies**
   ```bash
   npm install
   ```

3. **Create environment file**
   ```bash
   cp .env.example .env
   ```
   
   Update `.env` with your backend URL:
   ```
   VITE_API_URL=http://localhost:8080
   VITE_WS_URL=ws://localhost:8080
   ```

4. **Start the development server**
   ```bash
   npm run dev
   ```

5. **Open your browser**
   Navigate to `http://localhost:3000`

### Build for Production

```bash
npm run build
npm run preview
```

## Project Structure

```
src/
├── components/
│   ├── chat/          # Chat-specific components
│   ├── ui/            # Reusable UI components
│   └── ProtectedRoute.tsx
├── contexts/          # React contexts
├── hooks/             # Custom hooks
├── lib/               # Utilities and API
├── pages/             # Page components
├── types/             # TypeScript types
└── App.tsx
```

## Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build
- `npm run lint` - Run ESLint
- `npm run type-check` - Type checking without build

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `VITE_API_URL` | Backend API URL | `http://localhost:8080` |
| `VITE_WS_URL` | WebSocket URL | `ws://localhost:8080` |

## Features Overview

### Authentication
- JWT token-based authentication
- Automatic token refresh handling
- Protected routes
- Persistent login state

### Chat Interface
- Real-time messaging via WebSocket
- Message bubbles with timestamps
- Online user count
- Connection status indicator
- Auto-scroll to latest messages

### Responsive Design
- Mobile-first approach
- Adaptive layouts
- Touch-friendly interactions
- Modern animations and transitions
