package main

import (
	"backend/internal/auth"
	"backend/internal/chat"
	"backend/internal/database"

	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Connect to database
	database.Connect()
	defer database.DB.Close()

	// Setup Gin router
	r := gin.Default()

	// CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{os.Getenv("FRONTEND_URL")}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Chat API is running",
		})
	})

	// API routes
	api := r.Group("/api")
	{
		// Auth routes
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", auth.RegisterHandler)
			authGroup.POST("/login", auth.LoginHandler)
			authGroup.GET("/profile", auth.AuthMiddleware(), auth.ProfileHandler)
		}
		chatGroup := api.Group("/chat")
		{
			chatGroup.GET("/messages", chat.GetMessagesHandler)
			chatGroup.GET("/ws", auth.AuthMiddleware(), chat.WebSocketHandler)
			chatGroup.POST("/messages", auth.AuthMiddleware(), chat.SendMessageHandler)
		}

	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
