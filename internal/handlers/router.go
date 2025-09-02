package handlers

import (
	"github.com/anton1ks96/college-auth-svc/internal/handlers/v1"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.Use(gin.Recovery())
	InitAPI(r)
	return r
}

func InitAPI(r *gin.Engine) {
	v1.InitAPIRoutes(r.Group("/"))
}
