package auth

import (
	"errors"
	"time"

	"github.com/anton1ks96/college-auth-svc/internal/config"
	"github.com/anton1ks96/college-auth-svc/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Manager struct {
	cfg *config.Config
}

func NewManager(cfg *config.Config) *Manager {
	return &Manager{
		cfg: cfg,
	}
}

func (m *Manager) NewAccessToken(userId string) (string, error) {
	ttl, err := time.ParseDuration(m.cfg.JWT.AccessTokenTTL)
	if err != nil {
		logger.Error(errors.New("failed to parse access token TTL: " + err.Error()))
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(ttl).Unix(),
		"iat":     time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(m.cfg.JWT.SigningKey))
	if err != nil {
		logger.Error(errors.New("failed to sign token: " + err.Error()))
		return "", err
	}

	return tokenString, err
}

func (m *Manager) NewRefreshToken(userId string) (string, error) {
	jti := uuid.New().String()

	ttl, err := time.ParseDuration(m.cfg.JWT.RefreshTokenTTL)
	if err != nil {
		logger.Error(errors.New("failed to parse refresh token TTL: " + err.Error()))
		return "", err
	}

	claims := jwt.MapClaims{
		"user_id": userId,
		"jti":     jti,
		"exp":     time.Now().Add(ttl).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(m.cfg.JWT.SigningKey))
	if err != nil {
		logger.Error(errors.New("failed to sign token: " + err.Error()))
		return "", err
	}

	return tokenString, err
}
