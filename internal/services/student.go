package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/anton1ks96/college-auth-svc/internal/config"
	"github.com/anton1ks96/college-auth-svc/internal/domain"
	"github.com/anton1ks96/college-auth-svc/pkg/logger"
	"github.com/go-ldap/ldap/v3"
)

type StudentService interface {
	SearchStudents(ctx context.Context, query string) ([]domain.StudentInfo, error)
}

type StudentServiceImpl struct {
	cfg    *config.Config
	appCfg *config.App
}

func NewStudentService(cfg *config.Config, appCfg *config.App) *StudentServiceImpl {
	return &StudentServiceImpl{
		cfg:    cfg,
		appCfg: appCfg,
	}
}

func (s *StudentServiceImpl) SearchStudents(ctx context.Context, query string) ([]domain.StudentInfo, error) {
	if ctx.Err() != nil {
		return []domain.StudentInfo{}, nil
	}

	if query == "" {
		return []domain.StudentInfo{}, nil
	}

	if s.appCfg.Test {
		return []domain.StudentInfo{
			{ID: "i24s0291", Username: "Коломацкий Иван"},
			{ID: "i24s0002", Username: "Джапаридзе Артем"},
		}, nil
	}

	l, err := ldap.DialURL(s.cfg.LDAP.URL)
	if err != nil {
		logger.Error(fmt.Errorf("failed to connect to LDAP: %w", err))
		return nil, fmt.Errorf("LDAP connection failed")
	}
	defer l.Close()

	filter := fmt.Sprintf(
		"(&(objectClass=person)(!(uid=t*))(|(uid=*%s*)(cn=*%s*)))",
		ldap.EscapeFilter(query),
		ldap.EscapeFilter(query),
	)

	searchRequest := ldap.NewSearchRequest(
		"ou=people,dc=it-college,dc=ru",
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		50,
		0,
		false,
		filter,
		[]string{"uid", "cn"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		logger.Error(fmt.Errorf("LDAP search failed: %w", err))
		return nil, fmt.Errorf("search failed")
	}

	var students []domain.StudentInfo
	for _, entry := range sr.Entries {
		uid := entry.GetAttributeValue("uid")
		cn := entry.GetAttributeValue("cn")

		if !strings.HasPrefix(uid, "t") && uid != "" && cn != "" {
			students = append(students, domain.StudentInfo{
				ID:       uid,
				Username: cn,
			})
		}
	}

	return students, nil
}
