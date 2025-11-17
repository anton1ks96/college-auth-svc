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
			user.POST("/signout", h.signOut)
			user.POST("/refresh", h.refreshTokens)
		}

		authenticated := v1.Group("/auth")
		{
			authenticated.POST("/validate", h.validateToken)
		}

		app := v1.Group("/app")
		{
			app.POST("/signin", h.appSignIn)
			app.POST("/signout", h.appSignOut)
			app.POST("/refresh", h.appRefreshToken)
			app.POST("/validate", h.appValidateToken)
			app.POST("/access", h.appGetAccess)
		}

		search := v1.Group("/search")
		{
			search.POST("/students", h.searchStudents)
			search.POST("/teachers", h.searchTeachers)
		}
	}
}

func (h *Handler) ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
		"version": "v1",
	})
}
