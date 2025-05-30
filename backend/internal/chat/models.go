package chat

import "time"

type Message struct {
	Type      string    `json:"type"`
	UserID    int       `json:"user_id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

type WSMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type ChatMessage struct {
	Content string `json:"content"`
}

type UserJoined struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

type UserLeft struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}
