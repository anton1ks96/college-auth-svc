package service

import (
	"context"
	"fmt"
	"time"

	"github.com/anton1ks96/college-auth-svc/internal/domain"
	"github.com/anton1ks96/college-auth-svc/pkg/auth"
	"github.com/anton1ks96/college-auth-svc/pkg/logger"
)

type UserService struct {
	tokenManager *auth.Manager
	repos        Repositories

	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewUserService(tm auth.Manager, repos Repositories, accessTTL time.Duration, refreshTTL time.Duration) *UserService {
	return &UserService{
		tokenManager:    &tm,
		repos:           repos,
		accessTokenTTL:  accessTTL,
		refreshTokenTTL: refreshTTL,
	}
}

func (u *UserService) SignIn(ctx context.Context, input SignInInput) (Tokens, *domain.User, error) {
	if input.Username == "" {
		logger.Error(fmt.Errorf("empty login username field"))
	}

	if input.Password == "" {
		logger.Error(fmt.Errorf("empty login password field"))
	}

	if err := u.repos.UserRepo.Authentication(ctx, input.Username, input.Password); err != nil {
		logger.Error(fmt.Errorf("authentication failed for user %s: %w", input.Username, err))
		return Tokens{}, nil, fmt.Errorf("authentication failed: %w", err)
	}

	user, err := u.repos.UserRepo.GetByUsername(ctx, input.Username)
	if err != nil {
		logger.Error(fmt.Errorf("find user data failed for user %s: %w", input.Username, err))
		return Tokens{}, nil, fmt.Errorf("find user data failed: %w", err)
	}

	accessToken, err := u.tokenManager.NewAccessToken(input.Username)
	if err != nil {
		logger.Error(fmt.Errorf("failed to generate access token for user %s: %w", input.Username, err))
		return Tokens{}, nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := u.tokenManager.NewRefreshToken(input.Username)
	if err != nil {
		logger.Error(fmt.Errorf("failed to generate refresh token for user %s: %w", input.Username, err))
		return Tokens{}, nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	jti, err := u.tokenManager.ExtractJTI(refreshToken)
	if err != nil {
		logger.Error(fmt.Errorf("failed to extract jti from token for user %s: %w", input.Username, err))
		return Tokens{}, nil, fmt.Errorf("failed to extract jti: %w", err)
	}

	session := domain.RefreshSession{
		JTI:       jti,
		Username:  input.Username,
		ExpiresAt: time.Now().Add(u.refreshTokenTTL),
		CreatedAt: time.Now(),
	}

	if err := u.repos.SessionRepo.SaveRefreshToken(ctx, &session); err != nil {
		logger.Error(fmt.Errorf("failed to save refresh session for user %s: %w", input.Username, err))
		return Tokens{}, nil, fmt.Errorf("failed to save refresh session: %w", err)
	}

	return Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, user, nil
}

//func (u *UserService) SignOut(ctx context.Context, refreshToken string) error {
//	return nil
//}
//
//func (u *UserService) RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error) {
//	return Tokens{
//		AccessToken:  "",
//		RefreshToken: "",
//	}, nil
//}
//
//func (u *UserService) ValidateAccessToken(ctx context.Context, token string) (*domain.User, error) {
//	return nil, nil
//}
//
//func (u *UserService) GetInfo(ctx context.Context, input SignInInput) (domain.User, error) {
//	//TODO implement me
//	panic("implement me")
//}
