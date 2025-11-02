package services

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"slices"
	"time"

	"github.com/saifwork/mock-service/internal/core/config"
	database "github.com/saifwork/mock-service/internal/core/mongo"
	"github.com/saifwork/mock-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RecordService struct {
	coll           *mongo.Collection
	collectioncoll *mongo.Collection
	ctx            context.Context
	cfg            *config.Config
}

func NewRecordService(client *mongo.Client, cfg *config.Config) *RecordService {
	collection := client.Database(cfg.MongoDBName).Collection(database.Collections.Records)
	collectioncoll := client.Database(cfg.MongoDBName).Collection(database.Collections.Collection)
	return &RecordService{
		coll:           collection,
		collectioncoll: collectioncoll,
		ctx:            context.Background(),
		cfg:            cfg,
	}
}

// CreateRecord adds a new record under a specific collection
func (s *RecordService) CreateRecord(collectionID string, data map[string]interface{}) (*models.Record, error) {
	cid, err := primitive.ObjectIDFromHex(collectionID)
	if err != nil {
		return nil, errors.New("invalid collection id")
	}

	// Ensure the collection exists
	var collection models.Collection
	err = s.collectioncoll.FindOne(context.Background(), bson.M{"_id": cid}).Decode(&collection)
	if err != nil {
		return nil, errors.New("collection not found")
	}

	if err := validateRecordData(collection.Fields, data); err != nil {
		return nil, err
	}

	record := &models.Record{
		ID:           primitive.NewObjectID(),
		CollectionID: cid,
		Data:         data,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	_, err = s.coll.InsertOne(context.Background(), record)
	if err != nil {
		return nil, err
	}

	return record, nil
}

// GetRecordsByCollection fetches all records of a collection
func (s *RecordService) GetRecordsByCollection(collectionID string) ([]models.Record, error) {
	cid, err := primitive.ObjectIDFromHex(collectionID)
	if err != nil {
		return nil, errors.New("invalid collection id")
	}

	cursor, err := s.coll.Find(context.Background(), bson.M{"collectionId": cid})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var records []models.Record
	if err := cursor.All(context.Background(), &records); err != nil {
		return nil, err
	}

	return records, nil
}

// GetRecordByID returns a single record
func (s *RecordService) GetRecordByID(id string) (*models.Record, error) {
	rid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid record id")
	}

	var record models.Record
	err = s.coll.FindOne(context.Background(), bson.M{"_id": rid}).Decode(&record)
	if err != nil {
		return nil, errors.New("record not found")
	}

	return &record, nil
}

// UpdateRecord updates record data after validating it against collection fields
func (s *RecordService) UpdateRecord(id string, data map[string]interface{}) (*models.Record, error) {
	rid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid record id")
	}

	// Fetch the existing record
	var existingRecord models.Record
	err = s.coll.FindOne(context.Background(), bson.M{"_id": rid}).Decode(&existingRecord)
	if err != nil {
		return nil, errors.New("record not found")
	}

	// Fetch the associated collection
	var collection models.Collection
	err = s.collectioncoll.FindOne(context.Background(), bson.M{"_id": existingRecord.CollectionID}).Decode(&collection)
	if err != nil {
		return nil, errors.New("associated collection not found")
	}

	// Validate updated data
	if err := validateRecordData(collection.Fields, data); err != nil {
		return nil, err
	}

	// Apply update
	update := bson.M{
		"$set": bson.M{
			"data":      data,
			"updatedAt": time.Now(),
		},
	}

	res := s.coll.FindOneAndUpdate(
		context.Background(),
		bson.M{"_id": rid},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	if res.Err() != nil {
		return nil, res.Err()
	}

	var updated models.Record
	if err := res.Decode(&updated); err != nil {
		return nil, err
	}

	return &updated, nil
}

// DeleteRecord removes a record
func (s *RecordService) DeleteRecord(id string) error {
	rid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid record id")
	}

	res, err := s.coll.DeleteOne(context.Background(), bson.M{"_id": rid})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("record not found")
	}

	return nil
}

var emailRegex = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

func validateRecordData(fields []models.FieldDefinition, data map[string]interface{}) error {
	for _, field := range fields {
		val, exists := data[field.Name]

		if field.Required && !exists {
			return fmt.Errorf("missing required field: %s", field.Name)
		}

		if !exists {
			continue
		}

		switch field.Type {
		case "string":
			str, ok := val.(string)
			if !ok {
				return fmt.Errorf("field %s must be a string", field.Name)
			}
			if field.MinLength != nil && len(str) < *field.MinLength {
				return fmt.Errorf("field %s must be at least %d characters", field.Name, *field.MinLength)
			}
			if field.MaxLength != nil && len(str) > *field.MaxLength {
				return fmt.Errorf("field %s must be at most %d characters", field.Name, *field.MaxLength)
			}
			if field.Pattern != nil {
				match, _ := regexp.MatchString(*field.Pattern, str)
				if !match {
					return fmt.Errorf("field %s does not match required pattern", field.Name)
				}
			}
			if field.Type == "email" {
				if !emailRegex.MatchString(str) {
					return fmt.Errorf("field %s must be a valid email", field.Name)
				}
			}
			if field.Type == "url" {
				_, err := url.ParseRequestURI(str)
				if err != nil {
					return fmt.Errorf("field %s must be a valid URL", field.Name)
				}
			}

		case "number":
			num, ok := val.(float64)
			if !ok {
				return fmt.Errorf("field %s must be a number", field.Name)
			}
			if field.MinValue != nil && num < *field.MinValue {
				return fmt.Errorf("field %s must be >= %f", field.Name, *field.MinValue)
			}
			if field.MaxValue != nil && num > *field.MaxValue {
				return fmt.Errorf("field %s must be <= %f", field.Name, *field.MaxValue)
			}

		case "boolean":
			if _, ok := val.(bool); !ok {
				return fmt.Errorf("field %s must be a boolean", field.Name)
			}

		case "array":
			if _, ok := val.([]interface{}); !ok {
				return fmt.Errorf("field %s must be an array", field.Name)
			}

		case "object":
			if _, ok := val.(map[string]interface{}); !ok {
				return fmt.Errorf("field %s must be an object", field.Name)
			}

		case "enum":
			str, ok := val.(string)
			if !ok {
				return fmt.Errorf("field %s must be a string for enum type", field.Name)
			}
			valid := slices.Contains(field.EnumValues, str)
			if !valid {
				return fmt.Errorf("field %s must be one of %v", field.Name, field.EnumValues)
			}
		}
	}
	return nil
}
