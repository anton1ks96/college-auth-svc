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

type AppUser interface {
	SignIn(ctx context.Context, input SignInInput) (Tokens, *domain.UserExtended, error)
	SignOut(ctx context.Context, refreshToken string) error
	RefreshToken(ctx context.Context, refreshToken string) (string, error)
	GetAccessToken(ctx context.Context, refreshToken string) (string, *domain.UserExtended, error)
	ValidateAccessToken(ctx context.Context, token string) (*domain.UserExtended, error)
}

type AppUserService struct {
	tokenManager    *auth.Manager
	repos           Repositories
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	cfg             *config.App
}

func NewAppUserService(tm auth.Manager, repos Repositories, accessTTL time.Duration, refreshTTL time.Duration, appCfg *config.App) *AppUserService {
	return &AppUserService{
		tokenManager:    &tm,
		repos:           repos,
		accessTokenTTL:  accessTTL,
		refreshTokenTTL: refreshTTL,
		cfg:             appCfg,
	}
}

func (a *AppUserService) SignIn(ctx context.Context, input SignInInput) (Tokens, *domain.UserExtended, error) {
	if ctx.Err() != nil {
		return Tokens{}, nil, ctx.Err()
	}

	if input.UserID == "" || input.Password == "" {
		logger.Error(fmt.Errorf("empty login credentials"))
		return Tokens{}, nil, fmt.Errorf("empty login credentials")
	}

	var user *domain.User
	var userExtended *domain.UserExtended
	var err error

	if a.cfg.Test {
		userExtended = &domain.UserExtended{
			ID:            input.UserID,
			Username:      "i24s0291",
			Role:          "student",
			AcademicGroup: "ИТ25-11",
			Profile:       "BE",
		}
		user = &domain.User{
			ID:       userExtended.ID,
			Username: userExtended.Username,
		}
	} else {
		if err := a.repos.UserRepo.Authentication(ctx, input.UserID, input.Password); err != nil {
			logger.Error(fmt.Errorf("authentication failed for user %s: %w", input.UserID, err))
			return Tokens{}, nil, fmt.Errorf("authentication failed: %w", err)
		}

		user, err = a.repos.UserRepo.GetByID(ctx, input.UserID, input.Password)
		if err != nil {
			logger.Error(fmt.Errorf("find user data failed for user %s: %w", input.UserID, err))
			return Tokens{}, nil, fmt.Errorf("find user data failed: %w", err)
		}

		var academicGroup, profile string
		if user.Role != "teacher" && user.Role != "admin" {
			academicGroup, profile, err = a.repos.UserRepo.GetUserGroups(ctx, input.UserID, input.Password)
			if err != nil {
				logger.Warn(fmt.Sprintf("failed to get groups for user %s: %v", input.UserID, err))
			}
		}

		userExtended = &domain.UserExtended{
			ID:            user.ID,
			Username:      user.Username,
			Role:          user.Role,
			AcademicGroup: academicGroup,
			Profile:       profile,
		}
	}

	tokens, err := a.generateTokens(user)
	if err != nil {
		logger.Error(fmt.Errorf("failed to generate tokens for user %s: %w", input.UserID, err))
		return Tokens{}, nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	jti, err := a.tokenManager.ExtractClaim(tokens.RefreshToken, "jti")
	if err != nil {
		logger.Error(fmt.Errorf("failed to extract jti from token for user %s: %w", input.UserID, err))
		return Tokens{}, nil, fmt.Errorf("failed to extract jti: %w", err)
	}

	session := domain.RefreshSession{
		JTI:           jti,
		UserID:        input.UserID,
		Username:      user.Username,
		Role:          userExtended.Role,
		AcademicGroup: userExtended.AcademicGroup,
		Profile:       userExtended.Profile,
		ExpiresAt:     time.Now().Add(a.refreshTokenTTL),
		CreatedAt:     time.Now(),
	}

	if err := a.repos.SessionRepo.SaveRefreshToken(ctx, &session); err != nil {
		logger.Error(fmt.Errorf("failed to save refresh session for user %s: %w", input.UserID, err))
		return Tokens{}, nil, fmt.Errorf("failed to save refresh session: %w", err)
	}

	return tokens, userExtended, nil
}

func (a *AppUserService) SignOut(ctx context.Context, refreshToken string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if refreshToken == "" {
		logger.Error(fmt.Errorf("empty refresh token"))
		return fmt.Errorf("empty refresh token")
	}

	jti, err := a.tokenManager.ExtractClaim(refreshToken, "jti")
	if err != nil {
		logger.Error(fmt.Errorf("failed to extract jti from token: %w", err))
		return err
	}

	if err := a.repos.SessionRepo.RevokeRefreshToken(ctx, jti); err != nil {
		logger.Error(fmt.Errorf("failed to revoke refresh token: %w", err))
		return err
	}

	return nil
}

func (a *AppUserService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	if refreshToken == "" {
		logger.Error(fmt.Errorf("empty refresh token"))
		return "", fmt.Errorf("empty refresh token")
	}

	if err := a.tokenManager.ValidateRefreshToken(refreshToken); err != nil {
		logger.Warn(fmt.Sprintf("invalid refresh token structure: %v", err))
		return "", fmt.Errorf("invalid token")
	}

	userID, err := a.tokenManager.ExtractClaim(refreshToken, "user_id")
	if err != nil {
		logger.Error(fmt.Errorf("failed to extract user_id from token: %w", err))
		return "", fmt.Errorf("invalid token claims")
	}

	oldJti, err := a.tokenManager.ExtractClaim(refreshToken, "jti")
	if err != nil {
		logger.Error(fmt.Errorf("failed to extract jti from token: %w", err))
		return "", fmt.Errorf("invalid token claims")
	}

	tokenExists, err := a.repos.SessionRepo.TokenExists(ctx, oldJti)
	if err != nil {
		logger.Error(fmt.Errorf("failed to check token existence: %w", err))
		return "", fmt.Errorf("authentication service unavailable")
	}

	if !tokenExists {
		logger.Warn(fmt.Sprintf("attempt to use non-existent refresh token: jti=%s, user=%s",
			oldJti, userID))
		return "", fmt.Errorf("token not found or already used")
	}

	userExtended, err := a.repos.SessionRepo.GetExtendedUserByID(ctx, userID)
	if err != nil {
		logger.Error(fmt.Errorf("failed to get extended user data from session for ID %s: %w", userID, err))
		return "", fmt.Errorf("failed to get user data: %w", err)
	}

	newRefreshToken, err := a.tokenManager.NewRefreshToken(userID)
	if err != nil {
		logger.Error(fmt.Errorf("failed to generate refresh token for user %s: %w", userID, err))
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	newJti, err := a.tokenManager.ExtractClaim(newRefreshToken, "jti")
	if err != nil {
		logger.Error(fmt.Errorf("failed to extract jti from new token for user %s: %w", userID, err))
		return "", fmt.Errorf("failed to extract new token JTI: %w", err)
	}

	newSession := domain.RefreshSession{
		JTI:           newJti,
		UserID:        userID,
		Username:      userExtended.Username,
		Role:          userExtended.Role,
		AcademicGroup: userExtended.AcademicGroup,
		Profile:       userExtended.Profile,
		ExpiresAt:     time.Now().Add(a.refreshTokenTTL),
		CreatedAt:     time.Now(),
	}

	if err := a.repos.SessionRepo.ReplaceRefreshToken(ctx, oldJti, &newSession); err != nil {
		logger.Error(fmt.Errorf("failed to replace refresh token: %w", err))
		return "", fmt.Errorf("failed to rotate tokens")
	}

	return newRefreshToken, nil
}

func (a *AppUserService) ValidateAccessToken(ctx context.Context, accessToken string) (*domain.UserExtended, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if accessToken == "" {
		logger.Error(fmt.Errorf("empty access token"))
		return nil, fmt.Errorf("empty access token")
	}

	err := a.tokenManager.Validate(accessToken)
	if err != nil {
		logger.Error(fmt.Errorf("token validation failed: %w", err))
		return nil, fmt.Errorf("invalid token")
	}

	userID, err := a.tokenManager.ExtractClaim(accessToken, "user_id")
	if err != nil {
		logger.Error(fmt.Errorf("failed to extract user_id from token: %w", err))
		return nil, err
	}

	userExtended, err := a.repos.SessionRepo.GetExtendedUserByID(ctx, userID)
	if err != nil {
		logger.Error(fmt.Errorf("failed to get extended user data from session for ID %s: %w", userID, err))
		return nil, fmt.Errorf("failed to get user data: %w", err)
	}

	return userExtended, nil
}

func (a *AppUserService) GetAccessToken(ctx context.Context, refreshToken string) (string, *domain.UserExtended, error) {
	if ctx.Err() != nil {
		return "", nil, ctx.Err()
	}

	if refreshToken == "" {
		logger.Error(fmt.Errorf("empty refresh token"))
		return "", nil, fmt.Errorf("empty refresh token")
	}

	if err := a.tokenManager.ValidateRefreshToken(refreshToken); err != nil {
		logger.Warn(fmt.Sprintf("invalid refresh token structure: %v", err))
		return "", nil, fmt.Errorf("invalid token")
	}

	userID, err := a.tokenManager.ExtractClaim(refreshToken, "user_id")
	if err != nil {
		logger.Error(fmt.Errorf("failed to extract user_id from token: %w", err))
		return "", nil, fmt.Errorf("invalid token claims")
	}

	jti, err := a.tokenManager.ExtractClaim(refreshToken, "jti")
	if err != nil {
		logger.Error(fmt.Errorf("failed to extract jti from token: %w", err))
		return "", nil, fmt.Errorf("invalid token claims")
	}

	tokenExists, err := a.repos.SessionRepo.TokenExists(ctx, jti)
	if err != nil {
		logger.Error(fmt.Errorf("failed to check token existence: %w", err))
		return "", nil, fmt.Errorf("authentication service unavailable")
	}

	if !tokenExists {
		logger.Warn(fmt.Sprintf("attempt to use non-existent refresh token: jti=%s, user=%s", jti, userID))
		return "", nil, fmt.Errorf("token not found or already used")
	}

	userExtended, err := a.repos.SessionRepo.GetExtendedUserByID(ctx, userID)
	if err != nil {
		logger.Error(fmt.Errorf("failed to get extended user data from session for ID %s: %w", userID, err))
		return "", nil, fmt.Errorf("failed to get user data: %w", err)
	}

	user := &domain.User{
		ID:       userExtended.ID,
		Username: userExtended.Username,
		Role:     userExtended.Role,
	}

	accessToken, err := a.tokenManager.NewAccessToken(user.ID, user.Username, user.Role)
	if err != nil {
		logger.Error(fmt.Errorf("failed to generate access token for user %s: %w", user.ID, err))
		return "", nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	return accessToken, userExtended, nil
}

func (a *AppUserService) generateTokens(user *domain.User) (Tokens, error) {
	newAccess, err := a.tokenManager.NewAccessToken(user.ID, user.Username, user.Role)
	if err != nil {
		logger.Error(fmt.Errorf("failed to generate access token for user %s: %w", user.ID, err))
		return Tokens{}, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefresh, err := a.tokenManager.NewRefreshToken(user.ID)
	if err != nil {
		logger.Error(fmt.Errorf("failed to generate refresh token for user %s: %w", user.ID, err))
		return Tokens{}, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return Tokens{
		AccessToken:  newAccess,
		RefreshToken: newRefresh,
	}, nil
}
