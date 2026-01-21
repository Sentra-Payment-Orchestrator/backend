package helper

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	Argon2Version = argon2.Version
	Argon2Memory  = 64 * 1024
	Argon2Iter    = 1
	Argon2Thread  = 4
	Argon2SaltLen = 16
	Argon2KeyLen  = 32
)

func generateSalt(length uint32) ([]byte, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func HashPassword(password string) (string, error) {
	salt, err := generateSalt(Argon2SaltLen)
	if err != nil {
		return "", err
	}

	hashed := argon2.IDKey([]byte(password), salt, Argon2Iter, Argon2Memory, Argon2Thread, Argon2KeyLen)

	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	encodedHash := base64.RawStdEncoding.EncodeToString(hashed)
	log.Println(encodedHash)

	finalHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", Argon2Version, Argon2Memory, Argon2Iter, Argon2Thread, encodedSalt, encodedHash)

	return finalHash, nil
}

func VerifyPassword(encodedHash, password string) error {
	parts := strings.Split(encodedHash, "$")

	if len(parts) != 6 {
		return fmt.Errorf("invalid hash format")
	}

	var memory, iterations uint32
	var parallelism uint8
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)
	if err != nil {
		return fmt.Errorf("failed to parse parameters: %w", err)
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return fmt.Errorf("failed to decode salt: %w", err)
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return fmt.Errorf("failed to decode hash: %w", err)
	}

	computedHash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, uint32(len(hash)))

	fmt.Println(encodedHash, computedHash)
	if subtle.ConstantTimeCompare(computedHash, hash) == 1 {
		return nil
	}

	return fmt.Errorf("password does not match")
}
