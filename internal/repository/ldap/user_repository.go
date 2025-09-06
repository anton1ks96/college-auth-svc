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

// Authentication creates a new LDAP connection for each authentication request
// Also we can use only one LDAP service connection in the future

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
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	l, err := ldap.DialURL(u.cfg.LDAP.URL)
	if err != nil {
		logger.Error(err)
		return nil, fmt.Errorf("failed to create LDAP connection: %w", err)
	}
	defer l.Close()

	if ctx.Err() != nil {
		return nil, fmt.Errorf("operation cancelled: %w", ctx.Err())
	}

	dn := fmt.Sprintf("uid=%s,ou=people,dc=it-college,dc=ru", u.cfg.LDAP.BindUsername)

	if err := l.Bind(dn, u.cfg.LDAP.BindPassword); err != nil {
		logger.Error(err)
		return nil, fmt.Errorf("failed to bind user to LDAP: %w", err)
	}

	if ctx.Err() != nil {
		return nil, fmt.Errorf("operation cancelled after bind: %w", ctx.Err())
	}

	//var user domain.User
	userReq := ldap.NewSearchRequest(
		dn,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		5,
		false,
		fmt.Sprintf("(uid=%s)", ldap.EscapeFilter(username)),
		[]string{"uid", "cn", "mail", "employeeType"},
		nil,
	)

	sr, err := l.Search(userReq)
	if err != nil {
		logger.Error(err)
		return nil, fmt.Errorf("failed to find user in LDAP: %w", err)
	}

	if ctx.Err() != nil {
		return nil, fmt.Errorf("operation cancelled after search: %w", ctx.Err())
	}

	for _, entry := range sr.Entries {
		fmt.Printf("DN: %s\n", entry.DN)
		for _, attr := range entry.Attributes {
			fmt.Printf("  %s: %v\n", attr.Name, attr.Values)
		}
		fmt.Println()
	}

	return nil, nil
}
