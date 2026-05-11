package routes

import (
	"fmt"
	"net/http"
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

	now := time.Now()
	rts := viper.GetString("REFRESH_TOKEN_SECRET")

	rt, err := helper.CreateToken(
		[]byte(rts),
		paseto.JSONToken{
			Subject:    fmt.Sprintf("%d", user.Id),
			IssuedAt:   now,
			NotBefore:  now,
			Expiration: now.Add(24 * time.Hour),
		},
		"",
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create refresh token"})
		return
	}

	ats := viper.GetString("ACCESS_TOKEN_SECRET")

	at, err := helper.CreateToken(
		[]byte(ats),
		paseto.JSONToken{
			Subject:    fmt.Sprintf("%d", user.Id),
			IssuedAt:   now,
			NotBefore:  now,
			Expiration: now.Add(15 * time.Minute),
		},
		"",
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create access token"})
		return
	}

	domain := viper.GetString("DOMAIN")
	c.SetCookie("refresh_token", rt, 3600*24, "/", domain, false, true)

	err = h.Handlers.User.UpdateLastLogin(ctx, user.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update last login"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"data": gin.H{
			"access_token":  at,
			"refresh_token": rt,
		},
	})
}

func (h *AuthRoute) Logout(c *gin.Context) {
	domain := viper.GetString("DOMAIN")
	c.SetCookie("refresh_token", "", -1, "/", domain, false, true)
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
