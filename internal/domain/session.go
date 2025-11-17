package domain

import (
	"time"
)

type RefreshSession struct {
	JTI           string    `json:"jti" bson:"jti"`
	UserID        string    `json:"userid" bson:"userid"`
	Username      string    `json:"username" bson:"username"`
	Role          string    `json:"role" bson:"role"`
	AcademicGroup string    `json:"academic_group,omitempty" bson:"academic_group,omitempty"`
	Profile       string    `json:"profile,omitempty" bson:"profile,omitempty"`
	Subgroup      string    `json:"subgroup,omitempty" bson:"subgroup,omitempty"`
	EnglishGroup  string    `json:"english_group,omitempty" bson:"english_group,omitempty"`
	ExpiresAt     time.Time `json:"expires_at" bson:"expires_at"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at"`
}
