package main

import (
	"github.com/dwikie/sentra-payment-orchestrator/handler"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handlers struct {
	Auth *handler.AuthHandler
	User *handler.UserHandler
}

func InitializeHandlers(pool *pgxpool.Pool) *Handlers {
	UserHandler := handler.NewUserHandler(pool, &handler.UserHandlerDependencies{})
	AuthHandler := handler.NewAuthHandler(pool, &handler.AuthHandlerDependencies{User: UserHandler})

	return &Handlers{
		Auth: AuthHandler,
		User: UserHandler,
	}
}
