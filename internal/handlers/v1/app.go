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

	refreshTTL, err := time.ParseDuration(h.cfg.JWT.RefreshTokenTTL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "server configuration error",
		})
	}

	response := dto.AppSignInResponse{
		AccessToken:      tokens.AccessToken,
		RefreshToken:     tokens.RefreshToken,
		AccessExpiresIn:  int(accessTTL.Seconds()),
		RefreshExpiresIn: int(refreshTTL.Seconds()),
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

func (h *Handler) appGetAccess(c *gin.Context) {
	var req dto.AppGetAccessRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request format",
			"details": err.Error(),
		})
		return
	}

	accessToken, err := h.services.AppUserService.GetAccessToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "failed to get access token",
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

	response := dto.AppGetAccessTokenResponse{
		AccessToken: accessToken,
		ExpiresIn:   int(accessTTL.Seconds()),
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) appRefreshToken(c *gin.Context) {
	var req dto.AppRefreshRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request format",
			"details": err.Error(),
		})
		return
	}

	newRefreshToken, err := h.services.AppUserService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "failed to refresh token",
			"details": err.Error(),
		})
		return
	}

	refreshTTL, err := time.ParseDuration(h.cfg.JWT.RefreshTokenTTL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "server configuration error",
		})
		return
	}

	response := dto.AppRefreshResponse{
		RefreshToken: newRefreshToken,
		ExpiresIn:    int(refreshTTL.Seconds()),
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
