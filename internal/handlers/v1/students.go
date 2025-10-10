package v1

import (
	"net/http"

	"github.com/anton1ks96/college-auth-svc/internal/handlers/dto"
	"github.com/gin-gonic/gin"
)

func (h *Handler) searchStudents(c *gin.Context) {
	internalToken := c.GetHeader("X-Internal-Token")
	if internalToken == "" || internalToken != h.cfg.Tokens.InternalToken {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	var req dto.StudentSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	students, err := h.services.StudentService.SearchStudents(
		c.Request.Context(),
		req.Query,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"students": students,
		"total":    len(students),
	})
}

func (h *Handler) searchTeachers(c *gin.Context) {
	internalToken := c.GetHeader("X-Internal-Token")
	if internalToken == "" || internalToken != h.cfg.Tokens.InternalToken {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	var req dto.StudentSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	teachers, err := h.services.StudentService.SearchTeachers(
		c.Request.Context(),
		req.Query,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"students": teachers,
		"total":    len(teachers),
	})
}
