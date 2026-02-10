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
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve user"})
		}
		return
	}

	if err := helper.VerifyPassword(user.Password, payload.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	refreshTokenSecret := viper.GetString("REFRESH_TOKEN_SECRET")
	now := time.Now()

	refreshClaims := paseto.JSONToken{
		Subject:    fmt.Sprintf("%d", user.Id),
		IssuedAt:   now,
		NotBefore:  now,
		Expiration: now.Add(24 * time.Hour),
	}

	refreshToken, err := helper.CreateToken([]byte(refreshTokenSecret), refreshClaims, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create refresh token"})
		return
	}

	accessTokenSecret := viper.GetString("ACCESS_TOKEN_SECRET")
	accessClaims := paseto.JSONToken{
		Subject:    fmt.Sprintf("%d", user.Id),
		IssuedAt:   now,
		NotBefore:  now,
		Expiration: now.Add(15 * time.Minute),
	}

	accessToken, err := helper.CreateToken([]byte(accessTokenSecret), accessClaims, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create access token"})
		return
	}

	domain := viper.GetString("DOMAIN")
	c.SetCookie("refresh_token", refreshToken, 3600*24, "/", domain, false, true)

	err = h.UserHandler.updateLastLogin(ctx, user.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update last login"})
		return
	}

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

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing refresh token"})
		return
	}

	refreshTokenSecret := viper.GetString("REFRESH_TOKEN_SECRET")
	tokenValidator := map[string]func(string) error{
		"exp": func(value string) error {
			exp, err := time.Parse(time.RFC3339, value)
			if err != nil {
				return fmt.Errorf("invalid exp claim: %v", err)
			}
			if exp.Before(time.Now()) {
				return fmt.Errorf("invalid token: expired")
			}
			return nil
		},
	}

	claims, _, err := helper.DecodeToken([]byte(refreshTokenSecret), refreshToken, tokenValidator)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	userId := claims.Get("user_id")
	if userId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token: missing user_id claim"})
		return
	}

	accessTokenSecret := viper.GetString("ACCESS_TOKEN_SECRET")
	now := time.Now()
	accessClaims := paseto.JSONToken{
		IssuedAt:   now,
		NotBefore:  now,
		Expiration: now.Add(15 * time.Minute),
	}

	accessToken, err := helper.CreateToken([]byte(accessTokenSecret), accessClaims, "", map[string]string{
		"user_id": userId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create access token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Token refreshed successfully", "data": gin.H{
		"access_token": accessToken,
	}})
}
