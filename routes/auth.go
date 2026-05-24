package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dwikie/sentra-payment-orchestrator/handler"
	"github.com/dwikie/sentra-payment-orchestrator/helper"
	"github.com/dwikie/sentra-payment-orchestrator/model"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/o1egl/paseto"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type AuthRoute struct {
	Handlers *AuthRouteHandlers
}

type AuthRouteHandlers struct {
	Auth *handler.AuthHandler
	User *handler.UserHandler
}

func InitAuthRoute(pool *pgxpool.Pool, redis *redis.Client) *AuthRoute {
	return &AuthRoute{Handlers: &AuthRouteHandlers{
		Auth: handler.NewAuthHandler(pool, redis, nil),
		User: handler.NewUserHandler(pool, nil),
	}}
}

func (h *AuthRoute) Login(c *gin.Context) {
	ctx := c.Request.Context()
	payload := model.LoginRequest{}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Handlers.User.GetUserByEmail(ctx, payload.Email)
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

	rt, jti, err := h.Handlers.Auth.CreateRefreshToken(ctx, user.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	at, _, err := h.Handlers.Auth.CreateAccessToken(ctx, user.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.Handlers.Auth.StoreSession(ctx, user.Id, rt, jti)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store session"})
		return
	}

	domain := viper.GetString("DOMAIN")
	c.SetCookie("refresh_token", rt, 3600*24, "/", domain, false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"data": gin.H{
			"access_token":  at,
			"refresh_token": rt,
		},
	})
}

func (h *AuthRoute) Logout(c *gin.Context) {
	ctx := c.Request.Context()
	rt, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing refresh token"})
		return
	}

	domain := viper.GetString("DOMAIN")
	c.SetCookie("refresh_token", "", -1, "/", domain, false, true)

	rts := viper.GetString("REFRESH_TOKEN_SECRET")
	claims, _, _ := helper.DecodeToken([]byte(rts), rt, nil)

	userId, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad token"})
		return
	}

	err = h.Handlers.Auth.InvalidateSession(ctx, userId, claims.Jti)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to invalidate session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

func (h *AuthRoute) RefreshToken(c *gin.Context) {
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
