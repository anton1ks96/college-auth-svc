package app

import (
	"fmt"

	"github.com/anton1ks96/college-auth-svc/internal/config"
	"github.com/anton1ks96/college-auth-svc/internal/handlers"
	"github.com/anton1ks96/college-auth-svc/internal/repository"
	"github.com/anton1ks96/college-auth-svc/internal/server"
	service "github.com/anton1ks96/college-auth-svc/internal/services"
	"github.com/anton1ks96/college-auth-svc/pkg/auth"
	"github.com/anton1ks96/college-auth-svc/pkg/database/mongodb"
	"github.com/anton1ks96/college-auth-svc/pkg/logger"
)

func Run() {
	cfg, err := config.Init()
	if err != nil {
		logger.Fatal(err)
	}

	db, err := mongodb.NewClient(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	userRepo := repository.NewUserRepository(cfg)
	sessRepo := repository.NewSessionsRepository(cfg, db)

	tokenManager := auth.NewManager(cfg)

	services := service.NewServices(service.Deps{
		Repos: &service.Repositories{
			UserRepo:    userRepo,
			SessionRepo: sessRepo,
		},
		TokenManager: tokenManager,
		Config:       cfg,
	})

	handler := handlers.NewHandler(services, *tokenManager, cfg)

	router := handler.Init()

	srv := server.NewServer(cfg, router)

	logger.Info(fmt.Sprintf("College Auth Service started - PORT: %s", cfg.Server.Port))

	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
