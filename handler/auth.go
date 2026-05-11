package handler

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
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
