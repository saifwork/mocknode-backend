package mongo

import (
	"context"
	"log"
	"time"

	"github.com/saifwork/mock-service/internal/core/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InitMongo initializes the MongoDB client and sets up TTL indexes.
func InitMongo(cfg *config.Config) (*mongo.Client, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	log.Println("[MONGO] Connected successfully")

	// ✅ Setup TTL index (only once at startup)
	if err := setupTTLIndex(cfg, client); err != nil {
		return nil, err
	}

	return client, nil
}

// setupTTLIndex creates TTL index on "users" collection.
func setupTTLIndex(cfg *config.Config, client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := client.Database(cfg.MongoDBName) // 👈 change to your actual database name
	usersCol := db.Collection(Collections.Users)

	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "createdAt", Value: 1}}, // TTL on createdAt
		Options: options.Index().
			SetExpireAfterSeconds(24 * 60 * 60). // 24 hours
			SetPartialFilterExpression(bson.D{{Key: "isVerified", Value: false}}),
	}

	name, err := usersCol.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Printf("[MONGO] TTL index creation failed: %v", err)
		return err
	}

	log.Printf("[MONGO] TTL index created/exists: %s", name)

	return nil
}

// Helper to get DB handle cleanly
func GetDatabase(client *mongo.Client, dbName string) *mongo.Database {
	return client.Database(dbName)
}
