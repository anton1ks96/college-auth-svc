package repository

import (
	"context"

	"github.com/anton1ks96/college-auth-svc/internal/domain"
)

// UserLDAPRepository handles user authentication and data retrieval from LDAP
type UserLDAPRepository interface {
	Authentication(ctx context.Context, userID, userPass string) error
	GetByID(ctx context.Context, userID, userPass string) (*domain.User, error)
	GetUserGroups(ctx context.Context, userID, userPass string) (academicGroup, profile string, err error)
}

// SessionMongoRepository manages refresh tokens and user sessions in MongoDB
type SessionMongoRepository interface {
	SaveRefreshToken(ctx context.Context, session *domain.RefreshSession) error
	RevokeRefreshToken(ctx context.Context, jti string) error
	RevokeAllUserSessions(ctx context.Context, userID string) error
	TokenExists(ctx context.Context, jti string) (bool, error)
	ReplaceRefreshToken(ctx context.Context, oldJTI string, newSession *domain.RefreshSession) error
	GetUserByID(ctx context.Context, userID string) (*domain.User, error)
}
