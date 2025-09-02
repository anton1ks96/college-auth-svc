package ldap

import (
	"context"
	"fmt"

	"github.com/anton1ks96/college-auth-svc/internal/config"
	"github.com/anton1ks96/college-auth-svc/internal/domain"
	"github.com/anton1ks96/college-auth-svc/pkg/logger"
	"github.com/go-ldap/ldap/v3"
)

type Client struct {
	cfg *config.Config
}

func NewLDAPClient(cfg *config.Config) *Client {
	return &Client{
		cfg: cfg,
	}
}

func (c *Client) Authentication(ctx context.Context, username, password string) (*domain.User, error) {
	l, err := ldap.DialURL(c.cfg.LDAP.URL)
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

	logger.Info("Connection successfully created")
	return nil, nil
}
