package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/saifwork/mock-service/internal/core/config"
)

var (
	Rdb *redis.Client
	ctx = context.Background()
)

// InitRedis initializes a Redis client and verifies the connection
func InitRedis(cfg *config.Config) {
	addr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)

	Rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.RedisPassword, // empty string means no password
		DB:       0,                 // default DB
	})

	// Test connection
	pong, err := Rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("❌ Failed to connect to Redis at %s: %v", addr, err)
	}

	log.Printf("✅ Redis connected: %s", pong)
}

// CloseRedis closes the Redis connection (useful for graceful shutdown)
func CloseRedis() {
	if Rdb != nil {
		if err := Rdb.Close(); err != nil {
			log.Printf("⚠️ Error closing Redis: %v", err)
		} else {
			log.Println("🧹 Redis connection closed.")
		}
	}
}

// GetClient returns the global Redis client
func GetClient() *redis.Client {
	return Rdb
}
