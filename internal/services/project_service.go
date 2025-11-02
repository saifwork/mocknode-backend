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

type ProjectService struct {
	coll     *mongo.Collection
	usercoll *mongo.Collection
	ctx      context.Context
	cfg      *config.Config
}

func NewProjectService(client *mongo.Client, cfg *config.Config) *ProjectService {
	collection := client.Database(cfg.MongoDBName).Collection(database.Collections.Projects)
	usercollection := client.Database(cfg.MongoDBName).Collection(database.Collections.Users)
	return &ProjectService{
		coll:     collection,
		usercoll: usercollection,
		ctx:      context.Background(),
		cfg:      cfg,
	}
}

func (s *ProjectService) CreateProject(userID, name, desc string) (*models.Project, error) {
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	// 🧠 Step 1: Fetch user info
	var user models.User
	err = s.usercoll.FindOne(context.Background(), bson.M{"_id": uid}).Decode(&user)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// 🧩 Step 2: If not upgraded, check project count
	if !user.IsUpgraded {
		count, err := s.coll.CountDocuments(context.Background(), bson.M{"userId": uid})
		if err != nil {
			return nil, err
		}
		if count >= 2 {
			return nil, errors.New("upgrade required — you can only create 2 projects with a free account")
		}
	}

	// ✅ Step 3: Proceed to create project
	project := &models.Project{
		ID:          primitive.NewObjectID(),
		UserID:      uid,
		Name:        name,
		Description: desc,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err = s.coll.InsertOne(context.Background(), project)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (s *ProjectService) GetProjectByID(pid, userID string) (*models.Project, error) {
	oid, err := primitive.ObjectIDFromHex(pid)
	if err != nil {
		return nil, errors.New("invalid project id")
	}
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	var project models.Project
	err = s.coll.FindOne(context.Background(), bson.M{"_id": oid, "userId": uid}).Decode(&project)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("project not found or unauthorized")
		}
		return nil, err
	}

	return &project, nil
}

func (s *ProjectService) GetUserProjects(userID string) ([]models.Project, error) {
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	cur, err := s.coll.Find(context.Background(), bson.M{"userId": uid})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	var projects []models.Project
	if err := cur.All(context.Background(), &projects); err != nil {
		return nil, err
	}

	return projects, nil
}

func (s *ProjectService) UpdateProject(pid, userID, name, desc string) error {
	oid, err := primitive.ObjectIDFromHex(pid)
	if err != nil {
		return errors.New("invalid project id")
	}
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user id")
	}

	update := bson.M{
		"$set": bson.M{
			"name":        name,
			"description": desc,
			"updatedAt":   time.Now(),
		},
	}

	res, err := s.coll.UpdateOne(context.Background(), bson.M{"_id": oid, "userId": uid}, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("no project found or unauthorized")
	}

	return nil
}

func (s *ProjectService) DeleteProject(pid, userID string) error {
	oid, err := primitive.ObjectIDFromHex(pid)
	if err != nil {
		return errors.New("invalid project id")
	}
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user id")
	}

	res, err := s.coll.DeleteOne(context.Background(), bson.M{"_id": oid, "userId": uid})
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("no project found or unauthorized")
	}

	return nil
}
