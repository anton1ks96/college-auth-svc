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

// User handles user authentication, token generation and validation
type User interface {
	SignIn(ctx context.Context, input SignInInput) (Tokens, error)
	SignOut(ctx context.Context, refreshToken string) error
	RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)
	ValidateAccessToken(ctx context.Context, token string) (*domain.User, error)
	GetInfo(ctx context.Context, input SignInInput) (domain.User, error)
}

// Services contains all business logic services
type Services struct {
	UserService User
}

type Repositories struct {
	UserRepo    repository.UserLDAPRepository
	SessionRepo repository.SessionMongoRepository
}

// Deps contains external dependencies needed to initialize services
type Deps struct {
	Repos        *Repositories
	TokenManager *auth.Manager
	Config       *config.Config
}

func NewServices(deps Deps) *Services {
	return &Services{
		UserService: nil,
	}
}
