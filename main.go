package main

import (
	"log"

	"github.com/dwikie/sentra-payment-orchestrator/config"
	"github.com/gin-gonic/gin"
	"github.com/o1egl/paseto"
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

	token, err := app.Handlers.Auth.CreateToken([]byte("f9c346af754b4e0b0a7417ee9546c47c"), "access")
	if err != nil {
		log.Fatalf("Failed to create token: %v", err)
	}
	log.Printf("Generated Token: %s", token)

	var newJsonToken paseto.JSONToken
	var newFooter string
	err = paseto.NewV2().Decrypt(token, []byte("f9c346af754b4e0b0a7417ee9546c47c"), &newJsonToken, &newFooter)

	if err != nil {
		log.Fatalf("Failed to decrypt token: %v", err)
	}
	log.Printf("Decrypted Token Claims: %v", newJsonToken)
	log.Printf("Footer", newFooter)

	app.RegisterRoutes(r)

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
