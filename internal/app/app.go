package app

import (
	"github.com/anton1ks96/college-auth-svc/internal/config"
	"github.com/anton1ks96/college-auth-svc/internal/handlers"
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

	// Initialize MongoDB
	//client, err := mongodb.NewClient()
	//if err != nil {
	//	logger.Error(err)
	//	return
	//}
	//
	//db := client.Database(cfg.Mongo.DBName)

	// Router and server initialization
	router := handlers.NewRouter()
	err = router.Run(cfg.Server.Port)
	if err != nil {
		logger.Error(err)
		return
	}

}
