package v1

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/anton1ks96/college-auth-svc/internal/handlers/dto"
	service "github.com/anton1ks96/college-auth-svc/internal/services"
	"github.com/gin-gonic/gin"
)

func (h *Handler) signIn(c *gin.Context) {
	var loginReq dto.LoginRequest
	if err := c.BindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	tokens, user, err := h.services.UserService.SignIn(c.Request.Context(), service.SignInInput{
		UserID:   loginReq.Username,
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
			"error": "invalid access token TTL",
		})
		return
	}

	refreshTTL, err := time.ParseDuration(h.cfg.JWT.RefreshTokenTTL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid refresh token TTL",
		})
		return
	}

	c.SetCookie(
		"access_token",
		tokens.AccessToken,
		int(accessTTL.Seconds()),
		"/",
		"",
		false,
		true,
	)

	c.SetCookie(
		"refresh_token",
		tokens.RefreshToken,
		int(refreshTTL.Seconds()),
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
		"expires_in":    int(accessTTL.Seconds()),
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}

func (h *Handler) signOut(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "refresh token cookie not found",
		})
		return
	}

	if err := h.services.UserService.SignOut(c.Request.Context(), refreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to revoke session",
		})
		return
	}

	c.SetCookie("refresh_token", "", -1, "/", "", false, true)
	c.SetCookie("access_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "successfully signed out",
	})
}

func (h *Handler) refreshTokens(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "refresh token not found in cookies",
		})
		return
	}

	tokens, err := h.services.UserService.RefreshTokens(c.Request.Context(), refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "failed to refresh tokens",
		})
		return
	}

	accessTTL, err := time.ParseDuration(h.cfg.JWT.AccessTokenTTL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid access token TTL configuration",
		})
		return
	}

	c.SetCookie(
		"access_token",
		tokens.AccessToken,
		int(accessTTL.Seconds()),
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"access_token": tokens.AccessToken,
		"expires_in":   int(accessTTL.Seconds()),
	})
}

func (h *Handler) validateToken(c *gin.Context) {
	token, err := h.getFromHeader(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "auth header not found",
		})
		return
	}

	user, err := h.services.UserService.ValidateAccessToken(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to validate access token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}

func (h *Handler) getFromHeader(c *gin.Context) (string, error) {
	token := c.GetHeader("Authorization")
	if token != "" {
		if !strings.HasPrefix(token, "Bearer ") {
			return "", fmt.Errorf("invalid authorization header format")
		}
		extractedToken := strings.TrimPrefix(token, "Bearer ")
		if extractedToken == "" {
			return "", fmt.Errorf("empty token in authorization header")
		}
		return extractedToken, nil
	}

	cookieToken, err := c.Cookie("access_token")
	if err == nil && cookieToken != "" {
		return cookieToken, nil
	}

	return "", fmt.Errorf("authorization header is missing")
}
