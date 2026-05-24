package handler

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/dwikie/sentra-payment-orchestrator/helper"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/o1egl/paseto"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type TokenType string

const (
	AccessToken  TokenType = "access_token"
	RefreshToken TokenType = "refresh_token"
)

type AuthHandler struct {
	Pool     *pgxpool.Pool
	Redis    *redis.Client
	Handlers *AuthHandlerDependencies
}

type AuthHandlerDependencies struct {
	User *UserHandler
}

func NewAuthHandler(pool *pgxpool.Pool, redis *redis.Client, deps *AuthHandlerDependencies) *AuthHandler {
	return &AuthHandler{
		Pool:     pool,
		Redis:    redis,
		Handlers: deps,
	}
}

func (h *AuthHandler) CreateRefreshToken(ctx context.Context, userID int64) (token, jti string, err error) {
	now := time.Now()
	rts := viper.GetString("REFRESH_TOKEN_SECRET")
	uuid, err := uuid.NewRandom()
	if err != nil {
		return "", "", err
	}

	jti = uuid.String()

	rt, err := helper.CreateToken(
		[]byte(rts),
		paseto.JSONToken{
			Subject:    strconv.FormatInt(userID, 10),
			IssuedAt:   now,
			NotBefore:  now,
			Expiration: now.Add(24 * time.Hour),
			Jti:        jti,
		},
		"",
	)

	if err != nil {
		return "", "", err
	}

	return rt, jti, nil
}

func (h *AuthHandler) CreateAccessToken(ctx context.Context, userID int64) (token, jti string, err error) {
	now := time.Now()
	ats := viper.GetString("ACCESS_TOKEN_SECRET")
	uuid, err := uuid.NewV7()
	if err != nil {
		return "", "", err
	}

	jti = uuid.String()

	at, err := helper.CreateToken(
		[]byte(ats),
		paseto.JSONToken{
			Subject:    strconv.FormatInt(userID, 10),
			IssuedAt:   now,
			NotBefore:  now,
			Expiration: now.Add(15 * time.Minute),
			Jti:        jti,
		},
		"",
	)
	if err != nil {
		return "", "", err
	}

	return at, jti, nil
}

func (h *AuthHandler) StoreSession(ctx context.Context, userId int64, rt, jti string) error {
	tx := h.Redis.TxPipeline()

	sessionKey := fmt.Sprintf("session:%d", userId)
	sesinfoKey := fmt.Sprintf("session_info:%s", jti)

	hasher := sha256.New()
	hasher.Write([]byte(rt))
	hashedToken := hex.EncodeToString(hasher.Sum(nil))

	tx.HSet(ctx, sesinfoKey, map[string]string{
		"refresh_token": hashedToken,
		"created_at":    time.Now().Format(time.RFC3339),
	})
	tx.Expire(ctx, sesinfoKey, 24*time.Hour)

	tx.SAdd(ctx, sessionKey, jti)
	tx.Expire(ctx, sessionKey, 24*time.Hour)

	if _, err := tx.Exec(ctx); err != nil {
		return fmt.Errorf("error while storing session in Redis: %v", err)
	}

	return nil
}

func (h *AuthHandler) InvalidateSession(ctx context.Context, userId int64, jti string) error {
	tx := h.Redis.TxPipeline()
	fmt.Printf("Invalidating session for user %d with jti %s\n", userId, jti)

	sessionKey := fmt.Sprintf("session:%d", userId)
	sesinfoKey := fmt.Sprintf("session_info:%s", jti)

	tx.SRem(ctx, sessionKey, jti)
	tx.Del(ctx, sesinfoKey)

	if _, err := tx.Exec(ctx); err != nil {
		return fmt.Errorf("error while invalidating session in Redis: %v", err)
	}

	return nil
}

func (h *AuthHandler) InvalidateAllSessions(ctx context.Context, userId int64) error {
	sessionKey := fmt.Sprintf("session:%d", userId)

	jtis, err := h.Redis.SMembers(ctx, sessionKey).Result()
	if err != nil {
		return fmt.Errorf("error while retrieving sessions from Redis: %v", err)
	}

	tx := h.Redis.TxPipeline()

	for _, jti := range jtis {
		sesinfoKey := fmt.Sprintf("session_info:%s", jti)
		tx.Del(ctx, sesinfoKey)
	}

	tx.Del(ctx, sessionKey)

	if _, err := tx.Exec(ctx); err != nil {
		return fmt.Errorf("error while invalidating sessions in Redis: %v", err)
	}

	return nil
}

func (h *AuthHandler) RotateRefreshToken(ctx context.Context, userId int64, newJti, oldJti, newRt, oldRt string) error {
	tx := h.Redis.TxPipeline()

	sessionKey := fmt.Sprintf("session:%d", userId)
	oldSesinfoKey := fmt.Sprintf("session_info:%s", oldJti)
	newSesinfoKey := fmt.Sprintf("session_info:%s", newJti)

	hasher := sha256.New()
	hasher.Write([]byte(oldRt))
	hashedOldToken := hex.EncodeToString(hasher.Sum(nil))

	// Verify old refresh token before rotating
	storedHashedToken, err := h.Redis.HGet(ctx, oldSesinfoKey, "refresh_token").Result()
	if err != nil {
		return fmt.Errorf("error while retrieving old session info from Redis: %v", err)
	}
	if storedHashedToken != hashedOldToken {
		return fmt.Errorf("invalid refresh token: token mismatch")
	}

	hasher.Reset()
	hasher.Write([]byte(newRt))
	hashedNewToken := hex.EncodeToString(hasher.Sum(nil))

	tx.HSet(ctx, newSesinfoKey, map[string]string{
		"refresh_token": hashedNewToken,
		"created_at":    time.Now().Format(time.RFC3339),
	})
	tx.Expire(ctx, newSesinfoKey, 24*time.Hour)

	tx.SRem(ctx, sessionKey, oldJti)
	tx.SAdd(ctx, sessionKey, newJti)
	tx.Del(ctx, oldSesinfoKey)

	if _, err := tx.Exec(ctx); err != nil {
		return fmt.Errorf("error while rotating refresh token in Redis: %v", err)
	}

	return nil
}
