package handlers

import (
	"net/http"

	"github.com/anton1ks96/college-auth-svc/internal/config"
	v1 "github.com/anton1ks96/college-auth-svc/internal/handlers/v1"
	service "github.com/anton1ks96/college-auth-svc/internal/services"
	"github.com/anton1ks96/college-auth-svc/pkg/auth"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services     *service.Services
	tokenManager auth.Manager
	cfg          *config.Config
}

func NewHandler(services *service.Services, manager auth.Manager, cfg *config.Config) *Handler {
	return &Handler{
		services:     services,
		tokenManager: manager,
		cfg:          cfg,
	}
}

func (h *Handler) Init() *gin.Engine {
	router := gin.New()

	router.Use(
		gin.Recovery(),
		gin.Logger(),
		corsMiddleware,
	)

	h.initAPI(router)

	return router
}

func (h *Handler) ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
		"version": "v1",
	})
}

func (h *Handler) initAPI(router *gin.Engine) {
	handlerV1 := v1.NewHandler(h.services, h.tokenManager, h.cfg)
	api := router.Group("/api")
	router.GET("/ping", h.ping)
	handlerV1.Init(api)
}
