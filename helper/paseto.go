package helper

import (
	"fmt"
	"time"

	"github.com/o1egl/paseto"
)

type Validator struct {
	NotBefore        bool
	Expiration       bool
	CustomValidators map[string]func(string) error
}

func CreateToken(signature []byte, exp time.Time, jsonToken paseto.JSONToken, footer string, customClaims ...map[string]string) (string, error) {
	now := time.Now()
	jsonToken.IssuedAt = now
	jsonToken.NotBefore = now
	for _, claims := range customClaims {
		for k, v := range claims {
			jsonToken.Set(k, v)
		}
	}

	jsonToken.Expiration = exp

	token, err := paseto.NewV2().Encrypt(signature, jsonToken, footer)
	if err != nil {
		return "", err
	}

	return token, nil
}

func DecodeToken(signature []byte, token string, tokenValidator *Validator) (*paseto.JSONToken, *string, error) {
	jsonToken := paseto.JSONToken{}
	footer := ""

	err := paseto.NewV2().Decrypt(token, signature, &jsonToken, &footer)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid token: %v", err)
	}

	err = validateToken(&jsonToken, tokenValidator)
	if err != nil {
		return nil, nil, fmt.Errorf("token validation failed: %v", err)
	}

	return &jsonToken, &footer, nil
}

func validateToken(token *paseto.JSONToken, tokenValidator *Validator) error {
	if tokenValidator == nil {
		return nil
	}

	now := time.Now()
	if tokenValidator.Expiration {
		if token.Expiration.Before(now) {
			return fmt.Errorf("invalid token: token has expired")
		}
	}

	if tokenValidator.NotBefore {
		if token.NotBefore.After(now) {
			return fmt.Errorf("invalid token: token not valid yet")
		}
	}

	for claim, validateFunc := range tokenValidator.CustomValidators {
		value := token.Get(claim)
		if value == "" {
			return fmt.Errorf("invalid token: missing claim %s", claim)
		}
		if err := validateFunc(value); err != nil {
			return fmt.Errorf("invalid token: claim %s validation failed: %v", claim, err)
		}
	}

	return nil
}
