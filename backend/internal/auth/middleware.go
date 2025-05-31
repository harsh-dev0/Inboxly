package auth

import (
	"backend/pkg/utils"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Authorization header required", "missing_token")
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		if tokenString == authHeader {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid authorization format", "invalid_format")
			c.Abort()
			return
		}

		claims, err := ValidateToken(tokenString)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid token", err.Error())
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}

// WebSocketAuthMiddleware authenticates WebSocket connections that pass token as a query parameter
func WebSocketAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// For WebSocket, token is passed as a query parameter
		tokenParam := c.Query("token")
		log.Printf("WebSocket connection attempt with token: %s...", tokenParam[:10])
		
		if tokenParam == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Token query parameter required", "missing_token")
			c.Abort()
			return
		}

		claims, err := ValidateToken(tokenParam)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid token", err.Error())
			c.Abort()
			return
		}

		log.Printf("WebSocket auth successful for user: %s (ID: %d)", claims.Username, claims.UserID)
		
		// Set user info in context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
