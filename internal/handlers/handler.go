package handlers

import (
	"github.com/anton1ks96/college-auth-svc/internal/config"
	service "github.com/anton1ks96/college-auth-svc/internal/services"
	"github.com/anton1ks96/college-auth-svc/pkg/auth"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services     service.Services
	tokenManager auth.Manager
}

func NewHandler(services *service.Services, manager auth.Manager) *Handler {
	return &Handler{
		services:     *services,
		tokenManager: manager,
	}
}

func (h *Handler) Init(cfg *config.Config) *gin.Engine {
	router := gin.Default()
	router.Use(corsMiddleware)
	return router
}
