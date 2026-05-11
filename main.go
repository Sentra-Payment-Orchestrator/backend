package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type App struct {
	Pool     *pgxpool.Pool
	Redis    *redis.Client
	Handlers *Handlers
}

func main() {
	app, err := InitializeApp()
	if err != nil {
		log.Fatal(err)
	}

	defer app.Pool.Close()

	env := viper.GetString("APP_ENV")
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.Use(func(c *gin.Context) {
		c.Next()
	})

	app.InitializeHandlers(app.Pool, app.Redis)
	app.RegisterRoutes(r)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
