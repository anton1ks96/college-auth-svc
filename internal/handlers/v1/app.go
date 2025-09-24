package v1

import (
	"net/http"
	"strings"
	"time"

	"github.com/anton1ks96/college-auth-svc/internal/handlers/dto"
	service "github.com/anton1ks96/college-auth-svc/internal/services"
	"github.com/gin-gonic/gin"
)

func (h *Handler) appSignIn(c *gin.Context) {
	var loginReq dto.LoginRequest
	if err := c.BindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request format",
			"details": err.Error(),
		})
		return
	}

	tokens, user, err := h.services.AppUserService.SignIn(c.Request.Context(), service.SignInInput{
		UserID:   loginReq.Username,
		Password: loginReq.Password,
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "authentication failed",
			"details": err.Error(),
		})
		return
	}

	accessTTL, err := time.ParseDuration(h.cfg.JWT.AccessTokenTTL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "server configuration error",
		})
		return
	}

	response := dto.AppSignInResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    int(accessTTL.Seconds()),
		User: dto.AppUserInfo{
			ID:            user.ID,
			Username:      user.Username,
			Role:          user.Role,
			AcademicGroup: user.AcademicGroup,
			Profile:       user.Profile,
		},
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) appSignOut(c *gin.Context) {
	var req dto.AppSignOutRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request format",
			"details": err.Error(),
		})
		return
	}

	if err := h.services.AppUserService.SignOut(c.Request.Context(), req.RefreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to sign out",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "successfully signed out",
	})
}

func (h *Handler) appRefreshTokens(c *gin.Context) {
	var req dto.AppRefreshRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request format",
			"details": err.Error(),
		})
		return
	}

	tokens, err := h.services.AppUserService.RefreshTokens(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "failed to refresh tokens",
			"details": err.Error(),
		})
		return
	}

	accessTTL, err := time.ParseDuration(h.cfg.JWT.AccessTokenTTL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "server configuration error",
		})
		return
	}

	response := dto.AppRefreshResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    int(accessTTL.Seconds()),
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) appValidateToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "authorization header required",
		})
		return
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid authorization header format",
		})
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "token is empty",
		})
		return
	}

	user, err := h.services.AppUserService.ValidateAccessToken(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"valid": false,
			"error": "invalid or expired token",
		})
		return
	}

	response := dto.AppValidateResponse{
		Valid: true,
		User: dto.AppUserInfo{
			ID:            user.ID,
			Username:      user.Username,
			Role:          user.Role,
			AcademicGroup: user.AcademicGroup,
			Profile:       user.Profile,
		},
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) appUserInfo(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "authorization header required",
		})
		return
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid authorization header format",
		})
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "token is empty",
		})
		return
	}

	user, err := h.services.AppUserService.ValidateAccessToken(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid or expired token",
		})
		return
	}

	userInfo := dto.AppUserInfo{
		ID:            user.ID,
		Username:      user.Username,
		Role:          user.Role,
		AcademicGroup: user.AcademicGroup,
		Profile:       user.Profile,
	}

	c.JSON(http.StatusOK, gin.H{
		"user": userInfo,
	})
}
