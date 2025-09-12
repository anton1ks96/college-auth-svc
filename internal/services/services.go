package service

import (
	"context"
	"fmt"
	"time"

	"github.com/anton1ks96/college-auth-svc/internal/config"
	"github.com/anton1ks96/college-auth-svc/internal/domain"
	"github.com/anton1ks96/college-auth-svc/internal/repository"
	"github.com/anton1ks96/college-auth-svc/pkg/auth"
	"github.com/anton1ks96/college-auth-svc/pkg/logger"
)

type SignInInput struct {
	UserID   string `json:"userid" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type User interface {
	SignIn(ctx context.Context, input SignInInput) (Tokens, *domain.User, error)
	SignOut(ctx context.Context, refreshToken string) error
	RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)
	ValidateAccessToken(ctx context.Context, token string) (*domain.User, error)
}

type Services struct {
	UserService User
}

type Repositories struct {
	UserRepo    repository.UserLDAPRepository
	SessionRepo repository.SessionMongoRepository
}

type Deps struct {
	Repos        *Repositories
	TokenManager *auth.Manager
	Config       *config.Config
}

func NewServices(deps Deps) *Services {
	accessTTL, err := time.ParseDuration(deps.Config.JWT.AccessTokenTTL)
	if err != nil {
		logger.Fatal(fmt.Errorf("invalid access token TTL: %w", err))
	}

	refreshTTL, err := time.ParseDuration(deps.Config.JWT.RefreshTokenTTL)
	if err != nil {
		logger.Fatal(fmt.Errorf("invalid refresh token TTL: %w", err))
	}

	userService := NewUserService(*deps.TokenManager, *deps.Repos, accessTTL, refreshTTL, &deps.Config.App)

	return &Services{
		UserService: userService,
	}
}
