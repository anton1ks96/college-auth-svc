package repository

import (
	"context"
	"os"
	"testing"

	"github.com/anton1ks96/college-auth-svc/internal/config"
)

func TestGetUserGroups(t *testing.T) {
	ldapURL := os.Getenv("LDAP_URL")
	userID := os.Getenv("TEST_USER_ID")
	userPass := os.Getenv("TEST_USER_PASS")
	wantGroup := "ИТ24-11"
	wantProfile := "BE"
	wantSubgroup := ""
	wantEnglishGroup := "B1.21"

	if ldapURL == "" || userID == "" || userPass == "" {
		t.Skip("LDAP_URL, TEST_USER_ID и TEST_USER_PASS должны быть заданы через переменные окружения (export)")
	}

	cfg := &config.Config{
		LDAP: config.LDAPConfig{
			URL: ldapURL,
		},
	}

	repo := NewUserRepository(cfg)

	userGroups, err := repo.GetUserGroups(context.Background(), userID, userPass)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	t.Logf("got group=%q, profile=%q, subgroup=%q, english=%q",
		userGroups.AcademicGroup, userGroups.Profile, userGroups.Subgroup, userGroups.EnglishGroup)

	if userGroups.AcademicGroup != wantGroup {
		t.Errorf("expected group %q, got %q", wantGroup, userGroups.AcademicGroup)
	}

	if userGroups.Profile != wantProfile {
		t.Errorf("expected profile %q, got %q", wantProfile, userGroups.Profile)
	}

	if userGroups.Subgroup != wantSubgroup {
		t.Errorf("expected subgroup %q, got %q", wantSubgroup, userGroups.Subgroup)
	}

	if userGroups.EnglishGroup != wantEnglishGroup {
		t.Errorf("expected group %q, got %q", wantEnglishGroup, userGroups.EnglishGroup)
	}
}
