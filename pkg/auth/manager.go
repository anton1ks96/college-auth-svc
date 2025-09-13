package auth

import (
	"errors"
	"fmt"
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

func (m *Manager) NewAccessToken(userId, userName, role string) (string, error) {
	if userId == "" || userName == "" || role == "" {
		return "", errors.New("userId, userName and role cannot be empty")
	}

	ttl, err := time.ParseDuration(m.cfg.JWT.AccessTokenTTL)
	if err != nil {
		logger.Error(errors.New("failed to parse access token TTL: " + err.Error()))
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  userId,
		"username": userName,
		"role":     role,
		"exp":      time.Now().Add(ttl).Unix(),
		"iat":      time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(m.cfg.JWT.SigningKey))
	if err != nil {
		logger.Error(errors.New("failed to sign token: " + err.Error()))
		return "", err
	}

	return tokenString, nil
}

func (m *Manager) NewRefreshToken(userId string) (string, error) {
	if userId == "" {
		return "", errors.New("userId cannot be empty")
	}

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

	return tokenString, nil
}

func (m *Manager) ExtractClaim(tokenString string, claim string) (string, error) {
	if tokenString == "" {
		return "", errors.New("token cannot be empty")
	}
	if claim == "" {
		return "", errors.New("claim name cannot be empty")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.cfg.JWT.SigningKey), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		logger.Error(errors.New("failed to parse token: " + err.Error()))
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid claims format")
	}

	claimValue, exists := claims[claim]
	if !exists {
		return "", errors.New("claim not found in token")
	}

	claimString, ok := claimValue.(string)
	if !ok {
		return "", errors.New("claim is not a string")
	}

	return claimString, nil
}

func (m *Manager) Validate(tokenString string) error {
	if tokenString == "" {
		return errors.New("token cannot be empty")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.cfg.JWT.SigningKey), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		logger.Error(errors.New("failed to parse token: " + err.Error()))
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid claims format")
	}

	expAt, ok := claims["exp"].(float64)
	if !ok {
		return errors.New("expiration claim missing or invalid")
	}

	if time.Now().Unix() > int64(expAt) {
		logger.Error(errors.New("token expired"))
		return errors.New("token expired")
	}

	return nil
}

func (m *Manager) ValidateRefreshToken(tokenString string) error {
	if err := m.Validate(tokenString); err != nil {
		return err
	}

	_, err := m.ExtractClaim(tokenString, "jti")
	if err != nil {
		return errors.New("refresh token must contain JTI")
	}

	return nil
}
