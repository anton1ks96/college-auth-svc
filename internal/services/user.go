package service

import (
	"context"
	"fmt"
	"time"

	"github.com/anton1ks96/college-auth-svc/internal/config"
	"github.com/anton1ks96/college-auth-svc/internal/domain"
	"github.com/anton1ks96/college-auth-svc/pkg/auth"
	"github.com/anton1ks96/college-auth-svc/pkg/logger"
)

type UserService struct {
	tokenManager    *auth.Manager
	repos           Repositories
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	cfg             *config.App
}

func NewUserService(tm auth.Manager, repos Repositories, accessTTL time.Duration, refreshTTL time.Duration, appCfg *config.App) *UserService {
	return &UserService{
		tokenManager:    &tm,
		repos:           repos,
		accessTokenTTL:  accessTTL,
		refreshTokenTTL: refreshTTL,
		cfg:             appCfg,
	}
}

func (u *UserService) SignIn(ctx context.Context, input SignInInput) (Tokens, *domain.User, error) {
	if ctx.Err() != nil {
		return Tokens{}, nil, ctx.Err()
	}

	if input.UserID == "" {
		logger.Error(fmt.Errorf("empty login username field"))
		return Tokens{}, nil, fmt.Errorf("empty login username field")
	}

	if input.Password == "" {
		logger.Error(fmt.Errorf("empty login password field"))
		return Tokens{}, nil, fmt.Errorf("empty login password field")
	}

	if ctx.Err() != nil {
		return Tokens{}, nil, ctx.Err()
	}

	var user *domain.User
	var err error

	if u.cfg.Test {
		user = &domain.User{
			ID:       "i99s9999",
			Username: "Vasiliy Testov",
			Role:     "student",
		}
	} else {
		if err := u.repos.UserRepo.Authentication(ctx, input.UserID, input.Password); err != nil {
			logger.Error(fmt.Errorf("authentication failed for user %s: %w", input.UserID, err))
			return Tokens{}, nil, fmt.Errorf("authentication failed: %w", err)
		}

		if ctx.Err() != nil {
			return Tokens{}, nil, ctx.Err()
		}

		user, err = u.repos.UserRepo.GetByID(ctx, input.UserID)
		if err != nil {
			logger.Error(fmt.Errorf("find user data failed for user %s: %w", input.UserID, err))
			return Tokens{}, nil, fmt.Errorf("find user data failed: %w", err)
		}

		if ctx.Err() != nil {
			return Tokens{}, nil, ctx.Err()
		}
	}

	tokens, err := u.generateTokens(user)
	if err != nil {
		logger.Error(fmt.Errorf("failed to generate tokens for user %s: %w", input.UserID, err))
		return Tokens{}, nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	jti, err := u.tokenManager.ExtractClaim(tokens.RefreshToken, "jti")
	if err != nil {
		logger.Error(fmt.Errorf("failed to extract jti from token for user %s: %w", input.UserID, err))
		return Tokens{}, nil, fmt.Errorf("failed to extract jti: %w", err)
	}

	if ctx.Err() != nil {
		return Tokens{}, nil, ctx.Err()
	}

	session := domain.RefreshSession{
		JTI:       jti,
		UserID:    input.UserID,
		ExpiresAt: time.Now().Add(u.refreshTokenTTL),
		CreatedAt: time.Now(),
	}

	if err := u.repos.SessionRepo.RevokeAllUserSessions(ctx, input.UserID); err != nil {
		logger.Error(fmt.Errorf("failed to revoke old sessions: %w", err))
		return Tokens{}, nil, err
	}

	if err := u.repos.SessionRepo.SaveRefreshToken(ctx, &session); err != nil {
		logger.Error(fmt.Errorf("failed to save refresh session for user %s: %w", input.UserID, err))
		return Tokens{}, nil, fmt.Errorf("failed to save refresh session: %w", err)
	}

	return Tokens{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, user, nil
}

func (u *UserService) SignOut(ctx context.Context, refreshToken string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if refreshToken == "" {
		logger.Error(fmt.Errorf("empty refresh token"))
		return fmt.Errorf("empty refresh token")
	}

	jti, err := u.tokenManager.ExtractClaim(refreshToken, "jti")
	if err != nil {
		logger.Error(fmt.Errorf("failed to extract jti from token: %w", err))
		return err
	}

	if ctx.Err() != nil {
		return ctx.Err()
	}

	if err := u.repos.SessionRepo.RevokeRefreshToken(ctx, jti); err != nil {
		logger.Error(fmt.Errorf("failed to revoke refresh token: %w", err))
		return err
	}

	return nil
}

func (u *UserService) RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error) {
	if ctx.Err() != nil {
		return Tokens{}, ctx.Err()
	}

	if refreshToken == "" {
		logger.Error(fmt.Errorf("empty refresh token"))
		return Tokens{}, fmt.Errorf("empty refresh token")
	}

	if err := u.tokenManager.ValidateRefreshToken(refreshToken); err != nil {
		logger.Warn(fmt.Sprintf("invalid refresh token structure: %v", err))
		return Tokens{}, fmt.Errorf("invalid token")
	}

	userID, err := u.tokenManager.ExtractClaim(refreshToken, "user_id")
	if err != nil {
		logger.Error(fmt.Errorf("failed to extract user_id from token: %w", err))
		return Tokens{}, fmt.Errorf("invalid token claims")
	}

	oldJti, err := u.tokenManager.ExtractClaim(refreshToken, "jti")
	if err != nil {
		logger.Error(fmt.Errorf("failed to extract jti from token: %w", err))
		return Tokens{}, fmt.Errorf("invalid token claims")
	}

	tokenExists, err := u.repos.SessionRepo.TokenExists(ctx, oldJti)
	if err != nil {
		logger.Error(fmt.Errorf("failed to check token existence: %w", err))
		return Tokens{}, fmt.Errorf("authentication service unavailable")
	}

	if !tokenExists {
		logger.Warn(fmt.Sprintf("attempt to use non-existent refresh token: jti=%s, user=%s",
			oldJti, userID))
		return Tokens{}, fmt.Errorf("token not found or already used")
	}

	var user *domain.User

	if u.cfg.Test {
		user = &domain.User{
			ID:       "i99s9999",
			Username: "Vasiliy Testov",
			Role:     "student",
		}
	} else {
		user, err = u.repos.UserRepo.GetByID(ctx, userID)
		if err != nil {
			logger.Error(fmt.Errorf("failed to get user data for ID %s: %w", userID, err))
			return Tokens{}, fmt.Errorf("failed to get user data: %w", err)
		}
	}

	tokens, err := u.generateTokens(user)
	if err != nil {
		logger.Error(fmt.Errorf("failed to generate tokens for user %s: %w", user.ID, err))
		return Tokens{}, fmt.Errorf("failed to generate tokens: %w", err)
	}

	newJti, err := u.tokenManager.ExtractClaim(tokens.RefreshToken, "jti")
	if err != nil {
		logger.Error(fmt.Errorf("failed to extract jti from new token for user %s: %w", userID, err))
		return Tokens{}, fmt.Errorf("failed to extract new token JTI: %w", err)
	}

	newSession := domain.RefreshSession{
		JTI:       newJti,
		UserID:    userID,
		ExpiresAt: time.Now().Add(u.refreshTokenTTL),
		CreatedAt: time.Now(),
	}

	if err := u.repos.SessionRepo.ReplaceRefreshToken(ctx, oldJti, &newSession); err != nil {
		logger.Error(fmt.Errorf("failed to replace refresh token: %w", err))
		return Tokens{}, fmt.Errorf("failed to rotate tokens")
	}

	return Tokens{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (u *UserService) ValidateAccessToken(ctx context.Context, accessToken string) (*domain.User, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if accessToken == "" {
		logger.Error(fmt.Errorf("empty access token"))
		return nil, fmt.Errorf("empty access token")
	}

	err := u.tokenManager.Validate(accessToken)
	if err != nil {
		logger.Error(fmt.Errorf("token expired"))
		return nil, fmt.Errorf("token expired")
	}

	userID, err := u.tokenManager.ExtractClaim(accessToken, "user_id")
	if err != nil {
		logger.Error(fmt.Errorf("failed to extract user_id from token: %w", err))
		return nil, err
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	var user *domain.User

	if u.cfg.Test {
		user = &domain.User{
			ID:       "i99s9999",
			Username: "Vasiliy Testov",
			Role:     "student",
		}
	} else {
		user, err = u.repos.UserRepo.GetByID(ctx, userID)
		if err != nil {
			logger.Error(fmt.Errorf("failed to get user data for ID %s: %w", userID, err))
			return nil, fmt.Errorf("failed to get user data: %w", err)
		}
	}

	return user, nil
}

func (u *UserService) generateTokens(user *domain.User) (Tokens, error) {
	newAccess, err := u.tokenManager.NewAccessToken(user.ID, user.Username, user.Role)
	if err != nil {
		logger.Error(fmt.Errorf("failed to generate access token for user %s: %w", user.ID, err))
		return Tokens{}, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefresh, err := u.tokenManager.NewRefreshToken(user.ID)
	if err != nil {
		logger.Error(fmt.Errorf("failed to generate refresh token for user %s: %w", user.ID, err))
		return Tokens{}, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return Tokens{
		AccessToken:  newAccess,
		RefreshToken: newRefresh,
	}, nil
}
