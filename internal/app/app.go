package app

import (
	"github.com/anton1ks96/college-auth-svc/internal/config"
	"github.com/anton1ks96/college-auth-svc/internal/handlers"
	"github.com/anton1ks96/college-auth-svc/pkg/auth"
	"github.com/anton1ks96/college-auth-svc/pkg/logger"
)

func Run() {
	// Initialize configuration
	cfg, err := config.Init()
	if err != nil {
		logger.Error(err)
		return
	}

	//userRepo := ldap.NewUserRepository(cfg)
	//ctx := context.Background()
	//if err := userRepo.Authentication(ctx, cfg.LDAP.BindUsername, cfg.LDAP.BindPassword); err != nil {
	//	logger.Error(err)
	//	return
	//}

	// Initialize JWT token manager
	//tokenManager := auth.NewManager(cfg)

	//Initialize MongoDB
	//db, err := mongodb.NewClient(cfg)
	//if err != nil {
	//	logger.Error(err)
	//	return
	//}
	//
	//session := domain.RefreshSession{
	//	JTI:       "test",
	//	Username:  "test",
	//	ExpiresAt: time.Now(),
	//	CreatedAt: time.Now(),
	//}
	//
	//sessionRepo := mg.NewSessionsRepository(cfg, db)
	//if err := sessionRepo.SaveRefreshToken(context.TODO(), session); err != nil {
	//	logger.Error(err)
	//	return
	//}
	// Router and server initialization
	tokenManager := auth.NewManager(cfg)
	 := handlers.NewHandler(nil, *tokenManager)
}
