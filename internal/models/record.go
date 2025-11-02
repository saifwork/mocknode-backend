package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Record struct {
	ID           primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	CollectionID primitive.ObjectID     `bson:"collectionId" json:"collectionId"`
	Data         map[string]interface{} `bson:"data" json:"data"`
	CreatedAt    time.Time              `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time              `bson:"updatedAt" json:"updatedAt"`
}
