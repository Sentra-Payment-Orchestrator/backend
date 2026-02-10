package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dwikie/sentra-payment-orchestrator/helper"
	"github.com/gin-gonic/gin"
	"github.com/o1egl/paseto"
	"github.com/spf13/viper"
)

type ClaimsValidator func(t *paseto.JSONToken) error

func RequiredAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
			return
		}

		parts := strings.Fields(authHeader)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || !strings.HasPrefix(parts[1], "v2.local.") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid Authorization header format"})
			return
		}

		secret := viper.GetString("ACCESS_TOKEN_SECRET")
		tokenValidator := map[string]func(string) error{
			"nbf": func(value string) error {
				nbf, err := time.Parse(time.RFC3339, value)
				if err != nil {
					return fmt.Errorf("invalid nbf claim: %v", err)
				}
				if nbf.After(time.Now()) {
					return fmt.Errorf("invalid token: not valid yet")
				}
				return nil
			},
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

		claims, _, err := helper.DecodeToken([]byte(secret), parts[1], tokenValidator)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.Set("user_id", claims.Subject)
		c.Set("claims", claims)
		c.Next()
	}
}
