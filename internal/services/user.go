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
	tokenManager    *auth.Manager
	repos           Repositories
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
		return Tokens{}, nil, fmt.Errorf("empty login username field")
	}

	if input.Password == "" {
		logger.Error(fmt.Errorf("empty login password field"))
		return Tokens{}, nil, fmt.Errorf("empty login password field")
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

	jti, err := u.tokenManager.ExtractClaim(refreshToken, "jti")
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

func (u *UserService) SignInTest(ctx context.Context, input SignInInput) (Tokens, *domain.User, error) {
	if input.Username == "" {
		logger.Error(fmt.Errorf("empty login username field"))
		return Tokens{}, nil, fmt.Errorf("empty login username field")
	}

	if input.Password == "" {
		logger.Error(fmt.Errorf("empty login password field"))
		return Tokens{}, nil, fmt.Errorf("empty login password field")
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

	jti, err := u.tokenManager.ExtractClaim(refreshToken, "jti")
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

	var user = domain.User{
		ID:       "123",
		Username: "qwe",
		Mail:     "qwe",
		Role:     "qwe",
	}
	return Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, &user, nil
}

func (u *UserService) SignOut(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		logger.Error(fmt.Errorf("empty refresh token"))
		return fmt.Errorf("empty refresh token")
	}

	jti, err := u.tokenManager.ExtractClaim(refreshToken, "jti")
	if err != nil {
		logger.Error(fmt.Errorf("failed to extract jti from token: %w", err))
		return err
	}

	if err := u.repos.SessionRepo.RevokeRefreshToken(ctx, jti); err != nil {
		logger.Error(fmt.Errorf("failed to revoke refresh token: %w", err))
		return err
	}

	return nil
}

func (u *UserService) RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error) {
	if refreshToken == "" {
		logger.Error(fmt.Errorf("empty refresh token"))
		return Tokens{}, fmt.Errorf("empty refresh token")
	}

	userName, err := u.tokenManager.ExtractClaim(refreshToken, "username")
	if err != nil {
		logger.Error(fmt.Errorf("failed to extract username from token: %w", err))
		return Tokens{}, err
	}

	newAccess, err := u.tokenManager.NewAccessToken(userName)
	if err != nil {
		logger.Error(fmt.Errorf("failed to generate access token for user %s: %w", userName, err))
		return Tokens{}, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefresh, err := u.tokenManager.NewRefreshToken(userName)
	if err != nil {
		logger.Error(fmt.Errorf("failed to generate refresh token for user %s: %w", userName, err))
		return Tokens{}, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	jti, err := u.tokenManager.ExtractClaim(newRefresh, "jti")
	if err != nil {
		logger.Error(fmt.Errorf("failed to extract jti from token for user %s: %w", userName, err))
		return Tokens{}, fmt.Errorf("failed to extract jti: %w", err)
	}

	session := domain.RefreshSession{
		JTI:       jti,
		Username:  userName,
		ExpiresAt: time.Now().Add(u.refreshTokenTTL),
		CreatedAt: time.Now(),
	}

	if err := u.repos.SessionRepo.SaveRefreshToken(ctx, &session); err != nil {
		logger.Error(fmt.Errorf("failed to save refresh session for user %s: %w", userName, err))
		return Tokens{}, fmt.Errorf("failed to save refresh session: %w", err)
	}

	return Tokens{
		AccessToken:  newAccess,
		RefreshToken: newRefresh,
	}, nil
}

func (u *UserService) ValidateAccessToken(ctx context.Context, accessToken string) (*domain.User, error) {
	if accessToken == "" {
		logger.Error(fmt.Errorf("empty access token"))
		return nil, fmt.Errorf("empty access token")
	}

	_, err := u.tokenManager.Validate(accessToken)
	if err != nil {
		logger.Error(fmt.Errorf("token is not valid"))
		return nil, fmt.Errorf("token is not valid")
	}

	userName, err := u.tokenManager.ExtractClaim(accessToken, "username")
	if err != nil {
		logger.Error(fmt.Errorf("failed to extract username from token: %w", err))
		return nil, err
	}

	// TODO: add role field

	user, err := u.repos.UserRepo.GetByUsername(ctx, userName)
	if err != nil {
		logger.Error(fmt.Errorf("find user data failed for user %s: %w", userName, err))
		return nil, fmt.Errorf("find user data failed: %w", err)
	}

	return user, nil
}

func (u *UserService) GetInfo(ctx context.Context, input SignInInput) (*domain.User, error) {
	user, err := u.repos.UserRepo.GetByUsername(ctx, input.Username)
	if err != nil {
		logger.Error(fmt.Errorf("find user data failed for user %s: %w", input.Username, err))
		return nil, fmt.Errorf("find user data failed: %w", err)
	}

	return user, nil
}
