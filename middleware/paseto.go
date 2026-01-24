package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dwikie/sentra-payment-orchestrator/helper"
	"github.com/gin-gonic/gin"
	"github.com/o1egl/paseto"
	"github.com/spf13/viper"
)

type ClaimsValidator func(t *paseto.JSONToken) error

func PasetoAuth(validators ...ClaimsValidator) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
			return
		}

		parts := strings.Fields(authHeader)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid Authorization header format"})
			return
		}

		secret := viper.GetString("ACCESS_TOKEN_SECRET")
		tokenValidator := &helper.Validator{
			NotBefore:  true,
			Expiration: true,
			CustomValidators: map[string]func(string) error{
				"nbf": func(value string) error {
					return nil
				},
			},
		}
		claims, _, err := helper.DecodeToken([]byte(secret), parts[1], tokenValidator)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		for _, validate := range validators {
			if validate == nil {
				continue
			}
			if err := validate(claims); err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				return
			}
		}

		c.Set("paseto_claims", claims)
		c.Next()
	}
}

func RequireAudience(expected string) ClaimsValidator {
	return func(t *paseto.JSONToken) error {
		if t.Audience != expected {
			return fmt.Errorf("invalid token: incorrect audience")
		}

		return nil
	}
}

func RequireSubject(expected string) ClaimsValidator {
	return func(t *paseto.JSONToken) error {
		if t.Subject != expected {
			return fmt.Errorf("invalid token: incorrect subject")
		}

		return nil
	}
}

func RequireIssuer(expected string) ClaimsValidator {
	return func(t *paseto.JSONToken) error {
		if t.Issuer != expected {
			return fmt.Errorf("invalid token: incorrect issuer")
		}

		return nil
	}
}
