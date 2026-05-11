package config

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func InitRedis() (*redis.Client, error) {
	redisURL := "redis://dev:secretpassword@0.0.0.0:6379/0"
	if redisURL == "" {
		return nil, fmt.Errorf("REDIS_URL is not set in environment variables")
	}

	options, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing REDIS_URL: %v", err)
	}

	client := redis.NewClient(options)

	// Test the connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("error connecting to Redis: %v", err)
	}

	return client, nil
}
