package config

import (
	"log"

	"github.com/spf13/viper"
)

func LoadEnv() {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.SetDefault("DOMAIN", "localhost")
	viper.SetDefault("DATABASE_URL", "postgres://user:password@localhost:5432/mydb?sslmode=disable")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("No .env file found, relying on OS env vars/defaults")
		} else {
			log.Fatalf("fatal error config file: %v", err)
		}
	}

	log.Println("environment variables loaded successfully")
}
