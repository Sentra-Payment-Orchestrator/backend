package handler

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthHandler struct {
	Pool     *pgxpool.Pool
	Handlers *AuthHandlerDependencies
}

type AuthHandlerDependencies struct {
	User *UserHandler
}

func NewAuthHandler(pool *pgxpool.Pool, deps *AuthHandlerDependencies) *AuthHandler {
	return &AuthHandler{
		Pool:     pool,
		Handlers: deps,
	}
}
