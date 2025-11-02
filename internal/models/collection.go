package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FieldDefinition struct {
	Name        string   `bson:"name" json:"name"`
	Description string   `bson:"description,omitempty" json:"description,omitempty"`
	Type        string   `bson:"type" json:"type"` // string, number, boolean, array, object, date, email, url, enum
	Required    bool     `bson:"required" json:"required"`
	MinLength   *int     `bson:"minLength,omitempty" json:"minLength,omitempty"` // for strings
	MaxLength   *int     `bson:"maxLength,omitempty" json:"maxLength,omitempty"`
	MinValue    *float64 `bson:"minValue,omitempty" json:"minValue,omitempty"` // for numbers
	MaxValue    *float64 `bson:"maxValue,omitempty" json:"maxValue,omitempty"`
	Pattern     *string  `bson:"pattern,omitempty" json:"pattern,omitempty"` // regex for strings
	EnumValues  []string `bson:"enumValues,omitempty" json:"enumValues,omitempty"`
	Default     any      `bson:"default,omitempty" json:"default,omitempty"`
}

type Collection struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ProjectID primitive.ObjectID `bson:"projectId" json:"projectId"`
	Name      string             `bson:"name" json:"name"`
	Fields    []FieldDefinition  `bson:"fields" json:"fields"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}
