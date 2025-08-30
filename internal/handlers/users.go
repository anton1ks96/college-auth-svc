package handlers

import "github.com/gin-gonic/gin"

func initAPIRoutes(api *gin.RouterGroup) {
	api = api.Group("/api")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
	}
}
