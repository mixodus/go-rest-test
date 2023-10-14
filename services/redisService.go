package services

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var client *redis.Client

// InitializeRedis initializes the Redis client.
func InitializeRedis() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // No password
		DB:       0,                // Default database
	})

	c := &gin.Context{}
	_, err := client.Ping(c).Result()
	if err != nil {
		log.Fatal(err)
	}
}

// GetClient returns the Redis client for use in other parts of the application.
func GetRedisClient() *redis.Client {
	return client
}
