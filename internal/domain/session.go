package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	TokenHash string             `json:"tokenHash" bson:"token_hash"`
	UserID    string             `json:"userId" bson:"user_id"`
	CreatedAt time.Time          `json:"createdAt" bson:"created_at"`
	ExpiresAt time.Time          `json:"expiresAt" bson:"expires_at"`
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}
