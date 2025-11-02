package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all environment-based configuration
type Config struct {
	AppName         string
	AppPort         string
	AppEnv          string
	RedisHost       string
	RedisPort       string
	RedisPassword   string
	SessionTTL      time.Duration
	SessionReqLimit int
	// 🧩 MongoDB
	MongoURI    string
	MongoDBName string

	JWTSecret  string
	AppBaseURL string

	// SMTP
	GmailUser    string
	GmailPassKey string
}

// LoadConfig loads environment variables into the Config struct
func LoadConfig() *Config {
	// Load .env file if it exists
	_ = godotenv.Load()

	ttlSeconds := getEnvAsInt("SESSION_TTL_SECONDS", 604800) // default 7 days
	reqLimit := getEnvAsInt("SESSION_REQUEST_LIMIT", 500)

	cfg := &Config{
		AppName:         getEnv("APP_NAME", "Mock CRUD Service"),
		AppPort:         getEnv("APP_PORT", "8080"),
		AppEnv:          getEnv("APP_ENV", "development"),
		RedisHost:       getEnv("REDIS_HOST", "localhost"),
		RedisPort:       getEnv("REDIS_PORT", "6379"),
		RedisPassword:   getEnv("REDIS_PASSWORD", ""),
		SessionTTL:      time.Duration(ttlSeconds) * time.Second,
		SessionReqLimit: reqLimit,
		MongoURI:        getEnv("MONGO_URI", "mongodb+srv://mocknode_db_user:A3gl7aRY0ILpHF2i@cluster0.i3ygpz3.mongodb.net/?retryWrites=true&w=majority&tls=true"),
		MongoDBName:     getEnv("MONGO_DB_NAME", "mock_service"),
		JWTSecret:       getEnv("JWT_SECRET", "supersecret"),
		AppBaseURL:      getEnv("APP_BASE_URL", "http://localhost:8080"),
		GmailUser:       getEnv("GMAIL_USER", "mocknode@gmail.com"),
		GmailPassKey:    getEnv("GMAIL_PASSKEY", "pgtlcjrinxvywdss"),
	}

	log.Println("======== [CONFIG LOADED SUCCESSFULLY] ========")
	log.Printf("APP_NAME: %s", cfg.AppName)
	log.Printf("APP_PORT: %s", cfg.AppPort)
	log.Printf("APP_ENV: %s", cfg.AppEnv)
	log.Printf("REDIS_HOST: %s", cfg.RedisHost)
	log.Printf("REDIS_PORT: %s", cfg.RedisPort)
	log.Printf("REDIS_PASSWORD: %s", cfg.RedisPassword)
	log.Printf("SESSION_TTL: %v", cfg.SessionTTL)
	log.Printf("SESSION_REQUEST_LIMIT: %d", cfg.SessionReqLimit)
	log.Printf("MONGO_URI: %s", cfg.MongoURI)
	log.Printf("MONGO_DB_NAME: %s", cfg.MongoDBName)
	log.Printf("JWT_SECRET: %s", cfg.JWTSecret)
	log.Printf("APP_BASE_URL: %s", cfg.AppBaseURL)
	log.Printf("GMAIL_USER: %s", cfg.GmailUser)
	log.Printf("GMAIL_PASSKEY: %s", cfg.GmailPassKey)
	log.Println("=============================================")

	log.Printf("[CONFIG] Session TTL: %v | Request Limit: %d\n", cfg.SessionTTL, cfg.SessionReqLimit)
	return cfg
}

// Helpers
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if valueStr, exists := os.LookupEnv(key); exists {
		value, err := strconv.Atoi(valueStr)
		if err == nil {
			return value
		}
		log.Printf("[CONFIG] Invalid int for %s: %v (using default %d)\n", key, err, fallback)
	}
	return fallback
}
