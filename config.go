package main

import (
	"fmt"

	"github.com/dwikie/sentra-payment-orchestrator/config"
)

func InitializeApp() (*App, error) {
	config.LoadEnv()

	pool, err := config.InitDb()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	redis, err := config.InitRedis()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize redis: %w", err)
	}

	return &App{
		Pool:  pool,
		Redis: redis,
	}, nil
}
