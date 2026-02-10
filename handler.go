package main

import (
	"github.com/dwikie/sentra-payment-orchestrator/handler"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handlers struct {
	Auth         *handler.AuthHandler
	User         *handler.UserHandler
	Organization *handler.OrganizationHandler
}

func NewHandlers(pool *pgxpool.Pool) *Handlers {
	orgHandler := handler.NewOrganizationHandler(pool)
	userHandler := handler.NewUserHandler(pool, orgHandler)
	authHandler := handler.NewAuthHandler(pool, userHandler)

	return &Handlers{
		Auth:         authHandler,
		User:         userHandler,
		Organization: orgHandler,
	}
}
