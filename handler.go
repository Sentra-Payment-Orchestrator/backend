package main

import (
	"github.com/dwikie/sentra-payment-orchestrator/handler"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handlers struct {
	Auth *handler.AuthHandlers
}

// NewHandlers initializes all handlers
func NewHandlers(pool *pgxpool.Pool) *Handlers {
	return &Handlers{
		Auth: &handler.AuthHandlers{Pool: pool},
	}
}
