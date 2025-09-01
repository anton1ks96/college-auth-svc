package repository

import (
	"context"

	"github.com/anton1ks96/college-auth-svc/internal/domain"
)

type UserRepository interface {
	Authenticate(ctx context.Context, username, password string) (*domain.User, error)
	//GetByUsername(ctx context.Context, username string) (*domain.User, error)
	//List(ctx context.Context, filters UserFilters) ([]*User, error)
}

type SessionRepository interface {
	SaveRefreshToken(ctx context.Context, session *domain.RefreshSession) error
	GetRefreshToken(ctx context.Context, tokenHash string) (*domain.RefreshSession, error)
	RevokeRefreshToken(ctx context.Context, tokenHash string) error
}
