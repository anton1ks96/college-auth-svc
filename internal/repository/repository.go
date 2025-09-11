package repository

import (
	"context"

	"github.com/anton1ks96/college-auth-svc/internal/domain"
)

// UserLDAPRepository handles user authentication and data retrieval from LDAP
type UserLDAPRepository interface {
	Authentication(ctx context.Context, username, password string) error
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
}

// SessionMongoRepository manages refresh tokens and user sessions in MongoDB
type SessionMongoRepository interface {
	SaveRefreshToken(ctx context.Context, session *domain.RefreshSession) error
	RevokeRefreshToken(ctx context.Context, jti string) error
	// GetRefreshToken(ctx context.Context, jti string) (*domain.RefreshSession, error)
}
