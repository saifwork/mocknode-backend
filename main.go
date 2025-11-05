package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/saifwork/mock-service/internal/api"
	"github.com/saifwork/mock-service/internal/api/handlers"
	"github.com/saifwork/mock-service/internal/core/config"
	"github.com/saifwork/mock-service/internal/core/mongo"

	// redisClient "github.com/saifwork/mock-service/internal/core/redis"
	"github.com/saifwork/mock-service/internal/middlewares"
	"github.com/saifwork/mock-service/internal/services"
)

func main() {

	_ = godotenv.Load()

	val := os.Getenv("MONGO_URI")
	log.Printf("[DEBUG] Raw os.Getenv(MONGO_URI): '%s'\n", val)
	if val == "" {
		log.Println("[ERROR] MONGO_URI is empty! Check .env quoting or formatting.")
	}

	// --- Create global app context ---
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// --- Load configuration ---
	cfg := config.LoadConfig()

	// --- Init Redis ---
	// redisClient.InitRedis(cfg)
	// defer redisClient.CloseRedis()

	// --- Initialize Mongo ---
	mongoClient, err := mongo.InitMongo(cfg)
	if err != nil {
		log.Fatalf("Failed to connect MongoDB: %v", err)
	}
	defer func() {
		if err := mongoClient.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting MongoDB: %v", err)
		}
	}()

	authSvc := services.NewAuthService(mongoClient, cfg)
	projectSvc := services.NewProjectService(mongoClient, cfg)
	collectionSvc := services.NewCollectionService(mongoClient, cfg)
	recordSvc := services.NewRecordService(mongoClient, cfg)
	configSvc := services.NewConfigService()

	// init handlers
	authHandler := handlers.NewAuthHandler(authSvc, cfg)
	projectHandler := handlers.NewProjectHandler(projectSvc, cfg)
	collectionHandler := handlers.NewCollectionHandler(collectionSvc, cfg)
	recordHandler := handlers.NewRecordHandler(recordSvc, cfg)
	configHandler := handlers.NewConfigHandler(configSvc, cfg)
	healthHandler := handlers.NewHealthHandler(mongoClient)

	// --- Initialize Gin ---
	r := gin.New() // Use New() instead of Default() to control middleware order
	r.Use(
		gin.Recovery(),       // Gin recovery for panics
		middlewares.Logger(), // Custom request logger
		middlewares.CORS(),   // Handle cross-origin requests
	)

	// --- Register routes ---
	api.RegisterRoutes(r, cfg, authHandler, projectHandler, collectionHandler, recordHandler, healthHandler, configHandler)

	// --- Start server ---
	log.Printf("Starting %s on port %s...", cfg.AppName, cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
