package main

import (
	"fmt"

	"github.com/dwikie/sentra-payment-orchestrator/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Pool *pgxpool.Pool
}

func InitConfig() (*Config, error) {
	config.LoadEnv()

	pool, err := config.InitDb()
	if err != nil {
		return nil, fmt.Errorf("error while initialize database connection pool: %v", err)
	}

	return &Config{
		Pool: pool,
	}, nil
}
