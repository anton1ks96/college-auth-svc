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

	userDN, err := u.findUserDN(l, userID)
	if err != nil {
		logger.Error(fmt.Errorf("failed to find DN for user %s: %w", userID, err))
		return fmt.Errorf("user not found")
	}

	logger.Debug(fmt.Sprintf("Found DN for user %s: %s", userID, userDN))

	if err := l.Bind(userDN, userPass); err != nil {
		logger.Warn(fmt.Sprintf("LDAP authentication failed for user %s with DN %s", userID, userDN))
		return fmt.Errorf("authentication failed: %s", err.Error())
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

	dn, err := u.findUserDN(l, userID)
	if err != nil {
		logger.Error(fmt.Errorf("failed to find DN for user %s: %w", userID, err))
		return nil, fmt.Errorf("user not found")
	}

	logger.Debug(fmt.Sprintf("Found DN for user %s: %s", userID, dn))

	var baseDN string
	if !strings.HasPrefix(userID, "t") {
		baseDN = "ou=people,dc=it-college,dc=ru"
	} else {
		baseDN = "ou=people,ou=Teachers,dc=it-college,dc=ru"
	}

	if err := l.Bind(dn, userPass); err != nil {
		logger.Error(fmt.Errorf("failed to bind account %s to LDAP during user lookup for %s: %w", userID, userID, err))
		return nil, fmt.Errorf("service account bind failed")
	}

	searchFilter := fmt.Sprintf("(uid=%s)", ldap.EscapeFilter(userID))

	logger.Debug(fmt.Sprintf("Search filter: %s in baseDN: %s", searchFilter, baseDN))

	userReq := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		5,
		false,
		searchFilter,
		[]string{"uid", "cn", "memberOf"},
		nil,
	)

	sr, err := l.Search(userReq)
	if err != nil {
		logger.Error(fmt.Errorf("LDAP search failed for user %s with filter %s in baseDN %s: %w", userID, searchFilter, baseDN, err))
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

	uid := entry.GetAttributeValue("uid")
	cn := entry.GetAttributeValue("cn")
	memberOfValues := entry.GetAttributeValues("memberOf")

	logger.Debug(fmt.Sprintf("User %s memberOf: %v", userID, memberOfValues))

	role := u.determineRole(memberOfValues, dn)

	logger.Debug(fmt.Sprintf("User %s role determined as: %s", userID, role))

	user := &domain.User{
		ID:       uid,
		Username: cn,
		Role:     role,
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
		userDN = fmt.Sprintf("uid=%s,ou=people,ou=Teachers,dc=it-college,dc=ru", userID)
	}

	if err := l.Bind(userDN, userPass); err != nil {
		logger.Error(fmt.Errorf("failed to bind for group lookup with DN %s: %w", userDN, err))
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

func (u *UserRepository) findUserDN(l *ldap.Conn, userID string) (string, error) {
	var baseDN string
	if !strings.HasPrefix(userID, "t") {
		baseDN = "ou=people,dc=it-college,dc=ru"
	} else {
		baseDN = "ou=people,ou=Teachers,dc=it-college,dc=ru"
	}

	err := l.UnauthenticatedBind("")
	if err != nil {
		logger.Debug(fmt.Sprintf("Anonymous bind failed, trying without bind: %v", err))
	}

	searchFilter := fmt.Sprintf("(uid=%s)", ldap.EscapeFilter(userID))
	searchRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		1,
		5,
		false,
		searchFilter,
		[]string{"dn"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		logger.Error(fmt.Errorf("LDAP search failed for user %s in baseDN %s: %w", userID, baseDN, err))
		return "", fmt.Errorf("search failed: %w", err)
	}

	if len(sr.Entries) == 0 {
		logger.Warn(fmt.Sprintf("user %s not found in LDAP", userID))
		return "", fmt.Errorf("user not found")
	}

	if len(sr.Entries) > 1 {
		logger.Warn(fmt.Sprintf("multiple entries found for user %s", userID))
		return "", fmt.Errorf("multiple users found")
	}

	return sr.Entries[0].DN, nil
}

func (u *UserRepository) determineRole(memberOfValues []string, userDN string) string {
	userDNLower := strings.ToLower(userDN)

	isTeacherOU := strings.Contains(userDNLower, "ou=teachers")
	isPeopleOU := strings.Contains(userDNLower, "ou=people,dc=it-college,dc=ru") && !isTeacherOU

	logger.Debug(fmt.Sprintf("User DN: %s, isTeacherOU: %v, isPeopleOU: %v", userDN, isTeacherOU, isPeopleOU))

	for _, memberOf := range memberOfValues {
		if !strings.Contains(memberOf, "ou=groups,dc=it-college,dc=ru") {
			continue
		}

		parts := strings.Split(memberOf, ",")
		if len(parts) == 0 {
			continue
		}

		cnPart := parts[0]
		if !strings.HasPrefix(cnPart, "cn=") {
			continue
		}

		cn := strings.TrimPrefix(cnPart, "cn=")

		logger.Debug(fmt.Sprintf("Checking group cn=%s from memberOf", cn))

		if cn == "admin" {
			logger.Debug("User is member of admin group, role: admin")
			return "admin"
		}

		if cn == "teachers" {
			logger.Debug("User is member of teachers group, role: teacher")
			return "teacher"
		}

		if isPeopleOU && strings.HasPrefix(cn, "ИТ") {
			logger.Debug(fmt.Sprintf("User is in ou=people and member of academic group %s, role: student", cn))
			return "student"
		}
	}

	if isTeacherOU {
		logger.Debug("User is in ou=Teachers with no group membership, assigning teacher role by OU location")
		return "teacher"
	}

	logger.Warn("Role not determined from groups and DN")
	return ""
}
