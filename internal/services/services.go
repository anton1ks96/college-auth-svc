package service

import (
	"context"

	"github.com/anton1ks96/college-auth-svc/internal/config"
	"github.com/anton1ks96/college-auth-svc/internal/domain"
	"github.com/anton1ks96/college-auth-svc/internal/repository"
	"github.com/anton1ks96/college-auth-svc/pkg/auth"
)

type SignInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Auth handles user authentication, token generation and validation
type Auth interface {
	SignIn(ctx context.Context, input SignInInput) (Tokens, error)
	SignOut(ctx context.Context, refreshToken string) error
	RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)
	ValidateAccessToken(ctx context.Context, token string) (*domain.User, error)
}

// Users manages user data retrieval and operations
type Users interface {
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
}

// Sessions manages refresh token sessions and their lifecycle
type Sessions interface {
	CreateSession(ctx context.Context, session *domain.RefreshSession) error
	GetSession(ctx context.Context, jti string) (*domain.RefreshSession, error)
	RevokeSession(ctx context.Context, jti string) error
}

// Services contains all business logic services
type Services struct {
	Auth    Auth
	User    Users
	Session Sessions
}

type Repositories struct {
	UserRepo    repository.UserRepository
	SessionRepo repository.SessionRepository
}

// Deps contains external dependencies needed to initialize services
type Deps struct {
	Repos        *Repositories
	TokenManager *auth.Manager
	Config       *config.Config
}

func NewServices(deps Deps) *Services {
	return &Services{
		Auth:    nil,
		User:    nil,
		Session: nil,
	}
}
