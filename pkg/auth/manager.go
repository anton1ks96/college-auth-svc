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

func (m *Manager) NewAccessToken(userName string) (string, error) {
	ttl, err := time.ParseDuration(m.cfg.JWT.AccessTokenTTL)
	if err != nil {
		logger.Error(errors.New("failed to parse access token TTL: " + err.Error()))
		return "", err
	}

	var role string

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
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

	return tokenString, err
}

func (m *Manager) NewRefreshToken(userName string) (string, error) {
	jti := uuid.New().String()

	ttl, err := time.ParseDuration(m.cfg.JWT.RefreshTokenTTL)
	if err != nil {
		logger.Error(errors.New("failed to parse refresh token TTL: " + err.Error()))
		return "", err
	}

	claims := jwt.MapClaims{
		"username": userName,
		"jti":      jti,
		"exp":      time.Now().Add(ttl).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(m.cfg.JWT.SigningKey))
	if err != nil {
		logger.Error(errors.New("failed to sign token: " + err.Error()))
		return "", err
	}

	return tokenString, err
}

func (m *Manager) ExtractClaim(tokenString string, claim string) (string, error) {
	var token, err = jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(m.cfg.JWT.SigningKey), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		logger.Error(errors.New("failed to parse token: " + err.Error()))
		return "", err
	}

	claims := token.Claims.(jwt.MapClaims)
	extracted := claims[claim].(string)

	return extracted, err
}

func (m *Manager) Validate(tokenString string) (string, error) {
	var token, err = jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(m.cfg.JWT.SigningKey), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		logger.Error(errors.New("failed to parse token: " + err.Error()))
		return "", err
	}

	claims := token.Claims.(jwt.MapClaims)

	role := claims["role"].(string)
	expAt := claims["exp"].(float64)
	if time.Now().Unix() > int64(expAt) {
		logger.Error(errors.New("token expired"))
		return "", errors.New("token expired")
	}

	return role, nil
}
