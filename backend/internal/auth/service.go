package auth

import (
	"backend/internal/database"
	"backend/internal/models"
	"database/sql"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateToken(userID int, username string) (string, error) {
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func CreateUser(req models.RegisterRequest) (*models.User, error) {
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id, created_at`
	var user models.User
	err = database.DB.QueryRow(query, req.Username, req.Email, hashedPassword).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	user.Username = req.Username
	user.Email = req.Email
	return &user, nil
}

func AuthenticateUser(req models.LoginRequest) (*models.User, error) {
	query := `SELECT id, username, email, password_hash, created_at FROM users WHERE username = $1`
	var user models.User
	var passwordHash string

	err := database.DB.QueryRow(query, req.Username).Scan(
		&user.ID, &user.Username, &user.Email, &passwordHash, &user.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	if !CheckPasswordHash(req.Password, passwordHash) {
		return nil, errors.New("invalid credentials")
	}

	return &user, nil
}
