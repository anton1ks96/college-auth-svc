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

func (u *UserRepository) Authentication(ctx context.Context, username string, password string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	l, err := ldap.DialURL(u.cfg.LDAP.URL)
	if err != nil {
		logger.Error(err)
		return fmt.Errorf("failed to create LDAP connection: %w", err)
	}
	defer l.Close()

	if ctx.Err() != nil {
		return fmt.Errorf("operation cancelled: %w", ctx.Err())
	}

	dn := fmt.Sprintf("uid=%s,ou=people,dc=it-college,dc=ru", username)

	if err := l.Bind(dn, password); err != nil {
		logger.Error(err)
		return fmt.Errorf("failed to bind user to LDAP: %w", err)
	}

	if ctx.Err() != nil {
		return fmt.Errorf("operation cancelled after auth: %w", ctx.Err())
	}

	logger.Info("Connection successfully created")
	return nil
}

func (u *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	return nil, nil
}
