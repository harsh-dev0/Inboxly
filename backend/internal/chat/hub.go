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

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true

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

			log.Printf("Client registered: %s", client.username)

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

				log.Printf("Client unregistered: %s", client.username)
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
	query := `INSERT INTO messages (user_id, username, content) VALUES ($1, $2, $3)`
	_, err := database.DB.Exec(query, message.UserID, message.Username, message.Content)
	if err != nil {
		log.Printf("Error saving message: %v", err)
	}
}
