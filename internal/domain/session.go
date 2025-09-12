package domain

import (
	"time"
)

type RefreshSession struct {
	JTI       string    `json:"jti" bson:"jti"`
	UserID    string    `json:"userid" bson:"userid"`
	ExpiresAt time.Time `json:"expires_at" bson:"expires_at"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
