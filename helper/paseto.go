package helper

import (
	"fmt"

	"github.com/o1egl/paseto"
)

func CreateToken(signature []byte, jsonToken paseto.JSONToken, footer string, customClaims ...map[string]string) (string, error) {
	for _, claims := range customClaims {
		for k, v := range claims {
			jsonToken.Set(k, v)
		}
	}

	token, err := paseto.NewV2().Encrypt(signature, jsonToken, footer)
	if err != nil {
		return "", err
	}

	return token, nil
}

func DecodeToken(signature []byte, token string, tokenValidator map[string]func(string) error) (*paseto.JSONToken, *string, error) {
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

func validateToken(token *paseto.JSONToken, tokenValidator map[string]func(string) error) error {
	if tokenValidator == nil {
		return nil
	}

	for claim, validateFunc := range tokenValidator {
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
