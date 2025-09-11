package domain

import (
	"time"
)

type RefreshSession struct {
	JTI       string    `json:"jti" bson:"jti"`
	Username  string    `json:"username" bson:"username"`
	ExpiresAt time.Time `json:"expires_at" bson:"expires_at"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
