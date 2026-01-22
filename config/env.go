package config

import (
	"log"

	"github.com/spf13/viper"
)

func LoadEnv() {
	// load .env file
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Println("No .env file found, relying on OS env vars/defaults")
		} else {
			// Handle other errors
			log.Fatalf("fatal error config file: %v", err)
		}
	}

	log.Println("environment variables loaded successfully")
}
