package app

import (
	"context"

	"github.com/anton1ks96/college-auth-svc/internal/config"
	"github.com/anton1ks96/college-auth-svc/pkg/logger"
)

func Run() {
	// Initialize configuration
	cfg, err := config.Init()
	if err != nil {
		logger.Error(err)
		return
	}

	ctx := context.Background()
	//
}
