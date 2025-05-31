# Inboxly

A modern real-time chat application featuring authentication, persistent message history, and a sleek, responsive UI. Built with Go (Gin) for the backend and React (Vite + TypeScript + Tailwind CSS) for the frontend.

## Features

- Real-time messaging with WebSockets
- JWT authentication and protected routes
- PostgreSQL database integration (Supabase-ready)
- Modern, minimal UI with responsive design
- Persistent message history
- RESTful API

## Tech Stack

- **Backend:** Go, Gin, Gorilla WebSocket, PostgreSQL, JWT, Supabase
- **Frontend:** React 18, Vite, TypeScript, Tailwind CSS, WebSocket

## Getting Started

### Prerequisites
- Go 1.21+
- Node.js 18+
- PostgreSQL database (Supabase recommended)

### Backend Setup
1. Clone the repository and navigate to `backend/`
2. Install Go dependencies
3. Create a `.env` file (see `backend/Readme.md` for example)
4. Run the backend:
   ```bash
   go run cmd/server/main.go
   ```

### Frontend Setup
1. Navigate to `frontend/`
2. Install dependencies:
   ```bash
   npm install
   ```
3. Copy `.env.example` to `.env` and set your backend URLs
4. Start the frontend:
   ```bash
   npm run dev
   ```
5. Open [http://localhost:3000](http://localhost:3000)

## Project Structure

```
backend/   # Go backend (API, WebSocket, DB)
frontend/  # React frontend (UI)
```

## Deployment

Deployed app: [Add your deployment link here]

## License

This project is created for educational and portfolio purposes. 