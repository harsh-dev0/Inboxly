package auth

import (
	"backend/internal/models"
	"backend/pkg/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func RegisterHandler(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	// Check if user already exists
	user, err := CreateUser(req)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			utils.ErrorResponse(c, http.StatusConflict, "User already exists", "duplicate_user")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user", err.Error())
		return
	}

	// Generate token
	token, err := GenerateToken(user.ID, user.Username)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token", err.Error())
		return
	}

	response := models.LoginResponse{
		Token: token,
		User:  *user,
	}

	utils.CreatedResponse(c, "User registered successfully", response)
}

func LoginHandler(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	user, err := AuthenticateUser(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Authentication failed", err.Error())
		return
	}

	token, err := GenerateToken(user.ID, user.Username)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token", err.Error())
		return
	}

	response := models.LoginResponse{
		Token: token,
		User:  *user,
	}

	utils.SuccessResponse(c, "Login successful", response)
}

func ProfileHandler(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", "missing_user")
		return
	}

	username, _ := c.Get("username")

	userData := gin.H{
		"user_id":  userID,
		"username": username,
	}

	utils.SuccessResponse(c, "Profile retrieved successfully", userData)
}
