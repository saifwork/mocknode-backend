package services

import (
	"context"
	"errors"
	"time"

	"github.com/saifwork/mock-service/internal/core/config"
	database "github.com/saifwork/mock-service/internal/core/mongo"
	"github.com/saifwork/mock-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CollectionService struct {
	coll        *mongo.Collection
	projectColl *mongo.Collection
	userColl    *mongo.Collection
	ctx         context.Context
	cfg         *config.Config
}

func NewCollectionService(client *mongo.Client, cfg *config.Config) *CollectionService {
	collection := client.Database(cfg.MongoDBName).Collection(database.Collections.Collection)
	projectcollection := client.Database(cfg.MongoDBName).Collection(database.Collections.Projects)
	usercollection := client.Database(cfg.MongoDBName).Collection(database.Collections.Users)
	return &CollectionService{
		coll:        collection,
		projectColl: projectcollection,
		userColl:    usercollection,
		ctx:         context.Background(),
		cfg:         cfg,
	}
}

// CreateCollection creates a new collection under a specific project
func (s *CollectionService) CreateCollection(projectID, name string, fields []models.FieldDefinition) (*models.Collection, error) {
	pid, err := primitive.ObjectIDFromHex(projectID)
	if err != nil {
		return nil, errors.New("invalid project id")
	}

	ctx := context.Background()

	// ✅ Ensure project exists
	var project models.Project
	err = s.projectColl.FindOne(ctx, bson.M{"_id": pid}).Decode(&project)
	if err != nil {
		return nil, errors.New("project not found")
	}

	// ✅ Fetch user to check upgrade status
	var user models.User
	err = s.userColl.FindOne(ctx, bson.M{"_id": project.UserID}).Decode(&user)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// ✅ Count how many collections this user has under this project
	count, err := s.coll.CountDocuments(ctx, bson.M{"projectId": pid})
	if err != nil {
		return nil, err
	}

	// ⚠️ Restrict non-upgraded users to max 3 collections per project
	if !user.IsUpgraded && count >= 3 {
		return nil, errors.New("free users can only create up to 3 collections per project — upgrade to add more")
	}

	// ✅ Create collection
	collection := &models.Collection{
		ID:        primitive.NewObjectID(),
		ProjectID: pid,
		Name:      name,
		Fields:    fields,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = s.coll.InsertOne(ctx, collection)
	if err != nil {
		return nil, err
	}

	return collection, nil
}

// GetCollectionsByProject returns all collections under a project
func (s *CollectionService) GetCollectionsByProject(projectID string) ([]models.Collection, error) {
	pid, err := primitive.ObjectIDFromHex(projectID)
	if err != nil {
		return nil, errors.New("invalid project id")
	}

	cursor, err := s.coll.Find(context.Background(), bson.M{"projectId": pid})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var collections []models.Collection
	if err := cursor.All(context.Background(), &collections); err != nil {
		return nil, err
	}

	return collections, nil
}

// GetCollectionByID returns a specific collection
func (s *CollectionService) GetCollectionByID(id string) (*models.Collection, error) {
	cid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid collection id")
	}

	var collection models.Collection
	err = s.coll.FindOne(context.Background(), bson.M{"_id": cid}).Decode(&collection)
	if err != nil {
		return nil, errors.New("collection not found")
	}

	return &collection, nil
}

// DeleteCollection removes a collection
func (s *CollectionService) DeleteCollection(id string) error {
	cid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid collection id")
	}

	res, err := s.coll.DeleteOne(context.Background(), bson.M{"_id": cid})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("collection not found")
	}

	return nil
}
