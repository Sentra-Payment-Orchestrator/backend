package main

import (
	"github.com/dwikie/sentra-payment-orchestrator/handler"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Handlers struct {
	Auth *handler.AuthHandler
	User *handler.UserHandler
}

func (a *App) InitializeHandlers(pool *pgxpool.Pool, redis *redis.Client) {
	UserHandler := handler.NewUserHandler(pool, &handler.UserHandlerDependencies{})
	AuthHandler := handler.NewAuthHandler(pool, redis, &handler.AuthHandlerDependencies{User: UserHandler})

	a.Handlers = &Handlers{
		Auth: AuthHandler,
		User: UserHandler,
	}
}
