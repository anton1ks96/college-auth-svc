package v1

import "github.com/gin-gonic/gin"

func InitAPIRoutes(api *gin.RouterGroup) {
	api = api.Group("/api")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
	}
}
