package handlers

import (
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.Use(gin.Recovery())
	InitAPI(r)
	return r
}

func InitAPI(r *gin.Engine) {
	initAPIRoutes(r.Group("/"))
}
