package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dwikie/sentra-payment-orchestrator/helper"
	"github.com/dwikie/sentra-payment-orchestrator/model"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/o1egl/paseto"
	"github.com/spf13/viper"
)

type AuthHandler struct {
	Pool        *pgxpool.Pool
	UserHandler *UserHandler
}

func NewAuthHandler(pool *pgxpool.Pool, userHandler *UserHandler) *AuthHandler {
	return &AuthHandler{Pool: pool, UserHandler: userHandler}
}

func (h *AuthHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	payload := model.LoginRequest{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.UserHandler.getUserByEmail(ctx, payload.Email)
	if err != nil {
		fmt.Println(err)
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		}
		return
	}

	if err := helper.VerifyPassword(user.Password, payload.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	refreshTokenSecret := viper.GetString("REFRESH_TOKEN_SECRET")
	println(refreshTokenSecret)
	refreshToken, err := h.CreateToken([]byte(refreshTokenSecret), "refresh", paseto.JSONToken{}, "", map[string]string{
		"user_id": fmt.Sprintf("%d", user.Id),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create refresh token"})
		return
	}

	accessTokenSecret := viper.GetString("ACCESS_TOKEN_SECRET")
	accessToken, err := h.CreateToken([]byte(accessTokenSecret), "access", paseto.JSONToken{}, "", map[string]string{
		"user_id": fmt.Sprintf("%d", user.Id),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create access token"})
		return
	}

	domain := viper.GetString("DOMAIN")
	c.SetCookie("refresh_token", refreshToken, 3600*24, "/", domain, false, true)
	c.Header("Authorization", "Bearer "+accessToken)

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "data": gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	domain := viper.GetString("DOMAIN")
	c.SetCookie("refresh_token", "", -1, "/", domain, false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

func (h *AuthHandler) CreateToken(signature []byte, purpose string, jsonToken paseto.JSONToken, footer string, customClaims ...map[string]string) (string, error) {
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

func (h *AuthHandler) ParseToken(signature []byte, token string) (*paseto.JSONToken, error) {
	jsonToken := paseto.JSONToken{}

	err := paseto.NewV2().Decrypt(token, signature, &jsonToken, nil)
	if err != nil {
		return &jsonToken, fmt.Errorf("invalid token: %v", err)
	}

	now := time.Now()
	if jsonToken.Expiration.Before(now) {
		return &jsonToken, fmt.Errorf("invalid token: token has expired")
	}
	if jsonToken.NotBefore.After(now) {
		return &jsonToken, fmt.Errorf("invalid token: token not valid yet")
	}

	return &jsonToken, nil
}
