package chat

import (
	"backend/internal/database"
	"encoding/json"
	"log"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) broadcastUserCount() {
	count := len(h.clients)
	
	message := WSMessage{
		Type:    "online_count",
		Payload: count,
	}
	messageBytes, _ := json.Marshal(message)

	for c := range h.clients {
		select {
		case c.send <- messageBytes:
		default:
			close(c.send)
			delete(h.clients, c)
		}
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true

			// Send current user count immediately to the new client
			countMsg := WSMessage{
				Type:    "online_count",
				Payload: len(h.clients),
			}
			countBytes, _ := json.Marshal(countMsg)
			client.send <- countBytes

			// Notify all clients that a user joined
			userJoined := UserJoined{
				Username: client.username,
				Message:  client.username + " joined the chat",
			}
			message := WSMessage{
				Type:    "user_joined",
				Payload: userJoined,
			}
			messageBytes, _ := json.Marshal(message)

			for c := range h.clients {
				select {
				case c.send <- messageBytes:
				default:
					close(c.send)
					delete(h.clients, c)
				}
			}

			// Update online count after a user joins
			h.broadcastUserCount()

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)

				// Notify all clients that a user left
				userLeft := UserLeft{
					Username: client.username,
					Message:  client.username + " left the chat",
				}
				message := WSMessage{
					Type:    "user_left",
					Payload: userLeft,
				}
				messageBytes, _ := json.Marshal(message)

				for c := range h.clients {
					select {
					case c.send <- messageBytes:
					default:
						close(c.send)
						delete(h.clients, c)
					}
				}

				// Update online count after a user leaves
				h.broadcastUserCount()
			}

		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (h *Hub) saveMessage(message Message) {
	query := `INSERT INTO messages (user_id, username, content, created_at) VALUES ($1, $2, $3, $4)`
	_, err := database.DB.Exec(query, message.UserID, message.Username, message.Content, message.Timestamp)
	if err != nil {
		log.Printf("Error saving message: %v", err)
	}
}
