package chat

import (
	"backend/internal/database"
	"backend/internal/models"
	"backend/pkg/utils"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var hub = NewHub()

func init() {
	go hub.Run()
}

func WebSocketHandler(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", "missing_user")
		return
	}

	username, exists := c.Get("username")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Username not found", "missing_username")
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		userID:   userID.(int),
		username: username.(string),
	}

	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}

func GetMessagesHandler(c *gin.Context) {
	query := `
		SELECT id, user_id, username, content, created_at 
		FROM messages 
		ORDER BY created_at DESC 
		LIMIT 50
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch messages", err.Error())
		return
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		err := rows.Scan(&msg.ID, &msg.UserID, &msg.Username, &msg.Content, &msg.CreatedAt)
		if err != nil {
			log.Printf("Error scanning message: %v", err)
			continue
		}
		messages = append(messages, msg)
	}

	// Reverse to get chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	utils.SuccessResponse(c, "Messages retrieved successfully", messages)
}

func SendMessageHandler(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", "missing_user")
		return
	}

	username, exists := c.Get("username")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Username not found", "missing_username")
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	query := `INSERT INTO messages (user_id, username, content, created_at) VALUES ($1, $2, $3, CURRENT_TIMESTAMP)`
	result, err := database.DB.Exec(query, userID, username, req.Content)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to save message", err.Error())
		return
	}

	id, _ := result.LastInsertId()
	msg := models.Message{
		ID:        int(id),
		UserID:    userID.(int),
		Username:  username.(string),
		Content:   req.Content,
		CreatedAt: time.Now(),
	}

	utils.SuccessResponse(c, "Message sent successfully", msg)
}
