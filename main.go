package main

import (
	"log"

	"github.com/dwikie/sentra-payment-orchestrator/config"
	"github.com/dwikie/sentra-payment-orchestrator/router"
	"github.com/gin-gonic/gin"
)

func main() {
	// Create a Gin router with default middleware (logger and recovery)
	r := gin.Default()

	err := config.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer config.Pool.Close()

	// Register additional routes
	router.Register(r)

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
