package repository

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
		logger.Error(fmt.Errorf("failed to connect to LDAP server %s: %w", u.cfg.LDAP.URL, err))
		return err
	}
	defer l.Close()

	if ctx.Err() != nil {
		return ctx.Err()
	}

	dn := fmt.Sprintf("uid=%s,ou=people,dc=it-college,dc=ru", username)

	if err := l.Bind(dn, password); err != nil {
		logger.Error(fmt.Errorf("LDAP bind failed for user %s: %w", username, err))
		return err
	}

	if ctx.Err() != nil {
		return ctx.Err()
	}

	logger.Debug(fmt.Sprintf("LDAP authentication successful for user %s", username))
	return nil
}

func (u *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	l, err := ldap.DialURL(u.cfg.LDAP.URL)
	if err != nil {
		logger.Error(fmt.Errorf("failed to connect to LDAP server %s: %w", u.cfg.LDAP.URL, err))
		return nil, err
	}
	defer l.Close()

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	dn := fmt.Sprintf("uid=%s,ou=people,dc=it-college,dc=ru", u.cfg.LDAP.BindUsername)

	if err := l.Bind(dn, u.cfg.LDAP.BindPassword); err != nil {
		logger.Error(fmt.Errorf("failed to bind service account to LDAP: %w", err))
		return nil, err
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	userReq := ldap.NewSearchRequest(
		"ou=people,dc=it-college,dc=ru",
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
		logger.Error(fmt.Errorf("failed to search for user %s in LDAP: %w", username, err))
		return nil, err
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if len(sr.Entries) == 0 {
		logger.Debug(fmt.Sprintf("user %s not found in LDAP", username))
		return nil, fmt.Errorf("user not found")
	}

	if len(sr.Entries) > 1 {
		logger.Error(fmt.Errorf("multiple entries found for user %s in LDAP", username))
		return nil, fmt.Errorf("multiple users found")
	}

	// TODO: determine user role func

	entry := sr.Entries[0]
	user := &domain.User{
		ID:       entry.GetAttributeValue("uid"),
		Username: entry.GetAttributeValue("cn"),
		Mail:     entry.GetAttributeValue("mail"),
	}

	logger.Debug(fmt.Sprintf("successfully retrieved user %s from LDAP", username))
	return user, nil
}
