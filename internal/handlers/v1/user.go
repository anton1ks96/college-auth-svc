package v1

import "github.com/gin-gonic/gin"

//POST /api/v1/auth/signin
//GET /api/v1/auth/me
//POST /api/v1/auth/signout
//POST /api/v1/auth/refresh
//GET /api/v1/ping

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
