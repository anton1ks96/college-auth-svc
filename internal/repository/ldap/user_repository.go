package ldap

import (
	"context"
	"fmt"

	"github.com/anton1ks96/college-auth-svc/internal/config"
	"github.com/anton1ks96/college-auth-svc/internal/domain"
	"github.com/anton1ks96/college-auth-svc/pkg/logger"
	"github.com/go-ldap/ldap/v3"
)

type UserRepository struct {
	cfg *config.Config
}

func NewUserRepository(cfg *config.Config) *UserRepository {
	return &UserRepository{cfg: cfg}
}

func (u *UserRepository) Authentication(ctx context.Context, username string, password string) (*domain.User, error) {
	l, err := ldap.DialURL(u.cfg.LDAP.URL)
	if err != nil {
		logger.Error(err)
		return nil, fmt.Errorf("failed to create LDAP connection: %w", err)
	}
	defer l.Close()

	dn := fmt.Sprintf("uid=%s,ou=people,dc=it-college,dc=ru", username)

	if err := l.Bind(dn, password); err != nil {
		logger.Error(err)
		return nil, fmt.Errorf("failed to bind user to LDAP: %w", err)
	}

	user := domain.User{
		ID:       "123",
		Username: "asd",
		Role:     "asd",
	}

	logger.Info("Connection successfully created")
	return &user, nil
}

func (u *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UserRepository) GetUserRoles(ctx context.Context, username string) ([]string, error) {
	//TODO implement me
	panic("implement me")
}
