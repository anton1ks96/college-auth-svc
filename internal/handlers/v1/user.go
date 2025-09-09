package v1

import (
	"net/http"
	"time"

	"github.com/anton1ks96/college-auth-svc/internal/handlers/dto"
	service "github.com/anton1ks96/college-auth-svc/internal/services"
	"github.com/gin-gonic/gin"
)

//GET /api/v1/auth/me
//POST /api/v1/auth/signout
//POST /api/v1/auth/refresh
//GET /api/v1/ping

func (h *Handler) initUsersRoutes(api *gin.RouterGroup) {
	users := api.Group("/user")
	{
		users.POST("/signin", h.signIn)
		users.POST("/testsignin", h.signInTest)
	}
}

func (h *Handler) signIn(c *gin.Context) {
	var loginReq dto.LoginRequest
	if err := c.BindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	tokens, user, err := h.services.UserService.SignIn(c.Request.Context(), service.SignInInput{
		Username: loginReq.Username,
		Password: loginReq.Password,
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	accessTTL, err := time.ParseDuration(h.cfg.JWT.AccessTokenTTL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid token TTL",
		})
		return
	}

	c.SetCookie(
		"access_token",
		tokens.AccessToken,
		int(accessTTL.Seconds()),
		"/",
		"",
		true,
		true,
	)

	response := dto.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		User: dto.UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Role:     user.Role,
			Mail:     user.Mail,
		},
		ExpiresIn: int(accessTTL.Seconds()),
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) signInTest(c *gin.Context) {
	var loginReq dto.LoginRequest
	if err := c.BindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	tokens, user, err := h.services.UserService.SignInTest(c.Request.Context(), service.SignInInput{
		Username: loginReq.Username,
		Password: loginReq.Password,
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	accessTTL, err := time.ParseDuration(h.cfg.JWT.AccessTokenTTL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid token TTL",
		})
		return
	}

	c.SetCookie(
		"access_token",
		tokens.AccessToken,
		int(accessTTL.Seconds()),
		"/",
		"",
		true,
		true,
	)

	response := dto.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		User: dto.UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Role:     user.Role,
			Mail:     user.Mail,
		},
		ExpiresIn: int(accessTTL.Seconds()),
	}

	c.JSON(http.StatusOK, response)
}
