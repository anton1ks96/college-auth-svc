package repository

import (
	"context"

	"github.com/anton1ks96/college-auth-svc/internal/domain"
)

type Repositories struct {
	UserRepo    UserRepository
	SessionRepo SessionRepository
}

// UserRepository handles user authentication and data retrieval from LDAP
type UserRepository interface {
	Authentication(ctx context.Context, username, password string) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	GetUserRoles(ctx context.Context, username string) ([]string, error)
}

// SessionRepository manages refresh tokens and user sessions in MongoDB
type SessionRepository interface {
	SaveRefreshToken(ctx context.Context, session *domain.RefreshSession) error
	GetRefreshToken(ctx context.Context, tokenHash string) (*domain.RefreshSession, error)
	RevokeRefreshToken(ctx context.Context, tokenHash string) error
}
