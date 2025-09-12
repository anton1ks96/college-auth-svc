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

func (u *UserRepository) Authentication(ctx context.Context, userID string, password string) error {
	if ctx.Err() != nil {
		logger.Error(fmt.Errorf("context cancelled during authentication for user %s: %w", userID, ctx.Err()))
		return ctx.Err()
	}

	l, err := ldap.DialURL(u.cfg.LDAP.URL)
	if err != nil {
		logger.Error(fmt.Errorf("failed to connect to LDAP server %s for user %s: %w", u.cfg.LDAP.URL, userID, err))
		return fmt.Errorf("LDAP connection failed")
	}
	defer l.Close()

	logger.Debug(fmt.Sprintf("LDAP connection established for user: %s", userID))

	if ctx.Err() != nil {
		logger.Error(fmt.Errorf("context cancelled after LDAP connection for user %s: %w", userID, ctx.Err()))
		return ctx.Err()
	}

	dn := fmt.Sprintf("uid=%s,ou=people,dc=it-college,dc=ru", userID)

	if err := l.Bind(dn, password); err != nil {
		logger.Warn(fmt.Sprintf("LDAP authentication failed for user %s with DN %s", userID, dn))
		return fmt.Errorf("authentication failed")
	}

	if ctx.Err() != nil {
		logger.Error(fmt.Errorf("context cancelled after successful bind for user %s: %w", userID, ctx.Err()))
		return ctx.Err()
	}

	return nil
}

func (u *UserRepository) GetByID(ctx context.Context, userID string) (*domain.User, error) {
	if ctx.Err() != nil {
		logger.Error(fmt.Errorf("context cancelled during user retrieval for user %s: %w", userID, ctx.Err()))
		return nil, ctx.Err()
	}

	l, err := ldap.DialURL(u.cfg.LDAP.URL)
	if err != nil {
		logger.Error(fmt.Errorf("failed to connect to LDAP server %s during user lookup for %s: %w", u.cfg.LDAP.URL, userID, err))
		return nil, fmt.Errorf("LDAP connection failed")
	}
	defer l.Close()

	if ctx.Err() != nil {
		logger.Error(fmt.Errorf("context cancelled after LDAP connection during user lookup for %s: %w", userID, ctx.Err()))
		return nil, ctx.Err()
	}

	serviceDN := fmt.Sprintf("uid=%s,ou=people,dc=it-college,dc=ru", u.cfg.LDAP.BindUsername)

	if err := l.Bind(serviceDN, u.cfg.LDAP.BindPassword); err != nil {
		logger.Error(fmt.Errorf("failed to bind service account %s to LDAP during user lookup for %s: %w", u.cfg.LDAP.BindUsername, userID, err))
		return nil, fmt.Errorf("service account bind failed")
	}

	if ctx.Err() != nil {
		logger.Error(fmt.Errorf("context cancelled after service account bind during user lookup for %s: %w", userID, ctx.Err()))
		return nil, ctx.Err()
	}

	searchFilter := fmt.Sprintf("(uid=%s)", ldap.EscapeFilter(userID))

	userReq := ldap.NewSearchRequest(
		"ou=people,dc=it-college,dc=ru",
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		5,
		false,
		searchFilter,
		[]string{"uid", "cn", "mail", "employeeType"},
		nil,
	)

	sr, err := l.Search(userReq)
	if err != nil {
		logger.Error(fmt.Errorf("LDAP search failed for user %s with filter %s: %w", userID, searchFilter, err))
		return nil, fmt.Errorf("user search failed")
	}

	if ctx.Err() != nil {
		logger.Error(fmt.Errorf("context cancelled after LDAP search for user %s: %w", userID, ctx.Err()))
		return nil, ctx.Err()
	}

	if len(sr.Entries) == 0 {
		logger.Warn(fmt.Sprintf("user %s not found in LDAP directory", userID))
		return nil, fmt.Errorf("user not found")
	}

	if len(sr.Entries) > 1 {
		logger.Error(fmt.Errorf("multiple LDAP entries (%d) found for user %s - this should not happen", len(sr.Entries), userID))
		return nil, fmt.Errorf("multiple users found")
	}

	// TODO: determine user role func

	entry := sr.Entries[0]
	user := &domain.User{
		ID:       entry.GetAttributeValue("uid"),
		Username: entry.GetAttributeValue("cn"),
		Role:     entry.GetAttributeValue("employeeType"),
	}

	return user, nil
}
