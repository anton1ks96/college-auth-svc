package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/anton1ks96/college-auth-svc/internal/domain"
	"github.com/anton1ks96/college-auth-svc/internal/repository"
	"github.com/anton1ks96/college-auth-svc/pkg/logger"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}

	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		logger.Error(err)
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	if user.Username == "" {
		logger.Error(fmt.Errorf("user has empty username field"))
		return nil, errors.New("user has empty username field")
	}

	if user.ID == "" {
		logger.Error(fmt.Errorf("user has empty ID field"))
		return nil, errors.New("user has empty ID field")
	}

	if user.Role == "" {
		logger.Error(fmt.Errorf("user has empty role field"))
		return nil, errors.New("user has empty role field")
	}

	logger.Info(fmt.Sprintf("Successfully get user: %s with role: %s", user.Username, user.Role))

	return user, nil
}
