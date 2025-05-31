# Real-time Chat Application Backend

A modern real-time chat application backend built with Go, featuring WebSocket connections, JWT authentication, and PostgreSQL database integration.

## 🚀 Features

- **Real-time Messaging**: WebSocket-based chat with instant message delivery
- **User Authentication**: JWT-based authentication with secure password hashing
- **Database Integration**: PostgreSQL with Supabase support
- **RESTful API**: Clean API design following industry standards
- **CORS Support**: Configured for frontend integration
- **Message History**: Persistent message storage and retrieval

## 🛠️ Tech Stack

- **Framework**: Gin (Go web framework)
- **WebSockets**: Gorilla WebSocket
- **Database**: PostgreSQL (Supabase)
- **Authentication**: JWT tokens with bcrypt password hashing
- **Environment**: Go 1.21+

## 📋 Prerequisites

- Go 1.21 or higher
- PostgreSQL database (Supabase account recommended)
- Git

## 🔧 Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/harsh-dev0/Inboxly
   cd backend
   ```

2. **Initialize Go module and install dependencies**
   ```bash
   go mod init backend
   go get github.com/gin-gonic/gin
   go get github.com/gorilla/websocket
   go get github.com/lib/pq
   go get github.com/golang-jwt/jwt/v5
   go get github.com/joho/godotenv
   go get golang.org/x/crypto/bcrypt
   go get github.com/gin-contrib/cors
   ```

3. **Set up environment variables**
   Create a `.env` file in the root directory:
   ```env
   DB_HOST=db.your-supabase-project.supabase.co
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=your-supabase-password
   DB_NAME=postgres
   DB_SSLMODE=require
   JWT_SECRET=your-super-secret-jwt-key-here
   PORT=8080
   GIN_MODE=debug
   FRONTEND_URL=http://localhost:3000
   ```

4. **Run the application**
   ```bash
   go run cmd/server/main.go
   ```

## 🌐 API Endpoints

### Authentication
- `POST /api/auth/register` - Register a new user
- `POST /api/auth/login` - Login user
- `GET /api/auth/profile` - Get user profile (protected)

### Chat
- `GET /api/chat/messages` - Get message history (protected)
- `GET /api/chat/ws` - WebSocket connection for real-time chat (protected)

### Health Check
- `GET /health` - API health status

## 🔌 WebSocket Events

### Client to Server
```json
{
  "type": "chat_message",
  "payload": {
    "content": "Hello, everyone!"
  }
}
```

### Server to Client
```json
{
  "type": "chat_message",
  "user_id": 1,
  "username": "john_doe",
  "content": "Hello, everyone!",
  "timestamp": "2025-05-30T10:30:00Z"
}
```

## 📝 Usage Examples

### Register a new user
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "securepassword"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "securepassword"
  }'
```

### Get messages (with auth token)
```bash
curl -X GET http://localhost:8080/api/chat/messages \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 🏗️ Project Structure

```
backend/
├── cmd/server/main.go          # Application entry point
├── internal/
│   ├── auth/                   # Authentication logic
│   ├── chat/                   # Chat functionality
│   ├── database/               # Database connection
│   └── models/                 # Data models
├── pkg/utils/                  # Utility functions
├── .env                        # Environment variables
└── README.md                   # Project documentation
```

## 📄 License

This project is created for educational purposes.