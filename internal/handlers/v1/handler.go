package v1

import (
	"net/http"

	"github.com/anton1ks96/college-auth-svc/internal/config"
	service "github.com/anton1ks96/college-auth-svc/internal/services"
	"github.com/anton1ks96/college-auth-svc/pkg/auth"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services     *service.Services
	tokenManager auth.Manager
	cfg          *config.Config
}

func NewHandler(services *service.Services, tokenManager auth.Manager, cfg *config.Config) *Handler {
	return &Handler{
		services:     services,
		tokenManager: tokenManager,
		cfg:          cfg,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	api.GET("/ping", h.ping)
	v1 := api.Group("/v1")
	{
		user := v1.Group("/users")
		{
			user.POST("/signin", h.signIn)
			user.POST("/testsignin", h.signInTest)
			user.POST("/signout", h.signOut)
		}

		authenticated := v1.Group("/auth")
		{
			authenticated.POST("/validate", h.validateToken)
		}
	}
}

func (h *Handler) ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
		"version": "v1",
	})
}
