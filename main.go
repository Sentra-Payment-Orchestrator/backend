package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	cfg, err := InitConfig()
	if err != nil {
		log.Fatal(err)
	}
	defer cfg.Pool.Close()

	env := viper.GetString("APP_ENV")
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.Use(func(c *gin.Context) {
		c.Next()
	})

	InitializeHandlers(cfg.Pool)
	RegisterRoutes(r, cfg.Pool)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
