package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dwikie/sentra-payment-orchestrator/model"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/o1egl/paseto"
)

type AuthHandlers struct {
	Pool *pgxpool.Pool
}

func (h *AuthHandlers) Login(c *gin.Context) {
	ctx := c.Request.Context()
	payload := model.LoginRequest{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conn, err := h.Pool.Acquire(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}
	defer conn.Release()

	c.JSON(200, gin.H{"message": "Login successful"})
}

func (h *AuthHandlers) Logout(c *gin.Context) {
	// Implement logout logic here, e.g., invalidate tokens, clear sessions, etc.
	c.JSON(200, gin.H{"message": "Logout successful"})
}

func (h *AuthHandlers) Register(c *gin.Context) {
	ctx := c.Request.Context()
	payload := model.RegisterRequest{}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conn, err := h.Pool.Acquire(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	user := tx.QueryRow(ctx, `
	INSERT INTO users (username, password, status)
	VALUES ($1, $2, $3) RETURNING id
	`, payload.Username, payload.Password, 0)

	var userID uint8
	if err := user.Scan(&userID); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	_, err = tx.Exec(ctx, `
	INSERT INTO user_profile (user_id, first_name, last_name, email, phone_number)
	VALUES ($1, $2, $3, $4, $5)`,
		userID, payload.FirstName, payload.LastName, payload.Email, payload.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user profile"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Registration successful", "user_id": userID})
}

func (h *AuthHandlers) CreateToken(signature []byte, purpose string, jsonToken paseto.JSONToken, footer string, customClaims ...map[string]string) (string, error) {
	now := time.Now()
	jsonToken.IssuedAt = now
	jsonToken.NotBefore = now
	for _, claims := range customClaims {
		for k, v := range claims {
			jsonToken.Set(k, v)
		}
	}

	switch purpose {
	case "access":
		jsonToken.Expiration = now.Add(15 * time.Minute)
	case "refresh":
		jsonToken.Expiration = now.Add(24 * time.Hour)
	default:
		return "", fmt.Errorf("unknown token purpose: %s", purpose)
	}

	token, err := paseto.NewV2().Encrypt(signature, jsonToken, footer)
	if err != nil {
		return "", err
	}

	return token, nil
}
