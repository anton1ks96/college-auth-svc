package app

import (
	"github.com/anton1ks96/college-auth-svc/internal/config"
	"github.com/anton1ks96/college-auth-svc/internal/handlers"
	"github.com/anton1ks96/college-auth-svc/pkg/database/mongodb"
	"github.com/anton1ks96/college-auth-svc/pkg/logger"
)

func Run() {
	// Initialize configuration
	cfg, err := config.Init()
	if err != nil {
		logger.Error(err)
		return
	}

	client, err := mongodb.NewClient(cfg.MongoDB.URI)
	if err != nil {
		logger.Error(err)
		return
	}

	_ = client.Database(cfg.MongoDB.DBName)

	// Router and server initialization
	router := handlers.NewRouter()
	err = router.Run(cfg.Server.Port)
	if err != nil {
		logger.Error(err)
		return
	}

}
