package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user account stored in MongoDB.
type User struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FullName            string             `bson:"fullName" json:"fullName"`
	Email               string             `bson:"email" json:"email"`
	PasswordHash        string             `bson:"passwordHash,omitempty" json:"-"`
	IsUpgraded          bool               `bson:"isUpgraded" json:"isUpgraded"`
	IsVerified          bool               `bson:"isVerified" json:"isVerified"`
	VerificationToken   string             `bson:"verificationToken,omitempty" json:"-"`
	VerificationExpires time.Time          `bson:"verificationExpires,omitempty" json:"-"` // 👈 new field
	IsActive            bool               `bson:"isActive" json:"isActive"`
	CreatedAt           time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt           time.Time          `bson:"updatedAt" json:"updatedAt"`
}
