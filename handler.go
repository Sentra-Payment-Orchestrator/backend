package main

import (
	"github.com/dwikie/sentra-payment-orchestrator/handler"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handlers struct {
	Auth *handler.AuthHandler
	User *handler.UserHandler
}

func NewHandlers(pool *pgxpool.Pool) *Handlers {
	userHandler := handler.NewUserHandler(pool)
	authHandler := handler.NewAuthHandler(pool, userHandler)

	return &Handlers{
		Auth: authHandler,
		User: userHandler,
	}
}
