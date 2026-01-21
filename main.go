package main

import (
	"log"

	"github.com/dwikie/sentra-payment-orchestrator/config"
	"github.com/gin-gonic/gin"
)

type App struct {
	Handlers *Handlers
}

func main() {
	// Create a Gin router with default middleware (logger and recovery)
	r := gin.Default()

	err := config.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer config.Pool.Close()

	app := &App{
		Handlers: NewHandlers(config.Pool),
	}

	app.RegisterRoutes(r)

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
