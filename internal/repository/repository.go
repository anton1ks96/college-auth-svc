package repository

import (
	"context"

	"github.com/anton1ks96/college-auth-svc/internal/domain"
)

// UserRepository handles user authentication and data retrieval from LDAP
type UserRepository interface {
	Authentication(ctx context.Context, username, password string) error
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
}

// SessionRepository manages refresh tokens and user sessions in MongoDB
type SessionRepository interface {
	SaveRefreshToken(ctx context.Context, session *domain.RefreshSession) error
	GetRefreshToken(ctx context.Context, username string) (*domain.RefreshSession, error)
	RefreshTokens(ctx context.Context)
	RevokeRefreshToken(ctx context.Context, tokenHash string) error
}
