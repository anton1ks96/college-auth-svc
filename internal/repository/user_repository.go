package repository

import (
	"context"
	"fmt"
	"strings"

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

func (u *UserRepository) Authentication(ctx context.Context, userID, userPass string) error {
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

	var dn string

	if !strings.HasPrefix(userID, "t") {
		dn = fmt.Sprintf("uid=%s,ou=people,dc=it-college,dc=ru", userID)
	} else {
		dn = fmt.Sprintf("uid=%s,ou=teachers,dc=it-college,dc=ru", userID)
	}

	if err := l.Bind(dn, userPass); err != nil {
		logger.Warn(fmt.Sprintf("LDAP authentication failed for user %s with DN %s", userID, dn))
		return fmt.Errorf("authentication failed")
	}

	return nil
}

func (u *UserRepository) GetByID(ctx context.Context, userID, userPass string) (*domain.User, error) {
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

	var dn string

	if !strings.HasPrefix(userID, "t") {
		dn = fmt.Sprintf("uid=%s,ou=people,dc=it-college,dc=ru", userID)
	} else {
		dn = fmt.Sprintf("uid=%s,ou=Teachers,dc=it-college,dc=ru", userID)
	}

	if err := l.Bind(dn, userPass); err != nil {
		logger.Error(fmt.Errorf("failed to bind account %s to LDAP during user lookup for %s: %w", userID, userID, err))
		return nil, fmt.Errorf("service account bind failed")
	}

	searchFilter := fmt.Sprintf("(uid=%s)", ldap.EscapeFilter(userID))

	userReq := ldap.NewSearchRequest(
		dn,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		5,
		false,
		searchFilter,
		[]string{"uid", "cn", "employeeType"},
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

	entry := sr.Entries[0]
	user := &domain.User{
		ID:       entry.GetAttributeValue("uid"),
		Username: entry.GetAttributeValue("cn"),
		Role:     entry.GetAttributeValue("employeeType"),
	}

	return user, nil
}

func (u *UserRepository) GetUserGroups(ctx context.Context, userID, userPass string) (academicGroup, profile string, err error) {
	if ctx.Err() != nil {
		logger.Error(fmt.Errorf("context cancelled during group retrieval for user %s: %w", userID, ctx.Err()))
		return "", "", ctx.Err()
	}

	l, err := ldap.DialURL(u.cfg.LDAP.URL)
	if err != nil {
		logger.Error(fmt.Errorf("failed to connect to LDAP server %s for group lookup: %w", u.cfg.LDAP.URL, err))
		return "", "", fmt.Errorf("LDAP connection failed")
	}
	defer l.Close()

	var userDN string
	if !strings.HasPrefix(userID, "t") {
		userDN = fmt.Sprintf("uid=%s,ou=people,dc=it-college,dc=ru", userID)
	} else {
		userDN = fmt.Sprintf("uid=%s,ou=teachers,dc=it-college,dc=ru", userID)
	}

	if err := l.Bind(userDN, userPass); err != nil {
		logger.Error(fmt.Errorf("failed to bind for group lookup: %w", err))
		return "", "", fmt.Errorf("authentication failed")
	}

	searchFilter := fmt.Sprintf(
		"(&(|(objectClass=groupOfNames)(objectClass=posixGroup)(objectClass=group))"+
			"(|(member=%s)(memberUid=%s)))",
		ldap.EscapeFilter(userDN),
		ldap.EscapeFilter(userID),
	)

	searchRequest := ldap.NewSearchRequest(
		"ou=groups,dc=it-college,dc=ru",
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		searchFilter,
		[]string{"cn", "description"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		logger.Error(fmt.Errorf("LDAP group search failed for user %s: %w", userID, err))
		return "", "", fmt.Errorf("group search failed")
	}

	for _, entry := range sr.Entries {
		cn := entry.GetAttributeValue("cn")
		description := entry.GetAttributeValue("description")

		if description == "Академическая группа" && strings.HasPrefix(cn, "ИТ") {
			academicGroup = cn
			logger.Debug(fmt.Sprintf("found academic group for user %s: %s", userID, cn))
		}

		if description == "Профиль" {
			validProfiles := map[string]bool{
				"BE": true, "FE": true, "PM": true,
				"CD": true, "GD": true, "SA": true,
			}
			if validProfiles[cn] {
				profile = cn
				logger.Debug(fmt.Sprintf("found profile for user %s: %s", userID, cn))
			}
		}
	}

	if !strings.HasPrefix(userID, "t") && academicGroup == "" {
		logger.Warn(fmt.Sprintf("no academic group found for student %s", userID))
	}

	return academicGroup, profile, nil
}
